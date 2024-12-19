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
	"github.com/AlecAivazis/survey/v2"
	kass "barf/kass"
)

//Pt is the 50001th (sigh) 3d point struct floating around 
type Pt struct{
	X, Y, Z float64
}

//SubFrm stores struct fields for an rcc subframe
//check hulse section 2.4 and mosley/spencer section 7.1
type SubFrm struct {
	Title        string
	Foldr        string
	Report       string
	Id           int
	Nspans       int
	Nbays        int
	Code         int
	Opt          int
	Ax           int
	Exp          int //exposure condition (for cover params - mild/moderate/extreme)
	Lspans       []float64
	Lbays        []float64
	Lbay         float64
	Lx, Ly       float64
	Hs           []float64
	Hc, Dh       float64
	Drx, Dry     float64
	Drd          float64
	DL           float64
	LL           float64
	Sections     [][]float64 //section dims
	Styps        []int //section types
	Csec, Bsec   []int
	Fcks         []float64
	Fys          []float64
	Fyds         []float64
	DM           float64
	Pufac        float64
	Width        float64
	Dslb         float64
	Pg           float64
	PSFs         []float64
	Clvrs        [][]float64
	Cldim        [][]float64
	Clsec        []int
	Term         string
	Web          bool
	Noprnt       bool
	Nosec        bool
	Fltslb       bool
	Bfcalc       bool
	Dconst       bool
	Lclvr, Rclvr bool
	Tweak        bool
	Tweakb       bool
	Tweakc       bool
	Endrel       bool //if end support is a wall 
	Braced       bool
	Fixbase      bool
	Cred         bool //col live load reduction
	Verbose      bool `json:",omitempty"`
	Ljbase       bool //only if tis a ground (0) floor subframe
	Ujroof       bool //only if tis a nflrs-1 floor subframe
	Spam         bool
	Dz           bool
	Rekt         bool//is rekt sect
	Slbdz        bool
	Lcmem, Rcmem int
	Fop          int //optimize if > 0
	Ngrp         int //number of col/beam groups (uniform for now?)
	Nlp          int
	Selfwt       int
	Ldcalc       int
	Slbload      int
	Bmrel        int
	Nomcvr       float64
	Efcvr        float64
	Edgebm       []float64 `json:",omitempty"` //for flat slabs
	Cdim, Bdim   [][]float64 `json:",omitempty"`
	Tyb          float64 `json:",omitempty"`
	Kost         float64 `json:",omitempty"`
	Secmap       map[int][]float64 `json:",omitempty"`
	Mod          kass.Model `json:",omitempty"`
	X,Y,Z        []float64 `json:",omitempty"`
	Vec          [][]int `json:",omitempty"`
	Val          [][]float64 `json:",omitempty"`
	Flridx       []int `json:",omitempty"`
	Txtplots     []string `json:",omitempty"`
	DMs          []float64 `json:",omitempty"`
	Bloads       [][]float64 `json:",omitempty"`
	Cloads       [][]float64 `json:",omitempty"`
	Jloads       [][]float64 `json:",omitempty"`
	Nodemap      map[int][]int `json:",omitempty"`
	Nodes        map[Pt][]int `json:"-"`
	Members      map[int][][]int `json:",omitempty"`
	Cols, Beams  []int `json:",omitempty"`
	Mloadmap     map[int][][]float64 `json:",omitempty"`
	Jloadmap     map[int][][]float64 `json:",omitempty"`
	Advloads     map[int][][]float64 `json:",omitempty"`
	Benloads     map[int][][]float64 `json:",omitempty"`
	Loadcons     map[int][][]float64 `json:",omitempty"`
	Uniloads     map[int][][]float64 `json:",omitempty"`
	Bmenv        map[int]*kass.BmEnv `json:",omitempty"`
	Colenv       map[int]*kass.ColEnv `json:",omitempty"`
	Secprop      []kass.Secprop `json:",omitempty"`
	Sectin       []kass.SectIn `json:",omitempty"`
	Ms           map[int]*kass.Mem `json:",omitempty"`
	Mslmap       map[int]map[int][][]float64 `json:",omitempty"`
	RcBm         map[int][]*RccBm `json:",omitempty"`
	RcCol        map[int]*RccCol `json:",omitempty"`
	Pus          []float64 `json:",omitempty"`
	Mys          []float64 `json:",omitempty"`
	Nflr         float64 `json:",omitempty"` //use this as axial load multiplier?
	Kostin       []float64 `json:",omitempty"`
	Nmax         []float64 `json:",omitempty"`
	Dtyp         int //design ybays
	Rez          map[int][][]float64 `json:",omitempty"`
	Quants       []float64 `json:",omitempty"` 
}

//SubFrmRez was meant to store sub frame results
//but maybe this is not needed at all
type SubFrmRez struct {
	Bmenvs   []kass.BmEnv
	Colenvs  []kass.ColEnv
	Bmvec    []int
	Colvec   []int
	Flridx   int
	Txtplots []string
}

//Dump saves a SubFrm to a .json file
func (sf *SubFrm) Dump(name string) (filename string, err error){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	if name == "" {name = fmt.Sprintf("subframe_%v",sf.Id)}
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
	data, e := json.Marshal(sf)
	if e != nil{err = e; return}
	err = ioutil.WriteFile(filename, data, 0644)
	return

}

//AddSection adds a section to a SubFrm
func (sf *SubFrm) AddSection(styp int, sdim []float64) (sid int){
	sf.Sections = append(sf.Sections, sdim)
	sf.Styps = append(sf.Styps, styp)
	return len(sf.Sections)
}

//AddNode adds a node to a SubFrm at coords x, y
func (sf *SubFrm) AddNode(x,y float64) (err error){
	p := Pt{X:x,Y:y}
	if _, ok := sf.Nodes[p]; ok{
		return errors.New("node exists")
	} else {
		sf.Mod.Coords = append(sf.Mod.Coords, []float64{x,y})
		sf.Nodes[p] = []int{len(sf.Mod.Coords)}
	}
	return
}

