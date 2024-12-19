package barf


//Tupil is an int tuple struct
type Tupil struct {
	I int
	J int
}

//Vtx is a vertex defined by a Pt2d and an int tuple location
type Vtx struct {
	Pb Pt2d
	Loc Tupil
}

//Edge is an edge struct from vtxs vb,ve
type Edge struct {
	Vb, Ve Vtx
	Loc Tupil
	Dir string
	Typ int
	Pb, Pe Pt2d
}

//Face is a cell struct
type Face struct {
	Loc Tupil
	Dx, Dy float64
}

//Cx returns the x-coord of the face centroid
func (f * Face) Cx() (float64){
	return float64(f.Loc.J)*f.Dx + f.Dx/2 
}

//Cy returns the y-coord of the face centroid
func (f * Face) Cy() (float64) {
	return float64(f.Loc.I)*f.Dy + f.Dy/2
}

//Pt is a 3d point struct
//lmao strong ambitions these
type Pt struct {
	X,Y,Z float64
	I,J,K int
}

//Lout holds a grid and room map/node/edge data
type Lout struct {
	G *Grid
	Rmap map[int]*Rm
	Nodes map[Pt][]*Wall
}

//Grid is a grid of rooms with dimensions dx, dy
//and rooms from 1 - nr
type Grid struct {
	Nx, Ny, Nr int
	Vec [][]int
	Dx, Dy float64
	Xs, Ys []float64 
	Walls [][]int
	Weights map[Tuple]map[Tuple]float64
	Gstr [][]string
	Vals []float64 //can be anything
	Ni, Nj int //baah calling it this imposs confusion
}

//Cell is a cell struct 
type Cell struct {
	Row int
	Col int
	Pb *Pt
	Pe *Pt
	Dx,Dy,Area float64
	Centroid *Pt
	Room int
}

func IntInVec(vec []int, val int)(bool){
	for _, v := range vec{
		if v == val{
			return true
		}
	}
	return false
}

//CentroidCells returns the centroid of a list of cell indices
func CentroidCells(cs []Tupil, dx, dy float64) (*Pt) {
	var asum, cxsum, cysum float64
	for _, t := range cs {
		carea := dx*dy
		cx := float64(t.J)*dx + dx/2 
		cy := float64(t.I)*dy + dy/2
		asum += carea
		cxsum += carea * cx
		cysum += carea * cy
	}
	return &Pt{
		X: cxsum/asum,
		Y: cysum/asum,
	}	
}

//CentroidCalcs returns the centroid of a slice of cells
func CentroidCalc(cells []*Cell) (*Pt){
	var asum, cxsum, cysum float64
	for _, cell := range cells {
		asum += cell.Area
		cxsum += cell.Area*cell.Centroid.X
		cysum += cell.Area*cell.Centroid.Y
	}
	return &Pt{
		X: cxsum/asum,
		Y: cysum/asum,
	}
}


var (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func nocolor(){
	ColorReset  = ""
	ColorRed    = ""
	ColorGreen  = ""
	ColorYellow = ""
	ColorBlue   = ""
	ColorPurple = ""
	ColorCyan   = ""
	ColorWhite  = ""	
}
