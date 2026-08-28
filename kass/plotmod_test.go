package barf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDraw2d(t *testing.T){
	var examples = []string{"akms4.1"}
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples")
	t.Log("testing 2d model plots")
	for i, ex := range examples {
		t.Log("ex. no->",i,"file->",ex)
		fname := filepath.Join(datadir,ex+".json")
		ModInp(fname, "qt", false)
	}
}

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
		t.Log("ex. no->",i,"file->",ex)
		fname := filepath.Join(datadir,ex+".json")
		ModInp(fname, "qt", false)
		
	}
	
}
