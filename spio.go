package gomax7219

import (
	"fmt"
	"time"

	"github.com/fulr/spidev"
)

type SpiScreen struct {
	device         *spidev.SPIDevice
	cascadeCount   uint
	brightness     uint
	rotateCount    uint
	flipHorizontal bool
	flipVertical   bool
}

func (ss *SpiScreen) Close() {
	ss.device.Close()
}

func NewDeviceAndOpen(spibus, spidevice uint, cascadeCount uint, brightness uint, rotateCount uint, flipHorizontal, flipVertical bool) (*SpiScreen, error) {
	devstr := fmt.Sprintf("/dev/spidev%d.%d", spibus, spidevice)
	spi, err := spidev.NewSPIDevice(devstr)
	if err != nil {
		return nil, err
	}
	ss := SpiScreen{
		cascadeCount:   cascadeCount,
		brightness:     brightness,
		device:         spi,
		rotateCount:    rotateCount,
		flipHorizontal: flipHorizontal,
		flipVertical:   flipVertical,
	}
	scanLimitMessages := newControlMessageLine(MAX7219_REG_SCANLIMIT, 0x07, cascadeCount)
	encodingMessages := newControlMessageLine(MAX7219_REG_DECODEMODE, 0x00, cascadeCount) //0x00 = raw encoding, no 7 segment
	displayTestMessages := newControlMessageLine(MAX7219_REG_DISPLAYTEST, 0x00, cascadeCount)
	shutdownMessages := newControlMessageLine(MAX7219_REG_SHUTDOWN, 0x01, cascadeCount) //0x01 = no shutdown
	brightnessMessages := newControlMessageLine(MAX7219_REG_INTENSITY, byte(brightness*16), cascadeCount)
	messageLines := []messageLine{scanLimitMessages, encodingMessages, displayTestMessages, shutdownMessages, brightnessMessages}
	//send instructions
	for _, ml := range messageLines {
		err := ml.drawTo(&ss)
		if err != nil {
			return nil, err
		}
	}
	return &ss, nil
}

// shamelessly stolen from https://github.com/d2r2/go-max7219
type max7219Reg byte

const (
	MAX7219_REG_NOOP   max7219Reg = 0
	MAX7219_REG_DIGIT0            = iota
	MAX7219_REG_DIGIT1
	MAX7219_REG_DIGIT2
	MAX7219_REG_DIGIT3
	MAX7219_REG_DIGIT4
	MAX7219_REG_DIGIT5
	MAX7219_REG_DIGIT6
	MAX7219_REG_DIGIT7
	MAX7219_REG_DECODEMODE
	MAX7219_REG_INTENSITY
	MAX7219_REG_SCANLIMIT
	MAX7219_REG_SHUTDOWN
	MAX7219_REG_DISPLAYTEST = 0x0F
)

type square [8]byte

func (s square) rotateOnce() square {
	//convert [8]byte to [8][8]bool
	var boolSquare [8][8]bool
	for i, b := range s {
		boolSquare[i] = boolArrFromByte(b)
	}
	//90deg rotation
	var newBoolSquare [8][8]bool
	for i := range 8 {
		for j := range 8 {
			newBoolSquare[i][j] = boolSquare[j][7-i]
		}
	}
	//convert back to [8]byte
	var newSquare [8]byte
	for i := range 8 {
		newSquare[i] = byteFromBoolArr(newBoolSquare[i])
	}
	return newSquare
}

func boolArrFromByte(b byte) [8]bool {
	var boolArr [8]bool
	for i := range 8 {
		boolArr[i] = (1<<i)&b > 0
	}
	return boolArr
}
func byteFromBoolArr(boolArr [8]bool) byte {
	var b byte
	for i, val := range boolArr {
		if val {
			b += (1 << i)
		}
	}
	return b
}

type registerMessage [2]byte //2 bytes to send a message to a register
type messageLine []registerMessage

func newControlMessageLine(register, value byte, moduleCount uint) messageLine {
	message := registerMessage{register, value}
	var line messageLine
	for range moduleCount {
		line = append(line, message)
	}
	return line
}

func (sg StaticGrid) cutIntoSquares() []square {
	allBytes := []byte(sg)
	for len(allBytes)%8 != 0 {
		allBytes = append(allBytes, 0x00)
	}
	var squares []square
	for i := range len(allBytes) / 8 {
		squares = append(squares, square(allBytes[8*i:8*(i+1)]))
	}
	return squares
}

func (s square) getRegisterMessages() [8]registerMessage {
	var messages []registerMessage
	for regCount, content := range s {
		actualReg := MAX7219_REG_DIGIT0 + regCount
		message := [2]byte{byte(actualReg), content}
		messages = append(messages, message)
	}
	return [8]registerMessage(messages)
}

// const MIN_CLOCK_TIME = 20 * time.Microsecond //TODO get the exact time from the spec
const MIN_CLOCK_TIME = 80 * time.Microsecond //try slowing down the whole thing

func (sg StaticGrid) getMessageLines(rotateCount uint, flipHorizontal bool, flipVertical bool) []messageLine {
	//flip static grid contents
	if flipHorizontal {
		sg = sg.flippedAlongHorizontalAxis()
	}
	if flipVertical {
		sg = sg.flippedAlongVerticalAxis()
	}
	//rotate squares accordingly
	squares := sg.cutIntoSquares()
	for i := range squares {
		for range rotateCount {
			squares[i] = squares[i].rotateOnce()
		}
	}
	//should have dimensions [squareCount][registerInstruction]
	var allMessages [][8]registerMessage
	for _, s := range squares {
		allMessages = append(allMessages, s.getRegisterMessages())
	}
	var messageLines []messageLine
	for regLineIndex := range 8 {
		var messageLine []registerMessage
		for squareIndex := range squares {
			messageLine = append(messageLine, allMessages[squareIndex][regLineIndex])
		}
		messageLines = append(messageLines, messageLine)
	}
	return messageLines
}

func (ss *SpiScreen) Draw(r Renderer, delay time.Duration) error {
	for i := range r.GetFrameCount() {
		frame := r.Render(i)
		err := frame.drawTo(ss)
		if err != nil {
			return err
		}
		time.Sleep(delay)
	}
	return nil
}

func (ss *SpiScreen) Clear() error {
	empty := make(StaticGrid, ss.cascadeCount*8)
	return empty.drawTo(ss)
}

func (sg StaticGrid) drawTo(ss *SpiScreen) error {
	messageLines := sg.getMessageLines(ss.rotateCount, ss.flipHorizontal, ss.flipVertical)
	for _, messageLine := range messageLines {
		err := messageLine.drawTo(ss)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ml messageLine) drawTo(ss *SpiScreen) error {
	buf := make([]byte, 2*ss.cascadeCount)
	for i, regMsg := range ml {
		buf[2*i] = regMsg[0]
		buf[2*i+1] = regMsg[1]
	}
	_, err := ss.device.Xfer(buf)
	if err != nil {
		return err
	}
	return nil
}
