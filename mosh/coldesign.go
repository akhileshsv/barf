package barf

import (
	"fmt"
	"log"
	"math"
	"errors"
	"math/rand"
	kass"barf/kass"
	draw"barf/draw"
)

//RccCol is a struct to store rcc column fields
//see hulse chapter 5
type RccCol struct{
	Title           string
	Id              int
	Mid             int
	Fck             float64
	Fy              float64
	Nomcvr          float64
	Dtyp            int
	Styp            int
	Rtyp            int
	Lspan           float64 //c-c distance between supports
	Efcvr           float64
	Cvrt            float64
	Cvrc            float64
	Bw, Df          float64
	B, H            float64
	Dims            []float64
	Code            int
	Pu, Mux, Muy    float64
	Term            string
	Web             bool `json:",omitempty"`
	Dsgn            bool `json:",omitempty"`
	Subck           bool `json:",omitempty"`
	Ignore          bool `json:",omitempty"`
	L0              float64 `json:",omitempty"` //clear distance
	Lp              float64 `json:",omitempty"` //length of plastic hinge
	Lx              float64 `json:",omitempty"`
	Ly              float64 `json:",omitempty"`
	Leffx           float64 `json:",omitempty"`
	Leffy           float64 `json:",omitempty"`
	Ast             float64 `json:",omitempty"`
	Astl            float64 `json:",omitempty"`
	Asc             float64 `json:",omitempty"`
	Nbars           float64 `json:",omitempty"`
	Coords          [][]float64 `json:",omitempty"` 
	Dias            []float64 `json:",omitempty"`
	Dbars           []float64 `json:",omitempty"`
	Dybars          []float64 `json:",omitempty"`
	Barpts          [][]float64 `json:",omitempty`
	Ptie            float64 `json:",omitempty"`
	Dtie            float64 `json:",omitempty"`
	Ptiec           float64 `json:",omitempty"` //pitch of confining reinforcement
	Nties           float64 `json:",omitempty"`
	Tietyp          string `json:",omitempty"`
	Type            string `json:",omitempty"`
	Data            string `json:",omitempty"`
	Nmplot          string `json:",omitempty"`
	Report          string `json:",omitempty"`
	Txtplot         []string `json:",omitempty"`
	Txtplots        []string `json:",omitempty"`
	Plotdat         []string `json:",omitempty"`
	Sec             *kass.SectIn `json:",omitempty"`
	D1, D2          float64 `json:",omitempty"`
	N1, N2          float64 `json:",omitempty"`
	C1, C2          float64 `json:",omitempty"`
	Mt, Mb          float64 `json:",omitempty"`
	L1, L2, L3      float64 `json:",omitempty"` //confining reinforcement spans
	Ax              string `json:",omitempty"`
	Nlayers         int `json:",omitempty"`
	Xu              float64 `json:",omitempty"`
	Asteel          float64 `json:",omitempty"`
	Psteel          float64 `json:",omitempty"`
	Na, Nb          int `json:",omitempty"`
	Da, Db, Dc      float64 `json:",omitempty"`
	Pur, Murx, Mury float64 `json:",omitempty"`
	Dz              int `json:",omitempty"`
	Bdim            [][]float64 `json:",omitempty"`
	Bsec            []int `json:",omitempty"`
	Cdim            [][]float64 `json:",omitempty"`
	Csec            []int `json:",omitempty"`
	Lsec            kass.SectIn `json:",omitempty"`
	Rbrt, Rbrc      []float64 `json:",omitempty"`
	Rbr             []float64 `json:",omitempty"`
	Rbropt          [][]float64 `json:",omitempty"`
	Rbrtopt         [][]float64 `json:",omitempty"`
	Rbrcopt         [][]float64 `json:",omitempty"`
	Rbrlvl          int `json:",omitempty"`
	Braced          bool `json:",omitempty"`
	Verbose         bool `json:",omitempty"`
	Ljbase          bool `json:",omitempty"`
	Basemr          bool `json:",omitempty"`
	Slender         bool `json:",omitempty"`
	Approx          bool `json:",omitempty"` //design approximately for abs(mu - mur)/(pu - pur)
	Cdims           [][]float64 `json:",omitempty"`
	Bxdim           [][]float64 `json:",omitempty"` // lo-l, lo-r, hi-l, hi-r
	Bydim           [][]float64 `json:",omitempty"` // lo-l, lo-r, hi-l, hi-r
	Blens           [][]float64 `json:",omitempty"`
	Bxks            []float64 `json:",omitempty"`
	Byks            []float64 `json:",omitempty"`
	Csq             *RccCol `json:",omitempty"` //store csq for styp 0 (for now)
	Vrcc, Wstl      float64 `json:",omitempty"`
	Vtot, Afw       float64 `json:",omitempty"` //WHY DO YOU NEED VTOT isnt it Vrcc            
	Kost            float64 `json:",omitempty"`
	Kostin          []float64 `json:",omitempty"`
	Foldr           string `json:",omitempty"`
}

