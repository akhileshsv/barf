package barf

//Ref - Michael Laszlo - Computational geometry in c++

import (
	"math"
	"sort"
	"gonum.org/v1/gonum/mat"
)

//Pt2d is a 2d point struct (lurks all over this pc)
type Pt2d struct {
	X float64
	Y float64
}

func (p *Pt2d) SetTol(tol int){
	mulz := math.Pow(10, float64(tol))
	p.X = math.Round(p.X*mulz)/mulz
	p.Y = math.Round(p.Y*mulz)/mulz
}

func (p *Pt2d) InRect(p1, p2 Pt2d)(bool){
	return (p1.X <= p.X) && (p1.Y <= p.Y) && (p2.X >= p.X) && (p2.Y >= p.Y)
}

//OnEdge checks if p is on an edge of a polygon 
func (p *Pt2d) OnEdge(poly []Pt2d)(onedge bool){
	for i, p1 := range poly{
		p2 := poly[0]
		if i < len(poly) - 1{
			p2 = poly[i+1]
		}
		switch p.Lclass(p1, p2){
			case "origin","destination","between":
			onedge = true
			return
		}
	}
	return
}


//InPoly checks if p is inside a polygon
//https://stackoverflow.com/questions/217578/how-can-i-determine-whether-a-2d-point-is-within-a-polygon
//sec 7.4.2, rourke
//https://wrfranklin.org/Research/Short_Notes/pnpoly.html
func (p *Pt2d) InPoly(poly []Pt2d, b Pt2d)(inpoly bool){
	a := Pt2d{p.X, p.Y}
	//WHY DOES THIS WORK
	b.X = p.X
	intz := 0
	for i, p1 := range poly{
		p2 := poly[0]
		if i < len(poly) - 1{
			p2 = poly[i+1]
		}
		cls, _ := EdgeInt(a, b, p1, p2)
		switch cls{
			case "cross":
			intz += 1
		}
	}
	if intz % 2 == 1{
		inpoly = true
	}
	return
}


func (p1 *Pt2d) Add(p2 Pt2d) (p3 Pt2d){
	p3 = Pt2d{p1.X + p2.X, p1.Y + p2.Y}
	return
}

func (p1 *Pt2d) Sub(p2 Pt2d)(p3 Pt2d){
	p3 = Pt2d{p1.X - p2.X, p1.Y - p2.Y}
	return
}

func (p1 *Pt2d) Scale(v float64)(p2 Pt2d){
	p2 = Pt2d{p1.X * v, p1.Y * v}
	return
}

//Length returns the length/magnitude of a point/vector (page 78)
func (p *Pt2d) Length()(float64){
	return math.Sqrt(p.X * p.X + p.Y * p.Y) 
}

//Pequal checks for point equality
func Pequal(p1, p2 Pt2d)(bool){
	if p1.X == p2.X && p1.Y == p2.Y{
		return true
	}
	return false
}

//DotP (steve holt)!
func DotP(p1, p2 Pt2d)(float64){
	return p1.X * p2.X + p1.Y * p2.Y
}

//P3orient returns the orientation of three points - sec 4.2.4, pg. 75
func P3orient(p0, p1, p2 Pt2d)(ort int){
	a := Pt2d{p1.X - p0.X, p1.Y - p0.Y}
	b := Pt2d{p2.X - p0.X, p2.Y - p0.Y}
	tar := a.X * b.Y - b.X * a.Y
	switch{
		case tar == 0.0:
		//kollinear
		ort = 0
		case tar > 0.0:
		ort = 1
		case tar < 0.0:
		ort = -1
	}
	return
}

//Lclass returns the orientation of a point relative to a line section 4.2.5 pg 76
func (p *Pt2d) Lclass(p0, p1 Pt2d)(lcs string){
	p2 := Pt2d{p.X, p.Y}
	a := Pt2d{p1.X - p0.X, p1.Y - p0.Y}
	b := Pt2d{p2.X - p0.X, p2.Y - p0.Y}
	tar := a.X * b.Y - b.X * a.Y
	switch{
		case tar > 0.0:
		lcs = "left"
		case tar < 0.0:
		lcs = "right"
		case (a.X * b.X < 0.0) || (a.Y * b.Y < 0.0):
		lcs = "behind"
		case a.Length() < b.Length():
		lcs = "beyond"
		case Pequal(p0, p2):
		lcs = "orgin"
		case Pequal(p1, p2):
		lcs = "destination"
		default:
		lcs = "between"
	}
	return
}

//Rotedge rotates an edge from p1 (org) to p2 (dest) by 90 degrees cw about midpoint (sec 4.4.2)
func Rotedge(p1, p2 Pt2d)(p3, p4 Pt2d){
	a := p1.Add(p2)
	a = a.Scale(0.5)
	b := p2.Sub(p1)
	c := Pt2d{b.Y, -b.X}
	p3 = a.Sub(c.Scale(0.5))
	p4 = a.Add(c.Scale(0.5)) 
	return
}

