/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package rules

import (
	"strings"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
)

// groupPropsByRemarks will return the properties group by the same
// remark string. This is how properties are grouped to create rule sets.
func groupPropsByRemarks(props []oscalTypes.Property) map[string]set.Set[oscalTypes.Property] {
	grouped := map[string]set.Set[oscalTypes.Property]{}
	for _, prop := range props {
		if prop.Remarks == "" {
			continue
		}
		remarks := prop.Remarks
		propSet, ok := grouped[remarks]
		if !ok {
			propSet = set.New[oscalTypes.Property]()
		}
		propSet.Add(prop)
		grouped[remarks] = propSet
	}
	return grouped
}

// getProp finds a property in a set by the property name. This also implicitly checks the property is a
// trestle-defined property in the namespace.
func getProp(name string, props set.Set[oscalTypes.Property]) (oscalTypes.Property, bool) {
	for prop := range props {
		if prop.Name == name && strings.Contains(prop.Ns, extensions.TrestleNameSpace) {
			return prop, true
		}
	}
	return oscalTypes.Property{}, false
}
