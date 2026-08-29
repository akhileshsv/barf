package barf

import (
	"fmt"
	"testing"
)

func TestGenFrmG(t *testing.T){
	coords := [][]float64{
		{0,0},
		{8,0},
		{0,4},
		{4,4},
		{8,4},
		{0,8},
		{4,8},
		{8,8},
	}
	gstr := [][]int{
		{2,3},
		{5},
		{4,6},
		{5,7},
		{8},
		{7},
		{8},
		{},
	}
	f := Frm3d{
		Gcs:coords,
		Gstr:gstr,
		Nflrs:2,
	}
	f.GenFrm()
}

func TestGenFrm3dRcc(t *testing.T){
	fmt.Println("***TASTING***")
	f := Frm3d{
		X:[]float64{0,8,14},
		Y:[]float64{0,4},
		Z:[]float64{0,4.2},
		DL:1.5,
		LL:2.0,
		Plinth:0,
		Sections:[][]float64{{230,460},{230,230}},
		Csec:[]int{1},
		Bsec:[]int{1},
		Bysec:[]int{2},
		Cstyp:1,
		Bstyp:1,
		Bystyp:1,
		Wtwl:10,
		Clvrs:[][]float64{{},{},{},{}},
		Term:"svg",
		Verbose:true,
	}
	err := f.Calc()
	fmt.Println(err)
}

func TestXydx(t *testing.T){
	var xstep, ystep, colstep int
	xstep = 2; ystep = 2; colstep = xstep * ystep
	for _, i := range []int{1,2,3,4,5,6,7,8,9,10,11,12}{
		fmt.Println("i->",i)
		fmt.Println("xdx ydx fdx locdx")
		fmt.Println(getxydx(i,colstep, xstep, ystep))
	}
	
	xstep = 3; ystep = 2; colstep = xstep * ystep
	for _, i := range []int{1,2,3,4,5,6,7,8,9,10,11,12}{
		fmt.Println("i->",i)
		fmt.Println("xdx ydx fdx locdx")
		fmt.Println(getxydx(i,colstep, xstep, ystep))
	}
	
	//var colstep, xstep, ystep int
	xstep = 3; ystep = 2; colstep = xstep * ystep
	for _, i := range []int{1,2,3,4,5,6,7,8,9,10,11,12}{
		fmt.Println("i->",i)
		fmt.Println("xdx ydx fdx locdx")
		fmt.Println(getxydx(i,colstep, xstep, ystep))
	}
	
	xstep = 3; ystep = 3; colstep = xstep * ystep
	for _, i := range []int{1,2,3,4,5,6,7,8,9,10,11,12}{
		fmt.Println("i->",i)
		fmt.Println("xdx ydx fdx locdx")
		fmt.Println(getxydx(i,colstep, xstep, ystep))
	}
	
}
	/*


	f := F3d{
		Fck:[]float64{25.0},
		Fy:[]float64{550.0},
		X:[]float64{0,5,10,15,25},
		Y:[]float64{0,4,10,15},
		Z:[]float64{0,3.5,7.0},
		//Z:[]float64{0,1.5,4.5,7.5},
		Stair:[]int{1,2},
		SlbDL:[]float64{0,1.0,2.0},
		SlbLL:[]float64{0,3.0,4.0},
		WL:[]float64{0,0,0,0},
		Plinth:1,
		Cdim:[][]float64{{460,230}},
		Bdim:[][]float64{{230,460},{230,230}},
		Csec:[]int{1},
		Bsec:[]int{1,1},
		DLw:[]float64{10,10,10,10},
		Slbc:0,
		Styp:1,
		Slfwt:1,
		Flanged:true,
		Clvr:[][]float64{{2.4,1.0,0.75},{},{2.4,1.0,0.75},{}},//l,r,t,b
		Ncols:0,
		Nbms:0,
	}
	term := "svg"
	CalcRcFrm(&f, term, false)

func TestSlab2Dused(t * testing.T){
	ec := []int{0,0,1,0}
	lx := 7010.0
	ly := 3890.0
	fck := 25.0
	fy := 550.0
	fyd := 550.0
	nomcvr := 20.0
	slbc := 1
	wdl := 1.0
	wll := 3.0
	endc, dused, err := slab2dused(ec, wdl, wll, lx, ly, fck, fy, fyd, nomcvr, slbc)
	fmt.Println(endc, dused, err)
	ec = []int{0,0,0,1}
	endc, dused, err = slab2dused(ec, wdl, wll, lx, ly, fck, fy, fyd, nomcvr, slbc)
	fmt.Println(endc, dused, err)
	ec = []int{0,0,1,1}
	endc, dused, err = slab2dused(ec, wdl, wll, lx, ly, fck, fy, fyd, nomcvr, slbc)
	fmt.Println(endc, dused, err)
	ec = []int{0,0,0,0}
	endc, dused, err = slab2dused(ec, wdl, wll, lx, ly, fck, fy, fyd, nomcvr, slbc)
	fmt.Println(endc, dused, err)
}

	*/
