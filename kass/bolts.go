package barf

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"errors"
	"strings"
	"github.com/olekukonko/tablewriter"
	flay "barf/flay"
)

var BlkBltDias = []float64{12,14,16,18,20,22,24,27,30}
var Bctyps = map[int]string{
	1:"lap joint",
	2:"single cover butt joint",
	3:"double cover butt joint",
}

//table 5.2, duggal (pg. 164)
var Bgrds = map[float64][]float64{
	3.6 : []float64{180,330,25},
	4.6 : []float64{240,400,22},
	4.8 : []float64{320,420,14},
	5.6 : []float64{300,500,20},
	5.8 : []float64{400,520,10},
	6.8 : []float64{480,600,8},
	8.8 : []float64{640,800,12},
	9.8 : []float64{720,900,10},
	10.9 : []float64{940,1040,9},
	12.9 : []float64{1100,1220,8},
}

// var Bbtyps = map[int]string{
// 	1:"chain bolts",
// 	2:"staggered bolts",
// 	3:"diamond bolts",
	
// }

//Blt is/be/'tis a bolt group
//as seen in harrison section 4.1
//all 2d coords/forces
type Blt struct{
	Nc     []float64 //node coords
	Id     int //node id
	B0     []float64 //bolt origin
	Cs     [][]float64 //final coords (might be useful later)
	Bc     [][]float64 //bolt coords (x1,y1), (x2,y2)
	Bxs    []float64 
	Bys    []float64
	Ps     []float64  //pitches/gauges
	Gs     []float64 
	Bvec   [][]int
	Dia    float64  //set single dia for bolt group
	Hdia   float64
	D0     float64 //hole diameter
	Dias   []float64 //diameter (dia1,dia2)
	Typ    []float64 //single (1.0) or double (2.0) shear
	Frc    [][]float64 //(2d) force vectors of forces acting on bolt group (force x1, force y1, moment z1) 
	Fc     [][]float64 //coords of forces (x1, y1),(x2,y2)
	Xn     []float64
	Yn     []float64
	Res    []float64
	Strs   []float64
	Area   []float64 //slice of bolt areas
	Asb    float64 //nominal shank area of bolt
	Anb    float64 //net shear area at threads
	Bta    float64  //tensile stress area of bolt
	Ni     int //no. of bolt rows (lines along force - along y)
	Nj     int //no. of bolt columns (bolts/line - along x)
	Nb     int //no. of bolts
	Nb2    int //no. of bolts in supporting mem
	Stag   bool //is staggered
	Ctyp   int     //connection type - 1 - lap, 2 - single cover butt joint, 3 - double cover
	Ltyp   int     //load type - 1 - concentric (T/C), 2 - 
	Bltyp  int     //bolt layout type - 1 - chain, 2 - staggered, 3 - diamond
	Btyp   int     //bolt type 1 - ordinary black bolts, 2 - hsfg bolts
	Mtyp   int     //base material type - (wood-wood(33), steel-steel(22),wood-steel plate(32), etc) 
	Ljoint float64 //length of joint
	Bjoint float64
	Lgrip  float64 //grip length of bolt
	Pitch  float64 //pitch of bolts
	Grade  float64
	Fy     float64 //yield stress of bolt in N/mm2
	Fub    float64 //ult. tensile stress of bolt
	Fup    float64 //ult. tensile stress of plate
	Fyplt  float64 //yield stress of plate in N/mm2
	Gauge  float64
 	Endd   float64
	Edged  float64
	Stagd  float64
	Nsar   float64     //net sectional area
	Nn     float64     //no. of shear planes with threads 
	Ns     float64     //no. of shear planes without threads 
	Blj    float64
	Blg    float64
	Bpk    float64
	Bfx    float64
	Bfy    float64
	Bmt    float64
	Pu     float64
	Px     float64 
	Py     float64
	Pxy    float64
	Pyx    float64
	Pm     float64
	Psi    float64
	Kb     float64
	T1     float64 //thickness of mem1
	T2     float64 //thickness of mem2
	T3     float64 //thickness of plate/cleat
	Vdsb   float64 //design shear strength of single bolt
	Vdpb   float64 //design strength of single bolt in bearing
	Vnsf   float64 //design shear strength of single (HSFG) bolt
	Ymb    float64 //partial safety factor of material for bolt
	Uf     float64 //table 20, is 800 (coeff of friction/slip factor)
	Nef    float64 //no. of frictional resist. interfaces
	Ymf    float64 //partial safety factor of material for hsfg bolt
	Kh     float64 //hole factor for hsfg bolt
	F0     float64 //proof load for hsfg bolt
	Psec   SectIn //only rect plate for now
	Pt     float64 //plate thickness
	Psis   []float64
	T      float64 //design plate thickness
	Tmem   float64 //thickness of smallest mem
	Bmem   float64 //width of 
	Dmem   float64 //use for end plates/cleat depths, min. 0.6 Dmem
	Tms    []float64 //connected mem. thicknesses
	Dims   []float64 //connection dimensions
	Mdims  [][]float64 //list of member dims
	Mtyps  []int //member types (1,2 - beam, col)
	Mdxs   []string //member indices
	S1, S2 int    //styp 1, styp 2
	Cloc   string //"web","flange"
	Vdu    float64 //max design shear
	Mdu    float64 //max design moment
	Vb     float64 //design strength of single bolt
	Vb2    float64 //design strength of bolt for mem T2
	Folder string
	Report string
	Term   string
	Title  string
	Name   string
	Type   string //simple, framed, etx
	Weld   Wld 
	UniPig bool //uniform pitch and gauge
	Slip   bool //if slip is permitted for hsfg bolts
	Web    bool
	Print  bool
	Xchk   bool
	Ychk   bool
	Mchk   bool
	Verbose bool
	Chsfg  bool //column hsfg bolt for fep
}

