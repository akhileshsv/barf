package barf

import (
	"fmt"
	"strings"
	"github.com/olekukonko/tablewriter"
)


//Table prints a table summary of a designed subframe
//and plots col/beam sections
func (sf *SubFrm) Table(printz bool){
	cchn := make(chan string,1)
	bchn := make(chan string, 1)
	
	cpchn := make(chan []string, 1)
	bpchn := make(chan []string,1)
	
	go coltab(sf.Spam, sf.Cols, sf.RcCol, cchn)
	go bmtab(sf.Spam, sf.Beams, sf.RcBm, bchn)
	if sf.Term != ""{
		//go kass.Drawflps(f, lpchn)
		go drawcols(sf.Web, sf.Cols, sf.RcCol, cpchn)
		go drawbms(sf.Web, sf.Beams, sf.RcBm, bpchn)
		//go drawcbms(f.Web, f.Term, f.Title, f.CBeams, rcbm, cbchn)
	}
	rezstr := new(strings.Builder)
	rezstr.WriteString(ColorPurple)
	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"vol tot(m3)","vol rcc(m3)","wt stl(kg)","form area (m2)","cost (rs)"})
	table.SetCaption(true,"sub frame -> quantity take off")
	row := fmt.Sprintf("%.3f, %.3f, %.3f, %.3f, %.f\n",sf.Quants[0],sf.Quants[0],sf.Quants[1],sf.Quants[2], sf.Kost)
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorReset)
	r := rezstr.String()
	
	ctab := <- cchn
	sf.Report += ctab
	btab := <- bchn
	sf.Report += btab
	sf.Report += r
	
	if sf.Term != ""{
		cplts := <- cpchn
		sf.Txtplots = append(sf.Txtplots, cplts...)
		bplts := <-bpchn
		sf.Txtplots = append(sf.Txtplots, bplts...)
	}
	if printz{
		fmt.Println(sf.Report)
	}
}



