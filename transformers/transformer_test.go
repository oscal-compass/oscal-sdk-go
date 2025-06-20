/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package transformers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/validation"
)

func TestComponentDefinitionsToAssessmentPlan(t *testing.T) {
	testDataPath := filepath.Join("../testdata", "component-definition-test.json")

	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	definition, err := models.NewComponentDefinition(file, validation.NoopValidator{})
	require.NoError(t, err)
	require.NotNil(t, definition)
	require.NotNil(t, definition.Components)

	plan, err := ComponentDefinitionsToAssessmentPlan(context.TODO(), []oscalTypes.ComponentDefinition{*definition}, "cis")
	require.NoError(t, err)

	require.Len(t, *plan.LocalDefinitions.Activities, 2)
	require.Len(t, *plan.AssessmentAssets.Components, 2)
	require.Len(t, *plan.AssessmentSubjects, 1)
	require.Len(t, plan.ReviewedControls.ControlSelections, 1)
	require.Len(t, *plan.Tasks, 1)
	tasks := *plan.Tasks
	require.Len(t, *tasks[0].AssociatedActivities, 2)

	var activities []string
	for _, act := range *plan.LocalDefinitions.Activities {
		activities = append(activities, act.Title)
	}
	require.Contains(t, activities, "etcd_cert_file")
	require.Contains(t, activities, "etcd_key_file")

	// Backmatter check
	require.Len(t, *plan.BackMatter.Resources, 1)
	resources := *plan.BackMatter.Resources
	require.Equal(t, "cis", resources[0].Title)
	require.Equal(t, "CIS Profile", resources[0].Description)
	require.Len(t, *resources[0].Rlinks, 1)

	// Link check
	require.Len(t, *plan.ReviewedControls.Links, 1)
	links := *plan.ReviewedControls.Links
	require.Equal(t, "includes-controls-from-source", links[0].Rel)
	require.Equal(t, fmt.Sprintf("#%s", resources[0].UUID), links[0].Href)

	// Validate against the schema
	validator := validation.NewSchemaValidator()
	oscalModels := oscalTypes.OscalModels{
		AssessmentPlan: plan,
	}

	data, _ := json.MarshalIndent(oscalModels, "", " ")
	fmt.Println(string(data))
	require.NoError(t, validator.Validate(oscalModels))
}

func TestSSPToAssessmentPlan(t *testing.T) {
	testDataPath := filepath.Join("../testdata", "test-ssp.json")

	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	ssp, err := models.NewSystemSecurityPlan(file, validation.NoopValidator{})
	require.NoError(t, err)
	require.NotNil(t, ssp)

	plan, err := SSPToAssessmentPlan(context.TODO(), *ssp, "importPath")
	require.NoError(t, err)

	require.Len(t, *plan.LocalDefinitions.Activities, 2)
	require.Len(t, *plan.AssessmentAssets.Components, 1)
	require.Len(t, *plan.AssessmentSubjects, 1)
	require.Len(t, plan.ReviewedControls.ControlSelections, 1)
	require.Len(t, *plan.Tasks, 1)
	tasks := *plan.Tasks
	require.Len(t, *tasks[0].AssociatedActivities, 2)

	var activities []string
	for _, act := range *plan.LocalDefinitions.Activities {
		activities = append(activities, act.Title)
	}
	require.Contains(t, activities, "rule-1")
	require.Contains(t, activities, "rule-2")

	// Validate against the schema
	validator := validation.NewSchemaValidator()
	oscalModels := oscalTypes.OscalModels{
		AssessmentPlan: plan,
	}
	require.NoError(t, validator.Validate(oscalModels))
}

func TestAssessmentPlanToAssessmentResults(t *testing.T) {
	testDataPath := filepath.Join("../testdata", "test-ap.json")

	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	plan, err := models.NewAssessmentPlan(file, validation.NoopValidator{})
	require.NoError(t, err)
	require.NotNil(t, plan)

	results, err := AssessmentPlanToAssessmentResults(*plan, "importPath")
	require.NoError(t, err)

	// Validate against the schema
	validator := validation.NewSchemaValidator()
	oscalModels := oscalTypes.OscalModels{
		AssessmentResults: results,
	}
	require.NoError(t, validator.Validate(oscalModels))
}
