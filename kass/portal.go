package barf

import (
	"fmt"
	"log"
	"errors"
	// "math"
	// "path/filepath"
	// "runtime"
	// "os"
	// "github.com/go-gota/gota/dataframe"
	//"github.com/go-gota/gota/series"
)


//Portal is a struct to hold portal frame input fields
type Portal struct {
	//2d portal frame
	//Pz - wind pressure 0.6 * vz2
	//pzcs - wind load coeffs
	Title                        string
	Sname                        string
	Nbays, Nframes               int
	Span, Spacing, Slope, Height float64
	Pz, DL, LL                   float64
	Dh, Lh                       float64
	Wr, Wc, Wudr, Wudc, W        float64
	Pzcs                         []float64   //wind pressure coefficients
	Cp                           [][]float64 //if full, col/beam/brace
	Em                           [][]float64 //if full, col/beam/brace
	Sections                     [][]float64
	Hdims                        [][]float64 //bm haunch(lh, db, de); col haunch lh, db, de. if Config = 3; db, de
	Styps                        []int
	Sdxs                         []int
	Ldgen                        bool
	Selfwt                       bool
	Fixbase, Haunch, Gable, Mono bool
	Verbose                      bool
	Npsec                        bool
	Readsec                      bool
	Npmod                        bool //if true, generate memnp dims
	Code                         int
	Csec, Bsec                   int
	Config                       int //0 - uniform, 1- haunched rafter, 2 -haunched col/rafter, 3-uniform tapered
	Bstyp                        int //base section type
	Sectyp                       int
	Gentyp                       int //0 - single element per member, 1 - split mem by ndiv
	Opt                          int
	Ncols, Cdiv                  int
	Nbms, Bdiv                   int
	Ndiv                         int
	Nsecs                        int
	Grade                        float64
	Dmin                         float64
	Params                       []float64 //generic params float slice
	PSFs                         []float64
	Matprop                      []float64 //em
	Inpdims                      []float64 //generic input dims float slice
	Lx, Ly                       float64 //plan dimensions
	Vz, Cpi                      float64
	Prlnspc, Prlnwt              float64
	LR, Rise                     float64
	Nprlz, Ngrtz                 float64
	Grtspc, Grtwt                float64
	Wt, Kost                     float64
	Prlnsec                      int
	Term                         string
	Mod                          Model
	Mtyp                         int
	Prlntyp                      int
	Ldcalc                       int
	Ltyp                         int //load application type
	Cols, Beams                  []int
	Kbs                          []int
	Blnodes                      map[int]float64 //beam load nodes
	Clnodes                      map[int]float64 //col load nodes
	Nodes                        map[Pt][]int `json:"-"`
	Members                      map[int][][]int `json:",omitempty"`
	Advloads                     map[int][][]float64 `json:",omitempty"`
	Benloads                     map[int][][]float64 `json:",omitempty"`
	Loadcons                     map[int][][]float64 `json:",omitempty"`
	Uniloads                     map[int][][]float64 `json:",omitempty"`
	Jloadmap                     map[int][][]float64 `json:",omitempty"`
	Wadloads                     map[int][][]float64 `json:",omitempty"`

}

//Init initializes a portal frame
func (f *Portal) Init() (err error){
	//var x, y, fdl, fll, fulg, bmwt, colwt, lspan float64
	var lspan float64
	// f.Nframes = int(f.Ly/f.Spacing)
	// f.Spacing = f.Ly/float64(f.Nframes)
	// f.Nframes += 1
	if len(f.Sections) == 0{
		if len(f.Sdxs) == 0{
			err = errors.New("no sections specified")
			return
		}
		f.Readsec = true
	}
	if f.Mono{lspan = f.Span} else {lspan = f.Span/2.0}
	if f.Mtyp == 0{f.Mtyp = 2}
	if f.Verbose{
		log.Println("starting portal frame analysis of->",f.Title)
		log.Println("number of frames->",f.Nframes)
		log.Println("lspan->",lspan)
		log.Println("nbays->",f.Nbays)
	}
	if f.Rise == 0.0{
		if f.Slope > 0.0{
			f.Rise = lspan * f.Slope
		}
	}
	if f.Slope == 0.0{
		f.Slope = f.Rise/lspan
	}
	f.Nodes = make(map[Pt][]int)
	f.Members = make(map[int][][]int)
	f.Advloads = make(map[int][][]float64)
	f.Benloads = make(map[int][][]float64)
	f.Loadcons = make(map[int][][]float64)
	f.Uniloads = make(map[int][][]float64)
	f.Jloadmap = make(map[int][][]float64)
	f.Wadloads = make(map[int][][]float64)
	f.Blnodes = make(map[int]float64)
	f.Clnodes = make(map[int]float64)
	
	if f.Config > 0{
		f.Npsec = true
	}
	if f.Gentyp == 2 && f.Ndiv == 0{
		f.Ndiv = 5
	}
	switch f.Mtyp{
		case 1:
		//rcc portal frame
		case 2:
		//noble steel portal frame
		if f.Code == 0{f.Code = 2}
		if f.Grade == 0{f.Grade = 43.0}
		if f.Sectyp == 0{
			f.Sectyp = 12
			f.Sname = "i"
		}
		if f.Bstyp == 0{f.Bstyp = 12}
		switch f.Code{
			case 1:
			f.Mod.Em = [][]float64{{2.1e5,0.3}}
			case 2:
			f.Mod.Em = [][]float64{{2.0e5,0.3}}
		}
		f.Mod.Frmtyp = 3
		f.Mod.Units = "n-mm"
		f.Mod.Frmstr = "2df"
		f.Mod.Submems = make(map[int][]int)
		if f.Config == 2{f.Mod.Split = true}
		if f.Readsec{
			err = f.ReadStlDims()
			if err != nil{
				return
			}
			if f.Verbose{fmt.Println("read dims->",f.Sections)}
		}
		if f.Verbose{log.Println("setting mod.em->",f.Mod.Em," n/mm2")}
		
		case 3:
		//timber portal frame?
		case 4:
		//cfs?(see phan paper)

	}
	return
}

