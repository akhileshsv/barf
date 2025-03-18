package barf

import (
	//"fmt"
	"testing"
)

func TestCalcCableCl(t *testing.T){
	c := Cbl{
		Title:"harrison-1",
		Ucl:500.,
		Ar:0.007,
		Em:1.4e8,
		Sw:0.5,
		Xr:350.,
		Yr:100.,
		Alp:0.000011,
		Lds:[][]float64{
			{80,-60,150},
			{80, 60,300},
		},
		Tl:250.0,
		Theta:45.0,
		Verbose:true,
	}
	_ = c.CalcCl()
}
