package barf

import (
	"fmt"
)

//FORGET NEW STRUCTS JUST USE A MODEL
// //BmCs is a struct to store CBm fields
// //ADD ELASTIC FOUNDATION/SPRING SUPPORTS HERE
// //IF ONLY ONE WROTE CBm HERE
// type BmCs struct{
// 	Mstr    string
// 	Mtyp    int
// 	Em 	float64
// 	Cp      []float64
// 	Lspans 	[]float64
// 	Nspans 	float64
// 	Dims 	[]float64
// 	Styp 	int
// 	Ldcases []float64
// 	Selfwt  bool
// 	Lclvr   []float64
// 	Rclvr   []float64
// 	Term    string
// }

// {"Id":"2.3Hulse",
//  "Cmdz":["1db","mks","1","dl,ll","mredis,0.30"], 
//  "Ncjt":2, 
//  "Coords": [[0],[6],[10],[16]],
//  "Supports":[[1,-1,0],[2,-1, 0],[3,-1,0],[4,-1,0]],
//  "Em":[[25e9]],"Cp":[[1200e-9]],
//  "Mprp": [[1,2,1,1,0], [2,3,1,1,0], [3,4,1,1,0]],
//  "Jloads": [], 
//  "Msloads":[[1,3,25,0,0,0,1],[1,3,10,0,0,0,2],[2,3,25,0,0,0,1],[2,3,10,0,0,0,2], [3,3,25,0,0,0,1],[3,3,10,0,0,0,2]],
//  "PSFs":[1.4,1.0,1.6,0.0],
//  "Clvrs":[[0,0],[0,0]]	}

//InitBm
func (mod *Model) InitBm()(err error){
	if mod.Mstr == ""{		
		switch mod.Mtyp{
			case 1:
			mod.Mstr = "rcc"
			case 2:
			mod.Mstr = "stl"
			case 3:
			mod.Mstr = "tmbr"
		}	
	}
	if len(mod.Dims) == 0 && len(mod.Cp) == 0{
		return fmt.Errorf("missing cross section dims(%v)/prop(%v)-",mod.Dims, mod.Cp)
	}
	if mod.Sectyp == 0{
		mod.Sectyp = 1
	}
	return
}

func (mod *Model) GenBmGeom()(err error){
	if len(mod.Xs) == 0{
		return fmt.Errorf("no spans specified for beam-%v",mod.Xs)
	}
	if len(mod.Lspans) == 1{
		if mod.Nspans == 0{
			mod.Nspans = 1
		}
		lspan := mod.Lspans[0]
		for i := 0; i < mod.Nspans; i++{
			mod.Lspans = append(mod.Lspans, lspan)
		}
	}
	switch mod.Same{
		case true:
		//now gen mprp and supports
		if len(mod.Xs) == 0{
			x := 0.0
			mod.Xs = append(mod.Xs, x)
			for _, lspan := range mod.Lspans{
				x += lspan
				mod.Xs = append(mod.Xs, lspan)
			}
		}
	}
	for i, x := range mod.Xs{
		mod.Coords = append(mod.Coords, []float64{x,0.0})
		mod.Supports = append(mod.Supports, []int{i + 1, -1, 0})
		if i != len(mod.Xs)-1{
			mod.Mprp = append(mod.Mprp, []int{i+1,i+2,1,1,0})
		}
	}
	//mod.Msloads = mod.Ldcases
	return
}

func (mod *Model) GenEm() (err error){
	switch mod.Mstr{
		case "rcc":
		case "stl":
		mod.Units = "nmm"
		switch mod.Code{
			case 1:
			mod.Grade = 410.0
			mod.Fy = 250.0
			mod.Em = [][]float64{{200000.0}}
			case 2:
			mod.Grade = 43.0
			mod.Fy = 250.0
			mod.Em = [][]float64{{210000.0}}
		}
		case "tmbr":
		mod.Units = "nmm"
		if mod.Grade == 0{
			mod.Grade = 3
		}
		if mod.Group == 0{
			mod.Group = 2
		}
		switch mod.Group{
			case 1:
			mod.Em  = [][]float64{{12.6e3}}
			mod.Pg = 0.85
			case 2:
			mod.Em  = [][]float64{{9.8e3}}
			mod.Pg = 0.75
			case 3:
			mod.Em  = [][]float64{{5.6e3}}
			mod.Pg = 0.65
		}
	}
	return
}

// //BmAz analyzes/generates a beam model (for tmbr/bash funcs)
// func (mod *Model) BmAz()(err error){
// 	err = mod.InitBm()
// 	if err != nil{
// 		return
// 	}
// 	err = mod.GenBmGeom()
// 	if err != nil{
// 		return
// 	}
// 	err = mod.GenBmLoads()
// 	if err != nil{
// 		return
// 	}
// 	err = mod.BmCalc()
// 	if err != nil{
// 		return
// 	}
// 	return
// }

