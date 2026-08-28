package barf

import (
	"fmt"
	"math"
	"testing"
	//"github.com/go-gota/gota/dataframe"
)

func TestParFlgSec(t *testing.T){
	var ss StlSec
	ss,_ = GetStlSec("i",64,1)
	ss.Printz()
	ss,_ = GetStlSec("w",288,1)
	ss.Printz()
	ss,_ = GetStlSec("wpb",121,1)
	ss.Printz()
	ss, _ = GetStlSec("npb",69,1)
	ss.Printz()
	t.Log("sub ex. 13.2 section w 310x310x226")
	ss,_ = GetStlSec("w",228,1)
	ss.Printz()

}

func TestBoxSec(t *testing.T){
	df, err := GetStlDf("box")
	fmt.Println(df,err)
	ndx := StlSdxLims["box"]
	for idx := ndx; idx >= 0; idx--{
		ss, err := GetStlSec("box", idx, 1)
		if err != nil{
			t.Fatal(err)
		}
		fmt.Println("str-",ss.Sstr)
		fmt.Println("ixx, iyy, sxx, syy, zxx, zyy",ss.Ixx, ss.Iyy, ss.Sxx, ss.Syy, ss.Zxx, ss.Zyy)
	}
}

func TestBmColChk(t *testing.T){
	t.Log("starting beam-column check")
	t.Log("duggal ex. 10.1 ishb300 beam-column check")
	ss, _ := GetStlSec("i",20,1)
	ss.Pu = 1275.0e3
	ss.Vax = 150e3
	ss.Vbx = 150e3
	ss.Lspan = 3200.0
	err := ss.BmColChk800()
	if err != nil{
		t.Fatal("steel beam-column check failed")
	}
	t.Log("subramanian ex. 13.1 ishb250 beam-column check")
	ss, _ = GetStlSec("i",27,1)
	ss.Pu = 500e3
	ss.Max = 27e6
	ss.Mbx = 45e6
	ss.Lspan = 4000
	ss.Klx = 0.8
	ss.Kly = 0.8
	err = ss.BmColChk800()
	if err != nil{
		t.Fatal("steel beam-column check failed")
	}
	t.Log("duggal ex. 10.2 ishb450 beam-column check")
	ss, _ = GetStlSec("i",8,1)
	ss.Pu = 1035e3
	ss.Max = 0
	ss.Mbx = 81.25e6
	ss.May = 0
	ss.Mby = 22e6
	ss.Lspan = 3500
	err = ss.BmColChk800()
	if err != nil{
		t.Fatal("steel beam-column check failed")
	}
}

func TestPrlnChk(t *testing.T){
	t.Log("duggal ex 9.8 iswb150 check")
	pdl := 130.0
	pll := 0.0
	pwl := 2000.0
	theta := 30.0
	pspc := 1500.0
	tspc := 6000.0
	sname := "i"
	ss, err := GetStlSec(sname, 52, 1)
	if err != nil{
		t.Fatal("purlin check test failed")
	}
	p := (pdl + pll) * pspc/1000.0
	p += 100.0
	pn := p * math.Cos(theta * math.Pi/180.0)
	pp := p * math.Sin(theta * math.Pi/180.0)
	pn  += pwl * pspc/1000.0
	t.Log("unfactored pnormal P", pn, "pparallel H",pp)
	mux := 1.5 * pn * tspc * tspc/10.0/1000.0
	muy := 1.5 * pp * tspc * tspc/10.0/1000.0
	t.Log("mux, muy",mux, muy)
	ss.Ax = -2
	ss.Lspan = tspc
	ss.Mux = mux
	ss.Muy = muy
	ss.Lsb = true
	ss.Dmax = 5.0 * pn * math.Pow(tspc, 4)/ss.Em/ss.Ixx/384.0/1000.0
	t.Log("ss.Dmax-",ss.Dmax)
	ss.Vux = 0.6 * pn * tspc
	ss.Vuy = 0.6 * pp * tspc
	err = ss.BmChk800()
	t.Log(err)

}

