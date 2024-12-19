package barf


import (
	"os"
	//"fmt"
	"path/filepath"
	"testing"
)


func TestRcFrm2dOpt(t *testing.T){
	fname := "frm2d_rajeev2.json"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	filename := filepath.Join(datadir, fname)
	RcFrm2dOpt(filename)
	t.Errorf("frame 2d opt test failed")
	
}