//BoltSs does bolt group analysis - see harrision sec 4.1
//requires bolt coords, bolt dia, shear type, forces and force coords as input
//computes stresses in each bolt
//"the backbone" of bolt group calcs
func BoltSs(b *Blt) (err error){
	var atot, ax, ax2, ay, ay2, px, py, pxy, pyx, pm float64
	b.Px = 0.0; b.Py = 0.0; b.Pxy = 0.0; b.Pyx = 0.0; b.Pm = 0.0
	b.Area = make([]float64, len(b.Typ))
	b.Xn = make([]float64, len(b.Typ))
	b.Yn = make([]float64, len(b.Typ))
	b.Res = make([]float64, len(b.Typ))
	b.Strs = make([]float64, len(b.Typ))
	for i := range b.Bc{
		b.Area = append(b.Area, b.Typ[i]*math.Pi * math.Pow(b.Dias[i],2)/4.0)
		atot = atot + b.Area[i]
		ax += b.Area[i] * b.Bc[i][0]
		ax2 += math.Pow(b.Bc[i][0],2) * b.Area[i]
		ay += b.Area[i] * b.Bc[i][1]
		ay2 += math.Pow(b.Bc[i][1],2) * b.Area[i]
	}
	den := (math.Pow(ax, 2) + math.Pow(ay, 2))/atot - ax2 - ay2
	if den == 0.0{
		err = fmt.Errorf("zero division error->%f",den)
		return
	}
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
	for i, ar := range b.Area{
		b.Xn[i] = ar * (defx + theta * b.Bc[i][1])
		b.Yn[i] = ar * (defy - theta * b.Bc[i][0])
		b.Res[i] = math.Sqrt(math.Pow(b.Xn[i],2)+math.Pow(b.Yn[i],2))
		b.Strs[i] = b.Res[i]/b.Area[i]
	}
	//equilibrium check
	for i, xn := range b.Xn{
		yn := b.Yn[i]
		b.Bfx += xn
		b.Bfy += yn
		b.Bmt += xn * b.Bc[i][1] - yn * b.Bc[i][0]
	}
	b.Xchk = math.Abs(b.Bfx + b.Px) < 1e-3
	b.Ychk = math.Abs(b.Bfy + b.Py) < 1e-3
	b.Mchk = math.Abs(b.Bmt + b.Pm + b.Pxy - b.Pyx) < 1e-3
	if !b.Xchk || !b.Ychk || !b.Mchk{
		return errors.New("equilibrium check failed")
	}
	if b.Print{
		b.TableAz()
	}
	return nil
}

