package barf

//continuous slab design routines

import (
	"fmt"
	"log"
	"math"
	//"sort"
	"errors"
	kass"barf/kass"
)

//RSlb1Dims initializes dimensions for a ribbed slab
//NOT DONE/it is not werk
func RSlb1Dims(s *RccSlb){
	//init dims for a ribbed slab
	if s.Bw == 0.0{
		s.Bw = 65.0
	}
	if s.Df == 0.0{
		s.Df = 50.0
	}
	if s.Bf == 0.0{
		mrel := 0
		if s.Nspans == 1{mrel = 1}
		b1 := kass.GetBfRcc(s.Code, mrel, 1.0, s.Lspan, s.Df, s.Bw) 
		//minimum of 4 ribs per span? subramanian
		//supports are solid at least until 0.15l (span - 0.15l, 0.7l, 0.15l)
		b2 := math.Round(0.7 * s.Lspan/4.0/10.0) * 10.0
		log.Println(b1, b2)
		s.Bf = b2
		if b1 < b2{
			log.Println("ERRORE, errore->")
		}
	}
}

//RSlb2Chk performs design checks for a waffle slab using 2 way slab design coefficients
func RSlb2Chk(s *RccSlb) (err error, ok bool){
	//waffle slab design
	var dl, ll float64
	s.Dused = s.Dw
	s.Init()
	s.EndC2W()
	//check dimensions - rib width, rib depth, bflange > bf rib
	switch s.Ldcalc{
		case 0:
		//calc dead load
		//net wt of flange - area (bf2) x df
		dl = s.Bf * s.Bf * s.Df * 25.0
		//wt of rib - perimeter (4.0 * bf) * width (bw/2.0) * depth 
		dl += 2.0 * s.Bf * s.Bw * (s.Dw - s.Df)  * 25.0
		//divide by net area - bf * bf for udl in kn/m2
		dl = dl/s.Bf/s.Bf/1e3
	}
	psfs := []float64{1.5,1.0,1.5,0.0}
	switch s.Code{
		case 2:
		psfs = []float64{1.4,1.0,1.6,0.0}
		case 3:
		psfs = []float64{1.35,1.0,1.5,0.0}
		s.Code = 2
	}
	fmt.Println("self wt->",dl,"kn/m2")
	dl += s.DL; ll += s.LL
	wud := psfs[0]*dl + psfs[2]*ll
	fmt.Println("ult. load wud->",wud,"kn/m2")
	var mcxsup, mcxmid, mcysup, mcymid, vux, vuy float64
	switch s.Code{
		case 1:
		mcxsup, mcxmid, mcysup, mcymid = Slb2BMCoefIs(s.Endc, s.Lx, s.Ly)
		case 2:
		mcxsup, mcxmid, mcysup, mcymid = Slb2BMCoefBs(s.Endc, s.Ns, s.Nl, s.Lx, s.Ly)
	}
	//moments are per m width so multiply by bf for each rib vals (see mosley ec2)
	wlx2 := wud * math.Pow(s.Lx/1e3,2) * s.Bf/1e3
	fmt.Println("mcxsup, mcxmid, mcysup, mcymid coeff->",mcxsup, mcxmid, mcysup, mcymid)
	mcxsup, mcxmid, mcysup, mcymid = mcxsup * wlx2, mcxmid * wlx2, mcysup * wlx2, mcymid * wlx2
	fmt.Println("mcxsup, mcxmid, mcysup, mcymid->",mcxsup, mcxmid, mcysup, mcymid,"kn-m/m")
	//var vxsup, vxmid, vysup, vymid float64
	avxc, avxd, avyc, avyd := Slb2VfCoeff(s.Endc, s.Lx, s.Ly)
	fmt.Println("avxc, avxd, avyc, avyd coeff->",avxc, avxd, avyc, avyd)
	vux = avxc * wud * s.Lx * s.Bf/1e6
	if avxd > avxc{
		vux = avxd * wud * s.Lx * s.Bf/1e6
	}
	vuy = avyc * wud * s.Lx * s.Bf/1e6
	if avyd > avyc{
		vuy = avyd * wud * s.Lx * s.Bf/1e6
	}
	fmt.Println("vux, vuy->",vux, vuy)
	barr := make([]*RccBm, 4)
	dsup := []float64{s.Bw, s.Dused}
	dmid := []float64{s.Bf,s.Dw,s.Bw,s.Df}
	endc := 2
	if s.Endc == 10{
		endc = 1
	}
	for i, mu:= range []float64{mcxsup, mcxmid, mcysup, mcymid}{
		var lspan float64
		switch i{
			case 0:
			lspan = s.Lx * (s.S1 + s.S2)/1e3
			case 1:
			lspan = s.Lx/1e3
			case 2:
			lspan = s.Ly * (s.S3 + s.S4)/1e3
			case 3:
			lspan = s.Ly/1e3
		}
		barr[i] = &RccBm{
			Id:i+1,
			Mid:s.Id,
			Fck:s.Fck,
			Fy:s.Fy,
			Bf:s.Bf,
			Df:s.Df,
			Bw:s.Bw,
			Dused:s.Dused,
			Styp:1,
			Cvrt:s.Efcvr,
			Cvrc:s.Efcvr,
			Flip:true,
			Code:s.Code,
			Mu:mu,
			Dims:dsup,
			Verbose:s.Verbose,
			Lspan:lspan,
			Rslb:true,
			Endc:endc,
			D1: s.Diamain,
			//Lsx:bm.Lsx*1e3,
			//Rsx:bm.Rsx*1e3,
			//D1:d1,D2:d2,
			//Ldx:bm.Ldx, Rdx:bm.Rdx,
			//Dslb:dslb,
			//Ismid:ismid,
			//Term:bm.Term,
		}
		barr[i].Init()
		switch i{
			case 1, 3:
			barr[i].Tyb = 1.0
			barr[i].Styp = 6
			barr[i].Flip = false
			barr[i].Dims = dmid
			barr[i].Shrdz = true
			if i == 1{
				barr[i].Vu = vux
			} else{
				barr[i].Vu = vuy				
			}
		}
		err := BmDesign(barr[i])
		if err != nil{
			fmt.Println(err)
		}
		//barr[i].Printz()
		//barr[i].Draw()
		//PlotBmGeom(barr[i], "dumb")
		//fmt.Println(barr[i].Txtplot[0])
	}
	s.Dz = true
	
	s.R2Quant(barr)
	s.Dz = true
	s.R2Table(barr, s.Verbose)
	return
}

