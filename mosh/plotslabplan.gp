if (ARG2 eq 'dumb') set term dumb ansi enhanced size 79,79
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'mono') set term dumb mono size 79,79
if (ARG2 eq 'qt') set term qt persist font "Courier,5"; set mouse
if (ARG2 eq 'wxt') set term wxt persist
if (ARG2 eq 'svg') set term svg font "Kongtext,5"; set output ARG4
if (ARG2 eq 'dxf') set term dxf; set size 1000,100; set output ARG4
set tics textcolor rgb "magenta"
set noborder
set size ratio -1
set key outside right bottom
set xlabel 'mm' tc rgb "green"
set ylabel 'mm' tc rgb "green"
set ytics autofreq nomirror scale 0.25 tc lt 1
set xtics autofreq nomirror scale 0.25 rotate tc lt 1
set title ARG3
set offsets graph 0.1,0.1,0.1,0.1
set multiplot layout 1,2 rowsfirst

plot ARG1 index 0 using 1:2:3:4:5 w vectors nohead lc var title "concrete",\
     ARG1 index 1 using 1:2:3:4:5 w vectors nohead lc var dt 3 title "support",\
     ARG1 index 2 using 1:2:3:4:5 w vectors nohead lc var dt 2 title "rebar",\
     ARG1 index 3 using 1:2:3:4:5 w vectors heads notitle,\
     ARG1 index 4 using 1:2:3 w labels offset char 1,1 notitle
plot ARG1 index 0 using 1:2:3:4:5 w vectors nohead lc var title "concrete",\
     ARG1 index 0 using 6:7:8:9:10 w vectors nohead lc var dt 3 title "support",\
     ARG1 index 5 using 1:2:3:4:5 w vectors nohead lc var dt 2 title "rebar",\
     ARG1 index 4 using 1:2:3 w labels offset char 1,1 notitle,\
     ARG1 index 3 using 1:2:3:4:5 w vectors heads notitle,\
     ARG1 index 4 using 1:2:3:4:5 w vectors heads notitle,\
     ARG1 index 5 using 1:2:3 w labels offset char 1,1 notitle

# plot ARG1 index 0 using 1:2:3:4:5 w vectors nohead lc var title "concrete",\
#      ARG1 index  using 6:7:8:9:10 w vectors nohead lc var dt 3 title "support",\
#      ARG1 index 1 using 1:2:3:4:5 w vectors nohead lc var dt 2 title "rebar",\
#      ARG1 index 2 using 1:2:3:4:5 w vectors heads notitle,\
#      ARG1 index 3 using 1:2:3 w labels offset char 1,1 notitle
# plot ARG1 index 0 using 1:2:3:4:5 w vectors nohead lc var title "concrete",\
#      ARG1 index 0 using 6:7:8:9:10 w vectors nohead lc var dt 3 title "support",\
#      ARG1 index 4 using 1:2:3:4:5 w vectors nohead lc var dt 2 title "rebar",\
#      ARG1 index 3 using 1:2:3 w labels offset char 1,1 notitle,\
#      ARG1 index 2 using 1:2:3:4:5 w vectors heads notitle,\
#      ARG1 index 5 using 1:2:3:4:5 w vectors heads notitle,\
#      ARG1 index 6 using 1:2:3 w labels offset char 1,1 notitle
