package barf

import(
	"fmt"
	"math"
	"testing"
)


func TestBmDetail(t *testing.T){
	var rezstring string
	var err error
	rezstring += "allen 3.1 results detailing\n"
	b := &RccBm{
		Fck:30.0,
		Fy:460.0,
		Tyb:0.0,
		Bw:300.0,
		Dused:775.0,
		Cvrt:50.0,
		Cvrc:50.0,
		Code:1,
		D1:16.0,
		D2:20.0,
		Verbose:true,
	}
	mur := 603.0
	b.Mu = mur
	err = BmDesign(b)
	if err != nil{
		t.Errorf("beam rebar detailing (bs) test failed")
	}
	fmt.Println(err)
	//rezstring += b.BarPrint()
	PlotBmGeom(b, "dxf")
	fmt.Println(rezstring)
	//t.Log(b.Txtplot[0])
}

func TestBmAzGen(t *testing.T) {
	var rezstring string
	rezstring += "example 3.1.6.1 mosley\n"
	b := &RccBm{
		Fck:30.0,
		Fy:460.0,
		Bf:0.0,
		Df:0.0,
		Bw:280.0,
		Dused:560.0,
		Dias:[]float64{20,20,32,32,32},
		Dbars:[]float64{50,50,510,510,510},
		Tyb:0.0,
		Cvrt:50.0,
		Cvrc:50.0,
		Asc:628.0,
		Ast:2410.0,
		Styp:1,
		Code:2,
	}
	bscode := true
	rblk := true
	mr, x, xdrat, err := BmAzGen(b, bscode, rblk)
	if err != nil{
		t.Errorf(err.Error())
	}
	rezstring += "from func GenAz\n"
	rezstring += fmt.Sprintf("moment of resistance %.3f kNm, neutral axis depth %.3f mm x/d ratio %.3f\n", math.Ceil(100*mr)/100,math.Ceil(10*x)/10,math.Ceil(100*xdrat)/100)
	mr, x, xdrat, _ = SecAzBs(b)
	rezstring += "from func SecAz\n"
	rezstring += fmt.Sprintf("moment of resistance %.3f kNm, neutral axis depth %.3f mm x/d ratio %.3f\n", math.Ceil(100*mr)/100,math.Ceil(10*x)/10,math.Ceil(100*xdrat)/100)
	rezstring += "example 3.1.6.2 mosley\n"
	b = &RccBm{
		Fck:25.0,
		Fy:460.0,
		Bf:800.0,
		Df:150.0,
		Bw:300.0,
		Dused:470.0,
		Dias:[]float64{25,25,25},
		Dbars:[]float64{420,420,420},
		Tyb:1.0,
		Cvrt:50.0,
		Cvrc:50.0,
		Asc:0.0,
		Ast:1470.0,
		Code:2,
	}
	mr, x, xdrat, _ = BmAzGen(b, bscode, rblk)
	rezstring += "from func GenAz\n"
	rezstring += fmt.Sprintf("moment of resistance %.3f kNm, neutral axis depth %.3f mm x/d ratio %.3f\n", math.Ceil(100*mr)/100,math.Ceil(10*x)/10,math.Ceil(100*xdrat)/100)
	mr, x, xdrat, _ = SecAzBs(b)
	rezstring += "from func SecAz\n"
	rezstring += fmt.Sprintf("moment of resistance %.3f kNm, neutral axis depth %.3f mm x/d ratio %.3f\n", math.Ceil(100*mr)/100,math.Ceil(10*x)/10,math.Ceil(100*xdrat)/100)
	
	wantstring := `example 3.1.6.1 mosley
from func GenAz
moment of resistance 412.990 kNm, neutral axis depth 210.200 mm x/d ratio 0.380
from func SecAz
moment of resistance 411.990 kNm, neutral axis depth 209.700 mm x/d ratio 0.420
example 3.1.6.2 mosley
from func GenAz
moment of resistance 227.940 kNm, neutral axis depth 73.100 mm x/d ratio 0.160
from func SecAz
moment of resistance 227.950 kNm, neutral axis depth 72.700 mm x/d ratio 0.180
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("rcc beam section analysis test (bs) failed")
	}
}

func TestBmDzBs(t *testing.T){
	var rezstring string
	var err error
	rezstring += "hulse ex. 3.2.6.1 mur 145.0 kn-m\n"
	b := &RccBm{
		Fck:25.0,
		Fy:460,
		Tyb:0.0,
		Bw:230,
		Dused:540,
		Cvrt:50,
		Cvrc:50,
		Code:2,
	}
	mur := 145.0
	b.Mu = mur
	err = BmDesign(b)
	if err != nil{
		t.Errorf("beam design bs test failed")
	}
	//BmStlBs(b, mur)
	rezstring += b.BarPrint()
	PlotBmGeom(b, "dumb")
	t.Log(b.Txtplot[0])
	
	rezstring += "hulse ex. 3.2.6.2 mur 160.0 kn-m\n"
	b = &RccBm{
		Fck:25.0,
		Fy:460.0,
		Tyb:1.0,
		Bf:600.0,
		Df:150.0,
		Bw:250.0,
		Dused:580.0,
		Cvrt:50.0,
		Cvrc:50.0,
		Code:2,
	}
	mur = 160.0
	b.Mu = mur
	err = BmDesign(b)
	if err != nil{
		t.Errorf("beam rebar design (bs) test failed")
	}
	rezstring += b.BarPrint()
	PlotBmGeom(b, "dumb")
	t.Log(b.Txtplot[0])
	
	rezstring += "mosley/bungey ex. 4.1 mur 185.0 kn-m\n"
	b = &RccBm{
		Fck:30.0,
		Fy:460.0,
		Tyb:0.0,
		Bw:260.0,
		Dused:490.0,
		Cvrt:50.0,
		Cvrc:50.0,
		Code:2,
	}
	mur = 185.0
	b.Mu = mur
	err = BmDesign(b)
	if err != nil{
		t.Errorf("beam rebar design (bs) test failed")
	}

	rezstring += b.BarPrint()
	PlotBmGeom(b, "dumb")
	t.Log(b.Txtplot[0])

	rezstring += "mosley/bungey ex. 4.2 mur 263.0 kn-m\n"
	b = &RccBm{
		Fck:30.0,
		Fy:460.0,
		Tyb:0.0,
		Bw:300.0,
		Dused:570.0,
		Cvrt:50.0,
		Cvrc:50.0,
		Code:2,
	}
	mur = 263.0
	b.Mu = mur
	err = BmDesign(b)
	if err != nil{
		t.Errorf("beam rebar design (bs) test failed")
	}

	rezstring += b.BarPrint()
	PlotBmGeom(b, "dumb")
	t.Log(b.Txtplot[0])


	rezstring += "mosley/bungey ex. 4.3 mur 285\n"
	b = &RccBm{
		Fck:30.0,
		Fy:460.0,
		Tyb:0.0,
		Bw:260.0,
		Dused:490.0,
		Cvrt:50.0,
		Cvrc:50.0,
		Code:2,
	}
	mur = 285.0
	b.Mu = mur
	err = BmDesign(b)
	if err != nil{
		t.Errorf("beam rebar design (bs) test failed")
	}

	rezstring += b.BarPrint()
	PlotBmGeom(b, "dumb")
	t.Log(b.Txtplot[0])

	rezstring += "mosley/bungey ex. 4.5 mur 229\n"
	b = &RccBm{
		Fck:30.0,
		Fy:460.0,
		Tyb:1.0,
		Bf:800.0,
		Df:150.0,
		Bw:200.0,
		Dused:460.0,
		Cvrt:40.0,
		Cvrc:40.0,
		Code:2,
	}
	mur = 229.0
	b.Mu = mur
	err = BmDesign(b)
	if err != nil{
		t.Errorf("beam rebar design (bs) test failed")
	}

	rezstring += b.BarPrint()
	PlotBmGeom(b, "dumb")
	t.Log(b.Txtplot[0])

		
	rezstring += "mosley/bungey ex. 4.6 mur 180\n"
	b = &RccBm{
		Fck:30.0,
		Fy:460.0,
		Tyb:1.0,
		Bf:400.0,
		Df:100.0,
		Bw:200.0,
		Dused:390.0,
		Cvrt:40.0,
		Cvrc:40.0,
		Code:2,
	}
	mur = 180.0
	b.Mu = mur
	err = BmDesign(b)
	if err != nil{
		t.Errorf("beam rebar design (bs) test failed")
	}

	rezstring += b.BarPrint()

	PlotBmGeom(b, "dumb")
	t.Log(b.Txtplot[0])

	rezstring += "mosley/bungey ex. 4.7 mur 348\n"
	b = &RccBm{
		Fck:30.0,
		Fy:460.0,
		Tyb:1.0,
		Bw:300.0,
		Bf:450.0,
		Dused:490.0,
		Df:150.0,
		Cvrt:50.0,
		Cvrc:50.0,
		Code:2,
	}
	mur = 348.0
	b.Mu = mur
	err = BmDesign(b)
	if err != nil{
		t.Errorf("beam rebar design (bs) test failed")
	}

	rezstring += b.BarPrint()

	PlotBmGeom(b, "dumb")
	t.Log(b.Txtplot[0])

	rezstring += "pocket section mur "
	b = &RccBm{
		Fck:30.0,
		Fy:460.0,
		Tyb:1.0,
		Bw:300.0,
		Bf:200.0,
		Dused:490.0,
		Df:150.0,
		Cvrt:50.0,
		Cvrc:50.0,
		Styp:14,
		Code:2,
	}
	mur = 148.0
	b.Mu = mur
	err = BmDesign(b)
	if err != nil{
		t.Errorf("beam rebar design (bs) test failed")
	}

	rezstring += b.BarPrint()

	PlotBmGeom(b, "dumb")
	t.Log(b.Txtplot[0])

	wantstring := `hulse ex. 3.2.6.1 mur 145.0 kn-m
tension (bottom) steel
 number of layers 1
combo - 2 nos 16 mm dia 1 nos 25 mm dia  ast prov 893 mm2 ast req 855 mm2 a diff 38 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
hulse ex. 3.2.6.2 mur 160.0 kn-m
tension (bottom) steel
 number of layers 1
combo - 4 nos 16 mm dia 0 nos 0 mm dia  ast prov 804 mm2 ast req 789 mm2 a diff 15 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
mosley/bungey ex. 4.1 mur 185.0 kn-m
tension (bottom) steel
 number of layers 1
combo - 4 nos 20 mm dia 0 nos 0 mm dia  ast prov 1257 mm2 ast req 1255 mm2 a diff 2 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
mosley/bungey ex. 4.2 mur 263.0 kn-m
tension (bottom) steel
 number of layers 1
combo - 3 nos 25 mm dia 0 nos 0 mm dia  ast prov 1473 mm2 ast req 1469 mm2 a diff 4 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
mosley/bungey ex. 4.3 mur 285
tension (bottom) steel
 number of layers 1
combo - 3 nos 32 mm dia 0 nos 0 mm dia  ast prov 2413 mm2 ast req 2044 mm2 a diff 369 mm2
compression (top) steel
 number of layers 1
combo - 3 nos 12 mm dia 0 nos 0 mm dia  ast prov 339 mm2 ast req 308 mm2 a diff 32 mm2
mosley/bungey ex. 4.5 mur 229
tension (bottom) steel
 number of layers 1
combo - 2 nos 32 mm dia 0 nos 0 mm dia  ast prov 1608 mm2 ast req 1456 mm2 a diff 153 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
mosley/bungey ex. 4.6 mur 180
tension (bottom) steel
 number of layers 1
combo - 2 nos 32 mm dia 0 nos 0 mm dia  ast prov 1608 mm2 ast req 1540 mm2 a diff 68 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
mosley/bungey ex. 4.7 mur 348
tension (bottom) steel
 number of layers 1
combo - 3 nos 32 mm dia 0 nos 0 mm dia  ast prov 2413 mm2 ast req 2413 mm2 a diff -0 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
pocket section mur tension (bottom) steel
 number of layers 1
combo - 2 nos 16 mm dia 1 nos 28 mm dia  ast prov 1018 mm2 ast req 1014 mm2 a diff 4 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("beam rebar design (bs) test failed")
	}
}

