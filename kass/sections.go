package barf

import (
	"fmt"
	"math"
	"sort"
	"gonum.org/v1/gonum/mat"
)

var (
	//base section shapes
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
		12: "i",
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
		25: "built-i",//built up i section IS NOT EXISTS CHANGE CHANGE CHANGE
		26: "spaced",//(rect) spaced column
		27: "l2-ss",//double (eq/ueq) angles back to back (same side)
		28: "l2-os",//double (eq/ueq) angles on opp. sides
		29: "plate-i",//i section with cover plates
	}
)

//SectIn is a struct that holds section property fields
//see mosley spencer - general section properties calculation, chapter 4
type SectIn struct {
	Ncs        []int
	Wts        []float64
	Coords     [][]float64
	Solid      bool
	Styp       int
	Dims       []float64
	X          []float64 `json:",omitempty"`
	Y          []float64 `json:",omitempty"`
	Ym, Ymx    float64 `json:",omitempty"`
	Xm, Xmx    float64 `json:",omitempty"`
	Prop       Secprop `json:",omitempty"`
	Ds         [][]float64 `json:",omitempty"`
	Dbars      []float64 `json:",omitempty"`
	Barpts     [][]float64 `json:",omitempty"`
	Monoplot   string `json:",omitempty"`
	Txtplot    string `json:",omitempty"`
	Data       []string `json:",omitempty"`
}

//Secprop holds calculated section property fields (yeolde)
type Secprop struct {
	Area, Perimeter, Ixx, Iyy, Xc, Yc, Sxx, Syy, Zxx, Zyy, J, Rxx, Ryy float64
	Ixy, Iuu, Ivv, Ruu, Rvv, Pxangle                                   float64
	Qxx, Qyy, Vfx, Vfy                                                 float64
	Iww, Itt                                                           float64 
	Sectype                                                            string
	Sname                                                              string 
	Dims                                                               []float64
}

//SortCcw sorts points counter clockwise (about centroid xc, yc)
//actually is sortcw now? as long as area is +ve eh

func SortCcw(pts [][]float64, xc, yc float64) {
	sort.SliceStable(pts, func(i, j int) bool {
		if pts[i][0]-xc >= 0 && pts[j][0]-xc < 0 {
			return false
			//return true
		}
		if pts[i][0]-xc < 0 && pts[j][0]-xc >= 0 {
			return true
			//return false
		}
		if pts[i][0]-xc == 0 && pts[j][0]-xc == 0 {
			if pts[i][1]-yc >= 0 || pts[j][1]-yc >= 0 {
				//return pts[i][1] > pts[j][1]
				return pts[i][1] < pts[j][1]
			}
			//return pts[j][1] > pts[i][1]
			return pts[j][1] < pts[i][1]
		}
		det := (pts[i][0]-xc)*(pts[j][1]-yc) - (pts[j][0]-xc)*(pts[i][1]-yc)
		if det < 0 {
			//return true
			return false
		}
		if det > 0 {
			//return false
			return true
		}
		d1 := (pts[i][0]-xc)*(pts[i][0]-xc) + (pts[i][1]-yc)*(pts[i][1]-yc)
		d2 := (pts[j][0]-xc)*(pts[j][0]-xc) + (pts[j][1]-yc)*(pts[j][1]-yc)
		//return d1 > d2
		return d1 < d2
	})
}


//SecInit inits a section
func (sec *SectIn) SecInit() {
	//get xmax and ymax
	var n1 int
	var ym, ymx, xm, xmx float64
	//sec.Solid = true
	for idx, nc := range sec.Ncs {
		if sec.Wts[idx] < -1 {
			sec.Solid = false
		}
		if idx == 0 {
			n1 = 0
		} else {
			n1 = sec.Ncs[idx-1]
		}
		for i := range sec.Coords[n1 : n1+nc-1] {
			i = i + n1
			x1 := sec.Coords[i][0]
			y1 := sec.Coords[i][1]
			if ym > y1{
				ym = y1
			}
			if ymx < y1{
				ymx = y1
			}
			if xm > x1{
				xm = x1
			}
			if xmx < x1{
				xmx = y1
			}
		}
	}
	sec.Ym = ym; sec.Xm = xm
	sec.Ymx = ymx; sec.Xmx = xmx
}

//SortXY does what it says, sorts by x and then y
func SortXY(pts [][]float64){
	sort.Slice(pts, func(i, j int) bool {
		if pts[i][0] == pts[j][0] {
			return pts[i][1] < pts[j][1]
		}
		return pts[i][0] < pts[j][0]
	})
}

