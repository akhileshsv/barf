package barf

import (
	"os"
	"fmt"
	"path/filepath"
	"testing"
)


func TestSubFrmOpt(t *testing.T){
	var examples = []string{"mosley7.1","allen3.1","shah8.7"}
	for i, ex := range examples{
		if i != 2{continue}
		dirname,_ := os.Getwd()
		datadir := filepath.Join(dirname,"../data/examples/mosh/subfrm/")
		fname := filepath.Join(datadir,ex+".json")
		fmt.Println(ColorCyan,"example->","no.",i+1,"file->",fname,ColorReset)
		sf, err := ReadSubFrm(fname)
		if err != nil{
			fmt.Println(err)
			t.Errorf("subframe design test failed")
		}
		sf.Term = "svg"
		sf.Web = true
		sf.Opt = 21
		sf.Width = 225.0
		var sfrez SubFrm
		sfrez, err = OptSubFrm(sf)
		t.Log(err)
		if err == nil{
			_ = PlotSfDet(&sfrez)	
			for i, txtplot := range sfrez.Txtplots{
				fmt.Println("i,txtplot->",i, txtplot)
			}
		} else {
			fmt.Println("ERRORE, errore->",err)
		}
		// _ = PlotSfDet(&sfrez)
		// fmt.Println("textplots")
		// for i, txtplot := range sfrez.Txtplots{
		// 	fmt.Println("i,txtplot->",i, txtplot)
		// }
		//t.Errorf("NOTHING HAS BEEN DESIGNED what impossible laziness never have i seen")
	}
}


func TestFltSlb(t *testing.T){
	var examples = []string{"allen14.1","sub11.1"}
	//var examples = []string{"allen3.1"}
	//var rezstring string
	for i, ex := range examples{
		if i == 1{continue}
		dirname,_ := os.Getwd()
		datadir := filepath.Join(dirname,"../data/examples/mosh/subfrm/")
		fname := filepath.Join(datadir,ex+".json")
		fmt.Println(ColorCyan,"example->","no.",i+1,"file->",fname,ColorReset)
		sf, err := ReadSubFrm(fname)
		if err != nil{
			fmt.Println(err)
			t.Errorf("subframe design test failed")
		}
		//sf.Term = "qt"
		err = FltSlbDz(&sf)
		if err != nil{
			t.Errorf("flat slab design failed")
		}
	}
}

