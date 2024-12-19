package barf

import (
	"fmt"
	"testing"
)

func TestPrlnDz(t *testing.T){
	t.Log("testing purlin design")
	t.Log("duggal ex 9.9 i section")
	var tspc, pspc, theta, dl, ll, pwl float64
	sname := "i"
	nsecs := 5
	tspc = 6000.0
	pspc = 1500.0
	theta = 30.0
	dl = 130.0
	ll = 0.0
	pwl = 2000.0
	PrlnDz800(nsecs, sname, tspc, pspc, theta, dl, ll, pwl)
	t.Log("maity lec 55 i section")
	tspc = 5000.0
	dl = 120.0
	pspc = 2000.0
	pwl = 1500.0
	PrlnDz800(nsecs, sname, tspc, pspc, theta, dl, ll, pwl)
}

func TestBmDz(t *testing.T){
	t.Log("testing beam design redux")
	t.Log("duggal ex 9.4 ss beam")
	b := Bm{
		Lspan:4000.0,
		DL:15.0,
		Sname:"i",
		Nsecs:1,
		Dtyp:1,
		Code:1,
		Sdx:50,
		Lsb:true,
	}
	err := BmDz(&b)
	if err != nil{
		t.Fatal(err)
	}

	t.Log("duggal ex 9.5 dtyp 0")
	b = Bm{
		Lspan:6000.0,
		Sname:"i",
		Nsecs:1,
		Dtyp:0,
		Code:1,
		Sdx:32,
		Lsb:true,
		Verbose:false,
		Vu:210*1e3,
		Mu:150*1e6,
		Dmax:0.001,
	}
	err = BmDz(&b)
	if err != nil{
		fmt.Println(err)
	}
	
	
	// b := Bm{
	// 	Ldcases:[][]float64{
	// 		{1,3,100,0,0,0,1},
	// 		{1,3,160,0,2,1,2},
	// 	},
	// 	Lspan:6.0,
	// 	Ly:6.0,
	// 	Ty:1.0,
	// 	Lbr:200.0,
	// 	Tbr:20.0,
	// 	Sname:"ub",
	// 	Nsecs:8,
	// 	Brchk:true,
	// 	Yeolde:true,
	// 	Endc:1,
	// 	Dtyp:1,
	// 	Grd:43,
	// 	Code:2,
	// }
	// b.Spam = true
	// b.Verbose = true
	// err := BmDz(&b)
	// if err != nil{
	// 	fmt.Println(err)
	// }
}

//WORTHLESS (now)
func TestStlBmDBs(t *testing.T) {
	//var rezstring string
	ldcases := [][]float64{
		{1,3,100,0,0,0,1},
		{1,3,160,0,2,1,2},
	}
	lspan := 6.0
	ly := 6.0
	ty := 1.0
	lbr := 200.0
	tbr := 20.0
	nsecs := 5
	grd := 43
	sname := "ub"
	brchck := true
	yeolde := true
	StlBmDBs(lspan, ly, ty, lbr, tbr, ldcases, sname, grd, nsecs, brchck, yeolde)	
	for i := 0; i < 10; i++{
		fmt.Println("****")
	}
	ldcases = [][]float64{
		{1,3,3.01,0,0,0,1},
		{1,3,3.6,0,0,0,2},
	}
	lspan = 9.045
	ly = 9.045
	ty = 1.0
	lbr = 200.0
	tbr = 20.0
	nsecs = 5
	grd = 43
	brchck = false
	yeolde = false
	StlBmDBs(lspan, ly, ty, lbr, tbr, ldcases, sname, grd, nsecs, brchck, yeolde)
}
