package barf		  
			  
import (		  
	"math"		  
	"log"		  
)

//GenPloads generates nodal loads for ldcalc = -1 (all specified loads)
func (t *Trs2d) GenPloads()(err error){
	xmid := t.Span/2.0
	if t.Typ == 1 || t.Typ == 3{xmid = t.Span}
	var theta float64
	if t.Typ != 5 && t.Typ != 6{		
		theta = math.Atan(1.0/t.Slope)	
	}
	//log.Println(ColorCyan,"theta in rad-", theta, "in degrees-",theta*180.0/math.Pi,ColorReset)
	t.Mod.Jldcs = make(map[float64][][]float64)
	var nwl int
	switch t.Nlds{
		case 0, 1:
		//just pdl
		for _, node := range t.Tcns{
			ldcase := []float64{float64(node), 0.0, -t.Pdl}
			t.Mod.Jldcs[1] = append(t.Mod.Jldcs[1], ldcase)
		}
		for _, node := range []int{t.Bcns[0],t.Bcns[len(t.Bcns)-1]}{
			ldcase := []float64{float64(node), 0.0, -t.Pdl/2.0}
			t.Mod.Jldcs[1] = append(t.Mod.Jldcs[1], ldcase)
		}
		case 2:
		
		//pdl, pll
		for _, node := range t.Tcns{
			ldcase := []float64{float64(node), 0.0, -t.Pdl}
			t.Mod.Jldcs[1] = append(t.Mod.Jldcs[1], ldcase)
			ldcase = []float64{float64(node), 0.0, -t.Pll}
			t.Mod.Jldcs[2] = append(t.Mod.Jldcs[2], ldcase)
		}
		
		for _, node := range []int{t.Bcns[0],t.Bcns[len(t.Bcns)-1]}{
			ldcase := []float64{float64(node), 0.0, -t.Pdl/2.0}
			t.Mod.Jldcs[1] = append(t.Mod.Jldcs[1], ldcase)
			ldcase = []float64{float64(node), 0.0, -t.Pll/2.0}
			t.Mod.Jldcs[2] = append(t.Mod.Jldcs[2], ldcase)
		}
		case 3:
		//pdl, pll, pwl, pwr, pwp
		//theta := math.Atan(1.0/t.Slope)
		nwl = 1
		if t.Pwp > 0.0{nwl = 2}
		ldnodes := append(t.Tcns, []int{t.Bcns[0],t.Bcns[len(t.Bcns)-1]}...)
		for _, node := range ldnodes{
			ldcase := []float64{float64(node), 0.0, -t.Pdl}
			if node == t.Bcns[0] || node == t.Bcns[len(t.Bcns)-1]{
				ldcase[2] = ldcase[2]/2.0
			}
			t.Mod.Jldcs[1] = append(t.Mod.Jldcs[1], ldcase)
			
			ldcase = []float64{float64(node), 0.0, -t.Pll}
			if node == t.Bcns[0] || node == t.Bcns[len(t.Bcns)-1]{
				ldcase[2] = ldcase[2]/2.0
			}
			t.Mod.Jldcs[2] = append(t.Mod.Jldcs[2], ldcase)
			
			switch t.Typ{
				case 5, 6:
				//par truss, fun truss
				default:
				xn := t.Mod.Coords[node-1][0]
				//log.Println("xnode - ",xn, "xmid -",xmid)
				switch{
					case xn == xmid:
					//log.Println(ColorYellow,"APEX,",node,ColorReset)
					fx1 := -t.Pwl * math.Sin(theta)/2.0
					fy1 := t.Pwl * math.Cos(theta)/2.0
					ldcase = []float64{float64(node), fx1, fy1}
					t.Mod.Jldcs[3.1] = append(t.Mod.Jldcs[3.1], ldcase)

					fx1 = t.Pwr * math.Sin(theta)/2.0
					fy1 = t.Pwr * math.Cos(theta)/2.0
					ldcase = []float64{float64(node), fx1, fy1}
					t.Mod.Jldcs[3.1] = append(t.Mod.Jldcs[3.1], ldcase)

					if t.Pwp > 0.0{							
						fx1 = t.Pwp * math.Sin(theta)/2.0
						fy1 = -t.Pwp * math.Cos(theta)/2.0
						ldcase = []float64{float64(node), fx1, fy1}
						t.Mod.Jldcs[3.2] = append(t.Mod.Jldcs[3.2], ldcase)

						fx1 = -t.Pwp * math.Sin(theta)/2.0
						fy1 = -t.Pwp * math.Cos(theta)/2.0
						ldcase = []float64{float64(node), fx1, fy1}
						t.Mod.Jldcs[3.2] = append(t.Mod.Jldcs[3.2], ldcase)
					}
					case xn < xmid:
					fx1 := -t.Pwl * math.Sin(theta)
					fy1 := t.Pwl * math.Cos(theta)
					
					ldcase = []float64{float64(node), fx1, fy1}
					if node == t.Bcns[0]{
						ldcase[1] = ldcase[1]/2.0
						ldcase[2] = ldcase[2]/2.0
					}
					t.Mod.Jldcs[3.1] = append(t.Mod.Jldcs[3.1], ldcase)

					
					if t.Pwp > 0.0{
						fx1 = t.Pwp * math.Sin(theta)
						fy1 = -t.Pwp * math.Cos(theta)
						ldcase = []float64{float64(node), fx1, fy1}
						if node == t.Bcns[0]{
							ldcase[1] = ldcase[1]/2.0
							ldcase[2] = ldcase[2]/2.0
						}
						t.Mod.Jldcs[3.2] = append(t.Mod.Jldcs[3.2], ldcase)
						
					}
					case xn > xmid:
					fx1 := t.Pwr * math.Sin(theta)
					fy1 := t.Pwr * math.Cos(theta)
					ldcase = []float64{float64(node), fx1, fy1}
					if node == t.Bcns[len(t.Bcns)-1]{
						ldcase[1] = ldcase[1]/2.0
						ldcase[2] = ldcase[2]/2.0
					}
					t.Mod.Jldcs[3.1] = append(t.Mod.Jldcs[3.1], ldcase)
					if t.Pwp > 0.0{						
						fx1 = -t.Pwp * math.Sin(theta)
						fy1 = -t.Pwp * math.Cos(theta)
						ldcase = []float64{float64(node), fx1, fy1}
						if node == t.Bcns[len(t.Bcns)-1]{
							ldcase[1] = ldcase[1]/2.0
							ldcase[2] = ldcase[2]/2.0
						}
						t.Mod.Jldcs[3.2] = append(t.Mod.Jldcs[3.2], ldcase)	
					}
				}
			}
		}
	}
	t.Mod.Nwlc = nwl
	t.Mod.Nlds = t.Nlds
	return
}

