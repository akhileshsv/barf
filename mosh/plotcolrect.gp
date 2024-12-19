if (ARG2 eq 'dumb') set term dumb ansi enhanced size 79,49
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'qt') set term qt persist; set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb "black" behind
if (ARG2 eq 'svg') set term svg; set output ARG4
set tics textcolor rgb "magenta"
set size ratio -1
set noborder
set autoscale
unset key
set xlabel 'mm' tc rgb "green"
set ylabel 'mm' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
set title 'rectangular rcc column'
set offsets graph 0.1,0.1,0.1,0.1
plot ARG1 index 0 using 1:2:3:4:5 w boxxyerrorbars lc var,\
     ARG1 index 1 using 1:2:3:4:5 w ellipses lc var,\
     ARG1 index 1 using 1:2,\
