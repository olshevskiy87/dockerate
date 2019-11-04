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

func (l *List) fPrintContainer(w io.Writer, container types.Container) error {
	if err := l.fPrintID(w, container.ID); err != nil {
		return fmt.Errorf("could not display container's field \"ID\": %v", err)
	}
	if l.OptQuiet {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
		return nil
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
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}

func (l *List) fPrintID(w io.Writer, id string) error {
	var idOutput = id
	if !l.OptNoTrunc {
		idOutput = id[:IDMinWidth]
	}
	_, err := w.Write([]byte(colorer.Paint(
		colorer.ColorDarkGray,
		idOutput,
	)))
	return err
}

func (l *List) fPrintImage(w io.Writer, image string) error {
	var outputImage strings.Builder
	imageItems := strings.Split(image, ":")
	if len(imageItems) > 0 {
		outputImage.WriteString(colorer.Paint(
			colorer.ColorLightYellow,
			imageItems[0],
		))
	}
	if len(imageItems) > 1 {
		imageTag := imageItems[1]
		if !l.OptNoTrunc && len(imageTag) > ImageTagMinWidth {
			imageTag = fmt.Sprintf("%s…", imageTag[:ImageTagMinWidth])
		}
		outputImage.WriteString(colorer.Paintf(
			colorer.ColorLightGreen,
			":%s",
			imageTag,
		))
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
	_, err := w.Write([]byte(colorer.Paintf(
		colorer.ColorDarkGray,
		"\t\"%s\"",
		commandOutput,
	)))
	return err
}

func (l *List) fPrintCreated(w io.Writer, created int64) error {
	createdInterval := time.Now().Unix() - created
	createdColor := colorer.ColorLightGreen
	if createdInterval > unit.IntervalMonthSec {
		createdColor = colorer.ColorRed
	} else if createdInterval > unit.IntervalWeekSec {
		createdColor = colorer.ColorYellow
	}
	_, err := w.Write([]byte(colorer.Paintf(
		createdColor,
		"\t%s",
		humanize.Time(time.Unix(created, 0)),
	)))
	return err
}

func (l *List) fPrintStatus(w io.Writer, status string) error {
	statusColor := colorer.ColorDefault
	if strings.HasPrefix(status, "Up") {
		statusColor = colorer.ColorLightGreen
	}
	_, err := w.Write([]byte(colorer.Paintf(statusColor, "\t%s", status)))
	return err
}

func (l *List) fPrintPorts(w io.Writer, ports []types.Port) error {
	portsOutput := make([]string, len(ports))
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

		var portLine strings.Builder
		if hostIPPublicPort.Len() > 0 {
			portLine.WriteString(colorer.Paintf(
				colorer.ColorLightCyan,
				"%s->",
				hostIPPublicPort.String(),
			))
		}

		// PrivatePort and Type (required)
		portLine.WriteString(fmt.Sprintf(
			"%s/%s",
			strconv.Itoa(int(portInfo.PrivatePort)),
			portInfo.Type,
		))

		portsOutput[i] = portLine.String()
	}
	sort.Strings(portsOutput)
	_, err := w.Write([]byte(fmt.Sprintf(
		"\t%s",
		strings.Join(portsOutput, ", "),
	)))
	return err
}

func (l *List) fPrintNames(w io.Writer, names []string) error {
	namesOutput := make([]string, len(names))
	for i, name := range names {
		namesOutput[i] = strings.TrimLeft(name, "/")
	}
	_, err := w.Write([]byte(colorer.Paintf(
		colorer.ColorWhite,
		"\t%s",
		strings.Join(namesOutput, ", "),
	)))
	return err
}

func (l *List) fPrintSize(w io.Writer, sizeRw int64, sizeRootFs int64) error {
	_, err := w.Write([]byte(fmt.Sprintf(
		"\t%s (%s)",
		colorer.Paint(
			l.getSizeColor(sizeRw),
			humanize.Bytes(uint64(sizeRw)),
		),
		colorer.Paintf(
			l.getSizeColor(sizeRootFs),
			"virtual %s",
			humanize.Bytes(uint64(sizeRootFs)),
		),
	)))
	return err
}