//Draw plots a portal frame
func (f *Portal) Draw(){
	pltchn := make(chan string)
	switch f.Config{
		case 0:
		go PlotFrm2d(&f.Mod, f.Term, pltchn)
		default:
		go PlotNpFrm2d(&f.Mod, f.Term, pltchn)
	}
	<- pltchn
	PlotGenTrs(f.Mod.Coords, f.Mod.Mprp)
}

//Getndim generates ndims for portal frame optimization
func (f *Portal) Getndim() (nd int){
	switch f.Config{
		case 0:
		//col, beam
		nd = 2
		case 1:
		//col, beam, lh, dh
		nd = 4
		case 2:
		//col, beam, lh, dh, lh col, dh col 
		nd = 6
		case 3:
		//dc strt, end, db end
		nd = 3
	}
	return
}

//GenLoads generates loads and load combos for a portal frame
func (f *Portal) GenLoads(){
	//one with just load cases
	switch{
		case f.W > 0.0:
		//point loads
		switch f.Ltyp{
			case 0:
			//add point member loads
			var ldcase []float64
			for _, bm := range f.Beams{
				jb := f.Members[bm][0][0]
				je := f.Members[bm][0][1]
				coords := SplitCoords(f.Mod.Coords[jb-1], f.Mod.Coords[je-1], f.Bdiv)
				ldcase = []float64{0.0, f.W/2.0, 0.0}
				f.AddNodalLoad(jb, ldcase)
				f.AddNodalLoad(je, ldcase)
				for i, pt := range coords{
					if i == 0 || i == len(coords)-1{
						continue
					}
					la := Dist3d(f.Mod.Coords[jb-1],pt)
					ldcase = []float64{1.0, f.W, 0, la, 0, 1}
					f.AddMemLoad(bm, ldcase)
				}
			}
			case 1:
			//add nodal loads at end points (W/2)
			var ldcase []float64
			for n, lf := range f.Blnodes{
				ldcase = []float64{0.0, -f.W*lf, 1.0}
				f.AddNodalLoad(n, ldcase)
			}
		}
	}
}

//AddNodalLoad adds a nodal load to a node n of a portal frame
//is this necessary?not at all
func (f *Portal) AddNodalLoad(n int, ldcase []float64)(err error){
	ldvec := append([]float64{float64(n)},ldcase...)
	if f.Verbose{log.Println("adding load to node->",n, "ldcase->",ldcase)}
	f.Mod.Jloads = append(f.Mod.Jloads, ldvec)
	return
}

//AddMemLoad adds a member load to a mem (mdx) of a portal frame
func (f *Portal) AddMemLoad(mem int, ldcase []float64) (err error){
	if _, ok := f.Members[mem]; !ok{
		return errors.New("invalid member index")
	}
	ldvec := append([]float64{float64(mem)},ldcase...)
	f.Mod.Msloads = append(f.Mod.Msloads, ldvec)
	return
}

//CalcLoadEnv calculates load envelopes for a portal frame
func (f *Portal) CalcLoadEnv()(err error){
	//now it just calcs model
	f.Mod.Id = f.Title
	f.Mod.Units = "nmm"
	f.Mod.Frmstr = "2df"
	//fmt.Println("here be joint loads->",f.Mod.Jloads)
	switch f.Npmod{
		case false:
		_, err = CalcFrm2d(&f.Mod, 3)
		case true:
		_, err = CalcNp(&f.Mod, "2df",true)
	}
	if err != nil{return}
	fmt.Println("analysis finito (yee-haw)")
	return
}

//GenBasic generates a portal frame with a single div per member
//if config = 0 - basic mod else mod np
//use gensplit for all other configs
func (f *Portal) GenBasic(){
	var x, y float64
	switch f.Nbays{
		case 1:
		default:
		switch f.Mono{
			case true:
			return
		}
	}
	for i := 0; i < f.Nbays; i++{
		switch i{
			case 0:
			x += f.Span; y += f.Height
			switch f.Mono{
				case true:
				p1 := []float64{0,0}
				p2 := []float64{x,y-f.Height}
				p3 := []float64{x-f.Span,y}
				p4 := []float64{x,y+f.Rise}
				for _, pt := range [][]float64{p1, p2, p3, p4}{
					_ = f.AddNode(pt)
				}
				//add cols
				cns := [][]int{{1,3},{2,4}}
				for i, cn := range cns{
					switch f.Config{
						case 0:
						_ = f.AddMem(cn[0],cn[1],cn[0],cn[1],i+1,i+1)
						default:
						_ = f.AddMemNp(cn[0],cn[1],cn[0],cn[1],i+1,i+1)
					}
				}
				//add beam
				switch f.Config{
					case 0:
					_ = f.AddMem(3,4,3,4,3,3)
					default:
					_ = f.AddMemNp(3,4,3,4,3,3)

				}
				if f.Fixbase{
					f.Mod.Supports = append(f.Mod.Supports, []int{1,-1,-1,-1})
					f.Mod.Supports = append(f.Mod.Supports, []int{2,-1,-1,-1})
				} else {
					f.Mod.Supports = append(f.Mod.Supports, []int{1,-1,-1,0})
					f.Mod.Supports = append(f.Mod.Supports, []int{2,-1,-1,0})
				}
				case false:

				p1 :=  []float64{0,0}
				p2 :=  []float64{x,y-f.Height}
				p3 :=  []float64{x-f.Span,y}
				p4 :=  []float64{x,y}
				p5 :=  []float64{x/2.0,y+f.Rise}
				for _, pt := range [][]float64{p1, p2, p3, p4, p5}{
					_ = f.AddNode(pt)
				}
				cns := [][]int{{1,3},{2,4}}
				for j, cn := range cns{
					switch f.Config{
						case 0:
						_ = f.AddMem(cn[0],cn[1],cn[0],cn[1],j+1,j+1)
						default:
						_ = f.AddMemNp(cn[0],cn[1],cn[0],cn[1],j+1,j+1)
					}
				}
				bns := [][]int{{3,5},{4,5}}
				for j, bn := range bns{
					switch f.Config{
						case 0:
						_ = f.AddMem(bn[0],bn[1],bn[0],bn[1],j+3,j+3)
						default:
						_ = f.AddMemNp(bn[0],bn[1],bn[0],bn[1],j+3,j+3)
					}
				}
				if f.Fixbase{
					f.Mod.Supports = append(f.Mod.Supports, []int{1,-1,-1,-1})
					f.Mod.Supports = append(f.Mod.Supports, []int{2,-1,-1,-1})
				} else {
					f.Mod.Supports = append(f.Mod.Supports, []int{1,-1,-1,0})
					f.Mod.Supports = append(f.Mod.Supports, []int{2,-1,-1,0})
				}
			}
			default:
			switch f.Mono{
				case true:
				//FIGURE THIS. col r should be split at rise
				case false:
				x += f.Span
				p1 := []float64{x,y-f.Height}
				p2 := []float64{x, y}
				p3 := []float64{x-f.Span/2.0,y+f.Span*f.Slope/2.0}
				for _, pt := range [][]float64{p1, p2, p3}{
					_ = f.AddNode(pt)
				}
				idx := 6 + (i-1)*3
				//col, beams
				switch f.Config{
					case 0:
					//col
					_ = f.AddMem(idx, idx+1,idx, idx+1, 2, 2)
					//beams
					_ = f.AddMem(idx-2, idx+2,idx-2, idx+2, 3, 3)
					_ = f.AddMem(idx+1, idx+2,idx+1, idx+2, 4, 4)
					default:
					//col
					_ = f.AddMemNp(idx, idx+1,idx, idx+1, 2, 2)
					//beams
					_ = f.AddMemNp(idx-2, idx+2,idx-2, idx+2, 3, 3)
					_ = f.AddMemNp(idx+1, idx+2,idx+1, idx+2, 4, 4)
				}
				if f.Fixbase{
					f.Mod.Supports = append(f.Mod.Supports, []int{idx, -1,-1,-1})
				} else {
					f.Mod.Supports = append(f.Mod.Supports, []int{idx, -1,-1, 0})
				}
			}
		}
	}
	// log.Printf("done geom - %v\n",f.Mod)
}