//SecPrp calculates section properties given coords and weights
//as seen in mosley spencer section 4.3
func SecPrp(ncs []int, wts []float64, coords [][]float64) (area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle float64) {
	//section prop calculation (mosley sec
	var n1 int
	var mx, my float64
	for idx, nc := range ncs {
		if nc == 0 {continue}
		if idx == 0 {
			n1 = 0
		} else {
			n1 = ncs[idx-1]
		}
		wt := wts[idx]
		var xa, xb, ya, yb, xci, yci, areai, ixxi, iyyi, ixyi float64
		for i := range coords[n1 : n1+nc-1] {
			i = i + n1
			xa = coords[i][0]
			ya = coords[i][1]
			xb = coords[i+1][0]
			yb = coords[i+1][1]
			areai += (xa - xb) * (ya + yb) / 2.0
			yci += (xa - xb) * (math.Pow(ya, 2) + ya*yb + math.Pow(yb, 2)) / 6.0
			xci += (yb - ya) * (math.Pow(xa, 2) + xa*xb + math.Pow(xb, 2)) / 6.0
			ixxi += (xa - xb) * (math.Pow(ya, 3) + math.Pow(ya, 2)*yb + math.Pow(yb, 2)*ya + math.Pow(yb, 3)) / 12.0
			iyyi += (yb - ya) * (math.Pow(xa, 3) + math.Pow(xa, 2)*xb + math.Pow(xb, 2)*xa + math.Pow(xb, 3)) / 12.0
			ixyi += (xa - xb) * (xa*(9.0*math.Pow(ya, 2)+6.0*ya*yb+3.0*math.Pow(yb, 2)) + xb*(3.0*math.Pow(yb, 2)+6.0*ya*yb+9.0*math.Pow(yb, 2))) / 72.0
		}
		xci = xci / areai
		yci = yci / areai
		areai = wt * areai
		mx += areai * yci
		my += areai * xci
		ixx += wt * ixxi
		iyy += wt * iyyi
		ixy += wt * ixyi
		area += areai
	}
	xc = my / area
	yc = mx / area
	ixx -= area * math.Pow(yc, 2)
	iyy -= area * math.Pow(xc, 2)
	ixy -= area * xc * yc
	t1 := (ixx + iyy) / 2.0
	t2 := math.Sqrt(math.Pow(ixx-iyy, 2.0)+4.0*math.Pow(ixy, 2)) / 2.0
	iuu = t1 + t2
	ivv = t1 - t2
	pxangle = math.Atan((ixx-iuu)/ixy) * 180.0 / math.Pi
	return
}

//SecArea (sigh) does mosley spencer section 4.3 but now with a SectIn as input
//all these multiple similar function functions help improve code clarity
func SecArea(sec *SectIn, allprp bool) (area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle float64) {
	//GET PERIMETER (add euc dist of all coords)
	var nstrt int
	var mx, my float64
	for idx, nc := range sec.Ncs {
		wt := sec.Wts[idx]
		var xa, xb, ya, yb, xci, yci, areai, ixxi, iyyi, ixyi float64
		//fmt.Println("idx, nc, nstrt, nprev",idx, nc, nstrt, nprev)
		for i := range sec.Coords[nstrt : nstrt+nc-1]{
			i = i + nstrt
			xa = sec.Coords[i][0]
			ya = sec.Coords[i][1]
			xb = sec.Coords[i+1][0]
			yb = sec.Coords[i+1][1]
			areai += (xa - xb) * (ya + yb) / 2.0
			yci += (xa - xb) * (math.Pow(ya, 2) + ya*yb + math.Pow(yb, 2)) / 6.0
			xci += (yb - ya) * (math.Pow(xa, 2) + xa*xb + math.Pow(xb, 2)) / 6.0
			ixxi += (xa - xb) * (math.Pow(ya, 3) + math.Pow(ya, 2)*yb + math.Pow(yb, 2)*ya + math.Pow(yb, 3)) / 12.0
			iyyi += (yb - ya) * (math.Pow(xa, 3) + math.Pow(xa, 2)*xb + math.Pow(xb, 2)*xa + math.Pow(xb, 3)) / 12.0
			ixyi += (xa - xb) * (xa*(9.0*math.Pow(ya, 2)+6.0*ya*yb+3.0*math.Pow(yb, 2)) + xb*(3.0*math.Pow(yb, 2)+6.0*ya*yb+9.0*math.Pow(yb, 2))) / 72.0	
		}
		nstrt += nc
		xci = xci / areai
		yci = yci / areai
		areai = wt * areai
		mx += areai * yci
		my += areai * xci
		ixx += wt * ixxi
		iyy += wt * iyyi
		ixy += wt * ixyi
		area += areai
	}
	xc = my / area
	yc = mx / area
	if allprp {
		ixx -= area * math.Pow(yc, 2)
		iyy -= area * math.Pow(xc, 2)
		ixy -= area * xc * yc
		t1 := (ixx + iyy) / 2.0
		t2 := math.Sqrt(math.Pow(ixx-iyy, 2.0)+4.0*math.Pow(ixy, 2)) / 2.0
		iuu = t1 + t2
		ivv = t1 - t2
		pxangle = math.Atan((ixx-iuu)/ixy) * 180.0 / math.Pi
	}
	return
}

