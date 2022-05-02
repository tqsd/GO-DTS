package traffic

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type NBurstSource struct {
	NodeCount     int
	periodLengths []int
	periodTypes   []int

	// Pareto Distribution for on/off trains specific to this model
	T     int
	alpha float64
	theta float64
	gamma float64

	// Randomness objects
	on  TPT
	off distuv.Poisson
}

func NewNBurstTrafficSource(NodeCount, T int, alpha, theta, gamma float64) {
	on := NewTruncatedPowerTailDistribution(T, alpha, theta, gamma)

	var b1 [8]byte
	_, err := crypto_rand.Read(b1[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	off := distuv.Poisson{on.AverageDisc(), rand.NewSource(binary.LittleEndian.Uint64(b1[:]))}

	periodTypes := make([]int, NodeCount)
	periodLengths := make([]int, NodeCount)
	for i := 0; i < NodeCount; i++ {
		periodTypes[i] = rand.Intn(2)
		if periodTypes[i] == 1 {
			periodLengths[i] = int(on.Rand())
		} else {
			periodLengths[i] = int(off.Rand())
		}
	}
	return NBurstSource{
		NodeCount:     NodeCount,
		periodLengths: periodLengths,
		periodTypes:   periodTypes,
		T:             T,
		alpha:         alpha,
		theta:         theta,
		gamma:         gamma,
		on:            on,
		off:           off,
	}
}

// Determines amount of traffic in each moment, should be called many times to construt
func (link *NBurstSource) Tick() int {
	traffic := int(0)
	for i := 0; i < link.NodeCount; i++ {

		if link.periodTypes[i] == 1 {
			traffic++
		}
		link.periodLengths[i]--
		//At the end of the period decide on new period type and length
		if periodLengths <= 0 {
			r := rand.Float64()
			if r > on_prob {
				link.periodTypes[i] = 1
				link.periodLengths[i] = int(on.Rand())
			} else {
				link.periodTypes[i] = 0
				link.periodLengths[i] = int(off.Rand())
			}
		}
	}
	return traffic
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
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	//math_rand.Seed(time.Now().UnixNano() + 1)
	//fmt.Println(time.Now().UnixNano())
	values := make([]float64, 0)
	c := true
	Resolution := 0.01
	for x := float64(0); c; x += Resolution {
		values = append(values, r(T, x, alpha, theta, gamma))
		//fmt.Println(x, values[len(values)-1])
		if values[len(values)-1] < 0.0001 {
			c = false
		}
	}

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
