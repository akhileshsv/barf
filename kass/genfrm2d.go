package barf

import (
	"os"
	"io/ioutil"
	"log"
	"fmt"
	"math"
	"math/rand"
	"errors"
	"runtime"
	"encoding/json"
	"path/filepath"
)

//Frm2d is a 2d frame generation struct
//kinda broken 
type Frm2d struct {
	//it's a me Frm2d
	Id           int
	Nspans       int
	Nflrs        int
	Code         int
	X,Y,Z        []float64
	DL           float64
	LL           float64
	WL           []float64
	Sections     [][]float64
	Styps        []int 
	Csec, Bsec   []int
	Cstyp,Bstyp  int
	Clstyp       int
	Fcks         []float64
	Fys          []float64
	Fyds         []float64
	Lspans       []float64
	Lbays        []float64
	Lflrs        []float64
	DM           float64
	PSFs         []float64
	WLFs         []float64
	Clvrs        [][]float64
	Cldim        [][]float64
	Clsec        []int
	Term         string
	Title        string
	Foldr        string
	Verbose      bool
	Spam         bool
	Bfcalc       bool //flange width calculation
	Fltslb       bool
	Lclvr, Rclvr bool
	Rekt         bool //rect col/prism. beam
	Plinth       bool
	Braced       bool
	Fixbase      bool
	Tweak        bool
	Nosec        bool
	Web          bool
	Noprnt       bool
	Dconst       bool `json:",omitempty"`
	Nloads       int
	Ldcalc       int
	Wlcalc       int
	Wlvec        []float64
	Selfwt       int
	Mtyp         int
	Stair        int
	Ngrp         int //ISN'T THIS JUST FGEN and a rose by nother etc
	Fgen         int
	Bmrel        int
	Exp          int //exposure condition
	Slbload      int
	Maxstep      int
	Ceo          float64 //col edge offset ()
	Nomcvr       float64
	D1, D2       float64
	Lx, Ly       float64
	Hc, Dh       float64 //hc - effective diam of col head, dh- depth of head
	Drx          float64 //drop x dim
	Dry          float64 //drop y dim
	Drd          float64 //drop depth
	Dslb         float64
	Fyv          float64
	Lbay         float64
	Pg           float64 `json:",omitempty"`
	Matprop      []float64 `json:",omitempty"`
	Lcmem, Rcmem []int `json:",omitempty"`
	Ljs, Rjs     []int `json:",omitempty"`
	Nlp          int `json:",omitempty"`
	Cdim, Bdim   [][]float64 `json:",omitempty"`
	Tyb          float64 `json:",omitempty"`
	Secmap       map[int][]float64 `json:",omitempty"`
	Mod          Model `json:",omitempty"`
	Vec          [][]int `json:",omitempty"`
	Val          [][]float64 `json:",omitempty"`
	Flridx       []int `json:",omitempty"`
	Txtplots     []string `json:",omitempty"`
	DMs          []float64 `json:",omitempty"`
	Bloads       [][]float64 `json:",omitempty"`
	Cloads       [][]float64 `json:",omitempty"`
	Jloads       [][]float64 `json:",omitempty"`
	Wloadl       [][]float64 `json:",omitempty"`
	Wloadr       [][]float64 `json:",omitempty"`
	Nodemap      map[int][]int `json:",omitempty"`
	Nodes        map[Pt][]int `json:"-"`
	Members      map[int][][]int `json:",omitempty"`
	Cols, Beams  []int `json:",omitempty"`
	Mloadmap     map[int][][]float64 `json:",omitempty"`
	Jloadmap     map[int][][]float64 `json:",omitempty"`
	Advloads     map[int][][]float64 `json:",omitempty"`
	Benloads     map[int][][]float64 `json:",omitempty"`
	Wadloads     map[int][][]float64 `json:",omitempty"`
	Uniloads     map[int][][]float64 `json:",omitempty"`
	Loadcons     map[int][][]float64 `json:",omitempty"`
	Bmenv        map[int]*BmEnv `json:",omitempty"`
	Colenv       map[int]*ColEnv `json:",omitempty"`
	Secprop      []Secprop `json:",omitempty"`
	Ms           map[int]*Mem `json:",omitempty"`
	Mslmap       map[int]map[int][][]float64 `json:",omitempty"`
	CBeams       [][]int `json:",omitempty"`
	Reports      []string `json:",omitempty"`
	Nbms, Ncols  int `json:",omitempty"`
	Opt          int `json:",omitempty"`
	Fop          int `json:",omitempty"` //if 1, optimize
	Width        float64 `json:",omitempty"`
	Cvec         [][]float64 `json:",omitempty"`
	Bvec         [][]float64 `json:",omitempty"`
	Cdims        [][]float64 `json:",omitempty"`
	Bdims        [][]float64 `json:",omitempty"`
	Ncls         int `json:",omitempty"`
	Brel         int `json:",omitempty"`
	Df           float64 `json:",omitempty"`
	Cmap         map[int][]int `json:",omitempty"`
	Dz           bool `json:",omitempty"`
	Report       string `json:",omitempty"`
	Rez          []interface{} `json:",omitempty"`
	Kostin       []float64 `json:",omitempty"`
	Quants       []float64 `json:",omitempty"`
	Kosts        []float64 `json:",omitempty"`
	Kost         float64 `json:",omitempty"`
}

//Frm2dRez is worthless and will be deleted
type Frm2dRez struct {
	//maybe this is not needed DELETE
	Bmenvs   []BmEnv
	Colenvs  []ColEnv
	Bmvec    []int
	Colvec   []int
	Flridx   int
	Txtplots []string
}


//Init initializes default frame 2d values
func (f *Frm2d) Init() (err error){
	if f.Title == ""{
		f.Title = fmt.Sprintf("frame2d_%v",rand.Intn(666))
	}
	
	if f.Ngrp == 0{
		//single column and beam section (nsections = 2)
		f.Ngrp = 1
	}
	if (len(f.Sections) == 0 || len(f.Csec) == 0 || len(f.Bsec) == 0) && f.Opt ==0{
		return errors.New("no sections specified")
	}
	if len(f.X) == 0 && len(f.Lspans) == 0{
		return errors.New("frame spans (x) not specified")
	}
	if len(f.Y) == 0 && len(f.Lflrs) == 0{
		return errors.New("frame floor heights (y) not specified")
	}
	if f.Mtyp == 0{
		f.Mtyp = 1
	}
	//is code, bs code, (lol) euro code 
	if f.Code == 0 {f.Code = 1}
	if f.PSFs == nil{f.PSFs = []float64{1.5,1.0,1.5,0.0}}
	if f.WLFs == nil{f.WLFs = []float64{1.2,1.2,1.2}}
	f.Nodemap = make(map[int][]int)
	f.Nodes = make(map[Pt][]int)
	f.Members = make(map[int][][]int)
	f.Advloads = make(map[int][][]float64)
	f.Benloads = make(map[int][][]float64)
	f.Loadcons = make(map[int][][]float64)
	f.Uniloads = make(map[int][][]float64)
	f.Jloadmap = make(map[int][][]float64)
	f.Wadloads = make(map[int][][]float64)
	if len(f.X) == 0 || f.X == nil{
		f.X = genxs(f.Lspans)
	}
	if len(f.Y) == 0 || f.Y == nil{
		f.Y = genxs(f.Lflrs)
	}
	if len(f.X) == 1 && f.Nspans != 0{
		var x, xstep float64
		xstep = f.X[0]//not your average "xstep"
		if xstep == 0{return errors.New("frame spans (x) not specified")}
		f.X = []float64{x}
		for i := 0; i < f.Nspans; i++{
			x += xstep
			f.X = append(f.X, x)
		}
	}
	if len(f.Y) == 1 && f.Nflrs != 0{
		var y, ystep float64
		ystep = f.Y[0]
		if ystep == 0{return errors.New("frame floors (y) not specified")}
		f.Y = []float64{y}
		for i := 0; i < f.Nflrs; i++{
			y += ystep
			f.Y = append(f.Y, y)
		}

	}
	if len(f.Clvrs) > 0{
		if f.Clvrs[0][0] > 0.0 {f.Ncls++}
		if f.Clvrs[1][0] > 0.0 {f.Ncls++}
	}
	if f.Nflrs == 0{f.Nflrs = len(f.Y)-1}
	if f.Nspans == 0{f.Nspans = len(f.X)-1}

	if f.Styps == nil || len(f.Styps) == 0{
		f.Styps = make([]int, len(f.Sections))
		for i := range f.Sections{
			f.Styps[i] = 1
		}
	}
	for i := 1; i < len(f.X); i++{
		f.Lspans = append(f.Lspans, f.X[i]-f.X[i-1])
	}
	//define material prop
	if f.Mtyp > 0{
		f.InitMat()
	}
	if f.Mtyp == 1 && f.Fyv == 0{
		f.Fyv = f.Fys[2]
	}
	if f.Bfcalc && !f.Nosec{
		//fmt.Println("bf calcing as we speaks")
		var mrel, bstyp int
		if f.Dslb == 0.0{
			//fmt.Println(ColorRed,"setting flange depth to 125.0",ColorReset)
			f.Dslb = 150.0
		}
		//if f.Nspans == 1 && (len(f.Clvrs) == 0 || (f.Clvrs[0][0] + f.Clvrs[1][0] == 0.0)){mrel = 1}
		secvec, sts := CalcBf(f.Code, bstyp, mrel, f.Nspans, f.Dslb , f.Csec, f.Bsec, f.Styps, f.Lspans, f.Sections)
		f.Csec, f.Bsec = []int{}, []int{}
		f.Sections = make([][]float64, len(secvec))
		f.Styps = make([]int, len(sts))
		for i := range secvec{
			f.Sections[i] = make([]float64, len(secvec[i]))
			copy(f.Sections[i], secvec[i])
			f.Styps[i] = sts[i]
			if i < len(f.X){
				f.Csec = append(f.Csec, i+1)
			} else {
				f.Bsec = append(f.Bsec, i+1)
			}
			
		}
	}
	return
}

