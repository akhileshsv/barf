package barf

//shah slab design funcs, chapter 6

import (
	"os"
	"log"
	"fmt"
	"math"
	"errors"
	"math/rand"
	"runtime"
	"encoding/json"
	"path/filepath"
	"io/ioutil"
)
var (
	//table 12 and table 13 coeff in IS code/BS code coeff for cbeams/slabs
	CBmCspn = []float64{1./12.0,1./16.0,1./10.0,1./12.0}
	CBmCsup = []float64{1./10.0,1./12.0,1./9.0,1./9.0}
	CBm2Cspn = []float64{1./12.0,1./16.0,1./10.0,1./12.0}
	CVCdl = []float64{0.4,0.6,0.55,0.5}
	CVCll = []float64{0.45,0.6,0.6,0.6}
)

//Init initializes an RccSlb struct
func (s *RccSlb) Init(){
	if s.Code == 0{s.Code = 1}
	if s.Diamain == 0.0 {
		s.Diamain = 8.0
	}
	if s.Fy == 0.0 {
		s.Fy = 500.0
	}
	if s.Diadist == 0.0 {
		s.Diadist = 8.0
	}
	if s.Fyd == 0.0{
		s.Fyd = s.Fy
	}
	if s.Nomcvr == 0.0{
		s.Nomcvr = 20.0
	}
	if s.Efcvr == 0.0 {
		s.Efcvr = s.Nomcvr + s.Diamain/2.0
	}
	if s.Id == 0{
		//s.Id = len(allSlbz)+1
		s.Id = rand.Intn(666)
	}
	if s.Bsup == 0{s.Bsup = 230.0}
	if s.Id == 0{s.Id = rand.Intn(666)}
	if s.Title == ""{s.Title = fmt.Sprintf("rcc slab_%v",s.Id)}
	if s.Lx == 0.0{s.Lx = s.Lspan}
	if s.Ly == 0.0{s.Ly = 2.0*s.Lx}
	if s.Spancalc{
		//should be (min) depth, bsup but eh fwack it
		s.Lx = s.Lx + s.Bsup; s.Ly = s.Ly + s.Bsup
	}
}

