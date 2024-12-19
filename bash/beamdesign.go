package barf

import (
	//"os"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
	//"sort"
	"math/rand"
	kass "barf/kass"
	"github.com/go-gota/gota/dataframe"
	"github.com/olekukonko/tablewriter"
)

//Bm is a struct to store steel beam fields
//see mosley/spencer section 6.2
type Bm struct{
	Title                                            string
	Sname                                            string
	Name                                             string
	Bstyp                                            int
	Id                                               int
	Bf, H, Tf, Tw, B, D                              float64
	Bp, Tp                                           float64
	Grd, Styp, Nsecs, Sdx                            int
	Code, Endc, Dtyp                                 int
	H1, H2, Lx, Ly, Tx, Ty, Mx, My, Vx, Vy, Pu, Pfac float64
	Lspan, Lbr, Tbr, Klx, Kly                        float64
	Lbl, Rbl                                         float64
	Rez                                              []int
	Vals                                             [][]float64
	PSFs                                             []float64
	Maxs                                             []float64
	Vu, Mu                                           float64 
	Dmax                                             float64
	Ml, Mr                                           float64
	DL, LL, WL                                       float64
	Theta                                            float64
	Ldcases                                          [][]float64
	Nspans                                           float64
	Lspans                                           []float64
	Yeolde                                           bool
	Brchk                                            bool
	Verbose                                          bool
	Blurb                                            bool
	Spam                                             bool
	Braced                                           bool
	Lsb                                              bool
	Dz                                               bool
	Web                                              bool
	Selfwt                                           bool
	Isclvr                                           bool
	Prln                                             bool
	Store                                            bool
	Frame                                            bool
	Report                                           string
	Term                                             string
	Kostin                                           float64
	Mindx                                            int
	Ssec                                             kass.StlSec
	Ssecs                                            []kass.StlSec
	Sdxs                                             []int
	Dims                                             [][]float64
	Params                                           []float64
}

//Init inits default values
func (b *Bm) Init()(err error){
	if b.Sname == ""{
		if val, ok := kass.StlSnames[b.Styp]; !ok{
			err = fmt.Errorf("no base section type/name specified")
			return
		} else {
			b.Sname = val
		}
	}
	if b.Code == 0{
		b.Code = 1
	}
	if b.Nsecs == 0 || b.Nsecs > 10{b.Nsecs = 3}
	return
}

