package gomax7219

// NewBytesTextRender renders the given code points (byte per byte) using the given code page
func NewBytesTextRender(codePage [][]byte, codePoints []byte) StaticGrid {
	var render StaticGrid
	for _, codePoint := range codePoints {
		render = append(render, codePage[codePoint]...)
	}
	return render
}

// NewStringTextRender renders the given string (converting it to runes, then to bytes) using the given code page
func NewStringTextRender(codePage [][]byte, content string) StaticGrid {
	//convert runes to bytes one by one
	var asBytes []byte
	for _, r := range content {
		asBytes = append(asBytes, byte(r))
	}

	return NewBytesTextRender(codePage, asBytes)
}
