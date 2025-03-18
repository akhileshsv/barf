package barf

import(
	"log"
	"fmt"
	"math"
	"math/rand"
	kass"barf/kass"
)

//why define a separate CBm struct and functions this is the worst form of laziness

//CBm is a continuous beam struct
type CBm struct {
	Title        string
	Term         string
	Foldr        string
	Id           int
	DL           float64
	LL           float64
	Nspans       int
	Lspans       []float64
	Sections     [][]float64
	Sectypes     []int
	D1,D2        float64
	Fck          float64
	Fy           float64 
	Fyv          float64
	Nomcvr       float64
	Efcvr        float64
	Dslb         float64
	Code         int
	Bdx          []int
	DM           float64
	Clvrs        [][]float64
	Clmem, Crmem int
	Clsecs       [][]float64
	Flridx       int
	Web          bool
	Verbose      bool
	Dz           int
	Dvar         int
	Fop          int //optimize if > 0
	Allcons      int
	Csteel       int
	Bfcalc       bool
	Spam         bool
	Selfwt       bool
	Ldenvz       bool
	Lclvr, Rclvr bool
	Curtail      bool //curtail/detail span steel
	Noprnt       bool
	Width        float64 //if non empty, constant B
	Opt          int //1 - ga, 2 - pso
	Dconst       bool
	Rslb         bool //is ribbed slab
	Report       string
	Txtplots     []string `json:",omitempty"` 
	PSFs         []float64 `json:",omitempty"` 
	Msloads      [][]float64 `json:",omitempty"` 
	Jsloads      [][]float64 `json:",omitempty"` 
	Lsxs         []float64 `json:",omitempty"` 
	Cp           [][]float64 `json:",omitempty"` 
	Em           [][]float64 `json:",omitempty"`
	Bmvec        []int `json:",omitempty"`
	Quants       []float64 `json:",omitempty"`
	Base         *kass.Model `json:",omitempty"` //this should just be Mod 
	RcBm         [][]*RccBm `json:",omitempty"`
	BmEnv        map[int]*kass.BmEnv `json:",omitempty"`
	Kostin       []float64 `json:",omitempty"`
	Optpar       []float64 `json:",omitempty"` //optimization params
	Vec          []float64 `json:",omitempty"` //pso opt vec/zome
	Nlayers      int `json:",omitempty"`
	Kost         float64 `json:",omitempty"`
	Sec          kass.SectIn `json:",omitempty"`
}

