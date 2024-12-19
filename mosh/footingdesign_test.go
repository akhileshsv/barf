package barf

import(
	"os"
	"path/filepath"
	"fmt"
	"testing"
)

func TestFtngPadAz(t *testing.T) {
	//now all inputs are in METERS flip again
	//this feels like the mars orbiter
	//square footings for square columns alone? 
	f := &RccFtng{
		Colx:0.400,
		Coly:0.400,
		Based:0.600,
		Effd:0.520,
		Sbc:200.0,
		Pus:[]float64{1000,600,0},
		Mxs:[]float64{100,60,0},
		Lx:0.0,
		Typ:0,
		Shape:"square",
		Code:2,
	}
	var rezstring string
	rezstring += "hulse 6.1.6 square pad footing design\n"
	rezstring += fmt.Sprintf("column size %0.2f m\nbase depth %0.2f m\neffective depth %0.2f m\nsbc %0.2f kN/m2\n", f.Coly, f.Based, f.Effd, f.Sbc)
	rezstring += fmt.Sprintf("axial dl %0.2f kN moment dl %.2f kN-m\naxial ll %0.2f kN moment ll %.2f kN-m\naxial wl %0.2f kN moment wl %.2f kN-m\n",f.Pus[0],f.Mxs[0], f.Pus[1], f.Mxs[1], f.Pus[2], f.Mxs[2])
	err := FtngPadAz(f)
	if err != nil {
		rezstring += fmt.Sprintf("%s", err)
	}
	
	rezstring += fmt.Sprintf("base dimensions lx %.2f m ly %.2f m\nmax stress %.2f kN/m2 min stress %.2f kN/m2 basecy %.2f m\n", f.Lx, f.Ly, f.Qmax, f.Qmin, f.Basecy)
	rezstring += fmt.Sprintf("design moment %.2f kN-m design shear %.2f kN\npunching shear %.2f kN punching shear perimeter %.2f m\n", f.Mur, f.Vur, f.Vp, f.Vpp)
	wantstring := `hulse 6.1.6 square pad footing design
column size 0.40 m
base depth 0.60 m
effective depth 0.52 m
sbc 200.00 kN/m2
axial dl 1000.00 kN moment dl 100.00 kN-m
axial ll 600.00 kN moment ll 60.00 kN-m
axial wl 0.00 kN moment wl 0.00 kN-m
base dimensions lx 3.20 m ly 3.20 m
max stress 199.95 kN/m2 min stress 141.35 kN/m2 basecy 3.20 m
design moment 818.74 kN-m design shear 737.22 kN
punching shear 1474.63 kN punching shear perimeter 7.84 m
`
	t.Logf("%s",rezstring)
	if rezstring != wantstring {
		t.Errorf("pad footing (BS) analysis test failed")
	} 
}

func TestFtngBxOz(t *testing.T){
	var bx, by, pu, mx, my float64
	bx = 2.5; by = 1.50; pu = 400.0; mx = 120.0; my = 150.0
	FtngBxOz(bx, by, pu, mx, my)
}

func TestFtngDet(t *testing.T){
	var examples = []string{"shah8.1","sub15.3","rojas1","sub15.5"}
	var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/ftng")
	t.Log(ColorPurple,"testing footing design (rojas)\n",ColorReset)
	for i, ex := range examples {
		if i > 0{continue}
		fname := filepath.Join(datadir,ex+".json")
		rezstring += fmt.Sprintf("example-%v %v",i+1,fname)
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		f, err := ReadFtng(fname)
		if err != nil{
			fmt.Println(err)
			t.Errorf("footing design (rojas) test failed")
			t.Fail()
			os.Exit(1)
		}
		//f.Sloped = false
		f.Term = "svgmono"; f.Verbose = true
		//if i == 3{f.Term = "qt"}
		err = FtngDzRojas(&f)
		if err != nil{
			t.Errorf("footing design (rojas) test failed")
		}
		pltstr := PlotFtngDet(&f)
		fmt.Println(pltstr)
	}
	
}

