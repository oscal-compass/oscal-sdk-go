package components

import (
	oscaltypes112 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/models"
)

var _ Component = (*DefinedComponentAdapter)(nil)

// DefinedComponentAdapter wrapped an OSCAL DefinedComponent to
// provide methods for compatibility with Component.
type DefinedComponentAdapter struct {
	definedComp oscaltypes112.DefinedComponent
}

// NewDefinedComponentAdapter returns an initialized DefinedComponentAdapter from a given
// DefinedComponent.
func NewDefinedComponentAdapter(definedComponent oscaltypes112.DefinedComponent) *DefinedComponentAdapter {
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

func (d *DefinedComponentAdapter) AsDefinedComponent() (oscaltypes112.DefinedComponent, bool) {
	return d.definedComp, true
}

func (d *DefinedComponentAdapter) AsSystemComponent() (oscaltypes112.SystemComponent, bool) {
	return oscaltypes112.SystemComponent{
		Description:      d.definedComp.Description,
		Links:            d.definedComp.Links,
		Props:            d.definedComp.Props,
		Protocols:        d.definedComp.Protocols,
		Purpose:          d.definedComp.Purpose,
		Remarks:          d.definedComp.Remarks,
		ResponsibleRoles: d.definedComp.ResponsibleRoles,
		Status: oscaltypes112.SystemComponentStatus{
			State: models.DefaultRequiredString,
		},
		Title: d.definedComp.Title,
		Type:  d.definedComp.Type,
		UUID:  d.definedComp.UUID,
	}, true
}

func (d *DefinedComponentAdapter) Props() []oscaltypes112.Property {
	if d.definedComp.Props == nil {
		return []oscaltypes112.Property{}
	}
	return *d.definedComp.Props
}
