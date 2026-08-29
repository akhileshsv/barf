package barf

import (
	"fmt"
	draw "barf/draw"
)

//DrawDac draws a double angle connection
func (b *Blt) DrawDac(err error){
	return
}

//Draw draws a bolt group
func (b *Blt) Draw()(err error){
	switch b.Name{
		case "dac":
		var data, cdata, ldata, d1 string
		switch{
			case len(b.Dims) < 3:
			err = fmt.Errorf("invalid dims for cleat angles %.f",b.Dims)
			return
			case len(b.Mdims) < 2:
			err = fmt.Errorf("member dimensions not specified %.f",b.Mdims)
			return
			case len(b.Mdims[0]) < 4 || len(b.Mdims[1]) < 4:
			err = fmt.Errorf("invalid member dimensions for i section%.f",b.Mdims)
			return
			case len(b.Mdxs) < 2:
			b.Mdxs = []string{"m1","m2"}
		}
		//side view first (opposite day)
		//draw angle
		data += DrawRectView(3, b.Dims[0], []float64{0,b.Dims[0]/2.0},[]float64{b.Dims[1],b.Dims[0]/2.0})
		//support leg of angle
		data += DrawRectView(3, b.Dims[0], []float64{0-b.Dims[2],b.Dims[0]/2.0},[]float64{0,b.Dims[0]/2.0})
		ldata += fmt.Sprintf("%f -25.0 %.f\n",b.Dims[1]/2.0,b.Dims[1])
		ldata += fmt.Sprintf("%f %f %.f\n",b.Dims[1] + 25.0,b.Dims[0]/2.0,b.Dims[0])
		
		cx := -b.Dims[2]
		cy := b.Dims[0]/2.0
		//main mem web rect (500 mm long)
		h := b.Mdims[0][1]
		tf := b.Mdims[0][2]
		dw := h - 2.0 * tf
		ye := b.Dims[0] + 50.0
		yb := ye - dw
		yb = (ye + yb)/2.0
		xb := 10.0
		xe := 510.0
		data += DrawRectView(2, dw, []float64{xb,yb},[]float64{xe,yb})
		//member label
		ldata += fmt.Sprintf("%f %f %s\n",(xb+xe)/2.0,(yb+ye)/2.0,b.Mdxs[0])
		
		ye = b.Dims[0] + 50.0 + tf
		yb = ye - h
		yb = (ye + yb)/2.0
		data += DrawRectView(2, h, []float64{xb,yb},[]float64{xe,yb})
		//draw plate thickness
		switch b.Ctyp{
			case 1:
			//beam-col
			switch b.Cloc{
				case "flange":
				
				h := b.Mdims[1][1]
				tf := b.Mdims[1][2]
				dw := h - 2.0 * tf
				//draw flange rect 500 mm high
				csec := SecGen(1, []float64{h,500.0})
				tx := cx - csec.Prop.Xc - h/2.0
				ty := cy - csec.Prop.Yc
				csec = SecTranslate(csec, tx, ty)
				csec.Draw("")
				data += csec.Data[0]				
				//member label
				ldata += fmt.Sprintf("%f %f %s\n",csec.Prop.Xc,csec.Prop.Yc,b.Mdxs[1])

				//draw web rect
				csec = SecGen(1, []float64{dw,500.0})
				tx = cx - csec.Prop.Xc - dw/2.0 - tf
				ty = cy - csec.Prop.Yc
				csec = SecTranslate(csec, tx, ty)
				csec.Draw("")
				data += csec.Data[0]				
				
				case "web":
				//just draw col web
				tw := b.Mdims[1][3]
				
				csec := SecGen(1, []float64{tw,500.0})
				tx := cx - csec.Prop.Xc - tw/2.0
				ty := cy - csec.Prop.Yc
				csec = SecTranslate(csec, tx, ty)
				csec.Draw("")
				
				ldata += fmt.Sprintf("%f %f %s\n",csec.Prop.Xc,csec.Prop.Yc,b.Mdxs[1])
				data += csec.Data[0]								
			}
			case 2:
			//beam-beam
			//draw sec
			bsec := SecGen(12,b.Mdims[1])
			tx := cx - bsec.Prop.Xc - b.Mdims[1][2]/2.0
			ty := cy - bsec.Prop.Yc
			bsec = SecTranslate(bsec, tx, ty)
			bsec.Draw("")
			
			ldata += fmt.Sprintf("%f %f %s\n",bsec.Prop.Xc,bsec.Prop.Yc,b.Mdxs[1])
			data += bsec.Data[0]
		}
		xs := b.Edged
		ys := b.Edged
		for j := 0; j < b.Nj; j++{
			if j == 0{
				ldata += fmt.Sprintf("%f %f %.f\n",xs,ys-b.Edged/2.0-5.0,b.Edged)
			}
			for i := 0; i < b.Ni; i++{
				cdata += fmt.Sprintf("%f %f %f\n",xs,ys,b.Dia/2.0)
				ldata += fmt.Sprintf("%f %f M%.f\n",xs,ys,b.Dia)
				if i == 0{
					ldata += fmt.Sprintf("%f %f %.f\n",xs,ys-b.Edged/2.0-5.0,b.Edged)
				} else {
					ldata += fmt.Sprintf("%f %f %.f\n",xs,ys-b.Pitch/2.0,b.Pitch)
				}
				ys += b.Pitch
			}
			xs += float64(j) * b.Gauge
			ys = b.Edged	
		}
		
		data += "\n\n"; cdata += "\n\n"; ldata += "\n\n"
		data = cdata + data
		data += ldata

		d1 = data
		
		data = ""; cdata = ""; ldata = ""
		
		//front view
		//draw main member (section view)
		bsec := SecGen(12,b.Mdims[0])
		bsec.Draw("")
		data += bsec.Data[0]
		
		ldata += fmt.Sprintf("%f %f %s\n",bsec.Prop.Xc,bsec.Prop.Yc,b.Mdxs[0])

		cx = bsec.Prop.Xc
		cy = bsec.Prop.Yc
		//add plate thickness
		p1 := []float64{cx-tf/2.0,cy}
		p2 := []float64{cx-tf/2.0 -b.Dims[2],cy}
		p3 := []float64{cx+tf/2.0, cy}
		p4 := []float64{cx+tf/2.0+b.Dims[2],cy}
		
		data += DrawRectView(3,b.Dims[0], p1, p2)
		data += DrawRectView(3,b.Dims[0], p3, p4)
		
		//draw supporting mem cleats
		p1[0] = cx - tf -b.Dims[2]
		p2[0] = p1[0] - b.Dims[1]
		p3[0] = cx + tf + b.Dims[2]
		p4[0] = p3[0] + b.Dims[1]

		data += DrawRectView(3,b.Dims[0], p1, p2)
		data += DrawRectView(3,b.Dims[0], p3, p4)

		ldata += fmt.Sprintf("%f %f %s\n",p2[0]-50.0,p4[1],b.Mdxs[1])
		
		ldata += fmt.Sprintf("%f %f %s\n",p4[0]+50.0,p4[1],b.Mdxs[1])
		//draw supporting mem bolts (b.Nb2/2 on either side)
		nb := b.Nb2/2
		
		if b.Dims[0] < 2.0 * b.Edged + float64(nb - 1) * b.Pitch{
			err = fmt.Errorf("insufficent depth (%.f) for bolts (req %.f) in supporting mem",b.Dims[0],2.0 * b.Edged + float64(nb - 1) * b.Pitch)
		}
		
		//left
		xs = (p2[0] + p1[0])/2.0
		ys = p2[1] - b.Dims[0]/2.0 + b.Edged
		
		for i := 0; i < nb; i++{
			cdata += fmt.Sprintf("%f %f %f\n",xs,ys,b.Dia/2.0)
			ldata += fmt.Sprintf("%f %f M%.f\n",xs,ys,b.Dia)
			if i == 0{
				ldata += fmt.Sprintf("%f %f %.f\n",xs,ys-b.Edged/2.0-5.0,b.Edged)
			} else {
				ldata += fmt.Sprintf("%f %f %.f\n",xs,ys-b.Pitch/2.0,b.Pitch)
			}
			ys += b.Pitch
		}
		
		//right

		xs = (p3[0] + p4[0])/2.0
		ys = p3[1] - b.Dims[0]/2.0 + b.Edged
		
		for i := 0; i < nb; i++{
			cdata += fmt.Sprintf("%f %f %f\n",xs,ys,b.Dia/2.0)
			ldata += fmt.Sprintf("%f %f M%.f\n",xs,ys,b.Dia)
			if i == 0{
				ldata += fmt.Sprintf("%f %f %.f\n",xs,ys-b.Edged/2.0-5.0,b.Edged)
			} else {
				ldata += fmt.Sprintf("%f %f %.f\n",xs,ys-b.Pitch/2.0,b.Pitch)
			}
			ys += b.Pitch
		}

		
		data += "\n\n"; cdata += "\n\n"; ldata += "\n\n"
		data = cdata + data + ldata
		data = d1 + data
		if b.Web{
			b.Folder = "web"
		}
		var txtplot string
		txtplot, err = draw.Draw(data, "plotconn.gp", b.Term, b.Folder, b.Title, b.Title, "","","")
		if err != nil{
			return
		}
		if b.Term == "dumb"{fmt.Println(txtplot)}
		case "fep":
		
		case "fp":
	}
	return
}
