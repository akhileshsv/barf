package barf

import (
	//"fmt"
	kass"barf/kass"
	"testing"
	"path/filepath"
	"os"
)

func TestTrussOptRaPa(t *testing.T) {
	var examples = []string{"raka2"}
	var params = [][]float64{{2.0,25.0,0.1,10.0}}
	//var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for idx, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := kass.JsonInp(fname)
		dmax, pmax, dens := params[idx][0],params[idx][1],params[idx][2]
		nsecs := int(params[idx][3])
		TrussOptRaPa(mod,mod.Dims,dmax,pmax,dens,nsecs,1)
	}
}
