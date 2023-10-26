package barf

import (
	"fmt"
	"errors"
	"strings"
	"github.com/olekukonko/tablewriter"
	kass "barf/kass"
)

//coltab is the goroutine func for column reports and plotting from design funcs (frm2ddz, subfrm) 
func coltab(spam bool, cols []int, r map[int]*RccCol, cchn chan string){
	var rez string
	for _, i := range cols{
		c := r[i]
		c.Table(false)
		rez += c.Report
		if spam{
			if c.Term == ""{c.Term = "svg"}
			_ = c.Plot(c.Term)
			rez += c.Txtplot[0]
		}
	}
	cchn <- rez
}

//bmtab is the goroutine func for beam reports and plotting from design funcs (frm2ddz, subfrm) 
func bmtab(spam bool, bms []int, r map[int][]*RccBm, bchn chan string){
	var rez string
	for _, i := range bms{
		//barr[1].Table(false)
		r[i][1].Table(false)
		pltstr := ""
		if spam{	
			for _, b := range r[i]{
				if b.Term == ""{b.Term = "svg"}
				PlotBmGeom(b, b.Term)
				pltstr += b.Txtplot[0]
			}
		}
		//f.Reports = append(f.Reports, rcbm[i][1].Report)
		rez += r[i][1].Report
		rez += pltstr
	}
	bchn <- rez
}

//Frm2dTable prints a table summary of a designed 2d rcc frame
func Frm2dTable(f *kass.Frm2d, printz bool){
	rcol, _ := f.Rez[0].(map[int]*RccCol)
	rcbm, _ := f.Rez[1].(map[int][]*RccBm)
	
	cchn := make(chan string,1)
	bchn := make(chan string, 1)
	go coltab(f.Spam, f.Cols, rcol, cchn)
	go bmtab(f.Spam, f.Beams, rcbm, bchn)
	/*
	   var pltstr string
	for _, i := range f.Cols{
		c := rcol[i]
		c.Table(false)
		if f.Spam{
			pltstr = PlotColGeom(c, "dumb")
			f.Report += pltstr
		}
		//f.Reports = append(f.Reports, c.Report)
		f.Report += c.Report
	}
	for _, i := range f.Beams{
		//barr[1].Table(false)
		rcbm[i][1].Table(false)
		pltstr = ""
		if f.Spam{	
			for _, b := range rcbm[i]{
				PlotBmGeom(b, "dumb")
				pltstr += b.Txtplot[0]
			}
		}
		//f.Reports = append(f.Reports, rcbm[i][1].Report)
		f.Report += rcbm[i][1].Report
	}
	*/
	rezstr := new(strings.Builder)
	if f.Term != "mono"{rezstr.WriteString(ColorPurple)}
	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"vol tot(m3)","vol rcc(m3)","wt stl(kg)","form area (m2)","cost (rs)"})
	table.SetCaption(true,"frame -> quantity take off")
	row := fmt.Sprintf("%.3f, %.3f, %.3f, %.3f, %.f\n",f.Quants[0],f.Quants[0],f.Quants[1],f.Quants[2], f.Kost)
	table.Append(strings.Split(row,","))
	table.Render()
	if f.Term != "mono"{rezstr.WriteString(ColorReset)}
	r := fmt.Sprintf("%s",rezstr)
	//f.Reports = append(f.Reports, r)
	ctab := <- cchn
	f.Report += ctab
	btab := <- bchn
	f.Report += btab
	f.Report += r
	if printz{
		fmt.Println(f.Report)
		//fmt.Println("***XXX***XXX***LIVE NUDES***XXX***XXX")
		//fmt.Println(f.Kost," rupeesees")
		//fmt.Println("vrcc->",f.Quants[0],"m3") 
		//fmt.Println("wstl->",f.Quants[1],"kg") 
		//fmt.Println("afw->",f.Quants[2],"m2") 

	}

}

