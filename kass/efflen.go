package barf

import (
	"fmt"
	"math"
)

//EffHt calcs a column's effective height using the bs 8110 formula (hulse section 5.7)
//given a frmtyp (3 for both x and y directions, TODO)
func (c *ColEnv) EffHt(frmtyp int){
	var l1, l2, acl, acu, acmin float64
	acl = c.B1
	if c.Lbss{acl = 10.0}
	acu = c.B2
	if c.Ubss{acu = 10.0}
	switch c.Ljbase{
		case true:
		switch c.Fixbase{
			case true:
			acl = 1.0
			case false:
			acl = 10.0
		}
	}
	acmin = acl; if acl > acu{acmin = acu}
	//fmt.Println(acmin, acl, acu)
	switch c.Braced{
		case true:
		l1 = c.L0 * (0.7 + 0.05 * (acl + acu))
		l2 = c.L0 * (0.85 + 0.05 * acmin)
		if l1 <= l2{c.Le = l1} else {c.Le = l2}
		if c.Le > c.L0{c.Le = c.L0}
		case false:
		l1 = c.L0 * (1.0 + 0.15 * (acl + acu))
		l2 = c.L0 * (2.0 + 3.0 * acmin)
		if l1 <= l2{c.Le = l1} else {c.Le = l2}
	}
	switch frmtyp{
		case 3:
		//calc eff ht in y
		
	}
	//fmt.Println(ColorRed,c.Id,"->eff ht->",c.Le,c.B1, c.B2,ColorReset)
	//fmt.Println("lower beams->",c.Lbd)
	//fmt.Println("upper beams->",c.Ubd)
	return
}

//EffLen returns the effective length factor for a column
//as seen in harrison program efflen
func EffLen(sway bool, iter int, l0, ga, gb float64) (le float64){
	//returns the effective length FACTOR
	//ga - cstiff/sum of beam stiff above, gb - below
	var a, f, df float64
	le = l0
	for i := 0; i < iter; i++{
		a = math.Pi/le
		ct := math.Cos(a)/math.Sin(a)
		ct2 := math.Cos(a/2.0)/math.Sin(a/2.0)
		switch sway{
			case false:
			//sidesway prevented
			f = ga * gb * a * a/4.0 + 0.5 * (ga + gb) * (1.0 - a * ct) + 2.0/a/ct2
			f -= 1.0
			df = ga * gb * a/2.0 + 0.5 * (ga + gb) * (a/(math.Sin(a)*math.Sin(a)) - ct)
			df = df + 1.0/(a * math.Cos(a/2.0) * math.Cos(a/2.0)) - 2.0/(a * a * ct2)
			case true:
			//sways
			f = ga * gb * a * a - 36.0 - 6.0 * a * (ga + gb) * ct
			df = 2.0 * ga * gb * a - 6.0 * (ga + gb) * ct
			df = df + 6.0 * a * (ga + gb)/math.Sin(a) * math.Sin(a)
		}
		a -= f/df
		le = math.Pi/a
		fmt.Println("i, eff len",i, le)
	}
	return
}
