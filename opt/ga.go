package barf

import (
	"fmt"
	"time"
	"math"
	"math/rand"
	"sort"
	draw"barf/draw"
)

type Bat struct{
	Pos []int
	Fit float64
	Wt  float64
	Con float64
}

//okay forget Bobj - just take in a goddamn slice or interface
type Bobj func(*Bat, []interface{}) error

func (b *Bat) Init(nd []int, inp []interface{}, f Bobj) error{
	b.Pos = make([]int, len(nd))
	for i := range b.Pos{
		b.Pos[i] = rand.Intn(nd[i])
	}
	err := f(b, inp)
	return err
}

func (b *Bat) Create(pos []int, inp []interface{}, f Bobj, fit float64) error{
	b.Pos = make([]int, len(pos))
	copy(b.Pos, pos)
	if fit != -1.0{
		b.Fit = fit
		return nil
	}
	err := f(b, inp)
	return err
}

func binit(b *Bat, np int, nd []int, inp []interface{}, f Bobj, bchn chan error){
	//parallel init
	err := b.Init(nd, inp, f)
	bchn <- err
}

func bfit(b *Bat, f Bobj, inp []interface{}, bchn chan error){
	err := f(b, inp)
	bchn <- err
}

func Bpopgen(par bool, np int, nd []int, inp []interface{}, f Bobj) (pop []Bat, err error){
	pop = make([]Bat, np)
	if par{
		//paralell eval
		bchn := make(chan error,np)
		for i := range pop{
			go binit(&pop[i], np, nd, inp, f, bchn)
		}
		for _ = range pop{
			err = <- bchn
			if err != nil {
				//return
			}
		}
		return
	}
	for i := range pop{
		err = pop[i].Init(nd, inp, f)
		if err != nil{
			//return
		}
	}
	return
}

func Brepop(par, ndchk bool, np, mt, ct, cn, mx int, pmut, pcrs float64, nd []int, inp []interface{}, pool []Bat, f Bobj) (pop []Bat, err error){
	for len(pop) < np{
		if len(pool) < np{
			err = ErrSel
			return
		}
		var mutz float64
		x := rand.Intn(len(pool))
		y := rand.Intn(len(pool))
		if x == y {
			continue
		}
		a, b := Bcrs(pool[x], pool[y], ct, cn, pcrs, nd, ndchk)
		mutz = rand.Float64()
		if mutz < pmut {
			a.Mut(ndchk,nd,mt)
		}
		//f(&a,inp)		
		mutz = rand.Float64()
		if mutz < pmut {
			b.Mut(ndchk,nd,mt)
		}
		//f(&b,inp)
		pop = append(pop, a)
		pop = append(pop, b)
	}
	switch par{
		case true:
		bchn := make(chan error, np)
		for i := range pop{
			go bfit(&pop[i], f, inp, bchn)
		}
		for _ = range pop{
			_ = <- bchn
		}
		case false:
		
		for i := range pop{
			f(&pop[i],inp)
		}
	}
	return
}

func Popselec(pop []Bat, st, mx int) (pool []Bat, gb, gw, avg float64, gpos []int){
	switch st{
		case 1:
		//tournament selection
		pool, gb, gw, avg, gpos = selectour(pop, mx)
		case 2:
		//fitness proportional selection
		pool, gb, gw, avg, gpos = selecfprop(pop, mx)
		case 3:
		//stochastic universal sampling
		pool, gb, gw, avg, gpos = selecsus(pop, mx)
		case 4:
		//rank selection
		pool, gb, gw, avg, gpos = selecrank(pop, mx)
		case 5:
		//raka proportional selection
		//IT IS BRAKE
		pool, gb, gw, avg, gpos = selecprop(pop, mx)

	}
	return
}

func selecsus(pop []Bat, mx int) (pool []Bat, gb, gw, avg float64, gpos []int){
	//stochastic universal sampling
	gpos = make([]int, len(pop[0].Pos))
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
	avg = dist
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
			}
		}
	}
	return
}

