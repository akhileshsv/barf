package barf

import (
	"fmt"
	"math"
)

//NodeRez reads a load pattern's joint results and stores in base joint map
func (mod *Model) NodeRez(lp int, js map[int]*Node){
	for n, node := range js{
		if _, ok := mod.Js[n]; !ok{
			mod.Js[n] = &Node{}
			*mod.Js[n] = *node
			mod.Js[n].Dr = make([][]float64,mod.Nlp)
			mod.Js[n].Rr = make([][]float64,mod.Nlp)
		}
		mod.Js[n].Dr[lp-1] = node.Displ
		mod.Js[n].Rr[lp-1] = node.React
	}
	return
}


//NodeRez reads a load pattern's results (frmrez) and stores in base joint and member result maps
func (mod *Model) MemRez(lp int, frmrez []interface{}){
	//frm types - 1- beam, 2- 2d truss, 3 - 2d frame, 4 - 3d truss, 5 - 3d grid, 6 - 3d frame
	switch mod.Frmtyp{
		case 1:
		//beam
		case 2:
		//2d truss
		js, _ := frmrez[0].(map[int]*Node)
		ms, _ := frmrez[1].(map[int]*Mem)
		report,_ := frmrez[6].(string)
		mod.Reports[lp-1] = report
		
		//nodes
		for n, node := range js{
			if _, ok := mod.Js[n]; !ok{
				mod.Js[n] = &Node{}
				*mod.Js[n] = *node
				mod.Js[n].Dr = make([][]float64,mod.Nlp)
				mod.Js[n].Rr = make([][]float64,mod.Nlp)
			}
			mod.Js[n].Dr[lp-1] = node.Displ
			mod.Js[n].Rr[lp-1] = node.React
		}
		//mems
		for m, mem := range ms{
			if _, ok := mod.Ms[m]; !ok{
				mod.Ms[m] = &Mem{}
				*mod.Ms[m] = *mem
				mod.Ms[m].Qfr = make([][]float64,mod.Nlp)
				mod.Ms[m].Gfr = make([][]float64,mod.Nlp)
				mod.Ms[m].Pu = make([]float64,mod.Nlp)
			}
			mod.Ms[m].Qfr[lp-1] = mem.Qf
			mod.Ms[m].Gfr[lp-1] = mem.Gf
			mod.Ms[m].Pu[lp-1] =  mem.Qf[0]
			if mem.Qf[0] > 0.0{
				if mem.Cmax == 0.0{
					mem.Cmax = mem.Qf[0]
				} else if mem.Cmax < mem.Qf[0]{
					mem.Cmax = mem.Qf[0]
				}
			} else {
				if mem.Tmax == 0.0{
					mem.Tmax = mem.Qf[0]
				} else if math.Abs(mem.Tmax) < math.Abs(mem.Qf[0]){
					mem.Tmax = mem.Qf[0]
				}
			}
		}
		case 3:
		//2d frame
		//get shear and bending moment diagrams and stuff
		switch mod.Npsec{
			case false:
			case true:
		}
	}
	return
}

//SumFrcs konsolidates forces by location n typ
func (mod *Model) SumFrcs() (jl, ml [][]float64){
	jlmap := make(map[float64][]float64)
	mlmap := make(map[float64]map[float64][]float64)
	for _, load := range mod.Jloads{
		//fmt.Println("nodal load->no, load",i+1, load)
		n := load[0]
		if _, ok := jlmap[n]; !ok{
			jlmap[n] = make([]float64, len(load))
			jlmap[n][0] = n
		} 
		for i, val := range load[1:]{
			jlmap[n][1+i] += val
		}
	}
	for _, load := range mod.Msloads{
		//fmt.Println("member load->no, load",i+1, load)
		//fmt.Println("adding->",load)
		m := load[0]; lt := load[1]
		if _, ok := mlmap[m]; !ok{
			mlmap[m] = make(map[float64][]float64)
		}
		if _, ok := mlmap[m][lt]; !ok{
			mlmap[m][lt] = make([]float64, len(load))
			mlmap[m][lt][0] = m
			mlmap[m][lt][1] = lt 
		}	
		for i, val := range load[2:]{
			mlmap[m][lt][i+2] += val
		}
		//fmt.Println("added->",mlmap[m][lt])
	}
	for _, load := range jlmap{
		jl = append(jl, load)
	}
	for _, lt := range mlmap{
		for _, load := range lt{
			//fmt.Println("load vec->",load)
			ml = append(ml, load)
		}
	}
	return
}

//CalcRezNp processes np model results
func (mod *Model) CalcRezNp(frmtyp string, frmrez []interface{}){
	ms, _ := frmrez[1].(map[int]*MemNp)
	msloaded,_ := frmrez[5].(map[int][][]float64)
	for mem := 1; mem <= len(ms); mem++{
		var ldcases [][]float64
		var r BeamRez
		if val, ok := msloaded[mem]; ok{
			ldcases = val
		} else {
			continue
		}
		fmt.Println("loadcases->",ldcases)
		//get new func to calc deflections - THIS IS SO WRONG
		switch frmtyp{
			case "1db":
			r = Bmsfcalc(mem, ldcases, ms[mem].Lspan, ms[mem].Em, ms[mem].A0, ms[mem].I0, ms[mem].Qf[0], ms[mem].Qf[1], ms[mem].Qf[2], ms[mem].Qf[3], true)
			case "2df":
			
			r = Bmsfcalc(mem, ldcases, ms[mem].Lspan, ms[mem].Em, ms[mem].A0, ms[mem].I0, ms[mem].Qf[1], ms[mem].Qf[2], ms[mem].Qf[4], ms[mem].Qf[5], true)
		}
		fmt.Println("span->",mem)
		fmt.Println("mem->",ms[mem].Mprp)
		for j, bm := range r.BM{
			fmt.Println("section->",j,"moment->",bm," kn-m"," shear->",r.SF[j]," kn", " def->", 1e3 * r.Dxs[j]," mm")
		}
		//s := PlotBmSfBm(r.Xs, r.SF, r.BM, r.Dxs, ms[mem].Lspan, true)
		//fmt.Println(s)
	}
	reportz, _ := frmrez[6].(string)
	fmt.Println(reportz)
	return
}

/*
func (mem *Mem) Qfrez(lp int, qf []float64){
	
}

func (mod *Model) MRez(lp, ms map[int]*Mem){
	switch mod.Npsec{
		case true:
		//np sec rez
		case false:
		//basic rez
		for m, mem := range ms{
			if _, ok := mod.Ms[m]; !ok{
				mod.Ms[m] = &Mem{}
				*mod.Ms[m] = *mem
				mod.Ms[m].Qfr = make([][]float64,mod.Nlp)
				mod.Ms[m].Gfr = make([][]float64,mod.Nlp)
				mod.Ms[m].Pu = make([]float64,mod.Nlp)
			}
			mod.Ms[m].Qfr[lp-1] = mem.Qf
			mod.Ms[m].Gfr[lp-1] = mem.Gf
			mod.Ms[m].Pu[lp-1] =  mem.Qf[0]
			if mem.Qf[0] > 0.0{
				if mem.Cmax == 0.0{
					mem.Cmax = mem.Qf[0]
				} else if mem.Cmax < mem.Qf[0]{
					mem.Cmax = mem.Qf[0]
				}
			} else {
				if mem.Tmax == 0.0{
					mem.Tmax = mem.Qf[0]
				} else if math.Abs(mem.Tmax) < math.Abs(mem.Qf[0]){
					mem.Tmax = mem.Qf[0]
				}
			}
			
		}
	}
}

*/
