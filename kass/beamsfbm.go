package barf

import (
	"fmt"
	"log"
	"math"
	"os"
)

//BeamRez stores bm and sf values at 20 intervals over a member
//is the backbone of many funcs/fields
type BeamRez struct{
	Mem                                                               int
	SF, BM, Dxs, Dxs_BM, Dxs_SF, Maxs, Locs, Cfxs, Venv, Mnenv, Mpenv []float64
	Txtplot                                                           string
	Xs                                                                []float64
	Qf                                                                []float64
}

//BmFrcX interpolates shear and bending moment at x from vs ms arrays
func BmFrcX(lx, xdiv, l float64, xs, vs, ms []float64) (vl, ml float64){
	i := int(math.Ceil(lx/xdiv))
	x := xs[i]; vx := vs[i]
	switch{
		case x == lx:
		vl = vx
		ml = ms[i]
		case x > lx:
		switch i{
			case 0:
			vl = vs[i+1]
			ml = ms[i+1]
			default:
			vl = vx + (-vx + vs[i-1])*(x-lx)/xdiv
			ml = ms[i] + 0.5 * (x - lx)*(vl + vx)
			//ml = ms[i] + (ms[i] - ms[i-1])*(lsx - x)/xdiv
			
		}
		case x < lx:
		switch i{
			case 0:
			vl = vs[i+1]
			ml = ms[i+1]
			default:
			vl = vx + (- vx + vs[i-1])*(lx-x)/xdiv
			ml = ms[i] + 0.5 * (lx-x)*(vl + vx)
		}
	}
	return
}

//BeamFrc is the entry func for 2d beam sf and bm calcs from model calcs
//calls bmsfcalc
func BeamFrc(ncjt, member int, mem *Mem, ldcases [][]float64, spanchn chan BeamRez, plotbm bool) {
	var rl, ml, re, me, l, e, iz, ar float64
	var r BeamRez
	switch {
	case ncjt == 2://beam
		//memrel = 3
		rl = mem.Qf[0]
		ml = mem.Qf[1]
		re = mem.Qf[2]
		me = mem.Qf[3]
		l = mem.Geoms[0]
		e = mem.Geoms[1]
		iz = mem.Geoms[2]
		ar = mem.Geoms[3]
		r = Bmsfcalc(member, ldcases, l, e, ar, iz, rl, ml, re, me, plotbm, mem.Clvr)
	case ncjt == 3://frame
		rl = mem.Qf[1]
		ml = mem.Qf[2]
		re = mem.Qf[4]
		me = mem.Qf[5]
		l = mem.Geoms[0]
		e = mem.Geoms[1]
		ar = mem.Geoms[2]
		iz = mem.Geoms[3]
		r = Bmsfcalc(member, ldcases, l, e, ar, iz, rl, ml, re, me, plotbm, mem.Clvr)
	}
	spanchn <- r
}

//Bmsfcalc is basically hulse section 2.1 - calculates bending moments and shear force at span/20 divs
//backbone of many funcs
//takes in member (index), ldcases - list of member loads, l- length, e- young's modulus, a - area, iz - moment of inertia
//rl - left reaction, ml - left moment, re - end reaction (right), me - end moment (right)

