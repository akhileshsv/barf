package barf

import (
	"os"
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
)

//ReadBm reads a .json file and returns a Bm struct
func ReadBm(filename string) (Bm, error){
	var b Bm
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return b, err
	}
	err = json.Unmarshal([]byte(jsonfile), &b)
	if err !=nil{
		log.Println(err)
	}
	return b, err
}

//ReadCol reads a .json file and returns a Col struct
func ReadCol(filename string) (Col, error){
	var c Col
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return c, err
	}
	err = json.Unmarshal([]byte(jsonfile), &c)
	if err !=nil{
		log.Println(err)
	}
	return c, err
}

//CalcInp is the main entry func for steel design from flags/menu
func CalcInp(mtyp, cmdz, fname, term string, pipe bool) (err error){
	//entry func from flags/menu
	temp := os.Stdout
	if pipe{
		os.Stdout = nil
		
	}
	var jsonstr []byte
	switch mtyp{
		case "cbeam","cb","1db":
		//cbeam design
		case "beam","bm","rccbm":
		var b Bm
		b, err = ReadBm(fname)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		switch cmdz{
			case "design","dz","":
			//design
			//log.Println("beam design yo")
			err = BmDesign(&b)
			
			case "az","analyze":
			//analyze beam section
			log.Println("beam analysis yo")
			//err = BmAnalyze(&b)
		}
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return
		}
		if pipe{
			jsonstr, err = json.Marshal(&b)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				return
			}
		}
		case "col","column":
		var c Col
		c, err = ReadCol(fname)
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return				
		}
		switch cmdz{
			case "design","dz","":
			//design
			err = ColDesign(&c)
			case "az","analyze","calc":
			//analyze/check
			log.Println(ColorRed,"column analysis yo",ColorReset)
			//err = ColAnalyze(&c)
		}
		
		if pipe{
			jsonstr, err = json.Marshal(&c)
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				return
			}
		}
		case "truss","t2d","trs2d":
		case "portal","p2d","pfrm":
		case "frame2d","f2d","frm2d":
	}
	if pipe{
		os.Stdout = temp
		fmt.Print(string(jsonstr))
		
	}
	return
}
