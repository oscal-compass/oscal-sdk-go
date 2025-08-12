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

func ExampleSSPToAssessmentPlan() {
	file, err := os.Open("../testdata/test-ssp.json")
	if err != nil {
		log.Fatalf("failed to open system security plan, %v", err)
	}
	ssp, err := models.NewSystemSecurityPlan(file, validation.NoopValidator{})
	if err != nil {
		log.Fatalf("failed to read system security plan, %v", err)
	}

	if ssp != nil {
		assessmentPlan, err := transformers.SSPToAssessmentPlan(context.Background(), *ssp, "importPath")
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
	// [{"include-controls":[{"control-id":"ex-1"},{"control-id":"ex-2"}]}]

}

func ExampleAssessmentPlanToAssessmentResults() {
	file, err := os.Open("../testdata/test-ap.json")
	if err != nil {
		log.Fatalf("failed to open assessment plan, %v", err)
	}
	plan, err := models.NewAssessmentPlan(file, validation.NoopValidator{})

	if err != nil {
		log.Fatalf("failed to read assessment plan, %v", err)
	}

	if plan != nil {
		assessmentResults, err := transformers.AssessmentPlanToAssessmentResults(*plan, "importPath")
		if err != nil {
			log.Fatalf("failed to create assessment results, %v", err)
		}

		href := assessmentResults.ImportAp.Href
		fmt.Println(href)

		if len(assessmentResults.Results) == 0 {
			log.Fatalf("failed to find assessment results, %v", err)
		}
		reviewedControlsJson, err := json.Marshal(assessmentResults.Results[0].ReviewedControls.ControlSelections)
		if err != nil {
			log.Fatalf("failed to marshal reviewed controls, %v", err)
		}
		fmt.Println(string(reviewedControlsJson))
	}
	// Output:
	// importPath
	// [{"include-controls":[{"control-id":"ex-2"},{"control-id":"ex-1"}]},{"include-controls":[{"control-id":"ex-1"}]}]

}
