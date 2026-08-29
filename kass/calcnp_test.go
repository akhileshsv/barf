package barf

import (
	"testing"
	"os"
	"path/filepath"
	"fmt"
	"log"
)

func TestKFacMos(t *testing.T) {
	var rezstring string
	rezstring += "mosley ex 4.2\n"
	mem := &MemNp{
		Ts:[]int{2,0,2},
		Ls:[]float64{16,12,12},
		Ds:[]float64{3.6,1.2,3.2},
		Bs:[]float64{1.0,1.0,1.0},
		Styp:1,
		Dims:[]float64{1.0,1.2},
	}
	KFacMos(mem, false)
	rezstring += fmt.Sprintf("Ka %.2f Kb %.2f Kc %.2f Ca %.2f Cb %.2f\n",mem.Ka * mem.Lspan * mem.Lspan, mem.Kb * mem.Lspan * mem.Lspan, mem.Kc * mem.Lspan * mem.Lspan, mem.Ca, mem.Cb)
	rezstring += "hulse ex 2.6\nspan 1\n"
	mem = &MemNp{
		Ts:[]int{3,0,3},
		Ls:[]float64{3,2,3},
		Ds:[]float64{1.6,0.8,1.6},
		Bs:[]float64{1.0,1.0,1.0},
		Styp:1,
		Dims:[]float64{1.0,0.8},
	}
	KFacMos(mem, false)
	rezstring += fmt.Sprintf("Ka %.2f Kb %.2f Kc %.2f Ca %.2f Cb %.2f\n",mem.Ka * mem.Lspan * mem.Lspan, mem.Kb * mem.Lspan * mem.Lspan, mem.Kc * mem.Lspan * mem.Lspan, mem.Ca, mem.Cb)

	rezstring += "hulse ex 2.6\nspan 2\n"
	mem = &MemNp{
		Ts:[]int{3,0,3},
		Ls:[]float64{2,2.5,1.5},
		Ds:[]float64{1.6,0.6,1.0},
		Bs:[]float64{1.0,1.0,1.0},
		Styp:1,
		Dims:[]float64{1.0,0.6},
	}
	KFacMos(mem, false)
	rezstring += fmt.Sprintf("Ka %.2f Kb %.2f Kc %.2f Ca %.2f Cb %.2f\n",mem.Ka * mem.Lspan * mem.Lspan, mem.Kb * mem.Lspan * mem.Lspan, mem.Kc * mem.Lspan * mem.Lspan, mem.Ca, mem.Cb)
	
	wantstring := `mosley ex 4.2
Ka 24.42 Kb 18.20 Kc 16.77 Ca 0.69 Cb 0.92
hulse ex 2.6
span 1
Ka 9.52 Kb 9.52 Kc 6.49 Ca 0.68 Cb 0.68
hulse ex 2.6
span 2
Ka 9.72 Kb 7.16 Kc 5.55 Ca 0.57 Cb 0.77
`
	t.Logf(rezstring)
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("stiffness and carry over factor test failed")
	}
}

// func TestNpFrmPortal(t *testing.T){
// 	var examples = []string{"sp40_out"}
// 	var rezstring string
// 	//var frmrez []interface{}
// 	//rezstring += "\n"
// 	dirname,_ := os.Getwd()
// 	datadir := filepath.Join(dirname,"../data/out")
// 	for i, ex := range examples {
// 		fname := filepath.Join(datadir,ex+".json")
// 		t.Log(ColorRed,"example->",i+1,ColorYellow,ex,ColorReset)
// 		rezstring += fmt.Sprintf("example %v %s\n",i+1,ex)
// 		ModNpInp(fname, "qt", false)	
// 	}
	
// }

func TestNpFrmDraw(t *testing.T){
	var examples = []string{"hulse2.6", "phil1","leet11.13","phil4"}
	var rezstring string
	//var frmrez []interface{}
	//rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/kass/npfrm")
	for i, ex := range examples {
		if i != 0{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorRed,"example->",i+1,ColorYellow,ex,ColorReset)
		rezstring += fmt.Sprintf("example %v %s\n",i+1,ex)
		ModNpInp(fname, "qt", false)
		
	}
}

