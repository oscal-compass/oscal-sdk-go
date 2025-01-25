/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package components

import (
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/generators"
)

var _ Component = (*DefinedComponentAdapter)(nil)

// DefinedComponentAdapter wrapped an OSCAL DefinedComponent to
// provide methods for compatibility with Component.
type DefinedComponentAdapter struct {
	definedComp oscalTypes.DefinedComponent
}

// NewDefinedComponentAdapter returns an initialized DefinedComponentAdapter from a given
// DefinedComponent.
func NewDefinedComponentAdapter(definedComponent oscalTypes.DefinedComponent) *DefinedComponentAdapter {
	return &DefinedComponentAdapter{
		definedComp: definedComponent,
	}
}

func (d *DefinedComponentAdapter) UUID() string {
	return d.definedComp.UUID
}

func (d *DefinedComponentAdapter) Title() string {
	return d.definedComp.Title
}

func (d *DefinedComponentAdapter) Type() ComponentType {
	return ComponentType(d.definedComp.Type)
}

func (d *DefinedComponentAdapter) AsDefinedComponent() (oscalTypes.DefinedComponent, bool) {
	return d.definedComp, true
}

func (d *DefinedComponentAdapter) AsSystemComponent() (oscalTypes.SystemComponent, bool) {
	return oscalTypes.SystemComponent{
		Description:      d.definedComp.Description,
		Links:            d.definedComp.Links,
		Props:            d.definedComp.Props,
		Protocols:        d.definedComp.Protocols,
		Purpose:          d.definedComp.Purpose,
		Remarks:          d.definedComp.Remarks,
		ResponsibleRoles: d.definedComp.ResponsibleRoles,
		Status: oscalTypes.SystemComponentStatus{
			State: generators.SampleRequiredString,
		},
		Title: d.definedComp.Title,
		Type:  d.definedComp.Type,
		UUID:  d.definedComp.UUID,
	}, true
}

func (d *DefinedComponentAdapter) Props() []oscalTypes.Property {
	if d.definedComp.Props == nil {
		return []oscalTypes.Property{}
	}
	return *d.definedComp.Props
}
