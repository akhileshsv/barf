if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,99
if (ARG2 eq 'caca') set term caca driver utf8 enhanced inverted size 99,79  
if (ARG2 eq 'qt') set term qt enhanced persist font "Courier,3"#; set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb "black" behind
if (ARG2 eq 'svg') set term svg dynamic font "Kongtext,5"; set output ARG4
if (ARG2 eq 'qtmono') set term qt persist font 'Courier,8'; set monochrome
if (ARG2 eq 'svgmono') set term svg dynamic font 'Kongtext,5'; set output ARG4; set monochrome
if (ARG2 eq 'dxf') set term dxf; set output ARG4
set tics textcolor rgb "green"
set title ARG3
#set size ratio 0.5
set noborder
set autoscale
set multiplot layout 3,1 rowsfirst
#unset key
set xlabel 'mm' tc rgb "green"
set ylabel 'KN' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
#unset xtics
#set xtics autofreq nomirror scale 0.25 rotate tc lt 1
#set lmargin 1
#set rmargin 1
#set tmargin 1
#set bmargin 1
#set offsets graph 0.05,0.05,0.05,0.05
set mxtics 5
set mytics 5
if (ARG2 eq 'qt') || (ARG2 eq 'svg') set grid xtics ytics
if (ARG2 eq "dxf") unset xtics; unset ytics; set font "=0.01 m"
set title 'shear envelope'
set key bottom
plot ARG1 index 0 using 1:2 w lines title "shear",\
     ARG1 index 1 using 1:2:3:4:5 w boxxyerrorbars lc var notitle,\
     ARG1 index 2 using 1:2:3 w labels tc rgb "magenta" notitle,\
     ARG1 index 3 using 1:2 w points pt 8 ps 2 notitle,\
     ARG1 index 4 using 1:2 w points pt 1 notitle,\
     ARG1 index 4 using 1:2:3:4 w labels tc rgb "blue"  notitle
set title 'hogging moment envelope'
#set datafile missing '0.0'
#set offsets graph 0.05,0.05,0.05,0.05
plot ARG1 index 0 using 1:3 w lines title "hogging moment",\
     ARG1 index 1 using 1:2:3:4:5 w boxxyerrorbars lc var notitle,\
     ARG1 index 2 using 1:2:3 w labels tc rgb "magenta" notitle,\
     ARG1 index 3 using 1:2 w points pt 8 ps 2 notitle,\
     ARG1 index 5 using 1:2 w points pt 1 notitle,\
     ARG1 index 5 using 1:2:3 w labels offset 1 notitle
set title 'sagging moment envelope'
#set datafile missing '0.0'
#set offsets graph 0.05,0.05,0.05,0.05
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
plot ARG1 index 0 using 1:4 w lines title "sagging moment",\
     ARG1 index 1 using 1:2:3:4:5 w boxxyerrorbars lc var notitle,\
     ARG1 index 2 using 1:2:3 w labels tc rgb "magenta"  notitle,\
     ARG1 index 3 using 1:2 w points pt 8 ps 2 notitle,\
     ARG1 index 6 using 1:2 w points pt 1 notitle,\
     ARG1 index 6 using 1:2:3 w labels offset 1 tc rgb "blue" notitle
unset multiplot
