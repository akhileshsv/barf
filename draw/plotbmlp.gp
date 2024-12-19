if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,79
if (ARG2 eq 'caca') set term caca driver utf8 enhanced inverted size 99,79  
if (ARG2 eq 'qt') set term qt enhanced persist font "Courier,8"#; set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb "black" behind
set tics textcolor rgb "green"
#set size ratio 0.5
#set noborder
set autoscale
set multiplot layout 2,2 rowsfirst
#unset key
set xlabel 'mm' tc rgb "green"
set ylabel 'KN' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
unset xtics
#set xtics autofreq nomirror scale 0.25 rotate tc lt 1
#set lmargin 1
#set rmargin 1
#set tmargin 1
#set bmargin 1
set key bottom
set offsets graph 0.05,0.05,0.05,0.05
#set title 'shear envelope'
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
plot ARG1 index 1 using 1:2:3:4:5 w boxxyerrorbars lc var notitle,\
     ARG1 index 2 using 1:2:3 w labels tc rgb "white" notitle,\
     ARG1 index 3 using 1:2 w points pt 8 ps 2 notitle,\
     ARG1 index 4 using 1:2:3:4:5 w vectors nohead lc var title "dead load"
#set title 'hogging moment envelope'
#set datafile missing '0.0'
set offsets graph 0.05,0.05,0.05,0.05
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
plot ARG1 index 1 using 1:2:3:4:5 w boxxyerrorbars lc var notitle,\
     ARG1 index 2 using 1:2:3 w labels tc rgb "white" notitle,\
     ARG1 index 3 using 1:2 w points pt 8 ps 2 notitle,\
     ARG1 index 5 using 1:2:3:4:5 w vectors nohead lc var title "live load"
set offsets graph 0.05,0.05,0.05,0.05
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
plot ARG1 index 0 using 1:2 w lines title "shear force",\
     ARG1 index 1 using 1:2:3:4:5 w boxxyerrorbars lc var notitle,\
     ARG1 index 2 using 1:2:3 w labels tc rgb "white" notitle,\
     ARG1 index 3 using 1:2 w points pt 8 ps 2 notitle
#set title 'sagging moment envelope'
#set datafile missing '0.0'
set offsets graph 0.05,0.05,0.05,0.05
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
plot ARG1 index 0 using 1:3 w lines title "bending moment",\
     ARG1 index 1 using 1:2:3:4:5 w boxxyerrorbars lc var notitle,\
     ARG1 index 2 using 1:2:3 w labels tc rgb "white"  notitle,\
     ARG1 index 3 using 1:2 w points pt 8 ps 2 notitle

unset multiplot
