if (ARG2 eq 'dumb') set term dumb ansi enhanced size 79,49
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'mono') set term dumb mono size 79,49
if (ARG2 eq 'qt') set term qt persist
if (ARG2 eq 'wxt') set term wxt persist
set tics textcolor rgb "magenta"
# set size ratio 1
set noborder
set autoscale
set key outside right bottom
set xlabel 'mm' tc rgb "green"
set ylabel 'mm' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
set title ARG3
set offsets graph 0.1,0.1,0.1,0.1
plot ARG1 index 0 using 1:2:3:4:5 w vectors nohead lw 2 lc var notitle,\
     ARG1 index 0 using ($1+$3/2):($2+$4/2):6 w labels tc rgb "green" notitle,\
     ARG1 index 1 using 1:2:3 w labels offset char 1,1 tc rgb "red" notitle,\
     ARG1 index 2 using 1:2:3 w circles notitle

