/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package components

import oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

// Implementation is an interface representing a generic control implementation for
// a component that is present in an OSCAL Component Definition and an OSCAL SSP in
// different forms.
type Implementation interface {
	// Requirements returns a slice of implemented requirements associated to
	// the implementation.
	Requirements() []Requirement
	// SetParameters returns a list of OSCAL set-parameters associated with the implementation.
	SetParameters() []oscalTypes.SetParameter
	// Props returns a list of OSCAL properties associated with the implementation.
	Props() []oscalTypes.Property
}

// Requirement is an interface representing a generic implemented requirement
// for a component that is present in an OSCAL Component Definition and an OSCAL SSP in
// different forms.
type Requirement interface {
	// ControlID  returns the associated human-readable identifier for the requirement.
	ControlID() string
	// UUID returns the requirement assembly UUID
	UUID() string
	// SetParameters returns a list of OSCAL set-parameters associated with the requirement.
	SetParameters() []oscalTypes.SetParameter
	// Props returns a list of OSCAL properties associated with the requirement.
	Props() []oscalTypes.Property
	// Statements returns a slice of statements or implemented requirement parts associated
	// with an implemented requirement.
	Statements() []Statement
}

// Statement is an interface representing an implemented statement or
// sub-requirement for a component that is present in an OSCAL Component Definition and an OSCAL SSP in
// different forms.
type Statement interface {
	// StatementID returns the associated human-readable identifier for the statement.
	StatementID() string
	// UUID returns the statement assembly UUID
	UUID() string
	// Props returns a list of OSCAL properties associated with the statement.
	Props() []oscalTypes.Property
}
