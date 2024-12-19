package barf

import (
	"os"
	"fmt"
	"testing"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func TestCalcModRez(t *testing.T){
	var examples = []string{"aksa5.13","aksa5.14","aksa5.6"}
	//,"aksa5.14"
	var rezstring string
	//var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/kass")
	for i, ex := range examples{
		if i != 1{continue}
		fname := filepath.Join(datadir,ex+".json")
		err := ModInp(fname, "qt", false)
		if err != nil{
			t.Errorf("2d frame b.m and sf test failed")
		}
	}
	
	wantstring := `@`
	t.Logf(rezstring)
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("2d frame b.m and sf test failed")
	}	
}

func TestCalcModJson(t *testing.T){
	var examples = []string{"akms3.9"}
	//,"aksa5.14"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/")
	for _, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		frmtyp, mod, err := JsonInp(fname) 
		
		if err != nil{
			t.Errorf("model json output test failed")
		}
		err = CalcMod(mod, frmtyp, "", false)
		if err != nil{
			t.Errorf("model json output test failed")
		}
		
		file, _ := json.MarshalIndent(mod, "", " ")
		outfile := filepath.Join(dirname,"../data/out/",ex+"_out.json")
		_ = ioutil.WriteFile(outfile, file, 0644)
	}
}