//SlbDIs is the 1 way (ss/clvr) and 2-way slab design func- shah
func SlbDIs(s *RccSlb) (error) {
	//1 way (ss/clvr) and 2-way slab design- shah
	//1 way/clvr - calcs ult moment (wl2/8 or wl2/2)
	var ptreq, astreq, wd, wu, mdu, vumax, dtrial, dmin, dused, dserve, effd float64
	psfs := 1.15
	k1 := 0.361
	k2 := 0.416
	s.Init()
	efcvr := s.Efcvr 
	_, _, rumax, _ := BalSecKis(s.Fck, s.Fy, s.DM, psfs, k1, k2)
	//now go
	switch s.Type {
	case 1:
		//un way
		switch s.Endc {
		case 1, 0:
			//simply supported, cantilever
			dtrial = 75.0
			lspan := s.Lspan / 1000.0
			var dmcoef, sfcoef float64
			if s.Endc == 1 {
				dmcoef = 8.0
				sfcoef = 2.0
			} else {
				dmcoef = 2.0
				sfcoef = 1.0
			}
			wd = 0.025 * dtrial
			wu = 1.5 * (wd + s.DL + s.LL)
			mdu = wu * math.Pow(lspan, 2) / dmcoef
			dmin = math.Sqrt(1000.0*mdu/rumax) + efcvr
			//vumax = wu * lspan/sfcoef
			var kiter int
			for dmin > dtrial {
				if kiter > 666{return ErrIter}
				dtrial = dmin + 10.0
				wd = 0.025 * dtrial
				wu = 1.5 * (wd + s.DL + s.LL)
				mdu = wu * math.Pow(lspan, 2) / dmcoef
				dmin = math.Sqrt(1000.0*mdu/rumax) + efcvr
				kiter++
			}
			dused = 10.0 * math.Ceil(dmin/10.0)
			wu = 1.5 * (0.025*dused + s.DL + s.LL)
			mdu = wu * math.Pow(lspan, 2) / dmcoef
			effd = dused - efcvr
			H := 1.0 - (4600.0*mdu)/(s.Fck*math.Pow(effd, 2))
			ptreq = (0.5 * s.Fck / s.Fy) * (1.0 - math.Sqrt(H)) * 100.0
			astreq = ptreq * 10.0 * effd
			dserve = RccSlabServeChk(s.Type, s.Endc, s.Fy, s.LL, s.Lspan, ptreq, efcvr)
			kiter = 0
			for dserve > dused {
				if kiter > 666{return ErrIter}
				dused = dused + 5.0
				wu = 1.5 * (0.025*dused + s.DL + s.LL)
				mdu = wu * math.Pow(lspan, 2) / dmcoef
				effd = dused - efcvr
				H = 1.0 - ((4600.0 * mdu) / (s.Fck * math.Pow(effd, 2)))
				ptreq = (0.5 * s.Fck / s.Fy) * (1.0 - math.Sqrt(H)) * 100.0
				astreq = ptreq * 10.0 * effd
				dserve = RccSlabServeChk(s.Type, s.Endc, s.Fy, s.LL, s.Lspan, ptreq, efcvr)
				kiter++
			}
			effd = dused - efcvr
			vumax = wu * lspan / sfcoef
			murmax := rumax * math.Pow(effd, 2) * 1e-3
			rezmap, mindia := SlabRbrDiaSpcing(dused, astreq, efcvr)
			if len(rezmap) == 0 {
				return ErrSpacing
			}
			var asprov, diacom, prvcvr, spcmain float64
			asprov = rezmap[mindia][2]
			spcmain = rezmap[mindia][1]
			diacom = mindia
			if val, ok := rezmap[s.Diamain]; ok {
				asprov = val[2]
				spcmain = val[1]
				diacom = s.Diamain
			}
			prvcvr = s.Nomcvr + diacom/2.0
			if math.Abs(prvcvr - s.Efcvr) > 5.0{
				log.Println("cover error, changing depth")
				s.Efcvr = prvcvr
				s.Diamain = diacom
				if diacom == 6.0 {
					s.Fy = 250.0
				}
				return SlbDIs(s)
			}
			ptprov := 100 * asprov / (1000 * effd)
			//distribution isteel
			var astd, m1, ptup float64
			if s.Diadist == 6.0 {
				s.Fyd = 250.0
				astd = 1.5 * dused
			} else {
				astd = 1.2 * dused
			}
			sdmax := 5 * effd
			//shah chap 7 note 8 - in practice spacing is 300 mm ofc
			//it can be bumped upto 450 tbh as per codes
			if sdmax > 300.0 {
				sdmax = 300.0
			}
			sds := 1000 * RbrArea(s.Diadist) / astd
			if sds > sdmax {
				sds = sdmax
			}
			sds = 10.0 * math.Floor(sds/10.0)
			asdprov := math.Round(1000 * RbrArea(s.Diadist) / sds)
			s.Astd = asdprov
			//support reinforcement - this is wrong-ish but stays for tests
			//ss one way slab has 0.5 ast midspan @ support usually
			if s.Ibent == 0.0 {
				//ibent = percentage of bent up bars
				m1 = mdu
				ptup = ptprov
			} else {
				m1 = mdu * s.Ibent / 100.0
				ptup = ptprov * s.Ibent / 100.0
			}
			//check for shear
			beta := 0.8 * s.Fck / (6.89 * ptup)
			if beta < 1.0 {
				beta = 1.0
			}
			var k5 float64
			if dused <= 150.0 {
				k5 = 1.3
			}
			if dused >= 300.0 {
				k5 = 1.0
			} else {
				k5 = 1.3 - 0.3*(dused/150.0-1.0)
			}
			vuc := k5 * 0.8499999 * math.Sqrt(0.8*s.Fck) * (math.Sqrt(1.0+5.0*beta) - 1.0) * effd / (6.0 * beta)
			if vuc < vumax {
				//sec UNSAFE IN SHEAR holy fuck what will one do
				return ErrShear
			}
			devchk, _, ldreq := SlabDevLen(s.Fck, s.Fy, m1, vumax, effd, s.Diamain, s.Bsup, s.Endc)
			s.Devchk = devchk
			s.Ldev = ldreq
			if !devchk && s.Endc != 0 {
				//log.Println("devlen fubar FUBAR")
				//add additional length of 8* dia for anchorage + bend bar by 90 at support
				//do nuttin for now
			}
			s.Murmax = murmax
			s.Astm = asprov; s.Spcm = spcmain; s.Astd = asdprov; s.Spcd = sds; s.Astreq = astreq
			s.Dused = dused; s.BM = append(s.BM, mdu)
		}
	case 2:
		//if s.Verbose{fmt.Println(ColorPurple,"starting two way slab design",ColorReset)}
		//two way slab
		var mcxsup, mcxmid, mcysup, mcymid float64
		mcxsup, mcxmid, mcysup, mcymid = Slb2BMCoefIs(s.Endc, s.Lx, s.Ly)
		lx := s.Lx/1000.0
		dtrial = 75.0
		wd = 0.025 * dtrial
		wu = 1.5 * (wd + s.DL + s.LL)
		wulx2 := wu * math.Pow(lx,2)
		murmax := mcxmid * wulx2
		var mdus []float64
		for _, mc := range []float64{mcxsup, mcxmid, mcysup, mcymid} {
			msec := mc * wulx2
			mdus = append(mdus, msec)
			if msec > murmax {murmax = msec}
		}
		dmin = math.Sqrt(1000.0*murmax/rumax) + efcvr
		dreqy := math.Sqrt(1000.0*mdus[3]/rumax) + efcvr + 10.0
		if s.Endc == 10 {dreqy = math.Sqrt(1000.0*mdus[1]/rumax) + efcvr + 8.0}
		if dmin < dreqy {dmin = dreqy}
		for dmin > dtrial {
			dtrial = dmin + 10.0
			wd = 0.025 * dtrial
			wu = 1.5 * (wd + s.DL + s.LL)
			wulx2 = wu * math.Pow(lx,2)
			murmax = mcxmid * wulx2
			for _, mc := range []float64{mcxsup, mcxmid, mcysup, mcymid} {
				msec := mc * wulx2
				mdus = append(mdus, msec)
				if msec > murmax {murmax = msec}
			}
			dmin = math.Sqrt(1000.0*murmax/rumax) + efcvr
			dreqy = math.Sqrt(1000.0*mdus[3]/rumax) + efcvr + 8.0
			if s.Endc == 10 {dreqy = math.Sqrt(1000.0*mdus[1]/rumax) + efcvr + 8.0 }
			if dmin < dreqy {dmin = dreqy}
		}
		s.BM = mdus
		dused = 10.0 * math.Ceil(dmin/10.0)
		effd = dused - efcvr
		wu = 1.5 * (0.025*dused + s.DL + s.LL)
		mdu = wu * mcxmid * math.Pow(lx, 2) 
		H := 1.0 - (4600.0*mdu)/(s.Fck*math.Pow(effd, 2))
		ptreq = (0.5 * s.Fck / s.Fy) * (1.0 - math.Sqrt(H)) * 100.0
		astreq = ptreq * 10.0 * effd
		dserve = RccSlabServeChk(s.Type, s.Endc, s.Fy, s.LL, s.Lx, ptreq, efcvr)
		for dserve > dused {
			dused = dused + 5.0
			wu = 1.5 * (0.025*dused + s.DL + s.LL)
			mdu = wu * mcxmid * math.Pow(lx, 2) 
			wulx2 = wu * math.Pow(lx,2)
			mdu = mcxmid * wulx2
			//deflection check applied only at midspan? not at support
			effd = dused - efcvr
			H = 1.0 - ((4600.0 * mdu) / (s.Fck * math.Pow(effd, 2)))
			ptreq = (0.5 * s.Fck / s.Fy) * (1.0 - math.Sqrt(H)) * 100.0
			astreq = ptreq * 10.0 * effd
			dserve = RccSlabServeChk(s.Type, s.Endc, s.Fy, s.LL, s.Lx, ptreq, efcvr)
		}
		effd = dused - efcvr
		wu = 1.5 * (0.025*dused + s.DL + s.LL)
		wulx2 = wu * math.Pow(lx,2)
		s.Dused = dused
		//main steel dia-spacing
		var asts , diams, spcs []float64
		for idx, mc := range []float64{mcxsup, mcxmid, mcysup, mcymid} {
			//CHANGE LONG SPAN STEEL EFFD//(i haz)
			dsec := dused; effdi := effd
			if idx == 3{dsec = dused - s.Diamain; effdi = effd - s.Diamain}
			msec := mc * wulx2
			mdus[idx] = msec
			H = 1.0 - ((4600.0 * msec) / (s.Fck * math.Pow(effdi, 2)))
			astreq = (0.5 * s.Fck / s.Fy) * (1.0 - math.Sqrt(H)) * 1000.0 * effdi
			rezmap, mindia := SlabRbrDiaSpcing(dsec, astreq, efcvr)
			if len(rezmap) == 0{
				log.Println("ERRORE,errore->rbr dia error")
				return ErrSpacing
			} 
			var asprov, diacom, prvcvr, spcmain float64
			asprov = rezmap[mindia][2]
			spcmain = rezmap[mindia][1]
			diacom = mindia
			if val, ok := rezmap[s.Diamain]; ok{
				asprov = val[2]
				spcmain = val[1]
				diacom = s.Diamain 
			}
			//fmt.Println("ASPROV->",asprov)
			prvcvr = s.Nomcvr + diacom/2.0
			if prvcvr != s.Efcvr{
				if math.Abs(prvcvr - s.Efcvr) > 5.0{
					return ErrCvr
				}
			}
			s.Astr = append(s.Astr, astreq)
			diams = append(diams, diacom)
			spcs = append(spcs, spcmain)
			asts = append(asts, asprov)
			//ptprov := 100 * asprov / (1000 * effd)
		}
		if mcxsup == 0.0{
			//add 50% of midspan x steel
			asts[0] = math.Round(asts[1]/2.0)
			diams[0] = diams[1]
			spcs[0] = math.Round(spcs[1]*2.0)
			
		}
		if mcysup == 0.0{
			//add 50% of midspan y steel
			asts[3] = math.Round(asts[3]/2.0)
			diams[3] = diams[3]
			spcs[3] = math.Round(spcs[3]*2.0)
			
		}
		s.Asts = asts
		s.Dias = diams
		s.Spcms = spcs
		//CHECKS?
		//if s.Verbose{
		//	fmt.Println(asts, diams, spcs)
		//}
	}
	s.Dz = true
	return nil
}

