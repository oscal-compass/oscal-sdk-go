/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
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
func NewImplementationSettings(controlImplementation components.Implementation) *ImplementationSettings {
	implementation := &ImplementationSettings{
		implementedReqSettings: make(map[string]Settings),
		settings:               NewSettings(set.New[string](), make(map[string]string)),
		controlsByRules:        make(map[string]set.Set[string]),
		controlsById:           make(map[string]oscalTypes.AssessedControlsSelectControlById),
	}
	setParameters(controlImplementation.SetParameters(), implementation.settings.selectedParameters)

	for _, requirement := range controlImplementation.Requirements() {
		newRequirementForImplementation(requirement, implementation)
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
		// Skipped activity has a property named 'skipped'
		skipped, found := extensions.GetTrestleProp(extensions.SkippedRulesProperty, *activity.Props)
		if found && skipped.Value == "true" {
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
func newRequirementForImplementation(implementedReq components.Requirement, implementation *ImplementationSettings) {
	implementedControl := oscalTypes.AssessedControlsSelectControlById{
		ControlId: implementedReq.ControlID(),
	}
	requirement := settingsFromImplementedRequirement(implementedReq)

	// Do not add requirements without mapped rules
	if len(requirement.mappedRules) > 0 {
		for mappedRule := range requirement.mappedRules {
			controlSet, ok := implementation.controlsByRules[mappedRule]
			if !ok {
				controlSet = set.New[string]()
			}
			controlSet.Add(implementedReq.ControlID())
			implementation.controlsByRules[mappedRule] = controlSet
			implementation.controlsById[implementedReq.ControlID()] = implementedControl
			implementation.settings.mappedRules.Add(mappedRule)
		}

		implementation.implementedReqSettings[implementedReq.ControlID()] = requirement
	}
}

// settingsFromImplementedRequirement returns Settings populated with data from an
// OSCAL Implemented Requirement.
func settingsFromImplementedRequirement(implementedReq components.Requirement) Settings {
	requirement := NewSettings(set.New[string](), make(map[string]string))

	mappedRulesProps := extensions.FindAllProps(implementedReq.Props(), extensions.WithName(extensions.RuleIdProp))
	for _, mappedRule := range mappedRulesProps {
		requirement.mappedRules.Add(mappedRule.Value)
	}

	setParameters(implementedReq.SetParameters(), requirement.selectedParameters)

	for _, stm := range implementedReq.Statements() {
		mappedRulesStmProps := extensions.FindAllProps(stm.Props(), extensions.WithName(extensions.RuleIdProp))
		if len(mappedRulesStmProps) == 0 {
			continue
		}
		for _, mappedRule := range mappedRulesStmProps {
			requirement.mappedRules.Add(mappedRule.Value)
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
