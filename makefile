VERSION?=$(shell git describe --tags --dirty --always)
export VERSION

IMPORT_PATH=$(shell cat go.mod | head -n 1 | awk '{print $$2}')
BIN_NAME=$(notdir $(IMPORT_PATH))

export GO111MODULE=on
export GIT_TERMINAL_PROMPT=1

default: build

# sub-makefiles
# for build tools, docker build, deploy, static web files.
include scripts/makefile.d/*

build: build/$(BIN_NAME)

clean:
	rm -rf build/

run:
run: build/$(BIN_NAME)

.sources:
	@echo \
	makefile \
	$(GO_SOURCES) \
	| tr " " "\n"


.PHONY: default build clean run auto-run .sources test
