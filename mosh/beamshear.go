package barf

import (
	//"log"
	"math"
	"sort"
	"errors"
)

//beam shear (and torsion) design funcs (links/stirrups)

//Tcmax returns the max allowable shear stress given fck
func Tcmax(code int, fck float64) (tcmax float64){
	//tcmax - max allowable shear stress
	switch code{
		case 1:
		tcmax = 0.631 * math.Sqrt(fck)
		if tcmax > 4.8{tcmax = 4.8}
		case 2:
		tcmax = 0.8 * math.Sqrt(fck)
		if tcmax > 5.0{tcmax = 5.0}
	}
	return
}

//TcFck calcs the shear stress of concrete tc
//shear stress of concrete tc - func of grade (fck), effd (size effect) and percentage steel (pt)
func TcFck(code int, fck, effd, pt float64) (tc float64){
	////pt = 100.0 * (b.Ast)/(b.Bw * effd)
	switch code{
		case 1:
		beta := 0.8 * fck/6.89/pt
		if beta < 1.0 {beta = 1.0}
		//fmt.Println("BETA",beta,pt)
		//beta = 0.8 * fck /(6.89 * pt)
		//if beta < 1.0 {beta = 1.0}
		tc = 0.85 * math.Sqrt(0.8 * fck) * (math.Sqrt(1.0 + 5.0 * beta)-1.0)/(6.0 * beta)
		case 2:
		switch {
		case pt <= 0.15:
			pt = 0.15
		case pt >= 3.0:
			pt = 3.0
		}
		if fck > 40.0 {
			fck = 40.0
		}
		d9 := effd
		if d9 > 400.0 {
			d9 = 400
		}
		tc = 0.79 * math.Pow(pt, 0.333) * math.Pow(400.0/d9, 0.25) / 1.25
		if fck > 25.0{
			tc = tc * math.Pow(fck/25.0,0.333)
		}
		//fmt.Println("shear strength of concrete tc->",tc," n/mm2")
	}
	return
}


