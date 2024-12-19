package barf

//model plotting funcs via gnuplot

import (
	"net/http"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"runtime"
	"path/filepath"
	"math"
	"math/rand"
	"io/ioutil"
	"strconv"
	"strings"
	//"html/template"
	"encoding/base64"
	//draw"barf/draw" FORGET DRAW MAN
	//"gonum.org/v1/gonum/mat"
)

//slerp for sigmoid interpolation of a value
//l - max value, k - steepness , x0 - midpoint, x - is just "x"
func slerp(l, k, x0, x float64) (sval float64){
	return l/(1.0 - math.Exp(k *x0 - k*x))
}

//skriptpath returns the absolute path of a (gnuplot) script in the current folder
func skriptpath(skript string) (string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, skript)
}

//svgpath returns the absolute path for saving a plot as .svg (gnuplot)
func svgpath(foldr, fname, term string) (string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	
	if term == "svg" || term == "svgmono"{
		fname = fname + ".svg"
	}
	if term == "dxf"{
		fname = fname + ".dxf"
	}
	switch foldr{
		// case "web":
		// foldr = filepath.Join(basepath,"../data/web")
		case "","web","out":
		foldr = filepath.Join(basepath,"../data/out")
		default:
		foldr = filepath.Join(basepath, "../data/out", foldr)	
		if _, err := os.Stat(foldr); os.IsNotExist(err){
			var dirMod uint64
			if dirMod, err = strconv.ParseUint("0775", 8, 32); err == nil {
				err = os.Mkdir(foldr, os.FileMode(dirMod))
			}
			if err != nil{
				log.Println("error creating folder-",foldr,"error-",err)
				return ""
			}
		}
	}
	
	return filepath.Join(foldr,fname)
}

//skriptrun runs a gnuplot script and returns the path of svg/text plot/error string
func skriptrun(data, pltskript, term, title, foldr, fname string) (string){
	pltskript = skriptpath(pltskript)
	f, e1 := os.CreateTemp("", "kass")
	if e1 != nil {
		fmt.Println(e1)
	}
	defer f.Close()
	defer os.Remove(f.Name())	
	_, e1 = f.WriteString(data)
	if e1 != nil {
		fmt.Println(e1)
	}
	if fname == ""{
		fname = fmt.Sprintf("file_%v",rand.Intn(666))
	}
	fname = svgpath(foldr,fname, term)
	if fname == ""{
		return "error in folder creation"
	}
	//log.Println("HERE->",fname)
	cmd := exec.Command("gnuplot","-c",pltskript,f.Name(),term, title, fname)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	_ = cmd.Run()	
	// if err != nil {
	// 	//fmt.Println(err)
	// }
	outstr, errstr := stdout.String(), stderr.String()
	if term == "qt"{outstr = outstr + errstr}
	//log.Println("ran skript->",outstr)
	if term == "svg" || term == "dxf" || term == "svgmono"{
		outstr = fname
	}
	return outstr

}

//Svgkong reads an svg embeds kongtext in style defs and saves
//peak noob
func Svgkong(pltstr string){
	bytes, err := ioutil.ReadFile(pltstr)
	if err != nil {
		log.Println(err)
	}
	rez := string(bytes)
	//log.Println(rez)
	d , err := os.ReadFile(skriptpath("../data/out/kongdef.txt"))
	if err != nil{
		log.Fatal(err)
	}
	defs := string(d)
	title := "<title>Gnuplot</title>"
	ndefs := defs + title
	rez = strings.Replace(rez,title,ndefs,1)
	err = os.WriteFile(pltstr,[]byte(rez),0666)
	if err != nil{
		log.Fatal(err)
	}
}

//ToBase64 converts an image to a base64 string
func ToBase64(pltstr string) (rez string) {
	bytes, err := ioutil.ReadFile(pltstr)
	if err != nil {
		log.Println(err)
	}

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/svg":
		rez += "data:image/svg;base64,"
	case "image/png":
		rez += "data:image/png;base64,"
	default:
		rez += "data:image/svg;base64,"
	}
	rez += base64.StdEncoding.EncodeToString(bytes)
	return
}

//PlotFrm2d plots a 2d frame model (line plot)
func PlotFrm2d(mod *Model, term string, pltchn chan string) {
	data, ldata, mdata := GeomDat(mod)
	d1, l1, m1 := JloadDat(mod)
	data += d1; ldata += l1; mdata += m1
	d1, l1, m1 = MsLoad2Dat(mod)
	data += d1; ldata += l1; mdata += m1
	
	if mod.Drawsec{
		d1, l1, _ = View2dDat(mod)
		data += d1
		ldata += l1
	}
	data += ldata
	data += "\n\n"; data += mdata
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	title := fmt.Sprintf("%s-2d-frm",mod.Id)
	fname := title
	// switch term{
	// 	case "dxf":
	// 	fname = fname + ".dxf"
	// 	case "svg", "svgmono":
	// 	fname = fname + ".svg"
	// } 
	//fname = svgpath(mod.Foldr,fname, term)
	skript := "drawmod2d.gp"
	txtplot := skriptrun(data, skript, term, title,mod.Foldr,fname)
	if mod.Web{
		switch term{
			case "dxf":		
			txtplot = fname + ".dxf"
			case "svg", "svgmono":
			Svgkong(txtplot)
			txtplot = fname + ".svg"
		}
	}
	pltchn <- txtplot
}

//DrawRectPlan returns 2d coords of the top (plan) view of a rectangle eff. span lx, ly
//with support widths bsupx, y (in adata)

