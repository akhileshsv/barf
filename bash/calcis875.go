package barf

import (
	//"fmt"
	"math"
	"log"
	
)

//DELETE THIS FILE ITS IN KASS
//wltable5 gets wind loads on walls as per is875-3
func wltable5(vz, h, w, slope, cpi float64) (pd float64, cpos, cneg []float64){
	baseangles := []float64{0.,5.,10.,20.,30.,45.,60.}
	angle := math.Atan(slope) * 180./math.Pi
	cpw0, cpw90, cpl0, cpl90 := 0., 0., 0., 0.
	var ww0, ww90, lw0, lw90 []float64
	idx := -1
	for i, ang := range baseangles{
		if angle < ang{
			idx = i-1
			break
		}	
	}
	if idx == -1{
		log.Println("ERRORE,errore->roof angle calculation error")
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
	multiplier := (angle-baseangles[idx])/(baseangles[idx+1]-baseangles[idx])
	cpw0 = ww0[idx] + multiplier*(ww0[idx+1]-ww0[idx])
	cpl0 = lw0[idx] + multiplier*(lw0[idx+1]-lw0[idx])
	cpw90 = ww90[idx] + multiplier*(ww90[idx+1]-ww90[idx])
	cpl90 = lw90[idx] + multiplier*(lw90[idx+1]-lw90[idx])
	//vz := 1.*1.038*1*v
	pd = 0.6*vz*vz
	if h <= 10{pd = 0.75 * pd}
	cpes := []float64{cpw0,cpl0,cpw90,cpl90}
	cpos, cneg = make([]float64, 4), make([]float64, 4)
	for i, cp := range cpes{
		cpos[i] = (cp + cpi)
		cneg[i] = (cp - cpi)
	}
	return
}
