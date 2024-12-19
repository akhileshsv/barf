package barf

//Ref - Michael Laszlo - Computational geometry in c++

import (
	"fmt"
	"math"
)

//Pt is a 2d, floating point, point
type Pt struct {
	X float64
	Y float64
}

//Pt3d is a 3d, floating point point
type Pt3d struct {
	X float64
	Y float64
	Z float64
}

//Edge is an Edge defined by origin, dest
type Edge struct{
	Org, Dest Pt
}

//GetEdge returns an Edge struct from p1 (org) to p2(dest)
func GetEdge(p1, p2 []float64)(Edge){
	return Edge{
		Org:Pt{p1[0],p1[1]},
		Dest:Pt{p2[0],p2[1]},
	}
}

//Rot rotates an edge (by WHAT)
func (e *Edge) Rot()(Edge){
	m := AddPt(e.Org,e.Dest)
	m.Scale(0.5)
	v := SubPt(e.Dest,e.Org)
	n := Pt{v.Y, -v.X}
	o1 := n.Scale(0.5)
	o := SubPt(m, o1)
	d := AddPt(m,o1)
	return Edge{
		Org:o,
		Dest:d,
	}
}

//Dist2d returns the 2d distance between two points
func Dist2d(p1, p2 []float64) (float64){
	return math.Sqrt(math.Pow(p2[0]-p1[0],2)+math.Pow(p2[1]-p1[1],2))
}


//Dist3d is a more general 2d/3d distance between two points
//used for drawing labels and stuff lol
func Dist3d(p1, p2 []float64) (dist float64){
	switch len(p1){
		case 1:
		dist = math.Abs(p2[0] - p1[0])
		case 2:
		dist = math.Sqrt(math.Pow(p2[0]-p1[0],2)+math.Pow(p2[1]-p1[1],2))
		case 3:
		dist = math.Sqrt(math.Pow(p2[0]-p1[0],2)+math.Pow(p2[1]-p1[1],2)+math.Pow(p2[2]-p1[2],2))
	}
	return
}

//D2d returns the 2d distance, but now between two Pt structs
//wisdom is redundancy is strength
func D2d(p1, p2 Pt) (float64){
	p := SubPt(p2,p1)
	return p.Length()
}

//AddPt adds two point vectors
func AddPt(p1, p2 Pt)(Pt){
	return Pt{
		p1.X + p2.X,
		p1.Y + p2.Y,
	}
}

//SubPt subtracts two point vectors
func SubPt(p1, p2 Pt)(Pt){
	return Pt{
		p1.X - p2.X,
		p1.Y - p2.Y,
	}
}

//Length returns the length of a vector (from origin ofc)
func (p *Pt) Length() (float64){
	return math.Sqrt(p.X * p.X + p.Y * p.Y)
}

//Scale scales a vector by k
func (p *Pt) Scale(k float64)(Pt){
	return Pt{
		p.X * k,
		p.Y * k,
	}
}

//DotProd returns the dot product of two vectors
func DotProd(p1, p2 Pt) (float64){
	return p1.X * p2.X + p1.Y * p2.Y
}


//Orientation of p0, p1 and p2 - page 75 of ref.
//1 +ve, 0 collinear, -1 -ve
func Orientation(p0, p1, p2 Pt) (int){
	a := SubPt(p1,p0)
	b := SubPt(p2, p0)
	sa := a.X * b.Y - b.X * a.Y
	if sa > 0.0{
		return 1
	}
	if sa < 0.0{
		return -1
	}
	return 0
}

//Classify Pt p wrt edge formed by p0, p1 (see ref)
func (p *Pt) Classify(p0, p1 Pt) (int){
	a := SubPt(p1, p0)
	b := SubPt(*p, p0)
	sa := a.X * b.Y - b.X * a.Y
	if sa > 0.0{
		return 1
		//left
	}
	if sa < 0.0{
		return 2
		//right
	}
	if (a.X * b.X < 0.0) || (a.Y * b.Y < 0.0){
		return 3
		//behind
	}
	if a.Length() < b.Length(){
		return 4
		//beyond
	}
	if p0 == *p{
		return 5
		//origin
	}
	if p1 == *p{
		return 6
		//dest-
	}
	return 7
	//between
}

//Intersect returns the intersection point of two edges
func Intersect(e1, e2 Edge) (i int, t float64){
	a := e1.Org; b := e1.Dest
	c := e2.Org; d := e2.Dest
	n := Pt{SubPt(d,c).Y,SubPt(c,d).X}
	denom := DotProd(n, SubPt(b,a))
	if denom == 0.0{
		aclass := a.Classify(e2.Org, e2.Dest)
		fmt.Println("aclass-",aclass)
		if aclass == 1 || aclass == 2{
			//parallel
			i = 1
			return
		} else {
			i = 2
			//collinear
			return
		}
	}
	fmt.Println("denom",denom)
	num := DotProd(n, SubPt(a,c))
	fmt.Println("nom",num)
	t = -num/denom
	if  0.0 <= t && t <= 1.0{
		i = 3
		//skew cross
		return
	}
	i = 4
	//skew no cross
	return
}

//Normvec2d returns the unit normal to two vectors v1 v2 in two dimenshuns
func Normvec2d(v1, v2 []float64)(vn []float64){
	vu, _ := Unitvec(v1, v2)
	vn = make([]float64, 2)
	vn[0] = vu[1]; vn[1] = -vu[0]
	return
}

//DotPvec returns the dot product of v1 and v2 (2d vectors)
func DotPvec(v1, v2 []float64) (float64){
	return v1[0] * v2[0] + v1[1] * v2[1]
}

//EdgePlot plots a slice of edges
func EdgePlot(es []Edge){
	var data string
	for i, e := range es{
		data += fmt.Sprintf("%f %f %v\n",e.Org.X,e.Org.Y,i+1)
		data += fmt.Sprintf("%f %f %v\n",e.Dest.X,e.Dest.Y,i+1)
		data += "\n"
	}
	pltskript := "basic.gp"
	title := "edges yo"
	pltstr := skriptrun(data, pltskript, "dumb", title, "","")
	fmt.Println(pltstr)
	
}

//SplitEdge2d splits an edge between pts/verts v1, v2 in tol numbers
func SplitEdge2d(v1, v2 []float64, tol float64) (vs [][]float64){
	l := Dist3d(v1, v2)
	ndiv := math.Floor(l/tol)
	spc := l/ndiv
	nd := int(ndiv)
	//fmt.Println("l, ndiv, spc, nd",l, ndiv, spc, nd)
	for i := 0; i <= nd; i++{
		lseg := float64(i) * spc
		scale := lseg/l
		v3 := Lerpvec(scale, v1, v2)
		vs = append(vs, v3)
	} 
	return
}

//FindIntInf finds the intersection of the line segments from 2d pts 1, 2 to 3, 4
func FindIntInf(x1,y1,x2,y2,x3,y3,x4,y4 float64)(par bool, px, py float64){
	if (x1-x2)*(y3-y4)-(y1-y2)*(x3-x4) == 0.0{
		par = true
		return
	}
	px = ( (x1*y2-y1*x2)*(x3-x4)-(x1-x2)*(x3*y4-y3*x4) ) / ( (x1-x2)*(y3-y4)-(y1-y2)*(x3-x4) ) 
	py = ( (x1*y2-y1*x2)*(y3-y4)-(y1-y2)*(x3*y4-y3*x4) ) / ( (x1-x2)*(y3-y4)-(y1-y2)*(x3-x4) )
	return
}

//GenBmFunic generates a (1d bm) funicular for a loading given by ls and xs
func GenBmFunic(){}
