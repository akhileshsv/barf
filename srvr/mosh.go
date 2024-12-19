package barf

import (
	"fmt"
	"net/http"
	"encoding/json"
	mosh"barf/mosh"
	kass"barf/kass"
)

//rcftnghtml returns footing html
func rcftnghtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var err error
		if val, ok := r.Form["Sloped"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "true":
				err = trcftng.Execute(w,nil)
				case "false":
				err = trcftngpad.Execute(w,nil)
			}
			if err != nil{
				terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
			}
		}
	}
	
}

//rccolhtml returns col section html
func rccolhtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var err error
		if val, ok := r.Form["Dsgn"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "true":
				err = trccoldz.Execute(w,nil)
				case "false":
				err = trccolaz.Execute(w,nil)
			}
			if err != nil{
				terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
			}
		}
	}
	
}


//rccolstyphtml returns html for a column section type
func rccolstyphtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		hstr := `<div id="rcccol-dims">
    <label for="Dims">section dims(mm)</label>
    <input type="text" name="Dims" required class="nes-input" `
		if val, ok := r.Form["Styp"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "-1":
				hstr = `<div id="rcccol-dims">
    <label for="Coords">section coords(mm)</label>
    <input type="text" name="Coords" required class="nes-input" value="[-212,0],[0,-212],[212,0],[0,212],[-212,0]" placeholder=""/>
    </div>`
				case "0":
				hstr += `value="300" placeholder="300(dia)"/>
    </div>`
				case "1":
				hstr += `value="350,450" placeholder="350,450(bxh)"/>
    </div>`
				default:
				hstr += `value="151.5" placeholder="151.5(side)"/>
    </div>`
			}
			fmt.Fprint(w, hstr)
		}
	}
	
}

//rcbmstyphtml returns html for custom beam section types
func rcbmstyphtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var hstr string
		if val, ok := r.Form["Styp"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "-1":
				hstr = `<div id="rccbmsec-dims">
    <label for="Coords">section coords(mm)</label>
    <input type="text" name="Coords" class="nes-input" value="[0 0], [235 0], [317.5 550], [-82.5 550], [0 0]" placeholder="[0 0], [235 0], [317.5 550], [-82.5 550], [0 0](list of cw/ccw coords)"  onfocus = "this.value=''"/>
    </div>`
				default:
				hstr = `<div id="rccbmsec-dims">
    <label for="Dims">section dims(mm)</label>
    <input type="text" name="Dims" class="nes-input" value="400,550,235" placeholder="400,550,235"  onfocus = "this.value=''"/>
    </div>`
			}
			fmt.Fprint(w, hstr)
		}
	}
	
}


//rcbmsechtml returns beam section html
func rcbmsechtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var err error
		if val, ok := r.Form["Dsgn"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "true":
				err = trcbmdz.Execute(w,nil)
				case "false":
				err = trcbmaz.Execute(w,nil)
			}
			if err != nil{
				terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
			}
		}
	}
	
}

//rcsfhtml returns subframe html
func rcsfhtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var err error
		if val, ok := r.Form["Fop"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "0":
				err = trcsfdz.Execute(w,nil)
				case "1":
				err = trcsfopt.Execute(w,nil)
			}
			if err != nil{
				terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
			}
		}
	}
}

//rcbmcshtml returns beamcs html
func rcbmcshtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var err error
		if val, ok := r.Form["Fop"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "0":
				err = trcbmcsdz.Execute(w,nil)
				case "1":
				err = trcbmcsopt.Execute(w,nil)
			}
			if err != nil{
				terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
			}
		}
	}
}


//rcbmsshtml returns beamss html
func rcbmsshtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var err error
		if val, ok := r.Form["Fop"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "0":
				err = trcbmssdz.Execute(w,nil)
				case "1":
				err = trcbmssopt.Execute(w,nil)
			}
			if err != nil{
				terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
			}
		}
	}
}

//rcf2dhtml returns html for rcc 2d frame design/opt
func rcf2dhtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var err error
		if val, ok := r.Form["Fop"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "0":
				err = trcf2ddz.Execute(w,nil)
				case "1":
				err = trcf2dopt.Execute(w,nil)
			}
			if err != nil{
				terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
			}
		}
	}
}

