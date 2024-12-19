package barf

import (
	"fmt"
	"log"
	"encoding/json"
	kass"barf/kass"
	mosh"barf/mosh"
)

//mjsondxs stores choice dxs for read funcs (start index)
//slab, bm, col, ftng, cbm, sf, 2df, 3df 		
var mjsondxs = []int{0,5,7,10,13,15,19,21}

//moshread returns the json file either via menu text input (prints basefile for edits)
//or reads in .json from the input path provided
//basefiles := []string{"0 slb1","1 slb2","2 rslb","3 wslb","4 cslb1","5 bm","6 bmaz","7 col","8 colaz","9 colopt","10 ftngaz","11 ftngdz","12 ftngopt","13 cbm","14 cbmopt","15 sf","16 sfdz","17 fslb","18 sfopt","19 f2d","20 f2dopt"}
func moshread(choice, input int)(bytestr []byte, err error){
	var basefile string
	switch input{
		case 0:
		basefiles := []string{"slb1","slb2","rslb","wslb","cslb1","bm","bmaz","col","colaz","colopt","ftng","ftngaz","ftngopt","cbm","cbmopt","sf","sfdz","fslb","sfopt","f2ddz","f2dopt"}
		basefile = fmt.Sprintf("rc%s_base.json",basefiles[choice])
		bytestr, err = readjsontxt(basefile)
		case 1:
		bytestr, err = getjsonfile()
	}
	return
}

//rcslbmenu is the cli menu for rcc slab funcs
func rcslbmenu(term string){
	running := true
	cdx := mjsondxs[0]
	for running{
		choice := printmenu(icon_slab,[]string{"1-way","2-way","ribbed","waffle","1-way continuous","exit"})
		switch choice{
			case 5:
			running = false 
			break
			default:
			if choice == 2 || choice == 3{
				fmt.Println(ColorRed, icon_warning, ColorReset)
			}
			input := printmenu("choose input type", input_menus)
			bytestr, err := moshread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var s mosh.RccSlb
			err = json.Unmarshal(bytestr,&s)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			err = mosh.SlbDesign(&s)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
		}
	}
}

//rcbmmenu is the cli menu for rcc beam funcs
func rcbmmenu(term string){
	running := true
	cdx := mjsondxs[1]
	for running{
		choice := printmenu(icon_beam,[]string{"design","section analysis","exit"})
		switch choice{
			case 2:
			running = false
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := moshread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var b mosh.RccBm
			err = json.Unmarshal(bytestr,&b)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0:
				//beam design
				err = mosh.BmDesign(&b)
				case 1:
				//beam analysis
				err = mosh.BmAnalyze(&b)				
			}
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			b.Table(true)
			pltstr := mosh.PlotBmGeom(&b, b.Term) 
			fmt.Println(pltstr)
		}
	}
}

//rccolmenu is the cli menu for rcc col funcs
func rccolmenu(term string){
	running := true
	cdx := mjsondxs[2]
	for running{
		choice := printmenu(icon_col,[]string{"design","section analysis","optimize","exit"})
		switch choice{
			case 3:
			running = false
			break
			case 2:
			//col opt
			fmt.Println(ColorYellow,"column optimization is a month or two away \\('')/",ColorReset)
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := moshread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var c mosh.RccCol
			err = json.Unmarshal(bytestr,&c)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0:
				//col design				
				err = mosh.ColDesign(&c)
				if err != nil{
					log.Println(ColorRed,err,ColorReset)
					continue
				}
				c.Table(false)
				fmt.Println(c.Report)
				switch c.Term{
					case "dumb","mono":
					fmt.Println(c.Txtplot[0])	
					case "svg":
					fmt.Printf("%s svg plot file saved at %s %s\n",ColorCyan,c.Txtplot[0],ColorReset)
				}
				case 1:
				//col analysis
				err = mosh.ColAnalyze(&c)				
				if err != nil{
					continue
				}
			}
		}
	}
	return
}

//rcftngmenu is the cli menu for rcc ftng funcs
func rcftngmenu(term string){
	running := true
	cdx := mjsondxs[3]
	for running{
		choice := printmenu(icon_ftng,[]string{"pad/sloped footing design","pad footing analysis","optimize","exit"})
		switch choice{
			case 3:
			running = false
			break
			case 2:
			//ftng optimization
			fmt.Println(ColorYellow,"footing optimization is many months away m(_ _)m",ColorReset)
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := moshread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var f mosh.RccFtng
			err = json.Unmarshal(bytestr,&f)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0:
				//ftng design
				err = mosh.FtngDzRojas(&f)
				if err != nil{
					log.Println(ColorRed,err,ColorReset)
					continue
				}
				case 1:
				//ftng analysis
				err = mosh.FtngPadAz(&f)				
				if err != nil{
					log.Println(ColorRed,err,ColorReset)
					continue
				}
			}
		}
	}
	return
}

