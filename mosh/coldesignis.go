package barf

import (
	"fmt"
	"log"
	"math"
	"errors"
	//"sort"
	kass"barf/kass"
	draw"barf/draw"
)

//shah col design funcs

//ColDTypIs was meant to classify columns by design type (axially loaded, uniaxial bending, biaxial bending, slender)
//delete this? unused
func ColDtypIs(c *RccCol,pu float64, nrows int) (error){
	//var dtype int
	var msfs, pu400 float64
	//this seems wrong?
	//calc axial load capacity of a 400x400 column with 0.8% steel
	if c.Fy ==250.0 {msfs = 0.77} else {msfs = 0.67}
	pu400 = 0.64 * c.Fck + (msfs * c.Fy - 0.4 * c.Fck) * 1.28
	if pu > pu400 {
		c.Dtyp = 0
		return nil
	}
	return nil
}

//ColFstlIs returns fsc, fst for column as per shah (is code)
//subtracts fck from fstl for steel displaced area NO IT DOESN'T

func ColFstlIs(xu, effd, fck, fy, esc, est float64) (fsc, fst float64) {
	fsc = RbrFrcIs(fy, esc) 
	if xu == effd{
		fst = 0.0
	} else if xu < effd {		
		fst = -RbrFrcIs(fy, est)
	} else {
		fst = RbrFrcIs(fy, est) 
	}
	return
}

//ColEqSqr returns the equivalent square column of a circular column of dia c.B
func ColEqSqr(c *RccCol) (csq *RccCol){
	ds := c.B - c.Cvrt - c.Cvrc
	cvrc := (0.8 * c.B - 2.0 * ds/3.0)/2.0
	csq = &RccCol{
		Fck:c.Fck,
		Fy:c.Fy,
		Cvrc:cvrc,
		Cvrt:cvrc,
		B:0.8 * c.B,
		H:0.8 * c.B,
		Lspan:c.Lspan,
		Lx:c.Lx,
		Ly:c.Ly,
		Leffx:c.Leffx,
		Leffy:c.Leffy,
		Nlayers:2,
		Type:"rectangle",
		Styp:1,
		Rtyp:c.Rtyp,
		Pu:c.Pu,
		Mux:c.Mux,
		Muy:c.Muy,
		Code:c.Code,
		C1:cvrc,
		C2:cvrc,
		Dims:[]float64{0.8 * c.B, 0.8 * c.B},
	}
	csq.Init()
	return
}

//ColDzIs designs a column section for a given pur, mux, muy
//only rect sections (mix of shah/hulse routines)
func ColDzIs(c *RccCol) (error){
	/*
	   REPLACE pu, mux, muy float64, axy bool, plot string, opt int
	   axy - design around both minor and major axis ()
	   opt - optimize 0 - do nuthin' 1 - GA, 2 - PSO
	   plotstr - draw dumb ansi
	*/
	pu := c.Pu; mux := c.Mux; muy := c.Muy
	//make dtyp start from 1 - NEVER USE 0 as a value
	switch {
	case mux + muy == 0.0:
		c.Dtyp = 0
	case mux > 0.0 && muy > 0.0:
		c.Dtyp = 2
	case mux > 0.0 && muy == 0.0:
		c.Dtyp = 1; c.Ax = "x"
	case muy > 0.0 && mux == 0.0:
		c.Dtyp = 1; c.Ax = "y"
		//c.Flip()
	}
	var pur, murx float64
	var err error
	switch c.Styp{
		case 0:
		//circle
		csq := ColEqSqr(c)
		//log.Println(csq.B, csq.H)
		switch c.Dtyp{
			case 0,1:
			pur, murx, err = ColDIs(csq, pu, mux)			
			if err != nil{
				return err
			}
			//log.Println("REZZZ",pur, murx, csq.Asteel)
			err = ColRbrGen(csq)
			if err != nil{
				//log.Println(ColorRed,"ERRORE,errore->",err,ColorReset)
				return err
			}
			c.Dias = csq.Dias
			c.Pur = pur; c.Murx = murx
			c.Asteel = csq.Asteel
			//c.Asc = c.Asteel/2.0; c.Ast = c.Asteel/2.0
			c.Nties = csq.Nties
			c.Ptie = csq.Ptie
			c.Dtie = csq.Dtie
			c.Csq = csq
			case 2:
			ColBxRectAst(csq,pu,mux,muy)
			
			case 3:
		}
		case 1:
		//rectangle
		switch c.Dtyp{
			case 0,1:
			pur, murx, err = ColDIs(c, pu, mux)
			if err != nil{
				return err
			}
			//log.Println(ColorGreen,"calc pur,mur",pur, murx,ColorReset)
			c.Pur = pur; c.Murx = murx
			err = ColRbrGen(c)
			
			if err != nil{
				//log.Println(ColorRed,"ERRORE,errore->",err,ColorReset)
				return err
			}
			//log.Println(ColorGreen,"calc pur,mur",pur, murx,ColorReset)
			//make this part of the report
			//pur, murx, err = ColAzIs(c, pu)
			//pur, murx, err = ColAzGen(c, pu)
			//if err != nil{
			//	return err
			//}
			//log.Println(ColorWhite,"analyzed pur, murx -",pur, murx,ColorReset)
			case 2:
			ColBxRectAst(c,pu,mux,muy)
			case 3:
			//slender biaxial bending huh? what is case 3 (?_?)
		}
		default:
		//wot? now use ColStl, etc	
	}
	//fmt.Println(ColorGreen,c.Ast, c.Asc, c.Asteel,c.Dias,ColorReset)
	return nil
}

