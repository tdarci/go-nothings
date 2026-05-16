package config

import "github.com/tdarci/go-nothings/ifs/shared"

type Config struct {
	MaxIterations  int             `yaml:"MaxIterations"`
	DrawingDelayMs int             `yaml:"DrawingDelayMs"`
	Triangle       TriangleConfig  `yaml:"Triangle"`
	Carpet         CarpetConfig    `yaml:"Carpet"`
	Fern           FernConfig      `yaml:"Fern"`
	Rectangle      RectangleConfig `yaml:"Rectangle"`
	LSystem        LSystemConfig   `yaml:"LSystem"`
}

type TriangleConfig struct {
	EdgeLen     float64 `yaml:"EdgeLen"`
	NumVertices int     `yaml:"NumVertices"`
	Renderer    string  `yaml:"Renderer"`
}

type CarpetConfig struct {
	EdgeLen  float64 `yaml:"EdgeLen"`
	Renderer string  `yaml:"Renderer"`
}

type FernConfig struct {
	MaxWidth   float64        `yaml:"MaxWidth"`
	MaxHeight  float64        `yaml:"MaxHeight"`
	Thresholds FernThresholds `yaml:"Thresholds"`
}

type LSystemConfig struct {
	Axiom string        `yaml:"Axiom"`
	Angle int           `yaml:"Angle"`
	Rules []LSystemRule `yaml:"Rules"`
}

type LSystemRule struct {
	Match     string   `yaml:"Match"` // must be one-character long
	PreMatch  string `yaml:"PreMatch"`
	PostMatch string `yaml:"PostMatch"`
	RewriteAs string `yaml:"RewriteAs"`
}

// FernThresholds define the percentages at which each transformation "kicks in".
// Each time the generator generates its next point, it selects one of the four (StemBase, SmallerLeaflet,
// LeftBigLeaf, RightBigLeaf) to add to based on a random number between 0 and 100, and these thresholds.
// Note that they must get successively greater. And that the RightBigLeaf threshold is not defined, since
// it is everything above the threshold for LeftBigLeaf
type FernThresholds struct {
	StemBaseThresholdPct       int `yaml:"StemBaseThresholdPct"`       // canonical is 1
	SmallerLeafletThresholdPct int `yaml:"SmallerLeafletThresholdPct"` // canonical is 86
	LeftBigLeafletThresholdPct int `yaml:"LeftBigLeafletThresholdPct"` // canonical is 93
}

type RectangleConfig struct {
	MinPoint shared.Point `yaml:"MinPoint"`
	MaxPoint shared.Point `yaml:"MaxPoint"`
}