//Getndim returns ndims for 2d frame opt
func (f *Frm2d) Getndim()(nd int){
	if f.Ngrp > 3{f.Ngrp = 3}
	switch f.Ngrp{
		case 1:
		//one col, one beam
		switch f.Width{
			case 0.0:
			//bw, dc, db
			nd = 3
			default:
			//dc, dbm
			nd = 2
		}
		case 2:
		//edge col, mid col, edge beam , mid beam
		switch f.Width{
			case 0.0:
			//bw, dec, dmc, deb, dmb
			nd = 5
			default:
			//dec, dmc, deb, dmb
			switch f.Dconst{
				case true:
				nd = 3
				case false:
				nd = 4
			}
		}
		case 3:
		//ncols = len(f.X); nbms := len()
		switch f.Width{
			case 0.0:
			//log.Println(ColorRed,"ERRORE,errore-setting width to 225.0 mm",ColorReset)
			f.Width = 225.0
			switch f.Dconst{
				case true:
				nd = len(f.X) + 1
				case false:
				nd = 2 * len(f.X) - 1
			}
			default:
			switch f.Dconst{
				case true:
				nd = len(f.X) + 1
				case false:
				nd = 2 * len(f.X) - 1
			}
		}
	}
	return
}

//SecGenRcc generates rcc sections from a pos vec (for opt)
func (f *Frm2d) SecGenRcc(pos []float64)(err error){
	var bw, bc, dc, db, dec, dmc, deb, dmb float64
	step := 25.0
	strt := 225.0
	//HERE
	if f.Width == 0.0{
		f.Width = 225
		// if f.Ceo == 0.0{
		// 	f.Ceo = 25.0
		// }
	}
	switch f.Ngrp{
		case 1:
		//one col, one beam
		switch f.Width{
			case 0.0:
			//WILL NEBER HAPPENZ
			if len(pos) < 3{
				err = fmt.Errorf("invalid posvec length %f for ngrp %v",pos,f.Ngrp)
				return
			}
			bw = math.Round(pos[0]) * step + strt
			dc = math.Round(pos[1]) * step + strt
			db = math.Round(pos[2]) * step + strt
			default:
			if len(pos) < 2{
				err = fmt.Errorf("invalid posvec length %f for ngrp %v",pos,f.Ngrp)
				return
			}
			bw = f.Width
			dc = math.Round(pos[0]) * step + strt
			db = math.Round(pos[1]) * step + strt			
		}
		bc = bw + f.Ceo * 2.0
		f.Sections = make([][]float64, 2)
		f.Styps = make([]int, 2)
		f.Sections[0] = []float64{bc, dc}
		f.Sections[1] = []float64{bw, db}
		f.Csec = []int{1}
		f.Bsec = []int{2}
		f.Styps[0] = f.Cstyp
		f.Styps[1] = f.Bstyp
		case 2:
		//edge col, mid col, edge beam, mid beam?
		switch f.Width{
			case 0.0:
			default:
			bw = f.Width
			switch f.Dconst{
				case true:
				if len(pos) < 3{
					err = fmt.Errorf("invalid posvec length %f for ngrp %v",pos,f.Ngrp)
					return
				}
				dec = math.Round(pos[0]) * step + strt
				dmc = math.Round(pos[1]) * step + strt			
				deb = math.Round(pos[2]) * step + strt
				bc = f.Ceo * 2.0 + bw
				f.Sections = make([][]float64, 3)
				f.Sections[0] = []float64{bc, dec}
				f.Sections[1] = []float64{bc, dmc}
				f.Sections[2] = []float64{bw, deb}
				f.Styps = []int{f.Cstyp, f.Cstyp, f.Bstyp}
				for i := range f.X{
					switch i{
						case 0:
						f.Csec = append(f.Csec, 1)
						f.Bsec = append(f.Bsec, 3)
						case len(f.X)-1:
						f.Csec = append(f.Csec, 1)
						default:
						f.Csec = append(f.Csec, 2)
						f.Bsec = append(f.Bsec, 3)
					}
				}
				case false:
				if len(pos) < 4{
					err = fmt.Errorf("invalid posvec length %f for ngrp %v",pos,f.Ngrp)
					return
				}
				dec = math.Round(pos[0]) * step + strt
				dmc = math.Round(pos[1]) * step + strt			
				deb = math.Round(pos[2]) * step + strt
				dmb = math.Round(pos[3]) * step + strt
				bc = f.Ceo * 2.0 + bw
				f.Sections = make([][]float64, 4)
				f.Sections[0] = []float64{bc, dec}
				f.Sections[1] = []float64{bc, dmc}
				f.Sections[2] = []float64{bw, deb}
				f.Sections[3] = []float64{bw, dmb}
				f.Styps = []int{f.Cstyp, f.Cstyp, f.Bstyp, f.Bstyp}
				for i := range f.X{
					switch i{
						case 0:
						f.Csec = append(f.Csec, 1)
						f.Bsec = append(f.Bsec, 3)
						case len(f.X)-1:
						f.Csec = append(f.Csec, 1)
						default:
						f.Csec = append(f.Csec, 2)
						f.Bsec = append(f.Bsec, 4)
					}
				}
			}
		}
		case 3:
		//all individual secshuns 
		switch f.Width{
			case 0.0:
			//BAAAAAH
			default:
			bw = f.Width
			ncols := len(f.X)
			nbms := ncols -1
			switch f.Dconst{
				case true:
				//ncols + bd
				if len(pos) < ncols + 1{
					err = fmt.Errorf("invalid posvec length %f for ngrp %v",pos,f.Ngrp)
					return
				}
				for i := 0; i < ncols; i++{
					bc := bw + 2.0 * f.Ceo
					dc := math.Round(pos[i]) * step + strt
					cpvec := []float64{bc, dc}
					f.Sections = append(f.Sections, cpvec)
					f.Styps = append(f.Styps, f.Cstyp)
					f.Csec = append(f.Csec, i+1)
				}
				dbm := math.Round(pos[ncols]) * step + strt
				f.Sections = append(f.Sections, []float64{bw, dbm})
				f.Styps = append(f.Styps, f.Bstyp)
				for i := 0; i < nbms; i++{
					f.Bsec = append(f.Bsec, ncols+1)
				}
				case false:
				if len(pos) < ncols + nbms{
					err = fmt.Errorf("invalid posvec length %f for ngrp %v",pos,f.Ngrp)
					return
				}
				for i := 0; i < ncols; i++{
					bc := bw + 2.0 * f.Ceo
					dc := math.Round(pos[i]) * step + strt
					cpvec := []float64{bc, dc}
					f.Sections = append(f.Sections, cpvec)
					f.Styps = append(f.Styps, f.Cstyp)
					f.Csec = append(f.Csec, i+1)
				}
				for i := 0; i < nbms; i++{
					dbm := math.Round(pos[ncols+i]) * step + strt
					cpvec := []float64{bw, dbm}
					f.Sections = append(f.Sections, cpvec)
					f.Styps = append(f.Styps, f.Cstyp)
					f.Bsec = append(f.Bsec, ncols+i+1)
				}
			}
		}
	}
	return
}

//DrawMod plots the model in a Frm2d struct
//(todo) - make this plot load cases
func (f *Frm2d) DrawMod(term string)(txtplot string){
	//sort msloads for load type and member, sum up
	//OR - do dead and live load plots?
	//ldmap := make(map[int][]float64)
	plotchn := make(chan string)
	go PlotFrm2d(&f.Mod, term, plotchn)
	txtplot = <- plotchn
	//f.Txtplots = append(f.Txtplots, txtplot)
	//for _, flrvec := range f.CBeams{
	//	f.Txtplots = append(f.Txtplots, PlotBmEnv(f.Bmenv, flrvec, term))}
	//
	return
}


//DrawLp plots an individual load pattern in a Frm2d struct
func (f *Frm2d) DrawLp(lp int, term string)(txtplot string){
	f.Mod.Msloads = f.Loadcons[lp]
	if val, ok := f.Jloadmap[lp]; ok{
		f.Mod.Jloads = val
	} else {f.Mod.Jloads = [][]float64{}}
	//jl, ml := f.Mod.SumFrcs()
	//f.Mod.Msloads = ml; f.Mod.Jloads = jl
	f.Mod.Id = fmt.Sprintf("%s_lp_%v",f.Title,lp)
	f.Mod.Foldr = f.Foldr
	f.Mod.Web = f.Web
	txtplot = f.DrawMod(term)
	return 
}


//Drawflps is the goroutine func for load pattern plotting from design funcs
func Drawflps(f *Frm2d, lpchn chan []string){
	var rez []string
	for lp, lds := range f.Loadcons{
		if len(lds) == 0{continue}
		pltstr := f.DrawLp(lp, f.Term)
		rez = append(rez, pltstr)
	}
	lpchn <- rez
}


//Calc is the entry func for frm2d generation and calcs/analysis
func (f *Frm2d) Calc() (err error) {
	//gets bm and col envelopes
	err = f.Init()
	if err !=nil{
		return
	}
	f.GenCoords()
	f.GenMprp()
	f.GenLoads()
	f.InitMemRez()
	f.CalcLoadEnv()
	//f.CalcServLoads()
	//DO THIS IN MOSH
	//f.Mrd()
	//if plotlvl  = 0; 1; 2 etc
	//if f.Term != "" && f.Spam{
	//	for lp := range f.Loadcons{
	//		f.DrawLp(lp, f.Term)
	//	}
	//}
	//for _, txtplt := range f.Txtplots{
	//	fmt.Println(txtplt)
	//}
	return 
}

