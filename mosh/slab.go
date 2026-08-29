package barf

import (
	"fmt"
	"time"
	"math"
	"strings"
	"github.com/olekukonko/tablewriter"
	kass"barf/kass"
)

//RccSlb is a struct to store rcc slab fields
//see chapter 6, shah (design of slabs)
type RccSlb struct {
	//x shorter span alwayze
	//ss critical section - 0.1Lx (IS IT)
	//types - 0 clvr, 1 one way, 2 two way, 3 one way continuous, 4 one way ribbed continuous 5 flat slab
	//styp for irregularly shaped slabs (add coords)
	//yield line method - a - Lx, b - Ly; I1, I2, I3, I4 
	Id, Type, Endc, Sectype, Code, Ns, Nl         int
	Dtyp                                          int
	Fck, Fy, Fyd, Nomcvr, Diamain, Diadist        float64
	DL, LL, DM, Lspan, Ibent                      float64
	Lx, Ly, Dused, Astm, Astd, Spcm, Spcd         float64
	Efcvr, Bsup, Murmax, Ldev                     float64
	Rumax,Ptreq,Astreq, Ptsup                     float64
	D1, D2, L1, L2, L3, L4                        float64
	Rmyx                                          float64
	Title                                         string
	Term, Typstr                                  string
	Foldr                                         string
	I1, I2, I3, I4                                float64
	S1, S2, S3, S4                                float64 //support lengths for steel
	Pa, Pb                                        float64
	Bw, Dw, Df, Bf, Rspc                          float64 //ribbed slab params
	Ec                                            []int
	Slbc, Sdir                                    int      
	Asts, Asds, Spcms, Spcds                      []float64 `json:",omitempty"`
	Astop, Vfill                                  float64 `json:",omitempty"` //for ribbed slab topping mesh and filler quant
	Diatop, Spctop                                float64 `json:",omitempty"`
	Swt                                           float64 `json:",omitempty"` //self weight
	Dias, Dists, Astr                             []float64 `json:",omitempty"`
	Lspans, BM                                    []float64 `json:",omitempty"`
	Nspans                                        int       `json:",omitempty"`
	Ldcalc                                        int       `json:",omitempty"`
	Devchk                                        bool      `json:",omitempty"`
	Mindia                                        float64   `json:",omitempty"`
	Rezmap                                        map[float64][][]float64 `json:",omitempty"`
	Txtplots, Files                               []string `json:",omitempty"`
	Report                                        string `json:",omitempty"`
	Dz                                            bool `json:",omitempty"`
	Styp                                          int `json:",omitempty"`
	Coords                                        [][]float64 `json:",omitempty"`
	BMs                                           [][]float64 `json:",omitempty"`
	Astrs, Astps, Astds                           [][]float64 `json:",omitempty"`
	Vspns, Spcspns, Diaspns                       [][]float64 `json:",omitempty"`
	Distspns, Sdspns                              [][]float64 `json:",omitempty"`
	Nribs                                         float64 `json:",omitempty"`    
	Clvrs                                         [][]float64 `json:",omitempty"`
	Bsups                                         []float64 `json:",omitempty"`
	Fiti, Verbose, Spancalc, Web                  bool `json:",omitempty"`
	Vtot,Vrcc,Wstl,Afw,Kost                       float64 `json:",omitempty"`
	Kunit                                         float64 `json:",omitempty"`
	Bmap                                          map[float64][][]float64 `json:",omitempty"`
	Bsum                                          map[float64]float64 `json:",omitempty"`
	Cvec                                          []float64 `json:",omitempty"` //cost ck, fw, steel
	Kostin                                        []float64 `json:",omitempty"`
	Bmloads                                       []float64 `json:",omitempty"` //loads on supporting beams
	Wdl                                           float64 `json:",omitempty"`
	Kck, Ky                                       []float64 `json:",omitempty"` //cost ck grades, y grades
}

//SlbDesign is the main entry func for slab design from flags/menu
func SlbDesign(s *RccSlb) (err error){
	//slab design entry func
	//this should have been da way from ze beginning
	//use strings whenever possible
	switch s.Typstr{
		case "clvr":
		s.Type = 1
		s.Endc = 0
		s.Lx = s.Lspan
		case "1w":
		s.Type = 1
		s.Endc = 1
		s.Lx = s.Lspan
		case "1wcs":
		s.Type = 1
		s.Endc = 2
		case "2w":
		s.Type = 2
		s.Endc = 10
		case "2wcs":
		s.Type = 2
	}
	switch s.Type{
		case 1:
		switch s.Endc{
			case 0,1:
			switch s.Code{
				case 1:
				err = SlbDIs(s)
				case 2:
				err = SlbDBs(s)
			}
			case 2:
			switch s.Dtyp{
				case 1:
				//coefficients
				err = CSlb1DepthCs(s)
				if err != nil{
					return
				}
				err = CSlb1Stl(s)
				case 0:
				//envelope

				err = CSlb1Depth(s)
				if err != nil{
					return
				}
				err = CSlb1Stl(s)
			}
		}
		case 2:
		switch s.Code{
			case 1:
			err = SlbDIs(s)
			case 2:
			err = SlbDBs(s)
		}
		case 3:
		//one way ribbed slab
		err, _ = RSlb1Chk(s)
		if err != nil{
			return
		}
		case 4:
		//waffle slab
		err, _ = RSlb2Chk(s)
	}
	if err == nil{
		switch s.Type{
			case 1,2:
			s.Quant()
			s.Table(s.Verbose)	
			case 3:
			case 4:
		}
	}
	s.Draw(s.Term)
	return
}

