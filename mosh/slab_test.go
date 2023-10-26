package barf

import (
	//"fmt"
	"testing"
)

func TestSlbQuant(t *testing.T){
	var s RccSlb
	s = RccSlb{
		Fck:20,
		Fy:415,
		Dz:true,
		Lx:2500,
		Ly:6000,
		Lspan:2500,
		Type:1,
		Endc:1,
		Diamain:10.0,
		Diadist:8.0,
		Dused:150.0,
		Efcvr:25.0,
		Nomcvr:20.0,
		BM:[]float64{0,0,0,0},
		Spcm:150,
		Spcd:150,
	}
	s.Quant()
	s.Table(true)

	s = RccSlb{
		Fck:20,
		Fy:415,
		Dz:true,
		Lx:5400,
		Ly:6000,
		Type:2,
		Endc:4,
		Diamain:10.0,
		Diadist:8.0,
		Dused:150.0,
		Efcvr:25.0,
		Nomcvr:20.0,
		BM:[]float64{0,0,0,0},
		Spcms:[]float64{0,0,0,0},
		Spcds:[]float64{0,0,0,0},
		Spcm:150,
		Spcd:150,
	}
	s.Quant()
	s.Table(true)	
}

func TestSlbDraw(t *testing.T){
	var rezstring string
	rezstring += "shah 1 way slab examples section 6.2\n"
	s := &RccSlb{
		Fck:15.0,
		Fy:415.0,
		Diamain:8.0,
		Diadist:6.0,
		LL:1.5,
		DL:1.75,
		Bsup:200.0,
		Type:1,
		Endc:1,
		Lspan:2500.0,
		Nomcvr:16.0,
		Efcvr:20.0,
		Ibent:50.0,
	}

	err := SlbDIs(s)
	if err != nil{
		t.Errorf("slab draw test failed")
	}
	s.Draw("dumb")
	//s.Draw("qt")
	//s.Draw("svg")
	//mosley ex 8.5
	s = &RccSlb{
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
	err = Slb2DBs(s)
	if err != nil{
		t.Errorf("slab draw test failed")
	}
	s.Draw("dumb")

}
