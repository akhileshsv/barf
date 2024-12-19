package barf

import (
	"os"
	"fmt"	
	"testing"
	"path/filepath"
	kass"barf/kass"
)

func TestLinkSec(t *testing.T){
	c := &RccCol{
		Dims:[]float64{1000,1000,250,250},
		Styp:7,
	}
	c.Init()
	c.Rbrlvl = 0
	c.BarDbars()
	c.Sec.Draw("dumb")
	fmt.Println(c.Sec.Txtplot)
	c.Lsec.Draw("dumb")
	fmt.Println(c.Lsec.Txtplot)
}

func TestColDraw(t *testing.T){
	var examples = []string{"hulse5.6.6"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/col/")
	for _, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		c, err := ReadCol(fname)
		if err != nil{
			t.Errorf("column draw test failed")
		}
		c.Term = "qt"
		err = ColDesign(&c)
		if err != nil{
			fmt.Println(err)
			t.Errorf("column draw test failed")
		}
		pltstr := c.PlotColDet()
		
		fmt.Println(pltstr)
	}
}

func TestColBx(t *testing.T){
	//biaxial bending
	var examples = []string{"hulse5.6.6","sub14.9"}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/col/")
	for i, ex := range examples {
		if i == 0{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		c, err := ReadCol(fname)
		if err != nil{
			t.Errorf("column design (bs) test failed")
		}
		c.Term = "qt"
		err = ColDesign(&c)
		if err != nil{
			fmt.Println(err)
			t.Errorf("column design (bs) test failed")
		}
		//if i == 1{c.Pur = c.Pu}
		//mcomp, err := ColAzBs(&c, c.Pur)
		//if err != nil{
		//	t.Errorf("column design (bs) test failed")
		//}
		rezstring += fmt.Sprintf("%s axial load %.2f kn moment %.2f kn.m %.2f kn.m\n",ex, c.Pu,c.Mux,c.Muy)
		rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\nult.axial load %0.2f ult.moment %0.2f kN.m \ncomputed moment %0.2f kN.m\n", c.Ast, c.Asc, c.Ast+c.Asc, c.Pur, c.Murx, 0.0)
	}
}

func TestColWeirdBs(t *testing.T){
	//name says it all
	var rezstring string
	//subramanian circular column
	var examples = []string{"sub14.7","mosley4.11"}
	//mosley 4.11 triangular unsymmetrical steel (3 H25) section
	//pu = 354.0; mux = 68.9; muy = 0.0
	//get ast and asc

	//rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/col/")
	for i, ex := range examples {
		//if i != 0  {continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		c, err := ReadCol(fname)
		if err != nil{
			t.Errorf("non-rect column design test failed (file read error)")
		}
		switch i{
			case 0:
			rezstring += "sub 14.7 circular section pu 1400 kn mur 90 knm\n"
			case 1:
			rezstring += "mosley 4.11 triangular section pu 354 kn mur 68.9 knm\n"
		}
		rezstring += fmt.Sprintf("%s axial load %.2f kn moment %.2f kn.m %.2f kn.m\n",ex, c.Pu,c.Mux,c.Muy)
		c.Term = "dumb"
		err = ColDesign(&c)
		if err != nil{
			t.Errorf("non-rect column design test failed")
		}
		c.Table(false)
		t.Log(c.Report)
		t.Log(c.Txtplot[0])
		rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\ntie dia %.f pitch %.f\n", c.Ast, c.Asc, c.Asteel,c.Dtie,c.Ptie)
	}
	wantstring := ``
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("non-rect column design test failed")
	}
}

func TestColAzGen(t *testing.T){
	//hulse ex. 5.5 nm curve diamond section
	var rezstring string
	var pu, pur, mur float64
	var err error
	pu = 409.55
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
	c := &RccCol{
		Fck:25.0,
		Fy:460.0,
		Nbars:4.0,
		Dias:[]float64{25,25,25,25},
		Dbars:[]float64{80,212,212,344},
		Styp:-1,
		Sec:sec,
		Code:2,
		Subck:true,
	}
	pur, mur, err = ColAzGen(c, pu)
	if err != nil{
		t.Errorf("general ('-')7 column analysis test failed")
	}
	rezstring += fmt.Sprintf("pur %.2f kn mur %.2f kn.m",pur,mur)
	wantstring := `hulse 5.5.6
axial load (kN)	 moment (kN.m)
pur 409.85 kn mur 86.02 kn.m`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("general ('-')7 column analysis test failed")
	}
}