func DrawRectPlan(idx int, xs, ys, lx, ly, bsupx, bsupy float64) (data, sdata, adata, ldata string){
	id := float64(idx)
	xe := xs + lx; ye := ys + ly
	pts := [][]float64{
		{xs, ys, lx, 0.0, id},
		{xs, ys, 0.0, ly, id},
		{xe, ye, -lx, 0.0, id},
		{xe, ye, 0.0, -ly, id},
	}
	x0 := xs - bsupx/2.0; x1 := xs + bsupx/2.0
	y0 := ys - bsupy/2.0; y1 := ys + bsupy/2.0
	for _, pt := range pts{
		data += fmt.Sprintf("%f %f %f %f %f\n",pt[0],pt[1],pt[2],pt[3],pt[4])
	}
	adata += fmt.Sprintf("%f %f %f %f %f\n",xs,ys-0.5*ly,lx,0.0,id)
	ldata += fmt.Sprintf("%f %f %.2f\n",xs+lx/2.0,ys-0.4*ly,lx)
	adata += fmt.Sprintf("%f %f %f %f %f\n",xs-0.5*lx,ys,0.0,ly,id)
	ldata += fmt.Sprintf("%f %f %.2f\n",xs-0.4*lx,ys+ly/2.0,ly)
	//supports as dotted lines
	lx0 := lx + bsupx; ly0 := ly + bsupy 
	xe = x0 + lx0; ye = y0 + ly0
	pts = [][]float64{
		{x0, y0, lx0, 0.0, id},
		{x0, y0, 0.0, ly0, id},
		{xe, ye, -lx0, 0.0, id},
		{xe, ye, 0.0, -ly0, id},
	}
	for _, pt := range pts{
		sdata += fmt.Sprintf("%f %f %f %f %f\n",pt[0],pt[1],pt[2],pt[3],pt[4])
	}
	lx1 := lx - bsupx; ly1 := ly - bsupy 
	xe = x1 + lx1; ye = y1 + ly1
	pts = [][]float64{
		{x0, y0, lx0, 0.0, id},
		{x0, y0, 0.0, ly0, id},
		{xe, ye, -lx0, 0.0, id},
		{xe, ye, 0.0, -ly0, id},
	}
	for _, pt := range pts{
		sdata += fmt.Sprintf("%f %f %f %f %f\n",pt[0],pt[1],pt[2],pt[3],pt[4])
	}
	return
}

//DrawRectView draws a rect of (local x) length l and (local y) height d from pb, pe 
func DrawRectView(mdx int, d float64, pb, pe []float64) (data string){
	d = d/2.0
	l := Dist3d(pb, pe)
	var p0, p1, p2, p3, p4 []float64
	if len(pb) == 1{
		pb = append(pb, 0)
		pe = append(pe, 0)
	}
	p0 = Rotvec(90.0, pb, pe)
	p4 = Lerpvec(d/l, pb, p0)
	p1 = Lerpvec(-d/l, pb, p0)
	p0 = Rotvec(90.0, pe, pb)
	p3 = Lerpvec(-d/l, pe, p0)
	p2 = Lerpvec(d/l, pe, p0)
	for _, pt := range [][][]float64{{p1,p2},{p2,p3},{p3,p4},{p4,p1}}{
		p0 = pt[0]
		dx := pt[1][0] - pt[0][0]
		dy := pt[1][1] - pt[0][1]
		data += fmt.Sprintf("%f %f %f %f %v\n",p0[0],p0[1],dx, dy, mdx)
	}
	return
}


//WHY WONT IT WERK
//DrawLineOffst draws a line offset by +(top)/-(bot) d from pb, pe 
func DrawLineOffst(mdx int, d float64, pb, pe []float64) (data string){
	var p0, p3, p4 []float64
	l := Dist3d(pb, pe)
	p0 = Rotvec(90.0, pb, pe)
	p4 = Lerpvec(d/l, pb, p0)
	//p1 = Lerpvec(-d/l, pb, p0)
	p0 = Rotvec(90.0, pe, pb)
	p3 = Lerpvec(d/l, pe, p0)
	//p2 = Lerpvec(d/l, pe, p0)
	dx := p3[0] - p4[0]
	dy := p3[1] - p3[1]
	data += fmt.Sprintf("%f %f %f %f %v\n",p3[0],p3[1],-dx, -dy, mdx)
	return
}

//PlotBm1d plots a beam model
func PlotBm1d(mod *Model, term string, pltchn chan string) {
	data, ldata, mdata := GeomDat(mod)
	d1, l1, m1 := JloadDat(mod)
	//log.Println("jload data->d1->",d1,",l1->",l1,"m1->",m1)
	data += d1; ldata += l1; mdata += m1
	d1, l1, m1 = MsLoad2Dat(mod)

	data += d1; ldata += l1; mdata += m1
	d1, l1, _ = View2dDat(mod)
	data += d1
	ldata += l1
	data += ldata
	data += "\n\n"; data += mdata
	
	// data += d1; ldata += l1; mdata += m1
	// data += ldata
	// data += "\n\n"; data += mdata
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("1db-%s",mod.Id)
	title := fmt.Sprintf("1d beam %s",mod.Id)
	//fname = svgpath(mod.Foldr,fname, term)
	skript := "drawmod2d.gp"
	txtplot := skriptrun(data, skript, term, title,mod.Foldr,fname)
	//fmt.Println("txtplot-",txtplot)
	if mod.Web{
		switch term{
			case "dxf":		
			txtplot = fname + ".dxf"
			case "svg", "svgmono":
			Svgkong(txtplot)
			txtplot = fname + ".svg"
		}
	}

	pltchn <- txtplot
}


