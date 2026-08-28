package barf

import (
	"os"
	"testing"
	"path/filepath"
	flay"barf/flay"
)

func TestFcflr(t *testing.T){
	var examples = []string{"fpln",}
	var rezstring string
	rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/flay")
	t.Log(ColorPurple,"testing fc floor calc examples\n",ColorReset)
	for i, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		f, err := flay.ReadFcflr(fname)
		if err != nil{
			t.Fatal(err)
		}
		f.Term = "qt"
		FlrDz(f)
	}

}

func TestStlFlrRdark0(t *testing.T){
	f := StlFlr{
		Title:"rdark",
		Xs:[]float64{0,5510,11020,16530},
		Ys:[]float64{0,7100,14200,21300,28400,33600,38800},
		Zs:[]float64{0,3600,7200},
		Jspc:1800,
		Ohng:0.0,
		DL:6.0,
		LL:6.0,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Ht:3600.0,
		Dzcbs:true,
		Dzbs:true,
		Term:"dxf",
		Brc:[]int{1,3},
		Vz:52.5,
	}
	err := f.Calc()
	t.Log(err)
	t.Log(f.Report)
}

func TestStlFlrRdark1(t *testing.T){
	f := StlFlr{
		Title:"rdarkva",
		Xs:[]float64{0,5510,11020,16530,22040,27550},
		Ys:[]float64{0,7100,14200,21300,28400,35500},
		Zs:[]float64{0,3600,7200},
		Jspc:1800,
		Ohng:0.0,
		DL:6.0,
		LL:6.0,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Ht:3600.0,
		Dzcbs:true,
		Dzbs:true,
		Term:"svg",
		Vz:52.5,
		Matread:true,
		Mat:[][]int{
			{0,1,1,1,0},
			{0,1,1,1,0},
			{0,1,1,1,0},
			{0,1,1,1,0},
			{1,1,1,1,1},
		},
	}
	err := f.Calc()
	t.Log(err)
	t.Log(f.Report)
}

func TestStlFlrRdark2(t *testing.T){
	f := StlFlr{
		Title:"rdarkvb",
		Xs:[]float64{0,5510,11020,16530,22040,27550},
		Ys:[]float64{0,7100,14200,21300,28400,35500},
		Zs:[]float64{0,3600,7200},
		Jspc:1800,
		Ohng:0.0,
		DL:6.0,
		LL:6.0,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Ht:3600.0,
		Dzcbs:true,
		Dzbs:true,
		Term:"svg",
		Vz:52.5,
		Matread:true,
		Mat:[][]int{
			{0,1,1,1,0},
			{0,1,1,1,0},
			{0,1,1,1,0},
			{1,1,0,1,1},
			{1,1,0,1,1},
		},
	}
	err := f.Calc()
	t.Log(err)
	t.Log(f.Report)
}

func TestStlFlrRdark4(t *testing.T){
	f := StlFlr{
		Title:"rdarkvc",
		Xs:[]float64{0,7100,14200,21300,28400,33600,38800},
		Ys:[]float64{0,5510,11020,16530},
		Zs:[]float64{0,3600,7200},
		Jspc:1800,
		Ohng:0.0,
		DL:6.0,
		LL:6.0,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Ht:3600.0,
		Dzcbs:true,
		Dzbs:true,
		Term:"svg",
		Brc:[]int{1,1},
		Vz:52.5,
	}
	err := f.Calc()
	t.Log(err)
	t.Log(f.Report)
}

func TestStlFlrRdark5(t *testing.T){
	f := StlFlr{
		Title:"rdarkvd",
		Xs:[]float64{0,7500,16500},
		Ys:[]float64{0,3550,7100,10650,14200,17750,21300,24850,28400,31950,35500,39050},
		Zs:[]float64{0,3600,7200},
		Jspc:1800,
		Ohng:0.0,
		DL:6.0,
		LL:6.0,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Ht:3600.0,
		Dzcbs:true,
		Dzbs:false,
		Term:"svg",
		Brc:[]int{1,1},
		Vz:52.5,
	}
	err := f.Calc()
	t.Log(err)
	t.Log(f.Report)
}


func TestStlFlrRdark6(t *testing.T){
	f := StlFlr{
		Title:"rdarkvel",
		Xs:[]float64{0,8000,16000},
		Ys:[]float64{0,4025,8050,12075,16100,20125,24150,28175,32200},
		Zs:[]float64{0,3600,7200},
		Jspc:1800,
		Ohng:0.0,
		DL:6.0,
		LL:6.0,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Ht:3600.0,
		Dzcbs:true,
		Dzbs:true,
		Term:"svg",
		Brc:[]int{1,1},
		Matread:true,
		Mat:[][]int{
			{1,0},
			{1,0},
			{1,0},
			{1,0},
			{1,0},
			{1,0},
			{1,1},
			{1,1},
		},
		Vz:52.5,
	}
	err := f.Calc()
	t.Log(err)
	t.Log(f.Report)
}


