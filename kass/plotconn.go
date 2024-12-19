package barf

import (
	"fmt"
)


//Draw draws a bolt group
func (b *Blt) Draw()(data string){
	var cdata string
	for i, bc := range b.Bc{
		cdata += fmt.Sprintf("%f %f %f %v\n",bc[0],bc[1],b.Dias[i],i) 
	}
	data += cdata
	return
}
