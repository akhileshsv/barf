package barf

import (
	"testing"
)

func TestBltDac(t *testing.T){
	t.Log("testing dac duggal 14.4")
	b := Blt{
		Title:"dug14.4",
		Vdu:385e3,
		Dia: 24.0,
		Grade:4.6,
		Name:"dac",
		Cloc:"web",
		Ctyp:2,
		Fup:410.0,
		T1:6.7,
		T2:9.4,
		T3:10.0,
		Mdims:[][]float64{
			{150,300,9.4,6.7},
			{150,450,17.4,9.4},
		},
		Mdxs:[]string{"b1","b2"},
		Verbose:false,
		Term:"qt",
	}
	BltDz(&b)
	t.Log(b.Report)
	t.Log("testing dac sub 5.16")
	
	b = Blt{
		Title:"sub5.16",
		Vdu:140e3,
		Dia: 20.0,
		Grade:4.6,
		Name:"dac",
		Cloc:"flange",
		Ctyp:1,
		Fup:410.0,
		T1:6.7,
		T2:9.4,
		T3:10.0,
		Mdims:[][]float64{
			{150,300,9.4,6.7},
			{200,200,15,9},
		},
		Mdxs:[]string{"b1","c1"},
		Verbose:true,
		Term:"qt",
	}
	BltDz(&b)
	t.Log(b.Report)
	t.Fatal()
	
	b = Blt{
		Title:"bhav8.1",
		Vdu:300e3,
		Dia: 20.0,
		Grade:4.6,
		Name:"dac",
		Cloc:"web",
		Ctyp:2,
		Fup:410.0,
		T1:6.7,
		T2:9.4,
		T3:10.0,
		Mdims:[][]float64{
			{150,300,9.4,6.7},
			{150,450,17.4,9.4},
		},
		Mdxs:[]string{"b1","b2"},
		Verbose:true,
		Term:"qt",
	}
}

func TestBltFp(t *testing.T){
	b := Blt{
		Title:"sub5.18",
		Vdu:140e3,
		Dia: 20.0,
		Grade:8.8,
		Name:"fp",
		Cloc:"flange",
		Ctyp:1,
		T1:6.7,
		T2:15.0,
		T3:10.0,
		Dmem:400.0,
		Verbose:true,
	}
	BltDz(&b)
}

func TestBltFep(t *testing.T){
	b := Blt{
		Title:"sub5.17",
		Vdu:140e3,
		Dia: 20.0,
		Grade:4.6,
		Name:"fep",
		Cloc:"flange",
		Ctyp:1,
		T1:6.7,
		T2:15.0,
		T3:6.0,
		Verbose:true,
	}
	BltDz(&b)
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
