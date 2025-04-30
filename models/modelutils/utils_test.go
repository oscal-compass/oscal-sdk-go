/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package modelutils

import (
	"testing"
	"time"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/require"
)

func TestNilIfEmpty(t *testing.T) {
	tests := []struct {
		name      string
		slice     *[]string
		wantSlice *[]string
	}{
		{
			name:      "Valid/NilSlice",
			slice:     nil,
			wantSlice: nil,
		},
		{
			name:      "Valid/EmptySlice",
			slice:     &[]string{},
			wantSlice: nil,
		},
		{
			name:      "Valid/NonEmptySlice",
			slice:     &[]string{"test"},
			wantSlice: &[]string{"test"},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			gotSlice := NilIfEmpty(c.slice)
			require.Equal(t, c.wantSlice, gotSlice)
		})
	}
}

func TestFindValuesByName(t *testing.T) {
	tests := []struct {
		name  string
		model oscalTypes.OscalModels
	}{
		{
			name: "uuid",
			model: oscalTypes.OscalModels{
				ComponentDefinition: &oscalTypes.ComponentDefinition{
					Metadata: oscalTypes.Metadata{
						OscalVersion: "1.1.3",
						Version:      "0.1.0",
						LastModified: time.Now(),
					},
					UUID: "c14d8812-7098-4a9b-8f89-cba41b6ff0d8",
				},
			},
		},
	}
	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			values := FindValuesByName(&c.model, "UUID")
			require.Equal(t, []string{"c14d8812-7098-4a9b-8f89-cba41b6ff0d8"}, values)
		})
	}
}

func TestHasDuplicateValuesByName(t *testing.T) {
	tests := []struct {
		expectedValue bool
		name          string
		model         oscalTypes.OscalModels
	}{
		{
			expectedValue: false,
			name:          "uuid",
			model: oscalTypes.OscalModels{
				ComponentDefinition: &oscalTypes.ComponentDefinition{
					Metadata: oscalTypes.Metadata{
						OscalVersion: "1.1.3",
						Version:      "0.1.0",
						LastModified: time.Now(),
					},
					BackMatter: &oscalTypes.BackMatter{
						Resources: &[]oscalTypes.Resource{
							{
								UUID: "c14d8812-7098-4a9b-8f89-cba41b6ff0d8",
							},
						},
					},
					UUID: "c14d8812-7098-4a9b-8f89-cba41b6ff0d8",
				},
			},
		},
		{
			expectedValue: true,
			name:          "uuid",
			model: oscalTypes.OscalModels{
				ComponentDefinition: &oscalTypes.ComponentDefinition{
					Metadata: oscalTypes.Metadata{
						OscalVersion: "1.1.3",
						Version:      "0.1.0",
						LastModified: time.Now(),
					},
					BackMatter: &oscalTypes.BackMatter{
						Resources: &[]oscalTypes.Resource{
							{
								UUID: "c14d8812-xxxx-xxxx-xxxx-cba41b6ff0d8",
							},
						},
					},
					UUID: "c14d8812-7098-4a9b-8f89-cba41b6ff0d8",
				},
			},
		},
	}
	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			hasDupes := HasDuplicateValuesByName(&c.model, "UUID")
			require.Equal(t, c.expectedValue, hasDupes)
		})
	}
}
