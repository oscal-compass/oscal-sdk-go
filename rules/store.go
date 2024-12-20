/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package rules

import (
	"context"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
)

// Store provides methods for filtering and searching RuleSets generated from OSCAL rule/check extensions.
type Store interface {
	// GetByRuleID returns the RuleSet associated with the given rule ID.
	GetByRuleID(ctx context.Context, ruleID string) (extensions.RuleSet, error)
	// GetByCheckID returns the RuleSet associated with the given check ID.
	GetByCheckID(ctx context.Context, checkID string) (extensions.RuleSet, error)
	// FindByComponent returns RuleSets associated with the component ID.
	//
	// For validation components, only relevant checks are returned.
	// For non-validation or "target" components, all information is returned.
	FindByComponent(ctx context.Context, componentId string) ([]extensions.RuleSet, error)
}
