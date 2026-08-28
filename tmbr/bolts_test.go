package barf

import (
	"testing"
	//kass"barf/kass"
)

func TestNailDz(t *testing.T){
	
}

func TestBltDz(t *testing.T){
	t.Log(ColorCyan,"ramchandra ex. 13.16", ColorReset)
	ltyp := 1
	ptyp := 33
	lb := 100.0
	dia := 22.0
	fc := 7.8
	fcp := 2.65
	ang := 0.0
	pu := 46000.0
	bmem := 150.0
	kg := 0.0
	kpg := 0.0
	nrow := 2.0
	nblt, pblt, err := nbltsimp(ltyp, ptyp, pu, lb, dia, fc, fcp, ang)
	kg, kpg, err = bltfac(ltyp, ptyp, lb, dia)
	nsa := netsecchain(nrow, lb, dia, bmem)
	t.Log("nblt, pblt, err - ", nblt, pblt, err)
	t.Log("allow kg - ",kg * fc*lb*dia, " fpg - ", kpg * fc)
	t.Log("net sec area - ",nsa, "max loads->")
	t.Log("in plate - ", nsa * fc)
	t.Log("in bolting - par ",nblt * pblt, " perp ",nblt * pblt * kpg)
	//t.Fail()
	t.Log(ColorCyan,"ramchandra ex. 13.17",ColorReset)
	ltyp = 1
	ptyp = 33
	lb = 150.0
	dia = 22.0
	fc = 11.2
	ang = 0.0
	pu = 150000.0
	nblt, pblt, err = nbltsimp(ltyp, ptyp, pu, lb, dia, fc, fcp, ang)
	t.Log("nblt, pblt, err - ", nblt, pblt, err)
	dend, dedge, pitch, gauge := bltdims(1, ltyp, 2, lb, dia)
	t.Log("blt dims - dend, dedge, pitch, gauge",dend, dedge, pitch, gauge)

	t.Log(ColorCyan,"ramchandra ex 13.18 - right angled bolted joint",ColorReset)


	ltyp = 2
	ptyp = 33
	lb = 100.0
	dia = 19.0
	fc = 7.8
	fcp = 2.65
	ang = 0.0
	pu = 25000.0
	nblt, pblt, err = nbltsimp(ltyp, ptyp, pu, lb, dia, fc, fcp, ang)
	t.Log("nblt, pblt, err - ", nblt, pblt, err)
	dend, dedge, pitch, gauge = bltdims(1, ltyp, 2, lb, dia)
	t.Log("blt dims - dend, dedge, pitch, gauge",dend, dedge, pitch, gauge)

	
	t.Log(ColorCyan,"ramchandra ex 13.19 - 45 degree joint",ColorReset)


	ltyp = 3
	ptyp = 33
	lb = 100.0
	dia = 19.0
	fc = 7.8
	fcp = 2.65
	ang = 45.0
	pu = 66000.0
	nblt, pblt, err = nbltsimp(ltyp, ptyp, pu, lb, dia, fc, fcp, ang)
	t.Log("nblt, pblt, err - ", nblt, pblt, err)
	dend, dedge, pitch, gauge = bltdims(1, ltyp, 2, lb, dia)
	t.Log("blt dims - dend, dedge, pitch, gauge",dend, dedge, pitch, gauge)
	
}
