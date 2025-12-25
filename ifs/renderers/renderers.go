package renderers

import "github.com/tdarci/go-nothings/ifs/shared"

type PointRendererConfig struct {
	MinPoint shared.Point
	MaxPoint shared.Point
}

type PointRenderer interface {
	Renderer[shared.Point]

	Initialize(PointRendererConfig)
}

type Renderer[T any] interface{
	Draw(T)
	Shutdown()
}

