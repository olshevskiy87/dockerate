package docker

import (
	"github.com/docker/docker/client"
)

func GetClient(apiVer string) (*client.Client, error) {
	var apiOption client.Opt
	if apiVer != "" {
		apiOption = client.WithVersion(apiVer)
	} else {
		apiOption = client.WithAPIVersionNegotiation()
	}
	return client.NewClientWithOpts(apiOption)
}