func TestFtngDzRojas(t *testing.T){
	var examples = []string{"shah8.1","sub15.3","rojas1","sub15.5"}
	var rezstring string
	//rezstring += "\n"
	//okay which of this works? seriously or is there no difference in which case
	//just use this?
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/mosh/ftng")
	t.Log(ColorPurple,"testing footing design (rojas)\n",ColorReset)
	for i, ex := range examples {
		
		fname := filepath.Join(datadir,ex+".json")
		rezstring += fmt.Sprintf("example-%v %v\n",i+1,fname)
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		f, err := ReadFtng(fname)
		if err != nil{
			fmt.Println(err)
			t.Errorf("footing design (rojas) test failed")
			t.Fail()
			os.Exit(1)
		}
		f.Term = ""; f.Verbose = false
		//if i == 3{f.Term = "qt"}
		err = FtngDzRojas(&f)
		if err != nil{
			t.Errorf("footing design (rojas) test failed")
		}
		rezstring += f.Printz()
		//t.Log(f.Report)
	}
	wantstring := `example-1 c:\Users\Admin\junk\barf\data\examples\mosh\ftng\shah8.1.json
8.1Shah
sloped footing
dims- lx 1180.0, ly 1480.0 mm
grade of concrete M 15.0
steel - Fe 415
col dims - x 230.0 mm y 530.0 mm offset 25.0 mm
cover - nominal 60.0 mm effective 80.0 mm
footing depth 430 mm
min.edge depth 150 mm
mux, vux, vp - 72.38 knm, 77.16 knm, 672.00 kn
rbr x- dia 8 mm nx 14 nos spcx 100 mm asx 704 mm2 asx- astx 11 mm2
muy, vuy, vp - 90.78, 81.89, 672.00
rbr y- dia 8 mm ny 17 nos spcy 60 mm asy 855 mm2 asy- asty 47 mm2
example-2 c:\Users\Admin\junk\barf\data\examples\mosh\ftng\sub15.3.json
15.3Sub
sloped footing
dims- lx 2570.0, ly 2570.0 mm
grade of concrete M 20.0
steel - Fe 415
col dims - x 354.0 mm y 354.0 mm offset 346.0 mm
cover - nominal 60.0 mm effective 80.0 mm
footing depth 690 mm
min.edge depth 250 mm
mux, vux, vp - 429.92 knm, 277.76 knm, 1546.74 kn
rbr x- dia 8 mm nx 42 nos spcx 50 mm asx 2111 mm2 asx- astx 47 mm2
muy, vuy, vp - 429.92, 281.21, 1546.74
rbr y- dia 8 mm ny 42 nos spcy 50 mm asy 2111 mm2 asy- asty 15 mm2
example-3 c:\Users\Admin\junk\barf\data\examples\mosh\ftng\rojas1.json
1Rojas
pad footing
dims- lx 3100.0, ly 3100.0 mm
grade of concrete M 20.0
steel - Fe 415
col dims - x 400.0 mm y 400.0 mm offset 0.0 mm
cover - nominal 60.0 mm effective 80.0 mm
footing depth 690 mm
mux, vux, vp - 640.36 knm, 482.21 knm, 1465.91 kn
rbr x- dia 8 mm nx 60 nos spcx 50 mm asx 3016 mm2 asx- astx 49 mm2
muy, vuy, vp - 591.90, 529.36, 1465.91
rbr y- dia 8 mm ny 56 nos spcy 50 mm asy 2815 mm2 asy- asty 42 mm2
example-4 c:\Users\Admin\junk\barf\data\examples\mosh\ftng\sub15.5.json
15.5Sub
pad footing
dims- lx 2500.0, ly 2400.0 mm
grade of concrete M 20.0
steel - Fe 415
col dims - x 400.0 mm y 300.0 mm offset 0.0 mm
cover - nominal 60.0 mm effective 80.0 mm
footing depth 610 mm
mux, vux, vp - 358.05 knm, 286.62 knm, 1219.45 kn
rbr x- dia 8 mm nx 38 nos spcx 60 mm asx 1910 mm2 asx- astx 10 mm2
muy, vuy, vp - 308.59, 341.41, 1219.45
rbr y- dia 8 mm ny 33 nos spcy 70 mm asy 1659 mm2 asy- asty 0 mm2
`
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("footing design (rojas) test failed")
	}
	
}


