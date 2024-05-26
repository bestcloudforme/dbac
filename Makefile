LOCAL_BIN := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/.bin

DEFAULT_GO_RUN_ARGS ?= ""

GOLANGCI_LINT_VERSION := latest
REVIVE_VERSION := v1.3.7

.PHONY: all
all: clean tools lint build

.PHONY: clean
clean:
	rm -rf $(LOCAL_BIN)

.PHONY: tools
tools: golangci-lint-install revive-install vendor

.PHONY: golangci-lint-install
golangci-lint-install:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: revive-install
revive-install:
	GOBIN=$(LOCAL_BIN) go install github.com/mgechev/revive@$(REVIVE_VERSION)

.PHONY: lint
lint: tools run-lint

.PHONY: run-lint
run-lint: lint-golangci-lint lint-revive

.PHONY: lint-golangci-lint
lint-golangci-lint:
	$(info running golangci-lint...)
	$(LOCAL_BIN)/golangci-lint -v run ./... || (echo golangci-lint returned an error, exiting!; sh -c 'exit 1';)

.PHONY: lint-revive
lint-revive:
	$(info running revive...)
	$(LOCAL_BIN)/revive -formatter=stylish -config=build/ci/.revive.toml -exclude ./vendor/... ./... || (echo revive returned an error, exiting!; sh -c 'exit 1';)

.PHONY: upgrade-deps
upgrade-deps: vendor
	for item in `grep -v 'indirect' go.mod | grep '/' | cut -d ' ' -f 1`; do \
		echo "trying to upgrade direct dependency $$item" ; \
		go get -u $$item ; \
	done
	go mod tidy
	go mod vendor

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor: tidy
	go mod vendor

.PHONY: build
build: vendor
	$(info building binary...)
	go build -o .bin/dbac main.go || (echo an error while building binary, exiting!; sh -c 'exit 1';)

.PHONY: run
run: vendor
	go run main.go $(DEFAULT_GO_RUN_ARGS)