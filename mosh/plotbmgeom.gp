if (ARG2 eq 'dumb') set term dumb ansi enhanced size 79,49
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'mono') set term dumb mono size 79,49
if (ARG2 eq 'qt') set term qt persist
if (ARG2 eq 'wxt') set term wxt persist
if (ARG2 eq 'svg') set term svg dynamic font 'Kongtext,5'; set output ARG4
if (ARG2 eq 'qtmono') set term qt persist font 'Kongtext,5'; set monochrome
if (ARG2 eq 'svgmono') set term svg dynamic font 'Kongtext,5'; set output ARG4; set monochrome
if (ARG2 eq 'dxf') set term dxf; set output ARG4
set tics textcolor rgb "magenta"
set size ratio -1
set noborder
# set autoscale
set key outside right bottom
set xlabel 'mm' tc rgb "green"
set ylabel 'mm' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
set title ARG3
set offsets graph 0.1,0.1,0.1,0.1
set mouse
if (ARG2 eq "qt") set size ratio -1; set mouse #;set grid xtics ytics mxtics mytics
if (ARG2 eq "dxf") unset xtics; unset ytics; set font "=0.01 m"; set size ratio -1
if (ARG2 eq "svg") || (ARG2 eq 'svgmono') unset grid
if (ARG2 eq 'dumb') || (ARG2 eq 'mono'){
plot ARG1 index 0 using 1:2:3 w lines lw 2 lc var notitle,\
     ARG1 index 1 using 1:2:3 w labels offset char 1,1 tc rgb "green" notitle,\
     ARG1 index 2 using 1:2 w points pt 15 title "bars",\
     ARG1 index 2 using 1:2:(sprintf('%.f', $3*2.0)) w labels offset char 1,1 tc rgb "red" notitle,\
     ARG1 index 3 using 1:2 w lines dt 0 title "n.a"
} else {
plot ARG1 index 0 using 1:2:3 w lines lw 2 lc var notitle,\
     ARG1 index 1 using 1:2:3 w labels offset char 1,1 tc rgb "blue" notitle,\
     ARG1 index 2 using 1:2:($3 * 2):($3 * 2) w ellipses title "bars",\
     ARG1 index 2 using 1:2:(sprintf('%.f', $3*2.0)) w labels offset char 1,1 tc rgb "red" notitle,\
     ARG1 index 3 using 1:2 w lines dt 0 title "n.a"
}
if (ARG2 eq 'qt') pause mouse close