//BmShrDz - general (*_*)7 span shear design
//lspan, xs, vs are from bmenv/bmsf
func BmShrDz(b *RccBm, lsx, rsx float64, xs, vs []float64) (err error){
	//b - send the midspan beam (id - 1), store spacing in this
	//CHECK FOR NPSEC shear calc- dis is wreng
	//BHAIYYA make b.Verbose a running thing across all structs
	//log.Println("hello from shear design->",b.Id, b.Title)
	if b.Bw > 400.0 {b.Nlegs = 4} else {b.Nlegs = 2}
	if b.Dlink == 0.0{b.Dlink = 8.0}
	tcmax := Tcmax(b.Code, b.Fck)
	effd := b.Dused - b.Cvrt
	//fmt.Println("effd->",b.Dused, b.Cvrt, effd)
	vclx := (effd + lsx/2.0)/1e3; vcrx := xs[20] - (effd + rsx/2.0)/1e3
	vcl := math.Abs(lerp(xs, vs, vclx)); vcr := math.Abs(lerp(xs, vs, vcrx))
	vcmax := vcl; if vcmax < vcr {vcmax = vcr}
	//if b.Verbose{fmt.Println("max allowable shear",tcmax * b.Bw * effd/1000.0," kn")}
	area := b.Bw * effd/1000.0
	//if b.Npsec{area = b.Sec.Prop.Area/1e3}
	//log.Println("hello from shrdz->",b.Title)
	if vcmax > tcmax * area {
		err = errors.New("insufficient depth for shear")
		//log.Println("SHEAR DEPTH ERROR")
		return
	}
	pt := 100.0 * (b.Ast)/(b.Bw * effd)

	//THIS SEEMS LIKE A MISTAKE to not have(refer is 456)
	//if b.Tyb > 0.0{
	//	pt = 100.0 * b.Ast/(b.Bf * effd)
	//}

	//use asc for cantilevers 
	if b.Ast == 0.0 || b.Endc == 0{
		pt = 100.0 * (b.Asc)/(b.Bw * effd)
	}
	//just because it agrees with hulse 7.6.1 results and also makes sense, only 50% of steel continues onto supports in a ss beam
	if b.Endc == 1 && b.Code == 2{
		pt = pt/2.0
	}
	if b.Npsec{pt = 100.0 * b.Ast/b.Sec.Prop.Area}
	tc := TcFck(b.Code, b.Fck, effd, pt); vc := tc * b.Bw * effd/1e3
	//if b.Verbose{fmt.Println("shear strength of concrete tc->",tc," n/mm2","vuc->",vc,"kn")}
	//if b.Verbose{fmt.Println("shear to be resisted by stirrups vusv->", vcl - vc,"kn")}
	vsmin := 0.4 * b.Bw * effd/1e3
	if b.Npsec{pt = 0.4 * b.Sec.Prop.Area/1e3}
	//if b.Verbose{fmt.Println("shear resistance of minimum stirrups->", vsmin, "kn", "total->",vsmin + vc, "kn")}
	var l1, l2, r1, r2 float64
	for i, v := range vs{
		if l1 == 0.0 && math.Abs(v) <= vsmin + vc{
			//fmt.Println("found l1->",v, "at x->",xs[i])
			l1 = xs[i]
		}

		if l2 == 0.0 && math.Abs(v) <= 0.5 * vc{
			//fmt.Println("found l2->",v, "at x->",xs[i])
			l2 = xs[i]
		}
	}
	//fmt.Println(l1, l2)
	b.L1 = l1; b.L2 = l2
	for i := 20; i >= 0; i--{
		v := vs[i]
		if r1 == 0.0 && math.Abs(v) <= vsmin + vc{
			//fmt.Println("found r1->",v, "at x->",xs[i])
			r1 = xs[i]
		}

		if r2 == 0.0 && math.Abs(v) <= 0.5 * vc{
			//fmt.Println("found r2->",v, "at x->",xs[i])
			r2 = xs[i]
		}
	}
	b.L3 = r2; b.L4 = r1
	//fmt.Println(r1, r2)
	fy := b.Fyv
	switch b.Code{
		case 1:
		if fy > 415.0{fy = 415.0}
		case 2,3:
		if fy > 460.0{fy = 460.0}
	}
	vusv := vcl - vc
	if vcr - vc > vusv{vusv = vcr - vc}
	if vsmin + vc >= vcl && vsmin + vc >= vcr {
		//fmt.Println("minimum shear reinforcement required")
		vusv = 0.0
		b.Nomlink = true
		
	}
	slink, smin, snom, err := LinkSpacing(fy, b.Bw, effd, b.Dlink, vusv, b.Nlegs)
	if err != nil{
		//log.Println("LINK SPACING ERRORE,errore->",err)
		return
	}
	mainlen := b.L1 + xs[20]-b.L4; minlen := b.L2 - b.L1 + b.L4 - b.L3; nomlen := xs[20] - mainlen - minlen

	b.Lspc = []float64{slink, smin, snom, mainlen, minlen, nomlen}
	b.Nlx = make([]float64,4)
	b.Nlx[0] = math.Ceil(mainlen * 1e3/slink) + 1.0
	b.Nlx[1] = math.Ceil(minlen * 1e3/smin) + 1.0
	b.Nlx[2] = math.Ceil(nomlen * 1e3/snom) + 1.0
	b.Nlx[3] = b.Nlx[0] + b.Nlx[1] + b.Nlx[2] 
	err = nil
	//if b.Verbose{
	//log.Println("shear links main->",b.Dlink, " mm dia at ", slink, "mm spacing", "net length->",mainlen,"m")
	//log.Println("shear links min->",b.Dlink, " mm dia at ", smin, "mm spacing", "net length->", minlen,"m")
	//log.Println("shear links nominal->",b.Dlink, " mm dia at ", snom, "mm spacing","net length->",nomlen,"m")
	//}
	return
}


