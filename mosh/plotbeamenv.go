package barf

import (
	"fmt"
	kass"barf/kass"
	//draw"barf/draw"
)

//PlotLp plots beam results (loads/bm/sf) for a single load pattern 
//1 - dead load case, 2 - live load case, 3 - shear force env, 4 - bm env
//if term = svg txtplot = filename; if term = dumb txtplot = txtplot lol
func PlotLp(lp int, bmenv map[int]*kass.BmEnv, ms map[int]*kass.Mem, mslmap map[int]map[int][][]float64, bmvec []int, term, title, folder string, plotchn chan []interface{}) {
	var data, xdata, sdata, ldata, dldata, lldata, sfld, bmld, foldr string
	dscale := 300.0
	var allplot string
	//gaah.GAAH.
	if folder == "web"{
		foldr = ""
	} else {
		foldr = folder
	}
	for i, idx := range bmvec{
		bm := bmenv[idx]
		rez := bm.EnvRez[lp]
		allplot += rez.Txtplot
		loads := mslmap[lp][idx]
		//fmt.Println("****loads****",loads)
		c1, c2 := bm.Coords[0], bm.Coords[1]
		x1, x2 := c1[0], c2[0]
		y1, y2 := 0.0, 0.0
		//if len(bm.Coords[0]) > 1{
		//	y1 = c1[1]; y2 = c2[1]
		//}
		lspan := ms[idx].Geoms[0]
		//log.Println("LSPAN",lspan)
		div := lspan/5.0
		//geom
		xdata += fmt.Sprintf("%f %f %f %f %v\n",(x1+x2)/2.0, -bm.Dims[1]/dscale,(x2-x1)/2.0,bm.Dims[1]/dscale,idx)
		//supports
		if i != len(bmvec) -1 {
			sdata += fmt.Sprintf("%f %f %f %f %v\n",x1, y1, 0.0, -2.0*bm.Dims[1]/dscale, 1)
		} else {
			sdata += fmt.Sprintf("%f %f %f %f %v\n",x1, y1, 0.0, -2.0*bm.Dims[1]/dscale, 1)
			sdata += fmt.Sprintf("%f %f %f %f %v\n",x2, y2, 0.0, -2.0*bm.Dims[1]/dscale, 1)
		}
		//section label
		//ldata +=  fmt.Sprintf("%f %f B%v[%.0fx%.0f]\n",(x1+x2)/2,-2.0*bm.Dims[1]/dscale, idx, bm.Dims[0],bm.Dims[1])
		ldata +=  fmt.Sprintf("%f %f B%v\n",(x1+x2)/2,-bm.Dims[1]/dscale/0.75, idx)
		//shear force and bending moment
		for j, x := range bm.Xs{
			data += fmt.Sprintf("%f %0.1f %0.3f\n",x1+x,rez.SF[j],-rez.BM[j])
		}
		data += "\n"
		//sf and bm labels and POINTS OF YEAY INTEREST
		//rez.Maxs = []float64{vmax, mmax, dmax, mpmax}
		//rez.Locs = []float64{vmaxx, mmaxx, dmaxx, mpmaxx}
		//rez.Cfxs = cfxs

		//sfld += fmt.Sprintf("%f %f (%.1f)\n",x1+rez.Locs[0],rez.Maxs[0],rez.Maxs[0])
		sfld += fmt.Sprintf("%f %f %.1f\n",x1+rez.Xs[0],rez.SF[0],rez.SF[0])
		sfld += fmt.Sprintf("%f %f %.1f\n",x1+rez.Xs[20],rez.SF[20],rez.SF[20])
		
		bmld += fmt.Sprintf("%f %f %.1f\n",x1+rez.Locs[1],-rez.Maxs[1],rez.Maxs[1])
		bmld += fmt.Sprintf("%f %f %.1f\n",x1+rez.Locs[3],-rez.Maxs[3],rez.Maxs[3])
		bmld += fmt.Sprintf("%f %f %.1f\n",x1+rez.Xs[0],-rez.BM[0],rez.BM[0])
		bmld += fmt.Sprintf("%f %f %.1f\n",x1+rez.Xs[20],-rez.BM[20],rez.BM[20])
		//member loads
		for _, vec := range loads{
			lt, wa, wb, la, lb, lc := vec[1], vec[2], vec[3], vec[4], vec[5], vec[6]
			switch lt{
				case 1.0:
				switch lc{
					case 1.0:
					dldata += fmt.Sprintf("%f %f %0.1f %0.1f %0.1f\n",x1+la,y1,0.0,wa,lt)
					case 2.0:
					lldata += fmt.Sprintf("%f %f %0.1f %0.1f %0.1f\n",x1+la,y1,0.0,wa,lt)
				}
				case 3.0:
				nd := int((lspan-la-lb)/div)
				x0 := x1 + la
				switch lc{
					case 1.0:
					for j := 0; j < nd; j++{
						dx := float64(j)*div
						dldata += fmt.Sprintf("%f %f %0.1f %0.1f %0.1f\n",x0+dx,y1,0.0,wa,lt)
					}
					case 2.0:
					for j := 0; j < nd; j++{
						dx := float64(j)*div
						lldata += fmt.Sprintf("%f %f %0.1f %0.1f %0.1f\n",x0+dx,y1,0.0,wa,lt)
					}
				}
				case 4.0:
				nd := int((lspan-la-lb)/div)
				x0 := x1 + la
				wd := (wb - wa)/(lspan - la - lb)
				switch lc{
					case 1.0:
					for j := 0; j < nd; j++{
						dx := float64(j)*div; dw := wd * dx
						dldata += fmt.Sprintf("%f %f %0.1f %0.1f %0.1f\n",x0+dx,y1,0.0,wa+dw,lt)
					}
					case 2.0:
					for j := 0; j < nd; j++{
						dx := float64(j)*div; dw := wd * dx
						lldata += fmt.Sprintf("%f %f %0.1f %0.1f %0.1f\n",x0+dx,y1,0.0,wa+dw,lt)
					}
				}
			}
		}
	}
	data += "\n\n"; xdata += "\n\n"; ldata += "\n\n"; sdata += "\n\n"; dldata += "\n\n";lldata += "\n\n"; sfld += "\n\n"; bmld += "\n\n"
	data += xdata; data += ldata; data += sdata; data += dldata; data += lldata; data += sfld; data += bmld
	fn := fmt.Sprintf("%s_bm_lp_%v.svg",title,lp)
	if term == "dxf"{
		fn = fmt.Sprintf("%s_bm_lp_%v.dxf",title,lp)
	}
	fname := genfname(foldr,fn)
	pltskript := "plotbmlp.gp"
	txtplot := skriptrun(data, pltskript, term, fn, fname)
	if term == "svg" || term == "svgmono" || term == "dxf"{
		txtplot = fname
	}	
	if folder == "web"{
		if term != "dxf"{kass.Svgkong(txtplot)}
		txtplot = fn
	}
	rez := make([]interface{},2)
	rez[0] = lp
	rez[1] = txtplot
	plotchn <-rez
}

