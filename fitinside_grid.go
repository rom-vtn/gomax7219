package gomax7219

// NewFitInsideGrid returns a Renderer that shows r's content while ensuring it fits in the given size, padding (and centering) if too small and scrolling if too large
func NewFitInsideGrid(r Renderer, fitToSize uint) Renderer {
	widthOvershoot := int(fitToSize) - int(r.GetWidth())
	if widthOvershoot < 0 { //if too big, scroll
		return NewScrollingGrid(r, fitToSize)
	} else if widthOvershoot == 0 { //if size matches, return as raw
		return r
	} else { //if too small, center
		startPaddingSize := widthOvershoot / 2
		endPaddingSize := startPaddingSize
		if widthOvershoot%2 == 1 { //add one more if odd
			endPaddingSize++
		}
		startPadding := make(StaticGrid, startPaddingSize)
		endPadding := make(StaticGrid, endPaddingSize)
		return NewConcatenateGrid([]Renderer{startPadding, r, endPadding})
	}
}