//NetSecArea calculates the net section area of a bolt group
func (b *Blt) NetSecArea()(err error){
	var an float64
	switch b.Bltyp{
		case 1:
		//chain bolting
		n := float64(b.Ni)
		an = (b.Bmem - n * b.D0)*b.Tmem
		fmt.Println("b, n, d0, t -",b.Bmem, n, b.D0, b.Tmem)
		fmt.Println("net sec area->",an)
		case 2:
		//staggered/zig-zag bolting
		case 3:
		case 4:
	}
	return
}

//AzChk prints an ascii report of analysis results
func (b *Blt) TableAz(){
	
	//build analysis report
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
		row := fmt.Sprintf("%v, %.2f, %.2f, %.f, %.f, %.2f",i+1, bc[0], bc[1], b.Dias[i], b.Typ[i], b.Area[i])
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
		row := fmt.Sprintf("%v, %.f, %f, %f, %f, %f, %.f",i+1, b.Dias[i],b.Xn[i], b.Yn[i], b.Res[i], b.Strs[i], b.Typ[i])
		t1.Append(strings.Split(row,","))
	}
	t1.Render()
	rstr.WriteString(ColorPurple)
	t1 = tablewriter.NewWriter(rstr)
	t1.SetHeader([]string{"action","sum-bolt","sum-applied","ok?"})
	t1.SetCaption(true, "equilibrium check")
	row := fmt.Sprintf("%s, %f, %f, %t","x-force",b.Bfx, b.Px, b.Xchk)
	t1.Append(strings.Split(row, ","))
	row = fmt.Sprintf("%s, %f, %f, %t","y-force",b.Bfy, b.Py, b.Ychk)
	t1.Append(strings.Split(row, ","))
	row = fmt.Sprintf("%s, %f, %f, %t","moments-origin",b.Bmt, b.Pm + b.Pxy - b.Pyx, b.Mchk)
	t1.Append(strings.Split(row, ","))
	t1.Render()
	rstr.WriteString(ColorReset)
	b.Report = rstr.String()
}

//TableDz prints an ascii report of design results
//does it. where is it. 

//Init initializes default values
func (b *Blt) Init()(err error){
	if b.Btyp == 0{
		b.Btyp = 1
	}
	if b.Grade == 0.0{
		b.Grade = 4.6
	}
	if _, ok := Bgrds[b.Grade]; !ok{
		err = fmt.Errorf("invalid bolt grade - %.1f",b.Grade)
		return
	}
	b.Fy = Bgrds[b.Grade][0]
	b.Fub = Bgrds[b.Grade][1]
	switch{
		case b.Name == "dac":
		b.T = b.T1
		b.Tmem = b.T1
		if b.T3 == 0.0{
			b.T3 = 10.0
		}
		case b.Name == "fep":
		case b.Name == "fp":
	}
	//find smallest connected mem. thickness
	if b.Tmem == 0.0{
		for _, t := range b.Tms{
			if b.Tmem == 0.0{
				b.Tmem = t
			} else if b.Tmem > t{
				b.Tmem = t
			}
		}
	}
	//set design (min.) plate thickness
	if b.T == 0.0{
		b.T = b.Tmem
	}
	if b.Fup == 0.0{b.Fup = 410.0}
	switch b.Btyp{
		case 2:
		if b.Uf == 0.0{
			b.Uf = 0.5
		}
		if b.Kh == 0.0{
			b.Kh = 1.0
		}
		if b.Ymf == 0.0{
			b.Ymf = 1.25
		}
	}
	if b.Ni > 0 && b.Nj > 0{
		if b.Bltyp == 0{
			b.Bltyp = 1
		}
		err = b.BltVec()
	}
	return
}

