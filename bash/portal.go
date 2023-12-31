package barf

//TO DELETE - most of this is now in kass/portal.go
import (
	"fmt"
	"log"
	"math"
	"path/filepath"
	"runtime"
	"os"
	kass"barf/kass"
	//mosh"barf/mosh"
	"github.com/go-gota/gota/dataframe"
	//"github.com/go-gota/gota/series"
)

type Portal struct {
	//Pz - wind pressure 0.6 * vz2
	//pzcs - wind load coeffs
	Nbays, Nframes               int
	Span, Spacing, Slope, Height float64
	Pz, DL, LL                   float64
	Dh, Lh                       float64
	Pzcs                         []float64
	W                            float64
	Fixbs, Haunch, Gable, Mono   bool
	Csec, Bsec                   int
	Sectype                      int
	PSFs                         []float64
	Lx,Ly                        float64
	Vz,Cpi                       float64
	Drawtyp                      string
	Prlnspc, Prlnwt              float64
	LR, Rise                     float64
	Nprlz, Ngrtz                 float64
	Grtspc, Grtwt                float64
	Prlnsec                      int
	Term                         string
}


func PurlinInit(dl, ll, slope, spacing, span float64, psfs []float64, mono bool) (fdl, fll, fulg, rafterlen, rise, purlinspc float64, idx int){
	var nprz float64
	//var l float64
	if mono{
		rise = slope * span
		rafterlen = math.Sqrt(math.Pow(rise, 2) + math.Pow(span, 2))
		nprz = math.Ceil(rafterlen/1.4)
		purlinspc = math.Round(1000.0*rafterlen/nprz)/1000.0
	} else {
		rise = slope * span/2.0
		rafterlen = math.Sqrt(math.Pow(rise, 2) + math.Pow(span/2.0, 2))
		log.Println("rafterlen->",rafterlen)
		nprz = math.Ceil(rafterlen/1.4)
		purlinspc = math.Round(1000.0*rafterlen/nprz)/1000.0
	}
	prdl := purlinspc * dl; prll := purlinspc * ll
	var pul, dlsf float64
	if len(psfs) == 0{pul = 1.5 * prdl + 1.5 * prll; dlsf = 1.5} else {pul = psfs[0] * prdl + psfs[1] * prll; dlsf = psfs[0]}
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	sheet := filepath.Join(basepath,"../data/steel/","cfszsec.csv")
	csvfile, err := os.Open(sheet)
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(csvfile)
	//log.Println("purlin spacing-",purlinspc, "prdl kn/m", prdl, "prll kn/m",prll)
	for i, val := range []float64{1.83 , 2.44 , 3.05 , 3.66 , 4.27 , 4.88 , 5.49 , 6.1} {
		if val > spacing/2.0{
			idx = i
		}
	}
	var mul, w, wprz float64
	for i := 0; i < df.Nrow() - 1; i++{ //
		//log.Println("section->",df.Elem(i,0))
		//log.Println("mass->",df.Elem(i,16))
		//log.Println("pul->",pul)
		wprz = df.Elem(i,16).Float()*9.81/1e3 
		w = wprz*dlsf + pul
		//log.Println("w ->", w)
		mul = w * math.Pow(spacing/2.0,2)/10.0
		//log.Println("checking moment for braced length->",df.Elem(i,23+idx),"vs",mul, "OK?",df.Elem(i,23+idx).Float() > mul)
		if df.Elem(i,23+idx).Float() > mul{idx = i; break}
	}
	//fmt.Println("section->",df.Elem(idx,0))
	//fmt.Println("mass->",df.Elem(idx,16))
	//get final load in kn/m of portal frame
	//gable end load is that/2
	var llsf float64
	fdl = wprz * nprz * spacing/rafterlen + spacing * dl 
	if len(psfs) == 0{llsf = 1.5} else {llsf = psfs[1]}
	fll = ll * spacing
	fulg = fdl*dlsf + fll* llsf + dlsf * 0.1 * spacing
	log.Println(fdl, fll, fulg)
	return 
}

