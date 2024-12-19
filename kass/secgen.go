package barf

import (
	"fmt"
	"math"
)

//section generation funcs
//combo and input section entry func not written yet 
/*
	SectionMap = map[int]string{
		-2: "combo",
		-1: "generic/input",
		0:  "circle",
		1:  "rectangle",
		2:  "etri", 
		3:  "rtri", 
		4:  "box", 
		5:  "tube", //_HOW TO
		6:  "T",
		7:  "L",
		8:  "L-left",
		9:  "L-right",
		10: "L-eq",
		11: "plus",
		12: "ieq",
		13: "C",
		14: "T-pocket", //or inverted tee
		15: "pentagon", //reg
		16: "house",    //what else does one call this thing
		17: "hexagon", //reg
		18: "octagon", //reg
                19: "tapered pocket", //allen
                20: "trapezoid", //subramanian
                21: "diamond",
		22: "tapered t",
		23: "haunch2f",//bs5950 fig g2
		24: "haunch3f",//bs5950 fig g3
		25: "built-i",//built up i section - IS NOT EXISTS
		26: "spaced",//(rect) spaced column
		27: "l2-ss",//double (eq/ueq) angles back to back (same side)
		28: "l2-os",//double (eq/ueq) angles on opp. sides
                29: "plate-i",
	}
*/

//Draw plots a section view (gnuplot ofc)
func (s *SectIn) Draw(term string){
	//add if d1 and d2 to draw bars
	var data, ldata, cdata string
	if s.Styp == 0 || s.Styp == 5{
		switch s.Styp{
			case 0:
			//circle
			cdata += fmt.Sprintf("%f %f %f\n",s.Prop.Xc, s.Prop.Yc, s.Dims[0]/2.0)
			data += fmt.Sprintf("%f %f %f %f %f %.f\n",s.Prop.Xc, s.Prop.Yc,s.Dims[0]/2.0, 0.0,2.0, s.Dims[0]/2.0)
			//data += fmt.Sprintf("%f %f %f %f %f %.f\n",s.Prop.Xc, s.Prop.Yc,0.0,s.Prop.Dims[0]/2.0, 2.0, s.Prop.Dims[0]/2.0)
			ldata += fmt.Sprintf("%f %f CX\n",s.Prop.Xc,s.Prop.Yc)
			case 5:
			//tube
			cdata += fmt.Sprintf("%f %f %f\n",s.Prop.Xc, s.Prop.Yc, s.Prop.Dims[0]/2.0)
			cdata += fmt.Sprintf("%f %f %f\n",s.Prop.Xc, s.Prop.Yc, s.Prop.Dims[1]/2.0)
			data += fmt.Sprintf("%f %f %f %f %f %.f\n",s.Prop.Xc, s.Prop.Yc,s.Prop.Dims[0]/2.0, 0.0,2.0, s.Prop.Dims[0]/2.0)
			data += fmt.Sprintf("%f %f %f %f %f %.f\n",s.Prop.Xc, s.Prop.Yc,s.Prop.Dims[0]/2.0, 0.0,3.0, s.Prop.Dims[1]/2.0)
			//data += fmt.Sprintf("%f %f %f %f %f %.f\n",s.Prop.Xc, s.Prop.Yc,0.0,s.Prop.Dims[0]/2.0, 2.0, s.Prop.Dims[0]/2.0)
			//data += fmt.Sprintf("%f %f %f %f %f %.f\n",s.Prop.Xc, s.Prop.Yc,0.0,s.Prop.Dims[1]/2.0,3.0, s.Prop.Dims[1]/2.0)
			ldata += fmt.Sprintf("%f %f CX\n",s.Prop.Xc,s.Prop.Yc)
		}
	} else{
		var strt int
		for i, nc := range s.Ncs{
			wt := s.Wts[i]
			if wt == -1.0{wt = 2.0}
			for j := range s.Coords[strt:strt+nc-1]{
				p1 := s.Coords[strt+j]; p2 := s.Coords[strt+j+1]; l := Dist2d(p1,p2)
				data += fmt.Sprintf("%f %f %f %f %f %.f\n",p1[0],p1[1],p2[0]-p1[0], p2[1]-p1[1],wt, l)
			}
			strt += nc
		}
		ldata += fmt.Sprintf("%f %f CX\n",s.Prop.Xc, s.Prop.Yc)
		//switch s.Styp{
		//	case 0:
		//	//add circle radius
		//	cdata += fmt.Sprintf("%f %f %f\n",s.Coords[2][0]/2.0, s.Coords[2][1]/2.0, s.Prop.Dims[0]/2.0)
		//}
	}
	s.Data = []string{data, ldata, cdata}
	data += "\n\n"; ldata += "\n\n"; cdata += "\n\n"
	data += ldata; data += cdata
	pltskript := "plotsec.gp"
	if term == "dumb"{s.Txtplot = skriptrun(data, pltskript, "dumb", SectionMap[s.Styp], "", "")}
	if term == "mono"{s.Monoplot = skriptrun(data, pltskript, "mono", SectionMap[s.Styp], "", "")}
	if term == "qt"{skriptrun(data, pltskript, "qt", SectionMap[s.Styp], "", "")}
	return
}


