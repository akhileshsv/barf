if (ARG2 eq 'dumb') set term dumb ansi enhanced size 79,39; else set term ARG2 persist
if (ARG2 eq 'qt') set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb "black" behind
set tics textcolor rgb "green"
#set border lw 1 lc "magenta"
set noborder
#set grid lw 1 lc "blue"
set autoscale
set offsets graph 0.1,0.1,0.1,0.1
unset key
plot ARG1 index 0 using 1:2:3 w labels tc rgb "blue" point pt 7 offset char -1,-1 notitle,\
     ARG1 index 1 using 1:2:($3-$1):($4-$2):6 notitle w vectors nohead lc var,\
     ARG1 index 1 using ($3+$1)/2:($2+$4)/2:5 w labels tc rgb "white" notitle
#if (ARG2 eq 'wxt') pause mouse close
#exit
