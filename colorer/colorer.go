package colorer

import (
	"fmt"
	"io"
)

type ColorCode uint8

const (
	ColorBlack ColorCode = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorLightGray
	ColorDarkGray
	ColorLightRed
	ColorLightGreen
	ColorLightYellow
	ColorLightBlue
	ColorLightMagenta
	ColorLightCyan
	ColorWhite
)

func Paint(code ColorCode, str string) string {
	return fmt.Sprintf("<fg %d>%s<reset>", code, str)
}

func Paintf(code ColorCode, format string, a ...interface{}) string {
	return fmt.Sprintf(Paint(code, format), a...)
}

func Fpaintf(w io.Writer, code ColorCode, format string, a ...interface{}) (int, error) {
	p := Paintf(code, format, a...)
	return w.Write([]byte(p))
}
