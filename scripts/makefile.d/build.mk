##@ Build
## Go Sources
GO_SOURCES = $(shell find . -name "vendor"  -prune -o \
                            -type f -name "*.go" -print)

.watched_sources: $(GO_SOURCES) go.mod go.sum
build/docker-image: $(GO_SOURCES)

.PHONY: build
build: build/$(APP_NAME) ## Build web app

.PHONY: test
test: fmt vet ## Run test
	go test -v ./...

export CGO_ENABLED=0

build/$(APP_NAME): $(GO_SOURCES) $(MAKEFILE_LIST) fmt vet
	@$(MAKE) build/tools/go
	@mkdir -p build
	go build -v \
		-trimpath \
		-ldflags "\
			-X '$(IMPORT_PATH)/internal/buildinfo.AppName=$(APP_NAME)' \
			-X '$(IMPORT_PATH)/internal/buildinfo.Version=$(VERSION)' \
			-X '$(IMPORT_PATH)/internal/buildinfo.BuildTime=$(shell date --rfc-3339=ns)' \
		" \
		$(OPTIONAL_BUILD_ARGS) \
		-o $@ main.go

build-tools: build/tools/go #build/tools/rice
build/tools/go:
	@which $(notdir $@) || echo "see https://golang.org/doc/install"
build/tools/rice: build/tools/go
	@which $(notdir $@) || (go get -u github.com/GeertJohan/go.rice/rice)

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...
