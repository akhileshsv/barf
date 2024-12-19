#!/usr/bin/gnuplot
if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,49
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'qt') set term qt enhanced font 'Courier,8' persist
# set path of config snippets
#set loadpath '../draw/gnuplot-palettes'
#load 'noborder.cfg'
#load 'blues.pal'
#set tics textcolor rgb "magenta"
#set size ratio -1
#set noborder
#set autoscale
#unset key
set xlabel 'm' tc rgb "green"
set ylabel 'm' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
set title 'frame'
unset key
unset colorbox
set size ratio -1
set offsets graph 0.1,0.1,0.1,0.1
plot ARG1 index 0 using 1:2:3 w labels point pt 7 offset char 1,1 notitle,\
     ARG1 index 1 using 1:2:($3-$1):($4-$2):5 notitle w vectors nohead lw 2 dt 2 lc var,\
     ARG1 index 2 using 1:2 w points pt 9 ps 3 notitle,\
     ARG1 index 3 using 1:2:3:4 notitle w vectors,\
     ARG1 index 4 using 1:2:3:4:5 notitle w vectors lc var,\
     ARG1 index 5 using 1:2:3 w labels offset char 1,1 notitle,\
     ARG1 index 6 using 1:2:3 w circles notitle
