package barf

import (
	"fmt"
	"math"
	"time"
	"strings"
	"strconv"
	"github.com/olekukonko/tablewriter"
	flay"barf/flay"
	kass"barf/kass"
	mosh"barf/mosh"
	draw"barf/draw"
)

type StlFlr struct {
	Title     string
	Idx       int
	Xs        []float64
	Ys        []float64
	Zs        []float64  
	Vz        float64             //design wind speed vz
	Cpi       float64
	Ht        float64
	Hcol      float64             //total column height
	DL        float64
	LL        float64
	TDL       float64
	TLL       float64
	Df        float64             //depth of footing
	Wt        float64
	Bua       float64             //hi BUA
	Sbc       float64             //soil bearing capacity
	Fck       float64             //grade of pedestal/footing concrete
	Lablt     float64             //length of anchor bolt
	Sname     string
	Ftyp      string
	Term      string
	Folder    string
	Report    string
	Breport   string               //base report
	Txtplot   string
	Brc       []int                //xdx, ydx
	Coords    [][]float64
	Cols      []int
	Mat       [][]int
	Colon     map[string]bool
	Gon       map[string]bool
	Jon       map[string]bool
	Nmap      map[kass.Pt]int   //maps point to column
	Bms       map[string][]string //beams map to joists
	Idmap     map[string][]string //(ordered) map of members by category
	Jsts      [][]string
	Ojs       []string            //ohng joists
	Cgs       []string            //column clvr girders
	Jspc      float64
	Ohng      float64
	Jspcs     [][]float64
	Cmap      map[int][][]string //maps columns to girders/joists (x, y)
	Bmz       map[string]Bm
	Colz      map[string]Col
	Smap      map[string][]float64 //maps sections to total length
	Psmap     map[string][]float64 //maps plates to weight
	Qmap      map[string][]float64 
	Pdz       map[string]mosh.RccCol //pedestals
	Ftngz     map[string]mosh.RccFtng //ftngs.
	Cbz       map[string]Cbs          //map of column bases
	Jntz      map[string]kass.Blt     //map of joints
	Jnts      map[string][]string     //joint index main mem-X-supporting mem
	Jntord    []string                //for ordered joint table
	Fquant    []float64               //vrcc, wstl, afw
	Onism     bool
	Matread   bool
	Mapread   bool
	Web       bool
	Sub       bool
	Lsb       bool
	Cdz       bool
	Dzcbs     bool  //design column base plate
	Dzbs      bool  //design pedestal + footing
	Verbose   bool
	Yeolde    bool  //use (trusty n reliable) bs449 for column design
	Nflr      int

}

//FlrRez stores floor rez vals
type FlrRez struct{
	Cols map[int]Col
	Bms  map[int]Bm
}

//getbmdx returns a beam index string
func getbmdx(i, j int)(string){
	if i < j {
		return fmt.Sprintf("%d-%d",i,j)
	}
	return fmt.Sprintf("%d-%d",j,i)
}

//getybmdx returns a joist index string
func getybmdx(i,j,a,b int)(string){
	g1 := getbmdx(i,j)
	g2 := getbmdx(a,b)
	m1 := i; if j < i{
		m1 = j
	}
	m2 := a; if b < a{
		m2 = b
	}
	if m1 < m2{
		return fmt.Sprintf("%s-%s",g1,g2)
	} else {
		return fmt.Sprintf("%s-%s",g2,g1)
	}
}

//Init inits a StlFlr with defaults
func (f *StlFlr) Init()(err error){
	if len(f.Xs) < 2 || len(f.Ys) < 2 || len(f.Zs) < 2{
		err = fmt.Errorf("floor geom not specified xs %f ys %f zs %f",f.Xs, f.Ys, f.Zs)
	}
	if f.DL == 0.0{f.DL = 2.0}
	if f.LL == 0.0{f.LL = 3.0}
	if f.Jspc == 0.0{f.Jspc = 1800}
	if f.Df == 0.0{f.Df = 2000.0}
	if f.Fck == 0.0{f.Fck = 25.0}
	if f.Sbc == 0.0{f.Sbc = 150.0}
	if f.Matread{
		if len(f.Xs)-1 != len(f.Mat[0]){
			err = fmt.Errorf("error in x mat %v %v",f.Xs,f.Mat[0])
			return
		}
		
		if len(f.Ys)-1 != len(f.Mat){
			err = fmt.Errorf("error in y mat %v %v",f.Ys,f.Mat)
			return
		}
	}
	switch f.Sname{
		case "":
		//default ismb
		f.Sname = "i"
		f.Onism = true
		case "ub":
		f.Onism = false
	}
	if f.Vz == 0.0{f.Vz = 50.0}
	if f.Cpi == 0.0{f.Cpi = 0.5}
	return
}


