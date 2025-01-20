/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package plans

import (
	"context"
	"fmt"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscaltypes112 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/oscal-compass/oscal-sdk-go/settings"
)

const (
	defaultVersion     = "0.1.0"
	defaultSubjectType = "component"
	defaultTaskType    = "action"
)

type generateOpts struct {
	title     string
	importSSP string
}

func (g *generateOpts) defaults() {
	g.title = models.DefaultRequiredString
	g.importSSP = models.DefaultRequiredString
}

// GenerateOption define an options to tune the behavior of
// GenerateAssessmentPlan.
type GenerateOption func(opts *generateOpts)

// WithTitle is a GenerateOptions that sets the AssessmentPlan title
// in the metadata.
func WithTitle(title string) GenerateOption {
	return func(opts *generateOpts) {
		opts.title = title
	}
}

// WithImport is a GenerateOptions that sets the SystemSecurityPlan
// ImportSSP Href value.
func WithImport(importSSP string) GenerateOption {
	return func(opts *generateOpts) {
		opts.importSSP = importSSP
	}
}

// GenerateAssessmentPlan generates an AssessmentPlan for a set of Components and implementationSettings. The chosen inputs allow an Assessment Plan to be generated from
// a set of OSCAL Component Definitions or a System Security Plan.
func GenerateAssessmentPlan(ctx context.Context, components []components.Component, implementationSettings *settings.ImplementationSettings, opts ...GenerateOption) (*oscaltypes112.AssessmentPlan, error) {
	options := generateOpts{}
	options.defaults()
	for _, opt := range opts {
		opt(&options)
	}

	memoryStore := rules.NewMemoryStore()
	if err := memoryStore.IndexAll(components); err != nil {
		return nil, fmt.Errorf("failed processing components for assessment plan %q: %w", options.title, err)
	}

	ruleBasedTask := oscaltypes112.Task{
		UUID:                 uuid.NewUUID(),
		Title:                "Automated Assessment",
		Type:                 defaultTaskType,
		Description:          "Evaluation of defined rules for applicable components.",
		Subjects:             &[]oscaltypes112.AssessmentSubject{},
		AssociatedActivities: &[]oscaltypes112.AssociatedActivity{},
	}

	var allActivities []oscaltypes112.Activity
	var subjectSelectors []oscaltypes112.SelectSubjectById
	for _, comp := range components {
		compTitle := comp.Title()
		appliedRules, err := settings.ApplyToComponent(ctx, compTitle, memoryStore, implementationSettings)
		if err != nil {
			return nil, fmt.Errorf("error getting applied rules for component %s: %w", compTitle, err)
		}
		componentActivities, err := Activities(appliedRules, *implementationSettings)
		if err != nil {
			return nil, fmt.Errorf("error generating assessment activities for component %s: %w", compTitle, err)
		}
		allActivities = append(allActivities, componentActivities...)
		selector := oscaltypes112.SelectSubjectById{
			Type:        defaultSubjectType,
			SubjectUuid: comp.UUID(),
		}
		subjectSelectors = append(subjectSelectors, selector)
		assessmentSubject := oscaltypes112.AssessmentSubject{
			IncludeSubjects: &[]oscaltypes112.SelectSubjectById{selector},
		}

		associatedActivities := AssessmentActivities(assessmentSubject, componentActivities)
		*ruleBasedTask.AssociatedActivities = append(*ruleBasedTask.AssociatedActivities, associatedActivities...)
	}

	assessmentAssets := AssessmentAssets(components)
	localDefinitions := oscaltypes112.LocalDefinitions{
		Activities: &allActivities,
	}
	*ruleBasedTask.Subjects = append(*ruleBasedTask.Subjects, oscaltypes112.AssessmentSubject{IncludeSubjects: &subjectSelectors})

	assessmentPlan := &oscaltypes112.AssessmentPlan{
		UUID: uuid.NewUUID(),
		ImportSsp: oscaltypes112.ImportSsp{
			Href: options.importSSP,
		},
		Metadata: oscaltypes112.Metadata{
			Title:        options.title,
			LastModified: time.Now(),
			OscalVersion: models.OSCALVersion,
			Version:      defaultVersion,
		},
		AssessmentSubjects: &[]oscaltypes112.AssessmentSubject{
			{
				IncludeSubjects: &subjectSelectors,
				Type:            defaultSubjectType,
			},
		},
		LocalDefinitions: &localDefinitions,
		ReviewedControls: getAllReviewControls(localDefinitions),
		AssessmentAssets: &assessmentAssets,
		Tasks:            &[]oscaltypes112.Task{ruleBasedTask},
	}

	return assessmentPlan, nil
}

