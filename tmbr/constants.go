package barf

import (
	"errors"
)

var (
	//DELETE DIS goddamn

	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m" 
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"

)

var (
	ErrDim = errors.New("dimension error")
	ErrGrp = errors.New("invalid group/grade")
	ErrInp = errors.New("input param error")
)

var (
	matZ = []string{"timber","bamboo"}
	tmbrGrp = []string{"a","b","c"}
	rDims = [][]float64{
		{1,1.5},
		{1,2},
		{1,2.5},
		{1,3},
		{2,2},
		{2,2.5},
		{2,3},
		{2,3.5},
		{2,4},
		{2,5},
		{2,6},
		{3,3},
		{3,3.5},
		{3,4},
		{3,5},
		{3,6},
		{3.5,4},
		{4,4},
		{4,5},
		{4,6},
		{4,8},
	}

	circDims = [][]float64{{3.0},{4.0},{5.0},{6.0},{7.0},{8.0},{9.0},{10.0},{11.0},{12.0}}
	rectBs = []float64{10,15,20,25,30,40,50,60,80,100,120,140,160,180,200}
	rectDs = []float64{40,50,60,80,100,120,140,160,180,200}
	rectDims = [][]float64{}
	plyDs = []float64{6,9,12,16,19,25}
	deckDs = []float64{25,30,35,40,45,50}
	tmbrKost = []float64{50000,1000,20}//wood/m3 (incl fab), coating, cost/bolt 
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
