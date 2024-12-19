package barf

import (
	"log"
	"fmt"
	"math"
	opt"barf/opt"
)

func OptSubFrm(sf SubFrm) (sfrez SubFrm, err error){
	sf.Verbose = false
	sf.Noprnt = true
	sf.Sections = [][]float64{}
	sf.Kost = 0.0
	switch sf.Opt{
		case 1, 11, 12, 13:
		//ga
		sfrez, err = subfrmga(sf)
		case 2, 21, 22, 23:
		//pso
		sfrez, err = subfrmpso(sf)
	}
	return
}

func subfrmga(sf SubFrm)(sfrez SubFrm, err error){
	var drw, trm, title string
	var np, ng, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx int
	var pmut, pcrs, gap float64
	var nd []int
	var inp []interface{}
	var f opt.Bobj
	var ndchk, par bool

	inp = append(inp, sf)
	mx = -1
	f = sfobj
	drw = "gen";  title = sf.Title
	trm = sf.Term
	sf.Term = ""
	sf.Noprnt = true
	sf.Verbose = false
	var gpos []int
	var gb, gw float64
	var pltstr string
	np = 30; ngmn = 40; ngmx = 50; mt = 1; ct = 1; st = 1; ndchk = true
	npmn = 40; npmx = 60
	pmut = 0.02; pcrs = 0.75; gap = 0.01
	mxdx := (750-225)/25 
	ndim := sf.getndim()
	for i := 0; i < ndim; i++{
		nd = append(nd, mxdx)
	}
	ng = 70
	ncycle := 4
	var gmpos []int
	var gmw float64
	var plts []string
	step := 10
	par = true
	switch sf.Opt{
		case 1:
		gpos, gw, gb, pltstr, err = opt.Gabloop(sf.Web, par,ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, f, drw, trm, title)
		case 11:
		gpos, gw, gb, pltstr, err = opt.Gabaloop(sf.Web, par, ndchk, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step, pmut, pcrs, gap, nd, inp, f, drw, trm, title)
		case 12:
		ochn := make(chan []interface{}, ncycle)
		for i := 0; i < ncycle; i++{
			t1 := fmt.Sprintf("%s-run%v",title, i+1)
			go func(){
				gpos, gw, gb, pltstr, err = opt.Gabloop(sf.Web, par,ndchk, np, ng, mt, ct, cn, st, mx, pmut, pcrs, nd, inp, f, drw, trm, t1)
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
		//ncycle runs of adaptive ga
	}
	if err != nil{return}
	//getsf(pos []float64, inp []interface{})(sf SubFrm, err error){
	fpos := make([]float64, len(gpos))
	mindim := 225.0
	for i := range gpos{
		fpos[i] = float64(gpos[i]) * 25.0  + mindim
	}

	sfrez, err = getsf(fpos, inp)
	if err != nil{return}
	
	switch sf.Opt{
		case 1, 11:
		sfrez.Txtplots = append(sfrez.Txtplots, pltstr)
		default:
		sfrez.Txtplots = append(sfrez.Txtplots, plts...)
	}
	_, err = DzSubFrm(&sfrez)
	if err != nil{
		return
	}
	pltstr = PlotSfDet(&sfrez)
	sfrez.Txtplots = append(sfrez.Txtplots, pltstr)
	return
}

func subfrmpso(sf SubFrm)(sfrez SubFrm, err error){
	var w, c1, c2 float64
	var nd, ng, np int
	var mx, mn []float64
	var drw, trm, title string
	var inp []interface{}
	var par, impr bool
	sf.Nosec = true
	err = sf.Init()
	if err != nil{
		return
	}
	nd = sf.getndim()
	for i := 0; i < nd; i++{
		mx = append(mx, 21.0)
		mn = append(mn, 0.0)
	}
	//set pso params
	w = 0.4; c1 = 1.0; c2 = 1.0
	ng = 50; np = 30
	switch sf.Opt{
		case 2, 22:
		//simple pso
		case 21, 23:
		//improv pso
		impr = true
		w = 0.5
		c1 = 2.0
		c2 = 2.0
		ng = 70
	}
	var gpos, gmpos []float64
	var gb, gmb float64
	var pltstr string
	var plts []string
	ncycle := 4
	switch sf.Term{
		case "dumb","svg","svgmono", "dxf", "qt":
		default:
		sf.Term = "svg"
	}
	trm = sf.Term
	sf.Term = ""
	sf.Noprnt = true
	sf.Verbose = false
	par = true
	inp = append(inp, sf)
	switch sf.Opt{
		case 2, 21:
		//simple pso
		gpos, gb, pltstr = opt.Psoloop(sf.Web,par, impr, w, c1, c2, np, ng, nd,mx,mn,sffit, inp, drw, trm, title)
		
		
		case 22,23:
		//4 runs of pso w/improv
		//ncycle runs of simple pso, pso w/imprv
		ochn := make(chan []interface{}, ncycle)
		for i := 0; i < ncycle; i++{
			t1 := fmt.Sprintf("%s-run%v",title,i+1)
			go func(){
				gpos, gb, pltstr = opt.Psoloop(sf.Web,par, impr, w, c1, c2, np, ng, nd,mx,mn,sffit, inp, drw, trm, t1)
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
	sfrez, err = getsf(gpos, inp)
	if err != nil{return}
	sfrez.Term = trm
	sfrez.Noprnt = false
	sfrez.Verbose = true
	sfrez.Web = sf.Web
	switch sf.Opt{
		case 2, 21:
		sfrez.Txtplots = append(sfrez.Txtplots, pltstr)
		default:
		sfrez.Txtplots = append(sfrez.Txtplots, plts...)
	}
	_, err = DzSubFrm(&sfrez)
	if err != nil{return}
	pltstr = PlotSfDet(&sfrez)
	sfrez.Txtplots = append(sfrez.Txtplots, pltstr)
	switch len(sf.Kostin){
		case 3:
		sfrez.Kost = sfrez.Kost * sf.Kostin[0]
		default:
		sfrez.Kost = sfrez.Kost * CostRcc
	}
	sfrez.Report += fmt.Sprintf("optimized kost for sub frame- %f",sfrez.Kost)
	// fmt.Println(sfrez.Report)
	return
}

//sfobj is the objective function for evaluating an opt.Bat ([]int) as a slice 
func sfobj(b *opt.Bat, inp []interface{}) (err error) {
	pos := make([]float64, len(b.Pos))
	step := 25.0
	mindim := 225.0
	for i := range b.Pos{
		pos[i] = float64(b.Pos[i]) * step + mindim
	}
	b.Fit = sffit(pos, inp)
	return
}

//sffit is the objective function for evaluating an opt.Bat ([]int) as a slice of sl
func sffit(pos []float64, inp []interface{}) (fit float64){

	sf, err := getsf(pos,inp)
	if err != nil{
		log.Fatal(err)
	}
	var allcons int
	allcons, err = DzSubFrm(&sf)
	if err != nil || sf.Kost <= 0.0 || allcons > 0{
		fit = 1e10
		return
	} else {
		fit = sf.Kost
	}
	return
}


//getndim returns ndims/len of opt slice based on sf.Ngrp
func (sf *SubFrm) getndim() (nd int){
	//for now tis always dconst
	switch sf.Ngrp{
		case 0:
		switch sf.Width{
			case 0.0:
			//b, d
			nd = 2
			default:
			//d
			nd = 1
		}
		case 1:
		switch sf.Width{
			case 0.0:
			//b, dcol, dbm
			nd = 3
			default:
			//dcol, dbm
			nd = 2
		}
		case 2:
		//two cols, two beams
		switch sf.Width{
			case 0.0:
			//b, dc1, dc2, dbm
			nd = 4
			default:
			//dc1, dc2, dbm
			nd = 3
		}
		case 3:
		//nspan+1 cols, nspan beams
		switch sf.Width{
			case 0.0:
			nd = 1 + sf.Nspans + 1 + 1
			default:
			nd = sf.Nspans + 1 + 1
		}
	}
	return
}

//getrectsecs updates a subframe's sections/csec/bsec based on pos/sf.Ngrp
//only rect sects as it says so - maybe add more vars/funcs for other sections
func (sf *SubFrm) getrectsecs(step float64, pos []float64){
	var bw float64
	minval := 225.0
	// log.Println("pos in -> ", pos)
	fpos := make([]float64, len(pos))
	for i := range pos{
		fpos[i] = math.Round(pos[i]) * step  + minval
	}
	
	switch sf.Ngrp{
		case 0:
		var d float64
		//one col+beam
		switch sf.Width{
			case 0.0:
			// if len(pos) < 2{err = fmt.Errorf("invalid length of pos vec - %v for ngrp - %v",pos, sf.Ngrp); return}
			// bw = math.Round(pos[0]/step) * step
			// d = math.Round(pos[1]/step) * step
			bw = fpos[0]
			d = fpos[1]
			default:
			if len(pos) < 1{
				err := fmt.Errorf("invalid length of pos vec - %v for ngrp - %v",pos, sf.Ngrp)
				fmt.Println(err)
			}
			bw = sf.Width
			d = fpos[0]
			// d = math.Round(pos[0]/step) * step
		}
		sf.Sections = [][]float64{{bw, d}}
		case 1:
		var dc, db float64
		switch sf.Width{
			case 0.0:
			// bw = math.Round(pos[0]/step) * step
			// dc = math.Round(pos[1]/step) * step
			// db = math.Round(pos[2]/step) * step
			bw = fpos[0]
			dc = fpos[1]
			db = fpos[2]
			default:
			bw = sf.Width
			dc = fpos[0]
			db = fpos[1]
			// dc = math.Round(pos[0]/step) * step
			// db = math.Round(pos[1]/step) * step
		}
		sf.Sections = [][]float64{
			{bw, dc},
			{bw, db},
		}
		case 2:
		case 3:
	}
	
	// log.Println("pos out -> ", sf.Sections)
}


//getsf returns a subframe from a float slice
func getsf(pos []float64, inp []interface{})(sf SubFrm, err error){

	if len(pos) == 0{
		err = fmt.Errorf("invalid position vector - %v",pos)
		return
	}
	var ok bool
	sf, ok = inp[0].(SubFrm)
	if !ok{
		err = fmt.Errorf("missing subfrm in input slice inp - %v",inp)
		return
	}
	
	sf.Nosec = false
	//nd := sf.getndim()
	sf.getrectsecs(25.0, pos)
	return
}