//Init initializes a SubFrm struct
func (sf *SubFrm) Init() (err error){
	//folderr
	if sf.Title == ""{
		sf.Title = fmt.Sprintf("rcc_subframe_%v",rand.Intn(666))
	}
	// if sf.Foldr == ""{sf.Foldr = "out"}
	// e, fdir := kass.InitFolder(sf.Title, sf.Foldr)
	// var kiter int
	// for e != nil{
	// 	if kiter > 10{
	// 		err = errors.New("error creating folder; weird stuff")
	// 		return
	// 	}
	// 	sf.Title = fmt.Sprintf("%s_%v",sf.Title,rand.Intn(666))
	// 	e, fdir = kass.InitFolder(sf.Title, sf.Foldr)
	// 	kiter++
	// }
	// sf.Foldr = fdir
	sf.Mod = kass.Model{}
	if len(sf.Lspans) != sf.Nspans{
		sf.Nspans = len(sf.Lspans)
	}
	
	if len(sf.Lspans) == 0{
		return errors.New("no spans specified")
	}
	if sf.Clvrs == nil{
		sf.Clvrs = [][]float64{{0,0,0},{0,0,0}}
	} 
	if len(sf.Lspans) == 1 && sf.Nspans > 1{
		lspan := sf.Lspans[0]
		sf.Lspans = []float64{}
		for i := 0; i < sf.Nspans; i++{
			sf.Lspans = append(sf.Lspans, lspan)
		}
	}
	//if sf.Nspans == 0 {sf.Nspans = len(sf.Lspans)}
	if len(sf.Hs) == 0{
		return errors.New("floor-floor height not specified")
	}
	if len(sf.Clvrs) > 0 && len(sf.Clvrs[0]) > 0 && sf.Clvrs[0][0]  > 0{
		sf.Lclvr = true
	}
	if len(sf.Clvrs) > 1 && len(sf.Clvrs[1]) > 0 && sf.Clvrs[1][0]  > 0{
		sf.Rclvr = true
	}
	//is code, bs code, (lol) euro code 
	if sf.Code == 0 {sf.Code = 1}
	if sf.Pg == 0.0{
		switch sf.Code{
			case 1:
			sf.Pg = 25.0
			case 2:
			sf.Pg = 24.0
		}
	}
	if sf.PSFs == nil || len(sf.PSFs) < 4{
		switch sf.Code{
			case 1:
			sf.PSFs = []float64{1.5,1.0,1.5,0.0}
			case 2:
			sf.PSFs = []float64{1.4,1.0,1.6,0.0}
		}
	}
	if sf.Slbload > 0{
		if sf.Lbay == 0.0 && len(sf.Lbays)==0{
			return fmt.Errorf("length of bay lbay (%.f) - lbays (%v)required for slab load calc - %v",sf.Lbay, sf.Lbays, sf.Slbload)
		} else {
			sf.Lbay = sf.Lbays[0]
		}
	}
	if sf.Nosec{
		return
	}
	switch len(sf.Sections){
		case 1:
		sf.Csec = []int{1}
		sf.Bsec = []int{1}
		case 2:
		sf.Csec = []int{1}
		sf.Bsec = []int{2}
		
	}
	if (len(sf.Sections) == 0 || len(sf.Csec) == 0 || len(sf.Bsec) == 0){
		return errors.New("no sections specified")
	}
	sf.Nodemap = make(map[int][]int)
	sf.Nodes = make(map[Pt][]int)
	sf.Members = make(map[int][][]int)
	sf.Advloads = make(map[int][][]float64)
	sf.Benloads = make(map[int][][]float64)
	sf.Loadcons = make(map[int][][]float64)
	sf.Uniloads = make(map[int][][]float64)
	sf.RcBm = make(map[int][]*RccBm)
	sf.RcCol = make(map[int]*RccCol)
	//default styp if not specified is 1
	if sf.Styps == nil{
		sf.Styps = make([]int, len(sf.Sections))
		for i := range sf.Styps{
			sf.Styps[i] = 1
		}
	}
	sf.Beams, sf.Cols = []int{},[]int{}
	//csec,bsec default entry
	//if one sec, col/beam; if two secs- col, beam
	if sf.Bfcalc{
		var mrel int
		if sf.Nspans == 1 && (!sf.Lclvr && !sf.Rclvr){
			mrel = 1
		}
		if sf.Dslb == 0{
			log.Println("ERRORE,errore-> depth of slab required for flange calc")
			err = ErrDim
			return
		} else {
			secvec, sts := kass.CalcBf(sf.Code, 0, mrel, sf.Nspans, sf.Dslb, sf.Csec, sf.Bsec, sf.Styps, sf.Lspans, sf.Sections)
			sf.Sections = make([][]float64, len(secvec))
			sf.Styps = make([]int, len(sts))
			for i := range secvec{
				sf.Sections[i] = make([]float64, len(secvec[i]))
				copy(sf.Sections[i],secvec[i])
				sf.Styps[i] = sts[i]
			}
			sf.Csec = make([]int, sf.Nspans + 1)
			for i := range sf.Csec{
				sf.Csec[i] = i+1
			}
			sf.Bsec = make([]int, sf.Nspans)
			for i := range sf.Bsec{
				sf.Bsec[i] = sf.Nspans + 1 + i + 1
			}
		}
	}
	//log.Println("bfcalced\n", "sections\n", sf.Sections,"\nbsec\n", sf.Bsec,"\ncsec\n",sf.Csec,"\nstyps\n",sf.Styps)
	if sf.Pg == 0.0{sf.Pg = 25.0}
	return
}

//GenCoords generates SubFrm coords based on sf.Lspans and sf.Hs 
func (sf *SubFrm) GenCoords(){
	var x float64
	sf.X = append(sf.X, x)
	for _, lspan := range sf.Lspans{
		x += lspan
		sf.X = append(sf.X, x)
	}
	switch len(sf.Hs){
		case 1:
		///THIS IS NOT WERK
		y1 := sf.Hs[0]
		sf.Y = []float64{0,y1}
		case 2:
		y1 := sf.Hs[0]; y2 := sf.Hs[0] + sf.Hs[1]
		sf.Y = []float64{0,y1,y2}
	}
	for _, y := range sf.Y {
		for _, x := range sf.X {
			_ = sf.AddNode(x,y)
		}
	}
}

//AddMem adds a SubFrm member 
func (sf *SubFrm) AddMem(colstep, jb, je, em, cp, mrel, mtyp int) (err error){
	xdx := (jb -1) % colstep; fdx := (jb -1)/ colstep	
	if jb > len(sf.Mod.Coords) || je > len(sf.Mod.Coords){
		return errors.New("invalid member node")
	}
	mdx := len(sf.Mod.Mprp) + 1
	sf.Mod.Mprp = append(sf.Mod.Mprp,[]int{jb, je, em, cp, mrel})
	sf.Secmap[mdx] = sf.Sections[cp-1]
	styp := sf.Styps[cp-1]
	//sf.Secmap[mdx] = append(sf.Secmap[mdx],float64(sf.Styps[cp-1]))
	sf.Nodemap[jb] = append(sf.Nodemap[jb],mdx)
	sf.Nodemap[je] = append(sf.Nodemap[je],mdx)
	lex := 1; locx := 1
	switch mtyp{
		case 1:
		sf.Cols = append(sf.Cols, mdx)
		switch xdx{
			case 0:
			lex = 2; locx = 2
			case colstep:
			lex = 2; locx = 3
		}
		
		case 2,3,4:
		switch xdx{
			case 0:
			lex = 2; locx = 2
			case colstep-1:
			lex = 2; locx = 3
		}
		switch mtyp{
			case 2,4:
			sf.Beams = append(sf.Beams, mdx)
			case 3:
			//start beams with left cantilever
			sf.Beams = append([]int{mdx}, sf.Beams...)
		}
	}
	sf.Members[mdx] = append(sf.Members[mdx],[]int{jb, je, em, cp, mrel})
	sf.Members[mdx] = append(sf.Members[mdx],[]int{mtyp, xdx, fdx, locx, lex, styp})
	return
}

