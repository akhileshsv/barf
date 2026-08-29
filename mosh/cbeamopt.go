package barf

import (
	"fmt"
	"math"
	opt"barf/opt"
)

//CBmOpt calls ga and pso routines to optimize a continuous beam 
func CBmOpt(cb CBm)(copt CBm, err error){
	if cb.Nspans == 0{
		cb.Nspans = len(cb.Lspans)
	}
	if cb.Nspans == 0{
		err = fmt.Errorf("error in spans-> lspans %v\n nspans %v",cb.Lspans,cb.Nspans)
		return
	}
	switch cb.Opt{
		case 1, 11, 12, 13:
		//ga
		return cbmga(cb)
		case 2, 21, 22, 23:
		//pso
		return cbmpso(cb)
	}
	return
}

//cbmpso opts a cbm by pso
func cbmpso(cb CBm) (cbrez CBm, err error){
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
	trm = cb.Term
	cb.Term = ""
	cb.Noprnt = true
	inp = append(inp, cb)
	nd = cb.getndim()
	for i := 0; i < nd; i++{
		mx = append(mx, 21.0)
		mn = append(mn, 0.0)
	}
	//fmt.Println("calling PSO")
	w = 0.5; c1 = 1.0; c2 = 2.0
	ng = 30; np = 20
	ncycle = NCycles
	par = true
	drw = "gen"
	if !cb.Web{trm = "dumb"}
	title = cb.Title
	switch cb.Opt{
		case 21, 23:
		impr = true
		w = 0.5
		c1 = 2.0
		c2 = 2.0
		ng = 5		
	}
	ffnc := cbmfit
	switch cb.Opt{
		case 2, 21:		
		gpos, gb, pltstr = opt.Psoloop(cb.Web, par, impr, w, c1, c2, np, ng, nd, mx,mn, ffnc, inp, drw, trm, title)
		case 22,23:
		ochn := make(chan []interface{}, ncycle)
		for i := 0; i < ncycle; i++{
			t1 := fmt.Sprintf("%s-run%v",title,i+1)
			go func(){
				gpos, gb, pltstr = opt.Psoloop(cb.Web,par, impr, w, c1, c2, np, ng, nd,mx,mn,ffnc, inp, drw, trm, t1)
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
	cbrez = getcbm(gpos,inp)
	cbrez.Term = trm
	cbrez.Web = cb.Web
	switch cb.Opt{
		case 2, 21:
		cbrez.Txtplots = append(cbrez.Txtplots, pltstr)
		default:
		cbrez.Txtplots = append(cbrez.Txtplots, plts...)
	}
	//cbrez.Noprnt = true
	bmenv, e := CBeamEnvRcc(&cbrez, cbrez.Term, false)
	if e != nil{
		err = e
		return 
	}
	var allcons int
	allcons, err = CBmDz(&cbrez,bmenv)
	if err != nil{
		err = fmt.Errorf("failed constraints->%v\nerr->%s",allcons,err)
		return 
	}
	printz := false
	if !cb.Web && cb.Verbose{
		printz = true
	}
	pltstr = PlotCBmDet(cbrez.Web, cbrez.Bmvec, cbrez.RcBm, "", cbrez.Title, cbrez.Term)
	
	if cb.Web{
		switch cb.Term{
			case "dxf":
			pltstr = cb.Title + "-detail.dxf"
			case "svg","svgmono":
			pltstr = cb.Title + "-detail.svg"
		}
	}
	cbrez.Txtplots = append(cbrez.Txtplots, pltstr)
	cbrez.Table(printz)
	switch len(cb.Kostin){
		case 3:
		gb = gb * cb.Kostin[0]
		default:
		gb = gb * CostRcc
	}
	cbrez.Report += fmt.Sprintf("optimized minimum cost for cbeam -> %f rupees",gb)
	cbrez.Kost = gb
	cbrez.Vec = gpos
	return
}

// for _, span := range cb.RcBm{
// 	fmt.Println("i, span->",span[1].Title,span[1].Mid)
// 	//rezstring += fmt.Sprintf("span-%v\n",j+1)
// 	for _, bm := range span{
// 		if bm.Ignore{continue}
// 		//fmt.Println("is dz->",bm.Dz,"ast, asc->",bm.Ast, bm.Asc,"rbrt,rbrc->",bm.Rbrt,bm.Rbrc)
// 		//bm.BarPrint()
// 		bm.Table(false)
// 		//t.Log(bm.Report)
// 		//t.Log(bm.Report)
// 		//PlotBmGeom(bm, bm.Term)
// 		// if bm.Term == "dumb"{
// 		// 	fmt.Println(bm.Txtplot)
// 		// }
// 		//fmt.Println("dias->\n",bm.Dias,"\nbar points\n",bm.Barpts)
// 	}
// }
// fmt.Println("beamvec->",cb.Bmvec)
// fmt.Println("rcbm->",cb.RcBm)

//cbmga opts a cbm by ga
func cbmga(cb CBm)(cbrez CBm, err error){
	var drw, trm, title string
	var ndims, np, ng, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step int
	var pmut, pcrs, gap float64
	var nd []int
	var inp []interface{}
	var f opt.Bobj
	var ndchk, par bool
	mx = -1
	f = Cbmobj
	drw = "gen"
	trm = cb.Term
	title = cb.Title + "-g.a opt history"
	par = true
	cb.Noprnt = true
	cb.Verbose = false
	cb.Term = ""
	var gpos []int
	var gb, gw float64
	var pltstr string
	np = 40; ngmn = 40; ngmx = 75; mt = 1; ct = 1; st = 1; ndchk = true
	npmn = 30; npmx = 60
	pmut = 0.02; pcrs = 0.75; gap = 0.01; step = 10
	ndims = cb.getndim()
	maxval := 22
	for i := 0; i < ndims; i++{
			nd = append(nd, maxval)
	}
	ng = 75
	ncycle := NCycles
	var gmpos []int
	var gmw float64
	var plts []string
	inp = append(inp, cb)
	switch cb.Opt{
		case 1:
		//basic ga
		gpos, gw, gb, pltstr, err = opt.Gabloop(cb.Web, par,ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, f, drw, trm, title)
		case 11:
		//adaptive ga
		gpos, gw, gb, pltstr, err = opt.Gabaloop(cb.Web, par, ndchk, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, pmut, pcrs, gap, nd, inp, f, drw, trm, title)
		case 12:
		//n runs of basic ga
		ochn := make(chan []interface{}, ncycle)
		for i := 0; i < ncycle; i++{
			t1 := fmt.Sprintf("%s-run%v",title, i+1)
			go func(){
				gpos, gw, gb, pltstr, err = opt.Gabloop(cb.Web, par,ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, f, drw, trm, t1)
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
				gpos, gw, gb, pltstr, err = opt.Gabaloop(cb.Web, par, ndchk, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, pmut, pcrs, gap, nd, inp, f, drw, trm, t1)
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
	if !cb.Web{
		fmt.Println("opt finished- gpos, gw, gb, pltstr-\n",gpos, gw, gb, pltstr)
	}
	cbrez = gatcbm(gpos,inp)
	cbrez.Term = trm
	cbrez.Web = cb.Web
	
	switch cb.Opt{
		case 1, 11:
		cbrez.Txtplots = append(cbrez.Txtplots, pltstr)
		default:
		cbrez.Txtplots = append(cbrez.Txtplots, plts...)
	}
	//cbrez.Noprnt = true
	bmenv, e := CBeamEnvRcc(&cbrez, cbrez.Term, false)
	if e != nil{
		err = e
		return 
	}
	var allcons int
	allcons, err = CBmDz(&cbrez,bmenv)
	if err != nil{
		err = fmt.Errorf("failed constraints->%v\nerr->%s",allcons,err)
		return 
	}
	printz := false
	if !cb.Web && !cb.Verbose{
		printz = true
	}
	pltstr = PlotCBmDet(cbrez.Web, cbrez.Bmvec, cbrez.RcBm, "", cbrez.Title, cbrez.Term)
	
	if cb.Web{
		switch cb.Term{
			case "dxf":
			pltstr = cb.Title + "-detail.dxf"
			case "svg","svgmono":
			pltstr = cb.Title + "-detail.svg"
		}
	}
	cbrez.Txtplots = append(cbrez.Txtplots, pltstr)
	cbrez.Table(printz)
	switch len(cb.Kostin){
		case 3:
		gw = gw * cb.Kostin[0]
		default:
		gw = gw * CostRcc
	}
	cbrez.Report += fmt.Sprintf("optimized minimum cost for cbeam -> %.f rupees",gw)
	cbrez.Kost = gw
	return
}

//getndim returns ndims for cbm opt
func (cb *CBm) getndim()(nd int){
	//cb.Dconst = true
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
	cb.Noprnt = true
	//fmt.Println(ColorCyan,"oot->",cb.Sections,ColorReset)
	return

}

//getcbm returns a CBm given a floating point position vector
func getcbm(pos []float64, inp []interface{})(cb CBm){
	cb, _ = inp[0].(CBm)
	cb.Sections = make([][]float64,cb.Nspans)
	cb.Verbose = false
	//log.Println("title-",cb.Title)
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
	cb.Noprnt = true
	return
}

//cbmfit calcs a pos vec's(pso) fitness 
func cbmfit(pos []float64, inp []interface{}) (fit float64){
	cb := getcbm(pos,inp)
	cb.Noprnt = true
	bmenv, err := CBeamEnvRcc(&cb, cb.Term, false)
	if err != nil{
		fit = 1e6
		return
	}
	_, _ = CBmDz(&cb, bmenv)
	fit = cb.Kost
	if fit <= 0.0{
		fit = 1e6
		return
	}
	
	fit = fit * (10.0 * float64(cb.Allcons) + float64(cb.Csteel) + 1.0)
	// if cb.Allcons > 0{
	// 	fit = fit * math.Pow(10.0 * float64(cb.Allcons)+1.0,2)
	// }
	// if cb.Csteel > 0{
	// 	fit = fit * (float64(cb.Csteel)+1.0)
	// }
	return
}


//Cbmobj is an objective function for opt.Bat fitneess evaluation
func Cbmobj(b *opt.Bat, inp []interface{}) (err error){
	cb := gatcbm(b.Pos,inp)
	cb.Noprnt = true
	bmenv, err := CBeamEnvRcc(&cb, cb.Term, false)
	if err != nil{
		b.Wt = 1e7
		return
	}
	_, _ = CBmDz(&cb, bmenv)
	b.Wt = cb.Kost
	b.Fit = 1.0/b.Wt
	if b.Wt <= 0.0{
		b.Wt = 1e6
		b.Fit = 1.0/b.Wt
		return
	}
	b.Wt = b.Wt * (10.0 * float64(cb.Allcons) + float64(cb.Csteel) + 1.0)
	b.Fit = 1.0/b.Wt
	// if cb.Allcons > 0{
	// 	b.Wt = b.Wt * math.Pow(10.0 * float64(cb.Allcons)+1.0,2)
	// }
	// if cb.Csteel > 0{
	// 	b.Wt = b.Wt * (1.0+float64(cb.Csteel))
	// }
	return
}

//yeolde pso obj

// allcons, _ = CBmDz(&cb,bmenv)
// c1 := 0
// for _, barr := range cb.RcBm{
// 	fit += barr[1].Kost
// 	for _, bm := range barr{
// 		if bm.Csteel{
// 			c1 += 1
// 		}
// 	}
// }
// switch len(cb.Kostin){
// 	case 3:
// 	fit = fit/cb.Kostin[0]
// 	default:
// 	fit = fit/CostRcc
// }
// if fit == 0.0{
// 	fit = 1e6

// }
// if allcons > 0{
// 	fit = fit*(100.0 * float64(allcons)+1.0) 
// }
// if c1 > 0{
// 	fit = fit * (10.0 * float64(c1))
// }

//yeolde ga obj


// var allcons int
// bmenv, err := CBeamEnvRcc(&cb, cb.Term, false)
// if err != nil{
// 	b.Wt = 1e10
// 	b.Fit = 1.0/b.Wt
// 	return
// }
// allcons, _ = CBmDz(&cb,bmenv)
// c1 := 0
// for _, barr := range cb.RcBm{
// 	b.Wt += barr[1].Kost
// 	for _, bm := range barr{
// 		if bm.Csteel{
// 			c1 += 1
// 		}
// 	}
// }
// switch len(cb.Kostin){
// 	case 3:
// 	b.Wt = b.Wt/cb.Kostin[0]
// 	default:
// 	b.Wt = b.Wt/CostRcc
// }
// if b.Wt == 0.0{b.Wt = 1e6}

// if allcons > 0{
// 	b.Wt = b.Wt*(100.0 * float64(allcons)+1.0) 
// }
// if c1 > 0{
// 	b.Wt = b.Wt * (10.0 * float64(c1))
// }
