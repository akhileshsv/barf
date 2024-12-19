package barf

import (
	"log"
	"math"
)

//Slb2IsDused returns dused, astreq for a 2 way slab using is code coefficients
//YA DELETE THIS
func Slb2IsDused(s *RccSlb) (dused, astreq, ptreq float64,err error){
	//s is (usually) generated by slabgen in gen
	//returns dused, astreq, ptreq, error
	var wd, wu, mdu, dtrial, dmin, dserve, effd float64
	psfs := 1.15
	k1 := 0.361
	k2 := 0.416
	var kiter int
	//set defaults
	if s.Diamain == 0.0 {
		s.Diamain = 8.0
	}
	if s.Fy == 0.0 {
		s.Fy = 550.0
	}
	if s.Diadist == 0.0 {
		s.Diadist = 8.0
	}
	if s.Fyd == 0.0 {
		s.Diadist = 550.0
	}
	if s.Nomcvr == 0.0{
		s.Nomcvr = 20.0
	}
	efcvr := s.Efcvr //s.Nomcvr + s.Diamain/2.0
	
	if efcvr == 0.0 {
		efcvr = s.Nomcvr + s.Diamain/2.0 + s.Diadist/2.0
		s.Efcvr = efcvr
	}
	_, _, rumax, _ := BalSecKis(s.Fck, s.Fy, s.DM, psfs, k1, k2)
	
	//log.Println("TWO WAY SLAB")
	//two way slab
	mcxsup, mcxmid, mcysup, mcymid := Slb2BMCoefIs(s.Endc, s.Lx, s.Ly)
	lx := s.Lx/1000.0
	dtrial = 100.0
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
	if s.Endc == 10 {dreqy = math.Sqrt(1000.0*mdus[1]/rumax) + efcvr + 12.0}
	if dmin < dreqy {dmin = dreqy}
	for dmin > dtrial {
		if kiter > 3000{
			err = ErrIter
			return
		}
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
		dreqy = math.Sqrt(1000.0*mdus[3]/rumax) + efcvr + 10.0
		if s.Endc == 10 { dreqy = math.Sqrt(1000.0*mdus[1]/rumax) + efcvr + 12.0 }
		if dmin < dreqy {dmin = dreqy}
		kiter++
	}
	dused = 10.0 * math.Ceil(dmin/10.0)
	effd = dused - efcvr
	wu = 1.5 * (0.025*dused + s.DL + s.LL)
	mdu = wu * mcxmid * math.Pow(lx, 2) 
	H := 1.0 - (4600.0*mdu)/(s.Fck*math.Pow(effd, 2))
	ptreq = (0.5 * s.Fck / s.Fy) * (1.0 - math.Sqrt(H)) * 100.0
	//fmt.Println("PTREQ",ptreq)
	//astreq = ptreq * 10.0 * effd
	dserve = RccSlabServeChk(s.Type, s.Endc, s.Fy, s.LL, s.Lx, ptreq, efcvr)
	kiter = 0
	for dserve > dused {
		if kiter > 3000{
			err = ErrIter
			return
		}
		dused = dused + 5.0
		wu = 1.5 * (0.025*dused + s.DL + s.LL)
		//mdu = wu * mcxmid * math.Pow(lx, 2) 
		wulx2 = wu * math.Pow(lx,2)
		mdu = mcxmid * wulx2
		//deflection check applied only at midspan? not at support
		effd = dused - efcvr
		H = 1.0 - ((4600.0 * mdu) / (s.Fck * math.Pow(effd, 2)))
		ptreq = (0.5 * s.Fck / s.Fy) * (1.0 - math.Sqrt(H)) * 100.0
		//astreq = ptreq * 10.0 * effd
		dserve = RccSlabServeChk(s.Type, s.Endc, s.Fy, s.LL, s.Lx, ptreq, efcvr)
		kiter++
	}
	log.Println("ASTREQ ", astreq, ptreq)
	log.Printf("min depth for flexure %0.2f\n", dmin)
	log.Printf("min depth for deflection %0.2f\n", dserve)
	log.Printf("min depth assumed %0.2f against dserve %.2f dmin %.2f\n", dused, dserve, dmin)
	err = nil
	return 
	//effd = dused - efcvr
}

func Slb2IsDet(s *RccSlb, dused float64){
	//dused is in MM FOR GOD'S SAKE MM
	var asts , diams, spcs []float64
	mcxsup, mcxmid, mcysup, mcymid := Slb2BMCoefIs(s.Endc, s.Lx, s.Ly)
	wu := 1.5 * (0.025*dused + s.DL + s.LL)
	//mdu = wu * mcxmid * math.Pow(lx, 2) 
	wulx2 := wu * math.Pow(s.Lx,2)
	mdu := mcxmid * wulx2
	efcvr := s.Efcvr //s.Nomcvr + s.Diamain/2.0
	if efcvr == 0.0 {
		efcvr = s.Nomcvr + s.Diamain/2.0
	}
	effd := dused - efcvr
	H := 1.0 - ((4600.0 * mdu) / (s.Fck * math.Pow(effd, 2)))
	//ptreq := (0.5 * s.Fck / s.Fy) * (1.0 - math.Sqrt(H)) * 100.0
	mdus := make([]float64,4)
	for idx, mc := range []float64{mcxsup, mcxmid, mcysup, mcymid} {
		//CHANGE LONG SPAN STEEL EFFD
		if idx == 3{effd = effd - s.Diamain/2.0}
		msec := mc * wulx2
		mdus[idx] = msec
		H = 1.0 - ((4600.0 * msec) / (s.Fck * math.Pow(effd, 2)))
		astreq := (0.5 * s.Fck / s.Fy) * (1.0 - math.Sqrt(H)) * 1000.0 * effd
		rezmap, mindia := SlabRbrDiaSpcing(dused, astreq, efcvr)
		if len(rezmap) == 0 {
			log.Println("ERRORE,errore-> rebar spacing error")
			//return mosh.ErrSpacing
		}
		var asprov, diacom, prvcvr, spcmain float64
		asprov = rezmap[mindia][2]
		spcmain = rezmap[mindia][1]
		diacom = mindia
		log.Println("min area dia->",diacom, " at spacing",spcmain, "as provided",asprov)
		if val, ok := rezmap[s.Diamain]; ok {
			asprov = val[2]
			spcmain = val[1]
			diacom = s.Diamain
		}
		prvcvr = s.Nomcvr + diacom/2.0
		if prvcvr != s.Efcvr {
			if math.Abs(prvcvr - s.Efcvr) <= 2.0 {
				continue
			} else {
				log.Println("ERRORE,errore->effective cover error")
				//return ErrCvr
			}
		}
		diams = append(diams, diacom)
		spcs = append(spcs, spcmain)
		asts = append(asts, asprov)
		//ptprov := 100 * asprov / (1000 * effd)
	}
	idxdict := map[int]string{
		0:"short span support",
		1:"short span midspan",
		2:"long span support",
		3:"long span midspan",
	}
	for idx, ast := range asts {
		log.Println(idxdict[idx], "mdu ", mdus[idx], "dia ",diams[idx], " mm at ",spcs[idx], "mm spacing, astm",ast)
	}
	return 
}