//Calc calcs a StlFlr
func (f *StlFlr) Calc()(err error){
	//one way joist-girder system
	//is the only system it knows
	err = f.Init()
	if err != nil{
		return
	}
	//gen coords
	var errstr string
	var ecount int
	//make idmap for (hopefully) ordered tables
	f.Idmap = make(map[string][]string)
	for _, mtyp := range []string{"col","girder","j.girder","joist","o.joist","c.girder"}{
		f.Idmap[mtyp] = []string{}
	}
	f.Cmap = make(map[int][][]string)
	f.Nmap  = make(map[kass.Pt]int)
	f.Colon = make(map[string]bool)
	f.Gon = make(map[string]bool)
	f.Jon = make(map[string]bool)
	f.Bua = 0.0
	for _, y := range f.Ys{
		for _, x := range f.Xs{
			f.Coords = append(f.Coords,[]float64{x,y})
			f.Cols = append(f.Cols, len(f.Coords))
			f.Cmap[len(f.Coords)] = make([][]string,2)
			cdx := fmt.Sprintf("%v",len(f.Coords))
			f.Idmap["col"] = append(f.Idmap["col"],cdx)
			f.Nmap[kass.Pt{X:x,Y:y}] = len(f.Coords)
			f.Colon[cdx] = true
		}
	}

	//get joist spacing per bay
	if f.Jspc == 0.0{f.Jspc = 1800}
	f.Jspcs = make([][]float64,len(f.Ys)-1)
	for i := range f.Ys{
		if i == 0{
			continue
		}
		dy := f.Ys[i] - f.Ys[i-1]
		ndiv := math.Ceil(dy/f.Jspc)
		for j := 1; j < int(ndiv); j++{
			f.Jspcs[i-1] = append(f.Jspcs[i-1],float64(j)*dy/ndiv)
		}
	}
	f.Bms = make(map[string][]string)
	f.Bmz = make(map[string]Bm)
	f.Colz = make(map[string]Col)
	f.Cbz = make(map[string]Cbs)
	f.Psmap = make(map[string][]float64)
	//gen girders 
	for _, idx := range f.Cols{
		if idx + len(f.Xs) <= len(f.Cols){
			gdx := getbmdx(idx, idx+len(f.Xs))
			f.Bms[gdx] = []string{}
			f.Idmap["girder"] = append(f.Idmap["girder"],gdx)
		}
	}
	//gen joists - go by panels (ci, ci+1, ci+nx, ci+nx+1)
	f.Jsts = make([][]string,(len(f.Ys)-1))
	for i := range f.Jsts{
		f.Jsts[i] = make([]string,len(f.Xs)-1)
	}
	//dz joists
	for _, idx := range f.Cols[:len(f.Cols)-len(f.Xs)]{
		xdx := idx % len(f.Xs) //joist span index
		if xdx == 0{
			continue
		}
		ydx := idx/len(f.Xs) //joist bay index
		
		g1 := getbmdx(idx, idx+len(f.Xs))
		g2 := getbmdx(idx+1,idx+len(f.Xs)+1)
		jdx := getybmdx(idx, idx + len(f.Xs),idx+1,idx+len(f.Xs)+1)
		f.Jsts[ydx][xdx-1] = jdx
		jspc := f.Jspcs[ydx][0]

		j1 := getbmdx(idx, idx+1)
		j2 := getbmdx(idx+len(f.Xs),idx+len(f.Xs)+1)
		c1 := fmt.Sprintf("%v",idx)
		c2 := fmt.Sprintf("%v",idx+1)
		c3 := fmt.Sprintf("%v",idx+len(f.Xs))
		c4 := fmt.Sprintf("%v",idx+len(f.Xs)+1)
		
		if f.Matread{
			// fmt.Println("at cell idx, ydx, xdx-1, val",idx, ydx, xdx-1,f.Mat[ydx][xdx-1])
			// if f.Mat[ydx][xdx-1] == 0{
			// 	fmt.Println(ColorRed,"zero cell",c1,c2,c3,c4,ColorReset)
			// }
			// fmt.Println("girders g1 g2")
			// fmt.Println(g1)
 			// fmt.Println(g2)
			// fmt.Println("joists j1 j2")
			// fmt.Println(j1)
			// fmt.Println(j2)
			// fmt.Println("cols",c1,c2,c3,c4)
		 	//f.PanMar(idx, ydx, xdx-1)
			if _, ok := f.Gon[g1]; !ok{
				f.Gon[g1] = true
			}
			if _, ok := f.Gon[g2]; !ok{
				f.Gon[g2] = true
			}
			if _, ok := f.Gon[j1]; !ok{
				f.Gon[j1] = true
			}
			if _, ok := f.Gon[j2]; !ok{
				f.Gon[j2] = true
			}
			for _, jval := range []string{"bot","mid","top"}{
				f.Jon[jdx+"-"+jval] = true
			}
			f.Jon[jdx] = true
			//TODO
			//if cell val == -1 ?
			//if cell val == 0
			if f.Mat[ydx][xdx-1] == 0{
				f.Jon[jdx+"-"+"mid"] = false
				switch{
					case ydx == 0:
					//bottom row panel
					switch xdx{
						case 1:
						//bottom left
						//g1, j1 off
						f.Gon[g1] = false
						f.Gon[j1] = false
						f.Colon[c1] = false
						if f.Mat[ydx][xdx] == 0{
							f.Gon[g2] = false
						}
						if f.Mat[ydx+1][xdx-1] == 0{
							f.Gon[j2] = false
						}
						case len(f.Xs)-1:
						//bottom right
						f.Gon[g2] = false
						f.Gon[j1] = false
						f.Colon[c2] = false
						if f.Mat[ydx+1][xdx-1] == 0{
							f.Gon[j2] = false
						}
						default:
						//bottom
						f.Gon[j1] = false
						//right off
						if f.Mat[ydx][xdx] == 0{
							f.Gon[g2] = false
						}
						if f.Mat[ydx+1][xdx-1] == 0{
							f.Gon[j2] = false
						}
						if f.Mat[ydx][xdx-2] == 0{
							f.Gon[g1] = false
						}
						
						
					}
					case ydx == len(f.Ys)-2:
					
					//top row panel
					switch xdx{
						case 1:
						//top left
						f.Gon[g1] = false
						f.Gon[j2] = false
						f.Colon[c3] = false
						if f.Mat[ydx][xdx] == 0{
							f.Gon[g2] = false
						}
						if f.Mat[ydx-1][xdx-1] == 0{
							f.Gon[j1] = false
						}
						case len(f.Xs)-1:
						//top right
						f.Gon[g2] = false
						f.Gon[j2] = false
						f.Colon[c4] = false
						if f.Mat[ydx][xdx-2] == 0{
							f.Gon[g1] = false
						}
						if f.Mat[ydx-1][xdx-1] == 0{
							f.Gon[j1] = false
						}
						default:
						f.Gon[j2] = false
						if f.Mat[ydx][xdx-2] == 0{
							f.Gon[g1] = false
						}
					}
					case xdx == 1:
					//left column panel
					//left girder off
					f.Gon[g1] = false
					if f.Mat[ydx-1][xdx-1] == 0{
						f.Gon[j1] = false
					}
					case xdx == len(f.Xs)-1:
					//right column panel
					//right girder off
					f.Gon[g2] = false
					if f.Mat[ydx-1][xdx-1] == 0{
						f.Gon[j1] = false
					}
					default:
					//interior panel - mid interior joists off
					if f.Mat[ydx][xdx-2] == 0{
						f.Gon[g1] = false
					}
					if f.Mat[ydx][xdx] == 0{
						f.Gon[g2] = false
					}
					if f.Mat[ydx-1][xdx-1] == 0{
						f.Gon[j1] = false
					}
					if f.Mat[ydx+1][xdx-1] == 0{
						f.Gon[j2] = false
					}
					
				}
			}
			
		}
		//calc start and end point here
		v1 := f.Coords[idx-1]; v2 := f.Coords[idx+len(f.Xs)-1]
		v3 := f.Coords[idx]; v4 := f.Coords[idx+len(f.Xs)]
		d1 := kass.Dist2d(v1,v3)
		d2 := kass.Dist2d(v3, v4)
		if f.Matread && f.Mat[ydx][xdx-1] != 0{
			f.Bua += d1 * d2 * 10.76/1e6
		}
		strt := []float64{}
		end := []float64{}
		var lsup, rsup string
		//get tributary area
		for i := 0; i < 3; i++{
			jstr := jdx
			jtrib := jspc/2.0
			switch i{
				case 0:
				//bottom
				jstr += "-bot"
				lsup = c1
				rsup = c2
				switch ydx{
					case 0:
					//bottom trib is ohng/2
					jtrib += f.Ohng/2.0
					default:
					//add trib below (ydx-1) trib/2 
					jtrib += f.Jspcs[ydx-1][0]/2.0
					
				}
				strt = v1; end = v3
				case 1:
				//middle
				jstr += "-mid"
				lsup = g1
				rsup = g2
				jtrib += jspc/2.0
				f.Bms[g1] = append(f.Bms[g1],jstr)
				f.Bms[g2] = append(f.Bms[g2],jstr)
				//lerp by jspc at v1 and v3
				//strt = kass.Lerpvec(jspc/d1, v1, v2)
				//end = kass.Lerpvec(jspc/d2, v3, v4)
				strt = v1
				end = v3
				case 2:
				//top
				jstr += "-top"
				lsup = c3
				rsup = c4
				switch ydx{
					case len(f.Jspcs)-1:
					//upper trib is ohng/2
					jtrib += f.Ohng/2.0
					default:
					//add upper trib 
					jtrib += f.Jspcs[ydx+1][0]/2.0
					
				}
				strt = v2; end = v4
				
			}
			jtrib = jtrib/1000.0
			dl := f.DL*jtrib
			ll := f.LL*jtrib
			bm := Bm{
				Title:jstr,
				Sname:f.Sname,
				Lsb:true,
				Lspan:f.Xs[xdx]-f.Xs[xdx-1],
				DL:dl,
				LL:ll,
				Dtyp:1,
				Store:true,
				Dsgn:true,
				Onism:f.Onism,
				Selfwt:true,
				Term:"dumb",
				Start:strt,
				End:end,
				Name:"joist",
				Lsup:lsup,
				Rsup:rsup,
			}
			
			if f.Matread{	
				if jston, ok := f.Jon[jstr]; ok{
					if !jston{bm.Ignore = true}
				}
			}
			e := BmDz(&bm)
			f.Bmz[jstr] = bm
			if e != nil{
				
				errstr += fmt.Sprintf("joist %s error %s\n",jstr, fmt.Sprint(e))
				ecount++
			} // else {
			// 	f.Jerc[jstr] = bm.Ssecs[0].Vu
			// }
			//dz connections here?
			f.Idmap["joist"] = append(f.Idmap["joist"],jstr)
		}
	}
	//now do girders(bms)
	for bdx, jvec := range f.Bms{
		cvec := strings.Split(bdx,"-")
		c1, _ := strconv.Atoi(cvec[0])
		c2, _ := strconv.Atoi(cvec[1])
		lspan := kass.Dist2d(f.Coords[c1-1],f.Coords[c2-1])
		xdx := c1 % len(f.Xs) //span index
		ydx := c1/len(f.Xs) //bay index
		//when c1 is the last column of the penultimate bay, ydx changes
		if ydx == len(f.Jspcs) && xdx == 0{
			ydx = len(f.Jspcs)-1
		}
		f.Cmap[c1][1] = append(f.Cmap[c1][1],bdx)
		f.Cmap[c2][1] = append(f.Cmap[c2][1],bdx)
		//lspcs := Jspcs[ydx]
		var wu float64
		for _, jst := range jvec{
			//wu += math.Round(f.Jerc[jst])
			switch f.Matread{
				case true:	
				if f.Jon[jst]{
					wu += math.Round(f.Bmz[jst].Ssecs[0].Vu)
				}
				case false:
				wu += math.Round(f.Bmz[jst].Ssecs[0].Vu)
			}
		}
		ldcases := [][]float64{}
		for _, spc := range f.Jspcs[ydx]{
			ldcases = append(ldcases, []float64{1,1,wu,0,math.Round(spc),0,1})
		}
		bm := Bm{
			Title:bdx,
			Sname:f.Sname,
			Lsb:true,
			Lspan:lspan,
			DL:0.0,
			LL:0.0,
			Dtyp:1,
			Store:true,
			Dsgn:true,
			Onism:f.Onism,
			Selfwt:true,
			Ldcases:ldcases,
			Term:"dumb",
			Start:f.Coords[c1-1],
			End:f.Coords[c2-1],
			Name:"girder",
			Lsup:fmt.Sprintf("%v",c1),
			Rsup:fmt.Sprintf("%v",c2),
		}
		var dl, ll float64
		switch xdx{
			case 0,1:
			//add overhang/2 DL and LL
			dl = f.DL * f.Ohng/2000.0
			ll = f.LL * f.Ohng/2000.0
		}
		bm.DL = dl; bm.LL = ll

		if f.Matread && !f.Gon[bdx]{bm.Ignore = true}
		e := BmDz(&bm)
		if e != nil{
			errstr += fmt.Sprintf("girder %s error %s\n",bdx, fmt.Sprint(e))
			ecount++
			f.Bmz[bdx] = bm
			continue
		}
		f.Bmz[bdx] = bm
		//dz connections
	}
	if f.Ohng > 0{
		
		//overhang joists in x
		for i := range f.Xs{
			if i == 0{
				continue
			}
			dx := f.Xs[i] - f.Xs[i-1]
			jdx := fmt.Sprintf("%d-%d-x-o",i,i+1)
			
			f.Idmap["o.joist"] = append(f.Idmap["o.joist"],jdx)
			dl := f.DL * f.Ohng/2000.0
			ll := f.LL * f.Ohng/2000.0
			bm := Bm{
				Title:jdx,
				Sname:f.Sname,
				Lsb:true,
				Lspan:dx,
				DL:math.Round(dl*100.)/100.,
				LL:math.Round(ll*100.)/100.,
				Dtyp:1,
				Store:true,
				Dsgn:true,
				Onism:f.Onism,
				Selfwt:true,
				Term:"dumb",
				Start:[]float64{f.Coords[i-1][0],f.Coords[i-1][1]-f.Ohng},
				End:[]float64{f.Coords[i][0],f.Coords[i][1]-f.Ohng},
				Name:"o.joist",
			}
			
			e := BmDz(&bm)
			if e != nil{
				errstr += fmt.Sprintf("o.joist %s error %s\n",jdx, fmt.Sprint(e))
				ecount++
			}
			f.Bmz[jdx] = bm
			f.Ojs = append(f.Ojs,jdx)
		}
		//overhang joists in y
		for i := range f.Ys{
			
			if i == 0{
				continue
			}
			sc := i*len(f.Xs) + 1
			ec := (i-1)*len(f.Xs) + 1
			// fmt.Println("i, sc, ec",i, sc, ec)
			// fmt.Println(f.Coords[sc-1])
			// fmt.Println(f.Coords[ec-1])
			strt := []float64{f.Coords[sc-1][0]-f.Ohng,f.Coords[sc-1][1]}
			end := []float64{f.Coords[ec-1][0]-f.Ohng,f.Coords[ec-1][1]}
			dy := f.Ys[i] - f.Ys[i-1]
			//dy := kass.Dist2d(strt, end)
			jdx := fmt.Sprintf("%d-%d-y-o",i,i+1)
			f.Idmap["o.joist"] = append(f.Idmap["o.joist"],jdx)
			dl := f.DL * f.Ohng/2000.0
			ll := f.LL * f.Ohng/2000.0
			bm := Bm{
				Title:jdx,
				Sname:f.Sname,
				Lsb:true,
				Lspan:dy,
				DL:math.Round(dl*100.)/100.,
				LL:math.Round(ll*100.)/100.,
				Dtyp:1,
				Store:true,
				Dsgn:true,
				Onism:f.Onism,
				Selfwt:true,
				Term:"dumb",
				Start:strt,
				End:end,
				Name:"o.joist",
			}
			e := BmDz(&bm)
			if e != nil{
				errstr += fmt.Sprintf("o.joist %s error %s\n",jdx, fmt.Sprint(e))
				ecount++
			}
			f.Bmz[jdx] = bm
			f.Ojs = append(f.Ojs,jdx)
		}
	}
	
	//design column clvr girders (and add clvr loads to column)
	for _, idx := range f.Cols{
		xdx := idx%len(f.Xs)
		ydx := idx/len(f.Xs) + 1
		if xdx == 0{
			xdx = len(f.Xs)
			ydx--
		}
		if f.Matread && !f.Colon[fmt.Sprintf("%v",idx)]{
			continue
		}
		if f.Ohng > 0.0{
			strt := f.Coords[idx-1]
			end := []float64{f.Coords[idx-1][0]-f.Ohng,f.Coords[idx-1][1]}
			//only corner/edge columns have clvrs
			if xdx == 1 || xdx == len(f.Xs){
				//dsgn clvr in x
				//loaded by y ohng joist (xdx-)
				gdx := fmt.Sprintf("%d-x",idx)
				wu := 0.0
				switch{
					case ydx == 1:
					//load with ydx, ydx + 1 joist
					jdx := fmt.Sprintf("%d-%d-y-o",ydx, ydx+1)
					wu += f.Bmz[jdx].Ssecs[0].Vu
					case ydx == len(f.Ys):
					//load with ydx -1, ydx joist
					jdx := fmt.Sprintf("%d-%d-y-o",ydx-1, ydx)
					wu += f.Bmz[jdx].Ssecs[0].Vu
					default:
					//load with ydx-1, ydx and ydx, ydx+1 joists
					jdx := fmt.Sprintf("%d-%d-y-o",ydx, ydx+1)
					wu += f.Bmz[jdx].Ssecs[0].Vu
					jdx = fmt.Sprintf("%d-%d-y-o",ydx-1, ydx)
					wu += f.Bmz[jdx].Ssecs[0].Vu
				}
				if xdx == len(f.Xs){
					end = []float64{f.Coords[idx-1][0]+f.Ohng,f.Coords[idx-1][1]}
				}
				bm := Bm{
					Title:gdx,
					Sname:f.Sname,
					Lsb:true,
					Lspan:f.Ohng,
					Dtyp:-1,
					Store:true,
					Dsgn:true,
					Onism:f.Onism,
					Selfwt:true,
					Term:"dumb",
					Ldcases:[][]float64{{1,1,wu,0,f.Ohng-50.0,0,1}},
					Start:strt,
					End:end,
					Name:"c.girder",
				}
				e := BmDz(&bm)
				if e != nil{
					errstr += fmt.Sprintf("c.girder %s error %s\n",gdx, fmt.Sprint(e))
					ecount++
				} 
				f.Bmz[gdx] = bm
				f.Cgs = append(f.Cgs, gdx)
				f.Cmap[idx][0] = append(f.Cmap[idx][0],gdx)
				f.Idmap["c.girder"] = append(f.Idmap["c.girder"],gdx)
			}
			if ydx == 1 || ydx == len(f.Ys){
				
				if f.Ohng == 0.0{continue}
				//dsgn clvr in y
				//loaded by x ohng joist (xdx-)
				gdx := fmt.Sprintf("%d-y",idx)
				wu := 0.0
				end := []float64{f.Coords[idx-1][0],f.Coords[idx-1][1]-f.Ohng}
				strt := f.Coords[idx-1]
				if ydx == len(f.Ys){
					end = []float64{f.Coords[idx-1][0],f.Coords[idx-1][1]+f.Ohng}
				}
				switch{
					case xdx == 1:
					//load with xdx, xdx + 1 joist
					jdx := fmt.Sprintf("%d-%d-x-o",xdx, xdx+1)
					wu += f.Bmz[jdx].Ssecs[0].Vu
					case xdx == len(f.Xs):
					//load with xdx -1, xdx joist
					jdx := fmt.Sprintf("%d-%d-x-o",xdx-1, xdx)
					wu += f.Bmz[jdx].Ssecs[0].Vu
					default:
					//load with xdx-1, xdx and xdx, xdx+1 joists
					jdx := fmt.Sprintf("%d-%d-x-o",xdx, xdx+1)
					wu += f.Bmz[jdx].Ssecs[0].Vu
					jdx = fmt.Sprintf("%d-%d-x-o",xdx-1, xdx)
					wu += f.Bmz[jdx].Ssecs[0].Vu
				}
				bm := Bm{
					Title:gdx,
					Sname:f.Sname,
					Lsb:true,
					Lspan:f.Ohng,
					Dtyp:1,
					Store:true,
					Dsgn:true,
					Onism:f.Onism,
					Selfwt:true,
					Term:"dumb",
					Ldcases:[][]float64{{1,1,wu,0,f.Ohng-50.0,0,1}},
					Start:strt,
					End:end,
					Name:"c.girder",
				}
				e := BmDz(&bm)
				if e != nil{
					errstr += fmt.Sprintf("c.girder %s error %s\n",gdx, fmt.Sprint(e))
					ecount++
				} 
				f.Bmz[gdx] = bm
				f.Cgs = append(f.Cgs, gdx)
				f.Cmap[idx][1] = append(f.Cmap[idx][1],gdx)
				f.Idmap["c.girder"] = append(f.Idmap["c.girder"],gdx)
			}
		}
		//add intermediate joists between columns
		if xdx < len(f.Xs){
			switch ydx{
				case len(f.Ys):
				i2 := idx - len(f.Xs)
				jdx := getybmdx(i2, i2+len(f.Xs),i2+1,i2+len(f.Xs)+1)
				jdx += "-top"
				gdx := getbmdx(idx,idx+1)
				f.Bmz[gdx] = f.Bmz[jdx]
				f.Bms[gdx] = []string{}
				bm := f.Bmz[gdx]
				bm.Name = "j.girder"
				if f.Matread && !f.Gon[gdx]{bm.Ignore = true}
				f.Bmz[gdx] = bm
				f.Cmap[idx][0] = append(f.Cmap[idx][0],gdx)
				f.Cmap[idx+1][0] = append(f.Cmap[idx+1][0],gdx)
				f.Idmap["j.girder"] = append(f.Idmap["j.girder"],gdx)
				
				default:
				jdx := getybmdx(idx, idx + len(f.Xs),idx+1,idx+len(f.Xs)+1)
				jdx += "-bot"
				gdx := getbmdx(idx, idx+1)
				f.Bmz[gdx] = f.Bmz[jdx]
				f.Bms[gdx] = []string{}
				bm := f.Bmz[gdx]
				if f.Matread && !f.Gon[gdx]{bm.Ignore = true}
				bm.Name = "j.girder"
				f.Bmz[gdx] = bm
				f.Cmap[idx][0] = append(f.Cmap[idx][0],gdx)
				f.Cmap[idx+1][0] = append(f.Cmap[idx+1][0],gdx)	
				f.Idmap["j.girder"] = append(f.Idmap["j.girder"],gdx)
			}
		}
	}
	//what about overhanging clvr girder moments
	//TODO
	for col, cvec := range f.Cmap{
		// fmt.Println(ColorPurple,"at col",col,ColorReset)
		// fmt.Println("beams in x-",cvec[0])
		// fmt.Println("beams in y-",cvec[1])
		var pu, vx1, vy1, vx2, vy2 float64
		cname := "int"
		xdx := col%len(f.Xs)
		ydx := col/len(f.Xs) + 1
		if xdx == 0{
			xdx = len(f.Xs)
			ydx--
		}
		switch{
			case xdx == 1 || xdx == len(f.Xs):
			if ydx == 1 || ydx == len(f.Ys){
				cname = "corner"
			} else {
				cname = "edge"
			}
			case ydx == 1 || ydx == len(f.Ys):
			if xdx == 1 || xdx == len(f.Xs){
				cname = "corner"
			} else {
				cname = "edge"
			}
		}
		cname = fmt.Sprintf("%s-%v-%v",cname,xdx,ydx)
		var cdeg int
		for i, bxstr := range cvec[0]{
			
			if bmx, ok := f.Bmz[bxstr]; ok{
				if len(bmx.Ssecs) > 0 && !bmx.Ignore{
					cdeg++
					switch i{
						case 0:
						vx1 = bmx.Ssecs[0].Vu
						case 1:
						vx2 = bmx.Ssecs[0].Vu
					}
				}
			}
		}
		for i, bystr := range cvec[1]{
			if bmy, ok := f.Bmz[bystr]; ok{
				if len(bmy.Ssecs) > 0 && !bmy.Ignore{
					cdeg++
					switch i{
						case 0:
						vy1 = bmy.Ssecs[0].Vu
						case 1:
						vy2 = bmy.Ssecs[0].Vu
					}
				}
			}
		}
		// vx1 = f.Bmz[cvec[0][0]].Ssecs[0].Vu
		// vx2 = f.Bmz[cvec[0][1]].Ssecs[0].Vu
		// vy1 = f.Bmz[cvec[1][0]].Ssecs[0].Vu
		// vy2 = f.Bmz[cvec[1][1]].Ssecs[0].Vu
		
		// fmt.Println("vx1 vx2 vy1 vy2",vx1,vx2,vy1,vy2)
		pu  = math.Round(vx1 + vx2 + vy1 + vy2)
		vx := math.Round(math.Abs(vx1-vx2))
		vy := math.Round(math.Abs(vy1-vy2))
		pu = float64(len(f.Zs)-1)*pu

		//designing for the bottom col
		c := Col{
			Title:fmt.Sprintf("%d",col),
			Sname:f.Sname,
			Onism:f.Onism,
			Code:2,
			Term:"dumb",
			Store:true,
			Lspan:f.Ht,
			Lx:f.Ht,
			Ly:f.Ht,
			Tx:1.0,
			Ty:1.0,
			Pu:pu,
			Vx:vx,
			Vy:vy,
			Vax:0,
			Vbx:vx,
			Vay:0,
			Vby:vy,
			Dsgn:true,
			Pfac:1.0,
			H1:f.Ht,
			H2:f.Ht,
			Name:cname,
			Deg:cdeg,
			Verbose:false,
		}
		
		//mera joota hai japani
		//earlier only bs449 was written
		
		c.Code = 1; c.Dtyp = 1
		c.Calctyp = 1; c.Sname = "i";
		c.Onism = f.Onism

		//if bms are ub col is uc?
		if f.Sname == "ub"{c.Sname = "uc"}

		if vx + vy < 1.0{
			//design as an axially loaded column
			c.Dtyp = 0
		}
		if f.Yeolde{c.Dtyp = 0; c.Code = 2}
		e := ColDz(&c)
		
		if f.Matread && !f.Colon[c.Title] || cdeg == 0{c.Ignore = true}
		if e != nil{
			errstr += fmt.Sprintf("column %v error %s\n",col, fmt.Sprint(e))
			ecount++
		}
		f.Colz[c.Title] = c
		//design column base plate
		if f.Dzcbs && e == nil{
			cb := Cbs{
				Sname:f.Sname,
				Sdx:c.Ssecs[0].Sdx,
				Code:1,
				Fck:f.Fck,
				Pu:c.Pu/1e3,
				Title:c.Title,
			}
			e = SlbCbsDz(&cb)
			if e != nil{
				errstr += fmt.Sprintf("column %v base error %s\n",col, fmt.Sprint(e))
				ecount++
			} 
			f.Cbz[c.Title] = cb
		}
	}
	if ecount > 0{
		err = fmt.Errorf("error in steel floor design \n%s",errstr)
		return
	}
	
	f.Smap = make(map[string][]float64) 
	err = f.Draw()
	if err != nil{
		return
	}
	//joint design 
	f.Jntz = make(map[string]kass.Blt)
	f.Jnts = make(map[string][]string)
	//joist end connections
	for _,jdx := range f.Idmap["joist"]{
		jvec := strings.Split(jdx,"-")
		jloc := jvec[len(jvec)-1]
		if jloc == "bot" || jloc == "top"{
			continue
		}
		bm := f.Bmz[jdx]
		if bm.Ignore{continue}
		msec := bm.Ssecs[0]
		//fmt.Println(ColorRed,"jdx, vu, rl, rr",jdx,msec.Vu/1e3,msec.Rl/1e3,msec.Rr/1e3,ColorReset)
		lsec := f.Bmz[bm.Lsup].Ssecs[0]
		rsec := f.Bmz[bm.Rsup].Ssecs[0]
		kl := jdx + "-dac"
		t1 := msec.Tw
		t2 := lsec.Tw
		sdx := bm.Lsup
		sdims := lsec.Dims
		if rsec.Tw < t2{
			t2 = rsec.Tw
			sdims = rsec.Dims
			sdx = bm.Rsup
		}
		bk := kass.Blt{
			Title:kl,
			Vdu:msec.Vu,
			Dia: 20.0,
			Grade:4.6,
			Name:"dac",
			Cloc:"web",
			Ctyp:2,
			Fup:410.0,
			T1:t1,
			T2:t2,
			T3:10.0,
			Mdims:[][]float64{
				msec.Dims,
				sdims,
			},
			Mdxs:[]string{jdx,sdx},
			Verbose:false,
			Term:f.Term,
		}
		err = kass.BltDz(&bk)
		if err != nil{
			fmt.Println(ColorRed,err,ColorReset)
		} else {
			f.Jntz[kl] = bk
			f.Jntord = append(f.Jntord, kl)
		}
	}
	//int. column joists
	
	for _,jdx := range f.Idmap["j.girder"]{
		bm := f.Bmz[jdx]
		msec := bm.Ssecs[0]
		
		if bm.Ignore{continue}
		//fmt.Println(ColorRed,"jdx, vu, rl, rr",jdx,msec.Vu/1e3,msec.Rl/1e3,msec.Rr/1e3,ColorReset)
		lsec := f.Colz[bm.Lsup].Ssecs[0]
		rsec := f.Colz[bm.Rsup].Ssecs[0]
		kl := jdx + "-dac" 
		t1 := msec.Tw
		t2 := lsec.Tw
		sdx := bm.Lsup
		sdims := lsec.Dims
		if rsec.Tw < t2{
			t2 = rsec.Tw
			sdims = rsec.Dims
			sdx = bm.Rsup
		}
		bk := kass.Blt{
			Title:kl,
			Vdu:msec.Vu,
			Dia: 20.0,
			Grade:4.6,
			Name:"dac",
			Cloc:"web",
			Ctyp:1,
			Fup:410.0,
			T1:t1,
			T2:t2,
			T3:10.0,
			Mdims:[][]float64{
				msec.Dims,
				sdims,
			},
			Mdxs:[]string{jdx,sdx},
			Verbose:false,
			Term:f.Term,
		}
		err = kass.BltDz(&bk)
		if err != nil{
			fmt.Println(ColorRed,err,ColorReset)
		} else {
			f.Jntz[kl] = bk
			f.Jntord = append(f.Jntord, kl)
		}
	}
	//girder end connections - on column flange 
	
	for _,jdx := range f.Idmap["girder"]{
		bm := f.Bmz[jdx]
		msec := bm.Ssecs[0]
		
		if bm.Ignore{continue}
		//fmt.Println(ColorRed,"jdx, vu, rl, rr",jdx,msec.Vu/1e3,msec.Rl/1e3,msec.Rr/1e3,ColorReset)
		lsec := f.Colz[bm.Lsup].Ssecs[0]
		rsec := f.Colz[bm.Rsup].Ssecs[0]
		kl := jdx + "-dac" 
		t1 := msec.Tf
		t2 := lsec.Tf
		sdx := bm.Lsup
		sdims := lsec.Dims
		if rsec.Tf < t2{
			t2 = rsec.Tf
			sdims = rsec.Dims
			sdx = bm.Rsup
		}
		bk := kass.Blt{
			Title:kl,
			Vdu:msec.Vu,
			Dia: 20.0,
			Grade:4.6,
			Name:"dac",
			Cloc:"flange",
			Ctyp:1,
			Fup:410.0,
			T1:t1,
			T2:t2,
			T3:10.0,
			Mdims:[][]float64{
				msec.Dims,
				sdims,
			},
			Mdxs:[]string{jdx,sdx},
			Verbose:false,
			Term:f.Term,
		}
		err = kass.BltDz(&bk)
		if err != nil{
			fmt.Println(ColorRed,err,ColorReset)
		} else {
			f.Jntz[kl] = bk
			f.Jntord = append(f.Jntord, kl)
		}
	}
	//TODO overhanging joists/clvr girder moment joint
	
	

	//TODO add length of bracing to table
	//TODO bracing connection design
	//TODO subtract net area of bolts from cols/beams
	err = f.Table(f.Verbose)
	if err != nil{
		return
	}
	
	if len(f.Brc) > 0{
		err = f.BrcDz()
		if err != nil{
			return
		}
	}
	if f.Dzbs{
		//design pedestals + footing
		err = f.BsDz()
		if err != nil{
			return
		}
		
		//fmt.Println("footing quants- vrcc, wstl, afw",f.Fquant[0],f.Fquant[1],f.Fquant[2])
		//fkost := f.Fquant[0]*12500 + f.Fquant[1] * 100 + f.Fquant[2] * 75.0
		//fmt.Println("footing kost",fkost,"rupeeses")
	} 
	
	//reportz
	return
}

