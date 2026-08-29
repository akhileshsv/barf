package barf

import (
	"os"
	"testing"
	"path/filepath"
	kass"barf/kass"
)

func TestModOpt(t *testing.T){
	var examples = []string{"topping"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/tmbr/frame")
	for i, ex := range examples{
		t.Logf("starting ex no. %v - %s",i, ex)
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := kass.JsonInp(fname)
		mod.Opt = 12
		err := ModOpt(mod)
		t.Fatal(err)
	}

}
