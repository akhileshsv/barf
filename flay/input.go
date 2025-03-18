package barf

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

//ReadCrft reads a craft struct from a json file
func ReadCrft(fname string)(c Crft, err error){
	jsonfile, e := ioutil.ReadFile(fname)
	if e != nil {
		err = e
		return 
	}
	err = json.Unmarshal([]byte(jsonfile), &c)
	return
}

//CalcInp is the entry func from main/flags
func CalcInp(lay, cmdz, fname, term string, pipe bool) (err error){
	switch lay{
		case "craft", "crft":
		var c Crft
		c, err = ReadCrft(fname)
		if err != nil{return}
		err = c.Craft()
		if err != nil{return}
		fmt.Println(c.Report)
		case "sqr","squarify":
		fmt.Println("fuck u sir and have a shitty day")
	}
	return
}
