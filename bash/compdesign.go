package barf

import (
	"fmt"
	"math"
	kass"barf/kass"
	draw"barf/draw"
)

//CmpBm is a composite beam/slab struct
type CmpBm struct{
	Lspan  float64 //steve holt!
	Fck    float64
	Fy     float64
	Fyr    float64 //fy for rebar
	DL     float64 //slab DL
	LL     float64 //slab LL
	Pck    float64 
	PSFs   []float64
	Lslb   float64
	Dslb   float64 //total depth of slab
	Drib   float64 //depth of rib
	Hc     float64 //hc = dslb - drib
	Wuc    float64 //
	Wu     float64
	Muc    float64
	Mu     float64
	Vuc    float64
	Vu     float64
	Beff   float64 //breadth of flange
	Lrfac  float64 //lat. restraints at each 1/lrfac of span
	Qd     float64 //design strength of shear stud in kn
	Dstd   float64 //dia of shear stud
	Lstd   float64 //length/height of shear stud
	Sstd   float64 //spacing of shear studs
	Nstd   float64 //total number of shear studs
	Nrow   float64 //nstuds per row
	Dlk    float64 //const. stage DL (unfactored)
	Llk    float64 //const. stage LL (unfactored)
	Dls    float64 //service load DL (unfactored)
	Lls    float64 //service load LL (unfactored)
	Isdeck bool    //has profiled deck
	Ap     float64
 	Ip     float64
	Tp     float64
	Mpa    float64
	Ep     float64
	Dp     float64
	Chk    bool
	Dz     bool 
	Verbose bool
	Sdxs   []int
	Ssecs  []kass.StlSec
	Dtyp   int 
	Nspans int
	Nsecs  int
	Sdx    int
	Code   int
	Sname  string
}

//Init inits a CmpBm
func (c *CmpBm) Init(){
	if c.Code == 0{c.Code = 1}
	if c.Pck == 0{c.Pck = 24.0}
	if len(c.PSFs) < 2{c.PSFs = []float64{1.35,1.5}}
	if c.Nspans == 0{c.Nspans = 1}
	if c.Lstd == 0.0{c.Lstd = 100.0}
	if c.Dstd == 0.0{c.Dstd = 22.0}
	if c.Nrow == 0.0{c.Nrow = 2.0}
	if c.Lrfac == 0.0 || !c.Chk{c.Lrfac = 1.0}
	return
}

//Draw draws a CmpBm
func (c *CmpBm) Draw()(err error){
	fmt.Println("starting draw, ssec len-",len(c.Ssecs))
	switch len(c.Ssecs){
		case 0:
		err = fmt.Errorf("no sections found (arr len - %v)",len(c.Ssecs))
		default:
		ss := c.Ssecs[0]
		ss.Draw()
		x0 := - c.Beff/2.0
		x1 := c.Beff/2.0
		y0 := c.Ssecs[0].H
		y1 := y0 + c.Dslb
		data := fmt.Sprintf("%f %f %f\n",x0,y0,1.0)
		data = fmt.Sprintf("%f %f %f\n",x1,y1,1.0)
		data += "\n\n"
		cds := [][]float64{{x0,y0},{x1,y0},{x1,y1},{x0,y1},{x0,y0}}
		for i, p1 := range cds[:len(cds)-1]{
			p2 := cds[i+1]
			data += fmt.Sprintf("%f %f %f %f %f %.f\n",p1[0],p1[1],p2[0]-p1[0], p2[1]-p1[1],1.0, 1.0)
		}
		for i, p1 := range ss.Sec.Coords[:len(ss.Sec.Coords)-1]{
			p2 := ss.Sec.Coords[i+1]
			data += fmt.Sprintf("%f %f %f %f %f %.f\n",p1[0],p1[1],p2[0]-p1[0], p2[1]-p1[1],2.0, 2.0)
		}
		//draw shear studs
		//fmt.Println("no. of shear studs in a row-",c.Nrow, "len, dia of studs",c.Lstd,c.Dstd)
		x0 = ss.B/3.0
		x1 = 2.0*ss.B/3.0
		y0 = ss.H
		data += fmt.Sprintf("%f %f %f %f %f %.f\n",x0,y0,0.0,c.Lstd,3.0,3.0)
		data += fmt.Sprintf("%f %f %f %f %f %.f\n",x0-c.Dstd/2.0,y0+c.Lstd,c.Dstd,0.0,3.0,3.0)
		data += fmt.Sprintf("%f %f %f %f %f %.f\n",x1,y0,0.0,c.Lstd,3.0,3.0)
		data += fmt.Sprintf("%f %f %f %f %f %.f\n",x1-c.Dstd/2.0,y0+c.Lstd,c.Dstd,0.0,3.0,3.0)
		data += "\n\n"
		//add labels
		//slab depth
		
		data += fmt.Sprintf("%f %f %.1fmm\n",c.Beff/2.0+c.Dslb/2.0,y0+c.Dslb/2.0,c.Dslb)
		//section name
		data += fmt.Sprintf("%f %f %s\n",ss.B+c.Dslb/2.0,ss.H/2.0,ss.Sstr)
		//shear stud spacing
		data += fmt.Sprintf("%f %f %.1fmm\n",c.Beff/2.0+c.Dslb/2.0,y0+c.Dslb/2.0,c.Dslb)
		var txtplot string
		//data, skript, term, folder, fname, title, xl, yl, zl
		//txtplot, err = draw.Dumb(data, "basic2d.gp", "dumb", "compbm","mm","mm","mm")
		txtplot, err = draw.Draw(data, "basic2d.gp", "qt", "","cmpbm","compbm","mm","mm","mm")
		if err != nil{return}
		fmt.Println(txtplot)
	}
	return
}
//ShrConQd calcs the ultimate shear strength per stud
func ShrConQd(ds, hs, fck float64)(qd float64){
	//these are old vals from 1985 (for ex.)
	switch{
		case ds == 22.0, hs == 100.0:
		switch fck{
			case 20:
			qd = 49.0
			case 30:
			qd = 58.0
			case 35:
			case 40:
		}
	}
	return
}

