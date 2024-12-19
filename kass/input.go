package barf

//input funcs for reading json and stuff
//I SHALL NOT DOCUMENT THIS holy hell

import (
	"os"
	"path/filepath"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"encoding/json"
	"runtime"
	"io/ioutil"
)

//ReadFrm2d reads a frm2d from a filename
func ReadFrm2d(filename string) (Frm2d, error){
	var f Frm2d
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return f, err
	}
	err = json.Unmarshal([]byte(jsonfile), &f)
	return f, err

}


//ReadTrs2d reads a frm2d from a filename
func ReadTrs2d(filename string) (Trs2d, error){
	var t Trs2d
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal([]byte(jsonfile), &t)
	return t, err

}


//ReadPortal reads a Portal frame from a filename
func ReadPortal(filename string) (Portal, error){
	var pf Portal
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return pf, err
	}
	err = json.Unmarshal([]byte(jsonfile), &pf)
	return pf, err

}

//ln2vec is older than time and common sense.
//haha. HAHAHA. parsing .txt for input lmao
func ln2vec(line string, flt bool, istring bool) interface{} {
	rwstr := strings.Split(line, ";")

	if istring {
		return rwstr
	}

	if flt {
		arr := [][]float64{}
		rwarr := []float64{}
		for _, row := range rwstr {
			//fmt.Println(row)
			//for _, yax := range strings.Split(row, ",") {fmt.Println("xxx000xxx");fmt.Println(yax);}
			//if len(rwarr) == clz {arr = append(arr, rwarr); rwarr = nil}
			for _, ch := range strings.Split(row, ",") {
				//fmt.Println(ch)
				val, err := strconv.ParseFloat(string(ch), 64)
				if err != nil {
					log.Println(err)
					continue
				} else {
					rwarr = append(rwarr, val)
				}
				//fmt.Println(val)
			}

			if rwarr != nil {
				arr = append(arr, rwarr)
				rwarr = nil
			}
		}
		return arr

	} else {
		arr := [][]int{}
		rwarr := []int{}
		for _, row := range rwstr {
			for _, ch := range strings.Split(row, ",") {
				val, _ := strconv.Atoi(string(ch))
				rwarr = append(rwarr, val)
			}
			arr = append(arr, rwarr)
			rwarr = nil
		}
		return arr
	}

}

//JsonInp reads a json filename and returns frmtyp, model and an error
func JsonInp(filename string) (string, *Model, error) {
	var mod Model
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return "", &mod, err
	}
	err = json.Unmarshal([]byte(jsonfile), &mod)
	if err != nil {
		log.Println(err)
		return "", &mod, err
	}
	//CHANGE UNITS AND FRMTYP TO NAMED FIELDS dammit
	//here is seen the price of stupidity
	var frmstr, unitz string
	if len(mod.Cmdz) > 0{
		frmstr = mod.Cmdz[0]
	}
	if len(mod.Cmdz) > 1{
		unitz = mod.Cmdz[1]
	}
	frmstr = strings.ToLower(frmstr)
	frmstr = strings.TrimSpace(frmstr)
	
	if mod.Frmstr == ""{
		mod.Frmstr = frmstr
	} else {
		frmstr = mod.Frmstr
	}
	if mod.Units == ""{
		mod.Units = unitz
	}
	switch frmstr {
	case "2dt","3dt","1db","2df","3df","3dg":
		return frmstr, &mod, nil	
	default:
		log.Println("invalid frame type")
		return "", &mod, errors.New("invalid frame type")
	}
}

//TxtInp is a prehistoric man making a crude club to hit himself over the head
func TxtInp(f string) (string, []interface{}, *Model) {
	var frm Model
	var frmtyp string
	model := make([]interface{}, 9)
	for lno, line := range strings.Split(f, "\n") {
		line = strings.ReplaceAll(line, " ", "")
		line = strings.ReplaceAll(line, "\r", "")
		var cnt interface{}
		switch lno {
		case 1:
			cnt = ln2vec(line, false, true)
			cmdz := cnt.([]string)
			model[7] = cmdz
			frm.Cmdz = cmdz
			frmtyp = cmdz[0]
		case 3:
			cnt = ln2vec(line, true, false)
			cords := cnt.([][]float64)
			model[0] = cords
			frm.Coords = cords
		case 5:
			cnt = ln2vec(line, false, false)
			msup := cnt.([][]int)
			model[1] = msup
			frm.Supports = msup
		case 7:
			cnt = ln2vec(line, true, false)
			em := cnt.([][]float64)
			model[2] = em
			frm.Em = em
		case 9:
			cnt = ln2vec(line, true, false)
			cp := cnt.([][]float64)
			model[3] = cp
			frm.Cp = cp
		case 11:
			cnt = ln2vec(line, false, false)
			mprp := cnt.([][]int)
			model[4] = mprp
			frm.Mprp = mprp
		case 13:
			cnt = ln2vec(line, true, false)
			jp := cnt.([][]float64)
			model[5] = jp
			frm.Jloads = jp
		case 15:
			cnt = ln2vec(line, true, false)
			pm := cnt.([][]float64)
			model[6] = pm
			frm.Msloads = pm
		case 17:
			cnt = ln2vec(line, true, false)
			wng := cnt.([][]float64)
			model[8] = wng
			frm.Wng = wng
		}
	}
	return frmtyp, model, &frm
}

