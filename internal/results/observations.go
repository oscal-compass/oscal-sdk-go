/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package results

import (
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
)

// observationsManager indexes and manages OSCAL Observations
// to support Assessment Result generation.
type observationsManager struct {
	observationsByCheck map[string]oscalTypes.Observation
	actorsByCheck       map[string]string
}

// newObservationManager creates an observationManager struct loaded with
// actor information from the Assessment Plan Assessment Assets.
func newObservationManager(plan oscalTypes.AssessmentPlan) *observationsManager {
	// Index validation components to set the Actor information
	m := &observationsManager{
		observationsByCheck: make(map[string]oscalTypes.Observation),
		actorsByCheck:       make(map[string]string),
	}
	if plan.AssessmentAssets != nil && plan.AssessmentAssets.Components != nil {
		for _, comp := range *plan.AssessmentAssets.Components {
			if comp.Props == nil {
				continue
			}
			checkProps := extensions.FindAllProps(*comp.Props, extensions.WithName(extensions.CheckIdProp))
			for _, check := range checkProps {
				m.actorsByCheck[check.Value] = comp.UUID
			}
		}
	}
	return m
}

// load indexing and updates a set of given observations.
func (o *observationsManager) load(observations []oscalTypes.Observation) {
	for _, observation := range observations {
		o.updateObservation(&observation)
	}
}

// createOrGet return an existing observation or a newly created one.
func (o *observationsManager) createOrGet(checkId string) oscalTypes.Observation {
	for _, observation := range o.observationsByCheck {
		// Loop through the Props slice to find the AssessmentCheckIdProp
		if observation.Props == nil {
			continue
		}
		check, found := extensions.GetTrestleProp(extensions.AssessmentCheckIdProp, *observation.Props)
		if !found || check.Value != checkId {
			continue
		} else {
			return observation
		}
	}
	// Using observation title as a fallback
	observation, ok := o.observationsByCheck[checkId]
	if ok {
		return observation
	}

	emptyObservation := oscalTypes.Observation{
		UUID:      uuid.NewUUID(),
		Title:     checkId,
		Collected: time.Now(),
	}
	o.updateObservation(&emptyObservation)
	return emptyObservation
}

// updateObservation with Origin Actor information
func (o *observationsManager) updateObservation(observation *oscalTypes.Observation) {
	actor, found := o.actorsByCheck[observation.Title]
	if found {
		origins := []oscalTypes.Origin{
			{
				Actors: []oscalTypes.OriginActor{
					{
						Type:      defaultActor,
						ActorUuid: actor,
					},
				},
			},
		}
		observation.Origins = &origins
	}
	o.observationsByCheck[observation.Title] = *observation
}