//GenSplit generates a portal frame of gentyp 1 (b/cdiv per member) and gentyp 2 (b/cdiv * ndiv per member)
func (f *Portal) GenSplit(){
	var x, y float64
	var bmem, ns1, ns2 int
	for i := 0; i < f.Nbays; i++{
		switch i{
			case 0:
			x += f.Span; y += f.Height
			switch f.Mono{
				case false:
				//base coords
				p1 :=[]float64{0,0}
				p2 :=[]float64{x,y-f.Height}
				p3 :=[]float64{x-f.Span,y}
				p4 :=[]float64{x,y}
				p5 :=[]float64{x/2.0,y+f.Rise}
				var ndx []int
				
				//c1
				bmem++
				switch f.Gentyp{
					case 1:
					ndx = f.SplitCoords(p1, p3, f.Cdiv)
					for i, n := range ndx{
						if i < len(ndx) -1{
							f.AddMem(n, ndx[i+1],n, ndx[i+1],1, bmem)
						}
						switch i{
							case 0, len(ndx)-1:
							f.Clnodes[n] = 0.5
							default:
							f.Clnodes[n] = 1.0
						}
					}
					ns1 = ndx[0]
					case 2:
					//split mems by Ndiv
					coords := SplitCoords(p1, p3, f.Cdiv)
					for j, pt := range coords{
						if j < len(coords)-1{
							ndx = f.SplitCoords(pt, coords[j+1], f.Ndiv)
							for i, n := range ndx{
								if j == 0 && i == 0{
									ns1 = n
								}
								switch i{
									case 0, len(ndx)-1:
									switch j{
										case 0, len(coords)-1:
										f.Clnodes[n] = 0.5
										default:
										f.Clnodes[n] = 1.0
									}
								}
								if i < len(ndx) -1{
									f.AddMem(n, ndx[i+1], n, ndx[i+1],1, bmem)
								}
							}
						}
					}
				}
				//b1
				bmem++
				switch f.Gentyp{
					case 1:	
					ndx = f.SplitCoords(p3, p5, f.Bdiv)
					for i, n := range ndx{
						if i < len(ndx) -1{
							f.AddMem(n, ndx[i+1],n, ndx[i+1],3, bmem)
						}
						
						switch i{
							case 0:
							f.Blnodes[n] = 0.5
							default:
							f.Blnodes[n] = 1.0
						}
					}
					case 2:
					//split beams by bdiv and ndiv
					coords := SplitCoords(p3, p5, f.Bdiv)
					for j, pt := range coords{
						if j < len(coords)-1{
							ndx = f.SplitCoords(pt, coords[j+1], f.Ndiv)
							for i, n := range ndx{
								switch i{
									case 0, len(ndx)-1:
									switch j{
										case 0:
										f.Blnodes[n] = 0.5
										default:
										f.Blnodes[n] = 1.0
									}
								}	
								if i < len(ndx) -1{
									f.AddMem(n, ndx[i+1], n, ndx[i+1],1, bmem)
								}
							}
						}
					}
				}
				//b2
				bmem++
				switch f.Gentyp{
					case 1:
					ndx = f.SplitCoords(p5, p4, f.Bdiv)
					for i, n := range ndx{
						if i < len(ndx) -1{
							f.AddMem(n, ndx[i+1],n, ndx[i+1],4, bmem)
						}
						switch i{
							case len(ndx)-1:
							f.Blnodes[n] = 0.5
							default:
							f.Blnodes[n] = 1.0
						}
					}
					case 2:
					//split beams by bdiv and ndiv
					coords := SplitCoords(p5, p4, f.Bdiv)
					for j, pt := range coords{
						if j < len(coords)-1{
							ndx = f.SplitCoords(pt, coords[j+1], f.Ndiv)
							for i, n := range ndx{
								switch i{
									case 0:
									f.Blnodes[n] = 1.0
									case len(ndx)-1:
									switch j{
										case len(coords)-2:
										f.Blnodes[n] = 0.5
										default:
										f.Blnodes[n] = 1.0
									}
								}	
								if i < len(ndx) -1{
									f.AddMem(n, ndx[i+1], n, ndx[i+1],1, bmem)
								}
							}
						}
					}
				}
				//c2
				bmem++
				switch f.Gentyp{
					case 1:
					ndx = f.SplitCoords(p4, p2, f.Cdiv)
					for i, n := range ndx{
						if i < len(ndx) -1{
							f.AddMem(n, ndx[i+1],n, ndx[i+1],2, bmem)
						}
						
						switch i{
							case 0, len(ndx)-1:
							f.Clnodes[n] = 0.5
							default:
							f.Clnodes[n] = 1.0
						}
					}
					ns2 = ndx[len(ndx)-1]
					case 2:
					//split mems by cdiv/ndiv
					coords := SplitCoords(p4, p2, f.Cdiv)
					for j, pt := range coords{
						if j < len(coords)-1{
							ndx = f.SplitCoords(pt, coords[j+1], f.Ndiv)
							for i, n := range ndx{
								switch i{
									case 0, len(ndx)-1:
									switch j{
										case 0, len(coords)-1:
										f.Clnodes[n] = 0.5
										default:
										f.Clnodes[n] = 1.0
									}
								}
								if i < len(ndx) -1{
									f.AddMem(n, ndx[i+1], n, ndx[i+1],1, bmem)
								}
								if j == len(coords)-2 && i == len(ndx)-1{
									ns2 = n
								}
								
							}
						}
					}
				}
				if f.Fixbase{
					f.Mod.Supports = append(f.Mod.Supports, []int{ns1,-1,-1,-1})
					f.Mod.Supports = append(f.Mod.Supports, []int{ns2,-1,-1,-1})
				} else {
					f.Mod.Supports = append(f.Mod.Supports, []int{ns1,-1,-1,0})
					f.Mod.Supports = append(f.Mod.Supports, []int{ns2,-1,-1,0})
				}
			}
			default:
			x += f.Span
			switch f.Mono{
				case false:
				//base coords
				p1 :=[]float64{x-f.Span,y}
				p2 :=[]float64{x-f.Span/2.0,y+f.Rise}
				p3 :=[]float64{x,y}
				p4 :=[]float64{x,y-f.Height}
				var ndx []int
				//b1
				bmem++
				switch f.Gentyp{
					case 1:
					
					ndx = f.SplitCoords(p1,p2,f.Bdiv)
					for i, n := range ndx{
						if i < len(ndx) -1{
							f.AddMem(n, ndx[i+1],n, ndx[i+1],3,bmem)
						}
					}
					case 2:
					
				}
				//b2
				bmem++
				switch f.Gentyp{
					case 1:	
					ndx = f.SplitCoords(p2, p3, f.Bdiv)
					for i, n := range ndx{
						if i < len(ndx) -1{
							f.AddMem(n, ndx[i+1],n, ndx[i+1],4,bmem)
						}
					}
					case 2:
					
				}
				//c2
				bmem++
				switch f.Gentyp{
					case 1:
					ndx = f.SplitCoords(p3, p4, f.Cdiv)
					for i, n := range ndx{
						if i < len(ndx) -1{
							f.AddMem(n, ndx[i+1],n, ndx[i+1],2,bmem)
						}
					}
					ns2 = ndx[len(ndx)-1]
					case 2:
				}
				if f.Fixbase{
					f.Mod.Supports = append(f.Mod.Supports, []int{ns2,-1,-1,-1})
				} else {
					f.Mod.Supports = append(f.Mod.Supports, []int{ns2,-1,-1,0})
				}
			}
		}
	}
}

