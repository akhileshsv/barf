package barf

import (
	"os"
	"fmt"
	"path/filepath"
	"testing"
)

//sbsar beam design
func TestBmDzSar(t *testing.T){
	t.Log("testing sb sar beam design")
	//0,1,2,3,4,5,6,7,8,9,10,11,12,13,14
	var examples = []string{"sbsarj1","sbsarj2"}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/beam/")
	t.Log("starting steel beam design (is800) tests")
	for i, ex := range examples{
		//if i !=1{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		b, err := ReadBm(fname)
		if err != nil{
			t.Errorf("sbsar steel beam design (is800) test failed")
		}
		b.Onism = true
		err = BmDz(&b)
		if err != nil{
			t.Log(err)
			t.Errorf("sbsar steel beam design (is800) test failed")
		}
		rezstring += fmt.Sprintf("example %s\n",ex)
		for _, ss := range b.Ssecs{
			rezstring += fmt.Sprintf("%s vd %.2f kn vu %.2f kn mu %.2f knm mur %.2f knm def %.2f mm dperm %.2f mm fwc %.2f kn fwb %.2f kn\t\n",ss.Sstr,ss.Vd/1e3,ss.Vu/1e3,ss.Mu/1e6, ss.Mur/1e6, ss.Dmax/1.5, ss.Lspan/ss.Defrat, ss.Fwc/1e3, ss.Fwb/1e3)
		}
		rezstring += "\n"
	}
	wantstring := ``
	if rezstring != wantstring{
		t.Errorf("steel beam design (is800) failed")
		fmt.Println(rezstring)
	}
}

//rdarak beam design
func TestBmDzArk(t *testing.T){
	t.Log("testing rd ark beam design")
	//0,1,2,3,4,5,6,7,8,9,10,11,12,13,14
	var examples = []string{"rdarkj1","rdarkj2","rdarkj3",
		"rdarkg1","rdarkg2","rdarkg3","rdarkg4",
		"rdarkg5","rdarkg6","rdarkg7","rdarkg8",
		"rdarkg9","rdarkg10","rdarkg11","rdarkg12"}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/beam/")
	t.Log("starting steel beam design (is800) tests")
	for i, ex := range examples{
		//if i !=1{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		b, err := ReadBm(fname)
		if err != nil{
			t.Errorf("steel beam design (is800) test failed")
		}
		b.Onism = true
		err = BmDz(&b)
		if err != nil{
			t.Log(err)
			t.Errorf("steel beam design (is800) test failed")
		}
		rezstring += fmt.Sprintf("example %s\n",ex)
		for _, ss := range b.Ssecs{
			rezstring += fmt.Sprintf("%s vd %.2f kn vu %.2f kn mu %.2f knm mur %.2f knm def %.2f mm dperm %.2f mm fwc %.2f kn fwb %.2f kn\t\n",ss.Sstr,ss.Vd/1e3,ss.Vu/1e3,ss.Mu/1e6, ss.Mur/1e6, ss.Dmax/1.5, ss.Lspan/ss.Defrat, ss.Fwc/1e3, ss.Fwb/1e3)
		}
		rezstring += "\n"
	}
	wantstring := ``
	if rezstring != wantstring{
		t.Errorf("steel beam design (is800) failed")
		fmt.Println(rezstring)
	}
}

func TestBmDzBs(t *testing.T){
	t.Log("testing bs449 beam design")
	var examples = []string{"mosley6.2"}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/beam/")
	t.Log("starting steel beam design (bs449) tests")
	for i, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		b, err := ReadBm(fname)
		b.Drawrez = true
		b.Term = "svg"
		b.Nsecs = 5
		if err != nil{
			t.Errorf("steel beam design (bs449) test failed")
		}
		err = BmDz(&b)
		if err != nil{
			t.Log(err)
			t.Errorf("steel beam design (bs449) test failed")
		}
		for _, txtplt := range b.Txtplots{
			t.Log("txtplot",txtplt)
		}
		t.Log("report",b.Report)
		
		t.Log("txtplots",b.Txtplots)
		rezstring += fmt.Sprintf("example %s\n",ex)
		for _, ss := range b.Ssecs{
			rezstring += fmt.Sprintf("%s fm %.2f pm %.2f fs %.2f ps %.2f fwb %.2f pwb %.2f fwc %.2f pwc %.2f\t\n",ss.Sstr,ss.Fm,ss.Pm,ss.Fs, ss.Ps, ss.Fwb, ss.Pwb, ss.Fwc, ss.Pwc)
		}
		rezstring += "\n"
	}
	wantstring := ``
	if rezstring != wantstring{
		t.Errorf("steel beam design (bs449) failed")
		fmt.Println(rezstring)
	}	
}

