/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"os"
	"path/filepath"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/validation"
)

func TestGetFrameworkShortName(t *testing.T) {
	tests := []struct {
		name                string
		inputImplementation oscalTypes.ControlImplementationSet
		wantName            string
		wantFound           bool
	}{
		{
			name: "Valid/ShortNameFromProps",
			inputImplementation: oscalTypes.ControlImplementationSet{
				Props: &[]oscalTypes.Property{
					{
						Name:  extensions.FrameworkProp,
						Value: "propFramework",
						Ns:    extensions.TrestleNameSpace,
					},
				},
				Source: "profiles/framework/profile.json",
			},
			wantName:  "propFramework",
			wantFound: true,
		},
		{
			name: "Valid/ShortNameFromSource",
			inputImplementation: oscalTypes.ControlImplementationSet{
				Source: "profiles/sourceFramework/profile.json",
			},
			wantName:  "sourceFramework",
			wantFound: true,
		},
		{
			name:                "Valid/NoShortName",
			inputImplementation: oscalTypes.ControlImplementationSet{},
			wantName:            "",
			wantFound:           false,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			name, found := GetFrameworkShortName(c.inputImplementation)
			require.Equal(t, c.wantName, name)
			require.Equal(t, c.wantFound, found)
		})
	}
}

func TestByFramework(t *testing.T) {
	testDataPath := filepath.Join("../testdata", "component-definition-test-reqs.json")
	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	definition, err := models.NewComponentDefinition(file, validation.NoopValidator{})
	require.NoError(t, err)

	require.NotNil(t, definition.Components)

	var allImplementations []oscalTypes.ControlImplementationSet
	for _, component := range *definition.Components {
		if component.ControlImplementations == nil {
			continue
		}
		allImplementations = append(allImplementations, *component.ControlImplementations...)
	}

	implementationsMap, framework, err := ByFramework("cis", allImplementations)
	require.NoError(t, err)
	expectedSettings := &ImplementationSettings{
		settings: Settings{
			mappedRules: set.Set[string]{
				"etcd_cert_file": struct{}{},
				"etcd_key_file":  struct{}{},
			},
			selectedParameters: map[string]string{},
		},
		implementedReqSettings: map[string]Settings{
			"CIS-2.1": {
				mappedRules: set.Set[string]{
					"etcd_cert_file": struct{}{},
					"etcd_key_file":  struct{}{},
				},
				selectedParameters: map[string]string{},
			},
		},
		controlsByRules: map[string]set.Set[string]{
			"etcd_cert_file": {
				"CIS-2.1": struct{}{},
			},
			"etcd_key_file": {
				"CIS-2.1": struct{}{},
			},
		},
		controlsById: map[string]oscalTypes.AssessedControlsSelectControlById{
			"CIS-2.1": {
				ControlId: "CIS-2.1",
			},
		},
	}

	expectedFramework := FrameworkSource{
		Title:       "cis",
		Description: "CIS Profile",
		Href:        "profiles/cis/profile.json",
	}

	require.Equal(t, expectedFramework, framework)
	require.Equal(t, expectedSettings, implementationsMap)

	_, _, err = ByFramework("doesnotexist", allImplementations)
	require.EqualError(t, err, "framework doesnotexist is not in control implementations")
}
