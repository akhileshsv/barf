package barf

//THIS SHOULD BE WORTHLESS, CHECK AND DELETE

import (
	"fmt"
	"time"
	//"log"
	//"math"
	"math/rand"
)

type prat struct{
	id, dims, cons int
	pos, vel, vec []float64
	pbest, gbest, fitness float64
	pmin, pmax, bestpos []float64
	typ string
}

func (p *prat) printz(){
	fmt.Println("prat #",p.id)
	fmt.Println("pos",p.pos)
	//fmt.Println("vec",p.vec)
	fmt.Println("vel",p.vel)
	fmt.Println("pbest, gbest",p.pbest, p.gbest)
	fmt.Println("best pos",p.bestpos)
	fmt.Println("cons",p.cons)
	switch p.typ{
		case "ftng":
		fmt.Printf("lx %.2f ly %.2f d %.2f psx %.2f psy %.2f",p.pos[0],p.pos[1],p.pos[2],p.pos[3],p.pos[4])
	}
}

func makep(dims int, pmin, pmax []float64, typ string)(p prat){
	p.vec = make([]float64,dims)
	p.pos = make([]float64,dims)
	p.vel = make([]float64,dims)
	p.dims = dims
	p.pmin = pmin
	p.pmax = pmax
	//init with random velocities and zero position
	for i := range p.vec{
		vmax := 0.5 * (pmax[i] - pmin[i])
		p.vel[i] = 2.0 * vmax * rand.Float64() - vmax
		p.pos[i] = 0.0 + p.vel[i]
	}
	return
}

func (p *prat) evaluate(vecs [][]float64){
	switch p.typ{
		case "ftng":
		//p.evalftng(vecs)
	}
}
/*
func (p *prat) evalftng(vecs [][]float64){
	colx, coly, fck, fy, df, sbc, pgck, pgsoil, nomcvr := vecs[0][0], vecs[0][1], vecs[0][2], vecs[0][3], vecs[0][4], vecs[0][5], vecs[0][6], vecs[0][7], vecs[0][8]
	pus := vecs[1]; mxs := vecs[2]; mys := vecs[3]; psfs := vecs[4]
	FtngBxEval(colx, coly, fck, fy, df, sbc, pgck, pgsoil, nomcvr, pus, mxs, mys, psfs, p)
}
*/

func randswarm(dims, nps, ngens int, pmin, pmax []float64, vecs [][]float64, typ string){
	rand.Seed(time.Now().UnixNano())
	var s []prat
	for i := 0; i < nps; i++{
		p := makep(dims, pmin, pmax, typ)
		s = append(s, p)
	}
	fmt.Println("swarm created")
	for i := 0; i < ngens; i++{
		for _, p := range s{
			p.evaluate(vecs)
			p.printz()
		}
	}
}
