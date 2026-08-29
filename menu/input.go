package barf

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"runtime"
	"path/filepath"
	kass"barf/kass"
	mosh"barf/mosh"
)

//all these read funcs are not needed, one thinks?

func getFrm2d(inpt int)(f kass.Frm2d, err error){
	var filename string
	switch inpt{
		case 0:
		_, b, _, _:= runtime.Caller(0)
		basepath := filepath.Dir(b)
		filename = filepath.Join(basepath,"../data/json/rcfrm2d0.json")
		jsonfile, e := ioutil.ReadFile(filename)
		if e != nil{
			err = e
			return
		}
		err = json.Unmarshal([]byte(jsonfile),&f)
		jsonstr, _ := json.Marshal(f)
		fmt.Println(string(jsonstr))
		fmt.Println(ColorWhite,"edit struct below:")
		dec := json.NewDecoder(os.Stdin)
		for {
			e := dec.Decode(&f)
			if e == io.EOF {
				break
			}
			if e != nil {
				err = e
				return
				
			}
		}
		case 1:
		fmt.Println(ColorWhite,"enter json filename:",ColorReset)
		fmt.Scanln(&filename)
		filename = filepath.FromSlash(filename)
		jsonfile, e := ioutil.ReadFile(filename)
		if e != nil{
			err = e
			return
		}
		err = json.Unmarshal([]byte(jsonfile),&f)
	}
	return
	
}

func getSubFrm(inpt int)(sf mosh.SubFrm,err error){
	var filename string
	switch inpt{
		case 0:
		_, b, _, _:= runtime.Caller(0)
		basepath := filepath.Dir(b)
		filename = filepath.Join(basepath,"../data/json/rcsf0.json")
		sf, err = readSubFrm(filename)
		jsonstr, _ := json.Marshal(sf)
		fmt.Println(string(jsonstr))
		fmt.Println(ColorWhite,"enter edited json below:")
		dec := json.NewDecoder(os.Stdin)
		for {
			e := dec.Decode(&sf)
			if e == io.EOF {
				break
			}
			if e != nil {
				err = e
				return
			}
		}
		return
		case 1:
		bytestr, e := getjsonfile()
		if e != nil{err = e; return}
		err = json.Unmarshal(bytestr,&sf)
		return
	}
	return
}