//Sectbars returns 2d reinf data for plotting a slab section in gnuplot
func (s *RccSlb) Sectbars(sect string) (sdata, adata, ldata, cdata string){
	//returns bar coords for side view
	var xs, xi, xe, yb, yt float64
	switch s.Type{
		case 1:
		switch s.Endc{
			case 0,1:
			//clvr, ss one-way slab
			//clvr IS NOT COMPLETE IN ANY SHAPE OR FORM
			yb = s.Efcvr
			yt = s.Dused - s.Efcvr
			xi = xs - s.Bsup/2.0 + s.Efcvr
			xe = xs + s.Lspan + s.Bsup/2.0 - s.Efcvr
			switch s.Endc{
				case 0:
				xe = xs + s.Lspan - s.Efcvr
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yt, xe - xi, 0.0, 2.0)
				
				//extend top bar to ldev from the face of the support
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yt, -s.Ldev - s.Efcvr, 0.0, 7.0)
				
				ldata += fmt.Sprintf("%f %f M-T%.f-%.fc-c\n",(xi+xe)/2.0,-s.Dused*1.1,s.Diamain,s.Spcm)
				ldata += fmt.Sprintf("%f %f D-T%.f-%.fc-c\n",(xi+xe)/2.0,-s.Dused*1.0,s.Diadist,s.Spcd)
				
				// //label distance
				// adata += fmt.Sprintf("%f %f %f %f %f\n",xs, s.Dused*1.1, xe - xs, 0.0, 3.0)
				// ldata += fmt.Sprintf("%f %f %.fmm\n",(xi+xe)/2.0,s.Dused*1.2,0.1*s.Lspan+s.Bsup/2.0)

				//label dev len
				adata += fmt.Sprintf("%f %f %f %f %f\n",xi, s.Dused*1.1, -s.Ldev - s.Efcvr, 0.0, 3.0)
				ldata += fmt.Sprintf("%f %f %.fmm\n",xi-(s.Ldev + s.Efcvr)/2.0, s.Dused*1.2, s.Ldev)
				
				
				case 1:
				//main bottom reinf
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yb, xe - xi, 0.0, 2.0)
				ldata += fmt.Sprintf("%f %f M-T%.f-%.fc-c\n",(xi+xe)/2.0,-s.Dused*1.1,s.Diamain,s.Spcm)
				ldata += fmt.Sprintf("%f %f D-T%.f-%.fc-c\n",(xi+xe)/2.0,-s.Dused*1.0,s.Diadist,s.Spcd)
				//top left reinf
				xe = 0.1 * s.Lspan + s.Bsup/2.0
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yt, xe - xi, 0.0, 3.0)
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yt, 0.0, -yt/2.0, 3.0)
				ldata += fmt.Sprintf("%f %f M-T%.f-%.fc-c\n",(xi+xe)/2.0,s.Dused*1.6,s.Diamain,s.Spcm*2.0)
				ldata += fmt.Sprintf("%f %f D-T%.f-%.fc-c\n",(xi+xe)/2.0,s.Dused*1.5,s.Diadist,s.Spcd)
				//label distance
				adata += fmt.Sprintf("%f %f %f %f %f\n",xs, s.Dused*1.1, xe - xs, 0.0, 3.0)
				ldata += fmt.Sprintf("%f %f %.fmm\n",(xi+xe)/2.0,s.Dused*1.2,0.1*s.Lspan+s.Bsup/2.0)
				//top right reinf
				xi = xs + 0.9 * s.Lspan - s.Bsup/2.0
				xe = xs + s.Lspan + s.Bsup/2.0 - s.Efcvr
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yt, xe - xi, 0.0, 3.0)
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xe, yt, 0.0, -yt/2.0, 3.0)
				ldata += fmt.Sprintf("%f %f M-T%.f-%.fc-c\n",(xi+xe)/2.0,s.Dused*1.6,s.Diamain,s.Spcm*2.0)
				ldata += fmt.Sprintf("%f %f D-T%.f-%.fc-c\n",(xi+xe)/2.0,s.Dused*1.5,s.Diadist,s.Spcd)
				//label distance
				adata += fmt.Sprintf("%f %f %f %f %f\n",xi, s.Dused*1.1, xe - xi - s.Bsup/2.0, 0.0, 3.0)
				ldata += fmt.Sprintf("%f %f %.fmm\n",(xi+xe)/2.0,s.Dused*1.2,0.1*s.Lspan+s.Bsup/2.0)
			}
			ybd := yb + s.Diamain/2.0 + s.Diadist/2.0
			ytd := yt - s.Diamain/2.0 - s.Diadist/2.0
			switch s.Endc{
				case 0:
				//top dist. steel
				xs -= s.Efcvr
				for x := xs; x <= xs + s.Lspan; x += s.Spcd{
					cdata += fmt.Sprintf("%f %f %f %f %f\n",x,ytd,s.Diadist/2.0, s.Diadist, s.Diadist)
				}
				case 1:	
				//bottom dist. steel
				for x := xs; x <= xs + s.Lspan; x += s.Spcd{
					cdata += fmt.Sprintf("%f %f %f %f %f\n",x,ybd,s.Diadist/2.0, s.Diadist, s.Diadist)
				}
				//top left dist steel
				for x := xs; x <= xs + 0.1 * s.Lspan + s.Bsup/2.0; x += s.Spcd{
					cdata += fmt.Sprintf("%f %f %f %f %f\n",x,ytd,s.Diadist/2.0,s.Diadist, s.Diadist)
				}
				//top right dist steel
				for x := xs + 0.9 * s.Lspan - s.Bsup/2.0; x <= xs + s.Lspan + s.Bsup/2.0; x += s.Spcd{
					cdata += fmt.Sprintf("%f %f %f %f %f\n",x,ytd,s.Diadist/2.0, s.Diadist, s.Diadist)
				}
			}
			case 2:
			//cs slab
			switch sect{
				case "y":
				fmt.Println(ColorRed,"#TODO# #TODO# #TODO#",ColorReset)
				return
			}
			for i, lspan := range s.Lspans{
				s1 := s.Spcspns[i][0]; s2 := s.Spcspns[i][1]; s3 := s.Spcspns[i][2]
				s4 := s.Sdspns[i][0]; s5 := s.Sdspns[i][1]; s6 := s.Sdspns[i][2]
				//fmt.Println("here in sectbars, lspan, i, s1, s2, s3, s4, s5, s6-",i, lspan, s1,s2,s3,s4,s5,s6)
				var l1, l3, x1, x2 float64
				switch i{
					case 0:
					l1 = 0.1 * lspan//; l2 = lspan*(1.0-0.15-0.25)
					l3 = (1.0 - 0.3) * lspan	
					x1 = 0.15 * lspan
					x2 = 0.75 * lspan
					
					case s.Nspans - 1:
					l1 = 0.3 * lspan//; l2 = lspan*(1.0-0.25-0.15);
					l3 = (1.0 - 0.1) * lspan
					x1 = 0.25 * lspan
					x2 = 0.85 * lspan
					default:
					l1 = 0.3 * lspan//; l2 = lspan*(1.0-0.25-0.25);
					l3 = (1.0 - 0.3) * lspan
					x1 = 0.25 * lspan
					x2 = 0.75 * lspan
				}
				yb = s.Efcvr
				yt = s.Dused - s.Efcvr
				//xi = xs - s.Bsup/2.0 + s.Efcvr
				//xe = xs + lspan + s.Bsup/2.0 - s.Efcvr	
				//main bottom reinf
				xi = xs + x1
				xe = xs + x2
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yb + 0.1*s.Efcvr, xe - xi, 0.0, 2.0)				
				// //label tension steel curtailment from left 
				// adata += fmt.Sprintf("%f %f %f %f %f\n",xs, -1.3*s.Dused, xi - xs, 0.0, 4.0)
				// ldata += fmt.Sprintf("%f %f %.f\n",(xs+xi)/2.0,-s.Dused*1.25,xi - xs)
				
				// //label tension steel curtailment from right 
				// adata += fmt.Sprintf("%f %f %f %f %f\n",xe, -1.3*s.Dused, xs + lspan - xe, 0.0, 3.0)
				// ldata += fmt.Sprintf("%f %f %.f\n",(xe+ xs + lspan - xe)/2.0,-s.Dused*1.25, xs + lspan - xe)
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xs, yb, lspan, 0.0, 7.0)
				switch i{
					case 0:
					ldata += fmt.Sprintf("%f %f M-T%.f-%.f\n",(xi+xe)/2.0,-s.Dused*1.1,s.Diamain,s2)
					ldata += fmt.Sprintf("%f %f D-T%.f-%.f\n",(xi+xe)/2.0,-s.Dused*1.0,s.Diadist,s5)
					default:
					ldata += fmt.Sprintf("%f %f %.f\n",(xi+xe)/2.0,-s.Dused*1.1,s2)
					ldata += fmt.Sprintf("%f %f %.f\n",(xi+xe)/2.0,-s.Dused*1.0,s5)

				}
				//label distance
				
				//top left reinf
				xi = xs
				xe = xs + l1 + s.Bsup/2.0
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yt, xe - xi, 0.0, 3.0)
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yt, 0.0, -yt/2.0, 3.0)
				switch i{
					case 0:
					ldata += fmt.Sprintf("%f %f M-T%.f-%.f\n",(xi+xe)/2.0,s.Dused*1.6,s.Diamain,s1)
					ldata += fmt.Sprintf("%f %f D-T%.f-%.f\n",(xi+xe)/2.0,s.Dused*1.5,s.Diadist,s4)
					default:
					ldata += fmt.Sprintf("%f %f %.f\n",(xi+xe)/2.0,s.Dused*1.6,s1)
					ldata += fmt.Sprintf("%f %f %.f\n",(xi+xe)/2.0,s.Dused*1.5,s4)
				}
				//label distance
				adata += fmt.Sprintf("%f %f %f %f %f\n",xs, s.Dused*1.1, xe - xs, 0.0, 3.0)
				ldata += fmt.Sprintf("%f %f %.f\n",(xi+xe)/2.0,s.Dused*1.2,xe - xs)
				//top right reinf
				xi = xs + l3 - s.Bsup/2.0
				xe = xs + lspan 
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yt, xe - xi, 0.0, 3.0)
				sdata += fmt.Sprintf("%f %f %f %f %f\n",xe, yt, 0.0, -yt/2.0, 3.0)
				ldata += fmt.Sprintf("%f %f %.f\n",(xi+xe)/2.0,s.Dused*1.6,s3)
				ldata += fmt.Sprintf("%f %f %.f\n",(xi+xe)/2.0,s.Dused*1.5,s6)
				//label distance
				adata += fmt.Sprintf("%f %f %f %f %f\n",xi, s.Dused*1.1, xe - xi - s.Bsup/2.0, 0.0, 3.0)
				ldata += fmt.Sprintf("%f %f %.f\n",(xi+xe)/2.0,s.Dused*1.2,xe-xi)
				
				ybd := yb + s.Diamain/2.0 + s.Diadist/2.0
				ytd := yt - s.Diamain/2.0 - s.Diadist/2.0
				//bottom dist. steel:
				for x := xs; x <= xs + lspan; x += s5{
					cdata += fmt.Sprintf("%f %f %f %f %f\n",x,ybd,s.Diadist/2.0, s.Diadist, s.Diadist)
				}
				//top left dist steel
				for x := xs; x <= xs + l1 + s.Bsup/2.0; x += s4{
					cdata += fmt.Sprintf("%f %f %f %f %f\n",x,ytd,s.Diadist/2.0,s.Diadist, s.Diadist)
				}
				// //top right dist steel
				for x := xs + l3 - s.Bsup/2.0; x <= xs + lspan; x += s6{
					cdata += fmt.Sprintf("%f %f %f %f %f\n",x,ytd,s.Diadist/2.0, s.Diadist, s.Diadist)
				}
				xs += lspan
			}
		}
		case 2:
		//two way slab
		//order - short (x) sup, span, long (y) sup, span
		//groups - d0, s0, d2, s2, d1, s1, d3, s3
		//yb, yt, ydb, ydt - main bot, top, dist top, bot levels
		var lspan,sl1, sl2, dsp, ssp, ddsp, sdsp, ds, ss, dds, sds, ydb, ydt float64
		switch sect{
			case "y":
			lspan = s.Ly
			ds = s.Dias[3]
			ss = s.Spcms[3]
			dds = s.Dias[1]
			sds = s.Spcms[1]
			
			dsp = s.Dias[2]
			ssp = s.Spcms[2]
			ddsp = s.Dias[0]
			sdsp = s.Spcms[0]
			default:
			lspan = s.Lx
			
			ds = s.Dias[1]
			ss = s.Spcms[1]
			dds = s.Dias[3]
			sds = s.Spcms[3]
			
			dsp = s.Dias[0]
			ssp = s.Spcms[0]
			ddsp = s.Dias[2]
			sdsp = s.Spcms[2]
		}
		xi = xs - s.Bsup/2.0 + s.Efcvr
		xe = xs + lspan + s.Bsup/2.0 - s.Efcvr
		switch sect{
			case "y":
			ydb = s.Efcvr
			ydt = s.Dused - s.Efcvr
			yb = s.Efcvr + s.Diamain/2.0 + s.Diadist/2.0
			yt = s.Dused - s.Efcvr - s.Diamain/2.0 - s.Diadist/2.0
			sl1 = s.S3
			sl2 = s.S4
			default:
			yb = s.Efcvr
			yt = s.Dused - s.Efcvr
			ydb = s.Efcvr + s.Diamain/2.0 + s.Diadist/2.0
			ydt = s.Dused - s.Efcvr - s.Diamain/2.0 - s.Diadist/2.0
			sl1 = s.S1
			sl2 = s.S2
		}
		//main span bottom reinf ds, ss at yb
		sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yb, xe - xi, 0.0, 2.0)
		ldata += fmt.Sprintf("%f %f M-T%.f-%.fc-c\n",(xi+xe)/2.0,-s.Dused*1.1,ds,ss)
		ldata += fmt.Sprintf("%f %f D-T%.f-%.fc-c\n",(xi+xe)/2.0,-s.Dused*1.0,dds,sds)

		//top left reinf
		xe = sl1 * lspan + s.Bsup/2.0
		sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yt, xe - xi, 0.0, 3.0)
		sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yt, 0.0, -yt/2.0, 3.0)
		
		ldata += fmt.Sprintf("%f %f M-T%.f-%.fc-c\n",(xi+xe)/2.0,s.Dused*1.6,dsp,ssp)
		ldata += fmt.Sprintf("%f %f D-T%.f-%.fc-c\n",(xi+xe)/2.0,s.Dused*1.5,ddsp,sdsp)
		//label distance
		adata += fmt.Sprintf("%f %f %f %f %f\n",xs, s.Dused*1.1, xe - xs, 0.0, 3.0)
		ldata += fmt.Sprintf("%f %f %.f\n",(xi+xe)/2.0,s.Dused*1.2,sl1*lspan+s.Bsup/2.0)
		//top right reinf
		xi = xs + (1.0 - sl2)*lspan - s.Bsup/2.0
		xe = xs + lspan + s.Bsup/2.0 - s.Efcvr
		sdata += fmt.Sprintf("%f %f %f %f %f\n",xi, yt, xe - xi, 0.0, 3.0)
		sdata += fmt.Sprintf("%f %f %f %f %f\n",xe, yt, 0.0, -yt/2.0, 3.0)
		ldata += fmt.Sprintf("%f %f M-T%.f-%.fc-c\n",(xi+xe)/2.0,s.Dused*1.6,s.Diamain,s.Spcm*2.0)
		ldata += fmt.Sprintf("%f %f D-T%.f-%.fc-c\n",(xi+xe)/2.0,s.Dused*1.5,s.Diadist,s.Spcd)
		//label distance
		adata += fmt.Sprintf("%f %f %f %f %f\n",xi, s.Dused*1.1, xe - xi -s.Bsup/2.0, 0.0, 3.0)
		ldata += fmt.Sprintf("%f %f %.f\n",(xi+xe)/2.0,s.Dused*1.2,sl2*lspan+s.Bsup/2.0)

		//bottom dist. steel
		for x := xs; x <= xs + lspan; x += ss{
			cdata += fmt.Sprintf("%f %f %f %f %f\n",x,ydb,s.Diadist/2.0, s.Diadist, s.Diadist)
		}
		//top left dist steel
		for x := xs; x <= xs + sl1 * lspan + s.Bsup/2.0; x += ssp{
			cdata += fmt.Sprintf("%f %f %f %f %f\n",x,ydt,s.Diadist/2.0,s.Diadist, s.Diadist)
		}
		//top right dist steel
		for x := xs + (1.0 - sl2) * lspan - s.Bsup/2.0; x <= xs + lspan + s.Bsup/2.0; x += ssp{
			cdata += fmt.Sprintf("%f %f %f %f %f\n",x,ydt,s.Diadist/2.0, s.Diadist, s.Diadist)
		}
	}
	return
}

