package barf

import (
	//"fmt"
	"log"
	"math"
)

//RccFtng is an rcc (Rcc) footing (Ftng) field struct
//uses meters here for dimensions instead of mm 
type RccFtng struct {
	Colx, Coly                  float64
	Fck, Fy                     float64
	Df, Dmin                    float64
	Eo, Sbc                     float64
	Pgck, Pgsoil                float64
	Nomcvr                      float64
	Based                       float64
	Pus                         []float64
	Mxs                         []float64
	Mys                         []float64
	Psfs                        []float64
	Dused                       float64
	Effd                        float64
	Efcvr                       float64
	D1, D2                      float64
	Typ                         int
	Code                        int
	Shape                       string
	Selfwt                      bool
	Sloped                      bool
	Dlfac                       bool
	Verbose                     bool
	Dz                          bool
	Dsgn                        bool
	Web                         bool
	Brchk                       bool
	Term                        string
	Title                       string
	Astx, Asty                  float64                 `json:",omitempty"`
	Lx, Ly                      float64                 `json:",omitempty"`
	Rbropt                      [][]float64             `json:",omitempty"`
	Rbr                         []float64               `json:",omitempty"`
	Basecy                      float64                 `json:",omitempty"` //length of base in kompress
	Qmax, Qmin                  float64                 `json:",omitempty"`
	Mur, Vur                    float64                 `json:",omitempty"`
	Vp, Vpp                     float64                 `json:",omitempty"`
	Mux, Muy                    float64                 `json:",omitempty"`
	Vux, Vuy                    float64                 `json:",omitempty"`
	Vtot, Vrcc, Wstl, Afw, Kost float64                 `json:",omitempty"`
	Kunit                       float64                 `json:",omitempty"`
	Bmap                        map[float64][][]float64 `json:",omitempty"`
	Bsum                        map[float64]float64     `json:",omitempty"`
	Cvec                        []float64               `json:",omitempty"` //cost ck, fw, steel
	Report                      string                  `json:",omitempty"`
	Kostin                      []float64               `json:",omitempty"`
	Txtplots                    []string                `json:",omitempty"`
}

func FtngDesign(f *RccFtng) {}
//OKAY LET THIS STAYZ FOR NOW

//fsize returns a footing base size (lx, ly)
//one does not think this is used but who knows what lurks where
func fsize(typ int, pu, mx, my, smx, cx, cy float64)(lx, ly float64){
	//log.Println("smx-",smx)
	switch{
		case mx > 0.0 && my == 0:
		//single axis moment (sub ex 15.5)
		//(first try solve cubic) - nope error prone
		lx = math.Sqrt(2.0 * pu/smx)
		switch typ{
			case 0:
			ly = lx
			case 1:
			ly = lx - cx + cy
		}
		case mx > 0.0 && my > 0.0:
		//jarquio 1983
		ar := 2.0 * math.Pow(72.0*mx*my*smx,1.0/3.0)/smx
		switch typ{
			case 0:
			lx = math.Sqrt(ar); ly = lx
			case 1:
			ly = mx * math.Sqrt(ar)/my
			lx = my*ly/mx
		}
	}
	return
}

//FtngDz is an entry func for footing design
//just to propagate the noble house of Dz 
func FtngDz(f *RccFtng) (err error) {
	pus := make([]float64, 3)
	mxs := make([]float64, 3)
	mys := make([]float64, 3)
	psfs := make([]float64, 3)
	for i := range psfs{
		psfs[i] = 1.0
	}
	for i, pu := range f.Pus{
		pus[i] = pu
		if len(f.Mxs) > i{
			mxs[i] = f.Mxs[i]
		}
		if len(f.Mys) > i{
			mys[i] = f.Mys[i]
		}
		if len(f.Psfs) > i{
			psfs[i] = f.Psfs[i]
		} 
	}
	f.Pus = make([]float64, 3)
	f.Mxs = make([]float64, 3)
	f.Mys = make([]float64, 3)
	f.Psfs = make([]float64, 3)
	copy(f.Pus, pus)
	copy(f.Mxs, mxs)
	copy(f.Mys, mys)
	copy(f.Psfs, psfs)
	err = FtngDzRojas(f)
	return
}

