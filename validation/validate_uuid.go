package validation

import (
	"fmt"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/oscal-compass/oscal-sdk-go/models/modelutils"
)

type UuidValidator struct{}

func (d UuidValidator) Validate(model oscalTypes.OscalModels) error {
	if !modelutils.HasDuplicateValuesByName(&model, "UUID") {
		return fmt.Errorf("duplicate UUIDs found")
	}
	if model.Profile != nil {
		if !modelutils.HasDuplicateValuesByName(&model, "ParamId") {
			return fmt.Errorf("duplicate ParamIds found")
		}
	}
	return nil
}
