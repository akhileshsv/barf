package main

import (
	"flag"
	"io"
	"log"
	"os"
	"math/rand"
	"time"
	kass "barf/kass"
	mosh "barf/mosh"
	bash "barf/bash"
	tmbr "barf/tmbr"
	menu "barf/menu"
	srvr "barf/srvr"
	//flay "barf/flay"
)

var (
	//start menu
	tui = flag.Bool("tui",false,"start tui menu")
	//input file (.json) path
	inf = flag.String("inf", "", "input file path")
	//haha. this is never getting done (SETUP THE PI)
	srv = flag.Bool("srv", false, "start gin server")
	//term is also part of struct? use this to override in flags
	term = flag.String("term", "dumb", "gnuplot type (dumb,wxt,qt)")
	//calc routes to analysis via kass
	calc = flag.Bool("calc", false, "stiffness analysis")
	//global tweak is good to have?
	tweek = flag.Bool("tweek",false, "(global) tweak all")
	//rcc has to be string
	rcc = flag.String("rcc", "", "rcc design string")
	//so does steel
	stl = flag.String("stl", "", "steel design string")
	//so does timber lol nothing has been done here
	wood = flag.String("wood","","wood/timber design string")
	//additional string param thing
	cmdz = flag.String("cmdz","","commands/params string")
	chnf = flag.String("chnf","","chain file path")
	pipe = flag.Bool("pipe",false,"python pipe input")
	//is dis needed - check for os and change this?
	//drawterms = map[string]int{"mono":1,"dumb":2,"qt":3,"svg":4,"wxt":5}
)

func main() {
	log.SetFlags(log.Ltime | log.Ldate | log.Lmsgprefix)
	logfile, err := os.OpenFile("data/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(mw)
	flag.Parse()
	if *pipe{
		mw = io.MultiWriter(logfile)
		log.SetOutput(mw)
	}
	rand.Seed(time.Now().Unix())
	/*
	   
	_, okdraw := drawterms[*term]
	if !okdraw {
		log.Println("invalid terminal string")
		log.Println("switching to dumb terminal")
		*term = "dumb"
	}
	*/
	switch {
	case *tui:
		//start menu
		menu.InitMenu(*term,*cmdz)
		os.Exit(0)
	case *srv:
		kass.Nocolor()
		mosh.Nocolor()
		bash.Nocolor()
		tmbr.Nocolor()
		srvr.Srvr()
		//os.Exit(0)
	case *calc || *rcc != "" || *stl != "" || *wood != "":
		switch{
			case *inf == "":
			log.Println("no input file specified")
			flag.Usage()
			default:
			if *pipe{
				mw = io.MultiWriter(logfile)
				log.SetOutput(mw)
			}
			switch{
				case *calc:
				kass.CalcInp(*inf,*term,*cmdz,*pipe)
				case *rcc != "":
				mosh.CalcInp(*rcc,*cmdz,*inf,*chnf,*term,*tweek,*pipe)
				case *stl != "":
				bash.CalcInp(*stl,*cmdz,*inf,*term,*pipe)
				case *wood != "":
				tmbr.CalcInp(*wood,*cmdz,*inf,*term,*pipe)
			}
		}
	default:
		log.Println("no flag specified")
		flag.Usage()
	}
	os.Exit(0)
}