//TableOlde generates an ascii table report for a Bm
func (b *Bm) TableOlde(printz bool) {
	if b.Title == "" {
		if b.Id == 0 {
			b.Id = rand.Intn(666)
		}
		b.Title = fmt.Sprintf("stl_bm_%v", b.Id)
	}
	rezstr := new(strings.Builder)
	hdr := fmt.Sprintf("%s\nsteel beam report\ndate-%s\n%s\n%s\n", ColorYellow, time.Now().Format("2006-01-02"), b.Title, ColorReset)
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)
	table.SetCaption(true, "beam properties")
	table.SetHeader([]string{"grade", "section type", "span(m)", "dtyp(1-ss,2-frm)", "end condition"})
	row := fmt.Sprintf("%v, %s, %.3f, %v, %v", b.Grd, stlsecmap[b.Styp], b.Lspan,b.Dtyp, b.Endc)
	table.Append(strings.Split(row, ","))
	table.Render()
	if len(b.Ldcases) > 0{
		rezstr.WriteString(ColorRed)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true, "beam loads")
		table.SetHeader([]string{"load no.","load type","wa(kn)","wb(kn)","la(m)","lb(m)","dead(1)/imposed(2)"})
		for i, ldcase := range b.Ldcases{
			row = fmt.Sprintf("%v, %.f, %.2f, %.2f, %.2f, %.2f, %.f",i+1, ldcase[1],ldcase[2],ldcase[3],ldcase[4],ldcase[5],ldcase[6])
			table.Append(strings.Split(row,","))
		}
		table.Render()
	}
	if b.Dz{
		rezstr.WriteString(ColorPurple)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true, "section geometry")
		table.SetHeader([]string{"section", "wt\n(kg/m)", "depth\n(mm)", "t.web\n(mm)", "area\n(cm2)", "rxx\n(cm)", "ryy\n(cm)", "zxx\n(cm3)", "zyy\n(cm3)"})
		/*

		 */
		df, _ := kass.GetStlDf(b.Sname)
		for _, idx := range b.Rez{
			//fa, pa, px, py, fp, mx, my, sx, sy, s1, dtrat := b.Vals[i]
			sstr, wt, dw, tw, ar, rxx, ryy, zxx, zyy := df.Elem(idx, 1), df.Elem(idx, 2).Float(), df.Elem(idx, 3).Float(), df.Elem(idx, 6).Float(), df.Elem(idx, 23).Float(), df.Elem(idx, 13).Float(), df.Elem(idx, 14).Float(), df.Elem(idx, 15).Float(), df.Elem(idx, 16).Float()
			row = fmt.Sprintf("%s,%.3f,%.f,%.f,%.f,%.f,%.f,%.f,%.f", sstr, wt, dw, tw, ar, rxx, ryy, zxx, zyy)
			table.Append(strings.Split(row, ","))
		}
		table.Render()
		rezstr.WriteString(ColorGreen)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true, "section results")
		//vals = []float64{b.Vu, b.Mu, b.Dmax, px, fm, ps, fs, b.Lspan*1000.0/360., defrat, pc, fc, pb, fb}
		table.SetHeader([]string{"section","shear\n(kn)", "moment\n(kn-m)", "def.\n(mm)", "def.\nperm.(mm)","px\n(n/mm2)", "fx\n(n/mm2)", "ps\n(n/mm2)", "fs\n(n/mm2)", "pc\n(n/mm2)", "fc\n(n/mm2)", "pb\n(n/mm2)", "fb\n(n/mm2)"})
		for i, idx := range b.Rez{
			sstr := df.Elem(idx,1)
			//b.Vu, b.Mu, b.Dmax, b.Lspan*1000./360.,px, fm, ps, fs, pc, fc, pb, fb
			vu, mu, dmax, dperm, px, fx, ps, fs, pc, fc, pb, fb := b.Vals[i][0],b.Vals[i][1],b.Vals[i][2],b.Vals[i][3],b.Vals[i][4],b.Vals[i][5],b.Vals[i][6],b.Vals[i][7],b.Vals[i][8],b.Vals[i][9],b.Vals[i][10],b.Vals[i][11]
			row = fmt.Sprintf("%s,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f", sstr,vu, mu, dmax, dperm, px, fx, ps, fs, pc, fc, pb, fb)
			table.Append(strings.Split(row, ","))
		}
		table.Render()
		rezstr.WriteString(ColorCyan)
		if b.Kostin == 0.0{b.Kostin = 200.0}
		table = tablewriter.NewWriter(rezstr)
		wt := b.Vals[b.Mindx][12]
		minsec := df.Elem(b.Rez[b.Mindx],1)
		//fmt.Println("mindx",b.Mindx,df.Elem(b.Rez[b.Mindx],1),b.Vals[b.Mindx][12])
		table.SetCaption(true, "quantity take off")
		table.SetHeader([]string{"section","min. wt\n(kg)","cost\n(rs)","span\n(m)","total cost\n(rs)"})
		row = fmt.Sprintf("%s, %.3f, %.3f, %.3f, %.3f",minsec,wt,b.Kostin,b.Lspan,b.Kostin * wt*b.Lspan)
		table.Append(strings.Split(row, ","))
		table.Render()
		rezstr.WriteString(ColorReset)
	}
	b.Report = rezstr.String()
	if printz{
		fmt.Println(b.Report)
	}

}


