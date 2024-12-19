if (ARG2 eq 'dumb') set term dumb ansi enhanced size 79,79
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'mono') set term dumb mono size 79,79
if (ARG2 eq 'qt') set term qt persist font "Courier,3";set size ratio -1 
if (ARG2 eq 'wxt') set term wxt persist
if (ARG2 eq 'svg') set term svg font "Kongtext,5"; set output ARG4; set size ratio -1 
if (ARG2 eq 'svgmono') set term svg font "Kongtext,5"; set output ARG4; set mono; set size ratio -1 
if (ARG2 eq 'dxf') set term dxf; set size 1000,1000; set output ARG4; set size ratio -1 
set tics textcolor rgb "magenta"
set noborder
set key outside right bottom
set xlabel 'mm' tc rgb "green"
set ylabel 'mm' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
set title ARG3
set offsets graph 0.1,0.1,0.1,0.1
set multiplot layout 1,2 rowsfirst
if (ARG2 eq 'dumb') || (ARG2 eq 'mono'){
set multiplot layout 2,1 rowsfirst
plot ARG1 index 0 using 1:2:3:4:5 w vectors nohead lc var notitle,\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead lc var title "main steel",\
     ARG1 index 2 using 1:2:3:4:5 w vectors nohead dt 2 lw 0.4 lc var title "ties",\
     ARG1 index 3 using 1:2:3:4:5 w vectors heads lc var notitle,\
     ARG1 index 4 using 1:2:3 w labels notitle
plot ARG1 index 5 using 1:2:3:4:5 w vectors nohead lw 2 lc var notitle,\
     ARG1 index 5 using ($1+$3)/2:($2+$4)/2:6 w labels notitle,\
     ARG1 index 6 using 1:2:3 w labels offset char 1,1 notitle,\
     ARG1 index 7 using 1:2 w points pt 15 title "bars",\
     ARG1 index 8 using 1:2 w lines dt 0 title "n.a",\
     ARG1 index 9 using 1:2:3:4 w vectors nohead dt 0 title "tie"    
    } else {
plot ARG1 index 0 using 1:2:3:4:5 w vectors nohead lc var notitle,\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead lc var title "main steel",\
     ARG1 index 2 using 1:2:3:4:5 w vectors nohead dt 2 lw 0.4 lc var title "ties",\
     ARG1 index 3 using 1:2:3:4:5 w vectors heads lc var notitle,\
     ARG1 index 4 using 1:2:3 w labels notitle
plot ARG1 index 5 using 1:2:3:4:5 w vectors nohead lw 2 lc var notitle,\
     ARG1 index 5 using ($1+$3/2):($2+$4/2):6 w labels tc rgb "green" notitle,\
     ARG1 index 6 using 1:2:3 w labels offset char 1,1 tc rgb "red" notitle,\
     ARG1 index 7 using 1:2:3:4 w ellipses title "bars",\
     ARG1 index 8 using 1:2 w lines dt 0 title "n.a",\
     ARG1 index 9 using 1:2:3:4 w vectors nohead lw 1 dt 0 title "tie"
}
