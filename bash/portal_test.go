package barf

import (
	"os"
	"fmt"
	"testing"
	"path/filepath"
	kass"barf/kass"
)

func TestPortalOpt(t *testing.T){
	var examples = []string{"saka2003-1","saka97-3"}
	var rezstring string
	var pf kass.Portal
	var err error
	dirname,_ := os.Getwd()
	datadir := filepath.Join(dirname,"../data/examples/kass")
	t.Log(ColorPurple,"testing portal frame opt(saka)\n",ColorReset)
	for i, ex := range examples{
		//if i != 0{continue}
		fname := filepath.Join(datadir,ex+".json")
		t.Log(ColorCyan,"example->",i+1,"file->",fname,"\n",ColorReset)
		pf, err = kass.ReadPortal(fname)
		if err != nil{
			t.Logf(fmt.Sprintf("%s",err))
		}
		err = PortalOpt(pf)
		if err != nil{
			t.Logf(fmt.Sprintf("%s",err))
		}		
	}
	wantstring := ``
	if rezstring != wantstring{
		t.Errorf("portal frame test failed")
	}


}

func HassanEx5(){
	mod := &kass.Model{
		Ncjt:3,
		Cmdz:[]string{"2df","mks","1"},
		Coords:[][]float64{
			{0,0},
			{40,0},
			{0,5},
			{40,5},
			{20,10.36},
		},
		Supports:[][]int{{1,-1,-1,0},{2,-1,-1,0}},
		Mprp:[][]int{{1,3,1,1,0},{2,4,1,1,0},{3,5,1,2,0},{4,5,1,2,0}},
		Jloads:[][]float64{},
		Msloads:[][]float64{{3,3,10.32,0,0,0},{3,6,2.77,0,0,0},{4,3,10.32,0,0,0},{4,6,2.77,0,0,0},},
		Em:[][]float64{{2.1e8}},
		Cp:[][]float64{{105e-4,47540e-8},{129e-4,61520e-8}},
	}
	frmrez, err := kass.CalcFrm2d(mod, 3)
	if err != nil{
		fmt.Println(err)
		return
	}
	report, _ := frmrez[6].(string)
	fmt.Println(report)
}
/*

	//log.Println(coords)
	//log.Println(mprp)
	//log.Println(bmvec)
	//log.Println(msup)
	//log.Println("bmvec-",bmvec)
	//log.Println("colvec-",colvec)
	//pltchn := make(chan string, 1)
	//go kass.PlotGenTrs(coords, mprp, pltchn)
	//pltstr := <-pltchn
	//log.Println(pltstr)

   
	lspan := f.Span/2.0
	ly := lspan
	ty := 1.0
	lbr := 200.0
	tbr := 20.0
	nsecs := 1
	grd := 43
	sectyp := 0
	brchck := false
	yeolde := false
	ldcases := [][]float64{{1.0,3.0,wdl,0,0,0,1},{1.0,3.0,wll,0,0,0,1}}
	rez := StlBmDBs(lspan, ly, ty, lbr, tbr, ldcases, sectyp, grd, nsecs, brchck, yeolde)
	bdx := rez[0]
	df := StlSecBs(sectyp)
	wb, arb, ixb := df.Elem(bdx,2).Float(), df.Elem(bdx,23).Float(), df.Elem(bdx,11).Float()
	fil := df.Filter(
		dataframe.F{Colname:"ix", Comparator:series.Greater, Comparando:1.5*ixb},
	)
	//fmt.Println(fil.Nrow())
	//fmt.Println(fil.Subset([]int{fil.Nrow()-1}))
	cdx := fil.Nrow()-1
	wc, arc, ixc := df.Elem(cdx,2).Float(), df.Elem(cdx,23).Float(), df.Elem(cdx,11).Float()
	//log.Println(wc, arc, ixc, wb, arb, ixb)
	//START WITH SAME SECTION
	cp = [][]float64{{arc*1e-4, ixc*1e-8},{arb*1e-4, ixb*1e-8}}
	em = [][]float64{{2.1*1e8}}
	//log.Println("section->",df.Elem(i,1))
	//log.Println("depth, web thickness->",df.Elem(i,3), df.Elem(i,6))
	//log.Println("area, zx, zy->",df.Elem(i,23),df.Elem(i,15), df.Elem(i,16))
	//log.Println("rx, ry->",df.Elem(i,13), df.Elem(i,14))
*/
