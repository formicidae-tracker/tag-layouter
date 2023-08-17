all: apriltag/libapriltag.a oldtags/liboldtags.a

apriltag/libapriltag.a:
	$(MAKE) -C apriltag

oldtags/liboldtags.a:
	$(MAKE) -C oldtags

check:
	go test