//DrawNpMem returns 2d coords for drawing an np member struct
func DrawNpMem(frmtyp int, dmin float64, mem *MemNp, jb, je []float64)(data string){
	var tdata, bdata, edata string
	switch frmtyp{
		case 0:
		//1dbeam model
		if len(jb)==1{
			jb = append(jb, 0.0)
		}
		if len(je) == 1{
			je = append(je, 0.0)
		}
	}
	lspan := Dist3d(jb, je)
	d1 := dmin/2.0
	if mem.Mtyp == 1{
		d1 = mem.Ds[1]/2.0
	}
	// x1 := mem.Ls[0]; x2 := mem.Ls[0]+mem.Ls[1]
	// t1 := mem.Ts[0]; t2 := mem.Ts[2]
	var l1 float64
	for i, dx := range mem.Dxs{
		//d = d/1000.0
		x := mem.Xs[i]
		var p0, p1, p2, p3 []float64
		//plot these as lines
		if i == 20{
			p1 = Rotvec(-90.0,je, jb)
			p3 = Lerpvec(d1/lspan,je, p1)
			p2 = Lerpvec(-(dx-d1)/lspan, je, p1)
			
		} else {
			p0 = Lerpvec(x/lspan,jb,je)
			l1 = Dist3d(p0,je)
			p1 = Rotvec(90.0, p0, je)
			p3 = Lerpvec(d1/l1, p0, p1)
			p2 = Lerpvec(-(dx-d1)/l1, p0, p1)
		}
		bdata += fmt.Sprintf("%f %f %v %v\n",p2[0],p2[1],mem.Id,i)
		tdata += fmt.Sprintf("%f %f %v %v\n",p3[0],p3[1],mem.Id,i)
		if i == 0{ 
			edata += fmt.Sprintf("%f %f %v %v\n",p3[0],p3[1],mem.Id,i)
			edata += fmt.Sprintf("%f %f %v %v\n",p2[0],p2[1],mem.Id,i)
		}
		if i == 20{
			edata += "\n"
		 	edata += fmt.Sprintf("%f %f %v %v\n",p2[0],p2[1],mem.Id,i)	
			edata += fmt.Sprintf("%f %f %v %v\n",p3[0],p3[1],mem.Id,i)
		}
		// if x == x1 && t1 == 1{
		// 	p2 = Lerpvec(-d1/l1, p0, p1)
		// 	bdata += fmt.Sprintf("%f %f %v %v\n",p2[0],p2[1],mem.Id,i)
		// }
		// if x == x2 && t2 == 1{
		// 	p2 = Lerpvec(-d1/l1, p0, p1)
		// 	//bdata += fmt.Sprintf("%f %f %v %v\n",p2[0],p2[1],mem.Id,i)
		// }
		
		//fmt.Println("mem, i, x, dx->",mem.Id, i, x,dx, "x==x1?",x==x1, "x==x1?",x==x2)
	}
	edata += "\n"
	tdata += "\n"
	bdata += "\n"
	data = bdata + tdata + edata
	return
}

//PlotNpBm1d plots a non-prismatic beam model
func PlotNpBm1d(mod *Model, term string, pltchn chan string) {
	var hdata string
	for _, mem := range mod.Mnps{
		jb := mem.Mprp[0]
		je := mem.Mprp[1]
		pb := mod.Coords[jb-1]
		pe := mod.Coords[je-1]
		h1 := DrawNpMem(mod.Frmtyp,mod.Dmin,mem, pb, pe)
		hdata += h1
	}
	//fmt.Println(hdata)
	data, ldata, mdata := GeomDat(mod)
	d1, l1, m1 := JloadDat(mod)
	data += d1; ldata += l1; mdata += m1
	d1, l1, m1 = MsLoad2Dat(mod)
	data += d1; ldata += l1; mdata += m1
	data += ldata
	data += "\n\n"; data += mdata
	hdata += "\n\n"
	hdata += data
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("np-1db-%s",mod.Id)
	title := fmt.Sprintf("np-1d beam-%s",mod.Id)
	//fname = svgpath(mod.Foldr,fname, term)
	skript := "drawnpmod2d.gp"
	txtplot := skriptrun(hdata, skript, term, title,mod.Foldr,fname)
	if mod.Web{
		switch term{
			case "dxf":		
			txtplot = fname + ".dxf"
			case "svg", "svgmono":
			Svgkong(txtplot)
			txtplot = fname + ".svg"
		}
	}
	pltchn <- txtplot
}


//PlotNpFrm2d plots a non-prismatic beam model
func PlotNpFrm2d(mod *Model, term string, pltchn chan string) {
	var hdata string
	for _, mem := range mod.Mnps{
		jb := mem.Mprp[0]
		je := mem.Mprp[1]
		pb := mod.Coords[jb-1]
		pe := mod.Coords[je-1]
		h1 := DrawNpMem(mod.Frmtyp,mod.Dmin,mem, pb, pe)
		hdata += h1
	}
	data, ldata, mdata := GeomDat(mod)
	d1, l1, m1 := JloadDat(mod)
	data += d1; ldata += l1; mdata += m1
	d1, l1, m1 = MsLoad2Dat(mod)
	data += d1; ldata += l1; mdata += m1
	data += ldata
	data += "\n\n"; data += mdata
	hdata += "\n\n"
	hdata += data
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("np-2dfrm-%s",mod.Id)
	title := fmt.Sprintf("np frame2d-%s",mod.Id)
	skript := "drawnpmod2d.gp"
	txtplot := skriptrun(hdata, skript, term, title,mod.Foldr,fname)
	//if mod.Web{txtplot = fname}
	if mod.Web{
		switch term{
			case "dxf":		
			txtplot = fname + ".dxf"
			case "svg", "svgmono":
			Svgkong(txtplot)
			txtplot = fname + ".svg"
		}
	}
	pltchn <- txtplot
}

