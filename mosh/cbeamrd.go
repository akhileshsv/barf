package barf

import (
	//"log"
	"math"
	kass"barf/kass"
)

//CBeamDm redistributes support moments for a continuous beam
//see hulse section 2.5
func CBeamDM(ncjt int, bmvec []int, bmenv map[int]*kass.BmEnv, dm float64, ms map[int]*kass.Mem, mslmap map[int]map[int][][]float64){
	/*
	 moment redistribution for a continuous beam (wonky)
	*/
	//TODO if dmr == 1.0 -> get max possible red at supports without reducing max span moment
	if dm > 0.3 {dm = 0.3}
	
	var l , e, ar, iz float64
	var mem *kass.Mem
	var r kass.BeamRez
	nspans := len(bmvec)
	nlp := nspans+1
	bmsupz := make(map[int][]float64)
	//get redistributed support moments
	for idx := 2; idx <= nspans; idx++{
		lp := idx + 1
		//log.Println("at support->",idx)
		//log.Println("load pattern->",lp-1)
		bl := bmvec[idx-2]; br := bmvec[idx-1]
		//AI. THIS MUST CHANGES
		//switch ncjt{}
		var blml, blmr, brml, brmr float64
		switch ncjt{
			case 2:
			blml = bmenv[bl].EnvRez[lp].Qf[1]; blmr = bmenv[bl].EnvRez[lp].Qf[3]
			brml = bmenv[br].EnvRez[lp].Qf[1]; brmr = bmenv[br].EnvRez[lp].Qf[3]
			case 3:
			blml = bmenv[bl].EnvRez[lp].Qf[1]; blmr = bmenv[bl].EnvRez[lp].Qf[3]
			brml = bmenv[br].EnvRez[lp].Qf[1]; brmr = bmenv[br].EnvRez[lp].Qf[3]
		}
		//log.Println("end moments init->",blml, blmr,brml, brmr)
		delta := math.Abs(blmr)- math.Abs(brml)
		//log.Println("DELTA->",delta)
		switch{
			case delta <= 1.0:
			brml = (1.0-dm)*brml; blmr = (1.0-dm)*blmr
			case delta > 1.0:
			if math.Abs(blmr) < math.Abs(brml){
				//reduce blmr by rd
				blmr = blmr*(1.0-dm)
				brml = brml * (1.0 - math.Abs(dm * blmr)/math.Abs(brml))
			} else {
				brml = brml *(1.0-dm)
				blmr = blmr * (1.0 - math.Abs(dm*brml)/math.Abs(blmr))
			}
		}
		//log.Println("end moments final->",blml, blmr, brml, brmr)
		bmsupz[idx] = []float64{blml, blmr, brml,brmr}
		for _, bar := range []int{bl, br}{
			mem = ms[bar]
			bm := bmenv[bar]
			switch{
				case ncjt == 2://beam
				l = mem.Geoms[0]
				e = mem.Geoms[1]
				iz = mem.Geoms[2]
				if len(mem.Geoms) > 3{
					ar = mem.Geoms[3]
				}
				case ncjt == 3://frame
				l = mem.Geoms[0]
				e = mem.Geoms[1]
				ar = mem.Geoms[2]
				iz = mem.Geoms[3]
			}
			xdiv := l/20.0
			lsx := bmenv[bar].Lsx; rsx := l - bmenv[bar].Rsx
			il := int(math.Ceil(lsx/xdiv)); ir := int(math.Ceil(rsx/xdiv))
			if bar == bl{
				r = kass.BmCalcDM(bar, mslmap[0][bar], l, e, ar, iz, blml, blmr, false) //PREV-bm.EnvRez[0].Qf[1]
			}
			if bar == br{
				r = kass.BmCalcDM(bar, mslmap[0][bar], l, e, ar, iz, brml, brmr, false) //PREV-bm.EnvRez[0].Qf[3]
			}
			var vl, ml, vr, mr float64
			for i, vx := range r.SF{
				
				x := r.Xs[i]
				if math.Abs(bm.Vrd[i]) < math.Abs(vx) {
					bm.Vrd[i] = vx
				}
				if math.Abs(bm.Mnrd[i]) < math.Abs(r.BM[i]) && r.BM[i] < 0.0 {
					bm.Mnrd[i] = r.BM[i]
				}
				if r.BM[i] > 0.0 && math.Abs(bm.Mprd[i]) < math.Abs(r.BM[i]){
					bm.Mprd[i] = r.BM[i]
				}
				if math.Abs(bm.Vrmax) < math.Abs(vx) {
					bm.Vrmax = vx
					bm.Vrmaxx = x
				}
				if math.Abs(bm.Mnrmax) < math.Abs(r.BM[i]) && r.BM[i] < 0.0 {
					bm.Mnrmax = r.BM[i]
					bm.Mnrmaxx = x
				}
				if r.BM[i] > 0.0 && math.Abs(bm.Mprmax) < math.Abs(r.BM[i]) {
					bm.Mprmax = r.BM[i]
					bm.Mprmaxx = x
				}
				
				if (i == il || math.Abs(x-lsx) <= xdiv/2.0) && (vl + ml == 0.0){
					switch{
						case math.Abs(x-lsx) <= xdiv/2.0:
						vl = vx
						ml = r.BM[i]
						case x == lsx:
						vl = vx
						ml = r.BM[i]
						default:
						switch i{
							case 0:
							vl = r.SF[i+1]
							ml = r.BM[i+1]
							default:
							vl = vx + (vx - r.SF[i-1])*(lsx - x)/xdiv
							ml = r.BM[i] + 0.5 * (lsx - x)*(vl + vx)

						}
					}
					if math.Abs(bm.Vlrd) < math.Abs(vl){bm.Vlrd = vl}
					if math.Abs(bm.Mlrd) < math.Abs(ml){bm.Mlrd = ml}
				}
				if (i == ir || math.Abs(x-rsx) <= xdiv/2.0) && (vr + mr == 0.0){
					switch{
						case math.Abs(x-rsx) <= xdiv/2.0:
						vr = vx
						mr = r.BM[i]
						case x == rsx:
						vr = vx
						mr = r.BM[i]
						default:
						vr = vx + (vx - r.SF[i-1])*(rsx - x)/xdiv
						mr = r.BM[i] + 0.5 * (rsx - x)*(vr + vx) 
					}
					if math.Abs(bm.Vrrd) < math.Abs(vr){bm.Vrrd = vr}
					if math.Abs(bm.Mrrd) < math.Abs(mr){bm.Mrrd = mr}
				}
			}
		}
	}
	for lp := 1; lp <= nlp; lp++{
		var sup, bl, br int
		var blml, blmr, brml, brmr float64
		if lp > 2{
			sup = lp - 1
			bl = bmvec[sup-2]; br = bmvec[sup-1]
			blml, blmr, brml,brmr = bmsupz[sup][0], bmsupz[sup][1], bmsupz[sup][2], bmsupz[sup][3]
		}
		for _, bar := range bmvec{
			mem = ms[bar]
			bm := bmenv[bar]
			switch{
				case ncjt == 2://beam
				l = mem.Geoms[0]
				e = mem.Geoms[1]
				iz = mem.Geoms[2]
				if len(mem.Geoms) > 3{
					ar = mem.Geoms[3]
				}
				case ncjt == 3://frame
				l = mem.Geoms[0]
				e = mem.Geoms[1]
				ar = mem.Geoms[2]
				iz = mem.Geoms[3]
			}
			xdiv := l/20.0
			lsx := bmenv[bar].Lsx; rsx := l - bmenv[bar].Rsx
			il := int(math.Floor(lsx/xdiv)); ir := int(math.Floor(rsx/xdiv))
			r = bm.EnvRez[lp]
			if bar == bl{
				r = kass.BmCalcDM(bar, mslmap[lp][bar], l, e, ar, iz, blml, blmr, false)
			}
			if bar == br{
				r = kass.BmCalcDM(bar, mslmap[lp][bar], l, e, ar, iz, brml, brmr, false)
			}
			var vl, ml, vr, mr float64
			for i, vx := range r.SF{
				x := r.Xs[i]
				if math.Abs(bm.Vrd[i]) < math.Abs(vx) {
					bm.Vrd[i] = vx
				}
				if math.Abs(bm.Mnrd[i]) < math.Abs(r.BM[i]) && r.BM[i] < 0.0 {
					bm.Mnrd[i] = r.BM[i]
				}
				if r.BM[i] > 0.0 && math.Abs(bm.Mprd[i]) < math.Abs(r.BM[i]){
					bm.Mprd[i] = r.BM[i]
				}
				if math.Abs(bm.Vrmax) < math.Abs(vx) {
					bm.Vrmax = vx
					bm.Vrmaxx = x
				}
				if math.Abs(bm.Mnrmax) < math.Abs(r.BM[i]) && r.BM[i] < 0.0 {
					bm.Mnrmax = r.BM[i]
					bm.Mnrmaxx = x
				}
				if r.BM[i] > 0.0 && math.Abs(bm.Mprmax) < math.Abs(r.BM[i]) {
					bm.Mprmax = r.BM[i]
					bm.Mprmaxx = x
				}
				
				if (i == il || math.Abs(x-lsx) <= xdiv/2.0 && vl + ml == 0.0){
					switch{
						case math.Abs(x-lsx) <= xdiv/2.0:
						vl = vx
						ml = r.BM[i]
						case x == lsx:
						vl = vx
						ml = r.BM[i]
						default:
						switch i{
							case 0:
							vl = r.SF[i+1]
							ml = r.BM[i+1]
							default:
							vl = vx + (vx - r.SF[i-1])*(lsx - x)/xdiv
							ml = r.BM[i] + 0.5 * (lsx - x)*(vl + vx)
						}
					}
					if math.Abs(bm.Vlrd) < math.Abs(vl){bm.Vlrd = vl}
					if math.Abs(bm.Mlrd) < math.Abs(ml){bm.Mlrd = ml}
				}
				if (i == ir || math.Abs(x-rsx) <= xdiv/2.0) && (vr + mr == 0.0){
					switch{
						
						case math.Abs(x-rsx) <= xdiv/2.0:
						vr = vx
						mr = r.BM[i]
						case x == rsx:
						vr = vx
						mr = r.BM[i]
						default:
						vr = vx + (vx - r.SF[i-1])*(rsx - x)/xdiv
						mr = r.BM[i] + 0.5 * (rsx - x)*(vr + vx) 
					}
					if math.Abs(bm.Vrrd) < math.Abs(vr){bm.Vrrd = vr}
					if math.Abs(bm.Mrrd) < math.Abs(mr){bm.Mrrd = mr}
				}
			}
		}
	}
	return
}



