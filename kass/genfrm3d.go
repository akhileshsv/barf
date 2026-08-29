package barf

import (
	//"sort"
	"errors"
	"fmt"
	"log"
	"math"
	//"strconv"
	//"strings"
	//"encoding/json"
	//"github.com/gonum/floats"
)

//Frm3d generation struct
//it is woefully unused
type Frm3d struct {
	Id             int
	Title          string
	Foldr          string
	Code           int
	X, Y, Z        []float64
	DL             float64
	LL             float64
	WL             []float64
	EL             []float64
	CL             []float64
	Wtwl           float64
	Sections       [][]float64
	Styps          []int
	Csec, Bsec     []int
	Bysec          []int
	Stair          []int
	Bystyp         int
	Cstyp, Bstyp   int
	Clstyp         int
	Fcks           []float64
	Fys            []float64
	Fyds           []float64
	Lspans         []float64
	Lbays          []float64
	DM             float64
	PSFs           []float64
	WLFs           []float64
	Clvrs          [][]float64
	Cldim          [][]float64
	Clsec          []int
	Term           string
	Verbose        bool
	Spam           bool
	Fltslb         bool
	Lclvr, Rclvr   bool
	Braced         bool
	Fixbase        bool
	Nloads         int
	Ldcalc         int
	Wlcalc         int
	Selfwt         int
	Mtyp           int
	Noclvr         bool
	Fgen           int
	Bmrel          int
	Byrel          int
	Nomcvr         float64
	D1, D2         float64
	Dslb           float64
	Fyv            float64
	Bfcalc         bool                        //flange width calculation
	Bfycalc        bool                        //either x is flanged or y is? calculation
	Pg             float64                     `json:",omitempty"`
	Matprop        []float64                   `json:",omitempty"`
	Lcmem, Rcmem   []int                       `json:",omitempty"`
	Tcmem, Bcmem   []int                       `json:",omitempty"`
	Ljs, Rjs       []int                       `json:",omitempty"`
	Nlp            int                         `json:",omitempty"`
	Cdim, Bdim     [][]float64                 `json:",omitempty"`
	Bydim          [][]float64                 `json:",omitempty"`
	Tyb            float64                     `json:",omitempty"`
	Tyby           float64                     `json:",omitempty"`
	Secmap         map[int][]float64           `json:",omitempty"`
	Mod            Model                       `json:",omitempty"`
	Vec            [][]int                     `json:",omitempty"`
	Val            [][]float64                 `json:",omitempty"`
	Flridx         []int                       `json:",omitempty"`
	Txtplots       []string                    `json:",omitempty"`
	DMs            []float64                   `json:",omitempty"`
	Bloads         [][]float64                 `json:",omitempty"`
	Cloads         [][]float64                 `json:",omitempty"`
	Jloads         [][]float64                 `json:",omitempty"`
	Wloadl         [][]float64                 `json:",omitempty"`
	Wloadr         [][]float64                 `json:",omitempty"`
	Nodemap        map[int][][]int             `json:",omitempty"`
	Nodes          map[Pt3d][]int                `json:"-"`
	Members        map[int][][]int             `json:",omitempty"`
	Cols, Beams    []int                       `json:",omitempty"`
	Beamys         []int                       `json:",omitempty"`
	Mloadmap       map[int][][]float64         `json:",omitempty"`
	Jloadmap       map[int][][]float64         `json:",omitempty"`
	Advloads       map[int][][]float64         `json:",omitempty"`
	Benloads       map[int][][]float64         `json:",omitempty"`
	Sadloads       map[int][][]float64         `json:",omitempty"`
	Wadloads       map[int][][]float64         `json:",omitempty"`
	Uniloads       map[int][][]float64         `json:",omitempty"`
	Loadcons       map[int][][]float64         `json:",omitempty"`
	Bmenv          map[int]*BmEnv              `json:",omitempty"`
	Colenv         map[int]*ColEnv             `json:",omitempty"`
	Secprop        []Secprop                   `json:",omitempty"`
	Ms             map[int]*Mem                `json:",omitempty"`
	Mslmap         map[int]map[int][][]float64 `json:",omitempty"`
	CBeams         [][]int                     `json:",omitempty"`
	CBeamys        [][]int                     `json:",omitempty"`
	Reports        []string                    `json:",omitempty"`
	Opt            int                         `json:",omitempty"`
	Dconst         bool                        `json:",omitempty"`
	Width          float64                     `json:",omitempty"`
	Cvec           [][]float64                 `json:",omitempty"`
	Bvec           [][]float64                 `json:",omitempty"`
	Cdims          [][]float64                 `json:",omitempty"`
	Bdims          [][]float64                 `json:",omitempty"`
	Bdimy          [][]float64                 `json:",omitempty"`
	Ncols          int                         `json:",omitempty"`
	Nbms, Nbmy     int                         `json:",omitempty"`
	Ncls           int                         `json:",omitempty"`
	Nspans         int                         `json:",omitempty"`
	Nflrs          int                         `json:",omitempty"`
	Nbays          int                         `json:",omitempty"`
	Brel           int                         `json:",omitempty"`
	Df             float64                     `json:",omitempty"`
	Cmap           map[int][]int               `json:",omitempty"`
	Dzslb          bool                        `json:",omitempty"`
	Dz             bool                        `json:",omitempty"`
	Report         string                      `json:",omitempty"`
	Rez            []interface{}               `json:",omitempty"`
	Kostin         []float64                   `json:",omitempty"`
	Quants         []float64                   `json:",omitempty"`
	Kosts          []float64                   `json:",omitempty"`
	Kost           float64                     `json:",omitempty"`
	Plinth         int                         `json:",omitempty"`
	SlbDL, SlbLL   []float64                   `json:",omitempty"`
	Mstr           map[string]int              `json:",omitempty"`
	Slbnodes       [][]int                     `json:",omitempty"`
	Slbmems        [][]int                     `json:",omitempty"`
	Slbendc        [][]int                     `json:",omitempty"`
	Slbgrid        [][]int                     `json:",omitempty"`
	Mmap           map[tupil]int               `json:",omitempty"`
	Gstr           [][]int                     `json:",omitempty"`
	Gcs            [][]float64                 `json:",omitempty"`
}

