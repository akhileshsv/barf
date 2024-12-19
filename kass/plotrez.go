package barf

import (
	"fmt"
	"math/rand"
)


//PlotRez2d plots a 2d frame/beam model with bm/sf/def. results
func PlotRez2d(mod *Model, term string)(txtplot string){
	var data, rdata, cdata, ldata string
	//translate/rotate coords by nodal displ. - FORGET THIS FOR NOW
	var tx, ty, ang float64
	for i, node := range mod.Js{
		switch mod.Frmstr{
			case "2dt":
			tx, ty = node.Displ[0], node.Displ[1]
			case "1db":
			ty, ang = node.Displ[0], node.Displ[1]
			case "2df":
			tx, ty, ang = node.Displ[0], node.Displ[1], node.Displ[2]
		}
		ty = 10.0 * ty
		// tx = 10.0 * tx
		p0 := node.Coords
		if len(p0) == 1{p0 = append(p0,0.0)}
		p1 := Trans2d(p0, tx, ty)
		var p2 []float64
		if ang == 0.0{
			p2 = p1
		} else {
			p2 = RotZ(p1, ang)
		}
		mod.Js[i].Dcs = make([]float64, len(p2))
		copy(mod.Js[i].Dcs,p2)
	}
	for i, m := range mod.Mprp {
		//slices for filled curves
		jb := mod.Coords[m[0]-1]
		je := mod.Coords[m[1]-1]
		if len(jb) == 1{
			jb = append(jb, 0.0)
			je = append(je, 0.0)
		}
		p3 := mod.Js[m[0]].Dcs
		p4 := mod.Js[m[1]].Dcs
		lmem := Dist3d(jb, je)
		l2 := Dist3d(p3, p4)
		mem := mod.Ms[i+1]
		data += fmt.Sprintf("%f %f %f %f %v\n",jb[0], jb[1], je[0], je[1], i)
		//cdata += fmt.Sprintf("%f %f %f %f %v\n",jbn[0], jbn[1], jen[0], jen[1], i)
		ldata += fmt.Sprintf("%f %f %.1fKN %.1fKN.M %.3f(M)\n",(jb[0]+je[0])/2.0, (jb[1]+je[1])/2.0,mem.Rez.Maxs[0], mem.Rez.Maxs[1],mem.Rez.Maxs[2])
		//cdata += fmt.Sprintf("%f %f %f %f %v\n",jbn[0], jbn[1], jen[0]-jbn[0], jen[1]-jbn[1], i)
		
		for j, x := range mem.Rez.Xs{
			vu, bm, dx := mem.Rez.SF[j] * mod.Scales[0], mem.Rez.BM[j]* mod.Scales[1], mem.Rez.Dxs[j] * 10.0
			if j == 20{
				p1 := je
				p2 := Rotvec(90,je,jb)
				pv := Lerpvec(-vu/lmem, p1, p2)
				pb := Lerpvec(-bm/lmem, p1, p2)
				//ye olde
				// pd := Lerpvec(dx/lmem, p1, p2)
				
				//ye new
				p5 := p4
				p6 := Rotvec(90, p4, p3)
				pd := Lerpvec(dx/l2, p5, p6)
				//%f %f %f %f %f %f
				rdata += fmt.Sprintf("%f %f %f %f %f %f %f %f\n",pv[0], pv[1],pb[0], pb[1],pd[0], pd[1], p1[0], p1[1])
				continue
			}
			p1 := Lerpvec(x/lmem, jb, je)
			d0 := Dist3d(p1, je)
			p2 := Rotvec(90,p1,je)
			pv := Lerpvec(vu/d0, p1, p2)
			pb := Lerpvec(bm/d0, p1, p2)

			p5 := Lerpvec(x/l2, p3, p4)
			p6 := Rotvec(90, p5, p4)
			pd := Lerpvec(-dx/Dist3d(p5,p6),p5,p6)
			//ye olde
			// pd := Lerpvec(-dx/d0, p1, p2)
			rdata += fmt.Sprintf("%f %f %f %f %f %f %f %f\n",pv[0], pv[1],pb[0], pb[1],pd[0], pd[1], p1[0], p1[1])
		}
		rdata += "\n"	
	}
	rdata += "\n\n"; ldata += "\n\n"
	data += "\n\n"; data += rdata; data += ldata; data += cdata
	//data += ldata
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("results-%s",mod.Id)
	title := fmt.Sprintf("results-%s",mod.Id)
	//fname = svgpath(mod.Foldr,fname,term)
	skript := "drawmodrez.gp"
	txtplot = skriptrun(data, skript, term, title,mod.Foldr,fname)
	if term == "dumb" || term == "mono"{fmt.Println(txtplot)}
	if mod.Web{
		switch term{
			case "dxf":		
			txtplot = fname + ".dxf"
			case "svg", "svgmono":
			Svgkong(txtplot)
			txtplot = fname + ".svg"
		}
	}
	// log.Println("TXTPLAAAT->",txtplot)
	return
}

