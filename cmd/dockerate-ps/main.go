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

	"github.com/alexflint/go-arg"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	humanize "github.com/dustin/go-humanize"
	"github.com/reconquest/loreley"
)

const (
	DockerAPIVersionDefault  = "1.40"
	ContainerIDMinWidth      = 12
	ContainerCommandMinWidth = 20
)

type argsType struct {
	All     bool   `arg:"--all,-a" help:"show all containers"`
	NoTrunc bool   `arg:"--no-trunc" help:"don't truncate output"`
	Quiet   bool   `arg:"--quiet,-q" help:"only display containers IDs"`
	Verbose bool   `arg:"--verbose,-v" help:"output more information"`
	APIVer  string `arg:"env:DOCKER_API_VER" help:"docker server API version, env DOCKER_API_VER"`
	//Size    bool   `arg:"--size,-s" help:"display containers sizes"`
}

func (argsType) Description() string {
	return "Dockerate (decorate docker commands output): List containers"
}

var version = "0.1.1"

func (argsType) Version() string {
	return fmt.Sprintf("dockerate-ps %s", version)
}

func main() {

	var args argsType
	// TODO:
	// - detect docker server API version using command "docker version --format {{.Server.APIVersion}}"
	// - ask if user wants to set system environment variable to this version if it wasn't set yet
	args.APIVer = DockerAPIVersionDefault
	arg.MustParse(&args)

	if args.Verbose {
		fmt.Printf("docker server API version: %s\n", args.APIVer)
	}

	cli, err := client.NewClientWithOpts(client.WithVersion(args.APIVer))
	if err != nil {
		fmt.Printf("could not init docker client: %v\n", err)
		os.Exit(1)
	}

	containers, err := cli.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			All: args.All,
			//Size:  args.Size,
		},
	)
	if err != nil {
		fmt.Printf("could not get containers list: %v\n", err)
		os.Exit(1)
	}

	buffer := &bytes.Buffer{}

	const ColumnPadding = 5
	w := tabwriter.NewWriter(buffer, 0, 0, ColumnPadding, ' ', tabwriter.FilterHTML)

	if !args.Quiet {
		header := "CONTAINER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tPORTS\tNAMES"
		_, err = w.Write([]byte(fmt.Sprintf("<fg 12>%s<reset>\n", header)))
		if err != nil {
			fmt.Printf("could not write columns header to output buffer: %v\n", err)
			os.Exit(1)
		}
	}

	for _, container := range containers {

		containerID := container.ID
		if !args.NoTrunc {
			containerID = containerID[:ContainerIDMinWidth]
		}

		if args.Quiet {
			_, err = w.Write([]byte(fmt.Sprintf("%s\n", containerID)))
			if err != nil {
				fmt.Printf("could not write container info to output buffer: %v\n", err)
				os.Exit(1)
			}
			continue
		}

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
		if !args.NoTrunc && len(command) > ContainerCommandMinWidth {
			command = fmt.Sprintf("%s…", container.Command[:ContainerCommandMinWidth])
		}

		_, err = w.Write(
			[]byte(
				fmt.Sprintf(
					"%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
					containerID,
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
