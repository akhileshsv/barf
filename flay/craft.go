package barf

import (
	"fmt"
	"sort"
	"math"
	"gonum.org/v1/gonum/mat"
)

//Wall represents a wall/edge 
type Wall struct {
	//N = 1
	//W = 4
	//Given a tile q,r its four bordering edges are q,r,N ; q,r,W ; q,r+1,N ; q+1,r,W.
	Loc []int
	Typ int
	Pb,Pe *Pt
	Dir int
}

//Rm represents a cell/room 
type Rm struct {
	Id int
	Vertices []Pt
	Cells []*Cell
	Rms []*Rm
	Edges []int
	Touches []int
	Adj map[int][]*Cell
	Count map[int]int
	Walls map[int][]*Wall
	Area float64
	Centroid *Pt
	Imin, Jmin, Imax, Jmax int
}

//LoutGen generates a room map and wall map for a grid of rooms
func LoutGen(nrooms,nx,ny int, grid [][]int, dx,dy float64, dimvecx, dimvecy []float64) (map[int]*Rm, map[Pt][]*Wall){
	//generate room map and wall map 
	var x0, y0, x, y, xmax, ymax float64
	rmap := make(map[int]*Rm)
	nodemap := make(map[Pt][]*Wall)
	var rm *Rm
	prevdir := make([]int,4)
	if dx != 0.0 {
		xmax = dx * float64(nx)
	} else {
		for _, xc := range dimvecx {
			xmax += xc
		}
	}
	if dy != 0.0 {
		ymax = dy * float64(ny)
	} else {
		for _, yc := range dimvecy {
			ymax += yc
		}
	}
	for i, row := range grid {
		prevdir[0] = -1
		x0 = 0.0
		x = 0.0
		y0 = y
		if dy == 0.0 {
			y += dimvecy[i]
		} else {
			y += dy
		}
		for j, room := range row {
			if dx == 0.0 {
				x += dimvecx[j]
			} else {
				x += dx
			}
			p1 := Pt{X:x0,Y:y0,I:i,J:j}
			p2 := Pt{X:x0,Y:y,I:i+1,J:j}
			p3 := Pt{X:x,Y:y,I:i+1,J:j+1}
			p4 := Pt{X:x,Y:y0,I:i,J:j+1}
			width := x - x0
			height := y - y0
			area := width * height
			//left right top bottom : -1, -2, -3, -4
			//i,j,4 i+1
			cell := &Cell{
				Pb:&Pt{X:x0,Y:y0,I:i,J:j},
				Pe:&Pt{X:x,Y:y,I:i+1,J:j+1},
				Dx:width,
				Dy:height,
				Area:area,
				Centroid:&Pt{X:(x+x0)/2.0,Y:(y+y0)/2.0},
				Row:i, Col: j,
				Room:room,
			}
			wleft := Wall{Pb:&p1,Pe:&p2,Loc:[]int{i,j,4}}			
			wright := Wall{Pb:&p3,Pe:&p4,Loc:[]int{i+1,j,4}}
			wtop := Wall{Pb:&p4,Pe:&p1,Loc:[]int{i,j,1}}
			wbottom := Wall{Pb:&p2,Pe:&p3,Loc:[]int{i,j+1,1}}
			if x0 == 0.0 || j == 0{
				prevdir[0] = -1 //left
				wleft.Typ = 0 //external wole
				nodemap[p1] = append(nodemap[p1],&wleft)
				nodemap[p2] = append(nodemap[p2],&wleft)
			} else {
				prevdir[0] = grid[i][j-1]
				if prevdir[0] == room {
					wleft.Typ = -1
				} else {
					nodemap[p1] = append(nodemap[p4],&wleft)
					nodemap[p1] = append(nodemap[p3],&wleft)
				}
			}
			//right edge
			if x == xmax || j == nx-1{
				prevdir[1] = -2
				wright.Typ = 0
				nodemap[p4] = append(nodemap[p4], &wright)
				nodemap[p3] = append(nodemap[p3], &wright)
			} else {
				prevdir[1] = grid[i][j+1]
				if prevdir[1] == room {
					wright.Typ = -1
				} else {
					wright.Typ = 1
					nodemap[p4] = append(nodemap[p4], &wright)
					nodemap[p3] = append(nodemap[p3], &wright)
				}
			}
			//top edge
			if y0 == 0.0 || i == 0{
				prevdir[2] = -3
				wtop.Typ = 0
				nodemap[p1] = append(nodemap[p1], &wtop)
				nodemap[p4] = append(nodemap[p4], &wtop)
			} else {
				prevdir[2] = grid[i-1][j]
				if prevdir[2] == room {
					wtop.Typ = -1
				} else {
					wtop.Typ = 1
					nodemap[p1] = append(nodemap[p1], &wtop)
					nodemap[p4] = append(nodemap[p4], &wtop)
				}
			}
			//bottom edge
			if y == ymax || i == ny-1{
				prevdir[3] = -4
				wbottom.Typ = 0
				nodemap[p1] = append(nodemap[p1], &wbottom)
				nodemap[p2] = append(nodemap[p2], &wbottom)
			} else {
				prevdir[3] = grid[i+1][j]
				if prevdir[3] == room {
					wbottom.Typ = -1
				} else {
					wbottom.Typ = 1
					nodemap[p1] = append(nodemap[p1], &wtop)
					nodemap[p3] = append(nodemap[p3], &wtop)
				}
			}
			if val, ok := rmap[room]; !ok {
				rm = &Rm{
					Id:room,
					Walls:make(map[int][]*Wall),
					Centroid:&Pt{X:0.0,Y:0.0},
					Area:0.0,
					Count:make(map[int]int),
				}
			} else {
					rm = val
			}
			walls := []*Wall{&wleft,&wright,&wtop,&wbottom}
			for idx, dir := range prevdir {
				if dir != room {
					rm.Walls[dir] = append(rm.Walls[dir],walls[idx])
				}	
			}
			touchdir := [][]int{{-1,-1},{1,-1},{-1,1},{1,1}}
			for _, dir := range touchdir {
				rdx := i + dir[0]
				cdx := j + dir[1]
				if (rdx >= 0 && rdx < ny) && (cdx >=0 && cdx < nx) {
					rm.Touches = append(rm.Touches, grid[rdx][cdx])
					if grid[rdx][cdx] != room{rm.Count[grid[rdx][cdx]]++}
				} else {
					rm.Touches = append(rm.Touches, -1)
					rm.Count[-1]++
				}
			}
			rm.Edges = append(rm.Edges, prevdir...)
			rm.Cells = append(rm.Cells,cell)
			rm.Vertices = append(rm.Vertices,[]Pt{p1,p2,p3,p4}...)
			rm.Centroid.X = ((cell.Centroid.X*cell.Area) + (rm.Centroid.X*rm.Area))/(cell.Area+rm.Area)
			rm.Centroid.Y = ((cell.Centroid.Y*cell.Area) + (rm.Centroid.Y*rm.Area))/(cell.Area+rm.Area)
			rm.Area += cell.Area
			x0 = x
			rmap[room] = rm
		}
	}
	return rmap, nodemap
}

