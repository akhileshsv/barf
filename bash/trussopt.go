package barf

//TO DELETE - all old funcs must die
import (
	"fmt"
	"time"
	"math"
	"math/rand"
	"sort"
	kass "barf/kass"
	//"encoding/json"
)

type Rat struct {
	vec     []int
	fitness float64
	weight  float64
	pcrs    float64
	pmut    float64
	cons    int
	popsize int
	ngens   int
	mutyp   int
	ngen    int
	ctyp    int
	etyp    int
	//fitfunc func(*kass.Model, [][]float64, float64, float64, float64) float64
}

func MakeR(mod *kass.Model,dims [][]float64, nsecs int, dmax, pmax, dens float64) (r Rat) {
	ndims := len(dims)
	r.vec = make([]int, nsecs)
	for i := range r.vec{
		r.vec[i] = rand.Intn(ndims)
	}
	r.FitR(mod, dims, pmax, dmax, dens)
	return
}

func (r *Rat) FitR(mod *kass.Model, dims [][]float64, pmax, dmax, dens float64) {
	cp := make([][]float64, len(r.vec))
	for i, idx := range r.vec{
		cp[i] = make([]float64,1)
		cp[i][0] = dims[idx][0]
	}
	modr := kass.Model{
		Ncjt:mod.Ncjt,
		Coords:mod.Coords,
		Supports:mod.Supports,
		Mprp:mod.Mprp,
		Jloads:mod.Jloads,
		Em:mod.Em,
		Cp:cp,
	}
	var wt, C, gx float64
	var con int
	frmrez, err := kass.CalcTrs(&modr, mod.Ncjt)
	if err != nil {
		r.fitness = -1e20
		return
	}
	js, _ := frmrez[0].(map[int]*kass.Node)
	ms, _ := frmrez[1].(map[int]*kass.Mem)
	for _, node := range js{
		for _, d := range node.Displ {
			gx = math.Abs(d)/dmax - 1.0
			if gx > 0.0 {
				C += gx
				con += 1
			}
		}
	}
	for _, mem := range ms{
		wt += mem.Geoms[0] * mem.Geoms[2] * dens
		pmem := mem.Qf[0] / mem.Geoms[2]
		gx = math.Abs(pmem)/pmax - 1.0
		if gx > 0.0 {
			C += gx
			con += 1
		}
	}
	r.weight = wt
	r.fitness = wt * (1.0 + 10.0*C)
	r.cons = con
}

func CrsR(x, y *Rat, ctyp int, pcrs float64) (a,b Rat) {
	nsecs := len(x.vec)
	a.vec = make([]int, nsecs)
	b.vec = make([]int, nsecs)
	if nsecs < 3 {
		switch nsecs{
			case 1:
			//clone ofc
			a.vec = x.vec
			b.vec = y.vec
			case 2:
			a.vec[0], a.vec[1] = x.vec[0], y.vec[1]
			b.vec[0], b.vec[1] = y.vec[0], x.vec[1]
		}
		return
	}
	//clone if < pcrs
	if rand.Float64() > pcrs{
		ctyp = 0
	}
	switch ctyp{
		case 0:
		//clone
		copy(a.vec, x.vec)
		copy(b.vec, y.vec)
		case 1:	
		//1 point crossover
		cx := rand.Intn(len(x.vec)-2) + 1
		for i := range x.vec{
			if i < cx{
				a.vec[i] = x.vec[i]
				b.vec[i] = y.vec[i]
			} else {
				a.vec[i] = y.vec[i]
				b.vec[i] = x.vec[i]
			}
		}
		case 2:
		//n point crossover
		//might not break even if rand returns similar vals
		//n := rand.Intn(len(x.vec)-2) + 1
		n := 2
		var ps []int
		mps := make(map[int]bool)
		var i int
		for i < n{
			ndx := rand.Intn(len(x.vec)-2)+1
			if _, ok := mps[ndx]; !ok{
				ps = append(ps, ndx)
				mps[ndx] = true
			}
		}
		copy(a.vec, x.vec)
		copy(b.vec, y.vec)
		ps = append(ps, 0)
		ps = append(ps, len(x.vec)-1)
		sort.Slice(ps, func(i, j int) bool {
			return ps[i] < ps[j]
		})
		for i := 0; i < n+1; i++ {
			if i%2 == 0{
				continue
			}
			a.vec[ps[i]] = y.vec[ps[i]]
			a.vec[ps[i+1]] = y.vec[ps[i+1]]
			b.vec[ps[i]] = x.vec[ps[i]]
			b.vec[ps[i+1]] = x.vec[ps[i+1]]
		}
		case 3:
		//uniform crossover
		for i := range x.vec{
			if rand.Float64() < 0.5{
				a.vec[i], b.vec[i] = y.vec[i], x.vec[i]
			}
		}
		case 4:
		//ordered crossover (simple repair)
		n1 := rand.Intn(len(x.vec)-2)+1
		n2 := rand.Intn(len(x.vec)-2)+1
		switch{
			case n1 < n2:
			for i := range x.vec{
				switch{
					case i>=n1 && i<=n2:
					a.vec[i] = y.vec[i]
					b.vec[i] = x.vec[i]
					default:
					a.vec[i] = x.vec[i]
					b.vec[i] = y.vec[i]
				}
			} 
			case n1 > n2:
			for i := range x.vec{
				switch{
					case i>=n2 && i<=n1:
					a.vec[i] = y.vec[i]
					b.vec[i] = x.vec[i]
					default:
					a.vec[i] = x.vec[i]
					b.vec[i] = y.vec[i]
				}
			}
			default:
			//clone by law of dice
			for i := range x.vec{
				if i == n1{
					a.vec[i] = y.vec[i]
					b.vec[i] = x.vec[i]
				} else {
					a.vec[i] = x.vec[i]
					b.vec[i] = y.vec[i]
				}
			}
		}
	}
	return
}