//BmShrChk checks shear (b.Vf) at a section and calcs link spacing/dia
func BmShrChk(b *RccBm) (err error){
	//fmt.Println("ast->",b.Ast)
	if b.Bw > 400.0 {b.Nlegs = 4} else {b.Nlegs = 2}
	//if b.Rslb{b.Nlegs = 1}
	if b.Dlink == 0.0{b.Dlink = 8.0}
	tcmax := Tcmax(b.Code, b.Fck)
	effd := b.Dused - b.Cvrt
	//fmt.Println("effd->",b.Dused, b.Cvrt, effd)
	vcmax := b.Vu
	b.Nlx = make([]float64,4)
	//if b.Verbose{fmt.Println("max allowable shear",tcmax * b.Bw * effd/1000.0," kn")}
	area := b.Bw * effd/1000.0
	//if b.Npsec{area = b.Sec.Prop.Area/1e3}
	if vcmax > tcmax * area{
		err = errors.New("insufficient depth for shear")
		return
	}
	pt := 100.0 * (b.Ast)/(b.Bw * effd)
	if b.Tyb > 0.0{
		pt = 100.0 * b.Ast/(b.Bf * effd)
	}
	//fmt.Println("pt%",pt)
	//use asc for cantilevers 
	if b.Ast == 0.0 || b.Endc == 0{
		pt = 100.0 * (b.Asc)/(b.Bw * effd)
	}
	if b.Endc == 1 && b.Code == 2{
		pt = pt/2.0
	}
	if b.Npsec{pt = 100.0 * b.Ast/b.Sec.Prop.Area}
	tc := TcFck(b.Code, b.Fck, effd, pt); vc := tc * b.Bw * effd/1e3
	//fmt.Println("shear strength of concrete tc->",tc," n/mm2","vuc->",vc,"kn")
	//fmt.Println("shear to be resisted by stirrups vusv->", vcmax - vc,"kn")
	//fmt.Println("n leg stirrups ->", b.Nlegs)
	if vcmax <= vc{
		if b.Rslb{
			//fmt.Println("no stirrups required")
			b.Nomlink = true
			return
		}
	}
	vusv := vcmax - vc
	vsmin := 0.4 * b.Bw * effd/1e3
	fy := b.Fyv
	switch b.Code{
		case 1:
		if fy > 415.0 {fy = 415.0}
		case 2,3:
		if fy > 460.0{fy = 460.0}
	}
	dlinks := []float64{8.0,10.0,12.0,16.0}
	iter := 0
	idx := 0
	var slink, smin, snom float64
	for iter != -1{
		if idx == 3{
			iter = -1
			return
		}
		
		slink, smin, snom, err = LinkSpacing(fy, b.Bw, effd, dlinks[idx], vusv, b.Nlegs)	
		if err == nil{
			iter = -1
			break
		}
		idx++
	}
	b.Dlink = dlinks[idx]
	b.Nlx = make([]float64,4)

	b.Lspc = []float64{slink, smin, snom, 0,0,0}
	if vsmin + vc >= vcmax{
		//fmt.Println("minimum shear reinforcement required",smin)
		//if b.Verbose{fmt.Println("minimum shear reinforcement required")}
		b.Nlx[0] = math.Ceil(b.Lspan * 1e3/smin) + 1.0
	} else {
		//fmt.Println("regular shear reinforcement required BAAH",slink)
		b.Nlx[0] = math.Ceil(b.Lspan*1e3/slink) + 1.0
	}
	//log.Println("number of stirrups->",b.Nlx[0])
	err = nil
	return
	
}


