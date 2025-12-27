package harness

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/tdarci/go-nothings/ifs/config"
	"github.com/tdarci/go-nothings/ifs/engine"
	"github.com/tdarci/go-nothings/ifs/generators/fern"
	"github.com/tdarci/go-nothings/ifs/generators/negtest"
	"github.com/tdarci/go-nothings/ifs/generators/sierpinski"
	"github.com/tdarci/go-nothings/ifs/renderers"
	"github.com/tdarci/go-nothings/ifs/renderers/ebitrenderer"

	// "github.com/tdarci/go-nothings/ifs/renderers/listpointrenderer"
	// "github.com/tdarci/go-nothings/ifs/renderers/tcellrenderer"
	"github.com/tdarci/go-nothings/ifs/shared"
)

type System[T any] struct {
	Generator engine.Generator[T]
	Renderer  renderers.Renderer[T]
}

func Run[T any](sys *System[T], maxIterations int, drawingDelay time.Duration) {
	var wg sync.WaitGroup
	genChan := make(chan T, 50)
	ctx, canceFn := context.WithCancel(context.Background())

	wg.Go(func() {
		tckr := time.NewTicker(drawingDelay)
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
		engine.Run(ctx, sys.Generator, maxIterations, genChan)
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

func NewSierpinskiSystem(cfg config.TriangleConfig) *System[shared.Point] {
	gen := sierpinski.New(cfg)
	rendCfg := renderers.PointRendererConfig{
		MinPoint: shared.Point{X: 0, Y: 0},
		MaxPoint: shared.Point{X: cfg.EdgeLen, Y: cfg.EdgeLen},
	}
	var pr renderers.PointRenderer
	// pr = listpointrenderer.New()
	pr = ebitrenderer.New()
	// pr = tcellrenderer.New()
	pr.Initialize(rendCfg)

	out := &System[shared.Point]{
		Generator: gen,
		Renderer:  pr,
	}
	return out
}

func NewFernSystem(cfg config.FernConfig) *System[shared.Point] {
	gen := fern.New(cfg)
	rendCfg := renderers.PointRendererConfig{
		MinPoint: shared.Point{X: -1 * cfg.MaxWidth / 2, Y: 0},
		MaxPoint: shared.Point{X: cfg.MaxWidth / 2, Y: cfg.MaxHeight},
	}
	var pr renderers.PointRenderer
	pr = ebitrenderer.New()
	pr.Initialize(rendCfg)

	out := &System[shared.Point]{
		Generator: gen,
		Renderer:  pr,
	}
	return out
}

func NewNegTestSystem(cfg config.RectangleConfig) *System[shared.Point] {
	gen := negtest.New(cfg)
	rendCfg := renderers.PointRendererConfig{
		MinPoint: cfg.MinPoint,
		MaxPoint: cfg.MaxPoint,
	}
	var pr renderers.PointRenderer
	pr = ebitrenderer.New()
	pr.Initialize(rendCfg)

	out := &System[shared.Point]{
		Generator: gen,
		Renderer:  pr,
	}
	return out
}
