package barf

import (
	"fmt"
	"log"
	"encoding/json"
	"runtime"
	"path/filepath"
	"io/ioutil"
	kass"barf/kass"
)


//ReadRwall reads a .json file into an Rwall struct
func ReadRwall(filename string) (Rwall, error){
	var r Rwall
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return r, err
	}
	err = json.Unmarshal([]byte(jsonfile), &r)
	if err !=nil{
		log.Println(err)
	}
	return r, err
}

//ReadFtng reads a .json file into an RccFtng struct
func ReadFtng(filename string) (RccFtng, error){
	var f RccFtng
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return f, err
	}
	err = json.Unmarshal([]byte(jsonfile), &f)
	if err !=nil{
		log.Println(err)
	}
	return f, err
}

//ReadCBm reads a .json file into a CBm struct
func ReadCBm(filename string) (CBm, error){
	var cb CBm
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return cb, err
	}
	err = json.Unmarshal([]byte(jsonfile), &cb)
	if err !=nil{
		log.Println(err)
	}
	return cb, err
}

//ReadBm reads a .json file into an RccBm struct
func ReadBm(filename string) (RccBm, error){
	var b RccBm
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

//ReadCol reads a .json file into an RccCol struct
func ReadCol(filename string) (RccCol, error){
	var c RccCol
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

//ReadSlb reads a .json file into an RccSlb struct
func ReadSlb(filename string) (RccSlb, error){
	var s RccSlb
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return s, err
	}
	err = json.Unmarshal([]byte(jsonfile), &s)
	if err !=nil{
		log.Println(err)
	}
	return s, err
}

//ReadSubFrm reads a .json file into a SubFrm struct
func ReadSubFrm(filename string) (SubFrm, error){
	var sf SubFrm
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return sf, err
	}
	err = json.Unmarshal([]byte(jsonfile), &sf)
	if err !=nil{
		log.Println(err)
	}
	return sf, err
}

//CalcInp is the entry func from flags/menu for rcc design routines
func CalcInp(mtyp, cmdz, fname, chnf, term string, tweek, pipe bool) (err error){
	//entry func for mosh from flags/menu
	switch mtyp{
		case "slab","slb":
		//nope
		var s RccSlb
		s, err = ReadSlb(fname)
		if err != nil{			
			log.Println(ColorRed, err, ColorReset)
			return	
		}		
		if term != ""{s.Term = term}
		err = SlbDesign(&s)
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return				
		}
		if pipe{
			fname, e := s.JsonOut()
			if e != nil{
				err = e
				return
			}
			fmt.Printf("output file saved at -***->%s",fname)
		}
		case "cbeam","cb","1db","cbm":
		//nope
		var cb CBm
		cb, err = ReadCBm(fname)
		if err != nil{			
			log.Println(ColorRed, err, ColorReset)
			return	
		}
		err = CalcCBm(cb)
		if err != nil{			
			log.Println(ColorRed, err, ColorReset)
			return	
		}
		
		if pipe{
			fname, e := cb.JsonOut()
			if e != nil{
				err = e
				return
			}
			fmt.Printf("output file saved at -***->%s",fname)
		}
		case "beam","bm","rccbm":
		//nope
		var b RccBm
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
			case "az","analyze":
			//analyze beam section
			log.Println("beam analysis yo")
			err = BmAnalyze(&b)
		}
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return
		}
		
		if pipe{
			fname, e := b.JsonOut()
			if e != nil{
				err = e
				return
			}
			fmt.Printf("output file saved at -***->%s",fname)
		}
		case "col","column":
		var c RccCol
		c, err = ReadCol(fname)
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return				
		}
		switch cmdz{
			case "design","dz","":
			//design
			err = ColDesign(&c)
			//log.Println(c.Ast, c.Asc)
			if err != nil{
				log.Println(ColorRed,err,ColorReset)
				return
			}
			c.Table(false)
			fmt.Println(c.Report)
			if c.Term != ""{
				fmt.Println(c.Txtplot[0])	
			}
			case "az","analyze","calc":
			//analyze/NM curve
			//log.Println(ColorRed,"column analysis yo",ColorReset)
			err = ColAnalyze(&c)
		}
		
		if pipe{
			fname, e := c.JsonOut()
			if e != nil{
				err = e
				return
			}
			fmt.Printf("output file saved at -***->%s",fname)
		}
		case "footing","ftng":
		//footing design
		var f RccFtng
		f, err = ReadFtng(fname)
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return
		}
		switch cmdz{
			case "design","dz","":
			err = FtngDzRojas(&f)
			case "az","analyze","calc":
			err = FtngPadAz(&f)
		}
		
		if pipe{
			fname, e := f.JsonOut()
			if e != nil{
				err = e
				return
			}
			fmt.Printf("output file saved at -***->%s",fname)
		}
		case "subframe","sf","subfrm":
		//subframe design
		//nope
		var sf SubFrm
		sf, err = ReadSubFrm(fname)
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return
		}
		//if sf.Term == ""{sf.Term = "dumb"}
		if term != ""{sf.Term = term}
		if tweek{
			sf.Tweak = true
			sf.Tweakb = true
			sf.Tweakc = true
		}
		switch cmdz{
			case "fltslb","slb","flat slab","fs","fslb":
			err = FltSlbDz(&sf)
			case "design","dz":
			_, err = DzSubFrm(&sf)
			case "az","analyze","calc","":
			err = CalcSubFrm(&sf)
			case "chn","c","chain":
			//var s *RccSlb
			// var s RccSlb
			// s, err = ReadSlb(chnf)
			// if err != nil{			
			// 	log.Println(ColorRed, err, ColorReset)
			// 	return	
			// }
			// sf.ChainSlb(&s)
			sf.ChainSlab()
		}
		
		if pipe{
			fname, e := sf.JsonOut()
			if e != nil{
				err = e
				return
			}
			fmt.Printf("output file saved at -***->%s",fname)
		}
		case "frame2d","f2d","frm2d":
		//frame2d design
		//nope
		var f kass.Frm2d
		f, err = kass.ReadFrm2d(fname)
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return
		}
		switch cmdz{
			case "calc","dz","":
			Frm2dDz(&f)
			if f.Term != ""{
				for lp := range f.Loadcons{
					//fmt.Println(lp)
					f.DrawLp(lp, f.Term)
				}
			}
			case "opt":
			Frm2dOpt(f)
		}
		
		if pipe{
			fname, e := f.JsonOut()
			if e != nil{
				err = e
				return
			}
			fmt.Printf("output file saved at -***->%s",fname)
		}
		case "frame3d","f3d","frm3d":
		//frame3d design
		//nope
		/*
		var f kass.Frm3d
		f, err = kass.ReadFrm3d(fname)
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return
		}
		switch cmdz{
			case "calc","dz","":
			Frm2dDz(&f)
			case "opt":
			Frm2dOpt(&f)
		}
		*/
	}
	if err != nil{
		log.Println(ColorRed, err, ColorReset)
		return				
	}
	return
}


