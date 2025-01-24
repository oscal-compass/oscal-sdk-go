/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package plans

import (
	"context"
	"fmt"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/generators"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/oscal-compass/oscal-sdk-go/settings"
)

const (
	defaultSubjectType = "component"
	defaultTaskType    = "action"
)

type generateOpts struct {
	title           string
	importSSP       string
	localComponents []string
}

func (g *generateOpts) defaults() {
	g.title = generators.SampleRequiredString
	g.importSSP = generators.SampleRequiredString
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

// WithImport is a GenerateOption that sets the SystemSecurityPlan
// ImportSSP Href value.
func WithImport(importSSP string) GenerateOption {
	return func(opts *generateOpts) {
		opts.importSSP = importSSP
	}
}

// WithLocalComponents is a GenerateOptions that determine which components
// the `comps` argument in GenerateAssessmentPlan will be written to Components under
// LocalDefinitions. This would denote that these component are not defined in the SSP
// if the `WithImport` options is also used.
func WithLocalComponents(localComponents []string) GenerateOption {
	return func(opts *generateOpts) {
		opts.localComponents = localComponents
	}
}

// GenerateAssessmentPlan generates an AssessmentPlan for a set of Components and ImplementationSettings. The chosen inputs allow an Assessment Plan to be generated from
// a set of OSCAL ComponentDefinitions or a SystemSecurityPlan.
func GenerateAssessmentPlan(ctx context.Context, comps []components.Component, implementationSettings settings.ImplementationSettings, opts ...GenerateOption) (*oscalTypes.AssessmentPlan, error) {
	options := generateOpts{}
	options.defaults()
	for _, opt := range opts {
		opt(&options)
	}

	memoryStore := rules.NewMemoryStore()
	if err := memoryStore.IndexAll(comps); err != nil {
		return nil, fmt.Errorf("failed processing components for assessment plan %q: %w", options.title, err)
	}

	ruleBasedTask := oscalTypes.Task{
		UUID:                 uuid.NewUUID(),
		Title:                "Automated Assessment",
		Type:                 defaultTaskType,
		Description:          "Evaluation of defined rules for applicable comps.",
		Subjects:             &[]oscalTypes.AssessmentSubject{},
		AssociatedActivities: &[]oscalTypes.AssociatedActivity{},
	}

	var allActivities []oscalTypes.Activity
	var subjectSelectors []oscalTypes.SelectSubjectById
	for _, comp := range comps {
		if comp.Type() == components.Validation {
			continue
		}
		compTitle := comp.Title()
		componentActivities, err := ActivitiesForComponent(ctx, compTitle, memoryStore, implementationSettings)
		if err != nil {
			return nil, fmt.Errorf("error generating assessment activities for component %s: %w", compTitle, err)
		}
		allActivities = append(allActivities, componentActivities...)
		selector := oscalTypes.SelectSubjectById{
			Type:        defaultSubjectType,
			SubjectUuid: comp.UUID(),
		}
		subjectSelectors = append(subjectSelectors, selector)
		assessmentSubject := oscalTypes.AssessmentSubject{
			IncludeSubjects: &[]oscalTypes.SelectSubjectById{selector},
			Type:            defaultSubjectType,
		}

		associatedActivities := AssessmentActivities(assessmentSubject, componentActivities)
		*ruleBasedTask.AssociatedActivities = append(*ruleBasedTask.AssociatedActivities, associatedActivities...)
	}

	assessmentAssets := AssessmentAssets(comps)
	localDefinitions := oscalTypes.LocalDefinitions{
		Activities: &allActivities,
	}
	*ruleBasedTask.Subjects = append(*ruleBasedTask.Subjects, oscalTypes.AssessmentSubject{IncludeSubjects: &subjectSelectors})

	metadata := generators.NewSampleMetadata()
	metadata.Title = options.title

	assessmentPlan := &oscalTypes.AssessmentPlan{
		UUID: uuid.NewUUID(),
		ImportSsp: oscalTypes.ImportSsp{
			Href: options.importSSP,
		},
		Metadata: metadata,
		AssessmentSubjects: &[]oscalTypes.AssessmentSubject{
			{
				IncludeSubjects: &subjectSelectors,
				Type:            defaultSubjectType,
			},
		},
		LocalDefinitions: &localDefinitions,
		ReviewedControls: AllReviewedControls(implementationSettings),
		AssessmentAssets: &assessmentAssets,
		Tasks:            &[]oscalTypes.Task{ruleBasedTask},
	}

	return assessmentPlan, nil
}

// ActivitiesForComponent returns a list of activities with for a given component Title.
func ActivitiesForComponent(ctx context.Context, targetComponentID string, store rules.Store, implementationSettings settings.ImplementationSettings) ([]oscalTypes.Activity, error) {
	methodProp := oscalTypes.Property{
		Name:  "method",
		Value: "TEST",
	}

	appliedRules, err := settings.ApplyToComponent(ctx, targetComponentID, store, implementationSettings.AllSettings())
	if err != nil {
		return nil, fmt.Errorf("error getting applied rules for component %s: %w", targetComponentID, err)
	}

	var activities []oscalTypes.Activity
	for _, rule := range appliedRules {
		relatedControls, err := ReviewedControls(rule.Rule.ID, implementationSettings)
		if err != nil {
			return nil, err
		}

		var steps []oscalTypes.Step
		for _, check := range rule.Checks {
			checkStep := oscalTypes.Step{
				UUID:        uuid.NewUUID(),
				Title:       check.ID,
				Description: check.Description,
			}
			steps = append(steps, checkStep)
		}

		activity := oscalTypes.Activity{
			UUID:            uuid.NewUUID(),
			Description:     rule.Rule.Description,
			Props:           &[]oscalTypes.Property{methodProp},
			RelatedControls: &relatedControls,
			Title:           rule.Rule.ID,
			Steps:           &steps,
		}

		if rule.Rule.Parameter != nil {
			parameterProp := oscalTypes.Property{
				Name:  rule.Rule.Parameter.ID,
				Value: rule.Rule.Parameter.Value,
				Ns:    extensions.TrestleNameSpace,
				Class: "test-parameter",
			}
			*activity.Props = append(*activity.Props, parameterProp)
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

// AllReviewedControls returns ReviewControls with all the applicable controls ids in the implementation.
func AllReviewedControls(implementationSettings settings.ImplementationSettings) oscalTypes.ReviewedControls {
	applicableControls := implementationSettings.AllControls()
	return createReviewedControls(applicableControls)
}

// ReviewedControls returns ReviewedControls with controls ids that are associated with a given rule in ImplementationSettings.
func ReviewedControls(ruleId string, implementationSettings settings.ImplementationSettings) (oscalTypes.ReviewedControls, error) {
	applicableControls, err := implementationSettings.ApplicableControls(ruleId)
	if err != nil {
		return oscalTypes.ReviewedControls{}, fmt.Errorf("error getting applicable controls for rule %s: %w", ruleId, err)
	}
	return createReviewedControls(applicableControls), nil
}

func createReviewedControls(applicableControls []string) oscalTypes.ReviewedControls {
	var selectedControls []oscalTypes.AssessedControlsSelectControlById
	for _, control := range applicableControls {
		selector := oscalTypes.AssessedControlsSelectControlById{
			ControlId: control,
		}
		selectedControls = append(selectedControls, selector)
	}
	assessedControls := oscalTypes.AssessedControls{
		IncludeControls: &selectedControls,
	}

	return oscalTypes.ReviewedControls{
		ControlSelections: []oscalTypes.AssessedControls{
			assessedControls,
		},
	}
}

// AssessmentActivities returns an AssociatedActivity for addition to an Assessment Plan Task.
func AssessmentActivities(subject oscalTypes.AssessmentSubject, activities []oscalTypes.Activity) []oscalTypes.AssociatedActivity {
	var assocActivities []oscalTypes.AssociatedActivity
	for _, activity := range activities {
		assocActivity := oscalTypes.AssociatedActivity{
			ActivityUuid: activity.UUID,
			Subjects: []oscalTypes.AssessmentSubject{
				subject,
			},
		}
		assocActivities = append(assocActivities, assocActivity)
	}
	return assocActivities
}

// AssessmentAssets returns AssessmentAssets from validation components defined in the given DefinedComponents.
func AssessmentAssets(comps []components.Component) oscalTypes.AssessmentAssets {
	var systemComponents []oscalTypes.SystemComponent
	var usedComponents []oscalTypes.UsesComponent
	for _, component := range comps {
		if component.Type() == components.Validation {
			systemComponent, ok := component.AsSystemComponent()
			if ok {
				systemComponents = append(systemComponents, systemComponent)
				// This is an assumption that any validation components passed in
				// as input are part of a single Assessment Platform.
				usedComponent := oscalTypes.UsesComponent{
					ComponentUuid: systemComponent.UUID,
				}
				usedComponents = append(usedComponents, usedComponent)
			}

		}
	}
	// AssessmentPlatforms is a required field under AssessmentAssets
	assessmentPlatform := oscalTypes.AssessmentPlatform{
		UUID:           uuid.NewUUID(),
		Title:          generators.SampleRequiredString,
		UsesComponents: &usedComponents,
	}
	assessmentAssets := oscalTypes.AssessmentAssets{
		Components:          &systemComponents,
		AssessmentPlatforms: []oscalTypes.AssessmentPlatform{assessmentPlatform},
	}
	return assessmentAssets
}