//BltVals updates hole dia, min. edge dist, min. end dist, min. pitch for a bolt dia
func (b *Blt) BltVals() (err error){
	//dias {12,14,16,18,20,22,24,27,30}
	//hole diameters

	// hdias := []float64{13, 15, 18, 20, 22, 24, 27, 30, 33}
	// //sheared or hand flame cut edge
	// edgd1s := []float64{20, 26, 30, 34, 37, 40, 44, 51, 56}
	// //rolled or machine cut edge distance - add field for this
	// //{12,14,16,18,20,22,24,27,30}
	// if b.Dia > 30.0{
	// 	b.D0 = b.Dia + 3.0
	// 	b.Edged = 1.5 * b.D0 //use with sdx(standard secs)
	// 	//b.Edged = 1.7 * b.Hdia
	// } else {
	// 	for i, dia := range BlkBltDias{
	// 		if dia == b.Dia{
	// 			b.D0 = hdias[i]
	// 			b.Edged = edgd1s[i]
	// 			break
	// 		}
	// 	}
	// }
	//eh use this. maity does
	
	b.D0 = b.Dia + 2.0
	b.Edged = 1.5 * b.D0

	if b.D0 == 0.0{
		err = fmt.Errorf("invalid diameter -> %.2f for bolt",b.Dia)
		return
	}
	b.Asb = math.Pi * math.Pow(b.Dia, 2.0)/4.0
	b.Anb = 0.78 * b.Asb
	if b.Pitch == 0.0{b.Pitch = 2.5 * b.Dia}
	b.Kb = b.Edged/(3.0*b.D0)
	if b.Kb > b.Pitch/(3.0*b.D0) - 0.25{
		b.Kb = b.Pitch/(3.0*b.D0) - 0.25
	}
	
	if b.Kb > b.Fub/b.Fup{b.Kb = b.Fub/b.Fup}
	if b.Kb > 1.0{b.Kb = 1.0}
	//log.Println("b.Kb->",b.Kb)
	//set reduction factors
	b.Blj = 1.0; b.Blg = 1.0; b.Bpk = 1.0
	//for long joints
	if b.Ljoint > 15.0 * b.Dia{
		b.Blj = 1.075 - b.Ljoint/(200.0 * b.Dia)
		if b.Blj < 0.75{b.Blj = 0.75}
		if b.Blj > 1.0{b.Blj = 1.0}
	}
	//for large grip lengths
	if b.Lgrip == 0.0{
		switch b.Ctyp{
			case 1:
			if len(b.Tms) >= 2{b.Lgrip = b.Tms[0] + b.Tms[1]}
			case 2, 3:
			if len(b.Tms)>=2 {b.Lgrip =  b.Tms[0] + b.Tms[1] + b.Pt}	
			if b.Ctyp == 3{
				b.Lgrip += b.Pt
			}
			default:
		}
	}
	if b.Lgrip > 5.0 * b.Dia{
		b.Blg = 8.0 * b.Dia/(3.0 * b.Dia + b.Lgrip)
	}
	if b.Ctyp == 2{
		if b.T > b.Pt && b.Pt > 0.0{
			b.T = b.Pt
		}	
	}
	if b.Ctyp == 3{
		//packing plates ahoy		
		tpkg := math.Abs(b.Tms[0]-b.Tms[1])
		b.Bpk = 1.0 - 0.0125*tpkg
		if b.T > 2.0 * b.Pt{
			b.T = 2.0 * b.Pt
		}
	}
	//set max gauge
	if b.Gauge <= 0.0{
		b.Gauge = 200.0
	}
	if b.Gauge > 100.0 + 4.0 * b.T{
		b.Gauge = 100.0 + 4.0 * b.T
	}
	if b.Endd == 0.0{
		b.Endd = b.Edged
	}
	return
}

//BltVec generates the bolt grid vector/array/bvec
func (b *Blt) BltVec()(err error){
	b.Bvec = bltvec(b.Ni, b.Nj, b.Bltyp)
	if len(b.Bvec) == 0 || len(b.Bvec[0]) == 0{
		err = fmt.Errorf("invalid params for bolt grid array- %v %v %v",b.Ni, b.Nj, b.Bltyp)
	}
	return
}

//BltNsa returns the net sectional area of the bolt group
func (b *Blt) BltNsa()(err error){
	err = b.Init()
	
	if err != nil{
		return
	}
	err = b.BltVals()
	if err != nil{
		return
	}
	err = b.BltVec()
	if err != nil{
		return
	}
	_, _, nsamin, _ := flay.BltSecArea(b.Bvec, b.Bmem, b.Tmem,b.D0, b.Pitch, b.Gauge, b.Ps, b.Gs)
	b.Nsar = nsamin
	//tbh pitch doesn't vary for an equivalent plate so
	b.Ljoint = b.Pitch * float64(b.Nj - 1)
	if b.Ljoint == 0.0{b.Ljoint = b.Pitch}
	return
}