//CraftCombos builds a list of room swap combos for the craft algo
//rooms are swapped if adjacent or if they have equal areas
func CraftCombos(rmap map[int]*Rm) map[string][]int {
	combos := make(map[string][]int)
	for i:= 1; i <= len(rmap); i++ {
		for j := i+1; j <= len(rmap); j++ {
			cidx := fmt.Sprintf("%v-%v",i,j)
			if rmap[i].Area == rmap[j].Area {
				combos[cidx] = []int{i, j, 0}
			}
			for _, idx := range rmap[i].Edges {
				if _, ok := combos[cidx]; !ok {
					if idx == j {
						combos[cidx] = []int{i, j, 1}
					}
				}
			}
		}
	}
	return combos
}

//CraftCost computes the overall cost of a facility layout
func CraftCost(rcent []*Pt, cmat, fmat *mat.Dense) (float64) {
	dmat := mat.NewSymDense(len(rcent),nil)
	//fmt.Println(rcent)
	for i :=0; i < len(rcent); i++ {
		for j := i+1; j < len(rcent); j++ {
			dx := math.Abs(rcent[j].X - rcent[i].X)
			dy := math.Abs(rcent[j].Y - rcent[i].Y)
			dmat.SetSym(i,j,dx+dy)
		}
	}
	var tcmat mat.Dense
	tcmat.MulElem(dmat,cmat)
	tcmat.MulElem(fmat, &tcmat)
	return mat.Sum(&tcmat)
}

//AbsDiff returns the absolute difference between two ints
func AbsDiff(a,b int) (int){
	if a < b {
		return b -a
	} else {
		return a - b
	}
}

