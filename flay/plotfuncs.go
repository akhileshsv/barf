package barf

import (
	"fmt"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"	
	draw "barf/draw"
)

func EdgeDat(a, b Pt2d, idx int)(data string){
	dx := b.X - a.X
	dy := b.Y - a.Y
	data = fmt.Sprintf("%f %f %f %f %v\n",a.X,a.Y,dx,dy,idx)
	return
}

func (f *Flr) DrawPolys()(txtplot string, err error){
	var ndata, data, ldata string
	for i, poly := range f.Polys{
		for j, p1 := range poly{
			ndata += fmt.Sprintf("%f %f %v\n",p1.X,p1.Y,j)
			
			ldata += fmt.Sprintf("%f %f %v\n",p1.X,p1.Y,j)
			var p2 Pt2d
			if j == len(poly) - 1{
				p2 = poly[0]
			} else {
				p2 = poly[j+1]
			}
			data += EdgeDat(p1, p2, 1)
		}
		pc := Centroid2d(poly)
		ldata += fmt.Sprintf("%f %f %s\n",pc.X,pc.Y,f.Labels[i])
	}
	ndata += "\n\n"
	data = ndata + data
	data += "\n\n"
	data += ldata
	skript := "basic2d.gp"
	xl := f.Units; yl := f.Units; zl := f.Units
	term := f.Term
	folder := ""
	fname := ""
	title := "poly plan"
	txtplot, err = draw.Draw(data, skript, term, folder, fname, title, xl, yl, zl) 
	return	
}

//DrawCorg draws the corridor graph for a floor
func (f *Flr) DrawCorg()(txtplot string, err error){
	
	var data, ldata string
	//plot nodes	
	for i, pt := range f.Nodes{
		data += fmt.Sprintf("%f %f %v\n",pt.X,pt.Y,i+1)
		ldata += fmt.Sprintf("%f %f %v\n",pt.X, pt.Y,i+1)
	}
	data += "\n\n"
	//plot edges
	for _, edx := range f.Iwalls{
		val := f.Emap[edx]
		jb := val[0][0]
		je := val[0][1]
		ecls := val[0][2]
		eon := val[2][0]
		if eon == 1{			
			if ecls < 1{ecls = 2}
			a := f.Nodes[jb-1]
			b := f.Nodes[je-1]
			data += EdgeDat(a,b,ecls)	
		}
	}
	
	if len(f.Tmp) > 0{
		eds, _ := f.Tmp[0].([][]Pt2d)
		for _, ed := range eds{
			a := ed[0]; b := ed[1]
			data += EdgeDat(a, b, 3)
		}
	}
	data += "\n\n"
	//plot labels
	for i, r := range f.Rooms{
		ldata += fmt.Sprintf("%f %f %s\n",(r.Origin.X+r.End.X)/2.0,(r.Origin.Y+r.End.Y)/2.0,f.Labels[i])
	}
	data += ldata
	skript := "basic2d.gp"
	xl := f.Units; yl := f.Units; zl := f.Units
	term := f.Term
	folder := ""
	fname := ""
	title := "floor plan"
	txtplot, err = draw.Draw(data, skript, term, folder, fname, title, xl, yl, zl) 
	return
}

