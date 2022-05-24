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
	"os"
)

type parameters struct {
	// Pareto Distribution for on/off trains specific to this model
	Shape_on  float64
	Scale_on  float64
	Shape_off float64
	Scale_off float64

	//Traffic shape when traffic is generated
	Traffic_Poisson float64

	on_prob float64

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
	setups := make([]parameters, 0)
	on_prob := float64(0.5)
	burst_rate := float64(100)
	on_shapes := [3]float64{1.1, 1.5, 1.9}
	res := int(10)
	off_shapes := make([]float64, res-1)
	for i := 0; i < res-1; i++ {
		off_shapes[i] = 1 + float64(i+1)/float64(res)
	}
	scale := float64(8)
	mults := [3]float64{2, 32, 128}

	for _, on := range on_shapes {
		for _, off := range off_shapes {
			for _, m := range mults {
				cost := float64(1)
				gain := float64(1 / (m - 1))

				B := 10 * burst_rate
				E_on := on * scale / (on - 1)
				E_off := off * scale / (off - 1)
				rate := burst_rate / (((1-on_prob)/on_prob)*(E_off/E_on) + 1)

				setups = append(setups, parameters{
					Shape_on:        on,
					Scale_on:        scale,
					Shape_off:       off,
					Scale_off:       scale,
					Traffic_Poisson: burst_rate,
					on_prob:         on_prob,
					gain:            gain,
					cost:            cost,
					mult:            m,
					e:               0,
					E:               LENGTH_OF_TIMESERIES, // INFINITE
					gewi_rate:       rate,
					gewi_B:          B,
					link_rate:       rate,
					link_B:          B,
				})
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
	source := traffic.NewSingleParetoPoissonSource(
		setup.Scale_on,
		setup.Shape_on,
		setup.Scale_off,
		setup.Shape_off,
		setup.on_prob,
		setup.Traffic_Poisson)

	ts := source.GetTrafficTimeSeries(LENGTH_OF_TIMESERIES)

	fmt.Println("AVG:", average(ts), "Calculated Average", setup.link_rate, setup.Shape_on, setup.Shape_off)
	classic.ProcessGeneratedTraffic(ts)
	gewi.ProcessGeneratedTraffic(ts)
	//storables := []simulation.Storable{}{gewi, source}
	simulation.StoreResults("./results/", fileName, &classic, &gewi, &source)
}

// Set appropriate constants for your setup
const LENGTH_OF_TIMESERIES = 500000
const REPEAT_RUNS = 10
const MAX_NUM_OF_GO_ROUTINES = 100

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

func average(arr []int) float64 {
	avg := float64(0)
	for i, val := range arr {
		avg = (avg*float64(i) + float64(val)) / (float64(i) + 1)
	}
	return avg
}