//BltDiaCalc calculates single bolt values ctyp, force group and dia
func BltDiaCalc(b *Blt) (err error){
	b.Init()
	err = b.BltVals()
	if err != nil{
		return
	}
	if b.Verbose{
		log.Println("type of bolt->",b.Btyp,"1-ord., 2-hsfg")
		log.Println("checking dia->",b.Dia,"grade->",b.Grade,"for conn. typ->",b.Ctyp, b.Name)
		log.Println("pitch->",b.Pitch,"edge distance->",b.Edged, "t->",b.T)
		log.Println("reduction factors->",b.Blj, b.Blg, b.Bpk)		
		log.Println("joint name->",b.Name)
	}
	switch{
		case b.Name == "dac":
		//double angle cleat connection
		b.Nn = 2.0; b.Ns = 0.0; b.Ymb = 1.25; b.Nef = 1.0
		b.T = b.T1
		if b.T > b.T3{
			b.T = b.T3
		}
		case b.Name == "fep":
		//flexible end plate connection
		b.Nn = 1.0; b.Ns = 0.0; b.Ymb = 1.25; b.Nef = 1.0
		b.T = b.T1
		if b.T > b.T3{
			b.T = b.T3
		}
		//add b.T1 vs b.T3; b.T2 vs b.T3
		case b.Name == "fp":
		//fin plate connection
		b.Nn = 1.0; b.Ns = 0.0; b.Ymb = 1.25; b.Nef = 1.0
		b.T = b.T1
		if b.T > b.T3{
			b.T = b.T3
		}
		
		case b.Ctyp == 1:
		b.Nn = 1.0; b.Ns = 0.0; b.Ymb = 1.25; b.Nef = 1.0
		case b.Ctyp == 2:
		b.Nn = 1.0; b.Ns = 0.0; b.Ymb = 1.25; b.Nef = 1.0 	
		case b.Ctyp == 3:
		b.Nn = 2.0; b.Ns = 0.0; b.Ymb = 1.25; b.Nef = 1.0
	}
	b.Vdsb = b.Fub * (b.Nn * b.Anb + b.Ns * b.Asb)/math.Sqrt(3.0)/b.Ymb
	b.Vdsb = b.Vdsb * b.Blj * b.Blg * b.Bpk
	b.Vb = b.Vdsb
	b.Vdpb = 2.5 * b.Kb * b.Dia * b.T * b.Fup/b.Ymb
	if b.Vdpb < b.Vb{b.Vb = b.Vdpb}
	if b.Btyp == 2{
		b.F0 = 0.7 * b.Anb * b.Fub 
		b.Vnsf = b.F0 * b.Uf * b.Nef * b.Kh/b.Ymf
	}	
	if b.Verbose{
		log.Printf("design strength in bearing - dia %.3f ymb %.3f fup %.3f kb %.3f T %.3f\n",b.Dia, b.Ymb, b.Fup, b.Kb, b.T)
		log.Println("design strength in bearing of single bolt->",b.Vdpb/1e3,"kn")
		log.Println("design strength in shear of single bolt->",b.Vdsb/1e3,"kn")
		log.Printf("net area asb - %.2f x nn %.f anb - %.2f x ns %.f\n",b.Anb, b.Nn, b.Ns, b. Asb)
		log.Printf("proof load of bolt - %.3f kn\n",b.F0/1e3)
		log.Println("design strength of bolt->",b.Vb/1e3,"kn")
	}
	switch b.Btyp{
		case 2:
		if b.Verbose{
			log.Printf("design strength in shear (slip) - %.3f kn\n",b.Vnsf/1e3)
		}
	}
	switch b.Name{
		case "dac":
		//calc bearing on T2 (column), T1 is beam
		
		vdsb := b.Fub * (b.Nn/2.0 * b.Anb + b.Ns * b.Asb)/math.Sqrt(3.0)/b.Ymb
		vdsb = vdsb * b.Blj * b.Blg * b.Bpk
		b.Vb2 = vdsb
		vdpb := 2.5 * b.Kb * b.Dia * b.T2 * b.Fup/b.Ymb
		if b.T3 < b.T2{
			vdpb = 2.5 * b.Kb * b.Dia * b.T3 * b.Fup/b.Ymb
		
		}
		if vdpb < b.Vb2{b.Vb2 = vdpb}
		case "fep":
		
	}
	return
}


