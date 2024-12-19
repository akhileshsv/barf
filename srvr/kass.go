package barf

import (
	"os"
	"fmt"
	"net/http"
	"encoding/json"
	kass"barf/kass"
)


func modhtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var err error
		for k, v := range r.Form{
			if k == "Frmstr"{
				switch v[0]{
					case "1db":
					err = tmod1db.Execute(w,nil)
					case "2dt":
					err = tmod2dt.Execute(w,nil)
					case "2df":
					err = tmod2df.Execute(w,nil)
					case "3dt":
					err = tmod3dt.Execute(w,nil)
					case "3dg":
					err = tmod3dg.Execute(w,nil)
					case "3df":
					err = tmod3df.Execute(w,nil)
				}
				if err != nil{
					terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
				}
			}
		}
	}
}

func modnphtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		// var htext []byte
		// var err error
		fpath := "srvr/templates/ex/modnp"
		for k, v := range r.Form{
			if k == "Frmstr"{
				fpath = fpath + v[0] + ".html"
				htext, err := os.ReadFile(fpath)
				if err != nil{
					terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
				}
				hstr := string(htext)
				fmt.Fprint(w, hstr)
			}
		}
	}
}

func modephtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		// var htext []byte
		// var err error
		fpath := "srvr/templates/ex/modep"
		for k, v := range r.Form{
			if k == "Frmstr"{
				fpath = fpath + v[0] + ".html"
				htext, err := os.ReadFile(fpath)
				if err != nil{
					terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
				}
				hstr := string(htext)
				fmt.Fprint(w, hstr)
			}
		}
	}
}


func calcmod(w http.ResponseWriter, r *http.Request) {	
	if r.Method == "GET" {
		err := tcmod.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		}		
	} else {
		r.ParseForm()
		mod, err := parsemodel(r)
		if err != nil{
			terror.Execute(w, fmt.Errorf("model parse error->\n%s",err))
			return
		}
		err = kass.CalcMod(&mod, mod.Frmstr, mod.Term, false)
		if err != nil{			
			terror.Execute(w, fmt.Errorf("model calculation error->\n%s",err))
			return	
		}
		// log.Println("TEXPLOTS-",mod.Txtplots)
		err = tcmodrez.Execute(w, &mod)
		if err != nil{
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		}
	}
}

func calcnp(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		err := tcnp.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		}		
	} else {
		r.ParseForm()
		mod, err := parsemodel(r)
		mod.Calc = false
		if err != nil{
			terror.Execute(w, fmt.Errorf("model parse error->\n%s",err))
			return
		}
		err = kass.CalcModNp(&mod, mod.Frmstr, mod.Term, false)
		if err != nil{			
			terror.Execute(w, fmt.Errorf("model calculation error->\n%s",err))
			return	
		}
		err = tcmodrez.Execute(w, &mod)
		if err != nil{
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		}
	}

}


func calcep(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "GET" {
		err := tcep.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		}		
	} else {
		r.ParseForm()
		mod, err := parsemodel(r)
		if err != nil{
			terror.Execute(w, fmt.Errorf("model parse error->\n%s",err))
			return
		}
		err = kass.CalcEpFrm(&mod)
		if err != nil{			
			terror.Execute(w, fmt.Errorf("model calculation error->\n%s",err))
			return	
		}
		err = tcmodrez.Execute(w, &mod)
		if err != nil{
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		}
	}
}

//parsemodel parses/checks form data and returns a kass.Model
func parsemodel(r *http.Request)(mod kass.Model, err error){
	//must be a better way
	//a better day :(
	var jsonstr string
	strmap := map[string]int{
		"Id":1,"Term":1,"Units":1,"Frmstr":1,"Zup":2,"Fy":2,"Pg":2,"Dmax":2,"Ngrps":2,"Opt":2,"Json":-1,
	}
	
	if val, ok := r.Form["Json"]; ok && val[0] != ""{
		if json.Valid([]byte(val[0])){
			jsonstr = val[0]
		} else {
			err = fmt.Errorf("invalid json in jsonstr form field - [%v]",val[0])
			return
		}
	} else {
		jsonstr += "{"
		for key, val := range r.PostForm{
			// log.Println("processing key->",key," val->",val)
			
			if v, ok := strmap[key]; ok{
				switch v{
					case -1:
					continue
					case 1:
					jsonstr += fmt.Sprintf("\"%v\":\"%s\",",key,val[0])
					case 2:
					jsonstr += fmt.Sprintf("\"%v\":%v,",key,val[0])
				}
			} else {
				jsonstr += fmt.Sprintf("\"%v\":%v,",key,val)
			}
		}
		jsonstr = jsonstr[:len(jsonstr)-1]
		jsonstr += "}"
	}
	err = json.Unmarshal([]byte(jsonstr), &mod)
	switch mod.Frmstr{
		case "1db":
		mod.Ncjt = 2
		mod.Calc = true
		case "2dt":
		mod.Ncjt = 2
		case "2df":
		mod.Ncjt = 3
		mod.Calc = true
		case "3dt":
		mod.Ncjt = 3
		case "3dg":
		mod.Ncjt = 3
		case "3df":
		mod.Ncjt = 6
	}
	//mod.Term = "svg"
	mod.Web = true
	return
}


