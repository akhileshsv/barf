package barf		  
			  
import (		  
	"fmt"		  
	"math"		  
	//"log"		  
)			  
			  
//GenLd generates loads for a 2d Trs gen struct
func (t *Trs2d) GenLd() (err error){
	//lrftr, lspan, spacing, dl, ll float64, slfwt, mtyp, ftyp, rfmat,trstyp int) (wdl, wll float64){
	//returns total dead load w for a 2d truss
	//if slope == 0 flat roof
	//npr - number of purlins
	//slfwt of wooden truss/steel truss
	//fmt.Println("CALCing")
	var wdl, wll, llred, wtruss float64 
	al := t.Rftrl * t.Spacing
	if t.Typ == 2 || t.Typ == 4{al = 2.0 * al}
	if t.Typ == 5 {t.Slope = 0}
	ang := math.Atan(1.0/t.Slope) * 180.0/math.Pi
	//fmt.Println(ColorRed,"ROOFER ANGLE->",ang,"DEGREES",ColorReset)
	if ang > 10.0{
		llred = (ang - 10.0) * 0.02 * 1e3
	}
	switch t.Ldcalc{
		case 0:
		//dead, live
		if t.PSFs == nil{
			t.PSFs = [][]float64{{1.5,1.5,0.0}}
		}
		case 1:
		//dead, live, max wind
		case 2:
		//dead, live, wind
		case 3:
		//dead, live, wind, seismic
	}
	if t.Ldcalc == 0 && t.PSFs == nil{
		switch t.Code{
			case 1:	
			t.PSFs = [][]float64{
				{1.5,1.5,0.0},
				{0.9,0.0,1.5},
			}
		}
	}
	switch t.Mtyp{	  
		case 3:	  
		//"timber"
		if t.DL > 0.0{wdl += t.DL} else {wdl += rfWt[t.Roofmat-1][0] * 1e3}
		//if t.Cmdz[0] == "addld"{wdl += rfWt[t.Roofmat-1][0] * al * 1e3}
		//fmt.Println(rfWt[rfmat-1][0] * al * 1e3)
		if t.LL > 0.0 {
			wll += t.LL 
		} else {
			wll += 0.75 * 1e3  - llred
			//wll += 1.5 * 1e3  - llred
			if wll < 400.0{wll = 400.0}
		}	  
		switch t.Mtyp{
			case 3:
			//timber truss
			switch t.Clctyp{
				default:
				//john wight formula
				wtruss = 0.75 * t.Span/304.8 * t.Spacing/304.8 * (1.0 + t.Span/304.8/10.0) * 4.5
				fmt.Println("formula wt in n->",0.75 * t.Span/304.8 * t.Spacing/304.8 * (1.0 + t.Span/304.8/10.0) * 4.5)
				
			  	case 1:
			  	//ricker formula WRENG
			  	wtruss = 0.5 * (1.0 * t.Span/10.0) * 0.454 * 10.76
			  	case 2:
			  	//abel 160 n/m2
			  	wdl += 160.0
			  	case 3:
			  	//sp 33
				case 4:
				//arre. calc volume and get self weight
			} 
			case 2:
			//steel truss
			switch t.Clctyp{
			  	case 1:
			  	//welded steel truss
			  	wdl += 20.0 + 6.6 * t.Span
			  	case 2:
			  	//bolted truss
			  	wtruss = (53.7 + 0.53 * t.Span * t.Spacing)
				case 3:
				//ketchum formula
			} 
		}	  
	}		  
	var pd float64	  
	var cpos, cneg []float64
	var wlcs map[int][]float64
	if t.Cpi == 0{t.Cpi = 0.5}
	var slope float64 
	if t.Typ == 5 {slope =0} else {slope = 1.0/t.Slope}
	if t.Typ == 1 || t.Typ == 3{
		pd, cpos, cneg, wlcs = wltable7(t.Vb, t.Height, t.Span, slope, t.Cpi)
	} else {	  
		pd, cpos, cneg, wlcs = wltable6(t.Vb, t.Height, t.Span, slope, t.Cpi)
	}
	//n/m2 to n/mm2
	fmt.Println("wdl",wdl,"wtruss",wtruss)
	fmt.Println("wll",wll)
	fmt.Println("truss tributary area->",al*1e-6,"m2")
	fmt.Println("nodal tributary area->",t.Purlinspc*t.Spacing*1e-6,"m2")
	fmt.Println("total load per truss->",(wdl+wll)*al*1e-6+wtruss,"n")
	fmt.Println("total nodes->",len(t.Tcns))

	fmt.Println("pd",pd)	  
	fmt.Println("cpos",cpos) 	    
 	fmt.Println("cneg",cneg)
	
	fmt.Println("wload cases->",wlcs)
	
	dltot := (wdl)*al*1e-6+wtruss
	lltot := (wll)*al*1e-6
	ntot := (float64(len(t.Tcns))) * 2.0 + 2.0
	//fmt.Println("load per node-> int->",dln*2.0/ntot,"end node->",wtot/ntot)
	//fmt.Println("rftrl->",t.Rftrl)
	//add tc node load cases - dead, live
	//wlcases  - 1 left, 2 right
	fmt.Println(ColorRed)
	fmt.Println("TCNODES->",t.Tcns)
	fmt.Println(ColorCyan)
	fmt.Println("BCNODES->",t.Bcns)
	fmt.Println(ColorReset)
	var nwl, nsl int
	theta := math.Atan(1.0/t.Slope)  //angle of roof in radians
	for _, node := range t.Tcns{
		//dead load
		dl := dltot*2.0/ntot; ll := lltot*2.0/ntot
		t.Jldsrv[1] = append(t.Jldsrv[1],[]float64{float64(node), 0.0, -dl})
		t.Jldsrv[2] = append(t.Jldsrv[2],[]float64{float64(node), 0.0, -ll})
		//wind load
		switch t.Typ{
			case 1:
			//l type - 
			case 2:
			//a type
			nwl = len(wlcs[1])
			switch{
				case math.Abs(t.Coords[node-1][0] - t.Span/2.0) < t.Span/2.0/1e3:
				//midnode
				for i, wl1 := range wlcs[1]{
					wl2 := wlcs[2][i]
					lc := (float64(i+1)+30.0)/10.0
					fx1 := (wl1 * t.Purlinspc * t.Spacing * 1e-6/2.0) * math.Sin(theta)
					fy1 := -(wl1 * t.Purlinspc * t.Spacing * 1e-6/2.0) * math.Cos(theta)
					fx2 := -(wl2 * t.Purlinspc * t.Spacing * 1e-6/2.0) * math.Sin(theta)
					fy2 := -(wl2 * t.Purlinspc * t.Spacing * 1e-6/2.0) * math.Cos(theta)
					t.Jldsrv[lc] = append(t.Jldsrv[lc],[]float64{float64(node),fx1, fy1})
					t.Jldsrv[lc] = append(t.Jldsrv[lc],[]float64{float64(node),fx2, fy2})
				}
				case t.Coords[node-1][0] < t.Span/2.0:
				//left nodes
				
				for i, wl1 := range wlcs[1]{
					lc := (float64(i+1)+30.0)/10.0
					fx1 := (wl1 * t.Purlinspc * t.Spacing * 1e-6) * math.Sin(theta)
					fy1 := -(wl1 * t.Purlinspc * t.Spacing * 1e-6) * math.Cos(theta)
					t.Jldsrv[lc] = append(t.Jldsrv[lc],[]float64{float64(node),fx1, fy1})
				}
				case t.Coords[node-1][0] > t.Span/2.0:
				//right nodes
				for i, wl1 := range wlcs[2]{
					lc := (float64(i+1)+30.0)/10.0
					fx1 := -(wl1 * t.Purlinspc * t.Spacing * 1e-6) * math.Sin(theta)
					fy1 := -(wl1 * t.Purlinspc * t.Spacing * 1e-6) * math.Cos(theta)
					t.Jldsrv[lc] = append(t.Jldsrv[lc],[]float64{float64(node),fx1, fy1})
				}
			}
			default:
			//
		}
	}
	//end nodes
	for _, node := range []int{t.Bcns[0],t.Bcns[len(t.Bcns)-1]}{
		dl := dltot/ntot; ll := lltot/ntot
		t.Jldsrv[1] = append(t.Jldsrv[1],[]float64{float64(node), 0.0, -dl})
		t.Jldsrv[2] = append(t.Jldsrv[2],[]float64{float64(node), 0.0, -ll})
		//wind load
		switch t.Typ{
			case 1:
			//l types - 
			case 2:
			//a type
			//left node
			switch node{
				case t.Bcns[0]:
				for i, wl1 := range wlcs[1]{
					lc := (float64(i+1)+30.0)/10.0
					fx1 := (wl1 * (t.Purlinspc + 2.0 * t.Ovrhng) * t.Spacing/2.0 * 1e-6) * math.Sin(theta)
					fy1 := -(wl1 * (t.Purlinspc + 2.0 * t.Ovrhng) * t.Spacing/2.0 * 1e-6) * math.Cos(theta)
					t.Jldsrv[lc] = append(t.Jldsrv[lc],[]float64{float64(node),fx1, fy1})
				}
				default:
				for i, wl1 := range wlcs[2]{
					lc := (float64(i+1)+30.0)/10.0
					fx1 := -(wl1 * (t.Purlinspc + 2.0 * t.Ovrhng) * t.Spacing/2.0 * 1e-6) * math.Sin(theta)
					fy1 := -(wl1 * (t.Purlinspc + 2.0 * t.Ovrhng) * t.Spacing/2.0 * 1e-6) * math.Cos(theta)
					t.Jldsrv[lc] = append(t.Jldsrv[lc],[]float64{float64(node),fx1, fy1})
				}
			}
			case 3:
			case 4:
		}
	}
	//nchk := 0
	//for i, ldcases := range t.Jldsrv{
	//	if i < 3{continue}
		//nchk ++
	//	fmt.Println("ld typ",i,"ldcases",ldcases)
	//}
	mldcs, jldcs := GenLdCombos(nwl, nsl, t.Jldsrv, t.Mldsrv, t.PSFs)
	t.Mod.Mldcs = mldcs
	t.Mod.Jldcs = jldcs
	return		  	    
}			  	    

