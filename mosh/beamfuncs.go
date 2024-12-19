package barf

import (
	//"fmt"
	"math"
	kass"barf/kass"
)

//beam helper funcs

//var Bsk8 float64 = 0.45 //concrete stress factor
//var Bsk9 float64 = 0.9  //neutral axis depth

const floatEqBm = 0.0001

//almostEqualBm checks for floating point equality
func almostEqualBm(a, b float64) bool {
	if a == b {
		return true
	} else if a == 0 || b == 0 || (math.Abs(a)+math.Abs(b)) < floatEqBm {
		return math.Abs(a)-math.Abs(b) < floatEqBm*math.Float64frombits(0x0010000000000000)
	}
	return math.Abs(a-b)/(math.Abs(a)+math.Abs(b)) < floatEqBm
}

//lerp interpolates y at x given slices of ys and xs
func lerp(xs, ys []float64, xv float64) (yv float64){
	/*
	   note- only works if xs is sorted in ascending order
	*/
	for i, x := range xs{
		if i == len(xs)-1 && x < xv{yv = -6.66}
		if x == xv{
			yv = ys[i]
			return
		}
		if x > xv{
			switch i{
				case 0:
				yv = ys[i]
				default:
				dy := (ys[i] - ys[i-1])/(xs[i] - xs[i-1])
				dx := xv - xs[i-1]
				yv = ys[i-1] + dx * dy
				break
			}
		}
	}
	return
}

//rebarforce returns the force in rebar given stlarm (distance of rebar from top fiber)
func rebarforce(fy, stlarm float64) (fst float64) {
	emsteel := 200000.0
	eyield := (fy / 1.15) / emsteel
	
	est := 0.0035 * stlarm

	if est >= eyield {
		fst = fy / 1.15
	} else {
		fst = emsteel * est
	}
	return
}

//BmArXu returns the compressive force in the concrete block
func BmArXu(b *RccBm, k8, k9, x float64, rblk bool) (cck, area, xc, yc float64){
	if !rblk{
		//rectangular parabolic stress block (yeah right)
	}
	//rect stress block
	switch b.Npsec{
		case false:	
		if b.Tyb == 0.0{
			cck = k8 * b.Fck * b.Bw * k9 * x
		} else {
			cck = k8 * b.Fck * b.Bf * k9 * x
			if k9 * x > b.Df{
				cck = k8*b.Fck*b.Bf*b.Df + k8*b.Fck*b.Bw*(k9*x-b.Df)
			}
		}
		yc = k9 * x/2.0
		case true:
		area, xc, yc = ColSecArXu(b.Sec, k9 * x)
		cck = k8 * b.Fck * area
		yc = b.Sec.Ymx - yc
	}
	return
}

//BmRbrFrc returns stresses in steel bars for depth of neutral axis x
func BmRbrFrc(b *RccBm, x float64, bscode bool)(fst, fsc float64){
	//THIS IS NOT WERK FOR IS CODE
	for i, dia := range b.Dias{
		dbar := b.Dbars[i]
		if x >= dbar{
			switch bscode{
				case true:
				fsc += RbrArea(dia) * rebarforce(b.Fy, (x-dbar)/x)
				case false:
				fsc += RbrArea(dia) * RbrFrcIs(b.Fy, 0.0035 * (x-dbar)/x)
			}
			//subtract for added steel area in kompress
			//fsc -= b.Fck
		} else {
			switch bscode{
				case true:
				fst += RbrArea(dia) * rebarforce(b.Fy, (dbar-x)/x)
				case false:
				fst += RbrArea(dia) * RbrFrcIs(b.Fy, 0.0035 * (dbar-x)/x)
			}
		}
	}
	return
}

