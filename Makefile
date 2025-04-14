SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
.SILENT:

# go builds are fast enough that we can just build on demand instead of trying to do any fancy
# change detection
build: clean selfman
.PHONY: build

selfman:
	go build ./cmd/selfman

clean:
	rm -f ./selfman
.PHONY: clean

check:
	go test ./...
	go vet ./...
.PHONY: check
