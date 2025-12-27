package config

import "github.com/tdarci/go-nothings/ifs/shared"

type Config struct {
	MaxIterations  int             `yaml:"MaxIterations"`
	DrawingDelayMs int             `yaml:"DrawingDelayMs"`
	Triangle       TriangleConfig  `yaml:"Triangle"`
	Fern           FernConfig      `yaml:"Fern"`
	Rectangle      RectangleConfig `yaml:"Rectangle"`
}

type TriangleConfig struct {
	EdgeLen float64 `yaml:"EdgeLen"`
	Renderer string `yaml:"Renderer"`
}

type FernConfig struct {
	MaxWidth  float64 `yaml:"MaxWidth"`
	MaxHeight float64 `yaml:"MaxHeight"`
}

type RectangleConfig struct {
	MinPoint shared.Point `yaml:"MinPoint"`
	MaxPoint shared.Point `yaml:"MaxPoint"`
}