//GenMprp generates column and beam mprp slices of a SubFrm
func (sf *SubFrm) GenMprp() (err error){
	colstep := sf.Nspans + 1
	var colidx, flridx, csec, bsec, cedx, bedx int
	if len(sf.Fcks) == 1 {
		cedx = 1; bedx = 1
		sf.Mod.Em = append(sf.Mod.Em, []float64{FckEm(sf.Fcks[0])})
	} else {
		cedx = 1 ; bedx = 2
		sf.Mod.Em = append(sf.Mod.Em, []float64{FckEm(sf.Fcks[0])})
		sf.Mod.Em = append(sf.Mod.Em, []float64{FckEm(sf.Fcks[1])})
	}
	sf.Secmap = make(map[int][]float64)
	for i:=1; i <=len(sf.Mod.Coords); i++ {
		colidx = (i -1) % colstep; flridx = (i -1)/ colstep	
		if len(sf.Csec) == 1 {
			csec = sf.Csec[0]
		} else {
			csec = sf.Csec[colidx]
		}
		switch flridx {
		case 0:
			sf.Mod.Supports = append(sf.Mod.Supports,[]int{i, -1,-1,-1})
			_ = sf.AddMem(colstep, i, i+colstep, cedx, csec, 0, 1)
		case 1:
			_ = sf.AddMem(colstep, i, i+colstep, cedx, csec, 0, 1)
		case 2:
			sf.Mod.Supports = append(sf.Mod.Supports, []int{i, -1,-1,-1})
		}
	}
	//fmt.Println("colstep",colstep)
	//add beams to mprp
	for i := colstep + 1; i < 2*colstep; i++ {
		colidx = (i -1) % colstep
		//fmt.Println("colidx-",colidx)
		if len(sf.Bsec) == 1 {
			bsec = sf.Bsec[0]
		} else {
			bsec = sf.Bsec[colidx]
		}
		_ = sf.AddMem(colstep, i, i+1, bedx, bsec, sf.Bmrel, 2)
		//if sf.Endrel{
		//	switch i{
		//		case colstep + 1:
		//		//hinge at beginning
		//		_ = sf.AddMem(colstep, i, i+1, bedx, bsec, 1, 2)
		//		case 2 * colstep - 1:
		//		//hinge at end
		//		_ = sf.AddMem(colstep, i, i+1, bedx, bsec, 2, 2)
		//		default:
		//		_ = sf.AddMem(colstep, i, i+1, bedx, bsec, 0, 2)
		//	}
		//} else {
		//}
	}
	//left clvr
	if sf.Clvrs != nil && sf.Clvrs[0][0]  > 0{
		sf.Lclvr = true
		clen := sf.Clvrs[0][0]
		clsec := sf.Bsec[0]
		if len(sf.Clsec) > 0 {clsec = sf.Clsec[0]}
		_ = sf.AddNode(-clen, sf.Y[1])
		je := colstep + 1; jb := len(sf.Mod.Coords)
		_ = sf.AddMem(colstep, jb, je, bedx, clsec, 0, 3)
		sf.Lcmem = len(sf.Mod.Mprp)
	}
	//right clvr
	if sf.Clvrs != nil && sf.Clvrs[1][0]  > 0{
		//right cantilever on
		sf.Rclvr = true
		clen := sf.Clvrs[1][0]
		clsec := sf.Bsec[0]
		switch len(sf.Clsec){
			case 1:
			clsec = sf.Clsec[0]
			case 2:
			clsec = sf.Clsec[1]
			
		}
		if len(sf.Clsec) > 0 {clsec = sf.Clsec[0]}
		_ = sf.AddNode(clen+sf.X[len(sf.X)-1], sf.Y[1])
		jb := colstep * 2; je := len(sf.Mod.Coords)
		_ = sf.AddMem(colstep, jb, je, bedx, clsec, 0, 4)
		sf.Rcmem = len(sf.Mod.Mprp)
	}
	//get cp [][]
	for idx, sect := range sf.Sections {
		styp := 1
		if len(sf.Styps) == len(sf.Sections){
			styp = sf.Styps[idx]
		}
		//bar := kass.CalcSecProp(styp, sect)
		sectin := kass.SecGen(styp, sect)
		bar := sectin.Prop
		sf.Secprop = append(sf.Secprop, bar)
		sf.Mod.Cp = append(sf.Mod.Cp, []float64{bar.Area/1e6, bar.Ixx/1e12})
		sf.Sectin = append(sf.Sectin, sectin)
	}
	//check for self weight calc
	switch sf.Selfwt{
	//load vec - [self wt, nload cases]
		case 1,2:
		err := sf.AddSelfWeight(sf.Selfwt)
		if err != nil{
			log.Println("ERRORE, errore->",err)
		}
	}
	return
}

//AddSelfWeight adds self weight of subframe beams (dx = 1) and columns (dx = 2) as member loads
func (sf *SubFrm) AddSelfWeight(dx int) (err error){
	//REDO
	//density of concrete 25 KN/m3 - HAHA no.
	var wck float64
	switch sf.Code{
		case 1:
		wck = 25.0
		case 2:
		wck = 24.0
	}
	//if sf.Verbose{fmt.Println(ColorYellow,"adding self weight -> beams",ColorReset)}
	for _, bm := range sf.Beams{
		bdx := sf.Members[bm][0][3]
		bstyp := sf.Styps[bdx-1]
		var wdl float64
		switch bstyp{
			case 1:
			b := sf.Sections[bdx-1][0]; d := sf.Sections[bdx-1][1]
			if sf.Fltslb{
				wdl = wck * b * d/1e6
				//fmt.Println("wts->",wdl,wck*b*d/1e6)
			} else {
				wdl = wck * b * (d - sf.Dslb)/1e6
			}
			case 6,7,8,9,10,14:
			bdim := sf.Sections[bdx-1]; bw := bdim[2]; dw:= bdim[1] - sf.Dslb
			wdl = wck * bw * dw /1e6			
			default:
			wdl = wck * sf.Secprop[bdx-1].Area/1e6
		}
		//if sf.Verbose{fmt.Println(ColorCyan,"bm id ->",bm,"self wt->",wdl,"kn/m",ColorReset)}
		ldcase := []float64{1.0, 3.0, wdl, 0.0, 0.0, 0.0, 1.0}
		err = sf.AddMemLoad(bm, ldcase)
		if err != nil{
			fmt.Println("ERRORE,errore->",err)
			return
		}
	}
	if dx == 2{
		//if sf.Verbose{fmt.Println("adding self weight -> columns")}
		for _, col := range sf.Cols{
			bdx := sf.Members[col][0][3]
			wdl := wck * sf.Secprop[bdx-1].Area/1e6
			ldcase := []float64{1.0, 6.0, wdl, 0.0, 0.0, 0.0, 1.0}
			err = sf.AddMemLoad(col, ldcase)
			if sf.Verbose{fmt.Println(ColorCyan,"col id ->",col,"self wt->",wdl,"kn/m")}
			if err != nil{
				fmt.Println("ERRORE,errore->",err)
				return
			}
		}
	}
	return
}

