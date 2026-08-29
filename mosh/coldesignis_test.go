package barf

import (
	"os"
	"fmt"
	"path/filepath"
	"testing"
)

func TestColDzBasic(t *testing.T){
	var rezstring string
	var examples = []string{"shah7.4.1","shah7.4.2","hulse5.2.6","hulse5.3.6","hulse5.6.6","hulse5.8.6","sub14.7","sub14.11"}
	//rezstring += "\n"
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/col/")
	for i, ex := range examples {
		//if i != 7  {continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		c, err := ReadCol(fname)
		if err != nil{
			t.Errorf("column design (is) test failed (file read error)")
		}
		rezstring += fmt.Sprintf("%s axial load %.2f kn moment %.2f kn.m %.2f kn.m\n",ex, c.Pu,c.Mux,c.Muy)
		
		//c.Nlayers = 4
		//err = ColDzIs(&c)
		//_,_, err = ColStl(&c, c.Pu, c.Mux, c.Muy)
		c.Term = "dumb"
		err = ColDesign(&c)
		if err != nil{
			
			t.Errorf("column design (is) test failed")
		}
		t.Log(c.Txtplot[0])
		//c.Table(false)
		t.Log(c.Report)
		//var pur, mur, mcomp float64
		rezstring += fmt.Sprintf("ast %0.2f asc %0.2f astot %0.2f\ntie dia %.f pitch %.f\n", c.Ast, c.Asc, c.Asteel,c.Dtie,c.Ptie)
		pur, mur, err := ColAzGen(&c, c.Pu)
		if err != nil{
			fmt.Println(ColorRed, ex,ColorReset)
			fmt.Println("report",c.Report)
			t.Errorf("column design (is) test failed")
			break
		}
		rezstring += fmt.Sprintf("computed pur %.2f kn mur %.2f kn.m\n", pur, mur)
	}
	wantstring := `shah7.4.1 axial load 622.00 kn moment 184.00 kn.m 0.00 kn.m
ast 666.65 asc 666.65 astot 1608.50
tie dia 8 pitch 230
computed pur 625.68 kn mur 194.80 kn.m
shah7.4.2 axial load 1090.00 kn moment 19.00 kn.m 0.00 kn.m
ast 549.59 asc 549.59 astot 1206.37
tie dia 8 pitch 230
computed pur 1091.23 kn mur 26.97 kn.m
hulse5.2.6 axial load 3000.00 kn moment 450.00 kn.m 0.00 kn.m
ast 2227.50 asc 2227.50 astot 4474.00
tie dia 10 pitch 300
computed pur 3014.31 kn mur 457.31 kn.m
hulse5.3.6 axial load 1100.00 kn moment 230.00 kn.m 0.00 kn.m
ast 1051.69 asc 2159.46 astot 3257.00
tie dia 8 pitch 250
computed pur 1101.29 kn mur 236.47 kn.m
hulse5.6.6 axial load 1200.00 kn moment 80.00 kn.m 75.00 kn.m
ast 1357.26 asc 1357.26 astot 2768.00
tie dia 8 pitch 250
computed pur 1200.21 kn mur 118.94 kn.m
hulse5.8.6 axial load 3000.00 kn moment 0.00 kn.m 0.00 kn.m
ast 1596.76 asc 1596.76 astot 3196.00
tie dia 8 pitch 300
computed pur 3005.06 kn mur 361.67 kn.m
sub14.7 axial load 1400.00 kn moment 90.00 kn.m 0.00 kn.m
ast 0.00 asc 0.00 astot 3850.00
tie dia 8 pitch 300
computed pur 1400.99 kn mur 108.46 kn.m
sub14.11 axial load 4000.00 kn moment 750.00 kn.m 750.00 kn.m
ast 0.00 asc 0.00 astot 15310.55
tie dia 10 pitch 250
computed pur 4037.44 kn mur 1044.77 kn.m
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("column design (is) test failed")
	}

}


func TestColSizeIs(t *testing.T){
	var pu, fck, fy, pg, b float64
	var sectype int
	var rezstring string
	pu = 1.5 * 1650.0
	pg = 0.02
	fck = 25.0
	fy = 415.0
	sectype = 1
	b = 230.0
	dims := ColSizeIs(pu, fck, fy, pg, b, sectype)
	rezstring += fmt.Sprintf("%.f",dims)
	wantstring := `[230 700]`
	if rezstring != wantstring{
		t.Log(rezstring)
		t.Errorf("column init size (is) test failed")
	}
}


/*
   
func TestColDzIsOld(t *testing.T){
	var rezstring string
	rezstring += "shah ex 7.4.1\n"
	var pu, mux, muy float64
	var c *RccCol
	var axy bool
	var opt int
	var plot string
	var err error
	c = &RccCol{
		Fck:15.0,
		Fy:415.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:700.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
		Nlayers:4,
		Type:"rectangle",
	}
	pu = 622.0; mux = 184.0; muy = 0.0
	plot = ""
	err = ColDzIs(c, pu, mux, muy, axy, plot, opt)
	if err != nil{
		log.Println(err)
	}
	c.Printz()
	
	log.Println("shah ex 7.4.2")
	c = &RccCol{
		Fck:20.0,
		Fy:415.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:400.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
		Nlayers:3,
		Type:"rectangle",
	}
	pu = 1090.0; mux = 19.0; muy = 0.0
	err = ColDzIs(c, pu, mux, muy, axy, plot, opt)
	if err != nil{
		log.Println(err)
	}
	c.Printz()
	log.Println("***---***")
	log.Println("sub ex 14.9 - BIAXE")
	pu = 1400.0; mux = 130.0; muy = 60.0
   	c = &RccCol{
		Fck:30.0,
		Fy:500.0,
		Cvrt:65.5,
		Cvrc:65.5,
		B:300.0,
		H:500.0,
		Styp:1,
		Rtyp:0,
		Dtyp:2,
		Nlayers:0,
		Type:"rectangle",
	}
	err = ColDzIs(c, pu, mux, muy, axy, plot, opt)
	if err != nil{
		log.Println(err)
	}
}

func TestColDIs(t *testing.T){
	log.Println("shah ex 3 uniaxe start")
	log.Println("\n****")
	log.Println("dos layers->")
	c := &RccCol{
		Fck:15.0,
		Fy:415.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:700.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
		Nlayers:2,
	}
	pu := 622.0; mux := 184.0
	pur, mur, _ := ColDIs(c, pu, mux)
	log.Println("mur, pur - >",mur, pur)
	c.Printz()

	log.Println("\n****")
	log.Println("quatro layers->")
	c = &RccCol{
		Fck:15.0,
		Fy:415.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:700.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
		Nlayers:4,
	}
	pu = 622.0; mux = 184.0
	pur, mur, _ = ColDIs(c, pu, mux)
	log.Println("mur, pur - >",mur, pur)
	c.Printz()


	log.Println("\n****")
	log.Println("tres layers->")	
	c = &RccCol{
		Fck:15.0,
		Fy:415.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:700.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
		Nlayers:3,
	}
	pu = 622.0; mux = 184.0
	pur, mur, _ = ColDIs(c, pu, mux)
	log.Println("mur, pur - >",mur, pur)
	c.Printz()

	log.Println("\n****")
	log.Println("quatro layers Fe500->")
		
	c = &RccCol{
		Fck:15.0,
		Fy:500.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:700.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
	}
	pu = 622.0; mux = 184.0
	pur, mur, _ = ColDIs(c, pu, mux)
	log.Println("mur, pur - >",mur, pur)
	c.Printz()

	log.Println("\nnuttin")
	c = &RccCol{
		Fck:25.0,
		Fy:415.0,
		Cvrt:60.5,
		Cvrc:60.5,
		B:300.0,
		H:400.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
	}
	pu = 1400.0; mux = 90.0
	pur, mur, _ = ColDIs(c, pu, mux)
	log.Println("mur, pur - >",mur, pur)
	c.Printz()


	c = &RccCol{
		Fck:25.0,
		Fy:415.0,
		Cvrt:60.5,
		Cvrc:60.5,
		B:300.0,
		H:400.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
		Nlayers:4,
	}
	pu = 1400.0; mux = 90.0
	pur, mur, _ = ColDIs(c, pu, mux)
	log.Println("mur, pur - >",mur, pur)
	c.Printz()


	log.Println("\n****")
	log.Println("4 layers->")
	c = &RccCol{
		Fck:25.0,
		Fy:500.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:700.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
		Nlayers:4,
	}
	pu = 622.0; mux = 184.0
	pur, mur, _ = ColDIs(c, pu, mux)
	log.Println("mur, pur - >",mur, pur)
	c.Printz()


	log.Println("\n****")
	log.Println("2 layers->")
	c = &RccCol{
		Fck:35.0,
		Fy:500.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:700.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
		Nlayers:2,
	}
	pu = 622.0; mux = 184.0
	pur, mur, _ = ColDIs(c, pu, mux)
	log.Println("mur, pur - >",mur, pur)
	c.Printz()
}


   //vecin := []int{killr.d1, killr.d2, killr.d3, killr.na, killr.nb}
	//ColOptPsoSimp(c, pu, mux, vecin)
        
	log.Println("shah ex 1 short column")
	pu = 1680.0; mux = 0; muy = 0
	c = &RccCol{
		Fck:20.0,
		Fy:250.0,
		B:400.0,
		H:400.0,
		Dtyp:0,
		Cvrt:51.0,
		Cvrc:51.0,
		Rtyp:0,
		Typ:"rectangle",
	}
	err := ColDzIs(c, pu, mux , muy, axy, true, false)
	if err != nil{
		log.Println(err)
	}
	log.Println("->\nshah ex 2 short->")
	
	c = &RccCol{
		Fck:15.0,
		Fy:250.0,
		B:400.0,
		H:600.0,
		Dtyp:0,
		Cvrt:51.0,
		Rtyp:0,
		Styp:1,
	}
	
	pu = 1800.0
	//log.Println(ColSizeIs(pu, c.Fck, c.Fy, 0.02,c.B, c.Styp))
	pur, mur, _ := ColDIs(c, pu, 0 , 0)
	log.Println("->\nmur KN-m\tpur KN")
	log.Println(mur, pur)
	
	log.Println("->\nsub ex 1 uniaxe")
	pu = 1400 * 1.15
	pg := 0.02
	fck := 25.0
	fy := 415.0
	sectype := 1
	b := 300.0
	dims := ColSizeIs(pu, fck, fy, pg, b, sectype)

	log.Println("->\nel dims->",dims)
	
	c = &RccCol{
		Fck:fck,
		Fy:fy,
		Cvrt:60.5,
		Cvrc:60.5,
		B:300.0,
		H:400.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
	}
	pu = 1400.0; mux := 90.0; muy := 0.0
	pur, mur, _ = ColDIs(c, pu, mux , muy)
	log.Println("mur\tpur\n->",mur, pur)
	pur, mur, _ = ColAzIs(c, pu)
	log.Println("mur\tpur\n->",mur, pur)
	log.Println("***\n***\n***\n***")

	
	log.Println("->\nshah ex 3 uniaxe->")
	
	c = &RccCol{
		Fck:15.0,
		Fy:415.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:700.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
	}
	pu = 622.0; mux = 184.0; muy = 0.0
	pur, mur, _ = ColDIs(c, pu, mux , muy)
	log.Println("mur KN-m\tpur KN")
	log.Println(mur, pur)
	
	pur, mur, err := ColAzIs(c, pu)
	if err != nil{
		log.Println(err)
	}
	log.Println(mur, pur)
	log.Println("tres layers->")
	c = &RccCol{
		Fck:15.0,
		Fy:415.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:700.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
		Nlayers:3,
	}
	pu = 622.0; mux = 184.0; muy = 0.0
	pur, mur, _ = ColDIs(c, pu, mux , muy)
	log.Println("mur, pur - >",mur, pur)
	c.Printz()
        
func TestColEffHtIs(t *testing.T){
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
	lmin, slender, err := ColEffHt(true,ljbase, basemr, ljbeamss, ujbeamss, cbs, cds, cls, bbs, bds, bls, l0)
}

func TestLFlexRatIs(t *testing.T){
	var csecs, bsecs []int
	var cdims, bdims [][]float64
}










	rezstring += "shah ex 7.4.1\n"
	var pu, mux, muy float64
	var c *RccCol
	var axy bool
	var opt int
	var plot string
	var err error
	c = &RccCol{
		Fck:15.0,
		Fy:415.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:700.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
		Nlayers:4,
	}
	pu = 622.0; mux = 184.0; muy = 0.0
	plot = ""
	err = ColDzIs(c, pu, mux, muy, axy, plot, opt)
	if err != nil{
		log.Println(err)
	}
	c.Printz()
	
	log.Println("shah ex 7.4.2")
	c = &RccCol{
		Fck:20.0,
		Fy:415.0,
		Cvrt:48.0,
		Cvrc:48.0,
		B:230.0,
		H:400.0,
		Styp:1,
		Rtyp:0,
		Dtyp:1,
		Nlayers:3,
		Type:"rectangle",
	}
	pu = 1090.0; mux = 19.0; muy = 0.0
	err = ColDzIs(c, pu, mux, muy, axy, plot, opt)
	if err != nil{
		log.Println(err)
	}
	c.Printz()
	log.Println("***---***")
	log.Println("sub ex 14.9 - BIAXE")
	pu = 1400.0; mux = 130.0; muy = 60.0
   	c = &RccCol{
		Fck:30.0,
		Fy:500.0,
		Cvrt:65.5,
		Cvrc:65.5,
		B:300.0,
		H:500.0,
		Styp:1,
		Rtyp:0,
		Dtyp:2,
		Nlayers:0,
		Type:"rectangle",
	}
	err = ColDzIs(c, pu, mux, muy, axy, plot, opt)
	if err != nil{
		log.Println(err)
	}


*/
