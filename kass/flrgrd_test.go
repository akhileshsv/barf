package barf

import (
	"testing"
)

func TestFlrgrd(t *testing.T){
	fg := &Flrgrd{
		Fgrd:[][]int{
			{3,3,3,3,4},
			{3,3,3,3,4},
			{2,2,2,2,1},
		},
		Dimx:[]float64{4.0,4.0,4.0,4.0,4.5},
		Dimy:[]float64{5.0,5.0,2.5},
		Svec:[][]int{//type, endc, spandir (1-x,2-y), ix, iy, cx, cy
			{-1,1,1,0,0,0,0}, //1 stair typ -1
			{1,1,2,0,0,0,0}, //2 ss slab 
			{1,2,1,1,0,1,0}, //cs slab
			{2,1,0,0,1,0,0},
		},
	}
	fg.Init()
	t.Errorf("floor grid analysis test failed")
}
