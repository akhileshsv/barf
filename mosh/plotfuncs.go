package barf

import (
	"fmt"
	"log"
	"math"
	//"math/rand"
	"bytes"
	"os"
	"os/exec"	
	"runtime"
	"path/filepath"
	kass"barf/kass"
)

//skriptpath returns the absolute path of the (gnuplot) script named skript in current folder
func skriptpath(skript string) (string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, skript)
}

//genfname returns the absolute path of the svg/dxf output file for gnuplot
func genfname(folder, fn string) (string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	if folder != ""{
		return filepath.Join(folder,fn)
	}
	return filepath.Join(basepath, "../data/out",fn)
}

//skriptrun runs a gnuplot script and returns the text plot/path to svg
func skriptrun(data, pltskript, term, title, fname string) (string){
	pltskript = skriptpath(pltskript)
	f, e1 := os.CreateTemp("", "mosh")
	if e1 != nil {
		fmt.Println(e1)
	}
	defer f.Close()
	defer os.Remove(f.Name())	
	_, e1 = f.WriteString(data)
	if e1 != nil {
		fmt.Println(e1)
	}
	cmd := exec.Command("gnuplot","-c",pltskript,f.Name(),term, title, fname)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outstr, errstr := stdout.String(), stderr.String()
	outstr = fmt.Sprintf("%s",err) + errstr + outstr
	// if err != nil {
	// 	//fmt.Println(err)
	// }
	// if errstr != "" {
	// 	//fmt.Println(errstr)
	// }
	return outstr
}


//planrect returns a rect as headless vectors for gnuplot
func planrect(idx int, xs, ys, lx, ly, bsupx, bsupy float64, plotdim bool)(data, adata, ldata string){
	xb := xs - bsupx/2.0; yb := ys - bsupy/2.0
	lx1 := lx + bsupx; ly1 := ly + bsupy
	pts := [][]float64{
		{xb, yb, lx1, 0.0, float64(idx), 0.0},
		{xb, yb, 0.0, ly1, float64(idx), 0.0},
		{xb+lx1, yb+ly1, -lx1, 0.0, float64(idx), 0.0},
		{xb+lx1, yb+ly1, 0.0, -ly1, float64(idx), 0.0},
	}
	pts1 := [][]float64{
		{xs, ys, lx, 0.0, float64(idx), 3.0},
		{xs, ys, 0.0, ly, float64(idx), 3.0},
		{xs+lx, ys+ly, -lx, 0.0, float64(idx), 3.0},
		{xs+lx, ys+ly, 0.0, -ly, float64(idx), 3.0},
	}

	for i, pt := range pts{
		pt1 := pts1[i]
		data += fmt.Sprintf("%f %f %f %f %f %f %f %f %f %f\n",pt[0],pt[1],pt[2],pt[3],pt[4],pt1[0],pt1[1],pt1[2],pt1[3],pt1[4])
	}
	
	if plotdim{
		adata += fmt.Sprintf("%f %f %f %f %f\n",xs,ys-0.5*ly,lx,0.0,float64(idx))
		ldata += fmt.Sprintf("%f %f %.2f\n",xs+lx/2.0,ys-0.4*ly,lx)
		adata += fmt.Sprintf("%f %f %f %f %f\n",xs-0.5*lx,ys,0.0,ly,float64(idx))
		ldata += fmt.Sprintf("%f %f %.2f\n",xs-0.4*lx,ys+ly/2.0,ly)
	}
	return
}

//slabplanbars returns bottom and top bars for gnuplot
// func slabplanbars(s *RccSlb)(data, adata, ldata string){
// 	switch typ{
// 		case 1:
// 		//un way
// 		x0 := xs; y0 := ys
// 		nx := int(math.Round(ly/spcx))
// 		spcx1 := ly/float64(nx)
// 		dx := lx
// 		for i := 0; i <= nx; i++{
// 			x0 = xs
// 			dx = lx
// 			if i % 2 != 0{
// 				dx = lx - 0.2 * lx
// 				x0 = xs + 0.1 * lx
// 			}
// 			data += fmt.Sprintf("%f %f %f %f 1.0\n", x0, y0, dx, 0.0)
// 			y0 += spcx1
// 		}
// 		x0 = xs; y0 = ys
// 		ny := int(math.Round(lx/spcy)) 
// 		spcy1 := lx/float64(ny)
// 		for i := 0; i <= ny; i++{
// 			data += fmt.Sprintf("%f %f %f %f 2.0\n", x0, y0, 0.0, ly)
// 			x0 += spcy1
// 		}
// 		if plotdim{
// 			adata += fmt.Sprintf("%f %f %f %f %f\n", xs - 0.2 * lx, ys, 0.0, spcx, 1.0)
// 			ldata += fmt.Sprintf("%f %f X-T%.f-%.f\n", xs - 0.1 * lx, ys, diax, spcx)
// 			adata += fmt.Sprintf("%f %f %f %f %f\n", xs, ys - 0.2 * ly, spcy, 0.0, 1.0)
// 			ldata += fmt.Sprintf("%f %f Y-T%.f-%.f\n", xs, ys - 0.1 * ly, diay, spcy)
// 		}
// 		case 2:
		
// 	}
// 	return	
// }

//DrawPlan returns the plan data for a slab (top and bottom)
// func (s *RccSlb) DrawPlan() (pltstr string){
// 	if s.Id == 0{
// 		s.Id = rand.Intn(666)
// 	}
// 	var data, stdata, ldata, adata string
// 	var xs, ys float64 
// 	switch s.Type{
// 		case 1:
// 		//one way slab
// 		switch s.Endc{
// 			case 0,1:
// 			//ss slab/clvr
// 			if s.Lx == 0.0{s.Lx = s.Lspan}
// 			if s.Ly == 0.0{s.Ly = 2.0 * s.Lx}
// 			//if s.Bsup == 0.0{s.Bsup = 200.0}
// 			d1, a1, l1 := planrect(s.Id, xs, ys, s.Lx, s.Ly, s.Bsup, 0.0, true) 
// 			data += d1; adata += a1; ldata += l1
// 			//bottom mesh
// 			d1, a1, l1 = slabplanbars(s.Id, s.Type, s.Endc, xs + s.Bsup/2.0, ys + s.Bsup/2.0, s.Lx - 2.0*s.Efcvr - s.Bsup, s.Ly - 2.0*s.Efcvr - s.Bsup, s.Diamain, s.Diadist, s.Spcm, s.Spcd, true) 
// 			stdata += d1
// 			adata += a1; ldata += l1
// 			adata += "\n\n"; ldata += "\n\n"; data += "\n\n"; stdata += "\n\n"
// 			data += stdata; data += adata; data += ldata
// 			//top left mesh
// 			d1, a1, l1 = slabplanbars(s.Id, s.Type, s.Endc, xs + s.Efcvr, ys + s.Efcvr, 0.1 * s.Lx + s.Bsup/2.0, s.Ly- 2.0*s.Efcvr, s.Diamain, s.Diadist, s.Spcm, s.Spcd, false)
// 			data += d1; adata = a1; ldata = l1 
// 			//top right mesh
// 			d1, a1, l1 = slabplanbars(s.Id, s.Type, s.Endc, xs + 0.9 * s.Lx - s.Bsup/2.0, ys + 25.0, 0.1 * s.Lx + s.Bsup/2.0, s.Ly, s.Diamain, s.Diadist, s.Spcm, s.Spcd, false) 
// 			data += d1; data += "\n\n"; adata += "\n\n"; ldata += "\n\n"
// 			data += adata; data += ldata 
// 			default:
// 		} 
// 		case 2:
// 		//two way slab
// 		if s.Bsup == 0.0{s.Bsup = 200.0}
// 		d1, a1, l1 := planrect(s.Id, xs, ys, s.Lx, s.Ly, s.Bsup, 0.0, true) 
// 		data += d1; adata += a1; ldata += l1
// 		//bottom mesh
// 		d1, a1, l1 = slabplanbars(s.Id, s.Type, s.Endc, xs + s.Bsup/2.0, ys + s.Bsup/2.0, s.Lx - 2.0*s.Efcvr - s.Bsup, s.Ly - 2.0*s.Efcvr - s.Bsup, s.Diamain, s.Diadist, s.Spcm, s.Spcd, true) 
// 		stdata += d1
// 		adata += a1; ldata += l1
// 		adata += "\n\n"; ldata += "\n\n"; data += "\n\n"; stdata += "\n\n"
// 		data += stdata; data += adata; data += ldata
// 		//top left mesh
// 		d1, a1, l1 = slabplanbars(s.Id, s.Type, s.Endc, xs + s.Efcvr, ys + s.Efcvr, 0.1 * s.Lx + s.Bsup/2.0, s.Ly- 2.0*s.Efcvr, s.Diamain, s.Diadist, s.Spcm, s.Spcd, false)
// 		data += d1; adata = a1; ldata = l1 
// 		//top right mesh
// 		d1, a1, l1 = slabplanbars(s.Id, s.Type, s.Endc, xs + 0.9 * s.Lx - s.Bsup/2.0, ys + 25.0, 0.1 * s.Lx + s.Bsup/2.0, s.Ly, s.Diamain, s.Diadist, s.Spcm, s.Spcd, false) 
// 		data += d1; data += "\n\n"; adata += "\n\n"; ldata += "\n\n"
// 		data += adata; data += ldata 

