package barf

import (
	"fmt"
	"testing"
)

func TestPlateISec(t *testing.T){
	var styp int
	var dims []float64
	styp = 29
	dims = []float64{200,250,25,25,300,20}
	s := SecGen(styp, dims)
	s.Draw("dumb")
	rezstr := fmt.Sprintf("%s\nplate i section prop \narea %f \nixx %f iyy %f \nrxx %f ryy %f \ncx %f cy %f\n",s.Txtplot,s.Prop.Area, s.Prop.Ixx, s.Prop.Iyy, s.Prop.Rxx, s.Prop.Ryy, s.Prop.Xc, s.Prop.Yc)
	t.Log(rezstr)
}

func TestTubeSec(t *testing.T){
	var styp int
	var dims []float64
	styp = 5
	dims = []float64{21.3,21.3 - 2.0*2.3}
	s := SecGen(styp, dims)
	s.Draw("dumb")
	rezstr := fmt.Sprintf("%s\npipe section prop \narea %f \nixx %f iyy %f \nrxx %f ryy %f \ncx %f cy %f\n",s.Txtplot,s.Prop.Area, s.Prop.Ixx, s.Prop.Iyy, s.Prop.Rxx, s.Prop.Ryy, s.Prop.Xc, s.Prop.Yc)
	t.Log(rezstr)
	bar := CalcSecProp(styp, dims)
	t.Log(bar.Area)
}

func TestVecAng(t *testing.T){
	v1 := []float64{0,0}
	
	v2 := []float64{0,0}
	
	v3 := []float64{0,0}
	
	v4 := []float64{0,0}
	ang := VecAng(v1,v2,v3,v4)

	t.Log("angle-",ang*180.0/3.14)
}

func TestSecQ(t *testing.T){
	var rezstring string
	styps := []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,26}
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
		{300,500,300},
	}
	for i := range styps{
		if i > 1 && i != 12{continue}
		s := SecGen(styps[i], dims[i])
		rezstring += fmt.Sprintf("%v %s section\n",s.Styp, SectionMap[s.Styp])
		s.CalcQ()
		rezstring += fmt.Sprintf("qxx %.f  ixx %.f vfacx %f\n",s.Prop.Qxx, s.Prop.Ixx, s.Prop.Vfx)
		switch i{
			case 0:
			rezstring += fmt.Sprintf("Aq/ID = %.3f\n",s.Prop.Area * s.Prop.Qxx/s.Prop.Ixx/s.Prop.Dims[0])
			case 1:
			rezstring += fmt.Sprintf("dq/I = %.3f\n",dims[i][1] * s.Prop.Qxx/s.Prop.Ixx)
			case 12:
			hw := dims[i][1]
			rezstring += fmt.Sprintf("hq/I = %.3f\n",hw * s.Prop.Qxx/s.Prop.Ixx)
		}
	}
	wantstring := ``
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("sections shear factor test failed")
	}
}

func TestSecOffset(t *testing.T){
	s := SectIn{
		Coords:[][]float64{
			{0,0},
			{100,0},
			{60,50},
			{40,50},
			{0,0},
		},
		Ncs:[]int{5},
		Wts:[]float64{1.0},
		Styp:-1,
	}
	s.Draw("dumb")
	fmt.Println(s.Txtplot)
	s1 := SecOffset(s, 20.0, 1.0)
	s1.Draw("dumb")
	fmt.Println(s1.Txtplot)
	fmt.Println(s1.Coords)
}

func TestPerpVec(t *testing.T){
	v1 := []float64{10,10}
	v2 := []float64{20,20}
	v3 := Rotvec(45.0, v1, v2)
	v4 := Rotvec(90.0, v1, v2)
	v5 := Rotvec(120.0, v1, v2)
	
	fmt.Println(v3,v4,v5)
}

func TestMoment1(t *testing.T){
	styp := 1
	dims := []float64{50,100}
	s := SecGen(styp, dims)
	fmt.Println("qy->",s.Prop.Xc*s.Prop.Area*s.Prop.Area/s.Prop.Ixx/dims[1], "qx->",s.Prop.Yc*s.Prop.Area*s.Prop.Area/s.Prop.Iyy/dims[0])
	styp = 0
	dims = []float64{100}
	s = SecGen(styp, dims)
	fmt.Println("qy->",s.Prop.Xc*s.Prop.Area*s.Prop.Area/s.Prop.Ixx/dims[0]/2.0, "qx->",s.Prop.Yc*s.Prop.Area)
}

