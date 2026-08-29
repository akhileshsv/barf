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
		Flh:0.003,
		Flv:0.002,
		Fls:0.001,
		Frh:0.003,
		Frv:0.002,
		Frs:0.0,
		Tr:50.0,
		Tl:170.0,
		Theta:60.0,
		Verbose:true,
	}
	_ = c.CalcCl()
}


func TestCalcCableT(t *testing.T){
	c := Cbl{
		Title:"harrison-2.1",
		Ucl:500.0,
		Ar:0.007,
		Em:1.4e8,
		Sw:0.5,
		Xr:400.,
		Yr:0.,
		Tl:300.0,
		Theta:30.0,
		Verbose:true,
	}
	_ = c.CalcT()

	c = Cbl{
		Title:"harrison-2.2",
		Ucl:500.0,
		Ar:0.007,
		Em:1.4e8,
		Sw:0.5,
		Xr:400.,
		Yr:0.,
		Tl:300.0,
		Theta:30.0,
		Verbose:true,
		Lds:[][]float64{{100,0,100}},
	}
	_ = c.CalcT()


}
