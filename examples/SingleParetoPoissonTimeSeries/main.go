// Simple program to inspect the behavior of the link
package main

import (
	"bufio"
	"fmt"
	"github.com/tqsd/dts/link"
	"github.com/tqsd/dts/simulation"
	"github.com/tqsd/dts/traffic"
	"os"
	"strings"
)

type parameters struct {
	//Traffic Source Parameters
	on_shape  float64
	on_scale  float64
	off_shape float64
	off_scale float64
	on_prob   float64
	traffAvg  float64

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

type results struct {
	names               []string
	steps               []float64
	incoming            []float64
	gewi_transmitting   []float64
	gewi_dropping       []float64
	gewi_e_buffer_state []float64
	gewi_buffer_state   []float64
	clsc_transmitting   []float64
	clsc_dropping       []float64
	clsc_buffer_state   []float64
}

var names = []string{"Steps", "Incomming", "GEWI:Transmitting", "Gewi:Dropping", "Gewi:Entanglement.Buffer.state",
	"Gewi:Buffer.state",
	"Classic:Transmitting", "Classic:Dropping", "Classic:Buffer.State"}

// Writes data to file in a way that is understandable to gnuplot
func write_to_file_gnuplot_style(dirName, fileName string, r results) {
	simulation.CheckCreateDir(dirName)

	if exists, err := simulation.FileExists(dirName, fileName); err == nil {
		if exists {
			fmt.Println("File will be overwriten")
			os.Remove(dirName + "/" + fileName)
		}
	} else {
		panic("Error encountered: CSV File could not be created/read.")
	}

	f, errr := os.OpenFile(dirName+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errr != nil {
		panic("Error encountered: CSV File could not be created.")
	}

	datawriter := bufio.NewWriter(f)

	// column names
	fmt.Println(strings.Join(r.names[:], " ") + "\n")

	_, _ = datawriter.WriteString(strings.Join(r.names[:], " ") + "\n")

	for s := 0; s < len(r.steps); s++ {
		writeString := make([]string, 0)
		writeString = append(writeString, fmt.Sprintf("%f", r.steps[s]))
		writeString = append(writeString, fmt.Sprintf("%f", r.incoming[s]))
		writeString = append(writeString, fmt.Sprintf("%f", r.gewi_transmitting[s]))
		writeString = append(writeString, fmt.Sprintf("%f", r.gewi_dropping[s]))
		writeString = append(writeString, fmt.Sprintf("%f", r.gewi_e_buffer_state[s]))
		writeString = append(writeString, fmt.Sprintf("%f", r.gewi_buffer_state[s]))
		writeString = append(writeString, fmt.Sprintf("%f", r.clsc_transmitting[s]))
		writeString = append(writeString, fmt.Sprintf("%f", r.clsc_dropping[s]))
		writeString = append(writeString, fmt.Sprintf("%f", r.clsc_buffer_state[s]))
		_, _ = datawriter.WriteString(strings.Join(writeString[:], " ") + "\n")
	}

	datawriter.Flush()
	f.Close()
}

func main() {

	// Figure out the destination file for the results
	args := os.Args
	fileName := "results.data"
	if len(args) > 1 {
		fileName = args[1]
		if len(fileName) < 5 {
			fileName = fileName + ".data"
		} else if fileName[len(fileName)-5:] != ".data" {
			fileName = fileName + ".data"
		}
	}

	mult := float64(2)

	shape := 1.99
	avgTraff := float64(100)
	cost := float64(1)
	gain := float64(1 / (mult - 1))

	setup := parameters{
		on_scale:  8,
		on_shape:  shape,
		off_shape: 8,
		off_scale: shape,
		on_prob:   0.5,
		traffAvg:  avgTraff,
		gain:      gain,
		cost:      cost,
		gewi_rate: float64(avgTraff) / 2,
		mult:      mult,
		E:         avgTraff * 100,
		e:         0,
		gewi_B:    float64(10) * float64(avgTraff),
		link_rate: float64(avgTraff) / 2,
		link_B:    float64(10) * float64(avgTraff),
	}

	gewi := link.NewGewi(setup.gain, setup.cost, setup.gewi_rate,
		setup.mult, setup.E, setup.e, setup.gewi_B)
	classic := link.NewClassical(setup.link_rate, setup.link_B)

	// Creating Traffic source
	source := traffic.NewSingleParetoPoissonSource(setup.on_scale, setup.on_shape, setup.off_shape, setup.off_scale, setup.on_prob, setup.traffAvg)

	simulation_steps := 10000
	runway := int(float64(simulation_steps) * 0.0)
	result := results{
		names:               names,
		steps:               make([]float64, 0),
		incoming:            make([]float64, 0),
		gewi_transmitting:   make([]float64, 0),
		gewi_dropping:       make([]float64, 0),
		gewi_buffer_state:   make([]float64, 0),
		gewi_e_buffer_state: make([]float64, 0),
		clsc_transmitting:   make([]float64, 0),
		clsc_dropping:       make([]float64, 0),
		clsc_buffer_state:   make([]float64, 0),
	}
	//GENERATING TIMESTEPS
	for s := 0; s < simulation_steps; s++ {
		if s < runway {
			continue
		}
		gewi_transmitting_old := gewi.Transm
		gewi_dropping_old := gewi.Droppd
		clsc_transmitting_old := classic.Transm
		clsc_dropping_old := classic.Droppd
		incoming := source.Tick()
		fmt.Println(incoming)
		gewi.ProcessSingleStep(incoming)
		classic.ProcessSingleStep(incoming)
		result.steps = append(result.steps, float64(s))
		result.incoming = append(result.incoming, float64(incoming))
		result.gewi_transmitting = append(result.gewi_transmitting, gewi.Transm-gewi_transmitting_old)
		result.gewi_dropping = append(result.gewi_dropping, gewi.Droppd-gewi_dropping_old)
		result.gewi_buffer_state = append(result.gewi_buffer_state, gewi.CBuffS)
		result.gewi_e_buffer_state = append(result.gewi_e_buffer_state, gewi.EnBuffS)
		result.clsc_transmitting = append(result.clsc_transmitting, classic.Transm-clsc_transmitting_old)
		result.clsc_dropping = append(result.clsc_dropping, classic.Droppd-clsc_dropping_old)
		result.clsc_buffer_state = append(result.clsc_buffer_state, classic.CBuffS)
	}

	write_to_file_gnuplot_style("results", fileName, result)
}
