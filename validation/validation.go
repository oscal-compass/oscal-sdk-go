/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

// Package validation defines logic for validation OSCAL against OSCAL schema and custom validation for
// supported extensions.
package validation

import (
	"errors"
	"fmt"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
)

// Validator defines methods for semantic validation on a decoded OSCAL models.
// Structural validation is done through the decoding process.
type Validator interface {
	// Validate will take input OSCAL Models to perform validation.
	// Returns an error if data is invalid.
	Validate(oscalTypes.OscalModels) error
}

// ErrValidation is returned when data is not valid.
type ErrValidation struct {
	// Type returns the type of validator the error occurred on.
	Type string
	// Model returns the model type the error occurred on.
	Model string
	// Err return the error message.
	Err error
}

func (e *ErrValidation) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Err.Error())
}

var _ Validator = (*NoopValidator)(nil)

// NoopValidator for skipping validation.
type NoopValidator struct{}

func (n NoopValidator) Validate(_ oscalTypes.OscalModels) error {
	return nil
}

// ValidatorFunc implements the Validator interface
type ValidatorFunc func(models oscalTypes.OscalModels) error

func (f ValidatorFunc) Validate(models oscalTypes.OscalModels) error {
	return f(models)
}

// ValidateAll returns a func that will run multiple validators in sequence.
func ValidateAll(validators ...Validator) ValidatorFunc {
	return func(models oscalTypes.OscalModels) error {
		var valErrors []error
		for _, v := range validators {
			err := v.Validate(models)
			if err != nil {
				valErrors = append(valErrors, err)
			}
		}
		if len(valErrors) > 0 {
			return errors.Join(valErrors...)
		}
		return nil
	}
}