//Inits a 3d Frm and is also the path to madness
func (f *Frm3d) Init() (err error) {
	if f.Fgen == 0 {
		//single column and beamx, beamy section (nsections = 3)
		f.Fgen = 1
	}
	f.Noclvr = true
	if (len(f.Sections) == 0 || len(f.Csec) == 0 || len(f.Bsec) == 0 || len(f.Bysec) == 0) && f.Opt == 0 {
		return errors.New("no sections specified")
	}
	if len(f.X) == 0 {
		return errors.New("frame spans (x) not specified")
	}
	if len(f.Y) == 0 {
		return errors.New("frame floor heights (y) not specified")
	}
	if len(f.Z) == 0 {
		return errors.New("frame bays (z) not specified")
	}
	if f.Mtyp == 0 {
		f.Mtyp = 1
	}
	//is code, bs code, (lol) euro code
	if f.Code == 0 {
		f.Code = 1
	}
	if f.PSFs == nil{
		switch f.Code {
		case 1:
			f.PSFs = []float64{1.5, 1.0, 1.5, 0.0}
		case 2:
			f.PSFs = []float64{1.4, 1.0, 1.6, 0.0}
		}
	}
	if f.WLFs == nil {
		f.WLFs = []float64{1.2, 1.2, 1.2}
	}
	f.Nodemap = make(map[int][][]int) //now tis map[int] [mtyp-1][m1,m2]
	f.Nodes = make(map[Pt3d][]int)
	f.Members = make(map[int][][]int)
	f.Advloads = make(map[int][][]float64)
	f.Benloads = make(map[int][][]float64)
	f.Loadcons = make(map[int][][]float64)
	f.Uniloads = make(map[int][][]float64)
	f.Jloadmap = make(map[int][][]float64)
	f.Wadloads = make(map[int][][]float64)
	f.Sadloads = make(map[int][][]float64)
	f.Secmap = make(map[int][]float64)
	if len(f.X) == 1 && f.Nspans != 0 {
		var x, xstep float64
		xstep = f.X[0] //not your average "xstep"
		if xstep == 0 {
			return errors.New("frame spans (x) not specified")
		}
		f.X = []float64{x}
		for i := 0; i < f.Nspans; i++ {
			x += xstep
			f.X = append(f.X, x)
		}
	}
	if len(f.Y) == 1 && f.Nflrs != 0 {
		var y, ystep float64
		ystep = f.Y[0]
		if ystep == 0 {
			return errors.New("frame floors (y) not specified")
		}
		f.Y = []float64{y}
		for i := 0; i < f.Nflrs; i++ {
			y += ystep
			f.Y = append(f.Y, y)
		}

	}
	if len(f.Z) == 1 && f.Nbays != 0 {
		var z, zstep float64
		zstep = f.Z[0]
		if zstep == 0 {
			return errors.New("frame bays (z) not specified")
		}
		f.Z = []float64{z}
		for i := 0; i < f.Nbays; i++ {
			z += zstep
			f.Z = append(f.Z, z)
		}

	}
	f.Ncols = len(f.X)
	f.Nbms = len(f.X) - 1
	f.Nbmy = len(f.Z) - 1
	if len(f.Clvrs) >= 3{
		
		if len(f.Clvrs[0]) > 1 && f.Clvrs[0][0] > 0.0 {
			f.Ncls++
		}
		if len(f.Clvrs[0]) > 1 && f.Clvrs[1][0] > 0.0 {
			f.Ncls++
		}
		if len(f.Clvrs[0]) > 1 && f.Clvrs[2][0] > 0.0 {
			f.Ncls++
		}
		if len(f.Clvrs[0]) > 1 && f.Clvrs[3][0] > 0.0 {
			f.Ncls++
		}
	} else {
		f.Clvrs = make([][]float64, 4)
		for i := range f.Clvrs{
			f.Clvrs[i] = make([]float64, 4)
		}
	}
	if f.Nflrs == 0 {
		f.Nflrs = len(f.Y) - 1
	}
	if f.Nspans == 0 {
		f.Nspans = len(f.X) - 1
	}
	
	f.Slbnodes = [][]int{}
	f.Mmap = make(map[tupil]int)
	if f.Nbays == 0 {
		f.Nbays = len(f.Z) - 1
	}
	if f.Styps == nil || len(f.Styps) == 0 {
		f.Styps = make([]int, len(f.Sections))
		for i := range f.Sections {
			f.Styps[i] = 1
		}
	}
	for i := 1; i < len(f.X); i++ {
		f.Lspans = append(f.Lspans, f.X[i]-f.X[i-1])
	}
	for i := 1; i < len(f.Z); i++ {
		f.Lbays = append(f.Lbays, f.Z[i]-f.Z[i-1])
	}

	//define material prop
	if f.Mtyp > 0 {
		f.InitMat()
	}
	if f.Mtyp == 1 && f.Fyv == 0 {
		f.Fyv = f.Fys[2]
	}
	if f.Bfcalc{
		f.CalcBf3d()
	}
	//fmt.Println(f.Sections, f.Styps)
	return
}

func (f *Frm3d) InitMat(){
	//add G, poissons ratio etc
	switch f.Mtyp{
		case 1:
		//fmt.Println("rcc frame")
		f.Pg = 25.0 
		//if f.Code == 2{f.Pg = 24.0}
		if f.Fcks == nil && f.Fys == nil {
			//defaults - m25, fe500, fe500 (who is even selling fe415 lmao)
			f.Fcks = append(f.Fcks, 25.0)
			f.Fys = append(f.Fys, 500.0)
			f.Fyds = append(f.Fyds, 500.0)
		}
		if len(f.Fcks) == 1 || len(f.Fcks) < 4{
			//col, beam, footing, slab
			f.Fcks = []float64{f.Fcks[0],f.Fcks[0],f.Fcks[0],f.Fcks[0]}
		}
		if len(f.Fys) == 1 || len(f.Fys) < 4{
			f.Fys = []float64{f.Fys[0],f.Fys[0],f.Fys[0],f.Fys[0]}
		}
		if len(f.Fyds) == 1 || len(f.Fyds) < 4{
			f.Fyds = []float64{f.Fys[0],f.Fys[0],f.Fys[0],f.Fys[0]}
		}
		if f.Fyv == 0.0{
			//isnt shear reinf restricted to max 415 MPa in IS and 500 in BS? 
			f.Fyv = 415.0
		}
		f.Mod.Em = make([][]float64,2)
		poi_mu := 1.6
		e1 := FckEm(f.Fcks[0]); g1 := e1/(2.0 * (1.0 + poi_mu))
		e2 := FckEm(f.Fcks[1]); g2 := e2/(2.0 * (1.0 + poi_mu))
		f.Mod.Em[0] = []float64{e1, g1}
		f.Mod.Em[1] = []float64{e2, g2}
		f.Mod.Wng = make([][]float64, 2)
		f.Mod.Wng[0] = []float64{0,0}
		f.Mod.Wng[1] = []float64{1,90}
		case 2:
		//steel frame
		f.Pg = 7850.0 //or 7850? FINALIZE STEEL UNITZ
		f.Mod.Em = [][]float64{{2e5}}
	}
}

func (f *Frm3d) Calc() (err error) {
	err = f.Init()
	if err != nil{
		return
	}
	f.GenCoords()
	f.GenMprp()
	//f.GenLoads()
	return
}

func (f *Frm3d) AddNode(x,y,z float64) (err error) {
	p := Pt3d{X: x, Y: y, Z: z}
	if _, ok := f.Nodes[p]; ok {
		return errors.New("node exists")
	} else {
		f.Mod.Coords = append(f.Mod.Coords, []float64{x, y, z})
		ndx := len(f.Mod.Coords)
		f.Nodes[p] = []int{ndx}
		f.Nodemap[ndx] = make([][]int,7)
	}
	return
}

func (f *Frm3d) GenCoords() {
	for _, z := range f.Z {
		for _, y := range f.Y {
			for _, x := range f.X {
				f.AddNode(x,y,z)
			}
		}
	}
}