//BltLay lays out a bolt group
func BltLay(b *Blt)(err error){
	switch b.Ltyp{
		case 1:
		fmt.Println("pitch, gauge ->", b.Pitch, b.Gauge)
		fmt.Println("edge distance ->", b.Edged)
		fmt.Println("end distance ->",b.Endd)
		fmt.Println("b.Nj, b.Ni, width available -> ",b.Nj, b.Ni, b.Bmem)
		//set actual gauge
		var g float64
		if b.Ni == 1{
			g = b.Bmem - 2.0 * b.Edged
		} else {
			g = (b.Bmem - 2.0 * b.Edged)/float64(b.Ni - 1)
		}
		fmt.Println("actual gauge - ",g ,"mm")
		
	}
	return
}

//BltChk checks a bolt group for applied force
func BltChk(b *Blt) (err error){
	return BltLay(b)
}

//bltvec returns the bolt matrix given an ni (nrows), nj (ncols) and bltyp (bolt layout type) 
func bltvec(ni, nj, bltyp int)(bvec [][]int){
	switch bltyp{
		case 1:
		//chain bolting
		for i := 0; i < ni; i++{
			var vec []int
			for j := 0; j < nj; j++{
				vec = append(vec, 1)
			}
			bvec = append(bvec, vec)
		}
		case 2:
		//staggered bolting - even/odd
		for i := 0; i < ni; i++{
			var vec []int
			for j := 0; j < nj; j++{
				switch{
					case j % 2 == 0:
					switch{
						case i % 2 == 0:
						vec = append(vec, 1)
						case i % 2 == 1:
						vec = append(vec, 0)
					}
					case j % 2 == 1:
					switch{
						case i % 2 == 0:
						vec = append(vec, 0)
						case i % 2 == 1:
						vec = append(vec, 1)
					}
				}
			}
			bvec = append(bvec, vec)
		}
		case 3:
		//diamond
		if ni % 2 == 0{
			ni += 1
		}
		mi := ni/2
		bvec = make([][]int, ni)
		for i := range bvec{
			bvec[i] = make([]int, nj)
		}
		for j := 0; j < nj; j++{
			tp := mi - j
			bm := mi + j
			for i := 0; i < ni; i++{
				if i >= tp && i <= bm{
					bvec[i][j] = 1
				}
			}
		}
		case 4:
		//staggered diamond
		if ni % 2 == 0{
			ni += 1
		}
		mi := ni/2
		bvec = make([][]int, ni)
		for i := range bvec{
			bvec[i] = make([]int, nj)
		}
		for j := 0; j < nj; j++{
			tp := mi - j
			bm := mi + j
			for i := 0; i < ni; i++{
				if i >= tp && i <= bm{
					idx := bm - i
					if idx % 2 == 0{
						bvec[i][j] = 1
					} 
				}
			}
		}
	}
	return
}

//AbltLen returns the length of an anchor bolt for a given uplift force
func AbltLen(pulf, fck, dia float64)(lblt float64, err error){
	tbdz := map[float64]float64{20.0:1.2, 25.0:1.4, 30.0:1.5, 35.0:1.7, 40.0:1.9}
	anb := 0.78 * math.Pi * math.Pow(dia, 2.0)/4.0
	tdb := 0.9 * 330.0 * anb/1.25
	nblt := pulf/tdb
	//fmt.Println("nbolts-",nblt,"using 4")
	if nblt > 4{
		err = fmt.Errorf("revise dia, > 4 bolts required %f",nblt)
		return
	}
	//nblt = 4
	ublt := pulf/4
	lbltp := math.Pow(ublt/(15.5 * math.Sqrt(fck)),1.0/1.5)
	//fmt.Println("bolt len from pull out criteria-",lbltp,"mm")
	if _, ok := tbdz[fck]; !ok{
		err = fmt.Errorf("invalid fck %f",fck)
		return
	}
	tbd := tbdz[fck]
	lbltb := dia * 250.0/(4.0 * tbd * 1.1)
	//fmt.Println("bolt len from bond criteria-",lbltb,"mm")
	lblt = lbltp; if lbltb > lblt{
		lblt = lbltb
	}
	lblt = math.Ceil(lblt/25.0)*25.0
	return
}


