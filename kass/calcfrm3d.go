package barf

import (
	"fmt"
	"math"
	"log"
	"errors"
	"gonum.org/v1/gonum/mat"
)


//CalcFrm3d performs direct stiffness analysis of a space frame model
//chapter 8, kassimali
//returns frmrez (slice of empty interfaces) and an error
func CalcFrm3d(mod *Model, ncjt int) ([]interface{}, error){
	//init dicts, angle of roll stays in sep float Wng
	frmrez := make([]interface{}, 7)
	var err error
	ms := make(map[int]*Mem)
	for i, val := range mod.Mprp{
		if len(val) < 4{
			return frmrez, fmt.Errorf("invalid length of mat prop. slice for member %v",i+1)
		}
		jb := val[0]; je := val[1]; em := val[2]; cp := val[3]
		switch{
			case jb < 0 || jb > len(mod.Coords) || len(mod.Coords[jb-1])<3:
			err = fmt.Errorf("invalid node index for member %v -> %v",i+1, jb)
			case je < 0 || je > len(mod.Coords) || len(mod.Coords[je-1])<3:
			err = fmt.Errorf("invalid node index for member %v -> %v",i+1, je)
			case len(mod.Em) < em || len(mod.Em[em-1]) < 2:
			err = fmt.Errorf("invalid em index for member %v -> %v",i+1, em)
			case len(mod.Cp) < cp || len(mod.Cp[cp-1]) < 4:
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
		js[val[0]].Supports = val[1:]		
		for _, i := range val[1:] {
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
	//ADD SUPPORT DISPLACEMENTS *sigh
	//mark nodes with nscs
	for idx  := range js {
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
		go KmemF3d(member, m.Mprp, nsc, js, mod.Em, mod.Cp, mod.Wng, ndof, ncjt, kchn)
	}
	//call matrix assembly func 
	pffs := make([]float64, ndof)
	rezchn := make(chan [][]float64, len(ms))
	symchn := make(chan *mat.SymDense)
	go KassF3d(ndof, len(ms), rezchn, symchn) 
	//get member results from Kmem3dG and update dict
	for i := 0; i < len(ms); i++ {
		//0 - member, 1 - cx , cy, l, e, a, iz, 2 - bkmat
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
		ms[member].Tmat = tmat
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
	//receive smatsym (and smat) from Kass
	smatsym := <-symchn
	pchn := make(chan []float64, 1)
	pfchn := make(chan []interface{}, 1)
	//compute forces
	
	for i, val := range mod.Msloads{
		if len(val) < 8{
			return frmrez, fmt.Errorf("invalid member force slice #(%v)-(%v) for model",i+1, val)
		}
	}
	for i, val := range mod.Jloads{
		if len(val) < 7{
			return frmrez, fmt.Errorf("invalid nodal force slice #(%v)-(%v) for model",i+1, val)
		}
	}
	//MemFrc returns msqf (map of member end forces), pf (force rez; idx and force),
	//msloaded (map of member and individual force vectors)
	//NodeFrc assembles the nodal force vector p (of p-pf = sd fame)
	//mfrtyp = 3 (ND)
	go MemFrc(ms, mod.Msloads, nsc, ncjt, ndof, 3, pfchn)
	go NodeFrc(mod.Jloads, nsc, ncjt, ndof, pchn)
	
	pfrez := <-pfchn
	msqf, _ := pfrez[0].(map[int][]float64)
	pf, _ := pfrez[1].([]float64)
	msloaded, _ := pfrez[2].(map[int][][]float64)
	p := <-pchn
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
	//no?

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
	for member, memprp := range ms {
		go EndFrc(member, memprp, msqf, dglb, nsc, ncjt, ndof, rchn, fchn)
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
	for node := range js {
		//now these variable names should be clear enough?
		xidx := (node -1)*ncjt
		yidx := (node -1)*ncjt + 1
		zidx := (node -1)*ncjt + 2
		xridx := (node -1)*ncjt + 3
		yridx := (node -1)*ncjt + 4
		zridx := (node -1)*ncjt + 5
		xdis := dglb[xidx]
		ydis := dglb[yidx]
		zdis := dglb[zidx]
		xrot := dglb[xridx]
		yrot := dglb[yridx]
		zrot := dglb[zridx]
		js[node].Displ = []float64{xdis, ydis, zdis, xrot, yrot, zrot}
		//rx is the axial force, mx is torsion. TORSION.
		var rx, ry, rz, mx, my, mz float64
		//this is ..pathetic
		if nsc[xidx] > ndof {
			rx = rnode[nsc[xidx]-ndof-1]
		}
		if nsc[yidx] > ndof {
			ry =  rnode[nsc[yidx]-ndof-1]
		}

		if nsc[zidx] > ndof {
			rz =  rnode[nsc[zidx]-ndof-1]
		}
		if nsc[xridx] > ndof {
			mx = rnode[nsc[xridx]-ndof-1]
		}
		if nsc[yridx] > ndof {
			my = rnode[nsc[yridx]-ndof-1]
		}
		if nsc[zridx] > ndof {
			mz = rnode[nsc[zridx]-ndof-1]
		}
		js[node].React = []float64{rx, ry, rz, mx, my, mz}
	}	
	frmrez[0] = js
	frmrez[1] = ms
	frmrez[2] = dglb
	frmrez[3] = rnode
	frmrez[4] = msqf
	frmrez[5] = msloaded
	frmrez[6] = ""
	if !mod.Noprnt{frmrez[6] = Frm3dTable(js, ms, dglb, rnode, nsc, ndof, ncjt)}
	mod.Ncjt = ncjt
	mod.Frmtyp = 6
	mod.Frmstr = "3df"
	return frmrez, nil
}

//KmemF3d gets individual member stiffness matrices for members of a 3d frame model 
func KmemF3d(member int, memprp []int, nsc []int, js map[int]*Node, em [][]float64, cp [][]float64, wng [][]float64, ndof int, ncjt int, kchn chan []interface{}) {
	//sc numbers in nsc are y disp, rot x, rot z.
	//0 - member, 1 - cx , cy, cz, l, e, gu, iz, jv
	//2 - bkmat (local k mat)
	//3 - tmat, 4 - gkmat, 5 - srez
	var bkarr, tarr []float64
	rez := make([]interface{}, 9)
	jb := memprp[0]
	je := memprp[1]
	e := em[memprp[2]-1][0]
	gu := em[memprp[2]-1][1]
	a := cp[memprp[3]-1][0]
	iz := cp[memprp[3]-1][1]
	iy := cp[memprp[3]-1][2]
	jv := cp[memprp[3]-1][3]
	memrel := memprp[4]
	pa := js[jb].Coords
	pb := js[je].Coords
	xa := pa[0]
	ya := pa[1]
	za := pa[2]
	xb := pb[0]
	yb := pb[1]
	zb := pb[2]
	//wtyp 0-vertical, 1-others; wang - angle of roll in degrees (convert to rad)
	//assume Wng is in the same order as member order
	//else tis fubar
	wtyp := wng[member-1][0]
	wang := wng[member-1][1] * math.Pi / 180.0

	var mtyp int
	
	rez[0] = member
	//transformation matrix tmat = diag [r 0 0 0 ]
	var rxx, rxy, rxz, ryx, ryy, ryz, rzx, rzy, rzz, rden float64
	l := math.Sqrt(math.Pow(xb-xa, 2) + math.Pow(yb-ya, 2) + math.Pow(zb-za, 2))
	rxx = (xb - xa) / l
	rxy = (yb - ya) / l
	rxz = (zb - za) / l
	rden = math.Sqrt(math.Pow(rxx, 2) + math.Pow(rxz, 2))
	switch int(wtyp) {
	case 2: //general orientation
		ryx = (-rxx*rxy*math.Cos(wang) - rxz*math.Sin(wang)) / rden
		ryy = rden * math.Cos(wang)
		ryz = (-rxy*rxz*math.Cos(wang) + rxx*math.Sin(wang)) / rden
		rzx = (rxx*rxy*math.Sin(wang) - rxz*math.Cos(wang)) / rden
		rzy = -rden * math.Sin(wang)
		rzz = (rxy*rxz*math.Sin(wang) + rxx*math.Cos(wang)) / rden
		mtyp = 3
	case 1: //vertical members
		rxx = 0
		//rxy = rxy
		rxz = 0
		ryx = -rxy * math.Cos(wang)
		ryy = 0
		ryz = math.Sin(wang)
		rzx = rxy * math.Sin(wang)
		rzy = 0
		rzz = math.Cos(wang)
		mtyp = 1
	case 0:
		mtyp = 2
	}
	rez[1] = []float64{l, e, a, iz, iy, jv, rxx, rxy, rxz, wang, wtyp, gu}
	if wtyp == 1 || wtyp == 2 {
		tarr = []float64{
			rxx, rxy, rxz, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			ryx, ryy, ryz, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			rzx, rzy, rzz, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, rxx, rxy, rxz, 0, 0, 0, 0, 0, 0,
			0, 0, 0, ryx, ryy, ryz, 0, 0, 0, 0, 0, 0,
			0, 0, 0, rzx, rzy, rzz, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, rxx, rxy, rxz, 0, 0, 0,
			0, 0, 0, 0, 0, 0, ryx, ryy, ryz, 0, 0, 0,
			0, 0, 0, 0, 0, 0, rzx, rzy, rzz, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, rxx, rxy, rxz,
			0, 0, 0, 0, 0, 0, 0, 0, 0, ryx, ryy, ryz,
			0, 0, 0, 0, 0, 0, 0, 0, 0, rzx, rzy, rzz,
		}
	} else { //horizontal mem- identity matrix
		tarr = []float64{
			1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		}
	}
	tmat := mat.NewDense(2*ncjt, 2*ncjt, tarr)

	//set 12x12 member local stiffness matrix fuck
	zka := e / math.Pow(l, 3)
	zkb := gu * jv / l
	zkc := e * a / l
	switch memrel {
	case 0: //rigidly connected at both ends
		bkarr = []float64{
			zkc, 0, 0, 0, 0, 0, -zkc, 0, 0, 0, 0, 0,
			0, 12 * iz * zka, 0, 0, 0, 6 * l * iz * zka, 0, -12 * iz * zka, 0, 0, 0, 6 * l * iz * zka,
			0, 0, 12 * iy * zka, 0, -6 * l * iy * zka, 0, 0, 0, -12 * iy * zka, 0, -6 * l * iy * zka, 0,
			0, 0, 0, zkb, 0, 0, 0, 0, 0, -zkb, 0, 0,
			0, 0, -6 * l * iy * zka, 0, 4 * l * l * iy * zka, 0, 0, 0, 6 * l * iy * zka, 0, 2 * l * l * iy * zka, 0,
			0, 6 * l * iz * zka, 0, 0, 0, 4 * l * l * iz * zka, 0, -6 * l * iz * zka, 0, 0, 0, 2 * l * l * iz * zka,
			-zkc, 0, 0, 0, 0, 0, zkc, 0, 0, 0, 0, 0,
			0, -12 * iz * zka, 0, 0, 0, -6 * l * iz * zka, 0, 12 * iz * zka, 0, 0, 0, -6 * l * iz * zka,
			0, 0, -12 * iy * zka, 0, 6 * l * iy * zka, 0, 0, 0, 12 * iy * zka, 0, 6 * l * iy * zka, 0,
			0, 0, 0, -zkb, 0, 0, 0, 0, 0, zkb, 0, 0,
			0, 0, -6 * l * iy * zka, 0, 2 * l * l * iy * zka, 0, 0, 0, 6 * l * iy * zka, 0, 4 * l * l * iy * zka, 0,
			0, 6 * l * iz * zka, 0, 0, 0, 2 * l * l * iz * zka, 0, -6 * l * iz * zka, 0, 0, 0, 4 * l * l * iz * zka,
		}
	case 1: //hinge at beginning
		bkarr = []float64{
			zkc, 0, 0, 0, 0, 0, -zkc, 0, 0, 0, 0, 0,
			0, 3 * iz * zka, 0, 0, 0, 0, 0, -3 * iz * zka, 0, 0, 0, 3 * l * iz * zka,
			0, 0, 3 * iy * zka, 0, 0, 0, 0, 0, -3 * iy * zka, 0, -3 * l * iy * zka, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			-zkc, 0, 0, 0, 0, 0, zkc, 0, 0, 0, 0, 0,
			0, -3 * iz * zka, 0, 0, 0, 0, 0, 3 * iz * zka, 0, 0, 0, -3 * l * iz * zka,
			0, 0, -3 * iy * zka, 0, 0, 0, 0, 0, 3 * iy * zka, 0, 3 * l * iy * zka, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, -3 * l * iy * zka, 0, 0, 0, 0, 0, 3 * l * iy * zka, 0, 3 * l * l * iy * zka, 0,
			0, 3 * l * iz * zka, 0, 0, 0, 0, 0, -3 * l * iz * zka, 0, 0, 0, 3 * l * l * iz * zka,
		}
	case 2: //hinge at end
		bkarr = []float64{
			zkc, 0, 0, 0, 0, 0, -zkc, 0, 0, 0, 0, 0,
			0, 3 * iz * zka, 0, 0, 0, 3 * l * iz * zka, 0, -3 * iz * zka, 0, 0, 0, 0,
			0, 0, 3 * iy * zka, 0, -3 * l * iy * zka, 0, 0, 0, -3 * iy * zka, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, -3 * l * iy * zka, 0, -3 * l * l * iy * zka, 0, 0, 0, 3 * l * iy * zka, 0, 0, 0,
			0, 3 * l * iz * zka, 0, 0, 0, 3 * l * l * iz * zka, 0, -3 * l * iz * zka, 0, 0, 0, 0,
			-zkc, 0, 0, 0, 0, 0, zkc, 0, 0, 0, 0, 0,
			0, -3 * iz * zka, 0, 0, 0, -3 * l * iz * zka, 0, 3 * iz * zka, 0, 0, 0, 0,
			0, 0, -3 * iy * zka, 0, 3 * l * iy * zka, 0, 0, 0, 3 * iy * zka, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}
	case 3: //hinge at both ends, zero matrix like a n00b
		bkarr = []float64{
			zkc, 0, 0, 0, 0, 0, -zkc, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			-zkc, 0, 0, 0, 0, 0, zkc, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}
	}
	bkmat := mat.NewDense(2*ncjt, 2*ncjt, bkarr)

	//member global stiffness matrix gk = t.T @ k @ t
	gkmat := mat.NewDense(2*ncjt, 2*ncjt, nil)
	if wtyp == 0 {
		gkmat = bkmat
	} else {
		gkmat.Product(tmat.T(), bkmat, tmat)
	}

	rez[2] = bkmat
	rez[3] = tmat
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
					if r != 0.0 {
						srez = append(srez, []float64{float64(na), float64(nb), r})
					}
				}
			}
		}
	}
	rez[5] = srez
	rez[6] = mtyp

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

//KassF3d assembles structure stiffness matrix from rezchan

func KassF3d(ndof int, nms int, rezchn chan [][]float64, symchn chan *mat.SymDense) {
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
