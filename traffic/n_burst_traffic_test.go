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

func generateRandomTrafficSource() NBurstSource {
	fmt.Println("test")
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	T := rand.Intn(10) + 1 // t>0

	//Needs to be still figured out
	alpha := rand.Float64() * 2
	theta := float64(0.5)
	gamma := math.Pow((1 / theta), (1 / alpha))

	return NewNBurstTrafficSource(n, T, alpha, theta, gamma)
}

func TestNewNBurstTrafficSource(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	T := rand.Intn(10) + 1 // t>0

	//Needs to be still figured out
	alpha := rand.Float64() * 2
	theta := float64(0.5)
	gamma := math.Pow((1 / theta), (1 / alpha))

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
	source := generateRandomTrafficSource()

	fmt.Println("Testing Distributions")

	run_number := 100
	onPeriods := make([]int, 0)
	offPeriods := make([]int, 0)

	for i := 0; i < run_number; i++ {
		onPeriods = append(onPeriods, source.On.Rand())
		offPeriods = append(offPeriods, int(source.Off.Rand()))
	}
	fmt.Println(average(onPeriods))
	fmt.Println(average(offPeriods))

}

// Truncated Power Tail Tests

func TestNewTruncatedPowerTailDistribution(t *testing.T) {

	rand.Seed(time.Now().UnixNano())
	T := rand.Intn(10) + 1 // t>0
	//Needs to be still figured out
	alpha := rand.Float64() * 2
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