//sectcoords returns section coords for a slab span
func (s *RccSlb) Sectcoords(sect string, idx... int) (data, adata, ldata string){
	switch s.Type{
		case 1:
		//one way slab
		switch s.Endc{
			case 0:
			//clvr
			var xs, xe float64
			var pb, pe, dims []float64
			pb = []float64{xs - s.Bsup/2.0, s.Dused/2.0}
			dims = []float64{s.Dused, s.Dused}
			pe = []float64{xs + s.Lspan, s.Dused/2.0}
			
			data += kass.DrawMem2d(1, 1, pb, pe, dims)
			//add left support with depth Dused * 2.0
			dims = []float64{0.0,s.Bsup}
			pb = []float64{xs, 0.0}
			pe = []float64{xs, -s.Dused*2.0}
			data += kass.DrawMem2d(8, 1, pb, pe, dims)
			
			//draw span dimensions
			yar := -s.Dused/2.0
			adata += fmt.Sprintf("%f %f %f %f 1.0\n",xs, yar, s.Lspan, 0.0)
			ldata += fmt.Sprintf("%f %f %.fmm %f 1.0\n",xs+s.Lspan/2.0, yar - s.Dused/4.0, s.Lspan, 0.0)

			//indicate slab depth (add s.Bsup/2.0 for arrow loc)
			xe = xs + s.Lspan + s.Bsup 
			adata += fmt.Sprintf("%f %f %f %f 1.0\n",xe, 0.0, 0.0, s.Dused)
			ldata += fmt.Sprintf("%f %f %.fmm 1.0\n",xe, s.Dused/2.0, s.Dused)
			
			//indicate effective cover
			ldata += fmt.Sprintf("%f %f cover-%.fmm 1.0\n",xs, s.Efcvr/2.0, s.Efcvr)

			case 1:
			//ss
			//first draw (classic) section along short span
			var xs, xe float64
			var pb, pe, dims []float64
			pb = []float64{xs - s.Bsup/2.0, s.Dused/2.0}
			dims = []float64{s.Dused, s.Dused}
			switch s.Endc{
				case 0:
				pe = []float64{xs + s.Lspan, s.Dused/2.0}
				case 1:
				pe = []float64{xs + s.Lspan + s.Bsup/2.0, s.Dused/2.0}
			}
			data += kass.DrawMem2d(1, 1, pb, pe, dims)
			//add supports with depth Dused * 2.0
			dims = []float64{0.0,s.Bsup}
			pb = []float64{xs, 0.0}
			pe = []float64{xs, -s.Dused*2.0}
			data += kass.DrawMem2d(8, 1, pb, pe, dims)

			pb = []float64{xs + s.Lspan, 0.0}
			pe = []float64{xs + s.Lspan, -s.Dused*2.0}
			data += kass.DrawMem2d(8, 1, pb, pe, dims)

			//draw span dimensions
			yar := -s.Dused/2.0
			adata += fmt.Sprintf("%f %f %f %f 1.0\n",xs, yar, s.Lspan, 0.0)
			ldata += fmt.Sprintf("%f %f %.fmm %f 1.0\n",xs+s.Lspan/2.0, yar - s.Dused/4.0, s.Lspan, 0.0)

			//indicate slab depth (add s.Bsup/2.0 for arrow loc)
			xe = xs + s.Lspan + s.Bsup 
			adata += fmt.Sprintf("%f %f %f %f 1.0\n",xe, 0.0, 0.0, s.Dused)
			ldata += fmt.Sprintf("%f %f %.fmm 1.0\n",xe, s.Dused/2.0, s.Dused)
			
			//indicate effective cover
			ldata += fmt.Sprintf("%f %f cover-%.fmm 1.0\n",xs, s.Efcvr/2.0, s.Efcvr)

			case 2:
			//cs
			if s.Bsup == 0.0{s.Bsup = 230.0}
			switch sect{
				case "y":
				fmt.Println(ColorRed,"#TODO# #TODO# #TODO#",ColorReset)
				return
			}
			var xs, xe float64
			var pb, pe, dims []float64
			for i, lspan := range s.Lspans{
				xe = xs + lspan
				pb = []float64{xs, s.Dused/2.0}
				pe = []float64{xe, s.Dused/2.0}
				dims = []float64{s.Dused, s.Dused}
				switch i{
					case 0:
					pb[0] -= s.Bsup/2.0
					//indicate effective cover
					ldata += fmt.Sprintf("%f %f cover-%.fmm 1.0\n",xs, s.Efcvr/2.0, s.Efcvr)
					case s.Nspans-1:
					pe[0] += s.Bsup/2.0
					
					//indicate slab depth (add s.Bsup/2.0 for arrow loc)
					xar := xs + lspan + s.Bsup + s.Dused
					adata += fmt.Sprintf("%f %f %f %f 1.0\n",xar, 0.0, 0.0, s.Dused)
					ldata += fmt.Sprintf("%f %f %.fmm 1.0\n",xar, s.Dused/2.0, s.Dused)
				}
				data += kass.DrawMem2d(i, 1, pb, pe, dims)
				//add supports with depth Dused * 0.25 here
				dims = []float64{0.0,s.Bsup}
				pb = []float64{xs, 0.0}
				pe = []float64{xs, -s.Dused*0.25}
				data += kass.DrawMem2d(8, 1, pb, pe, dims)

				pb = []float64{xe, 0.0}
				pe = []float64{xe, -s.Dused*0.25}
				data += kass.DrawMem2d(8, 1, pb, pe, dims)
				
				//draw span dimensions
				yar := -s.Dused/2.0
				adata += fmt.Sprintf("%f %f %f %f 1.0\n",xs, yar, lspan, 0.0)
				ldata += fmt.Sprintf("%f %f %.fmm %f 1.0\n",xs+lspan/2.0, yar - s.Dused/4.0, lspan, 0.0)
				xs += lspan
			}
		}
		case 2:
		//two way slab
		var xs, xe, lspan float64
		var pb, pe, dims []float64
		switch sect{
			case "y":
			lspan = s.Ly
			default:
			//is always x(TO ME, MY X-)
			lspan = s.Lx
		}
		pb = []float64{xs - s.Bsup/2.0, s.Dused/2.0}
		dims = []float64{s.Dused, s.Dused}
		pe = []float64{xs + lspan + s.Bsup/2.0, s.Dused/2.0}
		data += kass.DrawMem2d(1, 1, pb, pe, dims)
		//add supports with depth Dused * 2.0
		dims = []float64{0.0,s.Bsup}
		pb = []float64{xs, 0.0}
		pe = []float64{xs, -s.Dused*2.0}
		data += kass.DrawMem2d(8, 1, pb, pe, dims)

		pb = []float64{xs + lspan, 0.0}
		pe = []float64{xs + lspan, -s.Dused*2.0}
		data += kass.DrawMem2d(8, 1, pb, pe, dims)

		//draw span dimensions
		yar := -s.Dused/2.0
		adata += fmt.Sprintf("%f %f %f %f 1.0\n",xs, yar, lspan, 0.0)
		ldata += fmt.Sprintf("%f %f %.fmm %f 1.0\n",xs+lspan/2.0, yar - s.Dused/4.0, lspan, 0.0)

		//indicate slab depth (add s.Bsup/2.0 for arrow loc)
		xe = xs + lspan + s.Bsup 
		adata += fmt.Sprintf("%f %f %f %f 1.0\n",xe, 0.0, 0.0, s.Dused)
		ldata += fmt.Sprintf("%f %f %.fmm 1.0\n",xe, s.Dused/2.0, s.Dused)
		//indicate effective cover
		ldata += fmt.Sprintf("%f %f cover-%.fmm 1.0\n",xs, s.Efcvr/2.0, s.Efcvr)
	}
	return
}

