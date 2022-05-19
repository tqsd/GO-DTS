set term x11
#set terminal png size 1000,400 enhanced font ",20"
#set output 'output.png'

#set term cairolatex pdf
#set output "mwt_pareto.tex"


#set terminal push
#set terminal lua tikz fulldoc createstyle
#set output 'sin.tex'

set key autotitle columnhead
set datafile separator ','
filename = "results/2.csv"


#set title "Mean Wait Time Difference"



set style circle radius 0.5
set format x ""
set yrange [-10:300]
set logscale y 10

set border back

set ylabel "$\\log(d^{ (C)}_W -  d^{(G)}_W)$" font ",12" offset 2
set xlabel "N" font ",12" offset 0,0.7
set xtics font ",5"
set xtics 0,50,200
set xtics font ",10"
set format "%.f"
set xtics
set ytics font ",10"
#set size sfillcolor "0x7faaaaaa"quare

set multiplot layout 1,3\
              margins 0.1,0.9,0.2,0.8 \
              spacing 0.02,0.1
unset key
set title "$\\alpha_{ON} < \\alpha_{OFF}$"
plot filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 128 && column("PARETO-Shape_on")<column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) w p pt 6 ps 0.3 lc rgb "#11ff0000" ,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 16 && column("PARETO-Shape_on")<column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) w p pt 6 ps 0.3 lc rgb "#1100ff00",\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 2 && column("PARETO-Shape_on")<column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0)  w p pt 6 ps 0.3 lc rgb "#110000ff",\
     NaN w p pt 6 ps 1 lc "red" title "Mult=128" , NaN w p pt 6 ps 1 lc "green" title "Mult=16" , NaN w p pt 6 ps 1 lc  "blue" title "Mult=2"



unset ylabel
#unset ytics
unset key
set format y "";

set title "$\\alpha_{ON} = \\alpha_{OFF}$"
plot filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 128 && column("PARETO-Shape_on")==column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) w p pt 6 ps 0.3 lc rgb "#11ff0000" ,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 16 && column("PARETO-Shape_on")==column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) w p pt 6 ps 0.3 lc rgb "#1100ff00",\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 2 && column("PARETO-Shape_on")==column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0)  w p pt 6 ps 0.3 lc rgb "#110000ff",\
     NaN w p pt 6 ps 1 lc "red" title "Mult=128" , NaN w p pt 6 ps 1 lc "green" title "Mult=16" , NaN w p pt 6 ps 1 lc  "blue" title "Mult=2"

set key at 195,250
set key box opaque fillcolor "0x33ffffff"
set key samplen 2 spacing 2 font ",10"


set title "$\\alpha_{ON} > \\alpha_{OFF}$"
plot filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 128 && column("PARETO-Shape_on")>column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) w p pt 6 ps 0.3 lc rgb "#11ff0000" ,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 16 && column("PARETO-Shape_on")>column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) w p pt 6 ps 0.3 lc rgb "#1100ff00",\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 2 && column("PARETO-Shape_on")>column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0)  w p pt 6 ps 0.3 lc rgb "#110000ff",\
     NaN w p pt 6 ps 1 lc "red" title "Mult=128" , NaN w p pt 6 ps 1 lc "green" title "Mult=16" , NaN w p pt 6 ps 1 lc  "blue" title "Mult=2"

unset multiplot
