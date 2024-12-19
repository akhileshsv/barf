package barf

//scratch functions for column optimization

import (
	"log"
	"time"
	"math/rand"
)

type rat struct{
	id  int
	vec []int
	pos []float64
	vel []float64
	pbest []float64
	//btopk int
	//rbr  cbr
	fitness float64
	cons int
	ast  float64
	typ string
	fbest float64
}

func maker(vecin []int, pin []float64, typ string)(r rat){
	r.vec = make([]int,len(vecin))
	r.pos = make([]float64,len(vecin))
	r.vel = make([]float64,len(vecin))
	r.pbest = make([]float64, len(vecin))
	r.fbest = -666.0
	for i := range r.vec{
		r.vel[i] = 2.0 * rand.Float64() - 1.0
		r.pos[i] = pin[i]
		r.vec[i] = vecin[i]
	}
	r.typ = typ
	return
}

func (r *rat) printz(){
	log.Println("rat # ", r.id)
	log.Println("vec", r.vec)
	log.Println("pos", r.pos)
	log.Println("vel", r.vel)
	log.Println("fitness", r.fitness)
	log.Println("cons", r.cons)
	log.Println("ast", r.ast)
}

func (r *rat) evalc(c *RccCol, pu, mux float64){
	r.cons = 0
	//log.Println("eval->", r.id)
	log.Println(r.vec)
	log.Println(r.pos)
	b := c.B - c.Cvrc; d := c.H - c.Cvrc
	var n2av, n2req, n3av, n3req, svc float64
	if cdias[r.vec[1]]/cdias[r.vec[0]] - 1.0 > 0{
		r.cons++
		svc += cdias[r.vec[1]]/cdias[r.vec[0]] - 1.0
	}
	if cdias[r.vec[2]]/cdias[r.vec[0]] - 1. > 0{
		r.cons++
		svc += cdias[r.vec[1]]/cdias[r.vec[0]]
	}
	if r.vec[3] > 0 {
		n2av = (b - 2.0*(cdias[r.vec[0]]+cldiv))
		n2req = float64(r.vec[3])*cdias[r.vec[1]] + float64(r.vec[3]-1)*cldiv
		if n2req/n2av-1.0 > 0 {
			r.cons++
			svc += n2req/n2av - 1.0
		}
	}
	if r.vec[4] > 0 {
		n3av = (d - 2.0*(cdias[r.vec[0]]+cldiv))
		n3req = float64(r.vec[4])*cdias[r.vec[2]] + float64(r.vec[4]-1)*cldiv
		if n3req/n3av-1.0 > 0 {
			r.cons++
			svc += n3req/n3av - 1.0
		}
		
	}
	ratbr := cbr{
		d1:r.vec[0],
		d2:r.vec[1],
		d3:r.vec[2],
		na:r.vec[3],
		nb:r.vec[4],
	}
	ratc := &RccCol{
		Fck:c.Fck,
		Fy:c.Fy,
		B:c.B,
		H:c.H,
		Dtyp:c.Dtyp,
		Cvrt:c.Cvrt,
		Rtyp:c.Rtyp,
		Styp:c.Styp,
	}
	ColRbrCbrDet(ratc, ratbr)
	pur, mur, err := ColAzIs(c, pu)
	if err != nil{
		r.cons++
		svc += 1.0
	}
	if pu/pur - 1.0 > 0.0{
		r.cons++
		svc += pu/pur - 1.0
	}
	if mux/mur - 1.0 > 0.0{
		r.cons++
		svc += mux/mur - 1.0
	}
	d1 := cdias[r.vec[0]]; d2 := cdias[r.vec[1]]; d3 := cdias[r.vec[2]]
	r.ast = 4.0 * RbrArea(d1) + 2.0 * RbrArea(d2) * float64(r.vec[3]) + 2.0 * RbrArea(d3) * float64(r.vec[4])
	r.fitness = 1.0/(r.ast*(1.0 + svc))
	//r.fitness = 1.0/((1.0 + 1000.0* float64(r.cons))*r.ast)
	if r.fbest > r.fitness || r.fbest == -666.0{
		r.fbest = r.fitness
	}
	//log.Println("fitness",r.fitness)
	//log.Println("end eval")
}

