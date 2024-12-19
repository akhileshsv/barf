package barf

import (
	kass"barf/kass"
	//draw"barf/draw"
)

//DELETE THIS
//worthless this func is

func PlotFrm2d(f *kass.Frm2d, term string, plotchn chan string){
	txtplot := kass.DrawMod2d(&f.Mod, f.Ms, term)
	plotchn <- txtplot
}
/*
func PlotFrm2dLp(f *Frm2d, lp int, term string){
	var data, mdata, ldata string
	var fsc, xmin, ymin, xmax, ymax float64
	fsc = 1e-2
	//index 0 nodes
	for idx, v := range f.Mod.Coords {
		data += fmt.Sprintf("%.1f %.1f %v\n", v[0], v[1], idx+1)
		if v[1] > ymax {ymax = v[1]}
		if v[0] > xmax {xmax = v[0]}
		if v[1] < ymin {ymin = v[1]}
		if v[0] < xmin {xmin = v[0]}
	}
	data += "\n\n"
	//index 1 members
	//ms := make(map[int][]int)
	for idx, mem := range f.Mod.Mprp {
		jb := f.Mod.Coords[mem[0]-1]
		je := f.Mod.Coords[mem[1]-1]
		data += fmt.Sprintf("%.1f %.1f %.1f %.1f %v\n", jb[0], jb[1], je[0], je[1], idx+1)
		ldata += fmt.Sprintf("%f %f %v\n",(jb[0]+je[0])/2.0, (jb[1]+je[1])/2.0, idx+1)
	}
	data += "\n\n"
	//index 2 supports
	for _, val := range f.Mod.Supports {
		pt := f.Mod.Coords[val[0]-1]
		if val[1]+val[2] != 0 {data += fmt.Sprintf("%.1f %.1f\n", pt[0],pt[1])}
	}
	data += "\n\n"
	//index 3 joint loads
	for _, val := range f.Mod.Jloads{
		pt := f.Mod.Coords[int(val[0])-1]
		if val[1] != 0.0 { 
			if pt[0] == xmax {
				data += fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f\n",pt[0],pt[1],fsc*val[1], 0.0, val[1])
			} else {
				data += fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f\n",pt[0],pt[1],-fsc*val[1], 0.0, val[1])
			}
		}
		if val[2] != 0.0 {
			if pt[1] == ymax {
				data += fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f\n",pt[0],pt[1],0.0, -fsc*val[2], val[2])
			} else {
				data += fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f\n",pt[0],pt[1],0.0, fsc*val[2], val[2])
			}
		}
		if val[3] != 0.0 {
			mdata += fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f\n",pt[0],pt[1],fsc, fsc, val[3])
		}
	}
	data += "\n\n"
	//index 4 member loads
	for _, val := range f.Mod.Msloads {
		m := int(val[0])
		mem := ms[m]
		jb := f.Mod.Mprp[m-1][0]
		xa, ya := f.Mod.Coords[jb-1][0], f.Mod.Coords[jb-1][1]
		ltyp := int(val[1])
		wa, wb, la, lb := val[2], val[3], val[4], val[5]
		cx := mem.Geoms[4]; cy := mem.Geoms[5]
		ldata += fmt.Sprintf("%f %f %.0f\n",xa+la*cx, ya+la*cy, wa)
		ya += fsc*5.0
		switch ltyp {
		case 1://point load at la
			data += fmt.Sprintf("%f %f %f %f %v\n",xa + la * cx, ya + la * cy, fsc*wa*cy, -fsc*wa*cx, ltyp)
		case 2:
			//moment at la 
			mdata += fmt.Sprintf("%f %f %f\n",xa + la * cx, ya + la * cy, fsc*wa)
		case 3://udl w from la to l - lb
			l := mem.Geoms[0]
			div := (l-lb-la)/5.0
			xa -= div * cx; ya -= div * cy
			for i:=0; i < 5; i++{
				xa += div * cx; ya += div * cy
				data += fmt.Sprintf("%f %f %f %f %.0f %v\n",xa,ya,wa*fsc*cy,wa*fsc*cx,wa, ltyp)
			}
		case 4://udl wa at la to wb at l - lb
			l := mem.Geoms[0]
			div := (l-lb-la)/5.0
			dw := (wb - wa)*fsc/5.0
			xa -= div * cx; ya -= div * cy
			for i:=0; i < 5; i++{
				xa += div * cx; ya += div * cy
				data += fmt.Sprintf("%f %f %f %f %v\n",xa,ya,dw*cy,-dw*cx, ltyp)
			}
		case 5:
			//point axial load at la
			data += fmt.Sprintf("%f %f %f %f %v\n",xa+fsc,ya+fsc,wa*cy*fsc,-wa*cx*fsc, ltyp)
		case 6:
			//uniform axial load w at la to l - lb
			l := mem.Geoms[0]
			div := (l-lb-la)/2.0
			//xa -= div * cx; ya -= div * cy
			for i:=0; i < 2; i++{
				xa += div * cx; ya += div * cy
				data += fmt.Sprintf("%f %f %f %f %v\n",xa+fsc*2.,ya+fsc*2.,wa*fsc*cx,wa*fsc*cy, ltyp)
			}
		case 7:
		}
	}
	data += "\n\n"; ldata += "\n\n"
	data += ldata; data += mdata
	skript := "drawmod2d.gp"
	txtplot := dumbplt(data, skript, term)
	return txtplot
        }
*/