//Table generates an ascii report for a stl flr
func (f *StlFlr) Table(printz bool)(err error){
	rezstr := new(strings.Builder)

	hdr := fmt.Sprintf("%s\n%s steel floor report\ndate-%s\n%s\n",ColorYellow,f.Title,time.Now().Format("2006-01-02"),ColorReset)
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorRed)
	hdr = ""
	hdr += fmt.Sprintf("floor geometry\n xs %.0f \n ys %0.0f\n zs % .f\n", f.Xs, f.Ys, f.Zs)
	hdr += fmt.Sprintf("joist spacing - %.0f mm \noverhang - % 0.0f mm\n",f.Jspc, f.Ohng)
	hdr += fmt.Sprintf("loads - dl %.3f kN/m2, ll %0.3f kN/m2\n", f.DL, f.LL)
	
	rezstr.WriteString(hdr)
	
	rezstr.WriteString(ColorPurple)
	table := tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"columns")
	
	table.SetHeader([]string{"id","section","length(m)","weight(kg/mt)","pu(kn)","vdx(kn)","vdy(kn)","flip?"})
	for _, idx := range f.Idmap["col"]{
		col := f.Colz[idx]
		row := fmt.Sprintf("%s ,%s, %.3f, %.0f, %.0f, %.0f, %.0f, %v",idx, col.Ssecs[0].Sstr, f.Hcol, col.Ssecs[0].Wt/9.81, col.Pu/1e3, col.Vx/1e3, col.Vy/1e3,col.Flip)
		table.Append(strings.Split(row,","))
	}
	table.Render()
	
	//girder, j.girder, joist, o.joist, c.girder

	//girder
	rezstr.WriteString(ColorBlue)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"girders")
	table.SetHeader([]string{"id","section","length(m)","weight(kg/mt)","vu(kn)","mu(kn.m)","dmax(mm)"})
	for _, idx := range f.Idmap["girder"]{
		bm := f.Bmz[idx]
		row := fmt.Sprintf("%s , %s, %.3f, %.0f, %.0f, %.0f, %.0f",idx,bm.Ssecs[0].Sstr, bm.Lspan, bm.Ssecs[0].Wt,bm.Ssecs[0].Vu/1e3,bm.Ssecs[0].Mu/1e6,bm.Ssecs[0].Dmax)
		table.Append(strings.Split(row,","))
	}
	table.Render()

	//j.girder
	rezstr.WriteString(ColorCyan)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"j.girder")
	table.SetHeader([]string{"id","section","length(m)","weight(kg/mt)","vu(kn)","mu(kn.m)","dmax(mm)"})
	for _, idx := range f.Idmap["j.girder"]{
		bm := f.Bmz[idx]
		row := fmt.Sprintf("%s , %s, %.3f, %.0f, %.0f, %.0f, %.0f",idx,bm.Ssecs[0].Sstr, bm.Lspan, bm.Ssecs[0].Wt,bm.Ssecs[0].Vu/1e3,bm.Ssecs[0].Mu/1e6,bm.Ssecs[0].Dmax)
		table.Append(strings.Split(row,","))
	}
	table.Render()

	//joist
	rezstr.WriteString(ColorPurple)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"joists")
	table.SetHeader([]string{"id","section","length(m)","weight(kg/mt)","vu(kn)","mu(kn.m)","dmax(mm)"})
	for _, idx := range f.Idmap["joist"]{
		bm := f.Bmz[idx]
		row := fmt.Sprintf("%s , %s, %.3f, %.0f, %.0f, %.0f, %.0f",idx,bm.Ssecs[0].Sstr, bm.Lspan, bm.Ssecs[0].Wt,bm.Ssecs[0].Vu/1e3,bm.Ssecs[0].Mu/1e6,bm.Ssecs[0].Dmax)
		table.Append(strings.Split(row,","))
	}
	table.Render()

	//o.joist
	rezstr.WriteString(ColorCyan)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"o.joists")
	table.SetHeader([]string{"id","section","length(m)","weight(kg/mt)","vu(kn)","mu(kn.m)","dmax(mm)"})
	for _, idx := range f.Idmap["o.joist"]{
		bm := f.Bmz[idx]
		row := fmt.Sprintf("%s , %s, %.3f, %.0f, %.0f, %.0f, %.0f",idx,bm.Ssecs[0].Sstr, bm.Lspan, bm.Ssecs[0].Wt,bm.Ssecs[0].Vu/1e3,bm.Ssecs[0].Mu/1e6,bm.Ssecs[0].Dmax)
		table.Append(strings.Split(row,","))
	}
	table.Render()

	//c.girder
	rezstr.WriteString(ColorBlue)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"c.girders")
	table.SetHeader([]string{"id","section","length(m)","weight(kg/mt)","vu(kn)","mu(kn.m)","dmax(mm)"})
	for _, idx := range f.Idmap["c.girder"]{
		bm := f.Bmz[idx]
		row := fmt.Sprintf("%s ,%s, %.3f, %.0f, %.0f, %.0f, %.0f",idx,bm.Ssecs[0].Sstr, bm.Lspan, bm.Ssecs[0].Wt,bm.Ssecs[0].Vu/1e3,bm.Ssecs[0].Mu/1e6,bm.Ssecs[0].Dmax)
		table.Append(strings.Split(row,","))
	}
	table.Render()

	//base plate
	if f.Dzbs{
		rezstr.WriteString(ColorYellow)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true,"base plates")
		table.SetHeader([]string{"column","length(c.H)","breath(c.B)","thickness"})
		for i := range f.Coords{
			idx := fmt.Sprintf("%v",i+1)
			cb := f.Cbz[idx]
			row := fmt.Sprintf("%s ,%.0f, %.0f, %.0f",idx, cb.L, cb.B,cb.Ts)
			table.Append(strings.Split(row,","))
		}
		table.Render()
	}
	//connections 'joints'
	rezstr.WriteString(ColorRed)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"joint summary")
	table.SetHeader([]string{"id","dia","grade","nb mem","nb sup","angle","depth(mm)","vu(kn)"})
	cwt := 0.0
	smap := make(map[string][]float64)
	for _, kdx := range f.Jntord{
		blt := f.Jntz[kdx]
		depth := blt.Dims[0]
		width := blt.Dims[1]
		thick := blt.Dims[2]
		astr := fmt.Sprintf("%.fx%.fx%.f",width,width,thick)
		wt := 4.0 * (width * thick * depth) * 7850 * 1e-9
		tsa := 8.0 * (width * depth) 
		if _, ok := smap[astr]; ok{
			smap[astr][0] += 2.0*depth/1e3
			smap[astr][1] += wt
			smap[astr][2] += tsa
		} else {
			smap[astr] = []float64{2.0*depth/1e3,wt,tsa}
		}
		cwt += wt
		row := fmt.Sprintf("%s, %.f, %.1f, %v, %v, %s, %.f, %.f",kdx,blt.Dia,blt.Grade,blt.Nb,blt.Nb2,astr,depth,blt.Vdu/1e3)
		table.Append(strings.Split(row,","))
	}
	table.Render()
	//summary
	rezstr.WriteString(ColorGreen)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"quantity summary")
	table.SetHeader([]string{"section","length","net weight(kg)","tsa(m2)"})
	for idx, vec := range f.Smap{
		row := fmt.Sprintf("%s ,%.0f, %.0f, %.0f",idx, vec[0], vec[1],vec[2])
		table.Append(strings.Split(row,","))
	}
	for idx, vec := range f.Psmap{
		row := fmt.Sprintf("%s ,%.3fm2, %.0f, %.0f",idx, vec[0], vec[1],vec[2])
		table.Append(strings.Split(row,","))
	}
	//forget conn for now
	// for idx, vec := range smap{
	// 	row := fmt.Sprintf("%s ,%.0f, %.0f, %.0f",idx, vec[0], vec[1],vec[2])
	// 	table.Append(strings.Split(row,","))
	// }
	table.Render()
	rezstr.WriteString(ColorPurple)
	table = tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"consumption")
	table.SetHeader([]string{"net weight(kg)","net BUA(sft)","nfloors","kg/sft"})
	row := fmt.Sprintf("%0.0f, %0.0f, %v, %.2f",f.Wt, f.Bua, len(f.Zs)-1,(f.Wt)/f.Bua)
	table.Append(strings.Split(row,","))
	table.Render()
	rezstr.WriteString(ColorReset)
	f.Report = rezstr.String()
	if printz{
		fmt.Println(f.Report)
	}
	//fmt.Println(f.Report)
	return
}

