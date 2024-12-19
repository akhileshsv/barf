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
	/*
	   styp 1 - abel 9.1, sp.33 - 2
	   styp 0 - ??
	   styp 26 - abel 10.1, 10.2
	   styp 4 - sp.33 - 3
	*/
	var rezstring string
	//abel rect col ex
	t.Log("abel ex. 1 rect col")
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
		Code:2,
		Dims:[]float64{50.0,75.0},
	}
	c.Init()
	ok, val, _ := ColChk(c)
	rezstring += fmt.Sprintf("abel ex. 9.1 rect. column\n")
	
	rezstring += fmt.Sprintf("perm stress - %.2f N/mm2 actual - %.2f N/mm2 ok? - %v\n",val[0], val[2], ok)
	
	_ = ColDz(c)	
	rezstring += fmt.Sprintf("sections from dz - %v\n",c.Rez)
	//sp.33 example 2 teak column
	
	t.Log("sp33 ex. 2 rect col")
	c = &WdCol{
		Styp:1,
		Lspan:3000.0,
		Prp:kass.Wdprp{
			Em:11020.0,
			Mc:0.5,
			Pg:0.5,
			Fc:9.6,
		},
		Pu:60000,
		Code:1,
		Dims:[]float64{150.0,150.0},
		
	}
	c.Init()
	ok, val, _  = ColChk(c)
	rezstring += fmt.Sprintf("sp.33 ex. 2 rect. column\n")
	rezstring += fmt.Sprintf("perm stress - %.2f N/mm2 actual - %.2f N/mm2 ok? - %v\n",val[0], val[2], ok)

	_ = ColDz(c)	
	rezstring += fmt.Sprintf("sections from dz - %v\n",c.Rez)

	t.Log("ramc ex. 13.3 rect col")
	c = &WdCol{
		Styp:1,
		Lspan:1500.0,
		Prp:kass.Wdprp{
			Em:10800.0,
			Pg:0.7,
			Fc:11.2,
		},
		Pu:252000,
		Code:1,
		Dims:[]float64{150.0,150.0},
	}
	c.Init()
	ok, val, _  = ColChk(c)
	rezstring += fmt.Sprintf("ramc ex. 13.3 rect. column\n")
	rezstring += fmt.Sprintf("perm stress - %.2f N/mm2 actual - %.2f N/mm2 ok? - %v\n",val[0], val[2], ok)
	
	_ = ColDz(c)	
	rezstring += fmt.Sprintf("sections from dz - %v\n",c.Rez)
	fmt.Println(rezstring)
	t.Fatal()
	t.Log("ramc ex. 13.6 round col")
	c = &WdCol{
		Styp:0,
		Lspan:1200.0,
		Prp:kass.Wdprp{
			Em:9500.0,
			Pg:0.6,
			Fc:7.0,
		},
		Pu:118000,
		Code:1,
		Dims:[]float64{150.0},
	}
	c.Init()
	ok, val, _  = ColChk(c)
	rezstring += fmt.Sprintf("ramc ex. 13.6 round column\n")
	rezstring += fmt.Sprintf("perm stress - %.2f N/mm2 actual - %.2f N/mm2 ok? - %v\n",val[0], val[2], ok)
	
	_ = ColDz(c)	
	rezstring += fmt.Sprintf("sections from dz - %v\n",c.Rez)

	//abel spaced column ex
	t.Log("abel ex. 10.1 spaced col")
	c = &WdCol{
		Styp:26,
		Lspan:3600.0,
		Prp:kass.Wdprp{
			Em:10600.0,
			Mc:0.5,
			Pg:0.5,
			Fc:9.0,
		},
		Pu:27000,
		Code:2,
		Dims:[]float64{75.0,100.0},
		Endc:3,
		Kerst:3.0,
		
	}
	c.Init()
	ok, val, _  = ColChk(c)	
	t.Log("rez - ok -> ",ok," val - > ",val)
	rezstring += fmt.Sprintf("abel ex. 10.1 spaced column\n")
	rezstring += fmt.Sprintf("perm stress - %.2f N/mm2 actual - %.2f N/mm2 ok? - %v\n",val[0], val[2], ok)
	
	t.Log("abel ex. 10.2 spaced col")
	c = &WdCol{
		Styp:26,
		Lspan:2500.0,
		Prp:kass.Wdprp{
			Em:12500.0,
			Mc:0.5,
			Pg:0.8,
			Fc:18.0,
		},
		Pu:10000,
		Code:2,
		Dims:[]float64{50.0,100.0},
		Endc:3,
		Kerst:3.0,
		
	}
	c.Init()
	ok, val, _  = ColChk(c)
	t.Log("rez - ok -> ",ok," val - > ",val)
	
	rezstring += fmt.Sprintf("abel ex. 10.2 spaced column\n")
	rezstring += fmt.Sprintf("perm stress - %.2f N/mm2 actual - %.2f N/mm2 ok? - %v\n",val[0], val[2], ok)

	
	t.Log("sp.33 ex. 3 box col")
	c = &WdCol{
		Styp:4,
		Lspan:3000.0,
		Grp:2,
		Tplnk:25.0,
		Pu:60000,
		Code:1,
		Dims:[]float64{125,125,75,75},
		Endc:1,
		
	}
	c.Init()
	ok, val, _  = ColChk(c)
	t.Log("rez - ok -> ",ok," val - > ",val)
	
	rezstring += fmt.Sprintf("sp.33 ex. 3 box column\n")
	rezstring += fmt.Sprintf("perm stress - %.2f N/mm2 actual - %.2f N/mm2 ok? - %v\n",val[0], val[2], ok)

	t.Log("ramc ex. 13.8 box col")
	c = &WdCol{
		Styp:4,
		Lspan:3000.0,
		Prp:kass.Wdprp{
			Em:12500.0,
			Mc:0.5,
			Pg:0.8,
			Fc:10.6,
		},
		Tplnk:50.0,
		Pu:940000,
		Code:1,
		Dims:[]float64{300,300,200,200},
		Endc:1,
		Solid:true,
	}
	c.Init()
	ok, val, _  = ColChk(c)
	t.Log("rez - ok -> ",ok," val - > ",val)
	rezstring += fmt.Sprintf("ramc ex. 13.8 solid box column\n")
	rezstring += fmt.Sprintf("perm stress - %.2f N/mm2 actual - %.2f N/mm2 ok? - %v\n",val[0], val[2], ok)
	
	wantstring := `abel ex. 9.1 rect. column
perm stress - 2.23 N/mm2 actual - 0.80 N/mm2 ok? - true
sp.33 ex. 2 rect. column
perm stress - 8.00 N/mm2 actual - 2.67 N/mm2 ok? - true
sp.33 ex. 2 rect. column
perm stress - 11.20 N/mm2 actual - 11.20 N/mm2 ok? - true
ramc ex. 13.6 round column
perm stress - 7.00 N/mm2 actual - 6.68 N/mm2 ok? - true
abel ex. 10.1 spaced column
perm stress - 4.14 N/mm2 actual - 1.80 N/mm2 ok? - true
abel ex. 10.2 spaced column
perm stress - 4.50 N/mm2 actual - 1.00 N/mm2 ok? - true
sp.33 ex. 3 box column
perm stress - 7.68 N/mm2 actual - 6.00 N/mm2 ok? - true
ramc ex. 13.8 solid box column
perm stress - 10.60 N/mm2 actual - 10.44 N/mm2 ok? - true
`
	if rezstring != wantstring{
		t.Log("output should be ->")
		fmt.Println(wantstring)
		t.Log("now is ->")
		fmt.Println(rezstring)
		t.Fatal("timber column section check test failed")
	}
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
