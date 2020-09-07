VERSION?=$(shell git describe --tags --dirty --always)
export VERSION

IMPORT_PATH=$(shell cat go.mod | head -n 1 | awk '{print $$2}')
BIN_NAME=$(notdir $(IMPORT_PATH))

export GO111MODULE=on
export GIT_TERMINAL_PROMPT=1

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

run: export INVENTORY_FILE=build/test.yaml
run: build/$(BIN_NAME)
	# log level: trace
	@rm $(INVENTORY_FILE) || true
	$< -vvvv
	$< template echo 'echo {{args}}'
	$< template ssh 'ssh -o "StrictHostKeyChecking=no" -n {{.user}}@{{.name}} -C {{args}}'
	$< template scp 'scp "{{arg 1}}" {{.name}}:"{{arg 2}}"'
	$< template rsh 'rsh -l {{.user}} {{.name}} {{args}}'
	$< template gce 'gcloud compute --project "{{.project}}" ssh --zone "{{.zone}}" "{{.name}}"'
	$< set test1.node
	$< set -l cluster=test test1.node test2.node
	$< set -l role=worker -l rack=1 test1.node
	$< set -l added=190201 1.node
	$< set -l added=190202 2.node
	$< set -l added=190203 3.node
	$< set -l name=a test2.node
	for rack in $$(seq -f rack%02g 1 8); do \
		seq -f "node.$$rack.%03g" 1 32 | xargs $< set -l rack=$$rack -l cluster=test2; \
	done
	$< get -e 'name=="test1.node"'
	$< get -o yaml
	$< get -o json | jq
	$< get -e 'rack>="rack05" && name contains "001"'
	$< run -e 'cluster=="test" && role=="worker"' -- echo hello world
	TIME="%E" time $< run -e 'rack>="rack05" && name contains "001"' -w 2 -- sleep 1
	$< run -e 'name contains "001"' -t echo -- hello world
	$< run -e 'name contains "001"' -t echo -o wide -- hello world
	$< run -e 'name contains "001"' -t echo -o text -- hello world
	$< run -e 'name contains "001"' -t echo -o json -- hello world
	$< get -e 'rack=="rack08"' | xargs $< del
	$< get | xargs $< set -l rack="-"
	$< shortcut set rack-new 'name contains "rack07" || name contains "rack06"'
	$< get -s rack-new
	$< run -s rack-new -t echo -- 1
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
	| tr " " "\n"

test:
	go test -v ./pkg/...

.PHONY: default build clean run auto-run .sources test
