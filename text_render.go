package gomax7219

// NewBytesTextRender renders the given code points (byte per byte) using the given code page.
// This requires all codepoints to fit in a signle byte (not compatible with >255 codepoints)
func NewBytesTextRender(codePage [][]byte, codePoints []byte) StaticGrid {
	var render StaticGrid
	for _, codePoint := range codePoints {
		render = append(render, codePage[codePoint]...)
	}
	return render
}

// NewRunesTextRender renders text as a series of runes, but allows runes to be >=256
func NewRunesTextRender(codePage [][]byte, codePoints []rune) StaticGrid {
	var render StaticGrid
	for _, codePoint := range codePoints {
		render = append(render, codePage[codePoint]...)
	}
	return render
}

// NewStringTextRender renders the given string (converting it to runes, then to bytes) using the given code page
func NewStringTextRender(codePage [][]byte, content string) StaticGrid {
	//convert runes to bytes one by one
	needsTruncating := false
	var asBytes []byte
	for _, r := range content {
		if len(codePage) <= int(r) {
			needsTruncating = true
		}
		asBytes = append(asBytes, byte(r))
	}

	if needsTruncating {
		return NewBytesTextRender(codePage, asBytes)
	}
	return NewRunesTextRender(codePage, []rune(content))
}
