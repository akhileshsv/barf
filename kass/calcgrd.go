package barf

import (
	"fmt"
	"log"
	"math"
	"errors"
	"gonum.org/v1/gonum/mat"
)

//CalcGrd is the entry func for grid/grillage direct stiffness analysis
//see chapter 8, kassimali
//returns frmrez (slice of empty interfaces) + error 
func CalcGrd(mod *Model, ncjt int) ([]interface{}, error) {
	//make member and joint maps
	frmrez := make([]interface{}, 7)
	var err error
	ms := make(map[int]*Mem)
	for i, val := range mod.Mprp{
		
		if len(val) < 4{
			return frmrez, fmt.Errorf("invalid length of mat prop. slice for member %v",i+1)
		}
		jb := val[0]; je := val[1]; em := val[2]; cp := val[3]
		switch{
			case jb < 1 || jb > len(mod.Coords) || len(mod.Coords[jb-1])<3:
			err = fmt.Errorf("invalid node index for member %v -> %v",i+1, jb)
			case je < 1 || je > len(mod.Coords) || len(mod.Coords[je-1])<3:
			err = fmt.Errorf("invalid node index for member %v -> %v",i+1, je)
			case em < 1 || len(mod.Em) < em || len(mod.Em[em-1]) < 2:
			err = fmt.Errorf("invalid em index for member %v -> %v",i+1, em)
			case cp < 1 || len(mod.Cp) < cp || len(mod.Cp[cp-1]) < 2:
			err = fmt.Errorf("invalid em index for member %v -> %v",i+1, cp)
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

	if ndof <= 0{
		return frmrez, fmt.Errorf("invalid degrees of freedom (%v) for model",ndof)
	}
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
			}
		}
	}
	for i, n := range nsc {
		if n == 0 {
			nsc[i] = ndis
			ndis++
		}
	}
	//ADD SUPPORT DISPLACEMENTS DAMMIT
	for idx := range js {
		nscidx := (idx-1) * ncjt
		js[idx].Nscs = nsc[nscidx:nscidx+ncjt]
		js[idx].Sdj = make([]float64,ncjt)
	}	
	for i, idx := range mod.Jsd{
		sdj := mod.Sdj[i]
		if len(sdj) == ncjt{
			copy(js[idx].Sdj, sdj)
			js[idx].Sd = true
		} else {
			log.Println("ERRORE,errore->error reading support displacement for node->",idx)
		}
	}
	kchn := make(chan []interface{}, len(ms))
	
	for member, mem := range ms {
		jb := mem.Mprp[0]; je := mem.Mprp[1]
		bcon := (jb < 0 || je < 0) || (jb > len(mod.Coords) || je > len(mod.Coords))
		if bcon{
			err := fmt.Errorf("invalid node index for member %v",member)
			return frmrez, err
		}
		js[jb].Mems = append(js[jb].Mems,member)
		js[je].Mems = append(js[je].Mems,member)
		go KmemG(member, mem.Mprp, nsc, js, mod.Em, mod.Cp, ndof, ncjt, kchn)
	}
	pffs := make([]float64, ndof)
	rezchn := make(chan [][]float64, len(ms))
	symchn := make(chan *mat.SymDense)
	go KassG(ndof, len(ms), rezchn, symchn) 
	
	//get member results from Kmem3dG and update dict
	for i := 0; i < len(ms); i++ {
		//0 - member, 1 - l, e, a, iz, cx, cy, cz 2 - bkmat
		//3 - tmat, 4 - gkmat, 5 - srez,
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
		ms[member].Tmat =  tmat
		gkmat, _ := x[4].(*mat.Dense)
		ms[member].Gkmat = gkmat
		mtyp, _ := x[6].(int)
		ms[member].Mtyp = mtyp
		vfs, _ := x[7].([]float64)
		ms[member].Vfs = vfs
		pfsrez, _ := x[8].([][]float64)
		if len(pfsrez) != 0{
			for _, v := range pfsrez{
				idx := int(v[0]); val := v[1]
				pffs[idx-1] += val
			}
		}

	}

	//receive smatsym from Kass
	
	smatsym := <-symchn
	pchn := make(chan []float64, 1)
	pfchn := make(chan []interface{}, 1)
	//compute forces
	for i, val := range mod.Msloads{
		if len(val) < 6{
			return frmrez, fmt.Errorf("invalid member force slice #(%v)-(%v) for model",i+1, val)
		}
	}
	for i, val := range mod.Jloads{
		if len(val) < 3{
			return frmrez, fmt.Errorf("invalid nodal force slice #(%v)-(%v) for model",i+1, val)
		}
	}
	//MemFrc returns msqf (map of member end forces), pf (force rez; idx and force),
	//msloaded (map of member and individual force vectors)
	//NodeFrc assembles the nodal force vector p (of p-pf = sd fame)

	go MemFrcG(ms, mod.Msloads, nsc, ncjt, ndof, pfchn)
	go NodeFrc(mod.Jloads, nsc, ncjt, ndof, pchn)

	p := <-pchn

	pfrez := <-pfchn
	msqf, _ := pfrez[0].(map[int][]float64)
	pf, _ := pfrez[1].([]float64)
	msloaded, _ := pfrez[2].(map[int][][]float64)

	for i, f := range pf {
		p[i] = p[i] - f - pffs[i]
	}
	pmat := mat.NewDense(ndof, 1, p)

	//solve by cholesky decomposition of smatsym

	var schmat mat.Cholesky
	var dchmat mat.Dense
	if ok := schmat.Factorize(smatsym); !ok {
		return frmrez, errors.New("non +ve definite stiffness matrix")
	}
	if err := schmat.SolveTo(&dchmat, pmat); err != nil {
		return frmrez, errors.New("near singular stiffness matrix")
	}



	//check if sd == p

	pchkchol := mat.NewDense(ndof, 1, nil)
	pchkchol.Product(&schmat, &dchmat)
	
	for i := 0; i < ndof; i++ {
		
		if math.Abs(pchkchol.At(i, 0)-pmat.At(i, 0)) >= 1e-5 {
			return frmrez, errors.New("Matrix multiplication check error")			
		} else {
			continue
		}

	}

	//build global displacement vector

	dglb := make([]float64, nj*ncjt)
	for i, sc := range nsc {
		if sc <= ndof {
			dglb[i] = dchmat.At(sc-1, 0)
		} else {
			dglb[i] = 0
		}
	}

	rchn := make(chan [][]float64, len(ms))
	fchn := make(chan []interface{}, len(ms))
	for member, mem := range ms {
		go EndFrc(member, mem, msqf, dglb, nsc, ncjt, ndof, rchn, fchn)
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

	for i :=0; i < len(ms); i++ {
		frez := <- fchn
		mem, _ := frez[0].(int)
		fvec,_ := frez[1].([]float64)
		msqf[mem] = fvec
		ms[mem].Qf = fvec
	}
	
	for node, _ := range js {
		yidx := (node -1)*ncjt
		xidx := (node -1)*ncjt + 1
		zidx := (node -1)*ncjt + 2
		ydis := dglb[yidx]
		xrot := dglb[xidx]
		zrot := dglb[zidx]
		js[node].Displ = []float64{ydis, xrot, zrot}
		var ry, mx, mz float64
		//mx is the torsional moment. TORSION.
		//IT IS TORSION fyi
		if nsc[yidx] > ndof {
			ry = rnode[nsc[yidx]-ndof-1]
		}
		if nsc[xidx] > ndof {
			mx =  rnode[nsc[xidx]-ndof-1]
		}
		if nsc[zidx] > ndof {
			mz = rnode[nsc[zidx]-ndof-1]
		}
		js[node].React = []float64{ry, mx, mz}
	}
	
	frmrez[0] = js
	frmrez[1] = ms
	frmrez[2] = dglb
	frmrez[3] = rnode
	frmrez[4] = msqf
	frmrez[5] = msloaded
	frmrez[6] = ""
	if !mod.Noprnt{frmrez[6] = Grd3dTable(js, ms, dglb, rnode, nsc, ndof, ncjt)}
	mod.Ncjt = ncjt
	mod.Frmtyp = 5
	mod.Frmstr = "3dg"
	return frmrez, nil
}

