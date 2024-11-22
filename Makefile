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
