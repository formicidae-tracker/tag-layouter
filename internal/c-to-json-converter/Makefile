all: apriltag/libapriltag.a oldtags/liboldtags.a

apriltag/libapriltag.a:
	$(MAKE) -C apriltag libapriltag.a

oldtags/liboldtags.a:
	$(MAKE) -C oldtags

c-to-json-converter:
	go build
