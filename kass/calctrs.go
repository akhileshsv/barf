package barf

import (
	"fmt"
	"log"
	"math"
	"errors"
	"gonum.org/v1/gonum/mat"
)

//CalcTrs performs direct stiffness analysis of 2d and 3d truss models
//see kassimali chapter 3/4 and chapter 8 
func CalcTrs(mod *Model, ncjt int) ([]interface{},error){
	frmrez := make([]interface{}, 10)
	var err error
	ms := make(map[int]*Mem)
	
	for i, val := range mod.Mprp {
		switch ncjt{
			case 2:	
			if len(val) < 4{
				return frmrez, fmt.Errorf("invalid length of mat prop. slice for member %v",i+1)
			}
			jb := val[0]; je := val[1]; em := val[2]; cp := val[3]
			switch{
				case jb < 1 || jb > len(mod.Coords)||len(mod.Coords[jb-1])<2:
				err = fmt.Errorf("invalid node index for member %v -> %v",i+1, jb)
				case je < 1 || je > len(mod.Coords)||len(mod.Coords[jb-1])<2:
				err = fmt.Errorf("invalid node index for member %v -> %v",i+1, je)
				case em < 1 || len(mod.Em) < em || len(mod.Em[em-1]) < 1:
				err = fmt.Errorf("invalid em index for member %v -> %v",i+1, em)
				case cp < 1 || len(mod.Cp) < cp || len(mod.Cp[cp-1]) < 1:
				err = fmt.Errorf("invalid cp index for member %v -> %v",i+1, cp)
			}
			case 3:	
			if len(val) < 4{
				return frmrez, fmt.Errorf("invalid length of mat prop. slice for member %v",i+1)
			}
			jb := val[0]; je := val[1]; em := val[2]; cp := val[3]
			switch{
				case jb < 1 || jb > len(mod.Coords) || len(mod.Coords[jb-1])<3:
				err = fmt.Errorf("invalid node index for member %v -> %v",i+1, jb)
				case je < 1 || je > len(mod.Coords) || len(mod.Coords[je-1])<3:
				err = fmt.Errorf("invalid node index for member %v -> %v",i+1, je)
				case em < 1 || len(mod.Em) < em || len(mod.Em[em-1]) < 1:
				err = fmt.Errorf("invalid em index for member %v -> %v",i+1, em)
				case cp < 1 || len(mod.Cp) < cp || len(mod.Cp[cp-1]) < 1:
				err = fmt.Errorf("invalid cp index for member %v -> %v",i+1, cp)
			}
		}	
		if err != nil{
			return frmrez, err
		}
		ms[i+1] = &Mem{Id:i+1,Mprp:val}
	}
	js := make(map[int]*Node)
	for i, val := range mod.Coords {
		js[i+1]= &Node{Id:i+1,Coords:val}

	}
	//calc nr, ndof and number structure coordinates
	nj := len(mod.Coords)
	nr := 0
	for _, val := range mod.Supports {
		for _, i := range val[1:] {
			js[val[0]].Supports = val[1:]
			if i == -1 {
				nr++
			}
		}
	}
	ndof := ncjt*nj - nr
	nsc := make([]int, nj*ncjt)
	ndis := 1
	nres := ndof + 1
	if ndof <= 0{
		err = fmt.Errorf("invalid degrees of freedom (%v) for model",ndof)
		return frmrez, err
	}
	for _, s := range mod.Supports{
		switch ncjt{
			case 2:	
			if len(s) < 3{
				err = fmt.Errorf("invalid support slice (%v) for model",s)
				return frmrez, err
			}
			
			case 3:
			if len(s) < 4{
				err = fmt.Errorf("invalid support slice (%v) for model",s)
				return frmrez, err
			}
		}
		idx := (s[0] - 1) * ncjt
		for i, val := range s[1:] {
			switch val {
			case -1:
				nsc[idx+i] = nres
				nres++
			}
		}
	}
	for i, val := range mod.Jloads{
		switch ncjt{
			case 2:	
			if len(val) < 3{
				return frmrez, fmt.Errorf("invalid nodal force slice #(%v)-(%v) for model",i+1, val)
			}
			case 3:
			if len(val) < 4{
				return frmrez, fmt.Errorf("invalid nodal force slice #(%v)-(%v) for model",i+1, val)
			}
		}
	}
	for i, n := range nsc {
		if n == 0 {
			nsc[i] = ndis
			ndis++
		}
	}
	
	//mark nodes with nscs
	for idx := range js {
		nscidx := (idx-1) * ncjt
		js[idx].Nscs = nsc[nscidx:nscidx+ncjt]
		js[idx].Sdj = make([]float64, ncjt)
	}
	//mark nodes with support displacements
	for i, idx := range mod.Jsd{
		sdj := mod.Sdj[i]
		if len(sdj) == ncjt{
			copy(js[idx].Sdj, sdj)
			js[idx].Sd = true
		} else {
			log.Println("ERRORE,errore->error reading support displacement for node->",idx)
		}
	}
	//get member geometry, stiffness matrices
	kchn := make(chan []interface{}, len(ms))
	for member, m := range ms{
		jb := m.Mprp[0]; je := m.Mprp[1]
		bcon := (jb < 0 || je < 0) || (jb > len(mod.Coords) || je > len(mod.Coords))
		if bcon{
			err := fmt.Errorf("invalid node index for member %v",member)
			return frmrez, err
		}
		js[jb].Mems = append(js[jb].Mems,member)
		js[je].Mems = append(js[je].Mems,member)
		if ncjt == 3 {
			go Kmem3dT(member, m.Mprp, nsc, js, mod.Em, mod.Cp, ndof, ncjt, kchn)
		}
		if ncjt == 2 {
			go Kmem2dT(member, m.Mprp, nsc, js, mod.Em, mod.Cp, ndof, ncjt, kchn)
		}
	}
	//call matrix assembly func (it waits, vacantly)
	//rezchan to receive rez, symchn for symmetric banded
	rezchn := make(chan [][]float64, len(ms))
	symchn := make(chan *mat.SymDense)
	go KassT(ndof, len(ms), rezchn, symchn) //removing schn
	
	//fixed force joint vector due to support displacements
	pffs := make([]float64, ndof)
	//get member results from KmemT and update dict
	for i := 0; i < len(ms); i++ {
		//0 - member, 1 - cx , cy, l, e, a, iz, 2 - bkmat
		//3 - tmat, 4 - gkmat, 5 - srez (stiffness matrix non zero elements),
		//6 - member type (End/Mid Bm/Col)
		x := <-kchn
		srez, _ := x[5].([][]float64)
		rezchn <- srez
		member, _ := x[0].(int)
		geos, _ := x[1].([]float64)
		ms[member].Geoms = geos
		bkmat, _ := x[2].(*mat.Dense)
		ms[member].Bkmat = bkmat
		tmat, _ := x[3].(*mat.Dense)
		ms[member].Tmat = tmat
		gkmat, _ := x[4].(*mat.Dense)
		ms[member].Gkmat = gkmat
		vfs, _ := x[6].([]float64)
		ms[member].Vfs = vfs
		pfsrez, _ := x[7].([][]float64)
		if len(pfsrez) != 0{
			for _, v := range pfsrez{
				j := int(v[0]); val := v[1]
				pffs[j-1] += val
			}
		}
		mtyp, _ := x[8].(int)
		ms[member].Mtyp = mtyp
	}
	//receive smatsym (and smat) from Kass
	//smat := <-schn
	smatsym := <-symchn
	pchn := make(chan []float64, 1)
	//compute forces
	go NodeFrc(mod.Jloads, nsc, ncjt, ndof, pchn)

	//NEW
	mfrtyp := -2
	pfchn := make(chan []interface{}, 1)
	switch ncjt{
		case 2:
		case 3:
		mfrtyp = -1
	}
	go MemFrc(ms, mod.Msloads, nsc, ncjt, ndof, mfrtyp, pfchn)
	p := <-pchn
	pfrez := <-pfchn
	
	msqf, _ := pfrez[0].(map[int][]float64)
	pf, _ := pfrez[1].([]float64)
	//log.Println(ColorCyan,"pffs",pffs,ColorReset)
	//log.Println(ColorCyan,"p",p,ColorReset)
	//log.Println(ColorCyan,"pf",pf,ColorReset)
	//subtract pffs from p
	for i := range p{
		p[i] = p[i] - pffs[i] -pf[i]
	}
	//make matriksees
	pmat := mat.NewDense(ndof, 1, p)
	
	// dmat := mat.NewDense(ndof, 1, nil)
	// err := dmat.Solve(smat, pmat)
	
	// if err != nil {
	// 	log.Println(ColorRed, "MATRIX SOLUTION ERROR\n", err)
	
	// 	return frmrez, errors.New("matrix solution error")
	// }
	//solve by cholesky decomposition of sym banded smatsym
	var schmat mat.Cholesky
	var dchmat mat.Dense
	if ok := schmat.Factorize(smatsym); !ok {
		log.Println(ColorRed, "STIFFN. MATRIX IS NOT A +VE SEMI_DEFINITE MATRIX (check supports, then members?)")
		return frmrez, errors.New("non +ve definite stiffness matrix")
	}
	if err := schmat.SolveTo(&dchmat, pmat); err != nil {
		log.Println(ColorRed, "NEAR-SINGULAR STIFFN. MATRIX (check supports?)", err)
		return frmrez, errors.New("near singular stiffness matrix")
	}

	//check if sd == p

	//pchk := mat.NewDense(ndof, 1, nil)
	//pchk.Product(smat, dmat)

	pchkchol := mat.NewDense(ndof, 1, nil)
	pchkchol.Product(&schmat, &dchmat)

	//check mat mul
	for i := 0; i < ndof; i++ {
		if math.Abs(pchkchol.At(i, 0)-pmat.At(i, 0)) >= 1e-5 {
			return frmrez, errors.New("matrix multiplication check error")
		} else {
			continue
		}
	}
	//}

	//build global displacement vector
	//d0 nl analysis displacement vector
	var d0 []float64
	dglb := make([]float64, nj*ncjt)
	for i, sc := range nsc { 
		if sc <= ndof {
			dglb[i] = dchmat.At(sc-1, 0)
			d0 = append(d0, dchmat.At(sc-1, 0))
		} else {
			dglb[i] = 0
		}
	}

	
	rchn := make(chan [][]float64, len(ms))
	
	//now endfrc is defo different
	fchn := make(chan []interface{}, len(ms))

	for member, memprp := range ms {
		go EndFrcT(member, memprp, msqf, dglb, nsc, ncjt, ndof, rchn, fchn)
	}
	rnode := make([]float64, nr)
	for i := 0; i < len(ms); i++ {
		rez := <-rchn
		for _, r := range rez {
			if len(r) != 0 {
				rnode[int(r[0])] = rnode[int(r[0])] + r[1]
			}
		}
	}
	//msqf := make(map[int][][]float64, len(ms))
	for i :=0; i < len(ms); i++ {
		frez := <- fchn
		mem, _ := frez[0].(int)
		qfvec,_ := frez[1].([]float64)
		fvec,_ := frez[2].([]float64)
		ms[mem].Qf = qfvec
		ms[mem].Gf = fvec
	}

	switch {
	case ncjt == 2:
		for node := range js {
			xidx := (node -1)*ncjt
			yidx := (node -1)*ncjt + 1
			xdis := dglb[xidx]
			ydis := dglb[yidx]
			js[node].Displ = []float64{xdis, ydis}
			var rx, ry float64
			if nsc[xidx] > ndof {
				rx = rnode[nsc[xidx]-ndof-1]
			}
			if nsc[yidx] > ndof {
				ry =  rnode[nsc[yidx]-ndof-1]
			}
			js[node].React = []float64{rx, ry}
		}
	case ncjt == 3:
		for node := range js {
			xidx := (node -1)*ncjt
			yidx := (node -1)*ncjt + 1
			zidx := (node -1)*ncjt + 2
			xdis := dglb[xidx]
			ydis := dglb[yidx]
			zdis := dglb[zidx]
			js[node].Displ = []float64{xdis, ydis, zdis}
			var rx, ry, rz float64
			if nsc[xidx] > ndof {
				rx = rnode[nsc[xidx]-ndof-1]
			}
			if nsc[yidx] > ndof {
				ry =  rnode[nsc[yidx]-ndof-1]
			}
			if nsc[zidx] > ndof {
				rz = rnode[nsc[zidx]-ndof-1]
			}
			js[node].React = []float64{rx, ry, rz}			
		}
	}
	
	frmrez[0] = js
	frmrez[1] = ms
	frmrez[2] = dglb
	frmrez[3] = rnode
	frmrez[4] = nsc
	frmrez[5] = ndof
	frmrez[6] = ""
	if ncjt == 2{
		if !mod.Noprnt{frmrez[6] = Trs2dTable(js,ms,dglb,rnode,nsc,ndof,ncjt)}
		mod.Frmtyp = 2
		mod.Frmstr = "2dt"
		mod.Ncjt = ncjt	

		// frmrez[6] = Trs2dTable(js,ms,dglb,rnode,nsc,ndof,ncjt)
		// mod.Frmtyp = 2
		// mod.Frmstr = "2dt"
		// mod.Ncjt = ncjt
		// tstr += FrcTable(mod)
		//mod.Frmtyp = 1
	}
	if ncjt == 3{
		if !mod.Noprnt{frmrez[6] = Trs3dTable(js,ms,dglb,rnode,nsc,ndof,ncjt)}
		mod.Frmtyp = 4
		mod.Frmstr = "3dt"
		mod.Ncjt = ncjt
		// tstr += FrcTable(mod)
		// frmrez[6] = tstr
	}
	frmrez[7] = p
	frmrez[8] = d0
	frmrez[9] = schmat
	mod.Ndof = ndof
	return frmrez, nil
}

