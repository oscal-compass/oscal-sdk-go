/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"fmt"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/oscal-compass/oscal-sdk-go/internal/set"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
)

// ImplementationSettings defines settings for RuleSets defined at the control
// implementation and control level.
type ImplementationSettings struct {
	// implementedReqSettings defines settings for RuleSets at the
	// implemented requirement/statement level.
	implementedReqSettings map[string]Settings
	// settings defines the settings for the
	// overall implementation of the requirements.
	settings Settings

	// Below are mapping with information about implemented
	// controls indexed

	// controls saves the control ID map keys, which are used with
	// the other fields to retrieve information about controls.
	controlsById map[string]oscalTypes.AssessedControlsSelectControlById
	// controlsByRules stores controlsIDs that have specific
	// rules mapped.
	controlsByRules map[string]set.Set[string]
}

// AllSettings returns all settings collected for the overall control implementation.
func (i *ImplementationSettings) AllSettings() Settings {
	return i.settings
}

// AllControls returns AssessedControlsSelectControlByID with all controls and associated statements that are applicable to
// the control implementation.
func (i *ImplementationSettings) AllControls() []oscalTypes.AssessedControlsSelectControlById {
	var allControls []oscalTypes.AssessedControlsSelectControlById
	for _, assessedControls := range i.controlsById {
		allControls = append(allControls, assessedControls)
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

// ApplicableControls finds controls and corresponding statements that are applicable to a given rule based in the control
// implementation.
func (i *ImplementationSettings) ApplicableControls(ruleId string) ([]oscalTypes.AssessedControlsSelectControlById, error) {
	controls, ok := i.controlsByRules[ruleId]
	if !ok {
		return nil, fmt.Errorf("rule id %s not found in settings", ruleId)
	}
	var assessedControls []oscalTypes.AssessedControlsSelectControlById
	for control := range controls {
		assessedControl, ok := i.controlsById[control]
		if !ok {
			return nil, fmt.Errorf("assessed control object %s not found for rule %s", control, ruleId)
		}
		assessedControls = append(assessedControls, assessedControl)
	}
	return assessedControls, nil
}

// merge another ImplementationSettings into the ImplementationSettings. Existing settings at the
// requirements level are also merged.
func (i *ImplementationSettings) merge(inputImplementation components.Implementation) {
	setParameters(inputImplementation.SetParameters(), i.settings.selectedParameters)

	for _, requirement := range inputImplementation.Requirements() {
		reqSettings, ok := i.implementedReqSettings[requirement.ControlID()]
		if !ok {
			newRequirementForImplementation(requirement, i)
		} else {

			inputRequirement := settingsFromImplementedRequirement(requirement)
			if len(inputRequirement.mappedRules) == 0 {
				continue
			}

			for mappedRule := range inputRequirement.mappedRules {
				controlSet, ok := i.controlsByRules[mappedRule]
				if !ok {
					controlSet = set.New[string]()
				}
				controlSet.Add(requirement.ControlID())
				i.controlsByRules[mappedRule] = controlSet
				i.settings.mappedRules.Add(mappedRule)
				reqSettings.mappedRules.Add(mappedRule)
			}
			for name, value := range inputRequirement.selectedParameters {
				reqSettings.selectedParameters[name] = value
			}
			i.implementedReqSettings[requirement.ControlID()] = reqSettings
		}
	}
}
