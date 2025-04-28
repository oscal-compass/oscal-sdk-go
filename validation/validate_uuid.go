package validation

import oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

type UuidValidator struct{}

func (n UuidValidator) Validate(models oscalTypes.OscalModels) error {

	models.ComponentDefinition.UUID = "c14d8812-7098-4a9b-8f89-cba41b6ff0d8"
	return nil
}
