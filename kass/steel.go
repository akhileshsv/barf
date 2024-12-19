package barf

import (
	"os"
	"log"
	"fmt"
	"math"
	"runtime"
	"errors"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"github.com/go-gota/gota/dataframe"
)

var (
	StlSnames = map[int]string{
		1:"i", 2:"c",3:"l",4:"ln",5:"t",6:"box",7:"tube",8:"ub",9:"uc",10:"flat",11:"bar",
		101:"built-i",102:"l2-ss",103:"l2-os",104:"ln2-ss",105:"ln2-os",106:"plate-i",
	}
	StlStyps = map[string]int{
		"i":1,"c":2,"l":3,"ln":4,
		"t":5,"box":6,"tube":7,
		"ub":8,"uc":9,"rect":10,
		"circle":11,"bar":11,
		"flat":10,"isa":3,
		"built-i":101,
		"l2-ss":102,"l2-os":103,
		"ln2-ss":104,"ln2-os":105,
		"plate-i":106,
	}
	StlBstyps = map[string][]int{
		"l":[]int{3,7},
		"ln":[]int{4,7},
		"built-i":[]int{1,12},
		"l2-ss":[]int{3,27},
		"l2-os":[]int{3,28},
		"ln2-ss":[]int{4,27},
		"ln2-os":[]int{4,28},
	}
	StlSdxLims = map[string]int{
		"i":64,"c":26,"l":71,"ln":64,
		"t":0,"box":0,"tube":81,
		"ub":71,"uc":30,"rect":0,
		"circle":0,"bar":0,
		"flat":0,"isa":0,
		"built-i":0,
		"l2-ss":0,"l2-os":0,
		"ln2-ss":0,"ln2-os":0,
		"plate-i":0,
	}
	Table2Bs = []float64{180,165,230,215,280,170,155,215,200,265,185}	
	//stlsecmap = map[int]string{7:"UB",8:"UC"}
	//pqcol - from mosley section 6
	PqCol = map[int][]float64{43:{155,165,165},50:{215,230,230},55:{265,280,280}}
	//pqbm - qx (bending), ps (shear), pc (web crushing)
	PqBm = map[int][]float64{43:{165,100,190},50:{230,140,260},55:{280,170,320}}
	PqCol89 = map[int][]float64{43:{170,180,180},50:{215,230,230},55:{265,280,280}}
	EStl = 210000.0
	EStl89 = 205000.0
)

//baah define stlsec struct to hold stlsec named params
//add stuff AT WILL
type StlSec struct{
	Sname                    string
	Sstr                     string
	Frmstr                   string
	Mstr                     string //memstr (col/beam/prln)
	Weld                     bool //is welded 
	Lsb, Lsbx, Lsby          bool //is laterally supported
	Scrit                    bool //shear critical
	Sbuck                    bool //check for shear buckling (if d/tw > 67e)
	Lvl                      int  //level (simple/basic/complex like justice)
	Sectyp                   int
	Endc                     int
	Ctyp                     int //connection type
	Jtyp                     int //joint type 
	Ltyp                     int
	Dtyp                     int //eh. D typ!
	Cleg                     int //con. leg. 0/1 - long leg, 2 - short leg
	Ncon                     int //connection number/nbolts/etc
	Sdx                      int
	Code                     int
	Ax                       int //0 - x, 1 - y, 2 - biaxe
	Pfac                     float64
	Grd                      float64
	Wt                       float64
	Bstyp                    int
	Em, Gm                   float64
	H, B, Tw, Tf             float64
	H1, H2                   float64 //height of col above and beyond
	Vbdx, Vbdy               float64 //diff of beam shears in x and y
	Lx, Ly, Tx, Ty           float64 //eh use (unbraced lengths) lx, ly, fixity factors 
	R1, R2                   float64
	Dw, Eps                  float64
	Area, Ixx, Iyy, Rxx, Ryy float64
	Cxx, Cyy, Zxx, Zyy       float64
	Zfd, Zfdx, Zfdy          float64
	Iww, Itt                 float64
	D1, D2, B1, B2           float64
	Bfrat, Hwrat             float64
	Defrat                   float64
	Fy, Fu                   float64
	Fcc, Fccx, Fccy          float64
	Fcdx, Fcdy, Fcd          float64
	Fwb, Fwc                 float64
	Lspan, Leff, Klx, Kly    float64
	Leffx, Leffy             float64
	Lbr, Tbr                 float64
	Fbd                      float64
	Tplate                   float64
	Bxx, Byy                 int
	Ocf, Icf, Wcf            int
	Sxx, Syy                 float64
	Srfx, Srfy               float64
	Pu, Pur                  float64
	Pud                      float64
	Mux, Muy, Mur, Mcr, Mdv  float64
	Mdx, Mdy                 float64
	Mndx, Mndy               float64
	Cbeq                     float64
	Vfac, Vux, Vuy, Vur      float64
	Vdx, Vdy, Vd             float64
	Vu, Mu, Dmax             float64
	Avx, Avy                 float64
	Tur                      float64
	Tdg, Tdn, Tdb1, Tdb2     float64
	Sleg, Anc, Ago, Bago     float64
	Avg, Atg, Avn, Atn       float64
	Ymo, Yml                 float64 
	Term                     string
	Sec                      SectIn
	Bg                       Blt
	Wg                       Wld
	Mem                      Mem
	Bdia                     float64
	Wsize                    float64
	Kbs                      [][]float64
	Dims                     []float64
}

//PbcBs returns permissible bending stresses as per bs449
func PbcBs(sname string, grade int) (vec [][]float64, err error){
	var mpbc map[string][][]float64
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	jsonin := filepath.Join(basepath,"../data/steel/bsteel","pbc.json")
	jsonfile, e := ioutil.ReadFile(jsonin)
	if e != nil {
		err = fmt.Errorf("error reading bs449 perm. stress data %s",e)
		return
	}
	err = json.Unmarshal([]byte(jsonfile), &mpbc)
	if err != nil {
		err = fmt.Errorf("error unmarshalling json %s",err)
		return
	}
	var query string
	switch sname{
		case "ub", "uc":
		//uc beams and columns
		switch grade{
			case 43:
			query = "3a" 
			case 50:
			query = "3b"
			case 55:
			query = "3c"
		}
	}
	vec = mpbc[query]
	return
}

//PbcLerp linearly interpolates the permissible bending stress given a slenderness ratio
//calls PbcBs for the table of permissible bending stresses
func PbcLerp(sname string, grd int, s1, dtrat float64) (pbc float64, err error){
	//log.Println("lerp in-> srat, drat->",s1, dtrat)
	pbvec, e := PbcBs(sname, grd)
	if e != nil{
		err = e
		return
	}
	pvec := PqCol[grd]
	var rdx, cdx int
	switch {
	case s1 <= 40.0:
		pbc = pvec[1]
		return
	case s1 <= 120.0:
		cdx = int((s1 - 40.0)/5.0)
	case s1 <= 300.0:
		cdx = int((s1-120.0)/10.0)+(120-40)/5
	}
	rdx = int((dtrat-5)/5.0)+1
	//log.Println("dxs->",rdx, cdx)
	sa := pbvec[0][cdx]; sb := pbvec[0][cdx+1]
	//log.Println("dtrats->",rdx*5, (rdx+1)*5)
	//log.Println("srats->",sa,sb)
	//log.Println("rdx, cdx->",rdx,cdx)
	pt0 := pbvec[rdx][cdx]; pt1 := pbvec[rdx+1][cdx]
	//log.Println("pts 1",pt0,pt1)
	p1 := pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	if cdx == 33{
		//log.Println(len(pbvec), len(pbvec[0]))
		pbc = p1
		return
	}
	pt0 = pbvec[rdx][cdx+1]; pt1 = pbvec[rdx+1][cdx+1]
	//log.Println("pts 2",pt0,pt1)
	p2 := pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	pbc = p1 + (s1 - sa) * (p2 - p1)/(sb - sa)
	return
}