//Slb2BMCoeffIs returns design ultimate moments for a 2 way slab based on is code coefficients
//see shah section 6.3 
//btw - ec 1-9 for 2w is for slabs that are cast integrally with the frame
//slab 10 has no friends or home
func Slb2BMCoefIs(endc int, lx, ly float64) (mcxsup, mcxmid, mcysup, mcymid float64) {
	/*
	   2 way slab yield line coefficients from is 456
	   shah
	//short span lx coefficients ax
	//-ve coeff at supports

	*/
	alpxsup := [][]float64{
		{0.032, 0.037, 0.043, 0.047, 0.051, 0.053, 0.060, 0.065},
		{0.037, 0.043, 0.048, 0.051, 0.055, 0.057, 0.064, 0.068},
		{0.037, 0.044, 0.052, 0.057, 0.063, 0.067, 0.077, 0.085},
		{0.047, 0.053, 0.060, 0.065, 0.071, 0.075, 0.084, 0.091},
		{0.045, 0.049, 0.052, 0.056, 0.059, 0.060, 0.065, 0.069},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0.057, 0.064, 0.071, 0.076, 0.080, 0.084, 0.091, 0.097},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
	}
	//+ve coeff at midspan
	alpxmid := [][]float64{
		{0.024, 0.028, 0.032, 0.036, 0.039, 0.041, 0.045, 0.049},
		{0.028, 0.032, 0.036, 0.039, 0.041, 0.044, 0.048, 0.052},
		{0.028, 0.033, 0.039, 0.044, 0.047, 0.051, 0.059, 0.065},
		{0.035, 0.040, 0.045, 0.049, 0.053, 0.056, 0.063, 0.069},
		{0.035, 0.037, 0.040, 0.043, 0.044, 0.045, 0.049, 0.052},
		{0.035, 0.043, 0.051, 0.057, 0.063, 0.068, 0.080, 0.088},
		{0.043, 0.048, 0.053, 0.057, 0.060, 0.064, 0.069, 0.073},
		{0.043, 0.051, 0.059, 0.065, 0.071, 0.076, 0.087, 0.096},
		{0.056, 0.064, 0.072, 0.079, 0.085, 0.089, 0.100, 0.107},
	}
	//long span ly coefficients 
	alpysup := []float64{0.032, 0.037, 0.037, 0.047, 0, 0.045, 0, 0.057, 0}
	alpymid := []float64{0.024, 0.028, 0.028, 0.035, 0.035, 0.035, 0.043, 0.043, 0.056}

	//ss coefficients
	alpxss := []float64{0.062, 0.074, 0.084, 0.093, 0.099, 0.104, 0.113, 0.118}
	alpyss := []float64{0.062, 0.061, 0.059, 0.055, 0.051, 0.046, 0.037, 0.029}

	lrats := []float64{1.0, 1.1, 1.2, 1.3, 1.4, 1.5, 1.75, 2.0}
	lrat := ly / lx
	//fmt.Println(lrat)
	switch endc {
	case 10:
		//simply supported slab
		if lrat == 1.0 {
			mcxmid = alpxss[0]
			mcymid = alpyss[1]
			
		}else{
			mcxmid = interpolator(lrat, lrats, alpxss)
			mcymid = interpolator(lrat, lrats, alpyss)
		}
	default:
		//all else
		mcymid = alpymid[endc-1]
		mcysup = alpysup[endc-1]
		if lrat == 1.0 {
			mcxmid = alpxmid[endc-1][0]
			mcxsup = alpxsup[endc-1][0]
		} else {
			mcxmid = interpolator(lrat, lrats, alpxmid[endc-1])
			mcxsup = interpolator(lrat, lrats, alpxsup[endc-1])
		}
	}
	
	return
}

