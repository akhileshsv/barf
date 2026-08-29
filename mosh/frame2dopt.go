package barf

import (
	"fmt"
	opt"barf/opt"
	kass"barf/kass"
)

//frm2dga optimizes a frame via good ol' ga
func frm2dga(f kass.Frm2d) (frez kass.Frm2d, err error){
	
	var drw, trm, title string
	var ndims, np, ng, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step int
	var pmut, pcrs, gap float64
	var nd []int
	var inp []interface{}
	var ffnc opt.Bobj
	var ndchk, par bool
	mx = -1
	
	//f = Cbmobj
	
	drw = "gen"
	trm = f.Term
	title = f.Title + "-g.a opt"
	if trm == "dxf"{trm = "svg"}
	f.Noprnt = true
	f.Term = ""
	f.Verbose = false
	f.Spam = false
	par = false
	var gpos []int
	var gb, gw float64
	var pltstr string
	
	np = 30; ngmn = 30; ngmx = 90
	mt = 5; ct = 1; st = 1; ndchk = true
	npmn = 20; npmx = 40
	pmut = 0.03; pcrs = 0.90; gap = 0.01; step = 10
	ng = 30
	
	
	maxval := f.Maxstep
	if maxval == 0{
		maxval = 22
	}
	ndims = f.Getndim()
	for i := 0; i < ndims; i++{
			nd = append(nd, maxval)
	}
	//fmt.Println(ColorRed,"NDIMS-",ndims,ColorReset)
	ncycle := NCycles
	var gmpos []int
	var gmw float64
	var plts []string
	inp = append(inp, f)
	ffnc = f2dobj
	
	switch f.Opt{
		case 1:
		//basic ga
		gpos, gw, gb, pltstr, err = opt.Gabloop(f.Web, par,ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, ffnc, drw, trm, title)
		//fmt.Println("gpos, gw",gpos, gw, gb, pltstr)
		case 11:
		//adaptive ga
		gpos, gw, gb, pltstr, err = opt.Gabaloop(f.Web, par, ndchk, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, pmut, pcrs, gap, nd, inp, ffnc, drw, trm, title)
		case 12:
		//n runs of basic ga
		ochn := make(chan []interface{}, ncycle)
		for i := 0; i < ncycle; i++{
			t1 := fmt.Sprintf("%s-run%v",title, i+1)
			go func(){
				gpos, gw, gb, pltstr, err = opt.Gabloop(f.Web, par,ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, ffnc, drw, trm, t1)
				rez := make([]interface{},5)
				rez[0] = gpos
				rez[1] = gw
				rez[2] = gb
				rez[3] = pltstr
				rez[4] = err
				ochn <- rez
			}()
		}
		for i := 0; i < ncycle; i++{
			rez := <- ochn
			gpos, _ = rez[0].([]int)
			gw, _ = rez[1].(float64)
			pltstr, _ = rez[3].(string)
			plts = append(plts, pltstr)
			if gmw == 0.0{
				gmw = gw
				gmpos = make([]int, len(gpos))
				copy(gmpos, gpos)
			} else {
				if gmw > gw{
					gmw = gw
					gmpos = make([]int, len(gpos))
					copy(gmpos, gpos)
				}
			}
		}
		gpos = gmpos
		gw = gmw
		case 13:
		//n runs of adaptive ga
		ochn := make(chan []interface{}, ncycle)
		for i := 0; i < ncycle; i++{
			t1 := fmt.Sprintf("%s-run%v",title, i+1)
			go func(){
				gpos, gw, gb, pltstr, err = opt.Gabaloop(f.Web, par, ndchk, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, pmut, pcrs, gap, nd, inp, ffnc, drw, trm, t1)
				rez := make([]interface{},5)
				rez[0] = gpos
				rez[1] = gw
				rez[2] = gb
				rez[3] = pltstr
				rez[4] = err
				ochn <- rez
			}()
		}
		for i := 0; i < ncycle; i++{
			rez := <- ochn
			gpos, _ = rez[0].([]int)
			gw, _ = rez[1].(float64)
			pltstr, _ = rez[3].(string)
			plts = append(plts, pltstr)
			if gmw == 0.0{
				gmw = gw
				gmpos = make([]int, len(gpos))
				copy(gmpos, gpos)
			} else {
				if gmw > gw{
					gmw = gw
					gmpos = make([]int, len(gpos))
					copy(gmpos, gpos)
				}
			}
		}
		gpos = gmpos
		gw = gmw
	
	}
	if err != nil{
		return
	}
	vec := make([]float64, len(gpos))
	for i, pos := range gpos{
		vec[i] = float64(pos)
	}
	
	frez = f2dpos(vec,inp)
	frez.Term = trm
	frez.Web = f.Web
	frez.Verbose = true
	frez.Noprnt = false
	switch f.Opt{
		case 1, 11:
		frez.Txtplots = append(frez.Txtplots, pltstr)
		default:
		frez.Txtplots = append(frez.Txtplots, plts...)
	}
	_, _, err = Frm2dDz(&frez)
	if err != nil{
		return 
	}
	frez.Report += fmt.Sprintf("optimized minimum cost (ga) for 2d frame -> %.f %.f rupees",gw, gw)
	return
}

