package barf

import (
	"fmt"
	"math/rand"
	"math"
	"time"
	draw"barf/draw"
)

type Sabobj func([]int, []interface{}) float64

func Bsamove(ndchk bool, nd []int, pos []int, mt int) (posn []int){
	posn = make([]int, len(pos))
	copy(posn, pos)
	switch mt{
		case 1:
		//swap
		
		i, j := rand.Intn(len(pos)), rand.Intn(len(pos))
		posn[i], posn[j] = posn[j], posn[i]
		case 2:
		//flip
		mutpt := rand.Intn(len(pos))
		mutdx := rand.Intn(nd[mutpt])
		posn[mutpt] = mutdx
		case 3:
		//shuffle
		idx := rand.Intn(len(pos)-1) 
		for i := len(pos) - 1; i > idx; i-- {
			j := rand.Intn(i + 1)
			posn[i], posn[j] = posn[j], posn[i]
		}
	}
	if ndchk{
		for i, v := range posn{
			if v > nd[i]{v = pos[i]}
		}
	}
	//fmt.Println("pos new-",posn, "pos olde-",pos)
	return
}

func Bsaloop(ndchk bool, temp, dmp, fa float64, posa, nd []int, ng, mt int, obj Sabobj, inp []interface{}, drw, trm, title string){
	rand.Seed(time.Now().UnixNano())
	var dat string
	var pf, gf float64
	pos, gpos := make([]int, len(nd)), make([]int, len(nd))
	if len(posa) == len(nd){
		fmt.Println("copying start point")
		copy(pos, posa)
		pf = fa
		gf = pf
	} else {
		for i := range pos{
			pos[i] = rand.Intn(nd[i])
		}
		pf = obj(pos, inp)
		gf = pf
	} 
	copy(gpos, pos)
	//var rez []float64
	for  gen := 0; gen < ng; gen++{
		switch drw{
			case "all":
			fmt.Println(ColorBlue, "gen->",gen,ColorWhite,"\ncurrent pos->\n",gpos, ColorGreen,"\nmin fitness->", gf,ColorReset)
			default:
			if gen % 10 == 0{
				fmt.Println(ColorBlue, "gen->",gen,ColorWhite, ColorGreen,"\tmin fitness->", gf,ColorReset)
			}
			
		}
		posn := Bsamove(ndchk, nd, pos, mt)
		nf := obj(posn, inp)
		if nf < gf {
			gf = nf
			copy(gpos, posn)
			//rez = append(rez, gf)
		}
		diff := nf - pf
		t := temp/(float64(gen)+1.)
		metro := math.Exp(-diff/t)
		if diff < 0.0 || rand.Float64() < metro {
			copy(pos, posn)
			pf = nf
		}
		dat += fmt.Sprintf("%v %f\n", gen, gf)
	}
	fmt.Println(ColorCyan, "min fitness->",gf,ColorRed,"position->",gpos, ColorReset)
	skript := "d2.gp"
	var dstr string
	switch drw{
		case "gen":
		dstr, _ = draw.Dumb(dat, skript, trm, "eval", "gen", "gbest", "")
		if trm != "qt" {fmt.Println(dstr)}
	}
}

