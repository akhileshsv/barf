package barf

import (
	"testing"
)

func TestBltDzFrm(t *testing.T){
	
}

func TestBltDz(t *testing.T){
	t.Log("maity ex. 1 - a")
	b := Blt{
		Pu:150e3,
		Dia: 16.0,
		Grade:4.6,
		Ctyp:1,
		Ltyp:1,
		Tms:[]float64{10,18},
		Pt:0.0,
		Fup:410.0,
		Verbose:true,
		Bmem:200.0,
	}
	BltDz(&b)
	t.Log(ColorRed,"nbolts->",b.Nj*b.Ni,ColorReset)
	t.Log("maity ex. 1 - b")
	b = Blt{
		Pu:150e3,
		Dia: 16.0,
		Grade:4.6,
		Ctyp:2,
		Ltyp:1,
		Tms:[]float64{10,18},
		Pt:8.0,
		Fup:410.0,
		Verbose:true,
		Bmem:200.0,
	}
	BltDz(&b)
	t.Log(ColorRed,"nbolts->",b.Nj*b.Ni,ColorReset)
	
	t.Log("maity ex. 1 - c")
	b = Blt{
		Pu:150e3,
		Dia: 16.0,
		Grade:4.6,
		Ctyp:3,
		Ltyp:1,
		Tms:[]float64{10,18},
		Pt:8.0,
		Fup:410.0,
		Verbose:true,
		Bmem:200.0,
	}
	BltDz(&b)
	t.Log(ColorRed,"nbolts->",b.Nj*b.Ni,ColorReset)

	t.Log("maity ex. 1 - hsfg")
	b = Blt{
		Pu:150e3,
		Dia: 20.0,
		Grade:8.8,
		Btyp:2,
		Ctyp:1,
		Ltyp:1,
		Tms:[]float64{10,10},
		Pt:10.0,
		Pitch:60.0,
		Fup:410.0,
		Verbose:true,
		Slip:true,
		Bmem:110.0,
	}
	BltDz(&b)
	t.Log(ColorRed,"nbolts->",b.Nj*b.Ni,ColorReset)
	
	
	t.Fatal()

	// t.Log("duggal ex. 7.1 net area calc")
	// b := Blt{
	// 	Dia: 18.0,
	// 	Grade:4.6,
	// 	Ctyp:1,
	// 	Ltyp:1,
	// 	Bltyp:1,
	// 	Bmem:300,
	// 	Tmem:8.0,
	// 	Fup:410.0,
	// 	Verbose:true,
	// 	Ni:4,
	// 	Nj:6,
	// }
	// err := b.BltNsa()
	// if err != nil{
	// 	t.Fatal(err)
	// }
	
}

func TestBltDiaCalc(t *testing.T){
	t.Log("maity ex. 1 - a")
	b := Blt{
		Pu:150e3,
		Dia: 16.0,
		Grade:4.6,
		Ctyp:1,
		Tms:[]float64{10,18},
		Pt:8.0,
		Fup:410.0,
		Pitch:50.0,
		Verbose:true,
	}
	BltDiaCalc(&b)
	t.Log(ColorRed,"nbolts->",b.Nj*b.Ni)
	t.Fatal()
	t.Log("maity ex. 1 - b")
	b = Blt{
		Pu:150e3,
		Dia: 16.0,
		Grade:4.6,
		Ctyp:2,
		Tms:[]float64{10,18},
		Pt:8.0,
		Fup:410.0,
		Pitch:50.0,
		Verbose:true,
	}
	BltDiaCalc(&b)
	t.Log("maity ex. 1 - c")
	b = Blt{
		Pu:150e3,
		Dia: 16.0,
		Grade:4.6,
		Ctyp:3,
		Tms:[]float64{10,18},
		Pt:8.0,
		Fup:410.0,
		Pitch:50.0,
		Verbose:true,
	}
	BltDiaCalc(&b)
	t.Log("duggal ex. 5.1 - a")	
	b = Blt{
		Pu:150e3,
		Dia: 20.0,
		Grade:4.6,
		Ctyp:1,
		Tms:[]float64{12,12},
		Pt:10.0,
		Fup:410.0,
		Verbose:true,
	}
	BltDiaCalc(&b)
	t.Log("duggal ex. 5.1 - b")
	b = Blt{
		Pu:150e3,
		Dia: 20.0,
		Grade:4.6,
		Ctyp:2,
		Tms:[]float64{12,12},
		Pt:10.0,
		Fup:410.0,
		Verbose:true,
	}
	
	BltDiaCalc(&b)
	t.Log("duggal ex. 5.1 - c")	
	b = Blt{
		Pu:150e3,
		Dia: 20.0,
		Grade:4.6,
		Ctyp:3,
		Tms:[]float64{12,12},
		Pt:8.0,
		Fup:410.0,
		Verbose:true,
	}
	BltDiaCalc(&b)
}

func TestBoltSs(t *testing.T){
	b := &Blt{
		Bc:[][]float64{{0,0},{400,0},{0,500},{400,500},{0,800},{400,800}},
		Dias:[]float64{20,20,20,20,20,20},
		Typ:[]float64{1,1,1,1,1,1},
		Frc:[][]float64{{2,-40,0}},
		Fc:[][]float64{{1200,1000}},
		Print:true,
	}
	err := BoltSs(b)
	if err != nil{
		t.Errorf("bolt group analysis test failed")
	}
	t.Log("\n\n",b.Report)
	
}

func TestBltvec(t *testing.T){
	var ni, nj int
	// ni = 6; nj = 6
	// for lout := 1; lout < 3; lout ++{
	// 	bvec := bltvec(ni, nj, lout)
	// 	t.Log(bvec)
	// }
	ni = 5; nj = 3
	bvec := bltvec(ni, nj, 4)
	t.Log(bvec)
}