func (f *StlFlr) Draw()(err error){
	var data, vdata, ldata string
	//all lengths in meters, METERS, el metere
	hcol := (f.Zs[len(f.Zs)-1] - f.Zs[0])/1000.0
	f.Hcol = hcol
	f.Wt = 0.0
	nflrs := float64(len(f.Zs)-1)
	lx := f.Xs[len(f.Xs)-1] + 2.0 * f.Ohng
	ly := f.Ys[len(f.Ys)-1] + 2.0 * f.Ohng
	if !f.Matread{
		f.Bua = lx * ly*10.76*nflrs/1e6
	} else {
		f.Bua = f.Bua * nflrs
	}
	//draw col nodes
	for i, idx := range f.Cols{
		cdx := fmt.Sprintf("%v",idx)
		if f.Colz[cdx].Ignore{continue}
		pt := f.Coords[idx-1]
		sec := f.Colz[cdx].Ssecs[0]
		ss := sec.Sec
		if f.Colz[cdx].Flip{ss = kass.FlipX(ss)}
		ss = kass.SecTranslate(ss,pt[0]-ss.Prop.Xc,pt[1]-ss.Prop.Yc)
		ss.Draw("")
		vdata += ss.Data[0]
		data += fmt.Sprintf("%f %f %v\n",pt[0],pt[1],i+1)
		
		ldata += fmt.Sprintf("%f %f %s\n",pt[0]+100,pt[1]+100,cdx)
		ldata += fmt.Sprintf("%f %f %s\n",pt[0]+200,pt[1]+200,sec.Sstr)
		if _, ok := f.Smap[sec.Sstr]; !ok{
			f.Smap[sec.Sstr] = []float64{0,0,0,0,0,0}
			
		}
		f.Smap[sec.Sstr][0] += hcol
		f.Smap[sec.Sstr][1] += hcol * sec.Wt/9.81
		f.Smap[sec.Sstr][2] += hcol * ss.Prop.Perimeter/1000.0
		f.Wt += hcol * sec.Wt/9.81
	}
	data += "\n\n"
	//draw girder/beam centerlines (dotted)
	for _, bm := range f.Bmz{
		if bm.Ignore{
			continue
		}
		sec := bm.Ssecs[0]
		bcol := 0
		lbm := bm.Lspan/1000.0
		if _, ok := f.Smap[sec.Sstr]; !ok{
			f.Smap[sec.Sstr] = []float64{0,0,0,0,0,0}
			
		}
		switch bm.Name{
			case "girder":
			bcol = 0
			data += fmt.Sprintf("%f %f %f %f %v\n",bm.Start[0],bm.Start[1],bm.End[0]-bm.Start[0],bm.End[1]-bm.Start[1],bcol)
			ldata += fmt.Sprintf("%f %f %s\n",(bm.Start[0]+bm.End[0])/2.0,(bm.Start[1]+bm.End[1])/2.0,sec.Sstr)
			//mdx int, d float64, pb, pe []float64
			vdata += kass.DrawRectView(bcol,sec.B,bm.Start,bm.End)
			f.Smap[sec.Sstr][0] += lbm
			f.Smap[sec.Sstr][1] += lbm * sec.Wt/9.81
			f.Smap[sec.Sstr][2] += lbm * sec.Sec.Prop.Perimeter/1000.0
			f.Wt += lbm * sec.Wt/9.81
			case "joist":
			bcol = 1
			jsp := strings.Split(bm.Title,"-")
			cdx, _ := strconv.Atoi(jsp[0])
			//idx, idx + len(f.Xs),idx+1,idx+len(f.Xs)+1
			if jsp[4] == "mid"{	
				nbm := 0.0
				ydx := cdx/len(f.Xs)
				x1 := bm.Start[0]
				y1 := bm.Start[1]
				x2 := bm.End[0]
				y2 := bm.End[1]
				//data += fmt.Sprintf("%f %f %f %f %v\n",x1,y1,x2-x1,y2-y1,bcol+3)
				//ldata += fmt.Sprintf("%f %f %s\n",(x1+x2)/2.0,(y1+y2)/2.0,sec.Sstr)
				for _, spc := range f.Jspcs[ydx]{
					data += fmt.Sprintf("%f %f %f %f %v\n",x1,y1+spc,x2-x1,y2-y1,bcol)
					ldata += fmt.Sprintf("%f %f %s\n",(x1+x2)/2.0,y1+spc,sec.Sstr)
					nbm += 1.0
				}
				f.Smap[sec.Sstr][0] += lbm
				f.Smap[sec.Sstr][1] += lbm * sec.Wt/9.81
				f.Smap[sec.Sstr][2] += lbm * sec.Sec.Prop.Perimeter/1000.0
				f.Wt += nbm * lbm * sec.Wt/9.81
			} 
			case "o.joist":
			bcol = 2
			data += fmt.Sprintf("%f %f %f %f %v\n",bm.Start[0],bm.Start[1],bm.End[0]-bm.Start[0],bm.End[1]-bm.Start[1],bcol)
			ldata += fmt.Sprintf("%f %f %s\n",(bm.Start[0]+bm.End[0])/2.0,(bm.Start[1]+bm.End[1])/2.0,sec.Sstr)
			//vdata += kass.DrawRectView(bcol,sec.B,bm.Start,bm.End)
			
			//multiply this by 2 (CHANGE WHEN IT DRAWS BOTH SIDES)
			f.Smap[sec.Sstr][0] += 2.0 * lbm 
			f.Smap[sec.Sstr][1] += 2.0 * lbm * sec.Wt/9.81
			f.Smap[sec.Sstr][2] += 2.0 * lbm * sec.Sec.Prop.Perimeter/1000.0
			f.Wt += 2.0 * lbm * sec.Wt/9.81
			case "c.girder":
			bcol = 3
			data += fmt.Sprintf("%f %f %f %f %v\n",bm.Start[0],bm.Start[1],bm.End[0]-bm.Start[0],bm.End[1]-bm.Start[1],bcol)
			ldata += fmt.Sprintf("%f %f %s\n",(bm.Start[0]+bm.End[0])/2.0,(bm.Start[1]+bm.End[1])/2.0,sec.Sstr)
			vdata += kass.DrawRectView(bcol,sec.B,bm.Start,bm.End)

			f.Smap[sec.Sstr][0] += 2.0 * lbm 
			f.Smap[sec.Sstr][1] += 2.0 * lbm * sec.Wt/9.81
			f.Smap[sec.Sstr][2] += 2.0 * lbm * sec.Sec.Prop.Perimeter/1000.0
			f.Wt += 2.0 * lbm * sec.Wt/9.81
			
			case "j.girder":
			bcol = 4
			data += fmt.Sprintf("%f %f %f %f %v\n",bm.Start[0],bm.Start[1],bm.End[0]-bm.Start[0],bm.End[1]-bm.Start[1],bcol)
			ldata += fmt.Sprintf("%f %f %s\n",(bm.Start[0]+bm.End[0])/2.0,(bm.Start[1]+bm.End[1])/2.0,sec.Sstr)
			vdata += kass.DrawRectView(bcol,sec.B,bm.Start,bm.End)
			f.Smap[sec.Sstr][0] += lbm 
			f.Smap[sec.Sstr][1] += lbm * sec.Wt/9.81
			f.Smap[sec.Sstr][2] += lbm * sec.Sec.Prop.Perimeter/1000.0
			f.Wt += lbm * sec.Wt/9.81
			
		}
	}
	//draw column base plates
	for id, cb := range f.Cbz{
		idx,_ := strconv.Atoi(id)
		dy := cb.L; dx := cb.B
		col := f.Colz[id]
		if col.Ignore{continue}
		if col.Flip{
			dy = cb.B; dx = cb.L
		}
		cbsec := kass.SecGen(1,[]float64{dx,dy})
		tx := f.Coords[idx-1][0] - cbsec.Prop.Xc
		ty := f.Coords[idx-1][1] - cbsec.Prop.Yc
		cbsec = kass.SecTranslate(cbsec, tx, ty)
		cbsec.Draw("")
		vdata += cbsec.Data[0]
		sstr := fmt.Sprintf("%vmm plt",cb.Ts)
		if _, ok := f.Psmap[sstr]; !ok{
			f.Psmap[sstr] = []float64{0,0,0,0,0,0}
		}
		f.Psmap[sstr][0] += dx * dy/1e6
		f.Psmap[sstr][1] += dy * dy * cb.Ts * 7850/1e9
		f.Psmap[sstr][2] += dy * dy * 2/1e6
		f.Wt += dy * dy * cb.Ts * 7850/1e9
	}
	data += "\n\n"; vdata += "\n\n"
	data += vdata; data += ldata
	if f.Web{
		f.Folder = "web"
	}
	txtplot, err := draw.Draw(data, "drawstlflr.gp", f.Term, f.Folder, f.Title, f.Title, "","","")
	if err != nil{
		return
	}
	if f.Term == "dumb"{fmt.Println(txtplot)}
	return
}


