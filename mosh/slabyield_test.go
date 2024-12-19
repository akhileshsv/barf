package barf

import (
	"fmt"
	"testing"
	"gonum.org/v1/gonum/stat/combin"
)

func TestSlbEndC2W(t *testing.T){
	s := &RccSlb{
		Fck: 25.0,
		Fy: 415.0,
		Lx: 4.5,
		Ly: 6.0,
	}
	s.Ec = make([]int,4)
	var labels = []string{"cl","cr","ct","cb"}
	n := len(labels)
	for k := 1; k <= n; k++{
		gen := combin.NewCombinationGenerator(n, k)
		idx := 0
		for gen.Next() {
			endc := make([]int, len(labels))
			fmt.Println(idx, gen.Combination(nil))
			for _, val := range gen.Combination(nil){
				endc[val] = 1
			}
			for i, val := range endc{
				fmt.Println(ColorCyan,labels[i], "\u2192",ColorRed,val,ColorReset)
			}
			copy(s.Ec, endc)
			s.Slbc = 1
			s.EndC()
			fmt.Println("end condition\u2192",s.Endc,s.Slbc)
			s.EndC2W()
			fmt.Println("I1 I2 I3 I4 Rmyx\u2192",s.I1, s.I2, s.I3, s.I4, s.Rmyx)
			s.Slbc = 2
			s.EndC()
			fmt.Println("end condition\u2192",s.Endc,s.Slbc)
			s.EndC2W()
			fmt.Println("I1 I2 I3 I4 Rmyx\u2192",s.I1, s.I2, s.I3, s.I4, s.Rmyx)
			idx++
		}
	}
}

func TestSlb2WCoefComp(t *testing.T){
	//compares bs code and is code b.m coefficients for 2 way slabs
	s := &RccSlb{
		Fck: 25.0,
		Fy: 415.0,
		Lx: 4000,
		Type:2,
	}
	for _, r := range []float64{1.0,1.1,1.2,1.3,1.4,1.5,1.75,2.0,2.5,3.0}{
		s.Ly = r * s.Lx
		fmt.Println(ColorPurple,"***/STARTING FOR R->",r,"***")
		for i := 1; i < 11; i++{
			s.Endc = i
			fmt.Println(ColorGreen,"end condition->",s.Endc)
			fmt.Println(ColorRed,"is code->")
			m1, m2, m3, m4 := Slb2BMCoefIs(s.Endc, s.Lx, s.Ly)
			fmt.Printf("mcx sup %.4f mcx mid % .4f mcy sup %.4f mcy mid %.4f\n",m1, m2, m3, m4)
			s.EndC2W()
			fmt.Println(ColorCyan,"bs code->")
			m5, m6, m7, m8 := Slb2BMCoefBs(s.Endc, s.Ns, s.Nl, s.Lx, s.Ly)
			fmt.Printf("mcx sup %.4f mcx mid % .4f mcy sup %.4f mcy mid %.4f\n",m5, m6, m7, m8)
			fmt.Println(ColorYellow,"delta")
			fmt.Printf("mcx sup %.4f mcx mid % .4f mcy sup %.4f mcy mid %.4f\n",m1-m5, m2-m6, m3-m7, m4-m8)
		}

	}
}

func TestSlbYld(t *testing.T){
	//hulse ex. 4.4.1 - three side supported slab
	dused := 0.0
	s := &RccSlb{
		Lx:5000.0,
		Ly:8000.0,
		Type:2,
		Endc:11,
		DL:1.0,
		LL:3.0,
		Rmyx:1.0,
		Code:2,
	}
	SlbYld(s,dused)


	//ss square slab
	s = &RccSlb{
		Lx:5000.0,
		Ly:5000.0,
		Type:2,
		Endc:10,
		DL:1.0,
		LL:3.0,
		Rmyx:1.0,
		Code:2,
	}
	SlbYld(s,dused)
}

func TestSlbYldRect(t *testing.T){
	s := &RccSlb{
		Fck:40.0,
		I1:1,I2:1,I3:0,I4:0,
		Lx:7500.0,
		Ly:9000.0,
		DL:20.0/1.4,
		LL:0.0,
		Code:2,
		Rmyx:1.0,
		Endc:4,
	}
	fmt.Printf("i1 %.2f i2 %.2f i3 %.2f i4 %.2f my/mx %.2f\n",s.I1, s.I2, s.I3, s.I4, s.Rmyx)
	dused := 0.0
	SlbYldRect(s, dused)
	fmt.Println("mur->",s.BM[0])
	s.FitI()
	s.I3 = 0.0; s.I4 = 0.0
	SlbYldRect(s, dused)
	fmt.Printf("i1 %.2f i2 %.2f i3 %.2f i4 %.2f my/mx %.2f\n",s.I1, s.I2, s.I3, s.I4, s.Rmyx)
	//fmt.Println("new ex")
	
}

/*
s = &RccSlb{
		Fck:30.0,
		Fy:460.0,
		Lx:4500.0,
		Ly:6300.0,
		DL:0.0,
		LL:10.0,
		Diamain:10.0,
		Diadist:10.0,
		Code:2,
		Type:2,
		Endc:10,
		Nomcvr:25.0,
	}
	s.EndC2W()
	dused = 220.0
	fmt.Printf("i1 %.2f i2 %.2f i3 %.2f i4 %.2f my/mx %.2f\n",s.I1, s.I2, s.I3, s.I4, s.Rmyx)
	
	SlbYldRect(s, dused)
	SlbYld(s, dused)
*/