// 		case 3:
// 		case 4:
// 	}
// 	//data += s.Draw()
// 	//fmt.Println("BEHOLD DATAPHYLE ->\n",data)
// 	if s.Title == ""{
// 		s.Title = fmt.Sprintf("rcc slab plan %v", s.Id)
// 		//fname = fmt.Sprintf("col%v-%v.svg",c.Mid, c.Id)
// 	}
// 	title := s.Title + ".svg"
// 	if s.Term == "dxf"{
// 		title = s.Title + ".dxf"
// 	}
// 	fname := genfname(s.Foldr,title)
// 	pltstr = skriptrun(data, "plotslabplan.gp", s.Term, s.Title, fname)
// 	if s.Term == "svg" || s.Term == "dxf"{
// 		pltstr = fname
// 	}
// 	return
// }

//PlotColGeom plots an rcc column section 
func PlotColGeom(c *RccCol, term string, multi bool) (pltstr string){
	var fname, data string
	if c.Data == ""{data = c.Draw()} else {data = c.Data}
	if c.Title == ""{
		c.Title = fmt.Sprintf("rcc column section %v-%v", c.Mid, c.Id)
		//fname = fmt.Sprintf("col%v-%v.svg",c.Mid, c.Id)
	}
	title := c.Title + ".svg"
	if term == "dxf"{
		title = c.Title + ".dxf"
	}
	fname = genfname(c.Foldr,title)
	pltstr = skriptrun(data, "plotcolgeom.gp", term, c.Title, fname)
	if term == "svg" || term == "svgmono" || term == "dxf"{
		pltstr = fname
	}
	if c.Web{
		if term != "dxf"{kass.Svgkong(pltstr)}
		pltstr = title
	}
	//c.Txtplot = append(c.Txtplot, pltstr)
	return
}

//PlotCBmDet plots the steel detailing drawing of a cbeam
func PlotCBmDet(web bool, bmvec []int, bms [][]*RccBm, folder, title, term string) (pltstr string){
	var xs float64
	var gdata, sdata, stdata, ldata, mdata, data string
	for _, id := range bmvec{
		bmarr := bms[id-1]
		g1, s1, st1, l1, m1 := BmSpanDraw(xs, bmarr)
		lspan := bmarr[1].Lspan
		xs += lspan * 1e3
		gdata += g1; sdata += s1; stdata += st1; ldata += l1; mdata += m1
	}
	gdata += "\n\n"; sdata += "\n\n"; stdata += "\n\n"
	ldata += "\n\n"
	data += gdata; data += sdata; data += stdata; data += ldata; data += mdata
	if term == "dxf"{
		title = title + "-detail.dxf"
	} else {
		title = title + "-detail.svg"
	}
	foldr := ""
	fname := genfname(foldr,title)
	pltstr = skriptrun(data, "plotcbmdet.gp", term, title, fname)
	
	if term == "svg" || term == "dxf" || term == "svgmono"{
		pltstr = fname
	}
	if web{
		if term != "dxf"{kass.Svgkong(pltstr)}
		pltstr = title
		
		//fmt.Println("Pltstr-",pltstr)
	}
	return

}


//PloFrmBmDet plots the steel detailing drawing of a (single) cbeam in a frame
func PlotFrmBmDet(web bool, bmvec []int, bms map[int][]*RccBm, folder, title, term string) (pltstr string){
	var xs float64
	var gdata, sdata, stdata, ldata, mdata, data string
	for _, id := range bmvec{
		bmarr := bms[id]
		g1, s1, st1, l1, m1 := BmSpanDraw(xs, bmarr)
		lspan := bmarr[1].Lspan
		xs += lspan * 1e3
		gdata += g1; sdata += s1; stdata += st1; ldata += l1; mdata += m1
	}
	gdata += "\n\n"; sdata += "\n\n"; stdata += "\n\n"
	ldata += "\n\n"
	data += gdata; data += sdata; data += stdata; data += ldata; data += mdata
	if term == "dxf"{
		title = title + "-detail.dxf"
	} else {
		title = title + "-detail.svg"
	}
	fname := genfname(folder,title)
	pltstr = skriptrun(data, "plotcbmdet.gp", term, title, fname)
	if term == "svg" || term == "dxf" || term == "svgmono"{
		pltstr = fname
	}
	if web{
		if term != "dxf"{kass.Svgkong(pltstr)}
		pltstr = title
	}
	return

}


//PlotSfDet plots the steel detailing drawing of a cbeam in a subframe
//you reap what you sow
func PlotSfDet(sf *SubFrm) (pltstr string){
	var xs float64
	var gdata, sdata, stdata, ldata, mdata, data string
	for _, id := range sf.Beams{
		bmarr := sf.RcBm[id]
		g1, s1, st1, l1, m1 := BmSpanDraw(xs, bmarr)
		lspan := bmarr[1].Lspan * 1e3
		xs += lspan
		gdata += g1; sdata += s1; stdata += st1; ldata += l1; mdata += m1
	}
	gdata += "\n\n"; sdata += "\n\n"; stdata += "\n\n"
	ldata += "\n\n"
	data += gdata; data += sdata; data += stdata; data += ldata; data += mdata
	title := sf.Title + "-detail.svg"
	if sf.Term == "dxf"{
		title = sf.Title + "-detail.dxf"
	}
	fname := genfname(sf.Foldr,title)
	pltstr = skriptrun(data, "plotcbmdet.gp", sf.Term,title, fname)
	if sf.Term == "svg" || sf.Term == "svgmono" || sf.Term == "dxf"{
		pltstr = fname
	}
	if sf.Web{
		if sf.Term != "dxf"{kass.Svgkong(pltstr)}
		sf.Txtplots = append(sf.Txtplots, title)
		pltstr = title
	}
	//b.Txtplot = append(b.Txtplot, pltstr)
	//fmt.Println("pltstr",pltstr,ColorReset)
	return

}


