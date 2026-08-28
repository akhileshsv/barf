package barf

import (
	"fmt"
	"log"
	"encoding/json"
	bash"barf/bash"
	"github.com/go-gota/gota/series"
)


//bjsondxs stores choice dxs for read funcs (start index)
var bjsondxs = []int{0,2,4,5}

//bashread returns the json file either via menu text input (prints basefile for edits)
//or reads in .json from the input path provided
func bashread(choice, input int) (bytestr []byte, err error){
	//0 - "bmis",1 "bmbs",2 "colis",3 "colbs",4 "cbs",5 "flr"
	var basefile string
	basefiles := []string{"bmis","bmbs","colis","colbs","cbs","flr"}
	basefile = fmt.Sprintf("stl%s_base.json",basefiles[choice])
	switch input{
		case 0:
		bytestr, err = readjsontxt(basefile)
		case 1:
		bytestr, err = getjsonfile()
		case 2:
		bytestr, err = readjsontxt(basefile)
		
	}
	return

}

//bbmmenu is the cli menu for beam design funcs in /bash
func bbmmenu(term string){
	running := true
	cdx := bjsondxs[0]
	for running{
		choice := printmenu(icon_beam,[]string{"is 800","bs 449","exit"})
		switch choice{
			case 2:
			running = false
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := bashread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var b bash.Bm
			b.Verbose = true
			err = json.Unmarshal(bytestr, &b)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0:
				b.Code = 1
				case 1:
				b.Code = 2
			}
			
			err = bash.BmDz(&b)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}			
		}
	}
	return
}


//bcolmenu is the cli menu for column design funcs in /bash
func bcolmenu(term string){
	running := true
	cdx := bjsondxs[1]
	for running{
		choice := printmenu(icon_col,[]string{"is 800","bs 449","exit"})
		switch choice{
			case 2:
			running = false
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := bashread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch input{
				case 2:
				bcolbch(bytestr)
				default:
				
				var c bash.Col
				err = json.Unmarshal(bytestr, &c)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					continue
				}
				
				c.Verbose = true
				if c.Code == 0{				
					switch choice{
						case 0:
						c.Code = 1
						case 1:
						c.Code = 2
					}
				}
				err = bash.ColDz(&c)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					continue
				}
				
			}	
		}
	}
	return
}


//bcbsmenu is the cli menu for base plate design funcs in bash
func bcbsmenu(term string){
	running := true
	cdx := bjsondxs[2]
	for running{
		choice := printmenu(icon_bp,[]string{"axially loaded base plate","exit"})
		switch choice{
			case 1:
			running = false
			case 0:
			input := printmenu("choose input type", input_menus)
			bytestr, err := bashread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			
			var c bash.Cbs
			err = json.Unmarshal(bytestr, &c)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			c.Verbose = true
			err = bash.SlbCbsDz(&c)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}			
		}
	}
}

//bflrmenu is the cli menu for floor design funcs in bash
func bflrmenu(term string){
	running := true
	cdx := bjsondxs[3]
	for running{
		choice := printmenu(icon_flr,[]string{"design steel floor grid","exit"})
		switch choice{
			case 1:
			running = false
			case 0:
			input := printmenu("choose input type", input_menus)
			bytestr, err := bashread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			
			var f bash.StlFlr
			err = json.Unmarshal(bytestr, &f)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			err = f.Calc()
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}			
		}
	}

}


//bashmenu is the cli menu for steel design funcs in bash
func bashmenu(term string){
	//var flrs []Rcflr0
	running := true
	//term := getterminal()
	for running{
		choice := printmenu(icon_bash,bash_menus)
		switch choice{
			case 4:
			//exit
			running = false
			case 0:
			//beam
			bbmmenu(term)
			case 1:
			//col
			bcolmenu(term)
			case 2:
			//column base plate
			bcbsmenu(term)
			case 3:
			//stlflr
			bflrmenu(term)
		}
	}
	return
}

func bcolbch (bytestr []byte){
	log.Println(ColorYellow,"reading steel column batch csv",ColorReset)
	df, err := readcsvfile()
	if err != nil{
		log.Println(ColorRed, err, ColorReset)
		return
	}
	var dtyp int
	var rez []string
	fmt.Println("df in",df)
	for i:=0; i < df.Nrow(); i++{
		dtyp, err = df.Elem(i,0).Int()
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return
		}
		switch dtyp{
			case 1:
			//pu, dx, dy
			if df.Ncol() < 5{
				log.Println(ColorRed, "not enough values (<5)", ColorReset)
				return	
			}
			var pu, dx, dy float64
			var cname string
			pu = df.Elem(i, 1).Float()
			dx = df.Elem(i, 2).Float()
			dy = df.Elem(i, 3).Float()
			cname = df.Elem(i,4).String()
			var c bash.Col
			err = json.Unmarshal(bytestr, &c)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				return
			}
			c.Pu = pu; c.Vx = dx; c.Vy = dy; c.Title = c.Title + cname
			err = bash.ColDz(&c)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				rez = append(rez, fmt.Sprintf("%v",err))
			} else {
				rez = append(rez, c.Ssecs[0].Sstr)
			}
		}
	}	
	fmt.Println(ColorGreen,"\n",rez,ColorReset)
	dfrez := df.Mutate(
		series.New(rez,series.String,"sec"),
	)
	fmt.Println(dfrez)
}


/*

//bashcalc calls steel design funcs from /bash based on menu choices and json input
func bashcalc(bytestr []byte, choice int, term string) (err error){
	switch choice{
		case 0:
		//beam
		var b bash.Bm
		err = json.Unmarshal(bytestr, &b)
		if err != nil{
			return
		}
		err = bash.BmDesign(&b)
		case 1:
		//col
		var c bash.Col
		err = json.Unmarshal(bytestr, &c)
		if err != nil{
			return
		}
		err = bash.ColDesign(&c)
		//case 3
		//case 4:
	}
	return
}

*/