//PDz designs pedestals + footing
func (f *StlFlr) BsDz()(err error){
	f.Pdz = make(map[string]mosh.RccCol)
	f.Ftngz = make(map[string]mosh.RccFtng)
	f.Fquant = []float64{0,0,0}
	ecount := 0
	var errstr string
	//mosh is in kn-m (apparently written by another person)
	for idx, cb := range f.Cbz{
		if f.Colz[idx].Ignore{
			//fmt.Println("ignorning",idx)
			continue
		}
		h := math.Round((cb.L + 50.0)/75.0)*75.0
		b := math.Round((cb.B + 50.0)/75.0)*75.0
		
		rc := mosh.RccCol{
			Title:f.Title+"-pd-"+idx,
			Fck:f.Fck,
			Fy:500.0,
			Cvrt:60.0,
			Cvrc:60.0,
			B:b,
			H:h,
			Styp:1,
			Dtyp:0,
			Pu:cb.Pu,
			Code:1,
			Lspan:f.Df/1000.0,
			Term:f.Term,
			Verbose:false,
		}
		e := mosh.ColDesign(&rc)
		if e != nil{
			ecount++
			errstr += fmt.Sprintf("column %v base pedestal error %s\n",idx, fmt.Sprint(e))
			
		} else {
			pltstr := rc.PlotColDet()

			f.Breport += rc.Report
			f.Breport += pltstr
			f.Pdz[idx] = rc
			f.Fquant[0] += rc.Vrcc
			f.Fquant[1] += rc.Wstl
			f.Fquant[2] += rc.Afw
			
			rf := mosh.RccFtng{
				Title:f.Title+"-ftng-"+idx,
				Colx: rc.B/1e3,
				Coly: rc.H/1e3,
				Df: f.Df/1e3,
				Eo: 0.075,
				Fck: f.Fck,
				Fy: 500.0,
				Sbc: f.Sbc,
				Pgck: 25.0,
				Pgsoil: 15.0,
				Nomcvr: 0.075,
				Dmin: 0.30,
				Pus: []float64{cb.Pu},
				Mxs: []float64{0.0},
				Mys: []float64{0.0},
				Psfs: []float64{1.5},
				Typ: 1,
				Shape: "rect",
				Sloped: false,
				Dlfac: false,
				Verbose:f.Verbose,
				Term:f.Term,
			}
			e = mosh.FtngDzRojas(&rf)
			if e != nil{
				fmt.Println(e)
				ecount++
				errstr += fmt.Sprintf("column %v footing error %s\n",idx, fmt.Sprint(e))
			} else {
				pltstr := mosh.PlotFtngDet(&rf)
				if rf.Term == "dumb"{
					fmt.Println(pltstr)
				}
				f.Breport += rf.Report
				f.Ftngz[idx] = rf
				f.Fquant[0] += rf.Vrcc
				f.Fquant[1] += rf.Wstl
				f.Fquant[2] += rf.Afw
			}
		}
	}
	if ecount > 0{
		err = fmt.Errorf("error in floor base design \n%s",errstr)
		return
	}
	return
}

