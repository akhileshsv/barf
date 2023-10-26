package barf

import (
	"os"
	"bufio"
	"bytes"
	//"encoding/json"
	"io"
	"strings"
)

//ReadMIn reads in multiline text input from stdin
func ReadMIn() (inmap map[string]string, err error){
	inmap = make(map[string]string)
	var buf bytes.Buffer
	reader := bufio.NewReader(os.Stdin)
	input := make(map[string]string)
	var line string
	for {	
		line, err = reader.ReadString('\n')
		if err == nil && line != ""{
			key, val := strings.Split(line," ")[0], strings.Split(line," ")[1]
			key, val = strings.ToLower(key), strings.ToLower(val)
			key, val = strings.TrimSpace(key), strings.TrimSpace(val)
			input[key] = val
		}
		if err != nil {
			if err == io.EOF {
				buf.WriteString(line)
				break
			} else {
				return
			}   
		}
	}
	err = nil
	return
}
