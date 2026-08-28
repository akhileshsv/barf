package barf

import (
	"os"
	"fmt"
	"math"
	"sort"
	"strings"
	"runtime"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"container/heap"
	"gonum.org/v1/gonum/mat"
	polygol"github.com/engelsjk/polygol"
)

//ReadFcflr reads an Fcflr from a file
func ReadFcflr(fname string)(f Fcflr, err error){
	jsonfile, e := ioutil.ReadFile(fname)
	if e != nil{
		err = e
		return
	} 
	err = json.Unmarshal([]byte(jsonfile),&f)
	return
}

//FlrJson saves flr labels and polygons to a .json file
func (f *Flr) FlrJson(){
	polys := PolyVec(f.Polys)
	jmap := make(map[string]interface{})
	jmap["Polys"] = polys
	jmap["Labels"] = f.Labels
	jmap["Width"] = f.Width
	jmap["Height"] = f.Height
	jmap["Grid"] = f.Grid
	jmap["Dx"] = f.Gx
	jmap["Walls"] = f.Walls
	jmap["Pts"] = f.Pts
	jmap["Wvec"] = f.Wvec
	jstr, _ := json.Marshal(jmap)
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	foldr := filepath.Join(basepath,"../data/out")
	fname := f.Title
	if fname == ""{fname = "squarify"}
	fname += ".json"
	fpath := filepath.Join(foldr, fname)
	ioutil.WriteFile(fpath, jstr, os.ModePerm)
	fmt.Println("floor json saved at",fpath)
	return
}


//DrawRectView draws a rect of (local x) length l and (local y) height d from pb, pe 
func DrawRectView(mdx int, d float64, pb, pe []float64) (data string){
	d = d/2.0
	l := Dist3d(pb, pe)
	var p0, p1, p2, p3, p4 []float64
	if len(pb) == 1{
		pb = append(pb, 0)
		pe = append(pe, 0)
	}
	p0 = Rotvec(90.0, pb, pe)
	p4 = Lerpvec(d/l, pb, p0)
	p1 = Lerpvec(-d/l, pb, p0)
	p0 = Rotvec(90.0, pe, pb)
	p3 = Lerpvec(-d/l, pe, p0)
	p2 = Lerpvec(d/l, pe, p0)
	for _, pt := range [][][]float64{{p1,p2},{p2,p3},{p3,p4},{p4,p1}}{
		p0 = pt[0]
		dx := pt[1][0] - pt[0][0]
		dy := pt[1][1] - pt[0][1]
		data += fmt.Sprintf("%f %f %f %f %v\n",p0[0],p0[1],dx, dy, mdx)
	}
	return
}

//Dist3d is a more general 2d/3d distance between two points
//used for drawing labels and stuff lol
func Dist3d(p1, p2 []float64) (dist float64){
	switch len(p1){
		case 1:
		dist = math.Abs(p2[0] - p1[0])
		case 2:
		dist = math.Sqrt(math.Pow(p2[0]-p1[0],2)+math.Pow(p2[1]-p1[1],2))
		case 3:
		dist = math.Sqrt(math.Pow(p2[0]-p1[0],2)+math.Pow(p2[1]-p1[1],2)+math.Pow(p2[2]-p1[2],2))
	}
	return
}

//Rotvec rotates v2 by ang about v1
//NOTE ALL THESE ARE 2D OPS - add switch len(v1) for matrixsees
func Rotvec(ang float64, v1, v2 []float64) (v3 []float64){
	//convert ang to radians
	v3 = make([]float64, len(v1))
	ang = ang * math.Pi/180.0
	//shift origin to v1
	tm1 := mat.NewDense(3,3, []float64{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		-v1[0], -v1[1], 1.0,
	})
	//rotate by ang
	rm := mat.NewDense(3,3, []float64{
		math.Cos(ang), math.Sin(ang), 0.0,
		-math.Sin(ang), math.Cos(ang), 0.0,
		0.0, 0.0, 1.0,
	})
	//shift origin to 0,0
	tm2 := mat.NewDense(3,3, []float64{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		v1[0], v1[1], 1.0,
	})
	//go forth and multiply
	a := mat.NewDense(1,3, []float64{v2[0],v2[1],1.0})
	a.Mul(a, tm1)
	a.Mul(a, rm)
	a.Mul(a, tm2)
	v3[0] = a.At(0,0)
	v3[1] = a.At(0,1)
	return
}

