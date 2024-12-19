package barf

import (
	"fmt"
	"math"
	"container/heap"
)

type Tuple struct{
	I   int
	J   int
}

type Item struct{
	Tup   Tuple
	Pri   float64
	Idx   int
}

type Gst struct{
	At    Tuple
	Path  []Tuple
}

//A Pque implements heap.Interface and holds Items.
type Pque []*Item

func (pq Pque) Len() int { return len(pq) }

func (pq Pque) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Pri > pq[j].Pri
}

func (pq Pque) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Idx = i
	pq[j].Idx = j
	
}

func (pq *Pque) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Idx = n
	*pq = append(*pq, item)
}

func (pq *Pque) Pop() (interface{}) {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.Idx = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority of an item given an index.
func (pq *Pque) Update(item *Item, priority float64) {
	item.Pri = priority
	heap.Fix(pq, item.Idx)
}

//Init inits a grid 
func (g *Grid) Init(ni, nj int, vec, walls [][]int, haswt bool){
	g.Ni = ni
	g.Nj = nj
	if len(vec) > 0{
		//init as per given vec
		g.Vec = make([][]int, len(vec))
		for i, val := range vec{
			g.Vec[i] = make([]int, len(val))
			copy(g.Vec[i],val)
			if haswt{	
				for j := range val{
					g.Weights[Tuple{I:i,J:j}] = make(map[Tuple]float64)
				}
			}
		}
		
		
	} else {
		//init binary vec
		g.Vec = make([][]int, g.Ni)
		for i := 0; i < g.Ni; i++{
			g.Vec[i] = make([]int, g.Nj)
			for j := 0; j < g.Nj; j++{
				g.Vec[i][j] = 1
				if haswt{
					g.Weights[Tuple{I:i,J:j}] = make(map[Tuple]float64)
				}
			}
		}
	}
	g.Walls = make([][]int, len(walls))
	g.Weights = make(map[Tuple]map[Tuple]float64)
	if len(walls) > 0{	
		for i, wall := range walls{
			g.Walls[i] = make([]int,2)
			copy(g.Walls[i], wall)
			i := wall[0]; j := wall[1]
			if g.Ingrid(Tuple{I:i,J:j}){
				g.Vec[i][j] = 0
			}
		}
	}
}

//Printbin prints a binary grid
func (g *Grid) Printbin(start, goal Tuple, printz bool)(gstr [][]string){
	gstr = make([][]string, len(g.Vec))
	for i, vec := range g.Vec{
		gstr[i] = make([]string,len(vec))
		for j, val := range vec{
			switch val{
				case 0:
				gstr[i][j] = "#"
				case 5:
				gstr[i][j] = "$"
				default:
				gstr[i][j] = "."
			}
		} 
	}
	if g.Ingrid(start){
		gstr[start.I][start.J] ="A"
	}
	if g.Ingrid(goal){
		gstr[goal.I][goal.J] ="B"
	}
	if printz{
		for _, v := range gstr{
			for _, val := range v{
				fmt.Printf("%s",val)
			}
			fmt.Printf("%s","\n")
		}
	}
	return
}

func (g *Grid) Draw(start, goal Tuple, cfrm map[Tuple]Tuple){
	gstr := make([][]string, len(g.Vec))
	for i, vec := range g.Vec{
		gstr[i] = make([]string,len(vec))
		for j, val := range vec{
			switch val{
				case 0:
				gstr[i][j] = "#"
				case 5:
				gstr[i][j] = "$"
				default:
				gstr[i][j] = "."
			}
		} 
	}
	for cell, parent := range cfrm{
		if parent.I != -1 && parent.J != -1{
			gstr[cell.I][cell.J] = Stampstr(cell, parent)
		}
	}
	
	if g.Ingrid(start){
		gstr[start.I][start.J] ="A"
	}
	
	if g.Ingrid(goal){
		gstr[goal.I][goal.J] ="B"
	}
	for _, v := range gstr{	
		for _, val := range v{
			fmt.Printf("%s",val)
		}
		fmt.Printf("%s","\n")
	}
}

func (g *Grid) Getpath(start, goal Tuple, cfrm map[Tuple]Tuple)(path []Tuple){
	p := []Tuple{}
	current := goal
	if _, ok := cfrm[goal]; !ok{
		return 
	}
	p = append(p, current)
	for{
		current = cfrm[current]
		
		stopcon := (current.I == start.I) && (current.J == start.J)
		p = append(p, current)
		if stopcon{
			//fmt.Println("breeeking")
			break
		}
	}
	path = make([]Tuple, len(p))
	for i, val := range p{
		path[len(p)-1-i] = val
	}
	return
}

func (g *Grid) Printpath(start, goal Tuple, path []Tuple){	
	gstr := g.Printbin(start, goal, false)
	for _, cell := range path{
		gstr[cell.I][cell.J] = "x"
	}
	if g.Ingrid(start){
		gstr[start.I][start.J] ="A"
	}
	
	if g.Ingrid(goal){
		gstr[goal.I][goal.J] ="B"
	}
	for _, v := range gstr{	
		for _, val := range v{
			fmt.Printf("%s",val)
		}
		fmt.Printf("%s","\n")
	}
}

func (g *Grid) Ingrid(cell Tuple) (bool){
	return (0 <= cell.I && cell.I < g.Nj) && (0 <= cell.J && cell.J < g.Ni)
}

func (g *Grid) Nbrs90(cell Tuple) (nbrs []Tuple){
	x := cell.I
	y := cell.J
	ns := []Tuple{
		{I:x+1,J:y},
		{I:x-1,J:y},
		{I:x,J:y-1},
		{I:x,J:y+1},
	}
	if (x+y)%2 == 0{
		ns = []Tuple{	
			{I:x,J:y+1},
			{I:x,J:y-1},
			{I:x-1,J:y},
			{I:x+1,J:y},
		}
	}
	for _, nbr := range ns{
		if g.Ingrid(nbr){
			nbrs = append(nbrs,nbr)
		}
	}
	return
}

func (g *Grid) NbrsBlt(cell Tuple) (nbrs []Tuple){
	x := cell.I
	y := cell.J
	ns := []Tuple{
		{I:x+1,J:y},
		{I:x+1,J:y-1},
		{I:x+1,J:y+1},
	}
	for _, nbr := range ns{
		if g.Ingrid(nbr){
			nbrs = append(nbrs,nbr)
		}
	}
	return
}

func (g *Grid) Costidx(fc, tc Tuple) (cost float64){
	v1 := g.Vec[fc.I][fc.J]
	v2 := g.Vec[tc.I][tc.J]
	cost = float64(v1 + v2)
	// fmt.Println("cost from", fc,"to", tc, "->",cost)
	return
}

func (g *Grid) CostBlt(fc, tc Tuple) (cost float64){
	//not checking for fc == 0 || tc == 0
	cost = g.Vals[0]
	if fc.I != tc.I{		
		ps := g.Vals[1]
		gg := g.Vals[2] * 4.0
		cost -= math.Pow(ps,2)/gg
	}
	return
}

func (g *Grid) Hrstic(crnt, goal Tuple)(hcost float64){
	hcost = math.Abs(float64(goal.I - crnt.I)) + math.Abs(float64(goal.J-crnt.J))
	return
}

func Stampstr(cell, parent Tuple)(mark string){
	switch{
		case cell.I == parent.I:
		if cell.J == parent.J - 1{
			mark = ">"
		} else {
			mark = "<"
		}
		case cell.J == parent.J:
		if cell.I == parent.I - 1{
			mark = "v"
		} else {
			mark = "^"
		}
	}
	return
}

func BfsGrid(g Grid, start, goal Tuple, eex bool, sdir string)(cfrm map[Tuple]Tuple){
	cfrm = make(map[Tuple]Tuple)
	q := []Tuple{}
	q = append(q, start)
	cfrm[start] = Tuple{I:-1,J:-1}
	for len(q) > 0{
		current  := q[0]
		
		stopcon := (current.I == goal.I) && (current.J == goal.J)
		if stopcon && eex{
			//early exit
			return
		}
		q = q[1:]
		var nbrs []Tuple
		switch sdir{
			case "90":
			nbrs = g.Nbrs90(current)
			case "45":
			//etc
		}
		for _, next := range nbrs{
			if _, ok := cfrm[next]; !ok{
				if g.Vec[next.I][next.J] != 0{
					q = append(q, next)
					cfrm[next] = current
				}
			}
		}
	}
	return
}

//DjkGrid performs a uniform/heuristic based cost search (djikstra's algo or astar) on a grid
func DjkGrid(g Grid, start, goal Tuple, sdir string, astar bool)(cfrm map[Tuple]Tuple, csf map[Tuple]float64){
	cfrm = make(map[Tuple]Tuple)
	csf = make(map[Tuple]float64)
	cfrm[start] = Tuple{I:-1,J:-1}
	csf[start] = 0.0
	var pq Pque
	pq = append(pq, &Item{Tup:start, Pri:0.0})
	heap.Init(&pq)
	iter := 0
	for len(pq) > 0 && iter == 0{
		citm := heap.Pop(&pq).(*Item)
		//fmt.Println("at current-",current,"goal-",goal)
		current := citm.Tup
		stopcon := (current.I == goal.I) && (current.J == goal.J)
		if stopcon{
			iter = -1
			return
		}
		var nbrs []Tuple
		switch sdir{
			case "90":
			nbrs = g.Nbrs90(current)
			case "45":
			//etc
		}
		for _, next := range nbrs{	
			if g.Vec[next.I][next.J] == 1{	
				newcost := csf[current] + g.Costidx(current, next)
				if _, ok := csf[next]; !ok || newcost < csf[next]{
					csf[next] = newcost
					priority := newcost
					if astar{
						priority += g.Hrstic(next, goal)
					}
					heap.Push(&pq,&Item{Tup:next,Pri:priority})
					cfrm[next] = current
				} 
			}
		}
	}
	return

}

//BltGrid performs a uniform/heuristic based cost search (djikstra's algo or astar) on a bolt group vec/grid
//for calculation of net sec area = b - nd + Σp2/4g
func BltGrid(g Grid, start Tuple, astar bool)(goal Tuple, cfrm map[Tuple]Tuple, csf map[Tuple]float64){
	cfrm = make(map[Tuple]Tuple)
	csf = make(map[Tuple]float64)
	cfrm[start] = Tuple{I:-1,J:-1}
	if g.Vec[start.I][start.J] == 1{csf[start] = g.Vals[0]}
	var pq Pque
	pq = append(pq, &Item{Tup:start, Pri:0.0})
	heap.Init(&pq)
	iter := 0
	fmt.Println("NI, NJ-", g.Ni, g.Nj)
	for len(pq) > 0 && iter == 0{
		citm := heap.Pop(&pq).(*Item)
		//fmt.Println("at current-",current,"goal-",goal)
		current := citm.Tup
		stopcon := (current.J == g.Nj-1) 
		fmt.Println(ColorCyan,"at current",current,ColorReset)
		fmt.Println("cost so far-",csf[current])
		if stopcon{
			fmt.Println("stop con reached")
			fmt.Println("cost so far",csf[current])
			goal = current
			iter = -1
			return
		}
		nbrs := g.NbrsBlt(current)
		for _, next := range nbrs{
			fmt.Println("checking nbr",next)
			switch g.Vec[next.I][next.J]{
				case 1:
				newcost := csf[current] + g.CostBlt(current, next)
				if _, ok := csf[next]; !ok || newcost < csf[next]{
					csf[next] = newcost
					priority := newcost
					heap.Push(&pq,&Item{Tup:next,Pri:priority})
					cfrm[next] = current
				} 
			}
			
		}
	}
	return

}

//BltGsts returns neighboring states for a cell in a bolt group grid
func BltGsts(grid [][]int, st Gst)(nsts []Gst){
	//down, l 45, r 45
	cur := st.At
	mr := len(grid)
	mc := len(grid[0])
	dwn := Tuple{cur.I+1, cur.J}
	//fmt.Println(ColorRed, "current-",cur,"down-",dwn,"max row, max col",mr, mc,ColorReset)
	lft := Tuple{cur.I+1, cur.J-1}
	rgt := Tuple{cur.I+1, cur.J+1}
	if cur.I + 1 < mr{
		if cur.J < mc{
			dst := Gst{
				At: dwn,
				Path:append(st.Path, dwn),
			}
			dcon := cur.I + 2 < mr && grid[dwn.I+1][dwn.J] == 1
			if grid[dwn.I][dwn.J] == 1{
				nsts = append(nsts, dst)
			} else if dcon{
					ndn := Tuple{dwn.I+1, dwn.J}
					dst = Gst{
						At: ndn,
						Path: append(st.Path, ndn),
					}
					nsts = append(nsts, dst)
			}
		}
		if cur.J-1 >= 0{
			lst := Gst{
				At:lft,
				Path:append(st.Path,lft),
			}
			if grid[lft.I][lft.J] == 1{nsts = append(nsts, lst)}
		}
		if cur.J+1 < mc{
			rst := Gst{
				At:rgt,
				Path:append(st.Path,rgt),
			}
			if grid[rgt.I][rgt.J] == 1{nsts = append(nsts, rst)}
		}
	}
	return

}

//BltPaths returns all possible failure paths from a start cell of a bolt grid
func BltPaths(grid [][]int, start Tuple)(paths [][]Tuple){
	st := Gst{
		At: start,
		Path:[]Tuple{start},
	}
	var iter int
	sq := []Gst{st}
	mr := len(grid)
	for iter != -1{
		if len(sq) == 0{
			iter = -1
			break
		}
		st = sq[0]
		if len(sq) == 1{
			sq  = []Gst{}
		} else {
			sq = sq[1:]
		}
		cur := st.At
		scon1 := cur.I == mr -1 && grid[cur.I][cur.J] == 1
		scon2 := cur.I == mr -3 && grid[cur.I][cur.J] == 1 && grid[cur.I+1][cur.J] == 0 && grid[cur.I+2][cur.J] == 1
		scon3 := cur.I == mr -2 && grid[cur.I][cur.J] == 1 && grid[cur.I+1][cur.J] == 0
		switch{
			case scon1:
			paths = append(paths, st.Path)
			case scon2:
			nxt := Tuple{cur.I+2, cur.J}
			path := append(st.Path, nxt)
			paths = append(paths, path)
			case scon3:
			paths = append(paths, st.Path)
		}
		nsts := BltGsts(grid, st)
		if len(nsts) > 0{
			sq = append(sq, nsts...)
		} 
	}
	return
}

func DrawBltPaths(grid [][]int, path []Tuple){
	fmt.Println(ColorCyan)
	start := path[0]
	goal := path[len(path)-1]
	ni := len(grid)
	nj := len(grid[0])
	var g Grid
	g.Init(ni, nj, grid, [][]int{},false)
	g.Printpath(start, goal, path)
	fmt.Println(ColorReset)
}

func BltNsa(grid [][]int, path []Tuple, bw, tp, dia, ps, gg float64)(nsa float64){
	var nh, ns float64
	//start := path[0]
	//goal := path[len(path)-1]
	for i, loc := range path{
		if grid[loc.I][loc.J] == 1{
			nh += 1.0
		}
		if i != 0{
			prev := path[i-1]
			if prev.J != loc.J{
				ns += 1.0
			}
		}
	}
	gg = 4.0 * gg
	nsa = tp * (bw - nh * dia + ns * ps * ps/gg)
	return
}

//BltNsaNon finds the net sectional area for a non uniform pitch/gauge bolt group 
func BltNsaNon(grid [][]int, path []Tuple, bw, tp, dia float64, pss, ggs []float64)(nsa float64){
	var nh, ns float64
	//start := path[0]
	//goal := path[len(path)-1]
	for i, loc := range path{
		if grid[loc.I][loc.J] == 1{
			nh += dia
		}
		if i != 0{
			prev := path[i-1]
			if prev.J != loc.J{
				ps := pss[loc.J]
				gg := ggs[loc.I]
				ns += ps * ps/(4.0 * gg)
			}
		}
	}
	nsa = tp * (bw - nh + ns)
	return
}


//BltSecArea checks for a range of nj/2 start vals to find min net sec area 
func BltSecArea(grid [][]int, bw, tp, dia, ps, gg float64, pss, ggs []float64)(paths [][]Tuple, nsas []float64, nsamin float64, mindx int){
	nj := len(grid[0])
	for j := 0; j < nj; j++{
		if j > nj/2 + 2{break}
		start := Tuple{0, j}
		if grid[start.I][start.J] == 0{continue}
		pjs := BltPaths(grid, start)
		paths = append(paths, pjs...)
	}
	var nsa float64
	uniform := true
	if len(pss) == len(grid[0]) && len(ggs) == len(grid){
		uniform = false
	}
	for i, path := range paths{
		if uniform{
			nsa = BltNsa(grid, path, bw, tp, dia, ps, gg)
		} else {
			nsa = BltNsaNon(grid, path, bw, tp, dia, pss, ggs)
		}
		if i == 0{
			nsamin = nsa
		}
		if nsamin > nsa{
			nsamin = nsa
			mindx = i

		}
		nsas = append(nsas, nsa)
	}
	return
}


func TestBfs1(){
	walls := [][]int{{0, 22}, {0, 23}, {1, 22}, {1, 23}, {2, 22}, {2, 23}, {3, 4}, {3, 5}, {3, 22}, {3, 23}, {4, 4}, {4, 5}, {4, 14}, {4, 15}, {4, 22}, {4, 23}, {5, 4}, {5, 5}, {5, 14}, {5, 15}, {5, 22}, {5, 23}, {5, 24}, {5, 25}, {5, 26}, {6, 4}, {6, 5}, {6, 14}, {6, 15}, {6, 22}, {6, 23}, {6, 24}, {6, 25}, {6, 26}, {7, 4}, {7, 5}, {7, 14}, {7, 15}, {8, 4}, {8, 5}, {8, 14}, {8, 15}, {9, 4}, {9, 5}, {9, 14}, {9, 15}, {10, 4}, {10, 5}, {10, 14}, {10, 15}, {11, 4}, {11, 5}, {11, 14}, {11, 15}, {12, 14}, {12, 15}, {13, 14}, {13, 15}, {14, 14}, {14, 15}}
	ni := 30
	nj := 15
	g := Grid{}
	vec := [][]int{}
	g.Init(ni, nj, vec, walls, false)
	start := Tuple{I:7,J:8}
	goal := Tuple{I:2,J:17}
	fmt.Println("bfs->")
	fmt.Println(ColorRed)
	g.Printbin(start, goal, true)
	fmt.Println(ColorCyan)
	cfrm := BfsGrid(g, start,goal,true,"90")
	g.Draw(start, goal, cfrm)
	fmt.Println("")
	path := g.Getpath(start, goal, cfrm)
	if len(path) == 0{
		fmt.Println("ERRORE - > no path found")
	}
	fmt.Println(ColorYellow)
	g.Printpath(start, goal,path)
	fmt.Println(ColorReset)
}

func TestDjk1(){
	walls := [][]int{}
	ni := 10
	nj := 10
	g := Grid{}
	vec := [][]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 5, 5, 1, 1, 1, 1},
		{1, 1, 1, 1, 5, 5, 5, 1, 1, 1},
		{1, 1, 1, 1, 5, 5, 5, 5, 1, 1},
		{1, 1, 1, 5, 5, 5, 5, 5, 1, 1},
		{1, 1, 1, 5, 5, 5, 5, 5, 1, 1},
		{1, 1, 1, 1, 5, 5, 5, 1, 1, 1},
		{1, 0, 0, 0, 5, 5, 5, 1, 1, 1},
		{1, 0, 0, 0, 5, 5, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}
	g.Init(ni, nj, vec, walls, false)
	start := Tuple{I:4,J:1}
	goal := Tuple{I:3,J:8}
	fmt.Println("djk->")
	fmt.Println(ColorRed)
	g.Printbin(start, goal, true)
	fmt.Println(ColorCyan)
	cfrm, _ := DjkGrid(g, start,goal,"90",false)
	path := g.Getpath(start, goal, cfrm)
	if len(path) == 0{
		fmt.Println("ERRORE - > no path found")
	}
	fmt.Println(ColorYellow)
	g.Printpath(start, goal,path)
	fmt.Println(ColorReset)
}

