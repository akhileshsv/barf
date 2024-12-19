package barf

import (
	"fmt"
	"math"
	"testing"
)

func TestWlRoof(t *testing.T){
	/*
	   func wltable5(pd, h, w, l, cpi float64)(cpos, cneg []float64)
	
	func wltable6(vz, h, w, slope, cpi float64) (pd float64, cpos, cneg []float64, wlcs map[int][]float64)   
	*/
	var rezstr string
	rezstr += "testing is code wind load cp generation\n"
	var h, w, l, slope, vz, cpi, theta float64
	rezstr += "testing example 1, sp.64 - 5° a type roof\n"
	theta = 5.0 * math.Pi/180.0
	h = 3500; w = 10000; vz = 44.42; slope = math.Tan(theta); cpi = 0.5
	pd, _, _, wlcs := wltable6(vz, h, w, slope, cpi)
	rezstr += fmt.Sprintf("design pressure pd %.2f n/m2\n",pd) 
	rezstr += "printing cpos and wload cases\n"
	for i := range wlcs[1]{
		rezstr += fmt.Sprintf("cpe left %.2f - cpe right %.2f\n",wlcs[1][i]/pd,wlcs[2][i]/pd)
		rezstr += fmt.Sprintf("udl left %.2f N/m2 udl right %.2f N/m2\n",wlcs[1][i],wlcs[2][i])
		fxl := wlcs[1][i] * math.Sin(theta); fyl := -wlcs[1][i] * math.Cos(theta)
		fxr := wlcs[2][i] * math.Sin(theta); fyr := -wlcs[2][i] * math.Cos(theta)
		if fxr > 0.0{fxr = -fxr}
		rezstr += fmt.Sprintf("nodal load left X %.2f N Y %.2f N right X %.2f N Y %.2f N\n",fxl, fyl, fxr, fyr)
	}
	rezstr += "testing example 2, sp.64 - 20° monoslope roof\n"
	
	theta = 20.0 * math.Pi/180.0
	h = 4000; w = 8000; vz = 35.88; slope = math.Tan(theta); cpi = 0.0
	pd, _, _, wlcs = wltable7(vz, h, w, l, slope, cpi)
	rezstr += fmt.Sprintf("design pressure pd %.2f n/m2\n",pd) 
	rezstr += "printing cpos and wload cases\n"
	for i := range wlcs[1]{
		rezstr += fmt.Sprintf("cpe left (H) %.2f - cpe right (L) %.2f\n",wlcs[1][i]/pd,wlcs[2][i]/pd)
		rezstr += fmt.Sprintf("udl left (H) %.2f N/m2 udl right (L) %.2f N/m2\n",wlcs[1][i],wlcs[2][i])
		fxl := wlcs[1][i] * math.Sin(theta); fyl := -wlcs[1][i] * math.Cos(theta)
		fxr := wlcs[2][i] * math.Sin(theta); fyr := -wlcs[2][i] * math.Cos(theta)
		if fxr > 0.0{fxr = -fxr}
		rezstr += fmt.Sprintf("nodal load left (H) X %.2f N Y %.2f N right (L) X %.2f N Y %.2f N\n",fxl, fyl, fxr, fyr)
	}
	t.Log(rezstr)
	t.Fatal("wind load generation test failed")
	
}
