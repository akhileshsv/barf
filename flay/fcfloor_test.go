package barf

import (
	"os"
	"testing"
	"path/filepath"
)

func TestFcflr(t *testing.T){
	var examples = []string{"fpln",}
	var rezstring string
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/flay")
	t.Log(ColorPurple,"testing fc floor calc examples\n",ColorReset)
	for i, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		f, err := ReadFcflr(fname)
		f.Term = "qt"
		if err != nil{
			t.Fatal(err)
		}
		f.Load()
		f.Draw()
	}

}

