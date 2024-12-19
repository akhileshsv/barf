if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,99
if (ARG2 eq 'caca') set term caca driver utf8 enhanced inverted size 99,79  
if (ARG2 eq 'qt') set term qt enhanced persist font "Courier,4"
if (ARG2 eq 'svg') set term svg font "Kongtext,5"; set output ARG4
if (ARG2 eq 'qtmono') set term qt persist font 'Courier,8'; set monochrome
if (ARG2 eq 'svgmono') set term svg font 'Kongtext,5'; set output ARG4; set monochrome
if (ARG2 eq 'dxf') set term dxf; set output ARG4
set tics textcolor rgb "green"
#set size ratio 0.5
set colorsequence classic
set noborder
set autoscale
set multiplot layout 3,1 rowsfirst
#unset key
set xlabel 'm' tc rgb "green"
set ylabel 'KN' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
#unset xtics
set mxtics 5
set mytics 5
if (ARG2 eq 'qt') || (ARG2 eq 'svg') set grid xtics ytics 
set offsets graph 0.05,0.05,0.05,0.05
if (ARG2 eq "dxf") unset xtics; unset ytics; set font "=0.01 m"
set title 'redis. shear envelope'
set key top
plot ARG1 index 0 using 1:2 w lines lt 1 lc 4 title "rd",\
     ARG1 index 0 using 1:5 w lines lt 0 lc 2 title "env",\
     ARG1 index 1 using 1:2:3:4:5 w boxxyerrorbars lc var notitle,\
     ARG1 index 2 using 1:2:3 w labels notitle,\
     ARG1 index 3 using 1:2 w points pt 1 ps 1 notitle,\
     ARG1 index 4 using 1:2:3:4 w labels offset char 1,1 notitle
set title 'redis. hogging moment envelope'
set offsets graph 0.05,0.05,0.05,0.05
plot ARG1 index 0 using 1:6 w lines lt 0 lc 2 title "env",\
     ARG1 index 0 using 1:8 w lines lt 0 lc 3 title "70%",\
     ARG1 index 0 using 1:3 w lines lt 1 lc 4 title "rd",\
     ARG1 index 1 using 1:2:3:4:5 w boxxyerrorbars lc var notitle,\
     ARG1 index 2 using 1:2:3 w labels notitle,\
     ARG1 index 3 using 1:2 w points pt 1 ps 1 notitle,\
     ARG1 index 5 using 1:2:3:4 w labels offset char 1,1 notitle
set title 'redis. sagging moment envelope'
set offsets graph 0.05,0.05,0.05,0.05
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
plot ARG1 index 0 using 1:4 w lines lt 1 lc 4 title "rd",\
     ARG1 index 0 using 1:7 w lines lt 0 lc 2 title "env",\
     ARG1 index 0 using 1:9 w lines lt 0 lc 3 title "70%",\
     ARG1 index 1 using 1:2:3:4:5 w boxxyerrorbars lc var notitle,\
     ARG1 index 2 using 1:2:3 w labels notitle,\
     ARG1 index 3 using 1:2 w points pt 1 ps 1 notitle,\
     ARG1 index 6 using 1:2:3 w labels offset char 1,1 notitle
#ARG1 index 6 using 1:2 w points pt 14 ps 1 notitle,\
unset multiplot