//CalcQ calculates statical moment of area Q about the neutral axis (now centroid) of a section
func (s *SectIn) CalcQ(){
	//first in y -
	switch s.Styp{
		case 0:
		//circle, neutral axis is at the centroid
		area := s.Prop.Area/2.0
		//yc = 4r/3pi = 2d/3pi
		dy := 2.0 * s.Prop.Dims[0]/3.0/math.Pi
		s.Prop.Qxx = area * dy
		s.Prop.Vfx = s.Prop.Qxx/s.Prop.Ixx/s.Prop.Dims[0]
		return
		case 5:
		//tube
		//HUH? HUH?
	}
	if s.Prop.Yc == 0.0{
		s.UpdateProp()
	}
	area, _, yc := SecArXu(s, s.Prop.Yc)
	dy := math.Abs(s.Prop.Yc - yc)
	_, dxs, _ := s.GetWidth(s.Prop.Yc)
	var dx float64
	for _, x := range dxs{
		dx += x
	}
	if dx == 0.0 || s.Prop.Ixx == 0.0{
		return
	}
	s.Prop.Qxx = area * dy
	s.Prop.Vfx = s.Prop.Qxx/s.Prop.Ixx/dx
	//now about Y- one must rotate, etc - TODO
	return
}


// //UpdateZxys updates the plastic section modulus about xx and yy
// func (s *SectIn) UpdateZxys(){
// 	//first sort coords to find xmax, ymax
// }

//CircArXu returns the area of a circular section at a depth dck from top
func CircArXu(r, dck float64) (area, xc, yc float64){
	r1 := r - dck
	area = math.Pow(r,2) * math.Acos(r1/r) - r1 * math.Sqrt(math.Pow(r,2)-math.Pow(r1,2))
	//fmt.Println("dck, area",dck, area)
	return
}

