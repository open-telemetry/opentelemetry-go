ALL_PKGS := $(shell go list ./...)

# All source code and documents. Used in spell check.
ALL_DOCS := $(shell find . -name '*.md' -type f | sort)

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

$(TOOLS_DIR)/misspell: go.mod go.sum tools.go
	go build -o $(TOOLS_DIR)/misspell github.com/client9/misspell/cmd/misspell

precommit: $(TOOLS_DIR)/goimports $(TOOLS_DIR)/golangci-lint  $(TOOLS_DIR)/misspell 
	$(TOOLS_DIR)/golangci-lint run --fix # TODO: Fix this on windows.
	$(TOOLS_DIR)/misspell -w $(ALL_DOCS)

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

all-docs:
	@echo $(ALL_DOCS) | tr ' ' '\n' | sort