//TableOlde generates an ascii table report for a Bm
func (b *Bm) Table(printz bool) {
	dtmap := map[int]string{
		0:"check",
		1:"ss beam",
		2:"w/end moments",
		3:"purlin",
	}
	if b.Title == "" {
		if b.Id == 0 {
			b.Id = rand.Intn(666)
		}
		b.Title = fmt.Sprintf("stl_bm_%v", b.Id)
	}
	rezstr := new(strings.Builder)
	hdr := fmt.Sprintf("%s\nsteel beam report\ndate-%s\n%s\n%s\n", ColorYellow, time.Now().Format("2006-01-02"), b.Title, ColorReset)
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)
	table.SetCaption(true, "beam properties")
	table.SetHeader([]string{"grade", "section name", "span(mm)", "dtyp", "end condition"})
	row := fmt.Sprintf("%v, %s, %.3f, %v, %v", b.Grd, b.Sname, b.Lspan, dtmap[b.Dtyp], b.Endc)
	table.Append(strings.Split(row, ","))
	table.Render()
	if len(b.Ldcases) > 0{
		rezstr.WriteString(ColorRed)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true, "beam loads")
		table.SetHeader([]string{"load no.","load type","wa(kn)","wb(kn)","la(m)","lb(m)","dead(1)/imposed(2)"})
		for i, ldcase := range b.Ldcases{
			row = fmt.Sprintf("%v, %.f, %.2f, %.2f, %.2f, %.2f, %.f",i+1, ldcase[1],ldcase[2],ldcase[3],ldcase[4],ldcase[5],ldcase[6])
			table.Append(strings.Split(row,","))
		}
		table.Render()
	}
	if b.Prln{
		rezstr.WriteString(ColorRed)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true, "purlin loads")
		table.SetHeader([]string{"angle(deg)","spacing(mm)","DL(n/m2)","LL(n/m2)","WL(n/m2)"})
		row = fmt.Sprintf("%.2f, %.2f, %.2f, %.2f, %.2f",b.Theta,b.Ly,b.DL,b.LL,b.WL)
		table.Append(strings.Split(row,","))
		table.Render()
		var mwt float64
		var mdx int
		if b.Dz{	
			rezstr.WriteString(ColorPurple)
			table = tablewriter.NewWriter(rezstr)
			table.SetCaption(true, "section geometry")
			table.SetHeader([]string{"no.","section", "wt\n(kg/m)", "depth\n(mm)", "t.web\n(mm)", "area\n(cm2)", "rxx\n(cm)", "ryy\n(cm)", "sxx\n(cm3)", "syy\n(cm3)"})
			for i, s := range b.Ssecs{
				row = fmt.Sprintf("%v, %s, %.2f, %.2f, %.2f, %.2f, %.2f, %.2f, %.2f, %.2f",i+1,s.Sstr,s.Wt,s.H,s.Tw,s.Area,s.Rxx,s.Ryy,s.Sxx,s.Syy)
				if mwt == 0.0 || mwt > s.Wt{
					mwt = s.Wt
					mdx = i
				}
			}
			table.Append(strings.Split(row, ","))			
		}
		table.Render()
		rezstr.WriteString(ColorCyan)
		if b.Kostin == 0.0{b.Kostin = 100.0}
		table = tablewriter.NewWriter(rezstr)
		smin := b.Ssecs[mdx]
		//fmt.Println("mindx",b.Mindx,df.Elem(b.Rez[b.Mindx],1),b.Vals[b.Mindx][12])
		table.SetCaption(true, "quantity take off")
		table.SetHeader([]string{"section","min. wt\n(kg)","cost\n(rs)","span\n(m)","total cost\n(rs)"})
		row = fmt.Sprintf("%s, %.2f, %.2f, %.2f, %.2f",smin.Sstr,mwt,b.Kostin,b.Lspan,b.Kostin * mwt*b.Lspan)
		table.Append(strings.Split(row, ","))
		table.Render()
		rezstr.WriteString(ColorReset)
		return

	}
	if b.Dz{
		rezstr.WriteString(ColorPurple)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true, "section geometry")
		table.SetHeader([]string{"section", "wt\n(kg/m)", "depth\n(mm)", "t.web\n(mm)", "area\n(cm2)", "rxx\n(cm)", "ryy\n(cm)", "zxx\n(cm3)", "zyy\n(cm3)"})
		/*

		 */
		df, _ := kass.GetStlDf(b.Sname)
		for _, idx := range b.Rez{
			//fa, pa, px, py, fp, mx, my, sx, sy, s1, dtrat := b.Vals[i]
			sstr, wt, dw, tw, ar, rxx, ryy, zxx, zyy := df.Elem(idx, 1), df.Elem(idx, 2).Float(), df.Elem(idx, 3).Float(), df.Elem(idx, 6).Float(), df.Elem(idx, 23).Float(), df.Elem(idx, 13).Float(), df.Elem(idx, 14).Float(), df.Elem(idx, 15).Float(), df.Elem(idx, 16).Float()
			row = fmt.Sprintf("%s,%.3f,%.f,%.f,%.f,%.f,%.f,%.f,%.f", sstr, wt, dw, tw, ar, rxx, ryy, zxx, zyy)
			table.Append(strings.Split(row, ","))
		}
		table.Render()
		rezstr.WriteString(ColorGreen)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true, "section results")
		//vals = []float64{b.Vu, b.Mu, b.Dmax, px, fm, ps, fs, b.Lspan*1000.0/360., defrat, pc, fc, pb, fb}
		table.SetHeader([]string{"section","shear\n(kn)", "moment\n(kn-m)", "def.\n(mm)", "def.\nperm.(mm)","px\n(n/mm2)", "fx\n(n/mm2)", "ps\n(n/mm2)", "fs\n(n/mm2)", "pc\n(n/mm2)", "fc\n(n/mm2)", "pb\n(n/mm2)", "fb\n(n/mm2)"})
		for i, idx := range b.Rez{
			sstr := df.Elem(idx,1)
			//b.Vu, b.Mu, b.Dmax, b.Lspan*1000./360.,px, fm, ps, fs, pc, fc, pb, fb
			vu, mu, dmax, dperm, px, fx, ps, fs, pc, fc, pb, fb := b.Vals[i][0],b.Vals[i][1],b.Vals[i][2],b.Vals[i][3],b.Vals[i][4],b.Vals[i][5],b.Vals[i][6],b.Vals[i][7],b.Vals[i][8],b.Vals[i][9],b.Vals[i][10],b.Vals[i][11]
			row = fmt.Sprintf("%s,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f", sstr,vu, mu, dmax, dperm, px, fx, ps, fs, pc, fc, pb, fb)
			table.Append(strings.Split(row, ","))
		}
		table.Render()
		rezstr.WriteString(ColorCyan)
		if b.Kostin == 0.0{b.Kostin = 200.0}
		table = tablewriter.NewWriter(rezstr)
		wt := b.Vals[b.Mindx][12]
		minsec := df.Elem(b.Rez[b.Mindx],1)
		//fmt.Println("mindx",b.Mindx,df.Elem(b.Rez[b.Mindx],1),b.Vals[b.Mindx][12])
		table.SetCaption(true, "quantity take off")
		table.SetHeader([]string{"section","min. wt\n(kg)","cost\n(rs)","span\n(m)","total cost\n(rs)"})
		row = fmt.Sprintf("%s, %.3f, %.3f, %.3f, %.3f",minsec,wt,b.Kostin,b.Lspan,b.Kostin * wt*b.Lspan)
		table.Append(strings.Split(row, ","))
		table.Render()
		rezstr.WriteString(ColorReset)
	}
	b.Report = rezstr.String()
	if printz{
		fmt.Println(b.Report)
	}

}


