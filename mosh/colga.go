package barf

//scratch funcs for column opt

import (
	"log"
	//"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	cdias   = []float64{12, 16, 18, 20, 22, 25, 28, 32}
	careas  = []float64{113.1, 201.1, 254.5, 314.2, 380.1, 490.8, 615.7, 804.2}
	cldiv  = 75.0
	nomcvr = 40.0
)



type cbr struct {
	d1, d2, d3 int
	na, nb     int
	fitness    float64
	ast        float64
	delta      float64
	n1, n2, n3 int
	cons       int
}

func (r *cbr) printz() {
	log.Println("rbr deetz->")
	log.Println(r.n1, " n1 bars of ", cdias[r.d1], " mm dia")
	log.Println(r.n2, " n2 bars of ", cdias[r.d2], " mm dia")
	log.Println(r.n3, " n3 bars of ", cdias[r.d3], " mm dia")
	log.Println("cons->", r.cons)
	log.Println("ast total-> ", r.ast, " mm2 \ndelta->", r.delta, " mm2")
}

func (r *cbr) Fitr(astr, astx, asty, b, d float64) {
	r.cons = 0
	r.ast = 4.0*careas[r.d1] + 2.0*careas[r.d2]*float64(r.na) + 2.0*careas[r.d3]*float64(r.nb)
	if r.ast/astr-1.0 < 0{
		r.cons++
	}
	var n2av, n2req, n3av, n3req float64
	if r.na > 0{
		n2av = (b - 2.0*(cdias[r.d1]+cldiv))
		n2req = float64(r.na)*cdias[r.d2] + float64(r.na-1)*cldiv
		if n2av/n2req - 1.0 < 0{
			r.cons++
		}
	}
	if astx > 0{
		asx := 4.0*careas[r.d1] + 2.0*careas[r.d2]*float64(r.na)
		if asx/astx - 1.0 < 0{
			r.cons++
		}
	}
	if asty > 0{
		asy := 4.0*careas[r.d1] + 2.0*careas[r.d3]*float64(r.nb)
		if asy/asty - 1.0 < 0{
			r.cons++
		}
	}
	if r.nb > 0{
		n3av = (d - 2.0*(cdias[r.d1]+cldiv))
		n3req = float64(r.nb)*cdias[r.d3] + float64(r.nb-1)*cldiv
		if n3av/n3req - 1.0 > 0 {
			r.cons++
		}
	}
	r.fitness = 1.0 / ((r.ast - astr) * (1000.0*float64(r.cons) + 1))
	r.n1 = 4
	r.n2 = 2 * r.na
	r.n3 = 2 * r.nb
	r.delta = r.ast - astr
}

func Maker(astr, astx, asty, b, d float64) (r cbr) {
	//SHAI-HULUD is the bringer of the water of life
	r.d1 = rand.Intn(len(cdias))
	r.d2 = rand.Intn(r.d1 + 1)
	r.d3 = rand.Intn(r.d1 + 1)
	r.na = rand.Intn(4)
	r.nb = rand.Intn(4)
	r.Fitr(astr, astx, asty, b, d)
	return
}
	
func Genpop(popsize int, astr, astx, asty, b, d float64) (popr []*cbr) {
	for i := 0; i < popsize; i++{
		r := Maker(astr, astx, asty, b, d)
		popr = append(popr, &r)
	}
	return
}

func Epop(popr []*cbr, popsize int) (fpop float64, midx int, matr []*cbr) {
	for x := range popr{
		if popr[x].fitness > fpop{
			fpop = popr[x].fitness
			midx = x
		}
		y := rand.Intn(popsize)
		if popr[y].fitness > popr[x].fitness{
			matr = append(matr, popr[y])
		} else {
			matr = append(matr, popr[x])
		}
	}
	return
}

func Repop(matr []*cbr, popsize int, astr, astx, asty, b, d, pcrs, pmut float64) (npop []*cbr) {
	for len(npop) < popsize{
		var mutz float64
		x, y := rand.Intn(len(matr)), rand.Intn(len(matr))
		if x == y {
			continue
		}
		x1, y1 := Crsr(matr[x], matr[y], pcrs)
		mutz = rand.Float64()
		if mutz < pmut {
			x1.Mutr()
		}
		x1.Fitr(astr, astx, asty, b, d)
		mutz = rand.Float64()
		if mutz < pmut {
			y1.Mutr()
		}
		y1.Fitr(astr, astx, asty, b, d)
		npop = append(npop, &x1)
		npop = append(npop, &y1)
	}
	return
}

