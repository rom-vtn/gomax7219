package gomax7219

type blinkingGrid struct {
	innerContent  Renderer
	visibleFrames uint
	totalFrames   uint
}

// returns a Renderer which will show r for visibleFrames frames, and then show a blank grid until totalFrames
func NewBlinkingGrid(r Renderer, visibleFrames, totalFrames uint) Renderer {
	return blinkingGrid{
		innerContent:  r,
		visibleFrames: visibleFrames,
		totalFrames:   totalFrames,
	}
}

func (bg blinkingGrid) GetFrameCount() uint {
	return bg.totalFrames
}
func (bg blinkingGrid) GetWidth() uint {
	return bg.innerContent.GetWidth()
}
func (bg blinkingGrid) Render(frame uint) StaticGrid {
	frame = frame % bg.totalFrames
	if frame >= bg.visibleFrames {
		return make(StaticGrid, bg.GetWidth()) //empty screen if not in visible timespan
	}
	//otherwise return normal inner render
	return bg.innerContent.Render(frame)
}