//frm2dpso optimizes a frame via pso
func frm2dpso(f kass.Frm2d)(frez kass.Frm2d, err error){
	var w, c1, c2 float64
	var nd, ng, np, ncycle int
	var mx, mn []float64
	var drw, trm, title string
	var inp []interface{}
	var par, impr bool
	var gpos, gmpos []float64
	var gb, gmb float64
	var pltstr string
	var plts []string
	trm = f.Term
	f.Term = ""
	f.Noprnt = true
	f.Verbose = false
	f.Spam = false
	par = false
	inp = append(inp, f)
	nd = f.Getndim()
	mxstp := float64(f.Maxstep)
	if mxstp == 0.0{mxstp = 21.0}
	for i := 0; i < nd; i++{
		mx = append(mx, mxstp)
		mn = append(mn, 0.0)
	}
	//fmt.Println("calling PSO")
	w = 0.5; c1 = 1.0; c2 = 2.0
	ng = 30; np = 20
	drw = "gen"
	ncycle = NCycles
	if !f.Web{trm = "dumb"}
	title = f.Title + "-pso opt"
	switch f.Opt{
		case 21, 23:
		impr = true
		w = 0.5
		c1 = 2.0
		c2 = 2.0
		ng = 60		
	}
	ffnc := f2dfit
	switch f.Opt{
		case 2, 21:		
		gpos, gb, pltstr = opt.Psoloop(f.Web, par, impr, w, c1, c2, np, ng, nd, mx,mn, ffnc, inp, drw, trm, title)
		case 22,23:
		ochn := make(chan []interface{}, ncycle)
		for i := 0; i < ncycle; i++{
			t1 := fmt.Sprintf("%s-run%v",title,i+1)
			go func(){
				gpos, gb, pltstr = opt.Psoloop(f.Web,par, impr, w, c1, c2, np, ng, nd,mx,mn,ffnc, inp, drw, trm, t1)
				rez := make([]interface{},3)
				rez[0] = gpos
				rez[1] = gb
				rez[2] = pltstr
				ochn <- rez
			}()	
		}
		for i := 0; i < ncycle; i++{
			rez := <- ochn
			gpos, _ = rez[0].([]float64)
			gb, _ = rez[1].(float64)
			pltstr, _ = rez[2].(string)
			plts = append(plts, pltstr)
			if gmb == 0.0{
				gmb = gb
				gmpos = make([]float64, len(gpos))
				copy(gmpos, gpos)
			} else {
				if gmb > gb{
					gmb = gb
					gmpos = make([]float64, len(gpos))
					copy(gmpos, gpos)
				}
			}
		}
		gpos = gmpos
		gb = gmb
	}
	frez = f2dpos(gpos,inp)
	frez.Term = trm
	frez.Web = f.Web
	frez.Verbose = true
	frez.Noprnt = false
	switch f.Opt{
		case 2, 21:
		frez.Txtplots = append(frez.Txtplots, pltstr)
		default:
		frez.Txtplots = append(frez.Txtplots, plts...)
	}
	_, _, err = Frm2dDz(&frez)
	if err != nil{
		return 
	}
	frez.Report += fmt.Sprintf("optimized minimum cost (pso) for 2d frame -> %.f rupees",gb)
	return
}

