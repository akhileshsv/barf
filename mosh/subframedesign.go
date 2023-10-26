package barf

import (
	"fmt"
	//"log"
	//"math"
	//kass"barf/kass"
	//"github.com/AlecAivazis/survey/v2"
)

//SubFlr stores a subframe and slice of slabs
//doesn't look like it'll be used
type SubFlr struct{
	Sf *SubFrm
	Slbz []*RccSlb
	Nbays int
	Slbtyp int
}

//DzSubFrm designs a SubFrm (columns/beams)
func DzSubFrm(sf *SubFrm) (err error){
	//first get x
	//then y
	//lx and ly are trib. widths
	//check ldtyps here
	switch sf.Dtyp{
		case 2:
		//start from roof
		case 1:
		case 0:
		err = CalcSubFrm(sf)
		if err != nil{
			return
		}
		bmchn := make(chan []interface{},len(sf.Beams))
		colchn := make(chan []interface{},len(sf.Cols))
		//Bmdz(code int, barr []*RccBm, bm *kass.BmEnv, term string, verbose bool, bmchn chan []interface{}){
		for i, bm := range sf.Bmenv{
			go Bmdz(sf.Code, sf.RcBm[i], bm, sf.Term, sf.Verbose, bmchn)
		}
		//ColDz(c *RccCol, colchn chan []interface{})
		//MIGHT BREAK - why
		for i := 1; i <= len(sf.Cols); i++{
		//for i := range sf.Colenv{
			if sf.Tweakc{
				sf.RcCol[i].Tweaks()
			}
			go ColDz(sf.RcCol[i], colchn)
		}
		for _ = range sf.Bmenv{
			_ = <- bmchn
			
		}
		for _ = range sf.Colenv{
			//ColDz(c *RccCol, colchn chan []interface{})
			_ = <- colchn
		}
		for _, bmarr := range sf.RcBm{
			for _, b := range bmarr{
				if sf.Tweakb{
					b.Tweaks()
				}
				PlotBmGeom(b, b.Term)
			}
		}
	}
	return
}

//getsx returns a subframe in either x (using Lspans) or y (using Lbays) directions
func(sf *SubFrm) getsx(ax int)(sx *SubFrm){
	switch ax{
		case 1:
		//x sub frame
		sx = &SubFrm{}
		*sx = *sf
		sx.Title = fmt.Sprintf("%s_x",sf.Title)
		sx.Lbay = sf.Ly
		sx.Ax = ax
		case 2:
		//y sub frame
		sx = &SubFrm{}
		*sx = *sf
		sx.Title = fmt.Sprintf("%s_y",sf.Title)
		sx.Lbay = sf.Lx
		sx.Ax = ax
		sx.Lspans = sf.Lbays
	}
	return
}

//ChainSlb designs a slab and uses the slab depth and loading
//to analyze and design a SubFrm
//use CHAIN more often; write more of such funcs
func (sf *SubFrm) ChainSlb(s *RccSlb)(err error){
	err = SlbDesign(s)
	if err != nil{
		fmt.Println(err)
		return
	}
	s.Table(true)
	fmt.Println("behold, dused and dl before->",sf.Dslb, sf.DL)
	switch s.Type{
		case 1, 2:
		sf.Dslb = s.Dused
		sf.Slbload = s.Type
		fmt.Println("after",sf.Dslb, sf.Slbload)
		case 3:
		case 4:
	}
	err = DzSubFrm(sf)
	return
}

//add - save beam env

/*
func DzSubFrmBm(sf *SubFrm, term string, menu bool){
	fmt.Println("starting subfrm beam design")
	bmchn := make(chan []interface{}, len(sf.Bmenv))
	for _, bm := range sf.Beams{
		fmt.Println("beam no->",bm)
		go Bmdz(sf.Code, sf.RcBm[bm], sf.Bmenv[bm], term, sf.Verbose, bmchn)
		_ = <- bmchn
	}
	//for _ = range sf.RcBm{
	//	_ = <- bmchn
	//}
	supz := []string{"left support","midspan","right support"}
	for _, bm := range sf.Beams{
		fmt.Println("bm id->", bm)
		switch sf.Bmenv[bm].Endc{
			case 0:
			fmt.Println(sf.RcBm[bm][1].Ast, sf.RcBm[bm][1].Asc)
			default:
			for j := 0; j < 3; j++{
				fmt.Println(supz[j])
				fmt.Println("ast->", sf.RcBm[bm][j].Ast, " mm2 asc->",sf.RcBm[bm][j].Asc," mm2")
				
			}

		}
	}
}

func DzSubFrmCol(sf *SubFrm, term string, menu bool){
	log.Println("subfrm column")
	colchn := make(chan []interface{}, len(sf.Cols))
	for _, c := range sf.Cols{
		go Coldz(c, sf, colchn)
		_ = <- colchn
	}
	//for _ = range sf.Cols{
	//	_ = <- colchn
	//}
	return
}
//DUDE. ALL SUB FRAMES ARE X AND Y FRAMES

func OptSubFrm(sf SubFrm){
	//first check for type of subframe
	//
}
*/