//GetSec returns a StlSec for a beam
func (b *Bm) GetSec()(ss kass.StlSec, err error){
	var bf, dt, tf, tw float64
	switch b.Sname{
		case "plate-i":
		bf = b.Params[0]
		dt = b.Params[1]
		tf = b.Params[2]
		tw = b.Params[3]
		case "l2-ss","l2-os","ln2-ss","ln2-os":
		bf = b.Params[0]
		case "built-i":
		bf = b.Params[0]
		dt = b.Params[1]	
	}
	ss, err = kass.GetStlSec(b.Sname, b.Sdx, b.Code, bf, dt, tf, tw)
	if err != nil{
		return
	}
	ss.Lsb = b.Lsb
	ss.Lspan = b.Lspan
	ss.Lbr = b.Lbr
	ss.Tbr = b.Tbr
	return
}

//Chk449 checks a beam section as per bs449
//as in section 6.2, mosley/spencer
//i is the df section index
func (b *Bm) Chk449(df dataframe.DataFrame, i int) (vals []float64, chk bool) {
	//only for bs449
	//many such checks to be written lol :(
	var b1, b2, fm, fs, fc, fb, y0, q4, q5 float64
	//check for bending stress
	log.Println("checking section->",df.Elem(i,1))
	pvec := pqbm[b.Grd]
	qx, ps, pc := pvec[0], pvec[1], pvec[2]
	log.Println("section modulus->",df.Elem(i,15))
	if b.Mu*1000.0/df.Elem(i, 15).Float() > qx {
		return
	}
	dtrat := df.Elem(i, 3).Float() / df.Elem(i, 6).Float()
	if dtrat < 5.0 {
		dtrat = 5.0
	}
	if b.Spam{log.Println("dtrat->", dtrat)}

	sdrat := b.Ly * 100.0 / df.Elem(i, 14).Float()
	
	if b.Spam{log.Println("sdrat->", dtrat)}
	var px float64
	if b.Yeolde{
		px = PbcYeolde(sdrat, dtrat)
	} else {
		px = PbcLerp(b.Sname, b.Grd, sdrat, dtrat)
	}
	fm = math.Abs(b.Mu) * 1e3 / df.Elem(i, 15).Float()
	if fm/px > 1.0 {
		return
	}
	if b.Spam{log.Println("section->", df.Elem(i, 1), ColorBlue, "bending o.k", ColorReset)}
	//check for shear stress
	fs = math.Abs(b.Vu) * 1e3 / df.Elem(i, 3).Float() / df.Elem(i, 5).Float()
	if fs/ps > 1.0 {
		return
	}

	if b.Spam{log.Println("section->", df.Elem(i, 1), kass.ColorBlue, "shear o.k", kass.ColorReset)}
	//check for deflection
	//dmax := b.Dmax
	defrat := math.Round(b.Lspan * 1000. / b.Dmax)
	if b.Spam{log.Println("deflection ->", b.Dmax, "mm vs perm. ->", b.Lspan*1000./360., "perm. ratio - 360.0, actual ->", defrat)}
	if b.Dmax > b.Lspan*1000./360.{
		return
	}
	
	//log.Println("section->",df.Elem(i,1),kass.ColorBlue,"deflection o.k",kass.ColorReset)
	//check for web crushing stress
	b1 = b.Lbr + (b.Tbr+0.5*(df.Elem(i, 3).Float()-df.Elem(i, 8).Float()))*math.Sqrt(3.0)
	fc = 1e3 * math.Abs(b.Vu) / df.Elem(i, 5).Float() / b1
	if b.Spam{log.Println("web crushing-> perm ",pc," vs ",fc,"actual (n/mm)")}
	if fc/pc > 1.0 {
		return
	}
	//check for web buckling stress
	b2 = b.Lbr + b.Tbr + df.Elem(i, 3).Float()/2.0
	fb = 1e3 * math.Abs(b.Vu) / df.Elem(i, 5).Float() / b2
	sweb := math.Sqrt(3.0) * df.Elem(i, 8).Float() / df.Elem(i, 5).Float()
	c0 := math.Pow(math.Pi, 2) * EStl / math.Pow(sweb, 2)
	n0 := 0.3 * math.Pow(sweb/100.0, 2)
	switch {
	case b.Grd == 43:
		y0 = 250.0
		q4 = 155.0
		q5 = 143.0
	case b.Grd == 50:
		y0 = 350.0
		q4 = 215.0
		q5 = 200.0
	case b.Grd == 55:
		y0 = 430.0
		q4 = 265.0
		q5 = 245.0
	}
	a0 := (y0 + c0*(n0+1.0)) / 2.0
	pb := (a0 - math.Sqrt((math.Pow(a0, 2) - y0*c0))) / 1.7
	if sweb <= 30 {
		fb = q4 - (q4-q5)*sweb/30.0
	}
	if b.Spam{log.Println("web buckling-> perm ",pb," vs ",fb,"actual (n/mm)")}
	if fb/pb > 1.0{
		return
	}
	chk = true
	//vals = append(maxs, []float64{}...)
	wt := df.Elem(i, 2).Float() 
	vals = []float64{b.Vu, b.Mu, b.Dmax, b.Lspan*1000./360.,px, fm, ps, fs, pc, fc, pb, fb, wt}
	// 0 Vu, 1 b.Mu, 2 b.Dmax,3 b.Lspan*1000./360.,4 px, 5 fm, 6 ps, 7 fs, 8 pc, 9 fc, 10 pb, 11 fb, 12 wt
	return
}