//FtngDzRojas rips off a.l. rojas - design of isolated rectangular footing
//designs an isolated rectangular footing (pad/sloped)
func FtngDzRojas(f *RccFtng) (err error){
	//a.l rojas - design of isloated rectangular footing
	//ALL DIMS IN METERS ALL DIMS IN METERS ALL DIMS IN METERS (sigh)
	
	var pu, pul, mx, my, mxl, myl, af, lx, ly, d, smx, efcvr, astx, asty, mux, muy, vux, vuy, vp float64
	for i := range f.Pus{
		pu += f.Pus[i]
		mx += f.Mxs[i]
		my += f.Mys[i]
		if f.Psfs[i] <= 0.0{
			switch len(f.Pus){
				case 1, 2:
				f.Psfs[i] = 1.5
				case 3:
				f.Psfs[i] = 1.2
			}
		}
		pul += f.Pus[i] * f.Psfs[i]
		mxl += f.Mxs[i] * f.Psfs[i]
		myl += f.Mys[i] * f.Psfs[i]
	}
	efcvr = f.Nomcvr + 0.020
	if f.D1 > 0.0 {efcvr = f.Nomcvr + f.D1}
	d = 0.075
	step := 0.01; iter := 1
	d -= step
	var kiter, mdx int
	var rez [][]float64
	switch f.Shape{
		case "rect":
		f.Typ = 1
		case "square":
		f.Typ = 0
	}
	if f.Pgck == 0.0{
		f.Pgck = 24.0
	}
	switch{
		case mx == 0 && my == 0:
		for iter != -1{
			kiter++
			if kiter > 666{
				//log.Println("ERRORE,errore->maximum iteration limit reached")
				err = ErrFtngDim
				return
			}
			d += step
			//if hf == 0.0 {hf = d + efcvr}
			if f.Sloped{
				//DAMMIT THIS NEEDS TO CHANGE
				smx = f.Sbc - f.Pgck * (d + efcvr) - (f.Df - d - efcvr) * f.Pgsoil
				if f.Eo == 0.0 {f.Eo = 0.025}
				if f.Df == 0.0 {smx = f.Sbc - f.Pgck * (d + efcvr)}
			}else{
				smx = f.Sbc - f.Pgck * (d + efcvr) - (f.Df - d - efcvr) * f.Pgsoil
				if f.Df == 0.0 {smx = f.Sbc - f.Pgck * (d + efcvr)}
			}
			af = pu/smx
			if f.Dlfac {af = 1.1 * pu/f.Sbc}
			switch f.Typ{
				case 0:
				//case "square":
				ly = math.Round(100.0*math.Sqrt(af))/100.0
				lx = ly
				case 1:
				//case "rect":
				ly = (f.Coly - f.Colx)/2.0 + math.Sqrt(math.Pow(f.Coly - f.Colx,2)/4.0 + af)
				ly = math.Round(ly*100.0)/100.0
				lx = ly - f.Coly + f.Colx
			}
			af = ly * lx
			wu := pul/af
			var vuc, d1, d2 float64
			if f.Sloped{
				//if no slope assume 1 in 3 BTW
				if f.Dmin == 0.0{f.Dmin = 0.15}
				var bx, by, k1, k2 float64
				bx = f.Colx + 2.0 * f.Eo; by = f.Coly + 2.0 * f.Eo
				//check depth for b.m
				mux = wu/24.0 * (2.0 * ly + f.Coly) * math.Pow(lx - f.Colx,2)		
				muy = wu/24.0 * (2.0 * lx + f.Colx) * math.Pow(ly - f.Coly,2)
				if f.Dlfac{
					mux = wu * lx * math.Pow((ly - f.Coly)/2,2)/2.0
					muy = wu * ly * math.Pow((lx - f.Colx)/2,2)/2.0
				}
				if f.Fy == 415 {
					k1 = 0.138; k2 = 0.025
				} else {k1 = 0.133; k2 = 0.023}
				d1 = math.Sqrt(mux/(k1 * bx * f.Fck + k2 * (lx - bx) * f.Fck)/1000.0)
				d2 = math.Sqrt(muy/(k1 * by * f.Fck + k2 * (ly - by) * f.Fck)/1000.0)
				if d2 > d1 {d1 = d2}
				if d1 > d {continue}
				//check for punching shear at d/2 (= dy) from COLUMN FACE (subramanian)
				lcx := f.Colx + d; lcy := f.Coly + d
				x1 := (lx - bx)/2.0; x2 := d/2.0 - f.Eo
				dcr := (d - f.Dmin) * x2/x1
				if dcr < 0 {dcr = 0}
				dcr = d - dcr
				pcr := 2.0 * (lcx + lcy)
				vp = wu * (lx * ly - lcx * lcy)
				ks := 0.5 + f.Colx/f.Coly
				if ks > 1.0 {ks = 1.0}
				tuc := 0.25 * math.Sqrt(f.Fck)*ks
				vuc = 1000.0 * tuc * pcr * dcr
				if vp > vuc {continue}
				//get area of steel in x and y
				dx := d + 0.008
				astx = 1e6 * (0.5 * f.Fck/f.Fy)*(1.0 - math.Sqrt(1.0 - 4.6* mux/f.Fck/bx/dx/dx/1000.0))*bx*dx
				asty = 1e6 * (0.5 * f.Fck/f.Fy)*(1.0 - math.Sqrt(1.0 - 4.6* muy/f.Fck/by/d/d/1000.0))*by*d
				//log.Println("asts moments",astx, asty, mux, muy, dx, dx, by, d)
				rez, mdx, err = FtngRbrDia(d, f.Fck, f.Fy, f.Colx, f.Coly, lx, ly, astx, asty, 16.0, f.Nomcvr)
				if err != nil{
					//log.Println("ERRORE, errore-> dev len check failed")
					continue
				}
				//check for one way shear in y 
				x1 = (ly - by)/2.0; x2 = d - f.Eo
				dcr = (d - f.Dmin) * x2/x1
				if dcr < 0.0{dcr = 0.0}
				dcr = d - dcr
				lcy = f.Coly + 2.0 * d 
				vuy = wu * ((ly - f.Coly)/2.0 - d) * (lx + lcy)/2.0
				acr := lcy * dcr
				pty := 100.0 * asty/acr/1e6
				beta := 0.8 * f.Fck/6.89/pty
				tuc = 0.85 * math.Sqrt(0.8 * f.Fck) * (math.Sqrt(1.0 + 5.0 * beta) - 1.0)/6.0/beta
				vuc = 1e3 * tuc * acr
				if vuc < vuy{
					continue
				}
				//check for one way shear in x
				x1 = (lx - bx)/2.0; x2 = dx - f.Eo
				dcr = (dx - f.Dmin) * x2/x1
				if dcr < 0.0{dcr = 0.0}
				dcr = dx - dcr
				lcx = f.Colx + 2.0 * dx
				vux = wu * ((lx - f.Colx)/2.0 - dx) * (ly + lcx)/2.0
				acr = lcx * dcr
				ptx := 100.0 * astx/acr/1e6
				beta = 0.8 * f.Fck/6.89/ptx
				tuc = 0.85 * math.Sqrt(0.8 * f.Fck) * (math.Sqrt(1.0 + 5.0 * beta) - 1.0)/6.0/beta
				vuc = 1e3 * tuc * acr
				if vuc < vux{
					continue
				}
				//check for bearing pressure
				a1 := (f.Colx + 4 * d) * (f.Coly + 4 * d)
				if a1 > lx * ly {a1 = lx * ly}
				arat := math.Sqrt(a1/f.Colx/f.Coly)
				if arat > 2.0{
					arat = 2.0
				}
				fbrf := 0.45 * f.Fck * arat
				//if fbrf > fbrc {fbrc = fbrf}
				sbr := 1e3 * fbrf * f.Colx * f.Coly
				f.Brchk = sbr > pul
				iter = -1
			} else {
				//check for bm
				mux = 0.5 * wu * lx * math.Pow((ly - f.Coly)/2.0,2)
				muy = 0.5 * wu * ly * math.Pow((lx - f.Colx)/2.0,2)
				var k1 float64
				if f.Fy == 415 {k1 = 0.138} else {k1 = 0.133}
				d1 = math.Sqrt(mux/k1/f.Fck/lx/1000.)
				d2 = math.Sqrt(muy/k1/f.Fck/ly/1000.)
				if d2 > d1 {d1 = d2}
				if d1 > d {continue}
				//check for punching shear 
				vp = pul * (lx * ly - (f.Coly + d) * (f.Colx + d))/lx/ly
				ks := 0.5 + f.Colx/f.Coly
				if ks > 1.0 {ks = 1.0}
				tuc := 0.25 * math.Sqrt(f.Fck)*ks
				vuc = 1000.0 * tuc * 2.0*(f.Colx + f.Coly + 2*d) * d
				if vp > vuc {continue}
				//get area of steel in x and y
				dx := d + 0.008
				astx = 1e6 * (0.5 * f.Fck/f.Fy)*(1.0 - math.Sqrt(1.0 - 4.6* mux/f.Fck/lx/dx/dx/1000.0))*lx*dx
				asty = 1e6 * (0.5 * f.Fck/f.Fy)*(1.0 - math.Sqrt(1.0 - 4.6* muy/f.Fck/ly/d/d/1000.0))*ly*d
				rez, mdx, err = FtngRbrDia(d, f.Fck, f.Fy, f.Colx, f.Coly, lx, ly, astx, asty, 16.0, f.Nomcvr)
				if err != nil{
					//log.Println("ERRORE, errore-> dev len check failed")
					continue
				}
				//check for one way shear about x
				vux = wu * lx * (ly/2.0 - f.Coly/2.0 - dx)
				ptx := 100.0 * astx/lx/1000./dx/1000.0
				beta := 0.8 * f.Fck/6.89/ptx
				tuc = 0.85 * math.Sqrt(0.8 * f.Fck) * (math.Sqrt(1.0 + 5.0 * beta) - 1.0)/6.0/beta
				vuc = 1e3 * tuc * lx * dx
				if vuc < vux{
					continue
				}
				vuy = wu * ly * (lx/2.0 - f.Colx/2.0 - d)
				pty := 100.0 * astx/ly/1000./d/1000.0
				beta = 0.8 * f.Fck/6.89/pty
				tuc = 0.85 * math.Sqrt(0.8 * f.Fck) * (math.Sqrt(1.0 + 5.0 * beta) - 1.0)/6.0/beta
				vuc = 1e3 * tuc * ly * d
				if vuc < vuy{
					continue
				}
				//check for bearing pressure
				//at column face WILL DEPEND ON COLUMN F.FCK IDIOT
				//fbrc := 0.45 * f.Fck
				a1 := (f.Colx + 4 * d) * (f.Coly + 4 * d)
				if a1 > lx * ly {a1 = lx * ly}
				arat := math.Sqrt(a1/f.Colx/f.Coly)
				if arat > 2.0{
					arat = 2
				}
				fbrf := 0.45 * f.Fck * arat
				//if fbrf > fbrc {fbrc = fbrf}
				sbr := 1e3 * fbrf * f.Colx * f.Coly
				f.Brchk = sbr > pul
				iter = -1
			}
		}
		default:
		for iter != -1{
			kiter++
			if kiter > 666{
				//log.Println("ERRORE,errore->maximum iteration limit reached")
				err = ErrFtngDim
				return
			}
			d += step
			smx = f.Sbc - f.Pgck * (d + efcvr) - f.Pgsoil * (f.Df - d - efcvr)
			if f.Df == 0.0{smx = f.Sbc - f.Pgck * (d + efcvr)}
			if mx > 0.0 && my > 0.0{
				lx, ly = fsize(f.Typ, pu, mx, my, smx, f.Colx, f.Coly)
			} else {
				f.Based = d + efcvr
				err = FtngPadAz(f)
				if err != nil{
					return
				}
				lx = f.Lx; ly = f.Ly
			}
			ly = math.Round(ly*100.0)/100.0
			lx = math.Round(lx*100.0)/100.0
			//log.Println(ColorCyan,ly,lx,ColorReset)
			var k1, d1, d2 float64
			if f.Fy == 415 {k1 = 0.138} else {k1 = 0.133}
			mux = (pul * math.Pow(ly, 2) + 2.0 * mxl * (2.0 * ly + f.Coly)) * math.Pow(ly - f.Coly, 2)/8.0/math.Pow(ly,3)
			muy = (pul * math.Pow(lx, 2) + 2.0 * myl * (2.0 * lx + f.Colx)) * math.Pow(lx - f.Colx, 2)/8.0/math.Pow(lx,3)
			d1 = math.Sqrt(mux/k1/f.Fck/lx/1000.)
			d2 = math.Sqrt(muy/k1/f.Fck/ly/1000.0)
			if d2 > d1 {d1 = d2}
			if d1 > d{continue}
			vp = pul * (lx * ly - (f.Colx + d)*(f.Coly + d))/lx/ly
			ks := 0.5 + f.Colx/f.Coly
			if ks > 1.0 {ks = 1.0}
			tuc := 0.25 * math.Sqrt(f.Fck)*ks
			vuc := 1000.0 * tuc * 2.0*(f.Colx + f.Coly + 2*d) * d
			if vp > vuc {continue}
			//get area of steel in x and y
			dx := d + 0.008
			astx = 1e6 * (0.5 * f.Fck/f.Fy)*(1.0 - math.Sqrt(1.0 - 4.6* mux/f.Fck/lx/dx/dx/1000.0))*lx*dx
			asty = 1e6 * (0.5 * f.Fck/f.Fy)*(1.0 - math.Sqrt(1.0 - 4.6* muy/f.Fck/ly/d/d/1000.0))*ly*d
			rez, mdx, err = FtngRbrDia(d, f.Fck, f.Fy, f.Colx, f.Coly, lx, ly, astx, asty, 16.0, f.Nomcvr)
			if err != nil{
				//log.Println(err,rez)
				//log.Println("ERRORE, errore-> dev len check failed")
				continue
			}
			//check for one way shear about x
			vux = (pul * math.Pow(lx, 2) + 3.0 * myl * (lx + f.Colx + 2.0 * dx))* (lx - f.Colx - 2.0 * dx)/2.0/math.Pow(lx,3)
			ptx := 100.0 * astx/lx/1000./dx/1000.0
			beta := 0.8 * f.Fck/6.89/ptx
			tuc = 0.85 * math.Sqrt(0.8 * f.Fck) * (math.Sqrt(1.0 + 5.0 * beta) - 1.0)/6.0/beta
			vuc = 1e3 * tuc * lx * dx
			if vuc < math.Abs(vux){
				continue
			}
			vuy = (pul * math.Pow(ly, 2) + 3.0 * mxl * (ly + f.Coly + 2.0 * d))* (ly - f.Coly - 2.0 * d)/2.0/math.Pow(ly,3)
			pty := 100.0 * asty/ly/1000./d/1000.0
			beta = 0.8 * f.Fck/6.89/pty
			tuc = 0.85 * math.Sqrt(0.8 * f.Fck) * (math.Sqrt(1.0 + 5.0 * beta) - 1.0)/6.0/beta
			vuc = 1e3 * tuc * ly * d
			if vuc < math.Abs(vuy){
				continue
			}
			//check for bearing pressure
			//at column face WILL DEPEND ON COLUMN F.FCK IDIOT
			//fbrc := 0.45 * f.Fck
			a1 := (f.Colx + 4 * d) * (f.Coly + 4 * d)
			if a1 > lx * ly {a1 = lx * ly}
			arat := math.Sqrt(a1/f.Colx/f.Coly)
			if arat > 2.0{
				arat = 2
			}
			fbrf := 0.45 * f.Fck * arat
			//if fbrf > fbrc {fbrc = fbrf}
			sbr := 1e3 * fbrf * f.Colx * f.Coly
			f.Brchk = sbr > pul
			iter = -1
		}
	}
	//log.Println("iter finito")
	//log.Println("lx",lx, "ly", ly,"dtot",math.Round(100.0*(d+efcvr))/100.0)
	//if f.Term != ""{err = PlotFtng(f.Colx, f.Coly, f.Fck, f.Fy, f.Df, f.Eo, d, f.Dmin, lx, ly, f.Nomcvr, f.Sloped, rez[mdx], f.Term)}
	f.Dused = d + efcvr
	f.Efcvr = efcvr
	f.Effd = d
	f.Lx = lx; f.Ly = ly
	//log.Println(rez[mdx])
	f.Rbr = rez[mdx]
	f.Rbropt = rez
	//for _, val := range rez{
	//	log.Println("dia",val[0],"nx",val[1],"spcx",val[2],"ny",val[5], "spcy",val[6])
	//}
	//log.Println("mux", mux, "kn-m muy", muy, "kn-m vux", vux, "kn vuy", vuy, "kn vp", vp, "kn")
	f.Mux = mux
	f.Muy = muy
	f.Vux = vux
	f.Vuy = vuy
	f.Vp = vp
	f.Dz = true
	f.Quant()
	if f.Web{f.Verbose = false}
	f.Table(f.Verbose)
	return
}

