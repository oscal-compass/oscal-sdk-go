package validation

import (
	"fmt"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/oscal-compass/oscal-sdk-go/models/modelutils"
)

// UuidValidator implements the Validator interface to check for duplicate UUIDs and ParamIds in OSCAL models.
// It ensures uniqueness of identifiers across the model structure.
type UuidValidator struct{}

func (d UuidValidator) Validate(model oscalTypes.OscalModels) error {
	if modelutils.HasDuplicateValuesByName(&model, "UUID") {
		return fmt.Errorf("duplicate UUIDs found")
	}
	if model.Profile != nil {
		if modelutils.HasDuplicateValuesByName(&model, "ParamId") {
			return fmt.Errorf("duplicate ParamIds found")
		}
	}
	return nil
}
