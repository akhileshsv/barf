package barf

//column rebar gen routines
//goddamn mess this whole file is

import (
	//"log"
	//"fmt"
	"errors"
	"sort"
	"math"
	kass"barf/kass"
	//"math/rand"
)

//GetLinkSec offsets column section by the average of cvrc and cvrt
//to generate c.Lsec
func (c *RccCol) GetLinkSec(){
	//equal offset is generally (-_-)7 wrong
	//offst := (c.Cvrt + c.Cvrc)/2.0
	offst := 40.0
	//if offst < 
	//offst := 50.0 //2 inches
	st := kass.SecOffset(*c.Sec, offst,-1)
	c.Lsec = st
	return
}

//BarDbars generates dbars for coltyp via splitsides (so far upto 3)
func (c *RccCol) BarDbars(){
	switch c.Rbrlvl{
		case -1:
		//split sides by max spacing (30 mm)
		c.Dbars = []float64{}
		c.GetLinkSec()
		s1 := c.Lsec.Splitmax(75.0)
		c.Lsec = s1
		nc := s1.Ncs[0]
		c.Barpts = s1.Coords[0:nc-1]
		for _, pt := range s1.Coords[0:nc-1]{
			c.Dbars = append(c.Dbars, math.Abs(c.Sec.Ymx - pt[1])) 
		}
		c.Dias = make([]float64, len(c.Dbars))

		case 0:		
		c.Dbars = []float64{}
		c.GetLinkSec()
		//if tis a hollow section IGNORE TEH inner box?
		nc := c.Lsec.Ncs[0]
		c.Barpts = c.Lsec.Coords[0:nc-1]
		//now these are level ZERO dbars
		for _, pt := range c.Lsec.Coords[0:nc-1]{
			c.Dbars = append(c.Dbars, math.Abs(c.Sec.Ymx - pt[1])) 
		}
		//c.Rbrlvl = 1
		c.Dias = make([]float64, len(c.Dbars))

		case 1:
		//get midpoint of sides of link sec
		c.Dbars = []float64{}
		if len(c.Lsec.Ncs) == 0 {c.GetLinkSec()}
		s1 := c.Lsec.SplitSides(75.0)
		c.Lsec = s1
		nc := s1.Ncs[0]
		c.Barpts = s1.Coords[0:nc-1]
		//now these are level une dbars
		for _, pt := range s1.Coords[0:nc-1]{
			c.Dbars = append(c.Dbars, math.Abs(c.Sec.Ymx - pt[1])) 
		}
		//c.Rbrlvl = 2
		c.Dias = make([]float64, len(c.Dbars))

		case 2:
		//split lsec again? lol
		c.Dbars = []float64{}
		if len(c.Lsec.Ncs) == 0 {c.GetLinkSec()}
		s1 := c.Lsec.SplitSides(75.0)
		
		c.Lsec = s1
		s2 := c.Lsec.SplitSides(75.0)
		nc := s2.Ncs[0]
		c.Barpts = s2.Coords[0:nc-1]
		//now these are level two dbars
		for _, pt := range s2.Coords[0:nc-1]{
			c.Dbars = append(c.Dbars, math.Abs(c.Sec.Ymx - pt[1])) 
		}
		//c.Rbrlvl = 3
		c.Dias = make([]float64, len(c.Dbars))
	}
	return
}