//BmSpanDraw details a single beam span - in MM
func BmSpanDraw(xs float64,bmarr []*RccBm)(gdata, sdata, stdata, ldata, adata string){
	//DrawMem2d returns 2d coords for drawing a member (front view)
 	//func DrawMem2d(mdx, styp int, pb, pe, dims []float64) (data string){
	//ALL DIMS IN MM ALL DIMS IN MM ALL DIMS IN MM
	bm := bmarr[1]
	lspan := bm.Lspan * 1e3
	pb := []float64{xs, bm.Dused/2.0}
	pe := []float64{xs + lspan, bm.Dused/2.0}
	
	if bm.Ldx == 1{
		pb[0] -= bm.Lsx/2.0
	}
	if bm.Rdx == 1{
		pe[0] += bm.Rsx/2.0
	}
	// var dims []float64
	// for _, dim := range bm.Dims{
	// 	dims = append(dims, dim)
	// }
	gdata += kass.DrawMem2d(bm.Mid, bm.Styp, pb, pe, bm.Dims)
	//add supports with depth Dused * 2.0
	dims := []float64{0.0,bm.Lsx}
	pb = []float64{xs, 0.0}
	pe = []float64{xs, -bm.Dused*2.0}
	gdata += kass.DrawMem2d(8, 1, pb, pe, dims)

	dims = []float64{0.0,bm.Rsx}
	pb = []float64{xs + lspan, 0.0}
	pe = []float64{xs + lspan, -bm.Dused*2.0}
	gdata += kass.DrawMem2d(8, 1, pb, pe, dims)

	//draw span dimensions
	yar := -bm.Dused*1.75
	adata += fmt.Sprintf("%f %f %f %f 1.0\n",xs, yar, lspan, 0.0)
	ldata += fmt.Sprintf("%f %f %.f %f 1.0\n",xs+lspan/2.0, yar, lspan, 0.0)
	
	//plot asts and ascs
	var x0, y0, x1, y1 float64
	//y0 = bm.Cvrt
	//using uniform cover
	y0 = 60.0
	y1 = (bm.Dused - bm.Cvrc)
	for i, b := range bmarr{
		// fmt.Println("bmarr->",i, bm.Printz())
		switch i{
			case 0:
			//x0 = xs //+ bm.Cvrt
			//x1 = xs + bm.Lsx	
			x0 = xs
			x1 = xs + bm.Lsx/2.0 + bm.CL[0]* 1e3
			if b.Asc > 0.0{
				n1 := int(b.Rbrc[0]); n2 := int(b.Rbrc[1]); d1 := b.Rbrc[2]; d2 := b.Rbrc[3]
				ldata += fmt.Sprintf("%f %f %vX%.f+%vX%.f 1.0\n",(x0+x1)/2.0, y1 + 0.1, n1, d1, n2, d2)
				sdata += fmt.Sprintf("%f %f %f %f 1.0 4.0\n",x0, y1, x1 - x0, 0.0)
				if bm.Ldx == 1{
					ldreq := bm.CL[4] * 1e3
					l0 := ldreq/3.0
					dx := bm.Lsx/2.0 - bm.Cvrt
					dy := l0 - dx
					x2 := xs - dx
					y2 := y1 - dy
					//add top left anchorage of l0					
					sdata += fmt.Sprintf("%f %f %f %f 1.0 4.0\n",x2, y1, dx, 0.0)
					sdata += fmt.Sprintf("%f %f %f %f 1.0 4.0\n",x2, y2, 0.0, dy)
					//ldata += fmt.Sprintf("%f %f %.f\n",x2 - 0.5, (y1 + y2)/2.0, l0)
				}
			}
			if b.Ast > 0.0{
				n1 := int(b.Rbrt[0]); n2 := int(b.Rbrt[1]); d1 := b.Rbrt[2]; d2 := b.Rbrt[3]
				ldata += fmt.Sprintf("%f %f %vX%.f-%vX%.f 2.0\n",(x0+x1)/2.0, y0 - 100.0, n1, d1, n2, d2)
				sdata += fmt.Sprintf("%f %f %f %f 2.0 4.0\n",x0, y0 + 50.0, x1 - x0, 0.0)
			}
			yar := b.Dused + 200.0
			adata += fmt.Sprintf("%f %f %f %f 1.0\n",x0, yar, x1 - x0, 0.0)
			ldata += fmt.Sprintf("%f %f %.f %f 1.0\n",(x0+x1)/2.0, yar, x1 - x0, 0.0)
			case 1:
			//x0 = xs + bm.Cvrt
			
			x0 = xs + b.CL[2]* 1e3
			//x1 = xs + lspan			
			x1 = xs + lspan - b.CL[3] * 1e3
			if b.Asc > 0.0{
				n1 := int(b.Rbrc[0]); n2 := int(b.Rbrc[1]); d1 := b.Rbrc[2]; d2 := b.Rbrc[3]
				ldata += fmt.Sprintf("%f %f %vX%.f+%vX%.f 3.0\n",(x0+x1)/2.0, y1 + 100, n1, d1, n2, d2)
				sdata += fmt.Sprintf("%f %f %f %f 3.0 4.0\n",x0, y1 - 25.0, x1 - x0, 0.0)
			}
			if b.Ast > 0.0{
				n1 := int(b.Rbrt[0]); n2 := int(b.Rbrt[1]); d1 := b.Rbrt[2]; d2 := b.Rbrt[3]
				ldata += fmt.Sprintf("%f %f %vX%.f+%vX%.f 4.0\n",(x0+x1)/2.0, y0 - 0.1, n1, d1, n2, d2)
				sdata += fmt.Sprintf("%f %f %f %f 4.0 4.0\n",x0, y0 + 25.0, x1 - x0, 0.0)
				//2nd line continues into the support
				sdata += fmt.Sprintf("%f %f %f %f 4.0 4.0\n",xs, y0 , lspan, 0.0)
			}
			yar := -b.Dused/2.0
			adata += fmt.Sprintf("%f %f %f %f 1.0\n",x0, yar, x1 - x0, 0.0)
			adata += fmt.Sprintf("%f %f %f %f 1.0\n",xs, yar, x0 - xs, 0.0)
			adata += fmt.Sprintf("%f %f %f %f 1.0\n",x1, yar, b.CL[3]* 1e3, 0.0)
			
			ldata += fmt.Sprintf("%f %f %.f %f 1.0\n",(x0+x1)/2.0, yar, x1 - x0, 0.0)
			//ldata += fmt.Sprintf("%f %f %.2fm %f 1.0\n",x0/2.0, yar, b.CL[2], 0.0)
			//ldata += fmt.Sprintf("%f %f %.2fm %f 1.0\n",(x1+b.CL[3])/2.0, yar, b.CL[3], 0.0)
			case 2:
			x0 = xs + lspan - bm.CL[1] * 1e3 - bm.Rsx/2.0
			x1 = xs + lspan
			
			if b.Asc > 0.0{
				n1 := int(b.Rbrc[0]); n2 := int(b.Rbrc[1]); d1 := b.Rbrc[2]; d2 := b.Rbrc[3]
				ldata += fmt.Sprintf("%f %f %vX%.f+%vX%.f 5.0\n",(x0+x1)/2.0, y1 + 0.1, n1, d1, n2, d2)
				sdata += fmt.Sprintf("%f %f %f %f 5.0 4.0\n",x0, y1, x1 - x0, 0.0)
				if bm.Rdx == 1{
					ldreq := bm.CL[5] * 1e3
					l0 := ldreq/3.0
					dx := bm.Rsx/2.0 - bm.Cvrt
					dy := l0 - dx
					x2 := x1 + dx
					y2 := y1 - dy
					//add top right anchorage of l0					
					sdata += fmt.Sprintf("%f %f %f %f 1.0 5.0\n",x1, y1, dx, 0.0)
					sdata += fmt.Sprintf("%f %f %f %f 1.0 5.0\n",x2, y2, 0.0, dy)
					//ldata += fmt.Sprintf("%f %f %.f\n",x2 - 0.5, (y1 + y2)/2.0, l0)
				}
			}
			if b.Ast > 0.0{
				n1 := int(b.Rbrt[0]); n2 := int(b.Rbrt[1]); d1 := b.Rbrt[2]; d2 := b.Rbrt[3]
				ldata += fmt.Sprintf("%f %f %vX%.f+%vX%.f 6.0\n",(x0+x1)/2.0, y0 - 100.0, n1, d1, n2, d2)
				sdata += fmt.Sprintf("%f %f %f %f 6.0 4.0\n",x0, y0 + 50.0, x1 - x0, 0.0)
			}
			
			yar := b.Dused + 200.0
			adata += fmt.Sprintf("%f %f %f %f 1.0\n",x0, yar, x1 - x0, 0.0)
			ldata += fmt.Sprintf("%f %f %.f %f 1.0\n",(x0+x1)/2.0, yar, x1 - x0, 0.0)
		}
		//plot links
		slink, smin, snom := bm.Lspc[0],bm.Lspc[1],bm.Lspc[2]
		yar := -bm.Dused
		for i := range []string{"main","min","nominal"}{
			spc := slink
			//nlx := bm.Nlx[i]
			//float64{slink, smin, snom, mainlen, minlen, nomlen}
			l1 := xs + 0.0; l2 := xs + bm.L1*1e3; l3 := xs + bm.L4*1e3; l4 := xs + bm.Xs[20]*1e3
			xl, yl := 0.0, b.Cvrt
			dy := b.Dused - b.Cvrc - yl
			switch i{
				case 0:
				n1 := int(math.Ceil((l2-l1)/spc))
				lx1 := xs + 1.1*bm.L1*1e3/2.0
				ldata += fmt.Sprintf("%f %f %v-%.f 1.0\n", lx1, yar + 25.0, n1, spc)
				ldata += fmt.Sprintf("%f %f %.f 1.0\n", lx1, yar - 100.0, l2-l1)
				adata += fmt.Sprintf("%f %f %.f %f 1.0 3.0\n", l1, yar, l2 - l1, 0.0)
				xl = l1
				for j := 1; j <= n1; j++{
					stdata += fmt.Sprintf("%f %f %f %f 1.0 %.f\n",xl, yl, 0.0, dy, b.Dlink)
					xl += spc
				}
				n2 := int(math.Ceil((l4-l3)/spc))
				if n2 <= 1{continue}
				xl = l3
				lx1 = xs + (bm.L4 + bm.Xs[20])*1e3/2.0
				ldata += fmt.Sprintf("%f %f %v-%.f 1.0\n", lx1, yar + 25.0, n2, spc)
				ldata += fmt.Sprintf("%f %f %.f 1.0\n", lx1, yar - 100, l4-l3)
				adata += fmt.Sprintf("%f %f %.f %f 1.0 3.0\n", l3, yar, l4 - l3, 0.0)
				for j := 1; j <= n2; j++{
					stdata += fmt.Sprintf("%f %f %f %f 1.0 %.f\n",xl, yl, 0.0, dy, b.Dlink)
					xl += spc
				}
				case 1:
				spc = smin
				l1 = xs + bm.L1*1e3; l2 = xs + bm.L2*1e3; l3 = xs + bm.L3*1e3; l4 = xs + bm.L4*1e3
				n1 := int(math.Ceil((l2-l1)/spc))
				
				xl = l1				
				lx1 := xs + (bm.L1 + bm.L2)*1e3/2.0				
				if l2 - l1 > 0.0{
					//fmt.Println("l2-l1")
					adata += fmt.Sprintf("%f %f %.f %f 2.0 %.f\n", l1, yar, l2 - l1, 0.0, b.Dlink)
					ldata += fmt.Sprintf("%f %f %.f 2.0\n", lx1, yar - 100.0, l2-l1)
					ldata += fmt.Sprintf("%f %f %v-%.f 1.0\n", lx1, yar + 25.0, n1, spc)
					for j := 1; j <= n1; j++{
						stdata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n",xl, yl, 0.0, dy)
						xl += spc
					}
				}
				n2 := int(math.Ceil((l4-l3)/spc))
				xl = l3
				if l4 - l3 > 0.0{
					
					//fmt.Println("l4-l3")
					lx1 = xs + (bm.L3 + bm.L4)*1e3/2.0
					adata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n", l3, yar, l4 - l3, 0.0)
					ldata += fmt.Sprintf("%f %f %v-%.f 1.0\n", lx1, yar + 25.0, n2, spc)
					ldata += fmt.Sprintf("%f %f %.f 2.0\n", lx1, yar - 100.0, l4-l3)
					for j := 1; j <= n2; j++{
						stdata += fmt.Sprintf("%f %f %f %f 2.0 %.f\n",xl, yl, 0.0, dy, b.Dlink)
						xl += spc
					}	
				}
				case 2:
				spc = snom
				l1 = xs + bm.L2*1e3; l2 = xs + bm.L3*1e3
				//l3 = xs + 0.0; l4 = xs + 0.0 
				n1 := int(math.Ceil((l2-l1)/spc))
				xl = l1
				lx1 := xs + (bm.L2 + bm.L3)*1e3/2.0
				adata += fmt.Sprintf("%f %f %.f %f 3.0 3.0\n", l1, yar, l2 - l1, 0.0)
				ldata += fmt.Sprintf("%f %f %v-%.f 1.0\n", lx1, yar + 25.0, n1, spc)
				ldata += fmt.Sprintf("%f %f %.f 2.0\n", lx1, yar - 100.0, l2-l1)
				for j := 1; j <= n1; j++{
					stdata += fmt.Sprintf("%f %f %f %f 3.0 %.f\n",xl, yl, 0.0, dy, b.Dlink)
					xl += spc
				}				
			}
		}
		//plot label/dim arrows
		//plot cfxs
		/*
		if len(bm.Cfxs) >= 4{		
			l1, l2, l3, l4 := xs + bm.Cfxs[0],xs + bm.Cfxs[1],xs + bm.Cfxs[2],xs + bm.Cfxs[3]
			yar = b.Dused + 0.1
			adata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n", xs, yar, l1 - xs, 0.0)
			adata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n", l2, yar, pe[0] - l2, 0.0)

			yar = -b.Dused/250.0 + 0.1
			adata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n", xs, yar, l3 - xs, 0.0)
			adata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n", l4, yar, pe[0] - l4, 0.0)
		}
		*/
	}
	return
}


