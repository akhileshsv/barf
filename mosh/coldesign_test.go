package barf

import(
	"os"
	"fmt"
	"math"
	"path/filepath"
	kass"barf/kass"
	"testing"
)

func TestColSecArXu(t *testing.T){
	//get area of section at y = dck from top face
	var rezstring string
	c := &RccCol{
		Fck:40.0,
		Fy:460.0,
		Asc:1610.0,
		Ast:982.0,
		Cvrt:60.0,
		Cvrc:60.0,
		B:200.0,H:400.0,
		Styp:1,
	}
	c.SecInit()
	dck := 200.0
	ack, xcy, ycy := ColSecArXu(c.Sec, dck)
	rezstring += fmt.Sprintf("ack %f xcy %f ycy %f\n",ack, xcy, ycy)

	c = &RccCol{
		Fck:40.0,
		Fy:460.0,
		Asc:1610.0,
		Ast:982.0,
		Cvrt:60.0,
		Cvrc:60.0,
		B:300.0,
		Styp:2,
	}
	c.SecInit()
	dck = 200.0
	ack, xcy, ycy = ColSecArXu(c.Sec, dck)
	rezstring += fmt.Sprintf("ack %f xcy %f ycy %f\n",ack, xcy, ycy)
	wantstring := `ack 40000.000000 xcy 100.000000 ycy 300.000000
ack 23094.010768 xcy 150.000000 ycy 126.474288
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("column section area function failed")
	}

}

func TestColAzBs(t *testing.T) {
	var pu float64
	var rezstring string
	c := &RccCol{
		Fck:40.0,
		Fy:460.0,
		Asc:1610.0,
		Ast:982.0,
		Cvrt:60.0,
		Cvrc:60.0,
		B:350.0,H:450.0,
		Pu:1760.0,
	}
	pu = 1760.0
	mur, err := ColAzBs(c, pu)
	rezstring += fmt.Sprintf("hulse 5.1.6 pu %.f kn mur %0.2f kn-m\n", pu, mur)
	if err != nil {
		t.Log(err)
		t.Errorf("column(BS) analysis test failed")
	}
	c = &RccCol{
		Fck:30.0,
		Fy:460.0,
		Asc:1383.3,
		Ast:1383.3,
		Cvrt:60.0,
		Cvrc:60.0,
		B:300.0,H:400.0,
		Pu:2084.0,
	}
	pu = 2084.0
	mur, err = ColAzBs(c, pu)
	if err != nil {
		t.Log(err)
		t.Errorf("column(BS) analysis test failed")
	}
	rezstring += fmt.Sprintf("mosley/bungey 9.2 pu %.f kn mur %0.2f kn-m\n", pu, mur)
	wantstring := `hulse 5.1.6 pu 1760 kn mur 322.89 kn-m
mosley/bungey 9.2 pu 2084 kn mur 102.53 kn-m
`
	t.Log(rezstring)
	if rezstring != wantstring {
		t.Errorf("column(BS) analysis test failed")
	}
}

func TestColDzBs(t *testing.T){	
	var examples = []string{"hulse5.2.6","hulse5.3.6","hulse5.6.6","hulse5.8.6"}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/col/")
	for i, ex := range examples {
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		c, err := ReadCol(fname)
		if err != nil{
			t.Errorf("column design (bs) test failed")
		}
		err = ColDzBs(&c)
		if err != nil{
			fmt.Println(err)
			t.Errorf("column design (bs) test failed")
		}
		if i == 1{c.Pur = c.Pu}
		mcomp, err := ColAzBs(&c, c.Pur)
		if err != nil{
			t.Errorf("column design (bs) test failed")
		}
		rezstring += fmt.Sprintf("%s axial load %.2f kn moment %.2f kn.m %.2f kn.m\n",ex, c.Pu,c.Mux,c.Muy)
		rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\nult.axial load %0.2f ult.moment %0.2f kN.m \ncomputed moment %0.2f kN.m\n", c.Ast, c.Asc, c.Ast+c.Asc, c.Pur, c.Murx, mcomp)
	}
	wantstring := `hulse5.2.6 axial load 3000.00 kn moment 450.00 kn.m 0.00 kn.m
ast 2227.50 asc 2227.50 astot 4455.00
ult.axial load 3054.88 ult.moment 450.00 kN.m 
computed moment 441.70 kN.m
hulse5.3.6 axial load 1100.00 kn moment 230.00 kn.m 0.00 kn.m
ast 1051.69 asc 2221.95 astot 3273.64
ult.axial load 1100.00 ult.moment 222.09 kN.m 
computed moment 224.02 kN.m
hulse5.6.6 axial load 1200.00 kn moment 80.00 kn.m 75.00 kn.m
ast 1357.26 asc 1357.26 astot 2714.53
ult.axial load 1223.18 ult.moment 115.71 kN.m 
computed moment 115.71 kN.m
hulse5.8.6 axial load 3000.00 kn moment 0.00 kn.m 0.00 kn.m
ast 1517.30 asc 1517.30 astot 3034.61
ult.axial load 2970.56 ult.moment 354.50 kN.m 
computed moment 361.21 kN.m
`
	t.Log(rezstring)
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("column design (bs) test failed")
	}
}

func TestColNMBs(t *testing.T){
	var rezstring string	
	rezstring += "hulse 5.4.6\naxial load (kN)\t moment (kN.m)\n"
	c := &RccCol{
		Fck:40.0,
		Fy:460.0,
		B:350.0,H:450.0,
		Nbars:4.0,
		Dias:[]float64{32,32,25,25},
		Dbars:[]float64{60,60,390,390},
		Styp:1,
	}
	pus, mus := ColNMBs(c)
	for idx, pu := range pus{	
		pu = math.Ceil(pu/10.0)/100.0; mu := math.Ceil(mus[idx]/1e4)/100.0
		rezstring += fmt.Sprintf("%0.5f kN %0.5f kN.m\n", pu, mu)
	}
	t.Log("ex 1")
	t.Log(c.Nmplot)
	rezstring += "hulse 5.5.6\naxial load (kN)\t moment (kN.m)\n"
	sec := &kass.SectIn{
		Ncs:[]int{5},
		Wts:[]float64{1.0},
		Coords:[][]float64{
			{-212,0},
			{0,-212},
			{212,0},
			{0,212},
			{-212,0},
		},
	}
	sec.SecInit()
	c = &RccCol{
		Fck:25.0,
		Fy:460.0,
		Nbars:4.0,
		Dias:[]float64{25,25,25,25},
		Dbars:[]float64{80,212,212,344},
		Styp:-1,
		Sec:sec,
	}
	pus, mus = ColNMBs(c)
	for idx, pu := range pus {	
		pu = math.Ceil(pu/10.0)/100.0; mu := math.Ceil(mus[idx]/1e4)/100.0
		rezstring += fmt.Sprintf("%0.5f kN %0.5f kN.m\n", pu, mu)
	}
	t.Log("ex 2")
	t.Log(c.Nmplot)
	rezstring += "mosley/bungey 4.11\n"
	c = &RccCol{
		Fck:30.0,
		Fy:460.0,
		Nbars:3.0,
		Dias:[]float64{25,25,25},
		Dbars:[]float64{100,296,296},
		Styp:2,
		B:400.0,
	}
	c.SecInit()
	pus, mus = ColNMBs(c)
	for idx, pu := range pus {	
		pu = math.Ceil(pu/10.0)/100.0; mu := math.Ceil(mus[idx]/1e4)/100.0
		rezstring += fmt.Sprintf("%0.5f kN %0.5f kN.m\n", pu, mu)		
	}
	t.Log("ex 3")
	t.Log(c.Nmplot)
	wantstring := `hulse 5.4.6
axial load (kN)	 moment (kN.m)
115.14000 kN 159.93000 kN.m
492.92000 kN 220.88000 kN.m
770.62000 kN 262.73000 kN.m
998.28000 kN 293.74000 kN.m
1143.73000 kN 308.60000 kN.m
1271.30000 kN 317.93000 kN.m
1398.88000 kN 324.67000 kN.m
1526.45000 kN 328.84000 kN.m
1654.03000 kN 330.41000 kN.m
1868.87000 kN 315.01000 kN.m
2072.80000 kN 298.82000 kN.m
2265.83000 kN 281.85000 kN.m
2450.13000 kN 263.74000 kN.m
2627.33000 kN 244.21000 kN.m
2798.70000 kN 223.06000 kN.m
2965.21000 kN 200.13000 kN.m
3127.61000 kN 175.30000 kN.m
3286.53000 kN 148.46000 kN.m
3442.47000 kN 119.52000 kN.m
3595.83000 kN 88.43000 kN.m
3647.72000 kN 78.23000 kN.m
3669.30000 kN 74.67000 kN.m
3689.15000 kN 71.39000 kN.m
3707.48000 kN 68.37000 kN.m
3724.45000 kN 65.57000 kN.m
3740.20000 kN 62.97000 kN.m
3754.87000 kN 60.55000 kN.m
3768.56000 kN 58.29000 kN.m
3781.37000 kN 56.18000 kN.m
3793.38000 kN 54.19000 kN.m
3804.66000 kN 52.33000 kN.m
3815.28000 kN 50.58000 kN.m
3825.29000 kN 48.93000 kN.m
3834.74000 kN 47.37000 kN.m
3843.68000 kN 45.89000 kN.m
3852.16000 kN 44.50000 kN.m
3860.19000 kN 43.17000 kN.m
3867.83000 kN 41.91000 kN.m
3871.10000 kN 41.37000 kN.m
hulse 5.5.6
axial load (kN)	 moment (kN.m)
75.49000 kN 78.78000 kN.m
255.38000 kN 84.19000 kN.m
409.56000 kN 86.57000 kN.m
591.13000 kN 83.19000 kN.m
773.15000 kN 76.98000 kN.m
932.58000 kN 70.16000 kN.m
1072.42000 kN 62.88000 kN.m
1194.95000 kN 55.29000 kN.m
1301.88000 kN 47.63000 kN.m
1394.51000 kN 40.12000 kN.m
1473.87000 kN 33.01000 kN.m
1540.76000 kN 26.57000 kN.m
1595.81000 kN 21.08000 kN.m
1639.56000 kN 16.81000 kN.m
1672.42000 kN 14.06000 kN.m
1697.22000 kN 12.57000 kN.m
1711.59000 kN 11.23000 kN.m
1720.88000 kN 10.01000 kN.m
1729.46000 kN 8.87000 kN.m
1737.40000 kN 7.82000 kN.m
1744.78000 kN 6.85000 kN.m
1751.64000 kN 5.94000 kN.m
1758.05000 kN 5.10000 kN.m
1764.05000 kN 4.31000 kN.m
1769.67000 kN 3.57000 kN.m
1774.95000 kN 2.87000 kN.m
1779.92000 kN 2.21000 kN.m
1784.60000 kN 1.59000 kN.m
1789.03000 kN 1.01000 kN.m
1793.21000 kN 0.46000 kN.m
mosley/bungey 4.11
12.00000 kN 72.86000 kN.m
159.56000 kN 71.29000 kN.m
294.91000 kN 69.76000 kN.m
415.92000 kN 67.30000 kN.m
526.77000 kN 63.93000 kN.m
634.42000 kN 60.11000 kN.m
740.10000 kN 55.66000 kN.m
844.77000 kN 50.38000 kN.m
949.19000 kN 44.12000 kN.m
1053.96000 kN 36.70000 kN.m
1159.58000 kN 28.00000 kN.m
1266.44000 kN 17.85000 kN.m
1308.26000 kN 14.23000 kN.m
1329.54000 kN 12.84000 kN.m
1349.11000 kN 11.57000 kN.m
1367.18000 kN 10.39000 kN.m
1383.91000 kN 9.30000 kN.m
1399.44000 kN 8.29000 kN.m
1413.91000 kN 7.35000 kN.m
1427.41000 kN 6.47000 kN.m
1440.04000 kN 5.65000 kN.m
1451.87000 kN 4.88000 kN.m
1463.00000 kN 4.16000 kN.m
1473.46000 kN 3.48000 kN.m
1483.33000 kN 2.84000 kN.m
1492.65000 kN 2.23000 kN.m
1501.47000 kN 1.66000 kN.m
1509.82000 kN 1.11000 kN.m
1517.75000 kN 0.60000 kN.m
1524.36000 kN 0.17000 kN.m
`
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("column(BS) NM curve test failed")
	}
}

func TestColEffHt(t *testing.T) {
	//now all these values are in METERS flip
	//3 columns (current, lower, upper)
	//4 beams (lower left, lower right, upper left, upper right) 
	braced := true
	ljbase := false
	basemr := false
	ljbeamss := false
	ujbeamss := false
	cbs := []float64{0.35,0.35,0.35}
	cds := []float64{0.35,0.35,0.35}
	cls := []float64{4.0,4.0,4.0}
	bbs := []float64{0.30,0.30,0.30,0.30}
	bds := []float64{0.60,0.60,0.60,0.60}
	bls := []float64{8.0,7.0,8.0,7.0}
	l0 := 3.45
	lmin, slender, err := ColEffHt(braced, ljbase, basemr, ljbeamss, ujbeamss, cbs, cds, cls, bbs, bds, bls, l0)
	if err != nil {
		fmt.Println(err)
	}
	var rezstring string
	rezstring += fmt.Sprintf("hulse 5.7.6 effective height %0.3f meters \n slender -> %v ",lmin, slender)
	wantstring := `hulse 5.7.6 effective height 2.564 meters 
 slender -> false `
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("column effective height (BS) test failed")
	}
}
//func ColSlmDBs(c *RccCol, pu, m1, m2 float64, braced bool) (float64, float64, error)
func TestColSlmDBs(t *testing.T) {
	c := &RccCol{
		Fck:40.0,
		Fy:460.0,
		Cvrt:37.5,
		Cvrc:37.5,
		B:400.0,H:500.0,
		Styp:1,
		Rtyp:0,
		Leffx:6.75,
		Leffy:6.75,
		D1:37.5,
		Braced:true,
		Pu:3000.0,
		Mt:225.0,
		Mb:-225.0,
		Slender:true,
	}
	var rezstring string
	pur, mur, err := ColSlmDBs(c)
	if err != nil {
		fmt.Println(err)
	}
	rezstring += fmt.Sprintf("hulse 5.8.6 axial load %.2f kN moment %.2f kN.m\n",pur,mur)
	mur, err = ColAzBs(c, pur)
	if err != nil {
		fmt.Println(err)
	}
	rezstring += fmt.Sprintf("computed moment %.2f kN.m\n",mur)
	rezstring += fmt.Sprintf("total area of steel %.2f sq.mm\n",c.Asc + c.Ast)
	wantstring := `hulse 5.8.6 axial load 3041.47 kN moment 356.37 kN.m
computed moment 363.32 kN.m
total area of steel 3193.53 sq.mm
`
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("symmetrical rebar slender column design (BS) test failed")
	}
}

/*
   YE OLDE
func TestColIs(t *testing.T) {
	c := &RccCol{
		Fck:20.0,
		Fy:250.0,
		Efcvr:51.0,
		Dimxcl:32.0,
		Dimincl:12.0,
	}
	var rezstring string
	var mur float64
	ColDIs(c)
	rezstring += fmt.Sprintf("shah 7.3-1 Mur %0.5f kN.m\n")
	wantstring := ``
	if rezstring != wantstring {
		t.Errorf("column design test failed")
	}
        }

   	
	var pu, mux, muy float64
	c := &RccCol{
		Fck:40.0,
		Fy:460.0,
		Cvrt:37.5,
		Cvrc:37.5,
		B:400.0,H:500.0,
		Styp:1,
		Rtyp:0,
		Pu:3000.0,
		Mux:450.0,
		Muy:0.0,
	}

	
	pu = 3000.0; mux = 450.0; muy = 0.0

	
	c = &RccCol{
		Fck:25.0,
		Fy:460.0,
		Cvrt:60.0,
		Cvrc:80.0,
		B:300.0,H:400.0,
		Styp:1,
		Rtyp:1,
	}
	pu = 1100.0; mux = 230.0; muy = 0.0
	rezstring += fmt.Sprintf("hulse 5.3.6 axial load %.2f kN moment %.2f kN.m\n",pu,mux)
	pur, mur, _ = ColDBs(c, pu, mux, muy)
	mcomp, err = ColAzBs(c, pu)
	if err != nil {
		rezstring += fmt.Sprintf("%s", err)
	}
	rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\n ult.axial load %0.2f ult.moment %0.2f kN.m \n computed moment %0.2f kN.m\n", c.Ast, c.Asc, c.Asc+c.Ast, pur, mur, mcomp)

	c = &RccCol{
		Fck:30.0,
		Fy:460.0,
		D1:60.0,
		D2:70.0,
		B:350.0,H:300.0,
		Styp:1,
		Rtyp:0,
	}
	pu = 1200.0; mux = 80.0; muy = 75.0
	rezstring += fmt.Sprintf("hulse 5.6.6 axial load %.2f kN moment x %.2f kN.m moment y %.2f kN.m\n",pu,mux,muy)
	pur, mur, _ = ColDBs(c, pu, mux, muy)
	mcomp, err = ColAzBs(c, pu)
	if err != nil {
		rezstring += fmt.Sprintf("%s", err)
	}
	rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\n ult.axial load %0.2f ult.moment %0.2f kN.m \n computed moment %0.2f kN.m\n", c.Ast, c.Asc, c.Asc+c.Ast, pur, mur, mcomp)

	c = &RccCol{
		Fck:15.0,
		Fy:415.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:700.0,
		Styp:1,
		Rtyp:0,
	}
	pu = 622.0; mux = 184.0; muy = 0.0
	rezstring += fmt.Sprintf("shah 7.4.1 axial load %.2f kN moment x %.2f kN.m moment y %.2f kN.m\n",pu,mux,muy)
	pur, mur, _ = ColDBs(c, pu, mux, muy)
	
	mcomp, err = ColAzBs(c, pu)
	if err != nil {
		rezstring += fmt.Sprintf("%s", err)
	}
	
	rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\n ult.axial load %0.2f ult.moment %0.2f kN.m \n computed moment %0.2f kN.m\n", c.Ast, c.Asc, c.Asc+c.Ast, pur, mur, mcomp)

	
	c = &RccCol{
		Fck:20.0,
		Fy:415.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:400.0,
		Styp:1,
		Rtyp:0,
	}
	pu = 1090.0; mux = 19.0; muy = 0.0
	rezstring += fmt.Sprintf("shah 7.4.2 axial load %.2f kN moment x %.2f kN.m moment y %.2f kN.m\n",pu,mux,muy)
	pur, mur, _ = ColDBs(c, pu, mux, muy)
	
	mcomp, err = ColAzBs(c, pu)
	if err != nil {
		rezstring += fmt.Sprintf("%s", err)
	}
	
	rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\n ult.axial load %0.2f ult.moment %0.2f kN.m \n computed moment %0.2f kN.m\n", c.Ast, c.Asc, c.Asc+c.Ast, pur, mur, mcomp)

	
	wantstring := ``
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("column(BS) design test failed")
	}

*/
