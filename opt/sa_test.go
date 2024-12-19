package barf

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	kass"barf/kass"
)

func TestSa(t *testing.T){
	var temp, dmp float64
	var nd, ng int
	var mx, mn []float64
	var drw, trm, title string
	temp = 10.0; dmp = 1.0
	mx = []float64{5.0}
	mn = []float64{-5.0}
	nd = 1; ng = 2000
	drw = "gen"; trm = "dumb"; title = "sphere"
	fmt.Println("sim annealing 1d sphere func")
	Saloop(sphere, temp, dmp, mx, mn, nd, ng, drw, trm, title)
	
	temp = 12.0; dmp = 1.0
	mx = []float64{5.0,5.0}
	mn = []float64{-5.0,-5.0}
	nd = 2; ng = 2000
	drw = "gen"; trm = "dumb"; title = "sphere"
	fmt.Println("sim annealing 2d sphere func")
	Saloop(sphere, temp, dmp, mx, mn, nd, ng, drw, trm, title)

	fmt.Println("sim annealing 2d drop wave func")
	Saloop(dwave, temp, dmp, mx, mn, nd, ng, drw, trm, title)

	fmt.Println("sim annealing 2d styblinski tang func")
	Saloop(stang, temp, dmp, mx, mn, nd, ng, drw, trm, title)
	
}

func TestSaTruss(t *testing.T){
	ex := "raka2"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	fname := filepath.Join(datadir,ex+".json")
	_, mod,_ := kass.JsonInp(fname)
	var inp []interface{}
	var nd []int
	var ndchk bool
	var ng, mt int
	modr := kass.Model{
		Ncjt:mod.Ncjt,
		Coords:mod.Coords,
		Supports:mod.Supports,
		Mprp:mod.Mprp,
		Jloads:mod.Jloads,
		Em:mod.Em,
	}
	inp = append(inp, modr)
	inp = append(inp, mod.Dims)
	fmt.Println(ColorCyan,"truss opt example ",ColorWhite,ex)
	for i := 0; i < 10; i++{
		nd = append(nd, 42)
	}
	inp = append(inp, 25.0); inp = append(inp, 2.0); inp = append(inp, 0.1)
	drw := "gen"; trm := "dumb"; title := "raka2"
	fmt.Println("starting sa run")
	temp := 12.0; dmp := 1.0
	ng = 500; mt = 3
	Bsaloop(ndchk, temp, dmp, 0.0, []int{}, nd, ng, mt, sabtrsobj, inp, drw, trm, title)		

}