//PbcYeolde returns the permissible bending stress as in table 6.1 of mosley/spencer (ye olde values)
func PbcYeolde(s1, dtrat float64) (pbc float64){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	sheet := filepath.Join(basepath,"../data/steel","hulsepbc43.csv")
	csvfile, err := os.Open(sheet)
	if err != nil {
		log.Fatal(err)
		return
	}
	df := dataframe.ReadCSV(csvfile)
	var rdx, cdx int
	switch{
		case s1 <= 90:
		pbc = df.Elem(0,1).Float()
		return
		case s1 <= 120:
		rdx = int((s1 - 90.)/5.0)
		default:
		rdx = int((s1 - 120.)/10.0) + (120 - 90)/5
	}
	if dtrat <= 10{dtrat = 10}
	if dtrat <= 40{
		cdx = int((dtrat - 10.)/5.0) + 1
	} else {
		cdx = 7
	}
	//var sa, sb float64
	sa := df.Elem(rdx,0).Float(); sb := df.Elem(rdx+1,0).Float()
	//log.Println("sa, sb->",sa, sb, rdx)
	var p1, p2 float64
	pt0 := df.Elem(rdx,cdx).Float(); pt1 := df.Elem(rdx,cdx+1).Float()
	if cdx < 7 {
		p1 = pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	} else {
		p1 = pt0 + (dtrat-40.)*(pt1 - pt0)/10.0
	}
	//log.Println(pt0, pt1)
	pt0 = df.Elem(rdx+1,cdx).Float(); pt1 = df.Elem(rdx+1,cdx+1).Float()
	if cdx < 7 {
		p2 = pt0 + math.Mod(dtrat,5.0)*(pt1 - pt0)/5.0
	} else {
		p2 = pt0 + (dtrat-40.)*(pt1 - pt0)/10.0
	}
	
	//log.Println(pt0, pt1)
	//log.Println(p1, p2)
	pbc = p1 + (s1 - sa) * (p2 - p1)/(sb - sa)
	return
}

//Printz printz 
func (ss *StlSec) Printz(){
	ss.Sec.Draw("dumb")
	rezstr := ss.Sec.Txtplot

	rezstr += fmt.Sprintf("section name - %s\ndims -  %f (mm)\narea - %f (mm2) weight - %f n/m\nixx %f - (mm4) iyy - %f (mm4)\nrxx - %f (mm3) ryy - %f (mm3)\nsxx - %f (mm4) syy - %f (mm4)\nzxx - %f (mm4) zyy - %f (mm4)\n",ss.Sstr,ss.Dims, ss.Area, ss.Wt,ss.Ixx, ss.Iyy, ss.Rxx, ss.Ryy, ss.Sxx, ss.Syy, ss.Zxx, ss.Zyy)
	fmt.Println(rezstr)
}

//Draw plots a StlSec w/dimensions
func (ss *StlSec) Draw()(err error){
	s := SecGen(ss.Bstyp, ss.Dims)
	if s.Prop.Area == 0.0{
		err = fmt.Errorf("invalid section specs styp - %v dims - %v",ss.Bstyp, ss.Dims)
		return
	}
	s.Draw(ss.Term)
	fmt.Println(s.Txtplot)
	return
}

//CalcZyy updates a sections plastic modulus about the minor axis
//SEE https://en.wikipedia.org/wiki/Section_modulus Z = Acyc + Atyt
func (ss *StlSec) CalcZyy()(err error){
	//rotate sec
	sr := SecRotate(ss.Sec, 90.0)
	switch ss.Sname{
		case "built-i":
		area, _, yc := SecArXu(&sr, sr.Ymx - sr.Prop.Yc)
		//is symmetric about centroid, so z = 2 * area * (sr.Prop.Yc - yc)
		ss.Zyy = 2.0 * area * math.Abs(sr.Prop.Yc - yc)
	}
	//profit?
	return
}

//SecGen generates a StlSec's section given a base sectype and dims
func (ss *StlSec) SecGen()(err error){
	ss.Sec = SecGen(ss.Bstyp, ss.Dims)
	if len(ss.Sec.Dims) == 0{
		err = fmt.Errorf("invalid section params - %v - %f",ss.Bstyp, ss.Dims)
		return
	}
	switch ss.Sname{
		case "built-i":
		//fmt.Println(ColorRed, ss.Dims,ColorReset)
		ss.B = ss.Sec.Dims[0]
		ss.H = ss.Sec.Dims[1]
		ss.Tf = ss.Sec.Dims[2]
		ss.Tw = ss.Sec.Dims[3]
		ss.Sstr = fmt.Sprintf("I-%f %f %f %f",ss.B,ss.H,ss.Tf,ss.Tw)
		ss.Area = ss.Sec.Prop.Area
		ss.Wt = ss.Area * 7850.0/1000000.0
		ss.Ixx = ss.Sec.Prop.Ixx
		ss.Iyy = ss.Sec.Prop.Iyy
		ss.Rxx = ss.Sec.Prop.Rxx
		ss.Ryy = ss.Sec.Prop.Ryy
		ss.Sxx = ss.Sec.Prop.Sxx
		ss.Syy = ss.Sec.Prop.Syy
		ss.Zxx = ss.Sec.Prop.Zxx
		ss.Zyy = ss.Sec.Prop.Zyy
		ss.Cxx = ss.Sec.Prop.Xc
		ss.Cyy = ss.Sec.Prop.Yc
		ss.Bfrat = ss.B/ss.Tf/2.0
		ss.Hwrat = (ss.H-2.0 * (ss.Tf+ss.R1))/ss.Tw
		//section modulus of flange(s) -  zp - aw * yw
		//yw = h/2/2 = h/4; duggal pg. 482	
		ss.Zfd = ss.Zxx - (ss.H * ss.Tw)* ss.H/4.0
		ss.Zfdx = ss.Zfd
		case "i","ub","uc":
		// ss.Sxx = ss.Sec.Prop.Sxx
		// ss.Syy = ss.Sec.Prop.Syy
		//ss.Zxx = ss.Sec.Prop.Zxx
		//ss.Zyy = ss.Sec.Prop.Zyy
		ss.Bfrat = ss.B/ss.Tf/2.0
		ss.Hwrat = (ss.H-2.0 * (ss.Tf+ss.R1))/ss.Tw
		//section modulus of flange(s) -  zp - aw * yw
		//yw = h/2/2 = h/4; duggal pg. 482	
		ss.Zfd = ss.Zxx - (ss.H * ss.Tw)* ss.H/4.0
		ss.Zfdx = ss.Zfd
		//TODO zfdy
		case "l", "ln":
		ss.Sxx = ss.Sec.Prop.Sxx
		ss.Syy = ss.Sec.Prop.Syy
		ss.Zxx = ss.Sec.Prop.Zxx
		ss.Zyy = ss.Sec.Prop.Zyy
		ss.Cxx = ss.Sec.Prop.Xc
		ss.Cyy = ss.Sec.Prop.Yc
		
		default:
		ss.Area = ss.Sec.Prop.Area
		ss.Wt = ss.Area * 7850.0*9.81/1000000.0
		ss.Ixx = ss.Sec.Prop.Ixx
		ss.Iyy = ss.Sec.Prop.Iyy
		ss.Rxx = ss.Sec.Prop.Rxx
		ss.Ryy = ss.Sec.Prop.Ryy
		ss.Sxx = ss.Sec.Prop.Sxx
		ss.Syy = ss.Sec.Prop.Syy
		ss.Zxx = ss.Sec.Prop.Zxx
		ss.Zyy = ss.Sec.Prop.Zyy
		ss.Cxx = ss.Sec.Prop.Xc
		ss.Cyy = ss.Sec.Prop.Yc
	}
	return	
}