func KassT(ndof int, nms int, rezchn chan [][]float64, symchn chan *mat.SymDense) {
	//assembles structure stiffness matrix from rezchan
	//s[na-1][nb-1] = s[na-1][nb-1] + rez
	
	smatsym := mat.NewSymDense(ndof, nil) //to build a symmetric matrix

	s := make([][]float64, ndof)
	//smat := mat.NewDense(ndof, ndof, nil)
	for i := 0; i < ndof; i++ {
		s[i] = make([]float64, ndof)
	}
	for i := 0; i < nms; i++ {
		srez := <-rezchn
		for _, rez := range srez {
			na := int(rez[0]) - 1
			nb := int(rez[1]) - 1
			s[na][nb] = s[na][nb] + rez[2]
			smatsym.SetSym(na, nb, s[na][nb])
		}
	}

	//for i, row := range s {
	//	smat.SetRow(i, row)
	//}

	//schn <- smat
	symchn <- smatsym
	//close(schn)
	//close(symchn)
}

//Kmem3dT generates individual member stiffness matrices for 3d truss members
func Kmem3dT(member int, memprp []int, nsc []int, js map[int]*Node, em [][]float64, cp [][]float64, ndof int, ncjt int, kchn chan []interface{}) {
	//local y and z end displacements are not evaluated? IS NOT
	//0 - member, 1 - l, e, a, cx , cy, cz 2 - bkmat (local k mat)
	//3 - tmat, 4 - gkmat, 5 - srez,
	rez := make([]interface{}, 9)
	jb := memprp[0]
	je := memprp[1]
	e := em[memprp[2]-1][0]
	a := cp[memprp[3]-1][0]
	pa := js[jb].Coords
	pb := js[je].Coords
	xa := pa[0]
	ya := pa[1]
	za := pa[2]
	xb := pb[0]
	yb := pb[1]
	zb := pb[2]

	//cx, cy, cz calcs
	l := math.Sqrt(math.Pow(xb-xa, 2) + math.Pow(yb-ya, 2) + math.Pow(zb-za, 2))
	cx := (xb - xa) / l
	cy := (yb - ya) / l
	cz := (zb - za) / l
	rez[0] = member
	rez[1] = []float64{l, e, a, cx, cy, cz}

	//set member local stiffness matrix
	zk := e * a / l
	bkmat := mat.NewDense(2, 2, []float64{zk, -zk, -zk, zk})
	rez[2] = bkmat

	//member transformation matrix t

	tmat := mat.NewDense(2, 6, []float64{cx, cy, cz, 0, 0, 0, 0, 0, 0, cx, cy, cz})

	rez[3] = tmat

	//member global stiffness matrix gk = t.T @ k @ t
	gkmat := mat.NewDense(6, 6, nil)
	gkmat.Product(tmat.T(), bkmat, tmat)

	rez[4] = gkmat

	//assemble structure stiffness matrix using nsc code numbers
	srez := [][]float64{}
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
					r := gkmat.At(i-1, j-1)
					if r != 0 {
						srez = append(srez, []float64{float64(na), float64(nb), r})
					}

				}
			}

		}
	}
	rez[5] = srez
	//structure fixed end forces due to end displacements
	vfs := make([]float64, 2 * ncjt)
	vf := make([]float64, 2* ncjt)
	memsd := false
	if js[jb].Sd{
		memsd = true
		for i, displ := range js[jb].Sdj{
			vfs[i] += displ
		}
	}
	if js[je].Sd{
		memsd = true
		for i, displ := range js[je].Sdj{
			vfs[i+ncjt] = displ
		}
	}
	copy(vf, vfs)
	pfsrez := [][]float64{}
	if memsd{
		vfmat := mat.NewDense(2*ncjt, 1, vf)
		vfmat.Product(gkmat, vfmat)
		for i := 1; i <= 2 * ncjt; i++{
			if i <= ncjt{
				ela = (jb-1)*ncjt + i
				na = nsc[(ela - 1)]
			} else {
				ela = (je-1)*ncjt + (i - ncjt)
				na = nsc[(ela - 1)]
			}
			if na <= ndof{
				pfsrez = append(pfsrez, []float64{float64(na),vfmat.At(i-1,0)})
			}
		}
	}
	rez[6] = vfs
	rez[7] = pfsrez
	var mtyp int
	switch {
	case cx == 0: //if column
		mtyp = 1
	case cx == 1 || cz == 1: //beam
		mtyp = 2
	default: //then weird thing
		mtyp = 3 
	}
	rez[8] = mtyp
	kchn <- rez
}

