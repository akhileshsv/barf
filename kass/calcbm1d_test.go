package barf

import (
	"fmt"
	"testing"
	"path/filepath"
	//"runtime"
	"os"
	//"path"
)

func TestBmJsd(t *testing.T){
	var examples = []string{"akms7.3","akms7.6"}
	var rezstring string
	var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	t.Log(ColorPurple,"starting beam support displacement test\n",ColorReset)
	for i, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		t.Log(ColorBlue,"ex->",i+1,"file->",fname,"\n",ColorReset)
		frmrez,_ = CalcBm1d(mod,2)
		dglb, _ := frmrez[2].([]float64)
		rnode,_ := frmrez[3].([]float64)
		report,_ := frmrez[6].(string)
		rezstring += fmt.Sprintf("%s\n",ex)
		rezstring += fmt.Sprintf("%.5f\n",dglb)
		rezstring += fmt.Sprintf("%.5f\n",rnode)
		t.Log(report)
	}
	wantstring := `
akms7.3
[0.00000 0.00000 0.00000 -0.00195 0.00000 -0.00906 0.00000 0.03256]
[58.69196 76.51190 121.46692 130.55427 49.28685]
akms7.6
[0.00000 0.00000 0.00000 0.00036 0.00000 -0.00145 0.00000 0.00545]
[0.24303 17.49849 -0.97214 3.40248 -2.67338]
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("beam support displacements test failed")
	}

}

func TestBm1d(t *testing.T) {
	var examples = []string{"akms5.6","akms5.7","akms5.8"}
	var rezstring string
	var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for _, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		frmrez,_ = CalcBm1d(mod,2)
		dglb, _ := frmrez[2].([]float64)
		rnode,_ := frmrez[3].([]float64)
		rezstring += fmt.Sprintf("%s\n",ex)
		rezstring += fmt.Sprintf("%.5f\n",dglb)
		rezstring += fmt.Sprintf("%.5f\n",rnode)
	}
	wantstring := `
akms5.6
[0.00000 0.00000 0.00000 0.00203 0.00000 -0.00162 0.00000 0.00000]
[18.12500 1150.00000 12.98611 11.38889 17.50000 -800.00000]
akms5.7
[0.00000 0.00000 -0.00447 0.00056 0.00000 -0.00068 0.00000 0.00323]
[146.32691 281.18656 243.46484 50.20825]
akms5.8
[0.00000 -0.00056 0.00000 -0.00172 0.00000 0.00162 0.00000 0.00000]
[-9.64354 29.69812 45.25999 -5.31057 272.42281]
`
	if rezstring != wantstring {
		t.Errorf("Beam analysis test failed")
	}
}
