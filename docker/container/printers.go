package container

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/olshevskiy87/dockerate/colorer"
	"github.com/olshevskiy87/dockerate/unit"

	"github.com/docker/docker/api/types"
	humanize "github.com/dustin/go-humanize"
)

type portsStrings struct {
	display string
	sort    string
}

func (l *List) fPrintHeader(w io.Writer) error {
	if l.OptQuiet {
		return nil
	}
	var header strings.Builder
	header.WriteString("CONTAINER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tPORTS\tNAMES")
	if l.OptSize {
		header.WriteString("\tSIZE")
	}
	header.WriteString("\n")
	var color = colorer.NoColor
	if l.Colorized {
		color = colorer.ColorLightBlue
	}
	_, err := colorer.Fpaintf(w, color, header.String())
	if err != nil {
		return fmt.Errorf("could not write header columns to output buffer: %v", err)
	}
	return nil
}

func (l *List) fPrintContainer(w io.Writer, container types.Container) error {
	if err := l.fPrintID(w, container.ID); err != nil {
		return fmt.Errorf("could not display container's field \"ID\": %v", err)
	}
	if l.OptQuiet {
		_, err := w.Write([]byte("\n"))
		return err
	}
	if err := l.fPrintImage(w, container.Image); err != nil {
		return fmt.Errorf("could not display container's field \"image\": %v", err)
	}
	if err := l.fPrintCommand(w, container.Command); err != nil {
		return fmt.Errorf("could not display container's field \"command\": %v", err)
	}
	if err := l.fPrintCreated(w, container.Created); err != nil {
		return fmt.Errorf("could not display container's field \"created\": %v", err)
	}
	if err := l.fPrintStatus(w, container.Status); err != nil {
		return fmt.Errorf("could not display container's field \"status\": %v", err)
	}
	if err := l.fPrintPorts(w, container.Ports); err != nil {
		return fmt.Errorf("could not display container's field \"ports\": %v", err)
	}
	if err := l.fPrintNames(w, container.Names); err != nil {
		return fmt.Errorf("could not display container's field \"names\": %v", err)
	}
	if l.OptSize {
		if err := l.fPrintSize(w, container.SizeRw, container.SizeRootFs); err != nil {
			return fmt.Errorf("could not display container's field \"size\": %v", err)
		}
	}
	_, err := w.Write([]byte("\n"))
	return err
}

func (l *List) fPrintID(w io.Writer, id string) error {
	var idOutput = id
	if !l.OptNoTrunc && len(idOutput) >= IDMinWidth {
		idOutput = idOutput[:IDMinWidth]
	}
	var color = colorer.NoColor
	if l.Colorized {
		color = colorer.ColorDarkGray
	}
	_, err := w.Write([]byte(colorer.Paint(color, idOutput)))
	return err
}

func (l *List) fPrintImage(w io.Writer, image string) error {
	var (
		outputImage strings.Builder
		imageItems  = strings.Split(image, ":")
	)
	if len(imageItems) > 0 {
		var color = colorer.NoColor
		if l.Colorized {
			color = colorer.ColorLightYellow
		}
		outputImage.WriteString(colorer.Paint(color, imageItems[0]))
	}
	if len(imageItems) > 1 {
		var imageTag = imageItems[1]
		if !l.OptNoTrunc && len(imageTag) > ImageTagMinWidth {
			imageTag = fmt.Sprintf("%s…", imageTag[:ImageTagMinWidth])
		}
		var color = colorer.NoColor
		if l.Colorized {
			color = colorer.ColorLightGreen
		}
		outputImage.WriteString(colorer.Paintf(color, ":%s", imageTag))
	}
	_, err := w.Write([]byte(fmt.Sprintf(
		"\t%s",
		outputImage.String(),
	)))
	return err
}

func (l *List) fPrintCommand(w io.Writer, command string) error {
	var commandOutput = command
	if !l.OptNoTrunc && len(command) > CommandMinWidth {
		commandOutput = fmt.Sprintf("%s…", command[:CommandMinWidth])
	}
	var color = colorer.NoColor
	if l.Colorized {
		color = colorer.ColorDarkGray
	}
	_, err := w.Write([]byte(colorer.Paintf(color, "\t\"%s\"", commandOutput)))
	return err
}

func (l *List) fPrintCreated(w io.Writer, created int64) error {
	var createdInterval = time.Now().Unix() - created
	var createdColor = colorer.NoColor
	if l.Colorized {
		switch {
		case createdInterval > unit.IntervalMonthSec:
			createdColor = colorer.ColorRed
		case createdInterval > unit.IntervalWeekSec:
			createdColor = colorer.ColorYellow
		default:
			createdColor = colorer.ColorLightGreen
		}
	}
	_, err := w.Write([]byte(colorer.Paintf(
		createdColor,
		"\t%s",
		humanize.Time(time.Unix(created, 0)),
	)))
	return err
}

func (l *List) fPrintStatus(w io.Writer, status string) error {
	var statusColor = colorer.NoColor
	if l.Colorized {
		if strings.HasPrefix(status, "Up") {
			statusColor = colorer.ColorLightGreen
		} else {
			statusColor = colorer.ColorDefault
		}
	}
	_, err := w.Write([]byte(colorer.Paintf(statusColor, "\t%s", status)))
	return err
}

func (l *List) fPrintPorts(w io.Writer, ports []types.Port) error {
	var portsItems = l.makePortsItems(ports)
	sort.Slice(
		portsItems,
		func(i, j int) bool {
			return portsItems[i].sort < portsItems[j].sort
		},
	)

	var portsOutput = make([]string, len(portsItems))
	for i, p := range portsItems {
		portsOutput[i] = p.display
	}

	_, err := w.Write([]byte(fmt.Sprintf(
		"\t%s",
		strings.Join(portsOutput, ", "),
	)))
	return err
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

func (l *List) fPrintNames(w io.Writer, names []string) error {
	var namesOutput = make([]string, len(names))
	for i, name := range names {
		namesOutput[i] = strings.TrimLeft(name, "/")
	}
	var color = colorer.NoColor
	if l.Colorized {
		color = colorer.ColorDefault
	}
	_, err := w.Write([]byte(colorer.Paintf(
		color,
		"\t%s",
		strings.Join(namesOutput, ", "),
	)))
	return err
}

func (l *List) fPrintSize(w io.Writer, sizeRw int64, sizeRootFs int64) error {
	var (
		colorSizeRw     = colorer.NoColor
		colorSizeRootFs = colorer.NoColor
	)
	if l.Colorized {
		colorSizeRw = l.getSizeColor(sizeRw)
		colorSizeRootFs = l.getSizeColor(sizeRw)
	}
	_, err := w.Write([]byte(fmt.Sprintf(
		"\t%s (%s)",
		colorer.Paint(
			colorSizeRw,
			humanize.Bytes(uint64(sizeRw)),
		),
		colorer.Paintf(
			colorSizeRootFs,
			"virtual %s",
			humanize.Bytes(uint64(sizeRootFs)),
		),
	)))
	return err
}

func (l *List) getSizeColor(size int64) colorer.ColorCode {
	if size >= 0 && size < 500*unit.Megabyte {
		return colorer.ColorDefault
	} else if size >= 500*unit.Megabyte && size < unit.Gigabyte {
		return colorer.ColorYellow
	}
	return colorer.ColorRed
}
