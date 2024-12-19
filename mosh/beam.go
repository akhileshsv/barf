package barf

import (
	"fmt"
	"time"
	"errors"
	"strings"
	kass"barf/kass"
	"github.com/olekukonko/tablewriter"
	"github.com/AlecAivazis/survey/v2"
)

var (
	//LMAO.fuckin' LMAO
	//1.48 - 28 percent gst (for rmc, 18 for steel/labor) + 20 percent margin (LMAO. fuckin' LMAO) 
	CostForm float64 = 80.0 * 10.76 * 1.38//0.42//
	CostRcc float64 = 5000.0 * 1.48//1.0//
	CostStl float64 = 80.0 * 1.38//0.01515

	//govindrajan/rajeev costs
	//CostForm float64 = 54.0 //0.42//
	//CostRcc float64 = 735.0//1.0//
	//CostStl float64 = 7.10//0.01515
)

//Quant takes off beam quantities from an array of left, mid, right sections 
//quantity of concrete, steel, formwork
//vol ck = lspan * area (check for styp - flange will be with slab formwork, etc)
//area of formwork afw -  perimeter * lspan - top or (2d + b)*lspan
func Quant(barr[]*RccBm)(err error){
	lspan := barr[1].Xs[20]
	if !barr[1].Dz{err = errors.New("strange design error"); return}
	barr[1].Vrcc = 0.0; barr[1].Wstl = 0.0;barr[1].Afw = 0.0
	astz := make([]float64,3)
	ascz := make([]float64,3)
	var l1, l2, l3, lxp float64
	barr[1].Vtot = barr[1].Sec.Prop.Area * 1e-6 * lspan
	//perimeter of shear stirrups - (WHAT ABOUT N LEGS)
	lxp = 2.0 * ((barr[1].Bw - 60.0)+(barr[1].Dused-60.0))
	switch barr[1].Nlegs{
		case 2:
		default:
		//subtract 2 from nlegs for lxp- this is wrong btw
		lxp = float64(barr[1].Nlegs) * ((barr[1].Bw - 60.0)+(barr[1].Dused-60.0))
	}
	lsup := barr[1].Lsx/1000.0 + barr[1].Rsx/1000.0 
	//fmt.Println("dims->",barr[1].Bf,barr[1].Bw,barr[1].Df,barr[1].Dused, barr[1].Sec.Prop.Perimeter)
	switch barr[1].Styp{
		case 1:
		barr[1].Vrcc = (lspan + lsup) * (barr[1].Bw * barr[1].Dused) * 1e-6 //- (lspan + lsup) * (barr[1].Bw * barr[1].Dslb) * 1e-6
		barr[1].Afw = (lspan + lsup) * (barr[1].Sec.Prop.Perimeter) * 1e-3 - barr[1].Bw * (lspan + lsup) * 1e-3 - barr[1].Dslb * (lspan + lsup) * 1e-3 * 2.0
		case 6,7,8,9,10:
		barr[1].Vrcc = (lspan + lsup) * (barr[1].Sec.Prop.Area) * 1e-6 - (lspan + lsup) * (barr[1].Bf * barr[1].Dslb) * 1e-6 
		//barr[1].Afw = lspan * (barr[1].Sec.Prop.Perimeter - 2.0 * barr[1].Df - barr[1].Bf) * 1e-3
		barr[1].Afw = (lspan + lsup) * (barr[1].Bw + 2.0 * (barr[1].Dused - barr[1].Df)) * 1e-3
		lxp = 2.0 * ((barr[1].Bw - 60.0)+(barr[1].Dused-60.0))	
		default:
		barr[1].Vrcc = (lspan + lsup) * (barr[1].Sec.Prop.Area) * 1e-6
		barr[1].Afw = (lspan + lsup) * (barr[1].Sec.Prop.Perimeter) * 1e-3
		//just use square stirrups? get bw sec at c.Cvrt, c.Cvrc
		//(or) get linksec and get perimeter is the way
		//lxp = 2.0 * ((barr[1].Bw - 60.0)+(barr[1].Dused-60.0))
	}
	//curtail lengths (bm.cls) as per govindraj paper (is super hacky, only for udls)
	//cl left top, cl right top, cl left bot, cl right bot, cl left anch, cl right anch
	var ldreq float64
	barr[1].CL = make([]float64, 6)
	dcl := (barr[1].Dused - barr[1].Cvrt)/1e3
	barr[1].CL[0] = 0.3 * lspan
	barr[1].CL[1] = 0.3 * lspan
	barr[1].CL[2] = 0.25 * lspan
	barr[1].CL[3] = 0.25 * lspan
	if barr[1].Ldx == 1{
		barr[1].CL[0] = 0.15 * lspan
		barr[1].CL[2] = 0.15 * lspan
	}
	if barr[1].Rdx == 1{
		barr[1].CL[1] = 0.15 * lspan
		barr[1].CL[3] = 0.15 * lspan
	}
	//adding just effd here, hoping tis' greater than 12.0 * diamax
	barr[1].CL[0] += dcl
	barr[1].CL[1] += dcl
	switch barr[1].Endc{
		case 0:
		dia := barr[1].Dia3
		if barr[1].Dia4 > dia{dia = barr[1].Dia4}
		ldreq = BarDevLen(barr[1].Fck, barr[1].Fy, dia)/1e3
		lsx := 0.0
		if barr[1].Lsx > 0.0{
			lsx = barr[1].Lsx
		} else {
			lsx = barr[1].Rsx
		}
		l1 = 0.0; l2 = lspan - lsx/2.0/1e3 + ldreq; l3 = 0.0
		astz[1] = barr[1].Ast; ascz[1] = barr[1].Asc
		case 1:
		dia := barr[1].Dia1
		if barr[1].Dia2 > dia{dia = barr[1].Dia2}
		
		ldreq = BarDevLen(barr[1].Fck, barr[1].Fy, dia)/1e3
		dia = barr[1].Dia3
		if barr[1].Dia4 > dia{dia = barr[1].Dia4}
		ldreq1 := BarDevLen(barr[1].Fck, barr[1].Fy, dia)/1e3
		//sheer disregard for any accuracy
		if ldreq1 > ldreq{ldreq = ldreq1}
		//l1 = 0.0; l2 = lspan - barr[1].Lsx/2.0/1e3 + -barr[1].Rsx/2.0/1e3 + 2.0*ldreq; l3 = 0.0
		//astz[1] = barr[1].Ast; ascz[1] = barr[1].Asc
		//astz[1] = barr[1].Ast; ascz[1] = barr[1].Asc
		l1 = 0.15 * lspan; l2 = lspan - barr[1].Lsx/2.0/1e3 -barr[1].Rsx/2.0/1e3; l3 = 0.15 * lspan
		for i := range barr{
			astz[i] = barr[i].Ast; ascz[i] = barr[i].Asc
		}
		barr[1].CL[4] = ldreq
		barr[1].CL[5] = ldreq
		
		case 2:
		
		dia := barr[1].Dia1
		if barr[1].Dia2 > dia{dia = barr[1].Dia2}
		ldreq = BarDevLen(barr[1].Fck, barr[1].Fy, dia)/1e3
		dia = barr[1].Dia3
		if barr[1].Dia4 > dia{dia = barr[1].Dia4}
		ldreq1 := BarDevLen(barr[1].Fck, barr[1].Fy, dia)/1e3
		if ldreq1 > ldreq{ldreq = ldreq1}
		l1 = barr[1].S1*lspan + ldreq; l2 = lspan * (1.0 - barr[1].S1 - barr[1].S3); l3 = barr[1].S3*lspan + ldreq
		for i := range barr{
			astz[i] = barr[i].Ast; ascz[i] = barr[i].Asc
		}
		if barr[1].Ldx == 1{
			barr[1].CL[4] = ldreq	
		}
		if barr[1].Rdx == 1{
			barr[1].CL[5] = ldreq	
		}
		
	}
	//vol rbr = (ast1 + asc1)*l1 + (ast2 + asc2) * l2 + (ast3 + asc3) * l3
	barr[1].Wstl = (astz[0] + ascz[0])*l1 + (astz[1] + ascz[1]) * l2 + (astz[0] + ascz[0])*l3
	//add 50% of midspan steel until face of support + anchorage length
	barr[1].Wstl += (lspan - l2 + 2.0 * ldreq - barr[1].Lsx/2.0/1e3 - barr[1].Rsx/2.0/1e3) * (astz[1])/2.0
	if barr[1].Asc == 0.0{
		//add two 12 mm dia bars at some l2 + ldreq
		barr[1].Wstl += 2.0 * RbrArea(12.0) * (l2 + BarDevLen(barr[1].Fck, barr[1].Fy, 12.0)/1e3)
	}
	//add side face rebar for deep beams
	if barr[1].Aside > 0.0{
		//b.Rbrside = []float64{nbars,0.0,dia,0.0,atot,atot,0.0, nrows, atot, 0.0, 0.0, b.Bw - b.Cvrc -b.Cvrt, spc, 2.0}
		barr[1].Wstl += barr[1].Rbrside[4] * (lspan/1e3)
	}
	//fmt.Println("linx->",barr[1].Nlx[3] * lxp * RbrArea(barr[1].Dlink)*1e-9*7850)
	barr[1].Wstl = barr[1].Wstl*1e-6*7850.0
	barr[1].Wstl += barr[1].Nlx[3] * lxp * RbrArea(barr[1].Dlink) * 1e-9 * 7850
	switch len(barr[1].Kostin){
		case 3:
		barr[1].Kost = barr[1].Vrcc * barr[1].Kostin[0] + barr[1].Afw * barr[1].Kostin[1] + barr[1].Wstl * barr[1].Kostin[2]
		default:
		barr[1].Kost = barr[1].Vrcc * CostRcc + barr[1].Afw * CostForm + barr[1].Wstl * CostStl
	}
	//barr[1].Kost = barr[1].Vrcc * CostRcc + barr[1].Afw * CostForm + barr[1].Wstl * CostStl
	//fmt.Println(ColorYellow)
	//vol links nlinks * perimeter * bar area
	//binding wire teehee - 15 kg / 1000 kg of rebar
	return
}

