if (ARG2 eq 'dumb') set term dumb ansi enhanced size 160,90
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'mono') set term dumb mono size 79,79
if (ARG2 eq 'qt') set term qt persist font "Courier,5";set size 1,1 
if (ARG2 eq 'wxt') set term wxt persist
if (ARG2 eq 'svg') set term svg dynamic lw 0.75 font "Kongtext,5"; set output ARG4; set size ratio 0.75
if (ARG2 eq 'svgmono') set term svg dynamic lw 0.4 font "Kongtext,5"; set output ARG4; set mono; set size ratio 0.75
if (ARG2 eq 'dxf') set term dxf; set output ARG4; set autoscale; set font "Arial,3"
set tics textcolor rgb "magenta"
set noborder
set key outside right bottom
set xlabel 'mm' tc rgb "blue"
set ylabel 'mm' tc rgb "blue"
set ytics autofreq nomirror scale 0.25 tc lt 1
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
set title ARG3
set offsets graph 0.1,0.1,0.1,0.1
if (ARG2 eq 'dumb') || (ARG2 eq 'mono'){
set size ratio 0.75
#set multiplot layout 2,1 rowsfirst
plot ARG1 index 0 using 1:2:3:4:5 w vectors nohead notitle,\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead lc var title "main bar",\
     ARG1 index 2 using 1:2:3:4 w vectors heads notitle,\
     ARG1 index 3 using 1:2:3 w labels notitle,\
     ARG1 index 4 using 1:2 w points pt 0 title "dist bar"
} else {
plot ARG1 index 0 using 1:2:3:4:5 w vectors nohead notitle,\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead lc var title "main (x) bar",\
     ARG1 index 2 using 1:2:3:4 w vectors heads lw 0.5 notitle,\
     ARG1 index 3 using 1:2:3 w labels notitle,\
     ARG1 index 4 using 1:2:3 w circles title "dist (y) bar"
}

if (ARG2 eq 'qt') pause mouse close