//YE OLDE METER VERSION
// func BmSpanDraw(xs float64,bmarr []*RccBm)(gdata, sdata, stdata, ldata, adata string){
// 	//DrawMem2d returns 2d coords for drawing a member (front view)
// 	//func DrawMem2d(mdx, styp int, pb, pe, dims []float64) (data string){
// 	bm := bmarr[1]
// 	pb := []float64{xs, bm.Dused/1000.0/2.0}
// 	pe := []float64{xs + bm.Lspan, bm.Dused/1000.0/2.0}
	
// 	if bm.Ldx == 1{
// 		pb[0] -= bm.Lsx/2.0/1000.0
// 	}
// 	if bm.Rdx == 1{
// 		pe[0] += bm.Rsx/2.0/1000.0
// 	}
// 	var dims []float64
// 	for _, dim := range bm.Dims{
// 		dims = append(dims, dim/1000.0)
// 	}
// 	gdata += kass.DrawMem2d(bm.Mid, bm.Styp, pb, pe, dims)
// 	//add supports with depth Dused * 2.0
// 	dims = []float64{0.0,bm.Lsx/1000.0}
// 	pb = []float64{xs, 0.0}
// 	pe = []float64{xs, -bm.Dused/500.0}
// 	gdata += kass.DrawMem2d(8, 1, pb, pe, dims)

// 	dims = []float64{0.0,bm.Rsx/1000.0}
// 	pb = []float64{xs + bm.Lspan, 0.0}
// 	pe = []float64{xs + bm.Lspan, -bm.Dused/500.0}
// 	gdata += kass.DrawMem2d(8, 1, pb, pe, dims)

// 	//draw span dimensions
// 	yar := -bm.Dused/475.0
// 	adata += fmt.Sprintf("%f %f %f %f 1.0\n",xs, yar, bm.Lspan, 0.0)
// 	ldata += fmt.Sprintf("%f %f %.2fm %f 1.0\n",xs+bm.Lspan/2.0, yar, bm.Lspan, 0.0)
	
// 	//plot asts and ascs
// 	var x0, y0, x1, y1 float64
// 	//y0 = bm.Cvrt/1000.0
// 	//using uniform cover
// 	y0 = 60.0/1000.0
// 	y1 = (bm.Dused - bm.Cvrc)/1000.0
// 	for i, b := range bmarr{
// 		switch i{
// 			case 0:
// 			//x0 = xs //+ bm.Cvrt/1000.0
// 			//x1 = xs + bm.Lsx	
// 			x0 = xs
// 			x1 = xs + bm.Lsx/2.0/1000.0 + bm.CL[0]
// 			if b.Asc > 0.0{
// 				n1 := int(b.Rbrc[0]); n2 := int(b.Rbrc[1]); d1 := b.Rbrc[2]; d2 := b.Rbrc[3]
// 				ldata += fmt.Sprintf("%f %f %vX%.f+%vX%.f 1.0\n",(x0+x1)/2.0, y1 + 0.1, n1, d1, n2, d2)
// 				sdata += fmt.Sprintf("%f %f %f %f 1.0 4.0\n",x0, y1, x1 - x0, 0.0)
// 				if bm.Ldx == 1{
// 					ldreq := bm.CL[4]
// 					l0 := ldreq/3.0
// 					dx := bm.Lsx/1000.0/2.0 - bm.Cvrt/1000.0
// 					dy := l0 - dx
// 					x2 := xs - dx
// 					y2 := y1 - dy
// 					//add top left anchorage of l0					
// 					sdata += fmt.Sprintf("%f %f %f %f 1.0 4.0\n",x2, y1, dx, 0.0)
// 					sdata += fmt.Sprintf("%f %f %f %f 1.0 4.0\n",x2, y2, 0.0, dy)
// 					//ldata += fmt.Sprintf("%f %f %.f\n",x2 - 0.5, (y1 + y2)/2.0, l0)
// 				}
// 			}
// 			if b.Ast > 0.0{
// 				n1 := int(b.Rbrt[0]); n2 := int(b.Rbrt[1]); d1 := b.Rbrt[2]; d2 := b.Rbrt[3]
// 				ldata += fmt.Sprintf("%f %f %vX%.f-%vX%.f 2.0\n",(x0+x1)/2.0, y0 - 0.1, n1, d1, n2, d2)
// 				sdata += fmt.Sprintf("%f %f %f %f 2.0 4.0\n",x0, y0 + 0.05, x1 - x0, 0.0)
// 			}
// 			yar := b.Dused/1000.0 + 0.2
// 			adata += fmt.Sprintf("%f %f %f %f 1.0\n",x0, yar, x1 - x0, 0.0)
// 			ldata += fmt.Sprintf("%f %f %.2fm %f 1.0\n",(x0+x1)/2.0, yar, x1 - x0, 0.0)
// 			case 1:
// 			//x0 = xs + bm.Cvrt/1000.0
			
// 			x0 = xs + b.CL[2]
// 			//x1 = xs + bm.Lspan			
// 			x1 = xs + b.Lspan - b.CL[3]
// 			if b.Asc > 0.0{
// 				n1 := int(b.Rbrc[0]); n2 := int(b.Rbrc[1]); d1 := b.Rbrc[2]; d2 := b.Rbrc[3]
// 				ldata += fmt.Sprintf("%f %f %vX%.f+%vX%.f 3.0\n",(x0+x1)/2.0, y1 + 0.1, n1, d1, n2, d2)
// 				sdata += fmt.Sprintf("%f %f %f %f 3.0 4.0\n",x0, y1 - 0.025, x1 - x0, 0.0)
// 			}
// 			if b.Ast > 0.0{
// 				n1 := int(b.Rbrt[0]); n2 := int(b.Rbrt[1]); d1 := b.Rbrt[2]; d2 := b.Rbrt[3]
// 				ldata += fmt.Sprintf("%f %f %vX%.f+%vX%.f 4.0\n",(x0+x1)/2.0, y0 - 0.1, n1, d1, n2, d2)
// 				sdata += fmt.Sprintf("%f %f %f %f 4.0 4.0\n",x0, y0 + 0.025, x1 - x0, 0.0)
// 				//2nd line continues into the support
// 				sdata += fmt.Sprintf("%f %f %f %f 4.0 4.0\n",xs, y0 , b.Lspan, 0.0)
// 			}
// 			yar := -b.Dused/2000.0
// 			adata += fmt.Sprintf("%f %f %f %f 1.0\n",x0, yar, x1 - x0, 0.0)
// 			adata += fmt.Sprintf("%f %f %f %f 1.0\n",xs, yar, x0 - xs, 0.0)
// 			adata += fmt.Sprintf("%f %f %f %f 1.0\n",x1, yar, b.CL[3], 0.0)
			
// 			ldata += fmt.Sprintf("%f %f %.2fm %f 1.0\n",(x0+x1)/2.0, yar, x1 - x0, 0.0)
// 			//ldata += fmt.Sprintf("%f %f %.2fm %f 1.0\n",x0/2.0, yar, b.CL[2], 0.0)
// 			//ldata += fmt.Sprintf("%f %f %.2fm %f 1.0\n",(x1+b.CL[3])/2.0, yar, b.CL[3], 0.0)
// 			case 2:
// 			x0 = xs + bm.Lspan - bm.CL[1] - bm.Rsx/1000.0/2.0
// 			x1 = xs + bm.Lspan
			
// 			if b.Asc > 0.0{
// 				n1 := int(b.Rbrc[0]); n2 := int(b.Rbrc[1]); d1 := b.Rbrc[2]; d2 := b.Rbrc[3]
// 				ldata += fmt.Sprintf("%f %f %vX%.f+%vX%.f 5.0\n",(x0+x1)/2.0, y1 + 0.1, n1, d1, n2, d2)
// 				sdata += fmt.Sprintf("%f %f %f %f 5.0 4.0\n",x0, y1, x1 - x0, 0.0)
// 				if bm.Rdx == 1{
// 					ldreq := bm.CL[5]
// 					l0 := ldreq/3.0
// 					dx := bm.Rsx/1000.0/2.0 - bm.Cvrt/1000.0
// 					dy := l0 - dx
// 					x2 := x1 + dx
// 					y2 := y1 - dy
// 					//add top right anchorage of l0					
// 					sdata += fmt.Sprintf("%f %f %f %f 1.0 5.0\n",x1, y1, dx, 0.0)
// 					sdata += fmt.Sprintf("%f %f %f %f 1.0 5.0\n",x2, y2, 0.0, dy)
// 					//ldata += fmt.Sprintf("%f %f %.f\n",x2 - 0.5, (y1 + y2)/2.0, l0)
// 				}
// 			}
// 			if b.Ast > 0.0{
// 				n1 := int(b.Rbrt[0]); n2 := int(b.Rbrt[1]); d1 := b.Rbrt[2]; d2 := b.Rbrt[3]
// 				ldata += fmt.Sprintf("%f %f %vX%.f+%vX%.f 6.0\n",(x0+x1)/2.0, y0 - 0.1, n1, d1, n2, d2)
// 				sdata += fmt.Sprintf("%f %f %f %f 6.0 4.0\n",x0, y0 + 0.05, x1 - x0, 0.0)
// 			}
			
