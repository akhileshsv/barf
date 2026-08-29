package barf

import (
	"fmt"
	"testing"
	"path/filepath"
	//"runtime"
	"os"
	//"path"
)

func TestFrm3d(t *testing.T) {
	var examples = []string{"akms8.4"}
	var rezstring string
	var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for _, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		frmrez,_ = CalcFrm3d(mod,6)
		dglb, _ := frmrez[2].([]float64)
		rnode,_ := frmrez[3].([]float64)
		rezstring += fmt.Sprintf("%s\n",ex)
		rezstring += fmt.Sprintf("%.5f\n",dglb)
		rezstring += fmt.Sprintf("%.5f\n",rnode)
	}
	wantstring := `
akms8.4
[-0.00135 -0.00280 -0.00181 -0.00300 0.00106 0.00650 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000]
[5.37574 44.10629 -0.74272 2.17215 58.98735 2330.51966 -4.62491 11.11738 -6.46065 -515.54573 -0.76472 369.67165 -0.75082 4.77633 7.20338 -383.50156 -60.16642 -4.70199]
`
	if rezstring != wantstring {
		t.Errorf("Space frame analysis test failed")
	}
}
