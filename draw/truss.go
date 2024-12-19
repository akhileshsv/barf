package barf

import (
	"fmt"
	"math"
)

//PlotGenTrs plots coords and mprp(ms) generated from TrussGen
//YEOLDE thing
func PlotGenTrs(coords [][]float64, ms [][]int, term, title string) (txtplot string, err error){
	var data string
	//all yr ys are belong to 1.0
	for idx, v := range coords {
		data += fmt.Sprintf("%v %v %v\n", v[0], v[1], idx+1)
	}
	data += "\n\n"
	for i, m := range ms {
		jb := coords[m[0]-1]
		je := coords[m[1]-1]
		data += fmt.Sprintf("%f %f %f %f %v %v %f\n",jb[0],jb[1],je[0],je[1],i+1,m[3],math.Cos((je[0]-jb[0])/(je[1]-jb[1])))
	}
	data += "\n\n"
	skript := "t2dgenplot.gp"
	txtplot, err = Dumb(data, skript, term, title, "", "", "")
	return 
} 