//RSlb1Chk performs design checks for a one way ribbed slab using CBm envelopes
func RSlb1Chk(s *RccSlb) (err error, ok bool){
	//checks a ribbed slab
	//first check dims
	//LATER
	//calc dead load
	var dl, ll float64
	s.Dused = s.Dw
	s.Init()
	//check dimensions - rib width, rib depth, bflange > bf rib
	switch s.Ldcalc{
		case 0:
		//calc dead load
		//wt of flange
		dl = s.Bf * s.Df * 25.0
		//wt of rib
		dl += (s.Dw - s.Df) * s.Bw * 25.0
		//divide by bf
		dl = dl/s.Bf/1e3
	}
	//log.Println("slab self weight->",dl, "kn/m2")
	s.Swt = dl
	dl += s.DL; ll = s.LL
	//log.Println("dl,ll->",dl, ll,"kn/m2")
	var lspans []float64
	if s.Nspans > 0 && len(s.Lspans) == 0{
		for i := 0; i < s.Nspans; i++{
			s.Lspans = append(s.Lspans, s.Lspan)
			lspans = append(lspans, s.Lspan/1000.0)
		}
	} else {
		for i := range s.Lspans{
			lspans = append(lspans, s.Lspans[i]/1000.0)
		}
	}
	if s.Clvrs == nil{s.Clvrs = [][]float64{{0,0,0},{0,0,0}}}
	if s.Bsup > 0.0{
		s.Bsups = make([]float64, s.Nspans + 1)
		for i := range s.Bsups{
			s.Bsups[i] = s.Bsup
		}
	}
	psfs := []float64{1.5,1.0,1.5,0.0}
	switch s.Code{
		case 2:
		psfs = []float64{1.4,1.0,1.6,0.0}
		case 3:
		//just to check examples NO WAY IS EC2 HAPPENING
		psfs = []float64{1.35,1.0,1.5,0.0}
		s.Code = 2
	}
	var cb *CBm
	//get dl and ll
	dl = dl * s.Bf/1e3 
	ll = ll * s.Bf/1e3 
	cb = &CBm{
		Fck:s.Fck,
		Fy:s.Fy,
		DL:dl,
		LL:ll,
		Nspans:len(s.Lspans),
		Lspans:lspans,
		Clvrs:s.Clvrs,
		Sectypes:[]int{6},
		Sections:[][]float64{{s.Bf,s.Dw,s.Bw,s.Df}},
		PSFs:psfs,
		DM:s.DM,
		Lsxs:s.Bsups,
		Rslb:true,
		Code:s.Code,
		Nomcvr:s.Nomcvr,
		Efcvr:s.Efcvr,
	}
	var bmenv map[int]*kass.BmEnv
	bmenv, err = CBeamEnvRcc(cb, cb.Term, true)
	if err != nil{
		return
	}
	_, err = CBmDz(cb,bmenv)
	if err != nil{
		return
	}
	s.RQuant(cb)
	s.Dz = true
	s.RTable(cb,true)
	return
}

