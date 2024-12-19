package barf

import (
	"os"
	"fmt"
	"log"
	"testing"
	"path/filepath"
)

func TestColDBs(t *testing.T){
	//h1, h2, lx, ly, tx, ty, mx, my, vx, vy, pn float64, grd, sectyp int
	var c Col
	//var h1, h2, lx, ly, tx, ty, mx, my, vx, vy, pn, pfac float64
	//var grd, sectyp, nsecs int
	
	c.H1 = 3.5; c.H2 = 4.0; c.Lx = 3.5; c.Ly = 3.5; c.Tx = 1.0; c.Ty = 1.0; c.Mx = 0.0; c.My = 0.0; c.Vx = 120.0; c.Vy = 40.0
	c.Pu = 1000.0; c.Pfac = 1.0
	c.Grd = 43; c.Sname = "uc"; c.Nsecs = 5
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
	c.Grd = 43; c.Sname = "uc"; c.Sdx = 23
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
	dt := 17.0; s := 52.0; sname := "ub"; grd := 43
	log.Println(PbcLerp(sname, grd, s, dt))
	dt = 41.785; s = 152.1; sname = "uc"; grd = 43
	log.Println(PbcLerp(sname, grd, s, dt))
	s = 155.543; dt = 15.83
	log.Println("yeolde->",PbcYeolde(s, dt))
	s = 104.543; dt = 25.83
	log.Println("yeolde->",PbcYeolde(s, dt))
}

func TestColDzIs(t *testing.T){
	var examples = []string{"duggal8.2","duggal8.4","duggal8.5",
		"duggal8.6","maity30.1","maity30.2",
		"maity32.1","maity32.2","maity32.3",
		"duggal8.9.1","duggal8.9.2","duggal8.9.3",
		"maity33.1","maity33.2",
	}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/col/")
	t.Log("starting axially loaded column design (is800) tests")
	for i, ex := range examples{
		if i < 6{continue}
		switch{
			case i <= 5:
			t.Log(ColorRed,"column i-section tests",ColorReset)
			default:
			t.Log(ColorGreen,"angle strut tests",ColorReset)
		}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		c, err := ReadCol(fname)
		if err != nil{
			t.Errorf("steel column design (is800) test failed")
		}
		err = ColDzIs(&c)
		if err != nil{
			t.Fatal(err)
		}
		rezstring += fmt.Sprintf("example %s\n",ex)
		for _, ss := range c.Ssecs{
			rezstring += fmt.Sprintf("%s fcd %.2f n/mm2 pu %.2f kn\t",ss.Sstr,ss.Fcd,ss.Fcd*ss.Area/1e3)
		}
		rezstring += "\n"
	}
	wantstring := ``
	if rezstring != wantstring{
		t.Errorf("steel column design (is800) failed")
		fmt.Println(rezstring)
	}
}

func TestColTIs(t *testing.T){	
	var examples = []string{"maity1"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/col")
	t.Log("starting tension member design tests")
	for i, ex := range examples{
		t.Logf("starting example no. %v %s\n",i+1,ex)
		fname := filepath.Join(datadir,ex+".json")
		c, err  := ReadCol(fname)
		if err != nil{
			t.Fatal("tension member design failed",err)
		}
		err = ColTCIs(&c)
		if err != nil{
			t.Fatal("tension member design failed",err)
		}
	}
}
