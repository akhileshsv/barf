package barf

import (
	"fmt"
	"math"
)

//section generation funcs
//combo and input section entry func not written yet 
/*
SectionMap:
		-2: "combo",
		-1: "generic/input",
		0:  "circle",
		1:  "rectangle",
		2:  "etri", 
		3:  "rtri", 
		4:  "box", 
		5:  "tube", 
		6:  "T",
		7:  "L",
		8:  "L-left",
		9:  "L-right",
		10: "L-eq",
		11: "plus",
		12: "ieq",
		13: "C",
		14: "T-pocket",
		15: "H-eq",
		16: "house", regular
		17: "hexagon", reg
		18: "octagon", reg
                19: "tapered pocket" allen
                20: "trapezoid", subramanian
                21: "diamond"
                22: "pentagon", reg
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
			if i == 0{
				strt = 0
			} else {strt = s.Ncs[i-1]}
			for j := range s.Coords[strt:strt+nc-1]{
				p1 := s.Coords[strt+j]; p2 := s.Coords[strt+j+1]; l := Dist2d(p1,p2)
				data += fmt.Sprintf("%f %f %f %f %f %.f\n",p1[0],p1[1],p2[0]-p1[0], p2[1]-p1[1],wt, l)
			}
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
	if term == "dumb"{s.Txtplot = skriptrun(data, pltskript, "dumb", SectionMap[s.Styp], "")}
	if term == "mono"{s.Monoplot = skriptrun(data, pltskript, "mono", SectionMap[s.Styp], "")}
	if term == "qt"{skriptrun(data, pltskript, "qt", SectionMap[s.Styp], "")}
	return
}

//RandSec generates a section given dims within a step - for opt routines
func RandSec(styp int, step float64, limits, dims []float64) (s SectIn){
	//TODO-for opt routines
	return
}

//SecGen returns a SectIn given a section type and dimensions
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
		b := dims[0]; d := dims[1]
		ncs = []int{5}
		wts = []float64{1}
		coords = [][]float64{{0,0},{b,0},{b,d},{0,d},{0,0}}
		bar = CalcSecProp(styp, dims)
		ymax = d
		case 2:
		//e tri
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
		ncs = []int{4}
		wts = []float64{1.0}
		coords = [][]float64{{0,0},{b, 0},{0, h},{0,0}}
		bar = CalcSecProp(styp, dims)
		ymax = h
		case 4:
		//box section aha
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
		//yes yes tis redundant
		bar = CalcSecProp(styp, dims)
		ymax = d
		case 8:
		//l - left
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
		b := dims[0]; d := dims[1]; t := dims[2]
		wts = []float64{1}; ncs = []int{7}
		coords = [][]float64{{0,0},{b,0},{b,t},{t,t},{t,d},{0,d},{0,0}}
		bar = CalcSecProp(styp, dims)
		ymax = d
		case 11:
		//plus
		b := dims[0]; d := dims[1]; s := dims[2]; t := dims[3]
		x1 := 0.0; x2 := (b-t)/2.0; x3 := (b+t)/2.0; x4 := b
		y1 := 0.0; y2 := (d-s)/2.0; y3 := (d+s)/2.0; y4 := d
		wts = []float64{1}; ncs = []int{13}
		coords = [][]float64{{x2,y1},{x3,y1},{x3,y2},{x4,y2},{x4,y3},{x3,y3},{x3,y4},{x2,y4},{x2,y3},{x1,y3},{x1,y2},{x2,y2},{x2,y1}}
		bar = CalcSecProp(styp, dims)
		ymax = d
		case 12:
		//equal i section I
		b := dims[0]
		h := dims[1]
		tf := dims[2]
		tw := dims[3]
		wts = []float64{1}
		ncs = []int{13}
		x1 := 0.0; x2 := (b - tw)/2.0; x3 := (b+tw)/2.0; x4 := b
		y1 := 0.0; y2 := tf; y3 := h-tf; y4 := h
		coords = [][]float64{{x1,y1},{x4,y1},{x4,y2},{x3,y2},{x3,y3},{x4,y3},{x4,y4},{x1,y4},{x1,y3},{x2,y3},{x2,y2},{x1,y2},{x1,y1}}
		bar = CalcSecProp(styp, dims)
		ymax = h
		case 13:
		//c ( [ ) section
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
