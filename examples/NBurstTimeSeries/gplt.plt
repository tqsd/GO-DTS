
set term x11
set key autotitle columnhead



set key outside
unset key
set ylabel "y-units" font "Times-Roman,8"
set key font ",8"
set xrange [999000:1000000]

#set xrange [0:2000]
unset border


set multiplot layout 5,1\
              margins 0.1,0.98,0.1,0.98 \
              spacing 0.005,0.005\
              title "NBURST"


set linetype 1 lc rgb "#2596be" lw 0.5 # GEWI PLOTS
set linetype 2 lc rgb "#063970" lw 0.5 # CLASSIC PLOTS
set linetype 3 lc rgb "#21130d" lw 0.5 # TRAFFIC
set linetype 4 lc rgb "#873e23" lw 2  # E-BUFFER
set linetype 5 lc rgb "#bb000000" lw 1  # E-BUFFER, TRAFFIC AVERAGE


#### THE PLOT WITH TRAFFIC
unset title
set format x ''
set ylabel "Incoming"
plot "results/results.data" u 1:2 with lines lt 3
### THE PLOT WITH DROPS
set ylabel "Droped"
plot "results/results.data" u 1:8 with lines lt 2
plot "results/results.data" u 1:4 with lines lt 1

set ylabel "Transmission rate"
plot "results/results.data" u 1:3 with lines lt 1 ,\
     "results/results.data" u 1:7 with lines lt 2

set xtics font ", 5"
set ylabel "Buffer state"
set xlabel "Timestep"
set format "%.f"
set xtics
plot "results/results.data" u 1:6 with lines lt 1,\
     "results/results.data" u 1:9 with lines lt 2,\
     "results/results.data" u 1:5 with lines linetype 5,\
     "results/results.data" u 1:5 with lines linetype 4 dashtype 2,\

unset multiplot
