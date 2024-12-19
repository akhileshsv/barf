package barf

import (
	"fmt"
	"strings"
	"github.com/olekukonko/tablewriter"
	kass "barf/kass"
)

//coltab is the goroutine func for column reports from design funcs (frm2ddz, subfrm) 
func coltab(spam bool, cols []int, r map[int]*RccCol, cchn chan string){
	var rez string
	for _, i := range cols{
		c := r[i]
		if c.Ignore{continue}
		c.Table(false)
		rez += c.Report
	}
	cchn <- rez
}


//drawcols is the goroutine func for column plotting from design funcs (frm2ddz, subfrm) 
func drawcols(web bool, cols []int, r map[int]*RccCol, cchn chan []string){
	var rez []string
	for _, i := range cols{
		c := r[i]
		c.Web = web
		if c.Ignore{
			continue
		}
		rez = append(rez, c.PlotColDet())
	}
	cchn <- rez
}


//bmtab is the goroutine func for beam reports from design funcs (frm2ddz, subfrm) 
func bmtab(spam bool, bms []int, r map[int][]*RccBm, bchn chan string){
	var rez string
	for _, i := range bms{
		//barr[1].Table(false)
		r[i][1].Table(false)
		rez += r[i][1].Report
	}
	bchn <- rez
}


//drawbms is the goroutine func for beam plotting from design funcs (frm2ddz, subfrm) 
func drawbms(web bool, bms []int, r map[int][]*RccBm, bchn chan []string){
	var rez []string
	for _, i := range bms{
		for _, b := range r[i]{
			b.Web = web
			rez = append(rez, PlotBmGeom(b, b.Term))
		}
	}
	// for _, bmvec := range ceams{
	// 	CBeamDM(3, bmvec, f.Bmenv, f.DM, f.Ms, f.Mslmap)
	// }
	bchn <- rez
}

//drawbms is the goroutine func for beam plotting from design funcs (frm2ddz) 
func drawcbms(web bool, term, title string, cbms [][]int, bms map[int][]*RccBm, bchn chan []string){
	var rez []string
	
	for i, bmvec := range cbms{
		title = fmt.Sprintf("%s-%v",title,i+1)
		rez = append(rez, PlotFrmBmDet(web, bmvec, bms, "", title, term))
	}
	// for _, bmvec := range ceams{
	// 	CBeamDM(3, bmvec, f.Bmenv, f.DM, f.Ms, f.Mslmap)
	// }
	bchn <- rez
}


//Frm2dTable prints (and plots) a table summary of a designed 2d rcc frame
func Frm2dTable(f *kass.Frm2d, printz bool){
	rcol, _ := f.Rez[0].(map[int]*RccCol)
	rcbm, _ := f.Rez[1].(map[int][]*RccBm)
	
	cchn := make(chan string,1)
	bchn := make(chan string, 1)
	lpchn := make(chan []string, 1)
	cpchn := make(chan []string, 1)
	bpchn := make(chan []string,1)	
	cbchn := make(chan []string, 1)
	go coltab(f.Spam, f.Cols, rcol, cchn)
	go bmtab(f.Spam, f.Beams, rcbm, bchn)
	f.Noprnt = false
	if f.Term != ""{
		go kass.Drawflps(f, lpchn)
		go drawcols(f.Web, f.Cols, rcol, cpchn)
		if !f.Web{go drawbms(f.Web, f.Beams, rcbm, bpchn)}
		go drawcbms(f.Web, f.Term, f.Title, f.CBeams, rcbm, cbchn)
	}
	rezstr := new(strings.Builder)
	if f.Term != "mono"{rezstr.WriteString(ColorPurple)}
	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"vol tot(m3)","vol rcc(m3)","wt stl(kg)","form area (m2)","cost (rs)"})
	row := fmt.Sprintf("%.3f, %.3f, %.3f, %.3f, %.f\n",f.Quants[0],f.Quants[0],f.Quants[1],f.Quants[2], f.Kost)
	table.SetCaption(true,"frame -> quantity take off")
	table.Append(strings.Split(row,","))
	table.Render()
	if f.Term != "mono"{rezstr.WriteString(ColorReset)}
	r := rezstr.String()
	ctab := <- cchn
	f.Report += ctab
	btab := <- bchn
	f.Report += btab
	f.Report += r
	if f.Term != ""{
		lplts := <- lpchn
		f.Txtplots = append(f.Txtplots, lplts...)
		cplts := <- cpchn
		f.Txtplots = append(f.Txtplots, cplts...)
		if !f.Web{
			bplts := <-bpchn
			f.Txtplots = append(f.Txtplots, bplts...)
		}
		cbplts := <-cbchn
		f.Txtplots = append(f.Txtplots, cbplts...)
	}
	if printz{
		fmt.Println(f.Report)
	}

}

