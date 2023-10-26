package barf

//model plotting funcs via gnuplot

import (
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
func svgpath(foldr, fname string) (string){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	fname = fname + ".svg"
	if foldr == ""{
		foldr = filepath.Join(basepath,"../data/out")
	}
	return filepath.Join(foldr,fname)
}

//skriptrun runs a gnuplot script and returns the path of svg/text plot/error string
func skriptrun(data, pltskript, term, title, fname string) (string){
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
		fname = svgpath("",fname)
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
		//fmt.Println(errstr)
	}
	return outstr

}

//PlotFrm2d plots a 2d frame model (line plot)
func PlotFrm2d(mod *Model, term string, pltchn chan string) {
	data, ldata, mdata := GeomDat(mod)
	d1, l1, m1 := JloadDat(mod)
	data += d1; ldata += l1; mdata += m1
	d1, l1, m1 = MsLoad2Dat(mod)
	data += d1; ldata += l1; mdata += m1
	data += ldata
	data += "\n\n"; data += mdata
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("2df-%s",mod.Id)
	title := fmt.Sprintf("2d frame-%s",mod.Id)
	fname = svgpath(mod.Foldr,fname)
	skript := "drawmod2d.gp"
	txtplot := skriptrun(data, skript, term, title,fname)
	pltchn <- txtplot
}


