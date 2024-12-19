package barf

import (
	"fmt"
	"testing"
)

func TestIntersect(t *testing.T){
	var a, b, c, d Pt
	a = Pt{0,0}; b = Pt{1,1}; c = Pt{10,10}; d = Pt{30,30}
	e1 := Edge{a,b}; e2 := Edge{c,d}
	fmt.Println(Intersect(e1,e2))
	//EdgePlot([]Edge{e1,e2})
}

func TestSplitEdge2d(t *testing.T){
	v1 := []float64{0,0}
	v2 := []float64{100,0}
	tol := 25.0
	vs := SplitEdge2d(v1, v2, tol)
	fmt.Println(vs)
}
