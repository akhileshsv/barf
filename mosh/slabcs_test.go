package barf

import (
	"fmt"
	"testing"
)

func TestRSlb2Chk(t *testing.T){
	var s *RccSlb
	//mosley ex. 8.9
	s = &RccSlb{
		Fck:25.0,
		Fy:500.0,
		Diamain:10.0,
		Diadist:8.0,
		Lx:6000.0,
		Ly:6000.0,
		LL:2.5,
		DL:6.0,
		Bsup:250.0,
		Type:4,
		Endc:1,
		Nomcvr:30.0,
		Efcvr:40.0,
		Ibent:0.0,
		Bf:400.0,
		Bw:125.0,
		Df:60.0,
		Dw:200.0,
		Code:3,
		Ldcalc:1,
		Verbose:true,
	}
	_, _ = RSlb2Chk(s)

	s = &RccSlb{
		Fck:25.0,
		Fy:500.0,
		Diamain:10.0,
		Diadist:8.0,
		Lx:6000.0,
		Ly:6000.0,
		LL:2.5,
		DL:2.0,
		Bsup:250.0,
		Type:2,
		Endc:1,
		Nomcvr:30.0,
		Efcvr:40.0,
		Ibent:0.0,
		Code:2,
		Ldcalc:1,
		Verbose:true,
	}
	
	err := SlbDesign(s)
	if err != nil{
		fmt.Println(err)
	}
	s = &RccSlb{
		Fck:25.0,
		Fy:415.0,
		Diamain:10.0,
		Diadist:8.0,
		Lx:3600.0,
		Ly:3900.0,
		LL:5.0/1.5,
		DL:5.0/1.5,
		Bsup:250.0,
		Type:4,
		Endc:4,
		Nomcvr:15.0,
		Efcvr:20.0,
		Ibent:0.0,
		Bf:1200.0,
		Bw:160.0,
		Df:25.0,
		Dw:150.0,
		Code:1,
		Ldcalc:1,
		Verbose:true,
	}
	//_, _ = RSlb2Chk(s)
	//varghese ex. 5.1
}

func TestRSlb1Chk(t *testing.T){
	//var rezstring string
	var s *RccSlb
	//first check for a simply supported slab
	//test sub example
	s = &RccSlb{
		Fck:25.0,
		Fy:415.0,
		Diamain:10.0,
		Diadist:8.0,
		Lx:3500.0,
		Ly:9000.0,
		LL:3.0,
		DL:2.0,
		Bsup:250.0,
		Type:3,
		Endc:1,
		Lspan:3500.0,
		Nomcvr:25.0,
		Efcvr:30.0,
		Ibent:0.0,
		Nspans:1,
		Bf:900.0,
		Bw:75.0,
		Df:75.0,
		Dw:300.0,
		Dtyp:1,
		Code:1,
		Verbose:true,
	}
	_, _ = RSlb1Chk(s)
	
	s.Type = 1

	err := SlbDesign(s)
	if err != nil{
		fmt.Println(err)
	}
	//then mosley example 8.8
	s = &RccSlb{
		Fck:30.0,
		Fy:250.0,
		Diamain:10.0,
		Diadist:8.0,
		Lx:5000.0,
		Ly:10000.0,
		LL:2.5,
		DL:4.5,
		Bsup:250.0,
		Type:3,
		Endc:2,
		Lspan:5000.0,
		Nomcvr:35.0,
		Efcvr:40.0,
		Ibent:0.0,
		Nspans:4,
		Bf:400.0,
		Bw:125.0,
		Df:60.0,
		Dw:200.0,
		Dtyp:1,
		Code:2,
		DM:0.2,
		Ldcalc:1,
	}
	//_, _ = RSlb1Chk(s)
}