//IsAdj checks if two cells are adjacent
func IsAdj(c1,c2 *Cell) (bool){
	i := c1.Row; j := c1.Col; i1:= c2.Row; j1 := c2.Col
	adjdir := [][]int{{-1,0},{1,0},{0,-1},{0,1}}
	for _, dir := range adjdir {
		idx := i + dir[0]
		jdx := j + dir[1]
		if idx == i1 && jdx == j1 {
			return true
		}
	}
	return false
}

//Rswap swaps two rooms 
func Rswap(combo []int, rmap map[int]*Rm, rcent []*Pt, grid [][]int, dx, dy float64) ([]*Pt, [][]int) {
	rcnew := make([]*Pt, len(rcent))
	for i := range rcnew {
		rcnew[i] = &Pt{X:rcent[i].X,Y:rcent[i].Y}
	}
	r1 := combo[0]; r2 := combo[1]
	c1 := rcnew[r1-1]; c2 := rcnew[r2-1]
	if combo[2] == 0 || rmap[r1].Area == rmap[r2].Area {
		rcnew[r1-1] = c2
		rcnew[r2-1] = c1
		for i, row := range grid {
			for j, room := range row {
				switch room {
				case r1:
					grid[i][j] = r2
				case r2:
					grid[i][j] = r1
				}
			}
		}
	}
	if combo[2] == 1 {
		//fmt.Println("adj swap between",combo[0],"-->",combo[1])
		var rsmol, rlarge, nsmol int
		if len(rmap[r1].Cells) < len(rmap[r2].Cells) {
			rsmol = r1
			rlarge = r2
		} else {
			rsmol = r2
			rlarge = r1
		}
		nsmol = len(rmap[rsmol].Cells)
		//nlarge := len(rmap[rlarge].Cells)
		//fmt.Println(nsmol,nlarge)
		distmap := make(map[float64][]*Cell)
		rm := rmap[rlarge]
		for _, cell := range rm.Cells {
			dist := math.Abs(cell.Centroid.X - rmap[rsmol].Centroid.X) + math.Abs(cell.Centroid.Y - rmap[rsmol].Centroid.Y)
			distmap[dist] = append(distmap[dist], cell)
		}
		//GOING BOLDLY AND BLINDLY BY CENTROIDAL CELL DISTANCE
		var dists []float64
		for dist := range distmap {dists = append(dists,dist)}
		sort.Sort(sort.Reverse(sort.Float64Slice(dists)))
		var ncomps int
		scs := []Tupil{}
		lcs := []Tupil{}
		for _, dist := range dists {
			if ncomps == 1 {
				break
			}
			for _, cell := range distmap[dist]{
				src := Tupil{cell.Row,cell.Col}
				grid, scs, lcs, ncomps = celldiv(nsmol,rsmol,rlarge,grid,src,rmap[rsmol].Cells,rmap[rlarge].Cells)
				if ncomps == 1 {
					break
				}
			}
		}
		//fmt.Println(len(scs),len(lcs))
		rcnew[rsmol-1] = CentroidCells(scs,dx,dy)
		rcnew[rlarge-1] = CentroidCells(lcs,dx,dy)
	}
	return rcnew, grid
}

//CraftCombosEval evaluates all room combos/facility layouts for optimal cost
func CraftCombosEval(rcent []*Pt, rmap map[int]*Rm, combos map[string][]int, cmat, fmat *mat.Dense, grid [][]int, dx, dy float64) (map[string]float64, string, map[string][][]int){
	costs := make(map[string]float64, len(combos))
	gridz := make(map[string][][]int, len(combos))
	var cmin float64
	var minidx string
	for idx, combo := range combos {
		//fmt.Println(idx, combo)
		gridnew := make([][]int,len(grid))
		for i := range grid {
			gridnew[i] = make([]int,len(grid[i]))
			copy(gridnew[i],grid[i])
		}
		rcnew, gridnew := Rswap(combo, rmap, rcent, gridnew, dx, dy)
		//for _, cent := range rcnew {fmt.Println(cent.X,cent.Y)}
		costs[idx] = CraftCost(rcnew, cmat, fmat)
		gridz[idx] = gridnew
		if cmin == 0.0 {
			cmin = costs[idx]
		} else {
			if cmin > costs[idx] {
				cmin = costs[idx]
				minidx = idx
			}
		}
	}
	return costs, minidx, gridz
}

