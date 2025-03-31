/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package rules

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/oscal-compass/oscal-sdk-go/models/components"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
)

var (
	// Store interface check
	_ Store = (*MemoryStore)(nil)

	// ErrRuleNotFound defines an error returned when rule queries fail.
	ErrRuleNotFound = errors.New("associated rule object not found")
	// ErrComponentsNotFound defines an error returned during MemoryStore creation when the input
	// is invalid.
	ErrComponentsNotFound = errors.New("no components not found")
)

// MemoryStore implements the Store interface using an in-memory map-based data structure.
// WARNING: This implementation is not thread safe.
type MemoryStore struct {
	// nodes saves the rule ID map keys, which are used with
	// the other fields.
	nodes map[string]extensions.RuleSet
	// ByCheck store a mapping between the checkId and its parent
	// ruleId
	byCheck map[string]string

	// Below contains maps that store information by component and
	// component types to form RuleSet with the correct context.

	// rulesByComponent stores the component title of any component
	// mapped to any relevant rules.
	rulesByComponent map[string]set.Set[string]
	// checksByValidationComponent store checkId mapped to validation
	// component title to filter check information on rules.
	checksByValidationComponent map[string]set.Set[string]
}

// NewMemoryStore creates a new memory-based Store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		nodes:                       make(map[string]extensions.RuleSet),
		byCheck:                     make(map[string]string),
		rulesByComponent:            make(map[string]set.Set[string]),
		checksByValidationComponent: make(map[string]set.Set[string]),
	}
}

// IndexAll indexes rule information from OSCAL Components.
func (m *MemoryStore) IndexAll(components []components.Component) error {
	if len(components) == 0 {
		return fmt.Errorf("failed to index components: %w", ErrComponentsNotFound)
	}
	for _, component := range components {

		// Catalog information here at the component in the MemoryStore at the
		// component level

		componentTitle := component.Title()
		extractedRules, extractedChecks := m.indexComponent(component)
		if len(extractedRules) != 0 {
			existingRules, ok := m.rulesByComponent[componentTitle]
			if ok {
				for rule := range existingRules {
					extractedRules.Add(rule)
				}
			}
			m.rulesByComponent[component.Title()] = extractedRules
		}
		existingRules, ok := m.rulesByComponent[componentTitle]
		if ok {
			for rule := range existingRules {
				extractedRules.Add(rule)
			}
		}

		if len(extractedChecks) != 0 {
			existingChecks, ok := m.checksByValidationComponent[componentTitle]
			if ok {
				for check := range existingChecks {
					extractedChecks.Add(check)
				}
			}
			m.checksByValidationComponent[componentTitle] = extractedChecks
		}
	}
	return nil
}

