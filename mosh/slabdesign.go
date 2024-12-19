package barf

import (
	"log"
	"math"
)

//SlbDBs is the entry func for bs code slab design
func SlbDBs(s *RccSlb) (err error){
	//bs kode slab design
	switch s.Type{
		case 1:
		switch s.Endc{
			case 0,1:
			err = Slb1DBs(s)
			case 2:
			switch s.Dtyp{
				case 0:
				//coefficients
				err = CSlb1DepthCs(s)
				case 1:
				//envelope analysis
				err = CSlb1Depth(s)
			}
			if err != nil{
				return
			}
			err = CSlb1Stl(s)
		}
		case 2:
		err = Slb2DBs(s)
	}
	if err != nil{
		log.Println("ERRORE,errore--",err)
	}
	return
}


//Slb2DBs is the general (-_-)7 2-way slab design func
//returns depth, astreq, ptreq of a two way slab for dl+ll
//s.Dtyp - 0 - code coeff, 1 - yield line
//s.Code - 1 (is 456), 2 (bs 8110)

func Slb2DBs(s *RccSlb) (err error){
	s.Init()
	s.EndC2W()
	var dused, astreq, ptreq float64
	var wd, wu, mdu, dtrial, kbm, dmin, dserve, dreqy, effd, rumax, mcxsup, mcxmid, mcysup, mcymid float64
	//is code stress block constants
	psfs := 1.15
	k1 := 0.361
	k2 := 0.416
	pgck := 0.025
	var kiter int
	//set defaults
	efcvr := s.Efcvr 
	psfd := 1.5; psfl := 1.5
	if s.Code == 2{psfd = 1.4; psfl = 1.6; pgck = 0.024}
	if efcvr == 0.0 {
		if s.Nomcvr == 0.0{s.Nomcvr = 20.0}
		efcvr = s.Nomcvr + s.Diamain/2.0 
		s.Efcvr = efcvr
	}
	switch s.Code{
		case 1:
		_, _, rumax, _ = BalSecKis(s.Fck, s.Fy, s.DM, psfs, k1, k2)
	}
	switch s.Code{
		case 1:
		//is code
		switch s.Dtyp{
			case 0:
			mcxsup, mcxmid, mcysup, mcymid = Slb2BMCoefIs(s.Endc, s.Lx, s.Ly)
			case 1:
			//mcxsup, mcxmid, mcysup, mcymid = SlbYldRect(s)
		}
		case 2:
		//bs code
		switch s.Dtyp{
			case 0:			
			mcxsup, mcxmid, mcysup, mcymid = Slb2BMCoefBs(s.Endc, s.Ns, s.Nl, s.Lx, s.Ly)
			case 1:
			//mcxsup, mcxmid, mcysup, mcymid = SlbYldRect(s)
		}
	}
	//ALL SPANS INPUT IN MM pls
	lx := s.Lx/1000.0
	dtrial = 75.0
	wd = pgck * dtrial
	wu = psfd * (wd + s.DL) + psfl * s.LL
	wulx2 := wu * math.Pow(lx,2)
	murmax := mcxmid * wulx2
	var mdus []float64
	for _, mc := range []float64{mcxsup, mcxmid, mcysup, mcymid} {
		msec := mc * wulx2
		mdus = append(mdus, msec)
		if msec > murmax {murmax = msec}
	}
	switch s.Code{
		case 1:
		dmin = math.Sqrt(1000.0*murmax/rumax) + efcvr
		dreqy = math.Sqrt(1000.0*mdus[3]/rumax) + efcvr + 10.0
		case 2:
		switch{
			case s.DM <= 0.1:
			kbm = 0.156
			case s.DM > 0.1:
			be := 1.0 - s.DM - 0.4
			kbm = 0.402 * be - 0.18 * math.Pow(be,2.0)
		}
		dmin = math.Sqrt(1000.0*murmax/kbm/s.Fck) + efcvr
		dreqy = math.Sqrt(1000.0*mdus[3]/kbm/s.Fck) + efcvr + s.Diamain
	}
	if dmin < dreqy {dmin = dreqy}
	for dmin > dtrial {
		if kiter > 666{
			err = ErrIter
			return
		}
		dtrial = dmin + 5.0
		wd = pgck * dtrial
		wu = psfd * (wd + s.DL) + psfl * s.LL
		wulx2 = wu * math.Pow(lx,2)
		murmax = mcxmid * wulx2
		for _, mc := range []float64{mcxsup, mcxmid, mcysup, mcymid}{
			msec := mc * wulx2
			mdus = append(mdus, msec)
			if msec > murmax {murmax = msec}
		}
		switch s.Code{
			case 1:
			dmin = math.Sqrt(1000.0*murmax/rumax) + efcvr
			dreqy = math.Sqrt(1000.0*mdus[3]/rumax) + efcvr + s.Diamain
			case 2:
			switch{
				case s.DM <= 0.1:
				kbm = 0.156
				case s.DM > 0.1:
				be := 1.0 - s.DM - 0.4
				kbm = 0.402 * be - 0.18 * math.Pow(be,2.0)
			}
			dmin = math.Sqrt(1000.0*murmax/kbm/s.Fck) + efcvr
			dreqy = math.Sqrt(1000.0*mdus[3]/kbm/s.Fck) + efcvr + s.Diamain
		}
		if dmin < dreqy {dmin = dreqy}
		kiter++
	}
	dused = 5.0 * math.Ceil(dmin/5.0)
	effd = dused - efcvr
	wu = psfd * (pgck*dused + s.DL) + psfl * s.LL
	mdu = wu * mcxmid * math.Pow(lx, 2)
	astreq = BalSecAst(murmax,effd,s.Fck,s.Fy,s.Code)
	ptreq = 100.0 * astreq/(1000.0 * effd)
	switch s.Code{
		case 1:
		dserve = RccSlabServeChk(s.Type, s.Endc, s.Fy, s.LL, s.Lx, ptreq, efcvr)
		case 2:
		_, _, dserve = SlbSdratBs(s, mdu, astreq, astreq, dused)
		dserve += efcvr
	}
	kiter = 0
	for dserve > dused{
		if kiter > 666{
			err = ErrIter
			return
		}
		dused = dused + 5.0
		wu = psfd * (0.025*dused + s.DL) + psfl * s.LL
		wulx2 = wu * math.Pow(lx,2)
		mdu = mcxmid * wulx2
		effd = dused - efcvr
		astreq = BalSecAst(mdu, effd, s.Fck, s.Fy, s.Code)
		ptreq = 100.0 * astreq/(1000.0 * effd)
		switch s.Code{
			case 1:
			dserve = RccSlabServeChk(s.Type, s.Endc, s.Fy, s.LL, s.Lx, ptreq, efcvr)
			case 2:
			_, _, dserve = SlbSdratBs(s, mdu, astreq, astreq, dused)
			dserve += efcvr
		}
		kiter++
	}
	//log.Println("ASTREQ ", astreq, ptreq)
	//log.Printf("min depth for flexure %0.2f\n", dmin)
	//log.Printf("min depth for deflection %0.2f\n", dserve)
	//log.Printf("min depth assumed %0.2f against dserve %.2f dmin %.2f\n", dused, dserve, dmin)
	//log.Printf("effective depth %f mm",effd)
	//get valiant steel
	s.Dused = dused
	s.BM = make([]float64,4)
	for idx, mc := range []float64{mcxsup, mcxmid, mcysup, mcymid}{
		s.BM[idx] = mc * wulx2
	}
	
	err = s.Slb2WRbr()
	if err != nil{log.Println(err)}
	//s.Table(s.Verbose)
	//if s.Verbose{log.Println(s.Rprt)}
	return 
}