//ColDzBs is a routing func for bs code column design routines 
func ColDzBs(c *RccCol) (err error) {
	//routing func similar to ColDzIs. does nothing
	var pur, murx, mury float64
	switch c.Styp{
		case 0:
		//circle
		csq := ColEqSqr(c)
		switch c.Slender{
			case true:
			pur, murx, err = ColSlmDBs(csq)
			case false:
			pur, murx, err = ColDBs(csq, csq.Pu,csq.Mux,csq.Muy)
			
		}	
		if err != nil{
			//if c.Verbose{log.Println(ColorRed,"ERRORE,errore->",err,ColorReset)}
			return
		}
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
		//fmt.Println("tiez->",c.Ptie,c.Dtie,c.Nties)
		//fmt.Println("diaze->",c.Dias)
		case 1:
		//good ol' rect col funcs
		switch c.Slender{
			case true:
			pur, murx, err = ColSlmDBs(c)
			//log.Println(ColorRed,pur,murx,ColorReset)
			case false:
			pur, murx, err = ColDBs(c, c.Pu,c.Mux,c.Muy)
		}	
		if err != nil{
			//if c.Verbose{log.Println(ColorRed,"ERRORE,errore->",err,ColorReset)}
			return
		}
		switch c.Rtyp{
			case 1:
			//wot absolute shite
			c.Pur = c.Pu
			default:
			c.Pur = pur
		}
		c.Murx = murx; c.Mury = mury
		
	}
	return
}

//ColSlmAdd calcs additional moments due to slenderness
//only for rectangular symmetrically reinforced section
//(simplified bs code method), hulse sec 5.8
func ColSlmMadd(pu, m1, m2, ecct, le, b, h, kred float64, braced bool) (mt float64){
	madd := pu * h * math.Pow(le/b,2)/2000.0
	m8 := math.Abs(m1); m9 := math.Abs(m2)
	m7 := m8; if m7 >= m9 {m7 = m9}
	m6 := m8; if m6 <= m9 {m6 = m9}
	if braced {
		switch {
		case m1 * m2 > 0:
			//double curvature
			m1 = -m7; m2 = m6
		default:
			//single curvature
			m1 = m7; m2 = m6
		}
		m1 = 0.4 * m1 + 0.6 * m2
		if m1 <= 0.4 * m2 {m1 = 0.4 * m2}
		mt = m1 + madd * kred
		if mt <= math.Abs(-m1 + madd * kred/2.0) {mt = math.Abs(-m1 + madd * kred/2.0)}
	} else {
		m1 = m6; m2 = m6
		mt = m1 + madd * kred
	}
	if mt <= m2 {mt = m2}
	if mt < pu * ecct {mt = pu * ecct}
	return
}

//ColSlmDBs designs a slender symmetrically reinforced rectangular column (hulse sec 5.8)
//calls ColSlmAdd to calculate additional moments due to slenderness 
//m1, m2 - top and bottom moments
//clockwise moments are +ve
//AH OKAY this is where effective height comes into play

