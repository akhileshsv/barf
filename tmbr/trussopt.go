package barf


// import (
// 	"fmt"
// 	// "math"
// 	// kass"barf/kass"
// 	// opt"barf/opt"
// )

// //TrussOpt optimizes a 2d truss struct
// func TrsOpt(trs kass.Trs2d)(trez kass.Trs2d, err error){
// 	switch trs.Opt{
// 		case 1, 11, 12, 13:
// 		//ga
// 		return trsga(trs)
// 		case 2, 21, 22, 23:
// 		//pso
// 		return trspso(trs)
// 	}
// 	return
// }

// //OptMod is the entry func for truss (axial force) model based optimization funcs
// func OptTrsMod(mod kass.Model)(mrez kass.Model, err error){
// 	switch mod.Frmstr{
// 		case "2dt","3dt":
// 		switch mod.Opt{
// 			case 1,11,12,13:
// 			//g.a, adapt. ga, nruns ga, nruns adapt ga 
// 			return trsoptga(mod)
// 			case 2,21,22,23:
// 			//pso, pso w improvement criteria, nruns pso, nruns pso w/impr
// 			return trsoptpso(mod)
// 		}
// 		default:
// 		err = fmt.Errorf("%s model opt is not written yet",mod.Frmstr)
// 	}
// 	return
// }

// //trsoptga optimizes a truss using rajeev-krishnamurthy g.a
// func trsoptga(mod kass.Model)(mrez kass.Model, err error){
// 	var drw, trm, title string
// 	var np, ng, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, nsecs int
// 	var pmut, pcrs, gap float64
// 	var nd []int
// 	var inp []interface{}
// 	var f opt.Bobj
// 	var ndchk, par bool
// 	modr := kass.Model{
// 		Id:mod.Id,
// 		Ncjt:mod.Ncjt,
// 		Frmstr:mod.Frmstr,
// 		Coords:mod.Coords,
// 		Supports:mod.Supports,
// 		Mprp:mod.Mprp,
// 		Jloads:mod.Jloads,
// 		Em:mod.Em,
// 		Fy:mod.Fy,
// 		Dmax:mod.Dmax,
// 		Pg:mod.Pg,
// 		Noprnt:true,
// 		Zup:mod.Zup,
// 	}
// 	inp = append(inp, modr)
// 	inp = append(inp, mod.Dims)
// 	mx = -1
// 	f = rakatrsobj
// 	drw = "gen"; trm = "dumb"; title = mod.Id
// 	if mod.Web{trm = mod.Term}
// 	var gpos []int
// 	var gb, gw float64
// 	var pltstr string
// 	np = 40; ngmn = 40; ngmx = 125; mt = 2; ct = 1; st = 1; ndchk = true
// 	npmn = 40; npmx = 60
// 	pmut = 0.02; pcrs = 0.75; gap = 0.01; step = 10
// 	switch len(mod.Dims){
// 		case 1:
// 		nsecs = len(mod.Dims[0])
// 	}
// 	for i := 0; i < mod.Ngrps; i++{
// 			nd = append(nd, nsecs)
// 	}
// 	np = npmx
// 	ng = 75
// 	ncycle := 4
// 	var gmpos []int
// 	var gmw float64
// 	var plts []string

