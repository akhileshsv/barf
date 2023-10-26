package barf

import (
	"os"
	//"log"
	//"fmt"
	"runtime"
	"path/filepath"
	"encoding/csv"
	"errors"
	"strings"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

var (
	//1 2 3 4 5 6 7 8
	StlStyps = []string{"angle","box","channel","i","pipe","tee","ub","uc"}
)

//StlDfBs takes in a section type index and returns the bs code sheet as a dataframe 
func StlDfBs(sectyp int) (dataframe.DataFrame, error){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	var sheet string
	switch sectyp{
		case 7:
		//ub sec
		sheet = filepath.Join(basepath,"../data/steel/bsteel","UB.csv")
		case 8:
		//uc sec
		sheet = filepath.Join(basepath,"../data/steel/bsteel","UC.csv")
		default:
		return dataframe.DataFrame{}, errors.New("invalid section type")
	}
	csvfile, err := os.Open(sheet)
	if err != nil {
		return dataframe.DataFrame{}, err
	}
	df := dataframe.ReadCSV(csvfile)
	return df, err
}

//StlDfIs takes in a section type index and returns the is code sheet as a dataframe 
func StlDfIs(sectyp int) (dataframe.DataFrame, error){
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	tabledir := filepath.Join(basepath,"../data/steel")
	f, _ := os.Open(filepath.Join(tabledir,"isteel/sefistl.csv"))
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		//fmt.Println("ERRORE,errore->unable to open csv file")
		return dataframe.DataFrame{}, err
	}
	if sectyp > 6 || sectyp < 0{
		return dataframe.DataFrame{}, errors.New("invalid section type")
	}
	sectype := StlStyps[sectyp-1]
	df := dataframe.LoadRecords(records)
	fil := df.Filter(
		dataframe.F{Colname:"Shape",
			Comparator: series.Eq,
			Comparando:strings.Title(sectype)},
	)
	return fil, nil
}

//GetStlDf gets stl df based on code and styp
func GetStlDf(sectyp int) (df dataframe.DataFrame, err error){
	switch {
	case sectyp < 7:
		//sefi steel sections
		df, err = StlDfIs(sectyp) 
	default:
		//bs steel sections
		df, err = StlDfBs(sectyp)
	}
	return df, err
}

//GetStlCp returns the cross section property slice given frm type, sectype and sheet index
func GetStlCp(frmtyp, sectyp, sdx, ax int) (cp []float64, err error){
	var df dataframe.DataFrame
	switch {
	case sectyp < 7:
		//sefi steel sections
		//area ixx iyy rxx ryy 
		df, err = StlDfIs(sectyp)	
		if err != nil{
			return 
		}
		if sdx > df.Nrow(){
			err = errors.New("invalid section index")
			return
		}
		//log.Println(df)
		switch frmtyp{
			case 1:
			//1d b - iz
			switch ax{
				case 1:
				cp = []float64{df.Elem(sdx,12).Float()}
				case 2:
				cp = []float64{df.Elem(sdx,13).Float()}
			}
			case 2:
			//2d t - a, iz
			switch ax{
				case 1:
				//major axis of bending - x
				cp = []float64{df.Elem(sdx,8).Float(),df.Elem(sdx,12).Float()}
				case 2:
				//y axis
				cp = []float64{df.Elem(sdx,8).Float(),df.Elem(sdx,13).Float()}
			}
			case 3:
			//2d f - a, iz
			switch ax{
				case 1:
				//major axis of bending - x
				cp = []float64{df.Elem(sdx,8).Float(),df.Elem(sdx,12).Float()}
				case 2:
				//y axis
				cp = []float64{df.Elem(sdx,8).Float(),df.Elem(sdx,13).Float()}
			}
			case 4:
			//3d t - a
			switch ax{
				case 1:
				//major axis of bending - x
				cp = []float64{df.Elem(sdx,8).Float(),df.Elem(sdx,12).Float()}
				case 2:
				//y axis
				cp = []float64{df.Elem(sdx,8).Float(),df.Elem(sdx,13).Float()}
			}
			case 5:
			//3d g - a, etc
			case 6:
			//3d f
		}
	case sectyp == 7 || sectyp == 8:
		//bs steel sections
		df, err = StlDfBs(sectyp)		
		
		if err != nil{
			return 
		}
		switch frmtyp{
			case 1:
		//1d b
			case 2:
		//2d t
			case 3:
		//2d f
			case 4:
		//3d t
			case 5:
		//3d g
			case 6:
			//3d f
		}
	}
	return
}
