package listpointrenderer

import (
	"fmt"
	"time"

	"github.com/tdarci/go-nothings/ifs/renderers"
	"github.com/tdarci/go-nothings/ifs/shared"
)

const maxDraw = 100
const maxRunTime = time.Second * 3

type ListPointRenderer struct {
	drawCount int
	startTime time.Time
}

func New() *ListPointRenderer {
	return &ListPointRenderer{}
}

func (l *ListPointRenderer) Initialize(cfg renderers.PointRendererConfig) {
	fmt.Println("I am list-based renderer. I will show you all the points I get in list order.")
	fmt.Printf("Min: %v Max: %v\n", cfg.MinPoint, cfg.MaxPoint)

	l.startTime = time.Now()
}

func (l *ListPointRenderer) Draw(p shared.Point) {
	l.drawCount++
	fmt.Printf("* x: %8.2f • y: %8.2f\n", p.X, p.Y)
}

func (l *ListPointRenderer) Run() {
	for {
		time.Sleep(time.Millisecond * 10)
		if l.drawCount >= maxDraw {
			return
		}
		if time.Since(l.startTime) > maxRunTime {
			return
		}
	}
}
