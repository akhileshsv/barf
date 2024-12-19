package barf

import (
	"fmt"
	"math"
	kass"barf/kass"
	"github.com/AlecAivazis/survey/v2"
)

//GetCol reads a ColEnv struct and grade params and returns a RccCol struct 
func GetCol(cenv *kass.ColEnv, kostin []float64,fck, fy, efcvr float64, frmtyp, code int) (c RccCol){
	//generates a column given params from anything
	//GetCol(f.Colenv[i], rcol[i], f.Fcks[1], f.Code, f.Styps[cp-1])
	dtyp := 2; rtyp := 0; rbrlvl := 0
	if frmtyp == 2 || frmtyp == 1{
		dtyp = 1
	}
	//efcvr = efcvr + 8.0
	//fmt.Println("efcvr->",efcvr)
	//l0 := kass.EffHt(cenv)
	c = RccCol{
		Id:cenv.Id,
		Mid:cenv.Id,
		Fck:fck,
		Fy:fy,
		Cvrt:efcvr,
		Cvrc:efcvr,
		Code:code,
		Styp:cenv.Styp,
		Lspan:cenv.Lspan,
		Dtyp:dtyp,
		Rtyp:rtyp,
		Rbrlvl:rbrlvl,
		Nlayers:2,
		Ast:0.0,
		Asc:0.0,
		Asteel:0.0,
		Subck:false,
		L0:cenv.L0,
		Leffx:cenv.Le,
		Term:cenv.Term,
		Foldr:cenv.Foldr,
		Title:fmt.Sprintf("col-%s-%v",cenv.Title,cenv.Id),
		Web:cenv.Web,
		Ignore:cenv.Ignore,
	}
	if cenv.Ignore{return}
	c.Dims = make([]float64, len(cenv.Dims))
	copy(c.Dims, cenv.Dims)
	
	c.Coords = make([][]float64, 2)
	copy(c.Coords, cenv.Coords)
	switch frmtyp{
		case 1,2:
		//subfrm, 2d frm
		c.Pu = math.Abs(cenv.Pumax)
		if cenv.Pufac > 0.0{
			c.Pu = c.Pu * cenv.Pufac
		}
		// if cenv.Puadd > 0.0{
		// 	c.Pu += cenv.Puadd
		// }
		//if frmtyp == 1 && cenv.Pufac > 1.0{c.Pu = c.Pu * cenv.Pufac}
		c.Mux = math.Abs(cenv.Mtmax)
		//fmt.Println(ColorYellow,"col->",c.Id,"pu->",c.Pu,"mux->",c.Mux,ColorReset)
		if cenv.Mbmax > c.Mux{
			c.Mux = math.Abs(cenv.Mbmax)
		}
		
		ecct := c.H/20.0
		if ecct > 20.0 {ecct = 20.0}
		if c.Mux < c.Pu * ecct * 1e-3 {c.Mux = c.Pu * ecct * 1e-3}
		//additional moment - ned*l0/400.0
		case 3:
		//3d frm
	}
	switch c.Styp{
		case 0:
		c.B = c.Dims[0]
		case 1:
		c.B = c.Dims[0]
		c.H = c.Dims[1]
	}
	//fmt.Println("hello from col d->\nc.Id,c.Mux, c.Pu, c.Dims, c.Pufac, cenv.Pumax\n",c.Id,c.Mux, c.Pu, c.Dims, cenv.Pufac, cenv.Pumax)

	//c.Init()
	//c.SecInit()
	return
}

//Tweaks are for column design tweaks
//TODO - this could be useful
func (c *RccCol) Tweaks(){
	var choice, iter int
	for iter != -1{
		prompt := &survey.Select{
			Message: fmt.Sprintf("%sChoose param to tweak for col_%v%s",ColorRed,c.Id,ColorReset),
			Options: []string{"pu", "pufac","muy","continue"},
		}
		survey.AskOne(prompt, &choice)
		switch choice{
			case 0:
			var pu float64
			pr1 := &survey.Input{
				Message: fmt.Sprintf("Enter pu in kn for col_%v",c.Id),
			}
			survey.AskOne(pr1, &pu)
			c.Pu = float64(pu)
			//iter = -1
			case 1:
			var puf float64
			pr1 := &survey.Input{
				Message: fmt.Sprintf("Enter pu multiplier for col_%v:",c.Id),
			}
			survey.AskOne(pr1, &puf)
			c.Pu = puf * c.Pu
			//iter = -1
			case 2:
			var muy float64
			pr1 := &survey.Input{
				Message: fmt.Sprintf("Enter moment about y-axis for col_%v:",c.Id),
			}
			survey.AskOne(pr1, &muy)
			c.Muy = muy
			//iter = -1
			case 3:
			iter = -1
		}
	}
	return
}