//BckCls800 sets the buckling class of a section as per table 10 of is 800 2007
func (ss *StlSec) BckCls800() (err error){
	switch ss.Sname{
		case "pipe","tube","box":
		ss.Bxx = 1
		ss.Byy = 1
		case "c","ismc":
		ss.Bxx = 3
		ss.Byy = 3
		case "l", "ln":
		ss.Bxx = 3
		ss.Byy = 3		
		case "l2-ss", "l2-os", "ln2-ss", "ln2-os":
		ss.Bxx = 3
		ss.Byy = 3
		case "plate-i":
		ss.Bxx = 3
		ss.Byy = 3
		case "built-i":
		switch{
			case ss.Tf <= 40.0:
			ss.Bxx = 2 //class b
			ss.Byy = 3 //class c
			case ss.Tf > 40.0:
			ss.Bxx = 3 //class c
			ss.Byy = 4 //class d
		}
		case "i","ub","uc","ismb":	
		switch{
			case ss.H/ss.B > 1.2:
			switch{
				case ss.Tf <= 40.0:
				ss.Bxx = 1 //class a
				ss.Byy = 2 //class b
				case ss.Tf <= 100.0:
				ss.Bxx = 2 
				ss.Byy = 3 //class c
				default:
				err = fmt.Errorf("invalid flange/dims - %f\n",ss.Dims)
			}
			case ss.H/ss.B <= 1.2:
			switch{
				case ss.Tf <= 100.0:
				ss.Bxx = 2
				ss.Byy = 3
				case ss.Tf > 100.0:
				ss.Bxx = 4 //class d
				ss.Byy = 4
				default:
				err = fmt.Errorf("invalid flange/dims - %f\n",ss.Dims)
			}
		}
	}
	return
}			

//SecCls800 sets the class of a section as per table 2 of is 800 2007 (limiting w to thickness ratio)
func (ss *StlSec) SecCls800() (err error){
	ysr := math.Sqrt(250.0/ss.Fy)
	var k1, k2, k3 float64
	switch ss.Sname{
		case "i", "ub", "uc":
		k1 = 9.4 * ysr
		k2 = 10.6 * ysr
		k3 = 15.7 * ysr
		switch{
			case ss.Bfrat < k1:
			ss.Ocf = 1
			case ss.Bfrat < k2:
			ss.Ocf = 2
			case ss.Bfrat < k3:
			ss.Ocf = 3
			default:
			//slender section
			err = fmt.Errorf("slender section class b/tf %f vs k3 %f",ss.Bfrat, k3)
		}
		k1 = 84.0 * ysr
		k2 = 105.0 * ysr 
		k3 = 126.0 * ysr
		switch{
			case ss.Hwrat < k1:
			ss.Icf = 1
			case ss.Hwrat < k2:
			ss.Icf = 2
			case ss.Hwrat < k3:
			ss.Icf = 3
			default:
			err = fmt.Errorf("slender section class d/tw %f vs k3 %f",ss.Hwrat, k3)
		}
		if ss.Hwrat > 67.0 * ysr{ss.Sbuck = true} 
		case "built-i":
		//TODO TODO
	}
	
	return
}

//CalcTu calculates the ultimate tensile strength of a section
func (ss *StlSec) CalcTu()(err error){
	err = ss.NetArea()
	if err != nil{
		return
	}
	switch ss.Sname{
		case "l", "ln", "l2-ss", "l2-os", "ln2-ss", "ln2-os":
		//duggal 7, eq 6,7,10,11
		ss.Tdg = ss.Area * ss.Fy/ss.Ymo
		ss.Tdn = 0.9 * ss.Anc * ss.Fu/ss.Yml + ss.Bago * ss.Ago * ss.Fy/ss.Ymo 
		ss.Tdb1 = ss.Avg * ss.Fy/ss.Ymo/math.Sqrt(3) + 0.9 * ss.Atn * ss.Fu/ss.Yml
		ss.Tdb2 = ss.Atg * ss.Fy/ss.Ymo + 0.9 * ss.Avn * ss.Fu/ss.Yml/math.Sqrt(3)
		//fmt.Println("gross section yeilding tdg-",ss.Tdg/1e3,"kn")
		//fmt.Println("gross section rupture tdn-",ss.Tdn/1e3,"kn")
		//fmt.Println("block shear - shear yield and tension fracture",ss.Tdb1/1e3,"kn")
		//fmt.Println("block shear - tension yield and shear fracture",ss.Tdb2/1e3,"kn")
		ss.Tur = MinVal(ss.Tdg, ss.Tdn, ss.Tdb1, ss.Tdb2)
		//fmt.Println("UTS",ss.Tur/1e3,"kn")
	}
	return
}

//CalcMur calculates the ultimate bending strength of a section
func (ss *StlSec) CalcMur()(err error){
	//TODO checks - pg. 442, duggal
	//effects of bolt holes in flanges- 0.9 fu anf/yml >= fy agf * ymo
	//shear lag effect 
	err = ss.SecCls800()
	if err != nil{
		return
	}
	//3.7.2, is 800 2007 pg 17 - most critical element governs
	cfac := ss.Ocf; if ss.Icf > cfac{
		cfac = ss.Icf
	}
	bfac := 1.0
	switch cfac{
		case 1:
		//plastic
		case 2:
		//compact
		case 3:
		//semi compact
		bfac = ss.Sxx/ss.Zxx
	}
	switch ss.Lsb{
		case true:
		//laterally supported
		switch ss.Scrit{
			case false:
			//vu < 0.6 * vd
			md := bfac * ss.Zxx * ss.Fy/ss.Ymo
			switch ss.Endc{
				case 0:
				//clvr
				if md  >= 1.5 * ss.Sxx * ss.Fy/ss.Ymo{
					md = 1.5 * ss.Sxx * ss.Fy/ss.Ymo
				}
				case 1:
				//ss
				if md  >= 1.2 * ss.Sxx * ss.Fy/ss.Ymo{
					md = 1.2 * ss.Sxx * ss.Fy/ss.Ymo
				}
			}
			ss.Mur = md
			case true:
			//vu > 0.6 * vd
			//fmt.Println(ColorRed,"hi shear here",ColorReset)
			switch cfac{
				case 1,2:
				//plastic/compact section
				
				mfd := ss.Zfd * ss.Fy/ss.Ymo
				md := ss.Zxx * ss.Fy/ss.Ymo
				bfac = math.Pow(2.0 * ss.Vu/ss.Vd - 1.0, 2.0)
				ss.Mdv = md - bfac * (md - mfd)
				if ss.Mdv > 1.2 * ss.Sxx * ss.Fy/ss.Ymo{
					ss.Mdv = 1.2 * ss.Sxx * ss.Fy/ss.Ymo
				}
				//fmt.Println("mfd, md, bfac, mdv", mfd/1e6, md/1e6, bfac, ss.Mdv/1e6)
				ss.Mur = ss.Mdv
				case 3:
				//semi-compact section
				ss.Mdv = ss.Sxx * ss.Fy/ss.Ymo
				ss.Mur = ss.Mdv
			}
		}
		case false:
		//find mcr
		err = ss.CalcFbd()
		ss.Mur = bfac * ss.Zxx * ss.Fbd
	}
	return
}

//ColChk800 checks a section for axial compression as per is800 (chap. 8, duggal)
func (ss *StlSec) ColChk800()(err error){
	err = ss.CalcPu()
	if err != nil{
		return
	}
	if ss.Area * ss.Fcd < ss.Pu{
		err = fmt.Errorf("section unsafe - perm. axial load %.2f kN vs actual %.2f kN", ss.Area * ss.Fcd/1e3, ss.Pu/1e3)
	}
	return
}

