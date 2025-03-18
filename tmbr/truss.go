package barf

import (
	"log"
	"fmt"
	"math"
	"strings"
	"github.com/olekukonko/tablewriter"
	kass"barf/kass"
)

func Table(t *kass.Trs2d) (err error){
	return
}

func Quant(t *kass.Trs2d) (err error){
	return
}

//ModDz designs a model assuming all members are pin connected columns (for now)
func ModDz(mod *kass.Model){
	colchn := make(chan []interface{}, len(mod.Ms))
	for i := range mod.Ms{
		cp := mod.Mprp[i-1][3]
		dims := mod.Dims[cp-1]
		styp := mod.Sts[cp-1]
		go ColD(mod.Ms[i], styp, mod.Group , dims, mod.Mtprp[0], colchn)
	}
	for range mod.Ms{
		rez := <- colchn
		fmt.Println(rez[0])
	}
}

// //NodeCalc sorts nodal forces and angles subtended at nodes by members
// func NodeCalc(t *kass.Trs2d) (err error){

// }

//NodeCalc designs bolts/nails at beginning and end of each member
func NodeCalc(t *kass.Trs2d) (err error){
	//fcp := t.Dzmem[0].Prp.Fcp
	var tprp kass.Wdprp
	var isapex bool
	tprp.Init(t.Group)
	err = fmt.Errorf("bhak chutiye")
	
	//each joint has to have two chord/rafter members?
	yapx := t.Rise + t.Height
	//var v1, v2, v3, v4 []float64
	for i := range t.Mod.Js{
		var mcs, mws []int
		node := t.Mod.Js[i]
		if t.Spam{fmt.Println("at node -", i, "mems->",node.Mems)}
		nmem := len(node.Mems)
		//classify nodes by number of tc/bc mems
		nbm := 0; ntm := 0
		mtyps := make([]int, nmem)
		frcs := make([]float64, len(node.Mems))
		if node.Coords[1] == yapx && yapx != 0{
			isapex = true
		} else {
			isapex = false
		}
		for j, mid := range node.Mems{
			mem := t.Mod.Ms[mid]
			mtyp := mem.Mprp[3]
			mtyps[j] = mtyp
			frc := mem.Cmax
			if frc < math.Abs(mem.Tmax){
				frc = math.Abs(mem.Tmax)
			}
			frcs[j] = frc
			fmt.Println("at node->",i,"member->",mid,"frc->",frc,"mtyp-",mtyp)
			//get angle at node? WAIT
			switch mtyp{
				case 1:
				//bottom chord
				nbm++
				mcs = append(mcs, mid)
				case 2:
				//top chord
				ntm++
				mcs = append(mcs, mid)
				default:
				mws = append(mws, mid)
			}
		}
		fmt.Println("at node -", i,  "mtyps -",mtyps)
		fmt.Println("nmems -",nmem,"nbot - ",nbm, "ntop - ",ntm)
		fmt.Println("chord members -",mcs)
		switch t.Ctyp{
			
		}
		switch{
			case isapex:
			fmt.Println(ColorYellow,"apex joint",ColorReset)
			//both top chord members are loaded at an angle to the grain
			//ntm == 2
			fmt.Println("number of web members-", len(mws)," said members-",mws)
			case ntm == 1 && nbm == 1:
			fmt.Println(ColorRed,"edge joint",ColorReset)
			//
			case ntm == 2:
			fmt.Println(ColorGreen,"top chord/rafter joint",ColorReset)
			//
			case nbm == 2:
			fmt.Println(ColorBlue,"bottom chord joint",ColorReset)
			default:
			fmt.Println(ColorGreen,"reg. WTF??? joint",ColorReset)
		}
		
	}	
	return
}