//SecCalc calcs section properties (calls SecArea)
func (s *SectIn) SecCalc(){
	//calcs via SecArea
	area, xc, yc, ixx, iyy, _, _, _, _ := SecArea(s, true)
	p := s.OutBound()
	s.Prop = Secprop{
		Area:area,
		Xc:xc,
		Yc:yc,
		Ixx:ixx,
		Iyy:iyy,
		Perimeter:p,
	}
	return
}

//SecTranslate translates a section by tx, ty
func SecTranslate(s SectIn, tx, ty float64) (st SectIn){
	coords := make([][]float64, len(s.Coords))
	wts := make([]float64, len(s.Wts))
	copy(wts, s.Wts)
	ncs := make([]int, len(s.Ncs))
	copy(ncs, s.Ncs)

	//translation matrix 
	tm := mat.NewDense(3,3, []float64{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		tx, ty, 1.0,
	})
	//loop thorugh coords and get new coords by mat mul
	for i, pt := range s.Coords{
		coords[i] = make([]float64, len(pt))
		a := mat.NewDense(1,3, []float64{pt[0],pt[1],1.0})
		a.Mul(a, tm)
		//copy to coords
		coords[i][0] = a.At(0,0)
		coords[i][1] = a.At(0,1)
	}
	st = SectIn{
		Coords:coords,
		Ncs:ncs,
		Wts:wts,
		Styp:s.Styp,
		Dims:s.Dims,
	}
	st.SecInit()
	st.SecCalc()
	return
}

//SecScale scales a section by sx, sy
func SecScale(s SectIn, sx, sy float64) (ss SectIn){
	coords := make([][]float64, len(s.Coords))
	wts := make([]float64, len(s.Wts))
	copy(wts, s.Wts)
	ncs := make([]int, len(s.Ncs))
	copy(ncs, s.Ncs)

	//scale matrix 
	sm := mat.NewDense(3,3, []float64{
		sx, 0.0, 0.0,
		0.0, sy, 0.0,
		0.0, 0.0,1.0,
	})

	for i, pt := range s.Coords{
		coords[i] = make([]float64, len(pt))
		a := mat.NewDense(1,3, []float64{pt[0],pt[1],1.0})
		a.Mul(a, sm)
		//copy to coords
		coords[i][0] = a.At(0,0)
		coords[i][1] = a.At(0,1)
	}
	ss = SectIn{
		Coords:coords,
		Ncs:ncs,
		Wts:wts,
		Styp:s.Styp,
		Dims:s.Dims,
	}
	ss.SecInit()
	ss.SecCalc()
	return
}

//SecRotate rotates a section by ang degrees (anticlockwise)
func SecRotate(s SectIn, ang float64) (sr SectIn){
	coords := make([][]float64, len(s.Coords))
	wts := make([]float64, len(s.Wts))
	copy(wts, s.Wts)
	ncs := make([]int, len(s.Ncs))
	copy(ncs, s.Ncs)

	//convert ang to radians
	ang = ang * math.Pi/180.0
	//rotation matrix 
	rm := mat.NewDense(3,3, []float64{
		math.Cos(ang), math.Sin(ang), 0.0,
		-math.Sin(ang), math.Cos(ang), 0.0,
		0.0, 0.0, 1.0,
	})
	//loop thorugh coords and get new coords by mat mul
	for i, pt := range s.Coords{
		coords[i] = make([]float64, len(pt))
		a := mat.NewDense(1,3, []float64{pt[0],pt[1],1.0})
		a.Mul(a, rm)
		//copy to coords
		coords[i][0] = a.At(0,0)
		coords[i][1] = a.At(0,1)
	}
	sr = SectIn{
		Coords:coords,
		Ncs:ncs,
		Wts:wts,
		Styp:s.Styp,
		Dims:s.Dims,
	}
	sr.SecInit()
	sr.SecCalc()
	return
}

