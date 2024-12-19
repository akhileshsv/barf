package barf

//load combination calculation functions

import (
	"fmt"
	"math"
	"errors"
)

var PsfsIs = map[int][][]float64{
	-1: [][]float64{{1.0,1.0,0.0}},
	0:  [][]float64{
		{1.5,1.5,0.0},
		{1.5,0.0,1.5},
	},
	1:  [][]float64{
		{1.5,1.5,0.0},
		{1.2,1.2,1.2},
		{1.5,0.0,1.5},
	},
}

//GetPSFs returns partial safety factors for load cases as per is/bs code
func GetPsfs(code, nlds, nwl, nsl int)(psfs [][]float64){
	switch code{
		case 1:
		
		switch nlds{
			case 1:
			//dead load
			psfs = [][]float64{{1.5,0.0,0.0}}
			case 2:
			//dead, live load
			psfs = [][]float64{
				{1.5, 1.5, 0.0},
			}
			case 3:
			//dead, live, wind
			psfs = [][]float64{
				{1.5, 1.5, 0.0},
			}
			for i := 0; i < nwl; i++{
				psfs = append(psfs, []float64{1.5, 0.0, 1.5})
			}
			
			for i := 0; i < nwl; i++{
				psfs = append(psfs, []float64{1.2, 1.2, 1.2})
			}
			case 4:
		}
		case 13:
		//timber/old is800(1973 methinks)
		
		switch nlds{
			case 1:
			//dead load
			psfs = [][]float64{{1.0,0.0,0.0}}
			case 2:
			//dead, live load
			psfs = [][]float64{
				{1.0, 1.0, 0.0},
			}
			case 3:
			//dead, live, wind
			psfs = [][]float64{
				{1.0, 1.0, 0.0},
			}
			
			// for i := 0; i < nwl; i++{
			// 	psfs = append(psfs, []float64{1.0, 0.0, 1.0})
			// }
			for i := 0; i < nwl; i++{
				psfs = append(psfs, []float64{1.0, 1.0, 1.0})
			}
		}
	}
	return
}

//CalcSrvLds superposes service load cases with PSFs
func CalcSrvLds(mod *Model) (err error){
	//change this when seismic happenz
	mod.Psfs = GetPsfs(1, mod.Nlds, mod.Nwlc, 0)
	ldmap := map[int]string{1:"dl",2:"ll",3:"wl",4:"sl"}
	if mod.Spam{fmt.Println(ColorRed,"number of load cases-",mod.Nlp,ColorReset)}
	mod.Nlp = len(mod.Jldcs) + len(mod.Mldcs) + len(mod.Psfs)
	mod.Ms = make(map[int]*Mem)
	mod.Js = make(map[int]*Node)
	mod.Reports = make([]string,mod.Nlp)
	mod.Txtplots = make([]string,mod.Nlp)
	mod.Lpmap = make(map[int]float64)
	mod.Grez = make(map[int][]float64)
	var frmrez []interface{}
	var pltstr string
	lpmap := make(map[float64]int)
	lp := 1
	switch mod.Frmstr{
		case "2dt":
		//first calc each srv load
		//then run thru psfs
		//for all of em, if wlon - for _, wlcs - range wlcs
		if len(mod.Jldcs) == 0{
			err = errors.New("zero force model???")
			return
		}
		for i := 1; i <= mod.Ngrps; i++{
			mod.Grez[i] = make([]float64, 2)
		}
		for i, jlds := range mod.Jldcs{
			//fmt.Println("calcing load case-",i,lp)
			modlp := Model{Id:mod.Id,Cmdz:mod.Cmdz,
				Term:mod.Term,
				Coords:mod.Coords,Mprp:mod.Mprp,
				Supports:mod.Supports,
				Em:mod.Em,Cp:mod.Cp,Dims:mod.Dims,
				Frmtyp:mod.Frmtyp, Frmstr:mod.Frmstr,
				Units:mod.Units}
			modlp.Id += fmt.Sprintf("lp-%v-%.1f-%s",lp,i,ldmap[int(i)])
			modlp.Jloads = jlds
			frmrez, pltstr, err = CalcModSer(&modlp)			
			if err == nil{
				mod.MemRez(lp,frmrez)
			}
			mod.Lpmap[lp] = i
			lpmap[i] = lp
			mod.Txtplots[lp-1] = pltstr
			lp++
			//HERE
			// if i > 0{
			// 	return
			// }
		}
	}
	for _, sfs := range mod.Psfs{
		//fmt.Println("sposing safety factor #",i+1, sfs)
		lds := []float64{1.0, 2.0}
		if sfs[2] > 0.0{
			for i := 0; i < mod.Nwlc; i++{
				val := 3.0 + float64(i+1)/10.0
				lds = append(lds, val)
			}
		}
		//fmt.Println("loads",lds)
		for m, mem := range mod.Ms{
			//cp := mem.Mprp[3]
			mod.Ms[m].Qfr[lp-1] = make([]float64, mod.Ncjt*2)
			for _, lpval := range lds{
				idx := int(lpval)
				lpdx := lpmap[lpval]
				for j, fval := range mem.Qfr[lpdx-1]{
					mod.Ms[m].Qfr[lp-1][j] += sfs[idx-1] *  fval
				}
			}
			switch mod.Frmstr{
				case "2dt":
				fval := mod.Ms[m].Qfr[lp-1][0]
				switch{
					case fval > 0.0:
					if fval > mod.Ms[m].Cmax{
						mod.Ms[m].Cmax = fval
					}
					case fval < 0.0:
					if math.Abs(fval) > mod.Ms[m].Tmax{
						mod.Ms[m].Tmax = math.Abs(fval)
					}
				}
			}
		}
		lp++
	}
	switch mod.Frmstr{
		case "2dt":	
		for _, mem := range mod.Ms{
			cp := mem.Mprp[3]
			if mod.Grez[cp][0] < mem.Cmax{
				mod.Grez[cp][0] = mem.Cmax
			}
			if mod.Grez[cp][1] < mem.Tmax{
				mod.Grez[cp][1] = mem.Tmax
			}
		}
		mod.Reports = append(mod.Reports,Trs2dSrvTable(mod))
	}
	return
}

