package barf

import (
	"fmt"
	"math"
	kass"barf/kass"
)

//Cbs is a column base plate struct
type Cbs struct{
	Pu    float64
	Fck   float64
	L     float64
	B     float64
	Oa    float64
	Ob    float64
	Ymo   float64
	Ymw   float64
	Fu    float64
	Fy    float64
	Ar    float64
	Ts    float64
	Sbc   float64
	Df    float64
	Bdia  float64 //bolt diameter
	Lblt  float64 //bolt length
	Pdz   bool    //design pedestal/footing
	Sdx   int
	Code  int
	Sname string
	Id    string
	Title string
	Verbose bool
}

// //ColBaseDz designs a column base plate, pedestal (based on df) and footing
// func ColBaseDz(c *Cbs)(err error){
	
// }

//SlbCbsDz designs an axially loaded slab base
func SlbCbsDz(c *Cbs)(err error){
	if c.Fu == 0.0{c.Fu = 410.0}
	if c.Fy == 0.0{c.Fy = 250.0}
	if c.Ymo == 0.0{c.Ymo = 1.1}
	if c.Ymw == 0.0{c.Ymw = 1.25}//shop welding
	ss, err := kass.GetStlSec(c.Sname, c.Sdx, c.Code)
	if err != nil{return}
	//ss.Printz()
	c.Ar = 1e3*c.Pu/(0.45 * c.Fck)
	qa := 1.0; qb := (ss.B+ss.H)/2.0; qc := (ss.B * ss.H -c.Ar)/4.0
	c.Oa = (-qb + math.Sqrt(qb * qb - 4.0*qa*qc))/(2.0*qa)

	c.Oa = math.Ceil(c.Oa/5.0) * 5.0
	//set min. val. why is it negative sometimes :-|
	if c.Oa < 50{
		c.Oa = 50.0
	}
	c.Ob = c.Oa
	c.L = math.Round(ss.H) + 2.0 * c.Oa
	c.B = math.Round(ss.B) + 2.0 * c.Oa
	
	wck := 1e3*c.Pu/(c.L * c.B)
	
	c.Ts = 2.5 * wck * (c.Oa * c.Oa - 0.3 * c.Ob * c.Ob) * c.Ymo/c.Fy
	c.Ts = math.Sqrt(c.Ts)
	if c.Ts < ss.Tf {c.Ts = ss.Tf}
	c.Ts = math.Ceil(c.Ts/2.0) * 2.0
	for _, ps := range CbsPlts{
		if c.Ts >= ps{
			c.Ts = ps
		}
	}
	if c.Ts > 30.0{
		err = fmt.Errorf("invalid thickness of base plate %f",c.Ts)
		return
	}
	if c.Verbose{
		fmt.Println("req area-",c.Ar,"mm2")
		fmt.Println("depth, width of section-",ss.B, ss.H)
		fmt.Printf("req offsets oa %.2f mm ob %.2f mm\n",c.Oa, c.Ob)
		fmt.Println("depth, width of section-",ss.B, ss.H)
		fmt.Println("depth, width of plate-",c.L, c.B)
		fmt.Printf("bearing pressure of concrete %.1f n/mm2 ok? %t\n",wck, wck < 0.45 * c.Fck)
		fmt.Printf("size of base %.1fx%.1f mm\n",c.L,c.B)
		fmt.Printf("thickness of plate %.0f mm\n",c.Ts)
	}
	return
}
