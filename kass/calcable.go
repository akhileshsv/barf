package barf

import (
	"fmt"
	"log"
	"math"
	"strings"
	"github.com/olekukonko/tablewriter"	
)

//Cbl is a cable struct
type Cbl struct{
	//cable 
	Title string
	Ucl float64 //unstretched cable length
	Em float64 //em
	Dia float64 //dia
	Ar float64 //area
	Alp float64 //coeff of expansion
	Xr float64 //x right (xl, yl = 0,0)
	Yr float64 //y right 
	Tr float64 //temp rise
	Sw float64 //self wt/unit length
	Lds [][]float64 //loads along length [fy, fx, dist]
	Nl int //no. of elements
	Lseg float64 //segment length
	Tl float64 //tension at left sup
	Theta float64 //slope at left sup
	Dip float64 //cable dip
	Idip int //max dip index
	Idep int //max dep index
	Dep float64 //cable depression
	Tsl float64 //tot. stretched length
	Lht float64 //left hand tension
	Lha float64 //left hand angle
	Rht float64 //right hand tension
	Rha float64 //right hand angle
	Xd []float64
	Yd []float64
	El []float64
	Ts []float64
	Pts [][]float64
	Web bool
	Dz bool
	Verbose bool	
	Report string
	Data string
	Term string
	Txtplot string
}


//Draw draws a cable
func (c *Cbl) Draw(){
	if c.Term == ""{c.Term = "qt"}
	var data, mdata string
	for i, xd := range c.Xd{
		yd := -c.Yd[i]
		data += fmt.Sprintf("%f %f\n",xd, yd)
		if i < c.Nl{
			xd2 := c.Xd[i+1]
			yd2 := -c.Yd[i+1]
			mdata += fmt.Sprintf("%f %f %f %f\n",xd, yd, xd2, yd2)
		}
	}
	data += "\n\n"; mdata += "\n\n"; data += mdata
	data += fmt.Sprintf("%f %f\n",c.Xd[0],c.Yd[0])
	data += fmt.Sprintf("%f %f\n",c.Xr,-c.Yr)
	data += "\n\n"
	for _, fvec := range c.Lds{
		fx := fvec[1]
		fy := fvec[0]
		x1 := fvec[2]
		i1 := x1/c.Lseg
		x1 = c.Xd[int(i1)]
		y1 := -c.Yd[int(i1)]
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x1,y1,fx,0.0)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x1,y1,0.0,-fy)
	}
	data += "\n\n"
	fname := fmt.Sprintf("cable-%s",c.Title)
	title := fname
	skript := "drawmod2d.gp"
	txtplot := skriptrun(data, skript, c.Term, title,"",fname)
	if c.Web{
		switch c.Term{
			case "dxf":		
			txtplot = fname + ".dxf"
			case "svg", "svgmono":
			Svgkong(txtplot)
			txtplot = fname + ".svg"
		}
	}
	if c.Term == "dumb"{fmt.Println(txtplot)}
}

//Table prints a cable table
func (c *Cbl) Table(printz bool){
	rstr := &strings.Builder{}
	rstr.WriteString(ColorYellow)
	hdr := fmt.Sprintf("\ncable analysis %s\n",c.Title)
	rstr.WriteString(hdr)
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
		row = fmt.Sprintf("%v, %f, %f, %f",i+1,frc[1],frc[0],frc[2])
		t0.Append(strings.Split(row,","))
	}
	t0.Render()
	rstr.WriteString(ColorPurple)
	t0 = tablewriter.NewWriter(rstr)
	t0.SetHeader([]string{"node","x","y","elem","tension"})
	t0.SetCaption(true, "cable profile")
	for i := 0; i < c.Nl+1; i += 9{
		row = fmt.Sprintf("%v, %f, %f, %v, %f",i+2,c.Xd[i],c.Yd[i],i+1,c.Ts[i])	
		t0.Append(strings.Split(row,","))
	}
	t0.Render()
	rstr.WriteString(ColorBlue)
	t0 = tablewriter.NewWriter(rstr)
	t0.SetHeader([]string{"t(lh)","ang(lh)","t(rh)","ang(rh)","dip","at","depr.","at","def.len"})
	t0.SetCaption(true, "cable vals")
	row = fmt.Sprintf("%f,%f,%f,%f,%f,%v,%f,%v,%f",c.Lht, c.Lha, c.Rht, c.Rha, c.Dip, c.Idip, c.Dep, c.Idep, c.Tsl)	
	t0.Append(strings.Split(row,","))
	t0.Render()
	rstr.WriteString(ColorReset)
	c.Report = rstr.String()//fmt.Sprintf("%s",rstr)
	if printz{
		fmt.Println(c.Report)
	}
}

