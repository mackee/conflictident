package consttest

const (
	enumA = iota
	enumB
	enumC
	enumD
)

func cnt() int {
	consttest := enumC
	if consttest == enumA {
		enumD := 1 // want "conflict identifier name of 'enumD' by testdata/src/consttest/consttest.go:7:2"
		_ = enumD
		return enumA
	}
	return enumB
}

func swh() {
	e := enumC
	switch e {
	case enumA:
		println(enumA)
	case enumC:
		enumD := 1 // want "conflict identifier name of 'enumD' by testdata/src/consttest/consttest.go:7:2"
		_ = enumD
		println(enumC)
	default:
		println(enumB)
	}
}