//BrcDz designs braced bays brc[0] = xdx, brc[1] = ydx 
func (f *StlFlr) BrcDz()(err error){
		if len(f.Brc) < 2{
		err = fmt.Errorf("brace locations not specified")
		return
	}
	//calc pdmax
	l := f.Xs[len(f.Xs)-1]
	w := f.Ys[len(f.Ys)-1]
	if l < w{
		w = f.Xs[len(f.Xs)-1]
		l = f.Ys[len(f.Ys)-1]
	}
	h := f.Zs[len(f.Zs)-1]
	xidx := f.Brc[0]; yidx := f.Brc[1]
	if xidx + 1 > len(f.Xs)-1 || yidx + 1 > len(f.Ys) - 1 || xidx < 0 || yidx < 0{
		err = fmt.Errorf("invalid brace indices xidx %v yidx %v",xidx, yidx)
		return
	}
	dx := f.Xs[xidx+1]-f.Xs[xidx]
	dy := f.Ys[yidx+1]-f.Ys[yidx]
	
	pd, _ , _ := kass.GetPdWall(f.Vz, h, w, l, f.Cpi)
	//calc wall area in x (trib area for y frame)
	triby := f.Xs[len(f.Xs)-1] * f.Ht/1e6
	tribx := f.Ys[len(f.Ys)-1] * f.Ht/1e6

	var kax, kay float64
	for i, trib := range []float64{triby, tribx}{
		ka := 0.0	
		switch {
		case trib < 10.0:
			ka = 1.0
		case trib < 25.0:
			darea := 25.0 - trib
			ka  = 0.9 + darea * 0.1/15.0
		case trib < 100.0:
			darea := 100.0 - trib
			ka = 0.8 + darea * 0.1/75.0 
		default:
			ka = 0.8
		}
		switch i{
			case 0:
			triby = triby * ka
			kay = ka
			case 1:
			tribx = tribx * ka
			kax = ka
		}
	}
			
	//build truss model in x
	//pick columns (start from 1)
	cx1 := f.Colz[fmt.Sprintf("%v",xidx+1)]
	cx2 := f.Colz[fmt.Sprintf("%v",xidx+2)]
	//beam connected from cx1, cx2
	g1 := getbmdx(xidx+1, xidx+2)
	//build truss model in y
	icy1 := fmt.Sprintf("%v",yidx*len(f.Xs)+1)
	icy2 := fmt.Sprintf("%v",(yidx+1)*len(f.Xs)+1)

	g2 := getbmdx(yidx*len(f.Xs)+1, (yidx+1)*len(f.Xs)+1)
	cy1 := f.Colz[icy1]
	cy2 := f.Colz[icy2]
	lx := f.Xs[len(f.Xs)-1] + 2.0 * f.Ohng
	ly := f.Ys[len(f.Ys)-1] + 2.0 * f.Ohng
	punx := lx*ly*(f.DL+f.LL)*0.5/1e5
	
	//use 2x75x75x10 angle sections for bracing, area = 1410*2 mm2
	//150x150x10 area = 2903 mm2
	//build x and y bracing analysis models
	tbx := kass.Model{
		Id:f.Title+"-tbx",
		Frmstr:"2dt",
		Term:f.Term,
		Web:f.Web,
		Units:"nmm",
		Calc:true,
		Supports:[][]int{
			{1, -1, -1},
			{2, -1, -1},
		},
		Cp:[][]float64{
			{cx1.Ssecs[0].Area},
			{cx2.Ssecs[0].Area},
			{f.Bmz[g1].Ssecs[0].Area},
			{2820},
		},
		Em:[][]float64{{200000.0}},
	}

	tby := kass.Model{
		Id:f.Title+"-tby",
		Frmstr:"2dt",
		Term:f.Term,
		Web:f.Web,
		Units:"nmm",
		Calc:true,
		Supports:[][]int{
			{1, -1, -1},
			{2, -1, -1},
		},
		Cp:[][]float64{
			{cy1.Ssecs[0].Area},
			{cy2.Ssecs[0].Area},
			{f.Bmz[g2].Ssecs[0].Area},
			{2820},
		},
		Em:[][]float64{{200000.0}},
	}
	

	for i, z := range f.Zs{
		p1 := []float64{0.0,z}
		p2 := []float64{dx, z}
		tbx.Coords = append(tbx.Coords, [][]float64{p1,p2}...)
		p1 = []float64{0.0,z}
		p2 = []float64{dy, z}
		tby.Coords = append(tby.Coords, [][]float64{p1,p2}...)
		
		if i > 0{
			//add c1, c2
			//add g1
			i1 := i * 2 + 1
			i2 := i1 + 1
			
			c1 := []int{i1-2, i1, 1, 1, 3}
			c2 := []int{i2-2, i2, 1, 2, 3}
			g1 := []int{i1,i2,1,3,3}
			b1 := []int{i1-2, i2, 1,4,3}
			if i % 2 != 0{
				b1 = []int{i2-2, i1,1,4,3}
			}
			tbx.Mprp = append(tbx.Mprp, [][]int{c1, c2, g1, b1}...)
			tby.Mprp = append(tby.Mprp, [][]int{c1, c2, g1, b1}...)
			
			//add wind pu/jload at i1
			dz := z - f.Zs[i-1]
			tribx := f.Ys[len(f.Ys)-1] * dz * kax/1e6/2.0
			triby := f.Xs[len(f.Xs)-1] * dz * kay/1e6/2.0
			
			//should it be factored?
			pux := pd * tribx * 1.5
			puy := pd * triby * 1.5

			if i == len(f.Zs)-1{
				pux = pux/2.0
				puy = puy/2.0
			}
			//fmt.Println(ColorGreen,"tribx,z,dz,pux fac",tribx,z,dz,pux+punx,"nmm")
			//fmt.Println(ColorCyan,"triby,z,dz,puy fac",triby,z,dz,puy+punx,"nmm",ColorReset)
			
			tbx.Jloads = append(tbx.Jloads, []float64{float64(i1),pux+punx,0.0})
			
			tby.Jloads = append(tby.Jloads, []float64{float64(i1),puy+punx,0.0})
		}
		
	}
	//calc x model
	err = kass.CalcMod(&tbx, "2dt",tbx.Term,false)
	if err != nil{
		err = fmt.Errorf("%v error in x bracing model analysis",err)
		return
	}
	
	//calc y model
	err = kass.CalcMod(&tby, "2dt",tby.Term,false)
	if err != nil{
		err = fmt.Errorf("%v error in y bracing model analysis",err)
		return
	}
	//get max joint reactions
	reacts := []float64{tbx.Js[1].React[1],tbx.Js[2].React[1],tby.Js[1].React[1],tby.Js[2].React[1]}
	//fmt.Println("reactions",reacts)
	cdls := []float64{cx1.Ssecs[0].Pu*f.DL/(f.DL+f.LL),cx2.Ssecs[0].Pu*f.DL/(f.DL+f.LL),cy1.Ssecs[0].Pu*f.DL/(f.DL+f.LL),cy2.Ssecs[0].Pu*f.DL/(f.DL+f.LL)}
	//fmt.Println("col loads", cdls)
	//get ult. loads for uplift
	ulds := []float64{reacts[0]+cdls[0],reacts[1]+cdls[1],reacts[2]+cdls[2],reacts[3]+cdls[3]}
	//fmt.Println("uplift loads",ulds)
	//fmt.Println("report-",tbx.Report)
	umax := 0.0
	for _, uld := range ulds{
		if umax > uld{
			umax = uld
		}
	}
	//fmt.Println("max uplift load",umax/1e3,"kn")
	//default dia of anchor bolts = 20mm
	var lblt float64
	lblt, err = kass.AbltLen(math.Abs(umax), 25.0, 20.0)
	if err != nil{return}

	// fmt.Println("len, width, ht", l, w, h)
	// fmt.Println("brace width in x",dx,"mm")
	// fmt.Println("brace width in y",dy,"mm")
	// fmt.Println("design (basic) wind pressure pd",pd/1e3,"kn/m2")
	// fmt.Println("tribx, triby basic",tribx, triby)
	// fmt.Println("tribx, triby, kax, kay",tribx, triby, kax, kay)
	fmt.Println(ColorYellow,"length of anchor bolts",lblt,ColorReset)
	f.Lablt = lblt
	//add this to base detail, rem bolts can be 300 mm?
	//min embedment of anchor bolts - 12*dia (ASCE)
	return
}

