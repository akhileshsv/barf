package barf

import (
	"fmt"
	"math"
)

type Fcflr struct {
	Coords   [][]float64
	Col      []int
	Xbm      [][]int
	Ybm      [][]int
	Slb      [][]int
	Slbl     []string
	Slbval   [][]float64
	Slbtyp   int
	DL, LL   float64
	Xbms     map[Tupil][]Tupil
	Xbmpts   map[Tupil][]int
	Bms      map[Tupil][][]float64
	Cols     map[Tupil][]Tupil
	Bmprp    map[Tupil][]int
	Bmdx     map[Tupil]int
	Ndx      map[Pt2d]int
	Nmap     map[int][]Tupil
	Ncmap    map[int]bool
	Term     string
	Bmat     string
	Ftyp     string
}

type StlFlr struct {
	Xs      []float64
	Ys      []float64
	Zs      []float64  
	DL      float64
	LL      float64
	Ftyp    string
	Coords  [][]float64
	Cols    []int
	Bmxs    [][]int
	Bmys    [][]int
	Jspc    float64
	Ohng    float64
	Jspcs   [][]float64
	
}

func BmTupil(i, j int)(bdx Tupil){
	if i < j{
		bdx = Tupil{i,j}
	} else {
		bdx = Tupil{j,i}
	}
	return
}

//slbbmsreturns left and right beams for a slab (if type 2 returns lbms)
func SlbBms(idxs []int,coords [][]float64, pts []Pt2d,sdir string, slbtyp int)(lbm, rbm []Tupil){
	pc := Centroid2d(pts)
	p0 := Pt2d{0,0}
	switch slbtyp{
		case 1:	
		for i, p1 := range pts{
			if i != len(pts)-1{
				idx := idxs[i]
				idx2 := idxs[i+1]
				bm := BmTupil(idx, idx2)//Tupil{idx, idx2}
				p2 := pts[i+1]
				//fmt.Println("p1, p2, pc-",p1, p2, pc)
				ort := P3orient(p1, p2, pc)
				//fmt.Println("ort-",ort)
				ort2 := P3orient(p1, p2, p0)
				//fmt.Println("ort2-",ort2)
				ort = ort * ort2
				//fmt.Println("final ort-",ort)
				switch sdir{
					case "x":
					if p1.Y == p2.Y{
						switch{
							case ort <= 0:
							lbm = append(lbm, bm)
							case ort > 0:
							rbm = append(rbm, bm)
						}
					}
					case "y":
					if p1.X == p2.X{
						switch{
							case ort <= 0:
							lbm = append(lbm, bm)
							case ort > 0:
							rbm = append(rbm, bm)
						}
					}
					case "xy":
					//wtf
				}
			}
		}
		case 2:
	}
	return
}