//GenMod generates a kass.Model for analysis
func (b *Bm) GenMod(ss kass.StlSec)(mod kass.Model, err error){
	mod.Frmstr = "1db"
	mod.Calc = true
	mod.Ncjt = 2
	mod.Id = b.Title
	mod.Em = [][]float64{{ss.Em}}
	mod.Cp = [][]float64{{ss.Ixx, ss.Area}}
	mod.Term = b.Term
	mod.Web = b.Web
	mod.Units = "nmm"
	switch b.Dtyp{
		case -1:
		//cantilever beam
		mod.Coords = [][]float64{{0,0},{b.Lspan,0}}
		mod.Supports = [][]int{{1, -1, -1}}
		mod.Mprp = [][]int{{1,2,1,1,0}}
		case 1:
		//ss beam
		mod.Coords = [][]float64{{0,0},{b.Lspan,0}}
		mod.Supports = [][]int{{1, -1, 0},{2,-1,0}}
		mod.Mprp = [][]int{{1,2,1,1,0}}
		case 2:
		//cs beam
		var x float64
		mod.Coords = append(mod.Coords, []float64{x,0})
		mod.Supports = append(mod.Supports, []int{1, -1, 0})
		for i, span := range b.Lspans{
			x += span
			mod.Coords = append(mod.Coords, []float64{x,0})
			mod.Mprp = append(mod.Mprp, []int{i+1, i+2, 1, 1, 0})
			mod.Supports = append(mod.Supports, []int{i+2, -1, 0})
		}
	}
	if len(mod.Mprp) == 0{
		err = fmt.Errorf("error generating model(wrong spans?) lspan %f lspans %f",b.Lspan, b.Lspans)
		return
	}
	//gen loads
	mod.Msloads = append(mod.Msloads, b.Ldcases...)

	if b.DL > 0.0{
		dl := b.DL
		if b.Selfwt{
			dl += ss.Wt/1e3
		}
		for i := range mod.Mprp{	
			mod.Msloads = append(mod.Msloads, []float64{
				float64(i+1), 3.0, b.DL, 0.0, 0.0, 0.0, 1.0,
			})
		}
	}
	if b.LL > 0.0{
		for i := range mod.Mprp{
			mod.Msloads = append(mod.Msloads, []float64{
				float64(i+1), 3.0, b.LL, 0.0, 0.0, 0.0, 2.0,
			})
		}
	}
	return
}

func (b *Bm) ChkSec()(err error){
	var ss kass.StlSec
	ss, err = b.GetSec()
	if err != nil{
		return
	}
	//fmt.Println("checking sec-",ss.Sstr)
	//uvals are max shear, bending moment, deflection
	if b.Dtyp != 0{
		var bmrez map[int]kass.BeamRez
		bmrez, err = b.Calc1d(ss)
		if err != nil{
			return
		}
		err = b.SortRez(bmrez)
		if err != nil{
			return
		}
	}
	ss.Vu = b.Vu; ss.Mu = b.Mu; ss.Dmax = b.Dmax
	switch b.Code{
		case 1:
		err = ss.BmChk800()
		case 2:
		err = ss.BmChk449()
	}
	if err != nil{b.Ssecs = append(b.Ssecs, ss)}
	return
}

