package barf

import (
	"fmt"
	"math"
	opt"barf/opt"
)

//CBmOpt calls ga and pso routines to optimize a continuous beam 
func CBmOpt(cb CBm){
	var gbest float64
	switch cb.Opt{
		case 1:
		//pso
		var w, c1, c2 float64
		var nd, ng, np int
		var mx, mn []float64
		var drw, trm, title string
		var inp []interface{}
		var par bool
		//par = true
		inp = append(inp, cb)
		switch{
			case cb.Width == 0.0:
			//opt for width and depth, width is always constant (sue me)
			switch cb.Dconst{
				case true:
				nd = 2
				case false:
				nd = 1 + cb.Nspans
			}
			case cb.Width > 0.0:
			switch cb.Dconst{
				case true:
				nd = 1
				case false:
				nd = cb.Nspans
			}
		}
		for i := 0; i < nd; i++{
			mx = append(mx, 31.0)
			mn = append(mn, 0.0)
		}
		w = 0.4; c1 = 2.0; c2 = 2.0
		ng = 30; np = 20
		drw = "gen"; trm = "dumb"; title = cb.Title
		gpos, gb := opt.Psoloop(par, w, c1, c2, np, ng, nd, mx,mn, cbmfit, inp, drw, trm, title)
		gbest = gb
		cb = getcbm(gpos,inp)
		case 2:
		//ga
		var drw, trm, title string
		var np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step int
		var pmut, pcrs, gap float64
		var nds []int
		var inp []interface{}
		//var f opt.Bobj
		var ndchk, par bool
		var nd int
		inp = append(inp, cb)
		mx = -1
		np = 2; ngmn = 3; ngmx = 7; mt = 3; ct = 4; st = 4; ndchk = true
		npmn = 2; npmx = 4
		pmut = 0.02; pcrs = 0.75; gap = 0.01; step = 10
		
		inp = append(inp, cb)
		switch{
			case cb.Width == 0.0:
			//opt for width and depth, width is always constant (sue me)
			switch cb.Dconst{
				case true:
				nd = 2
				case false:
				nd = 1 + cb.Nspans
			}
			case cb.Width > 0.0:
			switch cb.Dconst{
				case true:
				nd = 1
				case false:
				nd = cb.Nspans
			}
		}
		for i := 0; i < nd; i++{
			nds = append(nds, 31)
		}
		drw = "gen"; trm = "dumb"; title = cb.Title
		_, _, err := opt.Gabaloop(par, ndchk, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, pmut, pcrs, gap, nds, inp, Cbmobj, drw, trm, title)
		if err != nil{fmt.Println(err)}
		
		case 3:
	}
	cb.Verbose = true
	cb.Term = "dumb"
	bmenv, err := CBeamEnvRcc(&cb, cb.Term, false)
	if err != nil{
		fmt.Println("ERRORE,errore->",err)
		return
	}
	err = CBmDz(&cb,bmenv)
	
	if err != nil{
		fmt.Println("ERRORE,errore->",err)
		return
	}
	for _, span := range cb.RcBm{
		//rezstring += fmt.Sprintf("span-%v\n",j+1)
		for _, bm := range span{
			if bm.Ignore{continue}
			//bm.Printz()
			bm.Table(true)
			//t.Log(bm.Report)
			//t.Log(bm.Report)
			bm.Draw()
		}
	}
	fmt.Println("total cost->",gbest)

}

//Cbmobj is an objective function for opt.Bat fitneess evaluation
func Cbmobj(b *opt.Bat, inp []interface{}) (err error){
	cb := gatcbm(b.Pos,inp)
	bmenv, err := CBeamEnvRcc(&cb, cb.Term, false)
	if err != nil{
		b.Fit = 1e6
		return
	}
	err = CBmDz(&cb,bmenv)
	if err != nil{
		b.Fit = 1e6
		return
	}
	for _, barr := range cb.RcBm{
		b.Fit += barr[1].Kost
	}
	if b.Fit == 0.0{b.Fit = 1e6}
	return
}

//gatcbm returns a CBm struct given an int slice position vector
func gatcbm(pos []int, inp []interface{})(cb CBm){
	//what absolute noobery USE T(Any)
	cb, _ = inp[0].(CBm)
	cb.Sections = make([][]float64,cb.Nspans)
	var bw, dused, step, mindim float64
	step = 25.0
	mindim = 225.0
	for i := range cb.Sections{
		cb.Sections[i] = make([]float64, 2)
		switch{
			case cb.Width == 0.0:
			bw = float64(pos[0]) * step + mindim
			if cb.Dconst{
				dused = float64(pos[1]) * step + mindim				
			} else {
				dused = float64(pos[i+1]) * step + mindim				
			}
			case cb.Width > 0.0:
			bw = cb.Width
			if cb.Dconst{
				dused = float64(pos[0]) * step + mindim				
			} else {
				dused = float64(pos[i]) * step + mindim				
			}
		}
		cb.Sections[i][0] = bw; cb.Sections[i][1] = dused 
	}
	fmt.Println(ColorCyan,"oot->",cb.Sections,ColorReset)
	return

}

//getcbm returns a CBm given a floating point position vector
func getcbm(pos []float64, inp []interface{})(cb CBm){
	cb, _ = inp[0].(CBm)
	cb.Sections = make([][]float64,cb.Nspans)
	var bw, dused, step, mindim float64
	step = 25.0
	mindim = 225.0
	for i := range cb.Sections{
		cb.Sections[i] = make([]float64, 2)
		switch{
			case cb.Width == 0.0:
			bw = math.Round(pos[0]) * step + mindim
			if cb.Dconst{
				dused = math.Round(pos[1]) * step + mindim				
			} else {
				dused = math.Round(pos[i+1]) * step + mindim				
			}
			case cb.Width > 0.0:
			bw = cb.Width
			if cb.Dconst{
				dused = math.Round(pos[0]) * step + mindim				
			} else {
				dused = math.Round(pos[i]) * step + mindim				
			}
		}
		cb.Sections[i][0] = bw; cb.Sections[i][1] = dused 
	}
	return
}

//cbmfit calcs a pos vec's(pso) fitness 
func cbmfit(pos []float64, inp []interface{}) (fit float64){
	cb := getcbm(pos,inp)
	bmenv, err := CBeamEnvRcc(&cb, cb.Term, false)
	if err != nil{
		fit = 1e6
		return
	}
	err = CBmDz(&cb,bmenv)
	if err != nil{
		fit = 1e6
		return
	}
	for _, barr := range cb.RcBm{
		fit += barr[1].Kost
	}
	if fit == 0.0{fit = 1e6}
	return
}