// Activities returns a list of activities with the given list of applied rules.
func Activities(appliedRules []extensions.RuleSet, implementationSettings settings.ImplementationSettings) ([]oscaltypes112.Activity, error) {
	methodProp := oscaltypes112.Property{
		Name:  "method",
		Value: "TEST",
	}

	var activities []oscaltypes112.Activity

	for _, rule := range appliedRules {
		relatedControls, err := ReviewedControls(rule.Rule.ID, implementationSettings)
		if err != nil {
			return nil, err
		}

		activity := oscaltypes112.Activity{
			UUID:            uuid.NewUUID(),
			Description:     rule.Rule.Description,
			Props:           &[]oscaltypes112.Property{methodProp},
			RelatedControls: &relatedControls,
			Title:           rule.Rule.ID,
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func getAllReviewControls(localDefinitions oscaltypes112.LocalDefinitions) oscaltypes112.ReviewedControls {
	var allControls = set.New[string]()
	var allControlSelections []oscaltypes112.AssessedControlsSelectControlById
	if localDefinitions.Activities == nil {
		return oscaltypes112.ReviewedControls{}
	}
	for _, activity := range *localDefinitions.Activities {
		for _, assessedControls := range activity.RelatedControls.ControlSelections {
			if assessedControls.IncludeControls == nil {
				continue
			}
			for _, control := range *assessedControls.IncludeControls {
				if allControls.Has(control.ControlId) {
					continue
				}
				allControlSelections = append(allControlSelections, control)
				allControls.Add(control.ControlId)
			}
		}
	}
	return oscaltypes112.ReviewedControls{
		ControlSelections: []oscaltypes112.AssessedControls{
			{
				IncludeControls: &allControlSelections,
			},
		},
	}
}

// ReviewedControls returns ReviewedControls with controls ids that are associated with a given rule in ImplementationSettings.
func ReviewedControls(ruleId string, implementationSetting settings.ImplementationSettings) (oscaltypes112.ReviewedControls, error) {
	var selectedControls []oscaltypes112.AssessedControlsSelectControlById
	applicableControls, err := implementationSetting.ApplicableControls(ruleId)
	if err != nil {
		return oscaltypes112.ReviewedControls{}, err
	}

	for _, control := range applicableControls {
		selector := oscaltypes112.AssessedControlsSelectControlById{
			ControlId: control,
		}
		selectedControls = append(selectedControls, selector)
	}
	assessedControls := oscaltypes112.AssessedControls{
		IncludeControls: &selectedControls,
	}

	return oscaltypes112.ReviewedControls{
		ControlSelections: []oscaltypes112.AssessedControls{
			assessedControls,
		},
	}, nil
}

// AssessmentActivities returns an AssociatedActivity for addition to an Assessment Plan Task.
func AssessmentActivities(subject oscaltypes112.AssessmentSubject, activities []oscaltypes112.Activity) []oscaltypes112.AssociatedActivity {
	var assocActivities []oscaltypes112.AssociatedActivity
	for _, activity := range activities {
		assocActivity := oscaltypes112.AssociatedActivity{
			ActivityUuid: activity.UUID,
			Subjects: []oscaltypes112.AssessmentSubject{
				subject,
			},
		}
		assocActivities = append(assocActivities, assocActivity)
	}
	return assocActivities
}

// AssessmentAssets returns AssessmentAssets from validation components defined in the given DefinedComponents.
func AssessmentAssets(comps []components.Component) oscaltypes112.AssessmentAssets {
	var systemComponents []oscaltypes112.SystemComponent
	for _, component := range comps {
		if component.Type() == components.Validation {
			systemComponent, ok := component.AsSystemComponent()
			if ok {
				systemComponents = append(systemComponents, systemComponent)
			}

		}
	}
	assessmentAssets := oscaltypes112.AssessmentAssets{
		Components: &systemComponents,
	}
	return assessmentAssets
}
