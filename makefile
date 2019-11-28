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
	@rm $(TEST_FILE) || true
	$< -vvvv
	$< -i $(TEST_FILE) template echo 'echo {{args}}'
	$< -i $(TEST_FILE) template ssh 'ssh -o "StrictHostKeyChecking=no" -n {{.user}}@{{.name}} -C {{args}}'
	$< -i $(TEST_FILE) template scp 'scp "{{arg 1}}" {{.name}}:"{{arg 2}}"'
	$< -i $(TEST_FILE) template rsh 'rsh -l {{.user}} {{.name}} {{args}}'
	$< -i $(TEST_FILE) template gce 'gcloud compute --project "{{.project}}" ssh --zone "{{.zone}}" "{{.name}}"'
	$< -i $(TEST_FILE) set test1.node
	$< -i $(TEST_FILE) set -l cluster=test test1.node test2.node
	$< -i $(TEST_FILE) set -l role=worker -l rack=1 test1.node
	$< -i $(TEST_FILE) set -l added=190201 1.node
	$< -i $(TEST_FILE) set -l added=190202 2.node
	$< -i $(TEST_FILE) set -l added=190203 3.node
	$< -i $(TEST_FILE) set -l name=a test2.node
	for rack in $$(seq -f rack%02g 1 8); do \
		seq -f "node.$$rack.%03g" 1 32 | xargs $< -i $(TEST_FILE) set -l rack=$$rack -l cluster=test2; \
	done
	$< -i $(TEST_FILE) get -e 'name=="test1.node"'
	$< -i $(TEST_FILE) get -o yaml
	$< -i $(TEST_FILE) get -o json | jq
	$< -i $(TEST_FILE) get -e 'rack>="rack05" && name contains "001"'
	$< -i $(TEST_FILE) run -e 'cluster=="test" && role=="worker"' -- echo hello world
	TIME="%E" time $< -i $(TEST_FILE) run -e 'rack>="rack05" && name contains "001"' -w 2 -- sleep 1
	$< -i $(TEST_FILE) run -e 'name contains "001"' -t echo -- hello world
	$< -i $(TEST_FILE) run -e 'name contains "001"' -t echo -o wide -- hello world
	$< -i $(TEST_FILE) run -e 'name contains "001"' -t echo -o text -- hello world
	$< -i $(TEST_FILE) run -e 'name contains "001"' -t echo -o json -- hello world
	@echo "================ Done =============="

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
