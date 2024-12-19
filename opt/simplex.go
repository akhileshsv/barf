package barf

import (
	"log"
	"math"
	"sort"
	"errors"
)

func unitvec(v1, v2 []float64) (vu []float64, vmod float64){
	//find the unit vector between two points
	vu = make([]float64, len(v1))
	for i := range v1{
		vu[i] = v2[i] - v1[i]
		vmod += math.Pow(v2[i] - v1[i],2)
	}
	vmod = math.Sqrt(vmod)
	if vmod == 0{log.Println("ERRORE,errore->unit vector modulus == 0");return}
	for i := range vu{
		vu[i] = vu[i]/vmod
	}
	return
}

func lerpvec(scale float64, v1, v2 []float64) (v3 []float64){
	vu, vmod := unitvec(v1,v2)
	v3 = make([]float64, len(v1))
	for i := range v1{
		v3[i] = v1[i] + scale * vmod * vu[i]
	}
	return

}

func centroidsimp(v [][]float64) (vc []float64){
	//get the centroid of n points, pass slice from [1:]
	vc = make([]float64, len(v[0]))
	for _, vi := range v{
		for j, val := range vi{
			vc [j] += val
		}
	}
	for i := range vc{
		vc[i] = vc[i]/float64(len(v))
	}
	return
}

func initsimp(x0 []float64, c float64) (sx [][]float64){
	//generate n simplex points
	//WHAT SHOULD C BE (just increment (+) or fraction?)
	sx = append(sx, x0)
	for i := range x0{
		vi := make([]float64, len(x0))
		copy(vi,x0)
		vi[i] = x0[i] + c 
		sx = append(sx, vi)
	}
	return
}

func fitsimp(sx [][]float64, fobj func([]float64) float64) (fval []float64){
	for _, v := range sx{
		fi := fobj(v)
		fval = append(fval, fi)
	}
	return
}

func sortsimp(sx [][]float64, fval []float64){
	sort.Slice(sx, func(i, j int) bool {
		return fval[i] > fval[j]
	})
	
	sort.Slice(fval, func(i, j int) bool {
		return fval[i] > fval[j]
	})
	return
}

func reflect(h, c []float64, al float64) (r []float64){
	//scale h (x0) and c (x1) by 2 * alpha
	scale := 2.0 * al
	r = lerpvec(scale, h, c)
	return
}

func expand(c, r []float64, ga float64) (e []float64){
	//xpand r in dir cr by gamma
	scale := ga
	e = lerpvec(scale, c, r)
	return
}

func contract(h, c []float64, be float64) (s []float64){
	//contract c towards h by beta
	scale := be
	s = lerpvec(scale, h, c)
	return
}


func shrink(sx [][]float64, de float64) (nsx [][]float64){
	//keep best point and scale rest towards vb
	vb := sx[len(sx)-1]
	for i, vi := range sx{
		if i == len(sx) -1{
			nsx = append(nsx, vb)
		} else {
			nvi := lerpvec(de, vi, vb) 
			nsx = append(nsx, nvi)
		}
	}
	return
}

func checkconv(fval []float64, ep float64) (bool){
	//check if (fi - favg)2 < ep
	var fsum, rms float64
	for _, fi := range fval{
		fsum += fi
	}
	favg := fsum/float64(len(fval))
	for _, fi := range fval{
		rms += math.Pow(fi - favg, 2.0)/float64(len(fval)-1)
	}
	rms = math.Sqrt(rms)
	return rms < ep
}

func NeldMead(niter int, fobj func([]float64) float64, x0 [][]float64, params []float64) (gv []float64, gbest float64, err error){
	//nelder-mead (simplex) optimization method
	//convergence is a bitch
	var al, be, ga, de, ep float64
	var kiter, ndim int
	al, be, ga, de, ep = 1.0, 0.5, 2.0, 0.5, 1e-20
	if len(params) == 5{
		al, be, ga, de, ep = params[0], params[1], params[2], params[3], params[4]
	}
	ndim = len(x0[0])
	//init
	var sx [][]float64
	if len(x0) == 1{
		//initialize with radius = 0.1
		r := 0.1
		if len(params) == 6{r = params[5]}
		sx = initsimp(x0[0], r)
	} else {
		if len(x0) != ndim + 1{
			err = errors.New("starting simplex dimension error")
			return
		}
		sx = x0
	}
	fval := fitsimp(sx, fobj)
	//simplex loop
	for{
		if niter > 0 && kiter > niter{
			log.Println("max iterations reached, stopping")
			break
		}
		//sort
		sortsimp(sx, fval)
		gbest = fval[len(fval)-1]
		//calc centroid
		vc := centroidsimp(sx[1:])
		var vn []float64
		//reflect
		vr := reflect(sx[0], vc, al)
		if fobj(vr) < fval[len(fval)-2]{
			ve := expand(vc, vr, ga)
			if fobj(ve) < fobj(vr){
				vn = make([]float64, ndim)
				copy(vn, ve)
			} else {
				vn = make([]float64, ndim)
				copy(vn, vr)
			}
		} else {
			vs := contract(sx[0], vc, be)
			if fobj(vs) < fval[0]{
				vn = make([]float64, ndim)
				copy(vn, vs)
			} 
		}
		if len(vn) == 0{
			//shrink
			sx = shrink(sx, de)
			fval = fitsimp(sx, fobj)
			continue
		} else {
			copy(sx[0], vn)
			fval = fitsimp(sx, fobj)
		}
		//check
		if checkconv(fval, ep){
			log.Println("convergence criteria reached, stopping")
			break
		}
		kiter++
	}
	//final values
	sortsimp(sx, fval)
	gv = sx[len(sx)-1]
	gbest = fval[len(fval)-1]

	//s. rao in opt says "the centroid maybe taken as the optimal value"?
	//gv = centroidsimp(sx)
	//gbest = fobj(gv)
	return
}

func razaobj(v []float64) (float64) {
	//https://www.youtube.com/watch?v=tjyVqZu1t-s&ab_channel=RazaLatif
	//ex 1 f(x,y) = x2 - 4x + y2 -y -xy
	x := v[0]; y := v[1]
	return math.Pow(x,2.0) - 4.0 * x + math.Pow(y, 2.0) - y - x * y
}

func parobj(v []float64) (float64){
	//y = x2, minimum at 0
	return math.Pow(v[0],2.0)
}

func rosenpval(v []float64) (float64){
	//rosenbrock's parabolic valley (lmao)
	return 100.0 * math.Pow(v[1] - math.Pow(v[0],2),2) + math.Pow(1.0 - v[0],2)
}

func powelquart(v []float64) (float64){
	//powell's quartic function
	x1, x2, x3, x4 := v[0], v[1], v[2], v[3]
	return math.Pow(x1 + 10.0 * x2,2.0) + 5.0 * math.Pow(x3-x4,2.0) + math.Pow(x2 - 2.0 * x3,4.0) + 10.0 * math.Pow(x1 - x4, 4.0)
}

func datechobj(v []float64) (float64) {
	//https://www.datatechnotes.com/2022/01/nelder-mead-optimization-example-in.html
	//y = x2+2*sin(pi*x)
	return math.Pow(v[0],2) + 2.0 * math.Sin(math.Pi*v[0])
}
