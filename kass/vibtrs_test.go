package barf

import (
	"os"
	"fmt"
	"testing"
	"path/filepath"
)

func TestVibTrs2d(t *testing.T){
	var examples = []string{"ross2.1"}
	//var rezstring string
	//var frmrez []interface{}
	//rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for _, ex := range examples {
		fmt.Println("prob->",ex)
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		VibTrs2d(mod)
	}
	t.Errorf("truss vibration analysis test failed")
}