// indexComponent returns extracted rules and checks (respectively) from a given component.
func (m *MemoryStore) indexComponent(component components.Component) (set.Set[string], set.Set[string]) {
	// Catalog all registered rules for all components and check implementations by validation component for filtering in
	// `rules.FindByComponent`.
	rules := set.New[string]()
	checks := set.New[string]()

	if len(component.Props()) == 0 {
		return rules, checks
	}

	// Each rule set is linked by a group id in the property remarks
	byRemarks := groupPropsByRemarks(component.Props())
	for _, propSet := range byRemarks {
		ruleIdProp, ok := getProp(extensions.RuleIdProp, propSet)
		if !ok {
			continue
		}

		ruleSet, ok := m.nodes[ruleIdProp.Value]
		if !ok {
			ruleSet = extensions.RuleSet{}
		}

		// A check may or may not be registered depending on
		// component.Type().
		placeholderCheck := extensions.Check{}

		// paramMap stores a map of parameters to their suffix value
		// of the property name.  Multiple properties contain
		// values needed to populate each parameter.  As the loop
		// iterates over the properties the parameters stored
		// in the map are populated.
		paramMap := make(map[string]extensions.Parameter)

		for prop := range propSet {
			var propName string
			var propSuffix string

			// Matches any parameter properties that have a numerical suffix
			re := regexp.MustCompile(`^Parameter_.*\d+$`)
			// If the property name contains "Parameter" and has a numerical
			// suffix then extract the suffix as the key for the map.
			// Otherwise default to using "0" as the key (meaning there is
			// no numerical suffix and only one parameter contained
			// in the properties).
			if re.MatchString(prop.Name) {
				// Split the property name to handle properties that have
				// a numerical suffix.  e.g Parameter_Id_1
				propNameParts := strings.Split(prop.Name, "_")
				propPrefix := propNameParts[:len(propNameParts)-1]
				propName = strings.Join(propPrefix, "_")
				propSuffix = propNameParts[len(propNameParts)-1]
			} else {
				propName = prop.Name
				propSuffix = "0" // Used if property does not have numerical suffix
			}

			switch propName {
			case extensions.RuleIdProp:
				ruleSet.Rule.ID = prop.Value
			case extensions.RuleDescriptionProp:
				ruleSet.Rule.Description = prop.Value
			case extensions.CheckIdProp:
				placeholderCheck.ID = prop.Value
			case extensions.CheckDescriptionProp:
				placeholderCheck.Description = prop.Value
			case extensions.ParameterIdProp:
				p, ok := paramMap[propSuffix]
				if !ok {
					paramMap[propSuffix] = extensions.Parameter{
						ID: prop.Value,
					}
				} else {
					p.ID = prop.Value
					paramMap[propSuffix] = p
				}
			case extensions.ParameterDescriptionProp:
				p, ok := paramMap[propSuffix]
				if !ok {
					paramMap[propSuffix] = extensions.Parameter{
						Description: prop.Value,
					}
				} else {
					p.Description = prop.Value
					paramMap[propSuffix] = p
				}
			case extensions.ParameterDefaultProp:
				p, ok := paramMap[propSuffix]

				if !ok {
					paramMap[propSuffix] = extensions.Parameter{
						Value: prop.Value,
					}
				} else {
					p.Value = prop.Value
					paramMap[propSuffix] = p
				}
			}
		}

		// Add any parameters that were extracted from the
		// properties to the rule
		if len(paramMap) > 0 {
			ruleParams := make([]extensions.Parameter, 0, len(paramMap))
			ruleSet.Rule.Parameters = ruleParams
			for _, param := range paramMap {
				ruleSet.Rule.Parameters = append(ruleSet.Rule.Parameters, param)
			}
		}

		if placeholderCheck.ID != "" {
			ruleSet.Checks = append(ruleSet.Checks, placeholderCheck)
			m.byCheck[placeholderCheck.ID] = ruleSet.Rule.ID
			checks.Add(placeholderCheck.ID)
		}
		rules.Add(ruleSet.Rule.ID)
		m.nodes[ruleSet.Rule.ID] = ruleSet
	}
	return rules, checks
}

func (m *MemoryStore) GetByRuleID(_ context.Context, ruleId string) (extensions.RuleSet, error) {
	ruleSet, ok := m.nodes[ruleId]
	if !ok {
		return extensions.RuleSet{}, fmt.Errorf("rule %q: %w", ruleId, ErrRuleNotFound)
	}
	return ruleSet, nil
}

func (m *MemoryStore) GetByCheckID(ctx context.Context, checkId string) (extensions.RuleSet, error) {
	ruleId, ok := m.byCheck[checkId]
	if !ok {
		return extensions.RuleSet{}, fmt.Errorf("failed to find rule for check %q: %w", checkId, ErrRuleNotFound)
	}
	return m.GetByRuleID(ctx, ruleId)
}

func (m *MemoryStore) FindByComponent(ctx context.Context, componentId string) ([]extensions.RuleSet, error) {
	ruleIds, ok := m.rulesByComponent[componentId]
	if !ok {
		return nil, fmt.Errorf("failed to find rules for component %q", componentId)
	}

	var ruleSets []extensions.RuleSet
	var errs []error
	for ruleId := range ruleIds {
		ruleSet, err := m.GetByRuleID(ctx, ruleId)
		if err != nil {
			errs = append(errs, err)
		}

		// Make sure we are only returning the relevant checks for this
		// component.
		if checkIds, ok := m.checksByValidationComponent[componentId]; ok {
			filteredChecks := make([]extensions.Check, 0, len(ruleSet.Checks))
			for _, check := range ruleSet.Checks {
				if checkIds.Has(check.ID) {
					filteredChecks = append(filteredChecks, check)
				}
			}
			ruleSet.Checks = filteredChecks
		}

		ruleSets = append(ruleSets, ruleSet)
	}

	if len(errs) > 0 {
		joinedErr := errors.Join(errs...)
		return ruleSets, fmt.Errorf("failed to find rules for component %q: %w", componentId, joinedErr)
	}

	return ruleSets, nil
}
