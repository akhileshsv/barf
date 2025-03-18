package barf

import (
	"fmt"
	"log"
	"encoding/json"
	flay"barf/flay"
)


//flayread reads json input
func flayread(choice, input int)(bytestr []byte, err error){
	var basefile string
	switch input{
		case 0:
		//read base json
		basefiles := []string{"crft","sqr"}
		basefile = fmt.Sprintf("flr%s_base.json",basefiles[choice])
		bytestr, err = readjsontxt(basefile)
		case 1:
		//get file
		bytestr, err = getjsonfile()
	}
	return
}

//flaymenu is the cli menu for facility layout funcs in /flay
func flaymenu(term string){
	running := true
	for running{
		choice := printmenu(icon_flay, flay_menus)
		switch choice{
			case 2:
			running = false
			break
			default:
			input := printmenu("choose input type", input_menus)
			bytestr, err := flayread(choice, input)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				continue
			}
			switch choice{
				case 0:
				//craft 
				var c flay.Crft
				err = json.Unmarshal(bytestr, &c)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					continue
				}
				c.Verbose = true
				err = c.Craft()
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					continue
				}		
				case 1:
				//squarify
				var f flay.Flr
				err = json.Unmarshal(bytestr, &f)
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					continue
				}
				f.Verbose = true
				f.Round = true
				f.Cgrid = true
				f.Sort = true
				err = f.FlrLay()
				if err != nil{
					log.Println(ColorRed, err, ColorReset)
					continue
				}
				f.FlrJson()
			}
		}
	}
}
