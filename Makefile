# Makefile
BASE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
MKFILE_PATH := $(BASE_PATH)/Makefile

# Project specific top level directories
PROJ_TMP_PATH := ./tmp
ART_PATH := ./artifacts
REL_PATH := ./releases
PKG_PATH := ./pkg

# Target directories
BUILD_ROOT := $(ART_PATH)

# Artifacts definitions
SERVICE_NAME := ohas
ART_DARWIN_64 := $(SERVICE_NAME)-darwin-$$(git describe --abbrev=1)
ART_LINUX_64 := $(SERVICE_NAME)-linux-$$(git describe --abbrev=1)
ART_ARCHIVE := $(SERVICE_NAME)-$$(git describe --abbrev=1)-$$(git branch | grep \* | cut -d ' ' -f2)-$$(date +%y%m%d).zip

# default files to include
# GUIDELINE_JSON := guideline_hearts.json
# GUIDELINE_CONTENT_JSON := guideline_hearts_content.json
# GOALS_JSON := goals_hearts.json
# GOALS_CONTENT_JSON := goals_hearts_content.json
SAMPLE_REQUEST_JSON := sample-request.json
HELP_FILE := ./documentation
CONTENTS_FILE := ./contents
CONFIG_FILE := ohas.toml

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

clean_all: clean ## Clean all the targets at once

clean: clean_darwin_linux clean_releases clean_tmp ## Clean target paths

clean_darwin_linux: ## Clean target path
	-rm -fr $(BUILD_ROOT)

clean_releases: ## Clean target path for releases
	-rm -fr $(REL_PATH)/*

clean_tmp: ## Clean target path for tmp
	-rm -fr $(PROJ_TMP_PATH)/*

build: build_darwin build_linux ## Build binaries

build_darwin: ## Build binaries for darwin
	mkdir -p $(BUILD_ROOT)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -a -tags netgo -installsuffix cgo $(LDFLAGS) -o $(BUILD_ROOT)/$(ART_DARWIN_64) $(BASE_PATH)/main.go

build_linux: ## Build binaries for linux
	mkdir -p $(BUILD_ROOT)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix cgo $(LDFLAGS) -o $(BUILD_ROOT)/$(ART_LINUX_64) $(BASE_PATH)/main.go

artifacts: ## Create artifacts
	$(MAKE) -f $(MKFILE_PATH) house_keep
	$(MAKE) -f $(MKFILE_PATH) clean
	$(MAKE) -f $(MKFILE_PATH) build

zip_artifacts: ## Create a zip archive with artifacts
	$(MAKE) -f $(MKFILE_PATH) clean_releases
	mkdir -p $(REL_PATH)
	zip -j -v $(REL_PATH)/$(ART_ARCHIVE) $(BUILD_ROOT)/$(ART_DARWIN_64) $(BUILD_ROOT)/$(ART_LINUX_64) $(CONFIG_FILE) $(SAMPLE_REQUEST_JSON)
	zip -v -r -u $(REL_PATH)/$(ART_ARCHIVE) $(CONTENTS_FILE)	
	zip -v -r -u $(REL_PATH)/$(ART_ARCHIVE) $(HELP_FILE)

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

.PHONY: all clean_all \
	clean clean_darwin_linux clean_releases clean_tmp \
	build build_darwin build_linux \
	artifacts zip_artifacts test cover \
	house_keep