func Bmsfcalc(member int, ldcases [][]float64, l, e, a, iz, rl, ml, re, me float64, plotbm, clvr bool) (BeamRez){
	//mfrtyp := 0
	ndiv := 21
	nd := 20.0
	div := l/(nd)
	//div := l / 20.0
	xs := make([]float64, ndiv)
	for i := 0; i < ndiv; i++ {
		xs[i] = div * float64(i)
	}
	vxs := make([]float64, ndiv)
	mxs := make([]float64, ndiv)

	for i, x := range xs{
		vxs[i] = rl
		mxs[i] = rl*x - ml
	}
	//var fsb, fmb, fse float64
	for _, ldcase := range ldcases {
		ltyp := int(ldcase[1])
		la := ldcase[4]
		lb := ldcase[5]
		var w, wa, wb, ldcx, ldcvr, W float64
		switch ltyp {
		case 1:
			w = ldcase[2]
			ldcx = la
			lb = l - la
			W = w
		case 3:
			w = ldcase[2]
			ldcvr = (l - la - lb)
			ldcx = la + ldcvr/2.0
			W = w*ldcvr
		case 4:
			wa = ldcase[2]
			wb = ldcase[3]
			ldcvr = l - la - lb
			if wa >= wb {
				W = wb*ldcvr + (wa-wb)*ldcvr/2.0
				ldcx = la + (ldcvr)*((wa+2.0*wb)/(3.0*(wa+wb)))
			} else {
				W = wa*ldcvr + (wb -wa)*ldcvr/2.0
				lcx1 := (ldcvr)*((wb+2.0*wa)/(3.0*(wa+wb)))
				ldcx = la + ldcvr - lcx1
			}
		}
		for i, x := range xs {
			x1 := x - la
			if ltyp == 1{
				switch {
				case x < la://brez == -1: //before load start
					vxs[i] += 0.0
					mxs[i] += 0.0
				case x >= la:	
					vxs[i] -= w
					mxs[i] -= w*(x1)
				}
			} else {
				switch {
				case x <= la:
					vxs[i] += 0.0
					mxs[i] += 0.0
				case x >= la + ldcvr: //after load cover
					switch ltyp{
						case 4:
						if wa >= wb {
							vxs[i] -= W
							mxs[i] -= W*(x-ldcx)
						} else {
							vxs[i] -= W
							mxs[i] -= W*ldcx
						}
						case 3:
						vxs[i] -= W
						mxs[i] -= W*(x-ldcx)
					}
					
				default: //within load cover
					switch ltyp {
					case 3:
						w1 := w * x1
						ldcx2 := la + x1/2.0
						vxs[i] -= w1
						mxs[i] -= w1*(x-ldcx2)
						//above should work even if ldcx < x? who (or what) knows
					case 4:
						wm := math.Abs(wb - wa)/(l - lb - la)
						dw := wm * x1
						if wa >= wb{	
							wx := wa - dw
							w1 := (wx * x1) + (dw * x1 / 2.0)
							ldcx2 := (ldcvr)*(wa+2.0*wx)/(3.0*(wa+wx))
							//this is ze distance from side wa, so dist from x = x - la - 
							ldcx2 = x - la - ldcx2
							vxs[i] -= w1 
							//mxs[i] -= w1*(la + ldcx) //CHECK THIS EH
							mxs[i] -= w1*ldcx2

						} else {
							wx := wa + dw
							w1 := (wa * x1) + (dw * x1 / 2.0)
							ldcx2 := (ldcvr)*(wx+2.0*wa)/(3.0*(wa+wx))
							vxs[i] -= w1
							mxs[i] -= w1*ldcx2
						}
					}
				}

			}
		}		
	}
	//calc unit moment vector
	mdxs := make([]float64, ndiv)
	vdxs := make([]float64, ndiv)
	dxs := make([]float64, ndiv)
	var homerx, mx, vx float64
	for i := 0; i < ndiv; i++ {
		zd := float64(i) * div
		var mdxk, vdxk float64
		for x := 0; x < ndiv; x++ {
			z := float64(x) * div
			switch {
			case x == 0:
				homerx = 1
			case x == ndiv -1:
				homerx = 1
			case x%2 == 0:
				//homerx = 4
				homerx = 2
			case x%2 != 0:
				//homerx = 2
				homerx = 4
			}
			switch clvr{
				case false:	
				if z <= zd {
					mx = z * (l - zd) / l
					vx = (l - zd) / l
				} else {
					mx = zd * (l - z) / l
					vx = -zd / l
				}
				case true:
				//fmt.Println("clvr!",member)
				if z >= zd{
					mx = 0.0
					vx = 0.0
				} else {
					mx = z - zd
					vx = 1.0
				}
			}
			mdxk += mxs[x] * mx * homerx
			vdxk += vxs[x] * vx * homerx
		}
		mdxs[i] = mdxk * (l/nd) / (3.0 * e * iz)
		
		g := e / (2 * 1.15)
		if a > 0.0 {vdxs[i] = 1.5 * vdxk * (l / nd) / (3.0 * g * a)}
		dxs[i] = mdxs[i] 
	}
	var mmax, vmax, vmaxx, mmaxx, dmax, dmaxx, mpmax, mpmaxx float64
	var cfxs []float64
	for i, x := range xs {
		if i != 0 && mxs[i]*mxs[i-1] < 0.0 {
			cfxs = append(cfxs, x)
		}
		if math.Abs(vmax) < math.Abs(vxs[i]) {
			vmax = vxs[i]
			vmaxx = x
		}
		if math.Abs(mmax) < math.Abs(mxs[i]) {
			mmax = mxs[i]
			mmaxx = x
		}
		if math.Abs(dmax) < math.Abs(dxs[i]) {
			dmax = dxs[i]
			dmaxx = x
			//dxs[i] = mdxs[i] + vdxs[i]
		}
		if mxs[i] > 0.0 && mpmax < mxs[i]{
			mpmax = mxs[i]
			mpmaxx = x
		}
	}
	var s string
	if plotbm {s = PlotBmSfBm(xs, vxs, mxs, dxs, l, true)}
	//fmt.Println(s)
	var rez BeamRez
	rez.Mem = member
	rez.SF = vxs
	rez.BM = mxs
	rez.Dxs = dxs
	rez.Dxs_BM = mdxs
	rez.Dxs_SF = vdxs
	rez.Maxs = []float64{vmax, mmax, dmax, mpmax}
	rez.Locs = []float64{vmaxx, mmaxx, dmaxx, mpmaxx}
	rez.Cfxs = cfxs
	rez.Txtplot = s
	rez.Xs = xs
	rez.Qf = []float64{rl, ml, re, me}
	return rez
}

