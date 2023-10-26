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
	table2bs = []float64{180,165,230,215,280,170,155,215,200,265,185}	
	stlsecmap = map[int]string{7:"UB",8:"UC"}
	//pqcol - from mosley section 6
	pqcol = map[int][]float64{43:{155,165,165},50:{215,230,230},55:{265,280,280}}
	//pqbm - qx (bending), ps (shear), pc (web crushing)
	pqbm = map[int][]float64{43:{165,100,190},50:{230,140,260},55:{280,170,320}}
	pqcol89 = map[int][]float64{43:{170,180,180},50:{215,230,230},55:{265,280,280}}
	EStl = 210000.0
	EStl89 = 205000.0
)

//Col is a steel column struct
//mosley spencer section 6.1
type Col struct{
	Title      string
	Id         int
	H1, H2, Lx, Ly, Tx, Ty, Mx, My, Vx, Vy, Pu, Pfac float64
	Lspan      float64
	Grd, Styp  int
	Dtyp       int
	Nsecs, Sdx int
	Code, Endc int
	Rez        []int
	Vals       [][]float64
	Yeolde     bool
	Verbose    bool
	Blurb      bool
	Spam       bool
	Braced     bool
	Dz         bool
	Report     string
	Kostin     float64
	Mindx      int
}

//PbcBs returns permissible bending stresses as per bs449
func PbcBs(sectyp, grade int) (vec [][]float64){
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
	switch sectyp{
		case 7, 8:
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
func StlSecBs(sectyp int) (dataframe.DataFrame){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	var sheet string
	switch sectyp{
		case 7:
		//ub sec
		sheet = filepath.Join(basepath,"../data/steel/bsteel","UB.csv")
		case 8:
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
func PbcLerp(sectyp, grd int, s1, dtrat float64) (pbc float64){
	//log.Println("lerp in-> srat, drat->",s1, dtrat)
	pbvec := PbcBs(sectyp, grd)
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
	//log.Println("dtrats->",rdx*5, (rdx+1)*5)
	//log.Println("srats->",sa,sb)
	//log.Println("rdx, cdx->",rdx,cdx)
	pt0 := pbvec[rdx][cdx]; pt1 := pbvec[rdx+1][cdx]
	//log.Println("pts 1",pt0,pt1)
	p1 := pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	if cdx == 33{
		//log.Println(len(pbvec), len(pbvec[0]))
		pbc = p1
		return
	}
	pt0 = pbvec[rdx][cdx+1]; pt1 = pbvec[rdx+1][cdx+1]
	//log.Println("pts 2",pt0,pt1)
	p2 := pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	pbc = p1 + (s1 - sa) * (p2 - p1)/(sb - sa)
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
	df, err = kass.GetStlDf(c.Styp)
	
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
			px = PbcLerp(c.Styp, c.Grd, s1, dtrat)
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
	table.SetHeader([]string{"grade","section type","height(above)(m)","height(col)(m)","unb.len(lx)(m)","unb.len(ly)(m)","tx","ty"})
	row = fmt.Sprintf("%v, %s, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f",c.Grd,stlsecmap[c.Styp],c.H1,c.H2,c.Lx,c.Ly,c.Tx,c.Ty)
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorRed)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"ultimate loads")
	table.SetHeader([]string{"axial load(kn)","dtyp(0-b/1-m)","mx(kn-m)","my(kn-m)","vx(kn)","vy(kn)"})
	row = fmt.Sprintf("%.3f,%v,%.3f,%.2f,%.2f,%.2f",c.Pu,c.Dtyp, c.Mx, c.My, c.Vx, c.Vy)
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorPurple)
	if c.Dz{
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true,"section geometry")
		table.SetHeader([]string{"section","wt\n(kg/m)","depth\n(mm)","t.web\n(mm)","area\n(cm2)","rxx\n(cm)","ryy\n(cm)","zxx\n(cm3)","zyy\n(cm3)",})
		/*

		*/
		df, _ := kass.GetStlDf(c.Styp)
		for _, idx := range c.Rez{
			//fa, pa, px, py, fp, mx, my, sx, sy, s1, dtrat := c.Vals[i]
			sname, wt, dw, tw, ar, rxx, ryy, zxx, zyy := df.Elem(idx,1),df.Elem(idx,2).Float(),df.Elem(idx,3).Float(),df.Elem(idx,6).Float(),df.Elem(idx,23).Float(),df.Elem(idx,13).Float(),df.Elem(idx,14).Float(),df.Elem(idx,15).Float(),df.Elem(idx,16).Float() 
			row = fmt.Sprintf("%s,%.3f,%.f,%.f,%.f,%.f,%.f,%.f,%.f",sname, wt, dw, tw, ar, rxx, ryy, zxx, zyy)
			table.Append(strings.Split(row,","))
		}
		table.Render()
		rezstr.WriteString(ColorGreen)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true,"section results")
		table.SetHeader([]string{"section","mx\n(kn-m)","my\n(kn-m)","s1","sx","sy","dtrat","fp","fa\n(n/mm2)","pa\n(n/mm2)","fx\n(n/mm2)","px\n(n/mm2)","fy\n(n/mm2)","py\n(n/mm2)"})
		for i, idx := range c.Rez{
			sname, fa, pa, fx, px, fy, py, fp, mx, my, sx, sy, s1, dtrat := df.Elem(idx,1),c.Vals[i][0],c.Vals[i][1],c.Vals[i][2],c.Vals[i][3],c.Vals[i][4],c.Vals[i][5],c.Vals[i][6],c.Vals[i][7],c.Vals[i][8],c.Vals[i][9],c.Vals[i][10], c.Vals[i][11], c.Vals[i][12]
			row = fmt.Sprintf("%s,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f",sname,mx,my,s1,sx,sy,dtrat,fp,fa,pa,fx,px,fy,py)
			table.Append(strings.Split(row,","))
		}
		table.Render()
		

		if c.Mindx != -1{
			rezstr.WriteString(ColorCyan)
			if c.Kostin == 0.0{c.Kostin = 200.0}
			table = tablewriter.NewWriter(rezstr)
			wt := c.Vals[c.Mindx][13]
			minsec := df.Elem(c.Rez[c.Mindx],1)
			table.SetCaption(true, "quantity take off")
			table.SetHeader([]string{"section","min. wt\n(kg)","cost\n(rs)","span\n(m)","total cost\n(rs)"})
			row = fmt.Sprintf("%s, %.3f, %.3f, %.3f, %.3f",minsec,wt,c.Kostin,c.H1,c.Kostin * wt*c.H1)
			table.Append(strings.Split(row, ","))
			table.Render()
		}
	}
	rezstr.WriteString(ColorReset)
	c.Report = fmt.Sprintf("%s",rezstr)
	if printz{
		fmt.Println(c.Report)
	}

}

