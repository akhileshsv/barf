package barf

import (
	"fmt"
	"log"
	"encoding/json"
	kass"barf/kass"
)

//kjsondxs stores json indices of kassmenu items in basefiles
//"0 beam","1 2d truss","2 2d frame","3 3d truss","4 grid","5 3d frame","6 connections",

var kjsondxs = []int{0,8,13,20,22,23,25}

//kassread returns the json file either via menu text input (prints basefile for edits)
//or reads in .json from the input path provided
//basefiles := []string[]string{"0 bm1d","1 bmsd","2 bmtd","3 npbm","4 ssbm","5 bmsf","6 bmenv","7 bmep","8 t2d","9 t2sd","10 t2td","11 t2nl","12 t2gen","13f2d","14 f2rel","15 f2sd","16 f2td","17 f2np","18 f2ep","19 f2gen","20 t3d","21 t3sd","22 g3d","23 f3d","24 f3gen","25 blt","26 wld"}
func kassread(choice, input int)(bytestr []byte, err error){
	var basefile string
	basefiles := []string{"bm1d","bmsd","bmtd","npbm","ssbm","bmsf","bmenv","bmep","t2d","t2sd","t2td","t2nl","t2gen","f2d","f2rel","f2sd","f2td","f2np","f2ep","f2gen","t3d","t3sd","g3d","f3d","f3gen","blt","wld"}
	switch input{
		case 0:
		basefile = fmt.Sprintf("k%s_base.json",basefiles[choice])
		bytestr, err = readjsontxt(basefile)
		case 1:
		bytestr, err = getjsonfile()
	}
	return
}


//kassmenu is the cli menu for analysis and calc funcs in /kass
func kassmenu(term string){
	running := true
	for running{
		choice := printmenu(icon_kass, kass_menus)
		switch choice{
			case 7:
			running = false
			break
			case 0:
			kbmmenu(choice, term)
			case 1:
			kt2dmenu(choice,term)
			case 2:
			kf2dmenu(choice,term)
			case 3:
			kt3dmenu(choice,term)
			case 4:
			kgmenu(choice, term)
			case 5:
			kf3dmenu(choice,term)
			case 6:
			kconmenu(choice, term)
		}
	}
	return
}

//kbmmenu is the cli menu for beam analysis and calc funcs in /kass
func kbmmenu(choice int, term string){
	running := true
	cdx := kjsondxs[0]
	for running{
		choice := printmenu(icon_beam,[]string{"1d beam","beam (sup. displ.)","beam (temp. changes)","non uniform beam","s.s beam","beam sf/bm","beam envelopes","elastic-plastic beam","exit"})
		switch choice{
			case 8:
			running = false
			break
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := kassread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			//fmt.Println(bytestr)
			var mod *kass.Model
			err = json.Unmarshal(bytestr,&mod)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0,1,2:
				//calcbm1d
				err = kass.CalcMod(mod, "1db", term, false)
				case 3:
				//non uniform beam
				err = kass.CalcModNp(mod, "1db", term, false)
				case 4,5:
				//beam sf/bm
				err = kass.CalcModBmSf(mod)
				case 6:
				//beam envelopes
				//THIS IS STUPID CHANGE IT
				loadenvz, _ := kass.CalcBmEnv(mod) 
				for i := 1; i <= len(loadenvz); i++{
					bm := loadenvz[i]
					fmt.Printf(ColorGreen,"span %v\n",bm.Mem,ColorReset)
					//rezstring += "\tshear -kn \t moment-hog kn-m \t moment-sag kn-m\n"
					for j := 0; j < 21; j++ {
						fmt.Printf("sec %v\t%sshear %.3f\t%smoment-hog %.3f\tmoment-sag %.3f%s\n",j+1,ColorRed,bm.Venv[j],ColorCyan,bm.Mnenv[j],bm.Mpenv[j],ColorReset)
					}  
				}
				case 7:
				//ep beam
				kass.CalcEpFrm(mod)
			}
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
		}
	}
	return
}
	
