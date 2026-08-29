package barf

import (
	//"fmt"
	"math"
	"log"
	"gonum.org/v1/gonum/mat"
)

//CalcNp rips off section 5.4 of mosley/spencer for plane frame/beam stiffness analysis
//for members of variable cross section (np)

func CalcNp(mod *Model, frmtyp string, cmosh bool) ([]interface{}, error) {
	//make member and joint maps
	//cmosh - use mosley routine
	if !mod.Npsec{mod.Npsec = true}
	var ncjt, mfrtyp int
	switch frmtyp{
		case "1db":
		ncjt = 2; mfrtyp = 0
		mod.Ncjt = 2
		mod.Frmtyp = 0
		case "2df":
		ncjt = 3; mfrtyp = 1
		mod.Ncjt = 3
		mod.Frmtyp = 1
		//wot?
		//wtf these don't work
		case "3dg":
		ncjt = 3; mfrtyp = 2
		case "3df":
		ncjt = 6; mfrtyp = 3
	}
	frmrez := make([]interface{}, 7)
	js := make(map[int]*Node)
	for i, val := range mod.Coords {
		js[i+1]= &Node{Id:i+1,Coords:val}
	}
	ms := make(map[int]*MemNp)
	//introduce dmin (min depth) for ease of plotting
	mod.Dmin = 1e6
	for i, val := range mod.Mprp{
		var idx, tdx, styp int
		idx = val[3]-1
		if mod.Sectyp > 0{
			styp = mod.Sectyp
		} else {			
			switch len(mod.Sts){
				case 0:
				styp = 1
				default:
				if len(mod.Sts) >= idx{
					styp = mod.Sts[idx]
				} else {
					styp = 1
				}
			}	
		}
		var yb, ye float64
		if mod.Ncjt > 2{
			yb = mod.Coords[val[0]-1][1]; ye = mod.Coords[val[1]-1][1]
		}
		if len(mod.Ts) == 1{
			tdx = 0
		} else {
			tdx = val[5]-1
		}
		if mod.Ds[idx][1] < mod.Dmin && yb == ye{
			mod.Dmin = mod.Ds[idx][1]
		}
		ms[i+1] = &MemNp{Id:i+1, Mprp:val, Ts:mod.Ts[tdx], Ds:mod.Ds[idx], Bs:mod.Bs[idx], Ls:mod.Ls[idx], Dims:mod.Dims[idx], Em: mod.Em[val[2]-1][0], Vp:mod.Em[val[2]-1][1], Frmtyp:frmtyp, Styp:styp}
		jb, je := val[0], val[1]
		js[jb].Mems = append(js[jb].Mems, i+1)
		js[je].Mems = append(js[je].Mems, i+1)
		if cmosh{
			err := KFacMos(ms[i+1], mod.VDx)
			if err != nil{
				log.Println(err)
				return frmrez, err
			}
		}else{
			err := KFacToz(ms[i+1], mod.VDx)
			if err != nil{
				log.Println(err)
				return frmrez, err
			}
		}
	}
	//calc nr, ndof and number structure coordinates
	nj := len(mod.Coords)
	nr := 0
	for _, val := range mod.Supports {
		for _, i := range val[1:] {
			js[val[0]].Supports = val[1:]
			if i == -1{
				nr++
			}
		}
	}
	ndof := ncjt*nj - nr
	//log.Println("NDOF->",ndof)
	nsc := make([]int, nj*ncjt)
	ndis := 1
	nres := ndof + 1
	for _, s := range mod.Supports {
		idx := (s[0] - 1) * ncjt
		for i, val := range s[1:] {
			switch val {
			case -1:
				nsc[idx+i] = nres
				nres++
				js[s[0]].Nrs++
			}
		}
	}
	for i, n := range nsc{
		if n == 0{
			nsc[i] = ndis
			ndis++
		}
	}
	
	for idx  := range js{
		nscidx := (idx-1) * ncjt
		js[idx].Nscs = nsc[nscidx:nscidx+ncjt]
	}
	//log.Println("nsc vec->",nsc)

	kchn := make(chan [][]float64, len(ms))
	//individual member stiffness matrices
	for member, m := range ms{
		go KMemNp(member, m, ndof, ncjt, nsc, js, mod.VDx, cmosh, kchn)
	}
	rezchn := make(chan [][]float64, len(ms))
	symchn := make(chan *mat.SymDense)
	switch mfrtyp {
	case 0:
		go KassB(ndof, len(ms), rezchn, symchn)
		for i := 0; i < len(ms); i++ {
			srez := <-kchn
			rezchn <- srez
		}
	case 1:
		go Kass(ndof, len(ms), rezchn, symchn)
		for i := 0; i < len(ms); i++ {
			srez := <-kchn
			rezchn <- srez
		}
	}
	pchn := make(chan []float64, 1)
	pfchn := make(chan []interface{}, 1)
	//compute forces
	//MemFrc returns msqf (map of member end forces), pf (force rez; idx and force),
	//msloaded (map of member and individual force vectors)
	//NodeFrc assembles the nodal force vector p (of p-pf = sd fame)
	go MemFrcNp(ms, mod.Msloads, nsc, ncjt, ndof, mod.VDx, pfchn)
	go NodeFrc(mod.Jloads, nsc, ncjt, ndof, pchn)
	p := <-pchn
	pfrez := <-pfchn
	msqf, _ := pfrez[0].(map[int][]float64)
	pf, _ := pfrez[1].([]float64)
	//log.Println("pvec->",p, "pfvec->",pf)
	msloaded, _ := pfrez[2].(map[int][][]float64)
	//subtract pf from p
	for i, f := range pf{
		p[i] = p[i] - f
	}
	//log.Println("pvec->",p)
	//log.Println(msqf)
	//make matriksees
	pmat := mat.NewDense(ndof, 1, p)
	//NOW receive smatsym from Kass
	smatsym := <-symchn
	//solve by cholesky decomposition of sym banded smatsym
	//log.Println("stiffness matrix->")
	//log.Println("\n",mat.Formatted(smatsym))
	var schmat mat.Cholesky
	var dchmat mat.Dense
	if ok := schmat.Factorize(smatsym); !ok {
		return frmrez, ErrFact
	}
	if err := schmat.SolveTo(&dchmat, pmat); err != nil{
		return frmrez, ErrSolve
	}
	//build global displacement vector
	
	dglb := make([]float64, nj*ncjt)
	for i, sc := range nsc{
		if sc <= ndof {
			dglb[i] = dchmat.At(sc-1, 0)
		} else {
			dglb[i] = 0
		}
	}
	rchn := make(chan [][]float64, len(ms))
	fchn := make(chan []interface{}, len(ms))
	for member, mem := range ms{
		//END FRC NP is da same as endfrc, just mem typ changes
		go EndFrcNp(member, mem, msqf, dglb, nsc, mfrtyp, ncjt, ndof, rchn, fchn)
	}
	rnode := make([]float64, nr)
	for i := 0; i < len(ms); i++{
		rez := <-rchn
		for _, r := range rez{
			if len(r) != 0{
				rnode[int(r[0])] = rnode[int(r[0])] + r[1]
			}
		}
	}
	//log.Println(ColorRed,"RNODE->",rnode,ColorReset)
	for i :=0; i < len(ms); i++ {
		frez := <- fchn
		switch mfrtyp{
			case 0:
			//bheem
			mem, _ := frez[0].(int)
			fvec,_ := frez[1].([]float64)
			msqf[mem] = fvec
			ms[mem].Qf = fvec
			case 1:
			//2d frame
			mem, _ := frez[0].(int)
			fvec,_ := frez[1].([]float64)
			qfvec,_ := frez[2].([]float64)
			msqf[mem] = fvec
			ms[mem].Qf = qfvec
			ms[mem].Gf = fvec
		}
	}
	//HAHA THIS IS SO WRENG
	for node := range js {
		switch mfrtyp{
			case 0:
			yidx := (node -1)*ncjt
			midx := (node -1)*ncjt + 1
			ydis := dglb[yidx]
			mdis := dglb[midx]
			js[node].Displ = []float64{ydis, mdis}
			var ry, mz float64
			if nsc[yidx] > ndof {
				ry = rnode[nsc[yidx]-ndof-1]
			}
			if nsc[midx] > ndof {
				mz =  rnode[nsc[midx]-ndof-1]
			}
			js[node].React = []float64{ry, mz}
			case 1:
			xidx := (node -1)*ncjt
			yidx := (node -1)*ncjt + 1
			zidx := (node -1)*ncjt + 2
			xdis := dglb[xidx]
			ydis := dglb[yidx]
			zrot := dglb[zidx]
			js[node].Displ = []float64{xdis, ydis, zrot}
			var rx, ry, mz float64
			if nsc[xidx] > ndof {
				rx = rnode[nsc[xidx]-ndof-1]
			}
			if nsc[yidx] > ndof {
				ry =  rnode[nsc[yidx]-ndof-1]
			}

			if nsc[zidx] > ndof {
				mz =  rnode[nsc[zidx]-ndof-1]
			}
			js[node].React = []float64{rx, ry, mz}
		}
	}
	frmrez[0] = js
	frmrez[1] = ms
	frmrez[2] = dglb
	frmrez[3] = rnode
	frmrez[4] = msqf
	frmrez[5] = msloaded
	//log.Println("dglb->",dglb)
	switch mfrtyp{
		case 0:
		//bm
		frmrez[6] = BmNpTable(js,ms,dglb,rnode,nsc,ndof,ncjt)
		case 1:
		//frm 2d
		frmrez[6] = F2dNpTable(js,ms,dglb,rnode,nsc,ndof,ncjt)
		case 2:
		case 3:
	}
	//frmrez[6] = Bm1dTable(js,ms,dglb,rnode,nsc,ndof,ncjt)
	return frmrez, nil
}

