package barf

import (
	"os"
	"fmt"
	"testing"
	"path/filepath"
	"gonum.org/v1/gonum/mat"
)


func TestCsv(t *testing.T) {
	var rezstring string
	nrooms := 8
	dx := 20.0
	dy := 20.0
	grid, err := CsvIntMat("lout_prex.csv")
	if err != nil {fmt.Println(err)}
	fmt.Println(ColorCyan,grid,ColorReset)
	ny := len(grid)
	nx := len(grid[0])
	cmat := mat.NewDense(8, 8, []float64{
		0, 1, 1, 1, 1, 1, 1, 1,
		1, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1, 1,
		1, 1 ,1, 1, 0, 1, 1, 1,
		1, 1 ,1, 1, 1, 0, 1, 1,
		1, 1 ,1, 1, 1, 1, 0, 1,
		1, 1 ,1, 1, 1, 1, 1, 0,
	})
	fmat := mat.NewDense(8, 8, []float64{
		0, 45, 15, 25, 10, 5, 0, 0,
		0, 0, 0, 30, 25, 15, 0, 0,
		0, 0, 0, 0, 5, 10, 0, 0,
		0, 20, 0, 0, 35, 0, 0, 0,
		0, 0 ,0, 0, 0, 65, 35, 0,
		0, 5 ,0, 0, 25, 0, 65, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	})
	
	rmap, _, _, _, _ := LoutGen(nrooms,nx,ny, grid, dx, dy,[]float64{},[]float64{})

	outstr := PltLout(rmap)
	rezstring += "\nNEW PLOT\n"
	rezstring += outstr
	fmt.Println(outstr)
	//fmt.Println(rmap,nodemap)
	combos := CraftCombos(rmap)
	fmt.Println(combos)
	rcent := make([]*Pt,nrooms)
	for i:=1; i <=len(rmap); i++ {
		rcent[i-1] = rmap[i].Centroid
	}
	fmt.Println(rcent)
	rezstring += fmt.Sprintf("init cost --> \n %v\n",CraftCost(rcent, cmat, fmat))
	costs, minidx, gridz := CraftCombosEval(rcent, rmap, combos, cmat, fmat, grid, dx, dy)
	swpz := map[int]string{0:"area",1:"adj"}
	for idx, cost := range costs {
		if idx == minidx {
			rezstring += "dis MINIMUM COST\n"
		}
		rezstring += ColorCyan
		rezstring += fmt.Sprintf("combo - swap %v and %v cause %v total cost %v\n",combos[idx][0],combos[idx][1],swpz[combos[idx][2]],cost)
		rezstring += ColorPurple
		rezstring += fmt.Sprintf("%v\n",gridz[idx])
		rmap, _, _, _,_ = LoutGen(nrooms,nx,ny, gridz[idx], dx, dy,[]float64{},[]float64{})
		outstr := PltLout(rmap)
		rezstring += outstr
	}
	fmt.Println(rezstring)
	
}

func TestIARE(t *testing.T) {
	var rezstring string
	posvec := [][]int{
		{1,5,4},
		{2,3,3},
	}
	dimvecx := []float64{4.0,2.0,4.0}
	dimvecy := []float64{4.0,4.0}
	nrooms := 5
	cmat := mat.NewDense(5, 5, []float64{
		0, 1, 1, 1, 1,
		1, 0, 1, 1, 1,
		1, 1, 0, 1, 1,
		1, 1, 1, 0, 1,
		1, 1 ,1, 1, 0,
	})
	fmat := mat.NewDense(5, 5, []float64{
		0, 5, 2, 4, 0,
		0, 0, 2, 5, 0,
		2, 0, 0, 0, 5,
		3, 0, 1, 0, 0,
		0, 0, 2, 0, 0,
	})
	grid := GridGen(nrooms, posvec, dimvecx, dimvecy, 0.0, 0.0)
	//fmt.Println(grid.Vec)
	rezstring += fmt.Sprintf("grid x %v grid y %v \n GRID \n %v", grid.Dx, grid.Dy, grid.Vec)
	rmap, _, _, _,_ := LoutGen(grid.Nr,grid.Nx,grid.Ny, grid.Vec, grid.Dx, grid.Dy,[]float64{},[]float64{})
	outstr := PltLout(rmap)
	//rezstring += "\nNEW PLOT\n"
	rezstring += outstr
	fmt.Println(outstr)
	//fmt.Println(rmap,nodemap)
	combos := CraftCombos(rmap)
	//fmt.Println(combos)
	rcent := make([]*Pt,nrooms)
	for i:=1; i <=len(rmap); i++ {
		rcent[i-1] = rmap[i].Centroid
	}
	//fmt.Println(rcent)
	rezstring += fmt.Sprintf("init cost --> \n %v\n",CraftCost(rcent, cmat, fmat))
	costs, minidx, gridz := CraftCombosEval(rcent, rmap, combos, cmat, fmat, grid.Vec, grid.Dx, grid.Dy)
	swpz := map[int]string{0:"area",1:"adj"}
	for idx, cost := range costs{
		if idx == minidx {
			rezstring += "dis MINIMUM COST\n"
		}
		rezstring += ColorCyan
		rezstring += fmt.Sprintf("combo - swap %v and %v cause %v total cost %v\n",combos[idx][0],combos[idx][1],swpz[combos[idx][2]],cost)
		rezstring += ColorPurple
		//rezstring += fmt.Sprintf("%v\n",gridz[idx])
		//rmap, _ = LoutGen(nrooms,grid.Nx,grid.Ny, gridz[idx], grid.Dx, grid.Dy,[]float64{},[]float64{})
		outstr := Plotgrid(gridz[idx],grid.Dx,grid.Dy)
		rezstring += outstr
	}
	fmt.Println(rezstring)	
}

func TestCraft(t *testing.T){
	var examples = []string{"craft1","craft2"}
	//var rezstring string
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/flay")
	t.Log(ColorPurple,"testing CRAFT algo\n",ColorReset)
	for i, ex := range examples{
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		c, err := ReadCrft(fname)
		if err != nil{
			t.Log(err)
			t.Fatal("CRAFT algo test failed")
		}
		c.Verbose = true
		err = c.Craft()
		if err != nil{
			t.Log(err)
			t.Fatal("CRAFT algo test failed")
		}
	}
}
