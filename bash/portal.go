package barf

import (
	"fmt"
	//"math"
	//opt "barf/opt"
	kass"barf/kass"
)

//getportal returns a portal frame struct from a floating point position vector
func getportal(pos []float64, inp []interface{}) (f kass.Portal){
	f, _ = inp[0].(kass.Portal)
	f.Kost = 0.0
	f.Term = ""
	f.Verbose = false
	switch f.Config{
		case 0:
		//uniform section
		switch f.Readsec{
			case true:
			//pos - sdx, sdx
			case false:
			//pos - dcol, dbm
		}
		case 1:
		//haunched rafters
		switch f.Readsec{
			case false:
			case true:
		}
		case 2:
		//haunched rafters and cols
		case 3:
		//uniformly tapered
		//pos = ds col, de col/ds bm, de bm/apex
		switch f.Readsec{
			case true:
			//pos - db, de, dapex
			case false:
			//pos - db, de, dapex(dt)
			db := pos[0]; de := pos[1]; dt := pos[2]
			f.Hdims = [][]float64{{db, de, dt}}
		}
	}
	return
}

func PortalOpt(pf kass.Portal)(err error){
	fmt.Println("starting yeopt of ",pf.Title)
	switch pf.Opt{
		case 1,11,12,13:
		portalga(pf)
		case 2,21,22,23:
		portalpso(pf)
	}
	err = pf.Calc()
	return
}

func portalga(pf kass.Portal){
	
}

func portalpso(pf kass.Portal){
	
}

