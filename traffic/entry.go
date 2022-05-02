package traffic

type source interface {
	Tick() int
	Average() float64
}
