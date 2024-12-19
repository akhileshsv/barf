package barf


//bs code beam design funcs (hulse/mosley)

import (
	"os"
	//"log"
	"fmt"
	"math"
	"time"
	"errors"
	"math/rand"
	"runtime"
	"encoding/json"
	"path/filepath"
	"io/ioutil"
	kass"barf/kass"
)

//RccBm is a struct to store rcc beam design values
//see hulse/mosley chapter 3
type RccBm struct{
	Id                int
	Mid               int
	Title             string
	Term              string
	Foldr             string
	Fck               float64
	Fy                float64
	Fyv               float64
	Ast               float64
	Asc               float64
	Nomcvr            float64
	Cvrt              float64
	Cvrc              float64
	Bw, Bf, Df, D     float64
	Mu, Vu            float64 //design moment, design shear
	Dused             float64
	Typ,Endc          int
	Code              int
	DM                float64
	Lspan             float64
	D1, D2            float64
	Aside             float64 //side face reinforcement for deep beams
	Rbrside           []float64 //side face rebar opt
	Styp              int
	Nrows             int //if nrows == 1; only single layer of reinf per level
	//Barcat            int //if 0 -  only std (sail) dias, else govindrajan dias
	Lsx, Rsx          float64
	Dims              []float64
	Dslb              float64 `json:",omitempty"`
	Web               bool `json:",omitempty"` 
	Monolith          bool `json:",omitempty"` //causes chimp tribe water wars//monolithic support connection
	Verbose           bool `json:",omitempty"`
	Npsec             bool `json:",omitempty"`
	Shrdz             bool `json:",omitempty"`
	Ismid             bool `json:",omitempty"`
	Curtail           bool `json:",omitempty"` //calc curtailment/detail
	Dconst            bool `json:",omitempty"`
	Ignore            bool `json:",omitempty"` //haha then why have this
	Dsgn              bool `json:",omitempty"`
	Blurb             bool `json:",omitempty"`
	Ibent             bool `json:",omitempty"` //why bend bars it seems stressful
	Tweak             bool `json:",omitempty"`
	Rslb              bool `json:",omitempty"` //is ribbed slab
	Nlayers           int `json:",omitempty"` //n layers of reinforcement
	Csteel            bool `json:",omitempty"` //has compression steel
	Mps               []float64 `json:",omitempty"`
	Mns               []float64 `json:",omitempty"`
	Vs                []float64 `json:",omitempty"`
	Xs                []float64 `json:",omitempty"`
	Cfxs              []float64 `json:",omitempty"`
	Cs                []int     `json:",omitempty"`
	L0                float64 `json:",omitempty"`
	Lbd               float64 `json:",omitempty"`
	Flip              bool `json:",omitempty"` 
	Xu                float64 `json:",omitempty"` 
	Nlegs             int `json:",omitempty"` 
	Dlink             float64 `json:",omitempty"` 
	Nomlink           bool `json:",omitempty"` 
	Sec               *kass.SectIn `json:",omitempty"` 
	L1, L2, L3, L4    float64 `json:",omitempty"` //link distances from x0
	Nlx               []float64 `json:",omitempty"` //no. of main min and nominal links
	S1, S2, S3        float64 `json:",omitempty"` //span distances from x0
	CL                []float64 `json:",omitempty"` //curtailment lengths from x0
	LDs               []float64   `json:",omitempty"` //dev lengths - top left, top right, bottom left, bottom right
	Lspc              []float64 `json:",omitempty"`
	Txtplot           []string `json:",omitempty"` //BAAH DELETE THIS 
	Txtplots          []string `json:",omitempty"` 
	Rbrt              []float64 `json:",omitempty"` 
	Rbrtopt           [][]float64 `json:",omitempty"`
	Rbrc              []float64 `json:",omitempty"`
	Rbrcopt           [][]float64 `json:",omitempty"`
	Nbent             float64 `json:",omitempty"`
	Diabent           float64 `json:",omitempty"`
	Tyb               float64 `json:",omitempty"`
	Dias              []float64 `json:",omitempty"`
	Barpts            [][]float64 `json:",omitempty"`
	Dbars             []float64 `json:",omitempty"`
	Tds               [][]float64 `json:",omitempty"` //WHAT IS THIS
	Cds               [][]float64 `json:",omitempty"`
	Dz                bool `json:",omitempty"`
	Hdim, Ldim        [][]float64 `json:",omitempty"`
	Nr, Sdx           int `json:",omitempty"`
	Avr               [][]float64 `json:",omitempty"` //torsion links and add reinf
	Svr               []float64 `json:",omitempty"` //torsion link max spacing
	Report            string `json:",omitempty"`
	Ldx, Rdx          int `json:",omitempty"` //ldx - 0,1,2 - free/simple/continuous
	Vrcc, Wstl, Afw   float64 `json:",omitempty"`
	Vtot              float64 `json:",omitempty"`
	Kost              float64 `json:",omitempty"`
	N1, N2            float64 `json:",omitempty"` //ast bar numbers
	Dia1, Dia2        float64 `json:",omitempty"` //ast bar dias
	N3, N4            float64 `json:",omitempty"` //asc bar numbers
	Dia3, Dia4        float64 `json:",omitempty"` //asc bar dias
	Kostin            []float64 `json:",omitempty"`
	Kunit             float64 `json:",omitempty"`
}

