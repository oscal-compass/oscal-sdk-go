/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package results

import (
	"fmt"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/models"
)

type generateOpts struct {
	title        string
	importAP     string
	observations []oscalTypes.Observation
}

func (g *generateOpts) defaults() {
	g.title = models.SampleRequiredString
	g.importAP = models.SampleRequiredString
}

// GenerateOption defines an option to tune the behavior of the
// GenerateAssessmentPlan function.
type GenerateOption func(opts *generateOpts)

// WithTitle is a GenerateOption that sets the AssessmentPlan title
// in the metadata.
func WithTitle(title string) GenerateOption {
	return func(opts *generateOpts) {
		opts.title = title
	}
}

// WithImport is a GenerateOption that sets the AssessmentPlan
// ImportAP Href value.
func WithImport(importAP string) GenerateOption {
	return func(opts *generateOpts) {
		opts.importAP = importAP
	}
}

// WithObservations is a GenerateOption that adds pre-processed OSCAL Observations
// to the Assessment Result for associated to Assessment Plan Activities.
func WithObservations(observations []oscalTypes.Observation) GenerateOption {
	return func(opts *generateOpts) {
		opts.observations = observations
	}
}

// GenerateAssessmentResults generates an AssessmentPlan for a set of Components and ImplementationSettings. The chosen inputs allow an Assessment Plan to be generated from
// a set of OSCAL ComponentDefinitions or a SystemSecurityPlan.
//
// If the `WithImport` is not set, all input components are set as Components in the Local Definitions.
func GenerateAssessmentResults(plan oscalTypes.AssessmentPlan, opts ...GenerateOption) (*oscalTypes.AssessmentResults, error) {
	options := generateOpts{}
	options.defaults()
	for _, opt := range opts {
		opt(&options)
	}

	metadata := models.NewSampleMetadata()
	metadata.Title = options.title

	assessmentResults := &oscalTypes.AssessmentResults{
		UUID: uuid.NewUUID(),
		ImportAp: oscalTypes.ImportAp{
			Href: options.importAP,
		},
		Metadata: metadata,
	}

	if plan.Tasks == nil {
		return assessmentResults, nil
	}

	// FIXME(jpower432):  Identify the automated assessment task and create corresponding results
	// now just error if there is more than one task
	tasks := *plan.Tasks
	if len(tasks) != 1 {
		return assessmentResults, fmt.Errorf("the assessment plan should have one task")
	}
	task := tasks[0]
	result := oscalTypes.Result{
		// FIXME(jpower432): This may need to be an aggregated of
		// Review Controls in the Associated Task instead
		ReviewedControls: plan.ReviewedControls,
		Title:            "Automated Assessment Result",
		Description:      fmt.Sprintf("Assessment Results For Task %q", task.Title),
		Start:            time.Now(),
		UUID:             uuid.NewUUID(),
	}

	// Perform a couple checks and return early here if needed
	if len(options.observations) == 0 {
		assessmentResults.Results = []oscalTypes.Result{result}
		return assessmentResults, nil
	}

	if task.AssociatedActivities == nil || plan.LocalDefinitions == nil || plan.LocalDefinitions.Activities == nil {
		result.Observations = &options.observations
		assessmentResults.Results = []oscalTypes.Result{result}
		return assessmentResults, nil
	}

	result.Observations = &options.observations
	assessmentResults.Results = []oscalTypes.Result{result}

	return assessmentResults, nil
}

func LinkObservations(observations []oscalTypes.Observation, task oscalTypes.Task, activities []oscalTypes.Activity) {
	activitiesByRule := make(map[string]oscalTypes.Activity)
	for _, activity := range activities {
		activitiesByRule[activity.Title] = activity
	}

	associatedActivityByActivityUUID := make(map[string]oscalTypes.AssociatedActivity)
	for _, activity := range *task.AssociatedActivities {
		associatedActivityByActivityUUID[activity.ActivityUuid] = activity
	}

	for i, observation := range observations {
		if observations[i].Props == nil {
			continue
		}
		ruleIdProp, found := extensions.GetTrestleProp(extensions.AssessmentRuleIdProp, *observation.Props)
		if !found {
			continue
		}

		activity, found := activitiesByRule[ruleIdProp.Value]
		if !found {
			continue
		}

		// Should methods be attached here?
		if activity.Props != nil {
			methods := extensions.FindAllProps(*activity.Props, extensions.WithName("method"), extensions.WithNamespace(""))
			for _, method := range methods {
				observations[i].Methods = append(observations[i].Methods, method.Value)
			}
		}

		// Set the observation origins
		// QUESTION(jpower432): How should Actor be set?
		assocActivity, found := associatedActivityByActivityUUID[activity.UUID]
		if !found {
			continue
		}
		origins := []oscalTypes.Origin{
			{
				RelatedTasks: &[]oscalTypes.RelatedTask{
					{
						TaskUuid: task.UUID,
						Subjects: &assocActivity.Subjects,
					},
				},
			},
		}
		observation.Origins = &origins
	}
}