func (b *Bm) SortRez(bmrez map[int]kass.BeamRez)(err error){
	for i:= 1; i <= len(bmrez); i++{
		bm := bmrez[i]
		if i == 1 || b.Vu < bm.Maxs[0]{
			b.Vu = bm.Maxs[0]
		}
		if i == 1 || b.Mu < bm.Maxs[1]{
			b.Mu = bm.Maxs[1]
		}
		if i == 1 || b.Dmax < bm.Maxs[2]{
			b.Dmax = bm.Maxs[2]
		}
	}
	if b.Vu <= 0.0 || b.Mu <= 0.0 || b.Dmax <= 0.0{
		err = fmt.Errorf("error in results - SF %f BM %f Def %f",b.Vu, b.Mu, b.Dmax)
	}
	return
}

//GetUvalSs gets the max design values for a simply supported steel beam
func (b *Bm) GetUvalSs(df dataframe.DataFrame, sdx int){
	loads := b.Ldcases
	var wt, area, iz float64
	//bs code steel sheet (change this)
	
	wt = df.Elem(sdx, 2).Float() * 9.81 / 1e3
	area = df.Elem(sdx, 23).Float() * 1e-2
	iz = df.Elem(sdx, 11).Float() * 1e-4

	//log.Println("wt, area, iz-",wt,area,iz)
	loads = append(loads, []float64{float64(b.Id), 3, wt, 0, 0, 0, 1})
	maxs := SimpBmFrc(1, loads, b.Lspan, EStl, area, iz)
	b.Vu = math.Abs(maxs[0]); b.Mu = math.Abs(maxs[1]); b.Dmax = maxs[2]
	if math.Abs(maxs[3]) > b.Mu{b.Mu = math.Abs(maxs[3])}
	return 
}

//Calc1d calculates uvals for a beam section
func (b *Bm) Calc1d(ss kass.StlSec)(bmrez map[int]kass.BeamRez, err error){
	var mod kass.Model
	mod, err = b.GenMod(ss)
	if err != nil{
		return
	}
	frmrez, e := kass.CalcBm1d(&mod, 2)
	if e != nil{
		err = fmt.Errorf("beam analysis error - %s",err)
		return
	}
	bmrez = kass.CalcBmSf(&mod, frmrez,b.Verbose)
	return
}

func BmDz(b *Bm) (err error){
	err = b.Init()
	if err != nil{
		return
	}
	ndx := kass.StlSdxLims[b.Sname]
	if ndx == 0{
		err = fmt.Errorf("%s design functions not written",b.Sname)
		return
	}
	if b.Sdx > 0{ndx = b.Sdx}
	for idx := ndx; idx >= 0; idx--{
		if len(b.Sdxs) == b.Nsecs{
			break
		}
		//fmt.Println("checking ndx-", idx)
		b.Sdx = idx
		err = b.ChkSec()
		if err == nil{
			b.Sdxs = append(b.Sdxs, idx)
			
		}
	}
	return
}

//PrlnDz800 designs a purlin section as per duggal ex 9.9(is800)
func PrlnDz800(nsecs int, sname string, tspc, pspc, theta, pdl, pll, pwl float64)(sss []kass.StlSec, err error){
	//NOTE pdl,pll, pwl are in n/m2
	ndx := kass.StlSdxLims[sname]
	for idx := ndx; idx > 0; idx--{
		ss, e := kass.GetStlSec(sname, idx, 1)
		if e != nil{
			//fmt.Println(err)
			err = fmt.Errorf("index error %s",e)
			return
		}
		if len(sss) >= nsecs{break}
		p := (pdl + pll) * pspc/1000.0 + ss.Wt
		pn := p * math.Cos(theta * math.Pi/180.0)
		pp := p * math.Sin(theta * math.Pi/180.0)
		pn  += pwl * pspc/1000.0
		mux := 1.5 * pn * tspc * tspc/10.0/1000.0
		muy := 1.5 * pp * tspc * tspc/10.0/1000.0		
		ss.Ax = -2
		ss.Lspan = tspc
		ss.Mux = mux
		ss.Muy = muy
		ss.Lsb = true
		ss.Dmax = 5.0 * pn * math.Pow(tspc, 4)/ss.Em/ss.Ixx/384.0/1000.0
		ss.Vux = 0.6 * pn * tspc
		ss.Vuy = 0.6 * pp * tspc
		err = ss.BmChk800()
		if err == nil{
			fmt.Println("section phound-")
			fmt.Println(ss.Sstr)
			fmt.Println(idx)
			ss.Printz()
			sss = append(sss, ss)
		}
	}
	if len(sss) == 0{
		err = fmt.Errorf("no suitable section found")
	}
	return
}