//EDist returns the euclidean distance between jb and je
func EDist(jb, je []float64) (l float64){
	//WHY HERE TOO euc. distance between jb and je
	switch len(jb){
		case 1:
		l = math.Abs(je[0] - jb[0])
		case 2:
		l = math.Sqrt(math.Pow(je[0]-jb[0],2)+math.Pow(je[1]-jb[1],2))		
		case 3:
		l =  math.Sqrt(math.Pow(je[0]-jb[0],2)+math.Pow(je[1]-jb[1],2)+math.Pow(je[2]-jb[2],2))
	}
	return
}

//FckEm does what it does in kass too -> FCKEM
func FckEm(fck float64) (float64){
	//returns in KN/m2
	return 5e6 * math.Sqrt(fck) 
}

//getbflange returns the breadth of flange
func getbflange(code, mrel int, tyb, lspan, df, bw float64) (bf float64){
	//konvert lspan from meters to mm (MAYBE CHANGE LSPANS TO MM)
	lspan = lspan * 1e3
	if mrel == 0{lspan = 0.7 * lspan}
	switch code{
		case 1:
		//is code
		bf = math.Round(tyb*(6.0*df + lspan/6.0) + bw)
		case 2:
		//bs code
		bf = math.Round(bw + (lspan * tyb)/5.0)
		case 3:
		//calc using aci
		
	}
	return
}

//BeamDims generates dims (???)
//i have no idea why i wrote this
func BeamDims(styp, endc, ldcalc int, din []float64, lspan, efcvr float64) (dims []float64, wdl float64, err error){
	if efcvr == 0 {efcvr = 25.0 + 10.0 + 5.0 }
	switch styp{
		case 6:
		//"T"
		case 7:
		//generic L
		case 8:
		//l left -|
		case 9: 
		//l right |-
		case 1:
		var b, d, ldrat float64
		if len(din) == 2{
			dims = make([]float64, len(din))
			copy(dims, din)
			if ldcalc == 0{
				return
			} else {
				wdl = 25.0 * din[0] * din[1] * 1e-6
				return
			}
		}
		if len(din) == 0{
			b = 230.0
		} else {
			b = din[0]
		}
		switch endc{
			case 0:
			ldrat = 7.0
			case 1:
			ldrat = 12.0
			case 2:
			ldrat = 15.0
		}
		//DIS IS ok ish
		d = math.Round(lspan*1e3/ldrat)
		d += efcvr
		if d < 230.0{d = 230.0}
		dims = []float64{b,d}
		if ldcalc == 1{
			//wdl in KN/m
			wdl = 25.0 * b * d * 1e-6
		}
	}
	return
}

//BeamSecGen calls BeamDims to generate dims and cp slices
func BeamSecGen(styp, endc, ldcalc,nd int, din []float64, lspan, efcvr float64) (dims, cpvec []float64, wdl float64, err error) {
	//all dims here in MM. MM.
	//TODO- add moar sections
	switch styp{
		case 1:
		//rect
		dims, wdl, err = BeamDims(styp, endc, ldcalc, din, lspan, efcvr)
		case 2:
		//l - CHANGE styp
		case 3:
		//tee - CHANGE styp
		case 4:
		//pocket - TO DO - INCLUDE styp
	}
	if err != nil{return}
	bar := kass.CalcSecProp(styp, dims)
	switch nd{
		case 1:
		//1d beam
		cpvec = []float64{bar.Ixx*1e-12, bar.Area*1e-6}
		case 2:
		//2d frame
		cpvec = []float64{bar.Ixx*1e-12, bar.Area*1e-6}
	}
	return
}

/*
func BeamCp(ftyp, styp, ldcalc int, din []float64) (cp []float64, wdl float64){
	//DIN IS IN METERS GODDAMN METERS
	//returns cp for 1d beam
	switch ftyp{
		case 1:
		// 1-d beam
		bar := kass.CalcSecProp(styp, din)
		if ldcalc > 0{
			wdl = 25.0 * bar.Area
		}
		cp = []float64{bar.Ixx}
	}
	return
}
*/
