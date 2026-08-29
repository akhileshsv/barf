package barf

import (
	//"log"
	//"fmt"
	"math"
)

//rcc helper funcs

//FckEm takes in concrete grade fck and returns elastic modulus in KN/m2
func FckEm(fck float64) (float64){
	return 5e6 * math.Sqrt(fck) 
}

//GetBfRcc returns the flange width of an rcc beam
func GetBfRcc(code, mrel int, tyb, lspan, df, bw float64) (bf float64){
	//konvert lspan to mm (this will go horribly wrong somewhere)
	lspan = lspan * 1e3
	if mrel == 0{lspan = 0.7 * lspan}
	switch code{
		case 1:
		//is code
		bf = math.Round(tyb*(6.0*df + lspan/6.0) + bw)
		case 2:
		//bs code
		bf = math.Round(bw + (lspan * tyb)/5.0)
		case 3:
		//calc using aci
		
	}
	return
}

//CalcBf takes in model column/beam sections and data
//calculates bf for individual beams and returns updated section data
//breadth of flange calc for a bm
func CalcBf(code, bstyp, mrel, nspans int, dslb float64, cdx, bdx, styps []int, lspans []float64, sections [][]float64) (secvec[][]float64, sts []int){
	//calc breadth of flange 
	colvec := make([][]float64, nspans + 1)
	bmvec := make([][]float64, nspans)
	ctyps := make([]int, nspans+1)
	btyps := make([]int, nspans)
	sts = make([]int, 2*nspans+1)
	
	if len(cdx) == 1{
		idx := cdx[0]
		cdx = make([]int, nspans+1)
		for i := range cdx{
			cdx[i] = idx
		}
	}
	if len(cdx) == 0{
		colvec = [][]float64{}
		ctyps = []int{}
	} else {
		for i, idx := range cdx{
			colvec[i] = make([]float64, len(sections[idx-1]))
			copy(colvec[i], sections[idx-1])
			ctyps[i] = styps[idx-1]
		}
	}
	if len(bdx) == 1{
		idx := bdx[0]
		bdx = make([]int, nspans)
		for i := range bdx{
			bdx[i] = idx
		}
	}
	var bf, df, bw, dused float64
	for i, idx := range bdx{
		//log.Println("from bdx->i, idx->",i, idx,"sections-",sections)
		bdim := sections[idx-1]
		var btyp int
		var tyb float64
		if bstyp == 0 || bstyp == 1{
			bst := styps[idx-1]
			switch bst{
				case 6,7,8,9,10:
				btyp = bst
				default:
				btyp = 6
			}
		} else {
			btyp = bstyp
		}
		switch btyp{
			case 6:
			tyb = 1.0
			case 7,8,9,10:
			tyb = 0.5
		}
		btyps[i] = btyp
		bf = GetBfRcc(code, mrel, tyb, lspans[i],dslb,bdim[0])
		df = dslb; bw = bdim[0]; dused = bdim[1]
		bmvec[i] = []float64{bf,dused,bw,df}
	}
	secvec = append(colvec,bmvec...)
	sts = append(ctyps, btyps...)
	return
}