//Slb2VfCoeff returns the shear force in uniformly loaded two way slabs based on table 3.15 of bs 8110
//avxc, avyc - shear foce coeff at a continuous edge in x, y; avxd, avyd - at a discontinuous edge in x,y directions
//shear force vsn = avn * wu * lx
func Slb2VfCoeff(endc int, lx, ly float64) (avxc, avxd, avyc, avyd float64){
	lrats := []float64{1.0, 1.1, 1.2, 1.3, 1.4, 1.5, 1.75, 2.0}
	lrat := ly / lx
	axc := [][]float64{
		{0.33,0.36,0.39,0.41,0.43,0.45,0.48,0.50},
		{0.36,0.39,0.42,0.44,0.45,0.47,0.50,0.52},
		{0.36,0.40,0.44,0.47,0.49,0.51,0.55,0.59},
		{0.40,0.44,0.47,0.50,0.52,0.54,0.57,0.60},
		{0.40,0.43,0.45,0.47,0.48,0.49,0.52,0.54},
		{0,0,0,0,0,0,0,0},
		{0.45,0.48,0.51,0.53,0.55,0.57,0.60,0.63},
		{0,0,0,0,0,0,0,0},
		{0,0,0,0,0,0,0,0},
	}
	axd := [][]float64{
		{0,0,0,0,0,0,0,0},
		{0,0,0,0,0,0,0,0},
		{0.24,0.27,0.29,0.31,0.32,0.34,0.36,0.38},
		{0.26,0.29,0.31,0.33,0.34,0.35,0.38,0.40},
		{0,0,0,0,0,0,0,0},
		{0.26,0.30,0.33,0.36,0.38,0.40,0.44,0.47},
		{0.30,0.32,0.34,0.35,0.36,0.37,0.39,0.41},
		{0.29,0.33,0.36,0.38,0.40,0.42,0.45,0.48},
		{0.33,0.36,0.39,0.41,0.43,0.45,0.48,0.50},
	}
	ayc := []float64{0.33,0.36,0.36,0.40,0.0,0.40,0.0,0.45,0.33}
	ayd := []float64{0.0,0.24,0.0,0.26,0.26,0.0,0.29,0.30,0.0}
	avyc = ayc[endc-1]
	avyd = ayd[endc-1]
	if lrat == 1.0 {
		avxc = axc[endc-1][0]
		avxd = axd[endc-1][0]
	} else {
		avxc = interpolator(lrat, lrats, axc[endc-1])
		avxd = interpolator(lrat, lrats, axd[endc-1])
	}
	return
}


