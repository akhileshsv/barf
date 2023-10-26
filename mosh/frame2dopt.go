package barf

import (
	"fmt"
	"math"
	opt"barf/opt"
	kass"barf/kass"
)

//getd2d returns a 2d frame from a floating point position vector 
func getf2d(pos []float64, inp []interface{})(f kass.Frm2d){
	//gen f2d
	f, _ = inp[0].(kass.Frm2d)
	f.Kost = 0.0
	step := 25.0
	f.Sections = make([][]float64, f.Ncols+f.Nbms+f.Ncls)
	f.Csec = make([]int, f.Ncols)
	f.Bsec = make([]int, f.Nbms + f.Ncls)
	f.Verbose = false
	f.Term = ""
	switch f.Width{
		case 0.0:
		bw := math.Round(pos[0]/step) * step
		switch f.Dconst{
			case true:
			dbm := math.Round(pos[len(pos)-1]/step) * step
			for i, val := range pos[1:]{
				if i < f.Ncols{
					f.Sections[i] = make([]float64,2)
					f.Sections[i][0] = bw
					f.Sections[i][1] = math.Round(val/step) * step
					f.Styps[i] = f.Cstyp
					f.Csec[i] = i + 1
				}
			}
			for i := range f.Sections[f.Ncols:]{
				f.Sections[i+f.Ncols] = make([]float64, 2)
				f.Sections[i+f.Ncols][0] = bw
				f.Sections[i+f.Ncols][1] = dbm
				f.Styps[i+f.Ncols] = f.Bstyp
				f.Bsec[i] = i + f.Ncols + 1
			}
			case false:
			for i, val := range pos[1:]{
				f.Sections[i] = make([]float64,2)
				f.Sections[i][0] = bw
				f.Sections[i][1] = math.Round(val/step) * step 
				switch{
					case i < f.Ncols:
					f.Csec[i] = i + 1
					f.Styps[i] = f.Cstyp
					default:
					f.Bsec[i-f.Ncols] = i+1
					f.Styps[i] = f.Bstyp
				}
			}
		}
		default:
		//f.Verbose = true
		bw := f.Width
		switch f.Dconst{
			case true:
			dbm := math.Round(pos[len(pos)-1]/step) * step
			for i, val := range pos{
				if i < f.Ncols{
					f.Sections[i] = make([]float64,2)
					f.Sections[i][0] = bw
					f.Sections[i][1] = math.Round(val/step) * step 
					f.Csec[i] = i + 1
				}
			}
			for i := range f.Sections[f.Ncols:]{
				f.Sections[i+f.Ncols] = make([]float64, 2)
				f.Sections[i+f.Ncols][0] = bw
				f.Sections[i+f.Ncols][1] = dbm 
				f.Bsec[i] = i + f.Ncols + 1
			}
			case false:
			for i, val := range pos{
				f.Sections[i] = make([]float64,2)
				f.Sections[i][0] = bw
				f.Sections[i][1] = math.Round(val/step) * step 
				switch{
					case i < f.Ncols:
					f.Csec[i] = i + 1
					default:
					f.Bsec[i-f.Ncols] = i+1
				}
			}
		}
	}
	
	//fmt.Println(f.Csec)
	//fmt.Println(f.Bsec)
	//fmt.Println(f.Sections)
	return
}

