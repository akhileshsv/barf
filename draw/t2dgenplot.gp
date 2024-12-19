if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,49; else set term ARG2 persist font "Courier,4"
if (ARG2 eq 'qt') set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb "black" behind
set tics textcolor rgb "green" 
#set loadpath 'C:\Users\Admin\junk\barf\draw\config'
#load 'moreland.pal'
#set border lw 1 lc "magenta"
set noborder
if (ARG2 eq 'qt') set grid lw 1 lc "blue"
set autoscale
set offsets graph 0.1,0.1,0.1,0.1
set title ARG3 tc rgb "grey"
unset key
if (ARG2 eq 'qt'){
plot ARG1 index 0 using 1:2:3 w labels tc rgb "blue" point pt 7 offset char -1,-1 notitle,\
     ARG1 index 1 using 1:2:($3-$1):($4-$2):6 notitle w vectors nohead lc var lw 2,\
     ARG1 index 1 using ($3+$1)/2:($2+$4)/2:5 w labels tc rgb "white" notitle
} else {
plot ARG1 index 0 using 1:2 notitle,\
     ARG1 index 1 using 1:2:($3-$1):($4-$2):6 notitle w vectors nohead lc var,\
}
#if (ARG2 eq 'qt') pause mouse close
#exit
