package barf

import (
	//"log"
	"fmt"
	"math"
	"sort"
	"errors"
)

//rebar selection funcs for beams

//RbrSide calculates the side face reinforcement required for a deep beam
func (b *RccBm) RbrSide() (err error){
	if b.Dused <= 750.0{
		return
	}
	//0.1 percent of web area as per is code
	b.Aside = 0.1 * b.Dused * b.Bw/100.0
	spc := b.Bw
	switch b.Code{
		case 1:
		if spc > 300.0{spc = 300.0}
		case 2:
		if spc > 250.0{spc = 250.0}
	}
	dmin := math.Sqrt(spc * b.Bw/b.Fy)
	for _, dia := range []float64{10.0,12.0,16.0,20.0,25.0,28.0,32.0}{
		if dia >= dmin{
			dmin = dia
			break

		}
	}
	//fmt.Println("spcing, dia->",spc, dmin)
	nrows := math.Ceil((b.Dused - b.Cvrc - b.Cvrt)/spc)
	spc = math.Round((b.Dused - b.Cvrc - b.Cvrt)/nrows)
	nrows = nrows - 1.0
	nbars := nrows * 2.0
	//fmt.Println("nbarz->",nbars)
	for _, dia := range []float64{10.0,12.0,16.0,20.0,25.0,28.0,32.0}{
		atot := nbars * RbrArea(dia)
		if atot >= b.Aside && dia >= dmin{
			//rez = []float64{0 float64(n1), 1 float64(n2),2 d1,3 d2,4 astmin,5 ast,6 adiff}
			//rez = []float64{7 nlayer,8 astprov,9 efcvr,10 efdp,11 cldis,12 clvdis,13 nbarRow} 
			b.Rbrside = []float64{nbars,0.0,dia,0.0,atot,atot,0.0, nrows, atot, 0.0, 0.0, b.Bw - b.Cvrc -b.Cvrt, spc, 2.0}
			break
		}
	}
	if len(b.Rbrside) == 0{
		err = errors.New("side face rebar error wtf")
		return
	}
	var xt, yt float64

	switch b.Styp{
		case 1:
		xt = b.Cvrt; yt = b.Cvrt //xc = b.Cvrc; yc = b.Dused - b.Cvrc
		case 7:
		xt = b.Cvrt; yt = b.Cvrt//; xc = b.Cvrc; yc = b.Dused - b.Cvrc
		case 6:
		xt = b.Cvrt + b.Bf/2.0 - b.Bw/2.0; yt = b.Cvrt//; xc = xt; yc = b.Dused - b.Cvrc
		case 14:
		xt = b.Cvrt + b.Bw/2.0 - b.Bf/2.0; yt = b.Cvrt//; xc = xt; yc = b.Dused - b.Cvrc
		default:
		//HUH? HUH?
		xt = b.Cvrt; yt = b.Cvrt
	}
	xs := xt
	ys := yt
	for i := 1; i <= int(nrows); i++{
		ys += spc
		xe := xs + b.Bw - 2.0 * b.Cvrt
		b.Dias = append(b.Dias, dmin)
		b.Dias = append(b.Dias, dmin)
		b.Barpts = append(b.Barpts, []float64{xs,ys})
		b.Barpts = append(b.Barpts, []float64{xe,ys})
	}
	return
}

//BarNRows returns the nummber of bars/rows
//(capn obvious over and out)
func BarNRows(bw, dused float64, r []float64) ([]float64){
	n1, n2, d1, d2 := r[0], r[1], r[2], r[3]
	dmax := d2; if d2 < d1 || n2 == 0.0 {dmax = d1}
	nbar := n1 + n2
	rez := RbrNRows(bw, dused, nbar, dmax) 
	r = append(r, rez...)
	return r
}

