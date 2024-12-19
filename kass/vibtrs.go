package barf

import (
	"log"
	"fmt"
	//"math"
	//"math/cmplx"
	//"errors"
	"gonum.org/v1/gonum/mat"
)

//VibTrs2d performs vibrational analysis of a 2d truss
//as seen in VIBPT in weaver
//it is BRAKE/NOT WERK
func VibTrs2d(mod *Model){
	/*
	   vibpt from weaver 
	*/
	ncjt := 2
	frmrez, err := CalcTrs(mod, ncjt)
	if err != nil{
		log.Println(err)
	}
	js, _ := frmrez[0].(map[int]*Node)
	ms, _ := frmrez[1].(map[int]*Mem)
	nsc, _ := frmrez[4].([]int)
	nj := len(js)
	ndof, _ := frmrez[5].(int)
	schmat, _ := frmrez[9].(mat.Cholesky)
	log.Println(frmrez[6])
	var u mat.TriDense
	schmat.UTo(&u)
	log.Println("umat->")
	fmt.Println(mat.Formatted(&u))
	var na, nb int
	nscis := make(map[int]int)
	for i, sc := range nsc{
		nscis[sc] = i
	}
	mumat := mat.NewDense(nj*ncjt, nj*ncjt, nil)
	ssmat := mat.NewDense(nj*ncjt, nj*ncjt, nil)
	mmat := mat.NewDense(ndof, ndof, nil)
	//build overall nj * ncjt matrix
	for _, mem := range ms{
		l, a := mem.Geoms[0],mem.Geoms[2]
		pg := 1.0 //SO WRONG SO FUCKING WRONG
		mem.Mmat = mat.NewDense(2*ncjt, 2*ncjt, []float64{
			2,0,1,0,
			0,2,0,1,
			1,0,2,0,
			0,1,0,2,
		})
		mem.Mmat.Scale(pg * a * l/6.0,mem.Mmat)
		jb := js[mem.Mprp[0]]; je := js[mem.Mprp[1]]
		for i := 0; i < 2*ncjt; i++ {
			if i < ncjt {
				na = jb.Nscs[i]
			} else {
				na = je.Nscs[i-ncjt]
			}
			for j := 0; j < 2*ncjt; j++{
				if j < ncjt {
					nb = jb.Nscs[j]
				} else {
					nb = je.Nscs[j-ncjt]
				}
				ni, nj := nscis[na], nscis[nb]
				s0 := ssmat.At(ni, nj); s1 := mem.Gkmat.At(i,j)
				if s0 + s1 != 0.0{
					ssmat.Set(ni, nj, s0 + s1)
				}
				m0 := mumat.At(na-1, nb-1); m1 := mem.Mmat.At(i,j)
				if m0 + m1 != 0.0{
					mumat.Set(na-1, nb-1, m0 + m1)	
				}
				if na <= ndof && nb <= ndof{
					m0 := mmat.At(na-1, nb-1); m1 := mem.Mmat.At(i,j)
					if m0 + m1 != 0.0{mmat.Set(na-1, nb-1, m0+m1)}
				}
			}
		}
	}
	log.Println("mumat->")
	fmt.Println(mat.Formatted(mumat))
	log.Println("ssmat->")
	fmt.Println(mat.Formatted(ssmat))
	log.Println("schmat->")
	fmt.Println(mat.Formatted(&schmat))
	log.Println("mmat->")
	fmt.Println(mat.Formatted(mmat))
	log.Println("nsc->",nsc)
	
	var ssym mat.SymDense
	ssym.SymOuterK(1, ssmat)

	fmt.Printf("ssym = %0.4v\n", mat.Formatted(&ssym, mat.Prefix("    ")))
	var chol mat.Cholesky
	if ok := chol.Factorize(&ssym); !ok {
		log.Println("ssmat is not positive definite")
		var m mat.SymDense
		m.SymOuterK(1,mumat)
		if ok := chol.Factorize(&m); !ok{
			log.Println("mumat is not positive definite")
		} else {
			log.Println("***YOOO***")
		}
	}
	var eig mat.Eigen
	ok := eig.Factorize(&ssym, mat.EigenRight)
	if !ok {
		log.Fatal("Eigen decomposition of ssmat failed")
	}
	ok = eig.Factorize(&schmat, mat.EigenLeft)
	if !ok {
		log.Fatal("Eigen decomposition of schmat failed")
	}
	fmt.Printf("Eigenvalues of schmat:\n%v\n", eig.Values(nil))
	log.Println("inverting u-")
	//var uinv, utinv *mat.Dense
	uinv := mat.NewDense(ndof, ndof, nil)
	utinv := mat.NewDense(ndof, ndof, nil)
	err = uinv.Inverse(&u)
	if err != nil{
		log.Println(err)
	}
	err = utinv.Inverse(u.T())
	if err != nil{
		log.Println(err)
	}
	mv := mat.NewDense(ndof, ndof, nil)
	mv.Product(uinv.T(), mmat)
	mv.Product(mv, uinv)
	ok = eig.Factorize(mv, mat.EigenRight)
	if !ok {
		log.Fatal("Eigen decomposition of schmat failed")
	}
	fmt.Printf("Eigenvalues of mv:\n%v\n", eig.Values(nil))
	log.Println("***now basic dumb way***")
	kinv := mat.NewSymDense(ndof, nil)
	err = schmat.InverseTo(kinv)
	if err != nil{
		log.Println(err)
	} else{
		log.Println("calculating eigen values")
		var a mat.Dense
		a.Mul(kinv, mmat)	
		ok = eig.Factorize(&a, mat.EigenRight)
		if !ok {
			log.Fatal("Eigen decomposition of schmat failed")
		}else{
			fmt.Printf("Eigenvalues of mv:\n%v\n", eig.Values(nil))
		}
	}
}

/*
   #scratch

   	//pltchn := make(chan string)
	//go PlotTrs2d(mod, "dumb", pltchn)
	//pltstr := <- pltchn
	//log.Println(pltstr)
	//log.Println("kass smat->",ColorRed)
	//fmt.Println(mat.Formatted(&schmat),ColorReset)
	//mv := mat.NewDense(nj*ncjt, nj*ncjt, nil)
	//fmt.Printf("A = %v\n\n", mat.Formatted(a, mat.Prefix("    ")))
		//fmt.Println("mem->",mem.Id)
		//fmt.Println(mat.Formatted(mem.Gkmat))
		//fmt.Println(mat.Formatted(mem.Bkmat))


*/