//Frm2dOpt opts a 2d frame (entry func)
//IT IS NOT WERK
func Frm2dOpt(f kass.Frm2d) (frez kass.Frm2d, err error){
	//w const
	//pos := make([]int, 1+f.Ncols+f.Nbms)
	//all sep
	//err = f.Init()
	//if err != nil{return}
	//var gbest float64
	//fmt.Println(f.Ncols, f.Nbms, f.Ncls)
	//return
	//term := f.Term
	f.Init()
	term := f.Term
	f.Term = ""
	f.Verbose = false
	switch f.Opt{
		case 1:
		//pso
		var w, c1, c2 float64
		var nd, ng, np int
		var mx, mn []float64
		var drw, trm, title string
		var inp []interface{}
		var par bool
		//var allcons int
		par = true
		inp = append(inp, f)
		switch f.Width{
			case 0.0:
			switch f.Dconst{
				case true:
				//pos[0] - width, pos[1] - c1H, pos[xstep] - cnH, pos[xstep+1] - bdepth
				nd = 1 + f.Ncols + 1
				case false:
				//pos[0] - width, pos[1] - c1H, pos[xstep] - cnH, pos[xstep+1]
				nd = 1 + f.Ncols + f.Nbms + f.Ncls
			}
			
			for i := 0; i < nd; i++{
				mx = append(mx, 750.0)
				mn = append(mn, 250.0)
			}
			default:
			switch f.Dconst{
				case true:
				nd = f.Ncols + 1
				case false:
				nd = f.Ncols + f.Nbms + f.Ncls
			}
			
			for i := 0; i < nd; i++{
				mx = append(mx, math.Round(f.Width * 3.0/25.0)*25.0)
				mn = append(mn, math.Round(f.Width * 1.0/25.0)*25.0)
			}
		}
		w = 0.5; c1 = 1.0; c2 = 1.0
		
		//HERE
		ng = 30; np = 9
		drw = "gen"; trm = "dumb"; title = f.Title
		
		gpos, _ := opt.Psoloop(par, w, c1, c2, np, ng, nd, mx,mn, f2dfit, inp, drw, trm, title)
		f.Term = term
		f = getf2d(gpos,inp)
		f.Init()
		//f.Term = "qt"
		err, _, _ = Frm2dDz(&f)
		//jl, ml := f.Mod.SumFrcs()
		if err != nil{
			//fmt.Println("no. of warnings/errors->",allcons)
		} else{
			f.Spam = true
			fmt.Println("yee-haw")
			jl, ml := f.Mod.SumFrcs()
			f.Mod.Msloads = ml
			f.Mod.Jloads = jl
			f.DrawMod(f.Term)
			f.Mod.Msloads = f.Loadcons[1]
			f.Mod.Jloads = f.Jloadmap[1]
			jl, ml = f.Mod.SumFrcs()
			f.Mod.Msloads = ml
			f.Mod.Jloads = jl
			f.DrawMod(f.Term)
			//f.Term = "mono"
			Frm2dTable(&f, false)
			//fmt.Println("global best vector->",gb)
			//Frm2dTable(&f, true)
			//f.DrawMod("dumb")
			//fmt.Println(f.Txtplots[0])
			//for _, flrvec := range f.CBeams{
			//	f.Txtplots = append(f.Txtplots, PlotBmEnv(f.Bmenv, flrvec, "qt"))
			//}
		}

		//fmt.Println(par, w, c1, c2, np, ng, nd, mx,mn,  inp, drw, trm, title,f)
		case 2:
		var drw, trm, title string
		var np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step int
		var pmut, pcrs, gap float64
		var nds []int
		var inp []interface{}
		//var f opt.Bobj
		var ndchk, par bool
		var nd int
		mx = -1
		np = 2; ngmn = 3; ngmx = 7; mt = 3; ct = 4; st = 4; ndchk = true
		npmn = 2; npmx = 4
		pmut = 0.02; pcrs = 0.75; gap = 0.01; step = 10
		inp = append(inp, f)
		switch f.Width{
			case 0.0:
			switch f.Dconst{
				case true:
				//pos[0] - width, pos[1] - c1H, pos[xstep] - cnH, pos[xstep+1] - bdepth
				nd = 1 + f.Ncols + 1
				case false:
				//pos[0] - width, pos[1] - c1H, pos[xstep] - cnH, pos[xstep+1] - bdepth
				nd = 1 + f.Ncols + f.Nbms + f.Ncls
			}
			default:
			switch f.Dconst{
				case true:
				nd = f.Ncols + 1
				case false:
				nd = f.Ncols + f.Nbms + f.Ncls
			}
		}
		for i := 0; i < nd; i++{
			nds = append(nds, 15)
		}
		drw = "gen"; trm = "dumb"; title = f.Title
		_, _, err := opt.Gabaloop(par, ndchk, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, pmut, pcrs, gap, nds, inp, f2dobj, drw, trm, title)
		if err != nil{fmt.Println(err)}

	}
	frez = f
	return 
}