//Quant quantifies a single beam section
func (b *RccBm) Quant() (err error){
	//quantify a single section
	lspan := b.Lspan
	if lspan == 0.0{lspan = 1.0}
	//fmt.Println("LSPAN in beam->",lspan)
	b.Vrcc = 0.0; b.Wstl = 0.0;b.Afw = 0.0
	b.Vtot = b.Sec.Prop.Area * 1e-6 * lspan
	//perimeter of shear stirrups - (WHAT ABOUT N LEGS)
	lxp := 2.0 * ((b.Bw - b.Cvrt)+(b.Dused-b.Cvrt))

	switch b.Nlegs{
		case 2:
		default:
		//subtract 2 from nlegs for lxp- this is wrong btw
		lxp = float64(b.Nlegs) * ((b.Bw - 60.0)+(b.Dused-60.0))
	}
	//fmt.Println("dims->",b.Bf,b.Bw,b.Df,b.Dused, b.Sec.Prop.Perimeter)
	switch b.Styp{
		case 1:
		b.Vrcc = lspan * (b.Bw * b.Dused) * 1e-6 - lspan * (b.Bw * b.Dslb) * 1e-6
		b.Afw = lspan * (b.Sec.Prop.Perimeter) * 1e-3 - b.Bw * lspan * 1e-3 - b.Dslb * lspan * 1e-3 * 2.0
		case 6,7,8,9,10:
		b.Vrcc = lspan * (b.Sec.Prop.Area) * 1e-6 - lspan * (b.Bf * b.Dslb) * 1e-6 
		b.Afw = lspan * (b.Bf + 2.0 * (b.Dused - b.Df)) * 1e-3
		lxp = 2.0 * ((b.Bw - 60.0)+(b.Dused-60.0))	
		default:
		b.Vrcc = lspan * (b.Sec.Prop.Area) * 1e-6
		b.Afw = lspan * (b.Sec.Prop.Perimeter) * 1e-3
	}
	ldreq := BarDevLen(b.Fck, b.Fy, b.Dia1)/1e3
	b.Wstl = (lspan + ldreq) * (b.Ast + b.Asc) * 1e-6 * 7850.0
	if b.Shrdz{b.Wstl += b.Nlx[0] * lxp * RbrArea(b.Dlink) * 1e-9 * 7850}
	return
}