//BarGen generates column bars for c.Asteel 
func (c *RccCol) BarGen() (err error){
	//use for everything?
	//switch c.Styp{
	//	case 1:
	//use ColRbrGen
	//	default:
	switch c.Rtyp{
		case 2:
		if len(c.Dbars) == 0{c.BarDbars()}
		//get abar = asteel/len(c.Dbars); get dia, add to dias
		var diamain float64
		astot := c.Asteel; abar := astot/float64(len(c.Dbars))
		//see table 4.1 of allen for equiv bar sizes
		cdias := []float64{12.0,16.0,18.0,20.0,22.0,25.0,28.0,32.0,35.0,40.0,43.0,45.0,50.0,55.0,57.0,64.0,69.0}
		for _, dia := range cdias{
			adia := RbrArea(dia)
			
			if adia >= abar || math.Abs(adia - abar) < 1.0{
				for i := range c.Dbars{
					c.Dias[i] = dia
				}
				diamain = dia
				break
			}
		}
		//fmt.Println("bargen dia",ColorGreen, diamain, RbrArea(diamain),ColorReset)
		nmain := float64(len(c.Dbars))
		dialink := getdiatie(diamain)
		
		efcvr := 30.0 + dialink + diamain/2.0
		if efcvr < c.Cvrt{efcvr = c.Cvrt}
		astprov := RbrArea(diamain) * nmain
		opt := 	[]float64{diamain, 0.0, nmain, 0.0, astprov, c.Asteel,  astprov - c.Asteel,efcvr,dialink}
		c.Rbropt = [][]float64{opt}
		c.Asteel = astprov; c.Psteel = astprov * 100.0/c.Sec.Prop.Area
		c.Dtie = dialink
		c.D1 = diamain
		c.D2 = diamain
		c.TiePitch()
		c.Asteel = astprov; c.Psteel = astprov * 100.0/c.Sec.Prop.Area
		//fmt.Println("bargen",ColorRed, c.Asteel, ColorReset)
		default:
		//calc width at cvrc, cvrt
		c.GetLinkSec()
		if c.Cvrt == 0.0{c.Cvrt = c.C1}
		if c.Cvrt == 0.0{c.Cvrt = 50.0}
		yt := c.Sec.Ymx - c.Sec.Ym - c.Cvrt
		yc := c.Cvrc
		if yc == 0.0{yc = c.Cvrt}
		var diamax, diamin, astprov float64
		diamin = 40.0
		for i, dy := range []float64{yc, yt}{				
			var layr [][]float64
			_, dx, pts := c.Sec.GetWidth(dy)
			var ast float64
			switch len(dx){
				case 0:
				err = errors.New("error in section width calc")
				case 1:
				switch i{
					case 0:
					ast = c.Asc
					if ast < 225.0{ast = 225.0}
					case 1:
					ast = c.Ast
					if ast < 225.0{ast = 225.0}
				}
				//fmt.Println("i, ast",i,ast)
				rez, e := BmBarCombo(6, ast)
				if e != nil{err = e; return}
				
				for _, r := range rez{
					rlay := BarNRows(dx[0], (c.Sec.Ymx - c.Sec.Ym)/2.0, r)
					layr = append(layr, rlay)
					sort.Slice(layr, func(i, j int) bool {
						if layr[i][7] != layr[j][7]{
							return layr[i][7] < layr[j][7]
						}
						return layr[i][4] < layr[j][4]
					})
				}
				switch i{
					case 0:
					c.Rbrc = layr[0]; c.Rbrcopt = layr
					astprov += c.Rbrc[4]
					//log.Println("kompress\n",c.Rbrc)
					e := c.BarLay(dx[0], pts, 1)
					if e != nil{err = e; return}
					if layr[0][2] >= diamax{
						diamax = layr[0][2]
					}
					if layr[0][3] >= diamax{
						diamax = layr[0][3]
					}
					
					if layr[0][2] <= diamin{
						diamin = layr[0][2]
					}
					if layr[0][3] <= diamin{
						diamin = layr[0][3]
					}
					case 1:
					c.Rbrt = layr[0]; c.Rbrtopt = layr
					astprov += c.Rbrt[4]
					//log.Println("tenz\n",c.Rbrt)
					//fmt.Println("tens\n",c.Rbrt,"pts\n",pts)
					//fmt.Println("calling barlay->",dx[0])
					e := c.BarLay(dx[0], pts, 2)
					if e != nil{err = e; return}
					
					if layr[0][2] >= diamax{
						diamax = layr[0][2]
					}
					if layr[0][3] >= diamax{
						diamax = layr[0][3]
					}
					if layr[0][2] <= diamin{
						diamin = layr[0][2]
					}
					if layr[0][3] <= diamin{
						diamin = layr[0][3]
					}
				}
				case 2:
				//two flanges at dy
			}
		}
		//get dbars for x
		nbars := len(c.Barpts)
		c.Dbars = make([]float64, nbars)
		//since asc is called first, invert
		for i, pt := range c.Barpts{
			c.Dbars[nbars-i-1] = pt[1]
		}
		c.Dtie = getdiatie(diamax)
		c.D2 = diamin
		c.D1 = diamax
		c.TiePitch()
		c.Asteel = astprov; c.Psteel = astprov * 100.0/c.Sec.Prop.Area
	}
	
	return
}

//DBarpts gets c.Dbars (distance from top fiber) from c.Barpts (rebar coords)
func (c *RccCol) DBarpts(){
	c.Dbars = make([]float64, len(c.Barpts))
	for i, pt := range c.Barpts{
		c.Dbars[i] = pt[1]
	}
	return
}