//BmDesign designs a steel beam 
func BmDesign(b *Bm) (err error) {
	//return
	// if b.Dtyp > 0{
	// 	err = b.Calc1d()
	// }
	switch{
		case b.Prln:
		b.Init()
		sss, e := PrlnDz800(b.Nsecs, b.Sname, b.Lspan, b.Ly, b.Theta, b.DL, b.LL, b.WL)
		if e != nil{
			err = fmt.Errorf("error designing purlin section %s",e)
			return
		}
		b.Ssecs = sss
		b.Table(b.Verbose)
		default:
		err = BmDz(b)
		if err != nil{
			return
		}
		b.Table(b.Verbose)
	}
	return
}


//SimpBmFrc returns max. design values for a simply supported steel beam given a slice of load cases
//and geom data - lspan, elastic modulus, area, moment of inertia
func SimpBmFrc(member int, ldcases [][]float64, lspan, e, area, ix float64) []float64 {
	//memrel = 3 for hinge at both ends
	qfchn := make(chan []interface{}, 1)
	var rl, rr float64
	geoms := []float64{ix}
	for _, ldcase := range ldcases {
		ltyp := int(ldcase[1])
		memrel := 3
		go kass.FxdEndFrc(member, memrel, ltyp, lspan, ldcase[2:], 0, geoms, qfchn)
		r := <-qfchn
		qf, _ := r[1].([]float64)
		rl += qf[0]
		rr += qf[2]
	}
	//log.Println("reactions->", kass.ColorRed, rl, rr, kass.ColorReset)
	rez := kass.Bmsfcalc(1, ldcases, lspan, e, area, ix, rl, 0, rr, 0, false, false)
	//member int, ldcases [][]float64, l, e, a, iz, rl, ml, re, me float64, plotbm bool
	return rez.Maxs
}

