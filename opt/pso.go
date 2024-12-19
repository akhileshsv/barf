package barf

import (
	"fmt"
	"math/rand"
	"math"
	"time"
	draw"barf/draw"
	kass"barf/kass"
)


type Prat struct{
	Id    int
	Pos   []float64
	Vel   []float64
	Bpos  []float64
	Fit   float64
	Wt    float64
	Cons  int
	Typ   string
	Bfit  float64
	Lbest float64
}


func (p *Prat) Init(nd int, mn,mx []float64) {
	p.Pos = make([]float64, nd)
	p.Vel = make([]float64, nd)
	p.Bpos = make([]float64, nd)
	for i := 0; i < nd; i++{
		p.Pos[i] = (mx[i] - mn[i]) * rand.Float64() + mn[i]
		p.Vel[i] = 2.0 * rand.Float64() - 1.0
	}
	p.Bfit = -6.66
}

func (p *Prat) Eval(f func([]float64, []interface{}) float64, inp []interface{}){
	p.Fit = f(p.Pos, inp)
	if p.Bfit > p.Fit || p.Bfit == -6.66{
		p.Bfit = p.Fit
		copy(p.Bpos, p.Pos)
	}
}

func pinit(p *Prat, nd int, mn, mx []float64, pchn chan int){
	//parallel init
	p.Pos = make([]float64, nd)
	p.Vel = make([]float64, nd)
	p.Bpos = make([]float64, nd)
	for i := 0; i < nd; i++{
		p.Pos[i] = (mx[i] - mn[i]) * rand.Float64() + mn[i]
		p.Vel[i] = 2.0 * rand.Float64() - 1.0
	}
	p.Bfit = -6.66
	pchn <- 1
}

func pfit(i int, p *Prat, f func([]float64, []interface{}) float64, inp []interface{}, pchn chan []interface{}){
	//parallel fitness eval
	rez := make([]interface{}, 2)
	p.Fit = f(p.Pos, inp)
	if p.Bfit > p.Fit || p.Bfit == -6.66{
		p.Bfit = p.Fit
		copy(p.Bpos, p.Pos)
	}
	rez[0] = i
	rez[1] = p.Fit
	pchn <- rez
}

func pupdate(p *Prat, nd int, w, c1, c2 float64, gpos, mx, mn []float64, pchn chan int){
	for i := 0; i < nd; i++{
		//p.Vel[i] = w * p.Vel[i] + c1 * rand.Float64() * (p.Bpos[i] - p.Pos[i]) + c2 * rand.Float64() * (gpos[i] - p.Pos[i])
		
		p.Vel[i] = w * p.Vel[i] + c1 * rand.Float64() * (p.Bpos[i] - p.Pos[i])/(mx[i] - mn[i]) + c2 * rand.Float64() * (gpos[i] - p.Pos[i])/(mx[i] - mn[i])
		p.Pos[i] = p.Pos[i] + p.Vel[i]
		if p.Pos[i] > mx[i]{
			//just stay - ye olde
			p.Pos[i] = mx[i] 
		}
		if p.Pos[i] < mn[i] {
			//stay
			p.Pos[i] = mn[i] 
		}
	}
	pchn <- 1
}

func (p *Prat) Update(nd int, w, c1, c2 float64, gpos, mx, mn []float64){
	for i := 0; i < nd; i++{
		p.Vel[i] = w * p.Vel[i] + c1 * rand.Float64() * (p.Bpos[i] - p.Pos[i]) + c2 * rand.Float64() * (gpos[i] - p.Pos[i])
		p.Pos[i] = p.Pos[i] + p.Vel[i]
		if p.Pos[i] > mx[i] {
			//just stay - ye olde
			p.Pos[i] = mx[i]
		}
		if p.Pos[i] < mn[i] {
			//stay
			p.Pos[i] = mn[i]
		}
	}
}

func (p *Prat) Draw() (s string){
	for _, v := range p.Pos{
		s += fmt.Sprintf("%f ",v)
	}
	s += fmt.Sprintf("%f ",p.Fit)
	s += "\n"
	return
}