//AddMem adds a member to model
//mtyps - col, bm, lclvr, rclvr (1,2,3,4)
func (f *Frm2d) AddMem(xstep, jb, je, em, cp, mrel, mtyp int) (err error){
	xdx, fdx, locx, lex, mtx := getmemvec(f.Fgen,mtyp,jb,je,xstep)
	if jb > len(f.Mod.Coords) || je > len(f.Mod.Coords){
		return errors.New("invalid member node")
	}
	mdx := len(f.Mod.Mprp) + 1
	f.Mod.Mprp = append(f.Mod.Mprp,[]int{jb, je, em, cp, mrel})
	f.Secmap[mdx] = f.Sections[cp-1]
	//f.Secmap[mdx] = append(f.Secmap[mdx],float64(f.Styps[cp-1]))
	f.Nodemap[jb] = append(f.Nodemap[jb],mdx)
	f.Nodemap[je] = append(f.Nodemap[je],mdx)
	f.Members[mdx] = append(f.Members[mdx],[]int{jb, je, em, cp, mrel})
	f.Members[mdx] = append(f.Members[mdx],[]int{mtyp, xdx, fdx, locx, lex, mtx})
	switch mtyp{
		case 1:
		f.Cols = append(f.Cols, mdx)
		case 2:
		f.Beams = append(f.Beams, mdx)
		case 3:
		f.Lcmem = append(f.Lcmem, mdx)
		f.Beams = append(f.Beams, mdx)
		case 4:
		f.Rcmem = append(f.Rcmem, mdx)
		f.Beams = append(f.Beams, mdx)
	}
	//fmt.Println("added mem->",mdx,f.Members[mdx])
	return
}

//GenCoords generates (base) frame nodes
func (f *Frm2d) GenCoords(){
	for _, y := range f.Y {
		for _, x := range f.X {
			_ = f.AddNode(x,y)
		}
	}
	return
}

//GenMprp generates mprp and members
func (f *Frm2d) GenMprp() (err error){
	xstep := len(f.X); nflrs := len(f.Y)-1
	var csec, bsec, cedx, bedx int
	//switch f.Mtyp
	if len(f.Fcks) == 1 {
		cedx = 1; bedx = 1
		f.Mod.Em = append(f.Mod.Em, []float64{FckEm(f.Fcks[0])})
	} else {
		cedx = 1 ; bedx = 2
		f.Mod.Em = append(f.Mod.Em, []float64{FckEm(f.Fcks[0])})
		f.Mod.Em = append(f.Mod.Em, []float64{FckEm(f.Fcks[1])})
	}
	f.Secmap = make(map[int][]float64)
	var sup []int
	switch f.Fixbase{//footing vec
		case true:
		//fixed footing
		sup = []int{-1,-1,-1}
		case false:
		//pinned footing
		sup = []int{-1,-1,0}
	}
	//add supports and columns
	for i:=1; i <=len(f.Mod.Coords); i++{	
		if i <= xstep{
			//add supports
			supi := append([]int{i},sup...)
			f.Mod.Supports = append(f.Mod.Supports, supi)
		}
		if i + xstep <= len(f.Mod.Coords){
			//add columns
			xdx := (i-1)%xstep
			if len(f.Csec) == 1{csec = f.Csec[0]} else {csec = f.Csec[xdx]}
			_ = f.AddMem(xstep, i, i+xstep, cedx, csec, 0, 1)
			if i > xstep{
				switch xdx{
					case 0:
					f.Ljs = append(f.Ljs, i)
					case xstep - 1:
					f.Rjs = append(f.Rjs, i)
				}
			}
		}
	}
	//add beams
	f.Ljs = append(f.Ljs, xstep * nflrs + 1)
	f.Rjs = append(f.Rjs, xstep * (nflrs +1))
	for i := 1+xstep; i <= len(f.Mod.Coords); i++{
		if i + 1 <= len(f.Mod.Coords) && f.Mod.Coords[i-1][1] == f.Mod.Coords[i][1]{
			//add beams
			var brel int
			switch f.Bmrel{
				case 0,1:
				//regular cbeam
				brel = 0
				case 2:
				//both ends hinged (simply supported)
				brel = 3
			}
			xdx := (i-1)%xstep
			if len(f.Bsec) == 1{bsec = f.Bsec[0]} else {bsec = f.Bsec[xdx]}
			_ = f.AddMem(xstep, i, i+1, bedx, bsec, brel, 2)
		}		
	}
	//add left clvr
	if len(f.Clvrs) >= 1 && f.Clvrs[0][0] > 0.0{
		f.Lclvr = true
		clen := f.Clvrs[0][0]
		clsec := f.Bsec[0]
		if len(f.Clsec) > 0 {clsec = f.Clsec[0]}
		strt := 1
		if f.Plinth{
			strt = 2
		}
		for i, y := range f.Y[strt:]{
			_ = f.AddNode(-clen, y)
			jb := (strt+i) * xstep + 1; je := len(f.Mod.Coords)
			_ = f.AddMem(xstep, jb, je, bedx, clsec, 0, 3)
		}
	}
	//right clvr
	if len(f.Clvrs) > 1 && f.Clvrs[1][0]  > 0{
		//right cantilever on
		f.Rclvr = true
		clen := f.Clvrs[1][0]
		clsec := f.Bsec[0]
		switch len(f.Clsec){
			case 1:
			clsec = f.Clsec[0]
			case 2:
			clsec = f.Clsec[1]
			
		}
		xb := f.X[len(f.X)-1]+clen
		strt := 1
		if f.Plinth{
			strt = 2
		}
		for i, y := range f.Y[strt:]{
			_ = f.AddNode(xb, y)
			jb := (strt + 1 + i) * xstep; je := len(f.Mod.Coords)	
			_ = f.AddMem(xstep, jb, je, bedx, clsec, 0, 4)
		}
	}
	//get cp [][]
	for idx, sect := range f.Sections {
		styp := 1
		if len(f.Styps) == len(f.Sections){
			styp = f.Styps[idx]
		}
		bar := CalcSecProp(styp, sect)
		f.Secprop = append(f.Secprop, bar)
		f.Mod.Cp = append(f.Mod.Cp, []float64{bar.Area/1e6, bar.Ixx/1e12})
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
	//f.CBeams = matrix of [flrdx][spandx]; each floor vec being [b1, b2...bn]
	//lc has to be the first mem and rcmem has to be the last, figure - NOT DONE
	bmstrt := f.Nflrs* xstep
	f.CBeams = make([][]int, f.Nflrs)
	for i := 0; i < f.Nflrs; i++{
		f.CBeams[i] = make([]int, f.Nspans)
	}
	for j := 0; j < f.Nflrs; j++{
		for i := 0; i < f.Nspans; i++{
			f.CBeams[j][i] = bmstrt + j * f.Nspans + i + 1
		}
	}
	//fmt.Println("printing CBEAM vec->")
	//fmt.Println(f.CBeams)
	return
}

//calc funcs ()

//InitSections initializes frm2d sections
func (f *Frm2d) InitSections(cdims, bdims, cldims [][]float64){
	//WILL IT CALC FLANGE WIDTH AND REARRANGE DIMS (yes?)
	//cdims and bdims will usually be generated by opt
	//if both cstyp and bstyp are 0 then both are not set
	if f.Cstyp == 0 && f.Bstyp == 0{
		f.Cstyp = 1; f.Bstyp = 1
	}
	ncols := len(f.X); nbms := len(f.X) - 1; nclr := 2
	f.Sections = make([][]float64, ncols+nbms+nclr)
	f.Styps = make([]int, ncols+nbms+nclr)
	for i := 0; i < ncols; i++{
		f.Sections[i] = make([]float64, len(cdims[i]))
		copy(f.Sections[i],cdims[i])
		f.Styps[i] = f.Cstyp
	}
	for i := 0; i < nbms; i++{
		f.Sections[i+ncols] = make([]float64, len(bdims[i]))
		copy(f.Sections[i], bdims[i])
		f.Styps[i+ncols] = f.Bstyp
	}
	for i := 0; i < nclr; i++{
		f.Sections[i+ncols+nbms] = make([]float64, len(cldims[i]))
		copy(f.Sections[i+ncols+nbms], cldims[i])
		f.Styps[i+ncols+nbms] = f.Clstyp
	}
	return
}

//AddSection adds a section to a Frm2d struct
func (f *Frm2d) AddSection(styp int, sdim []float64) (sid int){
	//has this been used i think not
	f.Sections = append(f.Sections, sdim)
	f.Styps = append(f.Styps, styp)
	return len(f.Sections)
}

//AddNode adds a node at x, y; returns an error if node exists
func (f *Frm2d) AddNode(x,y float64) (err error){
	p := Pt{X:x,Y:y}
	if _, ok := f.Nodes[p]; ok{
		return errors.New("node exists")
	} else {
		f.Mod.Coords = append(f.Mod.Coords, []float64{x,y})
		f.Nodes[p] = []int{len(f.Mod.Coords)}
	}
	return
}

//AddSelfWeight adds frame self weight given dx (1 - add beam self wt, 2 - add beam + col self wt)
func (f *Frm2d) AddSelfWeight(dx int) (err error){
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
				if f.Fltslb{
					wdl = f.Pg * b * d/1e6
				} else {
					wdl = f.Pg * b * (d - f.Dslb)/1e6
				}
				case 6,7,8,9,10:
				//bdim := f.Sections[bdx-1]; 
				bw := bdim[2]; dw:= bdim[1] - bdim[3]
				wdl = f.Pg * bw * dw /1e6			
				default:
				//HAHA.
				wdl = f.Pg * f.Secprop[bdx-1].Area/1e6
				
			}
			default:
			wdl = f.Pg * f.Secprop[bdx-1].Area/1e6
			//log.Println(ColorYellow,"bdx,wdl->",bdx,wdl,ColorReset)
		}
		//log.Println(ColorYellow,"bdx,wdl->",bdx,wdl,ColorReset)
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
			//log.Println("col, wdl->",bdx, wdl)
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