// 			yar := b.Dused/1000.0 + 0.2
// 			adata += fmt.Sprintf("%f %f %f %f 1.0\n",x0, yar, x1 - x0, 0.0)
// 			ldata += fmt.Sprintf("%f %f %.2fm %f 1.0\n",(x0+x1)/2.0, yar, x1 - x0, 0.0)
// 		}
// 		//plot links
// 		slink, smin, snom := bm.Lspc[0],bm.Lspc[1],bm.Lspc[2]
// 		yar := -bm.Dused/1000.0
// 		for i := range []string{"main","min","nominal"}{
// 			spc := slink
// 			//nlx := bm.Nlx[i]
// 			//float64{slink, smin, snom, mainlen, minlen, nomlen}
// 			l1 := xs + 0.0; l2 := xs + bm.L1; l3 := xs + bm.L4; l4 := xs + bm.Xs[20]
// 			xl, yl := 0.0, b.Cvrt/1000.0
// 			dy := b.Dused/1000.0 - b.Cvrc/1000.0 - yl
// 			switch i{
// 				case 0:
// 				n1 := int(math.Ceil((l2-l1)*1000.0/spc))
// 				lx1 := xs + bm.L1/2.0
// 				ldata += fmt.Sprintf("%f %f %v-%.f 1.0\n", lx1, yar + 0.025, n1, spc)
// 				ldata += fmt.Sprintf("%f %f %.2fm 1.0\n", lx1, yar - 0.1, l2-l1)
// 				adata += fmt.Sprintf("%f %f %f %f 1.0 3.0\n", l1, yar, l2 - l1, 0.0)
// 				xl = l1
// 				for j := 1; j <= n1; j++{
// 					stdata += fmt.Sprintf("%f %f %f %f 1.0 %.f\n",xl, yl, 0.0, dy, b.Dlink)
// 					xl += spc/1000.0
// 				}
// 				n2 := int(math.Ceil((l4-l3)*1000.0/spc))
// 				//if n2 <= 1{continue}
// 				xl = l3
// 				lx1 = xs + (bm.L4 + bm.Xs[20])/2.0
// 				ldata += fmt.Sprintf("%f %f %v-%.f 1.0\n", lx1, yar + 0.025, n2, spc)
// 				ldata += fmt.Sprintf("%f %f %.2fm 1.0\n", lx1, yar - 0.1, l4-l3)
// 				adata += fmt.Sprintf("%f %f %f %f 1.0 3.0\n", l3, yar, l4 - l3, 0.0)
// 				for j := 1; j <= n2; j++{
// 					stdata += fmt.Sprintf("%f %f %f %f 1.0 %.f\n",xl, yl, 0.0, dy, b.Dlink)
// 					xl += spc/1000.0
// 				}
// 				case 1:
// 				spc = smin
// 				l1 = xs + bm.L1; l2 = xs + bm.L2; l3 = xs + bm.L3; l4 = xs + bm.L4
// 				n1 := int(math.Ceil((l2-l1)*1000.0/spc))
				
// 				xl = l1				
// 				lx1 := xs + (bm.L1 + bm.L2)/2.0				
// 				if l2 - l1 > 0.0{
// 					//fmt.Println("l2-l1")
// 					adata += fmt.Sprintf("%f %f %f %f 2.0 %.f\n", l1, yar, l2 - l1, 0.0, b.Dlink)
// 					ldata += fmt.Sprintf("%f %f %.2fm 2.0\n", lx1, yar - 0.1, l2-l1)
// 					ldata += fmt.Sprintf("%f %f %v-%.f 1.0\n", lx1, yar + 0.025, n1, spc)
// 					for j := 1; j <= n1; j++{
// 						stdata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n",xl, yl, 0.0, dy)
// 						xl += spc/1000.0
// 					}
// 				}
// 				n2 := int(math.Ceil((l4-l3)*1000.0/spc))
// 				xl = l3
// 				if l4 - l3 > 0.0{
					
// 					//fmt.Println("l4-l3")
// 					lx1 = xs + (bm.L3 + bm.L4)/2.0
// 					adata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n", l3, yar, l4 - l3, 0.0)
// 					ldata += fmt.Sprintf("%f %f %v-%.f 1.0\n", lx1, yar + 0.025, n2, spc)
// 					ldata += fmt.Sprintf("%f %f %.2fm 2.0\n", lx1, yar - 0.1, l4-l3)
// 					for j := 1; j <= n2; j++{
// 						stdata += fmt.Sprintf("%f %f %f %f 2.0 %.f\n",xl, yl, 0.0, dy, b.Dlink)
// 						xl += spc/1000.0
// 					}	
// 				}
// 				case 2:
// 				spc = snom
// 				l1 = xs + bm.L2; l2 = xs + bm.L3; l3 = xs + 0.0; l4 = xs + 0.0 
// 				n1 := int(math.Ceil((l2-l1)*1000.0/spc))
// 				xl = l1
// 				lx1 := xs + (bm.L2 + bm.L3)/2.0
// 				adata += fmt.Sprintf("%f %f %f %f 3.0 3.0\n", l1, yar, l2 - l1, 0.0)
// 				ldata += fmt.Sprintf("%f %f %v-%.f 1.0\n", lx1, yar + 0.025, n1, spc)
// 				ldata += fmt.Sprintf("%f %f %.2fm 2.0\n", lx1, yar - 0.1, l2-l1)
// 				for j := 1; j <= n1; j++{
// 					stdata += fmt.Sprintf("%f %f %f %f 3.0 %.f\n",xl, yl, 0.0, dy, b.Dlink)
// 					xl += spc/1000.0
// 				}				
// 			}
// 		}
// 		//plot label/dim arrows
// 		//plot cfxs
// 		/*
// 		if len(bm.Cfxs) >= 4{		
// 			l1, l2, l3, l4 := xs + bm.Cfxs[0],xs + bm.Cfxs[1],xs + bm.Cfxs[2],xs + bm.Cfxs[3]
// 			yar = b.Dused/1000.0 + 0.1
// 			adata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n", xs, yar, l1 - xs, 0.0)
// 			adata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n", l2, yar, pe[0] - l2, 0.0)

// 			yar = -b.Dused/250.0 + 0.1
// 			adata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n", xs, yar, l3 - xs, 0.0)
// 			adata += fmt.Sprintf("%f %f %f %f 2.0 3.0\n", l4, yar, pe[0] - l4, 0.0)
// 		}
// 		*/
// 	}
// 	return
// }

//PlotBmGeom plots an rcc beam section
func PlotBmGeom(b *RccBm, term string) (pltstr string){
	var title, fname string
	data, err := b.Draw()
	if err != nil{
		log.Println(err)
		return
	}
	
	b.Title = fmt.Sprintf("rcc-beam-%s-sec%v-%v", b.Title,b.Mid, b.Id)
	title = b.Title + ".svg"
	if term == "dxf"{
		title = b.Title + ".dxf"
	
	}
	fname = genfname(b.Foldr,title)
	pltstr = skriptrun(data, "plotbmgeom.gp", term, b.Title, fname)
	if term == "svg" || term == "svgmono" || term == "dxf"{
		pltstr = fname
	}
	if b.Web{
		//embed kongtext font in svg
		if term != "dxf"{kass.Svgkong(pltstr)}
		b.Txtplots = append(b.Txtplots, title)
		pltstr = title
	}
	//fmt.Println(ColorRed,"bm pltstr",pltstr,ColorReset)
	return
}

//PlotColNM plots column section n-m interaction curves
func PlotColNM(pus, mus []float64) (pltstr string) {
	//get plot script filepath
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	pltskript := filepath.Join(basepath,"/plotcolnm.gp")
	var data string
	for idx, pu := range pus {
		mu := mus[idx]
		data += fmt.Sprintf("%v %v\n", pu, mu)
	}
	f, e1 := os.CreateTemp("", "mosh")
	if e1 != nil {
		fmt.Println(e1)
	}
	defer f.Close()
	defer os.Remove(f.Name())	
	_, e1 = f.WriteString(data)
	if e1 != nil {
		fmt.Println(e1)
	}
	cmd := exec.Command("gnuplot","-c",pltskript,f.Name(),"dumb")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outstr, errstr := stdout.String(), stderr.String()
	if err != nil {
		fmt.Println(err)
	}
	if errstr != "" {
		fmt.Println(errstr)
	}
	return outstr
	
}

//DrawColCircle plots a circular rcc column section 
func DrawColCircle(c *RccCol, term string, multi bool) (pltstr string){
	//draw shape of column with dimensions
	c.Txtplot = []string{}
	if c.Title == ""{
		c.Title = fmt.Sprintf("rcc col section %v-%v", c.Mid, c.Id)
	}
	var data string
	nbars := len(c.Dias)
	//log.Println("nbars->",nbars)
	theta := math.Pi * 2.0/float64(nbars)
	rs := (c.B - c.Cvrc - c.Cvrt)/2.0
	rc := c.B/2.0
	data += fmt.Sprintf("0.0 0.0 %.1f 3 3\n",rc)
	data += fmt.Sprintf("0.0 0.0 %.1f 2 2\n",rs)
	data += "\n\n"
	var xc, yc, alpha float64
	alpha -= theta
	
	for i:=0;i < nbars;i++{
		alpha += theta
		xc = rs * math.Cos(alpha)
		yc = rs * math.Sin(alpha)
		rd := c.Dias[i]
		data += fmt.Sprintf("%.1f %.1f %.f 1 1\n", yc, xc, rd)
		//log.Println("rebar pts",xc, yc, rd)
	}
	data += "\n\n"
	if multi{
		pltstr = data
		return
	}
	title := c.Title + ".svg"
	if term == "dxf"{
		title = c.Title + ".dxf"
	}
	fname := genfname(c.Foldr,title)
	pltstr = skriptrun(data, "plotcolcircle.gp", term, c.Title, fname)
	if term == "svg" || term == "svgmono" || term == "dxf"{
		pltstr = fname
	}
	
	//fmt.Println("PLOTSTR->\n",pltstr)
	c.Txtplot = append(c.Txtplot, pltstr)
	return
}