//PlotTrs2d plots a 2d truss model
func PlotTrs2d(mod *Model, term string, pltchn chan string) {
	data, ldata, mdata := GeomDat(mod)
	d1, l1, m1 := JloadDat(mod)
	data += d1; ldata += l1; mdata += m1
	d1, l1, m1 = MsLoad2Dat(mod)
	data += d1; ldata += l1; mdata += m1
	if mod.Units != "" && mod.Drawsec{
		d1, l1, _ = View2dDat(mod)
		data += d1
		ldata += l1
	}
	data += ldata
	data += "\n\n"; data += mdata
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("2dt-%s",mod.Id)
	title := fmt.Sprintf("2d truss-%s",mod.Id)
	//fname = svgpath(mod.Foldr,fname, term)
	skript := "drawmod2d.gp"
	//log.Println("DATA-\n",data)
	txtplot := skriptrun(data, skript, term, title,mod.Foldr,fname)
	if mod.Web{
		switch term{
			case "dxf":		
			txtplot = fname + ".dxf"
			case "svg", "svgmono":
			Svgkong(txtplot)
			txtplot = fname + ".svg"
		}
	}
	//fmt.Println("txtplot->\n",txtplot, fname)
	pltchn <- txtplot
}

//PlotFrm3d plots a 3d frame model
func PlotFrm3d(mod *Model, term string, pltchn chan string){
	data, ldata, mdata := GeomDat(mod)
	d1, l1, m1 := JloadDat(mod)
	data += d1; ldata += l1; mdata += m1
	d1, l1, m1 = MsLoad2Dat(mod)
	data += d1; ldata += l1; mdata += m1
	//d1, l1, _ = View2dDat(mod)
	data += d1
	ldata += l1
	data += ldata
	data += "\n\n"; data += mdata
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("%s-3df",mod.Id)
	title := fmt.Sprintf("%s-3d frame",mod.Id)
	//fname = svgpath(mod.Foldr,fname, term)
	skript := "drawmod3d.gp"
	if mod.Zup{skript = "drawmod3d1.gp"}
	txtplot := skriptrun(data, skript, term, title,mod.Foldr,fname)
	if mod.Web{
		switch term{
			case "dxf":		
			txtplot = fname + ".dxf"
			case "svg", "svgmono":
			Svgkong(txtplot)
			txtplot = fname + ".svg"
		}
	}
	pltchn <- txtplot

}

//PlotTrs3d plots (does not at all) a 3d truss model
func PlotTrs3d(mod *Model, term string, pltchn chan string) {
	data, ldata, mdata := GeomDat(mod)
	d1, l1, m1 := JloadDat(mod)
	data += d1; ldata += l1; mdata += m1
	d1, l1, m1 = MsLoad2Dat(mod)
	data += d1; ldata += l1; mdata += m1
	//d1, l1, _ = View2dDat(mod)
	data += d1
	ldata += l1
	data += ldata
	data += "\n\n"; data += mdata
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("%s-3dt",mod.Id)
	title := fmt.Sprintf("%s-3d truss",mod.Id)
	//fname = svgpath(mod.Foldr,fname, term)
	skript := "drawmod3d.gp"
	if mod.Zup{skript = "drawmod3d1.gp"}
	txtplot := skriptrun(data, skript, term, title,mod.Foldr,fname)
	if mod.Web{
		switch term{
			case "dxf":		
			txtplot = fname + ".dxf"
			case "svg", "svgmono":
			Svgkong(txtplot)
			txtplot = fname + ".svg"
		}
	}

	pltchn <- txtplot

	// //3d truss plot (*edit- NOT)-SWAPPING Y AND Z VALUES
	// //as is the way of structures and things
	// var data string
	// //index 0 nodes
	// var frcscale , xmax, ymax, zmax float64
	// frcscale = 0.5
	// //switch mod.Cmdz[1] {
	// //case "kips":
	// //	frcscale = 30.0
	// //case "mks":
	// //	frcscale = 1.0
	// //case "mmks":
	// //	frcscale = 1000.0
	// //}
	// for idx, v := range mod.Coords {
	// 	data += fmt.Sprintf("%v %v %v %v\n", v[0], v[1], v[2], idx+1)
	// 	if v[2] > zmax {zmax = v[2]}
	// 	if v[1] > ymax {ymax = v[1]}
	// 	if v[0] > xmax {xmax = v[0]}
	// }
	// data += "\n\n"
	// //index 1 members
	// //ms := make(map[int][]int)
	// for idx, mem := range mod.Mprp {
	// 	jb := mod.Coords[mem[0]-1]
	// 	je := mod.Coords[mem[1]-1]
	// 	data += fmt.Sprintf("%v %v %v %v %v %v %v\n", jb[0], jb[1], jb[2], je[0], je[1], je[2], idx+1)
	// }
	// data += "\n\n"
	// //index 2 supports
	// for _, val := range mod.Supports {
	// 	pt := mod.Coords[val[0]-1]
	// 	if val[1]+val[2]+val[3] != 0 {data += fmt.Sprintf("%v %v %v\n", pt[0],pt[2],pt[1])}
	// }
	// data += "\n\n"
	// //index 3 joint loads
	// for _, val := range mod.Jloads {
	// 	//var delta float64
	// 	pt := mod.Coords[int(val[0])-1]
	// 	if val[1] != 0.0 { //X- force (assemble?)
	// 		if pt[0] == xmax {
	// 			//vector to the right
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],frcscale, 0, 0, val[1])
	// 		} else {
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],-frcscale, 0, 0, val[1])
	// 		}
	// 	}
	// 	if val[2] != 0.0 { //y force
	// 		if pt[2] == ymax {	
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1], 0, 0, frcscale, val[2])
	// 		} else {
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1]-frcscale,0, 0, frcscale,val[2])
	// 		}
	// 	}
	// 	if val[3] != 0.0 { //z force
	// 		if pt[1] == zmax {
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],0, 0,frcscale, val[3])
	// 		} else {
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1]-frcscale,0, 0, frcscale, val[3])
	// 		}	
	// 	}
	// }
	// data += "\n\n"
	// //create temp files
	// f, e1 := os.CreateTemp("", "barf")
	// if e1 != nil {
	// 	log.Println(e1)
	// }
	// defer f.Close()
	// defer os.Remove(f.Name())
	// _, e1 = f.WriteString(data)
	// if e1 != nil {
	// 	log.Println(e1)
	// }

	// var termstr string
	// switch term {
	// case "dumb":
	// 	termstr = "set term dumb ansi size 79,49"
	// case "dumbstr":
	// 	termstr = "set term dumb size 79,49"

	// case "caca":
	// 	termstr = "set term caca inverted size 79,49"
	// case "wxt":
	// 	termstr = "set term wxt"
	// }

	// setstr := "set autoscale; set key bottom; set title \"SPACE TRUSS\";set grid; set label;set tics;set view 60,30,1,1;set ticslevel 0; set linetype 1 lw 3 pt 5"
	// pltstr := fmt.Sprintf("splot '%s' index 0 using 1:2:3:4 w labels point pt 7 offset char 1,1 notitle,'' index 1 using 1:2:3:($4-$1):($5-$2):($6-$3) notitle w vectors lt 1 nohead, '' index 1 using ($4+$1)/2:($2+$5)/2:($3+$6)/2:7 w labels notitle,'' index 2 using 1:2:3 w points pointtype 19 notitle, '' index 3 using 1:2:3:4:5:6 notitle w vectors, '' index 3 u 1:2:3:5 notitle w labels left offset char 2,2,2", f.Name())
	// prg := "gnuplot"
	// arg0 := "-e"
	// arg2 := "--persist"
	// arg1 := fmt.Sprintf("%s; %s; %s", termstr, setstr, pltstr)
	// plotstr := exec_command(prg, arg2, arg0, arg1)
	// pltchn <- plotstr
	
}


