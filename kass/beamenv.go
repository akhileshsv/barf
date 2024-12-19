package barf

import (
	"log"
	"fmt"
	"math"
)

//BmEnv contains fields for member results for design funcs in other dirs
//env rez maps load patterns to individual beam rezs (sf, bm at 20 intervals)
//venv, mnenv, mpenv - shear, sagging and hogging moment envelopes

type BmEnv struct {
	Id                    int
	Cp                    int
	Title                 string
	Foldr                 string
	Kostin                []float64 //adding this here seems easier
	EnvRez                map[int]BeamRez
	Venv, Mnenv, Mpenv    []float64
	Vmax, Mnmax, Mpmax    float64
	Vmaxx, Mnmaxx, Mpmaxx float64
	Xs                    []float64
	Dims                  []float64
	Coords                [][]float64
	Vl, Vr, Ml, Mr        float64
	Vls, Vrs              float64 
	Vlrs, Vrrs            float64 
	Lsx, Rsx              float64
	Hc, Dh                float64
	Drx, Dry, Drd         float64
	DMl, DMr              float64
	Vrd, Mnrd, Mprd       []float64
	Vlrd, Vrrd            float64
	Mlrd, Mrrd            float64
	Vrmax                 float64
	Mnrmax, Mprmax        float64
	Vrmaxx                float64
	Mnrmaxx, Mprmaxx      float64
	Lspan                 float64
	Wspan                 float64 //wspan
	K                     float64 //ei/l
	Endc                  int
	Ldx, Rdx              int //continuity condition
	Qfrez                 map[int][]float64
	Pbmax, Pemax          float64
	Pumax                 float64
	Mbmax, Memax          float64
	Mtyp                  int
	Styp                  int
	Tmax, Cmax            float64 //max tension n kompresshun
	Nmax                  float64 //max design load/unit area (for flat slabs)
	DM                    float64
	Spandx                int //-1 - left, -2 - right end, 0 - interior span
	Frmtyp                int //1 - cbeam, 2 - t2d, 3-f2d, 4-t3d, 5-gr3d, 6-f3d
	Roof, Plinth          bool
	Endrel                bool
	Ignore                bool
	Tweak                 bool
	Rslb                  bool
	Web                   bool
	Term                  string
}

//ColEnv stores column results over multiple load cases
type ColEnv struct {
	//env rez maps load cases to qf
	Id           int
	Title        string
	Dims         []float64
	EnvRez       map[int][]float64
	Mtmax, Mbmax float64
	Ptmax, Pbmax float64
	Mtmay, Mbmay float64
	Pumax        float64
	Ptmay, Pbmay float64
	Pumay        float64
	Pufac        float64
	L0y, Ley     float64
	L0, Le       float64
	Lspan        float64
	Coords       [][]float64
	Lb, Ub       []int //lower beams (l, r); upper beams (l, r)
	Lbd, Ubd     [][]float64
	Lst, Ust     []int
	Lbk, Ubk     []float64
	Lbss, Ubss   bool
	Ljbase       bool
	Fixbase      bool
	Braced       bool
	Tweak        bool
	Styp         int
	Kcol         float64
	B1, B2       float64
	Lby, Uby     []int //lower beams (l, r); upper beams (l, r)
	Lbdy, Ubdy   [][]float64
	Lsty, Usty   []int
	Lbky, Ubky   float64
	Lbsy, Ubsy   bool
	Slender      bool
	Slendery     bool
	Ignore       bool
	Web          bool
	Stair        bool
	Term         string
	Foldr        string //why isn't it FOLDER
}