//DrawColRect plots a rectangular rcc column section
//if multi it returns a 
func DrawColRect(c *RccCol, term string, multi bool) (pltstr string){
	//draw shape of column with dimensions
	//draw steel bars with labels
	//draw ties with spacing
	var data string
	data += fmt.Sprintf("%.1f %.1f %.1f %.1f 3 3\n",c.B/2.0,c.H/2.0,c.B/2.0,c.H/2.0)
	data += fmt.Sprintf("%.1f %.1f %.1f %.1f 2 2\n",c.B/2.0,c.H/2.0,c.B/2.0 - c.Cvrc + c.Dtie,c.H/2.0 - c.Cvrt + c.Dtie)
	data += "\n\n"
	for _, pt := range c.Barpts{
		//fmt.Println(pt)
		data += fmt.Sprintf("%.1f %.1f %.1f 1 1\n",pt[0], pt[1], pt[2]/2.0)
	}
	data += "\n\n"
	if multi{
		pltstr = data
		return
	}
	if c.Title == ""{
		c.Title = fmt.Sprintf("rcc col section %v-%v", c.Mid, c.Id)
	}
	title := c.Title + ".svg"
	fname := genfname(c.Foldr,title)
	pltstr = skriptrun(data, "plotcolrect.gp", term, c.Title, fname)
	if term == "svg" || term == "svgmono" || term == "dxf"{
		pltstr = fname
	}
	c.Txtplot = append(c.Txtplot, pltstr)
	return
}

//PlotColDet plots a 2d (front) view of a column
func (c *RccCol) PlotColDet()(pltstr string){
	var data, stdata, tdata, adata, ldata string
	h := c.Dims[1]
	var x0, x1, xar, xspc, y0, y1, yspc float64
	x0 = c.Cvrt
	x1 = h - c.Cvrt
	y1 = c.Lspan * 1000.0
	pb := []float64{h/2.0,y0}
	pe := []float64{h/2.0, y1}
	data += kass.DrawMem2d(c.Id, c.Styp, pb, pe, c.Dims)
	xspc = (h - 2.0 * c.Cvrt)/float64(c.Nlayers-1)	
	//set this for non rect columns (for now)
	//later sort thru barpts for unique y values
	if c.Nlayers == 0{c.Nlayers = 2}
	//draw longitudinal steel
	for i := 0; i < c.Nlayers; i++{
		stdata += fmt.Sprintf("%f %f %f %f %f\n",x0,y0,0.0,y1,1.0)
		x0 += xspc
	}
	x0 = c.Cvrt
	//bot conf
	y1 = c.Lp * 1000.0
	xar = 1.2 * h
	adata += fmt.Sprintf("%f %f %f %f %f\n",xar,y0,0.0,y1,2.0)
	ldata += fmt.Sprintf("%f %f %.fmm\n",xar*2, (y0+y1)/2.0, y1-y0)

	//plot bot ties 
	yspc = c.Ptiec
	for y := y0; y <= y1; y += yspc{
		tdata += fmt.Sprintf("%f %f %f %f %f\n",x0,y,x1-x0,0.0,2.0)
	}
	//plot bot tie arrows 
	xar = -0.2 * h
	adata += fmt.Sprintf("%f %f %f %f %f\n",xar,y0,0.0,y1,2.0)
	ldata += fmt.Sprintf("%f %f T%.f-%.fmm\n",xar*2, (y0+y1)/2.0, c.Dtie,c.Ptiec)

	//mid conf
	y0 = c.Lp * 1000.0
	y1 = (c.Lspan - c.Lp) * 1000.0
	xar = 1.2 * h
	adata += fmt.Sprintf("%f %f %f %f %f\n",xar,y0,0.0,y1,1.0)
	ldata += fmt.Sprintf("%f %f %.fmm\n",xar*2, (y0+y1)/2.0, y1-y0)
	
	//plot mid ties
	yspc = c.Ptie
	for y := y0; y <= y1; y += yspc{
		tdata += fmt.Sprintf("%f %f %f %f %f\n",x0,y,x1-x0,0.0,3.0)
	}
	//plot mid tie arrows
	xar = -0.2 * h
	adata += fmt.Sprintf("%f %f %f %f %f\n",xar,y0,0.0,y1-y0,3.0)
	ldata += fmt.Sprintf("%f %f T%.f-%.fmm\n",xar*2, (y0+y1)/2.0, c.Dtie,c.Ptie)

	//top conf
	y0 = (c.Lspan - c.Lp) * 1000.0
	y1 = c.Lspan * 1000.0
	xar = 1.2 * h
	adata += fmt.Sprintf("%f %f %f %f %f\n",xar,y0,0.0,y1-y0,2.0)
	ldata += fmt.Sprintf("%f %f %.fmm\n",xar*2, (y0+y1)/2.0, y1-y0)
	
	//plot top ties
	yspc = c.Ptiec
	for y := y0; y <= y1; y += yspc{
		tdata += fmt.Sprintf("%f %f %f %f %f\n",x0,y,x1-x0,0.0,2.0)
	}
	//plot top tie arrows
	xar = -0.2 * h
	adata += fmt.Sprintf("%f %f %f %f %f\n",xar,y0,0.0,y1-y0,2.0)
	ldata += fmt.Sprintf("%f %f T%.f-%.fmm\n",xar*2, (y0+y1)/2.0, c.Dtie, c.Ptiec)

	data += "\n\n"; stdata += "\n\n"; tdata += "\n\n";adata += "\n\n"; ldata += "\n\n"
	data += stdata; data += tdata; data += adata; data += ldata
	switch c.Styp{
		case 0:
		data += DrawColCircle(c, c.Term, true)
		default:
		data += c.Draw()
	}
	if c.Title == ""{
		c.Title = fmt.Sprintf("rcc col view %v-%v", c.Mid, c.Id)
	}
	title := c.Title + ".svg"
	if c.Term == "dxf"{
		title = c.Title + ".dxf"
	}
	fname := genfname(c.Foldr,title)
	pltstr = skriptrun(data, "plotcolview.gp", c.Term, c.Title, fname)
	if c.Term == "svg" || c.Term == "dxf" || c.Term == "svgmono"{
		pltstr = fname
	}
	c.Txtplot = append(c.Txtplot, pltstr)

	if c.Web{
		if c.Term != "dxf"{kass.Svgkong(pltstr)}
		c.Txtplots = append(c.Txtplots, title)
		pltstr = title
	}
	return
}

//PlotFtng plots an rcc footing 
func PlotFtng(colx, coly, fck, fy, hf, eo, d, dmin, lx, ly, nomcvr float64, sloped bool, rez []float64, plot string) (err error){
	var data string
	if hf == 0.0{hf = 1.5 + d}
	switch sloped{
		case true:		
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0, -ly/2.0, 0.0)
		data += fmt.Sprintf("%f %f %f\n", lx/2.0, -ly/2.0, 0.0)
		data += fmt.Sprintf("%f %f %f\n", lx/2.0,  ly/2.0, 0.0)
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0,  ly/2.0, 0.0)
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0, -ly/2.0, 0.0)
		data += "\n"
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0, -ly/2.0, dmin)
		data += fmt.Sprintf("%f %f %f\n", lx/2.0, -ly/2.0, dmin)
		data += fmt.Sprintf("%f %f %f\n", lx/2.0,  ly/2.0, dmin)
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0,  ly/2.0, dmin)
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0, -ly/2.0, dmin)
		data += "\n"
		//column
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0, -coly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n", colx/2.0, -coly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n", colx/2.0,  coly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0,  coly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0, -coly/2.0, d)
		data += "\n"
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0, -coly/2.0, d + hf)
		data += fmt.Sprintf("%f %f %f\n", colx/2.0, -coly/2.0, d + hf)
		data += fmt.Sprintf("%f %f %f\n", colx/2.0,  coly/2.0, d + hf)
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0,  coly/2.0, d + hf)
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0, -coly/2.0, d + hf)
		data += "\n"
		case false:
		//plot base polygons
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0, -ly/2.0, 0.0)
		data += fmt.Sprintf("%f %f %f\n", lx/2.0, -ly/2.0, 0.0)
		data += fmt.Sprintf("%f %f %f\n", lx/2.0,  ly/2.0, 0.0)
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0,  ly/2.0, 0.0)
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0, -ly/2.0, 0.0)
		data += "\n"
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0, -ly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n", lx/2.0, -ly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n", lx/2.0,  ly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0,  ly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n",-lx/2.0, -ly/2.0, d)
		data += "\n"
		//column
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0, -coly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n", colx/2.0, -coly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n", colx/2.0,  coly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0,  coly/2.0, d)
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0, -coly/2.0, d)
		data += "\n"
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0, -coly/2.0, d + hf)
		data += fmt.Sprintf("%f %f %f\n", colx/2.0, -coly/2.0, d + hf)
		data += fmt.Sprintf("%f %f %f\n", colx/2.0,  coly/2.0, d + hf)
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0,  coly/2.0, d + hf)
		data += fmt.Sprintf("%f %f %f\n",-colx/2.0, -coly/2.0, d + hf)
		data += "\n"
		//plot connecting lines ?? IT CONNECTS BY ITSELF why also thanks
	}
	data += "\n\n"
	//rebar
	var step, x0, y0, z0 float64
	dia, nx, spcx, ny, spcy := rez[0], rez[1], rez[2], rez[5], rez[6]
	efcvr := dia/2000.0 + nomcvr
	y0 = -ly/2.0 + efcvr
	z0 = nomcvr + dia/2000.0
	n := int(nx)
	step = spcx/1000.0
	x0 = -lx/2.0 + nomcvr + dia/2000.0
	x0 -= step
	for i := 0; i < n; i++{
		x0 += step
		data += fmt.Sprintf("%f %f %f %f %f %f 1\n", x0, y0, z0, 0.0, ly - 2.0 * nomcvr, 0.0)
	}
	x0 = -lx/2.0 + nomcvr
	z0 = nomcvr + dia/1000.0
	n = int(ny)
	step = spcy/1000.0
	y0 = -ly/2.0 + nomcvr + dia/2.0/1000.0
	y0 -= step
	for i := 0; i < n; i++{
		y0 += step
		data += fmt.Sprintf("%f %f %f %f %f %f 2\n", x0, y0, z0, lx - 2.0 * nomcvr, 0.0, 0.0)
	}
	data += "\n\n"
	f, e1 := os.CreateTemp("", "barf")
	if e1 != nil {
		log.Println(e1)
	}	
	defer f.Close()
	defer os.Remove(f.Name())
	_, e1 = f.WriteString(data)
	if e1 != nil {
		log.Println(e1)
	}
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	pltskript := filepath.Join(basepath,"/plotfooting.gp")
	cmd := exec.Command("gnuplot","-c",pltskript,f.Name(),plot)
	//log.Println(pltskript)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outstr, errstr := stdout.String(), stderr.String()
	if err != nil{
		log.Println(err)
		log.Println(errstr)
		return
	}
	if errstr != "" {
		//err = ErrDraw
		log.Println(errstr)
		//xreturn
	}
	
	fmt.Println(outstr)
	//if plot == "dumb" || plot == "caca" {log.Println(outstr)}
	err = nil
	return
}

