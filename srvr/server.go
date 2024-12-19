package barf

import (
	"runtime"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

type Hmsg struct{
	Msg string
}

type Hpage struct{
	Title string
	Body  string
	Links []string
	Texts []string
	Prev  string
	Next  string
}

func about(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		err := tabout.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		var msgs string
		r.ParseForm()
		for key, val := range r.PostForm{
			log.Println("processing key->",key," val->",val)
			msgs += fmt.Sprintf("%s->\n%s\n",key,val[0])
		}
		hmsg := Hmsg{Msg:msgs}
		tmsg.Execute(w,hmsg)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	err := tindex.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func analyze(w http.ResponseWriter, r *http.Request) {
	err := tanalyze.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func hammeron(w http.ResponseWriter, r *http.Request){
// 	if r.Method == "GET"{
// 		hstr := `<img src="/svg/web/hammer.png">`
// 		fmt.Fprint(w, hstr)
// 	}
// }

func httpsRedir(w http.ResponseWriter, r *http.Request){
        http.Redirect(w, r, "https://barfcalc.in:443"+r.RequestURI,http.StatusMovedPermanently)

}

func sendrez(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	hstr := fmt.Sprintf(`<div "nes-container is-dark with-title" id="send-email">
  <p class="nes-text-is-success">email sent to %s with message %s</p>
</div>`,r.Form["Email"],r.Form["Message"])
	fmt.Fprint(w, hstr)
}


//Srvr serves pages. quite badly
func Srvr(){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	fs := http.FileServer(http.Dir(filepath.Join(basepath,"../data/out")))
	http.Handle("/svg/", http.StripPrefix("/svg/", fs))
	//index
	http.HandleFunc("/",index)
	//kass funcs
	http.HandleFunc("/analyze",analyze)
	http.HandleFunc("/analyze/basic", calcmod)
	http.HandleFunc("/analyze/np", calcnp)
	http.HandleFunc("/analyze/ep", calcep)
	//kass htmx funcs
	http.HandleFunc("/ex/models/basic",modhtml)
	http.HandleFunc("/ex/models/np",modnphtml)
	http.HandleFunc("/ex/models/ep",modephtml)
	
	//mosh funcs
	http.HandleFunc("/rcc",rcc)
	http.HandleFunc("/rcc/slab",calcrslab)	
		
	http.HandleFunc("/rcc/beam",rccbeam)
	http.HandleFunc("/rcc/beam/sec",calcrbeam)
	http.HandleFunc("/rcc/beam/ssbeam",calcrcbeam)
	http.HandleFunc("/rcc/beam/cbeam",calcrcbeam)
	
	http.HandleFunc("/rcc/column",calcrcol)
	http.HandleFunc("/rcc/footing",calcrftng)
	http.HandleFunc("/rcc/subframe",calcrsubfrm)
	http.HandleFunc("/rcc/frame2d",calcrfrm2d)

	
	//mosh htmx funcs
	http.HandleFunc("/ex/rccslab",rslabhtml)
	http.HandleFunc("/ex/rcc/ftng",rcftnghtml)
	http.HandleFunc("/ex/rcc/frm2d",rcf2dhtml)
	http.HandleFunc("/ex/rcc/bmcs",rcbmcshtml)
	http.HandleFunc("/ex/rcc/bmss",rcbmsshtml)
	http.HandleFunc("/ex/rcc/bmsec",rcbmsechtml)
	http.HandleFunc("/ex/rcc/bmsec/styp",rcbmstyphtml)
	http.HandleFunc("/ex/rcc/col",rccolhtml)
	http.HandleFunc("/ex/rcc/col/styp",rccolstyphtml)
	http.HandleFunc("/ex/rcc/sf",rcsfhtml)
	
	
	//steel funcs
	http.HandleFunc("/steel",steel)
	http.HandleFunc("/steel/beam",stlbeam)
	http.HandleFunc("/steel/column",stlcol)
	http.HandleFunc("/steel/column/frame",stlcolfrm)
	http.HandleFunc("/steel/column/strut",stlcolstrt)
	http.HandleFunc("/steel/truss",stltrs)
	http.HandleFunc("/steel/trussmodopt",stltrsmodopt)
	http.HandleFunc("/steel/trussgen",calcstltrsgen)
	
	//bash htmx funcs
	
	http.HandleFunc("/ex/steel/beam",stlbmhtml)
	http.HandleFunc("/ex/steel/col",stlcolstrthtml)
	http.HandleFunc("/ex/steel/col/frm",stlcolfrmhtml)
	http.HandleFunc("/ex/steel/trussmodopt",stltrsmodopthtml)
	http.HandleFunc("/ex/steel/trussgen",stltrsgenhtml)
	
	//timber funcs
	http.HandleFunc("/timber",timber)
	http.HandleFunc("/timber/beam",calctmbrbeam)
	http.HandleFunc("/timber/column",calctmbrcol)
	
	//tmbr htmx funcs
	http.HandleFunc("/ex/timber/beam",tmbrbmhtml)
	http.HandleFunc("/ex/timber/col",tmbrcolhtml)
	
	//about  
	http.HandleFunc("/about",about)
	
	//add help method
	http.HandleFunc("/docs/",hdocs)

	//email w/results
	
	//http.HandleFunc("/email",sendrez)

	
	//for localhost
	//print stuff to show seriousness
	log.Printf("Starting server at port 8080\n")
	http.ListenAndServe(":8080", nil)

	//for domain
	// log.Printf("Starting (tls) server at port 443\n")
	// go http.ListenAndServeTLS(":443","srvr/ssl/certificate.crt","srvr/ssl/certificate.key",nil)
	// http.ListenAndServe(":80",http.HandlerFunc(httpsRedir))
}


/*
		// switch key {
		// case "Id":	
		// 	jsonstr += fmt.Sprintf("\"%v\":\"%s\",",key,val[0])
		// case "Term":
		// 	jsonstr += fmt.Sprintf("\"%v\":\"%s\",",key,val[0])
		// case "Units":
		// 	jsonstr += fmt.Sprintf("\"%v\":\"%s\",",key,val[0])
		// case "Frmstr":
		// 	jsonstr += fmt.Sprintf("\"%v\":\"%s\",",key,val[0])
		// case "Cmdz":
		// 	jsonstr += fmt.Sprintf("\"%v\":%v,",key,val)
		// default:
		// 	jsonstr += fmt.Sprintf("\"%v\":%v,",key,val)
		// }


*/
