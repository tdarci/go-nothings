package sierpinski

import (
	"log"
	"math"
	"math/rand"

	"github.com/tdarci/go-nothings/ifs/config"
	"github.com/tdarci/go-nothings/ifs/shared"
)

type Sierpinski struct {
	config config.TriangleConfig
	state  State
}

func New(cfg config.TriangleConfig) *Sierpinski {
	return &Sierpinski{
		config: cfg,
	}
}

type State struct {
	Vertices []shared.Point
}

func (s *Sierpinski) Initialize() []shared.Point {

	s.state.Vertices = make([]shared.Point, s.config.NumVertices)
	var out []shared.Point
	switch s.config.NumVertices {
	case 3:
		out = s.prepareTriangle()
	case 4:
		out = s.prepareRectangle()
	default:
		log.Fatalf("Incorrect configuration for vertices. Must be between 3 and 4. %d were requested", s.config.NumVertices)
	}
	return out
}

func (s *Sierpinski) prepareRectangle() []shared.Point {
	legLen := s.config.EdgeLen
	vA := shared.Point{X: 0, Y: 0}
	vB := shared.Point{X: legLen, Y: 0}
	vC := shared.Point{X: 0, Y: legLen}
	vD := shared.Point{X: legLen, Y: legLen}

	s.state.Vertices[0] = vA
	s.state.Vertices[1] = vB
	s.state.Vertices[2] = vC
	s.state.Vertices[3] = vD

	// get a random point inside
	r1 := rand.Float64()
	r2 := rand.Float64()
	if r1+r2 > 1 {
		r1 = 1 - r1
		r2 = 1 - r2
	}

	initialDot := shared.Point{
		X: vA.X + r1*legLen,
		Y: vA.Y + r2*legLen,
	}

	out := append(s.state.Vertices, initialDot)
	return out
}

func (s *Sierpinski) prepareTriangle() []shared.Point {
	// todo: irregular triangles
	legLen := s.config.EdgeLen
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

	v := s.state.Vertices[rand.Intn(s.config.NumVertices)]

	out := shared.Point{
		X: (v.X + o.X) / 2,
		Y: (v.Y + o.Y) / 2,
	}
	return out
}