//FtngHmin returns the min horizontal dimension (WHAT)
func FtngHmin(pus, mxs []float64) (hmin float64) {
	//WHEN IS THIS USED
	var pu, mx float64
	for i := range pus {
		pu += pus[i]
		mx += mxs[i]
	}
	hmin = 12.0 * mx / pu
	return
}

//FtngPadAz is basically hulse sec 6.1, analysis of an isolated pad footing
func FtngPadAz(f *RccFtng) (err error){
	//hulse 6.1 pad footing analysis
	//possibly never gonna be used
	var ly, basecy, qmax, qmin, mur, vur, vp, vpp, psum, mu, pud,pul, puw, mud, mul, muw float64
	for i := range f.Pus{
		psum += f.Pus[i]
		mu += f.Mxs[i]
		switch i{
			case 0:
			pud = f.Pus[i]
			mud = f.Mxs[i]
			case 1:
			pul = f.Pus[i]
			mul = f.Mxs[i]
			case 2:
			puw = f.Pus[i]
			muw = f.Mxs[i]
		}
	}
	//psum := f.Pus[0] + f.Pus[1] + f.Pus[2]
	//mu := f.Mxs[0] + f.Mxs[1] + f.Mxs[2]
	//pud, pul, puw := f.Pus[0], f.Pus[1], f.Pus[2]
	//mud, mul, muw := f.Mxs[0], f.Mxs[1], f.Mxs[2]
	var step, lx float64
	var iter int
	//calc min base size
	lx = f.Lx
	ly = 0.5
	step = 0.5
	iter = 1
	pgck := 25.0
	if f.Code == 2{
		pgck = 24.0
	}
	for iter != -1 {
		ly += step
		if f.Typ == 0{
			lx = ly 
		} else {
			lx = ly - f.Coly + f.Colx
		}
		//log.Println(ColorRed,"lx, ly->",lx, ly,ColorReset)
		//if f.Shape == "square" {lx = ly}
		pu := psum + pgck*lx*ly*f.Based + (f.Df - f.Based) * f.Pgsoil
		
		switch iter {
		case 1:
			qmax = pu/ly/lx + 6.0*mu/lx/math.Pow(ly, 2)
			if qmax <= f.Sbc {
				ly -= step
				step = step / 10.0
				iter = 2
			}
		case 2:
			qmax = pu/ly/lx + 6.0*mu/lx/math.Pow(ly, 2)
			qmin = pu/ly/lx - 6.0*mu/lx/math.Pow(ly, 2)
			if qmin <= 0 {
				//part base is kompress
				ly = 0.5
				step = 0.5
				iter = 3
			} else if qmax <= f.Sbc {
				basecy = ly
				iter = -1
			}
		case 3:
			qmax = 2.0 * pu / 3.0 / lx / (ly/2.0 - mu/pu)
			if qmax > 0.0 && qmax <= f.Sbc {
				ly -= step
				step = step / 10.0
				iter = 4
			}
		case 4:
			qmax = 2.0 * pu / 3.0 / lx / (ly/2.0 - mu/pu)
			if qmax > 0.0 && qmax <= f.Sbc {
				basecy = 3 * (ly/2.0 - mu/pu)
				iter = -1
			}
		}
	}
	sfd := 1.2
	sfl := 1.2
	sfw := 1.2
	if puw == 0 && muw == 0 {
		switch f.Code {
		case 1:
			sfd = 1.5
			sfl = 1.5
			sfw = 0.0
		case 2:
			sfd = 1.4
			sfl = 1.6
			sfw = 0.0
		}
	}
	if pul == 0 && mul == 0 {
		switch f.Code {
		case 1:
			sfd = 1.5
			sfl = 0.0
			sfw = 1.5
		case 2:
			sfd = 1.4
			sfl = 0.0
			sfw = 1.4
		}
	}
	pu := pud*sfd + pul*sfl + puw*sfw + 24.0*lx*ly*f.Based*sfd
	mu = mud*sfd + mul*sfl + muw*sfw
	f1 := pu/ly/lx + 6.0*mu/lx/math.Pow(ly, 2)
	f2 := pu/ly/lx - 6.0*mu/lx/math.Pow(ly, 2)
	bcy := 0.0
	var y1, y2, f3, avp float64
	if f2 < 0.0 {
		//part el base kompres
		if ly/2.0-mu/pu < 0 {
			err = ErrFtngDim
			return
		}
		f1 = 2.0 * pu / 3.0 / lx / (ly/2.0 - mu/pu)
		bcy = 3.0 * (ly/2.0 - mu/pu)
		//bm at face of column
		y1 = ly/2.0 + f.Coly/2.0
		y2 = ly - y1
		if y2 < bcy {
			f3 = f1 * (bcy - y2) / bcy
			mur = (f1-f3)*lx*math.Pow(y2, 2)*1.0/3.0 + f3*lx*math.Pow(y2, 2)/2.0 - lx*y2*f.Based*y2*24.0*sfd/2.0
		} else {
			mur = f1*lx*bcy*(y2-bcy/2.0) - lx*y2*f.Based*y2*24.0*sfd/2.0
		}
		//shear at a distance 1.0 D from face of column
		y1 = y1 + f.Effd
		y2 = ly - y1
		switch {
		case y1 > ly:
			log.Println("critical shear section outside base")
			err = ErrFtngDim
			return
		case y2 > bcy:
			vur = f1 * lx * bcy
		default:
			f3 = f2 + (f1-f2)*y1/ly
			vur = (f1+f3)*lx*y2/2.0 - lx*y2*f.Based*24.0*sfd
		}
	} else {
		//all el base kompress
		//bm at face of column
		y1 = ly/2.0 + f.Coly/2.0
		y2 = ly - y1
		f3 = f2 + (f1-f2)*y1/ly
		mur = (f1-f3)*lx*math.Pow(y2, 2)/3.0 + f3*lx*math.Pow(y2, 2)/2.0 - lx*y2*f.Based*y2*24.0*sfd/2.0

		//shear at 1.0 * effd from face of column
		y1 = y1 + f.Effd
		y2 = ly - y1
		if y1 > ly {
			log.Println("critical shear section outside base")
			err = ErrFtngDim
			return
		}
		f3 = f2 + (f1-f2)*y1/ly
		vur = (f1+f3)*lx*y2/2.0 - lx*y2*f.Based*24.0*sfd
	}
	//punching shear
	sp := f.Coly/2.0 + 1.5*f.Effd
	if sp > ly/2.0 || sp > lx/2.0 {
		log.Println("punching shear perimeter outside base")
		err = ErrFtngDim
		return
	}
	vpp = 2.0*(f.Colx+f.Coly) + 12.0*f.Effd
	avp = (f.Colx + 3.0*f.Effd) * (f.Coly + 3.0*f.Effd)
	f1 = pu/ly/lx - f.Based*24.0*sfd
	vp = f1 * (ly*lx - avp)
	f.Ly = ly
	if f.Typ == 0{
		f.Lx = ly
	} else {
		f.Lx = ly - f.Coly + f.Colx
	}
	f.Basecy = basecy
	f.Qmax = qmax
	f.Qmin = qmin
	f.Mur = mur
	f.Vur = vur
	f.Vp = vp
	f.Vpp = vpp
	err = nil
	//log.Println("dims->",f.Lx, f.Ly)
	return
}

