package barf

import (
	//"fmt"
	"os"
	"math"
	"path/filepath"
	"runtime"
	"errors"
	
)

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
				

