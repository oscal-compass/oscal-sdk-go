/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package rules

import (
	"strings"

	oscal112 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
)

// groupPropsByRemarks will return the properties group by the same
// remark string. This is how properties are grouped to create rule sets.
func groupPropsByRemarks(props []oscal112.Property) map[string]set.Set[oscal112.Property] {
	grouped := map[string]set.Set[oscal112.Property]{}
	for _, prop := range props {
		if prop.Remarks == "" {
			continue
		}
		remarks := prop.Remarks
		propSet, ok := grouped[remarks]
		if !ok {
			propSet = set.New[oscal112.Property]()
		}
		propSet.Add(prop)
		grouped[remarks] = propSet
	}
	return grouped
}

// getProp finds a property in a set by the property name. This also implicitly checks the property is a
// trestle-defined property in the namespace.
func getProp(name string, props set.Set[oscal112.Property]) (oscal112.Property, bool) {
	for prop := range props {
		if prop.Name == name && strings.Contains(prop.Ns, extensions.TrestleNameSpace) {
			return prop, true
		}
	}
	return oscal112.Property{}, false
}
