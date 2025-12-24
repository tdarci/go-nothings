package engine

const (
	iterations = 20
)

type Generator[T any] interface {
	Next(T) T
	Initialize() []T
}

func Run[T any](gen Generator[T], outChan chan T) {
	var o T
	startOut := gen.Initialize()
	for _, o = range startOut {
		outChan <- o
	}

	for range iterations {
		o = gen.Next(o)
		outChan <- o
	}
}
