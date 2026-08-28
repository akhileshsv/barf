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
	//
	Loc []int
	Typ int
	Pb,Pe *Pt
	Dir int //x = 1, y = 2
	
}

//Rm represents a cell/room 
type Rm struct {
	Id int
	Vertices []Pt
	Cells []*Cell
	Rms []*Rm
	Edges []int
	Touches []int
	Nbrs    []int
	Adj map[int][]*Cell
	Count map[int]int
	Walls map[int][]*Wall
	Area float64
	Centroid *Pt
	Imin, Jmin, Imax, Jmax int
	
}

//Crft represents a craft algo struct
type Crft struct{
	Title   string
	Term    string
	Verbose bool
	Web     bool
	Nrooms  int
	Minvec  [][]int
	Posvec  [][]int
	Dx, Dy  float64
	Mincost float64
	Dimvecx []float64
	Dimvecy []float64
	Fmvec   []float64
	Cmvec   []float64
	Report  string
	Txtplot string
	Rmap    map[int]*Rm
	Grid    Grid
}


//Craft implements the CRAFT algo on a crft struct
func (c *Crft) Craft()(err error){
	swpz := map[int]string{0:"area",1:"adj"}
	//check bounds of cmvec and fmvec
	cmok := len(c.Cmvec) % c.Nrooms == 0 && len(c.Cmvec)/c.Nrooms == c.Nrooms
	fmok := len(c.Fmvec) % c.Nrooms == 0 && len(c.Fmvec)/c.Nrooms == c.Nrooms
	if !cmok || !fmok{
		err = fmt.Errorf("invalid cmvec %v fmvec %v arrays for nrooms %v",cmok, fmok, c.Nrooms)
		return
	}
	cmat := mat.NewDense(c.Nrooms, c.Nrooms, c.Cmvec)
	fmat := mat.NewDense(c.Nrooms, c.Nrooms, c.Fmvec)
	switch{
		case c.Dx + c.Dy == 0.0:
		//find dx, dy
		if len(c.Dimvecx) == 0 || len(c.Dimvecy) == 0{
			err = fmt.Errorf("invalid dimvecx %f dimvecy %f",c.Dimvecx,c.Dimvecy)
			return
		}
		dx := c.Dimvecx[0]
		for _, x := range c.Dimvecx{
			if dx > x{dx = x}
		}
		for _, y := range c.Dimvecy{
			if dx > y{dx = y}
		}
		c.Grid = GridGen(c.Nrooms, c.Posvec, c.Dimvecx, c.Dimvecy, dx, dx)
		default:
		//comes with grid
		c.Grid = Grid{Nx:len(c.Posvec[0]),Ny:len(c.Posvec),Nr:c.Nrooms,Vec:c.Posvec,Dx:c.Dx,Dy:c.Dy}
	}
	c.Rmap = make(map[int]*Rm)
	c.Rmap, _,_,_,_ = LoutGen(c.Grid.Nr,c.Grid.Nx,c.Grid.Ny, c.Grid.Vec, c.Grid.Dx, c.Grid.Dy,[]float64{},[]float64{})
	if c.Verbose{
		outstr := PltLout(c.Rmap)
		fmt.Println(outstr)
	}
	//get combos and room centroids
	combos := CraftCombos(c.Rmap)
	rcent := make([]*Pt,c.Nrooms)
	//room area map
	rcmap := make(map[int]int)
	ramap := make(map[int]float64)
	for i:=1; i <=len(c.Rmap); i++ {
		rcent[i-1] = c.Rmap[i].Centroid
		rcmap[i] = len(c.Rmap[i].Cells)
		ramap[i] = c.Rmap[i].Area
		//fmt.Println(ColorCyan,"room no.->",i,"cell count->",len(c.Rmap[i].Cells),"area->",ramap[i],ColorReset)
	}
	minc := CraftCost(rcent, cmat, fmat)
	mvec := c.Posvec
	//set init cost
	c.Mincost = minc
	c.Minvec = mvec
	var iter, kiter int
	mcosts := []float64{}
	for iter != -1{
		kiter++
		if c.Verbose{fmt.Println("iter no-",kiter)}
		costs, minidx, gridz := CraftCombosEval(rcent, c.Rmap, combos, cmat, fmat, c.Grid.Vec, c.Grid.Dx, c.Grid.Dy)
		for idx, cost := range costs{
			if idx == minidx{
				if c.Verbose{fmt.Println(ColorRed,"dis minimum cost",ColorReset)}
				if c.Verbose{
					fmt.Printf("combo - swap %v and %v cause %v total cost %v\n",combos[idx][0],combos[idx][1],swpz[combos[idx][2]],cost)
					outstr := Plotgrid(gridz[idx],c.Grid.Dx,c.Grid.Dy)
					fmt.Println(outstr)
				}
				mcosts = append(mcosts, cost)
				//c.Rmap = make(map[int]*Rm)
				rmap, _,_,_,_ := LoutGen(c.Grid.Nr,c.Grid.Nx,c.Grid.Ny, gridz[idx], c.Grid.Dx, c.Grid.Dy,[]float64{},[]float64{})		
				//fmt.Println("RMAP LEN",len(c.Rmap))
				for i:=1; i <=len(rmap); i++{
					//fmt.Println("ROOM NO-",i)
					rcent[i-1] = rmap[i].Centroid	
					//fmt.Println(ColorCyan,"room no.->",i,"cell count->",len(rmap[i].Cells),"init->",rcmap[i],"area->",rmap[i].Area,"init->",ramap[i],ColorReset)
					if rcmap[i] != len(c.Rmap[i].Cells){
						fmt.Println(ColorRed,"COUNT ERRORE from transylvania at room",i,ColorReset)
						return
					}
				}
				if minc > cost{
					minc = cost
					mvec = gridz[idx]
				}
			}
			if c.Verbose{fmt.Printf("combo - swap %v and %v cause %v total cost %v\n",combos[idx][0],combos[idx][1],swpz[combos[idx][2]],cost)}
		}
		
		if kiter > c.Nrooms{
			if c.Verbose{fmt.Println("stopping iter (niter > nrooms)")}
			iter = -1
			break
		}
		if len(mcosts) > 2{
			m1 := mcosts[len(mcosts)-1]
			m2 := mcosts[len(mcosts)-2]
			m3 := mcosts[len(mcosts)-3]
			if (m1 > m2 && m1 > m3) || (c.Mincost < m1 && c.Mincost < m2 && c.Mincost < m3){
				if c.Verbose{fmt.Println("stopping iter (no improvement)")}
				iter = -1
				break
			}
		}
	}
	if c.Verbose{
		fmt.Println("min vec-",mvec)
		fmt.Println("min cost-",minc)
	}
	c.Mincost = minc
	c.Minvec = mvec
	return	
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
	gnew := make([][]int, len(grid))
	for i := range rcnew {
		rcnew[i] = &Pt{X:rcent[i].X,Y:rcent[i].Y}
	}
	r1 := combo[0]; r2 := combo[1]
	c1 := rcnew[r1-1]; c2 := rcnew[r2-1]
	if combo[2] == 0 || rmap[r1].Area == rmap[r2].Area{
		rcnew[r1-1] = c2
		rcnew[r2-1] = c1
		for i, row := range grid{
			gnew[i] = make([]int, len(grid[i]))
			copy(gnew[i], grid[i])
			for j, room := range row{
				if room == r1{
					gnew[i][j] = r2
				}
				if room == r2{
					gnew[i][j] = r1
				}
			}
		}
	}
	if combo[2] == 1{
		
		//fmt.Println(ColorYellow, "swap r1 and r2-",r1,r2,ColorReset)
		var rsmol, rlarge, nsmol, nlarge int
		if len(rmap[r1].Cells) < len(rmap[r2].Cells) {
			rsmol = r1
			rlarge = r2
		} else {
			rsmol = r2
			rlarge = r1
		}
		
		//fmt.Println(ColorYellow, "rsmol and rlarge-",rsmol,rlarge,ColorReset)
		//fmt.Println(ColorRed,"init scs-",len(rmap[rsmol].Cells),"lcs-",len(rmap[rlarge].Cells),ColorReset)
		
		nsmol = len(rmap[rsmol].Cells)
		nlarge = len(rmap[rlarge].Cells)
		distmap := make(map[float64][]*Cell)
		rm := rmap[rlarge]
		for _, cell := range rm.Cells{
			dist := math.Pow(cell.Centroid.X - rmap[rsmol].Centroid.X,2) + math.Pow(cell.Centroid.Y - rmap[rsmol].Centroid.Y,2)
			distmap[dist] = append(distmap[dist], cell)
		}
		//GOING BOLDLY AND BLINDLY BY CENTROIDAL CELL DISTANCE
		var dists []float64
		for dist := range distmap {dists = append(dists,dist)}
		sort.Sort(sort.Reverse(sort.Float64Slice(dists)))
		var iter int
		scs := []Tupil{}
		lcs := []Tupil{}
		for _, dist := range dists{
			if iter == 1{
				break
			}
			for _, cell := range distmap[dist]{
				src := Tupil{cell.Row,cell.Col}
				gnew, scs, lcs = celldiv(nsmol,rsmol,rlarge,grid,src,rmap[rsmol].Cells,rmap[rlarge].Cells)
				if len(scs) == nsmol && len(lcs) == nlarge{
					iter = 1
					break
				}
			}
		}
		//fmt.Println(ColorCyan,"scs-",len(scs),"lcs-",len(lcs),ColorReset)
		rcnew[rsmol-1] = CentroidCells(scs,dx,dy)
		rcnew[rlarge-1] = CentroidCells(lcs,dx,dy)
		//fmt.Println(ColorRed, "ingrid\n",grid,ColorReset)
		//fmt.Println(ColorCyan, "outgrid\n",gnew,ColorReset)
	}
	return rcnew, gnew
}

//CraftCombosEval evaluates all room combos/facility layouts for optimal cost
func CraftCombosEval(rcent []*Pt, rmap map[int]*Rm, combos map[string][]int, cmat, fmat *mat.Dense, grid [][]int, dx, dy float64) (map[string]float64, string, map[string][][]int){
	costs := make(map[string]float64, len(combos))
	gridz := make(map[string][][]int, len(combos))
	var cmin float64
	var minidx string
	for idx, combo := range combos {
		//fmt.Println(idx, combo)
		rcnew, gnew := Rswap(combo, rmap, rcent, grid, dx, dy)
		//for _, cent := range rcnew {fmt.Println(cent.X,cent.Y)}
		costs[idx] = CraftCost(rcnew, cmat, fmat)
		gridz[idx] = gnew
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