//Kmem2dT generates individual member stiffness matrices and does a ton of heavy lifting
func Kmem2dT(member int, memprp []int, nsc []int, js map[int]*Node, em [][]float64, cp [][]float64, ndof int, ncjt int, kchn chan []interface{}) {
	//heavy lifting
	//0 - member, 1 - l, cx , cy, cz, l, e, a 2 - bkmat (local k mat)
	//3 - tmat, 4 - gkmat, 5 - srez,
	//6 - member type (Bm/Col E(L/R)nd/Mid TopBottomMid)
	//make mtyp - 1, 2, etc - FIGURE THIS FOR A TRUSS
	//ADD MEMBER MASS MATRIX HERE
	rez := make([]interface{}, 9)
	jb := memprp[0]
	je := memprp[1]
	e := em[memprp[2]-1][0]
	a := cp[memprp[3]-1][0]
	pa := js[jb].Coords
	pb := js[je].Coords
	xa := pa[0]
	ya := pa[1]
	xb := pb[0]
	yb := pb[1]
	//cx, cy calcs
	l := math.Sqrt(math.Pow(xb-xa, 2) + math.Pow(yb-ya, 2))
	cx := (xb - xa) / l
	cy := (yb - ya) / l
	rez[0] = member
	rez[1] = []float64{l, e, a, cx, cy}

	//set member local stiffness matrix
	zk := e * a / l
	bkmat := mat.NewDense(4, 4, []float64{zk, 0, -zk, 0, 0, 0, 0, 0, -zk, 0, zk, 0, 0, 0, 0, 0})
	rez[2] = bkmat

	//member transformation matrix t

	tmat := mat.NewDense(4, 4, []float64{cx, cy, 0, 0, -cy, cx, 0, 0, 0, 0, cx, cy, 0, 0, -cy, cx})

	rez[3] = tmat

	//member global stiffness matrix gk = t.T @ k @ t
	gkmat := mat.NewDense(2*ncjt, 2*ncjt, nil)
	gkmat.Product(tmat.T(), bkmat, tmat)

	rez[4] = gkmat
	//assemble structure stiffness matrix using nsc code numbers
	srez := [][]float64{}
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
					r := gkmat.At(i-1, j-1)
					if r != 0 {
						srez = append(srez, []float64{float64(na), float64(nb), r})
					}

				}
			}

		}
	}
	rez[5] = srez
	//structure fixed end forces due to end displacements
	vfs := make([]float64, 2 * ncjt)
	vf := make([]float64, 2* ncjt)
	memsd := false
	if js[jb].Sd{
		memsd = true
		//log.Println("member",member,"adding displacements for jb")
		for i, displ := range js[jb].Sdj{
			vfs[i] += displ
		}
	}
	if js[je].Sd{
		memsd = true
		//log.Println("member",member,"adding displacements for je")
		for i, displ := range js[je].Sdj{
			vfs[i+ncjt] = displ
		}
	}
	copy(vf, vfs)
	pfsrez := [][]float64{}
	if memsd{
		vfmat := mat.NewDense(2*ncjt, 1, vf)
		vfmat.Product(gkmat, vfmat)
		for i := 1; i <= 2 * ncjt; i++{
			if i <= ncjt{
				ela = (jb-1)*ncjt + i
				na = nsc[(ela - 1)]
			} else {
				ela = (je-1)*ncjt + (i - ncjt)
				na = nsc[(ela - 1)]
			}
			if na <= ndof{
				pfsrez = append(pfsrez, []float64{float64(na),vfmat.At(i-1,0)})
			}
		}
	}
	rez[6] = vfs
	rez[7] = pfsrez
	var mtyp int
	switch {
	case cx == 0: //if column
		mtyp = 1
	case cx == 1: //beam
		mtyp = 2
	default: //then weird thing (sling)
		mtyp = 3 
	}
	rez[8] = mtyp
	kchn <- rez
}


