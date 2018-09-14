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

PROJECT_REPO := 'github.com/tinakurian'
PROJECT_NAME := 'build-tool-detector'

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
	@goagen controller	-d $(PROJECT_REPO)/$(PROJECT_NAME)/design -o controllers
	@goagen app     	-d $(PROJECT_REPO)/$(PROJECT_NAME)/design
	@goagen swagger 	-d $(PROJECT_REPO)/$(PROJECT_NAME)/design
	@goagen schema  	-d $(PROJECT_REPO)/$(PROJECT_NAME)/design -o public
	@goagen client  	-d $(PROJECT_REPO)/$(PROJECT_NAME)/design
	

build:
	@go build -o app/build-tool-detector

deploy: 
	./app/build-tool-detector