/*			  	    
 			  	    
       +----------------+ 	  +--------------------------+
       |   dead	load   	       	  |          live load       |
       |self wt  purlin wt     	  | area load                |
       |sheet wt          	  +--------------------------+
       +----------------+ 	   +-------------------------+
   +-------------------------+ 	   | moar load
   |bent is a weird thing    |	   | 1. calc loads as per random formulae (dl esp)     	 
   +-------------------------+ 	   | 2. calc forces, design			  
       +-------------------+	   | 3. THEN ANALYZE AGAIN   |			  
       |  wind load        |   	   +-------------------------+ 	       	       	  
       |on roof            |	   p- delta analysis ?	       	      
       |func calcWHAAT            |    	       	       	 
       |                   |	  |    	       	             |
       |on column (girts)  |   	  |    	       	             |
       |func calcWHAAT     |   	  |                          |
       |                   |	  |                          |
       +-------------------+	  +--------------------------+
 				    		      
 				    	      X	      
 				    	     / \      
 				    	    /	\     
 		      /---	    	   /	 \    	       /--
 		  /---	  \---	    	  /	  \   	    /--	  \--
 	       /--	      \---  	 /	   \  	  /-	     \--
 	   /---			  \---  /	    \  /--		\--
      .	....	truss 1	     . 	.........  truss 2     .....truss 3	.	..............
 	  |	       ..	      |		     |		   	    |
 	  |	  		      |		     |		   	    |
      	  |	  		      |		     |		   	    |
      	  |	  	       	      |	       	     | 	       	       	    |  	    cheight
       	  |    	       	       	      |	       	     | 	       	       	    |  	       	   
 	  |	  		      |		     |		   	    |		  
 	--+---	  		    --+--	   --+---	   	 ---+---	  
 	  	  bent 1       	       	     2 	       	       	 3 	    |		  
 		  span 1               	     2 	       	       	 3     	  		  
       	       	       	       	       	       	       	       	   			  
     		       	       	       	       	       	       	       	       	       	  
		   				      		   
*/		       	       	       	       	       	       	   
