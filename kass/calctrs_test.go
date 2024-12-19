package barf

import (
	"fmt"
	"testing"
	"path/filepath"
	//"runtime"
	"os"
	//"path"
)

func TestTrsJsd(t *testing.T){
	var examples = []string{"akms7.2","akms7.5"}
	var rezstring string
	var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	t.Log(ColorPurple,"testing truss support displacement examples\n",ColorReset)
	for i, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		_, mod,_ := JsonInp(fname)
		frmrez,_ = CalcTrs(mod,mod.Ncjt)
		//ms, _ := frmrez[1].(map[int]*Mem)
		dglb, _ := frmrez[2].([]float64)
		rnode,_ := frmrez[3].([]float64)
		report,_ := frmrez[6].(string)
		rezstring += fmt.Sprintf("%s\n",ex)
		rezstring += fmt.Sprintf("%.5f\n",dglb)
		rezstring += fmt.Sprintf("%.5f\n",rnode)
		t.Log(report)
	}
	wantstring := `
akms7.2
[0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.33333 -0.14431]
[-49.04171 -65.38895 0.00000 130.77790 49.04171 -65.38895]
akms7.5
[0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.28068 -0.20193]
[-31.12542 -41.50056 0.00000 183.00113 -118.87458 158.49944]
`
	if rezstring != wantstring{
		//fmt.Println(rezstring)
		t.Errorf("truss support displacement calc test failed")
	}
}

func TestCalcTrs(t *testing.T) {
	var examples = []string{"akms3.8","akms3.9","akms4.1","akms8.1"}
	var rezstring string
	var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	t.Log(ColorPurple,"testing truss calc examples\n",ColorReset)
	for i, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		if i == 3 {
			frmrez,_ = CalcTrs(mod,3)
		} else {
			frmrez,_ = CalcTrs(mod,2)
		}
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		dglb, _ := frmrez[2].([]float64)
		rnode,_ := frmrez[3].([]float64)
		report,_ := frmrez[6].(string)
		rezstring += fmt.Sprintf("%s\n",ex)
		rezstring += fmt.Sprintf("%.5f\n",dglb)
		rezstring += fmt.Sprintf("%.5f\n",rnode)
		t.Log(report)
	}
	wantstring := `
akms3.8
[0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.21552 -0.13995]
[-10.06201 -13.41601 0.00000 126.83202 -139.93799 186.58399]
akms3.9
[0.00000 0.00000 0.00000 0.00000 0.00000 -9.18855 12.83651 -9.58441]
[-0.57761 320.82925 -298.38583 479.17075 -501.03657]
akms4.1
[0.00000 0.00000 0.07457 -0.20253 0.11362 0.00000 0.10487 0.00000 0.05782 -0.15268 0.02834 -0.07924]
[-25.00000 26.30140 112.34581 -3.64721]
akms8.1
[0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.10913 -0.12104 -0.57202]
[23.61617 47.23234 15.74411 -19.44191 77.76766 25.92255 -5.55809 -22.23234 7.41078 1.38383 -2.76766 0.92255]
`
	if rezstring != wantstring {
		t.Errorf("Truss analysis test failed")
	}
}
