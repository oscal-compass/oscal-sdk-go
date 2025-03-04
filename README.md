# oscal-sdk-go

[![Go Report Card](https://goreportcard.com/badge/github.com/oscal-compass/oscal-sdk-go)](https://goreportcard.com/report/github.com/oscal-compass/oscal-sdk-go)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/oscal-compass/oscal-sdk-go)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/oscal-compass/oscal-sdk-go/badge)](https://scorecard.dev/viewer/?uri=github.com/oscal-compass/oscal-sdk-go)

`oscal-sdk-go` complements the `compliance-trestle` SDK by providing the core SDK functionality in Go.

> WARNING: This project is currently under initial development. APIs may be changed incompatibly from one commit to another.

## Supported Functionality

Below is a table to show what is currently supported by the SDK.

| SDK Functionality                         | Supported          |
|-------------------------------------------|--------------------|
| OSCAL Types with Basic Trestle Extensions | :heavy_check_mark: |
| OSCAL Schema Validation                   | :heavy_check_mark: |
| Target Components Extension               | :x:                |
| Multiple Parameters per Rule              | :x:                |
| OSCAL to OSCAL Transformation             | :heavy_check_mark: |
| OSCAL Constraints Validation              | :x:                |


## Get Started

Add the module as dependency to your project:

```console
go get github.com/oscal-compass/oscal-sdk-go
```

### SDK Terms

[`Extensions`](https://github.com/oscal-compass/oscal-sdk-go/tree/main/extensions): `oscal-compass` uses OSCAL properties to [extend](https://pages.nist.gov/OSCAL/learn/tutorials/general/extension/#props) OSCAL.  
[`Rules`](https://github.com/oscal-compass/oscal-sdk-go/tree/main/rules): Rules are associated with Components and define a mechanism to verify the proper implementation of technical controls.  
[`Settings`](https://github.com/oscal-compass/oscal-sdk-go/tree/main/settings): Settings define adjustments to fine-tune pre-defined options in Rules for the implementation of a specific compliance framework.  

### Perform a Transformation

```go
import (
    ...
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/oscal-compass/oscal-sdk-go/validation"

	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/transformers"
)

func main() {
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
```

## Contributing

Our project welcomes external contributions. Please see `CONTRIBUTING.md` to get started.

## Code of Conduct

Participation in the OSCAL Compass community is governed by the [Code of Conduct](https://github.com/oscal-compass/community/blob/main/CODE_OF_CONDUCT.md).

## Acknowledgments

This project leverages [`go_oscal`](https://github.com/defenseunicorns/go-oscal) to provide Go types for the OSCAL schema.