//CalcCl - cable of constant length analysis
func (c *Cbl) CalcCl() (err error){
	//constant length cable analysis
	if c.Nl == 0{c.Nl = 100}
	c.Lseg = c.Ucl/float64(c.Nl)
	c.Xd = make([]float64, c.Nl+1)
	c.Yd = make([]float64, c.Nl+1)
	c.El = make([]float64, c.Nl+1)
	c.Ts = make([]float64, c.Nl+1)
	var iter, niter, kiter int
	//init vals? 
	tl := c.Tl
	theta := c.Theta
	sl := theta * math.Pi/180.0
	ea := c.Em * c.Ar
	swt := c.Lseg * c.Sw
	z := make([]float64,8)
	//init vals
	c.Ts[0] = tl	
	c.Xd[0] = 0.0
	c.Yd[0] = 0.0
	c.El[0] = c.Lseg * (1.0 + tl/ea + c.Alp * c.Tr)
	c.Xd[1] = c.El[0] * math.Cos(sl)
	c.Yd[1] = c.El[0] * math.Sin(sl)
	for iter == 0{
		if kiter > 30{
			err = fmt.Errorf("iteration error")
			return
		}
		//start from node no. 2
		for i := 2; i < c.Nl+1; i++{
			j := i - 1
			pt := 0.0; ph := 0.0
			dx := float64(j) * c.Lseg
			for k, lc := range c.Lds{
				if len(lc) < 3{
					log.Printf("invalid load case number %v - %v\n",k+1, lc)
					continue
				}
				dlc := lc[2]
				if c.Lseg - math.Abs(dx - dlc) > 0.0{
					pt = pt + (1.0 - math.Abs(dx - dlc)/c.Lseg) * lc[0]
					ph = ph + (1.0 - math.Abs(dx - dlc)/c.Lseg) * lc[1]
				}
			}
			tx := c.Ts[j-1] * (c.Xd[j] - c.Xd[j-1])/c.El[j-1] - ph 
			ty := c.Ts[j-1] * (c.Yd[j] - c.Yd[j-1])/c.El[j-1] - swt - pt
			tang := ty/tx
			ang := math.Atan(tang)
			c.Ts[j] = math.Sqrt(tx * tx + ty * ty)
			c.El[j] = c.Lseg * (1.0 + c.Ts[j]/ea + c.Alp * c.Tr)
			delx := c.El[j] * math.Cos(ang)
			dely := c.El[j] * math.Sin(ang)
			c.Xd[i] = c.Xd[j] + delx
			c.Yd[i] = c.Yd[j] + dely
		}
		xmis := c.Xr - c.Xd[c.Nl]
		ymis := c.Yr - c.Yd[c.Nl]
		errore := math.Sqrt(xmis*xmis + ymis*ymis)
		//fmt.Println(ColorCyan,"cycle",kiter,"error",errore,"lh tension",tl,"angle",theta,ColorReset)
		if errore < 0.0001{
			iter = -1
			break
		}
		switch niter{
			case 0:
			//first run
			z[0] = tl
			z[1] = theta
			z[2] = xmis
			z[3] = ymis	
			tl = tl * 1.01	
			niter = 1
			case 1:
			//second
			z[4] = xmis
			z[5] = ymis
			theta = 1.01 * theta
			tl = z[0]
			niter = 2
			case 2:
			//update
			z[6] = xmis
			z[7] = ymis
			dxt := (z[4] - z[2])/(0.01 * z[0])
			dyt := (z[5] - z[3])/(0.01 * z[0])
			dxa := (z[6] - z[2])/(0.01 * z[1])
			dya := (z[7] - z[3])/(0.01 * z[1])
			den := dxt * dya - dyt * dxa
			tl = z[0] - (dya * z[2] - dxa * z[3])/den
			theta = z[1] - (dxt * z[3] - dyt * z[2])/den
			niter = 0
		}
		sl = theta * math.Pi/180.0
		c.Ts[0] = tl	
		c.Xd[0] = 0.0
		c.Yd[0] = 0.0
		c.El[0] = c.Lseg * (1.0 + tl/ea + c.Alp * c.Tr)
		c.Xd[1] = c.El[0] * math.Cos(sl)
		c.Yd[1] = c.El[0] * math.Sin(sl)	
		kiter++
	}
	c.Lht = c.Ts[0]; c.Lha = theta
	c.Rht = c.Ts[c.Nl-1]
	c.Rha = (c.Yd[c.Nl-1] - c.Yd[c.Nl])/(c.Xd[c.Nl-1]-c.Xd[c.Nl])
	c.Rha = math.Atan(c.Rha) * 180.0/math.Pi
	for i, val := range c.El{
		c.Tsl += val
		if c.Dep < c.Yd[i]{
			c.Dep = c.Yd[i]; c.Idep = i
		}
		dip := c.Yd[i] - c.Xd[i] * c.Yr/c.Xr
		if c.Dip < dip{
			c.Dip = dip; c.Idip = i
		}
	}
	c.Table(c.Verbose)
	c.Draw()
	return
}

//CalcT - cable of constant tension analysis
func (c *Cbl) CalcT() (err error){
	//constant tension cable analysis
	return
}