/*
   YE OLDE

   
func CBeamDM(ncjt int, bmvec []int, bmenv map[int]*RccBmEnv, dm float64, ms map[int]*kass.Mem, mslmap map[int]map[int][][]float64){
	//TODO if dmr == 1.0 -> get max possible red at supports without reducing max span moment
	log.Println("HYARR SHE GOES moment redistribushun")
	log.Println(bmvec)
	var l , e, ar, iz float64
	var mem *kass.Mem
	var r BeamRez
	nspans := len(bmvec)
	bmsupz := make(map[int][]float64)
	
	for idx := 2; idx <= nspans; idx++{
		lp := idx + 1
		log.Println("at support->",idx)
		log.Println("load pattern->",lp-1)
		bl := bmvec[idx-2]; br := bmvec[idx-1]
		blml := bmenv[bl].EnvRez[lp].Qf[1]; blmr := bmenv[bl].EnvRez[lp].Qf[3]
		brml := bmenv[br].EnvRez[lp].Qf[1]; brmr := bmenv[br].EnvRez[lp].Qf[3]
		log.Println("end moments init->",blml, blmr,brml, brmr, )
		brml = (1.0-dm)*brml; blmr = (1.0-dm)*blmr
		log.Println("end moments final->",blml, blmr, brml, brmr )
		for _, bar := range []int{bl, br}{
			mem = ms[bar]
			switch{
				case ncjt == 2://beam
				l = mem.Geoms[0]
				e = mem.Geoms[1]
				iz = mem.Geoms[2]
				if len(mem.Geoms) > 3{
					ar = mem.Geoms[3]
				}
				case ncjt == 3://frame
				l = mem.Geoms[0]
				e = mem.Geoms[1]
				ar = mem.Geoms[2]
				iz = mem.Geoms[3]
			}
			//xdiv := l/20.0
			//lsx := bmenv[bar].Lsx; rsx := l - bmenv[bar].Rsx
			//il := int(math.Ceil(lsx/xdiv)); ir := int(math.Ceil(rsx/xdiv))
			switch{
				case bar == bl:
				ml := blml; mr := blmr
				//rl := (ml+mr)/l; rr := -rl
				//r = kass.BmCalcDM(bar, mslmap[lp][bar],l, e, ar, iz, blml, blmr)
				//log.Println("")
				//r = Bmsfcalc(bar, mslmap[lp][bar], l, e, ar, iz, rl, ml, rr, mr, false)
				r = kass.BmCalcDM(bar, mslmap[lp][bar], l, e, ar, iz, ml, mr, true)
				//j = idx - 1
				default:
				ml := brml; mr := brmr
				//rl := (ml+mr)/l; rr := -rl
				//r = Bmsfcalc(bar, mslmap[lp][bar], l, e, ar, iz, rl, ml, rr, mr, false)
				r = kass.BmCalcDM(bar, mslmap[lp][bar], l, e, ar, iz, ml, mr, true)
				//j = idx
			}
			//fmt.Println("shear",bar,r.SF)
			for i, vx := range r.SF{
				log.Println("section-",i)
				log.Println(ColorCyan)
				log.Println("before")
				log.Println("vx",bmenv[bar].EnvRez[lp].SF[i],"bm",bmenv[bar].EnvRez[lp].BM[i])
				log.Println(ColorGreen)
				log.Println("after")
				log.Println("vx",vx,"bm",r.BM[i])
				log.Println(ColorReset)
				if math.Abs(bmenv[bar].Vrd[i])<math.Abs(vx){
					bmenv[bar].Vrd[i] = vx
				}
				if math.Abs(bmenv[bar].Mnrd[i])<math.Abs(r.BM[i]) && r.BM[i] < 0{
					bmenv[bar].Mnrd[i] = r.BM[i]
				}
				if math.Abs(bmenv[bar].Mprd[i])<math.Abs(r.BM[i]) && r.BM[i] > 0{
					bmenv[bar].Mprd[i] = r.BM[i]
				}
			}
		}
	}
}



*/