func ColSlmDBs(c *RccCol) (pur, mur float64, err error){
	var k8, k9, ey, est, esc, fst, fsc, dck, ecct, xu, le, b, h float64
	le = c.Leffx; if le < c.Leffy {le = c.Leffy}
	b = c.B; h = c.H; if b > c.H {b = c.H; h = c.B}
	pu := c.Pu * 1e3; m1 := c.Mt * 1e6; m2 := c.Mb * 1e6; le = le * 1e3
	k8 = 0.45; k9 = 0.9
	ey = c.Fy/1.15/2e5
	ecct = c.H/20.0
	if ecct > 20.0 {ecct = 20.0}
	if c.Cvrt == 0.0{
		if c.C1 == 0.0{
			c.Cvrt = 37.5; c.Cvrc = 37.5
		} else {
			c.Cvrt = c.C1; c.Cvrc = c.C1
		}
	}
	switch c.Styp {
		case 0:
		//circular section
		//is eq.sqr
		//log.Println("getting equivalent square column")
		csq := ColEqSqr(c)
		return ColSlmDBs(csq)
	case 1:
		//rectangular section
		//okay what really are these conditions? read
		if h/b > 3.0 {
			//log.Println("ERRORE,errore -> h/b > 3")
			err = errors.New("h/b > 3.0, calculation method not appropriate")
			return 
		}
		if h > b && le/h > 20.0 {
			//log.Println("ERRORE,errore -> design for biaxial bending")
			err = errors.New("le/h>20.0, design as biaxially bent with zero initial moment about minor axis")
			return 
		}
		switch c.Rtyp {
		case 0:
			//symmetrical arrangement o' reinforcement
			kred := 1.0
			mu := ColSlmMadd(pu, m1, m2, ecct, le, b, h, kred, c.Braced)
			dck = pu/b/k8/c.Fck
			if pu * (h - dck)/2.0 > mu {
				//get min. steel
				c.Asc = 0.4 * c.B * h/2.0
				c.Ast = c.Asc
				return mu, pu, nil
			}
			var asprev, astot float64
			astot = 10.0
			effd := h - c.Cvrc
			for math.Abs(astot - asprev) >= 0.02 * asprev {
				//iterate
				asprev = astot
				xu = 2.33 * effd; step := 10.0
				iter := 0
				xu += h/step 
				for iter != -1 {
					//iterate
					switch iter {
					case 0:
						xu -= h/step
					case 1:
						xu += h/step
					}
					if xu < 0.9 * effd && iter == 0 {
						iter = 1
						xu = c.Cvrc + 1.0
						step = 10.0
					}
					if xu > h/k9 {
						dck = h
					} else {
 						dck = k9 * xu
					}
					esc = 0.0035 * (xu - c.Cvrc)/xu
					est = math.Abs(0.0035 * (effd - xu)/xu)
					fsc, fst = ColFstl(xu, effd, c.Fy, ey, esc, est)
					if iter == 0 {
						astot = 2.0 * (pu - k8 * c.Fck * b * dck)/(fsc + fst)
					}
					if iter == 1 {
						astot = 2.0 * (mu - k8 * c.Fck * b * dck * (h - dck)/2.0)/(fsc - fst)/(h/2.0 - c.Cvrc)
					}
					if astot < 0.0 {
						continue
					}
					if iter == 0 {
						mur = k8 * c.Fck * b * dck * (h - dck)/2.0 + astot * (fsc - fst) * (h/2.0 - c.Cvrc)/2.0
						if mur >= mu {
							if 1.02 * mu >= mur {
								//if mur - mu < 0.02 * mu || mur == mu {
								iter = -1
								break
							} else {
								xu += h/step 
								step = step * 10.0
								xu += h/step
							}
						}
						
						if math.Abs(mur-mu) < 0.02 * mu && step >= 100{
							iter = -1
							break
						}
					}
					if iter == 1 {
						pur = k8 * c.Fck * b * dck + astot * (fsc + fst)/2.0
						if (pu == 0 && pur - pu < 0.1) || (xu == c.Cvrc + 1.0 && pur > pu) {
							iter = -1
							break
						}
						if pur >= pu {
							if 1.02 * pu >= pur {
								//if (pur - pu < 0.02 * pu) || pur == pu {
								iter = -1
								//break
							} else {
								xu -= h/step
								step = step * 10.0
								xu -= h/step
							}
						}
						
						if math.Abs(pur-pu) < 0.02 * pu && step >= 1000.0{
							iter = -1
							break
						}
					}
				}
				puz := 0.45 * c.Fck * (b*h - astot) + 0.87 * c.Fy * astot
				pbal := 0.25 * c.Fck * b * effd
				kred = (puz - pu)/(puz - pbal)
				if kred > 1.0 {kred = 1.0}
				mu = ColSlmMadd(pu, m1, m2, ecct, le, b, h, kred, c.Braced)
			}
			pur = k8 * c.Fck * b * dck + astot * (fsc + fst) /2.0
			mur = k8 * c.Fck * b * dck * (h - dck)/2.0 + astot * (fsc - fst) * (h/2.0 -c.Cvrc)/2.0
			c.Asc = astot/2.0; c.Ast = c.Asc
		case 1:
			//unsymmetrical arrangement o' reinforcement
			//how, any leads?
		}
	case -1:
		//generic section
		//how mofo - see timoshenko i guess
	}
	pur, mur = pur/1e3, mur/1e6
	c.Asteel = c.Asc + c.Ast
	err = nil
	return 
}


