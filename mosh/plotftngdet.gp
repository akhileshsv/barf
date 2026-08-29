if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,99
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'mono') set term dumb mono size 99,99
if (ARG2 eq 'qt') set term qt persist font 'Courier,8' size 800,800
if (ARG2 eq 'wxt') set term wxt persist
if (ARG2 eq 'svg') set term svg font 'Kongtext,5'; set output ARG4
if (ARG2 eq 'qtmono') set term qt persist font 'Courier,8'; set monochrome
if (ARG2 eq 'svgmono') set term svg font 'Kongtext,5'; set output ARG4; set monochrome
if (ARG2 eq 'dxf') set term dxf; set output ARG4; set size 1000,1000
set tics textcolor rgb "magenta"
set size ratio -1
set noborder
#set autoscale
set key outside right bottom
set xlabel 'm' tc rgb "green"
set ylabel 'm' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
set mxtics 5
set mytics 5
set title ARG3
set multiplot layout 2,1 rows
set offsets graph 0.1,0.1,0.1,0.1
if (ARG2 eq "qt") set size ratio -1; set mouse #;set grid xtics ytics mxtics mytics
if (ARG2 eq "dxf") unset xtics; unset ytics; set font "=0.01 m"; set size ratio -1
if (ARG2 eq "svg") || (ARG2 eq 'svgmono') unset grid;unset xtics; unset ytics
set mouse
if (ARG2 eq 'svg') || (ARG2 eq 'svgmono'){
plot ARG1 index 0 using 1:2:3:4:5 w vectors nohead lc var lw 0.25 title "concrete",\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead lc var lw 0.25 title "rebar 1",\
     ARG1 index 2 using 1:2:3:4:5 w vectors heads lc var lw 0.2 title "dims(m)",\
     ARG1 index 3 using 1:2:3:4 w labels offset char 1,1 tc var notitle,\
     ARG1 index 4 using 1:2:3:4 w ellipses lw 0.5 title "rebar 2"
plot ARG1 index 5 using 1:2:3:4:5 w vectors nohead lc var lw 0.25 title "concrete",\
     ARG1 index 6 using 1:2:3:4:5 w vectors nohead lc var lw 0.5 title "rebar 1",\
     ARG1 index 7 using 1:2:3:4:5 w vectors heads lc var lw 0.2 title "dims(m)",\
     ARG1 index 8 using 1:2:3:4 w labels offset char 1,1 tc var notitle
} else {
plot ARG1 index 0 using 1:2:3:4:5 w vectors nohead lc var lw 1 title "concrete",\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead lc var lw 1 title "rebar 1",\
     ARG1 index 2 using 1:2:3:4:5 w vectors heads lc var lw 1 title "dims(m)",\
     ARG1 index 3 using 1:2:3:4 w labels offset char 1,1 tc var notitle,\
     ARG1 index 4 using 1:2:3:4 w ellipses lw 1 title "rebar 2"
plot ARG1 index 5 using 1:2:3:4:5 w vectors nohead lc var lw 1 title "concrete",\
     ARG1 index 6 using 1:2:3:4:5 w vectors nohead lc var lw 2.0 dt 2.0 title "rebar 1",\
     ARG1 index 7 using 1:2:3:4:5 w vectors heads lc var lw 1 title "dims(m)",\
     ARG1 index 8 using 1:2:3:4 w labels offset char 1,1 tc var notitle
}
     #ARG1 index 2 using 1:2:(sprintf('%.f', $3*2.0)) w labels offset char 1,1 tc rgb "red" notitle,\
     #ARG1 index 3 using 1:2 w lines dt 0 title "n.a"
if (ARG2 eq 'qt') || (ARG2 eq 'qtmono') pause mouse close

