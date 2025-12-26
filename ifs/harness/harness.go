package harness

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/tdarci/go-nothings/ifs/engine"
	"github.com/tdarci/go-nothings/ifs/generators/sierpinski"
	"github.com/tdarci/go-nothings/ifs/renderers"

	// "github.com/tdarci/go-nothings/ifs/renderers/ebitrenderer"
	// "github.com/tdarci/go-nothings/ifs/renderers/listpointrenderer"
	"github.com/tdarci/go-nothings/ifs/renderers/tcellrenderer"
	"github.com/tdarci/go-nothings/ifs/shared"
)

type System[T any] struct {
	Generator engine.Generator[T]
	Renderer  renderers.Renderer[T]
}

func RunSierpinski(iterations int) {
	Run(NewSierpinskiSystem(), iterations)
}

func Run[T any](sys *System[T], iterations int) {
	var wg sync.WaitGroup
	genChan := make(chan T, 50)
	ctx, canceFn := context.WithCancel(context.Background())

	wg.Go(func() {
		tckr := time.NewTicker(time.Millisecond * 10)
		defer tckr.Stop()
		for val := range genChan {
			select {
			case <-tckr.C:
			case <-ctx.Done():
				log.Println("[harness.Run] draw loop responding to closed context")
				return
			}
			sys.Renderer.Draw(val)
		}
		log.Println("[harness.Run] all generated points have been drawn")
	})

	wg.Go(func() {
		log.Println("[harness.Run] starting engine")
		engine.Run(ctx, sys.Generator, iterations, genChan)
		log.Println("[harness.Run] engine run complete")
		close(genChan)
	})

	sys.Renderer.Run()
	log.Println("Done running. Canceling context.")
	canceFn()
	log.Println("Waiting for goroutines to finish.")
	wg.Wait()
	log.Println("Done.")
}

func NewSierpinskiSystem() *System[shared.Point] {
	cfg := sierpinski.Config{
		ContainerEdgeLength: 200,
	}
	gen := sierpinski.New(cfg)
	rendCfg := renderers.PointRendererConfig{
		MinPoint: shared.Point{X: 0, Y: 0},
		MaxPoint: shared.Point{X: cfg.ContainerEdgeLength, Y: cfg.ContainerEdgeLength},
	}
	var pr renderers.PointRenderer
	// pr = listpointrenderer.New()
	// pr = ebitrenderer.New()
	pr = tcellrenderer.New()
	pr.Initialize(rendCfg)

	out := &System[shared.Point]{
		Generator: gen,
		Renderer:  pr,
	}
	return out
}
