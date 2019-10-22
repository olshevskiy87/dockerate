package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/olshevskiy87/dockerate/colorer"

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
	ColumnPadding            = 5
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

	tabBuffer := &bytes.Buffer{}

	w := tabwriter.NewWriter(tabBuffer, 0, 0, ColumnPadding, ' ', tabwriter.FilterHTML)

	if !args.Quiet {
		_, err = colorer.Fpaintf(w, colorer.ColorLightBlue, "CONTAINER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tPORTS\tNAMES\n")
		if err != nil {
			fmt.Printf("could not write columns header to output buffer: %v\n", err)
			os.Exit(1)
		}
	}

	for _, container := range containers {

		var containerLine strings.Builder

		// 1. container ID
		containerID := container.ID
		if !args.NoTrunc {
			containerID = containerID[:ContainerIDMinWidth]
		}

		containerLine.WriteString(colorer.Paint(colorer.ColorDarkGray, containerID))

		if args.Quiet {
			containerLine.WriteString("\n")
			_, err = w.Write([]byte(containerLine.String()))
			if err != nil {
				fmt.Printf("could not write container info to output buffer: %v\n", err)
				os.Exit(1)
			}
			continue
		}

		// 2. image
		var image strings.Builder
		imageItems := strings.Split(container.Image, ":")
		if len(imageItems) > 0 {
			image.WriteString(colorer.Paint(colorer.ColorLightYellow, imageItems[0]))
		}
		if len(imageItems) > 1 {
			image.WriteString(colorer.Paintf(colorer.ColorLightGreen, ":%s", imageItems[1]))
		}
		containerLine.WriteString(fmt.Sprintf("\t%s", image.String()))

		// 3. command
		command := container.Command
		if !args.NoTrunc && len(command) > ContainerCommandMinWidth {
			command = fmt.Sprintf("%s…", container.Command[:ContainerCommandMinWidth])
		}
		containerLine.WriteString(colorer.Paintf(colorer.ColorDarkGray, "\t\"%s\"", command))

		// 4. created
		created := humanize.Time(time.Unix(container.Created, 0))
		containerLine.WriteString(colorer.Paintf(colorer.ColorGreen, "\t%s", created))

		// 5. status
		containerLine.WriteString(colorer.Paintf(colorer.ColorLightGreen, "\t%s", container.Status))

		// 6. ports
		ports := make([]string, len(container.Ports))
		for i, portInfo := range container.Ports {
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
				portLine.WriteString(colorer.Paintf(colorer.ColorLightCyan, "%s->", hostIPPublicPort.String()))
			}

			// PrivatePort and Type (required)
			portLine.WriteString(fmt.Sprintf("%s/%s", strconv.Itoa(int(portInfo.PrivatePort)), portInfo.Type))

			ports[i] = portLine.String()
		}
		sort.Strings(ports)
		containerLine.WriteString(fmt.Sprintf("\t%s", strings.Join(ports, ", ")))

		// 7. names
		names := make([]string, len(container.Names))
		for i, name := range container.Names {
			names[i] = strings.TrimLeft(name, "/")
		}
		containerLine.WriteString(colorer.Paintf(colorer.ColorWhite, "\t%s", strings.Join(names, ", ")))

		// complete result string with container info
		containerLine.WriteString("\n")
		_, err = w.Write([]byte(containerLine.String()))
		if err != nil {
			fmt.Printf("could not write container info to output buffer: %v\n", err)
			os.Exit(1)
		}
	}
	w.Flush()

	loreley.DelimLeft = "<"
	loreley.DelimRight = ">"

	output, err := loreley.CompileAndExecuteToString(tabBuffer.String(), nil, nil)
	if err != nil {
		fmt.Printf("could not compile result output string: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(output)
}
