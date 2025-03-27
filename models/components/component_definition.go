/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package components

import oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

// Interface checks for the implementation for Component Definition structures
var (
	_ Component      = (*DefinedComponentAdapter)(nil)
	_ Implementation = (*ControlImplementationSetAdapter)(nil)
	_ Requirement    = (*ImplementedRequirementImplementationAdapter)(nil)
	_ Statement      = (*ControlStatementAdapter)(nil)
)

const defaultState = "operational"

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
			State: defaultState,
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

// ControlImplementationSetAdapter wraps an OSCAL ControlImplementationSet to provide
// methods for compatibility with Implementation.
type ControlImplementationSetAdapter struct {
	controlImp oscalTypes.ControlImplementationSet
}

// NewControlImplementationSetAdapter returns an initialized ControlImplementationAdapterSet from a given
// ControlImplementation from an OSCAL Component Definition.
func NewControlImplementationSetAdapter(controlImp oscalTypes.ControlImplementationSet) *ControlImplementationSetAdapter {
	return &ControlImplementationSetAdapter{
		controlImp: controlImp,
	}
}

func (c *ControlImplementationSetAdapter) Requirements() []Requirement {
	var requirements []Requirement
	for _, requirement := range c.controlImp.ImplementedRequirements {
		requirementAdapter := NewImplementedRequirementImplementationAdapter(requirement)
		requirements = append(requirements, requirementAdapter)
	}
	return requirements
}

func (c *ControlImplementationSetAdapter) SetParameters() []oscalTypes.SetParameter {
	if c.controlImp.SetParameters == nil {
		return []oscalTypes.SetParameter{}
	}
	return *c.controlImp.SetParameters
}

func (c *ControlImplementationSetAdapter) Props() []oscalTypes.Property {
	if c.controlImp.Props == nil {
		return []oscalTypes.Property{}
	}
	return *c.controlImp.Props
}

// ImplementedRequirementImplementationAdapter wraps an OSCAL ImplementedRequirementImplementation to provide
// methods for compatibility with Requirement.
type ImplementedRequirementImplementationAdapter struct {
	impReq oscalTypes.ImplementedRequirementControlImplementation
}

// NewImplementedRequirementImplementationAdapter returns an initialized ImplementedRequirementImplementationAdapter from a given
// ImplementedRequirementImplementation from an OSCAL Component Definition.
func NewImplementedRequirementImplementationAdapter(impReq oscalTypes.ImplementedRequirementControlImplementation) *ImplementedRequirementImplementationAdapter {
	return &ImplementedRequirementImplementationAdapter{
		impReq: impReq,
	}
}

func (i *ImplementedRequirementImplementationAdapter) ControlID() string {
	return i.impReq.ControlId
}

func (i *ImplementedRequirementImplementationAdapter) UUID() string {
	return i.impReq.UUID
}

func (i *ImplementedRequirementImplementationAdapter) SetParameters() []oscalTypes.SetParameter {
	if i.impReq.SetParameters == nil {
		return []oscalTypes.SetParameter{}
	}
	return *i.impReq.SetParameters
}

func (i *ImplementedRequirementImplementationAdapter) Props() []oscalTypes.Property {
	if i.impReq.Props == nil {
		return []oscalTypes.Property{}
	}
	return *i.impReq.Props
}

func (i *ImplementedRequirementImplementationAdapter) Statements() []Statement {
	var statements []Statement
	if i.impReq.Statements == nil {
		return statements
	}
	for _, stm := range *i.impReq.Statements {
		stmAdapter := NewControlStatementAdapter(stm)
		statements = append(statements, stmAdapter)
	}

	return statements
}

// ControlStatementAdapter wraps an OSCAL ControlStatement to provide
// methods for compatibility with Statement.
type ControlStatementAdapter struct {
	stm oscalTypes.ControlStatementImplementation
}

// NewControlStatementAdapter returns an initialized ControlStatementAdapter from a given
// ControlStatement from an OSCAL Component Definition.
func NewControlStatementAdapter(statement oscalTypes.ControlStatementImplementation) *ControlStatementAdapter {
	return &ControlStatementAdapter{
		stm: statement,
	}
}

func (c *ControlStatementAdapter) StatementID() string {
	return c.stm.StatementId
}

func (c *ControlStatementAdapter) UUID() string {
	return c.stm.UUID
}

func (c *ControlStatementAdapter) Props() []oscalTypes.Property {
	if c.stm.Props == nil {
		return []oscalTypes.Property{}
	}
	return *c.stm.Props
}