func TestAstar1(){
	walls := [][]int{
		{12,2},{12,3},{12,4},{12,5},{12,6},{12,7},{12,8},{12,9},{12,10},{12,11},{12,12},{11,12},{10,12},{9,12},{8,12},{7,12},{6,12},{5,12},{4,12},{3,12},{2,12},{2,11},{2,10},{2,9},{2,8},{2,7},{2,6},{2,5}}
	ni := 15
	nj := 15
	g := Grid{}
	vec := [][]int{}
	g.Init(ni, nj, vec, walls, false)
	start := Tuple{I:12,J:0}
	goal := Tuple{I:2,J:14}
	fmt.Println("astar->")
	fmt.Println(ColorRed)
	g.Printbin(start, goal, true)
	fmt.Println(ColorCyan)
	cfrm, _ := DjkGrid(g, start,goal,"90",true)
	// g.Draw(start, goal, cfrm)
	path := g.Getpath(start, goal, cfrm)
	if len(path) == 0{
		fmt.Println("ERRORE - > no path found")
		return
	}
	fmt.Println(ColorYellow)
	g.Printpath(start, goal,path)
	fmt.Println(ColorReset)
}

// func main(){
// 	//TestBfs1()
// 	TestAstar1()
// 	//TestDjk1()
// }


	// // Some items and their priorities.
	// items := []Tuple{
	// 	{I:1,J:17,P:3}, {I:2,J:3,P:4}, {I:17,J:21,P:5},
	// }
	
	// // Create a priority queue, put the items in it, and
	// // establish the priority queue (heap) invariants.
	// pq := make(Pque, len(items))
	// for i, cell := range items{
	// 	pq = append(pq,cell)
	// 	pq[i].Idx = i
		
	// }
	// heap.Init(&pq)
	// item := Tuple{
	// 	I:1,
	// 	J:1,
	// 	P:11,
	// }
	// heap.Push(&pq, item)
	// // Take the items out; they arrive in decreasing priority order.
	// for pq.Len() > 0 {
	// 	item := heap.Pop(&pq).(Tuple)
	// 	fmt.Println(item.P)
	// }	// Some items and their priorities.
	// items := []Tuple{
	// 	{I:1,J:17,P:3}, {I:2,J:3,P:4}, {I:17,J:21,P:5},
	// }
	
	// // Create a priority queue, put the items in it, and
	// // establish the priority queue (heap) invariants.
	// pq := make(Pque, len(items))
	// for i, cell := range items{
	// 	pq = append(pq,cell)
	// 	pq[i].Idx = i
		
	// }
	// heap.Init(&pq)
	// item := Tuple{
	// 	I:1,
	// 	J:1,
	// 	P:11,
	// }
	// heap.Push(&pq, item)
	// // Take the items out; they arrive in decreasing priority order.
	// for pq.Len() > 0 {
	// 	item := heap.Pop(&pq).(Tuple)
	// 	fmt.Println(item.P)
	// }
