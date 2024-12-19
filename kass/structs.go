package barf

import (
	"errors"
	"gonum.org/v1/gonum/mat"
)

//rez1d is a general struct to store int, float results
type rez1d struct{
	i int
	val float64
}

//rez2d stores int, int, float results
type rez2d struct{
	i, j int
	val float64
}

//ModRez stores model results
//what is it good for
//absolutely nothing
type ModRez struct{
	Report string
	Txtplot string
	
}

//Model stores a (kassimali type) model data fields
/*
   truss 2d:
   Commands/Parameters (model[7], cmdz []string) TYPE,UNITS,DESIGN,DESIGN TYPE(rcc/steel/wood)
   2dt;kips;0
   Coordinates(model[0],cords [][]float64) X1,Y1 ; X2,Y2..
   0,0; 144,0; 288,0; 288,196
   Supports(model[1],msup [][]int) J1,R1X,R1Y;J2,R2X,R2Y (-1 indicates restraint, else 0)
   1,-1,-1; 2,-1,-1; 3,-1,-1
   Elastic Modulus(model[2],em []float64) EM1;EM2.. 
   29000
   Cross-section(model[3],cp [][]float64) A1 ; A2
   8;6
   Members(model[4], mprp [][]int) JB1,JE1,EM1,CP1,MREL; JB2,JE2
   1,4,1,1; 2,4,1,2; 3,4,1,1
   Joint P-loads(model[5], jp [][]float64) J1,FX,FY;
   4,150,-300
   Member P-loads(model[6], msp [][]float64) M1,LTYP,W1,W2,L1,L2
*/
/*
   beam 1d:
   Commands/Parameters (model[7], cmdz []string) TYPE,UNITS,DESIGN,DESIGN TYPE(rcc/steel/wood)
   1db;kips;0
   Coordinates(model[0],cords [][]float64) X1; X2; ..
   0; 120; 360; 480
   Supports(model[1],msup [][]int) J1,R1X,R1Y;J2,R2X,R2Y (-1 indicates restraint, else 0)
   1,-1, 0; 2,-1, 0; 3,-1, 0; 4,-1,-1
   Elastic Modulus(model[2],em []float64) EM1;EM2.. 
   29000
   Cross-section(model[3],cp [][]float64) IZ1 ; IZ2
   350;500
   Members(model[4], mprp [][]int) JB1,JE1,EM1,CP1,MREL; JB2,JE2,...
   1,2,1,1,0; 2,3,1,1,0; 3,4,1,2,0
   Joint P-loads(model[5], jp [][]float64) J1,FX,FY;
   1,0,-480
   Member P-loads(model[6], msp [][]float64) M1,LTYP,W1,W2,L1,L2
   2,3,0.1667,0,0,120; 2,1,25,0,180,0; 3,4,0.25,0,0,0
*/
/*
   frame 2d:
   Commands/Parameters (model[7], cmdz []string) TYPE,UNITS,DESIGN,DESIGN TYPE(rcc/steel/wood)
   2df;mks;0
   Coordinates(model[0],cords [][]float64) X1,Y1 ; X2,Y2..
   -12,0,-8 ; 6,0,-8; 0,0,0 ; -6,0,8; 12,0,8;0,24,0
   Supports(model[1],msup [][]int) J1,R1X,R1Y,R1M; (-1 indicates restraint, else 0)
   1,-1,-1,-1; 2,-1,-1,-1; 3,-1,-1,-1; 4,-1,-1,-1
   Elastic Modulus(model[2],em []float64) EM1;EM2.. 
   29000.0
   Cross-section(model[3],cp [][]float64) A1,IZ1 ; A2,IZ2
   14.7,800.0
   Members(model[4], mprp [][]int) JB1,JE1,EM1,CP1,MREL; JB2,JE2
   1,3,1,1,2 ; 2,4,1,1,0 ; 3,4,1,1,1
   Joint P-loads(model[5], jp [][]float64) J1,FX,FY,MZ;
   2,25.0,0.0,0.0 
   Member P-loads(model[6], msp [][]float64) M1,LTYP,W1,W2,L1,L2
   1,3,0.1,0,0,0; 3,1,75,0,120,0
*/
/*
   truss 3d:
   Commands/Parameters (model[7], cmdz []string) TYPE,UNITS,DESIGN,DESIGN TYPE(rcc/steel/wood)
   3dt;kips;0
   Coordinates(model[0],cords [][]float64) X1,Y1 ; X2,Y2..
   -144,0,-96 ; 72,0,-96; -72,0,96; 144,0,96; 0,288,0
   Supports(model[1],msup [][]int) J1,R1X,R1Y,R1M; (-1 indicates restraint, else 0)
   1,-1,-1,-1; 2,-1,-1,-1; 3,-1,-1,-1; 4,-1,-1,-1
   Elastic Modulus(model[2],em []float64) EM1;EM2.. 
   10000.0
   Cross-section(model[3],cp [][]float64) A1,IZ1 ; A2,IZ2
   8.4
   Members(model[4], mprp [][]int) JB1,JE1,EM1,CP1,MREL; JB2,JE2
   1,5,1,1; 2,5,1,1; 3,5,1,1; 4,5,1,1
   Joint P-loads(model[5], jp [][]float64) J1,FX,FY,FZ;
   5,0,-100.0,-50.0 
   Member P-loads(model[6], msp [][]float64) M1,LTYP,W1,W2,L1,L2

*/
/*
   grillage/grid 3d:
   Commands/Parameters (model[7], cmdz []string) TYPE,UNITS,DESIGN,DESIGN TYPE(rcc/steel/wood)
   3dg;mks;0
   Coordinates(model[0],cords [][]float64) X1,Y1,Z1 ; X2,Y2, Z2..
   0,0,0;8,0,0;8,0,6;0,0,6
   Supports(model[1],msup [][]int) J1,R1Y(y disp.),R1X(x rotation),R1Z(z rotation); (-1 indicates restraint, else 0)
   1,-1,-1,-1; 2,-1,-1,-1; 4,-1,-1,-1
   Elastic Modulus, Shear Modulus(model[2],em []float64) EM1, G1; EM2, G2.. 
   200e6,76e6
   Cross-section(model[3],cp [][]float64) IZ1,J1 ; IZ2,J2..
   347e-6,115e-6
   Members(model[4], mprp [][]int) JB1,JE1,EM1,CP1,MREL; JB2,JE2
   1,3,1,1,0; 2,3,1,1,0; 4,3,1,1,0
   Joint P-loads(model[5], jp [][]float64) J1,FX,FY,FZ;
   Member P-loads(model[6], msp [][]float64) M1,LTYP,W1,W2,L1,L2
   2,3,20,0,0,0; 3,3,20,0,0,0

*/
/*
   frame 3d:
   Commands/Parameters (model[7], cmdz []string) TYPE,UNITS,DESIGN,DESIGN TYPE
   3df;kips;0
   Coordinates(model[0],cords [][]float64) X1,Y1,Z1 ; X2,Y2, Z2..
   0,0,0; -240,0,0; 0,-240,0; 0,0,-240
   Supports(model[1],msup [][]int) J1,RX,RY,RZ,MX,MY,MZ;   (-1 indicates restraint, else 0)
   2,-1,-1,-1,-1,-1,-1; 3,-1,-1,-1,-1,-1,-1; 4,-1,-1,-1,-1,-1,-1
   Elastic Modulus, Shear Modulus(model[2],em [][]float64) EM1, G1; EM2, G2.. 
   29000,11500
   Cross-section(model[3],cp [][]float64) A1,IZ1,IY1,J1 ; IZ2,J2..
   32.9,716,236,15.1
   Members(model[4], mprp [][]int) JB1,JE1,EM1,CP1,MREL,MANG; 
   2,1,1,1,0,1; 3,1,1,1,0,2; 4,1,1,1,0,3
   Joint P-loads(model[5], jp [][]float64) J1,FX,FY,FZ,MX,MY,MZ;
   1,0,0,0,-1800,0,1800
   Member P-loads(model[6], msp [][]float64) M1,LTYP,W1,W2,L1,L2,AX
   1,3,0.25,0,0,0,0 
   Angle of roll(model[7], wng [][]float64) W1,MEM(0-horizontal, 1-vertical, 2-others);W2,AX2...
   0,0;1,90;2,30
*/
type Model struct {
	//ze backbone
	//see old text input files above pls
	Id       string //title/ID
	Cmdz     []string //commands , this is not exactly frozen (TYPE,UNITS,DESIGN,DESIGN TYPE(rcc/steel/wood))
	Units    string //"mks","mmks","mns","mmns","kips"
	Frmstr   string //1db, 2dt, 2df, 3dt, 3dg, 3df      
	Ncjt     int //ndof per node 2 - truss, 2 - beam, 3 - frame 2d, 3 - grillage, 3 - truss 3d, 6 - space frame
	Frmtyp   int //1- beam, 2- 2d truss, 3 - 2d frame, 4 - 3d truss, 5 - 3d grid, 6 - 3d frame
	Ndof     int //no. of degrees of freedom
	Coords   [][]float64 //slice of nodal coordinates {p1(x, y, z), p2(x,y,z)}
	Supports [][]int //slice of restrained supports (nodal index, -1 if rest.in ndof 1, etc)
	Em       [][]float64 //slice of material properties (truss - elastic modulus, beam - , frm2d - , trs3d - ,grd3d - , frm3d - )
	Cp       [][]float64 //cross section properties (truss - area, beam - area, iz, frm2d - , trs3d - , grd3d - , frm3d - ) 
	Mprp     [][]int //slice of members      
	Jloads   [][]float64 //joint loads
	Msloads  [][]float64 //member loads
	Wng      [][]float64 //angle of roll
	Jsd      []int //joint support displacement indices
	Sdj      [][]float64 //support displacement slice
	Clvrs    [][]float64 //cantilevers - is this needed? YES
	PSFs     []float64 //lcat 1 - adv, ben, 2 - adv, ben , etc
	Psfs     [][]float64 //new one lol
	Ldfacs   [][]float64 //[lcase 1 - lcat 1, lcat 2...], [lcase 2 - lcat 1, 2...]
	Ldcs     [][][]float64 //list of loadcases
	Dims     [][]float64
	Ts       [][]int
	Bs       [][]float64
	Ds       [][]float64
	Ls       [][]float64
	DLs      [][]float64 //msloads
	LLs      [][]float64 //
	WLs      [][]float64 //
	NDLs     [][]float64 //nodal dead loads
	NLLs     [][]float64
	NWLs     [][]float64
	Mtprp    [][]float64 //material properties slice 
	Lsx      []float64 //support lengths
	Sts      []int
	Inp      []interface{} //use anything for inputs
	Code     int //1 - desi, 2 - bridesi
	Sectyp   int //if non zero, base section type
	Opt      int //if > 0, optimize using dims/sections/styps
	Ngrps    int //number of model section groups
	Rdinp    bool //if true, read vals from input - else check other params, etc
	Verbose  bool 
	Spam     bool //print a fuckton of messages
	VDx      bool //account for shear deformation
	Nocolor  bool //turn off ansi escape codes
	Web      bool
	Drawsec  bool //draw 2d section views
	Dtyp     int //design type (0 - nothing, 1 - cbm?) not used
	Mtyps    []int
	Mtyp     int
	Mstr     string //rcc, stl, tmbr (1, 2, 3)
	Group    int //use for rapid classification i guess?(wood group, etc)
	Nspans   int //use if cbeam
	Lspans   []float64 //use for cbeam/frm2d?
	Ldcalc   int //use ldcases if 1
	Nlp      int //same as below lmao
	Nlc      int //no. of load combos
	Nlds     int //no. of load types
	Nwlc     int //no. of wind load cases
	Nslc     int //no. of seismic load cases
	Same     bool //single section/cp
	Nonlin   bool //non linear analysis
	Vibe     bool //vibrational analysis
	Npsec    bool //non prismatic members
	Calc     bool //calc bm, sf, dxs
	Genwng   bool //generate member wng vectore
	Envcalc  bool //calc load envelopes if true
	Dz       bool //design?
	Noprnt   bool //use for not printing tables during opt
	Zup      bool //z axis vertical(up) in 3d model plots
	Split    bool //has split members
	Spose    bool //superpose load combos (calc once per load case, add rez * psfs)
	Cols     []int
	Beams    []int
	Slbs     [][]int
	Csec     []int
	Bsec     []int
	Xs       []float64
	Ys       []float64
	Zs       []float64
	Term     string
	Foldr    string
	Report   string
	Reports  []string
	Txtplot  string
	Txtplots []string
	Type     string
	Mldsrv   map[float64][][]float64 `json:"-"`
	Jldsrv   map[float64][][]float64 `json:"-"`
	Mldcs    map[float64][][]float64 `json:"-"`
	Jldcs    map[float64][][]float64 `json:"-"`
	Jmap     map[tupf]int `json:"-"`
	Grez     map[int][]float64 `json:"-"`
	Submems  map[int][]int `json:"-"`
	Bmenv    map[int]*BmEnv 
	Ms       map[int]*Mem 
	Mnps     map[int]*MemNp
	Js       map[int]*Node 
	Lpmap    map[int]float64 
	Frcscale float64
	Dmin     float64
	Grade    float64
	Fy       float64
	Pg       float64
	Dmax     float64
	Weight   float64
	Kost     float64
	Icols    []interface{}
	Ibms     []interface{}
	Matprop  []float64
	Scales   []float64
	Sects    []SectIn `json:"-"`
}

