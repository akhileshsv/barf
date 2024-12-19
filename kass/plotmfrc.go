package barf

import (
	"fmt"
	"log"
	"math"
)



//plotltyp1 plots a member load of type 1(point load)
func plotltyp1(val, pa, pb, wng []float64, fsx, l, fscale float64)(data, ldata string){
	//wa, wb, la, lb := lvec[2], lvec[3], lvec[4], lvec[5]
	wa := val[2]; la := val[4]
	lcase := 1
	x0 := fsx/2.0; k := 1.0
	if len(val) > 6{
		lcase = int(val[6])
	}
	yld := 0.0
	if lcase > 1{
		yld = fsx*float64(lcase-1)
	}
	if len(pa) == 1{
		pa = append(pa, 0.0)
		pb = append(pb, 0.0)
	}
	var wscl float64
	if fscale > 0.0{
		wscl = math.Abs(wa)/fscale
	} else {
		wscl = slerp(fsx, k, x0, math.Abs(wa))

	}
	var ivec []float64
	if len(pa) == 3{
		ax := 2
		if len(val) > 7{ax = int(val[7])+2}
		ivec = LocVec(ax, pa, pb, wng)
	}	
	var x, y, z, dx, dy, dz float64
	switch len(pa){
		case 2:		
		p1 := Lerpvec(la/l,pa,pb)
		p2 := Rotvec(90,p1,pb)
		d0 := Dist3d(p1, p2)
		p3 := Lerpvec(yld/d0, p1, p2)
		d0 = Dist3d(p3, p2)
		p4 := Lerpvec(wscl/d0, p3, p2)		
		if wa > 0.0{
			//draw arrow from p4 to p3 (negative direction of member axes)
			//see page 232, kassimali	
			x = p4[0]; y = p4[1]
			dx = -p4[0] + p3[0]
			dy = -p4[1] + p3[1]
		} else {
			//draw arrow from p3 t0 p4
			x = p3[0]; y = p3[1]
			dx = p4[0] - p3[0]
			dy = p4[1] - p3[1]
		}
		data += fmt.Sprintf("%f %f %f %f %v\n",x,y,dx,dy,lcase)
		ldata += fmt.Sprintf("%f %f %.2f %v\n",x+dx/2.0,y+dy/2.0,wa,lcase)
		case 3:
		
		p1 := Lerpvec(la/l,pa,pb)
		p1 = Lerp3d(yld, ivec, p1)
		p2 := Lerp3d(wscl, ivec, p1)
		if wa > 0.0{
			x = p2[0]; y = p2[1]; z = p2[2]
			dx = p1[0] - p2[0]; dy = p1[1] - p2[1]; dz = p1[2] - p2[2] 
		} else {
			x = p1[0]; y = p1[1]; z = p1[2]
			dx = p2[0] - p1[0]; dy = p2[1] - p1[1]; dz = p2[2] - p1[2] 
		}
		data += fmt.Sprintf("%f %f %f %f %f %f %v\n",x,y,z,dx,dy,dz,lcase)
		ldata += fmt.Sprintf("%f %f %f %.2f %v\n",x+dx/2.0,y+dy/2.0,z+dz/2.0,wa,lcase)
	}
	return
	
}

//plotltyp2 plots a member load of type 2(point moment)
//TODO - add arrow
func plotltyp2(val, pa, pb, wng []float64, fsx, l, fscale float64)(mdata, ldata string){
	//moment wa at la
	
	if len(pa) == 1{
		pa = append(pa, 0.0)
		pb = append(pb, 0.0)
	}
	wa, la := val[2], val[4]
	lcase := 1
	if len(val) > 6 {lcase = int(val[6])}
	x0 := fsx/2.0; k := 1.0
	p1 := Lerpvec(la/l,pa,pb)
	var wscl float64
	if fscale > 0.0{
		wscl = math.Abs(wa)/fscale
	} else {
		wscl = slerp(fsx/5.0,k,x0,wa)	
	}
	//circle at a point with diameter d
	switch len(pa){
		case 2:
		mdata += fmt.Sprintf("%f %f %f %v\n",p1[0], p1[1], wscl, lcase)
		ldata += fmt.Sprintf("%f %f %.1f %v\n",p1[0]+wscl/2.0, p1[1]+wscl/2.0, wa, lcase)
		case 3:
		mdata += fmt.Sprintf("%f %f %f %f %v\n",p1[0], p1[1], p1[2], wscl, lcase)
		ldata += fmt.Sprintf("%f %f %f %.1f %v\n",p1[0]+wscl/2.0, p1[1]+wscl/2.0, p1[2]+wscl/2.0, wa, lcase)
	}
	return
}

