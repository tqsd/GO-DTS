package traffic

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	math_rand "math/rand"
	"reflect"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type NBurstSource struct {
	NodeCount     int
	PeriodLengths []int
	PeriodTypes   []int

	// Pareto Distribution for on/off trains specific to this model
	T     int
	Alpha float64
	Theta float64
	Gamma float64

	// Randomness objects
	On  TPT
	Off distuv.Poisson

	// Checking
	Avg            float64
	Step           int
	OnPeriod       int
	OnPeriodCount  int
	OffPeriod      int
	OffPeriodCount int
	OffCounter     int
}

func NewNBurstTrafficSource(NodeCount, T int, alpha, theta, gamma float64) NBurstSource {

	on := NewTruncatedPowerTailDistribution(T, alpha, theta, gamma)

	var b1 [8]byte
	_, err := crypto_rand.Read(b1[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	off := distuv.Poisson{on.AverageDisc(), rand.NewSource(binary.LittleEndian.Uint64(b1[:]))}

	periodTypes := make([]int, NodeCount)
	periodLengths := make([]int, NodeCount)

	pType := 0
	// Decidint initial distribution
	for i := 0; i < NodeCount; i++ {
		periodTypes[i] = pType

		if periodTypes[i] == 1 {
			periodLengths[i] = int(on.Rand())
			pType = 1
		} else {
			periodLengths[i] = int(off.Rand())
			pType = 1
		}
	}
	return NBurstSource{
		NodeCount:     NodeCount,
		PeriodLengths: periodLengths,
		PeriodTypes:   periodTypes,
		T:             T,
		Alpha:         alpha,
		Theta:         theta,
		Gamma:         gamma,
		On:            on,
		Off:           off,
	}
}

// Determines amount of traffic in each moment, should be called many times to construt
func (link *NBurstSource) Tick() int {
	traffic := int(0)
	for i := 0; i < link.NodeCount; i++ {

		if link.PeriodTypes[i] == 1 {
			traffic++
		}
		link.PeriodLengths[i]--
		//At the end of the period decide on new period type and length
		if link.PeriodLengths[i] <= 0 {
			if link.PeriodTypes[i] == 0 {
				//Traffic is interchangable
				link.PeriodTypes[i] = 1
				link.PeriodLengths[i] = int(link.On.Rand())
				link.OnPeriod += link.PeriodLengths[i]
				link.OnPeriodCount++
			} else {
				link.PeriodTypes[i] = 0
				r := link.Off.Rand()
				link.PeriodLengths[i] = int(r)
				link.OffPeriod += link.PeriodLengths[i]
				link.OffPeriodCount++

				if r == 0 {
					// has to be bigger than 0,
					link.PeriodTypes[i] = 1
					link.PeriodLengths[i] = int(link.On.Rand())
					link.OnPeriod += link.PeriodLengths[i]
					link.OnPeriodCount++
				}
			}
		}
	}
	//fmt.Println(traffic, link.PeriodTypes, link.PeriodLengths, "OnC:", link.OnPeriodCount, "OffC:", link.OffPeriodCount,
	//	"OnPT", link.OnPeriod, "OffPT:", link.OffPeriod, "/", link.OffCounter)
	//fmt.Println(link.PeriodLengths, link.PeriodTypes, avrg(link.PeriodTypes))
	link.Avg = (link.Avg*float64(link.Step) + float64(traffic)) / (float64(link.Step) + 1)
	link.Step++
	if traffic == 0 {
		link.OffCounter++
	}
	return traffic
}

func (link *NBurstSource) GetTrafficTimeSeries(length int) []int {
	ts := make([]int, 0)
	for i := 0; i < length; i++ {
		ts = append(ts, link.Tick())
	}
	return ts
}

//Truncated Power Tail distribution ---------
type TPT struct {
	T          int
	alpha      float64
	theta      float64
	gamma      float64
	Resolution float64
	Pdf        []float64
	Cdf        []float64
}

// Truncated power tail distribution
func r(T int, x, alpha, theta, gamma float64) float64 {
	mu := (1 - theta) / (1 - gamma*theta)
	temp := float64(0)
	for j := float64(0); j < float64(T); j++ {
		temp += math.Pow(theta, j) * math.Exp(-mu*x/(math.Pow(gamma, j)))
	}
	return ((1 - theta) / (1 - math.Pow(theta, float64(T))) * temp)
}

func NewTruncatedPowerTailDistribution(T int, alpha, theta, gamma float64) TPT {

	if T < 1 {
		panic("Truncation factors must be greater than 0")
	}
	if gamma < 1 {
		panic("gamma should be > 1")
	}
	if theta > 1 || theta < 0 {
		panic("It should hold 0 < theta < 1")
	}

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	values := make([]float64, 0)
	c := true
	Resolution := 0.001

	for x := float64(0); c; x += Resolution {
		values = append(values, r(T, x, alpha, theta, gamma))
		if values[len(values)-1] < 0.1 {
			c = false
		}
		if values[len(values)-1] == math.Inf(+1) {
			panic("INFINITY")
		}
	}

	mu := (1 - theta) / (1 - gamma*theta)
	fmt.Println("MU: ", mu)
	valuesum := float64(0)
	for i := 0; i < len(values); i++ {
		valuesum += values[i]
	}
	for i := 0; i < len(values); i++ {
		values[i] = values[i] / valuesum
	}
	Cdf := make([]float64, len(values))
	Cdf[0] = values[0]
	for i := 1; i < len(Cdf); i++ {
		Cdf[i] = Cdf[i-1] + values[i]
	}
	//fmt.Println(values)
	//fmt.Println(cdf)
	return TPT{T: T, alpha: alpha, theta: theta, gamma: gamma, Pdf: values, Cdf: Cdf, Resolution: Resolution}

}

func (tpt *TPT) Rand() int {
	uniform_float := math_rand.Float64()
	result := int(0)
	for i := 0; i < len(tpt.Cdf)-1; i++ {
		if uniform_float < tpt.Cdf[i] {
			break
		}
		result++
	}

	//return float64(result) * tpt.Resolution
	return int(float64(result)*tpt.Resolution + 1)
}

func (tpt *TPT) Moments(l uint) float64 {
	temp1 := (1 - tpt.theta) / (1 - math.Pow(tpt.theta, float64(tpt.T)))
	temp2 := (1 - math.Pow(tpt.theta*math.Pow(tpt.gamma, float64(l)), float64(tpt.T))) / (1 - (tpt.theta * math.Pow(tpt.gamma, float64(l))))
	factorial := 1
	for i := 1; i <= int(l); i++ {
		factorial = factorial * i
	}
	mu := (1 - tpt.theta) / (1 - tpt.gamma*tpt.theta)
	temp3 := float64(factorial) / math.Pow(mu, float64(l))

	fmt.Println("Third Term:", temp3, "With T:", tpt.T, " ,theta:", tpt.theta, " ,gamma:", tpt.gamma)
	return temp1 * temp2 * temp3
}

// Calculated Expectation value for the distribution
func (tpt *TPT) AverageDisc() float64 {
	average := float64(0)
	pdf := make([]float64, 0)
	pdf = append(pdf, 0)
	x := 0
	for i := 0; i < len(tpt.Pdf); i++ {
		if i > (x+1)*int(1/tpt.Resolution) {
			x++
			pdf = append(pdf, 0)
		}
		pdf[x] += tpt.Pdf[i]
	}
	for i := 0; i < len(pdf); i++ {
		average += float64(i+1) * pdf[i]
	}
	return average
}

// For the csv logs Returns the list of the values
func (link *NBurstSource) GetValues() ([]string, []string) {

	e := reflect.ValueOf(&link).Elem().Elem()
	nameList := make([]string, 0)
	valueList := make([]string, 0)

	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varType := e.Type().Field(i).Type
		varValue := e.Field(i).Interface()

		if varType.Kind() != reflect.Int && varType.Kind() != reflect.Float64 {
			continue
		}
		nameList = append(nameList, fmt.Sprintf("NBURST-%v", varName))
		valueList = append(valueList, fmt.Sprintf("%v", varValue))
	}
	return nameList, valueList
}

func avrg(arr []int) float64 {
	avg := float64(0)
	for i, val := range arr {
		avg = (avg*float64(i) + float64(val)) / (float64(i) + 1)
	}
	return avg
}