//ColRbrGen generates rebar for a (RECT) column func 
//details a (RECT/CIRCULAR) column with c.Asc and c.Ast, c.Rtype and c.Sectype, c.Nlayers
//returns dia, depths
func ColRbrGen(c *RccCol) (err error){ //opts, dxs, dys [][]float64, 
	//var amin float64
	var opts, dxs, dys [][]float64
	if c.Nlayers == 0{c.Nlayers = 2}
	
	switch c.Styp{
		case 1:
		c.GetLinkSec()
		switch{
			case c.Rtyp == 0 && c.Nlayers == 2:
			
			err = c.BarGen()
			return
			case c.Rtyp == 0 && c.Nlayers > 2:
			opts, dxs, dys, err = ColRbrSelect(c)
			if err != nil{return}
			mindia := opts[0][0]
			astprov := opts[0][4]
			c.Asteel = astprov
			c.Psteel = 100.0 * astprov/c.Sec.Prop.Area
			//c.Asteel = opts[0][1]
			ycs := make(map[float64]bool)
			c.Dtie = opts[0][8]
			c.D1 = mindia; c.D2 = mindia
			c.TiePitch()
			//c.D1 = mindia
			for i := 0; i < c.Nlayers * 2; i++{
				c.Dias = append(c.Dias, mindia)
				c.Dbars = append(c.Dbars, dxs[0][i])
				c.Dybars = append(c.Dybars, dys[0][i])
				if _, ok := ycs[dxs[0][i]]; !ok{
					ycs[dxs[0][i]] = true
				}
			}
			efcvr := opts[0][7]
			for yc := range ycs{
				x1 := efcvr; x2 := c.B - efcvr
				c.Barpts = append(c.Barpts, []float64{x1, yc, mindia})
				c.Barpts = append(c.Barpts, []float64{x2, yc, mindia})	
			}
			c.Rbr = opts[0]
			c.Rbropt = opts
			
			default:
			//now use colbar gen type shit
			//THIS IS ALL WORTHLESS and possibly very confusing to see
			switch c.Rtyp{
				case 0:	
				
				opts, dxs, dys, err = ColRbrSelect(c)
				if err != nil{return}
				mindia := opts[0][0]
				astprov := opts[0][4]
				c.Asteel = astprov
				c.Psteel = 100.0 * astprov/c.Sec.Prop.Area
				//c.Asteel = opts[0][1]
				ycs := make(map[float64]bool)
				c.Dtie = opts[0][8]
				c.D1 = mindia; c.D2 = mindia
				c.TiePitch()
				//c.D1 = mindia
				for i := 0; i < c.Nlayers * 2; i++{
					c.Dias = append(c.Dias, mindia)
					c.Dbars = append(c.Dbars, dxs[0][i])
					c.Dybars = append(c.Dybars, dys[0][i])
					if _, ok := ycs[dxs[0][i]]; !ok{
						ycs[dxs[0][i]] = true
					}
				}
				efcvr := opts[0][4]
				for yc := range ycs{
					x1 := efcvr; x2 := c.B - efcvr
					c.Barpts = append(c.Barpts, []float64{x1, yc, mindia})
					c.Barpts = append(c.Barpts, []float64{x2, yc, mindia})	
				}
				c.Rbr = opts[0]
				c.Rbropt = opts
				
				case 1:
				//unsymm
				//get ast layer, get asc layer
				err = c.BarGen()
				return
			}
		}
	}
	return 
}

//ColBBxRbrGen is for biaxially bent column rebar generation 
func ColBxRbrGen(c *RccCol, astr float64) (rez[][]float64, midx int){
	mxlx := 2 + int((c.H - c.Cvrc - c.Cvrt)/75.0)
	mxly := 2 + int((c.B - c.Cvrc - c.Cvrt)/75.0)
	//astr = astr 
	var delta, deltamin, aprov float64
	deltamin = -6.6
	for i := mxlx; i >=2; i--{
		for j := 2; j <= mxly; j++{
			if len(rez) > 100{
				break
			}
			nbars := j * 2 + (i - 2) * 2; n2bars := nbars - 4
			dreq := 2.0 * math.Sqrt(astr/math.Pi/float64(nbars))
			for idx, dia := range cdias{
				if math.Abs(dia - dreq) < 6.0{
					aprov = RbrArea(dia) * float64(nbars); delta = math.Abs(aprov - astr)
					if delta < 200.0{						
						rez = append(rez, []float64{float64(i),float64(j),dia, float64(nbars), 0, 0, aprov, delta, float64(nbars)})
						//log.Println("found n0=>",dia, nbars, aprov, delta)
						if deltamin == -6.6{
							deltamin = delta
							midx = len(rez) - 1
						} else if deltamin > delta{
							deltamin = delta
							midx = len(rez) - 1
						}	
					}
				}
				if idx < len(cdias) -1 {
					for jdx := idx + 1; jdx < len(cdias) -1; jdx++{
						d1 := cdias[jdx]; a1 := 4.0 * RbrArea(d1)
						if a1 > astr{continue}
						dreq = 2.0 * math.Sqrt((astr - a1)/math.Pi/float64(n2bars))
						//if dreq < 12.0{continue}
						if math.Abs(dia - dreq) < 6.0{
							aprov = a1 + RbrArea(dia) * float64(n2bars); delta = math.Abs(aprov - astr)
							if delta < 200.0{
								rez = append(rez, []float64{float64(i),float64(j),float64(n2bars), dia, 4.0, d1, aprov, delta, float64(n2bars) + 4.0})
								if deltamin == -6.6{
									deltamin = delta
									midx = len(rez) - 1
								} else if deltamin > aprov - astr{
									deltamin = delta
									midx = len(rez) - 1
								}								
							}
						}
					}
				}
			}
		}
	}
	return rez, midx
}