//plotltyp3 plots a member load of type 3(udl)
func plotltyp3(val, pa, pb, wng []float64, fsx, l, fscale float64)(data, ldata string){
	//udl w from la to l - lb
	if len(pa) == 1{
		pa = append(pa, 0.0)
		pb = append(pb, 0.0)
	}
	wa, la, lb := val[2], val[4], val[5]
	lcase := 1
	if len(val) > 6{lcase = int(val[6])}
	x0 := fsx/2.0; k := 1.0
	step := (l-lb-la)/10.0
	var dw float64
	if fscale > 0.0{
		dw = math.Abs(wa)/fscale
	} else{
		dw = slerp(fsx,k,x0,math.Abs(wa))
	}
	yld := fsx*float64(lcase-1)
	var ivec []float64
	if len(pa) == 3{
		ax := 2
		if len(val) > 7{ax = int(val[7])+2}
		ivec = LocVec(ax, pa, pb, wng)
	}	
	var pla, plb []float64
	if yld > 0.0{
		switch len(pa){
			case 2:
			p1 := Rotvec(90, pa, pb)
			p2 := Rotvec(270,pb, pa)
			pla = Lerpvec(yld/Dist3d(pa,p1), pa, p1)
			plb = Lerpvec(yld/Dist3d(pb, p2), pb, p2)
			case 3:
			pla = Lerp3d(yld, ivec, pa)
			plb = Lerp3d(yld, ivec, pb)
		}
	} else {
		pla = pa
		plb = pb
	}
	p1 := make([]float64, len(pa))
	copy(p1, pla)
	if la > 0.0{
		p1 = Lerpvec(la/Dist3d(pla, plb), pla, plb)
	}
	d := la
	var x, y, z, dx, dy, dz float64
	for i:=0; i < 11; i++{
		// p2 := Rotvec(90, p1, plb)
		// p3 := Lerpvec(dw/Dist3d(p1,p2),p1, p2)
		//print
		switch len(pa){
			case 2:
			p2 := Rotvec(90, p1, plb)
			p3 := Lerpvec(dw/Dist3d(p1,p2),p1, p2)
			dx = p3[0]-p1[0]
			dy = p3[1]-p1[1]
			switch{
				case wa > 0.0:
				//arrow from p3 to p1
				x = p3[0]; y = p3[1]
				dx = -dx
				dy = -dy
				case wa < 0.0:
				//arrow from p1 to p3
				x = p1[0]; y = p1[1]
			}
			data += fmt.Sprintf("%f %f %f %f %v\n",x,y,dx,dy,lcase)
			if i == 4{
				ldata += fmt.Sprintf("%f %f %.2f %v\n",x+dx/2.0,y+dy/2.0,wa,lcase)
			}
			case 3:
			//p2 = p1 + dw * iy
			//iy = local y axis vector
			dx = ivec[0] * dw
			dy = ivec[1] * dw
			dz = ivec[2] * dw
			switch{
				case wa > 0.0:
				x = p1[0] + dx; y = p1[1] + dy; z = p1[2] + dz
				dx = -dx; dy = -dy; dz = -dz
				case wa < 0.0:
				x = p1[0]; y = p1[1]; z = p1[2]
			}
			data += fmt.Sprintf("%f %f %f %f %f %f %v\n",x,y,z,dx,dy,dz,lcase)
			if i == 4{
				ldata += fmt.Sprintf("%f %f %f %.2f %v\n",x+dx/2.0,y+dy/2.0,z+dz/2.0,wa,lcase)
			}
		}
		d += step
		p1 = Lerpvec(d/l,pla,plb)
		if i == 10{
			p1 = plb
		}		
	}
	return
}