//Edit edits/adds members to a stlfrm
func (f *StlFlr) Edit()(err error){
	running := true
	for running{
		choice := printmenu("edit floor",[]string{
			"add col between two cols",
			"add clvr beam to col",
			"delete panel",
			"exit",
		})
		switch choice{
			case 3:
			running = false
			case 0:
			f.AddColBet()
			
		}

	}
	return
}

func (f *StlFlr) AddColBet(){
	var c1, c2 int
	fmt.Println("enter c1 and c2:")
	_, err := fmt.Scan(&c1, &c2)
	if err != nil{
		fmt.Println(err)
		return
	} 
	fmt.Println("read c1, c2-",c1, c2)
	
	if c1 < 1 || c2 < 1 || c1 > len(f.Coords) || c2 > len(f.Coords){
		fmt.Println("error in column indices")
		return
	}
	fmt.Println("c1, c2")
	
}

//FlrDz designs an Fcflr struct
func FlrDz(f flay.Fcflr) (r FlrRez, err error){
	f.Load()
	f.Draw()
	r.Bms = make(map[int]Bm)
	r.Cols = make(map[int]Col)
	//build beam vec
	for b, bvec := range f.Bms{
		bdx := f.Bmdx[b]
		lspan := kass.Dist2d(f.Coords[b.I-1],f.Coords[b.J-1]) 
		bm := Bm{
			Lspan:lspan,
			Sname:"i",
			Code:1,
			Dtyp:1,
			Verbose:true,
		}
		btyp := f.Bmprp[b][0]
		fmt.Println(ColorRed,"beam idx",b,"type", btyp,ColorReset)
		
		switch btyp{
			case 0:
			//primary beam, get linked secondary beam reactions
			fmt.Println(ColorYellow,"primary beam",ColorReset)
			fmt.Println("list of secondary beam reactions-",f.Xbms[b])
			fmt.Println("beam points-",f.Xbmpts[b])
			case 1:
			//calc reactions and etx
			fmt.Println(ColorGreen,"secondary beam",ColorReset)
			for _, vec := range bvec{
				fmt.Println("M1,LTYP,W1,W2,L1,L2,LTYP")
				fmt.Println(vec)				
			}
			bm.Ldcases = append(bm.Ldcases, bvec...)
			err := BmDesign(&bm)
			if err != nil{
				fmt.Println("ERRORE,errore->",err)
			} 
			r.Bms[bdx] = bm
		}
		
	}
	fmt.Println("cols")
	return
}
			/*
	b := Bm{
		Lspan:4000.0,
		DL:15,
		Sname:"i",
		Nsecs:1,
		Dtyp:1,
		Code:1,
		Sdx:50,
		Lsb:true,
	}

			*/
