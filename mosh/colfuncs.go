package barf

//steve holt!
//colfuncs!

import (
	"fmt"
	"time"
	"math"
	"strings"
	"github.com/olekukonko/tablewriter"
	kass"barf/kass"
)

//pt is again, a point. WHY is this needed wtf
type pt struct{
	//why is this everywhere
	x,y float64
}

//Printz printz
func (c *RccCol) Printz() {
	fmt.Println(c.Type, " rcc column")
	fmt.Println("dims -", c.B, c.H)
	fmt.Println("grade of concrete, steel - ", c.Fck, c.Fy)
	fmt.Println("cover - compression face, tension face -", c.Cvrc, c.Cvrt)
	fmt.Println("number of rebar layers -",c.Nlayers)
	if len(c.Dias) != 0 {
		fmt.Println("percentage, area of steel -", 100.0*c.Asteel/c.B/c.H, c.Asteel)
		fmt.Println("xu -", c.Xu)
	}
}

//Table builds a column report ascii table
//in tablewriter it trusts (peace be upon tablewriter writer)
func (c *RccCol) Table(printz bool){
	rezstr := new(strings.Builder)
	//lspan := c.Lspan
	var hdr string
	if c.Term == "mono"{
		hdr = fmt.Sprintf("rcc col %s \ndate-%s\n",c.Title,time.Now().Format("2006-01-02"))
		
	} else {
		hdr = fmt.Sprintf("%s\nrcc col %s %s\ndate-%s\n%s\n",ColorYellow,c.Title,ColorGreen,time.Now().Format("2006-01-02"),ColorReset)

	}
	//rezstr.WriteString(hdr)

	//rezstr.WriteString(ColorCyan)
	/*
	hdr = ""
	hdr += fmt.Sprintf("grade of concrete M%.1f\nsteel - main Fe %.f, links Fe %.f\n", c.Fck, c.Fy, c.Fy)
	hdr += fmt.Sprintf("cover - top %0.1f mm bottom %0.1f mm x %.1f y %.1f\n", c.Cvrc, c.Cvrt, c.C1, c.C2)
	hdr += fmt.Sprintf("length - c-c %0.2f m eff.length (x) %.1f eff.length (y) %.1f\n", c.Lspan, c.Leffx, c.Leffy)
	hdr += fmt.Sprintf("slender? %t\n", c.Slender)
	hdr += fmt.Sprintf("design loads - pur %.1f kn murx %.1f kn.m mury %0.1f kn\n", c.Pu, c.Mux, c.Muy)
	hdr += fmt.Sprintf("section - %s\ndimensions - %.f mm\n",kass.SectionMap[c.Styp],c.Dims)
	*/
	var rtyp string
	switch c.Rtyp{
		case 0:
		rtyp = "symmetrical - 2 per layer"
		case 1:
		rtyp = "unsymmetrical - 2 layers"
		case 2:
		switch c.Rbrlvl{
			case 0,1,2:
			rtyp = "at corners"
		}
	}
	/*
	hdr += fmt.Sprintf("reinforcement type - %s \nno. of layers - %v\narea of steel - %.f mm2 percentage (steel)- %.3f\n",rtyp, c.Nlayers,c.Asteel,c.Psteel)
	hdr += fmt.Sprintf("lateral ties - dia %.f mm pitch %.f mm\n",c.Dtie,c.Ptie)
	*/
	hdr += ColorReset
	rezstr.WriteString(hdr)

	rezstr.WriteString(ColorCyan)
	
	table := tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"column specs")
	table.SetHeader([]string{"concrete","steel\n(main)","steel\n(links)","cover-top\n(mm)","cover-bot.\n(mm)","cover-x\n(mm)","cover-y\n(mm)"})
	row := fmt.Sprintf("M%.f,Fe%.f,Fe%.f,%.2f,%.2f,%.2f,%.2f",c.Fck,c.Fy,c.Fy,c.Cvrc,c.Cvrt,c.C1,c.C2)
	table.Append(strings.Split(row, ","))
	table.Render()

	
	rezstr.WriteString(ColorYellow)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"column geometry")
	table.SetHeader([]string{"section","dims\n(mm)","length\n(m)","eff.len(x)\n(m)","eff.len(y)\n(m)","slender?"})
	row = fmt.Sprintf("%s,%.2f,%.2f,%.2f,%.2f,%t",kass.SectionMap[c.Styp],c.Dims,c.Lspan,c.Leffx,c.Leffy,c.Slender)
	table.Append(strings.Split(row, ","))
	table.Render()

	rezstr.WriteString(ColorRed)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"design loads")
	table.SetHeader([]string{"pu\n(kn)","mx\n(kn-m)","my\n(kn-m)","pur\n(kn)","murx\n(kn-m)","mury\n(kn-m)"})
	row = fmt.Sprintf("%.2f,%.2f,%.2f,%.2f,%.2f,%.2f",c.Pu,c.Mux,c.Muy,c.Pur,c.Murx,c.Mury)
	table.Append(strings.Split(row, ","))
	table.Render()
	

	rezstr.WriteString(ColorCyan)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"rebar details")
	table.SetHeader([]string{"reinf.type","nlayers","steel area\n(mm2)","percent steel\n(mm)"})
	row = fmt.Sprintf("%s,%v,%.f,%.3f",rtyp,c.Nlayers,c.Asteel,c.Psteel)
	table.Append(strings.Split(row, ","))
	table.Render()	
	//HERE
	//pltstr := PlotColGeom(c, "dumb")
	//hdr += pltstr

	//rezstr.WriteString(hdr)
	rezstr.WriteString(ColorPurple)

	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"section reinforcement")
	table.SetHeader([]string{"location","dia1(mm)","nos","dia2(mm)","nos","ast prov(mm2)","ast req(mm2)","diff(mm2)"})
	var r []float64
	if len(c.Rbropt) >= 1{
		r = c.Rbropt[0]
		if len(r) > 6{
			row = fmt.Sprintf("%s, %.0f, %.0f, %.0f, %.0f, %.2f, %.2f, %.2f\n","2/layer",r[0],r[2],r[1],r[3],r[4],r[5],r[6])
			table.Append(strings.Split(row,","))
		}
	} else {
		//print from rbrcopt and rbrtopt
		r = c.Rbrc
		if len(r) > 6{
			row = fmt.Sprintf("%s, %.0f, %.0f, %.0f, %.0f, %.2f, %.2f, %.2f\n","top layer",r[0],r[2],r[1],r[3],r[4],r[5],r[6])
			table.Append(strings.Split(row,","))
		}
		
		r = c.Rbrt
		if len(r)> 6{
			row = fmt.Sprintf("%s, %.0f, %.0f, %.0f, %.0f, %.2f, %.2f, %.2f\n","bottom layer",r[0],r[2],r[1],r[3],r[4],r[5],r[6])
			table.Append(strings.Split(row,","))
		}

	}
	table.Render()
	var lxp float64
	switch c.Styp{
		case 0:
		lxp = (c.Dims[0] - 30.0) * math.Pi * 1e-3  
		default:
		lxp = c.Lsec.Prop.Perimeter * 1e-3
	}
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"column ties")
	table.SetHeader([]string{"dia(mm)","pitch/spacing(mm)","perimeter(m)","nties","net length(m)"})
	row = fmt.Sprintf("%.f,%.f,%.3f,%.f,%.3f",c.Dtie,c.Ptie,lxp,c.Nties,float64(c.Nties)*lxp)
	table.Append(strings.Split(row,","))
	table.Render()
	
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"plastic hinge ties")
	table.SetHeader([]string{"dia(mm)","pitch/spacing(mm)","length from joint(m)"})
	row = fmt.Sprintf("%.f,%.f,%.3f",c.Dtie,c.Ptiec,c.Lp)
	table.Append(strings.Split(row,","))
	table.Render()
	if c.Term != "mono"{rezstr.WriteString(ColorGreen)}
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"vol rcc(m3)","wt stl(kg)","form area (m2)","cost (rs?)"})
	table.SetCaption(true,"quantity take off")
	row = fmt.Sprintf("%.3f, %.3f, %.3f, %.f\n",c.Vrcc,c.Wstl,c.Afw, c.Kost)
	table.Append(strings.Split(row,","))
	table.Render()
	if c.Term != "mono"{rezstr.WriteString(ColorCyan)}
	hdr = ""
	hdr += fmt.Sprintf("total steel - %.f mm2 ast %.f mm2 asc %.f mm2 percentage %.3f\n", c.Asteel, c.Ast, c.Asc, c.Psteel)
	rezstr.WriteString(hdr)
	if c.Term != "mono"{rezstr.WriteString(ColorReset)}
	c.Report = fmt.Sprintf("%s",rezstr)
	if printz{
		fmt.Println(c.Report)
	}
	return
}

