package barf

//flat slab design using subfrm design routines 
//this is impossibly hard to write for some strange unholy reason?

import (
	"fmt"
	"math"
	//"errors"
	kass "barf/kass"
)

//FltSlbSf designs a single subframe struct as a flat slab
func FltSlbSf(sf *SubFrm) (err error){
	//flat slab span design - either x or y
	err = CalcSubFrm(sf)
	if err != nil{
		return
	}
	schn := make(chan []interface{},len(sf.Beams))
	for _, i := range sf.Beams{
		go fslbdz(sf.Lspans, sf.Code, sf.Bmenv[i], sf.RcBm[i],sf.Advloads[i], schn)
	}
	for _ = range sf.Beams{
		_ = <- schn
	}
	fmt.Println("YEE-HAW")
	return
}

//FltSlbDz designs a sub frame struct for (x dir) sf.Lspans and (y dir) sf.Lbays
func FltSlbDz(sf *SubFrm) (err error){
	//subframe design of flat slab
	//lspans := sf.Lspans; lbays := sf.Lbays
	
	//sf.Init()
	if sf.Verbose{
		fmt.Println("first l-r (x) dir, internal strip")
	}
	sx := &SubFrm{}
	*sx = *sf
	sx.Lbay = sf.Ly
	sx.Title = fmt.Sprintf("%s_x",sf.Title)
	sx.Ax = 1
	//err = EqSubFrm(sx)
	err = FltSlbSf(sx)
	if err != nil{
		fmt.Println(err)
		return
	}
	if sf.Verbose{
		fmt.Println("next t-b (y) dir, internal strip")		
	}
	
	sy := &SubFrm{}
	*sy = *sf
	sy.Lbay = sf.Lx
	sy.Ax = 2
	sy.Lspans = sf.Lbays
	sy.Lbays = sf.Lspans
	sy.Title = fmt.Sprintf("%s_y",sf.Title)
	
	err = FltSlbSf(sy)
	if err != nil{
		fmt.Println(err)
		return
	}
	return
}

//fsbmneg calculates negative design moments for a flat slab 
func fsbmneg(code int, bm *kass.BmEnv, ldcs [][]float64) (m1, mspn, m2 float64){
	//negative design moments for a flat slab
	switch code{
		//yeay is 456 - critical section at d/2 from column face
		//is hc/2 etc applicable?
		case 1, 2:
		//bs 8110 - check for hc/2
		var nmax, mp, ml, mr, m0, m20 float64
		xl := bm.Hc/2.0
		xr := bm.Lspan - bm.Hc/2.0
		//
		//if bm.Spandx <0{
		//	xl = bm.Lsy; xr = bm.Lspan - bm.Rsy
		//}
		fmt.Println(ColorRed,bm.Id,"xl,xr->",xl, xr,ColorReset)
		if bm.Endrel && bm.Spandx < 0{
			switch bm.Spandx{
				case -1:
				xl = bm.Lsx
				case -2:
				xr = bm.Lspan - bm.Rsx
			}
		}
		for _, ld := range ldcs{
			nmax += ld[2]
		}
		mp = bm.Mpmax
		if bm.DM > 0.0{
			mp = bm.Mprmax
		} 
		
		msum := nmax * math.Pow(bm.Lspan - 2.0 * bm.Hc/3.0,2)/8.0
		if bm.DM == 0.0{
			ml, _ = kass.BmSfX(bm.Xs,bm.Mnenv, bm.Venv,xl)
			if ml >= 0.0{
				ml = bm.Ml
			}
			mr, _ = kass.BmSfX(bm.Xs,bm.Mnenv, bm.Venv,xr)
			if mr >= 0.0{
				mr = bm.Mr
			}
			m0 = bm.Mnenv[0]; m20 = bm.Mnenv[20]
		} else {
			ml, _ = kass.BmSfX(bm.Xs,bm.Mnrd, bm.Vrd,xl)
			if ml >= 0.0{
				ml = bm.Mlrd
			}
			mr, _ = kass.BmSfX(bm.Xs,bm.Mnrd, bm.Vrd,xr)
			if mr >= 0.0{
				mr = bm.Mrrd
			}
			m0 = bm.Mnrd[0]; m20 = bm.Mnrd[20]
		}
		mavg := math.Abs(ml + mr)/2.0 + mp
		fmt.Println("beam",bm.Id, "msum min",msum, "mavg",mavg)
		fmt.Println("beam",bm.Id, "mp, ml, mr",mp, ml, mr)
		fmt.Println("beam",bm.Id, "m0, m20",m0, m20)
		fmt.Println("beam",bm.Id, "mavg > msum?",mavg > msum)
		mspn = mp
		m1 = ml; m2 = mr
		if mavg > msum{
			fmt.Println("final beam",bm.Id, "m1 (left), m2 (right)",m1, m2)
			return
		}
		delta := msum - mavg
		fmt.Println("adjusting negative moments->",bm.Id, m1, m2)
		m1 -= delta; m2 -= delta
		//if math.Abs(m1) < math.Abs(m2){
		//	m1 = -2.0 * delta + m1; m2 = m2
		//} else {
		//	m2 = -2.0 * delta + m2; m1 = m1
		//}
		fmt.Println("final adj. beam",bm.Id, delta, "m1 (left), m2 (right)?",m1, m2)
	}
	fmt.Println("yeehaw")
	return
}

//fslbdz designs a flat slab span (goroutine entry func)
func fslbdz(lspans []float64,code int, bm *kass.BmEnv, bmarr []*RccBm, ldcs [][]float64, schn chan []interface{}){
	//designs a flat slab span
	rez := make([]interface{},2)
	m1, mspn, m2 := fsbmneg(code, bm, ldcs) 
	
	rez[0] = m1; rez[1] = m2
	fmt.Println(rez)
	fmt.Println("design moments- left support->",m1,"span->",mspn,"right support->",m2)
	fmt.Println("dims->",bmarr[1].Dims, len(bmarr))
	ly := bm.Lspan
	lx := bm.Wspan; if lx > bm.Lspan{
		lx = bm.Lspan
	}
	fmt.Println("span->",bm.Lspan,"bay->",bm.Wspan,"lx->",lx,"ly->",ly)
	var cstrip, mstrip float64
	cstrip = lx/4.0; mstrip = ly - lx/2.0
	fmt.Println("cstrip, mstrip->",cstrip, mstrip)
	//var bms []float64
	for i := 0; i < 6; i++{
		switch i{
			case 0:
			//col strip left
			case 1:
			case 2:
			case 3:
			//mid strip left
			case 4:
			case 5:
		}
	}
	schn <- rez
	
}

//fsbmcoeff returns flat slab bending moment coefficients for simplified analysis
func fsbmcoeff(code, spandx int) (bmc []float64){
	switch code{
		case 1:
		case 2:
		bmc = []float64{0.65,0.35,0.65}
	}
	return
}
