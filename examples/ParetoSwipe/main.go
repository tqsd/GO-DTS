//
// EXAMPLE: Pareto based traffic source -> CLASSICAL LINK / GEWI
//
// Is running simulation in paralell number of Cores must be given
//
// modify function generate_setups() to run different setups
//

package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/tqsd/dts/link"
	"github.com/tqsd/dts/simulation"
	"github.com/tqsd/dts/traffic"
	"github.com/zenthangplus/goccm"
	"math"
	"os"
)

type parameters struct {
	//Traffic Source Parameters
	node_count int
	on_shape   float64
	on_scale   float64
	off_shape  float64
	off_scale  float64
	on_prob    float64

	//GEWI link parameters
	gain      float64
	cost      float64
	gewi_rate float64
	mult      float64
	E         float64
	e         float64
	gewi_B    float64

	//Classical link parameters
	link_rate float64
	link_B    float64
}

func generateSetups() []parameters {

	//shape parameter choosing
	res := int(5)
	on_shapes := make([]float64, res-1)
	off_shapes := make([]float64, res-1)
	for i := 0; i < res-1; i++ {
		// 1<shape<2
		on_shapes[i] = 1 + float64(i+1)/float64(res)
		off_shapes[i] = 1 + float64(i+1)/float64(res)
	}

	// Cluster Size
	nodeCounts := make([]int, 0)
	for i := 2; i < 60; i += 2 {
		nodeCounts = append(nodeCounts, i)
	}

	// Scale Parameters
	scale := float64(8)

	//Multiplication factors
	mults := make([]float64, 0)
	for i := float64(1); i < 6; i++ {
		mults = append(mults, math.Pow(2, i))
	}

	// Create Combinations of setups
	setups := make([]parameters, 0)
	for _, on := range on_shapes {
		for _, off := range on_shapes {
			for _, c := range nodeCounts {
				for _, m := range mults {

					cost := float64(1)
					gain := float64(1 / (m - 1))

					setups = append(setups, parameters{
						node_count: c,
						on_shape:   on,
						on_scale:   scale,
						off_shape:  off,
						off_scale:  8,
						on_prob:    0.5,
						gain:       gain,
						cost:       cost,
						gewi_rate:  float64(c) / 2,
						mult:       m,
						E:          2 * float64(c) * m * cost,
						e:          0,
						gewi_B:     2 * float64(c) * m,
						link_rate:  float64(c) / 2,
						link_B:     float64(2) * float64(c) * m,
					})
				}
			}
		}
	}
	return setups
}

func simulate(fileName string, setup parameters) {

	// Creating GEWI link
	gewi := link.NewGewi(setup.gain, setup.cost, setup.gewi_rate, setup.mult,
		setup.E, setup.e, setup.gewi_B)

	// Creating Classical link
	classic := link.NewClassical(setup.link_rate, setup.link_B)

	// Creating Traffic source
	source := traffic.NewParetoSelfSimilarSource(setup.node_count, setup.on_scale, setup.on_shape,
		setup.off_scale, setup.off_shape, setup.on_prob)

	ts := source.GetTrafficTimeSeries(LENGTH_OF_TIMESERIES)
	classic.ProcessGeneratedTraffic(ts)
	gewi.ProcessGeneratedTraffic(ts)
	//storables := []simulation.Storable{}{gewi, source}
	simulation.StoreResults("./results/", fileName, &classic, &gewi, &source)
}

// Set appropriate constants for your setup
const LENGTH_OF_TIMESERIES = 10000
const REPEAT_RUNS = 1
const MAX_NUM_OF_GO_ROUTINES = 1

func main() {

	// Figure out the destination file for the results
	args := os.Args
	fileName := "results.csv"
	if len(args) > 1 {
		fileName = args[1]
		if len(fileName) < 4 {
			fileName = fileName + ".csv"
		} else if fileName[len(fileName)-4:] != ".csv" {
			fileName = fileName + ".csv"
		}
	}
	setups := generateSetups()
	fmt.Println("Number of setups:", len(setups))
	fmt.Println("Number of runs to be executed:", len(setups)*REPEAT_RUNS)
	fmt.Println("Length of the timeseries:", LENGTH_OF_TIMESERIES)
	fmt.Println("Writing results to the: ./results/" + fileName)
	bar := progressbar.Default(int64(len(setups) * REPEAT_RUNS))
	c := goccm.New(MAX_NUM_OF_GO_ROUTINES)
	for r := 0; r < REPEAT_RUNS; r++ {
		for _, s := range setups {
			c.Wait()
			go func(fileName string, s parameters) {
				simulate(fileName, s)
				bar.Add(1)
				c.Done()
			}(fileName, s)
		}
	}
}