//AddMemLoad adds member loads from Frm2d loadcases
func (f *Frm2d) AddMemLoad(mem int, ldcase []float64) (err error){
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

//AddBeamLLDL adds beam live and dead (udl) loads
func (f *Frm2d) AddBeamLLDL(){
	if f.DL + f.LL > 0.0{
		for _, bm := range f.Beams{
			mvec := f.Members[bm]; mtyp := mvec[1][0]
			jb := f.Mod.Coords[mvec[0][0]-1]; je := f.Mod.Coords[mvec[0][1]-1]
			lspan := Dist2d(jb, je)
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
				switch f.Slbload{
					case 0:
					//add udl of DL as ldtyp 3
					_ = f.AddMemLoad(bm, []float64{1.0,3.0,f.DL,0.0,0.0,0.0,1.0})
					if f.LL > 0.0{
						_ = f.AddMemLoad(bm, []float64{1.0,3.0,f.LL,0.0,0.0,0.0,2.0})
					}
					case 1:
					//one way slab load - over lbay/2.0 * 2.0
					var dl float64
					if f.Dslb > 0.0{
						dl += f.Pg * f.Dslb * f.Lbay * 1e-3
					}
					dl += f.DL * f.Lbay; ll := f.LL * f.Lbay
					if dl > 0.0{
						_ = f.AddMemLoad(bm, []float64{1.0,3.0,dl,0.0,0.0,0.0,1.0})
					}
					if ll > 0.0{
						_ = f.AddMemLoad(bm, []float64{1.0,3.0,ll,0.0,0.0,0.0,2.0})
					}
					if f.Spam{fmt.Println("frame dl, ll->",dl, ll)}
					case 2:
					//two way slab load - 3 tri right, udl, tri left
					var dl float64
					if f.Dslb > 0.0{
						dl += f.Pg * f.Dslb * 1e-3
					}
					dl += f.DL; ll := f.LL
					switch{
						case f.Lbay >= lspan:
						//triangular load of 2.0 * wlspan/4.0
						dl = dl * lspan
						ll = ll * lspan
						_ = f.AddMemLoad(bm, []float64{1.0,4.0,0.0,dl,0.0,lspan/2.0,1.0})
						_ = f.AddMemLoad(bm, []float64{1.0,4.0,dl,0.0,lspan/2.0,0.0,1.0})
						if ll > 0.0{
							_ = f.AddMemLoad(bm, []float64{1.0,4.0,0.0,ll,0.0,lspan/2.0,2.0})
							_ = f.AddMemLoad(bm, []float64{1.0,4.0,ll,0.0,lspan/2.0,0.0,2.0})
						}					
						case f.Lbay < lspan:
						//
						//trap load of peak 2.0 * wlbay2/2.0
						dl = dl * f.Lbay; ll = ll * f.Lbay
						if dl > 0.0{
							_ = f.AddMemLoad(bm, []float64{1.0,4.0,0.0,dl,0.0,lspan-f.Lbay/2.0,1.0})
							_ = f.AddMemLoad(bm, []float64{1.0,3.0,dl,0.0,f.Lbay/2.0,f.Lbay/2.0,1.0})
							_ = f.AddMemLoad(bm, []float64{1.0,4.0,dl,0.0,lspan-f.Lbay/2.0,0.0,1.0})
						}
						if ll > 0.0{
							_ = f.AddMemLoad(bm, []float64{1.0,4.0,0.0,ll,0.0,lspan-f.Lbay/2.0,2.0})
							_ = f.AddMemLoad(bm, []float64{1.0,3.0,ll,0.0,f.Lbay/2.0,f.Lbay/2.0,2.0})
							_ = f.AddMemLoad(bm, []float64{1.0,4.0,ll,0.0,lspan-f.Lbay/2.0,0.0,2.0})
						}						
					}				
					case 11:
					//triangular one-way slab load (y-frame)
					case 22:
					//triangular two - way slab load (y-frame)
				}
			}
		}
	}
	return
}

//GenLoads generates load cases for Frm2d (load types 1- dl, 2-ll, 3-wl)
func (f *Frm2d) GenLoads() (err error){
	f.AddBeamLLDL()
	xstep := len(f.X); nflrs := len(f.Y)-1
	f.Nlp = f.Nspans + 1
	//nspans := f.Nspans
	//bmstrt := nflrs * xstep
	if len(f.Bloads) != 0 {
		//bload mems are (1)- beam no. 1, beam no. 2, etc
		for _, ldcase := range f.Bloads{
			mem := int(ldcase[0]) + nflrs*xstep
			_ = f.AddMemLoad(mem, ldcase)
		}
	}
	if len(f.Cloads) != 0 {
		for _, ldcase := range f.Cloads{
			mem := int(ldcase[0]) 
			_ = f.AddMemLoad(mem, ldcase)
		}
	}
	switch f.Ldcalc{
		case 1:
		//only load case 0, 1
		f.DM = 0.0
		f.Nlp = 2
		for i := 1; i <= f.Nspans; i++ {
			for j := 0; j < f.Nflrs; j++{
				mem := f.CBeams[j][i-1]
				f.Loadcons[0] = append(f.Loadcons[0], f.Advloads[mem]...)
			}
		}
		if f.Lclvr{
			for _, mem := range f.Lcmem{
				f.Loadcons[0] = append(f.Loadcons[0], f.Advloads[mem]...)
			}
		}
		if f.Rclvr{
			for _, mem := range f.Rcmem{
				f.Loadcons[0] = append(f.Loadcons[0], f.Advloads[mem]...)
			}
		}
		default:		
		//build load patterns
		for i := 1; i <= f.Nspans; i++ {
			for j := 0; j < f.Nflrs; j++{
				mem := f.CBeams[j][i-1]
				f.Loadcons[0] = append(f.Loadcons[0], f.Advloads[mem]...)
				if i % 2 == 0 {
					f.Loadcons[1] = append(f.Loadcons[1], f.Benloads[mem]...)
					f.Loadcons[2] = append(f.Loadcons[2], f.Advloads[mem]...)
				} else {
					f.Loadcons[1] = append(f.Loadcons[1], f.Advloads[mem]...)
					f.Loadcons[2] = append(f.Loadcons[2], f.Benloads[mem]...)
				}
			}
		}
		for i := 1; i <= f.Nspans - 1; i++ {
			lp := i + 2
			for j := 0; j < f.Nflrs; j++{
				flrbm := f.CBeams[j]; mind := flrbm[i-1]
				for _, mem := range flrbm {
					if mem == mind || mem == mind + 1 {
						f.Loadcons[lp] = append(f.Loadcons[lp], f.Advloads[mem]...)
					} else {
						f.Loadcons[lp] = append(f.Loadcons[lp], f.Benloads[mem]...)
					}
				}
			}
		}
		if f.Lclvr{
			f.Nlp++
			for lp := 0; lp <= 2; lp++{
				switch lp{
					case 0,2:
					for _, mem := range f.Lcmem{
						f.Loadcons[lp] = append(f.Loadcons[lp],f.Advloads[mem]...)
					}
					case 1:
					for _, mem := range f.Lcmem{
						f.Loadcons[lp] = append(f.Loadcons[lp], f.Benloads[mem]...)
					}
				}
			}
			for _, lcmem := range f.Lcmem{
				fdx := f.Members[lcmem][1][2]
				lmem := f.CBeams[fdx-1][0]
				f.Loadcons[f.Nspans+2] = append(f.Loadcons[f.Nspans+2], f.Advloads[lcmem]...)
				f.Loadcons[f.Nspans+2] = append(f.Loadcons[f.Nspans+2], f.Advloads[lmem]...)
				for _, flrmem := range f.CBeams[fdx-1][1:]{
					f.Loadcons[f.Nspans+2] = append(f.Loadcons[f.Nspans+2], f.Benloads[flrmem]...)
				}
			}
		}
		if f.Rclvr{
			f.Nlp++
			for lp := 0; lp <= 2; lp++{
				switch lp{
					case 0,1:
					for _, mem := range f.Rcmem{
						f.Loadcons[lp] = append(f.Loadcons[lp], f.Advloads[mem]...)
					}
					case 2:
					for _, mem := range f.Rcmem{
						f.Loadcons[lp] = append(f.Loadcons[lp], f.Benloads[mem]...)
					}
				}
			}
			for _, rcmem := range f.Rcmem{
				fdx := f.Members[rcmem][1][2]
				rmem := f.CBeams[fdx-1][len(f.CBeams[fdx-1])-1]
				f.Loadcons[f.Nspans+3] = append(f.Loadcons[f.Nspans+3], f.Advloads[rcmem]...)
				f.Loadcons[f.Nspans+3] = append(f.Loadcons[f.Nspans+3], f.Advloads[rmem]...)
				for _, flrmem := range f.CBeams[fdx-1][:f.Nspans-1]{
					f.Loadcons[f.Nspans+3] = append(f.Loadcons[f.Nspans+3], f.Benloads[flrmem]...)
				}
			}
		}
	}
	for lp := range f.Loadcons{
		for _, col := range f.Cols{
			//TO DO- CANNOT BE ADVLOADS ALL THE TIME
			f.Loadcons[lp] = append(f.Loadcons[lp], f.Advloads[col]...)
		}
	}
	if f.Nloads == 3 || len(f.WL) != 0{
		//add wind load
		err = f.AddWindLoad()
		if err != nil{
			fmt.Println(ColorRed,err,ColorReset)
			return
		}
		switch f.Ldcalc{
			case 1:
			f.Loadcons[0] = [][]float64{}
			for _, mem := range f.Beams{
				f.Loadcons[0] = append(f.Loadcons[0], f.Wadloads[mem]...)			
			}
			for _, col := range f.Cols{
				//TO DO- CANNOT BE ADVLOADS ALL THE TIME
				f.Loadcons[0] = append(f.Loadcons[0], f.Advloads[col]...)
			}
			f.Loadcons[1] = [][]float64{}
			for _, mem := range f.Beams{
				f.Loadcons[1] = append(f.Loadcons[1], f.Wadloads[mem]...)			
			}
			for _, col := range f.Cols{
				//TO DO- CANNOT BE ADVLOADS ALL THE TIME
				f.Loadcons[1] = append(f.Loadcons[1], f.Advloads[col]...)
			}
			default:
			for _, mem := range f.Beams{
				f.Loadcons[f.Nlp+1] = append(f.Loadcons[f.Nlp+1], f.Wadloads[mem]...)
				f.Loadcons[f.Nlp+2] = append(f.Loadcons[f.Nlp+2], f.Wadloads[mem]...)		
			}
			for _, col := range f.Cols{
				//TO DO- CANNOT BE ADVLOADS ALL THE TIME
				f.Loadcons[f.Nlp+1] = append(f.Loadcons[f.Nlp+1], f.Advloads[col]...)
				f.Loadcons[f.Nlp+2] = append(f.Loadcons[f.Nlp+2], f.Advloads[col]...)
			}
		}
		wsf1 := f.WLFs[2]
		for _, nl := range f.Wloadl{
			f.Jloadmap[-3] = append(f.Jloadmap[-3],nl)
			switch f.Ldcalc{
				case 1:
				//single load case
				f.Jloadmap[0] = append(f.Jloadmap[0], []float64{nl[0],nl[1]*wsf1,nl[2],nl[3]})
				
				default:
				//patterns + wl
				f.Jloadmap[f.Nlp+1] = append(f.Jloadmap[f.Nlp+1], []float64{nl[0],nl[1]*wsf1,nl[2],nl[3]})
			}			
			if f.Verbose{
				//fmt.Println(ColorRed,"adding wind load")
				//fmt.Printf("node - %.f load - %.2f\n",nl[0],nl[1])
				//fmt.Printf("vec - %.f\n",nl)
				//fmt.Println(ColorReset)				
			}
		}
		for _, nl := range f.Wloadr{
			//f.Uniloads[-4] = append(f.Uniloads[-4],nl)
			f.Jloadmap[-4] = append(f.Jloadmap[-4],nl)
			switch f.Ldcalc{
				case 1:
				//single load case
				f.Jloadmap[1] = append(f.Jloadmap[1], []float64{nl[0],nl[1]*wsf1,nl[2],nl[3]})
				
				default:
				//patterns + wl
				f.Jloadmap[f.Nlp+2] = append(f.Jloadmap[f.Nlp+2], []float64{nl[0],nl[1]*wsf1,nl[2],nl[3]})
			}			
			if f.Verbose{
				//fmt.Println(ColorRed,"adding wind load")
				//fmt.Printf("node - %.f load - %.2f\n",nl[0],nl[1])
				//fmt.Printf("vec - %.f\n",nl)
				//fmt.Println(ColorReset)				
			}
		}
	}
	f.Loadcons[-1] = f.Uniloads[1]
	f.Loadcons[-2] = f.Uniloads[2]
	f.Loadcons[-3] = [][]float64{}
	f.Loadcons[-4] = [][]float64{}
	return
}