// 	switch mod.Opt{
// 		case 1:
// 		//basic ga
// 		gpos, gw, gb, pltstr, err = opt.Gabloop(mod.Web, par,ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, f, drw, trm, title)
// 		case 11:
// 		//adaptive ga
// 		gpos, gw, gb, pltstr, err = opt.Gabaloop(mod.Web, par, ndchk, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, pmut, pcrs, gap, nd, inp, f, drw, trm, title)
// 		case 12:
// 		//n runs of basic ga
// 		ochn := make(chan []interface{}, ncycle)
// 		for i := 0; i < ncycle; i++{
// 			t1 := fmt.Sprintf("%s-run%v",title, i+1)
// 			go func(){
// 				gpos, gw, gb, pltstr, err = opt.Gabloop(mod.Web, par,ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, f, drw, trm, t1)
// 				rez := make([]interface{},5)
// 				rez[0] = gpos
// 				rez[1] = gw
// 				rez[2] = gb
// 				rez[3] = pltstr
// 				rez[4] = err
// 				ochn <- rez
// 			}()
// 		}
// 		for i := 0; i < ncycle; i++{
// 			rez := <- ochn
// 			gpos, _ = rez[0].([]int)
// 			gw, _ = rez[1].(float64)
// 			pltstr, _ = rez[3].(string)
// 			plts = append(plts, pltstr)
// 			if gmw == 0.0{
// 				gmw = gw
// 				gmpos = make([]int, len(gpos))
// 				copy(gmpos, gpos)
// 			} else {
// 				if gmw > gw{
// 					gmw = gw
// 					gmpos = make([]int, len(gpos))
// 					copy(gmpos, gpos)
// 				}
// 			}
// 		}
// 		gpos = gmpos
// 		gw = gmw
// 		case 13:
// 		//n runs of adaptive ga
// 		ochn := make(chan []interface{}, ncycle)
// 		for i := 0; i < ncycle; i++{
// 			t1 := fmt.Sprintf("%s-run%v",title, i+1)
// 			go func(){
// 				gpos, gw, gb, pltstr, err = opt.Gabaloop(mod.Web, par, ndchk, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, pmut, pcrs, gap, nd, inp, f, drw, trm, t1)
// 				rez := make([]interface{},5)
// 				rez[0] = gpos
// 				rez[1] = gw
// 				rez[2] = gb
// 				rez[3] = pltstr
// 				rez[4] = err
// 				ochn <- rez
// 			}()
// 		}
// 		for i := 0; i < ncycle; i++{
// 			rez := <- ochn
// 			gpos, _ = rez[0].([]int)
// 			gw, _ = rez[1].(float64)
// 			pltstr, _ = rez[3].(string)
// 			plts = append(plts, pltstr)
// 			if gmw == 0.0{
// 				gmw = gw
// 				gmpos = make([]int, len(gpos))
// 				copy(gmpos, gpos)
// 			} else {
// 				if gmw > gw{
// 					gmw = gw
// 					gmpos = make([]int, len(gpos))
// 					copy(gmpos, gpos)
// 				}
// 			}
// 		}
// 		gpos = gmpos
// 		gw = gmw
	
// 	}
	
// 	if err != nil{
// 		return
// 		//fmt.Println(err)
// 	}
// 	// fmt.Println("gpos, gb, err->", gpos, gb, err)
// 	if !mod.Web{
// 		fmt.Println(gpos, gw, gb, pltstr)
// 	}
// 	modr.Noprnt = false
// 	modr.Cp = make([][]float64, mod.Ngrps)
// 	for i, idx := range gpos{
// 		modr.Cp[i] = make([]float64,1)
// 		modr.Cp[i][0] = mod.Dims[0][idx]
// 	}
// 	modr.Web = mod.Web
// 	modr.Term = mod.Term
// 	modr.Units = mod.Units
// 	err = kass.CalcMod(&modr, modr.Frmstr, modr.Term, false)
// 	if err != nil{			
// 		return	
// 	}
// 	mrez = modr
// 	// fmt.Println("PLTSTR!->",pltstr)
// 	switch mod.Opt{
// 		case 1, 11:
// 		mrez.Txtplots = append(mrez.Txtplots, pltstr)
// 		case 12, 13:
// 		mrez.Txtplots = append(mrez.Txtplots, plts...)
// 	}
// 	mrez.Report += fmt.Sprintf("optimized minimum weight for truss -> %.2f (units-%s)",gw, mod.Units)
// 	return
// }

