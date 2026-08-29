package barf

import (
	"os"
	"fmt"
	"testing"
	"path/filepath"
)

func TestNlCalcTrs2d(t *testing.T){
	//var examples = []string{"akms10.2"}
	//var rezstring string
	//var frmrez []interface{}
	var examples = []string{}
	var tol float64
	//rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for _, ex := range examples {
		fmt.Println("prob->",ex)
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		NlCalcTrs2d(mod, tol)
	}
}
