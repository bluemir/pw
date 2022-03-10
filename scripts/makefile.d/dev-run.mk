
dev-run:
	while true; do \
		$(MAKE) .sources | \
		entr -rd $(MAKE) test;  \
		echo "hit ^C again to quit" && sleep 1  \
	; done

.PHONY: dev-run reset
