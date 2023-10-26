package barf

import (
	"fmt"
	"math"
	"math/rand"
	"errors"
	"strings"
	"github.com/olekukonko/tablewriter"	
)

//Blt is/be/'tis a bolt group
//as seen in harrison section 4.1
//all 2d coords/forces
type Blt struct{
	Bc     [][]float64 //bolt coords (x1,y1), (x2,y2)
	Dia    []float64 //diameter (dia1,dia2)
	Typ    []float64 //single (1.0) or double (2.0) shear
	Frc    [][]float64 //(2d) force vectors of forces acting on bolt group (force x1, force y1, moment z1) 
	Fc     [][]float64 //coords of forces (x1, y1),(x2,y2)
	Xn     []float64
	Yn     []float64
	Res    []float64
	Strs   []float64
	Report string
	Title  string
	Id     int
}

//BoltSs does bolt group analysis - see harrision sec 4.1
//requires bolt coords, bolt dia, shear type, forces and force coords as input
//computes stresses in each bolt
func BoltSs(b *Blt) (err error){
	var atot, ax, ax2, ay, ay2, px, py, pxy, pyx, pm float64
	area := []float64{}
	b.Xn = make([]float64, len(b.Typ))
	b.Yn = make([]float64, len(b.Typ))
	b.Res = make([]float64, len(b.Typ))
	b.Strs = make([]float64, len(b.Typ))
	for i := range b.Bc{
		area = append(area, b.Typ[i]*math.Pi * math.Pow(b.Dia[i],2)/4.0)
		atot = atot + area[i]
		ax += area[i] * b.Bc[i][0]
		ax2 += math.Pow(b.Bc[i][0],2) * area[i]
		ay += area[i] * b.Bc[i][1]
		ay2 += math.Pow(b.Bc[i][1],2) * area[i]
	}
	den := (math.Pow(ax, 2) + math.Pow(ay, 2))/atot - ax2 - ay2
	for i, frc := range b.Frc{
		px += frc[0]
		py += frc[1]
		pm += frc[2]
		pxy += frc[0] * b.Fc[i][1]
		pyx += frc[1] * b.Fc[i][0]
	}
	theta := (pxy - pyx + py * ax/atot - px * ay/atot + pm)/den
	defx := (-px - theta * ay)/atot
	defy := (-py + theta * ax)/atot
	for i, ar := range area{
		b.Xn[i] = ar * (defx + theta * b.Bc[i][1])
		b.Yn[i] = ar * (defy - theta * b.Bc[i][0])
		b.Res[i] = math.Sqrt(math.Pow(b.Xn[i],2)+math.Pow(b.Yn[i],2))
		b.Strs[i] = b.Res[i]/area[i]
	}
	//equilibrium check
	var bfx, bfy, bmt float64
	for i, xn := range b.Xn{
		yn := b.Yn[i]
		bfx += xn
		bfy += yn
		bmt += xn * b.Bc[i][1] - yn * b.Bc[i][0]
	}

	//build report
	rstr := &strings.Builder{}
	rstr.WriteString(ColorGreen)
	t0 := tablewriter.NewWriter(rstr)
	t0.SetHeader([]string{"bolt","x coord","y coord","dia","type","area"})
	if b.Title == ""{
		if b.Id == 0{
			b.Id = rand.Intn(666)
		}
		b.Title = fmt.Sprintf("no. %v",b.Id)
	}
	cptn := fmt.Sprintf("bolt group geometry-%s",b.Title)
	t0.SetCaption(true, cptn)
	for i, bc := range b.Bc{
		row := fmt.Sprintf("%v, %.2f, %.2f, %.f, %.f, %.2f",i+1, bc[0], bc[1], b.Dia[i], b.Typ[i], area[i])
		t0.Append(strings.Split(row,","))
	}
	t0.Render()
	t0 = tablewriter.NewWriter(rstr)
	t0.SetHeader([]string{"force","x-coord","y-coord","x-force","y-force","moment"})
	t0.SetCaption(true, "applied forces")
	for i, frc := range b.Frc{
		row := fmt.Sprintf("%v, %.f, %.f, %f, %f, %f",i+1,b.Fc[i][0], b.Fc[i][1], frc[0], frc[1], frc[2])
		t0.Append(strings.Split(row,","))
	}
	t0.Render()
	rstr.WriteString(ColorCyan)
	t1 := tablewriter.NewWriter(rstr)
	t1.SetHeader([]string{"bolt","dia","x-force","y-force","res. force","stress","type"})
	t1.SetCaption(true, "bolt group analysis results")
	for i := range b.Strs{
		row := fmt.Sprintf("%v, %.f, %f, %f, %f, %f, %.f",i+1, b.Dia[i],b.Xn[i], b.Yn[i], b.Res[i], b.Strs[i], b.Typ[i])
		t1.Append(strings.Split(row,","))
	}
	t1.Render()
	rstr.WriteString(ColorPurple)
	t1 = tablewriter.NewWriter(rstr)
	t1.SetHeader([]string{"action","sum-bolt","sum-applied","ok?"})
	t1.SetCaption(true, "equilibrium check")
	xchk := math.Abs(bfx + px) < 1e-3
	row := fmt.Sprintf("%s, %f, %f, %t","x-force",bfx, px, xchk)
	t1.Append(strings.Split(row, ","))
	ychk := math.Abs(bfy + py) < 1e-3
	row = fmt.Sprintf("%s, %f, %f, %t","y-force",bfy, py, ychk)
	t1.Append(strings.Split(row, ","))
	mchk := math.Abs(bmt + pm + pxy - pyx) < 1e-3
	row = fmt.Sprintf("%s, %f, %f, %t","moments-origin",bmt, pm + pxy - pyx, mchk)
	t1.Append(strings.Split(row, ","))
	if !xchk || !ychk || !mchk{
		return errors.New("equilibrium check failed")
	}
	t1.Render()
	rstr.WriteString(ColorReset)
	b.Report = fmt.Sprintf("%s",rstr)
	return nil
}
