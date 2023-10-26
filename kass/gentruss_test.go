package barf

import (
	"fmt"
	"testing"
	//"path/filepath"
	//"runtime"
	//"os"
	//"path"
)

func TestGenTrsSwiss(t *testing.T){
	//kingpost, queenpost etc
	//var rezstring string
	var trs *Trs2d
	//fmt.Println("fink a type truss abel")
	fmt.Println("swiss roof test")
	trs = &Trs2d{
		Typ:       2,
		Cfg:       1,
		Mtyp:      3,
		Span:      3660.0,
		Spacing:   2400.0,
		Slope:     1.5,
		Purlinspc: 2500,
		//Height:    2400,
		DL:        160,
		Clctyp:    0,
		Vb:        50,
		Frmtyp:    2,
		//Bent:      true,
		Term:      "dumb",
		//Bcslope:   1.0/12.0,
	}
	//GenTruss(trs)
	trs.Calc()
}

func TestGenWlTrs(t *testing.T){
	var trs *Trs2d
	//fmt.Println("fink a type truss abel")
	fmt.Println("swiss roof test")
	trs = &Trs2d{
		Typ:       2,
		Cfg:       1,
		Mtyp:      3,
		Span:      18000.0,
		Spacing:   3400.0,
		Slope:     11.43,
		Purlinspc: 6000,
		//Height:    2400,
		DL:        1000,
		Clctyp:    0,
		Vb:        44.42,
		Frmtyp:    2,
		//Bent:      true,
		Term:      "dumb",
	}
	//GenTruss(trs)
	trs.Calc()
}

func TestGenTrs(t *testing.T) {
	var rezstring string
	var trs *Trs2d
	//fmt.Println("fink a type truss abel")
	fmt.Println("civil girl ex")
	trs = &Trs2d{
		Typ:       2,
		Cfg:       2,
		Mtyp:      1,
		Span:      10000.0,
		Spacing:   4500.0,
		Slope:     5,
		Purlinspc: 2000,
		Height:    8000,
		DL:        3.0,
		LL:        0.75,
		Clctyp:    0,
		Vb:        50,
		Frmtyp:    2,
		Cmdz:      []string{"bent","addld","env","dl,ll,wl","iscode","opt","pso"},
		Bent:      true,
		Term:      "qt",
	}
	trs.Calc()
	fmt.Println("a type")
	//trs.Typ++
	//GenTruss(trs, term)
	//fmt.Println("trap l type")
	//trs.Typ++
	//GenTruss(trs,term)
	//fmt.Println("trap l type")
	//wdl, wll = calcldt2d(lrftr, lspan, spacing, dl, ll float64, slfwt, mtyp, ftyp, rfmat,trstyp int)
	wantstring := `#$#^&@#`
	if rezstring != wantstring {
		t.Errorf("Truss generation test failed")
	}
	//trs.Typ = 3
	//GenTruss(trs, term)
}