func TestSplitmax(t *testing.T){
	styp := 7
	dims := []float64{1000,1000,250,250}
	s := SecGen(styp, dims)
	s1 := s.Splitmax(100.0)
	for i, pt := range s1.Coords{
		fmt.Println("idx-",i, "pt-",pt)
	}
	//s1.Draw("qt")
	//fmt.Println(s1.Txtplot)
	//fmt.Println(s1.Wts, s1.Ncs, s1.Coords)
}

func TestSplitSides(t *testing.T){
	styp := 2
	dims := []float64{300.0}
	s := SecGen(styp, dims)
	s1 := s.SplitSides(0.0)
	//fmt.Println(s1.Wts, s1.Ncs, s1.Coords)
	s1.Draw("dumb")
	fmt.Println(s1.Txtplot)
	fmt.Println(s1.Wts, s1.Ncs, s1.Coords)
}

func TestSecWidth(t *testing.T){
	styp := 2
	dims := []float64{300.0}
	dy := 50.0
	s := SecGen(styp, dims)
	n, dx, pts := s.GetWidth(dy)
	s.Draw("dumb")
	fmt.Println(s.Txtplot)

	fmt.Println("triangular section")
	fmt.Println("no. of points->",n, "is div?",n%2==0)
	fmt.Println("points->",pts)
	fmt.Println("dxs->",dx)
	
	styp = 13
	dims = []float64{300.0,400.0,125.0,100.0}
	dy = 125.0
	s = SecGen(styp, dims)
	n, dx, pts = s.GetWidth(dy)

	s.Draw("dumb")
	fmt.Println(s.Txtplot)

	fmt.Println("c section")
	fmt.Println("no. of points->",n, "is div?",n%2==0)
	fmt.Println("points->",pts)
	fmt.Println("dxs->",dx)

	s = SecRotate(s, 90.0)
	s.UpdateProp()
	n, dx, pts = s.GetWidth(dy)

	s.Draw("dumb")
	fmt.Println(s.Txtplot)

	fmt.Println("rotated c section")
	fmt.Println("no. of points->",n, "is div?",n%2==0)
	fmt.Println("points->",pts)
	fmt.Println("dxs->",dx)
}

func TestSecRotate(t *testing.T){
	styps := []int{2}
	dims := [][]float64{
		{300},
	}
	angs := []float64{-45.0}
	scls := []float64{0.5}
	for i := range styps{
		s := SecGen(styps[i], dims[i])
		s.Draw("dumb")
		fmt.Println(s.Txtplot)
		sr := SecRotate(s, angs[i])
		sr.Draw("dumb")
		fmt.Println(sr.Txtplot)
		ss := SecScale(s, scls[i], scls[i])
		ss.Draw("dumb")
		fmt.Println(ss.Txtplot)
	}
}

func TestSecFlip(t *testing.T){
	styps := []int{1,2,3,4}
	dims := [][]float64{
		{200,300},
		{300},
		{200,300},
		{300,400,200,300},
	}
	t.Log("behold yeay flippening")
	for i := range styps{
		s := SecGen(styps[i], dims[i])
		s.Draw("dumb")
		fmt.Println(s.Txtplot)
		sf := FlipX(s)
		sf.Draw("dumb")
		fmt.Println(sf.Txtplot)
	}
}

func TestHaunchSec(t *testing.T){
	var rezstring string
	styps := []int{23,24}
	dims := [][]float64{
		{200,300,20,10,100},
		{200,300,20,10,100},
	}
	for i, styp := range styps{
		s := SecGen(styp, dims[i])
		rezstring += fmt.Sprintf("section %s area %.f ixx %.f\n",SectionMap[styp],s.Prop.Area,s.Prop.Ixx)
	}
	wantstring := ``
	fmt.Println(rezstring)
	if rezstring != wantstring{
		t.Errorf("haunch section generation test failed")
	}
}