//Frm2dOpt opts a 2d rcc frame (entry func)
func Frm2dOpt(f kass.Frm2d) (frez kass.Frm2d, err error){
	switch f.Opt{
		case 1, 11, 12, 13:
		//ga
		return frm2dga(f)
		case 2, 21, 22, 23:
		//pso
		return frm2dpso(f)
	}
	return 
}

//f2dobj is the objective function for evaluating an opt.Bat ([]int)
func f2dobj(b *opt.Bat, inp []interface{}) (err error) {
	vec := make([]float64, len(b.Pos))
	for i, val := range b.Pos{
		vec[i] = float64(val) 
	}
	fit := f2dfit(vec, inp)
	b.Wt = fit
	b.Fit = 1.0/fit
	return
}

//f2dfit is the objective function for evaluating a floating pos vec
func f2dfit(pos []float64, inp []interface{}) (fit float64){
	//fmt.Println("pos->",pos)

	//ye olde
	//f := getf2d(pos,inp)
	f := f2dpos(pos, inp) 
	f.Init()
	f.Noprnt = true
	allcons, _, err := Frm2dDz(&f)
	if err != nil || f.Kost <= 0.0{
		fit = 1e20
		return
	} else {
		fit = f.Kost * (1.0 + 10.0 * float64(allcons))
	}
	return
}

//f2dpos returns a 2d frame from a position vector
func f2dpos(pos []float64, inp []interface{})(f kass.Frm2d){
	f, _ = inp[0].(kass.Frm2d)
	f.Kost = 0.0
	err := f.SecGenRcc(pos)
	if err != nil{
		fmt.Println(ColorRed,"ERRORE generating frame sections",ColorReset)
		return
	}
	return
}


