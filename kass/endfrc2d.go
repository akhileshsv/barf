package barf

//bunch of 2d endforce computation functions, mostly from kassimali

import (
	//"log"
	"math"
	"gonum.org/v1/gonum/mat"
)

//MemFrc returns member end forces (msqf) and builds pf vec(p - pf = sd)
func MemFrc(ms map[int]*Mem, msp [][]float64, nsc []int, ncjt int, ndof int, mfrtyp int, pfchn chan []interface{}) {
	//map of loaded members msloaded[0] ex. - 1,1,40.25,0,134.16,0
	//msqf - map of local fixed end forces/end reactions
	//mfrtyp - 0 - 1db, 1- 3dg, 2- 2df, 3- 3df
	msloaded := make(map[int][][]float64)
	qfchn := make(chan []interface{}, len(msp))
	for _, ldcase := range msp{
		if len(ldcase) < 6{
			continue
		}
		member := int(ldcase[0])
		mem := ms[member]
		msloaded[member] = append(msloaded[member], ldcase)
		ms[member].Lds = append(ms[member].Lds, ldcase)
		ltyp := int(ldcase[1])
		var memrel int
		if len(mem.Mprp) > 4 {memrel = mem.Mprp[4]}
		l := mem.Geoms[0]
		go FxdEndFrc(member, memrel, ltyp, l, ldcase[2:], mfrtyp, mem.Geoms, qfchn)
	}
	msqf := make(map[int][]float64)
	for member := range ms {
		msqf[member] = make([]float64, 2*ncjt)
	}
	//get member FIXED END forces qf
	for i := 0; i < len(msp); i++ {
		r := <-qfchn
		member, _ := r[0].(int)
		qf, _ := r[1].([]float64)
		for i, f := range qf {
			msqf[member][i] += f
		}
	}
	//assemble global fixed end force vector
	
	
	ffchn := make(chan [][]float64, len(msloaded))
	for member, mem  := range ms {
		switch mfrtyp {
		case 0: //1d beam has no tmat no gkmat no roof or home
			_, ok := msloaded[member]
			if ok {
				qf := msqf[member]
				go ffmatb(mem.Mprp[0], mem.Mprp[1], qf, nsc, ncjt, ndof, ffchn)
			} else {
				continue
			}
		default:
			_, ok := msloaded[member]
			if ok {
				qf := msqf[member]
				go ffmat(mem.Mprp[0], mem.Mprp[1], qf, mem.Tmat, nsc, ncjt, ndof, ffchn)
			} else {
				continue
			}
		}
	}
	//get global fixed end force vector pf
	pf := make([]float64, ndof)
	for i := 0; i < len(msloaded); i++ {
		rez := <-ffchn
		for _, r := range rez {
			pf[int(r[0])] += r[1]
		}
	}
	//log.Println(ColorRed,"pf->",pf,ColorReset)
	rez := make([]interface{}, 3)
	rez[0] = msqf
	rez[1] = pf
	rez[2] = msloaded
	pfchn <- rez
}

//NodeFrc generates nodal force vector p from list of nodal/joint loads jp
func NodeFrc(jp [][]float64, nsc []int, ncjt int, ndof int, pchn chan []float64) {
	p := make([]float64, ndof)
	if len(jp) == 1 && len(jp[0]) == 0 {
		pchn <- p
		return
	} else {
		for _, pj := range jp{
			switch {
			case ncjt == 2 && len(pj) < 3:
				continue
			case ncjt == 3 && len(pj) < 4:
				continue
			case ncjt == 6 && len(pj) < 7:
				continue
			}
			j := int(pj[0])
			idx := (j - 1) * ncjt
			for i := 0; i < ncjt; i++ {
				n := nsc[idx+i]
				if n <= ndof {
					p[n-1] += pj[i+1]
				}
			}
		}
	}
	pchn <- p
}