//plot footing sec plots a footing section view from either dir "x" or "y"
func PlotFtngSec(f *RccFtng, dir string) (data, stdata, adata, ldata, cdata string){
	df := f.Df
	if df == 0.0{df = 1.5/2.0}
	var lx, colx, step float64
	rez := f.Rbr
	dia, nx, spcx, ny, spcy := rez[0], rez[1], rez[2], rez[5], rez[6]
	var n int
	switch dir{
		case "x":
		n = int(ny)
		step = spcy/1000.0
		lx = f.Lx; colx = f.Colx
		case "y":		
		n = int(nx)
		step = spcx/1000.0
		lx = f.Ly; colx = f.Coly
	}
	if f.Sloped{
		//data += fmt.Sprintf("%f %f %f %f\n")	
		var x0, x1, x2, x3, x4, x5, y0, y1, y2 float64
		x1 = lx/2.0 - colx/2.0 - f.Eo
		x2 = lx/2.0 - colx/2.0
		x3 = lx/2.0 + colx/2.0
		x4 = lx/2.0 + colx/2.0 + f.Eo
		x5 = lx
		y1 = f.Dmin
		y2 = f.Dused
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x0, y0, lx, 0.0)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x5, y0, 0.0, f.Dmin)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x0, y0, 0.0, y1-y0)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x1, y2, x2-x1, 0.0)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x3, y2, x4-x3, 0.0)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x2, y2, 0.0, f.Dused + df)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x3, y2, 0.0, f.Dused + df)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x0, y1, x1-x0, y2-y1)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x5, y1, x4-x5, y2-y1)
		
		
		x0 = f.Efcvr
		y0 = f.Efcvr + dia/1000.0		
		if dir == "y"{
			y0 = f.Efcvr
		}
		//stdata += fmt.Sprintf("%f %f %f %f 1.0\n",x0,f.Efcvr,x5-x0-f.Efcvr,0.0)
		adata += fmt.Sprintf("%f %f %f %f 1.0\n",0.0, y0 - 0.2, lx, 0.0)
		adata += fmt.Sprintf("%f %f %f %f 1.0\n",x5 + 0.2, 0.0,0.0, f.Dmin)
		
		adata += fmt.Sprintf("%f %f %f %f 1.0\n",x5 + 0.3, 0.0,0.0, f.Dused)
		adata += fmt.Sprintf("%f %f %f %f 1.0\n",x2, y2 + 0.2 ,x3-x2, 0.0)
		ldata += fmt.Sprintf("%f %f %.3fm 1.0\n",(x2+x3)/2.0, y2 + 0.2, colx)
		ldata += fmt.Sprintf("%f %f %.3fm 1.0\n",x5 + 0.25, f.Dmin/2.0, f.Dmin)
		ldata += fmt.Sprintf("%f %f %.3fm 1.0\n",x5 + 0.5, f.Dused/2.0, f.Dused)
		ldata += fmt.Sprintf("%f %f %.3fm 1.0\n",(x0+x5)/2.0, y0 - 0.25, lx)
		//adata += fmt.Sprintf("%f %f %f %f 1.0\n",x0,f.Efcvr,x5-x0-f.Efcvr,0.0)
		ldata += fmt.Sprintf("%f %f x-#%.f-Y%.f-%.fmm 3.0\n",x5+1.0,f.Efcvr-0.1,nx,dia,spcx)
		ldata += fmt.Sprintf("%f %f y-#%.f-Y%.f-%.fmm 4.0\n",x5+1.0,f.Efcvr-0.2,ny,dia,spcy)

		if dir == "y"{
			stdata += fmt.Sprintf("%f %f %f %f 3.0\n",x0,f.Efcvr+dia/1000.0,step*float64(n-1),0.0)

		} else {
			stdata += fmt.Sprintf("%f %f %f %f 3.0\n",x0,f.Efcvr,step*float64(n-1),0.0)
			
		}
		x0 -= step
		for i := 1; i <= n; i++{
			x0 += step
			cdata += fmt.Sprintf("%f %f %f %f 1\n", x0, y0, dia/1000.0, dia/1000.0)
		}
		data += "\n\n";stdata += "\n\n"; ldata += "\n\n"; adata += "\n\n";cdata += "\n\n"
	} else {
		var x0, x1, x2, x3, y0, y1 float64
		x1 = lx/2.0 - colx/2.0; x2 = lx/2.0 + colx/2.0; x3 = lx
		y1 = f.Dused
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x0, y0, lx, 0.0)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x3, y0, 0.0, f.Dused)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x3, y1, x2-x3, 0.0)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x1, y1, x0-x1, 0.0)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x1, y1, 0.0, f.Dused + df)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x2, y1, 0.0, f.Dused + df)
		data += fmt.Sprintf("%f %f %f %f 1.0\n",x0, y0, 0.0, f.Dused)
		adata += fmt.Sprintf("%f %f %f %f 1.0\n",x0, y0 - 0.2, lx, 0.0)
		adata += fmt.Sprintf("%f %f %f %f 1.0\n",x3 + 0.2, y0 ,0.0, f.Dused)
		adata += fmt.Sprintf("%f %f %f %f 1.0\n",x1, y1 + 0.2 ,x2-x1, 0.0)
		ldata += fmt.Sprintf("%f %f %.3fm 1.0\n",x3 + 0.25, f.Dused/2.0, f.Dused)
		ldata += fmt.Sprintf("%f %f %.3fm 1.0\n",(x0+x3)/2.0, y0 - 0.25, lx)
		ldata += fmt.Sprintf("%f %f %.3fm 1.0\n",(x0+x3)/2.0, y1 + 0.2, colx)
		ldata += fmt.Sprintf("%f %f x-#%.f-Y%.f-%.fmm 3.0\n",x3+1.0,f.Efcvr-0.1,nx,dia,spcx)
		ldata += fmt.Sprintf("%f %f y#%.f-Y%.f-%.fmm 4.0\n",x3+1.0,f.Efcvr-0.2,ny,dia,spcy)

		x0 = f.Efcvr
		y0 = f.Efcvr + dia/1000.0

		if dir == "y"{
			y0 = f.Efcvr
		}
		
		if dir == "y"{
			stdata += fmt.Sprintf("%f %f %f %f 3.0\n",x0,f.Efcvr+dia/1000.0,step*float64(n-1),0.0)
		} else {
			stdata += fmt.Sprintf("%f %f %f %f 1.0\n",x0,f.Efcvr,step*float64(n-1),0.0)
		}
		
		//stdata += fmt.Sprintf("%f %f %f %f 1.0\n",x0,f.Efcvr,x3-x0-f.Efcvr,0.0)
		
		//adata += fmt.Sprintf("%f %f %f %f 1.0\n",x0,f.Efcvr,x3-x0-f.Efcvr,0.0)
		x0 -= step
		for i := 1; i <= n; i++{
			x0 += step
			cdata += fmt.Sprintf("%f %f %f %f 1\n", x0, y0, dia/1000.0, dia/1000.0)
		}
		data += "\n\n";stdata += "\n\n"; ldata += "\n\n"; adata += "\n\n";cdata += "\n\n"	
	}
	return
}


