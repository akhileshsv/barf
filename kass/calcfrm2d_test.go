package barf

import (
	"fmt"
	"testing"
	"path/filepath"
	//"runtime"
	"os"
	//"path"
)

func TestFrm2dJsd(t *testing.T){
	var examples = []string{"akms7.1","akms7.4","akms7.7"}
	var rezstring string
	var frmrez []interface{}
	var err error
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	t.Log(ColorPurple,"starting plane frame support displacement test\n",ColorReset)
	for i, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod, _ := JsonInp(fname)
		t.Log(ColorBlue,"ex->",i+1,"file->",fname,"\n",ColorReset)
		frmrez, err = CalcFrm2d(mod,mod.Ncjt)
		if err != nil{
			fmt.Println(ColorRed, err, ColorReset)
		}
		dglb, _ := frmrez[2].([]float64)
		rnode,_ := frmrez[3].([]float64)
		report,_ := frmrez[6].(string)
		rezstring += fmt.Sprintf("%s\n",ex)
		rezstring += fmt.Sprintf("%.5f\n",dglb)
		rezstring += fmt.Sprintf("%.5f\n",rnode)
		t.Log(report)
	}
	wantstring := `
akms7.1
[0.00000 0.00000 0.00000 3.58003 -0.01212 0.00000 3.57103 -0.03011 -0.00166 0.00000 0.00000 -0.02149]
[-33.02445 21.52445 5045.86772 0.00000 -15.97555 53.47555]
akms7.4
[0.00000 0.00000 0.00000 0.01776 -1.05991 0.00074 0.00000 0.00000 0.00000]
[25.32274 97.42274 1431.70956 -25.32274 22.57676 -1536.99576]
akms7.7
[0.00000 0.00000 0.00000 -0.12199 -0.24965 0.00000 -0.00533 -0.00035 0.00053 0.00000 0.00000 0.00000]
[0.61417 -0.61417 -147.40164 0.00000 -0.61417 0.61417 0.00000]
`
	if rezstring != wantstring{
		//fmt.Println(rezstring)
		t.Errorf("frame (2d) support displacement test failed")
	}
}

func TestFrm2d(t *testing.T) {
	var examples = []string{"akms6.5","akms6.6","akms6.7","akms6.8"}
	var rezstring string
	var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	t.Log(ColorPurple,"starting plane frame support displacement test\n",ColorReset)
	for i, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		frmrez,_ = CalcFrm2d(mod,3)
		t.Log(ColorBlue,"ex->",i+1,"file->",fname,"\n",ColorReset)
		dglb, _ := frmrez[2].([]float64)
		rnode,_ := frmrez[3].([]float64)
		report,_ := frmrez[6].(string)
		t.Log(report)
		rezstring += fmt.Sprintf("%s\n",ex)
		rezstring += fmt.Sprintf("%.5f\n",dglb)
		rezstring += fmt.Sprintf("%.5f\n",rnode)
	}
	wantstring := `
akms6.5
[0.00000 0.00000 0.00000 0.00096 -0.00060 0.02539 0.00000 0.00000 0.02120]
[-126.81504 56.83140 222.80163 -113.18496 18.16860]
akms6.6
[0.00000 0.00000 0.00000 0.02130 -0.06732 -0.00255 0.00000 0.00000 0.00000]
[30.37118 102.08611 1215.98711 -30.37118 17.91339 -854.08504]
akms6.7
[0.00000 0.00000 0.00000 0.00000 0.00000 0.00000 0.18542 0.00042 -0.01762 0.18552 -0.00013 -0.02603 0.18662 0.00071 0.01789]
[-106.05122 -157.02711 360.44137 -85.94878 49.02711 320.31469]
akms6.8
[0.00000 0.00000 0.00000 3.44723 -0.00917 -0.01951 3.95204 -1.31523 0.00706 4.42471 -0.02116 -0.00927 0.00000 0.00000 -0.02302]
[-67.35552 33.01415 13788.66099 -33.50143 76.19511]
`
	if rezstring != wantstring {
		t.Errorf("Plane frame analysis test failed")
	}
}
