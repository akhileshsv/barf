package barf

import (
	"log"
	"fmt"
	"math"
)

//GenParTruss generates parallel chord truss models
//if height > 0.0 - bent
//member groups 1 to 5
func GenParTruss(t *Trs2d) (err error){
	var coords [][]float64
	var ms [][]int
	var bcns, tcns []int
	var prlnspc, rftrl float64
	var nrs, ngs int
	ngs = 4
	bcnodes := make(map[int]bool)
	tcnodes := make(map[int]bool)
	mcnodes := make(map[int]bool)
	mcns := []int{}
	var nrafters, x, y, y0, xstep float64
	y0 = t.Height
	prlnspc = t.Purlinspc
	rftrl = t.Span
	if rftrl > prlnspc {
		nrafters = math.Ceil(rftrl / prlnspc)
		if (t.Cfg == 8 || t.Cfg == 9) && int(nrafters)%2 != 0{
			nrafters += 1.0
		}
		prlnspc = math.Round(rftrl / nrafters)
	} else {
		prlnspc = rftrl
		nrafters = 1
	}
	
	nrs = int(nrafters)
	fmt.Println("nrafters->",nrs)
	if t.Endrat == 0{t.Endrat = 10.0}
	if t.Depth == 0{
		t.Depth = t.Span/t.Endrat
	}
	switch{
		case t.Cfg < 6 :
		xstep = math.Round(t.Span/nrafters)
		y = y0
		coords = append(coords, []float64{x, y})
		for i := 0; i < nrs; i++ {
			x += xstep
			coords = append(coords, []float64{x, y})
		}
		x = 0.0
		y = y0 + t.Depth

		coords = append(coords, []float64{x, y})
		for i := 0; i < nrs; i++ {
			x += xstep
			coords = append(coords, []float64{x, y})
		}
		for i := 1; i <= nrs+1; i++ {
			bcnodes[i] = true
			bcns = append(bcns, i)
		}
		for i := nrs + 2; i <= len(coords); i++ {
			tcns = append(tcns, i)
			tcnodes[i] = true
		}
		for _, jb := range bcns {
			if _, ok := bcnodes[jb+1]; ok {
				switch t.Bcr{
					case 0:
					ms = append(ms, []int{jb, jb + 1, 1, 1, 3})
					case 1:
					//bending for bottom chord
					ms = append(ms, []int{jb, jb + 1, 1, 1, 0})
				}
			}
		}
		for _, jb := range tcns {
			if _, ok := tcnodes[jb+1]; ok {
				switch t.Tcr{
					case 0:
					ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
					case 1:
					//bending 
					ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
				}
			}
		}
		case t.Cfg < 11:
		xstep = math.Round(t.Span/nrafters)
		y = y0
		coords = append(coords, []float64{x, y})
		for i := 0; i < nrs; i++ {
			x += xstep
			coords = append(coords, []float64{x, y})
		}
		x = 0.0
		y = y0 + t.Depth/2.0
		
		coords = append(coords, []float64{x, y})
		for i := 0; i < nrs; i++ {
			x += xstep
			coords = append(coords, []float64{x, y})
		}
		x = 0.0
		y = y0 + t.Depth
		
		coords = append(coords, []float64{x, y})
		for i := 0; i < nrs; i++ {
			x += xstep
			coords = append(coords, []float64{x, y})
		}
		for i := 1; i <= nrs+1; i++ {
			bcnodes[i] = true
			bcns = append(bcns, i)
		}
		for i := nrs + 2; i <= 2 * nrs + 2; i++ {
			mcns = append(mcns, i)
			mcnodes[i] = true
		}
		
		for i := 2*nrs + 3; i <= len(coords); i++ {
			tcns = append(tcns, i)
			tcnodes[i] = true
		}
		for i, jb := range bcns {
			if _, ok := bcnodes[jb+1]; ok {
				switch t.Bcr{
					case 0:
					ms = append(ms, []int{jb, jb + 1, 1, 1, 3})
					case 1:
					//bending for bottom chord
					ms = append(ms, []int{jb, jb + 1, 1, 1, 0})
				}
			}
			ms = append(ms, []int{jb, mcns[i],1,3,3})
		}
		for i, jb := range tcns {
			if _, ok := tcnodes[jb+1]; ok {
				switch t.Tcr{
					case 0:
					ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
					case 1:
					//bending 
					ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
				}
			}
			ms = append(ms, []int{mcns[i],jb,1,3,3})
		}		

	}
	switch t.Cfg {
	case 1:
		//forward brace
		ngs = 4
		for i, jb := range bcns{
			ms = append(ms, []int{jb, tcns[i],1,3,3})
			if _, ok := tcnodes[jb+nrs+2]; ok{
				ms = append(ms, []int{jb, tcns[i+1],1,4,3})
			}
		}
	case 2:
		//backward
		ngs = 4
		for i, jb := range bcns{
			ms = append(ms, []int{jb, tcns[i], 1, 3, 3})
			if _, ok := tcnodes[jb+nrs]; ok{
				ms = append(ms, []int{jb, tcns[i-1],1,4,3})
			}
		}
	case 3:
		//tri forward
		for i, jb := range bcns{
			ms = append(ms, []int{jb, tcns[i], 1, 3, 3})
			if jb % 2 == 0{
				if _, ok := tcnodes[jb+nrs+2]; ok{
					ms = append(ms, []int{jb, tcns[i+1],1 ,4, 3})
				}
				if _, ok := tcnodes[jb+nrs]; ok{
					ms = append(ms, []int{jb, tcns[i-1], 1, 4, 3})
				}			
			}
		}
	case 4:
		//tri backward
		for i, jb := range bcns{
			ms = append(ms, []int{jb, tcns[i], 1, 3, 3})
			if jb % 2 != 0{
				if _, ok := tcnodes[jb+nrs+2]; ok{
					ms = append(ms, []int{jb, tcns[i+1],1,4,3})
				}
				if _, ok := tcnodes[jb+nrs]; ok{
					ms = append(ms, []int{jb, tcns[i-1],1,4,3})
				}			
			}
		}
	case 5:
		//virendeel girder
		for i, jb := range bcns{
			ms = append(ms, []int{jb, tcns[i],1 ,3, 3})
		}
	case 6:
		//k forward
		for i, jb := range bcns{
			if _, ok := mcnodes[jb+nrs+2]; ok{
				ms = append(ms, []int{jb, mcns[i+1],1,4,3})
				ms = append(ms, []int{mcns[i+1],tcns[i],1,4,3})
			}
		}
	case 7:
		//k backward
		for i, jb := range bcns{
			if _, ok := mcnodes[jb+nrs]; ok{
				ms = append(ms, []int{jb, mcns[i-1],1,4,3})
				ms = append(ms, []int{mcns[i-1],tcns[i],1,4,3})
			}
		}
	case 8:
		//k (even rafters)
		for i, jb := range bcns{
			switch {
			case i < nrs/2:
				if _, ok := mcnodes[jb+nrs]; ok{
					ms = append(ms, []int{jb, mcns[i-1],1,4,3})
					ms = append(ms, []int{mcns[i-1],tcns[i],1,4,3})
				}
			case i == nrs/2:
				if _, ok := mcnodes[jb+nrs+2]; ok{
					ms = append(ms, []int{jb, mcns[i+1],1,4,3})
					ms = append(ms, []int{mcns[i+1],tcns[i],1,4,3})
				}
				if _, ok := mcnodes[jb+nrs]; ok{
					ms = append(ms, []int{jb, mcns[i-1],1,4,3})
					ms = append(ms, []int{mcns[i-1],tcns[i],1,4,3})
				}
			default:
				if _, ok := mcnodes[jb+nrs+2]; ok{
					ms = append(ms, []int{jb, mcns[i+1],1,4,3})
					ms = append(ms, []int{mcns[i+1],tcns[i],1,4,3})
				}

			}
		}
	case 11:
		//bailey truss
		xstep = math.Round(t.Span/nrafters)
		y = y0
		coords = append(coords, []float64{x, y})
		for i := 0; i < nrs; i++ {
			x += xstep
			coords = append(coords, []float64{x, y})
		}
		x = 0.0
		y = y0 + t.Depth/2.0
		
		coords = append(coords, []float64{x, y})
		for i := 0; i < nrs/2; i++ {
			x += xstep * 2.0
			coords = append(coords, []float64{x, y})
		}
		x = 0.0
		y = y0 + t.Depth
		
		coords = append(coords, []float64{x, y})
		for i := 0; i < nrs; i++ {
			x += xstep
			coords = append(coords, []float64{x, y})
		}
		for i := 1; i <= nrs+1; i++ {
			bcnodes[i] = true
			bcns = append(bcns, i)
		}
		for i := nrs + 2; i <= nrs + nrs/2 + 2; i++ {
			mcns = append(mcns, i)
			mcnodes[i] = true
		}
		for i := nrs + nrs/2 + 3; i <= len(coords); i++ {
			tcns = append(tcns, i)
			tcnodes[i] = true
		}
		for i, jb := range bcns {
			if _, ok := bcnodes[jb+1]; ok {
				ms = append(ms, []int{jb, jb + 1, 1, 1,3})
			}
			if i % 2 != 0{
				ms = append(ms, []int{jb, mcns[i/2+1],1,4,3})
				ms = append(ms, []int{jb, mcns[i/2],1,4,3})
			} else {
				ms = append(ms, []int{jb, mcns[i/2],1,3,3})
			}
		}
		for i, jb := range tcns {
			if _, ok := tcnodes[jb+1]; ok {
				ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
			}
			if i % 2 != 0{
				ms = append(ms, []int{mcns[i/2+1],jb,1,4,3})
				ms = append(ms, []int{mcns[i/2],jb,1,4,3})
			} else {
				ms = append(ms, []int{mcns[i/2],jb,1,3,3})
			}
		}
		
	case 9:
		//k forward wall frame
		ngs = 5
		for i, jb := range bcns{
			if _, ok := mcnodes[jb+nrs+2]; ok{
				ms = append(ms, []int{jb, mcns[i+1],1,4,3})
				ms = append(ms, []int{mcns[i+1],tcns[i],1,4,3})
			}
		}
		for i, jb := range mcns{
			if i < len(mcns) - 1{
				je := mcns[i+1]
				ms = append(ms, []int{jb, je, 1, 5, 3})
			} 
		}
	case 10:
		//k backward wall frame
		ngs = 5
		for i, jb := range bcns{
			if _, ok := mcnodes[jb+nrs]; ok{
				ms = append(ms, []int{jb, mcns[i-1],1,4,3})
				ms = append(ms, []int{mcns[i-1],tcns[i],1,4,3})
			}
		}
		
		for i, jb := range mcns{
			if i < len(mcns) - 1{
				je := mcns[i+1]
				ms = append(ms, []int{jb, je, 1, 5, 3})
			} 
		}

	}
	t.Coords = coords
	t.Ms = ms
	t.Bcns = bcns
	t.Tcns = tcns
	t.Purlinspc = prlnspc
	t.Rftrl = rftrl
	t.Nrs = nrs
	t.Ngs = ngs
	t.Mod.Coords = t.Coords; t.Mod.Mprp = t.Ms
	t.Mod.Supports = [][]int{{bcns[0],-1,-1},{bcns[len(bcns)-1],-1,-1}}

	log.Println("plotting truss->")

	PlotGenTrs(t.Coords, t.Ms)
	//HERE DO IF t.BENT, etc
	return
}
