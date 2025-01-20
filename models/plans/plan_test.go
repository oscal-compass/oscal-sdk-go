/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package plans

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/generators"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/settings"
)

func TestGenerateAssessmentPlan(t *testing.T) {
	testComp := readCompDef(t)
	defaultComponents := prepComponents(t, testComp)
	defaultSettings := prepSettings(t, testComp)

	tests := []struct {
		name            string
		inputComponents []components.Component
		inputSetting    *settings.ImplementationSettings
		inputOptions    []GenerateOption
		assertFunc      func(*testing.T, *oscalTypes.AssessmentPlan)
		expError        string
	}{
		{
			name:            "Success/Defaults",
			inputComponents: defaultComponents,
			inputSetting:    defaultSettings,
			inputOptions:    nil,
			assertFunc: func(t *testing.T, plan *oscalTypes.AssessmentPlan) {
				require.Len(t, *plan.LocalDefinitions.Activities, 4)
				require.Len(t, *plan.AssessmentAssets.Components, 2)
				require.Len(t, *plan.AssessmentSubjects, 1)
				require.Len(t, plan.ReviewedControls.ControlSelections, 1)
			},
			expError: "",
		},
		{
			name:            "Success/WithOptions",
			inputComponents: defaultComponents,
			inputSetting:    defaultSettings,
			inputOptions:    []GenerateOption{WithTitle("mytitle"), WithImport("myimport")},
			assertFunc: func(t *testing.T, plan *oscalTypes.AssessmentPlan) {
				require.Equal(t, plan.Metadata.Title, "mytitle")
				require.Equal(t, plan.ImportSsp.Href, "myimport")
			},
			expError: "",
		},
		{
			name:            "Failure/NoComponents",
			inputComponents: nil,
			inputSetting:    defaultSettings,
			inputOptions:    nil,
			expError:        "failed processing components for assessment plan \"REPLACE_ME\": failed to index components: no components not found",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.TODO()
			assessmentPlan, err := GenerateAssessmentPlan(ctx, c.inputComponents, c.inputSetting, c.inputOptions...)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				c.assertFunc(t, assessmentPlan)
			}
		})
	}
}

func readCompDef(t *testing.T) oscalTypes.ComponentDefinition {
	testDataPath := filepath.Join("../../testdata", "component-definition-test.json")

	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	definition, err := generators.NewComponentDefinition(file)
	require.NoError(t, err)
	require.NotNil(t, definition)
	return *definition
}

func prepComponents(t *testing.T, definition oscalTypes.ComponentDefinition) []components.Component {
	require.NotNil(t, definition.Components)
	var comps []components.Component
	for _, cp := range *definition.Components {
		adapters := components.NewDefinedComponentAdapter(cp)
		comps = append(comps, adapters)
	}
	return comps
}

func prepSettings(t *testing.T, definition oscalTypes.ComponentDefinition) *settings.ImplementationSettings {
	var allImplementations []oscalTypes.ControlImplementationSet
	for _, component := range *definition.Components {
		if component.ControlImplementations == nil {
			continue
		}
		allImplementations = append(allImplementations, *component.ControlImplementations...)
	}
	impSettings, err := settings.Framework("cis", allImplementations)
	require.NoError(t, err)
	return impSettings
}
