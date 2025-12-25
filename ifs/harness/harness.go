package harness

import (
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

	wg.Go(func() {
		tckr := time.NewTicker(time.Millisecond * 10)
		defer tckr.Stop()
		for val := range genChan {
			<- tckr.C
			sys.Renderer.Draw(val)
		}
	})

	wg.Go(func() {
		engine.Run(sys.Generator, iterations, genChan)
		close(genChan)
	})

	wg.Wait()
	time.Sleep(time.Second * 3)
	sys.Renderer.Shutdown()
}

func NewSierpinskiSystem() *System[shared.Point] {
	cfg := sierpinski.Config{
		ContainerEdgeLength: 200,
	}
	gen := sierpinski.New(cfg)
	rendCfg := renderers.PointRendererConfig{
		MinPoint: shared.Point{X:0, Y:0},
		MaxPoint: shared.Point{X: cfg.ContainerEdgeLength, Y:cfg.ContainerEdgeLength},
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