//Lerpvec linearly interpolates between v1 and v2 given a scale (0.0 at v1, 1.0 at v2)
func Lerpvec(scale float64, v1, v2 []float64) (v3 []float64){
	vu, vmod := Unitvec(v1,v2)
	v3 = make([]float64, len(v1))
	switch scale{
		case 0.0:
		copy(v3,v1)
		return
		case 1.0:
		copy(v3,v2)
		return
	}
	for i := range v1{
		v3[i] = v1[i] + scale * vmod * vu[i]
	}
	return

}


//Unitvec finds the unit vector between two points
func Unitvec(v1, v2 []float64) (vu []float64, vmod float64){
	vu = make([]float64, len(v1))
	for i := range v1{
		vu[i] = v2[i] - v1[i]
		vmod += math.Pow(v2[i] - v1[i],2)
	}
	vmod = math.Sqrt(vmod)
	if vmod == 0{return}
	for i := range vu{
		vu[i] = vu[i]/vmod
	}
	return
}

// //EdgeDx. EDGEDX. EdgeDx returns the edge index
func EdgeDx(i, j int)(edx Tupil){
	if i < j{
		edx = Tupil{
			I:i,
			J:j,
		}
	} else {
		edx = Tupil{
			I:j,
			J:i,
		} 
	}
	return
}

//PolyVec converts a slice of polys to a 3d float slice
func PolyVec(polys [][]Pt2d)(fpolys [][][]float64){
	fpolys = make([][][]float64, len(polys))
	for i, poly := range polys{
		for _, pt := range poly{
			vec := []float64{pt.X, pt.Y}
			fpolys[i] = append(fpolys[i],vec)
		}
	}
	return
}

//RmapPoly combines cells in rmap and returns polygons
func RmapPoly(rmap map[int]*Rm)(polys [][]Pt2d, err error){
	polys = make([][]Pt2d, len(rmap))
	for i := 1; i <= len(rmap); i++{
		
		rpoly := [][][][]float64{}
		if _, ok := rmap[i]; !ok{
			err = fmt.Errorf("rmap not in sequence %v",rmap)
			return
		} else {
			for j, cell := range rmap[i].Cells{
				p1, p2, p3, p4 := RectPts(Pt2d{cell.Pb.X,cell.Pb.Y}, Pt2d{cell.Pe.X,cell.Pe.Y})
				vec := Pt2Vec([]Pt2d{{p1.X,p1.Y},{p2.X, p2.Y},{p3.X,p3.Y},{p4.X,p4.Y}})
				
				if j == 0{
					rpoly = append(rpoly, [][][]float64{vec})
				} else {
					cpoly := [][][][]float64{}
					cpoly = append(rpoly, [][][]float64{vec})
					rpoly, err  = polygol.Union(rpoly, cpoly)
					if err != nil{
						return
					}
				}
			}
			for _, val := range rpoly{
				for _, vec := range val{
					for _, pt := range vec{
						polys[i-1] = append(polys[i-1],Pt2d{pt[0],pt[1]})
					} 
				}
			} 
		}
	}
	return
}

//GridNbrs returns the neighbors of a cell (4 dir)
func GridNbrs(grid [][]int, cell Tuple)(nbrs []Tuple){
	imax := len(grid); jmax := len(grid[0])
	allns := []Tuple{
		{cell.I, cell.J-1},
		{cell.I, cell.J+1},
		{cell.I-1, cell.J},
		{cell.I+1, cell.J},
	}
	for _, nbr := range allns{
		if nbr.I >= 0 && nbr.I < imax && nbr.J >=0 && nbr.J < jmax{
			nbrs = append(nbrs, nbr)
		}
	}
	return
}

