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
	file, err := os.Open("../testdata/component-definition-test.json")
	if err != nil {
		log.Fatalf("failed to open component definition, %v", err)
	}
	definition, err := models.NewComponentDefinition(file, validation.NoopValidator{})
	if err != nil {
		log.Fatalf("failed to read component definition, %v", err)
	}

	if definition != nil {
		assessmentPlan, err := transformers.ComponentDefinitionsToAssessmentPlan(context.Background(), []oscalTypes.ComponentDefinition{*definition}, "cis")
		if err != nil {
			log.Fatalf("failed to create assessment plan, %v", err)
		}
		reviewedControlsJson, err := json.Marshal(assessmentPlan.ReviewedControls.ControlSelections)
		if err != nil {
			log.Fatalf("failed to marshal reviewed controls, %v", err)
		}
		fmt.Println(string(reviewedControlsJson))
	}
	// Output:
	// [{"include-controls":[{"control-id":"CIS-2.1"}]}]
}
