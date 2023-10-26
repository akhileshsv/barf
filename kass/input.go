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
	//log.Println(mod)
	frmtyp := mod.Cmdz[0]
	frmtyp = strings.ToLower(frmtyp)
	frmtyp = strings.TrimSpace(frmtyp)
	switch frmtyp {
	case "2dt","3dt","1db","2df","3df","3dg":
		return frmtyp, &mod, nil	
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

//fileExists checks if a file exists
func fileExists(fname string) bool{
	info, err := os.Stat(fname)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//ModInp is the entry func from flag (calc) or menu (kass)
func ModInp(fname, term string) (err error){
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
	err = CalcMod(mod, frmtyp, term)
	return
}

//ModEpInp is the entry func for elastic plastic analyis (of frames and beams) from flag (calc) or menu (kass)
func ModEpInp(fname, term string) (err error){
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
	return
}

//ModNpInp is the entry func from flag (calc) or menu (kass)
func ModNpInp(fname, term string) (err error){
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
	err = CalcModNp(mod, frmtyp, term)
	
	return
}

//CalcInp is the entry func from main/flags
func CalcInp(fname, term, cmdz string) (err error){
	//add bolt and weld analysis
	switch cmdz{
		case "":
		err = ModInp(fname, term)
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
		err = ModEpInp(fname, term)
		case "np":
		err = ModNpInp(fname, term)
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
	}
	return
}
