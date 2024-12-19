package barf

import (
	"os"
	"log"
	"path/filepath"
	"testing"
)

func TestSvgkong(t *testing.T){
	pltstr := "../data/out/rcc slab_511_section_511.svg"
	Svgkong(pltstr)
}


func TestDraw3d(t *testing.T){
	var examples = []string{"akms8.1","akms8.2","akms8.4"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	for i, ex := range examples {
		//if i != 2{continue}
		log.Println("ex. no->",i,"file->",ex)
		fname := filepath.Join(datadir,ex+".json")
		ModInp(fname, "qt", false)
		
	}
	
}
