package barf

import (
	"log"
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

//Dimchks checks column dims - make this a recurring thingie
func (c *WdCol) Dimchks()(err error){
	return
}

//Init initalizes a WdCol struct
func (c *WdCol) Init() (err error){
	if c.Grp != 0 {err = c.Prp.Init(c.Grp)}
	if err != nil{return}
	//fmt.Println(c.Prp)
	//default pinned-pinned
	if c.Endc == 0{c.Endc = 3}
	//eff len factors c.Ke from abel, pg 187
	if c.Ke == 0{		
		switch c.Endc{
			case 1:
			//f-f, no sway
			c.Ke = 0.5
			// c.Kerst = 4.0
			case 2:
			//p-f, no sway
			c.Ke = 0.7
			case 3:
			//p-p, no sway
			c.Ke = 1.0
			case 4:
			//f-f, sway
			c.Ke = 1.0
			case 5:
			//f-free, sway
			c.Ke = 2.0
			case 6:
			//f-p, sway
			c.Ke = 2.0
		}	
	}
	//what is this lol (column fixity factor)
	if c.Kfac == 0.0{c.Kfac = 1.0}
	c.Le = c.Lspan * c.Ke
	switch c.Styp{
		case 0:
		// c.Cbi = 0.85
		case 1:
		// c.Cbi = 0.8
		case 4:
		if c.Tplnk == 0.0{c.Tplnk = 25.0}
	}
	// if len(c.Dims) != 0{
	// 	s := kass.SecGen(c.Styp, c.Dims)
	// 	c.Sec = s
	// }
	if c.Clctyp == 0{c.Clctyp = 1}
	return
}

//ColChk checks a solid/spaced column section for an axial load (pu)
func ColChk(c *WdCol) (ok bool, vals []float64, err error){
	//styp, conv, wcode, wtcalc int, dim []float64, pg, em, le, fcp, kfac, cbi, pul float64
	//solid column check
	var selr, sp, sdrat, k8, k9, k10, dmin, d1, d2 float64
	bstyp := c.Styp
	dims := make([]float64, len(c.Dims))
	copy(dims, c.Dims)
	c.Init()
	switch c.Styp{
		case 0:
		//get eq. sqr
		bstyp = 1
		side := math.Sqrt(math.Pi) * c.Dims[0]/2.0
		dims = []float64{side, side}
		case 4:
		B := dims[0]; D := dims[1]; b := dims[2]; d := dims[3]
		d1 = B; d2 = b
		if D < d1{d1 = D}
		if d < d2{d2 = d}
		// if B > 5.0 * c.Tplnk{
		// 	err = fmt.Errorf("invalid width (%.f) > 5 * plank thicknes (%.f)",B, c.Tplnk)
		// 	return 
		// }
		if c.Solid{
			dims = []float64{B, D}
			bstyp = 1
		}
		case 26:
		bstyp = 1
		default:
	}
	c.Sec = kass.SecGen(bstyp, dims)
	//s := kass.SecGen(b.Styp, dim)
	if c.Sec.Prop.Area <= 0.0{
		err = fmt.Errorf("invalid section properties (area - %.f)\n styp %v dims %v",c.Sec.Prop.Area, c.Styp, c.Dims)
		return 
	}
	wdl := c.Sec.Prop.Area * c.Prp.Pg * 9.8 * 1e-6 
	if c.Styp == 26{wdl = wdl * 2.0}
	pul := c.Pu
	if c.Selfwt {pul += wdl}
	if c.Styp == 26{
		//spaced col.
		pul = pul/2.0
		//set end restraint factore
		if c.Kerst == 0.0{
			c.Kerst = 2.5
		}
	}

	sp = c.Prp.Fc //comp. PARALLEL TO GRAIN HERE
	
	rmin := c.Sec.Prop.Rxx
	if rmin > c.Sec.Prop.Ryy {rmin = c.Sec.Prop.Ryy}
	dmin = c.Dims[0]
	switch c.Styp{
		case 0:
		dmin = dims[0]
		case 1,26:
		if c.Dims[1] < dmin{
			dmin = c.Dims[1]
		}
	}
	if c.Styp == 1 && c.Dims[1] < dmin{dmin = c.Dims[1]}
	//box columns as per is883; nothing in abel
	if c.Styp == 4{c.Code = 1}
	//calc euler buckling load
	switch c.Styp{
		case 0, 1, 26:
		//circle, rect, spaced col
		selr = 0.3 * c.Prp.Em/math.Pow(c.Le/dmin,2)
	sdrat = c.Le/dmin
		//note - check sp.33 sec 5.3.3.1 for c.Kerst
		case 4:
		//box column
		sdrat = c.Le/math.Sqrt(math.Pow(d1, 2) + math.Pow(d2,2))
	}
	switch c.Code{
		case 1:
		//sp 33/is 883
		if c.Spam{log.Printf("sdrat - %.2f euler stress - %.2f\n",sdrat, selr)}
		switch c.Styp{
			case 0, 1:
			//both sqr and rect seem to be the same here?
			if sdrat > 50.0{
				err = fmt.Errorf("too slender - %.2f",sdrat)
				return 
			}
			k8 = 0.702 * math.Sqrt(c.Prp.Em/sp)
			switch{
				case sdrat < 11.0:
			        //short col
				case sdrat < k8:
				//int col
				spfac := 1.0 - math.Pow(sdrat/k8,4.0)/3.0
				sp = sp * spfac
				case sdrat > k8:
				//long col
				sp = 0.329 * c.Prp.Em/math.Pow(sdrat,2.0) 
				
			}
			if c.Spam{
				log.Printf("k8 -> %.2f\n",k8)
				log.Printf("sp -> %.2f\n",sp)
			}
			case 4:
			//box column
			if sdrat > 50.0{
				err = fmt.Errorf("too slender - %.2f",sdrat)
				return 
			}
			var ufac, qfac float64
			//is.883 sec 7.6.2.5
			switch c.Tplnk{
				case 25.0:
				ufac = 0.8
				qfac = 1.0
				case 50.0:
				ufac = 0.6
				qfac = 1.0
			}
			k9 = math.Pi * math.Sqrt(ufac * c.Prp.Em/(5.0 * qfac * c.Prp.Fc))/2.0
			switch{
				case sdrat < 8.0:
				//short
				sp = qfac * sp
				c.Bclass = 1
				case sdrat < k9:
				sp = qfac * sp * (1.0 - math.Pow(sdrat/k9, 4)/3.0)
				c.Bclass = 2
				case sdrat > k9:
				//long
				sp = 0.329 * ufac * c.Prp.Em/math.Pow(sdrat, 2.0)
				c.Bclass = 3
			}
			if c.Spam{
				log.Printf("k9 -> %.2f\n",k10)
				log.Printf("sp -> %.2f\n",sp)
			}
			case 26:
			if sdrat > 80.0{
				
				err = fmt.Errorf("too slender - %.2f",sdrat)
				return 
			}
			k10 = 0.702 * math.Sqrt(c.Kerst * c.Prp.Em/c.Prp.Fc)
			switch{
				case sdrat < 11.0:
				case sdrat < k10:
				//int col
				spfac := 1.0 - math.Pow(sdrat/k10,4.0)/3.0
				sp = sp * spfac
				case sdrat > k10:
				//long col
				sp = 0.329 * c.Prp.Em * c.Kerst/math.Pow(sdrat,2.0) 
			}
			if c.Spam{
				log.Printf("k10 -> %.2f\n",k10)
				log.Printf("sp -> %.2f\n",sp)
			}
		}
		case 2:
		//madison approach
		switch c.Styp{
			case 0:
			selr = 3.619 * c.Prp.Em/math.Pow(c.Le/rmin,2)
			sdrat = c.Le/rmin
			k8 = 2.32 * math.Sqrt(c.Prp.Em/c.Prp.Fc)
			if sdrat > 173.0{
				err = fmt.Errorf("too slender - %.2f",sdrat)
				return
			}
			case 1, 26:
			if sdrat > 50.0{
				err = fmt.Errorf("too slender - %.2f",sdrat)
				return
			}
		}
		switch c.Styp{
			case 0:
			switch{
				case sdrat < 38.0:
				case sdrat < k8:
				spfac := 1.0 - math.Pow(sdrat/rmin, 4.0)/3.0
				sp = sp * spfac
				case sdrat > k8:
				sp = 3.619 * c.Prp.Em/math.Pow(sdrat,2.0)
				
			}
			case 1:
			k8 = 0.671 * math.Sqrt(c.Prp.Em/c.Prp.Fc)
			switch{
				case sdrat < 11.0:
				case sdrat < k8:
				spfac := 1.0 - math.Pow(sdrat/dmin, 4.0)/3.0
				sp = sp * spfac
				case sdrat > k8:
				sp = 0.3 * c.Prp.Em/math.Pow(sdrat, 2.0)
			}
			case 26:
			k8 = 0.671 * math.Sqrt(c.Prp.Em * c.Kerst/c.Prp.Fc)	
			switch{
				case sdrat < 11.0:
				case sdrat < k8:
				spfac := 1.0 - math.Pow(sdrat/dmin, 4.0)/3.0
				sp = sp * spfac
				case sdrat > k8:
				sp = 0.3 * c.Prp.Em * c.Kerst/math.Pow(sdrat, 2.0)
			}
		}
		if c.Spam{
			log.Printf("k8 -> %.2f\n",k8)
			log.Printf("sp -> %.2f\n",sp)
		}
		case 3:
		// //nds 1991, revised madison approach - TO CHECK
		// //fmt.Println("revised madison approach")
		z := (1.0 + selr/sp)/2.0/c.Cbi
		cp := z - math.Sqrt(math.Pow(z,2)/c.Cbi - selr/sp/c.Cbi)
		sp = cp * sp
	}
	switch c.Dtyp{
		// case 0:
		//axial load
		case 1:
		//axial load + flexure
		//return here
	}
	if c.Spam{
		log.Printf("sp -> %.2f vs actual -> %.2f\n",sp, pul/c.Sec.Prop.Area)
	}
	sact := pul/c.Sec.Prop.Area
	if sact <= sp * c.Kfac{
		ok = true
		vals = []float64{sp,selr,pul/c.Sec.Prop.Area,sdrat}
		return
	}
	//if c.Tensile{
	//	return c.Prp.Ft <= c.Pu/c.Sec.Prop.Area, []float64{sp,selr,pul/c.Sec.Prop.Area,sdrat}
	//}
	err = fmt.Errorf("permissible stress exceeded - %.3f N/mm2 vs %.3f N/mm2 actual", sp, sact)
	vals = []float64{sp,selr,pul/c.Sec.Prop.Area,sdrat}
	return
}

//ColDz designs a solid column section given design values (pu, wood properties)
//chapter 9, abel o. olorunnisola
//buckling interaction factor cbi for glulam = 0.9
func ColDz(c *WdCol) (err error){
	err = c.Init()
	if err != nil{return}
	if c.Nsecs == 0{c.Nsecs = 4}
	switch c.Styp{
		case 0, 1:
		var basedims [][]float64
		if c.Styp == 0{
			basedims = kass.TmbrDims0
			c.Cbi = 0.8
		} else {
			kass.GenTmbrDims1()
			basedims = kass.TmbrDims1
		}	
		for _, dim := range basedims{
			if len(dim) == 0{continue}
			if len(c.Rez) == c.Nsecs{break}
			switch c.Styp{
				case 1:
				if dim[0] > dim[1]{
					continue
				}
				case 4:
				
			}
			c.Dims = dim
			ok, val, _ := ColChk(c)
			if ok{
				c.Rez = append(c.Rez, dim)
				c.Vals = append(c.Vals, val)
			} 
			c.Dims = []float64{}
		}
		case 4:
		//box/built up
		if c.Tplnk == 0.0{
			c.Tplnk = 25.0
		}
		// if c.Dplnk == 0.0{
		// 	c.Dplnk = 300.0
		// }
		case 26:
		//spaced
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
