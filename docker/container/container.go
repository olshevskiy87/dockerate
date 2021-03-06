package container

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/olshevskiy87/dockerate/docker"

	"github.com/docker/docker/api/types"
	"github.com/reconquest/loreley"
)

const (
	IDMinWidth       = 12
	ImageTagMinWidth = 12
	CommandMinWidth  = 20

	StatusUpStr = "Up"
)

type List struct {
	tabBuffer    *bytes.Buffer
	writer       *tabwriter.Writer
	OptAll       bool
	OptSize      bool
	OptQuiet     bool
	OptNoTrunc   bool
	OptNameLike  string
	OptNameILike string
	Columns      []string
	Colorized    bool
}

const (
	ContainerIDColumnName = "CONTAINER ID"
	ImageColumnName       = "IMAGE"
	CommandColumnName     = "COMMAND"
	CreatedColumnName     = "CREATED"
	StatusColumnName      = "STATUS"
	PortsColumnName       = "PORTS"
	NamesColumnName       = "NAMES"
	SizeColumnName        = "SIZE"
)

var availableColumns = []string{
	ContainerIDColumnName,
	ImageColumnName,
	CommandColumnName,
	CreatedColumnName,
	StatusColumnName,
	PortsColumnName,
	NamesColumnName,
	SizeColumnName,
}

func NewList() *List {
	list := &List{}
	list.initWriter()
	return list
}

func (l *List) initWriter() {
	l.tabBuffer = &bytes.Buffer{}
	l.writer = tabwriter.NewWriter(l.tabBuffer, 0, 0, docker.ColumnPadding, ' ', tabwriter.FilterHTML)
}

func (l *List) write(buf []byte) error {
	_, err := l.writer.Write(buf)
	return err
}

func (l *List) getBufferString(flush bool) string {
	if flush {
		l.writer.Flush()
	}
	return l.tabBuffer.String()
}

func isColumnName(checkName string) bool {
	for _, column := range availableColumns {
		if column == checkName {
			return true
		}
	}
	return false
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

func (l *List) SetOptionNameLike(name string) {
	l.OptNameLike = name
	l.OptNameILike = ""
}

func (l *List) SetOptionNameILike(name string) {
	l.OptNameILike = name
	l.OptNameLike = ""
}

func (l *List) SetColumns(columns []string) error {
	columnsMap := make(map[string]bool, len(columns))
	l.Columns = make([]string, len(columns))
	for i, column := range columns {
		column := strings.ToUpper(strings.TrimSpace(column))
		if !isColumnName(column) {
			return fmt.Errorf("wrong column name: %s", column)
		}
		if _, ok := columnsMap[column]; ok {
			return fmt.Errorf("column \"%s\" specified more than once", column)
		}
		columnsMap[column] = true
		l.Columns[i] = column
	}
	return nil
}

func (l *List) SetColorize(colorized bool) {
	l.Colorized = colorized
}

func (l *List) isColumnSet(checkName string) bool {
	for _, column := range l.Columns {
		if column == checkName {
			return true
		}
	}
	return false
}

func (l *List) precompileCheck() error {
	switch {
	case len(l.Columns) == 0:
		columns := make([]string, 0, len(availableColumns))
		for _, c := range availableColumns {
			if !l.OptSize && c == SizeColumnName {
				continue
			}
			columns = append(columns, c)
		}
		l.Columns = make([]string, len(columns))
		copy(l.Columns, columns)
	case l.OptSize && !l.isColumnSet(SizeColumnName):
		l.SetOptionSize(false)
	case !l.OptSize && l.isColumnSet(SizeColumnName):
		return fmt.Errorf("header SIZE specified without option -s")
	}
	return nil
}

func (l *List) CompileOutput(cli *docker.Client) (string, error) {
	err := l.precompileCheck()
	if err != nil {
		return "", err
	}

	containers, err := cli.ContainerList(
		types.ContainerListOptions{
			All:  l.OptAll,
			Size: l.OptSize,
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not get container list: %v", err)
	}

	err = l.fPrintHeader()
	if err != nil {
		return "", fmt.Errorf("could not display header: %v", err)
	}

CONTAINERS_LOOP:
	for _, container := range containers {
		if l.OptNameLike != "" {
			for _, name := range container.Names {
				if !strings.Contains(name, l.OptNameLike) {
					continue CONTAINERS_LOOP
				}
			}
		} else if l.OptNameILike != "" {
			for _, name := range container.Names {
				if !strings.Contains(strings.ToLower(name), strings.ToLower(l.OptNameILike)) {
					continue CONTAINERS_LOOP
				}
			}
		}
		err := l.fPrintContainer(container)
		if err != nil {
			return "", fmt.Errorf("could not display container: %v", err)
		}
	}

	loreley.DelimLeft = "<"
	loreley.DelimRight = ">"

	output, err := loreley.CompileAndExecuteToString(l.getBufferString(true), nil, nil)
	if err != nil {
		return "", fmt.Errorf("could not compile result output string: %v", err)
	}
	return output, nil
}
