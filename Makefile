export GO111MODULE=on

EXAMPLES := $(shell ./get_main_pkgs.sh ./example)

# All source code and documents. Used in spell check.
ALL_DOCS := $(shell find . -name '*.md' -type f | sort)
# All directories with go.mod files. Used in go mod tidy.
ALL_GO_MOD_DIRS := $(shell find . -type f -name 'go.mod' -exec dirname {} \; | sort)

GOTEST_MIN = go test -v -timeout 30s
GOTEST = $(GOTEST_MIN) -race
GOTEST_WITH_COVERAGE = $(GOTEST) -coverprofile=coverage.txt -covermode=atomic

.DEFAULT_GOAL := precommit

.PHONY: precommit

TOOLS_DIR := ./.tools

$(TOOLS_DIR)/golangci-lint: go.mod go.sum tools.go
	go build -o $(TOOLS_DIR)/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

$(TOOLS_DIR)/misspell: go.mod go.sum tools.go
	go build -o $(TOOLS_DIR)/misspell github.com/client9/misspell/cmd/misspell

$(TOOLS_DIR)/stringer: go.mod go.sum tools.go
	go build -o $(TOOLS_DIR)/stringer golang.org/x/tools/cmd/stringer

precommit: lint generate build examples test

.PHONY: test-with-coverage
test-with-coverage:
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go test ./... + coverage in $${dir}"; \
	  (cd "$${dir}" && \
	    $(GOTEST_WITH_COVERAGE) ./... && \
	    go tool cover -html=coverage.txt -o coverage.html); \
	done

.PHONY: ci
ci: precommit check-clean-work-tree test-with-coverage test-386

.PHONY: check-clean-work-tree
check-clean-work-tree:
	@if ! git diff --quiet; then \
	  echo; \
	  echo 'Working tree is not clean, did you forget to run "make precommit"?'; \
	  echo; \
	  git status; \
	  exit 1; \
	fi

.PHONY: build
build:
	# TODO: Fix this on windows.
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "compiling all packages in $${dir}"; \
	  (cd "$${dir}" && \
	    go build ./... && \
	    go test -run xxxxxMatchNothingxxxxx ./... >/dev/null); \
	done

.PHONY: test
test:
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go test ./... + race in $${dir}"; \
	  (cd "$${dir}" && \
	    $(GOTEST) ./...); \
	done

.PHONY: test-386
test-386:
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go test ./... GOARCH 386 in $${dir}"; \
	  (cd "$${dir}" && \
	    GOARCH=386 $(GOTEST_MIN) ./...); \
	done

.PHONY: examples
examples:
	@set -e; for ex in $(EXAMPLES); do \
	  echo "Building $${ex}"; \
	  (cd "$${ex}" && \
	   go build .); \
	done

.PHONY: lint
lint: $(TOOLS_DIR)/golangci-lint $(TOOLS_DIR)/misspell $(TOOLS_DIR)/stringer
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "golangci-lint in $${dir}"; \
	  (cd "$${dir}" && \
	    $(abspath $(TOOLS_DIR))/golangci-lint run --fix); \
	done
	$(TOOLS_DIR)/misspell -w $(ALL_DOCS)
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go mod tidy in $${dir}"; \
	  (cd "$${dir}" && \
	    go mod tidy); \
	done

generate:
	PATH="$(abspath $(TOOLS_DIR)):$${PATH}" go generate ./...