//Mem is a struct to store model member fields
type Mem struct {
	Id     int
	Mprp   []int
	Geoms  []float64
	Qf     []float64
	Gf     []float64
	Vfs    []float64
	Bkmat  *mat.Dense
	Tmat   *mat.Dense
	Gkmat  *mat.Dense
	Memtyp []string
	Gmat   *mat.Dense
	Mmat   *mat.Dense
	Ktmat  *mat.Dense
	Mtyp   int
	Dtyp   int
	Styp   int
	Venv   []float64
	Mpenv  []float64
	Mnenv  []float64
	Pu     []float64
	Sec    *SectIn
	Tmax   float64
	Cmax   float64
	Qfr    [][]float64
	Gfr    [][]float64
	Lds    [][]float64
	Pmax   float64
	Vu     float64
	Mu     float64
	Dmax   float64
	Rez    BeamRez
	EnvRez map[string]BeamRez
	Xs     []float64
	Dz     bool
	Clvr   bool
	Dzinp  [][]float64
	Dzrez  [][]float64
	Dzvals [][]float64
	Params []float64 //design params
	Names  []string //design names
}

//Node is a struct to store... you guessed it, node struct fields
type Node struct {
	Id       int
	Coords   []float64
	Dcs      []float64
	Supports []int
	Nrs      int
	Nscs     []int
	Displ    []float64
	React    []float64
	Mems     []int
	Frcs     []float64
	Sd       bool
	Type     string  //could be anything
	Sdj      []float64
	Dmax     []float64
	Rmax     []float64
	Dr       [][]float64
	Rr       [][]float64
	Dzmem    []interface{}
	Fxtrs    [][]float64
	Nblts    float64
	Ctyp     int
	
}