//rslabhtml returns slab form html
func rslabhtml(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()
		var err error
		if val, ok := r.Form["Typstr"]; !ok{
			terror.Execute(w, fmt.Errorf("error parsing htmx form"))		
		} else {
			switch val[0]{
				case "clvr":
				err = tslbclvr.Execute(w,nil)
				case "1w":
				err = tslb1w.Execute(w,nil)
				case "1wcs":
				err = tslb1wcs.Execute(w,nil)
				case "2w":
				err = tslb2w.Execute(w,nil)
				case "2wcs":
				err = tslb2wcs.Execute(w,nil)
			}
			if err != nil{
				terror.Execute(w, fmt.Errorf("major (htmx) template error->\n%s",err))
			}
		}
	}
}

func rcc(w http.ResponseWriter, r *http.Request){
	err := trcc.Execute(w, nil)
	if err != nil {
		terror.Execute(w, fmt.Errorf("major template error\n%s",err))
	}
	
}

func calcrslab(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := trslab.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error\n%s",err))
		}
	} else {
		r.ParseForm()
		jsonstr := parsemosh(r)
		var s mosh.RccSlb
		err := json.Unmarshal([]byte(jsonstr), &s)
		if err != nil {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			terror.Execute(w, fmt.Errorf("json convert error\n%s",err))
			return
		}
		//time.Sleep(8 * time.Second)
		s.Web = true
		err = mosh.SlbDesign(&s)
		if err != nil {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			terror.Execute(w, err)
			return
		}
		err = trrez.Execute(w, &s)
		if err != nil {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			terror.Execute(w, err)
			return
		}
	}
	
}

func rccbeam(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := trbeam.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		} 
	}		
}

func calcrbeam(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := trbeamsec.Execute(w, nil)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		}
	} else {
		r.ParseForm()
		jsonstr := parsemosh(r)
		var bm mosh.RccBm
		err := json.Unmarshal([]byte(jsonstr), &bm)
		if err != nil {
			terror.Execute(w, fmt.Errorf("form read error->\n%s",err))
			return
		}
		bm.Web = true
		if bm.Vu > 0.0{
			bm.Shrdz = true
		}
		if bm.Styp == -1{
			bm.Sec.Ncs = []int{len(bm.Sec.Coords)}
			bm.Sec.Wts = []float64{1.0}
		}
		switch bm.Dsgn{
			case true:
			err = mosh.BmDesign(&bm)
			case false:
			err = mosh.BmAnalyze(&bm)
		}
		if err != nil {
			terror.Execute(w, fmt.Errorf("beam design error->\n%s",err))
			return
		}
		bm.Table(false)
		mosh.PlotBmGeom(&bm,bm.Term)
		err = trrez.Execute(w, &bm)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
			return
		}
	}
}

func calcrftng(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := trftng.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		
		r.ParseForm()
		jsonstr := parsemosh(r)
		var f mosh.RccFtng
		err := json.Unmarshal([]byte(jsonstr), &f)
		if err != nil{
			terror.Execute(w, fmt.Errorf("form read error->\n%s",err))
			return
		}
		f.Web = true
		switch f.Term{
			case "dxf","svg","svgmono":
			default:
			f.Term = "svg"
		}
		f.Verbose = false
		err = mosh.FtngDz(&f)
		if err != nil {
			terror.Execute(w, fmt.Errorf("footing design error->\n%s",err))
			return
		}
		f.Table(false)
		mosh.PlotFtngDet(&f)
		err = trrez.Execute(w, &f)
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
			return
		}
	}
}

func calcrcol(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{err := trcol.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		r.ParseForm()
		jsonstr := parsemosh(r)
		var c mosh.RccCol
		err := json.Unmarshal([]byte(jsonstr), &c)
		if err != nil {
			terror.Execute(w, fmt.Errorf("form read error->\n%s",err))
			return
		}
		c.Term = "svg"
		c.Web = true
		switch c.Dsgn{
			case true:
			err = mosh.ColDesign(&c)
			if err != nil {
				terror.Execute(w, fmt.Errorf("column design error->\n%s",err))
				return
			}
			c.Table(false)
			c.PlotColDet()
			
			case false:
			err = mosh.ColAnalyze(&c)
			if err != nil{
				terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
				return
			}
			c.Table(false)
		}
		err = trrez.Execute(w, &c)
		if err != nil{
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
			return
		}
	}
	
}