func selecrank(pop []Bat, mx int) (pool []Bat, gb, gw, avg float64, gpos []int){
	//rank selection
	sort.Slice(pop, func(i, j int) bool {
		return pop[i].Fit > pop[j].Fit
	})
	gpos = make([]int, len(pop[0].Pos))
	copy(gpos, pop[len(pop)-1].Pos)
	gb = pop[len(pop)-1].Fit; gw = pop[len(pop)-1].Wt
	rd := 1.0/float64(len(pop))
	var rsum float64
	rnks := make([]float64, len(pop))
	for i := 0; i < len(pop); i++{
		rsum += 1.0 - float64(i) * rd
		avg += pop[i].Fit
		rnks[i] = 1.0 - float64(i) * rd
	}
	avg = avg/float64(len(pop))
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

func selecfprop(pop []Bat, mx int) (pool []Bat, gb, gw, avg float64, gpos []int){
	//fitness proportional selection
	gpos = make([]int, len(pop[0].Pos))
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
	avg = fs/float64(len(pop))
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

func selectour(pop []Bat, mx int) (pool []Bat, gb, gw, avg float64, gpos []int){
	//tournament selection (coello 94)
	gpos = make([]int, len(pop[0].Pos))
	for i:=0; i< len(pop);i++{
		avg += pop[i].Fit
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
	avg = avg/float64(len(pop))
	for _ = range pop{
		x := rand.Intn(len(pop))
		y := rand.Intn(len(pop))
		z := rand.Intn(len(pop))
		w := x
		for _, i := range []int{y,z}{
			if pop[i].Fit > pop[w].Fit{
				w = i
			}
		}
		pool = append(pool, pop[w])
	}
	return
}

func selecprop(pop []Bat, mx int) (pool []Bat, gb, gw, avg float64, gpos []int){
	//THIS IS BRAKED
	//raka prop selection
	gpos = make([]int, len(pop[0].Pos)) 
	for i, b := range pop{
		avg += b.Fit
		if i == 0{
			gb = b.Fit
			gw = b.Wt
			copy(gpos, b.Pos)
		} else if gb < b.Fit{
			gb = b.Fit
			gw = b.Wt
			copy(gpos, b.Pos)
		}
	}
	avg = avg / float64(len(pop))
	for _, b := range pop{
		bcount := int(math.Ceil(b.Fit/avg))
		if bcount == 0{bcount = 1}
		for j := 0; j < bcount; j++{
			pool = append(pool, b)
		}
	}
	return
}

func popavg(gavs []float64, step int) (gavg float64){
	if len(gavs) < step{
		return 
	}
	for i := 0; i < step; i++{
		gavg += gavs[len(gavs)-1-i]
	}
	gavg = gavg/float64(step)
	return
}

func isimp(gavs []float64, step int, gap float64) (bool){
	gavg := popavg(gavs, step)
	if gavg == 0.0{return true}
	return gavs[len(gavs)-1] > gavg * (1.0 - gap)
}

func Gabloop(web, par, ndchk bool, np, ng, mt, ct, cn, st, mx int, pmut, pcrs float64, nd []int, inp []interface{},f Bobj, drw, trm, title string) ([]int, float64, float64, string, error){
	//basic ga loop
	rand.Seed(time.Now().UnixNano())	
	var pool []Bat
	var gb, gw, gavg float64
	var gpos []int
	var fdat, wdat, pltstr string
	var gavs []float64
	pop, err := Bpopgen(par, np, nd, inp, f)
	
	if err != nil{
		return gpos, gw, gb, pltstr, err
	}
	for gen := 0; gen < ng; gen++{
		//fmt.Println("gen-",gen)
		pool, gb, gw, gavg, gpos = Popselec(pop, st, mx)
		gavs = append(gavs, gavg)
		if mx == -1{
			wdat += fmt.Sprintf("%v %f\n",gen, gw) 
		}
		fdat += fmt.Sprintf("%v %f\n",gen, gb)
		switch drw{
			case "all":
			if !web{fmt.Println(ColorBlue, "gen->",gen,ColorWhite,"\nglobal best->\n",gpos, ColorGreen,"\tmax fitness\tweight->", gb, gw, ColorReset)}
			default:
			if gen % 10 == 0{
				if !web{fmt.Println(ColorBlue, "gen->",gen,ColorWhite, ColorGreen,"\tmax fitness",gb,"\tweight->", gw, "\tavg",gavg,ColorReset)}
			}
		}
		pop, err = Brepop(par, ndchk, np, mt, ct, cn, mx, pmut, pcrs, nd, inp, pool, f)
		if err != nil{
			fmt.Println("ERRORE-",err)
			return gpos, gw, gb, pltstr, err
		}
	}
	skript := "d2.gp"
	switch drw{
		case "gen":
		var folder string
		if web{folder = "web"}
		if mx == -1{
			pltstr, _ = draw.Draw(wdat, skript, trm, folder, title, title, "gen", "weight","")
		} else {
			pltstr, _ = draw.Draw(fdat, skript, trm, folder, title, title, "gen","fit","")

		}		
	}
	if !web{fmt.Println(ColorCyan, "best pos->\n", gpos, ColorRed, "\nmax fitness->", gb, ColorReset)}
	if mx == -1{
		if !web{fmt.Println(ColorCyan, "min weight->", gw, ColorReset)}
	}
	return gpos, gw, gb, pltstr, err
}


func Gabaloop(web, par, ndchk bool, np, npmn, npmx, ngmn, ngmx, mt, ct, cn, st, mx, step int, pmut, pcrs, gap float64, nd []int, inp []interface{},f Bobj, drw, trm, title string) ([]int, float64, float64, string, error){
	rand.Seed(time.Now().UnixNano())	
	var pool []Bat
	var gb, gw, avg float64
	var gpos []int
	var fdat, wdat, pltstr string
	var gavs []float64
	pop, err := Bpopgen(par, np, nd, inp, f)
	if err != nil{
		return gpos, gw, gb, pltstr, err
	}
	pcrs = 0.95
	for gen := 0; gen < ngmx; gen++{
		if !web{fmt.Println("gen->",gen,"gpos->",gpos,"gb->",gb,"gw->",gw,"pmut",pmut,"pcrs",pcrs)}
		if gen > ngmn && !isimp(gavs, 40, gap){
			if !web{fmt.Println("no improvement for 40, stopping")}
			break
		}
		if gen > step{
			if !isimp(gavs, step, gap){
				//fmt.Println(ColorGreen,"tweaking down",ColorReset)
				pmut = pmut * 0.9
				pcrs = pcrs * 0.9
				//sort.SliceStable(pop, func(i, j int) bool {
				//	return pop[i].Fit > pop[j].Fit
				//})
				//remove worst individuals
				pop = pop[:len(pop)-5]
				pop = append(pop[:5],pop...)
				
				sort.SliceStable(pop, func(i, j int) bool {
					return pop[i].Fit > pop[j].Fit
				})
				//np -= 2
			} else {
				//fmt.Println(ColorRed,"tweaking up",ColorReset)
				pmut = pmut * 1.03
				pcrs = pcrs * 1.03
				//add best individual thrice
				//pop = pop[:len(pop)-3]
				//pop = append([]Bat{pop[0],pop[0],pop[0]},pop...)
				
				//np += 2
				//pop = append(pop, []Bat{pop[0],pop[0],pop[0],pop[0],pop[0]}...)
				//pop = append(pop, pop[1])
			}
		}
		if np < npmn {np = npmn}
		if np > npmx {np = npmx}
		if pmut < 0.02{pmut = 0.02}
		//if pmut > 1.0{pmut = 1.0}
		if pcrs < 0.1{pcrs = 0.1}
		if pcrs > 1.0{pcrs = 0.98; pmut = 0.02}
		if pmut > 0.9{pmut = 0.9; pcrs = 0.1}
		//if gen > ngmn && !isimp(gavs, ngmn, gap){break}
		pool, gb, gw, avg, gpos = Popselec(pop, st, mx)
		gavs = append(gavs, avg)
		if mx == -1{
			wdat += fmt.Sprintf("%v %f\n",gen, gw) 
		}
		fdat += fmt.Sprintf("%v %f\n",gen, gb)
		// switch drw{
		// 	case "all":
		// 	fmt.Println(ColorBlue, "gen->",gen,ColorWhite,"\nglobal best->\n",gpos, ColorGreen,"\tmax fitness\tweight->", gb, gw, ColorReset)
		// 	default:
		// 	if gen % 10 == 0{
		// 		fmt.Println(ColorBlue, "gen->",gen,ColorWhite, ColorGreen,"\tmax fitness",gb,"\tweight->", gw, ColorReset)
		// 	}
		// }
		pop, err = Brepop(par, ndchk, np, mt, ct, cn, mx, pmut, pcrs, nd, inp, pool, f)
		if err != nil{return gpos, gw, gb, pltstr, err}
	}
	skript := "d2.gp"
	switch drw{
		case "gen":
		var folder string
		if web{folder = "web"}
		if mx == -1{
			pltstr, _ = draw.Draw(wdat, skript, trm, folder, title, title, "gen", "weight","")
		} else {
			pltstr, _ = draw.Draw(fdat, skript, trm, folder, title, title, "gen","fit","")

		}		
	}
	if !web{
		fmt.Println(ColorCyan, "best pos->\n", gpos, ColorRed, "\nmax fitness->", gb, ColorReset)
		if mx == -1{fmt.Println(ColorCyan, "min weight->", gw, ColorReset)}
	}
	return gpos, gw, gb, pltstr, err
}


/*


func galoop(ndchk bool, np, ng, mt, ct, cn, st, mx int, pmut, pcrs float64, nd []int, inp []interface{}, f Bobj, drw, trm, title string){
	rand.Seed(time.Now().UnixNano())	
	pop := bpopgen(np, nd, con, chc, f)
	var pool []Bat
	var gb, gw float64
	var gpos []int
	var dat string
	for gen := 0; gen < ng; gen++{
		pool, gb, gw, gpos = popselec(pop, st, mx)
		if mx == -1{
			dat += fmt.Sprintf("%v %f\n",gen, gw) 
		} else {dat += fmt.Sprintf("%v %f\n",gen, gb)}
		switch drw{
			case "all", "gen":
			fmt.Println(ColorBlue, "gen->",gen,ColorWhite,"\nglobal best->\n",gpos, ColorGreen,"\tmax fitness\tweight->", gb, gw, ColorReset)
			default:
			if gen % 10 == 0{
				fmt.Println(ColorBlue, "gen->",gen,ColorWhite, ColorGreen,"\tmax fitness\tweight->", gb, gw, ColorReset)
			}
		}
		pop = rebpop(ndchk, np, mt, ct, cn, mx, pmut, pcrs, nd, con, chc, pool, f)
	}
	skript := "d2.gp"
	var dstr string
	switch drw{
		case "gen":
		dstr = draw.Dumb(dat, skript, trm, title, "gen", "gbest", "")
		if trm != "qt" {fmt.Println(dstr)}
	}
	fmt.Println(ColorCyan, "best pos->\n", gpos, ColorRed, "\nmin fitness->", gb, ColorReset)

}

*/


/*
func repop(ndchk bool, np, mt, ct, cn, mx int, pmut, pcrs float64, nd []int, con []float64, chc [][]float64, pool []Bat, f Bobj) (pop []Bat){
	for len(pop) <= len(pool){
		if len(pool) == 0{
			fmt.Println("ERRORE, errore-> zero length pool wtf")
			return
		}
		var mutz float64
		x := rand.Intn(len(pool))
		y := rand.Intn(len(pool))
		if x == y {
			continue
		}
		a, b := rcrs(pool[x], pool[y], ct, cn, pcrs, mx, mn)
		mutz = rand.Float64()
		if mutz < pmut {
			a.Mut(nd,mt)
		}
		a.Fit, a.Wt = f(a.Pos,con,chc)
		mutz = rand.Float64()
		if mutz < pmut {
			b.Mut(nd,mt)
		}
		b.Fit, b.Wt = f(b.Pos, con, chc)
		pop = append(pop, a)
		pop = append(pop, b)
	}
	return
}
*/