//ColDBs calcs area of steel returns ultimate axial load pur, ultimate moment of resistance mur
//mostly hulse chapter 5 basic routines
//muy > 0 for biaxial bending
func ColDBs(c *RccCol, pu, mux, muy float64) (float64, float64, error){
	var k8, k9, ey, est, esc, fst, fsc, pur, dck, xu, step, mur, ecct, astot, effd, mu float64
	if muy == 0 {
		mu = mux
	} else if c.Styp == 1 {
		//get major axis of bending (simplified method)
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
			//WHAAT
			//h := c.B + c.D2; b := c.H + c.D1
			c.Cvrc = c.C2; c.Cvrt = c.C2
			c.H, c.B = c.B, c.H
			//c.H = h
			//c.B = b
		}
	}
	k8 = 0.45; k9 = 0.9
	ey = c.Fy/1.15/2e5
	pu = pu * 1e3; mu = mu * 1e6
	ecct = c.H/20.0
	if ecct > 20.0 {ecct = 20.0}
	if mu < pu * ecct {mu = pu * ecct}
	//log.Println(c.Id,"design moment, axial load ->",mu/1e6, pu/1e3)
	var kiter int
	xstep := 10.0
	switch c.Styp {
	case 1:
		//rectangular
		//fmt.Println("RECT SECKT DESIGN")
		switch c.Rtyp {
		case 0:
			//symmetrical arrangement of reinforcement
			effd = c.H - c.Cvrt
			//log.Println("effective depth ->",effd)
			dck = pu/c.B/k8/c.Fck
			//if pu * (c.H - dck)/2.0 > mu {
				//get min. steel
			//	c.Asc = 0.4 * c.B * c.H/2.0/100.0
			//	c.Ast = c.Asc
			//	return mu, pu, nil
			//}
			xu = 2.33 * effd; step = xstep
			iter := 0
			xu += c.H/step
			for iter != -1 {
				kiter++
				if kiter > 2999{
					return pur, mur, ErrIter
				}
				switch iter {
				case 0:
					xu -= c.H/step
				case 1:
					xu += c.H/step
				}
				if xu <= 0.9 * effd && iter == 0 {
					iter = 1
					xu = c.Cvrc + 1.0
					step = xstep
					//xu = xu - c.H/step
					//xu += c.H/step
				}
				//ADD DCK FUNC
				if xu > c.H/k9 {
					dck = c.H
				} else {
					dck = k9 * xu
				}
				esc = 0.0035 * (xu - c.Cvrc)/xu
				est = math.Abs(0.0035 * (effd - xu)/xu)
				fsc, fst = ColFstl(xu, effd, c.Fy, ey, esc, est)
				if iter == 0 {
					astot = 2.0 * (pu - k8 * c.Fck * c.B * dck)/(fsc + fst)
				}
				if iter == 1 {
					astot = 2.0 * (mu - k8 * c.Fck * c.B * dck * (c.H - dck)/2.0)/(fsc - fst)/(c.H/2.0 - c.Cvrc)
				}
				if astot < 0.0 {
					continue
				}
				if iter == 0 {
					mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * (fsc - fst) * (c.H/2.0 - c.Cvrc)/2.0
					if mur >= mu{
						//if xu == 2.33 * effd {
						//	iter = -1
						//}
						if step >= xstep * 10.0{
							iter = -1
							break
						}
						if 1.02 * mu >= mur{
						//if mur - mu < 0.02 * mu || mur == mu || 1.02 * mu >= mur{
							iter = -1
							break
						} else {
							xu += c.H/step 
							step = step * 10.0
							xu += c.H/step
						}
						
						//if math.Abs(mur-mu) < 0.03 * mu && step >= xstep * 10.0{
						//	iter = -1
						//	break
						//}
					} 
				}
				if iter == 1 {
					pur = k8 * c.Fck * c.B * dck + astot * (fsc + fst)/2.0
					if (pu == 0 && pur - pu < 0.1) || (xu == c.Cvrc + 1.0 && pur > pu) {
						iter = -1
						break
					}
					
					//if math.Abs(pur - pu) <= 0.03 * pu{
					//	iter = -1
					//	break
					//}
					if pur >= pu{
						if step >= xstep * 10.0{
							iter = -1
							break
						}
						if 1.02 * pu >= pur{
						//if (pur - pu < 0.02 * pu) || pur == pu {
							iter = -1
							break
						} else {
							xu -= c.H/step
							step = step * 10.0
							xu -= c.H/step
						} 
					}
					
					//if math.Abs(pur-pu) < 0.03 * pu && step >= xstep * 10.0{
					//	iter = -1
					//	break
					//}
					
				}
			}
			
			//fmt.Println(ColorYellow,astot,ColorReset)
			pur = k8 * c.Fck * c.B * dck + astot * (fsc + fst) /2.0
			mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * (fsc - fst) * (c.H/2.0 -c.Cvrc)/2.0
			c.Asc = astot/2.0; c.Ast = c.Asc
			pur = pur/1e3
			mur = mur/1e6
		case 1:
			//unsymmetrical arrangement o' reinforcement
			var ast, asc, astprev, ascprev float64
			effd = c.H - c.Cvrt
			dck = pu/c.B/k8/c.Fck
			if pu * (c.H - dck)/2.0 > mu {
				//get min. steel
				log.Println("providing minimum steel")
				c.Asc = 0.4 * c.B * c.H/2.0
				c.Ast = c.Asc
				return mu, pu, nil
			}
			if mu/pu < effd - c.H/2.0 {
				log.Println("switching to symmetrical section")
				c.Rtyp = 0
				return ColDBs(c, pu, mux, muy)
			}
			iter := 0; step = 100.0; ast = -1.0; asc = -1.0
			xu = c.Cvrc + 1.0
			xu -= c.H/step
			for iter != -1 {
				kiter++
				if kiter > 2999{
					return pur, mur, ErrIter
				}
				xu += c.H/step
				astprev = ast; ascprev = asc
				if xu > c.H/k9 {
					dck = c.H
				} else {
					dck = k9 * xu
				}
				esc = 0.0035 * (xu - c.Cvrc)/xu
				est = math.Abs(0.0035 * (effd - xu)/xu)
				fsc, fst = ColFstl(xu, effd, c.Fy, ey, esc, est)
				fsc = fsc - k8 * c.Fck
				asc = (mu + pu * (c.H/2.0 - c.Cvrt) - k8 * c.Fck * c.B * dck * (effd - dck/2.0))/fsc/(effd - c.Cvrc)
				ast = (pu - k8 * c.Fck * c.B * dck - asc* fsc)/fst
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
					}
				}
			}
			//DIS IS WRENG
			astot = ascprev + astprev
			pur = k8 * c.Fck * c.B * dck + astot * (fsc + fst) /2.0
			mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * (fsc - fst) * (c.H/2.0 -c.Cvrc)/2.0
			c.Asc = ascprev; c.Ast = astprev
			pur = pur/1e3
			mur = mur/1e6			
		}
	}
	c.Asteel = c.Asc + c.Ast
	//if c.Asteel < 4.0 * RbrArea(12.0){
	//	c.Asteel = 4.0 * RbrArea(12.0); c.Asc = c.Asteel/2.0; c.Ast = c.Asteel/2.0
	//}
	return pur, mur, nil
}

