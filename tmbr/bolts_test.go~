package barf

import (
	"testing"
	//kass"barf/kass"
)

func TestBltDz(t *testing.T){
	t.Log("ramchandra ex. 13.16")
	ltyp := 1
	ptyp := 33
	lb := 100.0
	dia := 22.0
	fc := 7.8
	ang := 0.0
	pu := 46000.0
	bmem := 150.0
	kg := 0.0
	kpg := 0.0
	nrow := 2.0
	nblt, pblt, err := nbltsimp(ltyp, ptyp, pu, lb, dia, fc, ang)
	kg, kpg, err = bltfac(ltyp, ptyp, lb, dia)
	nsa := netsecchain(nrow, lb, dia, bmem)
	t.Log("nblt, pblt, err - ", nblt, pblt, err)
	t.Log("allow kg - ",kg * fc*lb*dia, " fpg - ", kpg * fc)
	t.Log("net sec area - ",nsa, "max loads->")
	t.Log("in plate - ", nsa * fc)
	t.Log("in bolting - par ",nblt * pblt, " perp ",nblt * pblt * kpg)
	t.Fail()
	t.Log("ramchandra ex. 13.17")
	ltyp = 1
	ptyp = 33
	lb = 150.0
	dia = 22.0
	fc = 11.2
	ang = 0.0
	pu = 150000.0
	nblt, pblt, err = nbltsimp(ltyp, ptyp, pu, lb, dia, fc, ang)
	t.Log("nblt, pblt, err - ", nblt, pblt, err)
	dend, dedge, pitch, gauge := bltdims(1, ltyp, 2, lb, dia)
	t.Log("blt dims - dend, dedge, pitch, gauge",dend, dedge, pitch, gauge)

}
