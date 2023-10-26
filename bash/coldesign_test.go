package barf

import (
	"log"
	"testing"
)

func TestColDBs(t *testing.T){
	//h1, h2, lx, ly, tx, ty, mx, my, vx, vy, pn float64, grd, sectyp int
	var c Col
	//var h1, h2, lx, ly, tx, ty, mx, my, vx, vy, pn, pfac float64
	//var grd, sectyp, nsecs int
	
	c.H1 = 3.5; c.H2 = 4.0; c.Lx = 3.5; c.Ly = 3.5; c.Tx = 1.0; c.Ty = 1.0; c.Mx = 0.0; c.My = 0.0; c.Vx = 120.0; c.Vy = 40.0
	c.Pu = 1000.0; c.Pfac = 1.0
	c.Grd = 43; c.Styp = 8; c.Nsecs = 5
	c.Yeolde = true
	c.Spam = true
	c.Verbose = true
	err := ColDBs(&c)
	log.Println(err)
	c.Table(true)
}

func TestColCBs(t *testing.T){
	//h1, h2, lx, ly, tx, ty, mx, my, vx, vy, pn float64, grd, sectyp int
	//var h1, h2, lx, ly, tx, ty, mx, my, vx, vy, pn, pfac float64
	//var grd, sectyp, secdx int
	var c Col
	c.H1 = 3.5; c.H2 = 4.0; c.Lx = 3.5; c.Ly = 3.5; c.Tx = 1.0; c.Ty = 1.0; c.Mx = 0.0; c.My = 0.0; c.Vx = 120.0; c.Vy = 40.0
	c.Pu = 1000.0; c.Pfac = 1.0
	c.Grd = 43; c.Styp = 8; c.Sdx = 23
	fp, ok := ColCBs(&c)
	log.Println(fp, ok)
	c.H1 = 0.0; c.H2 = 0.0; c.Lx = 2.4; c.Ly = 4.8; c.Tx = 1.0; c.Ty = 0.9; c.Mx = 120.0; c.My = 45.0; c.Vx = 0.0; c.Vy = 0.0
	c.Pu = 500.0; c.Pfac = 1.25
	c.Grd = 43; c.Styp = 7; c.Sdx = 21
	fp, ok = ColCBs(&c)
	log.Println(fp, ok)
}


func TestVec(t *testing.T){
	/*
	vec := PbcBs(1,43)
	for i, v := range vec{
		log.Println(i*5,"->",v)
	}
	*/
	dt := 17.0; s := 52.0; sectyp := 0; grd := 43
	log.Println(PbcLerp(sectyp, grd, s, dt))
	dt = 41.785; s = 152.1; sectyp = 1; grd = 43
	log.Println(PbcLerp(sectyp, grd, s, dt))
	s = 155.543; dt = 15.83
	log.Println("yeolde->",PbcYeolde(s, dt))
	s = 104.543; dt = 25.83
	log.Println("yeolde->",PbcYeolde(s, dt))
}
