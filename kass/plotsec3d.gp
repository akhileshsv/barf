if (ARG2 eq 'dumb') set term dumb ansi enhanced size 79,49
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'mono') set term dumb mono size 79,49
if (ARG2 eq 'qt') set term qt persist
if (ARG2 eq 'wxt') set term wxt persist
if (ARG2 eq 'svg') set term svg; set output ARG4
set tics textcolor rgb "magenta"
#set view equal xyz
set noborder
set autoscale
set key outside right bottom
set ticslevel 0
set xlabel "X mm" offset 1,1,1 tc rgb 'green' rotate by 45
set ylabel "Y mm" offset 1,1,1 tc rgb 'green' rotate by 45
set zlabel "Z mm" offset 1,1,1 tc rgb 'green' rotate by 45
set hidden3d offset 2
set title ARG3
set offsets graph 0.1,0.1,0.1,0.1
splot ARG1 index 0 using 1:2:3 w lines notitle

