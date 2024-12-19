package barf

import (
	"fmt"
	"os"
	"math"
	"path/filepath"
	"runtime"
	"errors"
	
)

//where else does one put these basic funcs
func MinVal(val...float64)(mv float64){
	for i, v := range val{
		if i == 0{
			mv = v
		} else if v < mv{
			mv = v
		}
	}
	return
}
//MirrorX mirrors a 2d slice along x 
func MirrorX(vec [][]int)(mvec [][]int){
	width := len(vec[0])
	height := len(vec)
	mvec = make([][]int, height)
	for i := 0; i < height; i++{
		for j := 0; j < width; j++{
			mvec[i] = make([]int, width)
			mvec[i][j] = vec[i][width-j-1]
		}
	}
	return
}


//InitSects initalizes a model's sects slice (and cp vec)
func (mod *Model) InitSects()(err error){
	var styp int
	for i, dims := range mod.Dims{
		if len(mod.Sts) <= i{
			styp = mod.Sectyp
		} else {
			styp = mod.Sts[i]
		}
		s := SecGen(styp, dims)
		if s.Prop.Area <= 0.0{
			return fmt.Errorf("invalid section dims(%v) for styp(%v)",styp, dims)
		}
		mod.Sects = append(mod.Sects, s)
		switch mod.Type{
			case "cbm":
			switch mod.Same{
				case true: 
				if len(mod.Cp) == 0{
					mod.Cp = append(mod.Cp, []float64{s.Prop.Ixx, s.Prop.Area})
				}
				case false:
				mod.Cp = append(mod.Cp, []float64{s.Prop.Ixx, s.Prop.Area})
			}
		}
	}
	return
}
	
//SplitCoords splits the line between pb and pe into ndiv segments
func SplitCoords(pb, pe []float64, ndiv int)(coords [][]float64){
	p0 := make([]float64, len(pb))
	copy(p0, pb)
	coords = append(coords, p0)
	for i := 1; i <= ndiv; i++{
		p1 := Lerpvec(float64(i)/float64(ndiv), pb, pe)
		coords = append(coords, p1)
	}
	//fmt.Println("bet. pb-",pb, "and pe-",pe,"coords-\n",coords)
	return
}

//GetNpDepth returns the depth of a member at lx from start
func GetNpDepth(ts []int, ls, ds []float64, lx float64) (dx float64){
	var l1, l2, lspan float64
	l1 = ls[0]; l2 = ls[0] + ls[1]; lspan = ls[0] + ls[1] + ls[2]
	switch{
		case lx <= l1:
		switch ts[0]{
			case 0:
			//nothing
			dx = ds[0]
			case 1:
			//straight/flat
			case 2:
			//linear
			dx = ds[0] + (ds[1] - ds[0])*lx/l1
			case 3:
			//parabolic
		}
		case lx <= l2:
		switch ts[1]{
			case 0:
			dx = ds[1]
			case 1:
			case 2:
			case 3:
		}
		default:
		switch ts[2]{
			case 0:
			case 1:
			case 2:
			dx = ds[1] + (ds[2] - ds[1])*(lx-l2)/(lspan-l2)
			case 3:
		}
	}
	return
}

//AddMemLoad does what it says, adds a member load to a model
//one intends to use this for tweaks 
func (mod *Model) AddMemLoad(ld MLoad) (error){
	if ld.Mem < 1 || ld.Mem > len(mod.Mprp){
		return errors.New("non existent member index")
	}
	if ld.Ltyp < 1 || ld.Ltyp > 7{
		return errors.New("invalid load type")
	}
	ldvec := []float64{
		float64(ld.Mem),
		float64(ld.Ltyp),
		ld.Wa,
		ld.Wb,
		ld.La,
		ld.Lb,
		float64(ld.Lcon),
	}
	mod.Msloads = append(mod.Msloads, ldvec)
	return nil
}

//AddNodeLoad adds a nodal load to a model
func (mod *Model) AddNodeLoad(ld NLoad) (err error){
	if ld.Node < 1 || ld.Node > len(mod.Coords){
		return errors.New("non existent node index")
	}
	vec := []float64{float64(ld.Node)}
	vec = append(vec, ld.Vec...)
	mod.Jloads = append(mod.Jloads,vec) 
	return nil
}