func PlotFtngPlan(f *RccFtng) (data, stdata, adata, ldata string){
	//data += fmt.Sprintf("")
	var x0, y0, x1, y1 float64
	x1 = f.Ly; y1 = f.Lx
	data += fmt.Sprintf("%f %f %f %f 1.0\n",x0,y0,f.Ly,0.0)
	data += fmt.Sprintf("%f %f %f %f 1.0\n",x0,y0,0.0,f.Lx)
	data += fmt.Sprintf("%f %f %f %f 1.0\n",x1,y1,-f.Ly,0.0)
	data += fmt.Sprintf("%f %f %f %f 1.0\n",x1,y1,0.0,-f.Lx)

	
	data += fmt.Sprintf("%f %f %f %f 2.0\n",(f.Ly-f.Coly)/2.0,(f.Lx-f.Colx)/2.0,f.Coly,0.0)
	data += fmt.Sprintf("%f %f %f %f 2.0\n",(f.Ly-f.Coly)/2.0,(f.Lx-f.Colx)/2.0,0.0,f.Colx)
	data += fmt.Sprintf("%f %f %f %f 2.0\n",(f.Ly+f.Coly)/2.0,(f.Lx+f.Colx)/2.0,-f.Coly,0.0)
	data += fmt.Sprintf("%f %f %f %f 2.0\n",(f.Ly+f.Coly)/2.0,(f.Lx+f.Colx)/2.0,0.0,-f.Colx)
	
	
	
	adata += fmt.Sprintf("%f %f %f %f 1.0\n",x0,y0-0.2,f.Ly,0.0)
	adata += fmt.Sprintf("%f %f %f %f 1.0\n",x0-0.2,y0,0.0,f.Lx)

	
	var step float64
	rez := f.Rbr
	dia, nx, spcx, ny, spcy := rez[0], rez[1], rez[2], rez[5], rez[6]
	xs := f.Efcvr; ys := f.Efcvr
	step = spcx/1000.0
	xs -= step
	for i := 1; i <= int(nx); i++{
		xs += step
		stdata += fmt.Sprintf("%f %f %f %f 3.0\n",xs, ys, 0.0, f.Lx - 2.0 * f.Efcvr)
	}
	xs = f.Efcvr; ys = f.Efcvr
	step = spcy/1000.0
	ys -= step
	for i := 1; i <= int(ny); i++{
		ys += step
		stdata += fmt.Sprintf("%f %f %f %f 4.0\n",xs, ys, f.Ly - 2.0 * f.Efcvr, 0.0)
	}
	
	ldata += fmt.Sprintf("%f %f %f %f 1.0\n",(x0+f.Ly)/2.0,y0-0.2,f.Ly,1.0)
	ldata += fmt.Sprintf("%f %f %f %f 1.0\n",x0-0.2,(y0+f.Lx)/2.0,f.Lx,0.0)
	ldata += fmt.Sprintf("%f %f x-#%.f-Y%.f-%.fmm 3.0\n",x1+1.0,f.Efcvr-0.1,nx,dia,spcx)
	ldata += fmt.Sprintf("%f %f y#%.f-Y%.f-%.fmm 4.0\n",x1+1.0,f.Efcvr-0.2,ny,dia,spcy)
	
	data += "\n\n";stdata += "\n\n"; ldata += "\n\n"; adata += "\n\n"
	return
}

//PlotFtngDet plots rcc footing detail views 
func PlotFtngDet(f *RccFtng) (pltstr string){
	//colx, coly, fck, fy, hf, eo, d, dmin, lx, ly, nomcvr float64, sloped bool, rez []float64, plot string
	data, stdata, adata, ldata, cdata := PlotFtngSec(f, "x")
	data += stdata; data += adata; data += ldata; data += cdata
	//fmt.Println(data)
	
	d1, st1, a1, l1 := PlotFtngPlan(f)
	data += d1; data += st1; data += a1; data += l1
	fn := fmt.Sprintf("%s-%s.svg",f.Title,"detail")
	if f.Term == "dxf"{
		fn = fmt.Sprintf("%s-%s.dxf",f.Title,"detail")
	} 
	fname := genfname("",fn)
	pltskript := "plotftngdet.gp"
	pltstr = skriptrun(data, pltskript, f.Term, f.Title, fname)
	//err = fmt.Errorf("%s",pltstr)
	if f.Term == "svg" || f.Term == "dxf" || f.Term == "svgmono"{
		pltstr = fname
	}
	if f.Web{
		//embed kongtext font in svg
		//kass.Svgkong(pltstr)
		if f.Term != "dxf"{kass.Svgkong(pltstr)}
		f.Txtplots = append(f.Txtplots, fn)
	}
	//if plot == "dumb" || plot == "caca" {log.Println(outstr)}
	//err = nil
	return
}

//PlotSubFrm plots a sub frame line diagram
func PlotSubFrm(sf *SubFrm, mod *kass.Model, ms map[int]*kass.Mem, bmvec, colvec []int, bmenv map[int]*kass.BmEnv, colenv map[int]*kass.ColEnv, term string) (pltstr string){
	var data string
	//index 0 nodes
	var frcscale, xmax, ymax, xmin, ymin float64
	for idx, v := range mod.Coords{
		data += fmt.Sprintf("%v %v %v\n", v[0], v[1], idx+1)
		if v[1] > ymax {ymax = v[1]}
		if v[0] > xmax {xmax = v[0]}
		if v[1] < ymin {ymin = v[1]}
		if v[0] < xmin {xmin = v[0]}
	}
	data += "\n\n"
	//index 1 members
	for idx, mem := range mod.Mprp{
		jb := mod.Coords[mem[0]-1]
		je := mod.Coords[mem[1]-1]
		data += fmt.Sprintf("%v %v %v %v (%v) (%.fx%.f)\n", jb[0], jb[1], je[0], je[1], idx+1, sf.Sections[mem[3]-1][0], sf.Sections[mem[3]-1][1])
	}
	data += "\n\n"
	//index 2 supports
	for _, val := range mod.Supports{
		pt := mod.Coords[val[0]-1]
		data += fmt.Sprintf("%v %v %v\n", pt[0],pt[1],-int(val[1]+val[2]+val[3]))
	}
	data += "\n\n"
	//index 3 udl
	for mem, ldcases := range sf.Advloads {
		for _, val := range ldcases{
			switch sf.Members[mem][1][0]{
				case 2,3,4:
				//all beamz
				jb := mod.Mprp[mem-1][0]; xb := mod.Coords[jb-1][0]; yb := mod.Coords[jb-1][1]
				lspan := ms[mem].Geoms[0]
				//ye := val[2] //+ yb
				ye := sf.Hs[0] * (val[2])/(sf.DL*sf.PSFs[0])/40.0 + yb
				xb = xb + val[4]; xe := xb + lspan - val[5]
				data += fmt.Sprintf("%f %f %f %f %f %f %f\n",xb, ye, xe, ye, val[2], val[1], val[6])
			}
		}
	}
	data += "\n\n"
	//index 4 joint loads
	for _, val := range mod.Jloads{
		frcscale = frcscale+(ymin-xmin)/(ymax-xmax) //just to AVOID UNUSED VARIABLE 
		pt := mod.Coords[int(val[0])-1]
		if val[1] != 0.0 { //X- force (assemble?)
			//frcscale = 10.0//(val[1] - xmin)*(ymax - ymin)/(xmax - xmin) //horrible way to scale
			if pt[0] == xmax {
				//vector to the right
				data += fmt.Sprintf("%v %v %v %v %.1f\n",pt[0],pt[1],frcscale, 0, val[1])
			} else {
				data += fmt.Sprintf("%v %v %v %v %.1f\n",pt[0],pt[1],-frcscale, 0, val[1])
			}
		}
		if val[2] != 0.0 { //y force
			//frcscale = 10.0//(val[2] - ymin)*(xmax - xmin)/(ymax - ymin) 
			if pt[1] == ymax {
				data += fmt.Sprintf("%v %v %v %v %.1f\n",pt[0],pt[1],0, frcscale, val[2])
			} else {
				data += fmt.Sprintf("%v %v %v %v %.1f\n",pt[0],pt[1]-frcscale,0, frcscale*val[2], val[2])
			}
		}
		if val[3] != 0.0 {//moment z
			//frcscale = 10.0//(val[3] - ymin)*(xmax - xmin)/(ymax - ymin) 
			data += fmt.Sprintf("%v %v %v %v %.1fk-i\n",pt[0],pt[1],frcscale, frcscale, val[3])
		}
	}
	data += "\n\n"
	fn := fmt.Sprintf("%s_%s.svg",sf.Title,"subfrm")
	if term == "dxf"{fn = fmt.Sprintf("%s_%s.dxf",sf.Title,"subfrm")}
	fname := genfname(sf.Foldr,fn)
	//fmt.Println(fname)
	pltskript := "plotsubfrm.gp"
	outstr := skriptrun(data, pltskript, term, fn, fname)
	if term == "svg" || term == "svgmono" || term == "dxf"{
		outstr = fname
	}
	if sf.Web{
		if term != "dxf"{kass.Svgkong(outstr)}
		//sf.Txtplots = append(sf.Txtplots, fn)
		outstr = fn
	}
	return outstr
}

//dumbplt is an old plotting function
//god knows where it lurks
func dumbplt(data, skript, term string)(txtplot string){
	//create temp files
	f, e1 := os.CreateTemp("", "barf")
	if e1 != nil {
		log.Println(e1)
	}
	defer f.Close()
	defer os.Remove(f.Name())
	_, e1 = f.WriteString(data)
	if e1 != nil {
		log.Println(e1)
	}
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	pltskript := filepath.Join(basepath,skript)
	cmd := exec.Command("gnuplot","-c",pltskript,f.Name(),term)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outstr, errstr := stdout.String(), stderr.String()
	if err != nil {
		log.Println(err)
	}
	if errstr != "" {
		log.Println(errstr)
	}
	if term == "dumb" || term == "caca" || term == "svg"{
		txtplot = fmt.Sprintf("%v",outstr)
	}
	return
}


//PlotSlb is woefully incomplete, plots a slab
func PlotSlb(s *RccSlb, folder, term string) (txtplot string){
	var data, ldata string
	//draw plan, supports, main steel, dist
	switch s.Type{
		case 1:
		//une way slab
		case 2:
		fmt.Println("plsholdr",data, ldata)
	}
	return
	
}
