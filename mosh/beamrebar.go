package barf

import (
	//"log"
	"fmt"
	"math"
	"sort"
	"errors"
)

//rebar selection funcs for beams

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
func (b *RccBm) BarSort(rez [][]float64, tcdx int){
	layr := [][]float64{}
	for _, r := range rez{
		rlay := BarNRows(b.Bw, b.Dused, r)
		layr = append(layr, rlay)
	}
	sort.Slice(layr, func(i, j int) bool {
		if layr[i][7] != layr[j][7]{
			return layr[i][7] < layr[j][7]
		}
		return layr[i][4] < layr[j][4]
	})
	switch tcdx{
		case 1:
		b.Rbrt = layr[0]
		b.Rbrtopt = layr
		b.Ast = layr[0][4]
		n1, n2, d1, d2 := layr[0][0], layr[0][1], layr[0][2], layr[0][3]
		b.Dia1 = d1; b.Dia2 = d2
		b.N1 = n1; b.N2 = n2
		case 2:
		b.Rbrc = layr[0]
		b.Rbrcopt = layr
		b.Asc = layr[0][4]
		n1, n2, d1, d2 := layr[0][0], layr[0][1], layr[0][2], layr[0][3]
		b.Dia3 = d1; b.Dia4 = d2
		b.N3 = n1; b.N4 = n2
	}
	return
}

//RBarGen generates rebar templates for a ribbed slab
func (b *RccBm) RBarGen()(err error){
	//ribbed slab bar gen
	if b.Ast > 0.0{
		rez := RbrSingle(b.Ast)
		b.BarSort(rez, 1)
	}
	if b.Asc > 0.0{
		rez := RbrSingle(b.Asc)
		b.BarSort(rez, 2)
	}
	return
	
}

//BarGen generates rebar options for a beam section
func (b *RccBm) BarGen() (err error){
	var rez [][]float64
	if b.Ast == 0.0 && b.Asc == 0.0 {err = errors.New("no steel areas specified for beam")}
	if b.Asc == 0.0 && !b.Rslb{
		//add 2 12mm dia bars on top
		b.Asc = 225.0
	}
	if b.D1 + b.D2 > 0.0{

		if b.Ast > 0.0{
			if b.Ast < 225.0{
				b.Ast = 225.0
			}
			rez, err = BmBarDias(8,b.D1, b.D2, b.Ast)
			if err == nil{
				b.BarSort(rez, 1)
			} else {
				return
			}
		}
		if b.Asc > 0.0{
			if b.Asc < 225.0{
				b.Asc = 225.0
			}
			rez, err = BmBarDias(8,b.D1, b.D2, b.Asc)
			if err == nil{
				b.BarSort(rez, 2)
			} else {
				return
			}
		}
		
	} else {
		if b.Ast > 0.0{
			if b.Ast < 225.0 && !b.Rslb{
				b.Ast = 225.0
			}
			rez, err = BmBarCombo(8,b.Ast)
			if err == nil{
				b.BarSort(rez, 1)
			} else {
				return
			}
		}
		if b.Asc > 0.0{
			if b.Asc < 225.0 && !b.Rslb{
				b.Asc = 225.0
			}
			rez, err = BmBarCombo(8,b.Asc)
			if err == nil{
				b.BarSort(rez, 2)
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
		if n1 > maxbars {err = errors.New("max number of bars exceeded")}
		case 2:
		a1 := RbrArea(d1); a2 :=  RbrArea(d2)
		n1m := math.Ceil(ast/a1); n2m := math.Ceil(ast/a2)
		//if n1m == 1 {n1m = 2.0}
		//if n2m == 1 {n2m = 2.0}
		for i := int(n1m); i >1; i--{
			for j := 1; j <= int(n2m); j++{
				astot := math.Ceil(float64(i) * a1 + float64(j) * a2)
				if astot >= ast || ast - astot < 1.0{
					if i + j > maxbars{continue}
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
	if astmin == 0.0{err = errors.New(fmt.Sprintf("no rebar combos found for ast %.f",ast))}
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
