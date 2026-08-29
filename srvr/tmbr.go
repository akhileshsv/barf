package barf

import (
	"fmt"
	"net/http"
	"encoding/json"
	tmbr"barf/tmbr"
	//kass"barf/kass"
)


func timber(w http.ResponseWriter, r *http.Request){
	err := ttimber.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
}

func calctmbrbeam(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := ttmbrbeam.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error\n%s",err))
		}
	} else{
		
		r.ParseForm()
		jsonstr := parsetmbr(r)
		var b tmbr.WdBm
		err := json.Unmarshal([]byte(jsonstr), &b)
		if err != nil{
			terror.Execute(w, fmt.Errorf("form read error->\n%s",err))
			return
		}
		//b.Endc = 1
		
		
		err = tmbr.BmDesign(&b)
		if err != nil{
			terror.Execute(w, fmt.Errorf("beam design error->\n%s",err))
			return
		}
		
		err = trrez.Execute(w, &b)
		if err != nil{
			terror.Execute(w, fmt.Errorf("major template error\n%s",err))
			return
		}
	}
}

func calctmbrcol(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := ttmbrcol.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error\n%s",err))
		}
	} else {	
		r.ParseForm()
		jsonstr := parsetmbr(r)
		var c tmbr.WdCol
		err := json.Unmarshal([]byte(jsonstr), &c)
		if err != nil{
			terror.Execute(w, fmt.Errorf("form read error->\n%s",err))
			return
		}
		err = tmbr.ColDz(&c)
		if err != nil{
			terror.Execute(w, fmt.Errorf("timber column design error->\n%s",err))
			return
		}
		err = trrez.Execute(w, &c)
		if err != nil{
			terror.Execute(w, fmt.Errorf("major template error\n%s",err))
		}
	}
}

func tmbrbmhtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var err error
		//log.Println(r.Form)
		if val, ok := r.Form["Rdprp"]; !ok{
			terror.Execute(w, fmt.Errorf("major error parsing htmx form"))		
		} else {
			//log.Println(r.Form["Rdprp"])
			switch val[0]{
				case "false":
				err = ttmbmprp.Execute(w,nil)
				case "true":
				err = ttmbmgrp.Execute(w,nil)
			}
			if err != nil{
				terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
			}
		}
	}
}


func tmbrcolhtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var err error
		if val, ok := r.Form["Rdprp"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "true":
				err = ttmcolgrp.Execute(w,nil)
				case "false":
				err = ttmcolprp.Execute(w,nil)
			}
			if err != nil{
				terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
			}
		}
	}

}


//parsetmbr parses/checks form data and returns the json string
func parsetmbr(r *http.Request)(jsonstr string){
	if val, ok := r.Form["Json"]; ok{
		if val[0] != "" && json.Valid([]byte(val[0])){
			jsonstr = val[0]
			return
		}		
	}
	jsonstr += "{"
	var prpstr string
	strmap := map[string]int{
		"Title":1,
		"Term":1,
		"Pg":3,
		"Em":3,
		"Fcp":3,
		"Ftp":3,
		"Fc":3,
		"Ft":3,
		"Fcb":3,
		"Fv":3,
		"Fvp":3,
	}
	for key, val := range r.PostForm{
		fmt.Println("processing key->",key," val->",val, "len->", len(val))
		if len(val[0]) == 0{
			continue
		}
		
		if mval, ok := strmap[key];!ok{
			jsonstr += fmt.Sprintf("\"%v\":%v,",key,val[0])
		} else {
			switch mval{
				case 1:
				jsonstr += fmt.Sprintf("\"%v\":\"%s\",",key,val[0])
				case 2:
				jsonstr += fmt.Sprintf("\"%v\":%v,",key,val)
				case 3:
				if prpstr == ""{
					prpstr += "\"Prp\":{"
				}
				prpstr += fmt.Sprintf("\"%v\":%v,",key,val[0])
			}
		}
	}
	jsonstr = jsonstr[:len(jsonstr)-1]
	if prpstr != ""{
		
		prpstr = prpstr[:len(prpstr)-1]
		prpstr += "}"
		jsonstr += "," + prpstr
	}
	jsonstr += "}"
	fmt.Println("jsonstr->\n",jsonstr)
	return
}
