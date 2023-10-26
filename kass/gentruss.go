package barf

import (
	"fmt"
	//"math"
	//"os"
	//"strings"
	"errors"
	//draw"barf/draw"
	//kass "barf/kass"
	//bash "barf/bash"
)

//Trs2d holds 2d truss generation struct variables
type Trs2d struct {
	Id         string
	Cmdz       []string
	Mtyp       int //material type 
	Styp       int //section type
	Typ, Cfg   int //truss type, configuration
	Code       int //design code
	Group      int //for wood (but make this universal)
	Dzval      []float64 //dz val can be anything (might be useful )
	PSFs       [][]float64
	Roofmat    int //roofing material (see constants)
	Ngp        int //number of groups
	Frmtyp     int //truss or 2d frame 
	Pdz, Cdz   bool
	Ldcalc     int
	Tcr, Bcr   int
	Pmat, Cmat int
	DL, LL     float64
	DLb        float64 //bcnode dead load
	Pz         float64
	Cpi        float64
	Em, Pg     float64
	Term       string
	Span       float64
	Width      float64
	Height     float64
	Depth      float64
	Slope      float64 //dunno why but this is the "1 in slope" kinda slope
	Bcslope    float64
	Rise       float64
	Spacing    float64
	Rftrl      float64
	Vb         float64
	Girtspc    float64
	Purlinspc  float64
	Endrat     float64
	Ovrhng     float64
	Wtrs       float64
	Prms       []float64
	Prop       []interface{}
	Clctyp     int
	Opt        int
	Dzopt      bool
	Bcincl     bool
	Bent       bool
	Frame      bool
	Stltie     bool
	Verbose    bool
	Cols       []int
	Strt       []float64
	Coords     [][]float64
	Ms         [][]int
	Tcns, Bcns []int
	Nrs        int
	Ngs        int
	Mod        *Model
	Sections   [][]float64
	Styps      []int
	Sects      []SectIn
	Mldsrv     map[float64][][]float64 //map of ldtyp to service load case (1-dl, 2-ll, 3-wl)
	Jldsrv     map[float64][][]float64 //map of ldtyp to nodal load cases
}

//Init a 2d truss for gen 
func (t *Trs2d) Init() {
	//groups = rafter,
	if t.Id == ""{t.Id = "gentrs"}
	if t.Mtyp == 0 {
		t.Mtyp = 2
	}
	if t.Styp == 0 {
		t.Styp = 1
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
	t.Mod = &Model{Id:t.Id,Cmdz:t.Cmdz,Frmtyp:2,Ncjt:2,Term:t.Term,Verbose:t.Verbose}
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

	return
}

//GenGeom generates basic 2d truss geometry (coords, mprp) from templates
func (t *Trs2d) GenGeom() (err error){
	t.Init()
	switch t.Typ {
	case 1, 2: //"a", "l":
		err = GenATruss(t)
	case 3, 4://"trap-a", "trap-l":
		err = GenTrapTruss(t)
	case 5: //"parallel":
		err = GenParTruss(t)
	case 6:
		//bowstring
	case -1:
		//generate groundstruct
	default:
		//ERRORE,errore
		//return mod, errors.new("invalid truss string")
	}
	if err != nil{return}
	
	//plot- REMOVE DIS
	/*
	if t.Term != "" && t.Term != "svg"{
		txtplt, e := draw.PlotGenTrs(t.Coords, t.Ms, t.Term,"t2d")
		if e != nil{
			//log.Println("ERRORE, errore->",e)
			err = e
			return
		}
		if t.Term != "qt"{fmt.Println(txtplt)}
	}
	*/
	err = nil
	return
}

//Calc generates and calcs/analyzes a Trs2d
func (t *Trs2d) Calc() (err error){
	t.Init()
	err = t.GenGeom()
	if err != nil{return}
	err = t.GenCp()
	if err != nil{return}
	err = t.GenLd()
	if err != nil{return}
	err = CalcLdCombos(t.Mod)
	if err != nil{return}
	return
}

//GenCp generates initial cross section values for a 2d truss
func (t *Trs2d) GenCp() (err error){
	t.Sects = make([]SectIn, t.Ngs)
	t.Mod.Em = [][]float64{{t.Em}}
	t.Mod.Cp = make([][]float64,t.Ngs)
	if t.Sections == nil{
		t.InitCp()
	}
	if len(t.Sections) != t.Ngs{
		return errors.New("invalid number of truss groups")
	}
	for i, dims := range t.Sections{
		styp := t.Styps[i]
		fmt.Println(ColorRed,styp, dims,ColorReset)
		s := SecGen(styp, dims)
		t.Sects[i] = s
		t.Mod.Cp[i] = []float64{s.Prop.Area, s.Prop.Ixx}
	}
	return
}

//InitCp inits the base model's cp array
func (t *Trs2d) InitCp() {
	//get em
	t.Sections = make([][]float64, t.Ngs)
	t.Styps = make([]int, t.Ngs)
	switch t.Mtyp{
		case 2:
		//steel
		
		case 3:
		//wood
		for i := 0; i < t.Ngs; i++{
			var dims []float64
			switch i{
				case 0, 1:
				//top and bottom chords
				switch t.Styp{
						case 0:
						dims = []float64{4.0*25.4}
						case 1:
						dims = []float64{50.0,100.0}
				}
				default:
				switch t.Styp{
					case 0:
					dims = []float64{4.0*25.4}
					case 1:
					dims = []float64{50.0,80.0}
				}
			}	
			t.Sections[i] = dims
			t.Styps[i] = t.Styp
		}
	}
	t.Mod.Dims = t.Sections
	t.Mod.Sts = t.Styps
	return	
}

//wh-aat, what is this good for
//RIGHT?
func GenKp(t *Trs2d) (err error){
	return
}

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