//SecOffset offsets a section either outward (ccw +1) or inward (ccw -1)
//https://stackoverflow.com/questions/68104969/offset-a-parallel-line-to-a-given-line-python/68109283#68109283
//s1.Coords = make([][]float64, len(s.Coords))
func SecOffset(s SectIn, offst, ccw float64) (s1 SectIn){
	s1.Wts = make([]float64, len(s.Wts))
	copy(s1.Wts, s.Wts)
	s1.Ncs = make([]int, len(s.Ncs))
	copy(s1.Ncs, s.Ncs)
	//vc := []float64{s.Prop.Xc, s.Prop.Yc}
	var n1 int
	for idx, nc := range s.Ncs {
		//fmt.Println("idx, nc",idx, nc)
		if s.Wts[idx] < -1 {
			s1.Solid = false
		}
		s1.Wts[idx] = s.Wts[idx]
		if idx == 0 {
			n1 = 0
		} else {
			n1 = s.Ncs[idx-1]
		}
		pts := s.Coords[n1 : n1+nc]
		n2 := nc -1
		for i, pt := range pts{
			//j - prev, k - next vertex
			//fmt.Println(ColorCyan,i,pt,ColorReset)
			j := (i + n2 - 1)%n2
			k := (i + 1)%n2
			//fmt.Println(ColorRed, i, j, k, ColorReset)
			p1 := s.Coords[n1+j]; p2 := s.Coords[n1+k]
			//if i < len(pts)-1{
			vn1 := Normvec2d(pt, p1)
			vn2 := Normvec2d(p2, pt)
			bisx := ccw * (vn1[0] + vn2[0])
			bisy := ccw * (vn1[1] + vn2[1])
			var bvec []float64
			switch ccw{
				case 1.0:
				bvec, _ = Unitvec([]float64{0,0},[]float64{bisx,bisy})
				case -1.0:
				bvec, _ = Unitvec([]float64{bisx,bisy},[]float64{0,0})
			}
			dp := DotPvec(vn1, vn2)
			blen := offst/(math.Sqrt(1.0 + dp)/2.0)
			px := pt[0] + blen * bvec[0]
			py := pt[1] + blen * bvec[1]
			s1.Coords = append(s1.Coords, []float64{px,py})
		}
	}
	s1.Styp = s.Styp
	if s1.Styp == 0{s1.Styp = -1}
	s1.SecInit()
	s1.SecCalc()
	return
}


//FlipX rotates by 90 about origin
//then translate up by ymax
//change s.Dims (if it exists)
func FlipX(s SectIn) (sf SectIn){
	ang := 270.0
	sf = SecRotate(s, ang)
	sf = SecTranslate(sf, 0, math.Abs(sf.Ymx - sf.Ym))
	if s.Dims != nil{
		switch s.Styp{
			case 0:
			//lol wut
			case 1:
			//b = d, d = b
			b := s.Dims[1]; d := s.Dims[0]
			sf.Dims = []float64{b,d}
			case 2:
			sf.Dims = s.Dims
			case 3:
			b := s.Dims[1]; h := s.Dims[0]
			sf.Dims = []float64{b,h}
			case 4:
			D := s.Dims[0]; B := s.Dims[1]; d := s.Dims[2]; b := s.Dims[3]
			sf.Dims = []float64{B, D, b, d}
			default:
			sf.Dims = s.Dims
		}
	}
	return
}

//VecAng finds the angle between two vectors from v1, v2 and v3, v4
func VecAng(v1, v2, v3, v4 []float64)(ang float64){
	vu1, _ := Unitvec(v1, v2)
	vu2, _ := Unitvec(v3, v4)
	dotp := 0.0
	for i, val := range vu1{
		val2 := vu2[i]
		dotp += val * val2
	}
	//fmt.Println("dotp",dotp)
	ang = math.Acos(dotp)
	return
}

//Unitvec finds the unit vector between two points
func Unitvec(v1, v2 []float64) (vu []float64, vmod float64){
	vu = make([]float64, len(v1))
	for i := range v1{
		vu[i] = v2[i] - v1[i]
		vmod += math.Pow(v2[i] - v1[i],2)
	}
	vmod = math.Sqrt(vmod)
	if vmod == 0{return}
	for i := range vu{
		vu[i] = vu[i]/vmod
	}
	return
}

//Lerpvec linearly interpolates between v1 and v2 given a scale (0.0 at v1, 1.0 at v2)
func Lerpvec(scale float64, v1, v2 []float64) (v3 []float64){
	vu, vmod := Unitvec(v1,v2)
	v3 = make([]float64, len(v1))
	switch scale{
		case 0.0:
		copy(v3,v1)
		return
		case 1.0:
		copy(v3,v2)
		return
	}
	for i := range v1{
		v3[i] = v1[i] + scale * vmod * vu[i]
	}
	return

}


//LerpKz linearly interpolates between v1 and v2 given a scale (0.0 at v1, 1.0 at v2)
//with a constant z-value (from v1 for now)
func LerpKz(scale float64, v1, v2 []float64) (v3 []float64){
	z1 := v1[2]
	//; z2 := v2[2]
	vu, vmod := Unitvec(v1[:2],v2[:2])
	v3 = make([]float64, len(v1))
	switch scale{
		case 0.0:
		copy(v3,v1)
		return
		case 1.0:
		copy(v3,v2)
		return
	}
	for i := range v1[:2]{
		v3[i] = v1[i] + scale * vmod * vu[i]
	}
	v3[2] = z1
	return

}


