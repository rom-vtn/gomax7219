package gomax7219

//a grid composed of several other renderers
type concatenatedGrid struct {
	innerGrids []Renderer
	width      uint
	frameCount uint
}

// returns a Renderer which will concatenate the given Renderers (width = sum of widths, frame count = maximum frame count given)
func NewConcatenateGrid(innerGrids []Renderer) Renderer {
	var width uint
	var frameCount uint
	for _, r := range innerGrids {
		width += r.GetWidth()
		frameCount = max(frameCount, r.GetFrameCount())
	}
	return concatenatedGrid{
		innerGrids: innerGrids,
		width:      width,
		frameCount: frameCount,
	}
}
func (cg concatenatedGrid) Render(frame uint) StaticGrid {
	var contents StaticGrid
	for _, r := range cg.innerGrids {
		contents = append(contents, r.Render(frame)...)
	}
	return contents
}
func (cg concatenatedGrid) GetFrameCount() uint {
	return cg.frameCount
}
func (cg concatenatedGrid) GetWidth() uint {
	return cg.width
}
