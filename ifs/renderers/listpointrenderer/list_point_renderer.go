package listpointrenderer

import (
	"fmt"

	"github.com/tdarci/go-nothings/ifs/renderers"
	"github.com/tdarci/go-nothings/ifs/shared"
)

type ListPointRenderer struct {
}

func New() *ListPointRenderer {
	out := &ListPointRenderer{}
	return out
}

func (l *ListPointRenderer) Initialize(cfg renderers.PointRendererConfig) {
	fmt.Println("I am list-based renderer. I will show you all the points I get in list order.")
	fmt.Printf("Min: %v Max: %v\n", cfg.MinPoint, cfg.MaxPoint)
}

func (l *ListPointRenderer) Draw(p shared.Point) {
	fmt.Printf("* x: %5d • y: %5d\n", p.X, p.Y)
}
