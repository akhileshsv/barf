package barf

import (
	"fmt"
	"testing"
	//"github.com/go-gota/gota/dataframe"
)

func TestStlDf(t *testing.T){
	fmt.Println(StlDfIs(1))
	for i, styp := range StlStyps{
		df, _ := GetStlDf(i+1)
		fmt.Println("styp,len->",styp,df.Nrow())
	}
	cp, err := GetStlCp(1, 1, 0, 1) 
	fmt.Println(cp, err)
	cp, err = GetStlCp(1, 1, 1, 2)
	fmt.Println(cp, err)
}
