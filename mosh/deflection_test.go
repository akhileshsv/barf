package barf

import (
	"fmt"
	"testing"
)

func TestBmSdrat(t *testing.T){
	SpanDepRatBs(30.0, 460.0, 10.0, 0.0, 300.0, 600.0, 0.0, 600.0-45.0, 2036.0, 2036.0, 0.0, 368.5, 0.0, 1, 2)
}

func TestSlbSdratBs(t *testing.T) {
	//hulse section 7.1 test for span-depth ratio
	var rezstring string
	s := &RccSlb{
		Fck:20.0,
		Fy:250.0,
		Lspan:4500.0,
		Dused:170.0,
		Efcvr:25.0,
		Astm:1130.0,
		DM:0.0,
		Endc:1,
		Type:1,
	}
	astreq := 1044.0; mspan := 30.0
	rezstring += "hulse ex. 7.1\n"
	sdchk, sd, dserve := SlbSdratBs(s, mspan, astreq, s.Astm, s.Dused)
	rezstring += fmt.Sprintf("sdchk %v sd rat %f dserve %f\n", sdchk, sd, dserve)
	rezstring += "mosley ex. 8.5\n"
	
	s = &RccSlb{
		Fck:25.0,
		Fy:500.0,
		Lx:4500.0,
		Ly:6300.0,
		Dused:220.0,
		Efcvr:35.0,
		Astm:646.0,
		DM:0.0,
		Endc:10,
		Type:2,
	}
	astreq = 588.0
	mspan = 45.0
	sdchk, sd, dserve = SlbSdratBs(s, mspan, astreq, s.Astm, s.Dused)
	rezstring += fmt.Sprintf("sdchk %v sd rat %f dserve %f\n", sdchk, sd, dserve)
		
	wantstring := `hulse ex. 7.1
sdchk true sd rat 34.826084 dserve 129.213495
mosley ex. 8.5
sdchk true sd rat 25.490007 dserve 176.539768
`
	if rezstring != wantstring {
		t.Errorf("span depth ratio (bs) check test failed")
		fmt.Println(rezstring)
	}
}
