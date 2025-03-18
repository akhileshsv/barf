package barf

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"container/heap"
)



//Room is similar to a cell, which is in a room
//why is this not a Flr? unused anyway
type Room struct {
	Id int
	Label string
	Vtxs []Vtx
	Faces []Face
	Edges []*Edge
	Dirs []int
	Area float64
	Centroid Pt2d
}

//Block stores a floor and dimensions
//and has many other unused fields that were added on a whim
type Block struct {
	Floor Flr
	Shape string
	Idmap map[string]int
	Dims  []float64
	Nodes map[Pt2d]int
	Nval  map[int]Pt2d
	Edges map[Tupil][]int//wall jb, je, typ, dir
	Eval  map[int]Tupil
	Nmap  map[int][][]int//n1:[nadj1,wall1],[nadj2,wall2]
	Rmap  map[int][]int
	Rooms []string
}

//Flr is a floor struct, which could be a room?
type Flr struct{
	Name       string
	Title      string
	Area       float64
	Origin     Pt2d
	End        Pt2d
	Mid        Pt2d //centroid/mid point
	Bbo        Pt2d //minx, miny
	Bbe        Pt2d //maxx, maxy
	Rooms      []*Flr
	Cell_count int
	Cell_block []Pt2d
	Space      Pt2d
	Width      float64
	Height     float64
	Cwidth     float64 //corridor width
	Gx         float64 //(square) grid cell size
	Adj        []string//left right top bottom
	Dirs       []string//(ext) dirs - l r t b
	Edges      []Tupil
	Nodes      []Pt2d
	Polys      [][]Pt2d
	Wpolys     [][]Pt2d
	Nmap       map[Pt2d][][]int
	Emap       map[Tupil][][]int
	Lcmap      map[int]bool //living room connectivity map
	Rmap       map[int][]int //room connectivity map
	Walls      [][]int 
	Iwalls     []Tupil
	Wvec       [][]int
	Pts        [][]float64
	Nbrs       []int
	Grid       [][]int
	Colgrid    [][]float64
	Cxs, Cys   []float64
	Dirmat     [][]int //direction adjacency matrix
	Coords     [][]float64
	Togrid     bool //use grid rep of floor
	Cgrid      bool //generate corridor grid
	Sort       bool //if true, sort rooms and labels
	Bloc       bool //if true, switch rooms to blocks
	Isroot     bool
	Round      bool //if true, round all vals to tol
	Tomm       bool //convert all to mm
	Verbose    bool
	Sqrd       bool //squarified map/corridor generated
	Elib       bool //use ctessum/geom poly intersection (ext lib)
	Areas      []float64
	Labels     []string
	Tol        int //if > 0, round all vals to this
	Opt        int //opt. rules
	Ldx        int //living room index
	Units      string
	Term       string
	Txtplot    string
	Facing     string
	Ftyp       string //"res"/etc
	Tmp        []interface{}
}

//getdirvec returns left, right, top, bottom dirs  
func getdirvec(face string)(dirs []string){
	if face == ""{
		face = "e"
	}
	basevec := []string{"n","e","s","w"}
	strt := -1
	for i, dir := range basevec{
		if dir == face{
			strt = i
			break
		}
	}
	bottom := basevec[strt]
	top := strt + 2
	if top > 3{
		top = strt - 2
	}
	left := strt + 1
	if left > 3{
		left = strt - 3
	}
	right := left + 2
	if right > 3{
		right = left - 2
	}
	//this is LEFT RIGHT TOP BOTTOM 
	dirs = []string{
		basevec[left], basevec[right], basevec[top],bottom,
	}
	return
}

//Init inits a floor
func (f *Flr) Init(){
	switch f.Units{
		case "ft":
		f.Width = f.Width * 304.8
		f.Height = f.Height * 304.8
		f.Cwidth = f.Cwidth * 304.8
		f.Area = f.Width * f.Height
		aconv := make([]float64, len(f.Areas))
		copy(aconv, f.Areas)
		for i, ar := range aconv{
			f.Areas[i] = ar * 304.8 * 304.8
		}
		f.Units = "mm"
	}
	if f.Origin.X == f.End.X && f.Origin.Y == f.End.Y{	
		f.Origin = Pt2d{0,0}
		f.End = Pt2d{f.Width, f.Height}
		
	}
	if f.Sort{
		areas := Scalerooms(f, f.Areas, f.Round)
		asrt, lsrt := sortrooms(f.Units, areas, f.Labels)
		f.Areas = make([]float64, len(asrt))
		f.Labels = make([]string, len(lsrt))
		copy(f.Areas, asrt)
		copy(f.Labels, lsrt)
		
	}
	if f.Cwidth == 0.0{f.Cwidth = 900.0}
	if f.Gx == 0.0{f.Gx = f.Cwidth}
	f.Round = true
	//default floor type is residential
	if f.Ftyp == ""{f.Ftyp = "res"}
}

//sort rooms sorts labels and areas
func sortrooms(units string, areas []float64, labels []string)(asrt []float64, lsrt []string){
	amap := make(map[string]float64)
	for i, lbl := range labels {
		amap[lbl] = areas[i]
	}
	sort.Slice(labels, func(i, j int) bool {
		return amap[labels[i]] > amap[labels[j]]
	})
	for _, lbl := range labels{
		area := amap[lbl]
		asrt = append(asrt, area)
		lsrt = append(lsrt, lbl)
	}
	return
}

//Scalerooms scales a slice of room areas to match a floor's total area
func Scalerooms(f *Flr, r []float64, round bool) []float64{
	//first sort room slice
	//sort.Sort(sort.Reverse(sort.Float64Slice(r)))
	if f.Area == 0{
		f.Area = f.Width * f.Height
	}
	var tot_area float64
	for _, rm_area := range r{
		tot_area += rm_area
	}
	scale := f.Area / tot_area
	r_scale := []float64{}
	for _, rm_area := range r{
		if round {
			r_scale = append(r_scale, math.Round((rm_area)*scale))
		} else {
			r_scale = append(r_scale, (rm_area)*scale)
		}
	}
	return r_scale
}

//Clone returns a copy of a floor with basic fields filled in
func (f *Flr) Clone()(fn Flr){
	
	fn = Flr{
		Tomm:f.Tomm,
		Width:f.Width,
		Height:f.Height,
		Units:f.Units,
		Origin:f.Origin,
		End:f.End,
		Areas:f.Areas,
		Labels:f.Labels,
		Verbose:f.Verbose,
		Round:f.Round,
		Bloc:f.Bloc,
		Term:f.Term,
		Polys:f.Polys,
		Sqrd:f.Sqrd,
	}
	return
}

//Flrarea calcs the floor area
func (f *Flr) Flrarea(){
	f.Width = f.End.X - f.Origin.X
	f.Height = f.End.Y - f.Origin.Y
	f.Area = (f.End.X - f.Origin.X) * (f.End.Y - f.Origin.Y)
}