//Draw plots a slab using good ol' gnuplot
func (s *RccSlb) Draw(term string) (pltstr string){
	
	var data, adata, ldata string
	//bsup := s.Bsup
	//if bsup == 0.0{bsup = 230.0}
	d1, a1, l1 := s.Sectcoords("x")
	d1 += "\n\n"
	data += d1; adata += a1; ldata += l1
	
	s1, a1, l1, c1 := s.Sectbars("x")
	adata += a1; ldata += l1
	s1 += "\n\n"
	adata += "\n\n"; ldata += "\n\n"
	data += s1
	data += adata; data += ldata; data += c1
	
	
	var fname, fn string
	switch term{
		case "svg","svgmono":
		fn = fmt.Sprintf("%s_section_%v.svg",s.Title,s.Id)
		fname = genfname("",fn)
		case "dxf":
		fn = fmt.Sprintf("%s_section_%v.dxf",s.Title,s.Id)
		fname = genfname("",fn)
		case "png":
		fn = fmt.Sprintf("%s_section_%v.svg",s.Title,s.Id)
		fname = genfname("",fn)
	}
	if term != ""{
		pltstr = skriptrun(data, "plotslab.gp", term, s.Title, fname)
		switch term{
			case "dumb", "mono":
			fmt.Println(pltstr)
			case "svgmono", "dxf","svg":
			pltstr = fname
			//fmt.Println(pltstr)
		}
		if s.Web{
			//embed kongtext font in svg
			kass.Svgkong(pltstr)
			s.Txtplots = append(s.Txtplots, fn)
		} else {
			s.Txtplots = append(s.Txtplots, pltstr)
		}
	} 
	return
}