//CmpBmDz designs a section as a composite beam
func CmpBmDz(c *CmpBm)(err error){
	var ss kass.StlSec
	ndx := kass.StlSdxLims[c.Sname]
	if ndx == 0{
		err = fmt.Errorf("%s design functions not written",c.Sname)
		return
	}
	if c.Chk{
		ss, err = c.ChkSec(c.Sdx)
		if err != nil{
			c.Ssecs = append(c.Ssecs, ss)
		}
		return
	}
	if c.Sdx > 0{ndx = c.Sdx}
	for idx := ndx; idx >= 0; idx--{
		if len(c.Sdxs) == c.Nsecs{
			break
		}
		if c.Verbose{fmt.Println(ColorGreen,"checking ndx-", idx,ColorReset)}
		ss, err = c.ChkSec(idx)
		if err == nil{
			c.Sdxs = append(c.Sdxs, idx)
			c.Ssecs = append(c.Ssecs, ss)
			fmt.Println(ColorYellow,"section found->",ss.Sstr,ss.Wt,ColorReset)
		}else{
			fmt.Println(ColorRed, err, ColorReset)
		}
	}
	return
}

//ChkSec checks a section and slab depth as a composite beam
func (c *CmpBm) ChkSec(sdx int)(ss kass.StlSec, err error){
	c.Init()
	ss, err = kass.GetStlSec(c.Sname, sdx, c.Code)
	if err != nil{return}
	dlfac := c.PSFs[0]
	llfac := c.PSFs[1]	
	slbw := c.Dslb*c.Lslb*c.Pck/1e6
	dlc := c.Dlk * c.Lslb/1e3
	llc := c.Llk * c.Lslb/1e3
	dlc = dlfac * (dlc + ss.Wt/1e3 + slbw)
	llc = (c.Llk * c.Lslb/1e3) * llfac
	if c.Verbose{fmt.Printf("const. stage loads DL %.1f kn/m LL %.1f kn/m\n",dlc,llc)}
	c.Wuc = dlc + llc
	if c.Verbose{fmt.Printf("ult udl const. %.1f kn/m\n",c.Wuc)}
	dlc = dlfac * (slbw + c.DL * c.Lslb/1e3 + ss.Wt/1e3)
	llc = llfac * (c.LL * c.Lslb)/1e3
	//serivce load vars
	dls := (slbw + c.DL * c.Lslb/1e3 + ss.Wt/1e3)
	lls := (c.LL * c.Lslb)/1e3
	if c.Verbose{fmt.Printf("comp. stage loads DL %.1f kn/m LL %.1f kn/m\n",dlc,llc)}
	c.Wu = dlc + llc
	if c.Verbose{fmt.Printf("ult udl comp. %.1f kn/m\n",c.Wu)}
	//calc Dz vals
	c.Muc = c.Wuc * math.Pow(c.Lspan/1e3,2)/8.0
	c.Mu  = c.Wu * math.Pow(c.Lspan/1e3,2)/8.0
	if c.Nspans > 1{
		//TODO
		//continuous composite beam
	}
	fmt.Printf("ult b.m at const %.1f knm at comp %.1f knm\n",c.Muc,c.Mu)
	//mur at construction stage
	iter := 0
	for iter != -1{
		ss.Lspan = c.Lspan/c.Lrfac
		err = ss.CalcMur()
		if err != nil{
			iter = -1
			fmt.Println("ERRORE in stl sec-",err)
			break
			return
		}
		if ss.Mur/1e6 > c.Muc{
			iter = -1
			break
		} else {
			c.Lrfac += 1.0
		}
		if c.Lrfac > 6.0{
			fmt.Println("ERRORE in LRfac-",c.Lrfac)
			iter = -1
			break
		}
	}
	
	if c.Verbose{fmt.Printf("req. mu at const. %.1f section max. %.1f ok? %t\n",c.Muc, ss.Mur/1e6, c.Muc < ss.Mur/1e6)}
	if c.Verbose{fmt.Println("at lrfac",c.Lrfac)}
	if c.Lrfac > 6.0{
		err = fmt.Errorf("lat. restraint spacing error lrfac %.1f spacing %.1f mm",c.Lrfac,c.Lspan/c.Lrfac)
		return
	}
	//breadth of flange
	beff := c.Lslb
	if c.Lspan/4.0 < beff{
		beff = c.Lspan/4.0
	}
	c.Beff = beff
	lam := 0.8
	//nk := 1.0
	ak := 0.87 * c.Fy/(0.36 * c.Fck)
	Af := ss.Tf * ss.B
	dct := (c.Dslb + ss.H)/2.0
	//design moment capacity md
	murd := 0.0
	xu := 0.0
	if c.Verbose{fmt.Println("ak",ak,"area",ss.Area,"beff",beff)}
	if c.Verbose{fmt.Println("aA",ak*ss.Area,"beff",beff,"af",Af,"Dslb",c.Dslb,"dct",dct,"ymo",ss.Ymo,"fy",ss.Fy)}
	//from IS 11384-2022
	c1 := beff * c.Dslb > ak * ss.Area
	c3 := beff * c.Dslb + 2.0 * ak * Af < ak * ss.Area
	c2 := !c1 && !c3
	switch{
		case c1:
		if c.Verbose{fmt.Println("pna within slab")}
		xu = ak * ss.Area/beff
		murd = ss.Area * ss.Fy*(dct + 0.5 * c.Dslb - lam * xu/2.0)*0.87
		case c2:
		//CHECK THIS
		if c.Verbose{fmt.Println("pna in steel flange")}
		xu = c.Dslb + (ak * ss.Area - beff * c.Dslb)
		murd = 0.87 * ss.Fy * (ss.Area * (dct + 0.5 * c.Dslb*(1.0-lam)) - ss.B * (xu - c.Dslb)*(xu + (1.0 - lam)*c.Dslb))
		case c3:
		//AND CHECK THIS
		if c.Verbose{fmt.Println("pna in web")}
		xu = c.Dslb + ss.Tf + (ak * (ss.Area - 2.0 * Af) - beff * c.Dslb)/(2.0 * ak * ss.Tw)
		murd = 0.87 * ss.Fy * (ss.Area * (dct + 0.5 * c.Dslb * (1.0 -lam)) - 2.0 * Af * (0.5 * ss.Tf + (1.0 - lam/2.0)*c.Dslb) - ss.Tw *(xu - c.Dslb - ss.Tf)*(xu + (1.0 - lam) * c.Dslb + ss.Tf))
	}
	if c.Verbose{
		fmt.Println("xu",xu,"mm")
		fmt.Println("ult. moment capacity murd",murd/1e6,"knm")
	}
	if c.Qd == 0.0{
		c.Qd = ShrConQd(c.Dstd, c.Lstd, c.Fck)
	}
	if c.Verbose{fmt.Println("shear strength per stud",c.Qd,"kn")}
	fcc := 0.36 * c.Fck * beff * xu/1e3
	c.Nstd = math.Ceil(fcc/c.Qd) * 2.0
	c.Sstd = math.Floor(c.Lspan/c.Nstd/5.0)*5.0*c.Nrow
	if c.Verbose{fmt.Println("number of studs",c.Nstd,c.Nrow,"per row","spacing",c.Sstd,"mm")}
	//copy vals to ss
	ss.Cmp = true; ss.Lstd = c.Nstd; ss.Sstd = c.Sstd; ss.Dstd = c.Dstd
	//serviceability checks
	//get depth of neutral axis (live load)
	//modular ratios for dead and live load
	mdl := 30.0
	mll := 15.0
	dg := c.Dslb + ss.H/2.0
	//ixx of transformed section
	ixt := 0.0
	c1 = ss.Area * (dg - c.Dslb) < 0.5 * beff * math.Pow(c.Dslb,2.0)/mll
	if c1{
		//neutral axis lies in slab
		if c.Verbose{fmt.Println("neutral axis lies in slab")}
		a := beff*math.Pow(c.Dslb,2.0)*0.5/mll; b := ss.Area; c := - ss.Area * dg
		xu = kass.SolveQuad(a,b,c)
		ixt = ss.Ixx + ss.Area * math.Pow(dg - xu, 2.0) + math.Pow(xu, 3.0) * beff/3.0/mll
	} else {
		//neutral axis in beam
		if c.Verbose{fmt.Println("neutral axis lies in steel sec")}
		xu = (beff * math.Pow(c.Dslb,2) * 0.5/mll + ss.Area * dg)/(c.Dslb * beff/mll + ss.Area)
		ixt = ss.Ixx + ss.Area * math.Pow(dg - xu, 2) + beff * c.Dslb * (math.Pow(c.Dslb,2)/12.0 + math.Pow(xu-c.Dslb/2.0,2))/mll
	}
	if c.Verbose{
		fmt.Println("depth of neutral axis",xu,"mm")
		fmt.Println("moment of inertia of the transformed section",ixt,"mm4")
		fmt.Println("moment of inertia of steel section",ss.Ixx,"mm4")
		fmt.Printf("service loads dead %.2f kn/m live %.1f kn/m\n",dls,lls)
	}
	//deflection due to dead load
	ddl := 5.0 * dls * math.Pow(c.Lspan,4)/384.0/ss.Em/ss.Ixx 
	if c.Verbose{fmt.Printf("deflection due to dead load %.1f mm\n",ddl)}
	dll := 5.0 * lls * math.Pow(c.Lspan, 4.0)/384/ss.Em/ixt
	if c.Verbose{fmt.Printf("deflection due to live load %.1f mm\n",dll)}
	if ddl + dll > c.Lspan/325.0{
		if c.Verbose{fmt.Println("section fails deflection check (span/325)")}
		err = fmt.Errorf("section fails deflection check req. ratio - 325 vs actual %.1f",c.Lspan/(ddl+dll))
		return
	}
	if c.Verbose{fmt.Println("calculating stresses at service due to dead and live loads")}
	//calc stresses due to ll
	mull := lls * math.Pow(c.Lspan, 2.0)/8.0
	zll := ss.H + c.Dslb - xu
	sigll := mull * zll/ixt
	//calc stresses due to dl
	c1 = ss.Area * (dg - c.Dslb) < 0.5 * beff * math.Pow(c.Dslb,2.0)/mdl
	if c1{
		//neutral axis lies in slab
		if c.Verbose{fmt.Println("neutral axis lies in slab for DL")}
		a := beff*math.Pow(c.Dslb,2.0)*0.5/mdl; b := ss.Area; c := - ss.Area * dg
		xu = kass.SolveQuad(a,b,c)
		ixt = ss.Ixx + ss.Area * math.Pow(dg - xu, 2.0) + math.Pow(xu, 3.0) * beff/3.0/mdl
	} else {
		//neutral axis in beam
		if c.Verbose{fmt.Println("neutral axis lies in steel sec for DL")}
		xu = (beff * math.Pow(c.Dslb,2) * 0.5/mdl + ss.Area * dg)/(c.Dslb * beff/mdl + ss.Area)
		ixt = ss.Ixx + ss.Area * math.Pow(dg - xu, 2) + beff * c.Dslb * (math.Pow(c.Dslb,2)/12.0 + math.Pow(xu-c.Dslb/2.0,2))/mdl
	}
	if c.Verbose{fmt.Println("depth of neutral axis for DL",xu,"mm")}
	if c.Verbose{fmt.Println("moment of inertia of the transformed section for DL",ixt,"mm4")}
	mudl := dls * math.Pow(c.Lspan,2.0)/8.0
	zdl := ss.H + c.Dslb - xu
	sigdl := mudl * zdl/ixt 
	if c.Verbose{fmt.Printf("max stresses in steel flange dl %.1f n/mm2 ll %.1f n/mm2 net %.1f n/mm2\n",sigdl,sigll,sigdl+sigll)}
	return
}
