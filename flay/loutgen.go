package barf


//mergewalls merges internal and external wall edges

//LoutGen generates a room map and wall map for a grid of rooms
func LoutGen(nrooms,nx,ny int, grid [][]int, dx,dy float64, dimvecx, dimvecy []float64) (map[int]*Rm, map[Pt][]*Wall, map[Pt][]int, map[Tupil][]int, []Pt){
	//generate room map and wall map 
	var x0, y0, x, y, xmax, ymax float64
	rmap := make(map[int]*Rm)
	nodemap := make(map[Pt][]*Wall)
	//add a ptmap and wmap for ease
	ptmap := make(map[Pt][]int)
	wmap := make(map[Tupil][]int)
	pts := []Pt{}
	var rm *Rm
	prevdir := make([]int,4)
	if dx != 0.0 {
		xmax = dx * float64(nx)
	} else {
		for _, xc := range dimvecx{
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
			
			//adding ptmap shit here
			//ptmap[pt] = [pdx, loc(ext/int), degree, etx]
			for _, pt := range []Pt{p1, p2, p3, p4}{
				if _, ok := ptmap[pt]; !ok{
					pdx := len(ptmap) + 1
					ploc := 0
					switch{
						case pt.X == 0.0 && pt.Y == 0.0:
						ploc = 1
						case pt.X == 0.0 && pt.Y == ymax:
						ploc = 2
						case pt.X == xmax && pt.Y == 0.0:
						ploc = 4
						case pt.X == xmax && pt.Y == ymax:
						ploc = 3
						case (pt.X == xmax || pt.X == 0.0 || pt.Y == ymax || pt.Y == 0.0):
						ploc = 5
					}
					ptmap[pt] = []int{pdx,ploc,0,0}
					pts = append(pts, pt)
				} 
			}
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

			// n1 := ptmap[p1][0]; n2 := ptmap[p2][0]; n3 := ptmap[p3][0]; n4 := ptmap[p4][0]
			//left edge
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
					wleft.Typ = 1
					nodemap[p1] = append(nodemap[p1],&wleft)
					nodemap[p2] = append(nodemap[p2],&wleft)
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
				nodemap[p2] = append(nodemap[p2], &wbottom)
				nodemap[p3] = append(nodemap[p3], &wbottom)
			} else {
				prevdir[3] = grid[i+1][j]
				if prevdir[3] == room {
					wbottom.Typ = -1
				} else {
					wbottom.Typ = 1
					nodemap[p2] = append(nodemap[p2], &wbottom)
					nodemap[p3] = append(nodemap[p3], &wbottom)
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
				wall := walls[idx]
				jb := ptmap[*wall.Pb]; je := ptmap[*wall.Pe]
				edx := Edgedx(jb[0], je[0])
				if wall.Typ != -1{	
					if _, ok := wmap[edx]; !ok{
						wmap[edx] = []int{wall.Typ,dir,0,0}	
					}
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
			for _, nbr := range prevdir{
				if !IntInVec(rm.Nbrs, nbr){rm.Nbrs = append(rm.Nbrs, nbr)}
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
	return rmap, nodemap, ptmap, wmap, pts
}
