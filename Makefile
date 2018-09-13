#! /usr/bin/make
#
# Makefile for build-tool-detector
#
# Targets:
# - depend	  installs the project's dependencies
# - clean     delete all generated files
# - generate  (re)generate all goagen-generated files
# - build     compile executable
# - deploy 	  deploy to localhost:8080
#
# Meta targets:
# - all is the default target, it runs all the targets in the order above.
#
CURRENT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: depend clean generate build deploy

depend:
	# dep init
	dep ensure -v

clean:
	@rm -rf app
	@rm -rf client
	@rm -rf tool
	@rm -rf public/swagger
	@rm -rf public/schema
	@rm -rf public/js
	@rm -f build-tool-detector

generate:
	@goagen controller	-d build-tool-detector/design -o controllers
	@goagen app     	-d build-tool-detector/design
	@goagen swagger 	-d build-tool-detector/design
	@goagen schema  	-d build-tool-detector/design -o public
	@goagen client  	-d build-tool-detector/design
	

build:
	@go build -o app/build-tool-detector

deploy: 
	./app/build-tool-detector