//GridPathVal returns a path to val from start
func GridPathBasic(dx float64, grid [][]int, start Tuple, check int)(path []Tuple){
	var pq Pque
	cfrm := make(map[Tuple]Tuple)
	csf := make(map[Tuple]float64)
	cfrm[start] = Tuple{-1,-1}
	csf[start] = 0.0
	pq = append(pq, &Item{Tup:start, Pri:0.0})
	heap.Init(&pq)
	iter := 0
	goal := Tuple{-1,-1}
	//start djk/astar loop
	for len(pq) > 0 && iter == 0{
		current := heap.Pop(&pq).(*Item).Tup
		//fmt.Println("at kurrent->",current)
		stopcon := false
		prev := cfrm[current]
		if prev.I != -1{
			//not start, now check cell conn with other rooms
			nbrs := GridNbrs(grid, current)
			for _, nbr := range nbrs{
				rdx := grid[nbr.I][nbr.J]
				if rdx == check{
					fmt.Println("stopcon reached")
					stopcon = true
				}
			}
		}
		if stopcon{
			goal = current
			iter = -1
			break
		}
		nbrs := GridNbrs(grid, current)
		pcur := gridpt(dx, current)
		for _, next := range nbrs{
			pnxt := gridpt(dx, next)
			pdiff := pnxt.Sub(pcur)
			costn := pdiff.Length()
			newcost := csf[current] + costn
			if _, ok := csf[next]; !ok || newcost < csf[next]{
				csf[next] = newcost
				priority := newcost		
				heap.Push(&pq, &Item{Tup:next,Pri:priority})
				cfrm[next] = current
			}
		}
	}
	//build path
	p := []Tuple{}
	current := goal
	if _, ok := cfrm[goal]; !ok{
		fmt.Println("PATH ERRORE IN GOALE")
		return path 
	}
	p = append(p, current)
	for{
		current = cfrm[current]
		stopcon := current.I == start.I && current.J == start.J
		p = append(p, current)
		if stopcon{
			break
		}
	}
	path = make([]Tuple, len(p))
	for i, val := range p{
		path[len(p)-1-i] = val
	}
	return

}

//gridpt returns the centroid of a cell
func gridpt(dx float64, cell Tuple)(p Pt2d){
	xc := dx * float64(cell.J) + dx/2.0
	yc := dx * float64(cell.I) + dx/2.0
	p = Pt2d{xc, yc}
	return
}

//gridpthstic returns cell distance from a slice of uncon room centroids
func gridpthstic(grid [][]int, upts []Pt2d, cell Tuple, dx float64)(dist float64){
	for _, pt := range upts{
		rb := int(math.Round(pt.Y/dx))
		cb := int(math.Round(pt.X/dx))
		dist += math.Abs(float64(cell.I - rb-1)) + math.Abs(float64(cell.J-cb-1))
	}
	return
	
}

func Pt2Vec(pts []Pt2d)(vec [][]float64){
	for _, pt := range pts{
		vec = append(vec, []float64{pt.X, pt.Y})
	}
	return
}

func Vec2Pt(vec [][]float64)(pts []Pt2d){
	for _, pt := range vec{
		pts = append(pts, Pt2d{pt[0],pt[1]})
	}
	return
}

//flatpgol flattens polygol output
func flatpgol(mpoly [][][][]float64)(poly []Pt2d){
	for _, val := range mpoly{
		for _, vec := range val{
			for _, pt := range vec{
				poly = append(poly,Pt2d{pt[0],pt[1]})
			} 
		}
	}
	return
}

