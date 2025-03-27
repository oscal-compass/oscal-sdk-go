/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package components

import oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

// Interface checks for the implementations for SSP structures
var (
	_ Component      = (*SystemComponentAdapter)(nil)
	_ Implementation = (*ControlImplementationAdapter)(nil)
	_ Requirement    = (*ImplementedRequirementAdapter)(nil)
	_ Statement      = (*StatementAdapter)(nil)
)

// SystemComponentAdapter wraps an OSCAL SystemComponent to
// provide methods for compatibility with Component.
type SystemComponentAdapter struct {
	systemComp oscalTypes.SystemComponent
}

// NewSystemComponentAdapter returns an initialized SystemComponentAdapter from a given
// SystemComponent.
func NewSystemComponentAdapter(systemComponent oscalTypes.SystemComponent) *SystemComponentAdapter {
	return &SystemComponentAdapter{
		systemComp: systemComponent,
	}
}

func (s *SystemComponentAdapter) UUID() string {
	return s.systemComp.UUID
}

func (s *SystemComponentAdapter) Title() string {
	return s.systemComp.Title
}

func (s *SystemComponentAdapter) Type() ComponentType {
	return ComponentType(s.systemComp.Type)
}

func (s *SystemComponentAdapter) AsDefinedComponent() (oscalTypes.DefinedComponent, bool) {
	return oscalTypes.DefinedComponent{
		Description:      s.systemComp.Description,
		Links:            s.systemComp.Links,
		Props:            s.systemComp.Props,
		Protocols:        s.systemComp.Protocols,
		Purpose:          s.systemComp.Purpose,
		Remarks:          s.systemComp.Remarks,
		ResponsibleRoles: s.systemComp.ResponsibleRoles,
		Title:            s.systemComp.Title,
		Type:             s.systemComp.Type,
		UUID:             s.systemComp.UUID,
	}, true
}

func (s *SystemComponentAdapter) AsSystemComponent() (oscalTypes.SystemComponent, bool) {
	return s.systemComp, true
}

func (s *SystemComponentAdapter) Props() []oscalTypes.Property {
	if s.systemComp.Props == nil {
		return []oscalTypes.Property{}
	}
	return *s.systemComp.Props
}

// ControlImplementationAdapter wraps an OSCAL ControlImplementation to provide
// methods for compatibility with Implementation.
type ControlImplementationAdapter struct {
	controlImp oscalTypes.ControlImplementation
}

// NewControlImplementationAdapter returns an initialized ControlImplementationAdapter from a given
// ControlImplementation from an OSCAL SSP.
func NewControlImplementationAdapter(controlImp oscalTypes.ControlImplementation) *ControlImplementationAdapter {
	return &ControlImplementationAdapter{
		controlImp: controlImp,
	}
}

func (c *ControlImplementationAdapter) Requirements() []Requirement {
	var requirements []Requirement
	for _, requirement := range c.controlImp.ImplementedRequirements {
		requirementAdapter := NewImplementedRequirementAdapter(requirement)
		requirements = append(requirements, requirementAdapter)
	}
	return requirements
}

func (c *ControlImplementationAdapter) SetParameters() []oscalTypes.SetParameter {
	if c.controlImp.SetParameters == nil {
		return []oscalTypes.SetParameter{}
	}
	return *c.controlImp.SetParameters
}

func (c *ControlImplementationAdapter) Props() []oscalTypes.Property {
	// TODO: Where does this go in the SSP?
	return []oscalTypes.Property{}
}

// ImplementedRequirementAdapter wraps an OSCAL ImplementedRequirement to provide
// methods for compatibility with Requirement.
type ImplementedRequirementAdapter struct {
	impReq oscalTypes.ImplementedRequirement
}

// NewImplementedRequirementAdapter returns an initialized ImplementedRequirementAdapter from a given
// ImplementedRequirement from an OSCAL SSP.
func NewImplementedRequirementAdapter(impReq oscalTypes.ImplementedRequirement) *ImplementedRequirementAdapter {
	return &ImplementedRequirementAdapter{
		impReq: impReq,
	}
}

func (i *ImplementedRequirementAdapter) ControlID() string {
	return i.impReq.ControlId
}

func (i *ImplementedRequirementAdapter) UUID() string {
	return i.impReq.UUID
}

func (i *ImplementedRequirementAdapter) SetParameters() []oscalTypes.SetParameter {
	if i.impReq.SetParameters == nil {
		return []oscalTypes.SetParameter{}
	}
	return *i.impReq.SetParameters
}

func (i *ImplementedRequirementAdapter) Props() []oscalTypes.Property {
	var oscalProps []oscalTypes.Property
	if i.impReq.Props != nil {
		oscalProps = append(oscalProps, *i.impReq.Props...)
	}

	if i.impReq.ByComponents == nil {

		return oscalProps
	}

	for _, byComp := range *i.impReq.ByComponents {
		if byComp.Props != nil {
			oscalProps = append(oscalProps, *byComp.Props...)
		}
	}
	return oscalProps
}

func (i *ImplementedRequirementAdapter) Statements() []Statement {
	var statements []Statement
	if i.impReq.Statements == nil {
		return statements
	}
	for _, stm := range *i.impReq.Statements {
		stmAdapter := NewStatementAdapter(stm)
		statements = append(statements, stmAdapter)
	}

	return statements
}

// StatementAdapter wraps an OSCAL Statement to provide
// methods for compatibility with Statement.
type StatementAdapter struct {
	stm oscalTypes.Statement
}

// NewStatementAdapter returns an initialized StatementAdapter from a given
// Statement from an OSCAL SSP.
func NewStatementAdapter(statement oscalTypes.Statement) *StatementAdapter {
	return &StatementAdapter{
		stm: statement,
	}
}

func (s *StatementAdapter) StatementID() string {
	return s.stm.StatementId
}

func (s *StatementAdapter) UUID() string {
	return s.stm.UUID
}

func (s *StatementAdapter) Props() []oscalTypes.Property {
	var oscalProps []oscalTypes.Property
	if s.stm.Props != nil {
		oscalProps = append(oscalProps, *s.stm.Props...)
	}

	if s.stm.ByComponents == nil {
		return oscalProps
	}

	for _, byComp := range *s.stm.ByComponents {
		if byComp.Props != nil {
			oscalProps = append(oscalProps, *byComp.Props...)
		}
	}
	return oscalProps
}