//PlotGrd3d plots a grillage model (horribly)
func PlotGrd3d(mod *Model, term string, pltchn chan string) {
	data, ldata, mdata := GeomDat(mod)
	d1, l1, m1 := JloadDat(mod)
	data += d1; ldata += l1; mdata += m1
	d1, l1, m1 = MsLoad2Dat(mod)
	data += d1; ldata += l1; mdata += m1
	//d1, l1, _ = View2dDat(mod)
	data += d1
	ldata += l1
	data += ldata
	data += "\n\n"; data += mdata
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("%s-3dg",mod.Id)
	title := fmt.Sprintf("%s-3d grid",mod.Id)
	//fname = svgpath(mod.Foldr,fname, term)
	skript := "drawmod3d.gp"
	if mod.Zup{skript = "drawmod3d1.gp"}
	txtplot := skriptrun(data, skript, term, title,mod.Foldr,fname)
	if mod.Web{
		switch term{
			case "dxf":		
			txtplot = fname + ".dxf"
			case "svg", "svgmono":
			Svgkong(txtplot)
			txtplot = fname + ".svg"
		}
	}
	pltchn <- txtplot
	
	//CURSES gnuplot has the z axis as vertical
	//KERSES maybe screw gnuplot (heresy.UNBELIEVER)
	// var data string
	// //index 0 nodes
	// var frcscale , xmax, ymax, zmax float64
	// switch mod.Cmdz[1] {
	// case "kips":
	// 	frcscale = 30.0
	// case "mks":
	// 	frcscale = 1.0
	// case "mmks":
	// 	frcscale = 1000.0
	// }
	// for idx, v := range mod.Coords {
	// 	data += fmt.Sprintf("%v %v %v %v\n", v[0], v[2], v[1], idx+1)
	// 	if v[2] > zmax {zmax = v[1]}
	// 	if v[1] > ymax {ymax = v[2]}
	// 	if v[0] > xmax {xmax = v[0]}
	// }
	// data += "\n\n"
	// //index 1 members
	// //ms := make(map[int][]int)
	// for idx, mem := range mod.Mprp {
	// 	jb := mod.Coords[mem[0]-1]
	// 	je := mod.Coords[mem[1]-1]
	// 	data += fmt.Sprintf("%v %v %v %v %v %v %v\n", jb[0], jb[2], jb[1], je[0], je[2], je[1], idx+1)
	// }
	// data += "\n\n"
	// //index 2 supports
	// for _, val := range mod.Supports {
	// 	pt := mod.Coords[val[0]-1]
	// 	if val[1]+val[2]+val[3] != 0 {data += fmt.Sprintf("%v %v %v\n", pt[0],pt[2],pt[1])}
	// }
	// data += "\n\n"
	// //index 3 joint loads
	// for _, val := range mod.Jloads {
	// 	//var delta float64
	// 	pt := mod.Coords[int(val[0])-1]
	// 	if val[1] != 0.0 { //X- force (assemble?)
	// 		if pt[0] == xmax {
	// 			//vector to the right
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],frcscale, 0, 0, val[1])
	// 		} else {
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],-frcscale, 0, 0, val[1])
	// 		}
	// 	}
	// 	if val[2] != 0.0 { //y force
	// 		if pt[2] == ymax {
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1], 0,  0, frcscale, val[2])
	// 		} else {
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2], pt[1]-frcscale,0, 0, frcscale, val[2])
	// 		}
	// 	}
	// 	if val[3] != 0.0 { //z force
	// 		if pt[1] == zmax {
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],0, 0,frcscale, val[3])
	// 		} else {
	// 			data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1]-frcscale,0, 0, frcscale, val[3])
	// 		}	
	// 	}
	// }
	// data += "\n\n"
	// //create temp files
	// f, e1 := os.CreateTemp("", "barf")
	// if e1 != nil {
	// 	log.Println(e1)
	// }
	// defer f.Close()
	// defer os.Remove(f.Name())
	// _, e1 = f.WriteString(data)
	// if e1 != nil {
	// 	log.Println(e1)
	// }
	// var termstr string
	// switch term {
	// case "dumb":
	// 	termstr = "set term dumb ansi size 79,49"
	// case "caca":
	// 	termstr = "set term caca inverted size 79,49"
	// case "wxt":
	// 	termstr = "set term wxt"
	// case "dumbstr":
	// 	termstr = "set term dumb size 79,49"

	// }
	// setstr := "set autoscale; set key bottom; set title \"SPACE FRAME\";set grid; set label;set tics;set view 60,30,1,1;set ticslevel 0; set linetype 1 lw 0 pt 5"
	// pltstr := fmt.Sprintf("splot '%s' index 0 using 1:2:3:4 w labels point pt 7 offset char 1,1 notitle,'' index 1 using 1:2:3:($4-$1):($5-$2):($6-$3) notitle w vectors lt 1 nohead, '' index 1 using ($4+$1)/2:($2+$5)/2:($3+$6)/2:7 w labels notitle,'' index 2 using 1:2:3 w points pointtype 19 notitle, '' index 3 using 1:2:3:4:5:6 notitle w vectors, '' index 3 u 1:2:3:5 notitle w labels left offset char 2,2,2", f.Name())
	// prg := "gnuplot"
	// arg0 := "-e"
	// arg2 := "--persist"
	// arg1 := fmt.Sprintf("%s; %s; %s", termstr, setstr, pltstr)
	// plotstr := exec_command(prg, arg2, arg0, arg1)
	// pltchn <- plotstr
	
}