//DzSubFrm designs a SubFrm (columns/beams)
func DzSubFrm(sf *SubFrm) (allcons int, err error){
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
		// for range sf.Bmenv{
		// 	<- bmchn
			
		// }
		// for range sf.Colenv{
		// 	//ColDz(c *RccCol, colchn chan []interface{})
		// 	<- colchn
		// }
		
		// emap := make(map[int][]error)
		sf.Quants = make([]float64, 3)
		//fmt.Println(len(sf.Bmenv),len(sf.Colenv))
		
		var errstr string
		for range sf.Bmenv{
			rez := <- bmchn
			idx, _ := rez[0].(int)
			cons, _ := rez[1].(int)
			errz, _ := rez[2].([]error)
			allcons += cons
			if cons > 0{errstr += fmt.Sprintf("beam span no %v - cons - %v - \nerrz - %s\n",idx,cons,errz)}
			//emap[idx] = append(emap[idx],errz...)
			sf.Quants[0] += sf.RcBm[idx][1].Vrcc
			sf.Quants[1] += sf.RcBm[idx][1].Wstl
			sf.Quants[2] += sf.RcBm[idx][1].Afw
			
		}
		for range sf.Colenv{
			rez := <- colchn 		
			idx, _ := rez[0].(int)
			err, _ := rez[1].(error)
			//emap[idx] = append(emap[idx],err)
			sf.Quants[0] += sf.RcCol[idx].Vrcc
			sf.Quants[1] += sf.RcCol[idx].Wstl
			sf.Quants[2] += sf.RcCol[idx].Afw
			if err != nil{
				allcons++
				errstr += fmt.Sprintf("col span no %v - \nerr - %s\n",idx,err)
			}
		}
		//fmt.Println(allcons)
		if allcons != 0{
			// var errstr string
			// errstr += fmt.Sprintf("design error\nnumber of failed constraints-%v\n",allcons)
			// for i, evec := range emap{
			// 	errstr += fmt.Sprintf("member %v error list->\n",i)
			// 	for j, err := range evec{
			// 		errstr += fmt.Sprintf("span-%v\nerror->%s\n",j,err)
			// 	}
			// }
			// err = fmt.Errorf("%s",errstr)
			err = fmt.Errorf("total constraints violated - %v\n errors - %s",allcons, errstr)
			return
		}
		sf.Dz = true
		switch len(sf.Kostin){
			case 3:
			sf.Kost = sf.Quants[0] * sf.Kostin[0] + sf.Quants[1] * sf.Kostin[1] + sf.Quants[2] * sf.Kostin[2]
			if sf.Opt > 0{
				sf.Kost = sf.Kost/sf.Kostin[0]
			}
			default:
			sf.Kost = sf.Quants[0] * CostRcc + sf.Quants[1] * CostStl + sf.Quants[2] * CostForm
			if sf.Opt > 0{
				sf.Kost = sf.Kost/CostRcc
			}
		}
		// for _, bmarr := range sf.RcBm{
		// 	for _, b := range bmarr{
		// 		if sf.Tweakb{
		// 			b.Tweaks()
		// 		}
		// 		PlotBmGeom(b, b.Term)
		// 	}
		// }
	}
	if sf.Verbose{
		var printz bool
		if !sf.Web{
			printz = true
		}
		sf.Table(printz)
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

//ChainSlab designs a slab and adds slab loading to subframe
func (sf *SubFrm) ChainSlab()(s RccSlb, err error){
	s.Fck = sf.Fcks[0]
	s.Fy = sf.Fys[0]
	s.Fyd = s.Fy
	s.DL = sf.DL
	s.LL = sf.LL
	s.Title = sf.Title + "-slab"
	s.Web = sf.Web
	s.Term = sf.Term
	s.Code = sf.Code
	switch sf.Exp{
		case 0:
		//mild exposure
		s.Nomcvr = 15.0
		case 1:
		//moderate exposure
		s.Nomcvr = 20.0
		case 2:
		//severe exposure
		s.Nomcvr = 25.0
	}
	switch sf.Slbload{
		case 1:
		s.Type = 1
		s.Typstr = "1wcs"
		//one way slab
		switch len(sf.Lbays){
			case 0:
			if sf.Lbay == 0.0{
				err = fmt.Errorf("no bays specified for slab design- %v, %v", sf.Lbays, sf.Lbay)
				return
			} else {
				sf.Lbays = []float64{sf.Lbay}
			}
		}
		s.Lspans = make([]float64, len(sf.Lbays))
		for i, bay := range sf.Lbays{
			s.Lspans[i] = bay * 1000.0
		}
		s.Nspans = len(s.Lspans)
		s.Ly = 1000.0
		s.Endc = 2
		switch len(s.Lspans){
			case 1:
			s.Endc = 1
			s.Typstr = "1w"
		}
		case 2:
		//two way slab load
	}
	// fmt.Println("kalling slab dz")
	err = SlbDesign(&s)
	if err != nil{
		fmt.Println(err)
		return
	}
	
	// fmt.Println("done slab dz")
	//s.Table()
	//fmt.Println("behold, dused and dl before->",sf.Dslb, sf.DL)
	switch s.Type{
		case 1, 2:
		sf.Dslb = s.Dused
		sf.Lbay = sf.Lbays[0]
		//fmt.Println("after",sf.Dslb, sf.Slbload)
		case 3:
		case 4:
	}
	sf.Report += s.Report
	_, err = DzSubFrm(sf)
	return
}

// //SubFlr stores a subframe and slice of slabs
// //doesn't look like it'll be used
// type SubFlr struct{
// 	Sf *SubFrm
// 	Slbz []*RccSlb
// 	Nbays int
// 	Slbtyp int
// }

//add - save beam env

//Printz printz horribly
// func (sf *SubFrm) Printz(){
// 	//REDO DIS
// 	fmt.Println("subframe")
// 	fmt.Println("grade of concrete, steel-",sf.Fcks, sf.Fys)
// 	fmt.Println("nspans-",sf.Nspans)
// 	fmt.Println("dl, ll-",sf.DL, sf.LL)
// 	fmt.Println("clvrs-",sf.Clvrs)
// 	fmt.Println("h-",sf.Hs)
// }

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

//ChainSlb designs a slab and uses the slab depth and loading
//to analyze and design a SubFrm
//use CHAIN more often; write more of such funcs
// func (sf *SubFrm) ChainSlb(s *RccSlb)(err error){
// 	err = SlbDesign(s)
// 	if err != nil{
// 		fmt.Println(err)
// 		return
// 	}
// 	s.Table(true)
// 	fmt.Println("behold, dused and dl before->",sf.Dslb, sf.DL)
// 	switch s.Type{
// 		case 1, 2:
// 		sf.Dslb = s.Dused
// 		sf.Slbload = s.Type
// 		fmt.Println("after",sf.Dslb, sf.Slbload)
// 		case 3:
// 		case 4:
// 	}
// 	_, err = DzSubFrm(sf)
// 	return
// }