/*
  		//if math.Abs(blmr) < math.Abs(brml) {
		//	delta = math.Abs(blmr * dm)
		//} else {
		//	delta = math.Abs(brml * dm)
		//}
		//now actual lifting
		//blmr = blmr*(1.0 + delta/blmr)
		//brml = brml*(1.0 + delta/brml)
		//fmt.Println(blml, blmr, brml, brmr)
 
	x := r.Xs[i]
				   
	for i := range vs{
		vs[i] = make([]float64, 21)
		mns[i] = make([]float64, 21)
		mps[i] = make([]float64, 21)
	}
				if math.Abs(vs[j][i])<math.Abs(vx){
					vs[j][i] = vx
				}
				if math.Abs(mns[j][i])<math.Abs(r.BM[i]) && r.BM[i] < 0{
					mns[j][i] = r.BM[i]
				}
				if math.Abs(mps[j][i])<math.Abs(r.BM[i]) && r.BM[i] > 0{
					mps[j][i] = r.BM[i]
				}

	//var vl, vr, ml, mr float64
	//var j int
	//vs, mns, mps = make([][]float64, len(bmvec)), make([][]float64, len(bmvec)), make([][]float64, len(bmvec))
	//mspns, xspns, vls, vrs, mls, mrs := make([]float64, len(bmvec)), make([]float64, len(bmvec)), make([]float64, len(bmvec)), make([]float64, len(bmvec)), make([]float64, len(bmvec)), make([]float64, len(bmvec))
				*/
