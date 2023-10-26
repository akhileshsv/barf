package barf

import (
	"fmt"
	"math"
	"time"
	"strings"
	"math/rand"
	"github.com/olekukonko/tablewriter"
	kass"barf/kass"
)

//Table generates an ascii table report for a WdCol
func (c *WdCol) Table(printz bool){
	if c.Title == ""{
		if c.Id == 0{
			c.Id = rand.Intn(666)
		}
		c.Title = fmt.Sprintf("tmbr_col_%v",c.Id)
	}
	rezstr := new(strings.Builder)
	hdr := fmt.Sprintf("%s\ntimber column report\ndate-%s\n%s\n%s\n",ColorYellow,time.Now().Format("2006-01-02"),c.Title,ColorReset)
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)
	var row string
	table.SetCaption(true,"timber properties")
	table.SetHeader([]string{"group","section type","specific gravity","elastic modulus\n(n/mm2)"})
	row = fmt.Sprintf("%v, %v, %.2f, %.2f",c.Grp,c.Styp,c.Prp.Pg,c.Prp.Em)
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorBlue)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"allowable stresses (N/mm2)")
	table.SetHeader([]string{"bending","tension","comp.\npar.to grain","comp.\nperp.to grain","shear\npar.to grain","shear\nperp.to grain"})
	row = fmt.Sprintf("%.2f,%.2f,%.2f,%.2f,%.2f,%.2f",c.Prp.Fcb,c.Prp.Ft,c.Prp.Fc,c.Prp.Fcp,c.Prp.Fv,c.Prp.Fvp)
	table.Append(strings.Split(row,","))
	table.Render()
	if c.Dz{
		rezstr.WriteString(ColorPurple)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true,"section details")
		table.SetHeader([]string{"section","dims\n(mm)","perm.stress\n(n/mm2)","buckling stress\n(n/mm2)","actual stress\n(n/mm2)"})
		for i, dims := range c.Rez{
			row = fmt.Sprintf("%s,%.f,%.2f,%.2f,%.2f",kass.SectionMap[c.Styp],dims,c.Vals[i][0],c.Vals[i][1],c.Vals[i][2])
			table.Append(strings.Split(row,","))
		}
		table.Render()		
	}
	rezstr.WriteString(ColorReset)
	c.Report = fmt.Sprintf("%s",rezstr)
	if printz{
		fmt.Println(c.Report)
	}
	return
}

//Init initalizes a WdCol struct
func (c *WdCol) Init() (err error){
	if c.Grp != 0 {err = c.Prp.Init(c.Grp)}
	if err != nil{return}
	//fmt.Println(c.Prp)
	//default pinned
	if c.Ke == 0.0{c.Ke = 1.0}
	//what is this lol
	if c.Kfac == 0.0{c.Kfac = 1.0}
	c.Le = c.Lspan * c.Ke
	switch c.Styp{
		case 0:
		c.Cbi = 0.85
		case 1:
		c.Cbi = 0.8
	}
	if len(c.Dims) != 0{
		s := kass.SecGen(c.Styp, c.Dims)
		c.Sec = s
	}
	if c.Clctyp == 0{c.Clctyp = 1}
	return
}


//ColChk checks a solid column section for given design values (pu)
func ColChk(c *WdCol) (bool, []float64){
	//styp, conv, wcode, wtcalc int, dim []float64, pg, em, le, fcp, kfac, cbi, pul float64
	//solid column check
	var cp, selr, sp, z, sdrat, k8, dmin float64
	//s := kass.SecGen(b.Styp, dim)
	wdl := c.Sec.Prop.Area * c.Prp.Pg * 9.8 * 1e-6 
	//ei := s.Prop.Ixx * b.Prp.Em
	//bar, wdl := WdSecCp(styp, 1, dim, pg)
	//if wtcalc == 1 {pul += wdl}
	pul := c.Pu
	if c.Selfwt {pul += wdl}
	sp = c.Prp.Fc //PARALLEL TO GRAIN HERE
	dmin = c.Dims[0]
	rmin := c.Sec.Prop.Rxx
	if rmin > c.Sec.Prop.Ryy {rmin = c.Sec.Prop.Ryy}
	if c.Styp == 1 && c.Dims[1] < dmin{dmin = c.Dims[1]}
	//calc euler buckling load
	if c.Styp == 0{
		selr = 3.619 * c.Prp.Em/math.Pow(c.Le/rmin,2)
		sdrat = c.Le/rmin
	} else {
		selr = 0.3 * c.Prp.Em/math.Pow(c.Le/dmin,2)
		sdrat = c.Le/dmin
		//fmt.Printf("euler buckling stress %.3f N/mm2\n",selr)
	}
	//fmt.Println("s/d-",sdrat)
	switch c.Code{
		case 1:
		//madison approach, sp 33/is 883
		//fmt.Println("using madison approach")
		switch c.Styp{
			case 0, 1:
			sdrat = c.Le/dmin
			if sdrat > 50.0{
				return false, []float64{}
			}
			k8 = 0.702 * math.Sqrt(c.Prp.Em/sp)
			switch{
				case sdrat <= 11.0:
				//short col
				case sdrat <= k8:
				//intermediate col
				sp = sp * (1.0 - math.Pow(sdrat/k8,4.0)/3.0)
				case sdrat > k8:
				sp = 0.329 * c.Prp.Em/math.Pow(sdrat,2.0)
			}
		}
		case 2:
		//nds 1991, revised madison approach
		//fmt.Println("revised madison approach")
		z = (1.0 + selr/sp)/2.0/c.Cbi
		cp = z - math.Sqrt(math.Pow(z,2)/c.Cbi - selr/sp/c.Cbi)
		sp = cp * sp
	}
	if pul/c.Sec.Prop.Area <= sp * c.Kfac{
		return true, []float64{sp,selr,pul/c.Sec.Prop.Area,sdrat}
	}
	return false, []float64{sp,selr,pul/c.Sec.Prop.Area,sdrat}
}

//ColDz designs a solid column section given design values (pu, wood properties)
//chapter 9, abel o. olorunnisola
//buckling interaction factor cbi for glulam = 0.9
func ColDz(c *WdCol) (err error){
	var basedims [][]float64
	err = c.Init()
	if err != nil{return}
	if c.Nsecs == 0{c.Nsecs = 3}
	switch c.Styp{
		case 0:
		//circle
		basedims = kass.TmbrDims0
		c.Cbi = 0.85
		//sp.33 recos eq. sqr col
		//if c.Code == 1{
		//	sqDims := getEqSqdims()
		//	basedims = sqDims
		//	styp = 1
		//	c.Cbi = 0.8
		//}
		case 1:
		//rect/square
		if len(kass.TmbrDims1) == 0 {
			kass.GenTmbrDims1()
		}
		basedims = kass.TmbrDims1	
		c.Cbi = 0.8
		//I SEKSHUN
	}
	
	for _, dim := range basedims{
		s := kass.SecGen(c.Styp, dim)
		c.Dims = dim
		c.Sec = s
		if len(c.Rez) == c.Nsecs{break}
		ok, val := ColChk(c)
		if ok{
			c.Rez = append(c.Rez, dim)
			c.Vals = append(c.Vals, val)
		}
		c.Dims = []float64{}
	}
	if len(c.Rez) == 0 {
		err = ErrDim
		return
	}
	c.Dims = c.Rez[0]
	s := kass.SecGen(c.Styp, c.Dims)
	c.Sec = s
	err = nil
	c.Dz = true
	c.Table(c.Verbose)
	return
}
