package barf

//yield line design of slabs

import (
	"fmt"
	"math"
)

//Yslb stores struct fields for yield line analysis
//as in bauer_redwood_1987_yield_line_c&s.pdf
type Yslb struct{
	Np, Nedge, Nreg int
	Coords [][]float64
	Rmyx float64
	Cs, Vcs []float64
}


func SlbYld(s *RccSlb, dused float64){
	switch s.Endc{	
		case 10:
		//simply supported on 4 sides
		SlbYld10(s, dused)
		case 11:
		//simply supported on 3 sides, free on the fourth
		SlbYld11(s, dused)
		case 12:
		//simply supported on 2 sides
	}
}

func (s *RccSlb) FitI(){
	//based on end conditions figures out ideal i1, i2, i3, i4 vals
	//uses adam-hytham formulae
	r := s.Ly/s.Lx
	s.I1 = 0.21 + 1.68 * r -0.48 * math.Pow(r,2.0)
	s.I3 = s.I1
	s.I2 = -2.74 + 5.41 * r - 1.09 * math.Pow(r,2.0)
	s.I4 = s.I2
	s.Rmyx = -1.53 + 3.0 * r - 0.47 * math.Pow(r,2.0)
	s.Rmyx = 1/s.Rmyx
}

//SlbYldRect uses general formulae for rect supported slabs with udl/line loads
//if dused > 0.0 {add dused * pgck to udl load}
func SlbYldRect(s *RccSlb, dused float64) (mcxsup, mcxmid, mcysup, mcymid float64){
	switch {
	case s.Endc <= 10:
		ar := 2.0 * s.Lx/(math.Sqrt(1.0 + s.I2) + math.Sqrt(1.0 + s.I4))/1e3
		br := 2.0 * s.Ly/(math.Sqrt(1.0 + s.I1) + math.Sqrt(1.0 + s.I3))/1e3
		var n, mur, al, be float64
		switch s.Code{
			case 1:
			pgck := 0.025
			n = 1.5 * (pgck * dused + s.DL + s.LL)
			case 2:
			pgck := 0.024
			n = 1.4 * (pgck * dused + s.DL) + 1.6 * s.LL
		}
		//line load factors
		if s.Pa + s.Pb > 0.0 {
			al = s.Pa/n/s.Ly/1e3
			be = s.Pb/n/s.Lx/1e3
		}
		nadj := n * (1.0 + al + 2.0 * be)
		badj := br * math.Sqrt((1.0 + al + 2.0 * be)/(1.0 + 3.0 * be))
		fmt.Println(nadj, badj)
		mur = nadj * ar * badj/8.0/(1.0 + badj/ar + ar/badj)
		fmt.Println("mur->",mur,"kn-m/m")
		s.BM = []float64{}
		s.BM = append(s.BM, mur)
		for _, vi := range []float64{s.I1, s.I2, s.I3, s.I4}{
			s.BM = append(s.BM, vi * mur)
		}
		fmt.Println("bm->",s.BM,"kn-m/m")
	}
	return
}

//SlbYld10 uses yield line forumulae to calc ult. moment of a simply supported two way slab
func SlbYld10(s *RccSlb, dused float64){
	//simply supported (4 sides)
	var wul, ma, bf, lx, ly float64
	switch s.Code{
		case 1:
		wul = 1.5 * s.DL + 1.5 * s.LL + 1.5 * 0.025 * dused 
		case 2:
		wul = 1.4 * s.DL + 1.6 * s.LL + 1.4 * 0.024 * dused 		
	}
	lx = s.Lx/1e3; ly = s.Ly/1e3
	//first yl pattern
	for b := 0.05 * ly; b <= 0.5 * ly; b += 0.05 * ly{
		m1 := wul * math.Pow(lx, 2)*(3.0 * b - 2.0 * math.Pow(b,2))/12.0
		m1 = m1/(s.Rmyx * math.Pow(lx/ly,2) + 2.0 * b)
		if m1 > ma{
			ma = m1; bf = b 
		}
		if m1 < ma{
			break
			//return
		}
	}
	fmt.Println("pattern 1->",ma, bf)
	fmt.Println("ultimate moment in x direction->",ma,"kn.m/m")
	fmt.Println("ultimate moment in y direction->",ma*s.Rmyx, "kn.m/m")
}


//SlbYld11 calcs the ult.moment of a slab simply supported on 3 sides, free on the fourth
//hulse sec 4.4
func SlbYld11(s *RccSlb, dused float64){
	var wul, ma, mb, m1, m2, mul, xp, yp float64
	switch s.Code{
		case 1:
		wul = 1.5 * s.DL + 1.5 * s.LL + 1.5 * 0.025 * dused 
		case 2:
		wul = 1.4 * s.DL + 1.6 * s.LL + 1.4 * 0.024 * dused 		

	}
	lx := s.Lx/1000.0; ly := s.Ly/1000.0
	//first yl pattern
	for x := 0.05 * lx; x <= 0.5 * lx; x += 0.05 * lx{
		m1 = wul * (0.5 * lx * ly - ly * x/3.0)
		m2 = 2.0 * (s.Rmyx * ly/x + x/ly)
		m3 := m1/m2
		if m3 > ma{
			ma = m3; xp = x
		}
		if m3 < ma{
			break
			//return
		}
	}
	fmt.Println("pattern 1->",ma, xp, yp)
	mul = ma
	//xp = 0.5 * s.Lx
	for y := 0.05 * ly; y <= ly; y += 0.05 * ly{
		m1 = wul * (lx * ly/2.0 - lx * y/6.0)
		m2 = 4.0 * s.Rmyx * ly/lx + lx/y
		m3 := m1/m2
		if m3 > mb{
			mb = m3; yp = y
		}
		if m3 < mb{
			break
		}
	}
	fmt.Println("pattern 2->",mb, xp, yp)
	if mul < mb{mul = mb}
	fmt.Println("ultimate moment in x direction->",mul,"kn.m/m")
	fmt.Println("ultimate moment in y direction->",mul*s.Rmyx, "kn.m/m")
}
