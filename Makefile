include Make.inc

TARG=github.com/ha/doozer
GOFILES=\
	client.go\
	msg.pb.go\

include $(GOROOT)/src/Make.pkg

msg.pb.go: msg.proto
	mkdir -p _pb
	protoc --go_out=_pb $<
	cat _pb/$@\
	|sed s/\\bRequest/request/g\
	|sed s/\\bResponse/response/g\
	|sed s/\\bNewRequest/newRequest/g\
	|sed s/\\bNewResponse/newResponse/g\
	|gofmt >$@
	rm -rf _pb

CLEANFILES+=_pb
