package barf

import (
	"fmt"
	"testing"
	"math"
	kass"barf/kass"
)

func TestFc(t *testing.T){
	var ldrat float64
	var grp int
	ldrat = 17.0
	grp = 1
	fc, err := getFc(ldrat, grp)
	fmt.Println(fc, err)

	ldrat = 0.0
	grp = 2
	fc, err = getFc(ldrat, grp)
	fmt.Println(fc, err)

	ldrat = 57.0
	grp = 3
	fc, err = getFc(ldrat, grp)
	fmt.Println(fc, err)

}

func TestInit(t *testing.T){
	b := WdBm{
		Grp:1,
	}
	err := b.Prp.Init(b.Grp)
	fmt.Println(err)
	rez := b.printz()
	fmt.Println(rez)
}

func TestBmDzCs(t *testing.T){
	b := &WdBm{
		Styp:1,
		Lspan:3500.0,	
		Prp:kass.Wdprp{
			Em:9500,
			Fv:1.4,
			Fcp:2.8,
			Fcb:11.2,
			Ft:11.2,
			Pg:0.5,
		},
		Endc:1,
		DL:0.0,
		LL:0.5,
		Lbl:25.0,
		Rbl:25.0,
		Selfwt:true,
		Spam:true,
		Nsecs:3,
	}
	err := BmDzCs(b)
	if err != nil{
		fmt.Println(err)
	}
}

func TestPlyUdlSpn(t *testing.T){
	var d, wdl, wll float64
	wdl = 0.6; wll = 2.0
	d = 12.0
	fmt.Println("12 mm plywood")
	for scon :=0; scon < 4; scon++{
		lspan, err := PlyUdlSpn(d, wdl, wll, scon)
		if err == nil{fmt.Println(ColorRed,"safe span->",math.Round(math.Floor(lspan/25.0)*25.0),ColorReset)}
	}
	d = 19.0
	fmt.Println("19 mm plywood")
	for scon :=0; scon < 4; scon++{
		lspan, err := PlyUdlSpn(d, wdl, wll, scon)
		if err == nil{fmt.Println(ColorRed,"safe span->",math.Round(math.Floor(lspan/25.0)*25.0),ColorReset)}
	}
}