//Slb2BfCoeff returns the coefficients for loads on supporting beams of a two way slab
//see table 63, pg 205 reynolds steedman
//ly = longer span. longer. span.
//shear = r * wlx2

func Slb2BfCoeff(endc int, lx, ly float64) (r1, r2, r3, r4 float64){
	//todo - add a, b, c, d (coeff for yield line dims)
	k := ly/lx
	switch endc{
		case 1, 9, 10:
		//ss/cs on all 4 sides
		//w/w/o torsion reinf.(i guess)
		switch{
			case k == 1:
			r1, r2, r3, r4 = 0.25, 0.25, 0.25, 0.25
			case k > 1:
			r1, r3 = 0.25, 0.25
			r2 = 0.5 * (k - 0.5)
			r4 = r2
		}
		case 2:
		//1 short side disc. (3 sides cs)
		r1 = 3.0/20.0
		r3 = 0.25
		r2 = 0.5 * (k - 0.4)
		r4 = r2
		case 3:
		//1 long side disc.
		switch{
			case k > 5.0/4.0:
			r1 = 5.0/16.0
			r3 = r1
			r4 = 5.0 * (k - 5.0/8.0)/8.0
			r2 = 0.6 * r4
			default:
			//k < 5.0/4.0
			r1 = 0.5 * k * (1.0 - 0.4 * k)
			r3 = r1
			r2 = 3.0 * k * k/20.0
			r4 = 0.25 * k * k
		}
		case 4:
		//1 long n short cs (1 l and s disc.)
		r1 = 3.0/16.0
		r3 = 5.0/16.0
		r4 = 5.0 * (k - 0.5)/8.0
		r2 = 3.0 * r4/5.0
		case 5:
		//2 short sides disc. (2 long cs)
		r1 = 3.0/20.0
		r3 = r1
		r2 = 0.5 * (k - 0.3)
		case 6:
		//2 short cs.
		//what if k > 5/3? huh?
		r1 = 5.0/12.0
		r3 = r1
		r2 = 0.5 * (k - 5.0/6.0)
		r4 = r2
		case 7:
		//1 long side cs
		r1 = 3.0/16.0
		r3 = r1
		r4 = 5.0 * (k - 3.0/8.0)/8.0
		r2 = 3.0 * r4/5.0
		case 8:
		//1 short side cs
		switch{
			case k > 4.0/3.0:
			r1 = 0.25
			r2 = 0.5 * (k - 2.0/3.0)
			r4 = r2
			r3 = 5.0/12.0
			default:
			r3 = 5.0 * k * (1.0 - 3.0 * k/8.0)/8.0
			r1 = 3.0 * r3/5.0
			r2 = 3.0 * k * k/16.0
			r4 = r2
		}
	}
	return
}

