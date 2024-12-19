package barf    

import (
	"log"
	"fmt"
	"math"
	"bufio"
	"os"
)

//chapter 4 of shah routines
//steel bar dias 
var StlDia = []float64{6.0,8.0,10.0,12.0,16.0,20.0,25.0,28.0,32.0,40.0}
//steel bar areas 
var stlAreas = []float64{28.27,50.27,78.54,113.10,201.06,314.16,490.87,615.75,804.25,1256.64}
//steel density
var stlDens float64 = 7850
//rcc density (has never been used anywhere)
var rccDens float64 = 2500
//rcc grades
var fckarr = []float64{10,15,20,25,30,35,40,45,50,55,60}
//steel grades
var fyarr = []float64{415,500,550,600}
//base rcc m20
//var fck float64 = 20

//base steel
var fy float64 = 415
//young's modulus of steel
var Esteel float64 = 2e5
//safety factor for material for steel
var psfs float64 = 1.15

//safety factor for material for concrete
var psfck float64 = 1.5

//max percentage moment redistribution
var dMmaxallow = 0.3

//is code stress block constants k1 and k2
var isSBk1 float64 = 0.36
var isSBk2 float64 = 0.416
//bs code rct stress block constants 
//max size of coarse aggregate in mm
var msca float64 = 20

//interpolator interpolates the value of y at x given slices of ys and xs
func interpolator(x float64, xs []float64, ys []float64) (y float64) {
	//there are thousands of the same function floating around
	var y0, y1, x0, x1 float64
	if x > xs[len(xs)-1]{y = ys[len(ys)-1];return}
	for idx, val := range xs {
		if val > x && idx == 0{
			//fmt.Println("YEEHAW->",val,x,ys[0])
			y = ys[0]
			return
		}
		if val == x {
			y = ys[idx]
			return
		}
		if val > x && idx != 0 && xs[idx-1] < x{
			y0 = ys[idx-1]; y1 = ys[idx]
			x0 = xs[idx-1]; x1 = xs[idx]
			y = y0 + ((x-x0)*(y1-y0)/(x1-x0))
			return
		}
	}
	return
}

//RbrFrcIs interpolates fsc value for a given strain
//interpolates Fe250 and Fe415 curves for Fe550 curve ('')
func RbrFrcIs(fy, esc float64) (fsc float64) {
	var exs, fys []float64
	//fmt.Println("in esc->",esc)
	switch {
	case fy == 250.0 || esc <= 0.8*fy/230000.0:
		fsc= esc*230000.0
	case fy == 415.0:
		exs = []float64{0.0014435, 0.00163, 0.00192,0.00241,0.00276,0.00380}
		fys = []float64{288.7,306.7,324.8,342.8,351.8,360.9}
		fsc = interpolator(esc,exs,fys)
	case fy == 500.0, fy == 460.0:
		exs = []float64{0.001739,0.00195,0.00226,0.00277,0.00312,0.00417} 
		fys = []float64{347.8,369.6,391.3,413.0,423.9,434.8}
		fsc = interpolator(esc,exs,fys)
	case fy == 550.0:
		exs = []float64{0.001912824, 0.002138235, 0.00246,0.002981765,0.003331765,0.004387647} 
		fys = []float64{382.8,406.725,430.65,454.575,466.5375,478.5}
		fsc = interpolator(esc,exs,fys)		
	default:
		log.Println("ERRORE,errore->steel grade greater than Fe550")
		fsc = -99.0
	}
	if fsc > fy/1.15 {
		fsc = fy/1.15
	}
	//fmt.Println("out fsc->",fsc)
	return 
}

//RbrArea returns the area of a bar of diameter dia
func RbrArea(dia float64) (area float64) {
	area = math.Pi * math.Pow(dia,2)/4.0
	return
}

//RbrUnitWt returns the weight/running meter of a bar of dia dia
func RbrUnitWt(dia float64) (weight float64) {
	weight = RbrArea(dia) * stlDens/(1000*1000)
	return
}

//RbrSingle calcs the number of bars of a single dia for a (required) ast
func RbrSingle(ast float64)(rez [][]float64){
	//when only one dia is needed (use for ribbed slabs)
	
	ia := 2; ib := 9; n2 := 0.0; d2 := 0.0
	for i := ia; i < ib; i++{
		dia := StlDia[i]
		area := RbrArea(dia)
		n1 := int(math.Ceil(ast/area))
		asprov := float64(n1)*area
		adiff := asprov - ast
		rez = append(rez, []float64{float64(n1), float64(n2), dia, d2, asprov, ast, adiff})
	}
	return
}


