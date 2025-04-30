/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package plans

import (
	"context"
	"fmt"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/models/modelutils"
	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/oscal-compass/oscal-sdk-go/settings"
)

const (
	defaultSubjectType = "component"
	defaultTaskType    = "action"
)

type generateOpts struct {
	title     string
	importSSP string
}

func (g *generateOpts) defaults() {
	g.title = models.SampleRequiredString
	g.importSSP = models.SampleRequiredString
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

// GenerateAssessmentPlan generates an AssessmentPlan for a set of Components and ImplementationSettings. The chosen inputs allow an Assessment Plan to be generated from
// a set of OSCAL ComponentDefinitions or a SystemSecurityPlan.
//
// If the `WithImport` is not set, all input components are set as Components in the Local Definitions.
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

	var (
		allActivities    []oscalTypes.Activity
		subjectSelectors []oscalTypes.SelectSubjectById
		localComponents  []components.Component
		ruleBasedTask    = newTask()
	)

	for _, comp := range comps {

		// process components
		if comp.Type() == components.Validation {
			continue
		}
		compTitle := comp.Title()
		componentActivities, err := ActivitiesForComponent(ctx, compTitle, memoryStore, implementationSettings)
		if err != nil {
			return nil, fmt.Errorf("error generating assessment activities for component %s: %w", compTitle, err)
		}
		if len(componentActivities) == 0 {
			continue
		}

		// create assessment plan objects
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

		if options.importSSP == models.SampleRequiredString {
			// In this use case, there is no linked SSP, making specified Components
			// locally defined.
			localComponents = append(localComponents, comp)
		}
	}

	assessmentAssets := AssessmentAssets(comps)
	taskSubjects := oscalTypes.AssessmentSubject{
		IncludeSubjects: &subjectSelectors,
		Type:            defaultSubjectType,
	}
	*ruleBasedTask.Subjects = append(*ruleBasedTask.Subjects, taskSubjects)

	metadata := models.NewSampleMetadata()
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
		LocalDefinitions: createLocalDefinitions(allActivities, localComponents),
		ReviewedControls: AllReviewedControls(implementationSettings),
		AssessmentAssets: &assessmentAssets,
		Tasks:            &[]oscalTypes.Task{ruleBasedTask},
	}

	return assessmentPlan, nil
}

// newTask creates a new OSCAL Task with default values.
func newTask() oscalTypes.Task {
	return oscalTypes.Task{
		UUID:                 uuid.NewUUID(),
		Title:                "Automated Assessment",
		Type:                 defaultTaskType,
		Description:          "Evaluation of defined rules for components.",
		Subjects:             &[]oscalTypes.AssessmentSubject{},
		AssociatedActivities: &[]oscalTypes.AssociatedActivity{},
	}
}

// ActivitiesForComponent returns a list of activities with for a given component Title.
//
// The mapping between a RuleSet and Activity is as follows:
// Rule -> Activity
// ID -> Title
// Parameter -> Activity Property
// Check -> Activity Step
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
			Steps:           modelutils.NilIfEmpty(&steps),
		}

		for _, rp := range rule.Rule.Parameters {
			parameterProp := oscalTypes.Property{
				Name:  rp.ID,
				Value: rp.Value,
				Ns:    extensions.TrestleNameSpace,
				Class: extensions.TestParameterClass,
			}
			*activity.Props = append(*activity.Props, parameterProp)
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

// createLocationDefinitions for an AssessmentPlan from given Activities and components marked as local.
func createLocalDefinitions(activities []oscalTypes.Activity, localComps []components.Component) *oscalTypes.LocalDefinitions {
	localDefinitions := &oscalTypes.LocalDefinitions{
		Activities: &activities,
	}
	if len(localComps) == 0 {
		return localDefinitions
	}

	localDefinitions.Components = &[]oscalTypes.SystemComponent{}
	for _, comp := range localComps {
		sysComp, ok := comp.AsSystemComponent()
		if ok {
			*localDefinitions.Components = append(*localDefinitions.Components, sysComp)
		}
	}

	return localDefinitions
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

func createReviewedControls(selectedControls []oscalTypes.AssessedControlsSelectControlById) oscalTypes.ReviewedControls {
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
		Title:          models.SampleRequiredString,
		UsesComponents: modelutils.NilIfEmpty(&usedComponents),
	}

	assessmentAssets := oscalTypes.AssessmentAssets{
		Components:          &systemComponents,
		AssessmentPlatforms: []oscalTypes.AssessmentPlatform{assessmentPlatform},
	}
	return assessmentAssets
}
