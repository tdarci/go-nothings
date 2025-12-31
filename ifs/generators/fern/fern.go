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
		// notes are from wikipedia: https://en.wikipedia.org/wiki/Barnsley_fern
	case r < float32(f.config.Thresholds.StemBaseThresholdPct)/100:
		// Stem base.
		// This coordinate transformation is chosen 1% of the time and just maps any point to a point in the first 
		// line segment at the base of the stem. In the iterative generation, it acts as a reset to the base of the 
		// stem. Crucially it does not reset exactly to (0,0) which allows it to fill in the base stem which is 
		// translated and serves as a kind of "kernel" from which all other sections of the fern are generated via 
		// the three other transformations.
		newX = 0
		newY = 0.16 * oldY
	case r < float32(f.config.Thresholds.SmallerLeafletThresholdPct)/100:
		// Successively smaller leaflets.
		// This transformation encodes the self-similarity relationship of the entire fern with the sub-structure 
		// which consists of the fern with the removal of the section which includes the bottom two leaves. In the 
		// matrix representation, it can be seen to be a slight clockwise rotation, scaled to be slightly smaller a
		// nd translated in the positive y direction. In the iterative generation, this transformation is applied with 
		// probability 85% and is intuitively responsible for the generation of the main stem, and the successive vertical 
		// generation of the leaves on either side of the stem from their "original" leaves at the base.
		newX = (0.85*oldX + 0.04*oldY)
		newY = (-0.04*oldX + (0.85 * oldY) + 1.6)
	case r < float32(f.config.Thresholds.LeftBigLeafletThresholdPct)/100:
		// Largest left-side leaflet.
		// This transformation encodes the self-similarity of the entire fern with the bottom left leaf. In the matrix 
		// representation, it is seen to be a near-90° counterclockwise rotation, scaled down to approximately 30% size 
		// with a translation in the positive y direction. In the iterative generation, it is applied with probability 
		// 7% and is intuitively responsible for the generation of the lower-left leaf.
		newX = (0.2*oldX - 0.26*oldY)
		newY = (0.23*oldX + (0.22 * oldY) + 1.6)
	default:
		// Largest right-side leaflet.
		// Similarly, this transformation encodes the self-similarity of the entire fern with the bottom right leaf. From 
		// its determinant it is easily seen to include a reflection and can be seen as a similar transformation as f3 
		// albeit with a reflection about the y-axis. In the iterative-generation, it is applied with probability 7% and 
		// is responsible for the generation of the bottom right leaf.
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
