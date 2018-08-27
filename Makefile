#! /usr/bin/make
#
# Makefile for goa cellar example
#
# Targets:
# - clean     delete all generated files
# - generate  (re)generate all goagen-generated files.
# - build     compile executable
# - ae-build  build appengine
# - ae-dev    deploy to local (dev) appengine
# - ae-deploy deploy to appengine
#
# Meta targets:
# - all is the default target, it runs all the targets in the order above.
#
CURRENT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: depend clean generate build deploy

depend:
	# dep init
	dep ensure

clean:
	@rm -rf app
	@rm -rf client
	@rm -rf tool
	@rm -rf public/swagger
	@rm -rf public/schema
	@rm -rf public/js
	@rm -f build-tool-detector

generate:
	@goagen bootstrap 	-d build-tool-detector/design
	@goagen app     	-d build-tool-detector/design
	@goagen swagger 	-d build-tool-detector/design -o public
	@goagen schema  	-d build-tool-detector/design -o public
	@goagen client  	-d build-tool-detector/design
	

build:
	@go build -o build-tool-detector

deploy: 
	./build-tool-detector