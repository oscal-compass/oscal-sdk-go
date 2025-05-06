package transformers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/oscal-compass/oscal-sdk-go/validation"

	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/transformers"
)

func ExampleComponentDefinitionsToAssessmentPlan() {
	file, err := os.Open("example-component-definition.json")
	if err != nil {
		log.Fatalf("failed to open component definition, %v", err)
	}
	definition, err := models.NewComponentDefinition(file, validation.NoopValidator{})
	if err != nil {
		log.Fatalf("failed to read component definition, %v", err)
	}

	if definition != nil {
		assessmentPlan, err := transformers.ComponentDefinitionsToAssessmentPlan(context.Background(), []oscalTypes.ComponentDefinition{*definition}, "example-framework")
		if err != nil {
			log.Fatalf("failed to create assessment plan, %v", err)
		}
		assessmentPlanJSON, err := json.MarshalIndent(assessmentPlan, "", " ")
		if err != nil {
			log.Fatalf("failed to marshal assessment plan, %v", err)
		}
		fmt.Println(string(assessmentPlanJSON))
	}
}