//CBeamEnvRcc calculates shear and bending moment envelopes for a continuous beam
//hulse section 2.3
//holy hell rewrite this mess
//DAMMIT. just import mosh into bash and tmbr :(
func CBeamEnvRcc(cb *CBm, termstr string, report bool)(bmenv map[int]*kass.BmEnv, err error){
	/*
	   has to be broken down into funcs but '\(-o-)/
	   now with cantilevers, generates a cbeam given spans and sections
	   todo - add joint load envelopes
	*/
	//fmt.Println(ColorRed,"sections in->",cb.Sections,ColorReset)
	if cb.Title == ""{
		if cb.Id == 0{
			cb.Id = rand.Intn(666)
		}
		cb.Title = fmt.Sprintf("cbeam_%v",cb.Id)
	}
	cb.Quants = make([]float64, 3)
	if cb.Nspans == 0{cb.Nspans = len(cb.Lspans)}
	bmenv = make(map[int]*kass.BmEnv); bmvec := make([]int, cb.Nspans)
	var coords, em, cp [][]float64
	var mprp, msup [][]int
	var wtdl []float64
	var x float64
	//var sdl []float64
	//coord loop
	coords = append(coords, []float64{x})
	for i:=0; i<cb.Nspans; i++{
		if len(cb.Lspans) == 1{
			x += cb.Lspans[0]
		} else {
			x += cb.Lspans[i]
		}
		coords = append(coords, []float64{x})
	}
	if len(cb.Lspans) == 1 && cb.Nspans != 1{
		cb.Lspans = make([]float64, cb.Nspans)
		for i := 0; i < cb.Nspans; i++{
			cb.Lspans[i] = cb.Lspans[0]
		}
	}
	if len(cb.Em) == 0 {em = append(em, []float64{FckEm(cb.Fck)})} else {em = cb.Em}
	if len(cb.Sectypes) == 0{
		switch{
			case len(cb.Sections) == 1:
			cb.Sectypes = append(cb.Sectypes, 1)
			default:
			cb.Sectypes = make([]int, len(cb.Sections))
			for i := range cb.Sectypes{
				cb.Sectypes[i] = 1
			}
		}
	}
	if len(cb.Sections) == 0{
		err = fmt.Errorf("no sections specified-%v",cb.Sections)
		return
	}
	if len(cb.Lsxs) == 1{
		lsx := cb.Lsxs[0]
		cb.Lsxs = make([]float64,cb.Nspans+1)
		for i := range cb.Lsxs{
			cb.Lsxs[i] = lsx
		}
	}
	if cb.Efcvr == 0.0{
		if cb.Nomcvr == 0.0{
			cb.Nomcvr = 30.0
		}
		cb.Efcvr = cb.Nomcvr + 30.0/2.0 + 5.0
	}
	if cb.Code == 0{cb.Code = 1}
	if len(cb.PSFs) == 0{
		switch cb.Code{
			case 1:
			cb.PSFs = []float64{1.5,1.0,1.5,0.0}
			case 2:
			cb.PSFs = []float64{1.4,1.0,1.6,0.0}
		}
	}
	//fmt.Println(cb.PSFs)
	//calc breadth of flange (BAAH)
	switch cb.Bfcalc{
		case true:
		if cb.Dslb == 0.0{
			log.Println("ERRORE,errore-> depth of slab required for flange calc")
			err = ErrDim
			return
		}
		if len(cb.Bdx) == 0{
			cb.Bdx = make([]int, cb.Nspans)
			if len(cb.Sections) == 1{
				for i := range cb.Bdx{
					cb.Bdx[i] = 1
				}
			} else {
				for i := range cb.Bdx{
					cb.Bdx[i] = i+1
				}
			}
		}
		var mrel, bstyp int
		if cb.Nspans == 1 && (!cb.Lclvr && !cb.Rclvr){
			mrel = 1
		}
		cdx := []int{}
		secvec, sts := kass.CalcBf(cb.Code, bstyp, mrel, cb.Nspans, cb.Dslb, cdx, cb.Bdx, cb.Sectypes, cb.Lspans, cb.Sections)
		cb.Sections = make([][]float64, len(secvec))
		cb.Sectypes = make([]int, len(sts))
		for i := range secvec{
			cb.Sections[i] = make([]float64, len(secvec[i]))
			copy(cb.Sections[i],secvec[i])
			cb.Sectypes[i] = sts[i]
		}
	}
	if cb.Selfwt{
		wtdl = make([]float64, cb.Nspans)
	}
	//fmt.Println("after calc, sections->",cb.Sections)
	//fmt.Println("after calc, sectypes->",cb.Sectypes)
	var bsec, mrel int
	for i := 1; i <= cb.Nspans; i++ {
		bmvec[i-1] = i
		switch{
			case len(cb.Sections) == 1 || len(cb.Cp) == 1:
			bsec = 1
			case len(cb.Sections) == cb.Nspans:
			bsec = i
			case len(cb.Bdx) == cb.Nspans && len(cb.Sections) != cb.Nspans:
			//FALLS THRU SO SHOULD BE OKAY
			bsec = cb.Bdx[i-1]
			default:
			log.Println("ERRORE,errore->error in section vec")
			err = ErrDim
			return
		}
		mprp = append(mprp,[]int{i, i+1, 1, bsec, mrel})
	}
	for i := 1; i<= cb.Nspans +1; i++{
		msup = append(msup, []int{i,-1, 0})
	}
	if len(cb.Cp) == 0{
		if len(cb.Sections) == 0{
			log.Println("ERRORE,errore->specify sections")
			err = ErrDim
			return
		} else {
			for i, dims := range cb.Sections{
				var styp int
				if len(cb.Sectypes) == 1 {
					styp = cb.Sectypes[0]
				} else {
					styp = cb.Sectypes[i]
				}
				//"|***|" ALL DIMS IN MM. ALL DIMS IN MM. ALL DIMS IN MM...
				bar := kass.CalcSecProp(styp, dims)
				cp = append(cp, []float64{bar.Ixx*1e-12, bar.Area*1e-6})
				if cb.Selfwt{
					switch styp{
						case 1:
						wtdl[i] = dims[0] * (dims[1] - cb.Dslb) * 25.0 * 1e-6
						case 6, 7, 8, 9, 10:
						wtdl[i] = dims[2] * (dims[1] - cb.Dslb) * 25.0 * 1e-6
						default:
						//bar areas in m2
						wtdl[i] = cp[i][1] * 25.0
					}
				}
			}
		}
	} else {
		cb.Sections = make([][]float64, len(cb.Cp))
		cp = cb.Cp
	}
	//(change) what crap what absolute shite 
	if len(cb.Sections) == 1 && cb.Selfwt{
		wt := wtdl[0]
		for i := range wtdl{
			wtdl[i] = wt
		}
	}
	//for i := range wtdl{
	//	fmt.Println("span, dl->",i+1,wtdl[i])
	//}

	//build load patterns
	var nlp, nspans int
	var lclvr, rclvr bool
	var lcdl, lcll, rcdl, rcll float64
	nspans = cb.Nspans
	nlp = cb.Nspans + 1
	if len(cb.Clvrs) > 0 && cb.Clvrs[0][0]  > 0{
		//left cantilever OiNk
		lclvr = true; nlp++; nspans++
		if cb.Clvrs[0][1] + cb.Clvrs[0][2] > 0 {
			lcdl = cb.Clvrs[0][1]; lcll = cb.Clvrs[0][2]
		} else {
			lcdl = cb.DL; lcll = cb.LL
		}
		clen := cb.Clvrs[0][0]
		coords = append(coords, []float64{-clen})
		clsec := 1
		cb.Lclvr = true
		switch{
			case len(cb.Sections) > cb.Nspans:
			clsec = cb.Nspans+1
			case len(cb.Bdx) > cb.Nspans:
			clsec = cb.Bdx[cb.Nspans]
		}
		if cb.Selfwt{
			lcdl += cp[clsec-1][1] * 25.0
		}
		mprp = append(mprp,[]int{cb.Nspans+2,1,1,clsec,0})
		bmvec = append([]int{cb.Nspans+1},bmvec...)
		cb.Clmem = cb.Nspans + 1
	}
	if len(cb.Clvrs) > 0 && cb.Clvrs[1][0]  > 0{
		//right cantilever on
		rclvr = true; nlp++; nspans++; cb.Rclvr = true
		if cb.Clvrs[1][1] + cb.Clvrs[1][2] > 0 {
			rcdl = cb.Clvrs[1][1]; rcll = cb.Clvrs[1][2]
		} else {
			rcdl = cb.DL; rcll = cb.LL
		}
		clen := cb.Clvrs[1][0] + coords[cb.Nspans][0]
		coords = append(coords, []float64{clen})
		crsec := 1
		switch{ 
			case len(cb.Sections) > cb.Nspans + 1:
			crsec = cb.Nspans+2
			case len(cb.Bdx) > cb.Nspans+1:
			crsec = cb.Bdx[cb.Nspans]
		}
		if cb.Selfwt{
			rcdl += cp[crsec-1][1] * 25.0
		}
		mprp = append(mprp,[]int{cb.Nspans+1,cb.Nspans+3,1,crsec,0})
		bmvec = append(bmvec,cb.Nspans+2)
		cb.Crmem = cb.Nspans + 2
	}
	advloads := make(map[int][][]float64)
	benloads := make(map[int][][]float64)
	loadcons := make(map[int][][]float64)
	if len(cb.Msloads) != 0 {
		for _, ldcase := range cb.Msloads {
			ldcat := ldcase[6]; mem := int(ldcase[0])
			var w1a, w2a, w1b, w2b float64
			switch ldcat{
			case 1.0:
				w1a = cb.PSFs[0]*ldcase[2]
				w2a = cb.PSFs[0]*ldcase[3]
				w1b = cb.PSFs[1]*ldcase[2]
				w2b = cb.PSFs[1]*ldcase[3]
			case 2.0:
				w1a = cb.PSFs[2]*ldcase[2]
				w2a = cb.PSFs[2]*ldcase[3]
				w1b = cb.PSFs[3]*ldcase[2]
				w2b = cb.PSFs[3]*ldcase[3]
			}
			advloads[mem] = append(advloads[mem],[]float64{float64(mem),ldcase[1],w1a,w2a,ldcase[4],ldcase[5]})
			if w1b + w2b == 0.0 {continue}
			benloads[mem] = append(benloads[mem],[]float64{float64(mem),ldcase[1],w1b,w2b,ldcase[4],ldcase[5]})
		}
	}
	var wadl, wbdl, wall, wbll float64
	if cb.DL + cb.LL > 0.0{
		for i := 1; i <= cb.Nspans; i++ {
			dl := cb.DL
			if cb.Selfwt{
				dl += wtdl[i-1]
			}
			mem := i
			wadl = cb.PSFs[0] * dl
			wbdl = cb.PSFs[1] * dl
			wall = cb.PSFs[2] * cb.LL
			wbll = cb.PSFs[3] * cb.LL
			advloads[mem] = append(advloads[mem], []float64{float64(mem),3.0,wadl,0,0,0,1.0})
			advloads[mem] = append(advloads[mem], []float64{float64(mem),3.0,wall,0,0,0,2.0})
			benloads[mem] = append(benloads[mem], []float64{float64(mem),3.0,wbdl,0,0,0,1.0})
			if wbll > 0 {
				benloads[mem] = append(benloads[mem], []float64{float64(mem),3.0,wbll,0,0,0,2.0})
			}
			//if cb.Verbose{fmt.Printf("calc. loads-> wadl %.2f wall %.2f wbdl %.2f wbll %.2f dl %.2f ll %.2f\n",wadl, wall, wbdl, wbll, dl, cb.LL)}
		}	
		if lclvr{
			mem := cb.Nspans + 1
			wadl = cb.PSFs[0] * lcdl 
			wbdl = cb.PSFs[1] * lcdl 
			wall = cb.PSFs[2] * lcll 
			wbll = cb.PSFs[3] * lcll 
			advloads[mem] = append(advloads[mem], []float64{float64(mem),3.0,wadl,0,0,0,1.0})
			advloads[mem] = append(advloads[mem], []float64{float64(mem),3.0,wall,0,0,0,2.0})
			benloads[mem] = append(benloads[mem], []float64{float64(mem),3.0,wbdl,0,0,0,1.0})
			if wbll > 0 {
				benloads[mem] = append(benloads[mem], []float64{float64(mem),3.0,wbll,0,0,0,2.0})
			}
		}
		if rclvr{
			mem := cb.Nspans + 2
			wadl = cb.PSFs[0] * rcdl 
			wbdl = cb.PSFs[1] * rcdl 
			wall = cb.PSFs[2] * rcll 
			wbll = cb.PSFs[3] * rcll 
			advloads[mem] = append(advloads[mem], []float64{float64(mem),3.0,wadl,0,0,0,1.0})
			advloads[mem] = append(advloads[mem], []float64{float64(mem),3.0,wall,0,0,0,2.0})
			benloads[mem] = append(benloads[mem], []float64{float64(mem),3.0,wbdl,0,0,0,1.0})
			if wbll > 0 {
				benloads[mem] = append(benloads[mem], []float64{float64(mem),3.0,wbll,0,0,0,2.0})
			}
		}
	}
	for i := 1; i <= cb.Nspans; i++ {
		mem := i
		if i % 2 == 0 {
			loadcons[0] = append(loadcons[0], advloads[mem]...)
			loadcons[1] = append(loadcons[1], benloads[mem]...)
			loadcons[2] = append(loadcons[2], advloads[mem]...)
		} else {
			loadcons[0] = append(loadcons[0], advloads[mem]...)
			loadcons[1] = append(loadcons[1], advloads[mem]...)
			loadcons[2] = append(loadcons[2], benloads[mem]...)
		}
	}
	for i := 1; i <= cb.Nspans - 1; i++ {
		lp := i + 2
		mind := i
		for mem := range advloads {
			//if mem == -1 || mem == -2 {continue}
			if mem == mind || mem == mind + 1 {
				loadcons[lp] = append(loadcons[lp], advloads[mem]...)
			} else {
				loadcons[lp] = append(loadcons[lp], benloads[mem]...)
			}
		}
	}
	if lclvr {
		clmem := cb.Nspans+1
		for lp := 0; lp <= 2; lp++{
			switch lp{
				case 0,2:
				loadcons[lp] = append(loadcons[lp],advloads[clmem]...)
				case 1:
				loadcons[lp] = append(loadcons[lp],benloads[clmem]...)
			}
		} 
		for mem := range benloads {
			switch mem{
				case cb.Nspans+1, 1:
				loadcons[cb.Nspans+2] = append(loadcons[cb.Nspans+2], advloads[mem]...)
				default:
				loadcons[cb.Nspans+2] = append(loadcons[cb.Nspans+2], benloads[mem]...)
			}
		}
	}
	if rclvr {
		crmem := cb.Nspans+2
		for lp := 0; lp <= 2; lp++{
			switch lp{
				case 0,1:
				loadcons[lp] = append(loadcons[lp],advloads[crmem]...)
				case 2:
				loadcons[lp] = append(loadcons[lp],benloads[crmem]...)
			}
		} 
		for mem := range benloads {
			switch mem{
				case cb.Nspans+2, cb.Nspans:
				loadcons[cb.Nspans+3] = append(loadcons[cb.Nspans+3], advloads[mem]...)
				default:
				loadcons[cb.Nspans+3] = append(loadcons[cb.Nspans+3], benloads[mem]...)
			}	
		}
	}
	endc := 2; ldxb := 2; rdxb := 2
	if cb.Nspans == 1 && !cb.Lclvr && !cb.Rclvr{
		endc = 1
		ldxb = 1; rdxb = 1
	}
	//now demonstrating the use of one forgotten dupicate variable
	cb.RcBm = make([][]*RccBm,nspans)
	for i := 1; i <= cb.Nspans; i++{
		ldx := ldxb; rdx := rdxb
		if i == 1{ldx = 1}
		if i == cb.Nspans{rdx = 1}
		jb, je := mprp[i-1][0], mprp[i-1][1]
		c1, c2 := coords[jb-1], coords[je-1]
		cp := mprp[i-1][3]
		s := mprp[i-1][3]
		var lsx, rsx float64
		if len(cb.Lsxs) == cb.Nspans + 1{lsx, rsx = cb.Lsxs[i-1], cb.Lsxs[i]}
		bmenv[i] = &kass.BmEnv{
			Id:i,
			EnvRez:make(map[int]kass.BeamRez),
			Venv:make([]float64,21),
			Mpenv:make([]float64,21),
			Mnenv:make([]float64,21),
			Dims:cb.Sections[s-1],
			Coords:[][]float64{c1,c2},
			Lsx:lsx/1000.0,Rsx:rsx/1000.0,
			Vrd:make([]float64,21),
			Mnrd:make([]float64,21),
			Mprd:make([]float64,21),
			Endc:endc,
			Ldx:ldx,Rdx:rdx,
			Rslb:cb.Rslb,
			Frmtyp:1,
			Title:cb.Title,
			Kostin:cb.Kostin,
		}
		cb.RcBm[i-1] = make([]*RccBm,3)
		if len(cb.Sectypes) < cp{
			log.Println("ERRORE,errore->error in sectype vec")
			err = ErrDim
			return
		}
		//bstyp := cb.Sectypes[cp-1]
		GetBmArr(cb.RcBm[i-1], bmenv[i], cb.Kostin,cb.Fck, cb.Fy, cb.Fyv, cb.Efcvr, cb.DM, cb.D1, cb.D2, cb.Dslb, cb.Code, cb.Sectypes[cp-1], cb.Verbose)
	}
	if lclvr{
		i := cb.Nspans + 1
		jb, je := mprp[i-1][0], mprp[i-1][1]
		c1, c2 := coords[jb-1], coords[je-1]
		cp := mprp[i-1][3]
		s := mprp[i-1][3]
		ldx := 0; rdx := 2
		var lsx, rsx float64
		if len(cb.Lsxs) != 0{
			rsx = cb.Lsxs[0]
		}
		bmenv[i] = &kass.BmEnv{
			Id:i,
			EnvRez:make(map[int]kass.BeamRez),
			Venv:make([]float64,21),
			Mpenv:make([]float64,21),
			Mnenv:make([]float64,21),
			Dims:cb.Sections[s-1],
			Coords:[][]float64{c1,c2},
			Lsx:lsx/1000.0,Rsx:rsx/1000.0,
			Vrd:make([]float64,21),
			Mnrd:make([]float64,21),
			Mprd:make([]float64,21),
			Endc:0,
			Ldx:ldx,Rdx:rdx,
			Rslb:cb.Rslb,
		}
		//even cantilevers have three sections due to human laziness
		cb.RcBm[i-1] = make([]*RccBm,3)
		if len(cb.Sectypes) < cp{
			log.Println("ERRORE,errore->error in sectype vec")
			err = ErrDim
			return
		}
		GetBmArr(cb.RcBm[i-1], bmenv[i], cb.Kostin, cb.Fck, cb.Fy, cb.Fyv, cb.Efcvr, cb.DM, cb.D1, cb.D2, cb.Dslb, cb.Code, cb.Sectypes[cp-1], cb.Verbose)
	}
	if rclvr{
		i := cb.Nspans + 2
		jb, je := mprp[i-1][0], mprp[i-1][1]
		c1, c2 := coords[jb-1], coords[je-1]
		cp := mprp[i-1][3]
		s := mprp[i-1][3]
		ldx := 2; rdx := 0
		var lsx, rsx float64
		if len(cb.Lsxs) != 0{
			lsx = cb.Lsxs[len(cb.Lsxs)-1]
		}
		bmenv[i] = &kass.BmEnv{
			Id:i,
			EnvRez:make(map[int]kass.BeamRez),
			Venv:make([]float64,21),
			Mpenv:make([]float64,21),
			Mnenv:make([]float64,21),
			Dims:cb.Sections[s-1],
			Coords:[][]float64{c1,c2},
			Lsx:lsx/1000.0,Rsx:rsx/1000.0,
			Vrd:make([]float64,21),
			Mnrd:make([]float64,21),
			Mprd:make([]float64,21),
			Endc:0,
			Ldx:ldx, Rdx:rdx,
			Rslb:cb.Rslb,
		}
		cb.RcBm[i-1] = make([]*RccBm,3)
		if len(cb.Sectypes) < cp{
			log.Println("ERRORE,errore->error in sectype vec")
			err = ErrDim
			return
		}
		GetBmArr(cb.RcBm[i-1], bmenv[i], cb.Kostin, cb.Fck, cb.Fy, cb.Fyv, cb.Efcvr, cb.DM, cb.D1, cb.D2, cb.Dslb, cb.Code, cb.Sectypes[cp-1], cb.Verbose)
	}
	var modld *kass.Model	
	mod := &kass.Model{
		Ncjt:2,
		Cmdz:[]string{"1db","mks","1"},
		Frmstr:"1db",
		Units:"knm",
		Coords: coords,
		Supports: msup,
		Em: em,       
		Cp: cp,       
		Mprp:mprp,
	}
	ms0 := make(map[int]*kass.Mem)
	mslmap := make(map[int]map[int][][]float64)
	for lp, ldcons := range loadcons{
		switch lp{
			case 0:
			mod.Msloads = ldcons
		}
		modld = &kass.Model{
			Ncjt:2,
			Coords:coords,
			Supports:msup,
			Em:em,
			Cp:cp,
			Mprp:mprp,
			Msloads:ldcons,
			Frmstr:"1db",
			Units:"knm",
			Noprnt:cb.Noprnt,
			Web:cb.Web,
		}
		frmrez, e := kass.CalcBm1d(modld, 2)
		if e != nil{
			err = e
			return
		}
		ms,_ := frmrez[1].(map[int]*kass.Mem)
		if lp == 0 {ms0 = ms}
		msloaded, _ := frmrez[5].(map[int][][]float64)
		spanchn := make(chan kass.BeamRez,len(msloaded))
		mslmap[lp] = msloaded
		//cbeam.go.txt this is non parallel (wtf does this even mean)
		for id, ldcase := range msloaded{
			go kass.BeamFrc(2, id, ms[id], ldcase, spanchn, false)
		}
		for _ = range msloaded{
			r := <- spanchn
			id := r.Mem
			bm := bmenv[id]
			bm.EnvRez[lp] = r
			if len(bm.Xs) == 0 {
				bm.Xs = r.Xs
			}
			xdiv := ms[id].Geoms[0]/20.0
			lsx := bm.Lsx; rsx := ms[id].Geoms[0] - bm.Rsx
			
			//CEILING OR FLOOR?
			//il := int(math.Ceil(lsx/xdiv)); ir := int(math.Ceil(rsx/xdiv))

			il := int(math.Floor(lsx/xdiv)); ir := int(math.Floor(rsx/xdiv))
			//fmt.Println("beam-",id,"lsx-",lsx,"rsx-",rsx,"il-",il,"ir-",ir)
			var vl, vr, ml, mr float64
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
		}
	}
	if cb.DM > 0.0{
		CBeamDM(2, bmvec, bmenv, cb.DM, ms0, mslmap)
	}
	cb.Bmvec = make([]int, len(bmvec))
	copy(cb.Bmvec, bmvec)
	if cb.Web{cb.Foldr = "web"}
	if termstr != ""{
		//PLOT ALL LOAD PATTERNS - add new bool for this - doing it anyway
		// if termstr == "dumb"{
		plotchn := make(chan []interface{}, cb.Nspans+4)
		for i := 0; i <= nlp; i++{
			go PlotLp(i, bmenv, ms0, mslmap, bmvec, termstr, cb.Title, cb.Foldr,plotchn)
		}
		for i := 0; i<=nlp; i++{
			rez := <- plotchn
			//lp, _ := rez[0].(int)
			txtplot, _ := rez[1].(string)
			cb.Txtplots = append(cb.Txtplots,txtplot)			
			
		}
		txtplot := PlotBmEnv(bmenv, bmvec, termstr, cb.Title, cb.Foldr)
		//txtplot += fmt.Sprintf("\n%s\n",cb.Title)
		cb.Txtplots = append(cb.Txtplots, txtplot)
		if cb.DM > 0.0{
			rdtxt := PlotBmRdEnv(bmenv, bmvec, termstr, cb.Title, cb.Foldr)
			cb.Txtplots = append(cb.Txtplots, rdtxt)
		}
		// fmt.Println("finished plotting load cases->")
		// for i, txtplot := range cb.Txtplots{
		// 	fmt.Println("plot no.-",i+1,"plot data/loc->\n",txtplot)
		// }
	}
	err = nil
	return
}

