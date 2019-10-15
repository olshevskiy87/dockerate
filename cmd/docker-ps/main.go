package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	humanize "github.com/dustin/go-humanize"
	"github.com/reconquest/loreley"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.40"))
	if err != nil {
		fmt.Printf("could not init docker client: %v\n", err)
		os.Exit(1)
	}

	containers, err := cli.ContainerList(
		context.Background(),
		types.ContainerListOptions{All: true, Size: true},
	)
	if err != nil {
		fmt.Printf("could not get containers list: %v\n", err)
		os.Exit(1)
	}

	buffer := &bytes.Buffer{}

	const ColumnPadding = 5
	w := tabwriter.NewWriter(buffer, 0, 0, ColumnPadding, ' ', tabwriter.FilterHTML)

	_, err = w.Write([]byte("<fg 12>CONTAINER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tPORTS\tNAMES<reset>\n"))
	if err != nil {
		fmt.Printf("could not write columns header to output buffer: %v\n", err)
		os.Exit(1)
	}

	for _, container := range containers {
		var ports strings.Builder
		for _, port := range container.Ports {
			if port.PrivatePort != 0 {
				ports.WriteString(strconv.Itoa(int(port.PrivatePort)))
			}
			if port.Type != "" {
				ports.WriteString("/")
				ports.WriteString(port.Type)
			}
		}

		names := make([]string, len(container.Names))
		for i, name := range container.Names {
			names[i] = strings.TrimLeft(name, "/")
		}

		created := time.Unix(container.Created, 0)

		command := container.Command
		if len(command) > 20 {
			command = fmt.Sprintf("%s…", container.Command[:20])
		}

		_, err = w.Write(
			[]byte(
				fmt.Sprintf(
					"%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
					container.ID[:12],
					container.Image,
					fmt.Sprintf("\"%s\"", command),
					humanize.Time(created),
					container.Status,
					ports.String(),
					strings.Join(names, ", "),
				),
			),
		)
		if err != nil {
			fmt.Printf("could not write container info to output buffer: %v\n", err)
			os.Exit(1)
		}
	}
	w.Flush()

	loreley.DelimLeft = "<"
	loreley.DelimRight = ">"

	output, err := loreley.CompileAndExecuteToString(buffer.String(), nil, nil)
	if err != nil {
		fmt.Printf("could not compile result output string: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(output)
}