//KMemNp generates individual member stiffness matrices for members of variable cross section
func KMemNp(member int, mem *MemNp, ndof, ncjt int, nsc []int, js map[int]*Node, vdx,cmosh bool, kchn chan [][]float64){
	var jb, je int
	//memrel := mem.Mprp[4]
	jb = mem.Mprp[0]
	je = mem.Mprp[1]	
	switch mem.Frmtyp {
	case "1db":
		//rez[1] = []float64{mem.Lspan, mem.Em, mem.I0, mem.A0, mem.G}
		bk, _ := KVecNp(mem, ncjt, vdx, cmosh)
		bkmat := mat.NewDense(4, 4, nil)
		for i, row := range bk {
			bkmat.SetRow(i, row)
		}
		mem.Bkmat = bkmat
		//log.Println("\n",mat.Formatted(bkmat))
	case "2df":
		pa := js[jb].Coords
		pb := js[je].Coords
		xa := pa[0]
		ya := pa[1]
		xb := pb[0]
		yb := pb[1]
		l := math.Sqrt(math.Pow(xb-xa, 2) + math.Pow(yb-ya, 2))
		cx := (xb - xa) / l
		cy := (yb - ya) / l
		mem.Cx = cx; mem.Cy = cy
		switch {
		case cx == 0: //if column
			mem.Mtyp = 1
		case cx == 1: //beam
			mem.Mtyp = 2
		default: //then weird thing
			mem.Mtyp = 3 
		}
		//pmax := js[len(js)-1].Coords
		//xmax := pmax[0] //; ymax:= pmax[1]
		// //CHANGE MTYPS TO INTS (done)
		// switch {
		// case cx == 0: //if column
		// 	switch {
		// 	case xa == 0 || xb == 0: //left end column
		// 		mem.Memtyp = append(mem.Memtyp, "CEL,")
		// 	case xb == xmax || xa == xmax: //right end column
		// 		mem.Memtyp = append(mem.Memtyp, "CER,")
		// 	default:
		// 		mem.Memtyp = append(mem.Memtyp, "CM,")
		// 	}
		// default: //then beam
		// 	switch {
		// 	case xa == 0 || xb == 0:
		// 		mem.Memtyp = append(mem.Memtyp, "BEL,")
		// 	case xb == xmax || xa == xmax:
		// 		mem.Memtyp = append(mem.Memtyp, "BER,")
		// 	default:
		// 		mem.Memtyp = append(mem.Memtyp, "BM,")
		// 	}
		// }
		
		bk, t := KVecNp(mem, ncjt, vdx, cmosh)
		bkmat := mat.NewDense(2*ncjt, 2*ncjt, nil)
		for i, row := range bk {
			bkmat.SetRow(i, row)
		}
		mem.Bkmat = bkmat
		//member transformation matrix t
		tmat := mat.NewDense(2*ncjt, 2*ncjt, nil)
		for i, row := range t {
			tmat.SetRow(i, row)
		}
		mem.Tmat = tmat
		//member global stiffness matrix gk = t.T @ k @ t
		gkmat := mat.NewDense(2*ncjt, 2*ncjt, nil)
		gkmat.Product(tmat.T(), bkmat, tmat)
		mem.Gkmat = gkmat
	}
	//assemble structure stiffness matrix 
	srez := [][]float64{}
	var r float64
	var ela, elb, na, nb int
	for i := 1; i <= 2*ncjt; i++ {
		if i <= ncjt {
			ela = (jb-1)*ncjt + i
			na = nsc[(ela - 1)]
		} else {
			ela = (je-1)*ncjt + (i - ncjt)
			na = nsc[(ela - 1)]
		}
		if na <= ndof {
			for j := 1; j <= 2*ncjt; j++ {
				if j <= ncjt {
					elb = (jb-1)*ncjt + j
					nb = nsc[(elb - 1)]
				} else {
					elb = (je-1)*ncjt + (j - ncjt)
					nb = nsc[(elb - 1)]
				}
				if nb <= ndof {
					switch mem.Frmtyp {
					case "1db":
						r = mem.Bkmat.At(i-1, j-1)
					default:
						r = mem.Gkmat.At(i-1, j-1)
					}
					if r != 0.0 {
						srez = append(srez, []float64{float64(na), float64(nb), r})
					}
				}
			}
		}
	}
	if (js[jb].Nrs == 0 && len(js[jb].Mems) == 1) || (js[je].Nrs == 0 && len(js[je].Mems) == 1){
		mem.Clvr = true
	}
	kchn <- srez
}

