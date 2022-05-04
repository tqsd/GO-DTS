package link

import (
	"fmt"
	"reflect"
)

//Struct describing Unassisted link
type classical struct {
	//How many slots per one tick
	Rate float64
	// Entanglement buffer size
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

func NewClassical(R, B float64) classical {

	return classical{
		Rate:   R,
		CBuff:  B,
		CBuffS: 0,
		Recv:   0,
		Transm: 0,
		Droppd: 0,
		Step:   0,
		AvgC:   0,
		Wait:   make([]uint64, int(B)),
	}
}

func (link *classical) ProcessSingleStep(incomming int) {
	link.Step++
	link.CBuffS += float64(incomming)
	link.Recv += float64(incomming)

	// Deternime how much data and how much entanglement
	// can be sent
	stepTransmission := float64(0)

	if link.CBuffS > link.Rate {
		// If there is more data in the buffer than can be sent in one
		// step, send rate's worth
		stepTransmission += link.Rate

	} else {
		// If there is less data that can be sent, send what has to be sent
		stepTransmission += link.CBuffS
	}

	// Amount to transmit is calculated
	// Classical buffer state is reduced to appropriate amount
	link.UpdateMeanWaitingTime(stepTransmission)
	link.AdjustWaitQueue(stepTransmission)
	link.Transmit(stepTransmission)
	link.DropExcessData()
	link.IncreaseWaitTimes()
	link.UpdateAverageCBuffState()
}

func (link *classical) Transmit(stepTransmission float64) {
	link.Transm += stepTransmission
}

// Update average entanglement buffer state
func (link *classical) UpdateAverageCBuffState() {
	link.AvgC = ((link.AvgC * (float64(link.Step) - 1)) + link.CBuffS) / float64(link.Step)
}

//Increase the wait times for the data that is in the queue
func (link *classical) IncreaseWaitTimes() {
	for i := 0; i < int(link.CBuffS); i++ {
		link.Wait[i]++
	}
}

//Increase the drop counter
func (link *classical) DropExcessData() {
	if link.CBuffS > link.CBuff {
		link.Droppd += (link.CBuffS - link.CBuff)
		link.CBuffS = link.CBuff
	}
}

// Calulcate how long the data had to wait before sending this step
func (link *classical) UpdateMeanWaitingTime(transmission float64) {
	waitTime := float64(0)
	for i := 0; i < int(transmission); i++ {
		waitTime += float64(link.Wait[i])
	}

	link.MWT = (link.MWT*link.Transm + waitTime) / (transmission + link.Transm)
}

// Shift the wait queue
func (link *classical) AdjustWaitQueue(stepTransmission float64) {
	// If no data was sent don't update anything
	if stepTransmission == 0 {
		return
	}
	zeros := make([]uint64, int(stepTransmission))
	for i := 0; i < int(stepTransmission); i++ {
		zeros[i] = 0
	}
	link.Wait = append(link.Wait[int(stepTransmission):len(link.Wait)], zeros...)
}

func (link *classical) ProcessGeneratedTraffic(traffic []int) {
	for i := 0; i < len(traffic); i++ {
		link.ProcessSingleStep(traffic[i])
	}
}

// For the csv logs Returns the list of the values
func (obj *classical) GetValues() ([]string, []string) {

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
