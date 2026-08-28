package barf

import (
	"os"
	"fmt"
	"log"
	"math"
	"errors"
	"time"
	"strings"
	"runtime"
	"math/rand"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"github.com/olekukonko/tablewriter"
	"github.com/go-gota/gota/dataframe"
	kass"barf/kass"
)

var (
	//table2bs = []float64{180,165,230,215,280,170,155,215,200,265,185}	
	//stlsecmap = map[int]string{7:"UB",8:"UC"}
	//pqcol - from mosley section 6
	pqcol = map[int][]float64{43:{155,165,165},50:{215,230,230},55:{265,280,280}}
	//pqbm - qx (bending), ps (shear), pc (web crushing)
	pqbm = map[int][]float64{43:{165,100,190},50:{230,140,260},55:{280,170,320}}
	//pqcol89 = map[int][]float64{43:{170,180,180},50:{215,230,230},55:{265,280,280}}
	EStl = 210000.0
	EStl89 = 205000.0
	CbsPlts = []float64{10,12,16,20,25,30}
	
)

//Col is a steel column struct
//mosley spencer section 6.1
type Col struct{
	Title      string
	Sname      string
	Bstyp      string
	Term       string
	Frmstr     string
	Catstr     string
	Label      string
	Id         int
	H1, H2, Lx, Ly, Tx, Ty, Mx, My, Vx, Vy, Pu, Pfac float64
	Lspan      float64
	Max, Mbx   float64
	May, Mby   float64
	Vax, Vbx   float64
	Vay, Vby   float64
	Grd, Styp  int
	Dtyp       int
	Bmtyp      int
	Calctyp    int
	Nsecs, Sdx int
	Code, Endc int
	Deg        int
	Cleg       int  //connected leg index (1-longer, 2-shorter leg)
	Rez        []int
	Sdxs       []int
	Vals       [][]float64
	Bmx        []string
	Bmy        []string
	Ljbase     bool
	Web        bool
	Yeolde     bool
	Verbose    bool
	Blurb      bool
	Spam       bool
	Braced     bool
	Dz         bool
	Kondz      bool //design konnection
	Dsgn       bool //if false, check section
	Tz         bool //tensile member checks
	Tie        bool //is a tensile member/tie
	Frame      bool //is a frame member
	Store      bool //store stlsec
	Weld       bool
	Onism      bool //only ismb sections
	Flip       bool //flip/rotate by 90
	Ignore     bool
	Ctyp       int  //bolted/welded end connection
	Report     string
	Txtplots   []string
	Txtplot    string
	Name       string
	Kostin     float64
	Mindx      int
	Ssec       kass.StlSec
	Ssecs      []kass.StlSec
	Bg         kass.Blt
	Wg         kass.Wld
	Bcon, Tcon int
	Params     []float64
}