//SecArXu returns the area of a section at coord yc
//it computes the point of intersection of the line y = ymax - xu
//with all section polygon line segments etc IS BASICALLY ColSecArXu from mosh
func SecArXu(sec *SectIn, dck float64) (area, xc, yc float64) {
	var nc1s []int
	var wt1s []float64
	var coords [][]float64
	var n1, nc1, n int
	c1 := dck
	var xc1, yc1 float64
	for idx, nc := range sec.Ncs {
		var pts [][]float64
		ptmap := make(map[Pt]int)
		if idx == 0 {
			n1 = 0
		} else {
			n1 = sec.Ncs[idx-1]
		}
		nc1 = 0; n = 0
		for i := range sec.Coords[n1 : n1+nc-1] {
			i = i + n1
			pta := Pt{sec.Coords[i][0], sec.Coords[i][1]}
			ptb := Pt{sec.Coords[i+1][0], sec.Coords[i+1][1]}
			if (pta.Y < c1 && ptb.Y < c1) {
				continue
			}
			if pta.Y >= c1 {
				if _, ok := ptmap[pta]; !ok {
					ptmap[pta] = idx
					pts = append(pts, []float64{pta.X, pta.Y})
					nc1++
					xc1 += pta.X; yc1 += pta.Y; n++
				}
			}
			if ptb.Y >= c1 {
				if _, ok := ptmap[ptb]; !ok {
					ptmap[ptb] = idx
					pts = append(pts, []float64{ptb.X, ptb.Y})
					nc1++
					xc1 += ptb.X; yc1 += ptb.Y; n++
				}
			}
			if pta.Y - ptb.Y == 0 {
				continue
			}
			a2 := ptb.Y - pta.Y
			b2 := pta.X - ptb.X
			c2 := a2 * pta.X + b2 * pta.Y 
			xin := (c2 - b2*c1)/a2
			yin := c1
			ptx := Pt{xin, yin}
			if (pta.X <= xin && xin <= ptb.X) || (pta.X >= xin && xin >= ptb.X){
				if _, ok := ptmap[ptx]; !ok {
					pts = append(pts, []float64{ptx.X, ptx.Y})
					ptmap[ptx] = idx
					nc1++
					xc1 += ptx.X; yc1 += ptx.Y; n++
				}
			}
		}
		if len(pts) == 0{
			continue
		}
		SortCcw(pts, xc1/float64(n), yc1/float64(n))
		pts = append(pts, pts[0])
		nc1++
		nc1s = append(nc1s, nc1)
		coords = append(coords, pts...)
		wt1s = append(wt1s, sec.Wts[idx])
	}
	area, xc, yc, _, _, _, _, _, _ = SecPrp(nc1s, wt1s, coords)
	return
}


//RandSec generates a section given dims within a step - for opt routines
func RandSec(styp int, step float64, limits, dims []float64) (s SectIn){
	//TODO-for opt routines
	return
}