//GenLd generates loads for a 2d Trs gen struct
func (t *Trs2d) GenLd() (err error){
	//lrftr, lspan, spacing, dl, ll float64, slfwt, mtyp, ftyp, rfmat,trstyp int) (wdl, wll float64){
	//returns total dead load w for a 2d truss
	//if slope == 0 flat roof
	//npr - number of purlins
	//slfwt of wooden truss/steel truss
	//log.Println("CALCing")
	var wdl, wll, llred, wtruss float64 
	al := t.Rftrl * t.Spacing
	if t.Typ == 2 || t.Typ == 4{al = 2.0 * al}
	if t.Typ == 5 {t.Slope = 0}
	if t.Roofmat == 0{t.Roofmat = 1}
	// xmid := t.Span/2.0
	// if t.Typ == 1 || t.Typ == 3{xmid = t.Span}
	var theta float64
	if t.Typ != 5 && t.Typ != 6{		
		theta = math.Atan(1.0/t.Slope) * 180.0/math.Pi
		//if t.Spam{log.Println(ColorRed,"ROOFER ANGLE->",theta,"DEGREES",ColorReset)}
		if theta > 10.0{
			llred = (theta - 10.0) * 0.02 * 1e3
		}
		if t.Tang == 0.0{
			t.Tang = theta
		}
	}
	theta = math.Atan(1.0/t.Slope)
	switch t.Ldcalc{
		case -1:
		//all loads are given etc
		err = t.GenPloads()
		return
		// default:
		// if t.PSFs == nil{
		// 	t.PSFs = PsfsIs[t.Ldcalc]
		// }
	}
	
	
	if t.DL > 0.0{wdl += t.DL} else {wdl += rfWt[t.Roofmat-1][0] * 1e3}
	if t.LL > 0.0 {
		wll += t.LL - llred
	} else {
		wll += 750.0  - llred
		if wll < 400.0{wll = 400.0}
	}
	switch t.Mtyp{
		case 3:
		//timber truss
		switch t.Clctyp{
			case 2:
			//abel 160 n/m2
			wdl += 160.0
			default:
			//john wight formula for wood trusses
			wtruss = 0.75 * t.Span/304.8 * t.Spacing/304.8 * (1.0 + t.Span/304.8/10.0) * 0.454 * 9.81
			log.Println("in truss ld, tmbr truss weight",wtruss)
		}
		case 2:
		//steel truss
		switch t.Clctyp{
			case 0:
			//welded steel truss/bhavikatti formula
			wdl += 20.0 + 6.6 * t.Span/1000.0
			case 1:
			//general duggal formula
			wdl += (t.Span/1000.0 + 5.0) * 10.0
			case 2:
			//bolted truss
			wtruss = (53.7 + 0.53 * t.Span/1e3 * t.Spacing/1e3)
			case 3:
			//ketchum formula
		}
		
	}
	
	if t.Planld{
		//dl and ll on plan area
		wdl = wdl * math.Cos(theta)
		wll = wll * math.Cos(theta)
	}
	var pd float64	  
	var cpos, cneg []float64
	var wlcs map[int][]float64
	if t.Cpi == 0{t.Cpi = 0.2}
	var slope float64 
	if t.Typ == 5 || t.Typ == 6{slope =0} else {slope = 1.0/t.Slope}
	minw := t.Span
	lbld := t.Width
	if minw > t.Width{minw = t.Width; lbld = t.Span}
	switch t.Typ{
		case 1, 3:
		pd, cpos, cneg, wlcs = wltable7(t.Vb, t.Height, minw, lbld, slope, t.Cpi)
		case 2, 4:
		pd, cpos, cneg, wlcs = wltable6(t.Vb, t.Height, minw, slope, t.Cpi)
	}
	if t.Spam{
		log.Println("wdl",wdl,"wtruss",wtruss)
		log.Println("wll",wll)
		log.Println("truss tributary area->",al*1e-6,"m2")
		log.Println("nodal tributary area->",t.Purlinspc*t.Spacing*1e-6,"m2")
		log.Println("total load per truss->",(wdl+wll)*al*1e-6+wtruss,"n")
		log.Println("total nodes->",len(t.Tcns))

		log.Println("pd",pd)	  
		log.Println("cpos",cpos) 	    
		log.Println("cneg",cneg)
		log.Println("printing cpos and wload cases-")
		for i := range wlcs[1]{
			cpvec := []float64{wlcs[1][i]/pd,wlcs[2][i]/pd}
			t.Cps = append(t.Cps, cpvec)
			if t.Spam{		
				fxl := wlcs[1][i] * math.Sin(theta)*t.Purlinspc * t.Spacing * 1e-6; fyl := wlcs[1][i] * math.Cos(theta)*t.Purlinspc * t.Spacing * 1e-6
				fxr := -wlcs[2][i] * math.Sin(theta) * t.Purlinspc * t.Spacing * 1e-6; fyr := wlcs[2][i] * math.Cos(theta)*t.Purlinspc * t.Spacing * 1e-6
				
				log.Println(ColorRed,"cpe left->",wlcs[1][i]/pd,"cpe right->",wlcs[2][i]/pd, ColorReset)
				log.Println("left->",wlcs[1][i],"N/m2","right->",wlcs[2][i],"N/m2")
				log.Println(ColorCyan, "pwl\n","left",wlcs[1][i] * t.Purlinspc * t.Spacing * 1e-6,"right",wlcs[2][i] *t.Purlinspc * t.Spacing * 1e-6,ColorReset)
				log.Println("left nodal forces fxl", fxl, "fyl",fyl, " right fxr",fxr, "fyr",fyr)
				
			}
		}
	}
	dltot := (wdl)*al*1e-6+wtruss
	lltot := (wll)*al*1e-6
	ntot := (float64(len(t.Tcns))) * 2.0 + 2.0
	switch t.Typ{
		case 1:
		case 2:
		case 3:
		case 4:
	}
	//wtot := dltot + lltot
	if t.Spam{
		// log.Println("load per node-> int->",wtot*2.0/ntot,"end node->",wtot/ntot)
		// log.Println("rftrl->",t.Rftrl)
		// //add tc node load cases - dead, live
		// //wlcases  - 1 left, 2 right
		// //log.Println(ColorRed)
		// log.Println("TCNODES->",t.Tcns)
		// //log.Println(ColorCyan)
		// log.Println("BCNODES->",t.Bcns)
	}
	
	//log.Println(ColorReset)
	var nwl, nsl int
	t.Pdl = dltot*2.0/ntot + t.Purlinwt * t.Spacing/1e3
	t.Pll = lltot*2.0/ntot
	for i, val := range wlcs[1]{
		if val > 0.0{
			if t.Pwp == 0.0{
				t.Pwp = val 
			} else if t.Pwp < val{
				t.Pwp = val
			}
		}
		if t.Pwl == 0.0{
			t.Pwl = val
			t.Pwr = wlcs[2][i] 
		} else if math.Abs(t.Pwl) < val{
			t.Pwl = val
			t.Pwr = wlcs[2][i]
		}
	}
	t.Pwl = t.Pwl * t.Purlinspc * t.Spacing * 1e-6
	t.Pwr = t.Pwr * t.Purlinspc * t.Spacing * 1e-6
	t.Pwp = t.Pwp * t.Purlinspc * t.Spacing * 1e-6
	for _, node := range t.Tcns{
		//dead load
		dl := dltot*2.0/ntot; ll := lltot*2.0/ntot
		dl += t.Purlinwt * t.Spacing/1e3
		//log.Println(ColorCyan, "dl, purlin dl",dl, t.Purlinwt*t.Spacing/1e3,ColorReset)
		
		//log.Println(ColorGreen, "ll",ll,ColorReset)
		t.Jldsrv[1] = append(t.Jldsrv[1],[]float64{float64(node), 0.0, -dl})
		t.Jldsrv[2] = append(t.Jldsrv[2],[]float64{float64(node), 0.0, -ll})
		//wind load
		switch t.Typ{
			case 1, 3:
			//l type - 
			case 2, 4:
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
	// log.Println("jldsrv before end nodes")
	// for lc, val := range t.Jldsrv{
	// 	log.Println("lc, val-",lc, val)
	// }
	//end nodes
	for _, node := range []int{t.Bcns[0],t.Bcns[len(t.Bcns)-1]}{
		dl := dltot/ntot; ll := lltot/ntot
		dl += t.Purlinwt * t.Spacing/1e3
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
	if t.Spam{
		// for i, ldcases := range t.Jldsrv{
		// 	for _, ld := range ldcases{
		// 		log.Println("ld typ",i,"node",ld[0],"frcx",ld[1],"frcy",ld[2])
		// 	}
		// }
	}
	switch t.Spose{
		case false:
		mldcs, jldcs := GenLdCombos(nwl, nsl, t.Jldsrv, t.Mldsrv, t.PSFs)
		t.Mod.Mldcs = mldcs
		t.Mod.Jldcs = jldcs
		t.Mod.Nwlc = nwl
		case true:
		t.Mod.Jldcs = t.Jldsrv
		t.Mod.Nwlc = nwl
	}
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
/*

	if t.Ldcalc == 0 && t.PSFs == nil{
		switch t.Code{
			case 1:	
			t.PSFs = [][]float64{
				{1.5,1.5,0.0},
				{0.9,0.0,1.5},
			}
		}
	}

*/
/*


	if t.Ldcalc == -1{
		// t.Mod.Jldcs = make(map[float64][][]float64)
		// switch t.Nlds{
		// 	case 0, 1:
		// 	//just pdl
		// 	for _, node := range t.Tcns{
		// 		ldcase := []float64{float64(node), 0.0, -t.Pdl}
		// 		t.Mod.Jldcs[1] = append(t.Mod.Jldcs[1], ldcase)
		// 	}
		// 	case 2:
		// 	//pdl, pll
		// 	for _, node := range t.Tcns{
		// 		ldcase := []float64{float64(node), 0.0, -t.Pdl}
		// 		t.Mod.Jldcs[1] = append(t.Mod.Jldcs[1], ldcase)
		// 		ldcase = []float64{float64(node), 0.0, -t.Pll}
		// 		t.Mod.Jldcs[2] = append(t.Mod.Jldcs[2], ldcase)
		// 	}
		// 	case 3:
		// 	//pdl, pll, pwl, pwr, pwp
		// 	//theta := math.Atan(1.0/t.Slope)
		// 	for _, node := range t.Tcns{
		// 		ldcase := []float64{float64(node), 0.0, -t.Pdl}
		// 		t.Mod.Jldcs[1] = append(t.Mod.Jldcs[1], ldcase)
		// 		ldcase = []float64{float64(node), 0.0, -t.Pll}
		// 		t.Mod.Jldcs[2] = append(t.Mod.Jldcs[2], ldcase)
				
		// 		switch t.Typ{
		// 			case 5, 6:
		// 			//par truss, fun truss
		// 			default:
		// 			xn := t.Mod.Coords[node-1][0]
		// 			log.Println("xnode - ",xn, "xmid -",xmid)
		// 			switch{
		// 				case xn == xmid:
		// 				log.Println(ColorYellow,"APEX,",node,ColorReset)
		// 				fx1 := -t.Pwl * math.Sin(theta)/2.0
		// 				fy1 := t.Pwl * math.Cos(theta)/2.0
		// 				ldcase = []float64{float64(node), fx1, fy1}
		// 				t.Mod.Jldcs[3] = append(t.Mod.Jldcs[3], ldcase)
		// 				fx1 = t.Pwr * math.Sin(theta)/2.0
		// 				fy1 = t.Pwr * math.Cos(theta)/2.0
		// 				ldcase = []float64{float64(node), fx1, fy1}
		// 				t.Mod.Jldcs[3] = append(t.Mod.Jldcs[3], ldcase)

		// 				fx1 = t.Pwp * math.Sin(theta)/2.0
		// 				fy1 = -t.Pwp * math.Cos(theta)/2.0
		// 				ldcase = []float64{float64(node), fx1, fy1}
		// 				t.Mod.Jldcs[4] = append(t.Mod.Jldcs[4], ldcase)

						
		// 				fx1 = -t.Pwp * math.Sin(theta)/2.0
		// 				fy1 = -t.Pwp * math.Cos(theta)/2.0
		// 				ldcase = []float64{float64(node), fx1, fy1}
		// 				t.Mod.Jldcs[4] = append(t.Mod.Jldcs[4], ldcase)

		// 				case xn < xmid:
		// 				fx1 := -t.Pwl * math.Sin(theta)
		// 				fy1 := t.Pwl * math.Cos(theta)
						
		// 				ldcase = []float64{float64(node), fx1, fy1}
		// 				t.Mod.Jldcs[3] = append(t.Mod.Jldcs[3], ldcase)

						
		// 				fx1 = t.Pwp * math.Sin(theta)
		// 				fy1 = -t.Pwp * math.Cos(theta)
		// 				ldcase = []float64{float64(node), fx1, fy1}
		// 				t.Mod.Jldcs[4] = append(t.Mod.Jldcs[4], ldcase)

		// 				case xn > xmid:
		// 				fx1 := t.Pwr * math.Sin(theta)
		// 				fy1 := t.Pwr * math.Cos(theta)
		// 				ldcase = []float64{float64(node), fx1, fy1}
		// 				t.Mod.Jldcs[3] = append(t.Mod.Jldcs[3], ldcase)

						
		// 				fx1 = -t.Pwp * math.Sin(theta)
		// 				fy1 = -t.Pwp * math.Cos(theta)
		// 				ldcase = []float64{float64(node), fx1, fy1}
		// 				t.Mod.Jldcs[4] = append(t.Mod.Jldcs[4], ldcase)

		// 			}
		// 		}

		// 	}
		// }
		// if t.Spam{
		// 	log.Println("load case 3")
		// 	for i, ldcase := range t.Mod.Jldcs[3]{
		// 		log.Println("node i - ",i, " ldcase ", ldcase)
		// 	}
			
		// 	log.Println("load case 4")
		// 	for i, ldcase := range t.Mod.Jldcs[3]{
		// 		log.Println("node i - ",i, " ldcase ", ldcase)
		// 	}
		// }
		// return
	}
*/
