package renderers

import "github.com/tdarci/go-nothings/ifs/shared"

type PointRendererConfig struct {
	MinPoint shared.Point
	MaxPoint shared.Point
}

type PointRenderer interface {
	Initialize(PointRendererConfig)
	Draw(shared.Point)
}

type Renderer[T any] interface{
	Draw(T)
}

