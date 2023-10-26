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

                        v0 - (super rough demo)
```
# Table of Contents

1.  [Overview](#orgabfeceb)
2.  [Getting started](#orgb389b3e)
    1.  [Prerequisites](#org68372ef)
    2.  [Usage](#org09e35ec)
3.  [References](#org502a830)

BARF is a collection of programs for structural analysis and design written in Go.


<a id="orgabfeceb"></a>

# Overview

-   **Direct stiffness analysis (/kass):** Direct stiffness analysis of (2d/3d) bar member frameworks.
-   **RCC design (/mosh):** Design of rcc slabs, beams, columns and footings as per is456 and bs8110.
-   **Steel design (/bash):** Design of steel beam and column members as per bs449.
-   **Timber design (/tmbr):** Design of timber beam and column members as per is883.


<a id="orgb389b3e"></a>

# Getting started


<a id="org68372ef"></a>

## Prerequisites

-   Go
    - follow instructions at <https://go.dev/doc/install>
-   Gnuplot
    - follow instructions at <https://sourceforge.net/projects/gnuplot/files/gnuplot/>


<a id="org09e35ec"></a>

## Usage

Clone this repo    
```git clone https://github.com/akhileshsv/barf```

Goto folder    
```cd barf```

Run tui menu    
```go run main.go -tui```  

![]
(https://github.com/akhileshsv/barf/blob/main/install.gif)    

![]
(https://github.com/akhileshsv/barf/blob/main/mosh.gif)     

![]
(https://github.com/akhileshsv/barf/blob/main/bash.gif)     

![]
(https://github.com/akhileshsv/barf/blob/main/tmbr.gif)     


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

