package traffic

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"reflect"
)

type ParetoSelfSimilarSource struct {
	NodeCount     int
	PeriodLengths []int
	PeriodTypes   []int

	// Pareto Distribution for on/off trains specific to this model
	Shape_on  float64
	Scale_on  float64
	Shape_off float64
	Scale_off float64
	// Randomness objects
	On      distuv.Pareto
	Off     distuv.Pareto
	On_prob float64
}

// Creates Pareto Based Self Similar source
func NewParetoSelfSimilarSource(NodeCount int, Scale_on, Shape_on, Scale_off, Shape_off, On_prob float64) ParetoSelfSimilarSource {
	//Randomness setup
	var b1 [8]byte
	_, err := crypto_rand.Read(b1[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}

	var b2 [8]byte
	_, err = crypto_rand.Read(b2[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}

	on_pareto := distuv.Pareto{
		Xm:    Scale_on,
		Alpha: Shape_on,
		Src:   rand.NewSource(binary.LittleEndian.Uint64(b1[:])),
	}

	off_pareto := distuv.Pareto{
		Xm:    Scale_off,
		Alpha: Shape_off,
		Src:   rand.NewSource(binary.LittleEndian.Uint64(b2[:])),
	}

	//Generating initial node states
	PeriodTypes := make([]int, NodeCount)
	PeriodLengths := make([]int, NodeCount)
	for i := 0; i < NodeCount; i++ {
		r := rand.Float64()
		if r > On_prob {
			PeriodTypes[i] = 1
		} else {
			PeriodTypes[i] = 0
		}
		if PeriodTypes[i] == 1 {
			PeriodLengths[i] = int(on_pareto.Rand())
		} else {
			PeriodLengths[i] = int(off_pareto.Rand())
		}
	}

	return ParetoSelfSimilarSource{
		NodeCount:     NodeCount,
		PeriodLengths: PeriodLengths,
		PeriodTypes:   PeriodTypes,
		Shape_on:      Shape_on,
		Scale_on:      Scale_on,
		Shape_off:     Shape_off,
		Scale_off:     Scale_off,
		On:            on_pareto,
		Off:           off_pareto,
		On_prob:       On_prob,
	}
}

func (link *ParetoSelfSimilarSource) Tick() int {
	traffic := int(0)
	for i := 0; i < link.NodeCount; i++ {

		if link.PeriodTypes[i] == 1 {
			traffic++
		}
		link.PeriodLengths[i]--
		//At the end of the period decide on new period type and length
		if link.PeriodLengths[i] <= 0 {
			r := rand.Float64()
			if r > link.On_prob {
				link.PeriodTypes[i] = 1
				link.PeriodLengths[i] = int(link.On.Rand())
			} else {
				link.PeriodTypes[i] = 0
				link.PeriodLengths[i] = int(link.Off.Rand())
			}
		}
	}
	return traffic
}

// Uses the source to generate the traffic time series
func (link *ParetoSelfSimilarSource) GetTrafficTimeSeries(length int) []int {
	ts := make([]int, 0)
	for i := 0; i < length; i++ {
		ts = append(ts, link.Tick())
	}
	return ts
}

func (link *ParetoSelfSimilarSource) Average() (float64, error) {
	if link.Shape_on <= 1 || link.Shape_off <= 1 {
		return 0, errors.New("Expectation values infinite; Shape_on and Shape_off should be >1")
	}
	on := (link.Shape_on * link.Scale_on) / (link.Shape_on - 1)
	off := (link.Shape_off * link.Scale_off) / (link.Shape_off - 1)
	result := (link.On_prob * on) / (link.On_prob*on + (1-link.On_prob)*off)
	return result, nil
}

func (link *ParetoSelfSimilarSource) AverageEmpirical(rounds int) float64 {
	traffic_sum := int(0)
	for i := 0; i < rounds; i++ {
		traffic_sum += link.Tick()
	}
	return float64(traffic_sum) / float64(rounds)
}

// For the csv logs Returns the list of the values
func (obj *ParetoSelfSimilarSource) GetValues() ([]string, []string) {
	e := reflect.ValueOf(&obj).Elem().Elem()

	nameList := make([]string, 0)
	valueList := make([]string, 0)

	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varType := e.Type().Field(i).Type
		varValue := e.Field(i).Interface()

		if varType.Kind() != reflect.Int && varType.Kind() != reflect.Float64 {
			continue
		}

		nameList = append(nameList, fmt.Sprintf("PARETO-%v", varName))
		valueList = append(valueList, fmt.Sprintf("%v", varValue))
	}

	return nameList, valueList
}
