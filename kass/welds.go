package barf

import (
	"fmt"
	"math"
	"math/rand"
	"errors"
	"strings"
	"github.com/olekukonko/tablewriter"
)

//Wld is a struct to store weld group data
//see harrison program WELDSS, section 4.3
type Wld struct{
	Sizes  []float64 //fillet sizes of weld 
	Wcs    [][]float64 //weld coordinates start x, y, end x, y
	Frcs   [][]float64 //applied force - fx, fy, mz
	Fcs    [][]float64 // coordinates of point of force (x, y)
	Xns    []float64
	Yns    []float64
	Res    []float64
	Strs   []float64
	Report string
	Title  string
	Msname string  //base member section name
	Shop   bool    //is shop welded/site welded 
	Id     int
	Styp   int     //member section type
	Styps  []int   //if a slice is actually needed?
	Code   int     //1-dil hai 2-bs code
	Ctyp   int     //connection type - 1 - fillet , 2 - butt/groove 
	Ltyp   int     //load type - 1 - concentric (T/C), 2 - 
	Wltyp  int     //weld layout type - (1 - [ - | , 2 - [ , 3 - [] etc) 
	Grade  float64 //
	Ang    float64 //angle between fusion faces 
	Ks     float64 //K*S = eff. throat thickness
	Size   float64 //(uniform) weld size
	Lweld  float64 //net length of weld
	Ert    float64 //end return of weld
	Lop    float64 //overlap length
	L1, L2 float64 //lengths of weld (bottom, top)
	Tp     float64 //design plate thickness
	Tmem   float64 //thickness of smallest mem
	Dmem   float64 //depth of member/joint (the | in [ - 2lop + dmem, etc)
	Tms    []float64 //connected mem. thicknesses
	Dims   [][]float64 //connected mem dims
	Pu     float64 //design (ult.) axial force 
	Mu     float64 //design moment
	Vu     float64 //design shear
	Ymf    float64 //partial factor of safety
	Fu     float64 //ultimate stress 
	Fy     float64 //yield stress of parent metal
	Fyw    float64 //yield stress of weld
	Pdw    float64 //design strength of weld
}


