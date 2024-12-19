package barf

import (
	"log"
	"fmt"
	"math"
)

//plotfrc plots a (single) force vector as a vector :) in gnuplot
func plotfrc(pt []float64, fsx, frc float64, frmstr string, ax, ldc int) (data, ldata string){
	x0 := fsx/2.0; k := 1.0
	dw := slerp(fsx, k, x0, math.Abs(frc))
	var xs, ys, zs, dx, dy, dz float64
	//log.Println("ldcase fac->",ldc)
	for i, v := range pt{
		switch i{
			case 0:
			xs = v
			case 1:
			ys = v
			case 2:
			zs = v
		}
	}
	if ldc == 0{ldc = 1}
	switch ax{
		case 1:
		//frc x
		//add loadcase factor
		xs += float64(ldc-1)*fsx
		switch{
			case frc > 0.0:
			//arrow from left(xs-dw->xs)
			xs -= dw; dx = dw
			case frc < 0.0:
			//arrow towards left(xs->xs-dw)
			dx = -dw
		}
		case 2:
		//frc y
		ys += float64(ldc-1)*fsx		
		switch{
			case frc > 0.0:
			//arrow from bot(ys->ys+dw)
			dy = dw
			case frc < 0.0:
			//arrow from top(ys+dw->ys)
			ys += dw; dy = -dw
			
		}
		case 3:
		//frc z
		zs += float64(ldc-1)*fsx		
		switch{
			case frc > 0.0:
			//arrow from ? zs -> zs + dw
			dz = dw
			case frc < 0.0:
			//arrow from zs + dw -> zs
			zs += dw; dz = -dw
		}
		
	}
	//log.Println("xs, ys, zs, dx, dy, dz->",xs, ys, zs, dx, dy, dz)
	switch frmstr{
		case "1db","2dt","2df":
		data = fmt.Sprintf("%f %f %f %f %v\n",xs,ys,dx,dy,ldc+4)
		ldata = fmt.Sprintf("%f %f %.2f 2.0\n",xs+dx/2.0,ys+dy/2.0,frc)
		case "3dt", "3dg", "3df":			
		data = fmt.Sprintf("%f %f %f %f %f %f %v\n",xs,ys,zs,dx,dy,dz,ldc+4)
		ldata = fmt.Sprintf("%f %f %f %.2f 2.0\n",xs+dx/2.0,ys+dy/2.0,zs+dz/2.0,frc)
	}
	//log.Println("ldata->",ldata)
	return
}

//plotmz plots a moment as a circle/point(?) in gnuplot
func plotmz(pt []float64, fsx, frc float64, frmstr string, ax, ldc int)(mdata, ldata string){
	x0 := fsx/2.0; k := 1.0
	dw := slerp(fsx, k, x0, math.Abs(frc))
	//log.Println("pt->",pt,"mz->",frc, "dw->",dw)
	var xs, ys, zs float64
	for i, v := range pt{
		switch i{
			case 0:
			xs = v
			case 1:
			ys = v
			case 2:
			zs = v
		}
	}
	switch ax{
		// case 1:
		// //mx - circle in yz plane
		// case 2:
		//my - circle in xz plane
		case 1,2,3:
		//mz - circle in xy plane
		switch frmstr{
			case "1db","2dt","2df":
			mdata = fmt.Sprintf("%f %f %f 3.0\n",xs,ys,dw)
			ldata = fmt.Sprintf("%f %f %v 1.0\n",xs+dw/2.0,ys+dw/2.0,frc)
			default:
			mdata = fmt.Sprintf("%f %f %f %f 3.0\n",xs,ys,zs,dw)
			ldata = fmt.Sprintf("%f %f %f %v\n",xs+dw/2.0,ys+dw/2.0,zs+dw/2.0,frc)
		}
	}
	return
}

