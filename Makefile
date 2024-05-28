LOCAL_BIN := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/.bin

DEFAULT_GO_RUN_ARGS ?= ""

GOLANGCI_LINT_VERSION := latest
REVIVE_VERSION := v1.3.7

.PHONY: all
all: clean tools lint build
	@echo " # Completed all steps."

.PHONY: clean
clean:
	@echo " # Cleaning up the workspace..."
	@rm -fr $(LOCAL_BIN)
	@if [ -d "vendor" ]; then \
		echo " # Removing vendor directory..."; \
		rm -fr vendor; \
	fi
	@echo " # Clean complete."

.PHONY: tools
tools: golangci-lint-install revive-install vendor
	@echo " # All tools installed."

.PHONY: golangci-lint-install
golangci-lint-install:
	@echo " # Installing golangci-lint..."
	@GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@echo " # golangci-lint installation complete."

.PHONY: revive-install
revive-install:
	@echo " # Installing revive..."
	@GOBIN=$(LOCAL_BIN) go install github.com/mgechev/revive@$(REVIVE_VERSION)
	@echo " # Revive installation complete."

.PHONY: lint
lint: tools run-lint
	@echo " # Linting complete."

.PHONY: run-lint
run-lint: lint-golangci-lint lint-revive
	@echo " # Running all linters..."

.PHONY: lint-golangci-lint
lint-golangci-lint:
	@echo " # Running golangci-lint..."
	@$(LOCAL_BIN)/golangci-lint -v run ./... || (echo "golangci-lint returned an error, exiting!"; sh -c 'exit 1';)

.PHONY: lint-revive
lint-revive:
	@echo " # Running revive..."
	@$(LOCAL_BIN)/revive -formatter=stylish -config=build/ci/.revive.toml -exclude ./vendor/... ./... || (echo "Revive returned an error, exiting!"; sh -c 'exit 1';)

.PHONY: upgrade-deps
upgrade-deps: vendor
	@echo " # Upgrading dependencies..."
	@for item in `grep -v 'indirect' go.mod | grep '/' | cut -d ' ' -f 1`; do \
		echo "Trying to upgrade direct dependency $$item" ; \
		go get -u $$item ; \
	done
	@go mod tidy
	@go mod vendor
	@echo " # Dependencies upgraded."

.PHONY: tidy
tidy:
	@echo " # Tidying up the go.mod and go.sum files..."
	@go mod tidy
	@echo " # Go module tidy complete."

.PHONY: vendor
vendor: tidy
	@echo " # Vendoring dependencies..."
	@go mod vendor
	@echo " # Vendoring complete."

.PHONY: build
build: vendor
	@echo " # Building binary..."
	@go build -o .bin/dbac main.go || (echo "An error occurred while building the binary, exiting!"; sh -c 'exit 1';)
	@echo " # Binary built successfully."

.PHONY: run
run: vendor
	@echo " # Running the application..."
	@go run main.go $(DEFAULT_GO_RUN_ARGS)
	@echo " # Application run complete."
