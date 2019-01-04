# Makefile
BASE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
MKFILE_PATH := $(BASE_PATH)/Makefile

# Project specific top level directories
PROJ_TMP_PATH := ./tmp
CMD_PATH := ./cmd
ART_PATH := ./artifacts
REL_PATH := ./releases
ART_DEV_PATH := ./artifacts_dev
PKG_PATH := ./pkg

# Target directories
BUILD_ROOT := $(ART_PATH)
BUILD_DARWIN := $(BUILD_ROOT)/darwin
BUILD_LINUX := $(BUILD_ROOT)/linux

# Target dev directories
BUILD_DEV_ROOT := $(ART_DEV_PATH)
BUILD_DEV_DARWIN := $(BUILD_DEV_ROOT)/darwin
BUILD_DEV_LINUX := $(BUILD_DEV_ROOT)/linux

# Artifacts definitions
SERVICE_NAME := ohas
ART_DARWIN_64 := $(SERVICE_NAME)-darwin-amd64.bin
ART_LINUX_32 := $(SERVICE_NAME)-linux-386.bin
ART_LINUX_64 := $(SERVICE_NAME)-linux-amd64.bin
ART_ARCHIVE := $(SERVICE_NAME)-$$(git describe --abbrev=0)-$$(date +%Y%m%d%H%M%S).zip

# Additional tools
# Standalone algorithms
DEV_ALG_NAME := ohal
ART_ALG_DARWIN_64 := $(DEV_ALG_NAME)-darwin-amd64.bin
ART_ALG_LINUX_32 := $(DEV_ALG_NAME)-linux-386.bin
ART_ALG_LINUX_64 := $(DEV_ALG_NAME)-linux-amd64.bin

# default json files
GUIDELINE_JSON := guideline_hearts.json
GUIDELINE_CONTENT_JSON := guideline_hearts_content.json
GOALS_JSON := goals_hearts.json
GOALS_CONTENT_JSON := goals_hearts_content.json
SAMPLE_REQUEST_JSON := sample-request.json
HELP_FILE := INSTRUCTIONS.md

# tests
COVER_OUT := cover.out

# VERSION=1.0.0
# COMMIT=`git rev-parse HEAD`
# BRANCH=`git rev-parse --abbrev-ref HEAD`
# BUILD=`git rev-parse --short HEAD`

BUILD_ID := `git rev-parse --short HEAD`
LDFLAGS_DEV := -ldflags "-X main.appCommit=$(BUILD_ID)"
LDFLAGS := -ldflags "-X main.appCommit=$(BUILD_ID) -s -w"
LAST_TAG := `git describe --abbrev=0`

all:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Update vendored dependencies
	dep ensure -update

clean_all: clean clean_dev ## Clean all the targets at once

clean: clean_darwin clean_linux clean_releases clean_tmp ## Clean target paths

clean_darwin: ## Clean target path for darwin
	-rm -fr $(BUILD_DARWIN)

clean_linux: ## Clean target path for linux
	-rm -fr $(BUILD_LINUX)

clean_releases: ## Clean target path for releases
	-rm -fr $(REL_PATH)/*

clean_tmp: ## Clean target path for tmp
	-rm -fr $(PROJ_TMP_PATH)/*

clean_dev: clean_dev_darwin clean_dev_linux ## Clean dev target path

clean_dev_darwin: ## Clean dev target path for darwin
	-rm -fr $(BUILD_DEV_DARWIN)

clean_dev_linux: ## Clean dev target path for linux
	-rm -fr $(BUILD_DEV_LINUX)

build: build_darwin build_linux ## Build binaries

build_darwin: ## Build binaries for darwin
	mkdir -p $(BUILD_DARWIN)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -a -tags netgo -installsuffix cgo $(LDFLAGS) -o $(BUILD_DARWIN)/$(ART_DARWIN_64) $(CMD_PATH)/$(SERVICE_NAME)

build_linux: ## Build binaries for linux
	mkdir -p $(BUILD_LINUX)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix cgo $(LDFLAGS) -o $(BUILD_LINUX)/$(ART_LINUX_64) $(CMD_PATH)/$(SERVICE_NAME)
#	CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -a -tags netgo -installsuffix cgo $(LDFLAGS) -o $(BUILD_LINUX)/$(ART_LINUX_32) $(CMD_PATH)/$(SERVICE_NAME)

build_dev: build_dev_darwin build_dev_linux ## Build dev binaries

build_dev_darwin: ## Build dev binaries for darwin
	mkdir -p $(BUILD_DEV_DARWIN)
	GOOS=darwin GOARCH=amd64 go build -race $(LDFLAGS_DEV) -o $(BUILD_DEV_DARWIN)/$(ART_ALG_DARWIN_64) ./cmd/$(DEV_ALG_NAME)

build_dev_linux: ## Build dev binaries for linux
	mkdir -p $(BUILD_DEV_LINUX)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS_DEV) -o $(BUILD_DEV_LINUX)/$(ART_ALG_LINUX_64) ./cmd/$(DEV_ALG_NAME)
#	GOOS=linux GOARCH=386 go build $(LDFLAGS_DEV) -o $(BUILD_DEV_LINUX)/$(ART_ALG_LINUX_32) ./cmd/$(DEV_ALG_NAME)

artifacts: artifacts_darwin artifacts_linux ## Create artifacts
	$(MAKE) -f $(MKFILE_PATH) house_keep

artifacts_darwin: ## Create artifacts for darwin
	$(MAKE) -f $(MKFILE_PATH) clean_darwin
	$(MAKE) -f $(MKFILE_PATH) build_darwin

artifacts_linux: ## Create artifacts for linux
	$(MAKE) -f $(MKFILE_PATH) clean_linux
	$(MAKE) -f $(MKFILE_PATH) build_linux

zip_artifacts: ## Create a zip archive with artifacts
	$(MAKE) -f $(MKFILE_PATH) clean_releases
	mkdir -p $(REL_PATH)
	zip -j -v $(REL_PATH)/$(ART_ARCHIVE) $(BUILD_DARWIN)/$(ART_DARWIN_64) $(BUILD_LINUX)/$(ART_LINUX_64) $(GUIDELINE_JSON) $(GUIDELINE_CONTENT_JSON) $(GOALS_JSON) $(GOALS_CONTENT_JSON) $(SAMPLE_REQUEST_JSON) $(HELP_FILE)

house_keep: ## Remove any .DS_Store files
	find $(BASE_PATH) -name ".DS_Store" -depth -exec rm {} \;

test: ## Run tests
	go test ./... -coverpkg=./... -coverprofile=$(COVER_OUT)

cover: ## Show tests coverage
	@if [ -f $(COVER_OUT) ]; then \
		go tool cover -func=$(COVER_OUT); \
		rm -f $(COVER_OUT); \
	else \
		echo "$(COVER_OUT) is missing. Please run 'make test'"; \
	fi

.PHONY: all deps clean_all \
	clean clean_darwin clean_linux clean_releases clean_tmp \
	build build_darwin build_linux \
	clean_dev clean_dev_darwin clean_dev_linux \
	build_dev build_dev_darwin build_dev_linux \
	artifacts artifacts_darwin artifacts_linux \
	zip_artifacts test cover \
	house_keep