//PlotNpRez2d plots a 2d np frame/beam model with bm/sf/def. results
func PlotNpRez2d(mod *Model, term string)(txtplot string){
	var data, rdata, cdata, ldata string
	//translate/rotate coords by nodal displ. - FORGET THIS FOR NOW
	var tx, ty, ang float64
	for i, node := range mod.Js{
		switch mod.Frmstr{
			case "1db":
			ty, ang = node.Displ[0], node.Displ[1]
			case "2df":
			tx, ty, ang = node.Displ[0], node.Displ[1], node.Displ[2]
		}
		ty = 10.0 * ty
		// tx = 10.0 * tx
		p0 := node.Coords
		if len(p0) == 1{p0 = append(p0,0.0)}
		p1 := Trans2d(p0, tx, ty)
		var p2 []float64
		if ang == 0.0{
			p2 = p1
		} else {
			p2 = RotZ(p1, ang)
		}
		mod.Js[i].Dcs = make([]float64, len(p2))
		copy(mod.Js[i].Dcs,p2)
	}
	for i, m := range mod.Mprp {
		jb := mod.Coords[m[0]-1]
		je := mod.Coords[m[1]-1]
		if len(jb) == 1{
			jb = append(jb, 0.0)
			je = append(je, 0.0)
		}
		p3 := mod.Js[m[0]].Dcs
		p4 := mod.Js[m[1]].Dcs
		lmem := Dist3d(jb, je)
		l2 := Dist3d(p3, p4)
		mem := mod.Mnps[i+1]
		data += fmt.Sprintf("%f %f %f %f %v\n",jb[0], jb[1], je[0], je[1], i)
		//cdata += fmt.Sprintf("%f %f %f %f %v\n",jbn[0], jbn[1], jen[0], jen[1], i)
		ldata += fmt.Sprintf("%f %f %.1fKN %.1fKN.M %.3f(M)\n",(jb[0]+je[0])/2.0, (jb[1]+je[1])/2.0,mem.Rez.Maxs[0], mem.Rez.Maxs[1],mem.Rez.Maxs[2])
		//cdata += fmt.Sprintf("%f %f %f %f %v\n",jbn[0], jbn[1], jen[0]-jbn[0], jen[1]-jbn[1], i)
		for j, x := range mem.Rez.Xs{
			vu, bm, dx := mem.Rez.SF[j] * mod.Scales[0], mem.Rez.BM[j]* mod.Scales[1], mem.Rez.Dxs[j] * 10.0
			if j == 20{
				p1 := je
				p2 := Rotvec(90,je,jb)
				pv := Lerpvec(-vu/lmem, p1, p2)
				pb := Lerpvec(-bm/lmem, p1, p2)
				//ye olde
				// pd := Lerpvec(dx/lmem, p1, p2)
				
				//ye new
				p5 := p4
				p6 := Rotvec(90, p4, p3)
				pd := Lerpvec(dx/l2, p5, p6)
				rdata += fmt.Sprintf("%f %f %f %f %f %f\n",pv[0], pv[1],pb[0], pb[1],pd[0], pd[1])
				continue
			}
			p1 := Lerpvec(x/lmem, jb, je)
			d0 := Dist3d(p1, je)
			p2 := Rotvec(90,p1,je)
			pv := Lerpvec(vu/d0, p1, p2)
			pb := Lerpvec(bm/d0, p1, p2)

			p5 := Lerpvec(x/l2, p3, p4)
			p6 := Rotvec(90, p5, p4)
			pd := Lerpvec(-dx/Dist3d(p5,p6),p5,p6)
			//ye olde
			// pd := Lerpvec(-dx/d0, p1, p2)
			rdata += fmt.Sprintf("%f %f %f %f %f %f\n",pv[0], pv[1],pb[0], pb[1],pd[0], pd[1])
		}
		rdata += "\n"
	}
	rdata += "\n\n"; ldata += "\n\n"
	data += "\n\n"; data += rdata; data += ldata; data += cdata
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("results-%s",mod.Id)
	title := fmt.Sprintf("results-%s",mod.Id)
	//fname = svgpath(mod.Foldr,fname,term)
	skript := "drawmodrez.gp"
	txtplot = skriptrun(data, skript, term, title,mod.Foldr,fname)
	if term == "dumb" || term == "mono"{fmt.Println(txtplot)}
	if mod.Web{
		switch term{
			case "dxf":		
			txtplot = fname + ".dxf"
			case "svg", "svgmono":
			Svgkong(txtplot)
			txtplot = fname + ".svg"
		}
	}
	return
}

/*
   ye olde
   line 83
   
		var lpts []float64
		for k, val := range mem.Rez.Maxs{
			xloc := mem.Rez.Locs[k]
			p1 := Lerpvec(xloc/lmem,jb, je)
			p2 := Rotvec(90, p1, je)
			d0 := Dist3d(p1, je)
			if k == 3 || k == 2{val = -val}
			var scale float64
			switch k{
				case 0:
				scale = mod.Scales[0]
				case 1,3:
				scale = mod.Scales[1]
				case 2:
				scale = 20.0
			}
			pval := Lerpvec(val*scale/d0, p1, p2)
			lpts = append(lpts, pval[0],pval[1],val)
		}
		var row string
		for _, val := range lpts{
			row += fmt.Sprintf("%f ",val)
		}
		row += "\n"
		ldata += row		
		//rez.Maxs = []float64{vmax, mmax, dmax, mpmax}
		//rez.Locs = []float64{vmaxx, mmaxx, dmaxx, mpmaxx}
	

*/
