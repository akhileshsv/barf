package barf

import (
	"testing"
)



func TestWeldSs(t *testing.T){
	w := &Wld{
		Size:[]float64{10,10,10},
		Wc:[][]float64{{0,1200,800,1200},{0,1200,0,0},{0,0,800,0}},
		Frc:[][]float64{{2,-40,0}},
		Fc:[][]float64{{1400,1200}},
	}
	err := WeldSs(w)
	t.Log("\n\n",w.Report)
	if err != nil{
		t.Errorf("weld group analysis test failed")
	}
}