//BmChk800 checks a section in bending as per is800 (b.Vu, b.Mu, b.Dmax)
func (ss *StlSec) BmChk800()(err error){
	//set basic deflection ratio
	switch ss.Ax{
		case 0:
		//major (x-axis) bending	
		if ss.Vu == 0.0 || ss.Mu == 0.0 || ss.Dmax == 0.0{
			err = fmt.Errorf("invalid ultimate design values- vu %f mu %f dmax %f",ss.Vu, ss.Mu, ss.Dmax)
			return
		}
		//check for shear
		err = ss.ShrChk800()
		if err != nil{
			return
		}
		//check for bending
		err = ss.CalcMur()
		if err != nil{
			return
		}
		if ss.Mu > ss.Mur{
			err = fmt.Errorf("section unsafe in bending - actual %f nmm vs permitted %f nmm",ss.Mu, ss.Mur)
			return
		}
		//check for deflection
		if ss.Defrat == 0.0{
			ss.Defrat = 300.0
		}
		//CHANGE THIS FACTOR LATER
		if ss.Lspan/ss.Defrat <= ss.Dmax/1.5{
			err = fmt.Errorf("section unsafe in (service) deflection perm %f mm vs actual %f mm",ss.Lspan/ss.Defrat,ss.Dmax)
			return
		}
		//check for web buckling
		err = ss.CalcFwb()
		if err != nil{
			return
		}
		switch{
			case ss.Fwb < ss.Vu:
			err = fmt.Errorf("section unsafe in web buckling fwb %f vs vu %f",ss.Fwb, ss.Vu)
			case ss.Fwc < ss.Vu:
			err = fmt.Errorf("section unsafe in web bearing/crippling fwc %f vs vu %f",ss.Fwc, ss.Vu)
		
		}
		case -2:
		//(simple) purlin, biaxial bending
		if ss.Mux == 0.0 || ss.Muy == 0.0{
			err = fmt.Errorf("invalid ultimate design values- mux %f muy %f", ss.Mux, ss.Muy)
			return
		}
		//fmt.Println(ColorGreen, ss.Sstr, ColorReset)
		//check for shear - nope. 
		ss.Mu = ss.Mux
		//check for bending
		err = ss.CalcMur()
		if err != nil{
			//fmt.Println(ColorRed,"ERRORE",err,ColorReset)
			return
		}
		//check for x- axis beding
		ss.Mdx = ss.Mur
		if ss.Mux > ss.Mdx{
			err = fmt.Errorf("section unsafe in bending about x axis - actual %f nmm vs permitted %f nmm",ss.Mux, ss.Mdx)
			
			//fmt.Println(ColorRed,"ERRROE-",err,ColorReset)
			return
		}
		//get mdy
		ss.Mdy = ss.Zyy * ss.Fy/ss.Ymo
		if ss.Mdy > 1.5 * ss.Syy * ss.Fy/ss.Ymo{
			ss.Mdy = 1.5 * ss.Syy * ss.Fy/ss.Ymo
		}
		if ss.Muy > ss.Mdy{
			err = fmt.Errorf("section unsafe in bending about y axis - actual %f nmm vs permitted %f nmm",ss.Muy, ss.Mdy)
			
			//fmt.Println(ColorRed,"ERRROE-",err,ColorReset)
			return
		}
		if ss.Mux/ss.Mdx + ss.Muy/ss.Mdy > 1.0{
			err = fmt.Errorf("section unsafe (local capcity check) %f > 1.0",ss.Mux/ss.Mdx + ss.Muy/ss.Mdy)
			//fmt.Println(ColorRed,"ERRROE-",err,ColorReset)
		}
		//fmt.Println("mdx, mdy, local cap-",ss.Mdx/1e6, ss.Mdy/1e6, ss.Mux/ss.Mdx + ss.Muy/ss.Mdy)
		if ss.Dmax > ss.Lspan/180.0{
			err = fmt.Errorf("section unsafe in deflection perm. %f vs actual %f",ss.Lspan/180.0, ss.Dmax)
			//fmt.Println(ColorRed,"ERRROE-",err,ColorReset)
		}
	}
	return
}

//BmColChk800 checks a section for axial load+bending as per is800(duggal chap. 10)
func (ss *StlSec) BmColChk800()(err error){
	ss.Printz()
	fmt.Println("mux, pu-",ss.Mux, ss.Pu,"nmm")
	err = ss.CalcPu()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("allowable compressive load, fcd-",ss.Fcd,"n/mm2")
	err = ss.CalcMur()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("ult moment of resistance mur",ss.Mur,"nmm")
	fmt.Println("section class-",ss.Ocf, ss.Icf)

	err = ss.CalcFbd()
	if err != nil{
		fmt.Println(err)
	}
	
	fmt.Println("allowable stress in bending",ss.Fbd,"n/mm2")

	return
}

//CalcFwb calculates the web buckling (and crippling) strength of a section
func (ss *StlSec) CalcFwb()(err error){
	switch ss.Sname{
		case "i","built-i","ub","uc":
		if ss.Lbr == 0.0{ss.Lbr = 100.0}
		klr := math.Sqrt(6.0) * ss.Hwrat
		nwb := math.Sqrt(ss.Fy * klr * klr/ss.Em)/math.Pi
		alp := 0.49
		owb := 0.5 * (1.0 + alp * (nwb - 0.22) + math.Pow(nwb,2))
		dnx := owb + math.Pow(owb*owb - nwb*nwb,0.5)
		ss.Fwb = (ss.Lbr + ss.H/2.0) * ss.Tw* (ss.Fy/ss.Ymo/dnx)
		ss.Fwc = (ss.Lbr + 2.5 * (ss.Tf + ss.R1)) * ss.Tw * ss.Fy/ss.Ymo
		//fmt.Println(ColorGreen,"fwb - ",ss.Fwb/1e3, "kn fwc-",ss.Fwc/1e3,"kn",ColorReset)
	}
	return
}

//ShrChk800 checks a (beam section) for low and high shear cases
func (ss *StlSec) ShrChk800()(err error){
	ss.ShrAr()
	ss.Vdx = ss.Avx * ss.Fy/math.Sqrt(3)/ss.Ymo
	ss.Vdy = ss.Avy * ss.Fy/math.Sqrt(3)/ss.Ymo
	switch ss.Ax{
		case 0:
		ss.Vd = ss.Vdx
		if ss.Vu > ss.Vd{
			err = fmt.Errorf("section unsafe in shear - design %f n vs actual %f n",ss.Vd,ss.Vu)
			return
		}
		if ss.Vu >= 0.6 * ss.Vd{
			ss.Scrit = true
		}
	}
	return
}

//ShrAr calcs the shear area about x and y (maj/min) axes
//(for now) as per is800
func (ss *StlSec) ShrAr(){
	switch ss.Code{
		case 1:	
		switch ss.Sname{
			case "i", "c":
			ss.Avx = ss.H * ss.Tw
			ss.Avy = 2.0 * ss.B * ss.Tf
			case "built-i":
			ss.Avx = (ss.H - 2.0 * ss.Tf) * ss.Tw
			ss.Avy = 2.0 * ss.B * ss.Tf
			case "l","ln":
		}
	}
	return
}