//ColAzBs analyzes a rectangular 2 layer column section for pur and returns ult moment mur
//ONLY WORKS FOR REKT 2-layer SECTION
//ALL STYPS ARE 1 NEVER USE THIS (except when one uses this)
func ColAzBs(c *RccCol,pu float64) (float64, error) {
	var k8, k9, ey, est, esc, fst, fsc, pur, dck, xu, step, mur float64
	if pu == 0.0{
		return mur, errors.New("zero axial load")
	}
	//bs 8110 stress block factors
	k8 = 0.45; k9 = 0.9 
	ey = c.Fy/1.15/200000.0
	pu = pu * 1e3
	effd := c.H - c.Cvrc
	//iterate to find neutral axis depth xu
	step = 10.0
	xu = c.Cvrc + 1.0
	xu = xu - c.H/step
	for math.Abs(pur - pu) >= 0.02 * pu {
		if pur < pu {
			xu += c.H/step
		} else {
			xu -= c.H/step
			step = 10.0 * step
			xu -= c.H/step
		}
		if xu > c.H/k9 {
			dck = c.H
		} else {
			dck = k9 * xu
		}
		esc = 0.0035 * (xu - c.Cvrc)/xu
		est = math.Abs(0.0035 * (effd - xu)/xu)
		if xu == effd {
			//tensile stress zero
			fst = 0.0
			//compressive in compression face
			if esc > ey {
				fsc = c.Fy/1.15
			} else {
				fsc = 2e5 * esc
			}	
		} else if xu < effd {		
			//tensile stress in tension face steel
			if est > ey {
				fst = -c.Fy/1.15
			} else {
				fst = -2e5 * est
			}
			//compressive stress in compression face steel
			if esc > ey {
				fsc = c.Fy/1.15
			} else {
				fsc = 2e5 * esc
			}
		} else {
			//compressive stress in both steelz
			if esc > ey {
				fsc = c.Fy/1.15
			} else {
				fsc = 2e5 * esc
			}
			if est > ey {
				fst = c.Fy/1.15
			} else {
				fst = 2e5 * est
			}
		}
		pur = k8 * c.Fck * c.B * dck + c.Asc * fsc + c.Ast * fst
		mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + c.Asc * fsc * (c.H/2.0 - c.Cvrc) - c.Ast * fst * (effd - c.H/2.0) 
		//in kN m
		mur = mur/1e6 
		if xu > effd && fsc >= c.Fy/1.15 && fst >= c.Fy/1.15 && pur < pu {
			return mur, ErrCAxL
		}
		if mur <= 0.0 {
			return mur, ErrCAxL
		}
		if xu == c.Cvrc + 1.0 && pur > pu {
			return mur, nil
		}
	}
	return mur, nil
}

