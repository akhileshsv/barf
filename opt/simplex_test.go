package barf

import (
	"testing"
)

func TestLerpVec(t *testing.T){
	var v1, v2, v3, v4 []float64
	var scale float64
	v1 = []float64{1,1}
	v2 = []float64{2,2}
	scale = 0.5
	v3 = lerpvec(scale, v1, v2)
	scale = -1.5
	v4 = lerpvec(scale, v1, v2)
	if v3[0] != 1.5 && v3[1] != 1.5 && v4[0] != 0.5 && v4[1] != 0.5{
		t.Errorf("vector lerp test failed")
	}
}

func TestNeldMead(t *testing.T){
	var params []float64
	fobj := razaobj
	x0 := [][]float64{{0,0},{1.2,0},{0,0.8}}
	niter := 300
	gv, gbest, err := NeldMead(niter, fobj, x0, params)
	if err != nil{
		t.Log(err)
		t.Errorf("nelder mead opt test failed")
	}
	t.Log("f(x,y) = x2 - 4x + y2 -y -xy\n","global val->", gv, "best fitness->",gbest)

	fobj = parobj
	x0 = [][]float64{{10}}

	gv, gbest, err = NeldMead(niter, fobj, x0, params)
	if err != nil{
		t.Log(err)
		t.Errorf("nelder mead opt test failed")
	}	
	t.Log("y = x2\n","global val->", gv, "best fitness->",gbest)
	
	fobj = rosenpval
	x0 = [][]float64{{-1.2,1.0}}

	gv, gbest, err = NeldMead(niter, fobj, x0, params)
	if err != nil{
		t.Log(err)
		t.Errorf("nelder mead opt test failed")
	}	
	t.Log("rosenbrock's parabolic valley\n","global val->", gv, "best fitness->",gbest)



	fobj = datechobj
	x0 = [][]float64{{-1.0}}
	niter = 0
	gv, gbest, err = NeldMead(niter, fobj, x0, params)
	if err != nil{
		t.Log(err)
		t.Errorf("nelder mead opt test failed")
	}	
	t.Log("y = x2 + 2 * sin(pi * x)\n","global val->", gv, "best fitness->",gbest)
	
	fobj = powelquart
	x0 = [][]float64{{3.0,-1.0,0.0,1.0}}	
	niter = 0
	gv, gbest, err = NeldMead(niter, fobj, x0, params)
	if err != nil{
		t.Log(err)
		t.Errorf("nelder mead opt test failed")
	}	
	t.Log("powell's quartic function\n","global val->", gv, "best fitness->",gbest)
	
}