func (f *Frm3d) AddMem(xstep, ystep, colstep, jb, je, mtyp int) (err error) {
	//adds a member to model
	//mtyps - col, bm, lclvr, rclvr (1,2,3,4)
	var mrel int
	if jb < 0 || jb > len(f.Mod.Coords) || je < 0 || je > len(f.Mod.Coords) {
		return fmt.Errorf("invalid member node(s)-%v,%v", jb, je)
	}
	xdx, ydx, flrdx, locdx := getxydx(jb, colstep, xstep, ystep)
	mt, locx, lex := getmemtype(jb, colstep, xstep, ystep, xdx, ydx, locdx, flrdx, mtyp)
	mdx := len(f.Mod.Mprp) + 1
	var em, cp int
	switch mtyp{
		case 1:
		f.Cols = append(f.Cols, mdx)
		em = 1
		if len(f.Csec) >= xdx{
			cp = f.Csec[xdx-1]
		} else {
			cp = f.Csec[0]
		}
		case 2:
		mrel = f.Bmrel
		f.Beams = append(f.Beams, mdx)
		
		if len(f.Bsec) >= xdx{
			cp = f.Bsec[xdx-1]
		} else {
			cp = f.Bsec[0]
		}
		em = 2
		case 3:
		mrel = f.Byrel
		f.Beamys = append(f.Beamys, mdx)
		em = 2
		if len(f.Bysec) >= ydx{
			cp = f.Bysec[ydx-1]
		} else {
			cp = f.Bysec[0]
		}
		case 4:
		f.Lcmem = append(f.Lcmem, mdx)
		em = 2
		//cp = f.Clsec[0]
		case 5:
		f.Rcmem = append(f.Rcmem, mdx)
		em = 2
		//cp = f.Clsec[1]
		case 6:
		f.Tcmem = append(f.Tcmem, mdx)
		em = 2
		//cp = f.Clsec[2]
		case 7:
		f.Bcmem = append(f.Bcmem, mdx)
		em = 2
		//cp = f.Clsec[3]
	}

	f.Mod.Mprp = append(f.Mod.Mprp,[]int{jb, je, em, cp, mrel})
	if mtyp == 1 {
		f.Mod.Wng = append(f.Mod.Wng, []float64{0,0})
	} else {
		f.Mod.Wng = append(f.Mod.Wng, []float64{1,90})
	}
	f.Secmap[mdx] = f.Sections[cp-1]
	f.Members[mdx] = append(f.Members[mdx],[]int{jb, je, em, cp, mrel})
	f.Members[mdx] = append(f.Members[mdx],[]int{mtyp, xdx, ydx, flrdx, locx, lex, mt})
	f.Nodemap[jb][mtyp-1] = append(f.Nodemap[jb][mtyp-1],mdx)
	f.Nodemap[je][mtyp-1] = append(f.Nodemap[je][mtyp-1],mdx)
	f.Mmap[tupil{jb,je}] = mdx
	return
}

func (f *Frm3d) GenMprp() {
	colstep := (f.Nspans + 1) * (f.Nbays + 1)
	xstep := f.Nspans + 1
	ystep := f.Nflrs+1
	f.Mod.Ncjt = 6
	log.Println("steps->",colstep,xstep,ystep)
	var sup []int
	switch f.Fixbase { //footing vec
	case true:
		//fixed footing
		sup = []int{-1, -1, -1, -1, -1, -1}
	case false:
		//WOT
	}
	for i := 1; i <= len(f.Mod.Coords); i++ {
		if i <= colstep{
			//add supports
			supi := append([]int{i}, sup...)
			f.Mod.Supports = append(f.Mod.Supports, supi)
		}
		if i+colstep <= len(f.Mod.Coords) {
			//add col		
			f.AddMem(xstep, ystep, colstep,i, i+colstep, 1)
		}
		if i > colstep {
			if i+1 <= len(f.Mod.Coords) && f.Mod.Coords[i-1][1] == f.Mod.Coords[i][1] {
				//beam x
				f.AddMem(xstep, ystep, colstep,i, i+colstep, 2)
			}
			if i+xstep <= len(f.Mod.Coords) && f.Mod.Coords[i][2] == f.Mod.Coords[i+xstep][2] {
				//beam y
				f.AddMem(xstep, ystep, colstep,i, i+xstep, 3)
			}
		}
		if i+xstep+1 <= len(f.Mod.Coords) && f.Mod.Coords[i][2] == f.Mod.Coords[i+xstep+1][2] {
			f.AddSlb(xstep, i)
		}
	}
	for idx, sect := range f.Sections {
		styp := 1
		if len(f.Styps) == len(f.Sections){
			styp = f.Styps[idx]
		}
		bar := CalcSecProp(styp, sect)
		f.Secprop = append(f.Secprop, bar)
		f.Mod.Cp = append(f.Mod.Cp, []float64{bar.Area/1e6, bar.Ixx/1e12, bar.Iyy/1e12, bar.J/1e12})
	}
	//check for self weight calc
	switch f.Selfwt{
	//load vec - [self wt, nload cases]
		case 1,2:
		err := f.AddSelfWeight(f.Selfwt)
		if err != nil{
			log.Println("ERRORE, errore->",err)
		}
	}
}

func (f *Frm3d) GenFrm(){
	//generate from ground struct
	ncols := len(f.Gcs)
	bcs := make([][]float64, len(f.Gcs))
	for i, pt := range f.Gcs{
		bcs[i] = []float64{pt[0],0.0,pt[1]}
		//f.Mod.Coords = append(f.Mod.Coords, bcs[i])
	}
	ystep := 4.0; y := 0.0
	for i := 0; i < f.Nflrs; i++{
		for _, pt := range bcs{
			p2 := []float64{pt[0],pt[1]+y,pt[2]}
			f.Mod.Coords = append(f.Mod.Coords, p2)
		}
		y += ystep
	}
	for i := range f.Mod.Coords{
		idx := i 
		jb := i + 1
		if i >= ncols{
			idx = i - ncols 
		}
		if i + ncols < len(f.Mod.Coords){
			fmt.Println("adding col->",jb,i+ncols+1)
		}
		for _, j := range f.Gstr[idx]{
			je := j
			fmt.Println("adding bm->",jb, je)
		} 
	}
}

func (f *Frm3d) AddSelfWeight(dx int) (err error){
	//add beam self weight (ltyp 3)
	var wdl float64
	//log.Println("ADDING WDL")
	for _, bm := range f.Beams{
		bdx := f.Members[bm][0][3]
		bstyp := f.Styps[bdx-1]
		bdim := f.Sections[bdx-1]
		switch f.Mtyp{
			case 1:
			//rcc frame
			switch bstyp{
				case 1:
				b := bdim[0]; d := bdim[1]
				wdl = f.Pg * b * (d - f.Dslb)/1e6
				case 6,7,8,9,10:
				bdim := f.Sections[bdx-1]; bw := bdim[2]; dw:= bdim[1] - bdim[3]
				wdl = f.Pg * bw * dw /1e6			
				default:
				//HAHA.
				wdl = f.Pg * f.Secprop[bdx-1].Area/1e6
				
			}
			default:
			wdl = f.Pg * f.Secprop[bdx-1].Area/1e6
			log.Println(ColorYellow,"bdx,wdl->",bdx,wdl,ColorReset)
		}
		ldcase := []float64{1.0, 3.0, wdl, 0.0, 0.0, 0.0, 1.0}
		err = f.AddMemLoad(bm, ldcase)
		if err != nil{
			//log.Println("ERRORE,errore->",err)
			return
		}
	}
	if dx == 2{
		//add column self weight (ltyp 6)
		for _, col := range f.Cols{
			bdx := f.Members[col][0][3]
			wdl := f.Pg * f.Secprop[bdx-1].Area/1e6
			ldcase := []float64{1.0, 6.0, wdl, 0.0, 0.0, 0.0, 1.0}
			err = f.AddMemLoad(col, ldcase)
			if err != nil{
				//log.Println("ERRORE,errore->",err)
				return
			}
		}
	}
	return
}

