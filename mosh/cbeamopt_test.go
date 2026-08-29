package barf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOptCBmIEC(t *testing.T){
	//test via pso, ga for ncycle runs

	var examples = []string{"govr1","govr2","govr3"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/cbeam")
	//var gmin, pmin []float64
	opts := []int{13, 23}
	namez := []string{"ad. ga","ad. pso"}
	NCycles = 9
	for i, ex := range examples{
		for j, opt := range opts{	
			if j > 0{continue}
			fname := filepath.Join(datadir,ex+".json")
			t.Logf("starting ex no. %v - %s opt - %s\n",i+1, ex, namez[j])
			cb, err := ReadCBm(fname)
			if err != nil{
				t.Errorf("file read error - cbeam opt IEC test failed")
			}
			cb.D1 = 0.0; cb.D2 = 0.0
			cb.Opt = opt
			cb.Web = true
			cb.Term = "svg"
			cb.Dconst = true
			cbrez, e := CBmOpt(cb)
			if e != nil{
				t.Logf("opt finished, error - %s",e)
			}
			t.Logf("optimized cost - %.4f rs",cbrez.Kost)
			t.Logf("volume of  rcc %.4f m3 weight of steel %.2f kg area of formwork %.3f m2\n",cbrez.Quants[0],cbrez.Quants[1],cbrez.Quants[2])
			
			t.Logf("optimized section vec - %.f rs",cbrez.Sections)
		}
	}
}


func TestCBmOpt(t *testing.T){
	var examples = []string{"govr1","govr2","govr3"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/cbeam")
	exmap := make(map[int]float64)
	for i, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		t.Log("example->",i+1,"file->",fname,"opt - PSO\n")
		cb, err := ReadCBm(fname)

		if err != nil{
			t.Errorf("cbeam opt test failed")
		}
		cb.D1 = 0.0; cb.D2 = 0.0
		cb.Opt = 11
		//cb.Web = true
		
		//cb.Dconst = true
		cb.Term = "dumb"
		//CBmOpt(cb)
		
		//t.Log("example->",i+1,"file->",fname,"opt - GA\n")
		//cb, err = ReadCBm(fname)
		//if err != nil{
		//	t.Errorf("cbeam opt test failed")
		//}
		//cb.Opt = 2
		
		cbrez, e := CBmOpt(cb)
		if e != nil{
			t.Log(e)
		}
		exmap[i+1] = cbrez.Kost
		//t.Log("final kwast- ",cbrez.Kost," rupees")
		// t.Log(cbrez.Report)
	}
	for i, val := range exmap{
		t.Log("example #- ",i,"name - ",examples[i-1], "kwast-",val," rupeeses")
	}
}
