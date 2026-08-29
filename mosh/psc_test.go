package barf

import (
	"testing"
)

func TestCBeamPscBm(t *testing.T){
	cb := &CBm{
		Psc:true,
		Fck:25.0,
		Nspans:2,
		Lspans:[]float64{13.0,13.0},
		Clvrs:[][]float64{{0,0,0},{0,0,0}},
		Em:[][]float64{{25e6}},
		Cp:[][]float64{{5e10}},
		Pfrcs:[]float64{5000,5000},
		Pexs:[][]float64{{0.05,-0.14,0.25,8.0},{-0.14,0.0,0.25,6.5}},
		Para:true,
		Term:"dumb",
	}
	cb.PscBm()
	t.Fail()
}

func TestPscStress(t *testing.T){
	t.Log("starting psc fiber stress calc tests")
	pf := 1100e3
	ect := 125.
	bm := 300e6
	ar := 375e3
	zb := 50e6
	zt := 50e6
	ft, fb := PscStress(pf, ect, bm, ar, zt, zb)
	t.Logf("stress at bottom of section %.4f n/mm2 at top %.4f n/mm2",fb,ft)
}
