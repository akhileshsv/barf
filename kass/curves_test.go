package barf

import (
	"fmt"
	"testing"
	draw"barf/draw"
)

func TestParabolaK(t *testing.T){
	y1 := 6.25
	y2 := 6.25 
	yx := 0.0
	xe := 5.0
	xl := 10.0
	a, b, c, _ := ParabolaK(y1, y2, yx, xe, xl)
	t.Log("a,b,c->",a,b,c)
	y1 = -6.25
	y2 = -6.25 
	yx = 0.0
	xe = 5.0
	xl = 10.0
	a, b, c, _ = ParabolaK(y1, y2, yx, xe, xl)
	t.Log("a,b,c->",a,b,c)
}

func TestQuadK(t *testing.T){
	p1 := Pt{0, 6.25}
	p2 := Pt{5, 0}
	p3 := Pt{10, 6.25}
	a, b, c, _ := QuadK(p1, p2, p3)
	t.Log("a,b,c-",a,b,c)

	p1 = Pt{0, 4}
	p2 = Pt{1, 9}
	p3 = Pt{-2, 6}
	a, b, c, _ = QuadK(p1, p2, p3)
	t.Log("a,b,c-",a,b,c)
}

func TestGenCurve(t *testing.T){
	p1 := Pt{0, 6.25}
	p2 := Pt{5, 0}
	p3 := Pt{10, 6.25}
	a, b, c, _ := QuadK(p1, p2, p3)
	nsteps := 21
	xb := 0.
	xstep := 10.0/float64(nsteps-1)
	ctyp := "quad"
	cpts, _, err := GenCurve(ctyp, nsteps, xb, xstep, a, b, c)
	t.Log("errore-",err)
	var data string
	for _, pt := range cpts{
		data += fmt.Sprintf("%f %f\n",pt[0],pt[1])
	}
	skript := "d2.gp"
	term := "dumb"
	txtplot, _ := draw.Dumb(data, skript, term, "para bol saale", "", "", "")
	t.Log(txtplot)
}