//Printz printz
func (b *RccBm) Printz() (rezstring string){
	types := []string{"cantilever","single span","continuous"}
	rezstring += fmt.Sprintf("rcc beam type - %v %s\n", b.Endc, types[b.Endc])
	rezstring += fmt.Sprintf("%s - ast (bottom) %.f mm2 asc (top) %.f mm2\n",b.Title,b.Ast,b.Asc)	
	if b.Shrdz{
		rezstring += fmt.Sprintf("link dia - %.f mm spacing - main %.f mm min %.f nom %.f total %.f nos\n",b.Dlink,b.Lspc[0],b.Lspc[1],b.Lspc[2],b.Nlx[3])
	}
	
	if b.Asc > 0.0{
		r := b.Rbrc
		rezstring += fmt.Sprintf("%s - dia %.0f mm no %.0f mm dia 2 %.0f mm no %.0f \nreq astl prov %.2f mm2 astl req %.2f mm2 diff %.2f mm2\n","top",r[2],r[0],r[3],r[1],r[4],r[5],r[6])
	}
	
	if b.Ast > 0.0{
		r := b.Rbrt
		rezstring += fmt.Sprintf("%s - dia 1 %.0f mm no %.0f mm dia 2 %.0f mm no %.0f \nreq astl prov %.2f mm astl req %.2f mm2 diff %.2f mm2\n","bottom",r[2],r[0],r[3],r[1],r[4],r[5],r[6])
	}
	return
}

