package barf

import (
	"fmt"
	kass"barf/kass"
)


/*
"Id":"2.3Hulse",
 "Cmdz":["1db","mks","1","dl,ll","mredis,0.30"], 
 "Ncjt":2, 
 "Coords": [[0],[6],[10],[16]],
 "Supports":[[1,-1,0],[2,-1, 0],[3,-1,0],[4,-1,0]],
 "Em":[[25e9]],"Cp":[[1200e-9]],
 "Mprp": [[1,2,1,1,0], [2,3,1,1,0], [3,4,1,1,0]],
 "Jloads": [], 
 "Msloads":[[1,3,25,0,0,0,1],[1,3,10,0,0,0,2],[2,3,25,0,0,0,1],[2,3,10,0,0,0,2], [3,3,25,0,0,0,1],[3,3,10,0,0,0,2]],
 "PSFs":[1.4,1.0,1.6,0.0],
 "Clvrs":[[0,0],[0,0]]
*/

//CBeamPscBm calcs secondary bending moments due to prestressing in a continuous beam
func (cb *CBm) PscBm()(err error){
	//first - calc curve constants
	var mod kass.Model
	if cb.Pexctyp == ""{
		cb.Pexctyp = "parabolic"
	}
	mod.Em = cb.Em
	if len(cb.Em) == 0{
		if cb.Fck == 0.0{cb.Fck = 25.0}
		mod.Em = [][]float64{{kass.FckEm(cb.Fck)}}
	}
	if len(cb.Cp) != 1 && len(cb.Cp) != len(cb.Lspans){
		err = fmt.Errorf("invalid length of cp vec - %v for nspans %v",len(cb.Cp), len(cb.Lspans))
		return
	}
	var pexts [][]float64
	var x, iz, ar float64
	var cp int
	nms := make([]float64, len(cb.Lspans)+1)
	for i, lspan := range cb.Lspans{
		mod.Coords = append(mod.Coords, []float64{x})
		x += lspan
		if len(cb.Cp) == 1{
			iz = cb.Cp[0][0]
			cp = 1	
			if len(cb.Cp[0]) > 1{
				ar = cb.Cp[0][1]
			} 
		} else {
			iz = cb.Cp[i][0]
			cp = i+1
			if len(cb.Cp[i]) > 1{
				ar = cb.Cp[i][1]
			} 
		}
		mod.Cp = append(mod.Cp, []float64{iz, ar})
		pfrc := cb.Pfrcs[i]
		pexs := cb.Pexs[i]
		fmt.Println("i, lspan-",i, lspan)
		fmt.Println("p frc",pfrc," kn ext",pexs)
		//func PscBmExts(ctyp string, lspan, pfrc float64, pex []float64)(exts []float64, err error)
		pex, e := PscBmExts(cb.Pexctyp, lspan, pfrc, pexs)
		if e != nil{
			err = fmt.Errorf("error generating tendon profile for span %v %s",i+1, e)
			return
		}
		pexts = append(pexts, pex)
		//generate fixed end moments
		//PscBmAr(ndiv int, lspan, pfrc float64, pex []float64)(ar, xc, ml, mr float64)
		ar, xc, ml, mr := PscBmAr(21, lspan, pfrc, pex)
		fmt.Println("area under moment curve-",ar, "xc - ",xc)
		fmt.Println("fixed end moments - ml - ",ml, " mr - ",mr)
		//mlvec := []float64{float64(i+1), 0, 0, ml}
		//mrvec := []float64{float64(i+2), 0, 0, mr}
		nms[i] += ml
		nms[i+1] += mr
		mod.Mprp = append(mod.Mprp, []int{i+1, i+2, 1, cp, 0})
	}
	for i, nm := range nms{
		mdx := i+1
		mod.Supports = append(mod.Supports, []int{i+1, -1, 0})
		switch i{
			case len(nms)-1:
			lspn := cb.Lspans[len(cb.Lspans)-1]
			mod.Msloads = append(mod.Msloads, []float64{float64(mdx-1), 2, nm, 0, lspn, 0})
			default:
			mod.Msloads = append(mod.Msloads, []float64{float64(mdx), 2, nm, 0, 0, 0})
		}
		
	}
	fmt.Println("NODAL MOMENTS-",nms)
	mod.Coords = append(mod.Coords, []float64{x})
	mod.Frmstr = "1db"
	mod.Term = cb.Term
	mod.Web = cb.Web
	var frmrez []interface{}
	frmrez, err = kass.CalcBm1d(&mod, 2)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(mod.Report)
	_ = kass.CalcBmSf(&mod, frmrez,true)
	fmt.Println("finito-this is so horribly wreng")
	return
}