//KVecNp returns member local stiffness and transformation arrays for an np member
func KVecNp(mem *MemNp, ncjt int, vdx, cmosh bool) (bk, t [][]float64){
	//b0, t1, t2, t3 := 1.0, 1.0, 1.0, 1.0
	memrel := mem.Mprp[4]
	//if vdx{b0 = mem.B0; t1 = mem.T1; t2 = mem.T2; t3 = mem.T3}
	l := mem.Lspan
	//zk := mem.Em * mem.I0///math.Pow(mem.Lspan,3)
	bk = make([][]float64, 2 * ncjt)
	for i := 0; i < 2 * ncjt; i++{
		bk[i] = make([]float64, 2 * ncjt)
	}
	ei :=  mem.Em * mem.I0/l
	b := mem.Kc * math.Pow(l,2) * ei 
	ai := mem.Ka * math.Pow(l,2) * ei
	aj := mem.Kb * math.Pow(l,2) * ei
	ci := (ai + b)/l; cj := (aj + b)/l
	d := (ci + cj)/l
	switch mem.Frmtyp {
	case "1db":
		switch memrel {
		case 0:
		                    	
			bk[0][0] = d
			bk[0][1] = ci
			bk[0][2] = -d
			bk[0][3] = cj
			
			bk[1][0] = ci
			bk[1][1] = ai
			bk[1][2] = -ci
			bk[1][3] = b
			
			bk[2][0] = -d
			bk[2][1] = -ci
			bk[2][2] = d
			bk[2][3] = -cj
			
			bk[3][0] = cj
			bk[3][1] = b
			bk[3][2] = -cj
			bk[3][3] = aj
		}
	case "2df":
		//za := mem.N11 * mem.Em/l
		//from ozay - axial force term 
		za := mem.N11 * mem.Em * mem.A0/l
		t = make([][]float64, 2 * ncjt)
		for i := 0; i < 2 * ncjt; i++{
			t[i] = make([]float64, 2 * ncjt)
		}
		t[0][0] = mem.Cx
		t[0][1] = mem.Cy
		t[1][0] = -mem.Cy
		t[1][1] = mem.Cx
		t[2][2] = 1.0
		t[3][3] = mem.Cx
		t[3][4] = mem.Cy
		t[4][3] = -mem.Cy
		t[4][4] = mem.Cx
		t[5][5] = 1.0
		switch memrel {
		case 0:
			bk[0][0] = za
			bk[0][3] = -za

			bk[3][0] = -za
			bk[3][3] = za

			bk[1][1] = d 
			bk[1][2] = ci
			bk[1][4] = -d
			bk[1][5] = ci
			
			bk[2][1] = ci
			bk[2][2] = ai
			bk[2][4] = -ci
			bk[2][5] = b
			
			bk[4][1] = -d
			bk[4][2] = -ci
			bk[4][4] = d
			bk[4][5] = -cj
			
			bk[5][1] = cj
			bk[5][2] = b
			bk[5][4] = -cj
			bk[5][5] = aj

		}
	}
	return
}