//LerpKy linearly interpolates between v1 and v2 given a scale (0.0 at v1, 1.0 at v2)
//with a constant y-value (from v1 for now)
func LerpKy(scale float64, v1, v2 []float64) (v3 []float64){
	y1 := v1[1]
	//; z2 := v2[2]
	v1 = []float64{v1[0], v1[2]}
	
	v2 = []float64{v2[0], v2[2]}
	vu, vmod := Unitvec(v1,v2)
	v3 = make([]float64, len(v1))
	switch scale{
		case 0.0:
		copy(v3,v1)
		v3 = []float64{v3[0],y1, v3[1]}
		return
		case 1.0:
		copy(v3,v2)
		v3 = []float64{v3[0],y1, v3[1]}
		return
	}
	for i := range v1{
		v3[i] = v1[i] + scale * vmod * vu[i]
	}
	v3 = []float64{v3[0],y1, v3[1]}
	return

}

//Rotvec rotates v2 by ang about v1
//NOTE ALL THESE ARE 2D OPS - add switch len(v1) for matrixsees
func Rotvec(ang float64, v1, v2 []float64) (v3 []float64){
	//convert ang to radians
	v3 = make([]float64, len(v1))
	ang = ang * math.Pi/180.0
	//shift origin to v1
	tm1 := mat.NewDense(3,3, []float64{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		-v1[0], -v1[1], 1.0,
	})
	//rotate by ang
	rm := mat.NewDense(3,3, []float64{
		math.Cos(ang), math.Sin(ang), 0.0,
		-math.Sin(ang), math.Cos(ang), 0.0,
		0.0, 0.0, 1.0,
	})
	//shift origin to 0,0
	tm2 := mat.NewDense(3,3, []float64{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		v1[0], v1[1], 1.0,
	})
	//go forth and multiply
	a := mat.NewDense(1,3, []float64{v2[0],v2[1],1.0})
	a.Mul(a, tm1)
	a.Mul(a, rm)
	a.Mul(a, tm2)
	v3[0] = a.At(0,0)
	v3[1] = a.At(0,1)
	return
}


//Rotvec3d rotates v2 by ang (about axis ax (1-x,2-y,3-z)) about v1
func Rotvec3d(ax int, ang float64, v1, v2 []float64) (v3 []float64){
	//convert ang to radians
	v3 = make([]float64, len(v1))
	ang = ang * math.Pi/180.0
	//shift origin to v1
	tm1 := mat.NewDense(4,4, []float64{
		1.0, 0.0, 0.0,0.0,
		0.0, 1.0, 0.0,0.0,
		0.0, 0.0, 1.0,0.0,
		-v1[0], -v1[1], -v1[2],1.0,
	})
	rm := mat.NewDense(4, 4, nil)
	c0 := math.Cos(ang)
	s0 := math.Sin(ang)
	switch ax{
		case 1:
		//rotate about x
		rm = mat.NewDense(4,4,[]float64{
			1,0,0,0,
			0,c0,s0,0,
			0,-s0,c0,0,
			0,0,0,1,
		})
		case 2:
		//about y
		rm = mat.NewDense(4,4,[]float64{
			c0,0,-s0,0,
			0,1,0,0,
			s0,0,c0,0,
			0,0,0,1,
		})
		case 3:
		//about z
		rm = mat.NewDense(4,4,[]float64{
			c0,s0,0,0,
			-s0,c0,0,0,
			0,0,1,0,
			0,0,0,1,
		})
	}
	//shift origin to 0,0
	tm2 := mat.NewDense(4,4, []float64{
		1.0, 0.0, 0.0,0.0,
		0.0, 1.0, 0.0,0.0,
		0.0, 0.0, 1.0,0.0,
		v1[0], v1[1], v1[2],1.0,
	})
	//go forth and multiply
	a := mat.NewDense(1,4, []float64{v2[0],v2[1],v2[2],1.0})
	tm1.Mul(tm1,rm)
	tm1.Mul(tm1, tm2)
	a.Mul(a, tm1)
	v3[0] = a.At(0,0)
	v3[1] = a.At(0,1)
	v3[2] = a.At(0,2)
	return
}

//Lerp3d lerps v1 towards unit vec ivec by dist d
func Lerp3d(d float64, ivec, v1 []float64)([]float64){
	if d == 0.0{return v1}
	dx := d * ivec[0]
	dy := d * ivec[1]
	dz := d * ivec[2]
	return []float64{
		v1[0] + dx,
		v1[1] + dy,
		v1[2] + dz,
	}
}