//Table generates an ascii table report for a Col
func (c *Col) Table(printz bool){
	if c.Title == ""{
		if c.Id == 0{
			c.Id = rand.Intn(666)
		}
		c.Title = fmt.Sprintf("stl_col_%v",c.Id)
	}
	rezstr := new(strings.Builder)
	hdr := fmt.Sprintf("%s\nsteel column report\ndate-%s\n%s\n%s\n",ColorYellow,time.Now().Format("2006-01-02"),c.Title,ColorReset)
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)
	var row string
	table.SetCaption(true,"column properties")
	table.SetHeader([]string{"grade","section type","height(above)(mm)","height(col)(mm)","unb.len(lx)(mm)","unb.len(ly)(mm)","tx","ty","span(mm)"})
	row = fmt.Sprintf("%v, %s, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f",c.Grd,c.Sname,c.H2,c.H1,c.Lx,c.Ly,c.Tx,c.Ty,c.Lspan)
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorRed)
	table = tablewriter.NewWriter(rezstr)
	
	table.SetCaption(true,"ultimate loads")
	switch c.Code{
		case 0,1:
		//is800
		switch c.Dtyp{
			case 0:
			//axially loaded column
			table.SetHeader([]string{"code","dtyp","axial load(kn)"})
			row = fmt.Sprintf("is 800, col, %.3f",c.Pu/1e3)
			table.Append(strings.Split(row,","))
			table.Render()
			
			case 1:
			//beam-col
			table.SetHeader([]string{"code","dtyp(0-b/1-m)","axial load(kn)","mx bot(knm)","mx top(knm)","my top(knm)","my bot(knm)","vdx top(kn)","vdx bot(kn)","vdy top(kn)","vdy bot(kn)"})
			row = fmt.Sprintf("is 800, bm-col, %.3f, %.3f, %.3f, %.3f, %.3f,%.3f, %.3f, %.3f, %.3f",c.Pu/1e3, c.Max/1e6, c.Mbx/1e6, c.May/1e6, c.Mby/1e6,c.Vax/1e3, c.Vbx/1e3, c.Vay/1e3, c.Vby/1e3)
			table.Append(strings.Split(row,","))
			table.Render()
			
		}
		case 2:
		//bs449
		table.SetHeader([]string{"code","dtyp(0-b/1-m)","axial load(kn)","mx(knm)","my(knm)","vdx(kn)","vdy(kn)"})
		row = fmt.Sprintf("bs 449, %v, %.3f, %.3f, %.3f, %.3f, %.3f",c.Dtyp, c.Pu/1e3, c.Mx/1e6, c.My/1e6, c.Vx/1e3, c.Vy/1e3)
		table.Append(strings.Split(row,","))
		table.Render()
	}
	rezstr.WriteString(ColorPurple)
	if c.Dz{
		scls := map[int]string{1:"plastic",2:"compact",3:"semi-compact"}
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true,"section data")
		switch c.Code{
			case 0,1:
			switch c.Dtyp{
				case 0:
				table.SetHeader([]string{"sec","wt(n/m)","dims","fcdx(n/mm2)","fcdy(n/mm2)","fcd(n/mm2)","pur(kn)","class"})
				for _, ss := range c.Ssecs{
					row = fmt.Sprintf("%s, %.f, %.f, %.1f, %.1f, %.1f, %.3f, %s",ss.Sstr,ss.Wt,ss.Dims,ss.Fcdx,ss.Fcdy,ss.Fcd,ss.Area*ss.Fcd/1e3,scls[ss.Cfac])
					table.Append(strings.Split(row,","))
				}
				case 1:
				//
				table.SetHeader([]string{"sec","wt(n/m)","dims","mx top(knm)","mx bot(knm)","my top(knm)","my bot(knm)","fcdx(n/mm2)","fcdy(n/mm2)","fcd(n/mm2)","fbd(n/mm2)","lfac","ifaca","ifacb","class"})
				
				for _, ss := range c.Ssecs{
					row = fmt.Sprintf("%s, %.f, %.f, %.f, %.f, %.f, %.f, %.1f, %.1f, %.1f, %.3f, %.3f, %.3f, %.3f, %s",ss.Sstr,ss.Wt,ss.Dims,ss.Max/1e6, ss.Mbx/1e6, ss.May/1e6, ss.Mby/1e6, ss.Fcdx,ss.Fcdy,ss.Fcd,ss.Fbd,ss.Lfac,ss.Ifaca,ss.Ifacb,scls[ss.Cfac])
					table.Append(strings.Split(row,","))
				}
			}
			case 2:
			table.SetHeader([]string{"sec","wt","dims","fcd(n/mm2)","fcc(n/mm2)","fcdx(n/mm2)","fccx(n/mm2)","fcdy(n/mm2)","fccy(n/mm2)","frat"})
		
			for _, ss := range c.Ssecs{
				row = fmt.Sprintf("%s, %.f, %.f, %.1f, %.1f, %.1f, %.1f, %.1f, %.1f, %.3f",ss.Sstr,ss.Wt,ss.Dims,ss.Fcd,ss.Fcc,ss.Fcdx,ss.Fccx,ss.Fcdy,ss.Fccy,ss.Frat)
				table.Append(strings.Split(row,","))
			}
		}
		table.Render()
	}
	rezstr.WriteString(ColorReset)
	c.Report = rezstr.String()
	if printz{
		fmt.Println(c.Report)
	}
}


func (c *Col) Init() (err error){
	if c.Sname == ""{
		if val, ok := kass.StlSnames[c.Styp]; !ok{
			err = fmt.Errorf("no base section type/name specified")
			return
		} else {
			c.Sname = val
		}
	}
	if c.Code == 0{
		c.Code = 1
	}
	if c.Nsecs == 0 || c.Nsecs > 10{c.Nsecs = 3}
	if c.Frmstr == ""{c.Frmstr = "2dt"}
	if c.Pfac == 0.0{c.Pfac = 1.0}
	if c.Lx == 0{c.Lx = c.Lspan}
	if c.Ly == 0{c.Ly = c.Lspan}
	if c.Tx == 0{c.Tx = 1.0}
	if c.Ty == 0{c.Ty = 1.0}
	return

}