//RmSub returns the difference of (poly) r2 with r1
func RmSub(rm1, rm2 []Pt2d)(rm3 []Pt2d){
	r1 := Pt2Vec(rm1)
	r2 := Pt2Vec(rm2)
	p1 := [][][][]float64{[][][]float64{r1}}
	p2 := [][][][]float64{[][][]float64{r2}}
	p3, _ := polygol.Difference(p1,p2)
	rm3 = flatpgol(p3)
	return
}

//WallPoly returns the points of a wall given the centerline coords
func WallPoly(x0, y0, x1, y1, w float64)(pts [][]float64){
	switch{
		case x0 == x1:
		pts = [][]float64{
			{x0-w/2., y0},
			{x0-w/2., y1},
			{x0+w/2., y1},
			{x0+w/2., y0},
			{x0-w/2., y0},
		}
		case y0 == y1:
		pts = [][]float64{
			{x0, y0-w/2.},
			{x0, y0+w/2.},
			{x1, y0+w/2.},
			{x1, y0-w/2.},
			{x0, y0-w/2.},
		}		
	}
	return
}



//DirMatDesi returns a (desi) direction matrix given a list of rooms
func DirMatDesi(facing string,labels []string)(dirmat [][]int){
	nrooms := len(labels)
	dirmat = make([][]int, nrooms)
	for i := range dirmat{
		dirmat[i] = make([]int, 5)
	}
	switch facing{
		case "n":
		case "s":
		case "w":
		case "e":
		for i, val := range labels{
			vec := []int{0,0,0,0,0}
			rlbl := strings.Split(val,"-")[0]
			switch rlbl{
				case "kitchen":
				vec = []int{0,1,1,0,1}
				case "bath":
				vec = []int{1,0,0,1,1}
			}
			copy(dirmat[i],vec)
		}
	}
	return
}




func isum(vec []int)(isum int){
	for _, val := range vec{
		isum += val
	}
	return
}

//KostDir returns a cost vs dirmat adjacency
func (f *Flr) KostDir(grid [][]int)(cost float64, rmap map[int]*Rm, rcent []*Pt){
	dirz := []string{"n","e","s","w","ext"}
	dmap := make(map[int]string)
	for i, dir := range f.Dirs{
		switch i{
			case 0:
			dmap[-1] = dir
			case 1:
			dmap[-2] = dir
			case 2:
			dmap[-3] = dir
			case 3:
			dmap[-4] = dir
		}
	}
	//fmt.Println(ColorYellow,"GRID IN-",grid,ColorReset)
	rcent = make([]*Pt, len(f.Labels))
	rmap, _, _, _,_ = LoutGen(len(f.Labels),len(grid[0]),len(grid), grid, f.Gx,f.Gx, []float64{},[]float64{})
	for i, rm := range rmap{
		//fmt.Println("room->",i,"label-",f.Labels[i-1])
		rcent[i-1] = &Pt{X:rm.Centroid.X,Y:rm.Centroid.Y}
		dvec := f.Dirmat[i-1]
		//fmt.Println("dirvec->",dvec)
		if isum(dvec) == 0{
			//fmt.Println("skipping")
			continue
		}
		eadj := make(map[string]int)
		eadj = map[string]int{
			"n":-1,"s":-1,"e":-1,"w":-1,"ext":-1,
		}
		for _, nbr := range rm.Nbrs{
			if nbr < 0{
				dir := dmap[nbr]
				eadj[dir] = 1
				if eadj["ext"] == -1{eadj["ext"] = 1}
			}
		}
		for j, v := range dvec{
			dir := dirz[j]
			if v > 0 && eadj[dir] < 0{
				//fmt.Println(ColorRed,"rm-",f.Labels[i-1],"isnot konnekt to dir->",dir,ColorReset)
				cost += 1.0
			}
		}
	}
	return
}