//MaxOf returns the max of a bunch of float values
func MaxOf(vals []float64)(mval float64){
	mval = vals[0]
	for _, val := range vals{
		if mval < val{
			mval = val
		}
	}
	return
}

//Quant takes off rcc column quantities
//the calculation is WRENG but eh
func (c *RccCol) Quant(){
	//haha. HAHAHA
	lspan := c.Lspan
	if lspan == 0.0{
		lspan = 3.0
		c.Lspan = 3.0
	}
	if len(c.Coords) == 0{
		c.Coords = make([][]float64, 2)
		c.Coords[0] = []float64{0,0}
		c.Coords[1] = []float64{0, c.Lspan}
	}
	//fmt.Println("LSEC=LSEC+++++->")
	//fmt.Printf("+%v\n",c.Lsec.Prop)
	//fmt.Println(ColorWhite,"lspan->",lspan,"area->",c.Sec.Prop.Area,"m2",ColorReset)
	c.Vrcc = 0.0; c.Wstl = 0.0;c.Afw = 0.0
	var lxp, nlx float64
	c.Vrcc = c.Sec.Prop.Area * 1e-6 * lspan
	//perimeter of shear stirrups - (WHAT ABOUT N LEGS)
	c.Afw = c.Sec.Prop.Perimeter * 1e-3 * lspan
	switch c.Styp{
		case 0:
		lxp = (c.Dims[0] - 30.0) * math.Pi * 1e-3  
		default:
		lxp = c.Lsec.Prop.Perimeter * 1e-3 * math.Ceil(float64(len(c.Dias))/4.0)
	}
	//fmt.Println("lxp, perimeter->",c.Lsec.Prop.Perimeter,lxp)
	//fmt.Println("")
	if c.L0 == 0.0{
		c.L0 = c.Lspan //- 0.2 //easier this way
	}
	//length of plastic hinge lp = dmax or
	c.Lp = math.Ceil(c.L0/6.0)
	c.Ptiec = 6.0 * c.D1
	if c.Ptiec == 0.0 || c.D2 < c.D1{
		c.Ptiec = 6.0 * c.D2
	}
	
	for _, dim := range c.Dims{
		if c.Lp < dim/1000.0{
			c.Lp = dim/1000.0
		}
		if c.Ptiec > dim/4.0{
			c.Ptiec = math.Round(dim/4.0)
		}
	}
	if c.Lp > 0.45{c.Lp = 0.45}
	if c.Ptiec > 100.0{c.Ptiec = 100.0}
	if c.Ptiec < 75.0{c.Ptiec = 75.0}
	nlx = math.Ceil((lspan - 2.0 * c.Lp)*1e3/c.Ptie+1.0)
	nlx += math.Ceil((2.0 * c.Lp)*1e3/c.Ptiec + 1.0)
	c.Nties = nlx
	c.Wstl = c.Asteel * lspan * 1e-6 * 7850.0
	c.Wstl += nlx * lxp * RbrArea(c.Dtie) * 1e-6 * 7850.0
	switch len(c.Kostin){
		case 3:
		c.Kost = c.Vrcc * c.Kostin[0] + c.Afw * c.Kostin[1] + c.Wstl * c.Kostin[2]
		default:
		c.Kost = c.Vrcc * CostRcc + c.Afw * CostForm + c.Wstl * CostStl
	}
	//fmt.Println(ColorYellow)
	return

}

