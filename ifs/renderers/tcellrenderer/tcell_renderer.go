package tcellrenderer

import (
	"log"

	"github.com/gdamore/tcell/v3"

	"github.com/tdarci/go-nothings/ifs/renderers"
	"github.com/tdarci/go-nothings/ifs/shared"
)

type TCellRenderer struct {
	screen tcell.Screen
	min    shared.Point
	max    shared.Point
}

func New() *TCellRenderer {
	return &TCellRenderer{}
}

func (r *TCellRenderer) Initialize(cfg renderers.PointRendererConfig) {
	r.min = cfg.MinPoint
	r.max = cfg.MaxPoint

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("error creating screen: %s", err)
	}

	if err := s.Init(); err != nil {
		log.Fatalf("error initializing screen: %s", err)
	}

	s.Clear()
	s.Show()

	// this would get us mouse
	// s.EnableMouse()

	r.screen = s

}

func (r *TCellRenderer) Draw(p shared.Point) {
	if r.screen == nil {
		return
	}

	origScreenW, origScreenH := r.screen.Size()
	log.Printf("[tcellrenderer.Draw] screen width, height: %d, %d", origScreenW, origScreenH)
	origScreenH = origScreenH * 2 // each vertical "pixel" counts for 2
	
	screenW := origScreenW
	screenH := origScreenH
	screenAspect := float32(screenW) / float32(screenH)
	
	logW := r.max.X - r.min.X
	logH := r.max.Y - r.min.Y
	logAspect := float32(logW) / float32(logH)
	
	switch {
	case screenAspect > logAspect:
		screenW = int(float32(screenW) * logAspect / screenAspect)
	case screenAspect < logAspect:
		screenH = int(float32(screenH) * logAspect / screenAspect)
	}

	xOffset := -1.0 * r.min.X 
	yOffset := -1.0 * r.min.Y  

	x := p.X + xOffset
	y := p.Y + yOffset

	x = (x * float64(screenW) / logW) + ((float64(origScreenW - screenW)/2))
	y = (y * float64(screenH) / logH) + ((float64(origScreenH - screenH)/2))

	// y is more complicated
	// first, we have to subtract from our screen height since the coordinates are flipped.
	// then we have to divide by 2, since each "pixel" is 2 high
	y = (float64(screenH) - y) / 2

	if x < 0 || y < 0 {
	log.Printf("FOOBAR [tcellrenderer.Draw] NOT drawing at %6.2f, %6.2f (from %v)", x, y, p)
		return
	}

	log.Printf("FOOBAR [tcellrenderer.Draw] drawing at %6.2f, %6.2f (from %v)", x, y, p)
	r.screen.SetContent(int(x), int(y), '·', nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	r.screen.Show()
}

func (r *TCellRenderer) Run() {
	defer r.shutdown()

	ch := r.screen.EventQ()
	for ev := range ch {
		switch ev.(type){
		case *tcell.EventKey:
			// end when user presses a key
			return
		}
	}
}

func (r *TCellRenderer) shutdown() {
	if r.screen != nil {
		r.screen.Fini()
	}
}