//EndFrc gets member end forces from global displ list
//local end displ u = t (tmat ofc) *v (end global displacement),
//q = bk* u, ff = t.T * q
//returns rez for building reaction list
func EndFrc(member int, mem *Mem, msqf map[int][]float64, dglb []float64, nsc []int, ncjt int, ndof int, rchn chan [][]float64, fchn chan []interface{}) {
	qfmat := mat.NewDense(len(msqf[member]), 1, msqf[member])
	vmat := mat.NewDense(2*ncjt, 1, nil)
	jb := mem.Mprp[0]
	je := mem.Mprp[1]
	var n int
	for i := 1; i <= 2*ncjt; i++ {
		if i <= ncjt {
			n = (jb-1)*ncjt + i
		} else {
			n = (je-1)*ncjt + (i - ncjt)
		}
		vmat.Set(i-1, 0, dglb[n-1]+mem.Vfs[i-1])
	}
	u := mat.NewDense(2*ncjt, 1, nil)
	u.Mul(mem.Tmat, vmat)
	qmat := mat.NewDense(2*ncjt, 1, nil)
	qmat.Mul(mem.Bkmat,u)
	qmat.Add(qmat,qfmat)
	fmat := mat.NewDense(2*ncjt, 1, nil)
	fmat.Mul(mem.Tmat.T(),qmat)
	//return index and values of fmat (>ndof) for support reaction r
	rez := [][]float64{}
	fvec := make([]float64, 2*ncjt)
	qfvec := make([]float64, 2*ncjt)
	for i := 1; i <= 2*ncjt; i++ {
		fvec[i-1] = fmat.At(i-1,0)
		qfvec[i-1] = qmat.At(i-1,0)
		if i <= ncjt {
			n = (jb-1)*ncjt + i
		} else {
			n = (je-1)*ncjt + (i - ncjt)
		}
		sc := nsc[n-1]
		if sc > ndof {
			rez = append(rez, []float64{float64(sc - ndof - 1), fmat.At(i-1, 0)})
		}
	}
	frez := make([]interface{},3)
	frez[0] = member
	frez[1] = fvec
	frez[2] = qfvec
	
	rchn <- rez
	fchn <- frez	
}

//ffmat assembles the structure fixed end force vector elements from individual member matrices
func ffmat(jb int, je int, qf []float64, tmat *mat.Dense, nsc []int, ncjt int, ndof int, ffchn chan [][]float64) {
	rez := [][]float64{}
	qfmat := mat.NewDense(2*ncjt, 1, nil)
	for i, f := range qf {qfmat.Set(i,0,f)}
	//ffmat := mat.NewDense(2*ncjt,1,nil)
	qfmat.Product(tmat.T(), qfmat)
	var el int
	for i := 1; i <= 2*ncjt; i++ {
		if i <= ncjt {
			el = (jb-1)*ncjt + i
		} else {
			el = (je-1)*ncjt + (i - ncjt)
		}
		n := nsc[el-1]
		if n <= ndof {
			rez = append(rez, []float64{float64(n - 1), qfmat.At(i-1, 0)})
		}
	}
	ffchn <- rez
}

//ffmatb identifies the structure fixed end force vector elements from individual member matrices
func ffmatb(jb int, je int, qf []float64, nsc []int, ncjt int, ndof int, ffchn chan [][]float64) {
	rez := [][]float64{}
	var el int
	for i := 1; i <= 2*ncjt; i++ {
		if i <= ncjt {
			el = (jb-1)*ncjt + i
		} else {
			el = (je-1)*ncjt + (i - ncjt)
		}
		n := nsc[el-1]
		if n <= ndof {
			rez = append(rez, []float64{float64(n - 1), qf[i-1]})
		}
	}
	ffchn <- rez
}




