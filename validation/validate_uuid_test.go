/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
)

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name    string
		model   oscalTypes.OscalModels
		wantErr bool
	}{
		{
			name: "valid component definition with unique UUIDs",
			model: oscalTypes.OscalModels{
				ComponentDefinition: &oscalTypes.ComponentDefinition{
					UUID: "uuid-1",
					Components: &[]oscalTypes.DefinedComponent{
						{UUID: "uuid-2"},
						{UUID: "uuid-3"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid component definition with duplicate UUIDs",
			model: oscalTypes.OscalModels{
				ComponentDefinition: &oscalTypes.ComponentDefinition{
					UUID: "uuid-1",
					Components: &[]oscalTypes.DefinedComponent{
						{UUID: "uuid-1"},
						{UUID: "uuid-3"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "component definition with empty UUID",
			model: oscalTypes.OscalModels{
				ComponentDefinition: &oscalTypes.ComponentDefinition{
					UUID: "",
					Components: &[]oscalTypes.DefinedComponent{
						{UUID: "uuid-1"},
						{UUID: "uuid-3"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "component definition with multiple empty UUID",
			model: oscalTypes.OscalModels{
				ComponentDefinition: &oscalTypes.ComponentDefinition{
					UUID: "",
					Components: &[]oscalTypes.DefinedComponent{
						{UUID: ""},
						{UUID: "uuid-3"},
					},
				},
			},
			wantErr: true,
		},
	}
	validator := UuidValidator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("UuidValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
