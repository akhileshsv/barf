package barf

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	kass"barf/kass"
	bash"barf/bash"
)


func steel(w http.ResponseWriter, r *http.Request){
	err := tsteel.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
}


func stlsdxs(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := tstlsecs.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error\n%s",err))
		}
	}
}

func stlfloor(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := tstlsecs.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error\n%s",err))
		}
	}
}

func stlbeam(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := tstlbeam.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		}
	} else {
		r.ParseForm()
		jsonstr := parsebash(r)
		var bm bash.Bm
		err := json.Unmarshal([]byte(jsonstr), &bm)
		if err != nil {
			terror.Execute(w, fmt.Errorf("form read error->\n%s",err))
			return
		}
		bm.Web = true
		bm.Selfwt = true
		bm.Drawrez = true
		bm.Dsgn = true
		if bm.Sname == "ismb"{
			fmt.Println("changing onism")
			bm.Sname = "i"
			bm.Onism = true
			fmt.Println("sname, onism",bm.Sname,bm.Onism)
			
		}
		err = bash.BmDz(&bm)
		if err != nil {
			terror.Execute(w, fmt.Errorf("beam design error->\n%s",err))
			return
		}
		err = tstlrez.Execute(w, &bm)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
			return
		}
	}
}


func stlcol(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := tstlcol.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error\n%s",err))
		}
	}else {
		r.ParseForm()
		jsonstr := parsebash(r)
		var col bash.Col
		err := json.Unmarshal([]byte(jsonstr), &col)
		if err != nil {
			terror.Execute(w, fmt.Errorf("form read error->\n%s",err))
			return
		}
		col.Web = true
		// col.Drawrez = true
		if col.Sname == "ismb"{
			fmt.Println("changing onism")
			col.Sname = "i"
			col.Onism = true
			fmt.Println("sname, onism",col.Sname,col.Onism)
			
		}
		err = bash.ColDz(&col)
		if err != nil {
			terror.Execute(w, fmt.Errorf("column design error->\n%s",err))
			return
		}
		err = tstlrez.Execute(w, &col)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
			return
		}
	}
	
}

func stlcolhtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()	
		var err error
		fmt.Println("dsgn,code",r.Form["Dsgn"],r.Form["Code"])
		v1 := r.Form["Dsgn"][0]
		v2 := r.Form["Code"][0]
		switch{
			case v1 == "true" && v2 == "1":
			err = tstlcoldz.Execute(w, nil)
			case v1 == "true" && v2 == "2":
			err = tstlcoldzbs.Execute(w, nil)
			case v1 == "false" && v2 == "1":
			err = tstlcolchk.Execute(w, nil)
			case v1 == "false" && v2 == "2":
			err = tstlcolchkbs.Execute(w, nil)
		}
		if err != nil{
			terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
		} 
	
	}
	return
}

func stltrsmodopthtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		
		var err error
		if val, ok := r.Form["Frmstr"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "2dt":
				err = tsmod2dtopt.Execute(w,nil)	
				if err != nil{
					terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
				}
				case "3dt":
				err = tsmod3dtopt.Execute(w,nil)
				if err != nil{
					terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
				}
			}
		}
	}
	return
}

func stltrsmodopt(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := tstltrsmodopt.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error\n%s",err))
		}
	} else {
		r.ParseForm()
		mod, err := parsemodel(r)
		if err != nil{
			terror.Execute(w, fmt.Errorf("model parse error->\n%s",err))
			return
		}
		var modr kass.Model
		modr, err = bash.OptTrsMod(mod)
		if err != nil{			
			terror.Execute(w, fmt.Errorf("model calculation error->\n%s",err))
			return	
		}
		// log.Println("TEXPLOTS-",mod.Txtplots)
		err = tcmodrez.Execute(w, &modr)
		if err != nil{
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		}

	}
}

func calcstltrsgen(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := tstltrsgen.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		
		r.ParseForm()
		jsonstr := parsebash(r)
		var t kass.Trs2d
		err := json.Unmarshal([]byte(jsonstr), &t)
		if err != nil{
			terror.Execute(w, fmt.Errorf("form read error->\n%s",err))
			return
		}
		t.Web = true
		switch t.Term{
			case "svgmono","dxf","svg":
			default:
			t.Term = "svg"
		}
		if t.Fop == 0{
			t.Opt = 0
		}
		switch t.Opt{
			case 0:
			//add truss dz func here
			err = bash.TrussDz(&t)
			if err != nil{
				terror.Execute(w, fmt.Errorf("2d truss design error->\n%s",err))
				return
			}
			default:
			//add truss opt func here
			err = bash.TrussOpt(&t)
			
			if err != nil{
				terror.Execute(w, fmt.Errorf("2d truss design error->\n%s",err))
				return
			}
		}
		err = tstltrsrez.Execute(w, &t)
		if err != nil{
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		
		}
	}
}


func stltrsgenhtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		
		var err error
		if val, ok := r.Form["Fop"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "0":
				err = tstrsgen.Execute(w,nil)	
				//dz
				if err != nil{
					terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
				}
				case "1":
				//opt form
				err = tstrsgenopt.Execute(w,nil)
				if err != nil{
					terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
				}
			}
		}
	}
}


func stltrs(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := tstltrs.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error\n%s",err))
		}
	}
}

//parsebash parses/checks form data and returns the json string
func parsebash(r *http.Request)(jsonstr string){
	var bgstr, wgstr string
	strmap := map[string]int{
		"Title":1,
		"Id":1,
		"Sname":1,
		"Term":1,
		"Lspans":2,
		"Ldcases":2,
		"Dia":3,
		"Ni":3,
		"Nj":3,
		"Pitch":3,
		"Distance":3,
		"Edged":3,
		"Endd":3,
		"Wltyp":4,
		"L1":4,
		"L2":4,
	}
	if val, ok := r.Form["Json"]; ok{
		if val[0] != "" && json.Valid([]byte(val[0])){
			jsonstr = val[0]
			return
		}
	}
	jsonstr += "{"
	for key, val := range r.PostForm{
		log.Println("processing key->",key," val->",val, "len->", len(val))
		if len(val[0]) == 0{
			continue
		}
		if mval, ok := strmap[key];!ok{
			jsonstr += fmt.Sprintf("\"%v\":%v,",key,val[0])
		} else {
			switch mval{
				case -1:
				continue
				case 1:
				jsonstr += fmt.Sprintf("\"%v\":\"%s\",",key,val[0])
				case 2:
				jsonstr += fmt.Sprintf("\"%v\":%v,",key,val)
				case 3:
				if bgstr == ""{
					bgstr += "{\"Bg\":{"
				}
				bgstr += fmt.Sprintf("\"%v\":%v,",key,val)
				case 4:
				
				if wgstr == ""{
					wgstr += "{\"Wld\":{"
				}
				wgstr += fmt.Sprintf("\"%v\":%v,",key,val)
				
			}
		}
	}
	jsonstr = jsonstr[:len(jsonstr)-1]	
	if bgstr != ""{
		bgstr += "}"
		jsonstr += bgstr
	} else if wgstr != ""{
		wgstr += "}"
		jsonstr += wgstr
	}
	jsonstr += "}"
	log.Println("jsonstr->\n",jsonstr)
	return
}