func calcrcbeam(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		typstr := r.URL.Path[len("/rcc/beam/"):]
		// log.Println("typstr!",typstr)
		var err error
		switch typstr{
			case "ssbeam":
			err = trssbeam.Execute(w, nil)
			case "cbeam":
			err = trcsbeam.Execute(w, nil)			
		}
		if err != nil {
			terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
		}
	} else {
		r.ParseForm()
		jsonstr := parsemosh(r)
		//log.Println(jsonstr)
		var cb, copt mosh.CBm
		err := json.Unmarshal([]byte(jsonstr), &cb)
		if err != nil{
			terror.Execute(w, fmt.Errorf("error reading json ->\n%s",err))
			return
		}
		cb.Web = true
		cb.Dconst = true
		switch cb.Term{
			case "svgmono","dxf","svg":
			default:
			cb.Term = "svg"
		}
		switch cb.Opt{
			case 0:
			cb.Verbose = true
			cb.Foldr = ""
			var bmenv map[int]*kass.BmEnv
			bmenv, err = mosh.CBeamEnvRcc(&cb, cb.Term, true)
			if err != nil{
				terror.Execute(w, fmt.Errorf("cbeam analysis error->\n%s",err))
			}
			_, err = mosh.CBmDz(&cb,bmenv)
			if err != nil{
				terror.Execute(w, fmt.Errorf("cbeam design error->\n%s",err))
				
			}
			pltstr := mosh.PlotCBmDet(cb.Web, cb.Bmvec, cb.RcBm, cb.Foldr, cb.Title, cb.Term)
			cb.Txtplots = append(cb.Txtplots, pltstr)
			cb.Table(false)
			err = trrez.Execute(w, &cb)	
			if err != nil{
				terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
				return
			}
			default:
			copt, err = mosh.CBmOpt(cb)
			if err != nil{
				terror.Execute(w, fmt.Errorf("cbeam opt error->\n%s",err))	
				return
			}			
			err = trrez.Execute(w, &copt)
			if err != nil{
				terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
				return
			}
		}
	}
}

func calcrsubfrm(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := trcsubfrm.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		r.ParseForm()
		jsonstr := parsemosh(r)
		var sf, sfrez mosh.SubFrm
		err := json.Unmarshal([]byte(jsonstr), &sf)
		if err != nil{
			terror.Execute(w, fmt.Errorf("form read error->\n%s",err))
			return
		}
		sf.Web = true
		sf.Verbose = true
		switch sf.Term{
			case "svgmono","dxf","svg":
			default:
			sf.Term = "svg"
		}
		switch sf.Opt{
			case 0:
			switch sf.Slbdz{
				case true:
				_, err = sf.ChainSlab()
				case false:
				_, err = mosh.DzSubFrm(&sf)
			}
			if err != nil{
				terror.Execute(w, fmt.Errorf("subframe design error->\n%s",err))
				return
			}
			err = trrez.Execute(w, &sf)
			if err != nil{
				terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
				return
			}	
			default:
			sf.Dconst = true
			sf.Verbose = false
			sfrez, err = mosh.OptSubFrm(sf)
			if err != nil{
				terror.Execute(w, fmt.Errorf("subframe opt error->\n%s",err))
				return
			}
			err = trrez.Execute(w, &sfrez)
			if err != nil{
				terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
				return
			}	
		}
	}
}

