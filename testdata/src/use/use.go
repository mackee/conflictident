package use

type hoge struct{}

func main() {
	h := &hoge{}
	hh := hoge{}

	_, _ = h, hh
}
