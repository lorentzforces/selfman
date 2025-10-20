SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
.SILENT:

GO_BUILD_FLAGS := -buildvcs=true

help:
	echo '  make clean      remove generated files'
	echo '  make build      build the project from scratch'
	echo '  make selfman    build executable if not already built'
	echo '  make check      execute tests and checks'
.PHONY: help

# go builds are fast enough that we can just build on demand instead of trying to do any fancy
# change detection
build: clean selfman
.PHONY: build

selfman:
	go build ${GO_BUILD_FLAGS} ./cmd/selfman

clean:
	rm -f ./selfman
.PHONY: clean

check:
	go test ${GO_BUILD_FLAGS} ./...
	go vet ./...
.PHONY: check