//KmemG assembles grid/grillage member stiffness matrices
func KmemG(member int, memprp []int, nsc []int, js map[int]*Node, em [][]float64, cp [][]float64, ndof int, ncjt int, kchn chan []interface{}) {
	//sc numbers in nsc are y disp, rot x, rot z.
	//0 - member, 1 - cx , cy, cz, l, e, gu, iz, jv
	//2 - bkmat (local k mat)
	//3 - tmat, 4 - gkmat, 5 - srez
	var bkarr []float64
	rez := make([]interface{}, 9)
	jb := memprp[0]
	je := memprp[1]
	e := em[memprp[2]-1][0]
	gu := em[memprp[2]-1][1]
	iz := cp[memprp[3]-1][0]
	jv := cp[memprp[3]-1][1]
	memrel := memprp[4]
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
	rez[1] = []float64{l, e, gu, iz, jv, cx, cy, cz}

	//set member local stiffness matrix CHANGE THESE VARIABLES LOL
	za = e * iz / math.Pow(l, 3)
	zb = gu * jv / l
	switch memrel {
	case 0: //rigidly connected at both ends
		bkarr = []float64{12 * za, 0, 6 * l * za, -12 * za, 0, 6 * l * za,
			0, zb, 0, 0, -zb, 0,
			6 * l * za, 0, 4 * l * l * za, -6 * l * za, 0, 2 * l * l * za,
			-12 * za, 0, -6 * l * za, 12 * za, 0, -6 * l * za,
			0, -zb, 0, 0, zb, 0,
			6 * l * za, 0, 2 * l * l * za, -6 * l * za, 0, 4 * l * l * za}
	case 1: //hinge at beginning
		bkarr = []float64{3 * za, 0, 0, -3 * za, 0, 3 * l * za,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			-3 * za, 0, 0, 3 * za, 0, -3 * l * za,
			0, 0, 0, 0, 0, 0,
			3 * l * za, 0, 0, -3 * l * za, 0, 3 * l * l * za}
	case 2: //hinge at end
		bkarr = []float64{3 * za, 0, 3 * l * za, -3 * za, 0, 0,
			0, 0, 0, 0, 0, 0,
			3 * l * za, 0, 3 * l * l * za, -3 * l * za, 0, 0,
			-3 * za, 0, -3 * l * za, 3 * za, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0}
	case 3: //hinge at both ends, zero matrix like a n00b
		bkarr = []float64{0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0}

	}
	bkmat := mat.NewDense(2*ncjt, 2*ncjt, bkarr)
	rez[2] = bkmat

	//member transformation matrix t (here sin0 = cz)
	tarr := []float64{1, 0, 0, 0, 0, 0,
		0, cx, cz, 0, 0, 0,
		0, -cz, cx, 0, 0, 0,
		0, 0, 0, 1, 0, 0,
		0, 0, 0, 0, cx, cz,
		0, 0, 0, 0, -cz, cx,
	}
	tmat := mat.NewDense(2*ncjt, 2*ncjt, tarr)

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
	var mtyp int
	switch {
	case cx == 0: //if column
		mtyp = 1
	case cx == 1: //beam
		mtyp = 2
	default: //then weird thing
		mtyp = 3 
	}
	rez[6] = mtyp
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
	rez[7] = vfs
	rez[8] = pfsrez
	kchn <- rez
}

