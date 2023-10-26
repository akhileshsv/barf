
package barf

import (
	"fmt"
	"math"
	"errors"
	kass"barf/kass"
	"github.com/AlecAivazis/survey/v2"
	//"math"
)

//CalcCBm is the entry func from menu/flags for rcc continuous beam analysis and design
func CalcCBm(cb CBm) (err error){
	//entry func for menu and flags
	switch {
	case cb.Opt > 0:
		CBmOpt(cb)
	default:	
		var bmenv map[int]*kass.BmEnv
		bmenv, err = CBeamEnvRcc(&cb, cb.Term, true)
		if err != nil{
			return
		}
		err = CBmDz(&cb,bmenv)
		if err != nil{
			return
		}
	}
	return
}


//CBmDz designs a continuous beam at lsup, mspan, rsup
//dm - moment redistribution 0 - none, 1 - use
//dz - 0 - as is, 1 - envelope, 2 - +self wt 3 - envelope
func CBmDz(cb *CBm, bmenv map[int]*kass.BmEnv) (err error){
	bmchn := make(chan []interface{}, len(bmenv))

	for _,i := range cb.Bmvec{
		//fmt.Println(i)
		//fmt.Println(bmenv[i].Vl)
		//fmt.Println(bmenv[i].Vr)
		//fmt.Println(ColorCyan)
		//fmt.Printf("+%v\n",cb.RcBm[i-1][0])
		//fmt.Println(ColorReset)
		//fmt.Println("lsx,rsx->",cb.RcBm[i-1][1].Lsx,cb.RcBm[i-1][1].Rsx)
		go Bmdz(cb.Code, cb.RcBm[i-1], bmenv[i], cb.Term, cb.Verbose, bmchn)
	}
	var allcons int
	var evec [][]error
	evec = make([][]error, len(cb.Bmvec))
	for _ = range cb.Bmvec{
		rez := <- bmchn
		idx, _ := rez[0].(int)
		cons, _ := rez[1].(int)
		errz, _ := rez[2].([]error)
		allcons += cons
		if cb.Verbose{
			//fmt.Println("beam id->",rez[0])
			//fmt.Println("cons->",rez[1])
			//fmt.Println(ColorRed,"ERRORE,errore errors->",rez[2],ColorReset)
		}
		evec[idx-1] = append(evec[idx-1],errz...)
	}
	if allcons > 0{
		if cb.Verbose{
			for i, e := range evec{
				fmt.Println("span=>",i+1,"ERRORE,errorez->\n",e)
			}
		}
		err = errors.New("strange design error")
	} 
	//for _, txtplot := range cb.Txtplots{
	//	fmt.Println(txtplot)
	//}
	return
}

