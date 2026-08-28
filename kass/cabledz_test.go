package barf

import (
	//"fmt"
	"testing"
)

func TestCblAz(t *testing.T){
	c := Cbl{
		Title:"harrison-15.3.1",
		Ucl:486.,
		Ar:0.007,
		Em:1.4e8,
		Sw:0.5,
		Xr:350.,
		Yr:100.,
		Alp:0.000011,
		Flh:0.003,
		Flv:0.002,
		Fls:0.001,
		Frh:0.003,
		Frv:0.002,
		Frs:0.,
		Tl:170.0,
		Theta:60.0,
		Tr:50,
		Verbose:true,
	}
	_ = c.CblAz()


	c = Cbl{
		Title:"harrison-15.3.2",
		Ucl:485.5,
		Ar:0.007,
		Em:1.4e8,
		Sw:0.5,
		Xr:350.,
		Yr:100.,
		Alp:0.000011,
		Flh:0.003,
		Flv:0.002,
		Fls:0.001,
		Frh:0.003,
		Frv:0.002,
		Frs:0.,
		Tl:160.0,
		Theta:67.0,
		Tr:50,
		Verbose:true,
	}
	_ = c.CblAz()
	
	c = Cbl{
		Title:"harrison-15.3.3",
		Ucl:485.5,
		Ar:0.007,
		Em:1.4e8,
		Sw:0.5,
		Xr:350.,
		Yr:100.,
		Alp:0.000011,
		Flh:0.003,
		Flv:0.002,
		Fls:0.001,
		Frh:0.003,
		Frv:0.002,
		Frs:0.,
		Tl:300.0,
		Theta:70.0,
		Tr:50,
		Lds:[][]float64{
			{100,0,150,1},
			{100,0,250,2},
		},
		Verbose:true,
	}
	_ = c.CblAz()
}

func TestCblDzCs(t *testing.T){
	c := Cbl{
		Title:"harrison-15.4.1",
		Ucl:500.,
		Em:1.4e8,
		Dens:71.5,
		Xr:350.,
		Yr:100.,
		Alp:0.000011,
		Flh:0.,
		Flv:0.,
		Fls:0.,
		Frh:0.,
		Frv:0.,
		Frs:0.,
		Tr:0.,
		Lds:[][]float64{
			{100,0,150,1},
			{100,0,250,1},
		},
		Strs:40000.,
		Tl:300.0,
		Theta:65.0,
		Verbose:true,
	}
	_ = c.CblDzCs()

	//PLS DO THIS
	// t.Log("starting ex. 2")
	// c = Cbl{
	// 	Title:"harrison-15.4.2",
	// 	Ucl:500.,
	// 	Em:1.4e8,
	// 	Dens:71.5,
	// 	Xr:350.,
	// 	Yr:100.,
	// 	Alp:0.000011,
	// 	Flh:0.,
	// 	Flv:0.,
	// 	Fls:0.,
	// 	Frh:0.,
	// 	Frv:0.,
	// 	Frs:0.,
	// 	Tr:0.,
	// 	Lds:[][]float64{
	// 		{100,0,150,2},
	// 		{100,0,250,2},
	// 	},
	// 	Strs:40000.,
	// 	Tl:230.,
	// 	Theta:60.,
	// 	Verbose:true,
	// }
	// _ = c.CblDzCs()
	
}
