package barf

import (
	"fmt"
	"log"
	"math"
	"errors"
	kass"barf/kass"
)

//ColDzIter iterates over ColDesign and finds min ast
//varies the number of layers for a rect/circular section
//varies the rebar level (side split of link sec) for other sections
//TODO - unfinished/it is brake
//i guess this is in cold now?
func ColDzIter(c *RccCol) (err error){
	if c.Ignore{return}
	switch c.Rtyp{
		case 0:
		case 1:
		case 2:
	}
	var iter int
	//var astmin float64
	//var minval int
	switch c.Styp{
		case 0 , 1:		
		switch c.Dtyp{
			case 0:
			case 1:
		}
		for iter == 0 && c.Nlayers <= 6{
			err = ColDesign(c)
			if err == nil{
				iter = 1
				break
				//ast := c.Asteel
			}
			
			c.Nlayers += 1
		}
		default:
		for iter == 0 && c.Rbrlvl < 2{
			err = ColDesign(c)
			if err == nil{
				iter = 1
				break
			}
			c.Rbrlvl++
		}		
	}
	return
}

//ColDesign is the general (-_-)7 col design entry func
//from flags/menu
//check if analysis or design - DO THIS IN INPUT func
//calc eff height if there be data
//check if rectangular or else
func ColDesign(c *RccCol) (err error){
	if c.Ignore{return}
	c.Init()
	//fmt.Println(ColorCyan,c.Title,"dtyp",c.Dtyp,c.Styp,c.Dims,c.Code,c.Pu,c.Mux,c.Muy,ColorReset)
	switch c.Styp{
	//get eq sqr col in Dz func
		case 0:
		csq := ColEqSqr(c)
		_, _, err = ColStl(csq, csq.Pu, csq.Mux, csq.Muy)
		if err != nil{return}
		err = csq.BarGen()
		if err != nil{return}
		c.Dias = csq.Dias
		c.Asteel = csq.Asteel
		c.Nties = csq.Nties
		c.Ptie = csq.Ptie
		c.Dtie = csq.Dtie
		c.Csq = csq
		c.Rbropt = csq.Rbropt
		c.Rbrt = csq.Rbrt; c.Rbrc = csq.Rbrc
		c.Rbrtopt = csq.Rbrtopt
		c.Rbrcopt = csq.Rbrcopt
		c.Psteel = c.Asteel*100.0/c.Sec.Prop.Area
		
		default:
		_, _, err = ColStl(c, c.Pu, c.Mux, c.Muy)
		if err != nil{
			return
		}
		if c.Styp == 1{
			if c.Rtyp == 0 && c.Nlayers > 2{
				err = ColRbrGen(c)
			} else {
				err = c.BarGen()
			}
		} else {
			err = c.BarGen()
		}
	}
	//fmt.Println("title",c.Title)
	//fmt.Println(ColorCyan,"dias")
	//fmt.Println(c.Dias)
	//fmt.Println(ColorRed,"dbars")
	//fmt.Println(c.Dbars)
	//fmt.Println(ColorReset)
	if err != nil{
		//fmt.Println(ColorRed,err,ColorReset)
		return
	}
	// err = c.Plot(c.Term)
	// if err != nil{
	// 	return
	// }
	c.Quant()
	//c.Printz()
	
	if c.Verbose{
		c.Table(false)
	}
	//check if code = 1 (is456) or 2(bs8110)
	//do shit
	//print report
	//draw
	//return - DONE B-)
	return
}

//Plot plots a column section with rebar
func (c *RccCol) Plot(term string) (err error){
	if term != ""{
		switch {
		case c.Styp == 0:
			_ = DrawColCircle(c, term, false)
			//case c.Styp == 1 && c.Nlayers > 2:
			//err := DrawColRect(c, term)	
			//if err != nil{
			//	log.Println(ColorRed,"ERRORE, errore-> plot error", err,ColorReset)
			//}
		default:
			_ = PlotColGeom(c, term, false)
			//if term == "dumb" || c.Term == "mono"{
			//	fmt.Println(pltstr)
			//}
		}
	} 
	return
}

//ColAnalyze is the entry func from flags/menu for column analysis (b`_`)b
func ColAnalyze(c *RccCol) (err error){
	var pur, mur float64
	switch{
		case len(c.Dbars) == 0:
		switch c.Code{
			case 1:
			pur, mur, err = ColAzIs(c, c.Pu)
			
			_, _, _ = ColNMIs(c, 0.0, true)
			case 2:
			mur, err = ColAzBs(c, c.Pu)
			pur = c.Pu
			_, _ = ColNMBs(c)
			
		}
		default:
		pur, mur, err = ColAzGen(c,c.Pu)
		switch c.Code{
			case 1:
			_ , _,_ = ColNMIs(c,0.0,true)
			case 2:
			_, _ = ColNMBs(c)
		}
	}
	if !c.Web{
		fmt.Printf("analyzed pur - %.4f kn mur - %.4f kn.m\n",pur,mur)
		fmt.Println(c.Nmplot)
	}
	c.Mux = mur 
	return
}

