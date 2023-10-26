if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,79
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'qt') set term qt enhanced font 'Courier,3' persist
if (ARG2 eq 'mono') set term dumb mono size 99,79
if (ARG2 eq 'svg') set term svg dynamic font "Courier,3" #background rgb "black"
if (ARG2 eq 'svg') set output ARG4
set size ratio -1
set noborder
#set autoscale
#unset key
set ytics autofreq nomirror scale 0.25 tc lt 1
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
#set title ARG3 tc rgb "magenta" font "Courier,8"
set key bottom textcolor rgb "blue"
unset colorbox
#set size ratio -1
set xlabel 'm' tc rgb "green"
set ylabel 'm' tc rgb "green"
set offsets graph 0.1,0.1,0.1,0.1
set multiplot layout 3,1
set title sprintf('%s shear', ARG3) font "Courier,4"
plot ARG1 index 0 using 1:2:($3-$1):($4-$2):5 w vectors nohead lw 2 lc var title "mems",\
     ARG1 index 1 using 1:2 w linespoints pt 0 dt 1 title "sf"
     #ARG1 index 2 using 1:2:3 w labels tc rgb "red" notitle
set title sprintf('%s shear', ARG3) font "Courier,4"
plot ARG1 index 0 using 1:2:($3-$1):($4-$2):5 w vectors nohead lw 2 lc var title "mems",\
     ARG1 index 1 using 3:4 w linespoints pt 0 dt 1 title "bm"
     #ARG1 index 2 using 4:5:6 w labels tc rgb "red" notitle,\
     #ARG1 index 2 using 10:11:12 w labels tc rgb "blue" notitle
set title sprintf('%s shear', ARG3) font "Courier,4"
plot ARG1 index 0 using 1:2:($3-$1):($4-$2):5 w vectors nohead lw 2 lc var title "mems",\
     ARG1 index 1 using 5:6 w linespoints pt 0 dt 1 title "dx",\
     #ARG1 index 2 using 1:2:($3-$1):($4-$2) w vectors nohead lw 1 dt 1 title "def. shape"
#ENABLE THIS FOR INTERACTIVE PLOTS
unset multiplot
#if (ARG2 eq "qt")||(ARG2 eq "wxt") pause mouse close
