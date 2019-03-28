package multifile

import "path/filepath"

var _ = filepath.ErrBadPattern

type file struct{}

var path = "path"

const (
	cnt = "cnt"
)
