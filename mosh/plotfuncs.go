package barf

import (
	"fmt"
	"log"
	"math"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	kass"barf/kass"
	//draw"barf/draw"
)

//skriptpath returns the absolute path of the (gnuplot) script named skript in current folder
func skriptpath(skript string) (string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, skript)
}

//genfname returns the absolute path of the svg output file for gnuplot
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
	if err != nil {
		fmt.Println(err)
	}
	if errstr != "" {
		fmt.Println(errstr)
	}
	return outstr

}

//PlotColGeom plots an rcc column section 
func PlotColGeom(c *RccCol, term string) (pltstr string){
	var fname, data string
	if c.Data == ""{data = c.Draw()} else {data = c.Data}
	if c.Title == ""{
		c.Title = fmt.Sprintf("rcc column section %v-%v", c.Mid, c.Id)
		//fname = fmt.Sprintf("col%v-%v.svg",c.Mid, c.Id)
	}
	title := c.Title + ".svg"
	fname = genfname(c.Foldr,title)
	pltstr = skriptrun(data, "plotcolgeom.gp", term, c.Title, fname)
	if term == "svg"{
		pltstr = fname
	}
	c.Txtplot = append(c.Txtplot, pltstr)
	return
}

//PlotBmGeom plots an rcc beam section
func PlotBmGeom(b *RccBm, term string) (pltstr string){
	var title, fname string
	data, err := b.Draw()
	if err != nil{
		log.Println(err)
		return
	}
	if b.Title == ""{
		b.Title = fmt.Sprintf("rcc beam section %v-%v", b.Mid, b.Id)
	} else {
		b.Title = fmt.Sprintf("rcc beam section %s", b.Title)
		//title = fmt.Sprintf("rcc beam section %s %v-%v", b.Title, b.Mid, b.Id)
		//fname = fmt.Sprintf("%s_%v-%v.svg",b.Title,b.Mid,b.Id)
	}
	title = b.Title + ".svg"
	fname = genfname(b.Foldr,title)
	pltstr = skriptrun(data, "plotbmgeom.gp", term, b.Title, fname)
	if term == "svg"{
		pltstr = fname
	}
	b.Txtplot = append(b.Txtplot, pltstr)
	fmt.Println(ColorRed,"pltstr",pltstr,ColorReset)
	return
}

//PlotColNM plots column section n-m interaction curves
func PlotColNM(pus, mus []float64) (pltstr string) {
	//get plotscript filepath
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
func DrawColCircle(c *RccCol, term string) (err error){
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
	title := c.Title + ".svg"
	fname := genfname(c.Foldr,title)
	pltstr := skriptrun(data, "plotcolcircle.gp", term, c.Title, fname)
	if term == "svg"{
		pltstr = fname
	}
	//fmt.Println("PLOTSTR->\n",pltstr)
	c.Txtplot = append(c.Txtplot, pltstr)
	err = nil
	return
}

//DrawColRect plots a rectangular rcc column section 
func DrawColRect(c *RccCol, term string) (err error){
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
	if c.Title == ""{
		c.Title = fmt.Sprintf("rcc col section %v-%v", c.Mid, c.Id)
	}
	title := c.Title + ".svg"
	fname := genfname(c.Foldr,title)
	pltstr := skriptrun(data, "plotcolrect.gp", term, c.Title, fname)
	if term == "svg"{
		pltstr = fname
	}
	c.Txtplot = append(c.Txtplot, pltstr)
	err = nil
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
				ye := sf.H[0] * (val[2])/(sf.DL*sf.PSFs[0])/40.0 + yb
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
	fname := genfname(sf.Foldr,fn)
	//fmt.Println(fname)
	pltskript := "plotsubfrm.gp"
	outstr := skriptrun(data, pltskript, term, fn, fname)
	if term == "svg"{
		outstr = fname
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
