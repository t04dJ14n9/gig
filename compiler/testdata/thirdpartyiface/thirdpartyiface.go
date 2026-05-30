package thirdpartyiface

type Stringer interface {
	String() string
}

func AcceptStringer(Stringer) int {
	return 1
}