//ColBxRectAst calcs ast (pu) for a rectangular column
//iterates from 0.8 percent asteel in increments of 0.1 percent to find optimal percent steel
func ColBxRectAst(c *RccCol, pu, mux, muy float64){
	var step, pnz, pg, mnx, mny, alpha, mrat float64
	var err error
	kiter := 0
	//iter := 0
	pg = 0.008; step = 0.0001
	pg -= step
	cy := ColFlip(c)
	for kiter < 2999{
		kiter++
		pg += step
		astl := pg * c.B * c.H
		pnz = 0.45 * c.Fck * c.B * c.H + 0.75 * c.Fy * astl
		pnz = pnz/1e3
		prat := pu/pnz
		switch{
		case prat <= 0.2:
			alpha = 1.0
		case prat >= 0.8:
			alpha = 2.0
		default:
			alpha = 2.0/3.0 + 5.0 * pu/3.0/pnz
		}
		mnx, err = ColRectMur(c, pu, astl)
		if err != nil{
			return
		}
		mny, err = ColRectMur(cy, pu, astl)
		if err != nil{
			return
		}
		mrat = math.Pow(mux/mnx, alpha) + math.Pow(muy/mny, alpha)
		//log.Println("step, pg, mrat ->",step, pg * 100.0, mrat)
		if mrat > 1.0{
			continue
		}
		if mrat <= 1.0{
			if math.Abs(mrat - 1.0) <= 0.1{
				log.Println("optimal percentage of steel found")
				log.Println(mrat, 1.0 - mrat)
				log.Println(astl)
				break
			} else{
				pg -= step
				step = step/100.0
				pg -= step
				//iter = 1
			}
		}
	}
	astr := pg * c.B * c.H
	//log.Println("astreq",astr)
	rez, _ := ColBxRbrGen(c, astr)
	cx := RccCol{
		Fck:c.Fck,
		Fy:c.Fy,
		Cvrc:c.Cvrc,
		Cvrt:c.Cvrt,
		B:c.B,
		H:c.H,
		Styp:c.Styp,
		Dtyp:c.Dtyp,
		Rtyp:c.Rtyp,
	}
	cy = ColFlip(&cx)
	var crez []RccCol
	var mcol int
	for _, r := range rez{
		dias, dxs, dys, barpts, aprov := ColRbrDbarGen(c, r)
		cx.Dias = dias; cx.Dbars = dxs; cx.Dybars = dys; cy.Dias = dias; cy.Dbars = dys
		cx.Barpts = barpts; cx.Asteel = aprov
		//pg := aprov/c.B/c.H
		pnz = 0.45 * c.Fck * c.B * c.H + 0.75 * c.Fy * aprov
		pnz = pnz/1e3
		prat := pu/pnz
		switch{
			case prat <= 0.2:
			alpha = 1.0
			case prat >= 0.8:
			alpha = 2.0
			default:
			alpha = 2.0/3.0 + 5.0 * pu/3.0/pnz
		}
		_, mnx, err = ColAzIs(&cx, pu)
		if err != nil{
			continue
		}
		_, mny, err = ColAzIs(cy, pu)
		if err != nil{
			continue
		}
		mrat = math.Pow(mux/mnx, alpha) + math.Pow(muy/mny, alpha)
		if mrat < 1.0{
			crez = append(crez, cx)
			mcol = len(crez) - 1
		}
	}
	cmin := crez[mcol]
	cmin.Printz()
	//DrawColRect(&cmin, "qt")
	log.Println(cmin.Dias)
	//GET ASTS HERE
}

