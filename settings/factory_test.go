/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
)

func TestSettingsFromImplementedRequirements(t *testing.T) {
	tests := []struct {
		name             string
		inputRequirement oscalTypes.ImplementedRequirementControlImplementation
		wantSettings     Settings
	}{
		{
			name: "Valid/MappedRulesFound",
			inputRequirement: oscalTypes.ImplementedRequirementControlImplementation{
				Props: &[]oscalTypes.Property{
					{
						Name:  extensions.RuleIdProp,
						Ns:    extensions.TrestleNameSpace,
						Value: "rule-1",
					},
					{
						Name:  extensions.RuleIdProp,
						Ns:    extensions.TrestleNameSpace,
						Value: "rule-2",
					},
				},
			},
			wantSettings: Settings{
				mappedRules: set.Set[string]{
					"rule-1": struct{}{},
					"rule-2": struct{}{},
				},
				selectedParameters: map[string]string{},
			},
		},
		{
			name: "Valid/ParametersFound",
			inputRequirement: oscalTypes.ImplementedRequirementControlImplementation{
				Props: &[]oscalTypes.Property{
					{
						Name:  extensions.RuleIdProp,
						Ns:    extensions.TrestleNameSpace,
						Value: "rule-1",
					},
					{
						Name:  extensions.RuleIdProp,
						Ns:    extensions.TrestleNameSpace,
						Value: "rule-2",
					},
				},
				SetParameters: &[]oscalTypes.SetParameter{
					{
						ParamId: "param-1",
						Values: []string{
							"value",
						},
					},
				},
			},
			wantSettings: Settings{
				mappedRules: set.Set[string]{
					"rule-1": struct{}{},
					"rule-2": struct{}{},
				},
				selectedParameters: map[string]string{
					"param-1": "value",
				},
			},
		},
		{
			name:             "Valid/NoSettingsFound",
			inputRequirement: oscalTypes.ImplementedRequirementControlImplementation{},
			wantSettings: Settings{
				mappedRules:        map[string]struct{}{},
				selectedParameters: map[string]string{},
			},
		},
		{
			name: "Invalid/MultipleParametersValues",
			inputRequirement: oscalTypes.ImplementedRequirementControlImplementation{
				SetParameters: &[]oscalTypes.SetParameter{
					{
						ParamId: "param-1",
						Values: []string{
							"value-1",
							"value-2",
						},
					},
				},
			},
			wantSettings: Settings{
				mappedRules:        set.Set[string]{},
				selectedParameters: map[string]string{},
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			gotSettings := settingsFromImplementedRequirement(c.inputRequirement)
			require.Equal(t, c.wantSettings, gotSettings)
		})
	}
}
