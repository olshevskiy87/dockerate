package container

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/olshevskiy87/dockerate/colorer"

	"github.com/docker/docker/api/types"
	humanize "github.com/dustin/go-humanize"
)

func (l *List) fPrintHeader() error {
	if l.OptQuiet {
		return nil
	}
	var header strings.Builder
	if len(l.Columns) != 0 {
		header.WriteString(strings.Join(l.Columns, "\t"))
	}
	header.WriteString("\n")
	var color = colorer.NoColor
	if l.Colorized {
		color = colorer.ColorLightBlue
	}
	_, err := colorer.Fpaintf(l.writer, color, header.String())
	if err != nil {
		return fmt.Errorf("could not write header columns to output buffer: %v", err)
	}
	return nil
}

func (l *List) fPrintContainer(container types.Container) error {
	if l.OptQuiet {
		if err := l.fPrintID(container.ID); err != nil {
			return fmt.Errorf("could not display container's field \"ID\": %v", err)
		}
		return l.write([]byte("\n"))
	}
	for _, column := range l.Columns {
		if l.isColumnSet(column) {
			if err := l.fPrintColumn(container, column); err != nil {
				return fmt.Errorf("could not display column \"%s\"", column)
			}
		}
		if l.write([]byte("\t")) != nil {
			return fmt.Errorf("could not display column delimiter (tab)")
		}
	}
	return l.write([]byte("\n"))
}

func (l *List) fPrintColumn(container types.Container, column string) error {
	switch {
	case column == ContainerIDColumnName:
		if err := l.fPrintID(container.ID); err != nil {
			return fmt.Errorf("could not display container's field \"%s\": %v", ContainerIDColumnName, err)
		}
	case column == ImageColumnName:
		if err := l.fPrintImage(container.Image); err != nil {
			return fmt.Errorf("could not display container's field \"%s\": %v", ImageColumnName, err)
		}
	case column == CommandColumnName:
		if err := l.fPrintCommand(container.Command); err != nil {
			return fmt.Errorf("could not display container's field \"%s\": %v", CommandColumnName, err)
		}
	case column == CreatedColumnName:
		if err := l.fPrintCreated(container.Created); err != nil {
			return fmt.Errorf("could not display container's field \"%s\": %v", CreatedColumnName, err)
		}
	case column == StatusColumnName:
		if err := l.fPrintStatus(container.Status); err != nil {
			return fmt.Errorf("could not display container's field \"%s\": %v", StatusColumnName, err)
		}
	case column == PortsColumnName:
		if err := l.fPrintPorts(container.Ports); err != nil {
			return fmt.Errorf("could not display container's field \"%s\": %v", PortsColumnName, err)
		}
	case column == NamesColumnName:
		if err := l.fPrintNames(container.Names); err != nil {
			return fmt.Errorf("could not display container's field \"%s\": %v", NamesColumnName, err)
		}
	case column == SizeColumnName && l.OptSize:
		if err := l.fPrintSize(container.SizeRw, container.SizeRootFs); err != nil {
			return fmt.Errorf("could not display container's field \"%s\": %v", SizeColumnName, err)
		}
	}
	return nil
}

func (l *List) fPrintID(id string) error {
	var idOutput = id
	if !l.OptNoTrunc && len(idOutput) >= IDMinWidth {
		idOutput = idOutput[:IDMinWidth]
	}
	var color = colorer.NoColor
	if l.Colorized {
		color = colorer.ColorDarkGray
	}
	return l.write([]byte(colorer.Paint(color, idOutput)))
}

func (l *List) fPrintImage(image string) error {
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
	return l.write([]byte(outputImage.String()))
}

func (l *List) fPrintCommand(command string) error {
	var commandOutput = command
	if !l.OptNoTrunc && len(command) > CommandMinWidth {
		commandOutput = fmt.Sprintf("%s…", command[:CommandMinWidth])
	}
	var color = colorer.NoColor
	if l.Colorized {
		color = colorer.ColorDarkGray
	}
	return l.write([]byte(colorer.Paintf(color, "\"%s\"", commandOutput)))
}

func (l *List) fPrintCreated(created int64) error {
	return l.write([]byte(colorer.Paint(
		l.getCreatedColor(time.Now().Unix()-created),
		humanize.Time(time.Unix(created, 0)),
	)))
}

func (l *List) fPrintStatus(status string) error {
	var statusColor = colorer.NoColor
	if l.Colorized {
		if strings.HasPrefix(status, StatusUpStr) {
			statusColor = colorer.ColorLightGreen
		} else {
			statusColor = colorer.ColorDefault
		}
	}
	return l.write([]byte(colorer.Paint(statusColor, status)))
}

func (l *List) fPrintPorts(ports []types.Port) error {
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

	return l.write([]byte(strings.Join(portsOutput, ", ")))
}

func (l *List) fPrintNames(names []string) error {
	var namesOutput = make([]string, len(names))
	for i, name := range names {
		namesOutput[i] = strings.TrimLeft(name, "/")
	}
	var color = colorer.NoColor
	if l.Colorized {
		color = colorer.ColorDefault
	}
	return l.write([]byte(colorer.Paint(
		color,
		strings.Join(namesOutput, ", "),
	)))
}

func (l *List) fPrintSize(sizeRw int64, sizeRootFs int64) error {
	var (
		colorSizeRw     = colorer.NoColor
		colorSizeRootFs = colorer.NoColor
	)
	if l.Colorized {
		colorSizeRw = l.getSizeColor(sizeRw)
		colorSizeRootFs = l.getSizeColor(sizeRw)
	}
	return l.write([]byte(fmt.Sprintf(
		"%s (%s)",
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
}
