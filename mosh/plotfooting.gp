if (ARG2 eq 'dumb') set term dumb ansi enhanced size 149,49
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'qt') set term qt persist font "Courier,12"#set object 1 rectangle from screen 0,0 to screen 1,1 fillcolr rgb "black" behind
if (ARG2 eq 'svg') set term svg font "Kongtext,5"; set output ARG4 
set tics textcolor rgb "magenta" rotate
set ticslevel 0
set xlabel "X m" offset 1,1,1 tc rgb 'green' rotate by 45
set ylabel "Y m" offset 1,1,1 tc rgb 'green' rotate by 45
set zlabel "Z m" offset 1,1,1 tc rgb 'green' rotate by 45
#set multiplot layout 2,2 title 'rcc footing' 
#set hidden3d
set noborder
set autoscale
if (ARG2 eq 'qt') set mouse

unset key
if (ARG2 eq 'dumb')||(ARG2 eq 'mono') set view projection xy
splot ARG1 index 0 w lines dashtype 2 lt 2,\
      ARG1 index 1 w vectors nohead lc var lw 2
#set view projection xz
#replot
#set view projection yz
#replot
#set view projection xy
#replot
if (ARG2 eq 'qt') pause mouse close
exit