//Quant takes off quantities for an rcc slab
func (s *RccSlb) Quant(){
	//quantity take off
	//todo - add bbs - map[bar dia][length, type, total nos]
	//bbs types for slab - 1 - straight, 2 - l (top steel end spans) ???

	s.Bmap = make(map[float64][][]float64); s.Bsum = make(map[float64]float64)
	ld1 := BarDevLen(s.Fck, s.Fy, s.Diamain)
	ld2 := BarDevLen(s.Fck, s.Fy, s.Diadist)
	//nbars, cutting length, tot len
	var vstl, ar, nb, cl, tl float64
	switch s.Type{
		case 1:
		switch s.Endc{
			case 0,1:
			lx := s.Lx; ly := s.Ly
			if lx == 0{lx = s.Lspan}
			if ly == 0{ly = 2.0 * s.Ly; s.Ly = 2.0 * s.Lx}
			ar = lx * ly
			//lx += s.Bsup; ly += s.Bsup
			//fw,ck

			s.Vtot = s.Dused * (lx) * (ly) * 1e-9
			//so that one apple can be weighed against another orange 
			//WHY NOT THIS? covers total area, conservative, will do the same for a ribbed slab
			s.Vtot = s.Dused * (lx + s.Bsup) * (ly + s.Bsup) * 1e-9
			s.Vrcc = s.Vtot
			s.Afw = (lx+s.Bsup) * (ly+s.Bsup) * 1e-6

			//AGAIN.REMOVE THIS FOR NOW
			//s.Afw += s.Dused * 2.0 * (lx + 2.0 * s.Bsup + ly) * 1e-6

			//tensile main steel
			nb = math.Round(ly/s.Spcm) + 1.0; cl = (2.0 * ld1 + lx + s.Bsup - 2.0 * s.Efcvr)
			tl = math.Round(nb * cl)
			s.Bmap[s.Diamain] = append(s.Bmap[s.Diamain],[]float64{tl,nb,cl,1.0})
			s.Bsum[s.Diamain] += tl 
			vstl += tl * RbrArea(s.Diamain) * 1e-9
			//tensile dist steel
			nb = math.Round(lx/s.Spcm) + 1.0; cl = (2.0 * ld2 + ly +s.Bsup - 2.0 * s.Efcvr)
			tl = math.Round(nb * cl)
			s.Bmap[s.Diadist] = append(s.Bmap[s.Diadist],[]float64{tl,nb,cl,2.0})
			s.Bsum[s.Diadist] += tl			
			vstl += tl * RbrArea(s.Diadist) * 1e-9
			if s.Endc == 1{
				//THIS IS ONLY IF CAST INTEGRALLY WITH SUPPORTS ?? subramanian does it, why can't you
				//ADD TOP STEEL 50 % of ast at 0.1 lx over 2 supports 
				//kompressive main steel
				nb = math.Round(2.0 * 0.5 * ly/s.Spcm) + 1.0; cl = (2.0 * ld1 + 0.1*lx)
				tl = math.Round(nb * cl)
				s.Bmap[s.Diamain] = append(s.Bmap[s.Diamain],[]float64{tl,nb,cl,3.0})
				s.Bsum[s.Diamain] += tl
				vstl += tl * RbrArea(s.Diamain) * 1e-9
				//kompress dist steel
				nb = math.Round(0.2*lx/s.Spcm) + 1.0; cl = (2.0 * ld2 + ly - 2.0 * s.Efcvr)
				tl = math.Round(nb * cl)
				s.Bmap[s.Diadist] = append(s.Bmap[s.Diadist],[]float64{tl,nb,cl,4.0})
				s.Bsum[s.Diadist] += tl 
				vstl += tl * RbrArea(s.Diadist) * 1e-9
			}
			//vstl += (s.L1 + s.L2)
			s.Wstl = vstl * 7850.0
			//s.Wdl = s.Dused * 25.0 * 1e-3 * s.Lspan
			case 2:
			var vstl, vtot, ltot float64
			for i := 0; i < s.Nspans; i++{
				lx := s.Lspans[i]
				ly := s.Ly; if ly == 0{
					s.Ly = 2.0 * lx
					ly = 2.0 * lx
				}
				s1 := s.Spcspns[i][0]; s2 := s.Spcspns[i][1]; s3 := s.Spcspns[i][2]
				s4 := s.Sdspns[i][0]; s5 := s.Sdspns[i][1]; s6 := s.Sdspns[i][2]
				vtot += s.Dused * lx * ly * 1e-9
				ltot += lx
				var l1, l2, l3 float64
				switch i{
					case 0:
					l1 = 0.1 * lx; l2 = lx*(1.0-0.15-0.25); l3 = 0.3 * lx	
					case s.Nspans - 1:
					l1 = 0.3 * lx; l2 = lx*(1.0-0.25-0.15); l3 = 0.1 * lx

					default:
					l1 = 0.3 * lx; l2 = lx*(1.0-0.25-0.25); l3 = 0.3 * lx
				}
				//dist steel
				nb = math.Round(l1/s4 + 1.0) + math.Round(lx/s5 + 1.0) + math.Round(l3/s6 + 1.0)
				cl = ly + 2.0 * ld2; tl = math.Round(nb * cl)
				s.Bmap[s.Diadist] = append(s.Bmap[s.Diadist],[]float64{tl,nb,cl,float64(i+1)+0.1})
				s.Bsum[s.Diadist] += tl
				vstl += tl * RbrArea(s.Diadist)
				//main bottom steel; cl - l2
				nb = math.Round((ly/s2 + 1.0)/2.0)
				cl = l2; tl = math.Round(nb * cl)
				s.Bmap[s.Diamain] = append(s.Bmap[s.Diamain],[]float64{tl,nb,cl,float64(i+1)+0.2})
				s.Bsum[s.Diamain] += tl
				vstl += tl * RbrArea(s.Diamain)
				//main bottom steel; cl - lx + 2.0 ldev
				nb = math.Round((ly/s2 + 1.0)/2.0)
				cl = lx + 2.0 * ld1; tl = math.Round(nb * cl)
				s.Bmap[s.Diamain] = append(s.Bmap[s.Diamain],[]float64{tl,nb,cl,float64(i+1)+0.3})
				s.Bsum[s.Diamain] += tl
				vstl += tl * RbrArea(s.Diamain)
				//left top; cl1 = l1; cl2 = 0.5 l2 (forget this for now)
				nb = math.Round((ly/s1 + 1.0)/2.0)
				cl = l1+s.Bsup/2.0 + 2.0 * ld1; tl = math.Round(nb * cl)
				s.Bmap[s.Diamain] = append(s.Bmap[s.Diamain],[]float64{tl,nb,cl,float64(i+1)+0.4})
				s.Bsum[s.Diamain] += tl
				vstl += tl * RbrArea(s.Diamain)

				
				nb = math.Round((ly/s1 + 1.0)/2.0)
				cl = l1/2.0 + 2.0 * ld1; tl = math.Round(nb * cl)
				s.Bmap[s.Diamain] = append(s.Bmap[s.Diamain],[]float64{tl,nb,cl,float64(i+1)+0.4})
				s.Bsum[s.Diamain] += tl
				vstl += tl * RbrArea(s.Diamain)

				//right top
				nb = math.Round((ly/s3 + 1.0)*0.5)
				cl = l3+s.Bsup/2.0 + 2.0 * ld1; tl = math.Round(nb * cl)
				s.Bmap[s.Diamain] = append(s.Bmap[s.Diamain],[]float64{tl,nb,cl,float64(i+1)+0.4})
				s.Bsum[s.Diamain] += tl
				vstl += tl * RbrArea(s.Diamain)

				nb = math.Round((ly/s3 + 1.0)*0.5)
				cl = l3/2.0 + 2.0 * ld1; tl = math.Round(nb * cl)
				s.Bmap[s.Diamain] = append(s.Bmap[s.Diamain],[]float64{tl,nb,cl,float64(i+1)+0.4})
				s.Bsum[s.Diamain] += tl
				vstl += tl * RbrArea(s.Diamain)

			}
			s.Vtot = vtot
			s.Vrcc = s.Vtot
			ly := s.Ly; if ly == 0{ly = 1000.0}
			ar = ltot * ly
			s.Afw = (ltot) * (ly) * 1e-6

			//s.Afw += s.Dused * 2.0 * (ltot + ly + 2.0 * s.Bsup) * 1e-6

			s.Wstl = vstl * 7850.0 * 1e-9
		}
		case 2:
		lx := s.Lx; ly := s.Ly; ar = lx * ly
		s.Vtot = s.Dused * lx * ly * 1e-9
		s.Vrcc = s.Vtot
		s.Afw = lx * ly * 1e-6
		s1 := s.Spcms[0]; s2 := s.Spcms[1]; s3 := s.Spcms[2]; s4 := s.Spcms[3]
		//s.Afw += s.Dused * 2.0 * (lx + ly + 2.0 * s.Bsup) * 1e-6
		//short span support (lx)

		//YEOLDE
		//cl = 0.25 * lx + 2.0 * ld1
		
		sl := s.S1 + s.S2
		nb = 2.0 * math.Round(ly/s1 + 1.0); cl = sl * lx + 2.0 * ld1; tl = math.Round(nb * cl)
		s.Bmap[s.Diamain] = append(s.Bmap[s.Diamain],[]float64{tl,nb,cl,1.0})
		s.Bsum[s.Diamain] += tl
		vstl += tl * RbrArea(s.Diamain)
		
		//long span support (ly)
		sl = s.S3 + s.S4
		nb = 2.0 * math.Round(lx/s2 + 1.0); cl = sl * ly + 2.0 * ld2; tl = math.Round(nb * cl)
		s.Bmap[s.Diadist] = append(s.Bmap[s.Diadist],[]float64{tl,nb,cl,2.0})
		s.Bsum[s.Diadist] += tl
		vstl += tl * RbrArea(s.Diadist)
		
		//CURTAILMENT LENGTHS ARE ALL WRENG
		//WHAT ABOUT 50% of bars to continue until lx (support center)?
		
		//middle span - short (lx)
		nb = math.Round(0.5 * ly/s3 + 1.0); cl = 0.75 * lx + 2.0 * ld2; tl = math.Round(nb * cl)
		s.Bmap[s.Diamain] = append(s.Bmap[s.Diamain],[]float64{tl,nb,cl,3.0})
		s.Bsum[s.Diamain] += tl
		vstl += tl * RbrArea(s.Diamain)

		nb = math.Round(0.5 * ly/s3 + 1.0); cl = lx + s.Bsup; tl = math.Round(nb * cl)
		s.Bmap[s.Diamain] = append(s.Bmap[s.Diamain],[]float64{tl,nb,cl,3.0})
		s.Bsum[s.Diamain] += tl
		vstl += tl * RbrArea(s.Diamain)
		
		
		//middle span - long (ly)
		nb = math.Round(0.5 * lx/s4 + 1.0); cl = 0.75 * ly + 2.0 * ld2; tl = math.Round(nb * cl)
		s.Bmap[s.Diadist] = append(s.Bmap[s.Diadist],[]float64{tl,nb,cl,4.0})
		s.Bsum[s.Diadist] += tl
		vstl += tl * RbrArea(s.Diadist)

		
		nb = math.Round(0.5 * lx/s4 + 1.0); cl = ly + s.Bsup; tl = math.Round(nb * cl)
		s.Bmap[s.Diadist] = append(s.Bmap[s.Diadist],[]float64{tl,nb,cl,4.0})
		s.Bsum[s.Diadist] += tl
		vstl += tl * RbrArea(s.Diadist)

		//dist steel 300 c-c about 0.75 * lx *2 ; cl ly/8; 0.75 * ly * 2; cl lx/8 
		nb = math.Round(1.5 * lx/300.0 + 1.0); cl = ly/8.0; tl = math.Round(nb * cl)
		s.Bmap[s.Diadist] = append(s.Bmap[s.Diadist],[]float64{tl,nb,cl,5.0})
		s.Bsum[s.Diadist] += tl
		vstl += tl * RbrArea(s.Diadist)
		nb = math.Round(1.5 * ly/300.0 + 1.0); cl = lx/8.0; tl = math.Round(nb * cl)
		s.Bmap[s.Diadist] = append(s.Bmap[s.Diadist],[]float64{tl,nb,cl,6.0})
		s.Bsum[s.Diadist] += tl
		vstl += tl * RbrArea(s.Diadist)		
		vstl = vstl * 1e-9
		if s.Endc < 10{
			//ADD torsional reinforcement (for now, for all 4 sides)
			cl = lx/5.0; nb = 4.0 * math.Round(0.75 * lx/5.0/s3 + 1.0); tl = math.Round(cl * nb)
			s.Bmap[s.Diadist] = append(s.Bmap[s.Diadist],[]float64{tl,nb,cl,7.0})
		}
		s.Wstl = vstl * 7850.0
	}
	switch len(s.Kostin){
		case 3:
		s.Kost = s.Vrcc * s.Kostin[0] + s.Afw * s.Kostin[1] + s.Wstl * s.Kostin[2]
		default:
		s.Kost = s.Vrcc * CostRcc + s.Afw * CostForm + s.Wstl * CostStl
	}
	//s.Kost = s.Vrcc * CostRcc + s.Afw * CostForm + s.Wstl * CostStl
	s.Kunit = s.Kost/ar/1e-6
	return
}

//RQuant takes off quantities for a one way ribbed slab
func (s *RccSlb) RQuant(cb *CBm){
	//quantify a ribbed slab
	//quantity of top mesh/steel
	s.Bmap = make(map[float64][][]float64); s.Bsum = make(map[float64]float64)
	astmesh := 0.12 * s.Df * 1000.0/100.0
	//fmt.Println("ast/rmt->",astmesh)
	s.Astop = astmesh
	spcmin := 300.0
	//why 3.0 * (s.Df - 18.0) - subramanian
	if 3.0 * (s.Df - 18.0) < spcmin{
		spcmin = math.Round(math.Floor(3.0 * (s.Df - 18.0)/10.0)*10.0)
	}
	spc6 := math.Round(math.Floor(1000.0 * RbrArea(6.0)/astmesh/10.0)*10.0)
	if spcmin < spc6{
		spc6 = spcmin
	}
	s.Spctop = spc6; s.Diatop = 6.0
	//assume 6 mm dia steel for weight?
	var nribs, ly, ar, ltot float64
	ly = s.Ly; if ly == 0.0{
		ly = 1000.0
	}
	nribs = math.Round(ly/s.Bf) + 1.0
	fmt.Println("nribs-",nribs)
	//add rib quantities - each beam span is a t-beam of width bf
	s.Nribs = nribs
	s.BMs = make([][]float64, len(cb.RcBm))
	for _, barr:= range cb.RcBm{
		//s.BMs[i] = make([]float64, 3)
		s.Afw += barr[1].Afw * nribs
		s.Vrcc += barr[1].Vrcc * nribs
		s.Wstl += barr[1].Wstl * nribs
		s.Wstl += (barr[1].Lspan * (math.Ceil(ly/spc6) + 1.0) + ly * 1e-3 * (math.Ceil(barr[1].Lspan*1000.0/spc6) + 1.0)) * RbrArea(6.0) * 1e-6 * 7850.0
		//s.Vrcc += (ly - nribs*(s.Bf)) * barr[1].Dused * barr[1].Lspan * 1e-6

		//make slab solid at supports (0.1 lspan? earlier twas 0.15)
		s.Vrcc += (0.3 * barr[1].Lspan)* (s.Dw - s.Df) * (s.Bf - s.Bw) * 1e-6
		//add support flange pour?
		s.Vrcc += s.Bsup * s.Dw * ly * 1e-9

		ar += barr[1].Lspan * 1e3 * ly
		ltot += barr[1].Lspan
	}
	//add edge formwork for an extra dose of conservative paranoia
	//s.Afw += 2.0 * (ly * 1e-3 + ltot + 2.0 * s.Bsup * 1e-3) * s.Df * 1e-3
	//dw would be too damn high
	//s.Afw += 2.0 * (ly * 1e-3 + ltot + 2.0 * s.Bsup * 1e-3) * s.Dw * 1e-3
	s.Vtot = s.Vrcc
	s.Vfill = s.Dused * ar * 1e-9 - s.Vtot
	switch len(s.Kostin){
		case 3:
		s.Kost = s.Vrcc * s.Kostin[0] + s.Afw * s.Kostin[1] + s.Wstl * s.Kostin[2]
		default:
		s.Kost = s.Vrcc * CostRcc + s.Afw * CostForm + s.Wstl * CostStl
	}
	s.Kunit = s.Kost/ar/1e-6
	//s.Kost = s.Vrcc * CostRcc + s.Afw * CostForm + s.Wstl * CostStl
}

