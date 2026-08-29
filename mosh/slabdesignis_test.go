package barf

import (
	"os"
	"path/filepath"
	"fmt"
	"testing"
)

func TestSlb2DIs(t *testing.T){
	var examples = []string{"slab_sub10.1","slab_sub10.3"}
	var rezstring string
	//rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	t.Log(ColorPurple,"testing slab design (is code)\n",ColorReset)
	for i, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		s, err := ReadSlb(fname)
		if err != nil{
			fmt.Println(ColorRed,err,ColorReset)
			t.Errorf("slab design (is code) test failed")
		}	
		err = SlbDIs(&s)
		if err != nil{
			fmt.Println(ColorRed,err,ColorReset)
			t.Errorf("slab design (is code) test failed")
		}
		//s.Draw("dumb")
		//s.Report()
		rezstring += fmt.Sprintf("%s\n",s.Printz())
		//t.Log(s.Rprt)
	}
	wantstring := `10.1Sub
2 way ss slab 
lspan 0.0, lx 4530.0, ly 6230.0 mm
grade of concrete M 25.0, steel - main Fe 415, dist Fe 415
cover - nominal 15.0 mm, effective 20.0 mm
loads -dl 1.0 kN/m2, ll 2.0 kN/m2
design code is design type 0
short span support -> mdu 0.00 kn-m 8 mm dia at 300 mm spacing 
ast req - 0 mm2 ast - 168 mm2
short span midspan -> mdu 22.14 kn-m 8 mm dia at 120 mm spacing 
ast req - 414 mm2 ast - 419 mm2
long span support -> mdu 0.00 kn-m 8 mm dia at 300 mm spacing 
ast req - 0 mm2 ast - 168 mm2
long span midspan -> mdu 11.80 kn-m 8 mm dia at 220 mm spacing 
ast req - 228 mm2 ast - 228 mm2

10.3Sub
2 way endc 9 slab 
lspan 0.0, lx 4530.0, ly 6230.0 mm
grade of concrete M 25.0, steel - main Fe 415, dist Fe 415
cover - nominal 15.0 mm, effective 20.0 mm
loads -dl 1.0 kN/m2, ll 2.0 kN/m2
design code is design type 0
short span support -> mdu 0.00 kn-m 8 mm dia at 300 mm spacing 
ast req - 0 mm2 ast - 168 mm2
short span midspan -> mdu 18.64 kn-m 8 mm dia at 140 mm spacing 
ast req - 359 mm2 ast - 359 mm2
long span support -> mdu 0.00 kn-m 8 mm dia at 300 mm spacing 
ast req - 0 mm2 ast - 168 mm2
long span midspan -> mdu 12.50 kn-m 8 mm dia at 200 mm spacing 
ast req - 251 mm2 ast - 251 mm2
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("slab design (is code) test failed")
	}
	
}

func TestSlbSsShah(t *testing.T){
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
		t.Errorf("1 way ss slab test failed")
	}
	rezstring += s.Printz()
	s = &RccSlb{
		Fck:15.0,
		Fy:250.0,
		Diamain:10.0,
		Diadist:6.0,
		LL:4.0,
		DL:1.5,
		Bsup:230.0,
		Type:1,
		Endc:0,
		Lspan:1500.0,
		Nomcvr:15.0,
		Efcvr:20.0,
		Ibent:0.0,
	}
	err = SlbDIs(s)
	rezstring += s.Printz()
	t.Log(rezstring)
	wantstring := `shah 1 way slab examples section 6.2
