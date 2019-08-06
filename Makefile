tag-layouter: apriltag/libapriltag.a oldtags/liboldtags.a
	go test -coverprofile cover.out
	go build

apriltag/libapriltag.a:
	$(MAKE) -C apriltag

oldtags/liboldtags.a:
	$(MAKE) -C oldtags