//FtngBxOz was an attempted write of Ozmen, design of a rectangular isolated footing
func FtngBxOz(bx, by, pu, mx, my float64){
	/*
	   rips off g.ozmen 2011
	*/
	sxs := make([]float64, 4)
	ar := bx * by
	ix := bx * math.Pow(by, 3) / 12.0
	iy := by * math.Pow(bx, 3) / 12.0
	xv := -my / pu
	yv := -mx / pu
	ixy := 0.
	//compute stresses at corners
	sxs[0] = pu/ar + mx*by/2.0/ix + my*bx/2.0/iy
	sxs[1] = pu/ar - mx*by/2.0/ix + my*bx/2.0/iy
	sxs[2] = pu/ar - mx*by/2.0/ix - my*bx/2.0/iy
	sxs[3] = pu/ar + mx*by/2.0/ix - my*bx/2.0/iy
	if sxs[0] > 0 && sxs[1] > 0 && sxs[2] > 0 && sxs[3] > 0 {
		log.Println("all corner stresses are +ve")
		log.Println("l left, u left, u right, l right")
		log.Println(sxs[0], sxs[1], sxs[2], sxs[3])
		return
	}
	log.Println("l left, u left, u right, l right")
	log.Println(sxs[0], sxs[1], sxs[2], sxs[3])
	da := (sxs[0] - sxs[3]) / bx
	dc := (sxs[0] - sxs[1]) / by
	a := sxs[0] / da
	c := sxs[0] / dc
	//log.Println(a,c)
	tga := c / a
	var iter int
	var x0, y0, ug, vg, aprev, cprev float64
	var tol float64 = 0.001
	ugs := make([]float64, 3)
	vgs := make([]float64, 3)
	fas := make([]float64, 3)
	ixs := make([]float64, 3)
	iys := make([]float64, 3)
	ixys := make([]float64, 3)
	for {
		iter++
		ar = 0.0
		ix = 0.0
		iy = 0.0
		ixy = 0.0
		ug = 0.0
		vg = 0.0
		var a1, c1, a2, c2, a3, c3, tgb, xvi, yvi float64
		//get table 2 values from a and c
		switch {
		case a > bx && c > by && bx/a+by/c < 1.0:
			//zone 1
			a1 = bx
			c1 = by
			a2 = 0
			c2 = 0
			a3 = 0
			c3 = 0
		case a < bx && c > by:
			//zone 2
			a1 = (c - by) / tga
			c1 = by
			a2 = 0
			c2 = 0
			a3 = a - a1
			c3 = by
		case a > bx && c < by:
			//zone 3
			a1 = 0
			c1 = 0
			a2 = bx
			c2 = (a - bx) * tga
			a3 = bx
			c3 = c - c2
		case a > bx && c > by && bx/a+by/c > 1.0:
			//zone 4
			a1 = (c - by) / tga
			c1 = by
			a2 = bx - a1
			c2 = (a - bx) * tga
			a3 = bx - a1
			c3 = c2
		case a < bx && c < by:
			//zone 5
			a1 = 0
			c1 = 0
			a2 = 0
			c2 = 0
			a3 = a
			c3 = c
		}
		//get geometric values as per table 3
		for i := 0; i < 3; i++ {
			switch i {
			case 0:
				//rect part 1
				ugs[i] = a1 / 2.0
				vgs[i] = c1 / 2.0
				fas[i] = a1 * c1
				ixs[i] = a1 * math.Pow(c1, 3) / 12.0
				iys[i] = c1 * math.Pow(a1, 3) / 12.0
				ixys[i] = 0.0
			case 1:
				//rect part 2
				ugs[i] = a1 + a2/2.0
				vgs[i] = c2 / 2.0
				fas[i] = a2 * c2
				ixs[i] = a2 * math.Pow(c2, 3) / 12.0
				iys[i] = c2 * math.Pow(a2, 3) / 12.0
				ixys[i] = 0.0
			case 2:
				//tri part 3
				ugs[i] = a1 + a3/3.0
				vgs[i] = c2 + c3/3.0
				fas[i] = a3 * c3 / 2.0
				ixs[i] = a3 * math.Pow(c3, 3) / 36.0
				iys[i] = c3 * math.Pow(a3, 3) / 36.0
				ixys[i] = -math.Pow(a3, 2) * math.Pow(c3, 2) / 72.0
			}
		}
		ar = fas[0] + fas[1] + fas[2]
		ug = (ugs[0]*fas[0] + ugs[1]*fas[1] + ugs[2]*fas[2]) / ar
		vg = (vgs[0]*fas[0] + vgs[1]*fas[1] + vgs[2]*fas[2]) / ar
		for i := range ugs {
			ei := ug - ugs[i]
			fi := vg - vgs[i]
			ix += ixs[i] + fas[i]*math.Pow(fi, 2)
			iy += iys[i] + fas[i]*math.Pow(ei, 2)
			ixy += ixys[i] + fas[i]*ei*fi
		}
		xvi = xv + bx/2.0 - ug
		yvi = yv + by/2.0 - vg
		tgb = yvi / xvi
		tga = (ix - ixy*tgb) / (iy*tgb - ixy)
		x0 = -(iy + ixy/tga) / (xvi * ar)
		y0 = x0 * tga
		aprev = a
		cprev = c
		a = ug + x0 + vg/tga
		c = vg + y0 + ug*tga
		if math.Abs(aprev-a) < tol && math.Abs(cprev-c) < tol {
			log.Println("iteration converged")
			log.Println("a, c->", a, c)
			break
		}
		if iter > 3000 {
			log.Println("max iteration limit exceeded")
		}
	}
	cpts := [][]float64{{-bx / 2.0, -by / 2.0}, {-bx / 2.0, by / 2.0}, {bx / 2.0, by / 2.0}, {bx / 2.0, -by / 2.0}}
	for i, pt := range cpts {
		xi := pt[0] + bx/2.0 - ug
		yi := pt[1] + by/2.0 - vg
		sxs[i] = pu / ar * (1.0 - (xi*tga+yi)/y0)
		if sxs[i] < 0.0 {
			sxs[i] = 0.0
		}
	}
	log.Println("corner stresses->")
	for i := range sxs {
		log.Println(i, sxs[i])
	}
	log.Println("max stress->", a*pu/x0/ar)
}

