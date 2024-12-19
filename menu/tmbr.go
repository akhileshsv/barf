package barf

import (
	"fmt"
	"log"
	"encoding/json"
	kass"barf/kass"
	tmbr"barf/tmbr"
)

//tjsondxs stores timber choice dxs for read funcs(start index)
var tjsondxs = []int{0,1,2}

//tmbrread returns the json file either via menu text input (prints basefile for edits)
//or reads in .json from the input path provided
func tmbrread(choice, input int) (bytestr []byte, err error){
	var basefile string
	basefiles := []string{"bm","col","trs2d"}
	switch input{
		case 0:
		basefile = fmt.Sprintf("tm%s_base.json",basefiles[choice])
		bytestr, err = readjsontxt(basefile)
		case 1:
		bytestr, err = getjsonfile()
	}
	return

}

//tmbmmenu is the cli menu for timber beam funcs
func tmbmmenu(term string){
	running := true
	cdx := tjsondxs[0]
	for running{
		choice := printmenu(icon_beam,[]string{"design","exit"})
		switch choice{
			case 1:
			running = false
			break
			case 0:
			input := printmenu("choose input type", input_menus)
			bytestr, err := tmbrread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var b tmbr.WdBm
			err = json.Unmarshal(bytestr, &b)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			err = tmbr.BmDesign(&b)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}			
		}
	}
	return
}

//tmcolmenu is the cli menu for timber column funcs
func tmcolmenu(term string){
	running := true
	cdx := tjsondxs[1]
	for running{
		choice := printmenu(icon_col,[]string{"design","exit"})
		switch choice{
			case 1:
			running = false
			break
			case 0:
			input := printmenu("choose input type", input_menus)
			bytestr, err := tmbrread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var c tmbr.WdCol
			err = json.Unmarshal(bytestr, &c)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}						
			err = tmbr.ColDz(&c)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
		}
	}
	return
}

//tmt2dmenu is the cli menu for timber truss funcs
func tmt2dmenu(term string){
	running := true
	cdx := tjsondxs[2]
	for running{
		choice := printmenu(icon_truss,[]string{"design","exit"})
		switch choice{
			case 1:
			running = false
			break
			case 0:
			input := printmenu("choose input type", input_menus)
			bytestr, err := tmbrread(choice + cdx, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			var t kass.Trs2d
			err = json.Unmarshal(bytestr, &t)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}						
			err = tmbr.TrussDz(&t)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
		}
	}
	return
}

//tmbrmenu is the cli menu func for timber design funcs in /tmbr
func tmbrmenu(term string){
	running := true
	for running{  
		choice := printmenu(icon_tmbr,tmbr_menus)
		switch choice{
			case 0:
			tmbmmenu(term)
			case 1:
			tmcolmenu(term)
			case 2:
			tmt2dmenu(term)
			case 3:
			//exit
			running = false
			break
		}
	}
	return
}

/*

//tmbrcalc calls calc/analysis funcs in /tmbr from menu choices and json input
func tmbrcalc(bytestr []byte, choice int, term string) (err error){
	switch choice{
		case 0:
		//beam
		var b tmbr.WdBm
		err = json.Unmarshal(bytestr, &b)
		if err != nil{
			return
		}
		err = tmbr.BmDesign(&b)
		case 1:
		var c tmbr.WdCol
		err = json.Unmarshal(bytestr, &c)
		if err != nil{
			return
		}
		err = tmbr.ColDz(&c)
		case 2:
		fmt.Println(ColorRed,icon_warning,ColorReset)
		var t kass.Trs2d
		err = json.Unmarshal(bytestr, &t)
		if err != nil{
			return
		}
		err = tmbr.TrussDz(&t)
		//case 3
		//case 4:
	}
	return
}

*/
