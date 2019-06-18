.PHONY: precommit

TOOLS_DIR := ./.tools

$(TOOLS_DIR)/golangci-lint: go.mod go.sum tools.go
	go build -o $(TOOLS_DIR)/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

$(TOOLS_DIR)/goimports: go.mod go.sum tools.go
	go build -o $(TOOLS_DIR)/goimports golang.org/x/tools/cmd/goimports

precommit: $(TOOLS_DIR)/goimports $(TOOLS_DIR)/golangci-lint 
	$(TOOLS_DIR)/goimports -d -local github.com/open-telemetry/opentelemetry-go -w .
	$(TOOLS_DIR)/golangci-lint run