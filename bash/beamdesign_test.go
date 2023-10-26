package barf

import (
	"fmt"
	"testing"
)

func TestBmDz(t *testing.T){
	var b Bm
	b = Bm{
		Ldcases:[][]float64{
			{1,3,100,0,0,0,1},
			{1,3,160,0,2,1,2},
		},
		Lspan:6.0,
		Ly:6.0,
		Ty:1.0,
		Lbr:200.0,
		Tbr:20.0,
		Styp:7,
		Nsecs:8,
		Brchk:true,
		Yeolde:true,
		Endc:1,
		Dtyp:1,
		Grd:43,
		Code:2,
	}
	b.Spam = true
	b.Verbose = true
	err := BmDz(&b)
	if err != nil{
		fmt.Println(err)
	}
}

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
	sectyp := 7
	brchck := true
	yeolde := true
	StlBmDBs(lspan, ly, ty, lbr, tbr, ldcases, sectyp, grd, nsecs, brchck, yeolde)	
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
	sectyp = 7
	brchck = false
	yeolde = false
	StlBmDBs(lspan, ly, ty, lbr, tbr, ldcases, sectyp, grd, nsecs, brchck, yeolde)

}
