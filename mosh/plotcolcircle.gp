if (ARG2 eq 'dumb') set term dumb ansi enhanced size 79,49
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'qt') set term qt persist; set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb "black" behind
if (ARG2 eq 'svg') set term svg font 'Kongtext,5'; set output ARG4
set tics textcolor rgb "magenta"
set noborder
set autoscale
set size ratio -1
unset key #bottom
set xlabel 'mm'
set ylabel 'mm'
set ytics autofreq nomirror tc lt 1
set title 'circular rcc column'
#set offsets graph 0.3,0.3,0.3,0.3
plot ARG1 index 0 using 1:2:3:4 w circles lc var,\
     ARG1 index 1 using 1:2,\
     ARG1 index 1 using 1:2:3:4:5 w ellipses lc var,\
     ARG1 index 1 using 1:2:3 w labels
