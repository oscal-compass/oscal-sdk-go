/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
	"github.com/oscal-compass/oscal-sdk-go/rules"
)

func TestApplyToComponents(t *testing.T) {
	tests := []struct {
		name               string
		settings           Settings
		componentID        string
		expError           string
		wantRules          []extensions.RuleSet
		postValidationFunc func(store rules.Store) bool
	}{
		{
			name:        "Valid/WithMappedRules",
			componentID: "testComponent1",
			settings: Settings{
				mappedRules: set.Set[string]{
					"testRule1": struct{}{},
					"testRule2": struct{}{},
				},
				selectedParameters: map[string]string{},
			},
			wantRules: []extensions.RuleSet{testSet2},
		},
		{
			name:        "Valid/WithParameterOverrides",
			componentID: "testComponent2",
			settings: Settings{
				selectedParameters: map[string]string{
					"testParam1": "updatedValue",
				},
				mappedRules: set.Set[string]{
					"testRule1": struct{}{},
					"testRule2": struct{}{},
				},
			},
			wantRules: []extensions.RuleSet{
				{
					Rule: extensions.Rule{
						ID:          "testRule1",
						Description: "Test Rule",
						Parameter: &extensions.Parameter{
							ID:          "testParam1",
							Description: "Test Parameter",
							Value:       "updatedValue",
						},
					},
					Checks: []extensions.Check{
						{
							ID:          "testCheck1",
							Description: "Test Check",
						},
					},
				},
				testSet2,
			},
			postValidationFunc: func(store rules.Store) bool {
				ruleSet, _ := store.GetByRuleID(context.TODO(), "testRule1")
				return ruleSet.Rule.Parameter != nil && ruleSet.Rule.Parameter.Value == ""
			},
		},
		{
			name:        "Invalid/InvalidSettings",
			componentID: "testComponent1",
			settings: Settings{
				mappedRules: set.Set[string]{
					"doesnotexists": struct{}{},
				},
			},
			expError: "no rules found with criteria for component testComponent1",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			testCtx := context.Background()
			store := newFakeStore()

			gotRules, err := ApplyToComponent(testCtx, c.componentID, store, c.settings)
			sort.SliceStable(gotRules, func(i, j int) bool {
				return gotRules[i].Rule.ID < gotRules[j].Rule.ID
			})

			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.wantRules, gotRules)
			}

			if c.postValidationFunc != nil {
				require.True(t, c.postValidationFunc(store))
			}
		})
	}
}

var (
	testSet1 = extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "testRule1",
			Description: "Test Rule",
			Parameter: &extensions.Parameter{
				ID:          "testParam1",
				Description: "Test Parameter",
			},
		},
		Checks: []extensions.Check{
			{
				ID:          "testCheck1",
				Description: "Test Check",
			},
		},
	}
	testSet2 = extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "testRule2",
			Description: "Test Rule",
		},
		Checks: []extensions.Check{
			{
				ID:          "testCheck2",
				Description: "Test Check",
			},
		},
	}
	testSet3 = extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "testRule3",
			Description: "Test Rule",
			Parameter: &extensions.Parameter{
				ID:          "testParam3",
				Description: "Test Parameter",
				Value:       "default",
			},
		},
		Checks: []extensions.Check{
			{
				ID:          "testCheck3",
				Description: "Test Check",
			},
		},
	}
)

// fakeStore is a fake implementation of a rules.Store with static data
type fakeStore struct {
	staticRuleData map[string]extensions.RuleSet
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		staticRuleData: map[string]extensions.RuleSet{
			"testRule1": testSet1,
			"testRule2": testSet2,
			"testRule3": testSet3,
		},
	}
}

func (f fakeStore) GetByRuleID(ctx context.Context, ruleID string) (extensions.RuleSet, error) {
	ruleSet, ok := f.staticRuleData[ruleID]
	if !ok {
		return extensions.RuleSet{}, fmt.Errorf("rule %s not found", ruleID)
	}
	return ruleSet, nil
}

func (f fakeStore) GetByCheckID(ctx context.Context, checkID string) (extensions.RuleSet, error) {
	switch checkID {
	case "testCheck1":
		return f.staticRuleData["testRule1"], nil
	case "testCheck2":
		return f.staticRuleData["testRule2"], nil
	case "testCheck3":
		return f.staticRuleData["testRule3"], nil
	default:
		return extensions.RuleSet{}, fmt.Errorf("rule not found for %s", checkID)
	}
}

func (f fakeStore) FindByComponent(ctx context.Context, componentId string) ([]extensions.RuleSet, error) {
	switch componentId {
	case "testComponent1":
		return []extensions.RuleSet{testSet2, testSet3}, nil
	case "testComponent2":
		return []extensions.RuleSet{testSet1, testSet2}, nil
	default:
		return []extensions.RuleSet{}, fmt.Errorf("invalid component id: %s", componentId)
	}
}