//BarSort sorts bar rez by nlayers and astprov
func (b *RccBm) BarSort(rez [][]float64, tcdx int) (err error, mincvr float64){
	layr := [][]float64{}
	var nopts int
	mincvr = 666.0
	for _, r := range rez{
		rlay := BarNRows(b.Bw, b.Dused, r)
		//fmt.Println("nlayers->",rlay[7])
		//if rlay[0]
		if len(rlay) == 0{continue}
		efcvrc := rlay[9]
		if efcvrc <= 0.0 || efcvrc == 666.0{
			continue
		}
		efcvr := b.Cvrt
		switch tcdx{
			case 2:
			//efcvr = b.Cvrc
			//THIS HAS NO BEARING ON EFFECTIVE DEPTH right
			efcvr = efcvrc
		}
		if mincvr > efcvrc{
			mincvr = efcvrc
		}
		//fmt.Println("efcvrc->",efcvrc,"efcvr->",efcvr)
		//fmt.Println("dused->",b.Dused,)
		//if math.Abs(efcvrc - efcvr) < 5.0{layr = append(layr, rlay)}
		if efcvr - efcvrc >= 0.0 && rlay[7] <= 3.0{
			//fmt.Println("glorious SUCCESS!", efcvrc, "calc vs act.",efcvr)
			layr = append(layr, rlay)
			nopts++
		}
	}
	if len(layr) == 0 || nopts == 0{
		err = errors.New(fmt.Sprint("effective cover error",tcdx))
		return
	}
	sort.Slice(layr, func(i, j int) bool {
		if layr[i][7] != layr[j][7]{
			return layr[i][7] < layr[j][7]
		}
		return layr[i][4] < layr[j][4]
	})
	//fmt.Println("CHOICE->",layr[0])
	switch tcdx{
		case 1:
		b.Rbrt = layr[0]
		b.Rbrtopt = layr
		b.Ast = layr[0][4]
		n1, n2, d1, d2 := layr[0][0], layr[0][1], layr[0][2], layr[0][3]
		b.Dia1 = d1; b.Dia2 = d2
		b.N1 = n1; b.N2 = n2
		b.Nlayers += int(layr[0][7])
		case 2:
		b.Rbrc = layr[0]
		b.Rbrcopt = layr
		b.Asc = layr[0][4]
		n1, n2, d1, d2 := layr[0][0], layr[0][1], layr[0][2], layr[0][3]
		b.Dia3 = d1; b.Dia4 = d2
		b.N3 = n1; b.N4 = n2
		b.Nlayers += int(layr[0][7])
	}
	return
}

//RBarGen generates rebar templates for a ribbed slab
func (b *RccBm) RBarGen()(err error){
	//ribbed slab bar gen
	if b.Ast > 0.0{
		rez := RbrSingle(b.Ast)
		err, _ = b.BarSort(rez, 1)
	}
	if b.Asc > 0.0{
		rez := RbrSingle(b.Asc)
		err, _ = b.BarSort(rez, 2)
	}
	return
	
}

//BarGen generates rebar options for a beam section
func (b *RccBm) BarGen() (err error, mincvr float64){
	var rez [][]float64
	if b.Ast == 0.0 && b.Asc == 0.0 {err = errors.New("no steel areas specified for beam")}
	if b.Asc == 0.0 && !b.Rslb{
		//add 2 12mm dia bars on top
		b.Asc = 225.0
	}
	maxbars := 10
	if b.D1 + b.D2 > 0.0{

		if b.Ast > 0.0{
			if b.Ast < 225.0{
				b.Ast = 225.0
			}
			rez, err = BmBarDias(maxbars,b.D1, b.D2, b.Ast)
			if err == nil{
				err, mincvr = b.BarSort(rez, 1)
				if err != nil{
					return
				}
			} else {
				return
			}
		}
		if b.Asc > 0.0{
			if b.Asc < 225.0{
				b.Asc = 225.0
			}
			rez, err = BmBarDias(maxbars,b.D1, b.D2, b.Asc)
			if err == nil{
				err, mincvr = b.BarSort(rez, 2)
				if err != nil{
					return
				}
			} else {
				return
			}
		}
		
	} else {
		if b.Ast > 0.0{
			if b.Ast < 225.0 && !b.Rslb{
				b.Ast = 225.0
			}
			rez, err = BmBarCombo(maxbars,b.Ast)
			if err == nil{
				err, mincvr = b.BarSort(rez, 1)
				if err != nil{
					return
				}
			} else {
				return
			}
		}
		if b.Asc > 0.0{
			if b.Asc < 225.0 && !b.Rslb{
				b.Asc = 225.0
			}
			rez, err = BmBarCombo(maxbars,b.Asc)
			if err == nil{
				err, mincvr = b.BarSort(rez, 2)
				if err != nil{
					return
				}
			} else {
				return
			}
		}
	}
	return
}

//BmBarDias generates rebar dia combos for a given ast, dia 1 and 2
func BmBarDias(maxbars int, d1, d2 float64, ast float64) (rez [][]float64, err error){
	nbars := 1
	r, err := BmBarSizes(maxbars, nbars, d1, 0.0, ast)
	if err == nil{rez = append(rez, r)}
	r, err = BmBarSizes(maxbars, nbars, d2, 0.0, ast)
	if err == nil{rez = append(rez, r)}
	r, err = BmBarSizes(maxbars, nbars, d1, d2, ast)
	if err == nil{rez = append(rez, r)}
	if len(rez) == 0 {err = errors.New(fmt.Sprintf("no rebar combos found for ast %.f",ast))}
	return
}

