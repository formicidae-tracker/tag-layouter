tag-layouter: apriltag/libapriltag.a
	go test -coverprofile cover.out
	go build

apriltag/libapriltag.a:
	$(MAKE) -C apriltag
