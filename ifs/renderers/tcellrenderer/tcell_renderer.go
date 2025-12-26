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
	origScreenH = origScreenH * 2 // we only get to use half the height
	screenW := origScreenW
	screenH := origScreenH
	logW := r.max.X - r.min.X
	logH := r.max.Y - r.min.Y
	logAspect := float32(logW) / float32(logH)
	screenAspect := float32(screenW) / float32(screenH)

	switch {
	case screenAspect > logAspect:
		screenW = int(float32(screenW) * logAspect / screenAspect)
	case screenAspect < logAspect:
		screenH = int(float32(screenH) * logAspect / screenAspect)
	}

	x := ((p.X - r.min.X) * screenW / logW) + ((origScreenW - screenW)/2)
	y := (((p.Y - r.min.Y) * screenH / logH) + ((origScreenH - screenH)/2)) / 2 // we divide by 2 because our "pixels" are characters, which are rougly 2:1 h:w

	if x < 0 || y < 0 {
		return
	}

	log.Printf("FOOBAR [tcellrenderer.Draw] drawing at %d, %d (from %v)", x, y, p)
	r.screen.SetContent(x, y, '·', nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
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
