# Copyright The OpenTelemetry Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

EXAMPLES := $(shell ./get_main_pkgs.sh ./example)
TOOLS_MOD_DIR := ./tools

# All source code and documents. Used in spell check.
ALL_DOCS := $(shell find . -name '*.md' -type f | sort)
# All directories with go.mod files related to opentelemetry library. Used for building, testing and linting.
ALL_GO_MOD_DIRS := $(filter-out $(TOOLS_MOD_DIR), $(shell find . -type f -name 'go.mod' -exec dirname {} \; | sort))
ALL_COVERAGE_MOD_DIRS := $(shell find . -type f -name 'go.mod' -exec dirname {} \; | egrep -v '^./example|^$(TOOLS_MOD_DIR)' | sort)

# Mac OS Catalina 10.5.x doesn't support 386. Hence skip 386 test
SKIP_386_TEST = false
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	SW_VERS := $(shell sw_vers -productVersion)
	ifeq ($(shell echo $(SW_VERS) | egrep '^(10.1[5-9]|1[1-9]|[2-9])'), $(SW_VERS))
		SKIP_386_TEST = true
	endif
endif

GOCOVER_PKGS = $(shell go list ./... | grep -v 'internal/opentelemetry-proto' | paste -s -d, - )
GOTEST_MIN = go test -timeout 30s
GOTEST = $(GOTEST_MIN) -race
GOTEST_WITH_COVERAGE = $(GOTEST) -coverprofile=coverage.txt -covermode=atomic -coverpkg=$(GOCOVER_PKGS)

.DEFAULT_GOAL := precommit

.PHONY: precommit

TOOLS_DIR := $(abspath ./.tools)

$(TOOLS_DIR)/golangci-lint: $(TOOLS_MOD_DIR)/go.mod $(TOOLS_MOD_DIR)/go.sum $(TOOLS_MOD_DIR)/tools.go
	cd $(TOOLS_MOD_DIR) && \
	go build -o $(TOOLS_DIR)/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

$(TOOLS_DIR)/misspell: $(TOOLS_MOD_DIR)/go.mod $(TOOLS_MOD_DIR)/go.sum $(TOOLS_MOD_DIR)/tools.go
	cd $(TOOLS_MOD_DIR) && \
	go build -o $(TOOLS_DIR)/misspell github.com/client9/misspell/cmd/misspell

$(TOOLS_DIR)/stringer: $(TOOLS_MOD_DIR)/go.mod $(TOOLS_MOD_DIR)/go.sum $(TOOLS_MOD_DIR)/tools.go
	cd $(TOOLS_MOD_DIR) && \
	go build -o $(TOOLS_DIR)/stringer golang.org/x/tools/cmd/stringer

$(TOOLS_DIR)/gojq: $(TOOLS_MOD_DIR)/go.mod $(TOOLS_MOD_DIR)/go.sum $(TOOLS_MOD_DIR)/tools.go
	cd $(TOOLS_MOD_DIR) && \
	go build -o $(TOOLS_DIR)/gojq github.com/itchyny/gojq/cmd/gojq

precommit: generate build lint examples test

.PHONY: test-with-coverage
test-with-coverage:
	set -e; for dir in $(ALL_COVERAGE_MOD_DIRS); do \
	  echo "go test ./... + coverage in $${dir}"; \
	  (cd "$${dir}" && \
	    $(GOTEST_WITH_COVERAGE) ./... && \
	    go tool cover -html=coverage.txt -o coverage.html); \
	done

.PHONY: ci
ci: precommit check-clean-work-tree license-check test-with-coverage test-386

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
	if [ $(SKIP_386_TEST) = true ] ; then \
	  echo "skipping the test for GOARCH 386 as it is not supported on the current OS"; \
	else \
	  set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	    echo "go test ./... GOARCH 386 in $${dir}"; \
	    (cd "$${dir}" && \
	      GOARCH=386 $(GOTEST_MIN) ./...); \
	  done; \
	fi

.PHONY: examples
examples:
	@set -e; for ex in $(EXAMPLES); do \
	  echo "Building $${ex}"; \
	  (cd "$${ex}" && \
	   go build .); \
	done

.PHONY: lint
lint: $(TOOLS_DIR)/golangci-lint $(TOOLS_DIR)/misspell
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "golangci-lint in $${dir}"; \
	  (cd "$${dir}" && \
	    $(TOOLS_DIR)/golangci-lint run --fix && \
	    $(TOOLS_DIR)/golangci-lint run); \
	done
	$(TOOLS_DIR)/misspell -w $(ALL_DOCS)
	set -e; for dir in $(ALL_GO_MOD_DIRS) $(TOOLS_MOD_DIR); do \
	  echo "go mod tidy in $${dir}"; \
	  (cd "$${dir}" && \
	    go mod tidy); \
	done