//WeldSs performs analysis of a weld group given a set applied forces and weld geometry
//returns stresses at each start, stop point of each weld
func WeldSs(w *Wld) (err error){
	var a, b, c, d, e, h, v, hy, vx, sm, t, sfx, sfy, smt, cmx, cmy float64
	lenz := []float64{}
	w.Xns = make([]float64, len(w.Sizes) * 2)
	w.Yns = make([]float64, len(w.Sizes) * 2)
	w.Res = make([]float64, len(w.Sizes) * 2)
	w.Strs = make([]float64, len(w.Sizes) * 2)
	//kompute stress resultants
	for i, wc := range w.Wcs{
		lenz = append(lenz, math.Sqrt(math.Pow(wc[0] - wc[2],2) + math.Pow(wc[1]-wc[3],2)))
		a += w.Sizes[i] * lenz[i]
		b += w.Sizes[i] * lenz[i] * (wc[0] + wc[2])/2.0
		c += w.Sizes[i] * lenz[i] * (wc[1] + wc[3])/2.0
		t = math.Pow(wc[0],2) + wc[0] * wc[2] + math.Pow(wc[2],2)
		d += w.Sizes[i] * lenz[i] * t/3.0
		t = math.Pow(wc[1],2) + wc[1] * wc[3] + math.Pow(wc[3],2)
		e += w.Sizes[i] * lenz[i] * t/3.0
	}
	den := math.Pow(b,2) + math.Pow(c, 2) - a * (d + e)
	for i, frc := range w.Frcs{
		h += frc[0]
		hy += frc[0] * w.Fcs[i][1]
		v += frc[1]
		vx += frc[1] * w.Fcs[i][0]
		sm += frc[2]
	}
	theta := (a * hy - a * vx + a * sm + b * v - c * h)/den
	for i, wc := range w.Wcs{
		j := 2 * (i + 1) - 1
		w.Xns[j-1] = w.Sizes[i] * (-h/a + theta * (wc[1] - c/a))
		w.Xns[j] = w.Sizes[i] * (-h/a + theta * (wc[3] - c/a))
		w.Yns[j-1] = w.Sizes[i] * (-v/a + theta * (b/a - wc[0]))
		w.Yns[j] = w.Sizes[i] * (-v/a + theta * (b/a - wc[2]))
		w.Res[j-1] = math.Sqrt(math.Pow(w.Xns[j-1],2)+math.Pow(w.Yns[j-1],2))
		w.Res[j] = math.Sqrt(math.Pow(w.Xns[j],2)+math.Pow(w.Yns[j],2))
		w.Strs[j-1] = w.Res[j-1]/(0.7071 * w.Sizes[i])
		w.Strs[j] = w.Res[j]/(0.7071 * w.Sizes[i])
	}

	//equilibrium check
	for i, wc := range w.Wcs{
		j := 2 * (i + 1) - 1
		rmfx := lenz[i] * (w.Xns[j-1] + w.Xns[j])/2.0
		rmfy := lenz[i] * (w.Yns[j-1] + w.Yns[j])/2.0
		sfx += rmfx
		sfy += rmfy
		if wc[1] != wc[3]{
			t1 := (w.Xns[j-1] - w.Xns[j]) * lenz[i]/3.0
			t2 := math.Pow(wc[1],2) + wc[1] * wc[3] + math.Pow(wc[3],2)
			t3 := (w.Xns[j-1] * wc[3] - w.Xns[j] * wc[1]) * lenz[i]
			t4 := (wc[1] + wc[3])/2.0
			t5 := wc[1] - wc[3]
			cmx = (t1 * t2 - t3 * t4)/t5
		} else {
			cmx = lenz[i] * wc[1] * (w.Xns[j-1] + w.Xns[j])/2.0
		}
		if wc[0] != wc[2]{
			t1 := (w.Yns[j-1] - w.Yns[j]) * lenz[i]/3.0
			t2 := math.Pow(wc[0],2) + wc[0] * wc[2] + math.Pow(wc[2],2)
			t3 := (w.Yns[j-1] * wc[2] - w.Yns[j] * wc[0]) * lenz[i]
			t4 := (wc[0] + wc[2])/2.0
			t5 := wc[0] - wc[2]
			cmy = (t1 * t2 - t3 * t4)/t5
		} else {
			cmy = lenz[i] * wc[0] * (w.Yns[j-1] + w.Yns[j])/2.0
		}
		smt += cmx - cmy
	}

	//build report

	rstr := &strings.Builder{}
	rstr.WriteString(ColorGreen)
	t0 := tablewriter.NewWriter(rstr)
	t0.SetHeader([]string{"weld","fillet size","x start","y start","x end","y end","length"})
	if w.Title == ""{
		if w.Id == 0{
			w.Id = rand.Intn(666)
		}
		w.Title = fmt.Sprintf("no. %v",w.Id)
	}
	cptn := fmt.Sprintf("weld group geometry-%s",w.Title)
	t0.SetCaption(true, cptn)
	for i, wc := range w.Wcs{
		row := fmt.Sprintf("%v, %.f, %.f, %.f, %.f, %.f, %.2f",i+1, w.Sizes[i],wc[0], wc[1], wc[2], wc[3], lenz[i])
		t0.Append(strings.Split(row,","))
	}
	t0.Render()
	t0 = tablewriter.NewWriter(rstr)
	t0.SetHeader([]string{"force","x-coord","y-coord","x-force","y-force","moment"})
	t0.SetCaption(true, "applied forces")
	for i, frc := range w.Frcs{
		row := fmt.Sprintf("%v, %.f, %.f, %f, %f, %f",i+1,w.Fcs[i][0], w.Fcs[i][1], frc[0], frc[1], frc[2])
		t0.Append(strings.Split(row,","))
	}
	t0.Render()
	rstr.WriteString(ColorCyan)
	t1 := tablewriter.NewWriter(rstr)
	t1.SetHeader([]string{"weld","x coord","y coord","force/length","stress"})
	t1.SetCaption(true, "weld group analysis results")
	t0.SetCaption(true, "weld group geometry")
	
	for i, wc := range w.Wcs{
 		j := 2 * (i + 1) - 1
		row := fmt.Sprintf("%v, %.f, %f, %f, %f",i+1, wc[0], wc[1], w.Res[j-1], w.Strs[j-1])
		t1.Append(strings.Split(row,","))
		row = fmt.Sprintf("%v, %.f, %f, %f, %f",i+1, wc[2], wc[3], w.Res[j], w.Strs[j])
		t1.Append(strings.Split(row,","))
	}
	t1.Render()
	rstr.WriteString(ColorPurple)
	t1 = tablewriter.NewWriter(rstr)
	t1.SetHeader([]string{"action","sum-weld","sum-applied","ok?"})
	t1.SetCaption(true, "equilibrium check")
	xchk := math.Abs(sfx + h) < 1e-3
	row := fmt.Sprintf("%s, %f, %f, %t","x-force",sfx, h, xchk)
	t1.Append(strings.Split(row, ","))
	ychk := math.Abs(sfy + v) < 1e-3
	row = fmt.Sprintf("%s, %f, %f, %t","y-force",sfy, v, ychk)
	t1.Append(strings.Split(row, ","))
	mchk := math.Abs(smt + sm + hy - vx) < 1e-3
	row = fmt.Sprintf("%s, %f, %f, %t","moment-origin",smt, sm + hy - vx, mchk)
	t1.Append(strings.Split(row, ","))
	t1.Render()
	rstr.WriteString(ColorReset)
	w.Report = rstr.String()
	if !xchk || !ychk || !mchk{
		return errors.New("equilibrium check failed")
	}
	return nil
}