//AddWindLoad adds wind loads from frm.WL slice to frm wind load left, right slices
func (f *Frm2d) AddWindLoad() (err error){
	/*
	   wl - vec [0,1,2,3,4] - nflrs
	*/
	switch len(f.WL){
		case 0:
		return errors.New("no wind load specified")
		case 1:
		wl := f.WL[0]
		for i, node := range f.Ljs{
			rnode := f.Rjs[i]
			f.Wloadl = append(f.Wloadl, []float64{float64(node),wl,0,0})
			f.Wloadr = append(f.Wloadr, []float64{float64(rnode),-wl,0,0})
			
		}
		case f.Nflrs:
		//NOTE - 0 level load is 0
		for i, node := range f.Ljs{
			rnode := f.Rjs[i]
			wl := f.WL[i]
			f.Wloadl = append(f.Wloadl, []float64{float64(node),wl,0,0})
			f.Wloadr = append(f.Wloadr, []float64{float64(rnode),-wl,0,0})
		}
		default:
		return errors.New("invalid length of wind load slice")
	}
	return nil
}

//InitMemRez initializes member result maps (Bmenv and Colenv)
func (f *Frm2d) InitMemRez(){
	xstep := f.Nspans + 1
	bmstrt := xstep * f.Nflrs
	bmendc := 2
	if f.Bmrel == 4 || (xstep == 2 && f.Ncls == 0){
		bmendc = 1
	}
	f.Bmenv = make(map[int]*BmEnv)
	for _, i := range f.Beams{
		xdx, fdx := f.Members[i][1][1],f.Members[i][1][2]
		jb, je := f.Mod.Mprp[i-1][0], f.Mod.Mprp[i-1][1]
		c1, c2 := f.Mod.Coords[jb-1], f.Mod.Coords[je-1]
		cl := i - bmstrt
		cr := cl + 1
		var lsx, rsx float64
		mtyp := f.Members[i][1][0]
		cp := f.Members[i][0][3]
		styp := f.Styps[cp-1]
		switch mtyp{
			case 3:
			rsx = f.Secmap[1][1]/2.0
			case 4:
			lsx = f.Secmap[xstep][1]/2.0
			default:
			lsx = f.Secmap[cl][1]/2.0; rsx = f.Secmap[cr][1]/2.0
		}
		//fmt.Println("xstep->",xstep,"nspans-",f.Nspans,"floors-",f.Nflrs)
		//fmt.Println("beam->",i,"xdx-",xdx,"fdx-",fdx)
		spandx := 0; plinth := false; roof := false
		switch xdx {
		case 0:
			spandx = -1
		case f.Nspans - 1:
			spandx = -2
		}
		switch fdx{
			case 1:
			plinth = true
			case f.Nflrs:
			roof = true
		}
		f.Bmenv[i] = &BmEnv{
			Id:i,
			Foldr:f.Foldr,
			EnvRez:make(map[int]BeamRez),
			Venv:make([]float64,21),
			Mpenv:make([]float64,21),
			Mnenv:make([]float64,21),
			Vrd:make([]float64,21),
			Mnrd:make([]float64,21),
			Mprd:make([]float64,21),
			Dims:f.Secmap[i],
			Coords:[][]float64{c1,c2},
			Lsx:lsx/1000.0,Rsx:rsx/1000.0,
			Endc:bmendc,
			Lspan:Dist2d(c2,c1),
			Styp:styp,
			Cp:cp,
			Spandx:spandx,
			Plinth:plinth,
			Roof:roof,
			Hc:f.Hc,
			Dh:f.Dh,
			Drx:f.Drx,
			Dry:f.Dry,
			Drd:f.Drd,
			Term:f.Term,
		}
	}
	
	f.Colenv = make(map[int]*ColEnv)
	for _, i := range f.Cols{
		//fmt.Println("xdx,fdx,locdx,lex,mtyp->",f.Members[i][1])
		var ljbase,lbss, ubss bool
		jb, je := f.Mod.Mprp[i-1][0], f.Mod.Mprp[i-1][1]
		c1, c2 := f.Mod.Coords[jb-1], f.Mod.Coords[je-1]
		cp := f.Members[i][0][3]
		lspan := Dist2d(c2,c1)
		kcol := f.Secprop[cp-1].Ixx/lspan/1e3
		//f.Members[mdx] = append(f.Members[mdx],[]int{xdx, fdx, locdx, lex, mtyp})
		xdx, fdx := f.Members[i][1][1],f.Members[i][1][2]
		//fmt.Println("col->",i, jb, je)
		//fmt.Println("lspan",lspan,"kcol",kcol)
		var lbd, ubd [][]float64
		var lst, ust []int
		var lbk, ubk []float64
		var dl, dt float64
		var ct, cb, ub1, ub2, lb1, lb2 int
		lb1 = bmstrt + 1 + (fdx - 1) * (xstep - 1) + (xdx - 1)
		lb2 = lb1 + 1
		ub1 = bmstrt + 1 + (fdx) * (xstep - 1) + (xdx - 1)
		ub2 = ub1 + 1
		ct = i + xstep; cb = i - xstep
		switch xdx{
			case 0:
			lb1, ub1 = 0, 0
			case xstep -1:
			lb2, ub2 = 0, 0 
			
		}
		switch fdx{
			case 0:
			lb1, lb2 = 0, 0
			ljbase = true
			cb = 0
			case f.Nflrs -1:
			ct = 0

		}
		if f.Nflrs == 1{ct = 0}
		//fmt.Println("col->",i,"xdx, fdx-",xdx, fdx, "beams->",ub1,ub2,lb1,lb2)
		for i, bm := range []int{lb1, lb2, ub1, ub2}{
			var kb, bmd float64
			var styp int
			var dims []float64
			switch bm{
				case 0:
				dims = []float64{0,0,0,0,0,0}
				styp = -1
				bmd = 0.0
				default:
				dims = f.Bmenv[bm].Dims
				styp = f.Styps[cp-1]
				
			}
			if bm > 0{
				switch styp{
				//note - lspan is in METERS 
					case 6,7,8,9,10:
					//use rect beam stiffness - ONLY IF RCC (change this)
					kb = dims[2] * math.Pow(dims[1],3.0)/12.0
					kb = kb/f.Bmenv[bm].Lspan/1e3
					default:
					//use reg stiffness = i/lspan
					cp := f.Bmenv[bm].Cp
					kb = f.Secprop[cp-1].Ixx/f.Bmenv[bm].Lspan/1e3
				}
			}
			bmd = dims[0]; if len(dims)>1 {bmd = dims[1]}
			//fmt.Println("beam->",bm,bmd)
			if i < 2{
				lbd = append(lbd, dims)
				lbk = append(lbk, kb)
				lst = append(lst, styp)
				if dl <= bmd/2.0 {dl = bmd/2.0}
				if bm != 0 && f.Bmenv[bm].Endc < 2{lbss = true}
			} else {
				ubd = append(ubd, dims)
				ubk = append(ubk, kb)
				ust = append(ust, styp)
				if dt <= bmd/2.0 {dt = bmd/2.0}
				if bm != 0 && f.Bmenv[bm].Endc < 2{ubss = true}
			}
			
		}
		//fmt.Println("col->",i,"xdx, fdx-",xdx, fdx, "beams->",ub1,ub2,lb1,lb2,"ct, cb->",ct,cb)
		var ksum, kcsum float64
		kcsum = kcol
		
		switch f.Code{
			case 1:
			if f.Braced{
				ksum = kcol + 0.5 * (lbk[0] + lbk[1])
			} else {
				ksum = kcol + 1.5 * (lbk[0] + lbk[1])
			}	
			default:
			ksum = lbk[0] + lbk[1]
		}
		var b1, b2 float64
		if cb != 0{
			jb1, je1:= f.Mod.Mprp[cb-1][0], f.Mod.Mprp[cb-1][1]
			p1, p2 := f.Mod.Coords[jb1-1], f.Mod.Coords[je1-1]
			cp1 := f.Members[cb][0][3]
			kcb := f.Secprop[cp1-1].Ixx/Dist2d(p1,p2)/1e3
			kcsum += kcb
			if f.Code == 1 {ksum += kcb}
		}
		if ksum != 0{
			b1 = kcsum/ksum
		}
		switch f.Code{
			case 1:
			if f.Braced{
				ksum = kcol + 0.5 * (ubk[0] + ubk[1])
			} else {
				ksum = kcol + 1.5 * (ubk[0] + ubk[1])
			}
			
			default:
			ksum = ubk[0] + ubk[1]
		}
		kcsum = kcol
		if ct != 0{
			jb1, je1:= f.Mod.Mprp[ct-1][0], f.Mod.Mprp[ct-1][1]
			p1, p2 := f.Mod.Coords[jb1-1], f.Mod.Coords[je1-1]
			cp1 := f.Members[ct][0][3]
			kct := f.Secprop[cp1-1].Ixx/Dist2d(p1,p2)/1e3

			kcsum += kct
			if f.Code == 1{ksum += kct}
			//fmt.Println("col->",jb1,je1,cp1,kct, i)
		}
		if ksum != 0{b2 = kcsum/ksum}
		f.Colenv[i] = &ColEnv{
			Id:i,
			Foldr:f.Foldr,
			EnvRez:make(map[int][]float64),
			Dims:f.Secmap[i],
			Coords:[][]float64{c1,c2},
			Lb: []int{lb1, lb2},
			Ub: []int{ub1, ub2},
			Lbk:lbk,
			Ubk:ubk,
			Lbd: lbd,
			Ubd: ubd,
			Lst: lst,
			Ust: ust,
			Ljbase:ljbase,
			Styp:f.Styps[cp-1],
			Lspan:lspan,
			Fixbase:f.Fixbase,
			Braced:f.Braced,
			Kcol:kcol,
			L0:lspan - dl/1e3 - dt/1e3,
			Lbss:lbss,
			Ubss:ubss,
			B1:b1,
			B2:b2,
			Term:f.Term,
			Title:f.Title,
		}
		//if rcc calc effective height using bs8110
		switch f.Mtyp{
			case 1:
			//rcc
			f.Colenv[i].EffHt(2)
		}
	}
	f.Ms = make(map[int]*Mem)
	f.Mslmap = make(map[int]map[int][][]float64)
	//TODO check for wind load calcs
	return
}