//Table generates an ascii table report (god bless tablewriter)
func (b *RccBm) Table(printz bool){
	lspan := b.Lspan
	if len(b.Xs) > 19 {lspan = b.Xs[20]}
	if lspan == 0.0 {lspan = 1.0}
	rezstr := new(strings.Builder)
	hdr := fmt.Sprintf("%s\nrcc beam %s %s\ndate-%s\n%s\n",ColorYellow,b.Title,ColorGreen,time.Now().Format("2006-01-02"),ColorReset)
	rezstr.WriteString(hdr)
	//
	/*
	hdr = ""
	hdr += fmt.Sprintf("grade of concrete - M %.1f\nsteel - main Fe %.f links Fe %.f\n", b.Fck, b.Fy, b.Fyv)
	hdr += fmt.Sprintf("lspan - %.3f m cover - top %0.1f mm bottom %0.1f mm\n", lspan, b.Cvrc, b.Cvrt)
	hdr += fmt.Sprintf("design loads - mur %.1f kn.m, vur %0.1f kn\n", b.Mu, b.Vu)
	hdr += fmt.Sprintf("section - %s\ndimensions - %.f mm\n\n",kass.SectionMap[b.Styp],b.Dims)

	//HERE
	hdr += ColorReset
	*/
	//PlotBmGeom(b, "dumb")
	//hdr += fmt.Sprintf("%s",b.Txtplot[0])
	//rezstr.WriteString(hdr)
	psteel := 100.0 * (b.Ast + b.Asc)/b.Sec.Prop.Area
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)
	var row string
	table.SetCaption(true,"beam specs")
	table.SetHeader([]string{"concrete","steel(main)","steel(links)","cover-top\n(mm)","cover-bot.\n(mm)","steel %"})
	row = fmt.Sprintf("M%.f,Fe%.f,Fe%.f,%.2f,%.2f,%.2f",b.Fck,b.Fy,b.Fyv,b.Cvrc,b.Cvrt,psteel)
	table.Append(strings.Split(row,","))
	table.Render()

	rezstr.WriteString(ColorBlue)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"beam geometry")
	table.SetHeader([]string{"section","dimensions\n(mm)","span\n(m)","bf\n(mm)","d\n(mm)","bw\n(mm)","df\n(mm)"})
	row = fmt.Sprintf("%s,%.2f,%.3f,%.2f,%.2f,%.2f,%.2f",kass.SectionMap[b.Styp],b.Dims,b.Lspan,b.Bf,b.Dused,b.Bw,b.Df)
	table.Append(strings.Split(row,","))
	table.Render()


	rezstr.WriteString(ColorRed)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"beam loads")
	table.SetHeader([]string{"mu\n(kn-m)","vu\n(kn)"})
	row = fmt.Sprintf("%.2f,%.2f",b.Mu,b.Vu)
	table.Append(strings.Split(row,","))
	table.Render()
	
	rezstr.WriteString(ColorPurple)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"section reinforcement")
	r := b.Rbrt
	table.SetHeader([]string{"location","dia1(mm)","nos","dia2(mm)","nos","ast prov(mm2)","ast req(mm2)","diff(mm2)","nlayers"})
	if b.Ast > 0.0 && len(r) > 7{
		row = fmt.Sprintf("%s, %.0f, %.0f, %.0f, %.0f, %.2f, %.2f, %.2f, %.f\n","bottom",r[2],r[0],r[3],r[1],r[4],r[5],r[6],r[7])
		//rez += fmt.Sprintf("combo - %.f nos %.f mm dia %.f nos %.f mm dia  ast prov %.f mm2 ast req %.f mm2 a diff %.f mm2\n",r[0],r[2],r[1],r[3],r[4],r[5],r[6])
	} else {
		row = fmt.Sprintf("%s, %.0f, %.0f, %.0f, %.0f, %.2f ,%.2f ,%.2f, %.f\n","bottom",0.,0.,0.,0.0,0.0,0.0,0.0,0.0)
	}
	table.Append(strings.Split(row,","))
	row = ""
	r = b.Rbrc
	if b.Asc > 0.0 && len(r) > 7{
		row = fmt.Sprintf("%s, %.0f, %.0f, %.0f, %.0f, %.2f ,%.2f ,%.2f, %.f\n","top",r[2],r[0],r[3],r[1],r[4],r[5],r[6],r[7])
		//rez += fmt.Sprintf("combo - %.f nos %.f mm dia %.f nos %.f mm dia  ast prov %.f mm2 ast req %.f mm2 a diff %.f mm2\n",r[0],r[2],r[1],r[3],r[4],r[5],r[6])
	} else {
		row = fmt.Sprintf("%s, %.0f, %.0f, %.0f, %.0f, %.2f ,%.2f ,%.2f, %.f\n","top",0.,0.,0.,0.0,0.0,0.0,0.0,0.0)

	}
	table.Append(strings.Split(row,","))
	r = b.Rbrside
	if b.Aside > 0.0{
		row = fmt.Sprintf("%s, %.0f, %.0f, %.0f, %.0f, %.2f ,%.2f ,%.2f, %.f\n","side face",r[2],r[0],r[3],r[1],r[4],r[5],r[6],r[7])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rezstr.WriteString(ColorGreen)
	if b.Shrdz{
		table = tablewriter.NewWriter(rezstr)
		table.SetHeader([]string{"typ","dia(mm)","spacing","net length(m)","from","to","from","to","no."})
		table.SetCaption(true,"shear reinforcement")
		slink, smin, snom, mainlen, minlen, nomlen := b.Lspc[0],b.Lspc[1],b.Lspc[2],b.Lspc[3],b.Lspc[4],b.Lspc[5]
		for i, typ := range []string{"main","min","nominal"}{
			spc := slink; ltot := mainlen; nlx := b.Nlx[i]
			var l1, l2, l3, l4 float64
			l1 = 0.0; l2 = b.L1; l3 = b.L4
			if len(b.Xs) > 20{l4 = b.Xs[20]}
			switch i{
				case 0:
				case 1:
				ltot = minlen
				spc = smin
				l1 = b.L1; l2 = b.L2; l3 = b.L3; l4 = b.L4
				case 2:
				ltot = nomlen
				spc = snom
				l1 = b.L2; l2 = b.L3; l3 = 0.0; l4 = 0.0 
			}
			row = fmt.Sprintf("%s,%.f,%.f,%.2f,%.2f,%.2f,%.2f,%.2f,%.0f", typ, b.Dlink, spc, ltot, l1, l2, l3, l4,nlx)
			table.Append(strings.Split(row,","))
		}
		table.Render()
	}
	rezstr.WriteString(ColorBlue)
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"vol tot(m3)","vol rcc(m3)","wt stl(kg)","form area (m2)","cost (rs)","unit cost(rs/m)"})
	table.SetCaption(true,"quantity take off")
	row = fmt.Sprintf("%.3f, %.3f, %.3f, %.3f, %.f, %.2f\n",b.Vtot,b.Vrcc,b.Wstl,b.Afw, b.Kost, b.Kost/lspan)
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorReset)
	b.Report = rezstr.String()
	if printz{
		fmt.Println(b.Report)
	}
}

