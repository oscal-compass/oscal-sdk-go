// SPDX-License-Identifier: Apache-2.0

package extensions

// Property names for trestle-based extensions
const (
	RuleIdProp               = "Rule_Id"
	RuleDescriptionProp      = "Rule_Description"
	CheckIdProp              = "Check_Id"
	CheckDescriptionProp     = "Check_Description"
	ParameterIdProp          = "Parameter_Id"
	ParameterDescriptionProp = "Parameter_Description"
	ParameterDefaultProp     = "Parameter_Value_Default"
)

type RuleSet struct {
	// A single instance of a rule
	Rule Rule
	// A set of automation information called "checks"
	Checks []Check
}

// Rule defines a single rule with all associated metadata
type Rule struct {
	// Rule identification
	ID string
	// High level description
	RuleDescription string
	// Associated rule parameter information for tuning options.
	Parameter *Parameter
}

// Check defines a single check with all associated metadata.
type Check struct {
	// Associated check implementation identification
	ID string
	// High level check description
	Description string
}

// Parameter identifies a parameter or variable that can be used to alter rule logic
type Parameter struct {
	// Parameter Identification
	ID          string
	Description string
	// The selected value for the parameter
	Value string
}
