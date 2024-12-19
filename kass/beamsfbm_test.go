package barf

import (
	"os"
	"fmt"
	"testing"
	"path/filepath"
)

func TestBmSsMosley(t *testing.T){
	var examples = []string{"mosley4.1"}
	var rezstring string
	var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for _, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		frmrez, _ = CalcBm1d(mod, 2)
		bmresults := CalcBmSf(mod, frmrez,false)
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


func TestCalcBmEnv(t *testing.T) {
	var examples = []string{"hulse2.3"}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for _, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		loadenvz, _ := CalcBmEnv(mod) 
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
	//fmt.Println(rezstring)
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("Bending Moment and SF envelope test failed")
	}	
}


func TestCalcBmSf(t *testing.T) {
	var examples = []string{"hulse2.1","hulse2.2"}
	var rezstring string
	var frmrez []interface{}
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for _, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		_, mod,_ := JsonInp(fname)
		frmrez, _ = CalcBm1d(mod, 2)
		bmresults := CalcBmSf(mod, frmrez,false)
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
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("Bending Moment and SF analysis test failed")
	}	
}
