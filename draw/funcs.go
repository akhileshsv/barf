package barf

import (
	"os"
	"fmt"
	"log"
	"bytes"
	"os/exec"
	"runtime"
	"strings"
	"io/ioutil"
	"path/filepath"
)

//skriptpath returns the abs script path
func skriptpath(skript string) (string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, skript)
}

//svgpath returns the abs path to save .svg files
func svgpath(fname, term string) (string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	fname = fmt.Sprintf("%s.%s",fname,term)
	return filepath.Join(basepath, "../data/out",fname)
}


//svgkong reads an svg embeds kongtext in style defs and saves
//peak noob to be doing this TWICE (twice? it's everywhere)
func svgkong(pltstr string){
	bytes, err := ioutil.ReadFile(pltstr)
	if err != nil {
		log.Println(err)
	}
	rez := string(bytes)
	//log.Println(rez)
	d , err := os.ReadFile(skriptpath("../data/out/kongdef.txt"))
	if err != nil{
		log.Fatal(err)
	}
	defs := string(d)
	title := "<title>Gnuplot</title>"
	ndefs := defs + title
	rez = strings.Replace(rez,title,ndefs,1)
	err = os.WriteFile(pltstr,[]byte(rez),0666)
	if err != nil{
		log.Fatal(err)
	}
	
}

//Draw plots a gnuplot datafile "data" with a gnuplot script "skript"
func Draw(data, skript, term, folder, fname, title, xl, yl, zl string) (txtplot string, err error) {
	if xl == "" && yl == "" && zl == ""{
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
	
	if fname == ""{fname = "outfile"}
	//if term == "dxf"{term = "svg"}
	pltskript := skriptpath(skript)
	fpath := svgpath(fname, term)
	//log.Println("data filename",f.Name())
	//log.Println("plot name",p)
	//log.Println("plot params->pltskript,f.Name(),term,title,fn\n",pltskript,f.Name(),term,title,fpath)
	//exec.Command("gnuplot","-c",pltskript,f.Name(),term, title, fname)
	cmd := exec.Command("gnuplot","-c",pltskript,f.Name(),term,title,xl,yl,zl,fpath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	outstr, errstr := stdout.String(), stderr.String()
	txtplot = errstr+outstr
	// log.Println("frm gnuplot-",txtplot)
	if folder == "web"{
		
		// log.Println("calling svgkong at->",fpath)
		svgkong(fpath)
		txtplot = fname + ".svg"
	}
	err = nil
	return
}

//Dumb plots to the terminal - is this needs
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