func TestSecGen(t *testing.T) {
	var rezstring string
	styps := []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,26}
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
		{300,500,300},
	}
	for i := range styps{
		if i != len(styps)-1 {continue}
		s := SecGen(styps[i], dims[i])
		rezstring += fmt.Sprintf("%v %s section\n",s.Styp, SectionMap[s.Styp])
		rezstring += fmt.Sprintf("area %.f xc %.f yc %.f ixx %.f iyy %.f rxx %.f ryy %.f jz %.f\n",s.Prop.Area, s.Prop.Xc, s.Prop.Yc, s.Prop.Ixx, s.Prop.Iyy, s.Prop.Rxx, s.Prop.Ryy, s.Prop.J)
		s.Draw("dumb")
		fmt.Println(s.Txtplot)
		s.SecInit()
		area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle := SecArea(&s, true)
		rezstring += fmt.Sprintf("area %.f xc %.f yc %.f ixx %.f iyy %.f ixy %.f iuu %.f ivv %.f pxangle %.f\n",area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle)
	}
	wantstring := ``
	if rezstring != wantstring{
		fmt.Println(rezstring)
		t.Errorf("sections generation test (SectIn) failed")
	}
}

func TestSecInit(t *testing.T) {
	var rezstring string
	//SecArea(sec *SectIn, allprp bool) (area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle float64)

	ncs := []int{
		9,5,
	}
	wts := []float64{
		1.0,-1.0,
	}
	coords := [][]float64{
		{100,100},
		{200,100},
		{200,250},
		{700,250},
		{700,550},
		{850,550},
		{850,650},
		{100,650},
		{100,100},
		{200,300},
		{650,300},
		{650,550},
		{200,550},
		{200,300},
	}
	sec := &SectIn{
		Ncs:ncs,
		Wts:wts,
		Coords:coords,
	}
	sec.SecInit()
	area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle := SecArea(sec, true)
	rezstring += fmt.Sprintf("area %.3f, xc %.3f, yc %.3f, ixx %.3f, iyy %.3f, ixy %.3f, iuu %.3f, ivv %.3f, pxangle %.3f",area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle)

	wantstring := `area 157500.000, xc 394.048, yc 455.952, ixx 4050669642.857, iyy 8313169642.857, ixy 1950892857.143, iuu 9071246468.011, ivv 3292592817.704, pxangle -68.765`
	t.Log(rezstring)
	if rezstring != wantstring {
		t.Errorf("section properties calculation failed")
		fmt.Println(rezstring)
	}


}

func TestSecPrp(t *testing.T) {
	var rezstring string
	ncs := []int{
		9,5,
	}
	wts := []float64{
		1.0,-1.0,
	}
	coords := [][]float64{
		{100,100},
		{200,100},
		{200,250},
		{700,250},
		{700,550},
		{850,550},
		{850,650},
		{100,650},
		{100,100},
		{200,300},
		{650,300},
		{650,550},
		{200,550},
		{200,300},
	}
	area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle := SecPrp(ncs, wts, coords)
	rezstring += fmt.Sprintf("area %.3f, xc %.3f, yc %.3f, ixx %.3f, iyy %.3f, ixy %.3f, iuu %.3f, ivv %.3f, pxangle %.3f",area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle)

	wantstring := `area 157500.000, xc 394.048, yc 455.952, ixx 4050669642.857, iyy 8313169642.857, ixy 1950892857.143, iuu 9071246468.011, ivv 3292592817.704, pxangle -68.765`
	t.Log(rezstring)
	if rezstring != wantstring {
		t.Errorf("section properties calculation failed")
		fmt.Println(rezstring)
	}
}

func TestSecPropRccx(t *testing.T){
	//https://www.structx.com/
	//test secprop calcs for rcc sections
	//stypes - 0 (circle), 1 (rect), 6 (tee),  8 (l left), 9 (l right)
	//ADD - pentagon
	var rezstring string
	sts := []int{0,1,6,8,9}
	ds := [][]float64{
		{300},
		{230,360},
		{1000,690,230,130},
		{690,460,150,150},
		{600,360,75,75},
	}
	names := []string{"circle","rectangle","T","L-left","L-right"}
	for i, st := range sts{
		d := ds[i]
		bar := CalcSecProp(st, d)
		rezstring += fmt.Sprintf("%s section, dims - %v\n",names[i],d)
		rezstring += fmt.Sprintf("area %.1f xc %.1f yc %.1f ixx %.1f iyy %.1f J %.1f\n",
			bar.Area, bar.Xc, bar.Yc, bar.Ixx, bar.Iyy, bar.J)
	}
	wantstring := `circle section, dims - [300]
area 70685.8 xc 150.0 yc 150.0 ixx 397607820.2 iyy 397607820.2 J 795215640.4
rectangle section, dims - [230 360]
area 82800.0 xc 115.0 yc 180.0 ixx 894240000.0 iyy 365010000.0 J 880533159.3
T section, dims - [1000 690 230 130]
area 258800.0 xc 500.0 yc 453.3 ixx 11249808598.7 iyy 11401126666.7 J 468715230.0
L-left section, dims - [690 460 150 150]
area 150000.0 xc 428.7 yc 313.7 ixx 2263746500.0 iyy 6532546500.0 J 918528576.4
L-right section, dims - [600 360 75 75]
area 66375.0 xc 215.5 yc 264.5 ixx 635301205.0 iyy 2358576205.0 J 111376942.8
`
	t.Log(rezstring)
	if rezstring != wantstring{
		t.Errorf("rcc section prop test failed")
		fmt.Println(rezstring)
	}
}


