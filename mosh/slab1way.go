package barf

import (
	"math"
)
//DELTE THIS DAMMIT
func (s *RccSlb) Draw1W(term string){
	//s.DrawPlan()
	//s.Drawsectionx()
	//s.Drawsectiony()
}

//SlbSsIsDused returns dused, astreq, ptreq for a simply supported/cantilever one-way slab
//DELETE THIS - this is not needed at all
//thousands of slab funcs doing the same thing with different names dammit
func SlbSsIsDused(s *RccSlb) (dused, astreq, ptreq float64, err error){
	var wd, wu, mdu, dtrial, dmin, dserve, effd float64
	psfs := 1.15
	k1 := 0.361
	k2 := 0.416
	var kiter int
	//set defaults
	if s.Nomcvr == 0.0 {
		s.Nomcvr = 20.0
	}
	if s.Diamain == 0.0 {
		s.Diamain = 8.0
	}
	if s.Fy == 0.0 {
		s.Fy = 415.0
	}
	if s.Diadist == 0.0 {
		s.Diadist = 6.0
		s.Fyd = 250.0
	}
	efcvr := s.Efcvr //s.Nomcvr + s.Diamain/2.0
	if efcvr == 0.0 {
		efcvr = s.Nomcvr + s.Diamain/2.0
	}
	_, _, rumax, _ := BalSecKis(s.Fck, s.Fy, s.DM, psfs, k1, k2)
			//simply supported, cantilever
	dtrial = 100.0
	lspan := s.Lspan / 1000.0
	var dmcoef float64
	if s.Endc == 1 {
		dmcoef = 8.0
		//sfcoef = 2.0
	} else {
		dmcoef = 2.0
		//sfcoef = 1.0
	}
	wd = 0.025 * dtrial
	wu = 1.5 * (wd + s.DL + s.LL)
	mdu = wu * math.Pow(lspan, 2) / dmcoef
	dmin = math.Sqrt(1000.0*mdu/rumax) + efcvr
	//vumax = wu * lspan/sfcoef
	for dmin > dtrial {
		if kiter > 3000{
			err = ErrIter
			return
		}
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
		if kiter > 3000{
			err = ErrIter
			return
		}
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
	return 
}

