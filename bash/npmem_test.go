package barf

import (
	"fmt"
	"testing"
	kass"barf/kass"
)

func TestNpMem(t *testing.T){
	mem := &kass.MemNp{
		Ts:[]int{0,0,0},
		Ls:[]float64{1,1,1},
		Ds:[]float64{1.0,1.0,1.0,1.0},
		Bs:[]float64{1.0,1.0,1.0},
		Styp:1,
		Dims:[]float64{1.0,1.0},
		Frmtyp:"1db",
		Em:25e6,
		Lspan:3.0,
	}
	fmt.Println(mem)

}
