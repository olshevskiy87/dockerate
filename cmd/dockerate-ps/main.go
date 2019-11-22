package main

import (
	"fmt"
	"os"

	"github.com/olshevskiy87/dockerate/docker"
	"github.com/olshevskiy87/dockerate/docker/container"

	"github.com/alexflint/go-arg"
)

type argsType struct {
	All       bool   `arg:"--all,-a" help:"show all containers"`
	NoTrunc   bool   `arg:"--no-trunc" help:"don't truncate output"`
	Quiet     bool   `arg:"--quiet,-q" help:"only display containers IDs"`
	Size      bool   `arg:"--size,-s" help:"display containers sizes"`
	WhenColor string `arg:"--color" help:"when to use colors: always, auto, never" default:"auto"`
	APIVer    string `arg:"env:DOCKER_API_VERSION" help:"docker server API version, env DOCKER_API_VERSION"`
	Verbose   bool   `arg:"--verbose,-v" help:"output more information"`
}

func (argsType) Description() string {
	return "Dockerate (decorate docker commands output): List containers"
}

var (
	version   = "0.1.7"
	buildHash = ""
)

func (argsType) Version() string {
	var buildHashOut = ""
	if buildHash != "" {
		buildHashOut = fmt.Sprintf(" (build %s)", buildHash)
	}
	return fmt.Sprintf("dockerate-ps %s%s", version, buildHashOut)
}

func isColorModeAvailable(mode string) bool {
	for _, availableMode := range []string{"always", "auto", "never"} {
		if mode == availableMode {
			return true
		}
	}
	return false
}

func shouldBeColorized(mode string) bool {
	var isColorized = true
	if mode == "auto" {
		stdoutInfo, err := os.Stdout.Stat()
		if err != nil {
			fmt.Printf("could not get stdout info: %v\n", err)
			os.Exit(1)
		}
		// disable colors when piping
		isColorized = (stdoutInfo.Mode() & os.ModeCharDevice) != 0
	} else if mode == "never" {
		isColorized = false
	}
	return isColorized
}

func main() {
	var args argsType
	arg.MustParse(&args)

	if !isColorModeAvailable(args.WhenColor) {
		fmt.Printf("unknown color mode: %s\n", args.WhenColor)
		os.Exit(1)
	}

	cli, err := docker.NewClient(args.APIVer)
	if err != nil {
		fmt.Printf("could not init docker client: %v\n", err)
		os.Exit(1)
	}
	cli.SetVerbose(args.Verbose)

	if args.Verbose {
		fmt.Printf("docker client API version: %s\n", cli.GetVersion())
		fmt.Printf("color mode: %s\n", args.WhenColor)
	}

	list := container.NewList()

	list.SetOptionAll(args.All)
	list.SetOptionSize(args.Size)
	list.SetOptionQuiet(args.Quiet)
	list.SetOptionNoTrunc(args.NoTrunc)
	list.SetColorize(shouldBeColorized(args.WhenColor))

	output, err := list.CompileOutput(cli)
	if err != nil {
		fmt.Printf("could not display containers: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(output)
}