//RbrSelect calcs the number of bars of a single dia for a (required) asreq
//usr - true enables a menu
func RbrSelect(asreq float64, usr bool) (rez [][]float64, mindia float64) {
	var dimin, dimax, asprov, adiff , asvar float64
	var mxnbar, nbar , ia, ib int
	if usr{
		var r = bufio.NewReader(os.Stdin)
		fmt.Println("Enter min dia, max dia and max no. of bars")
		fmt.Fscanf(r, "%f %f %d",&dimin, &dimax, &mxnbar)
		for i := 0; i < len(StlDia); i++ {
			if StlDia[i] == dimin {ia = i}
			if StlDia[i] == dimax {ib = i}
		}
	} else {
		dimin = 10
		dimax = 32
		mxnbar = 15
		ia = 2
		ib = 9
	}
	//minidx = -1.0
	adiff = 9999.99
	for i := ia; i < ib; i++{
		dia := StlDia[i]
		area := RbrArea(dia)
		nbar = int(math.Ceil(asreq/area))
		asprov = float64(nbar)*area
		asvar = asprov - asreq
		if nbar < mxnbar{
			rez = append(rez, []float64{float64(i), dia, float64(nbar), asprov, asvar})
			if asvar < adiff{
				adiff = asvar
				mindia = StlDia[i]
			}
		}
	}
	return
}

//RbrNRows returns the number of rows and bars per row for an n-d combo
//rez = []float64{nlayer, astprov, efcvr, efdp, cldis, clvdis, nbarRow}
//CHANGE THIS TO INCLUDE NOM CVR
func RbrNRows(bw float64, dused float64, nbar float64, dia float64) (rez []float64){
	var cldis float64 = 25.0
	var clrCvr float64 = 30.0
	var clvdis, efdp, efcvr float64
	if dia  > cldis {cldis = dia}
	if dia  > clrCvr {clrCvr = dia}
	//HYARE
	//if nbar == 1.0{nbar = 2.0}
	clvdis = 2.0 * msca/3.0
	nbarRow := math.Floor((bw - 2.0 * clrCvr)/(dia + cldis))
	nlayer := math.Ceil(nbar/nbarRow)
	xstep := math.Floor(bw - 2.0 * clrCvr)/nbarRow
	if nlayer > 1.0{
		//clvdis = 2.0 * msca/3.0
		if clvdis < 15.0 {clvdis = 15.0}
		if dia > clvdis {clvdis = dia}
	}
	switch int(nlayer){
		case 1:
		efcvr = clrCvr + dia/2.0
		case 2:
		efcvr = clrCvr + dia + clvdis/2.0
		case 3:
		efcvr = clrCvr + 1.5 * dia + clvdis  
		case 4:
		efcvr = 666.0
	}
	efdp = dused - efcvr
	astprov := nbar * RbrArea(dia)
	rez = []float64{nlayer, astprov, efcvr, efdp, xstep, clvdis + dia, nbarRow}
	return
}

//FtngRbrDia returns rebar dia spacing combos for x (astx) and y (asty) footing directions 
func FtngRbrDia(dy, fck, fy, colx, coly, lx, ly, astx, asty, dreq, nomcvr float64) (rez [][]float64, mdx int, err error){
	dy, lx, ly, nomcvr =  dy * 1e3, lx * 1e3, ly * 1e3, nomcvr * 1e3
	var astd float64
	for _, dia := range StlDia{
		if dia < 8.0 {continue}
		//if dia > dreq{
		//	break
		//}
		nx := math.Ceil(astx/RbrArea(dia))
		efcvr := nomcvr + dia/2.0
		spcx := math.Floor((ly - 2.0 * efcvr)/(nx-1.)/10.0)*10.0
		if spcx < 50.0{
			continue
		}
		asx := nx * RbrArea(dia)
		ldx := (ly - nomcvr - dia/2.0)
		if ok := FtngDevLen(fck, fy, dia, ldx + dy); !ok{continue}
		ny := math.Ceil(asty/RbrArea(dia))
		//spcy := math.Round((ly - 50.)/(ny - 1.) - dia)
		spcy := math.Floor((lx - 2.0 * efcvr)/(ny-1.)/10.0)*10.0
		if spcy < 50.0{
			continue
		}
		ldy := (lx - nomcvr - dia/2.0)
		if ok := FtngDevLen(fck, fy, dia, ldy + dy); !ok{continue}
		asy := ny * RbrArea(dia)
		rez = append(rez, []float64{dia, nx, spcx, asx, asx- astx, ny, spcy, asy, asy- asty})
		//log.Println("appended rez->",rez, asx, asy)
		delta := asx - astx + asy - asty
		if len(rez) == 1{
			astd = delta
			mdx = len(rez) -1
		} else if astd >= delta{
			astd = delta
			mdx = len(rez) -1
		}
	}
	if len(rez) == 0{
		err = ErrFtngDim
		return
	}
	err = nil
	return
}

