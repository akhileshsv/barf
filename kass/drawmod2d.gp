if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,79
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'qt') set term qt enhanced font 'Courier,7' persist
if (ARG2 eq 'mono') set term dumb mono size 99,79
if (ARG2 eq 'svg') set term svg dynamic font "Courier,7" #background rgb "black"
if (ARG2 eq 'svg') set output ARG4
# set path of config snippets
#set loadpath '../draw/gnuplot-palettes'
#load 'noborder.cfg'
#load 'blues.pal'
#set tics textcolor rgb "magenta"
#set size ratio -1

set noborder
set autoscale
#unset key
set xlabel 'm' tc rgb "green"
set ylabel 'm' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
set title ARG3 tc rgb "magenta" font "Courier,8"
set key bottom textcolor rgb "blue"
unset colorbox
#set size ratio -1
set offsets graph 0.1,0.1,0.1,0.1
plot ARG1 index 0 using 1:2:3 w labels tc rgb "dark-cyan" point lc rgb "purple" pt 7 offset char -2,-2 title "nodes",\
     ARG1 index 1 using 1:2:($3-$1):($4-$2):6 w vectors nohead lw 1 lc var title "mems",\
     ARG1 index 2 using 1:2 w points pt 9 ps 2 title "sups",\
     ARG1 index 3 using 1:2:3:4:5 w vectors backhead lc var lw 1 title "load",\
     ARG1 index 4 using 1:2:3 w labels offset 2,2 tc rgb "green" notitle,\
     ARG1 index 5 using 1:2:3 w circles title "moment"
#ENABLE THIS FOR INTERACTIVE PLOTS
if (ARG2 eq "qt") pause mouse close