//KFacMos gets stiffness and carry over factors for an np member
//see mosley-spencer stiffness coeff calcs, section 4.2
func KFacMos(mem *MemNp, vdx bool) (err error){
	div := (mem.Ls[0]+ mem.Ls[1] + mem.Ls[2])/20.0
	if div == 0.0{return ErrDim}
	mem.Xs = make([]float64, 21)
	for i := 0; i < 21; i++ {
		mem.Xs[i] = div * float64(i)
	}
	mem.Ix = make([]float64, 21)
	mem.Ax = make([]float64, 21)
	//make plots easier
	mem.Bxs = make([]float64, 21)
	mem.Dxs = make([]float64, 21)
	//calc stiffness and carry over factors
	if len(mem.Ls) != len(mem.Ts){return ErrDim}
	var bx, dx float64
	for i, x := range mem.Xs {
		switch {
		case x < 1.001 * mem.Ls[0]:
			switch mem.Ts[0] {
			case 0:
				//nuthin
				if x < mem.Ls[0] {
					dx = mem.Ds[0]
					bx = mem.Bs[0]
				} else {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				}
			case 1:
				//prismatic (flat)
				if x < 0.999 * mem.Ls[0] {
					dx = mem.Ds[0]
					bx = mem.Bs[0]
				} else {
					//???HYARR???
					dx = (mem.Ds[0] + mem.Ds[1])/2.0
					bx = (mem.Bs[0] + mem.Bs[1])/2.0 //?????
					//dx = mem.Ds[1]
					//bx = mem.Bs[1]
				}
				
			case 2:
				//straight, 2f haunch, 3f haunch
				if x < mem.Ls[0] {
					dx = mem.Ds[1] + (mem.Ls[0] - x) * (mem.Ds[0] - mem.Ds[1])/mem.Ls[0]
					bx = mem.Bs[0]
				} else {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				}
			case 3:
				//parabolic
				if x < mem.Ls[0] {
					dx = mem.Ds[1] + math.Pow(mem.Ls[0] - x,2) * (mem.Ds[0] - mem.Ds[1])/math.Pow(mem.Ls[0],2)
					bx = mem.Bs[0]
				} else {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				}
			}
		case x >= mem.Ls[0] + mem.Ls[1]:
			switch mem.Ts[2] {
			case 0:
				//nuthin
				if x < mem.Ls[0] + mem.Ls[1] {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				} else {
					dx = mem.Ds[2]
					bx = mem.Bs[2]
				}	
			case 1:
				//prismatic
				if x < mem.Ls[0] + mem.Ls[1] * 1.001 {
					dx = (mem.Ds[1] + mem.Ds[2])/2.0
					bx = (mem.Bs[1] + mem.Bs[2])/2.0
				} else {
					dx = mem.Ds[2]
					bx = mem.Bs[2]
				}
			case 2:
				//straight
				if x > mem.Ls[0] + mem.Ls[1] {
					dx = mem.Ds[1] + (x - mem.Ls[0] - mem.Ls[1]) * (mem.Ds[2] - mem.Ds[1])/mem.Ls[2]
					bx = mem.Bs[2]
				} else {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				}
			case 3:
				//parabolic
				if x > mem.Ls[0] + mem.Ls[1] {
					dx = mem.Ds[1] + math.Pow(x - mem.Ls[0] - mem.Ls[1], 2) * (mem.Ds[2] - mem.Ds[1])/math.Pow(mem.Ls[2],2)
					bx = mem.Bs[2]
				} else {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				}
			}
		default:
			//always nuthing at center for this kind of beam
			dx = mem.Ds[1]
			bx = mem.Bs[1]
		}
		ar, ix, _ := PropNpBm(mem.Styp, bx, dx, mem.Dims)
		mem.Ix[i] = ix
		mem.Ax[i] = ar
		mem.Bxs[i] = bx
		mem.Dxs[i] = dx
	}
	a0, i0, _ := PropNpBm(mem.Styp, mem.Dims[0], mem.Dims[1], mem.Dims)
	mem.A0 = a0; mem.I0 = i0
	//simpson's ordinates
	homers := make([]float64, 21)
	for i := range homers {
		switch {
		case i == 0:
			homers[i] = 1
		case i == 20:
			homers[i] = 1
		case i % 2 == 0:
			homers[i] = 2
		case i % 2 != 0:
			homers[i] = 4
		}
	}
	//moment diagram ordinates
	lspan := mem.Ls[0] + mem.Ls[1] + mem.Ls[2]
	var m1, m2, m3, a, f1, f2, f3 float64
	mem.M4 = make([]float64, 21); mem.M5 = make([]float64, 21)
	for i, x := range mem.Xs {
		y := float64(i)/20.0
		m1 = math.Pow(lspan - x, 2)
		m2 = math.Pow(x,2)
		m3 = x * (lspan - x)
		//m1 = math.Pow(1.0-y,2)
		//m2 = math.Pow(y,2)
		//m3 = y * (1.0 - y)
		a = homers[i] * i0/ mem.Ix[i]/(3.0 * 20.0)
		//b := homers[i]/mem.Ix[i]/mem.Em/mem.Lspan/(60.0)
		f1 += m1 * a
		f2 += m2 * a
		f3 += m3 * a
		mem.M4[i] = a * (1.0 - y) 
		mem.M5[i] = a * y 
		mem.I1 += homers[i] * math.Pow(x, 2)/mem.Ix[i]/(3.0 * 20.0)
		mem.I2 += homers[i]/mem.Ix[i]/(3.0 * 20.0)
		mem.I3 += homers[i] * x /mem.Ix[i]/(3.0 * 20.0)
		mem.I4 += homers[i]/mem.Ax[i]/(3.0 * 20.0)
	}
	ka := f2/(f1 * f2 - math.Pow(f3,2))
	kb := f1/(f1 * f2 - math.Pow(f3,2))
	kc := ka * f3/f2
	ca := f3/f2
	cb := f3/f1
	mem.Ka = ka; mem.Kb = kb; mem.Kc = kc; mem.Ca = ca; mem.Cb = cb; mem.Lspan = lspan; mem.A0 = a0; mem.I0 = i0
	detc := -mem.I0 * (mem.I1 * mem.I2 - math.Pow(mem.I3,2))/lspan
	if detc == 0{return ErrDim}
	mem.M11 = -mem.I1/detc
	mem.M22 = -(lspan * mem.I3 - mem.I1)/detc
	mem.M12 = (2.0 * lspan * mem.I3 - math.Pow(lspan, 2) * mem.I2 - mem.I1)/detc
	mem.N11 = lspan/a0/mem.I4
	switch mem.Styp{
		case 0:
		//round bar
		mem.Fs = 10.0/9.0
		case 1,2,3:
		//rect, tri
		mem.Fs = 6.0/5.0
		case 4:
		//box
		b, d, B, D := mem.Dims[0],mem.Dims[1],mem.Dims[2],mem.Dims[3]
		aweb := (B- b)*(D-d)/2.0
		mem.Fs = (B*D - b*d)/aweb
		case 5:
		//tube
		mem.Fs = 2.0
		case 9:
		//i section
		b, d, tf, tw := mem.Dims[0],mem.Dims[1],mem.Dims[2],mem.Dims[3]
		aweb := (d - 2.0 * tf) * tw
		mem.Fs = (2.0 * b * tf + (d - 2.0 * tf) * tw)/aweb
	}
	//log.Println("member->",mem.Id,"\nka, kb, kc->",ka*lspan*lspan, kb*lspan*lspan, kc*lspan*lspan)
	//log.Println("mZZ->", mem.M11/lspan, mem.M22/lspan, mem.M12/lspan)
	if vdx{
		mem.G = mem.Em / (2.0 * mem.Vp)
		g := mem.G
		mem.B0 = g * math.Pow(lspan, 3)/((mem.M11 + mem.M22 + 2.0 * mem.M12) * mem.Fs * mem.Em * mem.I0 * mem.I4 + g * math.Pow(lspan, 3))
		mem.T1 = mem.Fs*mem.Em*mem.I0*mem.I4 * (mem.M11 * mem.M22 - math.Pow(mem.M12,2))/(mem.M11 * g * math.Pow(lspan,3))
		mem.T1 = (mem.T1 + 1.0) * mem.B0
		mem.T2 = -mem.Fs*mem.Em*mem.I0*mem.I4 * (mem.M11 * mem.M22 - math.Pow(mem.M12,2))/(mem.M12 * g * math.Pow(lspan,3))
		mem.T2 = (mem.T2 + 1.0) * mem.B0
		mem.T3 = mem.Fs*mem.Em*mem.I0*mem.I4 * (mem.M11 * mem.M22 - math.Pow(mem.M12,2))/(mem.M22 * g * math.Pow(lspan,3))
		mem.T3 = (mem.T3 + 1.0) * mem.B0
	}
	return
}