//Addroom adds a new room to a floor
//uses the squarify algo
//https://www.huy.dev/squarified-tree-map-reasonml-part-1-2019-03/
func Addroom(f *Flr, rm_area float64, label string) {
	var dx_new, dy_new, dx_add, dy_add, asr_new, asr_add float64
	//zero rooms
	if len(f.Rooms) == 0 {
		if f.Height > f.Width {
			dy_new = f.Origin.Y + rm_area/f.Width
			dx_new = f.End.X	
			r := &Flr{Origin: f.Origin, End: Pt2d{X: dx_new, Y: dy_new}, Area: rm_area, Name:label}
			f.Space = Pt2d{X: r.End.X, Y: f.Origin.Y}
			f.Cell_block = []Pt2d{r.Origin, r.End}
			f.Cell_count++	
			f.Rooms = append(f.Rooms, r)
			return
		}
		dx_new = f.Origin.X + rm_area/f.Height
		dy_new = f.End.Y
		r := &Flr{Origin: f.Origin, End: Pt2d{X: dx_new, Y: dy_new}, Area: rm_area, Name:label}
		f.Space = Pt2d{X: r.End.X, Y: f.Origin.Y}
		f.Cell_block = []Pt2d{r.Origin, r.End}
		f.Cell_count++
		f.Rooms = append(f.Rooms, r)
		return
	}
	//asr for new block
	y_rem := f.End.Y - f.Space.Y
	x_rem := f.End.X - f.Space.X
	switch y_rem <= x_rem{
	case true:
		dx_new = rm_area / y_rem
		dy_new = y_rem
		asr_new = math.Max(dy_new, dx_new) / math.Min(dy_new, dx_new)
	case false:
		dy_new = rm_area / x_rem
		dx_new = x_rem
		asr_new = math.Max(dy_new, dx_new) / math.Min(dy_new, dx_new)
	}
	//asr for add to current cell block
	cell_area := rm_area
	for _, room := range f.Rooms[len(f.Rooms)-f.Cell_count:] {
		cell_area += room.Area
	} //area of current cell block

	//either the cell block ends at f.end.y or it extends along x (or it gets the hose again)
	if f.End.Y == f.Cell_block[1].Y {
		//increase cell block x, scaling areas by cell block height(y)
		dy_add = (rm_area / cell_area) * (f.Cell_block[1].Y - f.Cell_block[0].Y)
		dx_add = rm_area / dy_add
		asr_add = math.Max(dy_add, dx_add) / math.Min(dy_add, dx_add)

	} else { //increase cell block y, scaling areas by cell block width(x)
		dx_add = (rm_area / cell_area) * (f.Cell_block[1].X - f.Cell_block[0].X)
		dy_add = rm_area / dx_add
		asr_add = math.Max(dy_add, dx_add) / math.Min(dy_add, dx_add)
	}
	//pick the smaller asr
	if asr_add < asr_new {
		//add to current block at origin (of) cell_block[0]
		r := &Flr{Area: rm_area, Origin: f.Cell_block[0],Name:label}
		//new row starts at origin, ends at dx_add, dy_add
		r.End.X = f.Cell_block[0].X + dx_add
		r.End.Y = f.Cell_block[0].Y + dy_add
		switch f.End.Y == f.Cell_block[1].Y {
		case true: //scale y by area ratio, origin x of all boxes stays constant
			for i, room := range f.Rooms[len(f.Rooms)-f.Cell_count:] {
				idx := (len(f.Rooms) - f.Cell_count) + i
				f.Rooms[idx].Origin.X = r.Origin.X
				f.Rooms[idx].End.X = r.End.X
				dy_room := (room.Area / cell_area) * (f.Cell_block[1].Y - f.Cell_block[0].Y) //area/total * cell y
				//dx_room := room.area/dy_room //should be equal to r.end.x
				f.Rooms[idx].Origin.Y = f.Cell_block[0].Y + dy_add //room(n).origin.y = dy0+dy1...dyn-1
				f.Rooms[idx].End.Y = f.Cell_block[0].Y + dy_add + dy_room
				dy_add += dy_room
			}
			f.Cell_block = []Pt2d{r.Origin, {X: r.End.X, Y: dy_add}} //dy_add should be = f.end.y
			f.Space = Pt2d{X: r.End.X, Y: r.Origin.Y}
		case false: //change x coordinates by area ratio, origin y stays the same
			for i, room := range f.Rooms[len(f.Rooms)-f.Cell_count:] {
				idx := (len(f.Rooms) - f.Cell_count) + i
				f.Rooms[idx].Origin.Y = r.Origin.Y
				f.Rooms[idx].End.Y = r.End.Y
				dx_room := (room.Area / cell_area) * (f.Cell_block[1].X - f.Cell_block[0].X)
				//dx_room := room.area/dy_room
				f.Rooms[idx].Origin.X = f.Cell_block[0].X + dx_add
				f.Rooms[idx].End.X = f.Cell_block[0].X + dx_room + dx_add
				dx_add += dx_room
			}
			f.Cell_block = []Pt2d{r.Origin, {X: dx_add, Y: r.End.Y}}
			f.Space = Pt2d{X: r.Origin.X, Y: r.End.Y}
		}
		r.Name = label
		f.Rooms = append(f.Rooms, r)
		f.Cell_count++
		
	} else {
		//add new block at origin f.space
		r := &Flr{Area: rm_area, Origin: f.Space, Name: label}
		r.End.X = f.Space.X + dx_new
		r.End.Y = f.Space.Y + dy_new
		f.Cell_count = 1
		f.Cell_block = []Pt2d{r.Origin, r.End}
		if r.End.X == f.End.X {
			f.Space = Pt2d{X: r.Origin.X, Y: r.End.Y}
		} else {
			f.Space = Pt2d{X: r.End.X, Y: r.Origin.Y}
		}
		f.Rooms = append(f.Rooms, r)
	}	
}


