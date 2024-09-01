package gomax7219

type scrollingGrid struct {
	innerContent Renderer
	width        uint
	frameCount   uint
}

// returns a Renderer which scrolls the content of r on the given width
func NewScrollingGrid(r Renderer, width uint) Renderer {
	frameCount := width + r.GetWidth() + 1 //scroll in, content, scroll out (fully, hence the +1)
	return scrollingGrid{
		innerContent: r,
		width:        width,
		frameCount:   frameCount,
	}
}

func (sg scrollingGrid) GetWidth() uint {
	return sg.width
}

func (sg scrollingGrid) GetFrameCount() uint {
	return sg.frameCount
}

func (sg scrollingGrid) Render(frame uint) StaticGrid {
	frame = frame % sg.frameCount
	innerRender := sg.innerContent.Render(frame)
	widthPad := make(StaticGrid, sg.width)
	totalView := append(widthPad, innerRender...)
	totalView = append(totalView, widthPad...)
	visiblePart := totalView[frame : frame+sg.width]
	return visiblePart
}
