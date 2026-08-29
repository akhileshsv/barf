package barf

import (
		"os"
"path/filepath"
	"fmt"
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

func TestSlbDxf(t *testing.T){
	var rezstring string
	rezstring += "shah 1 way slab examples section 6.2\n"
	s := &RccSlb{
		Title:"shah6.2-slab",
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
		Term:"dxf",
	}

	err := SlbDIs(s)
	if err != nil{
		t.Errorf("slab dxf test failed")
	}
	s.Draw(s.Term)
	//s.DrawPlan()
	//mosley ex 8.5
	
	s = &RccSlb{
		Title:"mosley8.5-slab",
		Fck:30.0,
		Fy:460.0,
		Lx:4500.0,
		Ly:6300.0,
		Bsup:230.0,
		DL:0.0,
		LL:10.0,
		Diamain:12.0,
		Diadist:10.0,
		Code:2,
		Type:2,
		Endc:10,
		Nomcvr:25.0,
		Verbose:true,
		Term: "dxf",
	}
	err = Slb2DBs(s)
	if err != nil{
		t.Errorf("slab dxf test failed")
	}
	s.Draw(s.Term)
	//s.DrawPlan()
	fmt.Println("done")
}

func TestSlbDraw(t *testing.T){
	var examples = []string{"slbs3"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/slab/")
	for _, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		s, err := ReadSlb(fname)
		if err != nil{
			t.Errorf("slab draw test failed")
		}
		s.Term = "qt"
		err = SlbDesign(&s)
		if err != nil{
			fmt.Println(err)
			t.Errorf("slab draw test failed")
		}
		//pltstr := s.Draw(s.Term)
		
		//fmt.Println(pltstr)
	}

}
