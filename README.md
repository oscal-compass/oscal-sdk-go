# oscal-sdk-go

`oscal-sdk-go` complements the `compliance-trestle` SDK by providing the core SDK functionality in Go.

> WARNING: This project is currently under initial development. APIs may be changed incompatibly from one commit to another.

## Supported Functionality

Below is a table to show what is currently supported by the SDK.

| SDK Functionality                         | Supported |
|-------------------------------------------|-----------|
| OSCAL Types with Basic Trestle Extensions | &#10003;  |
| OSCAL Schema Validation                   |           |
| Target Components Extension               |           |
| Multiple Parameters per Rule              |           |
| OSCAL to OSCAL Transformation             |           |
| OSCAL Constraints Validation              |           |

## Run tests

```bash
make test-unit
```

## Format and Style

**Requires [`golangci-lint`](https://golangci-lint.run/welcome/quick-start/)**

```bash
make format
# For issue identification
make vet
# Linting
make lint
```

# Acknowledgments

This project leverages [`go_oscal`](https://github.com/defenseunicorns/go-oscal) to provide Go types for the OSCAL schema.
