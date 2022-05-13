#Ploting Mean Waiting Times
set term wxt

set key autotitle columnhead
set datafile separator ','
filename = "results/2.csv"

#column("Pareto-Scale_on")<column("Pareto-Scale_off")
#plot "results/results.csv" using "PARETO-NodeCount":( column("GEWI-Mult")==2 ? (column("LINK-MWT")-column("GEWI-MWT"):1/0) : 1/0)
#plot "results/results.csv" using "PARETO-NodeCount":( column("Pareto-Scale_on") == 1.2 ? column("LINK-MWT")-column("GEWI-MWT") : 1/0)
#PARETO-Scale_on
#



set multiplot layout 2,3\
              margins 0.15,0.8,0.15,0.8 \
              spacing 0.08,0.08

set style circle radius 0.5
set format x ""
set yrange [0:300]
set logscale y 10
plot filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 128 && column("PARETO-Shape_on")<column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) with circles,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 16 && column("PARETO-Shape_on")<column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) with circles,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 2 && column("PARETO-Shape_on")<column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) with circles,\

plot filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 128 && column("PARETO-Shape_on")==column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) with circles,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 16 && column("PARETO-Shape_on")==column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) with circles,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 2 && column("PARETO-Shape_on")==column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) with circles,\

plot filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 128 && column("PARETO-Shape_on")>column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) with circles,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 16 && column("PARETO-Shape_on")>column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) with circles,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 2 && column("PARETO-Shape_on")>column("PARETO-Shape_off") ) ? column("LINK-MWT")-column("GEWI-MWT") : 1/0) with circles

unset logscale y


set format x "%.f"
set yrange [0:0.12]
set xtics 50
plot filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 128 && column("PARETO-Shape_on")<column("PARETO-Shape_off") ) ?\
     column("LINK-Droppd")/column("LINK-Recv") - column("GEWI-Droppd")/column("GEWI-Recv"): 1/0) with circles ,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 16 && column("PARETO-Shape_on")<column("PARETO-Shape_off") ) ?\
     column("LINK-Droppd")/column("LINK-Recv") - column("GEWI-Droppd")/column("GEWI-Recv"): 1/0) with circles,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 2 && column("PARETO-Shape_on")<column("PARETO-Shape_off") ) ?\
     column("LINK-Droppd")/column("LINK-Recv") - column("GEWI-Droppd")/column("GEWI-Recv"): 1/0) with circles

plot filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 128 && column("PARETO-Shape_on")==column("PARETO-Shape_off") ) ?\
     column("LINK-Droppd")/column("LINK-Recv") - column("GEWI-Droppd")/column("GEWI-Recv"): 1/0) with circles,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 16 && column("PARETO-Shape_on")==column("PARETO-Shape_off") ) ?\
     column("LINK-Droppd")/column("LINK-Recv") - column("GEWI-Droppd")/column("GEWI-Recv"): 1/0) with circles,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 2 && column("PARETO-Shape_on")==column("PARETO-Shape_off") ) ?\
     column("LINK-Droppd")/column("LINK-Recv") - column("GEWI-Droppd")/column("GEWI-Recv"): 1/0) with circles
set key
plot filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 128 && column("PARETO-Shape_on")>column("PARETO-Shape_off") ) ?\
     column("LINK-Droppd")/column("LINK-Recv") - column("GEWI-Droppd")/column("GEWI-Recv"): 1/0) with circles,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 16 && column("PARETO-Shape_on")>column("PARETO-Shape_off") ) ?\
     column("LINK-Droppd")/column("LINK-Recv") - column("GEWI-Droppd")/column("GEWI-Recv"): 1/0) with circles,\
     filename using "PARETO-NodeCount":( (column("GEWI-Mult") == 2 && column("PARETO-Shape_on")>column("PARETO-Shape_off") ) ?\
     column("LINK-Droppd")/column("LINK-Recv") - column("GEWI-Droppd")/column("GEWI-Recv"): 1/0) with circles