//JsonOut saves a slab struct to a json file in data/out
func (s *RccSlb) JsonOut() (fname string, err error){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	file, err := json.MarshalIndent(s, "", " ")
	if err != nil{
		return
	}
	fname = filepath.Join(basepath, "../data/out",s.Title+"_out.json")
	err = ioutil.WriteFile(fname, file, 0644)
	return
}


//JsonOut saves a cbm struct to a json file in data/out
func (cb *CBm) JsonOut() (fname string, err error){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	file, err := json.MarshalIndent(cb, "", " ")
	if err != nil{
		return
	}
	fname = filepath.Join(basepath, "../data/out",cb.Title+"_out.json")
	err = ioutil.WriteFile(fname, file, 0644)
	return
}


//JsonOut saves a bm struct to a json file in data/out
func (b *RccBm) JsonOut() (fname string, err error){
	_, base, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(base)
	file, err := json.MarshalIndent(b, "", " ")
	if err != nil{
		return
	}
	fname = filepath.Join(basepath, "../data/out",b.Title+"_out.json")
	err = ioutil.WriteFile(fname, file, 0644)
	return
}

//JsonOut saves a col struct to a json file in data/out
func (c *RccCol) JsonOut() (fname string, err error){
	_, base, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(base)
	file, err := json.MarshalIndent(c, "", " ")
	if err != nil{
		return
	}
	fname = filepath.Join(basepath, "../data/out",c.Title+"_out.json")
	err = ioutil.WriteFile(fname, file, 0644)
	return
}


//JsonOut saves a ftng struct to a json file in data/out
func (f *RccFtng) JsonOut() (fname string, err error){
	_, base, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(base)
	file, err := json.MarshalIndent(f, "", " ")
	if err != nil{
		return
	}
	fname = filepath.Join(basepath, "../data/out",f.Title+"_out.json")
	err = ioutil.WriteFile(fname, file, 0644)
	return
}


//JsonOut saves a subframe struct to a json file in data/out
func (sf *SubFrm) JsonOut() (fname string, err error){
	_, base, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(base)
	file, err := json.MarshalIndent(sf, "", " ")
	if err != nil{
		return
	}
	fname = filepath.Join(basepath, "../data/out",sf.Title+"_out.json")
	err = ioutil.WriteFile(fname, file, 0644)
	return
}

/*
func ReadFrm3d(filename string) (Frm3d, error) {
	var f Frm3d
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return f, err
	}
	err = json.Unmarshal([]byte(jsonfile), &f)
	return f, err
}

func ReadFrm2d(filename string) (Frm2d, error){
	var f Frm2d
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return f, err
	}
	err = json.Unmarshal([]byte(jsonfile), &f)
	return f, err

   }
   
*/
