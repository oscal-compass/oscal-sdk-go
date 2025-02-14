/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"os"
	"path/filepath"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/generators"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name                string
		inputImplementation oscalTypes.ControlImplementationSet
		wantSettings        ImplementationSettings
	}{
		{
			name: "Valid/ImplementationOnly",
			inputImplementation: oscalTypes.ControlImplementationSet{
				SetParameters: &[]oscalTypes.SetParameter{
					{
						ParamId: "my-test-param",
						Values:  []string{"test-value"},
					},
				},
				ImplementedRequirements: []oscalTypes.ImplementedRequirementControlImplementation{
					{
						ControlId: "ex-1",
						Props: &[]oscalTypes.Property{
							{
								Name:  extensions.RuleIdProp,
								Value: "my-test-rule",
								Ns:    extensions.TrestleNameSpace,
							},
						},
						Statements: &[]oscalTypes.ControlStatementImplementation{
							{
								Props: &[]oscalTypes.Property{
									{
										Name:  extensions.RuleIdProp,
										Value: "my-test-rule-2",
										Ns:    extensions.TrestleNameSpace,
									},
								},
							},
						},
					},
				},
			},
			wantSettings: ImplementationSettings{
				settings: Settings{
					mappedRules: set.Set[string]{
						"etcd_cert_file": struct{}{},
						"etcd_key_file":  struct{}{},
						"my-test-rule":   struct{}{},
						"my-test-rule-2": struct{}{},
					},
					selectedParameters: map[string]string{
						"my-test-param": "test-value",
					},
				},
				implementedReqSettings: map[string]Settings{
					"CIS-2.1": {
						mappedRules: set.Set[string]{
							"etcd_cert_file": struct{}{},
							"etcd_key_file":  struct{}{},
						},
						selectedParameters: map[string]string{},
					},
					"ex-1": {
						mappedRules: set.Set[string]{
							"my-test-rule":   struct{}{},
							"my-test-rule-2": struct{}{},
						},
						selectedParameters: map[string]string{},
					},
				},
				controlsByRules: map[string]set.Set[string]{
					"etcd_cert_file": {
						"CIS-2.1": struct{}{},
					},
					"etcd_key_file": {
						"CIS-2.1": struct{}{},
					},
					"my-test-rule": {
						"ex-1": struct{}{},
					},
					"my-test-rule-2": {
						"ex-1": struct{}{},
					},
				},
				controlsById: map[string]oscalTypes.AssessedControlsSelectControlById{
					"CIS-2.1": {
						ControlId: "CIS-2.1",
					},
					"ex-1": {
						ControlId: "ex-1",
					},
				},
			},
		},
		{
			name: "Valid/ExistingControl",
			inputImplementation: oscalTypes.ControlImplementationSet{
				ImplementedRequirements: []oscalTypes.ImplementedRequirementControlImplementation{
					{
						ControlId: "CIS-2.1",
						SetParameters: &[]oscalTypes.SetParameter{
							{
								ParamId: "my-test-param",
								Values:  []string{"test-value"},
							},
						},
						Props: &[]oscalTypes.Property{
							{
								Name:  extensions.RuleIdProp,
								Value: "my-test-rule",
								Ns:    extensions.TrestleNameSpace,
							},
						},
					},
				},
			},
			wantSettings: ImplementationSettings{
				settings: Settings{
					mappedRules: set.Set[string]{
						"etcd_cert_file": struct{}{},
						"etcd_key_file":  struct{}{},
						"my-test-rule":   struct{}{},
					},
					selectedParameters: map[string]string{},
				},
				implementedReqSettings: map[string]Settings{
					"CIS-2.1": {
						mappedRules: set.Set[string]{
							"etcd_cert_file": struct{}{},
							"etcd_key_file":  struct{}{},
							"my-test-rule":   struct{}{},
						},
						selectedParameters: map[string]string{
							"my-test-param": "test-value",
						},
					},
				},
				controlsByRules: map[string]set.Set[string]{
					"etcd_cert_file": {
						"CIS-2.1": struct{}{},
					},
					"etcd_key_file": {
						"CIS-2.1": struct{}{},
					},
					"my-test-rule": {
						"CIS-2.1": struct{}{},
					},
				},
				controlsById: map[string]oscalTypes.AssessedControlsSelectControlById{
					"CIS-2.1": {
						ControlId: "CIS-2.1",
					},
				},
			},
		},
		{
			name: "Valid/ExistingRule",
			inputImplementation: oscalTypes.ControlImplementationSet{
				ImplementedRequirements: []oscalTypes.ImplementedRequirementControlImplementation{
					{
						ControlId: "ex-1",
						Props: &[]oscalTypes.Property{
							{
								Name:  extensions.RuleIdProp,
								Value: "etcd_cert_file",
								Ns:    extensions.TrestleNameSpace,
							},
						},
					},
				},
			},
			wantSettings: ImplementationSettings{
				settings: Settings{
					mappedRules: set.Set[string]{
						"etcd_cert_file": struct{}{},
						"etcd_key_file":  struct{}{},
					},
					selectedParameters: map[string]string{},
				},
				implementedReqSettings: map[string]Settings{
					"CIS-2.1": {
						mappedRules: set.Set[string]{
							"etcd_cert_file": struct{}{},
							"etcd_key_file":  struct{}{},
						},
						selectedParameters: map[string]string{},
					},
					"ex-1": {
						mappedRules: set.Set[string]{
							"etcd_cert_file": struct{}{},
						},
						selectedParameters: map[string]string{},
					},
				},
				controlsByRules: map[string]set.Set[string]{
					"etcd_cert_file": {
						"CIS-2.1": struct{}{},
						"ex-1":    struct{}{},
					},
					"etcd_key_file": {
						"CIS-2.1": struct{}{},
					},
				},
				controlsById: map[string]oscalTypes.AssessedControlsSelectControlById{
					"CIS-2.1": {
						ControlId: "CIS-2.1",
					},
					"ex-1": {
						ControlId: "ex-1",
					},
				},
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			testSettings := prepSettings(t)
			testSettings.merge(c.inputImplementation)
			require.Equal(t, c.wantSettings, *testSettings)
		})
	}
}

func TestImplementationSettings_Controls(t *testing.T) {
	testSettings := prepSettings(t)
	expectedControlIds := []oscalTypes.AssessedControlsSelectControlById{
		{
			ControlId: "CIS-2.1",
		},
	}
	gotControlsIds := testSettings.AllControls()
	require.Equal(t, expectedControlIds, gotControlsIds)

	gotControlIds, err := testSettings.ApplicableControls("etcd_cert_file")
	require.NoError(t, err)
	require.Equal(t, expectedControlIds, gotControlIds)
}

func prepSettings(t *testing.T) *ImplementationSettings {
	testDataPath := filepath.Join("../testdata", "component-definition-test-reqs.json")

	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	definition, err := generators.NewComponentDefinition(file)
	require.NoError(t, err)
	require.NotNil(t, definition)
	var allImplementations []oscalTypes.ControlImplementationSet
	for _, component := range *definition.Components {
		if component.ControlImplementations == nil {
			continue
		}
		allImplementations = append(allImplementations, *component.ControlImplementations...)
	}
	impSettings, err := Framework("cis", allImplementations)
	require.NoError(t, err)
	return impSettings
}
