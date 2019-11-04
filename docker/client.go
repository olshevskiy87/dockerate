package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Client struct {
	cli     *client.Client
	verbose bool
}

func NewClient(apiVer string) (*Client, error) {
	var apiOption client.Opt
	if apiVer != "" {
		apiOption = client.WithVersion(apiVer)
	} else {
		apiOption = client.WithAPIVersionNegotiation()
	}
	client, err := client.NewClientWithOpts(apiOption)
	if err != nil {
		return nil, fmt.Errorf("could not create docker client: %v", err)
	}
	return &Client{
		cli:     client,
		verbose: false,
	}, nil
}

func (c *Client) SetVerbose(verbose bool) {
	c.verbose = verbose
}

func (c *Client) GetVersion() string {
	return c.cli.ClientVersion()
}

func (c *Client) ContainerList(options types.ContainerListOptions) ([]types.Container, error) {
	return c.cli.ContainerList(
		context.Background(),
		options,
	)
}
