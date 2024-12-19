package barf

import (
	"log"
	"fmt"
	"math"
)

//GenTrapTruss generates trapezoidal truss models
func GenTrapTruss(t *Trs2d) (err error){
	var coords [][]float64
	var ms [][]int
	var bcns, tcns []int
	var prlnspc, rftrl float64
	var nrs, ngs int

	bcnodes := make(map[int]bool)
	tcnodes := make(map[int]bool)
	var nrafters, x, y, y0, xstep, ystep float64
	prlnspc = t.Purlinspc
	var span float64
	var tmul int
	if t.Typ == 3{
		span = t.Span
		tmul = 1
	} else {
		tmul = 2
		span = t.Span/2.0
	}
	
	rise := t.Rise
	if rise == 0.0{
		switch{
			case t.Tang == 0.0 && t.Slope == 0.0:
			err = fmt.Errorf("roof rise/slope not specified")
			return
			case t.Slope == 0.0:
			rise = span * math.Tan(t.Tang * math.Pi/180)
			t.Slope = span/rise
			case t.Tang == 0.0:
			rise = span/t.Slope
		}
	}
	if t.Slope == 0.0{
		t.Slope = span/rise
	}
	rftrl = math.Round(math.Sqrt(math.Pow(span, 2) + math.Pow(rise, 2)))
	if rftrl > prlnspc {
		nrafters = math.Ceil(rftrl / prlnspc)
		prlnspc = math.Round(rftrl / nrafters)
	} else {
		prlnspc = rftrl 
		nrafters = 1
	}
	if t.Spam{log.Println("rafter len, span, rise->",rftrl, span, rise)}
	nrs = int(nrafters)
	fmt.Println("nrafters->", nrs)
	if nrs < 2{return}
	if t.Endrat == 0{t.Endrat = 10.0}
	if t.Depth == 0{
		t.Depth = span/t.Endrat
	}
	fmt.Println("depth-",t.Depth)
	if t.Bent{y0 = t.Height}
	y = y0
	switch t.Cfg {
	case 1:
		switch t.Typ{
			case 3:			
			//el typo trapo
			xstep = math.Round(span/nrafters)
			coords = append(coords, []float64{x, y})
			for i := 0; i < nrs; i++ {
				x += xstep
				coords = append(coords, []float64{x, y})
			}
			x = 0.0
			y = t.Depth + y0
			coords = append(coords, []float64{x,y})
			for i := 0; i < nrs; i++ {
				x += xstep
				y += xstep / t.Slope
				coords = append(coords, []float64{x, y})
			}
			for i := 1; i < nrs + 2; i++{
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
			ms = append(ms, []int{tcns[0], bcns[0], 1, 3, 3})

			for i, node := range bcns[1 : len(bcns)-1] {
				ms = append(ms, []int{node, tcns[i]+1,1,3,3})
			}
			ms = append(ms, []int{tcns[len(tcns)-1], bcns[len(bcns)-1],1,3,3})
			for _, node := range bcns{
				if _, ok := tcnodes[node + nrs + 2]; ok {
					ms = append(ms, []int{node, node + nrs + 2,1,4,3})
				}
			}			
			case 4:
			//a type
			xstep = math.Round(span/nrafters)
			coords = append(coords, []float64{x, y})
			for i := 0; i < 2*nrs; i++ {
				x += xstep
				coords = append(coords, []float64{x, y})
			}
			x = 0.0
			y = t.Depth + y0
			coords = append(coords, []float64{x,y})
			for i := 0; i < 2*nrs; i++ {
				if i < nrs {
					x += xstep
					y += xstep / t.Slope
				} else {
					x += xstep
					y -= xstep / t.Slope
				}
				coords = append(coords, []float64{x, y})
			}			
			for i := 1; i < 2*nrs+2; i++ {
				bcnodes[i] = true
				bcns = append(bcns, i)
			}
			for i := 2*nrs + 2; i <= len(coords); i++ {
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
			midnode := nrs + 1
			fmt.Println("midnode",midnode)
			var ns []int
			for _, jb := range bcns{
				switch {
				case jb < midnode:
					ns = []int{jb + 2 * nrs, jb + 2*nrs + 1}
				case jb > midnode:
					ns = []int{jb + 2 * nrs+1, jb + 2 * nrs + 2}
				case jb == midnode:
					ns = []int{jb + 2 * nrs, jb + 2* nrs + 1, jb + 2*nrs + 2}
				}
				for _, je := range ns{
					if _, ok := tcnodes[je]; ok{
						if coords[jb-1][0] == coords[je-1][0]{
							ms = append(ms, []int{jb, je, 1, 3, 3})
						} else {
							ms = append(ms, []int{jb, je, 1, 4, 3})
						}
					}
				}
			}
		}
	case 2:
		//howe fan CHANGE NAMES these are quite random
		switch t.Typ {
		case 3:
			//l
			//fmt.Println("n rafters",nrafters,"rafter len", rftrl)
			xns := nrs/2
			if nrs%2 != 0 {xns++}
			xstep = math.Round(t.Span/float64(xns))
			coords = append(coords, []float64{x, y})
			for i := 0; i < xns; i++ {
				if nrs%2 != 0 && (i == xns - 1){
					x += xstep/2.0
				} else {
					x += xstep
				}
				coords = append(coords, []float64{x, y})
			}
			x = 0.0
			y = t.Depth + y0
			//y= 0.0
			//sine := math.Sin(math.Atan(1.0 / t.Slope))
			//cosine := math.Cos(math.Atan(1.0 / t.Slope))
			coords = append(coords, []float64{x, y})
			for i := 0; i < nrs; i++ {
				x += xstep/2.0
				y += xstep/2.0/t.Slope
				coords = append(coords, []float64{x, y})
			}
			for i := 1; i <= xns + 1; i++ {
				bcnodes[i] = true
				bcns = append(bcns, i)
			}
			for i := xns + 2; i <= len(coords); i++ {
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
			for i, jb := range bcns{
				for j, je := range []int{jb + xns + i,jb + xns + i+1,jb + xns + i+2} {
					if _, ok := tcnodes[je]; ok{
						if i == len(bcns) - 1{
							ms = append(ms, []int{jb, je, 1, 3, 3})
							break
						}
						if j == 1{
							ms = append(ms, []int{jb, je, 1, 3, 3})
						} else {
							ms = append(ms, []int{jb, je, 1, 4, 3})
						}
					}
				}
			}
		case 4:
			//a			
			xstep = math.Round(t.Span/ (nrafters))
			coords = append(coords, []float64{x, y})
			for i := 0; i < nrs; i++ {
				x += xstep
				coords = append(coords, []float64{x, y})
			}
			x = 0.0
			//y = 0.0
			y = t.Depth + y0
			sine := math.Sin(math.Atan(1.0 / t.Slope))
			cosine := math.Cos(math.Atan(1.0 / t.Slope))
			coords = append(coords, []float64{x, y})
			for i := 0; i < 2*nrs; i++ {
				if i < nrs {
					x += prlnspc * cosine
					y += prlnspc * sine
				} else {
					x += prlnspc * cosine
					y -= prlnspc * sine
				}
				coords = append(coords, []float64{x, y})
			}
			for i := 1; i < nrs+2; i++ {
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
			//ms = append(ms, []int{tcns[0], bcns[0], 1, 3, 3})
			//ms = append(ms, []int{tcns[len(tcns)-1], bcns[len(bcns)-1], 1, 3, 3})
			for i, node := range bcns{
				for j, je := range []int{node+nrs+i, node+ nrs + i + 1, node+nrs+i+2}{
					if _, ok := tcnodes[je]; ok{
						switch{
							case j == 0,j == 2:
							ms = append(ms, []int{node, je,1,4,3})
							case j == 1:
							ms = append(ms, []int{node, je,1,3,3})
						}
					}
				}
			}

		}
	case 3:
		//"pratt"
		xstep = math.Round(span/nrafters)
		coords = append(coords, []float64{x, y})
		
		for i := 0; i < tmul*nrs; i++ {
			x += xstep
			coords = append(coords, []float64{x, y})
		}
		x = 0.0
		y = t.Depth + y0
		coords = append(coords, []float64{x, y})
		fmt.Println("nraftarrs",nrs)
		switch t.Typ{
			case 3:
			//el type
			for i := 0; i < nrs; i++ {
				x += xstep
				y += xstep / t.Slope
				coords = append(coords, []float64{x, y})
			}
			for i := 1; i < nrs + 2; i++{
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
			for i, jb := range bcns{
				ms = append(ms, []int{jb, tcns[i],1,3,3})
			}
			//ms = append(ms, []int{tcns[len(tcns)-1], bcns[len(bcns)-1],1,3,3})
			for _, jb := range bcns[1:]{
				if _, ok := tcnodes[jb + nrs]; ok {
					ms = append(ms, []int{jb, jb + nrs,1,4,3})
				}
			}			
			case 4:
			//a type		
			for i := 0; i < 2*nrs; i++ {
				if i < nrs {
					x += xstep
					y += xstep / t.Slope
				} else {
					x += xstep
					y -= xstep / t.Slope
				}
				coords = append(coords, []float64{x, y})
			}
			for i := 1; i < 2*nrs+2; i++ {
				bcnodes[i] = true
				bcns = append(bcns, i)
			}
			for i := 2*nrs + 2; i <= len(coords); i++ {
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
			
			//ms = append(ms, []int{tcns[0], bcns[0],1,3,3})
			//ms = append(ms, []int{tcns[len(tcns)-1], bcns[len(bcns)-1],1,3,3})
			for i, jb := range bcns{
				ms = append(ms, []int{jb, tcns[i],1,3,3})
			}
			for _, jb := range bcns{
				switch {
				case jb < nrs+1:
					if _, ok := tcnodes[jb+2*nrs+2]; ok {
						ms = append(ms, []int{jb, jb + 2*nrs+2,1,4,3})
					}
				case jb == nrs+1:
					//fmt.Println("xxx666xxx666")
					continue
				default:
					if _, ok := tcnodes[jb+2*nrs]; ok {
						ms = append(ms, []int{jb, jb + 2*nrs,1,4,3})
					}
				}
			}
		}
	case 4:
		//"pratt fan"
		switch t.Typ {
		case 3:
			xns := nrs/2
			if nrs%2 != 0 {xns++}
			xstep = math.Round(t.Span/float64(xns))
			coords = append(coords, []float64{x, y})
			for i := 0; i < xns; i++ {
				if nrs%2 != 0 && (i == xns - 1){
					x += xstep/2.0
				} else {
					x += xstep
				}
				coords = append(coords, []float64{x, y})
			}
			x = 0.0
			y = t.Depth + y0
			coords = append(coords, []float64{x, y})
			for i := 0; i < nrs; i++ {
				x += xstep/2.0
				y += xstep/2.0/t.Slope
				coords = append(coords, []float64{x, y})
			}
			for i := 1; i <= xns + 1; i++ {
				bcnodes[i] = true
				bcns = append(bcns, i)
			}
			for i := xns + 2; i <= len(coords); i++ {
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
	
			for i, jb := range bcns{
				for _, je := range []int{jb + xns + i + 1,jb + xns + i,jb + xns + i-1} {
					if _, ok := tcnodes[je]; ok{
						if coords[jb-1][0] == coords[je-1][0]{
							ms = append(ms, []int{jb, je, 1, 3, 3})
						} else {
							ms = append(ms, []int{jb, je, 1, 4, 3})
						}
					}
				}
			}
		case 4:
			xstep = math.Round(t.Span * 2/ (nrafters))
			coords = append(coords, []float64{x, y})
			for i := 0; i < nrs; i++ {
				x += xstep
				coords = append(coords, []float64{x, y})
			}
			x = 0.0
			y = t.Depth + y0
			coords = append(coords, []float64{x, y})
			fmt.Println("nraftarrs",nrs)
		
			for i := 0; i < 2*nrs; i++ {
				if i < nrs {
					x += xstep/2.0
					y += xstep / t.Slope
				} else {
					x += xstep/2.0
					y -= xstep / t.Slope
				}
				coords = append(coords, []float64{x, y})
			}
			for i := 1; i < nrs + 2; i++ {
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
				
			midnode := -1
			if nrs % 2 == 0{
				midnode = nrs/2 + 1
			}
			for i, jb := range bcns{
				var jbns []int
				switch {
				case jb <= nrs/2+1:
					jbns = []int{jb + nrs + i-1, jb + nrs+i,jb + nrs + i+1}
				default:
					jbns = []int{jb + nrs + i+1, jb + nrs+i+2,jb + nrs + i+3}
				}
				if jb == midnode{
					jbns = append(jbns, []int{jb + nrs + i + 2, jb + nrs + i+3}...)
				}
				for _, je := range jbns {
					if _, ok := tcnodes[je]; ok{
						if coords[jb-1][0] == coords[je-1][0]{
							ms = append(ms, []int{jb, je, 1, 3, 3})
						} else {
							ms = append(ms, []int{jb, je, 1, 4, 3})
						}
					}
				}	
			}
		}
	case 5:
		//"fink fan (sub 14.1)"
		ngs = 5
		switch t.Typ{
			case 3:
			//el l typ
			xns := nrs/2
			//fmt.Println("nrs init--->",nrs)
			if nrs % 2 != 0{xns++}
			sine := math.Sin(math.Atan(1.0 / t.Slope))
			cosine := math.Cos(math.Atan(1.0 / t.Slope))
			xstep = math.Round(prlnspc * cosine)
			ystep = math.Round(prlnspc * sine)
			coords = append(coords, []float64{x, y})
			for i := 0; i < xns; i++{
				if nrs % 2 != 0 && i == xns - 1{
					x += xstep
				} else {
					x += 2. * xstep
				}
				coords = append(coords, []float64{x, y})
			}
			x = 0.0
			y = t.Depth + y0
			coords = append(coords, []float64{x, y})
			for i := 0; i < nrs; i++{
				x += xstep; y += ystep
				coords = append(coords, []float64{x,y})
			}
			for i := 1; i <= xns+1; i++ {
				bcnodes[i] = true
				bcns = append(bcns, i)
			}
			for i := xns + 2; i <= len(coords); i++ {
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
			//fmt.Println("nodes->",bcns, tcns)
			ms = append(ms, []int{bcns[0], tcns[0],1,3,3})
			ms = append(ms, []int{bcns[len(bcns)-1],tcns[len(tcns)-1],1,3,3})
			ms = append(ms, []int{bcns[0], tcns[1],1,4,3})
			for i, jb := range bcns[1:len(bcns)-1]{
				if nrs % 2 != 0 && jb == xns+1{continue}
				for je := i + jb + xns + 1; je < i + jb + xns + 4; je++{
					if _, ok := tcnodes[je]; ok{
						if coords[jb-1][0] == coords[je-1][0]{
							//vertical web members
							ms = append(ms, []int{jb, je,1,3,3})
						} else {
							ms = append(ms, []int{jb, je,1,4,3})
						}
					}
				}
				//if i == 0 {continue}
				midx := (coords[i+jb+xns+3-1][0] + coords[jb-1][0])/2.0
				midy := (coords[i+jb+xns+3-1][1] + coords[jb-1][1])/2.0
				coords = append(coords, []float64{midx,midy})
				//ms = append(ms, []int{i+jb+xns,len(coords)})
				ms = append(ms, []int{i+jb+xns+2,len(coords),1,5,3})
			}
			case 4:
			//a type 
			x = 0.0; y = 0.0
			sine := math.Sin(math.Atan(1.0 / t.Slope))
			cosine := math.Cos(math.Atan(1.0 / t.Slope))
			xstep = math.Round(prlnspc * cosine)
			ystep = math.Round(prlnspc * sine)
			coords = append(coords, []float64{x, y})
			
			if nrs %2 == 0 && nrs > 2 {
				
				//fmt.Println("even n rafterz")
				for i := 1; i < 2*nrs + 1; i++ {
					x += xstep
					if i == nrs || i == nrs + 1 {
						continue
					}
					if i % 2 == 0 {coords = append(coords, []float64{x,y})}
				}
				x = 0.0
				coords = append(coords, []float64{x,y+t.Depth})
				for i := 1; i < 2*nrs + 1; i++ {
					if i <= nrs {
						x += xstep
						y += ystep
					} else {
						x += xstep
						y -= ystep
					}
					coords = append(coords, []float64{x,y+t.Depth})
				}
				
				for i := 1; i < nrs+1; i++ {
					bcnodes[i] = true
					bcns = append(bcns, i)
				}
				for i := nrs + 1; i < 3*nrs + 2; i++ {
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
				//ms = append(ms, []int{bcns[0], tcns[0]})
				//ms = append(ms, []int{bcns[len(bcns)-1],tcns[len(tcns)-1]})
				for i, node := range bcns{
					switch {
					case node < nrs/2:
						for je := i + node + nrs -1; je < i +node + nrs + 2; je++ {
							if _, ok := tcnodes[je]; ok{
								if coords[node-1][0] == coords[je-1][0]{
									ms = append(ms, []int{node, je, 1, 3, 3})
								} else {
									ms = append(ms, []int{node, je, 1, 4, 3})
								}
							}
						}
					case node == nrs/2:
						for je := i + node + nrs -1; je < i +node + nrs + 1; je++ {
							if _, ok := tcnodes[je]; ok{
								if coords[node-1][0] == coords[je-1][0]{
									ms = append(ms, []int{node, je, 1, 3, 3})
								} else {
									ms = append(ms, []int{node, je, 1, 4, 3})
								}
							}
						}
						//ms = append(ms, []int{node, i+node+nrs+2, 1, 1, 3})
						midx := (coords[i+node+nrs+2-1][0] + coords[node-1][0])/2.0
						midy := (coords[i+node+nrs+2-1][1] + coords[node-1][1])/2.0
						coords = append(coords, []float64{midx,midy})
						ms = append(ms, []int{i+node+nrs,len(coords), 1, 5, 3})
						ms = append(ms, []int{i+node+nrs+1,len(coords), 1, 5, 3})
						//inclined web member
						ms = append(ms, []int{i+node+nrs+2,len(coords), 1, 4, 3})
						ms = append(ms, []int{node,len(coords), 1, 4, 3})
					case node == (nrs/2) + 1:
						for _, je := range []int{i + node + nrs + 2, i + node + nrs + 3} {
							if _, ok := tcnodes[je]; ok{								
								if coords[node-1][0] == coords[je-1][0]{
									ms = append(ms, []int{node, je, 1, 3, 3})
								} else {
									ms = append(ms, []int{node, je, 1, 4, 3})
								}
							}
						}
						//ms = append(ms, []int{node, i+node+nrs+1, 1, 1, 3})
						midx := (coords[i+node+nrs-1][0] + coords[node-1][0])/2.0
						midy := (coords[i+node+nrs-1][1] + coords[node-1][1])/2.0
						coords = append(coords, []float64{midx,midy})
						ms = append(ms, []int{i+node+nrs+2,len(coords), 1, 5, 3})
						ms = append(ms, []int{i+node+nrs+1,len(coords), 1, 5, 3})
						//sling
						ms = append(ms, []int{i+node+nrs,len(coords), 1, 4, 3})				
						ms = append(ms, []int{node,len(coords), 1, 4, 3})				
					default:
						for je := i + node + nrs +1; je < i + node + nrs + 4; je++ {
							if _, ok := tcnodes[je]; ok{
								if coords[node-1][0] == coords[je-1][0]{
									ms = append(ms, []int{node, je, 1, 3, 3})
								} else {
									ms = append(ms, []int{node, je, 1, 4, 3})
								}
							}
						}
					}
				}
			} else if nrs %2 != 0 && nrs > 2 {
				//fmt.Println("odd n rafterz")
				for i := 1; i < 2*nrs + 1; i++ {
					x += xstep
					if i == nrs {
						continue
					}
					if i % 2 == 0 {coords = append(coords, []float64{x,y})}
				}
				x = 0.0
				coords = append(coords, []float64{x, y+t.Depth})
				for i := 1; i < 2*nrs + 1; i++ {
					if i <= nrs {
						x += xstep
						y += ystep
					} else {
						x += xstep
						y -= ystep
					}
					coords = append(coords, []float64{x,y+t.Depth})
				}
				for i := 1; i < nrs+2; i++ {
					bcnodes[i] = true
					bcns = append(bcns, i)
				}
				for i := nrs + 2; i < 3*nrs + 3; i++ {
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
				//ms = append(ms, []int{bcns[0],tcns[0]})
				//ms = append(ms, []int{bcns[len(bcns)-1],tcns[len(tcns)-1]})
				for i, node := range bcns {
					switch {
					case node <= nrs/2:
						for je := i + node + nrs; je <= i +node + nrs + 2; je++ {
							if _, ok := tcnodes[je]; ok{
								if coords[node-1][0] == coords[je-1][0]{
									ms = append(ms, []int{node, je, 1, 3, 3})
								} else {
									ms = append(ms, []int{node, je, 1, 4, 3})
								}
							}
						}
					case node == (nrs/2) + 1:
						for _, je := range []int{i + node + nrs,  i + node + nrs + 1} {
							if _, ok := tcnodes[je]; ok{
								if coords[node-1][0] == coords[je-1][0]{
									ms = append(ms, []int{node, je, 1, 3, 3})
								} else {
									ms = append(ms, []int{node, je, 1, 4, 3})
								}
							}
						}
						midx := (coords[i+node+nrs+2-1][0] + coords[node-1][0])/2.0
						midy := (coords[i+node+nrs+2-1][1] + coords[node-1][1])/2.0
						coords = append(coords, []float64{midx,midy})
						ms = append(ms, []int{node,len(coords), 1, 4, 3})
						ms = append(ms, []int{i+node+nrs+2,len(coords), 1, 4, 3})
						ms = append(ms, []int{i+node+nrs+1,len(coords), 1, 5, 3})
					case node == (nrs/2) + 2:
						for _, je := range []int{i + node + nrs+1, i + node + nrs + 2} {
							if _, ok := tcnodes[je]; ok{
								if coords[node-1][0] == coords[je-1][0]{
									ms = append(ms, []int{node, je, 1, 3, 3})
								} else {
									ms = append(ms, []int{node, je, 1, 4, 3})
								}
							}
						}
						midx := (coords[i+node+nrs-1][0] + coords[node-1][0])/2.0
						midy := (coords[i+node+nrs-1][1] + coords[node-1][1])/2.0
						coords = append(coords, []float64{midx,midy})
						ms = append(ms, []int{node,len(coords), 1, 4, 3})
						ms = append(ms, []int{i+node+nrs,len(coords), 1, 4, 3})
						ms = append(ms, []int{i+node+nrs+1,len(coords), 1, 5, 3})
					default:
						for je := i + node + nrs ; je < i + node + nrs + 3; je++ {
							if _, ok := tcnodes[je]; ok{
								if coords[node-1][0] == coords[je-1][0]{
									ms = append(ms, []int{node, je, 1, 3, 3})
								} else {
									ms = append(ms, []int{node, je, 1, 4, 3})
								}

							}
						}
					}
				}
			}
		}
	case 6://rendom pratt div
		//IT IS NOT WERK
		switch t.Typ {
		case 2:
			xstep = math.Round(t.Span / (nrafters))
			coords = append(coords, []float64{x, y})
			for i := 0; i < nrs; i++ {
				x += xstep
				coords = append(coords, []float64{x, y})
			}
			x = 0.0
			y = 0.0
			//sine := math.Sin(math.Atan(1/t.Slope))
			//cosine := math.Cos(math.Atan(1/t.Slope))
			for i := 0; i < 2*nrs-1; i++ {
				if i < nrs {
					x += xstep / 2.0
					y += (xstep / 2.0) / t.Slope
				} else {
					x += xstep / 2.0
					y -= (xstep / 2.0) / t.Slope
				}
				coords = append(coords, []float64{x, y})
			}
			for i := 1; i < nrs+2; i++ {
				bcnodes[i] = true
				bcns = append(bcns, i)
			}
			for i := nrs + 2; i <= len(coords); i++ {
				tcns = append(tcns, i)
				tcnodes[i] = true
			}
			for _, jb := range bcns {
				if _, ok := bcnodes[jb+1]; ok {
					ms = append(ms, []int{jb, jb + 1})
				}
			}
			for _, jb := range tcns {
				if _, ok := tcnodes[jb+1]; ok {
					ms = append(ms, []int{jb, jb + 1})
				}
			}
			ms = append(ms, []int{bcns[0], tcns[0]})
			ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1]})
			var midnode int
			midnode = -99
			nbcs := len(bcns[1 : len(bcns)-1])
			if nbcs%2 != 0 {
				midnode = bcns[len(bcns)/2]
			}
			for i, node := range bcns[1 : len(bcns)-1] {
				ms = append(ms, []int{node, node + nrs + i})
				ms = append(ms, []int{node, node + nrs + i + 1})
				ms = append(ms, []int{node, node + nrs + i + 2})

				if node == midnode {
					continue
				}
				if i < nbcs/2 {
					midx := (coords[node-1][0] + coords[node+nrs+i+2-1][0]) / 2.0
					midy := (coords[node-1][1] + coords[node+nrs+i+2-1][1]) / 2.0
					coords = append(coords, []float64{midx, midy})
					ms = append(ms, []int{node + nrs + i + 1, len(coords)})
				} else {
					midx := (coords[node-1][0] + coords[node+nrs+i-1][0]) / 2.0
					midy := (coords[node-1][1] + coords[node+nrs+i-1][1]) / 2.0
					coords = append(coords, []float64{midx, midy})
					ms = append(ms, []int{node + nrs + i + 1, len(coords)})

				}
			}
		}
	}
	t.Coords = coords
	t.Ms = ms
	t.Bcns = bcns
	t.Tcns = tcns
	t.Purlinspc = prlnspc
	t.Rftrl = rftrl
	t.Rise = rise
	t.Nrs = nrs
	t.Ngs = ngs
	t.Mod.Coords = coords; t.Mod.Mprp = ms
	t.Mod.Supports = [][]int{{bcns[0],-1,-1},{bcns[len(bcns)-1],-1,-1}}
	return
}

