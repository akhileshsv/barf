package barf

import (
	"os"
	"fmt"
	"path/filepath"
	"testing"
)

func TestFrmGen(t *testing.T){
	var examples = []string{"temp1"}
	//var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/gen/frame2d")
	for i, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		f, err := ReadFrm2d(fname)
		if err != nil{
			fmt.Println(err)
			t.Errorf("frame 2d test failed")
			return
		}
		err = f.Calc()
		if err != nil{
			fmt.Println(err)
			t.Errorf("frame 2d test failed")
			return
		}
		for lp := range f.Loadcons{f.DrawLp(lp, f.Term)}
	}
	//fmt.Println(rezstring)
	t.Errorf("frame 2d test failed")
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

/*
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
