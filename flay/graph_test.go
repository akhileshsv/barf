package barf

import (
	"fmt"
	"testing"
)


func TestNblt(t *testing.T){
	grid := [][]int{
		{1,1,1,0,1,0},
		{1,1,1,1,0,1},
		{1,1,1,1,0,1},
		{1,1,1,1,0,1},
	}
	check := 1
	rezcells, cells, ncomps := Nbcomponents(grid, check)
	fmt.Println(rezcells)
	fmt.Println(cells)
	fmt.Println(ncomps)
}

func TestGraph(t *testing.T){	
	grid := [][]int{
		{1,1,3,4,4},
		{2,2,3,3,3},
		{3,3,4,3,3},
		{2,2,4,5,3},
		{2,2,3,5,3},
		{2,2,3,5,3},
	}
	rms := []int{1,2,3,4,5}
	fmt.Println(grid)
	for _, rm := range rms{
		rezcells, cells, ncomps := Nbcomponents(grid,rm)
		fmt.Println(ColorYellow,"rm->",rm,ColorReset)
		fmt.Println(ColorCyan,"rezcells",len(rezcells),rezcells, ColorReset)
		fmt.Println(ColorRed, "ncells",len(cells), cells, ColorReset)
		fmt.Println("ncomps->",ncomps)
	}
	grid = [][]int{
		{1,1,2},
		{1,3,2},
	}
	var dx, dy float64
	var nx, ny int
	nrms := 3
	dimvecx := []float64{3.0,3.0,3.0}
	dimvecy := []float64{4.0,4.0}
	rmap, _, _, _, _ := LoutGen(nrms,nx,ny,grid,dx,dy,dimvecx, dimvecy)
	outstr := PltLout(rmap)
	fmt.Println(outstr)
}

func TestRSwap1(t *testing.T){
	grid := [][]int{
		{1,1,1,4,4},
		{1,1,2,4,4},
		{5,5,3,3,3},
		{5,5,3,3,3},
	}
	fmt.Println(ColorGreen,grid)
	rsmol := 4; rlarge := 3; nsmol := 0; nlarge := 0
	smolcs := []*Cell{}
	largecs := []*Cell{}
	//celldiv(nsmol,rsmol,rlarge int, grid [][]int, src Tupil, smolcs, largecs []*Cell) ([][]int,[]Tupil,bool)
	for i, row := range grid {
		for j, room := range row {
			if room == rsmol {
				smolcs = append(smolcs, &Cell{Row:i,Col:j})
				nsmol++
			}
			if room == rlarge {
				largecs = append(largecs, &Cell{Row:i,Col:j})
				nlarge++
			}
		}
	}
	src := Tupil{3,2}
	gridedit := make([][]int, len(grid))
	copy(gridedit, grid)
	gridedit, scs, lcs := celldiv(nsmol, rsmol, rlarge, gridedit, src, smolcs, largecs)
	//fmt.Println(ColorCyan,ncomps)
	fmt.Println(ColorRed,gridedit,ColorReset)
	fmt.Println(ColorGreen,scs)
	fmt.Println(ColorCyan,lcs,ColorReset)
}
