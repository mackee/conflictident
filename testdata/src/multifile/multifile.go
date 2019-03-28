package multifile

func main() {
	file := "file"         // want "conflict identifier name of 'file' by testdata/src/multifile/decl.go:7:6."
	path := 1              // want "conflict identifier name of 'path' by testdata/src/multifile/decl.go:9:5."
	cnt := true            // want "conflict identifier name of 'cnt' by testdata/src/multifile/decl.go:12:2."
	filepath := "filepath" // this is no warning, so import decl is file scope
	_, _, _, _ = file, path, cnt, filepath
}