//EdgeOverlap checks if two edges overlap 
func EdgeOverlap(a,b,c,d Pt2d)(bool){
	return (a.Lclass(c, d) == "between") || (b.Lclass(c,d) == "between") || (c.Lclass(a,b) == "between") || (d.Lclass(a,b) == "between")
}

//EdgeInt returns the point of intersection of two edges from (a,b) and (c,d) (pg. 93, sec. 4.4)
//switching from p1 etc (easier to copy)
func EdgeInt(a,b,c,d Pt2d)(cls string, px Pt2d){
	dc := d.Sub(c)
	cd := c.Sub(d)
	n := Pt2d{dc.Y, cd.X}
	denom := DotP(n, b.Sub(a))
	if denom == 0.0{
		lcs := a.Lclass(c,d)
		switch lcs{
			case "left", "right":
			cls = "parallel"
			default:
			cls = "collinear"
		}
		return
	}
	num := DotP(n, a.Sub(c))
	px = b.Sub(a)
	px = px.Scale(-num/denom)
	px = a.Add(px)
	cls = "cross"
	switch{
		case -num/denom > 1.0:
		cls = "skewf"
		case -num/denom < 0.0:
		cls = "skewb"
	}
	return
}

//Dist2d returns the 2d distance between two points
//barf has thousands of such functions.
func Dist2d(p1, p2 Pt2d)(dist float64){
	pd := p1.Sub(p2)
	dist = pd.Length()
	return
}

//Norm2d returns the unit normal between two Pt2ds
func Norm2d(p1, p2 Pt2d) (pn Pt2d){
	dist := Dist2d(p1, p2)
	if dist == 0{
		pn = Pt2d{p1.X, p1.Y}
		return
	}
	pn = Pt2d{(p2.X - p1.X)/dist, (p2.Y - p1.Y)/dist}
	return
}

//Lerp2d interpolates between two points
func Lerp2d(p1, p2 Pt2d, dist float64)(pe Pt2d){
	pn := Norm2d(p1, p2)
	ps := pn.Scale(dist)
	pe = p1.Add(ps)
	return
}

//Centroid2d returns the centroid of a []Pt2d
func Centroid2d(pts []Pt2d)(pc Pt2d){
	for _, pt := range pts{
		pc.X += pt.X
		pc.Y += pt.Y
	}
	np := float64(len(pts))
	pc.X = pc.X/np
	pc.Y = pc.Y/np
	return
}

//EdgeRot2d rotates p2 by ang about p1
func EdgeRot2d(ang float64, p1, p2 Pt2d) (p3 Pt2d){
	//convert ang to radians
	ang = ang * math.Pi/180.0
	//shift origin to v1
	tm1 := mat.NewDense(3,3, []float64{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		-p1.X, -p1.Y, 1.0,
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
		p1.X, p1.Y, 1.0,
	})
	//go forth and multiply
	a := mat.NewDense(1,3, []float64{p2.X,p2.Y,1.0})
	a.Mul(a, tm1)
	a.Mul(a, rm)
	a.Mul(a, tm2)
	p3.X = a.At(0,0)
	p3.Y = a.At(0,1)
	return
}

//EdgeOff2d offsets an edge from p1, p2 by dist
func EdgeOff2d(dist float64, p1, p2 Pt2d)(p3, p4 Pt2d){
	prot := EdgeRot2d(90.0, p1, p2)
	p3 = Lerp2d(p1, prot, dist)
	prot = EdgeRot2d(-90.0,p2, p1)
	p4 = Lerp2d(p2, prot, dist)
	return
}

//SortCw sorts a slice of points cw about xc, yc
//https://stackoverflow.com/questions/6989100/sort-points-in-clockwise-order
func SortCw(pts []Pt2d, pc Pt2d) {
	xc := pc.X; yc := pc.Y
	sort.SliceStable(pts, func(i,j int) bool {
		if pts[i].X-xc >= 0 && pts[j].X-xc < 0 {
			//return false
			return true
		}
		if pts[i].X-xc < 0 && pts[j].X-xc >= 0 {
			//return true
			return false
		}
		if pts[i].X-xc == 0 && pts[j].X-xc == 0 {
			if pts[i].Y-yc >= 0 || pts[j].Y-yc >= 0 {
				//return pts[i][1] > pts[j][1]
				return pts[i].Y > pts[j].Y
			}
			//return pts[j][1] > pts[i][1]
			return pts[j].Y > pts[i].Y
		}
		det := (pts[i].X-xc)*(pts[j].Y-yc) - (pts[j].X-xc)*(pts[i].Y-yc)
		if det < 0 {
			return true
		}
		if det > 0 {
			return false
			
		}
		d1 := (pts[i].X-xc)*(pts[i].X-xc) + (pts[i].Y-yc)*(pts[i].Y-yc)
		d2 := (pts[j].X-xc)*(pts[j].X-xc) + (pts[j].Y-yc)*(pts[j].Y-yc)
		//return d1 > d2
		return d1 < d2
	})
}
