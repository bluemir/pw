VERSION?=$(shell git describe --long --tags --dirty --always)
export VERSION

IMPORT_PATH=$(shell cat go.mod | head -n 1 | awk '{print $$2}')
BIN_NAME=$(notdir $(IMPORT_PATH))

export GO111MODULE=on
export GIT_TERMINAL_PROMPT=1

DOCKER_IMAGE_NAME=bluemir/$(BIN_NAME)

## Go Sources
GO_SOURCES = $(shell find . -name "vendor"  -prune -o \
                            -type f -name "*.go" -print)

default: build

build: build/$(BIN_NAME)

build/$(BIN_NAME): $(GO_SOURCES) makefile
	@mkdir -p build
	go build -v \
		-ldflags "-X main.VERSION=$(VERSION)" \
		$(OPTIONAL_BUILD_ARGS) \
		-o $@ main.go

clean:
	rm -rf build/

run: TEST_FILE=build/test.yaml
run: build/$(BIN_NAME)
	# log level: trace
	@rm build/test.yaml || true
	$< -vvvv
	$< -vvvv -i $(TEST_FILE) set test1.node
	$< -vvvv -i $(TEST_FILE) set -l cluster=test test1.node test2.node
	$< -vvvv -i $(TEST_FILE) set -l role=worker -l rack=1 test1.node
	$< -vvvv -i $(TEST_FILE) set -l added=190201 1.node
	$< -vvvv -i $(TEST_FILE) set -l added=190202 2.node
	$< -vvvv -i $(TEST_FILE) set -l added=190203 3.node
	$< -vvvv -i $(TEST_FILE) set -l _name=a test2.node
	$< -vvvv -i $(TEST_FILE) get -e 'name=="test1.node"'
	$< -vvvv -i $(TEST_FILE) get -e 'added>="190202"'
	$< -vvvv -i $(TEST_FILE) get -o yaml
	$< -vvvv -i $(TEST_FILE) get -o json
	$< -vvvv -i $(TEST_FILE) run -e 'cluster=="test" && role=="worker"' -- echo hello world
	#$< -vvvv -i $(TEST_FILE) run -e 'cluster=="test" && role=="worker"' -t echo -- hello world
	cat build/test.yaml

auto-run:
	while true; do \
		$(MAKE) .sources | \
		entr -rd $(MAKE) test run ;  \
		echo "hit ^C again to quit" && sleep 1  \
	; done

.sources:
	@echo \
	makefile \
	$(GO_SOURCES) \
	$(JS_SOURCES) \
	$(HTML_SOURCES) \
	$(CSS_SOURCES) \
	$(WEB_LIBS) \
	$(HTML_TEMPLATE) \
	| tr " " "\n"

test:
	go test -v ./pkg/...

.PHONY: default build clean run auto-run .sources test