//AddNode adds a node at pt and returns the index of the node
func (f *Portal) AddNode(pt []float64)(ndx int){
	p := Pt{pt[0],pt[1]}
	if val, ok := f.Nodes[p]; ok{
		ndx = val[0]
	} else {
		f.Mod.Coords = append(f.Mod.Coords, pt)
		f.Nodes[p] = []int{len(f.Mod.Coords)}
		ndx = len(f.Mod.Coords)
	}
	return
}

//ReadDims reads in dx col/dx beam from sdxs and updates f.Sections
func (f *Portal) ReadStlDims() (err error){
	f.Sections = [][]float64{}
	var dims []float64
	var bstyp int
	for _, sdx := range f.Sdxs{
		
		dims, bstyp, err = ReadSecDims(f.Sname, sdx, 1)
		if err != nil{
			log.Println("error reading stl dims",err)
			return
		}
		f.Sections = append(f.Sections, dims)
	}
	f.Bstyp = bstyp
	log.Println("done reading dims-", f.Sections, f.Bstyp)
	return
}

//ReadNpDims returns npmem params
func (f *Portal) ReadNpDims(jb, je, n1, n2, mdx, bmem int) (ts []int, ls, bs, ds, dims []float64, d1, d2 float64){
	p1 := f.Mod.Coords[jb-1]
	p2 := f.Mod.Coords[je-1]
	lspan := Dist3d(p1,p2)
	var db, de, dt, lh, lm, lx float64
	mtyp := 1
	if mdx > 2{
		mtyp = 2
	}
	dims = f.Sections[mtyp-1]
	b := dims[0]
	d := dims[1]
	bs = []float64{b,b,b}
	
	switch f.Config{
		case 1:
		//haunched rafters
		db = f.Hdims[0][0] + d
		lh = f.Hdims[0][1]
		lm = lspan - lh
		de = d
		switch mdx{
			case 3,4:
			ts = []int{2,0,0}
			ds = []float64{db, de, de}
			ls = []float64{lh, lm, 0}
			lx = Dist3d(p1, f.Mod.Coords[n1-1])
			if lx > 0.0{
				d1 = GetNpDepth(ts, ls, ds, lx)
			}
			lx = Dist3d(p1, f.Mod.Coords[n2-1])
			if lx > 0.0{
				d2 = GetNpDepth(ts, ls, ds, lx)
			}
			default:
			ts = []int{0,0,0}
			ds = []float64{de, de, de}
			ls = []float64{lspan, 0, 0}
			d1 = 0.0
			d2 = 0.0

		}
		case 2:
		//haunched everything
		db = f.Hdims[mtyp-1][0]
		lh = f.Hdims[mtyp-1][1]
		lm = lspan - lh
		switch mdx{
			case 3,4:
			ts = []int{2,0,0}
			ds = []float64{db, de, de}
			ls = []float64{lh, lm, 0}
			default:
			ts = []int{0,0,2}
			ds = []float64{de, de, db}
			ls = []float64{0, lm, lh}
		}
		lx = Dist3d(p1, f.Mod.Coords[n1-1])
		if lx > 0.0{
			d1 = GetNpDepth(ts, ls, ds, lx)
		}
		lx = Dist3d(p1, f.Mod.Coords[n2-1])
		if lx > 0.0{
			d2 = GetNpDepth(ts, ls, ds, lx)
		}
		case 3:
		//tapered
		db = f.Hdims[0][0]
		de = f.Hdims[0][1]
		dt = f.Hdims[0][2]
		ts = []int{2,0,0}
		ls = []float64{lspan, 0, 0}
		switch mdx{
			case 3,4:
			ds = []float64{de, dt, dt}
			default:
			ds = []float64{db, de, de}
		}
		lx = Dist3d(p1, f.Mod.Coords[n1-1])
		if lx > 0.0{
			d1 = GetNpDepth(ts, ls, ds, lx)
		}
		lx = Dist3d(p1, f.Mod.Coords[n2-1])
		if lx > 0.0{
			d2 = GetNpDepth(ts, ls, ds, lx)
		}
	}
	switch f.Gentyp{
		case 0:
		default:
		//split tis mem
		j1 := f.Mod.Coords[n1-1]
		j2 := f.Mod.Coords[n2-1]
		lsmol := Dist3d(j1, j2)
		ts = []int{2,0,0}
		ls = []float64{lsmol, 0, 0}
		ds = []float64{d1, d2, d2}
	}
	return
}

