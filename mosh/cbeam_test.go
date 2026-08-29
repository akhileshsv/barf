package barf

import(
	"fmt"
	"testing"
	"path/filepath"
	"os"
	//"strings"
	//"log"
	kass"barf/kass"	
)

func TestCalcBf(t *testing.T){
	var code, nspans, mrel, bstyp int
	var bdx, cdx, styps []int
	var sections [][]float64
	var lspans []float64
	var dslb float64
	code = 1; nspans = 3; mrel = 0; bstyp = 0
	bdx = []int{3,2,3}
	cdx = []int{1,2,2,1}
	sections = [][]float64{{300,500},{300,700},{200,300}}
	styps = []int{1,1,7}
	lspans = []float64{7.0,5.0,4.0}
	dslb = 125.0

	secvec, sts := kass.CalcBf(code, bstyp, mrel, nspans, dslb , cdx, bdx, styps, lspans, sections)
	
	for i := range secvec{
		if i < nspans + 1{fmt.Println("col sec")} else {fmt.Println("beam sec")}
		fmt.Println("section->",i+1,"dims->",secvec[i],"styp",sts[i])
	}
	cdx = []int{}
	secvec, sts = kass.CalcBf(code, bstyp, mrel, nspans, dslb , cdx, bdx, styps, lspans, sections)
	fmt.Println("cbeam->")
	for i := range secvec{
		fmt.Println("section->",i+1,"dims->",secvec[i],"styp",sts[i])
	}
	
}

