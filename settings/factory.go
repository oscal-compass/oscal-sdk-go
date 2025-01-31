/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
)

// NewSettings returns a new Settings instance with given rules and associated rule parameters.
func NewSettings(rules map[string]struct{}, parameters map[string]string) Settings {
	return Settings{
		selectedParameters: parameters,
		mappedRules:        rules,
	}
}

// NewImplementationSettings returns ImplementationSettings populated with data from an OSCAL Control Implementation
// Set and the nested Implemented Requirements.
func NewImplementationSettings(controlImplementation oscalTypes.ControlImplementationSet) *ImplementationSettings {
	implementation := &ImplementationSettings{
		implementedReqSettings: make(map[string]Settings),
		settings:               NewSettings(set.New[string](), make(map[string]string)),
		controlsByRules:        make(map[string]set.Set[string]),
		controlsById:           make(map[string]oscalTypes.AssessedControlsSelectControlById),
	}
	if controlImplementation.SetParameters != nil {
		setParameters(*controlImplementation.SetParameters, implementation.settings.selectedParameters)
	}

	for _, implementedReq := range controlImplementation.ImplementedRequirements {
		newRequirementForImplementation(implementedReq, implementation)
	}

	return implementation
}

// NewAssessmentActivitiesSettings returns a new Setting populate based on data from OSCAL Activities
//
// The mapping between a RuleSet and Activity is as follows:
// Activity -> Rule
// Title -> Rule ID
// Parameter -> Activity Property
func NewAssessmentActivitiesSettings(assessmentActivities []oscalTypes.Activity) Settings {
	rules := set.New[string]()
	parameters := make(map[string]string)
	for _, activity := range assessmentActivities {

		// Activities based on rules are expected to have at
		// least one property set
		if activity.Props == nil {
			continue
		}

		paramProps := extensions.FindAllProps(*activity.Props, extensions.WithClass(extensions.TestParameterClass))
		for _, param := range paramProps {
			parameters[param.Name] = param.Value
		}

		rules.Add(activity.Title)
	}
	return Settings{
		mappedRules:        rules,
		selectedParameters: parameters,
	}
}

//	newRequirementForImplementation adds a new Setting to an existing ImplementationSettings and updates all related
//
// fields.
func newRequirementForImplementation(implementedReq oscalTypes.ImplementedRequirementControlImplementation, implementation *ImplementationSettings) {
	implementedControl := oscalTypes.AssessedControlsSelectControlById{
		ControlId: implementedReq.ControlId,
	}
	requirement := settingsFromImplementedRequirement(implementedReq)

	// Do not add requirements without mapped rules
	if len(requirement.mappedRules) > 0 {
		for mappedRule := range requirement.mappedRules {
			controlSet, ok := implementation.controlsByRules[mappedRule]
			if !ok {
				controlSet = set.New[string]()
			}
			controlSet.Add(implementedReq.ControlId)
			implementation.controlsByRules[mappedRule] = controlSet
			implementation.controlsById[implementedReq.ControlId] = implementedControl
			implementation.settings.mappedRules.Add(mappedRule)
		}

		implementation.implementedReqSettings[implementedReq.ControlId] = requirement
	}
}

// settingsFromImplementedRequirement returns Settings populated with data from an
// OSCAL Implemented Requirement.
func settingsFromImplementedRequirement(implementedReq oscalTypes.ImplementedRequirementControlImplementation) Settings {
	requirement := NewSettings(set.New[string](), make(map[string]string))

	if implementedReq.Props != nil {
		mappedRulesProps := extensions.FindAllProps(*implementedReq.Props, extensions.WithName(extensions.RuleIdProp))
		for _, mappedRule := range mappedRulesProps {
			requirement.mappedRules.Add(mappedRule.Value)
		}
	}

	if implementedReq.SetParameters != nil {
		setParameters(*implementedReq.SetParameters, requirement.selectedParameters)
	}

	if implementedReq.Statements != nil {
		for _, stm := range *implementedReq.Statements {
			if stm.Props != nil {
				mappedRulesProps := extensions.FindAllProps(*stm.Props, extensions.WithName(extensions.RuleIdProp))
				if len(mappedRulesProps) == 0 {
					continue
				}
				for _, mappedRule := range mappedRulesProps {
					requirement.mappedRules.Add(mappedRule.Value)
				}
			}

		}
	}

	return requirement
}

// setParameters updates the paramMap with the input list of SetParameters.
func setParameters(parameters []oscalTypes.SetParameter, paramMap map[string]string) {
	for _, prm := range parameters {
		// Parameter values set for trestle Rule selection
		// should only map to a single value.
		if len(prm.Values) != 1 {
			continue
		}
		paramMap[prm.ParamId] = prm.Values[0]
	}
}
