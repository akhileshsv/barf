package barf

import (
	"os"
	"path/filepath"
	"fmt"
	"time"
	"sort"
	"math/rand"
	"testing"
	kass"barf/kass"
)

func TestMut(t *testing.T){
	//rand.Seed(time.Now().UnixNano())
	rand.Seed(42)
	var b Bat
	var nd []int
	var mt int
	var ndchk bool
	nd = []int{7,7,7,7,7,7}
	
	b = Bat{Pos:[]int{1,2,3,4,5,6}}
	fmt.Println(ColorCyan,"flip and shuffle mutation",ColorReset)
	fmt.Println("parent")
	fmt.Println(b.Pos)
	b.Mut(ndchk, nd, mt)
	fmt.Println("child")
	fmt.Println(ColorGreen,b.Pos,ColorReset)

	mt++
	
	b = Bat{Pos:[]int{1,2,3,4,5,6}}
	fmt.Println(ColorCyan,"flip mutation",ColorReset)
	fmt.Println("parent")
	fmt.Println(b.Pos)
	b.Mut(ndchk, nd, mt)
	fmt.Println("child")
	fmt.Println(ColorGreen,b.Pos,ColorReset)
	
	mt++
	b = Bat{Pos:[]int{1,2,3,4,5,6}}
	fmt.Println(ColorCyan,"shuffle mutation",ColorReset)
	fmt.Println("parent")
	fmt.Println(b.Pos)
	b.Mut(ndchk, nd, mt)
	fmt.Println("child")
	fmt.Println(ColorGreen,b.Pos,ColorReset)

	mt++
	b = Bat{Pos:[]int{1,2,3,4,5,6}}
	fmt.Println(ColorCyan,"exchange mutation",ColorReset)
	fmt.Println("parent")
	fmt.Println(b.Pos)
	b.Mut(ndchk, nd, mt)
	fmt.Println("child")
	fmt.Println(ColorGreen,b.Pos,ColorReset)

	
	mt++
	b = Bat{Pos:[]int{1,2,3,4,5,6}}
	fmt.Println(ColorCyan,"shift mutation",ColorReset)
	fmt.Println("parent")
	fmt.Println(b.Pos)
	b.Mut(ndchk, nd, mt)
	fmt.Println("child")
	fmt.Println(ColorGreen,b.Pos,ColorReset)

	mt++
	b = Bat{Pos:[]int{1,2,3,4,5,6}}
	fmt.Println(ColorCyan,"inversion mutation",ColorReset)
	fmt.Println("parent")
	fmt.Println(b.Pos)
	b.Mut(ndchk, nd, mt)
	fmt.Println("child")
	fmt.Println(ColorGreen,b.Pos,ColorReset)
	
	mt++
	for i := 0; i < 10; i++{
		b = Bat{Pos:[]int{1,2,3,4,5,6}}
		fmt.Println(ColorCyan,"power mutation",ColorReset)
		fmt.Println("parent")
		fmt.Println(b.Pos)
		b.Mut(ndchk, nd, mt)
		fmt.Println("child")
		fmt.Println(ColorGreen,b.Pos,ColorReset)
	}		
}

func TestCrs(t *testing.T){
	rand.Seed(time.Now().UnixNano())
	var ct, cn int
	var x, y, a, b Bat
	var pcrs float64
	var ndchk bool
	var nd []int
	pcrs = 1.0
	x = Bat{
		Pos:[]int{1,2,3,4,5},
	}
	y = Bat{
		Pos:[]int{10,20,30,40,50},
	}
	ct = 1
	fmt.Println(ColorYellow,"1 point crossover")
	fmt.Println(ColorWhite,"x-\n",x.Pos,"y-\n",y.Pos,ColorReset)
	a, b = Bcrs(x,y,ct,cn,pcrs,nd,ndchk)
	fmt.Println(ColorCyan,"a-\n",a.Pos,"b-\n",b.Pos,ColorReset)
	//x.Pos = append(x.Pos,[]int{6,7,8,9}...)
	//y.Pos = append(y.Pos,[]int{60,70,80,90}...)
	ct = 2; cn = 2
	fmt.Println(ColorYellow,cn,"point crossover")
	fmt.Println(ColorWhite,"x-\n",x.Pos,"y-\n",y.Pos,ColorReset)
	a, b = Bcrs(x,y,ct,cn,pcrs,nd,ndchk)
	fmt.Println(ColorCyan,"a-\n",a.Pos,"b-\n",b.Pos,ColorReset)

	ct = 3
	fmt.Println(ColorYellow,"uniform crossover")
	fmt.Println(ColorWhite,"x-\n",x.Pos,"y-\n",y.Pos,ColorReset)
	a, b = Bcrs(x,y,ct,cn,pcrs,nd,ndchk)
	fmt.Println(ColorCyan,"a-\n",a.Pos,"b-\n",b.Pos,ColorReset)

	x.Pos = append(x.Pos,[]int{6,7,8,9}...)
	y.Pos = append(y.Pos,[]int{60,70,80,90}...)

	ct = 4
	fmt.Println(ColorYellow,"ordered crossover")
	fmt.Println(ColorWhite,"x-\n",x.Pos,"y-\n",y.Pos,ColorReset)
	a, b = Bcrs(x,y,ct,cn,pcrs,nd,ndchk)
	fmt.Println(ColorCyan,"a-\n",a.Pos,"b-\n",b.Pos,ColorReset)

	x = Bat{
		Pos:[]int{1,2,3,4,5},
	}
	y = Bat{
		Pos:[]int{0,2,4,6,3},
	}

	ct = 5
	fmt.Println(ColorYellow,"laplace crossover")
	fmt.Println(ColorWhite,"x-\n",x.Pos,"y-\n",y.Pos,ColorReset)
	a, b = Bcrs(x,y,ct,cn,pcrs,nd,ndchk)
	fmt.Println(ColorCyan,"a-\n",a.Pos,"b-\n",b.Pos,ColorReset)

}