func TestSubFrmMosley(t *testing.T){
	var rezstring string
	sf := &SubFrm{
		Nspans:3,
		Lspans:[]float64{6.0,4.0,6.0},
		Hs:[]float64{4.0,3.5},
		DL:25.0,
		LL:10.0,
		Sections:[][]float64{{300,350},{300,600}},
		
		Csec:[]int{1},
		Bsec:[]int{2},
		Fcks:[]float64{25.0},
		Fys:[]float64{415.0},
		Code:2,
		PSFs:[]float64{1.4,1.0,1.6,0.0},
		DM:0.0,
		Selfwt:0,
		Term:"dumb",
	}
	err := CalcSubFrm(sf)
	if err != nil {
		fmt.Println(err)
		t.Errorf("subframe envelope test failed")
	}
	rezstring += "mosley7.1.5"
	rezstring += "\nbeam envelopes\n"
	for _, id := range sf.Beams {
		//fmt.Println("ID",id)
		rezstring += fmt.Sprintf("member id %v\n",id)
		bm := sf.Bmenv[id]
		//fmt.Println("BM",bm)
		for i := 0; i < 21; i++ {
			rezstring += fmt.Sprintf("sec %v\tshear %.3f\tmoment-hog %.3f\tmoment-sag %.3f\n",i+1,bm.Venv[i],bm.Mnenv[i],bm.Mpenv[i])	
		}
		rezstring += fmt.Sprintf("max hogging moment %.3f max sagging moment %.3f\n",bm.Mnmax,bm.Mpmax)
	}
	rezstring += "\ncolumn moment envelopes\n"
	for _, id := range sf.Cols {
		rezstring += fmt.Sprintf("member id %v\n",id)
		c := sf.Colenv[id]
		rezstring += fmt.Sprintf("max top moment %.3f max bottom moment %.3f\n",c.Mtmax,c.Mbmax)
	}
	//rezstring += sf.Txtplots[0]
	//rezstring += sf.Txtplots[1]
	t.Log(rezstring)
	//t.Log(sf.Txtplots[0])
	//t.Log(sf.Txtplots[1])
	wantstring := `mosley7.1.5
beam envelopes
member id 9
sec 1	shear 142.716	moment-hog -74.450	moment-sag 0.000
sec 2	shear 127.416	moment-hog -33.930	moment-sag 0.000
sec 3	shear 112.116	moment-hog 0.000	moment-sag 3.858
sec 4	shear 96.816	moment-hog 0.000	moment-sag 34.014
sec 5	shear 81.516	moment-hog 0.000	moment-sag 60.089
sec 6	shear 66.216	moment-hog 0.000	moment-sag 82.249
sec 7	shear 50.916	moment-hog 0.000	moment-sag 99.819
sec 8	shear 35.616	moment-hog 0.000	moment-sag 112.799
sec 9	shear 20.316	moment-hog 0.000	moment-sag 121.189
sec 10	shear 5.016	moment-hog 0.000	moment-sag 124.988
sec 11	shear -14.228	moment-hog 0.000	moment-sag 124.198
sec 12	shear -29.528	moment-hog 0.000	moment-sag 118.818
sec 13	shear -44.828	moment-hog 0.000	moment-sag 108.848
sec 14	shear -60.128	moment-hog 0.000	moment-sag 94.288
sec 15	shear -75.428	moment-hog 0.000	moment-sag 75.137
sec 16	shear -90.728	moment-hog 0.000	moment-sag 51.397
sec 17	shear -106.028	moment-hog -3.078	moment-sag 23.067
sec 18	shear -121.328	moment-hog -25.743	moment-sag 0.000
sec 19	shear -136.628	moment-hog -64.437	moment-sag 0.000
sec 20	shear -151.928	moment-hog -107.720	moment-sag 0.000
sec 21	shear -167.228	moment-hog -155.593	moment-sag 0.000
max hogging moment -155.593 max sagging moment 124.988
member id 10
sec 1	shear 116.436	moment-hog -123.519	moment-sag 0.000
sec 2	shear 106.236	moment-hog -101.252	moment-sag 0.000
sec 3	shear 96.036	moment-hog -81.025	moment-sag 0.000
sec 4	shear 85.836	moment-hog -63.390	moment-sag 0.000
sec 5	shear 75.636	moment-hog -56.890	moment-sag 0.000
sec 6	shear 65.436	moment-hog -51.390	moment-sag 0.000
sec 7	shear 55.236	moment-hog -46.890	moment-sag 8.097
sec 8	shear 45.036	moment-hog -43.390	moment-sag 15.237
sec 9	shear 34.836	moment-hog -40.890	moment-sag 20.337
sec 10	shear 24.636	moment-hog -39.390	moment-sag 23.397
sec 11	shear 14.436	moment-hog -38.890	moment-sag 24.417
sec 12	shear -24.636	moment-hog -39.390	moment-sag 23.397
sec 13	shear -34.836	moment-hog -40.890	moment-sag 20.337
sec 14	shear -45.036	moment-hog -43.390	moment-sag 15.237
sec 15	shear -55.236	moment-hog -46.890	moment-sag 8.097
sec 16	shear -65.436	moment-hog -51.390	moment-sag 0.000
sec 17	shear -75.636	moment-hog -56.890	moment-sag 0.000
sec 18	shear -85.836	moment-hog -63.390	moment-sag 0.000
sec 19	shear -96.036	moment-hog -81.025	moment-sag 0.000
sec 20	shear -106.236	moment-hog -101.252	moment-sag 0.000
sec 21	shear -116.436	moment-hog -123.519	moment-sag 0.000
max hogging moment -123.519 max sagging moment 24.417
member id 11
sec 1	shear 167.228	moment-hog -155.593	moment-sag 0.000
sec 2	shear 151.928	moment-hog -107.720	moment-sag 0.000
sec 3	shear 136.628	moment-hog -64.437	moment-sag 0.000
sec 4	shear 121.328	moment-hog -25.743	moment-sag 0.000
sec 5	shear 106.028	moment-hog -3.078	moment-sag 23.067
sec 6	shear 90.728	moment-hog 0.000	moment-sag 51.397
sec 7	shear 75.428	moment-hog 0.000	moment-sag 75.137
sec 8	shear 60.128	moment-hog 0.000	moment-sag 94.288
sec 9	shear 44.828	moment-hog 0.000	moment-sag 108.848
sec 10	shear 29.528	moment-hog 0.000	moment-sag 118.818
sec 11	shear 14.228	moment-hog 0.000	moment-sag 124.198
sec 12	shear -5.016	moment-hog 0.000	moment-sag 124.988
sec 13	shear -20.316	moment-hog 0.000	moment-sag 121.189
sec 14	shear -35.616	moment-hog 0.000	moment-sag 112.799
sec 15	shear -50.916	moment-hog 0.000	moment-sag 99.819
sec 16	shear -66.216	moment-hog 0.000	moment-sag 82.249
sec 17	shear -81.516	moment-hog 0.000	moment-sag 60.089
sec 18	shear -96.816	moment-hog 0.000	moment-sag 34.014
sec 19	shear -112.116	moment-hog 0.000	moment-sag 3.858
sec 20	shear -127.416	moment-hog -33.930	moment-sag 0.000
sec 21	shear -142.716	moment-hog -74.450	moment-sag 0.000
max hogging moment -155.593 max sagging moment 124.988

column moment envelopes
member id 1
max top moment -17.434 max bottom moment -34.810
member id 2
max top moment 11.022 max bottom moment 22.049
member id 3
max top moment -11.022 max bottom moment -22.049
member id 4
max top moment 17.434 max bottom moment 34.810
member id 5
max top moment -39.640 max bottom moment -19.782
member id 6
max top moment 25.214 max bottom moment 12.611
member id 7
max top moment -25.214 max bottom moment -12.611
member id 8
max top moment 39.640 max bottom moment 19.782
`
	t.Log(rezstring)
	//t.Logf(sf.Txtplots[0])	
	//t.Logf(sf.Txtplots[1])
	if rezstring != wantstring{
		//fmt.Println(rezstring)
		t.Errorf("subframe envelope test failed")
	}
}

