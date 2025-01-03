package gomax7219

import "errors"

type verticalSequenceGrid struct {
	contents       []Renderer
	frameDurations []uint
	totalFrames    uint
	width          uint
}

const VERTICAL_TRANSITION_FRAMECOUNT = 16

// NewVerticalSequenceGrid returns a Renderer that displays the given Renderers for the given frame durations with the appropriate vertical scroll transitions.
func NewVerticalSequenceGrid(contents []Renderer, frameDurations []uint) (Renderer, error) {
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
	totalFrames += uint(len(contents)) * VERTICAL_TRANSITION_FRAMECOUNT
	return verticalSequenceGrid{
		contents:       contents,
		frameDurations: frameDurations,
		totalFrames:    totalFrames,
		width:          width,
	}, nil
}

func (sg verticalSequenceGrid) GetFrameCount() uint {
	return sg.totalFrames
}
func (sg verticalSequenceGrid) GetWidth() uint {
	return sg.width
}
func (sg verticalSequenceGrid) Render(frame uint) StaticGrid {
	signedFrame := int(frame % sg.totalFrames)
	for i, duration := range sg.frameDurations {
		//subtract and keep going until we hit the frame
		signedFrame -= int(duration) + VERTICAL_TRANSITION_FRAMECOUNT
		if signedFrame > 0 {
			continue
		}
		if signedFrame > -VERTICAL_TRANSITION_FRAMECOUNT { //if in between renderers, change between last frame of current and 1st frame of next
			prevGrid := sg.contents[i].Render(sg.contents[i].GetFrameCount() - 1).padTo(sg.width)
			nextGrid := sg.contents[(i+1)%len(sg.contents)].Render(0).padTo(sg.width)
			var result []byte
			for pos := range sg.width {
				//note: LSB is up and MSB is down, and we're scrolling downwards
				currentByte := (int(prevGrid[pos])*256 + int(nextGrid[pos])) >> (-signedFrame * 8 / VERTICAL_TRANSITION_FRAMECOUNT)
				result = append(result, byte(currentByte))
			}
			return StaticGrid(result)
		}
		//if 100% inside current renderer
		signedFrame += int(duration) + VERTICAL_TRANSITION_FRAMECOUNT
		grid := sg.contents[i].Render(uint(signedFrame)).padTo(sg.width)
		return grid
	}
	panic(errors.New("uh oh, this isn't really expected")) //should never happen
}
