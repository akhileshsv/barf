package barf

import (
	"fmt"
	"log"
	"math"
	"gonum.org/v1/gonum/mat"
)


func NlCalcFrm2d(mod *Model,tol float64) {
	//non linear frame analysis (p-delta)
	if tol == 0.0 {tol = 0.001}
	frmrez, err := CalcFrm2d(mod, mod.Ncjt)
	if err != nil{
		log.Println(err)
	}
	log.Println(frmrez[6])
	log.Println("hear goeth nothingwa")
	js, _ := frmrez[0].(map[int]*Node)
	ms, _ := frmrez[1].(map[int]*Mem)
	ndof := mod.Ndof
	p, _ := frmrez[7].([]float64)
	d0, _ := frmrez[8].([]float64)
	log.Println("pvec->",p)
	log.Println(d0)
	log.Println("NDOF-",ndof)
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
		if kiter > 666{
			log.Println("ERRORE,errore->maximum iteration limit reached")
			break
		}
		fchn, schn, yochn := make(chan []rez1d, len(ms)), make(chan []rez2d, len(ms)), make(chan int, len(ms))
		for i, mem := range ms{
			// go NlKmem2df(i, ndof, mem, jb, je, dn, dchn, schn, yochn)
			jb, je := js[mem.Mprp[0]], js[mem.Mprp[1]]
			memrel := 0
			go NlKmem2df(i, ndof, memrel, mem, jb, je, dn, fchn, schn, yochn)
		}		
		//yochn is the signal channel (1=='yo')
		for i:=0; i< len(ms); i++{
			<- yochn
		}
		close(fchn); close(schn)
		umat := UassF2d(p, fchn)
		stmat := StassF2d(ndof, schn)
		fc := mat.Formatted(umat, mat.Prefix(" "), mat.Squeeze())
		log.Println("umat->\n",fc)
		fc = mat.Formatted(stmat, mat.Prefix(" "), mat.Squeeze())
		log.Println("stmat->\n", fc)
		var dmat mat.Dense
		err := dmat.Solve(stmat, umat)
		if err != nil{
			log.Println("ERRORE, errore-> matrix solution error")
			break
		}
		log.Println("dmat->\n",mat.Formatted(&dmat))
		for i := 0; i < ndof; i++{
			di := dmat.At(i,0)
			dsum += math.Pow(dn[i],2)
			difsum += math.Pow(di,2)
			dn[i] = dn[i] + di
		}
		log.Println("dnew->",dn)
		if math.Sqrt(difsum/dsum) > tol{
			log.Println(kiter, "rms rat->",math.Sqrt(difsum/dsum))
			continue
		} 
		log.Println("iteration converged", difsum, dsum)
		log.Println("dvec->",dn)
		iter = -1
	}
}


//Sf2dKima calcs majid stability funcs o1-o7
func Sf2dKima(pu, em, iz, lspan, ar float64)(o1, o2, o3, o4, o5, o6, o7 float64){
	fmt.Println("PU",pu)
	pelr := math.Pow(math.Pi, 2) * em * iz / math.Pow(lspan, 2)
	pr := pu/pelr
	if pr < 0.001{
		o1 = 1.0
		o2 = 1.0
		o3 = 1.0
		o4 = 1.0
		o5 = 1.0
		o6 = 1.0
		o7 = 1.0
		return
	}

	al := 0.5 * math.Pi * math.Sqrt(pr)
	a1 := 1.57973627
	a2 := 0.15858587
	a3 := 0.02748899
	a4 := 0.00547540
	a5 := 0.00115281
	a6 := 0.00024908
	a7 := 0.00005452
	afacs := []float64{a1,a2,a3,a4,a5,a6,a7}
	o1 = (64.0 - 60.0 * pr + 5.0 * pr * pr)/((16.0 - pr)*(4.0 - pr))
	for idx, af := range afacs{
		i := float64(idx+1)
		o1 -= math.Pow(af, 2) * math.Pow(pr, i)/math.Pow(2.0, 3.0 * i)	
	}
	//o1 = al/math.Tan(al)
	o2 = math.Pow(al, 2)/(1.0 - o1)/3.0
	o3 = (3.0 * o2 + o1)/4.0
	o4 = (3.0 * o2 - o1)/4.0
	o5 = o1 * o2
	o6 = o3/o2/(2.0 * o3 - o4)
	o7 = o4 * o6/o3
	return
}