//plotltyp4 plots a member load of type 4(trap)
func plotltyp4(val, pa, pb, wng []float64, fsx, l, fscale float64)(data, ldata string){
	//linear inc. load; wa at la from left -> wb at lb from right
	if len(pa) == 1{
		pa = append(pa, 0.0)
		pb = append(pb, 0.0)
	}
	wa, wb, la, lb := val[2], val[3], val[4], val[5]
	lcase := 1
	if len(val) > 6{lcase = int(val[6])}
	x0 := fsx/2.0; k := 1.0
	step := (l-lb-la)/10.0
	yld := 0.0
	if lcase > 1{
		yld = fsx*float64(lcase-1)
	}
	
	var ivec []float64
	if len(pa) == 3{
		ax := 2
		if len(val) > 7{ax = int(val[7])+2}
		ivec = LocVec(ax, pa, pb, wng)
	}
	var pla, plb []float64
	if yld > 0.0{
		switch len(pa){
			case 2:
			p1 := Rotvec(90, pa, pb)
			p2 := Rotvec(270,pb, pa)
			pla = Lerpvec(yld/Dist3d(pa,p1), pa, p1)
			plb = Lerpvec(yld/Dist3d(pb, p2), pb, p2)
			case 3:
			pla = Lerp3d(yld, ivec, pa)
			plb = Lerp3d(yld, ivec, pb)
		}
	} else {
		pla = pa
		plb = pb
	}
	p1 := make([]float64, len(pa))
	copy(p1, pla)
	// d := la
	
	var x, y, z, dx, dy, dz, wscl float64
	fscon := fscale > 0.0
	for i := 0; i < 11; i++{
		dl := la + float64(i) * step
		wx := wa + (wb-wa)*dl/(l-la-lb)
		//log.Println("wx",wx)
		if fscon{
			wscl = math.Abs(wx)/fscale
		} else{
			wscl = slerp(fsx,k,x0,math.Abs(wx))
		}
		p1 = Lerpvec(dl/l, pla, plb)
		switch len(pa){
			case 2:
			p2 := Rotvec(90, p1, plb)
			p3 := Lerpvec(wscl/Dist3d(p1,p2),p1,p2)
			dx = p3[0]-p1[0]
			dy = p3[1]-p1[1]
			switch{
				case wx > 0.0:
				//arrow from p3 to p1
				x = p3[0]; y = p3[1]
				dx = -dx
				dy = -dy
				case wx < 0.0:
				//arrow from p1 to p3
				x = p1[0]; y = p1[1]
			}
			
			data += fmt.Sprintf("%f %f %f %f %v\n",x,y,dx,dy,lcase)
			if i == 0{
				ldata += fmt.Sprintf("%f %f %.1f %v\n",x+dx/2.0,y+dy/2.0,wx,lcase)
			}
			if i == 10{
				ldata += fmt.Sprintf("%f %f %.1f %v\n",x+dx/2.0,y+dy/2.0,wx,lcase)
			}
			case 3:
			p3 := Lerp3d(wscl, ivec, p1)
			dx = p3[0]-p1[1]
			dy = p3[1]-p1[1]
			dz = p3[2]-p1[2]
			switch{
				case wx > 0.0:
				x = p3[0]; y = p3[1]; z = p3[2]
				dx = -dx; dy = -dy; dz = -dz
				case wx < 0.0:
				x = p1[0]; y = p1[1]; z = p1[2]
			}
			data += fmt.Sprintf("%f %f %f %f %f %f %v\n",x,y,z,dx,dy,dz,lcase)
			if i == 0{
				ldata += fmt.Sprintf("%f %f %f %.1f %v\n",x+dx/2.0,y+dy/2.0,z+dz/2.0,wx,lcase)
			}
			if i == 10{
				ldata += fmt.Sprintf("%f %f %f %.1f %v\n",x+dx/2.0,y+dy/2.0,z+dz/2.0,wx,lcase)
			}
		}
	}
	return
}

