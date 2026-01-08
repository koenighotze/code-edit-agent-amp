.DEFAULT_GOAL := build

.PHONY: all build vet clean

install.tools:
	go install golang.org/x/vuln/cmd/govulncheck@latest

clean:
	go clean -x -i

fmt:
	go fmt ./cmd/... ./internal/... ./pkg/...

vet: fmt
	go vet ./cmd/... ./internal/... ./pkg/...

lint:
	golangci-lint run ./...

deps.upgrade:
	go get -u ./...
	go mod tidy

deps.vulncheck:
	govulncheck ./...

get.dependencies:
	go mod tidy

build: get.dependencies
	go build -o code-edit-agent-amp .
