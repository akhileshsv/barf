package barf

import (
	"fmt"
	"log"
	"math"
	"time"
	"strings"
	"math/rand"
	"github.com/olekukonko/tablewriter"
	kass"barf/kass"
)

//NOTE - FOR DIFFERENT DURATIONS OF LOAD, MODIFICATION FACTOR (K2)
//perm. stress  = k2 * perm. stress
//normal, two months, seven days, wind/earthquake, instantaneous - 1.0, 1.15, 1.25, 1.33, 2.0

//Init initalizes a WdBm struct
func (b *WdBm) Init() (err error){
	if b.Grp != 0 {err = b.Prp.Init(b.Grp)}
	if err != nil{return}
	//fmt.Println(b.Prp)
	bstyp := b.Styp
	dims := b.Dims
	//default pinned
	switch b.Styp{
		case 0:
		//b.Cbi = 0.85
		case 1:
		//b.Cbi = 0.8
		case 4:
		switch b.Solid{
			case true:
			//built up solid beam as in ex. 13.13, ramc
			bstyp = 1
			dims  = []float64{b.Dims[0],b.Dims[1]}
		}
	}
	if len(b.Dims) != 0 && b.Sec.Prop.Area == 0.0{
		s := kass.SecGen(bstyp, dims)
		b.Sec = s
	}
	if b.Clctyp == 0{b.Clctyp = 1}
	if b.Nsecs == 0{b.Nsecs = 3}
	//rewrite dis
	if b.Ke == 0{
		if b.Kendc == 0{
			b.Kendc = 1
		}
		switch b.Kendc{
			case 1:
			//f-f, no sway
			b.Ke = 0.5
			// c.Kerst = 4.0
			case 2:
			//p-f, no sway
			b.Ke = 0.7
			case 3:
			//p-p, no sway
			b.Ke = 1.0
			case 4:
			//f-f, sway
			b.Ke = 1.0
			case 5:
			//f-free, sway
			b.Ke = 2.0
			case 6:
			//f-p, sway
			b.Ke = 2.0
		}	
	}

	return
}

//BmDesign is the entry func from menu/flags
//IS IT. WHY
func BmDesign(b *WdBm)(err error){
	//log.Println("HERE, Dtyp, endc",b.Dtyp)
	b.Init()
	var basedims [][]float64
	switch b.Styp{
		case 0:
		basedims = kass.TmbrDims0
		case 1:
		if !kass.TmbrSort{
			kass.GenTmbrDims1()
			basedims = kass.TmbrDims1
		}
	}
	switch b.Styp{
		case 0, 1:
		for _, dim := range basedims{
			if len(b.Rez) == b.Nsecs{
				break
			}
			b.Dims = dim
			ok, vals, _ := BmChk(b)
			if ok{
				b.Rez = append(b.Rez, dim)
				b.Vals = append(b.Vals, vals)
			}
			b.Dims = []float64{}
		}
		if len(b.Rez) == 0{
			err = fmt.Errorf("no suitable section found")
			return
		}
		
		b.Dims = b.Rez[0]
		s := kass.SecGen(b.Styp, b.Dims)
		b.Sec = s
		err = nil
		b.Dz = true
		b.Table(b.Verbose)
		return
	}
	switch b.Styp{
		case 4:
		case 26:
	}
	return
}

