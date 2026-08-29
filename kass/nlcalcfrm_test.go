package barf

import (
	"os"
	"testing"
	"path/filepath"
)

func TestNlCalcFrm2d(t *testing.T){
	var examples = []string{"kima4.3","kima4.4"}
	var tol float64
	tol = 0.1
	//rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for i, ex := range examples {
		if i != 0{continue}
		t.Log("testing non linear frame analysis ex. no",i,"file",ex)
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		NlCalcFrm2d(mod, tol)
	}
}

