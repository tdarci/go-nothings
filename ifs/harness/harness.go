package harness

import (
	"sync"

	"github.com/tdarci/go-nothings/ifs/engine"
	"github.com/tdarci/go-nothings/ifs/generators/sierpinski"
	"github.com/tdarci/go-nothings/ifs/renderers"

	// "github.com/tdarci/go-nothings/ifs/renderers/ebitrenderer"
	"github.com/tdarci/go-nothings/ifs/renderers/listpointrenderer"
	"github.com/tdarci/go-nothings/ifs/shared"
)

type System[T any] struct {
	Generator engine.Generator[T]
	Renderer  renderers.Renderer[T]
}

func RunSierpinski() {
	Run(NewSierpinskiSystem())
}

func Run[T any](sys *System[T]) {
	var wg sync.WaitGroup
	genChan := make(chan T, 50)

	wg.Go(func() {
		for val := range genChan {
			// log.Printf("FOOBAR: [harness.Run] rendering point: %v", val)
			sys.Renderer.Draw(val)
		}
	})

	wg.Go(func() {
		engine.Run(sys.Generator, genChan)
		close(genChan)
	})

	wg.Wait()
}

func NewSierpinskiSystem() *System[shared.Point] {
	cfg := sierpinski.Config{
		MinPoint: shared.Point{X: 0, Y: 0},
		MaxPoint: shared.Point{X: 500, Y: 500},
	}
	gen := sierpinski.New(cfg)
	rendCfg := renderers.PointRendererConfig{
		MinPoint: cfg.MinPoint,
		MaxPoint: cfg.MaxPoint,
	}
	var pr renderers.PointRenderer
	pr = listpointrenderer.New()
	// pr = ebitrenderer.New()
	pr.Initialize(rendCfg)

	out := &System[shared.Point]{
		Generator: gen,
		Renderer:  pr,
	}
	return out
}