//BmCalcDM is basically BmSfcalc without end shears
//needed to copy paste earlier bmsf calc to work with left moment and right moment input
//for moment redistribution funcs
//(the wise man writes multiple similar funtions)

func BmCalcDM(member int, ldcases [][]float64, l, e, ar, iz, ml, mr float64, plotbm bool) (BeamRez){
	ndiv := 21
	nd := float64(ndiv) - 1.0
	div := l/(nd)
	xs := make([]float64, 21)
	for i := 0; i < 21; i++ {
		xs[i] = div * float64(i)
	}
	vxs := make([]float64, 21)
	mxs := make([]float64, 21)
	rl := (ml + mr)/l
	for i, x := range xs {
		vxs[i] = rl
		mxs[i] = rl*x - ml
	}
	//var fsb, fmb, fse float64
	for _, ldcase := range ldcases{
		ltyp := int(ldcase[1])
		la := ldcase[4]
		lb := ldcase[5]
		var w, wa, wb, ldcx, ldcvr, W float64
		switch ltyp {
		case 1:
			w = ldcase[2]
			ldcx = la
			lb = l - la
			W = w
		case 3:
			w = ldcase[2]
			ldcvr = (l - la - lb)
			ldcx = la + ldcvr/2.0
			W = w*ldcvr
		case 4:
			wa = ldcase[2]
			wb = ldcase[3]
			ldcvr = l - la - lb
			if wa >= wb {
				W = wb*ldcvr + (wa-wb)*ldcvr/2.0
				ldcx = la + (ldcvr)*((wa+2.0*wb)/(3.0*(wa+wb)))
			} else {
				W = wa*ldcvr + (wb -wa)*ldcvr/2.0
				lcx1 := (ldcvr)*((wb+2.0*wa)/(3.0*(wa+wb)))
				ldcx = la + ldcvr - lcx1
			}
		}
		rl = W * (l-ldcx)/l
		for i, x := range xs {
			//brez := big.NewFloat(x).Cmp(big.NewFloat(la))
			//erez := big.NewFloat(x).Cmp(big.NewFloat(la + ldcvr))
			x1 := x - la
			if ltyp == 1{
				switch {
				case x < la://brez == -1: //before load start
					vxs[i] += rl//0.0
					mxs[i] += rl * x//0.0
				case x >= la:
					vxs[i] += rl - w
					mxs[i] += rl*x - w*(x1)
				}
			}
			switch {
			case x <= la:
				vxs[i] += rl
				mxs[i] += rl*x
			case x >= la + ldcvr: //after load cover
				switch ltyp{
					case 4:
					if wa >= wb {
						vxs[i] += rl - W
						mxs[i] += rl*x - W*(x-ldcx)
					} else {
						vxs[i] += rl - W
						mxs[i] += rl*x - W*ldcx
					}
					case 3:
					vxs[i] += rl - W
					mxs[i] += rl*x - W*(x-ldcx)
				}
				
			default: //within load cover
				switch ltyp {
				case 3:
					w1 := w * x1
					ldcx2 := la + x1/2.0
					vxs[i] += rl - w1
					mxs[i] += rl*x - w1*(x-ldcx2)
				//above should work even if ldcx < x? who (or what) knows
				case 4:
					wm := math.Abs(wb - wa)/(l - lb - la)
					dw := wm * x1
					if wa >= wb{	
						wx := wa - dw
						w1 := (wx * x1) + (dw * x1 / 2.0)
						ldcx2 := (ldcvr)*(wa+2.0*wx)/(3.0*(wa+wx))
						//this is ze distance from side wa, so dist from x = x - la - 
						ldcx2 = x - la - ldcx2
						vxs[i] += rl - w1 
						//mxs[i] -= w1*(la + ldcx) //CHECK THIS EH
						mxs[i] += rl*x - w1*ldcx2

					} else {
						wx := wa + dw
						w1 := (wa * x1) + (dw * x1 / 2.0)
						ldcx2 := (ldcvr)*(wx+2.0*wa)/(3.0*(wa+wx))
						vxs[i] += rl - w1
						mxs[i] += rl*x - w1*ldcx2
					}
				}
			}
		}
		
	}
	//calc unit moment vector
	mdxs := make([]float64, ndiv)
	vdxs := make([]float64, ndiv)
	dxs := make([]float64, ndiv)
	var homerx, mx, vx float64
	for i := 0; i < ndiv; i++ {
		zd := float64(i) * div
		var mdxk, vdxk float64
		for x := 0; x < ndiv; x++ {
			z := float64(x) * div
			//fmt.Println("i",i,"x",x,mxs[x],mx,homerx)
			switch {
			case x == 0:
				homerx = 1
			case x == ndiv -1:
				homerx = 1
			case x%2 == 0:
				//homerx = 4
				homerx = 2
			case x%2 != 0:
				//homerx = 2
				homerx = 4
			}
			if z <= zd {
				mx = z * (l - zd) / l
				vx = (l - zd) / l
			} else {
				mx = zd * (l - z) / l
				vx = -zd / l
			}
			//fmt.Println(mxs[x],mx,homerx)
			mdxk += mxs[x] * mx * homerx
			vdxk += vxs[x] * vx * homerx
		}
		mdxs[i] = mdxk * (l/nd) / (3.0 * e * iz)
		g := e / (2 * 1.15)
		if ar > 0.0 {vdxs[i] = 1.5 * vdxk * (l / nd) / (3.0 * g * ar)}
		dxs[i] = mdxs[i] + vdxs[i]
		//fmt.Println("i",i,"mdx",mdxs[i],"vdx",vdxs[i])
	}
	var mmax, vmax, vmaxx, mmaxx, dmax, dmaxx, mpmax, mpmaxx float64
	var cfxs []float64
	for i, x := range xs {
		if i != 0 && mxs[i]*mxs[i-1] <= 0.0 {
			cfxs = append(cfxs, x)
		}
		if math.Abs(vmax) < math.Abs(vxs[i]) {
			vmax = vxs[i]
			vmaxx = x
		}
		if math.Abs(mmax) < math.Abs(mxs[i]) {
			mmax = mxs[i]
			mmaxx = x
		}
		if math.Abs(dmax) < math.Abs(dxs[i]) {
			dmax = dxs[i]
			dmaxx = x
			dxs[i] = mdxs[i] + vdxs[i]
		}
		if mxs[i] > 0.0 && mpmax < mxs[i]{
			mpmax = mxs[i]
			mpmaxx = x
		}
	}
	//fmt.Println(s)
	var s string
	if plotbm {s = PlotBmSfBm(xs, vxs, mxs, dxs, l, true)}
	var rez BeamRez
	rez.Mem = member
	rez.SF = vxs
	rez.BM = mxs
	rez.Dxs = dxs
	rez.Dxs_BM = mdxs
	rez.Dxs_SF = vdxs
	rez.Maxs = []float64{vmax, mmax, dmax, mpmax}
	rez.Locs = []float64{vmaxx, mmaxx, dmaxx, mpmaxx}
	rez.Cfxs = cfxs
	rez.Txtplot = s
	rez.Xs = xs
	rez.Qf = []float64{-1, ml, -1, mr}
	return rez
}

