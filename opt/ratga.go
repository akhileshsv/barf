package barf

import (
	"math"
	"math/rand"
	"sort"
)

type Rat struct{
	Pos []float64
	Fit float64
	Wt  float64
	Con float64
}

type robj func(*Rat, []interface{}) error

func (r *Rat) init(mx, mn, stp []float64, inp []interface{}, f robj) (error){
	r.Pos = make([]float64, len(mx))
	for i := range r.Pos{
		r.Pos[i] = (mx[i] - mn[i]) * rand.Float64() + mn[i]
		if len(stp) > 0 && stp[i] > 0 {r.Pos[i] = math.Floor(r.Pos[i]/stp[i])*stp[i]}
	}
	err := f(r, inp)
	return err
}

func rpopgen(np int, mx, mn, stp []float64, inp []interface{}, f robj) (pop []Rat, err error){
	pop = make([]Rat, np)
	for i := range pop{
		err = pop[i].init(mx, mn, stp, inp, f)
		if err != nil{
			return
		}
	}
	return
}


func rpopselec(pop []Rat, st, mx int) (pool []Rat, gb, gw float64, gpos []float64){
	switch st{
		case 0:
		//tournament selection
		pool, gb, gw, gpos = rselectour(pop, mx)
		case 1:
		//raka proportional selection
		pool, gb, gw, gpos = rselecprop(pop, mx)
		case 2:
		//fitness proportional selection
		pool, gb, gw, gpos = rselecfprop(pop, mx)
		case 3:
		//stochastic universal sampling
		pool, gb, gw, gpos = rselecsus(pop, mx)
		case 4:
		//rank selection
		pool, gb, gw, gpos = rselecrank(pop, mx)
	}
	return
}

func rselecsus(pop []Rat, mx int) (pool []Rat, gb, gw float64, gpos []float64){
	//stochastic universal sampling
	gpos = make([]float64, len(pop[0].Pos))
	var fs float64
	sort.Slice(pop, func(i, j int) bool {
		return pop[i].Fit > pop[j].Fit
	})
	for i:=0; i< len(pop);i++{
		fs += pop[i].Fit
		if i == 0{
			gb = pop[i].Fit
			gw = pop[i].Wt
			copy(gpos, pop[i].Pos)
		} else if gb < pop[i].Fit{
			gb = pop[i].Fit
			gw = pop[i].Wt
			copy(gpos, pop[i].Pos)
		}
	}
	dist := fs/float64(len(pop))
	shift := dist * rand.Float64()
	var brdrs []float64
	for i := 0; i < len(pop); i++{
		brdrs = append(brdrs, shift + float64(i) * dist)
	}
	for _, brd := range brdrs{
		rs := 0.0
		for i := range pop{
			rs += pop[i].Fit
			if rs >= brd{	
				pool = append(pool, pop[i])
				break
			}
		}
	}
	return
}

func rselecrank(pop []Rat, mx int) (pool []Rat, gb, gw float64, gpos []float64){
	//rank selection
	sort.Slice(pop, func(i, j int) bool {
		return pop[i].Fit > pop[j].Fit
	})
	gpos = make([]float64, len(pop[0].Pos))
	copy(gpos, pop[len(pop)-1].Pos)
	gb = pop[len(pop)-1].Fit; gw = pop[len(pop)-1].Wt
	rd := 1.0/float64(len(pop))
	var rsum float64
	rnks := make([]float64, len(pop))
	for i := 0; i < len(pop); i++{
		rsum += 1.0 - float64(i) * rd
		rnks[i] = 1.0 - float64(i) * rd
	}
	var rsp float64
	for len(pool) < len(pop){
		rsp = 0.
		for i := range pop{
			rsp += rnks[i]
			if rand.Float64() <= rsp/rsum{
				pool = append(pool, pop[i])
				break
			}
		}
	}
	return
}

func rselecfprop(pop []Rat, mx int) (pool []Rat, gb, gw float64, gpos []float64){
	//fitness proportional selection
	gpos = make([]float64, len(pop[0].Pos))
	var fs, fp float64
	for i:=0; i< len(pop);i++{
		fs += pop[i].Fit
		if i == 0{
			gb = pop[i].Fit
			gw = pop[i].Wt
			copy(gpos, pop[i].Pos)
		} else if gb < pop[i].Fit{
			gb = pop[i].Fit
			gw = pop[i].Wt
			copy(gpos, pop[i].Pos)
		}
	}
	sort.SliceStable(pop, func(i, j int) bool {
		return pop[i].Fit > pop[j].Fit
	})
	for len(pool) < len(pop){
		fp = 0.
		shave := rand.Float64()
		for _, b := range pop{ 
			fp += b.Fit
			if fp/fs >= shave {
				pool = append(pool, b)
			}
			break
		}
	}
	return
}

func rselectour(pop []Rat, mx int) (pool []Rat, gb, gw float64, gpos []float64){
	//tournament selection (coello 94)
	gpos = make([]float64, len(pop[0].Pos))
	for i:=0; i< len(pop);i++{
		if i == 0{
			gb = pop[i].Fit
			gw = pop[i].Wt
			copy(gpos, pop[i].Pos)
		} else if gb < pop[i].Fit{
			gb = pop[i].Fit
			gw = pop[i].Wt
			copy(gpos, pop[i].Pos)
		}
		x := rand.Intn(len(pop))
		y := rand.Intn(len(pop))
		if pop[x].Fit > pop[y].Fit{
			pool = append(pool, pop[x])
		} else {pool = append(pool, pop[y])}
	}
	return
}

func rselecprop(pop []Rat, mx int) (pool []Rat, gb, gw float64, gpos []float64){
	gpos = make([]float64, len(pop[0].Pos))
	var fmax, fmin, favg float64
	if mx == -1{
		//minimize
		for i, b := range pop{
			if fmax < pop[i].Fit{
				fmax = b.Fit
			}
			if i == 0{
				fmin = b.Fit
			} else if fmin > b.Fit{
				fmin = b.Fit
			}
		}
	}
	for _, b := range pop{
		if mx == -1 {b.Fit = fmax + fmin - b.Fit}
		favg += b.Fit
	}
	favg = favg / float64(len(pop))
	for i, b := range pop{
		if i == 0{
			gb = b.Fit
			gw = b.Wt
			copy(gpos, b.Pos)
		} else if gb < b.Fit{
			gb = b.Fit
			gw = b.Wt
			copy(gpos, b.Pos)
		}
		bcount := int(math.Round(b.Fit/favg)) 
		if bcount == 0{
			continue
		}
		for i := 0; i < bcount; i++ {
			pool = append(pool, b)
		}
	}
	return
}
