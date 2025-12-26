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
	// Run initiates the UX and ends when the user
	// closes the window.
	Run()
}