func normvec(vecin []int, typ string) (pin []float64){
	pin = make([]float64,len(vecin))
	switch typ{
		case "col-r":		
		for i, val := range vecin{
			switch {
			case i > 2:
				pin[i] = float64(val+1)/float64(len(cdias)) 
			default:
				pin[i] = float64(val+1)/4.0
			}	
		}
	}
	return
}

func (r *rat) updatepos(){
	switch r.typ{
		case "col-r":
		for i, val := range r.pos{
			r.pos[i] = val + r.vel[i]			
		}
		
	}
}

func (r *rat) updatevec(){
	//convert pos to vec
	switch r.typ{
		case "col-r":
		for i, val := range r.pos{
			switch {
			case i > 2:
				r.vec[i] = int(4.0 * val - 1.0)
				if r.vec[i] > 3{r.vec[i] = 3}
				if r.vec[i] < 0{r.vec[i] = 0}
			default:
				r.vec[i] = int(float64(len(cdias)) * val - 1.0)
				if r.vec[i] > len(cdias) -1 {r.vec[i] = len(cdias) -1}
				if r.vec[i] < 0{r.vec[i] = 0}
			}
		}
	}
}

func (r *rat) updatevel(kbest []float64) {
	w := 0.5
	c1 := 1.0
	c2 := 2.0
	for i, vi := range r.vel{
		r1 := rand.Float64(); r2 := rand.Float64()
		velcog := c1 * r1 * (r.pbest[i] - r.pos[i])
		velsoc := c2 * r2 * (kbest[i] - r.pos[i])
		r.vel[i] = w * vi + velcog + velsoc
	}
}
	
func ColOptPsoSimp(c *RccCol, pu, mux float64, vecin []int){
	rand.Seed(time.Now().UnixNano())
	kmax :=500
	fmin := -666.
	fprev := 1.0
	kbest := make([]float64, 5)
	vbest := make([]int, 5)
	swarm := make([]*rat, 50)
	var bid int
	typ := "col-r"
	pin := normvec(vecin, "col-r")
	brat := maker(vecin, pin, typ)
	log.Println("init position-> ", pin)
	for i := range swarm{
		r := maker(vecin, pin, typ)
		r.id = i
		swarm[i] = &r
	}
	for i := 0; i < kmax; i++{
		log.Println("gen->",i,"delta fitness",(fprev - fmin/fprev)*100.0)
		for j, r:= range swarm{
			r.evalc(c, pu, mux)
			if r.fitness < fmin || fmin == -666.{
				fmin = r.fitness
				copy(kbest, r.pos)
				copy(vbest, r.vec)
				bid = j
			}
		}
		copy(brat.vec, swarm[bid].vec)
		copy(brat.pos, swarm[bid].pos)
		brat.fitness = swarm[bid].fitness
		brat.cons = swarm[bid].cons
		brat.ast = swarm[bid].ast
		//brat.printz()
		for _, r := range swarm{
			r.updatevel(kbest)
			r.updatepos()
			r.updatevec()
		}
		fprev = fmin
	}
	//d1, d2, d3 := vbest[0], vbest[1], vbest[2]
	//ast := 4.0 * RbrArea(cdias[d1]) + 2.0 * RbrArea(cdias[d2]) * float64(cdias[3]) + 2.0 * RbrArea(cdias[d3]) * float64(cdias[4])
	//log.Println(ast)
	log.Println(vbest)
	log.Println(kbest)
	brat.printz()
	log.Println("percent steel:",brat.ast*100.0/(c.B * c.H))
}