//Slb2WRbr calcs the steel (dia/spacing) required for a 2 way slab
//given s.BM - (x) mcx support, mcx mid, (y) mcy support, mcy mid
func (s *RccSlb) Slb2WRbr() (err error){
	//given mdus = mcxsup, mcxmid, mcysup, mcymid
	var asts ,diams, spcs []float64
	dused := s.Dused; effd := s.Dused - s.Efcvr
	for idx, msec := range s.BM{
		//log.Println("sec->",idx+1, msec,"kn-m/m")
		dsec := dused; effdi := effd
		if idx > 1{dsec = dused - s.Diamain; effdi = effd - s.Diamain/2.0 - s.Diadist/2.0}
		astreq := BalSecAst(msec, effdi, s.Fck, s.Fy, s.Code)
		rezmap, mindia := SlabRbrDiaSpcing(dsec, astreq, s.Efcvr)
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
	}
	
	if s.BM[0] == 0.0{
		//add 50% of midspan x steel
		asts[0] = math.Round(asts[1]/2.0)
		diams[0] = diams[1]
		spcs[0] = math.Round(spcs[1]*2.0)
		//log.Println("short zero->",spcs[0],spcs[1])
	}
	if s.BM[2] == 0.0{
		//add 50% of midspan y steel
		asts[2] = math.Round(asts[3]/2.0)
		diams[2] = diams[3]
		spcs[2] = math.Round(spcs[3]*2.0)
		
		//log.Println("long zero->",spcs[2],spcs[3])
	}
	s.Asts = asts
	s.Dias = diams
	s.Spcms = spcs
	s.Dz = true
	return
}