//CSlb1Depth gets the depth of a 1-way continuous slab (dused) using CBm envelope calcs
func CSlb1Depth(s *RccSlb) (err error){
	/*
	   get depth using pattern loading
	*/
	var lspans []float64
	if s.Nspans > 0 && len(s.Lspans) == 0{
		for i := 0; i < s.Nspans; i++{
			s.Lspans = append(s.Lspans, s.Lspan)
			lspans = append(lspans, s.Lspan/1000.0)
		}
	} else {
		lspans = make([]float64, len(s.Lspans))
		for i, v := range s.Lspans{
			lspans[i] = v/1000.0
		}
	}
	
	if s.Nspans == 0{s.Nspans = len(lspans)}
	if len(lspans) == 0{err = fmt.Errorf("invalid number of spans->%v",lspans); return}
	if s.Clvrs == nil{s.Clvrs = [][]float64{{0,0,0},{0,0,0}}}
	if s.Bsup > 0.0{
		s.Bsups = make([]float64, s.Nspans + 1)
		for i := range s.Bsups{
			s.Bsups[i] = s.Bsup
		}
	}
	var wdl, mdu, msup, vud, k2, ldrat, diamain, dused float64
	var iter, dchk, bmchk, schk int
	psfs := []float64{1.5,1.0,1.5,0.0}
	switch s.Code{
		case 2:
		psfs = []float64{1.4,1.0,1.6,0.0}
		case 3:
		psfs = []float64{1.35,1.0,1.5,0.0}
	}
	switch fy{
		case 250.0:
		k2 = 0.149
		case 415.0:
		k2 = 0.138
		default:
		k2 = 0.133
	}
	diamain = 10.0
	efcvr := s.Nomcvr + diamain/2.0
	fs := 5.0 * fy/8.0
	ldrat = 26.0
	dused = 75.0
	dused -= 5.0
	//log.Println("starting iter->")
	var cb *CBm
	var bmenv map[int]*kass.BmEnv
	for iter != -1{
		iter++
		dused += 5.0
		dchk, bmchk, schk = 0, 0 ,0
		effd := dused - efcvr
		if iter > 666{
			log.Println("ERRORE, errore-> iteration error")
			err = ErrIter
			return
		}
		wdl = s.DL + 0.025 * dused
		cb = &CBm{
			Fck:s.Fck,
			DL:wdl,
			LL:s.LL,
			Nspans:len(s.Lspans),
			Lspans:lspans,
			Clvrs:s.Clvrs,
			Sections:[][]float64{{1000.0,dused}},
			PSFs:psfs,
			DM:s.DM,
			Lsxs:s.Bsups,
		}
		bmenv, err = CBeamEnvRcc(cb, "", false)
		if err != nil{
			log.Println("ERRORE,errore->",err)
			return
		}
		for i := 1; i <= len(s.Lspans); i++{
			bm := bmenv[i]; lspan := s.Lspans[i-1]
			mdu = math.Abs(bm.Mpmax)
			if s.DM > 0.0{
				msup = math.Abs(bm.Mlrd)
				if math.Abs(bm.Mrrd) > msup{msup = bm.Mrrd}
				vud = math.Abs(bm.Vlrd)
				if math.Abs(bm.Vrrd) > vud {vud = math.Abs(bm.Vrrd)}
			} else {
				msup = math.Abs(bm.Ml)
				if math.Abs(bm.Mr) > msup{msup = bm.Mr}
				vud = math.Abs(bm.Vl)
				if math.Abs(bm.Vr) > vud {vud = math.Abs(bm.Vr)}
			}
			m1 := mdu*1e6/1000.0/math.Pow(effd,2)
			ldfac := 0.55 + (477.0 - fs)/120.0/m1			
			if ldfac > 2.0 {ldfac = 2.0}
			dreq := (lspan/ldfac/ldrat)
			astr := BalSecAst(mdu, effd, s.Fck, fy, s.Code)
			ptreq := astr*100.0/1000.0/effd
			d2 := RccSlabServeChk(1, 2, s.Fy, s.LL, lspan, ptreq, efcvr)
			if d2 - efcvr > dreq{
				dreq = d2 - efcvr
			}
			if dreq/effd - 1.0 >= 0.0{
				dchk = -1
				break
			}
			mmax := mdu
			//if msup > mdu{mmax = msup}
			dbm := math.Sqrt(mmax * 1e3/k2/s.Fck)
			if dbm >= effd{
				bmchk = -1
				break
			}
			astsup := BalSecAst(msup, effd, s.Fck, fy, s.Code)
			shrchk := SecShear(s.Fck, vud/effd, 1000.0, effd, dused, astsup, s.Code)
			if !shrchk{
				schk = -1
				break
			}
		}
		if dchk + bmchk + schk == 0 {
			//log.Println("effd, dused",effd,dused)
			iter = -1
			s.BMs = [][]float64{}
			var b1, b2, b3 float64
			for i := 0; i < s.Nspans; i++{
				bm := bmenv[i+1]
				b1 = bm.Ml; b3 = bm.Mr; b2 = bm.Mpmax
				if s.DM > 0.0 {
					b1 = bm.Mlrd; b3 = bm.Mrrd; b2 = bm.Mprmax
				}
				s.BMs = append(s.BMs, []float64{b1, b2, b3})
			}
			break
		}
	}
	s.Dused = dused
	//rstr := fmt.Sprintf("%.1f\n",s.BMs)
	//fmt.Println("lspanz",s.Lspans)
	//fmt.Println(rstr)
	return
}