func (c *Col) GetSec() (ss kass.StlSec, err error){
	var bf, dt, tf, tw float64
	switch c.Sname{
		case "built-i":
		if len(c.Params) < 4{err = fmt.Errorf("invalid params for %s section - %f",c.Sname,c.Params);return}
		bf = c.Params[0]
		dt = c.Params[1]
		tf = c.Params[2]
		tw = c.Params[3]
		c.Sdx = -1
		case "l2-ss","l2-os","ln2-ss","ln2-os":
		if len(c.Params) < 1{err = fmt.Errorf("invalid params for %s section - %f",c.Sname,c.Params);return}
		bf = c.Params[0]
		case "plate-i":
		if len(c.Params) < 2{err = fmt.Errorf("invalid params for %s section - %f",c.Sname,c.Params);return}
		bf = c.Params[0]
		dt = c.Params[1]	
	}
	ss, err = kass.GetStlSec(c.Sname, c.Sdx, c.Code, bf, dt, tf, tw)
	if err != nil{
		return
	}
	if c.Tx == 0.0{c.Tx = 1.0}
	if c.Ty == 0.0{c.Ty = 1.0}
	ss.Klx = c.Tx
	ss.Kly = c.Ty
	if c.Lx == 0.0{c.Lx = c.Lspan * c.Tx}
	if c.Ly == 0.0{c.Ly = c.Lspan * c.Ty}
	ss.Lx = c.Lx
	ss.Ly = c.Ly
	//if c.Lx > 0.0{ss.Leffx = c.Tx * c.Lx}
	//if c.Ly > 0.0{ss.Leffy = c.Ty * c.Ly}
	ss.Lspan = c.Lspan
	ss.Tx = c.Tx
	ss.Ty = c.Ty
	ss.Pu = c.Pu
	ss.Vbdx = c.Vx
	ss.Vbdy = c.Vy
	ss.Vax = c.Vax
	ss.Vbx = c.Vbx
	ss.Vay = c.Vay
	ss.Vby = c.Vby
	ss.Pfac = c.Pfac
	ss.H1 = c.H1
	ss.H2 = c.H2
	ss.Mux = c.Mx
	ss.Muy = c.My
	ss.Max = c.Max
	ss.Mbx = c.Mbx
	ss.May = c.May
	ss.Mby = c.Mby
	ss.Bg = c.Bg
	ss.Wg = c.Wg
	ss.Weld = c.Weld
	ss.Cleg = c.Cleg
	ss.Calctyp = c.Calctyp
	return
}

//ChkSec checks a column section for axial load/bm col
func (c *Col) ChkSec()(err error){
	var ss kass.StlSec
	ss, err = c.GetSec()
	if err != nil{
		return
	}
	if c.Tie{
		//is tensile member
		err = ss.TieChk800()
		if err == nil{
			c.Ssecs = append(c.Ssecs, ss)
		}
		return
	}
	if c.Onism{
		if !strings.HasPrefix(ss.Sstr, "ISMB"){
			err = fmt.Errorf("not an ismb section")
			return
		}
	}
	switch c.Code{
		case 1:
		//is800
		switch c.Dtyp{
			case 0:
			//axially loaded (pure) column
			err = ss.ColChk800()
			case 1:
			//beam-column
			err = ss.BmColChk800()
		}
		case 2:
		//bs449
		err = ss.ColChk449()
	}
	if err == nil{
		c.Ssecs = append(c.Ssecs, ss)
	}
	return
}

//ColDzBs designs a steel column section (using NEW AND IMPROVED kass.StlSec way) as in mosley/spencer section 6.1
func ColDzBs(c *Col) (err error){
	err = c.Init()
	if err != nil{
		return
	}
	if !c.Dsgn{
		err = c.ChkSec()
		return
	}
	ndx := kass.StlSdxLims[c.Sname]
	if ndx == 0{
		err = fmt.Errorf("%s design functions not written",c.Sname)
		return
	}
	if c.Sdx > 0{ndx = c.Sdx}
	for idx := ndx; idx >= 0; idx--{
		if len(c.Sdxs) == c.Nsecs{
			break
		}
		c.Sdx = idx
		err = c.ChkSec()
		if err == nil{
			c.Sdxs = append(c.Sdxs, idx)
		}
	}
	if c.Dsgn{
		if len(c.Ssecs) == 0{
			err = fmt.Errorf("no suitable section found")
		} else {
			err = nil
		}
	}
	return
}