//KFacToz calculates Topcu/Ozay factors for a member of variable cross section
func KFacToz(mem *MemNp, vdx bool) (error){
	//calcs topcu-ozay factors
	div := (mem.Ls[0]+ mem.Ls[1] + mem.Ls[2])/20.0
	if div == 0.0{return ErrDim}
	mem.Xs = make([]float64, 21)
	for i := 0; i < 21; i++ {
		mem.Xs[i] = div * float64(i)
	}
	mem.Ix = make([]float64, 21)
	mem.Ax = make([]float64, 21)
	
	mem.Bxs = make([]float64, 21)
	mem.Dxs = make([]float64, 21)
	if len(mem.Ls) != len(mem.Ts){return ErrDim}
	//calc stiffness and carry over factors
	var bx, dx float64
	for i, x := range mem.Xs {
		switch {
		case x < 1.001 * mem.Ls[0]:
			switch mem.Ts[0] {
			case 0:
				//nuthin
				if x < mem.Ls[0] {
					dx = mem.Ds[0]
					bx = mem.Bs[0]
				} else {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				}
			case 1:
				//prismatic
				if x < 0.999 * mem.Ls[0] {
					dx = mem.Ds[0]
					bx = mem.Bs[0]
				} else {
					dx = (mem.Ds[0] + mem.Ds[1])/2.0
					bx = (mem.Bs[0] + mem.Bs[1])/2.0 //?????
				}
				
			case 2:
				//straight
				if x < mem.Ls[0] {
					dx = mem.Ds[1] + (mem.Ls[0] - x) * (mem.Ds[0] - mem.Ds[1])/mem.Ls[0]
					bx = mem.Bs[0]
				} else {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				}
			case 3:
				//parabolic
				if x < mem.Ls[0] {
					dx = mem.Ds[1] + math.Pow(mem.Ls[0] - x,2) * (mem.Ds[0] - mem.Ds[1])/math.Pow(mem.Ls[0],2)
					bx = mem.Bs[0]
				} else {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				}
			}
		case x >= mem.Ls[0] + mem.Ls[1]:
			switch mem.Ts[2] {
			case 0:
				//nuthin
				if x < mem.Ls[0] + mem.Ls[1] {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				} else {
					dx = mem.Ds[2]
					bx = mem.Bs[2]
				}	
			case 1:
				//prismatic
				if x < mem.Ls[0] + mem.Ls[1] * 1.001 {
					dx = (mem.Ds[1] + mem.Ds[2])/2.0
					bx = (mem.Bs[1] + mem.Bs[2])/2.0
				} else {
					dx = mem.Ds[2]
					bx = mem.Bs[2]
				}
			case 2:
				//straight
				if x > mem.Ls[0] + mem.Ls[1] {
					dx = mem.Ds[1] + (x - mem.Ls[0] - mem.Ls[1]) * (mem.Ds[2] - mem.Ds[1])/mem.Ls[2]
					bx = mem.Bs[2]
				} else {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				}
			case 3:
				//parabolic
				if x > mem.Ls[0] + mem.Ls[1] {
					dx = mem.Ds[1] + math.Pow(x - mem.Ls[0] - mem.Ls[1], 2) * (mem.Ds[2] - mem.Ds[1])/math.Pow(mem.Ls[2],2)
					bx = mem.Bs[2]
				} else {
					dx = mem.Ds[1]
					bx = mem.Bs[1]
				}
			}
		default:
			//always nuthing at center for this kind of beam
			dx = mem.Ds[1]
			bx = mem.Bs[1]
		}
		ar, ix, _ := PropNpBm(mem.Styp, bx, dx, mem.Dims)
		mem.Ix[i] = ix
		mem.Ax[i] = ar

		mem.Bxs[i] = bx
		mem.Dxs[i] = dx
	}
	a0, i0, _ := PropNpBm(mem.Styp, mem.Bs[1], mem.Ds[1], mem.Dims)
	mem.I0 = i0; mem.A0 = a0
	//simpson's ordinates
	homers := make([]float64, 21)
	for i := range homers {
		switch {
		case i == 0:
			homers[i] = 1
		case i == 20:
			homers[i] = 1
		case i % 2 == 0:
			homers[i] = 2
		case i % 2 != 0:
			homers[i] = 4
		}
	}
	//moment diagram ordinates
	lspan := mem.Ls[0] + mem.Ls[1] + mem.Ls[2]
	var a , m1, m2, m3, f1, f2, f3 float64
	mem.M4, mem.M5 = make([]float64, 21), make([]float64, 21)
	for i, x := range mem.Xs {
		mem.I1 += homers[i] * math.Pow(x, 2)/mem.Ix[i]/(3.0 * 20.0)
		mem.I2 += homers[i]/mem.Ix[i]/(3.0 * 20.0)
		mem.I3 += homers[i] * x /mem.Ix[i]/(3.0 * 20.0)
		mem.I4 += homers[i]/mem.Ax[i]/(3.0 * 20.0)
		m1 = math.Pow(lspan - x, 2)
		m2 = math.Pow(x,2)
		m3 = x * (lspan - x)
		a = homers[i] * mem.I0/ mem.Ix[i]/ (3.0 * 20.0)
		f1 += m1 * a
		f2 += m2 * a
		f3 += m3 * a
		mem.M4[i] = homers[i] * mem.I0 * (lspan - x)/mem.Ix[i]/math.Pow(lspan,2)
		mem.M5[i] = homers[i] * mem.I0 * x /mem.Ix[i]/math.Pow(lspan,2)
	}
	detc := -mem.I0 * (mem.I1 * mem.I2 - math.Pow(mem.I3,2))/lspan
	if detc == 0{return ErrDim}
	mem.M11 = -mem.I1/detc
	mem.M22 = -(lspan * mem.I3 - mem.I1)/detc
	mem.M12 = (2.0 * lspan * mem.I3 - math.Pow(lspan, 2) * mem.I2 - mem.I1)/detc
	mem.N11 = lspan/a0/mem.I4
	mem.Ka = f2/(f1 * f2 - math.Pow(f3,2))
	mem.Kb = f1/(f1 * f2 - math.Pow(f3,2))
	mem.Kc = mem.Ka * f3/f2
	mem.Ca = f3/f2
	mem.Cb = f3/f1
	mem.Lspan = lspan
	switch mem.Styp{
		case 0:
		//round bar
		mem.Fs = 10.0/9.0
		case 1,2,3:
		//rect, tri
		mem.Fs = 6.0/5.0
		case 4:
		//box
		b, d, B, D := mem.Dims[0],mem.Dims[1],mem.Dims[2],mem.Dims[3]
		aweb := (B- b)*(D-d)/2.0
		mem.Fs = (B*D - b*d)/aweb
		case 5:
		//tube
		mem.Fs = 2.0
		case 9:
		//i section
		b, d, tf, tw := mem.Dims[0],mem.Dims[1],mem.Dims[2],mem.Dims[3]
		aweb := (d - 2.0 * tf) * tw
		mem.Fs = (2.0 * b * tf + (d - 2.0 * tf) * tw)/aweb
	}
	if vdx{
		g := mem.Em / (2 * mem.Vp)
		mem.B0 = g * math.Pow(lspan, 3)/((mem.M11 + mem.M22 + 2.0 * mem.M12) * mem.Fs * mem.Em * mem.I0 * mem.I4 + g * math.Pow(lspan, 3))
		mem.T1 = mem.Fs*mem.Em*mem.I0*mem.I4 * (mem.M11 * mem.M22 - math.Pow(mem.M12,2))/(mem.M11 * g * math.Pow(lspan,3))
		mem.T1 = (mem.T1 + 1.0) * mem.B0
		mem.T2 = -mem.Fs*mem.Em*mem.I0*mem.I4 * (mem.M11 * mem.M22 - math.Pow(mem.M12,2))/(mem.M12 * g * math.Pow(lspan,3))
		mem.T2 = (mem.T2 + 1.0) * mem.B0
		mem.T3 = mem.Fs*mem.Em*mem.I0*mem.I4 * (mem.M11 * mem.M22 - math.Pow(mem.M12,2))/(mem.M22 * g * math.Pow(lspan,3))
		mem.T3 = (mem.T3 + 1.0) * mem.B0
		mem.G = g
	}
	//log.Println("kZZZ->",mem.Ka, mem.Kb, mem.Kc, mem.M11, mem.M22, mem.M12)
	return nil
}