//TrussDz designs a truss gen struct assuming all members are pin connected columns
//TODO - design members as beams
func TrussDz(t *kass.Trs2d) (err error){
	err = t.Calc()
	if err != nil{
		log.Println(err)
		return
	}
	fmt.Println("starting node dz ->")
	err = NodeCalc(t)
	if err != nil{
		log.Println(err)
		return
	}
	fmt.Println("sorted nodes")
	
	//err = NodeDz(t)
	log.Println("now dezyneing->")
	colchn := make(chan []interface{}, len(t.Mod.Ms))
	
	for i := 1; i <= len(t.Mod.Ms); i++{
	//for i := range t.Mod.Ms{
		//log.Println("member->",i,"max comp-",t.Mod.Ms[i].Cmax,"max tens-",t.Mod.Ms[i].Tmax)
		cp := t.Mod.Mprp[i-1][3]
		dims := t.Sections[cp-1]
		styp := t.Styps[cp-1]
		go ColD(t.Mod.Ms[i], styp, t.Group , dims, t.Dzval, colchn)
	}
	t.Dzmem = make([]interface{}, len(t.Mod.Ms))
	dimap := make(map[int][]float64)
	for i := 1; i <= t.Ngs; i++{
		dimap[i] = t.Sections[i-1]
	}
	for range t.Mod.Ms{
		rez := <- colchn
		//fmt.Println(rez[0])
		err, _ := rez[0].(error)
		id, _ := rez[1].(int)
		
		if err != nil{
			log.Println(err,id)
			continue
		}
		c, _ := rez[2].(WdCol)
		t.Dzmem[id-1] = c
		cp := t.Mod.Mprp[id-1][3]
		dims := dimap[cp]
		switch c.Styp{
			case 0:
			if dims[0] < c.Rez[0][0]{
				dimap[cp] = c.Dims
			}
			case 1:
			if dims[0] * dims[1] < c.Sec.Prop.Area{
				log.Println("changing-",dims,"-to-",c.Rez[0])
				dimap[cp] = c.Dims
			}
		}
	}
	
	rezstr := new(strings.Builder)
	rezstr.WriteString(ColorCyan+"\n\n")
	
	table := tablewriter.NewWriter(rezstr)
	headers := []string{"member","mprp","length","em","cp","group","styp","dims","max comp","max tens"}
	table.SetHeader(headers)
	table.SetCaption(true,"member design table")
	var colreps string
	lmap := make(map[int]float64)
	amap := make(map[int]float64)
	vmap := make(map[int]float64)
	typmap := map[int]string{1:"chord",2:"rafter",3:"ties",4:"webs",5:"slings",6:"steel"}
	for idx := 1; idx <= len(t.Mod.Ms); idx++{
		//for idx, mem := range t.Mod.Ms{
		mem := t.Mod.Ms[idx]
		cp := mem.Mprp[3]
		em := mem.Mprp[2]
		lspan := mem.Geoms[0]
		mprp := mem.Mprp
		styp := t.Styps[cp-1]
		tsa := 0.0
		vol := 0.0
		dims := dimap[cp]
		switch styp{
			case 0:
			tsa = math.Pi * dims[0] * (lspan + dims[0]/2.0)
			vol = math.Pi * dims[0] * dims[0] * lspan/4.0
			case 1:
			tsa = 2.0 * (dims[0] * dims[1] + dims[0] * lspan + dims[1] * lspan)
			vol = lspan * dims[0] * dims[1]
		}
		lmap[cp] += lspan
		amap[cp] += tsa
		vmap[cp] += vol
		row := fmt.Sprintf("%v,%v,%.2f,%v,%v,%s,%v,%.0f,%.0f,%.0f",idx,mprp,lspan,em,cp,typmap[cp],styp,dims,mem.Cmax,mem.Tmax)
		table.Append(strings.Split(row,","))
		c, _ := t.Dzmem[idx-1].(WdCol)
		colreps += c.Report
	}
	table.Render()
	rezstr.WriteString(ColorGreen)
	table = tablewriter.NewWriter(rezstr)
	headers = []string{"group","dims","net length (m)","net area (m2)","net vol (m3)","nos(10 ft)","cost/fab","cost/coat"}
	table.SetHeader(headers)
	table.SetCaption(true,"member quantities")
	tcost := 0.0
	if len(t.Kostin) == 0{
		t.Kostin = tmbrKost
	}
	for i, grp := range typmap{
		ltot := lmap[i]
		vtot := vmap[i]/1e9
		if ltot == 0.0{continue}
		atot := amap[i]/1e6
		ntot := math.Ceil(lmap[i]/3048.0)
		row := fmt.Sprintf("%s,%.2f,%.3f,%.3f,%.3f,%.0f,%.0f,%.0f",grp,dimap[i],ltot/1e3,atot,vtot,ntot,vtot * t.Kostin[0] * 36.0,atot*t.Kostin[0])
		tcost += vtot * t.Kostin[0] * 36.0 + atot * t.Kostin[1]
		table.Append(strings.Split(row, ","))
	}
	vblt := 0.0
	
	for i := range t.Mod.Js{
		for _, bval := range t.Mod.Js[i].Fxtrs{
			vblt += bval[0] * bval[1] * bval[1] * math.Pi * bval[2]/4.0 
		}
	}
	vblt = vblt/1e9
	tcost += vblt * 7850.0 * t.Kostin[2]
	tcost = math.Ceil(tcost)
	table.Render()
	rezstr.WriteString(ColorReset)
	rezstring := rezstr.String()//fmt.Sprintf("%s",rezstr)
	fmt.Println(rezstring, "\nCOST PER TRUSS -> ",tcost, " rupeeses", "\nTOTAL COST -> ",tcost * t.Ntruss," rupeeses")
	t.Mod.Reports = append(t.Mod.Reports, rezstring)
	t.Mod.Reports = append(t.Mod.Reports, colreps)
	//if 1 == 2 {log.Println(colreps)}
	
	return
}

//ReadPrp reads design values for a column struct
//use when non standard group properties need to be read in
func (c *WdCol) ReadPrp(dzval []float64){
	//read dzvals
	c.Prp.Em = dzval[0]
	c.Prp.Pg = dzval[1]
	
	//etc
}

//ColD designs an axially loaded wooden column
//using results from a kass.Model (only a 2d truss so far)
func ColD(mem *kass.Mem, styp, grp int, dims, dzval []float64, colchn chan []interface{}){
	//first build column
	//then check/design
	c := WdCol{
		Id: mem.Id,
		Grp: grp,
		Styp:styp,
		Lspan:mem.Geoms[0],
	}
	if c.Grp == 0 && len(dzval) != 0{
		c.ReadPrp(dzval)
	}
	c.Init()
	//increase pu by 1.15 for net section in bolted trusses
	c.Pu = mem.Cmax * 1.15
	if math.Abs(mem.Tmax) > c.Pu{
		c.Pu = math.Abs(mem.Tmax) * 1.15
		c.Tensile = true
	}
	
	err := ColDz(&c)
	rez := make([]interface{},3)
	rez[0] = err
	rez[1] = mem.Id
	rez[2] = c
	colchn <- rez
}

