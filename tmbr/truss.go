package barf

import (
	"fmt"
	kass"barf/kass"
)

//TrussDz designs a truss gen struct assuming all members are pin connected columns
//TODO - design members as beams
func TrussDz(t *kass.Trs2d) (err error){
	err = t.Calc()
	if err != nil{
		fmt.Println(err)
		return
	}
	colchn := make(chan []interface{}, len(t.Mod.Ms))
	for i := range t.Mod.Ms{
		cp := t.Mod.Mprp[i-1][3]
		dims := t.Sections[cp-1]
		styp := t.Styps[cp-1]
		go ColD(t.Mod.Ms[i], styp, t.Group , dims, t.Dzval, colchn)
	}
	for _ = range t.Mod.Ms{
		rez := <- colchn
		fmt.Println(rez[0])
	}
	return
}

//ReadPrp reads design values for a column struct
//use when non standard group properties need to be read in
func (c *WdCol) ReadPrp(dzval []float64){
	//read dzvals
	c.Prp.Em = dzval[0]
	c.Prp.Pg = dzval[1]
	
	//etc
}

//ColD designs an axially loaded wooden column
//using results from a kass.Model (only a 2d truss so far)
func ColD(mem *kass.Mem, styp, grp int, dims, dzval []float64, colchn chan []interface{}){
	//first build column
	//then check/design
	c := WdCol{
		Id: mem.Id,
		Grp: grp,
		Dims:dims,
		Styp:styp,
		Lspan:mem.Geoms[0],
	}
	if c.Grp == 0 && len(dzval) != 0{
		c.ReadPrp(dzval)
	}
	c.Init()
	c.Pu = mem.Cmax
	ok, val := ColChk(&c)
	rez := make([]interface{},3)
	rez[0] = ok
	rez[1] = val
	rez[2] = c
	colchn <- rez
}

