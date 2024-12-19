package barf
/*
import (
	"fmt"
	"log"
	"math"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	kass"barf/kass"
)


   func PlotPortal(coords []float64) (pltstr string) {
	//get plotscript filepath
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	pltskript := filepath.Join(basepath,"/plotcolnm.gp")
	var data string
	for idx, pu := range pus {
		mu := mus[idx]
		data += fmt.Sprintf("%v %v\n", pu, mu)
	}
	f, e1 := os.CreateTemp("", "mosh")
	if e1 != nil {
		fmt.Println(e1)
	}
	defer f.Close()
	defer os.Remove(f.Name())	
	_, e1 = f.WriteString(data)
	if e1 != nil {
		fmt.Println(e1)
	}
	cmd := exec.Command("gnuplot","-c",pltskript,f.Name(),"dumb")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outstr, errstr := stdout.String(), stderr.String()
	if err != nil {
		fmt.Println(err)
	}
	if errstr != "" {
		fmt.Println(errstr)
	}
	return outstr
}
*/