// //rakatrsobj is the objective function for truss mod ga opt
// func rakatrsobj(b *opt.Bat, inp []interface{}) (err error){
// 	mod, _ := inp[0].(kass.Model)
// 	secs, _ := inp[1].([][]float64)
// 	// pmax,_ := inp[2].(float64)
// 	// dmax,_ := inp[3].(float64)
// 	// dens,_ := inp[4].(float64)
// 	pmax := mod.Fy
// 	dmax := mod.Dmax
// 	dens := mod.Pg
// 	cp := make([][]float64, len(b.Pos))
// 	for i, idx := range b.Pos{
// 		cp[i] = make([]float64,1)
// 		if len(secs) == 1{
// 			cp[i][0] = secs[0][idx]
// 		} else {
// 			//WRENG
// 			cp[i][0] = secs[i][idx]
// 		}
// 	}
// 	mod.Cp = cp
// 	var wt, C, gx, con float64
// 	frmrez, err := kass.CalcTrs(&mod, mod.Ncjt)
// 	if err != nil {
// 		b.Wt = 1e10
// 		b.Fit = 1.0/b.Wt
// 		return
// 	}
// 	js, _ := frmrez[0].(map[int]*kass.Node)
// 	ms, _ := frmrez[1].(map[int]*kass.Mem)
// 	for _, node := range js{
// 		for _, d := range node.Displ {
// 			gx = math.Abs(d)/dmax - 1.0
// 			if gx > 0.0 {
// 				C += gx
// 				con += 1.0
// 			}
// 		}
// 	}
// 	for _, mem := range ms{
// 		wt += mem.Geoms[0] * mem.Geoms[2] * dens
// 		pmem := mem.Qf[0] / mem.Geoms[2]
// 		gx = math.Abs(pmem)/pmax - 1.0
// 		if gx > 0.0 {
// 			C += gx
// 			con += 1.0
// 		}
// 	}
// 	b.Wt = wt*(1.0 + 10.0*C)
// 	b.Fit = 1.0/b.Wt
// 	b.Con = con
// 	err = nil
// 	return
// }

// //rakatrsfit is the objective function for trs mod pso
// func rakatrsfit(pos []float64, inp []interface{}) (fit float64){
// 	mod, _ := inp[0].(kass.Model)
// 	secs, _ := inp[1].([][]float64)
// 	pmax := mod.Fy
// 	dmax := mod.Dmax
// 	dens := mod.Pg
// 	cp := make([][]float64, len(pos))
// 	for i, val := range pos{
// 		cp[i] = make([]float64,1)
// 		idx := int(val)
// 		if len(secs) == 1{
// 			cp[i][0] = secs[0][idx]
// 		} else {
// 			cp[i][0] = secs[i][idx]
// 		}
// 	}
// 	mod.Cp = cp
// 	mod.Noprnt = true
// 	mod.Term = ""
	
// 	var wt, C, gx float64
// 	frmrez, err := kass.CalcTrs(&mod, mod.Ncjt)
// 	if err != nil {
// 		wt = 1e10
// 		fit = wt
// 		return
// 	}
// 	js, _ := frmrez[0].(map[int]*kass.Node)
// 	ms, _ := frmrez[1].(map[int]*kass.Mem)
// 	for _, node := range js{
// 		for _, d := range node.Displ {
// 			gx = math.Abs(d)/dmax - 1.0
// 			if gx > 0.0 {
// 				C += gx
// 				//con += 1.0
// 			}
// 		}
// 	}
// 	for _, mem := range ms{
// 		wt += mem.Geoms[0] * mem.Geoms[2] * dens
// 		pmem := mem.Qf[0] / mem.Geoms[2]
// 		gx = math.Abs(pmem)/pmax - 1.0
// 		if gx > 0.0 {
// 			C += gx
// 			//con += 1.0
// 		}
// 	}
// 	wt = wt*(1.0 + 10.0*C)
// 	fit = wt
// 	return
// }