//MLoad might be used sometime when tweak params are a thing
//member load struct
type MLoad struct{
	Mem, Ltyp int
	Wa, Wb, La, Lb float64
	Lcon int
	Lax  int
}

//NLoad might be used sometime
//to store/input nodal loads
type NLoad struct {
	Node int
	Fx, Fy, Fz, Mx, My, Mz float64
	Vec []float64
}

//MemNp for non uniform member struct fields 
type MemNp struct {
	//this is some kind of disaster
	Id     int
	Styp   int
	Mprp   []int
	Geoms  []float64
	Qf     []float64
	Gf     []float64
	Em     float64
	Vp     float64
	G      float64
	Bkmat  *mat.Dense
	Tmat   *mat.Dense
	Gkmat  *mat.Dense
	Memtyp []string
	Gmat   *mat.Dense
	Mmat   *mat.Dense
	Ts     []int
	Ls     []float64
	Ds     []float64
	Bs     []float64
	Dxs    []float64
	Bxs    []float64
	Txs    []int
	I0     float64
	I1     float64
	I2     float64
	I3     float64
	I4     float64
	I5     float64
	M11    float64
	M22    float64
	M12    float64
	N11    float64
	A0     float64
	Fs     float64
	B0     float64
	T1     float64
	T2     float64
	T3     float64
	Ka     float64
	Kb     float64
	Kc     float64
	Ca     float64
	Cb     float64
	Ix     []float64
	Ax     []float64
	Xs     []float64
	Dims   []float64
	M1     []float64
	M2     []float64
	M3     []float64
	M4     []float64
	M5     []float64
	Hmr    []float64
	Lspan  float64
	Frmtyp string
	Cx     float64
	Cy     float64
	Cz     float64
	Mx     []float64
	Vx     []float64
	Vfs    []float64
	Mtyp   int
	Venv   []float64
	Mpenv  []float64
	Mnenv  []float64
	Pu     []float64
	Tmax   float64
	Cmax   float64
	Lds    [][]float64
	Pmax   float64
	Vu     float64
	Mu     float64
	Dmax   float64
	Rez    BeamRez
	EnvRez map[string]BeamRez
	Dz     bool
	Clvr   bool
}