//CalcLdCombos calculates individual load combos with PSFs
func CalcLdCombos(mod *Model) (err error){
	ldmap := map[int]string{1:"dl",2:"ll",3:"wl",4:"sl"}
	nlp := len(mod.Mldcs)
	if nlp == 0{nlp = len(mod.Jldcs)}
	if nlp == 0{
		err = errors.New("zero force model???")
		return
	}
	mod.Nlp = nlp
	lp := 1
	mod.Ms = make(map[int]*Mem)
	mod.Js = make(map[int]*Node)
	mod.Reports = make([]string,nlp)
	mod.Txtplots = make([]string,nlp)
	mod.Lpmap = make(map[int]float64)
	var frmrez []interface{}
	var pltstr string
	if len(mod.Mldcs) == 0{
		for i, jlds := range mod.Jldcs{ 
			modlp := &Model{Id:mod.Id,Cmdz:mod.Cmdz,Term:mod.Term,Coords:mod.Coords,Mprp:mod.Mprp,Supports:mod.Supports,Em:mod.Em,Cp:mod.Cp,Dims:mod.Dims,Frmtyp:mod.Frmtyp,Ncjt:2,Noprnt:mod.Noprnt,Frmstr:"2dt",Units:mod.Units}
			//fmt.Println("lp->",lp, jlds)
			modlp.Id += fmt.Sprintf("lp-%v-%.1f-%s",lp,i,ldmap[int(i)])

			modlp.Jloads = jlds
			frmrez, pltstr, err = CalcModSer(modlp)
			mod.Lpmap[lp] = i
			if err == nil{
				mod.MemRez(lp,frmrez)
			}
			lp++
		}
	} else {	
		for i, mlds := range mod.Mldcs{
			jlds := mod.Jldcs[i]
			modlp := Model{Id:mod.Id,Cmdz:mod.Cmdz,Term:mod.Term,Coords:mod.Coords,Mprp:mod.Mprp,Supports:mod.Supports,Em:mod.Em,Cp:mod.Cp,Dims:mod.Dims,Frmtyp:mod.Frmtyp, Frmstr:mod.Frmstr, Units:mod.Units}
			modlp.Id += fmt.Sprintf("lp-%v-%.1f-%s",lp,i,ldmap[int(i)])
			modlp.Msloads = mlds
			modlp.Jloads = jlds
			frmrez, pltstr, err = CalcModSer(&modlp)
			mod.Lpmap[lp] = i			
			if err == nil{
				mod.MemRez(lp,frmrez)
			}
			mod.Txtplots[lp-1] = pltstr
			lp++
		}
	}
	return
}