func (r *Rat) MutR(ndims int, mutyp int){
	mumap := map[int]string{0:"flip,shuffle",1:"flip",2:"shuffle",3:"exchange",4:"bounded shift"}
	fmt.Println("mutyp->",mumap[mutyp])
	fmt.Println("in->",r.vec)
	switch mutyp{
	case 0:
		r.MutRFlip(ndims)
		r.MutRShuffle()
	case 1:
		r.MutRFlip(ndims)
	case 2:
		r.MutRShuffle()
	case 3:
		r.MutREx()
	case 4:
		r.MutRShift()
	}
	fmt.Println("oot->",r.vec)
}

func (r *Rat) MutRFlip(ndims int){
	//flip int
	mutpt := rand.Intn(len(r.vec))
	mutdx := rand.Intn(ndims)
	r.vec[mutpt] = mutdx	
}

func (r *Rat) MutRShuffle(){
	//shuffle mutation
	idx := rand.Intn(len(r.vec)-1) 
	for i := len(r.vec) - 1; i > idx; i-- {
		j := rand.Intn(i + 1)
		r.vec[i], r.vec[j] = r.vec[j], r.vec[i]
	}
}

func (r *Rat) MutREx(){
	//exchange mutation
	i, j := rand.Intn(len(r.vec)), rand.Intn(len(r.vec))
	r.vec[i], r.vec[j] = r.vec[j], r.vec[i]
}

func (r *Rat) MutRShift(){
	//bounded shift mutation
	//pick index and random int 
	x := rand.Intn(len(r.vec)) 
	y := rand.Intn(len(r.vec)) 
	switch {
	case x > y:
		//shift at x left by x - y 
		shift := x - y
		y0 :=  r.vec[y]
		for i := y; i < x; i++{
			r.vec[i] = r.vec[i + shift]
		}
		r.vec[x] = y0
	default:
		//shift at x right by y - x
		shift := y - x
		x0 := r.vec[x]
		for i:= x; i < y; i++{
			r.vec[i] = r.vec[i + shift]
		}
		r.vec[y] = x0
	}
}

func GenPop(mod *kass.Model, dims [][]float64, dmax, pmax, dens float64, nsecs, popsize int) ([]*Rat) {
	var popr []*Rat
	for i := 0; i < popsize; i++ {
		r := MakeR(mod, dims, nsecs, dmax, pmax, dens)
		popr = append(popr, &r)
	}
	return popr
}

func SPop(popr []*Rat, styp int) (matr []*Rat) {
	switch styp{
		case 0:
		//tournament selection
		for i:=0; i< len(popr);i++{
			x := rand.Intn(len(popr))
			y := rand.Intn(len(popr))
			if popr[x].fitness > popr[y].fitness {
				matr = append(matr, popr[x])
			} else {matr = append(matr, popr[y])}
		}
		case 1:
		//proportional selection
		case 2:
		case 3:
		
	}
	return
}

func EPopCru(popr []*Rat) (matr []*Rat, fpop, wmin float64, midx int, cz *Rat) {
	//coello fitness
	//tournament selection (2)
	var czmin float64
	for idx, r := range popr{
		r.fitness = 1.0/(r.weight*(1000.0*float64(r.cons) + 1.0))
		fmt.Println("cru fit",r.fitness)
		if fpop < r.fitness{
			fpop = r.fitness
			midx = idx
		}
		if r.cons == 0{
			if czmin == 0.0{
				cz = r
				czmin = r.weight
			} else if cz.weight > r.fitness{
				cz = r
				czmin = r.weight
			}
		}
	}
	//tournament selection
	for i:=0; i< len(popr);i++{
		x := rand.Intn(len(popr))
		y := rand.Intn(len(popr))
		fmt.Println("cru selec->",popr[x].fitness, popr[y].fitness)
		if popr[x].fitness > popr[y].fitness {
			matr = append(matr, popr[x])
		} else {matr = append(matr, popr[y])}
	}
	fmt.Println("cru pool->",len(matr))
	return
}

