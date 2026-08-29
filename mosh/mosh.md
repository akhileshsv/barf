
# Table of Contents

1.  [ref](#org91292ee)
2.  [v0 tests](#orgd3d87d1)
    1.  [slab](#org42b27c3)
    2.  [beam](#orge4856b5)
    3.  [column](#org30b9835)
    4.  [footing](#org93ec1e4)
    5.  [cbeam](#org9a6a2fd)
    6.  [subframe](#org776bf36)
    7.  [frame 2d](#org1ac7e06)
    8.  [optimization](#org9e2b6df)
3.  [v1 tests](#org60008e4)


<a id="org91292ee"></a>

# ref

Mosley,Spencer - Microcomputer appplications in structural engineering
Hulse, Mosley, Bungey - Reinforced concrete design by computer
Shah, Karve - Computer aided design in reinforced concrete


<a id="orgd3d87d1"></a>

# IN-PROGRESS v0 tests


<a id="org42b27c3"></a>

## IN-PROGRESS slab

1.  [X] SlbSdrat - slab span depth ratio tests (in deflection<sub>test</sub>)
2.  [ ] TestSlb2DIs - slab 1 way tests (subramanian) (in slabdesignis<sub>test</sub>)
3.  [X] TestSlbSsShah - shah sec. 6.2 tests
4.  [X] TestSlb2WShah - shah sec. 6.3 2w slb tests
5.  [X] BalSecAst - balanced section area of steel tests (slabdesign<sub>test</sub>)
6.  [ ] Slb1DBs - bs code one way slab design
7.  [ ] Slb2DBs - bs code two way slab design
8.  [ ] Slb2BmCoeffBs - bs code 2 way coeff checks
9.  [ ] RSlb2Chk - waffle slab check tests (in slabcs<sub>test</sub>)
10. [ ] RSlb1Chk - ribbed slab check
11. [ ] CSlb1Dz - cs 1 way slab test
12. [ ] CSlb1DepthCs - cs 1 way slab depth test
13. [ ] SlbQuant - slab quant/table test (slab<sub>test</sub>)
14. [ ] SlbDraw - slab draw test
15. [ ] SlbEndC2W - slab 2 way end condition test (slabyield<sub>test</sub>)
16. [ ] TestSlb2WCoefComp - 2 way coeff comparison (bs/is/yield line)
17. [ ] SlbYld - yield line analysis tests
18. [ ] SlbYldRect - rect 2 way slab yield line tests


<a id="orge4856b5"></a>

## IN-PROGRESS beam

1.  [X] BmSecAzIs - (is code) shah beam section analysis (styp 1, 6, 7, 14)
2.  [X] BmDIs - (is code) shah beam design
3.  [X] BmDzBs - (bs code) hulse beam rebar design (styp ze usual)
4.  [X] BmAzGen - (bs code) hulse beam section analysis (again the usual styps)
5.  [ ] BmAzTaper - tapered beam analysis
6.  [X] BmBarGen - (general) - beam rebar generation funcs
7.  [X] BmAsvrat - beam stirrup area/spacing ratio (hulse)
8.  [X] BmTorDBs - design for torsion (hulse)
9.  [ ] BmShrIs - beam shear design (shah) - WRITE THIS
10. [X] BmShrDz - might be redundant, tests for beam shear(shah)
11. [ ] BmSdrat - beam span depth ratio tests (in deflection<sub>test</sub>)


<a id="org30b9835"></a>

## IN-PROGRESS column

1.  [X] ColDzBasic - basic (styp == 1/0) column design tests (in coldesignis<sub>test</sub>)
2.  [X] ColSizeIs - column sizing funcs
3.  [X] ColSecArXu - col. section area test (in coldesign<sub>test</sub>)
4.  [X] ColAzBs - hulse col. analysis (rect sect alone)
5.  [X] ColDzBs - hulse col. design
6.  [X] ColNMBs - hulse nm curves
7.  [X] ColEffHt - hulse eff. height calcs
8.  [X] ColSlmDBs - hulse slender column
9.  [ ] ColFlip - column 'flip' test (in col<sub>test</sub>)
10. [ ] ColBx - col. biaxial bending test
11. [ ] ColWeirdBs - weird column tests (start with styp 0)
12. [ ] ColAzGen - col gen. analysis funcs
13. [ ] ColStl - is this needs
14. [ ] ColBarGen - column rebar gen test (in colrebar<sub>test</sub>)
15. [ ] ColOpt - NOTHING has been written damn


<a id="org93ec1e4"></a>

## IN-PROGRESS footing

1.  [X] FtngPadAz - pad footing analysis tests (hulse 6.1) (in footingdesign<sub>test</sub>)
2.  [ ] FtngBxOz - ozmen footing analysis tests
3.  [ ] FtngDzRojas - basic footing tests (rojas, hulse, subramanian)


<a id="org9a6a2fd"></a>

## IN-PROGRESS cbeam

isn't this just beam? all the stuff in cbeam<sub>test</sub>

1.  [ ] CBmOpt - cbeam opt test (govindraj ex 1 n 2)
2.  [ ] CBmDz - cbeam design funcs
    holy shit there's a ton of examples this will take a year


<a id="org776bf36"></a>

## IN-PROGRESS subframe

1.  [ ] TestFltSlb - flat slab design tests (subframe<sub>test</sub>)
2.  [X] SubFrmDz - subframe design tests
    -   [X] SubFrmMosley - subframe analysis test


<a id="org1ac7e06"></a>

## TODO frame 2d

1.  [ ] TestFrm2dOpt - raka opt tests (frame2d<sub>test</sub>)
2.  [ ] FrmInit - frame init test
3.  [ ] Frm2dAllen - allen sec. 3.1 frame2d test
4.  [ ] FrmDzMosley - hulse/mosley frame2d test


<a id="org9e2b6df"></a>

## TODO optimization


<a id="org60008e4"></a>

# TODO v1 tests