//EndFrcT gets individual member force vectors for a truss member (2d/3d)
func EndFrcT(member int, mem *Mem, msqf map[int][]float64, dglb []float64, nsc []int, ncjt int, ndof int, rchn chan [][]float64, fchn chan []interface{}) {
	
	vmat := mat.NewDense(2*ncjt, 1, mem.Vfs)
	jb := mem.Mprp[0]
	je := mem.Mprp[1]
	var n int
	for i := 1; i <= 2*ncjt; i++ {
		if i <= ncjt {
			n = (jb-1)*ncjt + i
		} else {
			n = (je-1)*ncjt + (i - ncjt)
		}
		val := vmat.At(i-1,0) + dglb[n-1]
		vmat.Set(i-1, 0, val)
	}
	qfmat := mat.NewDense(2*ncjt, 1, nil)
	if ncjt == 3 {
		qfmat = mat.NewDense(2,1,nil)
	}
	qfmat.Product(mem.Bkmat, mem.Tmat, vmat)

	switch ncjt{
		case 2:
		q1mat := mat.NewDense(len(msqf[member]),1,msqf[member])
		qfmat.Add(qfmat, q1mat)
		case 3:
		v1 := qfmat.At(0,0) + msqf[member][0]
		v2 := qfmat.At(1,0) + msqf[member][3]
		qfmat.Set(0,0,v1)
		qfmat.Set(1,0,v2)
	}
	fmat := mat.NewDense(2*ncjt, 1, nil)
	
	fmat.Product(mem.Tmat.T(), qfmat)
	
	
	//return index and values of fmat (>ndof) for support reaction r
	fvec := make([]float64,2*ncjt)
	qfvec := make([]float64,2*ncjt)
	if ncjt == 3 {
		qfvec = make([]float64, 2)
	}
	rez := [][]float64{}
	for i := 1; i <= 2*ncjt; i++ {
		fvec[i-1] = fmat.At(i-1,0)
		if ncjt == 2 {qfvec[i-1] = qfmat.At(i-1,0)}
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
	if ncjt == 3 {
		qfvec[0] = qfmat.At(0,0)
		qfvec[1] = qfmat.At(1,0)
	}
	//return qfvec (member local force vector) and fvec (global force vec) AND uvec (local end displacement vector)
	frez := make([]interface{},3)
	frez[0] = member
	frez[1] = qfvec
	frez[2] = fvec

	rchn <- rez
	fchn <- frez
}
