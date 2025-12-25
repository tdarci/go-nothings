package tcellrenderer

import (
	"github.com/tdarci/go-nothings/ifs/renderers"
	"github.com/tdarci/go-nothings/ifs/shared"
)

type TCellRenderer struct {

}

func New() *TCellRenderer {
	return &TCellRenderer{}
}

func (t *TCellRenderer) Initialize(cfg renderers.PointRendererConfig){

}

func (t *TCellRenderer) Draw (p shared.Point){
	
}