//LocVec returns the local unit x/y/z axis vector of a mem with jb(v1), je(v2) and wng(angle of roll)
func LocVec(ax int, v1, v2, wng []float64)(ivec []float64){
	ivec = []float64{0,0,0}
	switch wng[0]{
		case 2:
		dx := v2[0] - v1[0]
		dy := v2[1] - v1[1]
		dz := v2[2] - v1[2]
		l := Dist3d(v1, v2)
		rxx := dx/l; rxy := dy/l; rxz := dz/l
		ang := wng[1] * math.Pi/180.0
		cx := math.Cos(ang)
		sx := math.Sin(ang)
		switch ax{
			case 1:
			ivec = []float64{rxx, rxy, rxz}
			case 2:
			dnm := math.Sqrt(rxx * rxx + rxz * rxz)
			ryx := (-rxx * rxy * cx - rxz * sx)/dnm
			ryy := dnm * cx
			ryz := (-rxy * rxz * cx + rxx * sx)/dnm
			ivec =  []float64{ryx, ryy, ryz}
			case 3:
			dnm := math.Sqrt(rxx * rxx + rxz * rxz)
			rzx := (rxx * rxy * sx - rxz * cx)/dnm
			rzy := -dnm * sx
			rzz := (rxy * rxz * sx + rxx * cx)/dnm
			ivec = []float64{rzx, rzy, rzz}	
		}
		case 1:
		//vertical members
		dx := v2[0] - v1[0]
		dy := v2[1] - v1[1]
		dz := v2[2] - v1[2]
		l := Dist3d(v1, v2)
		rxx := dx/l; rxy := dy/l; rxz := dz/l
		ang := wng[1] * math.Pi/180.0
		cx := math.Cos(ang)
		sx := math.Sin(ang)
		switch ax{
			case 1:
			ivec = []float64{rxx, rxy, rxz}
			case 2:
			ivec = []float64{-rxy * cx, 0.0, sx}
			case 3:
			ivec = []float64{rxy * sx, 0.0, cx}
		}
		case 0:
		switch ax{
			case 2:
			ivec = []float64{0.0, 1.0, 0.0}
			case 3:
			ivec = []float64{0.0, 0.0, 1.0}
			case 1:
			dx := v2[0] - v1[0]
			dy := v2[1] - v1[1]
			dz := v2[2] - v1[2]
			l := Dist3d(v1, v2)
			rxx := dx/l; rxy := dy/l; rxz := dz/l
			ivec = []float64{rxx, rxy, rxz}
		}
	}
	return
}

//RotZ rotates v1 by ang (in radians) about the z-axis
func RotZ(v1 []float64,ang float64) (v3 []float64){
	v3 = make([]float64, len(v1))
	//rotate by ang
	rm := mat.NewDense(3,3, []float64{
		math.Cos(ang), -math.Sin(ang), 0.0,
		math.Sin(ang), math.Cos(ang), 0.0,
		0.0, 0.0, 1.0,
	})
	//go forth and multiply
	a := mat.NewDense(1,3, []float64{v1[0],v1[1],1.0})
	a.Mul(a, rm)
	v3[0] = a.At(0,0)
	v3[1] = a.At(0,1)
	return
}

//RotY rotates v1 by ang (in degrees) about the y-axis
func RotY(v1 []float64, ang float64)(v3 []float64){
	return
}

//Rot3dKz rotates v1 by ang keeping z constant
func RotKz(ang float64, v1, v2 []float64)(v3 []float64){
	v3 = Rotvec(90, v1[:2],v2[:2])
	v3 = append(v3, v1[2])
	return
}

//OutBound calculates the outer perimeter of a section
func (s *SectIn) OutBound() (p float64){
	//(..behold ze outer bounds!) outer perimeter
	nc := s.Ncs[0]
	for i:=0; i < nc - 1; i++{
		p1 := s.Coords[i]; p2 := s.Coords[i+1]
		p += Dist2d(p1,p2)
	}
	return
}


//GetWidth gets the available width of a section at dy
//n - no. of intersecting pts, dx - [p2-p1, p3-p4]

func (s *SectIn) GetWidth(dy float64) (n int, dx []float64, pts [][]float64){	
	//get points of intersection of line y = ymax - dy with polygon sides
	ptmap := make(map[Pt]int)
	c1 := s.Ymx - dy
	var n1 int
	for idx, nc := range s.Ncs {
		if idx == 0 {
			n1 = 0
		} else {
			n1 = s.Ncs[idx-1]
		}
		for i := range s.Coords[n1 : n1+nc-1] {
			i = i + n1
			pta := Pt{s.Coords[i][0], s.Coords[i][1]}
			ptb := Pt{s.Coords[i+1][0], s.Coords[i+1][1]}
			if (pta.Y < c1 && ptb.Y < c1) {
				continue
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
					n++
				}
			}
		}
		if len(pts) == 0 {
			continue
		}
	}
	//sort and calc distances
	SortXY(pts)
	for i := 0; i < len(pts) - 1; i++{
		if i % 2 == 0{
			p1 := pts[i]; p2 := pts[i+1]
			dx = append(dx, Dist2d(p1,p2))
		}
	}
	return
}

//UpdateProp (again) calcs sec prop and updates vals
//use after rotate/transform ops
//THIS IS THE SAME AS SecCalc goddamn

