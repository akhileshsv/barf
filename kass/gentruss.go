package barf

import (
	"fmt"
)

//Trs2d holds 2d truss generation struct variables
type Trs2d struct {
	Id         string
	Name       string //the ROOT OF ALL POWER is in naming
	Shade      int //shade
	Mtyp       int //material type 
	Styp       int //section type
	Sname      string //section name
	Cname      string //connection name
	Ctyp       int //connection type
	Typ, Cfg   int //truss type, configuration
	Code       int //design code
	Group      int //for wood (but make this universal)
	Bstyp      int //base section type
	Cmdz       []string //whaat, what is this good for
	Dzval      []float64 //dz val can be anything (might be useful )
	PSFs       [][]float64
	Roofmat    int //roofing material (see constants)
	Ngp        int //number of groups
	Frmtyp     int //truss or 2d frame
	Pdz, Cdz   bool //
	Clvr       bool //is cantilever
	Ldcalc     int
	Tcr, Bcr   int
	Pmat, Cmat int
	Brc        int //if 2 - brace both ways
	DL, LL     float64
	DLb        float64 //bcnode dead load
	Pz         float64
	Pdl, Pll   float64 //load/purlin or purlin dead load  
	Pwl, Pwr   float64 //load/purlin uplift case
	Pwp        float64 //load/purlin pressure/suction case
	Cpi        float64
	Em, Pg     float64
	Prln       string //section - "rect","tube",etc
	Prlntyp    string //cs/ss
	Units      string
	Term       string
	Loct       string //"roof", "floor"
	Span       float64
	Width      float64
	Height     float64
	Depth      float64
	Slope      float64 //dunno why but this is the "1 in slope" kinda slope
	Bcslope    float64
	Rise       float64
	Tang       float64 //truss angle (a struct is free real estate)
	Spacing    float64
	Rftrl      float64
	Vb         float64
	Girtspc    float64
	Purlinspc  float64
	Purlinwt   float64
	Endrat     float64
	Ovrhng     float64
	Wtrs       float64
	Nlfl       float64 //node (lat) strength per nail
	Nlfw       float64 //withdrawal strength per nail
	Prms       []float64
	Prop       []interface{}
	Nlds       int
	Clctyp     int
	Fop        int
	Opt        int
	Kpost      bool //add middle kingpost strut
	Spose      bool //make this default, all load combos are torture 
	Planld     bool //dead and live load calc on plan area
	Dzopt      bool
	Bcincl     bool
	Bent       bool
	Frame      bool
	Stltie     bool //if stltie, em = [timber, stl] for timber truss
	Verbose    bool
	Tweak      bool
	Spam       bool
	Web        bool
	Noprnt     bool
	Cols       []int
	Strt       []float64
	Coords     [][]float64
	Ms         [][]int
	Tcns, Bcns []int
	Nrs        int
	Ngs        int
	Mod        *Model
	Sections   [][]float64
	Cps        [][]float64
	Styps      []int
	Sects      []SectIn
	Ssecs      []StlSec
	Mldsrv     map[float64][][]float64 //ldtyp to service load case (1-dl, 2-ll, 3-wl)
	Jldsrv     map[float64][][]float64 //map of ldtyp to nodal load cases
	Dzmem      []interface{}
	Kostin     []float64
	Ntruss     float64
	// Condia     float64 //connection member dia
	// Conlen     float64 //connection member length
}

//Getndim returns ndims for truss opt
func (t *Trs2d) Getndim()(nd int){
	switch t.Mtyp{
		case 1:
		//rcc lmao
		case 2:
		//steel
		nd = t.Ngs
		case 3:
		//timber
		//ngrps = 4 (say)
		//tc(b,d); bc(b,d); struts/l/verts(b,d); ties/r/inclines(b,d)
		//now width of tc/bc = depth of webs
		//and depth of tc = depth of bc?
		//vars - chord width, chord depth, web width
		//for spaced - chord width, chord depth, web depth, web width
		switch t.Styp{
			case 0:
			//circular sect(roundwood)
			//dc, dl
			nd = 2
			case 1:
			//rect sect
			//bc, dc, bl, br
			//now is bl = br? YES
			nd = 3
			case 4:
			//tube sect
			//??
			case 26:
			//spaced col sect
			//bc, dc, dw, bl
			nd = 4
		}
		case 23:
		//timber truss + steel ties
		case 32:
		//see? and so on
	}
	return
}