//GeomDat returns node and member text data for plots
func GeomDat(mod *Model) (data, ldata, mdata string){
	//plot nodes and members
	//if len(mod.Dims) == len(mod.Cp) {draw member dims}
	//index 0 nodes
	for idx, v := range mod.Coords {
		switch len(v){
			case 1:
			//beam
			data += fmt.Sprintf("%.1f %.1f %v\n", v[0], 0.0, idx+1)
			case 2:
			//truss or frame
			data += fmt.Sprintf("%.1f %.1f %v\n", v[0], v[1], idx+1)
			case 3:
			//3d truss, grid or frame
			data += fmt.Sprintf("%.1f %.1f %.1f %v\n", v[0], v[1], v[2], idx+1)
		}
	}
	data += "\n\n"
	//index 1 members
	for idx, m := range mod.Mprp {
		jb := mod.Coords[m[0]-1]
		je := mod.Coords[m[1]-1]
		em := m[2]
		cp := m[3]
		
		switch len(jb){
			case 1:
			data += fmt.Sprintf("%.1f %.1f %.1f %.1f %v %v %v\n", jb[0], 0.0, je[0], 0.0, idx+1, cp, em)
			ldata += fmt.Sprintf("%f %f %v %v\n",(jb[0]+je[0])/2.0, 0.0, idx+1,idx+1)
			case 2:
			data += fmt.Sprintf("%.1f %.1f %.1f %.1f %v %v %v\n", jb[0], jb[1], je[0], je[1], idx+1, cp, em)
			ldata += fmt.Sprintf("%f %f %v %v\n",(jb[0]+je[0])/2.0, (jb[1]+je[1])/2.0, idx+1, idx+1)
			case 3:
			data += fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f %.1f %v %v %v\n", jb[0], jb[1], jb[2],je[0], je[1], je[2],idx+1, cp, em)
			ldata += fmt.Sprintf("%f %f %f %v %v\n",(jb[0]+je[0])/2.0, (jb[1]+je[1])/2.0, (jb[1]+je[1])/2.0, idx+1, idx+1)
		}
		
	}
	data += "\n\n"
	//index 2 supports
	for _, val := range mod.Supports{
		pt := mod.Coords[val[0]-1]
		switch len(val){
			case 3:
			//beam or truss
			switch len(pt){
				case 1:
				data += fmt.Sprintf("%.1f %.1f %v\n", pt[0],0.0,-(val[1]+val[2]))
				case 2:
				data += fmt.Sprintf("%.1f %.1f %v\n", pt[0],pt[1],-(val[1]+val[2]))
			}
			case 4:
			//frame or truss
			switch len(pt){
				case 2:
				//2d frame
				data += fmt.Sprintf("%.1f %.1f %v\n", pt[0],pt[1],-(val[1]+val[2]+val[3]))
				case 3:
				//3d truss or grid
				data += fmt.Sprintf("%.1f %.1f %.1f %v\n", pt[0],pt[1],pt[2],-(val[1]+val[2]+val[3]))
			}
			default:
			//space frame
			var sum int
			for _, v := range val[1:]{sum -= v}
			data += fmt.Sprintf("%.1f %.1f %.1f %v\n", pt[0],pt[1],pt[2],sum)
		}
	}
	data += "\n\n"
	return
}

//View2dDat returns 2d view coordinates (no it doesn't)
//it does nothing at all
func View2dDat(mod *Model)(data, ldata, mdata string){
	units := mod.Units
	
	// if _, ok := AUnits[mod.Units]; ok{
	// 	units = mod.Units
	// }
	//fmt.Println("UNITZ-?",mod.Units)
	for idx, m := range mod.Mprp{
		pb := mod.Coords[m[0]-1]
		pe := mod.Coords[m[1]-1]
		cp := m[3]
		dims := []float64{150.0,150.0}
		switch units{
			case "knm", "kn-m":
			dims = []float64{0.15,0.15}
			case "nmm","n-mm":
			dims = []float64{150,150}
			case "kpin","kp-in":
			dims = []float64{6.0,6.0}
		}
		styp := 1
		if len(mod.Dims) >= cp{
			dims = mod.Dims[cp-1]
		}
		if len(mod.Sts) >= cp{
			styp = mod.Sts[cp-1]
		}
		data += DrawMem2d(idx + 1, styp, pb, pe, dims)
		// pl := Lerpvec(0.25, pb, pe)
		// diml := "("
		// for _, dim := range dims{
		// 	diml += fmt.Sprintf("%.2f,",dim)
		// }
		// diml += ")"
		// ldata += fmt.Sprintf("%f %f %s\n",pl[0], pl[1], diml)
	}
	data += "\n\n"
	return
}