func (f *Frm3d) AddMemLoad(mem int, ldcase []float64) (err error){
	if _, ok := f.Members[mem]; !ok{
		return errors.New("invalid member index")
	}
	ldcat := ldcase[6]
	var w1a, w2a, w1b, w2b float64
	switch ldcat{
		case 1.0:
		w1a = f.PSFs[0]*ldcase[2]
		w2a = f.PSFs[0]*ldcase[3]
		w1b = f.PSFs[1]*ldcase[2]
		w2b = f.PSFs[1]*ldcase[3]
		f.Uniloads[1] = append(f.Uniloads[1], []float64{float64(mem), ldcase[1], ldcase[2], ldcase[3], ldcase[4], ldcase[5], 1.0})
		case 2.0:
		w1a = f.PSFs[2]*ldcase[2]
		w2a = f.PSFs[2]*ldcase[3]
		w1b = f.PSFs[3]*ldcase[2]
		w2b = f.PSFs[3]*ldcase[3]
		f.Uniloads[2] = append(f.Uniloads[2], []float64{float64(mem), ldcase[1], ldcase[2], ldcase[3], ldcase[4], ldcase[5], 2.0})
		
	}
	f.Advloads[mem] = append(f.Advloads[mem],[]float64{float64(mem),ldcase[1],w1a,w2a,ldcase[4],ldcase[5], ldcat})
	if w1b + w2b > 0.0 {
		f.Benloads[mem] = append(f.Benloads[mem],[]float64{float64(mem),ldcase[1],w1b,w2b,ldcase[4],ldcase[5], ldcat})
	}
	
	if (f.Nloads == 3 || len(f.WL) != 0){
		//build wind load patterns
		//ben loads are usually unity or zero so chill on those, they'll be the same?
		switch ldcat{
			case 1.0:
			w1a = f.WLFs[0]*ldcase[2]
			w2a = f.WLFs[0]*ldcase[3]			
			case 2.0:
			w1a = f.WLFs[1]*ldcase[2]
			w2a = f.WLFs[1]*ldcase[3]
		}
		f.Wadloads[mem] = append(f.Wadloads[mem],[]float64{float64(mem),ldcase[1],w1a,w2a,ldcase[4],ldcase[5], ldcat})
	}
	return
}

func (f *Frm3d) AddBeamLLDL(){
	if f.DL + f.LL > 0.0{
		for _, bm := range f.Beams{
			mvec := f.Members[bm]; mtyp := mvec[1][0]
			switch mtyp{
				case 3:
				//left clvr
				cdl := f.DL; cll := f.LL 
				if f.Clvrs[0][1]+f.Clvrs[0][2] > 0.0{
					cdl = f.Clvrs[0][1]; cll = f.Clvrs[0][2]
				}
				_ = f.AddMemLoad(bm, []float64{1.0,3.0,cdl,0.0,0.0,0.0,1.0})
				if cll > 0.0{
					_ = f.AddMemLoad(bm, []float64{1.0,3.0,cll,0.0,0.0,0.0,2.0})
				}
				case 4:
				//right clvr
				cdl := f.DL; cll := f.LL 
				if f.Clvrs[1][1]+f.Clvrs[1][2] > 0.0{
					cdl = f.Clvrs[1][1]; cll = f.Clvrs[1][2]
				}
				_ = f.AddMemLoad(bm, []float64{1.0,3.0,cdl,0.0,0.0,0.0,1.0})
				if cll > 0.0{
					_ = f.AddMemLoad(bm, []float64{1.0,3.0,cll,0.0,0.0,0.0,2.0})
				}
				default:
				_ = f.AddMemLoad(bm, []float64{1.0,3.0,f.DL,0.0,0.0,0.0,1.0})
				if f.LL > 0.0{
					_ = f.AddMemLoad(bm, []float64{1.0,3.0,f.LL,0.0,0.0,0.0,2.0})
				}
			}
		}
	}
	return
}

func (f *Frm3d) AddSlb(xstep, i int){	
	//add a slab at lower bottom node i
	
	fmt.Println("adding slab->i, xstep->",i, xstep)
	var lidx, ridx, tidx, bidx int
	ljb := i
	lje := i + xstep
	if val, ok := f.Mmap[tupil{ljb,lje}]; ok{
		lidx = val
	}
	rjb := i + 1
	rje := i + xstep + 1
	if val, ok := f.Mmap[tupil{rjb,rje}]; ok{
		ridx = val
	}
	tjb := i + xstep
	tje := i + xstep + 1
	if val, ok := f.Mmap[tupil{tjb,tje}]; ok{
		tidx = val
	}
	bjb := i
	bje := i + 1
	if val, ok := f.Mmap[tupil{bjb,bje}]; ok{
		bidx = val
	}
	f.Slbnodes = append(f.Slbnodes, []int{i, i + 1, i + xstep + 1, i + xstep})
	f.Slbmems = append(f.Slbmems, []int{lidx, ridx, tidx, bidx})
	return
}