//AddMemLoad adds a member load (ldcase) to a member (mem) of a SubFrm
func (sf *SubFrm) AddMemLoad(mem int, ldcase []float64) (err error){
	if _, ok := sf.Members[mem]; !ok{
		return errors.New("invalid member index")
	}
	ldcat := ldcase[6]
	var w1a, w2a, w1b, w2b float64
	switch ldcat{
		case 1.0:
		w1a = sf.PSFs[0]*ldcase[2]
		w2a = sf.PSFs[0]*ldcase[3]
		w1b = sf.PSFs[1]*ldcase[2]
		w2b = sf.PSFs[1]*ldcase[3]
		//wot horrible way to append (uniload = []float64{mem}; uniload = append(uniloads, ldcase[1:]...)
		sf.Uniloads[1] = append(sf.Uniloads[1], []float64{float64(mem),ldcase[1],ldcase[2],ldcase[3],ldcase[4],ldcase[5],ldcase[6]})
		case 2.0:
		w1a = sf.PSFs[2]*ldcase[2]
		w2a = sf.PSFs[2]*ldcase[3]
		w1b = sf.PSFs[3]*ldcase[2]
		w2b = sf.PSFs[3]*ldcase[3]
		sf.Uniloads[2] = append(sf.Uniloads[2], []float64{float64(mem),ldcase[1],ldcase[2],ldcase[3],ldcase[4],ldcase[5],ldcase[6]})
	}
	sf.Advloads[mem] = append(sf.Advloads[mem],[]float64{float64(mem),ldcase[1],w1a,w2a,ldcase[4],ldcase[5], ldcat})
	if w1b + w2b > 0.0 {
		sf.Benloads[mem] = append(sf.Benloads[mem],[]float64{float64(mem),ldcase[1],w1b,w2b,ldcase[4],ldcase[5], ldcat})
	}
	return
}

