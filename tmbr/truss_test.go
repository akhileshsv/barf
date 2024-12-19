package barf

import (
	"os"
	"fmt"
	"testing"
	"path/filepath"
	kass"barf/kass"
)

func TestTrsGen(t *testing.T){
	var examples = []string{"abel12.8","is2366","scissor","fink"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/tmbr/truss")
	//exmap := make(map[int]float64)
	for i, ex := range examples{
		if i != 1{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log("example->",i+1,"file->",fname,"opt - PSO\n")
		trs, err := kass.ReadTrs2d(fname)
		if err != nil{
			t.Fatal(err)
		}
		TrussDz(&trs)
	}
}

// func TestTrsOpt(t *testing.T){
// 	var examples = []string{"is2366"}
// 	dirname,_ := os.Getwd()
// 	datadir := filepath.Join(dirname,"../data/examples/tmbr/truss")
// 	//exmap := make(map[int]float64)
// 	for i, ex := range examples{
// 		fname := filepath.Join(datadir,ex+".json")
// 		t.Log("example->",i+1,"file->",fname,"opt - PSO\n")
// 		trs, err := kass.ReadTrs2d(fname)
// 		if err != nil{
// 			t.Fatal(err)
// 		}
// 		//_, _ = TrsOpt(trs)
// 	}
// }

func TestTrussDzIs(t *testing.T){
	var trs *kass.Trs2d
	fmt.Println("is code roof test")
	trs = &kass.Trs2d{
		Id:        "is.2366-B",
		Cmdz:      []string{"2dt","mmks"},
		Typ:       2,
		Cfg:       3,
		Mtyp:      3,
		Span:      12000.0,
		Spacing:   2500.0,
		Slope:     2.0,
		Purlinspc: 1400,
		DL:        300.0,
		LL:        750.0,
		Clctyp:    0,
		Ldcalc:    1,
		Vb:        50,
		Frmtyp:    2,
		Term:      "qt",
	}
	//GenTruss(trs)
	TrussDz(trs)

}

func TestTrussSwiss(t *testing.T){
	var trs *kass.Trs2d
	fmt.Println("swiss roof test")
	trs = &kass.Trs2d{
		Id:        "swiss roof",
		Cmdz:      []string{"2dt","mmks"},
		Typ:       2,
		Cfg:       1,
		Mtyp:      3,
		Group:     2,
		Span:      5200.0,
		Spacing:   1960.0,
		Slope:     2600.0/1400.0,
		Purlinspc: 3000,
		DL:        600.0,
		LL:        750.0,
		Clctyp:    0,
		Ldcalc:    1,
		Vb:        50*0.8,
		Frmtyp:    2,
		Term:      "svg",
		Bcslope:   2600.0/400.0,
		Cpi:       0.7,
		Kostin:    []float64{1000.0+1500.0,600.0,200.0},
		Ntruss:    4.0,
	}
	//GenTruss(trs)
	err := TrussDz(trs)
	fmt.Println(err)
}

func TestTrussCompSwiss(t *testing.T){
	var trs *kass.Trs2d
	fmt.Println("comp swiss roof test")
	trs = &kass.Trs2d{
		Id:        "comp swiss roof",
		Cmdz:      []string{"2dt","mmks"},
		Typ:       2,
		Cfg:       1,
		Mtyp:      3,
		Group:     2,
		Span:      5200.0,
		Spacing:   1960.0,
		Slope:     2600.0/1400.0,
		Purlinspc: 2400,
		DL:        600.0,
		LL:        750.0,
		Clctyp:    0,
		Ldcalc:    1,
		Vb:        50*0.8,
		Frmtyp:    2,
		Term:      "svg",
		Bcslope:   2600.0/400.0,
		Cpi:       0.7,
		Kostin:    []float64{1000.0+1500.0,600.0,200.0},
		Ntruss:    4.0,
	}
	//GenTruss(trs)
	err := TrussDz(trs)
	fmt.Println(err)
}

func TestParTruss(t *testing.T){
	var trs *kass.Trs2d
	fmt.Println("parallel eaves test")
	trs = &kass.Trs2d{
		Id:        "swizz",
		Cmdz:      []string{"2dt","mmks"},
		Typ:       5,
		Cfg:       1,
		Mtyp:      3,
		Group:     2,
		Span:      5880.0,
		Spacing:   5220.0,
		Slope:     0.0,
		Purlinspc: 1960.0,
		DL:        600.0,
		LL:        750.0,
		Clctyp:    0,
		Ldcalc:    1,
		Vb:        50*0.8,
		Frmtyp:    2,
		Term:      "svg",
		Bcslope:   2600.0/400.0,
		Cpi:       0.7,
		Kostin:    []float64{1000.0+1500.0,600.0,200.0},
		Ntruss:    4.0,
	}
	//GenTruss(trs)
	err := TrussDz(trs)
	fmt.Println(err)
}

func TestTrussDzIFink(t *testing.T){
	var trs *kass.Trs2d
	fmt.Println("fink roof test")
	trs = &kass.Trs2d{
		Id:        "fink",
		Cmdz:      []string{"2dt","mmks"},
		Typ:       2,
		Cfg:       4,
		Mtyp:      3,
		Span:      3300.0,
		Spacing:   2500.0,
		Slope:     1.0/0.577,
		Purlinspc: 1400,
		DL:        300.0,
		LL:        750.0,
		Clctyp:    0,
		Ldcalc:    0,
		Vb:        50,
		Frmtyp:    2,
		Term:      "svg",
	}
	//GenTruss(trs)
	TrussDz(trs)

}