func TestSsBmDz(t *testing.T){
	//single span beam design tests
	var rezstring string
	var examples = []string{"shah5.6.1","hulse7.6","govr1"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/cbeam")
	var bmenv map[int]*kass.BmEnv
	for i, ex := range examples{
		if i > 0{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		cb, err := ReadCBm(fname)
		if err != nil{
			t.Errorf("ss beam design test failed")
		}
		
		cb.Term = "dumb"
		bmenv, err = CBeamEnvRcc(&cb, cb.Term, true)
		if err != nil{
			fmt.Println(err)
			break
			t.Errorf("ss beam design test failed")
		}
		cb.Term = "dumb"
		_, err = CBmDz(&cb,bmenv)
		rezstring += fmt.Sprintf("\nexample-%s\n",ex)
		if err != nil{
			fmt.Println(err)
			break
		}
		cb.Term = "qt"
		_ = PlotCBmDet(cb.Web, cb.Bmvec, cb.RcBm, cb.Foldr, cb.Title, cb.Term)
		
		for j, span := range cb.RcBm{
			rezstring += fmt.Sprintf("span-%v\n",j+1)
			for _, bm := range span{
				if bm.Ignore{continue}
				rezstring += bm.Printz()
				bm.Table(true)
				//t.Log(bm.Report)
				//t.Log(bm.Report)
			}
		}
		//for _, txtplot := range cb.Txtplots{
		//	fmt.Println(txtplot)
		//}

	}
	wantstring := ``
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("ss beam design test failed")		
	}
}

func TestCBmDz(t *testing.T){
	var rezstring string
	var examples = []string{"hulse7.9","allen3.1"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/cbeam")
	//var bmenv map[int]*kass.BmEnv
	for i, ex := range examples{
		if i != 1{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		cb, err := ReadCBm(fname)
		if err != nil{
			t.Errorf("cbeam design test failed")
		}
		/*
		bmenv, err = CBeamEnvRcc(&cb, cb.Term, true)
		if err != nil{
			fmt.Println(err)
			t.Errorf("cbeam design test failed")
		}
		CBmDz(&cb,bmenv)
		*/
		CalcCBm(cb)
		rezstring += fmt.Sprintf("\nexample-%s\n",ex)
		for j, span := range cb.RcBm{
			rezstring += fmt.Sprintf("span-%v\n",j+1)
			for _, bm := range span{
				if bm.Ignore{continue}
				rezstring += bm.Printz()
				bm.Table(true)
				PlotBmGeom(bm,cb.Term)
				//t.Log(bm.Report)
				//t.Log(bm.Report)
			}
		}
	}
	wantstring := ``
	if rezstring != wantstring{
		//fmt.Println(rezstring)
		t.Errorf("cbeam design test failed")		
	}

}

func TestCBeamClvr(t *testing.T){
	cb := &CBm{
		Fck:25.0,
		DL:5.0,
		LL:2.0,
		Nspans:1,
		Lspans:[]float64{4.0},
		Clvrs:[][]float64{{0.8,0,0},{0.8,0,0}},
		PSFs:[]float64{1.5,1.0,1.5,0.0},
		DM:0.0,
		Em:[][]float64{{25e6}},
		Sections:[][]float64{{1000,130}},
		//Cp:[][]float64{{1.1,1.1}},
	}
	termstr := "dumb"
	var report bool
	_, err := CBeamEnvRcc(cb, termstr, report)
	if err != nil {
		fmt.Println(err)
	}
	for i, plt := range cb.Txtplots{
		fmt.Println("LOAD PATTERN->",i)
		fmt.Println(plt)
	}
}

func TestCBeamEnvShah(t *testing.T){
	/*
	   illustrated dz shah example 7.2.3
	*/
	cb := &CBm{
		Fck:25.0,
		DL:5.0,
		LL:1.5,
		Nspans:4,
		Lspans:[]float64{4.0,4.0,4.0,4.0},
		Clvrs:[][]float64{{0,0,0},{0,0,0}},
		PSFs:[]float64{1.5,1.0,1.5,0.0},
		DM:0.3,
		Em:[][]float64{{25e6}},
		Sections:[][]float64{{1000,130}},
		//Cp:[][]float64{{1.1,1.1}},
	}
	termstr := "dumb"
	var report bool
	_, err := CBeamEnvRcc(cb, termstr, report)
	if err != nil {
		fmt.Println(err)
	}
	for i, plt := range cb.Txtplots{
		fmt.Println("LOAD PATTERN->",i)
		fmt.Println(plt)
	}
	/*
	if cb.DM == 0.0{
		for i := 1; i <= cb.Nspans; i++{
			fmt.Println(ColorReset)
			fmt.Println("member-",i)
			
			for j, vx := range bmenv[i].Venv{
				fmt.Println(ColorCyan)
				fmt.Printf("section\t- %v shear %.1f hogging BM %.1f sagging BM %.1f\n",j, vx, bmenv[i].Mnenv[j], bmenv[i].Mpenv[j])
				//fmt.Println(ColorYellow)
				//fmt.Printf("redistributed\t- %v shear %.1f hogging BM %.1f sagging BM %.1f\n",j, bmenv[i].Vrd[j], bmenv[i].Mnrd[j], bmenv[i].Mprd[j])
				fmt.Println(ColorReset)
			}
			fmt.Println("mpmax, mnmax, vmax->",bmenv[i].Mpmax,bmenv[i].Mnmax,bmenv[i].Vmax)
			//fmt.Println("ml, mr, vl, vr->",bmenv[i].Ml,bmenv[i].Mr,bmenv[i].Vl,bmenv[i].Vr)
		}
	}
	*/
}

func TestCBeamEnvAllen(t *testing.T){
	//first check with ex hulse 2.3
	cb := &CBm{
		Fck:25.0,
		DL:23.2,
		LL:20.0,
		Nspans:3,
		Lspans:[]float64{10.0,10.0,10.0},
		Clvrs:[][]float64{{0,0,0},{0,0,0}},
		PSFs:[]float64{1.4,1.0,1.6,0.0},
		DM:0.3,
		Em:[][]float64{{25e6}},
		Sections:[][]float64{{300,600}},
		//Cp:[][]float64{{1.1,1.1}},
	}
	termstr := "dumb"
	var report bool
	_, err := CBeamEnvRcc(cb, termstr, report)
	if err != nil {
		fmt.Println(err)
	}
	for i, plot := range cb.Txtplots{
		fmt.Println("XXX---lp->",i,"---XXX")
		fmt.Println(plot)
	}
}

func TestCBeamEnvHulse(t *testing.T){
	//first check with ex hulse 2.3
	cb := &CBm{
		Fck:25.0,
		DL:25.0,
		LL:10.0,
		Nspans:3,
		Lspans:[]float64{6.0,4.0,6.0},
		Clvrs:[][]float64{{0,0,0},{0,0,0}},
		PSFs:[]float64{1.35,1.0,1.5,0.0},
		DM:0.20,
		Sections:[][]float64{{230,400}},
	}
	termstr := "dumb"
	var report bool
	bmenv, err := CBeamEnvRcc(cb, termstr, report)
	if err != nil {
		fmt.Println(err)
	}
	
	for i, plot := range cb.Txtplots{
		fmt.Println("load pattern->",i)
		fmt.Println(plot)
	}
	for i := 1; i <= cb.Nspans; i++{
		fmt.Println("member-",i)

		for j, vx := range bmenv[i].Venv{
			fmt.Println(ColorCyan)
			fmt.Printf("section - %v shear %.1f hogging BM %.1f sagging BM %.1f\n",j, vx, bmenv[i].Mnenv[j], bmenv[i].Mpenv[j])
			fmt.Println(ColorReset)
		}
		fmt.Println("mpmax, mnmax, vmax->",bmenv[i].Mpmax,bmenv[i].Mnmax,bmenv[i].Vmax)
		fmt.Println("ml, mr, vl, vr->",bmenv[i].Ml,bmenv[i].Mr,bmenv[i].Vl,bmenv[i].Vr)
	}
}

func TestBeamSecCp(t *testing.T){
	var dims, cpvec []float64
	var styp, endc, ldcalc, nd int
	var wdl, lspan, efcvr float64
	var err error
	dims = []float64{230,300}
	styp = 1; endc = 1; ldcalc = 1; nd = 1
	lspan = 4.0
	dims, cpvec, wdl, err = BeamSecGen(styp,endc, ldcalc, nd, dims, lspan, efcvr)
	
	fmt.Println(wdl, "kn/m",dims, cpvec, err)

	dims = []float64{230}
	styp = 1; endc = 1; ldcalc = 1
	lspan = 4.0; nd = 1
	dims, cpvec, wdl, err = BeamSecGen(styp,endc, ldcalc, nd, dims, lspan, efcvr)

	fmt.Println(wdl, "kn/m",dims, cpvec, err)

	styp = 1; endc = 2; ldcalc = 1
	lspan = 4.0; nd = 1

	dims = nil
	dims, cpvec, wdl, err = BeamSecGen(styp,endc, ldcalc, nd, dims, lspan, efcvr)
	fmt.Println(wdl, "kn/m",dims, cpvec, err)

}
/*
func TestBmSsMosley(t *testing.T){
	var examples = []string{"mosley4.1"}
	var rezstring string
	var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for _, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := kass.JsonInp(fname)
		frmrez, _ = kass.CalcBm1d(mod, 2)
		bmresults := CBeamBmSf(mod, frmrez,false)
		rezstring += fmt.Sprintf("%v\n",ex)
		for i:= 1; i <= len(bmresults); i++{
			bm := bmresults[i]
			rezstring += fmt.Sprintf("member %v\n",bm.Mem)
			rezstring += fmt.Sprintf("%.2f\n",bm.SF)
			rezstring += fmt.Sprintf("%.2f\n",bm.BM)
			rezstring += fmt.Sprintf("%.2f\n",bm.Dxs)
			rezstring += fmt.Sprintf("\nMax SF %.3f at %.3f Max BM %.3f at %.3f Max def %.3f at %.3f\n",bm.Maxs[0],bm.Locs[0],bm.Maxs[1],bm.Locs[1],bm.Maxs[2],bm.Locs[2])
		}
	}
	wantstring := `
mosley4.1
member 1
[60.60 60.60 60.60 50.60 28.60 18.60 8.60 -1.40 -11.40 -11.40 -11.40 -11.40 -12.07 -14.07 -17.40 -22.07 -28.07 -35.40 -35.40 -35.40 -35.40]
[0.00 60.60 121.20 176.80 222.40 246.00 259.60 263.20 256.80 245.40 234.00 222.60 209.87 194.47 176.40 155.67 132.27 -205.80 -217.20 -228.60 -240.00]
[0.00 0.00 0.00 0.01 0.01 0.01 0.01 0.01 0.01 0.01 0.01 0.01 0.01 0.01 0.01 0.01 0.00 0.00 0.00 0.00 0.00]

Max SF 60.600 at 0.000 Max BM 263.200 at 7.000 Max def 9.974 at 9.000
`
	t.Logf(rezstring)
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("Bending Moment and SF analysis test failed")
	}	
	
}


func TestCBeamRccDBs(t *testing.T) {
	var examples = []string{"hulse2.3"}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for _, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := kass.JsonInp(fname)
		loadenvz, _ := CBeamRccDBs(mod) 
		for i := 1; i <= len(loadenvz); i++{
			bm := loadenvz[i]
			rezstring += fmt.Sprintf("span %v\n",bm.Mem)
			//rezstring += "\tshear -kn \t moment-hog kn-m \t moment-sag kn-m\n"
			for j := 0; j < 21; j++ {
				rezstring += fmt.Sprintf("sec %v\tshear %.3f\tmoment-hog %.3f\tmoment-sag %.3f\n",j+1,bm.Venv[j],bm.Mnenv[j],bm.Mpenv[j])				
			}  
		}
	}
	wantstring := `span 1
sec 1	shear 131.097	moment-hog -0.000	moment-sag 0.000
sec 2	shear 115.797	moment-hog 0.000	moment-sag 37.034
sec 3	shear 100.497	moment-hog 0.000	moment-sag 69.478
sec 4	shear 85.197	moment-hog 0.000	moment-sag 97.332
sec 5	shear 69.897	moment-hog 0.000	moment-sag 120.597
sec 6	shear 54.597	moment-hog 0.000	moment-sag 139.271
sec 7	shear 39.297	moment-hog 0.000	moment-sag 153.355
sec 8	shear 23.997	moment-hog 0.000	moment-sag 162.849
sec 9	shear 8.697	moment-hog 0.000	moment-sag 167.753
sec 10	shear -11.929	moment-hog 0.000	moment-sag 168.067
sec 11	shear -27.229	moment-hog 0.000	moment-sag 163.792
sec 12	shear -42.529	moment-hog 0.000	moment-sag 154.926
sec 13	shear -57.829	moment-hog 0.000	moment-sag 141.470
sec 14	shear -73.129	moment-hog 0.000	moment-sag 123.424
sec 15	shear -88.429	moment-hog 0.000	moment-sag 100.788
sec 16	shear -103.729	moment-hog 0.000	moment-sag 73.562
sec 17	shear -119.029	moment-hog -0.200	moment-sag 41.747
sec 18	shear -134.329	moment-hog -21.824	moment-sag 5.341
sec 19	shear -149.629	moment-hog -64.417	moment-sag 0.000
sec 20	shear -164.929	moment-hog -111.601	moment-sag 0.000
sec 21	shear -180.229	moment-hog -163.375	moment-sag 0.000
span 2
sec 1	shear 123.938	moment-hog -163.375	moment-sag 0.000
sec 2	shear 113.737	moment-hog -139.607	moment-sag 0.000
sec 3	shear 103.537	moment-hog -117.880	moment-sag 0.000
sec 4	shear 93.338	moment-hog -105.917	moment-sag 0.000
sec 5	shear 83.138	moment-hog -99.417	moment-sag 0.000
sec 6	shear 72.938	moment-hog -93.917	moment-sag 0.000
sec 7	shear 62.737	moment-hog -89.417	moment-sag 0.000
sec 8	shear 52.537	moment-hog -85.917	moment-sag 2.570
sec 9	shear 42.337	moment-hog -83.417	moment-sag 7.670
sec 10	shear 32.138	moment-hog -81.917	moment-sag 10.730
sec 11	shear 21.938	moment-hog -81.417	moment-sag 11.750
sec 12	shear -32.138	moment-hog -81.917	moment-sag 10.730
sec 13	shear -42.338	moment-hog -83.417	moment-sag 7.670
sec 14	shear -52.538	moment-hog -85.917	moment-sag 2.570
sec 15	shear -62.738	moment-hog -89.417	moment-sag 0.000
sec 16	shear -72.938	moment-hog -93.917	moment-sag 0.000
sec 17	shear -83.138	moment-hog -99.417	moment-sag 0.000
sec 18	shear -93.338	moment-hog -105.917	moment-sag 0.000
sec 19	shear -103.537	moment-hog -117.880	moment-sag 0.000
sec 20	shear -113.738	moment-hog -139.608	moment-sag 0.000
sec 21	shear -123.938	moment-hog -163.375	moment-sag 0.000
span 3
sec 1	shear 180.229	moment-hog -163.375	moment-sag 0.000
sec 2	shear 164.929	moment-hog -111.601	moment-sag 0.000
sec 3	shear 149.629	moment-hog -64.418	moment-sag 0.000
sec 4	shear 134.329	moment-hog -21.824	moment-sag 5.341
sec 5	shear 119.029	moment-hog -0.200	moment-sag 41.747
sec 6	shear 103.729	moment-hog 0.000	moment-sag 73.562
sec 7	shear 88.429	moment-hog 0.000	moment-sag 100.788
sec 8	shear 73.129	moment-hog 0.000	moment-sag 123.424
sec 9	shear 57.829	moment-hog 0.000	moment-sag 141.470
sec 10	shear 42.529	moment-hog 0.000	moment-sag 154.926
sec 11	shear 27.229	moment-hog 0.000	moment-sag 163.792
sec 12	shear 11.929	moment-hog 0.000	moment-sag 168.067
sec 13	shear -8.697	moment-hog 0.000	moment-sag 167.753
sec 14	shear -23.997	moment-hog 0.000	moment-sag 162.849
sec 15	shear -39.297	moment-hog 0.000	moment-sag 153.355
sec 16	shear -54.597	moment-hog 0.000	moment-sag 139.271
sec 17	shear -69.897	moment-hog 0.000	moment-sag 120.597
sec 18	shear -85.197	moment-hog 0.000	moment-sag 97.333
sec 19	shear -100.497	moment-hog 0.000	moment-sag 69.478
sec 20	shear -115.797	moment-hog 0.000	moment-sag 37.034
sec 21	shear -131.097	moment-hog -0.000	moment-sag 0.000
`
	t.Logf(rezstring)
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("Bending Moment and SF envelope test failed")
	}	
}


func TestCBeamBmSf(t *testing.T) {
	var examples = []string{"hulse2.1","hulse2.2"}
	var rezstring string
	var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for _, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := kass.JsonInp(fname)
		frmrez, _ = kass.CalcBm1d(mod, 2)
		bmresults := CBeamBmSf(mod, frmrez,false)
		rezstring += fmt.Sprintf("%v\n",ex)
		for i:= 1; i <= len(bmresults); i++{
			bm := bmresults[i]
			rezstring += fmt.Sprintf("member %v\n",bm.Mem)
			rezstring += fmt.Sprintf("%.2f",bm.SF)
			rezstring += fmt.Sprintf("%.2f",bm.BM)
			rezstring += fmt.Sprintf("%.2f",bm.Dxs)
			rezstring += fmt.Sprintf("Max SF %.3f at %.3f Max BM %.3f at %.3f Max def %.3f at %.3f\n",bm.Maxs[0],bm.Locs[0],bm.Maxs[1],bm.Locs[1],bm.Maxs[2],bm.Locs[2])
		}
	}
	wantstring := `
hulse2.1
member 1
[14.00 14.00 14.00 14.00 14.00 14.00 14.00 14.00 4.00 1.50 -1.00 -3.50 -6.00 -8.50 -11.00 -13.50 -16.00 -16.00 -16.00 -16.00 -16.00][-0.00 14.00 28.00 42.00 56.00 70.00 84.00 98.00 112.00 114.75 115.00 112.75 108.00 100.75 91.00 78.75 64.00 48.00 32.00 16.00 0.00][0.00 0.00 0.01 0.01 0.02 0.02 0.02 0.02 0.03 0.03 0.03 0.03 0.03 0.02 0.02 0.02 0.02 0.01 0.01 0.00 0.00]Max SF -16.000 at 16.000 Max BM 115.000 at 10.000 Max def 26.614 at 10.000
hulse2.2
member 1
[131.10 115.80 100.50 85.20 69.90 54.60 39.30 24.00 8.70 -6.60 -21.90 -37.20 -52.50 -67.80 -83.10 -98.40 -113.70 -129.00 -144.30 -159.60 -174.90][-0.00 37.03 69.48 97.33 120.60 139.27 153.35 162.85 167.75 168.07 163.79 154.93 141.47 123.42 100.79 73.56 41.75 5.34 -35.65 -81.24 -131.42][0.00 0.00 0.01 0.01 0.01 0.01 0.02 0.02 0.02 0.02 0.02 0.02 0.02 0.02 0.01 0.01 0.01 0.01 0.00 0.00 0.00]Max SF -174.903 at 6.000 Max BM 168.068 at 2.700 Max def 18.996 at 2.700
member 2
[50.00 45.00 40.00 35.00 30.00 25.00 20.00 15.00 10.00 5.00 0.00 -5.00 -10.00 -15.00 -20.00 -25.00 -30.00 -35.00 -40.00 -45.00 -50.00][-131.42 -121.92 -113.42 -105.92 -99.42 -93.92 -89.42 -85.92 -83.42 -81.92 -81.42 -81.92 -83.42 -85.92 -89.42 -93.92 -99.42 -105.92 -113.42 -121.92 -131.42][0.00 -0.00 -0.00 -0.00 -0.00 -0.00 -0.01 -0.01 -0.01 -0.01 -0.01 -0.01 -0.01 -0.01 -0.01 -0.00 -0.00 -0.00 -0.00 -0.00 0.00]Max SF 50.000 at 0.000 Max BM -131.417 at 0.000 Max def -5.983 at 2.000
member 3
[174.90 159.60 144.30 129.00 113.70 98.40 83.10 67.80 52.50 37.20 21.90 6.60 -8.70 -24.00 -39.30 -54.60 -69.90 -85.20 -100.50 -115.80 -131.10][-131.42 -81.24 -35.65 5.34 41.75 73.56 100.79 123.42 141.47 154.93 163.79 168.07 167.75 162.85 153.36 139.27 120.60 97.33 69.48 37.03 -0.00][0.00 0.00 0.00 0.01 0.01 0.01 0.01 0.02 0.02 0.02 0.02 0.02 0.02 0.02 0.02 0.01 0.01 0.01 0.01 0.00 0.00]Max SF 174.903 at 0.000 Max BM 168.067 at 3.300 Max def 18.996 at 3.300
`
	fmt.Println(rezstring)
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("Bending Moment and SF analysis test failed")
	}	
}

*/

/*

	cb := &CBm{
		Fck:25.0,
		DL:23.2,
		LL:20.0,
		Nspans:3,
		Lspans:[]float64{10.0,10.0,10.0},
		Sections:[][]float64{{300,600}},
		Clvrs:[][]float64{{0,0,0},{0,0,0}},
		PSFs:[]float64{1.4,1.0,1.6,0.0},
		DM:0.0,
	}
	termstr := "caca"
	CBeamEnvRcc(cb, termstr)
	
	cb = &CBm{
		Fck:25.0,
		DL:5,
		LL:1.5,
		Nspans:4,
		Lspans:[]float64{4.0},
		Sections:[][]float64{{1000,130}},
		Clvrs:[][]float64{{0,0,0},{0,0,0}},
		PSFs:[]float64{1.5,1.0,1.5,0.0},
		DM:0.0,
	}
	termstr = "caca"
	CBeamEnvRcc(cb, termstr)
*/