//Slb1DBs returns the depth of a 1 way slab (bs code)
//pretty much a rewrite of SlbDIs
func Slb1DBs(s *RccSlb) (err error){
	var kbm, dmcoef, sfcoef, psfd, psfl, ptreq, astreq, ptsup, wd, wu, mdu, vumax, lspan, dtrial, dmin, dused, dserve, effd float64
	s.Init()
	efcvr := s.Efcvr
	lspan = s.Lspan/1000.0
	switch s.Endc{
		case 0:
		dmcoef = 2.0
		sfcoef = 1.0
		case 1:
		dmcoef = 8.0
		sfcoef = 2.0
	}
	wd = 0.025 * dtrial
	switch s.Code{
		case 1:
		psfd = 1.5; psfl = 1.5
		case 2:
		psfd = 1.4; psfl = 1.6
	}
	wu = psfd * (wd + s.DL ) + psfl * s.LL
	mdu = wu * math.Pow(lspan, 2) / dmcoef
	switch{
		case s.DM <= 0.1:
		kbm = 0.156
		case s.DM > 0.1:
		be := 1.0 - s.DM - 0.4
		kbm = 0.402 * be - 0.18 * math.Pow(be,2.0)
	}
	dmin = math.Sqrt(1000.0*mdu/kbm/s.Fck) + efcvr
	//log.Println("init depth->",dmin, kbm)
	var kiter int
	for dmin > dtrial{
		if kiter > 666{return ErrIter}		
		dtrial += 5.0
		//dtrial = dmin + 10.0
		wd = 0.025 * dtrial
		wu = psfd * (wd + s.DL) + psfl * s.LL
		mdu = wu * math.Pow(lspan, 2) / dmcoef
		dmin = math.Sqrt(1000.0*mdu/kbm/s.Fck) + efcvr
		kiter++
	}
	dused = 5.0 * math.Ceil(dmin/5.0)
	//log.Println("junk",dserve, dused, dmin, astreq, ptreq, sfcoef, effd, vumax)

	wu = 1.5 * (0.025*dused + s.DL + s.LL)
	mdu = wu * math.Pow(lspan, 2) / dmcoef
	effd = dused - efcvr
	astreq = BalSecAst(mdu, effd, s.Fck, s.Fy, s.Code)
	ptreq = 100.0 * astreq/(1000.0 * effd)
	switch s.Code{
		case 1:
		dserve = RccSlabServeChk(s.Type, s.Endc, s.Fy, s.LL, s.Lspan, ptreq, efcvr)
		case 2:
		_, _, dserve = SlbSdratBs(s, mdu/1e6, astreq, astreq, dused)
		//_, _, dserve = SpanDepRatBs(s.Fck, s.Fy, s.Lspan, 0.0, 1000.0, 1000.0, dused, effd, astreq, astreq, 0.0, mdu, s.DM, 1, s.Endc)
		//log.Println("dserve")
		//dserve += efcvr
		//log.Println("kiter, dserve, effd->",kiter, dserve, effd)
	}
	kiter = 0
	iter := 0
	for iter != -1 {
		if kiter > 666{return ErrIter}
		if dserve <= dused && SecShear(s.Fck, vumax, 1000.0, effd, dused, ptsup, s.Code){ 
			//log.Println("ptsup->", ptsup)
			//log.Println( SecShear(s.Fck, vumax, 1000.0, effd, dused, ptsup, s.Code))
			iter = -1
			break
		}
		dused += 5.0
		wu = psfd * (0.025*dused + s.DL) + psfl * s.LL
		mdu = wu * math.Pow(lspan, 2) / dmcoef
		effd = dused - efcvr
		astreq = BalSecAst(mdu, effd, s.Fck, s.Fy, s.Code)
		ptreq = 100.0 * astreq/(1000.0 * effd)
		switch s.Code{
			case 1:
			dserve = RccSlabServeChk(s.Type, s.Endc, s.Fy, s.LL, s.Lspan, ptreq, efcvr)
			case 2:
			_, _, dserve = SlbSdratBs(s, mdu/1e6, astreq, astreq, dused)
		}
		//check for shear
		vumax = wu * lspan / sfcoef
		//check for amount of steel at support (for one way slab)
		if s.Ibent == 0.0 {
			//ibent = percentage of bent up bars
			switch s.Endc{
				case 0:
				ptsup = 100.0 * astreq/(1000.0 * effd)
				case 1:
				ptsup = 50.0 * astreq/(1000.0 * effd)
			}
		} else {
			ptsup = s.Ibent * astreq/(1000.0 * effd)
		}
		kiter++
	}
	effd = dused - efcvr
	//log.Println("final depths-",dserve, dused)
	//find areas of steel
	//log.Println("ast req-",astreq,"mm2","mdu->",mdu)
	murmax := kbm * math.Pow(effd, 2) * 1e-3
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
	if math.Abs(prvcvr - efcvr) > 5.0{return ErrCvr}
	//distribution isteel
	var astd float64
	if s.Diadist == 6.0 {
		s.Fyd = 250.0
		astd = 1.5 * dused
	} else {
		astd = 1.2 * dused
	}
	sdmax := 5 * effd
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
	//development length check (usually at support, right?)
	var m1 float64
	switch s.Endc{
		case 0:
		m1 = mdu
		case 1:
		m1 = 0.5 * mdu
	}
	devchk, _, ldreq := SlabDevLen(s.Fck, s.Fy, m1, vumax, effd, s.Diamain, s.Bsup, s.Endc)
	s.Devchk = devchk
	s.Ldev = ldreq
	s.Murmax = murmax
	s.Astm = asprov; s.Spcm = spcmain; s.Astd = asdprov; s.Spcd = sds; s.Astreq = astreq
	s.Dused = dused; s.BM = append(s.BM, mdu)
	s.Dz = true
	if s.Lx == 0.0{
		s.Lx = s.Lspan
	}
	if s.Ly == 0.0{
		s.Ly = s.Lx * 2.0
	}
	return
}