//ColCBs checks a column section as in mosley/spencer sec. 6.1
func ColCBs(c *Col) (float64, bool){
	//iterate from end of df
	//df, err := kass.GetStlDf(c.Styp)
	c.Mindx = -1
	df := StlSecBs(c.Styp)
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
	var s1 float64
	sx := lx * 100.0/df.Elem(c.Sdx,13).Float()
	sy := ly * 100.0/df.Elem(c.Sdx,14).Float()
	s1 = sx; if sy > s1 {s1 = sy}
	//if s1 > 180.0 {continue}
	//log.Println("srats->",sx,sy,s1)
	
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
	//log.Println("paxial->",pa)
	//permissible stress bending(x) px
	var dtrat float64
	if dtrat = df.Elem(c.Sdx,3).Float()/df.Elem(c.Sdx,6).Float(); dtrat < 5.0 {dtrat = 5.0}	
	//log.Println("checking px->",s1,dtrat)
	px = PbcLerp(c.Styp, c.Grd, s1, dtrat)
	//log.Println("px->",px)
	//permissible stress in bending (y) py
	py = pvec[2]
	fp = fa/pa + fx/px + fy/py
	//log.Println("fp->",fp)
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

//ColDesign is the entry func for steel column design
func ColDesign(c *Col)(err error){
	//log.Println(ColorRed,"***insert col design idito**",ColorReset)
	switch c.Code{
		case 1:
		case 2:
		err = ColDBs(c)
		if err == nil && c.Verbose{
			c.Table(true)
		}
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
