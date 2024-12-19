package barf

import (
	"testing"
)

func TestWdFcAng(t *testing.T){
	angs := []float64{0,15,30,45,60,75,90}
	fc := 18.0; fcp := 4.5
	fcs := make([]float64, len(angs))
	for i, ang := range angs{
		fcs[i] = WdFcAng(fc, fcp, ang, true)
	}
	t.Logf("fcs-> %v",fcs)
}