//BmUvalCs contains (super useful) structx.com formulae for beams loaded by udl wul of eq spans lspans
//and uniform E*Iz ei
func BmUvalCs(wul, lspan, ei float64, nspans, endc int) (mul, vul, dmax float64){
	switch endc{
		case 0:
		mul = wul * math.Pow(lspan,2)/2.0
		vul = wul * lspan
		dmax = wul * math.Pow(lspan, 4.0)/8.0/ei
		case 1:
		mul = wul * math.Pow(lspan,2)/8.0
		vul = wul * lspan/2.0		
		dmax = 5.0 * wul * math.Pow(lspan, 4.0)/384.0/ei
		case 2:
		switch nspans{
			case 2:
			mul = wul * math.Pow(lspan,2)/8.0
			vul = 5.0 * wul * lspan/8.0
			dmax = wul * math.Pow(lspan, 4.0)/185.0/ei
			case 3:
			mul = 0.1 * wul * math.Pow(lspan,2)
			vul = 0.6 * wul * lspan
			dmax = 0.0069 * wul * math.Pow(lspan, 4.0)/ei
			default:
			mul = 0.1071 * wul * math.Pow(lspan,2)
			vul = 0.607 * wul * lspan
			dmax = 0.0065 * wul * math.Pow(lspan, 4.0)/ei
		}
	}
	return
}

//BmFrcX interpolate bm and shear at x
//written somewhere else too but again only the wise write redundant functions

func BmSfX(xs, bm, sf []float64, xr float64) (mr, vr float64){
	xdiv := xs[1] - xs[0]
	for i, vx := range sf{
		x := xs[i]
		if (math.Abs(x-xr) <= xdiv/5.0 || x >= xr){
			switch{
				case (x == xr || math.Abs(x-xr)<= xdiv/5.0) && vr + mr == 0.0:
				vr = vx
				mr = bm[i]
				//fmt.Println(ColorYellow,"AHA->x,xr,i,vr,mr",x, xr, i,vr,mr,ColorReset)
				return
				case x > xr && i == 0:
				//fmt.Println("i == 0->",x, xr, i)
				vr = vx
				mr = bm[i]
				return
				default:
				//BAAH. AVERAGE
				vr = (sf[i-1] + sf[i])/2.0
				mr = (bm[i-1]+ bm[i])/2.0
				return
			}
		}
	}
	return
}

//InitFolder creates a folder given a folder name, returns an error if it exists
func InitFolder(name, foldr string)(err error, fdir string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	if foldr == ""{foldr = "out"}
	switch foldr{
		case "out":
		fdir = filepath.Join(basepath,"../data/out",name)	
		case "web":
		fdir = filepath.Join(basepath,"../srvr/assets",name)
	}
	if err = os.Mkdir(fdir, 0755); os.IsExist(err) {
		err = errors.New("dir exists, change name")
	} else {
		err = nil
	}
	return
}


// //GetMemNpDepth returns the depth of a member at node/location n1
// func GetNpDepth(pb, pe, n1 []float64, lhnch, db, de float64, mdx, config int) (dy float64){
// 	lspan := Dist3d(pb, pe)
// 	lx := Dist3d(pb, n1)
// 	switch config{
// 		case 3:
// 		//uniformly tapered member
// 		dy = db + (de - db) * lx/lspan
// 		case 1,2:
// 		//haunched member
// 		var hstrt float64
// 		switch mdx{
// 			case 3:
// 			//left rafter has haunch at beginning
// 			if lx > lhnch{
// 				dy = de
// 			} else {
// 				dy = db + (de - db) * lx/lhnch
// 			}
// 			default:
// 			//all else have the haunch at the end
// 			hstrt = lspan - lhnch
// 			if lx > hstrt{
// 				dy = db + (de - db) * (lx - hstrt)/lhnch
// 			} else {
// 				dy = db 
// 			}
// 		}
// 	}
// 	return
// }


/*
   YEOLDE

   				//fmt.Println("FUBARRR->",x, xr, i)
				//xdiv := xs[i] - xs[i-1]
				/*
				if math.Abs(sf[i-1]) > math.Abs(vx){
					//sf decreases, usually left of span and positive sf
					vr = vx + math.Abs(sf[i-1]-vx)*(x - xr)/xdiv
					mr = bm[i] - math.Abs(vr+vx)*(x - xr)*0.5
					//fmt.Println("FUBARRR->",x, xr, i)
					fmt.Println(ColorYellow,"NEG SLOPE->x,xr",x, xr)
					fmt.Println("v1, v2, vr->",sf[i],sf[i-1],vr)
					fmt.Println("m1 m2 mr->",bm[i-1],bm[i],mr)
					fmt.Println(ColorReset)
				} else {
					//sf increases with x, on the  right
					vr = vx - math.Abs(vx - sf[i-1])*(x - xr)/xdiv
					mr = bm[i] + math.Abs(vr + vx)*(x - xr)*0.5
		
					fmt.Println(ColorRed,"POS SLOPE->x,xr",x, xr)
					fmt.Println("v1, v2, vr->",sf[i],sf[i-1],vr)
					fmt.Println("m1 m2 mr->",bm[i-1],bm[i],mr)
					fmt.Println(ColorReset)
				}
*/
				

