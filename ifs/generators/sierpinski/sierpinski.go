package sierpinski

import (
	"math"
	"math/rand"

	"github.com/tdarci/go-nothings/ifs/shared"
)

const numVertices = 3

type Sierpinski struct {
	config Config
	state  State
}

func New(cfg Config) *Sierpinski {
	return &Sierpinski{
		config: cfg,
	}
}

type Config struct {
	ContainerEdgeLength float64
}

type State struct {
	Vertices []shared.Point
}

func (s *Sierpinski) Initialize() []shared.Point {
	// todo: more vertices
	s.state.Vertices = make([]shared.Point, numVertices)
	// todo: irregular triangles
	legLen := s.config.ContainerEdgeLength
	vA := shared.Point{X: 0, Y: legLen}
	vB := shared.Point{X: legLen, Y: legLen}
	triangleHeight := math.Sqrt(math.Pow(legLen, 2) - math.Pow(legLen/2, 2))
	vC := shared.Point{
		X: legLen / 2,
		Y: legLen - triangleHeight,
	}

	s.state.Vertices[0] = vA
	s.state.Vertices[1] = vB
	s.state.Vertices[2] = vC

	// get a random point inside
	r1 := rand.Float64()
	r2 := rand.Float64()
	if r1+r2 > 1 {
		r1 = 1 - r1
		r2 = 1 - r2
	}

	initialDot := shared.Point{
		X: vA.X + r1*(vB.X-vA.X) + r2*(vC.X-vA.X),
		Y: vA.Y + r1*(vB.Y-vA.Y) + r2*(vC.Y-vA.Y),
	}

	out := append(s.state.Vertices, initialDot)
	return out
}

func (s *Sierpinski) Next(o shared.Point) shared.Point {

	v := s.state.Vertices[rand.Intn(3)]

	out := shared.Point{
		X: (v.X + o.X) / 2,
		Y: (v.Y + o.Y) / 2,
	}
	return out
}
