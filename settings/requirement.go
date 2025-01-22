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

var _ Settings = (*RequirementSettings)(nil)

// RequirementSettings defines settings for RuleSets defined by requirements.
type RequirementSettings struct {
	// MappedRules is a list of rule IDs that are mapped to this requirement.
	mappedRules set.Set[string]
	// SelectedParameters is a map of parameter names and their selected values for this requirement.
	selectedParameters map[string]string
}

// NewSettingsFromImplementedRequirement returns an initialized RequirementSettings from an
// OSCAL Implemented Requirement.
func NewSettingsFromImplementedRequirement(implementedReq oscalTypes.ImplementedRequirementControlImplementation) RequirementSettings {
	requirement := RequirementSettings{
		selectedParameters: make(map[string]string),
		mappedRules:        set.New[string](),
	}

	if implementedReq.Props != nil {
		mappedRulesProps := extensions.FindAllProps(extensions.RuleIdProp, *implementedReq.Props)
		for _, mappedRule := range mappedRulesProps {
			requirement.mappedRules.Add(mappedRule.Value)
		}
	}

	if implementedReq.SetParameters != nil {
		setParameters(*implementedReq.SetParameters, requirement.selectedParameters)
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

func (i RequirementSettings) ApplyParameterSettings(set extensions.RuleSet) extensions.RuleSet {
	if len(i.selectedParameters) > 0 && set.Rule.Parameter != nil {
		selectedValue, ok := i.selectedParameters[set.Rule.Parameter.ID]
		if ok {
			parameterCopy := *set.Rule.Parameter
			parameterCopy.Value = selectedValue
			set.Rule.Parameter = &parameterCopy
		}
	}
	return set
}

func (i RequirementSettings) ContainsRule(ruleId string) bool {
	return i.mappedRules.Has(ruleId)
}
