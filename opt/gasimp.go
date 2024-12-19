package barf

import (
	//"fmt"
	"strconv"
)

type Sat struct{
	Pos []string
	Fit float64
	Wt  float64
	Con float64
}

func int2bin(n int)(string){
	return strconv.FormatInt(int64(n), 2)
	
}

//chew on this for a while, see goldberg
