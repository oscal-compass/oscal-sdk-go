/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package components

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/generators"
)

func TestDefinedComponentAdapter(t *testing.T) {
	testDataPath := filepath.Join("../../testdata", "component-definition-test.json")

	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	definition, err := generators.NewComponentDefinition(file)
	require.NoError(t, err)
	require.NotNil(t, definition)
	require.NotNil(t, definition.Components)
	comps := *definition.Components
	require.Len(t, comps, 3)
	adapter := NewDefinedComponentAdapter(comps[0])
	require.Equal(t, "TestKubernetes", adapter.Title())
	require.Equal(t, Service, adapter.Type())
	require.Equal(t, "c8106bc8-5174-4e86-91a4-52f2fe0ed027", adapter.UUID())
	require.Len(t, adapter.Props(), 6)
	systemComp, ok := adapter.AsSystemComponent()
	require.True(t, ok)
	require.Equal(t, adapter.UUID(), systemComp.UUID)
	definedComp, ok := adapter.AsDefinedComponent()
	require.True(t, ok)
	require.Equal(t, adapter.UUID(), definedComp.UUID)
}