/*
	   ...beware the beast man for he is the devil's pawn?
	   YE OLDE
	colx := 0.4; coly := 0.4; based := 0.6; effd := 0.52; sbc := 200.0
	pud := 1000.0; mud := 100.0; pul := 600.0; mul := 60.0
	puw := 0.0; muw := 0.0; lx := 0.0
	shape := "square"
*/


/*

func TestFtngPso(t *testing.T){
	//footing = lx, ly, d, px, py
	var dims, nps, ngens int
	var pmin, pmax, pus, mxs, mys, psfs []float64
	var colx, coly, fck, fy, df, sbc, pgck, pgsoil, nomcvr float64
	var vecs [][]float64
	dims = 4
	pmin = []float64{0.15,0.12,0.12,0.0}
	pmax = []float64{4.0,2.0,2.0,10.0}
	nps = 30
	ngens = 2
	colx, coly, fck, fy, df, sbc, pgck, pgsoil, nomcvr = 0.4, 0.4, 21.0, 420.0, 1.5, 220.0, 24.0, 15.0, 80.0
	pus = []float64{700,500}
	mxs = []float64{140,100}
	mys = []float64{120, 80}
	hmin := FtngHmin(pus, mxs)
	pmin[3] = hmin
	psfs = []float64{1.2, 1.6}
	vecs = [][]float64{
		{colx, coly, fck, fy, df, sbc, pgck, pgsoil, nomcvr},
		pus,
		mxs,
		mys,
		psfs,
	}
	randswarm(dims, nps, ngens, pmin, pmax, vecs, "ftng")
}

*/

	/*
	   	var colx, coly, df, eo, dmin, fck, fy, sbc, pgck, pgsoil, nomcvr float64
	var pus, mxs, mys, psfs []float64
	var shape, plot string
	var sloped, dlfac bool
	//shah ex. 1
	fmt.Println("shah ex 8.1")
	colx = 0.23; coly = 0.53; df = 0.0; eo = 0.025
	fck = 15.0; fy = 415.0
	sbc = 400.0; pgck = 24.0; pgsoil = 15.0; nomcvr = 0.06; dmin = 0.15
	pus = []float64{633}
	mxs = []float64{0}
	mys = []float64{0}
	psfs = []float64{1.5}
	shape = "rect"
	sloped = true
	dlfac = true
	plot = "caca"
	FtngDzRojas(colx, coly, fck, fy, df, dmin, eo, sbc, pgck, pgsoil, nomcvr, pus, mxs, mys, psfs, shape, sloped, dlfac, plot)
	fmt.Println("sub ex 15.3")
	colx = 0.354; coly = 0.354; df = 0.0; eo = 0.173
	fck = 15.0; fy = 415.0
	sbc = 200.0; pgck = 24.0; pgsoil = 15.0; nomcvr = 0.06; dmin = 0.25
	pus = []float64{1200}
	mxs = []float64{0}
	mys = []float64{0}
	psfs = []float64{1.5}
	shape = "square"
	sloped = true
	dlfac = true
	FtngDzRojas(colx, coly, fck, fy, df, dmin, eo, sbc, pgck, pgsoil, nomcvr, pus, mxs, mys, psfs, shape, sloped, dlfac, plot)
	fmt.Println("rojas ex. 1")
	//rojas ex.
	colx = 0.4; coly = 0.4; df = 1.5; fck = 20.0; fy = 415.0
	sbc = 220.0; pgck = 24.0; pgsoil = 15.0
	pus = []float64{700,500}
	mxs = []float64{0,0}
	mys = []float64{0,0}
	psfs = []float64{1.2,1.6}
	FtngDzRojas(colx, coly, fck, fy, df, sbc, pgck, pgsoil, pus, mxs, mys, psfs, shape)
	
func TestFtngDzBxRojas(t *testing.T){
	var colx, coly, df, eo, dmin, fck, fy, sbc, pgck, pgsoil, nomcvr float64
	var pus, mxs, mys, psfs []float64
	var shape, plot string
	var sloped, dlfac bool
	fmt.Println("rojas ex. 1")
	//rojas ex.1
	colx = 0.4; coly = 0.4; df = 1.5; fck = 20.0; fy = 415.0
	sbc = 220.0; pgck = 24.0; pgsoil = 15.0; nomcvr = 0.06
	pus = []float64{700,500}
	mxs = []float64{140,100}
	mys = []float64{120,80}
	psfs = []float64{1.2,1.6}
	shape = "rect"
	plot = "caca"
	FtngDzRojas(colx, coly, fck, fy, df, dmin, eo, sbc, pgck, pgsoil, nomcvr, pus, mxs, mys, psfs, shape, sloped, dlfac, plot)
}

func TestFtngRskot(t *testing.T) {
	var colx, coly, df, eo, dmin, fck, fy, sbc, pgck, pgsoil, nomcvr float64
	var pus, mxs, mys, psfs []float64
	var shape, plot string
	var sloped, dlfac bool
	
	fmt.Println("RSKOTA FOOTING C1")
	colx = 0.46; coly = 0.23; df = 1.5; eo = 0.025
	fck = 15.0; fy = 550.0
	sbc = 150.0; pgck = 24.0; pgsoil = 15.0; nomcvr = 0.06; dmin = 0.15
	pus = []float64{382}
	mxs = []float64{0}
	mys = []float64{0}
	psfs = []float64{1.5}
	shape = "rect"
	sloped = true
	dlfac = true
	plot = ""
	FtngDzRojas(colx, coly, fck, fy, df, dmin, eo, sbc, pgck, pgsoil, nomcvr, pus, mxs, mys, psfs, shape, sloped, dlfac, plot)	

	fmt.Println("RSKOTA FOOTING C2")
	colx = 0.46; coly = 0.23; df = 1.5; eo = 0.025
	fck = 15.0; fy = 550.0
	sbc = 150.0; pgck = 24.0; pgsoil = 15.0; nomcvr = 0.06; dmin = 0.15
	pus = []float64{361}
	mxs = []float64{0}
	mys = []float64{0}
	psfs = []float64{1.5}
	shape = "rect"
	sloped = true
	dlfac = true
	plot = ""
	FtngDzRojas(colx, coly, fck, fy, df, dmin, eo, sbc, pgck, pgsoil, nomcvr, pus, mxs, mys, psfs, shape, sloped, dlfac, plot)	

	fmt.Println("RSKOTA FOOTING C3")
	colx = 0.46; coly = 0.23; df = 1.5; eo = 0.025
	fck = 15.0; fy = 550.0
	sbc = 150.0; pgck = 24.0; pgsoil = 15.0; nomcvr = 0.06; dmin = 0.15
	pus = []float64{264}
	mxs = []float64{0}
	mys = []float64{0}
	psfs = []float64{1.5}
	shape = "rect"
	sloped = true
	dlfac = true
	plot = ""
	FtngDzRojas(colx, coly, fck, fy, df, dmin, eo, sbc, pgck, pgsoil, nomcvr, pus, mxs, mys, psfs, shape, sloped, dlfac, plot)	

}
	   
	*/