//PlotGenTrs plots a truss given coords, members
func PlotGenTrs(coords [][]float64, ms [][]int, term... string){
	var data string
	//all yr ys are belong to 1.0
	for idx, v := range coords {
		data += fmt.Sprintf("%v %v %v\n", v[0], v[1], idx+1)
	}
	data += "\n\n"
	for i, m := range ms {
		jb := coords[m[0]-1]
		je := coords[m[1]-1]
		data += fmt.Sprintf("%f %f %f %f %v %f\n",jb[0],jb[1],je[0],je[1],i+1,math.Cos((je[0]-jb[0])/(je[1]-jb[1])))
	}
	data += "\n\n"
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	pltskript := filepath.Join(basepath,"/t2dgenplot.gp")
	
	f, e1 := os.CreateTemp("", "barf")
	if e1 != nil {
		fmt.Println(e1)
	}	
	_, e1 = f.WriteString(data)
	if e1 != nil {
		fmt.Println(e1)
	}
	
	defer f.Close()
	defer os.Remove(f.Name())
	var trmstr string
	if len(term) == 0{
		trmstr = "dumb"
	} else {
		trmstr = term[0]
	}
	cmd := exec.Command("gnuplot","-c",pltskript,f.Name(),trmstr)
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
	fmt.Println(outstr)
	//return
} 

//DrawMem2d returns 2d coords for drawing a member (front view)
func DrawMem2d(mdx, styp int, pb, pe, dims []float64) (data string){
	var d float64
	//var d1, d2 float64
	//fmt.Println("mdx, styp, pb, pe, dims",mdx, styp, pb, pe, dims)
	if len(dims) == 0{return}
	switch styp{
		case 0:
		d = dims[0]
		default:
		d = dims[1]
	}
	//d = d/1000.0
	data = DrawRectView(mdx, d, pb, pe) 
	switch styp{
		case 6:
		//t: sec
		//fmt.Println("here,mem-",mdx)
		//df := dims[3]
		//data += DrawLineOffst(mdx, d/2.0 - df, pb, pe)
		case 12:
		//eq. i section
		//add top n bottom flange lines
		tf := dims[2]
		data += DrawRectView(mdx, d - 2.0*tf, pb, pe)
	}
	//data += "\n"
	return
}



//DrawDetail plots a portal frame detail (see fig 1 saka 2003)
func (f *Portal) DrawDetail(){
	//first draw columns
	switch f.Config{
		case 0:
		case 1:
		case 2:
		case 3:
	}
}

//exec_command executes a shell command and returns output
func exec_command(program string, args ...string) string {
	cmd := exec.Command(program, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outstr, errstr := stdout.String(), stderr.String()
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
	}
	if errstr != "" {
		log.Println(errstr)
	}
	return outstr

}



//exec_wxt executes the wxt terminal of gnuplot
//it is a worthless function, here serves as padding
func exec_wxt(program string, args ...string) string {
	cmd := exec.Command(program, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	//outstr, errstr := stdout.String(), stderr.String()
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
	}
	return "done"

}

//Xsort sorts coords by x and then y values
func Xsort(cords [][]float64) {
	sort.Slice(cords[:], func(i, j int) bool {
		for x := range cords[i] {
			if cords[i][x] == cords[j][x] {
				continue
			}
			return cords[i][x] < cords[j][x]
		}
		return false
	})
}