//FtngDevLen checks for the development length required for a footing
//tbh this is just dev len? 
func FtngDevLen(fck, fy, dia, ldav float64) (bool){
	var tbd float64
	tbdvals := map[float64]float64{15.0:1.0,20.0:1.2,25.0:1.4,30.0:1.5,35.0:1.7,40.0:1.9}
	if val, ok := tbdvals[fck]; !ok {
		if fck > 40.0 {
			tbd = 1.9
		}
		if fck < 15.0 {
			tbd = 0.8//fuck this
		}
	} else {
		tbd = val                    
	}
	if fy > 250.0 {
		tbd = 1.6*tbd
	}
	ldfac := fy/(4.0*tbd*1.15)
	ldreq := ldfac * dia
	//log.Println(ColorRed,"ldreq,ldav",ldreq, ldav,ldreq<=ldav,ColorReset)
	return ldreq <= ldav
}

//SlabRbrDia returns slab rebar dia-spacing combos
func SlabRbrDia(dused float64, asreq float64) (rez [][]float64, minidx int) {
	//main steel alwayze 8mm min here
	var dimaxsl, smin, smax, efcvr , efdp float64
	var idxmax int
	dimaxsl = math.Floor(dused/8)
	//fmt.Println("max dia dimaxsl - ", dimaxsl)
	for i:= 1; i <len(StlDia) ; i++ {
	if StlDia[i] > dimaxsl {
			idxmax = i -1
			break
		}
	}
	//smin, smax - min and max spacing of bars
	smin = 80
	smax = 450
	efcvr = 20
	efdp = dused - efcvr
	if smax > 3 * efdp {smax = 3 * efdp}
	adiff := 1000.0
	for i := 1; i <= idxmax; i++ {
		dia := StlDia[i]
		area := stlAreas[i]
		spcing := math.Floor(1000*area/asreq)
		if spcing < smin {i++; continue}
		if spcing > smax {spcing = smax}
		asprov := 1000 * area / spcing
		prevdiff := asprov - asreq
		if asprov > asreq {
			rez = append(rez,[]float64{dia, spcing, asprov, prevdiff})
			if adiff > prevdiff {
				adiff = prevdiff
				minidx = i
			}
		}
		//fmt.Println("Bar dia ",dia," Spacing ",spcing," Area Provided ",asprov," Difference ",prevdiff)
		if (i == idxmax && len(rez) == 0) {
			fmt.Println("Change depth and ergo, area \n CONCORDANTLY "); return}
	}
	/*
	if len(rez) != 0 {
		fmt.Println("DEPTH AND AREA REQUIRED :", dused, " ",asreq )
		fmt.Println("min. index values : \n","Bar dia ",rez[minidx][0]," Spacing ",rez[minidx][1]," Area Provided ",rez[minidx][2]," Difference ",rez[minidx][3])
	}
	*/
	return	
}

//SlabRbrDiaSpacing returns slab rebar-dia combos as a map of dia-rez
func SlabRbrDiaSpcing(dused , asreq, efcvr float64) (rezmap map[float64][]float64, mindia float64) {
	//main steel alwayze 8mm min here
	var rdias = []float64{6.0,8.0,10.0,12.0,16.0,18.0,20.0,22.0,25.0,28.0,32.0,40.0}
	//steel bar areas
	var rareas = []float64{28.27,50.27,78.54,113.10,201.06,314.16,380.13,490.87,615.75,804.25,1256.64}

	var dimaxsl, smin, smax, efdp float64
	var idxmax int
	rezmap = make(map[float64][]float64)
	dimaxsl = math.Round(dused/8)
	
	for i:= 1; i <len(rdias) ; i++ {
	if rdias[i] > dimaxsl {
			idxmax = i -1
			break
		}
	}
	//WOT IF MIN SPACING CHANGES
	smin = 75
	smax = 300
	efdp = dused - efcvr
	if smax > 3.0 * efdp {smax = 3.0 * efdp}
	adiff := 1000.0
	for i := 1; i <= idxmax; i++ {
		dia := rdias[i]
		area := rareas[i]
		spcing := 1000.0*area/asreq
		if spcing < smin {continue}
		if spcing > smax {spcing = smax}
		spcing = 10.0 * math.Floor(spcing/10.0)
		asprov := 1000 * area / spcing
		prevdiff := asprov - asreq
		if asprov > asreq {
			rezmap[dia] = []float64{dia, spcing, asprov, prevdiff}
			if adiff > prevdiff {
				adiff = prevdiff
				mindia = dia
			}
		}
	}
	return
}

