#Ploting Mean Waiting Times
set term wxt

set key autotitle columnhead
set datafile separator ','


plot "results/results.csv" using "PARETO-NodeCount":(column("LINK-MWT")-column("GEWI-MWT"))