//AddBeamLLDL adds sf.DL and sf.LL as a udl to all SubFrm beams
func (sf *SubFrm) AddBeamLLDL(){
	if sf.DL + sf.LL > 0.0{
		for _, bm := range sf.Beams{
			mvec := sf.Members[bm]
			jb := sf.Mod.Coords[mvec[0][0]-1]; je := sf.Mod.Coords[mvec[0][1]-1]
			lspan := kass.Dist2d(jb, je)
			switch bm{
				case sf.Lcmem:
				cdl := sf.DL; cll := sf.LL 
				if sf.Clvrs[0][1]+sf.Clvrs[0][2] > 0.0{
					cdl = sf.Clvrs[0][1]; cll = sf.Clvrs[0][2]
				}
				_ = sf.AddMemLoad(bm, []float64{1.0,3.0,cdl,0.0,0.0,0.0,1.0})
				_ = sf.AddMemLoad(bm, []float64{1.0,3.0,cll,0.0,0.0,0.0,2.0})
				case sf.Rcmem:
				cdl := sf.DL; cll := sf.LL 
				if sf.Clvrs[1][1]+sf.Clvrs[1][2] > 0.0{
					cdl = sf.Clvrs[1][1]; cll = sf.Clvrs[1][2]
				}
				_ = sf.AddMemLoad(bm, []float64{1.0,3.0,cdl,0.0,0.0,0.0,1.0})
				_ = sf.AddMemLoad(bm, []float64{1.0,3.0,cll,0.0,0.0,0.0,2.0})
				default:
				if sf.Fltslb{
					_ = sf.AddMemLoad(bm, []float64{1.0,3.0,sf.Lbay * sf.DL,0.0,0.0,0.0,1.0})
					_ = sf.AddMemLoad(bm, []float64{1.0,3.0,sf.Lbay * sf.LL,0.0,0.0,0.0,2.0})
					
				} else {
					switch sf.Slbload{
						case 0:
						if sf.DL > 0.0{
							_ = sf.AddMemLoad(bm, []float64{1.0,3.0,sf.DL,0.0,0.0,0.0,1.0})
						}
						if sf.LL > 0.0{
							_ = sf.AddMemLoad(bm, []float64{1.0,3.0,sf.LL,0.0,0.0,0.0,2.0})
						}
						case 1, 3:
						//one way slab load - over lbay/2.0 * 2.0
						var dl float64
						if sf.Dslb > 0.0{
							switch sf.Slbload{
								case 1:
								dl += sf.Pg * sf.Dslb * sf.Lbay * 1e-3
								//fmt.Println("adding one way slab load->",sf.Pg * sf.Dslb * sf.Lbay * 1e-3)
								case 3:
								//this is already added by miracle of chaining
							}
						}
						dl += sf.DL * sf.Lbay; ll := sf.LL * sf.Lbay
						//fmt.Println("dead and live load on frame-",dl, ll)
						if dl > 0.0{
							_ = sf.AddMemLoad(bm, []float64{1.0,3.0,dl,0.0,0.0,0.0,1.0})
						}
						if ll > 0.0{
							_ = sf.AddMemLoad(bm, []float64{1.0,3.0,ll,0.0,0.0,0.0,2.0})
						}
						case 2:
						//two way slab load - 3 tri right, udl, tri left
						var dl float64
						if sf.Dslb > 0.0{
							dl += sf.Pg * sf.Dslb * 1e-3
						}
						dl += sf.DL; ll := sf.LL
						switch{
							case sf.Lbay >= lspan:
							//triangular load of 2.0 * wlspan/4.0
							dl = dl * lspan
							ll = ll * lspan
							_ = sf.AddMemLoad(bm, []float64{1.0,4.0,0.0,dl,0.0,lspan/2.0,1.0})
							_ = sf.AddMemLoad(bm, []float64{1.0,4.0,dl,0.0,lspan/2.0,0.0,1.0})
							if ll > 0.0{
								_ = sf.AddMemLoad(bm, []float64{1.0,4.0,0.0,ll,0.0,lspan/2.0,2.0})
								_ = sf.AddMemLoad(bm, []float64{1.0,4.0,ll,0.0,lspan/2.0,0.0,2.0})
							}					
							case sf.Lbay < lspan:
							//
							//trap load of peak 2.0 * wlbay2/2.0
							dl = dl * sf.Lbay; ll = ll * sf.Lbay
							if dl > 0.0{
								_ = sf.AddMemLoad(bm, []float64{1.0,4.0,0.0,dl,0.0,lspan-sf.Lbay/2.0,1.0})
								_ = sf.AddMemLoad(bm, []float64{1.0,3.0,dl,0.0,sf.Lbay/2.0,sf.Lbay/2.0,1.0})
								_ = sf.AddMemLoad(bm, []float64{1.0,4.0,dl,0.0,lspan-sf.Lbay/2.0,0.0,1.0})
							}
							if ll > 0.0{
								_ = sf.AddMemLoad(bm, []float64{1.0,4.0,0.0,ll,0.0,lspan-sf.Lbay/2.0,2.0})
								_ = sf.AddMemLoad(bm, []float64{1.0,3.0,ll,0.0,sf.Lbay/2.0,sf.Lbay/2.0,2.0})
								_ = sf.AddMemLoad(bm, []float64{1.0,4.0,ll,0.0,lspan-sf.Lbay/2.0,0.0,2.0})
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
	}
}

//GenLoads generates SubFrm loads/load patterns
func (sf *SubFrm) GenLoads() (err error){
	sf.AddBeamLLDL()
	colstep := sf.Nspans + 1
	if len(sf.Bloads) != 0 {
		for _, ldcase := range sf.Bloads{
			mem := int(ldcase[0]) + len(sf.Hs)*colstep
			_ = sf.AddMemLoad(mem, ldcase)
		}
	}
	if len(sf.Cloads) != 0 {
		for _, ldcase := range sf.Cloads{
			mem := int(ldcase[0]) 
			_ = sf.AddMemLoad(mem, ldcase)
		}
	}
	sf.Nlp = sf.Nspans + 1
	//build load patterns
	//IS THERE NO WAY TO TURN IT OFF
	for i := 1; i <= sf.Nspans; i++ {
		mem := i + len(sf.Hs)*colstep
		
		if i % 2 == 0 {
			sf.Loadcons[0] = append(sf.Loadcons[0], sf.Advloads[mem]...)
			sf.Loadcons[1] = append(sf.Loadcons[1], sf.Benloads[mem]...)
			sf.Loadcons[2] = append(sf.Loadcons[2], sf.Advloads[mem]...)
		} else {
			sf.Loadcons[0] = append(sf.Loadcons[0], sf.Advloads[mem]...)
			sf.Loadcons[1] = append(sf.Loadcons[1], sf.Advloads[mem]...)
			sf.Loadcons[2] = append(sf.Loadcons[2], sf.Benloads[mem]...)
		}
	}
	for i := 1; i <= sf.Nspans - 1; i++ {
		lp := i + 2
		mind := i + len(sf.Hs)*colstep
		for mem := range sf.Advloads {
			if mem == mind || mem == mind + 1 {
				sf.Loadcons[lp] = append(sf.Loadcons[lp], sf.Advloads[mem]...)
			} else {
				sf.Loadcons[lp] = append(sf.Loadcons[lp], sf.Benloads[mem]...)
			}
		}
	}
	if sf.Lclvr{
		sf.Nlp++
		for lp := 0; lp <= 2; lp++{
			switch lp{
				case 0,2:
				sf.Loadcons[lp] = append(sf.Loadcons[lp],sf.Advloads[sf.Lcmem]...)
				case 1:
				sf.Loadcons[lp] = append(sf.Loadcons[lp],sf.Benloads[sf.Lcmem]...)
			}
		} 
		for mem := range sf.Benloads {
			switch mem{
				case sf.Lcmem, 1+len(sf.Hs)*colstep:
				sf.Loadcons[sf.Nspans+2] = append(sf.Loadcons[sf.Nspans+2], sf.Advloads[mem]...)
				default:
				sf.Loadcons[sf.Nspans+2] = append(sf.Loadcons[sf.Nspans+2], sf.Benloads[mem]...)
			}
		}
	}
	if sf.Rclvr{
		sf.Nlp++
		for lp := 0; lp <= 2; lp++{
			switch lp{
				case 0,1:
				sf.Loadcons[lp] = append(sf.Loadcons[lp],sf.Advloads[sf.Rcmem]...)
				case 2:
				sf.Loadcons[lp] = append(sf.Loadcons[lp],sf.Benloads[sf.Rcmem]...)
			}
		} 
		for mem := range sf.Benloads {
			rmem := len(sf.Hs)*colstep + sf.Nspans
			switch mem{
				case sf.Rcmem, rmem:
				sf.Loadcons[sf.Nspans+3] = append(sf.Loadcons[sf.Nspans+3], sf.Advloads[mem]...)
				default:
				sf.Loadcons[sf.Nspans+3] = append(sf.Loadcons[sf.Nspans+3], sf.Benloads[mem]...)
			}	
		}
	}
	return
}

//InitMemRez initializes sf.Bmenv and sf.Colenv for beam and column results
func (sf *SubFrm) InitMemRez() (err error){
	//return error for non standard sections
	colstep := sf.Nspans + 1
	sf.Bmenv = make(map[int]*kass.BmEnv)
	bmendc := 2
	//NOTE - ADD THIS AND A FEW OTHER FUNCS TO FRM2D; if bmrel = 3 endc = 1
	
	for _, i := range sf.Beams{
		bmendc = 2
		//ldxb := 2; rdxb := 2
		//could be a portal frame so nope
		//if sf.Nspans == 1 && !sf.Lclvr && !sf.Rclvr{bmendc = 1}
		if sf.Bmrel == 3 {
			bmendc = 1
		}
		jb, je := sf.Mod.Mprp[i-1][0], sf.Mod.Mprp[i-1][1]
		c1, c2 := sf.Mod.Coords[jb-1], sf.Mod.Coords[je-1]
		bmlen:= EDist(c1, c2) 
		cl := i - len(sf.Hs) * colstep
		cr := cl + 1
		ldx := 2
		rdx := 2
		//fmt.Println("beam->",i, cl, cr)
		var lsx, rsx float64
		em := sf.Mod.Mprp[i-1][2]; cp := sf.Mod.Mprp[i-1][3]
		xdx := sf.Members[i][1][1]
		kbm := sf.Secprop[cp-1].Ixx/bmlen/1e3
		//fmt.Println(ColorRed,"kbm->",kbm,ColorReset)
		switch i{
			case sf.Lcmem:
			rsx = sf.Secmap[1][1]
			bmendc = 0
			case sf.Rcmem:
			lsx = sf.Secmap[colstep][1]
			bmendc = 0
			default:
			lsx = sf.Secmap[cl][1]; rsx = sf.Secmap[cr][1]
		}
		//fmt.Println("xstep->",xstep,"nspans-",f.Nspans,"floors-",f.Nflrs)
		//fmt.Println("beam->",i,"xdx-",xdx)
		spandx := 0
		switch xdx{
		case 0:
			spandx = -1
			ldx = 1
		case sf.Nspans - 1:
			spandx = -2
			rdx = 1
		}
		//fmt.Println(ColorYellow,"init mem rez beam->",i,lsx, rsx,ColorReset)
		sf.Bmenv[i] = &kass.BmEnv{
			Id:i,
			EnvRez:make(map[int]kass.BeamRez),
			Venv:make([]float64,21),
			Mpenv:make([]float64,21),
			Mnenv:make([]float64,21),
			Dims:sf.Secmap[i],
			Coords:[][]float64{c1,c2},
			Lsx:lsx/1000.0,Rsx:rsx/1000.0,
			Lspan:bmlen,
			Endc:bmendc,
			Vrd:make([]float64,21),
			Mprd:make([]float64,21),
			Mnrd:make([]float64,21),
			Hc:sf.Hc/1000.0,
			Dh:sf.Dh/1000.0,
			Drx:sf.Drx/1000.0,
			Dry:sf.Dry/1000.0,
			Drd:sf.Drd/1000.0,
			Spandx:spandx,
			Endrel:sf.Endrel,
			Foldr:sf.Foldr,
			Term:sf.Term,
			K:kbm,
			Kostin:sf.Kostin,
			Ldx:ldx,
			Rdx:rdx,
		}
		if sf.Fltslb{
			sf.Bmenv[i].Wspan = sf.Lbay
		}
		sf.RcBm[i] = make([]*RccBm,3)
		if sf.Fltslb{sf.RcBm[i] = make([]*RccBm,6)}
		bstyp := sf.Styps[cp-1]
		var bf, df, bw, dused, tyb float64
		var npsec bool
		tyb = 0.0
		switch bstyp{
			case 1:
			bw = sf.Secmap[i][0]; dused = sf.Secmap[i][1]
			case 6,7,8,9,10:
			bf = sf.Secmap[i][0]; dused = sf.Secmap[i][1]; bw = sf.Secmap[i][2]; df = sf.Secmap[i][3]
			case 14:
			bf = sf.Secmap[i][0]; dused = sf.Secmap[i][1]; bw = sf.Secmap[i][2]; df = sf.Secmap[i][3]
			default:
			npsec = true
		}
		switch bstyp{
			case 6, 14:
			tyb = 1.0
			case 7:
			tyb = 0.5
		}
		//mumbai beam nominal cover - 30 mm (m25), effcvr 30 + 20/2
		var sti int
		var tybi float64
		nbms := 3
		if sf.Fltslb{nbms = 6}
		for j := 0; j < nbms; j++{
			flip := true
			switch bstyp{
				case 1,6,7,8,9,10:
				sti = 1; tybi = 0.0
				case 14:
				sti = bstyp; tybi = tyb
			}
			//flange only at midspan
			if j == 1 || j == 4{
				sti = bstyp; tybi = tyb
				flip = false
			}
			sf.RcBm[i][j] = &RccBm{
				Id:j,
				Mid:i,
				Fck:sf.Fcks[em-1],
				Fy:sf.Fys[em-1],
				Bf:bf,
				Df:df,
				Bw:bw,
				Dused:dused,
				Styp:sti,
				Tyb:tybi,
				Cvrt:40.0,
				Cvrc:40.0,
				Flip:flip,
				Code:sf.Code,
				Endc:bmendc,
				Dims:sf.Secmap[i],
				Npsec:npsec,
				Verbose:sf.Verbose,
				DM:sf.DM,
				Kostin:sf.Kostin,
				Foldr:sf.Foldr,
				Tweak:sf.Tweak,
				Term:sf.Term,
				Lsx:lsx,
				Rsx:rsx,
				Ldx:ldx,
				Rdx:rdx,
			}
			sf.RcBm[i][j].Init()
		}
	}
	//column nominal cover - 40 mm (m25); effcvr 40 + 20/2
	sf.Colenv = make(map[int]*kass.ColEnv)
	bmstrt := colstep * len(sf.Hs)
	for _, i := range sf.Cols{
		//fmt.Println("xdx,fdx,locdx,lex,mtyp->",sf.Members[i][1])
		
		var ljbase, lbss, ubss, ignore bool
		jb, je := sf.Mod.Mprp[i-1][0], sf.Mod.Mprp[i-1][1]
		c1, c2 := sf.Mod.Coords[jb-1], sf.Mod.Coords[je-1]
		cp := sf.Members[i][0][3]
		lspan := kass.Dist2d(c2,c1)
		kcol := sf.Secprop[cp-1].Ixx/lspan/1e3
		//sf.Members[mdx] = append(sf.Members[mdx],[]int{xdx, fdx, locdx, lex, mtyp})
		xdx, fdx := sf.Members[i][1][1],sf.Members[i][1][2]
		//fmt.Println("col->",i, jb, je)
		//fmt.Println("lspan",lspan,"kcol",kcol)
		var lbd, ubd [][]float64
		var lst, ust []int
		var lbk, ubk []float64
		var dl, dt float64
		var ct, cb, ub1, ub2, lb1, lb2 int
		lb1 = bmstrt + 1 + (fdx - 1) * (colstep - 1) + (xdx - 1)
		lb2 = lb1 + 1
		ub1 = bmstrt + 1 + (fdx) * (colstep - 1) + (xdx - 1)
		ub2 = ub1 + 1
		ct = i + colstep; cb = i - colstep
		switch xdx{
			case 0:
			lb1, ub1 = 0, 0
			case colstep -1:
			lb2, ub2 = 0, 0 
			
		}
		switch fdx{
			case 0:
			lb1, lb2 = 0, 0
			ljbase = sf.Ljbase
			cb = 0
			case 1:
			ct = 0
			ub1, ub2 = 0, 0
			ignore = true
			
		}
		if len(sf.Hs) == 1{
			ct = 0
			ub1, ub2 = 0, 0
		}
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
				dims = sf.Bmenv[bm].Dims
				styp = sf.Styps[cp-1]
			}
			if bm > 0{
				switch styp{
				//note - lspan is in METERS 
					case 1:
					kb = dims[0] * math.Pow(dims[1],3.0)/12.0
					kb = kb/sf.Bmenv[bm].Lspan/1e3
					case 6,7,8,9,10:
					//use rect beam stiffness - ONLY IF RCC (change this)
					kb = dims[2] * math.Pow(dims[1],3.0)/12.0
					kb = kb/sf.Bmenv[bm].Lspan/1e3
					default:
					//use reg stiffness = i/lspan
					cp := sf.Bmenv[bm].Cp
					kb = sf.Secprop[cp-1].Ixx/sf.Bmenv[bm].Lspan/1e3
				}
			}
			bmd = dims[0]; if len(dims)>1 {bmd = dims[1]}
			//fmt.Println("beam->",bm,bmd)
			if i < 2{
				lbd = append(lbd, dims)
				lbk = append(lbk, kb)
				lst = append(lst, styp)
				if dl <= bmd/2.0 {dl = bmd/2.0}
				if bm != 0 && sf.Bmenv[bm].Endc < 2{lbss = true}
			} else {
				ubd = append(ubd, dims)
				ubk = append(ubk, kb)
				ust = append(ust, styp)
				if dt <= bmd/2.0 {dt = bmd/2.0}
				if bm != 0 && sf.Bmenv[bm].Endc < 2{ubss = true}
			}
		}
		//fmt.Println("col->",i,"xdx, fdx-",xdx, fdx, "beams->",ub1,ub2,lb1,lb2,"ct, cb->",ct,cb)
		var ksum, kcsum float64
		kcsum = kcol
		switch sf.Code{
			case 1:
			if sf.Braced{
				ksum = kcol + 0.5 * (lbk[0] + lbk[1])
			} else {
				ksum = kcol + 1.5 * (lbk[0] + lbk[1])
			}	
			default:
			ksum = lbk[0] + lbk[1]
		}
		var b1, b2 float64
		if cb != 0{
			jb1, je1:= sf.Mod.Mprp[cb-1][0], sf.Mod.Mprp[cb-1][1]
			p1, p2 := sf.Mod.Coords[jb1-1], sf.Mod.Coords[je1-1]
			cp1 := sf.Members[cb][0][3]
			kcb := sf.Secprop[cp1-1].Ixx/kass.Dist2d(p1,p2)/1e3
			kcsum += kcb
			if sf.Code == 1 {ksum += kcb}
		}
		if ksum != 0{
			b1 = kcsum/ksum
		}
		switch sf.Code{
			case 1:
			if sf.Braced{
				ksum = kcol + 0.5 * (ubk[0] + ubk[1])
			} else {
				ksum = kcol + 1.5 * (ubk[0] + ubk[1])
			}
			default:
			ksum = ubk[0] + ubk[1]
		}
		kcsum = kcol
		if ct != 0{
			jb1, je1:= sf.Mod.Mprp[ct-1][0], sf.Mod.Mprp[ct-1][1]
			p1, p2 := sf.Mod.Coords[jb1-1], sf.Mod.Coords[je1-1]
			cp1 := sf.Members[ct][0][3]
			kct := sf.Secprop[cp1-1].Ixx/kass.Dist2d(p1,p2)/1e3
			kcsum += kct
			if sf.Code == 1{ksum += kct}
		}
		if ksum != 0{b2 = kcsum/ksum}
		sf.Colenv[i] = &kass.ColEnv{
			Id:i,
			EnvRez:make(map[int][]float64),
			Dims:sf.Secmap[i],
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
			Styp:sf.Styps[cp-1],
			Lspan:lspan,
			Fixbase:sf.Fixbase,
			Braced:sf.Braced,
			Kcol:kcol,
			L0:lspan - dl/1e3 - dt/1e3,
			Lbss:lbss,
			Ubss:ubss,
			B1:b1,
			B2:b2,
			Pufac:sf.Pufac,
			Term:sf.Term,
			Tweak:sf.Tweak,
			Foldr:sf.Foldr,
			Ignore:ignore,
		}
		if !ignore{sf.Colenv[i].EffHt(1)}
	}
	sf.Ms = make(map[int]*kass.Mem)
	sf.Mslmap = make(map[int]map[int][][]float64)
	return
}

//CalcLoadEnv calcs SubFrm load envelopes (analysis via kass.Model)
//updates/sorts individual member results in sf.Bmenv and sf.Colenv
func (sf *SubFrm) CalcLoadEnv(){
	mod := &kass.Model{
		Frmstr:"2df",
		Units:"knm",
		Coords: sf.Mod.Coords,
		Supports: sf.Mod.Supports,
		Em: sf.Mod.Em,       
		Cp: sf.Mod.Cp,       
		Mprp:sf.Mod.Mprp,
		Noprnt:sf.Noprnt,
		Web:sf.Web,
		Term:sf.Term,
	}
	for lp, ldcons := range sf.Loadcons{
		if sf.Ldcalc == 1 && lp > 0{continue}
		//if lp >= 0{
		//	fmt.Println(ColorYellow,lp,"\n",ColorGreen,ldcons,"\n",ColorReset)
		//}
		mod.Msloads = ldcons
		frmrez, err := kass.CalcFrm2d(mod, 3)
		if err != nil{return}
		ms,_ := frmrez[1].(map[int]*kass.Mem)
		msloaded, _ := frmrez[5].(map[int][][]float64)
		spanchn := make(chan kass.BeamRez,len(msloaded))
		sf.Mslmap[lp] = msloaded
		switch lp{
			case 0:
			sf.Mod.Msloads = ldcons
			sf.Ms = ms
		}
		for _, id := range sf.Beams{
			var ldcase [][]float64
			if _, ok := msloaded[id];ok{ldcase = msloaded[id]}
			//ldcase := msloaded[id]
			go kass.BeamFrc(3, id, ms[id], ldcase, spanchn, false)
		}
		
		for _, id := range sf.Cols{
			mem := ms[id]
			cm := sf.Colenv[id]
			cm.EnvRez[lp] = mem.Qf
			if math.Abs(mem.Qf[2]) > math.Abs(cm.Mtmax) {cm.Mtmax = mem.Qf[2]}
			if math.Abs(mem.Qf[5]) > math.Abs(cm.Mbmax) {cm.Mbmax = mem.Qf[5]}
			if math.Abs(mem.Qf[0]) > math.Abs(cm.Pbmax) {cm.Pbmax = mem.Qf[0]}
			if math.Abs(mem.Qf[3]) > math.Abs(cm.Ptmax) {cm.Ptmax = mem.Qf[3]}
			if math.Abs(mem.Qf[0]) > math.Abs(cm.Pumax) {cm.Pumax = math.Abs(mem.Qf[0])}
			if math.Abs(mem.Qf[3]) > math.Abs(cm.Pumax) {cm.Pumax = math.Abs(mem.Qf[3])}
		}
		//for id, ldcase := range msloaded{
		//	go BeamFrc(3, id, ms[id], ldcase, spanchn, false)
		//}
		for range sf.Beams{
			r := <- spanchn
			id := r.Mem
			
			bm := sf.Bmenv[id]
			
			bm.EnvRez[lp] = r
			if len(bm.Xs) == 0 {
				bm.Xs = r.Xs
			}
			xdiv := ms[id].Geoms[0]/20.0

			lsx := bm.Lsx; rsx := ms[id].Geoms[0] - bm.Rsx
			//fmt.Println("lsx, rsx->",lsx,rsx," MEYTARSS")
			il := int(math.Ceil(lsx/xdiv)); ir := int(math.Ceil(rsx/xdiv))
			var vl, vr, ml, mr float64
			for i, vx := range r.SF{
				x := r.Xs[i]
				if i == il{
					//fmt.Println("i, x, il, lsx->",i, x, il, lsx," MEYTARSS")
					switch{
						case x == lsx:
						vl = vx
						ml = r.BM[i]
						case x > lsx && i == 0:
						vl = vx
						ml = r.BM[i]
						default:
						vl = (vx + r.SF[i-1])/2.0
						ml = (r.BM[i] + r.BM[i-1])/2.0
						//vl = vx + (vx - r.SF[i-1])*(lsx - x)/xdiv
						//ml = r.BM[i] + 0.5 * (lsx - x)*(vl + vx)
						//just lerp?
						//ml = r.BM[i] + (r.BM[i] - r.BM[i-1])*(lsx - x)/xdiv
					}
					if math.Abs(bm.Vl) < math.Abs(vl){bm.Vl = vl}
					if math.Abs(bm.Ml) < math.Abs(ml){bm.Ml = ml}
				}
				if i == ir{
					//fmt.Println("i, x, ir, rsx->",i, x, ir, rsx," MEYTARSS")
					switch{
						case x == rsx:
						vr = vx
						mr = r.BM[i]
						
						default:
						
						vr = (vx + r.SF[i-1])/2.0
						mr = (r.BM[i] + r.BM[i-1])/2.0
						//vr = vx + (vx - r.SF[i-1])*(rsx - x)/xdiv
						//mr = r.BM[i] + 0.5 * (rsx - x)*(vr + vx)
						//mr = r.BM[i] + (r.BM[i] - r.BM[i-1])*(rsx - x)/xdiv
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
	
	for i := range sf.Colenv{
		c := GetCol(sf.Colenv[i], sf.Kostin, sf.Fcks[0], sf.Fys[0], sf.Efcvr, 1, sf.Code)
		sf.RcCol[i] = &c
		sf.RcCol[i].Init()
	}
}

//Mrd redistributes beam support moments of a SubFrm
func (sf *SubFrm) Mrd(){
	if sf.DM != 0.0 && sf.Ldcalc != 1{
		//fmt.Println("moment redistribution DM->",sf.DM)
		CBeamDM(3, sf.Beams, sf.Bmenv, sf.DM, sf.Ms, sf.Mslmap)
	}
}

//DrawMod draws the SubFrm line diagram (PlotSubFrm) and plots load patterns/beam envelopes
func (sf *SubFrm) DrawMod(){
	folder := sf.Foldr
	if sf.Web{
		folder = "web"
	}
	pltstr := PlotSubFrm(sf, &sf.Mod, sf.Ms, sf.Beams, sf.Cols, sf.Bmenv, sf.Colenv, sf.Term)
	sf.Txtplots = append(sf.Txtplots, pltstr)
	sf.Txtplots = append(sf.Txtplots, PlotBmEnv(sf.Bmenv, sf.Beams, sf.Term, sf.Title, folder))
	if sf.DM > 0.0{
		sf.Txtplots = append(sf.Txtplots, PlotBmRdEnv(sf.Bmenv, sf.Beams, sf.Term, sf.Title, sf.Foldr))
	}
	plotchn := make(chan []interface{},sf.Nlp)
	if sf.Ldcalc == 1{sf.Nlp = 0}
	for i := 0; i <= sf.Nlp; i++{
		go PlotLp(i, sf.Bmenv, sf.Ms, sf.Mslmap, sf.Beams, sf.Term, sf.Title, folder,plotchn)
	}
	for i := 0; i<=sf.Nlp; i++{
		rez := <- plotchn
		//lp, _ := rez[0].(int)
		txtplot, _ := rez[1].(string)
		sf.Txtplots = append(sf.Txtplots, txtplot)
	}
}

//CalcSubFrm is the entry func for subframe analysis
//calls all the above funcs as seen below
//as above, so b...
func CalcSubFrm(sf *SubFrm) (err error) {
	err = sf.Init()
	if err != nil{
		return
	}
	sf.GenCoords()
	sf.GenMprp()
	sf.GenLoads()
	if sf.Tweak{
		sf.TweakMen()
	}
	sf.InitMemRez()
	sf.CalcLoadEnv()
	sf.Mrd()
	if sf.Term != "" && sf.Verbose{sf.DrawMod()}
	// if sf.Verbose{
	// 	for _, txtplt := range sf.Txtplots{
	// 		fmt.Println(txtplt)
	// 	}
	// }
	return 
}

//Draw3d tries to draw a 3d view of a SubFrm
//it does nothing, just fails
func (sf *SubFrm) Draw3d(){
	//first draw 2d view from top (plan)
}

func (sf *SubFrm) TweakMen(){
	//menu for tweaks
	running := true
	choice := -1
	for running{
		prompt := &survey.Select{
			Message: "choose item to tweak",
			Options: []string{"add/edit member","add member load","add nodal load","exit"},
		}
		survey.AskOne(prompt, &choice)
		switch choice{
			case 3:
			running = false 
			case 0:
			fmt.Println("eedit")
			case 1:
			c1 := -1
			for c1 != 1{	
				p1 := &survey.Select{
					Message: "choose item",
					Options: []string{"enter member load","exit"},
				}
				survey.AskOne(p1, &c1)
				
				switch c1{
					case 1:
					running = false
					case 0:
					var mem, ltyp, wa, wb, la, lb, lc float64
					fmt.Printf("enter member, load type(1-10),wa, wb, la, lb, load condition->\n")
					_, err := fmt.Scanln(&mem, &ltyp, &wa, &wb, &la, &lb, &lc)
					if err != nil{
						fmt.Println(ColorRed, err, ColorReset)
						continue
					}
					if _, ok := sf.Members[int(mem)]; !ok{
						fmt.Println(ColorRed, "member not found",ColorReset)
					} else {
						switch {
						case ltyp < 1 || ltyp > 10:
							fmt.Println(ColorRed,"invalid load type(1-10)",ColorReset)
						default:
							err = sf.AddMemLoad(int(mem), []float64{mem,ltyp,wa,wb,la,lb,lc})
							if err != nil{
								
								fmt.Println(ColorRed, err, ColorReset)
								continue
							}
						}
					}
				}
			}
			case 2:
			//nodal load
		}
	}
}

/*

func SubFrmBmSizes(sf *SubFrm, b, dbrat, cvrt float64) (bmsecs [][]float64) {
	//REDO DIS
	bmsecs = make([][]float64, len(sf.Lspans))
	//wmx := sf.DL * sf.PSFs[0] + sf.LL * sf.PSFs[2]
	//kulim := 805.0/(1265.0 + sf.Fys[0])
	for i, lspan := range sf.Lspans{
		//mud := wmx * math.Pow(lspan,2)/8.0
		effd := 1000.0*lspan/10.0
		d := effd + cvrt
		switch {
		case b == 0 && dbrat == 0:
			//set default dbrat and check for bd
			dbrat = 3.0
			b = effd/dbrat
		case b == 0:
			b = effd/dbrat
		}
		bmsecs[i] = []float64{b,d}
	}
	return
}

func SubFrmSizes(sf *SubFrm, bdrat, b float64){
	//TODO COLUMN SIZES BM SIZES
	//dlmax := sf.DL * sf.PSFs[0]
	//llmax := sf.LL * sf.PSFs[2]
	//var mod *kass.Model
	var coords [][]float64
	var x float64
	coords = append(coords, []float64{x})
	for _, lspan := range sf.Lspans{
		x += lspan
		coords = append(coords, []float64{x})
	}
	fmt.Println(coords)
}

   //earlier clvr loads could be specified but eh
   
			switch bm{
				case sf.Lcmem:
				cdl := sf.DL; cll := sf.LL 
				if sf.Clvrs[0][1]+sf.Clvrs[0][2] > 0.0{
					cdl = sf.Clvrs[0][1]; cll = sf.Clvrs[0][2]
				}
				_ = sf.AddMemLoad(bm, []float64{1.0,3.0,cdl,0.0,0.0,0.0,1.0})
				_ = sf.AddMemLoad(bm, []float64{1.0,3.0,cll,0.0,0.0,0.0,2.0})
				case sf.Rcmem:
				cdl := sf.DL; cll := sf.LL 
				if sf.Clvrs[1][1]+sf.Clvrs[1][2] > 0.0{
					cdl = sf.Clvrs[1][1]; cll = sf.Clvrs[1][2]
				}
				_ = sf.AddMemLoad(bm, []float64{1.0,3.0,cdl,0.0,0.0,0.0,1.0})
				_ = sf.AddMemLoad(bm, []float64{1.0,3.0,cll,0.0,0.0,0.0,2.0})
				default:
				if sf.Fltslb{
					_ = sf.AddMemLoad(bm, []float64{1.0,3.0,sf.Lbay * sf.DL,0.0,0.0,0.0,1.0})
					_ = sf.AddMemLoad(bm, []float64{1.0,3.0,sf.Lbay * sf.LL,0.0,0.0,0.0,2.0})
					
				} else {
					_ = sf.AddMemLoad(bm, []float64{1.0,3.0,sf.DL,0.0,0.0,0.0,1.0})
					_ = sf.AddMemLoad(bm, []float64{1.0,3.0,sf.LL,0.0,0.0,0.0,2.0})
				}
			}
		}

	for _, i := range sf.Cols{
		em := sf.Mod.Mprp[i-1][2]; cp := sf.Mod.Mprp[i-1][3]
		jb, je := sf.Mod.Mprp[i-1][0], sf.Mod.Mprp[i-1][1]
		c1, c2 := sf.Mod.Coords[jb-1], sf.Mod.Coords[je-1]
		collen := EDist(c1, c2)
		sf.Colenv[i] = &kass.ColEnv{
			Id:i,
			EnvRez:make(map[int][]float64),
			Lspan:collen,
		}
		sf.RcCol[i] = &RccCol{
			Fck:sf.Fcks[em-1],
			Fy:sf.Fys[em-1],
			Cvrt:50.0,
			Cvrc:50.0,
			B:0.0,
			H:0.0,
			Styp:sf.Styps[cp-1],
			Rtyp:0,
			Dtyp:1,
			Nlayers:3,
			Lspan:collen,
			Code:sf.Code,
			Dims:sf.Secmap[i],
			Verbose:sf.Verbose,
			Kostin:sf.Kostin,
			Foldr:sf.Foldr,
		}
		switch sf.Styps[cp-1]{
			case -1:
			case 0:
			case 1:
			sf.RcCol[i].B = sf.Secmap[i][0]
			sf.RcCol[i].H = sf.Secmap[i][1]
			default:
			err := sf.RcCol[i].SecInit()
			if err != nil{
				fmt.Println(err)
			}	
		}
	}

*/