//ColDBxIs (was meant to be the biaxial bendning design entry func) designs a column for biaxial bending
func ColDBxIs(c *RccCol, pu, mux, muy float64){
	/*
	   biaxial bending column design
	*/
	switch c.Styp{
		case 0:
		//circle
		case 1:
		//rectangle
		ColBxRectAst(c, pu, mux, muy)
		default:
		//wot?
	}
}

//ColDIs contains short column + uniaxial bending area of steel calc routines
//rips off hulse routines and uses is code rectangular stress block factors 
//ze backbone - only for sectype 0,1

func ColDIs(c *RccCol, pu, mux float64) (pur, murx float64, err error){
	var msfs, astot, k8, k9, mu, fst, fsc, est, esc, dck, ack, step, xu, mur, effd float64
	if c.Fy ==250.0 {msfs = 0.77} else {msfs = 0.67}
	if c.Code == 0{c.Code = 1}
	xstep := 10.0
	switch c.Code{
		case 1:
		k8 = 0.45; k9 = 0.8
		case 2:
		k8 = 0.45; k9 = 0.9
	}
	c.Init()
	switch c.Dtyp{
		case 0:
		//axially loaded short column
		pu = pu * 1e3
		switch c.Styp{
			case 0:
			//circular
			ack = math.Pi * math.Pow(c.B,2)/4.0
			case 1:
			//rectangular
			ack = c.B * c.H
			default:
			err = errors.New("invalid function section type")
			return
		}
		
		astot = (pu - 0.4 * c.Fck * ack)/(msfs*c.Fy - 0.4 * c.Fck)
		c.Asc = astot/2.0; c.Ast = c.Asc
		pur = pu
		murx = mux
		case 1:
		//axial load + uniaxial bending
		pu = pu * 1e3
		mu = mux* 1e6 //convert to N-mm
		var kiter int
		switch c.Styp{ 
			case 1:
			switch {
			case (c.Rtyp == 0 && c.Nlayers <= 2):
				//nlayers = 2
				effd = c.H - c.Cvrt
				dck = pu/c.B/k8/c.Fck
				if pu * (c.H - dck)/2.0 > mu{
					//get min. steel 
					c.Asc = 0.4 * c.B * c.H/2.0/100.0
					c.Ast = c.Asc
					err = nil
					return
				}
				xu = 2.33 * effd; step = xstep
				iter := 0; kiter = 0
				xu += c.H/step
				for iter != -1{
					kiter++; fsc = 0; fst = 0
					if kiter > 2999{
						log.Println("ERRORE, errore->max iterations reached")
						err = ErrD
						return 
					}
					switch iter{
					case 0:
						xu -= c.H/step
					case 1:
						xu += c.H/step
					}
					if xu <= 0.9 * effd && iter == 0{
						iter = 1
						xu = c.Cvrc + 1.0
						step = xstep
					}
					if xu >= c.H/k9 {
						dck = c.H
					} else {
						dck = k9 * xu
					}
					esc = 0.0035 * (xu - c.Cvrc)/xu
					est = math.Abs(0.0035 * (effd - xu)/xu)
					fsc, fst = ColFstlIs(xu, effd, c.Fck, c.Fy, esc, est)
					if iter == 0{
						astot = 2.0 * (pu - k8 * c.Fck * c.B * dck)/(fsc + fst)
					}
					if iter == 1{
						astot = 2.0 * (mu - k8 * c.Fck * c.B * dck * (c.H - dck)/2.0)/(fsc - fst)/(c.H/2.0 - c.Cvrc)
					}
					if astot < 0.0{
						continue
					}
					if iter == 0{
						mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * (fsc - fst) * (c.H/2.0 - c.Cvrc)/2.0
						if step >= xstep * 10.0{
							iter = -1
							break
						}
						if mur >= mu{
							if 1.02 * mu >= mur {
								iter = -1
								break
							} else {
								xu += c.H/step 
								step = step * 10.0
								xu += c.H/step
							}
						}
						
						if math.Abs(mur-mu) <= 0.02 * mu && step >= xstep*10.0{
							iter = -1
							break
						}
					}
					if iter == 1{
						pur = k8 * c.Fck * c.B * dck + astot * (fsc + fst)/2.0
						if (pu == 0 && pur - pu < 0.1) || (xu == c.Cvrc + 1.0 && pur > pu) {
							iter = -1
							break
						}
						if pur >= pu {	
							if step >= xstep * 10.0{
								iter = -1
								break
							}
							if 1.02 * pu >= pur {
								iter = -1
							} else {
								xu -= c.H/step
								step = step*10.0
								xu -= c.H/step
							}
						}
						
						if math.Abs(pur-pu) <= 0.02 * pu && step >= xstep*10.0{
							iter = -1
							break
						}
					}
				}
				pur = k8 * c.Fck * c.B * dck + astot * (fsc + fst) /2.0
				mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * (fsc - fst) * (c.H/2.0 -c.Cvrc)/2.0
				c.Asc = astot/2.0; c.Ast = c.Asc; c.Asteel = astot; c.Xu = xu
				pur = pur/1e3
				mur = mur/1e6
				murx = mur
				err = nil
				//fmt.Println(xu, c.Asc, c.Ast, astot, pur, mur, murx/1e6)
			case (c.Rtyp == 0 && c.Nlayers > 2):
				var dlevels []float64
				effd = c.H - c.Cvrt
				var fsum, zsum, ebar, fbar float64
				dlevels = append(dlevels,c.Cvrc)
				dlevels = append(dlevels, c.H - c.Cvrt)
				dstep := (c.H/2.0 - c.Cvrc)/float64(c.Nlayers/2)
				if (c.Nlayers - 2) % 2 != 0{
					dlevels = append(dlevels, c.H/2.0)
				}
				for i := 0; i < (c.Nlayers- 2)/2; i++{
					dlevels = append(dlevels, c.Cvrc + float64(i+1)*dstep)
					dlevels = append(dlevels, c.H/2.0 + float64(i+1)*dstep)
				}
				//log.Println("nlayers, dlevels->",c.Nlayers, dlevels)
				dck = pu/c.B/k8/c.Fck
				xu = 2.33 * effd; step = xstep
				iter := 0; kiter = 0
				xu += c.H/step
				for iter != -1{
					kiter++
					zsum, fsum = 0,0
					if kiter > 2999{
						log.Println("ERRORE, errore->max iterations reached")
						err = ErrD
						return
					}
					switch iter{
					case 0:
						xu -= c.H/step
					case 1:
						xu += c.H/step
					}
					if xu <= 0.9 * effd && iter == 0{
						iter = 1
						xu = c.Cvrc + 1.0
						step = xstep
					}
					
					if xu >= c.H/k9 {
						dck = c.H
					} else {
						dck = k9 * xu
					}		
					for _, dbar := range dlevels{
						ebar = math.Abs(0.0035 * (xu - dbar)/xu)
						switch {
						case xu >= dbar:
							fbar = RbrFrcIs(c.Fy, ebar)
							fsum += fbar
							zsum += (c.H/2.0 - dbar) * fbar
						case xu < dbar:
							fbar = -RbrFrcIs(c.Fy, ebar)
							fsum += fbar
							zsum += (c.H/2.0 - dbar) * fbar
						}
					}
					if iter == 0{
						astot = float64(c.Nlayers) * (pu - k8 * c.Fck * c.B * dck)/(fsum)
					}
					if iter == 1{
						astot = float64(c.Nlayers) * (mu - k8 * c.Fck * c.B * dck * (c.H - dck)/2.0)/zsum
					}
					if astot < 0.0{
						continue
					}
					if iter == 0{
						mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * zsum/float64(c.Nlayers)
						if mur >= mu{
							
							if step >= xstep * 10.0{
								iter = -1
								break
							}
							if 1.02 * mu >= mur {
								iter = -1
								break
							} else {
								xu += c.H/step 
								step = step * 10.0
								xu += c.H/step
							}
						}
						
						if math.Abs(mur-mu) <= 0.02 * mu && step >= xstep*10.0{
							iter = -1
							break
						}
					}
					if iter == 1{
						pur = k8 * c.Fck * c.B * dck + astot * fsum/float64(c.Nlayers)
						if (pu == 0 && pur - pu < 0.1) || (xu == c.Cvrc + 1.0 && pur > pu) {
							iter = -1
							break
						}
						if pur >= pu{	
							if step >= xstep * 10.0{
								iter = -1
								break
							}
							if 1.03 * pu >= pur {
								iter = -1
							} else {
								xu -= c.H/step
								step = step * 10.0
								xu -= c.H/step
							}
						}
						
						if math.Abs(pur-pu) <= 0.02 * pu && step >= xstep*10.0{
							iter = -1
							break
						}
					}
				}
				pur = k8 * c.Fck * c.B * dck + astot * fsum/float64(c.Nlayers)
				mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * zsum/float64(c.Nlayers)
				c.Asc = astot/2.0; c.Ast = c.Asc; c.Asteel = astot; c.Xu = xu
				pur = pur/1e3
				mur = mur/1e6
				murx = mur
				err = nil
			}
		}
	}
	
	//CHECK FOR THIS ALWAYZE (0.8 percent min steel)
	//astmin := 0.8 * c.B * c.H/100.0
	//if astot <  astmin{astot = astmin; c.Asteel = astot; c.Asc = astot/2.0; c.Ast = c.Asc}

	//why kind sir must one not generate rebar here 
	return 
}