//Init initializes a column section struct
func (c *RccCol) Init(){
	//etc? maybe have global fcks and add here
	if c.Ignore{return}
	if c.Code == 0{c.Code = 1}
	if c.Cvrt == 0.0{
		if c.C1 == 0.0{
			c.Cvrt = 37.5; c.Cvrc = 37.5
		} else {
			c.Cvrt = c.C1; c.Cvrc = c.C1
		}
	}
	if c.C2 == 0.0{c.C2 = c.Cvrc}
	if c.C1 == 0.0{
		c.C1 = c.C2
	}
	if c.Nlayers == 0{c.Nlayers = 2}
	c.SecInit()
	c.Dtyp = 2
	switch{
		case c.Mux == 0.0 && c.Muy == 0.0:
		c.Dtyp = 0
		case c.Muy == 0.0:
		c.Dtyp = 1
	}
	c.Ast = 0.0
	c.Asc = 0.0
	c.Asteel = 0.0
	c.Psteel = 0.0
	c.Ptie = 0.0
	c.Rbrc = []float64{}
	c.Rbrt = []float64{}
	c.Rbropt = [][]float64{}
	c.Rbrcopt = [][]float64{}
	c.Rbrtopt = [][]float64{}
	c.Dias = []float64{}
	c.Barpts = [][]float64{}
	// return
}

//SecInit initializes a column section section :p
//initializes the c.Sec struct
func (c *RccCol) SecInit() (err error){
	if c.Dims == nil{
		//just for the n holy basic types
		switch c.Styp{
			case 0:
			c.Dims = []float64{c.B,c.B}
			case 1:
			c.Dims = []float64{c.B, c.H}
			case 2:
			c.Dims = []float64{c.B,c.B*math.Sqrt(3.0)/2.0}
			case 3:
			c.Dims = []float64{c.B,c.H}
			case 6,7,8,9,10,11,12,14:
			c.Dims = []float64{c.B, c.H, c.Bw, c.Df}
			//add - plus col, etc
			default:
			err = errors.New("non standard col section")
		}
	}
	sect := kass.SecGen(c.Styp, c.Dims)
	c.Sec = &sect
	c.Sec.SecInit()
	c.Sec.Draw("")
	return
}

//ColAzGen is for general ('-')7 column analysis
//computes column area via c.Sec, steel forces via c.Dbars and c.Dias
//returns pur (ult axial load)/mur (ultimate moment of resistance)

func ColAzGen(c *RccCol, pu float64) (pur, mur float64, err error){
	var k8, k9, ebar, fbar, diabar, ack, asec, yg, ygsec, dck, xu, step float64
	fsts := make([]float64, len(c.Dias))	
	k8 = 0.45; k9 = 0.9
	if c.Code == 1{k8 = 0.45; k9 = 0.8}
	//if c.Code == 3{k8 = 0.567; k9 = 0.8} whaat?
	pu = pu * 1e3
	if c.Sec == nil{c.SecInit()}
	if c.Styp == 0{
		if c.Csq != nil{
			return ColAzGen(c.Csq, c.Pu)
		} else {
			err = errors.New("equivalent square column missing")
			return 
		}
	}
	if c.H == 0{
		c.H = math.Abs(c.Sec.Ymx - c.Sec.Ym)
		if c.H == 0.0{err = ErrDim; return}
	}
	if c.Dbars == nil || len(c.Dbars) == 0{
		log.Println("ERRORE,errore->depth of bars not specified")
		err = ErrDim
		return
	}
	//get section area and depth of centroid from kompression face
	asec = c.Sec.Prop.Area; yg = c.Sec.Prop.Yc
	yg = c.Sec.Ymx - c.Sec.Prop.Yc
	
	//iterate
	step = 20.0
	xu = c.Cvrc + 1.0
	xu = xu - c.H/step

	kiter := 0
	for { 
		if kiter > 2000{
			//log.Println("ERRORE, errore-> maximum iteration limit reached")
			return pur, mur, ErrCAxL
		}
		if pur < pu{
			xu += c.H/step
		} else {
			xu -= c.H/step
			step = 20.0 * step
			xu -= c.H/step
		}
		if k9 * xu >= c.H{
			dck = c.H
			ack = asec
			ygsec = yg
		} else {
			dck = k9 * xu
			ack, _ , ygsec = ColSecArXu(c.Sec, dck)
			ygsec = c.Sec.Ymx - ygsec
		}
		pur = k8 * c. Fck * ack
		for idx, dbar := range c.Dbars{
			diabar = c.Dias[idx]
			ebar = math.Abs(0.0035 * (dbar - xu)/xu)
			switch {
			case xu > dbar:
				switch c.Code{
					case 1:
					fbar = RbrFrcIs(c.Fy, ebar) 
					default:
					fbar = RbrFrcBs(c.Fy, ebar)
				}
				if c.Subck{
					fbar -= k8 * c.Fck
				}
			case xu == dbar:
				//ZERO FORCE? really
				fbar = 0
			case xu < dbar:
				switch c.Code{
					case 1:
					fbar = -RbrFrcIs(c.Fy, ebar) 
					default:
					fbar = -RbrFrcBs(c.Fy, ebar)
				}
			}
			fsts[idx] = fbar
			pur += RbrArea(diabar)*fbar
		}
		if xu == c.Cvrc + 1.0 && pur > pu{
			break
		}
		if pur > pu && pur <= 1.01 * pu{
			break
		}
		kiter++
	}
	mur = pur * yg - k8 * c.Fck * ack * ygsec
	for idx, fbar := range fsts{
		mur -= RbrArea(c.Dias[idx])*fbar*c.Dbars[idx]
	}
	return pur/1e3, mur/1e6, nil
}