//PlotBmEnv plots non redistributed beam envelopes
//for a continuous beam span
func PlotBmEnv(bmenv map[int]*kass.BmEnv, bmvec []int, term, title, folder string) (txtplot string){
	/*
	   plots non redistributed beam envelope 
	*/
	var data, xdata, sdata, ldata, vdata, mndata, mpdata, foldr string
	dscale := 50.0

	if folder == "web"{
		foldr = ""
	} else {
		foldr = folder
	}
	for _, idx := range bmvec{
		bm := bmenv[idx]
		c1, c2 := bm.Coords[0], bm.Coords[1]
		xdata += fmt.Sprintf("%f %f %f %f %v\n",(c2[0]+c1[0])/2.0,bm.Dims[1]/dscale,(c2[0]-c1[0])/2.0,bm.Dims[1]/dscale,idx)
		sdata += fmt.Sprintf("%f 0.0 %v\n",c1[0], 1)
		sdata += fmt.Sprintf("%f 0.0 %v\n",c2[0], 1)
		//section label
		ldata +=  fmt.Sprintf("%f %f B%v[%.0fx%.0f]\n",(c2[0]+c1[0])/2,bm.Dims[1]/dscale, idx, bm.Dims[0],bm.Dims[1])
		//max shear
		vdata += fmt.Sprintf("%f %f %.1f %v\n",c1[0]+bm.Vmaxx,bm.Vmax,bm.Vmax,idx)
		vdata += fmt.Sprintf("%f %f %.1f %v\n",c1[0]+bm.Xs[0],bm.Venv[0],bm.Venv[0],idx)
		vdata += fmt.Sprintf("%f %f %.1f %v\n",c1[0]+bm.Xs[20],bm.Venv[20],bm.Venv[20],idx)
		if bm.Lsx > 0.0 {vdata += fmt.Sprintf("%f %f (%.1f) %v\n",c1[0]+bm.Lsx,bm.Vl,bm.Vl,idx)}
		if bm.Rsx > 0.0 {vdata += fmt.Sprintf("%f %f (%.1f) %v\n",c2[0]-bm.Rsx,bm.Vr,bm.Vr,idx)}
		
		//max hogging bm
		mndata +=  fmt.Sprintf("%f %f %.1f %v\n",c1[0]+bm.Mnmaxx,-bm.Mnmax,bm.Mnmax,idx)
		mndata +=  fmt.Sprintf("%f %f %.1f %v\n",c1[0]+bm.Xs[0],-bm.Mnenv[0],bm.Mnenv[0],idx)
		mndata +=  fmt.Sprintf("%f %f %.1f %v\n",c1[0]+bm.Xs[20],-bm.Mnenv[20],bm.Mnenv[0],idx)
		//mndata += fmt.Sprintf("%f %f (%.1f) %v\n",c1[0]+bm.Lsx,-bm.Ml,-bm.Ml,idx)
		//mndata += fmt.Sprintf("%f %f (%.1f) %v\n",c2[0]-bm.Rsx,-bm.Mr,-bm.Mr,idx)
		//max sagging bm
		mpdata +=  fmt.Sprintf("%f %f %.1f %v\n",c1[0]+bm.Mpmaxx,-bm.Mpmax,bm.Mpmax,idx)
		//mpdata += fmt.Sprintf("%f %f (%.1f) %v\n",c1[0]+bm.Lsx,-bm.Ml,-bm.Ml,idx)
		//mpdata += fmt.Sprintf("%f %f (%.1f) %v\n",c2[0]-bm.Rsx,-bm.Mr,-bm.Mr,idx)
		for i, x := range bm.Xs{
			data += fmt.Sprintf("%f %0.1f %0.1f %0.1f\n",c1[0]+x,bm.Venv[i],-bm.Mnenv[i],-bm.Mpenv[i])
		}
	}
	data += "\n\n"; xdata += "\n\n"; ldata += "\n\n"; sdata += "\n\n"; mndata += "\n\n";mpdata += "\n\n";vdata += "\n\n"
	data += xdata; data += ldata; data += sdata; data += vdata; data += mndata; data += mpdata
	//fmt.Println("***data***\n",data)
	//txtplot = dumbplt(data,"plotbmenv.gp",term)
	fn := fmt.Sprintf("%s_elastic_env.svg",title)
	if term == "dxf"{
		fn = fmt.Sprintf("%s_elastic_env.dxf",title)
	}
	
	fname := genfname(foldr,fn)
	pltskript := "plotbmenv.gp"
	txtplot = skriptrun(data, pltskript, term, fn, fname)
	if term == "svg" || term == "svgmono" || term == "dxf"{
		txtplot = fname
	}
	if folder == "web"{
		if term != "dxf"{kass.Svgkong(txtplot)}
		txtplot = fn
	}
	return
}

