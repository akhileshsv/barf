package barf

import (
	"os"
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
	kass"barf/kass"
)

//ReadBm reads a .json file into a WdBm struct
func ReadBm(filename string) (WdBm, error){
	var b WdBm
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


//ReadCol reads a .json file into a WdCol struct
func ReadCol(filename string) (WdCol, error){
	var c WdCol
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

//CalcInp is the entry func for timber calcs from flags/menu
func CalcInp(mtyp, cmdz, fname, term string, pipe bool) (err error){
	//entry func from flags/menu
	temp := os.Stdout
	if pipe{
		os.Stdout = nil
		
	}
	switch mtyp{
		case "cbeam","cb","1db":
		//cbeam design - ARRE BHAIYYA
		case "beam","bm","rccbm":
		var b WdBm
		b, err = ReadBm(fname)
		if err != nil{
			log.Println(ColorRed,err,ColorReset)
			return
		}
		switch cmdz{
			case "design","dz","":
			//design
			log.Println("beam design yo")
			err = BmDesign(&b)
			case "az","analyze","chk":
			//check/analyze beam section for span/etc
		}
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return
		}
		if pipe{
			os.Stdout = temp
			jsonstr, e := json.Marshal(&b)
			if err != nil{
				log.Println(ColorRed, e, ColorReset)
				return
			}
			fmt.Print(string(jsonstr))
			
		}

		case "col","column":
		var c WdCol
		c, err = ReadCol(fname)
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return				
		}
		switch cmdz{
			case "design","dz","":
			//design
			err = ColDz(&c)
			
			if err != nil{
				log.Println(ColorRed, err, ColorReset)
				return				
			}
			case "az","analyze","calc":
			//analyze/check
			log.Println(ColorRed,"column analysis yo",ColorReset)
			//err = ColAnalyze(&c)
		}
		if pipe{
			os.Stdout = temp
			jsonstr, e := json.Marshal(&c)
			if err != nil{
				log.Println(ColorRed, e, ColorReset)
				return
			}
			fmt.Print(string(jsonstr))
			
		}

		case "portal","pfrm","p2d":
		//2d portal frame - do gables here too? lmao
		case "truss2d","t2d","trs2d":
		//truss design
		var t kass.Trs2d
		t, err = kass.ReadTrs2d(fname)
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return				
		}
		err = TrussDz(&t)
		if err != nil{
			log.Println(err)
			return
		}
		if pipe{
			os.Stdout = temp
			jsonstr, e := json.Marshal(&t)
			if err != nil{
				log.Println(ColorRed, e, ColorReset)
				return
			}
			fmt.Print(string(jsonstr))
			
		}
		case "frame2d","f2d","frm2d":
		//frame2d design
	}
	if err != nil{
		log.Println(ColorRed, err, ColorReset)
		return				
	}
	return
}
