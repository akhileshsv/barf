package barf

import (
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
		} else {
			if len(mod.Js[n].Dr) == 0{
				mod.Js[n].Dr = make([][]float64,mod.Nlp)
				mod.Js[n].Rr = make([][]float64,mod.Nlp)				
				
			}
		}
		mod.Js[n].Dr[lp-1] = node.Displ
		mod.Js[n].Rr[lp-1] = node.React
		switch mod.Frmtyp{
			case 1:
			case 2:
			case 3:
			case 4:
			case 5:
			case 6:
		}
	}
	//return
}

//MemRez reads a load pattern's results (frmrez) and stores in base joint and member result maps
func (mod *Model) MemRez(lp int, frmrez []interface{}){
	//frm types - 1- beam, 2- 2d truss, 3 - 2d frame, 4 - 3d truss, 5 - 3d grid, 6 - 3d frame
	switch mod.Frmtyp{
		case 1:
		//beam
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
				if mod.Ms[m].Cmax == 0.0{
					mod.Ms[m].Cmax = mem.Qf[0]
				} else if mod.Ms[m].Cmax < mem.Qf[0]{
					mod.Ms[m].Cmax = mem.Qf[0]
				}
			} else {
				if mod.Ms[m].Tmax == 0.0{
					mod.Ms[m].Tmax = math.Abs(mem.Qf[0])
				} else if mod.Ms[m].Tmax < math.Abs(mem.Qf[0]){
					mod.Ms[m].Tmax = math.Abs(mem.Qf[0])
				}
			}
		}

		case 2:
		//2d truss
		//plot axial force, nodal global displacement
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
				if mod.Ms[m].Cmax == 0.0{
					mod.Ms[m].Cmax = mem.Qf[0]
				} else if mod.Ms[m].Cmax < mem.Qf[0]{
					mod.Ms[m].Cmax = mem.Qf[0]
				}
			} else {
				if mod.Ms[m].Tmax == 0.0{
					mod.Ms[m].Tmax = math.Abs(mem.Qf[0])
				} else if mod.Ms[m].Tmax < math.Abs(mem.Qf[0]){
					mod.Ms[m].Tmax = math.Abs(mem.Qf[0])
				}
			}
		}
		case 3:
		//2d frame
		//get shear and bending moment diagrams and stuff

		js, _ := frmrez[0].(map[int]*Node)
		ms, _ := frmrez[1].(map[int]*Mem)
		report,_ := frmrez[6].(string)
		if len(mod.Reports) == 0{
			mod.Reports = append(mod.Reports, report)
		} else {
			mod.Reports[lp-1] = report
		}
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
				if mod.Ms[m].Cmax == 0.0{
					mod.Ms[m].Cmax = mem.Qf[0]
				} else if mod.Ms[m].Cmax < mem.Qf[0]{
					mod.Ms[m].Cmax = mem.Qf[0]
				}
			} else {
				if mod.Ms[m].Tmax == 0.0{
					mod.Ms[m].Tmax = math.Abs(mem.Qf[0])
				} else if mod.Ms[m].Tmax < math.Abs(mem.Qf[0]){
					mod.Ms[m].Tmax = math.Abs(mem.Qf[0])
				}
			}
		}
		// switch mod.Npsec{
		// 	case false:
		// 	case true:
		// }
	}
	//return
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