//Flrpln calls Addroom until it runs outta rooms
//tis the squarified algo - https://www.huy.dev/squarified-tree-map-reasonml-part-1-2019-03/
func FlrPln(f *Flr, r []float64, labels []string){
	if len(r) == 0 {
		return
	}
	Addroom(f, r[0], labels[0])
	FlrPln(f, r[1:], labels[1:])
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

//SqrRmap generates a room connectivity map from squarify output
func (f *Flr) SqrRmap(){
	// dirs := map[int]string{
	// 	-1:"left",-2:"right",-3:"top",-4:"bottom",1:"interior",
	// }	
	f.Nmap = make(map[Pt2d][][]int)
	f.Emap = make(map[Tupil][][]int)
	for i, room := range f.Rooms{
		//if f.Verbose{fmt.Println("room->",i+1, "label-",f.Labels[i],"points-",room.Origin, room.End)}
		
		p1, p2, p3, p4 := RectPts(room.Origin, room.End)
		for _, p := range []Pt2d{p1, p2, p3, p4}{
			if _, ok := f.Nmap[p]; !ok{
				idx := len(f.Nodes) + 1
				f.Nodes = append(f.Nodes, p)
				f.Nmap[p] = make([][]int,3)
				f.Nmap[p][0] = []int{idx}
				f.Nmap[p][1] = []int{i}
				f.Nmap[p][2] = []int{1}
				if (p.X == f.Origin.X) || (p.X == f.End.X) || (p.Y == f.Origin.Y) || (p.Y == f.End.Y){
					
					f.Nmap[p][2] = []int{2}
				}
			} else {
				if !IntInVec(f.Nmap[p][1],i){
					f.Nmap[p][1] = append(f.Nmap[p][1],i)
				}
			}
		}
		//edges are - p1p2, p2p3, p3p4, p4p1
		edges := [][]Pt2d{
			{p1,p2},
			{p2,p3},
			{p3,p4},
			{p4,p1},
		}
		
		for _, e := range edges{
			//each edge in map is - jb, je, ecls 
			//exterior - left right top bottom - -1, -2, -3, -4
			ecls := f.ClassEd(e[0],e[1])
			jb := f.Nmap[e[0]][0][0]; je := f.Nmap[e[1]][0][0]
			edx := Edgedx(jb, je)
			if _, ok := f.Emap[edx]; !ok{
				f.Emap[edx] = make([][]int,4)
				f.Emap[edx][0] = []int{jb, je, ecls}
				f.Emap[edx][1] = []int{i}
				f.Emap[edx][2] = []int{1}
				f.Edges = append(f.Edges, edx)
				f.Emap[edx][3] = []int{len(f.Edges)}
				if ecls == 1{
					f.Iwalls = append(f.Iwalls, edx)
				}
				if f.Labels[i] == "living"{
					f.Emap[edx][2] = []int{0}
				}
			} else {
				if !IntInVec(f.Emap[edx][1],i){
					f.Emap[edx][1] = append(f.Emap[edx][1],i)
				}
			}
			f.Rooms[i].Edges = append(f.Rooms[i].Edges, edx)
		}
	}
	f.Rmap = make(map[int][]int)
	for i, room := range f.Rooms{
		f.Rmap[i] = []int{}
		for _, edx := range room.Edges{
			edge := f.Emap[edx]
			// jb := edge[0][0]; je := edge[0][1]; ecls := edge[0][2]
			for _, nbr := range edge[1]{
				if nbr != i && !IntInVec(f.Rmap[i],nbr){
					f.Rmap[i] = append(f.Rmap[i],nbr)
				}
			}
		}
	}
	
	// for i, r1 := range f.Rooms{
	// 	// fmt.Println("checking room-",f.Labels[i])
	// 	nbrs := f.Rmap[i]
	// 	for _, edx := range r1.Edges{
	// 		jb := f.Emap[edx][0][0]
	// 		je := f.Emap[edx][0][1]
	// 		a := f.Nodes[jb-1]
	// 		b := f.Nodes[je-1]
	// 		for j, r2 := range f.Rooms{
	// 			if j != i{	
	// 				//fmt.Println("checking against-",f.Labels[j])
	// 				if !IntInVec(nbrs, j){	
	// 					for _, edy := range r2.Edges{
	// 						j1 := f.Emap[edy][0][0]
	// 						j2 := f.Emap[edy][0][1]
	// 						c := f.Nodes[j1-1]
	// 						d := f.Nodes[j2-1]
	// 						if EdgeOverlap(a,b,c,d){
	// 							f.Rmap[i] = append(f.Rmap[i], j)
	// 							if !IntInVec(f.Emap[edy][1],i){
	// 								f.Emap[edy][1] = append(f.Emap[edy][1],i)
	// 							}
	// 							if !IntInVec(f.Emap[edx][1],j){
	// 								f.Emap[edx][1] = append(f.Emap[edx][1],j)
	// 							}								
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	//new trimmed down loope
	
	rintz := make(map[Tupil]bool)
	for i, r1 := range f.Rooms[:len(f.Rooms)-1]{
		for j, r2 := range f.Rooms[i+1:]{
			j += i + 1
			
			rdx := Edgedx(i, j)
			if _, ok := rintz[rdx]; !ok{
				//process edges of both rooms
				for _, edx := range r1.Edges{
					jb := f.Emap[edx][0][0]
					je := f.Emap[edx][0][1]
					a := f.Nodes[jb-1]
					b := f.Nodes[je-1]
					for _, edy := range r2.Edges{
						j1 := f.Emap[edy][0][0]
						j2 := f.Emap[edy][0][1]
						c := f.Nodes[j1-1]
						d := f.Nodes[j2-1]
						
						if EdgeOverlap(a,b,c,d){
							f.Rmap[i] = append(f.Rmap[i], j)
							f.Rmap[j] = append(f.Rmap[j], i)
							if !IntInVec(f.Emap[edy][1],i){
								f.Emap[edy][1] = append(f.Emap[edy][1],i)
							}
							if !IntInVec(f.Emap[edx][1],j){
								f.Emap[edx][1] = append(f.Emap[edx][1],j)
							}								
						}
					} 
				}
				rintz[rdx] = true
			}
		}
	}
	//remove duplicates (baah)
	for i, nbrs := range f.Rmap{
		ns := []int{}
		for _, nbr := range nbrs{
			if !IntInVec(ns, nbr){
				ns = append(ns, nbr)
			}
		}
		f.Rmap[i] = ns
	}
	
	ldx := 0
	f.Lcmap = make(map[int]bool)
	//get a list of unkonnekted rooms
	for i, room := range f.Labels{
		if room == "living"{
			ldx = i
			f.Ldx = ldx
		}
		f.Lcmap[i] = false
	}
	
	for _, nbr := range f.Rmap[ldx]{
		f.Lcmap[nbr] = true 
	}
	f.Lcmap[ldx] = true
	
}

func nodehstic(p1 Pt2d, upts []Pt2d)(dist float64){
	for _, p2 := range upts{
			pd := p2.Sub(p1)
			dist += pd.Length()
	}
	return
}

//CorInt returns the points of intersection of each room with the corridor polygon
func (f *Flr) CorInt(cpts []Pt2d, e2s [][]Pt2d)(pts []Pt2d){
	//first build corridor edges
	eds := [][]Pt2d{}
	for i, p1 := range cpts{
		if i == len(cpts)-1{
			p2 := cpts[0]
			eds = append(eds, []Pt2d{p1, p2})
		} else {
			p2 := cpts[i+1]
			eds = append(eds, []Pt2d{p1, p2})
		}
	}
	pmap := make(map[Pt2d]bool)
	p1, p2, p3, p4 := RectPts(f.Origin, f.End)
	var intersects bool
	//get points of intersection of all edges
	//retain all points on edge of either polygon
	//sort cw
	for _, val := range [][]Pt2d{
		{p1, p2},
		{p2, p3},
		{p3, p4},
		{p4, p1},
	}{
		a := val[0]; b := val[1]
		for _, ced := range eds{
			c := ced[0]; d := ced[1]
			if c.InRect(f.Origin, f.End) || d.InRect(f.Origin,f.End){
				cls, px := EdgeInt(a, b, c, d)
				//fmt.Println("here-",f.Name, cls, px)
				if cls == "cross"{
					intersects = true
					for i, node := range []Pt2d{a, b, c, d, px}{
						if _, ok := pmap[node]; !ok{
							
							switch i{
								case 0, 1:
								//add these points if they are not on an edge of the corridor
								if !node.OnEdge(cpts){
									pmap[node] = true
								} else {
									pmap[node] = false
								}
								case 2, 3:
								//add em if in rectangle
								if node.InRect(f.Origin, f.End){
									pmap[node] = true
								} else {
									pmap[node] = false
								}
								case 4:
								//add px if in rectangle and if on cor edge
								if node.InRect(f.Origin, f.End) && node.OnEdge(cpts){
									pmap[node] = true
								} else {
									pmap[node] = false
								}
							}
						}	
					}
				}
			}
		}
	}
	
	if !intersects{
		fmt.Println(ColorYellow, f.Name, "does not intersekto",ColorReset)
		pts = []Pt2d{p1, p2, p3, p4}
		return
	} else {
		fmt.Println(ColorRed, f.Name, "intersekto korpus korridore",ColorReset)
		//now build pts
		for pt, val := range pmap{
			if val{
				right := false
				//remove pts on right of all new edges
				for _, ed := range e2s{
					lcs := pt.Lclass(ed[0],ed[1])
					if lcs == "right"{
						right = true
					}
				}
				if !right{
					pts = append(pts, pt)
				}
			}
		}
	}
	pc := Centroid2d(pts)
	SortCw(pts, pc)
	//fmt.Println("room points-",pts)
	return
}

//CorPolys builds either +ve or -ve offset corridor polygons
func (f *Flr) CorPolys(path []int, ldx int, neg bool)(){
	cw := f.Cwidth
	if cw == 0.0{
		cw = 750.0
	}
	if neg{
		cw = -cw
	}
	e2s := [][]Pt2d{}
	cps := make(map[Pt2d]bool)
	cpts := []Pt2d{}
	f.Polys = [][]Pt2d{}
	for i, jb := range path{
		p1 := f.Nodes[jb-1]
		if i != len(path)-1{
			je := path[i+1]
			p2 := f.Nodes[je-1]
			//eds = append(eds, []Pt2d{p1, p2})
			p3, p4 := EdgeOff2d(cw, p1, p2)
			e2s = append(e2s, []Pt2d{p3, p4})
			for _, pt := range []Pt2d{p1, p2, p3, p4}{
				if _, ok := cps[pt]; !ok{
					cps[pt] = true
					cpts = append(cpts, pt)
					//add p3/p4 to new corridor edges		
				} 
			}
		}
	}
	SortCw(cpts, Centroid2d(cpts))
	//fmt.Println("corridor points->",cpts)
	// fn = f.Clone()
	// fn.Sqrd = true
	//get points of intersection of cpts with all room (rect) polygons
	for i, rm := range f.Rooms{
		var rpoly []Pt2d
		var intersects bool
		//check if room intersects new corridor edges
		if i == ldx{
			p1, p2, p3, p4 := RectPts(rm.Origin, rm.End)
			//fmt.Println("living room points->",p1,p2,p3,p4)
			f.Polys = append(f.Polys, []Pt2d{p1, p2, p3, p4})				
			rpoly = []Pt2d{p1, p2, p3, p4}
		} else {	
			for _, ed := range e2s{
				if ed[0].InRect(rm.Origin, rm.End) || ed[1].InRect(rm.Origin,rm.End){
					intersects = true
				} 
			}
			if intersects{
				//fmt.Println("room-",rm.Name,"intersekto corpus corridore->")
				
				//p1, p2, p3, p4 := RectPts(rm.Origin, rm.End)
				
				//fmt.Println(rm.Name,"points->",p1,p2,p3,p4)
				//now kompute intersektions
				pts := rm.CorInt(cpts, e2s)
				rpoly = pts
				f.Polys = append(f.Polys, pts)
			} else {
				
				//fmt.Println("room-",rm.Name,"is not intersektomen->")
				p1, p2, p3, p4 := RectPts(rm.Origin, rm.End)
				//fmt.Println(rm.Name,"points->",p1,p2,p3,p4)
				
				f.Polys = append(f.Polys, []Pt2d{p1, p2, p3, p4})				
				rpoly = []Pt2d{p1, p2, p3, p4}
			}
		}
		f.Rooms[i].Polys = append(f.Rooms[i].Polys,rpoly)
	}
	f.Labels = append(f. Labels,"corridor")
	f.Polys = append(f.Polys, cpts)
	crm := Flr{
		Polys:[][]Pt2d{cpts},
		Name:"corridor",
	}
	f.Rooms = append(f.Rooms, &crm)
	return
}

//PolyRmap generates a room connectivity map from room polygons
func (f *Flr) PolyRmap(){
	f.Nmap = make(map[Pt2d][][]int)
	f.Emap = make(map[Tupil][][]int)
	f.Nodes = []Pt2d{}
	//fmt.Println("generating treemap from polygons (of all things)->")
	for i := range f.Rooms{
		edges := [][]Pt2d{}
		f.Rooms[i].Edges = []Tupil{}
		f.Rooms[i].Mid = Pt2d{0,0}
		var xmin, xmax, ymin, ymax float64
		f.Rooms[i].Area = 0.0
		f.Rooms[i].Area = 0.0
		xmin = -1.0; ymin = -1.0
		nn := float64(len(f.Polys[i]))
		for j, p := range f.Polys[i]{
			f.Rooms[i].Mid.X += p.X
			f.Rooms[i].Mid.Y += p.Y
			if xmin == -1.0{
				xmin = p.X
				ymin = p.Y
			}
			if xmin > p.X{
				xmin = p.X
			}
			if ymin > p.Y{
				ymin = p.Y
			}
			if xmax  < p.X{
				xmax = p.X
			}
			if ymax < p.Y{
				ymax = p.Y
			}
			p2 := f.Polys[i][0]
			if j < len(f.Polys[i])-1{
				p2 = f.Polys[i][j+1]
			}
			edges = append(edges, []Pt2d{p, p2})
			if _, ok := f.Nmap[p]; !ok{
				idx := len(f.Nodes) + 1
				f.Nodes = append(f.Nodes, p)
				f.Nmap[p] = make([][]int,3)
				f.Nmap[p][0] = []int{idx}
				f.Nmap[p][1] = []int{i}
				f.Nmap[p][2] = []int{1}
				if (p.X == f.Origin.X) || (p.X == f.End.X) || (p.Y == f.Origin.Y) || (p.Y == f.End.Y){
					f.Nmap[p][2] = []int{2}
				}
			} else {
				if !IntInVec(f.Nmap[p][1],i){
					f.Nmap[p][1] = append(f.Nmap[p][1],i)
				}
			}
		}
		
		f.Rooms[i].Mid.X = f.Rooms[i].Mid.X/nn
		f.Rooms[i].Mid.Y = f.Rooms[i].Mid.Y/nn
		f.Rooms[i].Bbo = Pt2d{xmin, ymin}
		f.Rooms[i].Bbe = Pt2d{xmax, ymax}
		
		for _, e := range edges{
			//each edge in map is - 
			//exterior - left right top bottom - -1, -2, -3, -4
			ecls := f.ClassEd(e[0],e[1])
			jb := f.Nmap[e[0]][0][0]; je := f.Nmap[e[1]][0][0]
			edx := Edgedx(jb, je)
			if _, ok := f.Emap[edx]; !ok{
				f.Emap[edx] = make([][]int,3)
				f.Emap[edx][0] = []int{jb, je, ecls}
				f.Emap[edx][1] = []int{i}
				f.Emap[edx][2] = []int{1}
				f.Edges = append(f.Edges, edx)
				if ecls == 1{
					f.Iwalls = append(f.Iwalls, edx)
					if f.Labels[i] == "living"{
						//should be zero here
						f.Emap[edx][2] = []int{1}
					}
				}
			} else {
				if !IntInVec(f.Emap[edx][1],i){
					f.Emap[edx][1] = append(f.Emap[edx][1],i)	
					if f.Labels[i] == "living" && f.Emap[edx][0][2] == 1{
						//again, should be zero here
						f.Emap[edx][2] = []int{1}
					}
				}
			}
			f.Rooms[i].Edges = append(f.Rooms[i].Edges, edx)
		}
	}
	f.Rmap = make(map[int][]int)
	for i, room := range f.Rooms{
		f.Rmap[i] = []int{}
		for _, edx := range room.Edges{
			edge := f.Emap[edx]
			// jb := edge[0][0]; je := edge[0][1]; ecls := edge[0][2]
			for _, nbr := range edge[1]{
				if nbr != i && !IntInVec(f.Rmap[i],nbr){
					f.Rmap[i] = append(f.Rmap[i],nbr)
				}
			}
		}
	}
	//if f.Verbose{fmt.Println("now komputing edge overlaps")}
	for i, r1 := range f.Rooms{
		// fmt.Println("checking room-",f.Labels[i])
		nbrs := f.Rmap[i]
		for _, edx := range r1.Edges{
			jb := f.Emap[edx][0][0]
			je := f.Emap[edx][0][1]
			a := f.Nodes[jb-1]
			b := f.Nodes[je-1]
			for j, r2 := range f.Rooms{
				if j != i{	
					//fmt.Println("checking against-",f.Labels[j])
					if !IntInVec(nbrs, j){	
						for _, edy := range r2.Edges{
							j1 := f.Emap[edy][0][0]
							j2 := f.Emap[edy][0][1]
							c := f.Nodes[j1-1]
							d := f.Nodes[j2-1]
							if EdgeOverlap(a,b,c,d){
								f.Rmap[i] = append(f.Rmap[i], j)
								if !IntInVec(f.Emap[edy][1],i){
									f.Emap[edy][1] = append(f.Emap[edy][1],i)
								}
								if !IntInVec(f.Emap[edx][1],j){
									f.Emap[edx][1] = append(f.Emap[edx][1],j)
								}
							}
						}
					}
				}
			}
		}
	}
	//remove duplicates (baah)
	for i, nbrs := range f.Rmap{
		ns := []int{}
		for _, nbr := range nbrs{
			if !IntInVec(ns, nbr){
				ns = append(ns, nbr)
			}
		}
		f.Rmap[i] = ns
	}
	f.Sqrd = true
	if f.Verbose{
		txtplot, _ := f.Draw()
		fmt.Println(txtplot)
	}	
}

//PolyGrid returns a grid rep. of a floor (from squarify output)
//all rooms are rects btw
func (f *Flr) PolyGrid()(grid [][]int){
	dx := f.Gx
	if dx == 0.0{
		switch f.Units{
			case "mm":
			dx = 900.0 
			case "in":
			dx = 36.0
			case "ft":
			dx = 3.0
			case "m":
			dx = 1.0
		}
		f.Gx = dx
	}
	dy := dx
	nc := int(math.Round(f.Width/dx))
	nr := int(math.Round(f.Height/dy))
	//fmt.Println(nr, nc)
	grid = make([][]int,nr)
	for i := range grid {
		grid[i] = make([]int, nc)
	}
	//set start point of ray at 2.0*f.End.X, 2.0*f.End.Y
	p0 := Pt2d{2.0*f.End.X, 2.0*f.End.Y}
	for idx, rm := range f.Rooms{		
		rb := int(math.Round(rm.Origin.Y/dy))
		re := int(math.Round(rm.End.Y/dy))
		cb := int(math.Round(rm.Origin.X/dx))
		ce := int(math.Round(rm.End.X/dx))
		p1, p2, p3, p4 := RectPts(rm.Origin, rm.End)
		//fmt.Println("at room->",idx, "label->",rm.Name,"rb, re, cb, ce",rb, re, cb, ce)
		poly := []Pt2d{p1, p2, p3, p4}
		for i := rb; i < re; i++ {
			for j := cb; j < ce; j++ {
				if i <= nr && j <= nc{
					xc := dx * float64(j) + dx/2.0
					yc := dy * float64(i) + dy/2.0
					pc := Pt2d{xc, yc}
					if pc.InPoly(poly,p0){
						grid[i][j] = idx + 1
					}
				}
			}
		}	   
	}
	//txtplot := Plotgrid(grid, dx, dy)
	//fmt.Println(txtplot)
	return
}

//CorGen generates a corridor for in-room connectivity
func (f *Flr) CorGen()(err error){
	//list of unconnected rooms
	uncon := []int{}
	ucmap := make(map[int]bool)
	upts := []Pt2d{}
	ldx := -1
	for i, room := range f.Labels{
		if !f.Lcmap[i]{
			// fmt.Println("LIVING IS NOT KONNECT->",room)
			uncon = append(uncon, i)
			ucmap[i] = false
			xc := (f.Rooms[i].Origin.X + f.Rooms[i].End.X)/2.0
			yc := (f.Rooms[i].Origin.Y + f.Rooms[i].End.Y)/2.0
			upts = append(upts, Pt2d{xc, yc})
		}
		if room == "living"{
			ldx = i
			f.Ldx = i
		}
	}
	graph := make(map[int][]int)
	//list of starting points (nodes connected to living room)
	strts := make(map[int]bool)
	for _, edx := range f.Iwalls{
		val := f.Emap[edx]
		jb := val[0][0]
		je := val[0][1]
		p1 := f.Nodes[jb-1]
		p2 := f.Nodes[je-1]
		if val[2][0] == 1{	
			if _, ok := graph[jb]; !ok{
				graph[jb] = []int{} 
			}
			if !IntInVec(graph[jb], je){
					graph[jb] = append(graph[jb],je)
			}
			if IntInVec(f.Nmap[p1][1],ldx){
					if _, ok := strts[jb]; !ok{
						strts[jb] = true
					}
			} else if f.Nodes[jb-1].InRect(f.Rooms[ldx].Origin,f.Rooms[ldx].End){
				strts[jb] = true
			}
			
			if _, ok := graph[je]; !ok{
				graph[je] = []int{} 
			}
			if !IntInVec(graph[je], jb){
				graph[je] = append(graph[je],jb)
			}
			if IntInVec(f.Nmap[p2][1],ldx){
				if _, ok := strts[je]; !ok{
					strts[je] = true
				} 
			} else if f.Nodes[je-1].InRect(f.Rooms[ldx].Origin,f.Rooms[ldx].End){
				strts[je] = true
			}	
		}
	}
	// fmt.Println("list of starting nodes->",strts)
	var start int
	var sdist float64
	for n1 := range strts{
		p1 := f.Nodes[n1-1]
		dist := nodehstic(p1, upts)
		// for _, p2 := range upts{
		// 	pd := p2.Sub(p1)
		// 	dist += pd.Length()
		// }
		if sdist == 0{
			start = n1
			sdist = dist
		} else if sdist > dist{
			sdist = dist
			start = n1
		}
	}
	
		
	var pq Pque
	cfrm := make(map[int]int)
	csf := make(map[int]float64)
	cfrm[start] = -1
	csf[start] = 0.0
	pq = append(pq, &Item{Tup:Tuple{start,0}, Pri:0.0})
	heap.Init(&pq)
	iter := 0
	goal := -1
	//start djk/astar loop
	for len(pq) > 0 && iter == 0{
		current := heap.Pop(&pq).(*Item).Tup.I
		stopcon := true
		prev := cfrm[current]
		if prev != -1{
			edx := Edgedx(current, prev)
			for _, rm := range f.Emap[edx][1]{
				if val, ok := ucmap[rm]; ok{
					if !val{
						ucmap[rm] = true
					}
				}
			}
		}
		for _, val := range ucmap{
			if !val{
				stopcon = false
			}
		}
		if stopcon{
			goal = current
			iter = -1
			break
		}
		nbrs := graph[current]
		for _, next := range nbrs{
			psub := f.Nodes[current-1].Sub(f.Nodes[next-1])
			costn := psub.Length()
			newcost := csf[current] + costn
			if _, ok := csf[next]; !ok || newcost < csf[next]{
				csf[next] = newcost
				//CHANGE HEURISTIC HERE (make it a method)
				priority := newcost + nodehstic(f.Nodes[next-1],upts)		
				//priority := newcost + f.NodeH(f.Nodes[next-1],upts,ucmap)
				heap.Push(&pq, &Item{Tup:Tuple{I:next,J:0},Pri:priority})
				cfrm[next] = current
			}
		}
	}
	//build path
	p := []int{}
	current := goal
	if _, ok := cfrm[goal]; !ok{
		return 
	}
	p = append(p, current)
	for{
		current = cfrm[current]
		stopcon := current == start
		p = append(p, current)
		if stopcon{
			break
		}
	}
	path := make([]int, len(p))
	for i, val := range p{
		path[len(p)-1-i] = val
	}
	psimp := f.SimpCor(ldx, path, uncon)
	switch f.Cgrid{
		case false:
		f.CorPolys2(psimp, ldx, false)
		f.PolyRmap()
		case true:
		//GAAAAAAAAAH
		//fmt.Println("APATHE->",psimp)
		grid := f.Gridpath(psimp, ldx+1)
		//txtplot := Plotgrid(grid, f.Gx, f.Gx)
		//fmt.Println(txtplot)
		rmap, _, _, _, _ := LoutGen(len(f.Labels)+1,len(grid[0]),len(grid), grid, f.Gx,f.Gx, []float64{},[]float64{})
		//outstr := PltLout(rmap)
		//fmt.Println(outstr)
		//fmt.Println(grid)
		//polys := f.GridPoly(grid)
		polys, e := RmapPoly(rmap)
		if e != nil{
			fmt.Println(ColorRed, e, ColorReset)
			err = e
			return
		}
		f.Polys = polys
		f.Labels = append(f.Labels, "corridor")
		//txtplot, err := f.DrawPolys()
		//if err != nil{
		//	fmt.Println(ColorRed, err, ColorReset)
		//	return
		//}
		//fmt.Println(txtplot)
		f.Grid = grid
		//now check intwalls, extwalls, room connections
	}
	return
}

//Gridpath marks a corridor grid based on a path
func (f *Flr) Gridpath(path []int, ldx int)([][]int){
	grid := f.PolyGrid()
	cdx := len(f.Labels)+1
	for i, i1 := range path[:len(path)-1]{
		
		i2 := path[i+1]
		p1 := f.Nodes[i1-1]
		p2 := f.Nodes[i2-1]
		//now r1, r2, c1, c2 are start and end indices
		rb := int(math.Round((p1.Y)/f.Gx))
		cb := int(math.Round((p1.X)/f.Gx))
		re := int(math.Round((p2.Y)/f.Gx))
		ce := int(math.Round((p2.X)/f.Gx))
		switch{
			case p1.X == p2.X:
			//move along y, ie, rb - re
			for j := rb; j <= re; j++{
				if j < len(grid){
					grid[j][cb] = cdx
				}
				
			}
			case p1.Y == p2.Y:
			//move along x, cb - ce
			for k := cb; k <= ce; k++{
				if k < len(grid[0]){
					grid[rb][k] = cdx
				}
			}
		}
	}
	
	//check if starting corridor cell borders/is connected to living room
	p1 := f.Nodes[path[0]-1]
	i1 := int(math.Round((p1.Y)/f.Gx))
	j1 := int(math.Round((p1.X)/f.Gx))
	nbrs := GridNbrs(grid, Tuple{i1,j1})
	nbrcon := false
	//fmt.Println("start cell-",i1, j1, "grid-",grid[i1][j1],"ldx-",ldx)
	
	if grid[i1][j1] == ldx{
		nbrcon = true
	} else {
		for _, nbr := range nbrs{
			if grid[nbr.I][nbr.J] == ldx{
				grid[i1][j1] = cdx
				nbrcon = true
			}
		}
	}
	
	
	//fmt.Println("loop done, nbr con-",nbrcon)
	if !nbrcon{
		f.ConCor(ldx, cdx, grid,Tuple{i1,j1})
	}
	//fmt.Println(ColorRed,"HYAAAAR",ColorReset)
	//txtplot := Plotgrid(grid, f.Gx, f.Gx)
	//fmt.Println(txtplot)
	return grid
	
}

//sigh. ConCor connects the corridor to living room in a grid
//run astar loop again 
func (f *Flr) ConCor(ldx, cdx int, grid [][]int, start Tuple)([][]int){
	path := GridPathBasic(f.Gx, grid, start, ldx)
	for _, val := range path{
		grid[val.I][val.J] = cdx
	}
	return grid
}

//SimpCor simplifies a floor's corridor/hall path
func (f *Flr) SimpCor(ldx int, path, uncon []int)(psimp []int){
	//first plot floor with nodes
	//txtplot, _ := f.Draw()
	//fmt.Println(txtplot)
	//fmt.Println("PATHE->",path)
	if len(path) < 3{
		return 
	}
	cw := f.Cwidth
	if cw == 0.0{cw = 750.0}
	//mark start and end as okay
	psimp = append(psimp, path[0])
	for i, j1 := range path[:len(path)-2]{
		j2 := path[i+1]
		j3 := path[i+2]
		if EdgeOverlap(f.Nodes[j1-1],f.Nodes[j2-1],f.Nodes[j2-1],f.Nodes[j3-1]){
			if cls, _ := EdgeInt(f.Nodes[j1-1],f.Nodes[j2-1],f.Nodes[j2-1],f.Nodes[j3-1]); cls == "collinear"{
				//fmt.Println(ColorRed,"errore here, removing center node",j2,ColorReset)	
				//pok[j2] = false
			} else {
				
				psimp = append(psimp, j2)
			}
		} else {
			psimp = append(psimp, j2)
		}
	}
	psimp = append(psimp, path[len(path)-1])
	//fmt.Println("new PATHe-",psimp)
	return

}

//SetTol rounds all room nodes to f.Tol
func (f *Flr) SetTol(){
	for _, room := range f.Rooms{
		room.Origin.SetTol(f.Tol)
		room.End.SetTol(f.Tol)
	}
}

//CorGrid generates a corridor given a floor grid
func (f *Flr) CorGrid(grid [][]int)([][]int){
	//list of unconnected rooms
	uncon := []int{}
	ucmap := make(map[int]bool)
	upts := []Pt2d{}
	ldx := -1
	for i, room := range f.Labels{
		if !f.Lcmap[i]{
			// fmt.Println("LIVING IS NOT KONNECT->",room)
			uncon = append(uncon, i+1)
			ucmap[i+1] = false
			xc := (f.Rooms[i].Origin.X + f.Rooms[i].End.X)/2.0
			yc := (f.Rooms[i].Origin.Y + f.Rooms[i].End.Y)/2.0
			upts = append(upts, Pt2d{xc, yc})
		}
		if room == "living"{
			ldx = i+1
		}
	}
	//the graph is the grid?
	//graph := make(map[int][]int)
	//check 4 extreme cells of living room
	//start from the cell with lowest dist
	p1, p2, p3, p4 := RectPts(f.Rooms[ldx].Origin,f.Rooms[ldx].End)
	var spt Pt2d
	var sdist float64
	for _, pt := range []Pt2d{p1,p2,p3,p4}{
		dist := nodehstic(pt, upts)
		if sdist == 0{
			spt = pt
			sdist = dist
		} else if sdist > dist{
			sdist = dist
			spt = pt
		}
	}
	//find the cell connected to this node in grid
	rb := int(math.Round(spt.Y/f.Gx))
	cb := int(math.Round(spt.X/f.Gx))
	stopcon := false
	var start Tuple
	for _, r := range []int{rb-1, rb, rb+1}{
		if stopcon{break}
		for _, c := range []int{cb-1, cb, cb+1}{
			if stopcon{break}
			if grid[r][c] == ldx{
				stopcon = true
				start = Tuple{r, c}
			}
		}
	}
	txtplot := Plotgrid(grid, f.Gx, f.Gx)
	fmt.Println(txtplot)
	fmt.Println(ColorCyan,"starting astar loop",ColorReset)
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
		stopcon := true
		prev := cfrm[current]
		if prev.I != -1{
			//not start, now check cell conn with other rooms
			nbrs := GridNbrs(grid, current)
			for _, nbr := range nbrs{
				rdx := grid[nbr.I][nbr.J]
				if _, ok := ucmap[rdx]; ok{
					//fmt.Println("cell-",nbr,"connected to-",rdx)
					ucmap[rdx] = true
				}
			}
		}
		for _, val := range ucmap{
			if !val{
				stopcon = false
			}
		}
		if stopcon{
			goal = current
			iter = -1
			break
		}
		nbrs := GridNbrs(grid, current)
		pcur := gridpt(f.Gx, current)
		for _, next := range nbrs{
			pnxt := gridpt(f.Gx, next)
			pdiff := pnxt.Sub(pcur)
			costn := pdiff.Length()
			//costn = 0.0
			newcost := csf[current] + costn + nodehstic(pnxt, upts)
			
			if _, ok := csf[next]; !ok || newcost < csf[next]{
				csf[next] = newcost
				priority := newcost //+ nodehstic(pnxt, upts)		
				heap.Push(&pq, &Item{Tup:next,Pri:priority})
				cfrm[next] = current
			}
		}
	}
	//build path
	p := []Tuple{}
	current := goal
	if _, ok := cfrm[goal]; !ok{
		fmt.Println(ColorRed,"ERRORE IN GOALE",ColorReset)
		return grid 
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
	path := make([]Tuple, len(p))
	for i, val := range p{
		path[len(p)-1-i] = val
	}
	cdx := len(f.Labels) + 1
	//fmt.Println("CDX-",cdx)
	for _, val := range path{
		grid[val.I][val.J] = cdx
	}
	txtplot = Plotgrid(grid, f.Gx, f.Gx)
	fmt.Println(txtplot)
	return grid
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

//FlrGen generates the room connectivity graph and edges
func (f *Flr) FlrGen()(err error){
	//f.Flrprint(true)
	if f.Round{
		f.SetTol()
	}
	if f.Verbose{
		GPlotFloors(f, true)
	}
	f.SqrRmap()
	switch f.Opt{
		case 1:
		//craft dir algo
		grid := f.CraftDir()
		f.Labels = append(f.Labels, "corridor")
		
		rmap, _, _, _, _ := LoutGen(len(f.Labels),len(grid[0]),len(grid), grid, f.Gx,f.Gx, []float64{},[]float64{})
		f.Polys, err = RmapPoly(rmap)
		if err != nil{
			fmt.Println(ColorRed, err, ColorReset)
			return
		}
		
		fmt.Println(grid)
		f.Grid = grid
		return
	}
	switch f.Togrid{
		case false:
		f.CorGen()
		rmap, nmap, ptmap, wmap, pts := LoutGen(len(f.Labels),len(f.Grid[0]),len(f.Grid), f.Grid, f.Gx,f.Gx, []float64{},[]float64{})
		f.WallGen(rmap, nmap, ptmap, wmap, pts)
		//f.ColGrid(rmap, nmap, ptmap, wmap, pts)
		//f.DrawPlan(rmap, nmap, ptmap, wmap, pts)
		case true:
		grid := f.PolyGrid()
		f.CorGrid(grid)
	}
	//f.Craft()
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

//WallGen generates walls/door openings
func (f *Flr) WallGen(rmap map[int]*Rm, nmap map[Pt][]*Wall, ptmap map[Pt][]int, wmap map[Tupil][]int, pts []Pt)(err error){
	
	//corridor index is last
	cdx := len(f.Labels)
	var kdx, ldx int
	for i , lbl := range f.Labels{
		if lbl == "kitchen"{kdx = i+1}
		if lbl == "living"{ldx = i+1}
	}
	for idx := range rmap{
		rm := rmap[idx]
		if _, ok := rm.Walls[cdx]; ok {
			cwall := rm.Walls[cdx][0]
			cwall.Typ = -1
			jb := ptmap[*cwall.Pb][0]; je := ptmap[*cwall.Pe][0]
			edx := EdgeDx(jb, je)
			wmap[edx][0] = 2
		} else {
			//check if conn. to living room
			if _, ok := rm.Walls[ldx]; ok {
				cwall := rm.Walls[ldx][0]
				cwall.Typ = -1
				jb := ptmap[*cwall.Pb][0]; je := ptmap[*cwall.Pe][0]
				edx := EdgeDx(jb, je)
				wmap[edx][0] = 2
			}
		}
		if idx == kdx{
			if _, ok := rm.Walls[ldx]; ok{
				for i, wall := range rm.Walls[ldx]{
					rm.Walls[ldx][i].Typ = 3
					jb := ptmap[*wall.Pb][0]; je := ptmap[*wall.Pe][0]
					edx := EdgeDx(jb, je)
					wmap[edx][0] = 3
				}
				for i, wall := range rmap[ldx].Walls[kdx]{
					rmap[ldx].Walls[kdx][i].Typ = 3
					jb := ptmap[*wall.Pb][0]; je := ptmap[*wall.Pe][0]
					edx := EdgeDx(jb, je)
					wmap[edx][0] = 3
				}
			}
		}
	}

	//mark walls on either side of each external node as non window (edx = 4)
	for pt, vec := range ptmap{
		switch vec[1]{
			case 0:
			case 1,2,3,4:
			if _, ok := nmap[pt]; ok{
				for _, wall := range nmap[pt]{
					jb := ptmap[*wall.Pb][0]; je := ptmap[*wall.Pe][0]
					edx := EdgeDx(jb, je)
					if _, ok := wmap[edx]; ok{
						wmap[edx][0] = 4
					}
				}
			}
			case 5:
			if len(nmap[pt]) > 2{	
				for _, wall := range nmap[pt]{
					jb := ptmap[*wall.Pb][0]; je := ptmap[*wall.Pe][0]
					edx := EdgeDx(jb, je)
					if _, ok := wmap[edx]; ok{
						if wmap[edx][0] == 0{wmap[edx][0] = 4}
					}
				}
			}
		}
	}
	//dump to struct
	//pts
	f.Pts = make([][]float64, len(pts))
	for i, pt := range pts{
		f.Pts[i] = []float64{pt.X,pt.Y}
	}
	//walls
	for edx, vec := range wmap{
		f.Walls = append(f.Walls, []int{edx.I, edx.J})
		f.Wvec = append(f.Wvec, vec)
	}
	//fmt.Println("nodemap-",nmap)

	//now join walls? traverse
	
	
	
	plotdx := 0
	txtplot, _ := f.DrawWalls(rmap, nmap, ptmap, wmap, pts, plotdx)
	if f.Term == "dumb"{fmt.Println(txtplot)}

	return
}

//FlrLay. FlrLay lays out a (res) floor using the squarified treemap algo
func (f *Flr) FlrLay()(err error){
	f.Init()
	//fmt.Println("f.Areas, f.Labels",f.Areas, f.Labels)
	//fmt.Println("f.Origin, f.End, f.Area",f.Origin, f.End, f.Area)
	switch f.Bloc{
		case false:
		f.Flrarea()
		if f.Name == ""{f.Name = "base"}
		f.Isroot = true
		//r := []float64{6,6,4,3,2,2,1}
		FlrPln(f, f.Areas, f.Labels)
		if f.Round{
			for _, room := range f.Rooms{
				room.SetTol()
			}
		}
		f.FlrGen()	
		case true:
		_ = f.ResBmap()
		
	}
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


//CorPolys2 builds either +ve or -ve offset corridor polygons using ctessum/geom clipping funcs
func (f *Flr) CorPolys2(path []int, ldx int, neg bool)(){
	cw := f.Cwidth
	if cw == 0.0{
		cw = 750.0
	}
	if neg{
		cw = -cw
	}
	e2s := [][]Pt2d{}
	cps := make(map[Pt2d]bool)
	cpts := []Pt2d{}
	f.Polys = [][]Pt2d{}
	for i, jb := range path{
		p1 := f.Nodes[jb-1]
		if i != len(path)-1{
			je := path[i+1]
			p2 := f.Nodes[je-1]
			//eds = append(eds, []Pt2d{p1, p2})
			p3, p4 := EdgeOff2d(cw, p1, p2)
			e2s = append(e2s, []Pt2d{p3, p4})
			for _, pt := range []Pt2d{p1, p2, p3, p4}{
				if _, ok := cps[pt]; !ok{
					cps[pt] = true
					cpts = append(cpts, pt)
					//add p3/p4 to new corridor edges		
				} 
			}
		}
	}
	SortCw(cpts, Centroid2d(cpts))
	// fn = f.Clone()
	// fn.Sqrd = true
	//get points of intersection of cpts with all room (rect) polygons
	for i, rm := range f.Rooms{
		var rpoly []Pt2d
		//var intersects bool
		//get living room as is
		p1, p2, p3, p4 := RectPts(rm.Origin, rm.End)
		if i == ldx{
			//fmt.Println("living room points->",p1,p2,p3,p4)
			f.Polys = append(f.Polys, []Pt2d{p1, p2, p3, p4})				
			rpoly = []Pt2d{p1, p2, p3, p4}
		} else {
			//get poly from points
			r1 := []Pt2d{p1, p2, p3, p4}
			rpoly = RmSub(r1, cpts)
			f.Polys = append(f.Polys, rpoly)				
		}
		f.Rooms[i].Polys = append(f.Rooms[i].Polys,rpoly)
	}
	f.Labels = append(f. Labels,"corridor")
	f.Polys = append(f.Polys, cpts)
	crm := Flr{
		Polys:[][]Pt2d{cpts},
		Name:"corridor",
	}
	f.Rooms = append(f.Rooms, &crm)
	return
}

/*
// var di, dj int
	// p1 := f.Nodes[path[0]-1]
	//rl := int(math.Round(p1.Y/f.Gx))
	//cl := int(math.Round(p1.X/f.Gx))
	// stopcon := false
	// for _, r := range []int{rl-1, rl, rl+1}{
	// 	if stopcon{break}
	// 	for _, c := range []int{cl-1, cl, cl+1}{
	// 		if stopcon{break}
	// 		if grid[r][c] == ldx{
	// 			stopcon = true
	// 			di = r - rl
	// 			dj = c - cl
	// 		}
	// 	}
	// }
	


//EvalRes evaluates a residential floor plan
func (f *Flr) EvalRes(grid [][]int)(rmap map[int]*Rm, cost float64){
	dirs := getdirvec(f.Facing)
	rmap, _ = LoutGen(len(f.Labels),len(grid[0]),len(grid), grid, f.Gx, f.Gx, []float64{},[]float64{})
	// outstr := PltLout(rmap)
	// fmt.Println(outstr)
	//fmt.Println("now rating residential plan")
	rnmap := make(map[int][]int)
	switch f.Opt{
		case 0:
		for i, rm := range rmap{
			rnmap[i] = []int{}
			//fmt.Println("at room-",i,f.Labels[i-1])
			for _, edge := range rm.Edges{
				if !IntInVec(rnmap[i], edge){
					rnmap[i] = append(rnmap[i],edge)
					//fmt.Println("edge-",edge)
				}
			}
		}
		for i, nbrs := range rnmap{
			rlbl := strings.Split(f.Labels[i-1],"-")[0]
			var eon, won, non, son bool
			for _, nbr := range nbrs{
				var lbl string
				switch nbr{
					case -1:
					//left
					lbl = dirs[0]
					case -2:
					//right
					lbl = dirs[1]
					case -3:
					//bottom
					lbl = dirs[3]
					case -4:
					//top
					lbl = dirs[2]
					default:
					lbl = f.Labels[i-1]
				}
				if nbr < 0{
					switch lbl{
						case "e":
						eon = true
						case "w":
						won = true
						case "n":
						son = true
						case "s":
						son = true
					}
				}
			}
			switch rlbl{
				case "bath":
				switch{
					case non && won:
					cost -= 100.0
					case son && eon:
					cost -= 100.0
					case non && eon:
					cost += 100.0
				}
				case "kitchen":
				switch{
					case son && eon:
					cost -= 100.0
					case non && eon:
					cost += 100.0
					case son && won:
					cost += 100.0
					case non:
					cost += 100.0
				}
				
			}
		}
	}
	return
}

//Craft applies the craft algo to a floor given a scoring func.
func (f *Flr) Craft(){
	if f.Facing == ""{
		f.Facing = "e"
	}
	
	grid := f.PolyGrid()

	rmap, cmin := f.EvalRes(grid)


	combos := make(map[string][]int)

	for i:= 1; i <= len(rmap); i++ {
		if i == len(f.Labels)-1{
			//is corridor
			continue
		}
		for j := i+1; j <= len(rmap); j++ {
			if j == len(f.Labels){
				//is corridor
				continue
			}
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

	
	costs := make(map[string]float64, len(combos))
	gridz := make(map[string][][]int, len(combos))

	rcent := make([]*Pt,len(f.Labels))
	for i:=1; i <=len(rmap); i++ {
		rcent[i-1] = rmap[i].Centroid
	}
	var minidx string
	for idx, combo := range combos {
		gridnew := make([][]int,len(grid))
		for i := range grid {
			gridnew[i] = make([]int,len(grid[i]))
			copy(gridnew[i],grid[i])
		}
		_, gridnew = Rswap(combo, rmap, rcent, gridnew, f.Gx, f.Gx)
		_, cnew := f.EvalRes(gridnew)
		costs[idx] = cnew
		gridz[idx] = gridnew
		if cmin > costs[idx]{
			cmin = costs[idx]
			minidx = idx
		}
	}
	rezstring := ""
	swpz := map[int]string{0:"area",1:"adj"}
	for idx, cost := range costs{
		if idx == minidx {
			rezstring += "dis MINIMUM COST\n"
		}
		rezstring += ColorCyan
		rezstring += fmt.Sprintf("combo - swap %v and %v cause %v total cost %v\n",combos[idx][0],combos[idx][1],swpz[combos[idx][2]],cost)
		rezstring += ColorPurple
		outstr := Plotgrid(gridz[idx],f.Gx,f.Gx)
		rezstring += outstr
	}
	fmt.Println(rezstring)	
}
// //Rmdat stores input data for squarify
// type Rmdat struct{
	Rooms []string
	Areas []float64
	Dims  []float64
	Units string
// }



*/
