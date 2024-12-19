package barf

import (
	"os"
	"testing"
	"path/filepath"
	kass"barf/kass"
)

func TestOptTrsIEC(t *testing.T){
	//test via pso, ga for ncycle runs
	var examples = []string{"raka2mks", "raka3mks"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/truss")
	opts := []int{12, 22}
	namez := []string{"ga","pso"}
	// opts := []int{12, 22}
	// namez := []string{"ga","pso"}
	NCycles = 9
	for i, ex := range examples{
		if i == 0{continue}
		for j, opt := range opts{	
			t.Logf("starting ex no. %v - %s opt - %s\n",i+1, ex, namez[j])
			fname := filepath.Join(datadir,ex+".json")
			_, mod,_ := kass.JsonInp(fname)
			mod.Web = true
			mod.Opt = opt
			mrez, err := OptTrsMod(*mod)
			t.Logf("opt finished with err-%v\n weight - %f\n",err, mrez.Weight)
			t.Logf("sections - %.2f\n",mrez.Cp)
		}
	}
}

func TestOptTrsMod(t *testing.T){
	var examples = []string{"raka2", "raka3"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/truss")
	for i, ex := range examples{
		if i != 0{continue}
		t.Logf("starting ex no. %v - %s",i, ex)
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := kass.JsonInp(fname)
		mod.Web = true
		mod.Opt = 12
		mrez, err := OptTrsMod(*mod)
		t.Logf("opt finished with err-%v\n-%s",err, mrez.Report)
		t.Log("txtplots\n",mrez.Txtplots)
	}
}

func TestOptTrsPso(t *testing.T){
	var examples = []string{"raka2", "raka3"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/truss")
	for i, ex := range examples{
		if i != 0{continue}
		t.Logf("starting ex no. %v - %s",i, ex)
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := kass.JsonInp(fname)
		mod.Opt = 23
		mod.Term = "svg"
		mod.Web = true
		mrez, err := OptTrsMod(*mod)
		t.Logf("opt finished with err-%v\n-%s",err, mrez.Report)
		
	}
}

func TestOptTrsTopo(t *testing.T){
	var examples = []string{"raka2", "raka3"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/truss")
	for i, ex := range examples{
		if i != 0{continue}
		t.Logf("starting ex no. %v - %s",i, ex)
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := kass.JsonInp(fname)
		mod.Opt = 2
		mod.Term = "svg"
		mod.Web = true
		mrez, err := trstopopso(*mod)
		t.Logf("opt finished with err-%v\n-%s",err, mrez.Report)
		
	}
}
