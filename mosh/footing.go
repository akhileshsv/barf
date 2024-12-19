package barf

//ah.footing funcs, what else

import (
	"fmt"
	"math"
	"time"
	"strings"
	"github.com/olekukonko/tablewriter"
)

//Quant takes off quantities for an rcc footing
func (f *RccFtng) Quant(){
	var cl, tl, vstl float64
	f.Bmap = make(map[float64][][]float64); f.Bsum = make(map[float64]float64)
	//rez = append(rez, []float64{dia, nx, spcx, asx, asx- astx, ny, spcy, asy, asy- asty})
	dia, nx, ny := f.Rbr[0], f.Rbr[1], f.Rbr[5]
	//volume of footing
	switch f.Sloped{
		case true:
		l := f.Eo + f.Colx; w := f.Eo + f.Coly
		f.Vtot = (f.Dused - f.Dmin) * (f.Lx * w + f.Ly * l + 2.0 * (f.Lx * f.Ly + (l+w))) + f.Dmin * f.Lx * f.Ly
		f.Afw = f.Lx * f.Ly + f.Dmin * 2.0 * (f.Lx + f.Ly)
		case false:
		f.Vtot = f.Dused * f.Lx * f.Ly
		f.Afw = f.Lx * f.Ly + f.Dused * 2.0 * (f.Lx + f.Ly)
	}
	ld := BarDevLen(f.Fck, f.Fy, dia)
	cl = f.Lx*1e3 - 2.0 * f.Nomcvr*1e3 + 2.0 * ld; tl = math.Round(cl * nx) 
	f.Bmap[dia] = append(f.Bmap[dia],[]float64{tl,nx,cl,1.0})
	f.Bsum[dia] += tl
	vstl += tl * RbrArea(dia) * 1e-9
	//fmt.Println("tl main",tl)
	cl = f.Ly * 1e3 - 2.0 * f.Nomcvr* 1e3 + 2.0 * ld; tl = math.Round(cl * ny) 
	f.Bmap[dia] = append(f.Bmap[dia],[]float64{tl,ny,cl,2.0})
	f.Bsum[dia] += tl
	//fmt.Println("tl dist",tl)
	vstl += tl * RbrArea(dia) * 1e-9
	f.Vrcc = f.Vtot
	f.Wstl = vstl * 7850.0
	switch len(f.Kostin){
		case 3:
		f.Kost = f.Vrcc * f.Kostin[0] + f.Afw * f.Kostin[1] + f.Wstl * f.Kostin[2]
		default:
		f.Kost = f.Vrcc * CostRcc + f.Afw * CostForm + f.Wstl * CostStl
	}
	f.Kunit = f.Kost/f.Lx/f.Ly	
	return
}