//Invert inverts a beam section (ast - asc, asc - ast and so on)
func (b *RccBm) Invert(){
	//WHAAT ABOUT RBR OPTS
	//b.ast = asc; b.asc = ast
	//for _, dbar := y = dused - y
	ast := b.Asc; asc := b.Ast
	var barpts [][]float64
	for _, pt := range b.Barpts{
		x := pt[0]; y := b.Dused - pt[1]
		barpts = append(barpts, []float64{x,y})
	}
	b.Ast = ast; b.Asc = asc
	b.Barpts = make([][]float64, len(barpts))
	for i := range barpts{
		b.Barpts[i] = make([]float64, 2)
		copy(b.Barpts[i],barpts[i])
	}
}

//Draw generates gnuplot data for a beam section
//change dis to something like col draw
func (b *RccBm) Draw() (data string, err error){
	var ldata string
	var coords, links [][]float64
	xue := b.Bw + 10.0
	lnky := b.Cvrt
	lnkx := 40.0
	switch b.Styp{
		//(will it ever be done) non standard section (similar to ColArXu)
		case 1:
		//rekt sekt
		coords = [][]float64{
			{0,0},
			{b.Bw,0},
			{b.Bw,b.Dused},
			{0,b.Dused},
			{0,0},
		}
		links = [][]float64{
			{lnkx,lnky},
			{b.Bw-lnkx,lnky},
			{b.Bw-lnkx,b.Dused-lnky},
			{lnkx,b.Dused-lnky},
			{lnkx,lnky},
		}
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bw, b.Dused/2.0, b.Dused)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bw/2.0, 0.0, b.Bw)
		ldata += fmt.Sprintf("%f %f %.fmm2\n",b.Bw/2.0, b.Cvrt + 60.0, b.Ast)
		ldata += fmt.Sprintf("%f %f %.fmm2\n",b.Bw/2.0, b.Dused - b.Cvrc - 60.0, b.Asc)
		case 7:
		//l right
		coords = [][]float64{
			{0,0},
			{b.Bw,0},
			{b.Bw,b.Dused - b.Df},
			{b.Bf,b.Dused - b.Df},
			{b.Bf,b.Dused},
			{0,b.Dused},
			{0,0},
		}
		links = [][]float64{
			{lnkx,lnky},
			{b.Bw-lnkx,lnky},
			{b.Bw-lnkx,b.Dused-lnky},
			{lnkx,b.Dused-lnky},
			{lnkx,lnky},
		}
		ldata += fmt.Sprintf("%f %f %.f\n",0.0, b.Dused/2.0, b.Dused)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bw/2.0, 0.0, b.Bw)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bw, (b.Dused - b.Df)/2.0, b.Dused - b.Df)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bf, b.Dused - b.Df/2.0, b.Df)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bf/2.0, b.Dused, b.Bf)
		ldata += fmt.Sprintf("%f %f %.fmm2\n",b.Bw/2.0, b.Cvrt + 60.0, b.Ast)
		ldata += fmt.Sprintf("%f %f %.fmm2\n",b.Bw/2.0, b.Dused - b.Cvrc - 60.0, b.Asc)
		case 6:
		//t section
		coords = [][]float64{
			{b.Bf/2.0 - b.Bw/2.0,0},
			{b.Bf/2.0 + b.Bw/2.0,0},
			{b.Bf/2.0 + b.Bw/2.0, b.Dused - b.Df},
			{b.Bf , b.Dused - b.Df},
			{b.Bf ,b.Dused},
			{0, b.Dused},
			{0, b.Dused - b.Df},
			{b.Bf/2.0 - b.Bw/2.0,b.Dused - b.Df},
			{b.Bf/2.0 - b.Bw/2.0,0},
		}
		links = [][]float64{
			{b.Bf/2.0 - b.Bw/2.0 + lnkx,lnky},
			{b.Bf/2.0 + b.Bw/2.0 - lnkx,lnky},
			{b.Bf/2.0 + b.Bw/2.0 - lnkx, b.Dused - lnky},
			{b.Bf/2.0 - b.Bw/2.0 + lnkx, b.Dused - lnky},
			{b.Bf/2.0 - b.Bw/2.0 + lnkx,lnky},
		}
		ldata += fmt.Sprintf("%f %f %.f\n",0.0, b.Dused/2.0, b.Dused)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bf/2.0, 0.0, b.Bw)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bw/2.0 + b.Bf/2.0, b.Dused/2.0, b.Dused - b.Df)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bf, b.Dused - b.Df/2.0, b.Df)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bf/2.0, b.Dused, b.Bf)
		ldata += fmt.Sprintf("%f %f %.fmm2\n",b.Bf/2.0, b.Cvrt + 60.0, b.Ast)
		ldata += fmt.Sprintf("%f %f %.fmm2\n",b.Bf/2.0, b.Dused - b.Cvrc - 60.0, b.Asc)
		xue = b.Bf + 10.0
		case 14:
		//pocket flange section
		coords = [][]float64{
			{0,0},
			{b.Bw,0},
			{b.Bw,b.Dused-b.Df},
			{b.Bw/2.0+b.Bf/2.0,b.Dused-b.Df},
			{b.Bw/2.0+b.Bf/2.0,b.Dused-b.Df},
			{b.Bw/2.0+b.Bf/2.0,b.Dused},
			{b.Bw/2.0-b.Bf/2.0,b.Dused},
			{b.Bw/2.0-b.Bf/2.0,b.Dused-b.Df},
			{0,b.Dused-b.Df},
			{0,0},
		}
		links = [][]float64{
			{b.Bw/2.0 - b.Bf/2.0+lnkx,lnky},
			{b.Bw/2.0 + b.Bf/2.0-lnkx,lnky},
			{b.Bw/2.0 + b.Bf/2.0-lnkx,b.Dused-lnky},
			{b.Bw/2.0 - b.Bf/2.0+lnkx,b.Dused-lnky},
			{b.Bw/2.0 - b.Bf/2.0+lnkx,lnky},
		}
		ldata += fmt.Sprintf("%f %f %.f\n",0.0, b.Dused/2.0, b.Dused)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bw/2.0, 0.0, b.Bw)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bw, (b.Dused - b.Df)/2.0, b.Dused - b.Df)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bw/2.0 + b.Bf/2.0, b.Dused - b.Df/2.0, b.Df)
		ldata += fmt.Sprintf("%f %f %.f\n",b.Bw/2.0, b.Dused, b.Bf)
		ldata += fmt.Sprintf("%f %f %.fmm2\n",b.Bw/2.0, b.Cvrt + 60.0, b.Ast)
		ldata += fmt.Sprintf("%f %f %.fmm2\n",b.Bw/2.0, b.Dused - b.Cvrc - 60.0, b.Asc)
		xue = b.Bw + 10.0
	}
	for _, c := range coords{
		data += fmt.Sprintf("%f %f %v\n",c[0],c[1],1)
	}
	data += "\n"
	for _, l := range links{
		data += fmt.Sprintf("%f %f %v\n",l[0],l[1],2)
	}
	data += "\n"
	//ldata += fmt.Sprintf("%f %f %.f\n",xue, b.Xu, b.Xu)
	data += "\n\n"
	switch b.Dz{
		case true:
		ldata += "\n\n"
		data += ldata
		for i, dia := range b.Dias{
			data += fmt.Sprintf("%f %f %f\n",b.Barpts[i][0],b.Barpts[i][1], dia/2.0) 
		}
		data += "\n\n"
		if b.Xu != 0.0{
			data += fmt.Sprintf("%f %f %v\n",0.0,b.Xu,0)
			data += fmt.Sprintf("%f %f %v\n",xue,b.Xu,0)	
		}
		case false:
		switch len(b.Dbars){
			case 0:
			//offset links by d/2 (say 10.0 mm) and draw 4 points	
			for _, l := range links{
				data += fmt.Sprintf("%f %f %v\n",l[0],l[1],2)
			}
			data += "\n\n"
			if b.Xu != 0.0{
				data += fmt.Sprintf("%f %f %v\n",0.0,b.Xu,0)
				data += fmt.Sprintf("%f %f %v\n",xue,b.Xu,0)	
			}
			
			default:
		}
	}
	return
}

