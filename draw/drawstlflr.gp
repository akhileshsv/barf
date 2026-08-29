if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,49
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'svg') set term svg dynamic font 'Kongtext,2'; set output ARG7
if (ARG2 eq 'qt') set term qt enhanced font 'Courier,4' persist
if (ARG2 eq 'svgmono') set term svg size 1000,1000 font 'Kongtext,2'; set monochrome; set output ARG7
if (ARG2 eq 'dxf') set term dxf; set output ARG7

# set path of config snippets
#set loadpath './gnuplot-palettes'
#load 'noborder.cfg'
#load 'blues.pal'
set palette rgbformulae 7,5,15
set mxtics 5
set mytics 5
unset grid
set title ARG3
set xlabel ARG4
set ylabel ARG5
#set ytics autofreq mirror scale 0.25 tc lt 1
#set xtics autofreq nomirror scale 0.25 rotate tc lt 1
unset key
unset colorbox
set size ratio -1
unset border
#set size ratio -1
set offsets graph 0.1,0.1,0.1,0.1
if (ARG2 eq 'dxf') {
plot ARG1 index 0 using 1:2:3 w points pt ".",\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead,\
     ARG1 index 2 using 1:2:3:4 w vectors nohead

} else {
plot ARG1 index 0 using 1:2:3 w points pt ".",\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead lw 0.3 lc var dashtype 1,\
     ARG1 index 2 using 1:2:3:4 w vectors nohead lw 0.2 dashtype 1,\
     ARG1 index 3 using 1:2:3 w labels offset char -1,-1 notitle
}