//PlotFrm2dRez plots a 2d frame model with bm/sf/def. results
func PlotFrm2dRez(mod *Model, term string){
	var data, rdata, cdata string
	//var ncs [][]float64
	//translate/rotate coords by nodal displ. - FORGET THIS FOR NOW
	/*
	org := []float64{0,0}
	for _, node := range mod.Js{
		tx, ty, ang := node.Displ[0], node.Displ[1], node.Displ[2] * 180.0/math.Pi
		p0 := node.Coords
		p1 := Trans2d(p0, tx, ty)
		_ = Rotvec(ang, org, p1) - this is not what it is, see fig 6.8 kassimali
		ncs = append(ncs, p1)
	}
	*/
	for i, m := range mod.Mprp {
		jb := mod.Coords[m[0]-1]
		je := mod.Coords[m[1]-1]
		//jbn := ncs[m[0]-1]
		//jen := ncs[m[1]-1]
		lmem := Dist3d(jb, je)
		mem := mod.Ms[i+1]
		data += fmt.Sprintf("%f %f %f %f %v\n",jb[0], jb[1], je[0], je[1], i)
		//cdata += fmt.Sprintf("%f %f %f %f %v\n",jbn[0], jbn[1], jen[0], jen[1], i)
		for j, x := range mem.Rez.Xs{
			vu, bm, dx := mem.Rez.SF[j] * mod.Scales[0], mem.Rez.BM[j]* mod.Scales[1], mem.Rez.Dxs[j] * 20.0
			if j == 20{
				p1 := je
				p2 := Rotvec(90,je,jb)
				pv := Lerpvec(-vu/lmem, p1, p2)
				pb := Lerpvec(-bm/lmem, p1, p2)
				pd := Lerpvec(-dx/lmem, p1, p2)
				rdata += fmt.Sprintf("%f %f %f %f %f %f\n",pv[0], pv[1],pb[0], pb[1],pd[0], pd[1])
				continue
			}
			p1 := Lerpvec(x/lmem, jb, je)
			d0 := Dist3d(p1, je)
			p2 := Rotvec(90,p1,je)
			pv := Lerpvec(vu/d0, p1, p2)
			pb := Lerpvec(bm/d0, p1, p2)
			pd := Lerpvec(-dx/d0, p1, p2)
			rdata += fmt.Sprintf("%f %f %f %f %f %f\n",pv[0], pv[1],pb[0], pb[1],pd[0], pd[1])			
		}
		rdata += "\n"
		/*
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
	}
	rdata += "\n\n"
	data += "\n\n"; data += rdata; data += cdata
	//data += ldata
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("2df-%s",mod.Id)
	title := fmt.Sprintf("2d frame-%s",mod.Id)
	fname = svgpath(mod.Foldr,fname)
	skript := "drawmodrez.gp"
	txtplot := skriptrun(data, skript, term, title,fname)
	if term == "dumb" || term == "mono"{fmt.Println(txtplot)}
	return
}

//PlotBm1d plots a beam model
func PlotBm1d(mod *Model, term string, pltchn chan string) {
	data, ldata, mdata := GeomDat(mod)
	d1, l1, m1 := JloadDat(mod)
	data += d1; ldata += l1; mdata += m1
	d1, l1, m1 = MsLoad2Dat(mod)
	data += d1; ldata += l1; mdata += m1
	data += ldata
	data += "\n\n"; data += mdata
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("1db-%s",mod.Id)
	title := fmt.Sprintf("1d beam-%s",mod.Id)
	fname = svgpath(mod.Foldr,fname)
	skript := "drawmod2d.gp"
	txtplot := skriptrun(data, skript, term, title,fname)
	pltchn <- txtplot
}

//PlotTrs2d plots a 2d truss model
func PlotTrs2d(mod *Model, term string, pltchn chan string) {
	data, ldata, mdata := GeomDat(mod)
	d1, l1, m1 := JloadDat(mod)
	data += d1; ldata += l1; mdata += m1
	d1, l1, m1 = MsLoad2Dat(mod)
	data += d1; ldata += l1; mdata += m1
	data += ldata
	data += "\n\n"; data += mdata
	//create temp files
	if mod.Id == ""{
		mod.Id = fmt.Sprintf("%v",rand.Intn(666))
	}
	fname := fmt.Sprintf("2dt-%s",mod.Id)
	title := fmt.Sprintf("2d truss-%s",mod.Id)
	fname = svgpath(mod.Foldr,fname)
	skript := "drawmod2d.gp"
	txtplot := skriptrun(data, skript, term, title,fname)
	pltchn <- txtplot

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
			ldata += fmt.Sprintf("%f %f %v 1\n",(jb[0]+je[0])/2.0, 0.0, idx+1)
			case 2:
			data += fmt.Sprintf("%.1f %.1f %.1f %.1f %v %v %v\n", jb[0], jb[1], je[0], je[1], idx+1, cp, em)
			ldata += fmt.Sprintf("%f %f %v 1\n",(jb[0]+je[0])/2.0, (jb[1]+je[1])/2.0, idx+1)
			case 3:
			data += fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f %.1f %v %v %v\n", jb[0], jb[1], jb[2],je[0], je[1], je[2],idx+1, cp, em)
			ldata += fmt.Sprintf("%f %f %f %v 1\n",(jb[0]+je[0])/2.0, (jb[1]+je[1])/2.0, (jb[1]+je[1])/2.0, idx+1)
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


//JLoadDat returns nodal force plot data
func JloadDat(mod *Model) (data, ldata, mdata string){
	fsx := 1.0
	if len(mod.Cmdz) > 1{
		switch mod.Cmdz[1]{
			case "kips":
			fsx = 40.0
			case "mmks":
			fsx = 1000.0
		}
	}
	
	//index 3 joint loads
	for _, val := range mod.Jloads{
		pt := mod.Coords[int(val[0])-1]
		switch len(pt){
			case 1:
			if val[1] != 0.0{
				data += fmt.Sprintf("%.1f %.1f %.1f %.1f %v\n",pt[0],0.0,slerp(fsx,1.0,1.0,val[1]), 0.0, 1)
				ldata += fmt.Sprintf("%f %f %.1f 1\n",pt[0], 0.0, val[1])
			}
			if val[2] != 0.0{
				mdata += fmt.Sprintf("%.1f %.1f %.1f %v\n",pt[0],0.0,slerp(fsx*3.0,1.0,1.0,math.Abs(val[2])), 2)
				ldata += fmt.Sprintf("%f %f %.1f 2\n",pt[0], 0.0, val[2])
			}
			case 2:
			switch len(val){
				case 3:
				//t2d
				if val[1] != 0.0{
					data += fmt.Sprintf("%.1f %.1f %.1f %.1f %v\n",pt[0],pt[1],slerp(fsx,1.0,1.0,val[1]), 0.0, 1)
					ldata += fmt.Sprintf("%f %f %.1f 1\n",pt[0]+slerp(fsx,1.0,1.0,val[1])/2.0, pt[1], val[1])
				}
				if val[2] != 0.0{
					data += fmt.Sprintf("%.1f %.1f %.1f %.1f %.v\n",pt[0],pt[1],0.0, slerp(fsx,1.0,1.0,-val[2]), 2)
					ldata += fmt.Sprintf("%f %f %.1f 2\n",pt[0], pt[1]+slerp(fsx,1.0,1.0,-val[2])/2.0, val[2])
				}
				case 4:
				//f2d
				if val[1] != 0.0{
					data += fmt.Sprintf("%.1f %.1f %.1f %.1f %v\n",pt[0],pt[1],-slerp(fsx,1.0,1.0,val[1]), 0.0, 1)
					ldata += fmt.Sprintf("%f %f %.1f 1\n",pt[0]-slerp(fsx,1.0,1.0,val[1])/2.0, pt[1], val[1])
				}
				if val[2] != 0.0{
					data += fmt.Sprintf("%.1f %.1f %.1f %.1f %v\n",pt[0],pt[1],0.0, slerp(fsx,1.0,1.0,val[2]), 2)
					ldata += fmt.Sprintf("%f %f %.1f 2\n",pt[0], pt[1], val[2])
				}
				
				if val[3] != 0.0{
					mdata += fmt.Sprintf("%.1f %.1f %.1f %v\n",pt[0],pt[1],slerp(fsx/5.0,1.0,1.0,val[3]), 3)
					ldata += fmt.Sprintf("%f %f %.1f 3\n",pt[0], pt[1], val[2])
				}
			}
			case 3:
			switch len(val){
				case 3:
				//t3d, g3d
				if val[1] != 0.0 {
					data += fmt.Sprintf("%v %v %v %v %v %v %v\n",pt[0],pt[1],pt[2],slerp(fsx,1.0,1.0,val[1]), 0, 0, 1)
				}
				if val[2] != 0.0 {
					data += fmt.Sprintf("%v %v %v %v %v %v %v\n",pt[0],pt[1],pt[2], 0, slerp(fsx,1.0,1.0,val[2]), 0, 2)
				}
				if val[3] != 0.0 {
					mdata += fmt.Sprintf("%v %v %v %v %v %v %v\n",pt[0],pt[1],pt[2],0, 0, slerp(fsx/5.0,1.0,1.0,val[3]), 3)
				}	
			}
			default:
			//f3d (LMAO)
			
		}
	}
	//data += "\n\n"
	return
}


//MsLoad2Dat returns (2d) member force plot data
func MsLoad2Dat(mod *Model) (data, ldata, mdata string){
	//index 4 member loads
	//var xa, ya, za float64
	yflr := 1.0
	if len(mod.Cmdz) > 1{
		switch mod.Cmdz[1]{
			case "kips":
			yflr = 0.25
			case "mmks":
			yflr = 10.0
		}
	}
	//if mod.Frcscale > 0.0{yflr = mod.Frcscale}
	for _, val := range mod.Msloads {

		m := int(val[0])
		//fmt.Println("mem->", m, "val->",val)
		//mem := ms[m]
		jb := mod.Mprp[m-1][0]
		je := mod.Mprp[m-1][1]
		pa := mod.Coords[jb-1]
		pb := mod.Coords[je-1]
		ltyp := int(val[1])
		wa, wb, la, lb := val[2], val[3], val[4], val[5]
		l := Dist3d(pa,pb)
		switch ltyp{
			case 1://point load at la
			p1 := Lerpvec(la/l,pa,pb)
			wscl := slerp(yflr,1.0,0.5,wa)
			
			switch len(pa){
				case 1:
				wscl = slerp(yflr/10.0,1.0,0.5,math.Abs(wa))
				data += fmt.Sprintf("%f %f %f %f %v\n",p1[0], 0.0, 0.0, wscl, ltyp)
				ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], wscl/5.0, wa, ltyp)
				default:
				p2 := Rotvec(90,p1,pb)
				d0 := Dist3d(p1,p2)
				p3 := Lerpvec(wscl/d0, p1, p2)
				switch len(pa){
					case 2:
					data += fmt.Sprintf("%f %f %f %f %v\n",p1[0], p1[1], p3[0]-p1[0], p3[1]-p1[1], ltyp)
					ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], p1[1]+wscl, wa, ltyp)
					case 3:
					//i assume z remains the same? I MIGHT BE WRONG AS USUAL
					data += fmt.Sprintf("%f %f %f %f %f %f %v\n",p1[0], p1[1], p1[2], p3[0]-p1[0], p3[1]-p1[1], p1[2],ltyp)
					ldata += fmt.Sprintf("%f %f %f %.2f %v\n",p1[0], p1[1]+wscl, p1[2], wa, ltyp)
				}
			}
			case 2:
			//moment at la
			p1 := Lerpvec(la/l,pa,pb)
			wscl := slerp(yflr/5.0,1.0,0.5,wa)
			switch len(pa){
				case 1:
				mdata += fmt.Sprintf("%f %f %f\n",p1[0], 0.0, wscl/5.0)
				case 2:
				mdata += fmt.Sprintf("%f %f %f\n",p1[0], p1[1], wscl)
				case 3:
				mdata += fmt.Sprintf("%f %f %f %f\n",p1[0], p1[1], p1[2], wscl)
			}
			case 3://udl w from la to l - lb
			step := (l-lb-la)/10.0
			p1 := Lerpvec(la/l,pa,pb)
			wscl := slerp(yflr,1.0,0.5,wa)
			//wscl := wa/yflr
			for i:=0; i < 10; i++{
				switch len(pa){
					case 1:		
					wscl = slerp(yflr/10.0,1.0,0.5,wa)
					data += fmt.Sprintf("%f %f %f %f %v\n",p1[0], 0.0, 0.0, -wscl, ltyp)
					switch i{
						case 4:
						ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], -wscl/5.0, wa, ltyp)
					}
					default:
					//wscl = wa/100.0
					p2 := Rotvec(90.0,p1, pb)
					d0 := Dist3d(p1,p2)
					p3 := Lerpvec(wscl/d0, p1, p2)
					switch len(pa){
						case 2:
						data += fmt.Sprintf("%f %f %f %f %v\n",p1[0], p1[1], p3[0]-p1[0], p3[1]-p1[1], ltyp)
						switch i{
							case 4:
							ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], p1[1], wa, ltyp)
						}
						case 3:
						data += fmt.Sprintf("%f %f %f %f %f %f %v\n",p1[0], p1[1], p1[2], p3[0]-p1[0], p3[1]-p1[1], p1[2],ltyp)
						switch i{
							case 4:
							ldata += fmt.Sprintf("%f %f %f %.2f %v\n",p1[0], p1[1], p1[2], wa, ltyp)
						}
					}
				}
				d1 := Dist3d(p1,pb)
				p1 = Lerpvec(step/d1,p1,pb)
			}
			case 4://udl wa at la to wb at l - lb
			step := (l-lb-la)/5.0
			p1 := Lerpvec(la/l,pa,pb)
			if la == 0.0{p1 = pa}
			dw := math.Abs((wb - wa)/(l-lb-la))/100.0
			//wscl := slerp(yflr, 1.0, 0.5, wa)
			//dw = slerp(1.0, 1.0, 0.0, dw/5.0)
			//if wscl == 0.0{wscl = 0.01}
			wscl := wa/100.0
			
			//fmt.Println("wscl init->",wscl,dw,l-lb-la)
			for i:=0; i < 5; i++{
				switch len(pa){
					case 1:
					data += fmt.Sprintf("%f %f %f %f %v\n",p1[0], 0.0, 0.0, wscl/5.0, ltyp)
					switch i{
						case 1:
						//ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], wscl/4.0, wa, ltyp)
						case 3:
						ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], wscl/4.0, wb, ltyp)
					}
					default:
					p2 := Rotvec(90.0, p1, pb)
					d0 := Dist3d(p1,p2)
					p3 := Lerpvec(wscl/d0, p1, p2)
					//fmt.Println("wscl,p1,p2,p3",wscl,p1,p2,p3)
					switch len(pa){
						case 2:
						data += fmt.Sprintf("%f %f %f %f %v\n",p1[0], p1[1], p3[0]-p1[0], p3[1]-p1[1], ltyp)
						switch i{
							case 1:
							//ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], p1[1], wa, ltyp)
							case 3:
							ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], p1[1], wb, ltyp)
						}
						case 3:
						data += fmt.Sprintf("%f %f %f %f %f %f %v\n",p1[0], p1[1], p1[2], p3[0]-p1[0], p3[1]-p1[1], p1[2],ltyp)
						switch i{
							case 1:
							//ldata += fmt.Sprintf("%f %f %f %.2f %v\n",p1[0], p1[1], p1[2], wa, ltyp)
							case 3:
							ldata += fmt.Sprintf("%f %f %f %.2f %v\n",p1[0], p1[1], p1[1], wb, ltyp)
						}
						//ldata += fmt.Sprintf("%f %f %f %.2f %v\n",p1[0], p1[1], p1[2], wa, ltyp)
					}
				}
				
				if wb < wa{
					wscl -= dw/5.0
				} else {
					wscl += dw/5.0
				}	
				d1 := Dist3d(p1,pb)
				p1 = Lerpvec(step/d1,p1,pb)
			}
			case 5:
			//point axial load at la
			p1 := Lerpvec(la/l,pa,pb)
			wscl := slerp(yflr,1.0,0.5,wa)
			switch len(pa){
				case 1:
				data += fmt.Sprintf("%f %f %f %f %v\n",p1[0], 0.5, wscl, 0.5, ltyp)
				ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], 0.5, wa, ltyp)
				default:
				d0 := Dist3d(p1,pb)
				p3 := Lerpvec(wscl/d0, p1, pb)
				switch len(pa){
					case 2:
					data += fmt.Sprintf("%f %f %f %f %v\n",p1[0], p1[1], p3[0]-p1[0], p3[1]-p1[1], ltyp)
					ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], p1[1]+wscl, wa, ltyp)
					case 3:
					//i assume z remains the same? I MIGHT BE WRONG AS USUAL
					data += fmt.Sprintf("%f %f %f %f %f %f %v\n",p1[0], p1[1], p1[2], p3[0]-p1[0], p3[1]-p1[1], p1[2],ltyp)
					ldata += fmt.Sprintf("%f %f %f %.2f %v\n",p1[0], p1[1]+wscl, p1[2], wa, ltyp)
				}
			}
			case 6:
			//TODO uniform axial load w at la to l - lb
			step := (l-lb-la)/10.0
			p1 := Lerpvec(la/l,pa,pb)
			wscl := slerp(yflr, 1.0, 0.5, wa)
			for i:=0; i < 10; i++{
				switch len(pa){
					case 1:
					data += fmt.Sprintf("%f %f %f %f %v\n",p1[0], 0.5, p1[0], wscl, ltyp)
					switch i{
						case 4:
						ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], 0.5, wa, ltyp)
					}					
					default:
					d0 := Dist3d(p1,pb)
					p3 := Lerpvec(wscl/d0, p1, pb)
					switch len(pa){
						case 2:
						data += fmt.Sprintf("%f %f %f %f %v\n",p1[0], p1[1], p3[0]-p1[0], p3[1]-p1[1], ltyp)
						switch i{
							case 4:
							ldata += fmt.Sprintf("%f %f %.2f %v\n",p1[0], p1[1], wa, ltyp)
						}						
						case 3:
						data += fmt.Sprintf("%f %f %f %f %f %f %v\n",p1[0], p1[1], p1[2], p3[0]-p1[0], p3[1]-p1[1], p1[2],ltyp)
						switch i{
							case 4:
							ldata += fmt.Sprintf("%f %f %f %.2f %v\n",p1[0], p1[1], p1[2], wa, ltyp)
						}
					}
				}
				d1 := Dist3d(p1,pb)
				p1 = Lerpvec(step/d1,p1,pb)
			}
			case 7:
			//TODO torsional moment?
			case 8:
			//HOWDO temperature changes?
			case 9:
			//FABRICATION ERRORS?
			//huh? HUH?
		}
	}
	data += "\n\n"
	return
}

//PlotGenTrs plots a truss given coords, members
func PlotGenTrs(coords [][]float64, ms [][]int, pltchn chan string){
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
	pltchn <- outstr
} 

//PlotGrd3d plots a grillage model (horribly)
func PlotGrd3d(mod *Model, term string, pltchn chan string) {
	//CURSES gnuplot has the z axis as vertical
	//KERSES maybe screw gnuplot
	var data string
	//index 0 nodes
	var frcscale , xmax, ymax, zmax float64
	switch mod.Cmdz[1] {
	case "kips":
		frcscale = 30.0
	case "mks":
		frcscale = 1.0
	case "mmks":
		frcscale = 1000.0
	}
	for idx, v := range mod.Coords {
		data += fmt.Sprintf("%v %v %v %v\n", v[0], v[2], v[1], idx+1)
		if v[2] > zmax {zmax = v[1]}
		if v[1] > ymax {ymax = v[2]}
		if v[0] > xmax {xmax = v[0]}
	}
	data += "\n\n"
	//index 1 members
	//ms := make(map[int][]int)
	for idx, mem := range mod.Mprp {
		jb := mod.Coords[mem[0]-1]
		je := mod.Coords[mem[1]-1]
		data += fmt.Sprintf("%v %v %v %v %v %v %v\n", jb[0], jb[2], jb[1], je[0], je[2], je[1], idx+1)
	}
	data += "\n\n"
	//index 2 supports
	for _, val := range mod.Supports {
		pt := mod.Coords[val[0]-1]
		if val[1]+val[2]+val[3] != 0 {data += fmt.Sprintf("%v %v %v\n", pt[0],pt[2],pt[1])}
	}
	data += "\n\n"
	//index 3 joint loads
	for _, val := range mod.Jloads {
		//var delta float64
		pt := mod.Coords[int(val[0])-1]
		if val[1] != 0.0 { //X- force (assemble?)
			if pt[0] == xmax {
				//vector to the right
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],frcscale, 0, 0, val[1])
			} else {
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],-frcscale, 0, 0, val[1])
			}
		}
		if val[2] != 0.0 { //y force
			if pt[2] == ymax {
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1], 0,  0, frcscale, val[2])
			} else {
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2], pt[1]-frcscale,0, 0, frcscale, val[2])
			}
		}
		if val[3] != 0.0 { //z force
			if pt[1] == zmax {
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],0, 0,frcscale, val[3])
			} else {
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1]-frcscale,0, 0, frcscale, val[3])
			}	
		}
	}
	data += "\n\n"
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
	var termstr string
	switch term {
	case "dumb":
		termstr = "set term dumb ansi size 79,49"
	case "caca":
		termstr = "set term caca inverted size 79,49"
	case "wxt":
		termstr = "set term wxt"
	case "dumbstr":
		termstr = "set term dumb size 79,49"

	}
	setstr := "set autoscale; set key bottom; set title \"SPACE FRAME\";set grid; set label;set tics;set view 60,30,1,1;set ticslevel 0; set linetype 1 lw 0 pt 5"
	pltstr := fmt.Sprintf("splot '%s' index 0 using 1:2:3:4 w labels point pt 7 offset char 1,1 notitle,'' index 1 using 1:2:3:($4-$1):($5-$2):($6-$3) notitle w vectors lt 1 nohead, '' index 1 using ($4+$1)/2:($2+$5)/2:($3+$6)/2:7 w labels notitle,'' index 2 using 1:2:3 w points pointtype 19 notitle, '' index 3 using 1:2:3:4:5:6 notitle w vectors, '' index 3 u 1:2:3:5 notitle w labels left offset char 2,2,2", f.Name())
	prg := "gnuplot"
	arg0 := "-e"
	arg2 := "--persist"
	arg1 := fmt.Sprintf("%s; %s; %s", termstr, setstr, pltstr)
	plotstr := exec_command(prg, arg2, arg0, arg1)
	pltchn <- plotstr
	
}

//PlotTrs3d plots (does not at all) a 3d truss model
func PlotTrs3d(mod *Model, term string, pltchn chan string) {
	//3d truss plot (*edit- NOT)-SWAPPING Y AND Z VALUES
	//as is the way of structures and things
	var data string
	//index 0 nodes
	var frcscale , xmax, ymax, zmax float64
	frcscale = 0.5
	//switch mod.Cmdz[1] {
	//case "kips":
	//	frcscale = 30.0
	//case "mks":
	//	frcscale = 1.0
	//case "mmks":
	//	frcscale = 1000.0
	//}
	for idx, v := range mod.Coords {
		data += fmt.Sprintf("%v %v %v %v\n", v[0], v[1], v[2], idx+1)
		if v[2] > zmax {zmax = v[2]}
		if v[1] > ymax {ymax = v[1]}
		if v[0] > xmax {xmax = v[0]}
	}
	data += "\n\n"
	//index 1 members
	//ms := make(map[int][]int)
	for idx, mem := range mod.Mprp {
		jb := mod.Coords[mem[0]-1]
		je := mod.Coords[mem[1]-1]
		data += fmt.Sprintf("%v %v %v %v %v %v %v\n", jb[0], jb[1], jb[2], je[0], je[1], je[2], idx+1)
	}
	data += "\n\n"
	//index 2 supports
	for _, val := range mod.Supports {
		pt := mod.Coords[val[0]-1]
		if val[1]+val[2]+val[3] != 0 {data += fmt.Sprintf("%v %v %v\n", pt[0],pt[2],pt[1])}
	}
	data += "\n\n"
	//index 3 joint loads
	for _, val := range mod.Jloads {
		//var delta float64
		pt := mod.Coords[int(val[0])-1]
		if val[1] != 0.0 { //X- force (assemble?)
			if pt[0] == xmax {
				//vector to the right
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],frcscale, 0, 0, val[1])
			} else {
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],-frcscale, 0, 0, val[1])
			}
		}
		if val[2] != 0.0 { //y force
			if pt[2] == ymax {	
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1], 0, 0, frcscale, val[2])
			} else {
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1]-frcscale,0, 0, frcscale,val[2])
			}
		}
		if val[3] != 0.0 { //z force
			if pt[1] == zmax {
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1],0, 0,frcscale, val[3])
			} else {
				data += fmt.Sprintf("%v %v %v %v %v %v %.1f\n",pt[0],pt[2],pt[1]-frcscale,0, 0, frcscale, val[3])
			}	
		}
	}
	data += "\n\n"
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

	var termstr string
	switch term {
	case "dumb":
		termstr = "set term dumb ansi size 79,49"
	case "dumbstr":
		termstr = "set term dumb size 79,49"

	case "caca":
		termstr = "set term caca inverted size 79,49"
	case "wxt":
		termstr = "set term wxt"
	}

	setstr := "set autoscale; set key bottom; set title \"SPACE TRUSS\";set grid; set label;set tics;set view 60,30,1,1;set ticslevel 0; set linetype 1 lw 3 pt 5"
	pltstr := fmt.Sprintf("splot '%s' index 0 using 1:2:3:4 w labels point pt 7 offset char 1,1 notitle,'' index 1 using 1:2:3:($4-$1):($5-$2):($6-$3) notitle w vectors lt 1 nohead, '' index 1 using ($4+$1)/2:($2+$5)/2:($3+$6)/2:7 w labels notitle,'' index 2 using 1:2:3 w points pointtype 19 notitle, '' index 3 using 1:2:3:4:5:6 notitle w vectors, '' index 3 u 1:2:3:5 notitle w labels left offset char 2,2,2", f.Name())
	prg := "gnuplot"
	arg0 := "-e"
	arg2 := "--persist"
	arg1 := fmt.Sprintf("%s; %s; %s", termstr, setstr, pltstr)
	plotstr := exec_command(prg, arg2, arg0, arg1)
	pltchn <- plotstr
	
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

*/
