package barf

import (
	//"fmt"
	"math"
)

//TODO
//genfun generates x and y funicular coords from the bending moment diagram
//leet ? page ?
func genfunic(xs, frcs, rise []float64)(){
	
}

//genarch generates x and y arch coords (from formula) from span and rise
func genarch(span, pspc, rise float64)(xs, ys []float64){
	var x float64
	h := rise
	l2 := math.Pow(span, 2.0)
	nbz := math.Round(span/pspc)
	dx := math.Round(span/nbz)
	xs = append(xs, x)
	for i := 0; i < int(nbz); i++{
		x += math.Round(dx)
		xs = append(xs, x)
	}
	ys = make([]float64,len(xs))
	for i, x := range xs{
		xa := span/2.0 - x 
		ya := 0.0 - 4.0 * h * xa * xa/l2
		ys[i] = math.Round(ya + rise)
	}
	return
	
}

//GenFunTruss generates a funicular/bowstring type truss
func GenFunTruss(t *Trs2d)(err error){
	xs , ys := genarch(t.Span, t.Purlinspc, t.Rise)
	var coords [][]float64
	var ms [][]int
	var bcns, tcns []int
	var prlnspc, rftrl float64
	var nrs, ngs int
	if t.Typ == 0{t.Typ = 1}
	bcnodes := make(map[int]bool)
	tcnodes := make(map[int]bool)
	nodes := make(map[Pt]int)
	var ndivs, y0 float64
	prlnspc = t.Purlinspc
	span := t.Span
	ndivs = span/prlnspc
	nrs = int(ndivs)
	//if t.Spam{log.Println("nrafters, purlinspc->",nrafters, prlnspc)}
	y0 = t.Height
	//y = y0
	ngs = 4
	apx := []float64{0,0}
	rftrl = t.Span
	var idx, nt1, nt2, napx int
	switch t.Cfg{
		case 1,2:
		//three hinged (funicular) flat top chord truss
		//first add bc nodes - arch coords
		for i, x := range xs{
			idx++
			y := ys[i] + y0
			if apx[1] < y{
				napx = idx
				apx[0] = x
				apx[1] = y
			}
			pt := Pt{x,y}
			coords = append(coords, []float64{x,y})
			bcns = append(bcns, idx)
			bcnodes[idx] = true
			nodes[pt] = idx
		}
		//then tc nodes. mark nodes adjacent to apex
		for _, x := range xs{
			y := apx[1]
			if x == apx[0]{
				nt1 = idx
				nt2 = idx + 1
			} else {
				idx++
				pt := Pt{x,y}
				coords = append(coords, []float64{x,y})
				tcns = append(tcns, idx)
				tcnodes[idx] = true
				nodes[pt] = idx
			}
		}
		//then join in between as per t.Cfg
		
		bc1 := []int{}
		//join -1 on left, + 1 on right or vice versa
		//add bottom chord
		for _, jb := range bcns{
			if jb != napx{
				bc1 = append(bc1, jb)
			}
			if _, ok := bcnodes[jb+1]; ok{
				//this can never be continuous (?) has to be pin ended
				mvec := []int{jb, jb+1, 1, 1, 3}
				ms = append(ms, mvec)
			}
		}
		for i, jb := range tcns{
			var mvec []int
			if jb != nt1{	
				if _, ok := tcnodes[jb+1]; ok{	
					switch t.Tcr{
						case 0:
						mvec = []int{jb, jb+1, 1, 2, 3} 
						case 1:
						mvec = []int{jb, jb+1, 1, 2, 0}
					}
					ms = append(ms, mvec)
				}
			}
			je := bc1[i]
			//add vertical struts/webs
			mvec = []int{je, jb, 1, 3, 3}
			ms = append(ms, mvec)
			x := coords[je-1][0]
			//add slings/ties
			//1 - slope left, 2 - slope right
			//t.Brc = 2; both ways
			switch{
				case t.Brc == 2:
				//brace both ways
				//get intersection of jb, opp. and je, -opp
				//FindIntInf(x1,y1,x2,y2,x3,y3,x4,y4 float64)(par bool, px, py float64)
				x1 := coords[je-1][0]; y1 := coords[je-1][1]
				if jb != nt1 && i != len(tcns)-1{
					n2 := tcns[i+1]
					p2 := coords[n2-1]
					x2 := p2[0]; y2 := p2[1]
					x3 := coords[jb-1][0]; y3 := coords[jb-1][1]
					n4 := bc1[i+1]
					p4 := coords[n4-1]
					x4 := p4[0]; y4 := p4[1]
					_, px, py := FindIntInf(x1, y1, x2, y2, x3, y3, x4, y4)
					coords = append(coords, []float64{px, py})
					je1 := len(coords)
					for _, vec := range [][]int{{je, je1},{je1, n2},{jb, je1},{je1,n4}}{
						mvec = []int{vec[0],vec[1],1,2,3}
						ms = append(ms, mvec)
					}
				}
				case t.Cfg == 1:	
				if i > 0 && i < len(tcns)-1{
					if x < apx[0]{
						jn := tcns[i-1]
						mvec = []int{je, jn, 1, 4, 3}
					} else {
						jn := tcns[i+1]
						mvec = []int{je, jn, 1, 4, 3}
					}
					ms = append(ms, mvec)
				}
				case t.Cfg == 2:				
				if jb != nt1 && jb != nt2{
					if x < apx[0]{
						jn := tcns[i+1]
						mvec = []int{je, jn, 1, 4, 3}
					} else {
						jn := tcns[i-1]
						mvec = []int{je, jn, 1, 4, 3}
					}
					ms = append(ms, mvec)
				}
			}
		}
		switch t.Tcr{
			case 0:
			mvec := []int{nt1, napx, 1, 2, 3}
			ms = append(ms, mvec)
			mvec = []int{napx, nt2, 1, 2, 3}
			ms = append(ms, mvec)
			case 1:
			mvec := []int{nt1, napx, 1, 2, 0}
			ms = append(ms, mvec)
			mvec = []int{napx, nt2, 1, 2, 0}
			ms = append(ms, mvec)
		}
	}
	tcns = append(tcns, napx)
	t.Coords = coords
	t.Ms = ms
	t.Bcns = bcns
	t.Tcns = tcns
	t.Purlinspc = prlnspc
	t.Rftrl = rftrl
	t.Nrs = nrs
	t.Ngs = ngs
	t.Mod.Coords = t.Coords; t.Mod.Mprp = t.Ms
	t.Mod.Supports = [][]int{{bcns[0],-1,-1},{bcns[len(bcns)-1],0,-1}}
	PlotGenTrs(coords, ms)
	return
}
