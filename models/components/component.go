/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package components

import (
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
)

// ComponentConverter defines methods to convert and/or retrieve OSCAL Component underlying types.
type ComponentConverter interface {
	// AsDefinedComponent returns a DefinedComponent and whether the
	// retrieval or conversion was successful.
	AsDefinedComponent() (oscalTypes.DefinedComponent, bool)
	// AsSystemComponent returns a SystemComponent and whether the
	// retrieve or conversion was successful.
	AsSystemComponent() (oscalTypes.SystemComponent, bool)
}

// Component defines methods to retrieve common information from OSCAL Component types.
type Component interface {
	// Title returns the title associated with the component
	Title() string
	// Type returns the type of component
	Type() ComponentType
	// UUID returns the component UUID
	UUID() string
	// Props returns a list of OSCAL properties associated with the component.
	Props() []oscalTypes.Property
	ComponentConverter
}

// ComponentType represents valid types of Components in OSCAL.
// Reference: https://pages.nist.gov/OSCAL-Reference/models/v1.1.3/assessment-plan/json-reference/#/assessment-plan/local-definitions/components/type
type ComponentType string

const (
	Validation       ComponentType = "validation"
	Software         ComponentType = "software"
	Service          ComponentType = "service"
	Interconnection  ComponentType = "interconnection"
	ThisSystem       ComponentType = "this-system"
	System           ComponentType = "system"
	Hardware         ComponentType = "hardware"
	Policy           ComponentType = "policy"
	Physical         ComponentType = "physical"
	ProcessProcedure ComponentType = "process-procedure"
	Plan             ComponentType = "plan"
	Guidance         ComponentType = "guidance"
	Standard         ComponentType = "standard"
	Network          ComponentType = "network"
)