//BarLay generates bar coords (barpts) and dia slices from rebar templates
func (b *RccBm) BarLay() (err error){
	//rez = []float64{0 float64(n1), 1 float64(n2),2 d1,3 d2,4 astmin,5 ast,6 adiff}
	//rez = []float64{7 nlayer,8 astprov,9 efcvr,10 efdp,11 cldis,12 clvdis,13 nbarRow} 
	var xc, xt, yc, yt, xs, ys float64
	if b.Styp == 0{
		switch b.Tyb{
			case 0.0:
			b.Styp = 1
			case 0.5:
			b.Styp = 7
			case 1.0:
			b.Styp = 6
		}
	}
	switch b.Styp{
		case 1:
		xt = b.Cvrt; yt = b.Cvrt; xc = b.Cvrc; yc = b.Dused - b.Cvrc
		case 7:
		xt = b.Cvrt; yt = b.Cvrt; xc = b.Cvrc; yc = b.Dused - b.Cvrc
		case 6:
		xt = b.Cvrt + b.Bf/2.0 - b.Bw/2.0; yt = b.Cvrt; xc = xt; yc = b.Dused - b.Cvrc
		case 14:
		xt = b.Cvrt + b.Bw/2.0 - b.Bf/2.0; yt = b.Cvrt; xc = xt; yc = b.Dused - b.Cvrc
		default:
		//HUH? HUH?
	}
	if b.Ast > 0.0{
		//nlay := int(b.Rbrt[7]); efcvrc := b.Rbrt[9]
		n1 := int(b.Rbrt[0]); n2 := int(b.Rbrt[1]); d1 := b.Rbrt[2]; d2 := b.Rbrt[3]//; asp := b.Rbrt[4];ast := b.Rbrt[5]
		xstep := b.Rbrt[11]; ystep := b.Rbrt[12]; nbarr := int(b.Rbrt[13])
		//fmt.Println("ast from bar lay-nlayers->",nlay,"nbarr/row->",nbarr,"efcvr calc->",efcvrc,"beam->",b.Cvrt, "xstep, ystep->",xstep, ystep)
		//here
		//rez = []float64{nlayer, astprov, efcvr, efdp, cldis, clvdis, nbarRow}
		switch n1 + n2{
			case 1:
			//HUH? HUH?
			case 2:
			switch b.Styp{
				case 14:
				xstep = b.Bf - 2.0 * b.Cvrt
				default:
				xstep = b.Bw - 2.0 * b.Cvrt
			}
			case 3:
			switch b.Styp{
				case 14:
				xstep = (b.Bf - 2.0 * b.Cvrt)/2.0
				default:
				xstep = (b.Bw - 2.0 * b.Cvrt)/2.0
			}
			default:
			/*
			switch b.Styp{
				case 14:
				xstep = (b.Bf - 2.0 * b.Cvrt)/float64(n1+n2-1)
				default:
				xstep = (b.Bw - 2.0 * b.Cvrt)/float64(n1+n2-1)
			}
			*/
                   
		}
		
		//dmax := d2
		xs  = xt ; ys = yt
		//if d1 > dmax || n2 == 0 {dmax = d1}
		for i := 1; i <= n1 + n2; i++{
			//fmt.Println("bar number->",i, xs, ys)
			if i <= n1{
				b.Dias = append(b.Dias, d1)
			} else {
				b.Dias = append(b.Dias, d2)
			}
			//fmt.Println("adding bar point->",xs,ys)
			b.Barpts = append(b.Barpts, []float64{xs, ys})
			//fmt.Println("adding xstep->",xstep," to xs->", xs)
			xs += xstep
			if i % nbarr == 0{
				//fmt.Println("resetting row->")
				xs = xt; ys = yt + ystep
			}
		}
	}
	if b.Asc > 0.0{
		//nlay := int(b.Rbrc[7]); efcvrc := b.Rbrc[9]
		n1 := int(b.Rbrc[0]); n2 := int(b.Rbrc[1]); d1 := b.Rbrc[2]; d2 := b.Rbrc[3]//; asp := b.Rbrc[4];ast := b.Rbrc[5]
		xstep := b.Rbrc[11]; ystep := b.Rbrc[12]; nbarr := int(b.Rbrc[13])
		//dmax := d2
		//fmt.Println("asc from bar lay-nlayers->",nlay,"nbarr/row->",nbarr,"efcvr calc->",efcvrc,"beam->",b.Cvrc,"xstep, ystep->",xstep, ystep)
		switch n1 + n2{
			case 1:
			//HUH? HUH?
			case 2:
			switch b.Styp{
				case 14:
				xstep = b.Bf - 2.0 * b.Cvrc
				default:
				xstep = b.Bw - 2.0 * b.Cvrc
			}
			case 3:
			switch b.Styp{
				case 14:
				xstep = (b.Bf - 2.0 * b.Cvrc)/2.0
				default:
				xstep = (b.Bw - 2.0 * b.Cvrc)/2.0
			}
			default:
			//switch b.Styp{
			//	case 14:
			//xstep = (b.Bf - 2.0 * b.Cvrc)/float64(n1+n2-1)
			//	default:
			//xstep = (b.Bw - 2.0 * b.Cvrc)/float64(n1+n2-1)
			//}
		}
		
		//xs = xc + d1/2.0; ys = yc
		xs = xc ; ys = yc
		//if d1 > dmax || n2 == 0 {dmax = d1}
		for i := 1; i <= n1 + n2; i++{
			if i <= n1{
				b.Dias = append(b.Dias, d1)
			} else {
				b.Dias = append(b.Dias, d2)
			}
			b.Barpts = append(b.Barpts, []float64{xs, ys})
			xs += xstep
			if i % nbarr == 0{
				xs = xc; ys = yc - ystep
			}
		}
	}
	//fmt.Println("dias, barpts->")
	//fmt.Println(b.Dias)
	//fmt.Println(b.Barpts)
	return
}


