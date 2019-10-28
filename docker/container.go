package docker

import (
	"github.com/olshevskiy87/dockerate/colorer"
	"github.com/olshevskiy87/dockerate/unit"
)

func GetContainerSizeColor(size int64) colorer.ColorCode {
	if size >= 0 && size < 500*unit.Megabyte {
		return colorer.ColorDefault
	} else if size >= 500*unit.Megabyte && size < unit.Gigabyte {
		return colorer.ColorYellow
	}
	return colorer.ColorRed
}
