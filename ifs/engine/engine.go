package engine

import (
	"context"
	"log"
)

type Generator[T any] interface {
	Next(T) T
	Initialize() []T
}

func Run[T any](ctx context.Context, gen Generator[T], iterations int, outChan chan T) {
	var o T
	startOut := gen.Initialize()
	for _, o = range startOut {
		outChan <- o
	}

	for range iterations {
		o = gen.Next(o)
		select {
		case outChan <- o:
		case <-ctx.Done():
			log.Println("[engine.Run] responding to closed context")
			return
		}
	}
	log.Println("[engine.Run] completed all iterations")
}