func TestPrlnDz(t *testing.T){
	t.Log("testing purlin design")
	t.Log("duggal ex 9.9 i section")
	var tspc, pspc, theta, dl, ll, pwl float64
	sname := "i"
	nsecs := 5
	tspc = 6000.0
	pspc = 1500.0
	theta = 30.0
	dl = 130.0
	ll = 0.0
	pwl = 2000.0
	PrlnDz800(nsecs, sname, tspc, pspc, theta, dl, ll, pwl)
	t.Log("maity lec 55 i section")
	tspc = 5000.0
	dl = 120.0
	pspc = 2000.0
	pwl = 1500.0
	PrlnDz800(nsecs, sname, tspc, pspc, theta, dl, ll, pwl)
}


func TestBmDzIs(t *testing.T){
	t.Log("testing is800 beam design")
	var examples = []string{"duggal9.4","duggal9.5","duggal9.6","duggal9.8"}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/bash/beam/")
	t.Log("starting steel beam design (is800) tests")
	for i, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		b, err := ReadBm(fname)
		if err != nil{
			t.Errorf("steel beam design (is800) test failed")
		}
		err = BmDz(&b)
		if err != nil{
			t.Log(err)
			t.Errorf("steel beam design (is800) test failed")
		}
		rezstring += fmt.Sprintf("example %s\n",ex)
		for _, ss := range b.Ssecs{
			rezstring += fmt.Sprintf("%s vd %.2f kn vu %.2f kn mu %.2f knm mur %.2f knm def %.2f mm dperm %.2f mm fwc %.2f kn fwb %.2f kn\t\n",ss.Sstr,ss.Vd/1e3,ss.Vu/1e3,ss.Mu/1e6, ss.Mur/1e6, ss.Dmax/1.5, ss.Lspan/ss.Defrat, ss.Fwc/1e3, ss.Fwb/1e3)
		}
		rezstring += "\n"
	}
	wantstring := ``
	if rezstring != wantstring{
		t.Errorf("steel beam design (is800) failed")
		fmt.Println(rezstring)
	}
}

/*
func TestBmDesign(t *testing.T){
	t.Log("testing beam design redux")
	t.Log("duggal ex 9.4 ss beam")
	b := Bm{
		Lspan:4000.0,
		DL:15,
		Sname:"i",
		Nsecs:1,
		Dtyp:1,
		Code:1,
		Sdx:50,
		Lsb:true,
	}
	err := BmDz(&b)
	if err != nil{
		t.Fatal(err)
	}
	b.Table(true)
}

func TestStlBmDBs(t *testing.T) {
	//var rezstring string
	ldcases := [][]float64{
		{1,3,100,0,0,0,1},
		{1,3,160,0,2,1,2},
	}
	lspan := 6.0
	ly := 6.0
	ty := 1.0
	lbr := 200.0
	tbr := 20.0
	nsecs := 5
	grd := 43
	sname := "ub"
	brchck := true
	yeolde := true
	StlBmDBs(lspan, ly, ty, lbr, tbr, ldcases, sname, grd, nsecs, brchck, yeolde)	
}

*/