//f2dobj is the objective function for evaluating an opt.Bat ([]int) as a slice of sl
func f2dobj(b *opt.Bat, inp []interface{}) (err error) {
	vec := make([]float64, len(b.Pos))
	step := 25.0
	mindim := 225.0
	for i := range b.Pos{
		vec[i] = float64(b.Pos[i]) * step + mindim
	}
	b.Fit = f2dfit(vec, inp)
	return
}

//f2dfit is the objective function for evaluating an opt.Bat ([]int) as a slice of sl
func f2dfit(pos []float64, inp []interface{}) (fit float64){
	f := getf2d(pos,inp)
	f.Init()
	_, allcons, _ := Frm2dDz(&f)
	//Frm2dDz(f *kass.Frm2d) (err error, allcons int, emap map[int][]error){
	
	if f.Kost <= 0.0{
		fit = 1e12
	} else {
		fit = f.Kost * (1.0 + 10.0 * float64(allcons))
	}
	return
}


/*
func Frm2dOpt(f kass.Frm2d){
	var gbest float64
	//pos - [bw, cd1,cd2,cdn, bd1, bd2, bdn]
	//or - [cd1, cd2, cd3, bd1, bd2, bd3
	//if dconst - [bw, cd1, cd2, cdn, bd]
	switch f.Opt{
		case 1:
		//pso
		var w, c1, c2 float64
		var nd, ng, np int
		var mx, mn []float64
		var drw, trm, title string
		var inp []interface{}
		var par bool
		//par = true
		inp = append(inp, f)
		switch{
			case f.Width == 0.0:
			//opt for width and depth, width is always constant (sue me)
			switch f.Dconst{
				case true:
				//
				case false:
				//d varies n.spans
				
			}
			case f.Width > 0.0:
			switch f.Dconst{
				case true:
				nd = 1
				case false:
				nd = f.Nspans
			}
		}
		for i := 0; i < nd; i++{
			mx = append(mx, 31.0)
			mn = append(mn, 0.0)
		}
		w = 0.4; c1 = 2.0; c2 = 2.0
		ng = 30; np = 20
		drw = "gen"; trm = "dumb"; title = f.Title
		gpos, gb := opt.Psoloop(par, w, c1, c2, np, ng, nd, mx,mn, fmfit, inp, drw, trm, title)
		gbest = gb
		f = getfm(gpos,inp)
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
		inp = append(inp, f)
		mx = -1
		np = 2; ngmn = 3; ngmx = 7; mt = 3; ct = 4; st = 4; ndchk = true
		npmn = 2; npmx = 4
		pmut = 0.02; pcrs = 0.75; gap = 0.01; step = 10
		
		inp = append(inp, f)
		switch{
			case f.Width == 0.0:
			//opt for width and depth, width is always constant (sue me)
			switch f.Dconst{
				case true:
				nd = 2
				case false:
				nd = 1 + f.Nspans
			}
			case f.Width > 0.0:
			switch f.Dconst{
				case true:
				nd = 1
				case false:
				nd = f.Nspans
			}
		}
		for i := 0; i < nd; i++{
			nds = append(nds, 31)
		}
		drw = "gen"; trm = "dumb"; title = f.Title
		_, _, err := opt.Gabaloop(par, ndchk, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, pmut, pcrs, gap, nds, inp, Fmobj, drw, trm, title)
		if err != nil{fmt.Println(err)}
		
		case 3:
	}
	f.Verbose = true
	f.Term = "dumb"
	bmenv, err := FeamEnvRcc(&f, f.Term, false)
	if err != nil{
		fmt.Println("ERRORE,errore->",err)
		return
	}
	err = kass.F2dDz(&f,bmenv)
	
	if err != nil{
		fmt.Println("ERRORE,errore->",err)
		return
	}
	for _, span := range f.Cbm{
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
*/
