##@ Run

run: build/$(APP_NAME) ## Run web app
	$< -vvv
dev-run: ## Run dev server. If detect file change, automatically rebuild&restart server
	@$(MAKE) build/tools/watcher
	watcher \
		--include "go.mod" \
		--include "go.sum" \
		--include "**.go" \
		--include "package.json" \
		--include "yarn.lock" \
		--include "assets/**" \
		--include "api/proto/**" \
		--include "Makefile" \
		--include "scripts/makefile.d/*.mk" \
		--exclude "build/**" \
		--exclude "**.sw*" \
		--exclude "assets/js/index.js" \
		-- \
	$(MAKE) test run


reset: ## Kill all make process. Use when dev-run stuck.
	ps -e | grep make | grep -v grep | awk '{print $$1}' | xargs kill


#tools: build/tools/entr
build/tools/entr:
	@which $(notdir $@) || (echo "see http://eradman.com/entrproject")

tools: build/tools/watcher
build/tools/watcher: build/tools/go
	@which $(notdir $@) || (./scripts/tools/install/go-tool.sh github.com/bluemir/watcher)


.PHONY: dev-run reset
