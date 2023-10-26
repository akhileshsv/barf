package barf

import (
	//"fmt"
	"testing"
)

func TestSecDraw3d(t *testing.T){
	styp := 1
	dims := []float64{50,100}
	s := SecGen(styp, dims)
	wng := []float64{0,0}
	p1 := []float64{0,0,0}
	p2 := []float64{0,10000,0}
	s.Draw3d(p1,p2,wng)
}