//AddMemNp adds an npmem bet n1 and n2 with basemem bmem at jb/je
func (f *Portal) AddMemNp(jb, je, n1, n2, mdx, bmem int)(err error){
	cp := len(f.Mod.Mprp)+1
	mvec := []int{n1, n2, 1, cp, 0, cp}
	mtyp := 1
	if mdx > 2{
		mtyp = 2
	}
	switch f.Gentyp{
	//1 - uniform member , 1 - single (np) member, 2 - tapered member
	//
		case -1:
		//single uniform member
		case 0:
		//single mem
		ts, ls, bs, ds, dims, _, _ := f.ReadNpDims(jb, je, jb, je, mdx, bmem)
		f.Mod.Bs = append(f.Mod.Bs, bs)
		f.Mod.Ls = append(f.Mod.Ls, ls)
		f.Mod.Ts = append(f.Mod.Ts, ts)
		f.Mod.Ds = append(f.Mod.Ds, ds)
		f.Mod.Dims = append(f.Mod.Dims, dims)
		f.Mod.Sts = append(f.Mod.Sts, f.Bstyp)
		f.Mod.Mprp = append(f.Mod.Mprp, mvec)
		fmt.Println("added mem->",ts, ls, bs, ds, dims, f.Members[cp])
		case 1:
		//all are tapered mems
		ts, ls, bs, ds, dims, _, _ := f.ReadNpDims(jb, je, jb, je, mdx, bmem)
		f.Mod.Bs = append(f.Mod.Bs, bs)
		f.Mod.Ls = append(f.Mod.Ls, ls)
		f.Mod.Ts = append(f.Mod.Ts, ts)
		f.Mod.Ds = append(f.Mod.Ds, ds)
		f.Mod.Dims = append(f.Mod.Dims, dims)
		f.Mod.Sts = append(f.Mod.Sts, f.Bstyp)
		f.Mod.Mprp = append(f.Mod.Mprp, mvec)
		fmt.Println("added mem->",ts, ls, bs, ds, dims, f.Members[cp])
	}
	f.Members[cp] = append(f.Members[cp], mvec)
	f.Members[cp] = append(f.Members[cp],[]int{mtyp, bmem, mdx})
	return
}

//AddMem adds a mem between node n1 and n2
func (f *Portal) AddMem(jb, je, n1, n2, mdx, bmem int)(err error){
	var mvec []int
	var mtyp int
	var dims, cpvec []float64
	cp := len(f.Mod.Mprp) + 1
	mvec = []int{n1, n2, 1, cp, 0}
	f.Mod.Mprp = append(f.Mod.Mprp, mvec)
	switch mdx{
		case 1,2:
		//col end, col int
		f.Cols = append(f.Cols, cp)
		mtyp = 1
		f.Members[cp] = append(f.Members[cp], mvec)
		f.Members[cp] = append(f.Members[cp],[]int{mtyp, bmem, mdx})
		case 3,4:
		//beaml, beam r
		mtyp = 2
		f.Beams = append(f.Beams, cp)
		f.Members[cp] = append(f.Members[cp], mvec)
		f.Members[cp] = append(f.Members[cp],[]int{mtyp, bmem, mdx})
	}
	switch f.Mtyp{
		case 2:
		//steel
		switch f.Npsec{
			case true:
			//get depth at n1 and n2
			_, _, _, _, dims, d1, d2 := f.ReadNpDims(jb, je, n1, n2, mdx, bmem)
			//calculate ar, ix at n1 and n2
			ar1, ix1, _ := PropNpBm(f.Bstyp, dims[0], d1, dims)
			ar2, ix2, _ := PropNpBm(f.Bstyp, dims[0], d2, dims)
			//avg ar, ix and return
			ar := (ar1 + ar2)/2.0
			ix := (ix1 + ix2)/2.0
			dav := (d1 + d2)/2.0
			dims[1] = dav
			cpvec = []float64{ar, ix}
			case false:
			//basic uniform member with sdx
			dims = f.Sections[mtyp-1]
			cpvec, err = GetStlCp(f.Sname, f.Sectyp, f.Sdxs[mtyp-1], 1)
			if err != nil{
				//fmt.Println("errore",err)
				return
			}
		}
		f.Mod.Dims = append(f.Mod.Dims, dims)
		f.Mod.Cp = append(f.Mod.Cp, cpvec)
	}
	switch f.Config{
		case 2:
		if _, ok := f.Mod.Submems[bmem]; !ok{
			f.Mod.Submems[bmem] = []int{}
		}
		f.Mod.Submems[bmem] = append(f.Mod.Submems[bmem],cp)
	}
	if f.Verbose{
		//fmt.Println("input mem->jb, je, n1, n2, mdx, bmem ->",jb, je, n1, n2, mdx, bmem)
		//fmt.Println("added mem->dims, mprp->",f.Mod.Dims[cp-1],f.Mod.Mprp[cp-1])
	}
	return
}

