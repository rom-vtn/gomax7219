# go-max7219

This is a project I've been working on to spice up the use of MAX7219 displays on a Raspberry Pi SPI interface and have near full control over what gets displayed without having to spend too much time tinkering with the low level SPI interface.

# Inspiration
- https://github.com/d2r2/go-max7219 (most of the inspiration comes from there)
- https://github.com/rm-hull/luma.core/discussions/226 (Atari font with diacritics)
- https://github.com/fulr/spidev (for the SPI interface)

# Sample usage
This should display a clock with the system time.
```golang
func main() {
	err := stuffies()
	if err != nil {
		panic(err)
	}
}

const (
	SPI_BUS         uint = 0
	SPI_DEVICE      uint = 0
	CASCADE_COUNT   uint = 8
	BRIGHTNESS      uint = 0
	ROTATE_COUNT    uint = 3
	FLIP_HORIZONTAL bool = false
	FLIP_VERTICAL   bool = false
)

func stuffies() error {
	ss, err := NewDeviceAndOpen(SPI_BUS, SPI_DEVICE, CASCADE_COUNT, BRIGHTNESS, ROTATE_COUNT, FLIP_HORIZONTAL, FLIP_VERTICAL)
	if err != nil {
		return err
	}
	defer ss.Close()

	for range time.Tick(time.Second) {
		clockIcon := NewRawGridFromPattern(ClockRefString)
		timeText := time.Now().Format("15:04:05")
		timeRenderUntrimmed := NewStringTextRender(ATARI_FONT, timeText)
		timeRenderTrimmed := NewFitInsideGrid(timeRenderUntrimmed, 8*CASCADE_COUNT-clockIcon.GetWidth())
		concat := NewConcatenateGrid([]Renderer{clockIcon, timeRenderTrimmed})
		err = ss.Draw(concat, 10*time.Millisecond)
		if err != nil {
			return err
		}
	}

	return nil
}

```

# Fonts
Currently I copy-pasted the CP437 and the Atari fonts. The Atari font has support for some latin diacritics (tested in French and German). Feel free to send a PR if you want to suggest more fonts and extend language support.