//BmAz analyzes a beam given a list of ldcases
//see mosley4.1
func (b *WdBm) BmAz()(err error){
	mod := kass.Model{}
	switch{
		case b.Nspans == 0:
		if len(b.Lspans) == 0{
			if b.Lspan == 0{
				err = fmt.Errorf("no geom data for beam -> lspan %.f lspans %v nspans %v",b.Lspan, b.Lspans, b.Nspans)
				return
			}
			b.Lspans = []float64{b.Lspan}
		}
	}
	b.Nspans = len(b.Lspans)
	dl := b.DL; ll := b.LL
	if b.Selfwt{
		wdl := b.Sec.Prop.Area * b.Prp.Pg * 9.8 * 1e-6 
		dl += wdl
	}
	x := 0.0
	mod.PSFs = []float64{1.0,1.0,0.7}
	mod.Coords = append(mod.Coords, []float64{0,0})
	mod.Supports = append(mod.Supports, []int{1, -1, 0})
	mod.Em = [][]float64{{b.Prp.Em}}
	mod.Cp = [][]float64{{b.Sec.Prop.Ixx, b.Sec.Prop.Area}}
	//"Em":[[10000]],"Cp":[[96]],
	for i, lspan := range b.Lspans{
		x += lspan
		jb := i + 1; je := i + 2
		mdx := i + 1
		mod.Coords = append(mod.Coords, []float64{x, 0})
		mod.Mprp = append(mod.Mprp, []int{jb, je, 1, 1, 0})
		mod.Supports = append(mod.Supports, []int{je, -1, 0})
		if dl > 0.0{mod.Msloads = append(mod.Msloads, []float64{float64(mdx), 3.0, dl, 0.0, 0.0, 0.0, 1.0})}
		if ll > 0.0{mod.Msloads = append(mod.Msloads, []float64{float64(mdx), 3.0, ll, 0.0, 0.0, 0.0, 2.0})}
	}
	mod.Msloads = append(mod.Msloads, b.Ldcases...)
	switch b.Nspans{
		case 1:
		frmrez, _ := kass.CalcBm1d(&mod, 2)
		bmrez := kass.CalcBmSf(&mod, frmrez,false)
		bm := bmrez[1]
		b.Vu = math.Abs(bm.Maxs[0])
		b.Mu = math.Abs(bm.Maxs[1])
		b.Dmax = math.Abs(bm.Maxs[2])
		default:
		//cbeam - use calcbmenv(mod)
		//or use mosh.cbeam? better plots
	}
	return
}

