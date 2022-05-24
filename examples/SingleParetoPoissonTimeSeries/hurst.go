// Simple program to inspect the behavior of the link
package main

import (
	"fmt"
	"github.com/tqsd/dts/analysis"
	"github.com/tqsd/dts/traffic"
)

func main() {

	// Creating Traffic source
	//
	scale := float64(10)
	shape := float64(1.991)
	on_prob := 0.5
	traffAvg := float64(100)
	source := traffic.NewSingleParetoPoissonSource(scale, shape, shape, scale, on_prob, traffAvg)

	simulation_steps := 10000

	traffic := make([]uint32, 0)
	for s := 0; s < simulation_steps; s++ {
		traffic = append(traffic, uint32(source.Tick()))
	}
	fmt.Println("Analyzing the hurst parameter")
	hurst, err := analysis.Hurst(traffic)
	if err != nil {
		fmt.Println("Some error", err)
	}
	fmt.Println("Alpha", scale)
	fmt.Println("Hurst Parameter", hurst)

}