//InitDims initializes beam geom fields from beam dims for basic section types
func (b *RccBm) InitDims(){
	//again for the n basic holy types
	var dims []float64
	switch b.Styp{
		case 1:
		dims = []float64{b.Bw, b.Dused}
		b.Dims = make([]float64, len(dims))
		copy(b.Dims, dims)
		return
		case 6,7,8,9,10,14:
		dims = []float64{b.Bf, b.Dused, b.Bw, b.Df}
		b.Dims = make([]float64, len(dims))
		copy(b.Dims, dims)
		return
	}
	
}

//Init initalizes an RccBm struct
func (b *RccBm) Init() (err error){
	if b.Id == 0 && b.Mid == 0{
		b.Id = rand.Intn(666)
	}
	
	if b.Fck == 0.0{
		switch b.Code{
			case 1:
			b.Fck = 25.0
			case 2:
			b.Fck = 30.0
		}
	}
	if b.Fy == 0.0{
		switch b.Code{
			case 1:
			b.Fy = 415.0
			case 2:
			b.Fy = 460.0
		}
	}
	
	if b.Fyv == 0.0{
		b.Fyv = b.Fy
	}
	if b.Styp == 0{
		switch b.Tyb{
			case 0.0:
			b.Styp = 1
			case 0.5:
			b.Styp = 7
			case 1.0:
			b.Styp = 6
		}
	}
	b.InitDims()
	switch b.Styp{
		case 1, 6, 7, 8, 9, 10, 14:
		//either rect or flanged sections
		s := kass.SecGen(b.Styp, b.Dims)
		b.Sec = &s
		default:
		b.Npsec = true
	}
	switch b.Npsec{
		case true:
		s := kass.SecGen(b.Styp, b.Dims)
		b.Sec = &s
		b.Dused = b.Sec.Ymx
	}
	if b.Cvrc == 0.0 && b.Cvrt == 0.0{
		b.Cvrc = 45.0; b.Cvrt = 45.0
	}
	if b.Dlink == 0.0{b.Dlink = 8.0}
	return
}