//RbrFrcBs returns stress in rebar at a strain ebar
func RbrFrcBs(fy, ebar float64) (fbar float64){
	//returns rebar force at a strain ebar
	ey := fy/1.15/2e5
	if ebar > ey {
		fbar = fy/1.15
	} else {
		fbar = 2e5 * ebar
	}
	return
}

//ColFstl returns fsc, fst for top and bottom column rebar layers
//at a neutral axis depth xu
//and cleans up ColD somewhat
func ColFstl(xu, effd, fy, ey, esc, est float64) (fsc, fst float64) {
	if xu == effd {
		fst = 0.0
		if esc > ey {
			fsc = fy/1.15
		} else {
			fsc = 2e5 * esc
		}	
	} else if xu < effd {
		if est > ey {
			fst = -fy/1.15
		} else {
			fst = -2e5 * est
		}
		if esc > ey {
			fsc = fy/1.15
		} else {
			fsc = 2e5 * esc
		}
	} else {
		if esc > ey {
			fsc = fy/1.15
		} else {
			fsc = 2e5 * esc
		}
		if est > ey {
			fst = fy/1.15
		} else {
			fst = 2e5 * est
		}
	}
	return
}

//ColNMBs generates n (p/axial load) - m (moment) interaction curves
//(actually m-n) c.Nmplot - txt plot of mn curve 
//hulse section 5.4
func ColNMBs(c *RccCol) (pus, mus []float64){
	var k8, k9, ey, ebar, fbar, diabar, mu, pu, dck, puprev float64
	k8 = 0.45; k9 = 0.9
	ey = c.Fy/1.15/2e5
	fsts := make([]float64, len(c.Dbars))
	var data string
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
					if ebar >= ey {
						fbar = c.Fy/1.15
					} else {
						fbar = 2e5 * ebar
					}	
				} else {
					if ebar >= ey {
						fbar = -c.Fy/1.15
					} else {
						fbar = -2e5 * ebar
					}	
				}
				fsts[idx] = fbar
				pu += RbrArea(diabar)*fbar
			}
			mu = pu * c.H/2.0 - k8 * c.Fck * c.B * math.Pow(dck,2)/2.0
			for idx, fbar := range fsts {
				mu -= RbrArea(c.Dias[idx])*fbar*c.Dbars[idx]
			}
			if pu <= 0 {
				continue
			}
			if mu < 0.1 || math.Abs(pu - puprev) < 0.01 {
				break
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
					if ebar >= ey {
						fbar = c.Fy/1.15
					} else {
						fbar = 2e5 * ebar
					}	
				} else {
					if ebar >= ey {
						fbar = -c.Fy/1.15
					} else {
						fbar = -2e5 * ebar
					}	
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
	//make this the way for generic point and line plots
	skript := "d2.gp"
	term := c.Term
	folder := ""
	if c.Web{folder="web"}
	switch term{
		case "svg","svgmono":
		default:
		term = "svg"
	}
	if c.Title == ""{
		c.Title = fmt.Sprintf("rcc-col-%v-az",rand.Intn(666))
	}
	
	fname := c.Title + "-m-n-curve"
	//Draw(data, skript, term, folder, fname, title, xl, yl, zl string) (txtplot string, err error)
	dstr, _ := draw.Draw(data, skript, term, folder, fname, "M-N interaction curve", "BM(kn.m)","N(kn)","")
	c.Nmplot = dstr
	
	c.Txtplots = append(c.Txtplots, dstr)
	return
}

//ColFstBs returns fsc, fst for column (bs code, hulse)
//WHY ARE THERE FIVE OF THESE why indeed
//subtracts fck from fstl for steel displaced area NO IT DOESN'T (yet)
func ColFstlBs(xu, effd, fck, fy, esc, est float64) (fsc, fst float64) {
	ey := fy/1.15/2e5
	if esc >= ey{
		fsc = fy/1.15
	} else {
		fsc = 2e5 * esc
	}
	if xu == effd{
		fst = 0.0
		return
	}
	if xu > effd{
		if est >= ey {
			fst = fy/1.15
		} else {
			fst = 2e5 * est
		}
	} else {
		if est >= ey{
			fst = -fy/1.15
		} else {
			fst = -2e5 * est
		}
	}
	return
}

