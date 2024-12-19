package barf

import(
	"fmt"
	"testing"
)

//shah chapter 5 tests 

func TestBmSecAzIs(t *testing.T) {
	b := &RccBm{
		Fck:18.0,
		Fy:290.0,
		Bf:230.0,
		Df:0.0,
		Bw:230.0,
		Dused:300.0,
		Ast:339.0,
		Asc:0.0,
		L0:3000.0,
		Tyb:0.0,
		Cvrt:31.0,
		Cvrc:0.0,
	}
	var rezstring string
	var mur, astmax float64
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("shah 5.3-1 Mur %0.5f kN.m\n",mur)

	b = &RccBm{
		Fck:20.0,
		Fy:500.0,
		Bf:200.0,
		Df:0.0,
		Bw:200.0,
		Dused:450.0,
		Ast:0.0,
		Asc:0.0,
		L0:0.0,
		Tyb:0.0,
		Cvrt:50.0,
		Cvrc:0.0,
		DM:0.25,
	}
	mur, astmax = BmSecAzIs(b)
	rezstring += fmt.Sprintf("shah 5.3-2 Mur %0.5f kN.m Ast max %0.5f mm2\n",mur, astmax)

	b = &RccBm{
		Fck:20.0,
		Fy:500.0,
		Bf:200.0,
		Df:0.0,
		Bw:200.0,
		Dused:450.0,
		Ast:0.0,
		Asc:0.0,
		L0:0.0,
		Tyb:0.0,
		Cvrt:50.0,
		Cvrc:0.0,
		DM:0.10,
	}
	mur, astmax = BmSecAzIs(b)
	rezstring += fmt.Sprintf("shah 5.3-3 Mur %0.5f kN.m Ast max %0.5f mm2\n",mur, astmax)
	
	b = &RccBm{
		Fck:15.0,
		Fy:415.0,
		Bf:230.0,
		Df:0.0,
		Bw:230.0,
		Dused:530.0,
		Ast:1570.0,
		Asc:1005.0,
		L0:0.0,
		Tyb:0.0,
		Cvrt:55.0,
		Cvrc:55.0,
		DM:0.0,
	}
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("shah 5.3-4 Mur %0.5f kN.m\n",mur)
	
	b = &RccBm{
		Fck:20.0,
		Fy:250.0,
		Bf:1000.0,
		Df:100.0,
		Bw:250.0,
		Dused:600.0,
		Ast:2513.0,
		Asc:0.0,
		L0:0.0,
		Tyb:1.0,
		Cvrt:50.0,
		Cvrc:0.0,
		DM:0.0,
	}
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("shah 5.3-5 Mur %0.5f kN.m\n",mur)
	
	b = &RccBm{
		Fck:15.0,
		Fy:415.0,
		Bf:1500.0,
		Df:100.0,
		Bw:300.0,
		Dused:600.0,
		Ast:2454.0,
		Asc:0.0,
		L0:0.0,
		Tyb:1.0,
		Cvrt:40.0,
		Cvrc:0.0,
		DM:0.0,
	}
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("shah 5.3-6 Mur %0.5f kN.m\n",mur)

	b = &RccBm{
		Fck:15.0,
		Fy:415.0,
		Bf:1500.0,
		Df:100.0,
		Bw:300.0,
		Dused:600.0,
		Ast:3373.0,
		Asc:0.0,
		L0:0.0,
		Tyb:1.0,
		Cvrt:65.0,
		Cvrc:0.0,
		DM:0.0,
	}
	
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("shah 5.3-7 Mur %0.5f kN.m\n",mur)

	b = &RccBm{
		Fck:15.0,
		Fy:415.0,
		Bf:1500.0,
		Df:100.0,
		Bw:300.0,
		Dused:550.0,
		Ast:2945.0,
		Asc:615.0,
		L0:0.0,
		Tyb:1.0,
		Cvrt:50.0,
		Cvrc:50.0,
		DM:0.0,
	}
	
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("shah 5.3-8 Mur %0.5f kN.m\n",mur)


	b = &RccBm{
		Fck:15.0,
		Fy:415.0,
		Bf:1500.0,
		Df:100.0,
		Bw:300.0,
		Dused:550.0,
		Ast:3927.0,
		Asc:509.0,
		L0:0.0,
		Tyb:1.0,
		Cvrt:50.0,
		Cvrc:50.0,
		DM:0.0,
	}
	
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("shah 5.3-9 Mur %0.5f kN.m\n",mur)
	
	wantstring := `shah 5.3-1 Mur 20.96183 kN.m
shah 5.3-2 Mur 69.09020 kN.m Ast max 464.96800 mm2
shah 5.3-3 Mur 85.38196 kN.m Ast max 605.90731 mm2
shah 5.3-4 Mur 233.69159 kN.m
shah 5.3-5 Mur 283.27147 kN.m
shah 5.3-6 Mur 456.22813 kN.m
shah 5.3-7 Mur 567.09047 kN.m
shah 5.3-8 Mur 483.55642 kN.m
shah 5.3-9 Mur 597.66994 kN.m
`
	if rezstring != wantstring {
		t.Errorf("beam section analysis (is code) failed")		
		fmt.Println(rezstring)
	}	
}

