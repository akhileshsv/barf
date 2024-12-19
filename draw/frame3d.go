package barf

// import (
// 	"fmt"
// 	//"math"
// )

// //FrmGeom plots frame geometry and plan from frm3dgen
// func FrmGeom(coords [][]float64, sups map[int][]int,mems map[string][][]int,
// 	sns [][]int, scat[][]int, ncols, nbms int, term, folder string) (txtplot string, err error){
// 	var data, ldata string
// 	//first coords
// 	for idx, v := range coords {
// 		data += fmt.Sprintf("%f %f %f\n", v[0], v[1], v[2])
// 		ldata += fmt.Sprintf("%f %f %f N(%v)\n", v[0], v[1], v[2], idx+1)
// 	}
// 	data += "\n\n"
// 	//then supports
// 	for s, val := range sups {
// 		pt := coords[s-1]
// 		if val[0]+val[1]+val[2]+val[3]+val[4]+val[5] != 0 {data += fmt.Sprintf("%f %f %f\n", pt[0],pt[1],pt[2])}
// 	}
// 	data += "\n\n"
// 	//then members
// 	//m[1] - xdx,ydx,flrdx,locdx,mtyp, locx, lex
// 	for mem, m := range mems{
// 		jb := coords[m[0][0]-1]
// 		je := coords[m[0][1]-1]
// 		mtyp := m[1][4]
// 		dtyp := m[0][2]
// 		lex := m[1][6]
// 		mcolor := dtyp; clf := ncols
// 		switch dtyp{
// 			case 2, 3:
// 			clf = nbms
// 		}
// 		switch clf{
// 			case 1:
// 			mcolor = lex
// 			case 2:
// 			mcolor = mtyp
// 		}
// 		data += fmt.Sprintf("%f %f %f %f %f %f %v %v\n",jb[0],jb[1],jb[2],je[0],je[1],je[2],mcolor, dtyp)
// 		ldata +=  fmt.Sprintf("%f %f %f %s\n",(jb[0]+je[0])/2.0,(jb[1]+je[1])/2.0,(jb[2]+je[2])/2.0,mem)
// 	}
// 	data += "\n\n"
// 	//slabs
// 	for i, sn := range sns{
// 		p1 := coords[sn[0]-1]; p2 := coords[sn[1]-1]; p3 := coords[sn[2]-1]; p4 := coords[sn[3]-1]
// 		//fdx := sn[4]
// 		//ldata += fmt.Sprintf()
// 		var xsum, ysum, zsum float64
// 		for _, v := range [][]float64{p1, p2, p3, p4, p1}{
// 			data += fmt.Sprintf("%f %f %f %v\n", v[0], v[1], v[2],i+1)
// 			xsum += v[0]; ysum += v[1]; zsum += v[2]
// 		}
// 		xsum = xsum/5; ysum = ysum/5; zsum = zsum/5
// 		data += "\n"
// 		lbl := fmt.Sprintf("%v,%v,%v",i,scat[i][0],scat[i][1])
// 		ldata += fmt.Sprintf("%f %f %f %s 3\n",xsum, ysum, zsum, lbl)
// 	}
// 	data += "\n\n"
// 	//text labels
// 	data += ldata
// 	skript := "frm3d.gp"
// 	//folder += "/framegnu.svg"
// 	Draw(fdat, skript, term, folder, title, title, "gen", "weight","")
// 	txtplot, err = Draw(data, skript, term, folder, "frm3d.svg","")
// 	//fmt.Println(txtplot,"\n****\n")
// 	//txtplot, err = Dumb(data, "dumbfrm3d.gp", "dumb", "FRAME", "", "", "")
// 	//fmt.Println(txtplot)
// 	//txtplot, err = Dumb(data, skript, term, title, "", "", "")
// 	return 
// } 

// //FrmLoads plots frame loads from frm3dgen
// func FrmLoads(coords [][]float64, sups map[int][]int,mems map[string][][]int,
// 	mloads map[string][][]float64,
// 	mverts map[string]map[int][][]float64,
// 	mareas map[string][]float64, term string, folder string) (txtplot string, err error){
// 	//SAME AS ABOVE EXCEPT FOR MEM VERTS
// 	var data, ldata string
// 	//first coords
// 	for idx, v := range coords {
// 		data += fmt.Sprintf("%f %f %f\n", v[0], v[1], v[2])
// 		ldata += fmt.Sprintf("%f %f %f N(%v)\n", v[0], v[1], v[2], idx+1)
// 	}
// 	data += "\n\n"
// 	//then supports
// 	for s, val := range sups {
// 		pt := coords[s-1]
// 		if val[0]+val[1]+val[2]+val[3]+val[4]+val[5] != 0 {data += fmt.Sprintf("%f %f %f\n", pt[0],pt[1],pt[2])}
// 	}
// 	data += "\n\n"
// 	//then members
// 	for mem, m := range mems{
// 		jb := coords[m[0][0]-1]
// 		je := coords[m[0][1]-1]
// 		mtyp := m[1][4]
// 		dtyp := m[0][2]
// 		data += fmt.Sprintf("%f %f %f %f %f %f %v %v\n",jb[0],jb[1],jb[2],je[0],je[1],je[2],mtyp,dtyp)
// 		ldata +=  fmt.Sprintf("%f %f %f %s\n",(jb[0]+je[0])/2.0,(jb[1]+je[1])/2.0,(jb[2]+je[2])/2.0,mem)
// 	}
// 	data += "\n\n"
// 	//mem trib areas
// 	for _, mvs := range mverts{
// 		for slab, mv := range mvs{
// 			var xsum, ysum, zsum, nv float64
// 			v0 := make([]float64,3)
// 			for k, v := range mv{
// 				nv += 1.0
// 				xsum += v[0];  ysum += v[1]; zsum += v[2]
// 				if k == 0{v0[0] = v[0]; v0[1] = v[1]; v0[2] = v[2]}
// 				data += fmt.Sprintf("%f %f %f %v\n", v[0], v[1], v[2], slab)
// 			}
// 			data += fmt.Sprintf("%f %f %f %v\n", v0[0], v0[1], v0[2], slab)
// 			data += "\n"
// 			ldata += fmt.Sprintf("%f %f %f %v\n", xsum/nv, ysum/nv, zsum/nv, slab)
// 		}
// 	}
// 	data += "\n\n"
// 	//text labels
// 	data += ldata
// 	skript := "frm3d.gp"
// 	//folder += "/framegnu.svg"
// 	txtplot, err = Draw(data, skript, term, folder, "frmloads.svg","")
// 	//txtplot, err = Dumb(data, skript, term, title, "", "", "")
// 	return 

// }

