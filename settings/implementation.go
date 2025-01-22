/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"fmt"
	"path/filepath"
	"strings"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
)

var _ Settings = (*ImplementationSettings)(nil)

// ImplementationSettings defines settings for RuleSets defined at the control
// implementation level.
type ImplementationSettings struct {
	// requirementSettings defines settings for RuleSets at the
	// implemented requirement/individual control level.
	requirementSettings map[string]RequirementSettings
	// settings defines the settings for the
	// overall implementation of the requirements.
	settings RequirementSettings
}

// NewImplementationSettings returned ImplementationSettings for an OSCAL ControlImplementationSet.
func NewImplementationSettings(controlImplementation oscalTypes.ControlImplementationSet) *ImplementationSettings {
	implementation := &ImplementationSettings{
		requirementSettings: make(map[string]RequirementSettings),
		settings: RequirementSettings{
			mappedRules:        map[string]struct{}{},
			selectedParameters: make(map[string]string),
		},
	}
	if controlImplementation.SetParameters != nil {
		setParameters(*controlImplementation.SetParameters, implementation.settings.selectedParameters)
	}

	for _, implementedReq := range controlImplementation.ImplementedRequirements {
		requirement := NewSettingsFromImplementedRequirement(implementedReq)
		for mappedRule := range requirement.mappedRules {
			implementation.settings.mappedRules.Add(mappedRule)
		}
		implementation.requirementSettings[implementedReq.ControlId] = requirement
	}

	return implementation
}

func (i *ImplementationSettings) ApplyParameterSettings(set extensions.RuleSet) extensions.RuleSet {
	return i.settings.ApplyParameterSettings(set)
}

func (i *ImplementationSettings) ContainsRule(ruleId string) bool {
	return i.settings.ContainsRule(ruleId)
}

// ByControlID returns the individual requirement settings for a given control id in the
// control implementation.
func (i *ImplementationSettings) ByControlID(controlId string) (RequirementSettings, error) {
	requirement, ok := i.requirementSettings[controlId]
	if !ok {
		return RequirementSettings{}, fmt.Errorf("control %s not found in settings", controlId)
	}
	return requirement, nil
}

// GetFrameworkShortName returns the human-readable short name for the control source in a
// control implementation set and whether this value is populated.
//
// This function checks the associated properties and falls back to the implementation
// Source reference.
func GetFrameworkShortName(implementation oscalTypes.ControlImplementationSet) (string, bool) {
	const (
		expectedPathParts = 3
		modelIDIndex      = 1
		filenameIndex     = 2
	)
	// Looks for the property, fallback to parsing it out of the control source href.
	if implementation.Props != nil {
		property, found := extensions.GetTrestleProp(extensions.FrameworkProp, *implementation.Props)
		if found {
			return property.Value, true
		}
	}

	// Fallback to the control source string based on trestle
	// workspace conventions of $MODEL/$MODEL_ID/$MODEL.json.
	cleanedSource := filepath.Clean(implementation.Source)
	parts := strings.Split(cleanedSource, "/")
	if len(parts) == expectedPathParts && strings.HasSuffix(parts[filenameIndex], ".json") {
		return parts[modelIDIndex], true
	}

	return "", false
}

// Framework returns ImplementationSettings from a list of OSCAL Control Implementations for a given framework.
func Framework(framework string, controlImplementations []oscalTypes.ControlImplementationSet) (*ImplementationSettings, error) {
	var implementationSettings *ImplementationSettings

	for _, controlImplementation := range controlImplementations {
		frameworkShortName, found := GetFrameworkShortName(controlImplementation)
		if found && frameworkShortName == framework {
			implementationSettings = NewImplementationSettings(controlImplementation)
			break
		}
	}

	if implementationSettings == nil {
		return implementationSettings, fmt.Errorf("framework %s is not in control implementations", framework)
	}
	return implementationSettings, nil
}
