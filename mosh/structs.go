package barf

import (
	kass "barf/kass"
	"fmt"
	//"encoding/json"
	//"log"
	//"math"
)

//rck is WHAAT?
type rck interface {
	Cord2d() [][]float64
	Cord3d() [][]float64
}

//RccMod will never be used, delete
type RccMod struct {
	Base             *kass.Model
	Fck, Fy          float64
	Grid             [][]int
	Xvec, Yvec, Zvec []float64
	Dx, Dy, Dz       float64
}

//FltSlb is not used; SubFrm is 
type FltSlb struct{
	//fck slab, fck col
	Fcks,Fys   []float64
	Rib        bool
	Dx, Dy, Dz float64
	Lx, Ly, Lz float64
	Clvrs      [][]float64
	Clvrz      []int
	Colsec     []float64
	Coltype    int
}

//Printz printz nothing of value here
func (c *CBm) Printz(){
	fmt.Println("continuous beam")
	fmt.Println("grade of concrete, steel-",c.Fck, c.Fy)
	fmt.Println("nspans -",c.Nspans)
	fmt.Println("lspans-",c.Lspans)
	fmt.Println("nominal cover-",c.Nomcvr)
}


//TODO FUNCS