//ColAzIs analyzes a column with multiple layers of rebar for a given axial load pu 
func ColAzIs(c *RccCol, pu float64) (float64, float64, error) {
	var k8, k9, ebar, fbar, diabar, pur, dck, xu, step, mur float64
	fsts := make([]float64, len(c.Dias))
	//is 456 stress block factors (subramanian)
	k8 = 0.45; k9 = 0.9
	pu = pu * 1e3
	//iterate to find neutral axis depth xu
	step = 100.0
	xu = c.Cvrc + 1.0
	//HERE c.Cvrc == depth of first bar
	xu = xu - c.H/step
	kiter := 0
	switch c.Styp{
		case 0:
		//circular section is sent in as equivalent square column
		//so do nuttin
		case 1:
		//rekt sekt
		for { //
			if kiter > 2999{
				log.Println("ERRORE, errore-> maximum iteration limit reached")
				return pur, mur, ErrCAxL
			}
			if pur < pu{
				xu += c.H/step
			} else {
				xu -= c.H/step
				step = 10.0 * step
				xu -= c.H/step
			}
			if xu > c.H/k9{
				dck = c.H
			} else {
				dck = k9 * xu
			}
			pur = k8 * c. Fck * c.B * dck
			for idx, dbar := range c.Dbars{
				diabar = c.Dias[idx]
				ebar = math.Abs(0.0035 * (dbar - xu)/xu)
				switch {
				case xu > dbar:
					fbar = RbrFrcIs(c.Fy, ebar)
				case xu == dbar:
					fbar = 0
				case xu < dbar:
					fbar = -RbrFrcIs(c.Fy, ebar)
				}
				fsts[idx] = fbar
				pur += RbrArea(diabar)*fbar
			}
			if xu == c.Cvrc + 1.0 && pur > pu{
				break
			}
			if pur > pu && pur <= 1.03 * pu{
				break
			}
			kiter++
		}
		mur = pur * c.H/2.0 - k8 * c.Fck * c.B * math.Pow(dck,2)/2.0
		for idx, fbar := range fsts {
			mur -= RbrArea(c.Dias[idx])*fbar*c.Dbars[idx]
		}
		case -1:
		//TODO generic section using *SectIn
		//generic
	}
	return pur/1e3, mur/1e6, nil
}

