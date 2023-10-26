package barf

import (
	"fmt"
	"testing"
	kass"barf/kass"
)

func TestTrussDz(t *testing.T){
	var trs *kass.Trs2d
	fmt.Println("swiss roof test")
	trs = &kass.Trs2d{
		Id:        "swizz",
		Cmdz:      []string{"2dt","mmks"},
		Typ:       2,
		Cfg:       1,
		Mtyp:      3,
		Span:      3660.0,
		Spacing:   2400.0,
		Slope:     1.5,
		Purlinspc: 2500,
		DL:        350,
		Clctyp:    0,
		Vb:        50,
		Frmtyp:    2,
		Term:      "mono",
		//Bcslope:   1.0/12.0,
	}
	//GenTruss(trs)
	TrussDz(trs)

}
