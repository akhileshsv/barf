package barf

import (
	"fmt"
	"log"
	"encoding/json"
	bash"barf/bash"
)


//bjsondxs stores choice dxs for read funcs (start index)
var bjsondxs = []int{0,2,5}

//bashread returns the json file either via menu text input (prints basefile for edits)
//or reads in .json from the input path provided
func bashread(choice, input int) (bytestr []byte, err error){
	var basefile string
	basefiles := []string{"bm","bmfrm","col","colfrm","colchk","t2d","t2dopt"}
	basefile = fmt.Sprintf("stl%s_base.json",basefiles[choice])
	switch input{
		case 0:
		bytestr, err = readjsontxt(basefile)
		case 1:
		bytestr, err = getjsonfile()
	}
	return

}

//bbmmenu is the cli menu for beam design funcs in /bash
func bbmmenu(term string){
	running := true
	cdx := bjsondxs[0]
	for running{
		choice := printmenu(icon_beam,[]string{"design-ss beam","design-end moments specified","exit"})
		switch choice{
			case 2:
			running = false
			break
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := bashread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var b bash.Bm
			err = json.Unmarshal(bytestr, &b)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			
			b.Verbose = true
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
		choice := printmenu(icon_col,[]string{"design-framing beams","design-end moments specified","check section","exit"})
		switch choice{
			case 3:
			running = false
			break
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := bashread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var c bash.Col
			err = json.Unmarshal(bytestr, &c)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0,1:
				err = bash.ColDesign(&c)
				case 2:
				c.Spam = true; c.Verbose = true
				val, ok := bash.ColCBs(&c)
				fmt.Println(val, ok)
			}
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}			
		}
	}
	return
}

//bt2dmenu is the cli menu for truss design funcs in bash
func bt2dmenu(term string){
	fmt.Println(ColorYellow,"steel 2d truss design might take forever (/'-')/",ColorReset)
	return
}

//bashmenu is the cli menu for steel design funcs in bash
func bashmenu(term string){
	//var flrs []Rcflr0
	running := true
	//term := getterminal()
	for running{  
		choice := printmenu(icon_bash,bash_menus)
		switch choice{
			case 3:
			//exit
			running = false
			break
			case 0:
			//beam
			bbmmenu(term)
			case 1:
			//col
			bcolmenu(term)
			case 2:
			//trs 2d
			bt2dmenu(term)
		}
	}
	return
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
