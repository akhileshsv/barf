package barf

import (
	"testing"
)

func TestCmpBmChk(t *testing.T){
	t.Log("testing insdag ex.1-ss beam")
	c := CmpBm{
		Fck:30,
		Fy:250,
		Dslb:125,
		Llk:0.75,
		DL:0.5,
		LL:4.5,
		Lslb:3000,
		Lspan:10000,
		Lrfac:1.0,
		PSFs:[]float64{1.35,1.5},
		Sdx:16,
		Sname:"i",
		Verbose:true,
	}
	_, err := c.ChkSec(16)
	t.Log(err)
}

func TestCmpBmDz(t *testing.T){
	t.Log("testing insdag ex.1-ss beam")
	c := CmpBm{
		Fck:30,
		Fy:250,
		Dslb:125,
		Llk:0.75,
		DL:0.5,
		LL:4.5,
		Lslb:3000,
		Lspan:10000,
		Lrfac:1.0,
		PSFs:[]float64{1.35,1.5},
		Sname:"i",
		Nsecs:3,
		Verbose:true,
	}
	err := CmpBmDz(&c)
	t.Log(err)
	t.Log(ColorRed,"testing draw",ColorReset)
	c.Draw()	
}

func TestCmpBmSmsh(t *testing.T){
	t.Log("testing smsh joist-ss beam")
	c := CmpBm{
		Fck:25,
		Fy:250,
		Dslb:100,
		Llk:0.75,
		DL:1.0,
		LL:3.0,
		Lslb:2250,
		Lspan:7000,
		Lrfac:1.0,
		PSFs:[]float64{1.35,1.5},
		Sname:"i",
		Nsecs:5,
		Verbose:true,
	}
	err := CmpBmDz(&c)
	t.Log(err)
	t.Log(ColorRed,"testing draw",ColorReset)
	c.Draw()
}
