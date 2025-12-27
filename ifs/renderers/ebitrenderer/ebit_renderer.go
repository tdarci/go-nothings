package ebitrenderer

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tdarci/go-nothings/ifs/renderers"
	"github.com/tdarci/go-nothings/ifs/shared"
)

type EbitRenderer struct {
	Game *Game
}

func New() *EbitRenderer {
	return &EbitRenderer{
		Game: NewGame(),
	}
}

// Game is the type that implements the ebitengine Game interface
type Game struct {
	Points    []shared.Point
	Dot       *ebiten.Image
	LastDrawn int // Index into Points slice of last-drawn point
	min       shared.Point
	max       shared.Point
}

func NewGame() *Game {
	dot := ebiten.NewImage(1, 1)
	dot.Fill(color.RGBA{R: 255, G: 255, B: 0, A: 255})
	return &Game{
		Points:    make([]shared.Point, 0, 1000),
		LastDrawn: -1,
		Dot:       dot,
	}
}

// Draw is called when the screen refreshes. It is meant for
// rendering only, with no update of game state.
func (g *Game) Draw(screen *ebiten.Image) {
	lastIdx := len(g.Points) - 1

	screenHeight := screen.Bounds().Size().Y
	xOffset := -1.0 * g.min.X
	yOffset := -1.0 * g.min.Y

	for g.LastDrawn < lastIdx {
		g.LastDrawn++
		curPt := g.Points[g.LastDrawn]

		op := &ebiten.DrawImageOptions{}
		x := curPt.X + xOffset
		y := float64(screenHeight) - curPt.Y - yOffset
		// log.Printf("FOOBAR. [ebitrenderer.Game.Draw] drawing point %.2f, %.2f on screen at %.2f, %.2f", curPt.X, curPt.Y, x, y)
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.Dot, op)
	}
}

// Layout sets the logical size of the screen
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w := int(g.max.X - g.min.X)
	h := int(g.max.Y - g.min.Y)
	return w, h
}

// Update updates internal game state.
// We are not using this function to do that. Instead, we are
// using the Draw() function of EbitRenderer, which is called by
// to our IFS system.
func (g *Game) Update() error {
	return nil
}

func (e *EbitRenderer) Initialize(cfg renderers.PointRendererConfig) {
	e.Game.min = cfg.MinPoint
	e.Game.max = cfg.MaxPoint

	ebiten.SetTPS(1) // we do not use the Update function, so set this very low
	xSize, ySize := ebiten.Monitor().Size()
	ebiten.SetWindowSize(int(float32(xSize)*.9), int(float32(ySize)*.9))
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowTitle("Fun Times")
}

// Our Draw function that is called by our IFS framework is essentially the way we
// know there is something to render. We use it here to mutate the state of our
// Game object.
func (e *EbitRenderer) Draw(p shared.Point) {
	e.Game.Points = append(e.Game.Points, p)
}

func (e *EbitRenderer) Run() {
	if err := ebiten.RunGame(e.Game); err != nil {
		log.Fatalf("Error running game: %s", err)
	}
}
