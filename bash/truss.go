package barf

import (
	"fmt"
	kass"barf/kass"
)

func TrussDz(t *kass.Trs2d) (err error){
	err = t.Calc()
	if err != nil{
		fmt.Println(err)
	}
	return
}

func TrussOpt(t *kass.Trs2d) (err error){
	err = fmt.Errorf("truss optimization not written yet (at all)")
	return
}