//Draw plots a floor
func (f *Flr) Draw()(txtplot string, err error){
	var data, ldata string
	//plot nodes	
	for i, pt := range f.Nodes{
		data += fmt.Sprintf("%f %f %v\n",pt.X,pt.Y,i+1)
		ldata += fmt.Sprintf("%f %f %v\n",pt.X, pt.Y,i+1)
	}
	data += "\n\n"
	//plot edges
	for _, val := range f.Emap{
		jb := val[0][0]
		je := val[0][1]
		ecls := val[0][2]
		if ecls < 1{ecls = 2}
		a := f.Nodes[jb-1]
		b := f.Nodes[je-1]
		if val[2][0] > 0{
			data += EdgeDat(a,b,ecls)
		}
	}
	data += "\n\n"
	//plot labels
	for i, r := range f.Rooms{
		switch f.Sqrd{
			case false:			
			ldata += fmt.Sprintf("%f %f %s\n",(r.Origin.X+r.End.X)/2.0,(r.Origin.Y+r.End.Y)/2.0,f.Labels[i])
			case true:			
			ldata += fmt.Sprintf("%f %f %s\n",r.Mid.X,r.Mid.Y,f.Labels[i])
		}
	}
	data += ldata
	skript := "basic2d.gp"
	xl := f.Units; yl := f.Units; zl := f.Units
	term := f.Term
	folder := ""
	fname := ""
	title := "floor plan"
	txtplot, err = draw.Draw(data, skript, term, folder, fname, title, xl, yl, zl) 
	return
	
}

//GpDatFloors returns the gnuplot datafile for a line/plan view of a floor with rooms
func GpDatFloors (f *Flr) (data , filename string) {
	//generates temp data file for gnuplot polygons
	//list of (n) vertices x1 y1 x2 y2 x3 y3 x4 y4
	//https://stackoverflow.com/questions/37607583/i-want-to-plot-a-rectangle-with-given-4-coordinates-in-a-text-file-in-gnuplot-t?rq=1
	//boxxyerror bars = x0,y0, xe, ye
	//https://stackoverflow.com/questions/28648740/plotting-rectangle-side-by-side-from-coordinates
	//list of x,y vertices separated by data block
	//https://stackoverflow.com/questions/32781536/gnuplot-how-to-draw-polygon-contour-from-its-vertices
	var x0, y0, xe, ye float64
	var name string
	x0 = f.Origin.X
	y0 = f.Origin.Y
	xe = f.End.X
	ye = f.End.Y
	name = f.Name
	data += fmt.Sprintf("%v %v %v %v %s\n",x0,y0,xe,ye,"")
	
	for _, r := range f.Rooms { 
		x0 = r.Origin.X
		y0 = r.Origin.Y
		xe = r.End.X
		ye = r.End.Y
		name = r.Name
		data += fmt.Sprintf("%v %v %v %v %s\n", x0,y0,xe,ye,name)
		if len(r.Rooms) > 0 {
			for _, r1 := range r.Rooms { 
				x0 = r1.Origin.X
				y0 = r1.Origin.Y
				xe = r1.End.X
				ye = r1.End.Y
				name = r1.Name
				data += fmt.Sprintf("%v %v %v %v %s\n", x0,y0,xe,ye,name)
			}
		}
 	}
	
	//create temp files
	file, e := os.CreateTemp("", "floorrect")
	if e != nil {
		fmt.Println(e)
	}
	defer file.Close()
	//defer os.Remove(f.Name())
	_, e = file.WriteString(data)
	if e != nil {
		fmt.Println(e)
	}
	filename = file.Name()
	return 
}

//GPlotFloors plots a floor struct
func GPlotFloors(f *Flr, dumb bool) { //
	//set loadpath for gnuplot
	prg := "gnuplot"
	arg0 := "-persist"
	arg1 := "-e"
	arg2 := " set autoscale; set key bottom; set title \"BC\";set offsets graph 0.1,0.1,0.1,0.1;"
	//set terminal
	if dumb {
		arg2 += "set term dumb ansi size 79,39;"
	} else {
		arg2 += "set terminal qt;"
	}
	
	//create data file for gnuplot
	_, filename := GpDatFloors(f)
	arg3 := fmt.Sprintf("plot '%v' using (($1+$3)/2):(($2+$4)/2):(($3-$1)/2):(($4-$2)/2) w boxxyerrorbars notitle,'' using (($1+$3)/2):(($2+$4)/2):5 with labels tc '#0000ff' notitle",filename)
	arg2 += arg3
	s := exec_command(prg, arg0,arg1,arg2)
	if dumb {fmt.Println(s)}
}

