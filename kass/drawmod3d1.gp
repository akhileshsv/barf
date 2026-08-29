if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,79
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'qt') set term qt enhanced font 'Kongtext,5' persist
if (ARG2 eq 'mono') set term dumb mono size 99,79
if (ARG2 eq 'svg') set term svg dynamic font "Kongtext,5"; set output ARG4 #background rgb "black"
if (ARG2 eq 'dxf') set term dxf; set output ARG4
if (ARG2 eq 'svgmono') set term svg dynamic font "Kongtext,5"; set monochrome; set output ARG4
# set path of config snippets
#set loadpath '../draw/gnuplot-palettes'
#load 'noborder.cfg'
#load 'blues.pal'
#set tics textcolor rgb "magenta"
#set size ratio -1
#unset key

#set size ratio -1
#set offsets graph 0.1,0.1,0.1,0.1
set ticslevel 0
set noborder
set autoscale
set xlabel "x m" offset 1,1,1 rotate by 45 
set ylabel "z m" offset 1,1,1 rotate by 45
set zlabel "y m" offset 1,1,1 rotate by 45
set title ARG3 tc rgb "magenta"
set key outside bottom textcolor rgb "blue"
unset colorbox
set view 60, 30, 1, 1
splot ARG1 index 0 using 1:2:3:4 w labels tc rgb "dark-cyan" point lc rgb "purple" pt 7 offset char 2,2 title "nodes",\
      ARG1 index 1 using 1:2:3:($4-$1):($5-$2):($6-$3):7 w vectors nohead lc var lw 3 title "mems",\
      ARG1 index 1 using ($1+$4)/2:($5+$2)/2:($6+$3)/2:7 w labels offset 1,1 title "mems",\
      ARG1 index 2 using 1:2:3 w points pt 9 ps 2 title "sups",\
      ARG1 index 3 using 1:2:3:4:5:6:7 w vectors lc var lw 1 title "load",\
      ARG1 index 4 using 1:3:3:4:5:6:7 w vectors nohead lc var lw 1 notitle,\
      ARG1 index 5 using 1:2:3:4 w labels tc var offset 2,2 notitle,\
      ARG1 index 6 using 1:2:3 w points ps 3 title "moment" 
#ENABLE THIS FOR INTERACTIVE PLOTS
if (ARG2 eq "qt") pause mouse close
