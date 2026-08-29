package barf

import (
	"errors"
)

var (
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
	ErrObj = errors.New("objective function error")
	ErrDim = errors.New("input dim error")
	ErrSel = errors.New("selection pool error")
)
