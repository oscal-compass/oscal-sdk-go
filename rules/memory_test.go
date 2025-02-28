/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package rules

import (
	"context"
	"os"
	"testing"

	oscaltypes112 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/generators"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
)

var (
	expectedCertFileRule = extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "etcd_cert_file",
			Description: "Ensure that the --cert-file argument is set as appropriate",
		},
		Checks: []extensions.Check{
			{
				ID:          "etcd_cert_file",
				Description: "Check that the --cert-file argument is set as appropriate",
			},
		},
	}
	expectedKeyFileRule = extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "etcd_key_file",
			Description: "Ensure that the --key-file argument is set as appropriate",
			Parameter: &extensions.Parameter{
				ID:          "file_name",
				Description: "A parameter for a file name",
			},
		},
		Checks: []extensions.Check{
			{
				ID:          "etcd_key_file",
				Description: "Check that the --key-file argument is set as appropriate",
			},
		},
	}
)

func TestMemoryStore_IndexAll(t *testing.T) {
	tests := []struct {
		name         string
		testDataPath string
		expError     string
		wantNodes    map[string]extensions.RuleSet
	}{
		{
			name:         "Valid/WithRules",
			testDataPath: "../testdata/component-definition-test.json",
			wantNodes: map[string]extensions.RuleSet{
				"etcd_key_file": expectedKeyFileRule,

				"etcd_cert_file": expectedCertFileRule,
			},
		},
		{
			name:         "Valid/NoRules",
			testDataPath: "../testdata/component-definition-no-rules.json",
			wantNodes:    map[string]extensions.RuleSet{},
		},
		{
			name:         "Failure/NoComponents",
			testDataPath: "../testdata/component-definition-no-components.json",
			expError:     "failed to index components: no components not found",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			file, err := os.Open(c.testDataPath)
			require.NoError(t, err)
			definition, err := generators.NewComponentDefinition(file)
			require.NoError(t, err)

			if definition.Components == nil {
				definition.Components = &[]oscaltypes112.DefinedComponent{}
			}

			var comps []components.Component
			for _, cp := range *definition.Components {
				adapters := components.NewDefinedComponentAdapter(cp)
				comps = append(comps, adapters)
			}

			testMemory := NewMemoryStore()
			err = testMemory.IndexAll(comps)

			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.wantNodes, testMemory.nodes)
			}
		})
	}
}

func TestMemoryStore_GetByRuleID(t *testing.T) {
	testMemory := prepMemoryStore(t)
	testCtx := context.Background()

	found, err := testMemory.GetByRuleID(testCtx, "etcd_cert_file")
	require.NoError(t, err)

	expectedRule := extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "etcd_cert_file",
			Description: "Ensure that the --cert-file argument is set as appropriate",
		},
		Checks: []extensions.Check{
			{
				ID:          "etcd_cert_file",
				Description: "Check that the --cert-file argument is set as appropriate",
			},
		},
	}
	require.Equal(t, expectedRule, found)

	_, err = testMemory.GetByRuleID(testCtx, "not_present")
	require.EqualError(t, err, "rule \"not_present\": associated rule object not found")

}

func TestMemoryStore_GetByCheckID(t *testing.T) {
	testMemory := prepMemoryStore(t)
	testCtx := context.Background()

	found, err := testMemory.GetByCheckID(testCtx, "etcd_key_file")
	require.NoError(t, err)
	require.Equal(t, expectedKeyFileRule, found)

	_, err = testMemory.GetByCheckID(testCtx, "not_present")
	require.EqualError(t, err, "failed to find rule for check \"not_present\": associated rule object not found")

}

func TestMemoryStore_FindByComponent(t *testing.T) {
	testMemory := prepMemoryStore(t)
	testCtx := context.Background()

	softwareRuleSet, err := testMemory.FindByComponent(testCtx, "TestKubernetes")
	require.NoError(t, err)

	require.Contains(t, softwareRuleSet, expectedCertFileRule)
	require.Contains(t, softwareRuleSet, expectedKeyFileRule)

	validator1RuleSet, err := testMemory.FindByComponent(testCtx, "Validator")
	require.NoError(t, err)
	require.Len(t, validator1RuleSet, 1)
	require.Contains(t, validator1RuleSet, expectedKeyFileRule)

	validator2RuleSet, err := testMemory.FindByComponent(testCtx, "Validator2")
	require.NoError(t, err)
	require.Contains(t, validator2RuleSet, expectedCertFileRule)

	_, err = testMemory.FindByComponent(testCtx, "not_a_component")
	require.EqualError(t, err, "failed to find rules for component \"not_a_component\"")

	// Add a new target component for "Validator"
	testDataPath := "../testdata/component-definition-test2.json"
	loadComponents(t, testMemory, testDataPath)

	expectedExampleRule := extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "example_rule_1",
			Description: "Example rule 1 description",
		},
		Checks: []extensions.Check{
			{
				ID:          "example_check_1",
				Description: "Example check 1 description",
			},
		},
	}

	validator1RuleSet, err = testMemory.FindByComponent(testCtx, "Validator")
	require.NoError(t, err)
	require.Len(t, validator1RuleSet, 2)
	require.Contains(t, validator1RuleSet, expectedExampleRule, expectedKeyFileRule)
}

func prepMemoryStore(t *testing.T) *MemoryStore {
	testDataPath := "../testdata/component-definition-test.json"
	testMemory := NewMemoryStore()
	loadComponents(t, testMemory, testDataPath)
	return testMemory
}

func loadComponents(t *testing.T, store *MemoryStore, testDataPath string) {
	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	definition, err := generators.NewComponentDefinition(file)
	require.NoError(t, err)

	var comps []components.Component
	for _, cp := range *definition.Components {
		adapters := components.NewDefinedComponentAdapter(cp)
		comps = append(comps, adapters)
	}
	err = store.IndexAll(comps)
	require.NoError(t, err)
}
