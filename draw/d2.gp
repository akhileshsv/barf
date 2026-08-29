if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,49
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'qt') set term qt enhanced font 'Courier,4' persist
if (ARG2 eq 'svg') set term svg dynamic font 'Kongtext,4'; set output ARG7
if (ARG2 eq 'svgmono') set term svg dynamic font 'Kongtext,4'; set monochrome; set output ARG7
if (ARG2 eq 'dxf') set term dxf; set size 1000,1000; set output ARG7; set size ratio -1 

# set path of config snippets
#set loadpath './gnuplot-palettes'
#load 'noborder.cfg'
#load 'blues.pal'
set mxtics 5
set mytics 5
set grid xtics ytics
set title ARG3
set xlabel ARG4
set ylabel ARG5
#set ytics autofreq mirror scale 0.25 tc lt 1
#set xtics autofreq nomirror scale 0.25 rotate tc lt 1
unset key
unset colorbox
set autoscale
unset border
#set size ratio -1
set offsets graph 0.1,0.1,0.1,0.1
plot ARG1 index 0 using 1:2 w lines lw 2,\
     ARG1 index 1 using 1:2:3 w labels,\
     ARG1 index 2 using 1:2:3 w points,\