//DesignInp is a dream from a different time
func DesignInp(f string) (data []interface{}) {
	data = make([]interface{}, 4)
	file := []string{}
	for lno, line := range strings.Split(f, "\n") {
		line = strings.ReplaceAll(line, " ", "")
		line = strings.ReplaceAll(line, "\r", "")
		file = append(file, line)
		var cnt interface{}
		switch lno {
		case 1:
			cnt = ln2vec(line, false, true)
			cmdz := cnt.([]string)
			data[0] = cmdz
		case 3:
			cnt = ln2vec(line, true, false)
			dims := cnt.([][]float64)
			data[1] = dims
		case 5:
			cnt = ln2vec(line, false, false)
			params := cnt.([][]int)
			data[2] = params
		}
	}

	//log.Println("XXXXX000000XXXXX")
	//log.Println(model)
	data[3] = file
	return
}

// //fileExists checks if a file exists
// func fileExists(fname string) bool{
// 	info, err := os.Stat(fname)
// 	if os.IsNotExist(err) {
// 		return false
// 	}
// 	return !info.IsDir()
// }

//ModInp is the entry func from flag (calc) or menu (kass)
func ModInp(fname, term string, pipe bool) (err error){
	temp := os.Stdout
	if pipe{
		os.Stdout = nil
		
	}
	fname, err = filepath.Abs(fname)
	if err != nil{
		log.Println(ColorRed, err, ColorReset)
		return
	}
	frmtyp, mod, e := JsonInp(fname)
	if e != nil {
		log.Println(ColorRed, e, ColorReset)
		err = e
		return
	}
	err = CalcMod(mod, frmtyp, term, pipe)
	if pipe{
		os.Stdout = temp
		jsonstr, e := json.Marshal(mod)
		if e != nil{
			log.Println(ColorRed, e, ColorReset)
			return
		}
		fmt.Print(string(jsonstr))
		
	}
	return
}

//ModEpInp is the entry func for elastic plastic analyis (of frames and beams) from flag (calc) or menu (kass)
func ModEpInp(fname, term string, pipe bool) (err error){
	temp := os.Stdout
	if pipe{
		os.Stdout = nil
		
	}
	fname, err = filepath.Abs(fname)
	if err != nil{
		log.Println(ColorRed, err, ColorReset)
		return
	}
	_, mod, e := JsonInp(fname)
	if e != nil {
		log.Println(ColorRed, e, ColorReset)
		err = e
		return
	}
	mod.Term = term
	err = CalcEpFrm(mod)
	if pipe{
		os.Stdout = temp
		jsonstr, e := json.Marshal(mod)
		if err != nil{
			log.Println(ColorRed, e, ColorReset)
			return
		}
		fmt.Print(string(jsonstr))
		
	}
	return
}

//ModNpInp is the entry func from flag (calc) or menu (kass)
func ModNpInp(fname, term string, pipe bool) (err error){
	fname, err = filepath.Abs(fname)
	if err != nil{
		log.Println(ColorRed, err, ColorReset)
		return
	}
	frmtyp, mod, e := JsonInp(fname)
	if e != nil {
		log.Println(ColorRed, e, ColorReset)
		err = e
		return
	}
	err = CalcModNp(mod, frmtyp, term, pipe)
	
	return
}

//CalcInp is the entry func from main/flags
func CalcInp(fname, term, cmdz string, pipe bool) (err error){
	//add bolt and weld analysis
	switch cmdz{
		case "":
		err = ModInp(fname, term, pipe)
		case "nl":
		//non linear analysis
		//only truss for now
		fname, err = filepath.Abs(fname)
		if err != nil{
			log.Println(ColorRed, err, ColorReset)
			return
		}
		_, mod, e := JsonInp(fname)
		if e != nil {
			log.Println(ColorRed, e, ColorReset)
			err = e
			return
		}
		NlCalcTrs2d(mod, 0.0)
		case "ep":
		//elastic-plastic analysis
		err = ModEpInp(fname, term, pipe)
		case "np":
		err = ModNpInp(fname, term, pipe)
		case "blt","bolt","bolts":
		var b Blt
		jsonfile, e := ioutil.ReadFile(fname)
		if e != nil {
			err = e
			return 
		}
		err = json.Unmarshal([]byte(jsonfile), &b)
		if err != nil{return}
		err = BoltSs(&b)
		fmt.Println(b.Report)
		case "wld", "weld", "welds":
		var w Wld
		jsonfile, e := ioutil.ReadFile(fname)
		if e != nil {
			err = e
			return 
		}
		err = json.Unmarshal([]byte(jsonfile), &w)
		if err != nil{return}
		err = WeldSs(&w)
		fmt.Println(w.Report)
		case "sec", "sect", "section", "secprp":
		
	}
	return
}

//JsonOut writes a model to a json file in /data/out
func (mod *Model) JsonOut() (fname string, err error){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	file, err := json.MarshalIndent(mod, "", " ")
	if err != nil{
		return
	}
	fname = filepath.Join(basepath, "../data/out",mod.Id+"_out.json")
	err = ioutil.WriteFile(fname, file, 0644)
	return
}


//JsonOut saves a frm2d struct to a json file in data/out
func (f *Frm2d) JsonOut() (fname string, err error){
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


//JsonOut saves a frm2d struct to a json file in data/out
func (t2d *Trs2d) JsonOut() (fname string, err error){
	_, base, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(base)
	file, err := json.MarshalIndent(t2d, "", " ")
	if err != nil{
		return
	}
	fname = filepath.Join(basepath, "../data/out",t2d.Id+"_out.json")
	err = ioutil.WriteFile(fname, file, 0644)
	return
}