/*
func (f *Frm3d) AddSlbLoad(){
	//calculates member loads for 1 and two way slabs as per is 456 clause 24.5
	//to do - add coefficients for continuous supports
	mloads = make(map[string][][]float64)
	mverts = make(map[string]map[int][][]float64)
	for mem := range mems{
		mverts[mem] = make(map[int][][]float64)
	}
	mareas = make(map[string][]float64)
	switch styp{
		case 1:
		//TODO (doingz)
		//one way slab
		for i, sns := range snodes{
			mems := smems[i]
			p1 := cdmap[sns[0]]; p2 := cdmap[sns[1]]; p3 := cdmap[sns[2]]; p4 := cdmap[sns[3]]
			fdx := sns[4]
			lx := geom.EDist(p1,p2); ly := geom.EDist(p2,p3)
			wdl := dl[fdx]; wll := ll[fdx]
			wdl += 0.025 * dused
			switch slbc{
			//case 0:
				//one way if ly/lx or lx/ly >= 2
				case 0, 1:
				//load on beamxs (l 0, r 0, t, b)
				lpt := []float64{(p1[0]+p4[0])/2.0, (p1[1]+p4[1])/2.0,(p1[2]+p4[2])/2.0}
				rpt := []float64{(p2[0]+p3[0])/2.0, (p2[1]+p3[1])/2.0,(p2[2]+p3[2])/2.0}
				verts1 := [][]float64{p1, p2, rpt, lpt}
				area1 := geom.PolyArea(verts1)
				verts2 := [][]float64{lpt, rpt, p3, p4}
				area2 := geom.PolyArea(verts2)
 				frcvec := [][]float64{
					{0,3,0,wdl*ly/2.0,0,0,1,1},
					{0,3,0,wll*ly/2.0,0,0,2,1},
				}
				for j, m := range mems{
					switch j{
						case 0:
						case 1:
						case 2:
						mverts[m][i] = verts1
						mareas[m] = append(mareas[m],[]float64{area1,4.0}...)
						mloads[m] = append(mloads[m],frcvec...)
						case 3:
						mverts[m][i] = verts2
						mareas[m] = append(mareas[m],[]float64{area2,4.0}...)
						mloads[m] = append(mloads[m],frcvec...)
					}
				}

				case 2:
				//load on beamys
				bpt := []float64{(p2[0]+p1[0])/2,(p1[1]+p2[1])/2, (p2[2]+p1[2])/2}
				tpt := []float64{(p3[0]+p4[0])/2,(p3[1]+p4[1])/2 , (p2[2]+p1[2])/2}
				verts1 := [][]float64{p1, bpt, tpt, p4}
				area1 := geom.PolyArea(verts1)
				verts2 := [][]float64{bpt, p2, p3, tpt}
				area2 := geom.PolyArea(verts2)
 				frcvec := [][]float64{
					{0,3,0,wdl*lx/2.0,0,0,1,1},
					{0,3,0,wll*lx/2.0,0,0,2,1},
				}
				for j, m := range mems{
					switch j{
						case 0:
						mverts[m][i] = verts1
						mareas[m] = append(mareas[m],[]float64{area1,4.0}...)
						mloads[m] = append(mloads[m],frcvec...)
						case 1:
						mverts[m][i] = verts2
						mareas[m] = append(mareas[m],[]float64{area2,4.0}...)
						mloads[m] = append(mloads[m],frcvec...)
						case 2:
						case 3:
					}
				}

			}
		}
		case 2:
		//two way slab
		for i  := range snodes{
			sns := snodes[i]
			mems := smems[i]
			p1 := cdmap[sns[0]]; p2 := cdmap[sns[1]]; p3 := cdmap[sns[2]]; p4 := cdmap[sns[3]]
			fdx := sns[4]
			lx := geom.EDist(p1,p2); ly := geom.EDist(p2,p3)
			wdl := dl[fdx]; wll := ll[fdx]
			wdl += 0.025 * dused
			switch{	
				case lx == ly:
				//square panel, all triangles
				//two triangular loads upto lx/2.0
				lpt := []float64{(p1[0]+p4[0])/2.0 + lx/2.0,(p1[1]+p4[1])/2.0,(p1[2]+p4[2])/2.0}
				verts := [][]float64{p1, p2, lpt}
				area := geom.PolyArea(verts)
				frcvec := [][]float64{
					{0,4,0,wdl*lx/2.0,0,lx/2.0,1,1},
					{0,4,0,wll*lx/2.0,0,lx/2.0,2,1},
					{0,4,wdl*lx/2.0,0,lx/2.0,0,1,1},
					{0,4,wll*lx/2.0,0,lx/2.0,0,2,1},
				}
				for j, m := range mems{
					mareas[m] = append(mareas[m],[]float64{area,3.0}...)
					mloads[m] = append(mloads[m],frcvec...)
					switch j{
						case 0:
						mverts[m][i] = [][]float64{p1,p4,lpt}
						case 1:
						mverts[m][i] = [][]float64{p2,p3,lpt}
						case 2:
						mverts[m][i] = [][]float64{p4,lpt,p3}
						case 3:
						mverts[m][i] = [][]float64{p1,p2,lpt}
					}
				}
				case lx > ly:
				//two sets - triangular frc vec -  2 ltyps 4 frcvec1
				//trap frc vec - 2 ltyp 4 until ly/2, 1 ltyp 3 la ly/2 lb ly/2 frcvec2
				lpt := []float64{(p1[0]+p4[0])/2.0 + ly/2.0, (p1[1]+p4[1])/2.0,(p1[2]+p4[2])/2.0}
				rpt := []float64{(p2[0]+p3[0])/2.0 - ly/2.0, (p2[1]+p3[1])/2.0,(p2[2]+p3[2])/2.0}
				frcvec1 := [][]float64{
					{0,4,0,wdl*ly/2.0,0,ly/2.0,1,1},
					{0,4,0,wll*ly/2.0,0,ly/2.0,2,1},
					{0,4,wdl*ly/2.0,0,ly/2.0,0,1,1},
					{0,4,wll*ly/2.0,0,ly/2.0,0,2,1},
				}
				frcvec2 := [][]float64{
					{0,4,0,wdl*ly/2.0,0,ly/2.0,1,1},
					{0,4,0,wll*ly/2.0,0,ly/2.0,2,1},
					{0,4,wdl*ly/2.0,0,ly/2.0,0,1,1},
					{0,4,wll*ly/2.0,0,ly/2.0,0,2,1},
					{0,3,wdl*ly/2.0,0,ly/2.0,ly/2.0,1,1},
					{0,3,wll*ly/2.0,0,ly/2.0,ly/2.0,2,1},
				}
				verts1 := [][]float64{p1, p4, lpt}
				area1 := geom.AreaPoly2(verts1)
				verts2 := [][]float64{p1, p2, rpt, lpt}
				area2 := geom.AreaPoly2(verts2)
				for j, m := range mems{
					switch j {
					case 0:
						mareas[m] = append(mareas[m],[]float64{area1,3.0}...)
						mloads[m] = append(mloads[m], frcvec1...)
						mverts[m][i] = [][]float64{p1, lpt, p4}
					case 1:
						mareas[m] = append(mareas[m],[]float64{area1,3.0}...)
						mloads[m] = append(mloads[m], frcvec1...)
						mverts[m][i] = [][]float64{p2, p3, rpt}
					case 2:
						mareas[m] = append(mareas[m],[]float64{area2, 4.0}...)
						mloads[m] = append(mloads[m], frcvec2...)
						mverts[m][i] = [][]float64{p4, lpt, rpt, p3}
					case 3:
						mareas[m] = append(mareas[m],[]float64{area2, 4.0}...)
						mloads[m] = append(mloads[m], frcvec2...)
						mverts[m][i] = [][]float64{p1, p2, rpt, lpt}
					}
				}		
				case ly > lx:
				bpt := []float64{(p2[0]+p1[0])/2,(p1[1]+p2[1])/2 + lx/2.0, (p2[2]+p1[2])/2}
				tpt := []float64{(p3[0]+p4[0])/2,(p3[1]+p4[1])/2 - lx/2.0, (p2[2]+p1[2])/2}
				verts1 := [][]float64{p1, p2, bpt}
				area1 := geom.AreaPoly2(verts1)
				verts2 := [][]float64{p1, bpt, tpt, p4}
				area2 := geom.AreaPoly2(verts2)
				frcvec1 := [][]float64{
					{0,4,0,wdl*lx/2.0,0,lx/2.0,1,1},
					{0,4,0,wll*lx/2.0,0,lx/2.0,2,1},
					{0,4,wdl*lx/2.0,0,lx/2.0,0,1,1},
					{0,4,wll*lx/2.0,0,lx/2.0,0,2,1},
				}
				frcvec2 := [][]float64{
					{0,4,0,wdl*lx/2.0,0,lx/2.0,1,1},
					{0,4,0,wll*lx/2.0,0,lx/2.0,2,1},
					{0,4,wdl*lx/2.0,0,lx/2.0,0,1,1},
					{0,4,wll*lx/2.0,0,lx/2.0,0,2,1},
					{0,3,wdl*lx/2.0,0,lx/2.0,lx/2.0,1,1},
					{0,3,wll*lx/2.0,0,lx/2.0,lx/2.0,2,1},
				}
				for j, m := range mems{
					switch j{
						case 0:
						mareas[m] = append(mareas[m],[]float64{area2,4.0}...)
						mloads[m] = append(mloads[m], frcvec2...)
						mverts[m][i] = [][]float64{p1, bpt, tpt, p4}
						case 1:
						mareas[m] = append(mareas[m],[]float64{area2,4.0}...)
						mloads[m] = append(mloads[m], frcvec2...)
						mverts[m][i] = [][]float64{p2, p3, tpt, bpt}
						case 2:
						mareas[m] = append(mareas[m],[]float64{area1, 3.0}...)
						mloads[m] = append(mloads[m], frcvec1...)
						mverts[m][i] = [][]float64{p4, tpt, p3}
						case 3:
						mareas[m] = append(mareas[m],[]float64{area1, 3.0}...)
						mloads[m] = append(mloads[m], frcvec1...)
						mverts[m][i] = [][]float64{p1, p2, bpt}
					}	
				}			
			}
		}
	}
	return
}
*/

