GO ?= $(shell command -v go 2> /dev/null)
MANIFEST_FILE ?= plugin.json
PLUGIN_ID := com.mattermost.meeting-bot
PLUGIN_VERSION := $(shell grep '"version"' $(MANIFEST_FILE) | cut -d'"' -f4)
BUNDLE_NAME ?= $(PLUGIN_ID)-$(PLUGIN_VERSION).tar.gz
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

export GO111MODULE=on

.PHONY: all
all: dist

.PHONY: dist
dist:	## Builds the plugin
	@echo Building plugin
	@mkdir -p dist
	@echo Building server for $(GOOS)-$(GOARCH)
	@GOWORK=off $(GO) build -o dist/plugin-$(GOOS)-$(GOARCH) .
ifeq ($(GOOS),windows)
	@mv dist/plugin-$(GOOS)-$(GOARCH) dist/plugin-$(GOOS)-$(GOARCH).exe
endif
	@echo Building plugin bundle
	@mkdir -p bundle/$(PLUGIN_ID)
	@cp $(MANIFEST_FILE) bundle/$(PLUGIN_ID)/
	@mkdir -p bundle/$(PLUGIN_ID)/server
	@cp -r dist bundle/$(PLUGIN_ID)/server/
	@if [ -d assets ]; then cp -r assets bundle/$(PLUGIN_ID)/; fi
	@cd bundle && tar czf ../$(BUNDLE_NAME) $(PLUGIN_ID)
	@rm -rf bundle
	@echo Plugin built at: $(BUNDLE_NAME)

.PHONY: deploy
deploy: dist	## Deploys the plugin to your local Mattermost server
	@echo Deploying plugin to local server
	@cp $(BUNDLE_NAME) ../../server/plugins/
	@echo Plugin deployed. Restart your Mattermost server to load the plugin.

.PHONY: clean
clean:	## Cleans build artifacts
	rm -rf dist
	rm -f $(BUNDLE_NAME)

.PHONY: help
help:	## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

