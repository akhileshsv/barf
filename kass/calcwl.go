package barf

//is code wind load generation funcs for sloped roof (truss and portal) models
//highly incomplete and possibly does not werk
import (
	"fmt"
	"math"
	"log"	
)

//wltable5 is supposed to generate wind load cases for building walls as per table 5 of is.875-3
func wltable5(pd, h, w, l, cpi float64){
	switch {
	case h/w < 0.5:
		switch{
			case l/w > 1.0 && l/w < 1.5:
			case l/w < 4.0:
		}
	case h/w < 1.5:
		switch{
			case l/w > 1.0 && l/w < 1.5:
			case l/w < 4.0:
		}
	case h/w < 6.0:
		switch{
			case l/w > 1.0 && l/w < 1.5:
			case l/w < 4.0:
		}
	default:
		//interpolate
	}
}

//wltable6 gets wind load cases for a type roofs as per table 6 of is 875 part 3
func wltable6(vz, h, w, slope, cpi float64) (pd float64, cpos, cneg []float64, wlcs map[int][]float64){
	/*
	var cpi float64
	switch inperm{
		case 1:
		cpi = 0.2
		case 2:
		cpi = 0.5
		case 3:
		cpi = 0.7
	}*/
	baseangles := []float64{0.,5.,10.,20.,30.,45.,60.}
	angle := math.Atan(slope) * 180./math.Pi
	fmt.Println("roof angle->", angle)
	cpw0, cpw90, cpl0, cpl90 := 0., 0., 0., 0.
	var ww0, ww90, lw0, lw90 []float64
	idx := -1
	var eqang bool
	for i, ang := range baseangles{
		if angle == ang{
			idx = i
			eqang = true
			break
		}
		if angle < ang{
			idx = i-1
			break
		}	
	}
	if idx == -1{
		log.Println("ERRORE,errore->roof angle calculation error")
		pd = -99.9
		return
	}
        switch{
		case h/w <= 0.5:
		ww0 = []float64{-0.8,-0.9,-1.2,-0.4,0,0.3,0.7}
		lw0 = []float64{-0.4,-0.4,-0.4,-0.4,-0.4,-0.5,-0.6}
		ww90 = []float64{-0.8,-0.8,-0.8,-0.7,-0.7,-0.7,-0.7}
		lw90 = []float64{-0.4,-0.4,-0.6,-0.6,-0.6,-0.6,-0.6,}
		case h/w <= 1.5:
		ww0 = []float64{-0.8,-0.9,-1.1,-0.7,-0.2,0.2,0.6}
		lw0 = []float64{-0.6,-0.6,-0.6,-0.5,-0.5,-0.5,-0.5}
		ww90 = []float64{-1.0,-0.9,-0.8,-0.8,-0.8,-0.8,-0.8}
		lw90 = []float64{-0.6,-0.6,-0.6,-0.6,-0.8,-0.8,-0.8}
		case h/w < 6:
		ww0 = []float64{-0.7,-0.7,-0.7,-0.8,-1.0,-0.2,0.2,0.5}
		lw0 = []float64{-0.6,-0.6,-0.6,-0.6,-0.5,-0.5,-0.5,-0.5}
		ww90 = []float64{-0.9,-0.8,-0.8,-0.8,-0.8,-0.8,-0.8,-0.8}
		lw90 = []float64{-0.7,-0.8,-0.8,-0.8,-0.7,-0.7,-0.7,-0.7}
	}
	pd = 0.6*vz*vz
	if eqang{
		cpw0 = ww0[idx] 
		cpl0 = lw0[idx] 
		cpw90 = ww90[idx] 
		cpl90 = lw90[idx] 
	} else {
		multiplier := (angle-baseangles[idx])/(baseangles[idx+1]-baseangles[idx])
		cpw0 = ww0[idx] + multiplier*(ww0[idx+1]-ww0[idx])
		cpl0 = lw0[idx] + multiplier*(lw0[idx+1]-lw0[idx])
		cpw90 = ww90[idx] + multiplier*(ww90[idx+1]-ww90[idx])
		cpl90 = lw90[idx] + multiplier*(lw90[idx+1]-lw90[idx])
	}	//vz := 1.*1.038*1*v
	cpes := []float64{cpw0,cpl0,cpw90,cpl90}
	cpos, cneg = make([]float64, 4), make([]float64, 4)
	for i, cp := range cpes{
		cpos[i] = (cp + cpi)
		cneg[i] = (cp - cpi)
	}
	wlcs = make(map[int][]float64)
	wlcs[1] = []float64{cpos[0]*pd, cpos[1]*pd, cneg[0]*pd, cneg[1]*pd, cpos[2] * pd, cpos[3]*pd, cneg[2]*pd, cneg[3]*pd}
	wlcs[2] = []float64{cpos[1]*pd, cpos[0]*pd, cneg[1]*pd, cneg[0]*pd, cpos[2] * pd, cpos[3]*pd, cneg[2]*pd, cneg[3]*pd}
	return
}

//wltable7 generates windload cases for lean to type roofs as per table 7 is-875, part 3
func wltable7(vz, h, w, slope, cpi float64) (pd float64, cpos, cneg []float64, wlcs map[int][]float64){
	pd = 0.6*vz*vz
	baseangles := []float64{0.,5.,10.,15.,20.,25.,30.}
	angle := math.Atan(slope) * 180./math.Pi
	fmt.Println(angle)
	var hl [][]float64
	var eqang bool
	idx := -1
	for i, ang := range baseangles{
		if ang >= angle{
			if ang == angle{
				idx = i
				eqang = true
			} else {
				idx = i-1
			}
			break
		}
	}
	
	if idx == -1{
		log.Println("ERRORE,errore->roof angle calculation error")
		return
	}
        switch{
		case h/w <= 2.0:
		hl = [][]float64{
			{-1.0,-0.5,-1.0,-0.9,-1.0,-0.5,-0.9,-1.0,-0.5,-1.0},
			{-1.0,-0.5,-1.0,-0.9,-1.0,-0.5,-0.9,-1.0,-0.5,-1.0},
			{-1.0,-0.5,-1.0,-0.8,-1.0,-0.5,-0.8,-1.0,-0.4,-1.0},
			{-0.9,-0.5,-1.0,-0.7,-1.0,-0.5,-0.6,-1.0,-0.3,-1.0},
			{-0.8,-0.5,-1.0,-0.6,-0.9,-0.5,-0.5,-1.0,-0.2,-1.0},
			{-0.7,-0.5,-1.0,-0.6,-0.8,-0.5,-0.3,-0.9,-0.1,-0.9},
			{-0.5,-0.5,-1.0,-0.6,-0.6,-0.5,-0.1,-0.6,-0.0,-0.6},
		}
	}
	multiplier := (angle-baseangles[idx])/(baseangles[idx+1]-baseangles[idx])
	var cpes []float64
	if eqang{
		for i, cp := range hl[idx]{
			if i == 4 || i == 5{
				cpes = append(cpes, cp)
				cpes = append(cpes, cp)
			} else {
				cpes = append(cpes, cp)
			}
		}
	} else {
		for i, cp := range hl[idx]{
			cp1 := hl[idx+1][i]
			cpn := cp + multiplier * (cp1 - cp)
			if i == 4 || i == 5{
				cpes = append(cpes, cpn)
				cpes = append(cpes, cpn)
			} else {
				cpes = append(cpes, cpn)
			}
		}
	}
	cpos = make([]float64, len(cpes))
	cneg = make([]float64, len(cpes))
	for i, cp := range cpes{
		cpos[i] = (cp + cpi)
		cneg[i] = (cp - cpi)
	}
	return
}