//ColNMIs calcs n-m interaction curves for a generic section with multiple layers of rebar
//column section is defined by c.Sec; dias and placement is defined by c.Dias and c.Dbars
func ColNMIs(c *RccCol, pureq float64, axx bool) (pus, mus []float64, mureq float64){
	//TODO - CHECK THIS AND ADD PLOT FUNC
	//n (p/axial load) - m (moment) interaction curves
	//ripping off hulse section 5.4
	//axx true - bending about x (major) axis
	var k8, k9, ebar, fbar, diabar, mu, pu, dck, puprev float64
	var data string
	k8 = 0.45; k9 = 0.9
	//ey = c.Fy/1.15/2e5
	fsts := make([]float64, len(c.Dbars))
	switch c.Styp {
	case 0:
		//circular sekt
	case 1:
		//rect sekt
		for xu := c.H/20.0; xu < 5.0 * c.H; xu += c.H/20.0 {
			if xu >= c.H/k9 {dck = c.H} else {dck = k9 * xu}
			puprev = pu
			pu = k8 * c. Fck * c.B * dck
			for idx, dbar := range c.Dbars {
				diabar = c.Dias[idx]
				ebar = math.Abs(0.0035 * (dbar - xu)/xu)
				if xu >= dbar {
					fbar = RbrFrcIs(fy, ebar)	
				} else {
					fbar = -RbrFrcIs(fbar, ebar)
				}
				fsts[idx] = fbar
				pu += RbrArea(diabar)*fbar
			}
			mu = pu * c.H/2.0 - k8 * c.Fck * c.B * math.Pow(dck,2)/2.0
			for idx, fbar := range fsts{
				mu -= RbrArea(c.Dias[idx])*fbar*c.Dbars[idx]
			}
			if pu <= 0{
				continue
			}
			if mu < 0.1 || math.Abs(pu - puprev) < 0.01{
				break
			}
			if pureq > 0.0 && math.Abs(pureq - pu) < 0.01{
				mureq = mu
			}
			pus = append(pus, pu)
			mus = append(mus, mu)
			data += fmt.Sprintf("%f %f\n",mu/1e4/100.0,pu/10.0/100.0)
		}
	default:
		//generic non rect sect
		//fmt.Println("NON REKT SEKT", c.Sec.ymx, c.Sec.ym)
		var ack, ycy, ygsec, yg, asec float64
		asec, _, ygsec, _, _, _, _, _, _ = kass.SecArea(c.Sec, false)
		yg = c.Sec.Ymx - ygsec
		if c.H == 0.0 {c.H = c.Sec.Ymx - c.Sec.Ym}
		for xu := c.H/20.0; xu < 5.0 * c.H; xu += c.H/20.0 {
			dck = k9 * xu
			if dck >= c.H {
				ack = asec
				ycy = ygsec
			} else {
				ack, _, ycy = ColSecArXu(c.Sec, dck)
			}
			ycy = c.Sec.Ymx - ycy
			puprev = pu
			pu = k8 * c.Fck * ack
			for idx, dbar := range c.Dbars {
				diabar = c.Dias[idx]
				ebar = math.Abs(0.0035 * (dbar - xu)/xu)
				if xu >= dbar {
					fbar = RbrFrcIs(fy, ebar)
				} else {
					fbar = -RbrFrcIs(fy, ebar)
				}
				fsts[idx] = fbar
				pu += RbrArea(diabar)*fbar
			}
			if pu <= 0 {
				continue
			}
			mu = pu * yg - k8 * c.Fck * ack * ycy
			for idx, fbar := range fsts {
				mu -= RbrArea(c.Dias[idx])*fbar*c.Dbars[idx]
			}
			if mu < 0.1 || math.Abs(pu - puprev) < 0.01 {
				break
			}
			pus = append(pus, pu)
			mus = append(mus, mu)
			data += fmt.Sprintf("%f %f\n",mu/1e4/100.0,pu/10.0/100.0)
		}		
	}
	skript := "d2.gp"
	term := c.Term
	folder := ""
	if c.Web{folder="web"}
	switch term{
		case "svg","svgmono":
		default:
		term = "svg"
	}
	title := c.Title + "-M-N curve.svg"
	fname := genfname("",title)
	//Draw(data, skript, term, folder, fname, title, xl, yl, zl string) (txtplot string, err error)
	dstr, _ := draw.Draw(data, skript, term, folder, fname, "M-N interaction curve", "BM(kn.m)","N(kn)","")
	c.Nmplot = dstr
	
	c.Txtplots = append(c.Txtplots, dstr)

	return
}

