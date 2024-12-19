package barf

import (
	//"fmt"
	"testing"
)

func TestCalcCableL(t *testing.T){
	var c Cbl
	c = Cbl{
		Ucl:500,
		Ar:0.007,
		Em:140e6,
		Sw:0.5,
		Xr:350,
		Yr:100,
		Alp:0.000011,
		Lds:[][]float64{
			{-60,-80,150},
			{60, -80,300},
		},
		Verbose:true,
	}
	_ = c.CalcL()
}