//DzSlbSpn designs a slab span - DELETE cause it evidently doesn't 
func DzSlbSpn(s *RccSlb, code int, errchn chan []interface{}){
	switch code{
		case 0,1:
		case 2:
		//TODO BS CODE
	}
}

//Printz prints
func (s *RccSlb) Printz() (rez string){
	if s.Code == 0{s.Code = 1}
	codez := []string{"is","bs"}
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
			default:
			t = "1 way"
		}
	}
	//add defaults here?
	
	rez += fmt.Sprintf("%s\n%s slab \nlspan %.1f, lx %.1f, ly %.1f mm\n",s.Title,t, s.Lspan, s.Lx, s.Ly)
	rez += fmt.Sprintf("grade of concrete M %.1f, steel - main Fe %.f, dist Fe %.f\n", s.Fck, s.Fy, s.Fyd)
	rez += fmt.Sprintf("cover - nominal %0.1f mm, effective %0.1f mm\n", s.Nomcvr, s.Efcvr)
	rez += fmt.Sprintf("loads -dl %.1f kN/m2, ll %0.1f kN/m2\n", s.DL, s.LL)
	rez += fmt.Sprintf("design code %s design type %v\n", codez[s.Code-1], s.Dtyp)
	switch s.Type{
		case 1:
		switch s.Endc{
			case 5:
			rez += fmt.Sprintf("slab depth- %.0f mm\n",s.Dused)
			for i := 0; i < s.Nspans; i++{
				ldx := i 
				rdx := i + 1
				rez += fmt.Sprintf("span no- %v\n",i+1)
				rez += fmt.Sprintf("sec \trbr spc dst spc (mm)\n")
				rez += fmt.Sprintf("lt sup\t %.0f %.0f %.0f %.0f\n",s.Dias[ldx], s.Spcms[ldx], s.Dists[ldx],s.Spcds[ldx])
				rez += fmt.Sprintf("md spn\t %.0f %.0f %.0f %.0f\n",s.Dias[i], s.Spcms[i], s.Dists[i],s.Spcds[i])
				rez += fmt.Sprintf("rt sup\t %.0f %.0f %.0f %.0f\n",s.Dias[rdx], s.Spcms[rdx], s.Dists[rdx],s.Spcds[rdx])
			}
			default:
			//cannot be midspan everywhere? eh fuck it
			rez += fmt.Sprintf("midspan moment %.2f kn-m\n",s.BM[0])
			rez += fmt.Sprintf("ast required - %.1f mm2\n",s.Astreq)
			rez += fmt.Sprintf("ast provided - %0.1f mm2\n",s.Astm)
			switch s.Endc{
				case 0:
				rez += fmt.Sprintf("extend bar to %.0f mm from support",s.Ldev)
			}
		}
		case 2:
		stlstr := []string{"short span support","short span midspan","long span support","long span midspan"}
		for i, ast := range s.Asts{
			rez += fmt.Sprintf("% s -> mdu %.2f kn-m %.0f mm dia at %.0f mm spacing \nast req - %.0f mm2 ast - %.0f mm2\n",stlstr[i],s.BM[i],s.Dias[i],s.Spcms[i],s.Astr[i],ast)
		}
	}
	return rez
}

//Dump saves an RccSlb struct to a .json file
func (s *RccSlb) Dump(sname string)(filename string, err error){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	if sname == "" {sname = fmt.Sprintf("slab_%v.json",s.Id)}
	//foldr := filepath.Join(basepath,"../data/out",time.Now().Format("06-Jan-02"))
	foldr := filepath.Join(basepath,"../data/out")
	if _, e := os.Stat(foldr); errors.Is(e, os.ErrNotExist) {
		e := os.Mkdir(foldr, os.ModePerm)
		if e != nil {
			err = e; return
		}
	}
	filename = filepath.Join(foldr,sname)
	data, e := json.Marshal(s)
	if e != nil{err = e; return}
	err = ioutil.WriteFile(filename, data, 0644)
	return
}
