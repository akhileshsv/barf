package barf

import (
	"os"
	"fmt"
	"path/filepath"
	"testing"
	opt"barf/opt"
	kass"barf/kass"
)

func TestFrm2dGen(t *testing.T){
	var examples = []string{"raka1","raka2"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/f2d")
	var a opt.Bat
	var inp []interface{}
	ffnc := f2dobj
	//t.Log("random-",a,b,x,y)
	for i, ex := range examples{
		if i >1{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Logf("starting ex no. %v - %s\n",i+1, ex)
		f, err := kass.ReadFrm2d(fname)
		if err != nil{
			t.Log(err)
			t.Errorf("frame 2d gen (or bust) test failed")
		}
		inp = append(inp, f)
		for j := 1; j < 4; j++{
			if j != 2{continue}
			f.Ngrp = j
			f.Bfcalc = true
			t.Log("setting ng to",f.Ngrp)
			ndims := f.Getndim()
			t.Log("ndims-",ndims)		
			var nd []int
			for k := 0; k < ndims; k++{
				nd = append(nd, 21)
			}
			t.Log("ndvec-",nd)
			err = a.Init(nd, inp, ffnc)
			if err != nil{
				t.Log(err)
				t.Errorf("frame 2d gen (or bust) test failed")
			}
			t.Log("kost-",f.Kost,a.Fit)
			t.Log("posvec-",a.Pos)
		}
		
	}
}

func TestOptFrm2dIEC(t *testing.T){
	var examples = []string{"raka1","raka2"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/f2d")
	opts := []int{1,2}
	namez := []string{"ga","pso"}
	NCycles = 3
	for i, ex := range examples{
		for j, opt := range opts{
			
			fname := filepath.Join(datadir,ex+".json")
			t.Logf("starting ex no. %v - %s opt - %s\n",i+1, ex, namez[j])
			f, err := kass.ReadFrm2d(fname)
			if err != nil{
				t.Log(err)
				t.Errorf("frame 2d optimization (IEC) test failed")
			}
			f.Width = 225.0
			//f.Bfcalc = true
			f.Term = "svg"
			f.Web = true
			f.Opt = opt
			f.Ngrp = 1
			frez, err := Frm2dOpt(f)
			if err != nil{
				t.Log(ColorRed,"ERRORE-",err,ColorReset)
			}
			
			t.Logf("optimized cost - %.4f rs",frez.Kost)
			t.Logf("volume of  rcc %.4f m3 weight of steel %.2f kg area of formwork %.3f m2\n",frez.Quants[0],frez.Quants[1],frez.Quants[2])
			
			t.Logf("optimized section vec - %.f rs",frez.Sections)
			
		}
	}
}


func TestFrm2dOpt(t *testing.T){
	//fname := "raka2.json"
	
	fname := "raka1.json"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/f2d")
	filename := filepath.Join(datadir, fname)
	f, err := kass.ReadFrm2d(filename)
	if err != nil{
		t.Log(err)
		t.Errorf("frame 2d test failed")
	}
	f.Term = "svg"
	f.Web = true
	f.Opt = 1
	f.Ngrp = 1
	var mqs []float64
	var mkost float64
	var msecs [][]float64
	mkost = 1e19
	for i := 0; i < 10; i++{		
		frez, err := Frm2dOpt(f)
		if err != nil{
			t.Fatal(err)
		}
		// for i, txtplot := range frez.Txtplots{
		// 	fmt.Println("i, txtplot ->",i, txtplot)
		// }
		t.Logf("optimized cost in run %v - %.4f rs",i+1,frez.Kost)
		t.Logf("volume of  rcc %.4f m3 weight of steel %.2f kg area of formwork %.3f m2\n",frez.Quants[0],frez.Quants[1],frez.Quants[2])
		t.Logf("optimized section vec - %.f rs",frez.Sections)
		
		if frez.Kost < mkost{
			mkost = frez.Kost
			msecs = frez.Sections
			mqs = frez.Quants
		}
	}
	fmt.Println("mkost",mkost)
	fmt.Println("msecs",msecs)
	fmt.Println("mqs",mqs)
}

func TestFrmInit(t *testing.T){	
	cvec := [][]float64{{200,250}}
	bvec := [][]float64{{200,300},{250, 350}}
	df := 0.0
	xs := []float64{0,6,10,14}
	code := 1
	//bsts := []int{1}
	bstyp := 1
	//fgs := []int{1}
	fgen := 1
	nbms := len(xs)-1; ncols := len(xs)
	brel := 0
	cdims := kass.InitCdims(cvec, fgen, ncols)
	bdims := kass.InitBdims(df, xs, bvec , code, fgen, bstyp, nbms, brel)
	fmt.Println(cdims)
	fmt.Println(bdims)
}

func TestFrm2dRaka(t *testing.T){
	fname := "raka1_sol1.json"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/f2d")
	filename := filepath.Join(datadir, fname)
	f, err := kass.ReadFrm2d(filename)
	if err != nil{
		t.Log(err)
		t.Errorf("frame 2d test failed")
	}
	//fmt.Println(f)
	allcons, _, err := Frm2dDz(&f)
	fmt.Println("error->",err, allcons)
	//for i := range emap{
	//	fmt.Println("i->",i, "emap->",emap[i])
	//}
	
	for lp := range f.Loadcons{
		//fmt.Println(lp)
		f.DrawLp(lp, "svg")
	}
}

func TestFrm2dDz(t *testing.T){
	fname := "allen3.2.json"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/f2d")
	filename := filepath.Join(datadir, fname)
	f, err := kass.ReadFrm2d(filename)
	if err != nil{
		t.Log(err)
		t.Errorf("frame 2d test failed")
	}
	//fmt.Println(f)
	f.Web = true
	f.Noprnt = true
	Frm2dDz(&f)
	for _, t := range f.Txtplots{
		fmt.Println(t)
	}
}

func TestFrameDzMosley(t *testing.T){
	fname := "frm2d_mosley7.5.json"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	filename := filepath.Join(datadir, fname)
	f, err := kass.ReadFrm2d(filename)
	if err != nil{
		t.Log(err)
		t.Errorf("frame 2d test failed")
	}
	Frm2dDz(&f)
}


/*
   
func TestFrm2dRaka(t *testing.T){
	var examples = []string{"raka1","raka2"}
	var rezstring string
	var err error
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/gen/frame2d/")
	//for i, ex := range examples {
	//	fname := filepath.Join(datadir,ex+".json")
	//}
	//f, err := kass.ReadFrm2d(fname)
	if err != nil{
		fmt.Println(err)
		t.Errorf("2d frm opt test failed")
	}
	fmt.Println(f)
	//Frm2dOpt(f)
}

func TestSelfWeight(t *testing.T){
	fname := "frm2d_selfwt.json"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	filename := filepath.Join(datadir, fname)
	f, err := ReadFrm2d(filename)
	if err != nil{
		t.Log(err)
		t.Errorf("frame 2d test failed")
	}
	err = f.Init()
	if err != nil{
		fmt.Println(err)
	} 
	f.GenCoords()
	f.GenMprp()
	f.GenLoads()
	f.InitMemRez()
	f.CalcLoadEnv()
	f.DrawMod("dumb")
	fmt.Println(f.Advloads)
	fmt.Println(f.Txtplots[0])
}

func TestInitSections(t *testing.T){
	bvec := [][]float64{{200,300},{200,400},{200,600},{200,300}}
	df := 125.0
	xs := []float64{0,6,10,16,20}
	code := 1
	bsts := []int{1,6,7}
	fgs := []int{1,2,3,4}
	nbms := 4
	brel := 0
	for _, bstyp := range bsts{
		for _, fgen := range fgs{
			fmt.Println(bstyp, fgen)
			bdims := InitBdims(df, xs, bvec , code, fgen, bstyp, nbms, brel)
			fmt.Println(bdims)
		}

	}
}

func TestFrame2dRajeev2(t *testing.T){
	fname := "frm2d_rajeev2.json"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	filename := filepath.Join(datadir, fname)
	f, err := ReadFrm2d(filename)
	if err != nil{
		t.Log(err)
		t.Errorf("frame 2d test failed")
	}
	err = CalcFrm2d(&f, f.Term)
	if err != nil{fmt.Println(err)}
	//for _, txtplot := range f.Txtplots{
	//	fmt.Println(txtplot)
	//}
	var rezstring string
	rezstring += fmt.Sprintf("%s\n",f.Txtplots[0])
	rezstring += "column envelopes->\n"
	for id := 1; id <= len(f.X) * f.Nflrs; id++{
		rezstring += fmt.Sprintf("member id %v floor level %v\n",id,(id-1)/len(f.X))
		c := f.Colenv[id]
		rezstring += fmt.Sprintf("max top moment %.3f max bottom moment %.3f\n",c.Mtmax,c.Mbmax)
		rezstring += fmt.Sprintf("max axial load %.3f\n", c.Pumax)
	}
	fmt.Println(rezstring)
	//fmt.Println("beams->",f.Beams, f.CBeams)
	//for i, report := range f.Reports{
	//	fmt.Println("load pattern->",i)
	//	fmt.Println(report)
	//}
	t.Errorf("frame 2d test failed")
}


func TestFrame2dRajeev1(t *testing.T){
	fname := "frm2d_rajeev1.json"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	filename := filepath.Join(datadir, fname)
	f, err := ReadFrm2d(filename)
	if err != nil{
		t.Log(err)
		t.Errorf("frame 2d test failed")
	}
	err = CalcFrm2d(&f, f.Term)
	if err != nil{fmt.Println(err)}
	//for _, txtplot := range f.Txtplots{
	//	fmt.Println(txtplot)
	//}
	var rezstring string
	rezstring += fmt.Sprintf("%s\n",f.Txtplots[0])
	rezstring += "column envelopes->\n"
	for id := 1; id <= len(f.X) * f.Nflrs; id++{
		rezstring += fmt.Sprintf("member id %v floor level %v\n",id,(id-1)/len(f.X))
		c := f.Colenv[id]
		rezstring += fmt.Sprintf("max top moment %.3f max bottom moment %.3f\n",c.Mtmax,c.Mbmax)
		rezstring += fmt.Sprintf("max axial load %.3f\n", c.Pumax)
	}
	fmt.Println(rezstring)
	//fmt.Println("beams->",f.Beams, f.CBeams)
	//for i, report := range f.Reports{
	//	fmt.Println("load pattern->",i)
	//	fmt.Println(report)
	//}
	t.Errorf("frame 2d test failed")
}


func TestFrame2dMosley92(t *testing.T){
	fname := "frm2d_mosley9.2.json"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	filename := filepath.Join(datadir, fname)
	f, err := ReadFrm2d(filename)
	if err != nil{
		t.Log(err)
		t.Errorf("frame 2d test failed")
	}
	err = CalcFrm2d(&f, f.Term)
	if err != nil{fmt.Println(err)}
	//for _, txtplot := range f.Txtplots{
	//	fmt.Println(txtplot)
	//}
	var rezstring string
	rezstring += fmt.Sprintf("%s\n",f.Txtplots[0])
	rezstring += "column envelopes->\n"
	for id := 1; id <= len(f.X) * f.Nflrs; id++{
		rezstring += fmt.Sprintf("member id %v floor level %v\n",id,(id-1)/len(f.X))
		c := f.Colenv[id]
		rezstring += fmt.Sprintf("max top moment %.3f max bottom moment %.3f\n",c.Mtmax,c.Mbmax)
		rezstring += fmt.Sprintf("max axial load %.3f\n", c.Pumax)
	}
	fmt.Println(rezstring)
	fmt.Println("beams->",f.Beams, f.CBeams)
	for i, report := range f.Reports{
		fmt.Println("load pattern->",i)
		fmt.Println(report)
	}
	t.Errorf("frame 2d test failed")
}


func TestFrame2dMosley33(t *testing.T){
	fname := "frm2d_mosley3.3.json"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	filename := filepath.Join(datadir, fname)
	f, err := ReadFrm2d(filename)
	if err != nil{
		t.Log(err)
		t.Errorf("frame 2d test failed")
	}
	err = CalcFrm2d(&f, f.Term)
	if err != nil{fmt.Println(err)}
	for _, txtplot := range f.Txtplots{
		fmt.Println(txtplot)
	}
	var rezstring string
	rezstring += "column envelopes->\n"
	for id := 1; id <= len(f.X) * f.Nflrs; id++{
		rezstring += fmt.Sprintf("member id %v floor level %v\n",id,(id-1)/len(f.X))
		c := f.Colenv[id]
		rezstring += fmt.Sprintf("max top moment %.3f max bottom moment %.3f\n",c.Mtmax,c.Mbmax)
		rezstring += fmt.Sprintf("max axial load %.3f\n", c.Pumax)
	}
	fmt.Println(rezstring)
	t.Errorf("frame 2d test failed")
}

func TestFrame2dAllen11(t *testing.T){
	fname := "frm2d_allen1.1.json"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	filename := filepath.Join(datadir, fname)
	f, err := ReadFrm2d(filename)
	if err != nil{
		t.Log(err)
		t.Errorf("frame 2d test failed")
	}
	err = CalcFrm2d(&f, f.Term)
	if err != nil{fmt.Println(err)}
	for _, txtplot := range f.Txtplots{
		fmt.Println(txtplot)
	}
	var rezstring string
	rezstring += "column envelopes->\n"
	for id := 1; id <= len(f.X) * f.Nflrs; id++{
		rezstring += fmt.Sprintf("member id %v floor level %v\n",id,(id-1)/len(f.X))
		c := f.Colenv[id]
		rezstring += fmt.Sprintf("max top moment %.3f max bottom moment %.3f\n",c.Mtmax,c.Mbmax)
		rezstring += fmt.Sprintf("max axial load %.3f\n", c.Pumax)
	}
	fmt.Println(rezstring)
	t.Errorf("frame 2d test failed")
}
*/