//SecGen returns a SectIn given a section type and dimensions
//TODO - ADD ERROR
func SecGen(styp int, dims []float64) (s SectIn){
	//all sections are 2-d point arrays
	//this might be redundant now but useful later? esp for non standard sections i guess
	var ncs []int
	var wts []float64
	var coords [][]float64
	var solid bool
	var bar Secprop
	var ymax float64
	solid = true
	switch styp{
		case 0:
		if len(dims) < 1{return}
		D := dims[0]
		bar = CalcSecProp(styp, dims)
		solid = true
		ymax = D
		s = SectIn{
			Prop:bar,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
		}
		return
		case 1:
		//rect
		if len(dims) < 2{return}
		b := dims[0]; d := dims[1]
		ncs = []int{5}
		wts = []float64{1}
		coords = [][]float64{{0,0},{b,0},{b,d},{0,d},{0,0}}
		bar = CalcSecProp(styp, dims)
		ymax = d
		case 2:
		//e tri
		
		if len(dims) < 1{return}
		b := dims[0]
		h := b * math.Tan(math.Pi/3.0)/2.0
		ncs = []int{4}
		wts = []float64{1.0}
		coords = [][]float64{{0,0},{b,0},{b/2.0,h},{0,0}}
		bar = CalcSecProp(styp, dims)
		ymax = h
		case 3:
		//r tri
		b := dims[0]
		h := dims[1]
		
		if len(dims) < 2{return}
		ncs = []int{4}
		wts = []float64{1.0}
		coords = [][]float64{{0,0},{b, 0},{0, h},{0,0}}
		bar = CalcSecProp(styp, dims)
		ymax = h
		case 4:
		//box section aha
		
		if len(dims) < 4{return}
		B := dims[0]; D := dims[1]; b := dims[2]; d := dims[3]
		ncs = []int{5,5}
		wts = []float64{1,-1}
		bsx := (B-b)/2.0; bsy := (D-d)/2.0
		coords = [][]float64{{0,0},{B,0},{B,D},{0,D},{0,0},{bsx,bsy},{bsx+b,bsy},{bsx+b,bsy+d},{bsx,bsy+d},{bsx,bsy}}
		bar = CalcSecProp(styp, dims)
		solid = false
		ymax = D
		case 5:
		//TUBE HOW
		//eh copy secprop for now
		
		if len(dims) < 2{return}
		D := dims[0]
		bar = CalcSecProp(styp, dims)
		solid = false
		ymax = D
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Prop:bar,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
			Ym:-D,
		}
		return
		case 6:
		//t- section
		
		if len(dims) < 4{return}
		bf := dims[0]; d := dims[1]; bw := dims[2]; df := dims[3]
		wts = []float64{1}
		coords = [][]float64{
			{bf/2.0 - bw/2.0,0},
			{bf/2.0 + bw/2.0,0},
			{bf/2.0 + bw/2.0, d - df},
			{bf , d - df},
			{bf ,d},
			{0, d},
			{0, d - df},
			{bf/2.0 - bw/2.0,d - df},
			{bf/2.0 - bw/2.0,0},
		}
		ncs = []int{len(coords)}
		bar = CalcSecProp(styp, dims)
		ymax = d
		case 7:
		//bottom left origin l
		
		if len(dims) < 4{return}
		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		ncs = []int{7}
		wts = []float64{1.0}
		coords = [][]float64{
			{0,0},
			{bf,0},
			{bf,df},
			{bw,df},
			{bw,d},
			{0,d},
			{0,0},
		}
		ymax = d
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
		}
		s.UpdateProp()
		return
		case 8:
		//l - left
		if len(dims) < 4{return}
		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		ncs = []int{7}
		wts = []float64{1.0}
		coords = [][]float64{
			{bf-bw,0},
			{bf,0},
			{bf,d},
			{0,d},
			{0,d-df},
			{bf-bw,d-df},
			{bf-bw,0},
		}
		bar = CalcSecProp(styp, dims)
		ymax = d
		case 9:
		//l- right
		
		if len(dims) < 4{return}
		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		ncs = []int{7}
		wts = []float64{1.0}
		coords = [][]float64{
			{0,0},
			{bw,0},
			{bw,d-df},
			{bf,d-df},
			{bf,d},
			{0,d},
			{0,0},
		}
		bar = CalcSecProp(styp, dims)
		ymax = d
		case 10:
		//l - eq
		
		if len(dims) < 3{return}
		b := dims[0]; d := dims[1]; t := dims[2]
		wts = []float64{1}; ncs = []int{7}
		coords = [][]float64{{0,0},{b,0},{b,t},{t,t},{t,d},{0,d},{0,0}}
		bar = CalcSecProp(styp, dims)
		ymax = d
		case 11:
		//plus
		
		if len(dims) < 4{return}
		b := dims[0]; d := dims[1]; s := dims[2]; t := dims[3]
		x1 := 0.0; x2 := (b-t)/2.0; x3 := (b+t)/2.0; x4 := b
		y1 := 0.0; y2 := (d-s)/2.0; y3 := (d+s)/2.0; y4 := d
		wts = []float64{1}; ncs = []int{13}
		coords = [][]float64{{x2,y1},{x3,y1},{x3,y2},{x4,y2},{x4,y3},{x3,y3},{x3,y4},{x2,y4},{x2,y3},{x1,y3},{x1,y2},{x2,y2},{x2,y1}}
		bar = CalcSecProp(styp, dims)
		ymax = d
		case 12:
		//equal i section I (ieq)// built-up i section
		
		if len(dims) < 4{return}
		b := dims[0]
		h := dims[1]
		tf := dims[2]
		tw := dims[3]
		wts = []float64{1}
		ncs = []int{13}
		x1 := 0.0; x2 := (b - tw)/2.0; x3 := (b+tw)/2.0; x4 := b
		y1 := 0.0; y2 := tf; y3 := h-tf; y4 := h
		coords = [][]float64{{x1,y1},{x4,y1},{x4,y2},{x3,y2},{x3,y3},{x4,y3},{x4,y4},{x1,y4},{x1,y3},{x2,y3},{x2,y2},{x1,y2},{x1,y1}}
		
		bar = CalcSecProp(12, dims)
		ymax = h
		
		case 13:
		//c ( [ ) section
		
		if len(dims) < 4{return}
		b := dims[0]
		h := dims[1]
		tf := dims[2]
		tw := dims[3]
		
		wts = []float64{1}
		ncs = []int{9}
		x1 := 0.0; x2 := tw; x3 := b
		y1 := 0.0; y2 := tf; y3 := h - tf; y4 := h
		coords = [][]float64{{x1,y1},{x3,y1},{x3,y2},{x2,y2},{x2,y3},{x3,y3},{x3,y4},{x1,y4},{x1,y1}}
		bar = CalcSecProp(styp, dims)
		ymax = h
		case 14:
		//t pocket
		
		if len(dims) < 4{return}
		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		ncs = []int{10}
		wts = []float64{1.0}
		coords = [][]float64{
			{0,0},
			{bw,0},
			{bw,d-df},
			{bw/2.0+bf/2.0,d-df},
			{bw/2.0+bf/2.0,d-df},
			{bw/2.0+bf/2.0,d},
			{bw/2.0-bf/2.0,d},
			{bw/2.0-bf/2.0,d-df},
			{0,d-df},
			{0,0},
		}
		bar = CalcSecProp(styp, dims)
		ymax = d
		case 15:
		//pentagon
		
		if len(dims) < 1{return}
		b := dims[0]
		wts = []float64{1}
		for i := 0; i < 6; i++{
			theta := float64(i) * 2.0 * math.Pi/5.0
			coords = append(coords, []float64{b * math.Cos(theta), b * math.Sin(theta)})
		}
		
		ncs = []int{len(coords)}
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
		}
		s = SecRotate(s, 90.0)
		s.UpdateProp()
		return
		case 16:
		//"house"
		
		if len(dims) < 1{return}
		b := dims[0]
		wts = []float64{1}
		coords = [][]float64{{0,0},{b,0},{b,b},{b/2,b+b*math.Sqrt(3.0)/2.0},{0,b},{0,0}}
		ncs = []int{len(coords)}
		ymax = b + b * math.Sqrt(3.0/2.0)
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
		}
		s.UpdateProp()
		return
		case 17:
		//hexagon
		
		if len(dims) < 1{return}
		b := dims[0]
		wts = []float64{1}
		for i := 0; i < 7; i++{
			theta := float64(i) * math.Pi * 2.0/6.0
			coords = append(coords, []float64{b*math.Cos(theta),b*math.Sin(theta)})
		}
		ncs = []int{len(coords)}
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
		}
		s.UpdateProp()
		return
		case 18:
		
		if len(dims) < 1{return}
		//octagon
		b := dims[0]
		r := b * math.Sqrt(1.0 + 1.0/math.Sqrt(2))
		wts = []float64{1}
		for i := 0; i < 9; i++{
			theta := float64(i) * math.Pi * 2.0/8.0
			coords = append(coords, []float64{r*math.Cos(theta),r*math.Sin(theta)})
		}
		ncs = []int{len(coords)}
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
		}
		s = SecRotate(s, 22.5)
		s.UpdateProp()
		return
		case 19:
		//tapered pocket section (allen 5.2)
		
		if len(dims) < 5{return}
		bf := dims[0]; d := dims[1]; bw := dims[2]; df := dims[3]; bp := dims[4]
		wts = []float64{1.0}
		x1 := 0.0; x2 := bf/2.0 + bp - bw/2.0; x3 := x2 + bw; x4 := x1 + bp; x5 := x4 + bf; x6 := x5 + bp
		y1 := 0.0; y2 := d - df; y3 := d
		coords = [][]float64{{x2,y1},{x3,y1},{x6,y2},{x5,y2},{x5,y3},{x4,y3},{x4,y2},{x1,y2},{x2,y1}}
		ncs = []int{len(coords)}
		ymax = d
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
		}
		s.UpdateProp()
		return
		case 20:
		//trapezoidal section (subramanian 5.7)
		//slopes over d from bw (bottom) to bf(top)
		
		if len(dims) < 3{return}
		bf := dims[0]; d := dims[1]; bw := dims[2]
		wts = []float64{1.0}
		coords = [][]float64{{0,0},{bw,0},{bw + (bf - bw)/2.0, d},{-(bf-bw)/2.0,d},{0,0}}
		ncs = []int{len(coords)}
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
		}
		s.UpdateProp()
		return
		case 21:
		//(square) diamond section
		
		if len(dims) < 1{return}
		b := dims[0]
		wts = []float64{1.0}
		
		for i := 0; i < 5; i++{
			theta := float64(i) * math.Pi/2.0
			coords = append(coords, []float64{b * math.Cos(theta), b * math.Sin(theta)})	
		}
		ncs = []int{len(coords)}
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
		}
		s.UpdateProp()
		return
		case 22:
		//tapered t section
		if len(dims) < 5 {return}
		bf := dims[0]; d := dims[1]; T := dims[2]; t := dims[3]; df := dims[4]
		wts = []float64{1.0}
		x1 := 0.0; x2 := (bf - T)/2.0; x3 := (bf - t)/2.0; x4 := (bf + t)/2.0; x5 := (bf + T)/2.0; x6 := bf
		y1 := 0.0; y2 := d - df; y3 := d
		coords = [][]float64{{x3,y1},{x4,y1},{x5,y2},{x6,y2},{x6,y3},{x1,y3},{x1,y2},{x2,y2},{x3,y1}}
		ncs = []int{len(coords)}
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
		}
		s.UpdateProp()
		return
		case 23:
		//two flange (ieq) haunch
		if len(dims) < 5{return}
		b := dims[0]; h := dims[1]; tf := dims[2]; tw := dims[3]; dy := dims[4]
		H := h - tf + dy
		d1 := []float64{b, H, tf, tw}
		s = SecGen(12,d1)
		s.Prop.Ixx = s.Prop.Ixx + s.Prop.Area * math.Pow(dy/2.0,2)
		return
		case 24:
		//three flange (ieq) haunch
		if len(dims) < 5{return}
		b := dims[0]; h := dims[1]; tf := dims[2]; tw := dims[3]; dy := dims[4]
		d1 := []float64{b, h, tf, tw}
		s  = SecGen(12,d1)
		d2 := []float64{b, dy, tf, tw}
		s2 := SecGen(6,d2)
		dy = s2.Prop.Yc + h/2.0
		s.Prop.Ixx = s.Prop.Ixx + s2.Prop.Area * math.Pow(dy/2.0,2)
		return
		case 26:
		//spaced column
		//b - width, d - depth, w - spacing
		if len(dims) < 3{return}
		b := dims[0]; d := dims[1]; w := dims[2]
		wts = []float64{1.0, 1.0}
		ncs = []int{5, 5}
		coords = [][]float64{{0,0},{b,0},{b,d},{0,d},{0,0},{0,d+w},{b,d+w},{b,2.0*d+w},{0,2.0*d+w},{0,d+w}}
		solid = false
		ymax = 2.0 * d + w
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
		}
		s.UpdateProp()
		return
		
		case 28:
		//l2-os
		if len(dims) < 5{return}
		bf := dims[0]; d := dims[1]; bw := dims[2]; tf := dims[3]; tp := dims[4]
		wts = []float64{1.0, 1.0}
		ncs = []int{7, 7}
		x1 := 0.0; x2 := bf - bw; x3 := bf
		y1 := 0.0; y2 := d - tf; y3 := d
		coords = [][]float64{
			{x1, y3},
			{x1, y2},
			{x2, y2},
			{x2, y1},
			{x3, y1},
			{x3, y3},
			{x1, y3},
		}
		x1 = bf + tp; x2 = x1 + bw; x3 = x1 + bf
		c2 := [][]float64{
			{x1, y1},
			{x2, y1},
			{x2, y2},
			{x3, y2},
			{x3, y3},
			{x1, y3},
			{x1, y1},
		}
		coords = append(coords, c2...)
		ymax = y3
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
			Xmx:x3,
		}
		s.UpdateProp()
		return
		case 27:
		//l2-ss
		if len(dims) < 4{return}
		bf := dims[0]; d := dims[1]; bw := dims[2]; tf := dims[3]
		wts = []float64{1.0, 1.0}
		ncs = []int{7, 7}
		x1 := 0.0; x2 := bw; x3 := bf
		y1 := 0.0; y2 := d - tf; y3 := d; y4 := d + tf; y5 := 2.0 * d
		coords = [][]float64{
			{x1, y1},
			{x2, y1},
			{x2, y2},
			{x3, y2},
			{x3, y3},
			{x1, y3},
			{x1, y1},
		}
		c2 := [][]float64{
			{x1, y3},
			{x3, y3},
			{x3, y4},
			{x2, y4},
			{x2, y5},
			{x1, y5},
			{x1, y3},
		}
		coords = append(coords, c2...)
		ymax = y3
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Solid:solid,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
			Xmx:x3,
		}
		s.UpdateProp()
		return
		case 29:
		//i section ()with cover plates of b, d
		if len(dims) < 6{return}
		//i section coords
		b := dims[0]
		h := dims[1]
		tf := dims[2]
		tw := dims[3]
		
		x1 := 0.0; x2 := (b - tw)/2.0; x3 := (b+tw)/2.0; x4 := b
		y1 := 0.0; y2 := tf; y3 := h-tf; y4 := h
		coords = [][]float64{{x1,y1},{x4,y1},{x4,y2},{x3,y2},{x3,y3},{x4,y3},{x4,y4},{x1,y4},{x1,y3},{x2,y3},{x2,y2},{x1,y2},{x1,y1}}
		//then top cover plate
		
		//cover plate width, thickness
		B := dims[4]
		D := dims[5]
		x1 = - (B - b)/2.0; x2 = b + (B - b)/2.0
		y1 = -D; y2 = 0.0; y3 = h; y4 = h + D
		cbs := [][]float64{{x1,y1},{x2,y1},{x2,y2},{x1,y2},{x1,y1}} 
		cts := [][]float64{{x1,y3},{x2,y3},{x2,y4},{x1,y4},{x1,y3}}
		wts = []float64{1,1,1}
		ncs = []int{13,5,5}
		coords = append(coords, cts...)
		coords = append(coords, cbs...)
		ymax = h + D
		s = SectIn{
			Ncs:ncs,
			Wts:wts,
			Coords:coords,
			Prop:bar,
			Solid:true,
			Styp:styp,
			Dims:dims,
			Ymx:ymax,
		}
		s.UpdateProp()
		return
	}
	s = SectIn{
		Ncs:ncs,
		Wts:wts,
		Coords:coords,
		Prop:bar,
		Solid:solid,
		Styp:styp,
		Dims:dims,
		Ymx:ymax,
	}
	return
}
