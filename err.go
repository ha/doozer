package doozer

import (
	"fmt"
	"goprotobuf.googlecode.com/hg/proto"
	"os"
)

var (
	ErrNoAddrs = os.NewError("no known address")
	ErrBadTag  = os.NewError("bad tag")
	ErrClosed  = os.NewError("closed")
)


var (
	ErrNotDir   response_Err = response_NOTDIR
	ErrIsDir    response_Err = response_ISDIR
	ErrNoEnt    response_Err = response_NOENT
	ErrRange    response_Err = response_RANGE
	ErrOldRev   response_Err = response_REV_MISMATCH
	ErrTooLate  response_Err = response_TOO_LATE
	ErrReadonly response_Err = response_READONLY
)


func (e response_Err) String() string {
	s := response_Err_name[int32(e)]
	if s == "" {
		s = fmt.Sprint("unknown error ", e)
	}
	return s
}


type Error struct {
	Err    os.Error
	Detail string
}


func newError(t *txn) *Error {
	return &Error{
		Err:    *t.resp.ErrCode,
		Detail: proto.GetString(t.resp.ErrDetail),
	}
}


func (e *Error) String() (s string) {
	s = e.Err.String()
	if e.Detail != "" {
		s += ": " + e.Detail
	}
	return s
}