//ColDBs designs a steel column section as in mosley/spencer section 6.1
func ColDBs(c *Col) (err error){
	//iterate from end of df
	c.Mindx = -1
	if c.Nsecs == 0{c.Nsecs = 5}
	c.Rez = []int{}
	//df := StlSecBs(c.Sectyp)
	var df dataframe.DataFrame
	df, err = kass.GetStlDf(c.Sname)
	
	var pa, px, py float64
	lx := c.Lx * c.Tx; ly := c.Ly * c.Ty
	pvec := pqcol[c.Grd]
	var vx, vy, mx, my float64
	if c.H2 > 0.0{
		vx = c.Vx * c.H2/(c.H1+c.H2)
		vy = c.Vy * c.H2/(c.H1+c.H2)
	}
	mx = c.Mx; my = c.My; c.Dtyp = 1
	for i := df.Nrow()-1; i > 0; i--{
		//log.Println("checking section->",df.Elem(i,1))
		if len(c.Rez) == c.Nsecs{
			break
		}
		if mx + my == 0.0 {
			c.Dtyp = 0 //member with framing beams
			mx = vx * (100.0 + df.Elem(i,3).Float()/2.0)/1000.0
			my = vy * (100.0 + df.Elem(i,5).Float()/2.0)/1000.0
		}
		fa := c.Pu*10.0/df.Elem(i,23).Float()
		fx := mx*1e3/df.Elem(i,15).Float()
		fy := my*1e3/df.Elem(i,16).Float()
		fp := fa/pvec[0] + fx/pvec[1] + fy/pvec[2]
		
		if fp > c.Pfac{
			continue
		}
		var s1 float64
		sx := lx * 100.0/df.Elem(i,13).Float()
		sy := ly * 100.0/df.Elem(i,14).Float()
		s1 = sx; if sy > s1 {s1 = sy}
		if s1 > 180.0 {continue}
		//permissible axial stress pa
		var y0, q4, q5 float64
		c0 := math.Pow(math.Pi,2) * EStl/math.Pow(s1,2)
		n0 := 0.3 * math.Pow(s1/100.0,2)
		switch{
			case c.Grd == 43:
			y0 = 250.0; q4 = 155.0; q5 = 143.0
			case c.Grd == 50:
			y0 = 350.0; q4 = 215.0; q5 = 200.0
			case c.Grd == 55:
			y0 = 430.0; q4 = 265.0; q5 = 245.0
		}
		if c.Grd == 50 && df.Elem(i,6).Float() >= 40.0{
			//CHEEECK THIS
			y0 = 325.0; q4 = 200.0; q5 = 185.0
		}
		a0 := (y0 + c0 * (n0 + 1.0))/2.0
		pa = (a0 - math.Sqrt((math.Pow(a0,2) - y0 * c0)))/1.7
		if s1 <=30 {
			pa = q4 - (q4 - q5) * s1/30.0
		}
		
		//permissible stress bending(x) px
		var dtrat float64
		if dtrat = df.Elem(i,3).Float()/df.Elem(i,6).Float(); dtrat < 5.0 {dtrat = 5.0}	
		//log.Println("checking px->",s1,dtrat)

		if c.Yeolde{
			px = PbcYeolde(s1, dtrat)
		} else {
			px = PbcLerp(c.Sname, c.Grd, s1, dtrat)
		}
		
		//permissible stress in bending (y) py
		py = pvec[2]
		fp = fa/pa + fx/px + fy/py

		//log.Println("***")
		if fp <= c.Pfac{
			wt := df.Elem(i,3).Float()
			c.Rez = append(c.Rez, i)
			c.Vals = append(c.Vals, []float64{fa, pa, fx,px, fy,py, fp, mx, my, sx, sy, s1, dtrat,wt})
			if c.Mindx == -1 || c.Vals[c.Mindx][13] > wt{
				c.Mindx = len(c.Rez)-1
			}
			if c.Spam{
				log.Println("section found->",df.Elem(i,1))
				log.Println("base fp->",fp)
				log.Println("srats->",sx,sy,s1)
				log.Println("paxial->",pa)
				log.Println("px->",px)
				log.Println("fp->",fp)
				log.Println("section->",df.Elem(i,1))
				log.Println("depth, web thickness->",df.Elem(i,3), df.Elem(i,6))
				log.Println("area, zx, zy->",df.Elem(i,23),df.Elem(i,15), df.Elem(i,16))
				log.Println("rx, ry->",df.Elem(i,13), df.Elem(i,14))
				log.Println("mx, my, s1, dtrat->", mx, my, s1, dtrat)
				log.Println("fa, pa, px, py, fp ->",fa, pa, px, py, fp)
				log.Println("***")
			}
		}
	}
	if len(c.Rez) == 0{err = errors.New("no suitable section found")}
	c.Dz = true
	return 
}