/*

	if 1 == 2{	
		var bw, bc, d, dc, db, dec, dmc, deb, dmb float64
		step := 25.0
		strt := 225.0
		switch f.Ngrp{
			case 0:
			switch f.Width{
				case 0.0:
				bw = math.Round(pos[0]) * step + strt
				d = math.Round(pos[1]) * step + strt
				default:
				bw = f.Width
				d = math.Round(pos[0]) * step + strt
			}
			bc = bw + f.Ceo * 2.0
			f.Sections = make([][]float64, 2)
			f.Sections[0] = []float64{bc, d}
			f.Sections[1] = []float64{bw, d}
			f.Csec = []int{1}
			f.Bsec = []int{2}
			case 1:
			switch f.Width{
				case 0.0:
				bw = math.Round(pos[0]) * step + strt
				dc = math.Round(pos[1]) * step + strt
				db = math.Round(pos[2]) * step + strt
				default:
				bw = f.Width
				dc = math.Round(pos[0]) * step + strt
				db = math.Round(pos[1]) * step + strt			
			}
			bc = bw + f.Ceo * 2.0
			f.Sections = make([][]float64, 2)
			f.Sections[0] = []float64{bc, dc}
			f.Sections[1] = []float64{bw, db}
			f.Csec = []int{1}
			f.Bsec = []int{2}
			case 2:
			switch f.Width{
				case 0.0:
			//TODO
				default:
				bw = f.Width
				switch f.Dconst{
					case true:
					dec = math.Round(pos[0]) * step + strt
					dmc = math.Round(pos[1]) * step + strt			
					deb = math.Round(pos[2]) * step + strt
					bc = f.Ceo * 2.0 + bw
					f.Sections = make([][]float64, 3)
					f.Sections[0] = []float64{bc, dec}
					f.Sections[1] = []float64{bc, dmc}
					f.Sections[2] = []float64{bw, deb}
					for i := range f.X{
						switch i{
							case 0:
							f.Csec = append(f.Csec, 1)
							f.Bsec = append(f.Bsec, 3)
							case len(f.X)-1:
							f.Csec = append(f.Csec, 1)
							default:
							f.Csec = append(f.Csec, 2)
							f.Bsec = append(f.Bsec, 3)
						}
					}
					case false:
					dec = math.Round(pos[0]) * step + strt
					dmc = math.Round(pos[1]) * step + strt			
					deb = math.Round(pos[2]) * step + strt
					dmb = math.Round(pos[3]) * step + strt
					bc = f.Ceo * 2.0 + bw
					f.Sections = make([][]float64, 4)
					f.Sections[0] = []float64{bc, dec}
					f.Sections[1] = []float64{bc, dmc}
					f.Sections[2] = []float64{bw, deb}
					f.Sections[3] = []float64{bw, dmb}
					for i := range f.X{
						switch i{
							case 0:
							f.Csec = append(f.Csec, 1)
							f.Bsec = append(f.Bsec, 3)
							case len(f.X)-1:
							f.Csec = append(f.Csec, 1)
							default:
							f.Csec = append(f.Csec, 2)
							f.Bsec = append(f.Bsec, 4)
						}
					}
				}
			}
			case 3:
			//individual columns and beams
			//TODO
		}
	}
*/
//getf2d returns a 2d frame from a floating point position vector
//old version
// func getf2d(pos []float64, inp []interface{})(f kass.Frm2d){
// 	//gen f2d
// 	f, _ = inp[0].(kass.Frm2d)
// 	f.Kost = 0.0
// 	step := 25.0
// 	f.Sections = make([][]float64, f.Ncols+f.Nbms+f.Ncls)
// 	f.Styps = make([]int, f.Ncols+f.Nbms+f.Ncls)
// 	f.Csec = make([]int, f.Ncols)
// 	f.Bsec = make([]int, f.Nbms + f.Ncls)
// 	f.Verbose = false
// 	f.Term = ""
// 	//fmt.Println("cstyp, bstyp->",f.Cstyp, f.Bstyp)
// 	switch f.Width{
// 		case 0.0:
// 		bw := math.Round(pos[0]/step) * step
// 		switch f.Dconst{
// 			case true:
// 			dbm := math.Round(pos[len(pos)-1]/step) * step
// 			for i, val := range pos[1:]{
// 				if i < f.Ncols{
// 					f.Sections[i] = make([]float64,2)
// 					f.Sections[i][0] = bw
// 					f.Sections[i][1] = math.Round(val/step) * step
// 					f.Styps[i] = f.Cstyp
// 					f.Csec[i] = i + 1
// 				}
// 			}
// 			for i := range f.Sections[f.Ncols:]{
// 				f.Sections[i+f.Ncols] = make([]float64, 2)
// 				f.Sections[i+f.Ncols][0] = bw
// 				f.Sections[i+f.Ncols][1] = dbm
// 				f.Styps[i+f.Ncols] = f.Bstyp
// 				f.Bsec[i] = i + f.Ncols + 1
// 			}
// 			case false:
// 			for i, val := range pos[1:]{
// 				f.Sections[i] = make([]float64,2)
// 				f.Sections[i][0] = bw
// 				f.Sections[i][1] = math.Round(val/step) * step 
// 				switch{
// 					case i < f.Ncols:
// 					f.Csec[i] = i + 1
// 					f.Styps[i] = f.Cstyp
// 					default:
// 					f.Bsec[i-f.Ncols] = i+1
// 					f.Styps[i] = f.Bstyp
// 				}
// 			}
// 		}
// 		default:
// 		//f.Verbose = true
// 		bw := f.Width
// 		switch f.Dconst{
// 			case true:
// 			dbm := math.Round(pos[len(pos)-1]/step) * step
// 			for i, val := range pos{
// 				if i < f.Ncols{
// 					f.Sections[i] = make([]float64,2)
// 					f.Sections[i][0] = bw
// 					f.Sections[i][1] = math.Round(val/step) * step 
// 					f.Csec[i] = i + 1
// 					f.Styps[i] = f.Cstyp
// 				}
// 			}
// 			for i := range f.Sections[f.Ncols:]{
// 				f.Sections[i+f.Ncols] = make([]float64, 2)
// 				f.Sections[i+f.Ncols][0] = bw
// 				f.Sections[i+f.Ncols][1] = dbm 
// 				f.Bsec[i] = i + f.Ncols + 1
// 				f.Styps[i+f.Ncols] = f.Bstyp
					
