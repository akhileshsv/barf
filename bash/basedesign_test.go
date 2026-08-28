package barf

import (
	"testing"
)

func TestSlbCbsDz(t *testing.T){
	t.Log("duggal ex. 11.1 ishb 350(2) pu 1500")
	c := &Cbs{
		Sname:"i",
		Sdx:15,
		Code:1,
		Fck:20,
		Pu:1500,
	}
	err := SlbCbsDz(c)
	t.Log(err)
	
	t.Log("bhavk ex. 6.10 ishb 300x577 pu 1000")
	c = &Cbs{
		Sname:"i",
		Sdx:22,
		Code:1,
		Fck:20,
		Pu:1000,
	}
	err = SlbCbsDz(c)
	t.Log(err)

}

func TestSlbCbsEx(t *testing.T){
	t.Log("uc 203x203x71, pu 600 kn")
	c := &Cbs{
		Sname:"uc",
		Sdx:24,
		Code:1,
		Fck:25,
		Pu:600,
	}
	err := SlbCbsDz(c)
	t.Log(err)

	t.Log("uc 203x203x46, pu 231 kn")

	c = &Cbs{
		Sname:"uc",
		Sdx:27,
		Code:1,
		Fck:25,
		Pu:231,
	}
	err = SlbCbsDz(c)
	t.Log(err)

	t.Log("uc 152x152x30, pu 231 kn")

	c = &Cbs{
		Sname:"uc",
		Sdx:30,
		Code:1,
		Fck:25,
		Pu:231,
	}
	err = SlbCbsDz(c)
	t.Log(err)

}
