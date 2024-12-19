set print "-"
set term ARG2 persist font "Courier,4"
set tics textcolor rgb "green"
#set border lw 1 lc "magenta"
#set palette cubehelix
set noborder
set tics textcolor rgb "magenta" rotate
#set view 
set ticslevel 0
#set grid lw 1 lc "light-gray"
set noborder
set autoscale
#set hidden3d
#set offsets graph 0.1,0.1,0.1,0.1,0.1,0.1
set title "FRAME" tc rgb "blue" font "Courier"
unset key
set mouse
splot ARG1 index 0 using 1:2:3 w points,\
      ARG1 index 1 using 1:2:3 w points pointtype 5,\
      ARG1 index 2 using 1:2:3:($4-$1):($5-$2):($6-$3):7 w vectors nohead lc var lw 3 notitle,\
      ARG1 index 3 using 1:2:3:4 w polygons fs transparent solid 0.1,\
      ARG1 index 4 using 1:2:3:4 w labels offset 7 char
#print("NOBLE GNUPLOT will now save svg")
#set view projection xz
#replot
#set view projection yz
#replot
#set view projection xy
#replot

#unset mouse
#unset
#set term png truecolor transparent crop 
#set output ARG3
#replot
if (ARG2 eq 'qt') pause mouse close
#exit
