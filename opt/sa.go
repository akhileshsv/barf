package barf

import (
	"fmt"
	"math/rand"
	"math"
	"time"
	draw"barf/draw"
)

func samove(nd int, pos, mx, mn []float64, dmp float64) (posn []float64){
	posn = make([]float64, nd)
	for i := range posn{
		posn[i] = pos[i] + dmp * (rand.Float64() - 0.5)
		if posn[i] > mx[i]{posn[i] = mx[i]}
		if posn[i] < mn[i]{posn[i] = mn[i]}
	}
	return
}

func Saloop(obj func([]float64) float64, temp, dmp float64, mx, mn []float64, nd, ng int, drw, trm, title string){
	rand.Seed(time.Now().UnixNano())
	var dat string
	pos, gpos := make([]float64, nd), make([]float64, nd)
	for i := range pos{
		pos[i] = (mx[i] - mn[i]) * rand.Float64() + mn[i]
	}
	pf := obj(pos)
	gf := pf
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
		posn := samove(nd, pos, mx, mn, dmp)
		nf := obj(posn)
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
		dstr,_ = draw.Dumb(dat, skript, trm, "eval", "gen", "gbest", "")
		if trm != "qt" {fmt.Println(dstr)}
	}
}