//Init inits default weld values(size, ym)
func (w *Wld) Init()(err error){
	if w.Code != 1{w.Code = 1}
	if w.Grade == 0.0{w.Grade = 410.0}
	switch w.Grade{
		case 410.0:
		w.Fu = 410.0
		w.Fy = 250.0
		w.Fyw = 250.0
	}
	w.Ymf = 1.5
	if w.Shop{
		w.Ymf = 1.25
	}
	var wmin, wmax float64
	//set weld size
	if w.Tmem > 0.0{
		tmem := w.Tmem
		if w.Tp > tmem{
			tmem = w.Tp
		}
		switch{
			case tmem <= 10.0:
			wmin = 3.0
			case tmem <= 20.0:
			wmin = 5.0
			case tmem <= 32.0:
			wmin = 6.0
			case tmem <= 50:
			//what is 8 first run?
			wmin = 10.0
			default:
			err = fmt.Errorf("super thick member do normal rules apply?")
			return
		}
		wmax = w.Tmem - 1.5
		if w.Tp - 1.5 < wmax{
			wmax = w.Tp - 1.5
		}
	} else {		
		tmin := w.Tmem
		wmin = 3.0
		
		tmems := append(w.Tms, w.Tp)
		for _ , tmem := range tmems{
			if tmin == 0.0{
				tmin = tmem
			} else if tmin > tmem{
				tmin = tmem
			}
			//max size - tmem - 1.5 mm, or 3/4 * tmem for rounded steel sections
			if wmax < tmem - 1.5{
				wmax = tmem - 1.5
			}
			//table 6.4 duggal(is the barf is800) - min. weld sizes
			switch{
				case tmem <= 10.0:
				wmin = 3.0
				case tmem <= 20.0:
				wmin = 5.0
				case tmem <= 32.0:
				wmin = 6.0
				case tmem <= 50:
				//what is 8 first run?
				wmin = 10.0
				default:
				err = fmt.Errorf("super thick member do normal rules apply?")
				return
			}
		}
	}
	fmt.Println("minimum and maximum size of weld in mm-", wmin, wmax)
	if wmax < wmin{
		err = fmt.Errorf("wmax %f < wmin %f something is wreng", wmax, wmin)
		return
	}
	switch{
		case w.Size == 0.0:
		w.Size = wmin
		case w.Size < wmin:
		w.Size = wmin
		case w.Size > wmax:
		w.Size = wmax
	}
	if w.Size == 0.0{
		err = fmt.Errorf("weld size %f /member thickness %f, %f not specified", w.Size, w.Tmem, w.Tms)
		return
	}
	switch{
		//what about "fillet weld not reco if wang < 60 deg? page 234 duggal"
		case w.Ang <= 90.0:
		w.Ks = 0.7
		case w.Ang <= 100.0:
		w.Ks = 0.65
		case w.Ang <= 106.0:
		w.Ks = 0.6
		case w.Ang <= 113.0:
		w.Ks = 0.55
		case w.Ang <= 120.0:
		w.Ks = 0.5
		default:
		err = fmt.Errorf("excessive angle for fillet weld-%f",w.Ang)
	}
	return
}