//Sf2dKass calculates kassimali stability and bowing functions
func Sf2dKass(memrel int, mpc, q0, t1, t2, delt, srat float64)(qf, c1, c2, cb, b1, b2, b11, b21, cb1 float64){
	var t11, t21 float64
	var kiter, iter int
	qi := q0
	tol := 0.1
	for iter != -1{
		//fmt.Println("qi, kiter",qi,kiter)
		kiter++
		if kiter > 60{
			log.Println(ColorRed, "iteration error", ColorReset)
			return
		}
		switch {
		case math.Abs(qi) <= 2.0:
			//use series expressions
			c1 = 4.0 - 2.0* math.Pow(math.Pi, 2)* qi/15.0 - 11.0 * math.Pow(math.Pi, 4) * math.Pow(qi, 2.0)/6300.0 - math.Pow(math.Pi, 6) * math.Pow(qi, 3)/27000.0
			c2 = 2.0 + math.Pow(math.Pi, 2) * qi/30.0 + 13.0 * math.Pow(math.Pi, 4) * math.Pow(qi, 2.0)/12600.0 + 11.0 * math.Pow(math.Pi, 6) * math.Pow(qi, 3)/378000.0
			b1 = 1.0/40.0 + math.Pow(math.Pi, 2) * qi/2800.0 + math.Pow(math.Pi, 4) * math.Pow(qi, 2.0)/168000.0 + 37.0 * math.Pow(math.Pi, 6) * math.Pow(qi, 3)/388080000.0
			b2 = 1.0/24.0 + math.Pow(math.Pi, 2) * qi/720.0 + math.Pow(math.Pi, 4) * math.Pow(qi, 2.0)/20160.0 + math.Pow(math.Pi, 6) * math.Pow(qi, 3)/604800.0
			cb = b1 * math.Pow(t1 + t2, 2.0) + b2 * math.Pow(t1 - t2,2.0)
			b11 = math.Pow(math.Pi, 2)/2800.0 + math.Pow(math.Pi, 4) * qi/84000.0 + 37.0 * math.Pow(math.Pi, 6) * math.Pow(qi, 2)/129360000.0
			b21 = math.Pow(math.Pi, 2)/720.0 + math.Pow(math.Pi, 4) * qi/10080.0 + math.Pow(math.Pi, 6) * math.Pow(qi, 2)/201600.0
		case qi > 0.0:
			thet := math.Sqrt(math.Pi * math.Pi * qi)
			c1 = thet * math.Sin(thet) - math.Pow(thet, 2) * math.Cos(thet)
			c1 = c1/(2.0 - 2.0 * math.Cos(thet) - thet * math.Sin(thet))
			c2 = math.Pow(thet, 2) - thet * math.Sin(thet)
			c2 = c2/(2.0 - 2.0 * math.Cos(thet) - thet * math.Sin(thet))
			b1 = (c1 + c2) * (c2 - 2.0)/(8.0 * math.Pi * math.Pi * qi)
			b2 = c2/(8.0 * (c1 + c2))
			cb = b1 * math.Pow(t1 + t2, 2.0) + b2 * math.Pow(t1 - t2,2.0)
			b11 = ((b1 - b2) * (c1 + c2) - 2.0 * c2 * b1)/4.0/qi
			b21 = math.Pi * math.Pi * (16.0 * b1 * b2 - b1 + b2)/4.0/(c1 + c2)
		
		case qi < 0.0:
			thet := math.Sqrt(-math.Pi * math.Pi * qi)
			c1 = math.Pow(thet, 2) * math.Cosh(thet) - thet * math.Sinh(thet)
			c1 = c1/(2.0 - 2.0 * math.Cosh(thet) - thet * math.Sinh(thet))
			c2 = thet * math.Sinh(thet) - math.Pow(thet, 2)
			c2 = c2/((2.0 - 2.0 * math.Cosh(thet) - thet * math.Sinh(thet)))
			b1 = (c1 + c2) * (c2 - 2.0)/(8.0 * math.Pi * math.Pi * qi)
			b2 = c2/(8.0 * (c1 + c2))
			cb = b1 * math.Pow(t1 + t2, 2.0) + b2 * math.Pow(t1 - t2,2.0)
			b11 = ((b1 - b2) * (c1 + c2) - 2.0 * c2 * b1)/4.0/qi
			b21 = math.Pi * math.Pi * (16.0 * b1 * b2 - b1 + b2)/4.0/(c1 + c2)
		}
		cb1 = b11 * math.Pow(t1 + t2, 2) + b21 * math.Pow(t1 - t2, 2) + 2.0 * b1 * (t1 + t2) * (t11 + t21) + 2.0 * b2 * (t1-t2) * (t11- t21)
		kqi := math.Pow(math.Pi, 2.0) * qi/math.Pow(srat, 2.0) + cb - delt
		kqi1 := math.Pow(math.Pi, 2.0)/math.Pow(srat, 2.0) + cb1
		fmt.Println(ColorYellow,"cb1, kqi,kqil, kqirat, qi, qn",cb1,kqi, kqi1,kqi/kqi1,qi,qi-kqi/kqi1,ColorReset)
		qn := qi - kqi/kqi1
		if math.Abs(qn - qi) <= tol || kiter == 1{
			qf = qi
			iter = -1
			break
		} else {
			//q0 = qi
			qi = qn
		}
	}
	fmt.Println(ColorCyan, "qi final-",qf,ColorReset)
	return
}


