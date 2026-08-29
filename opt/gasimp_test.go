package barf

import (
	"fmt"
	"testing"
)

func TestSimpGa(t *testing.T){
	n := 42
	fmt.Println("string rep-",int2bin(n))
	fmt.Println("string len-",len([]rune(int2bin(n))))
}