//RccSlabServeChk - fig 4.14 shah; serviceability check for l/d ratio of slab
//rounds to nearest 5 mm multiple

func RccSlabServeChk(slabtyp, endc int, fy, ll, lx, ptreq,efcvr float64) (dreq float64){
	var alphaarr []float64
	ptstlarr := []float64{0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1.0,1.2,1.4,1.6,1.8}
	switch {
	case fy == 250.0:
		alphaarr = []float64{2.0,2.0,2.0,2.0,1.90,1.75,1.63,1.54,1.47,1.41,1.32,1.25,1.2,1.16}		
	case fy >= 415.0:
		alphaarr = []float64{2.0,1.6,1.4,1.28,1.18,1.11,1.06,1.02,0.99,0.96,0.92,0.89,0.87,0.85}	
	}
	var slope, alpha, allowld float64
	switch {
	case slabtyp == 2:
		//two way slab
		if lx <= 3500.0 && ll <= 3.0 {
			if endc >= 9 {
				allowld = 35.0
			} else {
				allowld = 40.0
			}
			if fy > 250.0 {
				allowld = 0.8 * allowld
			}
		} else {
			if endc >=9 {
				allowld = 20.0
			} else {
				allowld = 26.0
			}		
		}
	case slabtyp == 1:
		//un way
		if endc == 1 {
			allowld = 20.0
		} else if endc == 2 {
			allowld = 26.0
		} else {
			allowld = 7.0
		}
	}
	for idx, val := range ptstlarr {
		if val == ptreq {
			alpha = alphaarr[idx]
			break
		}
		if val > ptreq {
			if idx > 0 {
				x1 := ptstlarr[idx-1]
				x2 := ptstlarr[idx]
				y1 := alphaarr[idx-1]
				y2 := alphaarr[idx]
				slope =  (y2-y1)/(x2-x1)
				alpha = y1 + slope*(ptreq - x1)
				break
			} //add for idx == 0
			if idx == 0 {
				alpha = alphaarr[idx]
			}
		}
	}
	allowld = allowld * alpha
	//dreq = 10.0 * math.Ceil(((lx/allowld) + efcvr)/10.0)
	dreq = lx/allowld + efcvr
	return 
}

//SlabDevLen checks slab rebar for development length
func SlabDevLen(fck, fy, m1, v, effd, dia, bs float64, endc int) (bool, float64, float64){
	var tbd float64
	tbdvals := map[float64]float64{15.0:1.0,20.0:1.2,25.0:1.4,30.0:1.5,35.0:1.7,40.0:1.9}
	if val, ok := tbdvals[fck]; !ok {
		if fck > 40.0 {
			tbd = 1.9
		}
		if fck < 15.0 {
			tbd = 0.8//lot wot CHECK THIS
		}
	} else {
		tbd = val                    
	}
	if fy > 250.0 {
		tbd = 1.6*tbd
	}
	ldfac := fy/(4.0*tbd*1.15)
	ldreq := ldfac * dia

	var mf, lex, ldav float64
	var rcom, barcut int
	if endc == 1 {
		rcom = 0; barcut = 1
	} else {
		rcom = 1; barcut = 1
	}
	switch rcom {
	case 0://compressive reaction
		mf = 1.3
	case 1://tensile reaction 
		mf = 1.0	
	}
	switch barcut {
	case 0://point of curtailment
		if dia > 12.0 {
			lex = dia
		} else {
			lex = 12.0
		}
	case 1:
		lex = ldreq/3.0 -bs/2.0
	}
	
	ldav = lex + 1000*m1*mf/v

	return ldav >= ldreq, ldav, 5.0 * math.Ceil(ldreq/5.0)
}

