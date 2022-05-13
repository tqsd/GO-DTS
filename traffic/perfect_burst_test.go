package traffic

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewPerfectTrafficBurstSource(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	onPeriodLen := rand.Intn(200)
	offPeriodLen := rand.Intn(200)
	burst := NewPerfectTrafficBurstSource(n, onPeriodLen, offPeriodLen)

	//Checking the traffic source parameters
	if burst.NodeCount != n {
		t.Errorf("Expected NodeCount to equal %v, but got %v \n", n, burst.NodeCount)
	}

	if burst.OnPeriodLen != onPeriodLen {
		t.Errorf("Expected OnPeriodLen to equal %v, but got %v \n", onPeriodLen, burst.OnPeriodLen)
	}

	if burst.OffPeriodLen != offPeriodLen {
		t.Errorf("Expected OffPeriodLen to equal %v, but got %v \n", offPeriodLen, burst.OffPeriodLen)
	}
}

func TestNextPeriodON(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	onPeriodLen := rand.Intn(200)
	offPeriodLen := rand.Intn(200)
	burst := NewPerfectTrafficBurstSource(n, onPeriodLen, offPeriodLen)

	burst.NextPeriodON()
	if burst.Start != 1 {
		t.Errorf("Expected start to equal 1, but got %v \n", burst.Start)
	}
}

func TestNextPeriodOFF(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	onPeriodLen := rand.Intn(200)
	offPeriodLen := rand.Intn(200)
	burst := NewPerfectTrafficBurstSource(n, onPeriodLen, offPeriodLen)
	burst.NextPeriodOFF()
	if burst.Start != 0 {
		t.Errorf("Expected start to equal 1, but got %v \n", burst.Start)
	}
}

func TestTick(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	onPeriodLen := rand.Intn(200)
	offPeriodLen := rand.Intn(200)
	burst := NewPerfectTrafficBurstSource(n, onPeriodLen, offPeriodLen)
	burst.NextPeriodON()

	for i := 0; i < 2*(onPeriodLen+offPeriodLen); i++ {
		m := burst.Tick()
		if i%(onPeriodLen+offPeriodLen) < onPeriodLen {
			if m != n {
				t.Errorf("Expected on period %v, but got %v", m, n)
			}
		} else {
			if m != 0 {
				t.Errorf("Expected offPeriod 0, but got %v", m)
			}
		}
	}

	burst = NewPerfectTrafficBurstSource(n, onPeriodLen, offPeriodLen)
	burst.NextPeriodOFF()

	for i := 0; i < 2*(onPeriodLen+offPeriodLen); i++ {
		m := burst.Tick()
		if i%(onPeriodLen+offPeriodLen) < offPeriodLen {
			if m != 0 {
				t.Errorf("Expected offPeriod 0, but got %v", m)
			}
		} else {
			if m != n {
				t.Errorf("Expected on period %v, but got %v", m, n)
			}
		}
	}
}

func TestGetTrafficTimeSeries(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	onPeriodLen := rand.Intn(200)
	offPeriodLen := rand.Intn(200)

	burst := NewPerfectTrafficBurstSource(n, onPeriodLen, offPeriodLen)
	burst.NextPeriodON()
	ts := burst.GetTrafficTimeSeries(10 * (onPeriodLen + offPeriodLen))

	for i, m := range ts {
		if i%(onPeriodLen+offPeriodLen) < onPeriodLen {
			if m != n {
				t.Errorf("Expected on period %v, but got %v", m, n)
			}
		} else {
			if m != 0 {
				t.Errorf("Expected offPeriod 0, but got %v", m)
			}
		}
	}

	burst = NewPerfectTrafficBurstSource(n, onPeriodLen, offPeriodLen)
	burst.NextPeriodOFF()
	ts = burst.GetTrafficTimeSeries(10 * (onPeriodLen + offPeriodLen))

	for i, m := range ts {
		if i%(onPeriodLen+offPeriodLen) < offPeriodLen {
			if m != 0 {
				t.Errorf("Expected offPeriod 0, but got %v", m)
				break
			}
		} else {
			if m != n {
				t.Errorf("Expected on period %v, but got %v", m, n)
				break
			}
		}
	}
}

func TestGetValues(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	onPeriodLen := rand.Intn(200)
	offPeriodLen := rand.Intn(200)

	burst := NewPerfectTrafficBurstSource(n, onPeriodLen, offPeriodLen)
	burst.NextPeriodON()

	nameList, valueList := burst.GetValues()

	if len(nameList) != len(valueList) {
		t.Errorf("Expected two lists to be of the same length but got %v:%v", len(nameList), len(valueList))
	}
	for _, name := range nameList {
		if name[0:7] != "PBURST-" {
			t.Errorf("Expected 'PBURST-' in the name of the values'")
			break
		}
	}
	// Check if it is in the correct order

}