/*
   

   YEOLDE
		for i := 1; i <= cb.Nspans; i++{
			//add REPORT
			if 1 == 0{
				fmt.Println(ColorReset)
				fmt.Println("member-",i)
				
				for j, vx := range bmenv[i].Venv{
					fmt.Println(ColorCyan)
					fmt.Printf("elastic\t- %v shear %.1f hogging BM %.1f sagging BM %.1f\n",j, vx, bmenv[i].Mnenv[j], bmenv[i].Mpenv[j])
					fmt.Println(ColorYellow)
					fmt.Printf("redistributed\t- %v shear %.1f hogging BM %.1f sagging BM %.1f\n",j, bmenv[i].Vrd[j], bmenv[i].Mnrd[j], bmenv[i].Mprd[j])
					fmt.Println(ColorReset)
				}
				fmt.Println("mpmax, mnmax, vmax->",bmenv[i].Mpmax,bmenv[i].Mnmax,bmenv[i].Vmax)
				fmt.Println("elastic ml, mr, vl, vr->",bmenv[i].Ml,bmenv[i].Mr,bmenv[i].Vl,bmenv[i].Vr)
				fmt.Println("redistribute ml, mr, vl, vr->",bmenv[i].Mlrd,bmenv[i].Mrrd,bmenv[i].Vlrd,bmenv[i].Vrrd)
			}
		}	
	}


*/