func (s *SectIn) UpdateProp(){
	area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle := SecArea(s, true)
	s.Prop = Secprop{
		Area:area,
		Xc:xc,
		Yc:yc,
		Ixx:ixx,
		Iyy:iyy,
		Rxx:math.Sqrt(ixx/area),
		Ryy:math.Sqrt(iyy/area),
		J:ixx+iyy,
		Iuu:iuu,
		Ixy:ixy,
		Ivv:ivv,
		Ruu:math.Sqrt(iuu/area),
		Rvv:math.Sqrt(ivv/area),
		Pxangle:pxangle,
	}
	s.SecInit()
	return
}

//SplitSides gets the midpoint of all sides of a section
//and basically, splits sides and returns a new section
func (s *SectIn) SplitSides(tol float64) (s1 SectIn){
	//get midpoint of each side of section
	//add to coords, sort and calc
	//if tol > 0.0 {}
	s1.Wts = make([]float64, len(s.Wts))
	copy(s1.Wts, s.Wts)
	s1.Ncs = make([]int, len(s.Ncs))
	s1.Coords = [][]float64{}
	var n1 int
	
	for idx, nc := range s.Ncs {
		//fmt.Println("idx, nc",idx, nc)
		if s.Wts[idx] < -1 {
			s1.Solid = false
		}
		s1.Wts[idx] = s.Wts[idx]
		if idx == 0 {
			n1 = 0
		} else {
			n1 = s.Ncs[idx-1]
		}
		pts := s.Coords[n1 : n1+nc]
		for i := range pts{
			if i < len(pts)-1{
				p1 := s.Coords[i+n1]; p2 := s.Coords[i+n1+1]
				//fmt.Println("p1, p2",p1, p2)
				p3 := MidPt(p1,p2)
				switch tol{
					case 0.0:
					s1.Coords = append(s1.Coords, p1)
					s1.Coords = append(s1.Coords, p3)
					s1.Ncs[idx] += 2	
					default:
					if Dist2d(p1,p3) < tol{
						s1.Coords = append(s1.Coords, p1)
						s1.Ncs[idx] += 1
					} else {
						s1.Coords = append(s1.Coords, p1)
						s1.Coords = append(s1.Coords, p3)
						s1.Ncs[idx] += 2	
					}
				}
			} else {
				s1.Coords = append(s1.Coords, s.Coords[i+n1])
				s1.Ncs[idx]++
			}
		}
	}
	s1.Styp = s.Styp
	s1.Dims = s.Dims
	s1.SecInit()
	return
	
}


//Splitmax splits a section with max. tol length sides
//uses split edge 2d to get max possible points within tol dist of each other
func (s *SectIn) Splitmax(tol float64) (s1 SectIn){
	s1.Wts = make([]float64, len(s.Wts))
	copy(s1.Wts, s.Wts)
	s1.Ncs = make([]int, len(s.Ncs))
	s1.Coords = [][]float64{}
	var n1 int
	
	for idx, nc := range s.Ncs {
		//fmt.Println("idx, nc",idx, nc)
		if s.Wts[idx] < -1 {
			s1.Solid = false
		}
		s1.Wts[idx] = s.Wts[idx]
		if idx == 0 {
			n1 = 0
		} else {
			n1 = s.Ncs[idx-1]
		}
		pts := s.Coords[n1 : n1+nc]
		for i := range pts{
			if i < len(pts)-1{
				p1 := s.Coords[i+n1]; p2 := s.Coords[i+n1+1]
				vts := SplitEdge2d(p1, p2, tol)
				for j, vtx := range vts{
					if j != len(vts)-1{
						s1.Coords = append(s1.Coords, vtx)
						s1.Ncs[idx]++
					}
				}
			} else {
				s1.Coords = append(s1.Coords, s.Coords[i+n1])
				s1.Ncs[idx]++
			}
		}
	}
	s1.Styp = s.Styp
	s1.Dims = s.Dims
	s1.SecInit()
	return
	
}

//Tolchk checks p1 and p2 for tolerance (distance) in x and y (for bar placement basically)
func Tolchk(p1,p2 []float64,tol float64) bool{
	switch len(p1){
		case 2:
		if math.Abs(p1[0] - p2[0]) >= tol && math.Abs(p1[1] - p2[1]) >= tol{
			return true
		}
		return false
		case 3:
		
	}
	return false
}

//MidPt gets the mid point of p1 and p2
func MidPt(p1, p2 []float64) (p3 []float64){
	p3 = make([]float64, len(p1))
	for i := range p3{
		p3[i] = (p1[i] + p2[i])/2.0
	}
	return
}

//Trans2d translates coords by tx, ty
func Trans2d(p1 []float64, tx, ty float64) (pt []float64){
	tm := mat.NewDense(3,3, []float64{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		tx, ty, 1.0,
	})
	
	pt = make([]float64, len(p1))
	a := mat.NewDense(1,3, []float64{p1[0],p1[1],1.0})
	a.Mul(a, tm)
	//copy to coords
	pt[0] = a.At(0,0)
	pt[1] = a.At(0,1)
	return
}

