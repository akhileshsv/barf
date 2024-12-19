package barf

import (
	//"time"
	"math"
	"math/rand"
	"sort"
	//draw"barf/draw"
)

func Rcrs(x, y Rat, ctyp int, pcrs, alp float64, mx, mn, stp []float64) (a, b Rat) {
	//ctyp : 0 - clone, 1- 1 pt crs, 2 - n pt crs, 3- uniform, 4 - ordered
	if rand.Float64() > pcrs{
		ctyp = 0
	}
	a.Pos = make([]float64, len(x.Pos))
	b.Pos = make([]float64, len(x.Pos))
	switch ctyp{
		case 0:
		//clone
		copy(a.Pos, x.Pos)
		copy(b.Pos, x.Pos)
		case 1:
		//linear crossover
		for i := range x.Pos{
			a.Pos[i] = math.Round((x.Pos[i] + alp * (y.Pos[i] - x.Pos[i]))*1000.0)/1000.0
			b.Pos[i] = math.Round((y.Pos[i] - alp * (y.Pos[i] - x.Pos[i]))*1000.0)/1000.0
			if len(stp) > 0 && stp[i] > 0{
				a.Pos[i] = math.Floor(a.Pos[i]/stp[i])*stp[i]
				b.Pos[i] = math.Floor(b.Pos[i]/stp[i])*stp[i]
			}
		}
		case 2:
		//blend crossover
		for i := range x.Pos{
			l := math.Min(x.Pos[i],y.Pos[i]) - alp * math.Abs(y.Pos[i] - x.Pos[i])
			h := math.Max(x.Pos[i],y.Pos[i]) + alp * math.Abs(y.Pos[i] - x.Pos[i])
			a.Pos[i] = math.Round((l + rand.Float64() * (h - l))*1000.0)/1000.0
			b.Pos[i] = math.Round((l + rand.Float64() * (h - l))*1000.0)/1000.0
			if len(stp) > 0 && stp[i] > 0{
				a.Pos[i] = math.Floor(a.Pos[i]/stp[i])*stp[i]
				b.Pos[i] = math.Floor(b.Pos[i]/stp[i])*stp[i]
			}
			if len(mx) != 0{
				if a.Pos[i] > mx[i] {a.Pos[i] = mx[i]}
				if b.Pos[i] > mx[i] {b.Pos[i] = mx[i]}
				if a.Pos[i] < mx[i] {a.Pos[i] = mn[i]}
				if b.Pos[i] < mx[i] {b.Pos[i] = mn[i]}
			}
		}
	}
	return
}

