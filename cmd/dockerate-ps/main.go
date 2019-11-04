package main

import (
	"fmt"
	"os"

	"github.com/olshevskiy87/dockerate/docker"
	"github.com/olshevskiy87/dockerate/docker/container"

	"github.com/alexflint/go-arg"
)

type argsType struct {
	All     bool   `arg:"--all,-a" help:"show all containers"`
	NoTrunc bool   `arg:"--no-trunc" help:"don't truncate output"`
	Quiet   bool   `arg:"--quiet,-q" help:"only display containers IDs"`
	Verbose bool   `arg:"--verbose,-v" help:"output more information"`
	Size    bool   `arg:"--size,-s" help:"display containers sizes"`
	APIVer  string `arg:"env:DOCKER_API_VERSION" help:"docker server API version, env DOCKER_API_VERSION"`
}

func (argsType) Description() string {
	return "Dockerate (decorate docker commands output): List containers"
}

var version = "0.1.5"

func (argsType) Version() string {
	return fmt.Sprintf("dockerate-ps %s", version)
}

func main() {
	var args argsType
	arg.MustParse(&args)

	cli, err := docker.NewClient(args.APIVer)
	if err != nil {
		fmt.Printf("could not init docker client: %v\n", err)
		os.Exit(1)
	}
	cli.SetVerbose(args.Verbose)

	if args.Verbose {
		fmt.Printf("docker client API version: %s\n", cli.GetVersion())
	}

	list := container.NewList()

	list.SetOptionAll(args.All)
	list.SetOptionSize(args.Size)
	list.SetOptionQuiet(args.Quiet)
	list.SetOptionNoTrunc(args.NoTrunc)

	output, err := list.CompileOutput(
		cli,
		true, // colorize
	)
	if err != nil {
		fmt.Printf("could not display containers: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(output)
}
