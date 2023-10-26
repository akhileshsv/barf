package barf

import (
	"fmt"
	"log"
	"math"
	"time"
	"errors"
	"strings"
	"math/rand"
	"github.com/olekukonko/tablewriter"
	kass"barf/kass"
)

//BmDz is the entry func for truss/frame design routines
//when input is a kass.Model
func BmDz(mem *kass.Mem, grp int, params []interface{}){
	
}

//Init initalizes a WdBm struct
func (b *WdBm) Init() (err error){
	if b.Grp != 0 {err = b.Prp.Init(b.Grp)}
	if err != nil{return}
	fmt.Println(b.Prp)
	//default pinned
	switch b.Styp{
		case 0:
		//b.Cbi = 0.85
		case 1:
		//b.Cbi = 0.8
	}
	if len(b.Dims) != 0{
		s := kass.SecGen(b.Styp, b.Dims)
		b.Sec = s
	}
	if b.Clctyp == 0{b.Clctyp = 1}
	if b.Nsecs == 0{b.Nsecs = 3}
	return
}

//BmDesign is the entry func from menu/flags
func BmDesign(b *WdBm)(err error){
	switch b.Dtyp{
		case 0:
		//use coeff to design
		err = BmDzCs(b)
		if err != nil{
			return
		}
		//b.Dz = true
		b.Table(b.Verbose)
		case 1:
		//check beam dims
		vals, chk := BmChk(b)
		if !chk {err = errors.New("beam dimension check failed")}
		if b.Verbose{
			fmt.Println(ColorGreen, "beam->",b.Dims, b.Styp, b.Mu, b.Vu, "\nOK?",chk,ColorReset)
			fmt.Println(ColorCyan, vals, ColorReset)
		}
		case 2:
		//cbeam design? dunno
	}
	
	return
}

//BmChk checks a solid beam section for moment, shear and deflection
//chapter 8, abel o.o
func BmChk(b *WdBm) (vals []float64, chk bool){
	var vfac float64
	//shear deformation constant - vdfac := 1.1
	switch b.Styp{
	//write a func to calc vfac for each styp
		case 0:
		vfac = 1.33
		case 1:
		vfac = 1.5
	}
	var rb, rk, fball float64
	fball = b.Prp.Fcb
	//check for lateral-torsional buckling
	//log.Println(ColorCyan,b.Dims,ColorReset)
	switch b.Styp{
		case 0:
		case 1:
		rb = math.Sqrt(b.Lspan * b.Dims[1]/math.Pow(b.Dims[0],2))
		rk = 0.775 * math.Sqrt(b.Prp.Em/b.Prp.Fcb)
		//log.Println(rb, rk)
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
		default:
		//wot?
	}
	
	fb := b.Mu/b.Sec.Prop.Sxx
	fv := vfac * b.Vu/b.Sec.Prop.Area
	dall := b.Lspan/b.Drat
	var brchk bool
	//check for bearing
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
	chk = fb <= fball && fv <= b.Prp.Fv && b.Dmax <= dall && brchk
	
	if b.Spam{
		log.Println("fball, fb, fv, vfac->",fball, fb, fv, vfac,"n/mm2")
		log.Println("dmax, dchk->",dall, b.Dmax <= dall)
	}
	vals = []float64{b.Vu, b.Mu, b.Dmax,dall,fb,fball, fv, b.Prp.Fv, b.Sec.Prop.Area}
	log.Println(ColorRed,chk,ColorReset)
	return
} 



//BmDzCs designs a continuous beam span using forumulae as per abel chapter 8
//uses coefficients to either choose a section or check span for a section
func BmDzCs(b *WdBm) (err error){
	var basedims [][]float64
	var mnarea float64
	if b.Drat == 0.0{
		if b.DL == 0.0{
			b.Drat = 360.0
		} else {
			b.Drat = 180.0
		}
	}
	//if b.Endc < 2 {b.Lspan += (b.Lbl + b.Rbl)/2.0}
	switch{
		case len(b.Dims) == 0:
		//get section rez
		switch b.Styp{
			case 0:
			//round timber
			basedims = kass.TmbrDims0
			case 1:
			//rect timber
			if len(kass.TmbrDims1) == 0{
				kass.GenTmbrDims1()
			}
			basedims = kass.TmbrDims1
			//basedims = [][]float64{{75,100},{75,150}}
		}
		//fmt.Println("HEAR")
		var mdx int
		for _, dim := range basedims{
			if len(b.Rez) == b.Nsecs{break}
			if b.Styp == 1{
				if dim[1]/dim[0] > 3.0 || dim[0] > dim[1]{continue}
			}
			b.Dims = dim
			wul := b.DL + b.LL
			s := kass.SecGen(b.Styp, dim)
			wdl := s.Prop.Area * b.Prp.Pg * 9.8 * 1e-6 
			ei := s.Prop.Ixx * b.Prp.Em
			if b.Selfwt{
				//log.Println("wdl->",wdl)
				wul = wdl + b.DL + b.LL
			}
			b.Sec = s
			mul, vul, dmax := kass.BmUvalCs(wul, b.Lspan, ei, b.Nspans, b.Endc)
			fmt.Println(ColorYellow,wul, mul, vul, dmax)
			b.Mu = mul; b.Vu = vul; b.Dmax = dmax
			if val, ok := BmChk(b); ok{
				b.Rez = append(b.Rez, dim)
				b.Vals = append(b.Vals, val)
				if len(b.Rez) == 1{
					mdx = 0
					mnarea = b.Sec.Prop.Area
				} else if mnarea > b.Sec.Prop.Area{
					mdx = len(b.Rez) - 1
					mnarea = b.Sec.Prop.Area
				}
			}
		}
		if len(b.Rez) == 0{
			err = errors.New("no suitable section found")
			return
		}
		b.Dims = b.Rez[mdx]
		s := kass.SecGen(b.Styp, b.Dims)
		b.Sec = s
		//log.Println("section->",b.Dims)
		default:
		//fuck get spans based on 3 condishuns
	}
	b.Dz = true
	b.Table(true)
	return
}

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