//plotfrmvec returns data for a single force array in gnuplot
func plotfrcvec(pt, fvec []float64, fsx float64, frmstr string)(data, ldata, mdata string){
	//plotfrc(pt []float64, fsx, frc float64, frmstr string, ax, ldc int)(data, ldata string){
	var ax, ldc int
	ldc = 1
	switch frmstr{
		case "1db","2dt":
		//len(fvec) = 3 (joint, fy/fx, mz/fx)
		if len(fvec) > 3{
			ldc = int(fvec[3])
		}
		for i, v := range fvec[1:]{
			switch i{
				case 0:
				switch frmstr{
					case "1db":
					//fy
					ax = 2
					default:
					//fx
					ax = 1
				}
				if v != 0.0{
					d1, l1 := plotfrc(pt, fsx, v, frmstr, ax, ldc)
					data += d1; ldata += l1	
				}
				case 1:
				switch frmstr{
					case "1db":
					//mz
					if v != 0.0{
						m1, l1 := plotmz(pt, fsx, v, frmstr, ax, ldc)
						mdata += m1; ldata += l1
					}
					default:
					//fy
					ax = 2
					if v != 0.0{						
						d1, l1 := plotfrc(pt, fsx, v, frmstr, ax, ldc)
						data += d1; ldata += l1	
					}
				}
			}
		}
		case "2df", "3dt", "3dg":
		//joint, fx/fx/fy, fy/fy/mx, mz/fz/mz
		if len(fvec) > 4{
			ldc = int(fvec[4])
		}
		for i, v := range fvec[1:]{
			switch i{
				case 0:
				switch frmstr{
					case "2df","3dt":
					//fx
					ax = 1
					default:
					//fy
					ax = 2
				}
				d1, l1 := plotfrc(pt, fsx, v, frmstr, ax, ldc)
				data += d1; ldata += l1
				case 1:
				switch frmstr{
					case "2df","3dt":
					//fy
					ax = 2
					d1, l1 := plotfrc(pt, fsx, v, frmstr, ax, ldc)
					data += d1; ldata += l1
					default:
					//mx(torsion)
					m1, l1 := plotmz(pt, fsx, v, frmstr, ax, ldc)
					mdata += m1; ldata += l1
				}
				case 2:
				switch frmstr{
					case "2df","3dg":
					//mz
					m1, l1 := plotmz(pt, fsx, v, frmstr, ax, ldc)
					mdata += m1; ldata += l1
					case "3dt":
					//fz
					ax = 3
					d1, l1 := plotfrc(pt, fsx, v, frmstr, ax, ldc)
					data += d1; ldata += l1
				}
			}
		}
		case "3df":
		if len(fvec) > 7{
			ldc = int(fvec[7]) 
		}
		for i, v := range fvec[1:]{
			switch i{
				case 0,1,2:
				//fx, fy, fz
				ax = i+1
				d1, l1 := plotfrc(pt, fsx, v, frmstr, ax, ldc)
				data += d1; ldata += l1	
				default:
				//mx, my, mz
				ax = i-2
				m1, l1 := plotmz(pt, fsx, v, frmstr, ax, ldc)
				mdata += m1; ldata += l1
			}
		}
	}
	return
}

//JLoadDat returns nodal force plot data
func JloadDat(mod *Model) (data, ldata, mdata string){
	fsx := 1.0
	//l - max value, k - steepness , x0 - midpoint, x - is just "x"
	//func slerp(l, k, x0, x float64)
	switch mod.Units{
		case "kp-in","kpin":
		fsx = 25.0
		case "n-mm","nmm":
		fsx = 1000.0
	}
	if len(mod.Cmdz) > 1{
		switch mod.Cmdz[1]{
			case "kips","kpin","kp-in":
			fsx = 25.0
			case "mmks","nmm","n-mm":
			fsx = 1000.0
		}
	}
	//log.Println("frmstr, fsx->",mod.Frmstr, fsx)
	//index 3 joint loads
	for i, val := range mod.Jloads{
		if len(val) < 3{
			log.Printf("error in nodal load no %v -> %v\n",i,val)
			continue
		}
		pt := mod.Coords[int(val[0])-1]
		d1, l1, m1 := plotfrcvec(pt, val, fsx, mod.Frmstr)
		data += d1; ldata += l1; mdata += m1
	}
	return
}