//SortBmEnv sorts a beamrez given a mem for max vals and updates ofc
func (bm *BmEnv) SortBmEnv(mem *Mem, r BeamRez){
	if len(bm.Xs) == 0 {
		bm.Xs = r.Xs
	}
	//teehee sort for col env
	if math.Abs(mem.Qf[2]) > math.Abs(bm.Mbmax) {bm.Mbmax = mem.Qf[2]}
	if math.Abs(mem.Qf[5]) > math.Abs(bm.Memax) {bm.Memax = mem.Qf[5]}
	if math.Abs(mem.Qf[0]) > math.Abs(bm.Pbmax) {bm.Pbmax = mem.Qf[0]}
	if math.Abs(mem.Qf[3]) > math.Abs(bm.Pemax) {bm.Pemax = mem.Qf[3]}

	//sort for bm env
	xdiv := r.Xs[1] - r.Xs[0]
	lsx := bm.Lsx; rsx := r.Xs[20] - bm.Rsx
	il := int(math.Ceil(lsx/xdiv)); ir := int(math.Ceil(rsx/xdiv))
	var vl, vr, ml, mr float64
	if math.Abs(r.SF[0]) > bm.Vls{
		bm.Vls = math.Abs(r.SF[0])
	}
	if math.Abs(r.SF[20]) > bm.Vrs{
		bm.Vrs = math.Abs(r.SF[20])
	}
	for i, vx := range r.SF{
		x := r.Xs[i]
		if i == il{
			switch{
				case i == 0:
				vl = vx
				ml = r.BM[i]					
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
	return
}

//CalcBmEnv generates bm and sf envelopes for a continuous beam (hulse section 2.3)		
//calcs bending moment and shear force envelopes for a model
//dead and live load envelopes
//bs 1.4 dl + 1.6 ll adverse, 1.0 dl beneficial
//loadcons : adv odd spans, ben even / vice versa; adv on span, span + 1
//UNUSED (test written for hulse ex 2.3 so it stays )

func CalcBmEnv(mod *Model) (loadenvz map[int]*BeamRez, bmenv map[int]*BmEnv){
	nspans := len(mod.Mprp)
	advloads := make(map[int][][]float64)
	benloads := make(map[int][][]float64)
	loadcons := make(map[int][][]float64)
	loadenvz = make(map[int]*BeamRez)	
	bmenv = make(map[int]*BmEnv)
	for i := 1; i <= nspans; i++ {
		loadenvz[i] = &BeamRez{
			Mem:i,
			Venv:make([]float64,21),
			Mnenv:make([]float64,21),
			Mpenv:make([]float64,21),
		}
		bmenv[i] = &BmEnv{
			Id:i,
			EnvRez:make(map[int]BeamRez),
			Venv:make([]float64,21),
			Mpenv:make([]float64,21),
			Mnenv:make([]float64,21),
			//Dims:cb.Sections[s-1],
			//Coords:[][]float64{c1,c2},
			//Lsx:lsx/1000.0,Rsx:rsx/1000.0,
		}
	}
	
	//var lcvlr, rclvr bool
	if len(mod.Clvrs[0]) > 2 && mod.Clvrs[0][1] + mod.Clvrs[0][2] > 0.0 {
		//add left clvr to model
		cldl := mod.Clvrs[0][1]; clll := mod.Clvrs[0][2]
		
		advloads[-1] = append(advloads[-1],[]float64{1.0,0.0,cldl*mod.PSFs[0]+clll*mod.PSFs[2]})
		benloads[-1] = append(benloads[-1],[]float64{1.0,0.0,cldl*mod.PSFs[1]+clll*mod.PSFs[3]})
	}
	if len(mod.Clvrs[1]) > 2 && mod.Clvrs[1][1] + mod.Clvrs[1][2] > 0.0 {
		crdl := mod.Clvrs[1][0]; crll := mod.Clvrs[1][1]
		advloads[-2] = append(advloads[-2],[]float64{float64(len(mod.Coords)),0.0,crdl*mod.PSFs[0]+crll*mod.PSFs[2]})
		benloads[-2] = append(benloads[-2],[]float64{float64(len(mod.Coords)),0.0,crdl*mod.PSFs[1]+crll*mod.PSFs[3]})
	}
	for _, ldcase := range mod.Msloads {
		ldcat := ldcase[6]; mem := int(ldcase[0])
		var w1a, w2a, w1b, w2b float64
		switch ldcat {
		case 1.0:
			w1a = mod.PSFs[0]*ldcase[2]
			w2a = mod.PSFs[0]*ldcase[3]
			w1b = mod.PSFs[1]*ldcase[2]
			w2b = mod.PSFs[1]*ldcase[3]
		case 2.0:
			w1a = mod.PSFs[2]*ldcase[2]
			w2a = mod.PSFs[2]*ldcase[3]
			w1b = mod.PSFs[3]*ldcase[2]
			w2b = mod.PSFs[3]*ldcase[3]
		}
		advloads[mem] = append(advloads[mem],[]float64{ldcase[0],ldcase[1],w1a,w2a,ldcase[4],ldcase[5]})
		if w1b + w2b == 0.0 {continue}
		benloads[mem] = append(benloads[mem],[]float64{ldcase[0],ldcase[1],w1b,w2b,ldcase[4],ldcase[5]})
	}
	for i := 1; i <= nspans; i++ {
		if i % 2 == 0 {
			loadcons[1] = append(loadcons[1], benloads[i]...)
			loadcons[2] = append(loadcons[2], advloads[i]...)
		} else {
			loadcons[1] = append(loadcons[1], advloads[i]...)
			loadcons[2] = append(loadcons[2], benloads[i]...)
		}
	}
	for i := 1; i <= nspans - 1; i++ {
		cind := i + 2
		for mem := range advloads {
			if mem == -1 || mem == -2 {continue}
			if mem == i || mem == i + 1 {
				loadcons[cind] = append(loadcons[cind], advloads[mem]...)
			} else {
				loadcons[cind] = append(loadcons[cind], benloads[mem]...)
			}
		}
	}
	for lp, ldcons := range loadcons {
		//fmt.Println("Loadcase-->",ldidx)
		modld := &Model{
			Ncjt:2,
			Coords:mod.Coords,
			Supports:mod.Supports,
			Em:mod.Em,
			Cp:mod.Cp,
			Mprp:mod.Mprp,
			Jloads:mod.Jloads,
			Msloads:ldcons,
		}
		frmrez, _ := CalcBm1d(modld, 2)
		ms,_ := frmrez[1].(map[int]*Mem)
		msloaded, _ := frmrez[5].(map[int][][]float64)
		spanchn := make(chan BeamRez,len(msloaded))
		//bmresults := make(map[int]*BeamRez)
		//var mem *BeamRez
		for member, ldcs := range msloaded {
			go BeamFrc(mod.Ncjt, member, ms[member], ldcs, spanchn, false)
			r := <- spanchn
			mem := loadenvz[member]
			bm := bmenv[member]
			bm.EnvRez[lp] = r
			for i, vx := range r.SF {
				if math.Abs(mem.Venv[i]) < math.Abs(vx) {
					mem.Venv[i] = vx
				}
				if math.Abs(mem.Mnenv[i]) < math.Abs(r.BM[i]) && r.BM[i] < 0.0 {
					mem.Mnenv[i] = r.BM[i]
				}
				if r.BM[i] > 0.0 && math.Abs(mem.Mpenv[i]) < math.Abs(r.BM[i]) {
					mem.Mpenv[i] = r.BM[i]
				}
			}
		}
	}
	return
}

//CalcBmSf calculates bending moment and shear force for spans of a continuous beam 
//hulse ex. 2.2

func CalcBmSf(mod *Model,frmrez []interface{},plotbm bool) (map[int]BeamRez){
	ncjt := mod.Ncjt
	ms,_ := frmrez[1].(map[int]*Mem)
	msloaded, _ := frmrez[5].(map[int][][]float64)
	spanchn := make(chan BeamRez,len(msloaded))
	bmresults := make(map[int]BeamRez)
	for member, ldcases := range msloaded {
		go BeamFrc(ncjt, member, ms[member], ldcases, spanchn,plotbm)
	}
	for i :=0; i < len(msloaded); i++ {
		r := <- spanchn
		bmresults[r.Mem] = r
	}
	if plotbm{
		for i:= 1; i <= len(bmresults); i++{
			bm := bmresults[i]
			//bm.TxtPlot = PlotBmSfBm(xs, vxs, mxs, dxs, l, true)
			//fmt.Println(ColorYellow,"member--",bm.Mem,ColorReset)
			for i, vx := range bm.SF {
				fmt.Println(ColorCyan)
				fmt.Printf("Div %d SF %.2f KN BM %.3f KN-m Def %.3f mm", i,vx, bm.BM[i],1000*bm.Dxs[i])
				log.Println("SPAN ",i,"SF ",vx, " Kn BM ", math.Ceil(bm.BM[i]*100)/100, " Kn-M DEF ", 1000*bm.Dxs[i]," mm")
			}
			fmt.Println(ColorPurple)
			fmt.Printf("Max SF %.3f at %.3f \nMax BM %.3f at %.3f\nMax def %.3f at %.3f",bm.Maxs[0],bm.Locs[0],bm.Maxs[1],bm.Locs[1],bm.Maxs[2],bm.Locs[2])
			fmt.Println(ColorGreen)
			for i, cfx := range bm.Cfxs{
				fmt.Println("counter flexsures",i, cfx)
			}
			fmt.Println(ColorReset)
			fmt.Println(bm.Txtplot)
		}
	}
	return bmresults
}

//CalcModBmSf is the entry func for bending moment and shear force calcs as in hulse ex. 2.2
//from menu/flags
func CalcModBmSf(mod *Model)(err error){
	//var rezstring string
	var frmrez []interface{}
	frmrez, err = CalcBm1d(mod, 2)
	if err != nil{
		return
	}
	_ = CalcBmSf(mod, frmrez,true)
	return
}