/*

// func TestBmDzIs(t *testing.T){
// 	t.Log("testing beam design redux")
// 	t.Log("duggal ex 9.4 ss beam")
// 	b := Bm{
// 		Lspan:4000.0,
// 		DL:15,
// 		Sname:"i",
// 		Nsecs:1,
// 		Dtyp:1,
// 		Code:1,
// 		Sdx:50,
// 		Lsb:true,
// 	}
// 	err := BmDz(&b)
// 	if err != nil{
// 		t.Fatal(err)
// 	}
// 	b.Table(true)
// 	t.Log("duggal ex 9.5 dtyp 0")
// 	b = Bm{
// 		Lspan:6000.0,
// 		Sname:"i",
// 		Nsecs:1,
// 		Dtyp:0,
// 		Code:1,
// 		Sdx:32,
// 		Lsb:true,
// 		Verbose:false,
// 		Vu:210*1e3,
// 		Mu:150*1e6,
// 		Dmax:0.001,
// 	}
// 	err = BmDz(&b)
// 	if err != nil{
// 		fmt.Println(err)
// 	}
	
	
// 	// b := Bm{
	// 	Ldcases:[][]float64{
	// 		{1,3,100,0,0,0,1},
	// 		{1,3,160,0,2,1,2},
	// 	},
	// 	Lspan:6.0,
	// 	Ly:6.0,
	// 	Ty:1.0,
	// 	Lbr:200.0,
	// 	Tbr:20.0,
	// 	Sname:"ub",
	// 	Nsecs:8,
	// 	Brchk:true,
	// 	Yeolde:true,
	// 	Endc:1,
	// 	Dtyp:1,
	// 	Grd:43,
	// 	Code:2,
	// // }
// 	// b.Spam = true
// 	// b.Verbose = true
// 	// err := BmDz(&b)
// 	// if err != nil{
	// 	fmt.Println(err)
	// // }
// }
   
//smash cafe secondary/primary beams
func TestSmshBmDz(t *testing.T) {
	//var rezstring string
	t.Log("smash cafe floor joist net udl 27 kn/m bs449")
	ldcases := [][]float64{
		{1,3,27,0,0,0,1},
	}
	lspan := 7.0
	ly := 7.0
	ty := 1.0
	lbr := 200.0
	tbr := 20.0
	nsecs := 5
	grd := 43
	sname := "ub"
	brchck := false
	yeolde := true
	StlBmDBs(lspan, ly, ty, lbr, tbr, ldcases, sname, grd, nsecs, brchck, yeolde)	

	t.Log("smash cafe primary beam point load 93 kn bs449")
	ldcases = [][]float64{
		{1,1,93,2.25,0,0,1},
	}
	lspan = 4.5
	ly = 4.5
	ty = 1.0
	nsecs = 5
	grd = 43
	sname = "ub"
	brchck = false
	yeolde = true
	StlBmDBs(lspan, ly, ty, lbr, tbr, ldcases, sname, grd, nsecs, brchck, yeolde)	

}

func TestBmDzIsSmshX(t *testing.T){
	t.Log("testing deckbmx non lsb clvr")
	b1 := Bm{
		Lspan:1800.0,
		Sname:"i",
		Nsecs:7,
		Dtyp:0,
		Code:1,
		Lsb:false,
		Mu:66e6,
		Vu:37e3,
		Dmax:0.5,
		Verbose:false,
	}
	err := BmDz(&b1)
	if err != nil{
		t.Fatal(err)
	}
	for i, ss := range b1.Ssecs{
		fmt.Printf("idx %v sdx %v sstr %s wt %.1f kn/m\n",i, ss.Sdx,ss.Sstr,ss.Wt)
	}
}

func TestBmDzIsSmshY(t *testing.T){
	t.Log("testing left beam non lsb")
	b1 := Bm{
		Lspan:4500.0,
		Sname:"i",
		Nsecs:3,
		Dtyp:1,
		Code:1,
		Lsb:false,
		Ldcases:[][]float64{{1,1,120e3,0,2250,0,1}},
		Verbose:false,
	}
	err := BmDz(&b1)
	if err != nil{
		t.Fatal(err)
	}
	for i, ss := range b1.Ssecs{
		fmt.Printf("idx %v sdx %v sstr %s wt %.1f kn/m\n",i, ss.Sdx,ss.Sstr,ss.Wt)
	}

	t.Log("testing right deck beam non lsb")
	b2 := Bm{
		Lspan:4500.0,
		Sname:"i",
		Nsecs:3,
		DL:4.1,
		LL:4.1,
		Dtyp:1,
		Code:1,
		Lsb:false,
		Verbose:false,
	}
	err = BmDz(&b2)
	if err != nil{
		t.Fatal(err)
	}
	for i, ss := range b2.Ssecs{
		fmt.Printf("idx %v sdx %v sstr %s wt %.1f kn/m\n",i, ss.Sdx,ss.Sstr,ss.Wt)
	}
	t.Log("testing right beam lsb")
	b3 := Bm{
		Lspan:4500.0,
		Sname:"i",
		Nsecs:3,
		DL:4.1,
		LL:4.1,
		Dtyp:1,
		Code:1,
		Lsb:false,
		Ldcases:[][]float64{{1,1,120e3,0,2250,0,1}},
		Verbose:false,
	}
	err = BmDz(&b3)
	if err != nil{
		t.Fatal(err)
	}
	for i, ss := range b3.Ssecs{
		fmt.Printf("idx %v sdx %v sstr %s wt %.1f kn/m\n",i, ss.Sdx,ss.Sstr,ss.Wt)
	}
}

func TestBmDzIsSmshJ(t *testing.T){
	t.Log("testing smsh floor joist lsb 20.25 kn/m ")

	b := Bm{
		Lspan:7000.0,
		DL:10.25,
		LL:10.25,
		Sname:"i",
		Nsecs:7,
		Dtyp:1,
		Code:1,
		Lsb:true,
	}
	err := BmDz(&b)
	if err != nil{
		t.Fatal(err)
	}
	for i, ss := range b.Ssecs{
		fmt.Printf("idx %v sdx %v sstr %s wt %.1f kn/m\n",i, ss.Sdx,ss.Sstr,ss.Wt)
	}

	t.Log("testing smsh end beam non lsb 10.125 kn/m ")

	b = Bm{
		Lspan:7000.0,
		DL:5.1,
		LL:5.1,
		Sname:"i",
		Nsecs:7,
		Dtyp:1,
		Code:1,
		Lsb:false,
	}
	err = BmDz(&b)
	if err != nil{
		t.Fatal(err)
	}
	for i, ss := range b.Ssecs{
		fmt.Printf("idx %v sdx %v sstr %s wt %.1f kn/m\n",i, ss.Sdx,ss.Sstr,ss.Wt)
	}
}

*/