//TestNpFrm tests npfrm and beam analysis results
func TestNpFrm(t *testing.T){
	var examples = []string{"hulse2.6","phil1","leet11.13","phil4"}
	//var examples = []string{"phil1"}
	var rezstring string
	//var frmrez []interface{}
	//rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/kass/npfrm")
	for i, ex := range examples {
		//if i != 3{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorRed,"example->",i+1,ColorYellow,ex,ColorReset)
		rezstring += fmt.Sprintf("example %v %s\n",i+1,ex)
		frmtyp, mod, err := JsonInp(fname)
		mod.Term = "dumb"
		//mod.Calc = true
		if err != nil{
			t.Errorf("non prismatic frame analysis test failed")
		}
		frmrez, err := CalcNp(mod, frmtyp, true)
		if err != nil{
			t.Errorf("non prismatic frame analysis test failed")
		}
		//dglb, _ := frmrez[2].([]float64)
		rnode, _ := frmrez[3].([]float64)
		//rezstring += fmt.Sprintf("global displacement -> %v\n",dglb)
		rezstring += fmt.Sprintf("nodal reactions  -> %.3f\n",rnode)
		//if mod.Calc{mod.CalcRezNp(frmtyp, frmrez)}
		
	}
	wantstring := `example 1 hulse2.6
nodal reactions  -> [74.324 234.911 40.765]
example 2 phil1
nodal reactions  -> [17.857 78.241 -16.099 72.411]
example 3 leet11.13
nodal reactions  -> [44.795 50.000 -238.872 -44.795 50.000 238.872]
example 4 phil4
nodal reactions  -> [4.730 44.488 -47.300 -4.730 44.488 47.300 -1.796 -4.488 1.796 -4.488]
`
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("non prismatic frame analysis test failed")
	}
}

//is worthless (for now)
func TestKFacToz(t *testing.T) {
	var rezstring string
	mem := &MemNp{
		Ts:[]int{2,0,2},
		Ls:[]float64{16,12,12},
		Ds:[]float64{3.6,1.2,3.2},
		Bs:[]float64{1.0,1.0,1.0},
		Styp:1,
		Dims:[]float64{1.0,1.2},
	}
	KFacToz(mem, false)
	rezstring += "toz toz"
	wantstring := ``
	if rezstring != wantstring {
		t.Errorf("stiffness and carry over factor test failed")
	}
}

func TestEndFrcNp(t *testing.T){
	mem := &MemNp{
		Ts:[]int{0,0,0},
		Ls:[]float64{1,1,1},
		Ds:[]float64{1.0,1.0,1.0,1.0},
		Bs:[]float64{1.0,1.0,1.0},
		Styp:1,
		Dims:[]float64{1.0,1.0},
		Frmtyp:"1db",
		Em:25e6,
		Lspan:3.0,
	}
	err := KFacMos(mem, false)
	if err != nil{
		log.Println(err)
	}
	qfchn := make(chan []interface{})
	vdx := false
	ncjt := 2
	member := 1
	lmap := map[int]string{
		1:"point load center",
		2:"point moment center",
		3:"udl",
		4:"trap udl",
		5:"point axial center",
		6:"dist axial",
		7:"torsional moment center",
		8:"tri left",
		9:"tri right",
	}
	load := map[int][]float64{
		1:{1,1,1,0,1.5,0},
		2:{1,2,1,0,1.5,0},
		3:{1,3,1,0,0,0},
		4:{1,4,1,2,0,0},
		5:{1,5,1,0,1.5,0},
		6:{1,6,1,0,0,0},
		7:{1,7,1,0,1.5,0},
		8:{1,4,1,0,0,0},
		9:{1,4,0,1,0,0},
	}
	for ltyp := 1; ltyp < 10; ltyp++{
		fmt.Println(ColorBlue,"load type-",lmap[ltyp],ColorReset)
		msloads := [][]float64{
			load[ltyp],
		}
		go FxdEndFrcNp(member, mem, msloads, ncjt, vdx, qfchn)
		r := <-qfchn
		qf, _ := r[1].([]float64)
		fmt.Println(ColorRed,qf,ColorReset)	
	}
}

func TestEndFrcNpTimo(t *testing.T){
	var rezstring string
	mem := &MemNp{
		Ts:[]int{2,0,2},
		Ls:[]float64{1,2,1},
		Ds:[]float64{1.7,1.0,1.7},
		Bs:[]float64{1.0,1.0,1.0},
		Styp:1,
		Dims:[]float64{1.0,1.0},
		Frmtyp:"1db",
		Em:25e6,
	}
	err := KFacMos(mem, false)
	if err != nil{
		log.Println(err)
	}
	rezstring += "timoshenko 9.4.1\n"
	ltyps := []string{"udl","point"}
	loads := [][][]float64{{{1,3,20,0,0,0}},{{1,1,20,0,2,0}}}
	qfchn := make(chan []interface{})
	for i, load := range loads{
		go FxdEndFrcNp(1, mem, load, 2, false, qfchn)
		r := <-qfchn
		qf, _ := r[1].([]float64)
		rezstring += fmt.Sprintf("load - %s end moments left %.4f right %.4f\n",ltyps[i],qf[1],qf[3])
		switch i{
			case 0:
			rezstring += fmt.Sprintf("ans- %f\n",0.1138 * load[0][2]*4.0*4.0)
			case 1:
			rezstring += fmt.Sprintf("ans- %f\n",0.1708 * load[0][2]*4.0)
		}
	}
	wantstring := ``
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("timoshenko end force np test failed")
	}
}
