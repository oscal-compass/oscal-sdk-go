/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package extensions

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/stretchr/testify/require"
)

func TestGetTrestleProp(t *testing.T) {
	tests := []struct {
		name       string
		inputProps []oscalTypes.Property
		inputName  string
		wantProp   oscalTypes.Property
		wantFound  bool
	}{
		{
			name:      "Valid/PropFound",
			inputName: "testProp1",
			inputProps: []oscalTypes.Property{
				{
					Name:  "testProp1",
					Value: "testValue",
				},
				{
					Name:  "testProp1",
					Value: "testValue",
					Ns:    TrestleNameSpace,
				},
			},
			wantProp: oscalTypes.Property{
				Name:    "testProp1",
				Value:   "testValue",
				Ns:      TrestleNameSpace,
				Group:   "",
				Class:   "",
				Remarks: "",
			},
			wantFound: true,
		},
		{
			name:      "Valid/PropNotFound",
			inputName: "testProp",
			inputProps: []oscalTypes.Property{
				{
					Name:  "testProp1",
					Value: "testValue",
				},
				{
					Name:  "testProp2",
					Value: "testValue",
					Ns:    TrestleNameSpace,
				},
			},
			wantProp:  oscalTypes.Property{},
			wantFound: false,
		},
		{
			name:      "Valid/PropNotFoundNs",
			inputName: "testProp1",
			inputProps: []oscalTypes.Property{
				{
					Name:  "testProp1",
					Value: "testValue",
				},
				{
					Name:  "testProp2",
					Value: "testValue",
				},
			},
			wantProp:  oscalTypes.Property{},
			wantFound: false,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			foundProp, found := GetTrestleProp(c.inputName, c.inputProps)
			require.Equal(t, c.wantProp, foundProp)
			require.Equal(t, c.wantFound, found)
		})
	}
}

func TestFindAllProps(t *testing.T) {
	tests := []struct {
		name         string
		inputOptions []FindOption
		inputProps   []oscalTypes.Property
		wantProps    []oscalTypes.Property
	}{
		{
			name: "Valid/Defaults",
			inputProps: []oscalTypes.Property{
				{
					Name:  "testProp1",
					Value: "testValue1",
					Ns:    TrestleNameSpace,
				},
				{
					Name:  "testProp2",
					Value: "testValue2",
					Ns:    TrestleNameSpace,
				},
				{
					Name:  "testProp3",
					Value: "testValue3",
				},
			},
			wantProps: []oscalTypes.Property{
				{
					Name:  "testProp1",
					Value: "testValue1",
					Ns:    TrestleNameSpace,
				},
				{
					Name:  "testProp2",
					Value: "testValue2",
					Ns:    TrestleNameSpace,
				},
			},
		},
		{
			name: "Valid/PropsFoundByName",
			inputOptions: []FindOption{
				WithName("testProp1"),
			},
			inputProps: []oscalTypes.Property{
				{
					Name:  "testProp1",
					Value: "testValue1",
					Ns:    TrestleNameSpace,
				},
				{
					Name:  "testProp1",
					Value: "testValue2",
					Ns:    TrestleNameSpace,
				},
				{
					Name:  "testProp1",
					Value: "testValue3",
				},
			},
			wantProps: []oscalTypes.Property{
				{
					Name:    "testProp1",
					Value:   "testValue1",
					Ns:      TrestleNameSpace,
					Group:   "",
					Class:   "",
					Remarks: "",
				},
				{
					Name:    "testProp1",
					Value:   "testValue2",
					Ns:      TrestleNameSpace,
					Group:   "",
					Class:   "",
					Remarks: "",
				},
			},
		},
		{
			name: "Valid/NoPropsFound",
			inputOptions: []FindOption{
				WithName("testProp3"),
			},
			inputProps: []oscalTypes.Property{
				{
					Name:  "testProp1",
					Value: "testValue1",
					Ns:    TrestleNameSpace,
				},
				{
					Name:  "testProp1",
					Value: "testValue2",
					Ns:    TrestleNameSpace,
				},
				{
					Name:  "testProp1",
					Value: "testValue3",
					Ns:    TrestleNameSpace,
				},
			},
			wantProps: []oscalTypes.Property(nil),
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			foundProps := FindAllProps(c.inputProps, c.inputOptions...)
			require.Equal(t, c.wantProps, foundProps)
		})
	}
}
