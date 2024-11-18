// SPDX-License-Identifier: Apache-2.0

package rules

import (
	oscal112 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	. "github.com/oscal-compass/oscal-sdk-go/rules/internal"
)

// groupPropsByRemarks will return the properties group by the same
// remark string. This is how properties are grouped to create rule sets.
func groupPropsByRemarks(props []oscal112.Property) map[string]Set[oscal112.Property] {
	grouped := map[string]Set[oscal112.Property]{}
	for _, prop := range props {
		if prop.Remarks == "" {
			continue
		}
		remarks := prop.Remarks
		set, ok := grouped[remarks]
		if !ok {
			set = NewSet[oscal112.Property]()
		}
		set.Add(prop)
		grouped[remarks] = set
	}
	return grouped
}

// findProp finds a property in a set by the property name.
func findProp(name string, props Set[oscal112.Property]) (oscal112.Property, bool) {
	for prop := range props {
		if prop.Name == name {
			return prop, true
		}
	}
	return oscal112.Property{}, false
}