func TestSecArea(t *testing.T) {
	var rezstring string
	ncs := []int{
		9,5,
	}
	wts := []float64{
		1.0,-1.0,
	}
	coords := [][]float64{
		{100,100},
		{200,100},
		{200,250},
		{700,250},
		{700,550},
		{850,550},
		{850,650},
		{100,650},
		{100,100},
		{200,300},
		{650,300},
		{650,550},
		{200,550},
		{200,300},
	}
	s := SectIn{Ncs:ncs, Wts:wts, Coords:coords}
	area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle := SecArea(&s, true)
	rezstring += fmt.Sprintf("area %.3f, xc %.3f, yc %.3f, ixx %.3f, iyy %.3f, ixy %.3f, iuu %.3f, ivv %.3f, pxangle %.3f",area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle)

	wantstring := `area 157500.000, xc 394.048, yc 455.952, ixx 4050669642.857, iyy 8313169642.857, ixy 1950892857.143, iuu 9071246468.011, ivv 3292592817.704, pxangle -68.765`
	t.Log(rezstring)
	if rezstring != wantstring {
		t.Errorf("section properties calculation failed")
		fmt.Println(rezstring)
	}
}

// func TestDraw3d(t *testing.T){	
// 	styps := []int{1,1}
// 	dims := [][]float64{{300,400},{300.0,300.0}}
// 	coords := [][]float64{{0,0,0},{0,0,3000},{2500,0,3000},{2500,0,0}}
// 	ms := [][]int{{1,2,1,1},{2,3,2,2},{4,3,1,1}}
// 	var ss []SectIn
// 	for i, styp := range styps{
// 		dim := dims[i]
// 		s := SecGen(styp, dim)
// 		ss = append(ss, s)
// 	}
// 	data := Draw3d(ms, coords, ss) 
// 	//data := s.Dat3d(pl, p0, lz)
// 	pltskript := "plotsec3d.gp"; term := "qt"; fname := ""; title := "bhak bc"
// 	_ = skriptrun(data, pltskript, term, title, "",fname)
// }


/*
func TestMove3d(t *testing.T){
	//given center points of section coords and []sections
	styp := 6
	dims := []float64{300.0,400.0,125.0,100.0}
	s := SecGen(styp, dims)
	
	secvec := []*SectIn{s1, s2}
}
*/

/*

	/*
	ncs = []int{
		5,
	}
	wts = []float64{
		1.0,
	}
	coords = [][]float64{
		{10,10},
		{20,10},
		{20,20},
		{10,20},
		{10,10},
	}
	area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle = SecPrp(ncs, wts, coords)
	rezstring += fmt.Sprintf("area %.3f, xc %.3f, yc %.3f, ixx %.3f, iyy %.3f, ixy %.3f, iuu %.3f, ivv %.3f, pxangle %.3f\n",area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle)

	ncs = []int{
		4,
	}
	wts = []float64{
		1.0,
	}
	coords = [][]float64{
		{0,0},
		{-10,-10},
		{10,-10},
		{0,0},
	}
	area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle = SecPrp(ncs, wts, coords)
	rezstring += fmt.Sprintf("area %.3f, xc %.3f, yc %.3f, ixx %.3f, iyy %.3f, ixy %.3f, iuu %.3f, ivv %.3f, pxangle %.3f\n",area, xc, yc, ixx, iyy, ixy, iuu, ivv, pxangle)
*/

