package barf

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"path/filepath"
	//"sort"
	//"math"
)

func PySubFrm(coords [][]float64, ms [][]int, foldr string){
	var data string
	data += "|["
	for _, c := range coords{
		data += fmt.Sprintf("[%v,%v],",c[0],c[1])
	}
	data += "]|["
	for _, m := range ms{
		data += fmt.Sprintf("[%v,%v,%v,%v,%v],",m[0],m[1],m[2],m[3],m[4])
	}
	data += "]|"	
	f, e := os.CreateTemp("", "barf")

	if e != nil {
		log.Println(e)
	}
	defer f.Close()
	//defer os.Remove(f.Name())
	fmt.Println(data)
	_, e = f.WriteString(data)
	if e != nil {
		log.Println(e)
	} 
	prg := "python"
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	arg0 := filepath.Join(basepath,"pysubfrm.py")
	arg1 := f.Name()
	arg2 := foldr
	cmd := exec.Command(prg,arg0,arg1,arg2)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	outstr, errstr := stdout.String(), stderr.String()
	//if errstr != "" {
	//	log.Println(errstr)
	//}
	outstr = outstr + errstr
	fmt.Println(outstr)
}


func PyFrmDat(fname, foldr string, plotchn chan string){
	prg := "python"
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	arg0 := filepath.Join(basepath,"pyfrmdat.py")
	//arg0 := "pyfrmdat.py"
	arg1 := fname
	arg2 := foldr
	cmd := exec.Command(prg,arg0,arg1,arg2)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	outstr, errstr := stdout.String(), stderr.String()
	//if errstr != "" {
	//	log.Println(errstr)
	//}
	//log.Println(outstr)
	outstr = outstr + errstr
	plotchn <- outstr	
}


func PyFrm3d(nodecords map[int][]float64,supports map[int][]int, cols, beamxs, beamys, slabnodes [][]int, foldr string, plotchn chan string) {
	var data string
	data += "|{"
	for node, coords := range nodecords {
		data += fmt.Sprintf("%v:[%v,%v,%v],", node, coords[0], coords[1], coords[2])
	}
	data += "}|{"
	for node, c := range supports {
		data += fmt.Sprintf("%v:[%v,%v,%v,%v,%v,%v],", node, c[0], c[1], c[2], c[3], c[4], c[5])
	}
	data += "}|["
	for _, c := range cols {
		data += fmt.Sprintf("[%v,%v],", c[0], c[1])
	}

	data += "]|["
	for _, beamx := range beamxs {
		data += fmt.Sprintf("[%v,%v],", beamx[0], beamx[1])
	}

	data += "]|["
	for _, beamy := range beamys {
		data += fmt.Sprintf("[%v,%v],", beamy[0], beamy[1])
	}
	data += "]|["
	for _, s := range slabnodes {
		data += fmt.Sprintf("[%v,%v,%v,%v],", s[0], s[1], s[2], s[3])
	}
	data += "]|"
	f, e := os.CreateTemp("", "barf")

	if e != nil {
		log.Println(e)
	}
	defer f.Close()
	defer os.Remove(f.Name())
	_, e = f.WriteString(data)
	if e != nil {
		log.Println(e)
	}
	//THIS WILL BREAK 
	prg := "python"//"C:\\Users\\Admin\\AppData\\Local\\Programs\\Python\\Python38\\python.exe"
	arg0 := "py_plot.py"
	//_, b, _, _:= runtime.Caller(0)
	//basepath := filepath.Dir(b)
	//fmt.Println("BASEPATH-->",basepath)
	//arg0 := filepath.Join(basepath,"pyfrm3d.py")
	arg1 := f.Name()
	arg2 := foldr
	cmd := exec.Command(prg,arg0,arg1,arg2)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	outstr, errstr := stdout.String(), stderr.String()
	//if errstr != "" {
	//	log.Println(errstr)
	//}
	outstr = outstr + errstr
	log.Println(outstr)
	plotchn <- outstr
}