func calcrfrm2d(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := trcfrm2d.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		
		r.ParseForm()
		jsonstr := parsemosh(r)
		var f, frez kass.Frm2d
		err := json.Unmarshal([]byte(jsonstr), &f)
		if err != nil{
			terror.Execute(w, fmt.Errorf("form read error->\n%s",err))
			return
		}
		f.Web = true
		switch f.Term{
			case "svgmono","dxf","svg":
			default:
			f.Term = "svg"
		}
		if f.Fop == 0{
			f.Opt = 0
		}
		f.Mtyp = 1
		f.Fixbase = true
		switch f.Opt{
			case 0:
			f.Verbose = true
			_, _, err = mosh.Frm2dDz(&f)
			
			if err != nil{
				terror.Execute(w, fmt.Errorf("frame design error->\n%s",err))
				return
			}
			err = trrez.Execute(w, &f)
			if err != nil{
				terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
				return
			}
			default:
			frez, err = mosh.Frm2dOpt(f) 
			if err != nil{
				terror.Execute(w, fmt.Errorf("frame design error->\n%s",err))
				return
			}
			err = trrez.Execute(w, &frez)
			if err != nil{
				terror.Execute(w, fmt.Errorf("major template error->\n%s",err))
				return
			}
		}
	}
}

//parsemosh parses/checks form data and returns the json string
func parsemosh(r *http.Request)(jsonstr string){
	//must be a better way
	//a better day :(
	var secstr string
	jsonstr += "{"
	strmap := map[string]int{
		"Json":-1,
		"Title":1,
		"Term":1,
		"Shape":1,
		"Typ":1,
		"Typstr":1,
		"Csec":2,
		"Bsec":2,
		"Styps":2,
		"Sections":2,
		"Lspans":2,
		"Bloads":2,
		"Lbays":2,
		"Lsxs":2,
		"Hs":2,
		"Fcks":2,
		"Fys":2,
		"Mys":2,
		"Mxs":2,
		"Pus":2,
		"Psfs":2,
		"PSFs":2,
		"WLFs":2,
		"Kostin":2,
		"X":2,
		"Y":2,
		"WL":2,
		"Dias":2,
		"Dbars":2,
		"Dims":2,
		"Coords":3,
	}
	if val, ok := r.Form["Json"]; ok && val[0] != ""{
		if json.Valid([]byte(val[0])){
			jsonstr = val[0]
			return
		} 
	}
	for key, val := range r.PostForm{
		//log.Println("processing key->",key," val->",val, "len->", len(val))
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
				if secstr == ""{
					secstr += "{\"Sec\":{"
				}
				secstr += fmt.Sprintf("\"%v\":%v,",key,val)
			}
		}
	}
	jsonstr = jsonstr[:len(jsonstr)-1]	
	if secstr != ""{
		secstr += "}"
		jsonstr += secstr
	}
	jsonstr += "}"
	fmt.Println(mosh.ColorRed,"jsonstr-",jsonstr,mosh.ColorReset)
	return
}

// func rcbeamstyp(w http.ResponseWriter, r *http.Request){
// 	if r.Method == "GET"{
// 		r.ParseForm()
// 		var hstr string
// 		for k, v := range r.Form{
// 			log.Println("key, val->",k,v)
// 			if k == "Dsgn"{
// 				switch v[0]{
// 					case "true":
// 					//design styps
// 					hstr = `
// 	<option value="1"  selected="selected">rectangle</option>
// 	<option value="6">T</option>
// 	<option value="7">L</option>
// 	<option value="14">T-pocket</option>`
// 					case "false":
// 					//all styps for analysis
// 					hstr = `
// 	<option value="1"  selected="selected">rectangle</option>
// 	<option value="4">box</option>
// 	<option value="6">T</option>
// 	<option value="7">L</option>
// 	<option value="12">I</option>
// 	<option value="13">C</option>
// 	<option value="14">T-pocket</option>
// 	<option value="19">tapered pocket</option>
// 	<option value="20">trapezoidal</option>
// `
// 				}
// 				fmt.Fprint(w, hstr)
// 			}
// 		}

// 	}
// }

//ye olde

// switch{
// 	case "Title","Term":	
// 	jsonstr += fmt.Sprintf("\"%v\":\"%s\",",key,val[0])
// 	case "Sections","Lspans","Lsxs","Hs","Fcks","Fys","Csec","Bsec","Mys","Mxs","Pus","X","Y","WL":
// 	jsonstr += fmt.Sprintf("\"%v\":%v,",key,val)
// 	default:
// 	jsonstr += fmt.Sprintf("\"%v\":%v,",key,val[0])
// }