func getRcBm(inpt int)(bm mosh.RccBm, err error){
	var filename string
	switch inpt{
		case 0:
		_, b, _, _:= runtime.Caller(0)
		basepath := filepath.Dir(b)
		filename = filepath.Join(basepath,"../data/json/rcbm0.json")
		jsonfile, e := ioutil.ReadFile(filename)
		if e != nil{
			err = e
			return
		}
		err = json.Unmarshal([]byte(jsonfile),&bm)
		jsonstr, _ := json.Marshal(bm)
		fmt.Println(string(jsonstr))
		fmt.Println(ColorWhite,"edit struct below:")
		dec := json.NewDecoder(os.Stdin)
		for {
			err := dec.Decode(&bm)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		case 1:
		fmt.Println("enter json filename:")
		fmt.Scanln(&filename)
		filename = filepath.FromSlash(filename)
		jsonfile, e := ioutil.ReadFile(filename)
		if e != nil{
			err = e
			return
		}
		err = json.Unmarshal([]byte(jsonfile),&bm)
	}
	return
	
}

func getRccSlb(inpt int)(s mosh.RccSlb, err error){
	var filename string
	switch inpt{
		case 0:
		_, b, _, _:= runtime.Caller(0)
		basepath := filepath.Dir(b)
		filename = filepath.Join(basepath,"../data/json/rcslb1way0.json")
		jsonfile, e := ioutil.ReadFile(filename)
		if e != nil{
			err = e
			return
		}
		//err = json.Unmarshal([]byte(jsonfile),&s)
		//jsonstr, _ := json.Marshal(s)
		fmt.Println(string(jsonfile))
		//fmt.Println(ColorWhite,"edit struct below:")
		bytestr, _ := readjsontxt("")
		err = json.Unmarshal(bytestr,&s)
		case 1:
		bytestr, e := getjsonfile()
		if e != nil{err = e; return}
		err = json.Unmarshal(bytestr,&s)
		return
	}
	return
	
}

func getRccCol(inpt int)(c mosh.RccCol, err error){
	var filename string
	switch inpt{
		case 0:
		_, b, _, _:= runtime.Caller(0)
		basepath := filepath.Dir(b)
		filename = filepath.Join(basepath,"../data/json/rccol0.json")
		jsonfile, e := ioutil.ReadFile(filename)
		if e != nil{
			err = e
			return
		}
		err = json.Unmarshal([]byte(jsonfile),&c)
		jsonstr, _ := json.Marshal(c)
		fmt.Println(string(jsonstr))
		fmt.Println(ColorWhite,"edit struct below:")
		dec := json.NewDecoder(os.Stdin)
		for {
			err := dec.Decode(&c)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		case 1:
		fmt.Println("enter json filename:")
		fmt.Scanln(&filename)
		filename = filepath.FromSlash(filename)
		jsonfile, e := ioutil.ReadFile(filename)
		if e != nil{
			err = e
			return
		}
		err = json.Unmarshal([]byte(jsonfile),&c)
	}
	return
	
}



func ReadJsonStr()(buf bytes.Buffer){
	//fmt.Printf("enter %s params in json below:\n",stype)
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				buf.WriteString(line)
				break 
			} else {
				log.Println(err)
			}   
		}   
		buf.WriteString(line)
	}   
	return
}

func ReadJsonMod() (mod *kass.Model, err error){
	err = json.NewDecoder(os.Stdin).Decode(mod)
	return
}

func ReadJsonCol() (c mosh.RccCol, err error){
	buf := ReadJsonStr()
	err = json.Unmarshal(buf.Bytes(), &c)
	return
}

func ReadJsonCBm() (c mosh.CBm, err error){
	buf := ReadJsonStr()
	err = json.Unmarshal(buf.Bytes(), &c)
	return
}

func ReadSlbFrmJson()(s mosh.RccSlb, err error){
	buf := ReadJsonStr()
	err = json.Unmarshal(buf.Bytes(), &s)
	return
}

func ReadSubFrmJson()(sf mosh.SubFrm, err error){
	buf := ReadJsonStr()
	err = json.Unmarshal(buf.Bytes(), &sf)
	return
}

func readSubFrm(filename string) (sf mosh.SubFrm, err error){
	jsonfile, e := ioutil.ReadFile(filename)
	if e != nil{
		err = e
		return
	}
	err = json.Unmarshal([]byte(jsonfile),&sf)
	return
}

/*
   YE OLDE
   
func getF3d(inpt int)(f gen.F3d, err error){
	var filename string
	switch inpt{
		case 0:
		_, b, _, _:= runtime.Caller(0)
		basepath := filepath.Dir(b)
		filename = filepath.Join(basepath,"../data/json/rcfrm3d0.json")
		jsonfile, e := ioutil.ReadFile(filename)
		if e != nil{
			err = e
			return
		}
		err = json.Unmarshal([]byte(jsonfile),&f)
		jsonstr, _ := json.Marshal(f)
		fmt.Println(string(jsonstr))
		fmt.Println(ColorWhite,"edit struct below:")
		dec := json.NewDecoder(os.Stdin)
		for {
			err := dec.Decode(&f)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		case 1:
		fmt.Println("enter json filename:")
		fmt.Scanln(&filename)
		filename = filepath.FromSlash(filename)
		jsonfile, e := ioutil.ReadFile(filename)
		if e != nil{
			err = e
			return
		}
		err = json.Unmarshal([]byte(jsonfile),&f)
	}
	return
	
}

*/