rcc slab
1 way slab 
lspan 2500.0, lx 0.0, ly 0.0 mm
grade of concrete M 15.0, steel - main Fe 415, dist Fe 250
cover - nominal 16.0 mm, effective 20.0 mm
loads -dl 1.8 kN/m2, ll 1.5 kN/m2
design code is design type 0
midspan moment 7.03 kn-m
ast required - 233.2 mm2
ast provided - 239.4 mm2
rcc slab
cantilever slab 
lspan 1500.0, lx 0.0, ly 0.0 mm
grade of concrete M 15.0, steel - main Fe 250, dist Fe 250
cover - nominal 15.0 mm, effective 20.0 mm
loads -dl 1.5 kN/m2, ll 4.0 kN/m2
design code is design type 0
midspan moment 15.19 kn-m
ast required - 638.9 mm2
ast provided - 654.5 mm2
extend bar to 545 mm from support`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("1 way ss slab (is) test failed")
	}

}

func TestSlb2WShah(t *testing.T){
	var rezstring string
	rezstring += "shah 2 way slab examples section 6.3.7\n"
	//shah ex 6.3.7
	s := &RccSlb{
		Fck:15.0,
		Fy:415.0,
		Fyd:250.0,
		Diamain:8.0,
		Diadist:6.0,
		LL:1.5,
		DL:1.75,
		Bsup:0.0,
		Type:2,
		Endc:4,
		Lspan:0.0,
		Lx:4500.0,
		Ly:5000.0,
		Nomcvr:16.0,
		Efcvr:20.0,
		Ibent:0.0,
	}
	err := SlbDIs(s)
	if err != nil{
		t.Log(err)
		t.Errorf("2 way slab (is) test failed")
	}
	rezstring += s.Printz()
	//fmt.Println("simply supported 2 way")
	s = &RccSlb{
		Fck:15.0,
		Fy:415.0,
		Fyd:250.0,
		Diamain:8.0,
		Diadist:6.0,
		LL:4.0,
		DL:1.0,
		Bsup:200.0,
		Type:2,
		Endc:10,
		Lspan:0.0,
		Lx:4120.0,
		Ly:5620.0,
		Nomcvr:16.0,
		Efcvr:20.0,
		Ibent:0.0,
	}
	err = SlbDIs(s)	

	if err != nil{
		t.Log(err)
		t.Errorf("2 way slab (is) test failed")
	}

	rezstring += s.Printz()
	
	//s.Report()
	//t.Log(s.Rprt)
	wantstring := `shah 2 way slab examples section 6.3.7
rcc slab
2 way endc 4 slab 
lspan 0.0, lx 4500.0, ly 5000.0 mm
grade of concrete M 15.0, steel - main Fe 415, dist Fe 250
cover - nominal 16.0 mm, effective 20.0 mm
loads -dl 1.8 kN/m2, ll 1.5 kN/m2
design code is design type 0
short span support -> mdu 10.62 kn-m 8 mm dia at 170 mm spacing 
ast req - 288 mm2 ast - 296 mm2
short span midspan -> mdu 8.01 kn-m 8 mm dia at 230 mm spacing 
ast req - 213 mm2 ast - 219 mm2
long span support -> mdu 9.28 kn-m 8 mm dia at 200 mm spacing 
ast req - 249 mm2 ast - 251 mm2
long span midspan -> mdu 6.91 kn-m 8 mm dia at 250 mm spacing 
ast req - 198 mm2 ast - 201 mm2
rcc slab
2 way ss slab 
lspan 0.0, lx 4120.0, ly 5620.0 mm
grade of concrete M 15.0, steel - main Fe 415, dist Fe 250
cover - nominal 16.0 mm, effective 20.0 mm
loads -dl 1.0 kN/m2, ll 4.0 kN/m2
design code is design type 0
short span support -> mdu 0.00 kn-m 8 mm dia at 300 mm spacing 
ast req - 0 mm2 ast - 168 mm2
short span midspan -> mdu 22.81 kn-m 8 mm dia at 100 mm spacing 
ast req - 460 mm2 ast - 503 mm2
long span support -> mdu 0.00 kn-m 8 mm dia at 300 mm spacing 
ast req - 0 mm2 ast - 168 mm2
long span midspan -> mdu 12.35 kn-m 8 mm dia at 190 mm spacing 
ast req - 254 mm2 ast - 265 mm2
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("2 way slab (is) test failed")
	}
}

