package barf

import (
	"os"
	"fmt"
	"testing"
	"path/filepath"
)

//okay. do this or bust
//truss load tests
func TestGenTrsLds(t *testing.T){
	var examples = []string{"sp64ex1","is2366.1","abel12.8","duggal15.1","bhav12.3","sub14.1"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/kass/truss")
	t.Log(ColorPurple,"testing truss load generation examples\n",ColorReset)
	var rezstring string
	for i, ex := range examples{
		//if i !=3{continue}
		fname := filepath.Join(datadir,ex+".json")
		typz := "timber"
		if i > 2{typz = "steel"}
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",typz," truss\n",ColorReset)
		trs, err := ReadTrs2d(fname)
		if err != nil{
			t.Fatal("error reading truss load generation example json-",fname, err)
		}
		trs.Spam = false
		err = trs.Init()
		if err != nil{
			t.Fatal(err,"truss load generation test failed")
		}
		err = trs.GenGeom()
		if err != nil{
			t.Log(err)
			t.Fatal(err,"truss load generation test failed")
		}
		err = trs.GenCp()
		if err != nil{
			t.Fatal(err,"truss load generation test failed")
		}
		err = trs.GenLd()
		if err != nil{
			t.Fatal(err,"truss load generation test failed")
		}
		rezstring += fmt.Sprintf("%s example %s\n",typz,ex)
		rezstring += fmt.Sprintf("truss nodal loads\ndl %.3f n ll %.3f n\nwind loads\npwl %.3f n pwr %.3f n pwpl %.3f n pwpr %.3f n\n",trs.Pdl, trs.Pll, trs.Pwl, trs.Pwr, trs.Pwp, trs.Pwp)
		tarea := trs.Purlinspc * trs.Spacing * 1e-6
		rezstring += fmt.Sprintf("wind pressures\npwl %.3f n/m2 pwr %.3f n/m2 pwp %.3f n/m2\nspacing\npurlin %.3f mm truss %.3f mm trib. area %.3f m2\n",trs.Pwl/tarea,trs.Pwr/tarea,trs.Pwp/tarea,trs.Purlinspc, trs.Spacing, tarea)
	}
	wantstring := ``
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Fatal("truss load generation test failed")
	}

}

//trussIEC24/v0 tests
func TestTrsCalcV0(t *testing.T){
	var examples = []string{"sp64.1","bhav12.3","duggal15.1","abel12.10","is2366"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/kass/truss")
	var rezstring string
	
	//exmap := make(map[int]float64)
	for i, ex := range examples{
		//if i == 0 || i == 3{continue}
		if i != 1{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log("example->",i+1,"file->",fname)
		trs, err := ReadTrs2d(fname)
		if err != nil{
			t.Fatal(err)
		}
		err = trs.Calc()
		if err != nil{
			t.Fatal(err)
		}
		// err = trs.Init()
		// if err != nil{
		// 	return
		// }
		// //t.Term = "dumb"
		// err = trs.GenGeom()
		// if err != nil{return}
		// err = trs.GenCp()
		// if err != nil{return}
		// PlotGenTrs(trs.Coords, trs.Ms, "qt")
	}
	wantstring := "nothingman"
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Fatal("truss load generation test failed")
	}
}

//truss purlin design tests
func TestPrlnDz(t *testing.T){
	
}

//stl truss load gen test
func TestStlTrsLds(t *testing.T){
	var examples = []string{"sp64.1","duggal15.1","bhav12.3","sub14.1"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/truss")
	var rezstring string
	
	//exmap := make(map[int]float64)
	for i, ex := range examples{
		//if i == 0 || i == 3{continue}
		if i != 2{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log("example->",i+1,"file->",fname)
		trs, err := ReadTrs2d(fname)
		if err != nil{
			t.Fatal(err)
		}
		err = trs.Calc()
		if err != nil{
			t.Fatal(err)
		}
		// err = trs.Init()
		// if err != nil{
		// 	return
		// }
		// //t.Term = "dumb"
		// err = trs.GenGeom()
		// if err != nil{return}
		// err = trs.GenCp()
		// if err != nil{return}
		// PlotGenTrs(trs.Coords, trs.Ms, "qt")
	}
	wantstring := "nothingman"
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Fatal("truss load generation test failed")
	}
}

func TestGenFunTrs(t *testing.T){
	t.Log("starting fun truss test")
	trs := &Trs2d{
		Typ:       6,
		Cfg:       1,
		Mtyp:      3,
		Span:      3200,
		Spacing:   2400,
		Rise:      3200,
		Purlinspc: 400,
		DL:        160,
		Vb:        50,
		Frmtyp:    2,
		Term:      "dumb",
	}
	//GenTruss(trs)
	trs.Calc()
}

func TestGenSwissTrs(t *testing.T){
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
		Ldcalc:    1,
		//Bent:      true,
		Term:      "dumb",
		Bcslope:   12.0,
	}
	//GenTruss(trs)
	trs.Calc()
}

func TestGenTrsCalc(t *testing.T) {
	var examples = []string{"is2366","cgrl","sub14.1","sp64ex1"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/kass/truss")
	t.Log(ColorPurple,"testing truss generation examples\n",ColorReset)
	for i, ex := range examples{
		if i != 0{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		trs, err := ReadTrs2d(fname)
		if err != nil{
			t.Fatal("error reading truss generation example json-",fname, err)
		}
		err = trs.Calc()
		if err != nil{
			t.Fatal("truss generation test failed")
		}
	}
}

func TestGenParTruss(t *testing.T){
	var rezstring string
	var trs *Trs2d
	//fmt.Println("fink a type truss abel")
	fmt.Println("civil girl ex")
	trs = &Trs2d{
		Typ:       5,
		Cfg:       9,
		Mtyp:      3,
		Span:      3600.0,
		Spacing:   3600.0,
		Slope:     0,
		Purlinspc: 600,
		Height:    3000,
		DL:        3000,
		LL:        750,
		Clctyp:    0,
		Vb:        50,
		Frmtyp:    2,
		Bent:      true,
		Term:      "qt",
	}
	trs.Calc()
	wantstring := `#$#^&@#`
	if rezstring != wantstring {
		t.Errorf("Truss generation test failed")
	}
	
}


/*
   YE OLDE

func TestGenTrs(t *testing.T) {
	//TEST FOR MAT, LOAD GENERATION AND CALCS HERE
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

*/
