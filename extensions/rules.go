/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package extensions

// RuleSet defines a Rule instance with associated
// Check implementation data.
type RuleSet struct {
	// Rule is a single rule instance associated with the set.
	Rule Rule
	// Checks include all associated check data registered for the rule.
	Checks []Check
}

// Rule defines a single compliance rule which can also be defined
// as a technical control or a way to validate implemented requirements.
type Rule struct {
	// ID is a string representation of the rule identifier.
	ID string
	// Description defines description of what the rule does.
	Description string
	// Parameters are optional information for tuning rule logic.
	Parameters []Parameter
}

// Check defines a single concrete implementation of a Rule.
type Check struct {
	// ID is a string representation of the check identifier.
	ID string
	// Description defines description of what the check does.
	Description string
}

// Parameter identifies a parameter or variable that can be used to alter rule logic.
type Parameter struct {
	// ID is a string representation of the parameter identifier.
	ID string
	// Description defines description of what the parameter does.
	Description string
	// Value is the selected value or option for the parameter.
	Value string
}