func Crsr(x, y *cbr, pcrs float64) (a, b cbr) {
	if rand.Float64() > pcrs{
		//clone
		a.d1, a.d2, a.d3 = x.d1, x.d2, x.d3
		a.na, a.nb = x.na, x.nb
		b.d1, b.d2, b.d3 = y.d1, y.d2, y.d3
		b.na, b.nb = y.na, y.nb
	} else {
		//switch n2 and n3 arrangement
		a.d1, a.d2, a.d3 = x.d1, x.d2, x.d3
		a.na, a.nb = y.na, y.nb
		b.d1, b.d2, b.d3 = y.d1, y.d2, y.d3
		b.na, b.nb = x.na, x.nb
	}
	return
}

func (r *cbr) Mutr() {
	r.d1 = rand.Intn(len(cdias))
	r.d2 = rand.Intn(r.d1+1)
	r.d3 = rand.Intn(r.d1+1)
	r.na = rand.Intn(3)
	r.nb = rand.Intn(3)
}

func ColRbrOpt(astr, astx, asty, b, d float64) (zerocons []cbr, killr cbr){
	rand.Seed(time.Now().UnixNano())
	popsize := 100
	ngens := 7
	pcrs := 0.75
	pmut := 0.01
	var popr, matr []*cbr
	var fpop, fprev, fmax, deltamin float64
	var midx int
	deltamin = 1000.0
	popr = Genpop(popsize, astr, astx, asty, b, d)
	for i := 0; i < ngens; i++{
		fpop, midx, matr = Epop(popr, popsize)
		//log.Println("gen-> ", i, "max fitness-> ", fpop)
		if i > 0{
			//log.Printf("fitness improvement -> %.6f", 100.0*(fpop-fprev)/fprev)
			if fprev < fpop{
				pmut = pmut * 2
				pcrs = pcrs * 2
			} else {
				pmut = pmut/2
				pcrs = pcrs/2
			}
		}
		if fpop > fmax{
			fmax = fpop
		}
		if popr[midx].cons == 0{
			m := popr[midx]
			mr := cbr{d1: m.d1, d2: m.d2, d3: m.d3, na: m.na, nb: m.nb, n1: m.n1, n2: m.n2, n3: m.n3, delta: m.delta, ast: m.ast}
			zerocons = append(zerocons, mr)
		}
		//popr[midx].printz()
		if popr[midx].delta < deltamin{
			m := popr[midx]
			deltamin = m.delta
			killr = cbr{d1: m.d1, d2: m.d2, d3: m.d3, na: m.na, nb: m.nb, n1: m.n1, n2: m.n2, n3: m.n3, delta: m.delta, ast: m.ast}
		}
		fprev = fpop
		popr = Repop(matr, popsize, astr, astx, asty, b, d, pcrs, pmut)
	}
	return 
}


func ColRbrCbrDet(c *RccCol, rbr cbr){
	switch c.Styp{
		case 1:
		for i:=0; i< 2; i++{
			c.Dias = append(c.Dias, cdias[rbr.d1])
			c.Dbars = append(c.Dbars, c.Cvrc)
			c.Dias = append(c.Dias, cdias[rbr.d1])
			c.Dbars = append(c.Dbars, c.H - c.Cvrt)
		}
		if rbr.na > 0{
			for i:=0; i < rbr.na; i++{
				c.Dias = append(c.Dias, cdias[rbr.d2])
				c.Dbars = append(c.Dbars, c.Cvrc)
				c.Dias = append(c.Dias, cdias[rbr.d2])
				c.Dbars = append(c.Dbars, c.H - c.Cvrt)
			}
		}
		if rbr.nb > 0{
			step := math.Round((c.H - c.Cvrc - c.Cvrt - cdias[rbr.d1])/float64(rbr.nb))
			start := c.Cvrc + cdias[rbr.d1]/2.0 - step
			for i := 0; i < rbr.nb; i++{
				start += step
				c.Dias = append(c.Dias, cdias[rbr.d3])
				c.Dbars = append(c.Dias, start)
			}
		}
	}
}