//BalSecAst returns the area of steel required for a balanced section to resist a moment mdu
func BalSecAst(mdu, effd, fck, fy float64, code int) (astr float64){
	//get area of steel required for a balanced section
	if code > 1{
		m1 := mdu*1e3/math.Pow(effd,2)/fck
		z := effd * (0.5 + math.Sqrt(0.25 - m1/1.134))
		astr = mdu * 1e6/fy/z/0.87
		
	} else {
		m1 := 4.6 * mdu * 1e3/fck/math.Pow(effd,2)
		ptreq := 0.5 * fck/fy * (1.0 - math.Sqrt(1.0 - m1)) * 100.0
		astr = ptreq * 1000.0 * effd/100.0
	}
	return
}


//Slb2BMCoefBs calcs moment coefficients at short span - support, midspan; long span - support, midspan
//ns, nl - no. of discontinuous short, long edges
func Slb2BMCoefBs(endc, ns, nl int, lx, ly float64) (mcxsup, mcxmid, mcysup, mcymid float64) {
	switch endc{
		case 10:
		//simply supported 2 way slab
		x := math.Pow(ly/lx,4); y := math.Pow(ly/lx,2)
		mcxmid = x/8.0/(1.0+x); mcymid = y/8.0/(1.0+x)
		default:
		var k1, k2, k3, k4, by, bx float64
		k1, k2, k3, k4  = 4.0/3.0, 4.0/3.0, 4.0/3.0, 4.0/3.0
		nd := float64(ns + nl)
		switch ns{
			case 1:
			k2 = 0
			case 2:
			k1 = 0; k2 = 0
		}
		switch nl{
			case 1:
			k4 = 0
			case 2:
			k3 = 0; k4 = 0
		}
		by = (24.0 + 2.0 * nd + 1.5 * math.Pow(nd, 2))/1000.0
		y := 2.0 * (3.0 - math.Sqrt(18.0) * lx/ly * (math.Sqrt(by + k1 * by)+math.Sqrt(by + k2 * by)))/9.0
		bx = y/math.Pow(math.Sqrt(1.0+k3)+math.Sqrt(1.0+k4),2)
		mcxsup = k3 * bx; mcxmid = bx
		mcysup = k1 * by; mcymid = by
	}
	return
}

