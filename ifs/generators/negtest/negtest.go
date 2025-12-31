package negtest

import (
	"math/rand/v2"

	"github.com/tdarci/go-nothings/ifs/config"
	"github.com/tdarci/go-nothings/ifs/shared"
)

const (
	edgeLeft   = 0
	edgeTop    = 1
	edgeRight  = 2
	edgeBottom = 3
)

type NegTest struct {
	config  config.RectangleConfig
	left    float64
	right   float64
	top     float64
	bottom  float64
	width   float64
	height  float64
	curEdge int
}

func New(cfg config.RectangleConfig) *NegTest {

	out := &NegTest{
		config: cfg,
	}

	scale := 0.8
	out.width = ((cfg.MaxPoint.X - cfg.MinPoint.X) * scale)
	xBorder := ((cfg.MaxPoint.X - cfg.MinPoint.X) - out.width) / 2
	out.left = cfg.MinPoint.X + xBorder
	out.right = cfg.MaxPoint.X - xBorder

	out.height = ((cfg.MaxPoint.Y - cfg.MinPoint.Y) * scale)
	yBorder := ((cfg.MaxPoint.Y - cfg.MinPoint.Y) - out.height) / 2
	out.bottom = cfg.MinPoint.Y + yBorder
	out.top = cfg.MaxPoint.Y - yBorder

	return out
}

func (n *NegTest) Initialize() []shared.Point {
	return nil
}

func (n *NegTest) Next(o shared.Point) shared.Point {

	out := shared.Point{}

	n.curEdge = (n.curEdge + 1) % 4
	r := rand.Float64()

	switch n.curEdge {
	case edgeLeft:
		out.X = n.left + (r * 10.0)
		out.Y = n.bottom + (r * n.height)
	case edgeRight:
		out.X = n.right - (r * 10.0)
		out.Y = n.bottom + (r * n.height)
	case edgeTop:
		out.X = n.left + (r * n.width)
		out.Y = n.top - (r * 10)
	case edgeBottom:
		out.X = n.left + (r * n.width)
		out.Y = n.bottom + (r * 10)
	default:
	}

	return out
}
