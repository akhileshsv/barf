package barf

import (
	"testing"
)

func TestPortalInit(t *testing.T){
	f := &Portal{
		Nbays:1,
		Span:18.0,
		Slope:0.1,
		Spacing:6.625,
		Height:9.0,
		Fixbs:false,
		Haunch:true,
		LL:0.6,
		DL:0.175,
		Vz:44.0,
		Cpi:0.2,
		Gable:true,
		Ly:26.5,
		Term:"qt",
	}
	PortalInit(f)
}

func TestHassan(t *testing.T){
	HassanEx5()
}
/*
   
   puf panel 0.15 kn/m2 tops
   services  0.1 - 0.4 kn/m2
   ac sheeting 0.04 kn/m2
   gi sheeting 0.085 kn/m2
   fixtures 0.025 kn/m2
   upvc gutter 0.04 kn/rmt
   sag rods 0.03 kn/m2
   cfs z section 200Z15 4.5 kg/rmt = 0.045 kn/rmt
*/
