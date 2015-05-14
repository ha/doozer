// Code generated by protoc-gen-go.
// source: msg.proto
// DO NOT EDIT!

package doozer

import proto "github.org/golang/protobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type request_Verb int32

const (
	request_GET    request_Verb = 1
	request_SET    request_Verb = 2
	request_DEL    request_Verb = 3
	request_REV    request_Verb = 5
	request_WAIT   request_Verb = 6
	request_NOP    request_Verb = 7
	request_WALK   request_Verb = 9
	request_GETDIR request_Verb = 14
	request_STAT   request_Verb = 16
	request_SELF   request_Verb = 20
	request_ACCESS request_Verb = 99
)

var request_Verb_name = map[int32]string{
	1:  "GET",
	2:  "SET",
	3:  "DEL",
	5:  "REV",
	6:  "WAIT",
	7:  "NOP",
	9:  "WALK",
	14: "GETDIR",
	16: "STAT",
	20: "SELF",
	99: "ACCESS",
}
var request_Verb_value = map[string]int32{
	"GET":    1,
	"SET":    2,
	"DEL":    3,
	"REV":    5,
	"WAIT":   6,
	"NOP":    7,
	"WALK":   9,
	"GETDIR": 14,
	"STAT":   16,
	"SELF":   20,
	"ACCESS": 99,
}

func (x request_Verb) Enum() *request_Verb {
	p := new(request_Verb)
	*p = x
	return p
}
func (x request_Verb) String() string {
	return proto.EnumName(request_Verb_name, int32(x))
}
func (x request_Verb) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
func (x *request_Verb) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(request_Verb_value, data, "request_Verb")
	if err != nil {
		return err
	}
	*x = request_Verb(value)
	return nil
}

type response_Err int32

const (
	response_OTHER        response_Err = 127
	response_TAG_IN_USE   response_Err = 1
	response_UNKNOWN_VERB response_Err = 2
	response_READONLY     response_Err = 3
	response_TOO_LATE     response_Err = 4
	response_REV_MISMATCH response_Err = 5
	response_BAD_PATH     response_Err = 6
	response_MISSING_ARG  response_Err = 7
	response_RANGE        response_Err = 8
	response_NOTDIR       response_Err = 20
	response_ISDIR        response_Err = 21
	response_NOENT        response_Err = 22
)

var response_Err_name = map[int32]string{
	127: "OTHER",
	1:   "TAG_IN_USE",
	2:   "UNKNOWN_VERB",
	3:   "READONLY",
	4:   "TOO_LATE",
	5:   "REV_MISMATCH",
	6:   "BAD_PATH",
	7:   "MISSING_ARG",
	8:   "RANGE",
	20:  "NOTDIR",
	21:  "ISDIR",
	22:  "NOENT",
}
var response_Err_value = map[string]int32{
	"OTHER":        127,
	"TAG_IN_USE":   1,
	"UNKNOWN_VERB": 2,
	"READONLY":     3,
	"TOO_LATE":     4,
	"REV_MISMATCH": 5,
	"BAD_PATH":     6,
	"MISSING_ARG":  7,
	"RANGE":        8,
	"NOTDIR":       20,
	"ISDIR":        21,
	"NOENT":        22,
}

func (x response_Err) Enum() *response_Err {
	p := new(response_Err)
	*p = x
	return p
}
func (x response_Err) Error() string {
	return x.String()
}
func (x response_Err) String() string {
	return proto.EnumName(response_Err_name, int32(x))
}
func (x response_Err) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
func (x *response_Err) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(response_Err_value, data, "response_Err")
	if err != nil {
		return err
	}
	*x = response_Err(value)
	return nil
}