func FrpsPrelim(f *Portal, fulg float64) (em, cp [][]float64, bmwt, colwt float64){
	//use pinned base formula from sci p 399
	//TODO check for alt formulae
	var l float64
	if f.Mono{l = f.Span} else {l = f.Span/2.0}
	theta := f.Rise/f.Height; m := 1.0 + theta
	k := f.Height/l/1.5
	b := 2.0 * (1.0 + k) + m; c := 1.0 + 2.0 * m
	n := b + m * c
	//moment at eaves
	me := fulg * math.Pow(l,2) * (3.0 + 5.0 * m)/16.0/n
	//moment at apex
	ma := fulg * math.Pow(l,2)/8.0 + m * me
	log.Println("moment at apex->",ma,"at eaves->",me, "kn-m")
	pbc := 165.0
	sxxr := ma*1e3/pbc
	log.Println("sec modulus req->",sxxr,"mm4")
	_, based, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(based)
	sheet := filepath.Join(basepath,"../data/steel/isteel","ISMZ.csv")
	csvfile, err := os.Open(sheet)
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(csvfile)
	var bmsecs, colsecs []int
	var mbm, mcol int
	bmwt = 0.; colwt = 0.
	for i := 0; i < df.Nrow() -1; i++{
		if len(bmsecs) > 10 {break}
		if df.Elem(i,11).Float() >= sxxr{
			bmsecs = append(bmsecs, i)
			//log.Println("section->", df.Elem(i,0), "weight kg/m->",df.Elem(i,1))
			if bmwt == 0{
				mbm = i
				bmwt = df.Elem(i,1).Float()
			} else if bmwt > df.Elem(i,1).Float(){
				mbm = i
				bmwt = df.Elem(i,1).Float()
			}
		}
	}
	log.Println("min bm->",bmwt, df.Elem(mbm,0), df.Elem(mbm,8))
	for i:=0; i < df.Nrow()-1;i++{
		if len(colsecs) > 10{break}
		if df.Elem(i,8).Float() >= df.Elem(mbm,8).Float()*1.5{
			colsecs = append(colsecs, i)
			if colwt == 0 {
				mcol = i
				colwt = df.Elem(i, 1).Float()
			} else if colwt > df.Elem(i, 1).Float(){
				mcol = i
				colwt = df.Elem(i,1).Float()
			}
		}
	}
	log.Println("min col->",colwt, df.Elem(mcol,0), df.Elem(mcol,8), df.Elem(mcol,8).Float()/df.Elem(mbm,8).Float())
	cp = make([][]float64, 2)
	for i:=0; i < 2; i++{
		cp[i] = make([]float64,2)
	}
	cp[0][0] = df.Elem(mcol,2).Float()*1e-4; cp[0][1] = df.Elem(mcol,7).Float()*1e-8
	cp[1][0] = df.Elem(mbm,2).Float()*1e-4; cp[1][1] = df.Elem(mbm,7).Float()*1e-8
	em = [][]float64{{2.1e8}}
	return
}