/*
   topcu ozay stiffness matrix
   IT IS NOT SYMMETRYS 
   1db
				bk[0][0] = (mem.M11 + mem.M22 + 2.0 * mem.M12) * zk * b0
				bk[0][1] = (mem.M11 + + mem.M12) * zk * l * b0
				bk[0][2] = -(mem.M11 + mem.M22 + 2.0 * mem.M12) * zk * b0
				bk[0][3] = (mem.M22 + + mem.M12) * zk * l * b0
				
				bk[1][0] = (mem.M11 + + mem.M12) * zk * l * b0
				bk[1][1] = mem.M11 * zk * math.Pow(l,2) * t1
				bk[1][2] = -(mem.M11 + mem.M12) * zk * l * b0
				bk[1][3] = mem.M12 * zk * math.Pow(l,2) * t2
				
				bk[2][0] = -(mem.M11 + mem.M22 + 2.0 * mem.M12) * zk * b0
				bk[2][1] = -(mem.M11 + + mem.M12) * zk * l * b0
				bk[2][2] = (mem.M11 + mem.M22 + 2.0 * mem.M12) * zk * b0
				bk[2][3] = -(mem.M22 + + mem.M12) * zk * l * b0
				
				bk[3][0] = (mem.M22 + mem.M12) * zk * l * b0
				bk[3][1] = mem.M12 * zk * math.Pow(l, 2) * t2
				bk[3][2] = -(mem.M22 + mem.M12) * zk * l * b0
				bk[3][3] = mem.M22 * zk * math.Pow(l, 2) * t3

   2df
				bk[0][0] = mem.N11 * za
				bk[0][3] = - mem.N11 * za

				bk[1][1] = (mem.M11 + mem.M22 + 2.0 * mem.M12) * zk * b0
				bk[1][2] = (mem.M11 + + mem.M12) * zk * l * b0
				bk[1][4] = -(mem.M11 + mem.M22 + 2.0 * mem.M12) * zk * b0
				bk[1][5] = (mem.M22 + + mem.M12) * zk * l * b0
				
				bk[2][1] = (mem.M11 + + mem.M12) * zk * l * b0
				bk[2][2] = mem.M11 * zk * math.Pow(l,2) * t1
				bk[2][4] = -(mem.M11 + mem.M12) * zk * l * b0
				bk[2][5] = mem.M12 * zk * math.Pow(l,2) * t2

				bk[3][0] = - mem.N11 * za
				bk[3][3] =  mem.N11 * za
				
				bk[4][1] = -(mem.M11 + mem.M22 + 2.0 * mem.M12) * zk * b0
				bk[4][2] = -(mem.M11 + + mem.M12) * zk * l * b0
				bk[4][4] = (mem.M11 + mem.M22 + 2.0 * mem.M12) * zk * b0
				bk[4][5] = -(mem.M22 + + mem.M12) * zk * l * b0
				
				bk[5][1] = (mem.M22 + mem.M12) * zk * l * b0
				bk[5][2] = mem.M12 * zk * math.Pow(l, 2) * t2
				bk[5][4] = -(mem.M22 + mem.M12) * zk * l * b0
				bk[5][5] = mem.M22 * zk * math.Pow(l, 2) * t3
*/
