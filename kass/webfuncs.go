package barf


import (
	"fmt"
)


func ShowSecs()(rezstring string){
	styps := []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22}
	dims := [][]float64{
		{300},
		{230,300},
		{400},
		{300,400},
		{300,300,200,200},
		{500,300},
		{600,460,230,75},
		{400,690,200,130},
		{400,690,200,130},
		{400,690,200,130},
		{600,490,125},
		{400,600,100,150},
		{390,490,75,125},
		{400,500,125,75},
		{230,460,375,150},
		{230,460,375,150},
		{460},
		{375},
		{230},
		{300,840,200,140,100},
		{400,500,305.8},
		{230},
		{690,600,375,230,125},
	}
	for i := range styps{
		//if i > 5{break}
		s := SecGen(styps[i], dims[i])
		rezstring += fmt.Sprintf("%v %s section\n",s.Styp, SectionMap[s.Styp])
		rezstring += fmt.Sprintf("from secgen funcs (formulae) area %.f xc %.f yc %.f ixx %.f iyy %.f rxx %.f ryy %.f jz %.f\n",s.Prop.Area, s.Prop.Xc, s.Prop.Yc, s.Prop.Ixx, s.Prop.Iyy, s.Prop.Rxx, s.Prop.Ryy, s.Prop.J)
		s.Draw("mono")
		rezstring += s.Txtplot
	}
	return
}
