ALL_SRC := $(shell find . -name '*.go' -type f | sort)
ALL_PKGS := $(shell go list $(sort $(dir $(ALL_SRC))))

GOTEST=go test
GOTEST_OPT?=-v -race -timeout 30s
GOTEST_OPT_WITH_COVERAGE = $(GOTEST_OPT) -coverprofile=coverage.txt -covermode=atomic

.DEFAULT_GOAL := precommit

.PHONY: precommit

TOOLS_DIR := ./.tools

$(TOOLS_DIR)/golangci-lint: go.mod go.sum tools.go
	go build -o $(TOOLS_DIR)/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

$(TOOLS_DIR)/goimports: go.mod go.sum tools.go
	go build -o $(TOOLS_DIR)/goimports golang.org/x/tools/cmd/goimports

precommit: $(TOOLS_DIR)/goimports $(TOOLS_DIR)/golangci-lint 
	$(TOOLS_DIR)/goimports -d -local github.com/open-telemetry/opentelemetry-go -w .
	$(TOOLS_DIR)/golangci-lint run # TODO: Fix this on windows.

.PHONY: test-with-coverage
test-with-coverage:
	$(GOTEST) $(GOTEST_OPT_WITH_COVERAGE) $(ALL_PKGS)
	go tool cover -html=coverage.txt -o coverage.html

.PHONY: circle-ci
circle-ci: precommit test-with-coverage test-386

.PHONY: test
test:
	$(GOTEST) $(GOTEST_OPT) $(ALL_PKGS)

.PHONY: test-386
test-386:
	GOARCH=386 $(GOTEST) -v -timeout 30s $(ALL_PKGS)

all-pkgs:
	@echo $(ALL_PKGS) | tr ' ' '\n' | sort

all-srcs:
	@echo $(ALL_SRC) | tr ' ' '\n' | sort
