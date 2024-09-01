package gomax7219

import "errors"

type sequenceGrid struct {
	contents       []Renderer
	frameDurations []uint
	totalFrames    uint
	width          uint
}

func NewSequenceGrid(contents []Renderer, frameDurations []uint) (Renderer, error) {
	if len(contents) != len(frameDurations) {
		return nil, errors.New("different amount of inner renderers and frame durations")
	}
	if len(contents) == 0 {
		return nil, errors.New("empty sequence grid")
	}

	var totalFrames, width uint
	for i := range contents {
		width = max(width, contents[i].GetWidth())
		totalFrames += frameDurations[i]
	}
	return sequenceGrid{
		contents:       contents,
		frameDurations: frameDurations,
		totalFrames:    totalFrames,
		width:          width,
	}, nil
}

func (sg sequenceGrid) GetFrameCount() uint {
	return sg.totalFrames
}
func (sg sequenceGrid) GetWidth() uint {
	return sg.width
}
func (sg sequenceGrid) Render(frame uint) StaticGrid {
	frame = frame % sg.totalFrames
	for i, duration := range sg.frameDurations {
		//keep going until we hit the frame
		frame -= duration
		if frame > 0 {
			continue
		}
		frame += duration //put into positive again
		grid := sg.contents[i].Render(frame).padTo(sg.width)
		return grid
	}
	return make(StaticGrid, sg.width) //should never happen
}