// //trsoptpso optimizes a truss by pso
// func trsoptpso(mod kass.Model)(mrez kass.Model, err error){
// 	var w, c1, c2 float64
// 	var nd, ng, np int
// 	var mx, mn []float64
// 	var drw, trm, title string
// 	var inp []interface{}
// 	var par, impr bool
// 	nsecs := len(mod.Dims[0])
// 	mxval := float64(nsecs - 1)
// 	mnval := 0.0
// 	nd = mod.Ngrps
// 	if nd == 0{
// 		err = fmt.Errorf("invalid ngroups for truss opt -> %v",mod.Ngrps)
// 		return
// 	}
// 	for i := 0; i < mod.Ngrps; i++{
// 		mx = append(mx, mxval)
// 		mn = append(mn, mnval)
// 	}
// 	par = false
// 	drw = "gen"
// 	trm = mod.Term
// 	title = mod.Id
// 	w = 0.4; c1 = 1.0; c2 = 1.0
// 	ng = 50; np = 25
// 	modr := kass.Model{
// 		Id:mod.Id,
// 		Ncjt:mod.Ncjt,
// 		Frmstr:mod.Frmstr,
// 		Coords:mod.Coords,
// 		Supports:mod.Supports,
// 		Mprp:mod.Mprp,
// 		Jloads:mod.Jloads,
// 		Em:mod.Em,
// 		Fy:mod.Fy,
// 		Dmax:mod.Dmax,
// 		Pg:mod.Pg,
// 		Noprnt:true,
// 		Zup:mod.Zup,
// 		Ngrps:mod.Ngrps,
// 		Dims:mod.Dims,
// 	}
// 	inp = append(inp, modr)
// 	inp = append(inp, mod.Dims)
// 	//fmt.Println("nd, mx, mn->",nd, mx, mn)
// 	switch mod.Opt{
// 		case 21, 23:
// 		//use li values
// 		impr = true
// 		w = 0.5
// 		c1 = 2.0
// 		c2 = 2.0
// 		ng = 75
// 	}
// 	var gpos, gmpos []float64
// 	var gb, gmb float64
// 	var pltstr string
// 	var plts []string
// 	ncycle := 4
// 	switch mod.Opt{
// 		case 2, 21:
// 		//simple pso, pso w/imprv
// 		gpos, gb, pltstr = opt.Psoloop(mod.Web,par, impr, w, c1, c2, np, ng, nd,mx,mn,rakatrsfit, inp, drw, trm, title)
// 		case 22, 23:
// 		//ncycle runs of simple pso, pso w/imprv
// 		ochn := make(chan []interface{}, ncycle)
// 		for i := 0; i < ncycle; i++{
// 			t1 := fmt.Sprintf("%s-run%v",title,i+1)
// 			go func(){
// 				gpos, gb, pltstr = opt.Psoloop(mod.Web,par, impr, w, c1, c2, np, ng, nd,mx,mn,rakatrsfit, inp, drw, trm, t1)
// 				rez := make([]interface{},3)
// 				rez[0] = gpos
// 				rez[1] = gb
// 				rez[2] = pltstr
// 				ochn <- rez
// 			}()
			
// 		}
// 		for i := 0; i < ncycle; i++{
// 			rez := <- ochn
// 			gpos, _ = rez[0].([]float64)
// 			gb, _ = rez[1].(float64)
// 			pltstr, _ = rez[2].(string)
// 			plts = append(plts, pltstr)
// 			if gmb == 0.0{
// 				gmb = gb
// 				gmpos = make([]float64, len(gpos))
// 				copy(gmpos, gpos)
// 			} else {
// 				if gmb > gb{
// 					gmb = gb
// 					gmpos = make([]float64, len(gpos))
// 					copy(gmpos, gpos)
// 				}
// 			}
// 		}
// 		gpos = gmpos
// 		gb = gmb
		
// 	}
// 	modr.Cp = make([][]float64, mod.Ngrps)
	
// 	for i, val := range gpos{
// 		modr.Cp[i] = make([]float64,1)
// 		idx := int(math.Round(val))
// 		if len(mod.Dims) == 1{
// 			modr.Cp[i][0] = mod.Dims[0][idx]
// 		} else {
// 			modr.Cp[i][0] = mod.Dims[i][idx]
// 		}
// 	}
// 	modr.Web = mod.Web
// 	modr.Term = mod.Term
// 	modr.Units = mod.Units
// 	modr.Noprnt = false
// 	err = kass.CalcMod(&modr, modr.Frmstr, modr.Term, false)
// 	if err != nil{			
// 		return	
// 	}
// 	mrez = modr
// 	switch mod.Opt{
// 		case 2, 21:
// 		mrez.Txtplots = append(mrez.Txtplots, pltstr)
// 		case 22, 23:
// 		mrez.Txtplots = append(mrez.Txtplots, plts...)
// 	}
// 	mrez.Report += fmt.Sprintf("optimized minimum weight for truss -> %.2f (units-%s)",gb, mod.Units)
// 	// for _, txtplot := range mrez.Txtplots{
// 	// 	fmt.Println(txtplot)
// 	// }
// 	//fmt.Println(mrez.Report)
// 	return
// }

