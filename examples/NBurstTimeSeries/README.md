# TIME SERIES ANALYSIS

In this example we analize a time series created with the _Pareto based_ traffic generator, which is processed by a gewi-based and classical link.

The parameters need to be adjusted in the file ```main.go```.
The simulation can be run by
```
go run main.go results/results.data
```

Then the plots can be generated with gnuplot
```
gnuplot --persist gplt.plt
```