//CSlb1DepthCs updates the depth of a continuous one way slab using bending moment coefficients
func CSlb1DepthCs(s *RccSlb) (error){
	//get depth of continuous one way slab using b.m coefficients
	var lspan float64
	ldrat := 26.0
	s.Init()
	//if err != nil{return err}
	//CHECK FOR RATIO of ll to dl < 0.75
	if s.Nspans < 3{
		return errors.New("invalid no. of spans for coefficient design")
	}
	for i := 0; i < s.Nspans; i++{
		s.Lspans = append(s.Lspans, s.Lspan)
	}
	lspan = s.Lspan/1000.0
	diamain := 10.0
	if s.Efcvr == 0.0{
		s.Efcvr = s.Nomcvr + diamain/2.0
	}
	efcvr := s.Efcvr
	fs := 5.0 * s.Fy/8.0/1.2
	var effd, wud, wul, mdu, msup, vud, dl, ll float64
	psfd := 1.5; psfl := 1.5
	if s.Code > 1{psfd = 1.4; psfl = 1.6}
	var iter int
	dused := 75.0
	dused -= 5.0
	var k2 float64
	switch s.Fy{
		case 250.0:
		k2 = 0.149
		case 415.0:
		k2 = 0.138
		default:
		k2 = 0.133
	}
	var astsup, ptreq float64
	for iter != -1{
		dused += 5.0
		wud = (0.025 * dused + s.DL) * psfd
		wul = s.LL * psfl
		dl = (0.025 * dused  + s.DL) * math.Pow(lspan, 2) * psfd
		ll = s.LL * math.Pow(lspan, 2) * psfl
		mdu = dl * CBmCspn[0] + ll * CBmCspn[2]
		msup = dl * CBmCsup[0] + ll * CBmCsup[2]		
		effd = dused - efcvr
		m1 := mdu*1e6/1000.0/math.Pow(effd,2)
		ldfac := 0.55 + (477.0 - fs)/120.0/m1
		if ldfac > 2.0 {ldfac = 2.0}
		dreq := (1000.0*lspan/ldfac/ldrat)
		astr := BalSecAst(mdu, effd, s.Fck, s.Fy, s.Code)
		ptreq = astr * 100.0/1000.0/effd
		d2 := RccSlabServeChk(1, 2, s.Fy, s.LL, lspan*1e3, ptreq,efcvr)
		if d2 - efcvr > dreq{dreq = d2 - efcvr}
		if dreq/effd - 1.0 > 0.0{
			continue
		}
		mmax := mdu
		if msup > mdu{mmax = msup}
		dbm := math.Sqrt(mmax * 1e3/k2/s.Fck)
		if dbm > effd{
			continue
		}
		astsup = BalSecAst(msup, effd, s.Fck, s.Fy, s.Code)
		vud = (wud * CVCdl[1] + wul * CVCll[1]) * lspan
		shrchk := SecShear(s.Fck, vud/effd, 1000.0, effd, dused, astsup, s.Code)
		if !shrchk{
			continue
		}
		iter = -1
		break
	}
	//fmt.Println("ptreq->",ptreq)
	s.Dused = dused
	//fmt.Println("final depth->",dused,"mm")
	
	dl = (0.025 * dused  + s.DL) * math.Pow(lspan, 2) * psfd
	ll = s.LL * math.Pow(lspan, 2) * psfl
	s.BMs = [][]float64{}
	for i := 0; i < s.Nspans; i++{
		bl, bm, br := CSlbBmCoeff(s.Nspans, i, dl, ll)
		s.BMs = append(s.BMs, []float64{bl, bm, br})
	}
	//rstr := fmt.Sprintf("%.1f\n",s.BMs)
	//fmt.Println(rstr)
	return nil
}

