package barf

import (
	"fmt"
	"math"
	kass"barf/kass"
)

//psc has prestressed concrete funcs
//ofc from pcdc, hulse

//PscStress calculates extreme fibre stresses in a psc section
//hulse 2.2
func PscStress(pf, ect, bm, ar, zt, zb float64)(ft, fb float64){
	ft = pf/ar - pf*ect/zt + bm/zt
	fb = pf/ar + pf*ect/zb - bm/zb
	return
}


//PscBmAr calculates the area under the free bending moment diagram 
func PscBmAr(ndiv int, lspan, pfrc float64, pex []float64)(ar, xc, ml, mr float64){
	var siga, fa, fb float64
	hs := kass.GenSimps(ndiv)
	for i, homerx := range hs{
		ar += homerx * pex[i] * lspan/20.0 * 1.0/3.0
		fb = pex[i] * lspan/20.0
		if i > 0 {
			fa = pex[i-1] * lspan/20.0
			siga += lspan * (fa + fb + 2.0 * (fa + fb))/6.0/20.0
			xca := float64(i-1) * lspan/20.0 + lspan/40.0
			xc += siga * xca
		}
	}
	xc = xc/ar
	ar = ar * pfrc
	ml = 6.0 * ar * (2.0 * lspan/3.0 - xc)/math.Pow(lspan,2) 
	mr = 6.0 * ar * (xc - lspan/3.0)/math.Pow(lspan,2) 
	return
}

//PscBmExts generates tendon eccentricities at 21th spans
func PscBmExts(ctyp string, lspan, pfrc float64, pex []float64)(exts []float64, err error){
	switch ctyp{
		case "parabolic","quadratic","para","quad":
		if len(pex) < 4{
			err = fmt.Errorf("invalid params for parabolic profile - %f",pex)
			return
		}
		var p1, p2, p3 kass.Pt
		p1 = kass.Pt{X:0, Y:pex[0]}
		
		p2 = kass.Pt{X:pex[3],Y:pex[2]}
		p3 = kass.Pt{X:lspan, Y:pex[1]}
		a, b, c, e := kass.QuadK(p1, p2, p3)
		if e != nil{
			err = fmt.Errorf("error generating parabolic curve constants - %s",e)
			return
		}
		xb := 0.0 //CHECK THIS
		_, exts, err = kass.GenCurve(ctyp, 21, xb, lspan/20.0, a, b, c)
		if err != nil{fmt.Println(err)}
		fmt.Println("EXTS-",exts)
		return
	}
	return
}
