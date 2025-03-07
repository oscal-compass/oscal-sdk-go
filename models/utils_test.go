/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package models

import (
	"testing"

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
