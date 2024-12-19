package barf

import (
	"fmt"
	"testing"
	"math"
	kass"barf/kass"
)

func TestPrlnDz(t *testing.T){
	
}

func TestBmChk(t *testing.T){
	var rezstring string
	//abel rect bm ex
	t.Log("abel ex. 8.1 rect beam")
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
		Dims:[]float64{75,150},
		DL:0.0,
		LL:0.5,
		Lbl:25.0,
		Rbl:25.0,
		Selfwt:false,
		Spam:true,
		Code:1,
	}
	ok, val, _ := BmChk(b)
	rezstring += fmt.Sprintf("abel ex. 8.1 rect. beam\n")

	rezstring += fmt.Sprintf("ok - %v val - %v\n",ok, val)
	_ = BmDesign(b)
	
	rezstring += fmt.Sprintf("from dz rez - %v\n",b.Rez)
	
	t.Log("ramc ex 13.12 rect deodar beam")
	b = &WdBm{
		Styp:1,
		Lspan:5300.0,	
		Grp:2,
		Endc:1,
		Dims:[]float64{240,400},
		DL:16.0,
		LL:0.0,
		Lbl:150.0,
		Rbl:150.0,
		Selfwt:false,
		Spam:true,
		Code:1,
		Brchk:false,
	}
	ok, val, _ = BmChk(b)
	rezstring += fmt.Sprintf("sp.33 ex. 1 rect. beam\n")
	t.Log(ColorRed,"rez - ok -> ",ok,"\n val - > ",val,ColorReset)
	t.Log(ColorYellow,"val - > b.Vu, b.Mu, b.Dmax,dall,fb,fball, fv, b.Prp.Fv, b.Sec.Prop.Area",ColorReset)
	// rezstring += fmt.Sprintf("perm stress - %.2f N/mm2 actual - %.2f N/mm2 ok? - %v\n",val[0], val[2], ok)

	
	_ = BmDesign(b)
	
	rezstring += fmt.Sprintf("from dz rez - %v\n",b.Rez)
	
	t.Fatal(rezstring)

	t.Log("ramc ex 13.13 built-up rect sal beam")
	b = &WdBm{
		Styp:4,
		Lspan:8000.0,	
		Grp:1,
		Endc:1,
		Dims:[]float64{200,300,50,50},
		Tplnk:50,
		DL:6.25,
		LL:0.0,
		Selfwt:false,
		Spam:true,
		Solid:true,
		Code:1,
		Brchk:false,
	}
	ok, val, _ = BmChk(b)
	rezstring += fmt.Sprintf("ramc ex. 13.3 built up solid rect. beam\n")
	t.Log(ColorRed,"rez - ok -> ",ok,"\n val - > ",val,ColorReset)
	t.Log(ColorYellow,"val - > b.Vu, b.Mu, b.Dmax,dall,fb,fball, fv, b.Prp.Fv, b.Sec.Prop.Area",ColorReset)

	t.Log("ramc ex 13.14 notched rect beam with point load")
	b = &WdBm{
		Styp:1,
		Lspan:5300.0,	
		Grp:1,
		Endc:1,
		Dims:[]float64{150,300},
		Selfwt:true,
		Spam:true,
		Notch:true,
		Code:1,
		Dtyp:1,
		Brchk:false,
		Ldcases:[][]float64{{1,1,25000,0,2650,0,1}},
	}
	ok, val, _ = BmChk(b)
	rezstring += fmt.Sprintf("ramc ex. 13.3 built up solid rect. beam\n")
	t.Log(ColorRed,"rez - ok -> ",ok,"\n val - > ",val,ColorReset)
	t.Log(ColorYellow,"val - > b.Vu, b.Mu, b.Dmax,dall,fb,fball, fv, b.Prp.Fv, b.Sec.Prop.Area",ColorReset)
	
	wantstring := ``
	if rezstring != wantstring{
		t.Errorf("timber beam check test failed")
	}
}

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

// func TestBmDzCs(t *testing.T){
// 	b := &WdBm{
// 		Styp:1,
// 		Lspan:3500.0,	
// 		Prp:kass.Wdprp{
// 			Em:9500,
// 			Fv:1.4,
// 			Fcp:2.8,
// 			Fcb:11.2,
// 			Ft:11.2,
// 			Pg:0.5,
// 		},
// 		Endc:1,
// 		DL:0.0,
// 		LL:0.5,
// 		Lbl:25.0,
// 		Rbl:25.0,
// 		Selfwt:true,
// 		Spam:true,
// 		Nsecs:3,
// 	}
// 	err := BmDzCs(b)
// 	if err != nil{
// 		fmt.Println(err)
// 	}
// }

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