//R2Quant takes off quantities for a waffle slab
func (s *RccSlb) R2Quant(barr []*RccBm){
	//quantify a waffle (2-w ribbed) slab
	//quantity of top mesh/steel
	s.Bmap = make(map[float64][][]float64); s.Bsum = make(map[float64]float64)
	s.Vtot, s.Vrcc, s.Wstl, s.Afw = 0.0, 0.0, 0.0, 0.0
	//astmesh := 0.12 * s.Df * 1000.0/100.0
	//use 0.15 % steel for flange action with beam (incase?)
	astmesh := 0.15 * s.Df * 1000.0/100.0
	fmt.Println("ast/rmt->",astmesh)
	s.Astop = astmesh
	spcmin := 300.0
	//why 3.0 * (s.Df - 18.0) - subramanian (CHECK i GUESS 3xeffd?)
	if 3.0 * (s.Df - 18.0) < spcmin{
		spcmin = math.Round(math.Floor(3.0 * (s.Df - 18.0)/10.0)*10.0)
	}
	spc6 := math.Round(math.Floor(1000.0 * RbrArea(6.0)/astmesh/10.0)*10.0)
	if spcmin < spc6{
		spc6 = spcmin
	}
	s.Spctop = spc6; s.Diatop = 6.0
	var nrx, nry, ar, ltot float64
	nrx = math.Round(s.Ly/s.Bf) + 1.0; nry = math.Round(s.Lx/s.Bf) + 1.0
	fmt.Println("nrx, nry",nrx, nry)
	//add rib quantities - each beam span is a t-beam of width bf
	s.Nribs = nrx + nry
	s.BM = make([]float64, len(barr))
	for i, bm := range barr{
		s.BM[i] = bm.Mu
		_ = bm.Quant()
		fmt.Println("i, quant->",i)
		fmt.Println("bm.Vtot, bm.Vrcc, bm.Wstl, bm.Afw-",bm.Vtot, bm.Vrcc, bm.Wstl, bm.Afw)
		fmt.Println("lspan-",bm.Lspan)
		switch i{
			case 0:
			//short span (ribs of length lx along lx)
			s.Wstl += nrx * (bm.Wstl * (s.S1 + s.S2))
			case 1:
			s.Vtot += bm.Vrcc * nrx
			s.Wstl += bm.Wstl * nrx
			s.Afw += bm.Afw * nrx
			case 2:
			s.Wstl += nry * (bm.Wstl * (s.S3 + s.S4))
			case 3:
			//long span
			s.Vtot += (bm.Vrcc - bm.Df * bm.Bf * bm.Lspan * 1e-6)* nry
			s.Wstl += bm.Wstl * nry
			s.Afw += (bm.Afw - bm.Lspan * (bm.Bf - bm.Bw)*1e-3) * nry
			
		}
		s.Vrcc = s.Vtot
		//s.BMs[i] = make([]float64, 3)
		/*
		s.Afw += barr[1].Afw * nribs
		s.Vrcc += barr[1].Vrcc * nribs
		s.Wstl += barr[1].Wstl * nribs
		s.Wstl += (barr[1].Lspan * (math.Ceil(ly/spc6) + 1.0) + ly * 1e-3 * (math.Ceil(barr[1].Lspan*1000.0/spc6) + 1.0)) * RbrArea(6.0) * 1e-6 * 7850.0
		s.Vrcc += (ly - nribs*(s.Bf)) * barr[1].Dused * barr[1].Lspan * 1e-6
		s.Vrcc += 0.3 * barr[1].Lspan * (s.Dw - s.Df) * (s.Bf - s.Bw) * 1e-6
		ar += barr[1].Lspan * 1e3 * ly
		ltot += barr[1].Lspan
		*/
	}
	ar = s.Lx * s.Ly
	//add edge formwork for an extra dose of conservative paranoia
	//s.Afw += 2.0 * (ly * 1e-3 + ltot + 2.0 * s.Bsup * 1e-3) * s.Df * 1e-3
	//dw would be too damn high
	//s.Afw += 2.0 * (ly * 1e-3 + ltot + 2.0 * s.Bsup * 1e-3) * s.Dw * 1e-3
	s.Vfill = s.Dused * ar * 1e-9 - s.Vtot
	switch len(s.Kostin){
		case 3:
		s.Kost = s.Vrcc * s.Kostin[0] + s.Afw * s.Kostin[1] + s.Wstl * s.Kostin[2]
		default:
		s.Kost = s.Vrcc * CostRcc + s.Afw * CostForm + s.Wstl * CostStl
	}
	fmt.Println(ltot)
	s.Kunit = s.Kost/ar/1e-6
	//s.Kost = s.Vrcc * CostRcc + s.Afw * CostForm + s.Wstl * CostStl
	fmt.Println("s.Vrcc, s.Afw, s.Wstl",s.Vrcc, s.Afw, s.Wstl)
}


//R2Table generates an ascii table report for a waffle slab
func (s *RccSlb) R2Table(barr []*RccBm, printz bool){
	//ribbed and waffle slab table
	rezstr := new(strings.Builder)
	var t string
	switch s.Endc{
		case 10:
		t = "2-way ss waffle slab"
		default:
		t = fmt.Sprintf("2-way waffle slab endc %v",s.Endc)
	}
	hdr := fmt.Sprintf("%s\nrcc slab report\ndate-%s\n%s\n",ColorYellow,time.Now().Format("2006-01-02"),ColorReset)
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	hdr = ""
	hdr += fmt.Sprintf("%s\n%s slab \nlspan %.1f, lx %.1f, ly %.1f mm b support %.1f mm\n",s.Title,t, s.Lspan, s.Lx, s.Ly, s.Bsup)
	hdr += fmt.Sprintf("grade of concrete M %.1f\nsteel - main Fe %.f, dist Fe %.f\n", s.Fck, s.Fy, s.Fyd)
	hdr += fmt.Sprintf("cover - nominal %0.1f mm effective %0.1f mm\n", s.Nomcvr, s.Efcvr)
	hdr += fmt.Sprintf("loads - dl %.3f kN/m2, ll %0.3f kN/m2\n", s.DL, s.LL)
	hdr += fmt.Sprintf("moment redistribution %.3f\n", s.DM)
	table := tablewriter.NewWriter(rezstr)
	if s.Dz{
		hdr += fmt.Sprintf("%s\nrib depth %.f mm width %.f mm spacing %.f mm\nslab (topping) depth %.f mm\n",ColorYellow,s.Dused, s.Bw, s.Bf, s.Df)
		hdr += fmt.Sprintf("number of ribs %.f\n",s.Nribs)
		
	}
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorReset)
	rezstr.WriteString(ColorPurple)
	var row string
	table.SetCaption(true,"rib reinforcement")
	table.SetHeader([]string{"span","len","loc","bm\n(kn-m)","top dia\n(mm)","no.","asc req\n(mm2)","asc prov","bottom dia\n(mm)","no.","ast req\n(mm2/m)","ast prov\n(mm2/m)"})

	for j, val := range []string{"short span support(t)","short span mid(b)","long span support(t)","long span mid(b)"}{
		var d1, n1, d2, n2, ascr, asc, astr, ast, lspan, mu float64
		b := barr[j]; mu = barr[j].Mu
		lspan = barr[j].Lspan
		if b.Asc > 0.0{
			r := b.Rbrc
			d1, n1 = r[2], r[0]
			asc, ascr = r[4], r[5]
			
		}
		if b.Ast > 0.0{
			r := b.Rbrt
			d2, n2 = r[2], r[0]
			ast, astr = r[4], r[5]				
		}
		row = fmt.Sprintf("%v, %.f, %s, %.2f, %.f, %.f, %.2f, %.2f,%.f, %.f, %.2f, %.2f",j+1,lspan,val,mu, d1, n1, ascr, asc, d2, n2, astr, ast)
		table.Append(strings.Split(row,","))
	}
	table.Render()
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"mesh reinforcement")
	table.SetHeader([]string{"dia\n(mm)","spacing\n(mm)","ast req\n(mm2)","ast prov\n(mm2)"})
	row = fmt.Sprintf("%.f, %.f, %.f, %.f",s.Diatop, s.Spctop, s.Astop, math.Round(RbrArea(s.Diatop)*1000.0/s.Spctop))
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorBlue)
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"vol tot(m3)","vol rcc(m3)","wt stl(kg)","form area (m2)","cost (rs)","unit cost(rs/m2)"})
	table.SetCaption(true,"quantity take off")
	row = fmt.Sprintf("%.3f, %.3f, %.3f, %.3f, %.f, %.2f\n",s.Vtot,s.Vrcc,s.Wstl,s.Afw, s.Kost, s.Kunit)
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorReset)
	s.Report = fmt.Sprintf("%s",rezstr)
	if printz{
		fmt.Println(s.Report)
	}

}