//Slb1RbrDet is a one way slab rebar detailing func
//is possibly not used, DELETE
func Slb1RbrDet(s *RccSlb, astz []float64, dused float64, bscode int) (error){
	//IZ DIS EVEN USE
	//one - way slab rebar detail func
	effd := dused - s.Efcvr
	for _, ast := range astz{
		rezmap, mindia := SlabRbrDiaSpcing(dused, ast, s.Efcvr)
		//CHECK FOR len(rezmap) = 0
		if len(rezmap) == 0 {
			return ErrSpacing
		}
		var asprov, diacom, spcmain float64
		asprov = rezmap[mindia][2]
		spcmain = rezmap[mindia][1]
		diacom = mindia

		if val, ok := rezmap[s.Diamain]; ok {
			asprov = val[2]
			spcmain = val[1]
			diacom = s.Diamain
		}
		//CHECK FOR COVER
		//distribution isteel
		var astd float64
		if s.Diadist == 6.0 {
			s.Fyd = 250.0
			astd = 1.5 * dused
		} else {
			astd = 1.2 * dused
		}
		sdmax := 5 * effd
		if sdmax > 300.0 {
			sdmax = 300.0
		}
		sds := 1000 * RbrArea(s.Diadist) / astd
		if sds > sdmax {
			sds = sdmax
		}
		sds = 10.0 * math.Floor(sds/10.0)
		astd = math.Round(1000 * RbrArea(s.Diadist) / sds)
		s.Asts = append(s.Asts, asprov)
		s.Asds = append(s.Asds, astd)
		s.Dias = append(s.Dias, diacom)
		s.Dists = append(s.Dists, s.Diadist)
		s.Spcms = append(s.Spcms, spcmain)
		s.Spcds = append(s.Spcds, sds)
	}
	return nil
}

//SecShear checks if a section is safe in shear 
func SecShear(fck, v, b, effd, dused, pt float64, bscode int) bool{
	/*
	   checks if section is safe in shear 
	*/
	var vuc float64
	//HERE PT IS ALREADY PERCENTAGE STEEL
	//pt := 100.0 * ast/b/effd
	vud := v * 1e3/b/effd
	if bscode > 1{
		//log.Println("bs KODE")
		m1 := 400.0/effd
		if effd > 400.0{m1 = 1.0}
		if pt > 3.0{pt = 3.0}
		if pt < 0.15{pt = 0.15}
		vuc = 0.79 * math.Pow(pt, 1./3.0) * math.Pow(m1, 0.25)/1.25
	} else {
		//log.Println("is KODE")
		//log.Println("b, dused, effd, ast",b, dused, effd, ast,"ptreq->",pt)
		beta := 0.8 * fck/6.89/pt
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
		vuc = k5 * 0.8499999 * math.Sqrt(0.8*fck) * (math.Sqrt(1.0+5.0*beta) - 1.0)/ (6.0 * beta)
	}
	//log.Println("vuc->",vuc,"v->",vud, "shear check->", vuc > vud)
	return vuc > vud
}
