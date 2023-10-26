package barf

import (
	"os"
	"fmt"
	"io/ioutil"
	"testing"
	"path/filepath"
	"encoding/json"
	kass"barf/kass"
)

func TestColJson(t *testing.T){
	var c WdCol
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/tmbr/col")
	filename := filepath.Join(datadir,"c1.json")
	fmt.Println(filename)
	jsonfile, err := ioutil.ReadFile(filename)
	if err != nil{fmt.Println(err)}
	err = json.Unmarshal([]byte(jsonfile), &c)
	if err != nil{fmt.Println(err)}
	fmt.Println(c.Prp)
}

func TestColChk(t *testing.T){
	c := &WdCol{
		Styp:1,
		Lspan:1500.0,
		Prp:kass.Wdprp{
			Em:6700.0,
			Mc:0.5,
			Pg:0.5,
			Fc:4.5,
		},
		Pu:3000,
		Code:1,
		Dims:[]float64{50.0,75.0},
	}
	c.Init()
	ok, val := ColChk(c)
	fmt.Println(ok,val)
}

func TestColDz(t *testing.T){
	c := &WdCol{
		Styp:1,
		Lspan:1500.0,
		Prp:kass.Wdprp{
			Em:6700.0,
			Fc:4.5,
			Mc:0.5,
			Pg:0.5,
		},
		Pu:3000,
		Code:1,
	}
	err := ColDz(c)
	if err != nil{
		fmt.Println(err)
	} else {
		for i := range c.Rez{
			switch c.Styp{
				case 0:
				fmt.Printf("%.1f dia\n",c.Rez[i][0]/25.4)
				case 1:
				fmt.Printf("%.1f x %.1f\n",c.Rez[i][0]/25.4, c.Rez[i][1]/25.4)
			}
			fmt.Printf("permissible - %0.2f N/mm2 euler - %0.2f N/mm2 actual- %0.2f N/mm2 s/d %0.2f\n",c.Vals[i][0], c.Vals[i][1], c.Vals[i][2], c.Vals[i][3])
		}
		mdx := 0
		fmt.Println("min-",c.Rez[mdx])
		fmt.Printf("stresses:\npermissible - %0.2f euler - %0.2f actual- %0.2f\n",c.Vals[mdx][0], c.Vals[mdx][1], c.Vals[mdx][2])
		switch c.Styp{
			case 0:
			fmt.Printf("%.1f dia\n",c.Rez[mdx][0]/25.4)
			case 1:
			fmt.Printf("%.1f x %.1f\n",c.Rez[mdx][0]/25.4, c.Rez[mdx][1]/25.4)
		}
	}
}