func plotltyp5(val, pa, pb, wng []float64, fsx, l, fscale float64)(data, ldata string){
	//point axial load at la
	wa, la := val[2], val[4]
	if len(pa) == 1{
		pa = append(pa, 0.0)
		pb = append(pb, 0.0)
	}
	x0 := fsx/2.0; k := 1.0
	lcase := 1
	if len(val) > 6{lcase = int(val[6])}
	yld := 0.0
	if lcase > 1{
		yld = fsx*float64(lcase-1)
	}
	if len(pa) == 3{yld = 0.0}
	var pla, plb []float64
	if yld > 0.0{
		p1 := Rotvec(90, pa, pb)
		p2 := Rotvec(270,pb, pa)
		pla = Lerpvec(yld/Dist3d(pa,p1), pa, p1)
		plb = Lerpvec(yld/Dist3d(pb, p2), pb, p2)
	} else {
		pla = pa
		plb = pb
	}

	p1 := Lerpvec(la/l,pla,plb)
	var wscl float64
	if fscale > 0.0{
		wscl = math.Abs(wa)/fscale
	} else {
		wscl = slerp(fsx, k, x0, math.Abs(wa))
	}
	
	p2 := Lerpvec(wscl/Dist3d(p1, plb),p1, plb)
	var x, y, z, dx, dy, dz float64
	switch len(pa){
		case 2:
		switch{
			case wa < 0.0:
			//draw arrow from p1 to p2
			x = p1[0]; y = p1[1]
			dx = p2[0] - p1[0]
			dy = p2[1] - p1[1]
			case wa > 0.0:
			//draw arrow from p2 to p1
			x = p2[0]; y = p2[1]
			dx = p2[0] - p1[0]
			dy = p2[1] - p1[1]
			dx = -dx
			dy = -dy
		}
		data += fmt.Sprintf("%f %f %f %f %v\n",x,y,dx,dy,lcase)
		ldata += fmt.Sprintf("%f %f %.1f %v\n",x+dx/2.0,y+dy/2.0,wa,lcase)
		case 3:
		switch{
			case wa > 0.0:
			//draw arrow from p1 to p2
			x = p1[0]; y = p1[1]; z = p1[2]
			dx = p2[0] - p1[0]
			dy = p2[1] - p1[1]
			dz = p2[2] - p1[2]
			case wa < 0.0:
			//draw arrow from p2 to p1
			x = p2[0]; y = p2[1]; z = p2[1]
			dx = p2[0] - p1[0]
			dy = p2[1] - p1[1]
			dz = p2[2] - p1[2]
			dx = -dx
			dy = -dy
			dz = -dz
		}
		data += fmt.Sprintf("%f %f %f %f %f %f %v\n",x,y,z,dx,dy,dz,lcase)
		ldata += fmt.Sprintf("%f %f %f %.1f %v\n",x+dx/2.0,y+dy/2.0,z+dz/2.0,wa,lcase)
	}
	return
}
//plotltyp6 plots a member load of type 6(u.axial load)
func plotltyp6(val, pa, pb, wng []float64, fsx, l, fscale float64)(data, ldata string){
	//uniform axial load w at la to l - lb
	wa, la, lb := val[2], val[4], val[5]
	if len(pa) == 1{
		pa = append(pa, 0.0)
		pb = append(pb, 0.0)
	}
	x0 := fsx/2.0; k := 1.0
	lcase := 1
	if len(val) > 6{lcase = int(val[6])}
	yld := 0.0
	if lcase > 1{
		yld = fsx*float64(lcase-1)
	}
	var pla, plb []float64
	if len(pa) == 3{yld = 0.0}
	if yld > 0.0{
		p1 := Rotvec(90, pa, pb)
		p2 := Rotvec(270,pb, pa)
		pla = Lerpvec(yld/Dist3d(pa,p1), pa, p1)
		plb = Lerpvec(yld/Dist3d(pb, p2), pb, p2)
	} else {
		pla = pa
		plb = pb
	}

	step := (l-lb-la)/10.0
	var wscl float64
	if fscale > 0.0{
		wscl = math.Abs(wa/fscale)
	} else {
		wscl = slerp(fsx, k, x0, math.Abs(wa))
	}
	var x, y, z, dx, dy, dz float64
	for i:=0; i < 11; i++{
		dl := step * float64(i)
		p1 := Lerpvec(dl/l, pla, plb)
		p2 := Lerpvec(wscl/Dist3d(p1,plb), p1, plb)
		switch len(pa){
			case 2:
			switch{
				case wa > 0.0:
				//draw arrow from p1 to p2
				x = p1[0]; y = p1[1]
				dx = p2[0] - p1[0]
				dy = p2[1] - p1[1]
				case wa < 0.0:
				//draw arrow from p2 to p1
				x = p2[0]; y = p2[1]
				dx = p2[0] - p1[0]
				dy = p2[1] - p1[1]
				dx = -dx
				dy = -dy
			}
			data += fmt.Sprintf("%f %f %f %f %v\n",x,y,dx,dy,lcase)
			ldata += fmt.Sprintf("%f %f %.1f %v\n",x+dx/2.0,y+dy/2.0,wa,lcase)
			case 3:
			switch{
				case wa > 0.0:
				//draw arrow from p1 to p2
				x = p1[0]; y = p1[1]; z = p1[2]
				dx = p2[0] - p1[0]
				dy = p2[1] - p1[1]
				dz = p2[2] - p1[2]
				case wa < 0.0:
				//draw arrow from p2 to p1
				x = p2[0]; y = p2[1]; z = p2[1]
				dx = p2[0] - p1[0]
				dy = p2[1] - p1[1]
				dz = p2[2] - p1[2]
				dx = -dx
				dy = -dy
				dz = -dz
			}
			data += fmt.Sprintf("%f %f %f %f %f %f %v\n",x,y,z,dx,dy,dz,lcase)
			ldata += fmt.Sprintf("%f %f %f %.1f %v\n",x+dx/2.0,y+dy/2.0,z+dz/2.0,wa,lcase)
		}

	}
	return
	
}


