package barf

import (
	"fmt"
	"math"
	//opt"barf/opt"
)

var Bfsoils = []string{"gravel1","gravel0","sand1","sand0","siltsand","claysand","insilts","inclays"}

//table 16.3, subramanian
var Thmap = map[string]float64{
	"gravel1":40.0,
	"gravel0":38.0,
	"sand1":38.0,
	"sand0":34.0,
	"siltsand":34.0,
	"claysand":32.0,
	"insilts":33.0,
	"inclays":27.0,
}

type Rwall struct{
	Title              string
	Typ                int     //1 - cantilever, 2 - counterfort, 3 - ret. wall w/platform
	Bfill              int     //type of backfill (have 3-6 types, see kaveh)
	Code               int
	H, Thw             float64 //total height of wall
	Sbc, Mu            float64 //sbc, coeff. of friction
	Th, Bet            float64 //angle of friction of backfill, ang. of fric bet wall and backfill(coulomb's theory)
	Thb                float64 //angle of backfill with horizontal
	Pgsoil             float64 //unit weight of soil
	Fck, Fy            float64
	Pos                []float64
	Df                 float64 //depth of footing
	Psfs               []float64 //fs overturning, 
	Dslb               float64 //footing thickness
	B                  float64 //length of footing
	Tstm, Ltoe, Lheel  float64 //thickness of stem, length of toe slab, length of heel slab
	Wschrg             float64 //weight/surcharge load
}

func RwallOpt(r Rwall) (rrez Rwall, err error){
	//input params
	//read vec
	//calc depth of footing
	if r.Df == 0.0{
		err = r.setDf()
	}
	if err != nil{
		return
	}
	//fmt.Println("dfooting - >", r.Df)
	switch r.Typ{
		case 1:
	}
	return
}

//getndim returns the number of dims for rwall optimization
func (r *Rwall) getndim()(nd int){
	return
}

//setDf sets the req. depth of footing using Rankine's formula
func (r *Rwall) setDf()(err error){
	if r.Sbc == 0.0 || r.Pgsoil == 0.0{
		err = fmt.Errorf("invalid value(s) of sbc - %f soil unit wt - %f", r.Sbc, r.Pgsoil)
		return
	}
	//convert angle of friction th to rad
	th := r.Th * math.Pi/180.0
	sine := math.Sin(th)
	r.Df = math.Pow((1.0 - sine)/(1.0 + sine),2.0)
	r.Df = r.Sbc * r.Df/r.Pgsoil
	
	return
}
