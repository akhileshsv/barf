package barf

import (
	"os"
	"path/filepath"
	"fmt"
	"testing"
	kass "barf/kass"
)

func TestBmAsvrat(t *testing.T){
	//hulse asv/sv rat test
	var rezstring string
	var b *RccBm
	b = &RccBm{
		Fck:25.0,
		Fy:250.0,
		Bw:300.0,
		Dused:600.0,
		Cvrt:50.0,
		Ast:982.0,
		Code:1,
	}
	var vf, asvrat float64
	vf = 196.0
	asvrat = AsvRatBs(b, vf)
	rezstring += fmt.Sprintf("hulse ex 3.3.6 vf - %f kn asv/sv - %f\n",vf, asvrat)
	wantstring := `hulse ex 3.3.6 vf - 196.000000 kn asv/sv - 0.905072
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("beam asv/sv ratio (bs) test failed")
	}
}

func TestBmTorDBs(t *testing.T){
	//very basic, needs more tests
	var rezstring string
	var b *RccBm
	b = &RccBm{
		Fck:25.0,
		Fy:250.0,
		Styp:6,
		Bf:800,
		Bw:300,
		Dused:500,
		Df:200,
		Cvrt:50.0,
		Cvrc:60.0,
	}
	b.TorRect()
	//fmt.Println(b.Hdim)
	//fmt.Println(b.Ldim)

	err := BmTorDBs(b,0.86,160,13)
	if err != nil{
		fmt.Println(err)
		t.Errorf("beam torsion design (bs) test failed")
	}
	for i, sv := range b.Svr{
		rezstring += fmt.Sprintf("component %v asv/sv (shear) %.2f asv/sv (total) %.2f max. link spacing %.0f mm additional steel %.0f mm2\n",i+1, b.Avr[i][0],b.Avr[i][1],sv, b.Avr[i][2])
	}
	wantstring := `component 1 asv/sv (shear) 0.00 asv/sv (total) 0.24 max. link spacing 125 mm additional steel 95 mm2
component 2 asv/sv (shear) 0.86 asv/sv (total) 1.41 max. link spacing 200 mm additional steel 371 mm2
component 3 asv/sv (shear) 0.00 asv/sv (total) 0.24 max. link spacing 125 mm additional steel 95 mm2
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("beam torsion design (bs) test failed")
	}
}

func TestBmShrDz(t *testing.T){
	//uses cbeam B-'
	//let's see
	var rezstring string
	var examples = []string{"shah5.6.1","shah5.6.2"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/cbeam/")
	var bmenv map[int]*kass.BmEnv
	for i, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		cb, err := ReadCBm(fname)
		if err != nil{
			t.Errorf("beam shear design test failed")
		}
		bmenv, err = CBeamEnvRcc(&cb, cb.Term, true)
		if err != nil{
			fmt.Println(err)
			t.Errorf("beam shear design test failed")
		}
		CBmDz(&cb,bmenv)
		bm := cb.RcBm[0][1]
		bm.Table(false)
		t.Log(bm.Report)
		//b.Lspc = []float64{slink, smin, snom, mainlen, minlen, nomlen}
		rezstring += fmt.Sprintf("link spacing - main %.f mm min %.f mm nominal %.f mm\nlength - main %.2f m min %.2f m nominal %.2f m\n",bm.Lspc[0],bm.Lspc[1],bm.Lspc[2],bm.Lspc[3],bm.Lspc[4],bm.Lspc[5])
	}
	wantstring := `link spacing - main 230 mm min 300 mm nominal 300 mm
length - main 4.00 m min 4.00 m nominal 2.00 m
link spacing - main 250 mm min 250 mm nominal 250 mm
length - main 0.17 m min 2.21 m nominal 1.02 m
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("beam shear design test failed")
	}
}

/*
   ye olde funcs
   
func TestBmShrIs(t *testing.T){
	var rezstring string
	
	dims := []float64{230,750}
	bar := kass.CalcSecProp(1, dims)

	cp := [][]float64{{bar.Ixx*1e-12, bar.Area*1e-6}}
	
	mod := &kass.Model{
		Ncjt:2,
		Cmdz:[]string{"1db","mks","1"},
		Coords:[][]float64{{0},{10.0}},
		Supports:[][]int{{1,-1,0},{2,-1,0}},
		Em:[][]float64{{15e6}},
		Cp:cp,
		Dims:[][]float64{dims},
		Mprp:[][]int{{1,2,1,1,0}},
		Msloads:[][]float64{{1,3,52,0,0,0}},
	}
	rezstring += "shah 5.5.1\n"
	lsx := 230.0; rsx := 230.0
	b := &RccBm{
		Fck:15.0,
		Fy:415.0,
		Bf:0.0,
		Df:140.0,
		Bw:230.0,
		Dused:750.0,
		Ast:2945.0,
		Asc:0.0,
		L0:0.0,
		Tyb:0.0,
		Cvrt:62.5,
		Cvrc:0.0,
		Lbd:0.0,
		Lspan:10.0,
		Endc:1,
		DM:0.0,
		Code:1,
		Lsx:lsx,
		Rsx:rsx,
	}
	frmrez, _ := kass.CalcBm1d(mod, 2)
	bmresults := kass.CalcBmSf(mod, frmrez, false)
	bmr := bmresults[1]
	idxs, nlegs, spacing, err := BmShrDIs(b, lsx, bmr.Xs, bmr.SF, false)
	if err != nil{fmt.Println("err-<,>",err)}
	for idx,didx := range idxs{
		rezstring += fmt.Sprintf("Bar dia %0.0f %v -legged stirrups @ %0.2f mm c-c \n",StlDia[didx],nlegs[idx],spacing[idx])
	}
	if 1 == 0{fmt.Println(rezstring)}
	//BmShrDz(b, lsx, rsx, bmr.Xs, bmr.SF)
}

*/