//CSlb1Stl calcs the steel required at left support, mid span, right support of each span of a continuous one way slab
func CSlb1Stl(s *RccSlb)(err error){
	//gets steel at b1, b2, b3 of a continuous slab
	s.Astrs = make([][]float64, len(s.BMs))
	s.Astps = make([][]float64, len(s.BMs))
	s.Spcspns = make([][]float64, len(s.BMs))
	s.Diaspns = make([][]float64, len(s.BMs))
	s.Sdspns = make([][]float64, len(s.BMs))
	s.Astds = make([][]float64, len(s.BMs))
	for i := 0; i < len(s.BMs); i++{
		s.Astrs[i] = make([]float64, 3)
		s.Astps[i] = make([]float64, 3)
		s.Spcspns[i] = make([]float64, 3)
		s.Diaspns[i] = make([]float64, 3)
		s.Sdspns[i] = make([]float64, 3)
		s.Astds[i] = make([]float64, 3)
	}
	effd := s.Dused - s.Efcvr
	//fmt.Println(effd)
	if s.Diadist == 0.0{
		s.Diadist = 8.0
	}
	for i := 0; i < s.Nspans; i++{
		for j := 0; j < 3; j++{
			bmr := math.Abs(s.BMs[i][j])
			if i == 0 && j == 0{
				//50 percent of span ast in top
				bmr = math.Abs(s.BMs[i][j+1])/2.0
			}
			if i == s.Nspans -1 && j == 2{
				//50 percent of span ast in top
				bmr = math.Abs(s.BMs[i][j-1])/2.0
			}
			ast := BalSecAst(bmr, effd, s.Fck, s.Fy, s.Code)
			s.Astrs[i][j] = ast
			rezmap, mindia := SlabRbrDiaSpcing(s.Dused, ast, s.Efcvr)
			if len(rezmap) == 0 {
				err = ErrSpacing
				return
			}
			asprov := rezmap[mindia][2]
			spcmain := rezmap[mindia][1]
			diacom := mindia
			if val, ok := rezmap[s.Diamain]; ok {
				asprov = val[2]
				spcmain = val[1]
				diacom = s.Diamain
			}
			s.Astps[i][j] = asprov
			s.Spcspns[i][j] = spcmain
			s.Diaspns[i][j] = diacom
			astd := 1.2 * s.Dused
			if s.Diadist == 6.0 {
				s.Fyd = 250.0
				astd = 1.5 * s.Dused
			} 
			sdmax := 5 * effd
			if sdmax > 300.0 {
				sdmax = 300.0
			}
			sds := 1000.0 * RbrArea(s.Diadist) / astd
			if sds > sdmax {
				sds = sdmax
			}
			sds = 5.0 * math.Floor(sds/5.0)
			asdprov := math.Round(1000.0 * RbrArea(s.Diadist) / sds)
			s.Astds[i][j] = asdprov
			s.Sdspns[i][j] = sds
		}
	}
	err = nil
	s.Dz = true
	return
}