//DrawSec is an attempt to make beam plot funcs more like col plot funcs (sigh)
func (b *RccBm) DrawSec(){
	
}

//BarPrint prints rebar details of a beam section
//does so horribly, just use tablewriter (may the Lord preserve)
func (b *RccBm) BarPrint() (rez string){
	r := b.Rbrt
	//fmt.Println("HERE->",r,b.Ast)
	if b.Ast > 0.0{
		rez += fmt.Sprintf("tension (bottom) steel\n number of layers %.f\n",r[7])
		rez += fmt.Sprintf("combo - %.f nos %.f mm dia %.f nos %.f mm dia  ast prov %.f mm2 ast req %.f mm2 a diff %.f mm2\n",r[0],r[2],r[1],r[3],r[4],r[5],r[6])

	}
	r = b.Rbrc
	if b.Asc > 0.0{
		rez += fmt.Sprintf("compression (top) steel\n number of layers %.f\n",r[7])
		rez += fmt.Sprintf("combo - %.f nos %.f mm dia %.f nos %.f mm dia  ast prov %.f mm2 ast req %.f mm2 a diff %.f mm2\n",r[0],r[2],r[1],r[3],r[4],r[5],r[6])
	}
	if b.Shrdz{
		rez += fmt.Sprintf("link spacing - main %.f mm min %.f mm nominal %.f mm\nlength - main %.2f m min %.2f m nominal %.2f m\n",b.Lspc[0],b.Lspc[1],b.Lspc[2],b.Lspc[3],b.Lspc[4],b.Lspc[5])
	}
	return
}

