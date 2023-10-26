package barf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCBmOpt(t *testing.T){
	var examples = []string{"govr1","govr2"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/cbeam")
	for i, ex := range examples{
		if i < 1 {continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		cb, err := ReadCBm(fname)
		if err != nil{
			t.Errorf("cbeam opt test failed")
		}
		//cb.Opt = 2
		CBmOpt(cb)
	}
}
