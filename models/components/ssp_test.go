/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package components

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/validation"
)

func TestSystemComponentAdapter(t *testing.T) {
	testDataPath := filepath.Join("../../testdata", "test-ssp.json")

	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	ssp, err := models.NewSystemSecurityPlan(file, validation.NoopValidator{})
	require.NoError(t, err)
	require.NotNil(t, ssp)

	require.Len(t, ssp.SystemImplementation.Components, 3)
	adapter := NewSystemComponentAdapter(ssp.SystemImplementation.Components[0])
	require.Equal(t, "Example Service", adapter.Title())
	require.Equal(t, Service, adapter.Type())
	require.Equal(t, "4e19131e-b361-4f0e-8262-02bf4456202e", adapter.UUID())
	require.Len(t, adapter.Props(), 7)
	systemComp, ok := adapter.AsSystemComponent()
	require.True(t, ok)
	require.Equal(t, adapter.UUID(), systemComp.UUID)
	definedComp, ok := adapter.AsDefinedComponent()
	require.True(t, ok)
	require.Equal(t, adapter.UUID(), definedComp.UUID)
}

func TestControlImplementationAdapter(t *testing.T) {
	testDataPath := filepath.Join("../../testdata", "test-ssp.json")

	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	ssp, err := models.NewSystemSecurityPlan(file, validation.NoopValidator{})
	require.NoError(t, err)
	require.NotNil(t, ssp)

	adapter := NewControlImplementationAdapter(ssp.ControlImplementation)
	require.Len(t, adapter.Requirements(), 2)
	require.Len(t, adapter.SetParameters(), 0)
	require.Len(t, adapter.Props(), 0)

	impReq := adapter.Requirements()[0]
	require.Len(t, impReq.SetParameters(), 0)
	require.Len(t, impReq.Props(), 1)
	require.Equal(t, "db7b97db-dadc-4afd-850a-245ca09cb811", impReq.UUID())
	require.Equal(t, "ex-1", impReq.ControlID())
	require.Len(t, impReq.Statements(), 1)

	statement := impReq.Statements()[0]
	require.Len(t, statement.Props(), 1)
	require.Equal(t, "7ad47329-dc55-4196-a19d-178a8fe7438e", statement.UUID())
	require.Equal(t, "ex-1_smt", statement.StatementID())
}
