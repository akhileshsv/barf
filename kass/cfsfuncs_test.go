package barf

import (
	"log"
	"testing"
)

func TestCfsCp(t *testing.T){
	sectyp := 1
	sdx := 5
	frmtyp := 2
	ax := 1
	log.Println(GetCfsCp(frmtyp, sectyp, sdx, ax))
}