func TestCSlb1Dz(t *testing.T){
	var rezstring string
	
	rezstring += "subramanian ex 9.2\n"
	var s *RccSlb
	fmt.Println("from coefficients->")
	s = &RccSlb{
		Fck:20.0,
		Fy:415.0,
		Diamain:10.0,
		Diadist:8.0,
		LL:3.0,
		DL:1.6,
		Bsup:250.0,
		Type:1,
		Endc:2,
		Lspan:4000.0,
		Nomcvr:25.0,
		Efcvr:30.0,
		Ibent:0.0,
		Nspans:4,
	}
	err := CSlb1DepthCs(s)
	if err != nil{
		t.Errorf("continuous slab design (coefficients) test failed")
	}
	err = CSlb1Stl(s)

	if err != nil{
		t.Errorf("continuous slab design (coefficients) test failed")
	}

	for i := 0; i < s.Nspans; i++{
		rezstring += fmt.Sprintf("span no %v \n",i+1)
		for j, sec := range []string{"left","middle","right"}{
			rezstring += fmt.Sprintf("%s section\n",sec)
			rezstring += fmt.Sprintf("dia - %.f mm spacing %.f mm dist %.f mm spacing %.f mm\n",s.Diaspns[i][j], s.Spcspns[i][j],s.Diadist, s.Sdspns[i][j])
			rezstring += fmt.Sprintf("ast req - %.f mm2 ast prov %.f mm2\n",s.Astrs[i][j], s.Astps[i][j])
		}
	}
	fmt.Println(rezstring)
	fmt.Println("from env->")
	s = &RccSlb{
		Fck:20.0,
		Fy:415.0,
		Diamain:10.0,
		Diadist:8.0,
		LL:3.0,
		DL:1.6,
		Bsup:0.0,
		Type:1,
		Endc:2,
		Lspan:4000.0,
		Nomcvr:25.0,
		Efcvr:30.0,
		Ibent:0.0,
		Nspans:4,
		DM:0.0,
	}
	/*
	err = CSlb1Depth(s)
	if err != nil{
		t.Log(err)
	}
	err = CSlb1Stl(s)
	*/

}

func TestCSlb1DepthCs(t *testing.T) {
	var rezstring string
	rezstring += "subramanian ex 9.2\n"
	var s *RccSlb
	fmt.Println("from coefficients->")
	s = &RccSlb{
		Fck:20.0,
		Fy:415.0,
		Diamain:10.0,
		Diadist:8.0,
		LL:3.0,
		DL:1.6,
		Bsup:250.0,
		Type:1,
		Endc:2,
		Lspan:4000.0,
		Nomcvr:25.0,
		Efcvr:30.0,
		Ibent:0.0,
		Nspans:4,
	}
	err := CSlb1DepthCs(s)
	if err != nil{
		t.Errorf("continuous slab design (coefficients) test failed")
	}

	fmt.Println("from env->")
	s = &RccSlb{
		Fck:20.0,
		Fy:415.0,
		Diamain:10.0,
		Diadist:8.0,
		LL:3.0,
		DL:1.6,
		Bsup:0.0,
		Type:1,
		Endc:2,
		Lspan:4000.0,
		Nomcvr:25.0,
		Efcvr:30.0,
		Ibent:0.0,
		Nspans:4,
		DM:0.0,
	}
	err = CSlb1Depth(s)
	if err != nil{
		t.Log(err)
	}
	//err = CSlb1Stl(s)


	
	rezstring += "\n\nmosley ex 8.4\n"
	fmt.Println("coeff again->")
	s = &RccSlb{
		Fck:25.0,
		Fy:500.0,
		Diamain:10.0,
		Diadist:6.0,
		LL:3.0,
		DL:1.0,
		Bsup:250.0,
		Type:1,
		Endc:2,
		Lspan:4500.0,
		Nomcvr:25.0,
		Efcvr:0.0,
		Ibent:0.0,
		Nspans:4,
		Code:2,
	}
	err = CSlb1DepthCs(s)

	fmt.Println("env->")	
	s = &RccSlb{
		Fck:25.0,
		Fy:500.0,
		Diamain:10.0,
		Diadist:6.0,
		LL:3.0,
		DL:1.0,
		Type:1,
		Endc:2,
		Lspan:4500.0,
		Nomcvr:25.0,
		Efcvr:0.0,
		Ibent:0.0,
		Nspans:4,
		Code:2,
		DM:0.0,
	}

	err = CSlb1Depth(s)
	wantstring := ``
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("continuous slab design (coefficients) test failed")
	}
}

