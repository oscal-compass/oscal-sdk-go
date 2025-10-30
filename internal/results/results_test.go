/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package results

import (
	"os"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/validation"
)

func TestGenerateAssessmentResults(t *testing.T) {
	file, err := os.Open("../../testdata/test-ap.json")
	require.NoError(t, err)
	defer file.Close()
	defaultAssessmentPlan, err := models.NewAssessmentPlan(file, validation.NoopValidator{})
	require.NoError(t, err)
	require.NotNil(t, defaultAssessmentPlan)

	tests := []struct {
		name           string
		inputOptions   []GenerateOption
		assessmentPlan *oscalTypes.AssessmentPlan
		assertFunc     func(*testing.T, *oscalTypes.AssessmentResults)
		expError       string
	}{
		{
			name:           "Success/Defaults",
			inputOptions:   nil,
			assessmentPlan: defaultAssessmentPlan,
			assertFunc: func(t *testing.T, results *oscalTypes.AssessmentResults) {
				require.Len(t, results.Results, 1)
				result := results.Results[0]
				require.Equal(t, result.Title, "Result For Task \"Automated Assessment\"")
				require.Equal(t, result.Description, "OSCAL Assessment Result For Task \"Automated Assessment\"")

				require.NotNil(t, result.Observations)
				require.Len(t, *result.Observations, 1)
				observations := *result.Observations
				observation := observations[0]

				expectedOrigins := []oscalTypes.Origin{
					{
						Actors: []oscalTypes.OriginActor{
							{
								ActorUuid: "701c70f1-482b-42b0-a419-9870158cd9e2",
								Type:      defaultActor,
							},
						},
						RelatedTasks: &[]oscalTypes.RelatedTask{
							{
								TaskUuid: "0733aaa9-9743-4971-967c-bbd951bb9026",
								Subjects: &[]oscalTypes.AssessmentSubject{
									{
										Type: "component",
										IncludeSubjects: &[]oscalTypes.SelectSubjectById{
											{
												SubjectUuid: "4e19131e-b361-4f0e-8262-02bf4456202e",
												Type:        "component",
											},
										},
									},
								},
							},
						},
					},
				}
				require.NotNil(t, observation.Origins)
				require.Len(t, *observation.Origins, 1)
				require.Equal(t, expectedOrigins, *observation.Origins)
			},
		},
		{
			name: "Success/WithObservationsMatchingCheckId",
			inputOptions: []GenerateOption{
				WithObservations([]oscalTypes.Observation{
					{
						Title: "check-1",
						Props: &[]oscalTypes.Property{
							{
								Name:  extensions.AssessmentRuleIdProp,
								Ns:    extensions.TrestleNameSpace,
								Value: "rule-1",
							},
							{
								Name:  extensions.AssessmentCheckIdProp,
								Ns:    extensions.TrestleNameSpace,
								Value: "check-1",
							},
						},
					},
				}),
			},
			assessmentPlan: defaultAssessmentPlan,
			assertFunc: func(t *testing.T, results *oscalTypes.AssessmentResults) {
				require.Len(t, results.Results, 1)
				result := results.Results[0]
				require.Equal(t, result.Title, "Result For Task \"Automated Assessment\"")
				require.Equal(t, result.Description, "OSCAL Assessment Result For Task \"Automated Assessment\"")

				expectedControls := oscalTypes.ReviewedControls{
					ControlSelections: []oscalTypes.AssessedControls{
						{
							IncludeControls: &[]oscalTypes.AssessedControlsSelectControlById{
								{
									ControlId: "ex-2",
								},
								{
									ControlId: "ex-1",
								},
							},
						},
						{
							IncludeControls: &[]oscalTypes.AssessedControlsSelectControlById{
								{
									ControlId: "ex-1",
								},
							},
						},
					},
				}
				require.Equal(t, expectedControls, result.ReviewedControls)
				require.NotNil(t, result.Observations)
				require.Len(t, *result.Observations, 1)
				observations := *result.Observations
				observation := observations[0]

				expectedOrigins := []oscalTypes.Origin{
					{
						Actors: []oscalTypes.OriginActor{
							{
								ActorUuid: "701c70f1-482b-42b0-a419-9870158cd9e2",
								Type:      defaultActor,
							},
						},
						RelatedTasks: &[]oscalTypes.RelatedTask{
							{
								TaskUuid: "0733aaa9-9743-4971-967c-bbd951bb9026",
								Subjects: &[]oscalTypes.AssessmentSubject{
									{
										Type: "component",
										IncludeSubjects: &[]oscalTypes.SelectSubjectById{
											{
												SubjectUuid: "4e19131e-b361-4f0e-8262-02bf4456202e",
												Type:        "component",
											},
										},
									},
								},
							},
						},
					},
				}
				require.NotNil(t, observation.Origins)
				require.Len(t, *observation.Origins, 1)
				require.Equal(t, expectedOrigins, *observation.Origins)
			},
		},
		{
			name: "Success/WithObservationsMatchingTitle",
			inputOptions: []GenerateOption{
				WithObservations([]oscalTypes.Observation{
					{
						Title: "check-1",
						Props: &[]oscalTypes.Property{
							{
								Name:  extensions.AssessmentRuleIdProp,
								Ns:    extensions.TrestleNameSpace,
								Value: "rule-1",
							},
							{
								Name:  extensions.AssessmentCheckIdProp,
								Ns:    extensions.TrestleNameSpace,
								Value: "check-2",
							},
						},
					},
				}),
			},
			assessmentPlan: defaultAssessmentPlan,
			assertFunc: func(t *testing.T, results *oscalTypes.AssessmentResults) {
				require.Len(t, results.Results, 1)
				result := results.Results[0]
				require.Equal(t, result.Title, "Result For Task \"Automated Assessment\"")
				require.Equal(t, result.Description, "OSCAL Assessment Result For Task \"Automated Assessment\"")

				expectedControls := oscalTypes.ReviewedControls{
					ControlSelections: []oscalTypes.AssessedControls{
						{
							IncludeControls: &[]oscalTypes.AssessedControlsSelectControlById{
								{
									ControlId: "ex-2",
								},
								{
									ControlId: "ex-1",
								},
							},
						},
						{
							IncludeControls: &[]oscalTypes.AssessedControlsSelectControlById{
								{
									ControlId: "ex-1",
								},
							},
						},
					},
				}
				require.Equal(t, expectedControls, result.ReviewedControls)
				require.NotNil(t, result.Observations)
				require.Len(t, *result.Observations, 1)
				observations := *result.Observations
				observation := observations[0]

				expectedOrigins := []oscalTypes.Origin{
					{
						Actors: []oscalTypes.OriginActor{
							{
								ActorUuid: "701c70f1-482b-42b0-a419-9870158cd9e2",
								Type:      defaultActor,
							},
						},
						RelatedTasks: &[]oscalTypes.RelatedTask{
							{
								TaskUuid: "0733aaa9-9743-4971-967c-bbd951bb9026",
								Subjects: &[]oscalTypes.AssessmentSubject{
									{
										Type: "component",
										IncludeSubjects: &[]oscalTypes.SelectSubjectById{
											{
												SubjectUuid: "4e19131e-b361-4f0e-8262-02bf4456202e",
												Type:        "component",
											},
										},
									},
								},
							},
						},
					},
				}
				require.NotNil(t, observation.Origins)
				require.Len(t, *observation.Origins, 1)
				require.Equal(t, expectedOrigins, *observation.Origins)
			},
		},
		{
			name: "Success/WithObservationsNoMatching",
			inputOptions: []GenerateOption{
				WithObservations([]oscalTypes.Observation{
					{
						Title: "check-2",
						Props: &[]oscalTypes.Property{
							{
								Name:  extensions.AssessmentRuleIdProp,
								Ns:    extensions.TrestleNameSpace,
								Value: "rule-1",
							},
						},
					},
				}),
			},
			assessmentPlan: defaultAssessmentPlan,
			assertFunc: func(t *testing.T, results *oscalTypes.AssessmentResults) {
				require.Len(t, results.Results, 1)
				result := results.Results[0]
				require.Equal(t, result.Title, "Result For Task \"Automated Assessment\"")
				require.Equal(t, result.Description, "OSCAL Assessment Result For Task \"Automated Assessment\"")

				expectedControls := oscalTypes.ReviewedControls{
					ControlSelections: []oscalTypes.AssessedControls{
						{
							IncludeControls: &[]oscalTypes.AssessedControlsSelectControlById{
								{
									ControlId: "ex-2",
								},
								{
									ControlId: "ex-1",
								},
							},
						},
						{
							IncludeControls: &[]oscalTypes.AssessedControlsSelectControlById{
								{
									ControlId: "ex-1",
								},
							},
						},
					},
				}
				require.Equal(t, expectedControls, result.ReviewedControls)
				require.NotNil(t, result.Observations)
				require.Len(t, *result.Observations, 1)
				observations := *result.Observations
				observation := observations[0]
				// Created a new observation when no matchings
				require.Equal(t, observation.Title, "check-1")
				expectedOrigins := []oscalTypes.Origin{
					{
						Actors: []oscalTypes.OriginActor{
							{
								ActorUuid: "701c70f1-482b-42b0-a419-9870158cd9e2",
								Type:      defaultActor,
							},
						},
						RelatedTasks: &[]oscalTypes.RelatedTask{
							{
								TaskUuid: "0733aaa9-9743-4971-967c-bbd951bb9026",
								Subjects: &[]oscalTypes.AssessmentSubject{
									{
										Type: "component",
										IncludeSubjects: &[]oscalTypes.SelectSubjectById{
											{
												SubjectUuid: "4e19131e-b361-4f0e-8262-02bf4456202e",
												Type:        "component",
											},
										},
									},
								},
							},
						},
					},
				}

				require.NotNil(t, observation.Origins)
				require.Len(t, *observation.Origins, 1)
				require.Equal(t, expectedOrigins, *observation.Origins)
			},
		},
		{
			name: "Success/VerifyCreateOrGetWithProps",
			inputOptions: []GenerateOption{
				WithObservations([]oscalTypes.Observation{
					{
						Title: "check-1",
						Props: &[]oscalTypes.Property{
							{
								Name:  extensions.AssessmentRuleIdProp,
								Ns:    extensions.TrestleNameSpace,
								Value: "rule-1",
							},
							{
								Name:  extensions.AssessmentCheckIdProp,
								Ns:    extensions.TrestleNameSpace,
								Value: "check-1",
							},
						},
					},
				}),
			},
			assessmentPlan: defaultAssessmentPlan,
			assertFunc: func(t *testing.T, results *oscalTypes.AssessmentResults) {
				require.Len(t, results.Results, 1)
				result := results.Results[0]
				require.NotNil(t, result.Observations)
				require.Len(t, *result.Observations, 1)

				observation := (*result.Observations)[0]

				// Verify that the observation has properties
				require.NotNil(t, observation.Props)
				props := *observation.Props

				var ruleIdFound, checkIdFound bool
				for _, prop := range props {
					if prop.Name == extensions.AssessmentRuleIdProp && prop.Ns == extensions.TrestleNameSpace {
						require.Equal(t, "rule-1", prop.Value)
						ruleIdFound = true
					}
					if prop.Name == extensions.AssessmentCheckIdProp && prop.Ns == extensions.TrestleNameSpace {
						require.Equal(t, "check-1", prop.Value)
						checkIdFound = true
					}
				}
				require.True(t, ruleIdFound)
				require.True(t, checkIdFound)
			},
		},
		{
			name:           "Failure/NoTasksPlan",
			assessmentPlan: &oscalTypes.AssessmentPlan{},
			expError:       "assessment plan tasks cannot be empty",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			assessmentResults, err := GenerateAssessmentResults(*c.assessmentPlan, c.inputOptions...)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				c.assertFunc(t, assessmentResults)
			}
		})
	}
}