//RmCombosEval evals floor room combos based on cost
func (f *Flr) RmCombosEval(opt int, rcent []*Pt, rmap map[int]*Rm, combos map[string][]int, grid [][]int, dx, dy float64) (map[string]float64, string, map[string][][]int){
	costs := make(map[string]float64, len(combos))
	gridz := make(map[string][][]int, len(combos))
	cmin := -1.0
	var minidx string
	for idx, combo := range combos {
		//fmt.Println(idx, combo)
		//for _, cent := range rcnew {fmt.Println(cent.X,cent.Y)}
		switch opt{
			case 1:
			_, gnew := Rswap(combo, rmap, rcent, grid, dx, dy)
			costs[idx],_,_ = f.KostDir(gnew)
			gridz[idx] = gnew
			if cmin == -1.0 {
				cmin = costs[idx]
				minidx = idx
			} else {
				if cmin > costs[idx] {
					cmin = costs[idx]
					minidx = idx
				}
			}	
		}
	}
	return costs, minidx, gridz
}


//CraftDir swaps rooms by area/adj until a min cost is reached
func (f *Flr) CraftDir()(grid [][]int){
	
	if f.Facing == ""{f.Facing = "e"}
	if len(f.Dirmat) == 0{
		f.Dirmat = DirMatDesi(f.Facing, f.Labels)
	}
	f.Dirs = getdirvec(f.Facing)
	fmt.Println(ColorGreen, "main dirs->",f.Dirs, ColorReset)
	
	grid = f.PolyGrid()
	txtplot := Plotgrid(grid, f.Gx, f.Gx)
	fmt.Println(txtplot)

	cost0, rmap, rcent := f.KostDir(grid)
	fmt.Println("room init kost->",cost0)
	var iter, kiter int
	cmin := cost0
	
	for iter != -1{
		kiter++
		fmt.Println("iter #",kiter,"min. cost->",cmin)
		combos := CraftCombos(rmap)
		fmt.Println("COMBOLEN-",len(combos))
		costs, mindx, gridz := f.RmCombosEval(f.Opt, rcent, rmap, combos, grid, f.Gx, f.Gx)
		fmt.Println("mindx, mingrid->",mindx, gridz[mindx])
		_, rmap, rcent = f.KostDir(gridz[mindx])
		if costs[mindx] < cmin{
			fmt.Println("smaller kost seen, copying gmin")
			cmin = costs[mindx]
			for i, val := range gridz[mindx]{
				grid[i] = make([]int, len(val))
				copy(grid[i],gridz[mindx][i])
			}
		}
		if kiter > 2{
			fmt.Println(ColorRed, "stopping iter",ColorReset)
			iter = -1
			break
		}
	}
	fmt.Println("min cost grid->",cmin,"not rupeeses")
	fmt.Println(grid)
	return
}

