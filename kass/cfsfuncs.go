package barf

import (
	"os"
	"log"
	//"fmt"
	"runtime"
	"path/filepath"
	"encoding/csv"
	"errors"
	//"strings"
	"github.com/go-gota/gota/dataframe"
	//"github.com/go-gota/gota/series"
)

var (
	//1
	CfsStyps = []string{"c-lipped"}
)

//GetCfsDf takes in a section type index and returns the is code sheet as a dataframe 
func GetCfsDf(sectyp int) (dataframe.DataFrame, error){
	if sectyp > 1 || sectyp < 0{
		return dataframe.DataFrame{}, errors.New("invalid section type")
	}
	_, b, _, _:= runtime.Caller(0)
	basepath := filepath.Dir(b)
	tabledir := filepath.Join(basepath,"../data/cfs")
	
	f, _ := os.Open(filepath.Join(tabledir,CfsStyps[sectyp-1]+".csv"))
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return dataframe.DataFrame{}, err
	}
	df := dataframe.LoadRecords(records)
	return df, nil
}

//GetStlCp returns the cross section property slice given frm type, sectype and sheet index
func GetCfsCp(frmtyp, sectyp, sdx, ax int) (cp, dims []float64, err error){
	var df dataframe.DataFrame
	df, err = GetCfsDf(sectyp)
	if err != nil{
		return
	}
	if sectyp == 1 && sdx > 25{
		err = errors.New("invalid section index")
		return
	}
	//area = 9, ix = 12, iy = 16
	switch frmtyp{
		case 1:
		//1d b - iz
		switch ax{
			case 1:
			cp = []float64{df.Elem(sdx,11).Float()*1e4}
			case 2:
			cp = []float64{df.Elem(sdx,15).Float()*1e4}
		}
		case 2:
		//2d t - a, iz
		switch ax{
			case 1:
			//major axis of bending - x
			cp = []float64{df.Elem(sdx,8).Float()*1e2,df.Elem(sdx,11).Float()*1e4}
			case 2:
			//y axis
			cp = []float64{df.Elem(sdx,8).Float()*1e2,df.Elem(sdx,15).Float()*1e4}
		}
		case 3:
		//2d f - a, iz
		switch ax{
			case 1:
			//major axis of bending - x
			cp = []float64{df.Elem(sdx,8).Float()*1e2,df.Elem(sdx,11).Float()*1e4}
			case 2:
			//y axis
			cp = []float64{df.Elem(sdx,8).Float()*1e2,df.Elem(sdx,15).Float()*1e4}
		}
		case 4:
		//3d t - a
		switch ax{
			case 1:
			//major axis of bending - x
			cp = []float64{df.Elem(sdx,8).Float()*1e2,df.Elem(sdx,11).Float()*1e4}
			case 2:
			//y axis
			cp = []float64{df.Elem(sdx,8).Float()*1e2,df.Elem(sdx,15).Float()*1e4}
		}
		case 5:
		//3d g - a, etc
		case 6:
		//3d f
	}
	dims = []float64{df.Elem(sdx,1).Float(),df.Elem(sdx,2).Float(),df.Elem(sdx,4).Float(),df.Elem(sdx,5).Float()}
	log.Println("calc wt->",df.Elem(sdx,8).Float()*1.0*7850.0/1e4," kg vs->",df.Elem(sdx,9).Float())
	return
}