//CalcLdCombosPar generates and calcs load combinations based on mod.PSFs / mod.Code
//in parallel
//it calls GenLdCombos for each lp index and chills
func CalcLdCombosPar(mod *Model) (err error){
	ldmap := map[int]string{1:"dl",2:"ll",3:"wl",4:"sl"}
	nlp := len(mod.Mldcs)
	if nlp == 0{nlp = len(mod.Jldcs)}
	if nlp == 0{
		err = errors.New("zero force model???")
		return
	}
	mod.Nlp = nlp
	rezchn := make(chan []interface{},nlp)
	lp := 1
	mod.Ms = make(map[int]*Mem)
	mod.Js = make(map[int]*Node)
	mod.Reports = make([]string,nlp)
	mod.Txtplots = make([]string,nlp)
	mod.Lpmap = make(map[int]float64)
	if len(mod.Mldcs) == 0{
		for i, jlds := range mod.Jldcs{ 
			modlp := Model{Cmdz:mod.Cmdz,Term:mod.Term,Coords:mod.Coords,Mprp:mod.Mprp,Supports:mod.Supports,Em:mod.Em,Cp:mod.Cp,Dims:mod.Dims,Frmtyp:mod.Frmtyp}
			//fmt.Println("lp->",lp, jlds)
			
			modlp.Id += fmt.Sprintf("lp-%v-%.1f-%s",lp,i,ldmap[int(i)])
			modlp.Jloads = jlds
			//fmt.Println("load pattern->",lp,i)
			//fmt.Println(FrcTable(&modlp))
			go CalcModPar(lp,modlp,rezchn)
			mod.Lpmap[lp] = i
			lp++
		}
	} else {	
		for i, mlds := range mod.Mldcs{
			jlds := mod.Jldcs[i]
			modlp := Model{Cmdz:mod.Cmdz,Term:mod.Term,Coords:mod.Coords,Mprp:mod.Mprp,Supports:mod.Supports,Em:mod.Em,Cp:mod.Cp,Dims:mod.Dims,Frmtyp:mod.Frmtyp}
			
			modlp.Id += fmt.Sprintf("lp-%v-%.1f-%s",lp,i,ldmap[int(i)])
			modlp.Msloads = mlds
			modlp.Jloads = jlds
			go CalcModPar(lp,modlp,rezchn)
			mod.Lpmap[lp] = i
			lp++		
		}
	}
	for i := 0; i < nlp; i++{
		rez := <- rezchn
		lp, _ := rez[0].(int)
		err,_ := rez[1].(error)
		frmrez,_ := rez[2].([]interface{})
		if mod.Term == "dumb" || mod.Term == "mono"{
			txtplt, _ := rez[3].(string)
			mod.Txtplots[lp-1] = txtplt
		}
		//fmt.Println("frmrez->",len(frmrez))
		if err == nil{
			mod.MemRez(lp,frmrez)
		}
		//fmt.Println(ColorCyan,"load pattern->",lp,ColorReset)
		//fmt.Println(ColorYellow,"ERRORE,errore->",ColorRed,err,ColorReset)
	}
	/*
	for idx := 1; idx <= len(mod.Mprp); idx++{
		log.Println(ColorYellow, "member->",idx)
		mem := mod.Ms[idx]
		log.Println("Max loads->\n",ColorRed,"tensile->",mem.Tmax, "N",ColorCyan,"kompressive->",mem.Cmax,"N",ColorReset)
		log.Println("Pf->",mem.Pu)
		log.Println(mem)
	}
	*/
	return
}