//T2dgen is not used now that we have Truss2d
type T2dgen struct {
	//import gen into kass?
	ID        string
	Type      []string
	Roofmat   int
	Section   []string
	Span      float64
	Width     float64
	Height    float64
	Slope     float64
	Spacing   float64
	Vb        float64
	Basemat   int
	Inperm    int
	Purlinspc float64
}

//Rcc2dInput is a horrible name and besides it is now Frm2d
type Rcc2dInput struct {
	Id     string
	Dims   []float64
	Bays   [][]float64
	Ws     []float64
	Ps     []float64
	Baystr string
	Ncjt   int
}


//Rcc2dFrame is again, a horrible name and besides it is now Frm2d
type Rcc2dFrame struct {
	Id       string
	Dims     []float64
	Bays     [][]float64
	W        []float64
	P        []float64
	Baystr   string
	Ncjt     int
	Coords   [][]float64
	Supports [][]int
	Em       [][]float64
	Cp       [][]float64
	Mprp     [][]int
	Jloads   [][]float64
	Msloads  [][]float64
}


type TrussParams struct {
	ID      string
	Type    int
	Config  int
	Roofmat int
	Section int
	Span    float64
	Width   float64
	Height  float64
	Slope   float64
	Vb      float64
	Basemat int
	Inperm  int
}