//Frm2dDz is the entry func for 2d frame design (column/beam)
//a chain of monkeys typing 
//it is worthless when all the designs are clearly wreng/does not werk/IS BRAKE
//RCC COLUMNS ARE ELEMENTS OF SATAN
func Frm2dDz(f *kass.Frm2d) (allcons int, emap map[int][]error, err error){
	//call kass
	err = f.Calc()
	var errstr string
	if err != nil{
		//log.Println(err)	
		return
	}
	if f.DM > 0.0 {
		Mrd(f)
	}
	rcbm := make(map[int][]*RccBm) //here tis map coz beams start wherever
	if f.Nomcvr == 0.0{f.Nomcvr = 30.0}
	efcvr := f.Nomcvr + 15.0
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
		//approximating for opt routines, ga goes wonky somehow
		if f.Opt > 0{
			c.Approx = true
		}
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
	for range f.Bmenv{
		rez := <- bmchn
		idx, _ := rez[0].(int)
		cons, _ := rez[1].(int)
		errz, _ := rez[2].([]error)
		allcons += cons
		if !f.Noprnt{emap[idx] = append(emap[idx],errz...)}
		f.Quants[0] += rcbm[idx][1].Vrcc
		f.Quants[1] += rcbm[idx][1].Wstl
		f.Quants[2] += rcbm[idx][1].Afw
		if cons > 0{
			errstr += fmt.Sprintf("error in beam mem %v\n",idx)
			for j, err := range errz{
				errstr += fmt.Sprintf("span %v-err-%s\n",j,err)
			}
		}
	}
	for range f.Colenv{
		rez := <- colchn 		
		idx, _ := rez[0].(int)
		err, _ := rez[1].(error)
		if !f.Noprnt{emap[idx] = append(emap[idx],err)}
		f.Quants[0] += rcol[idx].Vrcc
		f.Quants[1] += rcol[idx].Wstl
		f.Quants[2] += rcol[idx].Afw
		if err != nil{
			allcons++
			errstr += fmt.Sprintf("error in col dz mem %v \nerr %s\n",idx,err)
		}
	}
	//fmt.Println(allcons)
	if allcons != 0{
		//err = fmt.Errorf("frame design errors- %v\nerror map-%#v",allcons, emap)
		err = fmt.Errorf("frame design errors-%v \n error - %s",allcons, errstr)
		return
	}
	f.Dz = true
	f.Rez[0] = rcol
	f.Rez[1] = rcbm
	switch len(f.Kostin){
	//quants - vol rcc, wt stl, area form
		case 3:
		f.Kost = f.Quants[0] * f.Kostin[0] + f.Quants[1] * f.Kostin[1] + f.Quants[2] * f.Kostin[2]
		default:
		f.Kost = f.Quants[0] * CostRcc + f.Quants[1] * CostStl + f.Quants[2] * CostForm
	}
	if f.Noprnt{return}
	if f.Verbose{
		if f.Web{
			Frm2dTable(f, false)
		} else {
			Frm2dTable(f, true)
		}
	}
	//PlotCBmDet(bmvec []int, bms [][]*RccBm, folder, title, term string) (pltstr string)
	return
}

//Mrd redistributes 2d frame support moments
//calls CBeamDM multiple times 
func Mrd(f *kass.Frm2d){
	//get list of beams per floor
	if f.DM != 0.0{
		//fmt.Println("moment redistribution DMx->",f.DM)
		for _, bmvec := range f.CBeams{
			CBeamDM(3, bmvec, f.Bmenv, f.DM, f.Ms, f.Mslmap)
		}
		//CBeamDM(3, f.Beams, f.Bmenv, f.DM, f.Ms, f.Mslmap)
	}
	// return
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
