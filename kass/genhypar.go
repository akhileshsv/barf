package barf

import (
	"fmt"
	"math"
	// draw"barf/draw"
)

//Hypargen generates points for a 3-d hypar structure
//https://mathworld.wolfram.com/HyperbolicParaboloid.htm
func Hypargen(typ int, a, b, c, lx, ly float64)(coords [][]float64, ms [][]int){
	var data string
	var z float64
	//generate coords from x(-a/2,a/2), y (-b/2,b/2)
	for x := -a/2.0; x <= a/2.0; x += lx{
		for y := -b/2.0; y <= b/2.0; y += ly{
			switch typ{
				case 1:
				z = math.Pow(y,2)/math.Pow(b,2) - math.Pow(x,2)/math.Pow(a,2)
				case 2:
				z = x * y
			}
			coords = append(coords, []float64{x,y,z})
			data += fmt.Sprintf("%f %f %f %v\n",x,y,z,len(coords))
		}
	}
	data += "\n\n"
	nx := int(math.Round(a/lx) + 1.0)
	ny := int(math.Round(b/ly) + 1.0)
	fmt.Println("nx,ny,len(coords),nx*ny",nx,ny,len(coords),nx*ny)
	mems := make(map[string]bool)
	for i := range coords{
		idx := i + 1
		edges := [][]int{{idx, idx + 1},{idx, idx + ny},{idx + ny, idx + ny + 1},{idx + 1, idx + ny + 1}}
		for i, e := range edges{
			jb := e[0]; je := e[1]
			if jb < 1 || jb > len(coords) || je < 1 || je > len(coords){
				continue
			}
			mdx := getmemidx(jb,je)
			switch i{
				case 0:
				//left edge
				if Dist2d(coords[jb-1],coords[je-1]) == ly{
					if _, ok := mems[mdx]; !ok{
						ms = append(ms, []int{jb,je})
					}
				}
				case 1:
				//bottom edge
				
				if Dist2d(coords[jb-1],coords[je-1]) == lx{
					if _, ok := mems[mdx]; !ok{
						ms = append(ms, []int{jb,je})
					}
				}
				case 2:
				//right edge
				if Dist2d(coords[jb-1],coords[je-1]) == ly{
					if _, ok := mems[mdx]; !ok{
						ms = append(ms, []int{jb,je})
					}
				}
				case 3:
				//top edge
				
				if Dist2d(coords[jb-1],coords[je-1]) == lx{
					if _, ok := mems[mdx]; !ok{
						ms = append(ms, []int{jb,je})
					}
				}
			}
		}
	}
	for i, mem := range ms{
		pb := coords[mem[0]-1]; pe := coords[mem[1]-1]
		data += fmt.Sprintf("%f %f %f %f %f %f %v\n",pb[0],pb[1],pb[2],pe[0]-pb[0],pe[1]-pb[1],pe[2]-pb[2],i)
		//data += fmt.Sprintf("%f %f %f %v\n",pe[0],pe[1],pe[2],i)
		//data += "\n"
	}
	//draw.Draw(data, "d3.gp","qt","","","hypar")
	return
}