func (f *Frm3d) CalcBf3d() {
	//calc breadth of flange 
	colvec := make([][]float64, f.Nspans + 1)
	bmvec := make([][]float64, f.Nspans)
	byvec := make([][]float64, f.Nbays)
	//clvec := make([][]float64, 4)
	ctyps := make([]int, f.Nspans+1)
	btyps := make([]int, f.Nspans)
	bytyps := make([]int, f.Nbays)
	var cdx, bdx, bydx []int
	//sts = make([]int, nspans)
	if len(f.Csec) == 1{
		idx := f.Csec[0]
		cdx = make([]int, f.Nspans+1)
		for i := range cdx{
			cdx[i] = idx
		}
	}
	for i, idx := range cdx{
		colvec[i] = make([]float64, len(f.Sections[idx-1]))
		copy(colvec[i], f.Sections[idx-1])
		ctyps[i] = f.Styps[idx-1]
	}
	if len(f.Bsec) == 1{
		idx := f.Bsec[0]
		bdx = make([]int, f.Nspans)
		for i := range bdx{
			bdx[i] = idx
		}
	}
	if len(f.Bysec) == 1{
		idx := f.Bysec[0]
		bydx = make([]int, f.Nbays)
		for i := range bydx{
			bydx[i] = idx
		}
	}
	var bf, df, bw, dused float64
	for i, idx := range bdx{
		bdim := f.Sections[idx-1]
		var btyp int
		if f.Bstyp == 0 || f.Bstyp == 1{
			bst := f.Styps[idx-1]
			switch bst{
				case 6,7,8,9,10:
				btyp = bst
				default:
				btyp = 6
			}
		} else {
			btyp = f.Bstyp
		}
		switch btyp{
			case 6:
			f.Tyb = 1.0
			case 7,8,9,10:
			f.Tyb = 0.5
		}
		btyps[i] = btyp
		bf = GetBfRcc(f.Code, f.Bmrel, f.Tyb, f.Lspans[i],f.Dslb,bdim[0])
		df = f.Dslb; bw = bdim[0]; dused = bdim[1]
		bmvec[i] = []float64{bf,dused,bw,df}
	}
	
	for i, idx := range bydx{
		bdim := f.Sections[idx-1]
		var btyp int
		if f.Bystyp == 0 || f.Bystyp == 1{
			bst := f.Styps[idx-1]
			switch bst{
				case 6,7,8,9,10:
				btyp = bst
				default:
				btyp = 6
			}
		} else {
			btyp = f.Bystyp
		}
		if f.Bfycalc{
			switch btyp{
				case 6:
				f.Tyby = 1.0
				case 7,8,9,10:
				f.Tyby = 0.5
			}
			btyps[i] = btyp
			bf = GetBfRcc(f.Code, f.Byrel, f.Tyby, f.Lbays[i],f.Dslb,bdim[0])
			df = f.Dslb; bw = bdim[0]; dused = bdim[1]
			byvec[i] = []float64{bf,dused,bw,df}
		} else {
			byvec[i] = bdim
			bytyps[i] = btyp
		}
	}
	var secvec [][]float64
	var sts []int
	secvec = append(colvec,bmvec...)
	secvec = append(secvec,byvec...)
	sts = append(ctyps, btyps...)
	sts = append(sts, bytyps...)
	
	f.Csec, f.Bsec, f.Bysec = []int{}, []int{}, []int{}
	f.Sections = make([][]float64, len(secvec))
	f.Styps = make([]int, len(sts))
	for i := range secvec {
		f.Sections[i] = make([]float64, len(secvec[i]))
		copy(f.Sections[i], secvec[i])
		f.Styps[i] = sts[i]
		switch {
		case i < f.Nspans + 1:
			f.Csec = append(f.Csec, i+1)
		case i < f.Nspans + 1 + f.Nspans:
			f.Bsec = append(f.Bsec, i+1)
		default:
			f.Bysec = append(f.Bysec, i+1)
		}
	}
	return
}


func getmemtype(i, colstep, xstep, ystep, xdx, ydx, locdx, flrdx, mdx int) (mtyp, locx, lex int) {
	/*
	   locx - STARTS FROM 1; 1 to 9 (0 - 8)
	   strt - start index STARTS FROM 1 (1- col, 10 - beam x, 19 beamy (ends at 27))
	   lex - edge index (node count 1 - one beam (NOT HERE), 2 - corner, 3 - edge, 4 - interior)
	   classify by ex for symmetrical frame
	   classify by typ for unsym frame
	*/
	fmt.Println("mem-", i, xdx, ydx, mdx)
	var cl, cr, ct, cb, strt int
	switch mdx {
	case 1:
		//column
		lex = 4
		if xdx == 1 {
			cb = 1
		}
		if ydx == 1 {
			cl = 1
		}
		if xdx == ystep {
			ct = 1
		}
		if ydx == xstep {
			cr = 1
		}
		switch {
		case cb+cl == 2:
			//bot left corner
			locx = 1
			lex = 2
		case cb+cr == 2:
			//bot right corner col
			locx = 2
			lex = 2
		case ct+cl == 2:
			//top left corner
			locx = 3
			lex = 2
		case ct+cr == 2:
			//top right corner
			locx = 4
			lex = 2
		case cl == 1:
			//left edge col
			locx = 5
			lex = 3
		case cr == 1:
			//right edge col
			locx = 6
			lex = 3
		case ct == 1:
			//top edge col
			locx = 7
			lex = 3
		case cb == 1:
			//bottom edge col
			locx = 8
			lex = 3
		}
		strt = 1
	case 2:
		//beam x
		//lex == 1 is a cantilever span
		lex = 5
		switch {
		case xdx == 1 && ydx == 1:
			//bot left, edge beam edge frame
			locx = 1
			lex = 2
		case xdx == 1 && ydx == xstep-1:
			//bot right edge beam edge frame
			locx = 2
			lex = 2
		case xdx == ystep && ydx == 1:
			//top left edge beam edge frame
			locx = 3
			lex = 2
		case xdx == ystep && ydx == xstep-1:
			//top right edge beam edge frame
			locx = 4
			lex = 2
		case ydx == 1:
			//left edge beam edge beam int frame
			locx = 5
			lex = 4
		case ydx == xstep-1:
			//right edge beam edge beam int frame
			locx = 6
			lex = 4
		case xdx == ystep:
			//top edge beam int beam edge frame
			locx = 7
			lex = 3
		case xdx == 1:
			//bottom edge beam int beam edge frame
			locx = 8
			lex = 3
		}
		strt = 10
		//end = 18
	case 3:
		//beam y
		lex = 5
		switch {
		case xdx == 1 && ydx == 1:
			//bot left edge beam edge frame
			locx = 1
			lex = 2
		case xdx == 1 && ydx == xstep:
			//bot right edge beam edge frame
			locx = 2
			lex = 2
		case xdx == ystep-1 && ydx == 1:
			//top left edge beam edge frame
			locx = 3
			lex = 2
		case xdx == ystep-1 && ydx == xstep:
			//top right edge beam edge frame
			locx = 4
			lex = 2
		case ydx == 1:
			//left edge beam int beam edge frame
			locx = 5
			lex = 3
		case ydx == xstep:
			//right edge beam int beam edge frame
			locx = 6
			lex = 3
		case xdx == ystep-1:
			//top edge beam
			locx = 7
			lex = 4
		case xdx == 1:
			//bottom edge beam
			locx = 8
			lex = 4
		}
		strt = 19
		//end = 27
	}
	//fmt.Println("locx-",cl, cr, ct, cb, locx, strt)
	mtyp = locx + strt
	return
}