//rccbmmenu is the cli menu for rcc cbeam funcs
func rccbmmenu(term string){
	running := true
	cdx := mjsondxs[4]
	for running{
		choice := printmenu(icon_beam,[]string{"design","optimize","exit"})
		switch choice{
			case 2:
			running = false
			break
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := moshread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var cb mosh.CBm
			err = json.Unmarshal(bytestr,&cb)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			err = mosh.CalcCBm(cb)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
		}
	}
	return	
}

//rcsfmenu is the cli menu for rcc sub frame and flat slab funcs
func rcsfmenu(term string){
	running := true
	cdx := mjsondxs[5]
	for running{
		choice := printmenu(icon_subframe,[]string{"analyze","design","optimize","flat slab design","exit"})
		switch choice{
			case 4:
			running = false
			break
			case 2:
			fmt.Println(ColorYellow,"sub frm optimization will take eons ('//_ _)'//",ColorReset)
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := moshread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var sf mosh.SubFrm
			err = json.Unmarshal(bytestr,&sf)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0:
				err = mosh.CalcSubFrm(&sf)
				case 1:
				_, err = mosh.DzSubFrm(&sf)
				case 3:
				err = mosh.FltSlbDz(&sf)
			}
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
		}
	}
	return
}

//rcf2dmenu is the cli menu for rcc 2d frame funcs
func rcf2dmenu(term string){
	running := true
	cdx := mjsondxs[6]
	for running{
		choice := printmenu(icon_frame2d,[]string{"design","optimize","exit"})
		switch choice{
			case 2:
			running = false
			break
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := moshread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var f kass.Frm2d
			err = json.Unmarshal(bytestr,&f)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0:
				//log.Println("dz frame")
				_ , _, err = mosh.Frm2dDz(&f)
				case 1:
				_, err = mosh.Frm2dOpt(f)
				//log.Println("opt frame")
				//log.Println("report->",frmrez.Report)
				
			}
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
		}
	}
	return
}

//rcf3dmenu is the cli menu for rcc 3d frame funcs
func rcf3dmenu(term string){}

//moshmenu is the cli menu func for rcc design funcs in /mosh
func moshmenu(term string){
	//var flrs []Rcflr0
	running := true
	//term := getterminal()
	for running{  
		choice := printmenu(icon_mosh,mosh_menus)
		switch choice{
			case 0:
			//slab funcs
			rcslbmenu(term)
			case 1:
			//beam funcs
			rcbmmenu(term)
			case 2:
			//col funcs
			rccolmenu(term)
			case 3:
			//ftng funcs
			rcftngmenu(term)
			case 4:
			//cbeam	
			rccbmmenu(term)
			case 5:
			//sub frame
			rcsfmenu(term)
			case 6:
			//frame 2d
			rcf2dmenu(term)
			case 7:
			//frame 3d
			fmt.Println(ColorRed,"3d frame is light years away m(_ _)m",ColorReset)
			case 8:
			//EXIT
			running = false
			break
		}
	}
	return
}

/*
savez := savemenu()
	if savez{
		fname := getfilename()
		filename, e := sf.Dump(fname)
		if e != nil{err = e; return}
		fmt.Println("saved json at",filename)
	}

*/

/*
//moshcalc calls calc/analysis funcs in /mosh from menu choices and json input
func moshcalc(bytestr []byte, choice int, term string) (err error){
	//menu loop for rcc/mosh
	switch choice{
		case 0,1,2,3,4:
		//slab funcs
		var s mosh.RccSlb
		err = json.Unmarshal(bytestr, &s)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		err = mosh.SlbDesign(&s)
		
		case 5,6:
		var b mosh.RccBm
		err = json.Unmarshal(bytestr, &b)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
		}
		switch choice{
			case 5:
			
			case 6:
		}
		err = mosh.BmDesign(&b)

		case 7,8,9:
		var c mosh.RccCol
		err = json.Unmarshal(bytestr, &c)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
		}
		switch choice{
			case 7:
			err = mosh.ColDesign(&c)
			case 8:
			
			case 9:
		}
		case 3:
		var f mosh.RccFtng
		err = json.Unmarshal(bytestr, &f)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
		}
		err = mosh.FtngDzRojas(&f)

		case 4:
		case 5:
		case 6:
		var cb mosh.CBm
		err = json.Unmarshal(bytestr, &cb)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		err = mosh.CalcCBm(cb)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		case 7:
		var sf mosh.SubFrm
		err = json.Unmarshal(bytestr, &sf)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		err = mosh.CalcSubFrm(&sf)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		err = mosh.DzSubFrm(&sf)
		case 8:
		//2d frame
		var fr kass.Frm2d
		err = json.Unmarshal(bytestr, &fr)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		err, _, _ = mosh.Frm2dDz(&fr)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		case 9:
		//3d frame
		log.Println("3d frame will take a few years (at current speed)")
		case 10:
		//flat slab		
		log.Println("flat slab will take a few years (at current speed)")
	}
	if err != nil{
		log.Println(ColorRed,err,ColorReset)
		return
	}
	return
}
*/
