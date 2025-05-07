test-run: export INVENTORY_FILE=runtime/test.yaml
test-run: build/$(APP_NAME)
	# log level: trace
	@rm $(INVENTORY_FILE) || true
	$< -vvvv
	$< template echo -- echo {{args}}
	$< template ssh -- ssh -o "StrictHostKeyChecking=no" -n {{.user}}@{{.name}} -C {{args}}
	$< template scp -- scp "{{arg 1}}" {{.name}}:"{{arg 2}}"
	$< template rsh -- rsh -l {{.user}} {{.name}} {{args}}
	$< template gce -- gcloud compute --project "{{.project}}" ssh --zone "{{.zone}}" "{{.name}}"
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
	# multiple-line output
	rm runtime/test-file || true
	echo "1" >> runtime/test-file
	echo "2" >> runtime/test-file
	echo "3" >> runtime/test-file
	$< run -s rack-new -- podman
	@echo "================ Done =============="

