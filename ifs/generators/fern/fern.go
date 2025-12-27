package fern

import (
	"math/rand"

	"github.com/tdarci/go-nothings/ifs/config"
	"github.com/tdarci/go-nothings/ifs/shared"
)

// scaleFactor is the multiplier we use when turning our calculations into
// points to plot.
// A standard Barnsley's Fern exhibits these limits:
// x ≈ [-2.1820 , 2.6558]
// y ≈ [ 0      , 9.9983]

// We scale to fit to our bounds. Imagining a rounding like this:
// x ≈ [-3 ,  3]
// y ≈ [ 0 , 10]

// let's scale based on the Y axis
func (f *Fern) scaleFactor() float64{
	return f.config.MaxHeight/11
}

type Fern struct {
	config config.FernConfig
}

func New(cfg config.FernConfig) *Fern {
	return &Fern{
		config: cfg,
	}
}

func (f *Fern) Initialize() []shared.Point {
	return []shared.Point{{X: 0, Y: 0}}
}

func (f *Fern) Next(o shared.Point) shared.Point {
	var newX, newY float64

	r := rand.Float32()
	oldY := o.Y/f.scaleFactor()
	oldX := o.X/f.scaleFactor()
	switch {
	case r < 0.01:
		newX = 0
		newY = 0.16 * oldY
	case r < 0.86:
		newX = (0.85*oldX + 0.04*oldY)
		newY = (-0.04*oldX + (0.85 * oldY) + 1.6)
	case r < 0.93:
		newX = (0.2*oldX - 0.26*oldY)
		newY = (0.23*oldX + (0.22 * oldY) + 1.6)
	default:
		newX = (-0.15*oldX + 0.28*oldY)
		newY = (0.26*oldX + (0.24 * oldY) + 0.44)
	}
	newX = newX * f.scaleFactor()
	newY = newY * f.scaleFactor()

	if newX > f.config.MaxWidth {
		newX = f.config.MaxWidth
	}
	if newY > f.config.MaxHeight {
		newY = f.config.MaxHeight
	}

	return shared.Point{X: newX, Y: newY}
}
