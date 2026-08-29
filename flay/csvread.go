package barf

import (
	"os"
	"runtime"
	"encoding/csv"
	"log"
	"strconv"
	"path/filepath"
	//"github.com/go-gota/gota/dataframe"
	//"github.com/go-gota/gota/series"
	//"strings"
	//"errors"
)

//CsvIntMat reads a matrix of int values from a csv file
func CsvIntMat(filename string) ([][]int,error){
	mat := [][]int{}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return mat, err
	}
	lines, _ := csv.NewReader(file).ReadAll()

	for _, line := range lines {
		row := []int{}
		for _, val := range line {
			var v int
			var err error
			if val == "" {
				v = 0
			} else {
				v, err = strconv.Atoi(val)
			}
			if err != nil {
				//log.Fatal(err)
				//return mat, err
				continue
			}
			row = append(row, v)
		}
		mat = append(mat, row)
	}
	return mat, err
}

//ReadSqrMat calls CsvIntMat to read (squarify/basic) room connections csv
func ReadSqrMat()(mat [][]int, err error){
	//read data/flay/room_connections_squarify.csv
	//bhaiyya check out runtime.caller vs os.exec?
	//WILL runtime.caller WERK EVEN AFTER KOMPILE
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	fname := filepath.Join(basepath,"../data/flay","room_connections_squarify.csv")
	mat, err = CsvIntMat(fname)
	if err != nil{
		log.Fatal(err)
	}
	return
}