//Tweakz is a menu to tweak beam design options and etc
//if this works it would be quite functional
func (b *RccBm) Tweakz(){
	//menu to tweak rebar opts and what not

	running := true
	for running{
		var choice int
		prompt := &survey.Select{
		Message: "choose what to tweak",
			Options: []string{"rebar (top)","rebar (bottom)","exit"},
		}
		survey.AskOne(prompt, &choice)
		switch choice{
			case 2:
			return
			case 0:
			for i := range b.Rbrc{
				fmt.Println(i, b.Rbrc[i])
			}
			case 1:
			
			for i := range b.Rbrt{
				fmt.Println(i, b.Rbrt[i])
			}
		}
	}
}
/*
func BeamFrm2d(frmrez map[int]){
	//
	spanchn := make(chan BeamRez,len(msloaded))
	bmresults := make(map[int]BeamRez)
	for member, ldcases := range msloaded {
		go BeamFrc(ncjt, member, ms[member], ldcases, spanchn,plotbm)
	}
	for i :=0; i < len(msloaded); i++ {
		r := <- spanchn
		bmresults[r.Mem] = r

	}

}

func BeamFrm3d(){
	
        }
*/
/*
func CBeamBmSf(mod *kass.Model,frmrez []interface{},plotbm bool) (map[int]BeamRez){
	//continuous beam bm and sf along span
	//hulse 2.2
	ncjt := 2
	ms,_ := frmrez[1].(map[int]*kass.Mem)
	msloaded, _ := frmrez[5].(map[int][][]float64)
	spanchn := make(chan BeamRez,len(msloaded))
	bmresults := make(map[int]BeamRez)
	for member, ldcases := range msloaded {
		go BeamFrc(ncjt, member, ms[member], ldcases, spanchn,plotbm)
	}
	for i :=0; i < len(msloaded); i++ {
		r := <- spanchn
		bmresults[r.Mem] = r

	}
	if plotbm {
		for i:= 1; i <= len(bmresults); i++{
			bm := bmresults[i]
			//bm.TxtPlot = PlotBmSfBm(xs, vxs, mxs, dxs, l, true)
			//fmt.Println(kass.ColorYellow,"member--",bm.Mem,kass.ColorReset)
			for i, vx := range bm.SF {
				fmt.Println(kass.ColorCyan)
				fmt.Printf("Div %d SF %.2f KN BM %.3f KN-m Def %.3f mm", i,vx, bm.BM[i],1000*bm.Dxs[i])
				log.Println("SPAN ",i,"SF ",vx, " Kn BM ", math.Ceil(bm.BM[i]*100)/100, " Kn-M DEF ", 1000*bm.Dxs[i]," mm")
			}
			fmt.Println(kass.ColorPurple)
			fmt.Printf("Max SF %.3f at %.3f \nMax BM %.3f at %.3f\nMax def %.3f at %.3f",bm.Maxs[0],bm.Locs[0],bm.Maxs[1],bm.Locs[1],bm.Maxs[2],bm.Locs[2])
			fmt.Println(kass.ColorGreen)
			for i, cfx := range bm.Cfxs{
				fmt.Println("Countarr Flaxsures",i, cfx)
			}
			fmt.Println(kass.ColorReset)
			fmt.Println(bm.Txtplot)
		}
	}
	return bmresults
}


*/