//Init a 2d truss for gen 
func (t *Trs2d) Init() (err error){
	//groups = rafter,
	if t.Id == ""{t.Id = "gentrs"}
	if t.Units == ""{t.Units = "nmm"}
	//if len(t.Cmdz) < 2{t.Cmdz = []string{"2dt","mmks","1"}}
	if t.Mtyp == 0 {
		t.Mtyp = 2
	}
	if t.Styp == 0 {
		switch t.Mtyp{
			case 2:
			//okay now WTF
			t.Styp = 3
			t.Sname = "l"
			case 3:
			t.Styp = 1
		}
	}
	if t.Ctyp == 0{
		//decide what this is later
		t.Ctyp = 1
	}
	if t.Pz == 0 {
		t.Pz = 1.0
	}
	if t.Vb == 0 {
		t.Vb = 50.0
	}
	if t.Code == 0{t.Code = 1}
	t.Mldsrv = make(map[float64][][]float64)
	t.Jldsrv = make(map[float64][][]float64)
	if t.Em == 0.0{
		switch t.Mtyp{
			case 2:
			t.Em = 210000.0
			case 3:
			t.Em = 5600.0
		}
	}
	if t.Pg == 0{
		switch t.Mtyp{
			case 2:
			t.Pg = 7850.0
			case 3:
			t.Pg = 0.65
		}
	}
	if t.Mtyp == 3 && t.Group == 0{
		t.Group = 3
	}
	if t.Nlds == 0{
		switch t.Ldcalc{
			case -1:
			if t.Pdl == 0.0 && t.Pll == 0.0{
				err = fmt.Errorf("no dead or live loads specified- %f %f",t.DL, t.LL)
				return
			}
			if t.DL > 0.0{
				t.Nlds = 1
			}
			if t.LL > 0.0{
				t.Nlds += 1
			}
			if t.Pwl > 0.0{
				t.Nlds = 3
			}
		}
	}
	
	t.Mod = &Model{Id:t.Id,Cmdz:t.Cmdz,
		Frmtyp:2,Ncjt:2,Term:t.Term,
		Verbose:t.Verbose,Frmstr:"2dt",
		Noprnt:t.Noprnt,Web:t.Web,
		Units:t.Units,Spose:t.Spose,
		Mtyp: t.Mtyp,Ngrps:t.Ngs,
	}
	return
}

//GenGeom generates basic 2d truss geometry (coords, mprp) from templates
func (t *Trs2d) GenGeom() (err error){
	switch t.Typ {
	case 1, 2: //"a", "l"/cantilever:
		
		err = GenATruss(t)
	case 3, 4://"trap-a", "trap-l":
		err = GenTrapTruss(t)
	case 5: //"parallel":
		err = GenParTruss(t)
	case 6:
		//bowstring/funicular truss
		err = GenFunTruss(t)
	case -1:
		//generate groundstruct
	default:
		//ERRORE,errore
		//return mod, errors.new("invalid truss string")
	}
	if err != nil{return}
	
	// if t.Term != "" && t.Term != "svg"{
	// 	PlotGenTrs(t.Coords, t.Ms)
	// }
	t.Mod.Ngrps = t.Ngs
	err = nil
	return
}

//Calc generates and calcs/analyzes a Trs2d
func (t *Trs2d) Calc() (err error){
	err = t.Init()
	if err != nil{
		return
	}
	//t.Term = "dumb"
	err = t.GenGeom()
	if err != nil{return}
	err = t.GenCp()
	if err != nil{return}
	err = t.GenLd()
	if err != nil{return}
	switch t.Spose{
		case true:
		err = CalcSrvLds(t.Mod)
		case false:
		err = CalcLdCombos(t.Mod)
	}
	if err != nil{return}
	// if t.Spam{
	// 	// for lp, flp := range t.Mod.Lpmap{
	// 	// 	fmt.Println("load pattern->",lp," float->",flp)
	// 	// 	fmt.Println(t.Mod.Txtplots[lp-1])
	// 	// 	fmt.Println(t.Mod.Reports[lp-1])
	// 	// }
	// 	// fmt.Println(t.Mod.Reports[len(t.Mod.Reports)-1])
	// }
	return
}