//BltDz designs a bolt group given a ctyp and force group
//three ltyps as in bhavikatti, chap 3
func BltDz(b *Blt) (err error){
	switch b.Name{
		case "dac":
		//double angle cleats
		err = BltDac(b)
		return
		case "fep":
		//flexible end plate
		err = BltFep(b)
		return
		case "fp":
		//fin plate
		err = BltFp(b)
		return
	}
	err = BltDiaCalc(b)
	if err != nil{
		return
	}
	switch b.Ltyp{
		case 0:
		case 1:
		//axial/concentric load
		nb := math.Ceil(b.Pu/b.Vb)
		b.Nb = int(nb)
		if b.Nb <= 0{
			err = fmt.Errorf("invalid number of bolts - %v",b.Nb)
			return
		}
		err = BltLay(b)
		if err != nil{
			return
		}
		case 2:
		//eccentric bracket connection; in-plane moment
		
		case 3:
		//eccentric bracket connection; out-of plane moment
	
	}
	return
}



//BltFep designs a flexible end plate connection
//welded to the beam (main mem), bolted to the support
func BltFep(b *Blt)(err error){	
	if b.Ctyp == 0{b.Ctyp = 1}
	err = BltDiaCalc(b)
	if err != nil{
		return
	}
	fmt.Println("bolt vb1, vb2",b.Vb/1e3,b.Vb2/1e3)
	//calc nbolts at supporting member
	b.Nb = int(math.Ceil(b.Vdu/b.Vb/2.0)*2.0)
	b.Nj = 2
	b.Ni = b.Nb/b.Nj
	fmt.Println("nbolts",b.Nb,"ni,nj",b.Ni,b.Nj)
	
	//get length of end plate
	b.Edged = math.Ceil(b.Edged/10.0)*10.0
	gauge := 90.0
	if b.Gauge < gauge{
		gauge = math.Ceil(b.Gauge/10.0)*10.0
	}
	b.Gauge = gauge
	lplt := float64(b.Ni-1)*b.Pitch + 2.0 * b.Edged + 20.0
	wplt := float64(b.Nj-1)*gauge + 2.0 * b.Edged 
	//lplt := dplt
	//if wplt > lplt{lplt = wplt}
	a := 30.0 * b.T3
	if a > lplt{a = lplt}
	fmt.Println("lplt",lplt, "wplt",wplt,"pitch",b.Pitch,"edged",b.Edged,"a",a)
	fdw := 410.0/1.5/math.Sqrt(3) 
	//assume weld size of 6mm for lwld calcs
	lwld := 2.0 * (lplt - 2.0 * 6.0)
	wsize := b.Vdu/lwld/fdw
	fmt.Println("lwld, fdw, weld size",lwld/2.0, fdw, wsize)
	switch{
		case wsize < 6.0:
		wsize = 6.0
		case wsize < 8.0:
		wsize = 8.0
		
	}
	return
}

