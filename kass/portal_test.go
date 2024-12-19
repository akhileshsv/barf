package barf

import (
	"os"
	"fmt"
	"path/filepath"
	"testing"
)

func TestPortalInit(t *testing.T){
	var examples = []string{"saka2003-1","sp40"}
	var rezstring string
	var pf Portal
	var err error
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/kass/portal")
	t.Log(ColorPurple,"testing portal frame examples\n",ColorReset)
	for i, ex := range examples {
		if i != 0{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		pf, err = ReadPortal(fname)
		if err != nil{
			t.Logf(fmt.Sprintf("%s",err))
		}
		pf.Sdxs = []int{12,12}
		
		// pf.Gentyp = 2
		pf.Ndiv = 2
		err = pf.Calc()
		// cfg := []int{0,1,2,3}
		// for _, cf := range cfg{
		// 	pf.Config = cf
		// 	err = pf.Calc()
		// }
	}
	wantstring := ``
	if rezstring != wantstring{
		t.Errorf("portal frame test failed")
	}

}