//BmChk checks a solid beam section for moment, shear and deflection
//chapter 8, abel o.o
func BmChk(b *WdBm) (ok bool, vals []float64, err error){
	if b.Sec.Prop.Area == 0.0{
		err = b.Init()
		if err != nil{
			return
		}
	}
	if b.Mu == 0{
		//log.Println("CALCING AGAIN")
		switch b.Dtyp{
			case 0:
			//use coefficients
			wul := b.DL + b.LL
			wdl := b.Sec.Prop.Area * b.Prp.Pg * 9.8 * 1e-6 
			ei := b.Sec.Prop.Ixx * b.Prp.Em
			if b.Selfwt{
				wul = wdl + b.DL + b.LL
			}
			mul, vul, dmax := kass.BmUvalCs(wul, b.Lspan, ei, b.Nspans, b.Endc)
			if b.Spam{log.Println("wul, mul, vul, dmax->",wul, mul, vul, dmax)}
			b.Mu = mul; b.Vu = vul; b.Dmax = dmax
			case 1:
			//actual load cases (cbeam)
			err = b.BmAz()
			if err != nil{
				return
			}
		}
	}
	if b.Drat == 0.0{
		if b.DL == 0.0{
			b.Drat = 360.0
		} else {
			b.Drat = 180.0
		}
	}
	var rb, rk, vfac, fb, fball, fv, dall, vred float64
	//H = VQ/Ib
	//shear deformation constant - vdfac := 1.1
	switch b.Styp{
	//write a func to calc vfac for each styp
		case 0:
		vfac = 1.33
		case 1:
		vfac = 1.5
	}
	switch b.Code{
		case 1:
		//is883
		//min. width is max. of (lspan/50.0 & 50 mm)
		//if depth > 3.0 * width // lspan > 50 * width - include lateral restraints at every 50.0 * width
		//min. bearing of 75,, 
		var k3, k4, k5 float64
		fball = b.Prp.Fcb
		vred = 1.0
		//get form factor k3/k4
		switch b.Styp{
			case 0:
			//circular beam
			k5 = 1.18
			fball = fball * k5
			case 1:
			//rect beam
			d2 := math.Pow(b.Dims[1],2)
			k3 = 0.81 * (d2 + 89400.0)/(d2 + 55000.0)
			fball = fball * k3
			case 4:
			//box beam
			B := b.Dims[0]; D := b.Dims[1]; bw := b.Dims[2]; df := b.Dims[3]
			switch b.Solid{
				case false:
				//hollow core box beam
				d2 := math.Pow(b.Dims[1],2)
				p1 := (D - df)/2.0/D
				q1 := (B - bw)/2.0/B
				y := p1 * p1 * (6.0 - 8.0 * p1 + 3.0 * p1 * p1) * (1.0 - q1) + q1
				t1 := (d2 + 89400.0)/(d2 + 55000.0)
				k4 = 0.8 + 0.8 * y * (t1 - 1.0)
				fball = fball * k4
				case true:
				//solid beam of tplnk thickness
				d2 := math.Pow(b.Dims[1],2)
				k3 = 0.81 * (d2 + 89400.0)/(d2 + 55000.0)
				fball = fball * k3
			}
			case 12:
			//built up i beam
			bw := b.Dims[0]
			h := b.Dims[1]
			tf := b.Dims[2]
			tw := b.Dims[3]
			p1 := tf/h
			q1 := tw/bw
			d2 := math.Pow(b.Dims[1],2)
			y := p1 * p1 * (6.0 - 8.0 * p1 + 3.0 * p1 * p1) * (1.0 - q1) + q1
			t1 := (d2 + 89400.0)/(d2 + 55000.0)
			k4 = 0.8 + 0.8 * y * (t1 - 1.0)
			fball = fball * k4
		}
		fb = b.Mu/b.Sec.Prop.Sxx
		case 2:
		//abel o/madison/fao
		
		//check for lateral-torsional buckling
		rb = b.Ke * b.Lspan/b.Sec.Prop.Rxx
		vred = 1.0
		switch {
		case rb <= 10.0:
			fball = b.Prp.Fcb
		case rb <= rk:
			fball = b.Prp.Fcb * (1.0 - 0.33 * math.Pow(rb/rk,4))
		case rb <= 50.0:
			fball = 0.40 * b.Prp.Em /math.Pow(rb,2)
		default:
			fball = 0.40 * b.Prp.Em /math.Pow(rb,2)
		}
		switch b.Styp{
			case 0:
			case 1:
			//rb = math.Sqrt(b.Lspan * b.Dims[1]/math.Pow(b.Dims[0],2))
			rk = 0.775 * math.Sqrt(b.Prp.Em/b.Prp.Fcb)
			if b.Notch{vred = (b.Dims[1] - b.Dn)/b.Dims[1]}
			//log.Println(rb, rk)
			default:
			//wot?
		}
	}
	fb = b.Mu/b.Sec.Prop.Sxx
	fv = vfac * b.Vu/b.Sec.Prop.Area
	
	if b.Notch{
		d1 := b.Dims[1] - b.Dn
		fv = vfac * b.Vu/b.Dims[0]/d1
	}
	b.Sec.CalcQ()
	fvq := b.Vu * b.Sec.Prop.Vfx
	if b.Spam{log.Println("fv vs fvq -> fv ",fv, " fvq ",fvq, " fv - fvq ", fv - fvq)}
	dall = b.Lspan/b.Drat
	brchk := true
	//check for bearing
	if b.Brchk{
		//TODO ADD MOD FACTOR
		if b.Lbl > 0.0 && b.Rbl > 0.0{
			switch b.Styp{
				case 1:
				brchk = (b.Prp.Fcp >= b.Vu/b.Dims[0]/b.Lbl && b.Prp.Fcp >= b.Vu/b.Dims[0]/b.Rbl)
			}

		} else if b.Lbl > 0.0{
			switch b.Styp{
				case 1:
				brchk = b.Prp.Fcp >= b.Vu/b.Dims[0]/b.Lbl
			}
		}
	}
	ok = fb <= fball && fv <= b.Prp.Fv * vred && b.Dmax <= dall && brchk
	
	if b.Spam{
		log.Println("fball, fb, fv, vfac->",fball, fb, fv, vfac,"n/mm2")
		log.Println("dmax, dchk->",dall, b.Dmax <= dall)
	}
	vals = []float64{b.Vu, b.Mu, b.Dmax,dall,fb,fball, fv, b.Prp.Fv, b.Sec.Prop.Area}
	if b.Spam{log.Println(ColorRed,ok,ColorReset)}
	return
} 

