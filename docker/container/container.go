package container

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/olshevskiy87/dockerate/docker"

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
	Colorized  bool
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

func (l *List) SetColorize(colorized bool) {
	l.Colorized = colorized
}

func (l *List) CompileOutput(cli *docker.Client) (string, error) {
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

	err = l.fPrintHeader(w)
	if err != nil {
		return "", fmt.Errorf("could not display header: %v", err)
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