//GetBmArr generates an rcc beam slice for a beam span
//bmarr - > left sec, mid span, right sec
func GetBmArr(bmarr []*RccBm, bm *kass.BmEnv,kostin []float64,fck, fy, fyv, efcvr, dm, d1, d2, dslb float64,code, bstyp int, verbose bool){
	//for generating an rcc bm arr (updates vals and etc)
	var bf, df, bw, dused, tyb float64
	var npsec bool
	tyb = 0.0
	switch bstyp{
		case 1:
		bw = bm.Dims[0]; dused = bm.Dims[1]
		case 6,7,8,9,10:
		bf = bm.Dims[0]; dused = bm.Dims[1]; bw = bm.Dims[2]; df = bm.Dims[3]
		case 14:
		bf = bm.Dims[0]; dused = bm.Dims[1]; bw = bm.Dims[2]; df = bm.Dims[3]
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
	var ismid bool
	for j := 0; j < 3; j++{
		ismid = false
		flip := true
		switch bstyp{
			case 1,6,7,8,9,10:
			sti = 1; tybi = 0.0
			case 14:
			sti = bstyp; tybi = tyb
		}
		//flange only at midspan
		if j == 1 && bm.Endc > 0{
			sti = bstyp; tybi = tyb
			flip = false
			ismid = true
		}
		bmarr[j] = &RccBm{
			Id:j+1,
			Mid:bm.Id,
			Fck:fck,
			Fy:fy,
			Fyv:fyv,
			Bf:bf,
			Df:df,
			Bw:bw,
			Dused:dused,
			Styp:sti,
			Tyb:tybi,
			Cvrt:efcvr,
			Cvrc:efcvr,
			Flip:flip,
			Code:code,
			Endc:bm.Endc,
			Dims:bm.Dims,
			Npsec:npsec,
			Verbose:verbose,
			DM:dm,
			Lsx:bm.Lsx*1e3,
			Rsx:bm.Rsx*1e3,
			D1:d1,D2:d2,
			Ldx:bm.Ldx, Rdx:bm.Rdx,
			Dslb:dslb,
			Ismid:ismid,
			Term:bm.Term,
			Rslb:bm.Rslb,
		}
		if bm.Rslb{
			bmarr[j].Monolith = true
		}
		bmarr[j].Init()
	}
	return
}

//Bmdz is the routing func for (2d) beam design for subframe and frame 2d funcs
func Bmdz(code int, barr []*RccBm, bm *kass.BmEnv, term string, verbose bool, bmchn chan []interface{}){
	//FUTURE (imma tell you how it was) - BmDz - 3d beam design entry func 
	switch bm.Endc{
		case 0:
		//clvr - design for max (Mpmax, Mnmax) THIS IS so wrong
		ClvrSpanDz(code, barr, bm, term, verbose, bmchn)
		case 1:
		SsSpanDz(code, barr, bm, term, verbose, bmchn)
		default:
		//generic beam
		CsSpanDz(code, barr, bm, term, verbose, bmchn)
	}
}

//CsSpanDz designs a continuous span beam array
func CsSpanDz(code int, barr []*RccBm, bm *kass.BmEnv, term string, verbose bool, bmchn chan []interface{}){
	//continuous span design
	var cons int
	if verbose{fmt.Println(ColorCyan,"beam id->",bm.Id,ColorReset)}
	errz := []error{}
	rez := make([]interface{},3)
	rez[0] = bm.Id
	
	//CHECK - if l/r end span, adjust end support moment to 50% of midspan moment
	
	//get midspan moment, check if -ve moment is greater at that point
	//var mpspn, mnspn float64
	//var flipmid bool
	//if barr[0].DM == 0.0{
	//	mpspn = bm.Mpmax
	//	mnspn = math.Abs(bm.Mnenv[10])
		
	//} else {
	//	mpspn = bm.Mprmax
	//	mnspn = math.Abs(bm.Mnrd[10])
	//}
	//if mpspn > mnspn{
	//	mnspn = mpspn
	//	flipmid = true
	//}
	for i := range barr{
		barr[i].Mid = bm.Id
		switch i{
			case 0:
			barr[i].Id = 0 
			barr[i].Title = fmt.Sprintf("bm_%v_left",bm.Id)
			if barr[i].DM == 0.0{
				barr[i].Mu = math.Abs(bm.Ml)
				barr[i].Vu = math.Abs(bm.Vl)
				if barr[i].Ldx == 1 && bm.Frmtyp == 1{
					barr[i].Mu = math.Abs(bm.Mpmax)/2.0
				}
			} else {
				barr[i].Mu = math.Abs(bm.Mlrd)
				barr[i].Vu = math.Abs(bm.Vlrd)
				if barr[i].Ldx == 1 && bm.Frmtyp == 1{
					barr[i].Mu = math.Abs(bm.Mprmax)/2.0				
				}
			}
			barr[1].Lspan = bm.Xs[20]
			
			case 1:
			barr[i].Flip = false
			barr[i].Id = 1 
			barr[i].Title = fmt.Sprintf("bm_%v_mid",bm.Id)
			if barr[i].DM == 0.0{
				barr[i].Mu = math.Abs(bm.Mpmax)
				//if barr[i].Mu == 0.0{barr[i].Mu = math.Abs(bm.Mnmax)}
				barr[i].Xs = bm.Xs; barr[i].Vs = bm.Venv
				barr[i].Mns = bm.Mnenv
				barr[i].Mps = bm.Mpenv
				barr[1].Lspan = bm.Xs[20]
			} else {
				barr[i].Mu = math.Abs(bm.Mprmax)
				//if barr[i].Mu == 0.0{barr[i].Mu = math.Abs(bm.Mnrmax)}
				barr[i].Xs = bm.Xs; barr[i].Vs = bm.Venv			
				barr[i].Mns = bm.Mnrd
				barr[i].Mps = bm.Mprd
				barr[1].Lspan = bm.Xs[20]
			}
			barr[i].Shrdz = true
			barr[i].S1, barr[i].S3 = 0.3, 0.3
			switch {
			case barr[i].Ldx == 1:
				barr[i].S1 = 0.15
			case barr[i].Rdx == 1:
				barr[i].S3 = 0.15
			}
			case 2:
			barr[i].Id = 2 
			barr[i].Title = fmt.Sprintf("bm_%v_right",bm.Id)
			if barr[i].DM == 0.0{
				barr[i].Mu = math.Abs(bm.Mr)
				barr[i].Vu = math.Abs(bm.Vr)
				
				if barr[i].Rdx == 1 && bm.Frmtyp == 1{
					barr[i].Mu = math.Abs(bm.Mpmax)/2.0
					//barr[i].Vu = math.Abs(bm.Vl)					
				}
			} else {
				barr[i].Mu = math.Abs(bm.Mrrd)
				barr[i].Vu = math.Abs(bm.Vrrd)
				
				if barr[i].Rdx == 1 && bm.Frmtyp == 1{
					barr[i].Mu = math.Abs(bm.Mpmax)/2.0
					//barr[i].Vu = math.Abs(bm.Vl)					
				}
			}
			barr[1].Lspan = bm.Xs[20]
		}
		//if verbose{fmt.Printf("beam data->%+v\n",barr[i])}
		err := BmDesign(barr[i])
		
		errz = append(errz, err)
		if err != nil{
			//abbe.RETURN HERE
			
			//fmt.Println(ColorRed,"ERRORE,errore->",err,ColorReset)
			cons++
			break
		}
	}
	rez[1] = cons
	rez[2] = errz
	if cons == 0{
		Quant(barr)
		//barr[1].Table(false)
		//if verbose{
		//	//barr[1].Table(true)
		//}
	}
	bmchn <- rez
}

//SsSpanDz designs a simply supported beam array
func SsSpanDz(code int, barr []*RccBm, bm *kass.BmEnv, term string, verbose bool, bmchn chan []interface{}){
	//simply supported span design
	//moments at face of support are not applicable? 
	var cons int
	//if verbose{fmt.Println(ColorCyan,"beam id->",bm.Id,ColorReset)}
	errz := []error{}
	for i := range barr{
		barr[i].Mid = bm.Id
		barr[i].Shrdz = true
		switch i{
			case 0:
			//YEOLDE
			//
			barr[i].Shrdz = false
			barr[i].Flip = true
			barr[i].Id = 0 
			barr[i].Title = fmt.Sprintf("bm_%v_left",bm.Id)
			//page 199, item 3 mosley-hulse ec2 design
			//25 percent of span moment at supports to be considered when a simple support has been assumed
			
			if barr[i].Monolith{				
				//in monolithic construction
				if barr[i].DM == 0.0{
					barr[i].Mu = math.Abs(bm.Mpmax)/4.0
				} else {
					barr[i].Mu = math.Abs(bm.Mprmax)/4.0
				}
				
				err := BmDesign(barr[i])
				errz = append(errz, err)
				if err != nil{
					cons++
				}	
			} else{
				barr[i].Ignore = true			
			}
			case 1:
			barr[i].Id = 1 
			barr[i].Title = fmt.Sprintf("bm_%v_mid",bm.Id)
			barr[i].Shrdz = true
			if barr[i].DM == 0.0{
				barr[i].Mu = math.Abs(bm.Mpmax)
				barr[i].Xs = bm.Xs; barr[i].Vs = bm.Venv
				barr[i].Mns = bm.Mnenv
				barr[i].Mps = bm.Mpenv
				barr[1].Lspan = bm.Xs[20]
			} else {
				barr[i].Mu = math.Abs(bm.Mprmax)
				barr[i].Xs = bm.Xs; barr[i].Vs = bm.Venv			
				barr[i].Mns = bm.Mnrd
				barr[i].Mps = bm.Mprd
				barr[1].Lspan = bm.Xs[20]
			}	
			err := BmDesign(barr[i])
			errz = append(errz, err)
			if err != nil{
				cons++
			}
			case 2:
			barr[i].Shrdz = false
			barr[i].Flip = true
			barr[i].Id = 2 
			barr[i].Title = fmt.Sprintf("bm_%v_right",bm.Id)
			if barr[i].Monolith{	
				if barr[i].DM == 0.0{
					barr[i].Mu = math.Abs(bm.Mpmax)/4.0
				} else {
					barr[i].Mu = math.Abs(bm.Mprmax)/4.0
				}
				
				err := BmDesign(barr[i])
				errz = append(errz, err)
				if err != nil{
					cons++
				}
			} else{
				barr[i].Ignore = true
			}
		}
	}
	rez := make([]interface{},3)
	rez[0] = bm.Id
	rez[1] = cons
	rez[2] = errz
	if cons == 0{
		Quant(barr)
		if verbose{
			barr[1].Table(true)
		}
	}
	bmchn <- rez
}

//ClvrSpanDz designs a cantilever beam array
func ClvrSpanDz(code int, barr []*RccBm, bm *kass.BmEnv, term string, verbose bool, bmchn chan []interface{}){
	//if verbose{fmt.Println(ColorCyan,"beam id->",bm.Id,ColorReset)}
	var cons int
	if verbose{fmt.Println(ColorCyan,"beam id->",bm.Id,ColorReset)}
	errz := []error{}
	for i := range barr{
		barr[i].Mid = bm.Id
		barr[i].Shrdz = true
		switch i{
			case 0:
			barr[i].Ignore = true
			case 1:
			barr[i].Id = 1 
			barr[i].Title = fmt.Sprintf("bm_%v_mid",bm.Id)
			if barr[i].DM == 0.0{
				barr[i].Mu = math.Abs(bm.Mpmax)
				barr[i].Xs = bm.Xs; barr[i].Vs = bm.Venv
				barr[i].Mns = bm.Mnenv
				barr[i].Mps = bm.Mpenv
				barr[1].Lspan = bm.Xs[20]
			} else {
				barr[i].Mu = math.Abs(bm.Mprmax)
				barr[i].Xs = bm.Xs; barr[i].Vs = bm.Venv			
				barr[i].Mns = bm.Mnrd
				barr[i].Mps = bm.Mprd
				barr[1].Lspan = bm.Xs[20]
			}
			err := BmDesign(barr[i])
			errz = append(errz, err)
			if err != nil{cons++}
			case 2:
			barr[i].Ignore = true
		}
	}
	rez := make([]interface{},3)
	rez[0] = bm.Id
	rez[1] = cons
	rez[2] = errz
	if cons == 0{
		Quant(barr)
		if verbose{
			barr[1].Table(true)
		}
	}
	bmchn <- rez
}

//ViewRbr prints all rebar options as a slice of string
func (b *RccBm) ViewRbr(tcdx int) (opt []string){
	var rez [][]float64
	switch tcdx{
		case 1:
		rez = b.Rbrtopt
		case 2:
		rez = b.Rbrcopt
	}
	for i, r := range rez{
		
		row := fmt.Sprintf("%v. dia 1 %.0f no %.0f dia 2 %.0f no %.0f ast prov %.0f ast req %.2f diff %.2f",i+1,r[2],r[0],r[3],r[1],r[4],r[5],r[6])
		opt = append(opt, row)
	}
	return
}

//Tweaks - HOW MANY OF THESE ARE LYING AROUND
func (b *RccBm)Tweaks(){
	//tweak - beam steel, beam fck/fy
	//	
	var choice, iter int
	for iter != -1{
		prompt := &survey.Select{
			Message: fmt.Sprintf("choose param to tweak for beam_%v_%v",b.Mid,b.Id),
			Options: []string{"top steel", "bottom steel","enter","continue"},
		}
		survey.AskOne(prompt, &choice)
		switch choice{
			case 0,1:
			tcdx := 1
			if choice == 1{
				tcdx = 2
			}
			opt := b.ViewRbr(tcdx)
			pr1 := &survey.Select{
				Message: fmt.Sprintf("choose rbr opt for beam_%v_%v",b.Mid,b.Id),
				Options: opt,
			}
			var rc int
			survey.AskOne(pr1, &rc)
			fmt.Println(rc)
			//iter = -1
			case 2:
			case 3:
			iter = -1
			break
		}
	}
	return
}
