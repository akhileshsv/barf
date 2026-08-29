package barf

import (
	"log"
	"math"
	//"errors"
	"gonum.org/v1/gonum/mat"
)

//NlCalcTrs2d performs non linear analysis of a 2d truss model
//see kassimali chapter 10
func NlCalcTrs2d(mod *Model,tol float64) {
	if tol == 0.0 {tol = 0.001}
	frmrez, err := CalcTrs(mod, mod.Ncjt)
	if err != nil{
		log.Println(err)
	}
	log.Println(frmrez[6])
	js, _ := frmrez[0].(map[int]*Node)
	ms, _ := frmrez[1].(map[int]*Mem)
	//dglb, _ := frmrez[2].([]float64)
	ndof, _ := frmrez[5].(int)
	p, _ := frmrez[7].([]float64)
	d0, _ := frmrez[8].([]float64)
	//log.Println("pvec->",p)
	//log.Println(d0)
	//var pltchn chan string
	//go PlotTrs2d(mod, "dumb", pltchn)
	var iter, kiter int
	dn := make([]float64, len(d0))
	copy(dn,d0)
	//pltstr := <- pltchn
	//log.Println(pltstr)
	for iter != -1{	
		kiter++
		var dsum, difsum float64
		if kiter > 100{
			log.Println("ERRORE,errore->maximum iteration limit reached")
			break
		}
		fchn, schn, yochn := make(chan []rez1d, len(ms)), make(chan []rez2d, len(ms)), make(chan int, len(ms))
		for i, mem := range ms{
			jb, je := js[mem.Mprp[0]], js[mem.Mprp[1]]
			go NlKmem2dT(i, ndof, mem, jb, je, dn, fchn, schn, yochn)
		}
		for i:=0; i< len(ms); i++{
			<- yochn
		}
		close(fchn); close(schn)
		umat := Uass2dT(p, fchn)
		stmat := Stass2dT(ndof, schn)
		//fc := mat.Formatted(umat, mat.Prefix("    "), mat.Squeeze())
		//log.Println("umat\n",fc)
		//fc = mat.Formatted(stmat, mat.Prefix("    "), mat.Squeeze())
		//log.Println("stmat\n", fc)
		var dmat mat.Dense
		err := dmat.Solve(stmat, umat)
		if err != nil{
			log.Println("ERRORE, errore-> matrix solution error")
			break
		}
		//log.Println("dmat->\n",mat.Formatted(&dmat))
		for i := 0; i < ndof; i++{
			di := dmat.At(i,0)
			dsum += math.Pow(dn[i],2)
			difsum += math.Pow(di,2)
			dn[i] = dn[i] + di
		}
		//log.Println("dnew->",dn)
		if math.Sqrt(difsum/dsum) > tol{
			log.Println(kiter, "rms rat->",math.Sqrt(difsum/dsum))
			continue
		} 
		log.Println("iteration converged", difsum, dsum)
		log.Println("dvec->",dn)
		iter = -1
	}
}

//Uass2dT assembles the local force matrix umat
func Uass2dT(p []float64, fchn chan []rez1d) (umat *mat.Dense){
	f1 := make([]float64, len(p))
	copy(f1,p)
	for rez := range fchn{
		for _, r := range rez{
			f1[r.i] -= r.val
		}
	}
	umat = mat.NewDense(len(p), 1, f1)
	return
}

//Stass2dT assembles the structure tangent stiffness matrix for a 2d truss
func Stass2dT(ndof int, schn chan []rez2d) (stmat *mat.Dense){
	stmat = mat.NewDense(ndof, ndof, nil)
	for rez := range schn{
		for _, r := range rez{
			val := stmat.At(r.i-1,r.j-1) + r.val
			stmat.Set(r.i-1,r.j-1,val) 
		}
	}
	return
}

//NlKmem2dT generates individual member matrices for a 2d non linear truss member
func NlKmem2dT(i, ndof int, mem *Mem, jb, je *Node, dn []float64, fchn chan []rez1d, schn chan []rez2d, yochn chan int){
	ncjt := 2
	l0, e, a := mem.Geoms[0], mem.Geoms[1], mem.Geoms[2]
	xb, yb, xe, ye := jb.Coords[0], jb.Coords[1], je.Coords[0], je.Coords[1]
	nscs := []int{jb.Nscs[0],jb.Nscs[1],je.Nscs[0],je.Nscs[1]}
	vs := make([]float64,4)
	for i, sc := range nscs{
		if sc <= ndof{
			vs[i] = dn[sc-1]
		}
	}
	v1, v2, v3, v4 := vs[0], vs[1], vs[2], vs[3]
	//log.Println(i, "v1 init->",v1, v2, v3, v4)
	//log.Println(i, "v1 js->", jb.Displ[0], jb.Displ[1],je.Displ[0],je.Displ[1])
	l1 := math.Sqrt(math.Pow((xe + v3)-(xb + v1),2) + math.Pow((ye + v4)-(yb + v2),2))
	cx := (xe + v3 - xb - v1)/l1
	cy := (ye + v4 - yb - v2)/l1
	u := l0 - l1
	q := e * a * u/l0
	fvec := []float64{cx * q, cy * q, -cx * q, -cy * q}
	gmat := mat.NewDense(2*ncjt, 2*ncjt, []float64{
		-math.Pow(cy,2), cx * cy, math.Pow(cy,2), - cx * cy,
		 cx * cy, -math.Pow(cx,2), - cx * cy, math.Pow(cx,2),
		 math.Pow(cy,2), -cx * cy, -math.Pow(cy,2), cx * cy,
		-cx * cy, math.Pow(cx,2),  cx * cy, -math.Pow(cx, 2),
	})
	gmat.Scale(q/l1,gmat)
	tmat := mat.NewDense(1,4,[]float64{cx, cy, -cx, -cy})
	ktmat := mat.NewDense(2*ncjt, 2*ncjt, nil)
	ktmat.Product(tmat.T(),tmat)
	ktmat.Scale(e*a/l0, ktmat)
	ktmat.Add(ktmat, gmat)
	//log.Println("kmem deetz->")
	//log.Println(i,"\n",mat.Formatted(ktmat))
	//log.Println(i,fvec)
	var na, nb int
	var srez []rez2d
	for i := 0; i < 2*ncjt; i++ {
		if i < ncjt {
			na = jb.Nscs[i]
		} else {
			na = je.Nscs[i-ncjt]
		}
		if na <= ndof {
			for j := 0; j < 2*ncjt; j++ {
				if j < ncjt {
					nb = jb.Nscs[j]
				} else {
					nb = je.Nscs[j-ncjt]
				}
				if nb <= ndof {
					val := ktmat.At(i, j)
					if val != 0.0 {
						srez = append(srez, rez2d{na,nb,val})
					}
				}
			}

		}
	}
	var frez []rez1d
	for i, sc := range nscs{
		if sc <= ndof{
			if fvec[i] != 0.0{
				frez = append(frez, rez1d{sc -1, fvec[i]})
			}
		}
	}
	schn <- srez
	fchn <- frez
	yochn <- 1
}