//CalcFbd calculates the allowable strength in bending 
func (ss *StlSec) CalcFbd()(err error){
	cfac := ss.Ocf; if ss.Icf > cfac{
		cfac = ss.Icf
	}
	bfac := 1.0
	switch cfac{
		case 1:
		//plastic
		case 2:
		//compact
		case 3:
		//semi compact
		bfac = ss.Sxx/ss.Zxx
	}
	switch ss.Sname{
		case "i", "built-i", "ub", "uc":
		if ss.Klx == 0.0{ss.Klx = 1.0}
		ss.Leffx = ss.Lspan * ss.Klx
		ss.Iww = 0.5 * 0.5 * ss.Iyy * (ss.H - ss.Tf) * (ss.H - ss.Tf)
		ss.Itt = ss.Sec.Prop.J
		//fmt.Println("st venants const-",ss.Itt,"warp const-",ss.Iww)
		p2 := math.Pow(math.Pi, 2)
		t1 := p2 * ss.Em * ss.Iyy/math.Pow(ss.Leffx,2)
		t2 := ss.Gm * ss.Itt
		t3 := p2 * ss.Em * ss.Iww/math.Pow(ss.Leffx,2)
		ss.Mcr = math.Sqrt(t1 * (t2 + t3)) 
		//fmt.Println(ColorRed,"mcr-",ss.Mcr/1e6,"knm",ColorReset)
		nlt := math.Sqrt(bfac * ss.Zxx * ss.Fy/ss.Mcr)
		//fmt.Println("non dimensional srat-",nlt)
		alt := 0.21
		if ss.Sname == "built-i"{
			alt = 0.49
		}
		ylt := 1.0 //reduction factor for ltb
		if nlt >= 0.4{
			olt := 0.5 * (1.0 + alt * (nlt - 0.2) + nlt * nlt)
			ylt = olt + math.Pow(olt * olt - nlt * nlt, 0.5)
			ylt = 1.0/ylt
			if ylt > 1.0{
				ylt = 1.0
			}
		}
		ss.Fbd = ylt * ss.Fy/ss.Ymo
	}
	return
}

//NetArea calcs the net and gross area of a section
func (ss *StlSec) NetArea()(err error){
	switch ss.Sname{
		case "l", "ln", "l2-ss", "l2-os", "ln2-ss", "ln2-os":
		//angle sections
		if len(ss.Dims) < 3{
			err = fmt.Errorf("invalid section dimensions %f",ss.Dims)
			return
		}
		bf := ss.Dims[0]; dw := ss.Dims[1]; tf := ss.Dims[2]
		cleg := dw
		cout := bf
		amul := 1.0
		bagomax := ss.Fu * ss.Ymo/(ss.Fy * ss. Yml)
		switch ss.Sname{
			case "l2-ss", "l2-os", "ln2-ss", "ln2-os":
				amul = 2.0
		}
		if ss.Cleg == 2{
			//short leg connected
			cleg = bf
			cout = dw
		}
		
		switch ss.Weld{
			case true:
			ss.Anc = (cleg - tf/2.0) * tf * amul
			ss.Ago = (cout - tf/2.0) * tf * amul
			ss.Sleg = cout
			ss.Bago = 1.4 - 0.076 * (cout/tf) * (ss.Fy/ss.Fu) * ss.Sleg/ss.Wg.L1
			if ss.Bago > bagomax{ss.Bago = bagomax}
			if ss.Bago < 0.7{ss.Bago = 0.7}
			ss.Avg = (ss.Wg.L1 + ss.Wg.L2)* ss.Wg.Tp * amul
			ss.Avn = (ss.Wg.L1 + ss.Wg.L2) * ss.Wg.Tp * amul
			ss.Atg = cleg * ss.Wg.Tp * amul
			ss.Atn = cleg * ss.Wg.Tp * amul
			case false:	
			if ss.Cleg == -1{
				//get equivalent plate
				ss.Bg.Bmem = bf + dw - tf
				ss.Bg.Tmem = tf
			}
			ss.Bg.Bmem = cleg - tf/2.0
			ss.Bg.Tmem = tf
			//ni, nj are given
			err = ss.Bg.BltNsa()
			if err != nil{
				return
			}
			ss.Anc = ss.Bg.Nsar
			if ss.Cleg != -1{
				//TODO TODO TODO TODO
				//handle this here and return
				ss.Anc = ss.Bg.Nsar * amul
			}
			ss.Sleg = cout + ss.Bg.Endd - tf/2.0
			ss.Bago = 1.4 - 0.076 * (cout/tf) * (ss.Fy/ss.Fu) * ss.Sleg/ss.Bg.Ljoint

			ss.Ago = (cout - tf/2.0)*tf * amul
			
			//block shear areas
			nj := float64(ss.Bg.Nj)
			ss.Avg = ((nj - 1.0) * (ss.Bg.Pitch) + ss.Bg.Endd) * tf * amul 
			ss.Avn = ((nj - 1.0) * ss.Bg.Pitch + ss.Bg.Endd - (nj - 0.5) * ss.Bg.D0)*tf * amul
			ss.Atg = (cleg - ss.Bg.Edged) * tf * amul
			ss.Atn = (cleg - ss.Bg.Edged - ss.Bg.D0/2.0) * tf * amul
		}		
	}
	return
}

//CalcPu calculates the maximum axial load on a (column) section
func (ss *StlSec) CalcPu()(err error){
	switch ss.Code{
		case 1:
		err = ss.BckCls800()
	}
	//ss.Printz()
	err = ss.CalcFcd()
	return
}

//CalcFcd calculates the allowable compressive stress for a steel (column) section
func (ss *StlSec) CalcFcd()(err error){
	//default pin ended
	if ss.Klx == 0{ss.Klx = 1.0}
	if ss.Kly == 0{ss.Kly = 1.0}
	//imperfection factor alpha
	alps := []float64{0.21, 0.34,0.49,0.76}
	alpx := alps[ss.Bxx-1]
	alpy := alps[ss.Byy-1]
	//what about double anglezzz
	switch{
		case (ss.Sname == "l" || ss.Sname == "ln") && ss.Cleg != -1:
		var k1, k2, k3 float64
		switch{
			case ss.Weld:
			k1 = 0.2
			k2 = 0.35
			k3 = 20.0
			case ss.Bg.Nb == 1:
			k1 = 0.75
			k2 = 0.35
			k3 = 20.0
			case ss.Bg.Nb > 1:
			k1 = 0.2
			k2 = 0.35
			k3 = 20
			default:
			err = fmt.Errorf("no bolt group specified - nbolts %v ni(rows) %v nj(cols) %v vec %v",ss.Bg.Nb, ss.Bg.Ni, ss.Bg.Nj,ss.Bg.Bvec)
		}
		rvv := ss.Sec.Prop.Rvv; if ss.Sec.Prop.Ruu < rvv{
			rvv = ss.Sec.Prop.Ruu
		}
		ysr := math.Sqrt(250.0/ss.Fy)
		nvv := ss.Lspan/rvv/ysr/math.Sqrt(math.Pow(math.Pi,2)*ss.Em/250.0)
		noo := (ss.Dims[0]+ss.Dims[1])/(2.0*ss.Dims[2])/ysr/math.Sqrt(math.Pow(math.Pi,2)*ss.Em/250.0)
		nee := math.Sqrt(k1 + k2 * math.Pow(nvv, 2) + k3 * math.Pow(noo, 2))
		ovv := 0.5 * (1.0 + alpx * (nee - 0.2) + math.Pow(nee, 2))
		dvv := ovv + math.Pow(ovv * ovv - nee * nee, 0.5)
		ss.Fcd = ss.Fy/ss.Ymo/dvv
		default:
		if ss.Leffx == 0.0{
			ss.Leffx = ss.Lspan * ss.Klx
		}
		if ss.Leffy == 0.0{
			ss.Leffy = ss.Lspan * ss.Kly
		}
		ss.Fccx = math.Pow(math.Pi, 2) * ss.Em/math.Pow(ss.Leffx/ss.Rxx,2)
		ss.Fccy = math.Pow(math.Pi, 2) * ss.Em/math.Pow(ss.Leffy/ss.Ryy,2)
		nxx := math.Sqrt(ss.Fy/ss.Fccx)
		nyy := math.Sqrt(ss.Fy/ss.Fccy)
		oxx := 0.5 * (1.0 + alpx * (nxx - 0.22) + math.Pow(nxx,2))
		oyy := 0.5 * (1.0 + alpy * (nyy - 0.22) + math.Pow(nyy,2))
		dnx := oxx + math.Pow(oxx*oxx - nxx*nxx,0.5)
		dny := oyy + math.Pow(oyy*oyy - nyy*nyy,0.5)
		if dnx <= 0 || dny <= 0{
			err = fmt.Errorf("error in calculating fcd")
			return
		}
		ss.Fcdx = ss.Fy/ss.Ymo/dnx
		ss.Fcdy = ss.Fy/ss.Ymo/dny
		ss.Fcd = ss.Fcdy; if ss.Fcdx < ss.Fcd{ss.Fcd = ss.Fcdx}
	}
	//fmt.Println(ColorRed,"perm.load-",ss.Area*ss.Fcd/1e3,"kn",ColorReset)	
	return
}

