package barf

import(
	"fmt"
	"strings"
	"net/http"
	//"html/template"
)

func hdocs(w http.ResponseWriter, r *http.Request){
	var title, base, item string
	title = r.URL.Path[len("/docs/"):]
	// log.Println("looking up->",title)
	switch{
		case strings.Contains(title,"rcc"):
		item = r.URL.Path[len("/docs/rcc/"):]
		base = "rcc"
		switch item{
			case "slab":
			
		}
		//fmt.Println("looking up->",title, item)
		case strings.Contains(title,"calcmod"):
		item = r.URL.Path[len("/docs/calcmod/"):]
		base = "mod"
		case strings.Contains(title,"calcnp"):
		item = r.URL.Path[len("/docs/calcnp/"):]
		base = "modnp"
		case strings.Contains(title,"calcep"):
		item = r.URL.Path[len("/docs/calcep/"):]
		base = "modep"
		case strings.Contains(title,"steel"):
		case strings.Contains(title,"timber"):
	}
	tstr := fmt.Sprintf("%s%s.tmpl",base,item)
	tdc := tdocs.Lookup(tstr)
	if tdc == nil{
		terror.Execute(w, fmt.Errorf("%s - doc file not found",tstr))		
		return
	}
	//tdc = template.Must(tdc.Parse("srvr/templates/header.tmpl"))
	// tdc = template.Must(tdc.Parse("srvr/templates/footer.tmpl"))
	err := tdc.Execute(w, nil)
	if err != nil{
		terror.Execute(w, fmt.Errorf("%s - doc template execution error",err))		
	}
}