//ColCBs checks a column section as in mosley/spencer sec. 6.1
func ColCBs(c *Col) (float64, bool){
	//iterate from end of df
	//df, err := kass.GetStlDf(c.Styp)
	c.Mindx = -1
	df := StlSecBs(c.Sname)
	var pa, px, py, vx, vy, mx, my float64
	lx := c.Lx * c.Tx; ly := c.Ly * c.Ty
	pvec := pqcol[c.Grd]
	vx = c.Vx; vy = c.Vy; mx = c.Mx; my = c.My
	if c.H2 > 0.0{
		vx = c.Vx * c.H2/(c.H1+c.H2)
		vy = c.Vy * c.H2/(c.H1+c.H2)
	}
	if mx + my == 0.0 {
		mx = vx * (100.0 + df.Elem(c.Sdx,3).Float()/2.0)/1000.0
		my = vy * (100.0 + df.Elem(c.Sdx,5).Float()/2.0)/1000.0
	}
	fa := c.Pu*10.0/df.Elem(c.Sdx,23).Float()
	fx := mx*1e3/df.Elem(c.Sdx,15).Float()
	fy := my*1e3/df.Elem(c.Sdx,16).Float()
	
	fp := fa/pvec[0] + fx/pvec[1] + fy/pvec[2]
	if fp > 1.0{return fp, false}
	var s1 float64
	sx := lx * 100.0/df.Elem(c.Sdx,13).Float()
	sy := ly * 100.0/df.Elem(c.Sdx,14).Float()
	s1 = sx; if sy > s1 {s1 = sy}
	//if s1 > 180.0 {continue}
	
	//permissible axial stress pa
	var y0, q4, q5 float64
	c0 := math.Pow(math.Pi,2) * EStl/math.Pow(s1,2)
	n0 := 0.3 * math.Pow(s1/100.0,2)
	switch{
		case c.Grd == 43:
		y0 = 250.0; q4 = 155.0; q5 = 143.0
		case c.Grd == 50:
		y0 = 350.0; q4 = 215.0; q5 = 200.0
		case c.Grd == 55:
		y0 = 430.0; q4 = 265.0; q5 = 245.0
	}
	if c.Grd == 50 && df.Elem(c.Sdx,6).Float() >= 40.0{
		//CHEEECK THIS
		y0 = 325.0; q4 = 200.0; q5 = 185.0
	}
	a0 := (y0 + c0 * (n0 + 1.0))/2.0
	pa = (a0 - math.Sqrt((math.Pow(a0,2) - y0 * c0)))/1.7
	if s1 <=30 {
		pa = q4 - (q4 - q5) * s1/30.0
	}
	//permissible stress bending(x) px
	var dtrat float64
	if dtrat = df.Elem(c.Sdx,3).Float()/df.Elem(c.Sdx,6).Float(); dtrat < 5.0 {dtrat = 5.0}	
	//log.Println("checking px->",s1,dtrat)
	px = PbcLerp(c.Sname, c.Grd, s1, dtrat)
	//permissible stress in bending (y) py
	py = pvec[2]
	fp = fa/pa + fx/px + fy/py
	
	//log.Println("***")
	if c.Spam{
		log.Println("section->",df.Elem(c.Sdx,1))
		log.Println("depth, web thickness->",df.Elem(c.Sdx,3), df.Elem(c.Sdx,6))
		log.Println("area, zx, zy->",df.Elem(c.Sdx,23),df.Elem(c.Sdx,15), df.Elem(c.Sdx,16))
		log.Println("rx, ry->",df.Elem(c.Sdx,13), df.Elem(c.Sdx,14))
		log.Println("mx, my, s1, dtrat->", mx, my, s1, dtrat)
		log.Println("fa, pa, px, py, fp ->",fa, pa, px, py, fp)
		log.Println("***")
	}
	c.Rez = append(c.Rez, c.Sdx)
	c.Vals = append(c.Vals, []float64{fa, pa, fx,px, fy,py, fp, mx, my, sx, sy, s1, dtrat})
	c.Dz = true
	c.Table(true)
	fmt.Println(ColorYellow,"section",ColorReset)
	if fp > c.Pfac{
		fmt.Println(ColorRed,"over stressed",ColorReset)
	} else {
		fmt.Println(ColorGreen,"o.k",ColorReset)
	}
	return fp, fp < c.Pfac	
}