//NlKmem2df generates individual non linear analysis member matrices for a 2d frame member
func NlKmem2df(i, ndof, memrel int, mem *Mem, jb, je *Node, dn []float64, fchn chan []rez1d, schn chan []rez2d, yochn chan int){
	ncjt := 3
	lspan, em, ar, iz := mem.Geoms[0], mem.Geoms[1], mem.Geoms[2], mem.Geoms[3]
	x1, y1, x2, y2 := jb.Coords[0], jb.Coords[1], je.Coords[0], je.Coords[1]
	nscs := []int{jb.Nscs[0],jb.Nscs[1],jb.Nscs[2],je.Nscs[0],je.Nscs[1],je.Nscs[2]}
	vs := make([]float64,6)
	for i, sc := range nscs{
		if sc <= ndof{
			vs[i] = dn[sc-1]
		}
	}
	//global displacements in x and y
	v1, v2, v3, v4, v5, v6 := vs[0], vs[1], vs[2], vs[3], vs[4], vs[5]
	//angle in undeformed config (init angle)
	aud := math.Atan((y2-y1)/(x2-x1))
	//def. angle
	tade := (y2 + v5 - y1 -  v2)/(x2 + v4 - x1 - v1)
	ade := math.Atan(tade)
	pde := ade - aud
	//deformed length
	ldef := math.Sqrt(math.Pow(x2 +  v4 - x1 -  v1,2)+math.Pow(y2 +  v5 - y1 -  v2,2))
	u3 := lspan - ldef
	delt := u3/lspan
	m := math.Cos(ade)
	n := math.Sin(ade)
	t1 := v3 - pde
	t2 := v6 - pde
	qaxi := em * ar * (u3/lspan)
	qi := qaxi * lspan * lspan/math.Pi/math.Pi/em/iz
	srat := lspan/math.Sqrt(iz/ar)
	//fmt.Println("srat",srat,"qi",qi)
	//FIGURE PLASTIC MOMENT HERE
	mpc := 0.0
	//fmt.Println("memrel, mpc, qi, t1, t2, delt, srat",memrel, mpc, qi, t1, t2, delt, srat)
	_, c1, c2, cb, b1, b2, b11, b21, _ := Sf2dKass(memrel, mpc, qi, t1, t2, delt, srat)
	
	//log.Println("qf, c1, c2, cb-", qf, c1, c2, cb)

	tmat := mat.NewDense(6, 3, []float64{
		-n, -n, m*ldef,
		m, m, n*ldef,
		ldef, 0, 0,
		n, n, -m*ldef,
		-m, -m, -n*ldef,
		0, ldef, 0,
	})
	tmat.Scale(1.0/ldef,tmat)
	m1 := em * iz * (c1 * t1 + c2 * t2)/ldef
	m2 := em * iz * (c2 * t1 + c1 * t2)/ldef
	qax := em * ar * (u3/lspan - cb)
	lqmat := mat.NewDense(3, 1, []float64{m1, m2, qax})
	fmat := mat.NewDense(6, 1, nil)
	fmat.Product(tmat, lqmat)
	
	g1 := []float64{
		-2.0*m*n, m*m - n*n, 0, 2.0*m*n, -(m*m - n*n), 0,
		m*m - n*n, 2.0*m*n, 0, -(m*m - n*n), -2.0*m*n, 0,
		0, 0, 0, 0, 0, 0,
		2.0*m*n, -(m*m - n*n), 0, -2.0*m*n, m*m - n*n,0,
		-(m*m - n*n), -2.0*m*n, 0, (m*m - n*n), 2.0*m*n,0,
		0, 0, 0, 0, 0, 0,
	}
	g3 := []float64{
		-n*n, m*n, 0, n*n, -m*n, 0,
		m*n, -m*m, 0, -m*n, m*m, 0,
		0, 0, 0, 0, 0, 0,
		n*n, -m*n, 0, -n*n, m*n, 0,
		-m*n, m*m, 0, m*n, -m*m, 0,
		0, 0, 0, 0, 0, 0,
	}
	g1mat := mat.NewDense(6,6, g1)
	g2mat := mat.NewDense(6,6, g1)
	g3mat := mat.NewDense(6,6, g3)
	g1mat.Scale(m1/math.Pow(lspan,2.0), g1mat)
	g2mat.Scale(m2/math.Pow(lspan,2.0), g2mat)
	g3mat.Scale(qax*lspan, g3mat)
	c11 := -2.0 * math.Pow(math.Pi, 2) * (b1 + b2)
	c21 := -2.0 * math.Pow(math.Pi, 2) * (b1 - b2)
	
	G1 := c11 * t1 + c21 * t2
	G2 := c21 * t1 + c11 * t2
	H := math.Pow(math.Pi,2)/math.Pow(srat, 2) + b11 * math.Pow(t1 + t2,2) + b21 * math.Pow(t1 - t2, 2)

	tsmat := mat.NewDense(3, 3, []float64{
		c1 + math.Pow(G1, 2)/math.Pow(math.Pi, 2)/H, c2 + G1 * G2/math.Pow(math.Pi, 2)/H, G1/H,
		c2 + G1 * G2/math.Pow(math.Pi, 2)/H, c1 + math.Pow(G2, 2)/math.Pow(math.Pi, 2)/H, G2/H,
		G1/H, G2/H, math.Pow(math.Pi, 2)/H,
	})
	tsmat.Scale(em*iz/lspan, tsmat)
	bkmat := mat.NewDense(6, 6, nil)
	bkmat.Product(tmat, tsmat, tmat.T())
	bkmat.Add(bkmat, g1mat)
	bkmat.Add(bkmat, g2mat)
	bkmat.Add(bkmat, g3mat)
	
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
					val := bkmat.At(i, j)
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
			if fmat.At(i,0) != 0.0{
				frez = append(frez, rez1d{sc -1, fmat.At(i,0)})
			}
		}
	}
	schn <- srez
	fchn <- frez
	yochn <- 1
}

//UassF2d assembles the unbalanced force vector
func UassF2d(p []float64, fchn chan []rez1d) (umat *mat.Dense){
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

//StassF2d assembles the system tangent stiffness matrix
func StassF2d(ndof int, schn chan []rez2d) (stmat *mat.Dense){
	stmat = mat.NewDense(ndof, ndof, nil)
	for rez := range schn{
		for _, r := range rez{
			val := stmat.At(r.i-1,r.j-1) + r.val
			stmat.Set(r.i-1,r.j-1,val) 
		}
	}
	return
}

//Sf2dKass calcs kassimali stability/bowing funcs c1, c2
func SfuncsF2d(p, em, iz, lspan, ar float64)(){
	return
}