//ColRectMur returns the ultimate moment of resistance of a column section given pu/ast
//wh-aat, what is dis good for - BIAXE
func ColRectMur(c *RccCol, pu, astl float64) (float64, error) {
	var k8, k9, ebar, fbar, pur, dck, xu, step, mur float64
	//is 456 stress block factors (subramanian) IS 0.8x deep FOOL
	k8 = 0.45; k9 = 0.9
	pu = pu * 1e3
	//get dxlevels, alevelx
	dxs := []float64{c.Cvrc, c.H - c.Cvrc}
	alevs := []float64{astl/2.0, astl/2.0}
	fsts := make([]float64, 2)
	step = 20.0
	xu = c.Cvrc + 1.0
	xu = xu - c.H/step
	kiter := 0
	for { //
		if kiter > 2999{
			//if c.Verbose{//log.Println("ERRORE, errore-> maximum iteration limit reached")}
			return mur, ErrCAxL
		}
		if pur < pu{
			xu += c.H/step
		} else {
			xu -= c.H/step
			step = 20.0 * step
			xu -= c.H/step
		}
		if xu > c.H/k9{
			dck = c.H
		} else {
			dck = k9 * xu
		}
		pur = k8 * c. Fck * c.B * dck
		for idx, dx := range dxs{
			ebar = math.Abs(0.0035 * (dx - xu)/xu)
			switch {
			case xu > dx:
				fbar = RbrFrcIs(c.Fy, ebar)
			case xu == dx:
				fbar = 0
			case xu < dx:
				fbar = -RbrFrcIs(c.Fy, ebar)
			}
			fsts[idx] = fbar
			pur += alevs[idx]*fbar
		}
		if xu == c.Cvrc + 1.0 && pur > pu{
			break
		}
		if pur > pu && pur <= 1.03 * pu{
			break
		}
		kiter++
	}
	mur = pur * c.H/2.0 - k8 * c.Fck * c.B * math.Pow(dck,2)/2.0
	for idx, fbar := range fsts {
		mur -= alevs[idx] * fbar
	}
	return mur/1e6, nil
}


