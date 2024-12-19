#set loadpath 'C:\\Users\\Akhilesh SV\\Desktop\\go\\projex\\r2d3\\data\\gnuplot'
#load 'moreland.pal'
if (ARG2 eq 'dumb') set term dumb ansi enhanced size 99,79; else set term ARG2 persist
set tics textcolor rgb "green"
set colorsequence default
set border lw 1 lc "magenta"
#set multiplot layout 2,1
set autoscale
set offsets graph 0.1,0.1,0.1,0.1
unset key
unset colorbox
#plot ARG1 index 0 using 1:2:3:4 w lines lc var,\
#     ARG1 index 3 using 1:2:3 with labels tc rgb "dark-cyan"
plot ARG1 index 1 using 1:2:3:4 w lines lc var,\
     ARG1 index 2 using 1:2:3 with labels tc rgb "aquamarine"
if (ARG2 eq 'wxt') pause mouse close
exit
