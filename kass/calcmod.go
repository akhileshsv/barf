package barf

import (
	//"os"
	"fmt"
	"log"
	"math"
)

//CalcModNp - entry func from flags/menu for non prismatic beam-column frame analysis
//see hulse section 2.7/mosley spencer section 4.1
//IS UNDONE (not werk)
//basically what should it do? print stuff? call freecad, etc
func CalcModNp(mod *Model, frmtyp, term string, pipe bool) (err error){
	frmrez, err := CalcNp(mod, frmtyp, true)
	if err != nil{
		return
	}
	js, _ := frmrez[0].(map[int]*Node)
	ms, _ := frmrez[1].(map[int]*MemNp)
	mod.Js = js
	mod.Mnps = ms
	pltchn := make(chan string)
	//mod.Calc = true
	switch mod.Frmstr{
		case "1db":
		go PlotNpBm1d(mod, mod.Term, pltchn)
		case "2df":
		go PlotNpFrm2d(mod, mod.Term, pltchn)
	}

	if mod.Calc{
		mod.CalcRezNp(frmtyp)
	}
	pltstr := <- pltchn
	mod.Txtplots = append(mod.Txtplots, pltstr)
	reportz,_ := frmrez[6].(string)
	reportz += FrcTable(mod)
	if !pipe && !mod.Web{fmt.Println(reportz)}
	mod.Report = reportz
	if term != "" && mod.Calc{
		switch mod.Frmstr{
			case "1db", "2df":
			mod.Txtplots = append(mod.Txtplots, PlotNpRez2d(mod, term))
		}
	}
	return
}

//CalcMod - entry func from flags/menu for model direct stiffness analysis
//called after json - model read
//plots line diagram and prints report
func CalcMod(mod *Model, frmtyp, term string, pipe bool) (err error){
	pltchn := make(chan string, 1)
	var frmrez []interface{}
	okdraw := true
	if term == ""{okdraw = false}
	//JUST DO THIS HERE
	mod.Frmstr = frmtyp
	//else color?
	switch frmtyp{
		case "1db":
		mod.Ncjt = 2
		mod.Frmtyp = 1
		if okdraw {go PlotBm1d(mod, term, pltchn)}
		frmrez, err = CalcBm1d(mod, 2)
		case "2dt":
		mod.Ncjt = 2
		mod.Frmtyp = 2
		if okdraw {go PlotTrs2d(mod, term, pltchn)}
		frmrez, err = CalcTrs(mod, 2)
		case "2df":
		mod.Ncjt = 3
		mod.Frmtyp = 3
		if okdraw {go PlotFrm2d(mod, term, pltchn)}
		frmrez, err = CalcFrm2d(mod, 3)
		case "3dt":
		mod.Ncjt = 3
		mod.Frmtyp = 4
		if okdraw {go PlotTrs3d(mod, term, pltchn)}
		frmrez, err = CalcTrs(mod, 3)
		case "3dg":
		mod.Ncjt = 3
		mod.Frmtyp = 5
		if okdraw {go PlotGrd3d(mod, term, pltchn)}
		frmrez, err = CalcGrd(mod, 3)
		case "3df":
		mod.Ncjt = 6
		mod.Frmtyp = 6
		if okdraw {go PlotFrm3d(mod, term, pltchn)}
		frmrez, err = CalcFrm3d(mod, 6)
		default:
		log.Println(ColorRed,"invalid frame type specified",ColorReset)
	}
	if err != nil {
		log.Println(ColorRed,err,ColorReset)
		return
		//os.Exit(1)
	}
	//draw model
	var pltstr string
	if len(frmrez) != 0{
		if okdraw{
			pltstr = <-pltchn
			mod.Txtplots = append(mod.Txtplots, pltstr)
		}
		reportz,_ := frmrez[6].(string)
		mod.Report = reportz
		mod.Report += FrcTable(mod)
		// log.Println("REPORT",mod.Report)
		if mod.Calc{
			mod.CalcRez(frmtyp, frmrez)
		} else {
			ms, _ := frmrez[1].(map[int]*Mem)
			js, _ := frmrez[0].(map[int]*Node)
			mod.Ms = ms
			mod.Js = js
		}
		if term != "" && mod.Calc{
			switch mod.Frmstr{
				case "1db", "2df":
				mod.Txtplots = append(mod.Txtplots, PlotRez2d(mod, term))
			}
		}
		
		if !pipe && !mod.Web{
			fmt.Println(mod.Report)
			if term != ""{
				for i, txtplot := range mod.Txtplots{
					fmt.Println("plot #",i,"-\n",txtplot)
				}				
			}
		}
		if pipe{
			fname, err := mod.JsonOut()
			if err != nil{
				log.Println(err)
			} else {
				fmt.Printf("output file saved at ->%s",fname)
			}
		}
	}
	return
}