func getxydx(i, colstep, xstep, ystep int) (xdx, ydx, flrdx, locdx int){
	//edgex - returns 0 if not on an edge
	//1 - if on x edge
	//2 - if on y edge
	//3 - if on xy corner
	//4 - if on xy corner
	flrdx = i / colstep
	locdx = i % colstep
	if locdx == 0 {
		locdx = colstep
		flrdx -= 1
	}
	//if locdx == 0{flrdx = flrdx - 1}
	ydx = locdx % xstep
	xdx = (locdx-1)/xstep + 1
	if ydx == 0 {
		ydx = xstep
	}
	return
}

func getmemidx(jb, je int) (memidx string) {
	memidx = fmt.Sprintf("%v-%v", jb, je)
	return
}

func framegen(fck, X, Y, Z float64, xs, ys, zs []float64, xspans, yspans, zspans int) ([][]float64, map[int][]float64, map[int][]int, [][]int, [][]int, [][]int, [][]int) {
	//calc rcc E and G
	//density, poi_mu, E, G := rcfmaterialcalc(fck)
	//rcc_em := []float64{E,G,density,poi_mu}
	/*
		chn2d := make(chan []interface{},1)
		mod2d := make([]interface{},6)
		mod2d[0] = []float64{X,Y,Z}
		mod2d[1] = []int{xspans,yspans,zspans}
		mod2d[2] = xs
		mod2d[3] = ys
		mod2d[4] = zs
		mod2d[5] = rcc_em
		go rc2df(mod2d,chn2d)
	*/
	//cp := [][]float64{}
	coords := [][]float64{}
	for _, z := range zs {
		for _, y := range ys {
			for _, x := range xs {
				coords = append(coords, []float64{x, y, z})
			}
		}
	}
	colstep := (xspans + 1) * (yspans + 1)
	xstep := xspans + 1
	//ystep := yspans+1
	nodecords := make(map[int][]float64)
	for idx, vertex := range coords {
		nodecords[idx+1] = vertex
	}
	cols := [][]int{}
	supports := make(map[int][]int)
	members := make(map[string][]interface{})
	beamys := [][]int{}
	beamxs := [][]int{}
	slabnodes := [][]int{}
	slabmems := [][]string{}
	nodeadj := make(map[int][]int)
	var memidx string
	//members = start node jb end node je memtyp (0 - col, 1 - beamx, 2 - beamy)
	//members = memloc x 0 edge 1 interior memloc y 0 edge 1 interior
	//ncolgroups =
	for i := 1; i <= len(coords); i++ {
		if i+colstep <= len(coords) {
			cols = append(cols, []int{i, i + colstep})
			nodeadj[i] = append(nodeadj[i], i+colstep)
			nodeadj[i+colstep] = append(nodeadj[i+colstep], i)
			//[0] - typ (col -0, beamx - 1, beamy - -1)
			//memidx = i+((i+colstep)*(i+colstep))
			memidx = getmemidx(i, i+colstep)
			members[memidx] = append(members[memidx], []int{i, i + colstep, 1, 1, 0})
		}
		if i <= colstep {
			supports[i] = []int{-1, -1, -1, -1, -1, -1}
		} else {
			supports[i] = []int{0, 0, 0, 0, 0, 0}
		}
		if i > colstep {
			if i+1 <= len(coords) && nodecords[i][1] == nodecords[i+1][1] {
				beamxs = append(beamxs, []int{i, i + 1})
				nodeadj[i] = append(nodeadj[i], i+1)
				nodeadj[i+1] = append(nodeadj[i+1], i)
				memidx = getmemidx(i, i+1)
				members[memidx] = append(members[memidx], []int{i, i + 1, 1, 1, 1})
			}
			if i+xstep <= len(coords) && nodecords[i][2] == nodecords[i+xstep][2] {
				beamys = append(beamys, []int{i, i + xstep})
				nodeadj[i] = append(nodeadj[i], i+xstep)
				nodeadj[i+xstep] = append(nodeadj[i+xstep], i)
				memidx = getmemidx(i, i+xstep)
				members[memidx] = append(members[memidx], []int{i, i + xstep, 1, 1, 1})

			}
		}
		if i+xstep+1 <= len(coords) && nodecords[i][2] == nodecords[i+xstep+1][2] {
			ljb := i
			lje := i + xstep
			lidx := getmemidx(ljb, lje)
			rjb := i + 1
			rje := i + xstep + 1
			ridx := getmemidx(rjb, rje)
			tjb := i + xstep
			tje := i + xstep + 1
			tidx := getmemidx(tjb, tje)
			bjb := i
			bje := i + 1
			bidx := getmemidx(bjb, bje)
			slabnodes = append(slabnodes, []int{i, i + 1, i + xstep + 1, i + xstep})
			slabmems = append(slabmems, []string{lidx, ridx, tidx, bidx})
		}
	}

	//slab calc pipeline
	//go slabcalc(slabnodes,slabmems)
	//calc slab loads on beams
	return coords, nodecords, supports, cols, beamxs, beamys, slabnodes
}

func rectsectioncalc(b, d float64) []float64 {
	area := b * d
	iz := b * math.Pow(d, 3) / 12.0
	return []float64{area, iz}
}

