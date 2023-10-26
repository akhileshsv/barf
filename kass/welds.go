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
	Size   []float64 //fillet size of weld
	Wc     [][]float64 //weld coordinates start x, y, end x, y
	Frc    [][]float64 //applied force - fx, fy, mz
	Fc     [][]float64 // coordinates of point of force (x, y)
	Xn     []float64
	Yn     []float64
	Res    []float64
	Strs   []float64
	Report string
	Title  string
	Id     int
}

//WeldSs performs analysis of a weld group given a set applied forces and weld geometry
//returns stresses at each start, stop point of each weld
func WeldSs(w *Wld) (err error){
	var a, b, c, d, e, h, v, hy, vx, sm, t, sfx, sfy, smt, cmx, cmy float64
	lenz := []float64{}
	w.Xn = make([]float64, len(w.Size) * 2)
	w.Yn = make([]float64, len(w.Size) * 2)
	w.Res = make([]float64, len(w.Size) * 2)
	w.Strs = make([]float64, len(w.Size) * 2)
	//kompute stress resultants
	for i, wc := range w.Wc{
		lenz = append(lenz, math.Sqrt(math.Pow(wc[0] - wc[2],2) + math.Pow(wc[1]-wc[3],2)))
		a += w.Size[i] * lenz[i]
		b += w.Size[i] * lenz[i] * (wc[0] + wc[2])/2.0
		c += w.Size[i] * lenz[i] * (wc[1] + wc[3])/2.0
		t = math.Pow(wc[0],2) + wc[0] * wc[2] + math.Pow(wc[2],2)
		d += w.Size[i] * lenz[i] * t/3.0
		t = math.Pow(wc[1],2) + wc[1] * wc[3] + math.Pow(wc[3],2)
		e += w.Size[i] * lenz[i] * t/3.0
	}
	den := math.Pow(b,2) + math.Pow(c, 2) - a * (d + e)
	for i, frc := range w.Frc{
		h += frc[0]
		hy += frc[0] * w.Fc[i][1]
		v += frc[1]
		vx += frc[1] * w.Fc[i][0]
		sm += frc[2]
	}
	theta := (a * hy - a * vx + a * sm + b * v - c * h)/den
	for i, wc := range w.Wc{
		j := 2 * (i + 1) - 1
		w.Xn[j-1] = w.Size[i] * (-h/a + theta * (wc[1] - c/a))
		w.Xn[j] = w.Size[i] * (-h/a + theta * (wc[3] - c/a))
		w.Yn[j-1] = w.Size[i] * (-v/a + theta * (b/a - wc[0]))
		w.Yn[j] = w.Size[i] * (-v/a + theta * (b/a - wc[2]))
		w.Res[j-1] = math.Sqrt(math.Pow(w.Xn[j-1],2)+math.Pow(w.Yn[j-1],2))
		w.Res[j] = math.Sqrt(math.Pow(w.Xn[j],2)+math.Pow(w.Yn[j],2))
		w.Strs[j-1] = w.Res[j-1]/(0.7071 * w.Size[i])
		w.Strs[j] = w.Res[j]/(0.7071 * w.Size[i])
	}

	//equilibrium check
	for i, wc := range w.Wc{
		j := 2 * (i + 1) - 1
		rmfx := lenz[i] * (w.Xn[j-1] + w.Xn[j])/2.0
		rmfy := lenz[i] * (w.Yn[j-1] + w.Yn[j])/2.0
		sfx += rmfx
		sfy += rmfy
		if wc[1] != wc[3]{
			t1 := (w.Xn[j-1] - w.Xn[j]) * lenz[i]/3.0
			t2 := math.Pow(wc[1],2) + wc[1] * wc[3] + math.Pow(wc[3],2)
			t3 := (w.Xn[j-1] * wc[3] - w.Xn[j] * wc[1]) * lenz[i]
			t4 := (wc[1] + wc[3])/2.0
			t5 := wc[1] - wc[3]
			cmx = (t1 * t2 - t3 * t4)/t5
		} else {
			cmx = lenz[i] * wc[1] * (w.Xn[j-1] + w.Xn[j])/2.0
		}
		if wc[0] != wc[2]{
			t1 := (w.Yn[j-1] - w.Yn[j]) * lenz[i]/3.0
			t2 := math.Pow(wc[0],2) + wc[0] * wc[2] + math.Pow(wc[2],2)
			t3 := (w.Yn[j-1] * wc[2] - w.Yn[j] * wc[0]) * lenz[i]
			t4 := (wc[0] + wc[2])/2.0
			t5 := wc[0] - wc[2]
			cmy = (t1 * t2 - t3 * t4)/t5
		} else {
			cmy = lenz[i] * wc[0] * (w.Yn[j-1] + w.Yn[j])/2.0
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
	for i, wc := range w.Wc{
		row := fmt.Sprintf("%v, %.f, %.f, %.f, %.f, %.f, %.2f",i+1, w.Size[i],wc[0], wc[1], wc[2], wc[3], lenz[i])
		t0.Append(strings.Split(row,","))
	}
	t0.Render()
	t0 = tablewriter.NewWriter(rstr)
	t0.SetHeader([]string{"force","x-coord","y-coord","x-force","y-force","moment"})
	t0.SetCaption(true, "applied forces")
	for i, frc := range w.Frc{
		row := fmt.Sprintf("%v, %.f, %.f, %f, %f, %f",i+1,w.Fc[i][0], w.Fc[i][1], frc[0], frc[1], frc[2])
		t0.Append(strings.Split(row,","))
	}
	t0.Render()
	rstr.WriteString(ColorCyan)
	t1 := tablewriter.NewWriter(rstr)
	t1.SetHeader([]string{"weld","x coord","y coord","force/length","stress"})
	t1.SetCaption(true, "weld group analysis results")
	t0.SetCaption(true, "weld group geometry")
	
	for i, wc := range w.Wc{
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
	w.Report = fmt.Sprintf("%s",rstr)
	if !xchk || !ychk || !mchk{
		return errors.New("equilibrium check failed")
	}
	return nil
}
