package barf

import (
	"math"
	"fmt"
	"gonum.org/v1/gonum/mat"
	/*
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"runtime"
	"path/filepath"
	"math"
	"math/rand"
	//draw"barf/draw" FORGET DRAW MAN
	//
	*/
)

//aha new file feels fine

//TODO
func ViewDat(mod *Model){
	
}

func npsecdraw(dims []float64, ts, ds []float64){}
/*
func getd(dims []float64, styp int) (d float64){
	switch styp{
		case 0:
		d := dims[0]
		case 1:
		//rect
		b := dims[0]; d := dims[1]
		ymax = d
		case 2:
		//e tri
		b := dims[0]
		h := b * math.Tan(math.Pi/3.0)/2.0
		case 3:
		//r tri
		b := dims[0]
		h := dims[1]
		case 4:
		//box section aha
		B := dims[0]; D := dims[1]; b := dims[2]; d := dims[3]
		case 5:
		//TUBE HOW
		//eh copy secprop for now
		D := dims[0]
		case 6:
		//t- section
		bf := dims[0]; d := dims[1]; bw := dims[2]; df := dims[3]
		case 7:
		//bottom left origin l
		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		case 8:
		//l - left
		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		case 9:
		//l- right
		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		case 10:
		//l - eq
		b := dims[0]; d := dims[1]; t := dims[2]
		case 11:
		//plus
		b := dims[0]; d := dims[1]; s := dims[2]; t := dims[3]
		case 12:
		//equal i section I
		b := dims[0]
		h := dims[1]
		tf := dims[2]
		tw := dims[3]
		case 13:
		//c ( [ ) section
		b := dims[0]
		h := dims[1]
		tf := dims[2]
		tw := dims[3]
		case 14:
		//t pocket
		bf := dims[0]
		d := dims[1]
		bw := dims[2]
		df := dims[3]
		case 15:
		//pentagon
		b := dims[0]
		case 16:
		//"house"
		b := dims[0]
		case 17:
		//hexagon
		b := dims[0]
		return
		case 18:
		//octagon
		b := dims[0]
		case 19:
		//tapered pocket section (allen 5.2)
		bf := dims[0]; d := dims[1]; bw := dims[2]; df := dims[3]; bp := dims[4]
		case 20:
		//trapezoidal section (subramanian 5.7)
		//slopes over d from bw (bottom) to bf(top)
		bf := dims[0]; d := dims[1]; bw := dims[2]
		case 21:
		//(square) diamond section 
		b := dims[0]
		case 22:
		//tapered t section
		bf := dims[0]; d := dims[1]; T := dims[2]; t := dims[3]; df := dims[4]
	}
	return

}
*/

func secview2d(pb, pe, dims []float64, styp int) (pts [][]float64){
	//d, p0s := getd(dims, styp)
	return
}

//make this plot view?
func PlotEle2d(mod *Model, term string, pltchn chan string){
	//plot 2d mod elevation
	
}

func (s *SectIn) Draw3d(p1, p2, wng []float64){
	//draw a section in 3d between p1, p2
	var rxx, rxy, rxz, ryx, ryy, ryz, rzx, rzy, rzz, rden float64
	var rarr []float64
	wtyp := wng[0]
	wang := wng[1] * math.Pi / 180.0
	xb, yb, zb := p2[0], p2[1], p2[2]
	xa, ya, za := p1[0], p1[1], p1[2]
	l := math.Sqrt(math.Pow(xb-xa, 2) + math.Pow(yb-ya, 2) + math.Pow(zb-za, 2))
	rxx = (xb - xa) / l
	rxy = (yb - ya) / l
	rxz = (zb - za) / l
	rden = math.Sqrt(math.Pow(rxx, 2) + math.Pow(rxz, 2))
	switch int(wtyp) {
	case 2: //general orientation
		ryx = (-rxx*rxy*math.Cos(wang) - rxz*math.Sin(wang)) / rden
		ryy = rden * math.Cos(wang)
		ryz = (-rxy*rxz*math.Cos(wang) + rxx*math.Sin(wang)) / rden
		rzx = (rxx*rxy*math.Sin(wang) - rxz*math.Cos(wang)) / rden
		rzy = -rden * math.Sin(wang)
		rzz = (rxy*rxz*math.Sin(wang) + rxx*math.Cos(wang)) / rden
		
	case 1: //vertical members
		rxx = 0
		//rxy = rxy
		rxz = 0
		ryx = -rxy * math.Cos(wang)
		ryy = 0
		ryz = math.Sin(wang)
		rzx = rxy * math.Sin(wang)
		rzy = 0
		rzz = math.Cos(wang)
	case 0:
	}
	if wtyp == 1 || wtyp == 2 {
		rarr = []float64{
			rxx, rxy, rxz,
			ryx, ryy, ryz,
			rzx, rzy, rzz,
		}
	} else { //horizontal mem- identity matrix
		rarr = []float64{
			1,0,0,
			0,1,0,
			0,0,1,
		}
	}
	rmat := mat.NewDense(3, 3, rarr)
	p1m := mat.NewDense(1,3, p1)
	p2m := mat.NewDense(1,3, p2)
	p1m.Mul(p1m, rmat)
	p2m.Mul(p2m, rmat)
	fmt.Println(p1m)
	fmt.Println(p2m)
}