//DrawMod2d draws a 2d model, dunno where it's used
func DrawMod2d(mod *Model, ms map[int]*Mem, term string) (string){
	var data, mdata, ldata string
	var fsx, xmin, ymin, xmax, ymax, yflr float64
	//index 0 nodes
	for idx, v := range mod.Coords {
		data += fmt.Sprintf("%.1f %.1f %v\n", v[0], v[1], idx+1)
		if v[1] > ymax {ymax = v[1]}
		if v[0] > xmax {xmax = v[0]}
		if v[1] < ymin {ymin = v[1]}
		if v[0] < xmin {xmin = v[0]}
	}
	//ideally one should calc yflr
	yflr = 1.25
	//and fsx is just max x val
	fsx = 1.5
	data += "\n\n"
	//index 1 members
	for idx, mem := range mod.Mprp {
		jb := mod.Coords[mem[0]-1]
		je := mod.Coords[mem[1]-1]
		cp := mem[3]
		data += fmt.Sprintf("%.1f %.1f %.1f %.1f %v %v\n", jb[0], jb[1], je[0], je[1], idx+1, cp)
		ldata += fmt.Sprintf("%f %f %v 1\n",(jb[0]+je[0])/2.0, (jb[1]+je[1])/2.0, idx+1)
	}
	data += "\n\n"
	//index 2 supports
	for _, val := range mod.Supports {
		pt := mod.Coords[val[0]-1]
		if val[1]+val[2] != 0 {data += fmt.Sprintf("%.1f %.1f\n", pt[0],pt[1])}
	}
	data += "\n\n"
	//index 3 joint loads
	for _, val := range mod.Jloads{
		pt := mod.Coords[int(val[0])-1]
		if val[1] != 0.0 { 
			if pt[0] == xmax {
				data += fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f\n",pt[0],pt[1],slerp(fsx,1.0,1.0,val[1]), 0.0, val[1])
			} else {
				data += fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f\n",pt[0],pt[1],-slerp(fsx,1.0,1.0,val[1]), 0.0, val[1])
			}
			ldata += fmt.Sprintf("%f %f %.0f 2\n",pt[0], pt[1], val[1])
		}
		if val[2] != 0.0 {
			if pt[1] == ymax {
				data += fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f\n",pt[0],pt[1],0.0, -slerp(fsx, 1.0, 1.0, val[2]), val[2])
			} else {
				data += fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f\n",pt[0],pt[1],0.0, slerp(fsx, 1.0, 1.0, val[2]), val[2])
			}
			ldata += fmt.Sprintf("%f %f %.0f 2\n",pt[0], pt[1], val[2])
		}
		if val[3] != 0.0 {
			mdata += fmt.Sprintf("%.1f %.1f %.1f %.1f\n",pt[0],pt[1],slerp(fsx, 1.0, 1.0, val[3]), val[3])
			ldata += fmt.Sprintf("%f %f %.0f 2\n",pt[0], pt[1], val[3])
		}
	}
	data += "\n\n"
	//index 4 member loads
	for _, val := range mod.Msloads {
		m := int(val[0])
		mem := ms[m]
		jb := mod.Mprp[m-1][0]
		xa, ya := mod.Coords[jb-1][0], mod.Coords[jb-1][1]
		ltyp := int(val[1])
		wa, wb, la, lb := val[2], val[3], val[4], val[5]
		cx := mem.Geoms[4]; cy := mem.Geoms[5]
		//ldata += fmt.Sprintf("%f %f %.0f\n",xa+la*cx, ya+la*cy, wa)
		ya += 1.0
		switch ltyp {
		case 1://point load at la
			data += fmt.Sprintf("%f %f %f %f %v\n",xa + la * cx, ya + la * cy, slerp(yflr, 1.0, 0.5, wa)*cy, -slerp(yflr, 1.0, 0.5, wa)*cx, ltyp)
			ldata += fmt.Sprintf("%f %f %.0f 3\n",xa+la*cx, ya+la*cy, wa)
		case 2:
			//moment at la 
			mdata += fmt.Sprintf("%f %f %f\n",xa + la * cx, ya + la * cy, slerp(yflr, 1.0, 0.5, wa))
		case 3://udl w from la to l - lb
			l := mem.Geoms[0]
			div := (l-lb-la)/5.0
			xa += la * cx; ya += la * cy
			xa -= div * cx; ya -= div * cy
			for i:=0; i < 5; i++{
				xa += div * cx; ya += +div * cy
				data += fmt.Sprintf("%f %f %f %f %.0f %v\n",xa,ya,-slerp(yflr, 1.0, 0.5, wa)*cy,slerp(yflr, 1.0, 0.5, wa)*cx,wa, ltyp)
				if i == 2{ldata += fmt.Sprintf("%f %f %f 3\n",xa, ya, wa)}
			}
		case 4://udl wa at la to wb at l - lb
			l := mem.Geoms[0]
			div := (l-lb-la)/5.0
			dw := (wb - wa)/5.0
			xa -= div * cx; ya -= div * cy
			for i:=0; i < 5; i++{
				xa += div * cx ; ya += div * cy
				wx := wa + dw * float64(i)
				data += fmt.Sprintf("%f %f %f %f %v\n",xa,ya,-slerp(yflr, 1.0, 0.5, wx)*cy,slerp(yflr, 1.0, 0.5, wx)*cx, ltyp)
				if i == 2{ldata += fmt.Sprintf("%f %f %f 3\n",xa, ya, wa)}
			}
		case 5:
			//point axial load at la
			data += fmt.Sprintf("%f %f %f %f %v\n",xa+la*cx,ya+la*cy,slerp(yflr, 1.0, 0.5, wa)*cy,-slerp(yflr, 1.0, 0.5, wa)*cx, ltyp)
		case 6:
			//uniform axial load w at la to l - lb
			l := mem.Geoms[0]
			div := (l-lb-la)/3.0
			//xa -= div * cx; ya -= div * cy
			for i:=0; i < 3; i++{
				xa += div * cx; ya += div * cy
				data += fmt.Sprintf("%f %f %f %f %v\n",xa,ya,slerp(yflr, 1.0, 0.0, wa)*cx,slerp(yflr, 1.0, 0.0, wa)*cy, ltyp)
				if i == 2{ldata += fmt.Sprintf("%f %f %f 3\n",xa, ya, wa)}
			}
		case 7:
			//torsional moment?
		}
	}
	data += "\n\n"; ldata += "\n\n"
	//index 5 labels
	data += ldata
	//index 6 moments
	data += mdata
	//fname := fmt.Sprintf("m2d_%s",mod.Id)
	//title := "2d frame"
	//skript := "drawmod2d.gp"
	//txtplot, err := draw.Draw(data, skript, term, folder, fname, title) 
	txtplot := ""
	return txtplot
}


/*
   
You can pass arguments to a gnuplot script since version 5.0, with the flag -c. These arguments are accessed through the variables ARG0 to ARG9, ARG0 being the script, and ARG1 to ARG9 string variables. The number of arguments is given by ARGC.

For example, the following script ("script.gp")

#!/usr/local/bin/gnuplot --persist

THIRD=ARG3
print "script name        : ", ARG0
print "first argument     : ", ARG1
print "third argument     : ", THIRD 
print "number of arguments: ", ARGC 
can be called as:

$ gnuplot -c script.gp one two three four five
script name        : script.gp
first argument     : one
third argument     : three
number of arguments: 5
or within gnuplot as

gnuplot> call 'script.gp' one two three four five
script name        : script.gp
first argument     : one
third argument     : three
number of arguments: 5

gnuplot -e "datafile='${data}'; outputname='${output}'" foo.plg


   old 3d truss plot view
   
	// setstr := "set autoscale; set key bottom; set title \"SPACE TRUSS\";set grid; set label;set tics;set view 60,30,1,1;set ticslevel 0; set linetype 1 lw 3 pt 5"
	// pltstr := fmt.Sprintf("splot '%s' index 0 using 1:2:3:4 w labels point pt 7 offset char 1,1 notitle,'' index 1 using 1:2:3:($4-$1):($5-$2):($6-$3) notitle w vectors lt 1 nohead, '' index 1 using ($4+$1)/2:($2+$5)/2:($3+$6)/2:7 w labels notitle,'' index 2 using 1:2:3 w points pointtype 19 notitle, '' index 3 using 1:2:3:4:5:6 notitle w vectors, '' index 3 u 1:2:3:5 notitle w labels left offset char 2,2,2", f.Name())
	// prg := "gnuplot"
*/