//Trans3d translates coords by tx, ty, tz
func Trans3d(p1 []float64, tx, ty, tz float64) (pt []float64){	
	tm := mat.NewDense(4,4, []float64{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,		
		tx, ty, tz, 1.0,
	})
	
	pt = make([]float64, len(p1))
	a := mat.NewDense(1,3, []float64{p1[0],p1[1],p1[2],1.0})
	a.Mul(a, tm)
	//copy to coords
	pt[0] = a.At(0,0)
	pt[1] = a.At(0,1)
	pt[2] = a.At(0,2)
	return
}

//Draw3d is meant to draw a 3d view of a section
//it is brake (not werk)
func Draw3d(ms [][]int, coords [][]float64, ss []SectIn) (data string){
	for _, m := range ms{
		jb := m[0]; je := m[1]; cp := m[2]; pl := m[3]
		s := ss[cp-1]
		lz := Dist3d(coords[jb-1],coords[je-1])
		data += s.Dat3d(pl, coords[jb-1], lz)
		data += "\n"
	}
	return
}

//Dat3d prints 3d plot data for a set of edges
//it is not werk
func (s *SectIn) Dat3d(pl int, p0 []float64, lz float64) (data string){
	//
	var p1, p2 [][]float64
	switch pl{
		case 1:
		//xy plane	
		for _, pt := range s.Coords{
			//fmt.Println(i, pt)
			data += fmt.Sprintf("%f %f %f\n",p0[0]+pt[0],p0[1]+pt[1],p0[2]+0.0)
			p1 = append(p1, []float64{p0[0]+pt[0],p0[1]+pt[1],p0[2]+0.0})
		}
		data += "\n"
		for _, pt := range s.Coords{
			//fmt.Println(i, pt)
			data += fmt.Sprintf("%f %f %f\n",p0[0]+pt[0],p0[1]+pt[1],p0[2]+lz)
			p2 = append(p2, []float64{p0[0]+pt[0],p0[1]+pt[1],p0[2]+0.0})
		}
		case 2:
		//yz plane
		for _, pt := range s.Coords{
			//fmt.Println(i, pt)
			data += fmt.Sprintf("%f %f %f\n",p0[0]+0.0,p0[1]+pt[0],p0[2]+pt[1])
			p1 = append(p1, []float64{p0[0]+0.0,p0[1]+pt[0],p0[2]+pt[1]})
		}
		data += "\n"
		for _, pt := range s.Coords{
			//fmt.Println(i, pt)
			data += fmt.Sprintf("%f %f %f\n",p0[0]+lz,p0[1]+pt[0],p0[2]+pt[1])
			p2 = append(p2, []float64{p0[0]+lz,p0[1]+pt[0],p0[2]+pt[1]})
		}
		case 3:
		//xz plane
		for _, pt := range s.Coords{
			//fmt.Println(i, pt)
			data += fmt.Sprintf("%f %f %f\n",p0[0]+pt[0],p0[1]+0.0,p0[2]+pt[1])
			p1 = append(p1, []float64{p0[0]+pt[0],p0[1]+0.0,p0[2]+pt[1]})
		}
		data += "\n"
		for _, pt := range s.Coords{
			//fmt.Println(i, pt)
			data += fmt.Sprintf("%f %f %f\n",p0[0]+pt[0],p0[1]+lz,p0[2]+pt[1])
			p2 = append(p2, []float64{p0[0]+pt[0],p0[1]+lz,p0[2]+pt[1]})
		}
	}
	//data += "\n"
	for i, x := range p1{
		y := p2[i]
		if 1 == 0{
			data += fmt.Sprintf("%f %f %f\n",x[0],x[1],x[2])
			data += fmt.Sprintf("%f %f %f\n",y[0],y[1],y[2])
			data += "\n"
		}
	}
	return
} 

//SecOffsetyc offsets section using the centroid as origin
//ye olde - doesn't work for sections w/o centroid symmetry
func SecOffsetYc(s SectIn, offst float64) (st SectIn){
	coords := make([][]float64, len(s.Coords))
	wts := make([]float64, len(s.Wts))
	copy(wts, s.Wts)
	ncs := make([]int, len(s.Ncs))
	copy(ncs, s.Ncs)
	vc := []float64{s.Prop.Xc, s.Prop.Yc}
	
	for i, vi := range s.Coords{
		dist := Dist2d(vi, vc)
		scale := offst/dist
		vn := Lerpvec(scale, vi, vc)
		coords[i] = []float64{math.Round(vn[0]), math.Round(vn[1])}
	}
	st = SectIn{
		Coords:coords,
		Ncs:ncs,
		Wts:wts,
		Styp:s.Styp,
		Dims:s.Dims,
	}
	st.SecInit()
	st.SecCalc()
	return
}
