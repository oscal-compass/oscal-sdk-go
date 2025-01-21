/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"fmt"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
)

// ImplementationSettings defines settings for RuleSets defined at the control
// implementation and control level.
type ImplementationSettings struct {
	// implementedReqSettings defines settings for RuleSets at the
	// implemented requirement/individual control level.
	implementedReqSettings map[string]Settings
	// controlsByRules stores controlsIDs that have specific
	// rules mapped.
	controlsByRules map[string]set.Set[string]
	// settings defines the settings for the
	// overall implementation of the requirements.
	settings Settings
}

// NewImplementationSettings returned ImplementationSettings for an OSCAL ControlImplementationSet.
func NewImplementationSettings(controlImplementation oscalTypes.ControlImplementationSet) *ImplementationSettings {
	implementation := &ImplementationSettings{
		implementedReqSettings: make(map[string]Settings),
		settings: Settings{
			mappedRules:        map[string]struct{}{},
			selectedParameters: make(map[string]string),
		},
		controlsByRules: make(map[string]set.Set[string]),
	}
	if controlImplementation.SetParameters != nil {
		setParameters(*controlImplementation.SetParameters, implementation.settings.selectedParameters)
	}

	for _, implementedReq := range controlImplementation.ImplementedRequirements {
		requirement := NewSettingsFromImplementedRequirement(implementedReq)
		// Do not add requirements without mapped rules
		if len(requirement.mappedRules) == 0 {
			continue
		}
		for mappedRule := range requirement.mappedRules {
			controlSet, ok := implementation.controlsByRules[mappedRule]
			if !ok {
				controlSet = set.New[string]()
			}
			controlSet.Add(implementedReq.ControlId)
			implementation.controlsByRules[mappedRule] = controlSet
			implementation.settings.mappedRules.Add(mappedRule)
		}
		implementation.implementedReqSettings[implementedReq.ControlId] = requirement
	}

	return implementation
}

// AllSettings returns all settings collected for the overall control implementation.
func (i *ImplementationSettings) AllSettings() Settings {
	return i.settings
}

// AllControls returns a list of control ids found in the control implementation.
func (i *ImplementationSettings) AllControls() []string {
	var allControls []string
	for controlId := range i.implementedReqSettings {
		allControls = append(allControls, controlId)
	}
	return allControls
}

// ByControlID returns the individual requirement settings for a given control id in the
// control implementation.
func (i *ImplementationSettings) ByControlID(controlId string) (Settings, error) {
	requirement, ok := i.implementedReqSettings[controlId]
	if !ok {
		return Settings{}, fmt.Errorf("control %s not found in settings", controlId)
	}
	return requirement, nil
}

// ApplicableControls finds controls that are applicable to a given rule.
func (i *ImplementationSettings) ApplicableControls(ruleId string) ([]string, error) {
	controls, ok := i.controlsByRules[ruleId]
	if !ok {
		return nil, fmt.Errorf("rule id %s not found in settings", ruleId)
	}
	var controlsList []string
	for control := range controls {
		controlsList = append(controlsList, control)
	}
	return controlsList, nil
}

// merge another ControlImplementationSet into the ImplementationSettings. This also merged existing
// settings at the requirements level.
func (i *ImplementationSettings) merge(inputImplementation oscalTypes.ControlImplementationSet) {
	if inputImplementation.SetParameters != nil {
		setParameters(*inputImplementation.SetParameters, i.settings.selectedParameters)
	}

	for _, implementedReq := range inputImplementation.ImplementedRequirements {
		requirement, ok := i.implementedReqSettings[implementedReq.ControlId]
		if !ok {
			requirement = NewSettingsFromImplementedRequirement(implementedReq)
			// Do not add requirements without mapped rules
			if len(requirement.mappedRules) == 0 {
				continue
			}
			for mappedRule := range requirement.mappedRules {
				controlSet, ok := i.controlsByRules[mappedRule]
				if !ok {
					controlSet = set.New[string]()
				}
				controlSet.Add(implementedReq.ControlId)
				i.controlsByRules[mappedRule] = controlSet
				i.settings.mappedRules.Add(mappedRule)
			}
		} else {
			if implementedReq.Props != nil {
				mappedRulesProps := extensions.FindAllProps(extensions.RuleIdProp, *implementedReq.Props)
				for _, mappedRule := range mappedRulesProps {
					controlSet, ok := i.controlsByRules[mappedRule.Value]
					if !ok {
						controlSet = set.New[string]()
					}
					controlSet.Add(implementedReq.ControlId)
					i.controlsByRules[mappedRule.Value] = controlSet
					i.settings.mappedRules.Add(mappedRule.Value)
					requirement.mappedRules.Add(mappedRule.Value)
				}
			}

			if implementedReq.SetParameters != nil {
				setParameters(*implementedReq.SetParameters, requirement.selectedParameters)
			}
		}

		i.implementedReqSettings[implementedReq.ControlId] = requirement
	}
}