//ColRbrDbarGen calculates column rebar placement from an n1 d1 combo
//THIS IS ACTUALLY USEFUL write this into ColDesign
func ColRbrDbarGen(c *RccCol, r []float64)(dias, dxs, dys []float64, barpts [][]float64, aprov float64){
	//r haz x layers, y layers, n2, d2, n1 =4 , d1 aprov, delta, nbars
	//calculate ACTUAL EFFECTIVE COVER WITH DIALINKS
	xls := int(r[0]); yls := int(r[1])
	//n2 := int(r[2])
	d2 := r[3]
	//n1 := int(r[4])
	d1 := r[5]
	aprov = r[6]
	if d1 > 0.0{
		for i := 0; i < 4; i++{
			dias = append(dias, d1)
		}
	} else {
		for i := 0; i < 4; i++{
			dias = append(dias, d2)
		}
	}
	var clrdx, clrdy float64
	//if d1 > 0.0{clrdx = c.H - c.Cvrc - c.Cvrt - d1} else {clrdx = c.H - c.Cvrc - c.Cvrt - d2}
	clrdx = c.H - c.Cvrc - c.Cvrt
	xstep := math.Round(clrdx/float64(xls - 1))
	dxs = append(dxs, []float64{c.Cvrc, c.Cvrc, c.H - c.Cvrt, c.H - c.Cvrt}...)
	dys = append(dys, []float64{c.Cvrc, c.Cvrc, c.B - c.Cvrt, c.B - c.Cvrt}...)
	if d1 > 0.0{
		barpts = append(barpts,[][]float64{{c.Cvrc, c.Cvrt, d1},{c.B - c.Cvrc, c.Cvrt, d1},{c.Cvrc, c.H - c.Cvrt,d1},{c.B - c.Cvrt, c.H - c.Cvrt, d1}}...)
	} else {
		barpts = append(barpts,[][]float64{{c.Cvrc, c.Cvrt, d2},{c.B - c.Cvrc, c.Cvrt, d2},{c.Cvrc, c.H - c.Cvrt,d2},{c.B - c.Cvrt, c.H - c.Cvrt, d2}}...)
	}
	dxbar := c.Cvrc
	if xls > 2{
		for i := 0; i < xls -2; i++{
			dxbar += xstep
			dxs = append(dxs, dxbar)
			dxs = append(dxs, dxbar)
			dys = append(dys, c.Cvrc)
			dys = append(dys, c.B - c.Cvrt)
			dias = append(dias, d2)
			dias = append(dias, d2)
			barpts = append(barpts, []float64{c.Cvrc, dxbar, d2})
			barpts = append(barpts, []float64{c.B - c.Cvrt, dxbar, d2})
		}
	}
	//if d1 > 0.0{clrdy = c.B - c.Cvrc - c.Cvrt - d1} else {clrdy = c.B - c.Cvrc - c.Cvrt - d2}
	clrdy = c.B - c.Cvrc - c.Cvrt
	ystep := math.Round(clrdy/float64(yls -1))
	dybar := c.Cvrc
	if yls > 2{
		for i := 0; i < yls -2; i++{
			dybar += ystep
			dias = append(dias, d2)
			dias = append(dias, d2)
			dys = append(dys, dybar)
			dys = append(dys, dybar)
			dxs = append(dxs, c.Cvrc)
			dxs = append(dxs, c.H - c.Cvrt)
			barpts = append(barpts, []float64{dybar,c.Cvrc,d1})
			barpts = append(barpts, []float64{dybar,c.H - c.Cvrc,d1})
		}
	}
	return
}