func TestBmBarGen(t *testing.T){
	var rezstring string
	asts := []float64{452,678,804,904,1206,1256,1482,1658,1884,1964,2190,2366,2590,2964}
	for _, ast := range asts{
		rezstring += fmt.Sprintf("ast req-> %.0f\n",ast)
		b := &RccBm{Bw:230.0,Ast:ast, Cvrt:230.0}
		err, _ := b.BarGen()
		if err != nil{
			t.Log(err)
		} else {
			rezstring += b.BarPrint()
		}
	}
	wantstring := `ast req-> 452
tension (bottom) steel
 number of layers 1
combo - 4 nos 12 mm dia 0 nos 0 mm dia  ast prov 452 mm2 ast req 452 mm2 a diff 0 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 678
tension (bottom) steel
 number of layers 1
combo - 2 nos 16 mm dia 1 nos 20 mm dia  ast prov 717 mm2 ast req 678 mm2 a diff 39 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 804
tension (bottom) steel
 number of layers 1
combo - 4 nos 16 mm dia 0 nos 0 mm dia  ast prov 804 mm2 ast req 804 mm2 a diff 0 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 904
tension (bottom) steel
 number of layers 1
combo - 3 nos 20 mm dia 0 nos 0 mm dia  ast prov 942 mm2 ast req 904 mm2 a diff 38 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 1206
tension (bottom) steel
 number of layers 1
combo - 2 nos 28 mm dia 0 nos 0 mm dia  ast prov 1232 mm2 ast req 1206 mm2 a diff 26 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 1256
tension (bottom) steel
 number of layers 1
combo - 2 nos 25 mm dia 1 nos 20 mm dia  ast prov 1296 mm2 ast req 1256 mm2 a diff 40 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 1482
tension (bottom) steel
 number of layers 1
combo - 2 nos 28 mm dia 1 nos 20 mm dia  ast prov 1546 mm2 ast req 1482 mm2 a diff 64 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 1658
tension (bottom) steel
 number of layers 1
combo - 2 nos 28 mm dia 1 nos 25 mm dia  ast prov 1723 mm2 ast req 1658 mm2 a diff 65 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 1884
tension (bottom) steel
 number of layers 2
combo - 6 nos 20 mm dia 0 nos 0 mm dia  ast prov 1885 mm2 ast req 1884 mm2 a diff 1 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 1964
tension (bottom) steel
 number of layers 2
combo - 4 nos 25 mm dia 0 nos 0 mm dia  ast prov 1963 mm2 ast req 1964 mm2 a diff -1 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 2190
tension (bottom) steel
 number of layers 2
combo - 2 nos 12 mm dia 4 nos 25 mm dia  ast prov 2190 mm2 ast req 2190 mm2 a diff 0 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 2366
tension (bottom) steel
 number of layers 2
combo - 4 nos 25 mm dia 2 nos 16 mm dia  ast prov 2366 mm2 ast req 2366 mm2 a diff 0 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 2590
tension (bottom) steel
 number of layers 2
combo - 2 nos 25 mm dia 2 nos 32 mm dia  ast prov 2591 mm2 ast req 2590 mm2 a diff 1 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
ast req-> 2964
tension (bottom) steel
 number of layers 2
combo - 5 nos 25 mm dia 1 nos 28 mm dia  ast prov 3071 mm2 ast req 2964 mm2 a diff 107 mm2
compression (top) steel
 number of layers 1
combo - 2 nos 12 mm dia 0 nos 0 mm dia  ast prov 226 mm2 ast req 225 mm2 a diff 1 mm2
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("beam rebar generation test failed")
	}
}

func TestBmAzTaper(t *testing.T){
	var rezstring string
	//example 5.2 allen
	b := &RccBm{
		Fck:40.0,
		Fy:460.0,
		Dias:[]float64{25,25,40,40,40,40},
		Dbars:[]float64{50,50,760,760,760,760},		
		Styp:19,
		Dims:[]float64{300,840,200,140,100},
		Cvrc:50.0,
		Cvrt:80.0,
	}
	err := BmAnalyze(b)
	if err != nil{
		fmt.Println(err)
	}
	rezstring += "example 5.2 allen styp 19 (tapered pocket section)\n"
	rezstring += fmt.Sprintf("mur %f neutral axis depth x %f x/d ratio %f\n",b.Mu, b.Xu, b.Xu/b.Dused)

	//example 5.7 subramanian
	b = &RccBm{
		Fck:25.0,
		Fy:415.0,
		Dias:[]float64{28,18,28},
		Dbars:[]float64{500,500,500},		
		Styp:20,
		Dims:[]float64{400,550,235},
		Cvrt:50.0,
	}
	err = BmAnalyze(b)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("coords-\n",b.Sec.Coords)
	rezstring += "example 5.7 sub styp 20 (trapezoidal section)\n"
	rezstring += fmt.Sprintf("mur %f neutral axis depth x %f x/d ratio %f\n",b.Mu, b.Xu, b.Xu/b.Dused)
	

	
	wantstring := ``
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("tapered beam section analysis test failed")
		t.Errorf("CHANGE IS456 BEAM FUNCS GODDAMN")
	}
}
