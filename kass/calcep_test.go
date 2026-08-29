package barf

import (
	"os"
	"fmt"
	"path/filepath"
	"testing"
)

func TestCalcEp(t *testing.T){
	var examples = []string{"hson13.3","hson13.8","hson13.9"}
	//var rezstring string
	//var frmrez []interface{}
	//rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for i, ex := range examples {
		//if i != 0{continue}
		fname := filepath.Join(datadir,ex+".json")
		fmt.Println(ColorCyan,i+1,"-","example->",fname,ColorReset)
		_, mod,_ := JsonInp(fname)
		CalcEpFrm(mod)
		//rezstring += fmt.Sprintf("%s\n",ex)
		//rezstring += fmt.Sprintf("%.5f\n",dglb)
		//rezstring += fmt.Sprintf("%.5f\n",rnode))
	}
	t.Errorf("wot wot")
}
