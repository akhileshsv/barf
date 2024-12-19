package barf

import (
	"fmt"
	"log"
	"math"
	"errors"
	"gonum.org/v1/gonum/mat"
)

//CalcFrm2d is the entry func for plane frame direct stiffness analysis
//see chapter 6/7 kassimali
//returns frmrez (slice of empty interfaces) shamelessly + error 
func CalcFrm2d(mod *Model, ncjt int) ([]interface{},error) {
	//make member and joint maps
	frmrez := make([]interface{}, 10)
	var err error
	ms := make(map[int]*Mem)
	for i, val := range mod.Mprp {
		
		if len(val) < 4{
			return frmrez, fmt.Errorf("invalid length of mat prop. slice for member %v",i+1)
		}
		jb := val[0]; je := val[1]; em := val[2]; cp := val[3]
		switch{
			case jb < 0 || jb > len(mod.Coords) || len(mod.Coords[jb-1])<2:
			err = fmt.Errorf("invalid node index for member %v -> %v",i+1, jb)
			case je < 0 || je > len(mod.Coords) || len(mod.Coords[je-1])<2:
			err = fmt.Errorf("invalid node index for member %v -> %v",i+1, je)
			case em < 0 || len(mod.Em) < em || len(mod.Em[em-1]) < 1:
			err = fmt.Errorf("invalid em index for member %v -> %v",i+1, em)
			case cp < 0 || len(mod.Cp) < cp || len(mod.Cp[cp-1]) < 2:
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
	
	//get member geometry, stiffness matrices
	kchn := make(chan []interface{}, len(ms))

	for member, m := range ms {
		jb := m.Mprp[0]; je := m.Mprp[1]
		bcon := (jb < 0 || je < 0) || (jb > len(mod.Coords) || je > len(mod.Coords))
		if bcon{
			err := fmt.Errorf("invalid node index for member %v",member)
			return frmrez, err
		}
		js[jb].Mems = append(js[jb].Mems,member)
		js[je].Mems = append(js[je].Mems,member)
		go Kmem(member, m.Mprp, nsc, js, mod.Em, mod.Cp, ndof, ncjt, kchn)

	}

	//call matrix assembly func
	//rezchan to receive rez,  symchn for symmetric banded
	//good ol' pffs (support displacement fixed end force vector)
	pffs := make([]float64, ndof)
	rezchn := make(chan [][]float64, len(ms))
	symchn := make(chan *mat.SymDense)
	go Kass(ndof, len(ms), rezchn, symchn)
	//get member results from Kmem and update dict
	for i := 0; i < len(ms); i++ {
		//0 - member, 1 - cx , cy, l, e, a, iz, 2 - bkmat
		//3 - tmat, 4 - gkmat, 5 - srez,
		//6 - member type (End/Mid Bm/Col) 7 - vfs 8 - pfsrez
		x := <-kchn
		srez, _ := x[5].([][]float64)
		rezchn <- srez
		member, _ := x[0].(int)
		geos, _ := x[1].([]float64)
		ms[member].Geoms = geos
		bkmat, _ := x[2].(*mat.Dense)
		ms[member].Bkmat =  bkmat
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
		jb := ms[member].Mprp[0]; je := ms[member].Mprp[1]
		if (js[jb].Nrs == 0 && len(js[jb].Mems) == 1) || (js[je].Nrs == 0 && len(js[je].Mems) == 1){
			ms[member].Clvr = true
		}
	}
	//receive smatsym (and smat) from Kass
	//smat := <-schn
	smatsym := <-symchn

	pfchn := make(chan []interface{}, 1)
	pchn := make(chan []float64, 1)
	
	//compute forces
	for i, val := range mod.Msloads{
		if len(val) < 6{
			return frmrez, fmt.Errorf("invalid member force slice #(%v)-(%v) for model",i+1, val)
		}
	}
	for i, val := range mod.Jloads{
		if len(val) < 4{
			return frmrez, fmt.Errorf("invalid nodal force slice #(%v)-(%v) for model",i+1, val)
		}
	}
	//MemFrc returns msqf (map of member end forces), pf (force rez; idx and force),
	//msloaded (map of member and individual force vectors); mfrtyp for 2df = 2
	//NodeFrc assembles the nodal force vector p 
	go MemFrc(ms, mod.Msloads, nsc, ncjt, ndof, 2, pfchn)
	go NodeFrc(mod.Jloads, nsc, ncjt, ndof, pchn)
	pfrez := <-pfchn
	msqf, _ := pfrez[0].(map[int][]float64)
	pf, _ := pfrez[1].([]float64)
	msloaded, _ := pfrez[2].(map[int][][]float64)
	p := <-pchn
	
	//call RccD with msqf and msloaded
	//make matriksees
	
	pmat := mat.NewDense(ndof, 1, nil)
	for i, f := range pf {
		p[i] = p[i] - f -pffs[i]
		pmat.Set(i, 0, p[i])
	}
	
	var schmat mat.Cholesky
	var dchmat mat.Dense
	if ok := schmat.Factorize(smatsym); !ok {
		return frmrez, errors.New("non +ve definite stiffness matrix")
	}
	if err := schmat.SolveTo(&dchmat, pmat); err != nil {
		return frmrez, errors.New("near singular stiffness matrix")
	}

	//check if sd == p
	//does one need to really? IDIOT

	//build global displacement vector
	//d0 again is the nl analysis displacement vector
	dglb := make([]float64, nj*ncjt)
	var d0 []float64
	for i, sc := range nsc {
		if sc <= ndof {
			dglb[i] = dchmat.At(sc-1, 0)
			d0 = append(d0, dchmat.At(sc-1, 0))
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
		qfvec,_ := frez[2].([]float64)
		msqf[mem] = fvec
		ms[mem].Qf = qfvec
		ms[mem].Gf = fvec
	}
	/*
	switch dsgn {
	case 1:
		RccD(ncjt, dsgnchn, ms, msqf, msloaded, em, cp)
		rez := <- dsgnchn
		//close(chndexchn)
		bmrez,_ := rez[0].(map[int]BeamRez)
		frmrez[6] = bmrez
	}
	*/
	for node  := range js {
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
	frmrez[0] = js
	frmrez[1] = ms
	frmrez[2] = dglb
	frmrez[3] = rnode
	frmrez[4] = msqf
	frmrez[5] = msloaded
	frmrez[6] = ""
	if !mod.Noprnt{frmrez[6] = Frm2dTable(js, ms, dglb, rnode, nsc, ndof, ncjt)}
	frmrez[7] = p
	frmrez[8] = d0
	frmrez[9] = schmat
	mod.Ncjt = ncjt
	mod.Frmtyp = 3
	mod.Frmstr = "2df"
	mod.Ndof = ndof
	return frmrez, nil
}

//Kass assembles the overall structure stiffness matrix from rezchan (individual member stiffness matrix builder results)
//func Kass(ndof int, nms int, rezchn chan [][]float64, schn chan *mat.Dense, symchn chan *mat.SymDense)
func Kass(ndof int, nms int, rezchn chan [][]float64, symchn chan *mat.SymDense) {
	//assembles structure stiffness matrix from rezchan
	//s[na-1][nb-1] = s[na-1][nb-1] + rez

	smatsym := mat.NewSymDense(ndof, nil) // (there was an attempt) to build a symmetric matrix

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
	/*
	for i, row := range s {
		smat.SetRow(i, row)
	}
	*/
	//schn <- smat
	symchn <- smatsym
	//close(schn)
	//close(symchn)
}

//Kmem generates individual member stiffness matrices and etc from model data
func Kmem(member int, memprp []int, nsc []int, js map[int]*Node, em [][]float64, cp [][]float64, ndof int, ncjt int, kchn chan []interface{}){
	//0 - member, 1 - l, e, a, iz, cx, cy, 2 - bkmat
	//3 - tmat, 4 - gkmat, 5 - srez,
	//6 - member type (Bm/Col E(L/R)nd/Mid TopBottomMid)
	//
	rez := make([]interface{}, 9)
	jb := memprp[0]
	je := memprp[1]
	var e, a, iz float64
	//stupid bounds check
	if len(em[memprp[2]-1]) > 0{e = em[memprp[2]-1][0]}
	if len(cp[memprp[3]-1]) > 1{		
		a = cp[memprp[3]-1][0]
		iz = cp[memprp[3]-1][1]	
	}
	memrel := memprp[4]
	pa := js[jb].Coords
	pb := js[je].Coords
	xa := pa[0]
	ya := pa[1]
	xb := pb[0]
	yb := pb[1]
	l := math.Sqrt(math.Pow(xb-xa, 2) + math.Pow(yb-ya, 2))
	cx := (xb - xa) / l
	cy := (yb - ya) / l
	//MAKE MEMTYP INT for frame; make this func great again
	var mtyp int
	switch {
	case cx == 0: //if column
		mtyp = 1
	case cx == 1: //beam
		mtyp = 2
	default: //then weird thing
		mtyp = 3 
	}
	bk := make([][]float64, 2*ncjt)
	t := make([][]float64, 2*ncjt)
	for i := 0; i < 2*ncjt; i++ {
		bk[i] = make([]float64, 2*ncjt)
		t[i] = make([]float64, 2*ncjt)
	}
	//local stiffness matrix bk for each memrel type
	za := (e * a) / l
	zb := (e * iz) / math.Pow(l, 3)
	switch memrel {
	case 0: //fixed at both ends
		bk[0][0] = za
		bk[3][0] = -za
		bk[0][3] = -za
		bk[3][3] = za
		bk[1][1] = 12 * zb
		bk[1][2] = 6 * l * zb
		bk[1][4] = -12 * zb
		bk[1][5] = 6 * l * zb
		bk[2][1] = 6 * l * zb
		bk[2][2] = 4 * l * l * zb
		bk[2][4] = -6 * l * zb
		bk[2][5] = 2 * l * l * zb
		bk[4][1] = -12 * zb
		bk[4][2] = -6 * l * zb
		bk[4][4] = 12 * zb
		bk[4][5] = -6 * l * zb
		bk[5][1] = 6 * zb * l
		bk[5][2] = 2 * l * l * zb
		bk[5][4] = -6 * zb * l
		bk[5][5] = 4 * l * l * zb
	case 1: //hinge at beginning
		bk[0][0] = za
		bk[3][0] = -za
		bk[0][3] = -za
		bk[3][3] = za
		bk[1][1] = 3 * zb
		bk[1][4] = -3 * zb
		bk[1][5] = 3 * l * zb
		bk[4][1] = -3 * zb
		bk[4][4] = 3 * zb
		bk[4][5] = -3 * l * zb
		bk[5][1] = 3 * zb * l
		bk[5][4] = -3 * zb * l
		bk[5][5] = 3 * l * l * zb
	case 2: //hinge at end
		bk[0][0] = za
		bk[3][0] = -za
		bk[0][3] = -za
		bk[3][3] = za
		bk[1][1] = 3 * zb
		bk[1][2] = 3 * l * zb
		bk[1][4] = -3 * zb
		bk[2][1] = 3 * l * zb
		bk[2][2] = 3 * l * l * zb
		bk[2][4] = -3 * l * zb
		bk[4][1] = -3 * zb
		bk[4][2] = -3 * l * zb
		bk[4][4] = 3 * zb
	case 3: //hinge at both ends
		bk[0][0] = za
		bk[3][0] = -za
		bk[0][3] = -za
		bk[3][3] = za
	}

	bkmat := mat.NewDense(2*ncjt, 2*ncjt, nil)
	for i, row := range bk {
		bkmat.SetRow(i, row)
	}
	//member transformation matrix t

	t[0][0] = cx
	t[0][1] = cy
	t[1][0] = -cy
	t[1][1] = cx
	t[2][2] = 1
	t[3][3] = cx
	t[3][4] = cy
	t[4][3] = -cy
	t[4][4] = cx
	t[5][5] = 1
	tmat := mat.NewDense(2*ncjt, 2*ncjt, nil) //FOOL
	for i, row := range t { //foolish fool
		tmat.SetRow(i, row)
	}
	//member global stiffness matrix gk = t.T @ k @ t
	gkmat := mat.NewDense(2*ncjt, 2*ncjt, nil)
	gkmat.Product(tmat.T(), bkmat, tmat)
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
	//ADD - mass matrix, tangent stiffness matrix
	rez[0] = member
	rez[1] = []float64{l, e, a, iz, cx, cy}
	rez[2] = bkmat
	rez[3] = tmat
	rez[4] = gkmat
	rez[5] = srez
	rez[6] = mtyp
	rez[7] = vfs
	rez[8] = pfsrez
	kchn <- rez
}