func TestStlFlrRdark7(t *testing.T){
	f := StlFlr{
		Title:"rdarkvue",
		Xs:[]float64{0,8000,16000},
		Ys:[]float64{0,4025,8050,12075,16100,20125,24150,28175,32200},
		Zs:[]float64{0,3600,7200},
		Jspc:1800,
		Ohng:0.0,
		DL:6.0,
		LL:6.0,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Ht:3600.0,
		Dzcbs:true,
		Dzbs:true,
		Term:"svg",
		Brc:[]int{1,1},
		Matread:true,
		Mat:[][]int{
			{1,1},
			{1,1},
			{1,0},
			{1,0},
			{1,0},
			{1,0},
			{1,1},
			{1,1},
		},
		Vz:52.5,
	}
	err := f.Calc()
	t.Log(err)
	t.Log(f.Report)
}

func TestStlFlrRdark8(t *testing.T){
	f := StlFlr{
		Title:"rdarkvuf",
		Xs:[]float64{0,4025,8050,12075,16100,20125,24150,28175,32200},
		Ys:[]float64{0,8000,16000},
		Zs:[]float64{0,3600,7200},
		Jspc:1800,
		Ohng:0.0,
		DL:6.0,
		LL:6.0,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Ht:3600.0,
		Dzcbs:true,
		Dzbs:true,
		Term:"svg",
		Brc:[]int{1,1},
		Matread:true,
		Mat:[][]int{
			{1,1,0,0,0,0,1,1},
			{1,1,1,1,1,1,1,1},
		},
		Vz:52.5,
	}
	err := f.Calc()
	t.Log(err)
	t.Log(f.Report)
}


func TestFlrMat(t *testing.T){
	f := StlFlr{
		Title:"rdbog",
		Xs:[]float64{0,3048,6096,9144,12192,15240},
		Ys:[]float64{0,3048,6096,9144,12192,15240},
		Zs:[]float64{0,3600},
		Jspc:1800,
		Ohng:0.0,
		DL:1.5,
		LL:1.5,
		TDL:5.82,
		TLL:1.5,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Matread:true,
		Mat:[][]int{
			{0,1,1,1,0},
			{0,1,1,1,0},
			{0,1,1,1,0},
			{0,1,1,1,0},
			{1,1,1,1,1},
		},
		Ht:3600.0,
		Dzcbs:false,
		Dzbs:false,
		Term:"qt",
		Sname:"ub",
	}
	f.Calc()
}

func TestFlrSar(t *testing.T){
	f := StlFlr{
		Title:"sbsar",
		Xs:[]float64{0,6000,10000,14000,18000,22000,26000,32000},
		Ys:[]float64{0,13700},
		Zs:[]float64{0,5200,10400},
		Jspc:1800,
		Ohng:900.0,
		DL:5.82,
		LL:7.5,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Ht:5200.0,
		Dzcbs:false,
		Dzbs:false,
		Term:"qt",
		Sname:"ub",
	}
	err := f.Calc()
	t.Log(err)
}

func TestFlrBrc(t *testing.T){
	t.Log("testing steel floor bracing design")
	f := StlFlr{
		Title:"rdark",
		Xs:[]float64{0,6800,11500,16200},
		Ys:[]float64{0,4375,8475,15275,20475,24775,31575,38375},
		Zs:[]float64{0,3600,7200},
		Jspc:1800,
		Ohng:0.0,
		DL:5.82,
		LL:4.5,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Ht:3600.0,
		Dzcbs:false,
		Dzbs:false,
		Brc:[]int{1,3},
		Term:"svg",
		Vz:52.5,
	}
	err := f.Calc()
	t.Log(err)
}


func TestStlFlrBasic(t *testing.T){
	f := StlFlr{
		Title:"ark",
		Xs:[]float64{0,5510,11020,16530},
		Ys:[]float64{0,7100,14200},
		Zs:[]float64{0,3600,7200},
		Jspc:1800,
		Ohng:0.0,
		DL:6.0,
		LL:6.0,
		Df:2000.0,
		Sbc:125.0,
		Lsb:true,
		Onism:true,
		Ht:3600.0,
		Dzcbs:false,
		Dzbs:false,
		Term:"svg",
		Brc:[]int{1,1},
		Vz:52.5,
	}
	err := f.Calc()
	t.Log(err)
	t.Log(f.Report)
}

func TestStlFlrOpt(t *testing.T){
	t.Log("checking bay spans")
	var ltot, cw float64
	ltot = 16600
	ncs := []float64{2,3,4,5}
	cw = 300.0
	
	for _, nc := range ncs{
		lbay := (ltot - cw-cw)/(nc-1)
		
		t.Log("ncols",nc,"nbays",nc-1)
		fcbay := lbay
		t.Log("lbay",lbay,"fcbay",fcbay)
	}
	wtot := 32800.0
	
	ncs = []float64{6,7,8,9}
	for _, nc := range ncs{
		lbay := (wtot - cw-cw)/(nc-1)
		t.Log(ColorRed)
		t.Log("ncols",nc,"nbays",nc-1)
		fcbay := lbay
		t.Log("lbay",lbay,"fcbay",fcbay)
		t.Log(ColorReset)
	}
}
