package barf

import (
	"math"
	kass"barf/kass"
)


func sabtrsobj(pos []int, inp []interface{}) (fit float64){
	mod, _ := inp[0].(kass.Model)
	secs, _ := inp[1].([][]float64)
	pmax,_ := inp[2].(float64)
	dmax,_ := inp[3].(float64)
	dens,_ := inp[4].(float64)
	cp := make([][]float64, len(pos))
	for i, idx := range pos{
		cp[i] = make([]float64,1)
		if len(secs) == 1{
			cp[i][0] = secs[0][idx]
		} else {
			//WRENG
			cp[i][0] = secs[i][idx]
		}
	}
	mod.Cp = cp
	var wt, C, gx, con float64
	frmrez, err := kass.CalcTrs(&mod, mod.Ncjt)
	if err != nil {
		fit = 1e20
		return
	}
	js, _ := frmrez[0].(map[int]*kass.Node)
	ms, _ := frmrez[1].(map[int]*kass.Mem)
	for _, node := range js{
		for _, d := range node.Displ {
			gx = math.Abs(d)/dmax - 1.0
			if gx > 0.0 {
				C += gx
				con += 1.0
			}
		}
	}
	for _, mem := range ms{
		wt += mem.Geoms[0] * mem.Geoms[2] * dens
		pmem := mem.Qf[0] / mem.Geoms[2]
		gx = math.Abs(pmem)/pmax - 1.0
		if gx > 0.0 {
			C += gx
			con += 1.0
		}
	}
	fit = wt*(1.0 + 100.0*C)
	return
}
