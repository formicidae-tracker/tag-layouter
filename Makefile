all: tag check

tag:
	make -C internal/tag

check:
	make -C internal/tag check

.PHONY: all check