//KassG assembles structure stiffness matrix from rezchan
func KassG(ndof int, nms int, rezchn chan [][]float64, symchn chan *mat.SymDense) {

	smatsym := mat.NewSymDense(ndof, nil)

	s := make([][]float64, ndof)
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
	symchn <- smatsym
	close(symchn)
}

//MemFrcG assembles member fixed end force vectors for a grillage model
func MemFrcG(ms map[int]*Mem, msp [][]float64, nsc []int, ncjt int, ndof int, pfchn chan []interface{}) {
	//map of loaded members msloaded[0] ex. - 1,1,40.25,0,134.16,0
	//frmtyp - 0 2D frame, 1 3D grid, 2 3D frame

	msloaded := make(map[int][][]float64)

	qfchn := make(chan []interface{}, len(msp))
	for _, ldcase := range msp {
		member := int(ldcase[0])
		msloaded[member] = append(msloaded[member], ldcase)
		ltyp := int(ldcase[1])
		memrel := ms[member].Mprp[4]
		l := ms[member].Geoms[0]                    
		go FxdEndFrc(member, memrel, ltyp, l, ldcase[2:], 1, ms[member].Geoms, qfchn) //mfrtyp = 1
	}
	msqf := make(map[int][]float64)
	for member := range ms {
		msqf[member] = []float64{0, 0, 0, 0, 0, 0}
	}

	//get member local end forces qf
	for i := 0; i < len(msp); i++ {
		r := <-qfchn
		member, _ := r[0].(int)
		qf, _ := r[1].([]float64)
		for i, f := range qf {
			msqf[member][i] += f
		}
	}

	//get member global end force vector ff
	ffchn := make(chan [][]float64, len(msloaded))
	for member, mem := range ms {
		_, ok := msloaded[member]
		if ok {
			qf := msqf[member]
			//tmat, _ := ms[member][3].(*mat.Dense)
			//mprp, _ := ms[member][0].([]int)
			//jb := mprp[0]
			//je := mprp[1]
			go ffmatG(mem.Mprp[0], mem.Mprp[1], qf, mem.Tmat, nsc, ncjt, ndof, ffchn)
		} else {
			continue
			//fmt.Println("zero load member wtf")
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

	rez := make([]interface{}, 3)
	rez[0] = msqf
	rez[1] = pf
	rez[2] = msloaded
	pfchn <- rez
}

//ffmatG assembles structure fixed end force vector from a grillage member data
func ffmatG(jb int, je int, qf []float64, tmat *mat.Dense, nsc []int, ncjt int, ndof int, ffchn chan [][]float64) {
	rez := [][]float64{}
	qfmat := mat.NewDense(2*ncjt, 1, nil)
	for i, f := range qf {
		qfmat.Set(i, 0, f)
	}
	ffmat := mat.NewDense(2*ncjt, 1, nil)
	ffmat.Product(tmat.T(), qfmat)

	var el int
	for i := 1; i <= 2*ncjt; i++ {
		if i <= ncjt {
			el = (jb-1)*ncjt + i
		} else {
			el = (je-1)*ncjt + (i - ncjt)
		}
		n := nsc[el-1]
		if n <= ndof {
			rez = append(rez, []float64{float64(n - 1), ffmat.At(i-1, 0)})
		}
	}
	ffchn <- rez
}