//BltFp designs a fin plate connection
//welded to the support, bolted to the main mem
func BltFp(b *Blt)(err error){	
	if b.Ctyp == 0{b.Ctyp = 1}
	err = BltDiaCalc(b)
	if err != nil{
		return
	}
	fmt.Println("bolt vb1, vb2",b.Vb/1e3,b.Vb2/1e3)
	//calc nbolts at supporting member
	b.Nb = int(math.Ceil(b.Vdu/b.Vb))
	//if double row?
	b.Nj = 1
	b.Ni = b.Nb/b.Nj
	fmt.Println("nbolts",b.Nb,"ni,nj",b.Ni,b.Nj)
	
	//get length of end plate
	b.Edged = math.Ceil(b.Edged/10.0)*10.0
	gauge := 90.0
	if b.Gauge < gauge{
		gauge = math.Ceil(b.Gauge/10.0)*10.0
	}
	b.Gauge = gauge

	lmin := 0.75 * b.Dmem
	lplt := float64(b.Ni-1)*b.Pitch + 2.0 * b.Edged + 20.0
	if lplt < lmin && lmin != 0.0{
		lplt = lmin
	}
	fmt.Println("lmin, lplt",lmin, lplt)
	b.Pitch = (lplt - 2.0 * b.Edged - 20.0)/float64(b.Ni-1)
	cof := b.Edged + 20.0
	fmt.Println("edged, cof",b.Edged,cof, "rev pitch",b.Pitch)
	nb := float64(b.Nb) 
	bmb := b.Vdsb * (nb - 1.0) * b.Pitch * b.Pitch/((nb-1.0)*b.Pitch/2.0) 
	fmt.Println("bmb",bmb,"bmb",bmb/1e3)
	zplt := b.T3 * lplt * lplt/6.0
	bmplt := 1.2 * (250.0/1.1) * zplt
	fmt.Println("bmplt",bmplt,"zplt",zplt,"ok?",bmplt>zplt)
	wsize := math.Floor(0.6 * b.T3/0.7)
	fmt.Println("wsize",wsize)
	switch{
		case wsize <= 6.0:
		wsize = 6.0
		case wsize <= 8.0:
		wsize = 8.0
	}
	lwld := lplt - 2.0 * wsize
	zwld := lwld * lwld/3.0
	vh := bmb/zwld
	vv := b.Vdu/(2.0 * lwld)
	vrez := math.Sqrt(vh*vh + vv*vv)
	fmt.Println("vrez",vrez)
	fdw := 410.0/math.Sqrt(3)/1.25
	pdw := fdw * 0.7 * wsize
	fmt.Println("pdw",pdw,"pdw>vrez?",pdw>vrez)
	//w.Size * w.Fu/w.Ymf/math.Sqrt(3.0)
	return
}

//BltDac designs a double angle cleat connection 
func BltDac(b *Blt) (err error){
	//beam web to column flange/web joint
	//default - column beam dac
	//default s1,s2 - 12, 12 (i, i)
	//angle sizes - 100x100x10, 150x150x10, 150x150x12, 150x150x16
	//angz := [][]float64{{100,10},{150,10},{150,12},{150,16}}
	if b.Ctyp == 0{b.Ctyp = 1}
	err = BltDiaCalc(b)
	if err != nil{
		return
	}
	b.Nb = int(math.Ceil(b.Vdu/b.Vb/2.0)*2.0)
	b.Nb2 = int(math.Ceil(b.Vdu/b.Vb2/2.0)*2.0)
	
	b.Edged = math.Ceil(b.Edged/10.0)*10.0
	//b.Edged += 10.0
	if b.Nb > 4{
		b.Nj = 2
	} else {
		b.Nj = 1
	}
	if b.Nb2 < 4{b.Nb2 = 4}
	b.Ni = b.Nb/b.Nj
	//make gauge the pitch, why not
	depth := float64(b.Ni-1)*b.Pitch+2.0*b.Edged
	width := float64(b.Nj-1)*b.Pitch+2.0*b.Edged
	thick := b.Vdu * math.Sqrt(3) * 1.1/(250.0 * depth * 2.0)
	switch{
		// case thick < 6.0:
		// thick = 6.0
		case thick < 8.0:
		thick = 8.0
		case thick < 10.0:
		thick = 10.0
		case thick < 12.0:
		thick = 12.0
		case thick < 16.0:
		thick = 16.0
	}
	switch{
		case width < 75.0:
		width = 75.0
		case width < 100.0:
		width = 100.0
		case width < 150.0:
		width = 150.0
		case width < 200.0:
		width = 150.0
		
	}

	if b.Verbose{
		fmt.Println("bolt vb1, vb2",b.Vb/1e3,b.Vb2/1e3)
		fmt.Println("nbolts main mem-",b.Nb,"supporting mem",b.Nb2)
		fmt.Println("pitch-",b.Pitch,"edge d",b.Edged,"endd",b.Endd)
		fmt.Println("max gauge-",b.Gauge)
		fmt.Println("depth of angle-",depth)
		fmt.Println("width of angle-",width)
		fmt.Println("thickness of angle-",thick)
	}
	//if this is > T3 (10mm then wot?)
	//b.T3 = thick
	b.Dims = []float64{depth, width, thick} 
	//now consider checks depth = 0.6 * dbeam
	err = b.Draw()

	return
}
