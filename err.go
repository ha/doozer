package doozer

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
)

var (
	ErrNoAddrs = errors.New("no known address")
	ErrBadTag  = errors.New("bad tag")
	ErrClosed  = errors.New("closed")
)

var (
	ErrOther    response_Err = response_OTHER
	ErrNotDir   response_Err = response_NOTDIR
	ErrIsDir    response_Err = response_ISDIR
	ErrNoEnt    response_Err = response_NOENT
	ErrRange    response_Err = response_RANGE
	ErrOldRev   response_Err = response_REV_MISMATCH
	ErrTooLate  response_Err = response_TOO_LATE
	ErrReadonly response_Err = response_READONLY
)

type Error struct {
	Err    error
	Detail string
}

func newError(t *txn) *Error {
	return &Error{
		Err:    *t.resp.ErrCode,
		Detail: proto.GetString(t.resp.ErrDetail),
	}
}

func (e *Error) Error() (s string) {
	s = e.Err.Error()
	if e.Detail != "" {
		s += ": " + e.Detail
	}
	return s
}