//CalcLoadEnv loads model with individual load cases and sorts results 
func (f *Frm2d) CalcLoadEnv(){
	if f.Verbose{
		/*
		for lp, ldcons := range f.Loadcons{
			fmt.Println(ColorRed, "load patterns->",lp,ColorReset)
			for _, ldcase := range ldcons{
				fmt.Println("load->",ldcase)
			}
		}
		*/
	}
	mod := &Model{
		Cmdz:[]string{"2df","mks","1"},
		Units:"kn-m",
		Coords: f.Mod.Coords,
		Supports: f.Mod.Supports,
		Em: f.Mod.Em,       
		Cp: f.Mod.Cp,       
		Mprp:f.Mod.Mprp,
		Foldr:f.Foldr,
		Web:f.Web,
		Noprnt:f.Noprnt,
		Dims:f.Sections,
		Sts:f.Styps,
	}
	for lp, ldcons := range f.Loadcons{
		//calc only max load and unit load cases for Ldcalc = 1
		//add unit load calcs why not
		if f.Ldcalc == 1 && lp > 1{
			continue
		}
		//fmt.Println(lp, ldcons)
		if val, ok := f.Jloadmap[lp]; ok{
			mod.Jloads = val
		} else {
			mod.Jloads = [][]float64{}
		}
		mod.Id = fmt.Sprintf("%s_lp_%v",f.Title,lp)
		mod.Msloads = ldcons
		
		frmrez, err := CalcFrm2d(mod, 3)
		if err != nil{return}
		ms,_ := frmrez[1].(map[int]*Mem)
		msloaded, _ := frmrez[5].(map[int][][]float64)
		spanchn := make(chan BeamRez,len(msloaded))
		report,_ := frmrez[6].(string)
		f.Reports = append(f.Reports, report)
		f.Mslmap[lp] = msloaded
		switch lp{
			case 0:
			f.Mod.Msloads = ldcons
			f.Ms = ms
			if val, ok := f.Jloadmap[lp]; ok{
				f.Mod.Jloads = val
			}
			
		}
		for _, id := range f.Cols{
			//(todo) create a func to sort member calc results
			mem := ms[id]
			cm := f.Colenv[id]
			cm.EnvRez[lp] = mem.Qf
			if math.Abs(mem.Qf[2]) > math.Abs(cm.Mbmax) {cm.Mbmax = mem.Qf[2]}
			if math.Abs(mem.Qf[5]) > math.Abs(cm.Mtmax) {cm.Mtmax = mem.Qf[5]}
			if math.Abs(mem.Qf[0]) > math.Abs(cm.Pumax) {cm.Pumax = mem.Qf[0]}
		}
		for _, id := range f.Beams{
			ldcase := msloaded[id]
			go BeamFrc(3, id, ms[id], ldcase, spanchn, false)
		}
		//for id, ldcase := range msloaded{
		//	go BeamFrc(3, id, ms[id], ldcase, spanchn, false)
		//}
		for _ = range f.Beams{
			r := <- spanchn
			id := r.Mem
			
			bm := f.Bmenv[id]
			
			bm.EnvRez[lp] = r
			if len(bm.Xs) == 0 {
				bm.Xs = r.Xs
			}
			xdiv := ms[id].Geoms[0]/20.0
			lsx := bm.Lsx; rsx := bm.Lspan - bm.Rsx
			il := int(math.Ceil(lsx/xdiv)); ir := int(math.Ceil(rsx/xdiv))
			var vl, vr, ml, mr float64
			for i, vx := range r.SF{
				x := r.Xs[i]
				if i == il{
					switch{
						case x == lsx:
						vl = vx
						ml = r.BM[i]
						default:
						vl = vx + (vx - r.SF[i-1])*(lsx - x)/xdiv
						ml = r.BM[i] + 0.5 * (lsx - x)*(vl + vx)
					}
					if math.Abs(bm.Vl) < math.Abs(vl){bm.Vl = vl}
					if math.Abs(bm.Ml) < math.Abs(ml){bm.Ml = ml}
				}
				if i == ir{
					switch{
						case x == rsx:
						vr = vx
						mr = r.BM[i]
						default:
						vr = vx + (vx - r.SF[i-1])*(rsx - x)/xdiv
						mr = r.BM[i] + 0.5 * (rsx - x)*(vr + vx) 
					}
					if math.Abs(bm.Vr) < math.Abs(vr){bm.Vr = vr}
					if math.Abs(bm.Mr) < math.Abs(mr){bm.Mr = mr}
				}
				if math.Abs(bm.Venv[i]) < math.Abs(vx) {
					bm.Venv[i] = vx
					if math.Abs(bm.Vmax) < math.Abs(vx){
						bm.Vmax = vx
						bm.Vmaxx = r.Xs[i]
					}
				}
				if math.Abs(bm.Mnenv[i]) < math.Abs(r.BM[i]) && r.BM[i] < 0.0 {
					bm.Mnenv[i] = r.BM[i]
					if math.Abs(bm.Mnmax) < math.Abs(r.BM[i]){
						bm.Mnmax = r.BM[i]
						bm.Mnmaxx = r.Xs[i]
					}
				}
				if r.BM[i] > 0.0 && math.Abs(bm.Mpenv[i]) < math.Abs(r.BM[i]) {
					bm.Mpenv[i] = r.BM[i]
					if math.Abs(bm.Mpmax) < math.Abs(r.BM[i]){
						bm.Mpmax = r.BM[i]
						bm.Mpmaxx = r.Xs[i]
					}
				}
			}
		}
	}
}

