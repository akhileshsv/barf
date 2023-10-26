package barf

import (
	"testing"
)

func TestBoltSs(t *testing.T){
	b := &Blt{
		Bc:[][]float64{{0,0},{400,0},{0,500},{400,500},{0,800},{400,800}},
		Dia:[]float64{20,20,20,20,20,20},
		Typ:[]float64{1,1,1,1,1,1},
		Frc:[][]float64{{2,-40,0}},
		Fc:[][]float64{{1200,1000}},
	}
	err := BoltSs(b)
	t.Log("\n\n",b.Report)
	if err != nil{
		t.Errorf("bolt group analysis test failed")
	}
}
