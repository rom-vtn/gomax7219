package gomax7219

import "strings"

//TODO make a better function?
// Returns a StaticGrid which renders the given pattern as a multiline string (separated by \n)
func NewRawGridFromPattern(pattern string) StaticGrid {
	lines := strings.Split(pattern, "\n")
	var content []byte
	for _, line := range lines {
		var current byte
		for i, char := range line {
			if strings.TrimSpace(string(char)) != "" {
				current += (1 << (7 - i))
			}
		}
		content = append(content, current)
	}
	return StaticGrid(content)
}

const TramRefString string = " XXXX   \n XXXXX  \nXXX  X X\nXXX  XX \nXXXXXX X\n XX  X  \n XX  X  \n XXXXX  \n XX  X  \n XX  X  \nXXXXXX X\nXXX  XX \nXXX  X X\n XXXXX  \n XXXX   \n X      "
const TrainRefString string = " XXXXX  \n XX  X  \n XX  X  \n XXXXX  \nXXX  XX \nXXX  X X\n XXXXX X\nXXX  X  \nXXX X   \n XXX    \n XX     "
const ClockRefString string = "  XXXXX \n X     X\n X     X\n X  XXXX\n X X   X\n XX    X\n  XXXXX "