//ColDz is the goroutine entry func for rcc column design
//from frm2d/subframe design routines
func ColDz(c *RccCol, colchn chan []interface{}){
	var iter int
	var err error
	//fmt.Println("hello from dz->",c.Id,c.Mux, c.Pu, c.Dims)
	var astmin float64
	var minval int
	astmin = 6e66
	switch c.Styp{
		case 0 , 1:		
		for c.Nlayers <= 6{//iter == 0 && 
			if c.Nlayers < 2 {c.Nlayers = 2}
			err = ColDesign(c)
			if err == nil{
				//fmt.Println("id, nlay, asteel->",c.Id, c.Nlayers, c.Asteel)
				if astmin > c.Asteel{
					astmin = c.Asteel
					minval = c.Nlayers
				}
				//iter = 1
				//break
			}
			c.Nlayers += 1
			//fmt.Println(err)
		}
		
		//if err != nil && !c.Approx{
		//	c.Approx = true
			
		//}
		if astmin == 6e66 && !c.Approx{	
			//fmt.Println("setting approx->",c.Id)
			c.Approx = true
			c.Subck = false
			c.Nlayers = 2	
			for c.Nlayers <= 6{
				if c.Nlayers < 2 {c.Nlayers = 2}
				err = ColDesign(c)
				if err == nil{
					//iter = 1
					//break
					if astmin > c.Asteel{
						astmin = c.Asteel
						minval = c.Nlayers
					}
				} 
				c.Nlayers += 1
			}
			//fmt.Println("found via approx", astmin, c.Id)
			
		}
		if astmin != 6e66{
			c.Nlayers = minval
			c.Asteel = 0.0; c.Ast = 0.0; c.Asc = 0.0
			err = ColDesign(c)
		}
		default:
		for iter == 0 && c.Rbrlvl < 2{
			err = ColDesign(c)
			if err == nil{
				iter = 1
				break
			}
			c.Rbrlvl++
		}
		
		if err != nil && !c.Approx{
			c.Approx = true
			
		}
		for iter == 0 && c.Rbrlvl < 2{
			err = ColDesign(c)
			if err == nil{
				iter = 1
				break
			}
			c.Rbrlvl++
		}
		
	}
	//fmt.Println("hello from dz->",c.Id,c.Mux, c.Pu, c.Dims)
	
	//fmt.Println(ColorGreen,"astmin, minval->",c.Id,astmin,minval,ColorReset)
	rez := make([]interface{}, 2)
	rez[0] = c.Id
	rez[1] = err
	colchn <- rez
}

/*
func Coldz(c int, sf *SubFrm, colchn chan []interface{}){
	//might turn into main column entry func IT HASNT
	//colchn = [col dx, error, cons]
	//var braced, ljbase, basemr, ljbeamss, ujbeamss bool
	rez := make([]interface{},3)
	var err error
	var cons []int
	cenv := sf.Colenv[c]
	col := sf.RcCol[c]
	//xdx, fdx, locx, lex := sf.Members[col]
	fmt.Println(ColorBlue,"col->",cenv.Id,ColorReset)
	log.Println(ColorCyan,"moment top->",cenv.Mtmax,"moment bottom->",cenv.Mbmax)
	log.Println(ColorRed,"axial load",cenv.Ptmax, cenv.Pbmax, ColorReset)
	//calc effective height in X
	//switch xdx{
		
	//}
	switch sf.Code{
		case 1:
		switch col.Styp{
			case 0:
			//get equivalent square column
			
		}
		case 2:
	}
	rez[0] = c
	rez[1] = err
	rez[2] = cons
	colchn <- rez
}
*/
