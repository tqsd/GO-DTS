package link

import (
	"fmt"
	"reflect"
)

//Struct describing GEWI link
type gewi struct {
	//How much entanglement is gained per empty slot
	Gain float64
	//How much entanglement is used to assist with one slot
	Cost float64
	//How many slots per one tick
	Rate float64
	//Assistance multiplier
	Mult float64
	// Entanglement buffer size
	EnBuff float64
	// Entanglement buffer state
	EnBuffS float64
	// Classical buffer size
	CBuff float64
	// Classical buffer state
	CBuffS float64
	// Total packets received
	Recv float64
	// Total packets transmitted
	Transm float64
	// Total packets dropped
	Droppd float64
	// Total packets dropped
	Step int
	// Average Entanglement Buffer state
	AvgE float64
	// Average Classical Buffer state
	AvgC float64
	// Wait time
	Wait []uint64
	// Mean waiting time
	MWT float64
	// Mean waiting time standard deviation
	MWTstddev float64
	// Mean waiting time variance
	MWTvar float64
	// Mean waiting time square
	MWTsquare float64
}

// Returns a struct describing the state of the gewi link
// e-> Entanglement buffer fulness: determines how full the entanglement buffer
// should be at the beginnig e \in [0,1]
func NewGewi(G, C, R, M, E, e, B float64) gewi {
	return gewi{
		Gain:      G,
		Cost:      C,
		Rate:      R,
		Mult:      M,
		EnBuff:    E,
		EnBuffS:   float64(int(e * E)),
		CBuff:     B,
		CBuffS:    0,
		Recv:      0,
		Transm:    0,
		Droppd:    0,
		Step:      0,
		AvgE:      0,
		AvgC:      0,
		Wait:      make([]uint64, int(B)),
		MWT:       0,
		MWTstddev: 0,
		MWTvar:    0,
	}
}

func (link *gewi) ProcessSingleStep(incomming int) {
	link.Step++
	link.CBuffS += float64(incomming)
	link.Recv += float64(incomming)

	// Deternime how much data and how much entanglement
	// can be sent
	stepTransmission := float64(0)
	if link.CBuffS < link.Rate {
		// If after transmission unused rate exists
		stepTransmission = link.CBuffS
		link.DistributeEntanglement(link.Rate - stepTransmission)
		link.CBuffS = 0
	} else {
		for R := link.Rate; R > 0; R-- {
			if link.CBuffS < R {
				// More rate than data to send
				stepTransmission += link.CBuffS
				link.DistributeEntanglement(R - link.CBuffS)
				break
			} else if link.EnBuffS >= link.Cost && link.CBuffS >= link.Mult {
				// Enough entanglement and data left for assisted transmission
				stepTransmission += link.Mult
				link.EnBuffS -= link.Cost
				link.CBuffS -= link.Mult
			} else {
				// link.CBuffS > R || (link.EnBuffS < link.Cost)
				// More data left to send than unasisted rate permits OR
				// Not enough entanglement left  OR
				// Less data left to send than assisted could send -> No padding
				stepTransmission += R
				link.CBuffS -= R
				break
			}
		}
	}

	// Entanglement is Generated
	// Amount to transmit is calculated
	// Classical buffer state is reduced to appropriate amount
	link.UpdateMeanWaitingTime(stepTransmission)
	link.AdjustWaitQueue(stepTransmission)
	link.Transmit(stepTransmission)
	link.DropExcessData()
	link.IncreaseWaitTimes()
	link.UpdateAverageCBuffState()
	link.UpdateAverageEBuffState()

}

func (link *gewi) Transmit(stepTransmission float64) {
	link.Transm += stepTransmission
}

// Update average entanglement buffer state
func (link *gewi) UpdateAverageCBuffState() {
	link.AvgC = ((link.AvgC * (float64(link.Step) - 1)) + link.CBuffS) / float64(link.Step)
}

// Update average entantglement buffer state
func (link *gewi) UpdateAverageEBuffState() {
	link.AvgE = ((link.AvgE * (float64(link.Step) - 1)) + link.EnBuffS) / float64(link.Step)
}

//Increase the wait times for the data that is in the queue
func (link *gewi) IncreaseWaitTimes() {
	for i := 0; i < int(link.CBuffS); i++ {
		link.Wait[i]++
	}
}

//Increase the drop counter
func (link *gewi) DropExcessData() {
	if link.CBuffS > link.CBuff {
		link.Droppd += (link.CBuffS - link.CBuff)
		link.CBuffS = link.CBuff
	}
}

//Distributes the entanglement; given the unused rate
func (link *gewi) DistributeEntanglement(idleRate float64) {
	link.EnBuffS += idleRate * link.Gain
	if link.EnBuffS > link.EnBuff {
		link.EnBuffS = link.EnBuff
	}
}

// Calulcate how long the data had to wait before sending this step
func (link *gewi) UpdateMeanWaitingTime(transmission float64) {
	waitTime := float64(0)
	for i := 0; i < int(transmission); i++ {
		waitTime += float64(link.Wait[i])
	}

	link.MWT = (link.MWT*link.Transm + waitTime) / (transmission + link.Transm)
}

func (link *gewi) AdjustWaitQueue(stepTransmission float64) {
	// If nothing was sent no adjusting needed
	if stepTransmission == 0 {
		return
	}
	// Shift the wait queue
	zeros := make([]uint64, int(stepTransmission))
	for i := 0; i < int(stepTransmission); i++ {
		zeros[i] = 0
	}

	link.Wait = append(link.Wait[int(stepTransmission):len(link.Wait)], zeros...)
}

func (link *gewi) ProcessGeneratedTraffic(traffic []int) {
	for i := 0; i < len(traffic); i++ {
		link.ProcessSingleStep(traffic[i])
	}
}

// For the csv logs Returns the list of the values
func (obj *gewi) GetValues() ([]string, []string) {

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

		nameList = append(nameList, fmt.Sprintf("GEWI-%v", varName))
		valueList = append(valueList, fmt.Sprintf("%v", varValue))
	}

	return nameList, valueList
}
