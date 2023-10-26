if (ARG2 eq 'dumb') set term dumb ansi enhanced size 79,79
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'mono') set term dumb mono size 79,79
if (ARG2 eq 'qt') set term qt persist font "Courier,5"
if (ARG2 eq 'wxt') set term wxt persist
if (ARG2 eq 'svg') set term svg font "Courier,5"
set tics textcolor rgb "magenta"
set noborder
set autoscale
set key outside right bottom
set xlabel 'mm' tc rgb "green"
set ylabel 'mm' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
set title ARG3
set offsets graph 0.1,0.1,0.1,0.1
if (ARG2 eq 'svg') set output ARG4
if (ARG2 eq 'dumb') || (ARG2 eq 'mono'){
set size ratio 0.75
#set multiplot layout 2,1 rowsfirst
plot ARG1 index 0 using 1:2:3 w lines lw 2 lc var dt 1 notitle,\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead lc var title "main bar",\
     ARG1 index 2 using 1:2:3 w labels notitle,\
     ARG1 index 3 using 1:2 w points pt 0 title "dist bar"
#set title "plan"
#plot ARG1 index 4 using 1:2:3 w lines lc var dt 0 notitle,\
#     ARG1 index 5 using 1:2:3:4:5 w vectors heads lc var title "rebar"
} else {
#set multiplot layout 2,1 rowsfirst
#set title "section"
plot ARG1 index 0 using 1:2:3 w lines lw 2 lc var notitle,\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead lc var title "main bar",\
     ARG1 index 2 using 1:2:3 w labels notitle,\
     ARG1 index 3 using 1:2:3 w circles title "dist bar"
#set title "plan"
#plot ARG1 index 4 using 1:2:3 w lines lc var dt 0 notitle,\
#     ARG1 index 5 using 1:2:3:4:5 w vectors heads lc var title "rebar"
}

