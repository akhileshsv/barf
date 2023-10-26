package barf

import (
	"os"
	"fmt"
	"log"
	"math"
)

//CalcModNp - entry func from flags/menu for non prismatic beam-column frame analysis
//see hulse section 2.7/mosley spencer section 4.1
//IS UNDONE (not werk)
//basically what should it do? print stuff? call freecad, etc
func CalcModNp(mod *Model, frmtyp, term string) (err error){
	frmrez, err := CalcNp(mod, frmtyp, true)
	if err != nil{
		return
	}
	//frmrez[0] = js
	//frmrez[1] = ms
	//frmrez[2] = dglb
	//frmrez[3] = rnode
	//frmrez[4] = msqf
	//frmrez[5] = msloaded
	ms, _ := frmrez[1].(map[int]*MemNp)
	msloaded,_ := frmrez[5].(map[int][][]float64)
	//rezstring += fmt.Sprintf("%s\n",ex)
	//rezstring += fmt.Sprintf("%f\n",dglb)
	//rezstring += fmt.Sprintf("%.5f\n",rnode)
	//for m, mem := range ms{
	if mod.Calc{
		mod.CalcRezNp(frmtyp, frmrez)
	}
	for mem := 1; mem <= len(ms); mem++{
		r := Bmsfcalc(mem, msloaded[mem], ms[mem].Lspan, ms[mem].Em, ms[mem].A0, ms[mem].I0, ms[mem].Qf[0], ms[mem].Qf[1], ms[mem].Qf[2], ms[mem].Qf[3], true)
		log.Println("span->",mem)
		for j, bm := range r.BM{
			log.Println("section->",j,"moment->",bm," kn-m")
		}
	}
	reportz,_ := frmrez[6].(string)
	fmt.Println(reportz)
	return
}

//CalcMod - entry func from flags/menu for model direct stiffness analysis
//called after json - model read
//plots line diagram and prints report
func CalcMod(mod *Model, frmtyp, term string) (err error){
	pltchn := make(chan string, 1)
	var frmrez []interface{}
	okdraw := true
	if term == ""{okdraw = false}
	switch frmtyp{
		case "1db":
		if okdraw {go PlotBm1d(mod, term, pltchn)}
		frmrez, err = CalcBm1d(mod, 2)
		if mod.Calc{
			mod.CalcRez(frmtyp, frmrez)
		}
		case "2dt":
		if okdraw {go PlotTrs2d(mod, term, pltchn)}
		frmrez, err = CalcTrs(mod, 2)
		case "2df":
		if okdraw {go PlotFrm2d(mod, term, pltchn)}
		frmrez, err = CalcFrm2d(mod, 3)
		if mod.Calc{
			mod.CalcRez(frmtyp, frmrez)
			if okdraw{PlotFrm2dRez(mod, term)}
			os.Exit(1)
		}
		case "3dt":
		if okdraw {go PlotTrs3d(mod, term, pltchn)}
		frmrez, err = CalcTrs(mod, 3)
		case "3dg":
		if okdraw {go PlotGrd3d(mod, term, pltchn)}
		frmrez, err = CalcGrd(mod, 3)
		case "3df":
		if okdraw {go PlotGrd3d(mod, term, pltchn)}
		frmrez, err = CalcFrm3d(mod, 6)
		default:
		log.Println(ColorRed,"invalid frame type specified",ColorReset)
	}
	if err != nil {
		log.Println(ColorRed,err,ColorReset)
		os.Exit(1)
	}
	//draw model
	var pltstr string
	if len(frmrez) != 0 {
		if okdraw {
			pltstr = <-pltchn
			fmt.Println(pltstr)
		}
		reportz,_ := frmrez[6].(string)
		fmt.Println(reportz)
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
	return
}

//GetEndAct returns vals for bending moment and sf calcs
func (mem *Mem) CalcRez(frmtyp string){
	var rl, ml, re, me, l, e, ar, iz float64
	switch frmtyp{
		case "1db":
		case "2df":
		rl = mem.Qf[1]
		ml = mem.Qf[2]
		re = mem.Qf[4]
		me = mem.Qf[5]
		l = mem.Geoms[0]
		e = mem.Geoms[1]
		ar = mem.Geoms[2]
		iz = mem.Geoms[3]
		//TURN PLOT OFF later
		mem.Rez = Bmsfcalc(mem.Id, mem.Lds, l, e, ar, iz, rl, ml, re, me, false)
		//fmt.Println(mem.Rez.Txtplot)
		//fmt.Println(mem.Rez.Maxs)
		mem.Vu = math.Abs(mem.Rez.Maxs[0])
		mem.Dmax = math.Abs(mem.Rez.Maxs[2])
		mem.Mu = math.Abs(mem.Rez.Maxs[1])
		if math.Abs(mem.Rez.Maxs[3]) > mem.Mu{
			mem.Mu = math.Abs(mem.Rez.Maxs[3])
		}
		case "3dg":
		case "3df":
	}
	return
}