//FtngBxEval is the entry func for biaxial footing opt routines
//TODO - as in chaudhry/maity
func FtngBxEval(colx, coly, fck, fy, df, sbc, pgck, pgsoil, nomcvr float64, pus, mxs, mys, psfs []float64, p *prat) {
	var pu, mx, my, mux, muy, bx, by, d, psx, psy float64
	var cons int
	//get service level loads
	for i := range pus {
		pu += pus[i]
		mx += mxs[i]
		my += mys[i]
		mux += mxs[i] * psfs[i]
		muy += mys[i] * psfs[i]
	}
	by = 10.0 * math.Ceil(p.pos[0]/10.0)
	bx = my * by / mx
	bx = 10.0 * math.Ceil(bx/10.0)
	d = 10.0 * math.Ceil(p.pos[1]/10.0)
	psx = p.pos[2]
	psy = p.pos[3]
	//astx = psx * bx * by; asty = psy * bx * by
	log.Println(bx, psx, psy)
	//check for +ve pressure
	smx := sbc - pgck*(d+nomcvr) - (df-d-nomcvr)*pgsoil
	c1 := smx*my*math.Pow(by, 3) - pu*mx*by - 12.0*math.Pow(mx, 2)
	if c1 < 0 {
		cons++
	}

	//check for bending moment
	//mxr := (pu * math.Pow(by,2) + 2.0 * mux * (2.0 * by + coly))*math.Pow(by - coly,2)/8.0/math.Pow(by,3)
	//mxp := 0.87 * fy * astx * d * (1.0 - astx * fy/bx/d/fck)

}