//BmDesign is the general (-_-)7 beam section design entry func
func BmDesign(b *RccBm) (err error){
	b.Init()
	switch b.Styp{
		case 1,6,7,8,9,10,14:
		switch b.Code{
			case 2:
			err = BmStlBs(b, b.Mu)
			if b.D1 + b.D2 > 0.0{
				//fmt.Println("resetting beam")
				b.D1 = 0.0
				b.D2 = 0.0
				err = BmStlBs(b, b.Mu)
			}
			if err != nil{
				return
			}
			case 1:
			err, _ = BmDIs(b, b.Mu)
			
			if b.D1 + b.D2 > 0.0{
				
				//fmt.Println("resetting beam")
				b.D1 = 0.0
				b.D2 = 0.0
				err, _ = BmDIs(b, b.Mu)
			}
			
			if err != nil{
				return
			}
			//fmt.Println("ERRORE->",err)
		}
		default:
		//lmao 
		err = errors.New("npsec func not written (yet) :-|")
		return
	}
	if err != nil{
		return
	}
	if b.Ast + b.Asc <= 0.0{
		err = errors.New("steel calculation error")
		return
	}
	
	if b.Dused > 750.0{
		//deep beam
		b.RbrSide()
		
	}
	//fmt.Println(ColorRed,"beam->",b.Id, b.Mid,"calc steel-> ast, asc",b.Ast, b.Asc, b.Df,ColorReset)
	//this WILL BREAK for ribbed slabs
	/*
	switch b.Rslb{
		case true:
		err = b.RBarGen()
		case false:
		//err = b.BarGen()	
		//if err != nil{
			//fmt.Println("ERRORE,errore->",err)
		//	if b.D1 + b.D2 > 0.0{
		//		b.D1 = 0.0; b.D2 = 0.0
		//		err = b.BarGen()
		//		if err != nil{
					//fmt.Println("ERRORE,errore->",err)
		//			return
		//		}
		//	}
		//return
		//}
	}
	*/
	//err = b.BarLay()
	//if err != nil{
		//fmt.Println("ERRORE,errore->",err)
	//	return
	//}
	//if vu != 0 design for shear
	//fmt.Println("heare in beam design->",b.Rbrt,b.Rbrc,b.Id,b.Title)
	//fmt.Println("calling shear design")
	if b.Shrdz{
		
		switch{
			case len(b.Xs) == 0:
			err = BmShrChk(b)
			if err != nil{
				//fmt.Println("ERRORE,errore->",err)
				return
			}
			case len(b.Xs) == 21:
			err = BmShrDz(b, b.Lsx, b.Rsx, b.Xs, b.Vs)
				
			if err != nil{
				//fmt.Println("shear design ERRORE,errore->",err)
				return
			}
		}
	}
	//fmt.Println("end shear design")
	//CHECK FOR SPAN DEPTH RATIO CHECK FOR SPAN DEPTH RATIO CHECK FOR SPAN DEPTH RATIO
	//sdchk, sd, dserve  = b *RccBm, mspan, astprov, astreq, asc float64
	//no need to print the table if say, opt or whatever
	//b.Table(b.Verbose)
	if b.Ismid{
		//if b.Verbose{fmt.Println(ColorRed,"checking span depth ratio for ->",b.Title,ColorReset)}
		sdchk, _, dserve := BmSdratBs(b, b.Mu)
		//if b.Verbose{fmt.Println(ColorRed,"sdchk, sd, dserve ->",sdchk, sd, dserve)}
		if !sdchk{
			err = errors.New(fmt.Sprintf("span depth ratio check failed - required depth %.2f vs %.2f",dserve+b.Cvrt,b.Dused))
			return
		}
	}
	//if b.Verbose{b.Table(false)}
	b.Dz = true
	return
}

//BmAnalyze is the general entry func for beam section analysis (returns ultimate moment of resistance given steel bars/depths)
func BmAnalyze(b *RccBm) (err error){
	bscode := true
	if b.Code == 1{bscode = false}
	b.Dsgn = false
	b.Init()
	switch{
		case len(b.Dbars) == 0:
		switch b.Code{
			case 1:
			mur, _ := BmSecAzIs(b)
			//fmt.Println("mur->",mur, "astmax->",astmax)
			b.Mu = mur
			case 2:
			mr, x, _, _:= SecAzBs(b)
			//fmt.Println(mr, x, sdrat, err)		
			b.Mu = mr
			b.Xu = x
		}
		default:
		//var mr, x float64
		mr, x, _, _ := BmAzGen(b, bscode, true)
		//fmt.Printf("%s mr - > %f kn-m xd -> %f mm xdrat -> %f %s\n",ColorCyan,mr,x,xdrat,ColorReset)
		b.Mu = mr
		b.Xu = x
	}
	return
}