//Load loads an fc floor struct
func (f *Fcflr) Load(){
	f.Bms = make(map[Tupil][][]float64)
	f.Bmprp = make(map[Tupil][]int)
	f.Xbms = make(map[Tupil][]Tupil)
	f.Xbmpts = make(map[Tupil][]int)
	f.Bmdx = make(map[Tupil]int)
	f.Nmap = make(map[int][]Tupil)
	f.Ncmap = make(map[int]bool)
	f.Ndx = make(map[Pt2d]int)
	f.Cols = make(map[Tupil][]Tupil)
	//create node maps
	for i, pt := range f.Coords{
		f.Nmap[i+1] = []Tupil{}
		f.Ncmap[i+1] = false
		pt := Pt2d{X:pt[0],Y:pt[1]}
		f.Ndx[pt] = i+1
	}
	for _, idx := range f.Col{
		f.Ncmap[idx] = true
	}
	bdx := 1
	for _, xbm := range f.Xbm{
		bm := BmTupil(xbm[0],xbm[1])
		f.Bms[bm] = [][]float64{}
		f.Bmprp[bm] = []int{0}
		f.Bmdx[bm] = bdx
		f.Nmap[xbm[0]] = append(f.Nmap[xbm[0]], bm)
		f.Nmap[xbm[1]] = append(f.Nmap[xbm[1]], bm)
		bdx++
	}
	for _, ybm := range f.Ybm{
		bm := BmTupil(ybm[0],ybm[1])
		f.Bms[bm] = [][]float64{}
		f.Bmprp[bm] = []int{1}
		f.Bmdx[bm] = bdx
		f.Nmap[ybm[0]] = append(f.Nmap[ybm[0]], bm)
		f.Nmap[ybm[1]] = append(f.Nmap[ybm[1]], bm)
		bdx++
	}
	//read slabs
	for i, slb := range f.Slb{
		cds := [][]float64{}
		pts := []Pt2d{}
		for _, idx := range slb{
			cds = append(cds, f.Coords[idx-1])
			pts = append(pts, Pt2d{f.Coords[idx-1][0],f.Coords[idx-1][1]})
		}
		cds = append(cds, cds[0])
		pts = append(pts, pts[0])
		area := PolyArea(cds)/1e6
		slb = append(slb, slb[0])
		lbm, rbm := SlbBms(slb, cds, pts, f.Slbl[i], f.Slbtyp)
		switch f.Slbtyp{
			case 1:
			//global one way slab
			sdl := f.DL*area/2.0; sll := f.LL*area/2.0
			//all indices are 1 (change this?)
			dl1 := []float64{1.0, 3.0, sdl, 0, 0, 0, 1}
			ll1 := []float64{1.0, 3.0, sll, 0, 0, 0, 2}
			for _, bm := range lbm{				
				//idx := float64(f.Bmdx[bm])
				//dl1[0] = idx; ll1[0] = idx
				f.Bms[bm] = append(f.Bms[bm],dl1)
				f.Bms[bm] = append(f.Bms[bm],ll1) 
			}
			for _, bm := range rbm{
				//idx := float64(f.Bmdx[bm])
				//dl1[0] = idx; ll1[0] = idx
				f.Bms[bm] = append(f.Bms[bm],dl1)
				f.Bms[bm] = append(f.Bms[bm],ll1) 
			}
			case 2:
			//check if ly < lx, etx?
			
		}
	}
	//load primary beams with secondary beam point loads
	for _, ybm := range f.Ybm{
		a := NewPt2d(f.Coords[ybm[0]-1])
		b := NewPt2d(f.Coords[ybm[1]-1])
		ydx := Tupil{ybm[0],ybm[1]}
		//get pts of intersection with xbms
		for _, xbm := range f.Xbm{
			c := NewPt2d(f.Coords[xbm[0]-1])
			d := NewPt2d(f.Coords[xbm[1]-1])
			cls, px := EdgeInt(a,b,c,d)
			if cls == "cross"{
				xdx := Tupil{xbm[0],xbm[1]}
				f.Xbms[xdx] = append(f.Xbms[xdx],ydx)
				f.Xbmpts[xdx] = append(f.Xbmpts[xdx],f.Ndx[px])
			}
		}
	}
	return
}

//Init inits a StlFlr struct
func (f *StlFlr) Init(){
	//one way joist-girder system
	//gen coords
	for _, x := range f.Xs{
		for _, y := range f.Ys{
			f.Coords = append(f.Coords,[]float64{x,y})
			f.Cols = append(f.Cols, len(f.Coords))
		}
	}
	fmt.Println("coords->\n",f.Coords)
	fmt.Println("cols->\n",f.Cols)
	//get joist spacing per bay
	if f.Jspc == 0.0{f.Jspc = 1800}
	f.Jspcs = make([][]float64,len(f.Ys))
	for i := range f.Ys{
		if i == 0{
			continue
		}
		dy := f.Ys[i] - f.Ys[i-1]
		ndiv := math.Ceil(dy/f.Jspc)
		fmt.Println("bay-",i,"span-",dy,"jspc-",dy/ndiv,"mm girder nloads-",ndiv-1)
		for j := 1; j < int(ndiv); j++{
			f.Jspcs[i-1] = append(f.Jspcs[i-1],float64(j)*dy/ndiv)
		}
		fmt.Println("bay girder loads",f.Jspcs[i-1])
	}
	//gen bmxs (girders)
	
	
}
