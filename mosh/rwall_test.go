package barf

import (
	"os"
	"fmt"
	"path/filepath"
	"testing"
)

func TestRwallOpt(t *testing.T){
	var examples = []string{"sub16.3"}
	for i, ex := range examples{
		dirname,_ := os.Getwd()
		datadir := filepath.Join(dirname,"../data/examples/mosh/rwall/")
		fname := filepath.Join(datadir,ex+".json")
		fmt.Println(ColorCyan,"example->","no.",i+1,"file->",fname,ColorReset)
		r, err := ReadRwall(fname)
		if err != nil{
			t.Fatal("file read for rwall failed")
		}
		_, err = RwallOpt(r)
		if err != nil{
			t.Fatal("error in rwall opt-",err)
		}
	}
}