//getdiatie returns the min diameter of lateral ties (>dia/4)
//should be get tie die :-|
func getdiatie(dia float64) (dialink float64){
	dialink = dia/4.0
	
	switch {
		//case dialink < 6.0:
	//dialink = 6.0
	case dialink <= 8.0:
		dialink = 8.0
	case dialink <= 10.0:
		dialink = 10.0
	default:
		dialink = 10.0
	}
	if dia >= 32.0{dialink = 10.0}
	//if dia == 0.0 || dialink == 0.0{
	//	log.Println("abbe laude lag gaye")
	//}
	return
}

//TiePitch updates the pitch of lateral ties 
func (c *RccCol) TiePitch(){
	//log.Println("d1, d2",c.D1, c.D2)
	//get pitch of lateral ties
	c.Ptie = 300.0
	dmin := 300.0
	for _, dim := range c.Dims{
		if dim < dmin {
			dmin = dim
		}
	}
	if c.Ptie > dmin && dmin != 0{
		c.Ptie = dmin
	}
	if c.D2 != 0.0 && c.Ptie > 16.0 * c.D2{
		c.Ptie = math.Floor(16.0 * c.D2/10.0)*10.0
	}

	//if c.Ptie == 0.0{
	//	log.Println("pitch",c.Ptie,"id->",c.Id, "dlink->",)
	//}
	return
}

//ColRbrSelect is an old function for column rebar option generation
//OKAY THIS IS ONLY FOR A SINGLE DIA COMBO
//opts - slice of valid bars (len max 5)
//make opts to ["dia1(mm)","nos","dia2(mm)","nos","ast prov(mm2)","ast req(mm2)","diff(mm2)"]
//dets - [na, nb] converted to 4 + 2 * na + 2 * nb templates 
func ColRbrSelect(c *RccCol) (opts, dxs, dys [][]float64, err error){
	//var rcldis float64 = 75.0
	cdias := []float64{12.0,16.0,18.0,20.0,22.0,25.0,28.0,32.0,40.0}
	alayer := c.Asteel/float64(c.Nlayers)
	//log.Println(ColorRed,"asteel, nlayers, alayer",c.Asteel, alayer,ColorReset)
	clrd := c.H - c.Cvrt - c.Cvrc
	if clrd/float64(c.Nlayers-1) < 75.0{
		//log.Println(clrd, clrd/float64(c.Nlayers-1), c.Nlayers)
		err = errors.New("too many rebar layers")
		return
	}
	switch c.Nlayers{
		case 2:
		//get area per layer, use rbr n rows
		//fmt.Println(ColorRed,"herezzz",c.Asteel, c.Dias,ColorReset)
		err = c.BarGen()
		return
		default:		
		for _, dia := range cdias{
			if len(opts) == 5{
				return
			}
			albar := RbrArea(dia) * 2.0
			//fmt.Println("here,",dia, albar)
			//log.Println(dia,albar)
			if albar < alayer{
				continue
			}
			//fmt.Println("dia found->", dia, albar-alayer)
			dialink := getdiatie(dia)
			efcvr := 30.0 + dialink + dia/2.0
			nbars := c.Nlayers * 2
			//"dia1(mm)","nos","dia2(mm)","nos","ast prov (mm2)","ast req(mm2)","diff(mm2)"})
			opts = append(opts, []float64{dia, 0.0, float64(nbars), 0.0, albar*float64(c.Nlayers), c.Asteel, (albar - alayer)*float64(c.Nlayers),efcvr,dialink})
			
			
			clrdx := c.H - 2.0 * efcvr - dia
			var dxbars []float64
			step := math.Round(clrdx/float64(c.Nlayers - 1))
			//log.Println(step)
			dbar := efcvr 
			for i := 0; i < 2; i++{
				dxbars = append(dxbars, efcvr)
				dxbars = append(dxbars, c.H - efcvr)
			}
			//log.Println(dia, dxbars)
			if c.Nlayers > 2{
				for i := 0; i < c.Nlayers -2; i++{
					dbar += step
					dxbars = append(dxbars, dbar)
					dxbars = append(dxbars, dbar)
				}
			}
			var dybars []float64
			for i := 0; i < c.Nlayers; i++{
				dybars = append(dybars, efcvr)
				dybars = append(dybars, c.B - efcvr)
			}
			dxs = append(dxs, dxbars)
			dys = append(dys, dybars)
		}
	}
	if len(opts) == 0{err = ErrSpacing; return}
	err = nil
	return
}