//GetCritxs returns cfxs [l1, l2]
func GetCritxs(xs, vxs, mnxs, mpxs []float64, mtyp string) (cxs []float64, cs []int){
	var l1, l2, l3, l4 float64
	var i1, i2, i3, i4 int
	for i, x := range xs{
		mn := mnxs[i]
		mp := mpxs[i]
		
		if l1 == 0.0 && i != 0 {
			if mn >= 0.0 && i < len(xs)/2{
				l1 = x
				i1 = i
			}
		}
		if l2 == 0.0 && i > len(xs)/2{
			if mn >= 0.0 && i != len(xs)-1{
				l2 = x
				i2 = i
			}	
		}
		if l3 == 0.0 && i != 0{
			if mp <= 0.0 && i < len(xs)/2{
				l3 = x
				i3 = i
			}
		}
		
		if l4 == 0.0 && i > len(xs)/2{
			if mp <= 0.0 && i != len(xs)-1{
				l4 = x
				i4 = i
			}
		}
	}
	cxs = []float64{l1, l2, l3, l4}
	cs = []int{i1, i2, i3, i4}
	return
}

//PlotBmSfBm early func to plot bending moment and shear force rez
//wrote this can't delete it is now a part of me
func PlotBmSfBm(xs, vxs, mxs, dxs []float64, l float64, dumb bool) string {
	var dat, pltstr string
	for i, x := range xs {
		dat += fmt.Sprintf("%v %v %v %v\n", x, vxs[i], mxs[i], -1000.0*dxs[i])
	}
	//create temp files
	f, e1 := os.CreateTemp(".", "barf")
	
	if e1 != nil {
		log.Println(e1)
	}

	_, e1 = f.WriteString(dat)
	if e1 != nil {
		log.Println(e1)
	}
	var termstr string
	if dumb {
		termstr = "set term dumb ansi size 60,15 aspect 1; unset border;set colorsequence classic; set tics nomirror scale 0.5; set xlabel \"M\""
	} else {
		termstr = "set terminal qt persist; set output 'tst.png'"
	}
	prg := "gnuplot"
	arg0 := "-e"
	var s, arg1 string
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			setstr := "set autoscale; unset key; set title \"BEAM-SHEAR FORCE\"; set ylabel \"KN\"; set offsets graph 0.1 ,0.1 ,0.1 ,0.1; set xzeroaxis"
			pltstr = fmt.Sprintf("plot '%s' using 1:2 w lines lt 1", f.Name())
			arg1 = fmt.Sprintf("%s; %s; %s", termstr, setstr, pltstr)
			s += exec_command(prg, arg0, arg1)
		case 1:
			setstr := "set autoscale; unset key; set title \"BEAM-BENDING MOMENT\"; set ylabel \"KN M\"; set offsets graph 0.1 ,0.1 ,0.1 ,0.1; set xzeroaxis"
			pltstr = fmt.Sprintf("plot '%s' using 1:3 w lines lt 1", f.Name())
			arg1 = fmt.Sprintf("%s; %s; %s", termstr, setstr, pltstr)
			s += exec_command(prg, arg0, arg1)
		case 2:
			setstr := "set autoscale; unset key; set title \"BEAM-DEFLECTION\"; set ylabel \"MM\"; set offsets graph 0.1 ,0.1 ,0.1 ,0.1;set xzeroaxis"
			pltstr = fmt.Sprintf("plot '%s' using 1:4 w lines lt 1", f.Name())
			arg1 = fmt.Sprintf("%s; %s; %s", termstr, setstr, pltstr)
			s += exec_command(prg, arg0, arg1)
		}

	}
	//arg2 := "-persist"
	f.Close()
	os.Remove(f.Name())
	return s
}