//BmBarCombo returns rebar dia combos (calls BmBarDias for d1, d2, d1+d2)
func BmBarCombo(maxbars int, ast float64) (rez [][]float64, err error){
	//log.Println(ColorPurple,"ast in->",ast,ColorReset)
	nbars := 2
	//bmd := []float64{10.0,12.0,14.0,16.0,18.0,20.0,22.0,25.0,28.0,32.0}
	bmd := []float64{12.0,16.0,20.0,25.0,28.0,32.0}
	
	for _, d1 := range bmd[:len(bmd)-1]{
		for _, d2 := range bmd[1:]{
			if d1 == d2{continue}
			r, err := BmBarSizes(maxbars, nbars, d1, d2, ast)
			if err != nil{continue}
			rez = append(rez, r)
		}
	}
	nbars = 1
	for _, d1 := range bmd{
		r, err := BmBarSizes(maxbars, nbars, d1, 0.0, ast)
		if err != nil{continue}
		rez = append(rez, r)
	}
	sort.SliceStable(rez, func(i, j int) bool{
		return rez[i][4] < rez[j][4]
	})
	if len(rez) == 0 {err = errors.New(fmt.Sprintf("no rebar combos found for ast %.f",ast))}
	return
}


//BmBarSizes returns the minimum area n1 d1 n2 d2 combo 
func BmBarSizes(maxbars, nbars int, d1, d2, ast float64)(rez []float64, err error){
	var astmin, adiff float64
	var n1, n2 int
	switch nbars{
		case 1:
		a1 := RbrArea(d1)
		n := math.Ceil(ast/a1)
		if n == 1.0{n = 2.0}
		astmin = n * a1
		n1 = int(n); n2 = 0
		adiff = astmin - ast
		nf := math.Floor(ast/a1)
		
		if math.Abs(nf*a1 - ast) <= 1.0 && nf != 1.0{
			n1 = int(nf); astmin = nf * a1; adiff = astmin - ast
		}
		//TOOK THIS OUT FOR OPT ROUTINES (might be so wrong)
		//if n1 > maxbars {err = errors.New("max number of bars exceeded")}
		case 2:
		a1 := RbrArea(d1); a2 :=  RbrArea(d2)
		n1m := math.Ceil(ast/a1); n2m := math.Ceil(ast/a2)
		//if n1m == 1 {n1m = 2.0}
		//if n2m == 1 {n2m = 2.0}
		for i := int(n1m); i >1; i--{
			for j := 2; j <= int(n2m); j++{
				astot := math.Ceil(float64(i) * a1 + float64(j) * a2)
				if astot >= ast || ast - astot < 1.0{
					//if i + j > maxbars{continue}
					switch {
					case astmin == 0.0:
						astmin = astot
						n1 = i; n2 = j
						adiff = astot - ast
					default:
						if astmin > astot{
							astmin = astot
							n1 = i; n2 = j
							adiff = astot - ast
						}
					}
				}
			}
		}
	}
	if n1 + n2 > maxbars {
		err = errors.New("max number of bars exceeded")
		return
	}
	if astmin == 0.0{
		err = errors.New(fmt.Sprintf("no rebar combos found for ast %.f",ast))
		return
	}
	rez = []float64{float64(n1), float64(n2), d1, d2, astmin, ast, adiff}
	return 
}

/*

type BmRbr struct {
	Nt, Nc  []int
	Dt, Dc  []float64
	Ast float64
	Asc float64
 	Nbt, Nbc int
	Astd float64
	Ascd float64
	Astr float64
	Ascr float64
}

type BmLink struct {}

func (r *BmRbr) Report() (rez string){
	return
}

func (r *BmRbr) Printz() (rez string){
	if r.Nc != nil{
		rez += fmt.Sprintf("top layer (C)- n1 %v bars %.0f mm dia n2 %v bars %.0f mm dia\n",r.Nc[0],r.Dc[0], r.Nc[1], r.Dc[1])
		rez += fmt.Sprintf("area of steel- asc - %.2f\n",r.Asc)
	}
	if r.Nt != nil{
		rez += fmt.Sprintf("bottom layer (T)- n1 %v bars %.0f mm dia n2 %v bars %.0f mm dia\n",r.Nt[0],r.Dt[0], r.Nt[1], r.Dt[1])
		rez += fmt.Sprintf("area of steel- ast - %.2f\n",r.Ast)
	}
	return
}

*/
