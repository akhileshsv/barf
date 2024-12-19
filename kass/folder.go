package barf

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"path/filepath"
	"log"
	"math/rand"
	"time"
)

//GenFolder generates a random folder name from a list of 100 words ish?
func GenFolder(fname string) (foldr string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	werds := filepath.Join(basepath,"../data/out/werds.txt")

   // read the whole content of file and pass it to file variable, in case of error pass it to err variable
	file, err := os.Open(werds)
	if err != nil {
		fmt.Printf("Could not open the file due to this %s error \n", err)
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string
	
	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	if err = file.Close(); err != nil {
		fmt.Printf("Could not close the file due to this %s error \n", err)
	}
	
	if fname == ""{
		//generate random combo of three words
		rand.Seed(time.Now().Unix())
		for i := 0; i < 3; i++{
			fname += fileLines[rand.Intn(100)]
			if i < 2{fname += "_"}
		}
	}
	foldr = filepath.Join(basepath,"../data/out/",fname)
	if err := os.Mkdir(foldr, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	return 
}
