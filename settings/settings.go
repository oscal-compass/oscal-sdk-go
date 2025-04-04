/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"context"
	"errors"
	"fmt"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
	"github.com/oscal-compass/oscal-sdk-go/rules"
)

// ErrRulesNotFound defines an error returned when there are not intersecting ruleSet store
// for a component and in the given Settings.
var ErrRulesNotFound = errors.New("no rules found with criteria")

// Settings defines settings for RuleSets to tune options based in the
// target baseline or compliance goals.
type Settings struct {
	// mappedRules is a list of rule IDs that are mapped to this requirement.
	mappedRules set.Set[string]
	// selectedParameters is a map of parameter names and their selected values for this requirement.
	selectedParameters map[string]string
}

// ApplyParameterSettings returns the given rule set with update parameter values based on the implementation.
//
// If the implementation does have parameter values or the rule set does not have a parameter, the original rule set
// is returned.
// The parameter value is not altered on the original rule set, it is copied and returned with the new rule set.
func (i Settings) ApplyParameterSettings(set extensions.RuleSet) extensions.RuleSet {
	if len(i.selectedParameters) > 0 && len(set.Rule.Parameters) > 0 {
		sliceCopy := make([]extensions.Parameter, len(set.Rule.Parameters))
		copy(sliceCopy, set.Rule.Parameters)
		for idx := range sliceCopy {
			selectedValue, ok := i.selectedParameters[sliceCopy[idx].ID]
			if ok {
				sliceCopy[idx].Value = selectedValue
			}
		}
		set.Rule.Parameters = sliceCopy
	}
	return set
}

// ContainsRule returns whether the given rule id is defined in the Settings.
func (i Settings) ContainsRule(ruleId string) bool {
	return i.mappedRules.Has(ruleId)
}

// ApplyToComponent returns a list of RuleSets for a given component with options applied from the given Settings.
//
// Only the rules that overlap between the component and the mapped rules in the implementation are returned.
// Parameters will be applied as RuleSet selected parameter values.
func ApplyToComponent(ctx context.Context, componentId string, store rules.Store, settings Settings) ([]extensions.RuleSet, error) {
	var resolvedRules []extensions.RuleSet
	componentRuleSets, err := store.FindByComponent(ctx, componentId)
	if err != nil {
		return []extensions.RuleSet{}, err
	}

	for _, ruleSet := range componentRuleSets {
		if !settings.ContainsRule(ruleSet.Rule.ID) {
			continue
		}
		ruleSet = settings.ApplyParameterSettings(ruleSet)
		resolvedRules = append(resolvedRules, ruleSet)
	}
	if len(resolvedRules) == 0 {
		return []extensions.RuleSet{}, fmt.Errorf("component %s: %w", componentId, ErrRulesNotFound)
	}
	return resolvedRules, nil
}
