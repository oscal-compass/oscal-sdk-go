/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	oscalValidation "github.com/defenseunicorns/go-oscal/src/pkg/validation"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
)

var _ Validator = (*SchemaValidator)(nil)

// OSCALVersion is the default version of OSCAL supported.
const OSCALVersion = "1.1.3"

/*
SchemaValidator implementation a validation.Validator and
validates the OSCAL documents with OSCAL JSON schema using `go-oscal`.
*/
type SchemaValidator struct {
	id           string
	oscalVersion string
}

// NewSchemaValidator returns a new SchemaValidator with a
// default OSCAL Version.
func NewSchemaValidator() *SchemaValidator {
	return NewSchemaValidatorWithVersion(OSCALVersion)
}

// NewSchemaValidatorWithVersion returns a new versioned SchemaValidator.
func NewSchemaValidatorWithVersion(version string) *SchemaValidator {
	return &SchemaValidator{
		id:           "schema",
		oscalVersion: version,
	}
}

func (s *SchemaValidator) Validate(modelData oscalTypes.OscalModels) error {
	validator, err := oscalValidation.NewValidatorDesiredVersion(modelData, s.oscalVersion)
	if err != nil {
		return &ErrValidation{Type: s.id, Model: "", Err: err}
	}

	modelType := validator.GetModelType()
	err = validator.Validate()
	if err != nil {
		return &ErrValidation{Type: s.id, Model: modelType, Err: err}
	}

	return nil
}