type RcframeParams struct {
	ID      string
	Type    int
	Config  int
	Roofmat int
	Section int
	Span    float64
	Width   float64
	Height  float64
	Slope   float64
	Vb      float64
	Basemat int
	Inperm  int
}

type RccMod3d struct {
	X, Y, Z                         float64
	coords                          [][]float64
	nodecords                       map[int][]float64
	supports                        map[int][]int
	cols, beamxs, beamys, slabnodes [][]int
}


var (
	ErrDim = errors.New("invalid input dimensions")
	ErrFact = errors.New("non +ve definite stiffness matrix")
	ErrSolve = errors.New("near singular stiffness matrix")
)

var (
	// trsMat = []string{"timber", "steel", "rcc", "cfs", "timber+steel"}
	// //trsMat = map[int]string{1:"steel-250",2:"steel-cfs",3:"wood"}

	// trsSec = map[int][]string{
	// 	1: {"tube", ""},
	// 	2: {"s.angle", ""},
	// 	3: {"d.angle", ""},
	// }
	// trsTyp = map[int]string{
	// 	1: "l", 2: "a", 3: "trap-l", 4: "trap-a", 5: "parallel",
	// }
	// trsCfg = map[int]string{
	// 	1: "howe", 2: "howe fan", 3: "pratt", 4: "pratt fan", 5: "fink fan", 6: "bowstring",
	// }
	rfMat = []string{"Asbestos cement sheets", "Mangalore tiles with battens", "Double tiles with battens",
		"Copper sheet", "Bitumen", "Thatch with battens", "GI sheets",
		"PUF panels", "Steel sheet (1 mm)", "SP38 example", "Sub12.1",
		"25mm board", "16mm board", "abel ex"}
	rfWt = [][]float64{
		{0.17, 1200},
		{0.785, 750},
		{1.67, 750},
		{0.72, 1200},
		{0.102, 750},
		{0.49, 750},
		{1.60, 1400},
		{0.4, 1400},
		{0.08, 1200},
		{0.2, 1400},
		{0.21, 1400},
		{0.18, 1200},
		{0.12, 1000},
		{0.28, 2000},
	}
)

