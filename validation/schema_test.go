/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"
	"time"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/require"
)

func TestNewSchemaValidator(t *testing.T) {
	tests := []struct {
		name       string
		modelData  oscalTypes.OscalModels
		wantResult bool
		wantErr    string
	}{
		{
			name: "Valid/ValidComponentDefinition",
			modelData: oscalTypes.OscalModels{
				ComponentDefinition: &oscalTypes.ComponentDefinition{
					Metadata: oscalTypes.Metadata{
						OscalVersion: "1.1.3",
						Version:      "0.1.0",
						LastModified: time.Now(),
					},
					UUID: "c14d8812-7098-4a9b-8f89-cba41b6ff0d8",
				},
			},
			wantResult: true,
		},
		{
			name: "Invalid/CatalogUUID",
			modelData: oscalTypes.OscalModels{
				Catalog: &oscalTypes.Catalog{
					UUID: "not-a-uuid",
					Metadata: oscalTypes.Metadata{
						OscalVersion: "1.1.3",
						Version:      "0.1.0",
						LastModified: time.Now(),
					},
				},
			},
			wantResult: false,
			wantErr:    "not-a-uuid",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			validator := NewSchemaValidator()
			require.Equal(t, OSCALVersion, validator.oscalVersion)
			require.Equal(t, "schema", validator.id)
			err := validator.Validate(c.modelData)
			if c.wantErr != "" {
				require.ErrorContains(t, err, c.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
