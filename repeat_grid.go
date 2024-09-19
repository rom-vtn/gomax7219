package gomax7219

type repeatGrid struct {
	innerContent Renderer
	repeatCount  uint
}

// NewRepeatGrid returns a Renderer which behaved identically to r, but with an extended frame count to repeat it the given amount of times
func NewRepeatGrid(r Renderer, repeatCount uint) Renderer {
	return &repeatGrid{
		innerContent: r,
		repeatCount:  repeatCount,
	}
}

func (rg *repeatGrid) Render(frame uint) StaticGrid {
	return rg.innerContent.Render(frame)
}

func (rg *repeatGrid) GetWidth() uint {
	return rg.innerContent.GetWidth()
}

func (rg *repeatGrid) GetFrameCount() uint {
	return rg.repeatCount * rg.innerContent.GetFrameCount()
}