func PortalInit(f *Portal){
	//col sec 1, bmsec 2
	var coords, cp, em [][]float64
	var mprp, msup [][]int
	var x, y, fdl, fll, fulg, bmwt, colwt float64
	f.Nframes = int(f.Ly/f.Spacing)
	f.Spacing = f.Ly/float64(f.Nframes)
	f.Nframes += 1
	log.Println(f.Nframes)
	fdl, fll, fulg, f.LR, f.Rise, f.Prlnspc, f.Prlnsec = PurlinInit(f.DL, f.LL, f.Slope, f.Spacing, f.Span, f.PSFs, f.Mono)
	if f.W == 0{
		log.Println("w frame kn/m->",fulg)
		em, cp, bmwt, colwt = FrpsPrelim(f, fulg)
	}
	log.Println("em, cp",em, cp)
	bmwt = bmwt * 9.81/1e3; colwt = colwt * 9.81/1e3
	log.Println("bm wt, col wt", bmwt, colwt, "fdl,fll->", fdl, fll)
	var dlsf, llsf float64
	if len(f.PSFs) == 0{dlsf, llsf = 1.5, 1.5} else {dlsf, llsf = f.PSFs[0], f.PSFs[1]}
	fulg -= 0.1 * dlsf * f.Spacing
	fdl += bmwt
	angle := math.Atan(f.Slope)
	pd, cpos, cneg := wltable5(f.Vz, f.Height, f.Span, f.Slope, f.Cpi)
	//em := [][]float64{{210000}}
	coords = append(coords, []float64{x,y})
	var bmvec, colvec []int
	for i := 0; i < f.Nbays; i++{
		switch i{
			case 0:
			x += f.Span; y += f.Height
			coords = append(coords, []float64{x,y-f.Height})
			coords = append(coords, []float64{x-f.Span,y})
			coords = append(coords, []float64{x,y})
			coords = append(coords, []float64{x/2.0,y+f.Span*f.Slope/2.0})
			baymem := [][]int{{1,3,1,1,0},{2,4,1,1,0},{3,5,1,2,0},{4,5,1,2,0}}
			bmvec = append(bmvec, []int{3,4}...)
			colvec = append(colvec, []int{1,2}...)
			mprp = append(mprp, baymem...)
			if f.Fixbs{
				msup = append(msup, []int{1,-1,-1,-1})
				msup = append(msup, []int{2,-1,-1,-1})
			} else {
				msup = append(msup, []int{1,-1,-1,0})
				msup = append(msup, []int{2,-1,-1,0})				
			}
			default:
			x += f.Span
			coords = append(coords, []float64{x,y-f.Height})
			coords = append(coords, []float64{x, y})
			coords = append(coords, []float64{x-f.Span/2.0,y+f.Span*f.Slope/2.0})
			idx := 6 + (i-1)*3
			baymem := [][]int{{idx, idx + 1, 1, 1, 0},{idx+1, idx+2, 1,2,0},{idx-2, idx+2, 1, 2,0}}
			mprp = append(mprp, baymem...)
			bmvec = append(bmvec, []int{len(mprp),len(mprp)-1}...)
			colvec = append(colvec, len(mprp)-2)
			if f.Fixbs{
				msup = append(msup, []int{idx, -1,-1,-1})
			} else {
				msup = append(msup, []int{idx, -1,-1, 0})
			}
		}
	}
	var msloads, jsloads [][]float64
	var dirfac float64
	for _, bm := range bmvec{
		if mprp[bm-1][1] - mprp[bm-1][0] == 1{dirfac = -1.0} else {dirfac = 1.0}
		msloads = append(msloads, []float64{float64(bm), 3, dirfac * dlsf * fdl*math.Cos(angle), 0, 0, 0, 1})
		msloads = append(msloads, []float64{float64(bm), 3, dirfac * llsf * fll*math.Cos(angle), 0, 0, 0, 2})
		msloads = append(msloads, []float64{float64(bm), 3, dlsf * fdl*math.Sin(angle), 0, 0, 0, 1})
		msloads = append(msloads, []float64{float64(bm), 3, llsf * fll*math.Sin(angle), 0, 0, 0, 2})
	}
	//init mem sizes - column = 1.0 to 1.5 times beam stiffness/inertia Ixx
	for _, col := range colvec{
		msloads = append(msloads, []float64{float64(col), 6, dlsf * colwt, 0, 0, 0, 1})
	}
	mod := &kass.Model{
		Ncjt:3,
		Cmdz:[]string{"2df","mks","1"},
		Coords:coords,
		Mprp:mprp,
		Supports:msup,
		Msloads:msloads,
		Jloads:jsloads,
		Em:em,
		Cp:cp,
	}
	frmrez, err := kass.CalcFrm2d(mod,3)
	if err != nil{
		log.Println("ERRORE,errore->",err)
		return
	}
	report, _ := frmrez[6].(string)
	js, _ := frmrez[0].(map[int]*kass.Node)
	ms, _ := frmrez[1].(map[int]*kass.Mem)
	fmt.Println(report)
	pd = pd/1e3
	log.Println("summary of loads on frame->")
	//log.Println("purlins->",dlsf * wprz * nprz/rafterlen, "kn-m2")
	log.Println("fdl->",fdl, "kn-m")
	log.Println("fll->", fll, "kn-m")
	log.Println("fulg->", fulg, "kn-m")
	fmt.Println("design pressure pd-",pd,"kn/m2")
	fmt.Println("design wind load cases +ve cpi ", f.Cpi)
	fmt.Println(cpos)
	fmt.Println("design wind load cases -ve cpi ", -f.Cpi)
	fmt.Println(cneg)
	txtplot := kass.DrawMod2d(mod, ms, f.Term)
	fmt.Println(txtplot)
	for _, node := range js{
		if node.React[1] != 0{
			fmt.Println(node.React)			
			/*
			fmt.Println("footing design")
			colx := 0.45; coly := 0.45; df := 0.0; eo := 0.25
			fck := 25.0; fy := 500.0
			sbc := 100.0; pgck := 24.0; pgsoil := 15.0; nomcvr := 0.06; dmin := 0.25
			pus := []float64{node.React[1]}
			mxs := []float64{0}
			mys := []float64{0}
			psfs := []float64{1.0,1.0}
			shape := "square"
			sloped := true
			dlfac := false
			//mosh.FtngDzRojas(colx, coly, fck, fy, df, dmin, eo, sbc, pgck, pgsoil, nomcvr, pus, mxs, mys, psfs, shape, sloped, dlfac, "dumb")
			*/
		}
	}
}