//GenLdCombos generates load combos for 3 load cases now - 1 - dl, 2 - ll, 3 - wl
//nwl, nsl - number of wind load and seismic load cases
func GenLdCombos(nwl, nsl int, jldsrv, mldsrv map[float64][][]float64, psfs [][]float64)(mldcs, jldcs map[float64][][]float64){
	mldcs = make(map[float64][][]float64)
	jldcs = make(map[float64][][]float64)
	for i, sfs := range psfs{
		var wlon, slon bool
		//log.Println(sfs)
		lc := float64(i+1)
		if len(sfs) > 2{
			if sfs[2] > 0.0{
				wlon = true
			}
			if len(sfs) > 3{
				if sfs[3] > 0.0{
					slon = true
				}
			}
		}
		switch{
			case wlon:
			//wind load - gen cases for all wind load cs
			//log.Println("ncases->",nwl)
			//log.Println("lc number->",lc)
			for j := 0; j < nwl; j++{
				//if j > 2{continue}
				wlc := 3.0 + float64(j+1)/10.0
				lc = float64(i+1) + float64(j+1)/10.0
				sf := sfs[2]
				if len(jldsrv[wlc]) > 0{
					//log.Println("wind load case",wlc,"\n",jldsrv[wlc])
					
					for _, ld := range jldsrv[wlc]{
						ldvec := make([]float64, len(ld))
						ldvec[0] = ld[0]
						for k, val := range ld[1:]{
							val = sf * val
							ldvec[k+1] = val
						}
						//ldvec = append(ldvec, float64(3))
						jldcs[lc] = append(jldcs[lc],ldvec)
					}
				}
				if len(mldsrv[wlc]) > 0{
					for _, ld := range mldsrv[wlc]{	
						
						ldvec := make([]float64, len(ld))
						ldvec[0] = ld[0]
						for k, val := range ld[1:]{
							if k == 1 || k == 2{
								val = sf * val
							}
							ldvec[k+1] = val
						}
						//ldvec = append(ldvec, float64(3))
						mldcs[lc] = append(jldcs[lc],ldvec)
					}	
				}
				//now do other loads
				for ltyp, vec := range jldsrv{	
					lt := int(ltyp)
					sf := sfs[lt-1]
					if sf == 0.0{continue}
					if lt == 3{continue}
					for _, ld := range vec{
						ldvec := make([]float64, len(ld))
						for j, val := range ld{
							if j > 0{
								val = sf * val
							}
							ldvec[j] = val
						}
						//ldvec = append(ldvec, float64(lt))
						jldcs[lc] = append(jldcs[lc],ldvec)
					}
				}	
				for ltyp, vec := range mldsrv{
					lt := int(ltyp)
					sf := sfs[lt-1]
					if sf == 0.0{continue}
					if lt == 3{continue}
					for _, ld := range vec{
						ldvec := make([]float64, len(ld))
						for j, val := range ld{
							if j == 2 || j == 3{
								val = sf * val
							}
							ldvec[j] = val
						}
						//ldvec = append(ldvec, float64(lt))
						mldcs[lc] = append(mldcs[lc],ldvec)
					}
				}
			}
			case slon:
			//seismic load cases
			default:
			//either dead or live
			for ltyp, vec := range jldsrv{
				lc = float64(i+1.0)
				lt := int(ltyp)
				sf := sfs[lt-1]
				if sf == 0.0{continue}
				for _, ld := range vec{
					ldvec := make([]float64, len(ld))
					for j, val := range ld{
						if j > 0{
							val = sf * val
						}
						ldvec[j] = val
					}
					//ldvec = append(ldvec, float64(lt))
					jldcs[lc] = append(jldcs[lc],ldvec)
				}
			}	
			for ltyp, vec := range mldsrv{
				lt := int(ltyp)
				sf := sfs[lt-1]
				if sf == 0.0{continue}
				for _, ld := range vec{
					ldvec := make([]float64, len(ld))
					for j, val := range ld{
						if j == 2 || j == 3{
							val = sf * val
						}
						ldvec[j] = val
					}
					//ldvec = append(ldvec, float64(lt))
					mldcs[lc] = append(mldcs[lc],ldvec)
				}
			}
		}
	}

	return
}
/*

	for idx := 1; idx <= len(mod.Mprp); idx++{
		//log.Println(ColorYellow, "member->",idx)
		//mem := mod.Ms[idx]
		//log.Println("Max loads->\n",ColorRed,"tensile->",mem.Tmax, "N",ColorCyan,"kompressive->",mem.Cmax,"N",ColorReset)
		//log.Println("Pf->",mem.Pu)
	}


	// log.Println(ColorYellow)
	// for i, lc := range jldcs{
	// 	log.Println("lp->",i)
	// 	for _, ld := range lc{
	// 		log.Println("load->",ld)
	// 	}
	// }
	// log.Println(ColorReset)
*/
