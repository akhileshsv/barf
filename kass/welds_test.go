package barf

import (
	"testing"
)



func TestWeldSs(t *testing.T){
	w := &Wld{
		Sizes:[]float64{10,10,10},
		Wcs:[][]float64{{0,1200,800,1200},{0,1200,0,0},{0,0,800,0}},
		Frcs:[][]float64{{2,-40,0}},
		Fcs:[][]float64{{1400,1200}},
	}
	err := WeldSs(w)
	t.Log("\n\n",w.Report)
	if err != nil{
		t.Errorf("weld group analysis test failed")
	}
}

func TestWeldCalc(t *testing.T){
	t.Log("starting weld check calcs duggal ex. 7.3 case 1-[")
	w := &Wld{
		Size:6.0,
		Tmem:14.0,
		Tp:10.0,
		Lop:300.0,
		Dmem:250.0,
		Wltyp:2,
		Ltyp:1,
		Ctyp:1,
		Shop:true,
	}
	err := WeldCalc(w)
	if err != nil{
		t.Log(err)
		t.Errorf("weld calc test failed")
	}
	
	t.Log("starting weld check calcs duggal ex. 7.3 case 2-[]")
	w.Wltyp = 3
	err = WeldCalc(w)
	if err != nil{
		t.Log(err)
		t.Errorf("weld calc test failed")
	}
	
}

func TestWeldDz(t *testing.T){
	t.Log("starting weld dz calcs duggal ex. 6.4")
	w := &Wld{
		Pu:145e3,
		Tmem:8.0,
		Tp:12.0,
		Dmem:75.0,
		Ltyp:1,
		Ctyp:1,
		Shop:true,
		Msname:"flat",
	}
	for wl := 1; wl < 4; wl++{
		t.Log("case",wl)
		w.Wltyp = wl
		if wl == 3{
			w.Shop = false
		}
		err := WeldDz(w)
		if err != nil{
			t.Log(err)
			t.Errorf("weld calc test failed")
		}
	}
	t.Log("isa 80x50x8 duggal ex. 6.7/6.8")
	w = &Wld{
		Pu:223e3,
		Tmem:8.0,
		Tp:12.0,
		Dims:[][]float64{{50,80,8,8}},
		Wltyp:1,
		Ltyp:1,
		Ctyp:1,
		Shop:false,
		Msname:"l",
		Size:6.0,
	}
	for wl := 1; wl < 3; wl++{
		w.Wltyp = wl
		err := WeldDz(w)
		if err != nil{
			t.Log(err)
			t.Errorf("weld calc test failed")
		}
	}
	t.Log("isa 80x80x8 duggal ex. 6.9")
	
	w = &Wld{
		Pu:300e3,
		Tmem:8.0,
		Tp:12.0,
		Dims:[][]float64{{80,80,8,8,12}},
		Wltyp:1,
		Ltyp:1,
		Ctyp:1,
		Shop:true,
		Msname:"l2-os",
	}
	err := WeldDz(w)
	if err != nil{
		t.Log(err)
		t.Errorf("weld calc test failed")
	}

	t.Log("isa 100x75x8 maity ex 10.1")
	
	w = &Wld{
		Pu:303e3,
		Tmem:8.0,
		Tp:10.0,
		Dims:[][]float64{{75,100,8,8,12}},
		Wltyp:2,
		Ltyp:1,
		Ctyp:1,
		Shop:true,
		Msname:"ln",
		Size:5.0,
	}
	err = WeldDz(w)
	if err != nil{
		t.Log(err)
		t.Errorf("weld calc test failed")
	}

}


