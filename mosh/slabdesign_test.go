package barf

import (
	"fmt"
	"testing"
)

func TestBalSecAst(t *testing.T){
	var mdu, effd, fck, fy, astr float64
	var code int
	var rezstring string

	mdu = 18.05; effd = 145.0; fck = 25.0; fy = 415; code = 1
	astr = BalSecAst(mdu, effd, fck, fy , code)
	rezstring += fmt.Sprintf("mdu %f kn effd %f mm fck %f fy %f code %v\n",mdu, effd, fck, fy, code)
	rezstring += fmt.Sprintf("ast req -> %f mm2\n",astr)

	
	mdu = 12.03; effd = 135.0; fck = 25.0; fy = 415; code = 1
	astr = BalSecAst(mdu, effd, fck, fy , code)
	rezstring += fmt.Sprintf("mdu %f kn effd %f mm fck %f fy %f code %v\n",mdu, effd, fck, fy, code)
	rezstring += fmt.Sprintf("ast req -> %f mm2\n",astr)
	
	mdu = 31.9; effd = 170.0; fck = 25.0; fy = 500; code = 2
	astr = BalSecAst(mdu, effd, fck, fy , code)
	rezstring += fmt.Sprintf("mdu %f kn effd %f mm fck %f fy %f code %v\n",mdu, effd, fck, fy, code)
	rezstring += fmt.Sprintf("ast req -> %f mm2\n",astr)


	mdu = 20.18; effd = 140.0; fck = 25.0; fy = 500; code = 2
	astr = BalSecAst(mdu, effd, fck, fy , code)
	rezstring += fmt.Sprintf("mdu %f kn effd %f mm fck %f fy %f code %v\n",mdu, effd, fck, fy, code)
	rezstring += fmt.Sprintf("ast req -> %f mm2\n",astr)

	mdu = 45.0; effd = 185.0; fck = 25.0; fy = 500; code = 2
	astr = BalSecAst(mdu, effd, fck, fy , code)
	rezstring += fmt.Sprintf("mdu %f kn effd %f mm fck %f fy %f code %v\n",mdu, effd, fck, fy, code)
	rezstring += fmt.Sprintf("ast req -> %f mm2\n",astr)

	code = 1
	astr = BalSecAst(mdu, effd, fck, fy , code)
	rezstring += fmt.Sprintf("mdu %f kn effd %f mm fck %f fy %f code %v\n",mdu, effd, fck, fy, code)
	rezstring += fmt.Sprintf("ast req -> %f mm2\n",astr)

	wantstring :=`mdu 18.050000 kn effd 145.000000 mm fck 25.000000 fy 415.000000 code 1
ast req -> 359.770257 mm2
mdu 12.030000 kn effd 135.000000 mm fck 25.000000 fy 415.000000 code 1
ast req -> 254.925393 mm2
mdu 31.900000 kn effd 170.000000 mm fck 25.000000 fy 500.000000 code 2
ast req -> 449.618917 mm2
mdu 20.180000 kn effd 140.000000 mm fck 25.000000 fy 500.000000 code 2
ast req -> 344.359571 mm2
mdu 45.000000 kn effd 185.000000 mm fck 25.000000 fy 500.000000 code 2
ast req -> 587.840334 mm2
mdu 45.000000 kn effd 185.000000 mm fck 25.000000 fy 500.000000 code 1
ast req -> 598.137077 mm2
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("balanced section test failed")
	}
}

func TestSlb1DBs(t *testing.T){
	//mosley ex 8.3
	var rezstring string
	s := &RccSlb{
		Fck:30.0,
		Fy:460.0,
		Lspan:4500.0,
		DL:1.0,
		LL:3.0,
		Nomcvr:25.0,
		Type:1,
		Endc:1,
		Code:2,
		DM:0.0,
		Diamain:10.0,
		Diadist:10.0,
	}
	//err := Slb1DBs(s)
	rezstring += "mosley example 8.3\n"
	//s.Term = "dumb"
	err := SlbDesign(s)
	if err != nil{
		fmt.Println(err)
		t.Errorf("1 way slab design test (bs) failed")
	}
	rezstring += fmt.Sprintf("net cost %.2f\n",s.Kost)
	s.Table(false)
	t.Log(s.Report)
	rezstring += "mosley example 8.4\n"
	s = &RccSlb{
		Fck:25.0,
		Fy:500.0,
		Diamain:10.0,
		Diadist:6.0,
		LL:3.0,
		DL:0.5,
		Bsup:250.0,
		Type:1,
		Endc:2,
		Lspan:4500.0,
		Nomcvr:25.0,
		Efcvr:0.0,
		Ibent:0.0,
		Nspans:4,
		Code:2,
		DM:0.2,
	}
	err = SlbDesign(s)

	if err != nil{
		fmt.Println(err)
		t.Errorf("1 way slab design test (bs) failed")
	}
	rezstring += fmt.Sprintf("net cost %.2f\n",s.Kost)
	//s.Term = "mono"
	s.Table(false)
	t.Log(s.Report)
	//rezstring += s.Printz()
	//fmt.Println(s.Report)
	wantstring := `mosley example 8.3
net cost 170484.43
mosley example 8.4
net cost 532040.39
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("1 way slab design test (bs) failed")
	}
	//s.Draw("dumb")
	
	//s.Draw("svg")
}

func TestSlb2DBs(t *testing.T){
	//mosley ex 8.5
	s := &RccSlb{
		Fck:30.0,
		Fy:460.0,
		Lx:4500.0,
		Ly:6300.0,
		DL:0.0,
		LL:10.0,
		Diamain:12.0,
		Diadist:10.0,
		Code:2,
		Type:2,
		Endc:10,
		Nomcvr:25.0,
		Verbose:true,
	}
	err := Slb2DBs(s)
	if err != nil{fmt.Println(err)}
	//allen ex 12.1
	s = &RccSlb{
		Fck:35.0,
		Fy:460.0,
		Lx:6000.0,
		Ly:7200.0,
		DL:1.0,
		LL:5.0,
		Diamain:12.0,
		Diadist:10.0,
		Code:2,
		Type:2,
		Endc:4,
		Nomcvr:20,
		Verbose:true,
	}
	err = Slb2DBs(s)
	if err != nil{fmt.Println(err)}
}

func TestSlb2BmCoefBs(t *testing.T){
	endc := 0; ns := 0; nl := 1
	lx := 5.0; ly := 6.0
	udl := 1.4 * 3.0 + 1.6 * 5.0
	mcxsup, mcxmid, mcysup, mcymid := Slb2BMCoefBs(endc, ns, nl, lx, ly)
	var rezstr string
	tx := lx * lx * udl; ty := lx * lx * udl
	mxsup := tx * mcxsup; mxspan := tx * mcxmid; mysup := ty * mcysup; myspan := ty * mcymid
	rezstr += fmt.Sprintf("hulse example 4.3.6 two way slab ns %v nl %v lx %.f ly %.f\n",ns, nl, lx, ly)
	rezstr += fmt.Sprintf("short span edge moment %.2f kn.m, short midspan moment %.2f kn.m\n",mxsup, mxspan)
	rezstr += fmt.Sprintf("long span edge moment %.2f, long midspan moment %.2f\n",mysup, myspan)
	wantstr := ``
	if rezstr != wantstr{
		fmt.Println(rezstr)
		t.Errorf("2- way slab bending moment coefficient (bs) test failed")
	}
}
