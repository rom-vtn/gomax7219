package gomax7219

import (
	"math/bits"
	"slices"
)

// a Renderer is "something that can be displayed", consisting of 1 or more frames and a given width
type Renderer interface {
	Render(frame uint) StaticGrid //renders on a given frame
	GetFrameCount() uint          //gives the total number of frames
	GetWidth() uint               //number of columns the content takes on the screen
}

// StaticGrid represents static raw content; its width is its content's width and its frame count is 1
type StaticGrid []byte

func (sg StaticGrid) Render(frame uint) StaticGrid {
	return sg
}
func (sg StaticGrid) GetFrameCount() uint {
	return 1
}
func (sg StaticGrid) GetWidth() uint {
	return uint(len(sg))
}

// pad with empty bytes up to desired length
func (sg StaticGrid) padTo(length uint) StaticGrid {
	for len(sg) < int(length) {
		sg = append(sg, 0x00)
	}
	return sg
}

func (sg StaticGrid) flippedAlongVerticalAxis() StaticGrid {
	slices.Reverse(sg)
	return sg
}

func (sg StaticGrid) flippedAlongHorizontalAxis() StaticGrid {
	for i := range sg {
		sg[i] = bits.Reverse8(sg[i])
	}
	return sg
}
