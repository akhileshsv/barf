if (ARG2 eq 'dumb') set term dumb ansi enhanced; else set term ARG2 persist 
set tics textcolor rgb "green"
set border lw 1 lc "magenta"
set autoscale
set offsets graph 0.1,0.1,0.1,0.1
unset key
plot ARG1 index 0 using 1:2:3 w lines lc var, ARG1 index 1 using 1:2:3 with labels
exit