//CSlb1D is unused, DELETE
func CSlb1D(s *RccSlb) (error){
	/*
	   IS THIS USES
	   designs a continuous one-way slab
	*/
	switch s.Dtyp{
		case 0:
		//use coefficients
		case 1:
		//use envelope analysis
	}
	return nil
}

//CSlbBmCoeff returns the span bending moment coefficients (left support, midspan, right support)
func CSlbBmCoeff(nspans, idx int, dl, ll float64) (bl, bm, br float64){
	//returns span bm (l, m, r)
	switch idx{
		case 0:
		//always left span
		bl = 0.0
		bm = dl/12.0 + ll/10.0
		br = -dl/10.0 - ll/9.0
		return
		case nspans - 1:
		//always right span
		br = 0.0
		bm = dl/12.0 + ll/10.0
		bl = -dl/10.0 - ll/9.0
		return
		case 1:
		bm = dl/16.0 + ll/12.0
		bl = -dl/10.0 + -ll/9.0
		switch nspans{
			case 3:
			//both ends nultimate sup
			br = bl
			default:
			//right end interior
			br = -dl/12.0 - ll/9.0
		}
		return
		case nspans - 2:
		//left end interior sup, right end nultimate
		bm = dl/16.0 + ll/12.0
		bl = -dl/12.0 - ll/9.0
		br = -dl/10.0 - ll/9.0
		return
		default:
		//interior span
		bm = dl/16.0 + ll/12.0
		bl = -dl/12.0 - ll/9.0
		br = bl		
	}
	return
}

