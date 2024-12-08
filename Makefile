MAIN_PACKAGE_PATH := .
BINARY_NAME := csv-import
APP_VERSION ?= $(shell git describe --tags --always --dirty)
APP_GIT_COMMIT ?= $(shell git rev-parse HEAD)
APP_GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
APP_GIT_REPOSITORY ?= https://github.com/green-ecolution/tbz-csv-import-plugin
APP_BUILD_TIME ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
define GOFLAGS
-ldflags=" \
	-s -w \
  -X main.version=$(APP_VERSION) \
"
endef

.PHONY: all
all: build

.PHONY: generate
generate:
	@echo "Generating..."
	go generate 

.PHONY: build/ui
build/ui:
	@echo "Building UI..."
	@cd ui && yarn install && yarn build

.PHONY: build
build: generate
	@echo "Building..."
	@$(MAKE) build/ui
	go build $(GOFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

.PHONY: run
run: generate
	@echo "Running..."
	@$(MAKE) build/ui
	go run $(GOFLAGS) $(MAIN_PACKAGE_PATH)

