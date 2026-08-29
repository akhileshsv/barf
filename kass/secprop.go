package barf

import (
	"math"
)

//CalcSecProp takes in styp and dims and calculates section properties
//mostly rips off structx.com formulae
//okay so this has all formula based stuff

func CalcSecProp(styp int, dims []float64) (bar Secprop) {
	//returns a section 

	var area, perimeter, ixx, iyy, iww, xc, yc, sxx, syy, zxx, zyy, j, rxx, ryy float64
	switch SectionMap[styp]{
		
		case "rectangle":
		//styp 1
		b := dims[0]
		d := dims[1]
		area = b * d
		perimeter = 2 * (b + d)
		ixx = b * math.Pow(d, 3) / 12.0
		iyy = d * math.Pow(b, 3) / 12.0
		xc = b / 2.0
		yc = d / 2.0
		sxx = ixx / yc
		syy = iyy / xc
		j = math.Pow(b, 3) * d * ((1.0 / 3.0) - 0.21*(b/d)*(1.0-(1.0/12.0)*math.Pow((b/d), 4)))
		zxx = b * math.Pow(d, 2) / 2.0
		zyy = d * math.Pow(b, 2) / 2.0
		rxx = math.Sqrt(ixx / area)
		ryy = math.Sqrt(iyy / area)

		case "box":
		B := dims[0]
		D := dims[1]
		b := dims[2]
		d := dims[3]
		area = B*D - b*d
		perimeter = 2 * (B + D)
		ixx = (B * math.Pow(D, 3) / 12.0) - (b * math.Pow(d, 3) / 12.0)
		iyy = (D * math.Pow(B, 3) / 12.0) - (d * math.Pow(b, 3) / 12.0)
		xc = B / 2.0
		yc = D / 2.0
		sxx = ixx / yc
		syy = iyy / xc
		tf := D - d
		tw := B - b
		j = 2 * math.Pow(B, 2) * math.Pow(D, 2) / (B/tf + D/tw)
		zxx = (B*math.Pow(D, 2) - b*math.Pow(d, 2)) / 4.0
		zyy = (D*math.Pow(B, 2) - d*math.Pow(b, 2)) / 4.0
		rxx = math.Sqrt(ixx / area)
		ryy = math.Sqrt(iyy / area)

		case "circle":
		d := dims[0]
		area = math.Pi * math.Pow(d, 2) / 4.0
		perimeter = math.Pi * d
		ixx = math.Pi * math.Pow(d, 4) / 64.0
		iyy = math.Pi * math.Pow(d, 4) / 64.0
		xc = d / 2.0
		yc = d / 2.0
		sxx = ixx / yc
		syy = iyy / xc
		j = math.Pi * math.Pow(d, 4) / 32.0
		zxx = math.Pow(d, 3) / 6.0
		zyy = math.Pow(d, 3) / 6.0
		rxx = math.Sqrt(ixx / area)
		ryy = math.Sqrt(iyy / area)

		case "tube":

		D := dims[0]
		d := dims[1]
		area = math.Pi * (math.Pow(D, 2) - math.Pow(d, 2))/ 4.0
		perimeter = math.Pi * D
		ixx = math.Pi * (math.Pow(D, 4) - math.Pow(d, 4)) / 64.0
		iyy = math.Pi * (math.Pow(D, 4) - math.Pow(d, 4)) / 64.0
		xc = d / 2.0
		yc = d / 2.0
		sxx = ixx / yc
		syy = iyy / xc
		j = math.Pi * (math.Pow(D, 4) - math.Pow(d, 4)) / 32.0
		zxx = (math.Pow(D, 3) - math.Pow(d, 3)) / 6.0
		zyy = (math.Pow(D, 3) - math.Pow(d, 3)) / 6.0
		rxx = math.Sqrt(ixx / area)
		ryy = math.Sqrt(iyy / area)
		
		case "L-eq":
		//L - section of equal thickness

		b := dims[0]
		d := dims[1]
		t := dims[2]
		b1 := b - t
		d1 := d - t
		area = t * (b + d - t)
		perimeter = 2 * (b + d)
		xc = (math.Pow(b, 2) + d*t - math.Pow(t, 2)) / (2 * (b + d - t))
		yc = (math.Pow(d, 2) + b*t - math.Pow(t, 2)) / (2 * (b + d - t))
		ixx = (1.0/3.0)*(b*math.Pow(d, 3)-(b-t)*math.Pow((d-t), 3)) - area*math.Pow((d-yc), 2)
		iyy = (1.0/3.0)*(d*math.Pow(b, 3)-(d-t)*math.Pow((b-t), 3)) - area*math.Pow((b-xc), 2)
		sxx = ixx / yc
		syy = iyy / xc
		j = (1.0 / 3.0) * (b*math.Pow(t, 3) + (d-t)*math.Pow(t, 3))
		if t < area/(2*b) {
			zxx = t * (math.Pow(d1, 2) - math.Pow(b, 2) + 2*b*d) / 4.0
		} else {
			zxx = b*math.Pow(t, 2)/4.0 + d*t*d1/2.0 - math.Pow(t, 2)*math.Pow(d1, 2)/(4*b)
		}
		if t < area/(2*d) {
			zyy = t * (math.Pow(b1, 2) - math.Pow(d, 2) + 2*b*d) / 4.0
		} else {
			zyy = (d * math.Pow(t, 2) / 4.0) + b*t*b1/2.0 - math.Pow(t, 2)*math.Pow(b1, 2)/(4.0*d)
		}
		rxx = math.Sqrt(ixx / area)
		ryy = math.Sqrt(ixx / area)

		case "plus":

		b := dims[0]
		d := dims[1]
		s := dims[2]
		t := dims[3]
		area = d*t + s*(b-t)
		perimeter = 2 * (b + d)
		xc = b / 2.0
		yc = d / 2.0
		ixx = (t*math.Pow(d, 3) + math.Pow(s, 3)*(b-t)) / 12.0
		iyy = (s*math.Pow(b, 3) + math.Pow(t, 3)*(d-s)) / 12.0
		sxx = ixx / yc
		syy = iyy / xc
		j = (1.17 / 3.0) * (b*math.Pow(s, 3) + (d-s)*math.Pow(t, 3))
		rxx = math.Sqrt(ixx / area)
		ryy = math.Sqrt(ixx / area)
		zxx = sxx
		zyy = syy

		case "rtri":

		b := dims[0]
		h := dims[1]
		ixx = b * math.Pow(h,3)/36.0
		iyy = h * math.Pow(b,3)/36.0
		area = b * h /2.0
		perimeter = b + h + math.Sqrt(b*b + h*h)
		xc = b/3.0
		yc = h/3.0
		sxx = b * math.Pow(h,2)/24.0

		case "i":
		//https://calcresource.com/cross-section-doubletee.html
		b := dims[0]
		h := dims[1]
		tf := dims[2]
		tw := dims[3]
		hw := h - 2.0 * tf
		area = 2.0 * b * tf + (h - 2.0 * tf) * tw
		perimeter = 4.0 * b + 2.0 * h - 2.0 * tw
		ixx = b * math.Pow(h,3)/12.0 - (b - tw) * math.Pow(h - 2.0 * tf,3)/12.0
		iyy = tf * math.Pow(b, 3)/6.0 + (h - 2.0*tf) * math.Pow(tw, 3)/12.0
		xc = b/2.0
		yc = h/2.0 //CHECK THIS
		sxx = 2.0 * ixx /h
		syy = 2.0 * iyy/b
		rxx = math.Sqrt(ixx/area)
		ryy = math.Sqrt(iyy/area)
		zxx = b * h* h/4.0 - (b - tw) * hw * hw/4.0
		zyy = tf * b * b/2.0 + hw * tw * tw/4.0
		j = 1.0/3.0 * (2.0 * b * math.Pow(tf,3) + hw * math.Pow(tw,3))
		//j = 1.3/3.0 * (2.0 * b * math.Pow(tf,3) + hw * math.Pow(tw,3))
		case "C":
		b := dims[0]
		h := dims[1]
		tf := dims[2]
		tw := dims[3]
		hw := h - 2.0 * tf
		bf := b - tw
		area = 2.0 * b * tf + hw * tw
		xc = (hw * tw * tw/2.0 + tf * b * b)/area
		perimeter = 4.0 * b + 2.0 * h - 2.0 * tw
		iy0 := hw * math.Pow(tw,3)/3.0 + 2.0 * tf * math.Pow(b, 3)/3.0
		ixx = b * math.Pow(h, 3)/12.0 - (b - tw) * math.Pow(hw, 3)/12.0
		iyy = iy0 - area * math.Pow(xc,2)
		sxx = 2.0 * ixx/h
		syy = iyy / (b - xc)
		zxx = b * h * h/4.0 - bf * hw * hw/4.0
		if tw <= area/2.0/h{
			zyy = tf * bf * bf/2.0 + b * h * tw/2.0 - math.Pow(h,2) * math.Pow(tw,2)/8.0/tf
		} else {
			zyy = (4.0 * tf * b * b * (h- tf) + tw * tw * (math.Pow(h,2) - 4.0 * math.Pow(tf,2)) - 4.0 * b * tf * hw * tw)/4.0/h
		}
		j = 1.0/3.0 * (2.0 * b * math.Pow(tf, 3) + hw * math.Pow(tw, 3))
		
		case "L":
		//bottom left origin L
		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		dw := d - df
		ncs := []int{7}
		wts := []float64{1.0}
		coords := [][]float64{
			{0,0},
			{bf,0},
			{bf,df},
			{bw,df},
			{bw,d},
			{0,d},
			{0,0},
		}
		area, xc, yc, ixx, iyy, _, _, _, _ = SecPrp(ncs, wts, coords)
		j = 1.0/3.0 * (bf * math.Pow(df, 3) + dw * math.Pow(bw, 3))

		case "T":
		//https://amesweb.info/section/moment-of-inertia-of-t-section.aspx

		bf := dims[0] //B
		h := dims[1]
		tw := dims[2] //b
		tf := dims[3] //h (dused)
		hw := h - tf //H

		area = bf*tf + hw*tw
		perimeter = 2 * (bf + tf + hw)
		xc = bf / 2.0
		yc = ((hw+tf/2.0)*tf*bf + hw*hw*tw/2.0) / area
		ixx = tw*hw*math.Pow((yc-hw/2.0), 2) + tw*math.Pow(hw, 3)/12.0 + tf*bf*math.Pow(((hw+tf/2.0)-yc), 2) + bf*math.Pow(tf, 3)/12.0
		iyy = (math.Pow(tw, 3)*hw + math.Pow(bf, 3)*tf) / 12.0
		sxx = ixx / yc
		syy = iyy / xc
		//???
		//j = (1.12 / 3.0) * (bf*math.Pow(tf, 3) + hw*math.Pow(tf, 3))
		if tf < area/(2*bf) {
			zxx = tw*math.Pow(hw, 2)/4.0 + bf*h*tf/2.0 - math.Pow(bf, 2)*math.Pow(tf, 2)/(4.0*tw)
		} else {
			zxx = (tw * math.Pow(h, 2) / 2.0) + (bf * math.Pow(tf, 2) / 4.0) - (h * tf * tw / 2.0) - math.Pow(hw, 2)*math.Pow(tw, 2)/(4*bf)
		}
		zyy = (tf*math.Pow(bf, 2) + hw*math.Pow(tw, 2)) / 4.0
		rxx = math.Sqrt(ixx / area)
		ryy = math.Sqrt(ixx / area)
		//ref - j
		//https://www.structx.com/Shape_Formulas_006.html
		
		j = (bf * math.Pow(tf, 3) + h - (tf/2.0) * math.Pow(tw, 3))/3.0
		case "L-right":
		//flange on right |-

		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		dw := d - df
		ncs := []int{7}
		wts := []float64{1.0}
		coords := [][]float64{
			{0,0},
			{bw,0},
			{bw,d-df},
			{bf,d-df},
			{bf,d},
			{0,d},
			{0,0},
		}
		area, xc, yc, ixx, iyy, _, _, _, _ = SecPrp(ncs, wts, coords)
		
		//http://pont.ist/torsion-constant/
		var b1, t1, b2, t2, k1, k2 float64
		if bf >= df {b1 = bf; t1 = df} else {b1 = df; t1 = bf}
		if bw >= dw {b2 = bw; t2 = dw} else {b2 = dw; t2 = bw}
		k1 = 1.0 - 0.63 * t1/b1 + 0.052 * math.Pow(t1/b1,2)
		k2 = 1.0 - 0.63 * t2/b2 + 0.052 * math.Pow(t2/b2,2)
		j = 1.0/3.0 * (k1 * b1 * math.Pow(t1, 3) + k2 * b2 * math.Pow(t2, 3))
		
		case "L-left":
		//flange on left -|

		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		dw := d - df
		ncs := []int{7}
		wts := []float64{1.0}
		coords := [][]float64{
			{bf-bw,0},
			{bf,0},
			{bf,d},
			{0,d},
			{0,d-df},
			{bf-bw,d-df},
			{bf-bw,0},
		}
		area, xc, yc, ixx, iyy, _, _, _, _ = SecPrp(ncs, wts, coords)
		//j = 1.0/3.0 * (bf * math.Pow(df, 3) + dw * math.Pow(bw, 3))
		
		//j ref
		//http://pont.ist/torsion-constant/
		var b1, t1, b2, t2, k1, k2 float64
		if bf > df {b1 = bf; t1 = df} else {b1 = df; t1 = bf}
		if bw > dw {b2 = bw; t2 = dw} else {b2 = dw; t2 = bw}
		k1 = 1.0 - 0.63 * t1/b1 + 0.052 * math.Pow(t1/b1,2)
		k2 = 1.0 - 0.63 * t2/b2 + 0.052 * math.Pow(t2/b2,2)
		j = 1.0/3.0 * (k1 * b1 * math.Pow(t1, 3) + k2 * b2 * math.Pow(t2, 3))
		
		case "T-pocket":
		//bf < bw
		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		dw := d - df
		ncs := []int{10}
		wts := []float64{1.0}
		coords := [][]float64{
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
		area, xc, yc, ixx, iyy, _, _, _, _ = SecPrp(ncs, wts, coords)
		//j = 1.0/3.0 * (bf * math.Pow(df, 3) + dw * math.Pow(bw, 3))
		
		//THIS IS IMPOSSIBLY WRONG TO JUST COPY PASTE THE T SECTION FORMULA
		
		var b1, t1, b2, t2, k1, k2 float64
		if bf > df {b1 = bf; t1 = df} else {b1 = df; t1 = bf}
		if bw > dw {b2 = bw; t2 = dw} else {b2 = dw; t2 = bw}
		k1 = 1.0 - 0.63 * t1/b1 + 0.052 * math.Pow(t1/b1,2)
		k2 = 1.0 - 0.63 * t2/b2 + 0.052 * math.Pow(t2/b2,2)
		j = 1.0/3.0 * (k1 * b1 * math.Pow(t1, 3) + k2 * b2 * math.Pow(t2, 3))

		case "etri":
		b := dims[0]
		h := b * math.Tan(math.Pi/3.0)/2.0
		ncs := []int{4}
		wts := []float64{1.0}
		coords := [][]float64{{0,0},{b,0},{b/2.0,h},{0,0}}
		area, xc, yc, ixx, iyy, _, _, _, _ = SecPrp(ncs, wts, coords)
		j = 0.0216 * math.Pow(b,4)		
	}
	if rxx == 0.0{rxx = math.Sqrt(ixx/area)}
	if ryy == 0.0{ryy = math.Sqrt(iyy/area)}
	if j == 0.0{j = ixx + iyy}
	bar = Secprop{
		Area:      area,
		Perimeter: perimeter,
		Ixx:       ixx,
		Iyy:       iyy,
		Xc:        xc,
		Yc:        yc,
		Sxx:       sxx,
		Syy:       syy,
		Zxx:       zxx,
		Zyy:       zyy,
		J:         j,
		Rxx:       rxx,
		Ryy:       ryy,
		Dims:      dims,
		Iww:       iww,
	}
	return
}

//PropNpBm returns area, ix and iy for a non uniform beam/col/section
func PropNpBm(styp int, b, h float64, dims []float64) (ar, ix, iy float64) {
	switch SectionMap[styp]{
		case "rectangle":
		ar = b * h
		ix = b * math.Pow(h,3)/12.0
		iy = h * math.Pow(b,3)/12.0
		case "i","ieq":
		tf := dims[2]
		tw := dims[3]
		//uncomment this if b and h are given
		//b = dims[0]
		//h = dims[1]
		//hw := h - 2.0 * tf
		ar = 2.0 * b * tf + (h - 2.0 * tf) * tw
		//perimeter = 4.0 * b + 2.0 * h - 2.0 * tw
		ix = b * math.Pow(h,3)/12.0 - (b - tw) * math.Pow(h - 2.0 * tf,3)/12.0
		iy = tf * math.Pow(b, 3)/6.0 + (h - 2.0*tf) * math.Pow(tw, 3)/12.0
		case "haunch2f":
		//tf := dims[2]; tw := dims[3]
		//d1 := []float64{dims[0],h,dims[2],dims[3]}
		s := SecGen(styp, dims)
		ar = s.Prop.Area
		ix = s.Prop.Ixx
		iy = s.Prop.Iyy
		case "haunch3f":
		s := SecGen(styp, dims)
		ar = s.Prop.Area
		ix = s.Prop.Ixx
		iy = s.Prop.Iyy		
	}
	return
}

// //ComboSecProp calcs the properties of a compound section
// //specified by a list of styps, dims and z dist from axis
// func ComboSecProp(styps []int, dims []float64, zs []float64)(bar Secprop){
// 	return
// }
