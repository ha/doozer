package doozer

const (
	_ = 1 << iota
	_
	set
	del
)

type Event struct {
	Rev  int64
	Path string
	Body []byte
	Flag int32
}

func (e Event) IsSet() bool {
	return e.Flag&set > 0
}

func (e Event) IsDel() bool {
	return e.Flag&del > 0
}