//HOW. HOW?
//ColGrid generates a column and beam grid for a floor
func (f *Flr) ColGrid(rmap map[int]*Rm, nmap map[Pt][]*Wall, ptmap map[Pt][]int, wmap map[Tupil][]int, pts []Pt)(err error){
	//check for unique x and y values -
	//simplify
	//intersect and plot
	xmap := make(map[float64]bool)
	ymap := make(map[float64]bool)
	xs := []float64{}
	ys := []float64{}
	for pt, walls := range nmap{
		vec := ptmap[pt]
		x := pt.X; y := pt.Y
		
		if len(walls) > 2 && vec[0] > 0{
			if _, ok := xmap[x]; !ok{
				xmap[x] = false
				xs = append(xs, x)
			}
			if _, ok := ymap[y]; !ok{
				ymap[y] = false
				ys = append(ys, y)
			}
		}
	}
	sort.Slice(xs, func(i, j int) bool{
		return xs[i] < xs[j]
	})
	
	sort.Slice(ys, func(i, j int) bool{
		return ys[i] < ys[j]
	})
	f.Cxs = []float64{}
	f.Cys = []float64{}
	tol := 4000.0
	for i, x := range xs{
		switch i{
			case 0:
			f.Cxs = append(f.Cxs, x)
			xmap[x] = true
			case len(xs)-1:
			f.Cxs = append(f.Cxs, x)
			xmap[x] = true
			default:
			xp := xs[i-1]
			if xmap[xp] == false{
				xp = f.Cxs[len(f.Cxs)-1]
				if x - xp > tol{
					f.Cxs = append(f.Cxs, x)
					xmap[x] = true
				}
			} else {
				if x - xp > tol{
					f.Cxs = append(f.Cxs, x)
					xmap[x] = true
				}
			}			
		}
	}

	
	for i, y := range ys{
		switch i{
			case 0:
			f.Cys = append(f.Cys, y)
			ymap[y] = true
			case len(ys)-1:
			f.Cys = append(f.Cys, y)
			ymap[y] = true
			default:
			yp := ys[i-1]
			if ymap[yp] == false{
				yp = f.Cys[len(f.Cys)-1]
				if y - yp > tol{
					f.Cys = append(f.Cys, y)
					ymap[y] = true
				}
			} else {
				if y - yp > tol{
					f.Cys = append(f.Cys, y)
					ymap[y] = true
				}
			}			
		}
	}
	f.Colgrid = [][]float64{}
	for _, x := range f.Cxs{
		for _, y := range f.Cys{
			f.Colgrid = append(f.Colgrid, []float64{x, y})
		}
	}
	
	fmt.Println("f.Cxs->",f.Cxs)
	
	fmt.Println("f.Cys->",f.Cys)
	return
	
}

//ResBmap returns the (residential/house) block 
func (f *Flr) ResBmap()(fb Flr){
	rms := map[string]string{
		"out":"out",
		"kitchen":"service",
		"laundry":"service",
		"pantry":"service",
		"utility":"service",
		"toilet":"private",
		"bath":"private",
		"bed":"private",
		"living":"social",
		"dining":"social",
		"stairs":"social",
		"corridor":"social",
	}
	blocks := []string{"social", "service","private"}
	labels := make([][]string, len(blocks))
	areas := make([][]float64, len(blocks))
	var idx int
	for i, room := range f.Labels{
		area := f.Areas[i]
		bn := strings.Split(room, "_")[0]
		switch rms[bn]{
			case "social":
			idx = 0
			case "service":
			idx = 1
			case "private":
			idx = 2
		}
		labels[idx] = append(labels[idx], room)
		areas[idx] = append(areas[idx], area)
	}
	fb = f.BlockPlan(blocks, labels, areas)
	
	if f.Round{
		for i := range fb.Rooms{
			fb.Rooms[i].SetTol()
		}
	}
	if f.Verbose{
		GPlotFloors(&fb, true)
	}
	fb.SqrRmap()
	fb.CorGen()
	return 
}

//BlockPlan plans a floor as a tree of blocks 
func (f *Flr) BlockPlan(blocks []string, labels [][]string, areas [][]float64) (fb Flr){
	var blockareas []float64
	var sumarea float64
	for i := range blocks {
		for _, area := range areas[i] {
			sumarea += area
		}
		blockareas = append(blockareas, sumarea)
	}
	f1 := Flr{Origin: f.Origin, End: Pt2d{X: f.Width, Y: f.Height}, Name :"bloc"}
	f1.Flrarea()
	blockareas = Scalerooms(&f1, blockareas, false)
	FlrPln(&f1, blockareas, blocks)
	var rmareas []float64
	var roomsfinal []*Flr
	for i, room := range f1.Rooms {
		room.Flrarea()
		rmareas = Scalerooms(room,areas[i], false)
		FlrPln(room, rmareas, labels[i])
		roomsfinal = append(roomsfinal, room)
		
	}
	f1.Rooms = roomsfinal
	fb = Flr{
		Tomm:f.Tomm,
		Width:f.Width,
		Height:f.Height,
		Units:f.Units,
		Origin:f.Origin,
		End:f.End,
		Verbose:f.Verbose,
		Round:f.Round,
		Tol:f.Tol,
		Term:f.Term,
	}
	for _, block := range f1.Rooms{
		for i := range block.Rooms{
			if f.Round{
				block.Rooms[i].Tol = f.Tol
				block.Rooms[i].SetTol()
			}
			room := block.Rooms[i]
			fb.Rooms = append(fb.Rooms, room)
			fb.Labels = append(fb.Labels,room.Name)
			fb.Areas = append(fb.Areas,room.Area)
			
		}
	}
	return fb
}