func EPopRaka(popr []*Rat) (matr []*Rat, fpop, wmin float64, midx int, cz *Rat) {
	//rajeev fitness
	//proportional selection
	var fmax, fmin, favg, czmin float64
	for i, r := range popr{
		if fmax < r.fitness{
			fmax = r.fitness
		}
		if i == 0{
			fmin = r.fitness
		} else if fmin > r.fitness{
			fmin = r.fitness
		}
	}
	for _, r := range popr{
		r.fitness = fmax + fmin - r.fitness
		favg += r.fitness
	}
	favg = favg / float64(len(popr))
	for i, r := range popr{
		if fpop < r.fitness{
			fpop = r.fitness
			wmin = r.weight
			midx = i
		}
		rcount := int(math.Round(r.fitness/favg))
		if r.cons == 0{
			if czmin == 0.0{
				czmin = r.weight
				cz = r
			} else if czmin > r.weight{
				czmin = r.weight
				cz = r
			}
		}
		if rcount == 0{
			continue
		}
		for i := 0; i < rcount; i++ {
			matr = append(matr, r)
		}
	}
	return
}

func RePop(matr []*Rat, mod *kass.Model, dims [][]float64, dmax, pmax, dens, pmut, pcrs float64, popsize, mutyp, ctyp int) ([]*Rat) {
	var popr []*Rat
	for len(popr) <= popsize {
		if len(matr) == 0 {
			fmt.Println("ERRORE, errore->wtf zero mates")
		}
		var mutz float64
		x := rand.Intn(len(matr))
		y := rand.Intn(len(matr))
		if x == y {
			continue
		}
		a, b := CrsR(matr[x], matr[y], ctyp, pcrs)
		mutz = rand.Float64()
		if mutz < pmut {
			a.MutR(len(dims),mutyp)
		}
		a.FitR(mod,dims,pmax,dmax,dens)
		mutz = rand.Float64()
		if mutz < pmut {
			b.MutR(len(dims),mutyp)
		}
		b.FitR(mod,dims,pmax,dmax,dens)
		popr = append(popr, &a)
		popr = append(popr, &b)
	}
	return popr
}

func TrussOptSpa(mod *kass.Model, dims [][]float64, dmax, pmax, dens float64, nsecs int, rchan chan []*Rat) {
	//opt via SpaGhetti on wall
	rand.Seed(time.Now().UnixNano())
	//ngens := rand.Intn(35-25)+25
	ngens := 25
	pmut := 0.01 + rand.Float64() * (0.01)
	popsize := rand.Intn(50-30)+30
	popr := GenPop(mod, dims, dmax, pmax, dens, nsecs, popsize)
	mutyp := rand.Intn(4-3)+3
	etyp := 0
	ctyp := rand.Intn(2)+1
	//ctyp := 1
	pcrs := 0.8
	//var mdx, midx int
	//var fmax, wmin float64
	//var cdxs []int
	var zerocs []*Rat
	var cz *Rat
	for i := 0; i < ngens; i++ {
		if etyp == 0{
			//raka
			popr, _, _, _, cz = EPopRaka(popr)
		} else {
			//cru
			popr, _, _, _, cz = EPopCru(popr)
		}
		if cz != nil {
			zerocs = append(zerocs, cz)
			cz.pmut = pmut
			cz.popsize = popsize
			cz.ngens = ngens
			cz.mutyp = mutyp
			cz.ngen = i
			fmt.Println(kass.ColorBlue,"generation->", i, "min weight->", kass.ColorRed,cz.weight,kass.ColorReset)
		}
		//fmt.Println("generation->", i, "max fitness->", fmax, "weight->", wmin)
		popr = RePop(popr, mod, dims, dmax, pmax, dens, pmut, pcrs, popsize, mutyp, ctyp)
	}
	fmt.Println(kass.IconCubes,"\nfinito.")
	rchan <- zerocs
}

func TrussOptRaPa(mod *kass.Model, dims [][]float64, dmax, pmax, dens float64, nsecs, npa int) {
	pltchn := make(chan string, 1)
	go kass.PlotTrs2d(mod, "dumb", pltchn)
	pltstr := <-pltchn
	fmt.Println(pltstr)
	close(pltchn)
	rchan := make(chan []*Rat, npa)
	var minwt, pmut float64
	var popsize int
	minwt = 1e14
	for i := 0; i < npa; i++ {
		go TrussOptSpa(mod, dims, dmax, pmax, dens, nsecs, rchan)
	}
	var rats []*Rat
	for i := 0; i < npa; i++ {
		zerocs := <- rchan
		for _, r := range zerocs {
			if minwt > r.weight {minwt = r.weight; pmut = r.pmut; popsize = r.popsize}
		}
		rats = append(rats, zerocs...)
	}
	if len(rats) == 0 {
		return
	}
	sort.Slice(rats, func(i, j int) bool {
		return rats[i].weight < rats[j].weight
	})
	
	fmt.Println(kass.ColorYellow)
	fmt.Println("***minwt->",minwt,"<-***")
	fmt.Println("pmut->",pmut, popsize)
	
	for i:= range rats{
		if i > 5 {break}
		fmt.Println(kass.ColorRed)
		fmt.Println(rats[i].weight, rats[i].popsize, rats[i].ngen, rats[i].ngens, rats[i].pmut, rats[i].mutyp)
	}
}