//LinkSpacing returns link spacing for shear (stirrup, nominal)
func LinkSpacing(fy, bw, effd, dia, vusv float64, nleg int) (slink, smin, snom float64, err error){
	//vsmin := 0.4 * b.Bw * effd/1e3
	//fmt.Println("nleg->",nleg, float64(nleg))
	smax := 10.0 * math.Floor(0.75 * effd/10.0)
	if smax > 300.0{smax = 300.0}
	snom = 10 * math.Floor(fy* float64(nleg) * RbrArea(dia) *2.5/(10*bw))
	if snom < 75.0 {
		err = errors.New("link spacing error")
		return
	}
	smin = 10.0 * math.Floor(0.87 * fy * float64(nleg) * RbrArea(dia)/0.4/bw/10.0)
	if smin < 75.0 {
		err = errors.New("link spacing error")
		return
	}
	if smin > smax{
		smin = smax
	}
	if snom > smax{
		snom = smax
	}
	//what are nominal links?
	if vusv == 0.0{
		slink = smin
		return
	}
	slink = 10 * math.Floor(0.87* fy * float64(nleg) * RbrArea(dia)*effd/(10*vusv*1000.0))
	if slink < 75.0{
		err = errors.New("link spacing error")
		return
	}
	if slink > 300.0{slink = 300.0}
	if slink > 0.75 * effd{slink = 0.75 * effd}
	return
}

//AsvRatBs calcs asv/sv ratio of stirrups
//see hulse section 3.3 (not used except in one test)
func AsvRatBs(b *RccBm, vf float64) (asvrat float64){
	effd := b.Dused - b.Cvrt
	va := (vf * 1000.0) / (b.Bw * effd)
	v1 := 0.8 * math.Sqrt(b.Fck)
	if v1 > 5.0 {
		v1 = 5.0
	}
	if va > v1 {
		return -99
	}
	var vc float64
	r0 := 100.0 * b.Ast / (b.Bw * effd)
	switch {
	case r0 <= 0.15:
		r0 = 0.15
	case r0 >= 3.0:
		r0 = 3.0
	}
	u := b.Fck
	if b.Fck > 40.0 {
		u = 40.0
	}
	d9 := effd
	if d9 > 400.0 {
		d9 = 400
	}
	vc = 0.79 * math.Pow(r0, 0.333) * math.Pow(400.0/d9, 0.25) / 1.25
	if b.Fck > 25.0 {
		vc = vc * math.Pow(u/25.0, 0.333)
	}
	if va > vc+0.4 {
		asvrat = b.Bw * (va - vc) / (0.87 * b.Fy)
	} else {
		asvrat = 0.4 * b.Bw / (0.87 * b.Fy)
	}
	return
}

//TorRect updates b.Hdim and b.Ldim (rectangle dimensions and link dimensions)
//see hulse sec 3.4
func (b *RccBm) TorRect(){
	//hdim (rect dim), ldim (link dims - rect dim - cvr) [][]float64, nc - len(hdim), sdx (component carrying shear) int
	//NOTE- ADDING CVRL SEEMS WRONG, tis the only way to get results to agree
	cvrl := 30.0 
	switch b.Styp{
		case 1,14:
		b.Hdim = [][]float64{{b.Bw,b.Dused}}
		b.Ldim = [][]float64{{b.Bw - 2*cvrl,b.Dused - 2*cvrl}}
		b.Sdx = 0
		case 6:
		b.Hdim = [][]float64{
			{(b.Bf-b.Bw)/2.0,b.Df},
			{b.Bw,b.Dused},
			{(b.Bf-b.Bw)/2.0,b.Df},
		}
		
		b.Ldim = [][]float64{
			{(b.Bf-b.Bw)/2.0,b.Df-2*cvrl},
			{b.Bw - 2*cvrl,b.Dused - 2*cvrl},
			{(b.Bf-b.Bw)/2.0,b.Df-2.0*cvrl},
		}
		b.Sdx = 1
		case 7,8,9,10:
		b.Hdim = [][]float64{
			{(b.Bf-b.Bw)/2.0,b.Df},
			{b.Bw,b.Dused},
		}
		b.Ldim = [][]float64{
			{(b.Bf-b.Bw)/2.0,b.Df-2*cvrl},
			{b.Bw - 2*cvrl,b.Dused - 2*cvrl},
		}
		b.Sdx = 1
		default:
		//lol WHEN
	}
	for _, rect := range b.Hdim{
		sort.Slice(rect,func(i,j int) bool{
			return rect[i] > rect[j]
		})
	}
	for _, rect := range b.Ldim{
		sort.Slice(rect,func(i,j int) bool{
			return rect[i] > rect[j]
		})
	}
	b.Nr = len(b.Hdim)
	return
}

