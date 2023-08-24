all: tag check tag-layouter arena-layouter

tag:
	make -C internal/tag

tag-layouter:
	go build ./cmd/tag-layouter

arena-layouter:
	go build ./cmd/arena-layouter


check:
	make -C internal/tag check
	go test ./cmd/tag-layouter
	go test ./cmd/arena-layouter


.PHONY: all check tag tag-layouter arena-layouter
