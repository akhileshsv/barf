if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,79
if (ARG2 eq 'caca') set term caca driver utf8 inverted enhanced size 99,79 
if (ARG2 eq 'qt') set term qt persist font "Courier,4"
if (ARG2 eq 'svg') set term svg font "Kongtext,5"; set output ARG4
if (ARG2 eq 'svgmono') set term svg font "Kongtext,5"; set output ARG4; set monochrome
if (ARG2 eq 'dxf') set term dxf; set output ARG4

#set multiplot
set tics textcolor rgb "magenta"
set noborder
set autoscale
set key bottom 
set xlabel 'meters'
set ylabel 'meters'
#set y2label 'kn/m'
set ytics autofreq nomirror tc lt 1
#set y2tics autofreq nomirror tc lt 2
set title 'sub frame'
set mxtics 5
set mytics 5
if (ARG2 eq "qt" || ARG2 eq "svg") set grid xtics ytics 
set offsets graph 0.1,0.1,0.1,0.1
plot ARG1 index 0 using 1:2:3 w labels point pt 7 offset char 1,1 title "nodes",\
     ARG1 index 1 using 1:2:($3-$1):($4-$2) w vectors lw 2 nohead title "members", \
     ARG1 index 1 using ($3+$1)/2:($2+$4)/2:5 w labels offset char -1,-1 notitle,\
     ARG1 index 1 using ($3+$1)/2:($2+$4)/2:6 w labels offset char 1,1 notitle,\
     ARG1 index 2 using 1:2:3 w points pt variable title "supports",\
     ARG1 index 3 using 1:2:($3-$1):($4-$2):7 w vectors lc variable notitle,\
     ARG1 index 3 using ($3+$1)/2:($4+$2)/2:5 w labels offset char 2,2 title "load",\
     #ARG1 index 1 using 1:2:($3-$1):($4-$2):6 notitle w vectors nohead lc var,\
     #ARG1 index 1 using ($3+$1)/2:($2+$4)/2:5 w labels offset char -1,-1 tc rgb "green" notitle
     #ARG1 index 3 using 1:2:3 w 
#set loadpath 'C:\\Users\\Akhilesh SV\\Desktop\\go\\projex\\r2d3\\data\\gnuplot'
#load 'moreland.pal'
#if (ARG2 eq 'dumb') set term dumb ansi enhanced size 79,39; else set term ARG2 persist
#if (ARG2 eq 'qt') set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb "black" behind

#set border lw 1 lc "magenta"
#set noborder
#set grid lw 1 lc "blue"
#set autoscale
#set offsets graph 0.1,0.1,0.1,0.1
#unset key
#plot ARG1 index 0 using 1:2:3 w labels tc rgb "blue" point pt 7 offset char -1,-1 notitle,\
#if (ARG2 eq 'wxt') pause mouse close
#exit
