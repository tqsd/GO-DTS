package traffic

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type ParetoSelfSimilarSource struct {
	NodeCount     int
	periodLengths []int
	periodTypes   []int

	// Pareto Distribution for on/off trains specific to this model
	shape_on  float64
	scale_on  float64
	shape_off float64
	scale_off float64
	// Randomness objects
	on      distuv.Pareto
	off     distuv.Pareto
	on_prob float64
}

// Creates Pareto Based Self Similar source
func NewParetoSelfSimilarSource(NodeCount int, scale_on, shape_on, scale_off, shape_off, on_prob float64) ParetoSelfSimilarSource {
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
		Xm:    on.scale,
		Alpha: on.alpha,
		Src:   rand.NewSource(binary.LittleEndian.Uint64(b1[:])),
	}

	off_pareto := distuv.Pareto{
		Xm:    off.scale,
		Alpha: off.alpha,
		Src:   rand.NewSource(binary.LittleEndian.Uint64(b2[:])),
	}

	//Generating initial node states
	periodTypes := make([]int, NodeCount)
	periodLengths := make([]int, NodeCount)
	for i = 0; i < NodeCount; i++ {
		r := rand.Float64()
		if r > on_prob {
			periodTypes[i] = 1
		} else {
			periodTypes[i] = 0
		}
		if periodTypes[i] == 1 {
			periodLengths[i] = int(on_pareto.Rand())
		} else {
			periodLengths[i] = int(off_pareto.Rand())
		}
	}

	return ParetoSelfSimilarSource{
		NodeCount:     NodeCount,
		periodLengths: periodLengths,
		periodTypes:   periodTypes,
		shape_on:      shape_on,
		scale_on:      scale_on,
		shape_off:     shape_off,
		scale_off:     scale_off,
		on:            on_pareto,
		off:           off_pareto,
		on_prob:       on_prob,
	}
}

func (link *ParetoSelfSimilarSource) Tick() int {
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

func (link *ParetoSelfSimilarSource) Average() (float64, error) {
	if link.shape_on <= 1 || link.shape_off <= 1 {
		return 0, errors.New("Expectation values infinite; shape_on and shape_off should be >1")
	}
	on := (link.shape_on * link.scale_on) / (link.shape_on - 1)
	off := (link.shape_off * link.scale_off) / (link.shape_off - 1)
	result := (link.on_prob * on) / (link.on_prob*on + (1-link.on_prob)*off)
	return result, nil
}

func (link *ParetoSelfSimilarSource) AverageEmpirical(rounds int) float64 {
	traffic_sum := int(0)
	for i = 0; i < rounds; i++ {
		traffic_sum += link.Tick()
	}
	return float64(traffic_sum) / float64(rounds)
}
