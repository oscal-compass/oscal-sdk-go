/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package models

import (
	"encoding/json"
	"io"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/oscal-compass/oscal-sdk-go/validation"
)

// NewCatalog creates a new OSCAL-based control catalog using types from `go-oscal`.
func NewCatalog(reader io.Reader, validator validation.Validator) (catalog *oscalTypes.Catalog, err error) {
	var oscalModels oscalTypes.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	err = validator.Validate(oscalModels)
	if err != nil {
		return nil, err
	}

	return oscalModels.Catalog, nil
}

// NewProfile creates a new OSCAL-based profile using types from `go-oscal`.
func NewProfile(reader io.Reader, validator validation.Validator) (profile *oscalTypes.Profile, err error) {
	var oscalModels oscalTypes.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	err = validator.Validate(oscalModels)
	if err != nil {
		return nil, err
	}

	return oscalModels.Profile, nil
}

// NewComponentDefinition creates a new OSCAL-based component definition using types from `go-oscal`.
func NewComponentDefinition(reader io.Reader, validator validation.Validator) (componentDefinition *oscalTypes.ComponentDefinition, err error) {
	var oscalModels oscalTypes.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	err = validator.Validate(oscalModels)
	if err != nil {
		return nil, err
	}

	return oscalModels.ComponentDefinition, nil
}

// NewSystemSecurityPlan creates a new OSCAL-based system security plan using types from `go-oscal`.
func NewSystemSecurityPlan(reader io.Reader, validator validation.Validator) (systemSecurityPlan *oscalTypes.SystemSecurityPlan, err error) {
	var oscalModels oscalTypes.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	err = validator.Validate(oscalModels)
	if err != nil {
		return nil, err
	}

	return oscalModels.SystemSecurityPlan, nil
}

// NewAssessmentPlan creates a new OSCAL-based assessment plan using types from `go-oscal`.
func NewAssessmentPlan(reader io.Reader, validator validation.Validator) (assessmentPlan *oscalTypes.AssessmentPlan, err error) {
	var oscalModels oscalTypes.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	err = validator.Validate(oscalModels)
	if err != nil {
		return nil, err
	}

	return oscalModels.AssessmentPlan, nil
}

// NewAssessmentResults creates a new OSCAL-based assessment results set using types from `go-oscal`.
func NewAssessmentResults(reader io.Reader, validator validation.Validator) (assessmentResults *oscalTypes.AssessmentResults, err error) {
	var oscalModels oscalTypes.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	err = validator.Validate(oscalModels)
	if err != nil {
		return nil, err
	}

	return oscalModels.AssessmentResults, nil
}

// NewPOAM creates a new OSCAL-based plan of action and milestones using types from `go-oscal`.
func NewPOAM(reader io.Reader, validator validation.Validator) (pOAM *oscalTypes.PlanOfActionAndMilestones, err error) {
	var oscalModels oscalTypes.OscalModels
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return nil, err
	}

	err = validator.Validate(oscalModels)
	if err != nil {
		return nil, err
	}

	return oscalModels.PlanOfActionAndMilestones, nil
}