//gpos, gw, gb, pltstr, err
//w, c1, c2 - pso constants
//np - no. of particles, ng - no of gens/iter, nd - ndims/nparams (len(vec))
//mx, mn - []max and min vals, f - fitness func
//check out UPSO - add local best to modify velocity
func Psoloop(web, par, impr bool, w, c1, c2 float64, np, ng, nd int, mx,mn []float64, f func([]float64, []interface{}) float64, inp []interface{}, drw, trm, title string)(gpos []float64, gb float64, pltstr string){
	rand.Seed(time.Now().UnixNano())	
	var fdat, dat string
	swrm := make([]Prat, np)
	gpos = make([]float64, nd)
	gb = -6.66
	//adding here imp. criteria via li-peng-zhou 2013 paper (bool impr)
	//improvement of inertia w = wini - wend * (gmax - ng)/gmax + wend
	//there wi = 0.9; we - 0.2
	//improvement of acceleration c1 and c2
	var win, wend, wmin, wmax float64
	win = 0.5
	wend = 0.1
	wmin = 0.4
	wmax = 0.9
	cc := 2.0 //c = 2.0 (between 2.0 and 3.0)
	wavg := (wmin + wmax)/2.0
	switch par{
		case true:
		for i := 0; i < np; i++{
			swrm[i].Init(nd, mx, mn)
		}
		pchn := make(chan []interface{},np)
		for i := 0; i < np; i++{
			pfit(i, &swrm[i], f, inp, pchn)
		}
		for i := 0; i < np; i++{
			rez := <- pchn
			idx, _ := rez[0].(int)
			fit, _ := rez[1].(float64)
			if fit < gb || gb == -6.66{
				gb = fit
				copy(gpos, swrm[idx].Pos)
			}
			if drw == "all"{
				fdat += swrm[idx].Draw()
			}
		}
		case false:
		for i := 0; i < np; i++{		
			swrm[i].Init(nd, mx, mn)
			swrm[i].Eval(f, inp)
			if swrm[i].Fit < gb || gb == -6.66{
				gb = swrm[i].Fit
				copy(gpos, swrm[i].Pos)
			}
			if drw == "all"{
				fdat += swrm[i].Draw()
			}
		}
	}
	for gen := 0; gen < ng; gen++{
		switch drw{
			case "all":
			if !web{fmt.Println(ColorBlue, "gen->",gen,ColorWhite,"\nglobal best->\n",gpos, ColorGreen,"\nmin fitness->", gb,ColorReset)}
			default:
			if gen % 10 == 0{
				if !web{fmt.Println(ColorBlue, "gen->",gen,ColorWhite, ColorGreen,"\tmin fitness->", gb,ColorReset)}
			}
		}
		gf := (float64(ng) - float64(gen))/float64(ng)
		if gen !=0 && impr{w = (win - wend) * math.Pow(gf, 2.0) + wend}
		if impr && gen > 0{	
			switch {
			case wend < w && w <= wmin:
				c1 = cc * (wmin - w)/(wmin- wend)
			case wmin < w && w <= wavg:
				c1 = cc * ((wmin - w)/(wmin - wend) + 2.0 * (w - wmax)/(wmax - wmin) + (w - wmin)/(win - wmin))
			case wavg < w && w <= wmax:
				
				c1 = cc * ((wmax - w)/(wmax - wend) + 2.0 * (wmax - w)/(wmax - wmin) + (w - wmin)/(win - wmin))
			case wmax < w && w < win:
				c1 = cc * (w - wmax)/(win - wmax)
			}
			c2 =c1
		}
		//if gen > 0 && impr{c2 = c1}
		//*/
		switch par{
			case true:
			for i := 0; i < np; i++{
				swrm[i].Update(nd, w, c1, c2, gpos, mx, mn)
			}
			pchn := make(chan []interface{},np)
			for i := 0; i < np; i++{
				pfit(i, &swrm[i], f, inp, pchn)
			}
			for i := 0; i < np; i++{
				rez := <- pchn
				idx, _ := rez[0].(int)
				fit, _ := rez[1].(float64)
				if fit < gb || gb == -6.66{
					gb = fit
					copy(gpos, swrm[idx].Pos)
				}				
				if drw == "all"{
					fdat += swrm[idx].Draw()
				}
			
			}
			case false:
			for i := 0; i < np; i++{
				swrm[i].Update(nd, w, c1, c2, gpos, mx, mn)
				swrm[i].Eval(f, inp)	
				if swrm[i].Fit < gb || gb == -6.66{
					gb = swrm[i].Fit
					copy(gpos, swrm[i].Pos)
				}
				if drw == "all"{
					fdat += swrm[i].Draw()
				}
			}
		}
		dat += fmt.Sprintf("%v %f\n",gen, gb)
	}
	skript := "d2.gp"
	var folder string
	if web{folder = "web"}
	if trm == "dxf"{trm = "svg"}
	switch drw{
		case "gen":
		pltstr, _ = draw.Draw(dat, skript, trm, folder,title, title,"gen", "gbest","")		
		// if trm != "qt" && !web{fmt.Println(pltstr)}
		case "all":
		skript = "d3.gp"
		//fdat, skript, trm, folder, title, title, "gen","fit",""
		//fdat, skript, trm, folder, title, title, "gen","fit",""
		pltstr, _ = draw.Draw(fdat, skript, trm, folder,title,title, "gen", "gbest","")
		if trm != "qt" && !web{fmt.Println(pltstr)}
		skript = "d2.gp"
		pltstr, _ = draw.Draw(dat, skript, trm, folder,title, title,"gen", "gbest","")
		if trm != "qt" && !web{fmt.Println(pltstr)}
	}
	if !web{fmt.Println(ColorCyan, "best pos->\n", gpos, ColorRed, "\nmin fitness->", gb, ColorReset)}
	//fmt.Println("PLAATSTR-",pltstr)
	return 
}

