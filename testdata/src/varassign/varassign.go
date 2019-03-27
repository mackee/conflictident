package vardecl

type hoge struct{}

func main() {
	hoge, fuga := 1, 1 // want "conflict identifier name of 'hoge' by testdata/src/varassign/varassign.go:3:6."
	_ = hoge
	_ = fuga
}
