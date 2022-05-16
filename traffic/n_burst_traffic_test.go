package traffic

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func average(arr []int) float64 {

	avg := float64(0)
	for i, val := range arr {
		avg = (avg*float64(i) + float64(val)) / (float64(i) + 1)
	}
	return avg
}

func sum(arr []int) float64 {
	s := float64(0)
	for _, x := range arr {
		s += float64(x)
	}
	return s
}
func sumZero(arr []int) int {
	s := int(0)
	for _, x := range arr {
		if x == 0 {
			s += 1
		}
	}
	return s
}

func generateRandomTrafficSource() NBurstSource {
	rand.Seed(time.Now().UnixNano())
	//n := rand.Intn(100)
	n := 1
	T := rand.Intn(10) + 1 // t>0

	//Needs to be still figured out
	alpha := float64(1.4)
	theta := float64(0.5)
	gamma := math.Pow((1 / theta), (1 / alpha))

	return NewNBurstTrafficSource(n, T, alpha, theta, gamma)
}

func TestNewNBurstTrafficSource(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	T := rand.Intn(10) + 1 // t>0

	//Needs to be still figured out
	alpha := float64(1.4)
	theta := float64(0.5)
	gamma := math.Pow((1 / theta), (1 / alpha))
	fmt.Println(n, T, alpha, theta, gamma)
	source := NewNBurstTrafficSource(n, T, alpha, theta, gamma)

	if source.T != T {
		t.Errorf("Expected Truncation factor %v, but got %v", T, source.T)
	}

	if source.NodeCount != n {
		t.Errorf("Expected node count %v, but got %v", n, source.NodeCount)
	}

	if source.Alpha != alpha {
		t.Errorf("Expected alpha %v, but got %v", alpha, source.Alpha)
	}

	if source.Theta != theta {
		t.Errorf("Expected alpha %v, but got %v", theta, source.Theta)
	}

	if source.Gamma != gamma {
		t.Errorf("Expected alpha %v, but got %v", gamma, source.Gamma)
	}

}

func TestDistributions(t *testing.T) {
	return
	source := generateRandomTrafficSource()

	fmt.Println("Testing Distributions")

	run_number := 1000000
	onPeriods := make([]int, 0)
	offPeriods := make([]int, 0)

	for i := 0; i < run_number; i++ {
		onPeriods = append(onPeriods, int(source.On.Rand()))
		offPeriods = append(offPeriods, int(source.Off.Rand()))
	}
	fmt.Println(average(onPeriods))
	fmt.Println(average(offPeriods))

}

func TestAverage(t *testing.T) {
	source := generateRandomTrafficSource()

	fmt.Println("Average Traffic Rate:")

	run_number := 1000
	traffic := make([]int, 0)

	for i := 0; i < run_number; i++ {
		traffic = append(traffic, source.Tick())
	}
}

// Truncated Power Tail Tests

func TestNewTruncatedPowerTailDistribution(t *testing.T) {

	rand.Seed(time.Now().UnixNano())
	T := rand.Intn(10) + 1 // t>0
	//Needs to be still figured out
	alpha := rand.Float64() + 1
	theta := float64(0.5)
	gamma := math.Pow((1 / theta), (1 / alpha))

	tpt := NewTruncatedPowerTailDistribution(T, alpha, theta, gamma)

	if tpt.T != T {
		t.Errorf("Expected Truncation factor %v, but got %v", T, tpt.T)
	}

	if tpt.alpha != alpha {
		t.Errorf("Expected alpha %v, but got %v", alpha, tpt.alpha)
	}

	if tpt.theta != theta {
		t.Errorf("Expected alpha %v, but got %v", theta, tpt.theta)
	}

	if tpt.gamma != gamma {
		t.Errorf("Expected alpha %v, but got %v", gamma, tpt.gamma)
	}

}