//BmAzGen - hulse 3.1.3 section analysis modified
//analysis of beams with multiple levels of rebar
//rblk - rect stress block
//analysis of non- rect sections 
func BmAzGen(b *RccBm, bscode, rblk bool) (mr, x, xdrat float64, err error) {
	var fsc, fst, C, T, deltaprev, cck, yc float64
	var div float64 = 10.0
	
	if b.Styp == -1{
		//init sec
		b.Sec.SecInit()
		if b.Sec.Prop.Area <= 0.0{
			err = fmt.Errorf("invalid section area - %v\ndata - %+v",b.Sec.Prop.Area,b.Sec)
			return
		}
	}
	effd := b.Dused - b.Cvrt
	k8 := 0.45
	k9 := 0.9
	switch b.Code{
		case 1:
		k9 = 0.8
	}
	var iter, kiter int
	for iter != -1{
		kiter++
		x += effd/div
		fsc = 0.0; fst = 0.0
		if kiter > 666{
			err = ErrIter
			return
		}
		cck, _, _ , _ = BmArXu(b, k8, k9, x, rblk)
		fst, fsc = BmRbrFrc(b, x, bscode)
		T = fst; C = fsc + cck
		deltaprev = T - C
		if T == C {
			iter = -1
			break
		}
		if (T - C) > 0.0 && deltaprev > 0.0{
			continue
		} else {
			if iter == 0{
				x -= effd/div
				div = div * 10.
				iter = 1
				x -= effd/div
			} else if T == C || T - C < 1e-6{
				iter = -1
				break
			}
		} 
	}
	xdrat = x/b.Dused
	//get moment about top face
	cck, _, _, yc = BmArXu(b, k8, k9, x, rblk)
	mr = cck * yc
	if b.Tyb > 0.0 && k9 * x > b.Df{
		//stress block in web
		mr = b.Fck * b.Bf * b.Df * b.Df/2.0 + b.Fck * b.Bw * (k9 * x - b.Df) * (k9 * x - b.Df)/2.0
	}
	for i, dia := range b.Dias{
		dbar := b.Dbars[i]
		if x > dbar {
			//SHOULD BE fsc := rebarforce(b.Fy, (x-dbar)/x); fsc -= b.Fck	
			switch bscode{
				case true:
				mr += RbrArea(dia) * (rebarforce(b.Fy, (x-dbar)/x) - b.Fck) * dbar
				case false:
				esc := 0.0035* (1.0 - dbar/x)
				mr += RbrArea(dia) * (RbrFrcIs(b.Fy, esc)-b.Fck) * dbar
			}
			//mr += RbrArea(dia) * rebarforce(b.Fy, (x-dbar)/x) * dbar
		} else {
			switch bscode{
				case true:
				mr -= RbrArea(dia) * rebarforce(b.Fy, (dbar-x)/x) * dbar
				case false:
				mr -= RbrArea(dia) * RbrFrcIs(b.Fy, 0.0035*(dbar-x)/x) * dbar
			}
			//mr -= RbrArea(dia) * rebarforce(b.Fy, (dbar-x)/x) * dbar
		}
	}
	mr = math.Abs(mr / 1e6); err = nil
	return
}