// 			}
// 			case false:
// 			for i, val := range pos{
// 				f.Sections[i] = make([]float64,2)
// 				f.Sections[i][0] = bw
// 				f.Sections[i][1] = math.Round(val/step) * step 
// 				switch{
// 					case i < f.Ncols:
// 					f.Csec[i] = i + 1
// 					f.Styps[i] = f.Cstyp
					
// 					default:
// 					f.Bsec[i-f.Ncols] = i+1
// 					f.Styps[i] = f.Bstyp
					
// 				}
// 			}
// 		}
// 	}
// 	//fmt.Println("building frame")
// 	//fmt.Println(f.Csec)
// 	//fmt.Println(f.Bsec)
// 	//fmt.Println(f.Sections)
// 	//fmt.Println(f.Styps)
// 	return
// }

/*
   
			//jl, ml := f.Mod.SumFrcs()
			//f.Mod.Msloads = ml
			//f.Mod.Jloads = jl
			//f.DrawMod(f.Term)
			//f.Mod.Msloads = f.Loadcons[1]
			//f.Mod.Jloads = f.Jloadmap[1]
			//jl, ml = f.Mod.SumFrcs()
			//f.Mod.Msloads = ml
			//f.Mod.Jloads = jl
			//f.DrawMod(f.Term)
			//f.Term = "mono"
			//Frm2dTable(&f, true)
			//fmt.Println("global best vector->",gb)
			//Frm2dTable(&f, true)
			//f.DrawMod("dumb")
			//fmt.Println(f.Txtplots[0])
			//for _, flrvec := range f.CBeams{
			//	f.Txtplots = append(f.Txtplots, PlotBmEnv(f.Bmenv, flrvec, "qt"))
			//}
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

   		var w, c1, c2 float64
		var nd, ng, np int
		var mx, mn []float64
		var drw, trm, title string
		var inp []interface{}
		var gpos []float64
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
		
		ng = 40; np = 20
		//ng = 3; np = 2
		drw = "gen"; trm = term; title = f.Title
		gpos, gb, pltstr = opt.Psoloop(f.Web, par, true, w, c1, c2, np, ng, nd, mx,mn, f2dfit, inp, drw, trm, title)
		// f.Term = term
		
		// f.Verbose = true
		// f.Noprnt = false
		f = getf2d(gpos,inp)
		
		_, _, err = Frm2dDz(&f)
		if err != nil{
			return
			
		} else{
			f.Verbose = true
			f.Noprnt = false
			f.Term = term
			if f.Web{
				Frm2dTable(&f, false)
				
			} else {
				Frm2dTable(&f, true)

			}
			frez = f
		}
		
		//fmt.Println(par, w, c1, c2, np, ng, nd, mx,mn,  inp, drw, trm, title,f)

*/
