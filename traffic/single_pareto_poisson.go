package traffic

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type SingleParetoPoissonSource struct {
	PeriodLength int
	PeriodType   int

	// Pareto Distribution for on/off trains specific to this model
	Shape_on  float64
	Scale_on  float64
	Shape_off float64
	Scale_off float64
	// Randomness objects
	On      distuv.Pareto
	Off     distuv.Pareto
	Traffic distuv.Poisson

	On_prob float64
}

// Creates Pareto Based Self Similar source
func NewSingleParetoPoissonSource(Scale_on, Shape_on, Scale_off, Shape_off,
	On_prob, Traffic_Poisson float64) SingleParetoPoissonSource {
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

	var b3 [8]byte
	_, err = crypto_rand.Read(b3[:])
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
	Traffic := distuv.Poisson{
		Lambda: Traffic_Poisson,
		Src:    rand.NewSource(binary.LittleEndian.Uint64(b3[:])),
	}

	//Generating initial node states
	PeriodTypes := int(0)
	PeriodLengths := int(0)
	r := rand.Float64()
	if r > On_prob {
		PeriodTypes = 1
	} else {
		PeriodTypes = 0
	}
	if PeriodTypes == 1 {
		PeriodLengths = int(on_pareto.Rand())
	} else {
		PeriodLengths = int(off_pareto.Rand())
	}

	return SingleParetoPoissonSource{
		PeriodLength: PeriodLengths,
		PeriodType:   PeriodTypes,
		Shape_on:     Shape_on,
		Scale_on:     Scale_on,
		Shape_off:    Shape_off,
		Scale_off:    Scale_off,
		On:           on_pareto,
		Off:          off_pareto,
		Traffic:      Traffic,
		On_prob:      On_prob,
	}
}

func (link *SingleParetoPoissonSource) Tick() int {
	traffic := int(0)
	//At the end of the period decide on new period type and length
	if link.PeriodType > 0 {
		traffic = int(link.Traffic.Rand())
	}
	link.PeriodLength--
	if link.PeriodLength <= 0 {
		r := rand.Float64()
		if r > link.On_prob {
			link.PeriodType = 1
			link.PeriodLength = int(link.On.Rand())
		} else {
			link.PeriodType = 0
			link.PeriodLength = int(link.Off.Rand())
		}
	}
	return traffic
}

// Uses the source to generate the traffic time series
func (link *SingleParetoPoissonSource) GetTrafficTimeSeries(length int) []int {
	ts := make([]int, 0)
	for i := 0; i < length; i++ {
		ts = append(ts, link.Tick())
	}
	return ts
}

func (link *SingleParetoPoissonSource) Average() (float64, error) {
	if link.Shape_on <= 1 || link.Shape_off <= 1 {
		return 0, errors.New("Expectation values infinite; Shape_on and Shape_off should be >1")
	}
	on := (link.Shape_on * link.Scale_on) / (link.Shape_on - 1)
	off := (link.Shape_off * link.Scale_off) / (link.Shape_off - 1)
	result := (link.On_prob * on) / (link.On_prob*on + (1-link.On_prob)*off)
	return result, nil
}

func (link *SingleParetoPoissonSource) AverageEmpirical(rounds int) float64 {
	traffic_sum := int(0)
	for i := 0; i < rounds; i++ {
		traffic_sum += link.Tick()
	}
	return float64(traffic_sum) / float64(rounds)
}

// For the csv logs Returns the list of the values
func (obj *SingleParetoPoissonSource) GetValues() ([]string, []string) {
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

		nameList = append(nameList, fmt.Sprintf("SiParPoiss-%v", varName))
		valueList = append(valueList, fmt.Sprintf("%v", varValue))
	}

	return nameList, valueList
}