// //trstopopso performs topology opt of a truss by pso
// func trstopopso(mod kass.Model)(mrez kass.Model, err error){
// 	var w, c1, c2 float64
// 	var nd, ng, np int
// 	var mx, mn []float64
// 	var drw, trm, title string
// 	var inp []interface{}
// 	var par, impr bool
	
// 	nsecs := len(mod.Dims[0])
// 	mxval := float64(nsecs - 1)
// 	mnval := 0.0
// 	for i := 0; i < mod.Ngrps; i++{
// 		mx = append(mx, mxval)
// 		mn = append(mn, mnval)
// 	}
// 	mxval = 50.0
// 	mnval = -50.0
// 	nd = len(mod.Coords) 
// 	for i := 0; i < nd; i++{
// 		mx = append(mx, mxval)
// 		mn = append(mn, mnval)
// 	}
// 	par = false
// 	drw = "gen"
// 	trm = mod.Term
// 	title = mod.Id
// 	w = 0.4; c1 = 1.0; c2 = 1.0
// 	ng = 50; np = 25
// 	modr := kass.Model{
// 		Id:mod.Id,
// 		Ncjt:mod.Ncjt,
// 		Frmstr:mod.Frmstr,
// 		Coords:mod.Coords,
// 		Supports:mod.Supports,
// 		Mprp:mod.Mprp,
// 		Jloads:mod.Jloads,
// 		Em:mod.Em,
// 		Fy:mod.Fy,
// 		Dmax:mod.Dmax,
// 		Pg:mod.Pg,
// 		Noprnt:true,
// 		Zup:mod.Zup,
// 	}
// 	inp = append(inp, modr)
// 	//fmt.Println("nd, mx, mn->",nd, mx, mn)
// 	switch mod.Opt{
// 		case 21, 23:
// 		//use li values
// 		impr = true
// 		w = 0.5
// 		c1 = 2.0
// 		c2 = 2.0
// 		ng = 75
// 	}
// 	var gpos, gmpos []float64
// 	var gb, gmb float64
// 	var pltstr string
// 	var plts []string
// 	ncycle := 4
// 	switch mod.Opt{
// 		case 2, 21:
// 		//simple pso, pso w/imprv
// 		gpos, gb, pltstr = opt.Psoloop(mod.Web,par, impr, w, c1, c2, np, ng, nd,mx,mn,trstopofit, inp, drw, trm, title)
// 		case 22, 23:
// 		//ncycle runs of simple pso, pso w/imprv
// 		ochn := make(chan []interface{}, ncycle)
// 		for i := 0; i < ncycle; i++{
// 			t1 := fmt.Sprintf("%s-run%v",title,i+1)
// 			go func(){
// 				gpos, gb, pltstr = opt.Psoloop(mod.Web,par, impr, w, c1, c2, np, ng, nd,mx,mn,trstopofit, inp, drw, trm, t1)
// 				rez := make([]interface{},3)
// 				rez[0] = gpos
// 				rez[1] = gb
// 				rez[2] = pltstr
// 				ochn <- rez
// 			}()
			
// 		}
// 		for i := 0; i < ncycle; i++{
// 			rez := <- ochn
// 			gpos, _ = rez[0].([]float64)
// 			gb, _ = rez[1].(float64)
// 			pltstr, _ = rez[2].(string)
// 			plts = append(plts, pltstr)
// 			if gmb == 0.0{
// 				gmb = gb
// 				gmpos = make([]float64, len(gpos))
// 				copy(gmpos, gpos)
// 			} else {
// 				if gmb > gb{
// 					gmb = gb
// 					gmpos = make([]float64, len(gpos))
// 					copy(gmpos, gpos)
// 				}
// 			}
// 		}
// 		gpos = gmpos
// 		gb = gmb
		
// 	}
// 	modr.Cp = make([][]float64, mod.Ngrps)
	
