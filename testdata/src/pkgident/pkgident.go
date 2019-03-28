package pkgident

import (
	"path"
)

func hoge() {
	_ = path.ErrBadPattern
}

func main() {
	path, fuga := 1, 1 // want "conflict identifier name of 'path' by testdata/src/pkgident/pkgident.go:4:2."
	_, _ = path, fuga
}

func p(path string) { // want "conflict identifier name of 'path' by testdata/src/pkgident/pkgident.go:4:2."
}
