package barf

import (
	"log"
	"fmt"
	"math"
)

//harrison 15.1 (cable3) cable analysis/design
func (c *Cbl) CblAz()(err error){
	//constant length cable analysis; harrison 15.3
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
	f6 := c.Flh + c.Fls
	f7 := c.Flv + c.Fls
	f8 := c.Frh + c.Frs
	f9 := c.Frv + c.Frs
	for iter == 0{
		if kiter > 30{
			err = fmt.Errorf("iteration error")
			return
		}
		//init vals HERE
		c.Wt = 0.0
		c.Vol = 0.0
		sl = theta * math.Pi/180.0
		c.Ts[0] = tl	
		c.Xd[0] = tl * math.Cos(sl) * f6
		c.Yd[0] = tl * math.Sin(sl) * f7
		c.El[0] = c.Lseg * (1.0 + tl/ea + c.Alp * c.Tr)
		c.Xd[1] = c.El[0] * math.Cos(sl) + c.Xd[0]
		c.Yd[1] = c.El[0] * math.Sin(sl) + c.Yd[0]
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
				ltyp := 1
				switch len(lc){
					case 4:
					ltyp = int(lc[3])
				}
				switch ltyp{
					case 1:
					//fixed load
					dlc := lc[2]
					if c.Lseg - math.Abs(dx - dlc) > 0.0{
						pt = pt + (1.0 - math.Abs(dx - dlc)/c.Lseg) * lc[0]
						ph = ph + (1.0 - math.Abs(dx - dlc)/c.Lseg) * lc[1]
					}
					case 2:
					//rolling load	
					dx := c.Xd[j] - lc[2]
					dc := c.Xd[j] - c.Xd[j-1]
					
					if dc - math.Abs(dx) > 0.0{
						pt = pt + (1.0 - math.Abs(dx)/dc) * lc[0]
						ph = ph + (1.0 - math.Abs(dx)/dc) * lc[1]
					}
					default:
					log.Printf("invalid load type %v - %v",k+1, ltyp)
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
			c.Vol += c.El[j] * c.Ar
		}
		xmis := c.Xr - c.Xd[c.Nl] - f8 * c.Ts[c.Nl] * (c.Xd[c.Nl] - c.Xd[c.Nl-1])/c.El[c.Nl-1]
		ymis := c.Yr - c.Yd[c.Nl] + f9 * c.Ts[c.Nl] * (c.Yd[c.Nl-1] - c.Yd[c.Nl])/c.El[c.Nl-1]
		errore := math.Sqrt(xmis*xmis + ymis*ymis)
		//CHANGE THIS, percent of ucl or something
		if errore < 0.0005{
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
		// sl = theta * math.Pi/180.0
		// c.Ts[0] = tl	
		// c.Xd[0] = 0.0
		// c.Yd[0] = 0.0
		// c.El[0] = c.Lseg * (1.0 + tl/ea + c.Alp * c.Tr)
		// c.Xd[1] = c.El[0] * math.Cos(sl)
		// c.Yd[1] = c.El[0] * math.Sin(sl)	
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

//harrison 15.2 (cable4) constant stress suspension cable design
func (c *Cbl) CblDzCs()(err error){
	if c.Dtyp == 0{c.Dtyp = 3}
	if c.Nl == 0{c.Nl = 100}
	c.Lseg = c.Ucl/float64(c.Nl)
	c.Xd = make([]float64, c.Nl+1)
	c.Yd = make([]float64, c.Nl+1)
	c.El = make([]float64, c.Nl+1)
	c.Ts = make([]float64, c.Nl+1)
	c.Ars = make([]float64, c.Nl+1)
	var iter, niter, kiter int
	//init vals? 
	tl := c.Tl
	theta := c.Theta
	sl := theta * math.Pi/180.0
	z := make([]float64,8)
	fxl := c.Flh + c.Fls
	fyl := c.Flv + c.Fls
	fxr := c.Frh + c.Frs
	fyr := c.Frv + c.Frs
	for iter == 0{
		if kiter > 30{
			err = fmt.Errorf("iteration error")
			return
		}
		//init vals HERE
		c.Vol = 0.0
		sl = theta * math.Pi/180.0
		switch len(c.Lds){
			case 0:
			c.Ars[0] = 1.0
			c.Strs = tl
			default:
			if !c.Lockar{c.Ars[0] = tl/c.Strs}
		}
		c.Ts[0] = tl	
		c.Xd[0] = tl * math.Cos(sl) * fxl
		c.Yd[0] = tl * math.Sin(sl) * fyl
		c.El[0] = c.Lseg * (1.0 + tl/(c.Em*c.Ars[0]) + c.Alp * c.Tr)
		c.Xd[1] = c.El[0] * math.Cos(sl) + c.Xd[0]
		c.Yd[1] = c.El[0] * math.Sin(sl) + c.Yd[0]
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
				ltyp := 1
				switch len(lc){
					case 4:
					ltyp = int(lc[3])
				}
				switch ltyp{
					case 1:
					//fixed load
					dlc := lc[2]
					if c.Lseg - math.Abs(dx - dlc) > 0.0{
						pt = pt + (1.0 - math.Abs(dx - dlc)/c.Lseg) * lc[0]
						ph = ph + (1.0 - math.Abs(dx - dlc)/c.Lseg) * lc[1]
					}
					case 2:
					//rolling load	
					dx := c.Xd[j] - lc[2]
					dc := c.Xd[j] - c.Xd[j-1]
					
					if dc - math.Abs(dx) > 0.0{
						pt = pt + (1.0 - math.Abs(dx)/dc) * lc[0]
						ph = ph + (1.0 - math.Abs(dx)/dc) * lc[1]
					}
					default:
					log.Printf("invalid load type %v - %v",k+1, ltyp)
				}
			}
			tx := c.Ts[j-1] * (c.Xd[j] - c.Xd[j-1])/c.El[j-1] - ph 
			ty := c.Ts[j-1] * (c.Yd[j] - c.Yd[j-1])/c.El[j-1] - pt
			ty = ty - c.Dens * c.Ars[j-1] * c.Lseg
			tang := ty/tx
			ang := math.Atan(tang)
			c.Ts[j] = math.Sqrt(tx * tx + ty * ty)
			if !c.Lockar{c.Ars[j] = c.Ts[j]/c.Strs}
			c.El[j] = c.Lseg * (1.0 + c.Ts[j]/(c.Em*c.Ars[j]) + c.Alp * c.Tr)
			delx := c.El[j] * math.Cos(ang)
			dely := c.El[j] * math.Sin(ang)
			c.Xd[i] = c.Xd[j] + delx
			c.Yd[i] = c.Yd[j] + dely
			c.Vol += c.El[j] * c.Ars[j]
		}
		xmis := c.Xr - c.Xd[c.Nl] - fxr * c.Ts[c.Nl] * (c.Xd[c.Nl] - c.Xd[c.Nl-1])/c.El[c.Nl-1]
		ymis := c.Yr - c.Yd[c.Nl] + fyr * c.Ts[c.Nl] * (c.Yd[c.Nl-1] - c.Yd[c.Nl])/c.El[c.Nl-1]
		errore := math.Sqrt(xmis*xmis + ymis*ymis)
		//CHANGE THIS, percent of ucl or something
		if errore < 0.0005{
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