// 	for i := 0; i < mod.Ngrps; i++{
// 		modr.Cp[i] = make([]float64,1)
// 		idx := int(math.Round(gpos[i]))
// 		if len(mod.Dims) == 1{
// 			modr.Cp[i][0] = mod.Dims[0][idx]
// 		} else {
// 			modr.Cp[i][0] = mod.Dims[i][idx]
// 		}
// 	}
// 	jmap := make(map[int]bool)
// 	for _, val := range mod.Jloads{
// 		idx := int(val[0])
// 		jmap[idx] = true
// 	}
// 	coords := make([][]float64, len(mod.Coords))
// 	for i, pt := range mod.Coords{
// 		coords[i] = make([]float64, 2)
// 		if _, ok := jmap[i+1]; !ok{
// 			coords[i][0] = pt[0]
// 			coords[i][1] = pt[1] + gpos[i+mod.Ngrps]
// 		} else {
// 			copy(coords[i],pt)
// 		}
// 	}
// 	modr.Coords = coords
	
// 	modr.Web = mod.Web
// 	modr.Term = mod.Term
// 	modr.Units = mod.Units
// 	modr.Noprnt = false
// 	err = kass.CalcMod(&modr, modr.Frmstr, modr.Term, false)
// 	if err != nil{			
// 		return	
// 	}
// 	mrez = modr
// 	switch mod.Opt{
// 		case 2, 21:
// 		mrez.Txtplots = append(mrez.Txtplots, pltstr)
// 		case 22, 23:
// 		mrez.Txtplots = append(mrez.Txtplots, plts...)
// 	}
// 	mrez.Report += fmt.Sprintf("optimized minimum weight for truss -> %.2f (units-%s)",gb, mod.Units)
// 	// for _, txtplot := range mrez.Txtplots{
// 	// 	fmt.Println(txtplot)
// 	// }
// 	//fmt.Println(mrez.Report)
	
// 	return
// }


// //trstopofit is the objective function for trs topo pso
// func trstopofit(pos []float64, inp []interface{}) (fit float64){
// 	mod, _ := inp[0].(kass.Model)
// 	secs := mod.Dims
// 	pmax := mod.Fy
// 	dmax := mod.Dmax
// 	dens := mod.Pg
	
// 	cp := make([][]float64, mod.Ngrps)
// 	for i := 0; i < mod.Ngrps; i++{
// 		cp[i] = make([]float64,1)
// 		idx := int(math.Round(pos[i]))
// 		if len(secs) == 1{
// 			cp[i][0] = secs[0][idx]
// 		} else {
// 			cp[i][0] = secs[i][idx]
// 		}
// 	}
// 	mod.Cp = cp
// 	jmap := make(map[int]bool)
// 	for _, val := range mod.Jloads{
// 		idx := int(val[0])
// 		jmap[idx] = true
// 	}
// 	coords := make([][]float64, len(mod.Coords))
// 	for i, pt := range mod.Coords{
// 		coords[i] = make([]float64, 2)
// 		if _, ok := jmap[i+1]; !ok{
// 			coords[i][0] = pt[0]
// 			coords[i][1] = pt[1] + pos[i+mod.Ngrps]
// 		} else {
// 			copy(coords[i],pt)
// 		}
// 	}
// 	mod.Coords = coords
// 	mod.Noprnt = true
// 	mod.Term = ""
// 	var wt, C, gx float64
// 	frmrez, err := kass.CalcTrs(&mod, mod.Ncjt)
// 	if err != nil {
// 		wt = 1e10
// 		fit = wt
// 		return
// 	}
// 	js, _ := frmrez[0].(map[int]*kass.Node)
// 	ms, _ := frmrez[1].(map[int]*kass.Mem)
// 	for _, node := range js{
// 		for _, d := range node.Displ {
// 			gx = math.Abs(d)/dmax - 1.0
// 			if gx > 0.0 {
// 				C += gx
// 				//con += 1.0
// 			}
// 		}
// 	}
// 	for _, mem := range ms{
// 		wt += mem.Geoms[0] * mem.Geoms[2] * dens
// 		pmem := mem.Qf[0] / mem.Geoms[2]
// 		gx = math.Abs(pmem)/pmax - 1.0
// 		if gx > 0.0 {
// 			C += gx
// 			//con += 1.0
// 		}
// 	}
// 	wt = wt*(1.0 + 10.0*C)
// 	fit = wt
// 	return
// }