//ColStl estimates the area of steel as per bs and is codes (c.Code = 2, c.Code = 1)
//short column + uniaxial + biaxial bending area of steel
//dtyp 0 - axial load, 1 - axial load + uniaxial bending, 2 - biaxial bending
//rtyp 0 - symmetrical steel, nlayers (0-n)
//rtyp 1 - unsymmetrical steel, nlayers 2

func ColStl(c *RccCol, pu, mux, muy float64) (pur, murx float64, err error){
	var msfs, astot, k8, k9, mu, fst, fsc, est, esc, ey, dck, ack, step, xu, mur, effd float64
	//var astbckup float64
	if c.Fy ==250.0 {msfs = 0.77} else {msfs = 0.67}
	xstep := 20.0
	//default code bs (joota hai japani)
	if c.Code == 0{c.Code = 2}
	switch c.Code{
		case 1:
		k8 = 0.45; k9 = 0.8
		case 2:
		k8 = 0.45; k9 = 0.9
	}
	if c.Cvrt == 0.0{
		c.Cvrt = 30 + 6 + 16 + 3
		c.Cvrc = c.Cvrt
	}
	if c.Sec == nil{c.SecInit()}
	if c.H == 0.0{
		c.H = math.Abs(c.Sec.Ymx - c.Sec.Ym)
	}
	effd = c.H - c.Cvrt
	ey = c.Fy/1.15/2e5
	//maan. add biaxial bending and CLOSE
	if c.Slender{
		pur, murx, err = ColSlmDBs(c)
		c.Pur = pur; c.Murx = murx
		return
	}
	switch c.Dtyp{
		case 0:
		/*
		   axially loaded short column
		*/
		pu = pu * 1e3
		switch c.Styp{
			case 0:
			//circular
			ack = math.Pi * math.Pow(c.B,2)/4.0
			default:
			ack = c.Sec.Prop.Area
		}
		astot = (pu - 0.4 * c.Fck * ack)/(msfs*c.Fy - 0.4 * c.Fck)
		c.Asc = astot/2.0; c.Ast = c.Asc
		pur = pu
		murx = mux
		case 1:
		/*
		   axial load + uniaxial bending
		*/
		//fmt.Println("COLSTL!c.Rtyp, c.Nlayers,c.Code",c.Rtyp, c.Nlayers,c.Code,pu, mux)
		pu = pu * 1e3
		mu = mux* 1e6 //convert to N-mm
		var kiter int
		switch c.Rtyp{
			case 0:
			//symmetrical arrangement of rebar
			//get depth of centroid from top
			var ack, ycy, ygsec, yg, asec, zsc, zst, zck, fcc float64
			asec = c.Sec.Prop.Area; ygsec = c.Sec.Prop.Yc
			astmin := 0.4 * asec/100.0
			yg = c.Sec.Ymx - c.Sec.Prop.Yc
			switch c.Nlayers{
				case 0, 2:
				//two levels
				//CHECK FOR UNREINFORCED SECTION ??
				//fmt.Println("HYARR HYARR")
				xu = 2.33 * effd; step = xstep
				iter := 0; kiter = 0
				xu += c.H/step
				for iter != -1{
					kiter++; fsc = 0; fst = 0
					if kiter > 2999{
						//log.Println("ERRORE, errore->max iterations reached")
						err = ErrIter 
						//c.Astot = astbckup
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
						//xu -= c.H/step
						//continue
					}
					if xu * k9 >= c.H {
						dck = c.H
						ack = asec
						ycy = ygsec
					} else {
						dck = k9 * xu
						ack, _, ycy = ColSecArXu(c.Sec, dck)
					}
					ycy = c.Sec.Ymx - ycy
					esc = 0.0035 * (xu - c.Cvrc)/xu
					est = math.Abs(0.0035 * (effd - xu)/xu)
					//fmt.Println(ColorYellow,"esc,est->",xu, esc, est)
					switch c.Code{
						case 1:
						//is code
						fsc, fst = ColFstlIs(xu, effd, c.Fck, c.Fy, esc, est)
						case 2:
						//bs code
						fsc, fst = ColFstl(xu, effd, c.Fy, ey, esc, est)
					}
					if c.Subck{
						//SUBTRACT FCK FROM THE FORCE IN (noble) STEEL BARS
						fsc = fsc - k8*c.Fck
						//SUBTRACT IF (valiant) TENSILE STEEL IS IN KOMPRESSION
						if fst > 0.0 {fst = fst - k8*c.Fck}
					}
					zsc = yg - c.Cvrc; zst = yg - c.Cvrt
					fcc = k8 * ack * c.Fck; zck = yg - ycy
					
					if iter == 0{
						astot = 2.0 * (pu - fcc)/(fsc - fst)
					}
					if iter == 1{
						astot = 2.0 * (mu - fcc * zck)/(fsc * zsc - fst * zst)
					}
					//fmt.Println(ColorRed,"xu,esc,fsc,est,fst,astot",xu,esc,fsc,est,fst,astot,ColorReset)
					if astot < 0.0 || astot > 0.04 * asec{
						continue
					}
					if iter == 0{
						mur = fcc * zck + astot/2.0 * (fsc * zsc + fst * zst)
						if c.Approx && mur >= mu{
							iter = -1
							break
						}
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
						
					}
					if iter == 1{
						pur = fcc + astot * (fsc + fst)/2.0
						if (pu == 0 && pur - pu < 0.1) || (xu == c.Cvrc + 1.0 && pur > pu) {
							iter = -1
							break
						}
						
						if pur >= pu {
							
							if c.Approx{
								iter = -1
								break
							}
							if step >= xstep * 10.0{
								iter = -1
								break
							}
							if 1.02 * pu >= pur{
								iter = -1
								break
							} else {
								xu -= c.H/step
								step = step * 10.0
								xu -= c.H/step
							}
						}
					}
				}
				if astot < astmin {astot = astmin}
				pur = fcc + astot * (fsc + fst) /2.0
				mur = fcc * zck + astot/2.0 * (fsc * zsc - fst * zst)
				c.Asc = astot/2.0; c.Ast = c.Asc; c.Asteel = astot; c.Xu = xu
				pur = pur/1e3
				mur = mur/1e6
				murx = mur
				c.Psteel = c.Asteel*100.0/asec
				if c.Psteel > 4.0{
					err = errors.New("steel percentage greater than 4")
					return
				} 
				default:
				//AGAIN.FOR THE LAST goddamn time
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
						xu = c.Cvrc + c.H/step
						step = xstep
						//continue
					}
					if xu > c.H/k9 {
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
					if astot < 0.0 || astot > 0.04 * asec{
						continue
					}
					if iter == 0{
						mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * zsum/float64(c.Nlayers)
						if math.Round(mur) >= math.Round(mu){
	
							if c.Approx{
								iter = -1
								break
							}						
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
					}
					if iter == 1{
						pur = k8 * c.Fck * c.B * dck + astot * fsum/float64(c.Nlayers)
						if (pu == 0 && pur - pu < 0.1) || (xu == c.Cvrc + c.H/step && pur > pu) {
							iter = -1
							break
						}	
						if math.Round(pur) >= math.Round(pu){
							if c.Approx{
								iter = -1
								break
							}
							if step >= xstep * 10.0{
								iter = -1
								break
							}
							if 1.02 * pu >= pur {
								iter = -1
								break
							} else {
								xu -= c.H/step
								step = step * 10.0
								xu -= c.H/step
							}
						}
						
						//if math.Abs(pur-pu) < 0.02 * pu && c.Approx{
						//	iter = -1
						//	break
						//}
					}
				}
				if astot < astmin{astot = astmin}
				pur = fcc + astot * (fsc + fst) /2.0
				mur = fcc * zck + astot/2.0 * (fsc * zsc - fst * zst)
				c.Asc = astot/2.0; c.Ast = c.Asc; c.Asteel = astot; c.Xu = xu
				pur = pur/1e3
				mur = mur/1e6
				murx = mur
				c.Psteel = c.Asteel*100.0/asec
				if c.Psteel > 4.0{
					err = errors.New("steel percentage greater than 4")
					return
				} 
			}
			case 1:
			//unsymmetrical arrangement of rebar
			var ack, ycy, ygsec, yg, asec, zsc, fcc float64
			asec = c.Sec.Prop.Area; ygsec = c.Sec.Prop.Yc
			yg = c.Sec.Prop.Yc - c.Sec.Ym
			switch c.Nlayers{
				case 0,2:
				//two levels again
				var ast, asc, astprev, ascprev float64
				//log.Println("pu,mu->",pu,mu)
				effd = c.H - c.Cvrt
				iter := 0; step = 100.0; ast = -1.0; asc = -1.0
				xu = c.Cvrc + 1.0
				xu -= c.H/step
				var kiter int
				for iter != -1 {
					kiter++; fsc = 0.0; fst = 0.0
					if kiter > 2999{
						log.Println("ERRORE,errore->max no. of iterations reached")
						err = ErrD
						return
					}
					xu += c.H/step
					astprev = ast; ascprev = asc
					if xu * k9 >= c.H {
						dck = c.H
						ack = asec
						ycy = ygsec
					} else {
						dck = k9 * xu
						ack, _, ycy = ColSecArXu(c.Sec, dck)
					}
					esc = 0.0035 * (xu - c.Cvrc)/xu
					est = math.Abs(0.0035 * (effd - xu)/xu)
					switch c.Code{
						case 1:
						fsc, fst = ColFstlIs(xu, effd, c.Fck, c.Fy, esc, est)
						case 2:
						fsc, fst = ColFstl(xu, effd, c.Fy, ey, esc, est)
					}
					if c.Subck{
						fsc = fsc - k8 * c.Fck
						if fst > 0.0 {fst = fst - k8*c.Fck}
					}
					fcc = k8 * ack * c.Fck
					zsc = effd - c.Cvrc
					asc = (mu + pu * (yg - c.Cvrt) - fcc * (ycy - c.Cvrt))/fsc/zsc
					ast = (pu - fcc - asc* fsc)/fst
					if (ast < 0.0 || astprev < 0.0 || asc < 0.0 || ascprev < 0.0) {
						continue
					}
					if asc + ast >= ascprev + astprev {
						if iter == 0 {
							xu -= 2.0 * effd/step
							step = step * 10.0
							iter = 1
							xu -= effd/step
						} else {
							iter = -1
							break
						}
					}
				}
				astot = ascprev + astprev
				c.Asc = ascprev; c.Ast = astprev
				c.Asteel = astot
				pur = pu/1e3
				mur = mu/1e6
				c.Xu = xu
				
				c.Psteel = c.Asteel*100.0/asec
				if c.Psteel > 4.0{
					err = errors.New("steel percentage greater than 4")
					return
				} 
			}
			case 2:
			//LN level dbar gen
			//..and again lmao
			var ack,  ygsec, yg, asec float64
			asec = c.Sec.Prop.Area; yg = c.Sec.Prop.Yc
			astmin := 0.4 * asec/100.0
			//get dbars
			c.BarDbars()
			var fsum, zsum, ebar, fbar float64
			xu = 2.33 * effd; step = xstep
			iter := 0; kiter = 0
			xu += c.H/step
			for iter != -1{
				kiter++
				zsum, fsum = 0.0,0.0
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
					xu = c.Cvrc
					step = xstep
					continue
				}
				if k9 * xu >= c.H{
					dck = c.H
					ack = asec
					ygsec = yg
				} else {
					dck = k9 * xu
					ack, _, ygsec = ColSecArXu(c.Sec, dck)
				}
				for _, dbar := range c.Dbars{
					ebar = math.Abs(0.0035 * (xu - dbar)/xu)
					switch {
					case xu == dbar:
						fbar = 0.0//wont the force func take care of this, ebar == 0?
					case xu > dbar:
						switch c.Code{
							case 1:
							fbar = RbrFrcIs(c.Fy, ebar)
							case 2:
							fbar = RbrFrcBs(c.Fy, ebar)
						}
						if c.Subck{fbar -= k8 * c.Fck}
						fsum += fbar
						//lever arm about centerline of column (yg)
						zsum += (yg - dbar) * fbar
					case xu < dbar:
						switch c.Code{
							case 1:
							fbar = -RbrFrcIs(c.Fy, ebar)
							case 2:
							fbar = -RbrFrcBs(c.Fy, ebar)
						}
						fsum += fbar
						zsum += (yg - dbar) * fbar
					}
				}
				if iter == 0{
					astot = float64(len(c.Dbars)) * (pu - k8 * c.Fck * ack)/(fsum)
				}
				if iter == 1{
					astot = float64(len(c.Dbars)) * (mu - k8 * c.Fck * ack * math.Abs(yg - ygsec))/zsum
				}
				if astot < 0.0 || astot > 0.04*asec{
					continue
				}
				if iter == 0{
					mur = k8 * c.Fck * ack * math.Abs(yg - ygsec) + astot * zsum/float64(len(c.Dbars))
					if math.Round(mur) >= math.Round(mu){
						if c.Approx{
							iter = -1
							break
						}
						if 1.02 * mu >= mur{
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
					pur = k8 * c.Fck * ack + astot * fsum/float64(len(c.Dbars))
					if (pu == 0 && pur - pu < 0.1) || (xu == c.Cvrc && pur > pu) {
						iter = -1
						break
					}
					if math.Round(pur) >= math.Round(pu){
						
						if c.Approx{
							iter = -1
							break
						}
						if 1.02 * pu >= pur{
							iter = -1
							break
						} else {
							xu -= c.H/step
							step = step * 10.0
							xu -= c.H/step
						}
					}

					//if math.Abs(pur-pu) <= 0.03 * pu{
					//	iter = -1
					//	break						
					//}
				}
			}
			if astot < astmin{astot = astmin}
			//fmt.Println("iter finito->", astot)
			pur = k8 * c.Fck * ack + astot * fsum/float64(len(c.Dbars))
			mur = k8 * c.Fck * ack * math.Abs(yg - ygsec) + astot * zsum/float64(len(c.Dbars))
			c.Asc = astot/2.0; c.Ast = c.Asc; c.Asteel = astot; c.Xu = xu
			pur = pur/1e3
			mur = mur/1e6
			murx = mur
			c.Psteel = c.Asteel*100.0/asec
			if c.Psteel > 4.0{
				err = errors.New("steel percentage greater than 4")
				return
			} 
		}
		case 2:
		//now one has mux, muy
		switch{
			//case c.Styp == 1:
			//switch{
			case c.Code == 2  && c.Styp == 1 && c.Rtyp == 0:
			//use simplified bs code method
			c1 := pu/c.B/c.H/c.Fck*1e3
			var c2 float64
			if c1 >= 0.6 {c2 = 0.3} else {c2 = 0.3 + (0.7/0.6)*(0.6 - c1)}
			effb := c.B - c.C2
			effd = c.H - c.C1
			if mux/effd > muy/effb {
				c.Ax = "x"
				mu = mux + c2 * (effd/effb) * muy
				c.Cvrc = c.C1; c.Cvrt = c.C1
			} else {
				//flip
				c.Ax = "y"
				mu = muy + c2 * (effb/effd) * mux
				c.Cvrc = c.C2; c.Cvrt = c.C2
				c.H, c.B = c.B, c.H
			}
			c.Dtyp = 1
			return ColStl(c, pu, mu, 0.0)
			default:
			//here test random nonsense 
			//get pt and iterate as per bxrectast
			var step, pnz, pg, mnx, mny, alpha, mrat float64
			asec := c.Sec.Prop.Area
			kiter := 0; iter := 0
			pg = 0.6; step = 0.1
			pg -= step
			cy := ColFlip(c)
			c.Init()
			//c.Sec.Draw("dumb")
			//fmt.Println(c.Sec.Txtplot)
			//cy.Sec.Draw("dumb")
			//fmt.Println(cy.Sec.Txtplot)
			rbrlvl := c.Rbrlvl
			for iter != -1{
				kiter++
				switch iter{
					case 0:
					pg += step
					case 1:
					pg -= step
				}
				//pg += step
				astl := pg * asec * 1e-2
				
				pnz = 0.45 * c.Fck * asec + 0.75 * c.Fy * astl
				pnz = pnz/1e3
				prat := pu/pnz
				if kiter > 666{
					err = ErrIter
					return
				}
				switch{
					case prat <= 0.2:
					alpha = 1.0
					case prat >= 0.8:
					alpha = 2.0
					default:
					alpha = 2.0/3.0 + 5.0 * pu/3.0/pnz
				}
				//fmt.Println("prat-",pnz, pu,pu/pnz,"astl",astl)
				c.Asteel = astl; cy.Asteel = astl
				c.Dias = []float64{}; cy.Dias = []float64{}; c.Barpts = [][]float64{}
				c.Dbars = []float64{}; cy.Dbars = []float64{}; cy.Barpts = [][]float64{}
				c.Rbrlvl = rbrlvl; cy.Rbrlvl = rbrlvl; c.Lsec = kass.SectIn{}; cy.Lsec = kass.SectIn{}

				//c.Reset(); cy.Reset()
				if c.Styp == 1{
					if c.Rtyp == 0 && c.Nlayers > 2{
						err = ColRbrGen(c)
						if err != nil{
							continue
							//fmt.Println(err)
						}
						err = ColRbrGen(cy)
						
						if err != nil{
							continue
							//fmt.Println(err)
						}
					} else {
						c.Ast = astl/2.0; c.Asc = c.Ast
						cy.Ast = astl/2.0; cy.Asc = cy.Ast
						err = c.BarGen()
						
						if err != nil{
							continue
							//fmt.Println(err)
						}
						err = cy.BarGen()
						
						if err != nil{
							continue
							//fmt.Println(err)
						}
					}
				} else {
					c.BarDbars(); cy.BarDbars()
					err = c.BarGen()
					if err != nil{
						continue
						//fmt.Println(err)
					}
					err = cy.BarGen()
					if err != nil{
						continue
						//fmt.Println(err)
					}
				}
				//c.DBarpts(); cy.DBarpts()
				//c.BarLay()
				//cy.BarLay()
				_, mnx, err = ColAzGen(c,c.Pu)
				if err != nil{
					//fmt.Println(err)
					continue
				}
				_, mny, err = ColAzGen(cy,cy.Pu)
				if err != nil{
					//fmt.Println(err)
					continue
				}
				//fmt.Println("MNX, MNY->",mnx, mny)
				mrat = math.Pow(mux/mnx, alpha) + math.Pow(muy/mny, alpha)
				//log.Println("step, pg, mrat ->",step, pg , mrat)
				if mrat <= 1.0{
					
					c.Asteel = astl
					c.Dias = []float64{}; c.Barpts = [][]float64{}
					c.Dbars = []float64{}
					c.Rbrlvl = rbrlvl
					c.Lsec = kass.SectIn{}
					iter = -1
					break
				}
			}
		}
	}
	//c.Draw()
	
	return 
}

//Draw generates a gnuplot data string to plot a column section 
func (c *RccCol) Draw() (string){
	//get base section points
	data := c.Sec.Data[0]
	//fmt.Println(data)
	ldata := ""; cdata := ""
	switch c.Rtyp{
		case 2:
		for i, dia := range c.Dias{
			pt := c.Barpts[i]
			cdata += fmt.Sprintf("%f %f %.f %.f\n", pt[0], pt[1], dia, dia)
			ldata += fmt.Sprintf("%f %f %.f %.f\n", pt[0]+10.0, pt[1]+10.0, dia, dia)
		}
		default:
		for i, pt := range c.Barpts{
			cdata += fmt.Sprintf("%f %f %.f %.f\n", pt[0], pt[1], c.Dias[i], c.Dias[i])
			ldata += fmt.Sprintf("%f %f %.f %.f\n", pt[0]+10.0, pt[1]+10.0, c.Dias[i], c.Dias[i])
		}
		//ldata += fmt.Sprintf("%f %f %.fmm2\n",(c.Sec.Xmx - c.Sec.Xm)/2.0, c.Sec.Ymx - c.Cvrc - 50.0, c.Asc)
		//ldata += fmt.Sprintf("%f %f %.fmm2\n",(c.Sec.Xmx - c.Sec.Xm)/2.0, c.Sec.Ym + c.Cvrt - 25.0, c.Ast)

	}
	ldata += fmt.Sprintf("%f %f %.fmm2\n",c.Sec.Prop.Xc, c.Sec.Prop.Yc, c.Asteel)
	ldata += fmt.Sprintf("%f %f %.3f\n",c.Sec.Prop.Xc, c.Sec.Prop.Yc-25.0, c.Psteel)
	c.Plotdat = []string{data, ldata, cdata}
	data += "\n\n"; ldata += "\n\n"; cdata += "\n\n"
	data += ldata; data += cdata
	data += fmt.Sprintf("%f %f\n", c.Sec.Xm, c.Sec.Ymx - c.Xu)
	data += fmt.Sprintf("%f %f\n", c.Sec.Xmx, c.Sec.Ymx - c.Xu)
	data += "\n\n"
	c.Lsec.Draw("")
	data += c.Lsec.Data[0]
	//c.Draw2d()
	c.Data = data
	return data
}

//TODO
//Draw2d draws the front view of the column
//see fig 13.31 of subramanian
func (c *RccCol) Draw2d(){
	//pb := make([]float64,2)
	//pe := make([]float64,2)
	
}

//Reset is used in ColBx - resets all c.Dia and c.Dbar slices
func (c *RccCol) Reset(){
	//for colbx- resets all dia and dbar arrays
	rbrlvl := c.Rbrlvl
	c.Lsec = kass.SectIn{}
	c.Dias = []float64{}
	c.Barpts = [][]float64{}
	c.Dbars = []float64{}
	//c.Dxbars = []float64{}
	c.Dybars = []float64{}
	c.Rbrlvl = rbrlvl
	
}

//BarLay generates column barpts from a rebar template for rtyp 0 and 1 (rtyp 2 = c.Lsec.Coords[:-1])
func (c *RccCol) BarLay(dx float64, pts[][]float64, tcdx int) (err error){
	//rez = []float64{float64(n1), float64(n2), d1, d2, astmin, ast, adiff}
	//rez = []float64{nlayer, astprov, efcvr, efdp, cldis, clvdis, nbarRow} 
	var xi, yi, xs, ys, xstep, ystep, d1, d2 float64
	var n1, n2, nbarr, nlayer int
	if c.Rtyp == 2{return}
	switch tcdx{
		case 1:
		//asc
		n1 = int(c.Rbrc[0]); n2 = int(c.Rbrc[1]); d1 = c.Rbrc[2]; d2 = c.Rbrc[3]
		xstep = c.Rbrc[11]; ystep = c.Rbrc[12]; nbarr = int(c.Rbrc[13]); nlayer = int(c.Rbrc[7])
		case 2:
		//ast
		n1 = int(c.Rbrt[0]); n2 = int(c.Rbrt[1]); d1 = c.Rbrt[2]; d2 = c.Rbrt[3]
		xstep = c.Rbrt[11]; ystep = c.Rbrt[12]; nbarr = int(c.Rbrt[13]); nlayer = int(c.Rbrt[7])
	}
	dmax := d2
	xs  = pts[0][0] + 25.0; ys = pts[0][1]
	xi = xs; yi = ys
	if d1 > dmax || n2 == 0 {dmax = d1}
	switch nlayer{
		case 1:
		nbarr = n1 + n2
	}
	if nbarr > 1{
		xstep = (dx-50.0)/float64(nbarr - 1)
	}
	ystep += dmax
	for i := 0; i < n1 + n2; i++{
		if i > nbarr{
			xi = xs; yi = ys + ystep
		}
		if i < n1{
			c.Dias = append(c.Dias, d1)
		} else {
			c.Dias = append(c.Dias, d2)
		}
		if nbarr == 1{
			c.Barpts = append(c.Barpts, []float64{xs -25.0 + dx/2.0,yi})
		} else {
			c.Barpts = append(c.Barpts, []float64{xi, yi})
		}
		xi += xstep
	}
	return
}


//fmt.Println("MNX, MNY->",mnx, mny)
//fmt.Println(ColorRed)
//log.Println("optimal percentage of steel found")
//log.Println("alpha, mrat, 1.0 - mrat",alpha, mrat, 1.0 - mrat)
//log.Println("asteel, pg->",astl, pg)
//fmt.Println(ColorReset)
//fmt.Println(ColorRed)
//fmt.Println("dbarz->",c.Dbars, c.Dias)
//fmt.Println(cy.Dbars, cy.Dias)
//fmt.Println(ColorReset)

//c.BarDbars()
//err = c.BarGen()
//if err != nil{
//	return
//}
//c.Term = "qt"
//c.Plot("qt")
//fmt.Println(c.Txtplot)
//c.Dias = make([]float64, len(c.Dbars))
//c.Asteel = astl; cy.Asteel = astl
//c.Dias = []float64{}; cy.Dias = []float64{}; c.Barpts = [][]float64{}
//c.Dbars = []float64{}; cy.Dbars = []float64{}; cy.Barpts = [][]float64{}
//c.Rbrlvl = rbrlvl; cy.Rbrlvl = rbrlvl; c.Lsec = kass.SectIn{}; cy.Lsec = kass.SectIn{}
