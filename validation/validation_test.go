/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"testing"
	"time"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/require"
)

// Test validator that produces no errors.
type TestValidator struct{}

func (v TestValidator) Validate(models oscalTypes.OscalModels) error {
	return nil
}

// Test validator that produces errors.
type TestValidatorWithErr struct{}

func (v TestValidatorWithErr) Validate(models oscalTypes.OscalModels) error {
	return fmt.Errorf(("test error"))
}

func TestValidateAll(t *testing.T) {

	var modelData = oscalTypes.OscalModels{
		ComponentDefinition: &oscalTypes.ComponentDefinition{
			Metadata: oscalTypes.Metadata{
				OscalVersion: "1.1.3",
				Version:      "0.1.0",
				LastModified: time.Now(),
			},
			UUID: "c14d8812-7098-4a9b-8f89-cba41b6ff0d8",
		},
	}

	tests := []struct {
		name       string
		validators []Validator
		wantErr    bool
	}{
		{
			name:       "Valid/Success",
			validators: []Validator{TestValidator{}, TestValidator{}},
			wantErr:    false,
		},
		{
			name:       "Invalid/WithErrors",
			validators: []Validator{TestValidator{}, TestValidatorWithErr{}},
			wantErr:    true,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			validateAll := ValidateAll(c.validators...)
			err := validateAll(modelData)
			if !c.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