func HassanEx5(){
	mod := &kass.Model{
		Ncjt:3,
		Cmdz:[]string{"2df","mks","1"},
		Coords:[][]float64{
			{0,0},
			{40,0},
			{0,5},
			{40,5},
			{20,10.36},
		},
		Supports:[][]int{{1,-1,-1,0},{2,-1,-1,0}},
		Mprp:[][]int{{1,3,1,1,0},{2,4,1,1,0},{3,5,1,2,0},{4,5,1,2,0}},
		Jloads:[][]float64{},
		Msloads:[][]float64{{3,3,10.32,0,0,0},{3,6,2.77,0,0,0},{4,3,10.32,0,0,0},{4,6,2.77,0,0,0},},
		Em:[][]float64{{2.1e8}},
		Cp:[][]float64{{105e-4,47540e-8},{129e-4,61520e-8}},
	}
	frmrez, err := kass.CalcFrm2d(mod, 3)
	if err != nil{
		log.Println(err)
		return
	}
	report, _ := frmrez[6].(string)
	fmt.Println(report)
}
/*

	//log.Println(coords)
	//log.Println(mprp)
	//log.Println(bmvec)
	//log.Println(msup)
	//log.Println("bmvec-",bmvec)
	//log.Println("colvec-",colvec)
	//pltchn := make(chan string, 1)
	//go kass.PlotGenTrs(coords, mprp, pltchn)
	//pltstr := <-pltchn
	//log.Println(pltstr)

   
	lspan := f.Span/2.0
	ly := lspan
	ty := 1.0
	lbr := 200.0
	tbr := 20.0
	nsecs := 1
	grd := 43
	sectyp := 0
	brchck := false
	yeolde := false
	ldcases := [][]float64{{1.0,3.0,wdl,0,0,0,1},{1.0,3.0,wll,0,0,0,1}}
	rez := StlBmDBs(lspan, ly, ty, lbr, tbr, ldcases, sectyp, grd, nsecs, brchck, yeolde)
	bdx := rez[0]
	df := StlSecBs(sectyp)
	wb, arb, ixb := df.Elem(bdx,2).Float(), df.Elem(bdx,23).Float(), df.Elem(bdx,11).Float()
	fil := df.Filter(
		dataframe.F{Colname:"ix", Comparator:series.Greater, Comparando:1.5*ixb},
	)
	//fmt.Println(fil.Nrow())
	//fmt.Println(fil.Subset([]int{fil.Nrow()-1}))
	cdx := fil.Nrow()-1
	wc, arc, ixc := df.Elem(cdx,2).Float(), df.Elem(cdx,23).Float(), df.Elem(cdx,11).Float()
	//log.Println(wc, arc, ixc, wb, arb, ixb)
	//START WITH SAME SECTION
	cp = [][]float64{{arc*1e-4, ixc*1e-8},{arb*1e-4, ixb*1e-8}}
	em = [][]float64{{2.1*1e8}}
	//log.Println("section->",df.Elem(i,1))
	//log.Println("depth, web thickness->",df.Elem(i,3), df.Elem(i,6))
	//log.Println("area, zx, zy->",df.Elem(i,23),df.Elem(i,15), df.Elem(i,16))
	//log.Println("rx, ry->",df.Elem(i,13), df.Elem(i,14))
*/
