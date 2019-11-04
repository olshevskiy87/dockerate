package container

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/olshevskiy87/dockerate/colorer"
	"github.com/olshevskiy87/dockerate/docker"
	"github.com/olshevskiy87/dockerate/unit"

	"github.com/docker/docker/api/types"
	"github.com/reconquest/loreley"
)

const (
	IDMinWidth       = 12
	ImageTagMinWidth = 12
	CommandMinWidth  = 20
)

type List struct {
	OptAll     bool
	OptSize    bool
	OptQuiet   bool
	OptNoTrunc bool
}

func NewList() *List {
	return &List{}
}

func (l *List) SetOptionAll(all bool) {
	l.OptAll = all
}

func (l *List) SetOptionSize(size bool) {
	l.OptSize = size
}

func (l *List) SetOptionQuiet(quiet bool) {
	l.OptQuiet = quiet
}

func (l *List) SetOptionNoTrunc(noTrunc bool) {
	l.OptNoTrunc = noTrunc
}

func (l *List) CompileOutput(cli *docker.Client, colorize bool) (string, error) {
	containers, err := cli.ContainerList(
		types.ContainerListOptions{
			All:  l.OptAll,
			Size: l.OptSize,
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not get container list: %v", err)
	}

	tabBuffer := &bytes.Buffer{}

	w := tabwriter.NewWriter(tabBuffer, 0, 0, docker.ColumnPadding, ' ', tabwriter.FilterHTML)

	if !l.OptQuiet {
		var header strings.Builder
		header.WriteString("CONTAINER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tPORTS\tNAMES")
		if l.OptSize {
			header.WriteString("\tSIZE")
		}
		header.WriteString("\n")
		_, err = colorer.Fpaintf(w, colorer.ColorLightBlue, header.String())
		if err != nil {
			return "", fmt.Errorf("could not write columns header to output buffer: %v", err)
		}
	}

	for _, container := range containers {
		err := l.fPrintContainer(w, container)
		if err != nil {
			return "", fmt.Errorf("could not display container: %v", err)
		}
	}

	w.Flush()

	loreley.DelimLeft = "<"
	loreley.DelimRight = ">"

	output, err := loreley.CompileAndExecuteToString(tabBuffer.String(), nil, nil)
	if err != nil {
		return "", fmt.Errorf("could not compile result output string: %v", err)
	}
	return output, nil
}

func (l *List) getSizeColor(size int64) colorer.ColorCode {
	if size >= 0 && size < 500*unit.Megabyte {
		return colorer.ColorDefault
	} else if size >= 500*unit.Megabyte && size < unit.Gigabyte {
		return colorer.ColorYellow
	}
	return colorer.ColorRed
}
