
.PHONY: build
build:
	go build .

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: install
install:
	go install .

.PHONY: test
test:
	go test ./internal/...

.PHONY: all
all: build lint test install