//SecAzBs returns ult. moment of resistance of a beam section
//as seen in mosley 3.1.3 section
//two single layers of rebar (ast, asc)
func SecAzBs(b *RccBm) (mr, x, xdrat float64, err error) {
	var fsc, fst, C, T, deltaprev float64
	var div float64 = 100.0
	effd := b.Dused - b.Cvrt
	k8 := 0.45
	k9 := 0.9
	var kiter int
	switch b.Tyb {
	case 1.0,0.5: //flanged
		//b.Bf := dims[0]
		//hf := dims[1]
		//bw := dims[2]
		if b.Asc == 0.0 {
			//singly reinforced flanged beam
			x = effd / div
			fst = rebarforce(b.Fy, (effd-x)/x)
			T = b.Ast * fst
			C = k8 * b.Fck * b.Bf * k9 * x
			if k9*x > b.Df {
				C = (k8 * b.Fck * b.Bf * b.Df) + k8*b.Fck*b.Bw*(k9*x-b.Df)
			}
			deltaprev = T - C
			for (T-C) > 0 && deltaprev > 0 {
				if kiter > 3000{err = ErrIter; return}
				x += effd / div
				fst = rebarforce(fy, (effd-x)/x)
				T = b.Ast * fst
				C = k8 * b.Fck * b.Bf * k9 * x
				if k9*x > b.Df {
					C = k8*b.Fck*b.Bf*b.Df + k8*b.Fck*b.Bw*(k9*x-b.Df)
				}
				deltaprev = T - C
				kiter++
			}
			deltaprev = 1.0
			T = 1.0
			C = 0.0
			x -= effd / div
			div = 1000.0
			kiter = 0
			for (T-C) > 0 && deltaprev > 0 {
				if kiter > 3000{err = ErrIter; return}
				x += effd / div
				fst = rebarforce(b.Fy, (effd-x)/x)
				T = b.Ast * fst
				C = k8 * b.Fck * b.Bf * k9 * x
				if k9*x > b.Df {
					C = k8*b.Fck*b.Bf*b.Df + k8*b.Fck*b.Bw*(k9*x-b.Df)
				}
				deltaprev = T - C
				kiter++
			}
			mr = k8 * b.Fck * b.Bf * k9 * x * (effd - k9*0.5*x)
			if k9*x > b.Df {
				mr = k8*b.Fck*b.Bf*b.Df*(effd-0.5*b.Df) + k8*b.Fck*b.Bw*(k9*x-b.Df)*(effd-k9*0.5*x-0.5*b.Df)
			}
		} else {
			//doubly reinforced flanged beam
			x = effd / div
			fst = rebarforce(b.Fy, (effd-x)/x)
			fsc = rebarforce(b.Fy, (x-b.Cvrc)/x)
			T = b.Ast * fst
			C = k8*b.Fck*b.Bf*k9*x + fsc*b.Asc
			if k9*x > b.Df {
				C = k8*b.Fck*b.Bf*b.Df + k8*b.Fck*b.Bw*(k9*x-b.Df) + fsc*b.Asc
			}
			deltaprev = T - C
			kiter = 0
			for (T-C) > 0 && deltaprev > 0 {
				if kiter > 3000{err = ErrIter; return}
				x += effd / div
				fst = rebarforce(b.Fy, (effd-x)/x)
				fsc = rebarforce(b.Fy, (x-b.Cvrc)/x)
				T = b.Ast * fst
				C = k8*b.Fck*b.Bf*k9*x + fsc*b.Asc
				if k9*x > b.Df {
					C = k8*b.Fck*b.Bf*b.Df + k8*b.Fck*b.Bw*(k9*x-b.Df) + fsc*b.Asc
				}
				deltaprev = T - C
				kiter++
			}
			deltaprev = 1.0
			T = 1.0
			C = 0.0
			kiter =0 
			x -= effd / div
			div = 1000.0
			for (T-C) > 0 && deltaprev > 0 {
				if kiter > 3000{err = ErrIter; return}
				x += effd / div
				fst = rebarforce(b.Fy, (effd-x)/x)
				fsc = rebarforce(b.Fy, (x-b.Cvrc)/x)
				T = b.Ast * fst
				C = k8*b.Fck*b.Bf*k9*x + fsc*b.Asc
				if k9*x > b.Df {
					C = k8*b.Fck*b.Bf*b.Df + k8*b.Fck*b.Bw*(k9*x-b.Df) + fsc*b.Asc
				}
				deltaprev = T - C
				kiter++
			}
			//mr = k8*b.Fck*b.Bf*b.Df*(effd-0.5*b.Df)
			mr = k8*b.Fck*b.Bf*k9*x*(effd-k9*0.5*x) + b.Asc*fsc*(effd-b.Cvrc)
			if k9*x > b.Df {
				mr = k8*b.Fck*b.Bf*b.Df*(effd-0.5*b.Df) + k8*b.Fck*b.Bw*(k9*x-b.Df)*(effd-k9*0.5*x-0.5*b.Df) + b.Asc*fsc*(effd-b.Cvrc)
			}
		}

		xdrat = x / effd
	case 0.0:
		//rectangular beam
		if b.Asc == 0.0 {
			//singly reinforced beam
			x = effd / div
			fst = rebarforce(b.Fy, (effd-x)/x)
			T = b.Ast * fst
			C = k8 * b.Fck * b.Bw * k9 * x
			deltaprev = T - C
			kiter = 0
			for (T-C) > 0 && deltaprev > 0 {
				if kiter > 3000{err = ErrIter; return}
				x += effd / div
				fst = rebarforce(b.Fy, (effd-x)/x)
				T = b.Ast * fst
				C = k8 * b.Fck * b.Bw * k9 * x
				deltaprev = T - C
				kiter++
			}
			deltaprev = 1.0
			T = 1.0
			C = 0.0
			//fmt.Println("sign change between xprev", x-(effd/div), x, " x next")
			x -= effd / div
			div = 1000.0
			kiter = 0
			for (T-C) > 0 && deltaprev > 0 {
				if kiter > 3000{err = ErrIter; return}
				x += effd / div
				fst = rebarforce(b.Fy, (effd-x)/x)
				T = b.Ast * fst
				C = k8 * b.Fck * b.Bw * k9 * x
				deltaprev = T - C
				kiter++
			}
			mr = k8 * b.Fck * b.Bw * k9 * x * (effd - k9*0.5*x)
		} else {
			//var cnt int
			//incl. compression steel
			x = b.Cvrc
			fst = rebarforce(b.Fy, (effd-x)/x)
			fsc = rebarforce(b.Fy, (x-b.Cvrc)/x)
			T = b.Ast * fst
			C = (k8 * b.Fck * b.Bw * k9 * x) + (b.Asc * fsc)
			deltaprev = T - C
			for (T-C) > 0 && deltaprev > 0 {
				if kiter > 3000{err = ErrIter; return}
				x += effd / div
				fst = rebarforce(b.Fy, (effd-x)/x)
				fsc = rebarforce(b.Fy, (x-b.Cvrc)/x)
				T = b.Ast * fst
				C = (k8 * b.Fck * b.Bw * k9 * x) + b.Asc*fsc
				deltaprev = T - C
				kiter++
			}
			deltaprev = 1.0
			T = 1.0
			C = 0.0
			x -= effd / div
			div = 1000.0
			kiter = 0
			for (T-C) > 0 && deltaprev > 0 {
				if kiter > 3000{err = ErrIter; return}
				x += effd / div
				fst = rebarforce(b.Fy, (effd-x)/x)
				fsc = rebarforce(b.Fy, (x-b.Cvrc)/x)
				T = b.Ast * fst
				C = (k8 * b.Fck * b.Bw * k9 * x) + b.Asc*fsc
				deltaprev = T - C
				kiter++
			}
			mr = k8*b.Fck*b.Bw*k9*x*(effd-k9*0.5*x) + b.Asc*fsc*(effd-b.Cvrc)
		}
		xdrat = x / effd
	}
	mr = mr / 1e6
	err = nil
	return
}

