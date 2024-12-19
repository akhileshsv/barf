package barf

import (
	"fmt"
	kass"barf/kass"
)

func (b *WdBm) printz() (string){
	var rez string
	rez += fmt.Sprintf("timber beam group %v section %v \n",b.Grp,b.Styp)
	rez += fmt.Sprintf("allowable stress (N/mm2)\nbending %.1f tension %.1f \ncomp parallel %.1f perp %.1f \nshear parallel %.1f perp %.1f\nEm %.1f\n",b.Prp.Ft, b.Prp.Ft, b.Prp.Fc, b.Prp.Fcp, b.Prp.Fv, b.Prp.Fvp, b.Prp.Em)
	return rez
}

//WdBm is a struct that stores timber beam fields
//see chapter 8, abel o.o
type WdBm struct{
	Id       int
	Title    string
	Styp     int
	Grp      int
	Code     int
	Endc     int
	Kendc    int //end condition for slenderness ratio
	Nspans   int
	Dtyp     int
	Dims     []float64
	DL,LL    float64
	Lspan    float64
	Ke, Le   float64 //ke/effective length
	Dnl, Dnr float64 //depth of notch at left and right
	Lbl, Rbl float64 //bearing length at left and right
	Dn, Nle  float64 //depth of notch, length of notch
	Tplnk    float64 //thickness of plank
	Brchk    bool
	Notch    bool
	Solid    bool
	Prp      kass.Wdprp `json:"Prp"`
	Lspans   []float64
	Clvrs    [][]float64
	Selfwt   bool
	Lclvr    bool
	Rclvr    bool
	Verbose  bool
	Clctyp   int
	Ldcases  [][]float64
	Drat     float64
	Wu       float64
	Mu       float64
	Vu       float64
	Dmax     float64
	Kostin   float64
	Kost     float64
	Nsecs    int
	Sec      kass.SectIn
	Rez      [][]float64
	Vals     [][]float64
	Dz       bool
	Spam     bool
	Report   string
}


//WdCol is a struct that stores timber column fields
//see chapter 9, abel o.o
type WdCol struct{
	Id       int
	Title    string
	Styp     int
	Grp      int
	Endc     int
	Code     int
	Clctyp   int
	Dtyp     int
	Selfwt   bool
	Tensile  bool
	Solid    bool //if solid, box type solid column
	Prp      kass.Wdprp `json:"Prp"`
	Dims     []float64
	Kn       []float64
	Lspan    float64
	Pu       float64
	Le       float64
	Tplnk    float64
	Dplnk    float64
	Nplnk    float64
	Ke       float64
	Kce      float64
	Cbi      float64
	Kfac     float64
	Kerst    float64 //end restraint factor (for spaced col/etc)
	Spc      float64
	K1, K2   float64
	K8, K9   float64
	K10      float64
	Srat     float64
	Bdrat    float64 //limiting breadth to depth ratio
	Bsdims   [][]float64
	Rez      [][]float64
	Vals     [][]float64
	Sec      kass.SectIn
	Dz       bool
	Spam     bool
	Verbose  bool
	Nsecs    int
	Kostin   float64
	Kost     float64
	Report   string
	Txtplots []string
	Bclass    int //buckling class; 1 - short, 2 - int, 3 - long

}


/*
   column end conditions (Endc) - abel page 165
   1. Both ends fixed, no side sway 0.5 , 0.7
   2. One end pinned the other fixed, no side sway 0.7, 0.85
   3. Both ends pinned, no side sway 1.0, 1.0
   4. Both ends fixed, side sway allowed 0.7?, 0.7
   5. One end fixed, the other free, side sway allowed 2.0
   6. One end fixed the other pinned, side sway allowed 0.7, 0.85

   
   column/beam sectypes:
   0.solid circle
   1.solid rect
   2.tapered solid circle - NEW
   3.tapered solid rect - NEW
   4.hollow circ tube
   5.hollow rect box
   6.t (built up)
   7.i (built up)
   8.plywood box beam
   26. spaced col

*/
