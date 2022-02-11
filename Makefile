SHELL=bash
OK_MSG = \x1b[32m ✔\x1b[0m
GOLIST?=$$(go list ./... | grep -vE "(vendor|static)")
CGO_ENABLED=0
export GOFLAGS := -mod=vendor $(GOFLAGS)
default: build
build:
	go build -o ${GOPATH}/bin/assume-shell cmd/assume-shell/*.go

# tools fetches necessary dev requirements
tools-local:
	(GOFLAGS="" go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.0)
tools: tools-local

coverprofile:
	@go test $(TEST) -coverprofile coverage.out && go tool cover -html=coverage.out

lint:
	@echo -n "==> Checking that code complies with golint requirements..."
	@golangci-lint run
	@echo -e "$(OK_MSG)"

test: unittest

unittest:
	@echo "==> Checking that code complies with unit tests..."
	@go test -timeout 2s $(GOLIST) -cover
