package tools

import (
	"fmt"
	"testing"
)

type Base struct {
	Name string
}

type P1 struct {
	Age int
	B   Base
}

type P2 struct {
	Age int
	B   Base
}

func TestSimpleCopyProperties(t *testing.T) {
	p1 := P1{
		2,
		Base{"age"},
	}
	p2 := P2{}

	err := SimpleCopyProperties(&p2, &p1)
	fmt.Println(err)
}
