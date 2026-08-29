#!/usr/bin/gnuplot
if (ARG2 eq 'dumb') set term dumb ansi enhanced size 79,49
if (ARG2 eq 'caca') set term caca driver utf8 color inverted
if (ARG2 eq 'qt') set term qt font 'Courier,8' persist
if (ARG2 eq 'svg') set term svg dynamic font 'Kongtext,4'; set output ARG7
if (ARG2 eq 'svgmono') set term svg dynamic font 'Kongtext,4'; set monochrome; set output ARG7

set title ARG3
set xlabel ARG4 tc rgb "green"
set ylabel ARG5 tc rgb "green"
set zlabel ARG6 tc rgb "green"
#unset tics
set ticslevel 0
unset key
unset colorbox
set autoscale
splot ARG1 index 0 u 1:2:3 w points,\
      ARG1 index 0 u 1:2:3:4 w labels,\
      ARG1 index 1 u 1:2:3:4:5:6:7 w vectors nohead lc var lw 2
if (ARG2 eq 'qt') set mouse; pause mouse close