func TestSubFrmDz(t *testing.T){
	var examples = []string{"mosley7.1","allen3.1","shah10.3","shah8.7"}
	//var examples = []string{"allen3.1"}
	//var rezstring string
	for i, ex := range examples{
		if i != 2{continue}
		dirname,_ := os.Getwd()
		datadir := filepath.Join(dirname,"../data/examples/mosh/subfrm/")
		fname := filepath.Join(datadir,ex+".json")
		fmt.Println(ColorCyan,"example->","no.",i+1,"file->",fname,ColorReset)
		sf, err := ReadSubFrm(fname)
		if err != nil{
			fmt.Println(err)
			t.Errorf("subframe design test failed")
		}
		sf.Term = "svg"
		_, _ = sf.ChainSlab()
		_ = PlotSfDet(&sf)
		fmt.Println("textplots")
		for i, txtplot := range sf.Txtplots{
			fmt.Println("i,txtplot->",i, txtplot)
		}
		//t.Errorf("NOTHING HAS BEEN DESIGNED what impossible laziness never have i seen")
	}
}

func TestSubFrmChain(t *testing.T){
	var examples = []string{"shah10.3"}
	//var examples = []string{"allen3.1"}
	//var rezstring string
	for i, ex := range examples{
		dirname,_ := os.Getwd()
		datadir := filepath.Join(dirname,"../data/examples/mosh/subfrm/")
		fname := filepath.Join(datadir,ex+".json")
		fmt.Println(ColorCyan,"example->","no.",i+1,"file->",fname,ColorReset)
		sf, err := ReadSubFrm(fname)
		if err != nil{
			fmt.Println(err)
			t.Errorf("subframe chain test failed")
		}
		_, err = sf.ChainSlab()
		fmt.Println(err)
		//fmt.Println(sf.Report)
		//sf.Term = "dxf"
		_ = PlotSfDet(&sf)
		fmt.Println("textplots")
		for i, txtplot := range sf.Txtplots{
			fmt.Println("i,txtplot->",i, txtplot)
		}
		//t.Errorf("NOTHING HAS BEEN DESIGNED what impossible laziness never have i seen")
	}
}
/*
func TestSubFrmBmSizes(t *testing.T) {
	//var rezstring string
	sf := &SubFrm{
		Nspans:3,
		Lspans:[]float64{6.0,4.0,6.0},
		H0:4.0,
		H1:3.5,
		DL:25.0,
		LL:10.0,
		Sections:[][]float64{{300,350},{300,600}},
		Colidxs:[]int{1},
		Bmidxs:[]int{2},
		Fcks:[]float64{25.0},
		Fys:[]float64{415.0},
		Clvrs:[][]float64{{0},{0}},
		PSFs:[]float64{1.4,1.0,1.6,0.0},
		DM:0.0,
	}
	_ = SubFrmBmSizes(sf, 300.0, 0, 40.0)
	//t.Logf(bmsecs)
}

func TestCalcSubFrm(t *testing.T) {
	var rezstring string
	sf := &SubFrm{
		Nspans:3,
		Lspans:[]float64{6.0,4.0,6.0},
		H0:4.0,
		H1:3.5,
		DL:25.0,
		LL:10.0,
		Sections:[][]float64{{300,350},{300,600}},
		Colidxs:[]int{1},
		Bmidxs:[]int{2},
		Fcks:[]float64{25.0},
		Clvrs:[][]float64{{0,0,0},{0,0,0}},
		PSFs:[]float64{1.4,1.0,1.6,0.0},
		DM:0.0,
	}
	
	bmenv, colenv, bmvec, colvec, err := CalcSubFrm(sf,"dumb")
	if err != nil {
		fmt.Println(err)
		t.Errorf("subframe envelope test failed")
	}
	rezstring += "mosley7.1.5"
	rezstring += "\nbeam envelopes\n"
	for _, id := range bmvec {
		//fmt.Println("ID",id)
		rezstring += fmt.Sprintf("member id %v\n",id)
		bm := bmenv[id]
		//fmt.Println("BM",bm)
		for i := 0; i < 21; i++ {
			rezstring += fmt.Sprintf("sec %v\tshear %.3f\tmoment-hog %.3f\tmoment-sag %.3f\n",i+1,bm.Venv[i],bm.Mnenv[i],bm.Mpenv[i])	
		}
		rezstring += fmt.Sprintf("max hogging moment %.3f max sagging moment %.3f\n",bm.Mnmax,bm.Mpmax)
	}
	rezstring += "\ncolumn moment envelopes\n"
	for _, id := range colvec {
		rezstring += fmt.Sprintf("member id %v\n",id)
		c := colenv[id]
		rezstring += fmt.Sprintf("max top moment %.3f max bottom moment %.3f\n",c.Mtmax,c.Mbmax)
	}
	rezstring += sf.Txtplots[0]
	rezstring += sf.Txtplots[1]
	t.Logf(rezstring)
 	sf = &SubFrm{
		Nspans:3,
		Lspans:[]float64{10.0,10.0,10.0},
		H0:4.6,
		H1:3.6,
		DL:23.2,
		LL:20.0,
		Sections:[][]float64{{300,300},{400,400},{300,600}},
		Colidxs:[]int{1,2,2,1},
		Bmidxs:[]int{3},
		Fcks:[]float64{25.0},
		Clvrs:[][]float64{{0,0,0},{0,0,0}},
		PSFs:[]float64{1.4,1.0,1.6,0.0},
		DM:0.0,
	}
	
	bmenv, colenv, bmvec, colvec, err = CalcSubFrm(sf,"dumb")
	if err != nil {
		fmt.Println(err)
		t.Errorf("subframe envelope test failed")
	}
	rezstring += "allen3.1"
	rezstring += "\nbeam envelopes\n"
	for _, id := range bmvec {
		//fmt.Println("ID",id)
		rezstring += fmt.Sprintf("member id %v\n",id)
		bm := bmenv[id]
		//fmt.Println("BM",bm)
		for i := 0; i < 21; i++ {
			rezstring += fmt.Sprintf("sec %v\tshear %.3f\tmoment-hog %.3f\tmoment-sag %.3f\n",i+1,bm.Venv[i],bm.Mnenv[i],bm.Mpenv[i])	
		}
		rezstring += fmt.Sprintf("max hogging moment %.3f max sagging moment %.3f\n",bm.Mnmax,bm.Mpmax)
	}
	rezstring += "\ncolumn moment envelopes\n"
	for _, id := range colvec {
		rezstring += fmt.Sprintf("member id %v\n",id)
		c := colenv[id]
		rezstring += fmt.Sprintf("max top moment %.3f max bottom moment %.3f\n",c.Mtmax,c.Mbmax)
	}
	rezstring += sf.Txtplots[0]
	rezstring += sf.Txtplots[1]
	t.Logf(rezstring)	
}

	sf = &SubFrm{
		Nspans:3,
		Lspans:[]float64{10.0,10.0,10.0},
		H0:4.6,
		H1:3.6,
		DL:23.2,
		LL:20.0,
		Sections:[][]float64{{300,300},{400,400},{300,600}},
		Colidxs:[]int{1,2,2,1},
		Bmidxs:[]int{3},
		Fcks:[]float64{25.0},
		Clvrs:[][]float64{{0,0,0},{0,0,0}},
		PSFs:[]float64{1.4,1.0,1.6,0.0},
		DM:0.3,
	}
	
	bmenv, colenv, bmvec, colvec, err = CalcSubFrm(sf,"dumb")
	if err != nil {
		fmt.Println(err)
		t.Errorf("subframe envelope test failed")
	}
	rezstring += "allen3.1"
	rezstring += "\nbeam envelopes\n"
	for _, id := range bmvec {
		//fmt.Println("ID",id)
		rezstring += fmt.Sprintf("member id %v\n",id)
		bm := bmenv[id]
		//fmt.Println("BM",bm)
		for i := 0; i < 21; i++ {
			rezstring += fmt.Sprintf("sec %v\tshear %.3f\tmoment-hog %.3f\tmoment-sag %.3f\n",i+1,bm.Venv[i],bm.Mnenv[i],bm.Mpenv[i])	
		}
		rezstring += fmt.Sprintf("max hogging moment %.3f max sagging moment %.3f\n",bm.Mnmax,bm.Mpmax)
	}
	rezstring += "\ncolumn moment envelopes\n"
	for _, id := range colvec {
		rezstring += fmt.Sprintf("member id %v\n",id)
		c := colenv[id]
		rezstring += fmt.Sprintf("max top moment %.3f max bottom moment %.3f\n",c.Mtmax,c.Mbmax)
	}
	rezstring += sf.Txtplots[0]
	rezstring += sf.Txtplots[1]
	t.Logf(rezstring)

*/