//RTable generates an ascii table report for a ribbed slab
func (s *RccSlb) RTable(cb *CBm, printz bool){
	//ribbed and waffle slab table
	rezstr := new(strings.Builder)
	var t string
	switch s.Endc{
		case 0:
		t = "ribbed cantilever"
		case 1:
		t = "1 way ribbed ss"
		case 2:
		t = "1 way ribbed cs"
	}
	hdr := fmt.Sprintf("%s\nrcc slab report\ndate-%s\n%s\n",ColorYellow,time.Now().Format("2006-01-02"),ColorReset)
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	hdr = ""
	hdr += fmt.Sprintf("%s\n%s slab \nlspan %.1f, lx %.1f, ly %.1f mm b support %.1f mm\n",s.Title,t, s.Lspan, s.Lx, s.Ly, s.Bsup)
	hdr += fmt.Sprintf("grade of concrete M %.1f\nsteel - main Fe %.f, dist Fe %.f\n", s.Fck, s.Fy, s.Fyd)
	hdr += fmt.Sprintf("cover - nominal %0.1f mm effective %0.1f mm\n", s.Nomcvr, s.Efcvr)
	hdr += fmt.Sprintf("loads - dl %.3f kN/m2, ll %0.3f kN/m2\n", s.DL, s.LL)
	hdr += fmt.Sprintf("moment redistribution %.3f\n", s.DM)
	table := tablewriter.NewWriter(rezstr)
	if s.Dz{
		hdr += fmt.Sprintf("%s\nrib depth %.f mm width %.f mm spacing %.f mm\nslab (topping) depth %.f mm\n",ColorYellow,s.Dused, s.Bw, s.Bf, s.Df)
		switch s.Type{
			case 3:
			hdr += fmt.Sprintf("number of ribs (span dir) %.f\n",s.Nribs)
			case 4:
			//waffle slab
		}
	}
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorReset)
	rezstr.WriteString(ColorPurple)
	ls := make([]float64, 3)
	ls[1] = 1.0
	if s.Nspans == 1{
		ls[0] = 0.1; ls[2] = 0.1
	} else {
		ls[0] = 0.15; ls[2] = 0.15
	}
	var row string
	table.SetCaption(true,"rib reinforcement")
	table.SetHeader([]string{"span","len","loc","bm\n(kn-m)","top dia\n(mm)","no.","asc req\n(mm2)","asc prov","bottom dia\n(mm)","no.","ast req\n(mm2/m)","ast prov\n(mm2/m)"})
	for i := 0; i < s.Nspans; i++{
		bmarr := cb.RcBm[i]
		for j, val := range []string{"left(t)","mid(b)","right(t)"}{
			var d1, n1, d2, n2, ascr, asc, astr, ast, lspan, mu float64
			b := bmarr[j]; mu = bmarr[j].Mu
			lspan = s.Lspans[i] * ls[j]
			if b.Asc > 0.0{
				r := b.Rbrc
				d1, n1 = r[2], r[0]
				asc, ascr = r[4], r[5]
				if j != 1{
					ld1 := BarDevLen(s.Fck, s.Fy, d1)
					lspan += math.Ceil(ld1/10.0) * 10.0
				}
			}
			if b.Ast > 0.0{
				r := b.Rbrt
				d2, n2 = r[2], r[0]
				ast, astr = r[4], r[5]				
			}
			row = fmt.Sprintf("%v, %.f, %s, %.2f, %.f, %.f, %.2f, %.2f,%.f, %.f, %.2f, %.2f",i+1,lspan,val,mu, d1, n1, ascr, asc, d2, n2, astr, ast)
			table.Append(strings.Split(row,","))
		}
	}
	table.Render()
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"mesh reinforcement")
	table.SetHeader([]string{"dia\n(mm)","spacing\n(mm)","ast req\n(mm2)","ast prov\n(mm2)"})
	row = fmt.Sprintf("%.f, %.f, %.f, %.f",s.Diatop, s.Spctop, s.Astop, math.Round(RbrArea(s.Diatop)*1000.0/s.Spctop))
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorBlue)
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"vol tot(m3)","vol rcc(m3)","wt stl(kg)","form area (m2)","cost (rs)","unit cost(rs/m2)"})
	table.SetCaption(true,"quantity take off")
	row = fmt.Sprintf("%.3f, %.3f, %.3f, %.3f, %.f, %.2f\n",s.Vtot,s.Vrcc,s.Wstl,s.Afw, s.Kost, s.Kunit)
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorReset)
	s.Report = fmt.Sprintf("%s",rezstr)
	if printz{
		fmt.Println(s.Report)
	}

}

//Table generates an ascii table report for an rcc slab
func (s *RccSlb) Table(printz bool){
	rezstr := new(strings.Builder)
	var t string
	if s.Type == 2{
		switch s.Endc{
			case 10:
			t = "2 way ss"
			default:
			t = fmt.Sprintf("2 way endc %v",s.Endc)
		}
	} else {
		switch s.Endc{
			case 0:
			t = "cantilever"
			case 1:
			t = "1 way ss"
			case 2:
			t = "1 way cs"
		}
	}
	if s.Type == 3{t += "-ribbed"}
	hdr := fmt.Sprintf("%s\nrcc slab report\ndate-%s\n%s\n",ColorYellow,time.Now().Format("2006-01-02"),ColorReset)
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	hdr = ""
	hdr += fmt.Sprintf("%s\n%s slab \nlspan %.1f, lx %.1f, ly %.1f mm b support %.1f mm\n",s.Title,t, s.Lspan, s.Lx, s.Ly, s.Bsup)
	hdr += fmt.Sprintf("grade of concrete M %.1f\nsteel - main Fe %.f, dist Fe %.f\n", s.Fck, s.Fy, s.Fyd)
	hdr += fmt.Sprintf("cover - nominal %0.1f mm effective %0.1f mm\n", s.Nomcvr, s.Efcvr)
	hdr += fmt.Sprintf("loads - dl %.3f kN/m2, ll %0.3f kN/m2\n", s.DL, s.LL)
	if s.Type == 3{
		hdr += fmt.Sprintf("%s\nrib depth %.f mm width %.f mm spacing %.f mm\nslab (topping) depth %.f mm\n",ColorYellow,s.Dused, s.Bw, s.Bf, s.Df)
	} else{
		hdr += fmt.Sprintf("%s\nslab depth %.f mm\n\n",ColorYellow,s.Dused)
	}
	
	if (s.Type == 1 || s.Type == 3) && s.Endc == 2{
		hdr += fmt.Sprintf("moment redistribution %.3f\n", s.DM)
	}
	rezstr.WriteString(hdr)
	
	rezstr.WriteString(ColorBlue)
	var row string
	table := tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"slab specs")
	table.SetHeader([]string{"type","concrete","steel(main)","steel(dist.)","nom.cvr\n(mm)","eff.cvr\n(mm)","DM"})
	row = fmt.Sprintf("%s,M%.f,Fe%.f,Fe%.f,%.2f,%.2f,%.3f",t,s.Fck,s.Fy,s.Fyd,s.Nomcvr,s.Efcvr,s.DM)
	table.Append(strings.Split(row,","))
	table.Render()

	rezstr.WriteString(ColorCyan)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"slab geometry")
	table.SetHeader([]string{"span\n(mm)","lx\n(mm)","ly\n(mm)","b.sup.\n(mm)","depth\n(mm)"})
	row = fmt.Sprintf("%.2f,%.2f,%.2f,%.2f,%.2f",s.Lspan, s.Lx, s.Ly, s.Bsup, s.Dused)
	table.Append(strings.Split(row,","))
	table.Render()
	

	rezstr.WriteString(ColorRed)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"design loads")
	table.SetHeader([]string{"dl\n(kn/m2)","ll\n(kn/m2)"})
	row = fmt.Sprintf("%.2f,%.2f",s.DL, s.LL)
	table.Append(strings.Split(row,","))
	table.Render()
	
	table = tablewriter.NewWriter(rezstr)
	//area := s.Lx * s.Ly
	//if s.Ly == 0{area = s.Lspan * 1000.0}
	if s.Dz{
		rezstr.WriteString(ColorPurple)
		table.SetCaption(true,"reinforcement")
		switch s.Type{
			case 1:
			/*s.Astm = asprov; s.Spcm = spcmain; s.Astd = asdprov; s.Spcd = sds; s.Astreq = astreq
			s.Dused = dused; s.BM = append(s.BM, mdu)*/
			switch s.Endc{
				case 0,1:
				table.SetHeader([]string{"loc","bm\n(kn-m)","dia\n(mm)"," spacing\n(mm)","ast req\n(mm2/m)","ast prov\n(mm2/m)"})
				row = fmt.Sprintf("%s, %.2f, %.0f, %.0f, %.0f, %.0f mm2\n","main stl",s.BM[0],s.Diamain,s.Spcm,s.Astreq,s.Astm)
				table.Append(strings.Split(row,","))
				row = fmt.Sprintf("%s, %.2f, %.0f, %.0f, %.0f, %.0f mm2\n","dist stl",s.BM[0],s.Diadist,s.Spcd,s.Astd,s.Astd)
				table.Append(strings.Split(row,","))
				case 2:
				table.SetHeader([]string{"span","len","loc","bm\n(kn-m)","dia main\n(mm)","spacing\n(mm)","dia dist\n(mm)","spacing\n(mm)","ast req\n(mm2/m)","ast prov\n(mm2/m)"})
				for i := 0; i < s.Nspans; i++{
					for j, val := range []string{"left(t)","mid(b)","right(t)"}{
						//rezstring += fmt.Sprintf("dia - %.f mm spacing %.f mm dist %.f mm spacing %.f mm\n",s.Diaspns[i][j], s.Spcspns[i][j],s.Diadist, s.Sdspns[i][j])
						//rezstring += fmt.Sprintf("ast req - %.f mm2 ast prov %.f mm2\n",s.Astrs[i][j], s.Astps[i][j])
						row = fmt.Sprintf("%v, %.f, %s, %.2f, %.0f, %.0f, %.0f, %.0f, %.0f, %.0f",i+1,s.Lspan,val,s.BMs[i][j],s.Diaspns[i][j],s.Spcspns[i][j],s.Diadist,s.Sdspns[i][j],s.Astrs[i][j],s.Astps[i][j])
						table.Append(strings.Split(row,","))
					}
				}
			}
			case 2:
			table.SetHeader([]string{"loc","bm\n(kn-m)","dia\n(mm)"," spacing\n(mm)","ast req\n(mm2/m)","ast prov\n(mm2/m)"})
			for i, loc := range []string{"short span support","short span midspan","long span support","long span midspan"}{
				row = fmt.Sprintf("%s, %.2f, %.0f, %.0f, %.0f, %.0f mm2\n",loc,s.BM[i],s.Dias[i],s.Spcms[i],s.Astr[i],s.Asts[i])
				table.Append(strings.Split(row,","))
			}
		}
		table.Render()
		rezstr.WriteString(ColorBlue)
		table = tablewriter.NewWriter(rezstr)
		table.SetHeader([]string{"vol tot(m3)","vol rcc(m3)","wt stl(kg)","form area (m2)","cost (rs)","unit cost(rs/m2)"})
		table.SetCaption(true,"quantity take off")
		row = fmt.Sprintf("%.3f, %.3f, %.3f, %.3f, %.f, %.2f\n",s.Vtot,s.Vrcc,s.Wstl,s.Afw, s.Kost, s.Kunit)
		table.Append(strings.Split(row,","))
		table.Render()
		rezstr.WriteString(ColorReset)
	}
	s.Report = fmt.Sprintf("%s",rezstr)
	if printz{
		fmt.Println(s.Report)
	}
}

