if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,79
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'qt') set term qt enhanced font 'Kongtext,3' persist
if (ARG2 eq 'mono') set term dumb mono size 99,79
if (ARG2 eq 'dxf') set term dxf; set output ARG4
if (ARG2 eq 'svg') set term svg dynamic font "Kongtext,4"; set output ARG4
if (ARG2 eq 'svgmono') set term svg dynamic font "Kongtext,4"; set output ARG4; set monochrome
# set size ratio -1
set noborder
set autoscale
#unset key
# set ytics autofreq nomirror scale 0.25 tc lt 1
# set xtics autofreq nomirror scale 0.25 rotate tc lt 1
#set title ARG3 tc rgb "magenta" font "Kongtext,8"

set key bottom textcolor rgb "blue"
unset tics
unset colorbox
#set size ratio -1
set xlabel 'm' 
set ylabel 'm' 
set offsets graph 0.1,0.1,0.1,0.1
set style fill transparent solid 0.5 noborder
set multiplot layout 3,1
#set multiplot layout 2,1
set title sprintf('%s shear', ARG3) font "Kongtext,4"
plot ARG1 index 0 using 1:2:($3-$1):($4-$2):5 w vectors nohead lw 1 lc var title "mems",\
     ARG1 index 1 using 1:2 w lines dt 1 title "sf",\
     ARG1 index 1 using 7:8:($1-$7):($2-$8):5 w vectors nohead dt "." notitle,\
     ARG1 index 2 using 1:2:3 w labels offset char 1,1 notitle
set title sprintf('%s bending moment', ARG3) font "Kongtext,4"
plot ARG1 index 0 using 1:2:($3-$1):($4-$2):5 w vectors nohead lw 1 lc var title "mems",\
     ARG1 index 1 using 3:4 w lines dt 1 title "bm",\
     ARG1 index 1 using 7:8:($3-$7):($4-$8):5 w vectors nohead dt "." notitle,\
     ARG1 index 2 using 1:2:4 w labels offset char 1,1 notitle
     #ARG1 index 2 using 10:11:12 w labels tc rgb "blue" notitle
set title sprintf('%s deflection', ARG3) font "Kongtext,4"
plot ARG1 index 0 using 1:2:($3-$1):($4-$2):5 w vectors nohead lw 1 lc var title "mems",\
     ARG1 index 1 using 5:6 w lines dt 1 title "dx",\
     ARG1 index 1 using 7:8:($5-$7):($6-$8):5 w vectors nohead dt "." notitle,\
     # ARG1 index 3 using 1:2:3:4 w lines title "memdispl",\
     ARG1 index 2 using 1:2:5 w labels offset char 1,1 notitle
#ARG1 index 2 using 1:2:($3-$1):($4-$2) w vectors nohead lw 1 dt 1 title "def. shape"
#ENABLE THIS FOR INTERACTIVE PLOTS
unset multiplot