//WeldCalc calculates the design strength of a weld
func WeldCalc(w *Wld)(err error){
	err = w.Init()
	if err != nil{
		fmt.Println("ERRORE, errore",err)
		return
	}
	switch w.Ctyp{
		case 1:
		//fillet weld
		switch w.Ltyp{
			case 1:
			//axially loaded/truss member (end) fillet weld w/ gusset plate
			switch w.Wltyp{
				case 1:
				//=
				//check for equal overlap etc here
				case 2:
				//[
				w.Lweld = 2.0 * w.Lop + w.Dmem
				case 3:
				//[]
				w.Lweld = 2.0 * (w.Lop + w.Dmem)
			}		
			switch{
				case w.Pu == 0.0:
				//check weld
				w.Pdw = w.Lweld * w.Ks * w.Size * w.Fu/w.Ymf/math.Sqrt(3.0)
				fmt.Println("design istrengths pdw is",w.Pdw/1e3, "KN")		
				case w.Pu > 0.0:
				//calculate weld
				
			}
		}
	}
	return
}

//WeldDz designs a welded connection 
func WeldDz(w *Wld)(err error){
	err = w.Init()
	if err != nil{
		return
	}
	switch w.Ctyp{
		case 1:
		//fillet weld
		switch w.Ltyp{
			case 1:
			//axially loaded (truss) weld
			w.Lweld =  w.Pu /(w.Ks * w.Size * w.Fu/w.Ymf/math.Sqrt(3.0))
			w.Ert = 2.0 * w.Size
			switch w.Msname{
				case "l","ln","l2-ss","ln2-ss","l2-os","ln2-os":
				if w.Styp == 0{
					w.Styp = StlBstyps[w.Msname][1]
				}
				s := SecGen(w.Styp, w.Dims[0])
				if s.Prop.Area == 0.0{
					err = fmt.Errorf("error in section generation")
					return
				}
				pu := w.Pu
				var p1, p2, erat float64
				switch w.Msname{
					case "l2-ss","l2-os","ln2-ss","ln2-os":
					pu = w.Pu/2.0
				}
				h := w.Dims[0][1]
				if w.Dmem == 0.0{w.Dmem = h}
				h2 := h - s.Prop.Yc
				h1 := h - h2	
				switch w.Wltyp{
					case 1:
					//=
					p1 = pu * h1/h
					p2 = pu * h2/h
					erat = 4.0 * w.Ert
					case 2:
					//[
					p3 := h * w.Ks * w.Size * w.Fu/w.Ymf/math.Sqrt(3)
					p2 = (pu * h2 - p3 * h/2.0)/h
					p1 = (pu * h1 - p3 * h/2.0)/h
					erat = 2.0 * w.Ert
				}
				w.L1 = p1/(w.Ks * w.Size * w.Fu/w.Ymf/math.Sqrt(3.0))
				w.L2 = p2/(w.Ks * w.Size * w.Fu/w.Ymf/math.Sqrt(3.0))
				w.L1 = math.Ceil(w.L1/5.0)*5.0; w.L2 = math.Ceil(w.L2/5.0)*5.0
				switch w.Wltyp{
					case 1:
					w.Lweld = w.L1 + w.L2 + erat
					case 2:
					w.Lweld = w.L1 + w.L2 + w.Dmem + erat
				}
				fmt.Println("weld l1 and l2",w.L1,w.L2,"size",w.Size, "weld length-",w.Lweld)
				case "rect","flat":
				switch w.Wltyp{
					case 1:
					//=
					w.Lop = w.Lweld/2.0
					case 2:
					//[
					w.Lop = (w.Lweld - w.Dmem)/2.0
					case 3:
					//[]
					w.Lop = w.Lweld/2.0 - w.Dmem
				}
				w.Lop = math.Ceil(w.Lop/5.0)*5.0
				switch w.Wltyp{
					case 1:
					w.Lweld = 2.0 * (w.Lop + w.Ert)
					case 2:
					w.Lweld = 2.0 * (w.Lop + w.Ert) + w.Dmem 
					case 3:
					w.Lweld = 2.0 * (w.Lop + w.Dmem)
				}
				fmt.Println("weld overlap",w.Lop,"size",w.Size, "weld length-",w.Lweld)

			}
		}
	}
	return
}

//WeldDzFax designs an axially loaded fillet weld
func WeldDzFax(w *Wld) (err error){
	return
}