/*
func ColDIs(c *RccCol, pu, mux float64) (pur, murx float64, err error){
	//short column + uniaxial bending area of steel
	//ze backbone - only for sectype 0,1
	var msfs, astot, k8, k9, mu, fst, fsc, est, esc, dck, ack, step, xu, mur, effd float64
	if c.Fy ==250.0 {msfs = 0.77} else {msfs = 0.67}
	pu = pu * 1e3
	if c.Code == 0{c.Code = 1}
	switch c.Code{
		case 1:
		k8 = 0.45; k9 = 0.8
		case 2:
		k8 = 0.45; k9 = 0.9
	}
	c.Init()
	switch c.Dtyp{
		case 0:
		//axially loaded short column
		switch c.Styp{
			case 0:
			//circular
			ack = math.Pi * math.Pow(c.B,2)/4.0
			case 1:
			//rectangular
			ack = c.B * c.H
			default:
			err = errors.New("invalid function section type")
			return
		}
		astot = (pu - 0.4 * c.Fck * ack)/(msfs*c.Fy - 0.4 * c.Fck)
		c.Asc = astot/2.0; c.Ast = c.Asc
		pur = pu
		murx = mux
		case 1:
		//axial load + uniaxial bending
		mu = mux* 1e6 //convert to N-mm
		var kiter int
		switch c.Styp{ 
			case 1:
			switch {
			case (c.Rtyp == 0 && c.Nlayers == 0):
				//nlayers = 2
				effd = c.H - c.Cvrt
				dck = pu/c.B/k8/c.Fck
				if pu * (c.H - dck)/2.0 > mu{
					//get min. steel 
					c.Asc = 0.4 * c.B * c.H/2.0/100.0
					c.Ast = c.Asc
					err = nil
					return
				}
				xu = 2.33 * effd; step = 100.0
				iter := 0; kiter = 0
				xu += c.H/step
				for iter != -1{
					kiter++; fsc = 0; fst = 0
					if kiter > 2999{
						//if c.Verbose{log.Println("ERRORE, errore->max iterations reached")}
						err = ErrIter
						return 
					}
					switch iter{
					case 0:
						xu -= c.H/step
					case 1:
						xu += c.H/step
					}
					if xu <= 0.9 * effd && iter == 0{
						iter = 1
						xu = c.Cvrc + 1.0
						step = 100.0
					}
					if xu > c.H/k9 {
						dck = c.H
					} else {
						dck = k9 * xu
					}
					esc = 0.0035 * (xu - c.Cvrc)/xu
					est = math.Abs(0.0035 * (effd - xu)/xu)
					fsc, fst = ColFstlIs(xu, effd, c.Fck, c.Fy, esc, est)
					if iter == 0{
						astot = 2.0 * (pu - k8 * c.Fck * c.B * dck)/(fsc + fst)
					}
					if iter == 1{
						astot = 2.0 * (mu - k8 * c.Fck * c.B * dck * (c.H - dck)/2.0)/(fsc - fst)/(c.H/2.0 - c.Cvrc)
					}
					if astot < 0.0{
						continue
					}
					if iter == 0{
						mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * (fsc - fst) * (c.H/2.0 - c.Cvrc)/2.0
						//if math.Abs(mur - mu) <= 0.03 * mu{
						//	iter = -1
						//	break
						//}
						if mur >= mu {
							if 1.02 * mu >= mur{
								iter = -1
								break
							} else {
								xu += c.H/step 
								step = step 
								xu += c.H/step
							}
						}
					}
					if iter == 1{
						pur = k8 * c.Fck * c.B * dck + astot * (fsc + fst)/2.0
						if (pu == 0 && pur - pu < 0.1) || (xu == c.Cvrc + 1.0 && pur > pu) {
							iter = -1
							break
						}
						if pur >= pu {
							if 1.02 * pu >= pur{
								iter = -1
							} else if step > 1e12{
								iter = -1
								break
							}else {
								xu -= c.H/step
								step = 1000.0
								xu -= c.H/step
							}
						}	
						//if math.Abs(pur - pu) <= 0.03 * pu{
						//	iter = -1
						//	break
						//}
					}
				}
				pur = k8 * c.Fck * c.B * dck + astot * (fsc + fst) /2.0
				mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * (fsc - fst) * (c.H/2.0 -c.Cvrc)/2.0
				c.Asc = astot/2.0; c.Ast = c.Asc; c.Asteel = astot; c.Xu = xu
				pur = pur/1e3
				mur = mur/1e6
				murx = mur
				err = nil
			case (c.Rtyp == 0 && c.Nlayers > 0):
				var dlevels []float64
				effd = c.H - c.Cvrt
				var fsum, zsum, ebar, fbar float64
				dlevels = append(dlevels,c.Cvrc)
				dlevels = append(dlevels, c.H - c.Cvrt)
				dstep := (c.H/2.0 - c.Cvrc)/float64(c.Nlayers/2)
				if (c.Nlayers - 2) % 2 != 0{
					dlevels = append(dlevels, c.H/2.0)
				}
				for i := 0; i < (c.Nlayers- 2)/2; i++{
					dlevels = append(dlevels, c.Cvrc + float64(i+1)*dstep)
					dlevels = append(dlevels, c.H/2.0 + float64(i+1)*dstep)
				}
				
				dck = pu/c.B/k8/c.Fck
				xu = 2.33 * effd; step = 100.0
				iter := 0; kiter = 0
				xu += c.H/step
				for iter != -1{
					kiter++
					zsum, fsum = 0,0
					if kiter > 2999{
						//log.Println("ERRORE, errore->max iterations reached")
						err = ErrIter
						return
					}
					switch iter{
					case 0:
						xu -= c.H/step
					case 1:
						xu += c.H/step
					}
					if xu <= 0.9 * effd && iter == 0{
						iter = 1
						xu = c.Cvrc + 1.0
						step = 100.0
					}
					if xu >= c.H/k9 {
						dck = c.H
					} else {
						dck = k9 * xu
					}		
					for _, dbar := range dlevels{
						ebar = math.Abs(0.0035 * (xu - dbar)/xu)
						switch {
						case xu >= dbar:
							fbar = RbrFrcIs(c.Fy, ebar)
							fsum += fbar
							zsum += (c.H/2.0 - dbar) * fbar
						case xu < dbar:
							fbar = -RbrFrcIs(c.Fy, ebar)
							fsum += fbar
							zsum += (c.H/2.0 - dbar) * fbar
						}
					}
					if iter == 0{
						astot = float64(c.Nlayers) * (pu - k8 * c.Fck * c.B * dck)/(fsum)
					}
					if iter == 1{
						astot = float64(c.Nlayers) * (mu - k8 * c.Fck * c.B * dck * (c.H - dck)/2.0)/zsum
					}
					if astot < 0.0{
						continue
					}
					if iter == 0{
						mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * zsum/float64(c.Nlayers)
						if mur >= mu{
							if 1.02 * mu >= mur {
								iter = -1
								break
							} else {
								xu += c.H/step 
								step = step * 10.0
								xu += c.H/step
							}
						}
						//if math.Abs(mur-mu) <= 0.03 * mu{
						//	iter = -1
						//	break
						//}
					}
					if iter == 1{
						pur = k8 * c.Fck * c.B * dck + astot * fsum/float64(c.Nlayers)
						if (pu == 0 && pur - pu < 0.1) || (xu == c.Cvrc + 1.0 && pur > pu) {
							iter = -1
							break
						}	
						//if math.Abs(pur - pu) <= 0.03 * pu{
						//	iter = -1
						//	break
						//}
						if pur >= pu{
							if 1.02 * pu >= pur {
								iter = -1
							}  else {
								xu -= c.H/step
								step = step * 10.0
								xu -= c.H/step
							}
						}
					}
				}
				pur = k8 * c.Fck * c.B * dck + astot * fsum/float64(c.Nlayers)
				mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * zsum/float64(c.Nlayers)
				c.Asc = astot/2.0; c.Ast = c.Asc; c.Asteel = astot; c.Xu = xu
				pur = pur/1e3
				mur = mur/1e6
				murx = mur/1e6
				err = nil
				fmt.Println("ASTOT->",astot)
			}
		}
	}
	//CHECK FOR THIS ALWAYZE (0.8 percent min steel)
	//astmin := 0.8 * c.B * c.H/100.0
	//if astot <  astmin{astot = astmin; c.Asteel = astot; c.Asc = astot/2.0; c.Ast = c.Asc}

	//why kind sir must one not generate rebar here 
	return 
}

*/

