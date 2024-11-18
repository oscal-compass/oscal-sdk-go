// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"context"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
)

// Store defines methods for parsing, storing and searching for rule sets generated from
// OSCAL components.
type Store interface {
	// GetByRuleID return the rule object by the associated id or name.
	GetByRuleID(ctx context.Context, ruleID string) (extensions.RuleSet, error)
	// GetByCheckID returns rule object by the associated check id.
	GetByCheckID(ctx context.Context, checkID string) (extensions.RuleSet, error)

	// FindByComponent find rule objects by the associated component title.
	// Note that if a component is of type validation, this should only return
	// checks relevant to that validation component. Target component types (non-validation) should
	// return all information.
	FindByComponent(ctx context.Context, componentID string) ([]extensions.RuleSet, error)

	// All lists all indexed rules.
	All(ctx context.Context) ([]extensions.RuleSet, error)
}
