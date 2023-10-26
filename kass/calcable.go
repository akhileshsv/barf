package barf

import (
	"fmt"
	//"math"
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
	//var iter int
	
	c.Table(c.Verbose)
	return
}

//CalcT - cable of constant tension analysis
func (c *Cbl) CalcT() (err error){
	//constant tension cable analysis
	return
}
