vendor:
	go mod tidy
	go mod verify
	go mod vendor
.PHONY: vendor

test-unit:
	go test -race -v -coverprofile=coverage.out ./...
.PHONY: test-unit

sanity: vendor format vet
	git diff --exit-code
.PHONY: sanity

format:
	go fmt ./...
.PHONY: format

vet:
	go vet ./...
.PHONY: vet

lint:
	@golangci-lint run ./...
.PHONY: lint

# For testing only
update-schema:
	@curl -LJO https://raw.githubusercontent.com/oscal-compass/compliance-trestle/refs/heads/develop/release-schemas/oscal_complete_schema.json
.PHONY: update-schema

# For testing only
update-types:
	@npx quicktype -s schema oscal_complete_schema.json --package types --omit-empty --just-types-and-package -o types/oscal_core.go --top-level OSCALModels
.PHONY: update-types