//MsLoad2Dat returns (2d) member force plot data
func MsLoad2Dat(mod *Model) (data, ldata, mdata string){
	//index 4 member loads
	//var xa, ya, za float64
	
	fsx := 1.0
	mod.Frcscale = 50.0
	switch mod.Units{
		case "kpin","kips","kp-in":
		fsx = 25.0
		mod.Frcscale = 1.0
		case "nmm","n-mm":
		fsx = 1000.0
		mod.Frcscale = 1.0
	}
	if len(mod.Cmdz) > 1{
		switch mod.Cmdz[1]{
			case "kips","kpin":
			fsx = 25.0
			// mod.Frcscale = 10.0
			case "mmks","nmm":
			fsx = 1000.0
		}
	}
	for _, val := range mod.Msloads {
		if len(val) < 6{
			log.Println("error in member load->",val)
			return
		}
		m := int(val[0])
		jb := mod.Mprp[m-1][0]
		je := mod.Mprp[m-1][1]
		pa := mod.Coords[jb-1]
		pb := mod.Coords[je-1]
		lspan := Dist3d(pa,pb)
		ltyp := int(val[1])
		//read angle of roll
		wng := make([]float64,2)
		if len(pa) == 3{
			if len(mod.Wng)>=m{
				wng = mod.Wng[m-1]
			}
		}
		switch ltyp{
			case 1:
			d1, l1 := plotltyp1(val, pa, pb, wng, fsx, lspan, mod.Frcscale)
			data += d1; ldata += l1
			case 2:
			m1, l1 := plotltyp2(val, pa, pb, wng, fsx, lspan,mod.Frcscale)
			mdata += m1; ldata += l1
			case 3:
			d1, l1 := plotltyp3(val, pa, pb, wng, fsx, lspan,mod.Frcscale)
			data += d1; ldata += l1
			case 4:
			d1, l1 := plotltyp4(val, pa, pb, wng, fsx, lspan,mod.Frcscale)
			data += d1; ldata += l1
			case 5:
			d1, l1 := plotltyp5(val, pa, pb, wng, fsx, lspan,mod.Frcscale)
			data += d1; ldata += l1
			case 6:
			d1, l1 := plotltyp6(val, pa, pb, wng, fsx, lspan,mod.Frcscale)
			data += d1; ldata += l1
			case 7:
			//TODO torsional moment?
			case 8:
			//HOWDO temperature changes?
			case 9:
			//FABRICATION ERRORS?
		}
	}
	data += "\n\n"
	return
}
