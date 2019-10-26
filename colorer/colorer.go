package colorer

import (
	"fmt"
	"io"
)

type ColorCode int8

const (
	ColorDefault      ColorCode = -1
	ColorBlack        ColorCode = 0
	ColorRed          ColorCode = 1
	ColorGreen        ColorCode = 2
	ColorYellow       ColorCode = 3
	ColorBlue         ColorCode = 4
	ColorMagenta      ColorCode = 5
	ColorCyan         ColorCode = 6
	ColorLightGray    ColorCode = 7
	ColorDarkGray     ColorCode = 8
	ColorLightRed     ColorCode = 9
	ColorLightGreen   ColorCode = 10
	ColorLightYellow  ColorCode = 11
	ColorLightBlue    ColorCode = 12
	ColorLightMagenta ColorCode = 13
	ColorLightCyan    ColorCode = 14
	ColorWhite        ColorCode = 15
)

func Paint(code ColorCode, str string) string {
	if code == -1 {
		return fmt.Sprintf("<nofg>%s<reset>", str)
	}
	return fmt.Sprintf("<fg %d>%s<reset>", code, str)
}

func Paintf(code ColorCode, format string, a ...interface{}) string {
	return fmt.Sprintf(Paint(code, format), a...)
}

func Fpaintf(w io.Writer, code ColorCode, format string, a ...interface{}) (int, error) {
	p := Paintf(code, format, a...)
	return w.Write([]byte(p))
}