type request struct {
	Tag              *int32        `protobuf:"varint,1,opt,name=tag" json:"tag,omitempty"`
	Verb             *request_Verb `protobuf:"varint,2,opt,name=verb,enum=doozer.request_Verb" json:"verb,omitempty"`
	Path             *string       `protobuf:"bytes,4,opt,name=path" json:"path,omitempty"`
	Value            []byte        `protobuf:"bytes,5,opt,name=value" json:"value,omitempty"`
	OtherTag         *int32        `protobuf:"varint,6,opt,name=other_tag" json:"other_tag,omitempty"`
	Offset           *int32        `protobuf:"varint,7,opt,name=offset" json:"offset,omitempty"`
	Rev              *int64        `protobuf:"varint,9,opt,name=rev" json:"rev,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (this *request) Reset()         { *this = request{} }
func (this *request) String() string { return proto.CompactTextString(this) }
func (*request) ProtoMessage()       {}

func (this *request) GetTag() int32 {
	if this != nil && this.Tag != nil {
		return *this.Tag
	}
	return 0
}

func (this *request) GetVerb() request_Verb {
	if this != nil && this.Verb != nil {
		return *this.Verb
	}
	return 0
}

func (this *request) GetPath() string {
	if this != nil && this.Path != nil {
		return *this.Path
	}
	return ""
}

func (this *request) GetValue() []byte {
	if this != nil {
		return this.Value
	}
	return nil
}

func (this *request) GetOtherTag() int32 {
	if this != nil && this.OtherTag != nil {
		return *this.OtherTag
	}
	return 0
}

func (this *request) GetOffset() int32 {
	if this != nil && this.Offset != nil {
		return *this.Offset
	}
	return 0
}

func (this *request) GetRev() int64 {
	if this != nil && this.Rev != nil {
		return *this.Rev
	}
	return 0
}

type response struct {
	Tag              *int32        `protobuf:"varint,1,opt,name=tag" json:"tag,omitempty"`
	Flags            *int32        `protobuf:"varint,2,opt,name=flags" json:"flags,omitempty"`
	Rev              *int64        `protobuf:"varint,3,opt,name=rev" json:"rev,omitempty"`
	Path             *string       `protobuf:"bytes,5,opt,name=path" json:"path,omitempty"`
	Value            []byte        `protobuf:"bytes,6,opt,name=value" json:"value,omitempty"`
	Len              *int32        `protobuf:"varint,8,opt,name=len" json:"len,omitempty"`
	ErrCode          *response_Err `protobuf:"varint,100,opt,name=err_code,enum=doozer.response_Err" json:"err_code,omitempty"`
	ErrDetail        *string       `protobuf:"bytes,101,opt,name=err_detail" json:"err_detail,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (this *response) Reset()         { *this = response{} }
func (this *response) String() string { return proto.CompactTextString(this) }
func (*response) ProtoMessage()       {}

func (this *response) GetTag() int32 {
	if this != nil && this.Tag != nil {
		return *this.Tag
	}
	return 0
}

func (this *response) GetFlags() int32 {
	if this != nil && this.Flags != nil {
		return *this.Flags
	}
	return 0
}

func (this *response) GetRev() int64 {
	if this != nil && this.Rev != nil {
		return *this.Rev
	}
	return 0
}

func (this *response) GetPath() string {
	if this != nil && this.Path != nil {
		return *this.Path
	}
	return ""
}

func (this *response) GetValue() []byte {
	if this != nil {
		return this.Value
	}
	return nil
}

func (this *response) GetLen() int32 {
	if this != nil && this.Len != nil {
		return *this.Len
	}
	return 0
}

func (this *response) GetErrCode() response_Err {
	if this != nil && this.ErrCode != nil {
		return *this.ErrCode
	}
	return 0
}

func (this *response) GetErrDetail() string {
	if this != nil && this.ErrDetail != nil {
		return *this.ErrDetail
	}
	return ""
}

func init() {
	proto.RegisterEnum("doozer.request_Verb", request_Verb_name, request_Verb_value)
	proto.RegisterEnum("doozer.response_Err", response_Err_name, response_Err_value)
}
