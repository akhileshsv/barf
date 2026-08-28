package barf

import (
	"fmt"
	"math"
	"gonum.org/v1/gonum/mat"
)

//curves has curve interpolation/generation funcs

//ParabolaK calculates a, b, c for a parabolic profile
//given y1 (lhs), y2 (rhs), yx at xe and xl (net length)
//see pcdc hulse 
func ParabolaK(y1, y2, yx, xe, xl float64)(a,b,c float64, err error){
	d := math.Pow(xe,2) * xl - math.Pow(xl, 2) * xe
	if d == 0{
		err = fmt.Errorf("zero determinant-%f",d)
		return
	}
	a = (xl * (yx-y1) - xe * (y2-y1))/d
	// b = (math.Pow(xl,2)*(y2-y1)+math.Pow(xe,2)*(yx-y1))/d
	b = (-xl*xl*(yx-y1)+xe*xe*(y2-y1))/d
	c = y1
	return
}

//QuadK solves for quad constants given three points on the curve
func QuadK(p1, p2, p3 Pt) (a,b,c float64, err error) {
	//Create the augmented matrix
	x := mat.NewDense(3, 3, []float64{
		p1.X * p1.X, p1.X, 1,
		p2.X * p2.X, p2.X, 1,
		p3.X * p3.X, p3.X, 1,
	})

	y := mat.NewVecDense(3, []float64{p1.Y, p2.Y, p3.Y})
	
	var xinv mat.Dense
	err = xinv.Inverse(x)
	if err != nil {
		err = fmt.Errorf("matrix is not invertible: %v", err)
		return  
	}
	
	var cs mat.VecDense
	cs.MulVec(&xinv, y)
	
	a = cs.At(0, 0)
	b = cs.At(1, 0)
	c = cs.At(2, 0)
    return 
}

//GenCurve generates nstep points on a curve given params/curve type
func GenCurve(ctyp string, nsteps int, xb, xstep float64, params ...float64)(cpts [][]float64, pts []float64, err error){
	switch ctyp{
		case "line":
		//y = mx + c, params - m, c
		case "quad","parabolic","para","quadratic":
		//y = ax2 + bx + c
		if len(params) < 3{
			err = fmt.Errorf("insufficient params to gen quadratic-%f",params)
		}
		a := params[0]; b := params[1]; c := params[2]
		for i := 0; i < nsteps; i++{
			x := xb + float64(i) * xstep
			y := a * math.Pow(x,2) + b * x + c
			cpts = append(cpts, []float64{x, y})
			pts = append(pts, y)
		}
	}
	return
}