func TestCrsr(t *testing.T){
	rand.Seed(time.Now().UnixNano())
	var x, y, a, b Rat
	var ct int
	var pcrs, alp float64
	var mx, mn, stp []float64
	x = Rat{
		Pos:[]float64{1,2,3,4,5},
	}
	y = Rat{
		Pos:[]float64{10,20,30,40,50},
	}
	ct = 1; pcrs = 1.0; alp = 0.5
	fmt.Println(ColorYellow,"linear crossover")
	fmt.Println(ColorWhite,"x-\n",x.Pos,"y-\n",y.Pos,ColorReset)
	a,b =  Rcrs(x, y, ct, pcrs, alp, mx, mn, stp)
	fmt.Println(ColorCyan,"a-\n",a.Pos,"b-\n",b.Pos,ColorReset)

	ct = 2; pcrs = 1.0; alp = 0.5
	fmt.Println(ColorYellow,"blend crossover")
	fmt.Println(ColorWhite,"x-\n",x.Pos,"y-\n",y.Pos,ColorReset)
	a,b =  Rcrs(x, y, ct, pcrs, alp, mx, mn, stp)
	fmt.Println(ColorCyan,"a-\n",a.Pos,"b-\n",b.Pos,ColorReset)
}

func TestGaUnemax(t *testing.T){
	var drw, trm, title string
	var np, ng, mt, ct, cn, st, mx int
	var pmut, pcrs float64
	var nd []int
	var inp []interface{}
	var f Bobj
	var ndchk, par, web bool
	fmt.Println("one max problem dims 5")
	par = true
	np = 20; ng = 50; mt = 3; ct = 4; st = 4
	pmut = 0.02; pcrs = 0.80
	nd = []int{2,2,2,2,2}
	f = unemax
	drw = "gen"; trm = "dumb"; title = "unemax"
	_, _, _, _, err := Gabloop(web, par, ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, f, drw, trm, title)
	if err != nil{fmt.Println(err)}
}

func TestGaKnpsck(t *testing.T){
	items := []string{
		"lptp","book","radio","tv","potato","brick","bottle","cam","phone","pic","flwr","chair","wtch","bts","radiatr","tablet","printr",
	}
	vals := []float64{300,15,30,230,7,1,2,280,500,170,5,4,500,30,25,450,170}
	wts := []float64{3,2,1,6,5,3,1,0.5,0.1,1,2,3,0.05,1.5,5,0.5,4.5}
	var inp []interface{}
	inp = append(inp, vals)
	inp = append(inp, wts)
	var drw, trm, title string
	var np, ng, mt, ct, cn, st, mx int
	var pmut, pcrs float64
	var nd []int
	var f Bobj
	var ndchk, par, web bool
	fmt.Println("knapsack problem")
	par = true
	np = 8; ng = 20; mt = 1; ct = 1; st = 1
	pmut = 0.2; pcrs = 0.70
	nd = make([]int, len(items))
	for i := range items{
		nd[i] = 2
	}
	f = knpsck
	drw = "gen"; trm = "dumb"; title = "knapsack"
	gpos, _, _, _, _ := Gabloop(web, par, ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, f, drw, trm, title)
	wt := 0. 
	for i := range items{
		if gpos[i] == 1{
			fmt.Println("knapsack contains->",items[i])
			wt += wts[i]
		}
	}
	fmt.Println("weight->",wt)
}


func TestARaka(t *testing.T) {
	var examples = []string{"raka2"}
	//var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	var drw, trm, title string
	var np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step int
	var pmut, pcrs, gap float64
	var nd []int
	var inp []interface{}
	var f Bobj
	var ndchk, par bool
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
		mx = -1
		np = 45; ngmn = 300; ngmx = 500; mt = 6; ct = 5; st = 1; ndchk = true
		npmn = 20; npmx = 150
		pmut = 0.02; pcrs = 0.75; gap = 0.01; step = 10
		switch idx{
			case 0:	
			for i := 0; i < 10; i++{
				nd = append(nd, 42)
			}
			inp = append(inp, 25.0); inp = append(inp, 2.0); inp = append(inp, 0.1)
		}
		f = trsrakaobj
		drw = "gen"; trm = "dumb"; title = ex
		_, _, _, _, err := Gabaloop(false, par, ndchk, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, pmut, pcrs, gap, nd, inp, f, drw, trm, title)
		if err != nil{fmt.Println(err)}
	}
}


