package container

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/olshevskiy87/dockerate/colorer"
	"github.com/olshevskiy87/dockerate/unit"

	"github.com/docker/docker/api/types"
)

type portsStrings struct {
	display string
	sort    string
}

func (l *List) getSizeColor(size int64) colorer.ColorCode {
	if size >= 0 && size < 500*unit.Megabyte {
		return colorer.ColorDefault
	} else if size >= 500*unit.Megabyte && size < unit.Gigabyte {
		return colorer.ColorYellow
	}
	return colorer.ColorRed
}

func (l *List) getCreatedColor(interval int64) colorer.ColorCode {
	var createdColor = colorer.NoColor
	if l.Colorized {
		switch {
		case interval > unit.IntervalMonthSec:
			createdColor = colorer.ColorRed
		case interval > unit.IntervalWeekSec:
			createdColor = colorer.ColorYellow
		default:
			createdColor = colorer.ColorLightGreen
		}
	}
	return createdColor
}

func (l *List) makePortsItems(ports []types.Port) []portsStrings {
	var items = make([]portsStrings, len(ports))
	for i, portInfo := range ports {
		var hostIPPublicPort strings.Builder

		// IP
		if portInfo.IP != "" {
			hostIPPublicPort.WriteString(portInfo.IP)
		}
		// PublicPort
		if portInfo.PublicPort != 0 {
			if hostIPPublicPort.Len() > 0 {
				hostIPPublicPort.WriteString(":")
			}
			hostIPPublicPort.WriteString(strconv.Itoa(int(portInfo.PublicPort)))
		}

		var portLine, portLineSort strings.Builder
		if hostIPPublicPort.Len() > 0 {
			var color = colorer.NoColor
			if l.Colorized {
				color = colorer.ColorLightCyan
			}
			var hostIPPublicPortStr = fmt.Sprintf(
				"%s->",
				hostIPPublicPort.String(),
			)
			portLine.WriteString(colorer.Paint(color, hostIPPublicPortStr))
			portLineSort.WriteString(hostIPPublicPortStr)
		}

		// PrivatePort and Type (required)
		var privatePortStr = fmt.Sprintf(
			"%s/%s",
			strconv.Itoa(int(portInfo.PrivatePort)),
			portInfo.Type,
		)
		portLine.WriteString(privatePortStr)
		portLineSort.WriteString(privatePortStr)

		items[i] = portsStrings{display: portLine.String(), sort: portLineSort.String()}
	}
	return items
}