//ColDzIs designs a column as per is code(duggal, chap.8)
func ColDzIs(c *Col)(err error){
	err = c.Init()
	if err != nil{
		return
	}
	if !c.Dsgn{
		err = c.ChkSec()
		return
	}
	ndx := kass.StlSdxLims[c.Sname]
	if ndx == 0{
		err = fmt.Errorf("%s design functions not written",c.Sname)
		return
	}
	if c.Sdx > 0{ndx = c.Sdx}
	for idx := ndx; idx >= 0; idx--{
		if len(c.Sdxs) == c.Nsecs{
			break
		}
		//fmt.Println("checking ndx-", idx)
		c.Sdx = idx
		err = c.ChkSec()
		if err == nil{
			c.Sdxs = append(c.Sdxs, idx)
		}
	}
	return
}

//ColDz is the main dz entry func
func ColDz(c *Col)(err error){
	switch c.Code{
		case 0,1:
		err = ColDzIs(c)
		case 2:
		err = ColDzBs(c)
	}
	if err == nil{
		c.Dz = true
		c.Table(c.Verbose)
	}
	return
}

//ColTCIs checks a tension member as in duggal, chap.7
func ColTCIs(c *Col)(err error){
	fmt.Println("sname-",c.Sname)
	return
}

//PbcBs returns permissible bending stresses as per bs449
func PbcBs(sname string, grade int) (vec [][]float64){
	var mpbc map[string][][]float64
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	jsonin := filepath.Join(basepath,"../data/steel/bsteel","pbc.json")
	jsonfile, err := ioutil.ReadFile(jsonin)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal([]byte(jsonfile), &mpbc)
	if err != nil {
		log.Println(err)
	}
	var query string
	switch sname{
		case "ub", "uc":
		//uc beams and columns
		switch grade{
			case 43:
			query = "3a" 
			case 50:
			query = "3b"
			case 55:
			query = "3c"
		}
	}
	return mpbc[query]
}

//PbcIs returns the permissible compressive stress as per is800 (merchant-rankine formula)
func PbcIs(){
	//(insert) merchant rankine formula for permissible compressive stress here
	//victory
}

//StlSecBs returns the section type dataframe from csv data sheets
func StlSecBs(sname string) (dataframe.DataFrame){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	var sheet string
	switch sname{
		case "ub":
		//ub sec
		sheet = filepath.Join(basepath,"../data/steel/bsteel","UB.csv")
		case "uc":
		//uc sec
		sheet = filepath.Join(basepath,"../data/steel/bsteel","UC.csv")
	}
	//log.Println("sheet->",sheet)
	csvfile, err := os.Open(sheet)
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(csvfile)
	return df
}