//BmStlBs returns asc and ast given an ult. moment mur
//as seen in hulse chapter 3
func BmStlBs(b *RccBm, mur float64) (err error){
	//get asc and ast for an ultimate moment mur
	var mb, xb, xna float64
	//eyield := b.Fy/1.15/200000.0

	var kiter, iter int
	
	for iter != -1{
		kiter++
		effd := b.Dused - b.Cvrt
		xb = (1.0 - b.DM - 0.4) * effd
		if b.DM <= 0.1 {
			xb = 0.5 * effd
		}
		k8 := 0.45
		k9 := 0.9
		switch b.Code{
			case 1:
			k9 = 0.8
		}
		switch{
			case b.Npsec:
			//TODO - weird beam design
			b.Init()
		//fmt.Println("NPSEC",b.Npsec)
		//fmt.Println(b.Dused, effd)
		//fmt.Println("bal section depth->",xb)
		//fmt.Println(b.Sec.Prop.Area)
		//cck, area, xc , yc := BmArXu(b, k8, k9, xb, true)
			
		//fmt.Println(cck, area, xc, yc, b.Sec.Ymx-b.Sec.Ym)
		//mb = cck * (effd - yc)
		//fmt.Println(mb/1e6)
		//if mur * 1e6 < mb{
		//	fmt.Println("only T")
		//}
		//cck, _, _ , _ = BmArXu(b, k8, k9, x, rblk)
			case b.Tyb == 0.0:
			//rectangular beam
			mb = k8 * b.Fck * b.Bw * k9 * xb * (effd - k9*xb/2.0)
			if mur*1e6 > mb {
				//design with compression steel
				b.Csteel = true
				xna = xb
				fsc := rebarforce(b.Fy, (xna-b.Cvrc)/xna)
				b.Asc = (mur*1e6 - mb) /fsc/(effd - b.Cvrc)
				b.Ast = (b.Asc * fsc + k8 * b.Fck * b.Bw * k9 * xb)/ (0.87 * b.Fy)
			} else {
				//tension steel ftw
				xna = effd/k9 - math.Sqrt(math.Pow((effd/k9), 2)-(2.0*mur*1e6)/(k8*b.Fck*b.Bw*k9*k9))
				b.Ast = k8 * b.Fck * k9 * xna * b.Bw / (0.87 * b.Fy)
				//select nbars, nrows
			}
			case b.Tyb == 1.0, b.Tyb == 0.5:
			// flanged beam
			switch {
			case b.Df > k9*xb:
				//huge flange, design as rectangular section
				mb = k8 * b.Fck * b.Bf * k9 * xb * (effd - k9*xb/2)
				if mur*1e6 > mb {
					//design with compression steel
					b.Csteel = true
					xna = xb
					fsc := rebarforce(b.Fy, (xna-b.Cvrc)/xna)
					b.Asc = (mur*1e6 - mb) / (fsc * (effd - b.Cvrc))
					b.Ast = ((b.Asc * fsc) + (k8 * b.Fck * b.Bf * k9 * xb)) / (0.87 * b.Fy)
				} else {
					//tension steel ftw
					xna = effd/k9 - math.Sqrt(math.Pow((effd/k9), 2)-(2.0*mur*1e6)/(k8*b.Fck*b.Bf*k9*k9))
					b.Ast = k8 * b.Fck * k9 * xna * b.Bf / (0.87 * b.Fy)
					//select nbars, nrows
				}
			default:
				mb = k8*b.Fck*b.Bf*b.Df*(effd-b.Df/2.0) + k8*b.Fck*b.Bw*(k9*xb-b.Df)*(effd-(k9*xb+b.Df)/2.0)
				mf := k8 * b.Fck * b.Bf * b.Df * (effd - b.Df/2.0)
				switch {
				case mur*1e6 > mb:
					//flanged beam with tension plus compression steel
					b.Csteel = true
					xna = xb
					fsc := rebarforce(b.Fy, (xna-b.Cvrc)/xna)
					b.Asc = (mur*1e6 - mb) / (fsc * (effd - b.Cvrc))
					b.Ast = ((b.Asc * fsc) + (k8*b.Fck*b.Bf*b.Df + k8*b.Fck*b.Bw*(k9*xb-b.Df))) / (0.87 * b.Fy)
				case mur*1e6 > mf:
					//stress block below the flange
					var x1, x2, x3 float64
					x1 = math.Pow(k9, 2.0)/2.0
					x2 = k9 * effd
					x3 = mur * 1e6/(k8*b.Fck)/b.Bw
					x3 -= k8*b.Fck*b.Bf*b.Df*(effd-b.Df/2.0)/(k8*b.Fck)/b.Bw
					x3 += b.Df * effd - b.Df * b.Df/2.0
					xna = (x2 - math.Sqrt(x2 * x2 - 4.0 * x1 * x3))/2.0/x1
					b.Ast = ((k8 * b.Fck * b.Bf * b.Df) + (k8 * b.Fck * b.Bw * (k9*xna - b.Df))) / (0.87 * b.Fy)
				case mur*1e6 <= mf:
					//concrete stress block within the flange
					xna = effd/k9 - math.Sqrt(math.Pow((effd/k9), 2)-(2.0*mur*1e6)/(k8*b.Fck*b.Bf*k9*k9))
					b.Ast = k8 * b.Fck * b.Bf * k9 * xna / (0.87 * b.Fy)
				}
			}
		}
		b.Xu = xna
		//iter := 0		
		if b.Flip && !b.Ismid{
			//fmt.Println(ColorRed,"behold yeay flippening","ast",b.Ast, "asc",b.Asc,ColorReset)
			a1 := b.Asc; a2 := b.Ast
			b.Ast = a1; b.Asc = a2
			//fmt.Println(ColorYellow,"yeay flippening","ast",b.Ast, "asc",b.Asc,ColorReset)
		}
		
		if b.Rslb {
			err = b.RBarGen()
			return
		}
		//fmt.Println("kiter, bm, mur, ast, asc, tyb, bw, dused, ->",kiter,mur,b.Ast,b.Asc,b.Tyb,b.Bw,b.Dused)
		//fmt.Println("vals->",b.Ast,b.Asc)		
		e, mincvr := b.BarGen()
		//fmt.Println("error->",e,"mincvr->",mincvr, "act cover->",b.Rbrc)
		//fmt.Println("ast->",b.Ast)
		if e != nil{
			//fmt.Println("ERRORE,errore->",e)
			//equate both so that plot looks cleaner
			b.Cvrc = mincvr
			b.Cvrt = mincvr
			b.Ast = 0.0
			b.Asc = 0.0
			//if e.Error() == "effective cover error1"{
			//	b.Cvrc = mincvr
			//	b.Cvrt = mincvr
			//} else {
			//}
		} else {
			//fmt.Println("ast, asc, rbrc and rbrt->", b.Ast, b.Asc, b.Rbrc, b.Rbrt)
			err = b.BarLay()
			iter = -1
		}
		if kiter > 6{
			iter = -1
			err = errors.New(fmt.Sprint(e," and iteration error"))
		}
	}
	//fmt.Println("bm, mur, ast, asc, tyb, bw, dused, ->",b.Title,mur,b.Ast,b.Asc,b.Tyb,b.Bw,b.Dused)
	//b.BarGen()
	//b.BarLay()
	//fmt.Println("df stl->",b.Mid,b.Id,b.Df)
	return
}

//Dump saves an RccBm to a json file
func (b *RccBm) Dump()(filename string, err error){
	_, bp, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(bp)
	fname := fmt.Sprintf("beam_%v.json",b.Id)
	foldr := filepath.Join(basepath,"../data/out",time.Now().Format("06-Jan-02"))
	if _, e := os.Stat(foldr); errors.Is(e, os.ErrNotExist) {
		e := os.Mkdir(foldr, os.ModePerm)
		if e != nil {
			err = e; return
		}
	}
	filename = filepath.Join(foldr,fname)
	data, e := json.Marshal(b)
	if e != nil{err = e; return}
	err = ioutil.WriteFile(filename, data, 0644)
	return
}

