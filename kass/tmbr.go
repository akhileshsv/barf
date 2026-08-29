package barf

import (
	"sort"
	"errors"
)

var (
	TmbrGrp = []string{"a","b","c"}
	TmbrDims0 = [][]float64{{3.0*25.4},{4.0*25.4},{5.0*25.4},{6.0*25.4},{7.0*25.4},{8.0*25.4},{9.0*25.4},{10.0*25.4},{11.0*25.4},{12.0*25.4}}
	
	RectBs = []float64{25,30,35,40,50,60,80,100}
	RectDs = []float64{40,50,60,80,100,120,140,160,180,200}
	//good ol' IDLE to the rescue
	TmbrDims1 = [][]float64{{25, 40}, {25, 50}, {25, 60}, {25, 80}, {25, 100}, {25, 120}, {25, 140}, {25, 160}, {25, 180}, {25, 200}, {30, 40}, {30, 50}, {30, 60}, {30, 80}, {30, 100}, {30, 120}, {30, 140}, {30, 160}, {30, 180}, {30, 200}, {35, 40}, {35, 50}, {35, 60}, {35, 80}, {35, 100}, {35, 120}, {35, 140}, {35, 160}, {35, 180}, {35, 200}, {40, 40}, {40, 50}, {40, 60}, {40, 80}, {40, 100}, {40, 120}, {40, 140}, {40, 160}, {40, 180}, {40, 200}, {50, 40}, {50, 50}, {50, 60}, {50, 80}, {50, 100}, {50, 120}, {50, 140}, {50, 160}, {50, 180}, {50, 200}, {60, 40}, {60, 50}, {60, 60}, {60, 80}, {60, 100}, {60, 120}, {60, 140}, {60, 160}, {60, 180}, {60, 200}, {80, 40}, {80, 50}, {80, 60}, {80, 80}, {80, 100}, {80, 120}, {80, 140}, {80, 160}, {80, 180}, {80, 200}, {100, 40}, {100, 50}, {100, 60}, {100, 80}, {100, 100}, {100, 120}, {100, 140}, {100, 160}, {100, 180}, {100, 200}}
	PlyDs = []float64{6,9,12,16,19,25}
	DeckDs = []float64{25,30,35,40,45,50}
	TmbrSort = false
)

//Wdprp stores wood properties 
type Wdprp struct{
	Em, Pg    float64 //em, density
	G, U      float64 //shear modulus, poisson's ratio
	Ft, Fc    float64 //allowable stress (tensile, compressive) parallel to grain
	Ftp, Fcp  float64 //perpendicular to grain
	Fv, Fvp   float64 //shear parallel, perpendicular to grain
	Fcb       float64 //in bending
	Mc        float64 
	Mp        float64
	Ei        float64
	Wdl       float64
	Grp       int
	Mat       int
	Grade     int
	Loc       int     //1 - outside, 2 - wet
}

//Init initializes wood properties given a group
//as seen in is.883
func (p *Wdprp) Init(grp int) (error){
	//modification factor for grade of timber
	fg := 1.0
	//fg = 0.84
	if grp > 3{
		grp = 3
		p.Grade = grp - 3 
	}
	switch p.Grade{
	//is 883
		case -1:
		//"select" grade timber (wtf)
		fg = 1.16
		case 0:
		//default
		case 1:
		//grade 2 timber
		fg = 0.84
		case 2:
		//low durability timber
		fg = 0.84 * 0.8
	}
	switch p.Loc{
	//see table 3 note 2 is 883
		case 1:
		//outside
		fg = fg * 5.0/6.0
		case 2:
		//wet
		fg = fg * 2.0/3.0
	}
	switch grp{
		case 1:
		//group a timber
		p.Ft = 18.2 * fg
		p.Fcb = 18.2 * fg
		p.Fvp = 1.05 * fg
		p.Fv = 1.5 * fg
		p.Fc = 11.7 * fg
		p.Fcp = 4.0 * fg
		p.Em = 12.6e3
		p.Pg = 0.85
		case 2:
		//group b timber
		p.Ft = 12.3 * fg
		p.Fcb = 12.3 * fg
		p.Fvp = 0.64 * fg
		p.Fv = 0.91 * fg
		p.Fc = 7.8 * fg
		p.Fcp = 2.5 * fg
		p.Em = 9.8e3 
		p.Pg = 0.75
		case 3:
		fg = 0.84 * 0.8
		//group c timber
		p.Ft = 8.4 * fg
		p.Fcb = 8.4 * fg
		p.Fvp = 0.49 * fg
		p.Fv = 0.70 * fg
		p.Fc = 4.9 * fg
		p.Fcp = 1.1 * fg
		p.Em = 5.6e3 
		p.Pg = 0.65
		default:
		return errors.New("invalid timber group")
	}
	p.Grp = grp
	return nil
}

//GenTmbrDims1 generates timber dimensions from rectBs and rectDs
//from standard sizes in is.883
//this is really not needed now
func GenTmbrDims1() {
	//save
	//if len(TmbrDims1) > 3{return}
	// for i := range RectBs{
	// 	for j := range RectDs{
	// 		TmbrDims1 = append(TmbrDims1,[]float64{RectBs[i],RectDs[j]})
	// 	}
	// }
	if !TmbrSort{
		sort.Slice(TmbrDims1, func(i, j int) bool {
			return TmbrDims1[i][0]*TmbrDims1[i][1] < TmbrDims1[j][0]*TmbrDims1[j][1]
		})
		TmbrSort = true
	}
	return
}
