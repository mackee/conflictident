package receiver

import "path"

var _ = path.ErrBadPattern

type hoge struct{}

func (h hoge) String() string {
	return "hoge"
}

func (h hoge) Hogehoge(path string) { // want "conflict identifier name of 'path' by testdata/src/receiver/receiver.go:3:8."
}

func (h hoge) path() {}

type ppath struct{}

func (path ppath) String() string { // want "conflict identifier name of 'path' by testdata/src/receiver/receiver.go:3:8."
	return "ppath"
}