//ColEffHt rips off hulse 5.7 for column effective height calculation
//col dims cbs, cds, cls  -> lo hi
//beam dims bbs, bds, bls -> lol lor hil hir

func ColEffHt(braced, ljbase, basemr, ljbeamss, ujbeamss bool, cbs, cds, cls, bbs, bds, bls []float64, l0 float64) (float64, bool, error){
	var c1, c2, cmin, l1, l2, lmin, tl, tm float64
	var slender bool
	switch {
	case ljbase:
		//lower joint = base support
		c1 = 10.0 //s.s
		if basemr{
			//moment resisting base
			c1 = 1.0
		}
	case ljbeamss:
		//lower joint beams simply supported
		c1 = 10.0
	default:
		//sum of col stiffness/sum of beam stiffness at lower joint
		kcoll := cbs[0] * math.Pow(cds[0],3)/12.0/cls[0] + cbs[1] * math.Pow(cds[1],3)/12.0/cls[1]
		kbl := bbs[0] * math.Pow(bds[0],3)/12.0/bls[0] + bbs[1] * math.Pow(bds[1],3)/12.0/bls[1]
		c1 = kcoll/kbl
	}
	switch {
	case ujbeamss:
		//upper joint beams simply supported
		c2 = 10.0
	default:
		//sum of col stiffness/sum of beam stiffness at upper joint
		kcolu := cbs[0] * math.Pow(cds[0],3)/12.0/cls[0] + cbs[2] * math.Pow(cds[2],3)/12.0/cls[2]
		kbu := bbs[2] * math.Pow(bds[2],3)/12.0/bls[2] + bbs[3] * math.Pow(bds[3],3)/12.0/bls[3]
		c2 = kcolu/kbu
	}
	cmin = c1
	if cmin > c2 {
		cmin = c2
	}
	switch {
	case braced:
		l1 = l0 * (0.7 + 0.05 * (c1 + c2))
		l2 = l0 * (0.85 + 0.05 * cmin)
	default:
		l1 = l0 * (1.0 + 0.15 * (c1 + c2))
		l2 = l0 * (2.0 + 0.3 * cmin)
	}
	lmin = l1
	if lmin > l2 {lmin = l2}
	if braced  && lmin > l0 {lmin = l0}
	tm = cbs[0]; tl = cds[0]
	if tm > cds[0] {
		tm = cds[0]
	}
	if tl < cds[0] {
		tl = cds[0]
	}
	if l0 > 60.0 * tm {
		log.Println("ERRORE,errore-> clear height of column > 60 * minimum section dimension")
		return lmin, slender, ErrCEffHt
	}
	if !braced && l0 > 100.0 * tm * tm/tl {
		log.Println("ERRORE,errore-> clear height of column is excessive if one end is unrestrained")
		return lmin, slender, ErrCEffHt
	}
	switch {
	case braced:
		slender = lmin/cds[0] > 15.0
	default:
		slender = lmin/cds[0] > 10.0
	}
	return lmin, slender, nil
}

//ColEffHt was meant to calc column effective heights in x and y
func (c *RccCol) EffHt() (err error){
	//updates colum effective height lex and ley
	//stiffness = ei/l; e const if grades are ze same or k = sqrt(e) * i /l
	return
}