//SplitCoords splits coords between pb and pe and returns node indices
func (f *Portal) SplitCoords(pb, pe []float64, ndiv int) (ndx []int){
	coords := SplitCoords(pb, pe, ndiv)
	for _, pt := range coords{
		n := f.AddNode(pt)
		ndx = append(ndx, n)
	}
	return
}


//GenGeom generates portal frame geometry
func (f *Portal) GenGeom(){
	//one uniform column and beam section
	// var x, y float64
	// var ldx, rdx int
	// ldx = 2; rdx = 2
	//config, gentyp, type
	switch f.Gentyp{
		case 0:
		f.GenBasic()
		case 1, 2:
		f.GenSplit()
	}
	//fmt.Println("done,",f.Mod.Mprp, f.Mod.Coords, f.Mod.Cp)
	log.Println("done geom - .", "supports",f.Mod.Supports)
}

//GenHaunchR generates a haunched rafter portal frame
func (f *Portal) GenHaunchR(){
}

func (f *Portal) InitMemRez()(err error){
	return
}

//Calc analyzes a portal frame
func (f *Portal) Calc()(err error){
	//fmt.Println("hyarr goeth nuthin")
	f.Init()
	f.GenGeom()
	fmt.Println("hyar be mprp")
	fmt.Println(f.Mod.Mprp)
	fmt.Println(f.Mod.Coords)
	fmt.Println(f.Mod.Cp)
	f.GenLoads()
	fmt.Println("model loads->",f.Mod.Msloads,"\n jloads->",f.Mod.Jloads)
	f.CalcLoadEnv()
	f.Draw()
	// fmt.Printf("%+v\n", f)
	//fmt.Println(f.Mod.Reports[0])
	return
	//f.GenGeom()
}

