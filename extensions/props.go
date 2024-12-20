/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package extensions

import (
	"strings"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
)

// TrestleNameSpace is the generic namespace for trestle-defined property extensions.
const TrestleNameSpace = "https://oscal-compass.github.io/compliance-trestle/schemas/oscal"

// Below are defined oscal.Property names for compass-based extensions.
const (
	// RuleIdProp represents the property name for Rule ids.
	RuleIdProp = "Rule_Id"
	// RuleDescriptionProp represents the property name for Rule descriptions.
	RuleDescriptionProp = "Rule_Description"
	// CheckIdProp represents the property name for Check ids.
	CheckIdProp = "Check_Id"
	// CheckDescriptionProp represents the property name for Check descriptions.
	CheckDescriptionProp = "Check_Description"
	// ParameterIdProp represents the property name for Parameter ids.
	ParameterIdProp = "Parameter_Id"
	// ParameterDescriptionProp represents the property name for Parameter descriptions.
	ParameterDescriptionProp = "Parameter_Description"
	// ParameterDefaultProp represents the property name for Parameter default selected values.
	ParameterDefaultProp = "Parameter_Value_Default"
	// FrameworkProp represents the property name for the control source short name.
	FrameworkProp = "Framework_Short_Name"
)

// FindAllProps returns all properties with the given name. If no properties match, nil is returned.
// This function also implicitly checks that the property is a trestle-defined property in the namespace.
func FindAllProps(name string, props []oscalTypes.Property) []oscalTypes.Property {
	var matchingProps []oscalTypes.Property
	for _, prop := range props {
		if prop.Name == name && strings.Contains(prop.Ns, TrestleNameSpace) {
			matchingProps = append(matchingProps, prop)
		}
	}
	return matchingProps
}

// GetTrestleProp returned  the first property matching the given name and a match is found.
// This function also implicitly checks that the property is a trestle-defined property in the namespace.
func GetTrestleProp(name string, props []oscalTypes.Property) (oscalTypes.Property, bool) {
	for _, prop := range props {
		if prop.Name == name && strings.Contains(prop.Ns, TrestleNameSpace) {
			return prop, true
		}
	}
	return oscalTypes.Property{}, false
}