//Frm2dDz is the entry func for 2d frame design (column/beam)
//a chain of monkeys typing 
//it is worthless when all the designs are clearly wreng/does not werk/IS BRAKE
//RCC COLUMNS ARE ELEMENTS OF SATAN
func Frm2dDz(f *kass.Frm2d) (err error, allcons int, emap map[int][]error){
	//call kass
	err = f.Calc()
	if err != nil{
		//log.Println(err)	
		return
	}
	//if f.Ldcalc > 1{Mrd(f)}
	rcbm := make(map[int][]*RccBm) //here tis map coz beams start wherever
	if f.Nomcvr == 0.0{f.Nomcvr = 30.0}
	efcvr := f.Nomcvr + 20.0
	bmchn := make(chan []interface{}, len(f.Bmenv))
	colchn := make(chan []interface{}, len(f.Colenv))
	for i, bm := range f.Bmenv{
		//f.Members[mdx] = append(f.Members[mdx],[]int{jb, je, em, cp, mrel})
		cp := f.Members[i][0][3]
		rcbm[i] = make([]*RccBm, 3)
		GetBmArr(rcbm[i], bm, f.Kostin, f.Fcks[1], f.Fys[1], f.Fyv, efcvr, f.DM, f.D1, f.D2, f.Dslb, f.Code, f.Styps[cp-1], false)
	}
	for i, bm := range f.Bmenv{
		go Bmdz(f.Code, rcbm[i], bm, f.Term, false, bmchn)
	}
	rcol := make(map[int]*RccCol, len(f.Colenv))  //fair enough, map it is
	for i := range f.Colenv{
		c := GetCol(f.Colenv[i], f.Kostin, f.Fcks[0], f.Fys[0], efcvr, 2, f.Code)
		rcol[i] = &c
		rcol[i].Init()
	}
	for i := range f.Colenv{
		go ColDz(rcol[i], colchn)
	}
	emap = make(map[int][]error)
	f.Quants = make([]float64, 3)
	f.Rez = make([]interface{},2)
	//fmt.Println(len(f.Bmenv),len(f.Colenv))
	for _ = range f.Bmenv{
		rez := <- bmchn
		idx, _ := rez[0].(int)
		cons, _ := rez[1].(int)
		errz, _ := rez[2].([]error)
		allcons += cons
		emap[idx] = append(emap[idx],errz...)
		if cons == 0{
			//f.Kost += rcbm[idx][1].Kost
			f.Quants[0] += rcbm[idx][1].Vtot
			f.Quants[1] += rcbm[idx][1].Wstl
			f.Quants[2] += rcbm[idx][1].Afw
		}  else{
			if f.Verbose{
				//for _, err := range errz{
				//	fmt.Println(ColorRed, "bm->",idx,err, ColorReset)
				//}
			}
		}
	}
	for _ = range f.Colenv{
		rez := <- colchn 		
		idx, _ := rez[0].(int)
		err, _ := rez[1].(error)
		emap[idx] = append(emap[idx],err)
		if err != nil{
			if f.Verbose{
				//fmt.Println(ColorRed, "col->",rcol[idx].Dims, rcol[idx].Pu, rcol[idx].Mux,idx, err, ColorReset)
			}
			allcons++
			//fmt.Println(err)
		} else {
			//fmt.Println(idx, ColorGreen, err, ColorReset)
			//f.Kost += rcol[idx].Kost
			f.Quants[0] += rcol[idx].Vtot
			f.Quants[1] += rcol[idx].Wstl
			f.Quants[2] += rcol[idx].Afw

		}
	}
	//fmt.Println(allcons)
	if allcons != 0{
		err = errors.New("strange frame design error")
		return
	}
	f.Dz = true
	f.Rez[0] = rcol
	f.Rez[1] = rcbm
	switch len(f.Kostin){
		case 3:
		f.Kost = f.Quants[0] * f.Kostin[0] + f.Quants[1] * f.Kostin[1] + f.Quants[2] * f.Kostin[2]
		default:
		f.Kost = f.Quants[0] * CostRcc + f.Quants[1] * CostStl + f.Quants[2] * CostForm
	}
	if f.Verbose{
		Frm2dTable(f, true)
		//TableF2d(f, true)
		
		switch f.Ldcalc{
			case -1:
			fmt.Println("HERE")
			jl, ml := f.Mod.SumFrcs()
			f.Mod.Msloads = ml
			f.Mod.Jloads = jl
			f.DrawMod("dumb")
			f.Mod.Msloads = f.Loadcons[1]
			f.Mod.Jloads = f.Jloadmap[1]
			jl, ml = f.Mod.SumFrcs()
			f.Mod.Msloads = ml
			f.Mod.Jloads = jl
			f.DrawMod("dumb")
			//for _, pltstr := range f.Txtplots{
			//	fmt.Println(pltstr)
			//}
			default:
		}
	}
	
	return
}

//Mrd redistributes 2d frame support moments
//calls CBeamDM multiple times 
func Mrd(f *kass.Frm2d){
	//get list of beams per floor
	if f.DM != 0.0{
		//fmt.Println("moment redistribution DM->",f.DM)
		for _, bmvec := range f.CBeams{
			CBeamDM(3, bmvec, f.Bmenv, f.DM, f.Ms, f.Mslmap)
		}
		//CBeamDM(3, f.Beams, f.Bmenv, f.DM, f.Ms, f.Mslmap)
	}
	return
}
//Frm2dDu!
//what is a du, please du explain ser
//i think it was meant to be a one way cs slab frame design, get edge and interior frame details?
//it is now one line of code 
func Frm2dDu(x,y,z,sbc float64){
	//needs x (span), y (height), z(n spans)
	//calc 1 way cs slab
	//get slb depth
	//if wind loads - calc em
	//calc frm loads and frame
}