/*

	//check for bending moment
	t1 := pu * (ly - coly)/2.0/ly + 3.0 * mxl * (ly + coly) * (ly - coly)/2.0/math.Pow(ly,3)
	t2 := pu * math.Pow(ly, 2) * (math.Pow(ly, 2) - math.Pow(coly, 2)) + 4.0 * mxl * (math.Pow(ly, 3) - math.Pow(coly, 3))
	t3 := 4.0 * pu * math.Pow(ly, 2) * (ly - coly) + 12.0 * mxl * (math.Pow(ly, 2) - math.Pow(coly, 2))
	mux = t1 * (t2/t3 - coly/2.0)
	log.Println("mux->",mux)
	t1 = pu * (lx - colx)/2.0/lx + 3.0 * myl * (lx + colx) * (lx - colx)/2.0/math.Pow(lx,3)
	t2 = pu * math.Pow(lx, 2) * (math.Pow(lx, 2) - math.Pow(colx, 2)) + 4.0 * myl * (math.Pow(lx, 3) - math.Pow(colx, 3))
	t3 = 4.0 * pu * math.Pow(lx, 2) * (lx - colx) + 12.0 * myl * (math.Pow(lx, 2) - math.Pow(colx, 2))
	muy = t1 * (t2/t3 - colx/2.0)
	log.Println("muy->",muy)
*/
