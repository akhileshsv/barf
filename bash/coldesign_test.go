package barf

import (
	"os"
	"fmt"
	"log"
	"testing"
	"path/filepath"
)


func TestBmColDzIs(t *testing.T){
	var examples = []string{"duggal10.1","duggal10.2","sub13.1","sub13.2"}
	dirname,_ := os.Getwd()
	var rezstring string
	datadir := filepath.Join(dirname,"../data/examples/bash/col/")
	t.Log("starting beam-column design (is800) tests")
	for i, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		c, err := ReadCol(fname)
		if err != nil{
			t.Errorf("steel column design (is800) test failed")
		}
		err = ColDz(&c)
		if err != nil{
			t.Fatal(err)
		}
		fmt.Println(c.Report)
		rezstring += fmt.Sprintf("example %s\n",ex)
		for _, ss := range c.Ssecs{
			rezstring += fmt.Sprintf("%s \t lfac %.3f ifaca %.3f ifacb %.3f\n",ss.Sstr,ss.Lfac,ss.Ifaca,ss.Ifacb)
		}
		rezstring += "\n"
	}
	wantstring := `example duggal10.1
ISHB250(2) 	 lfac 0.935 ifaca 0.000 ifacb 0.986
ISHB300(1) 	 lfac 0.869 ifaca 0.000 ifacb 0.892
ISHB300(2) 	 lfac 0.815 ifaca 0.000 ifacb 0.838

example duggal10.2
ISHB400(2) 	 lfac 0.183 ifaca 0.971 ifacb 0.661
ISHB450(1) 	 lfac 0.135 ifaca 0.886 ifacb 0.602
ISHB450(2) 	 lfac 0.132 ifaca 0.856 ifacb 0.575

example sub13.1
ISHB250(1) 	 lfac 0.660 ifaca 0.000 ifacb 0.635
ISMB350 	 lfac 0.633 ifaca 0.000 ifacb 0.525
ISHB250(2) 	 lfac 0.629 ifaca 0.000 ifacb 0.604

example sub13.2
W310x310x226 	 lfac 0.400 ifaca 0.979 ifacb 0.648
W310x310x253 	 lfac 0.302 ifaca 0.874 ifacb 0.577
W310x310x283 	 lfac 0.238 ifaca 0.776 ifacb 0.511

`
	if rezstring != wantstring{
		t.Errorf("steel column design (is800) failed")
		fmt.Println(rezstring)
	}
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
		if i > 2{continue}
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

//tension member/tie check tests
func TestColTie(t *testing.T){	
	var examples = []string{"maity22.1","maity22.2","maity23.1","maity23.2","maity23.13"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/col")
	t.Log("starting tension member design tests")
	for i, ex := range examples{
		if i != 4{continue}
		t.Logf("starting example no. %v %s\n",i+1,ex)
		fname := filepath.Join(datadir,ex+".json")
		c, err  := ReadCol(fname)
		if err != nil{
			t.Fatal("tension member design failed",err)
		}
		err = ColDz(&c)
		t.Log("err",err)
		if i == 4{
			for i, ss := range c.Ssecs{
				fmt.Println("i, ss-",i, ss.Sstr)
			}
		}
	}
}

func TestColDzBs(t *testing.T){
	var examples = []string{"mosley6.1","mosley6.2"}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/col/")
	t.Log("starting axially loaded column design (bs449) tests")
	for i, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		c, err := ReadCol(fname)
		if err != nil{
			t.Errorf("steel column design (bs449) test failed")
		}
		err = ColDzBs(&c)
		if err != nil{
			t.Log(err)
			t.Errorf("steel column design (bs449) test failed")
		}
		rezstring += fmt.Sprintf("example %s\n",ex)
		for _, ss := range c.Ssecs{
			rezstring += fmt.Sprintf("%s frat %.2f\t",ss.Sstr,ss.Frat)
		}
		rezstring += "\n"
	}
	wantstring := ``
	if rezstring != wantstring{
		t.Errorf("steel column design (bs449) failed")
		fmt.Println(rezstring)
	}	
}

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
	// c.H1 = 3.5; c.H2 = 4.0; c.Lx = 3.5; c.Ly = 3.5; c.Tx = 1.0; c.Ty = 1.0; c.Mx = 0.0; c.My = 0.0; c.Vx = 120.0; c.Vy = 40.0
	// c.Pu = 1000.0; c.Pfac = 1.0
	// c.Grd = 43; c.Sname = "uc"; c.Sdx = 23
	// fp, ok := ColCBs(&c)
	// log.Println(fp, ok)
	c.H1 = 0.0; c.H2 = 0.0; c.Lx = 2.4; c.Ly = 4.8; c.Tx = 1.0; c.Ty = 0.9; c.Mx = 120.0; c.My = 45.0; c.Vx = 0.0; c.Vy = 0.0
	c.Pu = 500.0; c.Pfac = 1.25
	c.Grd = 43; c.Sname = "ub"; c.Sdx = 21
	fp, ok := ColCBs(&c)
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

/*

func TestColDzArk(t *testing.T){
	var examples = []string{"rdarkc1","rdarkc2","rdarkc3","rdarkc4",
		"rdarkc5","rdarkc6","rdarkc7","rdarkc8",
		"rdarkc9","rdarkc10","rdarkc11","rdarkc12",
		"rdarkc13","rdarkc14","rdarkc15","rdarkc16",
		"rdarkc17","rdarkc18","rdarkc19","rdarkc20",
		"rdarkc21","rdarkc22","rdarkc23","rdarkc24"}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/col/")
	t.Log("starting rdark col tests")
	for i, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		c, _ := ReadCol(fname)
		c.Onism = true
		c.Nsecs = 3
		err = ColDzBs(&c)
		if err != nil{
			t.Log(err)
			t.Errorf("steel column design (bs449) test failed")
		}
		rezstring += fmt.Sprintf("example %s\n",ex)
		for _, ss := range c.Ssecs{
			rezstring += fmt.Sprintf("%s frat %.2f\t",ss.Sstr,ss.Frat)
		}
		rezstring += "\n"
	}
	wantstring := ``
	if rezstring != wantstring{
		fmt.Println(rezstring)
	}	
}

*/
