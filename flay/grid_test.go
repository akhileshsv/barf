package barf

import (
	"fmt"
	"testing"
)


func TestBltSecArea(t *testing.T){
	
	t.Log("testing bolt net sectional area calcs")
	var nsas, pss, ggs []float64
	
	t.Log("duggal ex. 1(b)")
	grid := [][]int{
		{1,0,1,0,1,0},
		{0,1,0,1,0,1},
		{1,0,1,0,1,0},
		{0,1,0,1,0,1},
	}
	bw := 300.0
	tp := 8.0
	dia := 20.0
	ps := 65.0
	gg := 75.0
	paths, _, nsamin, mindx := BltSecArea(grid, bw, tp, dia, ps, gg, pss, ggs)

	fmt.Println("NSA MIN->",nsamin,"mm2")
	DrawBltPaths(grid, paths[mindx])

	
	t.Log("bhavikatti ex 5.4")
	grid = [][]int{
		{0,0,1,0,0},
		{0,1,0,1,0},
		{1,0,1,0,1},
		{0,1,0,1,0},
		{0,0,1,0,0},
	}

	
	bw = 160.0
	tp = 8.0
	dia = 18.0
	ps = 40.0
	gg = 25.0
	paths, _, nsamin, mindx = BltSecArea(grid, bw, tp, dia, ps, gg, pss, ggs)
	fmt.Println("NSA MIN->",nsamin,"mm2")
	DrawBltPaths(grid, paths[mindx])

	t.Log("maity ex. 20.1(a)")
	grid = [][]int{
		{1,1,1},
		{1,1,1},
	}
	bw = 200.0
	tp = 8.0
	dia = 18.0
	ps = 75.0
	gg = 50.0
	paths, _, nsamin, mindx = BltSecArea(grid, bw, tp, dia, ps, gg, pss, ggs)

	fmt.Println("NSA MIN->",nsamin,"mm2")
	DrawBltPaths(grid, paths[mindx])


	
	t.Log("maity ex. 20.1(b)")
	grid = [][]int{
		{0,1,1,0},
		{1,0,0,1},
		{0,1,1,0},
	}
	bw = 200.0
	tp = 8.0
	dia = 18.0
	ps = 75.0
	gg = 50.0
	paths, nsas, nsamin, mindx = BltSecArea(grid, bw, tp, dia, ps, gg, pss, ggs)
	for i, path := range paths{
		nsa := nsas[i]
		fmt.Println("nsa->",nsa,"mm2")
		DrawBltPaths(grid, path)
	}
	fmt.Println("NSA MIN->",nsamin,"mm2")
	DrawBltPaths(grid, paths[mindx])

	t.Log("duggal ex. 7.3")
	grid = [][]int{
		{1,0,1,0,1},
		{0,1,0,1,0},
		{1,0,1,0,1},
	}
	bw = 190.0
	tp = 10.0
	dia = 20.0
	ps = 0.0
	gg = 0.0
	pss = []float64{50,50,50,50,50}
	ggs = []float64{35,50,75}
	
	paths, nsas, nsamin, mindx = BltSecArea(grid, bw, tp, dia, ps, gg, pss, ggs)
	for i, path := range paths{
		nsa := nsas[i]
		fmt.Println("nsa->",nsa,"mm2")
		DrawBltPaths(grid, path)
	}
	fmt.Println("NSA MIN->",nsamin,"mm2")
	DrawBltPaths(grid, paths[mindx])

}

/*
//IS NOT WERK
func TestBltGrid(t *testing.T){
	t.Log("testing bolt group failure paths")
	var rezstr string
	bvecs := make([][][]int, 3)
	bvecs[0] = [][]int{
		{1,1,1},
		{1,1,1},
	}
	bvecs[1] = [][]int{
		{0,1,1,0},
		{1,0,0,1},
		{0,1,1,0},
	}
	bvecs[2] = [][]int{
		{1,0,1,0,1},
		{0,1,0,1,0},
		{1,0,1,0,1},
		{0,1,0,1,0},
	}
	for i, bvec := range bvecs{
		ni := len(bvec); nj := len(bvec[0])
		start := Tuple{0,0}
		var g Grid
		g.Init(ni, nj, bvec, [][]int{},false)
		switch i{
			case 0:
			rezstr += "maity lec 20, ex. 1-a, chain bolts"
			g.Vals = []float64{18,75,100}
			case 1:
			rezstr += "maity lec 20, ex. 1-a, staggered bolts"
			g.Vals = []float64{18,75,50}
			start.J = 1
			case 2:
			rezstr += "duggal ex. 7.1-b, staggered bolts"
			g.Vals = []float64{20,65,75}
			start.J = 0
		}
		if i == 2{
			for si := 0; si < 3; si += 2{
				start.J = si
				goal, cfrm, _ := BltGrid(g, start, false)
				path := g.Getpath(start, goal, cfrm)
				fmt.Println(ColorYellow)
				g.Printpath(start, goal,path)
				fmt.Println(ColorReset)
			}
		} else {
			
			goal, cfrm, _ := BltGrid(g, start, false)
			path := g.Getpath(start, goal, cfrm)
			fmt.Println(ColorYellow)
			g.Printpath(start, goal,path)
			fmt.Println(ColorReset)	
		}
	}
	
}

*/
