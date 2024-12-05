/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package generators

import (
	"encoding/json"
	"io"

	oscal112 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
)

// NewCatalog creates a new OSCAL-based control catalog using types from `go-oscal`.
func NewCatalog(reader io.Reader) (catalog *oscal112.Catalog, err error) {
	var oscalModels oscal112.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	return oscalModels.Catalog, nil
}

// NewProfile creates a new OSCAL-based profile using types from `go-oscal`.
func NewProfile(reader io.Reader) (profile *oscal112.Profile, err error) {
	var oscalModels oscal112.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	return oscalModels.Profile, nil
}

// NewComponentDefinition creates a new OSCAL-based component definition using types from `go-oscal`.
func NewComponentDefinition(reader io.Reader) (componentDefinition *oscal112.ComponentDefinition, err error) {
	var oscalModels oscal112.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	return oscalModels.ComponentDefinition, nil
}

// NewSystemSecurityPlan creates a new OSCAL-based system security plan using types from `go-oscal`.
func NewSystemSecurityPlan(reader io.Reader) (systemSecurityPlan *oscal112.SystemSecurityPlan, err error) {
	var oscalModels oscal112.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	return oscalModels.SystemSecurityPlan, nil
}

// NewAssessmentPlan creates a new OSCAL-based assessment plan using types from `go-oscal`.
func NewAssessmentPlan(reader io.Reader) (assessmentPlan *oscal112.AssessmentPlan, err error) {
	var oscalModels oscal112.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	return oscalModels.AssessmentPlan, nil
}

// NewAssessmentResults creates a new OSCAL-based assessment results set using types from `go-oscal`.
func NewAssessmentResults(reader io.Reader) (assessmentResult *oscal112.AssessmentResults, err error) {
	var oscalModels oscal112.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	return oscalModels.AssessmentResults, nil
}

// NewPOAM creates a new OSCAL-based plan of action and milestones using types from `go-oscal`.
func NewPOAM(reader io.Reader) (pOAM *oscal112.PlanOfActionAndMilestones, err error) {
	var oscalModels oscal112.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	return oscalModels.PlanOfActionAndMilestones, nil
}
