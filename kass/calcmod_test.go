package barf

import (
	"os"
	"fmt"
	"testing"
	"path/filepath"
)

func TestCalcModRez(t *testing.T){
	var examples = []string{"aksa5.13"}
	//,"aksa5.14"
	var rezstring string
	//var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/kass")
	for _, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		err := ModInp(fname, "dumb")
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