func TestStlDf(t *testing.T){
	//fmt.Println(StlDfIs(1))
	for sname, styp := range StlStyps{
		df, _ := GetStlDf(sname)
		fmt.Println("styp,sname,len->",styp,sname,df.Nrow())
	}
	cp, err := GetStlCp("i", 1, 0, 1) 
	fmt.Println(cp, err)
	cp, err = GetStlCp("i", 1, 1, 2)
	fmt.Println(cp, err)
}

func TestStlSec(t *testing.T){
	t.Log("testing steel sec prop - section type ISMB300 duggal ex 9.1")
	sname := "i"
	sdx := 31
	code := 1
	ss, err := GetStlSec(sname, sdx, code)
	
	fmt.Println(err)
	
	ss.Printz()

	
	t.Log("testing steel sec prop - section type ISMC200 duggal ex 9.1")
	sname = "c"
	sdx = 21
	code = 1
	ss, err = GetStlSec(sname, sdx, code)
	
	fmt.Println(err)
	
	ss.Printz()

	t.Log("testing steel sec prop - section type ISMC350 duggal ex 9.1")
	sdx = 25

	ss, err = GetStlSec(sname, sdx, code)
	
	fmt.Println(err)
	
	ss.Printz()

	t.Log("testing steel sec prop - built-up i section duggal ex 9.1")
	b := 380.0
	h := 1640.0
	tf := 20.0
	tw := 15.0
	sdx = -1
	sname = "built-i"
	ss, err = GetStlSec(sname, sdx, code, b, h, tf, tw)
	
	
	fmt.Println(err)
	
	ss.Printz()
	
	//t.Fatal()
	
	t.Log("testing steel sec prop maity- section type unqeual 2xL")
	sname = "ln2-ss"
	sdx = 41
	code = 1
	ss, err = GetStlSec(sname, sdx, code, 10.0)
	fmt.Println(err)
	ss.Printz()

	t.Log("testing steel sec prop maity- section type unqeual 2xL")
	sname = "ln2-os"
	ss, err = GetStlSec(sname, sdx, code, 10.0)
	fmt.Println(err)
	ss.Printz()

	sdx = 0
	t.Log("testing steel tube section")
	sname = "tube"
	ss, err = GetStlSec(sname, sdx, code)
	fmt.Println(err)
	ss.Printz()

	t.Log("testing steel sec prop - section type ISA65656")
	sname = "l"
	sdx = 33
	code = 1
	ss, err = GetStlSec(sname, sdx, code)
	
	fmt.Println(err)
	ss.Printz()
	
	t.Log("testing steel sec prop - section type ISMC100")
	sname = "c"
	sdx = 17
	code = 1
	ss, err = GetStlSec(sname, sdx, code)
	
	fmt.Println(err)
	ss.Printz()
	
}

func TestCalcTu(t *testing.T){
	t.Log("starting steel section tensile strength calcs")
	t.Log("maity lec 22 ex. 1(a) bolted 2xisa75508 longer leg")
	sname := "ln2-ss"
	code := 1
	sdx := 27
	ss, err := GetStlSec(sname, sdx, code)
	if err != nil{
		t.Fatal(err)
	}
	ss.Bg.Dia = 18.0
	ss.Bg.Ni = 1
	ss.Bg.Nj = 4
	ss.Bg.Pitch = 50.0
	ss.Bg.Edged = 30.0
	ss.Bg.Endd  = 30.0
	err = ss.CalcTu()
	if err != nil{
		t.Fatal(err)
	}
	t.Log("maity lec 22 ex. 1(b) bolted 2xisa75508 shorter leg")
	ss.Cleg = 2
	err = ss.CalcTu()
	if err != nil{
		t.Fatal(err)
	}
	t.Log("maity lec 23 ex. 1(a) welded isa 90606 longer leg")
	sname = "ln"
	sdx = 33
	ss, err = GetStlSec(sname, sdx, code)
	ss.Weld = true
	if err != nil{
		t.Fatal(err)
	}
	ss.Printz()
	ss.Wg.Wltyp = 2
	ss.Wg.L1 = 75.0
	ss.Wg.L2 = 75.0
	ss.Wg.Tp = 10.0
	err = ss.CalcTu()

	if err != nil{
		t.Fatal(err)
	}
	t.Log("maity lec 23 ex. 1(a) welded isa 90606 shorter leg")
	ss.Cleg = 2
	err = ss.CalcTu()
	if err != nil{
		t.Fatal(err)
	}
}

