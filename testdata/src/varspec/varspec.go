package varspec

type hoge struct{}

func main() {
	var hoge int // want "conflict identifier name of 'hoge' by testdata/src/varspec/varspec.go:3:6."
	_ = hoge
}
