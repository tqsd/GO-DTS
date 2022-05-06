package traffic

import (
	"fmt"
	"reflect"
)

type PerfectBurstSource struct {
	NodeCount    int
	OnPeriodLen  int
	OffPeriodLen int
	Step         int
	Start        int
}

func NewPerfectTrafficBurstSource(n, on, off int) PerfectBurstSource {
	return PerfectBurstSource{
		NodeCount:    n,
		OnPeriodLen:  on,
		OffPeriodLen: off,
		Step:         0,
		Start:        0,
	}
}

func (link *PerfectBurstSource) NextPeriodON() {
	link.Start = 1
}

func (link *PerfectBurstSource) NextPeriodOFF() {
	link.Start = 0
}

func (link *PerfectBurstSource) Tick() int {
	traffic := 0
	totalPeriodLen := link.OnPeriodLen + link.OffPeriodLen
	//It is on period
	if link.Start == 1 {
		if link.Step%(totalPeriodLen) < link.OnPeriodLen {
			traffic = link.NodeCount
		}
	} else {
		if link.Step%(totalPeriodLen) >= link.OffPeriodLen {
			traffic = link.NodeCount
		}
	}
	link.Step++
	return traffic
}

func (link *PerfectBurstSource) GetTrafficTimeSeries(length int) []int {
	ts := make([]int, 0)
	for i := 0; i < length; i++ {
		ts = append(ts, link.Tick())
	}
	return ts
}

// For the csv logs Returns the list of the values
func (obj *PerfectBurstSource) GetValues() ([]string, []string) {
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

		nameList = append(nameList, fmt.Sprintf("PBURST-%v", varName))
		valueList = append(valueList, fmt.Sprintf("%v", varValue))
	}

	return nameList, valueList
}