/*
// //Add member cp prop to model
// func (f *Portal) AddMemDims(mtyp, mem, cpdx int){
// 	//idx ==
// 	switch f.Mtyp{
// 		case 1:
// 		//rcc
// 		case 2:
// 		//honorable steel
// 		//dims := ReadSecDims(f.Sectyp,f.Sdxs[mtyp-1],1)
// 		f.Mod.Dims = append(f.Mod.Dims, dims)
// 		f.Mod.Sts = append(f.Mod.Sts, f.Bstyp)
// 		case 3:
// 		//wood sound good sound
// 	}
// 	switch cpdx{
// 		case 0:
// 		default:
// 	}
// }

// //GenGeom generates portal frame geometry and coords
// func (f *Portal) GenGeom(){
// 	switch f.Config{
// 		case 0:
// 		//uniform portal frame
// 		f.GenUniform()
// 		//done
// 		case 1:
// 		//haunched rafter portal frame
// 		f.GenUniform()
// 		case 2:
// 		//haunched col/rafter portal frame
// 		f.GenHaunchA()
// 		case 3:
// 		//tapered portal frame
// 		f.GenTaper()
// 	}
// 	return
// }

// //GenCp generates f.Mod.Cp and f.Mod.Dims, f.Mod.Styps - do this in AddMem
// func (f *Portal) GenCp(){
// 	switch f.Mtyp{
// 		case 2:
// 		//given - styp, sdx in styps, sections
// 		//read sec dims ReadSecDims(sectyp, sdx, ax int), store in mod.dims
// 		//calc sec with GetStlCp(frmtyp, sectyp, sdx, ax int), store in mod.Cp
// 	}
// 	for i, s := range f.Sections{}
// }

func (p *Portal) PurlinInit(){
	switch p.Prlntyp{
		case 1:
		//rcc
		case 2:
		//steel
		case 3:
		//timber
		case 4:
		//cfs
	}
}

func PurlinInit(dl, ll, slope, spacing, span float64, psfs []float64, mono bool) (fdl, fll, fulg, rafterlen, rise, purlinspc float64, idx int){
	var nprz float64
	//var l float64
	if mono{
		rise = slope * span
		rafterlen = math.Sqrt(math.Pow(rise, 2) + math.Pow(span, 2))
		nprz = math.Ceil(rafterlen/1.4)
		purlinspc = math.Round(1000.0*rafterlen/nprz)/1000.0
	} else {
		rise = slope * span/2.0
		rafterlen = math.Sqrt(math.Pow(rise, 2) + math.Pow(span/2.0, 2))
		log.Println("rafterlen->",rafterlen)
		nprz = math.Ceil(rafterlen/1.4)
		purlinspc = math.Round(1000.0*rafterlen/nprz)/1000.0
	}
	prdl := purlinspc * dl; prll := purlinspc * ll
	var pul, dlsf float64
	if len(psfs) == 0{pul = 1.5 * prdl + 1.5 * prll; dlsf = 1.5} else {pul = psfs[0] * prdl + psfs[1] * prll; dlsf = psfs[0]}
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	sheet := filepath.Join(basepath,"../data/steel/","cfszsec.csv")
	csvfile, err := os.Open(sheet)
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(csvfile)
	//log.Println("purlin spacing-",purlinspc, "prdl kn/m", prdl, "prll kn/m",prll)
	for i, val := range []float64{1.83 , 2.44 , 3.05 , 3.66 , 4.27 , 4.88 , 5.49 , 6.1} {
		if val > spacing/2.0{
			idx = i
		}
	}
	var mul, w, wprz float64
	for i := 0; i < df.Nrow() - 1; i++{ //
		//log.Println("section->",df.Elem(i,0))
		//log.Println("mass->",df.Elem(i,16))
		//log.Println("pul->",pul)
		wprz = df.Elem(i,16).Float()*9.81/1e3
		w = wprz*dlsf + pul
		//log.Println("w ->", w)
		mul = w * math.Pow(spacing/2.0,2)/10.0
		//log.Println("checking moment for braced length->",df.Elem(i,23+idx),"vs",mul, "OK?",df.Elem(i,23+idx).Float() > mul)
		if df.Elem(i,23+idx).Float() > mul{idx = i; break}
	}
	//fmt.Println("section->",df.Elem(idx,0))
	//fmt.Println("mass->",df.Elem(idx,16))
	//get final load in kn/m of portal frame
	//gable end load is that/2
	var llsf float64
	fdl = wprz * nprz * spacing/rafterlen + spacing * dl
	if len(psfs) == 0{llsf = 1.5} else {llsf = psfs[1]}
	fll = ll * spacing
	fulg = fdl*dlsf + fll* llsf + dlsf * 0.1 * spacing
	log.Println(fdl, fll, fulg)
	return
}

func PortalPrelim(f *Portal, fulg float64) (em, cp [][]float64, bmwt, colwt float64){
	//use pinned base formula from sci p 399
	//TODO check for alt formulae
	var l float64
	if f.Mono{l = f.Span} else {l = f.Span/2.0}
	theta := f.Rise/f.Height; m := 1.0 + theta
	k := f.Height/l/1.5
	b := 2.0 * (1.0 + k) + m; c := 1.0 + 2.0 * m
	n := b + m * c
	//moment at eaves
	me := fulg * math.Pow(l,2) * (3.0 + 5.0 * m)/16.0/n
	//moment at apex
	ma := fulg * math.Pow(l,2)/8.0 + m * me
	log.Println("moment at apex->",ma,"at eaves->",me, "kn-m")
	pbc := 165.0
	sxxr := ma*1e3/pbc
	log.Println("sec modulus req->",sxxr,"mm4")
	_, based, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(based)
	sheet := filepath.Join(basepath,"../data/steel/isteel","ISMZ.csv")
	csvfile, err := os.Open(sheet)
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(csvfile)
	var bmsecs, colsecs []int
	var mbm, mcol int
	bmwt = 0.; colwt = 0.
	for i := 0; i < df.Nrow() -1; i++{
		if len(bmsecs) > 10 {break}
		if df.Elem(i,11).Float() >= sxxr{
			bmsecs = append(bmsecs, i)
			//log.Println("section->", df.Elem(i,0), "weight kg/m->",df.Elem(i,1))
			if bmwt == 0{
				mbm = i
				bmwt = df.Elem(i,1).Float()
			} else if bmwt > df.Elem(i,1).Float(){
				mbm = i
				bmwt = df.Elem(i,1).Float()
			}
		}
	}
	log.Println("min bm->",bmwt, df.Elem(mbm,0), df.Elem(mbm,8))
	for i:=0; i < df.Nrow()-1;i++{
		if len(colsecs) > 10{break}
		if df.Elem(i,8).Float() >= df.Elem(mbm,8).Float()*1.5{
			colsecs = append(colsecs, i)
			if colwt == 0 {
				mcol = i
				colwt = df.Elem(i, 1).Float()
			} else if colwt > df.Elem(i, 1).Float(){
				mcol = i
				colwt = df.Elem(i,1).Float()
			}
		}
	}
	log.Println("min col->",colwt, df.Elem(mcol,0), df.Elem(mcol,8), df.Elem(mcol,8).Float()/df.Elem(mbm,8).Float())
	cp = make([][]float64, 2)
	for i:=0; i < 2; i++{
		cp[i] = make([]float64,2)
	}
	cp[0][0] = df.Elem(mcol,2).Float()*1e-4; cp[0][1] = df.Elem(mcol,7).Float()*1e-8
	cp[1][0] = df.Elem(mbm,2).Float()*1e-4; cp[1][1] = df.Elem(mbm,7).Float()*1e-8
	em = [][]float64{{2.1e8}}
	return
}

func (f *Portal) InitOld()(err error){
	//col sec 1, bmsec 2
	var x, y, fdl, fll, fulg, bmwt, colwt, lspan float64
	f.Nframes = int(f.Ly/f.Spacing)
	f.Spacing = f.Ly/float64(f.Nframes)
	f.Nframes += 1
	if f.Mono{lspan = f.Span} else {lspan = f.Span/2.0}

	if f.Verbose{
		log.Println("starting portal frame analysis of->",f.Title)
		log.Println("number of frames->",f.Nframes)
	}
	if f.Prlnwt == 0.0{
		fdl, fll, fulg, f.LR, f.Rise, f.Prlnspc, f.Prlnsec = PurlinInit(f.DL, f.LL, f.Slope, f.Spacing, f.Span, f.PSFs, f.Mono)
	}
	var coords, em, cp [][]float64
	var mprp, msup [][]int
	if f.W == 0{
		log.Println("w frame kn/m->",fulg)
		em, cp, bmwt, colwt = PortalPrelim(f, fulg)
	}
	log.Println("em, cp",em, cp)
	log.Println("lspan",lspan)
	bmwt = bmwt * 9.81/1e3; colwt = colwt * 9.81/1e3
	log.Println("bm wt, col wt", bmwt, colwt, "fdl,fll->", fdl, fll)
	var dlsf, llsf float64
	if len(f.PSFs) == 0{dlsf, llsf = 1.5, 1.5} else {dlsf, llsf = f.PSFs[0], f.PSFs[1]}
	fulg -= 0.1 * dlsf * f.Spacing
	fdl += bmwt
	angle := math.Atan(f.Slope)
	log.Println("roof angle in radians->",angle, "degrees->",angle * 180.0/math.Pi)
	pd, cpos, cneg, _ := wltable6(f.Vz, f.Height, f.Span, f.Slope, f.Cpi)
	//em := [][]float64{{210000}}
	coords = append(coords, []float64{x,y})
	var bmvec, colvec []int

	log.Println("bmvec->",bmvec)
	var msloads, jsloads [][]float64
	var dirfac float64
	for _, bm := range bmvec{
		dirfac = 1.0
		//if mprp[bm-1][1] - mprp[bm-1][0] == 1{dirfac = -1.0} else {dirfac = 1.0}
		msloads = append(msloads, []float64{float64(bm), 3, dirfac * dlsf * fdl*math.Cos(angle), 0, 0, 0, 1})
		msloads = append(msloads, []float64{float64(bm), 3, dirfac * llsf * fll*math.Cos(angle), 0, 0, 0, 2})
		//msloads = append(msloads, []float64{float64(bm), 6, dlsf * fdl*math.Sin(angle), 0, 0, 0, 1})
		//msloads = append(msloads, []float64{float64(bm), 6, llsf * fll*math.Sin(angle), 0, 0, 0, 2})
	}
	//init mem sizes - column = 1.0 to 1.5 times beam stiffness/inertia Ixx
	for _, col := range colvec{
		msloads = append(msloads, []float64{float64(col), 6, dlsf * colwt, 0, 0, 0, 1})
	}
	pd = pd/1e3
	log.Println("summary of loads on frame->")
	//log.Println("purlins->",dlsf * wprz * nprz/rafterlen, "kn-m2")
	log.Println("fdl->",fdl, "kn-m")
	log.Println("fll->", fll, "kn-m")
	log.Println("fulg->", fulg, "kn-m")
	log.Println("design pressure pd-",pd,"kn/m2")
	log.Println("design wind load cases +ve cpi ", f.Cpi)
	log.Println(cpos)
	log.Println("design wind load cases -ve cpi ", -f.Cpi)
	log.Println(cneg)
	log.Println(ColorRed,msloads,ColorReset)
	//term := "qt"
	//first do service loads, then combine?
	mod := &Model{
		Ncjt:3,
		Cmdz:[]string{"2df","mks","1"},
		Coords:coords,
		Mprp:mprp,
		Supports:msup,
		Msloads:msloads,
		Jloads:jsloads,
		Em:em,
		Cp:cp,

	}

	_, err = CalcFrm2d(mod,3)
	if err != nil{
		log.Println("ERRORE,errore->",err)
		return
	}
	//report, _ := frmrez[6].(string)
	//js, _ := frmrez[0].(map[int]*kass.Node)
	//ms, _ := frmrez[1].(map[int]*Mem)
	//log.Println(report)
	pltchn := make(chan string,1)
	mod.Frcscale = 2.0
	PlotFrm2d(mod, f.Term, pltchn)
	txtplot := <- pltchn
	log.Println(txtplot)
	return
}


	mod := &kass.Model{
		Ncjt:3,
		Cmdz:[]string{"2df","mks","1"},
		Coords:coords,
		Mprp:mprp,
		Supports:msup,
		Msloads:msloads,
		Jloads:jsloads,
		Em:em,
		Cp:cp,
	}
	frmrez, err := kass.CalcFrm2d(mod,3)
	if err != nil{
		log.Println("ERRORE,errore->",err)
		return
	}
	report, _ := frmrez[6].(string)
	js, _ := frmrez[0].(map[int]*kass.Node)
	ms, _ := frmrez[1].(map[int]*kass.Mem)
	log.Println(report)


	txtplot := kass.DrawMod2d(mod, ms, term)
	log.Println(txtplot)
	for _, node := range js{
		if node.React[1] != 0{
			log.Println(node.React)
			log.Println("footing design")
			colx := 0.45; coly := 0.45; df := 0.0; eo := 0.25
			fck := 25.0; fy := 500.0
			sbc := 100.0; pgck := 24.0; pgsoil := 15.0; nomcvr := 0.06; dmin := 0.25
			pus := []float64{node.React[1]}
			mxs := []float64{0}
			mys := []float64{0}
			psfs := []float64{1.0,1.0}
			shape := "square"
			sloped := true
			dlfac := false
			mosh.FtngDzRojas(colx, coly, fck, fy, df, dmin, eo, sbc, pgck, pgsoil, nomcvr, pus, mxs, mys, psfs, shape, sloped, dlfac, "dumb")
		}
	}

func HassanEx5(){
	mod := &kass.Model{
		Ncjt:3,
		Cmdz:[]string{"2df","mks","1"},
		Coords:[][]float64{
			{0,0},
			{40,0},
			{0,5},
			{40,5},
			{20,10.36},
		},
		Supports:[][]int{{1,-1,-1,0},{2,-1,-1,0}},
		Mprp:[][]int{{1,3,1,1,0},{2,4,1,1,0},{3,5,1,2,0},{4,5,1,2,0}},
		Jloads:[][]float64{},
		Msloads:[][]float64{{3,3,10.32,0,0,0},{3,6,2.77,0,0,0},{4,3,10.32,0,0,0},{4,6,2.77,0,0,0},},
		Em:[][]float64{{2.1e8}},
		Cp:[][]float64{{105e-4,47540e-8},{129e-4,61520e-8}},
	}
	frmrez, err := kass.CalcFrm2d(mod, 3)
	if err != nil{
		log.Println(err)
		return
	}
	report, _ := frmrez[6].(string)
	log.Println(report)
}

	//log.Println(coords)
	//log.Println(mprp)
	//log.Println(bmvec)
	//log.Println(msup)
	//log.Println("bmvec-",bmvec)
	//log.Println("colvec-",colvec)
	//pltchn := make(chan string, 1)
	//go kass.PlotGenTrs(coords, mprp, pltchn)
	//pltstr := <-pltchn
	//log.Println(pltstr)


	lspan := f.Span/2.0
	ly := lspan
	ty := 1.0
	lbr := 200.0
	tbr := 20.0
	nsecs := 1
	grd := 43
	sectyp := 0
	brchck := false
	yeolde := false
	ldcases := [][]float64{{1.0,3.0,wdl,0,0,0,1},{1.0,3.0,wll,0,0,0,1}}
	rez := StlBmDBs(lspan, ly, ty, lbr, tbr, ldcases, sectyp, grd, nsecs, brchck, yeolde)
	bdx := rez[0]
	df := StlSecBs(sectyp)
	wb, arb, ixb := df.Elem(bdx,2).Float(), df.Elem(bdx,23).Float(), df.Elem(bdx,11).Float()
	fil := df.Filter(
		dataframe.F{Colname:"ix", Comparator:series.Greater, Comparando:1.5*ixb},
	)
	//log.Println(fil.Nrow())
	//log.Println(fil.Subset([]int{fil.Nrow()-1}))
	cdx := fil.Nrow()-1
	wc, arc, ixc := df.Elem(cdx,2).Float(), df.Elem(cdx,23).Float(), df.Elem(cdx,11).Float()
	//log.Println(wc, arc, ixc, wb, arb, ixb)
	//START WITH SAME SECTION
	cp = [][]float64{{arc*1e-4, ixc*1e-8},{arb*1e-4, ixb*1e-8}}
	em = [][]float64{{2.1*1e8}}
	//log.Println("section->",df.Elem(i,1))
	//log.Println("depth, web thickness->",df.Elem(i,3), df.Elem(i,6))
	//log.Println("area, zx, zy->",df.Elem(i,23),df.Elem(i,15), df.Elem(i,16))
	//log.Println("rx, ry->",df.Elem(i,13), df.Elem(i,14))
*/
