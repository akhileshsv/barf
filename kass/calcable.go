package barf

import (
	"fmt"
	"log"
	"math"
	"strings"
	"github.com/olekukonko/tablewriter"	
)

//Cbl - THIS IS NOT DONE (it is not werk) avoid for now
//see harrison sec WHAAT
type Cbl struct{
	//cable 
	Ucl float64 //unstretched cable length
	Em float64 //em
	Dia float64 //dia
	Ar float64 //area
	Alp float64 //coeff of expansion
	Xr float64 //x right (xl, yl = 0,0)
	Yr float64 //y right 
	Tr float64 //temp rise
	Sw float64 //self wt/unit length
	Lds [][]float64 //loads along length [fx, fy, dist]
	Nl int //no. of elements
	Lseg float64 //segment length
	Tl float64 //tension at left sup
	Theta float64 //slope at left sup
	Xd []float64
	Yd []float64
	El []float64
	Ts []float64
	Pts [][]float64
	Dz bool
	Verbose bool	
	Report string
	Data string
}

//Table prints a cable table
func (c *Cbl) Table(printz bool){
	rstr := &strings.Builder{}
	rstr.WriteString(ColorGreen)
	t0 := tablewriter.NewWriter(rstr)
	t0.SetHeader([]string{"length","area","dia","xr","yr"})
	row := fmt.Sprintf("%f, %f, %f, %f, %f",c.Ucl,c.Ar,c.Dia,c.Xr,c.Yr)
	t0.Append(strings.Split(row,","))
	t0.Render()
	t0 = tablewriter.NewWriter(rstr)
	t0.SetHeader([]string{"em","alpha","self wt","temp delta"})
	t0.SetCaption(true, "cable geom/prop")
	row = fmt.Sprintf("%f, %f, %f, %f",c.Em,c.Alp,c.Sw,c.Tr)
	t0.Append(strings.Split(row,","))
	t0.Render()
	rstr.WriteString(ColorCyan)
	t0 = tablewriter.NewWriter(rstr)
	t0.SetCaption(true, "applied forces")
	t0.SetHeader([]string{"no.","fx","fy","dist"})
	for i, frc := range c.Lds{
		row = fmt.Sprintf("%v, %f, %f, %f",i+1,frc[0],frc[1],frc[2])
		t0.Append(strings.Split(row,","))
	}
	t0.Render()
	rstr.WriteString(ColorReset)
	c.Report = fmt.Sprintf("%s",rstr)
	if printz{
		fmt.Println(c.Report)
	}
	return
}

//CalcL - cable of constant length analysis
func (c *Cbl) CalcL() (err error){
	//constant length cable analysis
	c.Lseg = c.Ucl/float64(c.Nl)
	c.Xd = make([]float64,c.Nl+1)
	c.Yd = make([]float64, c.Nl+1)
	c.El = make([]float64, c.Nl+1)
	var iter, kiter int
	tl := c.Tl
	sl := c.Theta * math.Pi/180.0
	ea := c.Em * c.Ar
	swt := c.Lseg * c.Sw
	//init vals
	c.Ts[0] = tl	
	c.Xd[0] = 0.0
	c.Yd[0] = 0.0
	c.El[0] = c.Lseg * (1.0 + tl/ea + c.Alp * c.Tr)
	c.Xd[1] = c.El[0] * math.Cos(sl)
	c.Yd[1] = c.El[0] * math.Sin(sl)
	
	for iter == 0{
		kiter++
		if kiter > 666{
			err = fmt.Errorf("iteration error")
			return
		}
		fmt.Println("kiter, xr, yr")
		for i := 2; i <= c.Nl; i++{
			pt := 0.0; ph := 0.0
			dx := float64(i) * c.Lseg
			for j, lc := range c.Lds{
				if len(lc) < 3{
					log.Printf("invalid load case number %v - %v\n",j+1, lc)
					continue
				}
				dlc := lc[2]
				if c.Lseg - math.Abs(dx - dlc) > 0.0{
					pt = pt + (1.0 - math.Abs(dx - dlc)/c.Lseg) * lc[0]
					ph = pt + (1.0 - math.Abs(dx - dlc)/c.Lseg) * lc[1]
				}
			}
			tx := c.Ts[i-1] * (c.Xd[i] - c.Xd[i-1])/c.El[i-1] - ph 
			ty := c.Ts[i-1] * (c.Yd[i] - c.Yd[i-1])/c.El[i-1] - swt - pt
			tang := ty/tx
			ang := math.Atan(tang)
			c.Ts[i] = math.Sqrt(tx * tx + ty * ty)
			c.El[i] = c.Lseg * (1.0 + c.Ts[i]/ea + c.Alp * c.Tr)
			delx := c.El[i] * math.Cos(ang)
			dely := c.El[i] * math.Sin(ang)
			fmt.Println(delx, dely)
		}
	}
	c.Table(c.Verbose)
	return
}

//CalcT - cable of constant tension analysis
func (c *Cbl) CalcT() (err error){
	//constant tension cable analysis
	return
}