//ColCircArXu returns the area of a circle segement of radius r at a depth dck from extreme point
//https://en.wikipedia.org/wiki/Circular_segment ("imagine citing wikipedia")
//https://mathworld.wolfram.com/CircularSegment.html

func ColCircArXu(r, dck float64) (area, xc, yc float64){
	r1 := r - dck
	area = math.Pow(r,2) * math.Acos(r1/r) - r1 * math.Sqrt(math.Pow(r,2)-math.Pow(r1,2))
	fmt.Println("dck, area",dck, area)
	return
}

//ColSecArXu returns the area of the concrete block at a depth dck from top
//it computes the point of intersection of the line y = ymax - xu
//with all section polygon line segments
//calcs area from kass.SecPrp
func ColSecArXu(sec *kass.SectIn, dck float64) (area, xc, yc float64) {
	var nc1s []int
	var wt1s []float64
	var coords [][]float64
	var n1, nc1, n int
	c1 := sec.Ymx - dck
	var xc1, yc1 float64
	for idx, nc := range sec.Ncs {
		var pts [][]float64
		ptmap := make(map[pt]int)
		if idx == 0 {
			n1 = 0
		} else {
			n1 = sec.Ncs[idx-1]
		}
		nc1 = 0; n = 0
		for i := range sec.Coords[n1 : n1+nc-1] {
			i = i + n1
			pta := pt{sec.Coords[i][0], sec.Coords[i][1]}
			ptb := pt{sec.Coords[i+1][0], sec.Coords[i+1][1]}
			if (pta.y < c1 && ptb.y < c1) {
				continue
			}
			if pta.y >= c1 {
				if _, ok := ptmap[pta]; !ok {
					ptmap[pta] = idx
					pts = append(pts, []float64{pta.x, pta.y})
					nc1++
					xc1 += pta.x; yc1 += pta.y; n++
				}
			}
			if ptb.y >= c1 {
				if _, ok := ptmap[ptb]; !ok {
					ptmap[ptb] = idx
					pts = append(pts, []float64{ptb.x, ptb.y})
					nc1++
					xc1 += ptb.x; yc1 += ptb.y; n++
				}
			}
			if pta.y - ptb.y == 0 {
				continue
			}
			a2 := ptb.y - pta.y
			b2 := pta.x - ptb.x
			c2 := a2 * pta.x + b2 * pta.y 
			xin := (c2 - b2*c1)/a2
			yin := c1
			ptx := pt{xin, yin}
			if (pta.x <= xin && xin <= ptb.x) || (pta.x >= xin && xin >= ptb.x){
				if _, ok := ptmap[ptx]; !ok {
					pts = append(pts, []float64{ptx.x, ptx.y})
					ptmap[ptx] = idx
					nc1++
					xc1 += ptx.x; yc1 += ptx.y; n++
				}
			}
		}
		if len(pts) == 0 {
			continue
		}
		kass.SortCcw(pts, xc1/float64(n), yc1/float64(n))
		pts = append(pts, pts[0])
		nc1++
		nc1s = append(nc1s, nc1)
		coords = append(coords, pts...)
		wt1s = append(wt1s, sec.Wts[idx])
	}
	area, xc, yc, _, _, _, _, _, _ = kass.SecPrp(nc1s, wt1s, coords)
	return
}

