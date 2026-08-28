package barf

import (
	"log"
	"fmt"
	"sort"
	"github.com/fogleman/gg"
)

//pad grid pads a floor grid
func (f *Flr) PadGrid(dx float64)(grid [][]int){
	ncs := int(f.Gx/dx)
	grid = make([][]int, ncs * len(f.Grid))
	for i := range grid{
		grid[i] = make([]int, ncs * len(f.Grid[0]))
	}
	for i, row := range grid{
		for j := range row{
			gx := float64(i) * dx
			gy := float64(j) * dx
			oi := int(gx/f.Gx)
			oj := int(gy/f.Gx)
			grid[i][j] = f.Grid[oi][oj]
		}
	}
	return
	
}

func cartiso(cx, cy, winx, winy int)(scrx, scry int){
	ix := cx - cy
	iy := (cx + cy)/2
	scrx = winx/2.0 + ix
	scry = winy/4.0 + iy
	return
}

func printgrid(grid [][]int){
	log.Println("GRID-")
	for row := 0; row < len(grid); row++ {
		for column := 0; column < len(grid[0]); column++{
			fmt.Print(grid[row][column], " ")
		}
		fmt.Print("\n")
	}
}

//ViewFlat generates a top down view of a floor
func (f *Flr) ViewFlat(){
	
}

//ViewIso generates an isometric view of a floor
func (f *Flr) ViewIso(){
	grid := f.PadGrid(150.0)
	dx := 150.0
	//grid := f.Grid
	//dx := f.Cwidth
	twd := 32
	tht := 32
	twdh := twd/2
	thth := tht/2
	mx := len(grid)
	if len(grid[0]) > mx{
		mx  = len(grid[0])
	}
	winx := mx * twd * 2
	winy := mx * tht * 2
	 
	imf, err := gg.LoadPNG("../data/tiles/floor.png")
	if err != nil {
		panic(err)
	}
	dc := gg.NewContext(winx, winy)
	for i, row := range grid{
		for j := range row{
			cx := j * twdh
			cy := i * thth
			scrx, scry := cartiso(cx, cy, winx, winy) 
			dc.DrawImage(imf, scrx, scry)
		}
	}
	//now do walls
	wx, _ := gg.LoadPNG("../data/tiles/wall.png")
	//depth sort walls
	wsort := make(map[Tupil][]int)
	ws := []Tupil{}
	for idx, vec := range f.Walls{
		//nz := 2
		p1 := f.Pts[vec[0]-1]
		p2 := f.Pts[vec[1]-1]		
		j2 := int(p2[0]/dx)
		i2 := int(p2[1]/dx)
		j1 := int(p1[0]/dx)
		i1 := int(p1[1]/dx)
		wtyp := f.Wvec[idx][0]
		//forget all doors
		if wtyp == 2{
			continue
		}
		imin := i1
		imax := i2
		if i2 < imin{
			imin = i2
			imax = i1
		}
		jmin := j1
		jmax := j2
		if j2 < jmin{
			jmin = j2
			jmax = j1
		}
		switch{
			case p1[0] == p2[0]:
			for i := imin; i <= imax; i++{
				cx := j1 * twdh 
				cy := i * thth
				scrx, scry := cartiso(cx, cy, winx, winy) 
				tscrn := Tupil{scrx,scry}	
				if _, ok := wsort[tscrn]; !ok{
					wsort[tscrn] = []int{wtyp}
					ws = append(ws, tscrn)
				}
				
			}
			case p1[1] == p2[1]:
			
			for j := jmin; j <= jmax; j++{
				
				cx := j * twdh
				cy := i1 * thth
				scrx, scry := cartiso(cx, cy, winx, winy) 
				tscrn := Tupil{scrx, scry}
				if _, ok := wsort[tscrn]; !ok{
					wsort[tscrn] = []int{wtyp}
					ws = append(ws, tscrn)
				}
			}
		}
	}
	//depth sort
	sort.Slice(ws, func(i,j int) bool{
		return ws[i].J < ws[j].J
	})
	//now draw
	for _, w := range ws{
		scrx := w.I
		scry := w.J
		wtyp := wsort[w][0]
		switch wtyp{
			case 0:
			//windo
			for z := 0; z < 6; z++{
				dc.DrawImageAnchored(wx, scrx, scry, 0.5, 0.5)
				scry -= thth/2
			}
			scry = w.J - 19 * thth/2
			dc.DrawImageAnchored(wx, scrx, scry, 0.5, 0.5)
			case 1:
			//int wall
			for z := 0; z < 20; z++{
				dc.DrawImageAnchored(wx, scrx, scry, 0.5, 0.5)
				scry -= thth/2
			}
			case 4:
			//ext wall	
			for z := 0; z < 20; z++{
				dc.DrawImageAnchored(wx, scrx, scry, 0.5, 0.5)
				scry -= thth/2
			}
		}
	}
	dc.SavePNG("../data/out/out.png")
	log.Println("png savedz")
	//first draw grid
	//then walls?
}
