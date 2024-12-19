package barf

import (
	"os"
	"fmt"
	"log"
	"strings"
	"runtime"
	"io/ioutil"
	"path/filepath"
	"github.com/f1bonacc1/glippy"
	"github.com/AlecAivazis/survey/v2"
)

//suggestFiles suggests a list of files for survey menu
func suggestFiles(toComplete string) []string {
	files, _ := filepath.Glob(toComplete + "*")
	return files
}

//getjsonfile reads in a jsonfile to a string of bytes
func getjsonfile()(bytestr []byte, err error){
	var filename string
	var q = []*survey.Question{
		{
			Name: "file",
			Prompt: &survey.Input{
				Message: "enter path to json file:",
				Suggest: suggestFiles,
				Help:    "abs/relative filepath",
			},
			Validate: survey.Required,
		},
	}
	err = survey.Ask(q, &filename)
	if err != nil {
		return
	}
	jsonfile, e := ioutil.ReadFile(filename)
	if e != nil{
		err = e
		return
	}
	bytestr = []byte(jsonfile)
	return 
}

//readjsonstr reads in a base .json file and returns the string for further edits
func readjsonstr(basefile string) (jsonstr, helpstr string, err error){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	filename := filepath.Join(basepath,"../data/json/",basefile)
	jsonfile, e := ioutil.ReadFile(filename)
	if e != nil{
		err = e
		return
	}
	jsonstr = string(jsonfile)
	
	//copy json to clipboard (glippy is life)
	glippy.Set(jsonstr)
	//read helpfile
	txtfile, e := ioutil.ReadFile(strings.Replace(filename, ".json",".txt",-1))
	if e != nil{
		err = e
		return
	}
	helpstr = string(txtfile)
	return
}

//readjsontxt reads in a base .json file, reads in the edited struct via $EDITOR (notepad/nano/etc)
func readjsontxt(basefile string) (bytestr []byte, err error){
	var content string
	jsonstr, helpstr, e := readjsonstr(basefile)
	if e != nil{err = e; return}
	fmt.Println(ColorYellow, jsonstr, ColorReset)
	prompt := &survey.Editor{
		Message: "save struct on exit",
		FileName: "*.json",
		//Help: "copy/paste, edit and save",
		Help:helpstr,
	}
	survey.AskOne(prompt, &content)
	bytestr = []byte(content)
	return
}

//savemenu asks whether or not to save a file 
func savemenu() (savez bool){ 
	savez = true
	prompt := &survey.Confirm{
		Message: "save to disk?",
	}
	survey.AskOne(prompt, &savez)
	return
}

//save json saves slice of bytes to a jsonfile in /data/out/"filename"
func savejson(jsondat []byte, filename string) (err error){
	filename = filename + ".json"
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	dest := filepath.Join(basepath,"../data/out/",filename)
	if err = os.WriteFile(dest, jsondat, 0666); err != nil {
		log.Println(err)
	}
	return
}

//getfilename reads in filename from stdin
func getfilename() (filename string){
	prompt := &survey.Input{
		Message: "enter struct filename",
	}
	survey.AskOne(prompt, &filename)
	return

}

//printmenu prints a slice of menu items to the terminal
func printmenu(message string, menus []string) (choice int){
	prompt := &survey.Select{
		Message: message,
		Options: menus,
	}
	survey.AskOne(prompt, &choice)
	return
}

//getterminal displays a list of available gnuplot terminals
func getterminal() (term string){
	prompt := &survey.Select{
		Message: "set gnuplot terminal",
		Options: []string{"mono","dumb","qt","wxt","svg","none"},
	}
	survey.AskOne(prompt, &term)
	return
}

//InitMenu is the main entry func from menu/flags
func InitMenu(term, sub string){
	fmt.Println(ColorPurple,icon_barf,ColorReset)
	if sub == ""{
		running := true
		for running{
			choice := printmenu("choose module",main_menus)
			switch choice{
				case 0:
				kassmenu(term)
				case 1:
				moshmenu(term)
				case 2:
				bashmenu(term)
				case 3:
				tmbrmenu(term)
				//case 4:
				//flaymenu(term)
				case 4:
				running = false
				break
			}
		}
		return
	}
	switch sub{
		case "kass","calc":
		kassmenu(term)
		case "mosh","rcc":
		moshmenu(term)
		case "bash","stl":
		bashmenu(term)
		case "tmbr","wood":
		tmbrmenu(term)
		//case "flay","layout":
		//flaymenu(term)
		default:
		log.Println("invalid sub menu string")
	}
	return
}