//ColFlip flips a column from major to y-axis 
func ColFlip(c *RccCol) (cy *RccCol){
	//flips a col (cx to cy)
	cy = &RccCol{
		Fck:c.Fck,
		Fy:c.Fy,
		Cvrc:c.C2,
		Cvrt:c.C2,
		Lx:c.Lx,
		Ly:c.Ly,
		Nomcvr:c.Nomcvr,
		Efcvr:c.Efcvr,
		Nlayers:c.Nlayers,
		Styp:c.Styp,
		B:c.H,
		H:c.B,
		Pu:c.Pu,
		Mux:c.Muy,
		Muy:c.Mux,
		C1:c.C1,
		C2:c.C2,
		Dims:c.Dims,
		Subck:c.Subck,
		Code:c.Code,
		Dtyp:c.Dtyp,
		Rtyp:c.Rtyp,
		
	}
	if cy.Cvrc == 0.0 || cy.Cvrt == 0.0{
		cy.Cvrc = MaxOf([]float64{c.Cvrc,c.Cvrt})
	}
	switch c.Styp{
	//holy shit check for deep copying structs
		case 1:
		var dims []float64
		if len(c.Dims) > 1{
			dims = []float64{c.Dims[1],c.Dims[0]}
		} else{
			dims = []float64{c.H, c.B}
		}
		cy.Dims = dims
		cy.SecInit()
		default:
		cys := kass.FlipX(*c.Sec)
		cy.Sec = &cys
		//cy.SecInit()
		//cy.Sec.UpdateProp()
		//fmt.Println("SECT","%v\n",cy.Sec)
	}
	//cy.Init()
	return
}

//ColSizeIs sizes a column for estimated axial load pu and percentage of steel pg
//if b == 0 no limiting dimension
//increase pu by 15% if subjected to bending moment - BUT WHY (it was in shah?)
//actually increase pu by factors based on endc/degree of column

func ColSizeIs(pu, fck, fy, pg, b float64, sectype int) (dims []float64){
	areq := pu*1e3/0.4/(fck + 1.67*fy*pg)
	switch sectype{
	//CHECK THESE STYPS COMPLETELY RANDOM
		case 0:
		//circle
		dia := 2.0 * math.Sqrt(areq/math.Pi)
		dia = math.Ceil(dia/10.0)*10.0
		dims = append(dims, dia, dia)
		case 1:
		//rectangle
		if b > 0 {
			h := math.Ceil(areq/b/10.0)*10.0
			dims = append(dims, []float64{b, h}...)
		} else {
			//square
			b = math.Ceil(math.Sqrt(areq)/10.0)*10.0
			dims = append(dims, []float64{b, b}...)
		}
		case 7,8,9,10:
		//ela, elb, elc, eld, eh eh eh
		if b == 0.0{
			b = 230.0
		}
		h := math.Ceil(areq + math.Pow(b,2)/2.0/b)
		dims = append(dims, []float64{b, h}...)
		case 18:
		//diamond 
		case 19:
		//pentagon
		
	}
	return
} 
