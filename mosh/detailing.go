package barf

import (
	//"log"
	//"fmt"
	//"math"
)

//BarDevLen returns the development length for a given fck, fy and dia of rebar
func BarDevLen(fck, fy, dia float64) (ldreq float64){
	var tbd float64 //design bond stress
	tbdvals := map[float64]float64{15.0:1.0,20.0:1.2,25.0:1.4,30.0:1.5,35.0:1.7,40.0:1.9}
	if val, ok := tbdvals[fck]; !ok {
		if fck > 40.0 {
			tbd = 1.9
		}
		if fck < 15.0 {
			tbd = 0.6
		}
	} else {
		tbd = val
	}
	if fy > 250.0 {
		tbd = 1.6*tbd
	}
	ldfac := fy/(4.0*tbd*1.15)
	ldreq = ldfac * dia
	//log.Println("from bardev len-> fck, fy, dia->",fck, fy, dia, "rez->, ldfaq, ldreq->",ldfac, ldreq)
	return
}