//FxdEndFrc generates fixed end force vectors for each member for each individual load type
//see back cover of kassimali 
//mfrtyp - -2 - TRUSS 2D, -1 - TRUSS 3D 0 - BEAM, 1- GRID, 2 - FRAME, 3 - 3D FRAME
func FxdEndFrc(member int, memrel int, ltyp int, l float64, ldet []float64, mfrtyp int, geoms []float64, qfchn chan []interface{}) {
	//ldet = [wa, wb, la, lb, ltyp, ldaxis], mfrtyp - member frame type, 0 - 1d beam, 1 - 3d grid, 2 - 2d frame, 3 - 3dframe
	//memrel - member release type 0 - fixed, 1 - hinge at beginning, 2- hinge at end, 3- hinge at both ends
	w := ldet[0]
	wb := ldet[1]
	la := ldet[2]
	lb := ldet[3]
	wa := w
	var fab, fsb, fmb, fae, fse, fme, ftb, fte float64
	switch ltyp {
	case 1:
		//point load w at la from end a (lb = l -la)
		lb = l - la
		fsb = w * math.Pow(lb, 2) * (3*la + lb) / math.Pow(l, 3)
		fmb = w * la * math.Pow(lb, 2) / math.Pow(l, 2)
		fse = w * math.Pow(la, 2) * (la + 3*lb) / math.Pow(l, 3)
		fme = -w * math.Pow(la, 2) * lb / math.Pow(l, 2)
		fab = 0.0
		fae = 0.0
	case 2:
		//moment (m=w)  m at la (lb = l-la)
		lb = l - la
		fsb = -6 * w * la * lb / math.Pow(l, 3)
		fmb = w * lb * (lb - 2*la) / math.Pow(l, 2)
		fse = 6 * w * la * lb / math.Pow(l, 3)
		fme = w * la * (la - 2*lb) / math.Pow(l, 2)
		fab = 0.0
		fae = 0.0
	case 3:
		//udl w at la, lb (=0 -> all over)
		fsb = (1.0 / 2.0) * w * l * (1.0 - la*(1/math.Pow(l, 4))*(2*math.Pow(l, 3)-2*math.Pow(la, 2)*l+math.Pow(la, 3)) - math.Pow(lb, 3)*(1.0/math.Pow(l, 4))*(2*l-lb))
		fmb = (1.0 / 12.0) * w * math.Pow(l, 2) * (1.0 - (1.0/math.Pow(l, 4))*math.Pow(la, 2)*(6*math.Pow(l, 2)-8*la*l+3*math.Pow(la, 2)) - (1.0/math.Pow(l, 4))*math.Pow(lb, 3)*(4*l-3*lb))
		fse = (1.0 / 2.0) * w * l * (1.0 - (1.0/math.Pow(l, 4))*math.Pow(la, 3)*(2*l-la) - lb*(1/math.Pow(l, 4))*(2*math.Pow(l, 3)-2*math.Pow(lb, 2)*l+math.Pow(lb, 3)))
		fme = -(1.0 / 12.0) * w * math.Pow(l, 2) * (1.0 - (1/math.Pow(l, 4))*math.Pow(la, 3)*(4*l-3*la) - (1/math.Pow(l, 4))*math.Pow(lb, 2)*(6*math.Pow(l, 2)-(8*lb*l)+3*math.Pow(lb, 2)))
		fab = 0
		fae = 0

	case 4:
		//linearly varying udl (start,stop) wa,wb at la,lb
		var ta, tb, tx, ty, tz float64

		tz = 1.0 + (lb / (l - la)) + math.Pow(lb, 2)/math.Pow((l-la), 2)

		tx = (7*l + 8*lb) - (1.0/(l-la))*lb*(3*l+2*lb)*(tz) + 2*math.Pow(lb, 4)*(1.0/math.Pow((l-la), 3))
		ta = (1.0 / 20.0) * (1 / math.Pow(l, 3)) * wa * math.Pow((l-la), 3) * (tx)

		ty = (3*l+2*la)*tz - (1.0/math.Pow((l-la), 2))*math.Pow(lb, 3)*(2+(15*l-8*lb)/(l-la))
		tb = (1.0 / 20.0) * (1 / math.Pow(l, 3)) * wb * math.Pow((l-la), 3) * (ty)

		fsb = ta + tb
		//end fsb

		tx = 3*(l+4*lb) - (lb * (2*l + 3*lb) * tz / (l - la)) + 3*math.Pow(lb, 4)*(1.0/math.Pow((l-la), 3))
		ta = (1.0 / 60.0) * (1 / math.Pow(l, 2)) * wa * math.Pow((l-la), 3) * tx

		ty = (2*l+3*la)*tz - 3*math.Pow(lb, 3)*(1+(5*l-4*lb)/(l-la))/math.Pow((l-la), 2)
		tb = (1.0 / 60.0) * (1.0 / math.Pow(l, 2)) * wb * math.Pow((l-la), 3) * ty

		fmb = ta + tb
		//end fmb

		fse = (1.0/2.0)*(wa+wb)*(l-la-lb) - fsb
		//end fse

		fme = (1.0/6.0)*(l-la-lb)*(wa*(2*la-lb-2*l)-wb*(l-la+2*lb)) + fsb*l - fmb
		//end fmb

		fab = 0
		fae = 0
	case 5:
		//axial load w (start,stop) la, lb
		//lb = l - la
		lb = l - la
		fsb = 0
		fmb = 0
		fse = 0
		fme = 0
		fab = w * lb / l
		fae = w * la / l
	case 6:
		//axial load along axis w (start,stop) la, lb
		fsb = 0
		fmb = 0
		fse = 0
		fme = 0
		fab = (1.0 / 2 * l) * w * (l - la - lb) * (l - la + lb)
		fae = (1.0 / 2 * l) * w * (l - la - lb) * (l + la - lb)
	case 7:
		//torsional moment Mt in z w (start,stop) la, lb
		//
		fsb = 0
		ftb = w * (l - la) / l
		fse = 0
		fte = w * la / l
		fab = 0
		fae = 0
	case 8:
		//temp change ta, tb (wa, wb), alpha, d
		//if d == 0 or ta == tb 'TIS UNIFORM and no bending moments are induced
		ta := ldet[0]
		tb := ldet[1]
		alpha := ldet[2]
		depth := ldet[3]
		switch mfrtyp{
			case -2:
			//2d truss
			//{l, e, a, cx, cy}
			fab = geoms[1] * geoms[2] * alpha * (tb + ta)/2.0
			case -1:
			//{l, e, a, cx, cy, cz}
			fab = geoms[1] * geoms[2] * alpha * (tb + ta)/2.0
			case 0:
			//{l, e, iz, ar}
			//fab = 0.0 (as beams are free to expand axially)
			if tb == ta || depth == 0.0{
				fmb = 0
			} else {
				fmb = geoms[1] * geoms[2] * alpha * (tb - ta)/depth
			}
			default:
			//l, e, a, iz, cx, cy
			fab = geoms[1] * geoms[2] * alpha * (tb + ta)/2.0
			if tb == ta || depth == 0.0{
				fmb = 0.0
			} else {
				fmb = geoms[1] * geoms[3] * alpha * (tb - ta)/depth
			}
		}
		fae = -fab
		fme = -fmb
	case 9:
		//error in initial member length
		delta := ldet[0]
		switch mfrtyp{
			case 0:
			default:
			fab = delta * geoms[1] * geoms[2]/geoms[0]
			fae = -fab
		}
	case 10:
		//error in member straightness (bow)
		//only for beams and frames (?)
		delta := ldet[0]; l1 := ldet[2]; l2 := ldet[3]
		iz := geoms[2]
		switch mfrtyp{
			case 0:
			//{l, e, iz, ar}
			default:
			//l, e, a, iz, cx, cy
			iz = geoms[3]
		}
		n1 := 2.0 * geoms[1] * iz * delta/math.Pow(geoms[0],2)/l1/l2
		fsb = n1 * 3.0 * (l2 - l1)
		fse = -fsb
		fmb = n1 * geoms[0] * (2.0 * l2 - l1)
		fme = n1 * geoms[0] * (l2 - 2.0 * l1)
	}
	rez := make([]interface{}, 2)
	rez[0] = member
	//mfrtyp 0 - beam, 1 - grid, 2 - 2df, 3 - 3df
	switch memrel {
	case 0: //fixed end, standard
		switch mfrtyp {
		case -2:
			//2d truss
			rez[1] = []float64{fab, fsb, fae, fse}
		case -1:
			//3d truss
			rez[1] = []float64{fab, 0, 0, fae, 0, 0}
		case 0: //beam
			rez[1] = []float64{fsb, fmb, fse, fme}
		case 1: //grid
			rez[1] = []float64{fsb, ftb, fmb, fse, fte, fme}
		case 2: //2d frame
			rez[1] = []float64{fab, fsb, fmb, fae, fse, fme}
		case 3: //3d frame
			var fsby, fmby, fsey, fmey, fsbz, fmbz, fsez, fmez float64
			ldaxis := ldet[5]
			switch ldaxis {
			case 0: //along local y axis
				fsby = fsb
				fmby = 0
				fsey = fse
				fmey = 0
				fsbz = 0
				fmbz = fmb
				fsez = 0
				fmez = fme
			case 1: //along local z axis
				fsby = fsb
				fmby = 0
				fsey = fse
				fmey = 0
				fsbz = 0
				fmbz = fmb
				fsez = 0
				fmez = fme
			}
			rez[1] = []float64{fab, fsby, fsbz, ftb, fmby, fmbz, fae, fsey, fsez, fte, fmey, fmez}
		}
	case 1: //hinge at beginning
		switch mfrtyp {
		case 0: //beam
			rez[1] = []float64{fsb - 3.0*fmb/(2.0*l), 0, fse + 3.0*fmb/(2.0*l), fme - fmb/2.0}
		case 1: //grid
			rez[1] = []float64{fsb - 3.0*fmb/(2.0*l), 0, 0, fse + 3.0*fmb/(2.0*l), fte + ftb, fme - fmb/2.0}
		case 2: //2d frame
			rez[1] = []float64{fab, fsb - 3.0*fmb/(2.0*l), 0, fae, fse + 3.0*fmb/(2.0*l), fme - fmb/2.0}
		case 3: //3d frame
			var fsby, fmby, fsey, fmey, fsbz, fmbz, fsez, fmez float64
			ldaxis := ldet[5]
			switch ldaxis {
			case 0: //along local y axis
				fsby = fsb
				fmby = 0
				fsey = fse
				fmey = 0
				fsbz = 0
				fmbz = fmb
				fsez = 0
				fmez = fme
			case 1: //along local z axis
				fsby = fsb
				fmby = 0
				fsey = fse
				fmey = 0
				fsbz = 0
				fmbz = fmb
				fsez = 0
				fmez = fme
			}
			rez[1] = []float64{fab, fsby - 3.0*fmbz/2.0*l, fsbz + 3.0*fmby/2.0*l, 0, 0, 0, fae, fsey + 3.0*fmbz/2.0*l, fsez - 3.0*fmby/2.0*l, ftb + fte, fmey - fmby/2.0, fmez - fmbz/2.0}
		}
	case 2: //hinge at end
		switch mfrtyp {
		case 0: //beam
			rez[1] = []float64{fsb - 3*fme/(2.0*l), fmb - fme/2.0, fse + 3*fme/(2.0*l), 0}
		case 1: //grid
			rez[1] = []float64{fsb - 3*fme/(2.0*l), ftb + fte, fmb - fme/2.0, fse + 3*fme/(2.0*l), 0, 0}
		case 2: //2d frame
			rez[1] = []float64{fab, fsb - 3.0*fme/(2.0*l), fmb - fme/2.0, fae, fse + 3.0*fme/(2.0*l), 0}
		case 3: //3d frame
			var fsby, fmby, fsey, fmey, fsbz, fmbz, fsez, fmez float64
			ldaxis := ldet[5]
			switch ldaxis {
			case 0: //along local y axis
				fsby = fsb
				fmby = 0
				fsey = fse
				fmey = 0
				fsbz = 0
				fmbz = fmb
				fsez = 0
				fmez = fme
			case 1: //along local z axis
				fsby = fsb
				fmby = 0
				fsey = fse
				fmey = 0
				fsbz = 0
				fmbz = fmb
				fsez = 0
				fmez = fme
			}
			rez[1] = []float64{fab, fsby - 3.0*fmbz/2.0*l, fsbz + 3.0*fmby/2.0*l, 0, 0, 0, fae, fsey + 3.0*fmbz/2.0*l, fsez - 3.0*fmby/2.0*l, ftb + fte, fmey - fmby/2.0, fmez - fmbz/2.0}
		}
	case 3: //hinge at both ends
		switch mfrtyp {
		case 0: //beam
			rez[1] = []float64{fsb - (fmb+fme)/l, 0, fse + (fmb+fme)/l, 0}
		case 1: //grid
			rez[1] = []float64{fsb - (fmb+fme)/l, 0, 0, fse + (fmb+fme)/l, 0, 0}
		case 2: //2d frame
			rez[1] = []float64{fsb - (fmb+fme)/l, 0, 0, fse + (fmb+fme)/l, 0, 0}
		case 3: //3d frame
			var fsby, fmby, fsey, fmey, fsbz, fmbz, fsez, fmez float64
			ldaxis := ldet[5]
			switch ldaxis {
			case 0: //along local y axis
				fsby = fsb
				fmby = 0
				fsey = fse
				fmey = 0
				fsbz = 0
				fmbz = fmb
				fsez = 0
				fmez = fme
			case 1: //along local z axis
				fsby = fsb
				fmby = 0
				fsey = fse
				fmey = 0
				fsbz = 0
				fmbz = fmb
				fsez = 0
				fmez = fme
			}
			rez[1] = []float64{fab, fsby - (fmbz+fmez)/l, fsbz - (fmby+fmey)/l, 0, 0, 0, fae, fsey + (fmbz+fmez)/l, fsez - (fmby+fmey)/l, 0, 0, 0}
		}
	}

	qfchn <- rez
}