//bs449 funcs
//ColChk449 checks a column section as per bs449, 1969 (mosley section 6.1)
func (ss *StlSec) ColChk449()(err error){
	var pa, px, py float64
	lx := ss.Lx * ss.Tx ; ly := ss.Ly * ss.Ty
	//convert from iscode
	if ss.Code == 1 {
		switch ss.Grd{
			case 410.0:
			ss.Grd = 43.0
			case 500.0:
			ss.Grd = 50.0
		}
	}
	pvec := PqCol[int(ss.Grd)]
	var vx, vy, mx, my float64
	if ss.H2 > 0.0{
		vx = ss.Vbdx * ss.H2/(ss.H1+ss.H2)
		vy = ss.Vbdy * ss.H2/(ss.H1+ss.H2)
	}
	mx = ss.Mux; my = ss.Muy; ss.Dtyp = 1

	if mx + my == 0.0 {
		ss.Dtyp = 0 //member with framing beams
		mx = vx * (100.0 + ss.H)
		my = vy * (100.0 + ss.Tw/2.0)
	}
	fa := ss.Pu/ss.Area
	fx := mx/ss.Sxx
	fy := my/ss.Syy
	fp := fa/pvec[0] + fx/pvec[1] + fy/pvec[2]
	
	if fp > ss.Pfac{
		err = fmt.Errorf("fp %f > permimssible factor %f",fp, ss.Pfac)
		return
	}
	var s1 float64
	sx := lx/ss.Rxx
	sy := ly/ss.Ryy
	s1 = sx; if sy > s1 {s1 = sy}
	if s1 > 180.0 {
		err = fmt.Errorf("invalid slenderness ratio %f > 180",s1)
		return
	}
	//permissible axial stress pa
	var y0, q4, q5 float64
	c0 := math.Pow(math.Pi,2) * EStl/math.Pow(s1,2)
	n0 := 0.3 * math.Pow(s1/100.0,2)
	switch{
		case ss.Grd == 43:
		y0 = 250.0; q4 = 155.0; q5 = 143.0
		case ss.Grd == 50:
		y0 = 350.0; q4 = 215.0; q5 = 200.0
		case ss.Grd == 55:
		y0 = 430.0; q4 = 265.0; q5 = 245.0
	}
	if ss.Grd == 50 && ss.Tw >= 40.0{
		//CHEEECK THIS
		y0 = 325.0; q4 = 200.0; q5 = 185.0
	}
	a0 := (y0 + c0 * (n0 + 1.0))/2.0
	pa = (a0 - math.Sqrt((math.Pow(a0,2) - y0 * c0)))/1.7
	if s1 <=30 {
		pa = q4 - (q4 - q5) * s1/30.0
	}
	
	//permissible stress bending(x) px
	var dtrat float64
	if dtrat = ss.H/ss.Tf; dtrat < 5.0 {dtrat = 5.0}	
	//log.Println("checking px->",s1,dtrat)
	// if ss.Yeolde{
	// 	px = PbcYeolde(s1, dtrat)
	// } else {
	// 	px, _ = PbcLerp(ss.Sname, ss.Grd, s1, dtrat)
	// }
	px, err = PbcLerp(ss.Sname, int(ss.Grd), s1, dtrat)
	if err != nil{
		return
	}
	//permissible stress in bending (y) py
	py = pvec[2]
	fp = fa/pa + fx/px + fy/py

	//log.Println("***")
	if fp <= ss.Pfac{
		//wt := ss.Wt
		log.Println("section found->",ss.Sstr)
		log.Println("base fp->",fp)
		log.Println("srats->",sx,sy,s1)
		log.Println("paxial->",pa)
		log.Println("px->",px)
		log.Println("fp->",fp)
		// log.Println("depth, web thickness->",df.Elem(i,3), df.Elem(i,6))
		// log.Println("area, zx, zy->",df.Elem(i,23),df.Elem(i,15), df.Elem(i,16))
		// log.Println("rx, ry->",df.Elem(i,13), df.Elem(i,14))
		log.Println("mx, my, s1, dtrat->", mx, my, s1, dtrat)
		log.Println("fa, pa, px, py, fp ->",fa, pa, px, py, fp)
		log.Println("***")
		
	}
	
	return
}