//kt2dmenu is the cli menu for 2d truss analysis and calc funcs in /kass
func kt2dmenu(choice int, term string){
	running := true
	cdx := kjsondxs[1]
	for running{
		choice := printmenu(icon_truss,[]string{"2d truss","2d truss (sup. displ.)","2d truss (temp. changes)","non linear analysis","generate","exit"})
		switch choice{
			case 5:
			running = false
			break
			default:
			if choice == 3 || choice == 4{
				fmt.Println(ColorRed,icon_warning, ColorReset)
			}
			input := printmenu("choose input type", input_menus)
			bytestr, err := kassread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0,1,2,3:
				var mod kass.Model
				err = json.Unmarshal(bytestr,&mod)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					continue
				}
				switch choice{
					case 0,1,2:
					//call calc mod
					err = kass.CalcMod(&mod,"2dt",term, false)
					if err != nil{
						log.Println(ColorRed, err, ColorReset)
						continue
					}
					case 3:
					//call nl truss calc
					fmt.Println(ColorRed,icon_warning,ColorReset)
					kass.NlCalcTrs2d(&mod, 0.0)
				}
				case 4:
				fmt.Println(ColorRed,icon_warning,ColorReset)
				var trs kass.Trs2d
				err = json.Unmarshal(bytestr,&trs)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					continue
				}
				trs.Calc()
			}
		}
	}
	return
}
	
//kf2dmenu is the cli menu for analysis and calc funcs in /kass
func kf2dmenu(choice int, term string){
	running := true
	cdx := kjsondxs[2]
	for running{
		choice := printmenu(icon_frame2d,[]string{"2d frame","2d frame (mem. releases)","2d frame (sup. displ.)","2d frame (temp. changes)","non uniform frame","elastic-plastic frame","generate","exit"})
		switch choice{
			case 7:
			running = false
			break
			default:
			if choice == 6{
				fmt.Println(ColorRed,icon_warning, ColorReset)
			}
			input := printmenu("choose input type", input_menus)
			bytestr, err := kassread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0,1,2,3,4,5:
				var mod kass.Model
				err = json.Unmarshal(bytestr,&mod)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					continue
				}
				switch choice{
					case 0,1,2,3:
					err = kass.CalcMod(&mod,"2df",term, false)
					case 4:
					err = kass.CalcModNp(&mod, "2df", term, false)
					case 5:
					fmt.Println(ColorRed,icon_warning,ColorReset)
					kass.CalcEpFrm(&mod)
				}
				case 6:
				fmt.Println(ColorRed,icon_warning,ColorReset)
				
				var f kass.Frm2d
				err = json.Unmarshal(bytestr,&f)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					continue
				}
				
				err = f.Calc()
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					continue
				}
				for lp := range f.Loadcons{f.DrawLp(lp, term)}
			}
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
		}
	}
	return
}
	
//kt3dmenu is the cli menu for 3d truss analysis and calc funcs in /kass
func kt3dmenu(choice int, term string){
	running := true
	cdx := kjsondxs[3]
	for running{
		choice := printmenu(icon_truss3d,[]string{"3d truss","3d truss(sup. displ.)","exit"})
		switch choice{
			case 2:
			running = false
			break
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := kassread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var mod kass.Model
			err = json.Unmarshal(bytestr,&mod)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			err = kass.CalcMod(&mod, "3dt",term, false)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
		}
	}
	return
}

//kgmenu is the cli menu for (3d) grid analysis and calc funcs in /kass
func kgmenu(choice int, term string){
	running := true
	cdx := kjsondxs[4]
	for running{
		choice := printmenu(icon_grid,[]string{"3d grid analysis","exit"})
		switch choice{
			case 1:
			running = false
			break
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := kassread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var mod kass.Model
			err = json.Unmarshal(bytestr,&mod)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			err = kass.CalcMod(&mod, "3dg",term, false)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
		}
	}
	return
}
	