//PbcYeolde returns the permissible bending stress as in table 6.1 of mosley/spencer (ye olde values)
func PbcYeolde(s1, dtrat float64) (pbc float64){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	sheet := filepath.Join(basepath,"../data/steel","hulsepbc43.csv")
	csvfile, err := os.Open(sheet)
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(csvfile)
	var rdx, cdx int
	switch{
		case s1 <= 90:
		pbc = df.Elem(0,1).Float()
		return
		case s1 <= 120:
		rdx = int((s1 - 90.)/5.0)
		default:
		rdx = int((s1 - 120.)/10.0) + (120 - 90)/5
	}
	if dtrat <= 10{dtrat = 10}
	if dtrat <= 40{
		cdx = int((dtrat - 10.)/5.0) + 1
	} else {
		cdx = 7
	}
	//var sa, sb float64
	sa := df.Elem(rdx,0).Float(); sb := df.Elem(rdx+1,0).Float()
	//log.Println("sa, sb->",sa, sb, rdx)
	var p1, p2 float64
	pt0 := df.Elem(rdx,cdx).Float(); pt1 := df.Elem(rdx,cdx+1).Float()
	if cdx < 7 {
		p1 = pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	} else {
		p1 = pt0 + (dtrat-40.)*(pt1 - pt0)/10.0
	}
	//log.Println(pt0, pt1)
	pt0 = df.Elem(rdx+1,cdx).Float(); pt1 = df.Elem(rdx+1,cdx+1).Float()
	if cdx < 7 {
		p2 = pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	} else {
		p2 = pt0 + (dtrat-40.)*(pt1 - pt0)/10.0
	}
	
	//log.Println(pt0, pt1)
	//log.Println(p1, p2)
	pbc = p1 + (s1 - sa) * (p2 - p1)/(sb - sa)
	return
}

//PbcLerp linearly interpolates the permissible bending stress given a slenderness ratio
//calls PbcBs for the table of permissible bending stresses
func PbcLerp(sname string, grd int, s1, dtrat float64) (pbc float64){
	//log.Println("lerp in-> sname, grd, srat, drat->",sname,grd,s1, dtrat)
	pbvec := PbcBs(sname, grd)
	pvec := pqcol[grd]
	var rdx, cdx int
	switch {
	case s1 <= 40.0:
		pbc = pvec[1]
		return
	case s1 <= 120.0:
		cdx = int((s1 - 40.0)/5.0)
	case s1 <= 300.0:
		cdx = int((s1-120.0)/10.0)+(120-40)/5
	}
	rdx = int((dtrat-5)/5.0)+1
	//log.Println("dxs->",rdx, cdx)
	sa := pbvec[0][cdx]; sb := pbvec[0][cdx+1]
	// log.Println("dtrats->",rdx*5, (rdx+1)*5)
	// log.Println("srats->",sa,sb)
	// log.Println("rdx, cdx->",rdx,cdx)
	pt0 := pbvec[rdx][cdx]; pt1 := pbvec[rdx+1][cdx]
	// log.Println("pts 1",pt0,pt1)
	p1 := pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	if cdx == 33{
		//log.Println(len(pbvec), len(pbvec[0]))
		pbc = p1
		return
	}
	pt0 = pbvec[rdx][cdx+1]; pt1 = pbvec[rdx+1][cdx+1]
	// log.Println("pts 2",pt0,pt1)
	p2 := pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	pbc = p1 + (s1 - sa) * (p2 - p1)/(sb - sa)
	return
}

//ColDesign is the entry func for steel column design (IS IT)
//it isn't
func ColDesign(c *Col)(err error){
	//log.Println(ColorRed,"***insert col design idito**",ColorReset)
	if c.Dtyp != 0{
		//get uvals
		log.Println(ColorRed,"***insert col design idito**",ColorReset)
	}
	if c.Tz{
		switch c.Code{
			case 1:
			case 2:
		}
	}
	switch c.Code{
		case 1:
		// err = ColDIs(c)
		// if err == nil && c.Verbose{
		// 	c.Table(true)
		// }
		case 2:
		err = ColDBs(c)
		if err == nil && c.Verbose{
			c.Table(true)
		}
		//c.PlotSecs()
		return
	}
	return
}

/*
   
	//log.Println("dxs->",rdx, cdx)
	sa := pbvec[0][cdx]; sb := pbvec[0][cdx+1]
	//log.Println("dtrats->",rdx*5, (rdx+1)*5)
	//log.Println("srats->",sa,sb)
	//log.Println("rdx, cdx->",rdx,cdx)
	if cdx == 0{
		pt0 := pbvec[rdx][cdx]
		pt1 := pbvec[rdx+1][cdx]
		pbc = pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
		return
	}
	pt0 := pbvec[rdx][cdx-1]; pt1 := pbvec[rdx+1][cdx-1]
	//log.Println("pts 1",pt0,pt1)
	p1 := pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	pt0 = pbvec[rdx][cdx]; pt1 = pbvec[rdx+1][cdx]
	//log.Println("pts 2",pt0,pt1)
	p2 := pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	pbc = p1 + (s1 - sa) * (p2 - p1)/(sb - sa)
	return
*/
// 
