package barf

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"path/filepath"
)

//skriptpath returns the abs script path
func skriptpath(skript string) (string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, skript)
}

//svgpath returns the abs path to save .svg files
func svgpath(fname string) (string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	fname = fname + ".svg"
	return filepath.Join(basepath, "../data/out",fname)
}

//Draw plots a gnuplot datafile "data" with a gnuplot script "skript"
func Draw(data, skript, term, folder, fname, title string) (txtplot string, err error) {
	f, e1 := os.CreateTemp(folder, "barf")
	if e1 != nil {
		log.Println(e1)
		err = e1
		return
	}
	defer f.Close()
	defer os.Remove(f.Name())
	_, e1 = f.WriteString(data)
	if e1 != nil {
		log.Println(e1)
		err = e1
		return
	}
	pltskript := skriptpath(skript)
	fname = svgpath(fname)
	//log.Println("data filename",f.Name())
	//log.Println("plot name",p)
	cmd := exec.Command("gnuplot","-c",pltskript,f.Name(),term,fname)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
	outstr, errstr := stdout.String(), stderr.String()
	log.Println("outstr")
	log.Println(outstr, errstr)
	if errstr != "" {
		log.Println(errstr)
	
	}
	//outstr += errstr
	//if term == "dumb" || term == "caca"{
	txtplot = fmt.Sprintf("%s",outstr)
	//}
	err = nil
	return
}

//Dumb plots to the terminal
func Dumb(data, skript, term, title, xl, yl, zl string)(txtplot string, err error){
	//create temp files
	if xl == ""{
		xl = "x"; yl = "y"; zl = "z"
	}
	f, e1 := os.CreateTemp("", "barf")
	if e1 != nil {
		log.Println(e1)
		err = e1
		return
	}
	defer f.Close()
	defer os.Remove(f.Name())
	_, e1 = f.WriteString(data)
	if e1 != nil {
		log.Println(e1)
		err = e1
		return
	}
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	pltskript := filepath.Join(basepath,skript)
	cmd := exec.Command("gnuplot","-c",pltskript,f.Name(),term,title,xl,yl,zl)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return
	}
	outstr, errstr := stdout.String(), stderr.String()
	if errstr != "" {
		//log.Println(errstr)
	}
	if term == "dumb" || term == "caca"{
		txtplot = fmt.Sprintf("%v",outstr)
	}
	err = nil
	return
}

