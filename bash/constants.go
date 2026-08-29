package barf

/*
   from is 800
Unit mass of steel, p = 7850 kg/m3
Modulus
of elasticity E = 2.0 x 10e5 N/mm2 (MPa)
Poisson's
ratio, u = 0.3
Modulus
of rigidity,
G = 0.769 x 10e5 N/mm2 (MPa)
Co-efficient
of thermal
expansion
cx = 12 x 10e-6/deg c
*/

var (
	//DELETE DIS
	//USE THE COLOR PACKAGE https://github.com/fatih/color -  this is clearly not happening 
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m" 
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
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
