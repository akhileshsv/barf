package barf

import (
	"os"
	"fmt"
	"testing"
	"path/filepath"
	kass"barf/kass"
)

func TestPsoBasic(t *testing.T){
	var w, c1, c2 float64
	var nd, ng, np int
	var mx, mn []float64
	var drw, trm, title string
	var inp []interface{}
	var par bool
	par = true
	fmt.Println("rastagirin 3d")
	fmt.Println("dims-3, min- 0,0,0")
	w = 0.5; c1 = 1.0; c2 = 1.5
	mx = []float64{10.,10.,10.}
	mn = []float64{-10.,-10.,-10.}
	nd = 3; ng = 100; np = 50
	drw = "gen"; trm = "dumb"; title = "rasta 3d func"
	Psoloop(false,false,par, w, c1, c2, np, ng, nd, mx,mn, rastapso, inp, drw, trm, title)
	
	fmt.Println("sphere func")
	fmt.Println("dims-3, min- 0,0,0")
	w = 0.5; c1 = 1.0; c2 = 2.0
	mx = []float64{10.,10.,10.}
	mn = []float64{-10.,-10.,-10.}
	nd = 3; ng = 100; np = 50
	drw = "gen"; trm = "dumb"; title = "sphere func"
	Psoloop(false,false,par, w, c1, c2, np, ng, nd, mx,mn, spherepso, inp, drw, trm, title)
	
}


func TestTrsPso(t *testing.T) {
	var w, c1, c2 float64
	var nd, ng, np int
	var mx, mn []float64
	var drw, trm, title string
	var inp []interface{}
	var par bool
	var examples = []string{"raka2"}
	var ndims []int
	par = true
	//var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for idx, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := kass.JsonInp(fname)
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
		switch idx{
			case 0:	
			for i := 0; i < 10; i++{
				mx = append(mx, 41.0)
				mn = append(mn, 0.0)
				ndims = append(ndims, 42)
			}
			inp = append(inp, 25.0); inp = append(inp, 2.0); inp = append(inp, 0.1)
			nd = 10
		}
		drw = "gen"; trm = "dumb"; title = ex
		w = 0.4; c1 = 2.0; c2 = 2.0
		ng = 100; np = 100
		drw = "gen"; trm = "dumb"; title = "raka truss"
		gpos, gb, pltstr := Psoloop(false,false,par, w, c1, c2, np, ng, nd, mx,mn,trsrakapso, inp, drw, trm, title)
		fmt.Println(gpos, gb, pltstr)
 		//ipos := make([]int, nd)
		//for i, val := range gpos{
		//	ipos[i] = int(val)
		//}
		//fmt.Println("starting sa run")
		//temp := 12.0; dmp := 1.0
		//ng = 500; mt := 3; ndchk := true
		//bsaloop(ndchk, temp, dmp, gb, ipos, ndims, ng, mt, sabtrsobj, inp, drw, trm, title)		
	}
}
