package barf

import (
	//"log"
	//"fmt"
	"math"
)

//SlbSdratBs is the span-depth ratio (sdrat) entry func for a (singly reinforced) slab
func SlbSdratBs(s *RccSlb, mspan, astreq, astprov, dused float64) (sdchk bool, sd, dserve float64){
	lspan := 0.0; endc := s.Endc
	tyb := 0.0; bw := 1000.0; bf := 1000.0; df := dused; effd := dused - s.Efcvr
	//fmt.Println("EFFD->",effd,dused,s.Efcvr)
	//WHAAT.the.
	if effd < 0{effd = 0}
	asc := 0.0; mtyp := 1
	switch s.Type{
		case 1:
		lspan = s.Lspan
		case 2:
		lspan = s.Lx
		switch s.Endc{
			case 10:
			endc = 1
			default:
			endc = 2
		}
		case 3:
		//ribbed slab
		lspan = s.Lspan
		bw = s.Bw; bf = s.Bf; df = s.Df; effd = s.Dw - s.Efcvr; tyb = 1.0
		mtyp = 2
	}
	lspan = lspan/1000.0
	sdchk, sd, dserve = SpanDepRatBs(s.Fck, s.Fy, lspan, tyb, bw, bf, df, effd, astprov, astreq, asc, mspan, s.DM, mtyp, endc)
	//fmt.Println(sdchk, sd, dserve)
	return
}

//BmSdratBs is the entry func for span-depth ratio checks for a beam
func BmSdratBs(b *RccBm, mspan float64) (sdchk bool, sd, dserve float64){
	//span depth ratio checks for a beam
	//mspan- max midspan moment (bmenv mpmax) or max clvr moment at support
	mtyp := 2
	effd := b.Dused - b.Cvrt
	var asc, astprov, astreq float64
	switch b.Endc{
		case 0:		
		astprov = b.Rbrc[4]
		astreq = b.Rbrc[5]
		asc = b.Ast
		if len(b.Rbrt) > 4{asc = b.Rbrt[4]}
		default:
		astprov = b.Ast; astreq = astprov
		if len(b.Rbrt) > 6{
			astprov = b.Rbrt[4]
			astreq = b.Rbrt[5]
		}
		asc = b.Asc
		if len(b.Rbrc) > 4{asc = b.Rbrc[4]}
	}
	sdchk, sd, dserve = SpanDepRatBs(b.Fck, b.Fy, b.Lspan, b.Tyb, b.Bw, b.Bf, b.Df, effd, astprov, astreq, asc, mspan, b.DM, mtyp, b.Endc)
	// if b.Verbose{
	// 	log.Println("checking beam id-",b.Id, b.Mid)
	// 	log.Println("b.Fck, b.Fy, b.Lspan, b.Tyb, b.Bw, b.Bf, b.Df, effd, astprov, astreq, asc, mspan, b.DM, mtyp, b.Endc-\n",b.Fck, b.Fy, b.Lspan, b.Tyb, b.Bw, b.Bf, b.Df, effd, astprov, astreq, asc, mspan, b.DM, mtyp, b.Endc)
	// 	log.Println("mspan, astprov, asc->",mspan, astprov, asc)
	// 	log.Println("sdchk,sd,dserve,effd->",sdchk, sd, dserve, effd)
	// }
	return
}

//SpanDepRatBs is THE BACKBONE of deflection limit/span depth ratio calcs
//hulse sec 7.1, span- effective depth calculations
func SpanDepRatBs(fck, fy, lspan, tyb, bw, bf, df, effd, astprov, astreq, asc, mur, dm float64, mtyp, endc int) (bool, float64, float64){
	//span depth ratio check as per BS 8110
	var sd float64
	var sdchk bool
	switch endc{
	case 0:
		//cantilever
		sd = 7.0
	case 1:
		//simply supported
		sd = 20.0
	case 2:
		//continuous
		sd = 26.0
	}
	//lspan = lspan/1000.0
	if lspan > 10.0 {sd = sd * 10/lspan}
	mfactor := 1.0 - dm
	//tension rebar mod
	fs := 5.0 * fy * astreq/(astprov* mfactor)/8.0
	var m1 float64
	switch tyb{
		case 0.0:
		m1 = mur * 1e6/bw/math.Pow(effd,2)
		default:
		//m1 = mur * 1e6/bw/math.Pow(effd,2)
		m1 = mur * 1e6/bf/math.Pow(effd,2)
	}	
	if m1 > 6.0 {m1 = 6.0}
	if m1 < 0.5 {m1 = 0.5}
	u1 := 0.55 + (477.0 - fs)/(120.0*(0.9+m1))
	if u1 > 2.0 {u1 = 2.0}
	
	//compression rebar mod
	var r0 float64
	switch tyb{
		case 0.0:
		r0 = 100.0 * asc/(bw* effd)
		default:
		r0 = 100.0 * asc/(bf* effd)
	}
	
	if r0 > 3.0 {r0 = 3.0}
	u2 := 1.0 + r0/(3.0 + r0)
	
	sd = sd * u1 * u2
	
	if lspan*1000.0/effd <= sd {sdchk = true}

	if mtyp == 2 && tyb > 0.0 {
		//flanged beam
		if bw/bf < 0.3 {
			sd = sd * 0.8
		} else{
			sd = sd * (0.8 + 0.2/0.7 * (bw/bf - 0.3))
		}
	}
	return sdchk, sd, lspan*1000.0/sd 
}