func TestRakaTruss(t *testing.T) {
	var examples = []string{"raka2"}
	//var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	var drw, trm, title string
	var np, ng, mt, ct, cn, st, mx int
	var pmut, pcrs, gb, gw, bw float64
	var nd, gpos []int
	var inp []interface{}
	var f Bobj
	var ndchk, par bool
	var err error
	par = true
	mtz := map[int]string{
		1:"flip",
		2:"shuffle",
		3:"exchange",
		4:"shift",
		5:"invert",
	}
	ctz := map[int]string{
		1:"1 pt",2:"n pt",3:"uniform",4:"ordered",
	}
	stz := map[int]string{
		1:"tour",2:"fprop",3:"stoch",4:"rank",5:"raka",
	}
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
		mx = -1
		np = 50; ng = 50; mt = 1; ct = 1; st = 1
		fmt.Println(ColorRed)
		fmt.Println("params->",ColorCyan,"mut-",mtz[mt],"\tcrs-",ctz[ct],"\tsel-",stz[st])
		fmt.Println(ColorReset)
		pmut = 0.1; pcrs = 0.85
		switch idx{
			case 0:	
			for i := 0; i < 10; i++{
				nd = append(nd, 42)
			}
			inp = append(inp, 25.0); inp = append(inp, 2.0); inp = append(inp, 0.1)
			ndchk = true
		}
		f = trsrakaobj
		drw = "gen"; trm = "dumb"; title = ex
		bPos := make([]int, 10)
		copy(bPos,gpos)
		bb := -100.0
		for i := 0; i < 5; i++{
			gpos, gw, gb, _, err = Gabloop(false, par,ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, f, drw, trm, title)
			if bb < gb {bb = gb; bw = gw; copy(bPos, gpos)}
		}
		
		if err != nil{
			fmt.Println(err)
		} else {
			fmt.Println("best ga val")
			fmt.Println(bPos, bb, bw)
			//temp := 10.0; dmp := 1.0
			//ng = 5000; mt = 1
			//bsaloop(ndchk, temp, dmp, gb, gpos, nd, ng, mt, sabtrsobj, inp, drw, trm, title)		
		}
		
	}
}


func TestTrsStew(t *testing.T) {
	var examples = []string{"raka2"}
	//var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	var drw, trm, title string
	var np, ng, cn, mx int
	var pmut, pcrs, gb float64
	var nd, gpos []int
	var inp []interface{}
	var f Bobj
	var par, ndchk bool
	var err error
	par = true
	mtz := map[int]string{
		1:"flip",
		2:"shuffle",
		3:"exchange",
		4:"shift",
		5:"invert",
	}
	ctz := map[int]string{
		1:"1 pt",2:"n pt",3:"uniform",4:"ordered",
	}
	stz := map[int]string{
		1:"tour",2:"fprop",3:"stoch",4:"rank",
	}
	fname := filepath.Join(datadir,examples[0]+".json")
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
	fmt.Println(ColorCyan,"truss opt khichdi example",ColorWhite)
	mx = -1
	np = 25; ng = 30
	pmut = 0.2; pcrs = 0.8
	for i := 0; i < 10; i++{
		nd = append(nd, 42)
	}
	inp = append(inp, 25.0); inp = append(inp, 2.0); inp = append(inp, 0.1)
	f = trsrakaobj
	drw = "gen"; trm = "dumb"; title = "raka2"
	fitmap := make(map[float64][]int)
	var fitz []float64
	for mt := range mtz{
		for ct := range ctz{
			for st := range stz{
				fmt.Println("xxx---xxx---xxx")
				fmt.Println(ColorRed)
				fmt.Println("params->",ColorCyan,"mut-",mtz[mt],"\tcrs-",ctz[ct],"\tsel-",stz[st])
				fmt.Println(ColorReset)
				gpos, _, gb, _, err = Gabloop(false,par, ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, f, drw, trm, title)
				if err != nil{
					fmt.Println(err)
				} else {
					fmt.Println(gpos)
					fmt.Println(gb)
				}
				fitmap[gb] = []int{mt, ct, st}
				fitz = append(fitz, gb)
			}
		}
	}
	sort.Slice(fitz, func(i,j int) bool{
		return fitz[i]<fitz[j]
	})
	for i, fit := range fitz{
		pz := fitmap[fit]
		fmt.Println(i, fit, "params->",ColorCyan,"mut-",mtz[pz[0]],"\tcrs-",ctz[pz[1]],"\tsel-",stz[pz[2]])
	}
}

