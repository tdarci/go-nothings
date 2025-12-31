package carpet

import (
	"log"
	"math/rand/v2"

	"github.com/tdarci/go-nothings/ifs/config"
	"github.com/tdarci/go-nothings/ifs/shared"
)

// Carpet generates the Sierpinski Carpet.
type Carpet struct {
	config config.CarpetConfig
}

func New(cfg config.CarpetConfig) *Carpet {
	out := &Carpet{
		config: cfg,
	}
	return out
}

func (c *Carpet) Initialize() []shared.Point {
	// get point inside the square
	p := shared.Point{X: rand.Float64() * c.config.EdgeLen, Y: rand.Float64() * c.config.EdgeLen}

	// generate 100 points to throw away
	for range 100 {
		p = c.baseNext(p)
	}

	// generate one more, properly scaled
	p = c.Next(p)

	// send back points to draw
	out := []shared.Point{p}

	return out
}

func (c *Carpet) scaleFactor() float64 {
	return c.config.EdgeLen
}

func (c *Carpet) Next(in shared.Point) shared.Point {
	scale := c.scaleFactor()
	in.X = in.X / scale
	in.Y = in.Y / scale
	out := c.baseNext(in)
	out.X = out.X * scale
	out.Y = out.Y * scale
	// log.Printf("[carpet.Next] out: %+v", out)
	return out
}

func (c *Carpet) baseNext(in shared.Point) shared.Point {

	out := shared.Point{}
	var bucket int

	r := rand.Float32()
	switch {
	case r < .125:
		bucket = 1
		out.X = in.X / 3.0
		out.Y = in.Y / 3.0
	case r < .250:
		bucket = 2
		out.X = (in.X / 3) + (1.0 / 3.0)
		out.Y = in.Y / 3.0
	case r < .375:
		bucket = 3
		out.X = (in.X / 3.0) + (2.0 / 3.0)
		out.Y = in.Y / 3.0
	case r < .500:
		bucket = 4
		out.X = in.X / 3.0
		out.Y = (in.Y / 3.0) + (1.0 / 3.0)
	case r < .625:
		bucket = 5
		out.X = (in.X / 3.0) + (2.0 / 3.0)
		out.Y = (in.Y / 3.0) + (1.0 / 3.0)
	case r < .750:
		bucket = 6
		out.X = in.X / 3.0
		out.Y = (in.Y / 3.0) + (2.0 / 3.0)
	case r < .875:
		bucket = 7
		out.X = (in.X / 3) + (1.0 / 3.0)
		out.Y = (in.Y / 3.0) + (2.0 / 3.0)
	default:
		bucket = 8
		out.X = (in.X / 3.0) + (2.0 / 3.0)
		out.Y = (in.Y / 3.0) + (2.0 / 3.0)
	}

	if bucket == 0 {
		log.Fatalf("[carpet.baseNext] did not set bucket. hrm.")
	}

	// log.Printf("[carpet.baseNext] bucket: %d. in: %+v out: %+v", bucket, in, out)

	return out
}
