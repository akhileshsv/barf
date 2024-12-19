```
                  ___          ___         ___   
     _____       /  /\        /  /\       /  /\  
    /  /::\     /  /::\      /  /::\     /  /:/_ 
   /  /:/\:\   /  /:/\:\    /  /:/\:\   /  /:/ /\
  /  /:/~/::\ /  /:/~/::\  /  /:/~/:/  /  /:/ /:/
 /__/:/ /:/\:/__/:/ /:/\:\/__/:/ /:/__/__/:/ /:/ 
 \  \:\/:/~/:|  \:\/:/__\/\  \:\/:::::|  \:\/:/  
  \  \::/ /:/ \  \::/      \  \::/~~~~ \  \::/   
   \  \:\/:/   \  \:\       \  \:\      \  \:\   
    \  \::/     \  \:\       \  \:\      \  \:\  
     \__\/       \__\/        \__\/       \__\/  

                        v1 - (rough demo)
```
# Table of Contents

1.  [Overview](#intro)
2.  [Getting started](#start)
    1.  [Prerequisites](#req)
    2.  [Usage](#use)
3.  [References](#ref)

BARF is a collection of programs for structural analysis and design written in Go.


<a id="intro"></a>


# Overview

-   **Direct stiffness analysis (/kass):** Direct stiffness analysis of (2d/3d) bar member frameworks.
-   **RCC design (/mosh):** Design of rcc slabs, beams, columns and footings as per is456 and bs8110.
-   **Steel design (/bash):** Design of steel beam and column members as per bs449.
-   **Timber design (/tmbr):** Design of timber beam and column members as per is883.


<a id="start"></a>


# Getting started


<a id="req"></a>

## Prerequisites

-   Go
    - follow instructions at <https://go.dev/doc/install>
-   Gnuplot
    - follow instructions at <https://sourceforge.net/projects/gnuplot/files/gnuplot/>


<a id="use"></a>

## Usage

Clone and build this repo -   

```
git clone https://github.com/akhileshsv/barf  
cd barf  
go build -o barf.exe

```  
Run tui menu (uses https://github.com/AlecAivazis/survey) with -  

```
./barf.exe -tui

```
Note - using "read json txt" as an input option copies the base struct to clipboard and opens a new editor window to edit struct fields.  
Hit "?" to see an explanation of individual struct fields (as seen in the gifs below); paste the struct into the editor (ctrl + v) and save the edited file before exiting the editor window.

-   Flags
    -   inf (string) - input json file path
    -   term (string) - gnuplot terminal string (supported - "dumb","mono","qt","svg")
    -   calc (bool) - stiffness analysis flag (kass.Model json input)
    -   rcc (string) - rcc design flag
    -   stl (string) - steel design flag
    -   wood (string) - wood design flag
    -   cmdz (string) - flag sub-commands
    -   tui (bool) - start text menu

<table border="2" cellspacing="0" cellpadding="6" rules="groups" frame="hsides">


<colgroup>
<col  class="org-left" />

<col  class="org-left" />

<col  class="org-left" />

<col  class="org-left" />

<col  class="org-left" />

<col  class="org-left" />
</colgroup>
<thead>
<tr>
<th scope="col" class="org-left">flag</th>
<th scope="col" class="org-left">info</th>
<th scope="col" class="org-left">type</th>
<th scope="col" class="org-left">vals</th>
<th scope="col" class="org-left">cmdz</th>
<th scope="col" class="org-left">info</th>
</tr>
</thead>

<tbody>
<tr>
<td class="org-left">calc</td>
<td class="org-left">analysis</td>
<td class="org-left">bool</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">basic model direct stiffness analysis</td>
</tr>


<tr>
<td class="org-left">calc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">ep</td>
<td class="org-left">elastic plastic beam/frame analysis</td>
</tr>


<tr>
<td class="org-left">calc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">np</td>
<td class="org-left">non uniform beam/frame analysis</td>
</tr>


<tr>
<td class="org-left">calc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">blt, bolt</td>
<td class="org-left">bolt group analysis</td>
</tr>


<tr>
<td class="org-left">calc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">wld, weld</td>
<td class="org-left">weld group analysis</td>
</tr>


<tr>
<td class="org-left">rcc</td>
<td class="org-left">rcc design</td>
<td class="org-left">string</td>
<td class="org-left">slb,slab</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">rcc slab design</td>
</tr>


<tr>
<td class="org-left">rcc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">cb, cbeam</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">rcc continuous beam analysis and design</td>
</tr>


<tr>
<td class="org-left">rcc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">bm, beam</td>
<td class="org-left">az, analyze</td>
<td class="org-left">rcc beam section analysis</td>
</tr>


<tr>
<td class="org-left">rcc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">bm, beam</td>
<td class="org-left">dz, design</td>
<td class="org-left">rcc beam section design</td>
</tr>


<tr>
<td class="org-left">rcc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">col, column</td>
<td class="org-left">az, analyze</td>
<td class="org-left">rcc column section analysis</td>
</tr>


<tr>
<td class="org-left">rcc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">col, column</td>
<td class="org-left">dz, design</td>
<td class="org-left">rcc column section design</td>
</tr>


<tr>
<td class="org-left">rcc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">ftng, footing</td>
<td class="org-left">az, analyze</td>
<td class="org-left">rcc footing analysis</td>
</tr>


<tr>
<td class="org-left">rcc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">ftng, footing</td>
<td class="org-left">dz, design</td>
<td class="org-left">rcc footing design</td>
</tr>


<tr>
<td class="org-left">rcc</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">sf, subfrm</td>
<td class="org-left">az, analyze</td>
<td class="org-left">rcc subframe analysis</td>
</tr>


<tr>
<td class="org-left">stl</td>
<td class="org-left">steel design</td>
<td class="org-left">string</td>
<td class="org-left">col, column</td>
<td class="org-left">az, analyze</td>
<td class="org-left">steel column check</td>
</tr>


<tr>
<td class="org-left">stl</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">col, column</td>
<td class="org-left">dz, design</td>
<td class="org-left">steel column design</td>
</tr>


<tr>
<td class="org-left">stl</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">bm, beam</td>
<td class="org-left">az, analyze</td>
<td class="org-left">steel beam check</td>
</tr>


<tr>
<td class="org-left">stl</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">bm, beam</td>
<td class="org-left">dz, design</td>
<td class="org-left">steel beam design</td>
</tr>


<tr>
<td class="org-left">wood</td>
<td class="org-left">wood design</td>
<td class="org-left">string</td>
<td class="org-left">col, column</td>
<td class="org-left">dz, design</td>
<td class="org-left">solid timber column design</td>
</tr>


<tr>
<td class="org-left">wood</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">&#xa0;</td>
<td class="org-left">bm, beam</td>
<td class="org-left">dz, design</td>
<td class="org-left">solid timber beam design</td>
</tr>
</tbody>
</table>


<a id="org502a830"></a>

### Tests

#### kass (analysis)  

Navigate to the folder barf/kass and run the following command with the test name and the -v flag as below:  

```

cd barf/kass
go test -run "test name" -v
go test -run KFacMos -v

```

###### IN-PROGRESS calc np tests

-   [X] KFacMos - calc stiffness and distribution factors
-   [X] TestNpBm - non prismatic beam calc
-   [X] TestNpFrm - non prismatic frame calc


<a id="org727d411"></a>

###### IN-PROGRESS calc tests

-   [X] beamsfbm<sub>test</sub> - beam bm/sf funcs (hulse)
-   [X] Trs2d - 2d truss (chapter 3/4, kassimali)
-   [X] Trs2dJsd - 2d truss (chapter 3/4, kassimali)
-   [X] Bm1d - 1d beam (chapter 5, kassimali)
-   [X] BmJsd - joint support displacements calc (chapter 7, kassimali)
-   [X] Frm2d - 2d frame (chapter 6, kassimali)
-   [X] Frm2dJsd - 2d frame support disp./member releases (chapter 7, kassimali)
-   [X] Grd3d - 3d grid analysis (chapter 8, kassimali)
-   [X] Frm3d - 3d frame (chapter 8, kassimali)
-   [X] Trs3d - 3d truss (chapter 8, kassimali)
-   [X] add ons - chapter 7, kassimali examples (member releases, temp changes etc)
-   [X] NlCalcTrs2d - chapter 10 - special topics (non linear truss analysis) (kassimali)
-   [X] boltss - bolt group analysis (harrison)
-   [X] weldss - weld group analysis (harrison)
-   [X] CalcEp - elastic - plastic analysis (harrison)
    close the deflection bit (read)
-   [ ] EffLen - iterative column effective length calcs (harrison)
    finish ze test


<a id="org60bf48a"></a>

###### IN-PROGRESS section tests

-   [X] SecProp calc - mosley general section properties calculation
-   [X] SecGen - 22 sections so far


<a id="orga235c3e"></a>


#### mosh (rcc design)  

Navigate to the folder barf/mosh and run the following command with the test name and the -v flag as below:  

```


cd barf/mosh
go test -run "test name" -v
go test -run Slb1DBs -v

```


##### IN-PROGRESS v0 tests

<a id="org42b27c3"></a>

###### IN-PROGRESS slab

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

###### IN-PROGRESS beam

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

###### IN-PROGRESS column

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




<a id="org93ec1e4"></a>

###### IN-PROGRESS footing

1.  [X] FtngPadAz - pad footing analysis tests (hulse 6.1) (in footingdesign<sub>test</sub>)
2.  [ ] FtngBxOz - ozmen footing analysis tests
3.  [ ] FtngDzRojas - basic footing tests (rojas, hulse, subramanian)


<a id="org9a6a2fd"></a>

###### IN-PROGRESS cbeam

1.  [ ] CBmOpt - cbeam opt test (govindraj ex 1 n 2)
2.  [ ] CBmDz - cbeam design funcs


<a id="org776bf36"></a>

###### IN-PROGRESS subframe

1.  [ ] TestFltSlb - flat slab design tests (subframe<sub>test</sub>)
2.  [X] SubFrmDz - subframe design tests
    -   [X] SubFrmMosley - subframe analysis test

<a id="ref"></a>

# References

1.  Aslam Kassimali - Matrix Analysis of Structures , Second Edition - CL Engineering (2011)
2.  H. B. Harrison - Structural Analysis and Design, Some Microcomputer Applications-Elsevier Ltd, Pergamon (1990)
3.  W. H. Mosley, W. J. Spencer - Microcomputer Applications in Structural Engineering-Macmillan Education UK (1984)
4.  R. Hulse, W. H. Mosley - Reinforced Concrete Design by Computer-Macmillan Education UK (1986)
5.  Dr. V.L Shah - Computer Aided Design in Reinforced Concrete - Structures Publications (1998)
6.  Subramanian, Narayanan - Design of Reinforced Concrete Structures-Oxford University Press (2013)
7.  A. Allen - Reinforced Concrete Design to BS 8110 Simply Explained-CRC Press (1988)
8.  Abel O. Olorunnisola - Design of Structural Elements with Tropical Hardwoods - Springer (2017)
9.  Arnulfo Luevanos Rojas - Design of isolated rectangular footings of rectangular form using a new model (2013)