func Bcrs(x, y Bat, ctyp, cn int, pcrs float64, nd []int, ndchk bool) (a, b Bat) {
	//THIS IS FUCKED
	//add nd in params, check for valid zomes by dims
	ndims := len(x.Pos)
	a.Pos = make([]int, ndims)
	b.Pos = make([]int, ndims)
	if ndims < 3 {
		switch ndims{
			case 1:
			//clone ofc
			a.Pos[0] = x.Pos[0]
			b.Pos[0] = y.Pos[0]
			case 2:
			a.Pos[0], a.Pos[1] = x.Pos[0], y.Pos[1]
			b.Pos[0], b.Pos[1] = y.Pos[0], x.Pos[1]
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
		copy(a.Pos, x.Pos)
		copy(b.Pos, y.Pos)
		case 1:	
		//1 point crossover
		cx := rand.Intn(len(x.Pos)-2) + 1
		for i := range x.Pos{
			if i < cx{
				a.Pos[i] = x.Pos[i]
				b.Pos[i] = y.Pos[i]
				//if a.Pos[i] > nd[i]{a.Pos[i] = rand.Intn(nd[i])}
			} else {
				a.Pos[i] = y.Pos[i]
				b.Pos[i] = x.Pos[i]
			}
		}
		case 2:
		//n point crossover
		//still braked
		//n := rand.Intn(len(x.Pos)-2) + 1
		if cn == 0{cn = 2}
		if cn + 1 == len(x.Pos){cn = len(x.Pos) - 2}
		var ps []int
		mps := make(map[int]bool)
		var i int
		for i < cn{
			ndx := rand.Intn(len(x.Pos)-2)+1
			if _, ok := mps[ndx]; !ok{
				ps = append(ps, ndx)
				mps[ndx] = true
				i++
			}
		}
		copy(a.Pos, x.Pos)
		copy(b.Pos, y.Pos)
		ps = append(ps, 0)
		ps = append(ps, len(x.Pos)-1)
		sort.Slice(ps, func(i, j int) bool {
			return ps[i] < ps[j]
		})
		//fmt.Println(ps)
		for i := 0; i < cn+1; i++ {
			if i%2 == 0{
				continue
			}
			a.Pos[ps[i]] = y.Pos[ps[i]]
			a.Pos[ps[i+1]] = y.Pos[ps[i+1]]
			b.Pos[ps[i]] = x.Pos[ps[i]]
			b.Pos[ps[i+1]] = x.Pos[ps[i+1]]
		}
		case 3:
		//uniform crossover
		copy(a.Pos, x.Pos)
		copy(b.Pos, y.Pos)
		for i := range x.Pos{
			if rand.Float64() - 0.5 < 0{
				a.Pos[i], b.Pos[i] = y.Pos[i], x.Pos[i]
			}
		}
		case 4:
		//ordered crossover (simple repair)
		n1 := rand.Intn(len(x.Pos)-2)+1
		n2 := rand.Intn(len(x.Pos)-2)+1
		switch{
			case n1 < n2:
			for i := range x.Pos{
				switch{
					case i>=n1 && i<=n2:
					a.Pos[i] = y.Pos[i]
					b.Pos[i] = x.Pos[i]
					default:
					a.Pos[i] = x.Pos[i]
					b.Pos[i] = y.Pos[i]
				}
			} 
			case n1 > n2:
			for i := range x.Pos{
				switch{
					case i>=n2 && i<=n1:
					a.Pos[i] = y.Pos[i]
					b.Pos[i] = x.Pos[i]
					default:
					a.Pos[i] = x.Pos[i]
					b.Pos[i] = y.Pos[i]
				}
			}
			default:
			//clone by law of dice
			for i := range x.Pos{
				if i == n1{
					a.Pos[i] = y.Pos[i]
					b.Pos[i] = x.Pos[i]
				} else {
					a.Pos[i] = x.Pos[i]
					b.Pos[i] = y.Pos[i]
				}
			}
		}
		case 5:
		//laplace crossover
		//ndchk = true
		var une, dos, beta float64
		dos = 0.35
		u := rand.Float64()
		for i, x1 := range x.Pos{
			x1 := float64(x1)
			x2 := float64(y.Pos[i])
			if u <= 0.5{
				beta = une - dos * math.Log(u)
			} else {
				beta = une + dos * math.Log(u)
			}
			a1 := x1 + beta * math.Abs(x1 - x2)
			b1 := x2 + beta * math.Abs(x1 - x2)
			a.Pos[i] = int(a1); b.Pos[i] = int(b1)
		}
		
	}
	if ndchk{
		for i := range a.Pos{
			if a.Pos[i] > nd[i]-1{
				//a.Pos[i] = y.Pos[i]
				a.Pos[i] = nd[i]-1
			}
			if b.Pos[i] > nd[i]-1{
				//b.Pos[i] = x.Pos[i]
				b.Pos[i] = nd[i]-1
			}
			if a.Pos[i] < 0{a.Pos[i] = 0}
			if b.Pos[i] < 0{b.Pos[i] = 0}
		}
	}
	//fmt.Println(a.Pos, b.Pos)
	return

}

func (b *Bat) Mut(ndchk bool, nd []int, mt int){
	tmp := make([]int, len(b.Pos))
	if ndchk{copy(tmp, b.Pos)}
	if mt == 0{mt = 1}
	switch mt {
	case 0:
		b.Flip(nd)
		b.Shuffle()
	case 1:
		b.Flip(nd)
	case 2:
		b.Shuffle()
	case 3:
		b.Exchange()
	case 4:
		b.Shift()
	case 5:
		b.Invert()
	case 6:
		b.Powmut(nd)
	}
	if ndchk{
		for i, v := range b.Pos{
			if v > nd[i]-1{v = nd[i]-1}
			if v < 0{v = 0}
		}
	}
}

func (b *Bat) Powmut(nd []int){
	//power mutation, deep/singh
	pmdx := 4.0
	tmp := make([]int, len(b.Pos))
	s := math.Pow(rand.Float64(), pmdx)
	var x float64
	pr := rand.Float64()
	for i, v := range b.Pos{
		if v == nd[i] - 1{
			x = float64(v) - s * float64(v)
		} else {
			t := float64(v)/float64(nd[i]-1)
			if t < pr{
				x = float64(v) - s * float64(v)
			} else {
				x = float64(v) + float64(nd[i] - v - 1) 
			}

		}
		tmp[i] = int(x)
	}
	copy(b.Pos, tmp)
}

func (r *Rat) Mut(mt int, pmut float64, mx, mn []float64){
	switch mt{
		case 0:
		//random deviation mutation
		for i := range r.Pos{
			if rand.Float64() < pmut{
				r.Pos[i] = r.Pos[i] + r.Pos[i] * rand.Float64()
			}
			if len(mx) > 0{
				if r.Pos[i] < mn[i] {r.Pos[i] = mn[i]}
				if r.Pos[i] > mx[i] {r.Pos[i] = mx[i]}
			}
		}
	}
}

func (b *Bat) Flip(nd []int){
	//flip int
	mutpt := rand.Intn(len(b.Pos))
	mutdx := rand.Intn(nd[mutpt])
	b.Pos[mutpt] = mutdx	
}

func (b *Bat) Shuffle(){
	//shuffle mutation
	idx := rand.Intn(len(b.Pos)-1) 
	for i := len(b.Pos) - 1; i > idx; i-- {
		j := rand.Intn(i + 1)
		b.Pos[i], b.Pos[j] = b.Pos[j], b.Pos[i]
	}
}

func (b *Bat) Exchange(){
	//exchange mutation
	i, j := rand.Intn(len(b.Pos)), rand.Intn(len(b.Pos))
	b.Pos[i], b.Pos[j] = b.Pos[j], b.Pos[i]
}

func (b *Bat) Shift(){
	//bounded shift mutation
	x := rand.Intn(len(b.Pos)) 
	y := rand.Intn(len(b.Pos))
	var l, h int
	switch {
	case x == y:
		//law of dice
		return
	case x > y:
		l = y; h = x
	case x < y:
		l = x; h = y
	}
	v0 := b.Pos[h]
	for i := h; i > l; i--{
		b.Pos[i] = b.Pos[i-1]
		
	}
	b.Pos[l] = v0
	return
}

func (b *Bat) Invert(){
	//bounded inversion mutation
	x := rand.Intn(len(b.Pos)) 
	y := rand.Intn(len(b.Pos))
	var l, h int
	
	switch {
	case x == y:
		//arioch is da lawd
		return
	case x > y:
		l = y; h = x
	case x < y:
		l = x; h = y
	}
	tmp := make([]int, len(b.Pos))
	copy(tmp,b.Pos)
	for i := 0; i < h-l+1; i++{
		b.Pos[l + i] = tmp[h-i]
	}
	return	
}