func dwave(pos []float64) (float64){
	//drop wave function
	n := 1. + math.Cos(12.*math.Sqrt(math.Pow(pos[0], 2)+math.Pow(pos[1], 2)))
        d := 0.5*(math.Pow(pos[0], 2)+math.Pow(pos[1], 2)) + 2
	return -n/d
}

func rasta(pos []float64) (fit float64){
	//rastrigin function
	for _, x := range pos{
		fit += math.Pow(x, 2) - 10. * math.Cos(2. * math.Pi * x) + 10.0
	}
	return fit
}

func sphere(pos []float64) (fit float64){
	//sphere function
	for _, x := range pos{
		fit += math.Pow(x, 2)
	}
	return fit
}

func stang(pos []float64) (fit float64){
	for _, x := range pos{
		fit += math.Pow(x, 4) - 16*math.Pow(x, 2) + 5*x
	}
	fit = 0.5 * fit
	return
}


func dwavepso(pos []float64, inp []interface{}) (float64){
	//drop wave function
	n := 1. + math.Cos(12.*math.Sqrt(math.Pow(pos[0], 2)+math.Pow(pos[1], 2)))
        d := 0.5*(math.Pow(pos[0], 2)+math.Pow(pos[1], 2)) + 2
	return -n/d
}

func rastapso(pos []float64, inp []interface{}) (fit float64){
	//rastrigin function
	for _, x := range pos{
		fit += math.Pow(x, 2) - 10. * math.Cos(2. * math.Pi * x) + 10.0
	}
	return fit
}

func spherepso(pos []float64, inp []interface{}) (fit float64){
	//sphere function
	for _, x := range pos{
		fit += math.Pow(x, 2)
	}
	return fit
}

func stangpso(pos []float64, inp []interface{}) (fit float64){
	for _, x := range pos{
		fit += math.Pow(x, 4) - 16*math.Pow(x, 2) + 5*x
	}
	fit = 0.5 * fit
	return
}

func trsrakapso(pos []float64, inp []interface{}) (fit float64){
	mod, _ := inp[0].(kass.Model)
	secs, _ := inp[1].([][]float64)
	pmax,_ := inp[2].(float64)
	dmax,_ := inp[3].(float64)
	dens,_ := inp[4].(float64)
	cp := make([][]float64, len(pos))
	for i, idx := range pos{
		cp[i] = make([]float64,1)
		if len(secs) == 1{
			cp[i][0] = secs[0][int(idx)]
		} else {
		//	//WRENG
			cp[i][0] = secs[i][int(idx)]
		}
		
	}
	mod.Cp = cp
	var wt, C, gx, con float64
	frmrez, err := kass.CalcTrs(&mod, mod.Ncjt)
	if err != nil {
		fit = 1e6
		return
	}
	js, _ := frmrez[0].(map[int]*kass.Node)
	ms, _ := frmrez[1].(map[int]*kass.Mem)
	for _, node := range js{
		for _, d := range node.Displ {
			gx = math.Abs(d)/dmax - 1.0
			if gx > 0.0 {
				C += gx
				con += 1.0
			}
		}
	}
	for _, mem := range ms{
		wt += mem.Geoms[0] * mem.Geoms[2] * dens
		pmem := mem.Qf[0] / mem.Geoms[2]
		gx = math.Abs(pmem)/pmax - 1.0
		if gx > 0.0 {
			C += gx
			con += 1.0
		}
	}
	wt = wt*(1.0 + 10.0*C)
	fit = wt
	return
}
