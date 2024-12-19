package barf

import (
	//"fmt"
	"testing"
)

func TestDirvec(t *testing.T){
	faces := []string{"n","e","w","s"}
	for i, face := range faces{
		t.Log("checking",i,face)
		dirs := getdirvec(face)
		t.Log("dirs-",dirs)
	}
}

func TestSquarifyBasic(t *testing.T){
	t.Log("starting basic squarify example->")
	f := Flr{Origin: Pt2d{X: 0, Y: 0}, End: Pt2d{X: 6, Y: 4}, Units:"m"}
	f.Flrarea()
	f.Name = "base"
	//r := []float64{6,6,4,3,2,2,1}
	r := []float64{6, 6, 4, 3, 2, 2, 1}
	labels := []string{"r1-6","r2-6","r3-4","r4-3","r5-2","r6-2","r7-1"}
	r = Scalerooms(&f,r,false)
	FlrPln(&f, r, labels)
	//f.Flrprint(false)
	GPlotFloors(&f, true)	
	t.Log("now with unit conversion->")
}


func TestFlrCon(t *testing.T){
	f := Flr{Origin: Pt2d{X: 0, Y: 0}, End: Pt2d{X: 31, Y: 19},
		Isroot:true,
		Areas:[]float64{140, 120, 120, 75, 50, 50, 30},
		Labels:[]string{"living","bed1","bed2","kitchen","bath1","bath2","deck"}}
	f.FlrLay()
}

func TestFlrLay(t *testing.T){
	t.Log("starting marson/musse ex. 10a")
	var areas []float64
	var rooms []string
	areas = []float64{
		9e6,8e6,8e6,7e6,6e6,4e6,
	}
	rooms = []string{
		"living","kitchen","bed_m","bed","utility","bath",
	}
	f := Flr{
		Tomm:true,
		Width:6000,
		Height:7000,
		Units:"mm",
		Origin:Pt2d{0,0},
		End:Pt2d{6000,7000},
		Areas:areas,
		Labels:rooms,
		Verbose:true,
		Round:true,
		Term:"dumb",	
	}
	_ = f.FlrLay()

	
	t.Log("starting marson/musse ex. 12a")
	areas = []float64{
		22e6,14e6,14e6,12e6,12e6,10e6,
	}
	rooms = []string{
		"living","bed","bed","office","kitchen","bath",
	}
	f = Flr{
		Tomm:true,
		Width:7000,
		Height:12000,
		Cwidth:900,
		Units:"mm",
		Origin:Pt2d{0,0},
		End:Pt2d{7000,12000},
		Areas:areas,
		Labels:rooms,
		Verbose:true,
		Round:true,
		Term:"dumb",	
	}
	_ = f.FlrLay()

}