//GridGen generates a grid of rooms (rep by posvec, dimvecx and dimvecy) with minimum cell width of gridx, gridy
func GridGen(nrooms int, posvec [][]int, dimvecx, dimvecy []float64, gridx, gridy float64) (Grid){
	var xmin, ymin, xn, yn float64
	xs := make([]float64,len(dimvecx))
	ys := make([]float64, len(dimvecy))
	xmin = dimvecx[0]
	ymin = dimvecy[0]
	for idx, x := range dimvecx {
		if x < xmin {xmin = x}
		xn += x
		xs[idx] = xn
	}
	for idx, y := range dimvecy {
		if y < ymin {ymin = y}
		yn += y
		ys[idx] = yn
	}
	if gridx != 0.0 {
		xmin = gridx
	}
	if gridy != 0.0 {
		ymin = gridy
	}
	nx := int(math.Ceil(xn/xmin))
	ny := int(math.Ceil(yn/ymin))
	//row := make([]int, nx)
	grid := make([][]int,ny)
	for i := 0; i < ny; i++ {
		grid[i] = make([]int, nx)
	}
	var x0, y0, x, y float64
	for i, row := range posvec {
		x0 = 0.0
		x = 0.0
		y0 = y
		y = ys[i]
		for j, room := range row {
			x = xs[j]
			jstrt := int(x0/xmin)
			jfin := int(x/xmin)
			istrt := int(y0/ymin)
			ifin := int(y/ymin)
			for i :=istrt; i < ifin; i++ {
				for j := jstrt; j < jfin; j++ {
					grid[i][j] = room
				}
			}
			x0 = x
		}
	}
	gridx = xmin
	gridy = ymin
	return Grid{Nx:nx,Ny:ny,Nr:nrooms,Vec:grid,Dx:gridx,Dy:gridy}
}


/*
func RSwap(combo []int, rmap map[int]*Rm, rcent []*Pt, grid [][]int) ([]*Pt, map[int][][]int){
	rcnew := make([]*Pt, len(rcent))
	for i := range rcnew {
		rcnew[i] = &Pt{X:rcent[i].X,Y:rcent[i].Y}
	}
	r1 := combo[0]; r2 := combo[1]
	c1 := rcent[r1-1]; c2 := rcent[r2-1]
	cellmap := make(map[int][][]int,2)
	if combo[2] == 0 || rmap[r1].Area == rmap[r2].Area {
		rcnew[r1-1] = c2
		rcnew[r2-1] = c1
	}
	if combo[2] == 1 {
		var rsmol, rlarge, nsmol, nlarge int
		if len(rmap[r1].Cells) < len(rmap[r2].Cells) {
			rsmol = r1
			rlarge = r2
		} else {
			rsmol = r2
			rlarge = r1
		}
		nsmol = len(rmap[rsmol].Cells)
		nlarge = len(rmap[rlarge].Cells)
		distmap := make(map[float64][]*Cell)
		rm := rmap[rlarge]
		for _, cell := range rm.Cells {
			dist := math.Abs(cell.Centroid.X - rmap[rsmol].Centroid.X) + math.Abs(cell.Centroid.Y - rmap[rsmol].Centroid.Y)
			distmap[dist] = append(distmap[dist], cell)
		}
		//GOING BOLDLY AND BLINDLY BY CENTROIDAL CELL DISTANCE
		var dists []float64
		for dist := range distmap {dists = append(dists,dist)}
		sort.Sort(sort.Reverse(sort.Float64Slice(dists)))
		var gridedit [][]int
		var ncomps int
		for _, dist := range dists {
			if ncomps == 1 {
				break
			}
			for _, cell := range freemap[dist]{
				gridedit, scs, ncomps := celldiv(nsmol,rsmol,rlarge,gridedit,src,rmap[rsmol].Cells,rmap[rlarge].Cells)
				if ncomps == 1 {
					break
				}
			}
		}
		
		for _, cell := range smolcs{
			cellmap[rsmol] = append(cellmap[rsmol],[]int{cell.Row,cell.Col})
		}
		for _, cell := range largecs{
			//cell.Room = rsmol
			cellmap[rlarge] = append(cellmap[rlarge],[]int{cell.Row,cell.Col})	
		}
		rcnew[rsmol-1] = CentroidCalc(smolcs)
		rcnew[rlarge-1] = CentroidCalc(largecs)
	}
	return rcnew, cellmap
}
*/
