package barf

import (
	"fmt"
	"log"
	"math"
	"errors"
	"gonum.org/v1/gonum/mat"
)

//CalcBm1d is the entry func for direct stiffness analysis of a beam model 
//see kassimali chapter 5 for details
//returns frmrez as a slice of empty interfaces, err 
func CalcBm1d(mod *Model, ncjt int) ([]interface{}, error) {
	//make member and joint maps
	frmrez := make([]interface{}, 7)
	var err error
	ms := make(map[int]*Mem)
	for i, val := range mod.Mprp {
		//baah.bounds check.baah
		if len(val) < 4{
			return frmrez, fmt.Errorf("invalid length of mat prop. slice for member %v",i+1)
		}
		jb := val[0]; je := val[1]; em := val[2]; cp := val[3]
		switch{
			case jb < 1 || jb > len(mod.Coords) || len(mod.Coords[jb-1]) < 1 :
			err = fmt.Errorf("invalid node index/coords for member %v -> %v",i+1, jb)
			case je < 1 || je > len(mod.Coords)|| len(mod.Coords[je-1]) < 1 :
			err = fmt.Errorf("invalid node index/coords for member %v -> %v",i+1, je)
			case em < 1 || len(mod.Em) < em || len(mod.Em[em-1]) < 1:
			err = fmt.Errorf("invalid em index for member %v -> %v",i+1, em)
			case cp < 1 || len(mod.Cp) < cp || len(mod.Cp[cp-1]) < 1:
			err = fmt.Errorf("invalid cp index for member %v -> %v",i+1, cp)
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
				js[val[0]].Nrs++
			}
		}
	}
	ndof := ncjt*nj - nr
	if ndof <= 0 {
		return frmrez, fmt.Errorf("negative number of degrees of freedom (%v) for model",ndof)
	}
	nsc := make([]int, nj*ncjt)
	ndis := 1
	nres := ndof + 1
	for _, s := range mod.Supports {
		if len(s) < 3{
			return frmrez, fmt.Errorf("invalid length of support slice->%v",s)
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
		js[idx].Sdj = make([]float64,ncjt)
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
	
	kchn := make(chan []interface{}, len(ms))
	for member, m := range ms{
		jb := m.Mprp[0]
		je := m.Mprp[1]
		js[jb].Mems = append(js[jb].Mems,member)
		js[je].Mems = append(js[je].Mems,member)
		go KmemB(member, m.Mprp, nsc, js, mod.Em, mod.Cp, ndof, ncjt, kchn)
	}
	//call matrix assembly func (it waits, vacantly)
	//rezchan to receive rez, symchn for symmetric banded
	rezchn := make(chan [][]float64, len(ms))
	symchn := make(chan *mat.SymDense)
	go KassB(ndof, len(ms), rezchn, symchn) //removing schn

	//get member results
	//pffs - fixed force vector due to support settlements
	pffs := make([]float64, ndof)
	for i := 0; i < len(ms); i++ {
		//0 - member, 1 - l, e, iz, area 2 - bkmat
		//3 - srez 4 vfs (ufs) 5 pfsrez
		x := <-kchn
		srez, _ := x[3].([][]float64)
		rezchn <- srez
		member, _ := x[0].(int)
		geos, _ := x[1].([]float64)
		ms[member].Geoms =  geos
		bkmat, _ := x[2].(*mat.Dense)
		ms[member].Bkmat =  bkmat
		vfs, _ := x[4].([]float64)
		ms[member].Vfs =  vfs
		pfsrez, _ := x[5].([][]float64)
		if len(pfsrez) != 0{
			for _, v := range pfsrez{
				idx := int(v[0]); val := v[1]
				pffs[idx-1] += val
			}
		}
		jb := ms[member].Mprp[0]; je := ms[member].Mprp[1]
		if (js[jb].Nrs == 0 && len(js[jb].Mems) == 1) || (js[je].Nrs == 0 && len(js[je].Mems) == 1){
			ms[member].Clvr = true
		}

	}
	//receive smatsym from KassB
	smatsym := <-symchn
	pchn := make(chan []float64, 1)
	pfchn := make(chan []interface{}, 1)
	//compute forces
	//first check em
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
	go MemFrc(ms, mod.Msloads, nsc, ncjt, ndof, 0, pfchn)
	go NodeFrc(mod.Jloads, nsc, ncjt, ndof, pchn)
	p := <-pchn
	pfrez := <-pfchn
	msqf, _ := pfrez[0].(map[int][]float64)
	pf, _ := pfrez[1].([]float64)
	msloaded, _ := pfrez[2].(map[int][][]float64)
	//subtract pf and pffs from p
	for i, f := range pf {
		p[i] = p[i] - f - pffs[i]
	}

	pmat := mat.NewDense(ndof, 1, p)
	//solve by cholesky decomposition of sym banded smatsym

	var schmat mat.Cholesky
	var dchmat mat.Dense
	if ok := schmat.Factorize(smatsym); !ok {
		return frmrez, errors.New("non +ve definite stiffness matrix")
	}
	if err := schmat.SolveTo(&dchmat, pmat); err != nil {
		return frmrez, errors.New("near singular stiffness matrix")
	}
	//check if sd == p
	//? ?? ???
	/*
	pchkchol := mat.NewDense(ndof, 1, nil)
	pchkchol.Product(&schmat, &dchmat)
	for i := 0; i < ndof; i++ {
		if math.Abs(pchkchol.At(i, 0)-pmat.At(i, 0)) >= 1e-5 {
			return frmrez, errors.New("Matrix multiplication check error")
		} else {
			continue
		}

	}
	*/
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
		go EndFrcB(member, mem, msqf, dglb, nsc, ncjt, ndof, rchn, fchn)
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

	for node := range js{
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
	}	
	frmrez[0] = js
	frmrez[1] = ms
	frmrez[2] = dglb
	frmrez[3] = rnode
	frmrez[4] = msqf
	frmrez[5] = msloaded
	if mod.Noprnt{
		frmrez[6] = ""
	}else{
		frmrez[6] = Bm1dTable(js,ms,dglb,rnode,nsc,ndof,ncjt)
	}
	mod.Ncjt = ncjt
	mod.Frmtyp = 1
	mod.Frmstr = "1db"
	return frmrez, nil
}

//KassB assembles a beam structure stiffness matrix from rezchan
//s[na-1][nb-1] = s[na-1][nb-1] + rez
func KassB(ndof int, nms int, rezchn chan [][]float64, symchn chan *mat.SymDense) {
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
	//close(symchn)
}

//KmemB assembles local member stiffness matrices from a beam model data 
func KmemB(member int, memprp []int, nsc []int, js map[int]*Node, em [][]float64, cp [][]float64, ndof int, ncjt int, kchn chan []interface{}) {
	//local y and z end displacements are not evaluated?
	//0 - member, 1 - l, e, iz, a 2 - bkmat (local k mat)
	//3 - srez, 4 - vfs, 5 - pfsrez
	rez := make([]interface{}, 6)
	jb := memprp[0]
	je := memprp[1]
	e := em[memprp[2]-1][0]
	iz := cp[memprp[3]-1][0]
	var ar float64
	if len(cp[0]) > 1 {
		ar = cp[memprp[3]-1][1]
	}
	memrel := 0
	if len(memprp)>4{
		memrel = memprp[4]
	}
	pa := js[jb].Coords
	pb := js[je].Coords
	xa := pa[0]
	xb := pb[0]
	l := math.Abs(xb - xa)
	rez[0] = member
	rez[1] = []float64{l, e, iz, ar}
	//set member local stiffness matrix
	zk := e * iz / math.Pow(l, 3)
	var bkarr []float64
	switch memrel {
	case 0: //fixed at both ends
		bkarr = []float64{
			12 * zk, 6 * l * zk, -12 * zk, 6 * l * zk,
			6 * l * zk, 4 * math.Pow(l, 2) * zk, -6 * l * zk, 2 * math.Pow(l, 2) * zk,
			-12 * zk, -6 * l * zk, 12 * zk, -6 * l * zk,
			6 * l * zk, 2 * math.Pow(l, 2) * zk, -6 * l * zk, 4 * math.Pow(l, 2) * zk,
		}
	case 1: //hinge at the beginning
		bkarr = []float64{
			3 * zk, 0, -3 * zk, 3 * l * zk,
			0, 0, 0, 0,
			-3 * zk, 0, 3 * zk, -3 * l * zk,
			3 * l * zk, 0, -3 * l * zk, 3 * math.Pow(l, 2) * zk,
		}
	case 2: //hinge at end
		bkarr = []float64{
			3 * zk, 3 * l * zk, -3 * zk, 0,
			3 * l * zk, 3 * math.Pow(l, 2) * zk, -3 * l * zk, 0,
			-3 * zk, -3 * l * zk, 3 * zk, 0,
			0, 0, 0, 0,
		}
	case 3: //both ends hinged
		bkarr = []float64{
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
	}
	bkmat := mat.NewDense(4, 4, bkarr)
	rez[2] = bkmat
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
					r := bkmat.At(i-1, j-1)
					if r != 0 {
						srez = append(srez, []float64{float64(na), float64(nb), r})
					}
				}
			}
		}
	}
	rez[3] = srez
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
		ufmat := mat.NewDense(2*ncjt, 1, vf)
		ufmat.Product(bkmat,ufmat)
		for i := 1; i <= 2 * ncjt; i++{
			if i <= ncjt{
				ela = (jb-1)*ncjt + i
				na = nsc[(ela - 1)]
			} else {
				ela = (je-1)*ncjt + (i - ncjt)
				na = nsc[(ela - 1)]
			}
			if na <= ndof{
				pfsrez = append(pfsrez, []float64{float64(na),ufmat.At(i-1,0)})
			}
		}
	}
	rez[4] = vfs
	rez[5] = pfsrez
	kchn <- rez
}

//EndFrcB gets member end forces in global coordinates
func EndFrcB(member int, mem *Mem, msqf map[int][]float64, dglb []float64, nsc []int, ncjt int, ndof int, rchn chan [][]float64, fchn chan []interface{}) {
	//mp, _ := memprp[0].([]int)
	//bkmat := 
	qf := msqf[member]
	//mem.Qf = msqf[member]
	qfmat := mat.NewDense(len(qf), 1, qf)
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
	fmat := mat.NewDense(2*ncjt, 1, nil)
	fmat.Product(mem.Bkmat, vmat)
	fmat.Add(fmat, qfmat)
	//return index and values of fmat (>ndof) for support reaction r
	rez := [][]float64{}
	fvec := make([]float64,2*ncjt)
	for i := 1; i <= 2*ncjt; i++ {
		fvec[i-1] = fmat.At(i-1,0)
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
	//WHY NOT MODIFY THE STRUCT HERE *fool
	frez := make([]interface{},2)
	frez[0] = member
	frez[1] = fvec
	rchn <- rez
	fchn <- frez
}
