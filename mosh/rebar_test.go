package barf

import(
	"fmt"
	"testing"
	"strings"
	//"log"
		
)

func TestRebar(t *testing.T) {
	var rezstring string
	
	rezstring += "shah 4.2.1 steel area check :p\n"
	for _, dia := range StlDia {
		rezstring += fmt.Sprintf("dia %0.2f mm | area %0.2f | weight (kg/meter) %.2f\n",dia, RbrArea(dia), RbrUnitWt(dia))
	}
	
	rezstring += "\nshah 4.2.2 rebar combos\n"
	asreqs := []float64{2513,1243,12898}
	for _, asreq := range asreqs {
		rez, mindia := RbrSelect(asreq,false)
		if len(rez) != 0 {
			for _, r := range rez {
				rezstring += fmt.Sprintf("dia %0.1f nbar %0.1f asprov %0.1f diff %0.1f\n",r[1],r[2],r[3],r[4])
			}
			rezstring += fmt.Sprintf("min dia %0.1f \n",mindia)
		} else {
			rezstring += fmt.Sprintf("ERRORE,errore --> impossible area %0.1f",asreq)
		}
		
	}

	rezstring += "\nshah 4.2.3 n rows of bars (beam)\n"
	nbars := []float64{3,4,7}
	dia := 16.0
	for _, nbar := range nbars {
		rez := RbrNRows(150, 600, nbar, dia)
		rezstring += fmt.Sprintf("nbar %0.1f dia %0.1f nlayer %0.1f asprov %0.1f dused %0.1f efcvr %0.1f efdp %0.1f bw %0.1f\n", nbar, dia, rez[0], rez[1], 600.0, rez[2], rez[3], 150.0)
	}
	rezstring += "\nshah 4.2.4 slab rebar dia-spacing combos\n"
	rezstring += "\ndused 100, asreq 380\n"
	slbbarez, minidx := SlabRbrDia(100,380)
	rezstring += fmt.Sprintf("%0.1f %v",slbbarez,minidx)
	rezstring += fmt.Sprintf("\nMin index %v %0.1f\n", minidx, StlDia[minidx])
	rezstring += "\ndused 110, asreq 126\n"
	slbbarez, minidx = SlabRbrDia(110,126)
	rezstring += fmt.Sprintf("%0.1f %v",slbbarez,minidx)
	rezstring += fmt.Sprintf("\nMin index %v %0.1f\n", minidx, StlDia[minidx])
	rezstring += "\ndused 110, asreq 434\n"
	slbbarez, minidx = SlabRbrDia(110,434)
	rezstring += fmt.Sprintf("%0.1f %v",slbbarez,minidx)
	rezstring += fmt.Sprintf("\nMin index %v %0.1f\n", minidx, StlDia[minidx])
	rezstring += "\ndused 110, asreq 5068 (returns nothing, hence good news)\n"
	slbbarez, _ = SlabRbrDia(100,5068)
	if len(slbbarez) == 0 {
		rezstring += "ERRORE,errore-->change slab depth\n"
	}
	rezstring += "\nshah 4.3.1 design coffecients\n"
	kumax, zumax, rumax, ptmax := BalSecKis(15,415,0,1.15,0.361,0.416)
	rezstring += fmt.Sprintf("NA Factor :  %f LA Factor :  %f MR Factor:  %f  PT Max:  %f \n",kumax, zumax, rumax, ptmax)
	
	rezstring += "\nshah 4.3.2 steel stress-strain interpolation\n"
	
	fys := []float64{250,250,415,415,500,500}
	escs := []float64{0.001,0.002,0.00144,0.0014435,0.0017,0.0030}
	for idx, fy := range fys {
		rezstring += fmt.Sprintf("Strain in steel esc %0.6f fy %0.2f fsc %0.2f N/mm2 \n", escs[idx], fy, RbrFrcIs(fy, escs[idx]))
	}
	rezstring += "\nshah 4.3.3 slab serviceability check\n"
	slabtyp := 1; endc := 1; fy := 250.0; ll := 4.0; lx := 2500.0; ptreq := 0.16; efcvr := 30.0
	dreq := RccSlabServeChk(slabtyp, endc, fy, ll, lx, ptreq,efcvr) 
	rezstring += fmt.Sprintf("one-way slab\n slabtyp %v endc %v fy %0.1f ll %0.1f lx %0.1f ptreq %0.4f efcvr %0.1f mm dreq %0.0f mm\n", slabtyp, endc,fy,ll,lx, ptreq, efcvr, dreq)
	slabtyp = 2; fy = 415.0; endc = 10
	dreq = RccSlabServeChk(slabtyp, endc, fy, ll, lx, ptreq,efcvr)
	
	rezstring += fmt.Sprintf("two-way slab\n slabtyp %v endc %v fy %0.1f ll %0.1f lx %0.1f ptreq %0.4f efcvr %0.1f mm dreq %0.0f mm\n", slabtyp, endc,fy,ll,lx, ptreq, efcvr, dreq)

	rezstring += "\nshah 4.3.4 development length check\n"
	fck := 15.0 ; fy = 250.0 ; m1 := 44.18 ; v := 48.13 ; effd := 400.0; dia = 16.0; bs := 230.0; rcom := 1; barcut := 0 
	rezstring += fmt.Sprintf("fck %0.2f fy %0.2f m1 %0.2f v %0.2f effd %0.2f dia %0.2f bs %0.2f rcom %v barcut %v \n dev length ok- %v \n", fck, fy, m1, v, effd, dia, bs, rcom, barcut, RccDevLenChk(fck, fy, m1, v, effd, dia, bs, rcom, barcut))
	
	v = 44.36 ; rcom = 0; barcut = 1 
	rezstring += fmt.Sprintf("fck %0.2f fy %0.2f m1 %0.2f v %0.2f effd %0.2f dia %0.2f bs %0.2f rcom %v barcut %v \n dev length ok- %v \n", fck, fy, m1, v, effd, dia, bs, rcom, barcut, RccDevLenChk(fck, fy, m1, v, effd, dia, bs, rcom, barcut))
	
        rezstring += "\nshah 4.3.5 shear stirrup/link diameter/spacing selection\n"
	idxs, nlegs, spacing := RbrShearLink(415.0, 230.0, 685.0, 106.54, 150.0, 1, 2, 2, 1) 
	for i,didx := range idxs {
		rezstring += fmt.Sprintf("Bar dia %0.0f %v -legged stirrups @ %0.2f mm c-c \n",StlDia[didx],nlegs[i],spacing[i])
	}
	
	idxs, nlegs, spacing = RbrShearLink(415.0, 230.0, 685.0, 106.54, 150.0, 0, 2, 2, 0)
	for i,didx := range idxs {
		rezstring += fmt.Sprintf("Bar dia %0.0f %v -legged stirrups @ %0.2f mm c-c \n",StlDia[didx],nlegs[i],spacing[i])
	}

	wantstring := `shah 4.2.1 steel area check :p
dia 6.00 mm | area 28.27 | weight (kg/meter) 0.22
dia 8.00 mm | area 50.27 | weight (kg/meter) 0.39
dia 10.00 mm | area 78.54 | weight (kg/meter) 0.62
dia 12.00 mm | area 113.10 | weight (kg/meter) 0.89
dia 16.00 mm | area 201.06 | weight (kg/meter) 1.58
dia 18.00 mm | area 254.47 | weight (kg/meter) 2.00
dia 20.00 mm | area 314.16 | weight (kg/meter) 2.47
dia 22.00 mm | area 380.13 | weight (kg/meter) 2.98
dia 25.00 mm | area 490.87 | weight (kg/meter) 3.85
dia 28.00 mm | area 615.75 | weight (kg/meter) 4.83
dia 32.00 mm | area 804.25 | weight (kg/meter) 6.31
dia 40.00 mm | area 1256.64 | weight (kg/meter) 9.86

shah 4.2.2 rebar combos
dia 16.0 nbar 13.0 asprov 2613.8 diff 100.8
dia 18.0 nbar 10.0 asprov 2544.7 diff 31.7
dia 20.0 nbar 8.0 asprov 2513.3 diff 0.3
dia 22.0 nbar 7.0 asprov 2660.9 diff 147.9
dia 25.0 nbar 6.0 asprov 2945.2 diff 432.2
min dia 20.0
dia 12.0 nbar 11.0 asprov 1244.1 diff 1.1
dia 16.0 nbar 7.0 asprov 1407.4 diff 164.4
dia 18.0 nbar 5.0 asprov 1272.3 diff 29.3
dia 20.0 nbar 4.0 asprov 1256.6 diff 13.6
dia 22.0 nbar 4.0 asprov 1520.5 diff 277.5
dia 25.0 nbar 3.0 asprov 1472.6 diff 229.6
min dia 12.0
ERRORE,errore --> impossible area 12898.0
shah 4.2.3 n rows of bars (beam)
nbar 3.0 dia 16.0 nlayer 1.0 asprov 603.2 dused 600.0 efcvr 33.0 efdp 567.0 bw 150.0
nbar 4.0 dia 16.0 nlayer 2.0 asprov 804.2 dused 600.0 efcvr 49.0 efdp 551.0 bw 150.0
nbar 7.0 dia 16.0 nlayer 3.0 asprov 1407.4 dused 600.0 efcvr 65.0 efdp 535.0 bw 150.0

shah 4.2.4 slab rebar dia-spacing combos

dused 100, asreq 380
[[8.0 132.0 380.8 0.8] [10.0 206.0 381.3 1.3] [12.0 240.0 471.2 91.2]] 1
Min index 1 8.0

dused 110, asreq 126
[[8.0 270.0 186.2 60.2] [10.0 270.0 290.9 164.9] [12.0 270.0 418.9 292.9]] 1
Min index 1 8.0

dused 110, asreq 434
[[8.0 115.0 437.1 3.1] [10.0 180.0 436.3 2.3] [12.0 260.0 435.0 1.0]] 3
Min index 3 12.0

dused 110, asreq 5068 (returns nothing, hence good news)
ERRORE,errore-->change slab depth

shah 4.3.1 design coffecients
NA Factor :  0.479167 LA Factor :  0.800667 MR Factor:  2.077480  PT Max:  0.007190

shah 4.3.2 steel stress-strain interpolation
Strain in steel esc 0.001000 fy 250.00 fsc 200.00 N/mm2
Strain in steel esc 0.002000 fy 250.00 fsc 217.39 N/mm2
Strain in steel esc 0.001440 fy 415.00 fsc 288.00 N/mm2
Strain in steel esc 0.001443 fy 415.00 fsc 288.70 N/mm2
Strain in steel esc 0.001700 fy 500.00 fsc 340.00 N/mm2
Strain in steel esc 0.003000 fy 500.00 fsc 422.65 N/mm2

shah 4.3.3 slab serviceability check
one-way slab
 slabtyp 1 endc 1 fy 250.0 ll 4.0 lx 2.5 ptreq 0.1600 efcvr 30.0 mm dreq 125 mm
two-way slab
 slabtyp 2 endc 10 fy 415.0 ll 4.0 lx 2.5 ptreq 0.1600 efcvr 30.0 mm dreq 175 mm

shah 4.3.4 development length check
fck 15.00 fy 250.00 m1 44.18 v 48.13 effd 400.00 dia 16.00 bs 230.00 rcom 1 barcut 0
 dev length ok- true
fck 15.00 fy 250.00 m1 44.18 v 44.36 effd 400.00 dia 16.00 bs 230.00 rcom 0 barcut 1
 dev length ok- true

shah 4.3.5 shear stirrup/link diameter/spacing selection
Bar dia 8 2 -legged stirrups @ 240.00 mm c-c
Bar dia 10 2 -legged stirrups @ 370.00 mm c-c
Bar dia 6 2 -legged stirrups @ 160.00 mm c-c
Bar dia 8 2 -legged stirrups @ 280.00 mm c-c
Bar dia 10 2 -legged stirrups @ 430.00 mm c-c
`
	t.Logf(rezstring)
	if strings.Compare(rezstring, wantstring) != 1 {	
		t.Errorf("Rebar helper functions test failed")
		fmt.Println(rezstring)
	}

}