//Edgedx returns a unique Tupil value for each edge
func Edgedx(jb, je int) (edx Tupil){
	//get a unique value for each edge
	if jb < je{
		edx = Tupil{jb,je}
	} else {
		edx = Tupil{je,jb}
	}
	return 
}

//RectPts returns the vertices of the rectangle defined by origin and end points pb and pe
func RectPts(pb, pe Pt2d)(p1, p2, p3, p4 Pt2d){
	width := pe.X - pb.X
	height := pe.Y - pb.Y
	p1 = Pt2d{pb.X, pb.Y}
	p2 = Pt2d{pb.X+width, pb.Y}
	p3 = Pt2d{pb.X+width, pb.Y+height}
	p4 = Pt2d{pb.X,pb.Y+height}
	return
}


//sigh. RectPtz returns the vertices of the rectangle defined by origin and end points pb and pe
func RectPtz(pb, pe Pt)(p1, p2, p3, p4 Pt){
	width := pe.X - pb.X
	height := pe.Y - pb.Y
	p1 = Pt{X:pb.X, Y:pb.Y}
	p2 = Pt{X:pb.X+width, Y:pb.Y}
	p3 = Pt{X:pb.X+width, Y:pb.Y+height}
	p4 = Pt{X:pb.X,Y:pb.Y+height}
	return
}

//ClassEd classifies an edge as l/r/t/b/interior -1 -2 -3 -4 1
func (f *Flr) ClassEd(p1, p2 Pt2d)(ecls int){
	onleft := (p1.X == p2.X) && p1.X == f.Origin.X
	onright := (p1.X == p2.X) && p1.X == f.End.X
	ontop := (p1.Y == p2.Y) && p1.Y == f.End.Y
	onbot := (p1.Y == p2.Y) && p1.Y == f.Origin.Y
	switch{
		case onleft:
		ecls = -1
		case onright:
		ecls = -2
		case ontop:
		ecls = -3
		case onbot:
		ecls = -4
		default:
		ecls = 1
	}
	return
}




/*
//GetPoly returns a geom.Polygon from a slice of outer and inner ring points
func GetPoly(opts, ipts [][]float64)(poly geom.Polygon){
	p := []geom.Point{}
	for _, val := range opts{
		p = append(p, geom.Point{val[0], val[1]})
	}
	if len(ipts) == 0{	
		poly = geom.Polygon{p}
	} else {
		q := []geom.Point{}
		for _, val := range ipts{
			p = append(p, geom.Point{val[0], val[1]})
		}
		poly = geom.Polygon{p,q}
	}
	return

}

//RmSub returns the difference of (poly) r2 with r1
func RmSub(rm1, rm2 []Pt2d)(rm3 []Pt2d){
	r1 := Pt2Vec(rm1)
	r2 := Pt2Vec(rm2)
	p1 := GetPoly(r1, [][]float64{})
	p2 := GetPoly(r2, [][]float64{})
	p3 := p1.Difference(p2)
	for _, poly := range p3.Polygons(){
		for _, ring := range poly{
			for _, pt := range ring{
				rm3 = append(rm3, Pt2d{pt.X, pt.Y})
			}
		}		
	}
	rm3 = rm3[:len(rm3)-1]
	return
}


//RmSwp swaps two rooms 
func RmSwp(rm1, rm2 []Pt2d)(rm3, rm4 []Pt2d){
	return
}

*/