/*
func ColDBs(c *RccCol, pu, mu float64) (float64, float64, error){
	var k8, k9, ey, est, esc, fst, fsc, pur, dck, xu, step, mur, ecct, astot, effd float64
	k8 = 0.45; k9 = 0.9
	ey = c.Fy/1.15/2e5
	pu = pu * 1e3; mu = mu * 1e6
	ecct = c.H/20.0
	if ecct > 20.0 {ecct = 20.0}
	if mu < pu * ecct {mu = pu * ecct}

}
*/
/*
	switch c.Styp {
	case 0:
		//rectangular
		//fmt.Println("RECT SECKT DESIGN")
		switch c.Rtyp {
		case 0:
			//symmetrical arrangement of reinforcement
			effd = c.H - c.Cvrt
			dck = pu/c.B/k8/c.Fck
			if pu * (c.H - dck)/2.0 > mu {
				//get min. steel
				c.Asc = 0.4 * c.B * c.H/2.0
				c.Ast = c.Asc
				return mu, pu, nil
			}
			xu = 2.333 * effd; step = 10.0
			iter := 0
			xu += c.H/step
			for iter != -1 {
				if xu <= 0.9 * effd {
					iter = 1
					xu = c.Cvrc + 1.0
					step = 10.0
					xu -= c.H/step
				}
				if xu > c.H/k9 {
					dck = c.H
				} else {
					dck = k9 * xu
				}
				switch iter {
				case 0:
					xu -= c.H/step
				case 1:
					xu += c.H/step
				}
				esc = 0.0035 * (xu - c.Cvrc)/xu
				est = math.Abs(0.0035 * (effd - xu)/xu)
				fsc, fst = ColFstl(xu, effd, c.Fy, ey, esc, est)
				if iter == 0 {
					astot = 2.0 * (pu - k8 * c.Fck * c.B * dck)/(fsc + fst)
				}
				if iter == 1 {
					astot = 2.0 * (mu - k8 * c.Fck * c.B * dck * (c.H - dck)/2.0)/(fsc - fst)/(c.H/2.0 - c.Cvrc)
				}
				if astot < 0.0 {
					continue
				}
				if iter == 0 {
					mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * (fsc - fst) * (c.H/2.0 - c.Cvrc)/2.0
					if mur >= mu {
						if 1.02 * mu >= mu {
							iter = -1
							break
						} else {
							xu += c.H/step 
							step = 10.0*step
							xu += c.H/step
						}
					}
				}
				if iter == 1 {
					pur = k8 * c.Fck * c.B * dck + astot * (fsc + fst)/2.0
					if (pu == 0 && pur - pu < 0.1) || (xu == c.Cvrc + 1.0 && pur > pu) {
						iter = -1
						break
					}
					//mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * (fsc - fst) * (c.H/2.0 -c.Cvrc)/2.0
					if pur >= pu {
						if 1.02 * pu >= pur {
							iter = -1
							break
						} else {
							xu -= c.H/step
							step = 10.0 * step
							xu -= c.H/step
						}
						if step > 1e5 {iter = -1}
					}
					
				}
			}
			pur = k8 * c.Fck * c.B * dck + astot * (fsc + fst) /2.0
			mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * (fsc - fst) * (c.H/2.0 -c.Cvrc)/2.0
			c.Asc = astot/2.0; c.Ast = c.Asc
			pur = pur/1e3
			mur = mur/1e6
		case 1:
			//unsymmetrical arrangement o' reinforcement
			var ast, asc, astprev, ascprev float64
			effd = c.H - c.Cvrt
			dck = pu/c.B/k8/c.Fck
			if pu * (c.H - dck)/2.0 > mu {
				//get min. steel
				fmt.Println("providing minimum steel")
				c.Asc = 0.4 * c.B * c.H/2.0
				c.Ast = c.Asc
				return mu, pu, nil
			}
			if mu/pu < effd - c.H/2.0 {
				fmt.Println("switching to symmetrical section")
				c.Rtyp = 0
				return ColDBs(c, pu, mu, muy)
			}
			iter := 0; step = 100.0; ast = -1.0; asc = -1.0
			xu = c.Cvrc + 1.0
			xu -= c.H/step
			for iter != -1 {
				xu += c.H/step
				astprev = ast; ascprev = asc
				if xu > c.H/k9 {
					dck = c.H
				} else {
					dck = k9 * xu
				}
				esc = 0.0035 * (xu - c.Cvrc)/xu
				est = math.Abs(0.0035 * (effd - xu)/xu)
				fsc, fst = ColFstl(xu, effd, c.Fy, ey, esc, est)
				fsc = fsc - k8 * c.Fck
				asc = (mu + pu * (c.H/2.0 - c.Cvrt) - k8 * c.Fck * c.B * dck * (effd - dck/2.0))/fsc/(effd - c.Cvrc)
				ast = (pu - k8 * c.Fck * c.B * dck - asc* fsc)/fst
				if (ast < 0.0 || astprev < 0.0 || asc < 0.0 || ascprev < 0.0) {
					continue
				}
				if asc + ast >= ascprev + astprev {
					if iter == 0 {
						xu -= 2.0 * effd/step
						step = 1000.0
						iter = 1
						xu -= effd/step
					} else {
						iter = -1
					}
				}
			}
			//fmt.Println("asc ast", ascprev, astprev)
			c.Asc = ascprev; c.Ast = astprev
			//ALL THIS IS WRONG. SO WRONG.
			astot = ascprev + astprev
			pur = k8 * c.Fck * c.B * dck + astot * (fsc + fst) /2.0
			mur = k8 * c.Fck * c.B * dck * (c.H - dck)/2.0 + astot * (fsc - fst) * (c.H/2.0 -c.Cvrc)/2.0
			pur = pur/1e3
			mur = mur/1e6			
		}
	}
	return pur, mur, nil

*/
