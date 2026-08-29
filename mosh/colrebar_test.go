package barf

import (
	"fmt"
	"testing"
	
)

func TestColBarGen(t *testing.T){
	var c RccCol
	c = RccCol{Ast:2000,Asc:3000,Asteel:5000,Styp:1,B:300,H:500, Cvrc: 45.0, Cvrt:45.0}
	err := c.SecInit()
	err = c.BarGen()
	fmt.Println(err)
	c.Table(true)
	c.Draw()
}

func TestColOptPso(t *testing.T){}
