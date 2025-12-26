package fern

import (
	"math/rand"

	"github.com/tdarci/go-nothings/ifs/shared"
)

const scaleFactor = 100

type Config struct {
	MaxWidth  float64
	MaxHeight float64
}

type Fern struct {
	config Config
}

func New(cfg Config) *Fern {
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
	oldY := o.Y/scaleFactor
	oldX := o.X/scaleFactor
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
	newX = newX * scaleFactor
	newY = newY * scaleFactor

	if newX > f.config.MaxWidth {
		newX = f.config.MaxWidth
	}
	if newY > f.config.MaxHeight {
		newY = f.config.MaxHeight
	}

	if newX < 0 {
		newX = 0
	}
	return shared.Point{X: newX, Y: newY}
}