//PlotBmRdEnv plots redistributed bm and shear envelopes
//for a continuous beam span 
func PlotBmRdEnv(bmenv map[int]*kass.BmEnv, bmvec []int, term, title, folder string) (txtplot string){
	/*
	   plots redistributed bm and shear envelopes 
	*/
	var data, xdata, sdata, ldata, vdata, mndata, mpdata, foldr string
	dscale := 500.0
	if folder == "web"{
		foldr = ""
	} else {
		foldr = folder
	}
	for _, idx := range bmvec{
		bm := bmenv[idx]
		c1, c2 := bm.Coords[0], bm.Coords[1]
		//log.Println("coords->",c1,c2)
		xdata += fmt.Sprintf("%f %f %f %f %v\n",(c2[0]+c1[0])/2.0,bm.Dims[1]/dscale,(c2[0]-c1[0])/2.0,bm.Dims[1]/dscale,idx)
		sdata += fmt.Sprintf("%f 0.0 %v\n",c1[0], 1)
		sdata += fmt.Sprintf("%f 0.0 %v\n",c2[0], 1)
		//section label
		ldata +=  fmt.Sprintf("%f %f B%v[%.0fx%.0f]\n",(c2[0]+c1[0])/2,bm.Dims[1]/dscale, idx, bm.Dims[0],bm.Dims[1])
		//max shear
		vdata += fmt.Sprintf("%f %f %.1f %v\n",c1[0]+bm.Vrmaxx,bm.Vrmax,bm.Vrmax,idx)
		//vdata += fmt.Sprintf("%f %f (%.1f) %v\n",c1[0]+bm.Lsx,bm.Vlrd,bm.Vlrd,idx)
		//vdata += fmt.Sprintf("%f %f (%.1f) %v\n",c2[0]-bm.Rsx,bm.Vrrd,bm.Vrrd,idx)
		
		//max hogging bm
		mndata +=  fmt.Sprintf("%f %f %.1f %v\n",c1[0]+bm.Mnrmaxx,-bm.Mnrmax,bm.Mnrmax,idx)
		//mndata += fmt.Sprintf("%f %f (%.1f) %v\n",c1[0]+bm.Lsx,-bm.Ml,-bm.Ml,idx)
		//mndata += fmt.Sprintf("%f %f (%.1f) %v\n",c2[0]-bm.Rsx,-bm.Mr,-bm.Mr,idx)
		//max sagging bm
		mpdata +=  fmt.Sprintf("%f %f %.1f %v\n",c1[0]+bm.Mprmaxx,-bm.Mprmax,bm.Mprmax,idx)
		//mpdata += fmt.Sprintf("%f %f (%.1f) %v\n",c1[0]+bm.Lsx,-bm.Ml,-bm.Ml,idx)
		//mpdata += fmt.Sprintf("%f %f (%.1f) %v\n",c2[0]-bm.Rsx,-bm.Mr,-bm.Mr,idx)
		for i, x := range bm.Xs{
			data += fmt.Sprintf("%f %f %f %f %f %f %f %f %f\n",c1[0]+x,bm.Vrd[i],-bm.Mnrd[i],-bm.Mprd[i], bm.Venv[i],-bm.Mnenv[i],-bm.Mpenv[i],-0.7*bm.Mnenv[i],-0.7*bm.Mpenv[i])
		}
		data += "\n"
	}
	data += "\n\n"; xdata += "\n\n"; ldata += "\n\n"; sdata += "\n\n"; mndata += "\n\n";mpdata += "\n\n";vdata += "\n\n"
	data += xdata; data += ldata; data += sdata; data += vdata; data += mndata; data += mpdata
	//fmt.Println("***data***\n",data)
	//txtplot = dumbplt(data,"plotbmrdenv.gp",term)
	fn := fmt.Sprintf("%s_redist_env.svg",title)
	if term == "dxf"{
		fn = fmt.Sprintf("%s_redist_env.dxf",title)
	}
	fname := genfname(foldr,fn)
	pltskript := "plotbmrdenv.gp"
	txtplot = skriptrun(data, pltskript, term, fn, fname)
	
	if term == "svg" || term == "svgmono" || term == "dxf"{
		txtplot = fname
	}
	if folder == "web"{
		if term != "dxf"{kass.Svgkong(txtplot)}
		txtplot = fn
	}
	return
}
