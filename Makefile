#! /usr/bin/make

CURRENT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

PROJECT_REPO := 'github.com/tinakurian'
PROJECT_NAME := 'build-tool-detector'

.DEFAULT_GOAL := all

all: clean install generate build test check

# Build configuration
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BINARY_DIR:=${PWD}/bin
BINARY:=build-tool-detector
GITUNTRACKEDCHANGES:=$(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
  COMMIT := $(COMMIT)-dirty
endif
LDFLAGS="-X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}"

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

build: clean generate $(BINARY) ## Compiles executable
$(BINARY): $(SOURCES)
	@go build -ldflags ${LDFLAGS} -o ${BINARY_DIR}/${BINARY}

help: ## Hey! That's me!
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

.PHONY: install ## Fetches all dependencies using dep
install:
	dep ensure -v

.PHONY: update ## Updates all dependencies defined for dep
update:
	dep ensure -update -v

clean: ## Cleans up the project binaries and generated Goa files
	@rm -rf app
	@rm -rf client
	@rm -rf tool
	@rm -rf public/swagger
	@rm -rf public/schema
	@rm -rf public/js
	if [ -f ${BINARY_DIR} ] ; then rm ${BINARY_DIR} ; fi


generate: clean ## (re)generates all goagen-generated files
	@goagen controller	-d $(PROJECT_REPO)/$(PROJECT_NAME)/design -o controllers
	@goagen app     	-d $(PROJECT_REPO)/$(PROJECT_NAME)/design
	@goagen swagger 	-d $(PROJECT_REPO)/$(PROJECT_NAME)/design
	@goagen schema  	-d $(PROJECT_REPO)/$(PROJECT_NAME)/design -o public
	@goagen client  	-d $(PROJECT_REPO)/$(PROJECT_NAME)/design

.PHONY: test
test: build ## Executes all tests
	@ginkgo -r

.PHONY: format ## Removes unneeded imports and formats source code
format:
	@goimports -l -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: check
check: ## Concurrently runs a whole bunch of static analysis tools
	@gometalinter --enable=misspell --enable=gosimple --enable-gc --vendor --skip=app --skip=client --skip=tool --deadline 300s ./...

.PHONY: run
run: ## runs the service locally
	${BINARY_DIR}/${BINARY}

.PHONY: tools
tools: ## Installs all necessary tools
	@echo "Installing gometalinter"
	@go get -u github.com/alecthomas/gometalinter && gometalinter --install

	@echo "Installing ginkgo"
	@go get -u github.com/onsi/ginkgo/ginkgo

	@echo "Installing goimports"
	@go get -u golang.org/x/tools/cmd/goimports

	@echo "Installing goa"
	@go get -u github.com/goadesign/goa/...