//Table prints good ol' ascii table reports for an rcc footing
func (f *RccFtng) Table(printz bool){
	rezstr := new(strings.Builder)
	t := "pad"; if f.Sloped{t = "sloped"}
	hdr := fmt.Sprintf("%s\nrcc footing report\ndate-%s\n%s\n",ColorYellow,time.Now().Format("2006-01-02"),ColorReset)
	rezstr.WriteString(hdr)
	/*
	hdr = ""
	hdr += fmt.Sprintf("%s\n%s footing\ndims- lx %.1f, ly %.1f mm\n",f.Title,t, f.Lx*1e3, f.Ly*1e3)
	hdr += fmt.Sprintf("grade of concrete M %.1f\nsteel - Fe %.f\n", f.Fck, f.Fy)
	hdr += fmt.Sprintf("col dims - x %0.1f mm y %0.1f mm offset %.01f mm\n", f.Colx*1e3, f.Coly*1e3, f.Eo * 1e3)
	hdr += fmt.Sprintf("cover - nominal %0.1f mm effective %0.1f mm\n", f.Nomcvr*1e3, f.Efcvr*1e3)
	*/
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)
	table.SetCaption(true,"footing specs")
	table.SetHeader([]string{"type","concrete","steel","nom.cvr\n(mm)","eff.cvr\n(mm)"})
	var row string
	row = fmt.Sprintf("%s,M%.f,Fe%.f,%.2f,%.2f",t,f.Fck,f.Fy,f.Nomcvr*1e3,f.Efcvr*1e3)
	table.Append(strings.Split(row, ","))
	table.Render()

	if f.Dz{
		//hdr += fmt.Sprintf("%s\nfooting depth %.f mm\n",ColorYellow,f.Dused*1e3)
		//if f.Sloped{
		//	hdr += fmt.Sprintf("min.edge depth %.f mm\n",f.Dmin*1e3)
		//}
		//rezstr.WriteString(hdr)
		//rezstr.WriteString(ColorReset)
		rezstr.WriteString(ColorPurple)
		
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true, "footing geometry")
		table.SetHeader([]string{"col x\n(mm)","col y\n(mm)","edge offset\n(mm)","lx\n(mm)","ly\n(mm)","depth\n(mm)","min.edge depth\n(mm)"})
		row = fmt.Sprintf("%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f",f.Colx*1e3,f.Coly*1e3,f.Eo*1e3,f.Lx*1e3,f.Ly*1e3,f.Dused*1e3,f.Dmin*1e3)
		table.Append(strings.Split(row, ","))
		table.Render()
		
		rezstr.WriteString(ColorRed)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true,"design loads")
		table.SetHeader([]string{"load type","axial load\n(kn)","bm-x\n(kn-m)","bm-y\n(kn.m)","psf"})
		for i, pu := range f.Pus{
			switch i{
				case 0:
				t = "dead"
				case 1:
				t = "live"
				case 2:
				t = "wind"
			}
			row = fmt.Sprintf("%s,%.3f,%.3f,%.3f,%.2f",t,pu,f.Mxs[i],f.Mys[i],f.Psfs[i])
			table.Append(strings.Split(row,","))
		}
		table.Render()
		rezstr.WriteString(ColorBlue)
		table = tablewriter.NewWriter(rezstr)
		table.SetCaption(true,"reinforcement")
		//rez = append(rez, []float64{dia, nx, spcx, asx, asx- astx, ny, spcy, asy, asy- asty})
		table.SetHeader([]string{"loc","bm\n(kn-m)","shear\n(kn)","p.shear\n(kn)","dia\n(mm)"," spacing\n(mm)","ast prov\n(mm2/m)","delta\n(mm2/m)"})
		row = fmt.Sprintf("%s, %.2f, %.2f, %.2f, %.0f, %.0f, %.0f, %.0f","x",f.Mux,f.Vux,f.Vp,f.Rbr[0],f.Rbr[2],f.Rbr[3],f.Rbr[4])
		table.Append(strings.Split(row,","))
		row = fmt.Sprintf("%s, %.2f, %.2f, %.2f, %.0f, %.0f, %.0f, %.0f","y",f.Muy,f.Vuy,f.Vp,f.Rbr[0],f.Rbr[6],f.Rbr[7],f.Rbr[8])
		
		table.Append(strings.Split(row,","))
		table.Render()
		rezstr.WriteString(ColorCyan)
		table = tablewriter.NewWriter(rezstr)
		table.SetHeader([]string{"vol tot(m3)","vol rcc(m3)","wt stl(kg)","form area (m2)","cost (rs)","unit cost(rs/m2)"})
		table.SetCaption(true,"quantity take off")
		row = fmt.Sprintf("%.3f, %.3f, %.3f, %.3f, %.f, %.2f\n",f.Vtot,f.Vrcc,f.Wstl,f.Afw, f.Kost, f.Kunit)
		table.Append(strings.Split(row,","))
		table.Render()
		rezstr.WriteString(ColorReset)
	}
	f.Report = fmt.Sprintf("%s",rezstr)
	if printz{
		fmt.Println(f.Report)
	}
}

//Printz prints a footing
//printz printz!
func (f *RccFtng) Printz() (rez string){
	t := "pad"; if f.Sloped{t = "sloped"}
	
	rez += fmt.Sprintf("%s\n%s footing\ndims- lx %.1f, ly %.1f mm\n",f.Title,t, f.Lx*1e3, f.Ly*1e3)
	rez += fmt.Sprintf("grade of concrete M %.1f\nsteel - Fe %.f\n", f.Fck, f.Fy)
	rez += fmt.Sprintf("col dims - x %0.1f mm y %0.1f mm offset %.01f mm\n", f.Colx*1e3, f.Coly*1e3, f.Eo * 1e3)
	rez += fmt.Sprintf("cover - nominal %0.1f mm effective %0.1f mm\n", f.Nomcvr*1e3, f.Efcvr*1e3)

	

	if f.Dz{
		rez += fmt.Sprintf("footing depth %.f mm\n",f.Dused*1e3)
		if f.Sloped{
			rez += fmt.Sprintf("min.edge depth %.f mm\n",f.Dmin*1e3)
		}
		rez += fmt.Sprintf("mux, vux, vp - %.2f knm, %.2f knm, %.2f kn\n",f.Mux,f.Vux,f.Vp)
		rez += fmt.Sprintf("rbr x- dia %.f mm nx %.f nos spcx %.f mm asx %.f mm2 asx- astx %.f mm2\n",f.Rbr[0],f.Rbr[1],f.Rbr[2],f.Rbr[3],f.Rbr[4])
		rez += fmt.Sprintf("muy, vuy, vp - %.2f, %.2f, %.2f\n",f.Muy,f.Vuy,f.Vp)
		rez += fmt.Sprintf("rbr y- dia %.f mm ny %.f nos spcy %.f mm asy %.f mm2 asy- asty %.f mm2\n",f.Rbr[0],f.Rbr[5],f.Rbr[6],f.Rbr[7],f.Rbr[8])
	}
	return
}
