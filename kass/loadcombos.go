package barf

//load combination calculation functions

import (
	"fmt"
	"errors"
)

//CalcLdCombos generates and calcs load combinations based on mod.PSFs / mod.Code
//it calls GenLdCombos for each lp index and chills
func CalcLdCombos(mod *Model) (err error){
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
			modlp := *mod
			//fmt.Println("lp->",lp, jlds)
			modlp.Id += fmt.Sprintf("lp_%v_%.1f",lp,i)
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
			modlp := *mod
			modlp.Id += fmt.Sprintf("lp-%v-%.1f",lp,i)
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
	for idx := 1; idx <= len(mod.Mprp); idx++{
		//fmt.Println(ColorYellow, "member->",idx)
		//mem := mod.Ms[idx]
		//fmt.Println("Max loads->\n",ColorRed,"tensile->",mem.Tmax, "N",ColorCyan,"kompressive->",mem.Cmax,"N",ColorReset)
		//fmt.Println("Pf->",mem.Pu)
	}
	
	return
}

//GenLdCombos generates load combos for 3 load cases now - 1 - dl, 2 - ll, 3 - wl
//nwl, nsl - number of wind load and seismic load cases
func GenLdCombos(nwl, nsl int, jldsrv, mldsrv map[float64][][]float64, psfs [][]float64)(mldcs, jldcs map[float64][][]float64){
	mldcs = make(map[float64][][]float64)
	jldcs = make(map[float64][][]float64)
	for i, sfs := range psfs{
		var wlon, slon bool
		//fmt.Println(sfs)
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
			//fmt.Println("ncases->",nwl)
			//fmt.Println("lc number->",lc)
			for j := 0; j < nwl; j++{
				//if j > 2{continue}
				wlc := 3.0 + float64(j+1)/10.0
				lc = float64(i+1) + float64(j+1)/10.0
				sf := sfs[2]
				if len(jldsrv[wlc]) > 0{
					fmt.Println("wind load case",wlc,"\n",jldsrv[wlc])
					
					for _, ld := range jldsrv[wlc]{
						ldvec := make([]float64, len(ld))
						ldvec[0] = ld[0]
						for k, val := range ld[1:]{
							val = sf * val
							ldvec[k+1] = val
						}
						ldvec = append(ldvec, float64(3))
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
						ldvec = append(ldvec, float64(3))
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
						ldvec = append(ldvec, float64(lt))
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
						ldvec = append(ldvec, float64(lt))
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
					ldvec = append(ldvec, float64(lt))
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
					ldvec = append(ldvec, float64(lt))
					mldcs[lc] = append(mldcs[lc],ldvec)
				}
			}
		}
	}
	/*
	fmt.Println(ColorYellow)
	for i, lc := range jldcs{
		fmt.Println("lp->",i)
		for _, ld := range lc{
			fmt.Println("load->",ld)
		}
	}
	fmt.Println(ColorReset)
	*/
	return
}
