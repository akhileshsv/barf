package barf

import (
	"fmt"
	"math"
)

var (
	BltDias = []float64{5, 6, 8, 10, 12, 14, 16, 20, 22, 24, 30}
	//btw - ramchandra does dia + 3 mm for timber
	BltHoles = []float64{6, 7, 9, 11, 13, 15, 18, 22, 24, 26, 33}
	NailDias = []float64{10.0, 8.00, 6.30, 5.60, 5.00, 4.50, 4.00, 3.55, 3.15, 2.80, 2.50, 2.24, 2.00, 1.80, 1.60, 1.40, 1.25}
)

func bltholed(dia float64)(dhole float64){
	switch{
		case dia <= 14.0:
		dhole = dia + 1.0
		case dia <= 24.0:
		dhole = dia + 2.0
		default:
		dhole = dia + 3.0
	}
	return
}

//netsecchain returns the net sectional area of a chain bolt group
func netsecchain(nrow, lb, dia, bmem float64)(nsa float64){
	dhole := bltholed(dia)
	//dmem = length of blt lb
	nsa = (bmem - nrow * dhole) * lb
	return
}

//bltdims returns end dist., edge dist, pitch and gauge for lb/ dia and ltyp and wdtyp (1 - soft, 2 - hard wood)
func bltdims(frctyp, ltyp, wdtyp int, lb, dia float64)(dend, dedge, pitch, gauge float64){
	ldrat := lb/dia
	//pitch = 4d for any ltyp
	pitch = 4.0 * dia
	//gauge vs net sectional area?
	//will only influence if staggered?
	if ldrat < 3.0{
		gauge = 2.5 * dia
	} else {
		gauge = 5.0 * dia
	}
	switch ltyp{
		case 1:
		//parallel to grain
		switch frctyp{
			case 1:
			//kompressive
			dend = 4.0 * dia
			case 2:
			//tensile
			dend = 7.0 * dia
			if wdtyp == 2{
				dend = 5.0 * dia
			}
		}
		if ldrat < 6.0{
			dedge = 1.5 * dia
		} else {
			dedge = gauge/2.0
		}
		case 2, 3:
		//perp. to grain, angled
		dedge = 4.0 * dia
	}
	return
}

//nbltsimp returns the number of bolts required for a given ultimate (axial) load pu
//given ltyp, ptyp, pu, lb, dia, fc, ang
func nbltsimp(ltyp, ptyp int, pu, lb, dia, fc, ang float64)(nblt, pblt float64, err error){
	//does one need a ctyp for connection type?
	//first get kg, kpg
	var kg, kpg float64
	kg, kpg, err = bltfac(ltyp, ptyp, lb, dia)
	if err != nil{
		return
	}
	pblt = bltfbr(ltyp, lb, dia, kg, kpg, fc, ang)
	nblt = math.Ceil(pu/pblt)
	//fmt.Println("nblts -> ", nblt)
	
	return
}

//bltfbr returns the load carrying capcity of a bolt in bearing on timber
//given lb, dia, kg, kpg, ang
func bltfbr(ltyp int, lb, dia, kg, kpg, fc, ang float64)(pblt float64){
	var fblt float64
	//fmt.Println("in-",ltyp, kg, kpg)
	switch ltyp{
		case 1:
		fblt = kg * fc
		case 2:
		fblt = kpg * fc
		case 3:
		fblt = WdFcAng(kg*fc, kpg*fc, ang, false) 
	}
	pblt = fblt * lb * dia
	//fmt.Println("fblt->", fblt)
	//fmt.Println("pblt->", pblt, " N")
	return
}

//bltfac returns the bolt bearing stress factor kg (along grain), kpg (perp. to grain)
//given ltyp(1- par to grain, 2- perp to grain), bolt len, bolt dia and ptyp (plate type - 32, 33 - timber-metal, timber-timber)
func bltfac(ltyp, ptyp int, lb, dia float64)(kg, kpg float64, err error){
	//fmt.Println("inp->",ltyp, ptyp, lb, dia, lb/dia)
	k1 := 0.0
	k2 := 0.0
	k3 := 0.0
	d1 := 0.0
	ldrat := lb/dia
	//k1 - l/d ratio factor
	ls := []float64{1,1.5,2,2.5,3,3.5,4,4.5,5,5.5,6,6.5,7,7.5,8,8.5,9,9.5,10,10.5,11,11.5,12}
	r1s := []float64{100,100,100,100,100,100,96,90,80,72,65,58,52,46,40,36,34,32,30,0,0,0,0}
	r2s := []float64{100,96,88,80,72,66,60,56,52,49,46,43,40,39,38,36,34,33,31,31,30,30,28}
	for i, l := range ls{
		r1 := r1s[i]
		r2 := r2s[i]
		if k1 > 0.0 && k2 > 0.0{break}
		switch{
			case k1 == 0.0 && l == ldrat:
			k1 = r1
			case k1 == 0.0 && l > ldrat:
			if i == 0{
				k1 = r1
			} else {
				k1 = (ldrat - ls[i-1])*(r1 - r1s[i-1])/(l - ls[i-1]) + r1s[i-1]
			}
			case k2 == 0.0 && l == ldrat:
			k2 = r2
			case k2 == 0.0 && l > ldrat:
			if i == 0{
				k2 = r2
			} else {
				k2 = (ldrat - ls[i-1])*(r2 - r2s[i-1])/(l - ls[i-1]) + r2s[i-1]
			}
		} 
	}
	//fmt.Println("found ldrat - ", ldrat, " k1 - ", k1, " k2 ",k2)
	if ldrat > 12.0{
		k1 = 0.0
		k2 = 28.0 - 2.0 * (ldrat - 12.0)/0.5
		if k2 < 0.0{
			k2 = 0.0
		}
	}
	if k1 == 0.0 && k2 == 0.0{
		err = fmt.Errorf("invalid l/d ratio %.f for len - %.f dia %.f",ldrat, lb, dia)
		return
	}
	//get bolt dia factor d1
	//table 17 nbc
	bds := []float64{6, 10, 12, 16, 20, 22, 25}
	dfs := []float64{5.7, 3.6, 3.35, 3.15, 3.05, 3.0, 2.9}
	for i, bd := range bds{
		df := dfs[i]
		if d1 > 0.0{
			break
		}
		switch{
			case bd == dia:
			d1 = df
			case bd > dia:
			if i == 0{
				d1 = df
			} else {
				d1 = (df - dfs[i-1])*(dia - bd)/(bd - bds[i-1]) + dfs[i-1]
			}
			
		}
	}
	//reduce for par. to grain loading
	switch ptyp{
		case 32:
		//timber/metal plates
		k3 = 1.0
		case 33:
		//timber/timber plates
		//sp.33 doesn't have this???
		//WHAT IS RAMCHANDRA's SOURCE it works.figure later
		k3 = 0.8
		//k3 = 1.0
	}
	kg = k1 * k3/100.0
	kpg = k2 * d1/100.0
	return
}

