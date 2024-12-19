package barf

import (
	"math"
	kass"barf/kass"
)

func unemax(b *Bat, inp []interface{}) (err error){
	var fit int
	for _, val := range b.Pos{
		fit += val
	}
	b.Fit = float64(fit)
	err = nil
	return
}

func knpsck(b *Bat, inp []interface{}) (err error){
	var v, w, c float64
	val, _ := inp[0].([]float64)
	wt, _ := inp[1].([]float64)
	for i, idx := range b.Pos{
		if idx == 1{
			v += val[i]
			w += wt[i]
		}
	}
	if w > 10.0{c = w - 10.0}
	b.Fit = v / (1.0 + 100.0 * c)
	b.Wt = w
	b.Con = c
	err = nil
	return 
}

func trsrakaobj(b *Bat, inp []interface{}) (err error){
	mod, _ := inp[0].(kass.Model)
	secs, _ := inp[1].([][]float64)
	pmax,_ := inp[2].(float64)
	dmax,_ := inp[3].(float64)
	dens,_ := inp[4].(float64)
	cp := make([][]float64, len(b.Pos))
	for i, idx := range b.Pos{
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
		b.Wt = 1e5
		b.Fit = 1.0/b.Wt
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
	b.Wt = wt*(1.0 + 10.0*C)
	b.Fit = 1.0/b.Wt
	b.Con = con
	err = nil
	return
}
