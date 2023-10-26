package barf

//ya this is not happening, IGNORE
import (
	"fmt"
	flay"barf/flay"
)

//Flrgrd dreams of a day when every floor is a grid that is then mapped onto a kass.Model and so on
type Flrgrd struct {
	Fgrd [][]int
	Dimx,Dimy []float64
	Dimz []float64
	Rmap map[int]*flay.Rm
	Wmap map[flay.Pt][]*flay.Wall
	Svec [][]int
	Sgrd [][]int
}

func (f *Flrgrd) Init(){
	f.Rmap, f.Wmap = flay.LoutGen(0,0,0, f.Fgrd,0.0,0.0, f.Dimx, f.Dimy)
	outstr := flay.PltLout(f.Rmap)
	fmt.Println(outstr)
	f.ReadSlb()
}


func (f *Flrgrd) ReadSlb(){
	//read slab vec, update beams based on slab vec
	//okay slab ids HAVE TO BE a sekwens
	for i, vec := range f.Svec{
		fmt.Println("slab->",i+1,vec)
	}
}
