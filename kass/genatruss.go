package barf

import (
	"log"
	"fmt"
	"math"
	//"errors"
)

//GenATruss generates a type truss models
//there are five(?) types with different sling orientation (mostly that) 
func GenATruss(t *Trs2d) (err error){
	//fmt.Println("USE FREECAD reg idiota")
	var coords [][]float64
	var ms [][]int
	var bcns, tcns []int
	var prlnspc, rftrl float64
	var nrs, ngs int
	
	bcnodes := make(map[int]bool)
	tcnodes := make(map[int]bool)
	var nrafters, x, y, y0, xstep, ystep float64
	prlnspc = t.Purlinspc
	span := t.Span/float64(t.Typ)
	rise := t.Rise
	if t.Slope == 0.0{
		if rise == 0.0{
			if t.Tang == 0.0{
				err = fmt.Errorf("missing specs - slope %f rise %f angle %f",t.Slope, t.Rise, t.Tang)
				return
			}
			t.Rise = span * math.Tan(t.Tang*math.Pi/180.0)
			rise = t.Rise
		}
		t.Slope = span/t.Rise
	}
	if rise == 0.0 {rise = span/t.Slope}
	rftrl = math.Round(math.Sqrt(math.Pow(span, 2) + math.Pow(rise, 2)))
	
	if t.Spam{
		log.Println("rafter len, span, rise->",rftrl, span, rise)
		
	}
	if rftrl > prlnspc {
		nrafters = math.Round(rftrl / prlnspc)
		prlnspc = math.Round(rftrl / nrafters)
	} else {
		prlnspc = rftrl 
		nrafters = 1
		//err = errors.New("single rafter,use kingpost/etc funcs")
		//return
	}
	nrs = int(nrafters)
	if t.Spam{log.Println("nrafters, purlinspc->",nrafters, prlnspc)}
	if t.Bent{y0 = t.Height}
	y = y0
	ngs = 4
	if nrs <= 2{
		t.Cfg = 1
	}
	switch t.Cfg{
		case -1:
		//ground structure
		switch t.Typ{
			case 1:
			case 2:
		}
		case 1:
		//howe
		xstep = math.Round(span/nrafters)
		coords = append(coords, []float64{x, y})
		for i := 0; i < t.Typ*nrs; i++{
			x += xstep
			coords = append(coords, []float64{x, y})
		}
		x = 0.0
		y = y0
		switch t.Typ{
			case 1:		
			//el typo
			ngs = 4
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
				if _, ok := bcnodes[jb+1]; ok{
					//bottom chord
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
				if _, ok := tcnodes[jb+1]; ok{
					//top chord
					switch t.Tcr{
						case 0:
						ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
						case 1:
						//bending 
						ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
					}
				}
			}
			switch t.Tcr{
				case 0:
				ms = append(ms, []int{tcns[0], bcns[0], 1, 2, 3})
				case 1:
				//bending for top chord
				ms = append(ms, []int{tcns[0], bcns[0], 1, 2, 0})
			}
			

			for i, node := range bcns[1 : len(bcns)-1] {
				ms = append(ms, []int{node, tcns[i], 1, 3, 3})
			}
			ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1], 1, 3, 3})
			for _, node := range bcns[1 : len(bcns)-1] {
				if _, ok := tcnodes[node + nrs + 1]; ok {
					ms = append(ms, []int{node, node + nrs + 1, 1, 4, 3})
				}
			}
			case 2 :
			//a type
			ngs = 4
			for i := 0; i < 2*nrs-1; i++ {
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
				if _, ok := bcnodes[jb+1]; ok{
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
				if _, ok := tcnodes[jb+1]; ok{
					//top chord
					switch t.Tcr{
						case 0:
						ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
						case 1:
						//bending for top chord
						ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
					}
				}
			}
			switch t.Tcr{
				case 0:
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 3})
				ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1], 1, 2, 3})
				case 1:
				ms = append(ms, []int{tcns[0], bcns[0], 1, 2, 0})
				ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1], 1, 2, 0})
			}
			for i, node := range bcns[1 : len(bcns)-1] {
				ms = append(ms, []int{node, tcns[i], 1, 3, 3})
			}
			for _, node := range bcns[1 : len(bcns)-1] {
				switch {
				case node < nrs+1:
					if _, ok := tcnodes[node+2*nrs-1]; ok {
						ms = append(ms, []int{node, node + 2*nrs - 1, 1, 4, 3})
					}
				case node == nrs+1:
					if _, ok := tcnodes[node+2*nrs+1]; ok {
						ms = append(ms, []int{node, node + 1 + 2*nrs, 1, 4, 3})
					}
					if _, ok := tcnodes[node+2*nrs-1]; ok {
						ms = append(ms, []int{node, node + 2*nrs - 1, 1, 4, 3})
					}
				default:
					if _, ok := tcnodes[node+2*nrs+1]; ok {
						ms = append(ms, []int{node, node + 2*nrs + 1, 1, 4, 3})
					}
				}
			}
		}
	case 2:
		//howe fan
		ngs = 4
		switch t.Typ {
		case 1:
			//l
			//fmt.Println("n rafters",nrafters,"rafter len", rftrl)
			xns := nrs/2
			if nrs%2 != 0 {xns++}
			xstep = math.Round(t.Span/float64(xns))
			coords = append(coords, []float64{x, y})
			for i := 0; i < xns; i++ {
				if i == xns -1 && nrs%2 != 0 {
					x += xstep/2.0
				} else {
					x += xstep
				}
				coords = append(coords, []float64{x, y})
			}
			x = 0.0
			y = y0
			//sine := math.Sin(math.Atan(1.0 / t.Slope))
			//cosine := math.Cos(math.Atan(1.0 / t.Slope))
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
					//bottom chord
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
					//top chord
					switch t.Tcr{
						case 0:
						ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
						case 1:
						//.bending for top chord
						ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
					}
				}
			}
			switch t.Tcr{
				case 0:
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 3})
				//ms = append(ms, []int{bcns[len(bcns)-1],tcns[len(tcns)-1],1,2,3})
				case 1:
				//bending for top chord
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 0})
				//ms = append(ms, []int{bcns[len(bcns)-1],tcns[len(tcns)-1],1,2,0})
			}
			for i, jb := range bcns[1:]{
				for j, je := range []int{jb + xns + i,jb + xns + i+1,jb + xns + i+2} {
					if _, ok := tcnodes[je]; ok{
						switch j{
							case 0:
							ms = append(ms, []int{jb, je,1,4,3})
							case 1:
							ms = append(ms, []int{jb, je,1,3,3})
							case 2:
							ms = append(ms, []int{jb, je,1,4,3})
						}
					}
				}
			}
			
		case 2:
			//a			
			xstep = math.Round(t.Span/ (nrafters))
			coords = append(coords, []float64{x, y})
			for i := 0; i < nrs; i++ {
				x += xstep
				coords = append(coords, []float64{x, y})
			}
			x = 0.0
			y = y0
			sine := math.Sin(math.Atan(1.0 / t.Slope))
			cosine := math.Cos(math.Atan(1.0 / t.Slope))
			for i := 0; i < 2*nrs-1; i++ {
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
					//bottom chord
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
					//top chord
					switch t.Tcr{
						case 0:
						ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
						case 1:
						//bending 
						ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
					}
				}
			}
			switch t.Tcr{
				case 0:
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 3})
				ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1],1, 2, 3})
				case 1:
				//bending for top chord
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 0})
				ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1],1, 2, 0})
			}
			for i, node := range bcns[1 : len(bcns)-1] {
				ms = append(ms, []int{node, node + nrs + i, 1, 4, 3})
				ms = append(ms, []int{node, node + nrs + i + 1, 1, 3, 3})
				if i == len(bcns)-1{
					ms = append(ms, []int{node, node + nrs + i, 1, 4, 3})
				} else {ms = append(ms, []int{node, node + nrs + i + 2, 1, 4, 3})}
				
			}
		}
	case 3:
		//"pratt"
		ngs = 4
		xstep = math.Round(span/nrafters)
		coords = append(coords, []float64{x, y})
		for i := 0; i < t.Typ*nrs; i++ {
			x += xstep
			coords = append(coords, []float64{x, y})
		}
		x = 0.0
		y = y0
		switch t.Typ{
			case 1:
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
					//bottom chord
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
					//top chord
					switch t.Tcr{
						case 0:
						ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
						case 1:
						//.bending for top chord
						ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
					}
				}
			}
			switch t.Tcr{
				case 0:
				ms = append(ms, []int{tcns[0], bcns[0], 1, 2, 3})
				case 1:
				ms = append(ms, []int{tcns[0], bcns[0], 1, 2, 0})
			}
			for i, node := range bcns[1 : len(bcns)-1] {
				ms = append(ms, []int{node, tcns[i], 1, 3, 3})
			}
			ms = append(ms, []int{tcns[len(tcns)-1], bcns[len(bcns)-1], 1, 3, 3})
			for _, node := range bcns[2:]{
				if _, ok := tcnodes[node + nrs - 1]; ok {
					ms = append(ms, []int{node, node + nrs - 1, 1, 4, 3})
				}
			}			
			case 2:
			//a type		
			for i := 0; i < 2*nrs-1; i++ {
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
					//bottom chord
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
					//top chord
					switch t.Tcr{
						case 0:
						ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
						case 1:
						//bending 
						ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
					}
				}
			}
			
			switch t.Tcr{
				case 0:
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 3})
				ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1],1, 2, 3})
				case 1:
				//bending for top chord
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 0})
				ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1],1, 2, 0})
			}
			
			for i, node := range bcns[1 : len(bcns)-1] {
				ms = append(ms, []int{node, tcns[i],1, 3, 3})
			}
			for _, node := range bcns[1 : len(bcns)-1] {
				switch {
				case node < nrs+1:
					if _, ok := tcnodes[node+2*nrs+1]; ok {
						ms = append(ms, []int{node, node + 2*nrs + 1, 1, 4, 3})
					}
				case node == nrs+1:
					continue
				default:
					if _, ok := tcnodes[node+2*nrs-1]; ok {
						ms = append(ms, []int{node, node + 2*nrs - 1, 1, 4, 3})
					}
				}
			}
		}
	case 4:
		//"pratt fan"
		switch t.Typ {
		case 1:
			//l
			//fmt.Println("n rafters",nrafters,"rafter len", rftrl)
			xns := nrs/2
			if nrs % 2 !=0 {xns += 1}
			xstep = math.Round(t.Span/float64(xns))
			coords = append(coords, []float64{x, y})
			for i := 0; i < xns; i++ {
				if nrs %2 != 0 && i == xns - 1{
					x += xstep/2.0
					coords = append(coords, []float64{x, y})
				} else {
					x += xstep
					coords = append(coords, []float64{x, y})
				}
			}
			x = 0.; y = y0
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
					//bottom chord
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
					//top chord
					switch t.Tcr{
						case 0:
						ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
						case 1:
						//bending 
						ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
					}
				}
			}
			switch t.Tcr{
				case 0:
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 3})

				case 1:
				//bending for top chord
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 0})

			}
			for i, jb := range bcns[1:]{
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
			//fmt.Println("nraftars-",nrs,"raftarrs")
			/*
			if nrs % 2 != 0{
				jb := bcns[len(bcns)-1]
				je := jb + nrs - 1
				ms = append(ms, []int{jb, je, 1, 4, 3})
			}
			*/
		case 2:
			//fmt.Println("n rafters",nrafters,"rafter len", rftrl)
			xstep = math.Round(t.Span / (nrafters))
			coords = append(coords, []float64{x, y})
			for i := 0; i < nrs; i++ {
				x += xstep
				coords = append(coords, []float64{x, y})
			}
			x = 0.0
			y = y0
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
					//bottom chord
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
					//top chord
					switch t.Tcr{
						case 0:
						ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
						case 1:
						//bending 
						ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
					}
				}
			}
			switch t.Tcr{
				case 0:
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 3})
				ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1],1, 2, 3})
				case 1:
				//bending for top chord
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 0})
				ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1],1, 2, 0})
			}
			//var midnode int
			//midnode = -99
			//nbcs := len(bcns[1 : len(bcns)-1])
			//if nbcs%2 != 0 {
			//	midnode = bcns[len(bcns)/2]
			//}
			switch{
				case nrs % 2 == 0:
				//connect midnode
				// fmt.Println("nrs%2 == 0, case 1")
				for i, jb := range bcns[1 : len(bcns)-1] {
					var jbns []int
					switch {
					case jb <= nrs/2:
						jbns = []int{jb + nrs + i + 1,jb + nrs + i,jb + nrs + i-1}
					case jb == nrs/2 + 1:
						jbns = []int{jb + nrs + i + 1,jb + nrs + i,jb + nrs + i-1}
						jbns = append(jbns, []int{jb + nrs + i + 1,jb + nrs + i + 2,jb + nrs + i+3}...)
					default:
						jbns = []int{jb + nrs + i + 1,jb + nrs + i + 2,jb + nrs + i+3}
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
				default:
				//no midnode no cry
				//fmt.Println("weird one")
				for i, jb := range bcns[1 : len(bcns)-1] {
					var jbns []int
					switch {
					case jb <= nrs/2 + 1:
						jbns = []int{jb + nrs + i + 1,jb + nrs + i,jb + nrs + i-1}
						if jb == nrs/2 + 1{jbns = append(jbns, 2 * jb + nrs)}
					case jb == nrs/2 + 1:
						continue
					default:
						jbns = []int{jb + nrs + i + 1,jb + nrs + i + 2,jb + nrs + i+3}
						if jb == nrs/2 + 2 {
							jbns = append(jbns, 2*jb + nrs - 2)
						}
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
		}
		case 5://"fink fan 1", bhavikatti 5.3
		ngs = 5
		switch t.Typ{
			case 1:
			nbm := nrs/2
			if nrs %2 != 0{nbm = (nrs+1)/2}
			ang := math.Atan(t.Rise/t.Span)
			sine := math.Sin(ang)
			cosine := math.Cos(ang)
			coords = append(coords, []float64{x, y})
			xstep = t.Span/float64(nbm)
			for i := 0; i < nbm; i++{
				x += xstep
				coords = append(coords, []float64{x, y})
			}
			x = 0.0; y = y0
			xstep := math.Round(prlnspc * cosine)
			ystep := math.Round(prlnspc * sine)
			for i := 0; i < nrs; i++{
				x += xstep
				y += ystep
				coords = append(coords, []float64{x, y})
			}
			for i := 1; i <= nbm+1; i++ {
				bcnodes[i] = true
				bcns = append(bcns, i)
			}
			for i := nbm + 2; i <= len(coords); i++ {
				tcns = append(tcns, i)
				tcnodes[i] = true
			}
			for _, jb := range bcns {
				if _, ok := bcnodes[jb+1]; ok {					
					//bottom chord
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
					//top chord
					switch t.Tcr{
						case 0:
						ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
						case 1:
						//bending 
						ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
					}
				}
			}
			switch t.Tcr{
				case 0:
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 3})
				case 1:
				//bending for top chord
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 0})

			}
			ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1],1, 4, 3})
			bcwn := bcns[1:len(bcns)-1]

			if len(bcwn) == 1{
				jb := bcwn[0]
				switch nrs{
					case 3:
					ms = append(ms, []int{jb, jb+nbm, 1, 3, 3})
					//get midnode
					je := jb + nbm + 3
					x1 := (coords[jb-1][0] + coords[je-1][0])/2.0
					y1 := (coords[jb-1][1] + coords[je-1][1])/2.0
					coords = append(coords, []float64{x1, y1})
					je = len(coords)
					ms = append(ms, []int{jb, je, 1, 5, 3})
					ms = append(ms, []int{je, jb+nbm+2, 1, 5, 3})
					ms = append(ms, []int{jb+nbm+1,je, 1, 5, 3})
					case 4:		
					ms = append(ms, []int{jb, jb+nbm, 1, 3, 3})
					ms = append(ms, []int{jb, jb+nbm+1, 1, 4, 3})
					//get midnode
					je := jb + nbm + 3
					x1 := (coords[jb-1][0] + coords[je-1][0])/2.0
					y1 := (coords[jb-1][1] + coords[je-1][1])/2.0
					coords = append(coords, []float64{x1, y1})
					je = len(coords)
					ms = append(ms, []int{jb, je, 1, 5, 3})
					ms = append(ms, []int{je, jb+nbm+3, 1, 5, 3})
					ms = append(ms, []int{jb+nbm+2, je, 1, 5, 3})
					ms = append(ms, []int{jb+nbm+1,je, 1, 5, 3})
			
				}
				
			} else {			
				for i, jb := range bcwn{
					switch i{								
						case 0:
						if nrs % 2 != 0{
							for j, je := range []int{jb+nbm, jb+nbm+1}{
								switch j{
									case 0:
									ms = append(ms, []int{jb, je, 1, 3, 3})
									case 1:
									ms = append(ms, []int{jb, je, 1, 5, 3})
								}
							}
						} else {
							for j, je := range []int{jb+nbm, jb+nbm+1, jb+nbm+2}{
								switch j{
									case 0:
									ms = append(ms, []int{jb, je, 1, 3, 3})
									case 1:
									ms = append(ms, []int{jb, je, 1, 4, 3})
									case 2:
									ms = append(ms, []int{jb, je, 1, 5, 3})
								}
							}
						}
						case len(bcwn)-1:
						di := nbm + i
						if nrs % 2 != 0 && nbm > 2{
							di -= 1
						}
						je := jb + di + 3
						//get midnode
						x1 := (coords[jb-1][0] + coords[je-1][0])/2.0
						y1 := (coords[jb-1][1] + coords[je-1][1])/2.0
						coords = append(coords, []float64{x1, y1})
						je = len(coords)
						ms = append(ms, []int{jb, je, 1, 5, 3})
						ms = append(ms, []int{je, jb + di + 3, 1, 5, 3})
						ms = append(ms, []int{jb + di, je, 1, 3, 3})
						ms = append(ms, []int{jb + di + 1, je, 1, 5, 3})
						ms = append(ms, []int{jb + di + 2, je, 1, 5, 3})		
						ms = append(ms, []int{jb, jb + di, 1, 4, 3})

						default:
						di := nbm + i
						if nrs % 2 != 0{
							di -= 1
						}
						for j, je := range []int{jb + di, jb + di + 1, jb + di + 2}{
							switch j{
								case 0:
								ms = append(ms, []int{jb, je, 1, 3, 3})
								case 1:
								ms = append(ms, []int{jb, je, 1, 4, 3})
								case 2:
								ms = append(ms, []int{jb, je, 1, 5, 3})
								
							}
						}
					}
				}
			}
			case 2:
			ngs = 5
			nbm := nrs
			if nrs %2 != 0{nbm = nrs +1}
			ang := math.Atan(t.Rise/span)
			sine := math.Sin(ang)
			cosine := math.Cos(ang)
			coords = append(coords, []float64{x, y})
			xstep = t.Span/float64(nbm)
			for i := 0; i < nbm; i++{
				x += xstep
				coords = append(coords, []float64{x, y})
			}
			x = 0.0; y = y0
			xstep := math.Round(prlnspc * cosine)
			ystep := math.Round(prlnspc * sine)
			for i := 0; i <= 2 *(nrs-1); i++{
				x += xstep
				if i < nrs{
					y += ystep
				} else {
					y -= ystep
				}
				coords = append(coords, []float64{x, y})
			}
			for i := 1; i <= nbm+1; i++ {
				bcnodes[i] = true
				bcns = append(bcns, i)
			}
			for i := nbm + 2; i <= len(coords); i++ {
				tcns = append(tcns, i)
				tcnodes[i] = true
			}
			for _, jb := range bcns {
				if _, ok := bcnodes[jb+1]; ok {					
					//bottom chord
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
					//top chord
					switch t.Tcr{
						case 0:
						ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
						case 1:
						//bending 
						ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
					}
				}
			}
			switch t.Tcr{
				case 0:
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 3})
				ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1],1, 2, 3})
				case 1:
				//bending for top chord
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 0})
				ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1],1, 2, 0})
			}
			bcwn := bcns[1:len(bcns)-1]
			jtop := tcns[len(tcns)/2]
			jbot := nbm/2 + 1
			jl := jbot - 1
			jr := jbot + 1
			for i, jb := range bcwn{
				switch{
					case jb == jbot:
					ms = append(ms, []int{jbot, jtop, 1, 4, 3})
					case jb == jl:
					x1 := (coords[jb-1][0] + coords[jtop-1][0])/2.0
					y1 := (coords[jb-1][1] + coords[jtop-1][1])/2.0
					coords = append(coords, []float64{x1, y1})
					je := len(coords)
					ms = append(ms, []int{jb, je, 1, 5, 3})
					ms = append(ms, []int{je, jtop, 1, 5, 3})
					if _, ok := tcnodes[jtop-3]; ok{
						ms = append(ms, []int{jb, jtop-3, 1, 3, 3})
						ms = append(ms, []int{jtop-3, je, 1, 5, 3})
					} else {
						ms = append(ms, []int{jb, jtop-2, 1, 3, 3})
					}
					ms = append(ms, []int{jtop-2, je, 1, 5, 3})
					ms = append(ms, []int{jtop-1, je, 1, 5, 3})
					case jb == jr:
					x1 := (coords[jb-1][0] + coords[jtop-1][0])/2.0
					y1 := (coords[jb-1][1] + coords[jtop-1][1])/2.0
					coords = append(coords, []float64{x1, y1})
					je := len(coords)
					ms = append(ms, []int{jb, je, 1, 5, 3})
					ms = append(ms, []int{je, jtop, 1, 5, 3})
					if _, ok := tcnodes[jtop+3]; ok{
						ms = append(ms, []int{jb, jtop+3, 1, 3, 3})
						ms = append(ms, []int{jtop+3, je, 1, 5, 3})
					} else {
						ms = append(ms, []int{jb, jtop+2, 1, 3, 3})
					}
					ms = append(ms, []int{jtop+2, je, 1, 5, 3})
					ms = append(ms, []int{jtop+1, je, 1, 5, 3})
					case jb < jl:
					if i == 0 && nrs % 2 != 0{
						ms = append(ms, []int{jb, jb+nbm, 1, 3, 3})
						ms = append(ms, []int{jb, jb+nbm+1, 1, 4, 3})
					} else {
						
						ms = append(ms, []int{jb, jb+nbm, 1, 3, 3})
						ms = append(ms, []int{jb, jb+nbm+1, 1, 4, 3})
						ms = append(ms, []int{jb, jb+nbm+2, 1, 5, 3})
					}
					case jb > jr:
					if i == len(bcwn)-1 && nrs % 2 != 0{
						ms = append(ms, []int{jb, jb+nbm+i+2,1,4,3})
						ms = append(ms, []int{jb, jb+nbm+i+1,1,3,3})
					} else {	
						ms = append(ms, []int{jb, jb+nbm+i,1,5,3})
						ms = append(ms, []int{jb, jb+nbm+i+1,1,4,3})
						ms = append(ms, []int{jb, jb+nbm+i+2,1,3,3})
					}
				}
			}
		}
		case 6:
		//even numbered fink fan (duggal)
		switch t.Typ{
			case 1:
			nbm := nrs - 1
			ang := math.Atan(t.Rise/t.Span)
			x = 0.0
			y = y0
			xstep = math.Round(t.Span/(nrafters - 1.0))
			coords = append(coords, []float64{x, y})
			for i := 0; i < nrs-1; i++{
				x += xstep
				coords = append(coords, []float64{x, y})
			}
			x = 0.0; y = y0
			
			xstep = math.Cos(ang) * prlnspc
			ystep = math.Sin(ang) * prlnspc
			for i := 0; i < nrs; i++{
				x += xstep
				y += ystep
				coords = append(coords, []float64{x,y})
			}
			for i := 1; i <= nrs; i++ {
				bcnodes[i] = true
				bcns = append(bcns, i)
			}
			for i := nrs+1; i <= len(coords); i++ {
				tcns = append(tcns, i)
				tcnodes[i] = true
			}
			for _, jb := range bcns {
				if _, ok := bcnodes[jb+1]; ok {					
					//bottom chord
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
					//top chord
					switch t.Tcr{
						case 0:
						ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
						case 1:
						//bending 
						ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
					}
				}
			}
			switch t.Tcr{
				case 0:
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 3})
				case 1:
				//bending for top chord
				ms = append(ms, []int{bcns[0], tcns[0], 1, 2, 0})
			}
			ms = append(ms, []int{bcns[len(bcns)-1], tcns[len(tcns)-1],1, 3, 3})
			for _, jb := range bcns[1:len(bcns)-1]{
				if jb == nbm{
					je := jb + nbm
					ms = append(ms, []int{jb, je, 1, 3, 3})
					//get midnode of jb, jb + nbm + 1
					c1 := coords[jb-1]; c2 := coords[je+1]
					coords = append(coords, []float64{(c1[0]+c2[0])/2.0,(c1[1]+c2[1])/2.0})
					je = len(coords)
					ms = append(ms, []int{jb, je, 1, 3, 3})
					ms = append(ms, []int{je, tcns[len(tcns)-1], 1, 3, 3})
					ms = append(ms, []int{tcns[len(tcns)-2],je,1,4,3})
					ms = append(ms, []int{tcns[len(tcns)-3],je,1,4,3})
				} else {					
					for _, je := range []int{jb+nbm, jb+nbm+1}{
						if _, ok := tcnodes[je]; ok{
							ms = append(ms, []int{jb, je, 1, 3, 3})
						}
					}
				}
			}
			case 2:
			//WRITING TWO VERY SIMILAR LOOPS FOR ONE EXTRA MEMBER IS LAME
			if t.Kpost{
				nbm := 2 * (nrs-1)
				ang := math.Atan(t.Rise/span)
				x = 0.0
				xstep = math.Round(t.Span/(2.0 *(nrafters - 1.0)))
				coords = append(coords, []float64{x, y})
				for i := 0; i < nrs + 2; i++{
					x += xstep
					coords = append(coords, []float64{x, y})
				}			
				x = 0.0
				xstep = math.Cos(ang) * prlnspc
				ystep = math.Sin(ang) * prlnspc
				for i := 0; i <= 2*(nrs-1); i++{	
					x += xstep
					if i < nrs{
						y += ystep
					} else {
						y -= ystep
					}
					coords = append(coords, []float64{x,y})
				}
				for i := 1; i <= 2*(nrs-1)+1; i++ {
					bcnodes[i] = true
					bcns = append(bcns, i)
				}
				for i := 2*nrs; i <= len(coords); i++ {
					tcns = append(tcns, i)
					tcnodes[i] = true
				}
				for _, jb := range bcns {
					if _, ok := bcnodes[jb+1]; ok {					
						//bottom chord
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
						//top chord
						switch t.Tcr{
							case 0:
							ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
							case 1:
							//bending 
							ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
						}
					}
				}
				switch t.Tcr{
					case 0:
					jb := bcns[0]
					je := tcns[0]
					ms = append(ms, []int{jb, je, 1, 2, 3})
					jb = bcns[len(bcns)-1]
					je = tcns[len(tcns)-1]
					ms = append(ms, []int{jb, je, 1, 2, 3})
					case 1:
					jb := bcns[0]
					je := tcns[0]
					ms = append(ms, []int{jb, je, 1, 2, 0})
					jb = bcns[len(bcns)-1]
					je = tcns[len(tcns)-1]
					ms = append(ms, []int{jb, je, 1, 2, 0})
				}
				for _, jb := range bcns[1:len(bcns)-1] {
					switch{
						case jb == nrs:
						//add kingpost strut
						je := jb + nbm + 1
						ms = append(ms, []int{jb, je, 1, 5, 3})
						case jb == nrs-1:
						je := jb + nbm 
						ms = append(ms, []int{jb, je, 1, 3, 3})
						c1 := coords[jb-1]; c2 := coords[jb+nbm+2-1]
						coords = append(coords, []float64{(c1[0]+c2[0])/2.0,(c1[1]+c2[1])/2.0})
						je = len(coords)
						ms = append(ms, []int{jb, je, 1, 3, 3})
						ms = append(ms, []int{je, jb+nbm+2, 1, 3, 3})
						ms = append(ms, []int{jb+nbm+1,je,1,4,3})
						ms = append(ms, []int{jb+nbm,je,1,4,3})
						case jb == nrs+1:
						je := jb + nbm + 2				
						ms = append(ms, []int{jb, je, 1, 3, 3})
						//get midnode of jb, jb + nbm 
						c1 := coords[jb-1]; c2 := coords[jb+nbm-1]
						coords = append(coords, []float64{(c1[0]+c2[0])/2.0,(c1[1]+c2[1])/2.0})
						je = len(coords)
						ms = append(ms, []int{jb, je, 1, 3, 3})
						ms = append(ms, []int{je, jb+nbm, 1, 3, 3})
						ms = append(ms, []int{jb+nbm+1,je,1,4,3})
						ms = append(ms, []int{jb+nbm+2,je,1,4,3})
						case jb < nrs-1:
						for _, je := range []int{jb+nbm, jb+nbm+1}{
							if _, ok := tcnodes[je]; ok{
								ms = append(ms, []int{jb, je, 1, 3, 3})
							}
						}
						case jb > nrs+1:
						for _, je := range []int{jb+nbm+1, jb+nbm+2}{
							if _, ok := tcnodes[je]; ok{
								ms = append(ms, []int{jb, je, 1, 3, 3})
							}
						}
					}
				}	
			} else {
				nbm := 2 * (nrs-1) - 1
				//fmt.Println(ColorGreen,"NBM",nbm,ColorReset)
				ang := math.Atan(t.Rise/span)
				x = 0.0
				xstep = math.Round(t.Span/(2.0 *(nrafters - 1.0)))
				//fmt.Println("XSTEP-",xstep)
				coords = append(coords, []float64{x, y})
				for i := 0; i < nbm; i++{
					if i == nrs-2{
						x += 2.0 * xstep
					} else {
						x += xstep
					}
					coords = append(coords, []float64{x, y})
				}			
				x = 0.0
				xstep = math.Cos(ang) * prlnspc
				ystep = math.Sin(ang) * prlnspc
				for i := 0; i <= 2*(nrs-1); i++{	
					x += xstep
					if i < nrs{
						y += ystep
					} else {
						y -= ystep
					}
					coords = append(coords, []float64{x,y})
				}
				for i := 1; i <= 2*(nrs-1); i++ {
					bcnodes[i] = true
					bcns = append(bcns, i)
				}
				for i := 2*(nrs-1)+1; i <= len(coords); i++ {
					tcns = append(tcns, i)
					tcnodes[i] = true
				}
				for _, jb := range bcns {
					if _, ok := bcnodes[jb+1]; ok {					
						//bottom chord
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
						//top chord
						switch t.Tcr{
							case 0:
							ms = append(ms, []int{jb, jb + 1, 1, 2, 3})
							case 1:
							//bending 
							ms = append(ms, []int{jb, jb + 1, 1, 2, 0})
						}
					}
				}
				switch t.Tcr{
					case 0:
					jb := bcns[0]
					je := tcns[0]
					ms = append(ms, []int{jb, je, 1, 2, 3})
					jb = bcns[len(bcns)-1]
					je = tcns[len(tcns)-1]
					ms = append(ms, []int{jb, je, 1, 2, 3})
					case 1:
					jb := bcns[0]
					je := tcns[0]
					ms = append(ms, []int{jb, je, 1, 2, 0})
					jb = bcns[len(bcns)-1]
					je = tcns[len(tcns)-1]
					ms = append(ms, []int{jb, je, 1, 2, 0})
				}
				for _, jb := range bcns[1:len(bcns)-1] {
					switch{
						case jb == nrs:
						je := jb + nbm + 3				
						ms = append(ms, []int{jb, je, 1, 3, 3})
						//get midnode of jb, jb + nbm + 1
						c1 := coords[jb-1]; c2 := coords[jb+nbm]
						coords = append(coords, []float64{(c1[0]+c2[0])/2.0,(c1[1]+c2[1])/2.0})
						je = len(coords)
						ms = append(ms, []int{jb, je, 1, 3, 3})
						ms = append(ms, []int{je, jb+nbm+1, 1, 3, 3})
						ms = append(ms, []int{jb+nbm+2,je,1,4,3})
						ms = append(ms, []int{jb+nbm+3,je,1,4,3})
						case jb == nrs-1:
						je := jb + nbm
						ms = append(ms, []int{jb, je, 1, 3, 3})
						//get midnode of jb, jb + nbm + 2
						c1 := coords[jb-1]; c2 := coords[jb+nbm+1]
						coords = append(coords, []float64{(c1[0]+c2[0])/2.0,(c1[1]+c2[1])/2.0})
						je = len(coords)
						ms = append(ms, []int{jb, je, 1, 3, 3})
						ms = append(ms, []int{je, jb+nbm+2, 1, 3, 3})
						ms = append(ms, []int{jb+nbm+1,je,1,4,3})
						ms = append(ms, []int{jb+nbm,je,1,4,3})
						case jb < nrs-1:
						for _, je := range []int{jb+nbm, jb+nbm+1}{
							if _, ok := tcnodes[je]; ok{
								ms = append(ms, []int{jb, je, 1, 3, 3})
							}
						}
						case jb > nrs:
						for _, je := range []int{jb+nbm+2, jb+nbm+3}{
							if _, ok := tcnodes[je]; ok{
								ms = append(ms, []int{jb, je, 1, 3, 3})
							}
						}
					}
				}
			}
		}
	}
	if t.Bent{
		coords = append(coords, []float64{0,0})
		ms = append(ms, []int{len(coords),bcns[0],1,ngs+1,0})
		coords = append(coords, []float64{t.Span,0})
		ms = append(ms, []int{len(coords),bcns[len(bcns)-1],1,ngs+1,0})
		
		ms = append(ms, []int{len(coords),bcns[0],1,ngs+2,0})
		ms = append(ms, []int{len(coords)-1,bcns[len(bcns)-1],1,ngs+2,0})
		
	}
	var typfac float64
	typfac = 1.0
	if t.Typ == 2{
		typfac = 2.0
	}
	if t.Bcslope > 0.0{
		for _, node := range bcns[1:len(bcns)-1]{
			//fmt.Println("editing sloped node-",node)
			x := coords[node-1][0]; y := coords[node-1][1]
			dy := (t.Span/typfac - math.Abs(t.Span/typfac-x))/t.Bcslope
			coords[node-1] = []float64{x,y+dy}
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
	t.Mod.Coords = t.Coords; t.Mod.Mprp = t.Ms
	t.Mod.Supports = [][]int{{bcns[0],-1,-1},{bcns[len(bcns)-1],0,-1}}
	return
}
