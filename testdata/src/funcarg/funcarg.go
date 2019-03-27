package funcarg

type hoge struct{}

func f(hoge int, fuga bool) { // want "conflict identifier name of 'hoge' by testdata/src/funcarg/funcarg.go:3:6."
}
