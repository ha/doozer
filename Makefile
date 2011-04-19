include Make.inc

TARG=github.com/ha/doozer
GOFILES=\
	client.go\
	msg.pb.go\

include $(GOROOT)/src/Make.pkg
include $(GOROOT)/src/pkg/goprotobuf.googlecode.com/hg/Make.protobuf