// func (b *WdBm) ReadSpans()(err error){
// 	if 
// }

// func BmDz(b *WdBm) (err error){
// 	// var basedims [][]float64
// 	// var mnarea float64
	
// 	// if b.Drat == 0.0{
// 	// 	if b.DL == 0.0{
// 	// 		b.Drat = 360.0
// 	// 	} else {
// 	// 		b.Drat = 180.0
// 	// 	}
// 	// }
// 	if b.Mu == 0.0{
// 		//get uvals
// 		if len(b.Ldcases) == 0{
// 			//use coefficients
// 		} else {
// 			//
// 		}
// 	}
// 	return
// }

//BmDzCs designs a continuous beam span using forumulae as per abel chapter 8
//uses coefficients to either choose a section or check span for a section
// func BmDzCs(b *WdBm) (err error){
// 	var basedims [][]float64
// 	var mnarea float64
// 	if b.Drat == 0.0{
// 		if b.DL == 0.0{
// 			b.Drat = 360.0
// 		} else {
// 			b.Drat = 180.0
// 		}
// 	}
// 	//if b.Endc < 2 {b.Lspan += (b.Lbl + b.Rbl)/2.0}
// 	switch{
// 		case len(b.Dims) == 0:
// 		//get section rez
// 		switch b.Styp{
// 			case 0:
// 			//round timber
// 			basedims = kass.TmbrDims0
// 			case 1:
// 			//rect timber
// 			if len(kass.TmbrDims1) == 0{
// 				kass.GenTmbrDims1()
// 			}
// 			basedims = kass.TmbrDims1
// 			//basedims = [][]float64{{75,100},{75,150}}
// 		}
// 		fmt.Println("HEAR")
// 		var mdx int
// 		for _, dim := range basedims{
// 			if len(b.Rez) == b.Nsecs{break}
// 			if b.Styp == 1{
// 				if dim[1]/dim[0] > 3.0 || dim[0] > dim[1]{continue}
// 			}
// 			b.Dims = dim
// 			// wul := b.DL + b.LL
// 			// s := kass.SecGen(b.Styp, dim)
// 			// wdl := s.Prop.Area * b.Prp.Pg * 9.8 * 1e-6 
// 			// ei := s.Prop.Ixx * b.Prp.Em
// 			// if b.Selfwt{
// 			// 	//log.Println("wdl->",wdl)
// 			// 	wul = wdl + b.DL + b.LL
// 			// }
// 			// b.Sec = s
// 			// mul, vul, dmax := kass.BmUvalCs(wul, b.Lspan, ei, b.Nspans, b.Endc)
// 			// fmt.Println(ColorYellow,wul, mul, vul, dmax)
// 			// b.Mu = mul; b.Vu = vul; b.Dmax = dmax
// 			if ok, val, _ := BmChk(b); ok{
// 				b.Rez = append(b.Rez, dim)
// 				b.Vals = append(b.Vals, val)
// 				if len(b.Rez) == 1{
// 					mdx = 0
// 					mnarea = b.Sec.Prop.Area
// 				} else if mnarea > b.Sec.Prop.Area{
// 					mdx = len(b.Rez) - 1
// 					mnarea = b.Sec.Prop.Area
// 				}
// 			}
// 		}
// 		if len(b.Rez) == 0{
// 			err = errors.New("no suitable section found")
// 			return
// 		}
// 		b.Dims = b.Rez[mdx]
// 		s := kass.SecGen(b.Styp, b.Dims)
// 		b.Sec = s
// 		//log.Println("section->",b.Dims)
// 		default:
// 		//fuck get spans based on 3 condishuns
// 	}
// 	b.Dz = true
// 	b.Table(true)
// 	return
// }