.PHONY: generate
generate: stringer protobuf

.PHONY: stringer
stringer: $(TOOLS_DIR)/stringer
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "running generators in $${dir}"; \
	  (cd "$${dir}" && \
	    PATH="$(TOOLS_DIR):$${PATH}" go generate ./...); \
	done

.PHONY: license-check
license-check:
	@licRes=$$(for f in $$(find . -type f \( -iname '*.go' -o -iname '*.sh' \) ! -path './vendor/*' ! -path './internal/opentelemetry-proto/*') ; do \
	           awk '/Copyright The OpenTelemetry Authors|generated|GENERATED/ && NR<=3 { found=1; next } END { if (!found) print FILENAME }' $$f; \
	   done); \
	   if [ -n "$${licRes}" ]; then \
	           echo "license header checking failed:"; echo "$${licRes}"; \
	           exit 1; \
	   fi

# Find all .proto files.
OTEL_PROTO_SUBMODULE := internal/opentelemetry-proto
PROTO_GEN_DIR        := internal/opentelemetry-proto-gen
PROTOBUF_TEMP_DIR    := gen/pb-go
PROTOBUF_SOURCE_DIR  := gen/proto
ORIG_PROTO_FILES     := $(wildcard $(OTEL_PROTO_SUBMODULE)/opentelemetry/proto/*/v1/*.proto \
                           $(OTEL_PROTO_SUBMODULE)/opentelemetry/proto/collector/*/v1/*.proto)
SOURCE_PROTO_FILES   := $(subst $(OTEL_PROTO_SUBMODULE),$(PROTOBUF_SOURCE_DIR),$(ORIG_PROTO_FILES))

ifeq ($(CIRCLECI),true)
copy-protos-to:
	# create a dummy container to hold the protobufs
	docker create -v /defs --name proto-orig alpine:3.4 /bin/true
	# copy from here into the proto-src directory
	docker cp ./gen proto-orig:/defs

copy-protobuf-from:
	rm -fr ./$(PROTO_GEN_DIR)
	# copy from dummy volume back to our source directory
	docker cp proto-orig:/defs/$(PROTOBUF_TEMP_DIR)/go.opentelemetry.io/otel/$(PROTO_GEN_DIR) ./$(PROTO_GEN_DIR)

else
copy-protos-to:
	@:    # nop

copy-protobuf-from:
	rm -fr ./$(PROTO_GEN_DIR)
	mv ./$(PROTOBUF_TEMP_DIR)/go.opentelemetry.io/otel/$(PROTO_GEN_DIR) ./$(PROTO_GEN_DIR)
endif


# This step can be omitted assuming go_package changes are made in opentelemetry-proto repo
define exec-replace-pkgname
sed  's+go_package = "github.com/open-telemetry/opentelemetry-proto/gen/go+go_package = "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen+' < $(1) > $(2)

endef

.PHONY: protobuf make-protobufs
protobuf: protobuf-source copy-protos-to make-protobufs copy-protobuf-from

.PHONY: protobuf-source
protobuf-source: $(SOURCE_PROTO_FILES) | $(PROTOBUF_SOURCE_DIR)/

# replace opentelemetry-proto v0.4.0 package name by repo-local version
$(SOURCE_PROTO_FILES): $(PROTOBUF_SOURCE_DIR)/%.proto: $(OTEL_PROTO_SUBMODULE)/%.proto
	@mkdir -p $(@D)
	$(call exec-replace-pkgname,$<,$@)

ifeq ($(CIRCLECI),true)
VOLUMES_MOUNT=--volumes-from proto-orig
else
VOLUMES_MOUNT=-v `pwd`:/defs
endif

define exec-protoc-all
docker run $(VOLUMES_MOUNT) namely/protoc-all $(1)

endef

make-protobufs: $(SOURCE_PROTO_FILES)  | $(PROTOBUF_GEN_DIR)/
	$(foreach file,$(subst ${PROTOBUF_SOURCE_DIR}/,,$(SOURCE_PROTO_FILES)),$(call exec-protoc-all, -i $(PROTOBUF_SOURCE_DIR) -f ${file} -l go -o ${PROTOBUF_TEMP_DIR}))
	rm -fr ./gen/go


$(PROTOBUF_SOURCE_DIR)/ $(PROTO_GEN_DIR)/:
	mkdir -p $@

.PHONY: protobuf-clean
protobuf-clean:
	rm -rf ./gen ./$(PROTO_GEN_DIR)

