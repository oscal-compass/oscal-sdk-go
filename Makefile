vendor:
	go mod tidy
	go mod verify
	go mod vendor
.PHONY: vendor

test-unit:
	go test -race -v -coverprofile=coverage.out ./...
.PHONY: test-unit

format:
	go fmt ./...
.PHONY: format

lint:
	@golangci-lint run ./...
.PHONY: lint
