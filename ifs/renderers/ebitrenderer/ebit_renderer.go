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
	log.Println("FOOBAR")
	lastIdx := len(g.Points) - 1
	h := screen.Bounds().Dy()
	log.Printf("FOOBAR Draw lastidx: %d. lastdrawn: %d", lastIdx, g.LastDrawn)
	for g.LastDrawn < lastIdx {
		g.LastDrawn++
		curPt := g.Points[g.LastDrawn]
		log.Printf("FOOBAR. drawing point %v", curPt)
		screenX := curPt.X
		screenY := h - 1 - curPt.Y
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(screenX), float64(screenY))
		// screen.DrawImage(g.Dot, op)
	}
}

// Layout sets the logical size of the screen
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// Update updates internal game state.
// We are not using this function to do that. Instead, we are
// using the Draw() function of EbitRenderer, which is called by
// to our IFS system.
func (g *Game) Update() error {
	return nil
}

func (e *EbitRenderer) Initialize(cfg renderers.PointRendererConfig) {
	ebiten.SetTPS(1) // we do not use the Update function, so set this very low
	xSize := cfg.MaxPoint.X - cfg.MinPoint.X
	ySize := cfg.MaxPoint.Y - cfg.MinPoint.Y
	ebiten.SetWindowSize(xSize, ySize)
	ebiten.SetWindowTitle("Fun Times")
	go func() {
		// THIS DOES NOT WORK. ebitengine requires the main thread to do
		// its operations. But if we run RunGame in the main thread, it takes
		// over our program execution. : (
		// In order to use any GUI kind of renderer, will need to rework how
		// the IFS engine interacts with the renderer, since they all want some
		// sort of a loop on the main thread.
		if err := ebiten.RunGame(e.Game); err != nil {
			log.Fatalf("Error running game: %s", err)
		}
	}()
}

// Our Draw function that is called by our IFS framework is essentially the way we
// know there is something to render. We use it here to mutate the state of our
// Game object.
func (e *EbitRenderer) Draw(p shared.Point) {
	log.Printf("FOOBAR ebitrenderer draw %v", p)
	e.Game.Points = append(e.Game.Points, p)
}