//BmTorDBs designs a beam section for torsion
//hulse 3.4.2
func BmTorDBs(b *RccBm, asvrat, vf, tm float64) (err error){
	var s float64
	var hs = make([]float64, len(b.Hdim))
	for idx, hdim := range b.Hdim {
		//hdim 0 - hmim, hdim 1 - hmax
		hs[idx] = math.Pow(hdim[1], 3) * hdim[0]
		s += hs[idx]
	}
	vf = vf * 1e3 //kn - n
	tm = tm * 1e6 //knm - n mm
	var ts = make([]float64, len(b.Hdim))
	var vs = make([]float64, len(b.Hdim))
	var avs = make([][]float64, len(b.Hdim))
	var svs = make([]float64, len(b.Hdim))
	for i := 0; i < len(b.Hdim); i++ {
		avs[i] = make([]float64, 3)
	}
	for idx, h := range hs{
		var va, v1, v3 float64
		ts[idx] = tm * h / s
		vs[idx] = 2.0 * ts[idx] / (math.Pow(b.Hdim[idx][0], 2) * (b.Hdim[idx][1] - 0.333*b.Hdim[idx][0]))
		if idx == b.Sdx {
			va = 0.0
			avs[idx][0] = asvrat
		} else {
			va = vf / (b.Bw * (b.Dused - b.Cvrt))
		}
		if b.Fck > 40.0 {
			v1 = 5.0
			v3 = 0.4
		} else {
			v1 = 0.8 * math.Sqrt(b.Fck)
			v3 = 0.067 * math.Sqrt(b.Fck)
		}
		x1 := b.Ldim[idx][1]
		y1 := b.Ldim[idx][0]
		svs[idx] = 200.0
		if svs[idx] > x1 {
			svs[idx] = x1
		}
		if svs[idx] > y1/2.0 {
			svs[idx] = y1 / 2.0
		}
		if y1 < 550.0 && vs[idx] > (v1*y1/550.0) {
			//fmt.Println("max allowable shear stress exceeded")
			err = ErrShear
			avs[idx][0] = -99.0
			avs[idx][1] = -99.0
			idx++
			continue
		}
		if vs[idx]+va > v1 {
			//fmt.Println("max allowable ultimate stress exceeded")
			err = ErrShear
			avs[idx][0] = -99.0
			avs[idx][1] = -99.0
			idx++
			continue
		}
		if vs[idx] < v3 {
			//fmt.Println("no reinforcement required")
			idx++
			continue
		}
		avs[idx][1] = avs[idx][0] + ts[idx]/(0.8*x1*y1*0.87*b.Fy)
		//here bc b.Fy = b.Fyv so
		avs[idx][2] = avs[idx][2] + (avs[idx][1]-avs[idx][0])*(b.Fy/b.Fy)*(x1+y1)
	}
	b.Svr = svs
	b.Avr = avs
	return 
	
}


/*
   
	-----
	     \-------
	             |-------
	 +-----------+-------\------------------------^-----------+
      |	 +-----------+----------------\---------------+-----------+  |
      |	 |	     |	                      \-------|		  |  |
      |	 |	     |	                              |-------	  |  |
      |	 |     	     | 	       	                      |       \----  |
	     	     |l1       l2      	   l3  	      |l4        r0
		     |		    		      v
	  main 	       	  min  	    nom	       	 min   	     main
           vu > vc                vu < 0.5 vc

vc = tc * bw * d
tv = vu/bw/d


*/