//BmChk449 checks a beam section as per bs449, 1969 (mosley section 6.2)
func (ss *StlSec) BmChk449() (err error){
	// var b1, b2, fm, fs, fc, fb, y0, q4, q5 float64
	// log.Println("checking section->",df.Elem(i,1))
	// pvec := PqBm[ss.Grd]
	// qx, ps, pc := pvec[0], pvec[1], pvec[2]
	// log.Println("section modulus->",ss.Sxx)
	// if ss.Mu*1000.0/ss.Sxx > qx {
	// 	return
	// }
	// dtrat := df.Elem(i, 3).Float() / df.Elem(i, 6).Float()
	// if dtrat < 5.0 {
	// 	dtrat = 5.0
	// }
	// if ss.Spam{log.Println("dtrat->", dtrat)}

	// sdrat := ss.Ly * 100.0 / df.Elem(i, 14).Float()
	
	// if ss.Spam{log.Println("sdrat->", dtrat)}
	// var px float64
	// if ss.Yeolde{
	// 	px = PbcYeolde(sdrat, dtrat)
	// } else {
	// 	px = PbcLerp(ss.Sname, ss.Grd, sdrat, dtrat)
	// }
	// fm = math.Abs(ss.Mu) * 1e3 / df.Elem(i, 15).Float()
	// if fm/px > 1.0 {
	// 	return
	// }
	// if ss.Spam{log.Println("section->", df.Elem(i, 1), ColorBlue, "bending o.k", ColorReset)}
	// //check for shear stress
	// fs = math.Abs(ss.Vu) * 1e3 / df.Elem(i, 3).Float() / df.Elem(i, 5).Float()
	// if fs/ps > 1.0 {
	// 	return
	// }

	// if ss.Spam{log.Println("section->", df.Elem(i, 1), kass.ColorBlue, "shear o.k", kass.ColorReset)}
	// //check for deflection
	// //dmax := ss.Dmax
	// defrat := math.Round(ss.Lspan * 1000. / ss.Dmax)
	// if ss.Spam{log.Println("deflection ->", ss.Dmax, "mm vs perm. ->", ss.Lspan*1000./360., "perm. ratio - 360.0, actual ->", defrat)}
	// if ss.Dmax > ss.Lspan*1000./360.{
	// 	return
	// }
	
	// //log.Println("section->",df.Elem(i,1),kass.ColorBlue,"deflection o.k",kass.ColorReset)
	// //check for web crushing stress
	// b1 = ss.Lbr + (ss.Tbr+0.5*(df.Elem(i, 3).Float()-df.Elem(i, 8).Float()))*math.Sqrt(3.0)
	// fc = 1e3 * math.Abs(ss.Vu) / df.Elem(i, 5).Float() / b1
	// if ss.Spam{log.Println("web crushing-> perm ",pc," vs ",fc,"actual (n/mm)")}
	// if fc/pc > 1.0 {
	// 	return
	// }
	// //check for web buckling stress
	// b2 = ss.Lbr + ss.Tbr + df.Elem(i, 3).Float()/2.0
	// fb = 1e3 * math.Abs(ss.Vu) / df.Elem(i, 5).Float() / b2
	// sweb := math.Sqrt(3.0) * df.Elem(i, 8).Float() / df.Elem(i, 5).Float()
	// c0 := math.Pow(math.Pi, 2) * EStl / math.Pow(sweb, 2)
	// n0 := 0.3 * math.Pow(sweb/100.0, 2)
	// switch {
	// case ss.Grd == 43:
	// 	y0 = 250.0
	// 	q4 = 155.0
	// 	q5 = 143.0
	// case ss.Grd == 50:
	// 	y0 = 350.0
	// 	q4 = 215.0
	// 	q5 = 200.0
	// case ss.Grd == 55:
	// 	y0 = 430.0
	// 	q4 = 265.0
	// 	q5 = 245.0
	// }
	// a0 := (y0 + c0*(n0+1.0)) / 2.0
	// pb := (a0 - math.Sqrt((math.Pow(a0, 2) - y0*c0))) / 1.7
	// if sweb <= 30 {
	// 	fb = q4 - (q4-q5)*sweb/30.0
	// }
	// if ss.Spam{log.Println("web buckling-> perm ",pb," vs ",fb,"actual (n/mm)")}
	// if fb/pb > 1.0{
	// 	return
	// }
	// chk = true
	// //vals = append(maxs, []float64{}...)
	// wt := df.Elem(i, 2).Float() 
	// vals = []float64{ss.Vu, ss.Mu, ss.Dmax, ss.Lspan*1000./360.,px, fm, ps, fs, pc, fc, pb, fb, wt}
	// // 0 Vu, 1 ss.Mu, 2 ss.Dmax,3 ss.Lspan*1000./360.,4 px, 5 fm, 6 ps, 7 fs, 8 pc, 9 fc, 10 pb, 11 fb, 12 wt

	return
}

//connection design funcs

func (ss *StlSec) ConDz()(err error){
	switch ss.Frmstr{
		case "2dt":
		switch ss.Weld{
			case false:
			//blt dz
			switch ss.Sname{
				case "l","ln","l2-ss","l2-os","ln2-ss","ln2-os":
				case "tube":
				//TODO
			}
		}
	}
	return
}



// //purlin design funcs
// func StlPrlnDz(pdl, pll, pwl, theta, pspc, tspc float64, sname string, nsecs int)(sss []StlSec, err error){
// 	// pdl := 130.0
// 	// pll := 0.0
// 	// pwl := 2000.0
// 	// theta := 30.0
// 	// pspc := 1500.0
// 	// tspc := 6000.0
// 	// sname := "i"
// 	p := (pdl + pll) * pspc/1000.0
// 	ndx := StlSdxLims[sname]
// 	for idx := ndx; idx > 0; idx--{
// 		if len(sss) == nsecs{
// 			break
// 		}
// 		ss, err := GetStlSec(sname, idx, 1)
// 		if err != nil{
// 			log.Println("ERRORE,errore",err)
// 		}
// 		p += ss.Wt
// 		pn := p * math.Cos(theta * math.Pi/180.0)
// 		pp := p * math.Sin(theta * math.Pi/180.0)
// 		pn  += pwl * pspc/1000.0
		
// 		log.Println("unfactored pnormal P", pn, "pparallel H",pp)
// 		mux := 1.5 * pn * tspc * tspc/10.0/1000.0
// 		muy := 1.5 * pp * tspc * tspc/10.0/1000.0
// 		//t.Log("mux, muy",mux, muy)
// 		ss.Ax = -2
// 		ss.Lspan = tspc
// 		ss.Mux = mux
// 		ss.Muy = muy
// 		ss.Lsb = true
// 		ss.Dmax = 5.0 * pn * math.Pow(tspc, 4)/ss.Em/ss.Ixx/384.0/1000.0
// 		log.Println("ss.Dmax-",ss.Dmax)
// 		ss.Vux = 0.6 * pn * tspc
// 		ss.Vuy = 0.6 * pp * tspc
// 		err = ss.BmChk800()
// 		if err == nil{
// 			sss = append(sss, ss)
// 		}
// 	}
// 	// p += 100.0
// 	return
	
// }


//GetStlDf gets stl df based on NAME
func GetStlDf(sname string) (df dataframe.DataFrame, err error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	var sheet string
	switch sname{
		case "i","built-i","plate-i":
		sheet = filepath.Join(basepath, "../data/steel/isteel", "I.csv")
		case "c":
		sheet = filepath.Join(basepath, "../data/steel/isteel", "ISMC.csv")
		case "l","l2-ss","l2os":
		sheet = filepath.Join(basepath, "../data/steel/isteel", "ISA-eq.csv")
		case "ln","ln2-ss","ln2-os":
		sheet = filepath.Join(basepath, "../data/steel/isteel", "ISA-ueq.csv")
		case "tube", "pipe":
		sheet = filepath.Join(basepath, "../data/steel/isteel", "ISNB.csv")
		case "ub":
		sheet = filepath.Join(basepath, "../data/steel/bsteel", "UB.csv")
		case "uc":
		sheet = filepath.Join(basepath, "../data/steel/bsteel", "UC.csv")
		default:
		err = fmt.Errorf("%s-csv source not found",sname)
		return
	}
	csvfile, err := os.Open(sheet)
	if err != nil {
		return dataframe.DataFrame{}, err
	}
	df = dataframe.ReadCSV(csvfile)
	return df, err
}

