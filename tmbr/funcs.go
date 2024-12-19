package barf

import (
	"math"
	kass"barf/kass"
)

//WdSecCp calls kass.CalcSecProp to return section properties and calculates self-weight
func WdSecCp(styp, conv int, dims []float64, pg float64) (kass.Secprop, float64){
	if conv == 1{
		var dconv []float64
		for _, val := range dims{
			dconv = append(dconv, 25.4 * val)
		}
		copy(dims, dconv)
	}
	bar := kass.CalcSecProp(styp, dims)
	wdl := bar.Area * pg * 9800.0/1e9
	return bar, wdl
}


//WdFcAng uses hankinson's formula to calculate compressive strength (fca) at angle (ang)
func WdFcAng(fc, fcp, ang float64, conv bool) (fca float64) {
	if conv{ang = ang * math.Pi/180.0}
	fca = fc * fcp/(fc * math.Pow(math.Sin(ang),2) + fcp * math.Pow(math.Cos(ang),2))
	return
}

//getEqSddims returns equivalent square dimensions for circDims as per sp.33
func getEqSqdims() (sqDims [][]float64) {
	sqDims = make([][]float64, len(circDims))
	for i, dim := range circDims{
		sqDims[i] = append(sqDims[i], math.Sqrt(math.Pi) * dim[0]/2.0)
	}
	return
}

//getFc returns the compressive strength for a given slenderness ratio as per sp.33
func getFc(ldrat float64, grp int) (float64, error){
	fcvec := [][]float64{
		{10.6,6.3,5.6},
		{10.6,6.3,5.6},
		{10.6,6.3,5.6},
		{10.1,6.2,5.4},
		{9.0,5.9,5.1},
		{6.6,5.3,4.4},
		{4.6,4.2,2.8},
		{3.4,3.0,2.1},
		{2.6,2.3,1.6},
		{2.1,1.8,1.3},
		{1.7,1.5,1.0},
	}
	for idx, vec := range fcvec{
		v1 := vec[grp-1]
		ld1 := float64(idx * 5)
		if ld1 >= ldrat{
			if ld1 == ldrat || idx == 0{
				return v1, nil
			}
			if ld1 > ldrat{
				v0 := fcvec[idx-1][grp-1]
				fc := v0 + (ld1 - ldrat) * (v1 - v0)/5.0
				return fc, nil
			}
		}
	}
	return -99., ErrDim
}