//kf3dmenu is the cli menu for 3d frame analysis and calc funcs in /kass
func kf3dmenu(choice int, term string){
	running := true
	cdx := kjsondxs[5]
	for running{
		choice := printmenu(icon_frame3d,[]string{"3d frame","generate","exit"})
		switch choice{
			case 2:
			running = false
			break
			case 1:
			fmt.Println(ColorRed,"3d frame is light years away m(_ _)m",ColorReset)
			continue
			case 0:
			input := printmenu("choose input type", input_menus)
			bytestr, err := kassread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var mod kass.Model
			err = json.Unmarshal(bytestr,&mod)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			err = kass.CalcMod(&mod, "3df",term, false)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
		}
	}
	return
}

//kconnmenu is the cli menu for connection analysis and calc funcs in /kass
func kconmenu(choice int, term string){
	running := true
	cdx := kjsondxs[6]
	for running{
		choice := printmenu(icon_conn,[]string{"bolt group analysis","weld group analysis","exit"})
		switch choice{
			case 2:
			running = false
			break
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := kassread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0:
				//bolt group analysis
				var b kass.Blt
				err = json.Unmarshal(bytestr, &b)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					return
				}
				err = kass.BoltSs(&b)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					return
				}
				fmt.Println(b.Report)
				case 1:
				//weld group analysis
				var w kass.Wld
				err = json.Unmarshal(bytestr, &w)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
				}
				err = kass.WeldSs(&w)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					return
				}
				fmt.Println(w.Report)
			}
		}
	}
}
	

/*
   
//kasscalc0 calls calc/analysis funcs in /kass from menu choices and json input
func kasscalc0(choice int, term string) (err error){
	var basefile, frmtyp string
	var bytestr []byte
	
	switch choice{
		case 0:
		basefile = "b1d_base.json"
		frmtyp = "1db"
		case 1:
		basefile = "t2d_base.json"
		frmtyp = "2dt"
		case 2:
		basefile = "f2d_base.json"
		frmtyp = "2df"
		case 3:
		basefile = "t3d_base.json"
		frmtyp = "3dt"
		case 4:
		basefile = "g3d_base.json"
		frmtyp = "3dg"
		case 5:
		basefile = "f3d_base.json"
		frmtyp = "3df"
		case 6:
		basefile = "npbm_base.json"
		frmtyp = "1db"
		case 7:
		basefile = "npfrm_base.json"
		frmtyp = "2df"
		case 8:
		basefile = "boltss_base.json"
		frmtyp = "blt"
		case 9:
		basefile = "weldss_base.json"
		frmtyp = "wld"
	}
	switch input{
		case 0:
		bytestr, err = readjsontxt(basefile)
		case 1:
		bytestr, err = getjsonfile()
	}
	if err != nil{
		log.Println(ColorRed,err,ColorReset)
		return
	}
	switch{
		case choice < 6:	
		//basic kassimali models 
		var mod kass.Model
		err = json.Unmarshal(bytestr,&mod)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		err = kass.CalcMod(&mod, frmtyp, term)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		case choice < 8:
		//non uniform member beam and frame
		var mod kass.Model
		err = json.Unmarshal(bytestr,&mod)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		err = kass.CalcModNp(&mod, frmtyp, term)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}		
		case choice < 10:
		//bolt and weld group analysis
		switch choice{
			case 8:
			//bolt group analysis
			var b kass.Blt
			err = json.Unmarshal(bytestr, &b)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				return
			}
			err = kass.BoltSs(&b)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				return
			}
			fmt.Println(b.Report)
			case 9:
			//weld group analysis
			var w kass.Wld
			err = json.Unmarshal(bytestr, &w)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
			}
			err = kass.WeldSs(&w)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				return
			}
			fmt.Println(w.Report)
		}
	}
	return
}
kass_menus = []string{
		"0 beam",
		"1 2d truss",
		"2 2d frame",
                "3 3d truss",
		"4 grid",
		"5 3d frame",
		"6 connections",
		"7 exit",
	}
*/