/*


func CSlb1DCs(s *RccSlb) (error){
	var lspan float64
	ldrat := 26.0
	err := SlbInit(s)
	if err != nil{return err}
	//CHECK FOR RATIO of ll to dl < 0.75
	lspan = s.Lspan/1000.0
	efcvr := s.Efcvr
	fs := 5.0 * s.Fy/8.0/1.2
	var effd, wud, wul, mdu, msup, vud, dl, ll float64
	psfd := 1.5; psfl := 1.5
	if s.Code > 1{psfd = 1.4; psfl = 1.6}
	var iter int
	dused := 75.0
	dused -= 5.0
	var k2 float64
	switch s.Fy{
		case 250.0:
		k2 = 0.149
		case 415.0:
		k2 = 0.138
		default:
		k2 = 0.133
	}
	var astsup float64
	for iter != -1{
		dused += 5.0
		wud = (0.025 * dused + s.DL) * psfd
		wul = s.LL * psfl
		dl = (0.025 * dused  + s.DL) * math.Pow(lspan, 2) * psfd
		ll = s.LL * math.Pow(lspan, 2) * psfl
		mdu = dl * CBmCspn[0] + ll * CBmCspn[2]
		msup = dl * CBmCsup[0] + ll * CBmCsup[2]		
		effd = dused - efcvr
		m1 := mdu*1e6/1000.0/math.Pow(effd,2)
		ldfac := 0.55 + (477.0 - fs)/120.0/m1
		if ldfac > 2.0 {ldfac = 2.0}
		dreq := (1000.0*lspan/ldfac/ldrat)
		if dreq/effd - 1.0 > 0.0{
			continue
		}
		mmax := mdu
		if msup > mdu{mmax = msup}
		dbm := math.Sqrt(mmax * 1e3/k2/s.Fck)
		if dbm > effd{
			continue
		}
		//astr = BalSecAst(mdu, effd, s.Fck, s.Fy, s.Code)
		astsup = BalSecAst(msup, effd, s.Fck, s.Fy, s.Code)
		vud = (wud * CVCdl[1] + wul * CVCll[1]) * lspan
		shrchk := SecShear(s.Fck, vud/effd, 1000.0, effd, dused, astsup, s.Code)
		if !shrchk{
			continue
		}
		iter = -1
		break
	}
	fmt.Println(dused, astr, astsup, efcvr)
	//get moments at all sections - per span - lsup, midspan
	//first add midspan moments/steel, then support moments 
	var bml, bmm, bmr, astl, astm, astr float64
	for i := 0; i < s.Nspans; i++{
		switch i{
			case 0:
			//left end span
			bml = (dl * CBmCspn[0] + ll * CBmCspn[2]) * 0.5
			//right end span
			default:
			bmm = dl * CBmCspn[1] + ll * CBmCspn[3]
		}
		ast = BalSecAst(bm, effd, s.Fck, s.Fy, s.Code)
		astz = append(astz, ast)
	}
	for i := 0; i < s.Nspans; i++{
		switch i{
			case 0, s.Nspans:
			//end supports - 25% of span moment (mosley)
			//CHECK for flange action (60% of transverse steel)
			bm = (dl * CBmCspn[0] + ll * CBmCspn[2]) * 0.25
			case 1, s.Nspans-1:
			//penultimate supports
			bm = dl * CBmCsup[0] + ll * CBmCsup[2]
			default:
			//yadda yadda
			bm = dl * CBmCsup[1] + ll * CBmCsup[3]
		}
		ast = BalSecAst(bm, effd, s.Fck, s.Fy, s.Code)
		astz = append(astz, ast)
	}
	err = Slb1RbrDet(s, astz, dused, s.Code)
	if err != nil{
		log.Println("ERRORE,errore->rebar detailing error")
	}
	s.Dused = dused
	
	//PlotCSlb(s,"dumb")

	return nil
}

   //ribbed slab ss
   
func (s *RccSlb) RQuant(cb *Cbm){
	//quantify a ribbed slab
}

*/
