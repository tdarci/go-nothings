package harness

import (
	"sync"

	"github.com/tdarci/go-nothings/ifs/engine"
	"github.com/tdarci/go-nothings/ifs/generators/sierpinski"
	"github.com/tdarci/go-nothings/ifs/renderers"
	listpointrenderer "github.com/tdarci/go-nothings/ifs/renderers/ListPointRenderer"
	"github.com/tdarci/go-nothings/ifs/shared"
	"github.com/tdarci/go-nothings/ifs/utils/channelsplitter"
)

type System[T any] struct {
	Generator engine.Generator[T]
	Renderers []renderers.Renderer[T]
}

func RunSierpinski() {
	Run(NewSierpinskiSystem())
}

func Run[T any](sys *System[T]) {
	var wg sync.WaitGroup
	genChan := make(chan T, 50)
	channels := channelsplitter.Split(genChan, len(sys.Renderers), 50)
	for rIdx, r := range sys.Renderers {
		wg.Add(1)

		go func (i int, rd renderers.Renderer[T])  {
			for val := range channels[i] {
				r.Draw(val)
			}
			wg.Done()
		}(rIdx, r)
	}

	wg.Go(func ()  {
		engine.Run(sys.Generator, genChan)
		close(genChan)
	})

	wg.Wait()
}

func NewSierpinskiSystem() *System[shared.Point] {
	cfg := sierpinski.Config{
		MinPoint: shared.Point{X:0, Y:0},
		MaxPoint: shared.Point{X:100, Y:100},
	}
	gen := sierpinski.New(cfg)
	rendCfg := renderers.PointRendererConfig{
		MinPoint: cfg.MinPoint,
		MaxPoint: cfg.MaxPoint,
	}
	lpr := listpointrenderer.New()
	lpr.Initialize(rendCfg)
	rs := []renderers.Renderer[shared.Point]{lpr}

	out := &System[shared.Point]{
		Generator: gen,
		Renderers: rs,
	}
	return out
}