//genxs generates x values from lspans lmao
func genxs(lspans []float64)(xs []float64){
	x := 0.0
	xs = append(xs, x)
	for _, span := range lspans{
		x += span
		xs = append(xs, x)
	}
	return
}

//InitMat generates frame material prop/grade values
func (f *Frm2d) InitMat(){
	/*
	   fcks - [col, bm, ftng, slb]
	   //fys -  [ftng, col, bm, slb]
	   //fyds - [ftng, col, bm, slb]
	   //mod.Em - [ecol, ebeam]
	*/
	switch f.Mtyp{
		case 1:
		//fmt.Println("rcc frame")
		f.Pg = 25.0 
		//if f.Code == 2{f.Pg = 24.0}
		if f.Fcks == nil && f.Fys == nil {
			//defaults - m25, fe415, fe415
			f.Fcks = append(f.Fcks, 25.0)
			f.Fys = append(f.Fys, 415.0)
			f.Fyds = append(f.Fyds, 415.0)
		}
		if len(f.Fcks) == 1{
			//col, beam, footing, slab
			f.Fcks = []float64{f.Fcks[0],f.Fcks[0],f.Fcks[0],f.Fcks[0]}
		}
		if len(f.Fys) == 1{
			f.Fys = []float64{f.Fys[0],f.Fys[0],f.Fys[0],f.Fys[0]}
		}
		if len(f.Fyds) == 1{
			f.Fyds = []float64{f.Fys[0],f.Fys[0],f.Fys[0],f.Fys[0]}
		}
		f.Mod.Em = make([][]float64,2)
		f.Mod.Em[0] = []float64{FckEm(f.Fcks[0])}
		f.Mod.Em[1] = []float64{FckEm(f.Fcks[1])}
		f.Mod.Units = "knm"
		case 2:
		//steel frame
		f.Pg = 7850.0 //or 7850? FINALIZE STEEL UNITZ
		f.Mod.Em = [][]float64{{2e5}}
		case 3:
		//timber frame lmao
		//f.Fcks is dzvals
		case 4:
		//teehee cfs really
		case 5:
		//is bamboo even happens
	}
}


//getmemvec gets location indices for a 2d frame
//xdx (x- location), fdx - floor loc, locdx - location dx, lex - end condition dx, mtyp - beam/col/clvr
func getmemvec(fgen, mt, jb, je, xstep int) (xdx, fdx, locdx, lex, mtyp int){
	//the (future) backbone and source of all calcs
	//CANTILEVERS WILL BE THE DEATH OF THIS
	xdx = (jb-1)%xstep
	fdx = (jb-1)/xstep
	//fmt.Println("memtyp->",mt,"jb,je->",jb,je,"xdx,fdx->",xdx, fdx)
	switch mt{
		case 1:
		//col
		locdx = 1
		lex = 1
		switch xdx{
			case 0:
			locdx = 2
			lex = 2
			//case 1:
			//locdx = 4
			//case xstep-1:
			//locdx = 5
			case xstep-1:
			locdx = 3
			lex = 2
		}
		case 2:
		//beam - int, end left, end right
		locdx = 1
		lex = 1
		switch xdx{
			case 0:
			locdx = 2
			lex = 2
			case xstep-2:
			locdx = 3
			lex = 2
		}
		case 3:
		//left cantilever
		lex = 1
		switch xdx{
			case 0:
			locdx = 1
			case xstep-1:
			locdx = 2
		}
		case 4:
		lex = 1
		switch xdx{
			case 0:
			locdx = 1
			case xstep-1:
			locdx = 2
		}
	}
	switch fgen{
		case 0:
		//as is; just returning this randomly
		ncols := xstep
		mtyp = (mt-1)*ncols+locdx
		if mt == 3 || mt == 4 {mtyp = 2 * xstep - 1 + locdx}

		case 1:
		//one col, one beam, one cantilever
		mtyp = mt
		case 2:
		//classify by lex
		ncols := 2
		mtyp = (mt-1)*ncols+lex
		if mt == 3{mtyp = 4 + lex}
		// case 3:
		// //classify by locdx
		// ncols := 3
		// mtyp = (mt-1)*ncols+locdx
		// if mt == 3{mtyp = 6 + locdx}
		case 3:
		//all diff
		ncols := xstep
		mtyp = (mt-1)*ncols+locdx
		if mt == 3 || mt == 4 {mtyp = 2 * xstep - 1 + locdx}
	}
	return
}

//Printz prints horribly
func (f *Frm2d) Printz(){
	//printz should return (tablewriter) table/report
	//REDO DIS
	fmt.Println("frame 2d")
	fmt.Println("grade of concrete, steel-",f.Fcks, f.Fys)
	fmt.Println("nspans-",f.Nspans)
	fmt.Println("dl, ll-",f.DL, f.LL)
	fmt.Println("clvrs-",f.Clvrs)
}

//Dump saves frame json to data/out/frame2d_name.json
func (f *Frm2d) Dump(name string) (filename string, err error){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	if name == "" {name = fmt.Sprintf("frame2d_%v",f.Id)}
	name += ".json"
	//foldr := filepath.Join(basepath,"../data/out",time.Now().Format("06-Jan-02"))
	foldr := filepath.Join(basepath,"../data/out")
	if _, e := os.Stat(foldr); errors.Is(e, os.ErrNotExist) {
		e := os.Mkdir(foldr, os.ModePerm)
		if e != nil {
			err = e; return
		}
	}
	filename = filepath.Join(foldr,name)
	data, e := json.Marshal(f)
	if e != nil{err = e; return}
	err = ioutil.WriteFile(filename, data, 0644)
	return

}



//NOT NEEDED
//InitDims inits dimensions, use for opt routines (if needed (NOPE))
func (f *Frm2d) InitDims(){
	if f.Fgen == 0{
		f.Fgen = 1
	}
	if f.Cstyp == 0 && f.Bstyp == 0{
		f.Cstyp = 1; f.Bstyp = 1
	}
	if f.Bfcalc{
		if f.Dslb == 0.0{
			f.Dslb = 150.0
		}
		f.Fgen = 4
		if f.Bstyp == 0 || f.Bstyp == 1{
			f.Bstyp = 6
		}
	}
	f.Bdims = InitBdims(f.Df, f.X, f.Bvec, f.Code, f.Fgen, f.Bstyp, f.Nbms, f.Brel)
	f.Cdims = InitCdims(f.Cvec, f.Fgen, f.Ncols)
	//f.Cldims = InitCldims(f.Clvec, f.Nclr)
	f.Sections = make([][]float64,f.Nbms + f.Ncols)
	f.Styps = make([]int,f.Nbms + f.Ncols)
	for i := range f.Sections{
		switch {
		case i < f.Ncols:
			f.Styps[i] = f.Cstyp
			f.Sections[i] = make([]float64,len(f.Cdims[i]))
		default:
			f.Styps[i] = f.Bstyp
			f.Sections[i] = make([]float64,len(f.Bdims[i-f.Ncols-1]))
		}
	}
	// for i := range f.Sections{
	// 	if i < f.Ncols{
	// 		fmt.Println("col->",f.Styps[i],f.Sections[i])
	// 	} else {
	// 		fmt.Println("col->",f.Styps[i],f.Sections[i])
	// 	}
	// }
}

//InitCdims inits column dims
func InitCdims(cvec [][]float64, fgen, ncols int) (cdims [][]float64){
	cdims = make([][]float64, ncols)
	switch fgen{
		case 1:
		for i := 0; i < ncols; i++{
			cdims[i] = make([]float64, len(cvec[0]))
			copy(cdims[i],cvec[0])
		}
		case 2:
		for i := 0; i < ncols; i++{
			switch i {
			case 0, ncols -1:
				cdims[i] = make([]float64, len(cvec[0]))
				copy(cdims[i],cvec[0])
			default:
				cdims[i] = make([]float64, len(cvec[1]))
				copy(cdims[i],cvec[1])
			}
		}
		case 3:
		for i := 0; i < ncols; i++{
			switch i{
				case 0:
				cdims[i] = make([]float64, len(cvec[0]))
				copy(cdims[i],cvec[0])
				case ncols - 1:
				cdims[i] = make([]float64, len(cvec[1]))
				copy(cdims[i],cvec[1])
				default:
				cdims[i] = make([]float64, len(cvec[2]))
				copy(cdims[i],cvec[2])
			}
		}
		case 4:
		for i := 0; i < ncols; i++{
			cdims[i] = make([]float64, len(cvec[i]))
			copy(cdims[i],cvec[i])
		}
	}
	return
}