func TestColStl(t *testing.T){
	c := &RccCol{
		Fck:40.0,
		Fy:460.0,
		Cvrt:37.5,
		Cvrc:37.5,
		B:400.0,H:500.0,
		Styp:1,
		Dtyp:1,
		Rtyp:0,
		Code:2,
		Subck:true,
	}
	var pu, mux, muy float64
	var rezstring string
	pu = 3000.0; mux = 450.0
	rezstring += fmt.Sprintf("hulse 5.2.6 axial load %.2f kN moment %.2f kN.m\n",pu,mux)
	pur, mur, err := ColStl(c, pu, mux, muy)
	if err != nil{fmt.Println(err)}
	rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\n ult. axial load %0.2f ult.moment %0.2f kN.m\n", c.Ast, c.Asc, c.Ast+c.Asc, pur, mur)


	
	c = &RccCol{
		Fck:25.0,
		Fy:460.0,
		Cvrt:60,
		Cvrc:80,
		B:300.0,H:400.0,
		Styp:1,
		Dtyp:1,
		Rtyp:1,
		Code:2,
		Subck:true,
	}
	pu = 1100.0; mux = 230.0
	rezstring += fmt.Sprintf("hulse 5.3.6 axial load %.2f kN moment %.2f kN.m\n",pu,mux)
	pur, mur, err = ColStl(c, pu, mux, muy)
	if err != nil{fmt.Println(err)}
	rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\n ult. axial load %0.2f ult.moment %0.2f kN.m\n", c.Ast, c.Asc, c.Ast+c.Asc, pur, mur)



	
	c = &RccCol{
		Fck:30.0,
		Fy:460.0,
		Cvrt:50,
		Cvrc:100,
		B:400.0,
		Styp:2,
		Dtyp:1,
		Rtyp:1,
		Code:2,
		Subck:true,
	}
	pu = 354; mux = 68.9
	rezstring += fmt.Sprintf("mosley 4.11 unsym steel axial load %.2f kN moment %.2f kN.m\n",pu,mux)
	pur, mur, err = ColStl(c, pu, mux, muy)
	if err != nil{fmt.Println(err)}
	rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\n ult. axial load %0.2f ult.moment %0.2f kN.m\n", c.Ast, c.Asc, c.Ast+c.Asc, pur, mur)
	fmt.Println(rezstring)
	//fmt.Println(c.Xu)
	
	c = &RccCol{
		Fck:30.0,
		Fy:460.0,
		Cvrt:50,
		Cvrc:100,
		B:400.0,
		Styp:2,
		Dtyp:1,
		Rtyp:2,
		Code:2,
		Subck:true,
	}
	pu = 354; mux = 68.9
	rezstring += fmt.Sprintf("mosley 4.11 L0 steel axial load %.2f kN moment %.2f kN.m\n",pu,mux)
	pur, mur, err = ColStl(c, pu, mux, muy)
	if err != nil{fmt.Println(err)}
	rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\n ult. axial load %0.2f ult.moment %0.2f kN.m\n", c.Ast, c.Asc, c.Asteel, pur, mur)
}

/*
func TestColRect(t *testing.T){
	//test circular column dz via eq sqr for bs and is codes
	//entry func ColDesign
	examples := []string{"shah7.4.1"}
	for i, ex := range examples{
		t.Logf("example %v ")
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		c, err := ReadCol(fname)
		if err != nil{
			t.Errorf("column design (bs) test failed")
		}
		rezstring += fmt.Sprintf("%s axial load %.2f kn moment %.2f kn.m %.2f kn.m\n",ex, c.Pu,c.Mux,c.Muy)
		err = ColDzIs(&c)
		if err != nil{
			t.Errorf("column design (is) test failed")
		}
		c.Printz()
		fmt.Println("dia",c.D1,c.D2)
	}
}


func TestColFlip(t *testing.T){
	//test basic flipping
	c := &RccCol{
		Fck:40.0,
		Fy:460.0,
		Cvrt:37.5,
		Cvrc:37.5,
		B:400.0,H:600.0,
		Styp:2,
		Dtyp:1,
		Rtyp:0,
		Code:2,
		Subck:true,
	}
	c.Init()
	c.GetLinkSec()
	//fmt.Println("1->",c.Sec.Styp)
	cy := ColFlip(c)
	
	cy.GetLinkSec()
	//_ = PlotColGeom(c, "qt")	
	//fmt.Println(cy.Sec.Styp, cy.Styp)
	cy.Sec.Draw("dumb")
	fmt.Println(cy.Sec.Txtplot)
	//_ = PlotColGeom(cy, "qt")
	//fmt.Println(len(cy.Lsec.Coords))
	
	//test asteel gen

}

   
*/
