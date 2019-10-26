package docker

import (
	"github.com/olshevskiy87/dockerate/colorer"
)

func GetContainerSizeColor(size int64) colorer.ColorCode {
	if size >= 0 && size < 500*Megabyte {
		return colorer.ColorDefault
	} else if size >= 500*Megabyte && size < Gigabyte {
		return colorer.ColorYellow
	}
	return colorer.ColorRed
}
