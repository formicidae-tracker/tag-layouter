all: tag check tag-layouter arena-layouter

tag:
	make -C internal/tag

tag-layouter:
	go build ./cmd/tag-layouter

arena-layouter:
	go build ./cmd/arena-layouter


check:
	make -C internal/tag check
	make -C cmd/tag-layouter check
	make -C cmd/arena-layouter check

.PHONY: all check tag tag-layouter arena-layouter
