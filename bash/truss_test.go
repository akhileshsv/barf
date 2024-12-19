package barf

import (
	"os"
	"testing"
	"path/filepath"
	kass"barf/kass"
)

func TestTrussDz(t *testing.T){
	var examples = []string{"duggal15.1"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/truss")
	//exmap := make(map[int]float64)
	for i, ex := range examples{
		if i != 0{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log("example->",i+1,"file->",fname)
		trs, err := kass.ReadTrs2d(fname)
		if err != nil{
			t.Fatal(err)
		}
		TrussDz(&trs)
	}
}
