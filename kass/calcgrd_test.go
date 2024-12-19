package barf

import (
	"fmt"
	"testing"
	"path/filepath"
	//"runtime"
	"os"
	//"path"
)

func TestGrd(t *testing.T) {
	var examples = []string{"akms8.2"}
	var rezstring string
	var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for _, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		frmrez,_ = CalcGrd(mod,3)
		dglb, _ := frmrez[2].([]float64)
		rnode,_ := frmrez[3].([]float64)
		rezstring += fmt.Sprintf("%s\n",ex)
		rezstring += fmt.Sprintf("%.5f\n",dglb)
		rezstring += fmt.Sprintf("%.5f\n",rnode)
	}
	wantstring := `
akms8.2
[0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 -0.05595 0.01133 -0.00549 0.00000 0.00000 0.00000]
[0.01469 -50.66168 59.13979 144.66845 -445.05882 7.99072 135.31686 -12.37832 375.52188]
`
	if rezstring != wantstring {
		t.Errorf("Space grid analysis test failed")
	}
}
