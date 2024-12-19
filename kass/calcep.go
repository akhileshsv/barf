package barf

import (
	"math"
	"log"
	"fmt"
	"errors"
)

//CalcEpFrm performs elastic plastic analysis of plane frames
//see sec 13, harrison
//INCOMPLTE - to do deflection calcs
//also just prints stuff doesn't return anything
//plastic moment of a member - cpvec[-1]
//memrel - mod.Mprp[4]; 0 - regular, 1 - hinge at beginning, 2 - hinge at end, 3 - both ends hinged

func CalcEpFrm(mod *Model) (err error){
	var frmrez []interface{}
	var e error
	var m1, m2, m3, m4, mlf, clf float64 //min. load factor, cumulative load factor
	var ms map[int]*Mem
	var js map[int]*Node
	if mod.Ncjt == 2{mod.Frmtyp = 2} else {mod.Frmtyp = 3}
	//var report string
	hinges := [][]int{} //member, node
	rcap := make([][]float64, len(mod.Mprp))//residual capacity
	lfac := make([][]float64, len(mod.Mprp)) //load factor at each hinge
	elm := make([][]float64, len(mod.Mprp)) //elastic moments (mb, me) at each member
	mrez := make([][]float64, len(mod.Mprp)) //cumulative moments at each member
	mprev := make([][]float64, len(mod.Mprp)) //check for sign change from prev iteration
	jdisp := make([][]float64, len(mod.Coords)) //cumulative joint displacements
	elp := make([]float64, len(mod.Mprp)) //axial force at each member
	prez := make([]float64, len(mod.Mprp)) //res. axial force at each member
	plm := make([]float64, len(mod.Mprp)) //plastic moment of each member
	var mlfs [][]float64 //min. load factors
	for i := range mod.Coords{
		jdisp[i] = make([]float64, mod.Ncjt)
	}
	for i := range rcap{
		rcap[i] = make([]float64, 2)
		lfac[i] = make([]float64, 2)
		elm[i] = make([]float64, 2)
		mprev[i] = make([]float64, 2)
		mrez[i] = make([]float64, 2)
	}
	//get moment capacity (min. plastic moment) at each end of each member
	for mem, mprp := range mod.Mprp{
		
		if len(mprp) < 5{
			err = fmt.Errorf("for member - %v : invalid property slice %v",mem, mprp)
			return
		}
		var mp float64
		switch mod.Frmstr{
			case "1db":
			cp := mprp[3]
			if len(mod.Cp) < cp || len(mod.Cp[cp-1]) < 2{
				err = fmt.Errorf("for member - %v : invalid cp index/data %v in cpvec\n %v",mem, cp, mod.Cp)
				return
			}
			mp = mod.Cp[cp-1][1]
			case "2df":
			cp := mprp[3]
			if len(mod.Cp) < cp || len(mod.Cp[cp-1]) < 3{
				err = fmt.Errorf("for member - %v : invalid cp index/data %v in cpvec\n %v",mem, cp, mod.Cp)
				return
			}
			mp = mod.Cp[cp-1][2]
		}
		rcap[mem][0] = mp; rcap[mem][1] = mp
		plm[mem] = mp
		
	}
	//fmt.Println("\n***\n***\n***starting iter***\n***\n***\n***")
	iter, kiter := 0, 0
	for iter != -1{
		kiter++
		if mod.Verbose{fmt.Println(ColorYellow,"starting iter number->",kiter,ColorReset)}
		if kiter > 300{
			log.Println("ERRORE, errore->max iterations reached")
			err = errors.New("iteration error")
			iter = -1
			break
		}
		switch mod.Frmstr{
			case "1db":
			frmrez, e = CalcBm1d(mod, mod.Ncjt)
			case "2df":
			frmrez, e = CalcFrm2d(mod,mod.Ncjt)	
		}
		if e != nil{
			// err = e
			if mod.Verbose{
				fmt.Println(ColorRed,"***KOLLAPSE***",ColorReset)
				fmt.Println(ColorGreen,e,"\niteration terminated",ColorReset)
			}
			iter = -1
			break
		}
		//hinge formed at each stage - mh (mem), hj (node)
		var mh, hj int
		mlf = -6.66
		js , _ = frmrez[0].(map[int]*Node)
		ms , _ = frmrez[1].(map[int]*Mem)
		//report, _ = frmrez[6].(string)
		for mem := range ms{
			jb, je := ms[mem].Mprp[0], ms[mem].Mprp[1]
			switch mod.Frmstr{
				case "1db":
				elp[mem-1] = ms[mem].Qf[0]
				m1 = math.Abs(ms[mem].Qf[1]); m2 = math.Abs(ms[mem].Qf[3])
				m3 = ms[mem].Qf[1]; m4 = ms[mem].Qf[3]
				case "2df":
				elp[mem-1] = ms[mem].Qf[0]
				m1 = math.Abs(ms[mem].Qf[2]); m2 = math.Abs(ms[mem].Qf[5])
				m3 = ms[mem].Qf[2]; m4 = ms[mem].Qf[5]
			}
			if m1 != 0.0{
				lf := rcap[mem-1][0]/m1
				lfac[mem-1][0] = lf
				
				switch{
					//case m3 * mprev[mem-1][0] < 0.0:
					case kiter > 1 && (math.Abs(m3) < math.Abs(mprev[mem-1][0]) || m3 * mprev[mem-1][0] < 0.0):
					//do nothing at all - see page 12 of hson paper
					//sign change indicates a reversal of earlier moment direction so goot?
					case mlf == -6.66 && lf > 0.0:
					mlf = lf; mh = mem; hj = jb
					case mlf > lf && lf > 0.0:
					mlf = lf; mh = mem; hj = jb
					case mlf == lf && lf > 0.0:
					if rcap[mem-1][0] + rcap[mem-1][1] < rcap[mh-1][0] + rcap[mh-1][1]{
						mlf = lf; mh = mem; hj = jb
					}

					cpmem := ms[mem].Mprp[3]; cpmh := ms[mh].Mprp[3]
					switch{
						case cpmem == cpmh:	
						//compare rcap, choose the one with greater reserve capacity SAYS WHO OR WHAT
						if rcap[mem-1][0] + rcap[mem-1][1] < rcap[mh-1][0] + rcap[mh-1][1]{
							mlf = lf; mh = mem; hj = jb
						}
						case cpmem != cpmh:
						//now check for plastic moment capacity, choose the weak member
						switch mod.Frmstr{
							case "1db":
							if mod.Cp[cpmem-1][1] < mod.Cp[cpmh-1][1]{
								mlf = lf; mh = mem; hj = jb
							}
							case "2df":
							if mod.Cp[cpmem-1][2] < mod.Cp[cpmh-1][2]{
								mlf = lf; mh = mem; hj = jb
							}
						}
					}
				}
			} else {
				lfac[mem-1][0] = 0.0
			}
			if m2 != 0.0{
				lf := rcap[mem-1][1]/m2
				lfac[mem-1][1] = lf
				switch{
					case kiter > 1 && (math.Abs(m4) < math.Abs(mprev[mem-1][1]) || m4 * mprev[mem-1][1] < 0.0):
					//chew bubblegum or something
					case mlf == -6.66 && lf > 0.0:
					mlf = lf; mh = mem; hj = je
					case mlf > lf && lf > 0.0:
					mlf = lf; mh = mem; hj = je
					case mlf == lf && lf > 0.0:
					cpmem := ms[mem].Mprp[3]; cpmh := ms[mh].Mprp[3]
					switch{
						case cpmem == cpmh:	
						if rcap[mem-1][0] + rcap[mem-1][1] < rcap[mh-1][0] + rcap[mh-1][1]{
							mlf = lf; mh = mem; hj = je
						}
						case cpmem != cpmh:
						switch mod.Frmstr{
							case "1db":
							if mod.Cp[cpmem-1][1] < mod.Cp[cpmh-1][1]{
								mlf = lf; mh = mem; hj = je
							}
							case "2df":
							if mod.Cp[cpmem-1][2] < mod.Cp[cpmh-1][2]{
								mlf = lf; mh = mem; hj = je
							}
						}
					}
				}
			} else {
				lfac[mem-1][1] = 0.0
			}
			elm[mem-1] = []float64{m1, m2}
			mprev[mem-1] = []float64{m3,m4}
		}
		
		//fmt.Println("elastic moments->",elm)
		//fmt.Println("rcap in->",rcap)
		clf += mlf
		mlfs = append(mlfs, []float64{mlf,clf})
		for i, node := range js{
			//t.Println(ColorRed,"node->",i,ColorReset)
			for j, disp := range node.Displ{
				//fmt.Println("j-",j,"displ-",disp)
				jdisp[i-1][j] += disp * mlf
			}
		}
		if mod.Verbose{
			fmt.Println("iter no->",kiter,"minimum load factor->",mlf, "cumulative load factor->",clf)
			fmt.Println("plastic hinge formed in mem->", mh, "near node->",hj)
		}
		hinges = append(hinges, []int{mh, hj})
		
		//multiply elastic moments by min lf and subtract to get the residual capacity
		for m := range rcap{
			m1, m2 := elm[m][0], elm[m][1]
			if rcap[m][0] > 0.0{
				rcap[m][0] -= mlf * m1
			}
			if rcap[m][1] > 0.0{
				rcap[m][1] -= mlf * m2
			}
			
			switch mod.Frmstr{
				case "1db":
				mrez[m][0] += mlf * ms[m+1].Qf[1]
				mrez[m][1] += mlf * ms[m+1].Qf[3]
				case "2df":
				mrez[m][0] += mlf * ms[m+1].Qf[2]
				mrez[m][1] += mlf * ms[m+1].Qf[5]
			}
			prez[m] += mlf * elp[m]
		}
		//insert hinge in member mh at node hj
		var mrel int
		if hj == mod.Mprp[mh-1][0]{
			switch mod.Mprp[mh-1][4]{
				case 0:
				mrel = 1
				case 2:
				mrel = 3
			}
		} else {
			switch mod.Mprp[mh-1][4]{
				case 0:
				mrel = 2
				case 1:
				mrel = 3
			}
		}
		mod.Mprp[mh-1][4] = mrel
		
	}
	if mod.Spam{
		//fmt.Println(nodeh)
		fmt.Println("rcap->",rcap)
		//fmt.Println(mod.Mprp)
		fmt.Println("elastic moments->",elm)
		fmt.Println("terminal moments->",mrez)
		fmt.Println("plastic moments->",plm)
		fmt.Println("axial loads->",prez)
		fmt.Println("sequence of hinge formation")
		rstr := ""
		for i, h := range hinges{
			rstr += fmt.Sprintf("[member %v, node %v]",h[0],h[1])
			if i != len(hinges) - 1{
				rstr += "->"
			}
		}
		fmt.Println(rstr)
		//fmt.Println(report)
		for i, disp := range jdisp{
			fmt.Println("node->",i+1, "displacement->",disp)
			copy(js[i+1].Displ,disp)
			//fmt.Println("node->",i+1, "from js->",js[i+1].Displ)
			//for _, d := range js[i+1].Displ{
			//	fmt.Println("mul->", d * clf, d * mlf)
			//}
		}
	}
	pltchn := make(chan string, 1) 
	switch mod.Frmstr{
		case "1db":
		mod.Report = BmEpTable(hinges, mlfs, mrez, plm, prez, js, ms)
		if mod.Term != ""{
			go PlotBm1d(mod, mod.Term, pltchn)
		}
		// mod.Reports = append(mod.Reports, BmEpTable(hinges, mlfs, mrez, plm, prez, js, ms))
		case "2df":
		mod.Report = Frm2dEpTable(hinges, mlfs, mrez, plm, prez, js, ms)
		if mod.Term != ""{
			go PlotFrm2d(mod, mod.Term, pltchn)
		
		}
		// mod.Reports = append(mod.Reports, Frm2dEpTable(hinges, mlfs, mrez, plm, prez, js, ms))
		// if !mod.Web{fmt.Println(mod.Reports[0])}
	}
	if !mod.Web{fmt.Println(mod.Report)}
	if mod.Term != ""{
		pltstr := <- pltchn
		mod.Txtplots = append(mod.Txtplots, pltstr)
	}
	//frmrez[6] = Frm2dTable(js, ms, dglb, rnode, nsc, ndof, ncjt)
	return
}

/*
				if mlf == -6.66 && lf > 0.0{
					mlf = lf
					mh = mem; hj = jb
				} else if mlf > lf && lf > 0.0{
					mlf = lf
					mh = mem; hj = jb
				} else if mlf == lf && lf > 0.0{
					//kompare mh rcap sum and mem rcap sum, choose lower reserve value mem
					if m1 + m2 < rcap[mh-1][0] + rcap[mh-1][1]{
						mlf = lf; mh = mem; hj = jb
					}
				}
*/
