if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,49
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'qt') set term qt enhanced font 'Kongtext,4' persist
if (ARG2 eq 'svg') set term svg dynamic font 'Kongtext,4'; set output ARG7
if (ARG2 eq 'svgmono') set term svg dynamic font 'Kongtext,4'; set monochrome; set output ARG7
if (ARG2 eq 'dxf') set term dxf; set output ARG7 
#set mxtics 5
#set mytics 5
# set grid xtics ytics
set title ARG3
set xlabel ARG4
set ylabel ARG5
# set autoscale
# unset border
unset key
unset colorbox
set offsets graph 0.1,0.1,0.1,0.1
plot ARG1 index 0 using 1:2:3 w points pt 1,\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead lc var,\
     ARG1 index 2 using 1:2:3 w labels