//EndC returns the end condition of a two way rectangular slab
//based on s.Ec (edge continuity of slab supports)
//to be used with framegen funcs
func (s *RccSlb) EndC(){
	//get slab endc of 2 way rect slab
	//here (hear) lx and ly are global x and y lengths (use with framegen)
	cl := s.Ec[0]; cr := s.Ec[1]; ct := s.Ec[2]; cb := s.Ec[3]
	lx := s.Lx; ly := s.Ly
	switch s.Slbc{
		case 0:
		//simply supported
		s.Endc = 10
		switch ly <= lx{
			case true:
			s.Ly = lx
			s.Lx = ly
			s.Sdir = 1
			case false:
			s.Ly = ly
			s.Lx = lx
			s.Sdir = 2
		}
		case 1:
		//continuous along y
		//cl = 0; cr = 0
		switch ly <= lx{
			case true:
			s.Sdir = 1
			s.Ly = lx
			s.Lx = ly
			switch {
			case ct == 0 && cb == 0:
				s.Endc = 9
			case ct == 1 && cb == 0:
				s.Endc = 7
			case cb == 1 && ct == 0:
				s.Endc = 7
			case ct == 1 && cb == 1:
				s.Endc = 5
			}
			case false:
			s.Sdir = 2
			s.Ly = ly
			s.Lx = lx
			switch {
			case ct == 0 && cb == 0:
				s.Endc = 9
			case ct == 1 && cb == 0:
				s.Endc = 8
			case ct == 0 && cb == 1:
				s.Endc = 8
			case ct == 1 && cb == 1:
				s.Endc = 6
			}
		}
		case 2:
		//continuous along x
		switch ly <= lx{
			case true:
			s.Sdir = 1
			s.Ly = lx
			s.Lx = ly
			switch {
			case cl == 0 && cr == 0:
				s.Endc = 9
			case cl == 1 && cr == 0:
				s.Endc = 8
			case cl == 0 && cr == 1:
				s.Endc = 8
			case cl == 1 && cr == 1:
				s.Endc = 6
			}
			case false:
			s.Sdir = 2
			s.Ly = ly
			s.Lx = lx
			switch {
			case cl == 0 && cr == 0:
				s.Endc = 9
			case cl == 1 && cr == 0:
				s.Endc = 7
			case cl == 0 && cr == 1:
				s.Endc = 7
			case cl == 1 && cr == 1:
				s.Endc = 5
			}

		}
		case 3:
		//as and when lol
	}

}

//EndC2W updates ns nl, sets i1, i2, i3, i4 based on end conditions
//for a two way slab with corners restrained from torsion
//s1, s2, s3, s4 - top, b, l, r - (lx,lx, ly, ly top steel lengths)
func (s *RccSlb) EndC2W(){
	s.FitI()
	switch s.Endc{
		case 1:
		//interior panels - all sides continuous
		s.Ns = 0; s.Nl = 0
		s.S1, s.S2, s.S3, s.S4 = 0.3,0.3,0.3,0.3
		
		case 2:
		//one short edge discontinuous
		s.Ns = 1; s.Nl = 0
		s.I1 = 0.0
		s.S1, s.S2, s.S3, s.S4 = 0.1,0.3,0.3,0.3
		case 3:
		//one long edge discontinuous
		s.Ns = 0; s.Nl = 1
		s.I2 = 0.0
		s.S1, s.S2, s.S3, s.S4 = 0.3,0.3,0.3,0.1
		case 4:
		//two adjacent edges discontinuous
		s.Ns = 1; s.Nl = 1
		s.I1 = 0; s.I4 = 0
		s.S1, s.S2, s.S3, s.S4 = 0.1,0.3,0.1,0.3
		case 5:
		//two short edges discontinuous
		s.Ns = 2; s.Nl = 0
		s.I1 = 0; s.I3 = 0
		s.S1, s.S2, s.S3, s.S4 = 0.1,0.1,0.3,0.3
		case 6:
		//two long edges discontinuous
		s.Ns = 0; s.Nl = 2
		s.I2 = 0; s.I4 = 0
		s.S1, s.S2, s.S3, s.S4 = 0.3,0.3,0.1,0.1
		case 7:
		//3 edges discontinuous; one long edge continuous
		s.Ns = 2; s.Nl = 1
		s.I1 = 0; s.I2 = 0; s.I3 = 0
		s.S1, s.S2, s.S3, s.S4 = 0.1,0.1,0.1,0.3
		case 8:
		//3 edges discontinuous; one short edge continuous
		s.Ns = 1; s.Nl = 2
		s.I1 = 0; s.I2 = 0; s.I4 = 0
		s.S1, s.S2, s.S3, s.S4 = 0.3,0.1,0.1,0.1
		case 9:
		//4 edges discontinuous
		s.Ns = 2; s.Nl = 2
		s.I1 = 0; s.I2 = 0; s.I3 = 0; s.I4 = 0
		s.S1, s.S2, s.S3, s.S4 = 0.1,0.1,0.1,0.1
		case 10:
		//simply supported
		s.I1 = 0; s.I2 = 0; s.I3 = 0; s.I4 = 0
		s.S1, s.S2, s.S3, s.S4 = 0.1,0.1,0.1,0.1
	}
}

//Tweaks is for one way slab design tweaks
//TODO
func (s *RccSlb) Tweaks(){
}



/*
					   - ++
1 - way slab				     ++
----------------------  |----------------------  |

     +--------------------------------------------------------------------------+
     | --------------------........            	        .......----------------	|
     |                                                                         	|
     |  ---------------------------------------------------------------------- 	|
     +---------+--------------------------------------------------------+-------+
     |         |			      				|       |
     |         |			      				|       |
     |         |			      				|       |
     |         |			      				|       |
     +---------+			      				+-------+

     ............................	      		.........................


			case 3:
			//HUH? HUH?
			
			table.SetHeader([]string{"span","len","loc","bm\n(kn-m)","top dia\n(mm)","no.\n(mm)","dia top\n(mm)","no.\n(mm)","ast req\n(mm2/m)","ast prov\n(mm2/m)","asc req\n(mm2/m)","asc prov\n(mm2/m)"})
			for i := 0; i < s.Nspans; i++{
				for j, val := range []string{"left(t)","mid(b)","right(t)"}{
					row = fmt.Sprintf("%v, %.f, %s, %.2f",i+1,s.Lspans[i],val,s.BMs[i][j])
					table.Append(strings.Split(row,","))
				}
			}

*/