func TestCalcPu(t *testing.T){
	t.Log("starting steel section compressive strength calcs")
	t.Log("maity lec 30 ex. 1 ismb400")
	sname := "i"
	code := 1
	sdx := 21
	ss, err := GetStlSec(sname, sdx, code)
	if err != nil{
		t.Fatal(err)
	}
	ss.Lspan = 3500.0
	err = ss.CalcPu()
	if err != nil{
		t.Fatal(err)
	}
	
	t.Log("maity lec 30 ex. 1 ishb250+(300x16)")
	sname = "plate-i"
	sdx = 25
	ss, _ = GetStlSec(sname, sdx, code, 300, 16)
	ss.Lspan = 4000.0
	ss.Klx = 0.8
	ss.Kly = 0.8
	err = ss.CalcPu()
	if err != nil{
		t.Fatal(err)
	}

	t.Log("maity lec 32 ex. 1 isa15015012")
	sname = "l"
	sdx = 65
	for i := 0; i < 3; i++{
		ss, _ = GetStlSec(sname, sdx, code)
		ss.Lspan = 3000.0
		ss.Klx = 1.0
		ss.Kly = 1.0
		switch i{
			case 0:
			t.Log("case 1 one bolt at each end")
			ss.Bg.Nb = 1	
			case 1:
			t.Log("case 2 two bolts at each end")
			ss.Bg.Nb = 2
			case 2:
			t.Log("case 3 welded at each end")
			ss.Weld = true
		}
		err = ss.CalcPu()
		if err != nil{
			t.Fatal(err)
		}
	}
	t.Log("maity lec 33 ex 1 2xisa100758")
	sdx = 41
	for i := 0; i < 2; i++{
		switch i{
			case 0:
			t.Log("opposite side of gusset plate")
			sname = "ln2-os"
			case 1:
			t.Log("same side of gusset plate")
			sname = "ln2-ss"
		}
		ss, _ = GetStlSec(sname, sdx, code)
		ss.Lspan = 4000.0
		ss.Klx = 0.85
		ss.Kly = 0.85
		err = ss.CalcPu()
		if err != nil{
			t.Fatal(err)
		}
		t.Log("pu -",ss.Fcd*ss.Area/1e3)
	}
}

func TestCalcMur(t *testing.T){
	t.Log("starting steel section ult. moment calcs")

	t.Log("maity lec 48 islb500")
	sname := "i"
	code := 1
	sdx := 10
	ss, err := GetStlSec(sname, sdx, code)
	if err != nil{
		t.Fatal(err)
	}
	ss.Lsb = true
	err = ss.CalcMur()
	if err != nil{
		t.Fatal(err)
	}
	
	t.Log("duggal ex.9.3 islb350")
	sdx = 28
	ss, err = GetStlSec(sname, sdx, code)
	if err != nil{
		t.Fatal(err)
	}

	ss.Lsb = true
	err = ss.CalcMur()
	if err != nil{
		t.Fatal(err)
	}
	ss.Lsb = false
	ss.Lspan = 3000.0
	err = ss.CalcMur()
	if err != nil{
		t.Fatal(err)
	}
	t.Fatal()

}