/*
   f3d - 3d frame input struct
   fck - [ftng, slb, bm, col]
   fy - [ftng, slb, bm, col]
   dia - [ftng, slb main, slb dist, bm1, bm 2, c1, c2]
   stair - [l/r/t/b,n-s/e-w]
   clvr - [l,r,t,b] [l] - [len, dl, ll]
   ftngtyp - 0 - pad, 1 - sloped
   plinth - 0 - none, 1 - first level
   cdim -
   ftyp - 1 - 1 way slab, 2-2 way slabs
   flangex - 1 - x beams flanged
   flangey - 1 - y beams flanged
   dl - x-[dl 0, dl 1, dl 2], y-[dl 0, dl 1, dl 2]
   ll - same
   if TSlb > 0.0 - calc slab loads and add to beams
   if TSlb = -1.0 - calc dl
*/

//nodecords map[int][]float64,supports map[int][]int, cols, beamxs, beamys, slabnodes [][]int, mloads, mverts map[string][][]float64, xf, yf, zxf, zyf map[int][]string, foldr string, plotchn chan string

type tupil struct{
	i,j int
}

type tupf struct{
	i,j float64
}

type Frmdat struct{
	Nodecords map[int][]float64
	Supports map[int][]int
	Cols, Beamxs, Beamys, Slabnodes [][]int
	Mloads, Mverts map[string][][]float64
	Members map[string][][]int
	Xf, Yf map[int][]string
	Zxf, Zyf map[int]map[int][]string
	Foldr string
	Cdim, Bdim [][]float64
	Csec, Bsec []int
	X,Y,Z []float64
	Nodeadj map[int][]int
}


/*
type Rcc2dInput struct {
	Id     string
	Dims   []float64
	Bays   [][]float64
	Ws     []float64
	Ps     []float64
	Baystr string
	Ncjt   int
}

type Rcc2dFrame struct {
	Id       string
	Dims     []float64
	Bays     [][]float64
	W        []float64
	P        []float64
	Baystr   string
	Ncjt     int
	Coords   [][]float64
	Supports [][]int
	Em       [][]float64
	Cp       [][]float64
	Mprp     [][]int
	Jloads   [][]float64
	Msloads  [][]float64
}

type RccParams struct {
	Fcks   []float64
	Fys    []float64
	Fyvs   []float64
	Vb     float64
	Mredis float64
	Xs     []float64
	Ys     []float64
	Zs     []float64
	DLs    [][]float64
	LLs    [][]float64
	NFrms  int
	SBC    float64
	Kwa    map[string][]interface{}
}

type RccMod3d struct {
	X, Y, Z                         float64
	coords                          [][]float64
	nodecords                       map[int][]float64
	supports                        map[int][]int
	cols, beamxs, beamys, slabnodes [][]int
}

type Bm1d struct {
	ID              string
	Type            []string
	Fck, Fy, DL, LL float64
	Dims            []float64
}

type Frm3d struct {
	Dimx, Dimy []float64
	Colgrid    [][]int
	Slbgrid    [][]int
}

type FltSlb struct {
	X, Y, Z []float64
	Ctyp    int
	Cext    []float64
	Cint    []float64
	Cdim    []float64
	Ebm     []float64
	Lclvr   []float64
	Rclvr   []float64
	Dused   float64
}

*/