//GenCp generates initial cross section values for a 2d truss
func (t *Trs2d) GenCp() (err error){
	t.Mod.Em = [][]float64{{t.Em}}
	t.Sects = make([]SectIn, t.Ngs)
	t.Mod.Cp = make([][]float64,t.Ngs)
	if t.Sections == nil{
		switch t.Mtyp{
			case 2:
			//steel
			t.InitStlCp()
			case 3:
			//tmbr
			t.InitTmbrCp()
		}
	}
	// if len(t.Sections) != t.Ngs{
	// 	return errors.New("invalid number of truss groups")
	// }
	switch t.Mtyp{
		case 2:
		for i, dims := range t.Sections{
			styp := t.Styps[i]
			sdx := int(dims[0])
			tp := 0.0
			if len(dims) > 1{
				tp = dims[1]
			}
			sname := StlSnames[styp]
			ss, e := GetStlSec(sname, sdx, t.Code, tp)
			if e != nil{
				err = fmt.Errorf("error reading sec-%s",e)
				return
			}
			t.Ssecs = append(t.Ssecs, ss)
			t.Mod.Cp[i] = []float64{ss.Area, ss.Ixx}
			t.Mod.Dims = append(t.Mod.Dims, ss.Dims)
			t.Mod.Sts = append(t.Mod.Sts, ss.Bstyp)
		}
		case 3:	
		for i, dims := range t.Sections{
			styp := t.Styps[i]
			s := SecGen(styp, dims)
			t.Sects[i] = s
			t.Mod.Cp[i] = []float64{s.Prop.Area, s.Prop.Ixx}
		}
	}
	return
}

//InitStlCp initializes a stl truss's cp
func (t *Trs2d) InitStlCp(){
	if t.Ngs == 0{
		t.Ngs = 5
	}
	t.Ssecs = make([]StlSec, t.Ngs)
	t.Sections = make([][]float64, t.Ngs)
	t.Styps = make([]int, t.Ngs)
	switch t.Sname{
		case "l":
		case "tube":
		//tubular steel truss
	}
	//t.Sections for steel  - [sdx, params]
	//use bhav. minimum sections if span > 10.0
	if t.Span < 10.0{
		t.Sections = [][]float64{
			{33,8},
			{33,8},
			{23,8},
			{23,8},
			{19,8},
		}
		t.Styps = []int{105,105,105,105,105}
	} else {
		t.Sections = [][]float64{
			{26,8},
			{26,8},
			{19,8},
			{19,8},
			{19,8},
		}
		t.Styps = []int{4,4,4,4,4}
        }
	if t.Ngs == 4{
		t.Sections = t.Sections[:4]
		t.Styps = t.Styps[:4]
	}
	t.Mod.Dims = t.Sections
	t.Mod.Sts = t.Styps
	return
}

//InitTmbrCp initializes a tmbr truss's cp
//TODO - use snames
func (t *Trs2d) InitTmbrCp(){
	if t.Ngs == 0{
		t.Ngs = 5
	}
	//get em
	t.Sections = make([][]float64, t.Ngs)
	t.Styps = make([]int, t.Ngs)
	for i := 0; i < t.Ngs; i++{
		var dims []float64
		switch i{
			case 0, 1:
			//top and bottom chords
			switch t.Styp{
				case 0:
				dims = []float64{4.0*25.4}
				case 1:
				dims = []float64{50.0,125.0}
			}
			default:
			switch t.Styp{
				case 0:
				dims = []float64{4.0*25.4}
				case 1:
				dims = []float64{50.0,100.0}
			}
		}	
		t.Sections[i] = dims
		t.Styps[i] = t.Styp
	}
	t.Mod.Dims = t.Sections
	t.Mod.Sts = t.Styps
}

// //PrlnDz designs a purlin section for a roof truss
// //(simplified) uses b.m coeff based on Prlntyp("ss/cs")
// func (t *Trs2d) PrlnDz(){
// 	//var bmcf, vcf float64
// 	switch t.Prlntyp{
// 		case "clvr":
// 		//bmcf = 0.5
// 		case "cs":
// 		//bmcf = 0.1
// 		case "ss":
// 		//bmcf = 0.8
// 	}
// 	switch t.Mtyp{
// 		case 2:
// 		//steel
// 		switch t.Prln{
// 			case "i":
// 			case "c":
// 			case "l":
// 		}
// 		case 3:
// 		//tmbr
// 		switch t.Prln{
// 			case "tube":
// 			case "rect":
// 		}
// 	}
// }

/*
func (*T2din) getprlnspc (dispopt bool) {
	sheetlbls := []string{"asbestos","mangalore","double tiles","copper","bitumen","thatch","GI","PUF panels","Steel","SP38","Sub12.1"}
	sheetwts := [][]float64 {
		{0.17,1200},{0.785,750},{1.67,750},{0.72,1200},{0.102,750},{0.49,750},{1.60,1400},{0.4,1400},{0.08,1200},{0.2,1400},{0.21,1400},
		}
	if dispopt {
		fmt.Printf("%s",sheetlbls)
		fmt.Printf("%.3f", sheetwts)
	}

}
*/
