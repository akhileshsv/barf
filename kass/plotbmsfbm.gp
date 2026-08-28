if (ARG2 eq 'dumb') set term dumb ansirgb enhanced size 99,99
if (ARG2 eq 'caca') set term caca driver utf8 enhanced inverted size 99,79  
if (ARG2 eq 'qt') set term qt enhanced persist font "Kongtext,3"
if (ARG2 eq 'svg') set term svg dynamic font "Kongtext,3"; set output ARG4
if (ARG2 eq 'qtmono') set term qt persist font 'Kongtext,3'; set monochrome
if (ARG2 eq 'svgmono') set term svg dynamic font 'Kongtext,3'; set output ARG4; set monochrome
if (ARG2 eq 'dxf') set term dxf; set output ARG4
set multiplot layout 3,1
set border
unset tics 
unset colorbox
set offsets graph 0.1,0.1,0.1,0.1

if (ARG2 eq 'qt') set grid xtics ytics
set title 'shear force'
set xlabel 'mm'
set ylabel 'N'
plot ARG1 index 0 using 1:2 w lines title "shear",\
     ARG1 index 1 using 1:2:3 w labels offset char 1,1 title "max shear"
set title 'bending moment'
set xlabel 'mm'
set ylabel 'Nmm'
plot ARG1 index 0 using 1:3 w lines title "bending moment",\
     ARG1 index 2 using 1:2:3 w labels offset char 1,1 title "max BM"
set title 'deflection'
set xlabel 'mm'
set ylabel 'mm'
plot ARG1 index 0 using 1:4 w lines title "deflection",\
     ARG1 index 3 using 1:2:3 w labels offset char 1,1 title "max def."
unset multiplot