//CalcModPar analyses a model in a new goroutine - for opt routines
func CalcModPar(idx int, mod Model, rezchn chan []interface{}){
	//for parallel mod eval
	var frmrez, rez []interface{}
	var err error
	switch mod.Frmtyp{
		case 1:
		case 2:
		frmrez, err = CalcTrs(&mod,2)
		case 3:
		case 4:
		case 5:
		case 6:
	}
	
	pltchn := make(chan string)
	var pltstr string
	if mod.Term != ""{go PlotFrm2d(&mod, mod.Term, pltchn)}
	rez = make([]interface{},4)
	rez[0] = idx
	rez[1] = err
	rez[2] = frmrez
	if mod.Term != ""{pltstr = <- pltchn}
	rez[3] = pltstr
	//fmt.Println(pltstr)
	rezchn <- rez
}

//CalcRez calcs b.m/sf/deflections along each member
//stores max vals in each mem
func (mod *Model) CalcRez(frmtyp string, frmrez []interface{}){
	js, _ := frmrez[0].(map[int]*Node)
	ms, _ := frmrez[1].(map[int]*Mem)
	mod.Js = js
	mod.Ms = ms
	mod.Scales = make([]float64, 4)
	for _, mem := range mod.Ms{
		mem.CalcRez(frmtyp)
		//fmt.Println(mem.Vu, mem.Mu, mem.Dmax)
		if mod.Scales[3] < mem.Geoms[0]{
			mod.Scales[3] = mem.Geoms[0]
		}
		
		if mod.Scales[0] < mem.Vu{
			mod.Scales[0] = mem.Vu
		}
		
		if mod.Scales[1] < mem.Mu{
			mod.Scales[1] = mem.Mu
		}
		
		if mod.Scales[2] < mem.Dmax{
			mod.Scales[2] = mem.Dmax
		}
	}
	for i, val := range mod.Scales[:2]{
		//fmt.Println(i, val)
		mod.Scales[i] = mod.Scales[3]/val/4.0
	}
	//fmt.Println("scales->",mod.Scales)
	//return
}

func (mem *Mem) CalcRez(frmtyp string){
	var rl, ml, re, me, l, e, ar, iz float64
	switch frmtyp{
		case "1db":
		
		rl = mem.Qf[0]
		ml = mem.Qf[1]
		re = mem.Qf[2]
		me = mem.Qf[3]
		l = mem.Geoms[0]
		e = mem.Geoms[1]
		iz = mem.Geoms[2]
		if len(mem.Geoms) > 3{
			ar = mem.Geoms[3]
		}
		mem.Rez = Bmsfcalc(mem.Id, mem.Lds, l, e, ar, iz, rl, ml, re, me, false, mem.Clvr)
		mem.Vu = math.Abs(mem.Rez.Maxs[0])
		mem.Dmax = math.Abs(mem.Rez.Maxs[2])
		mem.Mu = math.Abs(mem.Rez.Maxs[1])
		if math.Abs(mem.Rez.Maxs[3]) > mem.Mu{
			mem.Mu = math.Abs(mem.Rez.Maxs[3])
		}
		case "2dt":
		case "2df":
		rl = mem.Qf[1]
		ml = mem.Qf[2]
		re = mem.Qf[4]
		me = mem.Qf[5]
		l = mem.Geoms[0]
		e = mem.Geoms[1]
		ar = mem.Geoms[2]
		iz = mem.Geoms[3]
		mem.Rez = Bmsfcalc(mem.Id, mem.Lds, l, e, ar, iz, rl, ml, re, me, false, mem.Clvr)
		mem.Vu = math.Abs(mem.Rez.Maxs[0])
		mem.Dmax = math.Abs(mem.Rez.Maxs[2])
		mem.Mu = math.Abs(mem.Rez.Maxs[1])
		if math.Abs(mem.Rez.Maxs[3]) > mem.Mu{
			mem.Mu = math.Abs(mem.Rez.Maxs[3])
		}
		case "3dt":
		case "3dg":
		case "3df":
	}
	//return
}