//InitBdims inits beam dims
func InitBdims(df float64, xs []float64, bvec [][]float64, code, fgen, bstyp, nbms, brel int) (bdims [][]float64){
	spans := make([]float64, len(xs)-1)
	var xl, xr, xint, tyb, xmin, xe float64
	for i := range spans{
		spans[i] = xs[i+1]-xs[i]
		switch i{
			case 0:
			xl = spans[i]; xmin = xl; xe = xl
			case nbms-1:
			xr = spans[i]
			if xmin > xr {xmin = xr}
			if xe > xr {xe = xr}
			default:
			if xint == 0{xint = spans[i]}
			if xint > spans[i] {xint = spans[i]}
			if xmin > xint {xmin = xint}
		}
	}
	bdims = make([][]float64, nbms)
	bvecf := make([][]float64, len(bvec))
	switch bstyp{
		case 1:
		for i := range bvec{
			bvecf[i] = make([]float64, len(bvec[i]))
			copy(bvecf[i], bvec[i])
		}
		case 6:
		//t section
		tyb = 1.0
		case 7,8,9,10:
		//l section
		tyb = 0.5
	}
	switch bstyp{
		case 6, 7, 8, 9, 10:		
		switch fgen{
			case 1:
			b := bvec[0][0]; d := bvec[0][1]; dfl := df; bw := b
			bf := GetBfRcc(code, brel, tyb, xmin, dfl, bw)
			for i := range bvecf{
				bvecf[i] = []float64{bf,d, b, df}
			}
			case 2:
			for i := range bvecf{
				switch i{
					case 0:
					b := bvec[0][0]; d := bvec[0][1]; dfl := df; bw := b
					bf := GetBfRcc(code, brel, tyb, xe, dfl, bw)
					bvecf[i] = []float64{bf,d, b, df}
					default:
					b := bvec[1][0]; d := bvec[1][1]; dfl := df; bw := b
					bf := GetBfRcc(code, brel, tyb, xint, dfl, bw)
					bvecf[i] = []float64{bf,d, b, df}
				}
			}
			case 3:
			for i := range bvecf{
				switch i{
					case 0:
					b := bvec[0][0]; d := bvec[0][1]; dfl := df; bw := b
					bf := GetBfRcc(code, brel, tyb, xl, dfl, bw)
					bvecf[i] = []float64{bf,d, b, df}
					case nbms-1:
					b := bvec[1][0]; d := bvec[1][1]; dfl := df; bw := b
					bf := GetBfRcc(code, brel, tyb, xr, dfl, bw)
					bvecf[i] = []float64{bf,d, b, df}
					default:
					b := bvec[2][0]; d := bvec[2][1]; dfl := df; bw := b
					bf := GetBfRcc(code, brel, tyb, xint, dfl, bw)
					bvecf[i] = []float64{bf,d, b, df}
				}
			}
			case 4:
			for i := range bvecf{
				b := bvec[i][0]; d := bvec[i][1]; dfl := df; bw := b
				bf := GetBfRcc(code, brel, tyb, spans[i], dfl, bw)
				bvecf[i] = []float64{bf,d, b, df}
			}
		}
	}
	switch fgen{
		case 1:
		for i := range bdims{
			bdims[i] = make([]float64, len(bvecf[0]))
			copy(bdims[i],bvecf[0])
		}
		case 2:
		for i := 0; i < nbms; i++{
			switch i {
			case 0, nbms -1:
				bdims[i] = make([]float64, len(bvecf[0]))
				copy(bdims[i],bvecf[0])
			default:
				bdims[i] = make([]float64, len(bvecf[1]))
				copy(bdims[i],bvecf[1])
			}
		}
		case 3:
		for i := 0; i < nbms; i++{
			switch i{
				case 0:
				bdims[i] = make([]float64, len(bvecf[0]))
				copy(bdims[i],bvecf[0])
				case nbms - 1:
				bdims[i] = make([]float64, len(bvecf[1]))
				copy(bdims[i],bvecf[1])
				default:
				bdims[i] = make([]float64, len(bvecf[2]))
				copy(bdims[i],bvecf[2])
			}
		}
		case 4:
		for i := 0; i < nbms; i++{
			bdims[i] = make([]float64, len(bvecf[i]))
			copy(bdims[i],bvecf[i])
		}
	}
	return
}


/*
   YEOLDE
   ...know ye o reader, eons past in the days of glory of the Hyboreans

//(todo) add this in mosh
func (f *Frm2d) Mrd(){
	//get list of beams per floor
	if f.DM != 0.0{
		fmt.Println(ColorRed,"moment redistribution DM->",f.DM, "DO THIS IN MOSH thanks",ColorReset)
		//CBeamDM(3, f.Beams, f.Bmenv, f.DM, f.Ms, f.Mslmap)
	}
	return
}

*/ 

//YE OLDE 




func (f *Frm2d) GenCp()(err error){
	//WHAT DO WE NEED THIS FOR
	//calcs section properties
	//adds self weight loads if f.Vec[1][1] == 1
	cplen := 3; ncol := 1; nbm := 1
	
	switch f.Ngrp{
		case 2:
		//2+2+1
		ncol = 2; nbm = 2
 		cplen = 5
		case 3:
		//3+3+2
		ncol = 3; nbm = 3
		cplen = 8
		case 4:
		//xstep + xstep - 1 + 2
		ncol = len(f.X); nbm = len(f.X)-1
		cplen = 2 * len(f.X) + 1
	}
	if len(f.Cdim) == 1{
		cdim := make([]float64, len(f.Cdim[0]))
		copy(cdim, f.Cdim[0])
		csec := f.Csec[0]
		f.Cdim = [][]float64{}
		f.Csec = []int{}
		for i := 0; i < ncol; i++{
			f.Cdim = append(f.Cdim, cdim)
			f.Csec = append(f.Csec, csec)
		}
	}
	if len(f.Bdim) == 1{
		bdim := make([]float64, len(f.Bdim[0]))
		copy(bdim, f.Bdim[0])
		bsec := f.Bsec[0]
		f.Bdim = [][]float64{}
		f.Bsec = []int{}
		for i := 0; i < nbm; i++{
			f.Bdim = append(f.Bdim, bdim)
			f.Bsec = append(f.Bsec, bsec)
		}
	}
	if len(f.Cdim) != ncol{
		return errors.New("mismatched length of column dims slice")
	}
	if len(f.Bdim) != nbm{
		return errors.New("mismatched length of beam dims slice")
	}
	f.Mod.Dims = append(f.Mod.Dims, f.Cdim...); f.Mod.Sts = append(f.Mod.Sts, f.Csec...)
	f.Mod.Dims = append(f.Mod.Dims, f.Bdim...); f.Mod.Sts = append(f.Mod.Sts, f.Bsec...)
	f.Mod.Cp = make([][]float64, cplen)
	return
}

func (f *Frm2d) GenMprpYeOlde() (err error){
	/*
	 generates mprp, msup
	*/
	xstep := len(f.X); nspans := xstep - 1
	switch nspans{
		case 1:
		//c-c, fgen = 1
		f.Ngrp = 1 
		case 2:
		//c-c-c, fgen = 2
		f.Ngrp = 2
	}
	var mdx int
	var sup []int
	switch f.Vec[2][0]{//footing vec
		case 1:
		sup = []int{-1,-1,-1}
		case 2:
		sup = []int{-1,-1,0}
	}
	for i := 1; i <= len(f.Mod.Coords); i++{
		if i <= xstep{
			//add supports
			supi := append([]int{i},sup...)
			f.Mod.Supports = append(f.Mod.Supports, supi)
		}
		if i + xstep <= len(f.Mod.Coords){
			//add columns
			mdx++
			xdx, fdx, locdx, lex, mtyp := getmemvec(f.Ngrp,1,i,i+xstep,xstep)
			f.Members[mdx] = append(f.Members[mdx],[]int{i, i + xstep, 1, mtyp, 0})
			f.Members[mdx] = append(f.Members[mdx],[]int{xdx, fdx, locdx, lex, mtyp})
			f.Mod.Mprp = append(f.Mod.Mprp, []int{i, i + xstep, 1, mtyp, 0})
		}
		if i + 1 <= len(f.Mod.Coords) && f.Mod.Coords[i-1][1] == f.Mod.Coords[i][1]{
			//add beams
			mdx++
			xdx, fdx, locdx, lex, mtyp := getmemvec(f.Ngrp,2,i,i+1,xstep)
			var brel int
			switch f.Bmrel{
				case 0, 1:
				brel = 0
				case 2:
				//both ends hinged (simply supported)
				brel = 4
			}
			f.Members[mdx] = append(f.Members[mdx],[]int{i, i + 1, 2, mtyp, brel})
			f.Members[mdx] = append(f.Members[mdx],[]int{xdx, fdx, locdx, lex, mtyp})
			f.Mod.Mprp = append(f.Mod.Mprp, []int{i, i + 1, 2, mtyp, brel})
		}
	}
	if len(f.Clvrs) >= 1 && f.Clvrs[0][0] > 0.0{
		f.Lclvr = true
		clen := f.Clvrs[0][0]
		strt := 1
		if f.Vec[0][1] == 1{
			//plinth on
			strt = 2
		}
		for i, y := range f.Y[strt:]{
			f.Mod.Coords = append(f.Mod.Coords, []float64{-clen, y})
			jb := (strt+i) * xstep + 1; je := len(f.Mod.Coords)
			mdx++; brel := 0
			xdx, fdx, locdx, lex, mtyp := getmemvec(f.Ngrp,3,jb,je,xstep)
			f.Members[mdx] = append(f.Members[mdx],[]int{jb, je, 2, mtyp, brel})
			f.Members[mdx] = append(f.Members[mdx],[]int{xdx, fdx, locdx, lex, mtyp})
			f.Mod.Mprp = append(f.Mod.Mprp, []int{jb, je, 2, mtyp, brel})			
		}
	}
	if len(f.Clvrs) > 2 && f.Clvrs[1][0] > 0.0{
		f.Rclvr = true
		clen := f.Clvrs[1][0]
		xb := f.X[len(f.X)-1]+clen
		strt := 1
		if f.Vec[0][1] == 1{
			//plinth on
			strt = 2
		}
		for i, y := range f.Y[strt:]{
			f.Mod.Coords = append(f.Mod.Coords, []float64{xb, y})
			jb := (strt + 1 + i) * xstep; je := len(f.Mod.Coords)
			mdx++; brel := 0
			xdx, fdx, locdx, lex, mtyp := getmemvec(f.Ngrp,3,jb,je,xstep)
			f.Members[mdx] = append(f.Members[mdx],[]int{jb, je, 2, mtyp, brel})
			f.Members[mdx] = append(f.Members[mdx],[]int{xdx, fdx, locdx, lex, mtyp})
			f.Mod.Mprp = append(f.Mod.Mprp, []int{jb, je, 2, mtyp, brel})			
		}
	}
	return
}