//YE OLDE
/*
func (f *F3d) Printz() (rez string) {
	fjson, err := json.Marshal(f)
	if err != nil {
		log.Println(err.Error())
	}
	rez = string(fjson)
	return
}

func rcc2dframegen(xs, ys, ws, ps []float64, fck_col, fck_bm, col_b, col_d, bm_b, bm_d float64) (coords [][]float64, nodes map[int][]float64, members, supports [][]int, jp, msp, em, cp [][]float64) {

	nodeidx := 1
	nodes = make(map[int][]float64)
	ymap := make(map[float64]int)
	xmap := make(map[float64]int)
	for yidx, y := range ys {
		if _, ok := ymap[y]; !ok {
			ymap[y] = yidx + 1
		}
		for xidx, x := range xs {
			if _, ok := xmap[x]; !ok {
				xmap[x] = xidx + 1
			}
			coords = append(coords, []float64{x, y})
			nodes[nodeidx] = []float64{x, y}
			nodeidx += 1
		}
	}
	memidx := 1
	xstep := len(xs)
	ystep := len(ys)
	em = make([][]float64, xstep+ystep-1)
	cp = make([][]float64, xstep+ystep-1)
	for i := 1; i <= xstep+ystep-1; i++ {
		if i <= xstep {
			em[i-1] = []float64{fck_col}
			cp[i-1] = rectsectioncalc(col_b, col_d)
		} else {
			em[i-1] = []float64{fck_bm}
			cp[i-1] = rectsectioncalc(bm_b, bm_d)
		}
	}
	//colgrpidx := 1
	//build members
	//CURRENT mrel = 0, add for type of structure (rigid frame vs hinged, shouldn't matter for a column?)
	//mprp 2d = [jb je em cp mrel colidx flridx]
	//em for col/beam, cp - as per colgroups (colidx) for column and flrgroups (flridx) for beam
	//ALL columns by x one group, all beams by floor one group (for now)
	for i, pt := range coords {
		ptx := pt[0]
		pty := pt[1]
		colidx := xmap[ptx]
		flridx := ymap[pty]
		idx := i + 1
		if idx+xstep <= len(coords) {
			//check for column index from xmap and assign cp
			//colidx = sort.SearchFloat64s(slice_1, f1)
			members = append(members, []int{idx, idx + xstep, colidx, colidx, 0, colidx, flridx})
			memidx++
		}
		if idx <= xstep {
			supports = append(supports, []int{idx, -1, -1, -1})
		} else {
			supports = append(supports, []int{idx, 0, 0, 0})
		}
		if idx >= xstep && idx < len(coords) && nodes[idx+1][1] == nodes[idx][1] {
			//beams = append(beams, []int{idx, idx + 1})
			members = append(members, []int{idx, idx + 1, xstep + flridx - 1, xstep + flridx - 1, 0, colidx, flridx})

			if len(ws) == 1 {
				msp = append(msp, []float64{float64(memidx), 3.0, ws[0], 0, 0, 0})
			} else {
				msp = append(msp, []float64{float64(memidx), 3.0, ws[flridx-1], 0, 0, 0})
			}

			memidx++
		}
	}
	//apply nodal loads (P)
	for i := 1 + xstep; i < len(nodes); i += xstep {
		idx := (i - xstep) / xstep
		if len(ps) == 1 {
			jp = append(jp, []float64{float64(i), ps[0], 0, 0})
		} else {
			jp = append(jp, []float64{float64(i), ps[idx], 0, 0})
		}
	}
	return
}

func parsebaystr(baystr string) (bx []float64, mbx []int) {
	//bx := []float64{}
	//s := strings.Split(baystr, ";")
	bxstr := strings.Split(baystr, "-")[0]
	mbxstr := strings.Split(baystr, "-")[1]
	for _, b := range bxstr {
		val, _ := strconv.ParseFloat(string(b), 64)
		bx = append(bx, val)
	}
	for _, mb := range mbxstr {
		val, _ := strconv.Atoi(string(mb))
		mbx = append(mbx, val)
	}
	return bx, mbx
}

func rcfinput(cmdz string) (float64, float64, float64, [][]float64, map[int][]float64, map[int][]int, [][]int, [][]int, [][]int, [][]int) {
	xyz := strings.Split(cmdz, ",")
	mod := xyz[0]
	X, _ := strconv.ParseFloat(xyz[1], 64)
	Y, _ := strconv.ParseFloat(xyz[2], 64)
	Z, _ := strconv.ParseFloat(xyz[3], 64)
	fck, _ := strconv.ParseFloat(xyz[7], 64)

	var xspans, yspans, zspans int
	var coords [][]float64
	var nodecords map[int][]float64
	var supports map[int][]int
	var cols, beamxs, beamys, slabnodes [][]int
	switch mod {
	case "0":
		xspans, _ = strconv.Atoi(xyz[4])
		yspans, _ = strconv.Atoi(xyz[5])
		zspans, _ = strconv.Atoi(xyz[6])
		xs := spangen(X, xspans)
		ys := spangen(Y, yspans)
		zs := spangen(Z, zspans)
		coords, nodecords, supports, cols, beamxs, beamys, slabnodes = framegen(fck, X, Y, Z, xs, ys, zs, xspans, yspans, zspans)
	case "1":
		bx, mbx := baycmdz(xyz[4])
		by, mby := baycmdz(xyz[5])
		bz, mbz := baycmdz(xyz[6])
		xs, xspans := baygen(X, bx, mbx)
		ys, yspans := baygen(Y, by, mby)
		zs, zspans := baygen(Z, bz, mbz)
		coords, nodecords, supports, cols, beamxs, beamys, slabnodes = framegen(fck, X, Y, Z, xs, ys, zs, xspans, yspans, zspans)
	}
	return X, Y, Z, coords, nodecords, supports, cols, beamxs, beamys, slabnodes
}

func varbaygen(bx []float64) []float64 {
	var xs []float64
	var baysum float64

	for _, x := range bx {
		xs = append(xs, baysum)
		baysum += x
	}
	xs = append(xs, xs[len(xs)-1]+bx[len(bx)-1])
	return xs
}

func baycmdz(baystr string) ([]float64, []int) {
	bx := []float64{}
	mbx := []int{}
	s := strings.Split(baystr, ";")
	bxstr := strings.Split(s[0], "-")
	mbxstr := strings.Split(s[0], "-")
	for _, b := range bxstr {
		val, _ := strconv.ParseFloat(b, 64)
		bx = append(bx, val)
	}
	for _, mb := range mbxstr {
		val, _ := strconv.Atoi(mb)
		mbx = append(mbx, val)
	}
	return bx, mbx
}

func spangen(xl float64, xspans int) []float64 {
	xs := make([]float64, xspans+1)
	xs = floats.Span(xs, 0, xl)
	return xs
}

func baygen(xl float64, bx []float64, mbx []int) ([]float64, int) {
	var xs []float64
	var nb, bsum, bspan, xmax float64
	var nspan int
	switch mbx[0] {
	case 0:
		for _, bn := range bx {
			bspan += bn
			nspan++
		}
		if math.Mod(xl, bspan) > 1e-10 {
			log.Println("Non integer bays")
			return xs, nspan
		}
		nb = xl / bspan
		for x := 0; x < int(nb); x++ {
			bsum = 0
			for _, bn := range bx {
				xs = append(xs, bspan*float64(x)+bsum)
				if xmax < bn {
					xmax = bn
				}
				bsum += bn
			}
		}
		xs = append(xs, xs[len(xs)-1]+bx[len(bx)-1])
	case 1:
		xs = []float64{}
	}
	return xs, int(nb)

}

func rcfmaterialcalc(fck float64) (density, poi_mu, E, G float64) {
	//g = e/(2(1+mu)), e = 5000 sqrt(fck)
	E = 5000.0 * math.Sqrt(fck)
	density = 2450
	poi_mu = 1.6
	G = E / (2.0 * (1.0 + poi_mu))
	return
}
	//Nbms, Ncols  int `json:",omitempty"`
	//Memwng    [][]float64 `json:",omitempty"`
	//EL, CL [][]float64 `json:",omitempty"`
	//Coords    [][]float64 `json:",omitempty"`
	//Cdmap     map[int][]float64 `json:",omitempty"`
	//Nodecords map[int][]float64 `json:",omitempty"`
	//Supports  map[int][]int `json:",omitempty"`
	//Nodeadj   map[int][]int `json:",omitempty"`

*/