func (mem *MemNp) CalcRezNp(frmtyp string){
	var rl, ml, re, me, l, e, ar, iz float64
	switch frmtyp{
		case "1db":
		rl = mem.Qf[0]
		ml = mem.Qf[1]
		re = mem.Qf[2]
		me = mem.Qf[3]
		l = mem.Lspan
		e = mem.Em
		iz = mem.I0
		ar = mem.A0
		// l = mem.Geoms[0]
		// e = mem.Geoms[1]
		// iz = mem.Geoms[2]
		// if len(mem.Geoms) > 3{
		// 	ar = mem.Geoms[3]
		// }
		mem.Rez = Bmsfcalc(mem.Id, mem.Lds, l, e, ar, iz, rl, ml, re, me, false, mem.Clvr)
		mem.Vu = math.Abs(mem.Rez.Maxs[0])
		mem.Dmax = math.Abs(mem.Rez.Maxs[2])
		mem.Mu = math.Abs(mem.Rez.Maxs[1])
		if math.Abs(mem.Rez.Maxs[3]) > mem.Mu{
			mem.Mu = math.Abs(mem.Rez.Maxs[3])
		}
		case "2df":
		rl = mem.Qf[1]
		ml = mem.Qf[2]
		re = mem.Qf[4]
		me = mem.Qf[5]
		l = mem.Lspan
		e = mem.Em
		iz = mem.I0
		ar = mem.A0

		// l = mem.Geoms[0]
		// e = mem.Geoms[1]
		// ar = mem.Geoms[2]
		// iz = mem.Geoms[3]
		mem.Rez = Bmsfcalc(mem.Id, mem.Lds, l, e, ar, iz, rl, ml, re, me, false, mem.Clvr)
		mem.Vu = math.Abs(mem.Rez.Maxs[0])
		mem.Dmax = math.Abs(mem.Rez.Maxs[2])
		mem.Mu = math.Abs(mem.Rez.Maxs[1])
		if math.Abs(mem.Rez.Maxs[3]) > mem.Mu{
			mem.Mu = math.Abs(mem.Rez.Maxs[3])
		}
	}
	//return
}

//CalcRezNp processes np model results
func (mod *Model) CalcRezNp(frmtyp string){
	mod.Scales = make([]float64, 4)
	for _, mem := range mod.Mnps{
		mem.CalcRezNp(frmtyp)
		//fmt.Println(mem.Vu, mem.Mu, mem.Dmax)
		if mod.Scales[3] < mem.Lspan{
			mod.Scales[3] = mem.Lspan
		}
		
		if mod.Scales[0] < mem.Vu{
			mod.Scales[0] = mem.Vu
		}
		
		if mod.Scales[1] < mem.Mu{
			mod.Scales[1] = mem.Mu
		}
		
		if mod.Scales[2] < mem.Dmax{
			mod.Scales[2] = mem.Dmax
		}
	}
	for i, val := range mod.Scales[:2]{
		//fmt.Println(i, val)
		mod.Scales[i] = mod.Scales[3]/val/4.0
	}	
}

//CalcModSer calculates model results serially
func CalcModSer(mod *Model)(frmrez []interface{}, pltstr string, err error){
	switch mod.Frmtyp{
		case 1:
		case 2:
		log.Println(ColorRed, IconCubes,ColorReset)
		frmrez, err = CalcTrs(mod,2)
		pltchn := make(chan string, 1)
		if mod.Term != ""{
			go PlotTrs2d(mod, mod.Term, pltchn)
			pltstr = <- pltchn
		}
		fmt.Println(pltstr)
		case 3:
		
		case 4:
		case 5:
		case 6:
	}
	return
}