//StlBmDBs designs a steel beam by iterating over stype df
//again, basically section 6.2 mosley/spencer
//old func, ignore
func StlBmDBs(lspan, ly, ty, lbr, tbr float64, ldcases [][]float64, sname string, grd, nsecs int, brchck, yeolde bool) []int {
	memid := ldcases[0][0]
	mem := int(memid)
	ldtyps := make(map[float64][][]float64)
	//how stupid. how dumb. ldtyp 1 = dl, ltyp 2 = ll. instead why this 99.9999 nonsense
	for _, ldcase := range ldcases {
		//add PSFs
		ldtyps[ldcase[6]] = append(ldtyps[ldcase[6]], ldcase)
		ldtyps[0.] = append(ldtyps[0.], ldcase)
		if ldcase[6] != 1.0 {
			ldtyps[99.9] = append(ldtyps[99.9], ldcase)
		}
	}
	df, err := kass.GetStlDf(sname)
	if err != nil {
		log.Println("ERRORE,errore->", err)
	}
	var rez []int
	var b1, b2, fm, fs, fc, fb, y0, q4, q5 float64
	pvec := pqbm[grd]
	qx, ps, pc := pvec[0], pvec[1], pvec[2]
	//log.Println("pvec->",kass.ColorCyan, qx, ps, pc)
	//UNNITTTTZZZZZZ!!!!!
	for i := df.Nrow() - 1; i > 0; i-- {
		//log.Println("checking section->",df.Elem(i,1))
		if len(rez) == nsecs {
			break
		}
		wt := df.Elem(i, 2).Float() * 9.81 / 1e3
		allloads := ldtyps[0.]
		allloads = append(allloads, []float64{memid, 3, wt, 0, 0, 0, 1})
		maxs := SimpBmFrc(mem, allloads, lspan, EStl, df.Elem(i, 23).Float()*1e2, df.Elem(i, 11).Float()*1e4)
		//check for bending stress
		if maxs[1]*1e3/df.Elem(i, 15).Float() > qx {
			continue
		}
		dtrat := df.Elem(i, 3).Float() / df.Elem(i, 6).Float()
		if dtrat < 5.0 {
			dtrat = 5.0
		}
		sdrat := ly * 100.0 / df.Elem(i, 14).Float()
		var px float64
		if yeolde {
			px = PbcYeolde(sdrat, dtrat)
		} else {
			px = PbcLerp(sname, grd, sdrat, dtrat)
		}
		fm = math.Abs(maxs[1]) * 1e3 / df.Elem(i, 15).Float()
		if fm/px > 1.0 {
			continue
		}
		//log.Println("section->",df.Elem(i,1),kass.ColorBlue,"bending o.k",kass.ColorReset)
		//check for shear stress
		fs = math.Abs(maxs[0]) * 1e3 / df.Elem(i, 3).Float() / df.Elem(i, 5).Float()
		if fs/ps > 1.0 {
			continue
		}

		//log.Println("section->",df.Elem(i,1),kass.ColorBlue,"shear o.k",kass.ColorReset)
		//check for deflection
		dmaxs := SimpBmFrc(mem, ldtyps[99.9], lspan*1e3, EStl, df.Elem(i, 23).Float()*1e2, df.Elem(i, 11).Float()*1e4)
		dmax := dmaxs[2] / 1e3
		//defrat := math.Round(lspan*1000./dmax)
		//log.Println("deflection->",dmax,"mm vs ->", lspan*1000./360., "actual->", defrat)
		if dmax > lspan*1000./360. {
			continue
		}

		//log.Println("section->",df.Elem(i,1),kass.ColorBlue,"deflection o.k",kass.ColorReset)
		//check for web crushing stress
		b1 = lbr + (tbr+0.5*(df.Elem(i, 3).Float()-df.Elem(i, 8).Float()))*math.Sqrt(3.0)
		fc = 1e3 * math.Abs(maxs[0]) / df.Elem(i, 5).Float() / b1
		
		if fc/pc > 1.0 {
			continue
		}
		//check for web buckling stress
		b2 = lbr + tbr + df.Elem(i, 3).Float()/2.0
		fb = 1e3 * math.Abs(maxs[0]) / df.Elem(i, 5).Float() / b2
		sweb := math.Sqrt(3.0) * df.Elem(i, 8).Float() / df.Elem(i, 5).Float()
		c0 := math.Pow(math.Pi, 2) * EStl / math.Pow(sweb, 2)
		n0 := 0.3 * math.Pow(sweb/100.0, 2)
		switch {
		case grd == 43:
			y0 = 250.0
			q4 = 155.0
			q5 = 143.0
		case grd == 50:
			y0 = 350.0
			q4 = 215.0
			q5 = 200.0
		case grd == 55:
			y0 = 430.0
			q4 = 265.0
			q5 = 245.0
		}
		a0 := (y0 + c0*(n0+1.0)) / 2.0
		pb := (a0 - math.Sqrt((math.Pow(a0, 2) - y0*c0))) / 1.7
		if sweb <= 30 {
			fb = q4 - (q4-q5)*sweb/30.0
		}
		if fb/pb > 1.0 {
			continue
		}
		rez = append(rez, i)
		log.Println(kass.ColorGreen, "section->", df.Elem(i, 1), kass.ColorReset)
		log.Println(kass.ColorCyan, "fb/pb->", fm/px, kass.ColorReset)
		log.Println("bending stress (px, fm)->", px, fm, fm/px)
		log.Println("shear stress (ps, fs)->", ps, fs, fs/ps)
		log.Println("web crushing stress->", pc, fc, fc/pc)
		log.Println("web buckling stress ->", pb, fb, fb/pb)
		log.Println("***\n***")
	}
	return rez
}


//BmDzOlde is yeolde entry func for steel beam design
//TODO - write continuous beam calcs
func BmDzOlde(b *Bm) (err error){
	b.Mindx = -1
	df, err := kass.GetStlDf(b.Sname)
	if err != nil {
		return
	}
	if b.Grd == 0 {
		b.Grd = 43
	}
	for i := df.Nrow() - 1; i > 0; i--{
		//log.Println("checking section->",df.Elem(i,1))
		if len(b.Rez) == b.Nsecs{
			break
		}
		switch b.Dtyp{
		case 0:
		//end moments specified in beam struct
		case 1:
			//ss beam
			switch b.Endc{
			case 0:
			//write this
			case 1:
				b.GetUvalSs(df, i)
				//log.Println(maxs)
				switch b.Code{
				case 2:
					if val, ok := b.Chk449(df, i); ok {
						//log.Println("maxs->", maxs)
						b.Rez = append(b.Rez, i)
						b.Vals = append(b.Vals, val)
						if b.Mindx == -1 || (b.Vals[b.Mindx][12] > val[12]){
							b.Mindx = len(b.Rez)-1
						}
						// if b.Verbose {
						// 	//do nuttin
						// }
					}
				}
			}
		case 2:
			//cs beam
		}
	}
	if len(b.Rez) == 0{
		err = errors.New("no suitable section found")
		return
	}
	/*
	sort.Slice(b.Rez, func(i,j int) bool{
		return b.Vals[i][12] < b.Vals[j][12]
	})
	
	sort.Slice(b.Vals, func(i,j int) bool{
		return b.Vals[i][12] < b.Vals[j][12]
	})
	*/
	b.Dz = true
	b.TableOlde(b.Verbose)
	return
}