func TestBmDIs(t *testing.T) {
	var rezstring string
	
	b := &RccBm{
		Fck:15.0,
		Fy:500.0,
		Bf:0.0,
		Df:0.0,
		Bw:250.0,
		Dused:600.0,
		Ast:0.0,
		Asc:0.0,
		L0:0.0,
		Tyb:0.0,
		Cvrt:35.0,
		Cvrc:0.0,
		Lbd:12.0,
		Lspan:6000.0,
	}
	mur := 142.59
	_,_ = BmDIs(b,mur)
	rezstring += fmt.Sprintf("shah 5.5.1 section %0.1f x dused %.0f mm ast %.4f mm2 asc %.4f mm2\n", b.Bw, b.Dused, b.Ast, b.Asc)
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("section analysis mur %0.5f kN.m\n",mur)

	b = &RccBm{
		Fck:15.0,
		Fy:500.0,
		Bf:0.0,
		Df:100.0,
		Bw:250.0,
		Dused:600.0,
		Ast:0.0,
		Asc:0.0,
		L0:0.0,
		Tyb:1.0,
		Cvrt:35.0,
		Cvrc:0.0,
		Lbd:12.0,
		Lspan:6000.0,
	}
	mur = 142.59
	_,_  = BmDIs(b,mur)
	rezstring += fmt.Sprintf("shah 5.5.2 section %0.1f x dused %.0f mm ast %.4f mm2 asc %.4f mm2\n", b.Bw, b.Dused, b.Ast, b.Asc)
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("section analysis mur %0.5f kN.m\n",mur)
	
	b = &RccBm{
		Fck:15.0,
		Fy:415.0,
		Bf:0.0,
		Df:0.0,
		Bw:230.0,
		Dused:600.0,
		Ast:0.0,
		Asc:0.0,
		L0:0.0,
		Tyb:0.0,
		Cvrt:35.0,
		Cvrc:31.0,
		Lbd:12.0,
		Lspan:6000.0,
		Endc:1,
	}
	mur = 225.79
	_,_ = BmDIs(b,mur)
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("shah 5.5.3 section %0.1f x dused %.0f mm ast %.4f mm2 asc %.4f mm2\n", b.Bw, b.Dused, b.Ast, b.Asc)
	rezstring += fmt.Sprintf("section analysis mur %0.5f kN.m\n",mur)
	
	b = &RccBm{
		Fck:20.0,
		Fy:415.0,
		Bf:1050.0,
		Df:120.0,
		Bw:300.0,
		Dused:680.0,
		Ast:0.0,
		Asc:0.0,
		L0:0.0,
		Tyb:0.5,
		Cvrt:65.0,
		Cvrc:35.0,
		Lbd:12.0,
		Lspan:0.0,
		Endc:1,
	}
	mur = 630.0
	_,_ = BmDIs(b,mur)
	rezstring += fmt.Sprintf("shah 5.5.4-a section %0.1f x dused %.0f mm ast %.4f mm2 asc %.4f mm2\n", b.Bw, b.Dused, b.Ast, b.Asc)
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("section analysis mur %0.5f kN.m\n",mur)


	b = &RccBm{
		Fck:20.0,
		Fy:415.0,
		Bf:1050.0,
		Df:120.0,
		Bw:300.0,
		Dused:680.0,
		Ast:0.0,
		Asc:0.0,
		L0:0.0,
		Tyb:0.5,
		Cvrt:65.0,
		Cvrc:35.0,
		Lbd:12.0,
		Lspan:0.0,
		Endc:1,
	}
	mur = 752.0
	_,_ = BmDIs(b,mur)
	rezstring += fmt.Sprintf("shah 5.5.4-b section %0.1f x dused %.0f mm ast %.4f mm2 asc %.4f mm2\n", b.Bw, b.Dused, b.Ast, b.Asc)
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("section analysis mur %0.5f kN.m\n",mur)


	b = &RccBm{
		Fck:20.0,
		Fy:415.0,
		Bf:1050.0,
		Df:120.0,
		Bw:300.0,
		Dused:680.0,
		Ast:0.0,
		Asc:0.0,
		L0:0.0,
		Tyb:0.5,
		Cvrt:65.0,
		Cvrc:35.0,
		Lbd:12.0,
		Lspan:0.0,
		Endc:1,
	}
	mur = 850.0
	_,_ = BmDIs(b,mur)
	rezstring += fmt.Sprintf("shah 5.5.4-c section %0.1f x dused %.0f mm ast %.4f mm2 asc %.4f mm2\n", b.Bw, b.Dused, b.Ast, b.Asc)
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("section analysis mur %0.5f kN.m\n",mur)


	b = &RccBm{
		Fck:15.0,
		Fy:250.0,
		Bf:0.0,
		Df:0.0,
		Bw:230.0,
		Dused:400.0,
		Ast:0.0,
		Asc:0.0,
		L0:0.0,
		Tyb:0.0,
		Cvrt:33.0,
		Cvrc:0.0,
		Lbd:0.0,
		Lspan:0.0,
		Endc:1,
		DM:0.3,
	}
	mur = 43.4
	_,_ = BmDIs(b,mur)
	rezstring += fmt.Sprintf("shah 5.5.5 section %0.1f x dused %.0f mm ast %.4f mm2 asc %.4f mm2\n", b.Bw, b.Dused, b.Ast, b.Asc)
	mur, _ = BmSecAzIs(b)
	rezstring += fmt.Sprintf("section analysis mur %0.5f kN.m\n",mur)



	
	wantstring := `shah 5.5.1 section 250.0 x dused 600 mm ast 694.1714 mm2 asc 0.0000 mm2
section analysis mur 142.53280 kN.m
shah 5.5.2 section 250.0 x dused 600 mm ast 591.6168 mm2 asc 0.0000 mm2
section analysis mur 142.58439 kN.m
shah 5.5.3 section 230.0 x dused 600 mm ast 1314.5097 mm2 asc 394.1266 mm2
section analysis mur 226.29600 kN.m
shah 5.5.4-a section 300.0 x dused 680 mm ast 3168.8109 mm2 asc 0.0000 mm2
section analysis mur 630.00000 kN.m
shah 5.5.4-b section 300.0 x dused 680 mm ast 3935.1612 mm2 asc 0.0000 mm2
section analysis mur 752.00000 kN.m
shah 5.5.4-c section 300.0 x dused 680 mm ast 4424.0748 mm2 asc 449.5600 mm2
section analysis mur 850.83552 kN.m
shah 5.5.5 section 230.0 x dused 400 mm ast 619.8379 mm2 asc 0.0000 mm2
section analysis mur 43.38761 kN.m
`
	if rezstring != wantstring {
		fmt.Println(rezstring)
		t.Errorf("beam rebar design (is code) failed")
	}	
}

