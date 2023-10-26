package barf

import (
	"fmt"
	//"log"
	//"encoding/json"
	//flay"barf/flay"
)

//flaymenu is the cli menu for facility layout funcs in /flay
func flaymenu(term string){
	running := true
	for running{
		choice := printmenu(icon_flay, flay_menus)
		switch choice{
			case 2:
			running = false
			break
			case 0:
			fmt.Println("you has choosen->",flay_menus[choice])
			case 1:
			fmt.Println("you has choosen->",flay_menus[choice])
		}
	}
}