func TestBmChk800(t *testing.T){
	t.Log("starting is800 beam checks")
	t.Log("duggal ex 9.4 islb200")
	sname := "i"
	code := 1
	sdx := 50
	ss, err := GetStlSec(sname, sdx, code)
	if err != nil{
		t.Fatal(err)
	}
	ss.Mu = 30e3
	ss.Vu = 30e3
	ss.Dmax = 9.1
	ss.Lsb = true
	ss.Lspan = 4000.0
	ss.Endc = 1
	err = ss.BmChk800()
	if err != nil{
		t.Fatal(err)
	}
	
	t.Log("starting biaxial beam (purlin) check(is800) tests")

	t.Log("duggal ex 9.9 islb150")
	sname = "i"
	code = 1
	sdx = 52
	ss, err = GetStlSec(sname, sdx, code)
	if err != nil{
		t.Fatal(err)
	}
	ss.Printz()
	ss.Ax = -2
	ss.Mux = 17.6 * 1e6
	ss.Muy = 0.8 * 1e6
	ss.Lsb = true
	err = ss.BmChk800()
	if err != nil{
		t.Fatal(err)
	}
}

func TestBmChk449(t *testing.T){
	t.Log("testing mosley bs449 beam examples")
	t.Log("example 6.2")
	sdxs := []int{5}
	for _, sdx := range sdxs{
		ss, _ := GetStlSec("ub", sdx, 2)
		ss.Mu = 970.3e6
		ss.Vu = 585e3
		ss.Dmax = 2.23
		ss.Lspan = 6000.0
		ss.Tx = 1.0
		ss.Lcf = 6000.0
		ss.Tcf = 1.0
		ss.Brchk = true
		ss.Lbr = 200.0
		ss.Tbr = 20.0
		err := ss.BmChk449()
		if err != nil{
			t.Log(err)
		}
	}
	t.Log("duggal ex 9.4 ub")
	sname := "ub"
	code := 1
	sdx := 58
	ss, err := GetStlSec(sname, sdx, code)
	if err != nil{
		t.Fatal(err)
	}
	ss.Printz()
	ss.Mu = 30e3
	ss.Vu = 30e3
	ss.Dmax = 9.1
	ss.Lsb = true
	ss.Lspan = 4000.0
	ss.Endc = 1
	err = ss.BmChk800()
	if err != nil{
		t.Fatal(err)
	}
}

func TestColChk449(t *testing.T){
	t.Log("testing mosley bs449 column examples")
	t.Log("example 6.1->")
	sdxs := []int{21,22,23}
	for _, sdx := range sdxs{
		ss, _ := GetStlSec("uc", sdx, 2)
		t.Log("checking section",ss.Sstr)
		ss.Lx = 3500.0
		ss.Ly = 3500.0
		ss.Tx = 1.0
		ss.Ty = 1.0
		ss.Pu = 1000.0*1e3
		ss.Vbdx = 120*1e3
		ss.Vbdy = 40*1e3
		ss.Pfac = 1.0
		ss.H1 = 3500.0
		ss.H2 = 4000.0
		err := ss.ColChk449()
		if err != nil{
			t.Log(err)
		}
	}
	t.Log("example 6.2->")
	sdxs = []int{16,21}
	//sdx = 17,22
	for _, sdx := range sdxs{
		ss, _ := GetStlSec("ub", sdx, 2)
		ss.H1 = 0.0
		ss.H2 = 0.0
		ss.Lx = 2400.0
		ss.Ly = 4800.0
		ss.Tx = 1.0
		ss.Ty = 0.9
		ss.Pu = 500.0*1e3
		ss.Mux = 120.0*1e6
		ss.Muy = 45.0*1e6
		ss.Pfac = 1.25
		ss.Dtyp = 1
		err := ss.ColChk449()
		if err != nil{
			t.Fatal(err)
		}
	}

}

func TestCmpBmChk(t *testing.T){
	sdx := 16
	ss, _ := GetStlSec("i", sdx, 1)
	ss.Printz()
	ss.Lspan = 10000
	err := ss.CalcMur()
	fmt.Println("ERRORE-",err)
	fmt.Println("mur-",ss.Mur/1e6,"knm")
	fmt.Println("dividing lspan/3")
	ss.Lspan = ss.Lspan/3.0
	err = ss.CalcMur()
	fmt.Println("ERRORE-",err)
	fmt.Println("mur-",ss.Mur/1e6,"knm")

}
