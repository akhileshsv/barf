
package barf

import (
	"errors"
)

var SecMapCol = map[string]int{
	"circle":    0,
	"rectangle": 1,
	"ela":       2,
	"elb":       3,
	"elc":       4,
	"eld":       5,
	"plus":      6,
	"hex":       7,
	"oct":       8,
	"tee":       9,
	"input":     -1,
}

//these are pretty much worthless? return errors.New()
var (
	ErrSpacing = errors.New("rebar spacing error, change depth and area")
	ErrShear   = errors.New("shear check failed")
	ErrDepth   = errors.New("insufficient depth")
	ErrCvr     = errors.New("effective cover error")
	ErrFtngDim = errors.New("base size inadequate at ULS")
	ErrCAxL    = errors.New("axial load too large")
	ErrCEffHt  = errors.New("effective height check failed")
	ErrDim     = errors.New("revise dimensions")
	ErrD       = errors.New("revise design category")
	ErrDraw    = errors.New("plot error")
	ErrParams  = errors.New("input param error")
	ErrIter    = errors.New("iteration error")
)

var (
	nomcvrIs  float64 = 40.0
	nomcvrCol float64 = 40.0
	nomcvrBm  float64 = 25.0
	nomcvrSlb float64 = 30.0
)

var (
	fckcol float64 = 15.0
	fycol float64 = 415.0
	efcvrcol float64 = 60.5
	nlayerscol int = 2
	dtypecol int = 1
	rtypecol int = 0
	coltype string = "rectangle"
	stypecol int = 1
	menus = map[int]string{
		1:"sub frame analysis",
		2:"2d frame analysis",
		3:"3d frame analysis",
		4:"column analysis",
		5:"beam analysis",
		6:"slab analysis",
	}
	f2dmenus = map[int]string{
		1:"exit",
	}
	colmenus = map[int]string{
		1:"design",
		2:"view/edit globals",
		3:"section input",
	}
	
	slbmenus = map[int]string{
		1:"single span design",
		2:"continuous slab",
	}
	allColz = make(map[int]*RccCol)
	allSlbz = make(map[int]RccSlb)
	slbTyps = map[int]string{
		1:"1 way/clvr",
		2:"2 way",
		3:"ribbed",
		4:"waffle",
	}
	plotOpts = map[int]string{
		1:"dumb",
		2:"caca",
		3:"qt",
	}
	allColSecs = map[int]string{
		0:"circle",
		1:"rectangle",
		2:"diamond",
		3:"l",
		4:"t",
		5:"plus",
	}
	allRtypes = map[int]string{
		0:"symmetrical",
		1:"unsymmetrical",
	}
	frmtypes = map[int]string{
		0:"sub frame",
		1:"one way slab",
		2:"two way slab",
		3:"flat slab",
	}
	NCycles = 4
)

var (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m" //this is actually magenta
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"

	/*
	ColorReset  = ""
	ColorRed    = ""
	ColorGreen  = ""
	ColorYellow = ""
	ColorBlue   = ""
	ColorPurple = "" //this is actually magenta
	ColorCyan   = ""
	ColorWhite  = ""
	*/
	IconMosh = `
                  __ 
  __ _  ___  ___ / / 
 /  ' \/ _ \(_-</ _ \
/_/_/_/\___/___/_//_/
          
          rcc design                                                          
`
	IconCol = `
           __              
 _______  / /_ ____ _  ___ 
/ __/ _ \/ / // /  ' \/ _ \
\__/\___/_/\_,_/_/_/_/_//_/
             
       rcc column routines                         
`
	IconSlb = `
      
     _      _    
  __| |__ _| |__ 
 (_-< / _  | '_ \
 /__/_\__,_|_.__/
                 

       rcc slab routines                         
`
)

func Nocolor(){
	ColorReset  = ""
	ColorRed    = ""
	ColorGreen  = ""
	ColorYellow = ""
	ColorBlue   = ""
	ColorPurple = "" 
	ColorCyan   = ""
	ColorWhite  = ""
	return
}
