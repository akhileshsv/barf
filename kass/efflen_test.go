package barf

import (
	"fmt"
	"testing"
)


func TestEffLen(t *testing.T){
	//EffLen(sway bool, iter int, l0, ga, gb float64)
	//var rezstring string
	var sway bool
	var iter int
	var l0, le, ga, gb float64
	iter = 10
	l0 = 0.6
	ga = 5.0
	gb = 1.0
	le = EffLen(sway, iter, l0, ga, gb)
	fmt.Println("effective len factor->",le)
	sway = true
	iter = 10
	l0 = 1.2
	ga = 5.0
	gb = 5.0
	le = EffLen(sway, iter, l0, ga, gb)
	fmt.Println("effective len factor->",le)

}
