package barf

import (
	"sort"
	"errors"
)

var (
	tmbrGrp = []string{"a","b","c"}
	TmbrDims0 = [][]float64{{3.0*25.4},{4.0*25.4},{5.0*25.4},{6.0*25.4},{7.0*25.4},{8.0*25.4},{9.0*25.4},{10.0*25.4},{11.0*25.4},{12.0*25.4}}
	rectBs = []float64{25,30,40,50,60,75,100,120,140,160,180,200}
	rectDs = []float64{40,50,60,75,100,120,140,160,180,200}
	TmbrDims1 = [][]float64{}
	plyDs = []float64{6,9,12,16,19,25}
	deckDs = []float64{25,30,35,40,45,50}
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
}

//Init initializes wood properties given a group
//as seen in is.883
func (p *Wdprp) Init(grp int) (error){
	//modification factor for grade of timber
	fg := 1.0
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
func GenTmbrDims1() {
	for i := range rectBs{
		for j := range rectDs{
			TmbrDims1 = append(TmbrDims1,[]float64{rectBs[i],rectDs[j]})
		}
	}
	
	sort.Slice(TmbrDims1, func(i, j int) bool {
		return TmbrDims1[i][0]*TmbrDims1[i][1] < TmbrDims1[j][0]*TmbrDims1[j][1]
	})
	return
}