//exec_command. good ol' traveller of local folders,
//hail fellow well met even to the bear
func exec_command(program string, args ...string) string {
	cmd := exec.Command(program, args...)
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



//Plotgrid plots a grid of rooms using ...
//gnuplot (' ')
//each rect has uniform dims (dx, dy)
func Plotgrid(grid [][]int, dx, dy float64) (string){
	var x0, y0, x, y, xc, yc float64
	var data, cdat string
	for i, row := range grid {
		x0 = 0.0
		x = 0.0
		y0 = y
		y += dy
		for j, room := range row {
			x += dx
			xc = float64(j)*dx + dx/2.0
			yc = float64(i)*dy + dy/2.0
			cdat += fmt.Sprintf("%v %v %v\n",xc, yc, room)
			
			data += fmt.Sprintf("%v %v %v\n",x0,y0,room)
			data += fmt.Sprintf("%v %v %v\n",x0,y,room)
			data += fmt.Sprintf("%v %v %v\n",x,y,room)
			data += fmt.Sprintf("%v %v %v\n",x,y0,room)
			data += fmt.Sprintf("%v %v %v\n",x0,y0,room)
			x0 = x
		}
		data += "\n"
	}
	data += "\n\n" + cdat
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	pltskript := filepath.Join(basepath,"/gridplot.gp")
	
	f, e1 := os.CreateTemp("", "flay")
	if e1 != nil {
		fmt.Println(e1)
	}	
	_, e1 = f.WriteString(data)
	if e1 != nil {
		fmt.Println(e1)
	}
	
	defer f.Close()
	defer os.Remove(f.Name())
	
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

//PltLout plots a layout (lout) which is a map of rooms 
func PltLout(rmap map[int]*Rm) (pltstr string){
	//get plotscript filepath
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	pltskript := filepath.Join(basepath,"/loutplot.gp")
	var data string
	for i := range rmap{
		rm := rmap[i]
		for _, walls := range rm.Walls {
			for _, wall := range walls {
				if wall.Typ != -1 {
					data += fmt.Sprintf("%v %v %v %v\n",wall.Pb.X,wall.Pb.Y,wall.Typ,rm.Id)
					data += fmt.Sprintf("%v %v %v %v\n",wall.Pe.X,wall.Pe.Y,wall.Typ,rm.Id)
				}
				data += "\n"
			}
		}
	}
	data += "\n\n"
	for i := range rmap{
		rm := rmap[i]
		for _, cell := range rm.Cells {
			data += fmt.Sprintf("%v %v %v\n",cell.Pb.X,cell.Pb.Y,rm.Id)
			data += fmt.Sprintf("%v %v %v\n",cell.Pb.X,cell.Pe.Y,rm.Id)
			data += fmt.Sprintf("%v %v %v\n",cell.Pe.X,cell.Pe.Y,rm.Id)
			data += fmt.Sprintf("%v %v %v\n",cell.Pe.X,cell.Pb.Y,rm.Id)
			data += fmt.Sprintf("%v %v %v\n",cell.Pb.X,cell.Pb.Y,rm.Id)
			data += "\n"
		}
	}
	data += "\n\n"
	for i := range rmap{
		rm := rmap[i]
		for _, cell := range rm.Cells {
			data += fmt.Sprintf("%v %v %v\n",cell.Centroid.X,cell.Centroid.Y,rm.Id)
		}
	}
	data += "\n\n"
	for i := range rmap {
		data += fmt.Sprintf("%v %v %v\n",rmap[i].Centroid.X,rmap[i].Centroid.Y,rmap[i].Id)
	}
	data += "\n\n"
	//fmt.Println("DATA->",data)
	f, e1 := os.CreateTemp("", "flay")
	if e1 != nil {
		fmt.Println(e1)
	}
	defer f.Close()
	defer os.Remove(f.Name())	
	_, e1 = f.WriteString(data)
	if e1 != nil {
		fmt.Println(e1)
	}
	cmd := exec.Command("gnuplot","-c",pltskript,f.Name(),"qt")
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