//GetStlSec returns a StlSec given sectyp and sdx
func GetStlSec(sname string, sdx, code int, params...float64) (ss StlSec, err error) {
	if _, ok := StlStyps[sname]; !ok{
		err = fmt.Errorf("invalid section name - %s",sname)
		return
	}
	
	grd := 410.0
	fy := 250.0
	fu := 410.0
	em := 200000.0
	ymo := 1.1
	yml := 1.25
	if code == 2{
		grd = 43.0
		em = 210000.0
	}
	if sdx == -1{
		switch sname{
			case "built-i":
			if len(params) < 4{
				err = fmt.Errorf("invalid params for built-i section %.f",params)
			}
			bf := params[0]
			d := params[1]
			tf := params[2]
			tw := params[3]
			ss.Dims = []float64{bf, d, tf, tw}
			ss.Bstyp = 12
			ss.Sname = sname
			err = ss.SecGen()
			if err != nil{
				return
			}
			ss.Sdx = sdx
			ss.Grd = grd
			ss.Fy = fy
			ss.Fu = fu
			ss.Ymo = ymo
			ss.Yml = yml
			ss.Code = code
			ss.Em = em
			ss.Gm = em/2.0/(1.3)
			return
		}
	}
	var df dataframe.DataFrame
	df, err = GetStlDf(sname)
	if err != nil{
		return
	}
	if sdx > df.Nrow() - 1{
		err = errors.New("invalid section index")
		return
	}
	//change grade and code later
	ss = StlSec{
		Sname:  sname,
		Sectyp: StlStyps[sname],
		Sstr:   df.Elem(sdx, 1).String(),
		Sdx:    sdx,
		Grd:    grd,
		Fy:     fy,
		Fu:     fu,
		Ymo:    ymo,
		Yml:    yml,
		Code:   code,
		Em:     em,
		Gm:     em/2.0/(1.3),
		Wt:     df.Elem(sdx, 2).Float()*9.81,
		H:      df.Elem(sdx, 3).Float(),
		B:      df.Elem(sdx, 4).Float(),
		Tw:     df.Elem(sdx, 5).Float(),
		Tf:     df.Elem(sdx, 6).Float(),
		R1:     df.Elem(sdx, 7).Float(),
		R2:     df.Elem(sdx, 8).Float(),
		Ixx:    df.Elem(sdx, 11).Float()*10000.0,
		Iyy:    df.Elem(sdx, 12).Float()*10000.0,
		Rxx:    df.Elem(sdx, 13).Float()*10.0,
		Ryy:    df.Elem(sdx, 14).Float()*10.0,
		Zxx:    df.Elem(sdx, 15).Float()*1000.0,
		Zyy:    df.Elem(sdx, 16).Float()*1000.0,
		Sxx:    df.Elem(sdx, 17).Float()*1000.0,
		Syy:    df.Elem(sdx, 18).Float()*1000.0,
		Area:   df.Elem(sdx, 23).Float()*100.0,
	}
	ss.Eps = math.Sqrt(250.0/ss.Fy)
	//ss.Sfmat = 1.1
	switch sname{
		case "i", "ub", "uc":
		ss.Bstyp = 12
		ss.Dims = []float64{ss.B, ss.H, ss.Tf, ss.Tw}
		case "c":
		ss.Bstyp = 13
		ss.Dims = []float64{ss.B, ss.H, ss.Tf, ss.Tw}
		case "rect", "flat", "isf":
		ss.Dims = []float64{ss.B, ss.H}
		case "tube":
		D := ss.H
		d := ss.H - 2.0 * ss.Tw
		ss.Bstyp = 5
		ss.Dims = []float64{D, d}
		case "l","ln","l2-ss","l2-os","ln2-ss","ln2-os":
		bf := ss.B// df.Elem(sdx, 4).Float()
		d := ss.H// df.Elem(sdx, 3).Float()
		bw := ss.Tw// df.Elem(sdx, 5).Float()
		tf := ss.Tf// df.Elem(sdx, 6).Float()
		tp := tf
		if len(params) > 0{
			tp = params[0]
		}
		ss.Dims = []float64{bf, d, bw, tf, tp}
		switch sname{
			case "l2-ss","ln2-ss":
			ss.Bstyp = 27
			ss.Sstr = fmt.Sprintf("2x%s-ss",ss.Sstr)
			case "l2-os","ln2-os":
			ss.Bstyp = 28
			ss.Sstr = fmt.Sprintf("2x%s-os",ss.Sstr)
			default:
			ss.Bstyp = 7
		}
		// case "built-i":
		// b := df.Elem(sdx, 4).Float()
		// h := df.Elem(sdx, 3).Float()
		// tf := df.Elem(sdx, 6).Float()
		// tw := df.Elem(sdx, 5).Float()
		// ss.Dims = []float64{b, h, tf, tw}
		// ss.Bstyp = 12
		//OR SHOULD THIS BE PARAMS?
		case "plate-i":
		if len(params) < 2{
			err = fmt.Errorf("plate dimensions not specified-%f",params)
			return
		}
		b := df.Elem(sdx, 4).Float()
		h := df.Elem(sdx, 3).Float()
		tf := df.Elem(sdx, 6).Float()
		tw := df.Elem(sdx, 5).Float()
		B := params[0]
		D := params[1]
		ss.Dims = []float64{b, h, tf, tw, B, D}
		ss.Bstyp = 29
		ss.Sstr = fmt.Sprintf("%s+2x(%.fx%.f)",ss.Sstr,B,D)
	}	
	err = ss.SecGen()
	return

}

//GetStlCp returns the cross section property slice given frm type, sectype and sheet index
func GetStlCp(sname string, frmtyp, sdx, ax int) (cp []float64, err error) {
	var df dataframe.DataFrame
	df, err = GetStlDf(sname)
	if err != nil {
		return
	}
	if sdx > df.Nrow() {
		err = errors.New("invalid section index")
		return
	}
	switch frmtyp {
	case 1:
		//1d b - iz
		switch ax {
		case 1:
			cp = []float64{df.Elem(sdx, 11).Float()}
		case 2:
			cp = []float64{df.Elem(sdx, 12).Float()}
		}
	case 2:
		//2d t - a, iz
		switch ax {
		case 1:
			//major axis of bending - x
			cp = []float64{df.Elem(sdx, 23).Float(), df.Elem(sdx, 11).Float()}
		case 2:
			//y axis
			cp = []float64{df.Elem(sdx, 23).Float(), df.Elem(sdx, 12).Float()}
		}
	case 3:
		//2d f - a, iz
		switch ax {
		case 1:
			//major axis of bending - x
			cp = []float64{df.Elem(sdx, 23).Float(), df.Elem(sdx, 11).Float()}
		case 2:
			//y axis
			cp = []float64{df.Elem(sdx, 23).Float(), df.Elem(sdx, 12).Float()}
		}
	case 4:
		//3d t - a
		switch ax {
		case 1:
			//major axis of bending - x
			cp = []float64{df.Elem(sdx, 23).Float(), df.Elem(sdx, 11).Float()}
		case 2:
			//y axis
			cp = []float64{df.Elem(sdx, 23).Float(), df.Elem(sdx, 12).Float()}
		}
	case 5:
	//3d g - a, etc
	case 6:
		//3d f
	}
	return
}

//ReadSecDims reads in base dims of a section
func ReadSecDims(sname string, sdx, ax int) (dims []float64, bstyp int, err error) {
	var df dataframe.DataFrame
	df, err = GetStlDf(sname)
	if err != nil {
		return
	}
	if sdx > df.Nrow() {
		err = errors.New("invalid section index")
		return
	}
	//x dims
	switch sname{
	case "i","ub","uc":
		//i/ieq
		b := df.Elem(sdx, 4).Float()
		h := df.Elem(sdx, 3).Float()
		tf := df.Elem(sdx, 5).Float()
		tw := df.Elem(sdx, 4).Float()
		dims = []float64{b, h, tf, tw}
		bstyp = 12
	}
	return
}

//GetHaunchDims returns dims at x from jb of a non-prismatic member
func GetHaunchDims(sname string, sdx, ax int, dy float64) (dims []float64, err error) {
	var df dataframe.DataFrame
	df, err = GetStlDf(sname)
	if err != nil {
		return
	}
	if sdx > df.Nrow() {
		err = errors.New("invalid section index")
		return
	}
	//x dims
	switch sname{
	case "built-i","i","ub","uc":
		//i/ieq
		//haunch2f/3f is the same as far as dims go?
		b := df.Elem(sdx, 4).Float()
		h := df.Elem(sdx, 3).Float()
		tf := df.Elem(sdx, 5).Float()
		tw := df.Elem(sdx, 4).Float()
		dims = []float64{b, h, tf, tw, dy}
	}
	return
}

//StlSecInit initalizes a (standard) steel section struct
//WHAAT WHAT IS IT GOOD FOR
//if kl has to be calc'ed this can be done separately
// func StlSecInit(sname string, sdx, code int, lspan, kl float64)(ss StlSec, err error){
// 	ss, err = GetStlSec(sname, sdx, code)
// 	ss.Lspan = lspan
// 	ss.Klx = kl
// 	// switch ss.Code{
// 	// 	case 1:
// 	// 	//is 800-2007
// 	// 	err = ss.BckClass800()
// 	// 	case 2:
// 	// 	//bs 5950
// 	// }
// 	return
// }