//RccDevLenCheck is a general function for development length checks
//as seen in shah chapter 4
func RccDevLenChk(fck, fy, m1, v, effd, dia, bs float64, rcom, barcut int) (bool){
	var tbd float64
	tbdvals := map[float64]float64{15.0:1.0,20.0:1.2,25.0:1.4,30.0:1.5,35.0:1.7,40.0:1.9}
	if val, ok := tbdvals[fck]; !ok {
		if fck > 40.0 {
			tbd = 1.9
		}
		if fck < 15.0 {
			tbd = 0.8//lot wot CHECK THIS
		}
	} else {
		tbd = val
	}
	if fy > 250.0 {
		tbd = 1.6*tbd
	}
	ldfac := fy/(4.0*tbd*1.15)
	ldreq := ldfac * dia

	var mf, lex, ldav float64

	switch rcom {
	case 0://tensile reaction
		mf = 1.3
	case 1://compressive reaction 
		mf = 1.0	
	}
	switch barcut {
	case 0://point of curtailment
		if dia > 12.0 {
			lex = dia
		} else {
			lex = 12.0
		}
	case 1:
		lex = ldreq/3 -bs/2
	}
	
	ldav = lex + 1000*m1*mf/v
	
	return ldav >= ldreq
}

//RbrShearLink returns shear link dia, nlegs and spacing as in chapter 4, shah
func RbrShearLink(fy, bw, effd, vusv, svmin float64, minidx, maxidx, mxleg, dnsr int) (idxs, nlegs []int, spacing []float64){
	/*
	   shear stirrups for vusv as per shah
	*/
	fyi := fy
	svmax := 0.75 * effd
	if svmax > 300.0 {
		svmax = 300.0
	}
	//
	var sv float64
	for j := 2; j <= mxleg; j += 2 {
		for i := minidx; i <= maxidx; i++ {
			if i == 0 {
				fy = 250.0
			} else {
				fy = fyi
			}
			switch dnsr {
			case 0: //design for min shear reinforcement	
				sv = 10 * math.Floor(fy* float64(j) * RbrArea(StlDia[i]) *2.5/(10*bw))
				//sv = 10.0*(sv/10)+(10 - int(sv)%10)
			case 1: //design for actual shear
				sv = 10 * math.Floor(0.87* fy * float64(j) * stlAreas[i]*effd/(10*vusv*1000.0))
				//sv = 10.0*(sv/10)+(10 - sv%10)
			}
			if sv > svmin && sv <= svmax {
				idxs = append(idxs,i)
				nlegs = append(nlegs,j)
				spacing = append(spacing,sv)
			}
			
		}
	}
	return
}

//SolveQuad returns the +ve root of  the quadratic ax2 + bx + c = 0
func SolveQuad(a, b, c float64) (xu float64){
	d := math.Pow(b,2.0) - 4.0*a*c
	if d < 0 {xu = -99.0} else {xu = (math.Sqrt(d) - b)/(2.0*a)}
	return
}

//SolveCubic solves the cubic equation WHAAT?
//this should be deleted just because one does not even remember what the equation is
func SolveCubic(a2, a1, a0 float64) (z1 float64){
	q := a1/3.0 - math.Pow(a2, 2)/9.0
	r := (a1 * a2 - 3.0 * a0)/6.0 - math.Pow(a2, 3)/27.0
	det := math.Pow(r,2) + math.Pow(q,3)
	if det >= 0{
		a := math.Pow(math.Abs(r) + math.Sqrt(math.Pow(r,2)+math.Pow(q,2)), 1.0/3.0)
		var t1 float64
		if r >= 0{
			t1 = a - q/a
		} else {
			t1 = q/a - a
		}
		z1 = t1 - a2/3.0
		
	} else {
		z1 = -99.9
	}
	return
}

//Solve DCubic returns the root of the depressed cubic equation x3 + px + q = 0
func SolveDcubic(p, q float64)(z float64){
	if 4.0 * math.Pow(p,3.0) + 27.0 * math.Pow(q,2.0) > 0.0 {
		t := math.Sqrt(q*q/4.0 + p*p*p/27.0)
		z = math.Pow(-q/2.0 + t, 1.0/3.0) + math.Pow(-q/2.0 - t, 1.0/3.0)
	}
	return
}

//BalSecKis returns design constants for a balanced section as per is code
//see section 4.3, shah
func BalSecKis(fck, fy, dm, psfs, k1, k2 float64) (kumax, zumax, rumax, ptmax float64) {
	//bal section constants section 4.3 shah
	//returns factors for neutral axis kumax, lever arm zumax, moment rumax, steel ptmax 
	fyd := fy/psfs
	kumax = 700.0/(1100.0+ fyd)
	kulim := 0.6 - (dm)
	if kulim < kumax {
		kumax = kulim
	}
	zumax = 1.0 - (k2 * kumax)
	rumax = k1*fck*kumax*zumax
	ptmax = k1*fck*kumax/fyd
	return 
}