//PlyUdlSpn returns max span of plywood for sheathing
//ABBE NOTE DOWN REFERENCES what are these strange calcs
func PlyUdlSpn(d, wdl, wll float64, scon int) (lspan float64, err error){
	var lm, lv, ld, dfrat, ei, mc, vc, wul float64
	if wll == 0.0{wll = 1.5}
	switch d{
		case 12.0:
		mc = 0.2; vc = 6.16; ei = 1.07; wdl += 9.24 * 9.8/1e3
		case 19.0:
		mc = 0.34; vc = 9.75; ei = 2.73; wdl += 13.86 * 9.8/1e3
		default:
		log.Println("invalid plywood thickness")
		err = ErrInp
		return
	}
	wul = wdl + wll
	dfrat = 360.0
	switch scon{
		case 0:
		//cantilever
		dfrat = 150.0
		lm = math.Sqrt(2.0 * mc/wul)
		lv = 2.0 * vc/wul
		ld = math.Pow(8.0 * ei/wul/dfrat, 1.0/3.0)
		case 1:
		//1 span ss
		lm = math.Sqrt(8.0*mc/wul)
		lv = 2.0 * vc/wul
		ld = math.Pow(384.0 * ei/wul/dfrat/5.0, 1.0/3.0)
		case 2:
		//2 span cs
		lm = math.Sqrt(8.0*mc/wul)
		lv = 8.0 * vc/wul/5.0
		ld = math.Pow(185.0 * ei/wul/dfrat, 1.0/3.0)
		case 3:
		//n span cs
		lm = math.Sqrt(10.0*mc/wul)
		lv = 8.0 * vc/wul/5.0
		ld = math.Pow(185.0 * ei/wul/dfrat, 1.0/3.0)
		default:
		log.Println("invalid support condition")
		err = ErrInp
		return 
	}
	log.Println("l moment",lm,"lshear",lv,"ldef",ld)
	lspan = ld
	if lspan > lm {lspan = lm}
	if lspan > lv {lspan = lv}
	lspan = 1000.0 * lspan
	return 
}


//Table generates an ascii table report for a WdBm
func (b *WdBm) Table(printz bool){
	if b.Title == ""{
		if b.Id == 0{
			b.Id = rand.Intn(666)
		}
		b.Title = fmt.Sprintf("tmbr_bm_%v",b.Id)
	}
	rezstr := new(strings.Builder)
	hdr := fmt.Sprintf("%s\ntimber beam report\ndate-%s\n%s\n%s\n",ColorYellow,time.Now().Format("2006-01-02"),b.Title,ColorReset)
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)
	var row string
	table.SetCaption(true,"timber properties")
	table.SetHeader([]string{"group","section type","specific gravity","elastic modulus\n(n/mm2)"})
	row = fmt.Sprintf("%v, %v, %.2f, %.2f",b.Grp,b.Styp,b.Prp.Pg,b.Prp.Em)
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorBlue)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"allowable stresses (N/mm2)")
	table.SetHeader([]string{"bending","tension","compression\npar.to grain","compression\nperp.to grain","shear\npar.to grain","shear\nperp.to grain"})
	row = fmt.Sprintf("%.2f,%.2f,%.2f,%.2f,%.2f,%.2f",b.Prp.Fcb,b.Prp.Ft,b.Prp.Fc,b.Prp.Fcp,b.Prp.Fv,b.Prp.Fvp)
	table.Append(strings.Split(row,","))
	table.Render()
	
	if b.Dz{
		rezstr.WriteString(ColorPurple)
	table = tablewriter.NewWriter(rezstr)
		table.SetHeader([]string{"dims","area\n(mm2)","vu\n(n)","mu\n(n-mm)", "dmax\n(mm)","d perm.\n(mm)","fb\n(n/mm2)","fb perm.\n(n/mm2)", "fv\n(n/mm2)", "fv perm.\n(n/mm2)"})	
		for i, dims := range b.Rez{
			vu, mu, dmax, dall, fb, fball, fv, prpv, ar := b.Vals[i][0], b.Vals[i][1],b.Vals[i][2],b.Vals[i][3],b.Vals[i][4],b.Vals[i][5],b.Vals[i][6],b.Vals[i][7], b.Vals[i][8]
			row = fmt.Sprintf("%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f",dims,ar,vu, mu, dmax, dall, fb, fball, fv, prpv)
			table.Append(strings.Split(row, ","))
		}
		table.Render()
	}
	rezstr.WriteString(ColorReset)
	b.Report = fmt.Sprintf("%s",rezstr)
	if printz{
		fmt.Println(b.Report)
	}
	return
}
